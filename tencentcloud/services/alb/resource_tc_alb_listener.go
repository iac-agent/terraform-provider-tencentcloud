package alb

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	albv20251030 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/alb/v20251030"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/fw"
	fwHelper "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/helper"
)

// Compile-time interface assertions: only the extension interfaces are
// asserted, the base resource.Resource interface is checked implicitly by the
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
// The cloud side addresses a listener by the pair (load balancer id, listener
// id): DescribeListeners, ModifyListenerAttributes and DeleteListener all
// require the load balancer id. The Terraform id is therefore the composite
// `<load_balancer_id>#<listener_id>` (parent first).
type AlbListenerResource struct {
	fw.ResourceWithConfigure
}

// AlbListenerDefaultActionModel maps one element of `default_actions`.
type AlbListenerDefaultActionModel struct {
	TargetGroupConfig types.Object `tfsdk:"target_group_config"`
	Type              types.String `tfsdk:"type"`
}

// AlbListenerTargetGroupConfigModel maps the `target_group_config` object.
type AlbListenerTargetGroupConfigModel struct {
	TargetGroups             types.List   `tfsdk:"target_groups"`
	TargetGroupStickySession types.Object `tfsdk:"target_group_sticky_session"`
}

// AlbListenerTargetGroupModel maps one element of `target_groups`.
type AlbListenerTargetGroupModel struct {
	TargetGroupId types.String `tfsdk:"target_group_id"`
	Weight        types.Int64  `tfsdk:"weight"`
}

// AlbListenerStickySessionModel maps the `target_group_sticky_session` object.
type AlbListenerStickySessionModel struct {
	Enabled types.Bool  `tfsdk:"enabled"`
	Timeout types.Int64 `tfsdk:"timeout"`
}

// AlbListenerTagModel maps one element of `tags`.
type AlbListenerTagModel struct {
	TagKey   types.String `tfsdk:"tag_key"`
	TagValue types.String `tfsdk:"tag_value"`
}

// AlbListenerXForwardedForConfigModel maps the `x_forwarded_for_config`
// object.
type AlbListenerXForwardedForConfigModel struct {
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

// AlbListenerResourceModel maps the schema attributes.
type AlbListenerResourceModel struct {
	Id                  types.String `tfsdk:"id"`
	ListenerId          types.String `tfsdk:"listener_id"`
	LoadBalancerId      types.String `tfsdk:"load_balancer_id"`
	ListenerPort        types.Int64  `tfsdk:"listener_port"`
	ListenerProtocol    types.String `tfsdk:"listener_protocol"`
	ListenerName        types.String `tfsdk:"listener_name"`
	DefaultActions      types.List   `tfsdk:"default_actions"`
	CaCertificateIds    types.List   `tfsdk:"ca_certificate_ids"`
	CaEnabled           types.Bool   `tfsdk:"ca_enabled"`
	CertificateIds      types.List   `tfsdk:"certificate_ids"`
	ClientToken         types.String `tfsdk:"client_token"`
	GzipEnabled         types.Bool   `tfsdk:"gzip_enabled"`
	Http2Enabled        types.Bool   `tfsdk:"http2_enabled"`
	IdleTimeout         types.Int64  `tfsdk:"idle_timeout"`
	RequestTimeout      types.Int64  `tfsdk:"request_timeout"`
	SecurityPolicyId    types.String `tfsdk:"security_policy_id"`
	Tags                types.List   `tfsdk:"tags"`
	XForwardedForConfig types.Object `tfsdk:"x_forwarded_for_config"`
}

// Attribute type maps mirroring the nested schema definitions, used by the
// expand/flatten helpers.
var (
	albListenerTagAttrTypes = map[string]attr.Type{
		"tag_key":   types.StringType,
		"tag_value": types.StringType,
	}
	albListenerXForwardedForConfigAttrTypes = map[string]attr.Type{
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
	}

	albListenerTagObjectType = types.ObjectType{AttrTypes: albListenerTagAttrTypes}
)

func (r *AlbListenerResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "tencentcloud_alb_listener"
}

func (r *AlbListenerResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a resource to create an ALB listener.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "ID of the resource, in the form `<load_balancer_id>#<listener_id>`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"listener_id": schema.StringAttribute{
				Computed:    true,
				Description: "Listener ID, in the form `lst-` followed by 8 letters or digits.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"load_balancer_id": schema.StringAttribute{
				Required:    true,
				Description: "Load balancer instance ID, in the form `alb-` followed by 8 letters or digits.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"listener_port": schema.Int64Attribute{
				Required:    true,
				Description: "Port used by the frontend of the load balancer instance. Value range: 1~65535.",
				PlanModifiers: []planmodifier.Int64{
					// The modify API does not accept the port, so a port change
					// is only reachable through replacement.
					int64planmodifier.RequiresReplace(),
				},
			},
			"listener_protocol": schema.StringAttribute{
				Required:    true,
				Description: "Listener protocol. Valid values: `HTTP`, `HTTPS`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"listener_name": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Custom listener name. Length: 1~255 characters, must be a Chinese or harmless string, may contain Chinese characters, letters, digits, hyphens (-), forward slashes (/), periods (.) and underscores (_).",
			},
			"default_actions": schema.ListNestedAttribute{
				Required:    true,
				Description: "List of default forwarding rule actions. Currently a listener supports adding only 1 default forwarding rule action.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"target_group_config": schema.SingleNestedAttribute{
							Required:    true,
							Description: "Forwarding target group configuration. When creating a listener, the target group configuration of a forwarding action supports only a single target group.",
							Attributes: map[string]schema.Attribute{
								"target_groups": schema.ListNestedAttribute{
									Required:    true,
									Description: "Target group list.",
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{
											"target_group_id": schema.StringAttribute{
												Required:    true,
												Description: "Target group ID, in the form `lbtg-` followed by 8 letters or digits.",
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
											Description: "Timeout in seconds. Value range: 1~86400. Default value: 1000.",
										},
									},
								},
							},
						},
						"type": schema.StringAttribute{
							Required:    true,
							Description: "Forwarding action type. When creating a listener, the default forwarding action type supports only forwarding to a target group.",
						},
					},
				},
			},
			"ca_certificate_ids": schema.ListAttribute{
				Optional:    true,
				ElementType: types.StringType,
				Description: "List of CA certificate IDs configured for the listener. Currently a listener supports adding only 1 CA certificate. Required when `ca_enabled` is `true`.",
			},
			"ca_enabled": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Whether to enable mutual authentication. Valid values: `true` (enabled), `false` (not enabled, default).",
			},
			"certificate_ids": schema.ListAttribute{
				Optional:    true,
				ElementType: types.StringType,
				Description: "List of server certificate IDs.",
			},
			"client_token": schema.StringAttribute{
				Optional:    true,
				Description: "Client token used to ensure idempotency of the request. Generate a value on your client side and make sure it is unique among different requests. Only ASCII characters are supported.",
			},
			"gzip_enabled": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Whether to enable Gzip compression. Valid values: `true` (enabled, default), `false` (not enabled).",
			},
			"http2_enabled": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Whether to enable the HTTP/2 feature. Default value: `false` for the HTTP protocol and `true` for the HTTPS protocol. Only the HTTPS protocol supports this parameter.",
			},
			"idle_timeout": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "Connection idle timeout in seconds. Value range: 1~600. Default value: 15. If no request is received within the timeout, the load balancer disconnects the current connection and creates a new one when the next request arrives.",
			},
			"request_timeout": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "Connection request timeout in seconds. Value range: 1~600. Default value: 60. If the backend server does not return a response within the timeout, the load balancer gives up waiting and returns HTTP 504 to the client.",
			},
			"security_policy_id": schema.StringAttribute{
				Optional:    true,
				Description: "Security policy ID, in the form `tls-` followed by 8 letters or digits.",
			},
			"tags": schema.ListNestedAttribute{
				Optional:    true,
				Description: "Tag list. Maximum: 20 tags.",
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
					// The modify API does not accept tags, so a tag change has
					// no in-place path.
					listplanmodifier.RequiresReplace(),
				},
			},
			"x_forwarded_for_config": schema.SingleNestedAttribute{
				Optional:    true,
				Description: "X-Forwarded-For configuration.",
				Attributes: map[string]schema.Attribute{
					"x_forwarded_for_alb_id_enabled": schema.BoolAttribute{
						Optional:    true,
						Description: "Whether to obtain the load balancer instance ID through the ALB-ID header. `true`: yes. `false`: no.",
					},
					"x_forwarded_for_client_src_port_enabled": schema.BoolAttribute{
						Optional:    true,
						Description: "Whether to obtain the port of the client that accesses the load balancer instance through the X-Forwarded-Client-srcport header. `true`: yes. `false`: no.",
					},
					"x_forwarded_for_host_enabled": schema.BoolAttribute{
						Optional:    true,
						Description: "Whether to enable obtaining the domain name of the client that accesses the load balancer instance through the X-Forwarded-Host header. `true`: yes. `false`: no.",
					},
					"x_forwarded_for_mode": schema.StringAttribute{
						Optional:    true,
						Description: "How to handle the X-Forwarded-For (XFF) HTTP header. `append`: append mode (default), appends the real IP address of the client to the end of the X-Forwarded-For header and keeps the original XFF chain. `remove`: remove mode, removes the X-Forwarded-For header and does not pass it to the backend server. `passthrough`: passthrough mode, keeps the X-Forwarded-For header unchanged and passes it to the backend server without modification.",
					},
					"x_forwarded_for_port_enabled": schema.BoolAttribute{
						Optional:    true,
						Description: "Whether to obtain the listening port of the load balancer instance through the X-Forwarded-Port header. `true`: yes. `false`: no.",
					},
					"x_forwarded_for_proto_enabled": schema.BoolAttribute{
						Optional:    true,
						Description: "Whether to obtain the listening protocol of the load balancer instance through the X-Forwarded-Proto header. `true`: yes. `false`: no.",
					},
					"x_tencent_client_idn_enabled": schema.BoolAttribute{
						Optional:    true,
						Description: "Whether to access the issuer of the client certificate `$ssl_client_i_dn` through the X-Tencent-Client-IDN header. `true`: yes. `false`: no.",
					},
					"x_tencent_client_sdn_enabled": schema.BoolAttribute{
						Optional:    true,
						Description: "Whether to access the subject `$ssl_client_s_dn` of the client certificate through the X-Tencent-Client-SDN header. `true`: yes. `false`: no.",
					},
					"x_tencent_client_serial_enabled": schema.BoolAttribute{
						Optional:    true,
						Description: "Whether to access the serial number `$ssl_client_serial` of the client certificate through the X-Tencent-Client-Serial header. `true`: yes. `false`: no.",
					},
					"x_tencent_client_verify_enabled": schema.BoolAttribute{
						Optional:    true,
						Description: "Whether to access the verification result `$ssl_client_verify` of the client certificate through the X-Tencent-Client-Verify header. `true`: yes. `false`: no.",
					},
				},
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

	createListenerRequest := albv20251030.NewCreateListenerRequest()
	createListenerRequest.LoadBalancerId = fwHelper.StringPointerFromValue(plan.LoadBalancerId)
	createListenerRequest.ListenerPort = fwHelper.Uint64PointerFromValue(plan.ListenerPort)
	createListenerRequest.ListenerProtocol = fwHelper.StringPointerFromValue(plan.ListenerProtocol)
	createListenerRequest.ListenerName = fwHelper.StringPointerFromValue(plan.ListenerName)
	createListenerRequest.CaEnabled = fwHelper.BoolPointerFromValue(plan.CaEnabled)
	createListenerRequest.ClientToken = fwHelper.StringPointerFromValue(plan.ClientToken)
	createListenerRequest.GzipEnabled = fwHelper.BoolPointerFromValue(plan.GzipEnabled)
	createListenerRequest.Http2Enabled = fwHelper.BoolPointerFromValue(plan.Http2Enabled)
	createListenerRequest.IdleTimeout = fwHelper.Uint64PointerFromValue(plan.IdleTimeout)
	createListenerRequest.RequestTimeout = fwHelper.Uint64PointerFromValue(plan.RequestTimeout)
	createListenerRequest.SecurityPolicyId = fwHelper.StringPointerFromValue(plan.SecurityPolicyId)

	defaultActions, actionDiags := expandAlbListenerDefaultActions(ctx, plan.DefaultActions)
	resp.Diagnostics.Append(actionDiags...)
	if resp.Diagnostics.HasError() {
		return
	}
	createListenerRequest.DefaultActions = defaultActions

	caCertificateIds := fwHelper.PointersFromList[string](ctx, plan.CaCertificateIds, &resp.Diagnostics)
	certificateIds := fwHelper.PointersFromList[string](ctx, plan.CertificateIds, &resp.Diagnostics)
	tags, tagDiags := expandAlbListenerTags(ctx, plan.Tags)
	resp.Diagnostics.Append(tagDiags...)
	xForwardedForConfig, xffDiags := expandAlbListenerXForwardedForConfig(ctx, plan.XForwardedForConfig)
	resp.Diagnostics.Append(xffDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createListenerRequest.CaCertificateIds = caCertificateIds
	createListenerRequest.CertificateIds = certificateIds
	createListenerRequest.Tags = tags
	createListenerRequest.XForwardedForConfig = xForwardedForConfig

	service := NewAlbService(r.Client())

	listenerId, err := service.CreateListener(ctx, createListenerRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating ALB listener",
			fmt.Sprintf("CreateListener failed: %v", err),
		)

		return
	}

	// Never write an empty id into state.
	if listenerId == "" {
		resp.Diagnostics.AddError(
			"Error creating ALB listener",
			"CreateListener returned an empty listener id.",
		)

		return
	}

	loadBalancerId := plan.LoadBalancerId.ValueString()
	plan.ListenerId = types.StringValue(listenerId)
	plan.Id = types.StringValue(fwHelper.IdFormat(loadBalancerId, listenerId))

	// Computed attributes are unknown at this point and a state cannot hold
	// unknown values, so they are normalized to null before the defensive
	// persist below.
	normalizeAlbListenerComputedUnknowns(&plan)

	// Persist the identity before reading back so a partial failure still
	// leaves a destroyable resource in state.
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	describeListenersRequest := albv20251030.NewDescribeListenersRequest()
	describeListenersRequest.LoadBalancerId = fwHelper.String(loadBalancerId)
	describeListenersRequest.ListenerIds = []*string{fwHelper.String(listenerId)}

	listeners, err := service.DescribeListeners(ctx, describeListenersRequest)
	if err != nil {
		resp.Diagnostics.AddError("Error reading ALB listener after create", err.Error())
		return
	}

	item := findAlbListener(listeners, listenerId)
	if item == nil {
		resp.Diagnostics.AddError(
			"ALB listener not readable after create",
			fmt.Sprintf("listener %s was created but cannot be read back", plan.Id.ValueString()),
		)

		return
	}

	resp.Diagnostics.Append(refreshAlbListenerModel(ctx, &plan, item)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Read handles `terraform refresh` and post-apply state refresh.
func (r *AlbListenerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state AlbListenerResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	parts, err := fwHelper.IdParseErr(state.Id.ValueString(), 2)
	if err != nil {
		resp.Diagnostics.AddError("Invalid resource id", err.Error())
		return
	}
	loadBalancerId, listenerId := parts[0], parts[1]

	describeListenersRequest := albv20251030.NewDescribeListenersRequest()
	describeListenersRequest.LoadBalancerId = fwHelper.String(loadBalancerId)
	describeListenersRequest.ListenerIds = []*string{fwHelper.String(listenerId)}

	service := NewAlbService(r.Client())

	listeners, err := service.DescribeListeners(ctx, describeListenersRequest)
	if err != nil {
		resp.Diagnostics.AddError("Error reading ALB listener", err.Error())
		return
	}

	item := findAlbListener(listeners, listenerId)

	// Confirmed "not found" on the cloud side => drop it from state so
	// Terraform plans a recreation.
	if item == nil {
		resp.Diagnostics.AddWarning(
			"Resource not found",
			fmt.Sprintf("ALB listener %s not found, removing from state.", state.Id.ValueString()),
		)

		resp.State.RemoveResource(ctx)
		return
	}

	state.LoadBalancerId = types.StringValue(loadBalancerId)
	state.ListenerId = types.StringValue(listenerId)

	resp.Diagnostics.Append(refreshAlbListenerModel(ctx, &state, item)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update handles in-place modifications.
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

	parts, err := fwHelper.IdParseErr(state.Id.ValueString(), 2)
	if err != nil {
		resp.Diagnostics.AddError("Invalid resource id", err.Error())
		return
	}
	loadBalancerId, listenerId := parts[0], parts[1]

	service := NewAlbService(r.Client())

	// The modify API is called only when one of the attributes it covers
	// changed: a plan differing from the state in the id alone must not send a
	// pointless no-change request.
	if albListenerAttributesChanged(plan, state) {
		modifyListenerAttributesRequest := albv20251030.NewModifyListenerAttributesRequest()
		modifyListenerAttributesRequest.LoadBalancerId = fwHelper.String(loadBalancerId)
		modifyListenerAttributesRequest.ListenerId = fwHelper.String(listenerId)
		modifyListenerAttributesRequest.ListenerName = fwHelper.StringPointerFromValue(plan.ListenerName)
		modifyListenerAttributesRequest.CaEnabled = fwHelper.BoolPointerFromValue(plan.CaEnabled)
		modifyListenerAttributesRequest.ClientToken = fwHelper.StringPointerFromValue(plan.ClientToken)
		modifyListenerAttributesRequest.GzipEnabled = fwHelper.BoolPointerFromValue(plan.GzipEnabled)
		modifyListenerAttributesRequest.Http2Enabled = fwHelper.BoolPointerFromValue(plan.Http2Enabled)
		modifyListenerAttributesRequest.IdleTimeout = fwHelper.Uint64PointerFromValue(plan.IdleTimeout)
		modifyListenerAttributesRequest.RequestTimeout = fwHelper.Uint64PointerFromValue(plan.RequestTimeout)
		modifyListenerAttributesRequest.SecurityPolicyId = fwHelper.StringPointerFromValue(plan.SecurityPolicyId)

		defaultActions, actionDiags := expandAlbListenerDefaultActions(ctx, plan.DefaultActions)
		resp.Diagnostics.Append(actionDiags...)
		caCertificateIds := fwHelper.PointersFromList[string](ctx, plan.CaCertificateIds, &resp.Diagnostics)
		certificateIds := fwHelper.PointersFromList[string](ctx, plan.CertificateIds, &resp.Diagnostics)
		xForwardedForConfig, xffDiags := expandAlbListenerXForwardedForConfig(ctx, plan.XForwardedForConfig)
		resp.Diagnostics.Append(xffDiags...)
		if resp.Diagnostics.HasError() {
			return
		}

		modifyListenerAttributesRequest.DefaultActions = defaultActions
		modifyListenerAttributesRequest.CaCertificateIds = caCertificateIds
		modifyListenerAttributesRequest.CertificateIds = certificateIds
		modifyListenerAttributesRequest.XForwardedForConfig = xForwardedForConfig

		// TODO(human): confirm the clear semantics of ModifyListenerAttributes —
		// the API documentation does not state whether an omitted optional
		// field keeps or clears the current value, so removing an attribute
		// from the configuration is currently sent as "omitted".
		if err := service.ModifyListenerAttributes(ctx, modifyListenerAttributesRequest); err != nil {
			resp.Diagnostics.AddError(
				"Error updating ALB listener",
				fmt.Sprintf("ModifyListenerAttributes failed: %v", err),
			)

			return
		}
	}

	describeListenersRequest := albv20251030.NewDescribeListenersRequest()
	describeListenersRequest.LoadBalancerId = fwHelper.String(loadBalancerId)
	describeListenersRequest.ListenerIds = []*string{fwHelper.String(listenerId)}

	listeners, err := service.DescribeListeners(ctx, describeListenersRequest)
	if err != nil {
		resp.Diagnostics.AddError("Error reading ALB listener after update", err.Error())
		return
	}

	item := findAlbListener(listeners, listenerId)
	if item == nil {
		resp.Diagnostics.AddError(
			"ALB listener not readable after update",
			fmt.Sprintf("listener %s was updated but cannot be read back", state.Id.ValueString()),
		)

		return
	}

	plan.Id = state.Id
	plan.ListenerId = state.ListenerId
	plan.LoadBalancerId = state.LoadBalancerId
	resp.Diagnostics.Append(refreshAlbListenerModel(ctx, &plan, item)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete handles resource deletion.
func (r *AlbListenerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state AlbListenerResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	parts, err := fwHelper.IdParseErr(state.Id.ValueString(), 2)
	if err != nil {
		resp.Diagnostics.AddError("Invalid resource id", err.Error())
		return
	}
	loadBalancerId, listenerId := parts[0], parts[1]

	deleteListenerRequest := albv20251030.NewDeleteListenerRequest()
	deleteListenerRequest.LoadBalancerId = fwHelper.String(loadBalancerId)
	deleteListenerRequest.ListenerIds = []*string{fwHelper.String(listenerId)}

	service := NewAlbService(r.Client())

	if err := service.DeleteListener(ctx, deleteListenerRequest); err != nil {
		resp.Diagnostics.AddError(
			"Error deleting ALB listener",
			fmt.Sprintf("DeleteListener failed: %v", err),
		)

		return
	}
}

// ImportState supports `terraform import` with the composite id
// `<load_balancer_id>#<listener_id>`.
func (r *AlbListenerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts, err := fwHelper.IdParseErr(req.ID, 2)
	if err != nil {
		resp.Diagnostics.AddError("Invalid resource id", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("load_balancer_id"), parts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("listener_id"), parts[1])...)
}

// albListenerAttributesChanged reports whether an attribute covered by
// ModifyListenerAttributes differs between the plan and the state.
func albListenerAttributesChanged(plan, state AlbListenerResourceModel) bool {
	return !plan.ListenerName.Equal(state.ListenerName) ||
		!plan.ClientToken.Equal(state.ClientToken) ||
		!plan.CaEnabled.Equal(state.CaEnabled) ||
		!plan.GzipEnabled.Equal(state.GzipEnabled) ||
		!plan.Http2Enabled.Equal(state.Http2Enabled) ||
		!plan.IdleTimeout.Equal(state.IdleTimeout) ||
		!plan.RequestTimeout.Equal(state.RequestTimeout) ||
		!plan.SecurityPolicyId.Equal(state.SecurityPolicyId) ||
		!plan.CaCertificateIds.Equal(state.CaCertificateIds) ||
		!plan.CertificateIds.Equal(state.CertificateIds) ||
		!plan.DefaultActions.Equal(state.DefaultActions) ||
		!plan.XForwardedForConfig.Equal(state.XForwardedForConfig)
}

// refreshAlbListenerModel copies the cloud-reported fields into the model.
//
// User-configurable attributes use the *OrPrev family (a cloud-omitted field
// keeps the previous value so Terraform does not fail the post-apply
// consistency check), while the identity attributes are set explicitly.
//
// Attributes the describe API does not return at all (default_actions,
// ca_certificate_ids, certificate_ids) are left untouched: their configured
// values stay in the model.
func refreshAlbListenerModel(ctx context.Context, m *AlbListenerResourceModel, item *albv20251030.ListenerOutput) diag.Diagnostics {
	if item == nil {
		return nil
	}

	var diags diag.Diagnostics

	// The describe response reports `CaEnable`, `Http2Enable` and
	// `TlsSecurityPolicyId` for the request fields `CaEnabled`, `Http2Enabled`
	// and `SecurityPolicyId`.
	m.CaEnabled = fwHelper.BoolValueOrPrev(item.CaEnable, m.CaEnabled)
	m.GzipEnabled = fwHelper.BoolValueOrPrev(item.GzipEnabled, m.GzipEnabled)
	m.Http2Enabled = fwHelper.BoolValueOrPrev(item.Http2Enable, m.Http2Enabled)
	m.IdleTimeout = fwHelper.Int64ValueFromUintOrPrev(item.IdleTimeout, m.IdleTimeout)
	m.ListenerName = fwHelper.StringValueOrPrev(item.ListenerName, m.ListenerName)
	m.RequestTimeout = fwHelper.Int64ValueFromUintOrPrev(item.RequestTimeout, m.RequestTimeout)
	m.SecurityPolicyId = fwHelper.StringValueOrPrev(item.TlsSecurityPolicyId, m.SecurityPolicyId)

	if item.Tags != nil {
		tags, tagDiags := flattenAlbListenerTags(ctx, item.Tags)
		diags.Append(tagDiags...)
		if !tagDiags.HasError() {
			m.Tags = tags
		}
	}

	if item.XForwardedForConfig != nil {
		xForwardedForConfig, xffDiags := flattenAlbListenerXForwardedForConfig(ctx, item.XForwardedForConfig)
		diags.Append(xffDiags...)
		if !xffDiags.HasError() {
			m.XForwardedForConfig = xForwardedForConfig
		}
	}

	return diags
}

// normalizeAlbListenerComputedUnknowns rewrites unknown computed attributes to
// null so the model can be persisted into state before the read-back fills
// them with the cloud-side values.
func normalizeAlbListenerComputedUnknowns(m *AlbListenerResourceModel) {
	if m.ListenerId.IsUnknown() {
		m.ListenerId = types.StringNull()
	}
	if m.ListenerName.IsUnknown() {
		m.ListenerName = types.StringNull()
	}
	if m.CaEnabled.IsUnknown() {
		m.CaEnabled = types.BoolNull()
	}
	if m.GzipEnabled.IsUnknown() {
		m.GzipEnabled = types.BoolNull()
	}
	if m.Http2Enabled.IsUnknown() {
		m.Http2Enabled = types.BoolNull()
	}
	if m.IdleTimeout.IsUnknown() {
		m.IdleTimeout = types.Int64Null()
	}
	if m.RequestTimeout.IsUnknown() {
		m.RequestTimeout = types.Int64Null()
	}
}

// findAlbListener picks this resource's own item out of a describe list.
func findAlbListener(listeners []*albv20251030.ListenerOutput, listenerId string) *albv20251030.ListenerOutput {
	for _, listener := range listeners {
		if listener != nil && listener.ListenerId != nil && *listener.ListenerId == listenerId {
			return listener
		}
	}

	return nil
}

// expandAlbListenerDefaultActions converts the `default_actions` list into the
// SDK form.
func expandAlbListenerDefaultActions(ctx context.Context, list types.List) ([]*albv20251030.DefaultAction, diag.Diagnostics) {
	if list.IsNull() || list.IsUnknown() {
		return nil, nil
	}

	var models []AlbListenerDefaultActionModel
	diags := list.ElementsAs(ctx, &models, false)
	if diags.HasError() {
		return nil, diags
	}

	actions := make([]*albv20251030.DefaultAction, 0, len(models))
	for _, model := range models {
		action := &albv20251030.DefaultAction{}
		if !model.Type.IsNull() && !model.Type.IsUnknown() {
			action.Type = model.Type.ValueStringPointer()
		}

		targetGroupConfig, configDiags := expandAlbListenerTargetGroupConfig(ctx, model.TargetGroupConfig)
		diags.Append(configDiags...)
		if configDiags.HasError() {
			return nil, diags
		}
		action.TargetGroupConfig = targetGroupConfig

		actions = append(actions, action)
	}

	return actions, diags
}

// expandAlbListenerTargetGroupConfig converts the `target_group_config` object
// into the SDK form.
func expandAlbListenerTargetGroupConfig(ctx context.Context, object types.Object) (*albv20251030.TargetGroupConfig, diag.Diagnostics) {
	if object.IsNull() || object.IsUnknown() {
		return nil, nil
	}

	var model AlbListenerTargetGroupConfigModel
	diags := object.As(ctx, &model, basetypes.ObjectAsOptions{})
	if diags.HasError() {
		return nil, diags
	}

	config := &albv20251030.TargetGroupConfig{}

	targetGroups, targetGroupDiags := expandAlbListenerTargetGroups(ctx, model.TargetGroups)
	diags.Append(targetGroupDiags...)
	if targetGroupDiags.HasError() {
		return nil, diags
	}
	config.TargetGroups = targetGroups

	stickySession, stickySessionDiags := expandAlbListenerStickySession(ctx, model.TargetGroupStickySession)
	diags.Append(stickySessionDiags...)
	if stickySessionDiags.HasError() {
		return nil, diags
	}
	config.TargetGroupStickySession = stickySession

	return config, diags
}

// expandAlbListenerTargetGroups converts the `target_groups` list into the SDK
// form.
func expandAlbListenerTargetGroups(ctx context.Context, list types.List) ([]*albv20251030.TargetGroupTuple, diag.Diagnostics) {
	if list.IsNull() || list.IsUnknown() {
		return nil, nil
	}

	var models []AlbListenerTargetGroupModel
	diags := list.ElementsAs(ctx, &models, false)
	if diags.HasError() {
		return nil, diags
	}

	targetGroups := make([]*albv20251030.TargetGroupTuple, 0, len(models))
	for _, model := range models {
		targetGroup := &albv20251030.TargetGroupTuple{}
		if !model.TargetGroupId.IsNull() && !model.TargetGroupId.IsUnknown() {
			targetGroup.TargetGroupId = model.TargetGroupId.ValueStringPointer()
		}
		if !model.Weight.IsNull() && !model.Weight.IsUnknown() {
			targetGroup.Weight = model.Weight.ValueInt64Pointer()
		}

		targetGroups = append(targetGroups, targetGroup)
	}

	return targetGroups, diags
}

// expandAlbListenerStickySession converts the `target_group_sticky_session`
// object into the SDK form.
func expandAlbListenerStickySession(ctx context.Context, object types.Object) (*albv20251030.TargetGroupStickySession, diag.Diagnostics) {
	if object.IsNull() || object.IsUnknown() {
		return nil, nil
	}

	var model AlbListenerStickySessionModel
	diags := object.As(ctx, &model, basetypes.ObjectAsOptions{})
	if diags.HasError() {
		return nil, diags
	}

	stickySession := &albv20251030.TargetGroupStickySession{}
	if !model.Enabled.IsNull() && !model.Enabled.IsUnknown() {
		stickySession.Enabled = model.Enabled.ValueBoolPointer()
	}
	if !model.Timeout.IsNull() && !model.Timeout.IsUnknown() {
		stickySession.Timeout = model.Timeout.ValueInt64Pointer()
	}

	return stickySession, diags
}

// flattenAlbListenerTags converts the cloud tag list into the `tags` list.
func flattenAlbListenerTags(ctx context.Context, items []*albv20251030.TagInfo) (types.List, diag.Diagnostics) {
	models := make([]AlbListenerTagModel, 0, len(items))
	for _, item := range items {
		if item == nil {
			continue
		}

		models = append(models, AlbListenerTagModel{
			TagKey:   fwHelper.StringPointerValueOrNull(item.TagKey),
			TagValue: fwHelper.StringPointerValueOrNull(item.TagValue),
		})
	}

	if len(models) == 0 {
		return types.ListNull(albListenerTagObjectType), nil
	}

	return types.ListValueFrom(ctx, albListenerTagObjectType, models)
}

// expandAlbListenerTags converts the `tags` list into the SDK form.
func expandAlbListenerTags(ctx context.Context, list types.List) ([]*albv20251030.TagInfo, diag.Diagnostics) {
	if list.IsNull() || list.IsUnknown() {
		return nil, nil
	}

	var models []AlbListenerTagModel
	diags := list.ElementsAs(ctx, &models, false)
	if diags.HasError() {
		return nil, diags
	}

	tags := make([]*albv20251030.TagInfo, 0, len(models))
	for _, model := range models {
		tag := &albv20251030.TagInfo{}
		if !model.TagKey.IsNull() && !model.TagKey.IsUnknown() {
			tag.TagKey = model.TagKey.ValueStringPointer()
		}
		if !model.TagValue.IsNull() && !model.TagValue.IsUnknown() {
			tag.TagValue = model.TagValue.ValueStringPointer()
		}

		tags = append(tags, tag)
	}

	return tags, diags
}

// expandAlbListenerXForwardedForConfig converts the `x_forwarded_for_config`
// object into the SDK form.
func expandAlbListenerXForwardedForConfig(ctx context.Context, object types.Object) (*albv20251030.XForwardedForConfig, diag.Diagnostics) {
	if object.IsNull() || object.IsUnknown() {
		return nil, nil
	}

	var model AlbListenerXForwardedForConfigModel
	diags := object.As(ctx, &model, basetypes.ObjectAsOptions{})
	if diags.HasError() {
		return nil, diags
	}

	return &albv20251030.XForwardedForConfig{
		XForwardedForAlbIdEnabled:         fwHelper.BoolPointerFromValue(model.XForwardedForAlbIdEnabled),
		XForwardedForClientSrcPortEnabled: fwHelper.BoolPointerFromValue(model.XForwardedForClientSrcPortEnabled),
		XForwardedForHostEnabled:          fwHelper.BoolPointerFromValue(model.XForwardedForHostEnabled),
		XForwardedForMode:                 fwHelper.StringPointerFromValue(model.XForwardedForMode),
		XForwardedForPortEnabled:          fwHelper.BoolPointerFromValue(model.XForwardedForPortEnabled),
		XForwardedForProtoEnabled:         fwHelper.BoolPointerFromValue(model.XForwardedForProtoEnabled),
		XTencentClientIDNEnabled:          fwHelper.BoolPointerFromValue(model.XTencentClientIdnEnabled),
		XTencentClientSDNEnabled:          fwHelper.BoolPointerFromValue(model.XTencentClientSdnEnabled),
		XTencentClientSerialEnabled:       fwHelper.BoolPointerFromValue(model.XTencentClientSerialEnabled),
		XTencentClientVerifyEnabled:       fwHelper.BoolPointerFromValue(model.XTencentClientVerifyEnabled),
	}, diags
}

// flattenAlbListenerXForwardedForConfig converts the cloud
// XForwardedForConfig into the `x_forwarded_for_config` object.
func flattenAlbListenerXForwardedForConfig(ctx context.Context, item *albv20251030.XForwardedForConfig) (types.Object, diag.Diagnostics) {
	model := AlbListenerXForwardedForConfigModel{
		XForwardedForAlbIdEnabled:         fwHelper.BoolValueOrNull(item.XForwardedForAlbIdEnabled),
		XForwardedForClientSrcPortEnabled: fwHelper.BoolValueOrNull(item.XForwardedForClientSrcPortEnabled),
		XForwardedForHostEnabled:          fwHelper.BoolValueOrNull(item.XForwardedForHostEnabled),
		XForwardedForMode:                 fwHelper.StringPointerValueOrNull(item.XForwardedForMode),
		XForwardedForPortEnabled:          fwHelper.BoolValueOrNull(item.XForwardedForPortEnabled),
		XForwardedForProtoEnabled:         fwHelper.BoolValueOrNull(item.XForwardedForProtoEnabled),
		XTencentClientIdnEnabled:          fwHelper.BoolValueOrNull(item.XTencentClientIDNEnabled),
		XTencentClientSdnEnabled:          fwHelper.BoolValueOrNull(item.XTencentClientSDNEnabled),
		XTencentClientSerialEnabled:       fwHelper.BoolValueOrNull(item.XTencentClientSerialEnabled),
		XTencentClientVerifyEnabled:       fwHelper.BoolValueOrNull(item.XTencentClientVerifyEnabled),
	}

	return types.ObjectValueFrom(ctx, albListenerXForwardedForConfigAttrTypes, model)
}
