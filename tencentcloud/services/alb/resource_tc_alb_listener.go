package alb

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	albv20251030 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/alb/v20251030"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/fw"
	fwhelper "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

var (
	_ resource.ResourceWithConfigure   = &AlbListenerResource{}
	_ resource.ResourceWithImportState = &AlbListenerResource{}

	errAlbListenerNotFound = errors.New("alb listener not found")
)

func NewAlbListenerResource() resource.Resource {
	return &AlbListenerResource{}
}

type AlbListenerResource struct {
	fw.ResourceWithConfigure
	fw.WithImportByID
}

type AlbListenerResourceModel struct {
	Id                   types.String `tfsdk:"id"`
	LoadBalancerId       types.String `tfsdk:"load_balancer_id"`
	ListenerPort         types.Int64  `tfsdk:"listener_port"`
	ListenerProtocol     types.String `tfsdk:"listener_protocol"`
	ListenerName         types.String `tfsdk:"listener_name"`
	DefaultActions       types.List   `tfsdk:"default_actions"`
	CertificateIds       types.List   `tfsdk:"certificate_ids"`
	CaEnabled            types.Bool   `tfsdk:"ca_enabled"`
	CaCertificateIds     types.List   `tfsdk:"ca_certificate_ids"`
	SecurityPolicyId     types.String `tfsdk:"security_policy_id"`
	GzipEnabled          types.Bool   `tfsdk:"gzip_enabled"`
	Http2Enabled         types.Bool   `tfsdk:"http2enabled"`
	IdleTimeout          types.Int64  `tfsdk:"idle_timeout"`
	RequestTimeout       types.Int64  `tfsdk:"request_timeout"`
	XForwardedForConfig  types.Object `tfsdk:"x_forwarded_for_config"`
	Tags                 types.Map    `tfsdk:"tags"`
	ListenerStatus       types.String `tfsdk:"listener_status"`
	CreateTime           types.String `tfsdk:"create_time"`
	ModifyTime           types.String `tfsdk:"modify_time"`
}

type albDefaultActionModel struct {
	Type              types.String `tfsdk:"type"`
	TargetGroupConfig types.Object `tfsdk:"target_group_config"`
}

type albTargetGroupConfigModel struct {
	TargetGroups             types.List   `tfsdk:"target_groups"`
	TargetGroupStickySession types.Object `tfsdk:"target_group_sticky_session"`
}

type albTargetGroupModel struct {
	TargetGroupId types.String `tfsdk:"target_group_id"`
	Weight        types.Int64  `tfsdk:"weight"`
}

type albTargetGroupStickySessionModel struct {
	Enabled types.Bool  `tfsdk:"enabled"`
	Timeout types.Int64 `tfsdk:"timeout"`
}

type albXForwardedForConfigModel struct {
	XForwardedForMode                 types.String `tfsdk:"x_forwarded_for_mode"`
	XForwardedForHostEnabled          types.Bool   `tfsdk:"x_forwarded_for_host_enabled"`
	XForwardedForPortEnabled          types.Bool   `tfsdk:"x_forwarded_for_port_enabled"`
	XForwardedForProtoEnabled         types.Bool   `tfsdk:"x_forwarded_for_proto_enabled"`
	XForwardedForAlbIdEnabled         types.Bool   `tfsdk:"x_forwarded_for_alb_id_enabled"`
	XForwardedForClientSrcPortEnabled types.Bool   `tfsdk:"x_forwarded_for_client_src_port_enabled"`
	XTencentClientVerifyEnabled       types.Bool   `tfsdk:"x_tencent_client_verify_enabled"`
	XTencentClientSDNEnabled          types.Bool   `tfsdk:"x_tencent_client_sdn_enabled"`
	XTencentClientIDNEnabled          types.Bool   `tfsdk:"x_tencent_client_idn_enabled"`
	XTencentClientSerialEnabled       types.Bool   `tfsdk:"x_tencent_client_serial_enabled"`
}

func (r *AlbListenerResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "tencentcloud_alb_listener"
}

func xForwardedForAttributeTypes(ctx context.Context) map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"x_forwarded_for_mode": schema.StringAttribute{
			Optional:    true,
			Computed:    true,
			Description: "Specifies how to handle the X-Forwarded-For (XFF) HTTP header field. Valid values: `append` (default), `remove`, `passthrough`.",
		},
		"x_forwarded_for_host_enabled": schema.BoolAttribute{
			Optional:    true,
			Computed:    true,
			Description: "Whether to obtain the client domain of the load balancer through the X-Forwarded-Host header field.",
		},
		"x_forwarded_for_port_enabled": schema.BoolAttribute{
			Optional:    true,
			Computed:    true,
			Description: "Whether to obtain the listener port of the load balancer through the X-Forwarded-Port header field.",
		},
		"x_forwarded_for_proto_enabled": schema.BoolAttribute{
			Optional:    true,
			Computed:    true,
			Description: "Whether to obtain the listener protocol of the load balancer through the X-Forwarded-Proto header field.",
		},
		"x_forwarded_for_alb_id_enabled": schema.BoolAttribute{
			Optional:    true,
			Computed:    true,
			Description: "Whether to obtain the load balancer instance ID through the ALB-ID header field.",
		},
		"x_forwarded_for_client_src_port_enabled": schema.BoolAttribute{
			Optional:    true,
			Computed:    true,
			Description: "Whether to obtain the client port through the X-Forwarded-Client-srcport header field.",
		},
		"x_tencent_client_verify_enabled": schema.BoolAttribute{
			Optional:    true,
			Computed:    true,
			Description: "Whether to access the client certificate verification result $ssl_client_verify through the X-Tencent-Client-Verify header.",
		},
		"x_tencent_client_sdn_enabled": schema.BoolAttribute{
			Optional:    true,
			Computed:    true,
			Description: "Whether to access the subject $ssl_client_s_dn of the client certificate through the X-Tencent-Client-SDN header.",
		},
		"x_tencent_client_idn_enabled": schema.BoolAttribute{
			Optional:    true,
			Computed:    true,
			Description: "Whether to access the issuer $ssl_client_i_dn of the client certificate through the X-Tencent-Client-IDN header.",
		},
		"x_tencent_client_serial_enabled": schema.BoolAttribute{
			Optional:    true,
			Computed:    true,
			Description: "Whether to access the serial number $ssl_client_serial of the client certificate through the X-Tencent-Client-Serial header.",
		},
	}
}

func (r *AlbListenerResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a resource to create an ALB (Application Load Balancer) listener.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "ID of the resource, in the format `load_balancer_id#listener_id`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"load_balancer_id": schema.StringAttribute{
				Required:    true,
				Description: "ID of the load balancer instance, in the format of `alb-` followed by 8 alphanumeric characters. Immutable: changing it forces resource replacement.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"listener_port": schema.Int64Attribute{
				Required:    true,
				Description: "Port used by the front end of the load balancer instance. Valid value range: 1~65535. Immutable: changing it forces resource replacement.",
				PlanModifiers: []planmodifier.Int64{
					int64RequiresReplace(),
				},
			},
			"listener_protocol": schema.StringAttribute{
				Required:    true,
				Description: "Listening protocol. Valid values: `HTTP`, `HTTPS`. Immutable: changing it forces resource replacement.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"listener_name": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Custom listener name. Length is 1~255 characters. It can contain Chinese characters, letters, digits, `-`, `/`, `.`, and `_`.",
			},
			"default_actions": schema.ListNestedAttribute{
				Required:    true,
				Description: "List of default forwarding rule actions. Currently only 1 default forwarding rule action is supported per listener.",
				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplace(),
				},
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"type": schema.StringAttribute{
							Required:    true,
							Description: "Forwarding action type. Only forwarding to target group is supported when creating a listener.",
						},
						"target_group_config": schema.SingleNestedAttribute{
							Required:    true,
							Description: "Forwarding target group configuration. Only a single target group is supported.",
							Attributes: map[string]schema.Attribute{
								"target_groups": schema.ListNestedAttribute{
									Required:    true,
									Description: "List of target groups.",
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{
											"target_group_id": schema.StringAttribute{
												Required:    true,
												Description: "Target group ID, in the format of `lbtg-` followed by 8 alphanumeric characters.",
											},
											"weight": schema.Int64Attribute{
												Optional:    true,
												Computed:    true,
												Description: "Weight. Valid value range: [0, 100]. Default to 10.",
											},
										},
									},
								},
								"target_group_sticky_session": schema.SingleNestedAttribute{
									Optional:    true,
									Computed:    true,
									Description: "Sticky session between target groups.",
									Attributes: map[string]schema.Attribute{
										"enabled": schema.BoolAttribute{
											Required:    true,
											Description: "Whether to enable sticky session. Default to false.",
										},
										"timeout": schema.Int64Attribute{
											Optional:    true,
											Computed:    true,
											Description: "Timeout in seconds. Valid value range: 1~86400. Default to 1000.",
										},
									},
								},
							},
						},
					},
				},
			},
			"certificate_ids": schema.ListAttribute{
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				Description: "List of server certificate IDs.",
			},
			"ca_enabled": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Whether to enable mutual authentication. Default to false.",
			},
			"ca_certificate_ids": schema.ListAttribute{
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				Description: "List of CA certificate IDs of the listener. Only 1 CA certificate is supported currently. Required when `ca_enabled` is true.",
			},
			"security_policy_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Security policy ID, in the format of `tls-` followed by 8 alphanumeric characters.",
			},
			"gzip_enabled": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Whether to enable Gzip compression. Default to true.",
			},
			"http2enabled": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Whether to enable HTTP/2. HTTP protocol defaults to false, HTTPS protocol defaults to true. Only HTTPS protocol supports this parameter.",
			},
			"idle_timeout": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "Connection idle timeout in seconds. Valid value range: 1~600. Default to 15.",
			},
			"request_timeout": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "Connection request timeout in seconds. Valid value range: 1~600. Default to 60.",
			},
			"x_forwarded_for_config": schema.SingleNestedAttribute{
				Optional:    true,
				Computed:    true,
				Description: "X-Forwarded-For configuration.",
				Attributes:  xForwardedForAttributeTypes(context.Background()),
			},
			"tags": schema.MapAttribute{
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				Description: "Tag list. Up to 20 tags are supported.",
			},
			"listener_status": schema.StringAttribute{
				Computed:    true,
				Description: "Listener status. Valid values: `Active`, `Provisioning`, `Configuring`, `ProvisionFailed`.",
			},
			"create_time": schema.StringAttribute{
				Computed:    true,
				Description: "Creation time of the listener in ISO 8601 format.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"modify_time": schema.StringAttribute{
				Computed:    true,
				Description: "Last modified time of the listener in ISO 8601 format.",
			},
		},
	}
}

func int64RequiresReplace() planmodifier.Int64 {
	return albInt64RequiresReplaceModifier{}
}

type albInt64RequiresReplaceModifier struct{}

func (m albInt64RequiresReplaceModifier) Description(_ context.Context) string {
	return "If the value of this attribute changes, Terraform will destroy and recreate the resource."
}

func (m albInt64RequiresReplaceModifier) MarkdownDescription(_ context.Context) string {
	return "If the value of this attribute changes, Terraform will destroy and recreate the resource."
}

func (m albInt64RequiresReplaceModifier) PlanModifyInt64(_ context.Context, req planmodifier.Int64Request, resp *planmodifier.Int64Response) {
	if req.StateValue.IsNull() || req.PlanValue.IsNull() {
		return
	}
	if !req.StateValue.Equal(req.PlanValue) {
		resp.RequiresReplace = true
	}
}

func (r *AlbListenerResource) service(diags *diag.Diagnostics) (AlbService, bool) {
	if r.Client() == nil {
		diags.AddError(
			"Provider not configured",
			"The provider client is not available. This is unexpected; please report it.",
		)
		return AlbService{}, false
	}
	return NewAlbService(r.Client()), true
}

func (r *AlbListenerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan AlbListenerResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	service, ok := r.service(&resp.Diagnostics)
	if !ok {
		return
	}

	loadBalancerId := plan.LoadBalancerId.ValueString()

	request := albv20251030.NewCreateListenerRequest()
	request.LoadBalancerId = helper.String(loadBalancerId)
	request.ListenerPort = helper.Uint64(uint64(plan.ListenerPort.ValueInt64()))
	request.ListenerProtocol = fwhelper.StringPointerFromValue(plan.ListenerProtocol)
	request.ListenerName = fwhelper.StringPointerFromValue(plan.ListenerName)
	request.GzipEnabled = fwhelper.BoolPointerFromValue(plan.GzipEnabled)
	request.Http2Enabled = fwhelper.BoolPointerFromValue(plan.Http2Enabled)
	request.IdleTimeout = uint64PointerFromInt64(plan.IdleTimeout)
	request.RequestTimeout = uint64PointerFromInt64(plan.RequestTimeout)
	request.SecurityPolicyId = fwhelper.StringPointerFromValue(plan.SecurityPolicyId)
	request.CaEnabled = fwhelper.BoolPointerFromValue(plan.CaEnabled)

	defaultActions, actionDiags := expandDefaultActions(ctx, plan.DefaultActions)
	resp.Diagnostics.Append(actionDiags...)
	if resp.Diagnostics.HasError() {
		return
	}
	request.DefaultActions = defaultActions

	if !plan.CertificateIds.IsNull() && !plan.CertificateIds.IsUnknown() {
		var certificateIds []string
		resp.Diagnostics.Append(plan.CertificateIds.ElementsAs(ctx, &certificateIds, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		request.CertificateIds = helper.Strings(certificateIds)
	}

	if !plan.CaCertificateIds.IsNull() && !plan.CaCertificateIds.IsUnknown() {
		var caCertificateIds []string
		resp.Diagnostics.Append(plan.CaCertificateIds.ElementsAs(ctx, &caCertificateIds, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		request.CaCertificateIds = helper.Strings(caCertificateIds)
	}

	if !plan.Tags.IsNull() && !plan.Tags.IsUnknown() {
		request.Tags = expandTags(ctx, plan.Tags)
	}

	request.ClientToken = helper.String(uuid.New().String())

	listenerId, err := service.CreateListener(ctx, request)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating ALB listener",
			fmt.Sprintf("CreateListener failed (load_balancer_id=%q): %s", loadBalancerId, err),
		)
		return
	}

	plan.Id = types.StringValue(strings.Join([]string{loadBalancerId, listenerId}, tccommon.FILED_SP))

	item, err := describeAlbListener(ctx, service, loadBalancerId, listenerId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading ALB listener after create",
			fmt.Sprintf("CreateListener succeeded but the listener cannot be read back (%s). "+
				"Import it with `terraform import tencentcloud_alb_listener.<name> %s` to keep it managed: %s",
				albListenerScope(loadBalancerId, listenerId), plan.Id.ValueString(), err),
		)
		return
	}

	resp.Diagnostics.Append(refreshAlbListenerModel(ctx, &plan, loadBalancerId, listenerId, item)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *AlbListenerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state AlbListenerResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	loadBalancerId, listenerId, ok := parseAlbListenerId(&resp.Diagnostics, state.Id.ValueString())
	if !ok {
		return
	}

	service, ok := r.service(&resp.Diagnostics)
	if !ok {
		return
	}

	item, err := describeAlbListener(ctx, service, loadBalancerId, listenerId)
	if err != nil {
		if errors.Is(err, errAlbListenerNotFound) {
			resp.Diagnostics.AddWarning(
				"Resource not found",
				fmt.Sprintf("ALB listener %s no longer exists, removing it from state.", state.Id.ValueString()),
			)
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(
			"Error reading ALB listener",
			fmt.Sprintf("DescribeListeners failed (%s): %s", albListenerScope(loadBalancerId, listenerId), err),
		)
		return
	}

	resp.Diagnostics.Append(refreshAlbListenerModel(ctx, &state, loadBalancerId, listenerId, item)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *AlbListenerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan AlbListenerResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state AlbListenerResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	loadBalancerId, listenerId, ok := parseAlbListenerId(&resp.Diagnostics, state.Id.ValueString())
	if !ok {
		return
	}

	service, ok := r.service(&resp.Diagnostics)
	if !ok {
		return
	}

	request := albv20251030.NewModifyListenerAttributesRequest()
	request.LoadBalancerId = helper.String(loadBalancerId)
	request.ListenerId = helper.String(listenerId)
	request.ListenerName = fwhelper.StringPointerFromValue(plan.ListenerName)
	request.GzipEnabled = fwhelper.BoolPointerFromValue(plan.GzipEnabled)
	request.Http2Enabled = fwhelper.BoolPointerFromValue(plan.Http2Enabled)
	request.IdleTimeout = uint64PointerFromInt64(plan.IdleTimeout)
	request.RequestTimeout = uint64PointerFromInt64(plan.RequestTimeout)
	request.SecurityPolicyId = fwhelper.StringPointerFromValue(plan.SecurityPolicyId)
	request.CaEnabled = fwhelper.BoolPointerFromValue(plan.CaEnabled)

	defaultActions, actionDiags := expandDefaultActions(ctx, plan.DefaultActions)
	resp.Diagnostics.Append(actionDiags...)
	if resp.Diagnostics.HasError() {
		return
	}
	request.DefaultActions = defaultActions

	if !plan.CertificateIds.IsNull() && !plan.CertificateIds.IsUnknown() {
		var certificateIds []string
		resp.Diagnostics.Append(plan.CertificateIds.ElementsAs(ctx, &certificateIds, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		request.CertificateIds = helper.Strings(certificateIds)
	}

	if !plan.CaCertificateIds.IsNull() && !plan.CaCertificateIds.IsUnknown() {
		var caCertificateIds []string
		resp.Diagnostics.Append(plan.CaCertificateIds.ElementsAs(ctx, &caCertificateIds, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		request.CaCertificateIds = helper.Strings(caCertificateIds)
	}

	xffConfig, xffDiags := expandXForwardedForConfig(ctx, plan.XForwardedForConfig)
	resp.Diagnostics.Append(xffDiags...)
	if resp.Diagnostics.HasError() {
		return
	}
	request.XForwardedForConfig = xffConfig

	if !plan.Tags.IsNull() && !plan.Tags.IsUnknown() {
		request.ClientToken = helper.String(uuid.New().String())
	}

	if err := service.ModifyListenerAttributes(ctx, request); err != nil {
		resp.Diagnostics.AddError(
			"Error updating ALB listener",
			fmt.Sprintf("ModifyListenerAttributes failed (%s): %s", albListenerScope(loadBalancerId, listenerId), err),
		)
		return
	}

	item, err := describeAlbListener(ctx, service, loadBalancerId, listenerId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading ALB listener after update",
			fmt.Sprintf("ModifyListenerAttributes succeeded but the listener cannot be read back (%s): %s", albListenerScope(loadBalancerId, listenerId), err),
		)
		return
	}

	resp.Diagnostics.Append(refreshAlbListenerModel(ctx, &plan, loadBalancerId, listenerId, item)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *AlbListenerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state AlbListenerResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	loadBalancerId, listenerId, ok := parseAlbListenerId(&resp.Diagnostics, state.Id.ValueString())
	if !ok {
		return
	}

	service, ok := r.service(&resp.Diagnostics)
	if !ok {
		return
	}

	request := albv20251030.NewDeleteListenerRequest()
	request.LoadBalancerId = helper.String(loadBalancerId)
	request.ListenerIds = []*string{helper.String(listenerId)}
	request.ClientToken = helper.String(uuid.New().String())

	if err := service.DeleteListener(ctx, request); err != nil {
		resp.Diagnostics.AddError(
			"Error deleting ALB listener",
			fmt.Sprintf("DeleteListener failed (%s): %s", albListenerScope(loadBalancerId, listenerId), err),
		)
		return
	}
}

func albListenerScope(loadBalancerId, listenerId string) string {
	if listenerId == "" {
		return fmt.Sprintf("load_balancer_id=%q", loadBalancerId)
	}
	return fmt.Sprintf("load_balancer_id=%q, listener_id=%q", loadBalancerId, listenerId)
}

func parseAlbListenerId(diags *diag.Diagnostics, id string) (loadBalancerId, listenerId string, ok bool) {
	parts := strings.Split(id, tccommon.FILED_SP)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		diags.AddError(
			"Invalid resource ID",
			fmt.Sprintf("Resource id %q is invalid, expected `<load_balancer_id>%s<listener_id>`.", id, tccommon.FILED_SP),
		)
		return "", "", false
	}
	return parts[0], parts[1], true
}

func describeAlbListener(ctx context.Context, service AlbService, loadBalancerId, listenerId string) (*albv20251030.ListenerOutput, error) {
	item, err := service.DescribeListenerById(ctx, loadBalancerId, listenerId)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, errAlbListenerNotFound
	}
	return item, nil
}

func uint64PointerFromInt64(v types.Int64) *uint64 {
	if v.IsNull() || v.IsUnknown() {
		return nil
	}
	return helper.Uint64(uint64(v.ValueInt64()))
}

func expandTags(ctx context.Context, tags types.Map) []*albv20251030.TagInfo {
	if tags.IsNull() || tags.IsUnknown() {
		return nil
	}
	tagMap := tags.Elements()
	out := make([]*albv20251030.TagInfo, 0, len(tagMap))
	for key, value := range tagMap {
		valueStr, ok := value.(types.String)
		if !ok {
			continue
		}
		out = append(out, &albv20251030.TagInfo{
			TagKey:   helper.String(key),
			TagValue: fwhelper.StringPointerFromValue(valueStr),
		})
	}
	return out
}

func flattenTags(ctx context.Context, tags []*albv20251030.TagInfo) (types.Map, diag.Diagnostics) {
	var diags diag.Diagnostics
	if len(tags) == 0 {
		return types.MapNull(types.StringType), diags
	}
	values := make(map[string]attr.Value, len(tags))
	for _, tag := range tags {
		key := ""
		value := ""
		if tag.TagKey != nil {
			key = *tag.TagKey
		}
		if tag.TagValue != nil {
			value = *tag.TagValue
		}
		values[key] = types.StringValue(value)
	}
	result, d := types.MapValue(types.StringType, values)
	diags.Append(d...)
	return result, diags
}

func expandDefaultActions(ctx context.Context, list types.List) ([]*albv20251030.DefaultAction, diag.Diagnostics) {
	var diags diag.Diagnostics
	if list.IsNull() || list.IsUnknown() {
		return nil, diags
	}

	var actions []albDefaultActionModel
	diags.Append(list.ElementsAs(ctx, &actions, false)...)
	if diags.HasError() {
		return nil, diags
	}

	out := make([]*albv20251030.DefaultAction, 0, len(actions))
	for _, action := range actions {
		item := &albv20251030.DefaultAction{
			Type: fwhelper.StringPointerFromValue(action.Type),
		}

		if !action.TargetGroupConfig.IsNull() && !action.TargetGroupConfig.IsUnknown() {
			var config albTargetGroupConfigModel
			diags.Append(action.TargetGroupConfig.As(ctx, &config, basetypes.ObjectAsOptions{})...)
			if diags.HasError() {
				return nil, diags
			}

			tgConfig := &albv20251030.TargetGroupConfig{}

			if !config.TargetGroups.IsNull() && !config.TargetGroups.IsUnknown() {
				var groups []albTargetGroupModel
				diags.Append(config.TargetGroups.ElementsAs(ctx, &groups, false)...)
				if diags.HasError() {
					return nil, diags
				}
				targetGroups := make([]*albv20251030.TargetGroupTuple, 0, len(groups))
				for _, group := range groups {
					targetGroups = append(targetGroups, &albv20251030.TargetGroupTuple{
						TargetGroupId: fwhelper.StringPointerFromValue(group.TargetGroupId),
						Weight:        int64PointerFromModel(group.Weight),
					})
				}
				tgConfig.TargetGroups = targetGroups
			}

			if !config.TargetGroupStickySession.IsNull() && !config.TargetGroupStickySession.IsUnknown() {
				var sticky albTargetGroupStickySessionModel
				diags.Append(config.TargetGroupStickySession.As(ctx, &sticky, basetypes.ObjectAsOptions{})...)
				if diags.HasError() {
					return nil, diags
				}
				tgConfig.TargetGroupStickySession = &albv20251030.TargetGroupStickySession{
					Enabled: fwhelper.BoolPointerFromValue(sticky.Enabled),
					Timeout: int64PointerFromModel(sticky.Timeout),
				}
			}

			item.TargetGroupConfig = tgConfig
		}

		out = append(out, item)
	}
	return out, diags
}

func int64PointerFromModel(v types.Int64) *int64 {
	if v.IsNull() || v.IsUnknown() {
		return nil
	}
	i := v.ValueInt64()
	return &i
}

func expandXForwardedForConfig(ctx context.Context, obj types.Object) (*albv20251030.XForwardedForConfig, diag.Diagnostics) {
	var diags diag.Diagnostics
	if obj.IsNull() || obj.IsUnknown() {
		return nil, diags
	}

	var config albXForwardedForConfigModel
	diags.Append(obj.As(ctx, &config, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return nil, diags
	}

	return &albv20251030.XForwardedForConfig{
		XForwardedForMode:                 fwhelper.StringPointerFromValue(config.XForwardedForMode),
		XForwardedForHostEnabled:          fwhelper.BoolPointerFromValue(config.XForwardedForHostEnabled),
		XForwardedForPortEnabled:          fwhelper.BoolPointerFromValue(config.XForwardedForPortEnabled),
		XForwardedForProtoEnabled:         fwhelper.BoolPointerFromValue(config.XForwardedForProtoEnabled),
		XForwardedForAlbIdEnabled:         fwhelper.BoolPointerFromValue(config.XForwardedForAlbIdEnabled),
		XForwardedForClientSrcPortEnabled: fwhelper.BoolPointerFromValue(config.XForwardedForClientSrcPortEnabled),
		XTencentClientVerifyEnabled:       fwhelper.BoolPointerFromValue(config.XTencentClientVerifyEnabled),
		XTencentClientSDNEnabled:          fwhelper.BoolPointerFromValue(config.XTencentClientSDNEnabled),
		XTencentClientIDNEnabled:          fwhelper.BoolPointerFromValue(config.XTencentClientIDNEnabled),
		XTencentClientSerialEnabled:       fwhelper.BoolPointerFromValue(config.XTencentClientSerialEnabled),
	}, diags
}

func refreshAlbListenerModel(ctx context.Context, m *AlbListenerResourceModel, loadBalancerId, listenerId string, item *albv20251030.ListenerOutput) diag.Diagnostics {
	var diags diag.Diagnostics

	var prevXff albXForwardedForConfigModel
	if !m.XForwardedForConfig.IsNull() && !m.XForwardedForConfig.IsUnknown() {
		diags.Append(m.XForwardedForConfig.As(ctx, &prevXff, basetypes.ObjectAsOptions{})...)
		if diags.HasError() {
			return diags
		}
	}

	m.LoadBalancerId = types.StringValue(loadBalancerId)
	m.Id = types.StringValue(strings.Join([]string{loadBalancerId, listenerId}, tccommon.FILED_SP))

	m.ListenerPort = int64ValueOrPrevUint(item.ListenerPort, m.ListenerPort)
	m.ListenerProtocol = fwhelper.StringValueOrPrev(item.ListenerProtocol, m.ListenerProtocol)
	m.ListenerName = fwhelper.StringValueOrPrev(item.ListenerName, m.ListenerName)
	m.GzipEnabled = boolValueOrPrev(item.GzipEnabled, m.GzipEnabled)
	m.Http2Enabled = boolValueOrPrev(item.Http2Enable, m.Http2Enabled)
	m.IdleTimeout = int64ValueOrPrevUint(item.IdleTimeout, m.IdleTimeout)
	m.RequestTimeout = int64ValueOrPrevUint(item.RequestTimeout, m.RequestTimeout)
	m.SecurityPolicyId = fwhelper.StringValueOrPrev(item.TlsSecurityPolicyId, m.SecurityPolicyId)
	m.CaEnabled = boolValueOrPrev(item.CaEnable, m.CaEnabled)
	m.ListenerStatus = types.StringPointerValue(item.ListenerStatus)
	m.CreateTime = types.StringPointerValue(item.CreateTime)
	m.ModifyTime = types.StringPointerValue(item.ModifyTime)

	tags, tagDiags := flattenTags(ctx, item.Tags)
	diags.Append(tagDiags...)
	m.Tags = tags

	// The Describe API does not echo certificate_ids, ca_certificate_ids or
	// default_actions back, so those values are preserved from plan/state.

	if item.XForwardedForConfig != nil {
		xff := albXForwardedForConfigModel{
			XForwardedForMode:                 fwhelper.StringValueOrPrev(item.XForwardedForConfig.XForwardedForMode, prevXff.XForwardedForMode),
			XForwardedForHostEnabled:          boolValueOrPrev(item.XForwardedForConfig.XForwardedForHostEnabled, prevXff.XForwardedForHostEnabled),
			XForwardedForPortEnabled:          boolValueOrPrev(item.XForwardedForConfig.XForwardedForPortEnabled, prevXff.XForwardedForPortEnabled),
			XForwardedForProtoEnabled:         boolValueOrPrev(item.XForwardedForConfig.XForwardedForProtoEnabled, prevXff.XForwardedForProtoEnabled),
			XForwardedForAlbIdEnabled:         boolValueOrPrev(item.XForwardedForConfig.XForwardedForAlbIdEnabled, prevXff.XForwardedForAlbIdEnabled),
			XForwardedForClientSrcPortEnabled: boolValueOrPrev(item.XForwardedForConfig.XForwardedForClientSrcPortEnabled, prevXff.XForwardedForClientSrcPortEnabled),
			XTencentClientVerifyEnabled:       boolValueOrPrev(item.XForwardedForConfig.XTencentClientVerifyEnabled, prevXff.XTencentClientVerifyEnabled),
			XTencentClientSDNEnabled:          boolValueOrPrev(item.XForwardedForConfig.XTencentClientSDNEnabled, prevXff.XTencentClientSDNEnabled),
			XTencentClientIDNEnabled:          boolValueOrPrev(item.XForwardedForConfig.XTencentClientIDNEnabled, prevXff.XTencentClientIDNEnabled),
			XTencentClientSerialEnabled:       boolValueOrPrev(item.XForwardedForConfig.XTencentClientSerialEnabled, prevXff.XTencentClientSerialEnabled),
		}
		obj, d := types.ObjectValueFrom(ctx, xForwardedForObjectTypes(ctx), xff)
		diags.Append(d...)
		m.XForwardedForConfig = obj
	}

	return diags
}

func int64ValueOrPrevUint(cloud *uint64, prev types.Int64) types.Int64 {
	if cloud != nil {
		return types.Int64Value(int64(*cloud))
	}
	if !prev.IsNull() && !prev.IsUnknown() {
		return prev
	}
	return types.Int64Null()
}

func boolValueOrPrev(cloud *bool, prev types.Bool) types.Bool {
	if cloud != nil {
		return types.BoolValue(*cloud)
	}
	if !prev.IsNull() && !prev.IsUnknown() {
		return prev
	}
	return types.BoolNull()
}




func xForwardedForObjectTypes(ctx context.Context) map[string]attr.Type {
	attrTypes := make(map[string]attr.Type, len(xForwardedForAttributeTypes(ctx)))
	for name, attribute := range xForwardedForAttributeTypes(ctx) {
		attrTypes[name] = attribute.GetType()
	}
	return attrTypes
}
