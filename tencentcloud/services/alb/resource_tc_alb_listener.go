package alb

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	albv20251030 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/alb/v20251030"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/fw"
	fwhelper "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

// Compile-time interface assertions. Only the extension interfaces are
// asserted; the base resource.Resource interface is checked implicitly by the
// factory return value.
var (
	_ resource.ResourceWithConfigure   = &AlbListenerResource{}
	_ resource.ResourceWithImportState = &AlbListenerResource{}
)

// NewAlbListenerResource is the factory referenced by
// tencentcloud/framework/registry.go to register this resource.
func NewAlbListenerResource() resource.Resource {
	return &AlbListenerResource{}
}

// AlbListenerResource implements resource.Resource for
// tencentcloud_alb_listener.
//
// An ALB listener is addressed by its load balancer plus the listener id on the
// cloud side (describe, modify and delete all require the load balancer id), so
// the Terraform id joins both parts with `#` and is split again for every CRUD
// call.
type AlbListenerResource struct {
	fw.ResourceWithConfigure
}

// AlbListenerResourceModel maps the top-level schema attributes. Every schema
// attribute appears here exactly once.
type AlbListenerResourceModel struct {
	Id                  types.String `tfsdk:"id"`
	LoadBalancerId      types.String `tfsdk:"load_balancer_id"`
	ListenerId          types.String `tfsdk:"listener_id"`
	ListenerName        types.String `tfsdk:"listener_name"`
	ListenerPort        types.Int64  `tfsdk:"listener_port"`
	ListenerProtocol    types.String `tfsdk:"listener_protocol"`
	DefaultActions      types.List   `tfsdk:"default_actions"`
	CaEnabled           types.Bool   `tfsdk:"ca_enabled"`
	CaCertificateIds    types.Set    `tfsdk:"ca_certificate_ids"`
	CertificateIds      types.Set    `tfsdk:"certificate_ids"`
	ClientToken         types.String `tfsdk:"client_token"`
	GzipEnabled         types.Bool   `tfsdk:"gzip_enabled"`
	Http2Enabled        types.Bool   `tfsdk:"http2enabled"`
	IdleTimeout         types.Int64  `tfsdk:"idle_timeout"`
	RequestTimeout      types.Int64  `tfsdk:"request_timeout"`
	SecurityPolicyId    types.String `tfsdk:"security_policy_id"`
	Tags                types.List   `tfsdk:"tags"`
	XForwardedForConfig types.List   `tfsdk:"x_forwarded_for_config"`
	ListenerStatus      types.String `tfsdk:"listener_status"`
	CreateTime          types.String `tfsdk:"create_time"`
	ModifyTime          types.String `tfsdk:"modify_time"`
}

// albDefaultActionModel maps one element of default_actions.
type albDefaultActionModel struct {
	TargetGroupConfig types.List   `tfsdk:"target_group_config"`
	Type              types.String `tfsdk:"type"`
}

// albTargetGroupConfigModel maps one element of
// default_actions[*].target_group_config.
type albTargetGroupConfigModel struct {
	TargetGroups             types.List `tfsdk:"target_groups"`
	TargetGroupStickySession types.List `tfsdk:"target_group_sticky_session"`
}

// albTargetGroupModel maps one element of
// default_actions[*].target_group_config[*].target_groups.
type albTargetGroupModel struct {
	TargetGroupId types.String `tfsdk:"target_group_id"`
	Weight        types.Int64  `tfsdk:"weight"`
}

// albTargetGroupStickySessionModel maps one element of
// default_actions[*].target_group_config[*].target_group_sticky_session.
type albTargetGroupStickySessionModel struct {
	Enabled types.Bool  `tfsdk:"enabled"`
	Timeout types.Int64 `tfsdk:"timeout"`
}

// albTagModel maps one element of tags.
type albTagModel struct {
	TagKey   types.String `tfsdk:"tag_key"`
	TagValue types.String `tfsdk:"tag_value"`
}

// albXForwardedForConfigModel maps one element of x_forwarded_for_config.
type albXForwardedForConfigModel struct {
	XForwardedForAlbIdEnabled         types.Bool   `tfsdk:"x_forwarded_for_alb_id_enabled"`
	XForwardedForClientSrcPortEnabled types.Bool   `tfsdk:"x_forwarded_for_client_src_port_enabled"`
	XForwardedForHostEnabled          types.Bool   `tfsdk:"x_forwarded_for_host_enabled"`
	XForwardedForMode                 types.String `tfsdk:"x_forwarded_for_mode"`
	XForwardedForPortEnabled          types.Bool   `tfsdk:"x_forwarded_for_port_enabled"`
	XForwardedForProtoEnabled         types.Bool   `tfsdk:"x_forwarded_for_proto_enabled"`
	XTencentClientIdnEnabled          types.Bool   `tfsdk:"x_tencent_client_idn_enabled"`
	XTencentClientSdnEnabled          types.Bool   `tfsdk:"x_tencent_client_sdn_enabled"`
	XTencentClientSerialEnabled       types.Bool   `tfsdk:"x_tencent_client_serial_enabled"`
	XTencentClientVerifyEnabled       types.Bool   `tfsdk:"x_tencent_client_verify_enabled"`
}

// Object types used to build the nested list/set values of the model.
var (
	albTagObjectType = types.ObjectType{AttrTypes: map[string]attr.Type{
		"tag_key":   types.StringType,
		"tag_value": types.StringType,
	}}

	albTargetGroupObjectType = types.ObjectType{AttrTypes: map[string]attr.Type{
		"target_group_id": types.StringType,
		"weight":          types.Int64Type,
	}}

	albTargetGroupStickySessionObjectType = types.ObjectType{AttrTypes: map[string]attr.Type{
		"enabled": types.BoolType,
		"timeout": types.Int64Type,
	}}

	albTargetGroupConfigObjectType = types.ObjectType{AttrTypes: map[string]attr.Type{
		"target_groups":               types.ListType{ElemType: albTargetGroupObjectType},
		"target_group_sticky_session": types.ListType{ElemType: albTargetGroupStickySessionObjectType},
	}}

	albDefaultActionObjectType = types.ObjectType{AttrTypes: map[string]attr.Type{
		"target_group_config": types.ListType{ElemType: albTargetGroupConfigObjectType},
		"type":                types.StringType,
	}}

	albXForwardedForConfigObjectType = types.ObjectType{AttrTypes: map[string]attr.Type{
		"x_forwarded_for_alb_id_enabled":          types.BoolType,
		"x_forwarded_for_client_src_port_enabled": types.BoolType,
		"x_forwarded_for_host_enabled":            types.BoolType,
		"x_forwarded_for_mode":                    types.StringType,
		"x_forwarded_for_port_enabled":            types.BoolType,
		"x_forwarded_for_proto_enabled":           types.BoolType,
		"x_tencent_client_idn_enabled":            types.BoolType,
		"x_tencent_client_sdn_enabled":            types.BoolType,
		"x_tencent_client_serial_enabled":         types.BoolType,
		"x_tencent_client_verify_enabled":         types.BoolType,
	}}
)

func (r *AlbListenerResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "tencentcloud_alb_listener"
}

func (r *AlbListenerResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a resource to manage an Application Load Balancer (ALB) listener.",
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
				Description: "ID of the ALB instance this listener belongs to, in the form `alb-xxxxxxxx`. It cannot be changed after creation.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"listener_id": schema.StringAttribute{
				Computed:    true,
				Description: "ID of the listener, in the form `lst-xxxxxxxx`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"listener_name": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Name of the listener. Length is 1 to 255 characters; it can contain Chinese characters, letters, digits, hyphens (-), slashes (/), periods (.) and underscores (_).",
			},
			"listener_port": schema.Int64Attribute{
				Required:    true,
				Description: "Front-end port used by the listener. Valid values: 1 to 65535. It cannot be changed after creation.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"listener_protocol": schema.StringAttribute{
				Required:    true,
				Description: "Protocol of the listener. Valid values: `HTTP`, `HTTPS` and `QUIC`. It cannot be changed after creation.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"default_actions": schema.ListNestedAttribute{
				Required:    true,
				Description: "List of default forwarding actions. Currently a listener supports exactly one default forwarding action, and the only supported action type is forwarding to a target group.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"target_group_config": schema.ListNestedAttribute{
							Required:    true,
							Description: "Target group configuration of the forwarding action. Only a single target group is supported when the listener is created.",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"target_groups": schema.ListNestedAttribute{
										Required:    true,
										Description: "List of target groups.",
										NestedObject: schema.NestedAttributeObject{
											Attributes: map[string]schema.Attribute{
												"target_group_id": schema.StringAttribute{
													Required:    true,
													Description: "ID of the target group, in the form `lbtg-xxxxxxxx`.",
												},
												"weight": schema.Int64Attribute{
													Optional:    true,
													Computed:    true,
													Description: "Weight of the target group. Valid values: 0 to 100. Default value: 10.",
												},
											},
										},
									},
									"target_group_sticky_session": schema.ListNestedAttribute{
										Optional:    true,
										Computed:    true,
										Description: "Session persistence between the client and the target group.",
										NestedObject: schema.NestedAttributeObject{
											Attributes: map[string]schema.Attribute{
												"enabled": schema.BoolAttribute{
													Required:    true,
													Description: "Whether session persistence is enabled. It is disabled by default.",
												},
												"timeout": schema.Int64Attribute{
													Optional:    true,
													Computed:    true,
													Description: "Session persistence timeout in seconds. Valid values: 1 to 86400. Default value: 1000.",
												},
											},
										},
									},
								},
							},
						},
						"type": schema.StringAttribute{
							Required:    true,
							Description: "Type of the forwarding action. Only `TargetGroup` (forward to a target group) is supported for the default action of a listener.",
						},
					},
				},
			},
			"ca_enabled": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Whether mutual authentication (two-way authentication) is enabled. Default value: `false`.",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"ca_certificate_ids": schema.SetAttribute{
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				Description: "List of CA certificate IDs configured on the listener. Currently only one CA certificate is supported, and it is required when `ca_enabled` is `true`.",
			},
			"certificate_ids": schema.SetAttribute{
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				Description: "List of server certificate IDs of the listener.",
			},
			"client_token": schema.StringAttribute{
				Optional:    true,
				Description: "Client token used to guarantee the idempotency of the create request. It must be unique among different requests and can only contain ASCII characters. It is only used when the listener is created.",
			},
			"gzip_enabled": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Whether Gzip compression is enabled. Default value: `true`.",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"http2enabled": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Whether HTTP/2 is enabled. It is only supported by the `HTTPS` protocol. Default value: `false` for HTTP and `true` for HTTPS.",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"idle_timeout": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "Connection idle timeout in seconds. Valid values: 1 to 600. Default value: 15.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"request_timeout": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "Connection request timeout in seconds. Valid values: 1 to 600. Default value: 60.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"security_policy_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "ID of the TLS security policy, in the form `tls-xxxxxxxx`.",
			},
			"tags": schema.ListNestedAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Tags of the listener. Up to 20 tags are supported. The listener modify API does not accept tags, so changing them forces the listener to be recreated.",
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
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
					listplanmodifier.RequiresReplace(),
				},
			},
			"x_forwarded_for_config": schema.ListNestedAttribute{
				Optional:    true,
				Computed:    true,
				Description: "X-Forwarded-For configuration of the listener.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"x_forwarded_for_alb_id_enabled": schema.BoolAttribute{
							Optional:    true,
							Computed:    true,
							Description: "Whether the ALB instance ID is obtained through the ALB-ID header.",
							PlanModifiers: []planmodifier.Bool{
								boolplanmodifier.UseStateForUnknown(),
							},
						},
						"x_forwarded_for_client_src_port_enabled": schema.BoolAttribute{
							Optional:    true,
							Computed:    true,
							Description: "Whether the client port accessing the ALB instance is obtained through the X-Forwarded-Client-srcport header.",
							PlanModifiers: []planmodifier.Bool{
								boolplanmodifier.UseStateForUnknown(),
							},
						},
						"x_forwarded_for_host_enabled": schema.BoolAttribute{
							Optional:    true,
							Computed:    true,
							Description: "Whether the domain name of the client accessing the ALB instance is obtained through the X-Forwarded-Host header.",
							PlanModifiers: []planmodifier.Bool{
								boolplanmodifier.UseStateForUnknown(),
							},
						},
						"x_forwarded_for_mode": schema.StringAttribute{
							Optional:    true,
							Computed:    true,
							Description: "How the X-Forwarded-For (XFF) HTTP header is handled. Valid values: `append` (default, append the client IP to the end of the header), `remove` (drop the header) and `passthrough` (forward the header unchanged).",
						},
						"x_forwarded_for_port_enabled": schema.BoolAttribute{
							Optional:    true,
							Computed:    true,
							Description: "Whether the listener port of the ALB instance is obtained through the X-Forwarded-Port header.",
							PlanModifiers: []planmodifier.Bool{
								boolplanmodifier.UseStateForUnknown(),
							},
						},
						"x_forwarded_for_proto_enabled": schema.BoolAttribute{
							Optional:    true,
							Computed:    true,
							Description: "Whether the listener protocol of the ALB instance is obtained through the X-Forwarded-Proto header.",
							PlanModifiers: []planmodifier.Bool{
								boolplanmodifier.UseStateForUnknown(),
							},
						},
						"x_tencent_client_idn_enabled": schema.BoolAttribute{
							Optional:    true,
							Computed:    true,
							Description: "Whether the issuer `$ssl_client_i_dn` of the client certificate is obtained through the X-Tencent-Client-IDN header.",
							PlanModifiers: []planmodifier.Bool{
								boolplanmodifier.UseStateForUnknown(),
							},
						},
						"x_tencent_client_sdn_enabled": schema.BoolAttribute{
							Optional:    true,
							Computed:    true,
							Description: "Whether the subject `$ssl_client_s_dn` of the client certificate is obtained through the X-Tencent-Client-SDN header.",
							PlanModifiers: []planmodifier.Bool{
								boolplanmodifier.UseStateForUnknown(),
							},
						},
						"x_tencent_client_serial_enabled": schema.BoolAttribute{
							Optional:    true,
							Computed:    true,
							Description: "Whether the serial number `$ssl_client_serial` of the client certificate is obtained through the X-Tencent-Client-Serial header.",
							PlanModifiers: []planmodifier.Bool{
								boolplanmodifier.UseStateForUnknown(),
							},
						},
						"x_tencent_client_verify_enabled": schema.BoolAttribute{
							Optional:    true,
							Computed:    true,
							Description: "Whether the verification result `$ssl_client_verify` of the client certificate is obtained through the X-Tencent-Client-Verify header.",
							PlanModifiers: []planmodifier.Bool{
								boolplanmodifier.UseStateForUnknown(),
							},
						},
					},
				},
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
			},
			"listener_status": schema.StringAttribute{
				Computed:    true,
				Description: "Status of the listener. Valid values: `Active` (running), `Provisioning` (being created), `Configuring` (being modified) and `ProvisionFailed` (failed to be created).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
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
				Description: "Last modification time of the listener in ISO 8601 format.",
			},
		},
	}
}

// Create handles `terraform apply` for new resources.
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

	createListenerRequest := albv20251030.NewCreateListenerRequest()
	createListenerRequest.LoadBalancerId = helper.String(loadBalancerId)
	createListenerRequest.ListenerPort = fwhelper.Uint64PointerFromValue(plan.ListenerPort)
	createListenerRequest.ListenerProtocol = fwhelper.StringPointerFromValue(plan.ListenerProtocol)
	createListenerRequest.ListenerName = fwhelper.StringPointerFromValue(plan.ListenerName)
	createListenerRequest.CaEnabled = fwhelper.BoolPointerFromValue(plan.CaEnabled)
	createListenerRequest.GzipEnabled = fwhelper.BoolPointerFromValue(plan.GzipEnabled)
	createListenerRequest.Http2Enabled = fwhelper.BoolPointerFromValue(plan.Http2Enabled)
	createListenerRequest.IdleTimeout = fwhelper.Uint64PointerFromValue(plan.IdleTimeout)
	createListenerRequest.RequestTimeout = fwhelper.Uint64PointerFromValue(plan.RequestTimeout)
	createListenerRequest.SecurityPolicyId = fwhelper.StringPointerFromValue(plan.SecurityPolicyId)
	// client_token is an idempotency key of the create call only: it is fixed
	// for the whole apply so a retried CreateListener cannot create the same
	// listener twice.
	createListenerRequest.ClientToken = fwhelper.StringPointerFromValue(plan.ClientToken)

	defaultActions, defaultActionDiags := expandAlbDefaultActions(ctx, plan.DefaultActions)
	resp.Diagnostics.Append(defaultActionDiags...)
	if resp.Diagnostics.HasError() {
		return
	}
	createListenerRequest.DefaultActions = defaultActions

	caCertificateIds, caDiags := expandStringSet(ctx, plan.CaCertificateIds)
	resp.Diagnostics.Append(caDiags...)
	if resp.Diagnostics.HasError() {
		return
	}
	createListenerRequest.CaCertificateIds = caCertificateIds

	certificateIds, certificateDiags := expandStringSet(ctx, plan.CertificateIds)
	resp.Diagnostics.Append(certificateDiags...)
	if resp.Diagnostics.HasError() {
		return
	}
	createListenerRequest.CertificateIds = certificateIds

	tags, tagDiags := expandAlbTags(ctx, plan.Tags)
	resp.Diagnostics.Append(tagDiags...)
	if resp.Diagnostics.HasError() {
		return
	}
	createListenerRequest.Tags = tags

	xForwardedForConfig, xffDiags := expandAlbXForwardedForConfig(ctx, plan.XForwardedForConfig)
	resp.Diagnostics.Append(xffDiags...)
	if resp.Diagnostics.HasError() {
		return
	}
	createListenerRequest.XForwardedForConfig = xForwardedForConfig

	listenerId, err := service.CreateListener(ctx, createListenerRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating ALB listener",
			fmt.Sprintf("CreateListener failed (load_balancer_id=%q): %s", loadBalancerId, err),
		)
		return
	}
	if listenerId == "" {
		resp.Diagnostics.AddError(
			"Error creating ALB listener",
			"CreateListener returned an empty listener id.",
		)
		return
	}

	plan.ListenerId = types.StringValue(listenerId)
	plan.Id = types.StringValue(albListenerId(loadBalancerId, listenerId))

	// The listener is created asynchronously, so wait until it is settled
	// before reading it back: a state written while the change is still in
	// flight makes Terraform report an inconsistent result after apply.
	if err := service.WaitListenerSettled(ctx, loadBalancerId, listenerId, tccommon.WriteRetryTimeout); err != nil {
		resp.Diagnostics.AddError(
			"Error waiting for ALB listener creation",
			fmt.Sprintf(
				"The listener was created but never became active (%s). "+
					"Import it with `terraform import tencentcloud_alb_listener.<name> %s` to keep it managed: %s",
				albListenerScope(loadBalancerId, listenerId), plan.Id.ValueString(), err,
			),
		)
		return
	}

	describeListenersRequest := albv20251030.NewDescribeListenersRequest()
	describeListenersRequest.LoadBalancerId = helper.String(loadBalancerId)
	describeListenersRequest.ListenerIds = []*string{helper.String(listenerId)}

	listeners, _, err := service.DescribeListeners(ctx, describeListenersRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading ALB listener after create",
			fmt.Sprintf("CreateListener succeeded but the listener cannot be read back (%s): %s",
				albListenerScope(loadBalancerId, listenerId), err),
		)
		return
	}

	item := findAlbListener(listeners, listenerId)
	if item == nil {
		resp.Diagnostics.AddError(
			"ALB listener not readable after create",
			fmt.Sprintf(
				"Listener %s was created but cannot be read back (%s). "+
					"Import it with `terraform import tencentcloud_alb_listener.<name> %s` to keep it managed.",
				listenerId, albListenerScope(loadBalancerId, listenerId), plan.Id.ValueString(),
			),
		)
		return
	}

	// DescribeListeners does not report default_actions, ca_certificate_ids or
	// certificate_ids; DescribeListenerDetail is the only read path that does.
	// A nil detail means the listener is already gone, which the caller handles
	// as "keep what is already in state".
	detail, err := service.DescribeListenerDetailById(ctx, loadBalancerId, listenerId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading ALB listener detail",
			fmt.Sprintf("DescribeListenerDetail failed (%s): %s", albListenerScope(loadBalancerId, listenerId), err),
		)
		return
	}

	resp.Diagnostics.Append(refreshAlbListenerModel(ctx, &plan, loadBalancerId, item, detail)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Read handles `terraform refresh` and post-apply state refresh.
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

	describeListenersRequest := albv20251030.NewDescribeListenersRequest()
	describeListenersRequest.LoadBalancerId = helper.String(loadBalancerId)
	describeListenersRequest.ListenerIds = []*string{helper.String(listenerId)}

	listeners, _, err := service.DescribeListeners(ctx, describeListenersRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading ALB listener",
			fmt.Sprintf("DescribeListeners failed (%s): %s", albListenerScope(loadBalancerId, listenerId), err),
		)
		return
	}

	item := findAlbListener(listeners, listenerId)
	// A confirmed "not found" means the listener was removed outside Terraform,
	// so it is dropped from state instead of failing the refresh.
	if item == nil {
		resp.Diagnostics.AddWarning(
			"Resource not found",
			fmt.Sprintf("ALB listener %s no longer exists, removing it from state.", listenerId),
		)
		resp.State.RemoveResource(ctx)
		return
	}

	// DescribeListeners does not report default_actions, ca_certificate_ids or
	// certificate_ids; DescribeListenerDetail is the only read path that does.
	// A nil detail means the listener is already gone, which the caller handles
	// as "keep what is already in state".
	detail, err := service.DescribeListenerDetailById(ctx, loadBalancerId, listenerId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading ALB listener detail",
			fmt.Sprintf("DescribeListenerDetail failed (%s): %s", albListenerScope(loadBalancerId, listenerId), err),
		)
		return
	}

	resp.Diagnostics.Append(refreshAlbListenerModel(ctx, &state, loadBalancerId, item, detail)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update handles in-place updates.
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

	modifyListenerAttributesRequest := albv20251030.NewModifyListenerAttributesRequest()
	modifyListenerAttributesRequest.ListenerId = helper.String(listenerId)
	modifyListenerAttributesRequest.LoadBalancerId = helper.String(loadBalancerId)
	modifyListenerAttributesRequest.ListenerName = fwhelper.StringPointerFromValue(plan.ListenerName)
	modifyListenerAttributesRequest.CaEnabled = fwhelper.BoolPointerFromValue(plan.CaEnabled)
	modifyListenerAttributesRequest.GzipEnabled = fwhelper.BoolPointerFromValue(plan.GzipEnabled)
	modifyListenerAttributesRequest.Http2Enabled = fwhelper.BoolPointerFromValue(plan.Http2Enabled)
	modifyListenerAttributesRequest.IdleTimeout = fwhelper.Uint64PointerFromValue(plan.IdleTimeout)
	modifyListenerAttributesRequest.RequestTimeout = fwhelper.Uint64PointerFromValue(plan.RequestTimeout)
	modifyListenerAttributesRequest.SecurityPolicyId = fwhelper.StringPointerFromValue(plan.SecurityPolicyId)

	defaultActions, defaultActionDiags := expandAlbDefaultActions(ctx, plan.DefaultActions)
	resp.Diagnostics.Append(defaultActionDiags...)
	if resp.Diagnostics.HasError() {
		return
	}
	modifyListenerAttributesRequest.DefaultActions = defaultActions

	caCertificateIds, caDiags := expandStringSet(ctx, plan.CaCertificateIds)
	resp.Diagnostics.Append(caDiags...)
	if resp.Diagnostics.HasError() {
		return
	}
	modifyListenerAttributesRequest.CaCertificateIds = caCertificateIds

	certificateIds, certificateDiags := expandStringSet(ctx, plan.CertificateIds)
	resp.Diagnostics.Append(certificateDiags...)
	if resp.Diagnostics.HasError() {
		return
	}
	modifyListenerAttributesRequest.CertificateIds = certificateIds

	xForwardedForConfig, xffDiags := expandAlbXForwardedForConfig(ctx, plan.XForwardedForConfig)
	resp.Diagnostics.Append(xffDiags...)
	if resp.Diagnostics.HasError() {
		return
	}
	modifyListenerAttributesRequest.XForwardedForConfig = xForwardedForConfig

	if err := service.ModifyListenerAttributes(ctx, modifyListenerAttributesRequest); err != nil {
		resp.Diagnostics.AddError(
			"Error updating ALB listener",
			fmt.Sprintf("ModifyListenerAttributes failed (%s): %s", albListenerScope(loadBalancerId, listenerId), err),
		)
		return
	}

	if err := service.WaitListenerSettled(ctx, loadBalancerId, listenerId, tccommon.WriteRetryTimeout); err != nil {
		resp.Diagnostics.AddError(
			"Error waiting for ALB listener update",
			fmt.Sprintf("The listener was modified but never became active (%s): %s",
				albListenerScope(loadBalancerId, listenerId), err),
		)
		return
	}

	describeListenersRequest := albv20251030.NewDescribeListenersRequest()
	describeListenersRequest.LoadBalancerId = helper.String(loadBalancerId)
	describeListenersRequest.ListenerIds = []*string{helper.String(listenerId)}

	listeners, _, err := service.DescribeListeners(ctx, describeListenersRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading ALB listener after update",
			fmt.Sprintf("ModifyListenerAttributes succeeded but the listener cannot be read back (%s): %s",
				albListenerScope(loadBalancerId, listenerId), err),
		)
		return
	}

	item := findAlbListener(listeners, listenerId)
	if item == nil {
		resp.Diagnostics.AddError(
			"ALB listener not readable after update",
			fmt.Sprintf("Listener %s was updated but cannot be read back (%s).",
				listenerId, albListenerScope(loadBalancerId, listenerId)),
		)
		return
	}

	// DescribeListeners does not report default_actions, ca_certificate_ids or
	// certificate_ids; DescribeListenerDetail is the only read path that does.
	// A nil detail means the listener is already gone, which the caller handles
	// as "keep what is already in state".
	detail, err := service.DescribeListenerDetailById(ctx, loadBalancerId, listenerId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading ALB listener detail",
			fmt.Sprintf("DescribeListenerDetail failed (%s): %s", albListenerScope(loadBalancerId, listenerId), err),
		)
		return
	}

	plan.Id = state.Id
	resp.Diagnostics.Append(refreshAlbListenerModel(ctx, &plan, loadBalancerId, item, detail)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete handles resource deletion.
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

	deleteListenerRequest := albv20251030.NewDeleteListenerRequest()
	deleteListenerRequest.LoadBalancerId = helper.String(loadBalancerId)
	deleteListenerRequest.ListenerIds = []*string{helper.String(listenerId)}

	if err := service.DeleteListener(ctx, deleteListenerRequest); err != nil {
		// A listener that is already gone is a successful deletion.
		if tccommon.IsExpectError(err, []string{
			"ResourceNotFound.Listener",
			"ResourceNotFound",
		}) {
			return
		}
		resp.Diagnostics.AddError(
			"Error deleting ALB listener",
			fmt.Sprintf("DeleteListener failed (%s): %s", albListenerScope(loadBalancerId, listenerId), err),
		)
	}
}

// ImportState supports `terraform import` using the
// `<load_balancer_id>#<listener_id>` form. The remaining attributes are filled
// in by the Read that follows.
func (r *AlbListenerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	loadBalancerId, listenerId, ok := splitAlbListenerId(req.ID)
	if !ok {
		resp.Diagnostics.AddError(
			"Invalid import ID",
			fmt.Sprintf("Import id %q is invalid, expected `<load_balancer_id>%s<listener_id>`.", req.ID, tccommon.FILED_SP),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("load_balancer_id"), loadBalancerId)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("listener_id"), listenerId)...)
}

// service returns the ALB service bound to the shared client. When the provider
// has not been configured it records a diagnostic and reports false, so the
// CRUD methods never dereference a nil SDK client.
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

// findAlbListener picks this resource's own item from a describe list. "No
// match" means the listener no longer exists.
func findAlbListener(listeners []*albv20251030.ListenerOutput, listenerId string) *albv20251030.ListenerOutput {
	for _, listener := range listeners {
		if listener != nil && listener.ListenerId != nil && *listener.ListenerId == listenerId {
			return listener
		}
	}
	return nil
}

// albListenerId builds the Terraform id from its two cloud parts.
func albListenerId(loadBalancerId, listenerId string) string {
	return strings.Join([]string{loadBalancerId, listenerId}, tccommon.FILED_SP)
}

// splitAlbListenerId parses `<load_balancer_id>#<listener_id>`.
func splitAlbListenerId(id string) (string, string, bool) {
	parts := strings.Split(id, tccommon.FILED_SP)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", false
	}
	return parts[0], parts[1], true
}

// parseAlbListenerId resolves the Terraform id into the pair of ids the ALB API
// expects, recording a diagnostic when the id is malformed.
func parseAlbListenerId(diags *diag.Diagnostics, id string) (loadBalancerId, listenerId string, ok bool) {
	loadBalancerId, listenerId, ok = splitAlbListenerId(id)
	if !ok {
		diags.AddError(
			"Invalid resource ID",
			fmt.Sprintf("Resource id %q is invalid, expected `<load_balancer_id>%s<listener_id>`.", id, tccommon.FILED_SP),
		)
	}
	return
}

// albListenerScope renders the cloud-side identity of a listener for
// diagnostics.
func albListenerScope(loadBalancerId, listenerId string) string {
	if listenerId == "" {
		return fmt.Sprintf("load_balancer_id=%q", loadBalancerId)
	}
	return fmt.Sprintf("load_balancer_id=%q, listener_id=%q", loadBalancerId, listenerId)
}

// refreshAlbListenerModel copies the cloud state into the model.
//
// DescribeListeners does not echo back default_actions, ca_certificate_ids and
// certificate_ids, so those are taken from the listener detail response when it
// is available. Fields the API omits keep their previous value so a configured
// value is never nulled out after apply.
func refreshAlbListenerModel(ctx context.Context, m *AlbListenerResourceModel, loadBalancerId string, item *albv20251030.ListenerOutput, detail *albv20251030.DescribeListenerDetailResponseParams) diag.Diagnostics {
	var diags diag.Diagnostics

	if item == nil {
		diags.AddError(
			"Malformed ALB listener response",
			fmt.Sprintf("DescribeListeners returned an empty listener (%s).", albListenerScope(loadBalancerId, "")),
		)
		return diags
	}

	m.LoadBalancerId = types.StringValue(loadBalancerId)
	m.ListenerId = fwhelper.StringValueOrPrev(item.ListenerId, m.ListenerId)
	m.Id = types.StringValue(albListenerId(loadBalancerId, m.ListenerId.ValueString()))

	m.ListenerName = fwhelper.StringValueOrPrev(item.ListenerName, m.ListenerName)
	m.ListenerPort = fwhelper.Int64ValueOrPrev(item.ListenerPort, m.ListenerPort)
	m.ListenerProtocol = fwhelper.StringValueOrPrev(item.ListenerProtocol, m.ListenerProtocol)
	m.CaEnabled = boolValueOrPrev(item.CaEnable, m.CaEnabled)
	m.GzipEnabled = boolValueOrPrev(item.GzipEnabled, m.GzipEnabled)
	m.Http2Enabled = boolValueOrPrev(item.Http2Enable, m.Http2Enabled)
	m.IdleTimeout = fwhelper.Int64ValueOrPrev(item.IdleTimeout, m.IdleTimeout)
	m.RequestTimeout = fwhelper.Int64ValueOrPrev(item.RequestTimeout, m.RequestTimeout)
	m.SecurityPolicyId = fwhelper.StringValueOrPrev(item.TlsSecurityPolicyId, m.SecurityPolicyId)
	m.ListenerStatus = fwhelper.StringValueOrPrev(item.ListenerStatus, m.ListenerStatus)
	m.CreateTime = fwhelper.StringValueOrPrev(item.CreateTime, m.CreateTime)
	m.ModifyTime = fwhelper.StringValueOrPrev(item.ModifyTime, m.ModifyTime)

	tags, tagDiags := flattenAlbTags(ctx, item.Tags, m.Tags)
	diags.Append(tagDiags...)
	m.Tags = tags

	xForwardedForConfig, xffDiags := flattenAlbXForwardedForConfig(ctx, item.XForwardedForConfig, m.XForwardedForConfig)
	diags.Append(xffDiags...)
	m.XForwardedForConfig = xForwardedForConfig

	// default_actions, ca_certificate_ids and certificate_ids are only
	// reported by DescribeListenerDetail. A plan may hold an unknown value for
	// an unconfigured Optional+Computed attribute and state must never hold
	// unknowns, so those are normalized to null when the cloud does not report
	// them.
	if detail != nil {
		defaultActions, defaultActionDiags := flattenAlbDefaultActions(ctx, detail.DefaultActions, m.DefaultActions)
		diags.Append(defaultActionDiags...)
		m.DefaultActions = defaultActions

		caCertificateIds, caDiags := flattenAlbStringSet(ctx, detail.CaCertificateIds, m.CaCertificateIds)
		diags.Append(caDiags...)
		m.CaCertificateIds = caCertificateIds

		certificateIds, certificateDiags := flattenAlbStringSet(ctx, detail.CertificateIds, m.CertificateIds)
		diags.Append(certificateDiags...)
		m.CertificateIds = certificateIds
	} else {
		if m.CaCertificateIds.IsUnknown() {
			m.CaCertificateIds = types.SetNull(types.StringType)
		}
		if m.CertificateIds.IsUnknown() {
			m.CertificateIds = types.SetNull(types.StringType)
		}
	}

	if diags.HasError() {
		return diags
	}

	// Required attributes must never be written as null, otherwise Terraform
	// rejects the refreshed object.
	if m.ListenerPort.IsNull() || m.ListenerProtocol.IsNull() || m.DefaultActions.IsNull() {
		diags.AddError(
			"Malformed ALB listener response",
			fmt.Sprintf("DescribeListeners returned no listener_port, listener_protocol or default_actions (%s).",
				albListenerScope(loadBalancerId, m.ListenerId.ValueString())),
		)
	}

	return diags
}

// flattenAlbDefaultActions converts the cloud action list into the
// default_actions list value.
//
// weight, target_group_sticky_session and timeout are Optional+Computed, so a
// value the practitioner configured is restored from prevList when the cloud
// omits it: writing null for a configured value makes Terraform report
// "provider produced inconsistent result after apply".
func flattenAlbDefaultActions(ctx context.Context, items []*albv20251030.DefaultAction, prevList types.List) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics

	if len(items) == 0 {
		if !prevList.IsNull() && !prevList.IsUnknown() {
			return prevList, diags
		}
		return types.ListNull(albDefaultActionObjectType), diags
	}

	var prevActions []albDefaultActionModel
	if !prevList.IsNull() && !prevList.IsUnknown() {
		diags.Append(prevList.ElementsAs(ctx, &prevActions, false)...)
		if diags.HasError() {
			return types.ListNull(albDefaultActionObjectType), diags
		}
	}

	models := make([]albDefaultActionModel, 0, len(items))
	for index, item := range items {
		if item == nil {
			continue
		}

		var prev albDefaultActionModel
		if index < len(prevActions) {
			prev = prevActions[index]
		}

		targetGroupConfigs, configDiags := flattenAlbTargetGroupConfigs(ctx, item.TargetGroupConfig, prev.TargetGroupConfig)
		diags.Append(configDiags...)
		if diags.HasError() {
			return types.ListNull(albDefaultActionObjectType), diags
		}

		models = append(models, albDefaultActionModel{
			TargetGroupConfig: targetGroupConfigs,
			Type:              fwhelper.StringValueOrPrev(item.Type, prev.Type),
		})
	}

	if len(models) == 0 {
		return types.ListNull(albDefaultActionObjectType), diags
	}

	list, listDiags := types.ListValueFrom(ctx, albDefaultActionObjectType, models)
	diags.Append(listDiags...)
	return list, diags
}

// flattenAlbTargetGroupConfigs converts the cloud target group configuration
// into the target_group_config list value.
func flattenAlbTargetGroupConfigs(ctx context.Context, config *albv20251030.TargetGroupConfig, prevList types.List) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics

	if config == nil {
		return types.ListNull(albTargetGroupConfigObjectType), diags
	}

	var prevConfigs []albTargetGroupConfigModel
	if !prevList.IsNull() && !prevList.IsUnknown() {
		diags.Append(prevList.ElementsAs(ctx, &prevConfigs, false)...)
		if diags.HasError() {
			return types.ListNull(albTargetGroupConfigObjectType), diags
		}
	}

	var prev albTargetGroupConfigModel
	if len(prevConfigs) > 0 {
		prev = prevConfigs[0]
	}

	targetGroups, targetGroupDiags := flattenAlbTargetGroups(ctx, config.TargetGroups, prev.TargetGroups)
	diags.Append(targetGroupDiags...)
	if diags.HasError() {
		return types.ListNull(albTargetGroupConfigObjectType), diags
	}

	stickySessions, stickyDiags := flattenAlbTargetGroupStickySessions(ctx, config.TargetGroupStickySession, prev.TargetGroupStickySession)
	diags.Append(stickyDiags...)
	if diags.HasError() {
		return types.ListNull(albTargetGroupConfigObjectType), diags
	}

	list, listDiags := types.ListValueFrom(ctx, albTargetGroupConfigObjectType, []albTargetGroupConfigModel{
		{
			TargetGroups:             targetGroups,
			TargetGroupStickySession: stickySessions,
		},
	})
	diags.Append(listDiags...)
	return list, diags
}

// flattenAlbTargetGroups converts the cloud target group list into the
// target_groups list value. An element the cloud omits keeps the configured
// weight from the previous value and, when the cloud reports no element at
// all, the whole previous list is kept.
func flattenAlbTargetGroups(ctx context.Context, items []*albv20251030.TargetGroupTuple, prevList types.List) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics

	var prevGroups []albTargetGroupModel
	if !prevList.IsNull() && !prevList.IsUnknown() {
		diags.Append(prevList.ElementsAs(ctx, &prevGroups, false)...)
		if diags.HasError() {
			return types.ListNull(albTargetGroupObjectType), diags
		}
	}

	models := make([]albTargetGroupModel, 0, len(items))
	for index, item := range items {
		if item == nil {
			continue
		}

		var prev albTargetGroupModel
		if index < len(prevGroups) {
			prev = prevGroups[index]
		}

		models = append(models, albTargetGroupModel{
			TargetGroupId: fwhelper.StringPointerValueOrNull(item.TargetGroupId),
			Weight:        int64ValueOrPrev(item.Weight, prev.Weight),
		})
	}

	if len(models) == 0 {
		if !prevList.IsNull() && !prevList.IsUnknown() {
			return prevList, diags
		}
		return types.ListNull(albTargetGroupObjectType), diags
	}

	list, listDiags := types.ListValueFrom(ctx, albTargetGroupObjectType, models)
	diags.Append(listDiags...)
	return list, diags
}

// flattenAlbTargetGroupStickySessions converts the cloud sticky session
// configuration into the target_group_sticky_session list value. When the cloud
// omits the object, the configured value is kept.
func flattenAlbTargetGroupStickySessions(ctx context.Context, session *albv20251030.TargetGroupStickySession, prevList types.List) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics

	if session == nil {
		if !prevList.IsNull() && !prevList.IsUnknown() {
			return prevList, diags
		}
		return types.ListNull(albTargetGroupStickySessionObjectType), diags
	}

	var prevSessions []albTargetGroupStickySessionModel
	if !prevList.IsNull() && !prevList.IsUnknown() {
		diags.Append(prevList.ElementsAs(ctx, &prevSessions, false)...)
		if diags.HasError() {
			return types.ListNull(albTargetGroupStickySessionObjectType), diags
		}
	}

	var prev albTargetGroupStickySessionModel
	if len(prevSessions) > 0 {
		prev = prevSessions[0]
	}

	list, listDiags := types.ListValueFrom(ctx, albTargetGroupStickySessionObjectType, []albTargetGroupStickySessionModel{
		{
			Enabled: boolValueOrPrev(session.Enabled, prev.Enabled),
			Timeout: int64ValueOrPrev(session.Timeout, prev.Timeout),
		},
	})
	diags.Append(listDiags...)
	return list, diags
}

// flattenAlbStringSet converts a cloud string list into a set value. When the
// cloud reports no element at all, the previous set is kept so a configured
// value is not nulled out.
func flattenAlbStringSet(_ context.Context, items []*string, prevSet types.Set) (types.Set, diag.Diagnostics) {
	var diags diag.Diagnostics

	values := make([]attr.Value, 0, len(items))
	for _, item := range items {
		if item == nil {
			continue
		}
		values = append(values, types.StringValue(*item))
	}

	if len(values) == 0 {
		if !prevSet.IsNull() && !prevSet.IsUnknown() {
			return prevSet, diags
		}
		return types.SetNull(types.StringType), diags
	}

	set, setDiags := types.SetValue(types.StringType, values)
	diags.Append(setDiags...)
	return set, diags
}

// expandAlbDefaultActions converts the planned default_actions list into the
// SDK action slice.
func expandAlbDefaultActions(ctx context.Context, list types.List) ([]*albv20251030.DefaultAction, diag.Diagnostics) {
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
		targetGroupConfigs, configDiags := expandAlbTargetGroupConfigs(ctx, action.TargetGroupConfig)
		diags.Append(configDiags...)
		if diags.HasError() {
			return nil, diags
		}

		defaultAction := &albv20251030.DefaultAction{
			Type: fwhelper.StringPointerFromValue(action.Type),
		}
		if len(targetGroupConfigs) > 0 {
			defaultAction.TargetGroupConfig = targetGroupConfigs[0]
		}
		out = append(out, defaultAction)
	}

	return out, diags
}

// expandAlbTargetGroupConfigs converts the planned target_group_config list
// into the SDK config slice. The API accepts a single config object, so the
// slice is only used to keep the caller symmetric with the list schema.
func expandAlbTargetGroupConfigs(ctx context.Context, list types.List) ([]*albv20251030.TargetGroupConfig, diag.Diagnostics) {
	var diags diag.Diagnostics
	if list.IsNull() || list.IsUnknown() {
		return nil, diags
	}

	var configs []albTargetGroupConfigModel
	diags.Append(list.ElementsAs(ctx, &configs, false)...)
	if diags.HasError() {
		return nil, diags
	}

	out := make([]*albv20251030.TargetGroupConfig, 0, len(configs))
	for _, config := range configs {
		targetGroups, targetDiags := expandAlbTargetGroups(ctx, config.TargetGroups)
		diags.Append(targetDiags...)
		if diags.HasError() {
			return nil, diags
		}

		stickySessions, stickyDiags := expandAlbTargetGroupStickySessions(ctx, config.TargetGroupStickySession)
		diags.Append(stickyDiags...)
		if diags.HasError() {
			return nil, diags
		}

		targetGroupConfig := &albv20251030.TargetGroupConfig{
			TargetGroups: targetGroups,
		}
		if len(stickySessions) > 0 {
			targetGroupConfig.TargetGroupStickySession = stickySessions[0]
		}
		out = append(out, targetGroupConfig)
	}

	return out, diags
}

// expandAlbTargetGroups converts the planned target_groups list into the SDK
// tuple slice.
func expandAlbTargetGroups(ctx context.Context, list types.List) ([]*albv20251030.TargetGroupTuple, diag.Diagnostics) {
	var diags diag.Diagnostics
	if list.IsNull() || list.IsUnknown() {
		return nil, diags
	}

	var groups []albTargetGroupModel
	diags.Append(list.ElementsAs(ctx, &groups, false)...)
	if diags.HasError() {
		return nil, diags
	}

	out := make([]*albv20251030.TargetGroupTuple, 0, len(groups))
	for _, group := range groups {
		out = append(out, &albv20251030.TargetGroupTuple{
			TargetGroupId: fwhelper.StringPointerFromValue(group.TargetGroupId),
			Weight:        fwhelper.Int64PointerFromValue(group.Weight),
		})
	}

	return out, diags
}

// expandAlbTargetGroupStickySessions converts the planned
// target_group_sticky_session list into the SDK slice.
func expandAlbTargetGroupStickySessions(ctx context.Context, list types.List) ([]*albv20251030.TargetGroupStickySession, diag.Diagnostics) {
	var diags diag.Diagnostics
	if list.IsNull() || list.IsUnknown() {
		return nil, diags
	}

	var sessions []albTargetGroupStickySessionModel
	diags.Append(list.ElementsAs(ctx, &sessions, false)...)
	if diags.HasError() {
		return nil, diags
	}

	out := make([]*albv20251030.TargetGroupStickySession, 0, len(sessions))
	for _, session := range sessions {
		out = append(out, &albv20251030.TargetGroupStickySession{
			Enabled: fwhelper.BoolPointerFromValue(session.Enabled),
			Timeout: fwhelper.Int64PointerFromValue(session.Timeout),
		})
	}

	return out, diags
}

// expandAlbTags converts the planned tags list into the SDK tag slice.
func expandAlbTags(ctx context.Context, list types.List) ([]*albv20251030.TagInfo, diag.Diagnostics) {
	var diags diag.Diagnostics
	if list.IsNull() || list.IsUnknown() {
		return nil, diags
	}

	var tags []albTagModel
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

// expandAlbXForwardedForConfig converts the planned x_forwarded_for_config list
// into the SDK object. The API accepts a single object, so the first element of
// the list is used.
func expandAlbXForwardedForConfig(ctx context.Context, list types.List) (*albv20251030.XForwardedForConfig, diag.Diagnostics) {
	var diags diag.Diagnostics
	if list.IsNull() || list.IsUnknown() {
		return nil, diags
	}

	var configs []albXForwardedForConfigModel
	diags.Append(list.ElementsAs(ctx, &configs, false)...)
	if diags.HasError() {
		return nil, diags
	}
	if len(configs) == 0 {
		return nil, diags
	}

	config := configs[0]
	xForwardedForConfig := &albv20251030.XForwardedForConfig{
		XForwardedForAlbIdEnabled:         fwhelper.BoolPointerFromValue(config.XForwardedForAlbIdEnabled),
		XForwardedForClientSrcPortEnabled: fwhelper.BoolPointerFromValue(config.XForwardedForClientSrcPortEnabled),
		XForwardedForHostEnabled:          fwhelper.BoolPointerFromValue(config.XForwardedForHostEnabled),
		XForwardedForMode:                 fwhelper.StringPointerFromValue(config.XForwardedForMode),
		XForwardedForPortEnabled:          fwhelper.BoolPointerFromValue(config.XForwardedForPortEnabled),
		XForwardedForProtoEnabled:         fwhelper.BoolPointerFromValue(config.XForwardedForProtoEnabled),
		XTencentClientIDNEnabled:          fwhelper.BoolPointerFromValue(config.XTencentClientIdnEnabled),
		XTencentClientSDNEnabled:          fwhelper.BoolPointerFromValue(config.XTencentClientSdnEnabled),
		XTencentClientSerialEnabled:       fwhelper.BoolPointerFromValue(config.XTencentClientSerialEnabled),
		XTencentClientVerifyEnabled:       fwhelper.BoolPointerFromValue(config.XTencentClientVerifyEnabled),
	}

	return xForwardedForConfig, diags
}

// expandStringSet converts a set attribute into the SDK pointer slice form. A
// null or unknown set is transmitted as "absent".
func expandStringSet(ctx context.Context, value types.Set) ([]*string, diag.Diagnostics) {
	var diags diag.Diagnostics
	if value.IsNull() || value.IsUnknown() {
		return nil, diags
	}

	var elements []types.String
	diags.Append(value.ElementsAs(ctx, &elements, false)...)
	if diags.HasError() {
		return nil, diags
	}

	out := make([]*string, 0, len(elements))
	for _, element := range elements {
		if element.IsNull() || element.IsUnknown() {
			continue
		}
		out = append(out, element.ValueStringPointer())
	}

	if len(out) == 0 {
		return nil, diags
	}

	return out, diags
}

// flattenAlbTags builds the tags list from the describe response. Tags attached
// by the cloud itself (the `qcs:` prefixed CAM cost allocation tags) are not
// managed by this resource and are skipped, otherwise Terraform reports a
// permanent diff.
func flattenAlbTags(ctx context.Context, items []*albv20251030.TagInfo, prevList types.List) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics

	models := make([]albTagModel, 0, len(items))
	for _, item := range items {
		if item == nil {
			continue
		}
		if item.TagKey != nil && strings.HasPrefix(*item.TagKey, "qcs:") {
			continue
		}
		models = append(models, albTagModel{
			TagKey:   fwhelper.StringPointerValueOrNull(item.TagKey),
			TagValue: fwhelper.StringPointerValueOrNull(item.TagValue),
		})
	}

	// Tags are only accepted when the listener is created: an empty response
	// must not clear the tags that were planned.
	if len(models) == 0 && !prevList.IsNull() && !prevList.IsUnknown() {
		return prevList, diags
	}
	if len(models) == 0 {
		return types.ListNull(albTagObjectType), diags
	}

	list, listDiags := types.ListValueFrom(ctx, albTagObjectType, models)
	diags.Append(listDiags...)
	return list, diags
}

// flattenAlbXForwardedForConfig builds the x_forwarded_for_config list from the
// describe response. The previous list is used to keep configured values the
// API omits.
func flattenAlbXForwardedForConfig(ctx context.Context, item *albv20251030.XForwardedForConfig, prevList types.List) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics

	if item == nil {
		return types.ListNull(albXForwardedForConfigObjectType), diags
	}

	var prev []albXForwardedForConfigModel
	if !prevList.IsNull() && !prevList.IsUnknown() {
		diags.Append(prevList.ElementsAs(ctx, &prev, false)...)
		if diags.HasError() {
			return types.ListNull(albXForwardedForConfigObjectType), diags
		}
	}

	var prior albXForwardedForConfigModel
	if len(prev) > 0 {
		prior = prev[0]
	}

	model := albXForwardedForConfigModel{
		XForwardedForAlbIdEnabled:         boolValueOrPrev(item.XForwardedForAlbIdEnabled, prior.XForwardedForAlbIdEnabled),
		XForwardedForClientSrcPortEnabled: boolValueOrPrev(item.XForwardedForClientSrcPortEnabled, prior.XForwardedForClientSrcPortEnabled),
		XForwardedForHostEnabled:          boolValueOrPrev(item.XForwardedForHostEnabled, prior.XForwardedForHostEnabled),
		XForwardedForMode:                 fwhelper.StringValueOrPrev(item.XForwardedForMode, prior.XForwardedForMode),
		XForwardedForPortEnabled:          boolValueOrPrev(item.XForwardedForPortEnabled, prior.XForwardedForPortEnabled),
		XForwardedForProtoEnabled:         boolValueOrPrev(item.XForwardedForProtoEnabled, prior.XForwardedForProtoEnabled),
		XTencentClientIdnEnabled:          boolValueOrPrev(item.XTencentClientIDNEnabled, prior.XTencentClientIdnEnabled),
		XTencentClientSdnEnabled:          boolValueOrPrev(item.XTencentClientSDNEnabled, prior.XTencentClientSdnEnabled),
		XTencentClientSerialEnabled:       boolValueOrPrev(item.XTencentClientSerialEnabled, prior.XTencentClientSerialEnabled),
		XTencentClientVerifyEnabled:       boolValueOrPrev(item.XTencentClientVerifyEnabled, prior.XTencentClientVerifyEnabled),
	}

	list, listDiags := types.ListValueFrom(ctx, albXForwardedForConfigObjectType, []albXForwardedForConfigModel{model})
	diags.Append(listDiags...)
	return list, diags
}

// int64ValueOrPrev returns the value reported by the cloud and, when the API
// omitted it, falls back to the previous plan/state value. It is the *int64
// counterpart of fwhelper.Int64ValueOrPrev, which only handles *uint64.
func int64ValueOrPrev(cloud *int64, prev types.Int64) types.Int64 {
	switch {
	case cloud != nil:
		return types.Int64Value(*cloud)
	case !prev.IsNull() && !prev.IsUnknown():
		return prev
	default:
		return types.Int64Null()
	}
}

// boolValueOrPrev returns the value reported by the cloud and, when the API
// omitted it, falls back to the previous plan/state value. It is the bool
// counterpart of fwhelper.StringValueOrPrev.
func boolValueOrPrev(cloud *bool, prev types.Bool) types.Bool {
	switch {
	case cloud != nil:
		return types.BoolValue(*cloud)
	case !prev.IsNull() && !prev.IsUnknown():
		return prev
	default:
		return types.BoolNull()
	}
}
