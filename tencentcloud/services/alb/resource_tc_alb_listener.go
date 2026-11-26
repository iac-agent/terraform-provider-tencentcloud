package alb

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	albv20251030 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/alb/v20251030"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/fw"
	fwHelper "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/helper"
)

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
type AlbListenerResource struct {
	fw.ResourceWithConfigure
	fw.WithImportByID
}

// AlbListenerResourceModel maps the schema attributes.
type AlbListenerResourceModel struct {
	Id             types.String `tfsdk:"id"`
	ListenerId     types.String `tfsdk:"listener_id"`
	LoadBalancerId types.String `tfsdk:"load_balancer_id"`
	ClientToken    types.String `tfsdk:"client_token"`
	DryRun         types.Bool   `tfsdk:"dry_run"`
	Rules          types.List   `tfsdk:"rules"`
}

type albRuleModel struct {
	Actions    types.List   `tfsdk:"actions"`
	Conditions types.List   `tfsdk:"conditions"`
	Priority   types.Int64  `tfsdk:"priority"`
	Direction  types.String `tfsdk:"direction"`
	RuleName   types.String `tfsdk:"rule_name"`
	Tags       types.List   `tfsdk:"tags"`
	RuleId     types.String `tfsdk:"rule_id"`
	Status     types.String `tfsdk:"status"`
	CreateTime types.String `tfsdk:"create_time"`
	ModifyTime types.String `tfsdk:"modify_time"`
}

type albRuleActionModel struct {
	Order               types.Int64  `tfsdk:"order"`
	Type                types.String `tfsdk:"type"`
	FixedResponseConfig types.Object `tfsdk:"fixed_response_config"`
	InsertHeaderConfig  types.Object `tfsdk:"insert_header_config"`
	RedirectConfig      types.Object `tfsdk:"redirect_config"`
	RemoveHeaderConfig  types.Object `tfsdk:"remove_header_config"`
	RewriteConfig       types.Object `tfsdk:"rewrite_config"`
	TargetGroupConfig   types.Object `tfsdk:"target_group_config"`
}

type albFixedResponseConfigModel struct {
	HttpCode    types.Int64  `tfsdk:"http_code"`
	Content     types.String `tfsdk:"content"`
	ContentType types.String `tfsdk:"content_type"`
}

type albInsertHeaderConfigModel struct {
	Key       types.String `tfsdk:"key"`
	Value     types.String `tfsdk:"value"`
	ValueType types.String `tfsdk:"value_type"`
}

type albRedirectConfigModel struct {
	HttpCode types.Int64  `tfsdk:"http_code"`
	Host     types.String `tfsdk:"host"`
	Path     types.String `tfsdk:"path"`
	Port     types.String `tfsdk:"port"`
	Protocol types.String `tfsdk:"protocol"`
	Query    types.String `tfsdk:"query"`
}

type albRemoveHeaderConfigModel struct {
	Key types.String `tfsdk:"key"`
}

type albRewriteConfigModel struct {
	Host  types.String `tfsdk:"host"`
	Path  types.String `tfsdk:"path"`
	Query types.String `tfsdk:"query"`
}

type albTargetGroupConfigModel struct {
	TargetGroups             types.List   `tfsdk:"target_groups"`
	TargetGroupStickySession types.Object `tfsdk:"target_group_sticky_session"`
}

type albTargetGroupTupleModel struct {
	TargetGroupId types.String `tfsdk:"target_group_id"`
	Weight        types.Int64  `tfsdk:"weight"`
}

type albTargetGroupStickySessionModel struct {
	Enabled types.Bool  `tfsdk:"enabled"`
	Timeout types.Int64 `tfsdk:"timeout"`
}

type albRuleConditionModel struct {
	Type              types.String `tfsdk:"type"`
	CookieConfig      types.List   `tfsdk:"cookie_config"`
	HeaderConfig      types.Object `tfsdk:"header_config"`
	HostConfig        types.List   `tfsdk:"host_config"`
	MethodConfig      types.List   `tfsdk:"method_config"`
	PathConfig        types.List   `tfsdk:"path_config"`
	QueryStringConfig types.List   `tfsdk:"query_string_config"`
	SourceIpConfig    types.List   `tfsdk:"source_ip_config"`
}

type albCookieConfigModel struct {
	Key   types.String `tfsdk:"key"`
	Value types.String `tfsdk:"value"`
}

type albHeaderConfigModel struct {
	Key    types.String `tfsdk:"key"`
	Values types.List   `tfsdk:"values"`
}

type albQueryStringConfigModel struct {
	Key   types.String `tfsdk:"key"`
	Value types.String `tfsdk:"value"`
}

type albTagModel struct {
	TagKey   types.String `tfsdk:"tag_key"`
	TagValue types.String `tfsdk:"tag_value"`
}

var albListenerFixedResponseConfigAttributes = map[string]schema.Attribute{
	"http_code": schema.Int64Attribute{
		Required:    true,
		Description: "The HTTP response code returned, supporting 2xx, 4xx, 5xx.",
	},
	"content": schema.StringAttribute{
		Optional:    true,
		Description: "The fixed content returned. Only ASCII characters are supported, up to 1 KB.",
	},
	"content_type": schema.StringAttribute{
		Optional:    true,
		Description: "The format of the fixed content. Valid values: text/plain, text/css, text/html, application/javascript, or application/json.",
	},
}

var albListenerInsertHeaderConfigAttributes = map[string]schema.Attribute{
	"key": schema.StringAttribute{
		Required:    true,
		Description: "The key of the HTTP Header to insert, 1 to 40 characters long, supporting a-z A-Z 0-9 - _ . Chinese characters and Cookie, Host, Content-Length, Connection, Upgrade, transfer-encoding, keep-alive, te, authority, x-forwarded-* are not supported.",
	},
	"value": schema.StringAttribute{
		Required:    true,
		Description: "The value of the HTTP Header. When ValueType is SystemDefined, valid values are ClientPort (client port), ClientIp (client IP address), Protocol (client request protocol), CLBPort (load balancer listening port). When ValueType is UserDefined, 1 to 128 printable characters.",
	},
	"value_type": schema.StringAttribute{
		Required:    true,
		Description: "The type of the HTTP Header value. Valid values: SystemDefined (system-defined header), UserDefined (user-defined header), ReferenceHeader (reference a header from the request).",
	},
}

var albListenerRedirectConfigAttributes = map[string]schema.Attribute{
	"http_code": schema.Int64Attribute{
		Required:    true,
		Description: "The HTTP code for redirection. Supported values: 301, 302, 303, 307, 308.",
	},
	"host": schema.StringAttribute{
		Optional:    true,
		Description: "The redirection host, default ${host}. 3 to 128 characters long, supporting a-z 0-9 _ . -.",
	},
	"path": schema.StringAttribute{
		Optional:    true,
		Description: "The redirection path, default ${path}. 1 to 128 characters long, supporting a-z A-Z 0-9 ? = _ . - / : .",
	},
	"port": schema.StringAttribute{
		Optional:    true,
		Description: "The redirection port, default ${port}. Valid values: 1 to 65535.",
	},
	"protocol": schema.StringAttribute{
		Optional:    true,
		Description: "The redirection protocol. Valid values: HTTP, HTTPS, default ${protocol}.",
	},
	"query": schema.StringAttribute{
		Optional:    true,
		Description: "The redirection query string, default ${query}. 1 to 128 printable characters, excluding #[ ]{}|<> & and spaces.",
	},
}

var albListenerRemoveHeaderConfigAttributes = map[string]schema.Attribute{
	"key": schema.StringAttribute{
		Required:    true,
		Description: "The key of the HTTP Header to remove, 1 to 40 characters long, supporting a-z A-Z 0-9 - _ . Chinese characters and Host, Cookie are not supported.",
	},
}

var albListenerRewriteConfigAttributes = map[string]schema.Attribute{
	"host": schema.StringAttribute{
		Optional:    true,
		Description: "The rewrite host, default ${host}. 3 to 128 characters long, supporting a-z 0-9 _ . -.",
	},
	"path": schema.StringAttribute{
		Optional:    true,
		Description: "The rewrite path, default ${path}. 1 to 128 characters long, supporting a-z A-Z 0-9 ? = _ . - / : .",
	},
	"query": schema.StringAttribute{
		Optional:    true,
		Description: "The rewrite query string, default ${query}. 1 to 128 printable characters, excluding #[ ]{}|<> & and spaces.",
	},
}

var albListenerTargetGroupStickySessionAttributes = map[string]schema.Attribute{
	"enabled": schema.BoolAttribute{
		Required:    true,
		Description: "Whether to enable sticky sessions between target groups. Disabled by default.",
	},
	"timeout": schema.Int64Attribute{
		Optional:    true,
		Description: "The timeout in seconds, ranging from 1 to 86400, default 1000.",
	},
}

var albListenerTargetGroupTupleAttributes = map[string]schema.Attribute{
	"target_group_id": schema.StringAttribute{
		Required:    true,
		Description: "The target group ID, in the format lbtg- followed by 8 alphanumeric characters.",
	},
	"weight": schema.Int64Attribute{
		Optional:    true,
		Description: "The weight, ranging from 0 to 100, default 10.",
	},
}

var albListenerTargetGroupConfigAttributes = map[string]schema.Attribute{
	"target_groups": schema.ListNestedAttribute{
		Required:    true,
		Description: "The target group list.",
		NestedObject: schema.NestedAttributeObject{
			Attributes: albListenerTargetGroupTupleAttributes,
		},
	},
	"target_group_sticky_session": schema.SingleNestedAttribute{
		Optional:    true,
		Description: "Sticky sessions between target groups.",
		Attributes:  albListenerTargetGroupStickySessionAttributes,
	},
}

var albListenerRuleActionAttributes = map[string]schema.Attribute{
	"order": schema.Int64Attribute{
		Required:    true,
		Description: "The execution order of the forwarding action. Orders cannot repeat; actions are executed in ascending order. Valid values: 1 to 50000.",
	},
	"type": schema.StringAttribute{
		Required:    true,
		Description: "The forwarding action type. Valid values: TargetGroup (forward to a target group), Redirect (redirect), FixedResponse (return a fixed response), Rewrite (rewrite), InsertHeader (write an HTTP header), RemoveHeader (remove an HTTP header). A forwarding action must include TargetGroup, Redirect, FixedResponse or Rewrite.",
	},
	"fixed_response_config": schema.SingleNestedAttribute{
		Optional:    true,
		Description: "The fixed response content configuration.",
		Attributes:  albListenerFixedResponseConfigAttributes,
	},
	"insert_header_config": schema.SingleNestedAttribute{
		Optional:    true,
		Description: "The inserted HTTP Header configuration.",
		Attributes:  albListenerInsertHeaderConfigAttributes,
	},
	"redirect_config": schema.SingleNestedAttribute{
		Optional:    true,
		Description: "The redirection configuration. Except for HttpCode, the other configurations cannot all use default values.",
		Attributes:  albListenerRedirectConfigAttributes,
	},
	"remove_header_config": schema.SingleNestedAttribute{
		Optional:    true,
		Description: "The removed HTTP Header configuration.",
		Attributes:  albListenerRemoveHeaderConfigAttributes,
	},
	"rewrite_config": schema.SingleNestedAttribute{
		Optional:    true,
		Description: "The rewrite configuration.",
		Attributes:  albListenerRewriteConfigAttributes,
	},
	"target_group_config": schema.SingleNestedAttribute{
		Optional:    true,
		Description: "The forwarding target group configuration.",
		Attributes:  albListenerTargetGroupConfigAttributes,
	},
}

var albListenerCookieConfigAttributes = map[string]schema.Attribute{
	"key": schema.StringAttribute{
		Required:    true,
		Description: "The key of the Cookie, 1 to 64 characters long, supporting letters, digits and underscores.",
	},
	"value": schema.StringAttribute{
		Required:    true,
		Description: "The value of the Cookie, 1 to 128 printable characters.",
	},
}

var albListenerHeaderConfigAttributes = map[string]schema.Attribute{
	"key": schema.StringAttribute{
		Required:    true,
		Description: "The key of the HTTP Header, 1 to 40 characters long, supporting a-z A-Z 0-9 - _ . Chinese characters, Host and Cookie are not supported.",
	},
	"values": schema.ListAttribute{
		Required:    true,
		ElementType: types.StringType,
		Description: "The values of the HTTP Header, 1 to 128 printable characters. Double quotes are not supported; leading and trailing spaces are not allowed, and a trailing backslash is not allowed.",
	},
}

var albListenerQueryStringConfigAttributes = map[string]schema.Attribute{
	"key": schema.StringAttribute{
		Required:    true,
		Description: "The key of the query string, 1 to 16 characters long. Printable characters are supported, excluding spaces and #[]{}\\|<>&. * matches multiple characters and ? matches a single character.",
	},
	"value": schema.StringAttribute{
		Required:    true,
		Description: "The value of the query string, 1 to 128 characters long. Printable characters are supported, excluding spaces and #[]{}\\|<>&. * matches multiple characters and ? matches a single character.",
	},
}

var albListenerRuleConditionAttributes = map[string]schema.Attribute{
	"type": schema.StringAttribute{
		Required:    true,
		Description: "The forwarding condition type. Valid values: Host, Path, Header, QueryString, Method, Cookie, SourceIp.",
	},
	"cookie_config": schema.ListNestedAttribute{
		Optional:    true,
		Description: "The Cookie configuration.",
		NestedObject: schema.NestedAttributeObject{
			Attributes: albListenerCookieConfigAttributes,
		},
	},
	"header_config": schema.SingleNestedAttribute{
		Optional:    true,
		Description: "The HTTP Header configuration.",
		Attributes:  albListenerHeaderConfigAttributes,
	},
	"host_config": schema.ListAttribute{
		Optional:    true,
		ElementType: types.StringType,
		Description: "The host name. A host configuration can appear only once in a rule, 3 to 128 characters long, supporting exact, regex and wildcard matching. It cannot start or end with a period (.) or underscore (_).",
	},
	"method_config": schema.ListAttribute{
		Optional:    true,
		ElementType: types.StringType,
		Description: "The request method. Valid values: HEAD, GET, POST, OPTIONS, PUT, PATCH, DELETE.",
	},
	"path_config": schema.ListAttribute{
		Optional:    true,
		ElementType: types.StringType,
		Description: "The forwarding path, 1 to 128 characters long, supporting exact, regex and wildcard matching.",
	},
	"query_string_config": schema.ListNestedAttribute{
		Optional:    true,
		Description: "The query string configuration.",
		NestedObject: schema.NestedAttributeObject{
			Attributes: albListenerQueryStringConfigAttributes,
		},
	},
	"source_ip_config": schema.ListAttribute{
		Optional:    true,
		ElementType: types.StringType,
		Description: "The source IP matching configuration in CIDR format, for example x.x.x.x/32 or x.x.x.x/24.",
	},
}

var albListenerTagAttributes = map[string]schema.Attribute{
	"tag_key": schema.StringAttribute{
		Required:    true,
		Description: "The key of the tag.",
	},
	"tag_value": schema.StringAttribute{
		Required:    true,
		Description: "The value of the tag.",
	},
}

var albListenerRuleAttributes = map[string]schema.Attribute{
	"actions": schema.ListNestedAttribute{
		Required:    true,
		Description: "The forwarding rule action list.",
		NestedObject: schema.NestedAttributeObject{
			Attributes: albListenerRuleActionAttributes,
		},
	},
	"conditions": schema.ListNestedAttribute{
		Optional:    true,
		Description: "The forwarding rule condition list.",
		NestedObject: schema.NestedAttributeObject{
			Attributes: albListenerRuleConditionAttributes,
		},
	},
	"priority": schema.Int64Attribute{
		Optional:    true,
		Computed:    true,
		Description: "The priority. A smaller value means a higher priority; values cannot repeat. Valid values: 1 to 10000.",
	},
	"direction": schema.StringAttribute{
		Optional:    true,
		Computed:    true,
		Description: "The direction of the forwarding rule. Request: from the client to the load balancer; Response: from the backend server to the load balancer. Default Request.",
	},
	"rule_name": schema.StringAttribute{
		Optional:    true,
		Description: "The forwarding rule name, 1 to 255 characters long, supporting digits, upper and lower case letters, Chinese characters, periods (.), underscores (_) and hyphens (-).",
	},
	"tags": schema.ListNestedAttribute{
		Optional:    true,
		Description: "The tags.",
		NestedObject: schema.NestedAttributeObject{
			Attributes: albListenerTagAttributes,
		},
	},
	"rule_id": schema.StringAttribute{
		Computed:    true,
		Description: "The forwarding rule ID, in the format rule- followed by 8 alphanumeric characters.",
	},
	"status": schema.StringAttribute{
		Computed:    true,
		Description: "The forwarding rule status. Provisioning: creating; Active: running; Configuring: configuring.",
	},
	"create_time": schema.StringAttribute{
		Computed:    true,
		Description: "The creation time.",
	},
	"modify_time": schema.StringAttribute{
		Computed:    true,
		Description: "The last modification time.",
	},
}

func (r *AlbListenerResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "tencentcloud_alb_listener"
}

func (r *AlbListenerResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a resource to manage a listener forwarding rule.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "ID of the resource, the ID of the first forwarding rule created.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"listener_id": schema.StringAttribute{
				Required:    true,
				Description: "The listener ID, in the format lst- followed by 8 alphanumeric characters.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"load_balancer_id": schema.StringAttribute{
				Required:    true,
				Description: "The load balancer instance ID, in the format alb- followed by 8 alphanumeric characters.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"client_token": schema.StringAttribute{
				Optional:    true,
				Description: "The client token used to ensure request idempotency. Generate a unique value on the client so different requests carry different tokens. Only ASCII characters are supported. If unspecified, the API uses the request ID as the client token.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"dry_run": schema.BoolAttribute{
				Optional:    true,
				Description: "Whether to only perform a pre-check for this request.",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
			},
			"rules": schema.ListNestedAttribute{
				Required:    true,
				Description: "The forwarding rule list. The resource manages a single rule; a list with exactly one rule is expected.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: albListenerRuleAttributes,
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

	createRulesRequest := albv20251030.NewCreateRulesRequest()
	createRulesRequest.ListenerId = plan.ListenerId.ValueStringPointer()
	createRulesRequest.LoadBalancerId = plan.LoadBalancerId.ValueStringPointer()
	if !plan.ClientToken.IsNull() && !plan.ClientToken.IsUnknown() {
		createRulesRequest.ClientToken = plan.ClientToken.ValueStringPointer()
	}
	if !plan.DryRun.IsNull() && !plan.DryRun.IsUnknown() {
		createRulesRequest.DryRun = plan.DryRun.ValueBoolPointer()
	}

	resp.Diagnostics.Append(expandAlbListenerRules(ctx, plan.Rules, &createRulesRequest.Rules)...)
	if resp.Diagnostics.HasError() {
		return
	}

	service := NewAlbService(r.Client())

	ruleIds, err := service.CreateRules(ctx, createRulesRequest)
	if err != nil {
		resp.Diagnostics.AddError("Error creating ALB listener forwarding rule", err.Error())
		return
	}

	if len(ruleIds) == 0 || ruleIds[0] == nil || *ruleIds[0] == "" {
		resp.Diagnostics.AddError("Error creating ALB listener forwarding rule", "CreateRules returned an empty rule id.")
		return
	}

	ruleId := *ruleIds[0]
	plan.Id = types.StringValue(ruleId)
	normalizeAlbListenerComputedUnknowns(ctx, &plan, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	describeRulesRequest := albv20251030.NewDescribeRulesRequest()
	describeRulesRequest.ListenerId = plan.ListenerId.ValueStringPointer()
	describeRulesRequest.LoadBalancerId = plan.LoadBalancerId.ValueStringPointer()
	describeRulesRequest.RuleIds = []*string{fwHelper.String(ruleId)}

	rules, err := service.DescribeRules(ctx, describeRulesRequest)
	if err != nil {
		resp.Diagnostics.AddError("Error reading ALB listener forwarding rule after create", err.Error())
		return
	}

	item := findAlbListenerRule(rules, ruleId)
	if item == nil {
		resp.Diagnostics.AddError("ALB listener forwarding rule not readable after create",
			fmt.Sprintf("rule %s was created but cannot be read back", ruleId))
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

	ruleId := state.Id.ValueString()

	describeRulesRequest := albv20251030.NewDescribeRulesRequest()
	describeRulesRequest.ListenerId = state.ListenerId.ValueStringPointer()
	describeRulesRequest.LoadBalancerId = state.LoadBalancerId.ValueStringPointer()
	describeRulesRequest.RuleIds = []*string{fwHelper.String(ruleId)}

	service := NewAlbService(r.Client())

	rules, err := service.DescribeRules(ctx, describeRulesRequest)
	if err != nil {
		resp.Diagnostics.AddError("Error reading ALB listener forwarding rule", err.Error())
		return
	}

	item := findAlbListenerRule(rules, ruleId)
	if item == nil {
		resp.Diagnostics.AddWarning("Resource not found",
			fmt.Sprintf("ALB listener forwarding rule %s not found, removing from state.", ruleId))
		resp.State.RemoveResource(ctx)
		return
	}

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

	ruleId := state.Id.ValueString()

	service := NewAlbService(r.Client())

	// TODO(human): ModifyRulesAttributes replaces all attributes of the rule via its Rules list; decide whether partial rules should be merged with the current state rather than sending the planned list as-is.
	if !plan.Rules.Equal(state.Rules) {
		modifyRulesAttributesRequest := albv20251030.NewModifyRulesAttributesRequest()
		modifyRulesAttributesRequest.ListenerId = plan.ListenerId.ValueStringPointer()
		modifyRulesAttributesRequest.LoadBalancerId = plan.LoadBalancerId.ValueStringPointer()
		if !plan.DryRun.IsNull() && !plan.DryRun.IsUnknown() {
			modifyRulesAttributesRequest.DryRun = plan.DryRun.ValueBoolPointer()
		}

		resp.Diagnostics.Append(expandAlbListenerRulesForModify(ctx, plan.Rules, ruleId, &modifyRulesAttributesRequest.Rules)...)
		if resp.Diagnostics.HasError() {
			return
		}

		if err := service.ModifyRulesAttributes(ctx, modifyRulesAttributesRequest); err != nil {
			resp.Diagnostics.AddError("Error updating ALB listener forwarding rule", err.Error())
			return
		}
	}

	describeRulesRequest := albv20251030.NewDescribeRulesRequest()
	describeRulesRequest.ListenerId = plan.ListenerId.ValueStringPointer()
	describeRulesRequest.LoadBalancerId = plan.LoadBalancerId.ValueStringPointer()
	describeRulesRequest.RuleIds = []*string{fwHelper.String(ruleId)}

	rules, err := service.DescribeRules(ctx, describeRulesRequest)
	if err != nil {
		resp.Diagnostics.AddError("Error reading ALB listener forwarding rule after update", err.Error())
		return
	}

	item := findAlbListenerRule(rules, ruleId)
	if item == nil {
		resp.Diagnostics.AddError("ALB listener forwarding rule not readable after update",
			fmt.Sprintf("rule %s was updated but cannot be read back", ruleId))
		return
	}

	plan.Id = state.Id
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

	deleteRulesRequest := albv20251030.NewDeleteRulesRequest()
	deleteRulesRequest.ListenerId = state.ListenerId.ValueStringPointer()
	deleteRulesRequest.LoadBalancerId = state.LoadBalancerId.ValueStringPointer()
	deleteRulesRequest.RuleIds = []*string{fwHelper.String(state.Id.ValueString())}
	if !state.DryRun.IsNull() && !state.DryRun.IsUnknown() {
		deleteRulesRequest.DryRun = state.DryRun.ValueBoolPointer()
	}

	service := NewAlbService(r.Client())

	if err := service.DeleteRules(ctx, deleteRulesRequest); err != nil {
		resp.Diagnostics.AddError("Error deleting ALB listener forwarding rule", err.Error())
	}
}

// normalizeAlbListenerComputedUnknowns rewrites unknown computed attributes to
// null so the model can be persisted into state before the read-back fills
// them with cloud-side values.
func normalizeAlbListenerComputedUnknowns(ctx context.Context, m *AlbListenerResourceModel, diags *diag.Diagnostics) {
	if m.Rules.IsNull() || m.Rules.IsUnknown() {
		return
	}

	var rules []albRuleModel
	diags.Append(m.Rules.ElementsAs(ctx, &rules, false)...)
	if diags.HasError() {
		return
	}

	for i := range rules {
		if rules[i].Priority.IsUnknown() {
			rules[i].Priority = types.Int64Null()
		}
		if rules[i].Direction.IsUnknown() {
			rules[i].Direction = types.StringNull()
		}
		if rules[i].RuleId.IsUnknown() {
			rules[i].RuleId = types.StringNull()
		}
		if rules[i].Status.IsUnknown() {
			rules[i].Status = types.StringNull()
		}
		if rules[i].CreateTime.IsUnknown() {
			rules[i].CreateTime = types.StringNull()
		}
		if rules[i].ModifyTime.IsUnknown() {
			rules[i].ModifyTime = types.StringNull()
		}
	}

	list, d := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: albListenerRuleAttrTypes()}, rules)
	diags.Append(d...)
	if !diags.HasError() {
		m.Rules = list
	}
}

// albListenerRuleAttrTypes returns the attr types of the rules list elements.
func albListenerRuleAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"actions":     types.ListType{ElemType: types.ObjectType{AttrTypes: albListenerRuleActionAttrTypes()}},
		"conditions":  types.ListType{ElemType: types.ObjectType{AttrTypes: albListenerRuleConditionAttrTypes()}},
		"priority":    types.Int64Type,
		"direction":   types.StringType,
		"rule_name":   types.StringType,
		"tags":        types.ListType{ElemType: types.ObjectType{AttrTypes: albListenerTagAttrTypes()}},
		"rule_id":     types.StringType,
		"status":      types.StringType,
		"create_time": types.StringType,
		"modify_time": types.StringType,
	}
}

func albListenerRuleActionAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"order":                 types.Int64Type,
		"type":                  types.StringType,
		"fixed_response_config": types.ObjectType{AttrTypes: albListenerFixedResponseConfigAttrTypes()},
		"insert_header_config":  types.ObjectType{AttrTypes: albListenerInsertHeaderConfigAttrTypes()},
		"redirect_config":       types.ObjectType{AttrTypes: albListenerRedirectConfigAttrTypes()},
		"remove_header_config":  types.ObjectType{AttrTypes: albListenerRemoveHeaderConfigAttrTypes()},
		"rewrite_config":        types.ObjectType{AttrTypes: albListenerRewriteConfigAttrTypes()},
		"target_group_config":   types.ObjectType{AttrTypes: albListenerTargetGroupConfigAttrTypes()},
	}
}

func albListenerRuleConditionAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"type":                types.StringType,
		"cookie_config":       types.ListType{ElemType: types.ObjectType{AttrTypes: albListenerCookieConfigAttrTypes()}},
		"header_config":       types.ObjectType{AttrTypes: albListenerHeaderConfigAttrTypes()},
		"host_config":         types.ListType{ElemType: types.StringType},
		"method_config":       types.ListType{ElemType: types.StringType},
		"path_config":         types.ListType{ElemType: types.StringType},
		"query_string_config": types.ListType{ElemType: types.ObjectType{AttrTypes: albListenerQueryStringConfigAttrTypes()}},
		"source_ip_config":    types.ListType{ElemType: types.StringType},
	}
}

func albListenerFixedResponseConfigAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{"http_code": types.Int64Type, "content": types.StringType, "content_type": types.StringType}
}

func albListenerInsertHeaderConfigAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{"key": types.StringType, "value": types.StringType, "value_type": types.StringType}
}

func albListenerRedirectConfigAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"http_code": types.Int64Type, "host": types.StringType, "path": types.StringType,
		"port": types.StringType, "protocol": types.StringType, "query": types.StringType,
	}
}

func albListenerRemoveHeaderConfigAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{"key": types.StringType}
}

func albListenerRewriteConfigAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{"host": types.StringType, "path": types.StringType, "query": types.StringType}
}

func albListenerTargetGroupConfigAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"target_groups":               types.ListType{ElemType: types.ObjectType{AttrTypes: albListenerTargetGroupTupleAttrTypes()}},
		"target_group_sticky_session": types.ObjectType{AttrTypes: albListenerTargetGroupStickySessionAttrTypes()},
	}
}

func albListenerTargetGroupTupleAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{"target_group_id": types.StringType, "weight": types.Int64Type}
}

func albListenerTargetGroupStickySessionAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{"enabled": types.BoolType, "timeout": types.Int64Type}
}

func albListenerCookieConfigAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{"key": types.StringType, "value": types.StringType}
}

func albListenerHeaderConfigAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{"key": types.StringType, "values": types.ListType{ElemType: types.StringType}}
}

func albListenerQueryStringConfigAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{"key": types.StringType, "value": types.StringType}
}

func albListenerTagAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{"tag_key": types.StringType, "tag_value": types.StringType}
}

// refreshAlbListenerModel copies cloud-owned fields into the model for the
// single rule tracked by this resource.
func refreshAlbListenerModel(ctx context.Context, m *AlbListenerResourceModel, item *albv20251030.RuleOutput) diag.Diagnostics {
	if item == nil {
		return nil
	}

	var diags diag.Diagnostics

	var rules []albRuleModel
	diags.Append(m.Rules.ElementsAs(ctx, &rules, false)...)
	if diags.HasError() {
		return diags
	}

	if len(rules) == 0 {
		rules = append(rules, albRuleModel{
			Actions:    types.ListNull(types.ObjectType{AttrTypes: albListenerRuleActionAttrTypes()}),
			Conditions: types.ListNull(types.ObjectType{AttrTypes: albListenerRuleConditionAttrTypes()}),
			Tags:       types.ListNull(types.ObjectType{AttrTypes: albListenerTagAttrTypes()}),
		})
	}

	idx := -1
	for i := range rules {
		if !rules[i].RuleId.IsNull() && rules[i].RuleId.ValueString() == fwHelper.StringOrEmpty(item.RuleId) {
			idx = i
			break
		}
	}
	if idx < 0 {
		idx = 0
	}

	rule := &rules[idx]
	rule.Actions = flattenAlbListenerActions(ctx, item.Actions, rule.Actions, &diags)
	rule.Conditions = flattenAlbListenerConditions(ctx, item.Conditions, rule.Conditions, &diags)
	rule.Priority = fwHelper.Int64ValueOrPrev(item.Priority, rule.Priority)
	rule.Direction = fwHelper.StringValueOrPrev(item.Direction, rule.Direction)
	rule.RuleName = fwHelper.StringValueOrPrev(item.RuleName, rule.RuleName)
	rule.Tags = flattenAlbListenerTags(ctx, item.Tags, rule.Tags, &diags)
	rule.RuleId = fwHelper.StringPointerValueOrNull(item.RuleId)
	rule.Status = fwHelper.StringPointerValueOrNull(item.Status)
	rule.CreateTime = fwHelper.StringPointerValueOrNull(item.CreateTime)
	rule.ModifyTime = fwHelper.StringPointerValueOrNull(item.ModifyTime)

	list, d := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: albListenerRuleAttrTypes()}, rules)
	diags.Append(d...)
	if !diags.HasError() {
		m.Rules = list
	}

	return diags
}

func findAlbListenerRule(rules []*albv20251030.RuleOutput, ruleId string) *albv20251030.RuleOutput {
	for _, rule := range rules {
		if rule != nil && rule.RuleId != nil && *rule.RuleId == ruleId {
			return rule
		}
	}
	return nil
}

// expandAlbListenerRules converts the planned rule list into the SDK RuleInput
// list used by CreateRules.
func expandAlbListenerRules(ctx context.Context, list types.List, out *[]*albv20251030.RuleInput) diag.Diagnostics {
	var diags diag.Diagnostics
	if list.IsNull() || list.IsUnknown() {
		return diags
	}

	var models []albRuleModel
	diags.Append(list.ElementsAs(ctx, &models, false)...)
	if diags.HasError() {
		return diags
	}

	result := make([]*albv20251030.RuleInput, 0, len(models))
	for i := range models {
		rule := &albv20251030.RuleInput{
			Priority:  fwHelper.Int64PointerFromValue(models[i].Priority),
			Direction: fwHelper.StringPointerFromValue(models[i].Direction),
			RuleName:  fwHelper.StringPointerFromValue(models[i].RuleName),
		}
		diags.Append(expandAlbListenerActions(ctx, models[i].Actions, &rule.Actions)...)
		diags.Append(expandAlbListenerConditions(ctx, models[i].Conditions, &rule.Conditions)...)
		diags.Append(expandAlbListenerTags(ctx, models[i].Tags, &rule.Tags)...)
		if diags.HasError() {
			return diags
		}
		result = append(result, rule)
	}

	*out = result
	return diags
}

// expandAlbListenerRulesForModify converts the planned rule list into the SDK
// RuleModify list used by ModifyRulesAttributes. The tracked rule id is always
// injected so the API addresses the correct rule.
func expandAlbListenerRulesForModify(ctx context.Context, list types.List, ruleId string, out *[]*albv20251030.RuleModify) diag.Diagnostics {
	var diags diag.Diagnostics
	if list.IsNull() || list.IsUnknown() {
		return diags
	}

	var models []albRuleModel
	diags.Append(list.ElementsAs(ctx, &models, false)...)
	if diags.HasError() {
		return diags
	}

	result := make([]*albv20251030.RuleModify, 0, len(models))
	for i := range models {
		rule := &albv20251030.RuleModify{
			RuleId:   fwHelper.String(ruleId),
			Priority: fwHelper.Int64PointerFromValue(models[i].Priority),
			RuleName: fwHelper.StringPointerFromValue(models[i].RuleName),
		}
		diags.Append(expandAlbListenerActions(ctx, models[i].Actions, &rule.Actions)...)
		diags.Append(expandAlbListenerConditions(ctx, models[i].Conditions, &rule.Conditions)...)
		if diags.HasError() {
			return diags
		}
		result = append(result, rule)
	}

	*out = result
	return diags
}

func expandAlbListenerActions(ctx context.Context, list types.List, out *[]*albv20251030.RuleAction) diag.Diagnostics {
	var diags diag.Diagnostics
	if list.IsNull() || list.IsUnknown() {
		return diags
	}

	var models []albRuleActionModel
	diags.Append(list.ElementsAs(ctx, &models, false)...)
	if diags.HasError() {
		return diags
	}

	result := make([]*albv20251030.RuleAction, 0, len(models))
	for i := range models {
		action := &albv20251030.RuleAction{
			Order: fwHelper.Int64PointerFromValue(models[i].Order),
			Type:  fwHelper.StringPointerFromValue(models[i].Type),
		}
		diags.Append(expandAlbListenerFixedResponseConfig(ctx, models[i].FixedResponseConfig, &action.FixedResponseConfig)...)
		diags.Append(expandAlbListenerInsertHeaderConfig(ctx, models[i].InsertHeaderConfig, &action.InsertHeaderConfig)...)
		diags.Append(expandAlbListenerRedirectConfig(ctx, models[i].RedirectConfig, &action.RedirectConfig)...)
		diags.Append(expandAlbListenerRemoveHeaderConfig(ctx, models[i].RemoveHeaderConfig, &action.RemoveHeaderConfig)...)
		diags.Append(expandAlbListenerRewriteConfig(ctx, models[i].RewriteConfig, &action.RewriteConfig)...)
		diags.Append(expandAlbListenerTargetGroupConfig(ctx, models[i].TargetGroupConfig, &action.TargetGroupConfig)...)
		if diags.HasError() {
			return diags
		}
		result = append(result, action)
	}

	*out = result
	return diags
}

func expandAlbListenerFixedResponseConfig(ctx context.Context, obj types.Object, out **albv20251030.FixedResponseInfo) diag.Diagnostics {
	var diags diag.Diagnostics
	if obj.IsNull() || obj.IsUnknown() {
		return diags
	}
	var m albFixedResponseConfigModel
	diags.Append(obj.As(ctx, &m, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return diags
	}
	*out = &albv20251030.FixedResponseInfo{
		HttpCode:    fwHelper.Int64PointerFromValue(m.HttpCode),
		Content:     fwHelper.StringPointerFromValue(m.Content),
		ContentType: fwHelper.StringPointerFromValue(m.ContentType),
	}
	return diags
}

func expandAlbListenerInsertHeaderConfig(ctx context.Context, obj types.Object, out **albv20251030.InsertHTTPHeaderInfo) diag.Diagnostics {
	var diags diag.Diagnostics
	if obj.IsNull() || obj.IsUnknown() {
		return diags
	}
	var m albInsertHeaderConfigModel
	diags.Append(obj.As(ctx, &m, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return diags
	}
	*out = &albv20251030.InsertHTTPHeaderInfo{
		Key:       fwHelper.StringPointerFromValue(m.Key),
		Value:     fwHelper.StringPointerFromValue(m.Value),
		ValueType: fwHelper.StringPointerFromValue(m.ValueType),
	}
	return diags
}

func expandAlbListenerRedirectConfig(ctx context.Context, obj types.Object, out **albv20251030.HTTPRedirectInfo) diag.Diagnostics {
	var diags diag.Diagnostics
	if obj.IsNull() || obj.IsUnknown() {
		return diags
	}
	var m albRedirectConfigModel
	diags.Append(obj.As(ctx, &m, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return diags
	}
	*out = &albv20251030.HTTPRedirectInfo{
		HttpCode: fwHelper.Int64PointerFromValue(m.HttpCode),
		Host:     fwHelper.StringPointerFromValue(m.Host),
		Path:     fwHelper.StringPointerFromValue(m.Path),
		Port:     fwHelper.StringPointerFromValue(m.Port),
		Protocol: fwHelper.StringPointerFromValue(m.Protocol),
		Query:    fwHelper.StringPointerFromValue(m.Query),
	}
	return diags
}

func expandAlbListenerRemoveHeaderConfig(ctx context.Context, obj types.Object, out **albv20251030.RemoveHTTPHeaderInfo) diag.Diagnostics {
	var diags diag.Diagnostics
	if obj.IsNull() || obj.IsUnknown() {
		return diags
	}
	var m albRemoveHeaderConfigModel
	diags.Append(obj.As(ctx, &m, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return diags
	}
	*out = &albv20251030.RemoveHTTPHeaderInfo{Key: fwHelper.StringPointerFromValue(m.Key)}
	return diags
}

func expandAlbListenerRewriteConfig(ctx context.Context, obj types.Object, out **albv20251030.HTTPRewriteInfo) diag.Diagnostics {
	var diags diag.Diagnostics
	if obj.IsNull() || obj.IsUnknown() {
		return diags
	}
	var m albRewriteConfigModel
	diags.Append(obj.As(ctx, &m, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return diags
	}
	*out = &albv20251030.HTTPRewriteInfo{
		Host:  fwHelper.StringPointerFromValue(m.Host),
		Path:  fwHelper.StringPointerFromValue(m.Path),
		Query: fwHelper.StringPointerFromValue(m.Query),
	}
	return diags
}

func expandAlbListenerTargetGroupConfig(ctx context.Context, obj types.Object, out **albv20251030.TargetGroupConfig) diag.Diagnostics {
	var diags diag.Diagnostics
	if obj.IsNull() || obj.IsUnknown() {
		return diags
	}
	var m albTargetGroupConfigModel
	diags.Append(obj.As(ctx, &m, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return diags
	}
	config := &albv20251030.TargetGroupConfig{}
	diags.Append(expandAlbListenerTargetGroupTuples(ctx, m.TargetGroups, &config.TargetGroups)...)
	diags.Append(expandAlbListenerStickySession(ctx, m.TargetGroupStickySession, &config.TargetGroupStickySession)...)
	if diags.HasError() {
		return diags
	}
	*out = config
	return diags
}

func expandAlbListenerTargetGroupTuples(ctx context.Context, list types.List, out *[]*albv20251030.TargetGroupTuple) diag.Diagnostics {
	var diags diag.Diagnostics
	if list.IsNull() || list.IsUnknown() {
		return diags
	}
	var models []albTargetGroupTupleModel
	diags.Append(list.ElementsAs(ctx, &models, false)...)
	if diags.HasError() {
		return diags
	}
	result := make([]*albv20251030.TargetGroupTuple, 0, len(models))
	for i := range models {
		result = append(result, &albv20251030.TargetGroupTuple{
			TargetGroupId: fwHelper.StringPointerFromValue(models[i].TargetGroupId),
			Weight:        fwHelper.Int64PointerFromValue(models[i].Weight),
		})
	}
	*out = result
	return diags
}

func expandAlbListenerStickySession(ctx context.Context, obj types.Object, out **albv20251030.TargetGroupStickySession) diag.Diagnostics {
	var diags diag.Diagnostics
	if obj.IsNull() || obj.IsUnknown() {
		return diags
	}
	var m albTargetGroupStickySessionModel
	diags.Append(obj.As(ctx, &m, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return diags
	}
	*out = &albv20251030.TargetGroupStickySession{
		Enabled: fwHelper.BoolPointerFromValue(m.Enabled),
		Timeout: fwHelper.Int64PointerFromValue(m.Timeout),
	}
	return diags
}

func expandAlbListenerConditions(ctx context.Context, list types.List, out *[]*albv20251030.RuleCondition) diag.Diagnostics {
	var diags diag.Diagnostics
	if list.IsNull() || list.IsUnknown() {
		return diags
	}
	var models []albRuleConditionModel
	diags.Append(list.ElementsAs(ctx, &models, false)...)
	if diags.HasError() {
		return diags
	}
	result := make([]*albv20251030.RuleCondition, 0, len(models))
	for i := range models {
		condition := &albv20251030.RuleCondition{
			Type:           fwHelper.StringPointerFromValue(models[i].Type),
			HostConfig:     fwHelper.PointersFromList[string](ctx, models[i].HostConfig, &diags),
			MethodConfig:   fwHelper.PointersFromList[string](ctx, models[i].MethodConfig, &diags),
			PathConfig:     fwHelper.PointersFromList[string](ctx, models[i].PathConfig, &diags),
			SourceIpConfig: fwHelper.PointersFromList[string](ctx, models[i].SourceIpConfig, &diags),
		}
		diags.Append(expandAlbListenerCookieConfigs(ctx, models[i].CookieConfig, &condition.CookieConfig)...)
		diags.Append(expandAlbListenerHeaderConfig(ctx, models[i].HeaderConfig, &condition.HeaderConfig)...)
		diags.Append(expandAlbListenerQueryStringConfigs(ctx, models[i].QueryStringConfig, &condition.QueryStringConfig)...)
		if diags.HasError() {
			return diags
		}
		result = append(result, condition)
	}
	*out = result
	return diags
}

func expandAlbListenerCookieConfigs(ctx context.Context, list types.List, out *[]*albv20251030.HTTPCookieInfo) diag.Diagnostics {
	var diags diag.Diagnostics
	if list.IsNull() || list.IsUnknown() {
		return diags
	}
	var models []albCookieConfigModel
	diags.Append(list.ElementsAs(ctx, &models, false)...)
	if diags.HasError() {
		return diags
	}
	result := make([]*albv20251030.HTTPCookieInfo, 0, len(models))
	for i := range models {
		result = append(result, &albv20251030.HTTPCookieInfo{
			Key:   fwHelper.StringPointerFromValue(models[i].Key),
			Value: fwHelper.StringPointerFromValue(models[i].Value),
		})
	}
	*out = result
	return diags
}

func expandAlbListenerHeaderConfig(ctx context.Context, obj types.Object, out **albv20251030.HTTPHeaderInfo) diag.Diagnostics {
	var diags diag.Diagnostics
	if obj.IsNull() || obj.IsUnknown() {
		return diags
	}
	var m albHeaderConfigModel
	diags.Append(obj.As(ctx, &m, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return diags
	}
	*out = &albv20251030.HTTPHeaderInfo{
		Key:    fwHelper.StringPointerFromValue(m.Key),
		Values: fwHelper.PointersFromList[string](ctx, m.Values, &diags),
	}
	return diags
}

func expandAlbListenerQueryStringConfigs(ctx context.Context, list types.List, out *[]*albv20251030.HTTPQueryStringInfo) diag.Diagnostics {
	var diags diag.Diagnostics
	if list.IsNull() || list.IsUnknown() {
		return diags
	}
	var models []albQueryStringConfigModel
	diags.Append(list.ElementsAs(ctx, &models, false)...)
	if diags.HasError() {
		return diags
	}
	result := make([]*albv20251030.HTTPQueryStringInfo, 0, len(models))
	for i := range models {
		result = append(result, &albv20251030.HTTPQueryStringInfo{
			Key:   fwHelper.StringPointerFromValue(models[i].Key),
			Value: fwHelper.StringPointerFromValue(models[i].Value),
		})
	}
	*out = result
	return diags
}

func expandAlbListenerTags(ctx context.Context, list types.List, out *[]*albv20251030.TagInfo) diag.Diagnostics {
	var diags diag.Diagnostics
	if list.IsNull() || list.IsUnknown() {
		return diags
	}
	var models []albTagModel
	diags.Append(list.ElementsAs(ctx, &models, false)...)
	if diags.HasError() {
		return diags
	}
	result := make([]*albv20251030.TagInfo, 0, len(models))
	for i := range models {
		result = append(result, &albv20251030.TagInfo{
			TagKey:   fwHelper.StringPointerFromValue(models[i].TagKey),
			TagValue: fwHelper.StringPointerFromValue(models[i].TagValue),
		})
	}
	*out = result
	return diags
}

// flatten helpers copy cloud values back into the model. User-configurable
// fields use the *OrPrev families; read-only computed fields use the *OrNull
// families.
func flattenAlbListenerActions(ctx context.Context, cloud []*albv20251030.RuleAction, prev types.List, diags *diag.Diagnostics) types.List {
	if cloud == nil && (prev.IsNull() || prev.IsUnknown()) {
		return types.ListNull(types.ObjectType{AttrTypes: albListenerRuleActionAttrTypes()})
	}

	var prevModels []albRuleActionModel
	if !prev.IsNull() && !prev.IsUnknown() {
		diags.Append(prev.ElementsAs(ctx, &prevModels, false)...)
		if diags.HasError() {
			return prev
		}
	}

	result := make([]albRuleActionModel, 0, len(cloud))
	for i, action := range cloud {
		m := albRuleActionModel{
			Order: fwHelper.Int64ValueOrNull(action.Order),
			Type:  fwHelper.StringPointerValueOrNull(action.Type),
		}
		if i < len(prevModels) {
			m.Order = fwHelper.Int64ValueOrPrev(action.Order, prevModels[i].Order)
			m.Type = fwHelper.StringValueOrPrev(action.Type, prevModels[i].Type)
		}
		m.FixedResponseConfig = flattenAlbListenerFixedResponseConfig(ctx, action.FixedResponseConfig, prevObjectAt(ctx, prevModels, i, "fixed_response_config"), diags)
		m.InsertHeaderConfig = flattenAlbListenerInsertHeaderConfig(ctx, action.InsertHeaderConfig, prevObjectAt(ctx, prevModels, i, "insert_header_config"), diags)
		m.RedirectConfig = flattenAlbListenerRedirectConfig(ctx, action.RedirectConfig, prevObjectAt(ctx, prevModels, i, "redirect_config"), diags)
		m.RemoveHeaderConfig = flattenAlbListenerRemoveHeaderConfig(ctx, action.RemoveHeaderConfig, prevObjectAt(ctx, prevModels, i, "remove_header_config"), diags)
		m.RewriteConfig = flattenAlbListenerRewriteConfig(ctx, action.RewriteConfig, prevObjectAt(ctx, prevModels, i, "rewrite_config"), diags)
		m.TargetGroupConfig = flattenAlbListenerTargetGroupConfig(ctx, action.TargetGroupConfig, prevObjectAt(ctx, prevModels, i, "target_group_config"), diags)
		result = append(result, m)
	}

	if len(result) == 0 && (prev.IsNull() || prev.IsUnknown()) {
		return types.ListNull(types.ObjectType{AttrTypes: albListenerRuleActionAttrTypes()})
	}

	list, d := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: albListenerRuleActionAttrTypes()}, result)
	diags.Append(d...)
	return list
}

func prevObjectAt(ctx context.Context, models []albRuleActionModel, idx int, field string) types.Object {
	nullObj := func(attrs map[string]attr.Type) types.Object { return types.ObjectNull(attrs) }
	if idx >= len(models) {
		switch field {
		case "fixed_response_config":
			return nullObj(albListenerFixedResponseConfigAttrTypes())
		case "insert_header_config":
			return nullObj(albListenerInsertHeaderConfigAttrTypes())
		case "redirect_config":
			return nullObj(albListenerRedirectConfigAttrTypes())
		case "remove_header_config":
			return nullObj(albListenerRemoveHeaderConfigAttrTypes())
		case "rewrite_config":
			return nullObj(albListenerRewriteConfigAttrTypes())
		case "target_group_config":
			return nullObj(albListenerTargetGroupConfigAttrTypes())
		}
		return types.ObjectNull(map[string]attr.Type{})
	}

	switch field {
	case "fixed_response_config":
		return models[idx].FixedResponseConfig
	case "insert_header_config":
		return models[idx].InsertHeaderConfig
	case "redirect_config":
		return models[idx].RedirectConfig
	case "remove_header_config":
		return models[idx].RemoveHeaderConfig
	case "rewrite_config":
		return models[idx].RewriteConfig
	case "target_group_config":
		return models[idx].TargetGroupConfig
	}
	return types.ObjectNull(map[string]attr.Type{})
}

func flattenAlbListenerFixedResponseConfig(ctx context.Context, cloud *albv20251030.FixedResponseInfo, prev types.Object, diags *diag.Diagnostics) types.Object {
	if cloud == nil {
		return types.ObjectNull(albListenerFixedResponseConfigAttrTypes())
	}
	var m albFixedResponseConfigModel
	if !prev.IsNull() && !prev.IsUnknown() {
		diags.Append(prev.As(ctx, &m, basetypes.ObjectAsOptions{})...)
		if diags.HasError() {
			return prev
		}
	}
	m.HttpCode = fwHelper.Int64ValueOrPrev(cloud.HttpCode, m.HttpCode)
	m.Content = fwHelper.StringValueOrPrev(cloud.Content, m.Content)
	m.ContentType = fwHelper.StringValueOrPrev(cloud.ContentType, m.ContentType)
	obj, d := types.ObjectValueFrom(ctx, albListenerFixedResponseConfigAttrTypes(), m)
	diags.Append(d...)
	return obj
}

func flattenAlbListenerInsertHeaderConfig(ctx context.Context, cloud *albv20251030.InsertHTTPHeaderInfo, prev types.Object, diags *diag.Diagnostics) types.Object {
	if cloud == nil {
		return types.ObjectNull(albListenerInsertHeaderConfigAttrTypes())
	}
	var m albInsertHeaderConfigModel
	if !prev.IsNull() && !prev.IsUnknown() {
		diags.Append(prev.As(ctx, &m, basetypes.ObjectAsOptions{})...)
		if diags.HasError() {
			return prev
		}
	}
	m.Key = fwHelper.StringValueOrPrev(cloud.Key, m.Key)
	m.Value = fwHelper.StringValueOrPrev(cloud.Value, m.Value)
	m.ValueType = fwHelper.StringValueOrPrev(cloud.ValueType, m.ValueType)
	obj, d := types.ObjectValueFrom(ctx, albListenerInsertHeaderConfigAttrTypes(), m)
	diags.Append(d...)
	return obj
}

func flattenAlbListenerRedirectConfig(ctx context.Context, cloud *albv20251030.HTTPRedirectInfo, prev types.Object, diags *diag.Diagnostics) types.Object {
	if cloud == nil {
		return types.ObjectNull(albListenerRedirectConfigAttrTypes())
	}
	var m albRedirectConfigModel
	if !prev.IsNull() && !prev.IsUnknown() {
		diags.Append(prev.As(ctx, &m, basetypes.ObjectAsOptions{})...)
		if diags.HasError() {
			return prev
		}
	}
	m.HttpCode = fwHelper.Int64ValueOrPrev(cloud.HttpCode, m.HttpCode)
	m.Host = fwHelper.StringValueOrPrev(cloud.Host, m.Host)
	m.Path = fwHelper.StringValueOrPrev(cloud.Path, m.Path)
	m.Port = fwHelper.StringValueOrPrev(cloud.Port, m.Port)
	m.Protocol = fwHelper.StringValueOrPrev(cloud.Protocol, m.Protocol)
	m.Query = fwHelper.StringValueOrPrev(cloud.Query, m.Query)
	obj, d := types.ObjectValueFrom(ctx, albListenerRedirectConfigAttrTypes(), m)
	diags.Append(d...)
	return obj
}

func flattenAlbListenerRemoveHeaderConfig(ctx context.Context, cloud *albv20251030.RemoveHTTPHeaderInfo, prev types.Object, diags *diag.Diagnostics) types.Object {
	if cloud == nil {
		return types.ObjectNull(albListenerRemoveHeaderConfigAttrTypes())
	}
	var m albRemoveHeaderConfigModel
	if !prev.IsNull() && !prev.IsUnknown() {
		diags.Append(prev.As(ctx, &m, basetypes.ObjectAsOptions{})...)
		if diags.HasError() {
			return prev
		}
	}
	m.Key = fwHelper.StringValueOrPrev(cloud.Key, m.Key)
	obj, d := types.ObjectValueFrom(ctx, albListenerRemoveHeaderConfigAttrTypes(), m)
	diags.Append(d...)
	return obj
}

func flattenAlbListenerRewriteConfig(ctx context.Context, cloud *albv20251030.HTTPRewriteInfo, prev types.Object, diags *diag.Diagnostics) types.Object {
	if cloud == nil {
		return types.ObjectNull(albListenerRewriteConfigAttrTypes())
	}
	var m albRewriteConfigModel
	if !prev.IsNull() && !prev.IsUnknown() {
		diags.Append(prev.As(ctx, &m, basetypes.ObjectAsOptions{})...)
		if diags.HasError() {
			return prev
		}
	}
	m.Host = fwHelper.StringValueOrPrev(cloud.Host, m.Host)
	m.Path = fwHelper.StringValueOrPrev(cloud.Path, m.Path)
	m.Query = fwHelper.StringValueOrPrev(cloud.Query, m.Query)
	obj, d := types.ObjectValueFrom(ctx, albListenerRewriteConfigAttrTypes(), m)
	diags.Append(d...)
	return obj
}

func flattenAlbListenerTargetGroupConfig(ctx context.Context, cloud *albv20251030.TargetGroupConfig, prev types.Object, diags *diag.Diagnostics) types.Object {
	if cloud == nil {
		return types.ObjectNull(albListenerTargetGroupConfigAttrTypes())
	}
	var m albTargetGroupConfigModel
	var prevTuples types.List
	var prevSticky types.Object
	if !prev.IsNull() && !prev.IsUnknown() {
		diags.Append(prev.As(ctx, &m, basetypes.ObjectAsOptions{})...)
		if diags.HasError() {
			return prev
		}
		prevTuples = m.TargetGroups
		prevSticky = m.TargetGroupStickySession
	}
	prevTuples = nullListIfUnset(ctx, prevTuples, types.ObjectType{AttrTypes: albListenerTargetGroupTupleAttrTypes()})
	if prevSticky.IsNull() || prevSticky.IsUnknown() {
		prevSticky = types.ObjectNull(albListenerTargetGroupStickySessionAttrTypes())
	}
	m.TargetGroups = flattenAlbListenerTargetGroupTuples(ctx, cloud.TargetGroups, prevTuples, diags)
	m.TargetGroupStickySession = flattenAlbListenerStickySession(ctx, cloud.TargetGroupStickySession, prevSticky, diags)
	obj, d := types.ObjectValueFrom(ctx, albListenerTargetGroupConfigAttrTypes(), m)
	diags.Append(d...)
	return obj
}

func flattenAlbListenerTargetGroupTuples(ctx context.Context, cloud []*albv20251030.TargetGroupTuple, prev types.List, diags *diag.Diagnostics) types.List {
	if cloud == nil {
		return prev
	}
	var prevModels []albTargetGroupTupleModel
	if !prev.IsNull() && !prev.IsUnknown() {
		diags.Append(prev.ElementsAs(ctx, &prevModels, false)...)
		if diags.HasError() {
			return prev
		}
	}
	result := make([]albTargetGroupTupleModel, 0, len(cloud))
	for i, tuple := range cloud {
		m := albTargetGroupTupleModel{
			TargetGroupId: fwHelper.StringPointerValueOrNull(tuple.TargetGroupId),
			Weight:        fwHelper.Int64ValueOrNull(tuple.Weight),
		}
		if i < len(prevModels) {
			m.TargetGroupId = fwHelper.StringValueOrPrev(tuple.TargetGroupId, prevModels[i].TargetGroupId)
			m.Weight = fwHelper.Int64ValueOrPrev(tuple.Weight, prevModels[i].Weight)
		}
		result = append(result, m)
	}
	list, d := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: albListenerTargetGroupTupleAttrTypes()}, result)
	diags.Append(d...)
	return list
}

func flattenAlbListenerStickySession(ctx context.Context, cloud *albv20251030.TargetGroupStickySession, prev types.Object, diags *diag.Diagnostics) types.Object {
	if cloud == nil {
		return prev
	}
	var m albTargetGroupStickySessionModel
	if !prev.IsNull() && !prev.IsUnknown() {
		diags.Append(prev.As(ctx, &m, basetypes.ObjectAsOptions{})...)
		if diags.HasError() {
			return prev
		}
	}
	m.Enabled = fwHelper.BoolValueOrPrev(cloud.Enabled, m.Enabled)
	m.Timeout = fwHelper.Int64ValueOrPrev(cloud.Timeout, m.Timeout)
	obj, d := types.ObjectValueFrom(ctx, albListenerTargetGroupStickySessionAttrTypes(), m)
	diags.Append(d...)
	return obj
}

func flattenAlbListenerConditions(ctx context.Context, cloud []*albv20251030.RuleCondition, prev types.List, diags *diag.Diagnostics) types.List {
	if cloud == nil && (prev.IsNull() || prev.IsUnknown()) {
		return types.ListNull(types.ObjectType{AttrTypes: albListenerRuleConditionAttrTypes()})
	}
	var prevModels []albRuleConditionModel
	if !prev.IsNull() && !prev.IsUnknown() {
		diags.Append(prev.ElementsAs(ctx, &prevModels, false)...)
		if diags.HasError() {
			return prev
		}
	}
	result := make([]albRuleConditionModel, 0, len(cloud))
	for i, condition := range cloud {
		m := albRuleConditionModel{
			Type:              fwHelper.StringPointerValueOrNull(condition.Type),
			HostConfig:        fwHelper.ListFromPointers[string](ctx, types.StringType, condition.HostConfig, diags),
			MethodConfig:      fwHelper.ListFromPointers[string](ctx, types.StringType, condition.MethodConfig, diags),
			PathConfig:        fwHelper.ListFromPointers[string](ctx, types.StringType, condition.PathConfig, diags),
			SourceIpConfig:    fwHelper.ListFromPointers[string](ctx, types.StringType, condition.SourceIpConfig, diags),
			CookieConfig:      flattenAlbListenerCookieConfigs(ctx, condition.CookieConfig, listObjectAt(ctx, prevModels, i, func(m albRuleConditionModel) types.List { return m.CookieConfig }, albListenerCookieConfigAttrTypes()), diags),
			QueryStringConfig: flattenAlbListenerQueryStringConfigs(ctx, condition.QueryStringConfig, listObjectAt(ctx, prevModels, i, func(m albRuleConditionModel) types.List { return m.QueryStringConfig }, albListenerQueryStringConfigAttrTypes()), diags),
			HeaderConfig:      flattenAlbListenerHeaderConfig(ctx, condition.HeaderConfig, objectAtCondition(ctx, prevModels, i), diags),
		}
		if i < len(prevModels) {
			m.Type = fwHelper.StringValueOrPrev(condition.Type, prevModels[i].Type)
			m.HostConfig = fwHelper.ListFromPointersOrPrev[string](ctx, types.StringType, condition.HostConfig, prevModels[i].HostConfig, diags)
			m.MethodConfig = fwHelper.ListFromPointersOrPrev[string](ctx, types.StringType, condition.MethodConfig, prevModels[i].MethodConfig, diags)
			m.PathConfig = fwHelper.ListFromPointersOrPrev[string](ctx, types.StringType, condition.PathConfig, prevModels[i].PathConfig, diags)
			m.SourceIpConfig = fwHelper.ListFromPointersOrPrev[string](ctx, types.StringType, condition.SourceIpConfig, prevModels[i].SourceIpConfig, diags)
		}
		result = append(result, m)
	}
	if len(result) == 0 && (prev.IsNull() || prev.IsUnknown()) {
		return types.ListNull(types.ObjectType{AttrTypes: albListenerRuleConditionAttrTypes()})
	}
	list, d := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: albListenerRuleConditionAttrTypes()}, result)
	diags.Append(d...)
	return list
}

func listObjectAt(ctx context.Context, models []albRuleConditionModel, idx int, getter func(albRuleConditionModel) types.List, attrs map[string]attr.Type) types.List {
	if idx >= len(models) {
		return types.ListNull(types.ObjectType{AttrTypes: attrs})
	}
	list := getter(models[idx])
	if list.IsNull() || list.IsUnknown() {
		return types.ListNull(types.ObjectType{AttrTypes: attrs})
	}
	return list
}

func objectAtCondition(ctx context.Context, models []albRuleConditionModel, idx int) types.Object {
	if idx >= len(models) {
		return types.ObjectNull(albListenerHeaderConfigAttrTypes())
	}
	obj := models[idx].HeaderConfig
	if obj.IsNull() || obj.IsUnknown() {
		return types.ObjectNull(albListenerHeaderConfigAttrTypes())
	}
	return obj
}

func flattenAlbListenerCookieConfigs(ctx context.Context, cloud []*albv20251030.HTTPCookieInfo, prev types.List, diags *diag.Diagnostics) types.List {
	if cloud == nil {
		if prev.IsNull() || prev.IsUnknown() {
			return types.ListNull(types.ObjectType{AttrTypes: albListenerCookieConfigAttrTypes()})
		}
		return prev
	}
	var prevModels []albCookieConfigModel
	if !prev.IsNull() && !prev.IsUnknown() {
		diags.Append(prev.ElementsAs(ctx, &prevModels, false)...)
		if diags.HasError() {
			return prev
		}
	}
	result := make([]albCookieConfigModel, 0, len(cloud))
	for i, cookie := range cloud {
		m := albCookieConfigModel{Key: fwHelper.StringPointerValueOrNull(cookie.Key), Value: fwHelper.StringPointerValueOrNull(cookie.Value)}
		if i < len(prevModels) {
			m.Key = fwHelper.StringValueOrPrev(cookie.Key, prevModels[i].Key)
			m.Value = fwHelper.StringValueOrPrev(cookie.Value, prevModels[i].Value)
		}
		result = append(result, m)
	}
	list, d := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: albListenerCookieConfigAttrTypes()}, result)
	diags.Append(d...)
	return list
}

func flattenAlbListenerQueryStringConfigs(ctx context.Context, cloud []*albv20251030.HTTPQueryStringInfo, prev types.List, diags *diag.Diagnostics) types.List {
	if cloud == nil {
		if prev.IsNull() || prev.IsUnknown() {
			return types.ListNull(types.ObjectType{AttrTypes: albListenerQueryStringConfigAttrTypes()})
		}
		return prev
	}
	var prevModels []albQueryStringConfigModel
	if !prev.IsNull() && !prev.IsUnknown() {
		diags.Append(prev.ElementsAs(ctx, &prevModels, false)...)
		if diags.HasError() {
			return prev
		}
	}
	result := make([]albQueryStringConfigModel, 0, len(cloud))
	for i, q := range cloud {
		m := albQueryStringConfigModel{Key: fwHelper.StringPointerValueOrNull(q.Key), Value: fwHelper.StringPointerValueOrNull(q.Value)}
		if i < len(prevModels) {
			m.Key = fwHelper.StringValueOrPrev(q.Key, prevModels[i].Key)
			m.Value = fwHelper.StringValueOrPrev(q.Value, prevModels[i].Value)
		}
		result = append(result, m)
	}
	list, d := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: albListenerQueryStringConfigAttrTypes()}, result)
	diags.Append(d...)
	return list
}

func flattenAlbListenerHeaderConfig(ctx context.Context, cloud *albv20251030.HTTPHeaderInfo, prev types.Object, diags *diag.Diagnostics) types.Object {
	if cloud == nil {
		return prev
	}
	var m albHeaderConfigModel
	if !prev.IsNull() && !prev.IsUnknown() {
		diags.Append(prev.As(ctx, &m, basetypes.ObjectAsOptions{})...)
		if diags.HasError() {
			return prev
		}
	}
	m.Key = fwHelper.StringValueOrPrev(cloud.Key, m.Key)
	values := nullListIfUnset(ctx, m.Values, types.StringType)
	m.Values = fwHelper.ListFromPointersOrPrev[string](ctx, types.StringType, cloud.Values, values, diags)
	obj, d := types.ObjectValueFrom(ctx, albListenerHeaderConfigAttrTypes(), m)
	diags.Append(d...)
	return obj
}

func flattenAlbListenerTags(ctx context.Context, cloud []*albv20251030.TagInfo, prev types.List, diags *diag.Diagnostics) types.List {
	if cloud == nil {
		if prev.IsNull() || prev.IsUnknown() {
			return types.ListNull(types.ObjectType{AttrTypes: albListenerTagAttrTypes()})
		}
		return prev
	}
	var prevModels []albTagModel
	if !prev.IsNull() && !prev.IsUnknown() {
		diags.Append(prev.ElementsAs(ctx, &prevModels, false)...)
		if diags.HasError() {
			return prev
		}
	}
	result := make([]albTagModel, 0, len(cloud))
	for i, tag := range cloud {
		m := albTagModel{TagKey: fwHelper.StringPointerValueOrNull(tag.TagKey), TagValue: fwHelper.StringPointerValueOrNull(tag.TagValue)}
		if i < len(prevModels) {
			m.TagKey = fwHelper.StringValueOrPrev(tag.TagKey, prevModels[i].TagKey)
			m.TagValue = fwHelper.StringValueOrPrev(tag.TagValue, prevModels[i].TagValue)
		}
		result = append(result, m)
	}
	list, d := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: albListenerTagAttrTypes()}, result)
	diags.Append(d...)
	return list
}

func nullListIfUnset(ctx context.Context, list types.List, elemType attr.Type) types.List {
	if list.IsNull() || list.IsUnknown() {
		return types.ListNull(elemType)
	}
	return list
}
