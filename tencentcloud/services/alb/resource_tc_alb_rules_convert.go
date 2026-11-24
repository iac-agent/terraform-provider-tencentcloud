package alb

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	albv20251030 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/alb/v20251030"

	fwhelper "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

// The object type constructors below feed the lists rebuilt during refresh;
// they must mirror the NestedObject attribute sets of the resource schema.
func albTagObjectType() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"tag_key":   types.StringType,
			"tag_value": types.StringType,
		},
	}
}

func albRuleObjectType() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"actions":    types.ListType{ElemType: albRuleActionObjectType()},
			"conditions": types.ListType{ElemType: albRuleConditionObjectType()},
			"priority":   types.Int64Type,
			"direction":  types.StringType,
			"rule_name":  types.StringType,
			"tags":       types.ListType{ElemType: albTagObjectType()},
		},
	}
}

func albRuleActionObjectType() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"order":                 types.Int64Type,
			"type":                  types.StringType,
			"fixed_response_config": types.ListType{ElemType: albFixedResponseConfigObjectType()},
			"insert_header_config":  types.ListType{ElemType: albInsertHeaderConfigObjectType()},
			"redirect_config":       types.ListType{ElemType: albRedirectConfigObjectType()},
			"remove_header_config":  types.ListType{ElemType: albRemoveHeaderConfigObjectType()},
			"rewrite_config":        types.ListType{ElemType: albRewriteConfigObjectType()},
			"target_group_config":   types.ListType{ElemType: albTargetGroupConfigObjectType()},
		},
	}
}

func albFixedResponseConfigObjectType() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"http_code":    types.Int64Type,
			"content":      types.StringType,
			"content_type": types.StringType,
		},
	}
}

func albInsertHeaderConfigObjectType() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"key":        types.StringType,
			"value":      types.StringType,
			"value_type": types.StringType,
		},
	}
}

func albRedirectConfigObjectType() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"http_code": types.Int64Type,
			"host":      types.StringType,
			"path":      types.StringType,
			"port":      types.StringType,
			"protocol":  types.StringType,
			"query":     types.StringType,
		},
	}
}

func albRemoveHeaderConfigObjectType() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"key": types.StringType,
		},
	}
}

func albRewriteConfigObjectType() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"host":  types.StringType,
			"path":  types.StringType,
			"query": types.StringType,
		},
	}
}

func albTargetGroupConfigObjectType() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"target_groups":               types.ListType{ElemType: albTargetGroupTupleObjectType()},
			"target_group_sticky_session": types.ListType{ElemType: albTargetGroupStickySessionObjectType()},
		},
	}
}

func albTargetGroupTupleObjectType() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"target_group_id": types.StringType,
			"weight":          types.Int64Type,
		},
	}
}

func albTargetGroupStickySessionObjectType() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"enabled": types.BoolType,
			"timeout": types.Int64Type,
		},
	}
}

func albRuleConditionObjectType() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"type":                types.StringType,
			"cookie_config":       types.ListType{ElemType: albCookieConfigObjectType()},
			"header_config":       types.ListType{ElemType: albHeaderConfigObjectType()},
			"host_config":         types.ListType{ElemType: types.StringType},
			"method_config":       types.ListType{ElemType: types.StringType},
			"path_config":         types.ListType{ElemType: types.StringType},
			"query_string_config": types.ListType{ElemType: albQueryStringConfigObjectType()},
			"source_ip_config":    types.ListType{ElemType: types.StringType},
		},
	}
}

func albCookieConfigObjectType() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"key":   types.StringType,
			"value": types.StringType,
		},
	}
}

func albHeaderConfigObjectType() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"key":    types.StringType,
			"values": types.ListType{ElemType: types.StringType},
		},
	}
}

func albQueryStringConfigObjectType() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"key":   types.StringType,
			"value": types.StringType,
		},
	}
}

// expandAlbRules converts the planned rules list into the SDK rule input
// slice. Exactly one rule is expected; anything else is reported as an error.
func expandAlbRules(ctx context.Context, list types.List) ([]*albv20251030.RuleInput, diag.Diagnostics) {
	var diags diag.Diagnostics
	if list.IsNull() || list.IsUnknown() {
		diags.AddError(
			"Invalid ALB forwarding rules configuration",
			"`rules` is required and must contain exactly one rule.",
		)
		return nil, diags
	}

	var rules []albRuleModel
	diags.Append(list.ElementsAs(ctx, &rules, false)...)
	if diags.HasError() {
		return nil, diags
	}

	if len(rules) != 1 {
		diags.AddError(
			"Invalid ALB forwarding rules configuration",
			fmt.Sprintf("Exactly one rule is expected in `rules`, got %d.", len(rules)),
		)
		return nil, diags
	}

	input, ruleDiags := expandAlbRuleInput(ctx, &rules[0])
	diags.Append(ruleDiags...)
	if diags.HasError() {
		return nil, diags
	}

	return []*albv20251030.RuleInput{input}, diags
}

// expandAlbRuleInput converts one planned rule into the SDK create input.
func expandAlbRuleInput(ctx context.Context, rule *albRuleModel) (*albv20251030.RuleInput, diag.Diagnostics) {
	var diags diag.Diagnostics

	actions, actionDiags := expandAlbRuleActions(ctx, rule.Actions)
	diags.Append(actionDiags...)
	if diags.HasError() {
		return nil, diags
	}

	conditions, conditionDiags := expandAlbRuleConditions(ctx, rule.Conditions)
	diags.Append(conditionDiags...)
	if diags.HasError() {
		return nil, diags
	}

	tags, tagDiags := expandAlbTags(ctx, rule.Tags)
	diags.Append(tagDiags...)
	if diags.HasError() {
		return nil, diags
	}

	return &albv20251030.RuleInput{
		Actions:    actions,
		Conditions: conditions,
		Priority:   fwhelper.Int64PointerFromValue(rule.Priority),
		Direction:  fwhelper.StringPointerFromValue(rule.Direction),
		RuleName:   fwhelper.StringPointerFromValue(rule.RuleName),
		Tags:       tags,
	}, diags
}

// expandAlbRuleActions converts the planned actions list into the SDK slice.
func expandAlbRuleActions(ctx context.Context, list types.List) ([]*albv20251030.RuleAction, diag.Diagnostics) {
	var diags diag.Diagnostics
	if list.IsNull() || list.IsUnknown() {
		diags.AddError(
			"Invalid ALB forwarding rule action configuration",
			"`actions` is required inside every rule.",
		)
		return nil, diags
	}

	var actions []albRuleActionModel
	diags.Append(list.ElementsAs(ctx, &actions, false)...)
	if diags.HasError() {
		return nil, diags
	}

	out := make([]*albv20251030.RuleAction, 0, len(actions))
	for i, action := range actions {
		out = append(out, &albv20251030.RuleAction{
			Order: fwhelper.Int64PointerFromValue(action.Order),
			Type:  fwhelper.StringPointerFromValue(action.Type),
		})

		fixedResponse, fixedDiags := expandAlbFixedResponseConfig(ctx, action.FixedResponseConfig)
		diags.Append(fixedDiags...)
		if diags.HasError() {
			return nil, diags
		}
		out[i].FixedResponseConfig = fixedResponse

		insertHeader, insertDiags := expandAlbInsertHeaderConfig(ctx, action.InsertHeaderConfig)
		diags.Append(insertDiags...)
		if diags.HasError() {
			return nil, diags
		}
		out[i].InsertHeaderConfig = insertHeader

		redirect, redirectDiags := expandAlbRedirectConfig(ctx, action.RedirectConfig)
		diags.Append(redirectDiags...)
		if diags.HasError() {
			return nil, diags
		}
		out[i].RedirectConfig = redirect

		removeHeader, removeDiags := expandAlbRemoveHeaderConfig(ctx, action.RemoveHeaderConfig)
		diags.Append(removeDiags...)
		if diags.HasError() {
			return nil, diags
		}
		out[i].RemoveHeaderConfig = removeHeader

		rewrite, rewriteDiags := expandAlbRewriteConfig(ctx, action.RewriteConfig)
		diags.Append(rewriteDiags...)
		if diags.HasError() {
			return nil, diags
		}
		out[i].RewriteConfig = rewrite

		targetGroup, targetDiags := expandAlbTargetGroupConfig(ctx, action.TargetGroupConfig)
		diags.Append(targetDiags...)
		if diags.HasError() {
			return nil, diags
		}
		out[i].TargetGroupConfig = targetGroup
	}

	return out, diags
}

// expandAlbFixedResponseConfig converts the planned fixed_response_config into
// the SDK structure.
func expandAlbFixedResponseConfig(ctx context.Context, list types.List) (*albv20251030.FixedResponseInfo, diag.Diagnostics) {
	var diags diag.Diagnostics
	if list.IsNull() || list.IsUnknown() {
		return nil, diags
	}

	var configs []albFixedResponseConfigModel
	diags.Append(list.ElementsAs(ctx, &configs, false)...)
	if diags.HasError() {
		return nil, diags
	}
	if len(configs) == 0 {
		return nil, diags
	}

	return &albv20251030.FixedResponseInfo{
		HttpCode:    fwhelper.Int64PointerFromValue(configs[0].HttpCode),
		Content:     fwhelper.StringPointerFromValue(configs[0].Content),
		ContentType: fwhelper.StringPointerFromValue(configs[0].ContentType),
	}, diags
}

// expandAlbInsertHeaderConfig converts the planned insert_header_config into
// the SDK structure.
func expandAlbInsertHeaderConfig(ctx context.Context, list types.List) (*albv20251030.InsertHTTPHeaderInfo, diag.Diagnostics) {
	var diags diag.Diagnostics
	if list.IsNull() || list.IsUnknown() {
		return nil, diags
	}

	var configs []albInsertHeaderConfigModel
	diags.Append(list.ElementsAs(ctx, &configs, false)...)
	if diags.HasError() {
		return nil, diags
	}
	if len(configs) == 0 {
		return nil, diags
	}

	return &albv20251030.InsertHTTPHeaderInfo{
		Key:       fwhelper.StringPointerFromValue(configs[0].Key),
		Value:     fwhelper.StringPointerFromValue(configs[0].Value),
		ValueType: fwhelper.StringPointerFromValue(configs[0].ValueType),
	}, diags
}

// expandAlbRedirectConfig converts the planned redirect_config into the SDK
// structure.
func expandAlbRedirectConfig(ctx context.Context, list types.List) (*albv20251030.HTTPRedirectInfo, diag.Diagnostics) {
	var diags diag.Diagnostics
	if list.IsNull() || list.IsUnknown() {
		return nil, diags
	}

	var configs []albRedirectConfigModel
	diags.Append(list.ElementsAs(ctx, &configs, false)...)
	if diags.HasError() {
		return nil, diags
	}
	if len(configs) == 0 {
		return nil, diags
	}

	return &albv20251030.HTTPRedirectInfo{
		HttpCode: fwhelper.Int64PointerFromValue(configs[0].HttpCode),
		Host:     fwhelper.StringPointerFromValue(configs[0].Host),
		Path:     fwhelper.StringPointerFromValue(configs[0].Path),
		Port:     fwhelper.StringPointerFromValue(configs[0].Port),
		Protocol: fwhelper.StringPointerFromValue(configs[0].Protocol),
		Query:    fwhelper.StringPointerFromValue(configs[0].Query),
	}, diags
}

// expandAlbRemoveHeaderConfig converts the planned remove_header_config into
// the SDK structure.
func expandAlbRemoveHeaderConfig(ctx context.Context, list types.List) (*albv20251030.RemoveHTTPHeaderInfo, diag.Diagnostics) {
	var diags diag.Diagnostics
	if list.IsNull() || list.IsUnknown() {
		return nil, diags
	}

	var configs []albRemoveHeaderConfigModel
	diags.Append(list.ElementsAs(ctx, &configs, false)...)
	if diags.HasError() {
		return nil, diags
	}
	if len(configs) == 0 {
		return nil, diags
	}

	return &albv20251030.RemoveHTTPHeaderInfo{
		Key: fwhelper.StringPointerFromValue(configs[0].Key),
	}, diags
}

// expandAlbRewriteConfig converts the planned rewrite_config into the SDK
// structure.
func expandAlbRewriteConfig(ctx context.Context, list types.List) (*albv20251030.HTTPRewriteInfo, diag.Diagnostics) {
	var diags diag.Diagnostics
	if list.IsNull() || list.IsUnknown() {
		return nil, diags
	}

	var configs []albRewriteConfigModel
	diags.Append(list.ElementsAs(ctx, &configs, false)...)
	if diags.HasError() {
		return nil, diags
	}
	if len(configs) == 0 {
		return nil, diags
	}

	return &albv20251030.HTTPRewriteInfo{
		Host:  fwhelper.StringPointerFromValue(configs[0].Host),
		Path:  fwhelper.StringPointerFromValue(configs[0].Path),
		Query: fwhelper.StringPointerFromValue(configs[0].Query),
	}, diags
}

// expandAlbTargetGroupConfig converts the planned target_group_config into the
// SDK structure.
func expandAlbTargetGroupConfig(ctx context.Context, list types.List) (*albv20251030.TargetGroupConfig, diag.Diagnostics) {
	var diags diag.Diagnostics
	if list.IsNull() || list.IsUnknown() {
		return nil, diags
	}

	var configs []albTargetGroupConfigModel
	diags.Append(list.ElementsAs(ctx, &configs, false)...)
	if diags.HasError() {
		return nil, diags
	}
	if len(configs) == 0 {
		return nil, diags
	}

	config := configs[0]

	targetGroups, groupsDiags := expandAlbTargetGroups(ctx, config.TargetGroups)
	diags.Append(groupsDiags...)
	if diags.HasError() {
		return nil, diags
	}

	stickySession, stickyDiags := expandAlbTargetGroupStickySession(ctx, config.TargetGroupStickySession)
	diags.Append(stickyDiags...)
	if diags.HasError() {
		return nil, diags
	}

	return &albv20251030.TargetGroupConfig{
		TargetGroups:             targetGroups,
		TargetGroupStickySession: stickySession,
	}, diags
}

// expandAlbTargetGroups converts the planned target_groups list into the SDK
// slice.
func expandAlbTargetGroups(ctx context.Context, list types.List) ([]*albv20251030.TargetGroupTuple, diag.Diagnostics) {
	var diags diag.Diagnostics
	if list.IsNull() || list.IsUnknown() {
		return nil, diags
	}

	var groups []albTargetGroupTupleModel
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

// expandAlbTargetGroupStickySession converts the planned
// target_group_sticky_session into the SDK structure.
func expandAlbTargetGroupStickySession(ctx context.Context, list types.List) (*albv20251030.TargetGroupStickySession, diag.Diagnostics) {
	var diags diag.Diagnostics
	if list.IsNull() || list.IsUnknown() {
		return nil, diags
	}

	var sessions []albTargetGroupStickySessionModel
	diags.Append(list.ElementsAs(ctx, &sessions, false)...)
	if diags.HasError() {
		return nil, diags
	}
	if len(sessions) == 0 {
		return nil, diags
	}

	return &albv20251030.TargetGroupStickySession{
		Enabled: fwhelper.BoolPointerFromValue(sessions[0].Enabled),
		Timeout: fwhelper.Int64PointerFromValue(sessions[0].Timeout),
	}, diags
}

// expandAlbRuleConditions converts the planned conditions list into the SDK
// slice.
func expandAlbRuleConditions(ctx context.Context, list types.List) ([]*albv20251030.RuleCondition, diag.Diagnostics) {
	var diags diag.Diagnostics
	if list.IsNull() || list.IsUnknown() {
		diags.AddError(
			"Invalid ALB forwarding rule condition configuration",
			"`conditions` is required inside every rule.",
		)
		return nil, diags
	}

	var conditions []albRuleConditionModel
	diags.Append(list.ElementsAs(ctx, &conditions, false)...)
	if diags.HasError() {
		return nil, diags
	}

	out := make([]*albv20251030.RuleCondition, 0, len(conditions))
	for i, condition := range conditions {
		out = append(out, &albv20251030.RuleCondition{
			Type: fwhelper.StringPointerFromValue(condition.Type),
		})

		cookies, cookieDiags := expandAlbCookieConfig(ctx, condition.CookieConfig)
		diags.Append(cookieDiags...)
		if diags.HasError() {
			return nil, diags
		}
		out[i].CookieConfig = cookies

		header, headerDiags := expandAlbHeaderConfig(ctx, condition.HeaderConfig)
		diags.Append(headerDiags...)
		if diags.HasError() {
			return nil, diags
		}
		out[i].HeaderConfig = header

		hosts, hostDiags := expandStringList(ctx, condition.HostConfig)
		diags.Append(hostDiags...)
		if diags.HasError() {
			return nil, diags
		}
		out[i].HostConfig = hosts

		methods, methodDiags := expandStringList(ctx, condition.MethodConfig)
		diags.Append(methodDiags...)
		if diags.HasError() {
			return nil, diags
		}
		out[i].MethodConfig = methods

		paths, pathDiags := expandStringList(ctx, condition.PathConfig)
		diags.Append(pathDiags...)
		if diags.HasError() {
			return nil, diags
		}
		out[i].PathConfig = paths

		queryStrings, queryDiags := expandAlbQueryStringConfig(ctx, condition.QueryStringConfig)
		diags.Append(queryDiags...)
		if diags.HasError() {
			return nil, diags
		}
		out[i].QueryStringConfig = queryStrings

		sourceIps, sourceDiags := expandStringList(ctx, condition.SourceIpConfig)
		diags.Append(sourceDiags...)
		if diags.HasError() {
			return nil, diags
		}
		out[i].SourceIpConfig = sourceIps
	}

	return out, diags
}

// expandAlbCookieConfig converts the planned cookie_config into the SDK slice.
func expandAlbCookieConfig(ctx context.Context, list types.List) ([]*albv20251030.HTTPCookieInfo, diag.Diagnostics) {
	var diags diag.Diagnostics
	if list.IsNull() || list.IsUnknown() {
		return nil, diags
	}

	var cookies []albCookieConfigModel
	diags.Append(list.ElementsAs(ctx, &cookies, false)...)
	if diags.HasError() {
		return nil, diags
	}

	out := make([]*albv20251030.HTTPCookieInfo, 0, len(cookies))
	for _, cookie := range cookies {
		out = append(out, &albv20251030.HTTPCookieInfo{
			Key:   fwhelper.StringPointerFromValue(cookie.Key),
			Value: fwhelper.StringPointerFromValue(cookie.Value),
		})
	}

	return out, diags
}

// expandAlbHeaderConfig converts the planned header_config into the SDK
// structure.
func expandAlbHeaderConfig(ctx context.Context, list types.List) (*albv20251030.HTTPHeaderInfo, diag.Diagnostics) {
	var diags diag.Diagnostics
	if list.IsNull() || list.IsUnknown() {
		return nil, diags
	}

	var headers []albHeaderConfigModel
	diags.Append(list.ElementsAs(ctx, &headers, false)...)
	if diags.HasError() {
		return nil, diags
	}
	if len(headers) == 0 {
		return nil, diags
	}

	values, valuesDiags := expandStringList(ctx, headers[0].Values)
	diags.Append(valuesDiags...)
	if diags.HasError() {
		return nil, diags
	}

	return &albv20251030.HTTPHeaderInfo{
		Key:    fwhelper.StringPointerFromValue(headers[0].Key),
		Values: values,
	}, diags
}

// expandAlbQueryStringConfig converts the planned query_string_config into the
// SDK slice.
func expandAlbQueryStringConfig(ctx context.Context, list types.List) ([]*albv20251030.HTTPQueryStringInfo, diag.Diagnostics) {
	var diags diag.Diagnostics
	if list.IsNull() || list.IsUnknown() {
		return nil, diags
	}

	var queryStrings []albQueryStringConfigModel
	diags.Append(list.ElementsAs(ctx, &queryStrings, false)...)
	if diags.HasError() {
		return nil, diags
	}

	out := make([]*albv20251030.HTTPQueryStringInfo, 0, len(queryStrings))
	for _, queryString := range queryStrings {
		out = append(out, &albv20251030.HTTPQueryStringInfo{
			Key:   fwhelper.StringPointerFromValue(queryString.Key),
			Value: fwhelper.StringPointerFromValue(queryString.Value),
		})
	}

	return out, diags
}

// expandAlbTags converts the planned tags list into the SDK slice.
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

// expandAlbRuleModifyRequest converts the planned rule plus its rule id into
// the SDK modify input. The modify API accepts the full actions and conditions
// lists, so they are re-expanded from the plan.
func expandAlbRuleModifyRequest(ctx context.Context, rule *albRuleModel, ruleId string) (*albv20251030.RuleModify, diag.Diagnostics) {
	var diags diag.Diagnostics

	actions, actionDiags := expandAlbRuleActions(ctx, rule.Actions)
	diags.Append(actionDiags...)
	if diags.HasError() {
		return nil, diags
	}

	conditions, conditionDiags := expandAlbRuleConditions(ctx, rule.Conditions)
	diags.Append(conditionDiags...)
	if diags.HasError() {
		return nil, diags
	}

	return &albv20251030.RuleModify{
		Actions:    actions,
		Conditions: conditions,
		Priority:   fwhelper.Int64PointerFromValue(rule.Priority),
		RuleId:     helper.String(ruleId),
		RuleName:   fwhelper.StringPointerFromValue(rule.RuleName),
	}, diags
}

// refreshAlbRulesModel copies the cloud state into the model. Attributes the
// Describe API may omit keep their plan/state value, so a value configured by
// the practitioner is never nulled out after apply.
func refreshAlbRulesModel(ctx context.Context, m *AlbRulesResourceModel, loadBalancerId, listenerId string, item *albv20251030.RuleOutput, prevTags types.List) diag.Diagnostics {
	var diags diag.Diagnostics

	m.LoadBalancerId = types.StringValue(loadBalancerId)
	m.ListenerId = types.StringValue(listenerId)
	m.Id = types.StringValue(albRulesId(loadBalancerId, listenerId, stringValueOf(item.RuleId)))

	rules, ruleDiags := flattenAlbRules(ctx, []*albv20251030.RuleOutput{item}, prevTags)
	diags.Append(ruleDiags...)
	m.Rules = rules

	return diags
}

// flattenAlbRules builds the rules list from the Describe response. The
// previous tags list is carried over because DescribeRules does not echo the
// rule tags.
func flattenAlbRules(ctx context.Context, rules []*albv20251030.RuleOutput, prevTags types.List) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics
	objType := albRuleObjectType()

	out := make([]albRuleModel, 0, len(rules))
	for _, rule := range rules {
		if rule == nil {
			continue
		}

		actions, actionDiags := flattenAlbRuleActions(ctx, rule.Actions)
		diags.Append(actionDiags...)
		if diags.HasError() {
			return types.ListNull(objType), diags
		}

		conditions, conditionDiags := flattenAlbRuleConditions(ctx, rule.Conditions)
		diags.Append(conditionDiags...)
		if diags.HasError() {
			return types.ListNull(objType), diags
		}

		out = append(out, albRuleModel{
			Actions:    actions,
			Conditions: conditions,
			Priority:   int64PtrValueOrPrev(rule.Priority),
			Direction:  fwhelper.StringValueOrPrev(rule.Direction, types.StringNull()),
			RuleName:   fwhelper.StringValueOrPrev(rule.RuleName, types.StringNull()),
			Tags:       prevTags,
		})
	}

	list, listDiags := types.ListValueFrom(ctx, objType, out)
	diags.Append(listDiags...)
	if diags.HasError() {
		return types.ListNull(objType), diags
	}
	return list, diags
}

// flattenAlbRuleActions builds the actions list from the Describe response.
func flattenAlbRuleActions(ctx context.Context, actions []*albv20251030.RuleAction) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics
	objType := albRuleActionObjectType()

	out := make([]albRuleActionModel, 0, len(actions))
	for _, action := range actions {
		if action == nil {
			continue
		}

		fixedResponse, fixedDiags := flattenAlbFixedResponseConfig(ctx, action.FixedResponseConfig)
		diags.Append(fixedDiags...)
		if diags.HasError() {
			return types.ListNull(objType), diags
		}

		insertHeader, insertDiags := flattenAlbInsertHeaderConfig(ctx, action.InsertHeaderConfig)
		diags.Append(insertDiags...)
		if diags.HasError() {
			return types.ListNull(objType), diags
		}

		redirect, redirectDiags := flattenAlbRedirectConfig(ctx, action.RedirectConfig)
		diags.Append(redirectDiags...)
		if diags.HasError() {
			return types.ListNull(objType), diags
		}

		removeHeader, removeDiags := flattenAlbRemoveHeaderConfig(ctx, action.RemoveHeaderConfig)
		diags.Append(removeDiags...)
		if diags.HasError() {
			return types.ListNull(objType), diags
		}

		rewrite, rewriteDiags := flattenAlbRewriteConfig(ctx, action.RewriteConfig)
		diags.Append(rewriteDiags...)
		if diags.HasError() {
			return types.ListNull(objType), diags
		}

		targetGroup, targetDiags := flattenAlbTargetGroupConfig(ctx, action.TargetGroupConfig)
		diags.Append(targetDiags...)
		if diags.HasError() {
			return types.ListNull(objType), diags
		}

		out = append(out, albRuleActionModel{
			Order:               int64PtrValueOrPrev(action.Order),
			Type:                fwhelper.StringValueOrPrev(action.Type, types.StringNull()),
			FixedResponseConfig: fixedResponse,
			InsertHeaderConfig:  insertHeader,
			RedirectConfig:      redirect,
			RemoveHeaderConfig:  removeHeader,
			RewriteConfig:       rewrite,
			TargetGroupConfig:   targetGroup,
		})
	}

	list, listDiags := types.ListValueFrom(ctx, objType, out)
	diags.Append(listDiags...)
	if diags.HasError() {
		return types.ListNull(objType), diags
	}
	return list, diags
}

// flattenAlbFixedResponseConfig builds the fixed_response_config list.
func flattenAlbFixedResponseConfig(ctx context.Context, config *albv20251030.FixedResponseInfo) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics
	objType := albFixedResponseConfigObjectType()
	if config == nil {
		return types.ListNull(objType), diags
	}

	out := []albFixedResponseConfigModel{{
		HttpCode:    int64PtrValueOrPrev(config.HttpCode),
		Content:     fwhelper.StringPointerValueOrNull(config.Content),
		ContentType: fwhelper.StringPointerValueOrNull(config.ContentType),
	}}

	list, listDiags := types.ListValueFrom(ctx, objType, out)
	diags.Append(listDiags...)
	if diags.HasError() {
		return types.ListNull(objType), diags
	}
	return list, diags
}

// flattenAlbInsertHeaderConfig builds the insert_header_config list.
func flattenAlbInsertHeaderConfig(ctx context.Context, config *albv20251030.InsertHTTPHeaderInfo) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics
	objType := albInsertHeaderConfigObjectType()
	if config == nil {
		return types.ListNull(objType), diags
	}

	out := []albInsertHeaderConfigModel{{
		Key:       fwhelper.StringPointerValueOrNull(config.Key),
		Value:     fwhelper.StringPointerValueOrNull(config.Value),
		ValueType: fwhelper.StringPointerValueOrNull(config.ValueType),
	}}

	list, listDiags := types.ListValueFrom(ctx, objType, out)
	diags.Append(listDiags...)
	if diags.HasError() {
		return types.ListNull(objType), diags
	}
	return list, diags
}

// flattenAlbRedirectConfig builds the redirect_config list.
func flattenAlbRedirectConfig(ctx context.Context, config *albv20251030.HTTPRedirectInfo) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics
	objType := albRedirectConfigObjectType()
	if config == nil {
		return types.ListNull(objType), diags
	}

	out := []albRedirectConfigModel{{
		HttpCode: int64PtrValueOrPrev(config.HttpCode),
		Host:     fwhelper.StringPointerValueOrNull(config.Host),
		Path:     fwhelper.StringPointerValueOrNull(config.Path),
		Port:     fwhelper.StringPointerValueOrNull(config.Port),
		Protocol: fwhelper.StringPointerValueOrNull(config.Protocol),
		Query:    fwhelper.StringPointerValueOrNull(config.Query),
	}}

	list, listDiags := types.ListValueFrom(ctx, objType, out)
	diags.Append(listDiags...)
	if diags.HasError() {
		return types.ListNull(objType), diags
	}
	return list, diags
}

// flattenAlbRemoveHeaderConfig builds the remove_header_config list.
func flattenAlbRemoveHeaderConfig(ctx context.Context, config *albv20251030.RemoveHTTPHeaderInfo) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics
	objType := albRemoveHeaderConfigObjectType()
	if config == nil {
		return types.ListNull(objType), diags
	}

	out := []albRemoveHeaderConfigModel{{
		Key: fwhelper.StringPointerValueOrNull(config.Key),
	}}

	list, listDiags := types.ListValueFrom(ctx, objType, out)
	diags.Append(listDiags...)
	if diags.HasError() {
		return types.ListNull(objType), diags
	}
	return list, diags
}

// flattenAlbRewriteConfig builds the rewrite_config list.
func flattenAlbRewriteConfig(ctx context.Context, config *albv20251030.HTTPRewriteInfo) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics
	objType := albRewriteConfigObjectType()
	if config == nil {
		return types.ListNull(objType), diags
	}

	out := []albRewriteConfigModel{{
		Host:  fwhelper.StringPointerValueOrNull(config.Host),
		Path:  fwhelper.StringPointerValueOrNull(config.Path),
		Query: fwhelper.StringPointerValueOrNull(config.Query),
	}}

	list, listDiags := types.ListValueFrom(ctx, objType, out)
	diags.Append(listDiags...)
	if diags.HasError() {
		return types.ListNull(objType), diags
	}
	return list, diags
}

// flattenAlbTargetGroupConfig builds the target_group_config list.
func flattenAlbTargetGroupConfig(ctx context.Context, config *albv20251030.TargetGroupConfig) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics
	objType := albTargetGroupConfigObjectType()
	if config == nil {
		return types.ListNull(objType), diags
	}

	targetGroups, groupsDiags := flattenAlbTargetGroups(ctx, config.TargetGroups)
	diags.Append(groupsDiags...)
	if diags.HasError() {
		return types.ListNull(objType), diags
	}

	stickySession, stickyDiags := flattenAlbTargetGroupStickySession(ctx, config.TargetGroupStickySession)
	diags.Append(stickyDiags...)
	if diags.HasError() {
		return types.ListNull(objType), diags
	}

	out := []albTargetGroupConfigModel{{
		TargetGroups:             targetGroups,
		TargetGroupStickySession: stickySession,
	}}

	list, listDiags := types.ListValueFrom(ctx, objType, out)
	diags.Append(listDiags...)
	if diags.HasError() {
		return types.ListNull(objType), diags
	}
	return list, diags
}

// flattenAlbTargetGroups builds the target_groups list.
func flattenAlbTargetGroups(ctx context.Context, groups []*albv20251030.TargetGroupTuple) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics
	objType := albTargetGroupTupleObjectType()

	out := make([]albTargetGroupTupleModel, 0, len(groups))
	for _, group := range groups {
		if group == nil {
			continue
		}
		out = append(out, albTargetGroupTupleModel{
			TargetGroupId: fwhelper.StringPointerValueOrNull(group.TargetGroupId),
			Weight:        int64PtrValueOrPrev(group.Weight),
		})
	}

	list, listDiags := types.ListValueFrom(ctx, objType, out)
	diags.Append(listDiags...)
	if diags.HasError() {
		return types.ListNull(objType), diags
	}
	return list, diags
}

// flattenAlbTargetGroupStickySession builds the target_group_sticky_session
// list.
func flattenAlbTargetGroupStickySession(ctx context.Context, session *albv20251030.TargetGroupStickySession) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics
	objType := albTargetGroupStickySessionObjectType()
	if session == nil {
		return types.ListNull(objType), diags
	}

	out := []albTargetGroupStickySessionModel{{
		Enabled: fwhelper.BoolValueOrNull(session.Enabled),
		Timeout: int64PtrValueOrPrev(session.Timeout),
	}}

	list, listDiags := types.ListValueFrom(ctx, objType, out)
	diags.Append(listDiags...)
	if diags.HasError() {
		return types.ListNull(objType), diags
	}
	return list, diags
}

// flattenAlbRuleConditions builds the conditions list from the Describe
// response.
func flattenAlbRuleConditions(ctx context.Context, conditions []*albv20251030.RuleCondition) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics
	objType := albRuleConditionObjectType()

	out := make([]albRuleConditionModel, 0, len(conditions))
	for _, condition := range conditions {
		if condition == nil {
			continue
		}

		cookies, cookieDiags := flattenAlbCookieConfig(ctx, condition.CookieConfig)
		diags.Append(cookieDiags...)
		if diags.HasError() {
			return types.ListNull(objType), diags
		}

		header, headerDiags := flattenAlbHeaderConfig(ctx, condition.HeaderConfig)
		diags.Append(headerDiags...)
		if diags.HasError() {
			return types.ListNull(objType), diags
		}

		hosts, hostDiags := flattenStringList(ctx, condition.HostConfig)
		diags.Append(hostDiags...)
		if diags.HasError() {
			return types.ListNull(objType), diags
		}

		methods, methodDiags := flattenStringList(ctx, condition.MethodConfig)
		diags.Append(methodDiags...)
		if diags.HasError() {
			return types.ListNull(objType), diags
		}

		paths, pathDiags := flattenStringList(ctx, condition.PathConfig)
		diags.Append(pathDiags...)
		if diags.HasError() {
			return types.ListNull(objType), diags
		}

		queryStrings, queryDiags := flattenAlbQueryStringConfig(ctx, condition.QueryStringConfig)
		diags.Append(queryDiags...)
		if diags.HasError() {
			return types.ListNull(objType), diags
		}

		sourceIps, sourceDiags := flattenStringList(ctx, condition.SourceIpConfig)
		diags.Append(sourceDiags...)
		if diags.HasError() {
			return types.ListNull(objType), diags
		}

		out = append(out, albRuleConditionModel{
			Type:              fwhelper.StringValueOrPrev(condition.Type, types.StringNull()),
			CookieConfig:      cookies,
			HeaderConfig:      header,
			HostConfig:        hosts,
			MethodConfig:      methods,
			PathConfig:        paths,
			QueryStringConfig: queryStrings,
			SourceIpConfig:    sourceIps,
		})
	}

	list, listDiags := types.ListValueFrom(ctx, objType, out)
	diags.Append(listDiags...)
	if diags.HasError() {
		return types.ListNull(objType), diags
	}
	return list, diags
}

// flattenAlbCookieConfig builds the cookie_config list.
func flattenAlbCookieConfig(ctx context.Context, cookies []*albv20251030.HTTPCookieInfo) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics
	objType := albCookieConfigObjectType()

	out := make([]albCookieConfigModel, 0, len(cookies))
	for _, cookie := range cookies {
		if cookie == nil {
			continue
		}
		out = append(out, albCookieConfigModel{
			Key:   fwhelper.StringPointerValueOrNull(cookie.Key),
			Value: fwhelper.StringPointerValueOrNull(cookie.Value),
		})
	}

	list, listDiags := types.ListValueFrom(ctx, objType, out)
	diags.Append(listDiags...)
	if diags.HasError() {
		return types.ListNull(objType), diags
	}
	return list, diags
}

// flattenAlbHeaderConfig builds the header_config list.
func flattenAlbHeaderConfig(ctx context.Context, header *albv20251030.HTTPHeaderInfo) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics
	objType := albHeaderConfigObjectType()
	if header == nil {
		return types.ListNull(objType), diags
	}

	values, valuesDiags := flattenStringList(ctx, header.Values)
	diags.Append(valuesDiags...)
	if diags.HasError() {
		return types.ListNull(objType), diags
	}

	out := []albHeaderConfigModel{{
		Key:    fwhelper.StringPointerValueOrNull(header.Key),
		Values: values,
	}}

	list, listDiags := types.ListValueFrom(ctx, objType, out)
	diags.Append(listDiags...)
	if diags.HasError() {
		return types.ListNull(objType), diags
	}
	return list, diags
}

// flattenAlbQueryStringConfig builds the query_string_config list.
func flattenAlbQueryStringConfig(ctx context.Context, queryStrings []*albv20251030.HTTPQueryStringInfo) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics
	objType := albQueryStringConfigObjectType()

	out := make([]albQueryStringConfigModel, 0, len(queryStrings))
	for _, queryString := range queryStrings {
		if queryString == nil {
			continue
		}
		out = append(out, albQueryStringConfigModel{
			Key:   fwhelper.StringPointerValueOrNull(queryString.Key),
			Value: fwhelper.StringPointerValueOrNull(queryString.Value),
		})
	}

	list, listDiags := types.ListValueFrom(ctx, objType, out)
	diags.Append(listDiags...)
	if diags.HasError() {
		return types.ListNull(objType), diags
	}
	return list, diags
}

// expandStringList converts a planned list of strings into the SDK slice.
func expandStringList(ctx context.Context, list types.List) ([]*string, diag.Diagnostics) {
	var diags diag.Diagnostics
	if list.IsNull() || list.IsUnknown() {
		return nil, diags
	}

	var values []types.String
	diags.Append(list.ElementsAs(ctx, &values, false)...)
	if diags.HasError() {
		return nil, diags
	}

	out := make([]*string, 0, len(values))
	for _, value := range values {
		if value.IsNull() || value.IsUnknown() {
			continue
		}
		val := value.ValueString()
		out = append(out, &val)
	}

	return out, diags
}

// flattenStringList builds a list of strings from the SDK slice. A nil or
// empty slice stays a null list so omitted condition blocks do not show up as
// spurious diffs.
func flattenStringList(_ context.Context, values []*string) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics
	if len(values) == 0 {
		return types.ListNull(types.StringType), diags
	}

	out := make([]types.String, 0, len(values))
	for _, value := range values {
		if value == nil {
			continue
		}
		out = append(out, types.StringValue(*value))
	}

	if len(out) == 0 {
		return types.ListNull(types.StringType), diags
	}

	list, listDiags := types.ListValueFrom(context.Background(), types.StringType, out)
	diags.Append(listDiags...)
	if diags.HasError() {
		return types.ListNull(types.StringType), diags
	}
	return list, diags
}

// stringValueOf safely reads a string pointer.
func stringValueOf(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

// int64PtrValueOrPrev converts an *int64 reported by the ALB SDK into
// types.Int64, keeping null when the API omitted the field.
func int64PtrValueOrPrev(value *int64) types.Int64 {
	if value == nil {
		return types.Int64Null()
	}
	return types.Int64Value(*value)
}
