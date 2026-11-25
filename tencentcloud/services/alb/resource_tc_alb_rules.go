package alb

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	albv20251030 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/alb/v20251030"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/fw"
	fwHelper "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/helper"
)

var (
	_ resource.ResourceWithConfigure   = &AlbRulesResource{}
	_ resource.ResourceWithImportState = &AlbRulesResource{}
)

// NewAlbRulesResource is the factory referenced by tencentcloud/framework/registry.go.
func NewAlbRulesResource() resource.Resource { return &AlbRulesResource{} }

// AlbRulesResource implements resource.Resource for tencentcloud_alb_rules.
type AlbRulesResource struct {
	fw.ResourceWithConfigure
	fw.WithImportByID
}

// AlbRulesResourceModel maps the schema attributes.
type AlbRulesResourceModel struct {
	Id             types.String `tfsdk:"id"`
	ListenerId     types.String `tfsdk:"listener_id"`
	LoadBalancerId types.String `tfsdk:"load_balancer_id"`
	Rules          types.List   `tfsdk:"rules"`
	ClientToken    types.String `tfsdk:"client_token"`
	DryRun         types.Bool   `tfsdk:"dry_run"`
}

type RulesActionsFixedResponseConfigModel struct {
	HttpCode    types.Int64  `tfsdk:"http_code"`
	Content     types.String `tfsdk:"content"`
	ContentType types.String `tfsdk:"content_type"`
}

type RulesActionsInsertHeaderConfigModel struct {
	Key       types.String `tfsdk:"key"`
	Value     types.String `tfsdk:"value"`
	ValueType types.String `tfsdk:"value_type"`
}

type RulesActionsRedirectConfigModel struct {
	HttpCode types.Int64  `tfsdk:"http_code"`
	Host     types.String `tfsdk:"host"`
	Path     types.String `tfsdk:"path"`
	Port     types.String `tfsdk:"port"`
	Protocol types.String `tfsdk:"protocol"`
	Query    types.String `tfsdk:"query"`
}

type RulesActionsRemoveHeaderConfigModel struct {
	Key types.String `tfsdk:"key"`
}

type RulesActionsRewriteConfigModel struct {
	Host  types.String `tfsdk:"host"`
	Path  types.String `tfsdk:"path"`
	Query types.String `tfsdk:"query"`
}

type RulesActionsTargetGroupConfigTargetGroupsModel struct {
	TargetGroupId types.String `tfsdk:"target_group_id"`
	Weight        types.Int64  `tfsdk:"weight"`
}

type RulesActionsTargetGroupConfigTargetGroupStickySessionModel struct {
	Enabled types.Bool  `tfsdk:"enabled"`
	Timeout types.Int64 `tfsdk:"timeout"`
}

type RulesActionsTargetGroupConfigModel struct {
	TargetGroups             types.List   `tfsdk:"target_groups"`
	TargetGroupStickySession types.Object `tfsdk:"target_group_sticky_session"`
}

type RulesActionsModel struct {
	Order               types.Int64  `tfsdk:"order"`
	Type                types.String `tfsdk:"type"`
	FixedResponseConfig types.Object `tfsdk:"fixed_response_config"`
	InsertHeaderConfig  types.Object `tfsdk:"insert_header_config"`
	RedirectConfig      types.Object `tfsdk:"redirect_config"`
	RemoveHeaderConfig  types.Object `tfsdk:"remove_header_config"`
	RewriteConfig       types.Object `tfsdk:"rewrite_config"`
	TargetGroupConfig   types.Object `tfsdk:"target_group_config"`
}

type RulesConditionsCookieConfigModel struct {
	Key   types.String `tfsdk:"key"`
	Value types.String `tfsdk:"value"`
}

type RulesConditionsHeaderConfigModel struct {
	Key    types.String `tfsdk:"key"`
	Values types.List   `tfsdk:"values"`
}

type RulesConditionsQueryStringConfigModel struct {
	Key   types.String `tfsdk:"key"`
	Value types.String `tfsdk:"value"`
}

type RulesConditionsModel struct {
	Type              types.String `tfsdk:"type"`
	CookieConfig      types.List   `tfsdk:"cookie_config"`
	HeaderConfig      types.Object `tfsdk:"header_config"`
	HostConfig        types.List   `tfsdk:"host_config"`
	MethodConfig      types.List   `tfsdk:"method_config"`
	PathConfig        types.List   `tfsdk:"path_config"`
	QueryStringConfig types.List   `tfsdk:"query_string_config"`
	SourceIpConfig    types.List   `tfsdk:"source_ip_config"`
}

type RulesTagsModel struct {
	TagKey   types.String `tfsdk:"tag_key"`
	TagValue types.String `tfsdk:"tag_value"`
}

type RulesModel struct {
	Actions    types.List   `tfsdk:"actions"`
	Conditions types.List   `tfsdk:"conditions"`
	Priority   types.Int64  `tfsdk:"priority"`
	Direction  types.String `tfsdk:"direction"`
	RuleName   types.String `tfsdk:"rule_name"`
	Tags       types.List   `tfsdk:"tags"`
}

func (r *AlbRulesResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "tencentcloud_alb_rules"
}

func (r *AlbRulesResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a resource to manage a set of Application Load Balancer forwarding rules.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:      true,
				Description:   "ID of the resource.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"listener_id": schema.StringAttribute{
				Required:      true,
				Description:   "The listener ID, in the format lst- followed by 8 alphanumeric characters.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"load_balancer_id": schema.StringAttribute{
				Required:      true,
				Description:   "The load balancer instance ID, in the format alb- followed by 8 alphanumeric characters.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"rules": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"actions": schema.ListNestedAttribute{
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"order": schema.Int64Attribute{
										Required:      true,
										Description:   "The order in which forwarding actions are executed. It cannot be repeated and actions are executed in ascending order. Valid values: 1 ~ 50000.",
										PlanModifiers: []planmodifier.Int64{int64planmodifier.RequiresReplace()},
									},
									"type": schema.StringAttribute{
										Required:      true,
										Description:   "The forwarding action type. Valid values: TargetGroup: forward to a target group. Redirect: redirect. FixedResponse: return a fixed response. Rewrite: rewrite. InsertHeader: write an HTTP header. RemoveHeader: delete an HTTP header. A forwarding action must include one of TargetGroup, Redirect and FixedResponse, and it must be executed last.",
										PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
									},
									"fixed_response_config": schema.SingleNestedAttribute{
										Attributes: map[string]schema.Attribute{
											"http_code": schema.Int64Attribute{
												Required:      true,
												Description:   "The HTTP response code to return. 2xx, 4xx and 5xx are supported.",
												PlanModifiers: []planmodifier.Int64{int64planmodifier.RequiresReplace()},
											},
											"content": schema.StringAttribute{
												Optional:      true,
												Description:   "The fixed content to return. Only ASCII characters are supported, up to 1 KB.",
												PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
											},
											"content_type": schema.StringAttribute{
												Optional:      true,
												Description:   "The format of the fixed content to return. Valid values: text/plain, text/css, text/html, application/javascript or application/json.",
												PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
											},
										},
										Optional:      true,
										Description:   "Fixed response content configuration.",
										PlanModifiers: []planmodifier.Object{objectplanmodifier.RequiresReplace()},
									},
									"insert_header_config": schema.SingleNestedAttribute{
										Attributes: map[string]schema.Attribute{
											"key": schema.StringAttribute{
												Required:      true,
												Description:   "The key of the HTTP header to insert, 1 ~ 40 characters long. Supported characters: a-z A-Z 0-9 - _. Chinese characters are not supported, nor are Cookie, Host, Content-Length, Connection, Upgrade, transfer-encoding, keep-alive, te, authority, x-forwarded-for, x-forwarded-proto, x-forwarded-host and x-forwarded-port.",
												PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
											},
											"value": schema.StringAttribute{
												Required:      true,
												Description:   "The value of the HTTP header. When ValueType is SystemDefined, valid values are ClientPort: client port, ClientIp: client IP address, Protocol: protocol of the client request, CLBPort: load balancer listener port. When ValueType is UserDefined, it is 1 ~ 128 printable characters, does not support \", cannot start or end with a space, and cannot end with \\. When ValueType is ReferenceHeader, it references a header in the request, 1 ~ 128 printable characters, does not support \", cannot start or end with a space, and cannot end with \\.",
												PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
											},
											"value_type": schema.StringAttribute{
												Required:      true,
												Description:   "The value type of the HTTP header. Valid values: SystemDefined: a system-defined header. UserDefined: a user-defined header. ReferenceHeader: a header referenced from the request.",
												PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
											},
										},
										Optional:      true,
										Description:   "Insert HTTP header configuration.",
										PlanModifiers: []planmodifier.Object{objectplanmodifier.RequiresReplace()},
									},
									"redirect_config": schema.SingleNestedAttribute{
										Attributes: map[string]schema.Attribute{
											"http_code": schema.Int64Attribute{
												Required:      true,
												Description:   "The HTTP response code to return. 2xx, 4xx and 5xx are supported.",
												PlanModifiers: []planmodifier.Int64{int64planmodifier.RequiresReplace()},
											},
											"host": schema.StringAttribute{
												Optional:      true,
												Description:   "The host address for redirect/rewrite, with a default value of $${host}. It is 3 ~ 128 characters long. Supported characters: a-z 0-9 _ . -.",
												PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
											},
											"path": schema.StringAttribute{
												Optional:      true,
												Description:   "The path for redirect/rewrite, with a default value of $${path}. It is 1 ~ 128 characters long. Supported characters: a-z A-Z 0-9 ? = _ . - / :.",
												PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
											},
											"port": schema.StringAttribute{
												Optional:      true,
												Description:   "The port for redirection, with a default value of $${port}. Valid values: 1 ~ 65535.",
												PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
											},
											"protocol": schema.StringAttribute{
												Optional:      true,
												Description:   "The protocol for redirection. Valid values: HTTP, HTTPS. Default value: $${protocol}.",
												PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
											},
											"query": schema.StringAttribute{
												Optional:      true,
												Description:   "The query string for redirect/rewrite, with a default value of $${query}. It is 1 ~ 128 characters long and supports printable characters, but not #[]{}|<>& or spaces.",
												PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
											},
										},
										Optional:      true,
										Description:   "Redirection configuration. Except for HttpCode, the other configurations cannot all use default values.",
										PlanModifiers: []planmodifier.Object{objectplanmodifier.RequiresReplace()},
									},
									"remove_header_config": schema.SingleNestedAttribute{
										Attributes: map[string]schema.Attribute{
											"key": schema.StringAttribute{
												Required:      true,
												Description:   "The key of the HTTP header to insert, 1 ~ 40 characters long. Supported characters: a-z A-Z 0-9 - _. Chinese characters are not supported, nor are Cookie, Host, Content-Length, Connection, Upgrade, transfer-encoding, keep-alive, te, authority, x-forwarded-for, x-forwarded-proto, x-forwarded-host and x-forwarded-port.",
												PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
											},
										},
										Optional:      true,
										Description:   "Remove HTTP header configuration.",
										PlanModifiers: []planmodifier.Object{objectplanmodifier.RequiresReplace()},
									},
									"rewrite_config": schema.SingleNestedAttribute{
										Attributes: map[string]schema.Attribute{
											"host": schema.StringAttribute{
												Optional:      true,
												Description:   "The host address for redirect/rewrite, with a default value of $${host}. It is 3 ~ 128 characters long. Supported characters: a-z 0-9 _ . -.",
												PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
											},
											"path": schema.StringAttribute{
												Optional:      true,
												Description:   "The path for redirect/rewrite, with a default value of $${path}. It is 1 ~ 128 characters long. Supported characters: a-z A-Z 0-9 ? = _ . - / :.",
												PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
											},
											"query": schema.StringAttribute{
												Optional:      true,
												Description:   "The query string for redirect/rewrite, with a default value of $${query}. It is 1 ~ 128 characters long and supports printable characters, but not #[]{}|<>& or spaces.",
												PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
											},
										},
										Optional:      true,
										Description:   "Rewrite configuration.",
										PlanModifiers: []planmodifier.Object{objectplanmodifier.RequiresReplace()},
									},
									"target_group_config": schema.SingleNestedAttribute{
										Attributes: map[string]schema.Attribute{
											"target_groups": schema.ListNestedAttribute{
												NestedObject: schema.NestedAttributeObject{
													Attributes: map[string]schema.Attribute{
														"target_group_id": schema.StringAttribute{
															Required:      true,
															Description:   "The target group ID, in the format lbtg- followed by 8 alphanumeric characters.",
															PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
														},
														"weight": schema.Int64Attribute{
															Optional:      true,
															Description:   "The weight. Valid values: [0, 100]. Default value: 10.",
															PlanModifiers: []planmodifier.Int64{int64planmodifier.RequiresReplace()},
														},
													},
												},
												Required:      true,
												Description:   "Target group list.",
												PlanModifiers: []planmodifier.List{listplanmodifier.RequiresReplace()},
											},
											"target_group_sticky_session": schema.SingleNestedAttribute{
												Attributes: map[string]schema.Attribute{
													"enabled": schema.BoolAttribute{
														Required:      true,
														Description:   "Whether to enable sticky sessions. Disabled by default.",
														PlanModifiers: []planmodifier.Bool{boolplanmodifier.RequiresReplace()},
													},
													"timeout": schema.Int64Attribute{
														Optional:      true,
														Description:   "The timeout period in seconds. Valid values: 1 ~ 86400. Default value: 1000.",
														PlanModifiers: []planmodifier.Int64{int64planmodifier.RequiresReplace()},
													},
												},
												Optional:      true,
												Description:   "Sticky session between target groups.",
												PlanModifiers: []planmodifier.Object{objectplanmodifier.RequiresReplace()},
											},
										},
										Optional:      true,
										Description:   "Forwarding target group configuration.",
										PlanModifiers: []planmodifier.Object{objectplanmodifier.RequiresReplace()},
									},
								},
							},
							Required:      true,
							Description:   "Forwarding rule action list.",
							PlanModifiers: []planmodifier.List{listplanmodifier.RequiresReplace()},
						},
						"conditions": schema.ListNestedAttribute{
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"type": schema.StringAttribute{
										Required:      true,
										Description:   "The forwarding action type. Valid values: TargetGroup: forward to a target group. Redirect: redirect. FixedResponse: return a fixed response. Rewrite: rewrite. InsertHeader: write an HTTP header. RemoveHeader: delete an HTTP header. A forwarding action must include one of TargetGroup, Redirect and FixedResponse, and it must be executed last.",
										PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
									},
									"cookie_config": schema.ListNestedAttribute{
										NestedObject: schema.NestedAttributeObject{
											Attributes: map[string]schema.Attribute{
												"key": schema.StringAttribute{
													Required:      true,
													Description:   "The key of the HTTP header to insert, 1 ~ 40 characters long. Supported characters: a-z A-Z 0-9 - _. Chinese characters are not supported, nor are Cookie, Host, Content-Length, Connection, Upgrade, transfer-encoding, keep-alive, te, authority, x-forwarded-for, x-forwarded-proto, x-forwarded-host and x-forwarded-port.",
													PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
												},
												"value": schema.StringAttribute{
													Required:      true,
													Description:   "The value of the HTTP header. When ValueType is SystemDefined, valid values are ClientPort: client port, ClientIp: client IP address, Protocol: protocol of the client request, CLBPort: load balancer listener port. When ValueType is UserDefined, it is 1 ~ 128 printable characters, does not support \", cannot start or end with a space, and cannot end with \\. When ValueType is ReferenceHeader, it references a header in the request, 1 ~ 128 printable characters, does not support \", cannot start or end with a space, and cannot end with \\.",
													PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
												},
											},
										},
										Optional:      true,
										Description:   "Cookie configuration.",
										PlanModifiers: []planmodifier.List{listplanmodifier.RequiresReplace()},
									},
									"header_config": schema.SingleNestedAttribute{
										Attributes: map[string]schema.Attribute{
											"key": schema.StringAttribute{
												Required:      true,
												Description:   "The key of the HTTP header to insert, 1 ~ 40 characters long. Supported characters: a-z A-Z 0-9 - _. Chinese characters are not supported, nor are Cookie, Host, Content-Length, Connection, Upgrade, transfer-encoding, keep-alive, te, authority, x-forwarded-for, x-forwarded-proto, x-forwarded-host and x-forwarded-port.",
												PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
											},
											"values": schema.ListAttribute{
												ElementType:   types.StringType,
												Required:      true,
												Description:   "The value of the HTTP header, 1 ~ 128 characters long, supporting printable characters. \" is not supported; it cannot start or end with a space and cannot end with \\.",
												PlanModifiers: []planmodifier.List{listplanmodifier.RequiresReplace()},
											},
										},
										Optional:      true,
										Description:   "HTTP header configuration.",
										PlanModifiers: []planmodifier.Object{objectplanmodifier.RequiresReplace()},
									},
									"host_config": schema.ListAttribute{
										ElementType:   types.StringType,
										Optional:      true,
										Description:   "Host name. A host configuration can appear only once in a rule and is 3 ~ 128 characters long, supporting exact, regular expression and wildcard matching. It cannot start or end with a period (.) or underscore (_). Exact match supported character set: a-z 0-9 . - _. Regular expression match: starts with a tilde (~); supported characters: a-z 0-9 . - ? = ~ _ - + \\ ^ * ! $ & | ( ) [ ]. Wildcard match: * matches multiple characters and ? matches one character; supported characters: a-z 0-9 . - _ * ?.",
										PlanModifiers: []planmodifier.List{listplanmodifier.RequiresReplace()},
									},
									"method_config": schema.ListAttribute{
										ElementType:   types.StringType,
										Optional:      true,
										Description:   "The request method. Valid values: HEAD, GET, POST, OPTIONS, PUT, PATCH, DELETE.",
										PlanModifiers: []planmodifier.List{listplanmodifier.RequiresReplace()},
									},
									"path_config": schema.ListAttribute{
										ElementType:   types.StringType,
										Optional:      true,
										Description:   "The forwarding path, 1 ~ 128 characters long, supporting exact, regular expression and wildcard matching. Exact match supported character set: a-z A-Z 0-9 . - _ / = :. Regular expression match: starts with ~ (case-sensitive) or ~* (case-insensitive); supported characters: a-z A-Z 0-9 . - _ / = ? ~ ^ * $ : ( ) [ ] + |. Wildcard match: * matches multiple characters and ? matches one character; supported characters: a-z A-Z 0-9 . - _ / = :.",
										PlanModifiers: []planmodifier.List{listplanmodifier.RequiresReplace()},
									},
									"query_string_config": schema.ListNestedAttribute{
										NestedObject: schema.NestedAttributeObject{
											Attributes: map[string]schema.Attribute{
												"key": schema.StringAttribute{
													Required:      true,
													Description:   "The key of the HTTP header to insert, 1 ~ 40 characters long. Supported characters: a-z A-Z 0-9 - _. Chinese characters are not supported, nor are Cookie, Host, Content-Length, Connection, Upgrade, transfer-encoding, keep-alive, te, authority, x-forwarded-for, x-forwarded-proto, x-forwarded-host and x-forwarded-port.",
													PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
												},
												"value": schema.StringAttribute{
													Required:      true,
													Description:   "The value of the HTTP header. When ValueType is SystemDefined, valid values are ClientPort: client port, ClientIp: client IP address, Protocol: protocol of the client request, CLBPort: load balancer listener port. When ValueType is UserDefined, it is 1 ~ 128 printable characters, does not support \", cannot start or end with a space, and cannot end with \\. When ValueType is ReferenceHeader, it references a header in the request, 1 ~ 128 printable characters, does not support \", cannot start or end with a space, and cannot end with \\.",
													PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
												},
											},
										},
										Optional:      true,
										Description:   "Query string configuration.",
										PlanModifiers: []planmodifier.List{listplanmodifier.RequiresReplace()},
									},
									"source_ip_config": schema.ListAttribute{
										ElementType:   types.StringType,
										Optional:      true,
										Description:   "Source IP matching configuration, in CIDR format: IP address x.x.x.x/32 or IP range x.x.x.x/24.",
										PlanModifiers: []planmodifier.List{listplanmodifier.RequiresReplace()},
									},
								},
							},
							Required:      true,
							Description:   "Forwarding rule condition list.",
							PlanModifiers: []planmodifier.List{listplanmodifier.RequiresReplace()},
						},
						"priority": schema.Int64Attribute{
							Required:      true,
							Description:   "The priority. A smaller value means higher priority. Priorities cannot be duplicated. Valid values: 1 ~ 10000.",
							PlanModifiers: []planmodifier.Int64{int64planmodifier.RequiresReplace()},
						},
						"direction": schema.StringAttribute{
							Optional:      true,
							Description:   "The direction of the forwarding rule. Request: requests from the client to the load balancer. Response: responses from the backend server to the load balancer. Default value: Request.",
							PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
						},
						"rule_name": schema.StringAttribute{
							Optional:      true,
							Description:   "The name of the forwarding rule, 1 ~ 255 characters long. It can contain digits, uppercase and lowercase letters, Chinese characters, periods (.), underscores (_) and hyphens (-).",
							PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
						},
						"tags": schema.ListNestedAttribute{
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"tag_key": schema.StringAttribute{
										Required:      true,
										Description:   "The tag key.",
										PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
									},
									"tag_value": schema.StringAttribute{
										Required:      true,
										Description:   "The tag value.",
										PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
									},
								},
							},
							Optional:      true,
							Description:   "Tags.",
							PlanModifiers: []planmodifier.List{listplanmodifier.RequiresReplace()},
						},
					},
				},
				Required:      true,
				Description:   "Forwarding rule list.",
				PlanModifiers: []planmodifier.List{listplanmodifier.RequiresReplace()},
			},
			"client_token": schema.StringAttribute{
				Optional:      true,
				Description:   "The client token used to guarantee request idempotency. Generate a unique value on your client so that it is unique across requests. ClientToken supports only ASCII characters. If it is not specified, the API request RequestId is used as the ClientToken. Each API request has a different RequestId.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"dry_run": schema.BoolAttribute{
				Optional:      true,
				Description:   "Whether to only perform a pre-check of this request.",
				PlanModifiers: []planmodifier.Bool{boolplanmodifier.RequiresReplace()},
			},
		},
	}
}

func expandRulesActionsFixedResponseConfig(ctx context.Context, src *RulesActionsFixedResponseConfigModel, diags *diag.Diagnostics) *albv20251030.FixedResponseInfo {
	if src == nil {
		return nil
	}
	return &albv20251030.FixedResponseInfo{
		HttpCode:    fwHelper.Int64PointerFromValue(src.HttpCode),
		Content:     fwHelper.StringPointerFromValue(src.Content),
		ContentType: fwHelper.StringPointerFromValue(src.ContentType),
	}
}

func expandRulesActionsInsertHeaderConfig(ctx context.Context, src *RulesActionsInsertHeaderConfigModel, diags *diag.Diagnostics) *albv20251030.InsertHTTPHeaderInfo {
	if src == nil {
		return nil
	}
	return &albv20251030.InsertHTTPHeaderInfo{
		Key:       fwHelper.StringPointerFromValue(src.Key),
		Value:     fwHelper.StringPointerFromValue(src.Value),
		ValueType: fwHelper.StringPointerFromValue(src.ValueType),
	}
}

func expandRulesActionsRedirectConfig(ctx context.Context, src *RulesActionsRedirectConfigModel, diags *diag.Diagnostics) *albv20251030.HTTPRedirectInfo {
	if src == nil {
		return nil
	}
	return &albv20251030.HTTPRedirectInfo{
		HttpCode: fwHelper.Int64PointerFromValue(src.HttpCode),
		Host:     fwHelper.StringPointerFromValue(src.Host),
		Path:     fwHelper.StringPointerFromValue(src.Path),
		Port:     fwHelper.StringPointerFromValue(src.Port),
		Protocol: fwHelper.StringPointerFromValue(src.Protocol),
		Query:    fwHelper.StringPointerFromValue(src.Query),
	}
}

func expandRulesActionsRemoveHeaderConfig(ctx context.Context, src *RulesActionsRemoveHeaderConfigModel, diags *diag.Diagnostics) *albv20251030.RemoveHTTPHeaderInfo {
	if src == nil {
		return nil
	}
	return &albv20251030.RemoveHTTPHeaderInfo{
		Key: fwHelper.StringPointerFromValue(src.Key),
	}
}

func expandRulesActionsRewriteConfig(ctx context.Context, src *RulesActionsRewriteConfigModel, diags *diag.Diagnostics) *albv20251030.HTTPRewriteInfo {
	if src == nil {
		return nil
	}
	return &albv20251030.HTTPRewriteInfo{
		Host:  fwHelper.StringPointerFromValue(src.Host),
		Path:  fwHelper.StringPointerFromValue(src.Path),
		Query: fwHelper.StringPointerFromValue(src.Query),
	}
}

func expandRulesActionsTargetGroupConfigTargetGroups(ctx context.Context, src *RulesActionsTargetGroupConfigTargetGroupsModel, diags *diag.Diagnostics) *albv20251030.TargetGroupTuple {
	if src == nil {
		return nil
	}
	return &albv20251030.TargetGroupTuple{
		TargetGroupId: fwHelper.StringPointerFromValue(src.TargetGroupId),
		Weight:        fwHelper.Int64PointerFromValue(src.Weight),
	}
}

func expandRulesActionsTargetGroupConfigTargetGroupsList(ctx context.Context, v types.List, diags *diag.Diagnostics) []*albv20251030.TargetGroupTuple {
	if v.IsNull() || v.IsUnknown() {
		return nil
	}
	var models []RulesActionsTargetGroupConfigTargetGroupsModel
	diags.Append(v.ElementsAs(ctx, &models, false)...)
	if diags.HasError() {
		return nil
	}
	result := make([]*albv20251030.TargetGroupTuple, 0, len(models))
	for i := range models {
		result = append(result, expandRulesActionsTargetGroupConfigTargetGroups(ctx, &models[i], diags))
	}
	return result
}

func expandRulesActionsTargetGroupConfigTargetGroupStickySession(ctx context.Context, src *RulesActionsTargetGroupConfigTargetGroupStickySessionModel, diags *diag.Diagnostics) *albv20251030.TargetGroupStickySession {
	if src == nil {
		return nil
	}
	return &albv20251030.TargetGroupStickySession{
		Enabled: fwHelper.BoolPointerFromValue(src.Enabled),
		Timeout: fwHelper.Int64PointerFromValue(src.Timeout),
	}
}

func expandRulesActionsTargetGroupConfig(ctx context.Context, src *RulesActionsTargetGroupConfigModel, diags *diag.Diagnostics) *albv20251030.TargetGroupConfig {
	if src == nil {
		return nil
	}
	return &albv20251030.TargetGroupConfig{
		TargetGroups:             expandRulesActionsTargetGroupConfigTargetGroupsList(ctx, src.TargetGroups, diags),
		TargetGroupStickySession: expandRulesActionsTargetGroupConfigTargetGroupStickySession(ctx, objRulesActionsTargetGroupConfigTargetGroupStickySession(ctx, src.TargetGroupStickySession, diags), diags),
	}
}

func objRulesActionsTargetGroupConfigTargetGroupStickySession(ctx context.Context, v types.Object, diags *diag.Diagnostics) *RulesActionsTargetGroupConfigTargetGroupStickySessionModel {
	if v.IsNull() || v.IsUnknown() {
		return nil
	}
	var m RulesActionsTargetGroupConfigTargetGroupStickySessionModel
	diags.Append(v.As(ctx, &m, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return nil
	}
	return &m
}

func expandRulesActions(ctx context.Context, src *RulesActionsModel, diags *diag.Diagnostics) *albv20251030.RuleAction {
	if src == nil {
		return nil
	}
	return &albv20251030.RuleAction{
		Order:               fwHelper.Int64PointerFromValue(src.Order),
		Type:                fwHelper.StringPointerFromValue(src.Type),
		FixedResponseConfig: expandRulesActionsFixedResponseConfig(ctx, objRulesActionsFixedResponseConfig(ctx, src.FixedResponseConfig, diags), diags),
		InsertHeaderConfig:  expandRulesActionsInsertHeaderConfig(ctx, objRulesActionsInsertHeaderConfig(ctx, src.InsertHeaderConfig, diags), diags),
		RedirectConfig:      expandRulesActionsRedirectConfig(ctx, objRulesActionsRedirectConfig(ctx, src.RedirectConfig, diags), diags),
		RemoveHeaderConfig:  expandRulesActionsRemoveHeaderConfig(ctx, objRulesActionsRemoveHeaderConfig(ctx, src.RemoveHeaderConfig, diags), diags),
		RewriteConfig:       expandRulesActionsRewriteConfig(ctx, objRulesActionsRewriteConfig(ctx, src.RewriteConfig, diags), diags),
		TargetGroupConfig:   expandRulesActionsTargetGroupConfig(ctx, objRulesActionsTargetGroupConfig(ctx, src.TargetGroupConfig, diags), diags),
	}
}

func expandRulesActionsList(ctx context.Context, v types.List, diags *diag.Diagnostics) []*albv20251030.RuleAction {
	if v.IsNull() || v.IsUnknown() {
		return nil
	}
	var models []RulesActionsModel
	diags.Append(v.ElementsAs(ctx, &models, false)...)
	if diags.HasError() {
		return nil
	}
	result := make([]*albv20251030.RuleAction, 0, len(models))
	for i := range models {
		result = append(result, expandRulesActions(ctx, &models[i], diags))
	}
	return result
}

func objRulesActionsFixedResponseConfig(ctx context.Context, v types.Object, diags *diag.Diagnostics) *RulesActionsFixedResponseConfigModel {
	if v.IsNull() || v.IsUnknown() {
		return nil
	}
	var m RulesActionsFixedResponseConfigModel
	diags.Append(v.As(ctx, &m, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return nil
	}
	return &m
}

func objRulesActionsInsertHeaderConfig(ctx context.Context, v types.Object, diags *diag.Diagnostics) *RulesActionsInsertHeaderConfigModel {
	if v.IsNull() || v.IsUnknown() {
		return nil
	}
	var m RulesActionsInsertHeaderConfigModel
	diags.Append(v.As(ctx, &m, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return nil
	}
	return &m
}

func objRulesActionsRedirectConfig(ctx context.Context, v types.Object, diags *diag.Diagnostics) *RulesActionsRedirectConfigModel {
	if v.IsNull() || v.IsUnknown() {
		return nil
	}
	var m RulesActionsRedirectConfigModel
	diags.Append(v.As(ctx, &m, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return nil
	}
	return &m
}

func objRulesActionsRemoveHeaderConfig(ctx context.Context, v types.Object, diags *diag.Diagnostics) *RulesActionsRemoveHeaderConfigModel {
	if v.IsNull() || v.IsUnknown() {
		return nil
	}
	var m RulesActionsRemoveHeaderConfigModel
	diags.Append(v.As(ctx, &m, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return nil
	}
	return &m
}

func objRulesActionsRewriteConfig(ctx context.Context, v types.Object, diags *diag.Diagnostics) *RulesActionsRewriteConfigModel {
	if v.IsNull() || v.IsUnknown() {
		return nil
	}
	var m RulesActionsRewriteConfigModel
	diags.Append(v.As(ctx, &m, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return nil
	}
	return &m
}

func objRulesActionsTargetGroupConfig(ctx context.Context, v types.Object, diags *diag.Diagnostics) *RulesActionsTargetGroupConfigModel {
	if v.IsNull() || v.IsUnknown() {
		return nil
	}
	var m RulesActionsTargetGroupConfigModel
	diags.Append(v.As(ctx, &m, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return nil
	}
	return &m
}

func expandRulesConditionsCookieConfig(ctx context.Context, src *RulesConditionsCookieConfigModel, diags *diag.Diagnostics) *albv20251030.HTTPCookieInfo {
	if src == nil {
		return nil
	}
	return &albv20251030.HTTPCookieInfo{
		Key:   fwHelper.StringPointerFromValue(src.Key),
		Value: fwHelper.StringPointerFromValue(src.Value),
	}
}

func expandRulesConditionsCookieConfigList(ctx context.Context, v types.List, diags *diag.Diagnostics) []*albv20251030.HTTPCookieInfo {
	if v.IsNull() || v.IsUnknown() {
		return nil
	}
	var models []RulesConditionsCookieConfigModel
	diags.Append(v.ElementsAs(ctx, &models, false)...)
	if diags.HasError() {
		return nil
	}
	result := make([]*albv20251030.HTTPCookieInfo, 0, len(models))
	for i := range models {
		result = append(result, expandRulesConditionsCookieConfig(ctx, &models[i], diags))
	}
	return result
}

func expandRulesConditionsHeaderConfig(ctx context.Context, src *RulesConditionsHeaderConfigModel, diags *diag.Diagnostics) *albv20251030.HTTPHeaderInfo {
	if src == nil {
		return nil
	}
	return &albv20251030.HTTPHeaderInfo{
		Key:    fwHelper.StringPointerFromValue(src.Key),
		Values: fwHelper.PointersFromList[string](ctx, src.Values, diags),
	}
}

func expandRulesConditionsQueryStringConfig(ctx context.Context, src *RulesConditionsQueryStringConfigModel, diags *diag.Diagnostics) *albv20251030.HTTPQueryStringInfo {
	if src == nil {
		return nil
	}
	return &albv20251030.HTTPQueryStringInfo{
		Key:   fwHelper.StringPointerFromValue(src.Key),
		Value: fwHelper.StringPointerFromValue(src.Value),
	}
}

func expandRulesConditionsQueryStringConfigList(ctx context.Context, v types.List, diags *diag.Diagnostics) []*albv20251030.HTTPQueryStringInfo {
	if v.IsNull() || v.IsUnknown() {
		return nil
	}
	var models []RulesConditionsQueryStringConfigModel
	diags.Append(v.ElementsAs(ctx, &models, false)...)
	if diags.HasError() {
		return nil
	}
	result := make([]*albv20251030.HTTPQueryStringInfo, 0, len(models))
	for i := range models {
		result = append(result, expandRulesConditionsQueryStringConfig(ctx, &models[i], diags))
	}
	return result
}

func expandRulesConditions(ctx context.Context, src *RulesConditionsModel, diags *diag.Diagnostics) *albv20251030.RuleCondition {
	if src == nil {
		return nil
	}
	return &albv20251030.RuleCondition{
		Type:              fwHelper.StringPointerFromValue(src.Type),
		CookieConfig:      expandRulesConditionsCookieConfigList(ctx, src.CookieConfig, diags),
		HeaderConfig:      expandRulesConditionsHeaderConfig(ctx, objRulesConditionsHeaderConfig(ctx, src.HeaderConfig, diags), diags),
		HostConfig:        fwHelper.PointersFromList[string](ctx, src.HostConfig, diags),
		MethodConfig:      fwHelper.PointersFromList[string](ctx, src.MethodConfig, diags),
		PathConfig:        fwHelper.PointersFromList[string](ctx, src.PathConfig, diags),
		QueryStringConfig: expandRulesConditionsQueryStringConfigList(ctx, src.QueryStringConfig, diags),
		SourceIpConfig:    fwHelper.PointersFromList[string](ctx, src.SourceIpConfig, diags),
	}
}

func expandRulesConditionsList(ctx context.Context, v types.List, diags *diag.Diagnostics) []*albv20251030.RuleCondition {
	if v.IsNull() || v.IsUnknown() {
		return nil
	}
	var models []RulesConditionsModel
	diags.Append(v.ElementsAs(ctx, &models, false)...)
	if diags.HasError() {
		return nil
	}
	result := make([]*albv20251030.RuleCondition, 0, len(models))
	for i := range models {
		result = append(result, expandRulesConditions(ctx, &models[i], diags))
	}
	return result
}

func objRulesConditionsHeaderConfig(ctx context.Context, v types.Object, diags *diag.Diagnostics) *RulesConditionsHeaderConfigModel {
	if v.IsNull() || v.IsUnknown() {
		return nil
	}
	var m RulesConditionsHeaderConfigModel
	diags.Append(v.As(ctx, &m, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return nil
	}
	return &m
}

func expandRulesTags(ctx context.Context, src *RulesTagsModel, diags *diag.Diagnostics) *albv20251030.TagInfo {
	if src == nil {
		return nil
	}
	return &albv20251030.TagInfo{
		TagKey:   fwHelper.StringPointerFromValue(src.TagKey),
		TagValue: fwHelper.StringPointerFromValue(src.TagValue),
	}
}

func expandRulesTagsList(ctx context.Context, v types.List, diags *diag.Diagnostics) []*albv20251030.TagInfo {
	if v.IsNull() || v.IsUnknown() {
		return nil
	}
	var models []RulesTagsModel
	diags.Append(v.ElementsAs(ctx, &models, false)...)
	if diags.HasError() {
		return nil
	}
	result := make([]*albv20251030.TagInfo, 0, len(models))
	for i := range models {
		result = append(result, expandRulesTags(ctx, &models[i], diags))
	}
	return result
}

func expandRules(ctx context.Context, src *RulesModel, diags *diag.Diagnostics) *albv20251030.RuleInput {
	if src == nil {
		return nil
	}
	return &albv20251030.RuleInput{
		Actions:    expandRulesActionsList(ctx, src.Actions, diags),
		Conditions: expandRulesConditionsList(ctx, src.Conditions, diags),
		Priority:   fwHelper.Int64PointerFromValue(src.Priority),
		Direction:  fwHelper.StringPointerFromValue(src.Direction),
		RuleName:   fwHelper.StringPointerFromValue(src.RuleName),
		Tags:       expandRulesTagsList(ctx, src.Tags, diags),
	}
}

func expandRulesList(ctx context.Context, v types.List, diags *diag.Diagnostics) []*albv20251030.RuleInput {
	if v.IsNull() || v.IsUnknown() {
		return nil
	}
	var models []RulesModel
	diags.Append(v.ElementsAs(ctx, &models, false)...)
	if diags.HasError() {
		return nil
	}
	result := make([]*albv20251030.RuleInput, 0, len(models))
	for i := range models {
		result = append(result, expandRules(ctx, &models[i], diags))
	}
	return result
}

func flattenRulesActionsFixedResponseConfig(ctx context.Context, src *albv20251030.FixedResponseInfo, diags *diag.Diagnostics) *RulesActionsFixedResponseConfigModel {
	if src == nil {
		return nil
	}
	return &RulesActionsFixedResponseConfigModel{
		HttpCode:    fwHelper.Int64ValueOrPrev(src.HttpCode, types.Int64Null()),
		Content:     fwHelper.StringValueOrPrev(src.Content, types.StringNull()),
		ContentType: fwHelper.StringValueOrPrev(src.ContentType, types.StringNull()),
	}
}

func flattenRulesActionsInsertHeaderConfig(ctx context.Context, src *albv20251030.InsertHTTPHeaderInfo, diags *diag.Diagnostics) *RulesActionsInsertHeaderConfigModel {
	if src == nil {
		return nil
	}
	return &RulesActionsInsertHeaderConfigModel{
		Key:       fwHelper.StringValueOrPrev(src.Key, types.StringNull()),
		Value:     fwHelper.StringValueOrPrev(src.Value, types.StringNull()),
		ValueType: fwHelper.StringValueOrPrev(src.ValueType, types.StringNull()),
	}
}

func flattenRulesActionsRedirectConfig(ctx context.Context, src *albv20251030.HTTPRedirectInfo, diags *diag.Diagnostics) *RulesActionsRedirectConfigModel {
	if src == nil {
		return nil
	}
	return &RulesActionsRedirectConfigModel{
		HttpCode: fwHelper.Int64ValueOrPrev(src.HttpCode, types.Int64Null()),
		Host:     fwHelper.StringValueOrPrev(src.Host, types.StringNull()),
		Path:     fwHelper.StringValueOrPrev(src.Path, types.StringNull()),
		Port:     fwHelper.StringValueOrPrev(src.Port, types.StringNull()),
		Protocol: fwHelper.StringValueOrPrev(src.Protocol, types.StringNull()),
		Query:    fwHelper.StringValueOrPrev(src.Query, types.StringNull()),
	}
}

func flattenRulesActionsRemoveHeaderConfig(ctx context.Context, src *albv20251030.RemoveHTTPHeaderInfo, diags *diag.Diagnostics) *RulesActionsRemoveHeaderConfigModel {
	if src == nil {
		return nil
	}
	return &RulesActionsRemoveHeaderConfigModel{
		Key: fwHelper.StringValueOrPrev(src.Key, types.StringNull()),
	}
}

func flattenRulesActionsRewriteConfig(ctx context.Context, src *albv20251030.HTTPRewriteInfo, diags *diag.Diagnostics) *RulesActionsRewriteConfigModel {
	if src == nil {
		return nil
	}
	return &RulesActionsRewriteConfigModel{
		Host:  fwHelper.StringValueOrPrev(src.Host, types.StringNull()),
		Path:  fwHelper.StringValueOrPrev(src.Path, types.StringNull()),
		Query: fwHelper.StringValueOrPrev(src.Query, types.StringNull()),
	}
}

func flattenRulesActionsTargetGroupConfigTargetGroups(ctx context.Context, src *albv20251030.TargetGroupTuple, diags *diag.Diagnostics) *RulesActionsTargetGroupConfigTargetGroupsModel {
	if src == nil {
		return nil
	}
	return &RulesActionsTargetGroupConfigTargetGroupsModel{
		TargetGroupId: fwHelper.StringValueOrPrev(src.TargetGroupId, types.StringNull()),
		Weight:        fwHelper.Int64ValueOrPrev(src.Weight, types.Int64Null()),
	}
}

func flattenRulesActionsTargetGroupConfigTargetGroupsList(ctx context.Context, src []*albv20251030.TargetGroupTuple, prev types.List, diags *diag.Diagnostics) types.List {
	elemType := types.ObjectType{AttrTypes: RulesActionsTargetGroupConfigTargetGroupsAttrTypes}
	if len(src) == 0 {
		if prev.IsNull() || len(prev.Elements()) == 0 {
			return types.ListNull(elemType)
		}
		return prev
	}
	vals := make([]attr.Value, 0, len(src))
	for _, item := range src {
		obj, d := types.ObjectValueFrom(ctx, RulesActionsTargetGroupConfigTargetGroupsAttrTypes, flattenRulesActionsTargetGroupConfigTargetGroups(ctx, item, diags))
		diags.Append(d...)
		vals = append(vals, obj)
	}
	return types.ListValueMust(elemType, vals)
}

func flattenRulesActionsTargetGroupConfigTargetGroupStickySession(ctx context.Context, src *albv20251030.TargetGroupStickySession, diags *diag.Diagnostics) *RulesActionsTargetGroupConfigTargetGroupStickySessionModel {
	if src == nil {
		return nil
	}
	return &RulesActionsTargetGroupConfigTargetGroupStickySessionModel{
		Enabled: fwHelper.BoolValueOrPrev(src.Enabled, types.BoolNull()),
		Timeout: fwHelper.Int64ValueOrPrev(src.Timeout, types.Int64Null()),
	}
}

func flattenRulesActionsTargetGroupConfig(ctx context.Context, src *albv20251030.TargetGroupConfig, diags *diag.Diagnostics) *RulesActionsTargetGroupConfigModel {
	if src == nil {
		return nil
	}
	return &RulesActionsTargetGroupConfigModel{
		TargetGroups:             flattenRulesActionsTargetGroupConfigTargetGroupsList(ctx, src.TargetGroups, types.ListNull(types.ObjectType{AttrTypes: RulesActionsTargetGroupConfigTargetGroupsAttrTypes}), diags),
		TargetGroupStickySession: objFromRulesActionsTargetGroupConfigTargetGroupStickySession(ctx, flattenRulesActionsTargetGroupConfigTargetGroupStickySession(ctx, src.TargetGroupStickySession, diags)),
	}
}

func objFromRulesActionsTargetGroupConfigTargetGroupStickySession(ctx context.Context, m *RulesActionsTargetGroupConfigTargetGroupStickySessionModel) types.Object {
	if m == nil {
		return types.ObjectNull(RulesActionsTargetGroupConfigTargetGroupStickySessionAttrTypes)
	}
	v, d := types.ObjectValueFrom(ctx, RulesActionsTargetGroupConfigTargetGroupStickySessionAttrTypes, m)
	if d.HasError() {
		return types.ObjectNull(RulesActionsTargetGroupConfigTargetGroupStickySessionAttrTypes)
	}
	return v
}

func flattenRulesActions(ctx context.Context, src *albv20251030.RuleAction, diags *diag.Diagnostics) *RulesActionsModel {
	if src == nil {
		return nil
	}
	return &RulesActionsModel{
		Order:               fwHelper.Int64ValueOrPrev(src.Order, types.Int64Null()),
		Type:                fwHelper.StringValueOrPrev(src.Type, types.StringNull()),
		FixedResponseConfig: objFromRulesActionsFixedResponseConfig(ctx, flattenRulesActionsFixedResponseConfig(ctx, src.FixedResponseConfig, diags)),
		InsertHeaderConfig:  objFromRulesActionsInsertHeaderConfig(ctx, flattenRulesActionsInsertHeaderConfig(ctx, src.InsertHeaderConfig, diags)),
		RedirectConfig:      objFromRulesActionsRedirectConfig(ctx, flattenRulesActionsRedirectConfig(ctx, src.RedirectConfig, diags)),
		RemoveHeaderConfig:  objFromRulesActionsRemoveHeaderConfig(ctx, flattenRulesActionsRemoveHeaderConfig(ctx, src.RemoveHeaderConfig, diags)),
		RewriteConfig:       objFromRulesActionsRewriteConfig(ctx, flattenRulesActionsRewriteConfig(ctx, src.RewriteConfig, diags)),
		TargetGroupConfig:   objFromRulesActionsTargetGroupConfig(ctx, flattenRulesActionsTargetGroupConfig(ctx, src.TargetGroupConfig, diags)),
	}
}

func flattenRulesActionsList(ctx context.Context, src []*albv20251030.RuleAction, prev types.List, diags *diag.Diagnostics) types.List {
	elemType := types.ObjectType{AttrTypes: RulesActionsAttrTypes}
	if len(src) == 0 {
		if prev.IsNull() || len(prev.Elements()) == 0 {
			return types.ListNull(elemType)
		}
		return prev
	}
	vals := make([]attr.Value, 0, len(src))
	for _, item := range src {
		obj, d := types.ObjectValueFrom(ctx, RulesActionsAttrTypes, flattenRulesActions(ctx, item, diags))
		diags.Append(d...)
		vals = append(vals, obj)
	}
	return types.ListValueMust(elemType, vals)
}

func objFromRulesActionsFixedResponseConfig(ctx context.Context, m *RulesActionsFixedResponseConfigModel) types.Object {
	if m == nil {
		return types.ObjectNull(RulesActionsFixedResponseConfigAttrTypes)
	}
	v, d := types.ObjectValueFrom(ctx, RulesActionsFixedResponseConfigAttrTypes, m)
	if d.HasError() {
		return types.ObjectNull(RulesActionsFixedResponseConfigAttrTypes)
	}
	return v
}

func objFromRulesActionsInsertHeaderConfig(ctx context.Context, m *RulesActionsInsertHeaderConfigModel) types.Object {
	if m == nil {
		return types.ObjectNull(RulesActionsInsertHeaderConfigAttrTypes)
	}
	v, d := types.ObjectValueFrom(ctx, RulesActionsInsertHeaderConfigAttrTypes, m)
	if d.HasError() {
		return types.ObjectNull(RulesActionsInsertHeaderConfigAttrTypes)
	}
	return v
}

func objFromRulesActionsRedirectConfig(ctx context.Context, m *RulesActionsRedirectConfigModel) types.Object {
	if m == nil {
		return types.ObjectNull(RulesActionsRedirectConfigAttrTypes)
	}
	v, d := types.ObjectValueFrom(ctx, RulesActionsRedirectConfigAttrTypes, m)
	if d.HasError() {
		return types.ObjectNull(RulesActionsRedirectConfigAttrTypes)
	}
	return v
}

func objFromRulesActionsRemoveHeaderConfig(ctx context.Context, m *RulesActionsRemoveHeaderConfigModel) types.Object {
	if m == nil {
		return types.ObjectNull(RulesActionsRemoveHeaderConfigAttrTypes)
	}
	v, d := types.ObjectValueFrom(ctx, RulesActionsRemoveHeaderConfigAttrTypes, m)
	if d.HasError() {
		return types.ObjectNull(RulesActionsRemoveHeaderConfigAttrTypes)
	}
	return v
}

func objFromRulesActionsRewriteConfig(ctx context.Context, m *RulesActionsRewriteConfigModel) types.Object {
	if m == nil {
		return types.ObjectNull(RulesActionsRewriteConfigAttrTypes)
	}
	v, d := types.ObjectValueFrom(ctx, RulesActionsRewriteConfigAttrTypes, m)
	if d.HasError() {
		return types.ObjectNull(RulesActionsRewriteConfigAttrTypes)
	}
	return v
}

func objFromRulesActionsTargetGroupConfig(ctx context.Context, m *RulesActionsTargetGroupConfigModel) types.Object {
	if m == nil {
		return types.ObjectNull(RulesActionsTargetGroupConfigAttrTypes)
	}
	v, d := types.ObjectValueFrom(ctx, RulesActionsTargetGroupConfigAttrTypes, m)
	if d.HasError() {
		return types.ObjectNull(RulesActionsTargetGroupConfigAttrTypes)
	}
	return v
}

func flattenRulesConditionsCookieConfig(ctx context.Context, src *albv20251030.HTTPCookieInfo, diags *diag.Diagnostics) *RulesConditionsCookieConfigModel {
	if src == nil {
		return nil
	}
	return &RulesConditionsCookieConfigModel{
		Key:   fwHelper.StringValueOrPrev(src.Key, types.StringNull()),
		Value: fwHelper.StringValueOrPrev(src.Value, types.StringNull()),
	}
}

func flattenRulesConditionsCookieConfigList(ctx context.Context, src []*albv20251030.HTTPCookieInfo, prev types.List, diags *diag.Diagnostics) types.List {
	elemType := types.ObjectType{AttrTypes: RulesConditionsCookieConfigAttrTypes}
	if len(src) == 0 {
		if prev.IsNull() || len(prev.Elements()) == 0 {
			return types.ListNull(elemType)
		}
		return prev
	}
	vals := make([]attr.Value, 0, len(src))
	for _, item := range src {
		obj, d := types.ObjectValueFrom(ctx, RulesConditionsCookieConfigAttrTypes, flattenRulesConditionsCookieConfig(ctx, item, diags))
		diags.Append(d...)
		vals = append(vals, obj)
	}
	return types.ListValueMust(elemType, vals)
}

func flattenRulesConditionsHeaderConfig(ctx context.Context, src *albv20251030.HTTPHeaderInfo, diags *diag.Diagnostics) *RulesConditionsHeaderConfigModel {
	if src == nil {
		return nil
	}
	return &RulesConditionsHeaderConfigModel{
		Key:    fwHelper.StringValueOrPrev(src.Key, types.StringNull()),
		Values: fwHelper.ListFromPointers(ctx, types.StringType, src.Values, diags),
	}
}

func flattenRulesConditionsQueryStringConfig(ctx context.Context, src *albv20251030.HTTPQueryStringInfo, diags *diag.Diagnostics) *RulesConditionsQueryStringConfigModel {
	if src == nil {
		return nil
	}
	return &RulesConditionsQueryStringConfigModel{
		Key:   fwHelper.StringValueOrPrev(src.Key, types.StringNull()),
		Value: fwHelper.StringValueOrPrev(src.Value, types.StringNull()),
	}
}

func flattenRulesConditionsQueryStringConfigList(ctx context.Context, src []*albv20251030.HTTPQueryStringInfo, prev types.List, diags *diag.Diagnostics) types.List {
	elemType := types.ObjectType{AttrTypes: RulesConditionsQueryStringConfigAttrTypes}
	if len(src) == 0 {
		if prev.IsNull() || len(prev.Elements()) == 0 {
			return types.ListNull(elemType)
		}
		return prev
	}
	vals := make([]attr.Value, 0, len(src))
	for _, item := range src {
		obj, d := types.ObjectValueFrom(ctx, RulesConditionsQueryStringConfigAttrTypes, flattenRulesConditionsQueryStringConfig(ctx, item, diags))
		diags.Append(d...)
		vals = append(vals, obj)
	}
	return types.ListValueMust(elemType, vals)
}

func flattenRulesConditions(ctx context.Context, src *albv20251030.RuleCondition, diags *diag.Diagnostics) *RulesConditionsModel {
	if src == nil {
		return nil
	}
	return &RulesConditionsModel{
		Type:              fwHelper.StringValueOrPrev(src.Type, types.StringNull()),
		CookieConfig:      flattenRulesConditionsCookieConfigList(ctx, src.CookieConfig, types.ListNull(types.ObjectType{AttrTypes: RulesConditionsCookieConfigAttrTypes}), diags),
		HeaderConfig:      objFromRulesConditionsHeaderConfig(ctx, flattenRulesConditionsHeaderConfig(ctx, src.HeaderConfig, diags)),
		HostConfig:        fwHelper.ListFromPointers(ctx, types.StringType, src.HostConfig, diags),
		MethodConfig:      fwHelper.ListFromPointers(ctx, types.StringType, src.MethodConfig, diags),
		PathConfig:        fwHelper.ListFromPointers(ctx, types.StringType, src.PathConfig, diags),
		QueryStringConfig: flattenRulesConditionsQueryStringConfigList(ctx, src.QueryStringConfig, types.ListNull(types.ObjectType{AttrTypes: RulesConditionsQueryStringConfigAttrTypes}), diags),
		SourceIpConfig:    fwHelper.ListFromPointers(ctx, types.StringType, src.SourceIpConfig, diags),
	}
}

func flattenRulesConditionsList(ctx context.Context, src []*albv20251030.RuleCondition, prev types.List, diags *diag.Diagnostics) types.List {
	elemType := types.ObjectType{AttrTypes: RulesConditionsAttrTypes}
	if len(src) == 0 {
		if prev.IsNull() || len(prev.Elements()) == 0 {
			return types.ListNull(elemType)
		}
		return prev
	}
	vals := make([]attr.Value, 0, len(src))
	for _, item := range src {
		obj, d := types.ObjectValueFrom(ctx, RulesConditionsAttrTypes, flattenRulesConditions(ctx, item, diags))
		diags.Append(d...)
		vals = append(vals, obj)
	}
	return types.ListValueMust(elemType, vals)
}

func objFromRulesConditionsHeaderConfig(ctx context.Context, m *RulesConditionsHeaderConfigModel) types.Object {
	if m == nil {
		return types.ObjectNull(RulesConditionsHeaderConfigAttrTypes)
	}
	v, d := types.ObjectValueFrom(ctx, RulesConditionsHeaderConfigAttrTypes, m)
	if d.HasError() {
		return types.ObjectNull(RulesConditionsHeaderConfigAttrTypes)
	}
	return v
}

func flattenRulesTags(ctx context.Context, src *albv20251030.TagInfo, diags *diag.Diagnostics) *RulesTagsModel {
	if src == nil {
		return nil
	}
	return &RulesTagsModel{
		TagKey:   fwHelper.StringValueOrPrev(src.TagKey, types.StringNull()),
		TagValue: fwHelper.StringValueOrPrev(src.TagValue, types.StringNull()),
	}
}

func flattenRulesTagsList(ctx context.Context, src []*albv20251030.TagInfo, prev types.List, diags *diag.Diagnostics) types.List {
	elemType := types.ObjectType{AttrTypes: RulesTagsAttrTypes}
	if len(src) == 0 {
		if prev.IsNull() || len(prev.Elements()) == 0 {
			return types.ListNull(elemType)
		}
		return prev
	}
	vals := make([]attr.Value, 0, len(src))
	for _, item := range src {
		obj, d := types.ObjectValueFrom(ctx, RulesTagsAttrTypes, flattenRulesTags(ctx, item, diags))
		diags.Append(d...)
		vals = append(vals, obj)
	}
	return types.ListValueMust(elemType, vals)
}

func flattenRules(ctx context.Context, src *albv20251030.RuleInput, diags *diag.Diagnostics) *RulesModel {
	if src == nil {
		return nil
	}
	return &RulesModel{
		Actions:    flattenRulesActionsList(ctx, src.Actions, types.ListNull(types.ObjectType{AttrTypes: RulesActionsAttrTypes}), diags),
		Conditions: flattenRulesConditionsList(ctx, src.Conditions, types.ListNull(types.ObjectType{AttrTypes: RulesConditionsAttrTypes}), diags),
		Priority:   fwHelper.Int64ValueOrPrev(src.Priority, types.Int64Null()),
		Direction:  fwHelper.StringValueOrPrev(src.Direction, types.StringNull()),
		RuleName:   fwHelper.StringValueOrPrev(src.RuleName, types.StringNull()),
		Tags:       flattenRulesTagsList(ctx, src.Tags, types.ListNull(types.ObjectType{AttrTypes: RulesTagsAttrTypes}), diags),
	}
}

func flattenRulesList(ctx context.Context, src []*albv20251030.RuleInput, prev types.List, diags *diag.Diagnostics) types.List {
	elemType := types.ObjectType{AttrTypes: RulesAttrTypes}
	if len(src) == 0 {
		if prev.IsNull() || len(prev.Elements()) == 0 {
			return types.ListNull(elemType)
		}
		return prev
	}
	vals := make([]attr.Value, 0, len(src))
	for _, item := range src {
		obj, d := types.ObjectValueFrom(ctx, RulesAttrTypes, flattenRules(ctx, item, diags))
		diags.Append(d...)
		vals = append(vals, obj)
	}
	return types.ListValueMust(elemType, vals)
}

var RulesActionsFixedResponseConfigAttrTypes = map[string]attr.Type{
	"http_code":    types.Int64Type,
	"content":      types.StringType,
	"content_type": types.StringType,
}

var RulesActionsInsertHeaderConfigAttrTypes = map[string]attr.Type{
	"key":        types.StringType,
	"value":      types.StringType,
	"value_type": types.StringType,
}

var RulesActionsRedirectConfigAttrTypes = map[string]attr.Type{
	"http_code": types.Int64Type,
	"host":      types.StringType,
	"path":      types.StringType,
	"port":      types.StringType,
	"protocol":  types.StringType,
	"query":     types.StringType,
}

var RulesActionsRemoveHeaderConfigAttrTypes = map[string]attr.Type{
	"key": types.StringType,
}

var RulesActionsRewriteConfigAttrTypes = map[string]attr.Type{
	"host":  types.StringType,
	"path":  types.StringType,
	"query": types.StringType,
}

var RulesActionsTargetGroupConfigTargetGroupsAttrTypes = map[string]attr.Type{
	"target_group_id": types.StringType,
	"weight":          types.Int64Type,
}

var RulesActionsTargetGroupConfigTargetGroupStickySessionAttrTypes = map[string]attr.Type{
	"enabled": types.BoolType,
	"timeout": types.Int64Type,
}

var RulesActionsTargetGroupConfigAttrTypes = map[string]attr.Type{
	"target_groups":               types.ListType{ElemType: types.ObjectType{AttrTypes: RulesActionsTargetGroupConfigTargetGroupsAttrTypes}},
	"target_group_sticky_session": types.ObjectType{AttrTypes: RulesActionsTargetGroupConfigTargetGroupStickySessionAttrTypes},
}

var RulesActionsAttrTypes = map[string]attr.Type{
	"order":                 types.Int64Type,
	"type":                  types.StringType,
	"fixed_response_config": types.ObjectType{AttrTypes: RulesActionsFixedResponseConfigAttrTypes},
	"insert_header_config":  types.ObjectType{AttrTypes: RulesActionsInsertHeaderConfigAttrTypes},
	"redirect_config":       types.ObjectType{AttrTypes: RulesActionsRedirectConfigAttrTypes},
	"remove_header_config":  types.ObjectType{AttrTypes: RulesActionsRemoveHeaderConfigAttrTypes},
	"rewrite_config":        types.ObjectType{AttrTypes: RulesActionsRewriteConfigAttrTypes},
	"target_group_config":   types.ObjectType{AttrTypes: RulesActionsTargetGroupConfigAttrTypes},
}

var RulesConditionsCookieConfigAttrTypes = map[string]attr.Type{
	"key":   types.StringType,
	"value": types.StringType,
}

var RulesConditionsHeaderConfigAttrTypes = map[string]attr.Type{
	"key":    types.StringType,
	"values": types.ListType{ElemType: types.StringType},
}

var RulesConditionsQueryStringConfigAttrTypes = map[string]attr.Type{
	"key":   types.StringType,
	"value": types.StringType,
}

var RulesConditionsAttrTypes = map[string]attr.Type{
	"type":                types.StringType,
	"cookie_config":       types.ListType{ElemType: types.ObjectType{AttrTypes: RulesConditionsCookieConfigAttrTypes}},
	"header_config":       types.ObjectType{AttrTypes: RulesConditionsHeaderConfigAttrTypes},
	"host_config":         types.ListType{ElemType: types.StringType},
	"method_config":       types.ListType{ElemType: types.StringType},
	"path_config":         types.ListType{ElemType: types.StringType},
	"query_string_config": types.ListType{ElemType: types.ObjectType{AttrTypes: RulesConditionsQueryStringConfigAttrTypes}},
	"source_ip_config":    types.ListType{ElemType: types.StringType},
}

var RulesTagsAttrTypes = map[string]attr.Type{
	"tag_key":   types.StringType,
	"tag_value": types.StringType,
}

var RulesAttrTypes = map[string]attr.Type{
	"actions":    types.ListType{ElemType: types.ObjectType{AttrTypes: RulesActionsAttrTypes}},
	"conditions": types.ListType{ElemType: types.ObjectType{AttrTypes: RulesConditionsAttrTypes}},
	"priority":   types.Int64Type,
	"direction":  types.StringType,
	"rule_name":  types.StringType,
	"tags":       types.ListType{ElemType: types.ObjectType{AttrTypes: RulesTagsAttrTypes}},
}

// Create handles `terraform apply` for new resources.
func (r *AlbRulesResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan AlbRulesResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createRulesRequest := albv20251030.NewCreateRulesRequest()
	createRulesRequest.ListenerId = fwHelper.StringPointerFromValue(plan.ListenerId)
	createRulesRequest.LoadBalancerId = fwHelper.StringPointerFromValue(plan.LoadBalancerId)
	createRulesRequest.Rules = expandRulesList(ctx, plan.Rules, &resp.Diagnostics)
	createRulesRequest.ClientToken = fwHelper.StringPointerFromValue(plan.ClientToken)
	createRulesRequest.DryRun = fwHelper.BoolPointerFromValue(plan.DryRun)
	if resp.Diagnostics.HasError() {
		return
	}

	service := NewAlbService(r.Client())
	ruleIds, err := service.CreateRules(ctx, createRulesRequest)
	if err != nil {
		resp.Diagnostics.AddError("Error creating ALB rules", fmt.Sprintf("CreateRules failed: %v", err))
		return
	}
	if len(ruleIds) == 0 || ruleIds[0] == nil || *ruleIds[0] == "" {
		resp.Diagnostics.AddError("Error creating ALB rules", "CreateRules returned an empty rule id.")
		return
	}
	firstRuleId := *ruleIds[0]

	// The state id is the first created rule id (resource_id = RuleIds, a
	// list response; the first element is used as the single resource id).
	plan.Id = types.StringValue(firstRuleId)

	// Persist the identity before the read-back so a partial failure still
	// leaves a destroyable resource in state.
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	describeRulesRequest := albv20251030.NewDescribeRulesRequest()
	describeRulesRequest.ListenerId = fwHelper.StringPointerFromValue(plan.ListenerId)
	describeRulesRequest.LoadBalancerId = fwHelper.StringPointerFromValue(plan.LoadBalancerId)
	describeRulesRequest.RuleIds = []*string{fwHelper.String(firstRuleId)}
	rulesOut, err := service.DescribeRules(ctx, describeRulesRequest)
	if err != nil {
		resp.Diagnostics.AddError("Error reading ALB rules after create", err.Error())
		return
	}
	item := findAlbRule(rulesOut, firstRuleId)
	if item == nil {
		resp.Diagnostics.AddError("ALB rule not readable after create", fmt.Sprintf("rule %s was created but cannot be read back", plan.Id.ValueString()))
		return
	}

	var diags diag.Diagnostics
	plan.Rules = flattenRulesFromOutputs(ctx, []*albv20251030.RuleOutput{item}, plan.Rules, &diags)
	resp.Diagnostics.Append(diags...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Read handles `terraform refresh` and post-apply state refresh.
func (r *AlbRulesResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state AlbRulesResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ruleId := state.Id.ValueString()
	describeRulesRequest := albv20251030.NewDescribeRulesRequest()
	describeRulesRequest.ListenerId = fwHelper.StringPointerFromValue(state.ListenerId)
	describeRulesRequest.LoadBalancerId = fwHelper.StringPointerFromValue(state.LoadBalancerId)
	describeRulesRequest.RuleIds = []*string{fwHelper.String(ruleId)}

	service := NewAlbService(r.Client())
	rulesOut, err := service.DescribeRules(ctx, describeRulesRequest)
	if err != nil {
		resp.Diagnostics.AddError("Error reading ALB rules", err.Error())
		return
	}

	item := findAlbRule(rulesOut, ruleId)
	if item == nil {
		resp.Diagnostics.AddWarning("Resource not found", fmt.Sprintf("ALB rule %s not found, removing from state.", ruleId))
		resp.State.RemoveResource(ctx)
		return
	}

	var diags diag.Diagnostics
	state.Rules = flattenRulesFromOutputs(ctx, []*albv20251030.RuleOutput{item}, state.Rules, &diags)
	resp.Diagnostics.Append(diags...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update handles in-place modifications of the tracked rule.
func (r *AlbRulesResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state AlbRulesResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ruleId := state.Id.ValueString()

	// ModifyRulesAttributes carries the full rule attributes including its id;
	// the write API is called only when the configured rule attributes change.
	if !plan.Rules.Equal(state.Rules) {
		var ruleModels []RulesModel
		resp.Diagnostics.Append(plan.Rules.ElementsAs(ctx, &ruleModels, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		if len(ruleModels) == 0 {
			resp.Diagnostics.AddError("Error updating ALB rules", "Planned rules list is empty.")
			return
		}
		var diags diag.Diagnostics
		modifyRules := expandRulesList(ctx, plan.Rules, &diags)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() || len(modifyRules) == 0 {
			return
		}

		modifyRulesAttributesRequest := albv20251030.NewModifyRulesAttributesRequest()
		modifyRulesAttributesRequest.ListenerId = fwHelper.StringPointerFromValue(plan.ListenerId)
		modifyRulesAttributesRequest.LoadBalancerId = fwHelper.StringPointerFromValue(plan.LoadBalancerId)
		modifyRulesAttributesRequest.Rules = []*albv20251030.RuleModify{{
			Actions:    convertActionsToModify(modifyRules[0].Actions),
			Conditions: convertConditionsToModify(modifyRules[0].Conditions),
			Priority:   modifyRules[0].Priority,
			RuleId:     fwHelper.String(ruleId),
			RuleName:   modifyRules[0].RuleName,
		}}

		service := NewAlbService(r.Client())
		if err := service.ModifyRulesAttributes(ctx, modifyRulesAttributesRequest); err != nil {
			resp.Diagnostics.AddError("Error updating ALB rules", fmt.Sprintf("ModifyRulesAttributes failed: %v", err))
			return
		}
	}

	describeRulesRequest := albv20251030.NewDescribeRulesRequest()
	describeRulesRequest.ListenerId = fwHelper.StringPointerFromValue(state.ListenerId)
	describeRulesRequest.LoadBalancerId = fwHelper.StringPointerFromValue(state.LoadBalancerId)
	describeRulesRequest.RuleIds = []*string{fwHelper.String(ruleId)}
	service := NewAlbService(r.Client())
	rulesOut, err := service.DescribeRules(ctx, describeRulesRequest)
	if err != nil {
		resp.Diagnostics.AddError("Error reading ALB rules after update", err.Error())
		return
	}
	item := findAlbRule(rulesOut, ruleId)
	if item == nil {
		resp.Diagnostics.AddError("ALB rule not readable after update", fmt.Sprintf("rule %s was updated but cannot be read back", ruleId))
		return
	}

	var diags diag.Diagnostics
	plan.Id = state.Id
	plan.Rules = flattenRulesFromOutputs(ctx, []*albv20251030.RuleOutput{item}, plan.Rules, &diags)
	resp.Diagnostics.Append(diags...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete handles resource deletion.
func (r *AlbRulesResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state AlbRulesResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	deleteRulesRequest := albv20251030.NewDeleteRulesRequest()
	deleteRulesRequest.ListenerId = fwHelper.StringPointerFromValue(state.ListenerId)
	deleteRulesRequest.LoadBalancerId = fwHelper.StringPointerFromValue(state.LoadBalancerId)
	deleteRulesRequest.RuleIds = []*string{fwHelper.String(state.Id.ValueString())}
	deleteRulesRequest.DryRun = fwHelper.BoolPointerFromValue(state.DryRun)

	service := NewAlbService(r.Client())
	if err := service.DeleteRules(ctx, deleteRulesRequest); err != nil {
		resp.Diagnostics.AddError("Error deleting ALB rules", fmt.Sprintf("DeleteRules failed: %v", err))
		return
	}
}

// findAlbRule picks the rule with the given id from a describe list.
func findAlbRule(rules []*albv20251030.RuleOutput, ruleId string) *albv20251030.RuleOutput {
	for _, item := range rules {
		if item != nil && item.RuleId != nil && *item.RuleId == ruleId {
			return item
		}
	}
	return nil
}

// flattenRulesFromOutputs maps a describe list (RuleOutput) back to the
// RuleInput-shaped schema model. RuleOutput carries read-only fields
// (RuleId, Status, CreateTime, ModifyTime) that are not part of the schema.
func flattenRulesFromOutputs(ctx context.Context, src []*albv20251030.RuleOutput, prev types.List, diags *diag.Diagnostics) types.List {
	inputs := make([]*albv20251030.RuleInput, 0, len(src))
	for _, o := range src {
		if o == nil {
			continue
		}
		inputs = append(inputs, &albv20251030.RuleInput{
			Actions:    o.Actions,
			Conditions: o.Conditions,
			Priority:   o.Priority,
			Direction:  o.Direction,
			RuleName:   o.RuleName,
			Tags:       o.Tags,
		})
	}
	return flattenRulesList(ctx, inputs, prev, diags)
}

// convertActionsToModify converts RuleInput actions to the RuleModify shape.
// TODO(human): ModifyRulesAttributes uses RuleModify whose action/condition
// element types are the same RuleAction/RuleCondition structs; if the cloud
// API requires a different nested shape for updates, adjust these helpers.
func convertActionsToModify(actions []*albv20251030.RuleAction) []*albv20251030.RuleAction {
	return actions
}

// convertConditionsToModify converts RuleInput conditions to the RuleModify
// shape.
func convertConditionsToModify(conditions []*albv20251030.RuleCondition) []*albv20251030.RuleCondition {
	return conditions
}
