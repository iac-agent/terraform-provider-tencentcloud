package alb

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	albv20251030 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/alb/v20251030"

	fwHelper "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/helper"
)

// Schema attribute-type mirrors used by expand/flatten for nested objects.

var albRuleActionAttrTypes = map[string]attr.Type{
	"order":                 types.Int64Type,
	"type":                  types.StringType,
	"fixed_response_config": types.ObjectType{AttrTypes: albFixedResponseAttrTypes},
	"insert_header_config":  types.ObjectType{AttrTypes: albInsertHeaderAttrTypes},
	"redirect_config":       types.ObjectType{AttrTypes: albRedirectAttrTypes},
	"remove_header_config":  types.ObjectType{AttrTypes: albRemoveHeaderAttrTypes},
	"rewrite_config":        types.ObjectType{AttrTypes: albRewriteAttrTypes},
	"target_group_config":   types.ObjectType{AttrTypes: albTargetGroupConfigAttrTypes},
}

var albFixedResponseAttrTypes = map[string]attr.Type{
	"http_code":    types.Int64Type,
	"content":      types.StringType,
	"content_type": types.StringType,
}

var albInsertHeaderAttrTypes = map[string]attr.Type{
	"key":        types.StringType,
	"value":      types.StringType,
	"value_type": types.StringType,
}

var albRedirectAttrTypes = map[string]attr.Type{
	"http_code": types.Int64Type,
	"host":      types.StringType,
	"path":      types.StringType,
	"port":      types.StringType,
	"protocol":  types.StringType,
	"query":     types.StringType,
}

var albRemoveHeaderAttrTypes = map[string]attr.Type{
	"key": types.StringType,
}

var albRewriteAttrTypes = map[string]attr.Type{
	"host":  types.StringType,
	"path":  types.StringType,
	"query": types.StringType,
}

var albTargetGroupConfigAttrTypes = map[string]attr.Type{
	"target_groups":               types.ListType{ElemType: types.ObjectType{AttrTypes: albTargetGroupTupleAttrTypes}},
	"target_group_sticky_session": types.ObjectType{AttrTypes: albStickySessionAttrTypes},
}

var albTargetGroupTupleAttrTypes = map[string]attr.Type{
	"target_group_id": types.StringType,
	"weight":          types.Int64Type,
}

var albStickySessionAttrTypes = map[string]attr.Type{
	"enabled": types.BoolType,
	"timeout": types.Int64Type,
}

var albRuleConditionAttrTypes = map[string]attr.Type{
	"type":                types.StringType,
	"cookie_config":       types.ListType{ElemType: types.ObjectType{AttrTypes: albCookieAttrTypes}},
	"header_config":       types.ObjectType{AttrTypes: albHeaderAttrTypes},
	"host_config":         types.ListType{ElemType: types.StringType},
	"method_config":       types.ListType{ElemType: types.StringType},
	"path_config":         types.ListType{ElemType: types.StringType},
	"query_string_config": types.ListType{ElemType: types.ObjectType{AttrTypes: albQueryStringAttrTypes}},
	"source_ip_config":    types.ListType{ElemType: types.StringType},
}

var albCookieAttrTypes = map[string]attr.Type{
	"key":   types.StringType,
	"value": types.StringType,
}

var albHeaderAttrTypes = map[string]attr.Type{
	"key":    types.StringType,
	"values": types.ListType{ElemType: types.StringType},
}

var albQueryStringAttrTypes = map[string]attr.Type{
	"key":   types.StringType,
	"value": types.StringType,
}

var albRuleAttrTypes = map[string]attr.Type{
	"actions":     types.ListType{ElemType: types.ObjectType{AttrTypes: albRuleActionAttrTypes}},
	"conditions":  types.ListType{ElemType: types.ObjectType{AttrTypes: albRuleConditionAttrTypes}},
	"priority":    types.Int64Type,
	"direction":   types.StringType,
	"rule_name":   types.StringType,
	"tags":        types.ListType{ElemType: types.ObjectType{AttrTypes: albTagAttrTypes}},
	"rule_id":     types.StringType,
	"create_time": types.StringType,
	"modify_time": types.StringType,
	"status":      types.StringType,
}

var albTagAttrTypes = map[string]attr.Type{
	"tag_key":   types.StringType,
	"tag_value": types.StringType,
}

// ---- expand: Terraform model -> SDK request ----

func expandAlbRules(ctx context.Context, list types.List, diags *diag.Diagnostics) []*albv20251030.RuleInput {
	var models []albRuleModel
	diags.Append(list.ElementsAs(ctx, &models, false)...)
	if diags.HasError() {
		return nil
	}

	rules := make([]*albv20251030.RuleInput, 0, len(models))
	for i := range models {
		rules = append(rules, &albv20251030.RuleInput{
			Actions:    expandAlbActions(ctx, models[i].Actions, diags),
			Conditions: expandAlbConditions(ctx, models[i].Conditions, diags),
			Priority:   fwHelper.Int64PointerFromValue(models[i].Priority),
			Direction:  fwHelper.StringPointerFromValue(models[i].Direction),
			RuleName:   fwHelper.StringPointerFromValue(models[i].RuleName),
			Tags:       expandAlbTags(ctx, models[i].Tags, diags),
		})
	}

	return rules
}

func expandAlbActions(ctx context.Context, list types.List, diags *diag.Diagnostics) []*albv20251030.RuleAction {
	var models []albRuleActionModel
	diags.Append(list.ElementsAs(ctx, &models, false)...)
	if diags.HasError() {
		return nil
	}

	actions := make([]*albv20251030.RuleAction, 0, len(models))
	for i := range models {
		action := &albv20251030.RuleAction{
			Order: fwHelper.Int64PointerFromValue(models[i].Order),
			Type:  fwHelper.StringPointerFromValue(models[i].Type),
		}
		action.FixedResponseConfig = expandAlbFixedResponse(ctx, models[i].FixedResponseConfig, diags)
		action.InsertHeaderConfig = expandAlbInsertHeader(ctx, models[i].InsertHeaderConfig, diags)
		action.RedirectConfig = expandAlbRedirect(ctx, models[i].RedirectConfig, diags)
		action.RemoveHeaderConfig = expandAlbRemoveHeader(ctx, models[i].RemoveHeaderConfig, diags)
		action.RewriteConfig = expandAlbRewrite(ctx, models[i].RewriteConfig, diags)
		action.TargetGroupConfig = expandAlbTargetGroupConfig(ctx, models[i].TargetGroupConfig, diags)
		actions = append(actions, action)
	}

	return actions
}

func expandAlbFixedResponse(ctx context.Context, obj types.Object, diags *diag.Diagnostics) *albv20251030.FixedResponseInfo {
	if obj.IsNull() || obj.IsUnknown() {
		return nil
	}

	var m struct {
		HttpCode    types.Int64  `tfsdk:"http_code"`
		Content     types.String `tfsdk:"content"`
		ContentType types.String `tfsdk:"content_type"`
	}
	diags.Append(obj.As(ctx, &m, basetypesObjectAsOptions())...)
	if diags.HasError() {
		return nil
	}

	return &albv20251030.FixedResponseInfo{
		HttpCode:    fwHelper.Int64PointerFromValue(m.HttpCode),
		Content:     fwHelper.StringPointerFromValue(m.Content),
		ContentType: fwHelper.StringPointerFromValue(m.ContentType),
	}
}

func expandAlbInsertHeader(ctx context.Context, obj types.Object, diags *diag.Diagnostics) *albv20251030.InsertHTTPHeaderInfo {
	if obj.IsNull() || obj.IsUnknown() {
		return nil
	}

	var m struct {
		Key       types.String `tfsdk:"key"`
		Value     types.String `tfsdk:"value"`
		ValueType types.String `tfsdk:"value_type"`
	}
	diags.Append(obj.As(ctx, &m, basetypesObjectAsOptions())...)
	if diags.HasError() {
		return nil
	}

	return &albv20251030.InsertHTTPHeaderInfo{
		Key:       fwHelper.StringPointerFromValue(m.Key),
		Value:     fwHelper.StringPointerFromValue(m.Value),
		ValueType: fwHelper.StringPointerFromValue(m.ValueType),
	}
}

func expandAlbRedirect(ctx context.Context, obj types.Object, diags *diag.Diagnostics) *albv20251030.HTTPRedirectInfo {
	if obj.IsNull() || obj.IsUnknown() {
		return nil
	}

	var m struct {
		HttpCode types.Int64  `tfsdk:"http_code"`
		Host     types.String `tfsdk:"host"`
		Path     types.String `tfsdk:"path"`
		Port     types.String `tfsdk:"port"`
		Protocol types.String `tfsdk:"protocol"`
		Query    types.String `tfsdk:"query"`
	}
	diags.Append(obj.As(ctx, &m, basetypesObjectAsOptions())...)
	if diags.HasError() {
		return nil
	}

	return &albv20251030.HTTPRedirectInfo{
		HttpCode: fwHelper.Int64PointerFromValue(m.HttpCode),
		Host:     fwHelper.StringPointerFromValue(m.Host),
		Path:     fwHelper.StringPointerFromValue(m.Path),
		Port:     fwHelper.StringPointerFromValue(m.Port),
		Protocol: fwHelper.StringPointerFromValue(m.Protocol),
		Query:    fwHelper.StringPointerFromValue(m.Query),
	}
}

func expandAlbRemoveHeader(ctx context.Context, obj types.Object, diags *diag.Diagnostics) *albv20251030.RemoveHTTPHeaderInfo {
	if obj.IsNull() || obj.IsUnknown() {
		return nil
	}

	var m struct {
		Key types.String `tfsdk:"key"`
	}
	diags.Append(obj.As(ctx, &m, basetypesObjectAsOptions())...)
	if diags.HasError() {
		return nil
	}

	return &albv20251030.RemoveHTTPHeaderInfo{Key: fwHelper.StringPointerFromValue(m.Key)}
}

func expandAlbRewrite(ctx context.Context, obj types.Object, diags *diag.Diagnostics) *albv20251030.HTTPRewriteInfo {
	if obj.IsNull() || obj.IsUnknown() {
		return nil
	}

	var m struct {
		Host  types.String `tfsdk:"host"`
		Path  types.String `tfsdk:"path"`
		Query types.String `tfsdk:"query"`
	}
	diags.Append(obj.As(ctx, &m, basetypesObjectAsOptions())...)
	if diags.HasError() {
		return nil
	}

	return &albv20251030.HTTPRewriteInfo{
		Host:  fwHelper.StringPointerFromValue(m.Host),
		Path:  fwHelper.StringPointerFromValue(m.Path),
		Query: fwHelper.StringPointerFromValue(m.Query),
	}
}

func expandAlbTargetGroupConfig(ctx context.Context, obj types.Object, diags *diag.Diagnostics) *albv20251030.TargetGroupConfig {
	if obj.IsNull() || obj.IsUnknown() {
		return nil
	}

	var m struct {
		TargetGroups             types.List   `tfsdk:"target_groups"`
		TargetGroupStickySession types.Object `tfsdk:"target_group_sticky_session"`
	}
	diags.Append(obj.As(ctx, &m, basetypesObjectAsOptions())...)
	if diags.HasError() {
		return nil
	}

	result := &albv20251030.TargetGroupConfig{}

	var tuples []struct {
		TargetGroupId types.String `tfsdk:"target_group_id"`
		Weight        types.Int64  `tfsdk:"weight"`
	}
	diags.Append(m.TargetGroups.ElementsAs(ctx, &tuples, false)...)
	if !diags.HasError() {
		result.TargetGroups = make([]*albv20251030.TargetGroupTuple, 0, len(tuples))
		for i := range tuples {
			result.TargetGroups = append(result.TargetGroups, &albv20251030.TargetGroupTuple{
				TargetGroupId: fwHelper.StringPointerFromValue(tuples[i].TargetGroupId),
				Weight:        fwHelper.Int64PointerFromValue(tuples[i].Weight),
			})
		}
	}

	if !m.TargetGroupStickySession.IsNull() && !m.TargetGroupStickySession.IsUnknown() {
		var sticky struct {
			Enabled types.Bool  `tfsdk:"enabled"`
			Timeout types.Int64 `tfsdk:"timeout"`
		}
		diags.Append(m.TargetGroupStickySession.As(ctx, &sticky, basetypesObjectAsOptions())...)
		if !diags.HasError() {
			result.TargetGroupStickySession = &albv20251030.TargetGroupStickySession{
				Enabled: fwHelper.BoolPointerFromValue(sticky.Enabled),
				Timeout: fwHelper.Int64PointerFromValue(sticky.Timeout),
			}
		}
	}

	return result
}

func expandAlbConditions(ctx context.Context, list types.List, diags *diag.Diagnostics) []*albv20251030.RuleCondition {
	var models []albRuleConditionModel
	diags.Append(list.ElementsAs(ctx, &models, false)...)
	if diags.HasError() {
		return nil
	}

	conditions := make([]*albv20251030.RuleCondition, 0, len(models))
	for i := range models {
		condition := &albv20251030.RuleCondition{
			Type:              fwHelper.StringPointerFromValue(models[i].Type),
			CookieConfig:      expandAlbCookies(ctx, models[i].CookieConfig, diags),
			HeaderConfig:      expandAlbHeader(ctx, models[i].HeaderConfig, diags),
			HostConfig:        fwHelper.PointersFromList[string](ctx, models[i].HostConfig, diags),
			MethodConfig:      fwHelper.PointersFromList[string](ctx, models[i].MethodConfig, diags),
			PathConfig:        fwHelper.PointersFromList[string](ctx, models[i].PathConfig, diags),
			QueryStringConfig: expandAlbQueryStrings(ctx, models[i].QueryStringConfig, diags),
			SourceIpConfig:    fwHelper.PointersFromList[string](ctx, models[i].SourceIpConfig, diags),
		}
		conditions = append(conditions, condition)
	}

	return conditions
}

func expandAlbCookies(ctx context.Context, list types.List, diags *diag.Diagnostics) []*albv20251030.HTTPCookieInfo {
	if list.IsNull() || list.IsUnknown() {
		return nil
	}

	var models []struct {
		Key   types.String `tfsdk:"key"`
		Value types.String `tfsdk:"value"`
	}
	diags.Append(list.ElementsAs(ctx, &models, false)...)
	if diags.HasError() {
		return nil
	}

	result := make([]*albv20251030.HTTPCookieInfo, 0, len(models))
	for i := range models {
		result = append(result, &albv20251030.HTTPCookieInfo{
			Key:   fwHelper.StringPointerFromValue(models[i].Key),
			Value: fwHelper.StringPointerFromValue(models[i].Value),
		})
	}

	return result
}

func expandAlbHeader(ctx context.Context, obj types.Object, diags *diag.Diagnostics) *albv20251030.HTTPHeaderInfo {
	if obj.IsNull() || obj.IsUnknown() {
		return nil
	}

	var m struct {
		Key    types.String `tfsdk:"key"`
		Values types.List   `tfsdk:"values"`
	}
	diags.Append(obj.As(ctx, &m, basetypesObjectAsOptions())...)
	if diags.HasError() {
		return nil
	}

	return &albv20251030.HTTPHeaderInfo{
		Key:    fwHelper.StringPointerFromValue(m.Key),
		Values: fwHelper.PointersFromList[string](ctx, m.Values, diags),
	}
}

func expandAlbQueryStrings(ctx context.Context, list types.List, diags *diag.Diagnostics) []*albv20251030.HTTPQueryStringInfo {
	if list.IsNull() || list.IsUnknown() {
		return nil
	}

	var models []struct {
		Key   types.String `tfsdk:"key"`
		Value types.String `tfsdk:"value"`
	}
	diags.Append(list.ElementsAs(ctx, &models, false)...)
	if diags.HasError() {
		return nil
	}

	result := make([]*albv20251030.HTTPQueryStringInfo, 0, len(models))
	for i := range models {
		result = append(result, &albv20251030.HTTPQueryStringInfo{
			Key:   fwHelper.StringPointerFromValue(models[i].Key),
			Value: fwHelper.StringPointerFromValue(models[i].Value),
		})
	}

	return result
}

func expandAlbTags(ctx context.Context, list types.List, diags *diag.Diagnostics) []*albv20251030.TagInfo {
	if list.IsNull() || list.IsUnknown() {
		return nil
	}

	var models []struct {
		TagKey   types.String `tfsdk:"tag_key"`
		TagValue types.String `tfsdk:"tag_value"`
	}
	diags.Append(list.ElementsAs(ctx, &models, false)...)
	if diags.HasError() {
		return nil
	}

	tags := make([]*albv20251030.TagInfo, 0, len(models))
	for i := range models {
		tags = append(tags, &albv20251030.TagInfo{
			TagKey:   fwHelper.StringPointerFromValue(models[i].TagKey),
			TagValue: fwHelper.StringPointerFromValue(models[i].TagValue),
		})
	}

	return tags
}

func expandAlbRuleModifies(ctx context.Context, list types.List, ruleId string, diags *diag.Diagnostics) []*albv20251030.RuleModify {
	var models []albRuleModel
	diags.Append(list.ElementsAs(ctx, &models, false)...)
	if diags.HasError() {
		return nil
	}

	rules := make([]*albv20251030.RuleModify, 0, len(models))
	for i := range models {
		rules = append(rules, &albv20251030.RuleModify{
			Actions:    expandAlbActions(ctx, models[i].Actions, diags),
			Conditions: expandAlbConditions(ctx, models[i].Conditions, diags),
			Priority:   fwHelper.Int64PointerFromValue(models[i].Priority),
			RuleId:     fwHelper.String(ruleId),
			RuleName:   fwHelper.StringPointerFromValue(models[i].RuleName),
		})
	}

	return rules
}

// ---- flatten: SDK response -> Terraform model ----

func flattenAlbRules(ctx context.Context, item *albv20251030.RuleOutput, diags *diag.Diagnostics) types.List {
	value, objectDiags := types.ObjectValue(albRuleAttrTypes, map[string]attr.Value{
		"actions":     flattenAlbActions(ctx, item.Actions, diags),
		"conditions":  flattenAlbConditions(ctx, item.Conditions, diags),
		"priority":    fwHelper.Int64ValueOrNull(item.Priority),
		"direction":   fwHelper.StringPointerValueOrNull(item.Direction),
		"rule_name":   fwHelper.StringPointerValueOrNull(item.RuleName),
		"tags":        flattenAlbTags(ctx, item.Tags, diags),
		"rule_id":     fwHelper.StringPointerValueOrNull(item.RuleId),
		"create_time": fwHelper.StringPointerValueOrNull(item.CreateTime),
		"modify_time": fwHelper.StringPointerValueOrNull(item.ModifyTime),
		"status":      fwHelper.StringPointerValueOrNull(item.Status),
	})
	diags.Append(objectDiags...)

	return types.ListValueMust(types.ObjectType{AttrTypes: albRuleAttrTypes}, []attr.Value{value})
}

func flattenAlbRuleModels(ctx context.Context, models []albRuleModel, diags *diag.Diagnostics) types.List {
	values := make([]attr.Value, 0, len(models))
	for i := range models {
		value, objectDiags := types.ObjectValue(albRuleAttrTypes, map[string]attr.Value{
			"actions":     models[i].Actions,
			"conditions":  models[i].Conditions,
			"priority":    models[i].Priority,
			"direction":   models[i].Direction,
			"rule_name":   models[i].RuleName,
			"tags":        models[i].Tags,
			"rule_id":     models[i].RuleId,
			"create_time": models[i].CreateTime,
			"modify_time": models[i].ModifyTime,
			"status":      models[i].Status,
		})
		diags.Append(objectDiags...)
		values = append(values, value)
	}

	return types.ListValueMust(types.ObjectType{AttrTypes: albRuleAttrTypes}, values)
}

func flattenAlbActions(ctx context.Context, actions []*albv20251030.RuleAction, diags *diag.Diagnostics) types.List {
	values := make([]attr.Value, 0, len(actions))
	for _, action := range actions {
		if action == nil {
			continue
		}

		value, objectDiags := types.ObjectValue(albRuleActionAttrTypes, map[string]attr.Value{
			"order":                 fwHelper.Int64ValueOrNull(action.Order),
			"type":                  fwHelper.StringPointerValueOrNull(action.Type),
			"fixed_response_config": flattenAlbFixedResponse(ctx, action.FixedResponseConfig, diags),
			"insert_header_config":  flattenAlbInsertHeader(ctx, action.InsertHeaderConfig, diags),
			"redirect_config":       flattenAlbRedirect(ctx, action.RedirectConfig, diags),
			"remove_header_config":  flattenAlbRemoveHeader(ctx, action.RemoveHeaderConfig, diags),
			"rewrite_config":        flattenAlbRewrite(ctx, action.RewriteConfig, diags),
			"target_group_config":   flattenAlbTargetGroupConfig(ctx, action.TargetGroupConfig, diags),
		})
		diags.Append(objectDiags...)
		values = append(values, value)
	}

	return types.ListValueMust(types.ObjectType{AttrTypes: albRuleActionAttrTypes}, values)
}

func flattenAlbFixedResponse(_ context.Context, cfg *albv20251030.FixedResponseInfo, diags *diag.Diagnostics) types.Object {
	if cfg == nil {
		return types.ObjectNull(albFixedResponseAttrTypes)
	}

	value, objectDiags := types.ObjectValue(albFixedResponseAttrTypes, map[string]attr.Value{
		"http_code":    fwHelper.Int64ValueOrNull(cfg.HttpCode),
		"content":      fwHelper.StringPointerValueOrNull(cfg.Content),
		"content_type": fwHelper.StringPointerValueOrNull(cfg.ContentType),
	})
	diags.Append(objectDiags...)

	return value
}

func flattenAlbInsertHeader(_ context.Context, cfg *albv20251030.InsertHTTPHeaderInfo, diags *diag.Diagnostics) types.Object {
	if cfg == nil {
		return types.ObjectNull(albInsertHeaderAttrTypes)
	}

	value, objectDiags := types.ObjectValue(albInsertHeaderAttrTypes, map[string]attr.Value{
		"key":        fwHelper.StringPointerValueOrNull(cfg.Key),
		"value":      fwHelper.StringPointerValueOrNull(cfg.Value),
		"value_type": fwHelper.StringPointerValueOrNull(cfg.ValueType),
	})
	diags.Append(objectDiags...)

	return value
}

func flattenAlbRedirect(_ context.Context, cfg *albv20251030.HTTPRedirectInfo, diags *diag.Diagnostics) types.Object {
	if cfg == nil {
		return types.ObjectNull(albRedirectAttrTypes)
	}

	value, objectDiags := types.ObjectValue(albRedirectAttrTypes, map[string]attr.Value{
		"http_code": fwHelper.Int64ValueOrNull(cfg.HttpCode),
		"host":      fwHelper.StringPointerValueOrNull(cfg.Host),
		"path":      fwHelper.StringPointerValueOrNull(cfg.Path),
		"port":      fwHelper.StringPointerValueOrNull(cfg.Port),
		"protocol":  fwHelper.StringPointerValueOrNull(cfg.Protocol),
		"query":     fwHelper.StringPointerValueOrNull(cfg.Query),
	})
	diags.Append(objectDiags...)

	return value
}

func flattenAlbRemoveHeader(_ context.Context, cfg *albv20251030.RemoveHTTPHeaderInfo, diags *diag.Diagnostics) types.Object {
	if cfg == nil {
		return types.ObjectNull(albRemoveHeaderAttrTypes)
	}

	value, objectDiags := types.ObjectValue(albRemoveHeaderAttrTypes, map[string]attr.Value{
		"key": fwHelper.StringPointerValueOrNull(cfg.Key),
	})
	diags.Append(objectDiags...)

	return value
}

func flattenAlbRewrite(_ context.Context, cfg *albv20251030.HTTPRewriteInfo, diags *diag.Diagnostics) types.Object {
	if cfg == nil {
		return types.ObjectNull(albRewriteAttrTypes)
	}

	value, objectDiags := types.ObjectValue(albRewriteAttrTypes, map[string]attr.Value{
		"host":  fwHelper.StringPointerValueOrNull(cfg.Host),
		"path":  fwHelper.StringPointerValueOrNull(cfg.Path),
		"query": fwHelper.StringPointerValueOrNull(cfg.Query),
	})
	diags.Append(objectDiags...)

	return value
}

func flattenAlbTargetGroupConfig(ctx context.Context, cfg *albv20251030.TargetGroupConfig, diags *diag.Diagnostics) types.Object {
	if cfg == nil {
		return types.ObjectNull(albTargetGroupConfigAttrTypes)
	}

	tuples := make([]attr.Value, 0, len(cfg.TargetGroups))
	for _, tuple := range cfg.TargetGroups {
		if tuple == nil {
			continue
		}

		value, objectDiags := types.ObjectValue(albTargetGroupTupleAttrTypes, map[string]attr.Value{
			"target_group_id": fwHelper.StringPointerValueOrNull(tuple.TargetGroupId),
			"weight":          fwHelper.Int64ValueOrNull(tuple.Weight),
		})
		diags.Append(objectDiags...)
		tuples = append(tuples, value)
	}

	var sticky types.Object
	if cfg.TargetGroupStickySession == nil {
		sticky = types.ObjectNull(albStickySessionAttrTypes)
	} else {
		value, objectDiags := types.ObjectValue(albStickySessionAttrTypes, map[string]attr.Value{
			"enabled": fwHelper.BoolValueOrNull(cfg.TargetGroupStickySession.Enabled),
			"timeout": fwHelper.Int64ValueOrNull(cfg.TargetGroupStickySession.Timeout),
		})
		diags.Append(objectDiags...)
		sticky = value
	}

	value, objectDiags := types.ObjectValue(albTargetGroupConfigAttrTypes, map[string]attr.Value{
		"target_groups":               types.ListValueMust(types.ObjectType{AttrTypes: albTargetGroupTupleAttrTypes}, tuples),
		"target_group_sticky_session": sticky,
	})
	diags.Append(objectDiags...)

	return value
}

func flattenAlbConditions(ctx context.Context, conditions []*albv20251030.RuleCondition, diags *diag.Diagnostics) types.List {
	values := make([]attr.Value, 0, len(conditions))
	for _, condition := range conditions {
		if condition == nil {
			continue
		}

		value, objectDiags := types.ObjectValue(albRuleConditionAttrTypes, map[string]attr.Value{
			"type":                fwHelper.StringPointerValueOrNull(condition.Type),
			"cookie_config":       flattenAlbCookies(ctx, condition.CookieConfig, diags),
			"header_config":       flattenAlbHeader(ctx, condition.HeaderConfig, diags),
			"host_config":         fwHelper.ListFromPointers(ctx, types.StringType, condition.HostConfig, diags),
			"method_config":       fwHelper.ListFromPointers(ctx, types.StringType, condition.MethodConfig, diags),
			"path_config":         fwHelper.ListFromPointers(ctx, types.StringType, condition.PathConfig, diags),
			"query_string_config": flattenAlbQueryStrings(ctx, condition.QueryStringConfig, diags),
			"source_ip_config":    fwHelper.ListFromPointers(ctx, types.StringType, condition.SourceIpConfig, diags),
		})
		diags.Append(objectDiags...)
		values = append(values, value)
	}

	return types.ListValueMust(types.ObjectType{AttrTypes: albRuleConditionAttrTypes}, values)
}

func flattenAlbCookies(_ context.Context, cookies []*albv20251030.HTTPCookieInfo, diags *diag.Diagnostics) types.List {
	values := make([]attr.Value, 0, len(cookies))
	for _, cookie := range cookies {
		if cookie == nil {
			continue
		}

		value, objectDiags := types.ObjectValue(albCookieAttrTypes, map[string]attr.Value{
			"key":   fwHelper.StringPointerValueOrNull(cookie.Key),
			"value": fwHelper.StringPointerValueOrNull(cookie.Value),
		})
		diags.Append(objectDiags...)
		values = append(values, value)
	}

	return types.ListValueMust(types.ObjectType{AttrTypes: albCookieAttrTypes}, values)
}

func flattenAlbHeader(_ context.Context, cfg *albv20251030.HTTPHeaderInfo, diags *diag.Diagnostics) types.Object {
	if cfg == nil {
		return types.ObjectNull(albHeaderAttrTypes)
	}

	value, objectDiags := types.ObjectValue(albHeaderAttrTypes, map[string]attr.Value{
		"key":    fwHelper.StringPointerValueOrNull(cfg.Key),
		"values": fwHelper.ListFromPointers(context.Background(), types.StringType, cfg.Values, diags),
	})
	diags.Append(objectDiags...)

	return value
}

func flattenAlbQueryStrings(_ context.Context, items []*albv20251030.HTTPQueryStringInfo, diags *diag.Diagnostics) types.List {
	values := make([]attr.Value, 0, len(items))
	for _, item := range items {
		if item == nil {
			continue
		}

		value, objectDiags := types.ObjectValue(albQueryStringAttrTypes, map[string]attr.Value{
			"key":   fwHelper.StringPointerValueOrNull(item.Key),
			"value": fwHelper.StringPointerValueOrNull(item.Value),
		})
		diags.Append(objectDiags...)
		values = append(values, value)
	}

	return types.ListValueMust(types.ObjectType{AttrTypes: albQueryStringAttrTypes}, values)
}

func flattenAlbTags(_ context.Context, tags []*albv20251030.TagInfo, diags *diag.Diagnostics) types.List {
	values := make([]attr.Value, 0, len(tags))
	for _, tag := range tags {
		if tag == nil {
			continue
		}

		value, objectDiags := types.ObjectValue(albTagAttrTypes, map[string]attr.Value{
			"tag_key":   fwHelper.StringPointerValueOrNull(tag.TagKey),
			"tag_value": fwHelper.StringPointerValueOrNull(tag.TagValue),
		})
		diags.Append(objectDiags...)
		values = append(values, value)
	}

	return types.ListValueMust(types.ObjectType{AttrTypes: albTagAttrTypes}, values)
}

// basetypesObjectAsOptions returns the framework options used by Object.As when
// decoding nested objects into local structs.
func basetypesObjectAsOptions() basetypes.ObjectAsOptions {
	return basetypes.ObjectAsOptions{UnhandledUnknownAsEmpty: false, UnhandledNullAsEmpty: false}
}
