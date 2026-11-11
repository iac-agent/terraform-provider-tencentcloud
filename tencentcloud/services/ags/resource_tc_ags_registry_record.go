package ags

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	sdk "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ags/v20250920"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/fw"
)

// Extended-interface assertions: the base resource.Resource interface is
// checked implicitly by the factory return value.
var (
	_ resource.ResourceWithConfigure   = &AgsRegistryRecordResource{}
	_ resource.ResourceWithImportState = &AgsRegistryRecordResource{}
)

// NewAgsRegistryRecordResource is the factory referenced by
// tencentcloud/framework/registry.go to register this resource.
func NewAgsRegistryRecordResource() resource.Resource {
	return &AgsRegistryRecordResource{}
}

// AgsRegistryRecordResource implements resource.Resource for
// tencentcloud_ags_registry_record.
type AgsRegistryRecordResource struct {
	fw.ResourceWithConfigure
}

// agsMCPSource mirrors the mcp_source single-nested attribute.
type agsMCPSource struct {
	Type        types.String `tfsdk:"type"`
	Descriptors types.String `tfsdk:"descriptors"`
	EndpointURL types.String `tfsdk:"endpoint_url"`
}

// agsAgentSource mirrors the agent_source single-nested attribute.
type agsAgentSource struct {
	Type        types.String `tfsdk:"type"`
	Descriptors types.String `tfsdk:"descriptors"`
	EndpointURL types.String `tfsdk:"endpoint_url"`
}

// agsSkillSource mirrors the skill_source single-nested attribute.
type agsSkillSource struct {
	Type    types.String `tfsdk:"type"`
	SkillMd types.String `tfsdk:"skill_md"`
}

// agsLabelMutation mirrors one element of the label_mutations list attribute.
type agsLabelMutation struct {
	Operation types.String `tfsdk:"operation"`
	Name      types.String `tfsdk:"name"`
	VersionId types.String `tfsdk:"version_id"`
	Reason    types.String `tfsdk:"reason"`
}

// AgsRegistryRecordResourceModel maps the schema attributes.
type AgsRegistryRecordResourceModel struct {
	Id                types.String `tfsdk:"id"`
	RegistryId        types.String `tfsdk:"registry_id"`
	Name              types.String `tfsdk:"name"`
	DescriptorType    types.String `tfsdk:"descriptor_type"`
	Description       types.String `tfsdk:"description"`
	VersionName       types.String `tfsdk:"version_name"`
	MCPSource         types.Object `tfsdk:"mcp_source"`
	AgentSource       types.Object `tfsdk:"agent_source"`
	SkillSource       types.Object `tfsdk:"skill_source"`
	CustomDescriptors types.String `tfsdk:"custom_descriptors"`
	LabelMutations    types.List   `tfsdk:"label_mutations"`
	Reason            types.String `tfsdk:"reason"`
	RecordId          types.String `tfsdk:"record_id"`
	Record            types.Object `tfsdk:"record"`
	Version           types.Object `tfsdk:"version"`
	ResolvedBy        types.String `tfsdk:"resolved_by"`
	ResolvedLabel     types.String `tfsdk:"resolved_label"`
}

// Attribute type maps used to build the computed nested object values.
var agsRecordAttrTypes = map[string]attr.Type{
	"record_id":               types.StringType,
	"registry_id":             types.StringType,
	"name":                    types.StringType,
	"description":             types.StringType,
	"descriptor_type":         types.StringType,
	"lifecycle_status":        types.StringType,
	"app_id":                  types.Int64Type,
	"creator_uin":             types.StringType,
	"creator_sub_account_uin": types.StringType,
	"create_time":             types.StringType,
	"update_time":             types.StringType,
	"version_count":           types.Int64Type,
	"label_set":               types.SetType{ElemType: types.StringType},
}

var agsApprovalActionAttrTypes = map[string]attr.Type{
	"action_id":             types.StringType,
	"action_type":           types.StringType,
	"actor_type":            types.StringType,
	"actor_uin":             types.StringType,
	"actor_sub_account_uin": types.StringType,
	"comment":               types.StringType,
	"create_time":           types.StringType,
	"request_id":            types.StringType,
}

var agsHarnessVersionAttrTypes = map[string]attr.Type{
	"harness_engine":         types.StringType,
	"harness_engine_version": types.StringType,
}

var agsMetadataEntryAttrTypes = map[string]attr.Type{
	"name":  types.StringType,
	"value": types.StringType,
}

var agsVersionAttrTypes = map[string]attr.Type{
	"version_id":                 types.StringType,
	"record_id":                  types.StringType,
	"revision":                   types.Int64Type,
	"status":                     types.StringType,
	"approval_mode":              types.StringType,
	"app_id":                     types.Int64Type,
	"creator_uin":                types.StringType,
	"create_time":                types.StringType,
	"update_time":                types.StringType,
	"version_name":               types.StringType,
	"descriptors":                types.StringType,
	"source_type":                types.StringType,
	"source_config":              types.StringType,
	"content_status":             types.StringType,
	"content_sha256":             types.StringType,
	"content_size_bytes":         types.Int64Type,
	"config_sha256":              types.StringType,
	"creator_sub_account_uin":    types.StringType,
	"approval_actions":           types.ListType{ElemType: types.ObjectType{AttrTypes: agsApprovalActionAttrTypes}},
	"content_ready_time":         types.StringType,
	"change_log":                 types.StringType,
	"label_set":                  types.SetType{ElemType: types.StringType},
	"supported_harness_versions": types.ListType{ElemType: types.ObjectType{AttrTypes: agsHarnessVersionAttrTypes}},
	"metadata":                   types.ListType{ElemType: types.ObjectType{AttrTypes: agsMetadataEntryAttrTypes}},
	"name":                       types.StringType,
	"description":                types.StringType,
	"descriptor_type":            types.StringType,
}

func (r *AgsRegistryRecordResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "tencentcloud_ags_registry_record"
}

func (r *AgsRegistryRecordResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a resource to create a record in an AGS (Agent Governance Service) registry.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Resource identifier in the form `<registry_id>#<record_id>`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"registry_id": schema.StringAttribute{
				Required:    true,
				Description: "ID of the registry the record belongs to.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Record name, length 1..255. Unique per tenant and registry by normalized name (case-insensitive); reusable after soft delete. Cannot be changed after creation.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"descriptor_type": schema.StringAttribute{
				Required:    true,
				Description: "Protocol descriptor type. Valid values: MCP, A2A, AGUI, CUSTOM, AGENT_SKILLS. Immutable after record creation.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Record description, up to 4096 characters. Set an empty string to clear it on update; removing the attribute keeps the previous value.",
			},
			"version_name": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Display name of Revision 1. Only used at creation time.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"mcp_source": schema.SingleNestedAttribute{
				Optional:    true,
				Description: "Source when descriptor_type is MCP; required for MCP and forbidden for other types.",
				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{
						Required:    true,
						Description: "Source type. Valid values: MANUAL (submit MCP descriptors JSON text directly), URL_IMPORT (import from a remote MCP server.json URL).",
					},
					"descriptors": schema.StringAttribute{
						Optional:    true,
						Description: "Required when type is MANUAL: full MCP server.json object as a JSON string. Validation is performed server side.",
					},
					"endpoint_url": schema.StringAttribute{
						Optional:    true,
						Description: "Remote MCP server.json URL (HTTPS). Required when type is URL_IMPORT.",
					},
				},
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.RequiresReplace(),
				},
			},
			"agent_source": schema.SingleNestedAttribute{
				Optional:    true,
				Description: "Source when descriptor_type is A2A or AGUI; required for those types and forbidden for others.",
				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{
						Required:    true,
						Description: "Source type. Valid values: MANUAL (submit Agent descriptors JSON text directly), URL_IMPORT (import from a remote Agent Card or AGUI endpoint).",
					},
					"descriptors": schema.StringAttribute{
						Optional:    true,
						Description: "Required when type is MANUAL: generic JSON object text. Validation is performed server side.",
					},
					"endpoint_url": schema.StringAttribute{
						Optional:    true,
						Description: "Agent Card URL (A2A) or Runtime Endpoint URL (AGUI). Required when type is URL_IMPORT, HTTPS.",
					},
				},
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.RequiresReplace(),
				},
			},
			"skill_source": schema.SingleNestedAttribute{
				Optional:    true,
				Description: "Source when descriptor_type is AGENT_SKILLS; required for AGENT_SKILLS and forbidden for other types.",
				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{
						Required:    true,
						Description: "Source type. Valid values: MANUAL (submit SKILL.md text directly), TAR_PACKAGE (server issues a COS PUT pre-signed URL, the client uploads and the server validates asynchronously).",
					},
					"skill_md": schema.StringAttribute{
						Optional:    true,
						Description: "SKILL.md text; required and non-empty when type is MANUAL, must not be provided when type is TAR_PACKAGE.",
					},
				},
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.RequiresReplace(),
				},
			},
			"custom_descriptors": schema.StringAttribute{
				Optional:    true,
				Description: "Required when descriptor_type is CUSTOM and forbidden for other types. Content must be a JSON object string.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"label_mutations": schema.ListNestedAttribute{
				Optional:    true,
				Description: "Label changes applied in Record update mode, at most 32 entries; label names must be unique within one request.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"operation": schema.StringAttribute{
							Required:    true,
							Description: "Operation type. Valid values: SET (create or move a label), DELETE (delete a custom label; reserved labels stable and latest cannot be deleted).",
						},
						"name": schema.StringAttribute{
							Required:    true,
							Description: "Label name, length 1..63, matching ^[a-z][a-z0-9._-]{0,62}$; stable and latest are reserved.",
						},
						"version_id": schema.StringAttribute{
							Optional:    true,
							Description: "Target version ID. Optional for SET: omitted means unbound (allowed for custom labels, forbidden for stable); forbidden for DELETE.",
						},
						"reason": schema.StringAttribute{
							Optional:    true,
							Description: "Change reason, up to 1024 characters.",
						},
					},
				},
			},
			"reason": schema.StringAttribute{
				Optional:    true,
				Description: "Deletion reason, up to 1024 characters. Only sent by the Delete operation.",
			},
			"record_id": schema.StringAttribute{
				Computed:    true,
				Description: "Record ID, format `rec-` followed by 8 lowercase alphanumeric characters.",
			},
			"record": schema.SingleNestedAttribute{
				Computed:    true,
				Description: "Record metadata and all labels returned by the cloud.",
				Attributes: map[string]schema.Attribute{
					"record_id": schema.StringAttribute{
						Computed:    true,
						Description: "Record ID, format `rec-` followed by 8 lowercase alphanumeric characters.",
					},
					"registry_id": schema.StringAttribute{
						Computed:    true,
						Description: "ID of the owning registry.",
					},
					"name": schema.StringAttribute{
						Computed:    true,
						Description: "Record name; may repeat within one registry.",
					},
					"description": schema.StringAttribute{
						Computed:    true,
						Description: "Record description.",
					},
					"descriptor_type": schema.StringAttribute{
						Computed:    true,
						Description: "Protocol descriptor type; immutable after creation.",
					},
					"lifecycle_status": schema.StringAttribute{
						Computed:    true,
						Description: "Lifecycle status. Valid values: ACTIVE, DELETED (soft-delete tombstone).",
					},
					"app_id": schema.Int64Attribute{
						Computed:    true,
						Description: "AppId of the owning tenant.",
					},
					"creator_uin": schema.StringAttribute{
						Computed:    true,
						Description: "Creator root account UIN.",
					},
					"creator_sub_account_uin": schema.StringAttribute{
						Computed:    true,
						Description: "Creator sub-account UIN; empty when created by the root account.",
					},
					"create_time": schema.StringAttribute{
						Computed:    true,
						Description: "Creation time, ISO 8601 UTC.",
					},
					"update_time": schema.StringAttribute{
						Computed:    true,
						Description: "Last update time.",
					},
					"version_count": schema.Int64Attribute{
						Computed:    true,
						Description: "Number of non-deleted versions under the record.",
					},
					"label_set": schema.SetAttribute{
						Computed:    true,
						ElementType: types.StringType,
						Description: "All label names under the record, including system labels (stable, latest) and custom labels.",
					},
				},
			},
			"version": schema.SingleNestedAttribute{
				Computed:    true,
				Description: "Full version resolved by the Describe operation (defaults to the stable label).",
				Attributes: map[string]schema.Attribute{
					"version_id": schema.StringAttribute{
						Computed:    true,
						Description: "Version ID, format `rv-` followed by 8 lowercase alphanumeric characters.",
					},
					"record_id": schema.StringAttribute{
						Computed:    true,
						Description: "ID of the owning record.",
					},
					"revision": schema.Int64Attribute{
						Computed:    true,
						Description: "Incrementing version number, starting from 1.",
					},
					"status": schema.StringAttribute{
						Computed:    true,
						Description: "Version status.",
					},
					"approval_mode": schema.StringAttribute{
						Computed:    true,
						Description: "Approval mode; locked at creation time.",
					},
					"app_id": schema.Int64Attribute{
						Computed:    true,
						Description: "AppId of the owning tenant.",
					},
					"creator_uin": schema.StringAttribute{
						Computed:    true,
						Description: "Creator root account UIN.",
					},
					"create_time": schema.StringAttribute{
						Computed:    true,
						Description: "Creation time.",
					},
					"update_time": schema.StringAttribute{
						Computed:    true,
						Description: "Last update time.",
					},
					"version_name": schema.StringAttribute{
						Computed:    true,
						Description: "Version display name.",
					},
					"descriptors": schema.StringAttribute{
						Computed:    true,
						Description: "Protocol descriptor object (JSON string).",
					},
					"source_type": schema.StringAttribute{
						Computed:    true,
						Description: "Content source type.",
					},
					"source_config": schema.StringAttribute{
						Computed:    true,
						Description: "Normalized source config object (JSON string).",
					},
					"content_status": schema.StringAttribute{
						Computed:    true,
						Description: "Content status.",
					},
					"content_sha256": schema.StringAttribute{
						Computed:    true,
						Description: "SHA-256 of READY content.",
					},
					"content_size_bytes": schema.Int64Attribute{
						Computed:    true,
						Description: "Size of READY content in bytes.",
					},
					"config_sha256": schema.StringAttribute{
						Computed:    true,
						Description: "SHA-256 of the normalized config content (used for idempotent deduplication).",
					},
					"creator_sub_account_uin": schema.StringAttribute{
						Computed:    true,
						Description: "Creator sub-account UIN.",
					},
					"approval_actions": schema.ListNestedAttribute{
						Computed:    true,
						Description: "Historical approval actions of the version.",
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"action_id": schema.StringAttribute{
									Computed:    true,
									Description: "Action ID.",
								},
								"action_type": schema.StringAttribute{
									Computed:    true,
									Description: "Action type.",
								},
								"actor_type": schema.StringAttribute{
									Computed:    true,
									Description: "Actor type. Valid values: USER, SYSTEM.",
								},
								"actor_uin": schema.StringAttribute{
									Computed:    true,
									Description: "Actor root account UIN.",
								},
								"actor_sub_account_uin": schema.StringAttribute{
									Computed:    true,
									Description: "Actor sub-account UIN.",
								},
								"comment": schema.StringAttribute{
									Computed:    true,
									Description: "Action comment.",
								},
								"create_time": schema.StringAttribute{
									Computed:    true,
									Description: "Occurrence time.",
								},
								"request_id": schema.StringAttribute{
									Computed:    true,
									Description: "RequestId of the corresponding cloud API request.",
								},
							},
						},
					},
					"content_ready_time": schema.StringAttribute{
						Computed:    true,
						Description: "Content ready time, UTC RFC3339. Only returned for TAR_PACKAGE or FILE sources with content_status READY.",
					},
					"change_log": schema.StringAttribute{
						Computed:    true,
						Description: "Change reason of the version, up to 4096 characters; immutable.",
					},
					"label_set": schema.SetAttribute{
						Computed:    true,
						ElementType: types.StringType,
						Description: "Label names currently bound to this version.",
					},
					"supported_harness_versions": schema.ListNestedAttribute{
						Computed:    true,
						Description: "Harness engines and versions supported by this asset version.",
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"harness_engine": schema.StringAttribute{
									Computed:    true,
									Description: "Harness engine identifier, for example pi.",
								},
								"harness_engine_version": schema.StringAttribute{
									Computed:    true,
									Description: "Specific version of the harness engine, for example 1.3.0.",
								},
							},
						},
					},
					"metadata": schema.ListNestedAttribute{
						Computed:    true,
						Description: "Version metadata Name/Value entries, sorted by name ascending.",
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"name": schema.StringAttribute{
									Computed:    true,
									Description: "Metadata name, non-empty lowercase, unique within the array.",
								},
								"value": schema.StringAttribute{
									Computed:    true,
									Description: "Metadata string value.",
								},
							},
						},
					},
					"name": schema.StringAttribute{
						Computed:    true,
						Description: "Current name of the owning record. Only returned by the version list query.",
					},
					"description": schema.StringAttribute{
						Computed:    true,
						Description: "Current description of the owning record. Only returned by the version list query.",
					},
					"descriptor_type": schema.StringAttribute{
						Computed:    true,
						Description: "Asset type of the owning record. Only returned by the version list query.",
					},
				},
			},
			"resolved_by": schema.StringAttribute{
				Computed:    true,
				Description: "Resolution method. Valid values: DEFAULT_STABLE, LABEL, VERSION_ID.",
			},
			"resolved_label": schema.StringAttribute{
				Computed:    true,
				Description: "Label name returned when resolution goes through a label (DEFAULT_STABLE or LABEL).",
			},
		},
	}
}

// Create handles `terraform apply` for new resources.
func (r *AgsRegistryRecordResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan AgsRegistryRecordResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	request := sdk.NewCreateRegistryRecordRequest()
	request.RegistryId = plan.RegistryId.ValueStringPointer()
	request.Name = plan.Name.ValueStringPointer()
	request.DescriptorType = plan.DescriptorType.ValueStringPointer()
	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		request.Description = plan.Description.ValueStringPointer()
	}
	if !plan.VersionName.IsNull() && !plan.VersionName.IsUnknown() {
		request.VersionName = plan.VersionName.ValueStringPointer()
	}
	if !plan.CustomDescriptors.IsNull() && !plan.CustomDescriptors.IsUnknown() {
		request.CustomDescriptors = plan.CustomDescriptors.ValueStringPointer()
	}

	if !plan.MCPSource.IsNull() && !plan.MCPSource.IsUnknown() {
		var m agsMCPSource
		resp.Diagnostics.Append(plan.MCPSource.As(ctx, &m, basetypes.ObjectAsOptions{})...)
		if resp.Diagnostics.HasError() {
			return
		}
		request.MCPSource = &sdk.CloudMCPSourceInput{
			Type:        m.Type.ValueStringPointer(),
			Descriptors: m.Descriptors.ValueStringPointer(),
			EndpointURL: m.EndpointURL.ValueStringPointer(),
		}
	}

	if !plan.AgentSource.IsNull() && !plan.AgentSource.IsUnknown() {
		var a agsAgentSource
		resp.Diagnostics.Append(plan.AgentSource.As(ctx, &a, basetypes.ObjectAsOptions{})...)
		if resp.Diagnostics.HasError() {
			return
		}
		request.AgentSource = &sdk.CloudAgentSourceInput{
			Type:        a.Type.ValueStringPointer(),
			Descriptors: a.Descriptors.ValueStringPointer(),
			EndpointURL: a.EndpointURL.ValueStringPointer(),
		}
	}

	if !plan.SkillSource.IsNull() && !plan.SkillSource.IsUnknown() {
		var s agsSkillSource
		resp.Diagnostics.Append(plan.SkillSource.As(ctx, &s, basetypes.ObjectAsOptions{})...)
		if resp.Diagnostics.HasError() {
			return
		}
		request.SkillSource = &sdk.CloudSkillSourceInput{
			Type:    s.Type.ValueStringPointer(),
			SkillMd: s.SkillMd.ValueStringPointer(),
		}
	}

	response, err := r.Client().UseAgsClient().CreateRegistryRecordWithContext(ctx, request)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating ags registry record",
			fmt.Sprintf("CreateRegistryRecord failed: %v", err),
		)
		return
	}
	if response == nil || response.Response == nil || response.Response.RecordId == nil || *response.Response.RecordId == "" {
		resp.Diagnostics.AddError(
			"Error creating ags registry record",
			"CreateRegistryRecord returned an empty response or an empty record id.",
		)
		return
	}

	// Persist the identity BEFORE any further work so `terraform destroy`
	// can still clean the resource up if later steps fail.
	plan.Id = types.StringValue(plan.RegistryId.ValueString() + "#" + *response.Response.RecordId)
	plan.RecordId = types.StringValue(*response.Response.RecordId)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	describeResp, exists, err := r.read(ctx, plan.RegistryId.ValueString(), plan.RecordId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading ags registry record after create", err.Error())
		return
	}
	if !exists {
		resp.Diagnostics.AddError(
			"ags registry record not readable after create",
			fmt.Sprintf("record %s was created but cannot be read back", plan.RecordId.ValueString()),
		)
		return
	}
	r.refreshModel(&plan, describeResp.Response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Read handles `terraform refresh` and post-apply state refresh.
func (r *AgsRegistryRecordResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state AgsRegistryRecordResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	describeResp, exists, err := r.read(ctx, state.RegistryId.ValueString(), state.RecordId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading ags registry record", err.Error())
		return
	}
	if !exists {
		// Confirmed gone on the cloud side: drop it from state so Terraform
		// plans a recreation.
		resp.Diagnostics.AddWarning(
			"Resource not found",
			fmt.Sprintf("ags registry record %s not found, removing from state.", state.RecordId.ValueString()),
		)
		resp.State.RemoveResource(ctx)
		return
	}

	r.refreshModel(&state, describeResp.Response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update handles in-place updates. Only the Record update mode of
// UpdateRegistryRecord is used: description and label mutations. All other
// request fields of the dual-mode API belong to the Version creation mode,
// which this resource does not model (those attributes force replacement
// instead, so a change recreates the record).
func (r *AgsRegistryRecordResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan AgsRegistryRecordResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state AgsRegistryRecordResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	descriptionChanged := !plan.Description.IsNull() && !plan.Description.IsUnknown() &&
		!plan.Description.Equal(state.Description)
	labelsChanged := !plan.LabelMutations.IsNull() && !plan.LabelMutations.IsUnknown() &&
		!plan.LabelMutations.Equal(state.LabelMutations)

	if !descriptionChanged && !labelsChanged {
		// Only replace-triggering attributes changed; nothing updatable to send.
		return
	}

	request := sdk.NewUpdateRegistryRecordRequest()
	request.RegistryId = state.RegistryId.ValueStringPointer()
	request.RecordId = state.RecordId.ValueStringPointer()

	if descriptionChanged {
		request.Description = plan.Description.ValueStringPointer()
	}

	if labelsChanged {
		elements := plan.LabelMutations.Elements()
		mutations := make([]*sdk.CloudRecordLabelMutation, 0, len(elements))
		for _, elem := range elements {
			obj, ok := elem.(types.Object)
			if !ok {
				resp.Diagnostics.AddError(
					"Error updating ags registry record",
					"label_mutations contains a non-object element.",
				)
				return
			}
			var m agsLabelMutation
			resp.Diagnostics.Append(obj.As(ctx, &m, basetypes.ObjectAsOptions{})...)
			if resp.Diagnostics.HasError() {
				return
			}
			mutations = append(mutations, &sdk.CloudRecordLabelMutation{
				Operation: m.Operation.ValueStringPointer(),
				Name:      m.Name.ValueStringPointer(),
				VersionId: m.VersionId.ValueStringPointer(),
				Reason:    m.Reason.ValueStringPointer(),
			})
		}
		request.LabelMutations = mutations
	}

	_, err := r.Client().UseAgsClient().UpdateRegistryRecordWithContext(ctx, request)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating ags registry record",
			fmt.Sprintf("UpdateRegistryRecord failed: %v", err),
		)
		return
	}

	describeResp, exists, err := r.read(ctx, state.RegistryId.ValueString(), state.RecordId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading ags registry record after update", err.Error())
		return
	}
	if !exists {
		resp.Diagnostics.AddError(
			"ags registry record not readable after update",
			fmt.Sprintf("record %s was updated but cannot be read back", state.RecordId.ValueString()),
		)
		return
	}
	plan.Id = state.Id
	r.refreshModel(&plan, describeResp.Response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete handles resource deletion. VersionId is intentionally never set so
// the whole record (soft delete tombstone) is removed; passing an empty
// string would be rejected by the API as InvalidParameter.VersionId.
func (r *AgsRegistryRecordResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state AgsRegistryRecordResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	request := sdk.NewDeleteRegistryRecordRequest()
	request.RegistryId = state.RegistryId.ValueStringPointer()
	request.RecordId = state.RecordId.ValueStringPointer()
	if !state.Reason.IsNull() && !state.Reason.IsUnknown() {
		request.Reason = state.Reason.ValueStringPointer()
	}

	_, err := r.Client().UseAgsClient().DeleteRegistryRecordWithContext(ctx, request)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting ags registry record",
			fmt.Sprintf("DeleteRegistryRecord failed: %v", err),
		)
		return
	}
}

// ImportState supports `terraform import <addr> <registry_id>#<record_id>`.
func (r *AgsRegistryRecordResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.SplitN(req.ID, "#", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		resp.Diagnostics.AddError(
			"Invalid import id for ags registry record",
			fmt.Sprintf("Expected import id in the form `<registry_id>#<record_id>`, got: %s", req.ID),
		)
		return
	}
	_ = resp.State.SetAttribute(ctx, path.Root("id"), req.ID)
	_ = resp.State.SetAttribute(ctx, path.Root("registry_id"), parts[0])
	_ = resp.State.SetAttribute(ctx, path.Root("record_id"), parts[1])
}

// read queries the record by registry_id + record_id. Three outcomes stay
// distinguishable: API failure (error), queried but not found (exists=false),
// and a complete response (exists=true).
func (r *AgsRegistryRecordResource) read(ctx context.Context, registryId, recordId string) (*sdk.DescribeRegistryRecordResponse, bool, error) {
	request := sdk.NewDescribeRegistryRecordRequest()
	request.RegistryId = &registryId
	request.RecordId = &recordId
	// VersionId and Label are omitted on purpose: the default stable version
	// is the readback target.

	response, err := r.Client().UseAgsClient().DescribeRegistryRecordWithContext(ctx, request)
	if err != nil {
		return nil, false, err
	}
	if response == nil || response.Response == nil {
		return nil, false, fmt.Errorf("DescribeRegistryRecord returned an empty response for %s", recordId)
	}
	if response.Response.Record == nil {
		return nil, false, nil
	}
	return response, true, nil
}

// refreshModel copies cloud-owned fields into the model. Configuration fields
// without a readback path (mcp_source, agent_source, skill_source,
// custom_descriptors, label_mutations, reason) keep their state values.
func (r *AgsRegistryRecordResource) refreshModel(m *AgsRegistryRecordResourceModel, params *sdk.DescribeRegistryRecordResponseParams) {
	if params == nil {
		return
	}
	if rec := params.Record; rec != nil {
		m.RecordId = types.StringPointerValue(rec.RecordId)
		m.RegistryId = types.StringPointerValue(rec.RegistryId)
		m.Name = types.StringPointerValue(rec.Name)
		m.Description = types.StringPointerValue(rec.Description)
		m.DescriptorType = types.StringPointerValue(rec.DescriptorType)
		if rec.RecordId != nil {
			m.Id = types.StringValue(m.RegistryId.ValueString() + "#" + *rec.RecordId)
		}
		m.Record = agsRecordObject(rec)
	}
	if params.Version != nil {
		m.Version = agsVersionObject(params.Version)
		if params.Version.VersionName != nil {
			m.VersionName = types.StringPointerValue(params.Version.VersionName)
		}
	}
	m.ResolvedBy = types.StringPointerValue(params.ResolvedBy)
	m.ResolvedLabel = types.StringPointerValue(params.ResolvedLabel)
}

// agsStringValueSet builds a set of strings from a nullable []*string slice.
func agsStringValueSet(items []*string) types.Set {
	if len(items) == 0 {
		return types.SetNull(types.StringType)
	}
	elems := make([]attr.Value, 0, len(items))
	for _, s := range items {
		elems = append(elems, types.StringPointerValue(s))
	}
	set, _ := types.SetValue(types.StringType, elems)
	return set
}

// agsStringList builds a list of objects from nullable SDK struct slices.
func agsObjectList[T any](items []*T, build func(*T) types.Object, elemTypes map[string]attr.Type) types.List {
	if len(items) == 0 {
		return types.ListNull(types.ObjectType{AttrTypes: elemTypes})
	}
	elems := make([]attr.Value, 0, len(items))
	for _, item := range items {
		if item == nil {
			continue
		}
		elems = append(elems, build(item))
	}
	list, _ := types.ListValue(types.ObjectType{AttrTypes: elemTypes}, elems)
	return list
}

// agsRecordObject converts a CloudRecord into the nested computed object.
func agsRecordObject(rec *sdk.CloudRecord) types.Object {
	if rec == nil {
		return types.ObjectNull(agsRecordAttrTypes)
	}
	obj, _ := types.ObjectValue(agsRecordAttrTypes, map[string]attr.Value{
		"record_id":               types.StringPointerValue(rec.RecordId),
		"registry_id":             types.StringPointerValue(rec.RegistryId),
		"name":                    types.StringPointerValue(rec.Name),
		"description":             types.StringPointerValue(rec.Description),
		"descriptor_type":         types.StringPointerValue(rec.DescriptorType),
		"lifecycle_status":        types.StringPointerValue(rec.LifecycleStatus),
		"app_id":                  types.Int64PointerValue(rec.AppId),
		"creator_uin":             types.StringPointerValue(rec.CreatorUin),
		"creator_sub_account_uin": types.StringPointerValue(rec.CreatorSubAccountUin),
		"create_time":             types.StringPointerValue(rec.CreateTime),
		"update_time":             types.StringPointerValue(rec.UpdateTime),
		"version_count":           types.Int64PointerValue(rec.VersionCount),
		"label_set":               agsStringValueSet(rec.LabelSet),
	})
	return obj
}

// agsVersionObject converts a CloudRecordVersion into the nested computed object.
func agsVersionObject(v *sdk.CloudRecordVersion) types.Object {
	if v == nil {
		return types.ObjectNull(agsVersionAttrTypes)
	}
	obj, _ := types.ObjectValue(agsVersionAttrTypes, map[string]attr.Value{
		"version_id":              types.StringPointerValue(v.VersionId),
		"record_id":               types.StringPointerValue(v.RecordId),
		"revision":                types.Int64PointerValue(v.Revision),
		"status":                  types.StringPointerValue(v.Status),
		"approval_mode":           types.StringPointerValue(v.ApprovalMode),
		"app_id":                  types.Int64PointerValue(v.AppId),
		"creator_uin":             types.StringPointerValue(v.CreatorUin),
		"create_time":             types.StringPointerValue(v.CreateTime),
		"update_time":             types.StringPointerValue(v.UpdateTime),
		"version_name":            types.StringPointerValue(v.VersionName),
		"descriptors":             types.StringPointerValue(v.Descriptors),
		"source_type":             types.StringPointerValue(v.SourceType),
		"source_config":           types.StringPointerValue(v.SourceConfig),
		"content_status":          types.StringPointerValue(v.ContentStatus),
		"content_sha256":          types.StringPointerValue(v.ContentSHA256),
		"content_size_bytes":      types.Int64PointerValue(v.ContentSizeBytes),
		"config_sha256":           types.StringPointerValue(v.ConfigSHA256),
		"creator_sub_account_uin": types.StringPointerValue(v.CreatorSubAccountUin),
		"approval_actions": agsObjectList(v.ApprovalActions, func(a *sdk.CloudVersionApprovalAction) types.Object {
			obj, _ := types.ObjectValue(agsApprovalActionAttrTypes, map[string]attr.Value{
				"action_id":             types.StringPointerValue(a.ActionId),
				"action_type":           types.StringPointerValue(a.ActionType),
				"actor_type":            types.StringPointerValue(a.ActorType),
				"actor_uin":             types.StringPointerValue(a.ActorUin),
				"actor_sub_account_uin": types.StringPointerValue(a.ActorSubAccountUin),
				"comment":               types.StringPointerValue(a.Comment),
				"create_time":           types.StringPointerValue(a.CreateTime),
				"request_id":            types.StringPointerValue(a.RequestId),
			})
			return obj
		}, agsApprovalActionAttrTypes),
		"content_ready_time": types.StringPointerValue(v.ContentReadyTime),
		"change_log":         types.StringPointerValue(v.ChangeLog),
		"label_set":          agsStringValueSet(v.LabelSet),
		"supported_harness_versions": agsObjectList(v.SupportedHarnessVersions, func(h *sdk.HarnessVersionReference) types.Object {
			obj, _ := types.ObjectValue(agsHarnessVersionAttrTypes, map[string]attr.Value{
				"harness_engine":         types.StringPointerValue(h.HarnessEngine),
				"harness_engine_version": types.StringPointerValue(h.HarnessEngineVersion),
			})
			return obj
		}, agsHarnessVersionAttrTypes),
		"metadata": agsObjectList(v.Metadata, func(md *sdk.MetadataEntry) types.Object {
			obj, _ := types.ObjectValue(agsMetadataEntryAttrTypes, map[string]attr.Value{
				"name":  types.StringPointerValue(md.Name),
				"value": types.StringPointerValue(md.Value),
			})
			return obj
		}, agsMetadataEntryAttrTypes),
		"name":            types.StringPointerValue(v.Name),
		"description":     types.StringPointerValue(v.Description),
		"descriptor_type": types.StringPointerValue(v.DescriptorType),
	})
	return obj
}
