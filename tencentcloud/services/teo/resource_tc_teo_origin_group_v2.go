package teo

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/fw"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	teo "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"
)

var (
	_ resource.ResourceWithConfigure   = &TeoOriginGroupV2Resource{}
	_ resource.ResourceWithImportState = &TeoOriginGroupV2Resource{}
)

// NewTeoOriginGroupV2Resource is the factory referenced by
// tencentcloud/framework/registry.go to register this resource.
func NewTeoOriginGroupV2Resource() resource.Resource {
	return &TeoOriginGroupV2Resource{}
}

// TeoOriginGroupV2Resource implements resource.Resource for
// tencentcloud_teo_origin_group_v2.
type TeoOriginGroupV2Resource struct {
	fw.ResourceWithConfigure
}

// TeoOriginGroupV2ResourceModel maps the schema attributes.
type TeoOriginGroupV2ResourceModel struct {
	Id          types.String `tfsdk:"id"`
	ZoneId      types.String `tfsdk:"zone_id"`
	Name        types.String `tfsdk:"name"`
	Type        types.String `tfsdk:"type"`
	HostHeader  types.String `tfsdk:"host_header"`
	Records     types.List   `tfsdk:"records"`
	CreateTime  types.String `tfsdk:"create_time"`
	UpdateTime  types.String `tfsdk:"update_time"`
	References  types.List   `tfsdk:"references"`
}

// TeoOriginGroupV2RecordModel maps one origin record of the records list.
type TeoOriginGroupV2RecordModel struct {
	Record           types.String `tfsdk:"record"`
	Type             types.String `tfsdk:"type"`
	RecordId         types.String `tfsdk:"record_id"`
	Weight           types.Int64  `tfsdk:"weight"`
	Private          types.Bool   `tfsdk:"private"`
	PrivateParameters types.List  `tfsdk:"private_parameters"`
}

// TeoOriginGroupV2PrivateParameterModel maps one private auth parameter.
type TeoOriginGroupV2PrivateParameterModel struct {
	Name  types.String `tfsdk:"name"`
	Value types.String `tfsdk:"value"`
}

// TeoOriginGroupV2ReferenceModel maps one reference entry of the references list.
type TeoOriginGroupV2ReferenceModel struct {
	InstanceType  types.String `tfsdk:"instance_type"`
	InstanceId    types.String `tfsdk:"instance_id"`
	InstanceName  types.String `tfsdk:"instance_name"`
	ZoneId        types.String `tfsdk:"zone_id"`
	ZoneName      types.String `tfsdk:"zone_name"`
	AliasZoneName types.String `tfsdk:"alias_zone_name"`
}

// teoOriginGroupV2RecordAttrTypes returns the attribute types of one records element.
func teoOriginGroupV2RecordAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"record":            types.StringType,
		"type":              types.StringType,
		"record_id":         types.StringType,
		"weight":            types.Int64Type,
		"private":           types.BoolType,
		"private_parameters": types.ListType{ElemType: types.ObjectType{
			AttrTypes: teoOriginGroupV2PrivateParameterAttrTypes(),
		}},
	}
}

// teoOriginGroupV2PrivateParameterAttrTypes returns the attribute types of one
// private_parameters element.
func teoOriginGroupV2PrivateParameterAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"name":  types.StringType,
		"value": types.StringType,
	}
}

// teoOriginGroupV2ReferenceAttrTypes returns the attribute types of one
// references element.
func teoOriginGroupV2ReferenceAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"instance_type":   types.StringType,
		"instance_id":     types.StringType,
		"instance_name":   types.StringType,
		"zone_id":         types.StringType,
		"zone_name":       types.StringType,
		"alias_zone_name": types.StringType,
	}
}

func (r *TeoOriginGroupV2Resource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "tencentcloud_teo_origin_group_v2"
}

func (r *TeoOriginGroupV2Resource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a resource to create a TEO origin group.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Resource identifier in the form `<zone_id>#<group_id>`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"zone_id": schema.StringAttribute{
				Required: true,
				Description: "Site ID.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Description: "Origin group name. It contains 1 to 200 characters, and only supports a - z, A - Z, 0 - 9, _, -. " +
					"When omitted in the configuration, keep the original name: a null value in the configuration is sent as an " +
					"explicit empty string, which keeps the group name unchanged; remove the attribute to keep the current name.",
			},
			"type": schema.StringAttribute{
				Required: true,
				Description: "Origin group type. Valid values: `GENERAL`: universal origin group, which only supports IP/domain origins and can be " +
					"referenced by domain service, rule engine, layer-4 proxy, general load balancing and HTTP-specific load balancing; " +
					"`HTTP`: HTTP-specific origin group, which supports IP/domain and object storage origins and cannot be referenced " +
					"by the layer-4 proxy.",
			},
			"host_header": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Description: "Back-to-origin Host Header, which only takes effect when type is HTTP. The rule engine Host Header " +
					"configuration has a higher priority than the origin group Host Header.",
			},
		"records": schema.ListNestedAttribute{
			Required: true,
			Description: "Origin records of the origin group. The whole list is sent to the API on update, so a change " +
				"replaces the whole record set.",
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"record": schema.StringAttribute{
						Required: true,
						Description: "Origin record value without the port. It can be an IPv4, IPv6 or domain name.",
					},
					"type": schema.StringAttribute{
						Optional: true,
						Computed: true,
						Description: "Origin type. Valid values: `IP_DOMAIN`: IPv4, IPv6 or domain name origin; `COS`: COS origin; " +
							"`AWS_S3`: AWS S3 object storage origin.",
					},
					"record_id": schema.StringAttribute{
						Optional: true,
						Computed: true,
						Description: "Origin record ID, which is server-assigned and back-filled by Read.",
					},
					"weight": schema.Int64Attribute{
						Optional: true,
						Computed: true,
						Description: "Origin weight for traffic scheduling, ranging from 0 to 100 (integer). Leave it unset and the " +
							"system schedules by the default policy; set it to 0 and no traffic is scheduled to this origin (at least " +
							"one origin must have a weight greater than 0). A null value in the configuration means \"not set\".",
					},
					"private": schema.BoolAttribute{
						Optional: true,
						Computed: true,
						Description: "Whether to use private authentication. It takes effect when the origin type is COS or AWS_S3. " +
							"Valid values: `true`: use private authentication; `false`: not use.",
					},
					"private_parameters": schema.ListNestedAttribute{
						Optional: true,
						Description: "Private authentication parameters, which take effect only when private is true.",
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"name": schema.StringAttribute{
									Required: true,
									Description: "Private authentication parameter name. Valid values: `AccessKeyId`: Access Key ID; " +
										"`SecretAccessKey`: Secret Access Key; `SignatureVersion`: signature version, v2 or v4; `Region`: bucket region.",
								},
								"value": schema.StringAttribute{
									Required: true,
									Sensitive: true,
									Description: "Private authentication parameter value.",
								},
							},
						},
					},
				},
			},
		},
			"create_time": schema.StringAttribute{
				Computed: true,
				Description: "Origin group creation time.",
			},
			"update_time": schema.StringAttribute{
				Computed: true,
				Description: "Origin group update time.",
			},
			"references": schema.ListNestedAttribute{
				Computed: true,
				Description: "List of instances referencing the origin group.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"instance_type": schema.StringAttribute{
							Computed: true,
							Description: "Reference service type. Valid values: `acceleration-domain`; `rule-engine`; `load-balancer`; `application-proxy`.",
						},
						"instance_id": schema.StringAttribute{
							Computed: true,
							Description: "Instance ID of the reference type.",
						},
						"instance_name": schema.StringAttribute{
							Computed: true,
							Description: "Instance name of the reference type.",
						},
						"zone_id": schema.StringAttribute{
							Computed: true,
							Description: "Zone ID of the referenced instance.",
						},
						"zone_name": schema.StringAttribute{
							Computed: true,
							Description: "Zone name of the referenced instance.",
						},
						"alias_zone_name": schema.StringAttribute{
							Computed: true,
							Description: "Zone alias of the referenced instance.",
						},
					},
				},
			},
		},
	}
}

// Create handles `terraform apply` for new resources.
func (r *TeoOriginGroupV2Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan TeoOriginGroupV2ResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	service := NewTeoService(r.Client())
	request := service.BuildCreateOriginGroupRequest(ctx, &plan)
	if request == nil {
		resp.Diagnostics.AddError(
			"Error building TEO origin group create request",
			"Failed to convert the plan into a CreateOriginGroup request.",
		)
		return
	}

	originGroupId, err := service.CreateTeoOriginGroupV2(ctx, request)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating TEO origin group",
			fmt.Sprintf("CreateOriginGroup failed: %v", err),
		)
		return
	}

	plan.Id = types.StringValue(strings.Join([]string{plan.ZoneId.ValueString(), originGroupId}, tccommon.FILED_SP))

	// Persist the identity BEFORE any further work so that a later failure
	// still leaves the id in state for `terraform destroy`.
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Read back to align state with the cloud side.
	item, err := service.DescribeTeoOriginGroup(ctx, plan.ZoneId.ValueString(), originGroupId)
	if err != nil {
		resp.Diagnostics.AddError("Error reading TEO origin group after create", err.Error())
		return
	}
	if item == nil {
		resp.Diagnostics.AddError(
			"TEO origin group not readable after create",
			fmt.Sprintf("origin group %s was created but cannot be read back", originGroupId),
		)
		return
	}
	r.refreshModel(ctx, &plan, item)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Read handles `terraform refresh` and post-apply state refresh.
func (r *TeoOriginGroupV2Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state TeoOriginGroupV2ResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	zoneId, originGroupId, ok := splitTeoOriginGroupV2Id(state.Id.ValueString())
	if !ok {
		resp.Diagnostics.AddError(
			"Invalid resource ID",
			fmt.Sprintf("TEO origin group id is broken, id is %q (expected `<zone_id>%s<group_id>`)", state.Id.ValueString(), tccommon.FILED_SP),
		)
		return
	}

	service := NewTeoService(r.Client())
	item, err := service.DescribeTeoOriginGroup(ctx, zoneId, originGroupId)
	if err != nil {
		if strings.Contains(err.Error(), "ResourceNotFound") {
			resp.Diagnostics.AddWarning(
				"Resource not found",
				fmt.Sprintf("TEO origin group %s not found, removing from state.", state.Id.ValueString()),
			)
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading TEO origin group", err.Error())
		return
	}

	// Queried successfully but nothing matched: the group is confirmed gone.
	if item == nil {
		resp.Diagnostics.AddWarning(
			"Resource not found",
			fmt.Sprintf("TEO origin group %s not found, removing from state.", state.Id.ValueString()),
		)
		resp.State.RemoveResource(ctx)
		return
	}

	state.ZoneId = types.StringValue(zoneId)
	r.refreshModel(ctx, &state, item)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update handles in-place modifications.
//
// All the mutable attributes (name, type, host_header, records) are sent in a
// single ModifyOriginGroup call, matching the sibling SDKv2 implementation:
// the API replaces the whole record set on each call, and omitting an
// attribute keeps its previous cloud-side value.
func (r *TeoOriginGroupV2Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan TeoOriginGroupV2ResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state TeoOriginGroupV2ResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	zoneId, originGroupId, ok := splitTeoOriginGroupV2Id(state.Id.ValueString())
	if !ok {
		resp.Diagnostics.AddError(
			"Invalid resource ID",
			fmt.Sprintf("TEO origin group id is broken, id is %q", state.Id.ValueString()),
		)
		return
	}

	service := NewTeoService(r.Client())
	request := service.BuildModifyOriginGroupRequest(ctx, &plan)
	if request == nil {
		// Nothing changed; keep the state and return.
		plan.Id = state.Id
		resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
		return
	}

	request.ZoneId = helper.String(zoneId)
	request.GroupId = helper.String(originGroupId)

	if err := service.ModifyTeoOriginGroupV2(ctx, request); err != nil {
		resp.Diagnostics.AddError(
			"Error updating TEO origin group",
			fmt.Sprintf("ModifyOriginGroup failed: %v", err),
		)
		return
	}

	item, err := service.DescribeTeoOriginGroup(ctx, zoneId, originGroupId)
	if err != nil {
		resp.Diagnostics.AddError("Error reading TEO origin group after update", err.Error())
		return
	}
	if item == nil {
		resp.Diagnostics.AddError(
			"TEO origin group not readable after update",
			fmt.Sprintf("origin group %s was updated but cannot be read back", originGroupId),
		)
		return
	}
	plan.Id = state.Id
	r.refreshModel(ctx, &plan, item)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete handles resource deletion.
func (r *TeoOriginGroupV2Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state TeoOriginGroupV2ResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	zoneId, originGroupId, ok := splitTeoOriginGroupV2Id(state.Id.ValueString())
	if !ok {
		resp.Diagnostics.AddError(
			"Invalid resource ID",
			fmt.Sprintf("TEO origin group id is broken, id is %q", state.Id.ValueString()),
		)
		return
	}

	service := NewTeoService(r.Client())
	if err := service.DeleteTeoOriginGroupById(ctx, zoneId, originGroupId); err != nil {
		resp.Diagnostics.AddError(
			"Error deleting TEO origin group",
			fmt.Sprintf("DeleteOriginGroup failed: %v", err),
		)
		return
	}
}

// ImportState supports `terraform import` using the `<zone_id>#<group_id>`
// form: the composite id is validated and stored as the id attribute.
func (r *TeoOriginGroupV2Resource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if _, _, ok := splitTeoOriginGroupV2Id(req.ID); !ok {
		resp.Diagnostics.AddError(
			"Invalid import ID",
			fmt.Sprintf("Expected import ID in the form `<zone_id>%s<group_id>`, got %q", tccommon.FILED_SP, req.ID),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
}

// splitTeoOriginGroupV2Id parses `<zone_id>#<group_id>`.
func splitTeoOriginGroupV2Id(id string) (string, string, bool) {
	parts := strings.Split(id, tccommon.FILED_SP)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", false
	}
	return parts[0], parts[1], true
}

// refreshModel copies cloud-owned fields and echoed fields into the model,
// keeping user-provided values for fields the cloud side echoes normalized.
func (r *TeoOriginGroupV2Resource) refreshModel(ctx context.Context, m *TeoOriginGroupV2ResourceModel, item *teo.OriginGroup) {
	if item == nil {
		return
	}
	m.Name = types.StringPointerValue(item.Name)
	m.Type = types.StringPointerValue(item.Type)
	m.HostHeader = types.StringPointerValue(item.HostHeader)
	m.CreateTime = types.StringPointerValue(item.CreateTime)
	m.UpdateTime = types.StringPointerValue(item.UpdateTime)

	records := make([]TeoOriginGroupV2RecordModel, 0, len(item.Records))
	for _, record := range item.Records {
		if record == nil {
			continue
		}
		rm := TeoOriginGroupV2RecordModel{
			Record:   types.StringPointerValue(record.Record),
			Type:     types.StringPointerValue(record.Type),
			RecordId: types.StringPointerValue(record.RecordId),
			Private:  types.BoolPointerValue(record.Private),
		}
		if record.Weight != nil {
			rm.Weight = types.Int64Value(int64(*record.Weight))
		}
		privateParameters := make([]TeoOriginGroupV2PrivateParameterModel, 0, len(record.PrivateParameters))
		for _, pp := range record.PrivateParameters {
			if pp == nil {
				continue
			}
			privateParameters = append(privateParameters, TeoOriginGroupV2PrivateParameterModel{
				Name:  types.StringPointerValue(pp.Name),
				Value: types.StringPointerValue(pp.Value),
			})
		}
		ppList, diags := types.ListValueFrom(ctx, types.ObjectType{
			AttrTypes: teoOriginGroupV2PrivateParameterAttrTypes(),
		}, privateParameters)
		_ = diags
		rm.PrivateParameters = ppList
		records = append(records, rm)
	}
	recordsList, diags := types.ListValueFrom(ctx, types.ObjectType{
		AttrTypes: teoOriginGroupV2RecordAttrTypes(),
	}, records)
	_ = diags
	m.Records = recordsList

	references := make([]TeoOriginGroupV2ReferenceModel, 0, len(item.References))
	for _, ref := range item.References {
		if ref == nil {
			continue
		}
		references = append(references, TeoOriginGroupV2ReferenceModel{
			InstanceType:  types.StringPointerValue(ref.InstanceType),
			InstanceId:    types.StringPointerValue(ref.InstanceId),
			InstanceName: types.StringPointerValue(ref.InstanceName),
			ZoneId:        types.StringPointerValue(ref.ZoneId),
			ZoneName:      types.StringPointerValue(ref.ZoneName),
			AliasZoneName: types.StringPointerValue(ref.AliasZoneName),
		})
	}
	refsList, diags := types.ListValueFrom(ctx, types.ObjectType{
		AttrTypes: teoOriginGroupV2ReferenceAttrTypes(),
	}, references)
	_ = diags
	m.References = refsList
}
