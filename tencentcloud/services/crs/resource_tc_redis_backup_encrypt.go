package crs

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	sdkResource "github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	redis "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/redis/v20180412"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/sharedmeta"
)

var (
	_ resource.Resource                = &CrsRedisBackupEncryptConfigResource{}
	_ resource.ResourceWithConfigure   = &CrsRedisBackupEncryptConfigResource{}
	_ resource.ResourceWithImportState = &CrsRedisBackupEncryptConfigResource{}
)

// NewCrsRedisBackupEncryptConfigResource is the factory referenced by
// tencentcloud/framework/registry.go to register this resource.
func NewCrsRedisBackupEncryptConfigResource() resource.Resource {
	return &CrsRedisBackupEncryptConfigResource{}
}

// CrsRedisBackupEncryptConfigResource implements resource.Resource for
// tencentcloud_redis_backup_encrypt. It manages the backup encryption
// configuration of an existing Redis instance and never touches the lifecycle
// of the host instance itself.
type CrsRedisBackupEncryptConfigResource struct {
	client RedisService
}

// CrsRedisBackupEncryptConfigResourceModel maps the schema attributes.
// Id equals the host instance ID.
type CrsRedisBackupEncryptConfigResourceModel struct {
	Id             types.String `tfsdk:"id"`
	InstanceId     types.String `tfsdk:"instance_id"`
	BackupEncrypt  types.Bool   `tfsdk:"backup_encrypt"`
	EncryptKeyType types.Int64  `tfsdk:"encrypt_key_type"`
	KmsKeyId       types.String `tfsdk:"kms_key_id"`
	KmsRegion      types.String `tfsdk:"kms_region"`
}

func (r *CrsRedisBackupEncryptConfigResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "tencentcloud_redis_backup_encrypt"
}

func (r *CrsRedisBackupEncryptConfigResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a resource to manage the backup encryption configuration of an existing TencentDB for Redis instance. " +
			"This resource never creates or deletes the host instance: destroying it disables backup encryption on the cloud side " +
			"(`backup_encrypt = false`), the Redis instance itself is left untouched.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "ID of the backup encryption configuration, same as the host instance ID.",
				PlanModifiers: []planmodifier.String{
					stringUseStateForUnknown{},
				},
			},
			"instance_id": schema.StringAttribute{
				Required:    true,
				Description: "ID of the Redis instance whose backup encryption is managed. Changing it forces a new resource.",
				PlanModifiers: []planmodifier.String{
					stringRequiresReplace{},
				},
			},
			"backup_encrypt": schema.BoolAttribute{
				Required: true,
				Description: "Whether to enable backup encryption. `true`: backups are encrypted; `false`: backups are not encrypted. " +
					"Destroying this resource sets it to `false`.",
			},
			"encrypt_key_type": schema.Int64Attribute{
				Optional: true,
				Computed: true,
				Description: "Encryption key type. `0`: Tencent Cloud managed key (cloud default); `1`: user-defined KMS key. " +
					"Only meaningful when `backup_encrypt` is `true`.",
				Validators: []validator.Int64{
					int64OneOf{allowed: []int64{0, 1}},
				},
			},
			"kms_key_id": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Description: "ID of the KMS key, format `xxxxxxx-xxxx-xxxx-xxxx-xxxxxxx`. " +
					"Only meaningful when `encrypt_key_type` is `1`.",
			},
			"kms_region": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Description: "Region of the KMS key, format `ap-xxxx`. " +
					"Only meaningful when `encrypt_key_type` is `1`.",
			},
		},
	}
}

// Configure receives the *sharedmeta.ProviderMeta populated by
// framework.Provider.Configure and constructs the SDKv2-style RedisService
// from the shared *connectivity.TencentCloudClient.
func (r *CrsRedisBackupEncryptConfigResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	meta, ok := req.ProviderData.(*sharedmeta.ProviderMeta)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected ProviderData type",
			fmt.Sprintf("Expected *sharedmeta.ProviderMeta, got: %T. "+
				"This is a bug in the provider; please report it.", req.ProviderData),
		)
		return
	}

	r.client = NewRedisService(meta.Client)
}

// Create enables or configures backup encryption for the first time.
func (r *CrsRedisBackupEncryptConfigResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	defer tccommon.LogElapsed("resource.tencentcloud_redis_backup_encrypt.create")()

	var plan CrsRedisBackupEncryptConfigResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	instanceId := plan.InstanceId.ValueString()
	if instanceId == "" {
		resp.Diagnostics.AddError(
			"Invalid instance_id",
			"redis_backup_encrypt requires a non-empty instance_id",
		)
		return
	}

	if err := r.modifyBackupEncrypt(ctx, &plan); err != nil {
		resp.Diagnostics.AddError(
			"Error setting Redis backup encryption",
			fmt.Sprintf("ModifyBackupEncrypt failed for redis_backup_encrypt instance %s: %v", instanceId, err),
		)
		return
	}

	plan.Id = plan.InstanceId
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Refresh so that the state mirrors the cloud after the synchronous modify.
	r.refreshState(ctx, &plan, &resp.State, &resp.Diagnostics)
}

// Read reflects the cloud-side backup encryption configuration.
func (r *CrsRedisBackupEncryptConfigResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	defer tccommon.LogElapsed("resource.tencentcloud_redis_backup_encrypt.read")()

	var state CrsRedisBackupEncryptConfigResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// After import only `id` is known; backfill the schema attribute from it.
	if state.InstanceId.IsNull() || state.InstanceId.ValueString() == "" {
		state.InstanceId = state.Id
	}

	r.refreshState(ctx, &state, &resp.State, &resp.Diagnostics)
}

// Update applies configuration changes through the same modify API used by Create.
func (r *CrsRedisBackupEncryptConfigResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	defer tccommon.LogElapsed("resource.tencentcloud_redis_backup_encrypt.update")()

	var plan CrsRedisBackupEncryptConfigResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state CrsRedisBackupEncryptConfigResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// instance_id is RequiresReplace, so it cannot change here; keep the state value as anchor.
	if plan.InstanceId.IsNull() || plan.InstanceId.ValueString() == "" {
		plan.InstanceId = state.InstanceId
	}

	if err := r.modifyBackupEncrypt(ctx, &plan); err != nil {
		resp.Diagnostics.AddError(
			"Error updating Redis backup encryption",
			fmt.Sprintf("ModifyBackupEncrypt failed for redis_backup_encrypt instance %s: %v", plan.InstanceId.ValueString(), err),
		)
		return
	}

	plan.Id = plan.InstanceId
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	r.refreshState(ctx, &plan, &resp.State, &resp.Diagnostics)
}

// Delete disables backup encryption. The host Redis instance is never deleted.
func (r *CrsRedisBackupEncryptConfigResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	defer tccommon.LogElapsed("resource.tencentcloud_redis_backup_encrypt.delete")()

	var state CrsRedisBackupEncryptConfigResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	instanceId := state.InstanceId.ValueString()
	if instanceId == "" {
		instanceId = state.Id.ValueString()
	}
	if instanceId == "" {
		resp.Diagnostics.AddError(
			"Invalid resource state",
			"redis_backup_encrypt has neither instance_id nor id in state, cannot disable backup encryption",
		)
		return
	}

	log.Printf("[DEBUG] redis_backup_encrypt id=%s, destroying resource disables backup encryption", instanceId)

	// Disabling encryption is idempotent on the cloud side: it is issued even when
	// the current state already has backup_encrypt = false.
	var err error
	retryErr := sdkResource.Retry(tccommon.WriteRetryTimeout, func() *sdkResource.RetryError {
		err = r.client.ModifyBackupEncrypt(ctx, instanceId, false, nil, nil, nil)
		if err != nil {
			return tccommon.RetryError(err)
		}
		return nil
	})
	if retryErr != nil {
		resp.Diagnostics.AddError(
			"Error disabling Redis backup encryption",
			fmt.Sprintf("ModifyBackupEncrypt failed for redis_backup_encrypt instance %s: %v", instanceId, retryErr),
		)
		return
	}
}

// ImportState supports `terraform import` using the host instance ID.
func (r *CrsRedisBackupEncryptConfigResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if req.ID == "" {
		resp.Diagnostics.AddError(
			"Invalid import ID",
			"Expected the host Redis instance ID for redis_backup_encrypt, got an empty string",
		)
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
}

// modifyBackupEncrypt is shared by Create and Update — both call the same modify API.
// Only attributes explicitly configured by the user are forwarded.
func (r *CrsRedisBackupEncryptConfigResource) modifyBackupEncrypt(ctx context.Context, m *CrsRedisBackupEncryptConfigResourceModel) error {
	instanceId := m.InstanceId.ValueString()
	backupEncrypt := m.BackupEncrypt.ValueBool()

	var encryptKeyType *int64
	if !m.EncryptKeyType.IsNull() && !m.EncryptKeyType.IsUnknown() {
		v := m.EncryptKeyType.ValueInt64()
		encryptKeyType = &v
	}

	var kmsKeyId *string
	if !m.KmsKeyId.IsNull() && !m.KmsKeyId.IsUnknown() {
		v := m.KmsKeyId.ValueString()
		kmsKeyId = &v
	}

	var kmsRegion *string
	if !m.KmsRegion.IsNull() && !m.KmsRegion.IsUnknown() {
		v := m.KmsRegion.ValueString()
		kmsRegion = &v
	}

	var err error
	retryErr := sdkResource.Retry(tccommon.WriteRetryTimeout, func() *sdkResource.RetryError {
		err = r.client.ModifyBackupEncrypt(ctx, instanceId, backupEncrypt, encryptKeyType, kmsKeyId, kmsRegion)
		if err != nil {
			return tccommon.RetryError(err)
		}
		return nil
	})
	if retryErr != nil {
		return retryErr
	}

	return nil
}

// refreshState queries the cloud configuration and writes it into the model and state.
// The state id is only cleared when the host instance no longer exists.
func (r *CrsRedisBackupEncryptConfigResource) refreshState(
	ctx context.Context,
	m *CrsRedisBackupEncryptConfigResourceModel,
	state *tfsdk.State,
	diagnostics *diag.Diagnostics,
) {
	instanceId := m.InstanceId.ValueString()
	if instanceId == "" {
		instanceId = m.Id.ValueString()
	}
	if instanceId == "" {
		log.Printf("[CRUD] redis_backup_encrypt id=%s", m.Id.ValueString())
		diagnostics.AddError(
			"Invalid resource state",
			"redis_backup_encrypt has neither instance_id nor id in state",
		)
		return
	}

	var config *redis.DescribeBackupEncryptConfigResponseParams
	var err error
	retryErr := sdkResource.Retry(tccommon.ReadRetryTimeout, func() *sdkResource.RetryError {
		config, err = r.client.DescribeBackupEncryptConfig(ctx, instanceId)
		if err != nil {
			if isRedisInstanceNotFound(err) {
				return sdkResource.NonRetryableError(err)
			}
			return tccommon.RetryError(err)
		}
		return nil
	})

	if retryErr != nil {
		if isRedisInstanceNotFound(retryErr) {
			log.Printf("[CRUD] redis_backup_encrypt id=%s, instance not exists, remove from state", instanceId)
			state.RemoveResource(ctx)
			return
		}
		diagnostics.AddError(
			"Error reading Redis backup encryption",
			fmt.Sprintf("DescribeBackupEncryptConfig failed for redis_backup_encrypt instance %s: %v", instanceId, retryErr),
		)
		return
	}

	if config == nil {
		log.Printf("[CRUD] redis_backup_encrypt id=%s", instanceId)
		diagnostics.AddError(
			"Error reading Redis backup encryption",
			fmt.Sprintf("DescribeBackupEncryptConfig returned an empty response for redis_backup_encrypt instance %s", instanceId),
		)
		return
	}

	applyBackupEncryptConfig(m, instanceId, config)

	diagnostics.Append(state.Set(ctx, m)...)
}

// applyBackupEncryptConfig maps the cloud response onto the resource model.
// A nil SDK pointer becomes a null attribute instead of a zero value.
func applyBackupEncryptConfig(m *CrsRedisBackupEncryptConfigResourceModel, instanceId string, config *redis.DescribeBackupEncryptConfigResponseParams) {
	if config == nil {
		return
	}

	if config.BackupEncrypt != nil {
		m.BackupEncrypt = types.BoolPointerValue(config.BackupEncrypt)
	}
	m.EncryptKeyType = types.Int64PointerValue(config.EncryptKeyType)
	m.KmsKeyId = types.StringPointerValue(config.KmsKeyId)
	m.KmsRegion = types.StringPointerValue(config.KmsRegion)
	m.InstanceId = types.StringValue(instanceId)
	m.Id = types.StringValue(instanceId)
}

// isRedisInstanceNotFound reports whether the error means the host Redis instance is gone.
func isRedisInstanceNotFound(err error) bool {
	if err == nil {
		return false
	}
	var sdkErr *sdkErrors.TencentCloudSDKError
	if !errors.As(err, &sdkErr) {
		return false
	}
	return sdkErr.Code == RedisInstanceNotFound
}

// stringRequiresReplace forces resource replacement when a string attribute changes.
type stringRequiresReplace struct{}

func (stringRequiresReplace) Description(_ context.Context) string {
	return "If the value of this attribute changes, Terraform will destroy and recreate the resource."
}

func (stringRequiresReplace) MarkdownDescription(ctx context.Context) string {
	return stringRequiresReplace{}.Description(ctx)
}

func (stringRequiresReplace) PlanModifyString(_ context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	if req.StateValue.IsNull() {
		// Resource is being created; nothing to replace.
		return
	}
	if req.PlanValue.IsNull() {
		// Resource is being destroyed; nothing to replace.
		return
	}
	if !req.PlanValue.Equal(req.StateValue) {
		resp.RequiresReplace = true
	}
}

// stringUseStateForUnknown copies the prior state value into the plan when the new
// plan value is unknown — used for the computed `id` attribute so it stays stable
// across refreshes.
type stringUseStateForUnknown struct{}

func (stringUseStateForUnknown) Description(_ context.Context) string {
	return "Once set, the value of this attribute in state will not change."
}

func (stringUseStateForUnknown) MarkdownDescription(ctx context.Context) string {
	return stringUseStateForUnknown{}.Description(ctx)
}

func (stringUseStateForUnknown) PlanModifyString(_ context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	if req.StateValue.IsNull() {
		return
	}
	if !req.PlanValue.IsUnknown() {
		return
	}
	resp.PlanValue = req.StateValue
}

// int64OneOf validates that an int64 attribute is one of the allowed values.
type int64OneOf struct {
	allowed []int64
}

func (v int64OneOf) Description(_ context.Context) string {
	return fmt.Sprintf("value must be one of %v", v.allowed)
}

func (v int64OneOf) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v int64OneOf) ValidateInt64(_ context.Context, req validator.Int64Request, resp *validator.Int64Response) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	value := req.ConfigValue.ValueInt64()
	for _, allowed := range v.allowed {
		if value == allowed {
			return
		}
	}

	resp.Diagnostics.AddAttributeError(
		req.Path,
		"Invalid Attribute Value",
		fmt.Sprintf("Attribute %s %s, got: %d", req.Path, v.Description(context.Background()), value),
	)
}
