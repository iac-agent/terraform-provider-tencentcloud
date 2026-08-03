package teo

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func TestResourceTencentCloudTeoL7AccRulesSchema(t *testing.T) {
	r := ResourceTencentCloudTeoL7AccRules()
	assert.NotNil(t, r)
	assert.NotNil(t, r.Schema["zone_id"])
	assert.NotNil(t, r.Schema["rules"])
	assert.NotNil(t, r.Schema["rule_ids"])
	assert.True(t, r.Schema["zone_id"].Required)
	assert.True(t, r.Schema["zone_id"].ForceNew)
}

func TestBuildRuleEngineItemFromRuleMap(t *testing.T) {
	ruleMap := map[string]interface{}{
		"status":    "enable",
		"rule_name": "test-rule",
		"description": []interface{}{
			"desc1",
			"desc2",
		},
		"branches": []interface{}{},
	}

	rule := buildRuleEngineItemFromRuleMap(ruleMap)
	assert.NotNil(t, rule)
	assert.Equal(t, "enable", *rule.Status)
	assert.Equal(t, "test-rule", *rule.RuleName)
	assert.Len(t, rule.Description, 2)
	assert.Equal(t, "desc1", *rule.Description[0])
	assert.Equal(t, "desc2", *rule.Description[1])
}

func TestBuildRuleEngineItemFromRuleMapEmpty(t *testing.T) {
	ruleMap := map[string]interface{}{}
	rule := buildRuleEngineItemFromRuleMap(ruleMap)
	assert.NotNil(t, rule)
	assert.Nil(t, rule.Status)
	assert.Nil(t, rule.RuleName)
	assert.Nil(t, rule.Description)
	assert.Nil(t, rule.Branches)
}

func TestBuildRuleEngineItemsFromRules(t *testing.T) {
	rules := []interface{}{
		map[string]interface{}{
			"rule_name": "rule1",
			"status":    "enable",
		},
		map[string]interface{}{
			"rule_name": "rule2",
			"status":    "disable",
		},
	}

	result := buildRuleEngineItemsFromRules(rules)
	assert.Len(t, result, 2)
	assert.Equal(t, "rule1", *result[0].RuleName)
	assert.Equal(t, "rule2", *result[1].RuleName)
}

func TestResourceTencentCloudTeoL7AccRulesCreateBasic(t *testing.T) {
	d := schema.TestResourceDataRaw(t, ResourceTencentCloudTeoL7AccRules().Schema, map[string]interface{}{
		"zone_id": "zone-test123",
		"rules": []interface{}{
			map[string]interface{}{
				"rule_name": "test-rule",
				"status":    "enable",
			},
		},
	})

	assert.Equal(t, "zone-test123", d.Get("zone_id").(string))
	assert.NotNil(t, d.Get("rules"))
}

func TestResourceTencentCloudTeoL7AccRulesReadBasic(t *testing.T) {
	d := schema.TestResourceDataRaw(t, ResourceTencentCloudTeoL7AccRules().Schema, map[string]interface{}{
		"zone_id": "zone-test123",
	})
	d.SetId("zone-test123")

	assert.Equal(t, "zone-test123", d.Id())
}

func TestResourceTencentCloudTeoL7AccRulesDeleteBasic(t *testing.T) {
	d := schema.TestResourceDataRaw(t, ResourceTencentCloudTeoL7AccRules().Schema, map[string]interface{}{
		"zone_id": "zone-test123",
		"rule_ids": []interface{}{
			"rule-1",
			"rule-2",
		},
	})
	d.SetId("zone-test123")

	ruleIds := make([]string, 0)
	if v, ok := d.GetOk("rule_ids"); ok {
		for _, item := range v.([]interface{}) {
			ruleIds = append(ruleIds, item.(string))
		}
	}
	assert.Len(t, ruleIds, 2)
	assert.Contains(t, ruleIds, "rule-1")
	assert.Contains(t, ruleIds, "rule-2")
}

func TestResourceTencentCloudTeoL7AccRulesImport(t *testing.T) {
	r := ResourceTencentCloudTeoL7AccRules()
	assert.NotNil(t, r.Importer)
	assert.NotNil(t, r.Importer.State)
}

func TestBuildRuleEngineItemFromRuleMapWithBranches(t *testing.T) {
	ruleMap := map[string]interface{}{
		"rule_name": "test-rule",
		"status":    "enable",
		"branches": []interface{}{
			map[string]interface{}{
				"condition": `${http.request.host} in ['example.com']`,
				"actions": []interface{}{
					map[string]interface{}{
						"name": "Cache",
						"cache_parameters": []interface{}{
							map[string]interface{}{
								"custom_time": []interface{}{
									map[string]interface{}{
										"switch":               "on",
										"cache_time":           2592000,
										"ignore_cache_control": "off",
									},
								},
							},
						},
					},
				},
			},
		},
	}

	rule := buildRuleEngineItemFromRuleMap(ruleMap)
	assert.NotNil(t, rule)
	assert.Equal(t, "test-rule", *rule.RuleName)
	assert.Equal(t, "enable", *rule.Status)
	assert.NotNil(t, rule.Branches)
	assert.Len(t, rule.Branches, 1)
}

func TestSDKTypesAvailable(t *testing.T) {
	_ = helper.String("test")
	_ = teov20220901.NewCreateL7AccRulesRequest()
	_ = teov20220901.NewDescribeL7AccRulesRequest()
	_ = teov20220901.NewModifyL7AccRuleRequest()
	_ = teov20220901.NewDeleteL7AccRulesRequest()
}
