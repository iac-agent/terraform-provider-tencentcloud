package teo

import (
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func flattenTeoFunctionFilters(filters []interface{}) []*teov20220901.Filter {
	result := make([]*teov20220901.Filter, 0, len(filters))
	for _, item := range filters {
		filterMap := item.(map[string]interface{})
		filter := &teov20220901.Filter{}
		if v, ok := filterMap["name"].(string); ok && v != "" {
			filter.Name = helper.String(v)
		}
		if v, ok := filterMap["values"]; ok {
			valuesSet := v.(*schema.Set).List()
			for _, val := range valuesSet {
				filter.Values = append(filter.Values, helper.String(val.(string)))
			}
		}
		result = append(result, filter)
	}
	return result
}

func flattenTeoFunctions(functions []*teov20220901.Function) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(functions))
	for _, f := range functions {
		functionMap := map[string]interface{}{}
		if f.FunctionId != nil {
			functionMap["function_id"] = *f.FunctionId
		}
		if f.ZoneId != nil {
			functionMap["zone_id"] = *f.ZoneId
		}
		if f.Name != nil {
			functionMap["name"] = *f.Name
		}
		if f.Remark != nil {
			functionMap["remark"] = *f.Remark
		}
		if f.Content != nil {
			functionMap["content"] = *f.Content
		}
		if f.Domain != nil {
			functionMap["domain"] = *f.Domain
		}
		if f.CreateTime != nil {
			functionMap["create_time"] = *f.CreateTime
		}
		if f.UpdateTime != nil {
			functionMap["update_time"] = *f.UpdateTime
		}
		result = append(result, functionMap)
	}
	return result
}
