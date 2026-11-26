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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	alb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/alb/v20251030"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/fw"
	fwhelper "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

// Compile-time interface assertions. The base resource.Resource interface is
// checked implicitly by the factory return type.
var (
	_ resource.ResourceWithConfigure   = &AlbRuleResource{}
	_ resource.ResourceWithImportState = &AlbRuleResource{}

	// errAlbRuleNotFound marks a confirmed "rule does not exist" result of
	// DescribeRules. It is distinct from a transport or API failure.
	errAlbRuleNotFound = errors.New("alb rule not found")

	// errAlbRuleStatusPending marks a rule that DescribeRules reports as still
	// being created or configured. The documented statuses are Provisioning,
	// Active and Configuring.
	errAlbRuleStatusPending = errors.New("alb rule status not converged")
)

// albRuleAdoptKey marks the first Read after an import. With no prior
// configuration to preserve, cloud-echoed optional values are adopted into
// state; on a regular refresh an omitted optional attribute must stay null.
type albRuleAdoptKey struct{}

// NewAlbRuleResource is the factory referenced by
// tencentcloud/framework/registry.go to register this resource.
func NewAlbRuleResource() resource.Resource {
	return &AlbRuleResource{}
}

// AlbRuleResource implements resource.Resource for tencentcloud_alb_rule.
//
// CreateRules accepts a list of rules and returns a list of ids. This resource
// manages exactly one forwarding rule: both lists are required to have exactly
// one element. The three identifiers every read/update/delete call needs
// (LoadBalancerId, ListenerId, RuleId) form the Terraform id.
type AlbRuleResource struct {
	fw.ResourceWithConfigure
}

// AlbRuleResourceModel maps the top-level schema attributes. Every schema
// attribute appears here exactly once.
type AlbRuleResourceModel struct {
	Id             types.String `tfsdk:"id"`
	ListenerId     types.String `tfsdk:"listener_id"`
	LoadBalancerId types.String `tfsdk:"load_balancer_id"`
	Rules          types.List   `tfsdk:"rules"`
	ClientToken    types.String `tfsdk:"client_token"`
	DryRun         types.Bool   `tfsdk:"dry_run"`
	NextToken      types.String `tfsdk:"next_token"`
}

// albRuleInputModel maps one element of rules.
type albRuleInputModel struct {
	Actions    types.List   `tfsdk:"actions"`
	Conditions types.List   `tfsdk:"conditions"`
	Priority   types.Int64  `tfsdk:"priority"`
	Direction  types.String `tfsdk:"direction"`
	RuleName   types.String `tfsdk:"rule_name"`
	Tags       types.List   `tfsdk:"tags"`
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

type albRuleFixedResponseModel struct {
	HttpCode    types.Int64  `tfsdk:"http_code"`
	Content     types.String `tfsdk:"content"`
	ContentType types.String `tfsdk:"content_type"`
}

type albRuleInsertHeaderModel struct {
	Key       types.String `tfsdk:"key"`
	Value     types.String `tfsdk:"value"`
	ValueType types.String `tfsdk:"value_type"`
}

type albRuleRedirectModel struct {
	HttpCode types.Int64  `tfsdk:"http_code"`
	Host     types.String `tfsdk:"host"`
	Path     types.String `tfsdk:"path"`
	Port     types.String `tfsdk:"port"`
	Protocol types.String `tfsdk:"protocol"`
	Query    types.String `tfsdk:"query"`
}

type albRuleRemoveHeaderModel struct {
	Key types.String `tfsdk:"key"`
}

type albRuleRewriteModel struct {
	Host  types.String `tfsdk:"host"`
	Path  types.String `tfsdk:"path"`
	Query types.String `tfsdk:"query"`
}

type albRuleTargetGroupConfigModel struct {
	TargetGroups            types.List   `tfsdk:"target_groups"`
	TargetGroupStickySession types.Object `tfsdk:"target_group_sticky_session"`
}

type albRuleTargetGroupModel struct {
	TargetGroupId types.String `tfsdk:"target_group_id"`
	Weight        types.Int64  `tfsdk:"weight"`
}

type albRuleStickySessionModel struct {
	Enabled types.Bool  `tfsdk:"enabled"`
	Timeout types.Int64 `tfsdk:"timeout"`
}

type albRuleConditionModel struct {
	Type              types.String `tfsdk:"type"`
	CookieConfig      types.List   `tfsdk:"cookie_config"`
	HeaderConfig      types.Object `tfsdk:"header_config"`
	HostConfig        types.Set    `tfsdk:"host_config"`
	MethodConfig      types.Set    `tfsdk:"method_config"`
	PathConfig        types.Set    `tfsdk:"path_config"`
	QueryStringConfig types.List   `tfsdk:"query_string_config"`
	SourceIpConfig    types.Set    `tfsdk:"source_ip_config"`
}

type albRuleCookieModel struct {
	Key   types.String `tfsdk:"key"`
	Value types.String `tfsdk:"value"`
}

type albRuleHeaderModel struct {
	Key    types.String `tfsdk:"key"`
	Values types.Set    `tfsdk:"values"`
}

type albRuleQueryStringModel struct {
	Key   types.String `tfsdk:"key"`
	Value types.String `tfsdk:"value"`
}

type albRuleTagModel struct {
	TagKey   types.String `tfsdk:"tag_key"`
	TagValue types.String `tfsdk:"tag_value"`
}

func (r *AlbRuleResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "tencentcloud_alb_rule"
}

func (r *AlbRuleResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a resource to create and manage an ALB forwarding rule.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Resource identifier in the form `<load_balancer_id>#<listener_id>#<rule_id>`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"listener_id": schema.StringAttribute{
				Required:    true,
				Description: "Listener ID in the format `lst-` followed by 8 alphanumeric characters. Immutable: changing it forces resource replacement.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"load_balancer_id": schema.StringAttribute{
				Required:    true,
				Description: "Load balancer instance ID in the format `alb-` followed by 8 alphanumeric characters. Immutable: changing it forces resource replacement.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"rules": schema.ListNestedAttribute{
				Required:    true,
				Description: "Forwarding rule list. This resource manages exactly one rule, so the list must contain exactly one element.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"actions": schema.ListNestedAttribute{
							Required:    true,
							Description: "Forwarding action list.",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"order": schema.Int64Attribute{
										Required:    true,
										Description: "Action execution order. Values must not repeat and are executed ascending. Valid range: 1 to 50000.",
									},
									"type": schema.StringAttribute{
										Required:    true,
										Description: "Action type. Valid values: `TargetGroup`, `Redirect`, `FixedResponse`, `Rewrite`, `InsertHeader`, `RemoveHeader`. One of `TargetGroup`, `Redirect` or `FixedResponse` is required and must be placed last.",
									},
									"fixed_response_config": schema.SingleNestedAttribute{
										Optional:    true,
										Description: "Fixed response configuration.",
										Attributes: map[string]schema.Attribute{
											"http_code": schema.Int64Attribute{
												Required:    true,
												Description: "Returned HTTP response code. Supports 2xx, 4xx and 5xx.",
											},
											"content": schema.StringAttribute{
												Optional:    true,
												Description: "Returned fixed content. ASCII characters only, at most 1KB.",
											},
											"content_type": schema.StringAttribute{
												Optional:    true,
												Description: "Format of the returned content. Valid values: `text/plain`, `text/css`, `text/html`, `application/javascript`, `application/json`.",
											},
										},
									},
									"insert_header_config": schema.SingleNestedAttribute{
										Optional:    true,
										Description: "Insert HTTP header configuration.",
										Attributes: map[string]schema.Attribute{
											"key": schema.StringAttribute{
												Required:    true,
												Description: "HTTP header key, 1 to 40 characters. Supports a-z, A-Z, 0-9, `-` and `_`.",
											},
											"value": schema.StringAttribute{
												Required:    true,
												Description: "HTTP header value. Its allowed values depend on `value_type`.",
											},
											"value_type": schema.StringAttribute{
												Required:    true,
												Description: "Type of the HTTP header value. Valid values: `SystemDefined`, `UserDefined`, `ReferenceHeader`.",
											},
										},
									},
									"redirect_config": schema.SingleNestedAttribute{
										Optional:    true,
										Description: "Redirect configuration. Except for `http_code`, not all other settings may use their default values.",
										Attributes: map[string]schema.Attribute{
											"http_code": schema.Int64Attribute{
												Required:    true,
												Description: "Redirect HTTP code. Valid values: 301, 302, 303, 307, 308.",
											},
											"host": schema.StringAttribute{
												Optional:    true,
												Description: "Redirect host. Defaults to `${host}`.",
											},
											"path": schema.StringAttribute{
												Optional:    true,
												Description: "Redirect path. Defaults to `${path}`.",
											},
											"port": schema.StringAttribute{
												Optional:    true,
												Description: "Redirect port. Defaults to `${port}`. Valid range: 1 to 65535.",
											},
											"protocol": schema.StringAttribute{
												Optional:    true,
												Description: "Redirect protocol. Valid values: `HTTP`, `HTTPS`. Defaults to `${protocol}`.",
											},
											"query": schema.StringAttribute{
												Optional:    true,
												Description: "Redirect query string. Defaults to `${query}`.",
											},
										},
									},
									"remove_header_config": schema.SingleNestedAttribute{
										Optional:    true,
										Description: "Remove HTTP header configuration.",
										Attributes: map[string]schema.Attribute{
											"key": schema.StringAttribute{
												Required:    true,
												Description: "HTTP header key to remove, 1 to 40 characters.",
											},
										},
									},
									"rewrite_config": schema.SingleNestedAttribute{
										Optional:    true,
										Description: "Rewrite configuration.",
										Attributes: map[string]schema.Attribute{
											"host": schema.StringAttribute{
												Optional:    true,
												Description: "Rewrite host. Defaults to `${host}`.",
											},
											"path": schema.StringAttribute{
												Optional:    true,
												Description: "Rewrite path. Defaults to `${path}`.",
											},
											"query": schema.StringAttribute{
												Optional:    true,
												Description: "Rewrite query string. Defaults to `${query}`.",
											},
										},
									},
									"target_group_config": schema.SingleNestedAttribute{
										Optional:    true,
										Description: "Forwarding target group configuration.",
										Attributes: map[string]schema.Attribute{
											"target_groups": schema.ListNestedAttribute{
												Required:    true,
												Description: "Target group list.",
												NestedObject: schema.NestedAttributeObject{
													Attributes: map[string]schema.Attribute{
														"target_group_id": schema.StringAttribute{
															Required:    true,
															Description: "Target group ID in the format `lbtg-` followed by 8 alphanumeric characters.",
														},
														"weight": schema.Int64Attribute{
															Optional:    true,
															Description: "Weight. Valid range: 0 to 100. Defaults to 10.",
														},
													},
												},
											},
											"target_group_sticky_session": schema.SingleNestedAttribute{
												Optional:    true,
												Description: "Session affinity across target groups.",
												Attributes: map[string]schema.Attribute{
													"enabled": schema.BoolAttribute{
														Required:    true,
														Description: "Whether session affinity is enabled. Disabled by default.",
													},
													"timeout": schema.Int64Attribute{
														Optional:    true,
														Description: "Session timeout in seconds. Valid range: 1 to 86400. Defaults to 1000.",
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
							Description: "Forwarding condition list.",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"type": schema.StringAttribute{
										Required:    true,
										Description: "Condition type. Valid values: `Host`, `Path`, `Header`, `QueryString`, `Method`, `Cookie`, `SourceIp`.",
									},
									"cookie_config": schema.ListNestedAttribute{
										Optional:    true,
										Description: "Cookie configuration.",
										NestedObject: schema.NestedAttributeObject{
											Attributes: map[string]schema.Attribute{
												"key": schema.StringAttribute{
													Required:    true,
													Description: "Cookie key, 1 to 64 characters.",
												},
												"value": schema.StringAttribute{
													Required:    true,
													Description: "Cookie value, 1 to 128 characters.",
												},
											},
										},
									},
									"header_config": schema.SingleNestedAttribute{
										Optional:    true,
										Description: "HTTP header condition configuration.",
										Attributes: map[string]schema.Attribute{
											"key": schema.StringAttribute{
												Required:    true,
												Description: "HTTP header key, 1 to 40 characters.",
											},
											"values": schema.SetAttribute{
												Required:    true,
												ElementType: types.StringType,
												Description: "HTTP header values, 1 to 128 characters each.",
											},
										},
									},
									"host_config": schema.SetAttribute{
										Optional:    true,
										ElementType: types.StringType,
										Description: "Host names. Supports exact, regular expression and wildcard matching.",
									},
									"method_config": schema.SetAttribute{
										Optional:    true,
										ElementType: types.StringType,
										Description: "Request methods. Valid values: `HEAD`, `GET`, `POST`, `OPTIONS`, `PUT`, `PATCH`, `DELETE`.",
									},
									"path_config": schema.SetAttribute{
										Optional:    true,
										ElementType: types.StringType,
										Description: "Forwarding paths. Supports exact, regular expression and wildcard matching.",
									},
									"query_string_config": schema.ListNestedAttribute{
										Optional:    true,
										Description: "Query string configuration.",
										NestedObject: schema.NestedAttributeObject{
											Attributes: map[string]schema.Attribute{
												"key": schema.StringAttribute{
													Required:    true,
													Description: "Query string key, 1 to 16 characters.",
												},
												"value": schema.StringAttribute{
													Required:    true,
													Description: "Query string value, 1 to 128 characters.",
												},
											},
										},
									},
									"source_ip_config": schema.SetAttribute{
										Optional:    true,
										ElementType: types.StringType,
										Description: "Source IP CIDRs, for example `x.x.x.x/32` or `x.x.x.x/24`.",
									},
								},
							},
						},
						"priority": schema.Int64Attribute{
							Required:    true,
							Description: "Rule priority. Lower values have higher priority. Valid range: 1 to 10000.",
						},
						"direction": schema.StringAttribute{
							Optional:    true,
							Description: "Rule direction. Valid values: `Request` (client to load balancer), `Response` (backend to load balancer). Defaults to `Request`. Create-only: changing it forces resource replacement.",
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.RequiresReplace(),
							},
						},
						"rule_name": schema.StringAttribute{
							Optional:    true,
							Description: "Forwarding rule name, 1 to 255 characters.",
						},
						"tags": schema.ListNestedAttribute{
							Optional:    true,
							Description: "Tags of the rule. Create-only: the Modify API does not accept tags, so changing them forces resource replacement.",
							PlanModifiers: []planmodifier.List{
								listplanmodifier.RequiresReplace(),
							},
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
					},
				},
			},
			"client_token": schema.StringAttribute{
				Optional:    true,
				Description: "Client token used to make the CreateRules request idempotent. ASCII characters only. If omitted, the cloud service uses the request id. Create-only: changing it forces resource replacement.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"dry_run": schema.BoolAttribute{
				Optional:    true,
				Description: "Whether to only pre-check the request. When true the cloud API does not apply the change.",
			},
			"next_token": schema.StringAttribute{
				Computed:    true,
				Description: "Pagination token returned by the last DescribeRules call. Empty on the last page.",
			},
		},
	}
}

// Create handles `terraform apply` for new resources.
func (r *AlbRuleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan AlbRuleResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if plan.DryRun.ValueBool() {
		resp.Diagnostics.AddError(
			"Invalid ALB rule configuration",
			"dry_run must be false when creating a managed rule: a pre-check call never creates the cloud object, so Terraform would have no rule to manage.",
		)
		return
	}

	client, ok := r.sdkClient(&resp.Diagnostics)
	if !ok {
		return
	}

	inputs := albRuleDecodeList[albRuleInputModel](ctx, plan.Rules, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	if len(inputs) != 1 {
		resp.Diagnostics.AddError(
			"Invalid ALB rule configuration",
			fmt.Sprintf("This resource manages exactly one forwarding rule, but %d rules were configured.", len(inputs)),
		)
		return
	}

	input := inputs[0]
	request := alb.NewCreateRulesRequest()
	request.ListenerId = albRuleStringPtr(plan.ListenerId)
	request.LoadBalancerId = albRuleStringPtr(plan.LoadBalancerId)
	request.ClientToken = albRuleStringPtr(plan.ClientToken)
	request.DryRun = albRuleBoolPtr(plan.DryRun)
	request.Rules = []*alb.RuleInput{
		{
			Actions:    albRuleExpandActions(ctx, input.Actions, &resp.Diagnostics),
			Conditions: albRuleExpandConditions(ctx, input.Conditions, &resp.Diagnostics),
			Priority:   albRuleInt64Ptr(input.Priority),
			Direction:  albRuleStringPtr(input.Direction),
			RuleName:   albRuleStringPtr(input.RuleName),
			Tags:       albRuleExpandTags(ctx, input.Tags, &resp.Diagnostics),
		},
	}
	if resp.Diagnostics.HasError() {
		return
	}

	response, err := client.CreateRulesWithContext(ctx, request)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating ALB rule",
			fmt.Sprintf("CreateRules failed: %v", err),
		)
		return
	}
	if response == nil || response.Response == nil {
		resp.Diagnostics.AddError(
			"Error creating ALB rule",
			"CreateRules returned an empty response.",
		)
		return
	}
	ids := response.Response.RuleIds
	if len(ids) != 1 || ids[0] == nil || *ids[0] == "" {
		resp.Diagnostics.AddError(
			"Error creating ALB rule",
			fmt.Sprintf("CreateRules returned %d rule ids; exactly one non-empty id is required. Set dry_run to false to create a rule.", len(ids)),
		)
		return
	}

	ruleId := *ids[0]
	loadBalancerId := plan.LoadBalancerId.ValueString()
	listenerId := plan.ListenerId.ValueString()
	plan.Id = types.StringValue(albRuleComposeId(loadBalancerId, listenerId, ruleId))
	// A Computed attribute must not stay Unknown in state, so it is pinned to
	// Null until the read-back below fills it in.
	plan.NextToken = types.StringNull()

	// Persist the identity before any wait so a failed read-back still leaves
	// a destroyable resource in state.
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	item, nextToken, err := albRuleWaitAvailable(ctx, client, loadBalancerId, listenerId, ruleId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading ALB rule after create",
			fmt.Sprintf("The rule %q was created but is not readable yet. Import it with `terraform import tencentcloud_alb_rule.<name> %s` to keep it managed: %v", ruleId, plan.Id.ValueString(), err),
		)
		return
	}

	plan.NextToken = types.StringPointerValue(nextToken)
	plan.Rules = albRuleFlattenInputs(ctx, []*alb.RuleOutput{item}, plan.Rules, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Read handles `terraform refresh` and post-apply state refresh.
func (r *AlbRuleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state AlbRuleResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	loadBalancerId, listenerId, ruleId, ok := albRuleParseId(state.Id.ValueString())
	if !ok {
		resp.Diagnostics.AddError(
			"Invalid ALB rule resource ID",
			fmt.Sprintf("Resource id %q is invalid, expected `<load_balancer_id>%s<listener_id>%s<rule_id>`.", state.Id.ValueString(), tccommon.FILED_SP, tccommon.FILED_SP),
		)
		return
	}

	client, ok := r.sdkClient(&resp.Diagnostics)
	if !ok {
		return
	}

	item, nextToken, err := albRuleDescribe(ctx, client, loadBalancerId, listenerId, ruleId)
	if err != nil {
		if errors.Is(err, errAlbRuleNotFound) {
			resp.Diagnostics.AddWarning(
				"ALB rule not found",
				fmt.Sprintf("ALB rule %q no longer exists, removing it from state.", state.Id.ValueString()),
			)
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error reading ALB rule",
			fmt.Sprintf("DescribeRules failed for %s: %v", state.Id.ValueString(), err),
		)
		return
	}

	// Right after an import there is no practitioner configuration to keep, so
	// flattening adopts cloud-reported values instead of holding nulls. A
	// normal refresh must preserve null optional attributes to avoid both
	// spurious diffs and the impression that a field was cleared.
	if state.Rules.IsNull() || state.Rules.IsUnknown() {
		ctx = context.WithValue(ctx, albRuleAdoptKey{}, true)
	}

	state.LoadBalancerId = types.StringValue(loadBalancerId)
	state.ListenerId = types.StringValue(listenerId)
	state.Id = types.StringValue(state.Id.ValueString())
	state.NextToken = types.StringPointerValue(nextToken)
	state.Rules = albRuleFlattenInputs(ctx, []*alb.RuleOutput{item}, state.Rules, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update handles in-place updates through ModifyRulesAttributes.
func (r *AlbRuleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan AlbRuleResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if plan.DryRun.ValueBool() {
		resp.Diagnostics.AddError(
			"Invalid ALB rule configuration",
			"dry_run must be false when updating a managed rule: a pre-check call does not change the cloud object.",
		)
		return
	}

	var state AlbRuleResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Omitting a field in ModifyRulesAttributes is undocumented: it may either
	// keep or clear the cloud value. Removing a configured optional attribute
	// is therefore rejected instead of guessing a direction.
	cleared := albRuleUpdateClearedFields(state, plan)
	if len(cleared) > 0 {
		resp.Diagnostics.AddError(
			"ALB rule clear semantics are not confirmed",
			fmt.Sprintf("The API does not document whether omitting an attribute keeps or clears the cloud value. "+
				"Removed or emptied attributes: %s. Keep the attribute explicitly set, or recreate the resource.", strings.Join(cleared, ", ")),
		)
		return
	}

	loadBalancerId, listenerId, ruleId, ok := albRuleParseId(state.Id.ValueString())
	if !ok {
		resp.Diagnostics.AddError(
			"Invalid ALB rule resource ID",
			fmt.Sprintf("Resource id %q is invalid, expected `<load_balancer_id>%s<listener_id>%s<rule_id>`.", state.Id.ValueString(), tccommon.FILED_SP, tccommon.FILED_SP),
		)
		return
	}

	client, ok := r.sdkClient(&resp.Diagnostics)
	if !ok {
		return
	}

	inputs := albRuleDecodeList[albRuleInputModel](ctx, plan.Rules, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	if len(inputs) != 1 {
		resp.Diagnostics.AddError(
			"Invalid ALB rule configuration",
			fmt.Sprintf("This resource manages exactly one forwarding rule, but %d rules were configured.", len(inputs)),
		)
		return
	}

	input := inputs[0]
	request := alb.NewModifyRulesAttributesRequest()
	request.LoadBalancerId = albRuleStringPtr(plan.LoadBalancerId)
	request.ListenerId = albRuleStringPtr(plan.ListenerId)
	request.DryRun = albRuleBoolPtr(plan.DryRun)
	// RuleModify has no Direction and no Tags fields; both are create-only and
	// carry RequiresReplace plan modifiers.
	request.Rules = []*alb.RuleModify{
		{
			Actions:    albRuleExpandActions(ctx, input.Actions, &resp.Diagnostics),
			Conditions: albRuleExpandConditions(ctx, input.Conditions, &resp.Diagnostics),
			Priority:   albRuleInt64Ptr(input.Priority),
			RuleId:     helper.String(ruleId),
			RuleName:   albRuleStringPtr(input.RuleName),
		},
	}
	if resp.Diagnostics.HasError() {
		return
	}

	if _, err := client.ModifyRulesAttributesWithContext(ctx, request); err != nil {
		resp.Diagnostics.AddError(
			"Error updating ALB rule",
			fmt.Sprintf("ModifyRulesAttributes failed: %v", err),
		)
		return
	}

	item, nextToken, err := albRuleWaitAvailable(ctx, client, loadBalancerId, listenerId, ruleId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading ALB rule after update",
			fmt.Sprintf("ModifyRulesAttributes succeeded but the rule %q cannot be read back: %v", ruleId, err),
		)
		return
	}

	plan.Id = state.Id
	plan.NextToken = types.StringPointerValue(nextToken)
	plan.Rules = albRuleFlattenInputs(ctx, []*alb.RuleOutput{item}, plan.Rules, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete handles resource deletion.
func (r *AlbRuleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state AlbRuleResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	loadBalancerId, listenerId, ruleId, ok := albRuleParseId(state.Id.ValueString())
	if !ok {
		resp.Diagnostics.AddError(
			"Invalid ALB rule resource ID",
			fmt.Sprintf("Resource id %q is invalid, expected `<load_balancer_id>%s<listener_id>%s<rule_id>`.", state.Id.ValueString(), tccommon.FILED_SP, tccommon.FILED_SP),
		)
		return
	}

	client, ok := r.sdkClient(&resp.Diagnostics)
	if !ok {
		return
	}

	request := alb.NewDeleteRulesRequest()
	request.LoadBalancerId = helper.String(loadBalancerId)
	request.ListenerId = helper.String(listenerId)
	request.RuleIds = []*string{helper.String(ruleId)}
	// Delete must actually remove the rule. Do not forward a pre-check flag
	// from configuration: a dry-run delete would leave the cloud object behind
	// while Terraform drops it from state.

	_, err := client.DeleteRulesWithContext(ctx, request)
	if err != nil {
		// Only a documented ResourceNotFound code means the object is already
		// gone. Network, permission and nil-response failures must keep the id.
		if tccommon.IsExpectError(err, []string{"ResourceNotFound"}) {
			return
		}
		resp.Diagnostics.AddError(
			"Error deleting ALB rule",
			fmt.Sprintf("DeleteRules failed for %s: %v", state.Id.ValueString(), err),
		)
	}
}

// ImportState supports `terraform import` using the
// `<load_balancer_id>#<listener_id>#<rule_id>` form. The remaining attributes
// are filled in by the Read that follows.
func (r *AlbRuleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	loadBalancerId, listenerId, _, ok := albRuleParseId(req.ID)
	if !ok {
		resp.Diagnostics.AddError(
			"Invalid import ID",
			fmt.Sprintf("Import id %q is invalid, expected `<load_balancer_id>%s<listener_id>%s<rule_id>`.", req.ID, tccommon.FILED_SP, tccommon.FILED_SP),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("load_balancer_id"), loadBalancerId)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("listener_id"), listenerId)...)
}

// sdkClient returns the ALB SDK client bound to the shared provider client.
func (r *AlbRuleResource) sdkClient(diags *diag.Diagnostics) (*alb.Client, bool) {
	if r.Client() == nil {
		diags.AddError(
			"Provider not configured",
			"The provider client is not available. This is unexpected; please report it.",
		)
		return nil, false
	}
	client := r.Client().UseAlbV20251030Client()
	if client == nil {
		diags.AddError(
			"Provider not configured",
			"The ALB SDK client is not available. This is unexpected; please report it.",
		)
		return nil, false
	}
	return client, true
}

// albRuleDescribe locates exactly one rule by its id. The DescribeRules
// response is a list, so the returned entries are matched by RuleId instead of
// blindly taking the first element. A nil response is an API failure, while a
// successfully executed query with no match maps to errAlbRuleNotFound.
func albRuleDescribe(ctx context.Context, client *alb.Client, loadBalancerId, listenerId, ruleId string) (*alb.RuleOutput, *string, error) {
	var token *string
	seen := make(map[string]bool)

	for page := 0; page < 100; page++ {
		request := alb.NewDescribeRulesRequest()
		request.LoadBalancerId = helper.String(loadBalancerId)
		request.ListenerId = helper.String(listenerId)
		request.RuleIds = []*string{helper.String(ruleId)}
		if token != nil && *token != "" {
			request.NextToken = token
		}

		response, err := client.DescribeRulesWithContext(ctx, request)
		if err != nil {
			if tccommon.IsExpectError(err, []string{"ResourceNotFound"}) {
				return nil, nil, errAlbRuleNotFound
			}
			return nil, nil, err
		}
		if response == nil || response.Response == nil {
			return nil, nil, fmt.Errorf("DescribeRules returned an empty response for %q", ruleId)
		}

		next := response.Response.NextToken
		var match *alb.RuleOutput
		matches := 0
		for _, item := range response.Response.Rules {
			if item == nil || item.RuleId == nil || *item.RuleId != ruleId {
				continue
			}
			matches++
			match = item
		}
		if matches > 1 {
			return nil, next, fmt.Errorf("DescribeRules returned %d rules with id %q", matches, ruleId)
		}
		if matches == 1 {
			return match, next, nil
		}
		if next == nil || *next == "" {
			return nil, next, errAlbRuleNotFound
		}
		if seen[*next] {
			return nil, next, fmt.Errorf("DescribeRules returned a repeated NextToken for %q", ruleId)
		}
		seen[*next] = true
		token = next
	}
	return nil, nil, fmt.Errorf("DescribeRules pagination exceeded 100 pages for %q", ruleId)
}

// albRuleWaitAvailable polls DescribeRules until the rule is visible and no
// longer reports a transitional status. Only the documented empty-result and
// transitional-status outcomes are retried.
func albRuleWaitAvailable(ctx context.Context, client *alb.Client, loadBalancerId, listenerId, ruleId string) (*alb.RuleOutput, *string, error) {
	var rule *alb.RuleOutput
	var nextToken *string

	_, err := fwhelper.RetryFramework(ctx, tccommon.WriteRetryTimeout,
		func(err error) bool {
			return errors.Is(err, errAlbRuleNotFound) || errors.Is(err, errAlbRuleStatusPending)
		},
		func() (string, error) {
			item, token, describeErr := albRuleDescribe(ctx, client, loadBalancerId, listenerId, ruleId)
			if describeErr != nil {
				return "", describeErr
			}
			rule, nextToken = item, token
			status := ""
			if item.Status != nil {
				status = *item.Status
			}
			switch status {
			case "Provisioning", "Configuring":
				return status, fmt.Errorf("%w: rule %q is in status %q", errAlbRuleStatusPending, ruleId, status)
			default:
				return status, nil
			}
		})
	if err != nil {
		if fwhelper.IsTimeoutError(err) {
			return nil, nil, fmt.Errorf("timed out waiting for ALB rule %q: %w", ruleId, err)
		}
		return nil, nil, err
	}
	if rule == nil {
		return nil, nil, fmt.Errorf("DescribeRules did not return ALB rule %q", ruleId)
	}
	return rule, nextToken, nil
}

// albRuleComposeId builds the Terraform id from the three cloud locators.
func albRuleComposeId(loadBalancerId, listenerId, ruleId string) string {
	return strings.Join([]string{loadBalancerId, listenerId, ruleId}, tccommon.FILED_SP)
}

// albRuleParseId splits `<load_balancer_id>#<listener_id>#<rule_id>`.
func albRuleParseId(id string) (loadBalancerId, listenerId, ruleId string, ok bool) {
	parts := strings.Split(id, tccommon.FILED_SP)
	if len(parts) != 3 || parts[0] == "" || parts[1] == "" || parts[2] == "" {
		return "", "", "", false
	}
	return parts[0], parts[1], parts[2], true
}

// albRuleStringPtr converts a framework string into a request pointer while
// keeping Null and Unknown distinct from an empty string.
func albRuleStringPtr(v types.String) *string {
	if v.IsNull() || v.IsUnknown() {
		return nil
	}
	return helper.String(v.ValueString())
}

// albRuleInt64Ptr converts a framework int64 into a request pointer while
// keeping Null and Unknown distinct from zero.
func albRuleInt64Ptr(v types.Int64) *int64 {
	if v.IsNull() || v.IsUnknown() {
		return nil
	}
	return helper.Int64(v.ValueInt64())
}

// albRuleBoolPtr converts a framework bool into a request pointer while
// keeping Null and Unknown distinct from false.
func albRuleBoolPtr(v types.Bool) *bool {
	if v.IsNull() || v.IsUnknown() {
		return nil
	}
	return helper.Bool(v.ValueBool())
}

// albRuleDecodeList decodes a framework list into typed element models.
func albRuleDecodeList[T any](ctx context.Context, list types.List, diags *diag.Diagnostics) []T {
	if list.IsNull() || list.IsUnknown() {
		return nil
	}
	var out []T
	diags.Append(list.ElementsAs(ctx, &out, false)...)
	if diags.HasError() {
		return nil
	}
	return out
}

// albRuleDecodeObject decodes a framework object into a typed model. A Null
// object returns nil, which maps to an omitted cloud configuration.
func albRuleDecodeObject[T any](ctx context.Context, object types.Object, diags *diag.Diagnostics) *T {
	if object.IsNull() || object.IsUnknown() {
		return nil
	}
	var out T
	diags.Append(object.As(ctx, &out, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return nil
	}
	return &out
}

// albRuleDecodeObjectPrior decodes a previous state object, returning the zero
// model for Null or Unknown values.
func albRuleDecodeObjectPrior[T any](ctx context.Context, object types.Object, diags *diag.Diagnostics) T {
	var out T
	if object.IsNull() || object.IsUnknown() {
		return out
	}
	diags.Append(object.As(ctx, &out, basetypes.ObjectAsOptions{})...)
	return out
}

// albRuleDecodeStringSet decodes a framework string set into request pointers.
func albRuleDecodeStringSet(ctx context.Context, set types.Set, diags *diag.Diagnostics) []*string {
	if set.IsNull() || set.IsUnknown() {
		return nil
	}
	var values []types.String
	diags.Append(set.ElementsAs(ctx, &values, false)...)
	if diags.HasError() {
		return nil
	}
	out := make([]*string, 0, len(values))
	for _, value := range values {
		if value.IsNull() || value.IsUnknown() {
			continue
		}
		out = append(out, helper.String(value.ValueString()))
	}
	return out
}

// albRuleExpandActions converts the planned actions into SDK actions.
func albRuleExpandActions(ctx context.Context, list types.List, diags *diag.Diagnostics) []*alb.RuleAction {
	models := albRuleDecodeList[albRuleActionModel](ctx, list, diags)
	if diags.HasError() {
		return nil
	}
	out := make([]*alb.RuleAction, 0, len(models))
	for _, model := range models {
		action := &alb.RuleAction{
			Order: albRuleInt64Ptr(model.Order),
			Type:  albRuleStringPtr(model.Type),
		}
		if fixed := albRuleDecodeObject[albRuleFixedResponseModel](ctx, model.FixedResponseConfig, diags); fixed != nil {
			action.FixedResponseConfig = &alb.FixedResponseInfo{
				HttpCode:    albRuleInt64Ptr(fixed.HttpCode),
				Content:     albRuleStringPtr(fixed.Content),
				ContentType: albRuleStringPtr(fixed.ContentType),
			}
		}
		if insert := albRuleDecodeObject[albRuleInsertHeaderModel](ctx, model.InsertHeaderConfig, diags); insert != nil {
			action.InsertHeaderConfig = &alb.InsertHTTPHeaderInfo{
				Key:       albRuleStringPtr(insert.Key),
				Value:     albRuleStringPtr(insert.Value),
				ValueType: albRuleStringPtr(insert.ValueType),
			}
		}
		if redirect := albRuleDecodeObject[albRuleRedirectModel](ctx, model.RedirectConfig, diags); redirect != nil {
			action.RedirectConfig = &alb.HTTPRedirectInfo{
				HttpCode: albRuleInt64Ptr(redirect.HttpCode),
				Host:     albRuleStringPtr(redirect.Host),
				Path:     albRuleStringPtr(redirect.Path),
				Port:     albRuleStringPtr(redirect.Port),
				Protocol: albRuleStringPtr(redirect.Protocol),
				Query:    albRuleStringPtr(redirect.Query),
			}
		}
		if remove := albRuleDecodeObject[albRuleRemoveHeaderModel](ctx, model.RemoveHeaderConfig, diags); remove != nil {
			action.RemoveHeaderConfig = &alb.RemoveHTTPHeaderInfo{
				Key: albRuleStringPtr(remove.Key),
			}
		}
		if rewrite := albRuleDecodeObject[albRuleRewriteModel](ctx, model.RewriteConfig, diags); rewrite != nil {
			action.RewriteConfig = &alb.HTTPRewriteInfo{
				Host:  albRuleStringPtr(rewrite.Host),
				Path:  albRuleStringPtr(rewrite.Path),
				Query: albRuleStringPtr(rewrite.Query),
			}
		}
		if target := albRuleDecodeObject[albRuleTargetGroupConfigModel](ctx, model.TargetGroupConfig, diags); target != nil {
			action.TargetGroupConfig = &alb.TargetGroupConfig{
				TargetGroups:            albRuleExpandTargetGroups(ctx, target.TargetGroups, diags),
				TargetGroupStickySession: albRuleExpandStickySession(ctx, target.TargetGroupStickySession, diags),
			}
		}
		out = append(out, action)
	}
	return out
}

func albRuleExpandTargetGroups(ctx context.Context, list types.List, diags *diag.Diagnostics) []*alb.TargetGroupTuple {
	models := albRuleDecodeList[albRuleTargetGroupModel](ctx, list, diags)
	if diags.HasError() {
		return nil
	}
	out := make([]*alb.TargetGroupTuple, 0, len(models))
	for _, model := range models {
		out = append(out, &alb.TargetGroupTuple{
			TargetGroupId: albRuleStringPtr(model.TargetGroupId),
			Weight:        albRuleInt64Ptr(model.Weight),
		})
	}
	return out
}

func albRuleExpandStickySession(ctx context.Context, object types.Object, diags *diag.Diagnostics) *alb.TargetGroupStickySession {
	model := albRuleDecodeObject[albRuleStickySessionModel](ctx, object, diags)
	if model == nil {
		return nil
	}
	return &alb.TargetGroupStickySession{
		Enabled: albRuleBoolPtr(model.Enabled),
		Timeout: albRuleInt64Ptr(model.Timeout),
	}
}

// albRuleExpandConditions converts the planned conditions into SDK conditions.
func albRuleExpandConditions(ctx context.Context, list types.List, diags *diag.Diagnostics) []*alb.RuleCondition {
	models := albRuleDecodeList[albRuleConditionModel](ctx, list, diags)
	if diags.HasError() {
		return nil
	}
	out := make([]*alb.RuleCondition, 0, len(models))
	for _, model := range models {
		condition := &alb.RuleCondition{
			Type:           albRuleStringPtr(model.Type),
			HostConfig:     albRuleDecodeStringSet(ctx, model.HostConfig, diags),
			MethodConfig:   albRuleDecodeStringSet(ctx, model.MethodConfig, diags),
			PathConfig:     albRuleDecodeStringSet(ctx, model.PathConfig, diags),
			SourceIpConfig: albRuleDecodeStringSet(ctx, model.SourceIpConfig, diags),
		}
		condition.CookieConfig = albRuleExpandCookies(ctx, model.CookieConfig, diags)
		condition.QueryStringConfig = albRuleExpandQueryStrings(ctx, model.QueryStringConfig, diags)
		if header := albRuleDecodeObject[albRuleHeaderModel](ctx, model.HeaderConfig, diags); header != nil {
			condition.HeaderConfig = &alb.HTTPHeaderInfo{
				Key:    albRuleStringPtr(header.Key),
				Values: albRuleDecodeStringSet(ctx, header.Values, diags),
			}
		}
		out = append(out, condition)
	}
	return out
}

func albRuleExpandCookies(ctx context.Context, list types.List, diags *diag.Diagnostics) []*alb.HTTPCookieInfo {
	models := albRuleDecodeList[albRuleCookieModel](ctx, list, diags)
	if diags.HasError() {
		return nil
	}
	out := make([]*alb.HTTPCookieInfo, 0, len(models))
	for _, model := range models {
		out = append(out, &alb.HTTPCookieInfo{
			Key:   albRuleStringPtr(model.Key),
			Value: albRuleStringPtr(model.Value),
		})
	}
	return out
}

func albRuleExpandQueryStrings(ctx context.Context, list types.List, diags *diag.Diagnostics) []*alb.HTTPQueryStringInfo {
	models := albRuleDecodeList[albRuleQueryStringModel](ctx, list, diags)
	if diags.HasError() {
		return nil
	}
	out := make([]*alb.HTTPQueryStringInfo, 0, len(models))
	for _, model := range models {
		out = append(out, &alb.HTTPQueryStringInfo{
			Key:   albRuleStringPtr(model.Key),
			Value: albRuleStringPtr(model.Value),
		})
	}
	return out
}

func albRuleExpandTags(ctx context.Context, list types.List, diags *diag.Diagnostics) []*alb.TagInfo {
	models := albRuleDecodeList[albRuleTagModel](ctx, list, diags)
	if diags.HasError() {
		return nil
	}
	out := make([]*alb.TagInfo, 0, len(models))
	for _, model := range models {
		out = append(out, &alb.TagInfo{
			TagKey:   albRuleStringPtr(model.TagKey),
			TagValue: albRuleStringPtr(model.TagValue),
		})
	}
	return out
}

// albRuleFlattenInputs rebuilds the rules list from the Describe response.
// Cloud returned values refresh configured attributes, while omitted optional
// attributes keep their previous value. This prevents cloud defaults from
// turning into a permanent diff and prevents a Required attribute from ever
// being nulled out by a partial response.
func albRuleFlattenInputs(ctx context.Context, cloud []*alb.RuleOutput, prev types.List, diags *diag.Diagnostics) types.List {
	attrTypes := albRuleInputAttrTypes()
	objType := types.ObjectType{AttrTypes: attrTypes}
	if prev.IsNull() || prev.IsUnknown() {
		prev = types.ListNull(objType)
	}
	priors := albRuleDecodeList[albRuleInputModel](ctx, prev, diags)
	if diags.HasError() {
		return types.ListNull(objType)
	}

	out := make([]albRuleInputModel, 0, len(cloud))
	for index, item := range cloud {
		if item == nil {
			continue
		}
		var prior albRuleInputModel
		if index < len(priors) {
			prior = priors[index]
		}
		model := albRuleInputModel{
			Actions:    albRuleFlattenActions(ctx, item.Actions, prior.Actions, diags),
			Conditions: albRuleFlattenConditions(ctx, item.Conditions, prior.Conditions, diags),
			Priority:   albRuleRequiredInt64(item.Priority, prior.Priority),
			Direction:  albRuleOptionalString(ctx, item.Direction, prior.Direction),
			RuleName:   albRuleOptionalString(ctx, item.RuleName, prior.RuleName),
			Tags:       albRuleFlattenTags(ctx, item.Tags, prior.Tags, diags),
		}
		if model.Priority.IsNull() || model.Actions.IsNull() || model.Conditions.IsNull() {
			diags.AddError(
				"Malformed ALB rule response",
				fmt.Sprintf("DescribeRules returned no priority, actions or conditions for rules[%d].", index),
			)
			return types.ListNull(objType)
		}
		out = append(out, model)
	}
	// A collection the practitioner never configured stays Null. Writing an
	// empty (non-null) collection into state would cause a permanent diff.
	if len(out) == 0 {
		return prev
	}
	return albRuleMakeList(ctx, attrTypes, out, diags)
}

func albRuleFlattenActions(ctx context.Context, cloud []*alb.RuleAction, prev types.List, diags *diag.Diagnostics) types.List {
	attrTypes := albRuleActionAttrTypes()
	objType := types.ObjectType{AttrTypes: attrTypes}
	if prev.IsNull() || prev.IsUnknown() {
		prev = types.ListNull(objType)
	}
	if cloud == nil {
		return prev
	}
	priors := albRuleDecodeList[albRuleActionModel](ctx, prev, diags)
	if diags.HasError() {
		return types.ListNull(objType)
	}

	out := make([]albRuleActionModel, 0, len(cloud))
	for index, item := range cloud {
		if item == nil {
			continue
		}
		var prior albRuleActionModel
		if index < len(priors) {
			prior = priors[index]
		}
		model := albRuleActionModel{
			Order:               albRuleRequiredInt64(item.Order, prior.Order),
			Type:                albRuleRequiredString(item.Type, prior.Type),
			FixedResponseConfig: albRuleFlattenFixedResponse(ctx, item.FixedResponseConfig, prior.FixedResponseConfig, diags),
			InsertHeaderConfig:  albRuleFlattenInsertHeader(ctx, item.InsertHeaderConfig, prior.InsertHeaderConfig, diags),
			RedirectConfig:      albRuleFlattenRedirect(ctx, item.RedirectConfig, prior.RedirectConfig, diags),
			RemoveHeaderConfig:  albRuleFlattenRemoveHeader(ctx, item.RemoveHeaderConfig, prior.RemoveHeaderConfig, diags),
			RewriteConfig:       albRuleFlattenRewrite(ctx, item.RewriteConfig, prior.RewriteConfig, diags),
			TargetGroupConfig:   albRuleFlattenTargetGroupConfig(ctx, item.TargetGroupConfig, prior.TargetGroupConfig, diags),
		}
		if model.Order.IsNull() || model.Type.IsNull() {
			diags.AddError(
				"Malformed ALB rule response",
				fmt.Sprintf("DescribeRules returned no order or type for actions[%d].", index),
			)
			return types.ListNull(objType)
		}
		out = append(out, model)
	}
	// A collection the practitioner never configured stays Null. Writing an
	// empty (non-null) collection into state would cause a permanent diff.
	if len(out) == 0 {
		return prev
	}
	return albRuleMakeList(ctx, attrTypes, out, diags)
}

func albRuleFlattenFixedResponse(ctx context.Context, cloud *alb.FixedResponseInfo, prev types.Object, diags *diag.Diagnostics) types.Object {
	attrTypes := albRuleFixedResponseAttrTypes()
	if prev.IsNull() || prev.IsUnknown() {
		prev = types.ObjectNull(attrTypes)
	}
	if cloud == nil {
		return prev
	}
	prior := albRuleDecodeObjectPrior[albRuleFixedResponseModel](ctx, prev, diags)
	model := albRuleFixedResponseModel{
		HttpCode:    albRuleRequiredInt64(cloud.HttpCode, prior.HttpCode),
		Content:     albRuleOptionalString(ctx, cloud.Content, prior.Content),
		ContentType: albRuleOptionalString(ctx, cloud.ContentType, prior.ContentType),
	}
	if model.HttpCode.IsNull() {
		diags.AddError("Malformed ALB rule response", "DescribeRules returned no http_code for fixed_response_config.")
		return types.ObjectNull(attrTypes)
	}
	return albRuleMakeObject(ctx, attrTypes, model, diags)
}

func albRuleFlattenInsertHeader(ctx context.Context, cloud *alb.InsertHTTPHeaderInfo, prev types.Object, diags *diag.Diagnostics) types.Object {
	attrTypes := albRuleInsertHeaderAttrTypes()
	if prev.IsNull() || prev.IsUnknown() {
		prev = types.ObjectNull(attrTypes)
	}
	if cloud == nil {
		return prev
	}
	prior := albRuleDecodeObjectPrior[albRuleInsertHeaderModel](ctx, prev, diags)
	model := albRuleInsertHeaderModel{
		Key:       albRuleRequiredString(cloud.Key, prior.Key),
		Value:     albRuleRequiredString(cloud.Value, prior.Value),
		ValueType: albRuleRequiredString(cloud.ValueType, prior.ValueType),
	}
	if model.Key.IsNull() || model.Value.IsNull() || model.ValueType.IsNull() {
		diags.AddError("Malformed ALB rule response", "DescribeRules returned an incomplete insert_header_config.")
		return types.ObjectNull(attrTypes)
	}
	return albRuleMakeObject(ctx, attrTypes, model, diags)
}

func albRuleFlattenRedirect(ctx context.Context, cloud *alb.HTTPRedirectInfo, prev types.Object, diags *diag.Diagnostics) types.Object {
	attrTypes := albRuleRedirectAttrTypes()
	if prev.IsNull() || prev.IsUnknown() {
		prev = types.ObjectNull(attrTypes)
	}
	if cloud == nil {
		return prev
	}
	prior := albRuleDecodeObjectPrior[albRuleRedirectModel](ctx, prev, diags)
	model := albRuleRedirectModel{
		HttpCode: albRuleRequiredInt64(cloud.HttpCode, prior.HttpCode),
		Host:     albRuleOptionalString(ctx, cloud.Host, prior.Host),
		Path:     albRuleOptionalString(ctx, cloud.Path, prior.Path),
		Port:     albRuleOptionalString(ctx, cloud.Port, prior.Port),
		Protocol: albRuleOptionalString(ctx, cloud.Protocol, prior.Protocol),
		Query:    albRuleOptionalString(ctx, cloud.Query, prior.Query),
	}
	if model.HttpCode.IsNull() {
		diags.AddError("Malformed ALB rule response", "DescribeRules returned no http_code for redirect_config.")
		return types.ObjectNull(attrTypes)
	}
	return albRuleMakeObject(ctx, attrTypes, model, diags)
}

func albRuleFlattenRemoveHeader(ctx context.Context, cloud *alb.RemoveHTTPHeaderInfo, prev types.Object, diags *diag.Diagnostics) types.Object {
	attrTypes := albRuleRemoveHeaderAttrTypes()
	if prev.IsNull() || prev.IsUnknown() {
		prev = types.ObjectNull(attrTypes)
	}
	if cloud == nil {
		return prev
	}
	prior := albRuleDecodeObjectPrior[albRuleRemoveHeaderModel](ctx, prev, diags)
	model := albRuleRemoveHeaderModel{Key: albRuleRequiredString(cloud.Key, prior.Key)}
	if model.Key.IsNull() {
		diags.AddError("Malformed ALB rule response", "DescribeRules returned no key for remove_header_config.")
		return types.ObjectNull(attrTypes)
	}
	return albRuleMakeObject(ctx, attrTypes, model, diags)
}

func albRuleFlattenRewrite(ctx context.Context, cloud *alb.HTTPRewriteInfo, prev types.Object, diags *diag.Diagnostics) types.Object {
	attrTypes := albRuleRewriteAttrTypes()
	if prev.IsNull() || prev.IsUnknown() {
		prev = types.ObjectNull(attrTypes)
	}
	if cloud == nil {
		return prev
	}
	prior := albRuleDecodeObjectPrior[albRuleRewriteModel](ctx, prev, diags)
	model := albRuleRewriteModel{
		Host:  albRuleOptionalString(ctx, cloud.Host, prior.Host),
		Path:  albRuleOptionalString(ctx, cloud.Path, prior.Path),
		Query: albRuleOptionalString(ctx, cloud.Query, prior.Query),
	}
	return albRuleMakeObject(ctx, attrTypes, model, diags)
}

func albRuleFlattenTargetGroupConfig(ctx context.Context, cloud *alb.TargetGroupConfig, prev types.Object, diags *diag.Diagnostics) types.Object {
	attrTypes := albRuleTargetGroupConfigAttrTypes()
	if prev.IsNull() || prev.IsUnknown() {
		prev = types.ObjectNull(attrTypes)
	}
	if cloud == nil {
		return prev
	}
	prior := albRuleDecodeObjectPrior[albRuleTargetGroupConfigModel](ctx, prev, diags)
	model := albRuleTargetGroupConfigModel{
		TargetGroups:            albRuleFlattenTargetGroups(ctx, cloud.TargetGroups, prior.TargetGroups, diags),
		TargetGroupStickySession: albRuleFlattenStickySession(ctx, cloud.TargetGroupStickySession, prior.TargetGroupStickySession, diags),
	}
	return albRuleMakeObject(ctx, attrTypes, model, diags)
}

func albRuleFlattenTargetGroups(ctx context.Context, cloud []*alb.TargetGroupTuple, prev types.List, diags *diag.Diagnostics) types.List {
	attrTypes := albRuleTargetGroupAttrTypes()
	objType := types.ObjectType{AttrTypes: attrTypes}
	if prev.IsNull() || prev.IsUnknown() {
		prev = types.ListNull(objType)
	}
	if cloud == nil {
		return prev
	}
	priors := albRuleDecodeList[albRuleTargetGroupModel](ctx, prev, diags)
	if diags.HasError() {
		return types.ListNull(objType)
	}
	out := make([]albRuleTargetGroupModel, 0, len(cloud))
	for index, item := range cloud {
		if item == nil {
			continue
		}
		var prior albRuleTargetGroupModel
		if index < len(priors) {
			prior = priors[index]
		}
		out = append(out, albRuleTargetGroupModel{
			TargetGroupId: albRuleRequiredString(item.TargetGroupId, prior.TargetGroupId),
			Weight:        albRuleOptionalInt64(ctx, item.Weight, prior.Weight),
		})
	}
	// A collection the practitioner never configured stays Null. Writing an
	// empty (non-null) collection into state would cause a permanent diff.
	if len(out) == 0 {
		return prev
	}
	return albRuleMakeList(ctx, attrTypes, out, diags)
}

func albRuleFlattenStickySession(ctx context.Context, cloud *alb.TargetGroupStickySession, prev types.Object, diags *diag.Diagnostics) types.Object {
	attrTypes := albRuleStickySessionAttrTypes()
	if prev.IsNull() || prev.IsUnknown() {
		prev = types.ObjectNull(attrTypes)
	}
	if cloud == nil {
		return prev
	}
	prior := albRuleDecodeObjectPrior[albRuleStickySessionModel](ctx, prev, diags)
	model := albRuleStickySessionModel{
		Enabled: albRuleRequiredBool(cloud.Enabled, prior.Enabled),
		Timeout: albRuleOptionalInt64(ctx, cloud.Timeout, prior.Timeout),
	}
	if model.Enabled.IsNull() {
		diags.AddError("Malformed ALB rule response", "DescribeRules returned no enabled for target_group_sticky_session.")
		return types.ObjectNull(attrTypes)
	}
	return albRuleMakeObject(ctx, attrTypes, model, diags)
}

func albRuleFlattenConditions(ctx context.Context, cloud []*alb.RuleCondition, prev types.List, diags *diag.Diagnostics) types.List {
	attrTypes := albRuleConditionAttrTypes()
	objType := types.ObjectType{AttrTypes: attrTypes}
	if prev.IsNull() || prev.IsUnknown() {
		prev = types.ListNull(objType)
	}
	if cloud == nil {
		return prev
	}
	priors := albRuleDecodeList[albRuleConditionModel](ctx, prev, diags)
	if diags.HasError() {
		return types.ListNull(objType)
	}
	out := make([]albRuleConditionModel, 0, len(cloud))
	for index, item := range cloud {
		if item == nil {
			continue
		}
		var prior albRuleConditionModel
		if index < len(priors) {
			prior = priors[index]
		}
		model := albRuleConditionModel{
			Type:              albRuleRequiredString(item.Type, prior.Type),
			CookieConfig:      albRuleFlattenCookies(ctx, item.CookieConfig, prior.CookieConfig, diags),
			HeaderConfig:      albRuleFlattenHeader(ctx, item.HeaderConfig, prior.HeaderConfig, diags),
			HostConfig:        albRuleFlattenStringSet(item.HostConfig, prior.HostConfig, diags),
			MethodConfig:      albRuleFlattenStringSet(item.MethodConfig, prior.MethodConfig, diags),
			PathConfig:        albRuleFlattenStringSet(item.PathConfig, prior.PathConfig, diags),
			QueryStringConfig: albRuleFlattenQueryStrings(ctx, item.QueryStringConfig, prior.QueryStringConfig, diags),
			SourceIpConfig:    albRuleFlattenStringSet(item.SourceIpConfig, prior.SourceIpConfig, diags),
		}
		if model.Type.IsNull() {
			diags.AddError("Malformed ALB rule response", fmt.Sprintf("DescribeRules returned no type for conditions[%d].", index))
			return types.ListNull(objType)
		}
		out = append(out, model)
	}
	// A collection the practitioner never configured stays Null. Writing an
	// empty (non-null) collection into state would cause a permanent diff.
	if len(out) == 0 {
		return prev
	}
	return albRuleMakeList(ctx, attrTypes, out, diags)
}

func albRuleFlattenCookies(ctx context.Context, cloud []*alb.HTTPCookieInfo, prev types.List, diags *diag.Diagnostics) types.List {
	attrTypes := albRuleCookieAttrTypes()
	objType := types.ObjectType{AttrTypes: attrTypes}
	if prev.IsNull() || prev.IsUnknown() {
		prev = types.ListNull(objType)
	}
	if cloud == nil {
		return prev
	}
	priors := albRuleDecodeList[albRuleCookieModel](ctx, prev, diags)
	if diags.HasError() {
		return types.ListNull(objType)
	}
	out := make([]albRuleCookieModel, 0, len(cloud))
	for index, item := range cloud {
		if item == nil {
			continue
		}
		var prior albRuleCookieModel
		if index < len(priors) {
			prior = priors[index]
		}
		out = append(out, albRuleCookieModel{
			Key:   albRuleRequiredString(item.Key, prior.Key),
			Value: albRuleRequiredString(item.Value, prior.Value),
		})
	}
	// A collection the practitioner never configured stays Null. Writing an
	// empty (non-null) collection into state would cause a permanent diff.
	if len(out) == 0 {
		return prev
	}
	return albRuleMakeList(ctx, attrTypes, out, diags)
}

func albRuleFlattenQueryStrings(ctx context.Context, cloud []*alb.HTTPQueryStringInfo, prev types.List, diags *diag.Diagnostics) types.List {
	attrTypes := albRuleQueryStringAttrTypes()
	objType := types.ObjectType{AttrTypes: attrTypes}
	if prev.IsNull() || prev.IsUnknown() {
		prev = types.ListNull(objType)
	}
	if cloud == nil {
		return prev
	}
	priors := albRuleDecodeList[albRuleQueryStringModel](ctx, prev, diags)
	if diags.HasError() {
		return types.ListNull(objType)
	}
	out := make([]albRuleQueryStringModel, 0, len(cloud))
	for index, item := range cloud {
		if item == nil {
			continue
		}
		var prior albRuleQueryStringModel
		if index < len(priors) {
			prior = priors[index]
		}
		out = append(out, albRuleQueryStringModel{
			Key:   albRuleRequiredString(item.Key, prior.Key),
			Value: albRuleRequiredString(item.Value, prior.Value),
		})
	}
	// A collection the practitioner never configured stays Null. Writing an
	// empty (non-null) collection into state would cause a permanent diff.
	if len(out) == 0 {
		return prev
	}
	return albRuleMakeList(ctx, attrTypes, out, diags)
}

func albRuleFlattenHeader(ctx context.Context, cloud *alb.HTTPHeaderInfo, prev types.Object, diags *diag.Diagnostics) types.Object {
	attrTypes := albRuleHeaderAttrTypes()
	if prev.IsNull() || prev.IsUnknown() {
		prev = types.ObjectNull(attrTypes)
	}
	if cloud == nil {
		return prev
	}
	prior := albRuleDecodeObjectPrior[albRuleHeaderModel](ctx, prev, diags)
	model := albRuleHeaderModel{
		Key:    albRuleRequiredString(cloud.Key, prior.Key),
		Values: albRuleFlattenRequiredStringSet(cloud.Values, prior.Values, diags),
	}
	if model.Key.IsNull() {
		diags.AddError("Malformed ALB rule response", "DescribeRules returned no key for header_config.")
		return types.ObjectNull(attrTypes)
	}
	return albRuleMakeObject(ctx, attrTypes, model, diags)
}

func albRuleFlattenTags(ctx context.Context, cloud []*alb.TagInfo, prev types.List, diags *diag.Diagnostics) types.List {
	attrTypes := albRuleTagAttrTypes()
	objType := types.ObjectType{AttrTypes: attrTypes}
	if prev.IsNull() || prev.IsUnknown() {
		prev = types.ListNull(objType)
	}
	if cloud == nil {
		return prev
	}
	out := make([]albRuleTagModel, 0, len(cloud))
	for _, item := range cloud {
		if item == nil {
			continue
		}
		out = append(out, albRuleTagModel{
			TagKey:   albRuleRequiredString(item.TagKey, types.StringNull()),
			TagValue: albRuleRequiredString(item.TagValue, types.StringNull()),
		})
	}
	// A collection the practitioner never configured stays Null. Writing an
	// empty (non-null) collection into state would cause a permanent diff.
	if len(out) == 0 {
		return prev
	}
	return albRuleMakeList(ctx, attrTypes, out, diags)
}

func albRuleFlattenStringSet(cloud []*string, prev types.Set, diags *diag.Diagnostics) types.Set {
	if prev.IsNull() || prev.IsUnknown() {
		prev = types.SetNull(types.StringType)
	}
	if cloud == nil {
		return prev
	}
	if len(cloud) == 0 && !prev.IsNull() {
		return prev
	}
	elements := make([]attr.Value, 0, len(cloud))
	for _, item := range cloud {
		if item != nil {
			elements = append(elements, types.StringValue(*item))
		}
	}
	// An omitted optional set must stay null (especially right after an
	// import), while a configured set that the cloud emptied must land as an
	// empty set so the drift is visible.
	if len(elements) == 0 && (prev.IsNull() || prev.IsUnknown()) {
		return prev
	}
	set, setDiags := types.SetValue(types.StringType, elements)
	diags.Append(setDiags...)
	if setDiags.HasError() {
		return types.SetNull(types.StringType)
	}
	return set
}

func albRuleFlattenRequiredStringSet(cloud []*string, prev types.Set, diags *diag.Diagnostics) types.Set {
	set := albRuleFlattenStringSet(cloud, prev, diags)
	if set.IsNull() || set.IsUnknown() {
		diags.AddError("Malformed ALB rule response", "DescribeRules returned no values for header_config.")
		empty, emptyDiags := types.SetValue(types.StringType, []attr.Value{})
		diags.Append(emptyDiags...)
		return empty
	}
	return set
}

// albRuleMakeObject converts a typed model into a framework object value.
func albRuleMakeObject[T any](ctx context.Context, attrTypes map[string]attr.Type, model T, diags *diag.Diagnostics) types.Object {
	object, objectDiags := types.ObjectValueFrom(ctx, attrTypes, model)
	diags.Append(objectDiags...)
	if objectDiags.HasError() {
		return types.ObjectNull(attrTypes)
	}
	return object
}

// albRuleMakeList converts typed models into a framework list value.
func albRuleMakeList[T any](ctx context.Context, attrTypes map[string]attr.Type, models []T, diags *diag.Diagnostics) types.List {
	objType := types.ObjectType{AttrTypes: attrTypes}
	elements := make([]attr.Value, 0, len(models))
	for _, model := range models {
		object, objectDiags := types.ObjectValueFrom(ctx, attrTypes, model)
		diags.Append(objectDiags...)
		if objectDiags.HasError() {
			return types.ListNull(objType)
		}
		elements = append(elements, object)
	}
	list, listDiags := types.ListValue(objType, elements)
	diags.Append(listDiags...)
	if listDiags.HasError() {
		return types.ListNull(objType)
	}
	return list
}

// albRuleOptionalString refreshes an optional attribute. A value that was not
// configured stays Null even if the cloud echoes a default back, unless the
// first post-import Read is adopting the cloud configuration.
func albRuleOptionalString(ctx context.Context, cloud *string, prev types.String) types.String {
	if (prev.IsNull() || prev.IsUnknown()) && !albRuleAdopting(ctx) {
		return prev
	}
	if cloud == nil {
		return prev
	}
	return types.StringValue(*cloud)
}

func albRuleOptionalInt64(ctx context.Context, cloud *int64, prev types.Int64) types.Int64 {
	if (prev.IsNull() || prev.IsUnknown()) && !albRuleAdopting(ctx) {
		return prev
	}
	if cloud == nil {
		return prev
	}
	return types.Int64Value(*cloud)
}

// albRuleAdopting reports whether cloud values should be adopted into state.
func albRuleAdopting(ctx context.Context) bool {
	adopt, ok := ctx.Value(albRuleAdoptKey{}).(bool)
	return ok && adopt
}

func albRuleRequiredString(cloud *string, prev types.String) types.String {
	if cloud != nil {
		return types.StringValue(*cloud)
	}
	return prev
}

func albRuleRequiredInt64(cloud *int64, prev types.Int64) types.Int64 {
	if cloud != nil {
		return types.Int64Value(*cloud)
	}
	return prev
}

func albRuleRequiredBool(cloud *bool, prev types.Bool) types.Bool {
	if cloud != nil {
		return types.BoolValue(*cloud)
	}
	return prev
}

// albRuleClearedPaths reports nested attribute paths whose configured value
// disappeared from the plan. Removing an optional value has undocumented cloud
// semantics, so Update refuses to proceed when such a path exists.
func albRuleClearedPaths(prev, next attr.Value, current string) []string {
	if !albRuleValuePresent(prev) {
		return nil
	}
	if !albRuleValuePresent(next) {
		return []string{current}
	}

	switch previous := prev.(type) {
	case types.Object:
		following, ok := next.(types.Object)
		if !ok {
			return nil
		}
		nextAttributes := following.Attributes()
		var cleared []string
		for name, value := range previous.Attributes() {
			cleared = append(cleared, albRuleClearedPaths(value, nextAttributes[name], current+"."+name)...)
		}
		return cleared
	case types.List:
		following, ok := next.(types.List)
		if !ok {
			return nil
		}
		nextElements := following.Elements()
		var cleared []string
		for index, value := range previous.Elements() {
			elementPath := fmt.Sprintf("%s[%d]", current, index)
			if index >= len(nextElements) {
				cleared = append(cleared, elementPath)
				continue
			}
			cleared = append(cleared, albRuleClearedPaths(value, nextElements[index], elementPath)...)
		}
		return cleared
	case types.Set:
		following, ok := next.(types.Set)
		if !ok {
			return nil
		}
		if len(previous.Elements()) > 0 && len(following.Elements()) < len(previous.Elements()) {
			return []string{current}
		}
	}
	return nil
}

func albRuleValuePresent(value attr.Value) bool {
	return value != nil && !value.IsNull() && !value.IsUnknown()
}

func albRuleUpdateClearedFields(state, plan AlbRuleResourceModel) []string {
	pairs := []struct {
		name           string
		before, after  attr.Value
	}{
		{"load_balancer_id", state.LoadBalancerId, plan.LoadBalancerId},
		{"listener_id", state.ListenerId, plan.ListenerId},
		{"rules", state.Rules, plan.Rules},
		{"client_token", state.ClientToken, plan.ClientToken},
		{"dry_run", state.DryRun, plan.DryRun},
	}
	var cleared []string
	for _, pair := range pairs {
		cleared = append(cleared, albRuleClearedPaths(pair.before, pair.after, pair.name)...)
	}
	return cleared
}

func albRuleInputAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"actions":    types.ListType{ElemType: types.ObjectType{AttrTypes: albRuleActionAttrTypes()}},
		"conditions": types.ListType{ElemType: types.ObjectType{AttrTypes: albRuleConditionAttrTypes()}},
		"priority":   types.Int64Type,
		"direction":  types.StringType,
		"rule_name":  types.StringType,
		"tags":       types.ListType{ElemType: types.ObjectType{AttrTypes: albRuleTagAttrTypes()}},
	}
}

func albRuleActionAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"order":                 types.Int64Type,
		"type":                  types.StringType,
		"fixed_response_config": types.ObjectType{AttrTypes: albRuleFixedResponseAttrTypes()},
		"insert_header_config":  types.ObjectType{AttrTypes: albRuleInsertHeaderAttrTypes()},
		"redirect_config":       types.ObjectType{AttrTypes: albRuleRedirectAttrTypes()},
		"remove_header_config":  types.ObjectType{AttrTypes: albRuleRemoveHeaderAttrTypes()},
		"rewrite_config":        types.ObjectType{AttrTypes: albRuleRewriteAttrTypes()},
		"target_group_config":   types.ObjectType{AttrTypes: albRuleTargetGroupConfigAttrTypes()},
	}
}

func albRuleFixedResponseAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"http_code":    types.Int64Type,
		"content":      types.StringType,
		"content_type": types.StringType,
	}
}

func albRuleInsertHeaderAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"key":        types.StringType,
		"value":      types.StringType,
		"value_type": types.StringType,
	}
}

func albRuleRedirectAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"http_code": types.Int64Type,
		"host":      types.StringType,
		"path":      types.StringType,
		"port":      types.StringType,
		"protocol":  types.StringType,
		"query":     types.StringType,
	}
}

func albRuleRemoveHeaderAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"key": types.StringType,
	}
}

func albRuleRewriteAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"host":  types.StringType,
		"path":  types.StringType,
		"query": types.StringType,
	}
}

func albRuleTargetGroupConfigAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"target_groups":             types.ListType{ElemType: types.ObjectType{AttrTypes: albRuleTargetGroupAttrTypes()}},
		"target_group_sticky_session": types.ObjectType{AttrTypes: albRuleStickySessionAttrTypes()},
	}
}

func albRuleTargetGroupAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"target_group_id": types.StringType,
		"weight":          types.Int64Type,
	}
}

func albRuleStickySessionAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"enabled": types.BoolType,
		"timeout": types.Int64Type,
	}
}

func albRuleConditionAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"type":                types.StringType,
		"cookie_config":       types.ListType{ElemType: types.ObjectType{AttrTypes: albRuleCookieAttrTypes()}},
		"header_config":       types.ObjectType{AttrTypes: albRuleHeaderAttrTypes()},
		"host_config":         types.SetType{ElemType: types.StringType},
		"method_config":       types.SetType{ElemType: types.StringType},
		"path_config":         types.SetType{ElemType: types.StringType},
		"query_string_config": types.ListType{ElemType: types.ObjectType{AttrTypes: albRuleQueryStringAttrTypes()}},
		"source_ip_config":    types.SetType{ElemType: types.StringType},
	}
}

func albRuleCookieAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"key":   types.StringType,
		"value": types.StringType,
	}
}

func albRuleHeaderAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"key":    types.StringType,
		"values": types.SetType{ElemType: types.StringType},
	}
}

func albRuleQueryStringAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"key":   types.StringType,
		"value": types.StringType,
	}
}

func albRuleTagAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"tag_key":   types.StringType,
		"tag_value": types.StringType,
	}
}
