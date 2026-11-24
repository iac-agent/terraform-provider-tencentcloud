package alb

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	albv20251030 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/alb/v20251030"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/fw"
	fwHelper "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/helper"
)

// Compile-time interface assertions. The base resource.Resource interface is
// checked implicitly by the factory return value.
var (
	_ resource.ResourceWithConfigure   = &AlbRulesResource{}
	_ resource.ResourceWithImportState = &AlbRulesResource{}
)

// NewAlbRulesResource is the factory referenced by
// tencentcloud/framework/registry.go to register this resource.
func NewAlbRulesResource() resource.Resource {
	return &AlbRulesResource{}
}

// AlbRulesResource implements resource.Resource for
// tencentcloud_alb_rules.
//
// Each state instance owns the set of forwarding rules created in one apply
// on a fixed load balancer listener. The listener binding is immutable, so
// load_balancer_id and listener_id force replacement; every other attribute
// is updated in place through ModifyRulesAttributes.
type AlbRulesResource struct {
	fw.ResourceWithConfigure
	fw.WithImportByID
}

// AlbRulesFixedResponseConfigModel maps the actions.fixed_response_config
// nested attribute.
type AlbRulesFixedResponseConfigModel struct {
	HttpCode    types.Int64  `tfsdk:"http_code"`
	Content     types.String `tfsdk:"content"`
	ContentType types.String `tfsdk:"content_type"`
}

// AlbRulesInsertHeaderConfigModel maps the actions.insert_header_config
// nested attribute.
type AlbRulesInsertHeaderConfigModel struct {
	Key       types.String `tfsdk:"key"`
	Value     types.String `tfsdk:"value"`
	ValueType types.String `tfsdk:"value_type"`
}

// AlbRulesRedirectConfigModel maps the actions.redirect_config nested
// attribute.
type AlbRulesRedirectConfigModel struct {
	HttpCode types.Int64  `tfsdk:"http_code"`
	Host     types.String `tfsdk:"host"`
	Path     types.String `tfsdk:"path"`
	Port     types.String `tfsdk:"port"`
	Protocol types.String `tfsdk:"protocol"`
	Query    types.String `tfsdk:"query"`
}

// AlbRulesRemoveHeaderConfigModel maps the actions.remove_header_config
// nested attribute.
type AlbRulesRemoveHeaderConfigModel struct {
	Key types.String `tfsdk:"key"`
}

// AlbRulesRewriteConfigModel maps the actions.rewrite_config nested
// attribute.
type AlbRulesRewriteConfigModel struct {
	Host  types.String `tfsdk:"host"`
	Path  types.String `tfsdk:"path"`
	Query types.String `tfsdk:"query"`
}

// AlbRulesTargetGroupModel maps the
// actions.target_group_config.target_groups nested attribute.
type AlbRulesTargetGroupModel struct {
	TargetGroupId types.String `tfsdk:"target_group_id"`
	Weight        types.Int64  `tfsdk:"weight"`
}

// AlbRulesTargetGroupStickySessionModel maps the
// actions.target_group_config.target_group_sticky_session nested attribute.
type AlbRulesTargetGroupStickySessionModel struct {
	Enabled types.Bool  `tfsdk:"enabled"`
	Timeout types.Int64 `tfsdk:"timeout"`
}

// AlbRulesTargetGroupConfigModel maps the actions.target_group_config
// nested attribute.
type AlbRulesTargetGroupConfigModel struct {
	TargetGroups             types.List   `tfsdk:"target_groups"`
	TargetGroupStickySession types.Object `tfsdk:"target_group_sticky_session"`
}

// AlbRulesActionModel maps the rules.actions nested attribute.
type AlbRulesActionModel struct {
	Order               types.Int64  `tfsdk:"order"`
	Type                types.String `tfsdk:"type"`
	FixedResponseConfig types.Object `tfsdk:"fixed_response_config"`
	InsertHeaderConfig  types.Object `tfsdk:"insert_header_config"`
	RedirectConfig      types.Object `tfsdk:"redirect_config"`
	RemoveHeaderConfig  types.Object `tfsdk:"remove_header_config"`
	RewriteConfig       types.Object `tfsdk:"rewrite_config"`
	TargetGroupConfig   types.Object `tfsdk:"target_group_config"`
}

// AlbRulesCookieConfigModel maps the conditions.cookie_config nested
// attribute.
type AlbRulesCookieConfigModel struct {
	Key   types.String `tfsdk:"key"`
	Value types.String `tfsdk:"value"`
}

// AlbRulesHeaderConfigModel maps the conditions.header_config nested
// attribute.
type AlbRulesHeaderConfigModel struct {
	Key    types.String `tfsdk:"key"`
	Values types.Set    `tfsdk:"values"`
}

// AlbRulesQueryStringConfigModel maps the conditions.query_string_config
// nested attribute.
type AlbRulesQueryStringConfigModel struct {
	Key   types.String `tfsdk:"key"`
	Value types.String `tfsdk:"value"`
}

// AlbRulesConditionModel maps the rules.conditions nested attribute.
type AlbRulesConditionModel struct {
	Type              types.String `tfsdk:"type"`
	CookieConfig      types.List   `tfsdk:"cookie_config"`
	HeaderConfig      types.Object `tfsdk:"header_config"`
	HostConfig        types.Set    `tfsdk:"host_config"`
	MethodConfig      types.Set    `tfsdk:"method_config"`
	PathConfig        types.Set    `tfsdk:"path_config"`
	QueryStringConfig types.List   `tfsdk:"query_string_config"`
	SourceIpConfig    types.Set    `tfsdk:"source_ip_config"`
}

// AlbRulesTagModel maps the rules.tags nested attribute.
type AlbRulesTagModel struct {
	TagKey   types.String `tfsdk:"tag_key"`
	TagValue types.String `tfsdk:"tag_value"`
}

// AlbRulesRuleModel maps the rules nested attribute.
type AlbRulesRuleModel struct {
	Actions    types.List   `tfsdk:"actions"`
	Conditions types.List   `tfsdk:"conditions"`
	Priority   types.Int64  `tfsdk:"priority"`
	Direction  types.String `tfsdk:"direction"`
	RuleName   types.String `tfsdk:"rule_name"`
	Tags       types.List   `tfsdk:"tags"`
	RuleId     types.String `tfsdk:"rule_id"`
	Status     types.String `tfsdk:"status"`
}

// AlbRulesResourceModel maps every top-level schema attribute; each one
// appears here exactly once.
type AlbRulesResourceModel struct {
	Id             types.String `tfsdk:"id"`
	LoadBalancerId types.String `tfsdk:"load_balancer_id"`
	ListenerId     types.String `tfsdk:"listener_id"`
	Rules          types.List   `tfsdk:"rules"`
	NextToken      types.String `tfsdk:"next_token"`
}

// albRulesRuleIdAttrType is the single attribute of the rule_id object; it
// is declared once so the schema-facing shape never drifts.
var albRulesRuleIdAttrType = map[string]attr.Type{
	"rule_id": types.StringType,
}

// albRulesFixedResponseConfigAttrTypes returns the attribute types of the
// actions.fixed_response_config object.
func albRulesFixedResponseConfigAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"http_code":    types.Int64Type,
		"content":      types.StringType,
		"content_type": types.StringType,
	}
}

// albRulesInsertHeaderConfigAttrTypes returns the attribute types of the
// actions.insert_header_config object.
func albRulesInsertHeaderConfigAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"key":        types.StringType,
		"value":      types.StringType,
		"value_type": types.StringType,
	}
}

// albRulesRedirectConfigAttrTypes returns the attribute types of the
// actions.redirect_config object.
func albRulesRedirectConfigAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"http_code": types.Int64Type,
		"host":      types.StringType,
		"path":      types.StringType,
		"port":      types.StringType,
		"protocol":  types.StringType,
		"query":     types.StringType,
	}
}

// albRulesRemoveHeaderConfigAttrTypes returns the attribute types of the
// actions.remove_header_config object.
func albRulesRemoveHeaderConfigAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"key": types.StringType,
	}
}

// albRulesRewriteConfigAttrTypes returns the attribute types of the
// actions.rewrite_config object.
func albRulesRewriteConfigAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"host":  types.StringType,
		"path":  types.StringType,
		"query": types.StringType,
	}
}

// albRulesTargetGroupAttrTypes returns the attribute types of the
// actions.target_group_config.target_groups elements.
func albRulesTargetGroupAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"target_group_id": types.StringType,
		"weight":          types.Int64Type,
	}
}

// albRulesTargetGroupStickySessionAttrTypes returns the attribute types of
// the actions.target_group_config.target_group_sticky_session object.
func albRulesTargetGroupStickySessionAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"enabled": types.BoolType,
		"timeout": types.Int64Type,
	}
}

// albRulesTargetGroupConfigAttrTypes returns the attribute types of the
// actions.target_group_config object.
func albRulesTargetGroupConfigAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"target_groups":               types.ListType{ElemType: types.ObjectType{AttrTypes: albRulesTargetGroupAttrTypes()}},
		"target_group_sticky_session": types.ObjectType{AttrTypes: albRulesTargetGroupStickySessionAttrTypes()},
	}
}

// albRulesActionAttrTypes returns the attribute types of the rules.actions
// elements.
func albRulesActionAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"order":                 types.Int64Type,
		"type":                  types.StringType,
		"fixed_response_config": types.ObjectType{AttrTypes: albRulesFixedResponseConfigAttrTypes()},
		"insert_header_config":  types.ObjectType{AttrTypes: albRulesInsertHeaderConfigAttrTypes()},
		"redirect_config":       types.ObjectType{AttrTypes: albRulesRedirectConfigAttrTypes()},
		"remove_header_config":  types.ObjectType{AttrTypes: albRulesRemoveHeaderConfigAttrTypes()},
		"rewrite_config":        types.ObjectType{AttrTypes: albRulesRewriteConfigAttrTypes()},
		"target_group_config":   types.ObjectType{AttrTypes: albRulesTargetGroupConfigAttrTypes()},
	}
}

// albRulesCookieConfigAttrTypes returns the attribute types of the
// conditions.cookie_config elements.
func albRulesCookieConfigAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"key":   types.StringType,
		"value": types.StringType,
	}
}

// albRulesHeaderConfigAttrTypes returns the attribute types of the
// conditions.header_config object.
func albRulesHeaderConfigAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"key":    types.StringType,
		"values": types.SetType{ElemType: types.StringType},
	}
}

// albRulesQueryStringConfigAttrTypes returns the attribute types of the
// conditions.query_string_config elements.
func albRulesQueryStringConfigAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"key":   types.StringType,
		"value": types.StringType,
	}
}

// albRulesConditionAttrTypes returns the attribute types of the
// rules.conditions elements.
func albRulesConditionAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"type":                types.StringType,
		"cookie_config":       types.ListType{ElemType: types.ObjectType{AttrTypes: albRulesCookieConfigAttrTypes()}},
		"header_config":       types.ObjectType{AttrTypes: albRulesHeaderConfigAttrTypes()},
		"host_config":         types.SetType{ElemType: types.StringType},
		"method_config":       types.SetType{ElemType: types.StringType},
		"path_config":         types.SetType{ElemType: types.StringType},
		"query_string_config": types.ListType{ElemType: types.ObjectType{AttrTypes: albRulesQueryStringConfigAttrTypes()}},
		"source_ip_config":    types.SetType{ElemType: types.StringType},
	}
}

// albRulesTagAttrTypes returns the attribute types of the rules.tags
// elements.
func albRulesTagAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"tag_key":   types.StringType,
		"tag_value": types.StringType,
	}
}

// albRulesAttrTypes returns the attribute types of the rules elements.
func albRulesAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"actions":    types.ListType{ElemType: types.ObjectType{AttrTypes: albRulesActionAttrTypes()}},
		"conditions": types.ListType{ElemType: types.ObjectType{AttrTypes: albRulesConditionAttrTypes()}},
		"priority":   types.Int64Type,
		"direction":  types.StringType,
		"rule_name":  types.StringType,
		"tags":       types.ListType{ElemType: types.ObjectType{AttrTypes: albRulesTagAttrTypes()}},
		"rule_id":    types.StringType,
		"status":     types.StringType,
	}
}

// Metadata reports the Terraform type name.
func (r *AlbRulesResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "tencentcloud_alb_rules"
}

// Schema declares the ALB forwarding rules resource.
func (r *AlbRulesResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a resource to create forwarding rules for an ALB (Application Load Balancer) listener.",
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
				Description: "ID of the load balancer instance, in the format of `alb-` followed by 8 alphanumeric characters.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"listener_id": schema.StringAttribute{
				Required:    true,
				Description: "ID of the listener, in the format of `lst-` followed by 8 alphanumeric characters.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"rules": schema.ListNestedAttribute{
				Required:    true,
				Description: "List of the forwarding rules.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"actions": schema.ListNestedAttribute{
							Required:    true,
							Description: "List of the forwarding rule actions. The actions must contain exactly one of `TargetGroup`, `Redirect` or `FixedResponse`, and it must be placed last by execution order.",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"order": schema.Int64Attribute{
										Required:    true,
										Description: "Execution order of the forwarding action. It cannot be duplicated and actions are executed in ascending order. Value range: 1 ~ 50000.",
									},
									"type": schema.StringAttribute{
										Required:    true,
										Description: "Type of the forwarding action. Valid values: `TargetGroup` (forward to target group), `Redirect` (redirect), `FixedResponse` (return fixed content), `Rewrite` (rewrite), `InsertHeader` (insert HTTP header), `RemoveHeader` (remove HTTP header). The actions must contain exactly one of `TargetGroup`, `Redirect` or `FixedResponse`, and it must be placed last by execution order.",
									},
									"fixed_response_config": schema.SingleNestedAttribute{
										Optional:    true,
										Description: "Configuration of the fixed response. It is required when `type` is `FixedResponse`.",
										Attributes: map[string]schema.Attribute{
											"http_code": schema.Int64Attribute{
												Required:    true,
												Description: "HTTP response code to return. Supported values: 2xx, 4xx, 5xx.",
											},
											"content": schema.StringAttribute{
												Optional:    true,
												Description: "Fixed content to return. Only ASCII characters are supported, with a maximum size of 1KB.",
											},
											"content_type": schema.StringAttribute{
												Optional:    true,
												Description: "Format of the fixed content. Valid values: `text/plain`, `text/css`, `text/html`, `application/javascript`, `application/json`.",
											},
										},
									},
									"insert_header_config": schema.SingleNestedAttribute{
										Optional:    true,
										Description: "Configuration of the inserted HTTP header. It is required when `type` is `InsertHeader`.",
										Attributes: map[string]schema.Attribute{
											"key": schema.StringAttribute{
												Required:    true,
												Description: "Key of the inserted HTTP header, 1 ~ 40 characters long, with the supported character set of `a-z A-Z 0-9 - _`. Chinese characters and the headers `Cookie`, `Host`, `Content-Length`, `Connection`, `Upgrade`, `transfer-encoding`, `keep-alive`, `te`, `authority`, `x-forwarded-for`, `x-forwarded-proto`, `x-forwarded-host`, `x-forwarded-port` are not supported.",
											},
											"value": schema.StringAttribute{
												Required:    true,
												Description: "Value of the inserted HTTP header. When `value_type` is `SystemDefined`, valid values: `ClientPort` (client port), `ClientIp` (client IP address), `Protocol` (protocol of the client request), `CLBPort` (listener port of the load balancer instance). When `value_type` is `UserDefined` or `ReferenceHeader`, it is 1 ~ 128 printable characters long, without `\"`, without leading or trailing spaces, and must not end with `\\`.",
											},
											"value_type": schema.StringAttribute{
												Required:    true,
												Description: "Type of the HTTP header value. Valid values: `SystemDefined` (system-defined header), `UserDefined` (user-defined header), `ReferenceHeader` (reference one header of the request).",
											},
										},
									},
									"redirect_config": schema.SingleNestedAttribute{
										Optional:    true,
										Description: "Configuration of the redirect. Except `http_code`, the other configurations cannot all use the default values. It is required when `type` is `Redirect`.",
										Attributes: map[string]schema.Attribute{
											"http_code": schema.Int64Attribute{
												Required:    true,
												Description: "HTTP code of the redirect. Supported values: 301, 302, 303, 307, 308.",
											},
											"host": schema.StringAttribute{
												Optional:    true,
												Description: "Host of the redirect, default value `${host}`. It is 3 ~ 128 characters long, with the supported character set of `a-z 0-9 _ . -`.",
											},
											"path": schema.StringAttribute{
												Optional:    true,
												Description: "Path of the redirect, default value `${path}`. It is 1 ~ 128 characters long, with the supported character set of `a-z A-Z 0-9 ? = _ . - / :`.",
											},
											"port": schema.StringAttribute{
												Optional:    true,
												Description: "Port of the redirect, default value `${port}`. Value range: 1 ~ 65535.",
											},
											"protocol": schema.StringAttribute{
												Optional:    true,
												Description: "Protocol of the redirect, default value `${protocol}`. Valid values: `HTTP`, `HTTPS`.",
											},
											"query": schema.StringAttribute{
												Optional:    true,
												Description: "Query string of the redirect, default value `${query}`. It is 1 ~ 128 characters long, supports printable characters, and does not support `#`, `[`, `]`, `{`, `}`, `|`, `<`, `>`, `&` and spaces.",
											},
										},
									},
									"remove_header_config": schema.SingleNestedAttribute{
										Optional:    true,
										Description: "Configuration of the removed HTTP header. It is required when `type` is `RemoveHeader`.",
										Attributes: map[string]schema.Attribute{
											"key": schema.StringAttribute{
												Required:    true,
												Description: "Key of the removed HTTP header, 1 ~ 40 characters long, with the supported character set of `a-z A-Z 0-9 - _`. The headers `Cookie`, `Host`, `Content-Length`, `Connection`, `Upgrade`, `transfer-encoding`, `keep-alive`, `te`, `authority`, `x-forwarded-for`, `x-forwarded-proto`, `x-forwarded-host`, `x-forwarded-port`, `server` are not supported.",
											},
										},
									},
									"rewrite_config": schema.SingleNestedAttribute{
										Optional:    true,
										Description: "Configuration of the rewrite. It is required when `type` is `Rewrite`.",
										Attributes: map[string]schema.Attribute{
											"host": schema.StringAttribute{
												Optional:    true,
												Description: "Host of the rewrite, default value `${host}`. It is 3 ~ 128 characters long, with the supported character set of `a-z 0-9 _ . -`.",
											},
											"path": schema.StringAttribute{
												Optional:    true,
												Description: "Path of the rewrite, default value `${path}`. It is 1 ~ 128 characters long, with the supported character set of `a-z A-Z 0-9 ? = _ . - / :`.",
											},
											"query": schema.StringAttribute{
												Optional:    true,
												Description: "Query string of the rewrite, default value `${query}`. It is 1 ~ 128 characters long, supports printable characters, and does not support `#`, `[`, `]`, `{`, `}`, `|`, `<`, `>`, `&` and spaces.",
											},
										},
									},
									"target_group_config": schema.SingleNestedAttribute{
										Optional:    true,
										Description: "Configuration of the target group forwarding. It is required when `type` is `TargetGroup`.",
										Attributes: map[string]schema.Attribute{
											"target_groups": schema.ListNestedAttribute{
												Required:    true,
												Description: "List of the target groups.",
												NestedObject: schema.NestedAttributeObject{
													Attributes: map[string]schema.Attribute{
														"target_group_id": schema.StringAttribute{
															Required:    true,
															Description: "ID of the target group, in the format of `lbtg-` followed by 8 alphanumeric characters.",
														},
														"weight": schema.Int64Attribute{
															Optional:    true,
															Description: "Weight of the target group. Value range: 0 ~ 100, default value: 10.",
														},
													},
												},
											},
											"target_group_sticky_session": schema.SingleNestedAttribute{
												Optional:    true,
												Description: "Sticky session configuration between the target groups.",
												Attributes: map[string]schema.Attribute{
													"enabled": schema.BoolAttribute{
														Required:    true,
														Description: "Whether to enable sticky session, disabled by default.",
													},
													"timeout": schema.Int64Attribute{
														Optional:    true,
														Description: "Timeout of the sticky session, in seconds. Value range: 1 ~ 86400, default value: 1000.",
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
							Description: "List of the forwarding rule conditions.",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"type": schema.StringAttribute{
										Required:    true,
										Description: "Type of the forwarding condition. Valid values: `Host` (host), `Path` (path), `Header` (HTTP header field), `QueryString` (HTTP query string), `Method` (request method), `Cookie` (cookie), `SourceIp` (source IP).",
									},
									"cookie_config": schema.ListNestedAttribute{
										Optional:    true,
										Description: "Configuration of the cookie. It is required when `type` is `Cookie`.",
										NestedObject: schema.NestedAttributeObject{
											Attributes: map[string]schema.Attribute{
												"key": schema.StringAttribute{
													Required:    true,
													Description: "Key of the cookie, 1 ~ 64 characters long, supporting letters, digits and underscores.",
												},
												"value": schema.StringAttribute{
													Required:    true,
													Description: "Value of the cookie, 1 ~ 128 characters long, supporting printable characters.",
												},
											},
										},
									},
									"header_config": schema.SingleNestedAttribute{
										Optional:    true,
										Description: "Configuration of the HTTP header. It is required when `type` is `Header`.",
										Attributes: map[string]schema.Attribute{
											"key": schema.StringAttribute{
												Required:    true,
												Description: "Key of the HTTP header, 1 ~ 40 characters long, with the supported character set of `a-z A-Z 0-9 - _`. Chinese characters and the headers `Host`, `Cookie` are not supported.",
											},
											"values": schema.SetAttribute{
												Required:    true,
												ElementType: types.StringType,
												Description: "Values of the HTTP header, each 1 ~ 128 characters long, supporting printable characters. `\"` is not supported, leading and trailing spaces are not supported, and `\\` is not supported at the end.",
											},
										},
									},
									"host_config": schema.SetAttribute{
										Optional:    true,
										ElementType: types.StringType,
										Description: "Host names. The host configuration can appear only once in a rule, each 3 ~ 128 characters long, supporting exact matching, regular expression matching and wildcard matching. It cannot start or end with a period (`.`) or an underscore (`_`). Exact matching supports the character set of `a-z 0-9 . - _`. Regular expression matching starts with a tilde (`~`) and supports the character set of `a-z 0-9 . - ? = ~ _ - + \\ ^ * ! $ & | ( ) [ ]`. Wildcard matching uses an asterisk (`*`) for multi-character wildcards and a question mark (`?`) for single-character wildcards, supporting the character set of `a-z 0-9 . - _ * ?`.",
									},
									"method_config": schema.SetAttribute{
										Optional:    true,
										ElementType: types.StringType,
										Description: "Request methods. Valid values: `HEAD`, `GET`, `POST`, `OPTIONS`, `PUT`, `PATCH`, `DELETE`.",
									},
									"path_config": schema.SetAttribute{
										Optional:    true,
										ElementType: types.StringType,
										Description: "Forwarding paths, each 1 ~ 128 characters long, supporting exact matching, regular expression matching and wildcard matching. Exact matching supports the character set of `a-z A-Z 0-9 . - _ / = :`. Regular expression matching starts with `~` for case-sensitive matching or `~*` for case-insensitive matching, supporting the character set of `a-z A-Z 0-9 . - _ / = ? ~ ^ * $ : ( ) [ ] + |`. Wildcard matching uses an asterisk (`*`) for multi-character wildcards and a question mark (`?`) for single-character wildcards, supporting the character set of `a-z A-Z 0-9 . - _ / = :`.",
									},
									"query_string_config": schema.ListNestedAttribute{
										Optional:    true,
										Description: "Configuration of the query string. It is required when `type` is `QueryString`.",
										NestedObject: schema.NestedAttributeObject{
											Attributes: map[string]schema.Attribute{
												"key": schema.StringAttribute{
													Required:    true,
													Description: "Key of the query string, 1 ~ 16 characters long, supporting printable characters, and does not support spaces and `#`, `[`, `]`, `{`, `}`, `\\`, `|`, `<`, `>`, `&`. Multi-character wildcard `*` and single-character wildcard `?` are supported.",
												},
												"value": schema.StringAttribute{
													Required:    true,
													Description: "Value of the query string, 1 ~ 128 characters long, supporting printable characters, and does not support spaces and `#`, `[`, `]`, `{`, `}`, `\\`, `|`, `<`, `>`, `&`. Multi-character wildcard `*` and single-character wildcard `?` are supported.",
												},
											},
										},
									},
									"source_ip_config": schema.SetAttribute{
										Optional:    true,
										ElementType: types.StringType,
										Description: "Source IP matching configuration. CIDR format, such as `x.x.x.x/32` for an IP address and `x.x.x.x/24` for an IP range.",
									},
								},
							},
						},
						"priority": schema.Int64Attribute{
							Required:    true,
							Description: "Priority of the forwarding rule. A smaller value indicates a higher priority. It cannot be duplicated. Value range: 1 ~ 10000.",
						},
						"direction": schema.StringAttribute{
							Optional:    true,
							Computed:    true,
							Description: "Direction of the forwarding rule. Valid values: `Request` (request direction from the client to the load balancer), `Response` (response direction from the backend server to the load balancer). Default value: `Request`.",
						},
						"rule_name": schema.StringAttribute{
							Optional:    true,
							Computed:    true,
							Description: "Name of the forwarding rule. It is 1 ~ 255 characters long and can contain digits, upper and lower case letters, Chinese characters, periods (`.`), underscores (`_`) and hyphens (`-`).",
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
						"rule_id": schema.StringAttribute{
							Computed:    true,
							Description: "ID of the forwarding rule, in the format of `rule-` followed by 8 alphanumeric characters.",
						},
						"status": schema.StringAttribute{
							Computed:    true,
							Description: "Status of the forwarding rule. Valid values: `Provisioning` (creating), `Active` (running), `Configuring` (configuring).",
						},
					},
				},
			},
			"next_token": schema.StringAttribute{
				Computed:    true,
				Description: "Token of the next query. It is empty when the current page is the last one.",
			},
		},
	}
}

// Create handles `terraform apply` for new resources. The whole planned rule
// set is created in one CreateRules call; each created rule is then read
// back with DescribeRules so rule_id and status are persisted.
func (r *AlbRulesResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan AlbRulesResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Warn instead of error when the client is missing; keep the guard
	// exactly like the framework templates do.
	if r.Client() == nil {
		resp.Diagnostics.AddError(
			"Provider not configured",
			"The provider client is not available. This is unexpected; please report it.",
		)

		return
	}

	service := NewAlbService(r.Client())

	loadBalancerId := plan.LoadBalancerId.ValueString()
	listenerId := plan.ListenerId.ValueString()

	createRulesRequest := albv20251030.NewCreateRulesRequest()
	createRulesRequest.LoadBalancerId = albSdkString(loadBalancerId)
	createRulesRequest.ListenerId = albSdkString(listenerId)

	ruleInputs, diags := expandAlbRules(ctx, plan.Rules)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	createRulesRequest.Rules = ruleInputs

	ruleIds, err := service.CreateRules(ctx, createRulesRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating ALB listener rules",
			fmt.Sprintf("CreateRules failed: %v", err),
		)

		return
	}

	if len(ruleIds) != len(ruleInputs) {
		resp.Diagnostics.AddError(
			"Error creating ALB listener rules",
			fmt.Sprintf("CreateRules returned %d rule ids for %d rules.", len(ruleIds), len(ruleInputs)),
		)

		return
	}

	plan.Id = types.StringValue(fmt.Sprintf("%s#%s", loadBalancerId, listenerId))

	// Persist the identity before reading back so a partial failure still
	// leaves a destroyable resource in state.
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Read-back: same request form as Read; the caller picks its own items
	// from the list.
	describeRulesRequest := albv20251030.NewDescribeRulesRequest()
	describeRulesRequest.LoadBalancerId = albSdkString(loadBalancerId)
	describeRulesRequest.ListenerId = albSdkString(listenerId)
	describeRulesRequest.RuleIds = albRuleIdPointers(ruleIds)

	rules, _, err := service.DescribeRules(ctx, describeRulesRequest)
	if err != nil {
		resp.Diagnostics.AddError("Error reading ALB listener rules after create", err.Error())
		return
	}

	refreshAlbRulesResourceModel(ctx, &plan, rules, &resp.Diagnostics)
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

	if r.Client() == nil {
		resp.Diagnostics.AddError(
			"Provider not configured",
			"The provider client is not available. This is unexpected; please report it.",
		)

		return
	}

	service := NewAlbService(r.Client())

	loadBalancerId := state.LoadBalancerId.ValueString()
	listenerId := state.ListenerId.ValueString()

	describeRulesRequest := albv20251030.NewDescribeRulesRequest()
	describeRulesRequest.LoadBalancerId = albSdkString(loadBalancerId)
	describeRulesRequest.ListenerId = albSdkString(listenerId)

	if !state.Rules.IsNull() && !state.Rules.IsUnknown() {
		var ruleModels []AlbRulesRuleModel
		resp.Diagnostics.Append(state.Rules.ElementsAs(ctx, &ruleModels, false)...)
		if resp.Diagnostics.HasError() {
			return
		}

		ruleIds := make([]string, 0, len(ruleModels))
		for _, ruleModel := range ruleModels {
			if !ruleModel.RuleId.IsNull() && !ruleModel.RuleId.IsUnknown() && ruleModel.RuleId.ValueString() != "" {
				ruleIds = append(ruleIds, ruleModel.RuleId.ValueString())
			}
		}
		if len(ruleIds) > 0 {
			describeRulesRequest.RuleIds = albRuleIdPointers(ruleIds)
		}
	}

	rules, nextToken, err := service.DescribeRules(ctx, describeRulesRequest)
	if err != nil {
		resp.Diagnostics.AddError("Error reading ALB listener rules", err.Error())
		return
	}

	// Keep only the rules this state instance owns so a partial removal is
	// treated as not-found.
	if len(rules) == 0 {
		resp.Diagnostics.AddWarning(
			"Resource not found",
			fmt.Sprintf("ALB listener rules %s not found, removing from state.", state.Id.ValueString()),
		)

		resp.State.RemoveResource(ctx)
		return
	}

	state.NextToken = fwHelper.StringValueOrNull(nextToken)

	refreshAlbRulesResourceModel(ctx, &state, rules, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update handles in-place modifications through ModifyRulesAttributes. The
// listener binding is immutable (it forces replacement in the schema).
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

	if r.Client() == nil {
		resp.Diagnostics.AddError(
			"Provider not configured",
			"The provider client is not available. This is unexpected; please report it.",
		)

		return
	}

	service := NewAlbService(r.Client())

	loadBalancerId := plan.LoadBalancerId.ValueString()
	listenerId := plan.ListenerId.ValueString()

	modifyRulesAttributesRequest := albv20251030.NewModifyRulesAttributesRequest()
	modifyRulesAttributesRequest.LoadBalancerId = albSdkString(loadBalancerId)
	modifyRulesAttributesRequest.ListenerId = albSdkString(listenerId)

	ruleModifies, diags := expandAlbRuleModifies(ctx, plan.Rules, state.Rules)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	modifyRulesAttributesRequest.Rules = ruleModifies

	if err := service.ModifyRulesAttributes(ctx, modifyRulesAttributesRequest); err != nil {
		resp.Diagnostics.AddError(
			"Error updating ALB listener rules",
			fmt.Sprintf("ModifyRulesAttributes failed: %v", err),
		)

		return
	}

	describeRulesRequest := albv20251030.NewDescribeRulesRequest()
	describeRulesRequest.LoadBalancerId = albSdkString(loadBalancerId)
	describeRulesRequest.ListenerId = albSdkString(listenerId)
	describeRulesRequest.RuleIds = albRuleIdPointers(albRuleIdsFromModifies(ruleModifies))

	rules, _, err := service.DescribeRules(ctx, describeRulesRequest)
	if err != nil {
		resp.Diagnostics.AddError("Error reading ALB listener rules after update", err.Error())
		return
	}

	plan.Id = state.Id
	refreshAlbRulesResourceModel(ctx, &plan, rules, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete handles resource deletion. Every rule id this state instance owns
// is deleted with one DeleteRules call.
func (r *AlbRulesResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state AlbRulesResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if r.Client() == nil {
		resp.Diagnostics.AddError(
			"Provider not configured",
			"The provider client is not available. This is unexpected; please report it.",
		)

		return
	}

	service := NewAlbService(r.Client())

	loadBalancerId := state.LoadBalancerId.ValueString()
	listenerId := state.ListenerId.ValueString()

	ruleIds := albStateRuleIds(state.Rules, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	deleteRulesRequest := albv20251030.NewDeleteRulesRequest()
	deleteRulesRequest.LoadBalancerId = albSdkString(loadBalancerId)
	deleteRulesRequest.ListenerId = albSdkString(listenerId)
	deleteRulesRequest.RuleIds = albRuleIdPointers(ruleIds)

	if err := service.DeleteRules(ctx, deleteRulesRequest); err != nil {
		resp.Diagnostics.AddError(
			"Error deleting ALB listener rules",
			fmt.Sprintf("DeleteRules failed: %v", err),
		)

		return
	}
}

// refreshAlbRulesResourceModel copies cloud-owned fields into the model
// while keeping the user-supplied values for fields the cloud API does not
// echo back.
func refreshAlbRulesResourceModel(ctx context.Context, m *AlbRulesResourceModel, rules []*albv20251030.RuleOutput, diags *diag.Diagnostics) {
	ruleModels := make([]AlbRulesRuleModel, 0, len(rules))
	for _, rule := range rules {
		ruleModel := AlbRulesRuleModel{}

		actions := make([]AlbRulesActionModel, 0, len(rule.Actions))
		for _, action := range rule.Actions {
			actionModel := AlbRulesActionModel{
				Order: albInt64PointerValue(action.Order),
				Type:  fwHelper.StringValueOrPrev(action.Type, types.StringNull()),
			}

			if action.FixedResponseConfig != nil {
				actionModel.FixedResponseConfig = albFixedResponseConfigObject(action.FixedResponseConfig)
			} else {
				actionModel.FixedResponseConfig = types.ObjectNull(albRulesFixedResponseConfigAttrTypes())
			}

			if action.InsertHeaderConfig != nil {
				actionModel.InsertHeaderConfig = albInsertHeaderConfigObject(action.InsertHeaderConfig)
			} else {
				actionModel.InsertHeaderConfig = types.ObjectNull(albRulesInsertHeaderConfigAttrTypes())
			}

			if action.RedirectConfig != nil {
				actionModel.RedirectConfig = albRedirectConfigObject(action.RedirectConfig)
			} else {
				actionModel.RedirectConfig = types.ObjectNull(albRulesRedirectConfigAttrTypes())
			}

			if action.RemoveHeaderConfig != nil {
				actionModel.RemoveHeaderConfig = albRemoveHeaderConfigObject(action.RemoveHeaderConfig)
			} else {
				actionModel.RemoveHeaderConfig = types.ObjectNull(albRulesRemoveHeaderConfigAttrTypes())
			}

			if action.RewriteConfig != nil {
				actionModel.RewriteConfig = albRewriteConfigObject(action.RewriteConfig)
			} else {
				actionModel.RewriteConfig = types.ObjectNull(albRulesRewriteConfigAttrTypes())
			}

			if action.TargetGroupConfig != nil {
				actionModel.TargetGroupConfig = albTargetGroupConfigObject(ctx, action.TargetGroupConfig, diags)
				if diags.HasError() {
					return
				}
			} else {
				actionModel.TargetGroupConfig = types.ObjectNull(albRulesTargetGroupConfigAttrTypes())
			}

			actions = append(actions, actionModel)
		}

		conditions := make([]AlbRulesConditionModel, 0, len(rule.Conditions))
		for _, condition := range rule.Conditions {
			conditionModel := AlbRulesConditionModel{
				Type: fwHelper.StringValueOrPrev(condition.Type, types.StringNull()),
			}

			if len(condition.CookieConfig) > 0 {
				cookieConfigs := make([]AlbRulesCookieConfigModel, 0, len(condition.CookieConfig))
				for _, cookieConfig := range condition.CookieConfig {
					cookieConfigs = append(cookieConfigs, AlbRulesCookieConfigModel{
						Key:   fwHelper.StringValueOrPrev(cookieConfig.Key, types.StringNull()),
						Value: fwHelper.StringValueOrPrev(cookieConfig.Value, types.StringNull()),
					})
				}
				cookieConfigList, listDiags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: albRulesCookieConfigAttrTypes()}, cookieConfigs)
				diags.Append(listDiags...)
				if diags.HasError() {
					return
				}
				conditionModel.CookieConfig = cookieConfigList
			} else {
				conditionModel.CookieConfig = types.ListNull(types.ObjectType{AttrTypes: albRulesCookieConfigAttrTypes()})
			}

			if condition.HeaderConfig != nil {
				headerValues, setDiags := types.SetValueFrom(ctx, types.StringType, albStringPointers(condition.HeaderConfig.Values))
				diags.Append(setDiags...)
				if diags.HasError() {
					return
				}
				headerConfig, objectDiags := types.ObjectValue(albRulesHeaderConfigAttrTypes(), map[string]attr.Value{
					"key":    fwHelper.StringValueOrPrev(condition.HeaderConfig.Key, types.StringNull()),
					"values": headerValues,
				})
				diags.Append(objectDiags...)
				if diags.HasError() {
					return
				}
				conditionModel.HeaderConfig = headerConfig
			} else {
				conditionModel.HeaderConfig = types.ObjectNull(albRulesHeaderConfigAttrTypes())
			}

			conditionModel.HostConfig = albStringSetValue(ctx, condition.HostConfig, diags)
			conditionModel.MethodConfig = albStringSetValue(ctx, condition.MethodConfig, diags)
			conditionModel.PathConfig = albStringSetValue(ctx, condition.PathConfig, diags)
			conditionModel.SourceIpConfig = albStringSetValue(ctx, condition.SourceIpConfig, diags)
			if diags.HasError() {
				return
			}

			if len(condition.QueryStringConfig) > 0 {
				queryStringConfigs := make([]AlbRulesQueryStringConfigModel, 0, len(condition.QueryStringConfig))
				for _, queryStringConfig := range condition.QueryStringConfig {
					queryStringConfigs = append(queryStringConfigs, AlbRulesQueryStringConfigModel{
						Key:   fwHelper.StringValueOrPrev(queryStringConfig.Key, types.StringNull()),
						Value: fwHelper.StringValueOrPrev(queryStringConfig.Value, types.StringNull()),
					})
				}
				queryStringConfigList, listDiags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: albRulesQueryStringConfigAttrTypes()}, queryStringConfigs)
				diags.Append(listDiags...)
				if diags.HasError() {
					return
				}
				conditionModel.QueryStringConfig = queryStringConfigList
			} else {
				conditionModel.QueryStringConfig = types.ListNull(types.ObjectType{AttrTypes: albRulesQueryStringConfigAttrTypes()})
			}

			conditions = append(conditions, conditionModel)
		}

		tags := make([]AlbRulesTagModel, 0, len(rule.Tags))
		for _, tag := range rule.Tags {
			tags = append(tags, AlbRulesTagModel{
				TagKey:   fwHelper.StringValueOrPrev(tag.TagKey, types.StringNull()),
				TagValue: fwHelper.StringValueOrPrev(tag.TagValue, types.StringNull()),
			})
		}

		ruleModel.Actions = albListValueFrom(ctx, albRulesActionAttrTypes(), actions, diags)
		ruleModel.Conditions = albListValueFrom(ctx, albRulesConditionAttrTypes(), conditions, diags)
		ruleModel.Tags = albListValueFrom(ctx, albRulesTagAttrTypes(), tags, diags)
		if diags.HasError() {
			return
		}

		ruleModel.Priority = albInt64PointerValue(rule.Priority)
		ruleModel.Direction = fwHelper.StringValueOrPrev(rule.Direction, types.StringNull())
		ruleModel.RuleName = fwHelper.StringValueOrPrev(rule.RuleName, types.StringNull())
		ruleModel.RuleId = fwHelper.StringValueOrPrev(rule.RuleId, types.StringNull())
		ruleModel.Status = fwHelper.StringValueOrPrev(rule.Status, types.StringNull())

		ruleModels = append(ruleModels, ruleModel)
	}

	m.Rules = albListValueFrom(ctx, albRulesAttrTypes(), ruleModels, diags)
}

// albListValueFrom converts a model slice into a types.List; an empty slice
// becomes an empty list, never null.
func albListValueFrom(ctx context.Context, attrTypes map[string]attr.Type, models any, diags *diag.Diagnostics) types.List {
	if models == nil {
		return types.ListNull(types.ObjectType{AttrTypes: attrTypes})
	}

	list, listDiags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: attrTypes}, models)
	diags.Append(listDiags...)
	return list
}

// albStringSetValue converts nullable string pointers into a set of strings;
// a nil or empty slice becomes an empty set.
func albStringSetValue(ctx context.Context, values []*string, diags *diag.Diagnostics) types.Set {
	set, setDiags := types.SetValueFrom(ctx, types.StringType, albStringPointers(values))
	diags.Append(setDiags...)
	return set
}

// albStringPointers dereferences string pointers, skipping nils.
func albStringPointers(values []*string) []string {
	result := make([]string, 0, len(values))
	for _, value := range values {
		if value != nil {
			result = append(result, *value)
		}
	}
	return result
}

// albFixedResponseConfigObject converts the cloud-side fixed response
// configuration into the nested object value.
func albFixedResponseConfigObject(config *albv20251030.FixedResponseInfo) types.Object {
	object, _ := types.ObjectValue(albRulesFixedResponseConfigAttrTypes(), map[string]attr.Value{
		"http_code":    albInt64PointerValue(config.HttpCode),
		"content":      fwHelper.StringValueOrPrev(config.Content, types.StringNull()),
		"content_type": fwHelper.StringValueOrPrev(config.ContentType, types.StringNull()),
	})
	return object
}

// albInsertHeaderConfigObject converts the cloud-side insert header
// configuration into the nested object value.
func albInsertHeaderConfigObject(config *albv20251030.InsertHTTPHeaderInfo) types.Object {
	object, _ := types.ObjectValue(albRulesInsertHeaderConfigAttrTypes(), map[string]attr.Value{
		"key":        fwHelper.StringValueOrPrev(config.Key, types.StringNull()),
		"value":      fwHelper.StringValueOrPrev(config.Value, types.StringNull()),
		"value_type": fwHelper.StringValueOrPrev(config.ValueType, types.StringNull()),
	})
	return object
}

// albRedirectConfigObject converts the cloud-side redirect configuration
// into the nested object value.
func albRedirectConfigObject(config *albv20251030.HTTPRedirectInfo) types.Object {
	object, _ := types.ObjectValue(albRulesRedirectConfigAttrTypes(), map[string]attr.Value{
		"http_code": albInt64PointerValue(config.HttpCode),
		"host":      fwHelper.StringValueOrPrev(config.Host, types.StringNull()),
		"path":      fwHelper.StringValueOrPrev(config.Path, types.StringNull()),
		"port":      fwHelper.StringValueOrPrev(config.Port, types.StringNull()),
		"protocol":  fwHelper.StringValueOrPrev(config.Protocol, types.StringNull()),
		"query":     fwHelper.StringValueOrPrev(config.Query, types.StringNull()),
	})
	return object
}

// albRemoveHeaderConfigObject converts the cloud-side remove header
// configuration into the nested object value.
func albRemoveHeaderConfigObject(config *albv20251030.RemoveHTTPHeaderInfo) types.Object {
	object, _ := types.ObjectValue(albRulesRemoveHeaderConfigAttrTypes(), map[string]attr.Value{
		"key": fwHelper.StringValueOrPrev(config.Key, types.StringNull()),
	})
	return object
}

// albRewriteConfigObject converts the cloud-side rewrite configuration into
// the nested object value.
func albRewriteConfigObject(config *albv20251030.HTTPRewriteInfo) types.Object {
	object, _ := types.ObjectValue(albRulesRewriteConfigAttrTypes(), map[string]attr.Value{
		"host":  fwHelper.StringValueOrPrev(config.Host, types.StringNull()),
		"path":  fwHelper.StringValueOrPrev(config.Path, types.StringNull()),
		"query": fwHelper.StringValueOrPrev(config.Query, types.StringNull()),
	})
	return object
}

// albTargetGroupConfigObject converts the cloud-side target group
// configuration into the nested object value.
func albTargetGroupConfigObject(ctx context.Context, config *albv20251030.TargetGroupConfig, diags *diag.Diagnostics) types.Object {
	targetGroups := make([]AlbRulesTargetGroupModel, 0, len(config.TargetGroups))
	for _, targetGroup := range config.TargetGroups {
		targetGroups = append(targetGroups, AlbRulesTargetGroupModel{
			TargetGroupId: fwHelper.StringValueOrPrev(targetGroup.TargetGroupId, types.StringNull()),
			Weight:        albInt64PointerValue(targetGroup.Weight),
		})
	}

	targetGroupList, listDiags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: albRulesTargetGroupAttrTypes()}, targetGroups)
	diags.Append(listDiags...)
	if diags.HasError() {
		return types.ObjectNull(albRulesTargetGroupConfigAttrTypes())
	}

	stickySession := types.ObjectNull(albRulesTargetGroupStickySessionAttrTypes())
	if config.TargetGroupStickySession != nil {
		stickySession, _ = types.ObjectValue(albRulesTargetGroupStickySessionAttrTypes(), map[string]attr.Value{
			"enabled": albBoolValueOrPrev(config.TargetGroupStickySession.Enabled),
			"timeout": albInt64PointerValue(config.TargetGroupStickySession.Timeout),
		})
	}

	object, objectDiags := types.ObjectValue(albRulesTargetGroupConfigAttrTypes(), map[string]attr.Value{
		"target_groups":               targetGroupList,
		"target_group_sticky_session": stickySession,
	})
	diags.Append(objectDiags...)
	return object
}

// albBoolValueOrPrev keeps the cloud value when it is known and falls back
// to null otherwise.
func albBoolValueOrPrev(cloud *bool) types.Bool {
	if cloud == nil {
		return types.BoolNull()
	}
	return types.BoolValue(*cloud)
}

// expandAlbRules converts the planned rules list into the SDK request
// structs. Fields left unset in the plan are omitted from the request.
func expandAlbRules(ctx context.Context, rules types.List) ([]*albv20251030.RuleInput, diag.Diagnostics) {
	var diags diag.Diagnostics

	var ruleModels []AlbRulesRuleModel
	diags.Append(rules.ElementsAs(ctx, &ruleModels, false)...)
	if diags.HasError() {
		return nil, diags
	}

	ruleInputs := make([]*albv20251030.RuleInput, 0, len(ruleModels))
	for _, ruleModel := range ruleModels {
		ruleInput := &albv20251030.RuleInput{}

		var actionModels []AlbRulesActionModel
		diags.Append(ruleModel.Actions.ElementsAs(ctx, &actionModels, false)...)
		if diags.HasError() {
			return nil, diags
		}

		actions := make([]*albv20251030.RuleAction, 0, len(actionModels))
		for _, actionModel := range actionModels {
			action := &albv20251030.RuleAction{
				Order: albInt64Pointer(actionModel.Order),
				Type:  albStringPointer(actionModel.Type),
			}

			if !actionModel.FixedResponseConfig.IsNull() && !actionModel.FixedResponseConfig.IsUnknown() {
				var fixedResponseConfig AlbRulesFixedResponseConfigModel
				diags.Append(actionModel.FixedResponseConfig.As(ctx, &fixedResponseConfig, basetypes.ObjectAsOptions{})...)
				if diags.HasError() {
					return nil, diags
				}
				action.FixedResponseConfig = &albv20251030.FixedResponseInfo{
					HttpCode:    albInt64Pointer(fixedResponseConfig.HttpCode),
					Content:     albStringPointer(fixedResponseConfig.Content),
					ContentType: albStringPointer(fixedResponseConfig.ContentType),
				}
			}

			if !actionModel.InsertHeaderConfig.IsNull() && !actionModel.InsertHeaderConfig.IsUnknown() {
				var insertHeaderConfig AlbRulesInsertHeaderConfigModel
				diags.Append(actionModel.InsertHeaderConfig.As(ctx, &insertHeaderConfig, basetypes.ObjectAsOptions{})...)
				if diags.HasError() {
					return nil, diags
				}
				action.InsertHeaderConfig = &albv20251030.InsertHTTPHeaderInfo{
					Key:       albStringPointer(insertHeaderConfig.Key),
					Value:     albStringPointer(insertHeaderConfig.Value),
					ValueType: albStringPointer(insertHeaderConfig.ValueType),
				}
			}

			if !actionModel.RedirectConfig.IsNull() && !actionModel.RedirectConfig.IsUnknown() {
				var redirectConfig AlbRulesRedirectConfigModel
				diags.Append(actionModel.RedirectConfig.As(ctx, &redirectConfig, basetypes.ObjectAsOptions{})...)
				if diags.HasError() {
					return nil, diags
				}
				action.RedirectConfig = &albv20251030.HTTPRedirectInfo{
					HttpCode: albInt64Pointer(redirectConfig.HttpCode),
					Host:     albStringPointer(redirectConfig.Host),
					Path:     albStringPointer(redirectConfig.Path),
					Port:     albStringPointer(redirectConfig.Port),
					Protocol: albStringPointer(redirectConfig.Protocol),
					Query:    albStringPointer(redirectConfig.Query),
				}
			}

			if !actionModel.RemoveHeaderConfig.IsNull() && !actionModel.RemoveHeaderConfig.IsUnknown() {
				var removeHeaderConfig AlbRulesRemoveHeaderConfigModel
				diags.Append(actionModel.RemoveHeaderConfig.As(ctx, &removeHeaderConfig, basetypes.ObjectAsOptions{})...)
				if diags.HasError() {
					return nil, diags
				}
				action.RemoveHeaderConfig = &albv20251030.RemoveHTTPHeaderInfo{
					Key: albStringPointer(removeHeaderConfig.Key),
				}
			}

			if !actionModel.RewriteConfig.IsNull() && !actionModel.RewriteConfig.IsUnknown() {
				var rewriteConfig AlbRulesRewriteConfigModel
				diags.Append(actionModel.RewriteConfig.As(ctx, &rewriteConfig, basetypes.ObjectAsOptions{})...)
				if diags.HasError() {
					return nil, diags
				}
				action.RewriteConfig = &albv20251030.HTTPRewriteInfo{
					Host:  albStringPointer(rewriteConfig.Host),
					Path:  albStringPointer(rewriteConfig.Path),
					Query: albStringPointer(rewriteConfig.Query),
				}
			}

			if !actionModel.TargetGroupConfig.IsNull() && !actionModel.TargetGroupConfig.IsUnknown() {
				targetGroupConfig, targetGroupDiags := expandAlbTargetGroupConfig(ctx, actionModel.TargetGroupConfig)
				diags.Append(targetGroupDiags...)
				if diags.HasError() {
					return nil, diags
				}
				action.TargetGroupConfig = targetGroupConfig
			}

			actions = append(actions, action)
		}
		ruleInput.Actions = actions

		var conditionModels []AlbRulesConditionModel
		diags.Append(ruleModel.Conditions.ElementsAs(ctx, &conditionModels, false)...)
		if diags.HasError() {
			return nil, diags
		}

		conditions := make([]*albv20251030.RuleCondition, 0, len(conditionModels))
		for _, conditionModel := range conditionModels {
			condition := &albv20251030.RuleCondition{
				Type: albStringPointer(conditionModel.Type),
			}

			if !conditionModel.CookieConfig.IsNull() && !conditionModel.CookieConfig.IsUnknown() {
				var cookieConfigs []AlbRulesCookieConfigModel
				diags.Append(conditionModel.CookieConfig.ElementsAs(ctx, &cookieConfigs, false)...)
				if diags.HasError() {
					return nil, diags
				}
				cookieInfos := make([]*albv20251030.HTTPCookieInfo, 0, len(cookieConfigs))
				for _, cookieConfig := range cookieConfigs {
					cookieInfos = append(cookieInfos, &albv20251030.HTTPCookieInfo{
						Key:   albStringPointer(cookieConfig.Key),
						Value: albStringPointer(cookieConfig.Value),
					})
				}
				condition.CookieConfig = cookieInfos
			}

			if !conditionModel.HeaderConfig.IsNull() && !conditionModel.HeaderConfig.IsUnknown() {
				var headerConfig AlbRulesHeaderConfigModel
				diags.Append(conditionModel.HeaderConfig.As(ctx, &headerConfig, basetypes.ObjectAsOptions{})...)
				if diags.HasError() {
					return nil, diags
				}
				condition.HeaderConfig = &albv20251030.HTTPHeaderInfo{
					Key:    albStringPointer(headerConfig.Key),
					Values: albStringSetToSDK(headerConfig.Values, &diags),
				}
				if diags.HasError() {
					return nil, diags
				}
			}

			condition.HostConfig = albStringSetToSDK(conditionModel.HostConfig, &diags)
			condition.MethodConfig = albStringSetToSDK(conditionModel.MethodConfig, &diags)
			condition.PathConfig = albStringSetToSDK(conditionModel.PathConfig, &diags)
			condition.SourceIpConfig = albStringSetToSDK(conditionModel.SourceIpConfig, &diags)
			if diags.HasError() {
				return nil, diags
			}

			if !conditionModel.QueryStringConfig.IsNull() && !conditionModel.QueryStringConfig.IsUnknown() {
				var queryStringConfigs []AlbRulesQueryStringConfigModel
				diags.Append(conditionModel.QueryStringConfig.ElementsAs(ctx, &queryStringConfigs, false)...)
				if diags.HasError() {
					return nil, diags
				}
				queryStringInfos := make([]*albv20251030.HTTPQueryStringInfo, 0, len(queryStringConfigs))
				for _, queryStringConfig := range queryStringConfigs {
					queryStringInfos = append(queryStringInfos, &albv20251030.HTTPQueryStringInfo{
						Key:   albStringPointer(queryStringConfig.Key),
						Value: albStringPointer(queryStringConfig.Value),
					})
				}
				condition.QueryStringConfig = queryStringInfos
			}

			conditions = append(conditions, condition)
		}
		ruleInput.Conditions = conditions

		ruleInput.Priority = albInt64Pointer(ruleModel.Priority)
		ruleInput.Direction = albStringPointer(ruleModel.Direction)
		ruleInput.RuleName = albStringPointer(ruleModel.RuleName)

		if !ruleModel.Tags.IsNull() && !ruleModel.Tags.IsUnknown() {
			var tagModels []AlbRulesTagModel
			diags.Append(ruleModel.Tags.ElementsAs(ctx, &tagModels, false)...)
			if diags.HasError() {
				return nil, diags
			}
			tags := make([]*albv20251030.TagInfo, 0, len(tagModels))
			for _, tagModel := range tagModels {
				tags = append(tags, &albv20251030.TagInfo{
					TagKey:   albStringPointer(tagModel.TagKey),
					TagValue: albStringPointer(tagModel.TagValue),
				})
			}
			ruleInput.Tags = tags
		}

		ruleInputs = append(ruleInputs, ruleInput)
	}

	return ruleInputs, diags
}

// expandAlbTargetGroupConfig converts the planned target group configuration
// object into the SDK request struct.
func expandAlbTargetGroupConfig(ctx context.Context, value types.Object) (*albv20251030.TargetGroupConfig, diag.Diagnostics) {
	var diags diag.Diagnostics

	var targetGroupConfig AlbRulesTargetGroupConfigModel
	diags.Append(value.As(ctx, &targetGroupConfig, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return nil, diags
	}

	config := &albv20251030.TargetGroupConfig{}

	var targetGroups []AlbRulesTargetGroupModel
	diags.Append(targetGroupConfig.TargetGroups.ElementsAs(ctx, &targetGroups, false)...)
	if diags.HasError() {
		return nil, diags
	}
	targetGroupTuples := make([]*albv20251030.TargetGroupTuple, 0, len(targetGroups))
	for _, targetGroup := range targetGroups {
		targetGroupTuples = append(targetGroupTuples, &albv20251030.TargetGroupTuple{
			TargetGroupId: albStringPointer(targetGroup.TargetGroupId),
			Weight:        albInt64Pointer(targetGroup.Weight),
		})
	}
	config.TargetGroups = targetGroupTuples

	if !targetGroupConfig.TargetGroupStickySession.IsNull() && !targetGroupConfig.TargetGroupStickySession.IsUnknown() {
		var stickySession AlbRulesTargetGroupStickySessionModel
		diags.Append(targetGroupConfig.TargetGroupStickySession.As(ctx, &stickySession, basetypes.ObjectAsOptions{})...)
		if diags.HasError() {
			return nil, diags
		}
		config.TargetGroupStickySession = &albv20251030.TargetGroupStickySession{
			Enabled: albBoolPointer(stickySession.Enabled),
			Timeout: albInt64Pointer(stickySession.Timeout),
		}
	}

	return config, diags
}

// expandAlbRuleModifies converts the planned rules into the SDK
// ModifyRulesAttributes structs. For rules whose identity is unchanged the
// planned rule_id is reused; otherwise the state rule id at the same index
// is carried over, so updates address the existing cloud objects.
func expandAlbRuleModifies(ctx context.Context, rules types.List, prevRules types.List) ([]*albv20251030.RuleModify, diag.Diagnostics) {
	var diags diag.Diagnostics

	ruleInputs, expandDiags := expandAlbRules(ctx, rules)
	diags.Append(expandDiags...)
	if diags.HasError() {
		return nil, diags
	}

	var prevRuleModels []AlbRulesRuleModel
	if !prevRules.IsNull() && !prevRules.IsUnknown() {
		diags.Append(prevRules.ElementsAs(ctx, &prevRuleModels, false)...)
		if diags.HasError() {
			return nil, diags
		}
	}

	ruleModifies := make([]*albv20251030.RuleModify, 0, len(ruleInputs))
	for i, ruleInput := range ruleInputs {
		ruleModify := &albv20251030.RuleModify{
			Actions:    ruleInput.Actions,
			Conditions: ruleInput.Conditions,
			Priority:   ruleInput.Priority,
			RuleName:   ruleInput.RuleName,
		}

		var planRuleModels []AlbRulesRuleModel
		diags.Append(rules.ElementsAs(ctx, &planRuleModels, false)...)
		if diags.HasError() {
			return nil, diags
		}

		if i < len(planRuleModels) && !planRuleModels[i].RuleId.IsNull() && !planRuleModels[i].RuleId.IsUnknown() && planRuleModels[i].RuleId.ValueString() != "" {
			ruleModify.RuleId = albSdkString(planRuleModels[i].RuleId.ValueString())
		} else if i < len(prevRuleModels) && !prevRuleModels[i].RuleId.IsNull() && !prevRuleModels[i].RuleId.IsUnknown() && prevRuleModels[i].RuleId.ValueString() != "" {
			ruleModify.RuleId = albSdkString(prevRuleModels[i].RuleId.ValueString())
		} else {
			diags.AddError(
				"Error updating ALB listener rules",
				fmt.Sprintf("rule index %d has no rule_id in plan or state; the rule set cannot be updated in place.", i),
			)
			return nil, diags
		}

		ruleModifies = append(ruleModifies, ruleModify)
	}

	return ruleModifies, diags
}

// albStateRuleIds collects the rule ids recorded in state; every id is
// required, otherwise the state cannot be destroyed reliably.
func albStateRuleIds(rules types.List, diags *diag.Diagnostics) []string {
	if rules.IsNull() || rules.IsUnknown() {
		diags.AddError(
			"Error deleting ALB listener rules",
			"The state records no forwarding rule ids; the resource cannot be destroyed reliably.",
		)
		return nil
	}

	ctx := context.Background()
	var ruleModels []AlbRulesRuleModel
	diags.Append(rules.ElementsAs(ctx, &ruleModels, false)...)
	if diags.HasError() {
		return nil
	}

	ruleIds := make([]string, 0, len(ruleModels))
	for _, ruleModel := range ruleModels {
		if ruleModel.RuleId.IsNull() || ruleModel.RuleId.IsUnknown() || ruleModel.RuleId.ValueString() == "" {
			diags.AddError(
				"Error deleting ALB listener rules",
				"The state records a forwarding rule without rule_id; the resource cannot be destroyed reliably.",
			)
			return nil
		}
		ruleIds = append(ruleIds, ruleModel.RuleId.ValueString())
	}
	return ruleIds
}

// albRuleIdsFromModifies collects the rule ids of the modify request.
func albRuleIdsFromModifies(ruleModifies []*albv20251030.RuleModify) []string {
	ruleIds := make([]string, 0, len(ruleModifies))
	for _, ruleModify := range ruleModifies {
		if ruleModify.RuleId != nil {
			ruleIds = append(ruleIds, *ruleModify.RuleId)
		}
	}
	return ruleIds
}

// albRuleIdPointers converts rule ids into SDK string pointers.
func albRuleIdPointers(ruleIds []string) []*string {
	ruleIdPointers := make([]*string, 0, len(ruleIds))
	for _, ruleId := range ruleIds {
		ruleIdPointers = append(ruleIdPointers, albSdkString(ruleId))
	}
	return ruleIdPointers
}

// albStringSetToSDK converts a set of strings into SDK string pointers.
func albStringSetToSDK(value types.Set, diags *diag.Diagnostics) []*string {
	if value.IsNull() || value.IsUnknown() {
		return nil
	}

	var strings []types.String
	diags.Append(value.ElementsAs(context.Background(), &strings, false)...)
	if diags.HasError() {
		return nil
	}

	result := make([]*string, 0, len(strings))
	for _, v := range strings {
		if v.IsNull() || v.IsUnknown() {
			continue
		}
		result = append(result, albSdkString(v.ValueString()))
	}
	return result
}

// albInt64PointerValue converts a nullable int64 into types.Int64.
func albInt64PointerValue(v *int64) types.Int64 {
	if v == nil {
		return types.Int64Null()
	}
	return types.Int64Value(*v)
}

// albSdkString returns a pointer to a copy of the string, mirroring the
// repo-standard internal/helper.String helper for SDK request fields.
func albSdkString(v string) *string {
	return &v
}

// albStringPointer returns a pointer to a copy of the value; null and
// unknown values become nil so the field is omitted from the request.
func albStringPointer(value types.String) *string {
	if value.IsNull() || value.IsUnknown() {
		return nil
	}
	return albSdkString(value.ValueString())
}

// albInt64Pointer returns a pointer to a copy of the value; null and
// unknown values become nil so the field is omitted from the request.
func albInt64Pointer(value types.Int64) *int64 {
	if value.IsNull() || value.IsUnknown() {
		return nil
	}
	v := value.ValueInt64()
	return &v
}

// albBoolPointer returns a pointer to a copy of the value; null and unknown
// values become nil so the field is omitted from the request.
func albBoolPointer(value types.Bool) *bool {
	if value.IsNull() || value.IsUnknown() {
		return nil
	}
	v := value.ValueBool()
	return &v
}

// ImportState imports the resource with the composite id
// load_balancer_id#listener_id.
func (r *AlbRulesResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
}
