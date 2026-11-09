package teo

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	teo "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/fw"
)

// Compile-time interface assertions (repo convention: only assert the extension
// interfaces; the base interface is checked implicitly by the factory return value).
var (
	_ resource.ResourceWithConfigure   = &TeoAliasDomain110904Resource{}
	_ resource.ResourceWithImportState = &TeoAliasDomain110904Resource{}
)

// NewTeoAliasDomain110904Resource is the factory referenced by
// tencentcloud/framework/registry.go to register this resource.
func NewTeoAliasDomain110904Resource() resource.Resource {
	return &TeoAliasDomain110904Resource{}
}

// TeoAliasDomain110904Resource implements resource.Resource for
// tencentcloud_teo_alias_domain_1109_04.
type TeoAliasDomain110904Resource struct {
	fw.ResourceWithConfigure
}

// TeoAliasDomain110904ResourceModel maps the schema attributes.
//
// CertType and CertId are write-only on the cloud side: DescribeAliasDomains
// does not return them, so Read keeps the values already stored in state
// instead of overwriting them (drift detection for these two attributes is
// limited; documented in resource_tc_teo_alias_domain_1109_04.md).
type TeoAliasDomain110904ResourceModel struct {
	Id         types.String `tfsdk:"id"`
	ZoneId     types.String `tfsdk:"zone_id"`
	AliasName  types.String `tfsdk:"alias_name"`
	TargetName types.String `tfsdk:"target_name"`
	CertType   types.String `tfsdk:"cert_type"`
	CertId     types.Set    `tfsdk:"cert_id"`
	Status     types.String `tfsdk:"status"`
}

func (r *TeoAliasDomain110904Resource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "tencentcloud_teo_alias_domain_1109_04"
}

func (r *TeoAliasDomain110904Resource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a resource to create a TEO alias domain. " +
			"The alias domain maps to a target domain under the same zone; deletion only removes the alias domain itself " +
			"and does not affect the target domain or the zone.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Resource identifier in the form `<zone_id>#<alias_name>`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"zone_id": schema.StringAttribute{
				Required:    true,
				Description: "Site ID.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"alias_name": schema.StringAttribute{
				Required:    true,
				Description: "Alias domain name.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"target_name": schema.StringAttribute{
				Required:    true,
				Description: "Target domain name.",
			},
			"cert_type": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Description: "Certificate configuration. Valid values: `none` (no certificate), `hosting` (SSL hosted certificate). " +
					"The cloud side defaults to `none` when omitted. The field is not returned by the query API, so the configured " +
					"value is kept in state (limited drift detection).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"cert_id": schema.SetAttribute{
				Optional:    true,
				ElementType: types.StringType,
				Description: "Certificate IDs. Required when `cert_type` is `hosting`. " +
					"The field is not returned by the query API, so the configured value is kept in state (limited drift detection).",
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"status": schema.StringAttribute{
				Computed: true,
				Description: "Alias domain status. Valid values: `active` (effective), `pending` (deploying), " +
					"`conflict` (reclaimed), `stop` (disabled).",
			},
		},
	}
}

// Create handles `terraform apply` for new resources.
func (r *TeoAliasDomain110904Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan TeoAliasDomain110904ResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	zoneId := plan.ZoneId.ValueString()
	aliasName := plan.AliasName.ValueString()

	certType, certId, diags := teoAliasDomain110904CertParams(ctx, plan.CertType, plan.CertId)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	service := TeoService{client: r.Client()}
	err := service.CreateTeoAliasDomain110904(ctx, zoneId, aliasName, plan.TargetName.ValueString(), certType, certId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating TEO alias domain",
			fmt.Sprintf("CreateAliasDomain failed for zone %s alias %s: %v", zoneId, aliasName, err),
		)
		return
	}

	// Persist the identity BEFORE any asynchronous wait: if the wait fails, the
	// resource is still in state and `terraform destroy` can clean it up.
	plan.Id = types.StringValue(strings.Join([]string{zoneId, aliasName}, tccommon.FILED_SP))
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Creation is asynchronous on the cloud side: poll until the alias domain
	// becomes active. Poll only on the expected "not visible yet" outcome.
	if err := service.WaitTeoAliasDomainActive110904(ctx, zoneId, aliasName); err != nil {
		resp.Diagnostics.AddError(
			"Error waiting for TEO alias domain to become active",
			fmt.Sprintf("alias domain %s (zone %s) was created but did not become active: %v", aliasName, zoneId, err),
		)
		return
	}

	item, err := service.DescribeTeoAliasDomainById(ctx, zoneId, aliasName)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading TEO alias domain after create",
			fmt.Sprintf("DescribeAliasDomains failed for zone %s alias %s: %v", zoneId, aliasName, err),
		)
		return
	}
	if item == nil {
		resp.Diagnostics.AddError(
			"TEO alias domain not readable after create",
			fmt.Sprintf("alias domain %s (zone %s) was created but cannot be read back", aliasName, zoneId),
		)
		return
	}
	// cert_type is Optional+Computed: when omitted from the configuration the
	// plan value stays unknown at apply time. The query API does not return
	// CertType, so resolve it to null (not configured) here; a final state value
	// may never stay unknown.
	if plan.CertType.IsUnknown() {
		plan.CertType = types.StringNull()
	}
	refreshTeoAliasDomain110904Model(&plan, item)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Read handles `terraform refresh` and post-apply state refresh.
func (r *TeoAliasDomain110904Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state TeoAliasDomain110904ResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	zoneId, aliasName, ok := splitTeoAliasDomain110904Id(state.Id.ValueString())
	if !ok {
		resp.Diagnostics.AddError(
			"Invalid resource ID",
			fmt.Sprintf("TEO alias domain id is broken, id is %q (expected `<zone_id>#<alias_name>`)", state.Id.ValueString()),
		)
		return
	}

	service := TeoService{client: r.Client()}
	item, err := service.DescribeTeoAliasDomainById(ctx, zoneId, aliasName)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading TEO alias domain",
			fmt.Sprintf("DescribeAliasDomains failed for zone %s alias %s: %v", zoneId, aliasName, err),
		)
		return
	}
	// Queried successfully but nothing matched => the alias domain is confirmed
	// gone on the cloud side; drop it from state so Terraform plans a recreation.
	if item == nil {
		resp.Diagnostics.AddWarning(
			"Resource not found",
			fmt.Sprintf("TEO alias domain %s (zone %s) not found, removing from state.", aliasName, zoneId),
		)
		resp.State.RemoveResource(ctx)
		return
	}

	state.ZoneId = types.StringValue(zoneId)
	state.AliasName = types.StringValue(aliasName)
	refreshTeoAliasDomain110904Model(&state, item)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update handles in-place modifications.
func (r *TeoAliasDomain110904Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan TeoAliasDomain110904ResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state TeoAliasDomain110904ResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	zoneId, aliasName, ok := splitTeoAliasDomain110904Id(state.Id.ValueString())
	if !ok {
		resp.Diagnostics.AddError(
			"Invalid resource ID",
			fmt.Sprintf("TEO alias domain id is broken, id is %q", state.Id.ValueString()),
		)
		return
	}

	// Cert contract (confirmed): ModifyAliasDomain keeps the previous
	// configuration when CertType is omitted; clearing the certificate requires
	// explicitly setting cert_type = "none". cert_type is Optional+Computed, so
	// removing it from the configuration leaves the plan value unknown (or, when
	// a state value exists, copied from state via UseStateForUnknown); it never
	// becomes null on update. The diff against the current state is therefore the
	// only reliable "user changed the certificate configuration" signal:
	//   - plan unknown or plan == state -> send nothing, the cloud keeps the
	//     previous configuration (which already matches state);
	//   - plan != state -> send the new values from the plan (each field only
	//     when it is set), including the explicit cert_type = "none" clear path.
	var certId []*string
	certChanged := !plan.CertType.IsUnknown() && !plan.CertType.Equal(state.CertType)
	if certChanged || (!plan.CertId.IsUnknown() && !plan.CertId.Equal(state.CertId)) {
		certChanged = true
	}

	certType := ""
	if certChanged {
		if !plan.CertType.IsNull() && !plan.CertType.IsUnknown() {
			certType = plan.CertType.ValueString()
		}
		var diags diag.Diagnostics
		certId, diags = teoAliasDomain110904CertIdList(ctx, plan.CertId)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	service := TeoService{client: r.Client()}
	if err := service.ModifyTeoAliasDomain110904(ctx, zoneId, aliasName, plan.TargetName.ValueString(),
		certChanged, certType, certId); err != nil {
		resp.Diagnostics.AddError(
			"Error updating TEO alias domain",
			fmt.Sprintf("ModifyAliasDomain failed for zone %s alias %s: %v", zoneId, aliasName, err),
		)
		return
	}

	item, err := service.DescribeTeoAliasDomainById(ctx, zoneId, aliasName)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading TEO alias domain after update",
			fmt.Sprintf("DescribeAliasDomains failed for zone %s alias %s: %v", zoneId, aliasName, err),
		)
		return
	}
	if item == nil {
		resp.Diagnostics.AddError(
			"TEO alias domain not readable after update",
			fmt.Sprintf("alias domain %s (zone %s) was updated but cannot be read back", aliasName, zoneId),
		)
		return
	}
	plan.Id = state.Id
	// cert_type is Optional+Computed and the query API does not return CertType:
	// resolve an unknown plan value to the prior state value (or null when it was
	// never configured); a final state value may never stay unknown.
	if plan.CertType.IsUnknown() {
		plan.CertType = state.CertType
	}
	if plan.CertId.IsUnknown() {
		plan.CertId = state.CertId
	}
	refreshTeoAliasDomain110904Model(&plan, item)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete handles resource deletion. Only the alias domain itself is removed;
// the target domain and the zone are not affected.
func (r *TeoAliasDomain110904Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state TeoAliasDomain110904ResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	zoneId, aliasName, ok := splitTeoAliasDomain110904Id(state.Id.ValueString())
	if !ok {
		resp.Diagnostics.AddError(
			"Invalid resource ID",
			fmt.Sprintf("TEO alias domain id is broken, id is %q", state.Id.ValueString()),
		)
		return
	}

	service := TeoService{client: r.Client()}
	// Confirm the alias domain is gone before declaring the delete finished;
	// only a confirmed absence (item == nil) short-circuits the delete call.
	item, err := service.DescribeTeoAliasDomainById(ctx, zoneId, aliasName)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading TEO alias domain before delete",
			fmt.Sprintf("DescribeAliasDomains failed for zone %s alias %s: %v", zoneId, aliasName, err),
		)
		return
	}
	if item == nil {
		return
	}

	if err := service.DeleteTeoAliasDomain110904(ctx, zoneId, aliasName); err != nil {
		resp.Diagnostics.AddError(
			"Error deleting TEO alias domain",
			fmt.Sprintf("DeleteAliasDomain failed for zone %s alias %s: %v", zoneId, aliasName, err),
		)
		return
	}

	// Deletion is asynchronous on the cloud side: poll until the alias domain
	// is no longer returned by the query API.
	if err := service.WaitTeoAliasDomainDeleted110904(ctx, zoneId, aliasName); err != nil {
		resp.Diagnostics.AddError(
			"Error waiting for TEO alias domain to be deleted",
			fmt.Sprintf("alias domain %s (zone %s) delete was accepted but the domain is still visible: %v", aliasName, zoneId, err),
		)
		return
	}
}

// ImportState supports `terraform import` using the `<zone_id>#<alias_name>` form.
func (r *TeoAliasDomain110904Resource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if _, _, ok := splitTeoAliasDomain110904Id(req.ID); !ok {
		resp.Diagnostics.AddError(
			"Invalid import ID",
			fmt.Sprintf("Expected import ID in the form `<zone_id>%s<alias_name>`, got %q", tccommon.FILED_SP, req.ID),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
}

// refreshTeoAliasDomain110904Model copies cloud-owned fields into the model.
//
// Only fields that DescribeAliasDomains actually returns are refreshed.
// CertType/CertId are not part of the response: the user-supplied values in
// the model are left untouched to avoid wiping them into null and producing a
// permanent diff (drift detection for them is limited by the cloud API).
func refreshTeoAliasDomain110904Model(m *TeoAliasDomain110904ResourceModel, item *teo.AliasDomain) {
	if item == nil {
		return
	}
	if item.TargetName != nil {
		m.TargetName = types.StringPointerValue(item.TargetName)
	}
	m.Status = types.StringPointerValue(item.Status)
}

// splitTeoAliasDomain110904Id parses `<zone_id>#<alias_name>`.
func splitTeoAliasDomain110904Id(id string) (string, string, bool) {
	parts := strings.Split(id, tccommon.FILED_SP)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", false
	}
	return parts[0], parts[1], true
}

// teoAliasDomain110904CertParams converts the framework cert_type/cert_id
// attributes into the plain Go values consumed by the service layer. An empty
// certType means "not configured" (nothing is sent and the cloud side applies
// its own default).
func teoAliasDomain110904CertParams(ctx context.Context, certType types.String, certId types.Set) (string, []*string, diag.Diagnostics) {
	var diags diag.Diagnostics
	typ := ""
	if !certType.IsNull() && !certType.IsUnknown() {
		typ = certType.ValueString()
	}
	if typ == "" {
		return "", nil, diags
	}
	ids, idDiags := teoAliasDomain110904CertIdList(ctx, certId)
	diags.Append(idDiags...)
	return typ, ids, diags
}

// teoAliasDomain110904CertIdList converts the cert_id set into []*string.
// It returns nil when the set is null/unknown (nothing to send).
func teoAliasDomain110904CertIdList(ctx context.Context, certId types.Set) ([]*string, diag.Diagnostics) {
	var diags diag.Diagnostics
	if certId.IsNull() || certId.IsUnknown() {
		return nil, diags
	}
	var elements []types.String
	diags.Append(certId.ElementsAs(ctx, &elements, false)...)
	if diags.HasError() {
		return nil, diags
	}
	ids := make([]*string, 0, len(elements))
	for _, element := range elements {
		ids = append(ids, element.ValueStringPointer())
	}
	return ids, diags
}
