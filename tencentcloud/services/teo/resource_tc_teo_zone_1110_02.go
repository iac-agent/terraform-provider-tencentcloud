package teo

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/fw"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

var (
	_ resource.ResourceWithConfigure   = &TeoZone1110_02Resource{}
	_ resource.ResourceWithImportState = &TeoZone1110_02Resource{}
)

// NewTeoZone1110_02Resource is the factory referenced by
// tencentcloud/framework/registry.go to register this resource.
func NewTeoZone1110_02Resource() resource.Resource {
	return &TeoZone1110_02Resource{}
}

// TeoZone1110_02Resource implements resource.Resource for
// tencentcloud_teo_zone_1110_02 (TencentCloud EdgeOne (TEO) zone).
type TeoZone1110_02Resource struct {
	fw.ResourceWithConfigure
	fw.WithImportByID
}

// Nested attribute maps shared by the top-level ownership_verification
// attribute and the read-only zones detail tree.

var teoZone1110_02TagAttributes = map[string]schema.Attribute{
	"tag_key": schema.StringAttribute{
		Required:    true,
		Description: "Tag key.",
	},
	"tag_value": schema.StringAttribute{
		Required:    true,
		Description: "Tag value.",
	},
}

var teoZone1110_02DnsVerificationAttributes = map[string]schema.Attribute{
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
}

var teoZone1110_02FileVerificationAttributes = map[string]schema.Attribute{
	"path": schema.StringAttribute{
		Computed:    true,
		Description: "URL path of the verification file.",
	},
	"content": schema.StringAttribute{
		Computed:    true,
		Description: "Content of the verification file.",
	},
}

var teoZone1110_02NsVerificationAttributes = map[string]schema.Attribute{
	"name_servers": schema.SetAttribute{
		ElementType: types.StringType,
		Computed:    true,
		Description: "DNS servers assigned by EdgeOne for NS access; switch the domain NameServer to these addresses.",
	},
}

var teoZone1110_02OwnershipVerificationAttributes = map[string]schema.Attribute{
	"dns_verification": schema.SingleNestedAttribute{
		Computed:    true,
		Attributes:  teoZone1110_02DnsVerificationAttributes,
		Description: "Information required for DNS resolution verification.",
	},
	"file_verification": schema.SingleNestedAttribute{
		Computed:    true,
		Attributes:  teoZone1110_02FileVerificationAttributes,
		Description: "Information required for file verification.",
	},
	"ns_verification": schema.SingleNestedAttribute{
		Computed:    true,
		Attributes:  teoZone1110_02NsVerificationAttributes,
		Description: "Information required for switching DNS servers in NS access mode.",
	},
}

var teoZone1110_02VanityAttributes = map[string]schema.Attribute{
	"switch": schema.StringAttribute{
		Required:    true,
		Description: "Custom NS switch. Valid values: on, off.",
	},
	"servers": schema.SetAttribute{
		ElementType: types.StringType,
		Optional:    true,
		Description: "List of custom NS servers.",
	},
}

var teoZone1110_02VanityDetailAttributes = map[string]schema.Attribute{
	"switch": schema.StringAttribute{
		Computed:    true,
		Description: "Custom NS switch. Valid values: on, off.",
	},
	"servers": schema.SetAttribute{
		ElementType: types.StringType,
		Computed:    true,
		Description: "List of custom NS servers.",
	},
}

var teoZone1110_02VanityIpAttributes = map[string]schema.Attribute{
	"name": schema.StringAttribute{
		Computed:    true,
		Description: "Custom name server name.",
	},
	"i_pv4": schema.StringAttribute{
		Computed:    true,
		Description: "Custom name server IPv4 address.",
	},
}

var teoZone1110_02SvAttributes = map[string]schema.Attribute{
	"key": schema.StringAttribute{
		Computed:    true,
		Description: "Inquiry parameter key.",
	},
	"value": schema.StringAttribute{
		Computed:    true,
		Description: "Inquiry parameter value.",
	},
	"pack": schema.StringAttribute{
		Computed:    true,
		Description: "Quota mapped by the inquiry parameter. Valid values: zone, custom-rule, rate-limiting-rule, l4-proxy-instance.",
	},
	"instance_id": schema.StringAttribute{
		Computed:    true,
		Description: "Layer-4 proxy instance ID mapped by the inquiry parameter.",
	},
	"protection_specs": schema.StringAttribute{
		Computed:    true,
		Description: "Protection level mapped by the inquiry parameter.",
	},
}

var teoZone1110_02ResourceAttributes = map[string]schema.Attribute{
	"id": schema.StringAttribute{
		Computed:    true,
		Description: "Resource ID.",
	},
	"pay_mode": schema.Int64Attribute{
		Computed:    true,
		Description: "Billing mode. Valid values: 0 (postpaid).",
	},
	"create_time": schema.StringAttribute{
		Computed:    true,
		Description: "Creation time.",
	},
	"enable_time": schema.StringAttribute{
		Computed:    true,
		Description: "Effective time.",
	},
	"expire_time": schema.StringAttribute{
		Computed:    true,
		Description: "Expiration time.",
	},
	"status": schema.StringAttribute{
		Computed:    true,
		Description: "Plan status. Valid values: normal, isolated, destroyed.",
	},
	"sv": schema.ListNestedAttribute{
		Computed:    true,
		NestedObject: schema.NestedAttributeObject{Attributes: teoZone1110_02SvAttributes},
		Description: "Inquiry parameters.",
	},
	"auto_renew_flag": schema.Int64Attribute{
		Computed:    true,
		Description: "Whether auto renewal is enabled. Valid values: 0 (default), 1 (auto renew), 2 (no auto renew).",
	},
	"plan_id": schema.StringAttribute{
		Computed:    true,
		Description: "Associated plan resource ID.",
	},
	"area": schema.StringAttribute{
		Computed:    true,
		Description: "Region. Valid values: mainland, overseas, global.",
	},
	"group": schema.StringAttribute{
		Computed:    true,
		Description: "Resource type. Valid values: plan, pay-as-you-go, value-added.",
	},
	"zone_number": schema.Int64Attribute{
		Computed:    true,
		Description: "Number of zones bound to the resource.",
	},
	"type": schema.StringAttribute{
		Computed:    true,
		Description: "Resource marker type.",
	},
}

var teoZone1110_02NSDetailAttributes = map[string]schema.Attribute{
	"cname_speed_up": schema.StringAttribute{
		Computed:    true,
		Description: "Whether CNAME acceleration is enabled. Valid values: enabled, disabled.",
	},
	"is_fake": schema.Int64Attribute{
		Computed:    true,
		Description: "Whether a zone with the same name exists. Valid values: 0 (no), 1 (yes).",
	},
	"ownership_verification": schema.SingleNestedAttribute{
		Computed:    true,
		Attributes:  teoZone1110_02OwnershipVerificationAttributes,
		Description: "Ownership verification information.",
	},
	"original_name_servers": schema.SetAttribute{
		ElementType: types.StringType,
		Computed:    true,
		Description: "NS servers currently in use detected by EdgeOne.",
	},
	"name_servers": schema.SetAttribute{
		ElementType: types.StringType,
		Computed:    true,
		Description: "NS servers assigned by EdgeOne.",
	},
	"vanity_name_servers": schema.SingleNestedAttribute{
		Computed:    true,
		Attributes:  teoZone1110_02VanityDetailAttributes,
		Description: "Custom NS server information.",
	},
	"vanity_name_servers_ips": schema.ListNestedAttribute{
		Computed:    true,
		NestedObject: schema.NestedAttributeObject{Attributes: teoZone1110_02VanityIpAttributes},
		Description: "IP address information of custom NS servers.",
	},
}

var teoZone1110_02CNAMEDetailAttributes = map[string]schema.Attribute{
	"is_fake": schema.Int64Attribute{
		Computed:    true,
		Description: "Whether the zone is a fake zone. Valid values: 0 (no), 1 (yes).",
	},
	"ownership_verification": schema.SingleNestedAttribute{
		Computed:    true,
		Attributes:  teoZone1110_02OwnershipVerificationAttributes,
		Description: "Ownership verification information.",
	},
}

var teoZone1110_02DNSPodDetailAttributes = map[string]schema.Attribute{
	"is_fake": schema.Int64Attribute{
		Computed:    true,
		Description: "Whether the zone is a fake zone. Valid values: 0 (no), 1 (yes).",
	},
}

var teoZone1110_02WorkModeInfoAttributes = map[string]schema.Attribute{
	"config_group_type": schema.StringAttribute{
		Computed:    true,
		Description: "Configuration group type. Valid values: l7_acceleration, edge_functions, web_security.",
	},
	"work_mode": schema.StringAttribute{
		Computed:    true,
		Description: "Work mode. Valid values: immediate_effect, version_control.",
	},
}

var teoZone1110_02DetailAttributes = map[string]schema.Attribute{
	"zone_id": schema.StringAttribute{
		Computed:    true,
		Description: "Zone ID.",
	},
	"zone_name": schema.StringAttribute{
		Computed:    true,
		Description: "Zone name.",
	},
	"alias_zone_name": schema.StringAttribute{
		Computed:    true,
		Description: "Alias zone name for zones sharing the same name.",
	},
	"area": schema.StringAttribute{
		Computed:    true,
		Description: "Acceleration area of the zone. Valid values: global, mainland, overseas.",
	},
	"type": schema.StringAttribute{
		Computed:    true,
		Description: "Zone access type. Valid values: full, partial, noDomainAccess, dnsPodAccess, pages, ai.",
	},
	"tags": schema.ListNestedAttribute{
		Computed:     true,
		NestedObject: schema.NestedAttributeObject{Attributes: teoZone1110_02TagAttributes},
		Description:  "Tags associated with the zone.",
	},
	"resources": schema.ListNestedAttribute{
		Computed:     true,
		NestedObject: schema.NestedAttributeObject{Attributes: teoZone1110_02ResourceAttributes},
		Description:  "Billing resource list.",
	},
	"ns_detail": schema.SingleNestedAttribute{
		Computed:    true,
		Attributes:  teoZone1110_02NSDetailAttributes,
		Description: "NS access zone details. Returned only when Type = full.",
	},
	"cname_detail": schema.SingleNestedAttribute{
		Computed:    true,
		Attributes:  teoZone1110_02CNAMEDetailAttributes,
		Description: "CNAME access zone details. Returned only when Type = partial.",
	},
	"dns_pod_detail": schema.SingleNestedAttribute{
		Computed:    true,
		Attributes:  teoZone1110_02DNSPodDetailAttributes,
		Description: "DNSPod hosted zone details. Returned only when Type = dnsPodAccess.",
	},
	"created_on": schema.StringAttribute{
		Computed:    true,
		Description: "Zone creation time.",
	},
	"modified_on": schema.StringAttribute{
		Computed:    true,
		Description: "Zone modification time.",
	},
	"status": schema.StringAttribute{
		Computed:    true,
		Description: "Zone status. Valid values: active, pending, moved, deactivated, initializing.",
	},
	"cname_status": schema.StringAttribute{
		Computed:    true,
		Description: "CNAME access status. Valid values: finished, pending.",
	},
	"active_status": schema.StringAttribute{
		Computed:    true,
		Description: "Display status. Valid values: active, inactive, paused.",
	},
	"lock_status": schema.StringAttribute{
		Computed:    true,
		Description: "Lock status. Valid values: enable, disable, plan_migrate.",
	},
	"paused": schema.BoolAttribute{
		Computed:    true,
		Description: "Whether the zone is paused.",
	},
	"is_fake": schema.Int64Attribute{
		Computed:    true,
		Description: "Whether the zone is a fake zone (legacy field, no longer maintained). Valid values: 0 (no), 1 (yes).",
	},
	"cname_speed_up": schema.StringAttribute{
		Computed:    true,
		Description: "Whether CNAME acceleration is enabled (legacy field, no longer maintained). Valid values: enabled, disabled.",
	},
	"ownership_verification": schema.SingleNestedAttribute{
		Computed:    true,
		Attributes:  teoZone1110_02OwnershipVerificationAttributes,
		Description: "Ownership verification information (legacy field, no longer maintained).",
	},
	"original_name_servers": schema.SetAttribute{
		ElementType: types.StringType,
		Computed:    true,
		Description: "NS servers currently used by the zone (legacy field, no longer maintained).",
	},
	"name_servers": schema.SetAttribute{
		ElementType: types.StringType,
		Computed:    true,
		Description: "NS servers assigned by Tencent Cloud (legacy field, no longer maintained).",
	},
	"vanity_name_servers": schema.SingleNestedAttribute{
		Computed:    true,
		Attributes:  teoZone1110_02VanityDetailAttributes,
		Description: "Custom NS information (legacy field, no longer maintained).",
	},
	"vanity_name_servers_ips": schema.ListNestedAttribute{
		Computed:     true,
		NestedObject: schema.NestedAttributeObject{Attributes: teoZone1110_02VanityIpAttributes},
		Description:  "Custom NS IP information (legacy field, no longer maintained).",
	},
	"work_mode_infos": schema.ListNestedAttribute{
		Computed:     true,
		NestedObject: schema.NestedAttributeObject{Attributes: teoZone1110_02WorkModeInfoAttributes},
		Description:  "Work mode of configuration groups for version management.",
	},
}

func (r *TeoZone1110_02Resource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "tencentcloud_teo_zone_1110_02"
}

func (r *TeoZone1110_02Resource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a resource to create a zone in TencentCloud EdgeOne (TEO).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "ID of the zone in the form `zone_id`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"type": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Zone access type. Valid values: partial (CNAME access, default), full (NS access), noDomainAccess, dnsPodAccess (requires the domain hosted in DNSPod), ai (edge inference access).",
			},
			"zone_name": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Zone name. For CNAME/NS access, use the second-level domain (example.com). Leave empty for noDomainAccess access.",
			},
			"area": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Acceleration area of layer-7 domains when Type is partial/full. Valid values: global, mainland, overseas (default). Leave empty for noDomainAccess access.",
			},
			"plan_id": schema.StringAttribute{
				Optional: true,
				Description: "ID of the plan to bind. If omitted, the zone is created in the `init` (inactive) state; bind a plan via BindZoneToPlan to activate it.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"alias_zone_name": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Alias zone name for zones sharing the same name. Only digits, letters, `.`, `-` and `_` are allowed, up to 200 characters.",
			},
			"tags": schema.ListNestedAttribute{
				Optional:     true,
				NestedObject: schema.NestedAttributeObject{Attributes: teoZone1110_02TagAttributes},
				Description:  "Tags used for permission control and billing allocation of the zone.",
				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplace(),
				},
			},
			"allow_duplicates": schema.BoolAttribute{
				Optional:    true,
				Description: "Whether to allow duplicate access. Valid values: true, false (default).",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
			},
			"jump_start": schema.BoolAttribute{
				Optional:    true,
				Description: "Whether to skip scanning existing DNS records of the zone. Default value: false.",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
			},
			"vanity_name_servers": schema.SingleNestedAttribute{
				Optional:    true,
				Computed:    true,
				Attributes:  teoZone1110_02VanityAttributes,
				Description: "Custom NS information to replace the system-assigned name servers. Not allowed for noDomainAccess zones.",
			},
			"ownership_verification": schema.SingleNestedAttribute{
				Computed:    true,
				Attributes:  teoZone1110_02OwnershipVerificationAttributes,
				Description: "Zone ownership verification information returned by CreateZone.",
			},
			"zones": schema.ListNestedAttribute{
				Computed:     true,
				NestedObject: schema.NestedAttributeObject{Attributes: teoZone1110_02DetailAttributes},
				Description:  "Zone detail list returned by DescribeZones. Always contains one element for this resource.",
			},
		},
	}
}

// TeoZone1110_02ResourceModel maps the schema attributes.
type TeoZone1110_02ResourceModel struct {
	Id                    types.String         `tfsdk:"id"`
	Type                  types.String         `tfsdk:"type"`
	ZoneName              types.String         `tfsdk:"zone_name"`
	Area                  types.String         `tfsdk:"area"`
	PlanId                types.String         `tfsdk:"plan_id"`
	AliasZoneName         types.String         `tfsdk:"alias_zone_name"`
	Tags                  []teoZoneTagModel    `tfsdk:"tags"`
	AllowDuplicates       types.Bool           `tfsdk:"allow_duplicates"`
	JumpStart             types.Bool           `tfsdk:"jump_start"`
	VanityNameServers     *teoZoneVanityModel  `tfsdk:"vanity_name_servers"`
	OwnershipVerification types.Object         `tfsdk:"ownership_verification"`
	Zones                 types.List           `tfsdk:"zones"`
}

// Nested model types. Attributes consumed from the cloud are expressed with
// typed values so plan/state round-trips stay consistent.

type teoZoneTagModel struct {
	TagKey   types.String `tfsdk:"tag_key"`
	TagValue types.String `tfsdk:"tag_value"`
}

type teoZoneVanityModel struct {
	Switch  types.String `tfsdk:"switch"`
	Servers types.Set    `tfsdk:"servers"`
}

type teoZoneDnsVerificationModel struct {
	Subdomain   types.String `tfsdk:"subdomain"`
	RecordType  types.String `tfsdk:"record_type"`
	RecordValue types.String `tfsdk:"record_value"`
}

type teoZoneFileVerificationModel struct {
	Path    types.String `tfsdk:"path"`
	Content types.String `tfsdk:"content"`
}

type teoZoneNsVerificationModel struct {
	NameServers types.Set `tfsdk:"name_servers"`
}

type teoZoneOwnershipVerificationModel struct {
	DnsVerification  *teoZoneDnsVerificationModel  `tfsdk:"dns_verification"`
	FileVerification *teoZoneFileVerificationModel `tfsdk:"file_verification"`
	NsVerification   *teoZoneNsVerificationModel   `tfsdk:"ns_verification"`
}

type teoZoneVanityIpModel struct {
	Name types.String `tfsdk:"name"`
	IPv4 types.String `tfsdk:"i_pv4"`
}

type teoZoneSvModel struct {
	Key             types.String `tfsdk:"key"`
	Value           types.String `tfsdk:"value"`
	Pack            types.String `tfsdk:"pack"`
	InstanceId      types.String `tfsdk:"instance_id"`
	ProtectionSpecs types.String `tfsdk:"protection_specs"`
}

type teoZoneResourceModel struct {
	Id            types.String      `tfsdk:"id"`
	PayMode       types.Int64       `tfsdk:"pay_mode"`
	CreateTime    types.String      `tfsdk:"create_time"`
	EnableTime    types.String      `tfsdk:"enable_time"`
	ExpireTime    types.String      `tfsdk:"expire_time"`
	Status        types.String      `tfsdk:"status"`
	Sv            []teoZoneSvModel  `tfsdk:"sv"`
	AutoRenewFlag types.Int64       `tfsdk:"auto_renew_flag"`
	PlanId        types.String      `tfsdk:"plan_id"`
	Area          types.String      `tfsdk:"area"`
	Group         types.String      `tfsdk:"group"`
	ZoneNumber    types.Int64       `tfsdk:"zone_number"`
	Type          types.String      `tfsdk:"type"`
}

type teoZoneNSDetailModel struct {
	CnameSpeedUp          types.String          `tfsdk:"cname_speed_up"`
	IsFake                types.Int64           `tfsdk:"is_fake"`
	OwnershipVerification types.Object          `tfsdk:"ownership_verification"`
	OriginalNameServers   types.Set             `tfsdk:"original_name_servers"`
	NameServers           types.Set             `tfsdk:"name_servers"`
	VanityNameServers     types.Object          `tfsdk:"vanity_name_servers"`
	VanityNameServersIps  []teoZoneVanityIpModel `tfsdk:"vanity_name_servers_ips"`
}

type teoZoneCNAMEDetailModel struct {
	IsFake                types.Int64  `tfsdk:"is_fake"`
	OwnershipVerification types.Object `tfsdk:"ownership_verification"`
}

type teoZoneDNSPodDetailModel struct {
	IsFake types.Int64 `tfsdk:"is_fake"`
}

type teoZoneWorkModeInfoModel struct {
	ConfigGroupType types.String `tfsdk:"config_group_type"`
	WorkMode        types.String `tfsdk:"work_mode"`
}

type teoZoneDetailModel struct {
	ZoneId                types.String             `tfsdk:"zone_id"`
	ZoneName              types.String             `tfsdk:"zone_name"`
	AliasZoneName         types.String             `tfsdk:"alias_zone_name"`
	Area                  types.String             `tfsdk:"area"`
	Type                  types.String             `tfsdk:"type"`
	Tags                  []teoZoneTagModel        `tfsdk:"tags"`
	Resources             []teoZoneResourceModel   `tfsdk:"resources"`
	NSDetail              *teoZoneNSDetailModel    `tfsdk:"ns_detail"`
	CNAMEDetail           *teoZoneCNAMEDetailModel `tfsdk:"cname_detail"`
	DNSPodDetail          *teoZoneDNSPodDetailModel `tfsdk:"dns_pod_detail"`
	CreatedOn             types.String             `tfsdk:"created_on"`
	ModifiedOn            types.String             `tfsdk:"modified_on"`
	Status                types.String             `tfsdk:"status"`
	CnameStatus           types.String             `tfsdk:"cname_status"`
	ActiveStatus          types.String             `tfsdk:"active_status"`
	LockStatus            types.String             `tfsdk:"lock_status"`
	Paused                types.Bool               `tfsdk:"paused"`
	IsFake                types.Int64              `tfsdk:"is_fake"`
	CnameSpeedUp          types.String             `tfsdk:"cname_speed_up"`
	OwnershipVerification types.Object             `tfsdk:"ownership_verification"`
	OriginalNameServers   types.Set                `tfsdk:"original_name_servers"`
	NameServers           types.Set                `tfsdk:"name_servers"`
	VanityNameServers     types.Object             `tfsdk:"vanity_name_servers"`
	VanityNameServersIps  []teoZoneVanityIpModel   `tfsdk:"vanity_name_servers_ips"`
	WorkModeInfos         []teoZoneWorkModeInfoModel `tfsdk:"work_mode_infos"`
}

// Create handles `terraform apply` for new resources.
func (r *TeoZone1110_02Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan TeoZone1110_02ResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	request := buildTeoZone1110_02CreateRequest(&plan)

	response, err := r.Client().UseTeoClient().CreateZoneWithContext(ctx, request)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating teo zone resource",
			fmt.Sprintf("CreateZone failed: %v", err),
		)
		return
	}
	if response == nil || response.Response == nil || response.Response.ZoneId == nil || *response.Response.ZoneId == "" {
		resp.Diagnostics.AddError(
			"Error creating teo zone resource",
			"CreateZone returned an empty response or an empty zone id.",
		)
		return
	}

	// Persist the identity before any read-back so terraform destroy can
	// still clean up the zone when the refresh below fails.
	plan.Id = types.StringValue(*response.Response.ZoneId)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	zone, err := r.readZone(ctx, plan.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading teo zone resource after create", err.Error())
		return
	}
	if zone == nil {
		resp.Diagnostics.AddError(
			"teo zone resource not readable after create",
			fmt.Sprintf("zone %s was created but cannot be read back", plan.Id.ValueString()),
		)
		return
	}

	if err := refreshTeoZone1110_02Model(ctx, &plan, zone); err != nil {
		resp.Diagnostics.AddError("Error refreshing teo zone resource after create", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Read handles `terraform refresh` and post-apply state refresh.
func (r *TeoZone1110_02Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state TeoZone1110_02ResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	zone, err := r.readZone(ctx, state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading teo zone resource", err.Error())
		return
	}
	// Queried successfully but nothing matched: the zone is gone on the cloud
	// side, so drop it from state to plan a recreation.
	if zone == nil {
		resp.Diagnostics.AddWarning(
			"Resource not found",
			fmt.Sprintf("teo zone %s not found, removing from state.", state.Id.ValueString()),
		)
		resp.State.RemoveResource(ctx)
		return
	}

	if err := refreshTeoZone1110_02Model(ctx, &state, zone); err != nil {
		resp.Diagnostics.AddError("Error refreshing teo zone resource", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update handles in-place updates via ModifyZone.
//
// ModifyZone documents "values not passed keep the current configuration",
// so Null plan values are skipped. Fields removed from the configuration are
// absorbed by Optional+Computed (the plan keeps the prior state value) and
// re-sent unchanged, which converges.
func (r *TeoZone1110_02Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan TeoZone1110_02ResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state TeoZone1110_02ResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	request := teov20220901.NewModifyZoneRequest()
	request.ZoneId = state.Id.ValueStringPointer()
	if !plan.Type.IsNull() && !plan.Type.IsUnknown() {
		request.Type = plan.Type.ValueStringPointer()
	}
	if !plan.ZoneName.IsNull() && !plan.ZoneName.IsUnknown() {
		request.ZoneName = plan.ZoneName.ValueStringPointer()
	}
	if !plan.Area.IsNull() && !plan.Area.IsUnknown() {
		request.Area = plan.Area.ValueStringPointer()
	}
	if !plan.AliasZoneName.IsNull() && !plan.AliasZoneName.IsUnknown() {
		request.AliasZoneName = plan.AliasZoneName.ValueStringPointer()
	}
	if plan.VanityNameServers != nil {
		vanity := &teov20220901.VanityNameServers{}
		vanity.Switch = plan.VanityNameServers.Switch.ValueStringPointer()
		if !plan.VanityNameServers.Servers.IsNull() && !plan.VanityNameServers.Servers.IsUnknown() {
			var servers []string
			d := plan.VanityNameServers.Servers.ElementsAs(ctx, &servers, false)
			resp.Diagnostics.Append(d...)
			if resp.Diagnostics.HasError() {
				return
			}
			vanity.Servers = make([]*string, 0, len(servers))
			for _, s := range servers {
				v := s
				vanity.Servers = append(vanity.Servers, &v)
			}
		}
		request.VanityNameServers = vanity
	}

	_, err := r.Client().UseTeoClient().ModifyZoneWithContext(ctx, request)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating teo zone resource",
			fmt.Sprintf("ModifyZone failed: %v", err),
		)
		return
	}

	zone, err := r.readZone(ctx, state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading teo zone resource after update", err.Error())
		return
	}
	if zone == nil {
		resp.Diagnostics.AddError(
			"teo zone resource not readable after update",
			fmt.Sprintf("zone %s was updated but cannot be read back", state.Id.ValueString()),
		)
		return
	}

	plan.Id = state.Id
	if err := refreshTeoZone1110_02Model(ctx, &plan, zone); err != nil {
		resp.Diagnostics.AddError("Error refreshing teo zone resource after update", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete handles `terraform destroy`.
func (r *TeoZone1110_02Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state TeoZone1110_02ResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.deleteZone(ctx, state.Id.ValueString()); err != nil {
		resp.Diagnostics.AddError(
			"Error deleting teo zone resource",
			fmt.Sprintf("DeleteZone failed for %s: %v", state.Id.ValueString(), err),
		)
	}
}

// deleteZone removes the zone by ID.
func (r *TeoZone1110_02Resource) deleteZone(ctx context.Context, zoneId string) error {
	request := teov20220901.NewDeleteZoneRequest()
	request.ZoneId = helper.String(zoneId)

	_, err := r.Client().UseTeoClient().DeleteZoneWithContext(ctx, request)
	return err
}

// readZone fetches the zone by ID through the zone-id filter of DescribeZones.
//
// Two outcomes stay distinguishable:
//   - API failure (nil response / nil Response / error) -> error
//   - queried successfully but nothing matched         -> (nil, nil)
func (r *TeoZone1110_02Resource) readZone(ctx context.Context, zoneId string) (*teov20220901.Zone, error) {
	request := teov20220901.NewDescribeZonesRequest()
	request.Filters = []*teov20220901.AdvancedFilter{
		{
			Name:   helper.String("zone-id"),
			Values: []*string{helper.String(zoneId)},
		},
	}

	response, err := r.Client().UseTeoClient().DescribeZonesWithContext(ctx, request)
	if err != nil {
		return nil, err
	}
	if response == nil || response.Response == nil {
		return nil, fmt.Errorf("DescribeZones returned an empty response for zone %s", zoneId)
	}
	if len(response.Response.Zones) == 0 {
		return nil, nil
	}
	return response.Response.Zones[0], nil
}

// buildTeoZone1110_02CreateRequest maps the plan model onto CreateZone
// request fields. Create-only fields (plan_id, allow_duplicates, jump_start)
// have no ModifyZone counterpart and are guarded by RequiresReplace.
func buildTeoZone1110_02CreateRequest(plan *TeoZone1110_02ResourceModel) *teov20220901.CreateZoneRequest {
	request := teov20220901.NewCreateZoneRequest()
	request.Type = plan.Type.ValueStringPointer()
	request.ZoneName = plan.ZoneName.ValueStringPointer()
	request.Area = plan.Area.ValueStringPointer()
	request.PlanId = plan.PlanId.ValueStringPointer()
	request.AliasZoneName = plan.AliasZoneName.ValueStringPointer()
	request.AllowDuplicates = plan.AllowDuplicates.ValueBoolPointer()
	request.JumpStart = plan.JumpStart.ValueBoolPointer()
	if len(plan.Tags) > 0 {
		request.Tags = make([]*teov20220901.Tag, 0, len(plan.Tags))
		for _, tag := range plan.Tags {
			request.Tags = append(request.Tags, &teov20220901.Tag{
				TagKey:   tag.TagKey.ValueStringPointer(),
				TagValue: tag.TagValue.ValueStringPointer(),
			})
		}
	}
	return request
}

// refreshTeoZone1110_02Model copies cloud-owned fields into the model.
// Create-only fields (plan_id, allow_duplicates, jump_start) are left
// untouched: DescribeZones does not return them, so the plan/state values
// must be preserved instead of being overwritten with null.
func refreshTeoZone1110_02Model(ctx context.Context, m *TeoZone1110_02ResourceModel, zone *teov20220901.Zone) error {
	if zone == nil {
		return nil
	}

	m.Type = types.StringPointerValue(zone.Type)
	m.ZoneName = types.StringPointerValue(zone.ZoneName)
	m.Area = types.StringPointerValue(zone.Area)
	m.AliasZoneName = types.StringPointerValue(zone.AliasZoneName)

	tags := teoZoneTagModelsFromSdk(zone.Tags)
	if tags == nil {
		m.Tags = nil
	} else {
		m.Tags = tags
	}

	if zone.VanityNameServers == nil {
		m.VanityNameServers = nil
	} else {
		servers, serr := teoZoneStringSetFromSdk(ctx, zone.VanityNameServers.Servers)
		if serr != nil {
			return serr
		}
		m.VanityNameServers = &teoZoneVanityModel{
			Switch:  types.StringPointerValue(zone.VanityNameServers.Switch),
			Servers: servers,
		}
	}

	ownership, err := teoZoneOwnershipObjectFromSdk(ctx, zone.OwnershipVerification)
	if err != nil {
		return err
	}
	m.OwnershipVerification = ownership

	detail, err := teoZoneDetailModelFromSdk(ctx, zone)
	if err != nil {
		return err
	}
	zonesElemType := types.ObjectType{AttrTypes: teoZoneAttrTypes(teoZone1110_02DetailAttributes)}
	zones, err := types.ListValueFrom(ctx, zonesElemType, []teoZoneDetailModel{*detail})
	if err != nil {
		return err
	}
	m.Zones = zones
	return nil
}

// teoZoneAttrTypes derives the attr.Type map from a schema attribute map.
func teoZoneAttrTypes(attrs map[string]schema.Attribute) map[string]attr.Type {
	result := make(map[string]attr.Type, len(attrs))
	for name, a := range attrs {
		result[name] = a.GetType()
	}
	return result
}

func teoZoneTagModelsFromSdk(tags []*teov20220901.Tag) []teoZoneTagModel {
	if len(tags) == 0 {
		return nil
	}
	result := make([]teoZoneTagModel, 0, len(tags))
	for _, tag := range tags {
		result = append(result, teoZoneTagModel{
			TagKey:   types.StringPointerValue(tag.TagKey),
			TagValue: types.StringPointerValue(tag.TagValue),
		})
	}
	return result
}

func teoZoneStringSetFromSdk(ctx context.Context, values []*string) (types.Set, error) {
	if len(values) == 0 {
		return types.SetNull(types.StringType), nil
	}
	strs := make([]string, 0, len(values))
	for _, v := range values {
		if v == nil {
			continue
		}
		strs = append(strs, *v)
	}
	return types.SetValueFrom(ctx, types.StringType, strs)
}

func teoZoneOwnershipObjectFromSdk(ctx context.Context, ov *teov20220901.OwnershipVerification) (types.Object, error) {
	if ov == nil {
		return types.ObjectNull(teoZoneAttrTypes(teoZone1110_02OwnershipVerificationAttributes)), nil
	}

	model := teoZoneOwnershipVerificationModel{}
	if ov.DnsVerification != nil {
		model.DnsVerification = &teoZoneDnsVerificationModel{
			Subdomain:   types.StringPointerValue(ov.DnsVerification.Subdomain),
			RecordType:  types.StringPointerValue(ov.DnsVerification.RecordType),
			RecordValue: types.StringPointerValue(ov.DnsVerification.RecordValue),
		}
	}
	if ov.FileVerification != nil {
		model.FileVerification = &teoZoneFileVerificationModel{
			Path:    types.StringPointerValue(ov.FileVerification.Path),
			Content: types.StringPointerValue(ov.FileVerification.Content),
		}
	}
	if ov.NsVerification != nil {
		servers, err := teoZoneStringSetFromSdk(ctx, ov.NsVerification.NameServers)
		if err != nil {
			return types.ObjectNull(teoZoneAttrTypes(teoZone1110_02OwnershipVerificationAttributes)), err
		}
		model.NsVerification = &teoZoneNsVerificationModel{NameServers: servers}
	}
	return types.ObjectValueFrom(ctx, teoZoneAttrTypes(teoZone1110_02OwnershipVerificationAttributes), &model)
}

func teoZoneVanityObjectFromSdk(ctx context.Context, vanity *teov20220901.VanityNameServers) (types.Object, error) {
	if vanity == nil {
		return types.ObjectNull(teoZoneAttrTypes(teoZone1110_02VanityDetailAttributes)), nil
	}
	servers, err := teoZoneStringSetFromSdk(ctx, vanity.Servers)
	if err != nil {
		return types.ObjectNull(teoZoneAttrTypes(teoZone1110_02VanityDetailAttributes)), err
	}
	model := teoZoneVanityModel{
		Switch:  types.StringPointerValue(vanity.Switch),
		Servers: servers,
	}
	return types.ObjectValueFrom(ctx, teoZoneAttrTypes(teoZone1110_02VanityDetailAttributes), &model)
}

func teoZoneDetailModelFromSdk(ctx context.Context, zone *teov20220901.Zone) (*teoZoneDetailModel, error) {
	detail := &teoZoneDetailModel{
		ZoneId:        types.StringPointerValue(zone.ZoneId),
		ZoneName:      types.StringPointerValue(zone.ZoneName),
		AliasZoneName: types.StringPointerValue(zone.AliasZoneName),
		Area:          types.StringPointerValue(zone.Area),
		Type:          types.StringPointerValue(zone.Type),
		Tags:          teoZoneTagModelsFromSdk(zone.Tags),
		CreatedOn:     types.StringPointerValue(zone.CreatedOn),
		ModifiedOn:    types.StringPointerValue(zone.ModifiedOn),
		Status:        types.StringPointerValue(zone.Status),
		CnameStatus:   types.StringPointerValue(zone.CnameStatus),
		ActiveStatus:  types.StringPointerValue(zone.ActiveStatus),
		LockStatus:    types.StringPointerValue(zone.LockStatus),
		Paused:        types.BoolPointerValue(zone.Paused),
		IsFake:        types.Int64PointerValue(zone.IsFake),
		CnameSpeedUp:  types.StringPointerValue(zone.CnameSpeedUp),
	}

	for _, res := range zone.Resources {
		if res == nil {
			continue
		}
		resourceModel := teoZoneResourceModel{
			Id:            types.StringPointerValue(res.Id),
			PayMode:       types.Int64PointerValue(res.PayMode),
			CreateTime:    types.StringPointerValue(res.CreateTime),
			EnableTime:    types.StringPointerValue(res.EnableTime),
			ExpireTime:    types.StringPointerValue(res.ExpireTime),
			Status:        types.StringPointerValue(res.Status),
			AutoRenewFlag: types.Int64PointerValue(res.AutoRenewFlag),
			PlanId:        types.StringPointerValue(res.PlanId),
			Area:          types.StringPointerValue(res.Area),
			Group:         types.StringPointerValue(res.Group),
			ZoneNumber:    types.Int64PointerValue(res.ZoneNumber),
			Type:          types.StringPointerValue(res.Type),
		}
		for _, sv := range res.Sv {
			if sv == nil {
				continue
			}
			resourceModel.Sv = append(resourceModel.Sv, teoZoneSvModel{
				Key:             types.StringPointerValue(sv.Key),
				Value:           types.StringPointerValue(sv.Value),
				Pack:            types.StringPointerValue(sv.Pack),
				InstanceId:      types.StringPointerValue(sv.InstanceId),
				ProtectionSpecs: types.StringPointerValue(sv.ProtectionSpecs),
			})
		}
		detail.Resources = append(detail.Resources, resourceModel)
	}

	ownership, err := teoZoneOwnershipObjectFromSdk(ctx, zone.OwnershipVerification)
	if err != nil {
		return nil, err
	}
	detail.OwnershipVerification = ownership

	originalServers, err := teoZoneStringSetFromSdk(ctx, zone.OriginalNameServers)
	if err != nil {
		return nil, err
	}
	detail.OriginalNameServers = originalServers

	nameServers, err := teoZoneStringSetFromSdk(ctx, zone.NameServers)
	if err != nil {
		return nil, err
	}
	detail.NameServers = nameServers

	vanity, err := teoZoneVanityObjectFromSdk(ctx, zone.VanityNameServers)
	if err != nil {
		return nil, err
	}
	detail.VanityNameServers = vanity

	for _, ip := range zone.VanityNameServersIps {
		if ip == nil {
			continue
		}
		detail.VanityNameServersIps = append(detail.VanityNameServersIps, teoZoneVanityIpModel{
			Name: types.StringPointerValue(ip.Name),
			IPv4: types.StringPointerValue(ip.IPv4),
		})
	}

	for _, info := range zone.WorkModeInfos {
		if info == nil {
			continue
		}
		detail.WorkModeInfos = append(detail.WorkModeInfos, teoZoneWorkModeInfoModel{
			ConfigGroupType: types.StringPointerValue(info.ConfigGroupType),
			WorkMode:        types.StringPointerValue(info.WorkMode),
		})
	}

	if zone.NSDetail != nil {
		ns := &teoZoneNSDetailModel{
			CnameSpeedUp: types.StringPointerValue(zone.NSDetail.CnameSpeedUp),
			IsFake:       types.Int64PointerValue(zone.NSDetail.IsFake),
		}
		ov, err := teoZoneOwnershipObjectFromSdk(ctx, zone.NSDetail.OwnershipVerification)
		if err != nil {
			return nil, err
		}
		ns.OwnershipVerification = ov

		original, err := teoZoneStringSetFromSdk(ctx, zone.NSDetail.OriginalNameServers)
		if err != nil {
			return nil, err
		}
		ns.OriginalNameServers = original

		assigned, err := teoZoneStringSetFromSdk(ctx, zone.NSDetail.NameServers)
		if err != nil {
			return nil, err
		}
		ns.NameServers = assigned

		vanityObj, err := teoZoneVanityObjectFromSdk(ctx, zone.NSDetail.VanityNameServers)
		if err != nil {
			return nil, err
		}
		ns.VanityNameServers = vanityObj

		for _, ip := range zone.NSDetail.VanityNameServersIps {
			if ip == nil {
				continue
			}
			ns.VanityNameServersIps = append(ns.VanityNameServersIps, teoZoneVanityIpModel{
				Name: types.StringPointerValue(ip.Name),
				IPv4: types.StringPointerValue(ip.IPv4),
			})
		}
		detail.NSDetail = ns
	}

	if zone.CNAMEDetail != nil {
		cname := &teoZoneCNAMEDetailModel{
			IsFake: types.Int64PointerValue(zone.CNAMEDetail.IsFake),
		}
		ov, err := teoZoneOwnershipObjectFromSdk(ctx, zone.CNAMEDetail.OwnershipVerification)
		if err != nil {
			return nil, err
		}
		cname.OwnershipVerification = ov
		detail.CNAMEDetail = cname
	}

	if zone.DNSPodDetail != nil {
		detail.DNSPodDetail = &teoZoneDNSPodDetailModel{
			IsFake: types.Int64PointerValue(zone.DNSPodDetail.IsFake),
		}
	}

	return detail, nil
}
