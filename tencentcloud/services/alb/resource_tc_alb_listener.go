package alb

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
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	sdkv2resource "github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	albv20251030 "github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/alb/v20251030"

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
}

type AlbListenerResourceModel struct {
	Id                  types.String `tfsdk:"id"`
	LoadBalancerId      types.String `tfsdk:"load_balancer_id"`
	ListenerId          types.String `tfsdk:"listener_id"`
	ListenerPort        types.Int64  `tfsdk:"listener_port"`
	ListenerProtocol    types.String `tfsdk:"listener_protocol"`
	ListenerName        types.String `tfsdk:"listener_name"`
	DefaultActions      types.List   `tfsdk:"default_actions"`
	CertificateIds      types.List   `tfsdk:"certificate_ids"`
	CaCertificateIds    types.List   `tfsdk:"ca_certificate_ids"`
	CaEnabled           types.Bool   `tfsdk:"ca_enabled"`
	SecurityPolicyId    types.String `tfsdk:"security_policy_id"`
	GzipEnabled         types.Bool   `tfsdk:"gzip_enabled"`
	Http2Enabled        types.Bool   `tfsdk:"http2enabled"`
	IdleTimeout         types.Int64  `tfsdk:"idle_timeout"`
	RequestTimeout      types.Int64  `tfsdk:"request_timeout"`
	Tags                types.List   `tfsdk:"tags"`
	XForwardedForConfig types.Object `tfsdk:"x_forwarded_for_config"`
	ClientToken         types.String `tfsdk:"client_token"`
}

type albListenerDefaultActionModel struct {
	Type              types.String `tfsdk:"type"`
	TargetGroupConfig types.Object `tfsdk:"target_group_config"`
}

type albListenerTargetGroupConfigModel struct {
	TargetGroups             types.List   `tfsdk:"target_groups"`
	TargetGroupStickySession types.Object `tfsdk:"target_group_sticky_session"`
}

type albListenerTargetGroupModel struct {
	TargetGroupId types.String `tfsdk:"target_group_id"`
	Weight        types.Int64  `tfsdk:"weight"`
}

type albListenerStickySessionModel struct {
	Enabled types.Bool  `tfsdk:"enabled"`
	Timeout types.Int64 `tfsdk:"timeout"`
}

type albListenerTagModel struct {
	TagKey   types.String `tfsdk:"tag_key"`
	TagValue types.String `tfsdk:"tag_value"`
}

type albListenerXForwardedForConfigModel struct {
	XForwardedForMode                 types.String `tfsdk:"x_forwarded_for_mode"`
	XForwardedForAlbIdEnabled         types.Bool   `tfsdk:"x_forwarded_for_alb_id_enabled"`
	XForwardedForClientSrcPortEnabled types.Bool   `tfsdk:"x_forwarded_for_client_src_port_enabled"`
	XForwardedForHostEnabled          types.Bool   `tfsdk:"x_forwarded_for_host_enabled"`
	XForwardedForPortEnabled          types.Bool   `tfsdk:"x_forwarded_for_port_enabled"`
	XForwardedForProtoEnabled         types.Bool   `tfsdk:"x_forwarded_for_proto_enabled"`
	XTencentClientIDNEnabled          types.Bool   `tfsdk:"x_tencent_client_idn_enabled"`
	XTencentClientSDNEnabled          types.Bool   `tfsdk:"x_tencent_client_sdn_enabled"`
	XTencentClientSerialEnabled       types.Bool   `tfsdk:"x_tencent_client_serial_enabled"`
	XTencentClientVerifyEnabled       types.Bool   `tfsdk:"x_tencent_client_verify_enabled"`
}

func (r *AlbListenerResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "tencentcloud_alb_listener"
}

func (r *AlbListenerResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a resource to create a listener of an application load balancer (ALB).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Resource identifier in the form `<load_balancer_id>#<listener_id>`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"load_balancer_id": schema.StringAttribute{
				Required:    true,
				Description: "Load balancer instance ID, in the format of alb- followed by 8 alphanumeric characters. Immutable: changing it forces resource replacement.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"listener_id": schema.StringAttribute{
				Computed:    true,
				Description: "Listener ID, in the format of lst- followed by 8 alphanumeric characters.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"listener_port": schema.Int64Attribute{
				Required:    true,
				Description: "Port used by the load balancer frontend. Value range: 1-65535. Immutable: changing it forces resource replacement.",
				PlanModifiers: []planmodifier.Int64{
					int64RequiresReplace{},
				},
			},
			"listener_protocol": schema.StringAttribute{
				Required:    true,
				Description: "Listening protocol. Valid values: HTTP, HTTPS or QUIC. Immutable: changing it forces resource replacement.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"listener_name": schema.StringAttribute{
				Optional:    true,
				Description: "Custom listener name, 1-255 characters. It can contain Chinese characters, letters, digits, dashes (-), forward slashes (/), periods (.) and underscores (_).",
			},
			"default_actions": schema.ListNestedAttribute{
				Required:    true,
				Description: "Default forwarding rule action list. Currently, a listener supports adding only 1 default forwarding rule action.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"type": schema.StringAttribute{
							Required:    true,
							Description: "Forward action type. When a listener is created, the default forward action type only supports forwarding to a target group.",
						},
						"target_group_config": schema.SingleNestedAttribute{
							Required:    true,
							Description: "Forward target group configuration. Only a single target group is supported when creating a listener.",
							Attributes: map[string]schema.Attribute{
								"target_groups": schema.ListNestedAttribute{
									Required:    true,
									Description: "Target group list.",
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{
											"target_group_id": schema.StringAttribute{
												Required:    true,
												Description: "Target group ID, in the format of lbtg- followed by 8 alphanumeric characters.",
											},
											"weight": schema.Int64Attribute{
												Optional:    true,
												Description: "Weight. Value range: [0, 100]. Default value: 10.",
											},
										},
									},
								},
								"target_group_sticky_session": schema.SingleNestedAttribute{
									Optional:    true,
									Description: "Session persistence between target groups.",
									Attributes: map[string]schema.Attribute{
										"enabled": schema.BoolAttribute{
											Required:    true,
											Description: "Whether to enable session persistence. Disabled by default.",
										},
										"timeout": schema.Int64Attribute{
											Optional:    true,
											Description: "Timeout period in seconds. Value range: 1-86400. Default value: 1000.",
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
				ElementType: types.StringType,
				Description: "List of server certificate IDs.",
			},
			"ca_certificate_ids": schema.ListAttribute{
				Optional:    true,
				ElementType: types.StringType,
				Description: "List of CA certificate IDs configured for the listener. Currently, a listener supports adding only 1 CA certificate. This parameter is required when ca_enabled is true.",
			},
			"ca_enabled": schema.BoolAttribute{
				Optional:    true,
				Description: "Whether mutual authentication is enabled. Default value: false.",
			},
			"security_policy_id": schema.StringAttribute{
				Optional:    true,
				Description: "Security policy ID, in the format of tls- followed by 8 alphanumeric characters.",
			},
			"gzip_enabled": schema.BoolAttribute{
				Optional:    true,
				Description: "Whether Gzip compression is enabled. Default value: true.",
			},
			"http2enabled": schema.BoolAttribute{
				Optional:    true,
				Description: "Whether HTTP/2 is enabled. Default value: false for HTTP and true for HTTPS. Only the HTTPS protocol supports this parameter.",
			},
			"idle_timeout": schema.Int64Attribute{
				Optional:    true,
				Description: "Connection idle timeout, in seconds. Value range: 1-600. Default value: 15.",
			},
			"request_timeout": schema.Int64Attribute{
				Optional:    true,
				Description: "Connection request timeout period, in seconds. Value range: 1-600. Default value: 60.",
			},
			"tags": schema.ListNestedAttribute{
				Optional:    true,
				Description: "Tag list. Supports up to 20 tags.",
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
			"x_forwarded_for_config": schema.SingleNestedAttribute{
				Optional:    true,
				Description: "X-Forwarded-For configuration.",
				Attributes: map[string]schema.Attribute{
					"x_forwarded_for_mode": schema.StringAttribute{
						Optional:    true,
						Description: "Specify how to handle the X-Forwarded-For (XFF) HTTP header field. Valid values: append (default), remove, passthrough.",
					},
					"x_forwarded_for_alb_id_enabled": schema.BoolAttribute{
						Optional:    true,
						Description: "Whether to get the load balancer instance ID through the ALB-ID header field.",
					},
					"x_forwarded_for_client_src_port_enabled": schema.BoolAttribute{
						Optional:    true,
						Description: "Whether to get the client port through the X-Forwarded-Client-srcport header field.",
					},
					"x_forwarded_for_host_enabled": schema.BoolAttribute{
						Optional:    true,
						Description: "Whether to get the client domain name through the X-Forwarded-Host header field.",
					},
					"x_forwarded_for_port_enabled": schema.BoolAttribute{
						Optional:    true,
						Description: "Whether to get the listening port of the load balancer instance through the X-Forwarded-Port header field.",
					},
					"x_forwarded_for_proto_enabled": schema.BoolAttribute{
						Optional:    true,
						Description: "Whether to get the listening protocol of the load balancer instance through the X-Forwarded-Proto header field.",
					},
					"x_tencent_client_idn_enabled": schema.BoolAttribute{
						Optional:    true,
						Description: "Whether to access the issuer of the client certificate ($ssl_client_i_dn) through the X-Tencent-Client-IDN header.",
					},
					"x_tencent_client_sdn_enabled": schema.BoolAttribute{
						Optional:    true,
						Description: "Whether to access the subject of the client certificate ($ssl_client_s_dn) through the X-Tencent-Client-SDN header.",
					},
					"x_tencent_client_serial_enabled": schema.BoolAttribute{
						Optional:    true,
						Description: "Whether to access the serial number of the client certificate ($ssl_client_serial) through the X-Tencent-Client-Serial header.",
					},
					"x_tencent_client_verify_enabled": schema.BoolAttribute{
						Optional:    true,
						Description: "Whether to access the verification result of the client certificate ($ssl_client_verify) through the X-Tencent-Client-Verify header.",
					},
				},
			},
			"client_token": schema.StringAttribute{
				Optional:    true,
				Description: "Client token, used to ensure the idempotency of requests.",
			},
		},
	}
}

func (r *AlbListenerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan AlbListenerResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if r.Client() == nil {
		resp.Diagnostics.AddError("Provider not configured", "The provider client is not available. This is unexpected; please report it.")
		return
	}

	request := albv20251030.NewCreateListenerRequest()
	request.LoadBalancerId = helper.String(plan.LoadBalancerId.ValueString())
	request.ListenerPort = fwhelper.Uint64PointerFromValue(plan.ListenerPort)
	request.ListenerProtocol = helper.String(plan.ListenerProtocol.ValueString())
	request.ListenerName = fwhelper.StringPointerFromValue(plan.ListenerName)
	request.CertificateIds = expandStringList(plan.CertificateIds)
	request.CaCertificateIds = expandStringList(plan.CaCertificateIds)
	request.CaEnabled = fwhelper.BoolPointerFromValue(plan.CaEnabled)
	request.SecurityPolicyId = fwhelper.StringPointerFromValue(plan.SecurityPolicyId)
	request.GzipEnabled = fwhelper.BoolPointerFromValue(plan.GzipEnabled)
	request.Http2Enabled = fwhelper.BoolPointerFromValue(plan.Http2Enabled)
	request.IdleTimeout = fwhelper.Uint64PointerFromValue(plan.IdleTimeout)
	request.RequestTimeout = fwhelper.Uint64PointerFromValue(plan.RequestTimeout)
	request.ClientToken = fwhelper.StringPointerFromValue(plan.ClientToken)

	defaultActions, diags := expandDefaultActions(ctx, plan.DefaultActions)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	request.DefaultActions = defaultActions

	tags, tagDiags := expandTags(ctx, plan.Tags)
	resp.Diagnostics.Append(tagDiags...)
	if resp.Diagnostics.HasError() {
		return
	}
	request.Tags = tags

	xff, xffDiags := expandXForwardedForConfig(plan.XForwardedForConfig)
	resp.Diagnostics.Append(xffDiags...)
	if resp.Diagnostics.HasError() {
		return
	}
	request.XForwardedForConfig = xff

	var response *albv20251030.CreateListenerResponse
	var err error
	err = sdkv2resource.RetryContext(ctx, tccommon.WriteRetryTimeout, func() *sdkv2resource.RetryError {
		response, err = r.Client().UseAlbClient().CreateListenerWithContext(ctx, request)
		if err != nil {
			return tccommon.RetryError(err)
		}
		return nil
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating ALB listener",
			fmt.Sprintf("CreateListener failed (load_balancer_id=%q): %s", plan.LoadBalancerId.ValueString(), err),
		)
		return
	}
	if response == nil || response.Response == nil || response.Response.ListenerId == nil || *response.Response.ListenerId == "" {
		resp.Diagnostics.AddError(
			"Error creating ALB listener",
			fmt.Sprintf("CreateListener returned no listener_id (load_balancer_id=%q).", plan.LoadBalancerId.ValueString()),
		)
		return
	}

	listenerId := *response.Response.ListenerId
	plan.ListenerId = types.StringValue(listenerId)
	plan.Id = types.StringValue(strings.Join([]string{plan.LoadBalancerId.ValueString(), listenerId}, tccommon.FILED_SP))

	item, err := describeAlbListener(ctx, r.Client(), plan.LoadBalancerId.ValueString(), listenerId)
	if err != nil {
		if errors.Is(err, errAlbListenerNotFound) {
			resp.Diagnostics.AddError(
				"Error reading ALB listener after create",
				fmt.Sprintf("CreateListener succeeded but the listener %s cannot be read back yet: %s", listenerId, err),
			)
			return
		}
		resp.Diagnostics.AddError(
			"Error reading ALB listener after create",
			fmt.Sprintf("CreateListener succeeded but DescribeListenerDetail failed (%s): %s", listenerId, err),
		)
		return
	}

	resp.Diagnostics.Append(refreshAlbListenerModel(ctx, &plan, item)...)
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

	if r.Client() == nil {
		resp.Diagnostics.AddError("Provider not configured", "The provider client is not available. This is unexpected; please report it.")
		return
	}

	loadBalancerId, listenerId, ok := parseAlbListenerId(&resp.Diagnostics, state.Id.ValueString())
	if !ok {
		return
	}
	state.LoadBalancerId = types.StringValue(loadBalancerId)

	item, err := describeAlbListener(ctx, r.Client(), loadBalancerId, listenerId)
	if err != nil {
		if errors.Is(err, errAlbListenerNotFound) {
			resp.Diagnostics.AddWarning(
				"Resource not found",
				fmt.Sprintf("ALB listener %s no longer exists, removing it from state.", listenerId),
			)
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error reading ALB listener",
			fmt.Sprintf("DescribeListenerDetail failed (%s): %s", listenerId, err),
		)
		return
	}

	resp.Diagnostics.Append(refreshAlbListenerModel(ctx, &state, item)...)
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

	if r.Client() == nil {
		resp.Diagnostics.AddError("Provider not configured", "The provider client is not available. This is unexpected; please report it.")
		return
	}

	_, listenerId, ok := parseAlbListenerId(&resp.Diagnostics, plan.Id.ValueString())
	if !ok {
		return
	}

	request := albv20251030.NewModifyListenerAttributesRequest()
	request.ListenerId = helper.String(listenerId)
	request.LoadBalancerId = helper.String(plan.LoadBalancerId.ValueString())
	request.ListenerName = fwhelper.StringPointerFromValue(plan.ListenerName)
	request.CertificateIds = expandStringList(plan.CertificateIds)
	request.CaCertificateIds = expandStringList(plan.CaCertificateIds)
	request.CaEnabled = fwhelper.BoolPointerFromValue(plan.CaEnabled)
	request.SecurityPolicyId = fwhelper.StringPointerFromValue(plan.SecurityPolicyId)
	request.GzipEnabled = fwhelper.BoolPointerFromValue(plan.GzipEnabled)
	request.Http2Enabled = fwhelper.BoolPointerFromValue(plan.Http2Enabled)
	request.IdleTimeout = fwhelper.Uint64PointerFromValue(plan.IdleTimeout)
	request.RequestTimeout = fwhelper.Uint64PointerFromValue(plan.RequestTimeout)
	request.ClientToken = fwhelper.StringPointerFromValue(plan.ClientToken)

	defaultActions, diags := expandDefaultActions(ctx, plan.DefaultActions)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	request.DefaultActions = defaultActions

	xff, xffDiags := expandXForwardedForConfig(plan.XForwardedForConfig)
	resp.Diagnostics.Append(xffDiags...)
	if resp.Diagnostics.HasError() {
		return
	}
	request.XForwardedForConfig = xff

	err := sdkv2resource.RetryContext(ctx, tccommon.WriteRetryTimeout, func() *sdkv2resource.RetryError {
		_, err := r.Client().UseAlbClient().ModifyListenerAttributesWithContext(ctx, request)
		if err != nil {
			return tccommon.RetryError(err)
		}
		return nil
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating ALB listener",
			fmt.Sprintf("ModifyListenerAttributes failed (%s): %s", listenerId, err),
		)
		return
	}

	item, err := describeAlbListener(ctx, r.Client(), plan.LoadBalancerId.ValueString(), listenerId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading ALB listener after update",
			fmt.Sprintf("ModifyListenerAttributes succeeded but DescribeListenerDetail failed (%s): %s", listenerId, err),
		)
		return
	}

	resp.Diagnostics.Append(refreshAlbListenerModel(ctx, &plan, item)...)
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

	if r.Client() == nil {
		resp.Diagnostics.AddError("Provider not configured", "The provider client is not available. This is unexpected; please report it.")
		return
	}

	loadBalancerId, listenerId, ok := parseAlbListenerId(&resp.Diagnostics, state.Id.ValueString())
	if !ok {
		return
	}

	request := albv20251030.NewDeleteListenerRequest()
	request.ListenerIds = helper.Strings([]string{listenerId})
	request.LoadBalancerId = helper.String(loadBalancerId)

	err := sdkv2resource.RetryContext(ctx, tccommon.WriteRetryTimeout, func() *sdkv2resource.RetryError {
		_, err := r.Client().UseAlbClient().DeleteListenerWithContext(ctx, request)
		if err != nil {
			if tccommon.IsExpectError(err, []string{"ResourceNotFound", "ResourceInexistent"}) {
				return nil
			}
			return tccommon.RetryError(err)
		}
		return nil
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting ALB listener",
			fmt.Sprintf("DeleteListener failed (%s): %s", listenerId, err),
		)
	}
}

func (r *AlbListenerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if _, _, ok := splitAlbListenerId(req.ID); !ok {
		resp.Diagnostics.AddError(
			"Invalid import ID",
			fmt.Sprintf("Import id %q is invalid, expected `<load_balancer_id>%s<listener_id>`.", req.ID, tccommon.FILED_SP),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
}

func describeAlbListener(ctx context.Context, client albClientInterface, loadBalancerId, listenerId string) (*albv20251030.DescribeListenerDetailResponseParams, error) {
	request := albv20251030.NewDescribeListenerDetailRequest()
	request.ListenerId = helper.String(listenerId)
	request.LoadBalancerId = helper.String(loadBalancerId)

	var response *albv20251030.DescribeListenerDetailResponse
	var err error
	err = sdkv2resource.RetryContext(ctx, tccommon.ReadRetryTimeout, func() *sdkv2resource.RetryError {
		response, err = client.UseAlbClient().DescribeListenerDetailWithContext(ctx, request)
		if err != nil {
			if tccommon.IsExpectError(err, []string{"ResourceNotFound"}) {
				return sdkv2resource.NonRetryableError(errAlbListenerNotFound)
			}
			return tccommon.RetryError(err)
		}
		return nil
	})
	if err != nil {
		if errors.Is(err, errAlbListenerNotFound) {
			return nil, errAlbListenerNotFound
		}
		return nil, err
	}
	if response == nil || response.Response == nil || response.Response.ListenerId == nil || *response.Response.ListenerId == "" {
		return nil, errAlbListenerNotFound
	}
	return response.Response, nil
}

func refreshAlbListenerModel(ctx context.Context, m *AlbListenerResourceModel, item *albv20251030.DescribeListenerDetailResponseParams) diag.Diagnostics {
	var diags diag.Diagnostics

	if item == nil {
		diags.AddError("Malformed ALB listener response", "DescribeListenerDetail returned an empty response.")
		return diags
	}

	m.Id = types.StringValue(strings.Join([]string{*item.LoadBalancerId, *item.ListenerId}, tccommon.FILED_SP))
	m.LoadBalancerId = fwhelper.StringValueOrPrev(item.LoadBalancerId, m.LoadBalancerId)
	m.ListenerId = fwhelper.StringValueOrPrev(item.ListenerId, m.ListenerId)
	m.ListenerPort = fwhelper.Int64ValueOrPrev(item.ListenerPort, m.ListenerPort)
	m.ListenerProtocol = fwhelper.StringValueOrPrev(item.ListenerProtocol, m.ListenerProtocol)
	m.ListenerName = fwhelper.StringValueOrPrev(item.ListenerName, m.ListenerName)
	m.CaEnabled = fwhelper.BoolValueOrNull(item.CaEnabled)
	m.GzipEnabled = fwhelper.BoolValueOrNull(item.GzipEnabled)
	m.Http2Enabled = fwhelper.BoolValueOrNull(item.Http2Enabled)
	m.IdleTimeout = fwhelper.Int64ValueOrPrev(item.IdleTimeout, m.IdleTimeout)
	m.RequestTimeout = fwhelper.Int64ValueOrPrev(item.RequestTimeout, m.RequestTimeout)
	m.SecurityPolicyId = fwhelper.StringValueOrPrev(item.SecurityPolicyId, m.SecurityPolicyId)

	if len(item.CertificateIds) > 0 {
		m.CertificateIds = types.ListValueMust(types.StringType, stringElements(item.CertificateIds))
	} else if m.CertificateIds.IsNull() {
		m.CertificateIds = types.ListNull(types.StringType)
	}

	if len(item.CaCertificateIds) > 0 {
		m.CaCertificateIds = types.ListValueMust(types.StringType, stringElements(item.CaCertificateIds))
	} else if m.CaCertificateIds.IsNull() {
		m.CaCertificateIds = types.ListNull(types.StringType)
	}

	if m.ClientToken.IsNull() {
		m.ClientToken = types.StringNull()
	}

	actions, actionDiags := flattenDefaultActions(ctx, item.DefaultActions, m.DefaultActions)
	diags.Append(actionDiags...)
	m.DefaultActions = actions

	tags, tagDiags := flattenTags(ctx, item.Tags, m.Tags)
	diags.Append(tagDiags...)
	m.Tags = tags

	xff, xffDiags := flattenXForwardedForConfig(ctx, item.XForwardedForConfig, m.XForwardedForConfig)
	diags.Append(xffDiags...)
	m.XForwardedForConfig = xff

	if m.ListenerPort.IsNull() || m.ListenerProtocol.IsNull() || m.DefaultActions.IsNull() {
		diags.AddError(
			"Malformed ALB listener response",
			fmt.Sprintf("DescribeListenerDetail returned no listener_port, listener_protocol or default_actions (%s).", *item.ListenerId),
		)
	}
	return diags
}

func int64FromPointer(v *int64) types.Int64 {
	if v == nil {
		return types.Int64Null()
	}
	return types.Int64Value(*v)
}

func stickyTimeout(timeout *int64) types.Int64 {
	if timeout == nil {
		return types.Int64Null()
	}
	return types.Int64Value(*timeout)
}

func expandDefaultActions(ctx context.Context, list types.List) ([]*albv20251030.DefaultAction, diag.Diagnostics) {
	var diags diag.Diagnostics
	if list.IsNull() || list.IsUnknown() {
		return nil, diags
	}

	var actions []albListenerDefaultActionModel
	diags.Append(list.ElementsAs(ctx, &actions, false)...)
	if diags.HasError() {
		return nil, diags
	}

	out := make([]*albv20251030.DefaultAction, 0, len(actions))
	for _, action := range actions {
		config, configDiags := expandTargetGroupConfig(ctx, action.TargetGroupConfig)
		diags.Append(configDiags...)
		if diags.HasError() {
			return nil, diags
		}
		out = append(out, &albv20251030.DefaultAction{
			Type:              fwhelper.StringPointerFromValue(action.Type),
			TargetGroupConfig: config,
		})
	}
	return out, diags
}

func expandTargetGroupConfig(ctx context.Context, obj types.Object) (*albv20251030.TargetGroupConfig, diag.Diagnostics) {
	var diags diag.Diagnostics
	if obj.IsNull() || obj.IsUnknown() {
		return nil, diags
	}

	var config albListenerTargetGroupConfigModel
	diags.Append(obj.As(ctx, &config, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return nil, diags
	}

	var groups []albListenerTargetGroupModel
	if !config.TargetGroups.IsNull() && !config.TargetGroups.IsUnknown() {
		diags.Append(config.TargetGroups.ElementsAs(ctx, &groups, false)...)
		if diags.HasError() {
			return nil, diags
		}
	}

	targetGroups := make([]*albv20251030.TargetGroupTuple, 0, len(groups))
	for _, group := range groups {
		targetGroups = append(targetGroups, &albv20251030.TargetGroupTuple{
			TargetGroupId: fwhelper.StringPointerFromValue(group.TargetGroupId),
			Weight:        fwhelper.Int64PointerFromValue(group.Weight),
		})
	}

	var sticky *albv20251030.TargetGroupStickySession
	if !config.TargetGroupStickySession.IsNull() && !config.TargetGroupStickySession.IsUnknown() {
		var stickyModel albListenerStickySessionModel
		diags.Append(config.TargetGroupStickySession.As(ctx, &stickyModel, basetypes.ObjectAsOptions{})...)
		if diags.HasError() {
			return nil, diags
		}
		sticky = &albv20251030.TargetGroupStickySession{
			Enabled: fwhelper.BoolPointerFromValue(stickyModel.Enabled),
			Timeout: fwhelper.Int64PointerFromValue(stickyModel.Timeout),
		}
	}

	return &albv20251030.TargetGroupConfig{
		TargetGroups:             targetGroups,
		TargetGroupStickySession: sticky,
	}, diags
}

func expandTags(ctx context.Context, list types.List) ([]*albv20251030.TagInfo, diag.Diagnostics) {
	var diags diag.Diagnostics
	if list.IsNull() || list.IsUnknown() {
		return nil, diags
	}

	var tags []albListenerTagModel
	diags.Append(list.ElementsAs(ctx, &tags, false)...)
	if diags.HasError() {
		return nil, diags
	}

	out := make([]*albv20251030.TagInfo, 0, len(tags))
	for _, tag := range tags {
		out = append(out, &albv20251030.TagInfo{
			TagKey:   fwhelper.StringPointerFromValue(tag.TagKey),
			TagValue: fwhelper.StringPointerFromValue(tag.TagValue),
		})
	}
	return out, diags
}

func expandXForwardedForConfig(obj types.Object) (*albv20251030.XForwardedForConfig, diag.Diagnostics) {
	var diags diag.Diagnostics
	if obj.IsNull() || obj.IsUnknown() {
		return nil, diags
	}

	var xff albListenerXForwardedForConfigModel
	diags.Append(obj.As(context.Background(), &xff, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return nil, diags
	}

	return &albv20251030.XForwardedForConfig{
		XForwardedForMode:                 fwhelper.StringPointerFromValue(xff.XForwardedForMode),
		XForwardedForAlbIdEnabled:         fwhelper.BoolPointerFromValue(xff.XForwardedForAlbIdEnabled),
		XForwardedForClientSrcPortEnabled: fwhelper.BoolPointerFromValue(xff.XForwardedForClientSrcPortEnabled),
		XForwardedForHostEnabled:          fwhelper.BoolPointerFromValue(xff.XForwardedForHostEnabled),
		XForwardedForPortEnabled:          fwhelper.BoolPointerFromValue(xff.XForwardedForPortEnabled),
		XForwardedForProtoEnabled:         fwhelper.BoolPointerFromValue(xff.XForwardedForProtoEnabled),
		XTencentClientIDNEnabled:          fwhelper.BoolPointerFromValue(xff.XTencentClientIDNEnabled),
		XTencentClientSDNEnabled:          fwhelper.BoolPointerFromValue(xff.XTencentClientSDNEnabled),
		XTencentClientSerialEnabled:       fwhelper.BoolPointerFromValue(xff.XTencentClientSerialEnabled),
		XTencentClientVerifyEnabled:       fwhelper.BoolPointerFromValue(xff.XTencentClientVerifyEnabled),
	}, diags
}

func flattenDefaultActions(ctx context.Context, actions []*albv20251030.DefaultAction, prevList types.List) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics
	objType := defaultActionObjectType()

	out := make([]albListenerDefaultActionModel, 0, len(actions))
	for _, action := range actions {
		if action == nil {
			continue
		}

		var config types.Object
		if action.TargetGroupConfig != nil {
			groups := make([]albListenerTargetGroupModel, 0)
			for _, group := range action.TargetGroupConfig.TargetGroups {
				if group == nil {
					continue
				}
				groups = append(groups, albListenerTargetGroupModel{
					TargetGroupId: fwhelper.StringValueOrPrev(group.TargetGroupId, types.StringNull()),
					Weight:        int64FromPointer(group.Weight),
				})
			}

			groupList, groupDiags := types.ListValueFrom(ctx, targetGroupObjectType(), groups)
			diags.Append(groupDiags...)

			var sticky types.Object
			if action.TargetGroupConfig.TargetGroupStickySession != nil {
				stickyModel := albListenerStickySessionModel{
					Enabled: fwhelper.BoolValueOrNull(action.TargetGroupConfig.TargetGroupStickySession.Enabled),
					Timeout: stickyTimeout(action.TargetGroupConfig.TargetGroupStickySession.Timeout),
				}
				var stickyDiags diag.Diagnostics
				sticky, stickyDiags = types.ObjectValueFrom(ctx, stickySessionAttrTypes(), stickyModel)
				diags.Append(stickyDiags...)
			}

			configModel := albListenerTargetGroupConfigModel{
				TargetGroups:             groupList,
				TargetGroupStickySession: sticky,
			}
			var configDiags diag.Diagnostics
			config, configDiags = types.ObjectValueFrom(ctx, targetGroupConfigAttrTypes(), configModel)
			diags.Append(configDiags...)
		}

		out = append(out, albListenerDefaultActionModel{
			Type:              fwhelper.StringValueOrPrev(action.Type, types.StringNull()),
			TargetGroupConfig: config,
		})
	}

	if len(out) == 0 && !prevList.IsNull() && !prevList.IsUnknown() {
		return prevList, diags
	}

	list, listDiags := types.ListValueFrom(ctx, objType, out)
	diags.Append(listDiags...)
	if diags.HasError() {
		return types.ListNull(objType), diags
	}
	return list, diags
}

func flattenTags(ctx context.Context, tags []*albv20251030.TagInfo, prevList types.List) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics
	objType := tagObjectType()

	out := make([]albListenerTagModel, 0, len(tags))
	for _, tag := range tags {
		if tag == nil {
			continue
		}
		out = append(out, albListenerTagModel{
			TagKey:   fwhelper.StringValueOrPrev(tag.TagKey, types.StringNull()),
			TagValue: fwhelper.StringValueOrPrev(tag.TagValue, types.StringNull()),
		})
	}

	if len(out) == 0 && !prevList.IsNull() && !prevList.IsUnknown() {
		return prevList, diags
	}

	list, listDiags := types.ListValueFrom(ctx, objType, out)
	diags.Append(listDiags...)
	if diags.HasError() {
		return types.ListNull(objType), diags
	}
	return list, diags
}

func flattenXForwardedForConfig(ctx context.Context, xff *albv20251030.XForwardedForConfig, prev types.Object) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics
	if xff == nil {
		if prev.IsNull() || prev.IsUnknown() {
			return types.ObjectNull(xForwardedForConfigAttrTypes()), diags
		}
		return prev, diags
	}

	model := albListenerXForwardedForConfigModel{
		XForwardedForMode:                 fwhelper.StringValueOrPrev(xff.XForwardedForMode, types.StringNull()),
		XForwardedForAlbIdEnabled:         fwhelper.BoolValueOrNull(xff.XForwardedForAlbIdEnabled),
		XForwardedForClientSrcPortEnabled: fwhelper.BoolValueOrNull(xff.XForwardedForClientSrcPortEnabled),
		XForwardedForHostEnabled:          fwhelper.BoolValueOrNull(xff.XForwardedForHostEnabled),
		XForwardedForPortEnabled:          fwhelper.BoolValueOrNull(xff.XForwardedForPortEnabled),
		XForwardedForProtoEnabled:         fwhelper.BoolValueOrNull(xff.XForwardedForProtoEnabled),
		XTencentClientIDNEnabled:          fwhelper.BoolValueOrNull(xff.XTencentClientIDNEnabled),
		XTencentClientSDNEnabled:          fwhelper.BoolValueOrNull(xff.XTencentClientSDNEnabled),
		XTencentClientSerialEnabled:       fwhelper.BoolValueOrNull(xff.XTencentClientSerialEnabled),
		XTencentClientVerifyEnabled:       fwhelper.BoolValueOrNull(xff.XTencentClientVerifyEnabled),
	}

	obj, objDiags := types.ObjectValueFrom(ctx, xForwardedForConfigAttrTypes(), model)
	diags.Append(objDiags...)
	return obj, diags
}

func expandStringList(list types.List) []*string {
	if list.IsNull() || list.IsUnknown() {
		return nil
	}
	var values []string
	_ = list.ElementsAs(context.Background(), &values, false)
	if len(values) == 0 {
		return nil
	}
	return helper.Strings(values)
}

func stringElements(values []*string) []attr.Value {
	out := make([]attr.Value, 0, len(values))
	for _, value := range values {
		if value == nil {
			continue
		}
		out = append(out, types.StringValue(*value))
	}
	return out
}

func defaultActionObjectType() attr.Type {
	return types.ObjectType{AttrTypes: defaultActionAttrTypes()}
}

func defaultActionAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"type":                types.StringType,
		"target_group_config": types.ObjectType{AttrTypes: targetGroupConfigAttrTypes()},
	}
}

func targetGroupConfigObjectType() attr.Type {
	return types.ObjectType{AttrTypes: targetGroupConfigAttrTypes()}
}

func targetGroupConfigAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"target_groups":               types.ListType{ElemType: targetGroupObjectType()},
		"target_group_sticky_session": types.ObjectType{AttrTypes: stickySessionAttrTypes()},
	}
}

func targetGroupObjectType() attr.Type {
	return types.ObjectType{AttrTypes: targetGroupAttrTypes()}
}

func targetGroupAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"target_group_id": types.StringType,
		"weight":          types.Int64Type,
	}
}

func stickySessionAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"enabled": types.BoolType,
		"timeout": types.Int64Type,
	}
}

func tagObjectType() attr.Type {
	return types.ObjectType{AttrTypes: tagAttrTypes()}
}

func tagAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"tag_key":   types.StringType,
		"tag_value": types.StringType,
	}
}

func xForwardedForConfigAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"x_forwarded_for_mode":                    types.StringType,
		"x_forwarded_for_alb_id_enabled":          types.BoolType,
		"x_forwarded_for_client_src_port_enabled": types.BoolType,
		"x_forwarded_for_host_enabled":            types.BoolType,
		"x_forwarded_for_port_enabled":            types.BoolType,
		"x_forwarded_for_proto_enabled":           types.BoolType,
		"x_tencent_client_idn_enabled":            types.BoolType,
		"x_tencent_client_sdn_enabled":            types.BoolType,
		"x_tencent_client_serial_enabled":         types.BoolType,
		"x_tencent_client_verify_enabled":         types.BoolType,
	}
}

func albListenerId(loadBalancerId, listenerId string) string {
	return strings.Join([]string{loadBalancerId, listenerId}, tccommon.FILED_SP)
}

func splitAlbListenerId(id string) (string, string, bool) {
	parts := strings.Split(id, tccommon.FILED_SP)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", false
	}
	return parts[0], parts[1], true
}

func parseAlbListenerId(diags *diag.Diagnostics, id string) (string, string, bool) {
	loadBalancerId, listenerId, ok := splitAlbListenerId(id)
	if !ok {
		diags.AddError(
			"Invalid resource ID",
			fmt.Sprintf("Resource id %q is invalid, expected `<load_balancer_id>%s<listener_id>`.", id, tccommon.FILED_SP),
		)
	}
	return loadBalancerId, listenerId, ok
}

type albClientInterface interface {
	UseAlbClient() *albv20251030.Client
}

type int64RequiresReplace struct{}

func (int64RequiresReplace) Description(_ context.Context) string {
	return "Changing this attribute forces resource replacement."
}

func (int64RequiresReplace) MarkdownDescription(ctx context.Context) string {
	return int64RequiresReplace{}.Description(ctx)
}

func (int64RequiresReplace) PlanModifyInt64(_ context.Context, req planmodifier.Int64Request, resp *planmodifier.Int64Response) {
	if req.StateValue.IsNull() {
		return
	}
	if req.PlanValue.Equal(req.StateValue) {
		return
	}
	resp.RequiresReplace = true
}
