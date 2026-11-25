package alb

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	albv20251030 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/alb/v20251030"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/fw"
	fwHelper "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/helper"
)

// Compile-time interface assertions.
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
// A single resource manages one forwarding rule of one ALB listener. Every API
// addresses the rule by the triple (load balancer id, listener id, rule id), so
// the Terraform id is the composite `<load_balancer_id>#<listener_id>#<rule_id>`
// while the rule id itself comes from the CreateRules response.
type AlbRulesResource struct {
	fw.ResourceWithConfigure
}

// AlbRulesResourceModel maps the schema attributes.
//
// The CreateRules/ModifyRulesAttributes APIs embed the rule body in the `Rules`
// array; because one Terraform resource owns exactly one rule, the fields of
// that array element are exposed as top-level attributes.
type AlbRulesResourceModel struct {
	Id             types.String `tfsdk:"id"`
	LoadBalancerId types.String `tfsdk:"load_balancer_id"`
	ListenerId     types.String `tfsdk:"listener_id"`
	Actions        types.List   `tfsdk:"actions"`
	Conditions     types.List   `tfsdk:"conditions"`
	Priority       types.Int64  `tfsdk:"priority"`
	Direction      types.String `tfsdk:"direction"`
	RuleName       types.String `tfsdk:"rule_name"`
	Tags           types.List   `tfsdk:"tags"`
}

// AlbRuleActionModel maps one element of the `actions` attribute.
type AlbRuleActionModel struct {
	Order               types.Int64  `tfsdk:"order"`
	Type                types.String `tfsdk:"type"`
	FixedResponseConfig types.Object `tfsdk:"fixed_response_config"`
	InsertHeaderConfig  types.Object `tfsdk:"insert_header_config"`
	RedirectConfig      types.Object `tfsdk:"redirect_config"`
	RemoveHeaderConfig  types.Object `tfsdk:"remove_header_config"`
	RewriteConfig       types.Object `tfsdk:"rewrite_config"`
	TargetGroupConfig   types.Object `tfsdk:"target_group_config"`
}

// AlbFixedResponseConfigModel maps the `fixed_response_config` object.
type AlbFixedResponseConfigModel struct {
	HttpCode    types.Int64  `tfsdk:"http_code"`
	Content     types.String `tfsdk:"content"`
	ContentType types.String `tfsdk:"content_type"`
}

// AlbInsertHeaderConfigModel maps the `insert_header_config` object.
type AlbInsertHeaderConfigModel struct {
	Key       types.String `tfsdk:"key"`
	Value     types.String `tfsdk:"value"`
	ValueType types.String `tfsdk:"value_type"`
}

// AlbRedirectConfigModel maps the `redirect_config` object.
type AlbRedirectConfigModel struct {
	HttpCode types.Int64  `tfsdk:"http_code"`
	Host     types.String `tfsdk:"host"`
	Path     types.String `tfsdk:"path"`
	Port     types.String `tfsdk:"port"`
	Protocol types.String `tfsdk:"protocol"`
	Query    types.String `tfsdk:"query"`
}

// AlbRemoveHeaderConfigModel maps the `remove_header_config` object.
type AlbRemoveHeaderConfigModel struct {
	Key types.String `tfsdk:"key"`
}

// AlbRewriteConfigModel maps the `rewrite_config` object.
type AlbRewriteConfigModel struct {
	Host  types.String `tfsdk:"host"`
	Path  types.String `tfsdk:"path"`
	Query types.String `tfsdk:"query"`
}

// AlbTargetGroupConfigModel maps the `target_group_config` object.
type AlbTargetGroupConfigModel struct {
	TargetGroups             types.List   `tfsdk:"target_groups"`
	TargetGroupStickySession types.Object `tfsdk:"target_group_sticky_session"`
}

// AlbTargetGroupTupleModel maps one element of `target_group_config.target_groups`.
type AlbTargetGroupTupleModel struct {
	TargetGroupId types.String `tfsdk:"target_group_id"`
	Weight        types.Int64  `tfsdk:"weight"`
}

// AlbTargetGroupStickySessionModel maps the `target_group_sticky_session`
// object.
type AlbTargetGroupStickySessionModel struct {
	Enabled types.Bool  `tfsdk:"enabled"`
	Timeout types.Int64 `tfsdk:"timeout"`
}

// AlbRuleConditionModel maps one element of the `conditions` attribute.
type AlbRuleConditionModel struct {
	Type              types.String `tfsdk:"type"`
	CookieConfig      types.List   `tfsdk:"cookie_config"`
	HeaderConfig      types.Object `tfsdk:"header_config"`
	HostConfig        types.List   `tfsdk:"host_config"`
	MethodConfig      types.List   `tfsdk:"method_config"`
	PathConfig        types.List   `tfsdk:"path_config"`
	QueryStringConfig types.List   `tfsdk:"query_string_config"`
	SourceIpConfig    types.List   `tfsdk:"source_ip_config"`
}

// AlbCookieConfigModel maps one element of `cookie_config`.
type AlbCookieConfigModel struct {
	Key   types.String `tfsdk:"key"`
	Value types.String `tfsdk:"value"`
}

// AlbHeaderConfigModel maps the `header_config` object.
type AlbHeaderConfigModel struct {
	Key    types.String `tfsdk:"key"`
	Values types.List   `tfsdk:"values"`
}

// AlbQueryStringConfigModel maps one element of `query_string_config`.
type AlbQueryStringConfigModel struct {
	Key   types.String `tfsdk:"key"`
	Value types.String `tfsdk:"value"`
}

// AlbTagModel maps one element of the `tags` attribute.
type AlbTagModel struct {
	TagKey   types.String `tfsdk:"tag_key"`
	TagValue types.String `tfsdk:"tag_value"`
}

var albFixedResponseConfigAttributes = map[string]schema.Attribute{
	"http_code": schema.Int64Attribute{
		Required:    true,
		Description: "HTTP response code returned. 2xx, 4xx and 5xx are supported.",
	},
	"content": schema.StringAttribute{
		Optional:    true,
		Description: "Fixed content returned. Only ASCII characters are supported, with a maximum of 1KB.",
	},
	"content_type": schema.StringAttribute{
		Optional:    true,
		Description: "Format of the fixed content returned. Valid values: text/plain, text/css, text/html, application/javascript, application/json.",
	},
}

var albInsertHeaderConfigAttributes = map[string]schema.Attribute{
	"key": schema.StringAttribute{
		Required:    true,
		Description: "Key of the HTTP header to insert, 1 ~ 40 characters. Supported characters: a-z A-Z 0-9 - _ . Chinese characters are not supported, and Cookie, Host, Content-Length, Connection, Upgrade, transfer-encoding, keep-alive, te, authority, x-forwarded-for, x-forwarded-proto, x-forwarded-host, x-forwarded-port are not supported.",
	},
	"value": schema.StringAttribute{
		Required:    true,
		Description: "Type of the HTTP header value. When ValueType is SystemDefined, valid values are ClientPort: client port, ClientIp: client IP address, Protocol: protocol of the client request, CLBPort: listening port of the load balancer instance. When ValueType is UserDefined, 1 ~ 128 printable characters, double quotes are not supported, the value cannot start or end with a space, and cannot end with \\\\. When ValueType is ReferenceHeader, references a header of the request, 1 ~ 128 printable characters, double quotes are not supported, the value cannot start or end with a space, and cannot end with \\\\.",
	},
	"value_type": schema.StringAttribute{
		Required:    true,
		Description: "Type of the HTTP header value. Valid values: SystemDefined: system defined header; UserDefined: user defined header; ReferenceHeader: references a header of the request.",
	},
}

var albRedirectConfigAttributes = map[string]schema.Attribute{
	"http_code": schema.Int64Attribute{
		Required:    true,
		Description: "HTTP code of the redirect. 301, 302, 303, 307 and 308 are supported.",
	},
	"host": schema.StringAttribute{
		Optional:    true,
		Description: "Host address of the redirect, default value ${host}. 3 ~ 128 characters, supported characters: a-z 0-9 _ . -.",
	},
	"path": schema.StringAttribute{
		Optional:    true,
		Description: "Path of the redirect, default value ${path}. 1 ~ 128 characters, supported characters: a-z A-Z 0-9 ? = _ . - / : .",
	},
	"port": schema.StringAttribute{
		Optional:    true,
		Description: "Port of the redirect, default value ${port}. Value range: 1 ~ 65535.",
	},
	"protocol": schema.StringAttribute{
		Optional:    true,
		Description: "Protocol of the redirect. Valid values: HTTP, HTTPS. Default value: ${protocol}.",
	},
	"query": schema.StringAttribute{
		Optional:    true,
		Description: "Query string of the redirect, default value ${query}. 1 ~ 128 characters, printable characters are supported, #[]{}|<> and spaces are not supported.",
	},
}

var albRemoveHeaderConfigAttributes = map[string]schema.Attribute{
	"key": schema.StringAttribute{
		Required:    true,
		Description: "Key of the HTTP header to remove, 1 ~ 40 characters. Supported characters: a-z A-Z 0-9 - _ . Cookie, Host, Content-Length, Connection, Upgrade, transfer-encoding, keep-alive, te, authority, x-forwarded-for, x-forwarded-proto, x-forwarded-host, x-forwarded-port and server are not supported.",
	},
}

var albRewriteConfigAttributes = map[string]schema.Attribute{
	"host": schema.StringAttribute{
		Optional:    true,
		Description: "Host address of the rewrite, default value ${host}. 3 ~ 128 characters, supported characters: a-z 0-9 _ . -.",
	},
	"path": schema.StringAttribute{
		Optional:    true,
		Description: "Path of the rewrite, default value ${path}. 1 ~ 128 characters, supported characters: a-z A-Z 0-9 ? = _ . - / : .",
	},
	"query": schema.StringAttribute{
		Optional:    true,
		Description: "Query string of the rewrite, default value ${query}. 1 ~ 128 characters, printable characters are supported, #[]{}|<> and spaces are not supported.",
	},
}

var albTargetGroupTupleAttributes = map[string]schema.Attribute{
	"target_group_id": schema.StringAttribute{
		Required:    true,
		Description: "Target group ID, in the format of lbtg- followed by 8 alphanumeric characters.",
	},
	"weight": schema.Int64Attribute{
		Optional:    true,
		Description: "Weight. Value range: [0, 100]. Default value: 10.",
	},
}

var albTargetGroupStickySessionAttributes = map[string]schema.Attribute{
	"enabled": schema.BoolAttribute{
		Required:    true,
		Description: "Whether to enable session persistence. Disabled by default.",
	},
	"timeout": schema.Int64Attribute{
		Optional:    true,
		Description: "Timeout period in seconds. Value range: 1 ~ 86400. Default value: 1000.",
	},
}

var albTargetGroupConfigAttributes = map[string]schema.Attribute{
	"target_groups": schema.ListNestedAttribute{
		Required:     true,
		Description:  "Target group list.",
		NestedObject: schema.NestedAttributeObject{Attributes: albTargetGroupTupleAttributes},
	},
	"target_group_sticky_session": schema.SingleNestedAttribute{
		Optional:    true,
		Description: "Session persistence between target groups.",
		Attributes:  albTargetGroupStickySessionAttributes,
	},
}

var albRuleActionAttributes = map[string]schema.Attribute{
	"order": schema.Int64Attribute{
		Required:    true,
		Description: "Execution order of the forwarding action. It must be unique and actions are executed in ascending order. Value range: 1 ~ 50000.",
	},
	"type": schema.StringAttribute{
		Required:    true,
		Description: "Type of the forwarding action. Valid values: TargetGroup: forward to the target group; Redirect: redirect; FixedResponse: return fixed content; Rewrite: rewrite; InsertHeader: insert an HTTP header; RemoveHeader: remove an HTTP header. A forwarding action must contain one of TargetGroup, Redirect and FixedResponse, and it must be executed last.",
	},
	"fixed_response_config": schema.SingleNestedAttribute{
		Optional:    true,
		Description: "Configuration of the fixed response content.",
		Attributes:  albFixedResponseConfigAttributes,
	},
	"insert_header_config": schema.SingleNestedAttribute{
		Optional:    true,
		Description: "Configuration of the inserted HTTP header.",
		Attributes:  albInsertHeaderConfigAttributes,
	},
	"redirect_config": schema.SingleNestedAttribute{
		Optional:    true,
		Description: "Redirect configuration. Except for HttpCode, the other configurations cannot all use the default values.",
		Attributes:  albRedirectConfigAttributes,
	},
	"remove_header_config": schema.SingleNestedAttribute{
		Optional:    true,
		Description: "Configuration of the removed HTTP header.",
		Attributes:  albRemoveHeaderConfigAttributes,
	},
	"rewrite_config": schema.SingleNestedAttribute{
		Optional:    true,
		Description: "Rewrite configuration.",
		Attributes:  albRewriteConfigAttributes,
	},
	"target_group_config": schema.SingleNestedAttribute{
		Optional:    true,
		Description: "Configuration of the forwarding target group.",
		Attributes:  albTargetGroupConfigAttributes,
	},
}

var albCookieConfigAttributes = map[string]schema.Attribute{
	"key": schema.StringAttribute{
		Required:    true,
		Description: "Cookie key, 1 ~ 64 characters. Letters, digits and underscores are supported.",
	},
	"value": schema.StringAttribute{
		Required:    true,
		Description: "Cookie value, 1 ~ 128 characters. Printable characters are supported.",
	},
}

var albHeaderConfigAttributes = map[string]schema.Attribute{
	"key": schema.StringAttribute{
		Required:    true,
		Description: "Key of the HTTP header, 1 ~ 40 characters. Supported characters: a-z A-Z 0-9 - _ . Chinese characters, Host and Cookie are not supported.",
	},
	"values": schema.ListAttribute{
		Required:    true,
		ElementType: types.StringType,
		Description: "Values of the HTTP header, 1 ~ 128 characters. Printable characters are supported. Double quotes are not supported, the value cannot start or end with a space, and cannot end with \\\\.",
	},
}

var albQueryStringConfigAttributes = map[string]schema.Attribute{
	"key": schema.StringAttribute{
		Required:    true,
		Description: "Key of the query string, 1 ~ 16 characters. Printable characters are supported, spaces and #[]{}\\|<> are not supported. * matches multiple characters and ? matches a single character.",
	},
	"value": schema.StringAttribute{
		Required:    true,
		Description: "Value of the query string, 1 ~ 128 characters. Printable characters are supported, spaces and #[]{}\\|<> are not supported. * matches multiple characters and ? matches a single character.",
	},
}

var albRuleConditionAttributes = map[string]schema.Attribute{
	"type": schema.StringAttribute{
		Required:    true,
		Description: "Type of the forwarding condition. Valid values: Host: host; Path: path; Header: HTTP header field; QueryString: HTTP query string; Method: request method; Cookie: cookie; SourceIp: source IP.",
	},
	"cookie_config": schema.ListNestedAttribute{
		Optional:     true,
		Description:  "Cookie configuration.",
		NestedObject: schema.NestedAttributeObject{Attributes: albCookieConfigAttributes},
	},
	"header_config": schema.SingleNestedAttribute{
		Optional:    true,
		Description: "HTTP header configuration.",
		Attributes:  albHeaderConfigAttributes,
	},
	"host_config": schema.ListAttribute{
		Optional:    true,
		ElementType: types.StringType,
		Description: "Host name. Host configuration can appear only once in a rule, 3 ~ 128 characters, exact match, regular expression match and wildcard match are supported. It cannot start or end with a period (.) or an underscore (_). Exact match supported characters: a-z 0-9 . - _ . Regular expression match starts with a tilde (~), supported characters: a-z 0-9 . - ? = ~ _ - + \\ ^ * ! $ & | ( ) [ ] . Wildcard match uses an asterisk (*) for multiple characters and a question mark (?) for a single character, supported characters: a-z 0-9 . - _ * ?.",
	},
	"method_config": schema.ListAttribute{
		Optional:    true,
		ElementType: types.StringType,
		Description: "Request method. Valid values: HEAD, GET, POST, OPTIONS, PUT, PATCH, DELETE.",
	},
	"path_config": schema.ListAttribute{
		Optional:    true,
		ElementType: types.StringType,
		Description: "Forwarding path, 1 ~ 128 characters, exact match, regular expression match and wildcard match are supported. Exact match supported characters: a-z A-Z 0-9 . - _ / = :. Regular expression match must start with ~, where ~ means case-sensitive and ~* means case-insensitive, supported characters: a-z A-Z 0-9 . - _ / = ? ~ ^ * $ : ( ) [ ] + |. Wildcard match uses * for multiple characters and ? for a single character, supported characters: a-z A-Z 0-9 . - _ / = :.",
	},
	"query_string_config": schema.ListNestedAttribute{
		Optional:     true,
		Description:  "Query string configuration.",
		NestedObject: schema.NestedAttributeObject{Attributes: albQueryStringConfigAttributes},
	},
	"source_ip_config": schema.ListAttribute{
		Optional:    true,
		ElementType: types.StringType,
		Description: "Source IP matching configuration. CIDR format, IP address x.x.x.x/32, IP range x.x.x.x/24.",
	},
}

var albTagAttributes = map[string]schema.Attribute{
	"tag_key": schema.StringAttribute{
		Required:    true,
		Description: "Key of the tag.",
	},
	"tag_value": schema.StringAttribute{
		Required:    true,
		Description: "Value of the tag.",
	},
}

var (
	albFixedResponseConfigAttrTypes = albAttrTypesOf(albFixedResponseConfigAttributes)
	albInsertHeaderConfigAttrTypes  = albAttrTypesOf(albInsertHeaderConfigAttributes)
	albRedirectConfigAttrTypes      = albAttrTypesOf(albRedirectConfigAttributes)
	albRemoveHeaderConfigAttrTypes  = albAttrTypesOf(albRemoveHeaderConfigAttributes)
	albRewriteConfigAttrTypes       = albAttrTypesOf(albRewriteConfigAttributes)
	albTargetGroupTupleAttrTypes    = albAttrTypesOf(albTargetGroupTupleAttributes)
	albTargetGroupStickyAttrTypes   = albAttrTypesOf(albTargetGroupStickySessionAttributes)
	albTargetGroupConfigAttrTypes   = albAttrTypesOf(albTargetGroupConfigAttributes)
	albRuleActionAttrTypes          = albAttrTypesOf(albRuleActionAttributes)
	albCookieConfigAttrTypes        = albAttrTypesOf(albCookieConfigAttributes)
	albHeaderConfigAttrTypes        = albAttrTypesOf(albHeaderConfigAttributes)
	albQueryStringConfigAttrTypes   = albAttrTypesOf(albQueryStringConfigAttributes)
	albRuleConditionAttrTypes       = albAttrTypesOf(albRuleConditionAttributes)
	albTagAttrTypes                 = albAttrTypesOf(albTagAttributes)

	albTargetGroupTupleObjectType  = types.ObjectType{AttrTypes: albTargetGroupTupleAttrTypes}
	albCookieConfigObjectType      = types.ObjectType{AttrTypes: albCookieConfigAttrTypes}
	albQueryStringConfigObjectType = types.ObjectType{AttrTypes: albQueryStringConfigAttrTypes}
	albRuleActionObjectType        = types.ObjectType{AttrTypes: albRuleActionAttrTypes}
	albRuleConditionObjectType     = types.ObjectType{AttrTypes: albRuleConditionAttrTypes}
	albTagObjectType               = types.ObjectType{AttrTypes: albTagAttrTypes}
)

// albAttrTypesOf mirrors the value types of a nested attribute map so the
// conversion helpers stay in sync with the schema definition.
func albAttrTypesOf(attributes map[string]schema.Attribute) map[string]attr.Type {
	result := make(map[string]attr.Type, len(attributes))
	for name, attribute := range attributes {
		result[name] = attribute.GetType()
	}
	return result
}

func (r *AlbRulesResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "tencentcloud_alb_rules"
}

func (r *AlbRulesResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a resource to manage a forwarding rule of an ALB listener.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "ID of the resource, in the form `<load_balancer_id>#<listener_id>#<rule_id>`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"load_balancer_id": schema.StringAttribute{
				Required:    true,
				Description: "Load balancer instance ID, in the format of alb- followed by 8 alphanumeric characters.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"listener_id": schema.StringAttribute{
				Required:    true,
				Description: "Listener ID, in the format of lst- followed by 8 alphanumeric characters.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"actions": schema.ListNestedAttribute{
				Required:     true,
				Description:  "Forwarding action list.",
				NestedObject: schema.NestedAttributeObject{Attributes: albRuleActionAttributes},
			},
			"conditions": schema.ListNestedAttribute{
				Required:     true,
				Description:  "Forwarding condition list.",
				NestedObject: schema.NestedAttributeObject{Attributes: albRuleConditionAttributes},
			},
			"priority": schema.Int64Attribute{
				Required:    true,
				Description: "Priority. A smaller value means a higher priority. Value range: 1 ~ 10000.",
			},
			"direction": schema.StringAttribute{
				Optional:    true,
				Description: "Direction of the forwarding rule. Request: request direction from the client to the load balancer, Response: response direction from the backend server to the load balancer. Default value: Request.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"rule_name": schema.StringAttribute{
				Optional:    true,
				Description: "Name of the forwarding rule, 1 ~ 255 characters. Digits, upper and lower case letters, Chinese characters, periods (.), underscores (_) and hyphens (-) are supported.",
			},
			"tags": schema.ListNestedAttribute{
				Optional:     true,
				Description:  "Tag list.",
				NestedObject: schema.NestedAttributeObject{Attributes: albTagAttributes},
				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplace(),
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

	createRulesRequest := albv20251030.NewCreateRulesRequest()
	createRulesRequest.LoadBalancerId = plan.LoadBalancerId.ValueStringPointer()
	createRulesRequest.ListenerId = plan.ListenerId.ValueStringPointer()
	createRulesRequest.Rules = []*albv20251030.RuleInput{albRuleInputFromPlan(ctx, &plan, &resp.Diagnostics)}
	if resp.Diagnostics.HasError() {
		return
	}

	service := NewAlbService(r.Client())
	ruleId, err := service.CreateRules(ctx, createRulesRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating ALB rule",
			fmt.Sprintf("CreateRules failed: %v", err),
		)

		return
	}

	// Never write an empty id into state.
	if ruleId == "" {
		resp.Diagnostics.AddError(
			"Error creating ALB rule",
			"CreateRules returned an empty rule id.",
		)

		return
	}

	plan.Id = types.StringValue(fwHelper.IdFormat(plan.LoadBalancerId.ValueString(), plan.ListenerId.ValueString(), ruleId))

	// Persist the identity before reading back so a partial failure still
	// leaves a destroyable resource in state.
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	describeRulesRequest := albv20251030.NewDescribeRulesRequest()
	describeRulesRequest.LoadBalancerId = plan.LoadBalancerId.ValueStringPointer()
	describeRulesRequest.ListenerId = plan.ListenerId.ValueStringPointer()
	describeRulesRequest.RuleIds = []*string{fwHelper.String(ruleId)}

	rules, err := service.DescribeRules(ctx, describeRulesRequest)
	if err != nil {
		resp.Diagnostics.AddError("Error reading ALB rule after create", err.Error())
		return
	}

	item := findAlbRule(rules, ruleId)
	if item == nil {
		resp.Diagnostics.AddError(
			"ALB rule not readable after create",
			fmt.Sprintf("rule %s was created but cannot be read back", plan.Id.ValueString()),
		)

		return
	}

	resp.Diagnostics.Append(refreshAlbRulesModel(ctx, &plan, item)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Read handles `terraform refresh` and post-apply state refresh.
func (r *AlbRulesResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state AlbRulesResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	loadBalancerId, listenerId, ruleId, err := albRulesParseId(state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid resource id", err.Error())
		return
	}

	describeRulesRequest := albv20251030.NewDescribeRulesRequest()
	describeRulesRequest.LoadBalancerId = fwHelper.String(loadBalancerId)
	describeRulesRequest.ListenerId = fwHelper.String(listenerId)
	describeRulesRequest.RuleIds = []*string{fwHelper.String(ruleId)}

	service := NewAlbService(r.Client())

	rules, err := service.DescribeRules(ctx, describeRulesRequest)
	if err != nil {
		resp.Diagnostics.AddError("Error reading ALB rule", err.Error())
		return
	}

	item := findAlbRule(rules, ruleId)

	// Confirmed "not found" on the cloud side => drop it from state so
	// Terraform plans a recreation.
	if item == nil {
		resp.Diagnostics.AddWarning(
			"Resource not found",
			fmt.Sprintf("ALB rule %s not found, removing from state.", state.Id.ValueString()),
		)

		resp.State.RemoveResource(ctx)
		return
	}

	resp.Diagnostics.Append(refreshAlbRulesModel(ctx, &state, item)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update handles in-place modifications through ModifyRulesAttributes.
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

	loadBalancerId, listenerId, ruleId, err := albRulesParseId(state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid resource id", err.Error())
		return
	}

	service := NewAlbService(r.Client())

	// Only the attributes ModifyRulesAttributes covers go through the write
	// API; direction and tags are RequiresReplace, and a plan that differs in
	// the `timeouts` block alone would be a pointless no-change request.
	//
	// TODO(human): the modify API does not document whether an omitted
	// `Actions`/`Conditions` means "keep" or "clear", so removing those
	// attributes from the configuration is passed through to the cloud as an
	// omitted field; confirm the clear semantics.
	if !plan.Actions.Equal(state.Actions) || !plan.Conditions.Equal(state.Conditions) ||
		!plan.Priority.Equal(state.Priority) || !plan.RuleName.Equal(state.RuleName) {
		modifyRulesAttributesRequest := albv20251030.NewModifyRulesAttributesRequest()
		modifyRulesAttributesRequest.LoadBalancerId = fwHelper.String(loadBalancerId)
		modifyRulesAttributesRequest.ListenerId = fwHelper.String(listenerId)

		ruleModify := &albv20251030.RuleModify{
			RuleId: fwHelper.String(ruleId),
		}
		if !plan.Actions.IsNull() && !plan.Actions.IsUnknown() {
			ruleModify.Actions = albExpandRuleActions(ctx, plan.Actions, &resp.Diagnostics)
		}
		if !plan.Conditions.IsNull() && !plan.Conditions.IsUnknown() {
			ruleModify.Conditions = albExpandRuleConditions(ctx, plan.Conditions, &resp.Diagnostics)
		}
		if !plan.Priority.IsNull() && !plan.Priority.IsUnknown() {
			ruleModify.Priority = plan.Priority.ValueInt64Pointer()
		}
		if !plan.RuleName.IsNull() && !plan.RuleName.IsUnknown() {
			ruleModify.RuleName = plan.RuleName.ValueStringPointer()
		}
		if resp.Diagnostics.HasError() {
			return
		}

		modifyRulesAttributesRequest.Rules = []*albv20251030.RuleModify{ruleModify}

		if err := service.ModifyRulesAttributes(ctx, modifyRulesAttributesRequest); err != nil {
			resp.Diagnostics.AddError(
				"Error updating ALB rule",
				fmt.Sprintf("ModifyRulesAttributes failed: %v", err),
			)

			return
		}
	}

	describeRulesRequest := albv20251030.NewDescribeRulesRequest()
	describeRulesRequest.LoadBalancerId = fwHelper.String(loadBalancerId)
	describeRulesRequest.ListenerId = fwHelper.String(listenerId)
	describeRulesRequest.RuleIds = []*string{fwHelper.String(ruleId)}

	rules, err := service.DescribeRules(ctx, describeRulesRequest)
	if err != nil {
		resp.Diagnostics.AddError("Error reading ALB rule after update", err.Error())
		return
	}

	item := findAlbRule(rules, ruleId)
	if item == nil {
		resp.Diagnostics.AddError(
			"ALB rule not readable after update",
			fmt.Sprintf("rule %s was updated but cannot be read back", state.Id.ValueString()),
		)

		return
	}

	plan.Id = state.Id
	resp.Diagnostics.Append(refreshAlbRulesModel(ctx, &plan, item)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete handles resource deletion.
func (r *AlbRulesResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state AlbRulesResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	loadBalancerId, listenerId, ruleId, err := albRulesParseId(state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid resource id", err.Error())
		return
	}

	deleteRulesRequest := albv20251030.NewDeleteRulesRequest()
	deleteRulesRequest.LoadBalancerId = fwHelper.String(loadBalancerId)
	deleteRulesRequest.ListenerId = fwHelper.String(listenerId)
	deleteRulesRequest.RuleIds = []*string{fwHelper.String(ruleId)}

	service := NewAlbService(r.Client())
	if err := service.DeleteRules(ctx, deleteRulesRequest); err != nil {
		resp.Diagnostics.AddError(
			"Error deleting ALB rule",
			fmt.Sprintf("DeleteRules failed: %v", err),
		)

		return
	}
}

// ImportState parses the composite
// `<load_balancer_id>#<listener_id>#<rule_id>` id.
func (r *AlbRulesResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	loadBalancerId, listenerId, _, err := albRulesParseId(req.ID)
	if err != nil {
		resp.Diagnostics.AddError("Invalid resource id", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("load_balancer_id"), loadBalancerId)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("listener_id"), listenerId)...)
}

// albRulesParseId splits the composite id into its three segments.
func albRulesParseId(id string) (loadBalancerId string, listenerId string, ruleId string, err error) {
	parts, err := fwHelper.IdParseErr(id, 3)
	if err != nil {
		return "", "", "", fmt.Errorf("the expected id format is <load_balancer_id>#<listener_id>#<rule_id>: %v", err)
	}

	return parts[0], parts[1], parts[2], nil
}

// findAlbRule picks this resource's own rule out of a describe answer.
func findAlbRule(rules []*albv20251030.RuleOutput, ruleId string) *albv20251030.RuleOutput {
	for _, rule := range rules {
		if rule != nil && rule.RuleId != nil && *rule.RuleId == ruleId {
			return rule
		}
	}

	return nil
}

// albRuleInputFromPlan builds the single rule the create API receives.
func albRuleInputFromPlan(ctx context.Context, m *AlbRulesResourceModel, diags *diag.Diagnostics) *albv20251030.RuleInput {
	rule := &albv20251030.RuleInput{
		Actions:    albExpandRuleActions(ctx, m.Actions, diags),
		Conditions: albExpandRuleConditions(ctx, m.Conditions, diags),
	}
	if diags.HasError() {
		return nil
	}

	if !m.Priority.IsNull() && !m.Priority.IsUnknown() {
		rule.Priority = m.Priority.ValueInt64Pointer()
	}
	if !m.Direction.IsNull() && !m.Direction.IsUnknown() {
		rule.Direction = m.Direction.ValueStringPointer()
	}
	if !m.RuleName.IsNull() && !m.RuleName.IsUnknown() {
		rule.RuleName = m.RuleName.ValueStringPointer()
	}
	rule.Tags = albExpandTags(ctx, m.Tags, diags)

	return rule
}

// refreshAlbRulesModel copies the cloud-reported rule into the model. Fields
// the cloud omits keep their previous value so a configured attribute never
// turns into null after apply.
func refreshAlbRulesModel(ctx context.Context, m *AlbRulesResourceModel, item *albv20251030.RuleOutput) diag.Diagnostics {
	var diags diag.Diagnostics

	if item == nil {
		return diags
	}

	m.Priority = fwHelper.Int64ValueOrPrev(item.Priority, m.Priority)
	m.Direction = fwHelper.StringValueOrPrev(item.Direction, m.Direction)
	m.RuleName = fwHelper.StringValueOrPrev(item.RuleName, m.RuleName)

	if item.Actions != nil {
		m.Actions = albFlattenRuleActions(ctx, item.Actions, &diags)
	}
	if item.Conditions != nil {
		m.Conditions = albFlattenRuleConditions(ctx, item.Conditions, &diags)
	}
	if item.Tags != nil {
		m.Tags = albFlattenTags(ctx, item.Tags, &diags)
	}

	return diags
}

// albExpandRuleActions converts the `actions` attribute into SDK rule actions.
func albExpandRuleActions(ctx context.Context, list types.List, diags *diag.Diagnostics) []*albv20251030.RuleAction {
	if list.IsNull() || list.IsUnknown() {
		return nil
	}

	var models []AlbRuleActionModel
	d := list.ElementsAs(ctx, &models, false)
	diags.Append(d...)
	if d.HasError() {
		return nil
	}

	actions := make([]*albv20251030.RuleAction, 0, len(models))
	for _, model := range models {
		action := &albv20251030.RuleAction{
			Order:               fwHelper.Int64PointerFromValue(model.Order),
			Type:                fwHelper.StringPointerFromValue(model.Type),
			FixedResponseConfig: albExpandFixedResponseConfig(ctx, model.FixedResponseConfig, diags),
			InsertHeaderConfig:  albExpandInsertHeaderConfig(ctx, model.InsertHeaderConfig, diags),
			RedirectConfig:      albExpandRedirectConfig(ctx, model.RedirectConfig, diags),
			RemoveHeaderConfig:  albExpandRemoveHeaderConfig(ctx, model.RemoveHeaderConfig, diags),
			RewriteConfig:       albExpandRewriteConfig(ctx, model.RewriteConfig, diags),
			TargetGroupConfig:   albExpandTargetGroupConfig(ctx, model.TargetGroupConfig, diags),
		}
		if diags.HasError() {
			return nil
		}

		actions = append(actions, action)
	}

	return actions
}

// albFlattenRuleActions converts SDK rule actions into the `actions`
// attribute.
func albFlattenRuleActions(ctx context.Context, items []*albv20251030.RuleAction, diags *diag.Diagnostics) types.List {
	models := make([]AlbRuleActionModel, 0, len(items))
	for _, item := range items {
		if item == nil {
			continue
		}

		models = append(models, AlbRuleActionModel{
			Order:               fwHelper.Int64ValueOrNull(item.Order),
			Type:                fwHelper.StringPointerValueOrNull(item.Type),
			FixedResponseConfig: albFlattenFixedResponseConfig(ctx, item.FixedResponseConfig, diags),
			InsertHeaderConfig:  albFlattenInsertHeaderConfig(ctx, item.InsertHeaderConfig, diags),
			RedirectConfig:      albFlattenRedirectConfig(ctx, item.RedirectConfig, diags),
			RemoveHeaderConfig:  albFlattenRemoveHeaderConfig(ctx, item.RemoveHeaderConfig, diags),
			RewriteConfig:       albFlattenRewriteConfig(ctx, item.RewriteConfig, diags),
			TargetGroupConfig:   albFlattenTargetGroupConfig(ctx, item.TargetGroupConfig, diags),
		})
	}

	list, d := types.ListValueFrom(ctx, albRuleActionObjectType, models)
	diags.Append(d...)

	return list
}

// albExpandFixedResponseConfig converts the `fixed_response_config` attribute.
func albExpandFixedResponseConfig(ctx context.Context, object types.Object, diags *diag.Diagnostics) *albv20251030.FixedResponseInfo {
	if object.IsNull() || object.IsUnknown() {
		return nil
	}

	var model AlbFixedResponseConfigModel
	d := object.As(ctx, &model, basetypes.ObjectAsOptions{})
	diags.Append(d...)
	if d.HasError() {
		return nil
	}

	return &albv20251030.FixedResponseInfo{
		HttpCode:    fwHelper.Int64PointerFromValue(model.HttpCode),
		Content:     fwHelper.StringPointerFromValue(model.Content),
		ContentType: fwHelper.StringPointerFromValue(model.ContentType),
	}
}

// albFlattenFixedResponseConfig converts an SDK fixed response configuration.
func albFlattenFixedResponseConfig(ctx context.Context, item *albv20251030.FixedResponseInfo, diags *diag.Diagnostics) types.Object {
	if item == nil {
		return types.ObjectNull(albFixedResponseConfigAttrTypes)
	}

	value, d := types.ObjectValueFrom(ctx, albFixedResponseConfigAttrTypes, AlbFixedResponseConfigModel{
		HttpCode:    fwHelper.Int64ValueOrNull(item.HttpCode),
		Content:     fwHelper.StringPointerValueOrNull(item.Content),
		ContentType: fwHelper.StringPointerValueOrNull(item.ContentType),
	})
	diags.Append(d...)

	return value
}

// albExpandInsertHeaderConfig converts the `insert_header_config` attribute.
func albExpandInsertHeaderConfig(ctx context.Context, object types.Object, diags *diag.Diagnostics) *albv20251030.InsertHTTPHeaderInfo {
	if object.IsNull() || object.IsUnknown() {
		return nil
	}

	var model AlbInsertHeaderConfigModel
	d := object.As(ctx, &model, basetypes.ObjectAsOptions{})
	diags.Append(d...)
	if d.HasError() {
		return nil
	}

	return &albv20251030.InsertHTTPHeaderInfo{
		Key:       fwHelper.StringPointerFromValue(model.Key),
		Value:     fwHelper.StringPointerFromValue(model.Value),
		ValueType: fwHelper.StringPointerFromValue(model.ValueType),
	}
}

// albFlattenInsertHeaderConfig converts an SDK insert header configuration.
func albFlattenInsertHeaderConfig(ctx context.Context, item *albv20251030.InsertHTTPHeaderInfo, diags *diag.Diagnostics) types.Object {
	if item == nil {
		return types.ObjectNull(albInsertHeaderConfigAttrTypes)
	}

	value, d := types.ObjectValueFrom(ctx, albInsertHeaderConfigAttrTypes, AlbInsertHeaderConfigModel{
		Key:       fwHelper.StringPointerValueOrNull(item.Key),
		Value:     fwHelper.StringPointerValueOrNull(item.Value),
		ValueType: fwHelper.StringPointerValueOrNull(item.ValueType),
	})
	diags.Append(d...)

	return value
}

// albExpandRedirectConfig converts the `redirect_config` attribute.
func albExpandRedirectConfig(ctx context.Context, object types.Object, diags *diag.Diagnostics) *albv20251030.HTTPRedirectInfo {
	if object.IsNull() || object.IsUnknown() {
		return nil
	}

	var model AlbRedirectConfigModel
	d := object.As(ctx, &model, basetypes.ObjectAsOptions{})
	diags.Append(d...)
	if d.HasError() {
		return nil
	}

	return &albv20251030.HTTPRedirectInfo{
		HttpCode: fwHelper.Int64PointerFromValue(model.HttpCode),
		Host:     fwHelper.StringPointerFromValue(model.Host),
		Path:     fwHelper.StringPointerFromValue(model.Path),
		Port:     fwHelper.StringPointerFromValue(model.Port),
		Protocol: fwHelper.StringPointerFromValue(model.Protocol),
		Query:    fwHelper.StringPointerFromValue(model.Query),
	}
}

// albFlattenRedirectConfig converts an SDK redirect configuration.
func albFlattenRedirectConfig(ctx context.Context, item *albv20251030.HTTPRedirectInfo, diags *diag.Diagnostics) types.Object {
	if item == nil {
		return types.ObjectNull(albRedirectConfigAttrTypes)
	}

	value, d := types.ObjectValueFrom(ctx, albRedirectConfigAttrTypes, AlbRedirectConfigModel{
		HttpCode: fwHelper.Int64ValueOrNull(item.HttpCode),
		Host:     fwHelper.StringPointerValueOrNull(item.Host),
		Path:     fwHelper.StringPointerValueOrNull(item.Path),
		Port:     fwHelper.StringPointerValueOrNull(item.Port),
		Protocol: fwHelper.StringPointerValueOrNull(item.Protocol),
		Query:    fwHelper.StringPointerValueOrNull(item.Query),
	})
	diags.Append(d...)

	return value
}

// albExpandRemoveHeaderConfig converts the `remove_header_config` attribute.
func albExpandRemoveHeaderConfig(ctx context.Context, object types.Object, diags *diag.Diagnostics) *albv20251030.RemoveHTTPHeaderInfo {
	if object.IsNull() || object.IsUnknown() {
		return nil
	}

	var model AlbRemoveHeaderConfigModel
	d := object.As(ctx, &model, basetypes.ObjectAsOptions{})
	diags.Append(d...)
	if d.HasError() {
		return nil
	}

	return &albv20251030.RemoveHTTPHeaderInfo{
		Key: fwHelper.StringPointerFromValue(model.Key),
	}
}

// albFlattenRemoveHeaderConfig converts an SDK remove header configuration.
func albFlattenRemoveHeaderConfig(ctx context.Context, item *albv20251030.RemoveHTTPHeaderInfo, diags *diag.Diagnostics) types.Object {
	if item == nil {
		return types.ObjectNull(albRemoveHeaderConfigAttrTypes)
	}

	value, d := types.ObjectValueFrom(ctx, albRemoveHeaderConfigAttrTypes, AlbRemoveHeaderConfigModel{
		Key: fwHelper.StringPointerValueOrNull(item.Key),
	})
	diags.Append(d...)

	return value
}

// albExpandRewriteConfig converts the `rewrite_config` attribute.
func albExpandRewriteConfig(ctx context.Context, object types.Object, diags *diag.Diagnostics) *albv20251030.HTTPRewriteInfo {
	if object.IsNull() || object.IsUnknown() {
		return nil
	}

	var model AlbRewriteConfigModel
	d := object.As(ctx, &model, basetypes.ObjectAsOptions{})
	diags.Append(d...)
	if d.HasError() {
		return nil
	}

	return &albv20251030.HTTPRewriteInfo{
		Host:  fwHelper.StringPointerFromValue(model.Host),
		Path:  fwHelper.StringPointerFromValue(model.Path),
		Query: fwHelper.StringPointerFromValue(model.Query),
	}
}

// albFlattenRewriteConfig converts an SDK rewrite configuration.
func albFlattenRewriteConfig(ctx context.Context, item *albv20251030.HTTPRewriteInfo, diags *diag.Diagnostics) types.Object {
	if item == nil {
		return types.ObjectNull(albRewriteConfigAttrTypes)
	}

	value, d := types.ObjectValueFrom(ctx, albRewriteConfigAttrTypes, AlbRewriteConfigModel{
		Host:  fwHelper.StringPointerValueOrNull(item.Host),
		Path:  fwHelper.StringPointerValueOrNull(item.Path),
		Query: fwHelper.StringPointerValueOrNull(item.Query),
	})
	diags.Append(d...)

	return value
}

// albExpandTargetGroupConfig converts the `target_group_config` attribute.
func albExpandTargetGroupConfig(ctx context.Context, object types.Object, diags *diag.Diagnostics) *albv20251030.TargetGroupConfig {
	if object.IsNull() || object.IsUnknown() {
		return nil
	}

	var model AlbTargetGroupConfigModel
	d := object.As(ctx, &model, basetypes.ObjectAsOptions{})
	diags.Append(d...)
	if d.HasError() {
		return nil
	}

	config := &albv20251030.TargetGroupConfig{}
	if !model.TargetGroups.IsNull() && !model.TargetGroups.IsUnknown() {
		var groups []AlbTargetGroupTupleModel
		groupDiags := model.TargetGroups.ElementsAs(ctx, &groups, false)
		diags.Append(groupDiags...)
		if groupDiags.HasError() {
			return nil
		}

		targetGroups := make([]*albv20251030.TargetGroupTuple, 0, len(groups))
		for _, group := range groups {
			targetGroups = append(targetGroups, &albv20251030.TargetGroupTuple{
				TargetGroupId: fwHelper.StringPointerFromValue(group.TargetGroupId),
				Weight:        fwHelper.Int64PointerFromValue(group.Weight),
			})
		}

		config.TargetGroups = targetGroups
	}

	if !model.TargetGroupStickySession.IsNull() && !model.TargetGroupStickySession.IsUnknown() {
		var sticky AlbTargetGroupStickySessionModel
		stickyDiags := model.TargetGroupStickySession.As(ctx, &sticky, basetypes.ObjectAsOptions{})
		diags.Append(stickyDiags...)
		if stickyDiags.HasError() {
			return nil
		}

		config.TargetGroupStickySession = &albv20251030.TargetGroupStickySession{
			Enabled: fwHelper.BoolPointerFromValue(sticky.Enabled),
			Timeout: fwHelper.Int64PointerFromValue(sticky.Timeout),
		}
	}

	return config
}

// albFlattenTargetGroupConfig converts an SDK target group configuration.
func albFlattenTargetGroupConfig(ctx context.Context, item *albv20251030.TargetGroupConfig, diags *diag.Diagnostics) types.Object {
	if item == nil {
		return types.ObjectNull(albTargetGroupConfigAttrTypes)
	}

	model := AlbTargetGroupConfigModel{
		TargetGroups:             types.ListNull(albTargetGroupTupleObjectType),
		TargetGroupStickySession: types.ObjectNull(albTargetGroupStickyAttrTypes),
	}

	if item.TargetGroups != nil {
		groups := make([]AlbTargetGroupTupleModel, 0, len(item.TargetGroups))
		for _, group := range item.TargetGroups {
			if group == nil {
				continue
			}

			groups = append(groups, AlbTargetGroupTupleModel{
				TargetGroupId: fwHelper.StringPointerValueOrNull(group.TargetGroupId),
				Weight:        fwHelper.Int64ValueOrNull(group.Weight),
			})
		}

		list, listDiags := types.ListValueFrom(ctx, albTargetGroupTupleObjectType, groups)
		diags.Append(listDiags...)
		if listDiags.HasError() {
			return types.ObjectNull(albTargetGroupConfigAttrTypes)
		}

		model.TargetGroups = list
	}

	if item.TargetGroupStickySession != nil {
		sticky, stickyDiags := types.ObjectValueFrom(ctx, albTargetGroupStickyAttrTypes, AlbTargetGroupStickySessionModel{
			Enabled: fwHelper.BoolValueOrNull(item.TargetGroupStickySession.Enabled),
			Timeout: fwHelper.Int64ValueOrNull(item.TargetGroupStickySession.Timeout),
		})
		diags.Append(stickyDiags...)
		if stickyDiags.HasError() {
			return types.ObjectNull(albTargetGroupConfigAttrTypes)
		}

		model.TargetGroupStickySession = sticky
	}

	value, d := types.ObjectValueFrom(ctx, albTargetGroupConfigAttrTypes, model)
	diags.Append(d...)

	return value
}

// albExpandRuleConditions converts the `conditions` attribute into SDK rule
// conditions.
func albExpandRuleConditions(ctx context.Context, list types.List, diags *diag.Diagnostics) []*albv20251030.RuleCondition {
	if list.IsNull() || list.IsUnknown() {
		return nil
	}

	var models []AlbRuleConditionModel
	d := list.ElementsAs(ctx, &models, false)
	diags.Append(d...)
	if d.HasError() {
		return nil
	}

	conditions := make([]*albv20251030.RuleCondition, 0, len(models))
	for _, model := range models {
		condition := &albv20251030.RuleCondition{
			Type:              fwHelper.StringPointerFromValue(model.Type),
			CookieConfig:      albExpandCookieConfig(ctx, model.CookieConfig, diags),
			HeaderConfig:      albExpandHeaderConfig(ctx, model.HeaderConfig, diags),
			HostConfig:        fwHelper.PointersFromList[string](ctx, model.HostConfig, diags),
			MethodConfig:      fwHelper.PointersFromList[string](ctx, model.MethodConfig, diags),
			PathConfig:        fwHelper.PointersFromList[string](ctx, model.PathConfig, diags),
			QueryStringConfig: albExpandQueryStringConfig(ctx, model.QueryStringConfig, diags),
			SourceIpConfig:    fwHelper.PointersFromList[string](ctx, model.SourceIpConfig, diags),
		}
		if diags.HasError() {
			return nil
		}

		conditions = append(conditions, condition)
	}

	return conditions
}

// albFlattenRuleConditions converts SDK rule conditions into the `conditions`
// attribute.
func albFlattenRuleConditions(ctx context.Context, items []*albv20251030.RuleCondition, diags *diag.Diagnostics) types.List {
	models := make([]AlbRuleConditionModel, 0, len(items))
	for _, item := range items {
		if item == nil {
			continue
		}

		models = append(models, AlbRuleConditionModel{
			Type:              fwHelper.StringPointerValueOrNull(item.Type),
			CookieConfig:      albFlattenCookieConfig(ctx, item.CookieConfig, diags),
			HeaderConfig:      albFlattenHeaderConfig(ctx, item.HeaderConfig, diags),
			HostConfig:        fwHelper.ListFromPointers[string](ctx, types.StringType, item.HostConfig, diags),
			MethodConfig:      fwHelper.ListFromPointers[string](ctx, types.StringType, item.MethodConfig, diags),
			PathConfig:        fwHelper.ListFromPointers[string](ctx, types.StringType, item.PathConfig, diags),
			QueryStringConfig: albFlattenQueryStringConfig(ctx, item.QueryStringConfig, diags),
			SourceIpConfig:    fwHelper.ListFromPointers[string](ctx, types.StringType, item.SourceIpConfig, diags),
		})
	}

	list, d := types.ListValueFrom(ctx, albRuleConditionObjectType, models)
	diags.Append(d...)

	return list
}

// albExpandCookieConfig converts the `cookie_config` attribute.
func albExpandCookieConfig(ctx context.Context, list types.List, diags *diag.Diagnostics) []*albv20251030.HTTPCookieInfo {
	if list.IsNull() || list.IsUnknown() {
		return nil
	}

	var models []AlbCookieConfigModel
	d := list.ElementsAs(ctx, &models, false)
	diags.Append(d...)
	if d.HasError() {
		return nil
	}

	items := make([]*albv20251030.HTTPCookieInfo, 0, len(models))
	for _, model := range models {
		items = append(items, &albv20251030.HTTPCookieInfo{
			Key:   fwHelper.StringPointerFromValue(model.Key),
			Value: fwHelper.StringPointerFromValue(model.Value),
		})
	}

	return items
}

// albFlattenCookieConfig converts an SDK cookie configuration list.
func albFlattenCookieConfig(ctx context.Context, items []*albv20251030.HTTPCookieInfo, diags *diag.Diagnostics) types.List {
	if items == nil {
		return types.ListNull(albCookieConfigObjectType)
	}

	models := make([]AlbCookieConfigModel, 0, len(items))
	for _, item := range items {
		if item == nil {
			continue
		}

		models = append(models, AlbCookieConfigModel{
			Key:   fwHelper.StringPointerValueOrNull(item.Key),
			Value: fwHelper.StringPointerValueOrNull(item.Value),
		})
	}

	list, d := types.ListValueFrom(ctx, albCookieConfigObjectType, models)
	diags.Append(d...)

	return list
}

// albExpandHeaderConfig converts the `header_config` attribute.
func albExpandHeaderConfig(ctx context.Context, object types.Object, diags *diag.Diagnostics) *albv20251030.HTTPHeaderInfo {
	if object.IsNull() || object.IsUnknown() {
		return nil
	}

	var model AlbHeaderConfigModel
	d := object.As(ctx, &model, basetypes.ObjectAsOptions{})
	diags.Append(d...)
	if d.HasError() {
		return nil
	}

	return &albv20251030.HTTPHeaderInfo{
		Key:    fwHelper.StringPointerFromValue(model.Key),
		Values: fwHelper.PointersFromList[string](ctx, model.Values, diags),
	}
}

// albFlattenHeaderConfig converts an SDK header configuration.
func albFlattenHeaderConfig(ctx context.Context, item *albv20251030.HTTPHeaderInfo, diags *diag.Diagnostics) types.Object {
	if item == nil {
		return types.ObjectNull(albHeaderConfigAttrTypes)
	}

	value, d := types.ObjectValueFrom(ctx, albHeaderConfigAttrTypes, AlbHeaderConfigModel{
		Key:    fwHelper.StringPointerValueOrNull(item.Key),
		Values: fwHelper.ListFromPointers[string](ctx, types.StringType, item.Values, diags),
	})
	diags.Append(d...)

	return value
}

// albExpandQueryStringConfig converts the `query_string_config` attribute.
func albExpandQueryStringConfig(ctx context.Context, list types.List, diags *diag.Diagnostics) []*albv20251030.HTTPQueryStringInfo {
	if list.IsNull() || list.IsUnknown() {
		return nil
	}

	var models []AlbQueryStringConfigModel
	d := list.ElementsAs(ctx, &models, false)
	diags.Append(d...)
	if d.HasError() {
		return nil
	}

	items := make([]*albv20251030.HTTPQueryStringInfo, 0, len(models))
	for _, model := range models {
		items = append(items, &albv20251030.HTTPQueryStringInfo{
			Key:   fwHelper.StringPointerFromValue(model.Key),
			Value: fwHelper.StringPointerFromValue(model.Value),
		})
	}

	return items
}

// albFlattenQueryStringConfig converts an SDK query string configuration list.
func albFlattenQueryStringConfig(ctx context.Context, items []*albv20251030.HTTPQueryStringInfo, diags *diag.Diagnostics) types.List {
	if items == nil {
		return types.ListNull(albQueryStringConfigObjectType)
	}

	models := make([]AlbQueryStringConfigModel, 0, len(items))
	for _, item := range items {
		if item == nil {
			continue
		}

		models = append(models, AlbQueryStringConfigModel{
			Key:   fwHelper.StringPointerValueOrNull(item.Key),
			Value: fwHelper.StringPointerValueOrNull(item.Value),
		})
	}

	list, d := types.ListValueFrom(ctx, albQueryStringConfigObjectType, models)
	diags.Append(d...)

	return list
}

// albExpandTags converts the `tags` attribute.
func albExpandTags(ctx context.Context, list types.List, diags *diag.Diagnostics) []*albv20251030.TagInfo {
	if list.IsNull() || list.IsUnknown() {
		return nil
	}

	var models []AlbTagModel
	d := list.ElementsAs(ctx, &models, false)
	diags.Append(d...)
	if d.HasError() {
		return nil
	}

	items := make([]*albv20251030.TagInfo, 0, len(models))
	for _, model := range models {
		items = append(items, &albv20251030.TagInfo{
			TagKey:   fwHelper.StringPointerFromValue(model.TagKey),
			TagValue: fwHelper.StringPointerFromValue(model.TagValue),
		})
	}

	return items
}

// albFlattenTags converts SDK tags into the `tags` attribute.
func albFlattenTags(ctx context.Context, items []*albv20251030.TagInfo, diags *diag.Diagnostics) types.List {
	if items == nil {
		return types.ListNull(albTagObjectType)
	}

	models := make([]AlbTagModel, 0, len(items))
	for _, item := range items {
		if item == nil {
			continue
		}

		models = append(models, AlbTagModel{
			TagKey:   fwHelper.StringPointerValueOrNull(item.TagKey),
			TagValue: fwHelper.StringPointerValueOrNull(item.TagValue),
		})
	}

	list, d := types.ListValueFrom(ctx, albTagObjectType, models)
	diags.Append(d...)

	return list
}
