package teo

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	teo "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/fw"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

// _ resource.Resource is intentionally omitted: the factory function
// NewTeoZone1110_01Resource returns resource.Resource, which already performs
// the implicit compile-time check at the return statement. Only the extended
// interface assertions are kept here.
var (
	_ resource.ResourceWithConfigure   = &TeoZone111001Resource{}
	_ resource.ResourceWithImportState = &TeoZone111001Resource{}
)

// NewTeoZone1110_01Resource is the factory referenced by
// tencentcloud/framework/registry.go to register this resource.
func NewTeoZone1110_01Resource() resource.Resource {
	return &TeoZone111001Resource{}
}

// TeoZone111001Resource implements resource.Resource for
// tencentcloud_teo_zone_1110_01.
type TeoZone111001Resource struct {
	fw.ResourceWithConfigure
}

// TeoZone1110_01TagModel maps one element of the tags list.
type TeoZone1110_01TagModel struct {
	TagKey   types.String `tfsdk:"tag_key"`
	TagValue types.String `tfsdk:"tag_value"`
}

// TeoZone1110_01DnsVerificationModel maps DNS record based verification information.
type TeoZone1110_01DnsVerificationModel struct {
	Subdomain   types.String `tfsdk:"subdomain"`
	RecordType  types.String `tfsdk:"record_type"`
	RecordValue types.String `tfsdk:"record_value"`
}

// TeoZone1110_01FileVerificationModel maps file based verification information.
type TeoZone1110_01FileVerificationModel struct {
	Path    types.String `tfsdk:"path"`
	Content types.String `tfsdk:"content"`
}

// TeoZone1110_01NsVerificationModel maps the name servers assigned to the zone.
type TeoZone1110_01NsVerificationModel struct {
	NameServers types.Set `tfsdk:"name_servers"`
}

// TeoZone1110_01OwnershipVerificationModel maps the ownership verification block.
type TeoZone1110_01OwnershipVerificationModel struct {
	DnsVerification  types.List `tfsdk:"dns_verification"`
	FileVerification types.List `tfsdk:"file_verification"`
	NsVerification   types.List `tfsdk:"ns_verification"`
}

// TeoZone111001ResourceModel maps the schema attributes.
//
// AllowDuplicates and JumpStart are accepted by CreateZone only and are never
// returned by DescribeZones. They are kept verbatim in state, which means no
// drift is reported for them: neither an out-of-band change nor a plain
// `terraform import` can be detected for these two attributes.
//
// OwnershipVerification is a CreateZone output. It is populated once at creation
// time and never refreshed, because the verification material is only meaningful
// for the pending verification flow.
type TeoZone111001ResourceModel struct {
	Id                    types.String `tfsdk:"id"`
	ZoneId                types.String `tfsdk:"zone_id"`
	Type                  types.String `tfsdk:"type"`
	ZoneName              types.String `tfsdk:"zone_name"`
	Area                  types.String `tfsdk:"area"`
	PlanId                types.String `tfsdk:"plan_id"`
	AliasZoneName         types.String `tfsdk:"alias_zone_name"`
	Tags                  types.List   `tfsdk:"tags"`
	AllowDuplicates       types.Bool   `tfsdk:"allow_duplicates"`
	JumpStart             types.Bool   `tfsdk:"jump_start"`
	OwnershipVerification types.List   `tfsdk:"ownership_verification"`
}

// Attribute types of the nested blocks. They are declared at package level
// because the schema, plan and state conversions have to agree on them exactly.
var (
	teoZone1110_01TagAttrTypes = map[string]attr.Type{
		"tag_key":   types.StringType,
		"tag_value": types.StringType,
	}

	teoZone1110_01DnsVerificationAttrTypes = map[string]attr.Type{
		"subdomain":    types.StringType,
		"record_type":  types.StringType,
		"record_value": types.StringType,
	}

	teoZone1110_01FileVerificationAttrTypes = map[string]attr.Type{
		"path":    types.StringType,
		"content": types.StringType,
	}

	teoZone1110_01NsVerificationAttrTypes = map[string]attr.Type{
		"name_servers": types.SetType{ElemType: types.StringType},
	}

	teoZone1110_01OwnershipVerificationAttrTypes = map[string]attr.Type{
		"dns_verification":  types.ListType{ElemType: types.ObjectType{AttrTypes: teoZone1110_01DnsVerificationAttrTypes}},
		"file_verification": types.ListType{ElemType: types.ObjectType{AttrTypes: teoZone1110_01FileVerificationAttrTypes}},
		"ns_verification":   types.ListType{ElemType: types.ObjectType{AttrTypes: teoZone1110_01NsVerificationAttrTypes}},
	}

	teoZone1110_01TagObjectType                   = types.ObjectType{AttrTypes: teoZone1110_01TagAttrTypes}
	teoZone1110_01DnsVerificationObjectType       = types.ObjectType{AttrTypes: teoZone1110_01DnsVerificationAttrTypes}
	teoZone1110_01FileVerificationObjectType      = types.ObjectType{AttrTypes: teoZone1110_01FileVerificationAttrTypes}
	teoZone1110_01NsVerificationObjectType        = types.ObjectType{AttrTypes: teoZone1110_01NsVerificationAttrTypes}
	teoZone1110_01OwnershipVerificationObjectType = types.ObjectType{AttrTypes: teoZone1110_01OwnershipVerificationAttrTypes}
)

func (r *TeoZone111001Resource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "tencentcloud_teo_zone_1110_01"
}

func (r *TeoZone111001Resource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a resource to create a TEO zone. A freshly created zone only serves traffic after a " +
			"plan is bound to it and after the ownership verification is finished. Note that `allow_duplicates` and " +
			"`jump_start` are create-only switches that the cloud API does not return, so out-of-band changes of " +
			"these two attributes cannot be detected.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "ID of the resource. It is identical to `zone_id`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"zone_id": schema.StringAttribute{
				Computed:    true,
				Description: "Zone ID, for example `zone-2noz78a8ev6k`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"type": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Description: "Access type of the zone. Valid values: `partial` (CNAME access), `full` (NS access), " +
					"`noDomainAccess` (access without a domain), `dnsPodAccess` (DNSPod hosted access, the domain has " +
					"to be hosted in DNSPod already) and `ai` (edge inference access). Default value: `partial`. " +
					"`noDomainAccess` is a create-only value: an existing zone cannot be switched back to it. " +
					"When it is used, both `zone_name` and `area` have to be empty.",
			},
			"zone_name": schema.StringAttribute{
				Optional: true,
				Description: "Zone name. Pass the second-level domain (for example `example.com`) for CNAME/NS access " +
					"and leave it empty for access without a domain. The name can only be changed when the zone is " +
					"switched from access without a domain to CNAME access, so changing it in any other case replaces " +
					"the zone.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplaceIfConfigured(),
				},
			},
			"area": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Description: "Acceleration area of the zone when `type` is `partial` or `full`. Valid values: `global`, " +
					"`mainland`, `overseas`. Default value: `overseas`. It has to be empty when `type` is " +
					"`noDomainAccess`.",
			},
			"plan_id": schema.StringAttribute{
				Optional: true,
				Description: "ID of the plan to be bound to the zone. When a plan already exists in the account this " +
					"parameter binds the zone to it directly. When it is omitted, a zone in the `init` state is created: " +
					"such a zone is not activated and does not show up in the console until a plan is bound to it. " +
					"Changing this value replaces the zone.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"alias_zone_name": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Description: "Alias of the zone name, used to distinguish zones that share the same zone name. " +
					"Up to 200 characters, only digits, letters, `.`, `-` and `_` are allowed. Leave it empty when this " +
					"feature is not needed. Removing the attribute from the configuration is rejected, because the API " +
					"does not document whether an omitted value keeps or clears the existing alias.",
			},
			"tags": schema.ListNestedAttribute{
				Optional: true,
				Computed: true,
				Description: "Tags of the zone, used for permission and billing management. The tags have to be created " +
					"in the tag console before they can be used here.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"tag_key": schema.StringAttribute{
							Required:    true,
							Description: "Tag key.",
						},
						"tag_value": schema.StringAttribute{
							Required:    true,
							Description: "Tag value.",
						},
					},
				},
			},
			"allow_duplicates": schema.BoolAttribute{
				Optional: true,
				Description: "Whether duplicating the access of an existing zone is allowed. `true` allows duplicate " +
					"access, `false` does not. Default value: `false`. This create-only switch is not returned by the " +
					"cloud API, so the configured value is kept in state.",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
			},
			"jump_start": schema.BoolAttribute{
				Optional: true,
				Description: "Whether the existing DNS record scan of the zone is skipped. Default value: `false`. " +
					"This create-only switch is not returned by the cloud API, so the configured value is kept in state.",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
			},
			"ownership_verification": schema.ListNestedAttribute{
				Computed: true,
				Description: "Ownership verification information returned by CreateZone. The zone only serves traffic " +
					"after the ownership verification is finished. It is a creation-time output: `terraform import` " +
					"leaves it empty.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"dns_verification": schema.ListNestedAttribute{
							Computed:    true,
							Description: "Information needed to verify the ownership by adding a DNS record.",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"subdomain": schema.StringAttribute{
										Computed:    true,
										Description: "Host record.",
									},
									"record_type": schema.StringAttribute{
										Computed:    true,
										Description: "Record type.",
									},
									"record_value": schema.StringAttribute{
										Computed:    true,
										Description: "Record value.",
									},
								},
							},
						},
						"file_verification": schema.ListNestedAttribute{
							Computed:    true,
							Description: "Information needed to verify the ownership by uploading a file.",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"path": schema.StringAttribute{
										Computed: true,
										Description: "URL path of the verification file, which is fetched by the EdgeOne " +
											"backend via `http://{Host}{Path}`.",
									},
									"content": schema.StringAttribute{
										Computed:    true,
										Description: "Content that has to be written into the file located at `path`.",
									},
								},
							},
						},
						"ns_verification": schema.ListNestedAttribute{
							Computed:    true,
							Description: "Information needed to verify the ownership by switching the name servers.",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"name_servers": schema.SetAttribute{
										Computed:    true,
										ElementType: types.StringType,
										Description: "Name servers assigned to the zone. The name servers of the domain " +
											"have to be switched to these addresses.",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

// Create handles `terraform apply` for new zones.
func (r *TeoZone111001Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan TeoZone111001ResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	request, diags := buildTeoZone1110_01CreateZoneRequest(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	service := NewTeoService(r.Client())
	response, err := service.CreateTeoZoneBy1110_01(ctx, request)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating TEO zone",
			fmt.Sprintf("CreateZone failed: %v", err),
		)
		return
	}

	// A nil response (or a nil zone id) is an API failure, not an empty result.
	// Never write an empty id into state: the zone would become unmanageable.
	if response == nil || response.Response == nil || response.Response.ZoneId == nil || *response.Response.ZoneId == "" {
		resp.Diagnostics.AddError(
			"Error creating TEO zone",
			"CreateZone returned an empty response or an empty zone id.",
		)
		return
	}

	zoneId := *response.Response.ZoneId
	plan.Id = types.StringValue(zoneId)
	plan.ZoneId = types.StringValue(zoneId)
	plan.OwnershipVerification = buildTeoZone1110_01OwnershipVerification(ctx, response.Response.OwnershipVerification, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	// Values owned by the cloud side are only known after the read back below.
	plan.Type = types.StringNull()
	plan.Area = types.StringNull()
	plan.AliasZoneName = types.StringNull()
	plan.Tags = types.ListNull(teoZone1110_01TagObjectType)

	// Persist the identity BEFORE reading the zone back: if the read fails the
	// zone is still tracked in state and `terraform destroy` can clean it up.
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	zone, err := service.DescribeTeoZoneBy1110_01(ctx, zoneId)
	if err != nil {
		resp.Diagnostics.AddError("Error reading TEO zone after create", err.Error())
		return
	}
	if zone == nil {
		resp.Diagnostics.AddError(
			"TEO zone not readable after create",
			fmt.Sprintf("zone %s was created but cannot be read back", zoneId),
		)
		return
	}

	resp.Diagnostics.Append(refreshTeoZone1110_01Model(ctx, &plan, zone)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Read handles `terraform refresh` and post-apply state refresh.
func (r *TeoZone111001Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state TeoZone111001ResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	zoneId := state.ZoneId.ValueString()
	if zoneId == "" {
		zoneId = state.Id.ValueString()
	}
	if zoneId == "" {
		resp.State.RemoveResource(ctx)
		return
	}

	service := NewTeoService(r.Client())
	zone, err := service.DescribeTeoZoneBy1110_01(ctx, zoneId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading TEO zone",
			fmt.Sprintf("DescribeZones failed for zone %s: %v", zoneId, err),
		)
		return
	}

	// Confirmed "not found" on the cloud side => drop the zone from state so
	// Terraform plans a recreation. A nil response is an API failure and is
	// reported as an error above instead of landing here.
	if zone == nil {
		log.Printf("[WARN] TEO zone %s not found, removing it from state", zoneId)
		resp.State.RemoveResource(ctx)
		return
	}

	state.Id = types.StringValue(zoneId)
	resp.Diagnostics.Append(refreshTeoZone1110_01Model(ctx, &state, zone)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update handles in-place modifications.
func (r *TeoZone111001Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan TeoZone111001ResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state TeoZone111001ResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	zoneId := state.Id.ValueString()
	request, diags := buildTeoZone1110_01ModifyZoneRequest(&plan, &state, zoneId)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	service := NewTeoService(r.Client())
	if err := service.ModifyTeoZoneBy1110_01(ctx, request); err != nil {
		resp.Diagnostics.AddError(
			"Error updating TEO zone",
			fmt.Sprintf("ModifyZone failed for zone %s: %v", zoneId, err),
		)
		return
	}

	zone, err := service.DescribeTeoZoneBy1110_01(ctx, zoneId)
	if err != nil {
		resp.Diagnostics.AddError("Error reading TEO zone after update", err.Error())
		return
	}
	if zone == nil {
		resp.Diagnostics.AddError(
			"TEO zone not readable after update",
			fmt.Sprintf("zone %s was updated but cannot be read back", zoneId),
		)
		return
	}

	plan.Id = state.Id
	resp.Diagnostics.Append(refreshTeoZone1110_01Model(ctx, &plan, zone)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete handles zone deletion.
func (r *TeoZone111001Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state TeoZone111001ResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	zoneId := state.Id.ValueString()
	if zoneId == "" {
		zoneId = state.ZoneId.ValueString()
	}
	if zoneId == "" {
		return
	}

	service := NewTeoService(r.Client())
	if err := service.DeleteTeoZoneBy1110_01(ctx, zoneId); err != nil {
		// Only an error code documented as "already gone" may be swallowed here;
		// anything else must fail so the zone stays in state for a retry.
		if !tccommon.IsExpectError(err, []string{"ResourceNotFound", "InvalidParameter.ZoneNotExists"}) {
			resp.Diagnostics.AddError(
				"Error deleting TEO zone",
				fmt.Sprintf("DeleteZone failed for zone %s: %v", zoneId, err),
			)
			return
		}
		log.Printf("[WARN] TEO zone %s is already gone, ignoring the delete error: %v", zoneId, err)
	}
	// After a successful delete there is no need to explicitly clear state:
	// the Framework removes it automatically.
}

// ImportState is implemented explicitly because the imported state has to carry
// both id and zone_id; the bare ID passthrough of fw.WithImportByID would leave
// zone_id unknown and the first plan would try to replace the zone.
func (r *TeoZone111001Resource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("zone_id"), req.ID)...)
}

// buildTeoZone1110_01CreateZoneRequest converts the plan into a CreateZone
// request. Only attributes that are neither null nor unknown are sent, so the
// cloud side applies its own defaults for the remaining ones.
func buildTeoZone1110_01CreateZoneRequest(ctx context.Context, plan *TeoZone111001ResourceModel) (*teo.CreateZoneRequest, diag.Diagnostics) {
	var diags diag.Diagnostics

	request := teo.NewCreateZoneRequest()
	if !plan.Type.IsNull() && !plan.Type.IsUnknown() {
		request.Type = plan.Type.ValueStringPointer()
	}
	if !plan.ZoneName.IsNull() && !plan.ZoneName.IsUnknown() {
		request.ZoneName = plan.ZoneName.ValueStringPointer()
	}
	if !plan.Area.IsNull() && !plan.Area.IsUnknown() {
		request.Area = plan.Area.ValueStringPointer()
	}
	if !plan.PlanId.IsNull() && !plan.PlanId.IsUnknown() {
		request.PlanId = plan.PlanId.ValueStringPointer()
	}
	if !plan.AliasZoneName.IsNull() && !plan.AliasZoneName.IsUnknown() {
		request.AliasZoneName = plan.AliasZoneName.ValueStringPointer()
	}
	if !plan.AllowDuplicates.IsNull() && !plan.AllowDuplicates.IsUnknown() {
		request.AllowDuplicates = plan.AllowDuplicates.ValueBoolPointer()
	}
	if !plan.JumpStart.IsNull() && !plan.JumpStart.IsUnknown() {
		request.JumpStart = plan.JumpStart.ValueBoolPointer()
	}

	if !plan.Tags.IsNull() && !plan.Tags.IsUnknown() {
		var tags []TeoZone1110_01TagModel
		diags.Append(plan.Tags.ElementsAs(ctx, &tags, false)...)
		if diags.HasError() {
			return request, diags
		}
		for i := range tags {
			request.Tags = append(request.Tags, &teo.Tag{
				TagKey:   tags[i].TagKey.ValueStringPointer(),
				TagValue: tags[i].TagValue.ValueStringPointer(),
			})
		}
	}

	return request, diags
}

// buildTeoZone1110_01ModifyZoneRequest converts the plan into a ModifyZone
// request.
//
// The API does not document whether an omitted value keeps or clears a field, so
// every mutable attribute is handled explicitly:
//   - unknown          -> the value is not decided yet, nothing is sent
//   - null + old value -> the attribute was removed from the configuration; the
//     update is refused instead of guessing a direction
//   - new value        -> sent
//
// zone_name is only sent when it actually changed, because ModifyZone applies it
// exclusively when the zone is switched from access without a domain to CNAME
// access (any other change is expressed as a replacement, not as an update).
func buildTeoZone1110_01ModifyZoneRequest(plan, state *TeoZone111001ResourceModel, zoneId string) (*teo.ModifyZoneRequest, diag.Diagnostics) {
	request := teo.NewModifyZoneRequest()
	request.ZoneId = helper.String(zoneId)

	if !plan.Type.IsNull() && !plan.Type.IsUnknown() {
		request.Type = plan.Type.ValueStringPointer()
	}
	if !plan.Area.IsNull() && !plan.Area.IsUnknown() {
		request.Area = plan.Area.ValueStringPointer()
	}
	if !plan.ZoneName.IsNull() && !plan.ZoneName.IsUnknown() &&
		plan.ZoneName.ValueString() != state.ZoneName.ValueString() {
		request.ZoneName = plan.ZoneName.ValueStringPointer()
	}

	if plan.AliasZoneName.IsNull() && !state.AliasZoneName.IsNull() {
		var diags diag.Diagnostics
		diags.AddError(
			"Removing alias_zone_name is not supported",
			"The TEO ModifyZone API does not document whether an omitted alias zone name keeps or clears the "+
				"existing value. Keep the attribute in the configuration, or replace the zone to remove the alias.",
		)
		return request, diags
	}
	if !plan.AliasZoneName.IsNull() && !plan.AliasZoneName.IsUnknown() {
		request.AliasZoneName = plan.AliasZoneName.ValueStringPointer()
	}

	return request, nil
}

// refreshTeoZone1110_01Model copies the cloud side values into the model.
//
// allow_duplicates, jump_start and ownership_verification are deliberately left
// untouched: the cloud API never returns the first two, and the verification
// material is only meaningful right after creation.
func refreshTeoZone1110_01Model(ctx context.Context, m *TeoZone111001ResourceModel, zone *teo.Zone) diag.Diagnostics {
	var diags diag.Diagnostics
	if zone == nil {
		return diags
	}

	m.ZoneId = types.StringPointerValue(zone.ZoneId)
	m.Type = types.StringPointerValue(zone.Type)
	m.ZoneName = types.StringPointerValue(zone.ZoneName)
	m.Area = types.StringPointerValue(zone.Area)
	m.AliasZoneName = types.StringPointerValue(zone.AliasZoneName)

	tags := make([]TeoZone1110_01TagModel, 0, len(zone.Tags))
	for _, tag := range zone.Tags {
		if tag == nil {
			continue
		}
		tags = append(tags, TeoZone1110_01TagModel{
			TagKey:   types.StringPointerValue(tag.TagKey),
			TagValue: types.StringPointerValue(tag.TagValue),
		})
	}
	tagList, tagDiags := types.ListValueFrom(ctx, teoZone1110_01TagObjectType, tags)
	diags.Append(tagDiags...)
	if tagDiags.HasError() {
		return diags
	}
	m.Tags = tagList

	return diags
}

// buildTeoZone1110_01OwnershipVerification converts the CreateZone ownership
// verification output into the state representation. A nil input yields a null
// list instead of an empty one, so "no verification needed" (for example access
// without a domain) stays distinguishable from an empty response.
func buildTeoZone1110_01OwnershipVerification(
	ctx context.Context,
	verification *teo.OwnershipVerification,
	diags *diag.Diagnostics,
) types.List {
	if verification == nil {
		return types.ListNull(teoZone1110_01OwnershipVerificationObjectType)
	}

	dnsVerifications := make([]TeoZone1110_01DnsVerificationModel, 0, 1)
	if verification.DnsVerification != nil {
		dnsVerifications = append(dnsVerifications, TeoZone1110_01DnsVerificationModel{
			Subdomain:   types.StringPointerValue(verification.DnsVerification.Subdomain),
			RecordType:  types.StringPointerValue(verification.DnsVerification.RecordType),
			RecordValue: types.StringPointerValue(verification.DnsVerification.RecordValue),
		})
	}
	dnsList, dnsDiags := types.ListValueFrom(ctx, teoZone1110_01DnsVerificationObjectType, dnsVerifications)
	diags.Append(dnsDiags...)
	if dnsDiags.HasError() {
		return types.ListNull(teoZone1110_01OwnershipVerificationObjectType)
	}

	fileVerifications := make([]TeoZone1110_01FileVerificationModel, 0, 1)
	if verification.FileVerification != nil {
		fileVerifications = append(fileVerifications, TeoZone1110_01FileVerificationModel{
			Path:    types.StringPointerValue(verification.FileVerification.Path),
			Content: types.StringPointerValue(verification.FileVerification.Content),
		})
	}
	fileList, fileDiags := types.ListValueFrom(ctx, teoZone1110_01FileVerificationObjectType, fileVerifications)
	diags.Append(fileDiags...)
	if fileDiags.HasError() {
		return types.ListNull(teoZone1110_01OwnershipVerificationObjectType)
	}

	nsVerifications := make([]TeoZone1110_01NsVerificationModel, 0, 1)
	if verification.NsVerification != nil {
		nameServers, nsDiags := types.SetValueFrom(
			ctx,
			types.StringType,
			flattenTeoZone1110_01Strings(verification.NsVerification.NameServers),
		)
		diags.Append(nsDiags...)
		if nsDiags.HasError() {
			return types.ListNull(teoZone1110_01OwnershipVerificationObjectType)
		}
		nsVerifications = append(nsVerifications, TeoZone1110_01NsVerificationModel{
			NameServers: nameServers,
		})
	}
	nsList, nsDiags := types.ListValueFrom(ctx, teoZone1110_01NsVerificationObjectType, nsVerifications)
	diags.Append(nsDiags...)
	if nsDiags.HasError() {
		return types.ListNull(teoZone1110_01OwnershipVerificationObjectType)
	}

	verificationValue, verificationDiags := types.ObjectValue(
		teoZone1110_01OwnershipVerificationAttrTypes,
		map[string]attr.Value{
			"dns_verification":  dnsList,
			"file_verification": fileList,
			"ns_verification":   nsList,
		},
	)
	diags.Append(verificationDiags...)
	if verificationDiags.HasError() {
		return types.ListNull(teoZone1110_01OwnershipVerificationObjectType)
	}

	list, listDiags := types.ListValue(
		teoZone1110_01OwnershipVerificationObjectType,
		[]attr.Value{verificationValue},
	)
	diags.Append(listDiags...)
	if listDiags.HasError() {
		return types.ListNull(teoZone1110_01OwnershipVerificationObjectType)
	}
	return list
}

// flattenTeoZone1110_01Strings converts a slice of string pointers into a slice
// of strings, skipping nil entries.
func flattenTeoZone1110_01Strings(values []*string) []string {
	result := make([]string, 0, len(values))
	for _, value := range values {
		if value == nil {
			continue
		}
		result = append(result, *value)
	}
	return result
}
