package alb

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	albv20251030 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/alb/v20251030"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/fw"
	fwHelper "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/helper"
)

var (
	_ resource.ResourceWithConfigure   = &AlbListenerResource{}
	_ resource.ResourceWithImportState = &AlbListenerResource{}
)

// NewAlbListenerResource is the factory referenced by registry.go.
func NewAlbListenerResource() resource.Resource {
	return &AlbListenerResource{}
}

// AlbListenerResource manages the forwarding rules of an ALB listener.
type AlbListenerResource struct {
	fw.ResourceWithConfigure
	fw.WithImportByID
}

// AlbListenerResourceModel maps the resource schema attributes.
type AlbListenerResourceModel struct {
	Id             types.String `tfsdk:"id"`
	ListenerId     types.String `tfsdk:"listener_id"`
	LoadBalancerId types.String `tfsdk:"load_balancer_id"`
	Rules          types.List   `tfsdk:"rules"`
}

// albRuleModel maps one element of the rules list.
type albRuleModel struct {
	Actions    types.List   `tfsdk:"actions"`
	Conditions types.List   `tfsdk:"conditions"`
	Priority   types.Int64  `tfsdk:"priority"`
	Direction  types.String `tfsdk:"direction"`
	RuleName   types.String `tfsdk:"rule_name"`
	Tags       types.List   `tfsdk:"tags"`
	RuleId     types.String `tfsdk:"rule_id"`
	CreateTime types.String `tfsdk:"create_time"`
	ModifyTime types.String `tfsdk:"modify_time"`
	Status     types.String `tfsdk:"status"`
}

// albRuleActionModel maps one element of a rule's actions list.
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

// albRuleConditionModel maps one element of a rule's conditions list.
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

func (r *AlbListenerResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "tencentcloud_alb_listener"
}

func (r *AlbListenerResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a resource to manage the forwarding rules of an ALB listener.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "ID of the resource, taken from the first rule ID returned by CreateRules.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"listener_id": schema.StringAttribute{
				Required:    true,
				Description: "Listener ID, in the format of lst- followed by 8 alphanumeric characters.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"load_balancer_id": schema.StringAttribute{
				Required:    true,
				Description: "Load balancer instance ID, in the format of alb- followed by 8 alphanumeric characters.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"rules": schema.ListNestedAttribute{
				Required:    true,
				Description: "Forwarding rule list.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: albRuleSchemaAttributes(),
				},
			},
		},
	}
}

func albRuleSchemaAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"actions": schema.ListNestedAttribute{
			Required:    true,
			Description: "Forwarding rule action list.",
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"order": schema.Int64Attribute{
						Required:    true,
						Description: "Action execution order. Values cannot be duplicated and actions are executed in ascending order. Value range: 1 ~ 50000.",
					},
					"type": schema.StringAttribute{
						Required:    true,
						Description: "Forwarding action type. Valid values: TargetGroup, Redirect, FixedResponse, Rewrite, InsertHeader, RemoveHeader. An action list must contain one of TargetGroup, Redirect or FixedResponse, placed last in execution order.",
					},
					"fixed_response_config": schema.SingleNestedAttribute{
						Optional:    true,
						Description: "Fixed response content configuration.",
						Attributes: map[string]schema.Attribute{
							"http_code": schema.Int64Attribute{
								Required:    true,
								Description: "Returned HTTP response code. 2xx, 4xx and 5xx are supported.",
							},
							"content": schema.StringAttribute{
								Optional:    true,
								Description: "Returned fixed content. Only ASCII characters are supported, up to 1 KB.",
							},
							"content_type": schema.StringAttribute{
								Optional:    true,
								Description: "Format of the returned fixed content. Valid values: text/plain, text/css, text/html, application/javascript, application/json.",
							},
						},
					},
					"insert_header_config": schema.SingleNestedAttribute{
						Optional:    true,
						Description: "Insert HTTP header configuration.",
						Attributes: map[string]schema.Attribute{
							"key": schema.StringAttribute{
								Required:    true,
								Description: "Key of the HTTP header to insert, 1 ~ 40 characters, supporting a-z A-Z 0-9 - _. Chinese characters and Cookie, Host, Content-Length, Connection, Upgrade, transfer-encoding, keep-alive, te, authority, x-forwarded-for, x-forwarded-proto, x-forwarded-host and x-forwarded-port are not supported.",
							},
							"value": schema.StringAttribute{
								Required:    true,
								Description: "HTTP header value. When ValueType is SystemDefined, valid values are ClientPort, ClientIp, Protocol and CLBPort. When ValueType is UserDefined, 1 ~ 128 printable characters without double quotes, leading or trailing spaces, or a trailing backslash. When ValueType is ReferenceHeader, references a request header with the same constraints.",
							},
							"value_type": schema.StringAttribute{
								Required:    true,
								Description: "HTTP header value type. Valid values: SystemDefined, UserDefined, ReferenceHeader.",
							},
						},
					},
					"redirect_config": schema.SingleNestedAttribute{
						Optional:    true,
						Description: "Redirect configuration. Except for HttpCode, the other settings cannot all use default values.",
						Attributes: map[string]schema.Attribute{
							"http_code": schema.Int64Attribute{
								Required:    true,
								Description: "HTTP redirect code. 301, 302, 303, 307 and 308 are supported.",
							},
							"host": schema.StringAttribute{
								Optional:    true,
								Description: "Redirect host address, default ${host}. 3 ~ 128 characters, supporting a-z 0-9 _ . -.",
							},
							"path": schema.StringAttribute{
								Optional:    true,
								Description: "Redirect path, default ${path}. 1 ~ 128 characters, supporting a-z A-Z 0-9 ? = _ . - / : .",
							},
							"port": schema.StringAttribute{
								Optional:    true,
								Description: "Redirect port, default ${port}. Value range: 1 ~ 65535.",
							},
							"protocol": schema.StringAttribute{
								Optional:    true,
								Description: "Redirect protocol. Valid values: HTTP, HTTPS. Default ${protocol}.",
							},
							"query": schema.StringAttribute{
								Optional:    true,
								Description: "Redirect query string, default ${query}. 1 ~ 128 printable characters; #[]{}<> and spaces are not supported.",
							},
						},
					},
					"remove_header_config": schema.SingleNestedAttribute{
						Optional:    true,
						Description: "Remove HTTP header configuration.",
						Attributes: map[string]schema.Attribute{
							"key": schema.StringAttribute{
								Required:    true,
								Description: "Key of the HTTP header to remove, 1 ~ 40 characters, supporting a-z A-Z 0-9 - _.",
							},
						},
					},
					"rewrite_config": schema.SingleNestedAttribute{
						Optional:    true,
						Description: "Rewrite configuration.",
						Attributes: map[string]schema.Attribute{
							"host": schema.StringAttribute{
								Optional:    true,
								Description: "Rewrite host address, default ${host}. 3 ~ 128 characters, supporting a-z 0-9 _ . -.",
							},
							"path": schema.StringAttribute{
								Optional:    true,
								Description: "Rewrite path, default ${path}. 1 ~ 128 characters, supporting a-z A-Z 0-9 ? = _ . - / : .",
							},
							"query": schema.StringAttribute{
								Optional:    true,
								Description: "Rewrite query string, default ${query}. 1 ~ 128 printable characters; #[]{}<> and spaces are not supported.",
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
											Description: "Target group ID, in the format of lbtg- followed by 8 alphanumeric characters.",
										},
										"weight": schema.Int64Attribute{
											Optional:    true,
											Description: "Weight, value range [0, 100], default 10.",
										},
									},
								},
							},
							"target_group_sticky_session": schema.SingleNestedAttribute{
								Optional:    true,
								Description: "Sticky session between target groups.",
								Attributes: map[string]schema.Attribute{
									"enabled": schema.BoolAttribute{
										Required:    true,
										Description: "Whether to enable sticky sessions. Disabled by default.",
									},
									"timeout": schema.Int64Attribute{
										Optional:    true,
										Description: "Timeout in seconds, value range 1 ~ 86400, default 1000.",
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
			Description: "Forwarding rule condition list.",
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{
						Required:    true,
						Description: "Forwarding condition type. Valid values: Host, Path, Header, QueryString, Method, Cookie, SourceIp.",
					},
					"cookie_config": schema.ListNestedAttribute{
						Optional:    true,
						Description: "Cookie configuration.",
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"key": schema.StringAttribute{
									Required:    true,
									Description: "Cookie key, 1 ~ 64 characters, supporting letters, digits and underscores.",
								},
								"value": schema.StringAttribute{
									Required:    true,
									Description: "Cookie value, 1 ~ 128 printable characters.",
								},
							},
						},
					},
					"header_config": schema.SingleNestedAttribute{
						Optional:    true,
						Description: "HTTP header configuration.",
						Attributes: map[string]schema.Attribute{
							"key": schema.StringAttribute{
								Required:    true,
								Description: "HTTP header key, 1 ~ 40 characters, supporting a-z A-Z 0-9 - _.",
							},
							"values": schema.ListAttribute{
								Required:    true,
								ElementType: types.StringType,
								Description: "HTTP header values, 1 ~ 128 printable characters each.",
							},
						},
					},
					"host_config": schema.ListAttribute{
						Optional:    true,
						ElementType: types.StringType,
						Description: "Host names. A host can be configured only once per rule; 3 ~ 128 characters, supporting exact, regular-expression and wildcard matching.",
					},
					"method_config": schema.ListAttribute{
						Optional:    true,
						ElementType: types.StringType,
						Description: "Request methods. Valid values: HEAD, GET, POST, OPTIONS, PUT, PATCH, DELETE.",
					},
					"path_config": schema.ListAttribute{
						Optional:    true,
						ElementType: types.StringType,
						Description: "Forwarding paths, 1 ~ 128 characters, supporting exact, regular-expression and wildcard matching.",
					},
					"query_string_config": schema.ListNestedAttribute{
						Optional:    true,
						Description: "Query string configuration.",
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"key": schema.StringAttribute{
									Required:    true,
									Description: "Query string key, 1 ~ 16 characters; printable characters except spaces and #[]{}|\\<>&. * and ? are supported as wildcards.",
								},
								"value": schema.StringAttribute{
									Required:    true,
									Description: "Query string value, 1 ~ 128 characters; printable characters except spaces and #[]{}|\\<>&. * and ? are supported as wildcards.",
								},
							},
						},
					},
					"source_ip_config": schema.ListAttribute{
						Optional:    true,
						ElementType: types.StringType,
						Description: "Source IP match configuration in CIDR format, for example x.x.x.x/32 or x.x.x.x/24.",
					},
				},
			},
		},
		"priority": schema.Int64Attribute{
			Required:    true,
			Description: "Priority. Smaller values indicate higher priority; values cannot be duplicated. Value range: 1 ~ 10000.",
		},
		"direction": schema.StringAttribute{
			Optional:    true,
			Description: "Forwarding rule direction. Request: client to load balancer; Response: backend server to load balancer. Default Request.",
		},
		"rule_name": schema.StringAttribute{
			Optional:    true,
			Description: "Forwarding rule name, 1 ~ 255 characters, allowing digits, letters, Chinese characters, periods (.), underscores (_) and hyphens (-).",
		},
		"tags": schema.ListNestedAttribute{
			Optional:    true,
			Description: "Tags.",
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
		"rule_id": schema.StringAttribute{
			Computed:    true,
			Description: "Forwarding rule ID, in the format of rule- followed by 8 alphanumeric characters.",
		},
		"create_time": schema.StringAttribute{
			Computed:    true,
			Description: "Creation time.",
		},
		"modify_time": schema.StringAttribute{
			Computed:    true,
			Description: "Last modification time.",
		},
		"status": schema.StringAttribute{
			Computed:    true,
			Description: "Forwarding rule status. Provisioning: creating; Active: running; Configuring: configuring.",
		},
	}
}

func (r *AlbListenerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan AlbListenerResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createRulesRequest := albv20251030.NewCreateRulesRequest()
	createRulesRequest.ListenerId = plan.ListenerId.ValueStringPointer()
	createRulesRequest.LoadBalancerId = plan.LoadBalancerId.ValueStringPointer()
	createRulesRequest.Rules = expandAlbRules(ctx, plan.Rules, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	service := NewAlbService(r.Client())
	ruleIds, err := service.CreateRules(ctx, createRulesRequest)
	if err != nil {
		resp.Diagnostics.AddError("Error creating ALB listener rules", fmt.Sprintf("CreateRules failed: %v", err))
		return
	}

	if len(ruleIds) == 0 || ruleIds[0] == "" {
		resp.Diagnostics.AddError("Error creating ALB listener rules", "CreateRules returned an empty rule id.")
		return
	}

	ruleId := ruleIds[0]
	plan.Id = types.StringValue(ruleId)
	normalizeAlbListenerComputedUnknowns(ctx, &plan, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	readBackAlbListener(ctx, &plan, service, &resp.Diagnostics)
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

	service := NewAlbService(r.Client())
	item := describeAlbListenerRule(ctx, &state, service, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	if item == nil {
		resp.Diagnostics.AddWarning(
			"Resource not found",
			fmt.Sprintf("ALB listener rule %s not found, removing from state.", state.Id.ValueString()),
		)
		resp.State.RemoveResource(ctx)
		return
	}

	state.Rules = flattenAlbRules(ctx, item, &resp.Diagnostics)
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

	ruleId := state.Id.ValueString()
	service := NewAlbService(r.Client())

	if !plan.Rules.Equal(state.Rules) {
		modifyRulesAttributesRequest := albv20251030.NewModifyRulesAttributesRequest()
		modifyRulesAttributesRequest.ListenerId = plan.ListenerId.ValueStringPointer()
		modifyRulesAttributesRequest.LoadBalancerId = plan.LoadBalancerId.ValueStringPointer()
		modifyRulesAttributesRequest.Rules = expandAlbRuleModifies(ctx, plan.Rules, ruleId, &resp.Diagnostics)
		if resp.Diagnostics.HasError() {
			return
		}

		if err := service.ModifyRulesAttributes(ctx, modifyRulesAttributesRequest); err != nil {
			resp.Diagnostics.AddError("Error updating ALB listener rules", fmt.Sprintf("ModifyRulesAttributes failed: %v", err))
			return
		}
	}

	item := describeAlbListenerRule(ctx, &plan, service, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	if item == nil {
		resp.Diagnostics.AddError("ALB listener rule not readable after update", fmt.Sprintf("rule %s cannot be read back", ruleId))
		return
	}

	plan.Id = state.Id
	plan.Rules = flattenAlbRules(ctx, item, &resp.Diagnostics)
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

	ruleId := state.Id.ValueString()
	deleteRulesRequest := albv20251030.NewDeleteRulesRequest()
	deleteRulesRequest.ListenerId = state.ListenerId.ValueStringPointer()
	deleteRulesRequest.LoadBalancerId = state.LoadBalancerId.ValueStringPointer()
	deleteRulesRequest.RuleIds = []*string{fwHelper.String(ruleId)}

	service := NewAlbService(r.Client())
	if err := service.DeleteRules(ctx, deleteRulesRequest); err != nil {
		resp.Diagnostics.AddError("Error deleting ALB listener rules", fmt.Sprintf("DeleteRules failed: %v", err))
		return
	}
}

// describeAlbListenerRule queries DescribeRules for the rule in state/plan and
// returns the matching item, or nil when no item matches.
func describeAlbListenerRule(ctx context.Context, m *AlbListenerResourceModel, service AlbService, diags *diag.Diagnostics) *albv20251030.RuleOutput {
	describeRulesRequest := albv20251030.NewDescribeRulesRequest()
	describeRulesRequest.ListenerId = m.ListenerId.ValueStringPointer()
	describeRulesRequest.LoadBalancerId = m.LoadBalancerId.ValueStringPointer()
	describeRulesRequest.RuleIds = []*string{fwHelper.String(m.Id.ValueString())}

	rules, err := service.DescribeRules(ctx, describeRulesRequest)
	if err != nil {
		diags.AddError("Error reading ALB listener rules", err.Error())
		return nil
	}

	return findAlbRule(rules, m.Id.ValueString())
}

// readBackAlbListener refreshes the whole rule list from the cloud.
func readBackAlbListener(ctx context.Context, m *AlbListenerResourceModel, service AlbService, diags *diag.Diagnostics) {
	item := describeAlbListenerRule(ctx, m, service, diags)
	if diags.HasError() {
		return
	}

	if item == nil {
		diags.AddError("ALB listener rule not readable after create", fmt.Sprintf("rule %s cannot be read back", m.Id.ValueString()))
		return
	}

	m.Rules = flattenAlbRules(ctx, item, diags)
}

// findAlbRule picks the rule with the given ID from a DescribeRules list.
func findAlbRule(rules []*albv20251030.RuleOutput, ruleId string) *albv20251030.RuleOutput {
	for _, rule := range rules {
		if rule != nil && rule.RuleId != nil && *rule.RuleId == ruleId {
			return rule
		}
	}

	return nil
}

// normalizeAlbListenerComputedUnknowns rewrites unknown computed attributes to
// null so the model can be persisted into state before the read-back.
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
		if rules[i].RuleId.IsUnknown() {
			rules[i].RuleId = types.StringNull()
		}

		if rules[i].CreateTime.IsUnknown() {
			rules[i].CreateTime = types.StringNull()
		}

		if rules[i].ModifyTime.IsUnknown() {
			rules[i].ModifyTime = types.StringNull()
		}

		if rules[i].Status.IsUnknown() {
			rules[i].Status = types.StringNull()
		}
	}

	m.Rules = flattenAlbRuleModels(ctx, rules, diags)
}
