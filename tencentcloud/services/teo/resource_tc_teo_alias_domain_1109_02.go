package teo

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	teo "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/fw"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

// _ resource.Resource is intentionally omitted: the factory function
// NewTeoAliasDomain110902Resource returns resource.Resource, which already
// performs the implicit compile-time check at the return statement. Only the
// extended-interface assertions are kept here.
var (
	_ resource.ResourceWithConfigure   = &TeoAliasDomain110902Resource{}
	_ resource.ResourceWithImportState = &TeoAliasDomain110902Resource{}
)

// NewTeoAliasDomain110902Resource is the factory referenced by
// tencentcloud/framework/registry.go to register this resource.
func NewTeoAliasDomain110902Resource() resource.Resource {
	return &TeoAliasDomain110902Resource{}
}

// TeoAliasDomain110902Resource implements resource.Resource for
// tencentcloud_teo_alias_domain_1109_02.
type TeoAliasDomain110902Resource struct {
	fw.ResourceWithConfigure
}

// TeoAliasDomainItem is the read-only detail of one alias domain returned by
// DescribeAliasDomains.
type TeoAliasDomainItem struct {
	AliasName  types.String `tfsdk:"alias_name"`
	ZoneId     types.String `tfsdk:"zone_id"`
	TargetName types.String `tfsdk:"target_name"`
	Status     types.String `tfsdk:"status"`
	ForbidMode types.Int64  `tfsdk:"forbid_mode"`
	CreatedOn  types.String `tfsdk:"created_on"`
	ModifiedOn types.String `tfsdk:"modified_on"`
}

// TeoAliasDomain110902ResourceModel maps the schema attributes.
type TeoAliasDomain110902ResourceModel struct {
	Id           types.String         `tfsdk:"id"`
	ZoneId       types.String         `tfsdk:"zone_id"`
	AliasName    types.String         `tfsdk:"alias_name"`
	TargetName   types.String         `tfsdk:"target_name"`
	CertType     types.String         `tfsdk:"cert_type"`
	CertId       types.Set            `tfsdk:"cert_id"`
	AliasDomains []TeoAliasDomainItem `tfsdk:"alias_domains"`
}

func (r *TeoAliasDomain110902Resource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "tencentcloud_teo_alias_domain_1109_02"
}

func (r *TeoAliasDomain110902Resource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a resource to create a TEO alias domain. " +
			"The alias domain resolves to the target domain and shares its acceleration configuration.",
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
				Description: "Target domain name that the alias domain points to.",
			},
			"cert_type": schema.StringAttribute{
				Optional: true,
				Description: "Certificate configuration. Valid values: `none` (no configuration), `hosting` (SSL hosted certificate). " +
					"Default value: `none`. On update, omitting this attribute keeps the previous configuration.",
			},
			"cert_id": schema.SetAttribute{
				Optional:    true,
				ElementType: types.StringType,
				Description: "Certificate ID list, required when `cert_type` is `hosting`. " +
					"DescribeAliasDomains does not return this field, so its value is kept from the configuration and drift detection for it is not possible.",
			},
			"alias_domains": schema.ListAttribute{
				Computed:    true,
				Description: "Alias domain details. Contains exactly one element describing this alias domain.",
				ElementType: aliasDomain110902ObjectType(),
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

// Create handles `terraform apply` for new resources.
func (r *TeoAliasDomain110902Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan TeoAliasDomain110902ResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if plan.ZoneId.IsNull() || plan.ZoneId.ValueString() == "" {
		resp.Diagnostics.AddError("Missing zone_id", "zone_id is required and cannot be empty.")
		return
	}
	if plan.AliasName.IsNull() || plan.AliasName.ValueString() == "" {
		resp.Diagnostics.AddError("Missing alias_name", "alias_name is required and cannot be empty.")
		return
	}
	if plan.TargetName.IsNull() || plan.TargetName.ValueString() == "" {
		resp.Diagnostics.AddError("Missing target_name", "target_name is required and cannot be empty.")
		return
	}

	request := teo.NewCreateAliasDomainRequest()
	request.ZoneId = plan.ZoneId.ValueStringPointer()
	request.AliasName = plan.AliasName.ValueStringPointer()
	request.TargetName = plan.TargetName.ValueStringPointer()
	if !plan.CertType.IsNull() && !plan.CertType.IsUnknown() {
		request.CertType = plan.CertType.ValueStringPointer()
	}
	if !plan.CertId.IsNull() && !plan.CertId.IsUnknown() {
		var certIds []types.String
		resp.Diagnostics.Append(plan.CertId.ElementsAs(ctx, &certIds, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		for _, id := range certIds {
			request.CertId = append(request.CertId, id.ValueStringPointer())
		}
	}

	service := NewTeoService(r.Client())
	if err := service.CreateTeoAliasDomain(ctx, request); err != nil {
		resp.Diagnostics.AddError(
			"Error creating TEO alias domain",
			fmt.Sprintf("CreateAliasDomain failed for zone %q alias %q: %v",
				plan.ZoneId.ValueString(), plan.AliasName.ValueString(), err),
		)
		return
	}

	// CreateAliasDomain returns no resource-level id, so the identity is the
	// zone id + alias name composite. Persist the identity before any read
	// back so a partial failure keeps the resource destroyable.
	plan.Id = types.StringValue(strings.Join(
		[]string{plan.ZoneId.ValueString(), plan.AliasName.ValueString()}, tccommon.FILED_SP))
	plan.AliasDomains = []TeoAliasDomainItem{}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Read back to align state with the cloud side. cert_type and cert_id are
	// not returned by DescribeAliasDomains, so the configured values stay in
	// state as-is.
	item, err := service.DescribeTeoAliasDomainById(ctx, plan.ZoneId.ValueString(), plan.AliasName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading TEO alias domain after create",
			fmt.Sprintf("DescribeAliasDomains failed for zone %q alias %q: %v",
				plan.ZoneId.ValueString(), plan.AliasName.ValueString(), err),
		)
		return
	}
	if item == nil {
		resp.Diagnostics.AddError(
			"TEO alias domain not readable after create",
			fmt.Sprintf("alias domain %q was created but cannot be read back", plan.Id.ValueString()),
		)
		return
	}
	plan.AliasDomains = []TeoAliasDomainItem{newTeoAliasDomain110902Item(item)}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Read handles `terraform refresh` and post-apply state refresh.
func (r *TeoAliasDomain110902Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state TeoAliasDomain110902ResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	zoneId, aliasName, ok := SplitTeoAliasDomain110902Id(state.Id.ValueString())
	if !ok {
		resp.Diagnostics.AddError(
			"Invalid resource ID",
			fmt.Sprintf("TEO alias domain id is broken, id is %q (expected `<zone_id>%s<alias_name>`)",
				state.Id.ValueString(), tccommon.FILED_SP),
		)
		return
	}

	service := NewTeoService(r.Client())
	item, err := service.DescribeTeoAliasDomainById(ctx, zoneId, aliasName)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading TEO alias domain",
			fmt.Sprintf("DescribeAliasDomains failed for zone %q alias %q: %v", zoneId, aliasName, err),
		)
		return
	}

	// DescribeTeoAliasDomainById returns a confirmed nil only when the query
	// succeeded and no matching alias domain exists; that is the only path
	// which removes the resource from state.
	if item == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	// zone_id and alias_name are the identity members: keep the configured
	// values so a normalized cloud-side echo never produces a diff.
	state.ZoneId = types.StringValue(zoneId)
	state.AliasName = types.StringValue(aliasName)
	state.AliasDomains = []TeoAliasDomainItem{newTeoAliasDomain110902Item(item)}
	// cert_type and cert_id are not returned by DescribeAliasDomains: keep the
	// state values instead of clearing them (drift detection for them is not
	// possible, documented in the resource description).
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update handles in-place updates.
func (r *TeoAliasDomain110902Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan TeoAliasDomain110902ResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state TeoAliasDomain110902ResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	zoneId, aliasName, ok := SplitTeoAliasDomain110902Id(state.Id.ValueString())
	if !ok {
		resp.Diagnostics.AddError(
			"Invalid resource ID",
			fmt.Sprintf("TEO alias domain id is broken, id is %q", state.Id.ValueString()),
		)
		return
	}

	request := teo.NewModifyAliasDomainRequest()
	request.ZoneId = helper.String(zoneId)
	request.AliasName = helper.String(aliasName)
	// target_name is required by ModifyAliasDomain and is updatable.
	request.TargetName = plan.TargetName.ValueStringPointer()

	// Contract (documented): omitting cert_type keeps the previous cloud-side
	// configuration ("an omitted value keeps the original configuration"). A
	// null plan value therefore maps to "do not send", and the final state
	// takes the plan value (null), so removing the attribute never produces a
	// permanent diff and never clears the cloud-side certificate. To switch to
	// "no certificate", set cert_type = "none" explicitly.
	if !plan.CertType.IsNull() && !plan.CertType.IsUnknown() {
		request.CertType = plan.CertType.ValueStringPointer()
	}

	// Same keep semantics for cert_id: a null plan value is not sent, and the
	// final state takes the plan value. The API never echoes cert_id back, so
	// drift detection for it is not possible (documented in the description).
	if !plan.CertId.IsNull() && !plan.CertId.IsUnknown() {
		var certIds []types.String
		resp.Diagnostics.Append(plan.CertId.ElementsAs(ctx, &certIds, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		for _, id := range certIds {
			request.CertId = append(request.CertId, id.ValueStringPointer())
		}
	}

	service := NewTeoService(r.Client())
	if err := service.ModifyTeoAliasDomain(ctx, request); err != nil {
		resp.Diagnostics.AddError(
			"Error updating TEO alias domain",
			fmt.Sprintf("ModifyAliasDomain failed for zone %q alias %q: %v", zoneId, aliasName, err),
		)
		return
	}

	item, err := service.DescribeTeoAliasDomainById(ctx, zoneId, aliasName)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading TEO alias domain after update",
			fmt.Sprintf("DescribeAliasDomains failed for zone %q alias %q: %v", zoneId, aliasName, err),
		)
		return
	}
	if item == nil {
		resp.Diagnostics.AddError(
			"TEO alias domain not readable after update",
			fmt.Sprintf("alias domain %q was updated but cannot be read back", state.Id.ValueString()),
		)
		return
	}

	plan.Id = state.Id
	if plan.ZoneId.IsNull() {
		plan.ZoneId = state.ZoneId
	}
	if plan.AliasName.IsNull() {
		plan.AliasName = state.AliasName
	}
	plan.AliasDomains = []TeoAliasDomainItem{newTeoAliasDomain110902Item(item)}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete handles resource deletion.
func (r *TeoAliasDomain110902Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state TeoAliasDomain110902ResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	zoneId, aliasName, ok := SplitTeoAliasDomain110902Id(state.Id.ValueString())
	if !ok {
		resp.Diagnostics.AddError(
			"Invalid resource ID",
			fmt.Sprintf("TEO alias domain id is broken, id is %q", state.Id.ValueString()),
		)
		return
	}

	service := NewTeoService(r.Client())

	// Delete only the alias domain itself; the target domain and any hosted
	// certificate are left untouched. A confirmed-missing alias domain is
	// already deleted, so delete is idempotent in that case.
	item, err := service.DescribeTeoAliasDomainById(ctx, zoneId, aliasName)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading TEO alias domain before delete",
			fmt.Sprintf("DescribeAliasDomains failed for zone %q alias %q: %v", zoneId, aliasName, err),
		)
		return
	}
	if item == nil {
		return
	}

	request := teo.NewDeleteAliasDomainRequest()
	request.ZoneId = helper.String(zoneId)
	// DeleteAliasDomain takes a list; a single alias domain is expressed as a
	// one-element list, and an empty list would perform no deletion at all.
	request.AliasNames = []*string{helper.String(aliasName)}

	if err := service.DeleteTeoAliasDomain(ctx, request); err != nil {
		resp.Diagnostics.AddError(
			"Error deleting TEO alias domain",
			fmt.Sprintf("DeleteAliasDomain failed for zone %q alias %q: %v", zoneId, aliasName, err),
		)
		return
	}
}

// ImportState supports `terraform import` using the `<zone_id>#<alias_name>` form.
func (r *TeoAliasDomain110902Resource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if _, _, ok := SplitTeoAliasDomain110902Id(req.ID); !ok {
		resp.Diagnostics.AddError(
			"Invalid import ID",
			fmt.Sprintf("Expected import ID in the form `<zone_id>%s<alias_name>`, got %q", tccommon.FILED_SP, req.ID),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
}

// SplitTeoAliasDomain110902Id parses `<zone_id>#<alias_name>`.
func SplitTeoAliasDomain110902Id(id string) (string, string, bool) {
	parts := strings.Split(id, tccommon.FILED_SP)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", false
	}
	return parts[0], parts[1], true
}

// newTeoAliasDomain110902Item converts a cloud-side AliasDomain into the
// read-only list element stored in state.
func newTeoAliasDomain110902Item(item *teo.AliasDomain) TeoAliasDomainItem {
	return TeoAliasDomainItem{
		AliasName:  types.StringPointerValue(item.AliasName),
		ZoneId:     types.StringPointerValue(item.ZoneId),
		TargetName: types.StringPointerValue(item.TargetName),
		Status:     types.StringPointerValue(item.Status),
		ForbidMode: types.Int64PointerValue(item.ForbidMode),
		CreatedOn:  types.StringPointerValue(item.CreatedOn),
		ModifiedOn: types.StringPointerValue(item.ModifiedOn),
	}
}

// aliasDomain110902AttrTypes returns the attribute types of one
// alias_domains list element.
func aliasDomain110902AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"alias_name":  types.StringType,
		"zone_id":     types.StringType,
		"target_name": types.StringType,
		"status":      types.StringType,
		"forbid_mode": types.Int64Type,
		"created_on":  types.StringType,
		"modified_on": types.StringType,
	}
}

// aliasDomain110902ObjectType returns the object type of one alias_domains
// list element.
func aliasDomain110902ObjectType() types.ObjectType {
	return types.ObjectType{AttrTypes: aliasDomain110902AttrTypes()}
}
