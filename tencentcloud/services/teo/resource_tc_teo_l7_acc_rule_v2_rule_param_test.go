package teo_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/teo"
)

// go test ./tencentcloud/services/teo/ -run "TestL7AccRuleV2RuleParam" -v -count=1 -gcflags="all=-l"

// TestL7AccRuleV2RuleParam_CreateWithRuleParam tests Create with the rule parameter
func TestL7AccRuleV2RuleParam_CreateWithRuleParam(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)

	// Mock CreateL7AccRules
	patches.ApplyMethodFunc(teoClient, "CreateL7AccRules", func(request *teov20220901.CreateL7AccRulesRequest) (*teov20220901.CreateL7AccRulesResponse, error) {
		resp := teov20220901.NewCreateL7AccRulesResponse()
		resp.Response = &teov20220901.CreateL7AccRulesResponseParams{
			RuleIds:   []*string{ptrString("rule-test123")},
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	// Mock DescribeTeoL7AccRuleById for the Read after Create
	patches.ApplyMethodFunc(&teo.TeoService{}, "DescribeTeoL7AccRuleById", func(ctx context.Context, zoneId string, ruleId string) (*teov20220901.DescribeL7AccRulesResponseParams, error) {
		resp := &teov20220901.DescribeL7AccRulesResponseParams{
			TotalCount: ptrInt64(1),
			Rules: []*teov20220901.RuleEngineItem{
				{
					RuleId:      ptrString("rule-test123"),
					Status:      ptrString("enable"),
					RuleName:    ptrString("test-rule"),
					Description: []*string{ptrString("test-desc")},
				},
			},
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoL7AccRuleV2()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-1234567890",
		"rule": []interface{}{
			map[string]interface{}{
				"status":      "enable",
				"rule_name":   "test-rule",
				"description": []interface{}{"test-desc"},
			},
		},
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-1234567890#rule-test123", d.Id())

	// Verify rule parameter is populated in Read
	ruleList := d.Get("rule").([]interface{})
	assert.Equal(t, 1, len(ruleList))
}

// TestL7AccRuleV2RuleParam_CreateWithTopLevelFields tests backward compatibility - Create with top-level fields
func TestL7AccRuleV2RuleParam_CreateWithTopLevelFields(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "CreateL7AccRules", func(request *teov20220901.CreateL7AccRulesRequest) (*teov20220901.CreateL7AccRulesResponse, error) {
		resp := teov20220901.NewCreateL7AccRulesResponse()
		resp.Response = &teov20220901.CreateL7AccRulesResponseParams{
			RuleIds:   []*string{ptrString("rule-test456")},
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(&teo.TeoService{}, "DescribeTeoL7AccRuleById", func(ctx context.Context, zoneId string, ruleId string) (*teov20220901.DescribeL7AccRulesResponseParams, error) {
		resp := &teov20220901.DescribeL7AccRulesResponseParams{
			TotalCount: ptrInt64(1),
			Rules: []*teov20220901.RuleEngineItem{
				{
					RuleId:      ptrString("rule-test456"),
					Status:      ptrString("enable"),
					RuleName:    ptrString("top-level-rule"),
					Description: []*string{ptrString("top-desc")},
				},
			},
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoL7AccRuleV2()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":     "zone-1234567890",
		"status":      "enable",
		"rule_name":   "top-level-rule",
		"description": []interface{}{"top-desc"},
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-1234567890#rule-test456", d.Id())
	assert.Equal(t, "enable", d.Get("status"))
	assert.Equal(t, "top-level-rule", d.Get("rule_name"))
}

// TestL7AccRuleV2RuleParam_ReadWithRuleParam tests Read populates rule parameter
func TestL7AccRuleV2RuleParam_ReadWithRuleParam(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(&teo.TeoService{}, "DescribeTeoL7AccRuleById", func(ctx context.Context, zoneId string, ruleId string) (*teov20220901.DescribeL7AccRulesResponseParams, error) {
		resp := &teov20220901.DescribeL7AccRulesResponseParams{
			TotalCount: ptrInt64(1),
			Rules: []*teov20220901.RuleEngineItem{
				{
					RuleId:       ptrString("rule-abc123"),
					Status:       ptrString("enable"),
					RuleName:     ptrString("my-rule"),
					Description:  []*string{ptrString("my-desc")},
					RulePriority: ptrInt64(10),
				},
			},
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoL7AccRuleV2()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-1234567890",
	})
	d.SetId("zone-1234567890#rule-abc123")

	err := res.Read(d, meta)
	assert.NoError(t, err)

	// Verify top-level fields
	assert.Equal(t, "enable", d.Get("status"))
	assert.Equal(t, "my-rule", d.Get("rule_name"))
	assert.Equal(t, "rule-abc123", d.Get("rule_id"))

	// Verify rule parameter is populated
	ruleList := d.Get("rule").([]interface{})
	assert.Equal(t, 1, len(ruleList))
}

// TestL7AccRuleV2RuleParam_UpdateWithRuleParam tests Update with the rule parameter
func TestL7AccRuleV2RuleParam_UpdateWithRuleParam(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "ModifyL7AccRule", func(request *teov20220901.ModifyL7AccRuleRequest) (*teov20220901.ModifyL7AccRuleResponse, error) {
		resp := teov20220901.NewModifyL7AccRuleResponse()
		resp.Response = &teov20220901.ModifyL7AccRuleResponseParams{
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(&teo.TeoService{}, "DescribeTeoL7AccRuleById", func(ctx context.Context, zoneId string, ruleId string) (*teov20220901.DescribeL7AccRulesResponseParams, error) {
		resp := &teov20220901.DescribeL7AccRulesResponseParams{
			TotalCount: ptrInt64(1),
			Rules: []*teov20220901.RuleEngineItem{
				{
					RuleId:      ptrString("rule-abc123"),
					Status:      ptrString("disable"),
					RuleName:    ptrString("updated-rule"),
					Description: []*string{ptrString("updated-desc")},
				},
			},
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoL7AccRuleV2()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-1234567890",
		"rule": []interface{}{
			map[string]interface{}{
				"status":      "disable",
				"rule_name":   "updated-rule",
				"description": []interface{}{"updated-desc"},
			},
		},
	})
	d.SetId("zone-1234567890#rule-abc123")

	err := res.Update(d, meta)
	assert.NoError(t, err)
}

// TestL7AccRuleV2RuleParam_UpdateWithTopLevelFields tests backward compatibility - Update with top-level fields
func TestL7AccRuleV2RuleParam_UpdateWithTopLevelFields(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "ModifyL7AccRule", func(request *teov20220901.ModifyL7AccRuleRequest) (*teov20220901.ModifyL7AccRuleResponse, error) {
		resp := teov20220901.NewModifyL7AccRuleResponse()
		resp.Response = &teov20220901.ModifyL7AccRuleResponseParams{
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(&teo.TeoService{}, "DescribeTeoL7AccRuleById", func(ctx context.Context, zoneId string, ruleId string) (*teov20220901.DescribeL7AccRulesResponseParams, error) {
		resp := &teov20220901.DescribeL7AccRulesResponseParams{
			TotalCount: ptrInt64(1),
			Rules: []*teov20220901.RuleEngineItem{
				{
					RuleId:      ptrString("rule-abc123"),
					Status:      ptrString("disable"),
					RuleName:    ptrString("updated-top-level-rule"),
					Description: []*string{ptrString("updated-top-desc")},
				},
			},
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoL7AccRuleV2()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":     "zone-1234567890",
		"status":      "disable",
		"rule_name":   "updated-top-level-rule",
		"description": []interface{}{"updated-top-desc"},
	})
	d.SetId("zone-1234567890#rule-abc123")

	err := res.Update(d, meta)
	assert.NoError(t, err)
}

// TestL7AccRuleV2RuleParam_CreateAPIError tests Create handles API error
func TestL7AccRuleV2RuleParam_CreateAPIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "CreateL7AccRules", func(request *teov20220901.CreateL7AccRulesRequest) (*teov20220901.CreateL7AccRulesResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=InvalidParameter, Message=Invalid zone_id")
	})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoL7AccRuleV2()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-invalid",
		"rule": []interface{}{
			map[string]interface{}{
				"status":    "enable",
				"rule_name": "test-rule",
			},
		},
	})

	err := res.Create(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "InvalidParameter")
}

// TestL7AccRuleV2RuleParam_ReadNotFound tests Read handles resource not found
func TestL7AccRuleV2RuleParam_ReadNotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)

	// Return nil to simulate not found
	patches.ApplyMethodFunc(&teo.TeoService{}, "DescribeTeoL7AccRuleById", func(ctx context.Context, zoneId string, ruleId string) (*teov20220901.DescribeL7AccRulesResponseParams, error) {
		return nil, nil
	})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoL7AccRuleV2()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-1234567890",
	})
	d.SetId("zone-1234567890#rule-nonexistent")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}

// TestL7AccRuleV2RuleParam_Schema validates the rule parameter schema definition
func TestL7AccRuleV2RuleParam_Schema(t *testing.T) {
	res := teo.ResourceTencentCloudTeoL7AccRuleV2()

	assert.NotNil(t, res)
	assert.Contains(t, res.Schema, "rule")

	ruleField := res.Schema["rule"]
	assert.Equal(t, schema.TypeList, ruleField.Type)
	assert.True(t, ruleField.Optional)
	assert.Equal(t, 1, ruleField.MaxItems)

	// Verify nested schema fields
	ruleElem := ruleField.Elem.(*schema.Resource)
	assert.Contains(t, ruleElem.Schema, "rule_id")
	assert.Contains(t, ruleElem.Schema, "status")
	assert.Contains(t, ruleElem.Schema, "rule_name")
	assert.Contains(t, ruleElem.Schema, "description")
	assert.Contains(t, ruleElem.Schema, "branches")

	// rule_id should be Computed
	assert.True(t, ruleElem.Schema["rule_id"].Computed)
	// status should be Optional
	assert.True(t, ruleElem.Schema["status"].Optional)
	// rule_name should be Optional
	assert.True(t, ruleElem.Schema["rule_name"].Optional)
	// description should be Optional TypeList
	assert.Equal(t, schema.TypeList, ruleElem.Schema["description"].Type)
	assert.True(t, ruleElem.Schema["description"].Optional)
	// branches should be Optional TypeList
	assert.Equal(t, schema.TypeList, ruleElem.Schema["branches"].Type)
	assert.True(t, ruleElem.Schema["branches"].Optional)
}

func ptrInt64(i int64) *int64 {
	return &i
}
