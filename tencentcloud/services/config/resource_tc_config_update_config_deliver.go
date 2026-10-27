package config

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	configv20220802 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/config/v20220802"

	tccommonhelper "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/sharedmeta"
)

// configUpdateConfigDeliverId is the fixed sentinel state ID of this
// singleton resource. The Config delivery setting has no cloud-side
// identifier, so a stable constant is used instead of a synthetic token;
// it doubles as the only accepted `terraform import` ID.
const configUpdateConfigDeliverId = "config_update_config_deliver"

var (
	_ resource.Resource                = &ConfigUpdateConfigDeliverResource{}
	_ resource.ResourceWithConfigure   = &ConfigUpdateConfigDeliverResource{}
	_ resource.ResourceWithImportState = &ConfigUpdateConfigDeliverResource{}
)

// NewConfigUpdateConfigDeliverResource is the factory referenced by
// tencentcloud/framework/registry.go to register this resource.
func NewConfigUpdateConfigDeliverResource() resource.Resource {
	return &ConfigUpdateConfigDeliverResource{}
}

// ConfigUpdateConfigDeliverResource implements resource.Resource for
// tencentcloud_config_update_config_deliver.
type ConfigUpdateConfigDeliverResource struct {
	client ConfigService
}

// ConfigUpdateConfigDeliverResourceModel maps the schema attributes onto the
// DescribeConfigDeliver / UpdateConfigDeliver payloads.
type ConfigUpdateConfigDeliverResourceModel struct {
	Id                 types.String `tfsdk:"id"`
	Status             types.Int64  `tfsdk:"status"`
	DeliverName        types.String `tfsdk:"deliver_name"`
	TargetArn          types.String `tfsdk:"target_arn"`
	DeliverPrefix      types.String `tfsdk:"deliver_prefix"`
	DeliverType        types.String `tfsdk:"deliver_type"`
	DeliverContentType types.Int64  `tfsdk:"deliver_content_type"`
	CreateTime         types.String `tfsdk:"create_time"`
}

func (r *ConfigUpdateConfigDeliverResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "tencentcloud_config_update_config_deliver"
}

func (r *ConfigUpdateConfigDeliverResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a resource to manage the Tencent Cloud Config log delivery setting. " +
			"The delivery setting is a global singleton per account and region: it always exists and cannot be " +
			"created or deleted, only read and updated, so this resource adopts the existing setting in place on " +
			"first apply. Exactly one instance of this resource should be declared per provider configuration. " +
			"Destroying this resource disables Config log delivery by setting status to 0 on the cloud side: " +
			"it never deletes the Config service, nor any COS bucket or CLS topic referenced by target_arn.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "ID of the resource. It is the fixed value `config_update_config_deliver` because the delivery setting is a singleton.",
				PlanModifiers: []planmodifier.String{
					stringUseStateForUnknownModifier{},
				},
			},
			"status": schema.Int64Attribute{
				Required: true,
				Description: "Delivery switch of Config log delivery. `0` means disabled, `1` means enabled. " +
					"This argument is required because the switch decides whether audit logs are delivered, " +
					"so it must be declared explicitly rather than inherited from cloud state.",
			},
			"deliver_name": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Name of the delivery service.",
			},
			"target_arn": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Description: "Six-segment resource ARN of the delivery target. For COS: " +
					"`qcs::cos:$region:$account:prefix/$appid/$BucketName`. For CLS: `qcs::cls:$region:$account:cls/topicId`.",
			},
			"deliver_prefix": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Log prefix of the delivered files.",
			},
			"deliver_type": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Delivery type. Valid values: `COS`, `CLS`.",
				Validators: []validator.String{
					stringOneOfValidator{allowed: []string{"COS", "CLS"}},
				},
			},
			"deliver_content_type": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "Content type of the delivered logs. Valid values: `1` (configuration change), `2` (resource list), `3` (all).",
				Validators: []validator.Int64{
					int64BetweenValidator{min: 1, max: 3},
				},
			},
			"create_time": schema.StringAttribute{
				Computed:    true,
				Description: "Creation time of the delivery setting.",
			},
		},
	}
}

// Configure receives the *sharedmeta.ProviderMeta populated by
// framework.Provider.Configure and builds the SDKv2-style ConfigService from
// the shared *connectivity.TencentCloudClient.
func (r *ConfigUpdateConfigDeliverResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = NewConfigService(meta.Client)
}

// Create handles `terraform apply`. The delivery setting has no create API,
// so Create performs the first UpdateConfigDeliver call and adopts whatever
// setting the account already holds.
func (r *ConfigUpdateConfigDeliverResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ConfigUpdateConfigDeliverResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.modifyConfig(ctx, &plan); err != nil {
		resp.Diagnostics.AddError(
			"Error creating config_update_config_deliver",
			fmt.Sprintf("UpdateConfigDeliver failed for config_update_config_deliver: %v", err),
		)
		return
	}

	plan.Id = types.StringValue(configUpdateConfigDeliverId)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Read handles `terraform refresh` and the post-apply state refresh.
func (r *ConfigUpdateConfigDeliverResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ConfigUpdateConfigDeliverResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.Id.ValueString()
	if id == "" {
		id = configUpdateConfigDeliverId
	}

	result, err := r.client.DescribeConfigDeliver(ctx)
	if err != nil {
		log.Printf("[CRITICAL] describe config_update_config_deliver failed, id=%s, reason=%s", id, err.Error())
		resp.Diagnostics.AddError(
			"Error reading config_update_config_deliver",
			fmt.Sprintf("DescribeConfigDeliver failed for config_update_config_deliver (id=%s): %v", id, err),
		)
		return
	}

	if result == nil {
		log.Printf("[CRITICAL] describe config_update_config_deliver returned an empty body, id=%s", id)
		resp.Diagnostics.AddError(
			"Empty response from DescribeConfigDeliver",
			fmt.Sprintf("DescribeConfigDeliver returned an empty response body for config_update_config_deliver (id=%s); "+
				"this is an API anomaly, so the resource is kept in state", id),
		)
		return
	}

	state.Id = types.StringValue(id)
	flattenConfigDeliver(result, &state)

	log.Printf("[DEBUG] describe config_update_config_deliver success, id=%s", id)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update handles in-place modifications. Every attribute accepted by
// UpdateConfigDeliver is mutable in place, so no replacement ever occurs.
func (r *ConfigUpdateConfigDeliverResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan ConfigUpdateConfigDeliverResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state ConfigUpdateConfigDeliverResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.modifyConfig(ctx, &plan); err != nil {
		resp.Diagnostics.AddError(
			"Error updating config_update_config_deliver",
			fmt.Sprintf("UpdateConfigDeliver failed for config_update_config_deliver: %v", err),
		)
		return
	}

	plan.Id = state.Id
	if plan.Id.IsNull() || plan.Id.ValueString() == "" {
		plan.Id = types.StringValue(configUpdateConfigDeliverId)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete handles `terraform destroy`. This resource is an enable switch, so
// Delete turns delivery off instead of deleting anything: leaving delivery
// enabled after destroy would silently keep shipping audit logs.
func (r *ConfigUpdateConfigDeliverResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state ConfigUpdateConfigDeliverResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.Id.ValueString()

	request := configv20220802.NewUpdateConfigDeliverRequest()
	request.Status = tccommonhelper.Uint64(0)

	if err := r.client.UpdateConfigDeliver(ctx, request); err != nil {
		log.Printf("[CRITICAL] disable config_update_config_deliver failed, id=%s, reason=%s", id, err.Error())
		resp.Diagnostics.AddError(
			"Error deleting config_update_config_deliver",
			fmt.Sprintf("UpdateConfigDeliver with Status=0 failed for config_update_config_deliver (id=%s): %v", id, err),
		)
		return
	}
}

// ImportState supports `terraform import`, which must use the sentinel ID
// because the delivery setting is a singleton. The following Read hydrates
// every attribute from the cloud response.
func (r *ConfigUpdateConfigDeliverResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if req.ID != configUpdateConfigDeliverId {
		resp.Diagnostics.AddError(
			"Invalid import ID",
			fmt.Sprintf("%q is the only supported import ID for tencentcloud_config_update_config_deliver, got %q",
				configUpdateConfigDeliverId, req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
}

// modifyConfig sends one UpdateConfigDeliver call built from the planned
// model. It is shared by Create and Update. Status is always sent because it
// is Required; every other field is sent only when the plan carries a known,
// non-null value.
func (r *ConfigUpdateConfigDeliverResource) modifyConfig(ctx context.Context, plan *ConfigUpdateConfigDeliverResourceModel) error {
	// status is Required, so an unknown or null value here would silently be
	// coerced to 0 by ValueInt64() and turn log delivery off. Reject it
	// explicitly instead of sending a wrong switch state to the cloud.
	if plan.Status.IsNull() || plan.Status.IsUnknown() {
		return fmt.Errorf("config_update_config_deliver: status must be known, got a null or unknown value")
	}

	request := configv20220802.NewUpdateConfigDeliverRequest()
	request.Status = tccommonhelper.Uint64(uint64(plan.Status.ValueInt64()))

	if !plan.DeliverName.IsNull() && !plan.DeliverName.IsUnknown() {
		request.DeliverName = tccommonhelper.String(plan.DeliverName.ValueString())
	}

	if !plan.TargetArn.IsNull() && !plan.TargetArn.IsUnknown() {
		request.TargetArn = tccommonhelper.String(plan.TargetArn.ValueString())
	}

	if !plan.DeliverPrefix.IsNull() && !plan.DeliverPrefix.IsUnknown() {
		request.DeliverPrefix = tccommonhelper.String(plan.DeliverPrefix.ValueString())
	}

	if !plan.DeliverType.IsNull() && !plan.DeliverType.IsUnknown() {
		request.DeliverType = tccommonhelper.String(plan.DeliverType.ValueString())
	}

	if !plan.DeliverContentType.IsNull() && !plan.DeliverContentType.IsUnknown() {
		request.DeliverContentType = tccommonhelper.Uint64(uint64(plan.DeliverContentType.ValueInt64()))
	}

	if err := r.client.UpdateConfigDeliver(ctx, request); err != nil {
		return fmt.Errorf("UpdateConfigDeliver failed: %v", err)
	}

	return nil
}

// flattenConfigDeliver maps a DescribeConfigDeliver payload onto the model.
// Nil response pointers become null state values so that Optional+Computed
// attributes do not produce spurious diffs, and never the Go zero value.
func flattenConfigDeliver(src *configv20220802.DescribeConfigDeliverResponseParams, dst *ConfigUpdateConfigDeliverResourceModel) {
	dst.Status = uint64ToInt64Value(src.Status)
	dst.DeliverName = types.StringPointerValue(src.DeliverName)
	dst.TargetArn = types.StringPointerValue(src.TargetArn)
	dst.DeliverPrefix = types.StringPointerValue(src.DeliverPrefix)
	dst.DeliverType = types.StringPointerValue(src.DeliverType)
	dst.DeliverContentType = uint64ToInt64Value(src.DeliverContentType)
	dst.CreateTime = types.StringPointerValue(src.CreateTime)
}

// uint64ToInt64Value converts the *uint64 fields used by the Config SDK to
// types.Int64. A nil pointer becomes types.Int64Null() so that an absent
// field does not produce a spurious diff against an Optional+Computed
// attribute. The framework-side helper package is not importable from the
// services subtree because it lives under tencentcloud/framework/internal.
func uint64ToInt64Value(p *uint64) types.Int64 {
	if p == nil {
		return types.Int64Null()
	}
	return types.Int64Value(int64(*p))
}

// stringUseStateForUnknownModifier copies the prior state value into the plan
// when the planned value is unknown. It keeps the computed `id` stable across
// refreshes, mirroring planmodifier.UseStateForUnknown.
type stringUseStateForUnknownModifier struct{}

func (stringUseStateForUnknownModifier) Description(_ context.Context) string {
	return "once set, the value of this attribute in state will not change"
}

func (m stringUseStateForUnknownModifier) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}

func (stringUseStateForUnknownModifier) PlanModifyString(_ context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	if req.StateValue.IsNull() {
		return
	}
	if !req.PlanValue.IsUnknown() {
		return
	}
	resp.PlanValue = req.StateValue
}

// stringOneOfValidator rejects values outside the documented enum, so that
// an invalid value fails at plan time instead of surfacing as an API error.
type stringOneOfValidator struct {
	allowed []string
}

func (v stringOneOfValidator) Description(_ context.Context) string {
	return fmt.Sprintf("value must be one of %v", v.allowed)
}

func (v stringOneOfValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v stringOneOfValidator) ValidateString(_ context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	value := req.ConfigValue.ValueString()
	for _, allowed := range v.allowed {
		if value == allowed {
			return
		}
	}

	resp.Diagnostics.Append(diag.NewAttributeErrorDiagnostic(
		req.Path,
		"Invalid attribute value",
		fmt.Sprintf("expected deliver_type to be one of %v, got %q", v.allowed, value),
	))
}

// int64BetweenValidator rejects values outside the inclusive [min, max] range.
type int64BetweenValidator struct {
	min int64
	max int64
}

func (v int64BetweenValidator) Description(_ context.Context) string {
	return fmt.Sprintf("value must be between %d and %d inclusive", v.min, v.max)
}

func (v int64BetweenValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v int64BetweenValidator) ValidateInt64(_ context.Context, req validator.Int64Request, resp *validator.Int64Response) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	value := req.ConfigValue.ValueInt64()
	if value >= v.min && value <= v.max {
		return
	}

	resp.Diagnostics.Append(diag.NewAttributeErrorDiagnostic(
		req.Path,
		"Invalid attribute value",
		fmt.Sprintf("expected deliver_content_type to be between %d and %d inclusive, got %d", v.min, v.max, value),
	))
}
