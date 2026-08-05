package teo_test

import (
	"fmt"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/teo"
)

// TestL7AccRulesCreate_Success tests successful rule creation
func TestL7AccRulesCreate_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "CreateL7AccRules", func(request *teov20220901.CreateL7AccRulesRequest) (*teov20220901.CreateL7AccRulesResponse, error) {
		assert.Equal(t, "zone-test-123", *request.ZoneId)
		assert.Len(t, request.Rules, 1)
		assert.Equal(t, "test-rule", *request.Rules[0].RuleName)
		assert.Equal(t, "enable", *request.Rules[0].Status)

		resp := teov20220901.NewCreateL7AccRulesResponse()
		resp.Response = &teov20220901.CreateL7AccRulesResponseParams{
			RuleIds:   []*string{ptrString("rule-test-456")},
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	// Mock DescribeL7AccRules for the Read call after Create
	patches.ApplyMethodFunc(teoClient, "DescribeL7AccRules", func(request *teov20220901.DescribeL7AccRulesRequest) (*teov20220901.DescribeL7AccRulesResponse, error) {
		resp := teov20220901.NewDescribeL7AccRulesResponse()
		resp.Response = &teov20220901.DescribeL7AccRulesResponseParams{
			TotalCount: ptrInt64(1),
			Rules: []*teov20220901.RuleEngineItem{
				{
					RuleId:       ptrString("rule-test-456"),
					RuleName:     ptrString("test-rule"),
					Status:       ptrString("enable"),
					RulePriority: ptrInt64(1),
					Description:  []*string{ptrString("test description")},
				},
			},
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoL7AccRules()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-test-123",
		"rules": []interface{}{
			map[string]interface{}{
				"rule_name":   "test-rule",
				"status":      "enable",
				"description": []interface{}{"test description"},
			},
		},
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-test-123#rule-test-456", d.Id())
}

// TestL7AccRulesCreate_EmptyResponse tests empty rule IDs response
func TestL7AccRulesCreate_EmptyResponse(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "CreateL7AccRules", func(request *teov20220901.CreateL7AccRulesRequest) (*teov20220901.CreateL7AccRulesResponse, error) {
		resp := teov20220901.NewCreateL7AccRulesResponse()
		resp.Response = &teov20220901.CreateL7AccRulesResponseParams{
			RuleIds:   []*string{},
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoL7AccRules()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-test-123",
		"rules": []interface{}{
			map[string]interface{}{
				"rule_name": "test-rule",
				"status":    "enable",
			},
		},
	})

	err := res.Create(d, meta)
	assert.Error(t, err)
}

// TestL7AccRulesCreate_APIError tests API error handling
func TestL7AccRulesCreate_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "CreateL7AccRules", func(request *teov20220901.CreateL7AccRulesRequest) (*teov20220901.CreateL7AccRulesResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=InvalidParameter, Message=Invalid zone ID")
	})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoL7AccRules()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-invalid",
		"rules": []interface{}{
			map[string]interface{}{
				"rule_name": "test-rule",
				"status":    "enable",
			},
		},
	})

	err := res.Create(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "InvalidParameter")
}

// TestL7AccRulesRead_Success tests successful rule read
func TestL7AccRulesRead_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeL7AccRules", func(request *teov20220901.DescribeL7AccRulesRequest) (*teov20220901.DescribeL7AccRulesResponse, error) {
		assert.Equal(t, "zone-test-123", *request.ZoneId)
		resp := teov20220901.NewDescribeL7AccRulesResponse()
		resp.Response = &teov20220901.DescribeL7AccRulesResponseParams{
			TotalCount: ptrInt64(1),
			Rules: []*teov20220901.RuleEngineItem{
				{
					RuleId:       ptrString("rule-test-456"),
					RuleName:     ptrString("test-rule"),
					Status:       ptrString("enable"),
					RulePriority: ptrInt64(1),
					Description:  []*string{ptrString("test description")},
				},
			},
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoL7AccRules()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-test-123",
	})
	d.SetId("zone-test-123#rule-test-456")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-test-123", d.Get("zone_id"))

	rules := d.Get("rules").([]interface{})
	assert.Len(t, rules, 1)
	ruleMap := rules[0].(map[string]interface{})
	assert.Equal(t, "test-rule", ruleMap["rule_name"])
	assert.Equal(t, "enable", ruleMap["status"])
	assert.Equal(t, "rule-test-456", ruleMap["rule_id"])
	assert.Equal(t, int64(1), ruleMap["rule_priority"])
}

// TestL7AccRulesRead_NotFound tests rule not found (removed from state)
func TestL7AccRulesRead_NotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeL7AccRules", func(request *teov20220901.DescribeL7AccRulesRequest) (*teov20220901.DescribeL7AccRulesResponse, error) {
		resp := teov20220901.NewDescribeL7AccRulesResponse()
		resp.Response = &teov20220901.DescribeL7AccRulesResponseParams{
			TotalCount: ptrInt64(0),
			Rules:      []*teov20220901.RuleEngineItem{},
			RequestId:  ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoL7AccRules()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-test-123",
	})
	d.SetId("zone-test-123#rule-deleted-456")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}

// TestL7AccRulesRead_BrokenID tests broken resource ID
func TestL7AccRulesRead_BrokenID(t *testing.T) {
	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoL7AccRules()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-test-123",
	})
	d.SetId("broken-id-without-separator")

	err := res.Read(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "id is broken")
}

// TestL7AccRulesUpdate_Success tests successful rule update
func TestL7AccRulesUpdate_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "ModifyL7AccRule", func(request *teov20220901.ModifyL7AccRuleRequest) (*teov20220901.ModifyL7AccRuleResponse, error) {
		assert.Equal(t, "zone-test-123", *request.ZoneId)
		assert.Equal(t, "rule-test-456", *request.Rule.RuleId)
		assert.Equal(t, "updated-rule", *request.Rule.RuleName)
		assert.Equal(t, "disable", *request.Rule.Status)

		resp := teov20220901.NewModifyL7AccRuleResponse()
		resp.Response = &teov20220901.ModifyL7AccRuleResponseParams{
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	// Mock DescribeL7AccRules for the Read call after Update
	patches.ApplyMethodFunc(teoClient, "DescribeL7AccRules", func(request *teov20220901.DescribeL7AccRulesRequest) (*teov20220901.DescribeL7AccRulesResponse, error) {
		resp := teov20220901.NewDescribeL7AccRulesResponse()
		resp.Response = &teov20220901.DescribeL7AccRulesResponseParams{
			TotalCount: ptrInt64(1),
			Rules: []*teov20220901.RuleEngineItem{
				{
					RuleId:       ptrString("rule-test-456"),
					RuleName:     ptrString("updated-rule"),
					Status:       ptrString("disable"),
					RulePriority: ptrInt64(1),
				},
			},
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoL7AccRules()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-test-123",
		"rules": []interface{}{
			map[string]interface{}{
				"rule_name": "updated-rule",
				"status":    "disable",
			},
		},
	})
	d.SetId("zone-test-123#rule-test-456")

	// Mark that rules has changed to trigger update
	d.MarkNewResource()

	err := res.Update(d, meta)
	assert.NoError(t, err)
}

// TestL7AccRulesDelete_Success tests successful rule deletion
func TestL7AccRulesDelete_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DeleteL7AccRules", func(request *teov20220901.DeleteL7AccRulesRequest) (*teov20220901.DeleteL7AccRulesResponse, error) {
		assert.Equal(t, "zone-test-123", *request.ZoneId)
		assert.Len(t, request.RuleIds, 1)
		assert.Equal(t, "rule-test-456", *request.RuleIds[0])

		resp := teov20220901.NewDeleteL7AccRulesResponse()
		resp.Response = &teov20220901.DeleteL7AccRulesResponseParams{
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoL7AccRules()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-test-123",
	})
	d.SetId("zone-test-123#rule-test-456")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
}

// TestL7AccRulesDelete_APIError tests API error during deletion
func TestL7AccRulesDelete_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DeleteL7AccRules", func(request *teov20220901.DeleteL7AccRulesRequest) (*teov20220901.DeleteL7AccRulesResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=ResourceNotFound, Message=Rule not found")
	})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoL7AccRules()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-test-123",
	})
	d.SetId("zone-test-123#rule-nonexistent")

	err := res.Delete(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ResourceNotFound")
}

// TestL7AccRules_Schema tests schema definition
func TestL7AccRules_Schema(t *testing.T) {
	res := teo.ResourceTencentCloudTeoL7AccRules()

	assert.NotNil(t, res)
	assert.NotNil(t, res.Create)
	assert.NotNil(t, res.Read)
	assert.NotNil(t, res.Update)
	assert.NotNil(t, res.Delete)
	assert.NotNil(t, res.Importer)

	assert.Contains(t, res.Schema, "zone_id")
	assert.Contains(t, res.Schema, "rules")

	// Check zone_id
	zoneId := res.Schema["zone_id"]
	assert.Equal(t, schema.TypeString, zoneId.Type)
	assert.True(t, zoneId.Required)
	assert.True(t, zoneId.ForceNew)

	// Check rules
	rules := res.Schema["rules"]
	assert.Equal(t, schema.TypeList, rules.Type)
	assert.True(t, rules.Optional)

	// Check nested rule fields
	ruleResource := rules.Elem.(*schema.Resource)
	assert.Contains(t, ruleResource.Schema, "rule_name")
	assert.Contains(t, ruleResource.Schema, "status")
	assert.Contains(t, ruleResource.Schema, "description")
	assert.Contains(t, ruleResource.Schema, "branches")
	assert.Contains(t, ruleResource.Schema, "rule_id")
	assert.Contains(t, ruleResource.Schema, "rule_priority")

	// Check rule_id is computed
	ruleId := ruleResource.Schema["rule_id"]
	assert.Equal(t, schema.TypeString, ruleId.Type)
	assert.True(t, ruleId.Computed)

	// Check rule_priority is computed
	rulePriority := ruleResource.Schema["rule_priority"]
	assert.Equal(t, schema.TypeInt, rulePriority.Type)
	assert.True(t, rulePriority.Computed)
}
