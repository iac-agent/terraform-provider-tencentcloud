package teo

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/fw"
)

// _ resource.Resource is intentionally omitted: the factory function
// NewTeoAliasDomain110903Resource returns resource.Resource, which already
// performs the implicit compile-time check at the return statement. Only the
// extended-interface assertions are kept here.
var (
	_ resource.ResourceWithConfigure   = &TeoAliasDomain110903Resource{}
	_ resource.ResourceWithImportState = &TeoAliasDomain110903Resource{}
)

// teoAliasDomain110903PollInterval is the interval between two status polls
// while waiting for a newly created alias domain to leave the pending state.
var teoAliasDomain110903PollInterval = 5 * time.Second

// NewTeoAliasDomain110903Resource is the factory referenced by
// tencentcloud/framework/registry.go to register this resource.
func NewTeoAliasDomain110903Resource() resource.Resource {
	return &TeoAliasDomain110903Resource{}
}

// TeoAliasDomain110903Resource implements resource.Resource for
// tencentcloud_teo_alias_domain_1109_03.
type TeoAliasDomain110903Resource struct {
	fw.ResourceWithConfigure
}

// TeoAliasDomain110903ResourceModel maps the schema attributes.
type TeoAliasDomain110903ResourceModel struct {
	Id         types.String `tfsdk:"id"`
	ZoneId     types.String `tfsdk:"zone_id"`
	AliasName  types.String `tfsdk:"alias_name"`
	TargetName types.String `tfsdk:"target_name"`
	CertType   types.String `tfsdk:"cert_type"`
	CertId     types.Set    `tfsdk:"cert_id"`
	// Read-only outputs of DescribeAliasDomains.
	Status     types.String `tfsdk:"status"`
	ForbidMode types.Int64  `tfsdk:"forbid_mode"`
	CreatedOn  types.String `tfsdk:"created_on"`
	ModifiedOn types.String `tfsdk:"modified_on"`
}

func (r *TeoAliasDomain110903Resource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "tencentcloud_teo_alias_domain_1109_03"
}

func (r *TeoAliasDomain110903Resource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a resource to create a TEO alias domain.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Resource ID in the form `<zone_id>#<alias_name>`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"zone_id": schema.StringAttribute{
				Required: true,
				Description: "Site ID. Changing it forces a new resource: ModifyAliasDomain cannot relocate " +
					"an alias domain to another site.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"alias_name": schema.StringAttribute{
				Required: true,
				Description: "Alias domain name. Changing it forces a new resource: ModifyAliasDomain uses " +
					"AliasName to locate the alias domain, so it cannot be renamed in place.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"target_name": schema.StringAttribute{
				Required:    true,
				Description: "Target domain name. Always sent on update because ModifyAliasDomain requires it.",
			},
			"cert_type": schema.StringAttribute{
				Optional: true,
				Description: "Certificate configuration. Valid values: `none` (no certificate), `hosting` (SSL " +
					"hosted certificate). Default value: `none`. When omitted on update, the cloud side keeps the " +
					"previous configuration; pass `none` explicitly to remove the certificate binding. Not returned " +
					"by DescribeAliasDomains, so drift cannot be detected on this field.",
			},
			"cert_id": schema.SetAttribute{
				Optional:    true,
				ElementType: types.StringType,
				Description: "Certificate ID list, required when `cert_type` is `hosting`. When omitted on update " +
					"the cloud side keeps the previous configuration. Not returned by DescribeAliasDomains, so " +
					"drift cannot be detected on this field.",
			},
			"status": schema.StringAttribute{
				Computed: true,
				Description: "Alias domain status. Valid values: `active` (effective), `pending` (deploying), " +
					"`conflict` (reclaimed), `stop` (disabled).",
			},
			"forbid_mode": schema.Int64Attribute{
				Computed: true,
				Description: "Block mode. Valid values: `0` (not blocked), `11` (compliance blocked), " +
					"`14` (not-registered blocked).",
			},
			"created_on": schema.StringAttribute{
				Computed:    true,
				Description: "Alias domain creation time.",
			},
			"modified_on": schema.StringAttribute{
				Computed:    true,
				Description: "Alias domain modification time.",
			},
		},
	}
}

// splitTeoAliasDomain110903Id parses `<zone_id>#<alias_name>`.
func splitTeoAliasDomain110903Id(id string) (string, string, bool) {
	parts := strings.Split(id, tccommon.FILED_SP)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", false
	}
	return parts[0], parts[1], true
}

// setToStringPtrs converts a types.Set of strings to a string pointer slice.
// A null or unknown set yields nil; null or unknown elements are skipped.
func setToStringPtrs(v types.Set) []*string {
	if v.IsNull() || v.IsUnknown() {
		return nil
	}
	out := make([]*string, 0, len(v.Elements()))
	for _, e := range v.Elements() {
		s, ok := e.(types.String)
		if !ok || s.IsNull() || s.IsUnknown() {
			continue
		}
		val := s.ValueString()
		out = append(out, &val)
	}
	return out
}

// Create handles `terraform apply` for new resources.
func (r *TeoAliasDomain110903Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan TeoAliasDomain110903ResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	request := teov20220901.NewCreateAliasDomainRequest()
	request.ZoneId = plan.ZoneId.ValueStringPointer()
	request.AliasName = plan.AliasName.ValueStringPointer()
	request.TargetName = plan.TargetName.ValueStringPointer()
	if !plan.CertType.IsNull() && !plan.CertType.IsUnknown() {
		request.CertType = plan.CertType.ValueStringPointer()
	}
	request.CertId = setToStringPtrs(plan.CertId)

	service := NewTeoService(r.Client())
	if err := service.CreateTeoAliasDomain(ctx, request); err != nil {
		resp.Diagnostics.AddError(
			"Error creating TEO alias domain",
			fmt.Sprintf("CreateAliasDomain failed: %v", err),
		)
		return
	}

	zoneId := plan.ZoneId.ValueString()
	aliasName := plan.AliasName.ValueString()
	plan.Id = types.StringValue(strings.Join([]string{zoneId, aliasName}, tccommon.FILED_SP))

	// Persist the identity BEFORE waiting for the status to settle: if the wait
	// fails, the resource is still in state and `terraform destroy` can clean it up.
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Creation is asynchronous: the alias domain shows up as "pending" until the
	// deployment finishes. Wait until the status leaves "pending".
	item, err := r.waitTeoAliasDomain110903Settled(ctx, zoneId, aliasName)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error waiting for TEO alias domain creation",
			fmt.Sprintf("alias domain %s was created but did not settle: %v", plan.Id.ValueString(), err),
		)
		return
	}

	// cert_type / cert_id are not returned by DescribeAliasDomains: the values
	// from the plan are kept (drift detection is limited on these two fields).
	refreshTeoAliasDomain110903Model(&plan, item)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Read handles `terraform refresh` and post-apply state refresh.
func (r *TeoAliasDomain110903Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state TeoAliasDomain110903ResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	zoneId, aliasName, ok := splitTeoAliasDomain110903Id(state.Id.ValueString())
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
			fmt.Sprintf("DescribeAliasDomains failed for %s: %v", state.Id.ValueString(), err),
		)
		return
	}

	// Confirmed "not found" on the cloud side => drop it from state so Terraform
	// plans a recreation. A nil response is an API failure and must not land here.
	if item == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	// cert_type / cert_id keep the state values: DescribeAliasDomains does not
	// return them, and clearing them would silently drop the user's binding.
	refreshTeoAliasDomain110903Model(&state, item)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update handles in-place updates.
func (r *TeoAliasDomain110903Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan TeoAliasDomain110903ResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state TeoAliasDomain110903ResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	zoneId, aliasName, ok := splitTeoAliasDomain110903Id(state.Id.ValueString())
	if !ok {
		resp.Diagnostics.AddError(
			"Invalid resource ID",
			fmt.Sprintf("TEO alias domain id is broken, id is %q (expected `<zone_id>%s<alias_name>`)",
				state.Id.ValueString(), tccommon.FILED_SP),
		)
		return
	}

	request := teov20220901.NewModifyAliasDomainRequest()
	request.ZoneId = &zoneId
	request.AliasName = &aliasName
	// TargetName is required by ModifyAliasDomain and is always carried.
	request.TargetName = plan.TargetName.ValueStringPointer()

	// Three states must stay distinguishable -- do NOT collapse them:
	//   Unknown -> value not decided yet, leave it alone
	//   Null    -> user removed the attribute from HCL; confirmed contract
	//              (evaluation "ClearSemantics: keep"): the cloud side keeps the
	//              previous configuration, so nothing is sent
	//   value   -> send it (an explicit `none` clears the certificate binding)
	if !plan.CertType.IsNull() && !plan.CertType.IsUnknown() {
		request.CertType = plan.CertType.ValueStringPointer()
	}
	request.CertId = setToStringPtrs(plan.CertId)

	service := NewTeoService(r.Client())
	if err := service.ModifyTeoAliasDomain(ctx, request); err != nil {
		resp.Diagnostics.AddError(
			"Error updating TEO alias domain",
			fmt.Sprintf("ModifyAliasDomain failed for %s: %v", state.Id.ValueString(), err),
		)
		return
	}

	item, err := service.DescribeTeoAliasDomainById(ctx, zoneId, aliasName)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading TEO alias domain after update",
			fmt.Sprintf("DescribeAliasDomains failed for %s: %v", state.Id.ValueString(), err),
		)
		return
	}
	if item == nil {
		resp.Diagnostics.AddError(
			"TEO alias domain not readable after update",
			fmt.Sprintf("alias domain %s was updated but cannot be read back", state.Id.ValueString()),
		)
		return
	}

	plan.Id = state.Id
	// cert_type / cert_id keep the planned values (not returned by Describe).
	refreshTeoAliasDomain110903Model(&plan, item)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete handles resource deletion. It only removes the alias domain itself,
// never the target domain or the site.
func (r *TeoAliasDomain110903Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state TeoAliasDomain110903ResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	zoneId, aliasName, ok := splitTeoAliasDomain110903Id(state.Id.ValueString())
	if !ok {
		resp.Diagnostics.AddError(
			"Invalid resource ID",
			fmt.Sprintf("TEO alias domain id is broken, id is %q (expected `<zone_id>%s<alias_name>`)",
				state.Id.ValueString(), tccommon.FILED_SP),
		)
		return
	}

	service := NewTeoService(r.Client())

	// Query before delete: DeleteAliasDomain does not document its behavior on
	// a missing alias domain, so the absence is confirmed with a query first.
	item, err := service.DescribeTeoAliasDomainById(ctx, zoneId, aliasName)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading TEO alias domain before delete",
			fmt.Sprintf("DescribeAliasDomains failed for %s: %v", state.Id.ValueString(), err),
		)
		return
	}
	if item == nil {
		// Already gone on the cloud side: nothing to delete.
		return
	}

	if err := service.DeleteTeoAliasDomainById(ctx, zoneId, aliasName); err != nil {
		resp.Diagnostics.AddError(
			"Error deleting TEO alias domain",
			fmt.Sprintf("DeleteAliasDomain failed for %s: %v", state.Id.ValueString(), err),
		)
		return
	}
}

// ImportState supports `terraform import` using the `<zone_id>#<alias_name>` form.
func (r *TeoAliasDomain110903Resource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	zoneId, aliasName, ok := splitTeoAliasDomain110903Id(req.ID)
	if !ok {
		resp.Diagnostics.AddError(
			"Invalid import ID",
			fmt.Sprintf("Expected import ID in the form `<zone_id>%s<alias_name>`, got %q", tccommon.FILED_SP, req.ID),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("zone_id"), zoneId)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("alias_name"), aliasName)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
}

// waitTeoAliasDomain110903Settled polls DescribeAliasDomains until the alias
// domain is visible and its status is no longer "pending", and returns the
// last observed item.
func (r *TeoAliasDomain110903Resource) waitTeoAliasDomain110903Settled(ctx context.Context, zoneId, aliasName string) (*teov20220901.AliasDomain, error) {
	service := NewTeoService(r.Client())
	deadline := time.Now().Add(tccommon.WriteRetryTimeout)
	lastStatus := ""
	for {
		item, err := service.DescribeTeoAliasDomainById(ctx, zoneId, aliasName)
		if err != nil {
			return nil, err
		}
		if item != nil {
			lastStatus = ""
			if item.Status != nil {
				lastStatus = *item.Status
			}
			if lastStatus != "pending" {
				return item, nil
			}
		}
		if !time.Now().Before(deadline) {
			return nil, fmt.Errorf("alias domain %s did not settle within %v (last status: %q)",
				aliasName, tccommon.WriteRetryTimeout, lastStatus)
		}
		time.Sleep(teoAliasDomain110903PollInterval)
	}
}

// refreshTeoAliasDomain110903Model copies cloud-owned fields into the model and
// keeps the user-supplied values for fields the API does not return.
func refreshTeoAliasDomain110903Model(m *TeoAliasDomain110903ResourceModel, item *teov20220901.AliasDomain) {
	if item == nil {
		return
	}
	m.ZoneId = types.StringPointerValue(item.ZoneId)
	m.AliasName = types.StringPointerValue(item.AliasName)
	m.TargetName = types.StringPointerValue(item.TargetName)
	m.Status = types.StringPointerValue(item.Status)
	m.ForbidMode = types.Int64PointerValue(item.ForbidMode)
	m.CreatedOn = types.StringPointerValue(item.CreatedOn)
	m.ModifiedOn = types.StringPointerValue(item.ModifiedOn)
	// cert_type / cert_id: intentionally kept from the current model. The
	// DescribeAliasDomains response has no counterpart fields.
}
