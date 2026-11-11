package teo

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	teosdk "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/fw"
	fwhelper "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

// Compile-time interface assertions (repo convention: only assert the
// extension interfaces; the base interface is checked implicitly by the
// factory return value).
var (
	_ resource.ResourceWithConfigure   = &TeoFunction111101Resource{}
	_ resource.ResourceWithImportState = &TeoFunction111101Resource{}
)

// errTeoFunctionNotFound marks an empty DescribeFunctions response filtered by
// the exact function id. The case must be told apart from a transport or API
// failure before state may be dropped on refresh.
var errTeoFunctionNotFound = errors.New("teo function not found")

// NewTeoFunction111101Resource is the factory referenced by
// tencentcloud/framework/registry.go to register this resource.
func NewTeoFunction111101Resource() resource.Resource {
	return &TeoFunction111101Resource{}
}

// TeoFunction111101Resource implements resource.Resource for
// tencentcloud_teo_function_1111_01, managing the lifecycle of one TEO edge
// function (CreateFunction / DescribeFunctions / ModifyFunction /
// DeleteFunction).
//
// The cloud side addresses a function with a single FunctionId, so the
// Terraform id is the FunctionId itself. Because DescribeFunctions requires
// the owning ZoneId, the import id is `<zone_id>#<function_id>` and both
// parts are stored back into state before the follow-up Read runs.
type TeoFunction111101Resource struct {
	fw.ResourceWithConfigure
}

// TeoFunction111101ResourceModel maps the top-level schema attributes. Every
// schema attribute appears here exactly once.
type TeoFunction111101ResourceModel struct {
	Id                           types.String `tfsdk:"id"`
	ZoneId                       types.String `tfsdk:"zone_id"`
	Name                         types.String `tfsdk:"name"`
	Content                      types.String `tfsdk:"content"`
	Remark                       types.String `tfsdk:"remark"`
	FunctionId                   types.String `tfsdk:"function_id"`
	Domain                       types.String `tfsdk:"domain"`
	DomainComplianceRestrictions types.List   `tfsdk:"domain_compliance_restrictions"`
	CreateTime                   types.String `tfsdk:"create_time"`
	UpdateTime                   types.String `tfsdk:"update_time"`
}

// TeoFunction111101RestrictionModel maps one element of
// domain_compliance_restrictions.
type TeoFunction111101RestrictionModel struct {
	Reason types.String `tfsdk:"reason"`
	Region types.String `tfsdk:"region"`
}

func (r *TeoFunction111101Resource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "tencentcloud_teo_function_1111_01"
}

func (r *TeoFunction111101Resource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a resource to create and manage an edge function of TEO (Tencent Cloud EdgeOne).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Resource identifier. It is the function id.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"zone_id": schema.StringAttribute{
				Required:    true,
				Description: "ID of the zone the function belongs to. Immutable: changing it forces resource replacement.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Required: true,
				Description: "Function name. It can only contain lowercase letters, digits and hyphens, must start and end " +
					"with a letter or digit, and supports up to 30 characters. Immutable: ModifyFunction cannot change it, " +
					"so changing it forces resource replacement.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"content": schema.StringAttribute{
				Required:    true,
				Description: "Function content. Only JavaScript code is supported, up to 5MB.",
			},
			"remark": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Function description, up to 60 characters. The cloud keeps the existing remark when the " +
					"field is not sent to ModifyFunction, so removing it from configuration does not clear the remote value.",
			},
			"function_id": schema.StringAttribute{
				Computed:    true,
				Description: "Function ID.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"domain": schema.StringAttribute{
				Computed:    true,
				Description: "Default domain of the function.",
			},
			"domain_compliance_restrictions": schema.ListNestedAttribute{
				Computed:    true,
				Description: "List of region access restrictions on the default domain caused by compliance issues.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"reason": schema.StringAttribute{
							Computed: true,
							Description: "Reason why the access restriction is enforced. " +
								"Valid values: `ICP_RECORD_REQUIRED`, `GOVERNMENT_ORDER`.",
						},
						"region": schema.StringAttribute{
							Computed:    true,
							Description: "Country/region code of the restricted access area, following the ISO 3166 standard.",
						},
					},
				},
			},
			"create_time": schema.StringAttribute{
				Computed:    true,
				Description: "Creation time in UTC, following the ISO 8601 standard.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"update_time": schema.StringAttribute{
				Computed:    true,
				Description: "Last modification time in UTC, following the ISO 8601 standard.",
			},
		},
	}
}

// Create handles `terraform apply` for new resources.
func (r *TeoFunction111101Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan TeoFunction111101ResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, ok := r.client(&resp.Diagnostics)
	if !ok {
		return
	}

	request := teosdk.NewCreateFunctionRequest()
	request.ZoneId = fwhelper.StringPointerFromValue(plan.ZoneId)
	request.Name = fwhelper.StringPointerFromValue(plan.Name)
	request.Content = fwhelper.StringPointerFromValue(plan.Content)
	request.Remark = fwhelper.StringPointerFromValue(plan.Remark)

	response, err := client.UseTeoClient().CreateFunctionWithContext(ctx, request)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating TEO function",
			fmt.Sprintf("CreateFunction failed (zone_id=%q, name=%q): %s",
				plan.ZoneId.ValueString(), plan.Name.ValueString(), err),
		)
		return
	}
	if response == nil || response.Response == nil || response.Response.FunctionId == nil || *response.Response.FunctionId == "" {
		resp.Diagnostics.AddError(
			"Error creating TEO function",
			fmt.Sprintf("CreateFunction returned no function id (zone_id=%q, name=%q).",
				plan.ZoneId.ValueString(), plan.Name.ValueString()),
		)
		return
	}

	// The reliable identity is available, so persist it before any follow-up
	// read so a failure cannot leave the resource unmanaged.
	plan.FunctionId = types.StringValue(*response.Response.FunctionId)
	plan.Id = types.StringValue(*response.Response.FunctionId)

	fn, err := describeTeoFunction111101(ctx, client, plan.ZoneId.ValueString(), plan.FunctionId.ValueString())
	if err != nil {
		if errors.Is(err, errTeoFunctionNotFound) {
			resp.Diagnostics.AddError(
				"Error reading TEO function after create",
				fmt.Sprintf("CreateFunction succeeded (function_id=%q) but DescribeFunctions cannot find it yet. "+
					"Run `terraform import tencentcloud_teo_function_1111_01.<name> %s%s%s` to keep it managed.",
					plan.FunctionId.ValueString(), plan.ZoneId.ValueString(), tccommon.FILED_SP, plan.FunctionId.ValueString()),
			)
			return
		}
		resp.Diagnostics.AddError(
			"Error reading TEO function after create",
			fmt.Sprintf("CreateFunction succeeded (function_id=%q) but DescribeFunctions failed: %s",
				plan.FunctionId.ValueString(), err),
		)
		return
	}

	resp.Diagnostics.Append(refreshTeoFunction111101Model(ctx, &plan, fn)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Read handles `terraform refresh` and post-apply state refresh.
func (r *TeoFunction111101Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state TeoFunction111101ResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if state.Id.IsNull() || state.ZoneId.IsNull() {
		resp.Diagnostics.AddError(
			"Invalid resource state",
			fmt.Sprintf("TEO function state is missing id or zone_id (id=%q, zone_id=%q).",
				state.Id.ValueString(), state.ZoneId.ValueString()),
		)
		return
	}

	client, ok := r.client(&resp.Diagnostics)
	if !ok {
		return
	}

	fn, err := describeTeoFunction111101(ctx, client, state.ZoneId.ValueString(), state.Id.ValueString())
	if err != nil {
		// A confirmed "not found" means the function was removed outside
		// Terraform, so it is dropped from state instead of failing refresh.
		if errors.Is(err, errTeoFunctionNotFound) {
			resp.Diagnostics.AddWarning(
				"Resource not found",
				fmt.Sprintf("TEO function %s no longer exists, removing it from state.", state.Id.ValueString()),
			)
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error reading TEO function",
			fmt.Sprintf("DescribeFunctions failed (zone_id=%q, function_id=%q): %s",
				state.ZoneId.ValueString(), state.Id.ValueString(), err),
		)
		return
	}

	resp.Diagnostics.Append(refreshTeoFunction111101Model(ctx, &state, fn)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update handles in-place updates. Only `content` and `remark` are writable
// through ModifyFunction; zone_id and name are guarded by RequiresReplace.
// The ModifyFunction API keeps the existing remark when the field is omitted,
// so a null planned remark is skipped instead of being sent as an empty value.
func (r *TeoFunction111101Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan TeoFunction111101ResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state TeoFunction111101ResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if state.Id.IsNull() || state.ZoneId.IsNull() {
		resp.Diagnostics.AddError(
			"Invalid resource state",
			fmt.Sprintf("TEO function state is missing id or zone_id (id=%q, zone_id=%q).",
				state.Id.ValueString(), state.ZoneId.ValueString()),
		)
		return
	}

	client, ok := r.client(&resp.Diagnostics)
	if !ok {
		return
	}

	request := teosdk.NewModifyFunctionRequest()
	request.ZoneId = helper.String(state.ZoneId.ValueString())
	request.FunctionId = helper.String(state.Id.ValueString())
	request.Content = fwhelper.StringPointerFromValue(plan.Content)
	if !plan.Remark.IsNull() && !plan.Remark.IsUnknown() {
		request.Remark = fwhelper.StringPointerFromValue(plan.Remark)
	}

	if _, err := client.UseTeoClient().ModifyFunctionWithContext(ctx, request); err != nil {
		resp.Diagnostics.AddError(
			"Error updating TEO function",
			fmt.Sprintf("ModifyFunction failed (zone_id=%q, function_id=%q): %s",
				state.ZoneId.ValueString(), state.Id.ValueString(), err),
		)
		return
	}

	fn, err := describeTeoFunction111101(ctx, client, state.ZoneId.ValueString(), state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading TEO function after update",
			fmt.Sprintf("ModifyFunction succeeded (function_id=%q) but DescribeFunctions failed: %s",
				state.Id.ValueString(), err),
		)
		return
	}

	resp.Diagnostics.Append(refreshTeoFunction111101Model(ctx, &plan, fn)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete handles resource deletion.
func (r *TeoFunction111101Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state TeoFunction111101ResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if state.Id.IsNull() || state.ZoneId.IsNull() {
		resp.Diagnostics.AddError(
			"Invalid resource state",
			fmt.Sprintf("TEO function state is missing id or zone_id (id=%q, zone_id=%q).",
				state.Id.ValueString(), state.ZoneId.ValueString()),
		)
		return
	}

	client, ok := r.client(&resp.Diagnostics)
	if !ok {
		return
	}

	request := teosdk.NewDeleteFunctionRequest()
	request.ZoneId = helper.String(state.ZoneId.ValueString())
	request.FunctionId = helper.String(state.Id.ValueString())

	if _, err := client.UseTeoClient().DeleteFunctionWithContext(ctx, request); err != nil {
		// A function that is already gone is a successful deletion.
		if tccommon.IsExpectError(err, []string{"ResourceNotFound", "ResourceInexistent"}) {
			return
		}
		resp.Diagnostics.AddError(
			"Error deleting TEO function",
			fmt.Sprintf("DeleteFunction failed (zone_id=%q, function_id=%q): %s",
				state.ZoneId.ValueString(), state.Id.ValueString(), err),
		)
	}
}

// ImportState supports `terraform import` with the `<zone_id>#<function_id>`
// form. Both parts are written into state so the Read that follows can query
// DescribeFunctions, which requires the owning zone id.
func (r *TeoFunction111101Resource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	zoneId, functionId, ok := splitTeoFunction111101Id(req.ID)
	if !ok {
		resp.Diagnostics.AddError(
			"Invalid import ID",
			fmt.Sprintf("Import id %q is invalid, expected `<zone_id>%s<function_id>`.", req.ID, tccommon.FILED_SP),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), functionId)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("zone_id"), zoneId)...)
}

// client returns the shared TencentCloud SDK client. When the provider has not
// been configured it records a diagnostic and reports false, so the CRUD
// methods never dereference a nil SDK client.
func (r *TeoFunction111101Resource) client(diags *diag.Diagnostics) (*connectivity.TencentCloudClient, bool) {
	if r.Client() == nil {
		diags.AddError(
			"Provider not configured",
			"The provider client is not available. This is unexpected; please report it.",
		)
		return nil, false
	}
	return r.Client(), true
}

// describeTeoFunction111101 reads one function by its exact id within a zone.
// An empty result is reported as errTeoFunctionNotFound; API and transport
// failures are returned verbatim.
func describeTeoFunction111101(ctx context.Context, client *connectivity.TencentCloudClient, zoneId, functionId string) (*teosdk.Function, error) {
	request := teosdk.NewDescribeFunctionsRequest()
	request.ZoneId = helper.String(zoneId)
	request.FunctionIds = []*string{helper.String(functionId)}

	response, err := client.UseTeoClient().DescribeFunctionsWithContext(ctx, request)
	if err != nil {
		return nil, err
	}
	if response == nil || response.Response == nil {
		return nil, errTeoFunctionNotFound
	}
	for _, fn := range response.Response.Functions {
		if fn != nil && fn.FunctionId != nil && *fn.FunctionId == functionId {
			return fn, nil
		}
	}
	return nil, errTeoFunctionNotFound
}

// refreshTeoFunction111101Model copies the cloud state into the model.
// Attributes the Describe API may omit keep their previous plan/state value,
// so a value configured by the practitioner is never nulled out after apply.
func refreshTeoFunction111101Model(ctx context.Context, m *TeoFunction111101ResourceModel, fn *teosdk.Function) diag.Diagnostics {
	var diags diag.Diagnostics

	if fn == nil {
		diags.AddError(
			"Malformed TEO function response",
			"DescribeFunctions returned an empty function entry.",
		)
		return diags
	}

	m.ZoneId = fwhelper.StringValueOrPrev(fn.ZoneId, m.ZoneId)
	m.FunctionId = fwhelper.StringValueOrPrev(fn.FunctionId, m.FunctionId)
	m.Name = fwhelper.StringValueOrPrev(fn.Name, m.Name)
	m.Content = fwhelper.StringValueOrPrev(fn.Content, m.Content)
	m.Remark = fwhelper.StringValueOrPrev(fn.Remark, m.Remark)
	m.Domain = types.StringPointerValue(fn.Domain)
	m.CreateTime = types.StringPointerValue(fn.CreateTime)
	m.UpdateTime = types.StringPointerValue(fn.UpdateTime)

	if !m.FunctionId.IsNull() && !m.FunctionId.IsUnknown() {
		m.Id = m.FunctionId
	}

	restrictions, restrictionDiags := flattenTeoFunction111101Restrictions(ctx, fn.DomainComplianceRestrictions)
	diags.Append(restrictionDiags...)
	m.DomainComplianceRestrictions = restrictions

	// Required attributes must never be written as null, otherwise Terraform
	// rejects the refreshed object.
	if !diags.HasError() && (m.ZoneId.IsNull() || m.Name.IsNull() || m.Content.IsNull()) {
		diags.AddError(
			"Malformed TEO function response",
			fmt.Sprintf("DescribeFunctions returned no zone_id, name or content for function %q.", m.Id.ValueString()),
		)
	}
	return diags
}

// flattenTeoFunction111101Restrictions builds the domain_compliance_restrictions
// list from the Describe response.
func flattenTeoFunction111101Restrictions(ctx context.Context, restrictions []*teosdk.ComplianceRestriction) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics
	objType := teoFunction111101RestrictionObjectType()

	out := make([]TeoFunction111101RestrictionModel, 0, len(restrictions))
	for _, restriction := range restrictions {
		if restriction == nil {
			continue // a malformed response must never panic the provider
		}
		out = append(out, TeoFunction111101RestrictionModel{
			Reason: types.StringPointerValue(restriction.Reason),
			Region: types.StringPointerValue(restriction.Region),
		})
	}

	list, listDiags := types.ListValueFrom(ctx, objType, out)
	diags.Append(listDiags...)
	if diags.HasError() {
		return types.ListNull(objType), diags
	}
	return list, diags
}

// teoFunction111101RestrictionObjectType is the object type of the
// domain_compliance_restrictions elements.
func teoFunction111101RestrictionObjectType() attr.Type {
	return types.ObjectType{AttrTypes: map[string]attr.Type{
		"reason": types.StringType,
		"region": types.StringType,
	}}
}

// splitTeoFunction111101Id parses `<zone_id>#<function_id>`.
func splitTeoFunction111101Id(id string) (string, string, bool) {
	parts := strings.Split(id, tccommon.FILED_SP)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", false
	}
	return parts[0], parts[1], true
}
