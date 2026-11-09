package teo

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/fw"
)

// _ resource.Resource is intentionally omitted: the factory function
// NewTeoFunctionReplica110901Resource returns resource.Resource, which already
// performs the implicit compile-time check at the return statement. Only the
// extended-interface assertions are kept here.
var (
	_ resource.ResourceWithConfigure   = &TeoFunctionReplica110901Resource{}
	_ resource.ResourceWithImportState = &TeoFunctionReplica110901Resource{}
)

// NewTeoFunctionReplica110901Resource is the factory referenced by
// tencentcloud/framework/registry.go to register this resource.
func NewTeoFunctionReplica110901Resource() resource.Resource {
	return &TeoFunctionReplica110901Resource{}
}

// TeoFunctionReplica110901Resource implements resource.Resource for
// tencentcloud_teo_function_replica_1109_01.
type TeoFunctionReplica110901Resource struct {
	fw.ResourceWithConfigure
}

// TeoFunctionReplica110901ResourceModel maps the schema attributes.
type TeoFunctionReplica110901ResourceModel struct {
	Id              types.String `tfsdk:"id"`
	ZoneId          types.String `tfsdk:"zone_id"`
	FunctionId      types.String `tfsdk:"function_id"`
	ReplicaName     types.String `tfsdk:"replica_name"`
	Content         types.String `tfsdk:"content"`
	Remark          types.String `tfsdk:"remark"`
	FunctionReplica types.List   `tfsdk:"function_replicas"`
}

// teoFunctionReplica110901ItemModel maps one element of the computed
// function_replicas list.
type teoFunctionReplica110901ItemModel struct {
	FunctionId  types.String `tfsdk:"function_id"`
	ReplicaName types.String `tfsdk:"replica_name"`
	Content     types.String `tfsdk:"content"`
	Remark      types.String `tfsdk:"remark"`
	CreatedOn   types.String `tfsdk:"created_on"`
	ModifiedOn  types.String `tfsdk:"modified_on"`
}

// teoFunctionReplica110901ItemType returns the element type of the computed
// function_replicas list.
func teoFunctionReplica110901ItemType() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"function_id":  types.StringType,
			"replica_name": types.StringType,
			"content":      types.StringType,
			"remark":       types.StringType,
			"created_on":   types.StringType,
			"modified_on":  types.StringType,
		},
	}
}

func (r *TeoFunctionReplica110901Resource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "tencentcloud_teo_function_replica_1109_01"
}

func (r *TeoFunctionReplica110901Resource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a resource to create a TEO edge function replica.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Resource identifier in the form `<zone_id>#<function_id>#<replica_name>`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"zone_id": schema.StringAttribute{
				Required:    true,
				Description: "Zone ID.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"function_id": schema.StringAttribute{
				Required:    true,
				Description: "Function ID.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"replica_name": schema.StringAttribute{
				Required: true,
				Description: "Edge function replica name. Limited to 1-50 characters, allowed characters are a-z, 0-9 and -, " +
					"and - cannot be used alone or consecutively, nor at the beginning or end. Replica names must be unique " +
					"under the same FunctionId.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"content": schema.StringAttribute{
				Required:    true,
				Description: "Edge function replica content. Currently only JavaScript code is supported, up to 5MB.",
			},
			"remark": schema.StringAttribute{
				Optional:    true,
				Description: "Edge function replica description. Maximum 50 characters.",
			},
			"function_replicas": schema.ListNestedAttribute{
				Computed:    true,
				Description: "Edge function replica list.",
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"function_id": schema.StringAttribute{
							Computed:    true,
							Description: "Function ID.",
						},
						"replica_name": schema.StringAttribute{
							Computed:    true,
							Description: "Edge function replica name.",
						},
						"content": schema.StringAttribute{
							Computed:    true,
							Description: "Edge function replica content, in JavaScript code format.",
						},
						"remark": schema.StringAttribute{
							Computed:    true,
							Description: "Edge function replica description.",
						},
						"created_on": schema.StringAttribute{
							Computed:    true,
							Description: "Edge function replica creation time.",
						},
						"modified_on": schema.StringAttribute{
							Computed:    true,
							Description: "Edge function replica update time.",
						},
					},
				},
			},
		},
	}
}

// Create handles `terraform apply` for new resources.
func (r *TeoFunctionReplica110901Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan TeoFunctionReplica110901ResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	zoneId := plan.ZoneId.ValueString()
	functionId := plan.FunctionId.ValueString()
	replicaName := plan.ReplicaName.ValueString()

	service := NewTeoService(r.Client())
	if err := service.CreateTeoFunctionReplica(ctx, zoneId, functionId, replicaName,
		plan.Content.ValueString(), plan.Remark.ValueStringPointer()); err != nil {
		resp.Diagnostics.AddError(
			"Error creating TEO function replica",
			fmt.Sprintf("CreateFunctionReplica failed for zone %q function %q replica %q: %v",
				zoneId, functionId, replicaName, err),
		)
		return
	}

	plan.Id = types.StringValue(joinTeoFunctionReplica110901Id(zoneId, functionId, replicaName))

	// Persist the identity BEFORE the read-back: a partial failure here keeps
	// the id in state so `terraform destroy` can clean the replica up.
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	item, err := service.DescribeTeoFunctionReplica(ctx, zoneId, functionId, replicaName)
	if err != nil {
		resp.Diagnostics.AddError("Error reading TEO function replica after create", err.Error())
		return
	}
	if item == nil {
		resp.Diagnostics.AddError(
			"TEO function replica not readable after create",
			fmt.Sprintf("replica %q was created but cannot be read back", replicaName),
		)
		return
	}
	resp.Diagnostics.Append(r.refreshModel(ctx, &plan, item)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Read handles `terraform refresh` and post-apply state refresh.
func (r *TeoFunctionReplica110901Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state TeoFunctionReplica110901ResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	zoneId, functionId, replicaName, ok := SplitTeoFunctionReplica110901Id(state.Id.ValueString())
	if !ok {
		resp.Diagnostics.AddError(
			"Invalid resource ID",
			fmt.Sprintf("TEO function replica id is broken, id is %q (expected `<zone_id>#<function_id>#<replica_name>`)",
				state.Id.ValueString()),
		)
		return
	}

	service := NewTeoService(r.Client())
	item, err := service.DescribeTeoFunctionReplica(ctx, zoneId, functionId, replicaName)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading TEO function replica",
			fmt.Sprintf("DescribeFunctionReplicas failed for zone %q function %q replica %q: %v",
				zoneId, functionId, replicaName, err),
		)
		return
	}

	// Queried successfully but nothing matched => confirmed not found on the
	// cloud side; drop it from state so Terraform plans a recreation.
	if item == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	// zone_id is not echoed back by DescribeFunctionReplicas; keep the value
	// from state (documented as limited drift detection in the .md file).
	state.ZoneId = types.StringValue(zoneId)
	state.FunctionId = types.StringValue(functionId)
	state.ReplicaName = types.StringValue(replicaName)
	resp.Diagnostics.Append(r.refreshModel(ctx, &state, item)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update handles in-place modifications.
func (r *TeoFunctionReplica110901Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan TeoFunctionReplica110901ResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state TeoFunctionReplica110901ResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	zoneId, functionId, replicaName, ok := SplitTeoFunctionReplica110901Id(state.Id.ValueString())
	if !ok {
		resp.Diagnostics.AddError(
			"Invalid resource ID",
			fmt.Sprintf("TEO function replica id is broken, id is %q", state.Id.ValueString()),
		)
		return
	}

	// content is Required so the plan always carries a value; send it as-is.
	//
	// remark is Optional and its clear semantics is not documented: whether
	// ModifyFunctionReplica treats an omitted Remark as "keep the old value" or
	// "clear it" is unknown. Fail loudly on the ambiguous transition (value in
	// state, null in plan) instead of guessing a direction: guessing "keep"
	// leaves a permanent diff on every later plan, guessing "clear" silently
	// destroys the description. An explicit empty string is passed through and
	// validated by the cloud API.
	var remark *string
	switch {
	case plan.Remark.IsNull():
		if !state.Remark.IsNull() {
			resp.Diagnostics.AddError(
				"Clear semantics of remark is not confirmed",
				"Removing the remark attribute from the configuration requires knowing whether "+
					"ModifyFunctionReplica treats an omitted Remark as \"keep\" or \"clear\"; this is not documented. "+
					"Keep a non-empty remark value, or set an explicit empty string to ask the cloud API to clear it.",
			)
			return
		}
	default:
		remark = plan.Remark.ValueStringPointer()
	}

	service := NewTeoService(r.Client())
	if err := service.ModifyTeoFunctionReplica(ctx, zoneId, functionId, replicaName,
		plan.Content.ValueString(), remark); err != nil {
		resp.Diagnostics.AddError(
			"Error updating TEO function replica",
			fmt.Sprintf("ModifyFunctionReplica failed for zone %q function %q replica %q: %v",
				zoneId, functionId, replicaName, err),
		)
		return
	}

	item, err := service.DescribeTeoFunctionReplica(ctx, zoneId, functionId, replicaName)
	if err != nil {
		resp.Diagnostics.AddError("Error reading TEO function replica after update", err.Error())
		return
	}
	if item == nil {
		resp.Diagnostics.AddError(
			"TEO function replica not readable after update",
			fmt.Sprintf("replica %q was updated but cannot be read back", replicaName),
		)
		return
	}
	plan.Id = state.Id
	resp.Diagnostics.Append(r.refreshModel(ctx, &plan, item)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete handles resource deletion.
func (r *TeoFunctionReplica110901Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state TeoFunctionReplica110901ResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	zoneId, functionId, replicaName, ok := SplitTeoFunctionReplica110901Id(state.Id.ValueString())
	if !ok {
		resp.Diagnostics.AddError(
			"Invalid resource ID",
			fmt.Sprintf("TEO function replica id is broken, id is %q", state.Id.ValueString()),
		)
		return
	}

	service := NewTeoService(r.Client())
	if err := service.DeleteTeoFunctionReplica(ctx, zoneId, functionId, replicaName); err != nil {
		resp.Diagnostics.AddError(
			"Error deleting TEO function replica",
			fmt.Sprintf("DeleteFunctionReplica failed for zone %q function %q replica %q: %v",
				zoneId, functionId, replicaName, err),
		)
		return
	}
	// After a successful delete there is no need to explicitly clear state:
	// the Framework removes it automatically.
}

// ImportState supports `terraform import` using the
// `<zone_id>#<function_id>#<replica_name>` form. Composite identity cannot be
// passed through as a plain id attribute: the segment count is validated here
// and each segment is written back to its own attribute.
func (r *TeoFunctionReplica110901Resource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	zoneId, functionId, replicaName, ok := SplitTeoFunctionReplica110901Id(req.ID)
	if !ok {
		resp.Diagnostics.AddError(
			"Invalid import ID",
			fmt.Sprintf("Expected import ID in the form `<zone_id>%s<function_id>%s<replica_name>`, got %q",
				tccommon.FILED_SP, tccommon.FILED_SP, req.ID),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("zone_id"), zoneId)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("function_id"), functionId)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("replica_name"), replicaName)...)
}

// refreshModel copies cloud-owned fields into the model and rebuilds the
// computed function_replicas list from the Describe response. The identity
// fields keep the values from state (zone_id is not echoed back by the API;
// documented as limited drift detection).
func (r *TeoFunctionReplica110901Resource) refreshModel(ctx context.Context, m *TeoFunctionReplica110901ResourceModel, item *teov20220901.FunctionReplica) diag.Diagnostics {
	var diags diag.Diagnostics
	m.Content = types.StringPointerValue(item.Content)
	m.Remark = types.StringPointerValue(item.Remark)
	list, listDiags := types.ListValueFrom(ctx, teoFunctionReplica110901ItemType(), []teoFunctionReplica110901ItemModel{{
		FunctionId:  types.StringPointerValue(item.FunctionId),
		ReplicaName: types.StringPointerValue(item.ReplicaName),
		Content:     types.StringPointerValue(item.Content),
		Remark:      types.StringPointerValue(item.Remark),
		CreatedOn:   types.StringPointerValue(item.CreatedOn),
		ModifiedOn:  types.StringPointerValue(item.ModifiedOn),
	}})
	diags.Append(listDiags...)
	if diags.HasError() {
		return diags
	}
	m.FunctionReplica = list
	return diags
}

// joinTeoFunctionReplica110901Id builds the composite resource id.
func joinTeoFunctionReplica110901Id(zoneId, functionId, replicaName string) string {
	return strings.Join([]string{zoneId, functionId, replicaName}, tccommon.FILED_SP)
}

// SplitTeoFunctionReplica110901Id parses `<zone_id>#<function_id>#<replica_name>`.
// It is exported for use by the external test package.
func SplitTeoFunctionReplica110901Id(id string) (string, string, string, bool) {
	parts := strings.Split(id, tccommon.FILED_SP)
	if len(parts) != 3 || parts[0] == "" || parts[1] == "" || parts[2] == "" {
		return "", "", "", false
	}
	return parts[0], parts[1], parts[2], true
}
