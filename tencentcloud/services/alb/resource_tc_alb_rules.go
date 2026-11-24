package alb

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	albv20251030 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/alb/v20251030"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/fw"
	fwhelper "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	svctag "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tag"
)

// Compile-time interface assertions (repo convention: only assert the
// extension interfaces; the base interface is checked implicitly by the
// factory return value).
var (
	_ resource.ResourceWithConfigure   = &AlbRulesResource{}
	_ resource.ResourceWithImportState = &AlbRulesResource{}
)

// NewAlbRulesResource is the factory referenced by
// tencentcloud/framework/registry.go to register this resource.
func NewAlbRulesResource() resource.Resource {
	return &AlbRulesResource{}
}

// AlbRulesResource implements resource.Resource for tencentcloud_alb_rules.
//
// A forwarding rule is addressed by the load balancer id, the listener id and
// the rule id, so the Terraform id joins all three parts with `#`
// (tccommon.FILED_SP) and splits them again for every CRUD call. One TF
// instance manages exactly one forwarding rule even though the ALB APIs
// exchange lists.
type AlbRulesResource struct {
	fw.ResourceWithConfigure
}

// AlbRulesResourceModel maps the top-level schema attributes. Every schema
// attribute appears here exactly once.
type AlbRulesResourceModel struct {
	Id             types.String `tfsdk:"id"`
	ListenerId     types.String `tfsdk:"listener_id"`
	LoadBalancerId types.String `tfsdk:"load_balancer_id"`
	Rules          types.List   `tfsdk:"rules"`
}

// albRuleModel maps one element of rules.
type albRuleModel struct {
	Actions    types.List   `tfsdk:"actions"`
	Conditions types.List   `tfsdk:"conditions"`
	Priority   types.Int64  `tfsdk:"priority"`
	Direction  types.String `tfsdk:"direction"`
	RuleName   types.String `tfsdk:"rule_name"`
	Tags       types.List   `tfsdk:"tags"`
}

// albRuleActionModel maps one element of rules[*].actions.
type albRuleActionModel struct {
	Order               types.Int64  `tfsdk:"order"`
	Type                types.String `tfsdk:"type"`
	FixedResponseConfig types.List   `tfsdk:"fixed_response_config"`
	InsertHeaderConfig  types.List   `tfsdk:"insert_header_config"`
	RedirectConfig      types.List   `tfsdk:"redirect_config"`
	RemoveHeaderConfig  types.List   `tfsdk:"remove_header_config"`
	RewriteConfig       types.List   `tfsdk:"rewrite_config"`
	TargetGroupConfig   types.List   `tfsdk:"target_group_config"`
}

// albFixedResponseConfigModel maps one element of
// rules[*].actions[*].fixed_response_config.
type albFixedResponseConfigModel struct {
	HttpCode    types.Int64  `tfsdk:"http_code"`
	Content     types.String `tfsdk:"content"`
	ContentType types.String `tfsdk:"content_type"`
}

// albInsertHeaderConfigModel maps one element of
// rules[*].actions[*].insert_header_config.
type albInsertHeaderConfigModel struct {
	Key       types.String `tfsdk:"key"`
	Value     types.String `tfsdk:"value"`
	ValueType types.String `tfsdk:"value_type"`
}

// albRedirectConfigModel maps one element of
// rules[*].actions[*].redirect_config.
type albRedirectConfigModel struct {
	HttpCode types.Int64  `tfsdk:"http_code"`
	Host     types.String `tfsdk:"host"`
	Path     types.String `tfsdk:"path"`
	Port     types.String `tfsdk:"port"`
	Protocol types.String `tfsdk:"protocol"`
	Query    types.String `tfsdk:"query"`
}

// albRemoveHeaderConfigModel maps one element of
// rules[*].actions[*].remove_header_config.
type albRemoveHeaderConfigModel struct {
	Key types.String `tfsdk:"key"`
}

// albRewriteConfigModel maps one element of rules[*].actions[*].rewrite_config.
type albRewriteConfigModel struct {
	Host  types.String `tfsdk:"host"`
	Path  types.String `tfsdk:"path"`
	Query types.String `tfsdk:"query"`
}

// albTargetGroupConfigModel maps one element of
// rules[*].actions[*].target_group_config.
type albTargetGroupConfigModel struct {
	TargetGroups             types.List `tfsdk:"target_groups"`
	TargetGroupStickySession types.List `tfsdk:"target_group_sticky_session"`
}

// albTargetGroupTupleModel maps one element of
// rules[*].actions[*].target_group_config[*].target_groups.
type albTargetGroupTupleModel struct {
	TargetGroupId types.String `tfsdk:"target_group_id"`
	Weight        types.Int64  `tfsdk:"weight"`
}

// albTargetGroupStickySessionModel maps one element of
// rules[*].actions[*].target_group_config[*].target_group_sticky_session.
type albTargetGroupStickySessionModel struct {
	Enabled types.Bool  `tfsdk:"enabled"`
	Timeout types.Int64 `tfsdk:"timeout"`
}

// albRuleConditionModel maps one element of rules[*].conditions.
type albRuleConditionModel struct {
	Type              types.String `tfsdk:"type"`
	CookieConfig      types.List   `tfsdk:"cookie_config"`
	HeaderConfig      types.List   `tfsdk:"header_config"`
	HostConfig        types.List   `tfsdk:"host_config"`
	MethodConfig      types.List   `tfsdk:"method_config"`
	PathConfig        types.List   `tfsdk:"path_config"`
	QueryStringConfig types.List   `tfsdk:"query_string_config"`
	SourceIpConfig    types.List   `tfsdk:"source_ip_config"`
}

// albCookieConfigModel maps one element of
// rules[*].conditions[*].cookie_config.
type albCookieConfigModel struct {
	Key   types.String `tfsdk:"key"`
	Value types.String `tfsdk:"value"`
}

// albHeaderConfigModel maps one element of
// rules[*].conditions[*].header_config.
type albHeaderConfigModel struct {
	Key    types.String `tfsdk:"key"`
	Values types.List   `tfsdk:"values"`
}

// albQueryStringConfigModel maps one element of
// rules[*].conditions[*].query_string_config.
type albQueryStringConfigModel struct {
	Key   types.String `tfsdk:"key"`
	Value types.String `tfsdk:"value"`
}

// albTagModel maps one element of rules[*].tags.
type albTagModel struct {
	TagKey   types.String `tfsdk:"tag_key"`
	TagValue types.String `tfsdk:"tag_value"`
}

func (r *AlbRulesResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "tencentcloud_alb_rules"
}

func (r *AlbRulesResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a resource to create a forwarding rule of an ALB listener.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Resource identifier in the form `<load_balancer_id>#<listener_id>#<rule_id>`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"listener_id": schema.StringAttribute{
				Required: true,
				Description: "ID of the listener, in the format `lst-` followed by 8 alphanumeric characters." +
					" Immutable: changing it forces resource replacement.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"load_balancer_id": schema.StringAttribute{
				Required: true,
				Description: "ID of the load balancer instance, in the format `alb-` followed by 8 alphanumeric characters." +
					" Immutable: changing it forces resource replacement.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"rules": schema.ListNestedAttribute{
				Required:    true,
				Description: "List of forwarding rules. Exactly one rule is managed by one resource instance.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"actions": schema.ListNestedAttribute{
							Required:    true,
							Description: "List of forwarding rule actions. At least one of `TargetGroup`, `Redirect` or `FixedResponse` type must be present and executed last.",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"order": schema.Int64Attribute{
										Required:    true,
										Description: "Execution order of the forwarding action, must be unique and executed in ascending order. Range: 1 ~ 50000.",
									},
									"type": schema.StringAttribute{
										Required: true,
										Description: "Type of the forwarding action. Valid values: `TargetGroup`, `Redirect`, `FixedResponse`, `Rewrite`, `InsertHeader`, `RemoveHeader`." +
											" The action list must contain one of `TargetGroup`, `Redirect`, `FixedResponse` and it must be executed last.",
									},
									"fixed_response_config": schema.ListNestedAttribute{
										Optional:    true,
										Description: "Configuration of the fixed response content. Required when `type` is `FixedResponse`.",
										NestedObject: schema.NestedAttributeObject{
											Attributes: map[string]schema.Attribute{
												"http_code": schema.Int64Attribute{
													Required:    true,
													Description: "HTTP response code to return. Supports 2xx, 4xx and 5xx codes.",
												},
												"content": schema.StringAttribute{
													Optional:    true,
													Description: "Fixed content to return. Only ASCII characters are supported, up to 1KB.",
												},
												"content_type": schema.StringAttribute{
													Optional:    true,
													Description: "Format of the fixed content. Valid values: `text/plain`, `text/css`, `text/html`, `application/javascript`, `application/json`.",
												},
											},
										},
									},
									"insert_header_config": schema.ListNestedAttribute{
										Optional:    true,
										Description: "Configuration of the HTTP header to insert. Required when `type` is `InsertHeader`.",
										NestedObject: schema.NestedAttributeObject{
											Attributes: map[string]schema.Attribute{
												"key": schema.StringAttribute{
													Required:    true,
													Description: "Key of the HTTP header to insert, 1 ~ 40 characters, charset `a-z A-Z 0-9 - _`. Chinese characters and the headers `Cookie`, `Host`, `Content-Length`, `Connection`, `Upgrade`, `transfer-encoding`, `keep-alive`, `te`, `authority`, `x-forwarded-for`, `x-forwarded-proto`, `x-forwarded-host`, `x-forwarded-port` are not supported.",
												},
												"value": schema.StringAttribute{
													Required:    true,
													Description: "Value of the HTTP header. When `value_type` is `SystemDefined`, valid values: `ClientPort`, `ClientIp`, `Protocol`, `CLBPort`. Otherwise 1 ~ 128 printable characters, `\"` is not supported, leading/trailing spaces and a trailing `\\` are not allowed.",
												},
												"value_type": schema.StringAttribute{
													Required:    true,
													Description: "Type of the HTTP header value. Valid values: `SystemDefined`, `UserDefined`, `ReferenceHeader`.",
												},
											},
										},
									},
									"redirect_config": schema.ListNestedAttribute{
										Optional:    true,
										Description: "Redirect configuration. Required when `type` is `Redirect`. Except `http_code`, not all other fields may keep their default values.",
										NestedObject: schema.NestedAttributeObject{
											Attributes: map[string]schema.Attribute{
												"http_code": schema.Int64Attribute{
													Required:    true,
													Description: "HTTP code of the redirect. Valid values: `301`, `302`, `303`, `307`, `308`.",
												},
												"host": schema.StringAttribute{
													Optional:    true,
													Description: "Host of the redirect, defaults to `${host}`. 3 ~ 128 characters, charset `a-z 0-9 _ . -`.",
												},
												"path": schema.StringAttribute{
													Optional:    true,
													Description: "Path of the redirect, defaults to `${path}`. 1 ~ 128 characters, charset `a-z A-Z 0-9 ? = _ . - / :`.",
												},
												"port": schema.StringAttribute{
													Optional:    true,
													Description: "Port of the redirect, defaults to `${port}`. Range: 1 ~ 65535.",
												},
												"protocol": schema.StringAttribute{
													Optional:    true,
													Description: "Protocol of the redirect, defaults to `${protocol}`. Valid values: `HTTP`, `HTTPS`.",
												},
												"query": schema.StringAttribute{
													Optional:    true,
													Description: "Query string of the redirect, defaults to `${query}`. 1 ~ 128 characters, printable characters except `#[]{}|<>` and spaces.",
												},
											},
										},
									},
									"remove_header_config": schema.ListNestedAttribute{
										Optional:    true,
										Description: "Configuration of the HTTP header to remove. Required when `type` is `RemoveHeader`.",
										NestedObject: schema.NestedAttributeObject{
											Attributes: map[string]schema.Attribute{
												"key": schema.StringAttribute{
													Required:    true,
													Description: "Key of the HTTP header to remove, 1 ~ 40 characters, charset `a-z A-Z 0-9 - _`. The headers `Cookie`, `Host`, `Content-Length`, `Connection`, `Upgrade`, `transfer-encoding`, `keep-alive`, `te`, `authority`, `x-forwarded-for`, `x-forwarded-proto`, `x-forwarded-host`, `x-forwarded-port`, `server` are not supported.",
												},
											},
										},
									},
									"rewrite_config": schema.ListNestedAttribute{
										Optional:    true,
										Description: "Rewrite configuration. Required when `type` is `Rewrite`.",
										NestedObject: schema.NestedAttributeObject{
											Attributes: map[string]schema.Attribute{
												"host": schema.StringAttribute{
													Optional:    true,
													Description: "Host of the rewrite, defaults to `${host}`. 3 ~ 128 characters, charset `a-z 0-9 _ . -`.",
												},
												"path": schema.StringAttribute{
													Optional:    true,
													Description: "Path of the rewrite, defaults to `${path}`. 1 ~ 128 characters, charset `a-z A-Z 0-9 ? = _ . - / :`.",
												},
												"query": schema.StringAttribute{
													Optional:    true,
													Description: "Query string of the rewrite, defaults to `${query}`. 1 ~ 128 characters, printable characters except `#[]{}|<>` and spaces.",
												},
											},
										},
									},
									"target_group_config": schema.ListNestedAttribute{
										Optional:    true,
										Description: "Forwarding target group configuration. Required when `type` is `TargetGroup`.",
										NestedObject: schema.NestedAttributeObject{
											Attributes: map[string]schema.Attribute{
												"target_groups": schema.ListNestedAttribute{
													Required:    true,
													Description: "List of target groups.",
													NestedObject: schema.NestedAttributeObject{
														Attributes: map[string]schema.Attribute{
															"target_group_id": schema.StringAttribute{
																Required:    true,
																Description: "ID of the target group, in the format `lbtg-` followed by 8 alphanumeric characters.",
															},
															"weight": schema.Int64Attribute{
																Optional:    true,
																Description: "Weight of the target group. Range: [0, 100], defaults to 10.",
															},
														},
													},
												},
												"target_group_sticky_session": schema.ListNestedAttribute{
													Optional:    true,
													Description: "Sticky session between target groups.",
													NestedObject: schema.NestedAttributeObject{
														Attributes: map[string]schema.Attribute{
															"enabled": schema.BoolAttribute{
																Required:    true,
																Description: "Whether to enable sticky session, disabled by default.",
															},
															"timeout": schema.Int64Attribute{
																Optional:    true,
																Description: "Sticky session timeout in seconds. Range: 1 ~ 86400, defaults to 1000.",
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
						"conditions": schema.ListNestedAttribute{
							Required:    true,
							Description: "List of forwarding rule conditions.",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"type": schema.StringAttribute{
										Required:    true,
										Description: "Type of the forwarding condition. Valid values: `Host`, `Path`, `Header`, `QueryString`, `Method`, `Cookie`, `SourceIp`.",
									},
									"cookie_config": schema.ListNestedAttribute{
										Optional:    true,
										Description: "Cookie configuration. Required when `type` is `Cookie`.",
										NestedObject: schema.NestedAttributeObject{
											Attributes: map[string]schema.Attribute{
												"key": schema.StringAttribute{
													Required:    true,
													Description: "Key of the cookie, 1 ~ 64 characters, letters, digits and underscores are supported.",
												},
												"value": schema.StringAttribute{
													Required:    true,
													Description: "Value of the cookie, 1 ~ 128 characters, printable characters are supported.",
												},
											},
										},
									},
									"header_config": schema.ListNestedAttribute{
										Optional:    true,
										Description: "HTTP header configuration. Required when `type` is `Header`.",
										NestedObject: schema.NestedAttributeObject{
											Attributes: map[string]schema.Attribute{
												"key": schema.StringAttribute{
													Required:    true,
													Description: "Key of the HTTP header, 1 ~ 40 characters, charset `a-z A-Z 0-9 - _`. Chinese characters and the headers `Host`, `Cookie` are not supported.",
												},
												"values": schema.ListAttribute{
													Required:    true,
													ElementType: types.StringType,
													Description: "Values of the HTTP header, each 1 ~ 128 characters with printable characters. `\"` is not supported, leading/trailing spaces and a trailing `\\` are not allowed.",
												},
											},
										},
									},
									"host_config": schema.ListAttribute{
										Optional:    true,
										ElementType: types.StringType,
										Description: "Host names. May appear only once in one rule, 3 ~ 128 characters, exact, regex and wildcard matching are supported. Must not start or end with `.` or `_`. Exact matching charset: `a-z 0-9 . - _`. Regex matching starts with `~`, charset: `a-z 0-9 . - ? = ~ _ - + \\ ^ * ! $ & | ( ) [ ]`. Wildcard matching: `*` for multi-character and `?` for single-character wildcards, charset: `a-z 0-9 . - _ * ?`.",
									},
									"method_config": schema.ListAttribute{
										Optional:    true,
										ElementType: types.StringType,
										Description: "Request methods. Valid values: `HEAD`, `GET`, `POST`, `OPTIONS`, `PUT`, `PATCH`, `DELETE`. Required when `type` is `Method`.",
									},
									"path_config": schema.ListAttribute{
										Optional:    true,
										ElementType: types.StringType,
										Description: "Forwarding paths, each 1 ~ 128 characters, exact, regex and wildcard matching are supported. Exact matching charset: `a-z A-Z 0-9 . - _ / = :`. Regex matching starts with `~` (case-sensitive) or `~*` (case-insensitive), charset: `a-z A-Z 0-9 . - _ / = ? ~ ^ * $ : ( ) [ ] + |`. Wildcard matching: `*` for multi-character and `?` for single-character wildcards, charset: `a-z A-Z 0-9 . - _ / = :`.",
									},
									"query_string_config": schema.ListNestedAttribute{
										Optional:    true,
										Description: "Query string configuration. Required when `type` is `QueryString`.",
										NestedObject: schema.NestedAttributeObject{
											Attributes: map[string]schema.Attribute{
												"key": schema.StringAttribute{
													Required:    true,
													Description: "Key of the query string, 1 ~ 16 characters, printable characters except spaces and `#[]{}\\|<>&`. `*` matches multiple characters and `?` matches a single character.",
												},
												"value": schema.StringAttribute{
													Required:    true,
													Description: "Value of the query string, 1 ~ 128 characters, printable characters except spaces and `#[]{}\\|<>&`. `*` matches multiple characters and `?` matches a single character.",
												},
											},
										},
									},
									"source_ip_config": schema.ListAttribute{
										Optional:    true,
										ElementType: types.StringType,
										Description: "Source IP matching configuration in CIDR format, e.g. `x.x.x.x/32` for a single IP address or `x.x.x.x/24` for an IP range. Required when `type` is `SourceIp`.",
									},
								},
							},
						},
						"priority": schema.Int64Attribute{
							Required:    true,
							Description: "Priority of the rule, the smaller the value the higher the priority. Must be unique. Range: 1 ~ 10000.",
						},
						"direction": schema.StringAttribute{
							Optional:    true,
							Computed:    true,
							Description: "Direction of the forwarding rule. Valid values: `Request` (from client to load balancer, default), `Response` (from backend server to load balancer).",
						},
						"rule_name": schema.StringAttribute{
							Optional:    true,
							Computed:    true,
							Description: "Name of the forwarding rule, 1 ~ 255 characters, digits, letters, Chinese characters, `.`, `_` and `-` are supported.",
						},
						"tags": schema.ListNestedAttribute{
							Optional:    true,
							Description: "Tags of the forwarding rule.",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"tag_key": schema.StringAttribute{
										Required:    true,
										Description: "Key of the tag.",
									},
									"tag_value": schema.StringAttribute{
										Required:    true,
										Description: "Value of the tag.",
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

// Create handles `terraform apply` for new resources.
func (r *AlbRulesResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan AlbRulesResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	loadBalancerId := plan.LoadBalancerId.ValueString()
	listenerId := plan.ListenerId.ValueString()

	service, ok := r.service(&resp.Diagnostics)
	if !ok {
		return
	}

	request := albv20251030.NewCreateRulesRequest()
	request.LoadBalancerId = fwhelper.StringPointerFromValue(plan.LoadBalancerId)
	request.ListenerId = fwhelper.StringPointerFromValue(plan.ListenerId)

	rules, ruleDiags := expandAlbRules(ctx, plan.Rules)
	resp.Diagnostics.Append(ruleDiags...)
	if resp.Diagnostics.HasError() {
		return
	}
	request.Rules = rules

	ruleIds, err := service.CreateRules(ctx, request)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating ALB forwarding rule",
			fmt.Sprintf("CreateRules failed (%s): %s", albRulesScope(loadBalancerId, listenerId), err),
		)
		return
	}

	if len(ruleIds) != 1 {
		resp.Diagnostics.AddError(
			"Error creating ALB forwarding rule",
			fmt.Sprintf("CreateRules returned %d rule ids while exactly one rule is expected (%s): %v",
				len(ruleIds), albRulesScope(loadBalancerId, listenerId), ruleIds),
		)
		return
	}

	plan.Id = types.StringValue(albRulesId(loadBalancerId, listenerId, ruleIds[0]))

	// Persist the identity before reading back so a partial failure still
	// leaves a destroyable resource in state.
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	item := albReadRule(ctx, service, loadBalancerId, listenerId, ruleIds[0], &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	if item == nil {
		resp.Diagnostics.AddError(
			"ALB forwarding rule not readable after create",
			fmt.Sprintf("CreateRules succeeded but the rule cannot be read back (%s). Import it with `terraform import tencentcloud_alb_rules.<name> %s` to keep it managed.",
				albRulesScope(loadBalancerId, listenerId, ruleIds[0]), plan.Id.ValueString()),
		)
		return
	}

	// Tags are carried on the create request but not echoed by
	// DescribeRules, so the tags of the very first state are read from the
	// tag service.
	prevTags := types.ListNull(albTagObjectType())
	if !plan.Rules.IsNull() && !plan.Rules.IsUnknown() {
		var plannedRules []albRuleModel
		resp.Diagnostics.Append(plan.Rules.ElementsAs(ctx, &plannedRules, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		if len(plannedRules) == 1 {
			prevTags = plannedRules[0].Tags
		}
	}

	resp.Diagnostics.Append(refreshAlbRulesModel(ctx, &plan, loadBalancerId, listenerId, item, prevTags)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Read handles `terraform refresh` and post-apply state refresh.
func (r *AlbRulesResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state AlbRulesResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	loadBalancerId, listenerId, ruleId, ok := parseAlbRulesId(&resp.Diagnostics, state.Id.ValueString())
	if !ok {
		return
	}

	service, ok := r.service(&resp.Diagnostics)
	if !ok {
		return
	}

	item := albReadRule(ctx, service, loadBalancerId, listenerId, ruleId, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	// A confirmed "not found" means the rule was removed outside Terraform,
	// so it is dropped from state instead of failing the refresh.
	if item == nil {
		resp.Diagnostics.AddWarning(
			"Resource not found",
			fmt.Sprintf("ALB forwarding rule %s no longer exists, removing it from state.", state.Id.ValueString()),
		)
		resp.State.RemoveResource(ctx)
		return
	}

	// Tags are not echoed by DescribeRules, so the tags recorded in state
	// are refreshed from the tag service.
	prevTags := types.ListNull(albTagObjectType())
	if !state.Rules.IsNull() && !state.Rules.IsUnknown() {
		var stateRules []albRuleModel
		resp.Diagnostics.Append(state.Rules.ElementsAs(ctx, &stateRules, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		if len(stateRules) == 1 {
			prevTags = stateRules[0].Tags
		}
	}

	resp.Diagnostics.Append(refreshAlbRulesModel(ctx, &state, loadBalancerId, listenerId, item, prevTags)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update handles in-place updates through ModifyRulesAttributes. Listener and
// load balancer are marked RequiresReplace, so they never reach this method.
func (r *AlbRulesResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan AlbRulesResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state AlbRulesResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	loadBalancerId, listenerId, ruleId, ok := parseAlbRulesId(&resp.Diagnostics, state.Id.ValueString())
	if !ok {
		return
	}

	service, ok := r.service(&resp.Diagnostics)
	if !ok {
		return
	}

	request := albv20251030.NewModifyRulesAttributesRequest()
	request.LoadBalancerId = fwhelper.StringPointerFromValue(plan.LoadBalancerId)
	request.ListenerId = fwhelper.StringPointerFromValue(plan.ListenerId)

	var plannedRules []albRuleModel
	resp.Diagnostics.Append(plan.Rules.ElementsAs(ctx, &plannedRules, false)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if len(plannedRules) != 1 {
		resp.Diagnostics.AddError(
			"Error updating ALB forwarding rule",
			fmt.Sprintf("Exactly one rule is expected in `rules`, got %d (%s).", len(plannedRules), albRulesScope(loadBalancerId, listenerId, ruleId)),
		)
		return
	}

	ruleModify, modifyDiags := expandAlbRuleModifyRequest(ctx, &plannedRules[0], ruleId)
	resp.Diagnostics.Append(modifyDiags...)
	if resp.Diagnostics.HasError() {
		return
	}
	request.Rules = []*albv20251030.RuleModify{ruleModify}

	if err := service.ModifyRulesAttributes(ctx, request); err != nil {
		resp.Diagnostics.AddError(
			"Error updating ALB forwarding rule",
			fmt.Sprintf("ModifyRulesAttributes failed (%s): %s", albRulesScope(loadBalancerId, listenerId, ruleId), err),
		)
		return
	}

	// Tags are not covered by ModifyRulesAttributes, so a change is applied
	// through the tag service.
	var stateRules []albRuleModel
	resp.Diagnostics.Append(state.Rules.ElementsAs(ctx, &stateRules, false)...)
	if resp.Diagnostics.HasError() {
		return
	}
	stateTags := types.ListNull(albTagObjectType())
	if len(stateRules) == 1 {
		stateTags = stateRules[0].Tags
	}
	if albTagsChanged(plannedRules[0].Tags, stateTags) {
		if err := albApplyTagChange(ctx, r.Client(), loadBalancerId, listenerId, ruleId, plannedRules[0].Tags, stateTags); err != nil {
			resp.Diagnostics.AddError(
				"Error updating ALB forwarding rule tags",
				fmt.Sprintf("Modifying tags failed (%s): %s", albRulesScope(loadBalancerId, listenerId, ruleId), err),
			)
			return
		}
	}

	item := albReadRule(ctx, service, loadBalancerId, listenerId, ruleId, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	if item == nil {
		resp.Diagnostics.AddError(
			"ALB forwarding rule not readable after update",
			fmt.Sprintf("ModifyRulesAttributes succeeded but the rule cannot be read back (%s).", albRulesScope(loadBalancerId, listenerId, ruleId)),
		)
		return
	}

	plan.Id = state.Id
	resp.Diagnostics.Append(refreshAlbRulesModel(ctx, &plan, loadBalancerId, listenerId, item, plannedRules[0].Tags)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete handles resource deletion.
func (r *AlbRulesResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state AlbRulesResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	loadBalancerId, listenerId, ruleId, ok := parseAlbRulesId(&resp.Diagnostics, state.Id.ValueString())
	if !ok {
		return
	}

	service, ok := r.service(&resp.Diagnostics)
	if !ok {
		return
	}

	request := albv20251030.NewDeleteRulesRequest()
	request.LoadBalancerId = fwhelper.StringPointerFromValue(state.LoadBalancerId)
	request.ListenerId = fwhelper.StringPointerFromValue(state.ListenerId)
	request.RuleIds = []*string{helper.String(ruleId)}

	if err := service.DeleteRules(ctx, request); err != nil {
		// An already deleted rule is a successful deletion.
		if tccommon.IsExpectError(err, []string{"ResourceNotFound", "InvalidParameter"}) {
			return
		}
		resp.Diagnostics.AddError(
			"Error deleting ALB forwarding rule",
			fmt.Sprintf("DeleteRules failed (%s): %s", albRulesScope(loadBalancerId, listenerId, ruleId), err),
		)
	}
}

// ImportState supports `terraform import` using the
// `<load_balancer_id>#<listener_id>#<rule_id>` form. The remaining attributes
// are filled in by the Read that follows.
func (r *AlbRulesResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if _, _, _, ok := splitAlbRulesId(req.ID); !ok {
		resp.Diagnostics.AddError(
			"Invalid import ID",
			fmt.Sprintf("Import id %q is invalid, expected `<load_balancer_id>%s<listener_id>%s<rule_id>`.", req.ID, tccommon.FILED_SP, tccommon.FILED_SP),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
}

// service returns the ALB service bound to the shared client. When the
// provider has not been configured it records a diagnostic and reports false,
// so the CRUD methods never dereference a nil SDK client.
func (r *AlbRulesResource) service(diags *diag.Diagnostics) (AlbService, bool) {
	if r.Client() == nil {
		diags.AddError(
			"Provider not configured",
			"The provider client is not available. This is unexpected; please report it.",
		)
		return AlbService{}, false
	}
	return NewAlbService(r.Client()), true
}

// albReadRule reads one forwarding rule from the cloud. A nil item means the
// rule no longer exists; API and transport failures are recorded as
// diagnostics.
func albReadRule(ctx context.Context, service AlbService, loadBalancerId, listenerId, ruleId string, diags *diag.Diagnostics) *albv20251030.RuleOutput {
	request := albv20251030.NewDescribeRulesRequest()
	request.LoadBalancerId = helper.String(loadBalancerId)
	request.ListenerId = helper.String(listenerId)
	request.RuleIds = []*string{helper.String(ruleId)}

	rules, _, err := service.DescribeRules(ctx, request)
	if err != nil {
		diags.AddError(
			"Error reading ALB forwarding rule",
			fmt.Sprintf("DescribeRules failed (%s): %s", albRulesScope(loadBalancerId, listenerId, ruleId), err),
		)
		return nil
	}

	for _, rule := range rules {
		if rule != nil && rule.RuleId != nil && *rule.RuleId == ruleId {
			return rule
		}
	}
	return nil
}

// albRulesId builds the Terraform id from its three cloud parts.
func albRulesId(loadBalancerId, listenerId, ruleId string) string {
	return strings.Join([]string{loadBalancerId, listenerId, ruleId}, tccommon.FILED_SP)
}

// splitAlbRulesId parses `<load_balancer_id>#<listener_id>#<rule_id>`.
func splitAlbRulesId(id string) (string, string, string, bool) {
	parts := strings.Split(id, tccommon.FILED_SP)
	if len(parts) != 3 || parts[0] == "" || parts[1] == "" || parts[2] == "" {
		return "", "", "", false
	}
	return parts[0], parts[1], parts[2], true
}

// parseAlbRulesId resolves the Terraform id into the triple of ids the ALB
// APIs expect, recording a diagnostic when the id is malformed.
func parseAlbRulesId(diags *diag.Diagnostics, id string) (loadBalancerId, listenerId, ruleId string, ok bool) {
	loadBalancerId, listenerId, ruleId, ok = splitAlbRulesId(id)
	if !ok {
		diags.AddError(
			"Invalid resource ID",
			fmt.Sprintf("Resource id %q is invalid, expected `<load_balancer_id>%s<listener_id>%s<rule_id>`.", id, tccommon.FILED_SP, tccommon.FILED_SP),
		)
	}
	return
}

// albRulesScope renders the cloud-side identity of a rule for diagnostics, so
// every failure can be traced back to a concrete load balancer, listener and
// rule.
func albRulesScope(parts ...string) string {
	return strings.Join(parts, ", ")
}

// albTagsChanged reports whether the planned tags differ from the tags of the
// rule recorded in state.
func albTagsChanged(plannedTags, stateTags types.List) bool {
	if plannedTags.IsNull() || plannedTags.IsUnknown() {
		return false
	}
	if plannedTags.Equal(stateTags) {
		return false
	}
	return true
}

// albApplyTagChange aligns the tags recorded by the tag service with the
// planned tags: keys absent from the plan are deleted, the rest are set to the
// planned values.
func albApplyTagChange(ctx context.Context, client *connectivity.TencentCloudClient, loadBalancerId, listenerId, ruleId string, plannedTags, stateTags types.List) error {
	tagService := svctag.NewTagService(client)

	planned := map[string]string{}
	if !plannedTags.IsNull() && !plannedTags.IsUnknown() {
		var tags []albTagModel
		if diags := plannedTags.ElementsAs(ctx, &tags, false); diags.HasError() {
			return fmt.Errorf("cannot read planned tags: %v", diags)
		}
		for _, tag := range tags {
			planned[tag.TagKey.ValueString()] = tag.TagValue.ValueString()
		}
	}

	current := map[string]string{}
	if !stateTags.IsNull() && !stateTags.IsUnknown() {
		var tags []albTagModel
		if diags := stateTags.ElementsAs(ctx, &tags, false); diags.HasError() {
			return fmt.Errorf("cannot read current tags: %v", diags)
		}
		for _, tag := range tags {
			current[tag.TagKey.ValueString()] = tag.TagValue.ValueString()
		}
	}

	replaceTags := map[string]string{}
	deleteKeys := make([]string, 0)
	for key, value := range planned {
		if current[key] != value {
			replaceTags[key] = value
		}
	}
	for key := range current {
		if _, exists := planned[key]; !exists {
			deleteKeys = append(deleteKeys, key)
		}
	}

	if len(replaceTags) == 0 && len(deleteKeys) == 0 {
		return nil
	}

	resourceName := tccommon.BuildTagResourceName("clb", "listener", client.Region, listenerId)
	return tagService.ModifyTags(ctx, resourceName, replaceTags, deleteKeys)
}
