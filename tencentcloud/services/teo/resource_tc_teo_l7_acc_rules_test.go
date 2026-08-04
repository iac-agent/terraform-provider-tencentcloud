package teo_test

import (
	"fmt"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/teo"
)

// go test ./tencentcloud/services/teo/ -run "TestTeoL7AccRules" -v -count=1 -gcflags="all=-l"

// TestTeoL7AccRulesCreate_Success tests Create calls CreateL7AccRules and sets ID
func TestTeoL7AccRulesCreate_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)

	// Patch CreateL7AccRules to return success
	patches.ApplyMethodFunc(teoClient, "CreateL7AccRules", func(request *teov20220901.CreateL7AccRulesRequest) (*teov20220901.CreateL7AccRulesResponse, error) {
		assert.Equal(t, "zone-test-123", *request.ZoneId)
		assert.Len(t, request.Rules, 1)
		assert.Equal(t, "test-rule", *request.Rules[0].RuleName)
		resp := teov20220901.NewCreateL7AccRulesResponse()
		resp.Response = &teov20220901.CreateL7AccRulesResponseParams{
			RuleIds:   []*string{ptrString("rule-test-123")},
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	// Patch DescribeL7AccRules to return the created rule
	patches.ApplyMethodFunc(teoClient, "DescribeL7AccRules", func(request *teov20220901.DescribeL7AccRulesRequest) (*teov20220901.DescribeL7AccRulesResponse, error) {
		assert.Equal(t, "zone-test-123", *request.ZoneId)
		resp := teov20220901.NewDescribeL7AccRulesResponse()
		resp.Response = &teov20220901.DescribeL7AccRulesResponseParams{
			TotalCount: ptrInt64(1),
			Rules: []*teov20220901.RuleEngineItem{
				{
					RuleId:       ptrString("rule-test-123"),
					RuleName:     ptrString("test-rule"),
					Status:       ptrString("enable"),
					RulePriority: ptrInt64(1),
					Description:  []*string{ptrString("test description")},
				},
			},
			RequestId: ptrString("fake-request-id-read"),
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
				"branches":    []interface{}{},
			},
		},
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-test-123", d.Id())
}

// TestTeoL7AccRulesCreate_EmptyRules tests Create returns error for empty rules
func TestTeoL7AccRulesCreate_EmptyRules(t *testing.T) {
	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoL7AccRules()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-test-123",
		"rules":   []interface{}{},
	})

	err := res.Create(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot be empty")
}

// TestTeoL7AccRulesCreate_APIError tests Create handles API error
func TestTeoL7AccRulesCreate_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "CreateL7AccRules", func(request *teov20220901.CreateL7AccRulesRequest) (*teov20220901.CreateL7AccRulesResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=InvalidParameter, Message=Invalid zone_id")
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
}

// TestTeoL7AccRulesCreate_NilResponse tests Create returns error when response is nil
func TestTeoL7AccRulesCreate_NilResponse(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "CreateL7AccRules", func(request *teov20220901.CreateL7AccRulesRequest) (*teov20220901.CreateL7AccRulesResponse, error) {
		resp := teov20220901.NewCreateL7AccRulesResponse()
		// Response is nil
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
	assert.Contains(t, err.Error(), "response is nil")
}

// TestTeoL7AccRulesRead_Success tests Read fetches and sets rules
func TestTeoL7AccRulesRead_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeL7AccRules", func(request *teov20220901.DescribeL7AccRulesRequest) (*teov20220901.DescribeL7AccRulesResponse, error) {
		assert.Equal(t, "zone-test-123", *request.ZoneId)
		resp := teov20220901.NewDescribeL7AccRulesResponse()
		resp.Response = &teov20220901.DescribeL7AccRulesResponseParams{
			TotalCount: ptrInt64(2),
			Rules: []*teov20220901.RuleEngineItem{
				{
					RuleId:       ptrString("rule-1"),
					RuleName:     ptrString("rule-one"),
					Status:       ptrString("enable"),
					RulePriority: ptrInt64(1),
					Description:  []*string{ptrString("desc 1")},
				},
				{
					RuleId:       ptrString("rule-2"),
					RuleName:     ptrString("rule-two"),
					Status:       ptrString("disable"),
					RulePriority: ptrInt64(2),
					Description:  []*string{ptrString("desc 2")},
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
	d.SetId("zone-test-123")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-test-123", d.Id())

	rules := d.Get("rules").([]interface{})
	assert.Len(t, rules, 2)

	ruleIds := d.Get("rule_ids").([]interface{})
	assert.Len(t, ruleIds, 2)
}

// TestTeoL7AccRulesRead_EmptyRules tests Read with empty rules response
func TestTeoL7AccRulesRead_EmptyRules(t *testing.T) {
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
	d.SetId("zone-test-123")

	err := res.Read(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "rules is empty")
}

// TestTeoL7AccRulesRead_WithFilters tests Read with filters
func TestTeoL7AccRulesRead_WithFilters(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeL7AccRules", func(request *teov20220901.DescribeL7AccRulesRequest) (*teov20220901.DescribeL7AccRulesResponse, error) {
		assert.Len(t, request.Filters, 1)
		assert.Equal(t, "rule-id", *request.Filters[0].Name)
		assert.Len(t, request.Filters[0].Values, 1)
		assert.Equal(t, "rule-1", *request.Filters[0].Values[0])

		resp := teov20220901.NewDescribeL7AccRulesResponse()
		resp.Response = &teov20220901.DescribeL7AccRulesResponseParams{
			TotalCount: ptrInt64(1),
			Rules: []*teov20220901.RuleEngineItem{
				{
					RuleId:       ptrString("rule-1"),
					RuleName:     ptrString("rule-one"),
					Status:       ptrString("enable"),
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
		"filters": []interface{}{
			map[string]interface{}{
				"name": "rule-id",
				"values": []interface{}{
					"rule-1",
				},
			},
		},
	})
	d.SetId("zone-test-123")

	err := res.Read(d, meta)
	assert.NoError(t, err)
}

// TestTeoL7AccRulesUpdate_AddRule tests Update creates new rules
func TestTeoL7AccRulesUpdate_AddRule(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)

	createCalled := false
	patches.ApplyMethodFunc(teoClient, "CreateL7AccRules", func(request *teov20220901.CreateL7AccRulesRequest) (*teov20220901.CreateL7AccRulesResponse, error) {
		createCalled = true
		assert.Equal(t, "rule-new", *request.Rules[0].RuleName)
		resp := teov20220901.NewCreateL7AccRulesResponse()
		resp.Response = &teov20220901.CreateL7AccRulesResponseParams{
			RuleIds:   []*string{ptrString("rule-new-id")},
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(teoClient, "DescribeL7AccRules", func(request *teov20220901.DescribeL7AccRulesRequest) (*teov20220901.DescribeL7AccRulesResponse, error) {
		resp := teov20220901.NewDescribeL7AccRulesResponse()
		resp.Response = &teov20220901.DescribeL7AccRulesResponseParams{
			TotalCount: ptrInt64(2),
			Rules: []*teov20220901.RuleEngineItem{
				{
					RuleId:   ptrString("rule-existing-id"),
					RuleName: ptrString("rule-existing"),
					Status:   ptrString("enable"),
				},
				{
					RuleId:   ptrString("rule-new-id"),
					RuleName: ptrString("rule-new"),
					Status:   ptrString("enable"),
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
			// New rule set: has rule-existing (unchanged) and rule-new (added)
			map[string]interface{}{
				"rule_name": "rule-existing",
				"status":    "enable",
				"rule_id":   "rule-existing-id",
			},
			map[string]interface{}{
				"rule_name": "rule-new",
				"status":    "enable",
			},
		},
	})
	d.SetId("zone-test-123")

	// Store old state with just rule-existing
	_ = d.Set("rules", []interface{}{
		map[string]interface{}{
			"rule_name": "rule-existing",
			"status":    "enable",
			"rule_id":   "rule-existing-id",
		},
	})

	err := res.Update(d, meta)
	assert.NoError(t, err)
	assert.True(t, createCalled)
}

// TestTeoL7AccRulesUpdate_DeleteRule tests Update deletes removed rules
func TestTeoL7AccRulesUpdate_DeleteRule(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)

	deleteCalled := false
	patches.ApplyMethodFunc(teoClient, "DeleteL7AccRules", func(request *teov20220901.DeleteL7AccRulesRequest) (*teov20220901.DeleteL7AccRulesResponse, error) {
		deleteCalled = true
		assert.Equal(t, "rule-removed-id", *request.RuleIds[0])
		resp := teov20220901.NewDeleteL7AccRulesResponse()
		resp.Response = &teov20220901.DeleteL7AccRulesResponseParams{
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(teoClient, "DescribeL7AccRules", func(request *teov20220901.DescribeL7AccRulesRequest) (*teov20220901.DescribeL7AccRulesResponse, error) {
		resp := teov20220901.NewDescribeL7AccRulesResponse()
		resp.Response = &teov20220901.DescribeL7AccRulesResponseParams{
			TotalCount: ptrInt64(1),
			Rules: []*teov20220901.RuleEngineItem{
				{
					RuleId:   ptrString("rule-kept-id"),
					RuleName: ptrString("rule-kept"),
					Status:   ptrString("enable"),
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
			// New rule set: only rule-kept, rule-removed is gone
			map[string]interface{}{
				"rule_name": "rule-kept",
				"status":    "enable",
				"rule_id":   "rule-kept-id",
			},
		},
	})
	d.SetId("zone-test-123")

	// Store old state with both rules
	_ = d.Set("rules", []interface{}{
		map[string]interface{}{
			"rule_name": "rule-kept",
			"status":    "enable",
			"rule_id":   "rule-kept-id",
		},
		map[string]interface{}{
			"rule_name": "rule-removed",
			"status":    "enable",
			"rule_id":   "rule-removed-id",
		},
	})

	err := res.Update(d, meta)
	assert.NoError(t, err)
	assert.True(t, deleteCalled)
}

// TestTeoL7AccRulesUpdate_ModifyRule tests Update modifies existing rules
func TestTeoL7AccRulesUpdate_ModifyRule(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)

	modifyCalled := false
	patches.ApplyMethodFunc(teoClient, "ModifyL7AccRule", func(request *teov20220901.ModifyL7AccRuleRequest) (*teov20220901.ModifyL7AccRuleResponse, error) {
		modifyCalled = true
		assert.Equal(t, "rule-1-id", *request.Rule.RuleId)
		assert.Equal(t, "disable", *request.Rule.Status)
		resp := teov20220901.NewModifyL7AccRuleResponse()
		resp.Response = &teov20220901.ModifyL7AccRuleResponseParams{
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(teoClient, "DescribeL7AccRules", func(request *teov20220901.DescribeL7AccRulesRequest) (*teov20220901.DescribeL7AccRulesResponse, error) {
		resp := teov20220901.NewDescribeL7AccRulesResponse()
		resp.Response = &teov20220901.DescribeL7AccRulesResponseParams{
			TotalCount: ptrInt64(1),
			Rules: []*teov20220901.RuleEngineItem{
				{
					RuleId:   ptrString("rule-1-id"),
					RuleName: ptrString("rule-1"),
					Status:   ptrString("disable"),
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
			// New rule set: status changed from enable to disable
			map[string]interface{}{
				"rule_name": "rule-1",
				"status":    "disable",
				"rule_id":   "rule-1-id",
			},
		},
	})
	d.SetId("zone-test-123")

	// Store old state with status=enable
	_ = d.Set("rules", []interface{}{
		map[string]interface{}{
			"rule_name": "rule-1",
			"status":    "enable",
			"rule_id":   "rule-1-id",
		},
	})

	err := res.Update(d, meta)
	assert.NoError(t, err)
	assert.True(t, modifyCalled)
}

// TestTeoL7AccRulesDelete_Success tests Delete calls DeleteL7AccRules
func TestTeoL7AccRulesDelete_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)

	deleteCalled := false
	patches.ApplyMethodFunc(teoClient, "DeleteL7AccRules", func(request *teov20220901.DeleteL7AccRulesRequest) (*teov20220901.DeleteL7AccRulesResponse, error) {
		deleteCalled = true
		assert.Equal(t, "zone-test-123", *request.ZoneId)
		assert.Len(t, request.RuleIds, 2)
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
		"rule_ids": []interface{}{
			"rule-1",
			"rule-2",
		},
	})
	d.SetId("zone-test-123")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
	assert.True(t, deleteCalled)
	assert.Equal(t, "", d.Id())
}

// TestTeoL7AccRulesDelete_NoRules tests Delete with no rule_ids
func TestTeoL7AccRulesDelete_NoRules(t *testing.T) {
	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoL7AccRules()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":  "zone-test-123",
		"rule_ids": []interface{}{},
	})
	d.SetId("zone-test-123")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}

// TestTeoL7AccRulesImport tests Import by zone_id
func TestTeoL7AccRulesImport(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeL7AccRules", func(request *teov20220901.DescribeL7AccRulesRequest) (*teov20220901.DescribeL7AccRulesResponse, error) {
		assert.Equal(t, "zone-import-test", *request.ZoneId)
		resp := teov20220901.NewDescribeL7AccRulesResponse()
		resp.Response = &teov20220901.DescribeL7AccRulesResponseParams{
			TotalCount: ptrInt64(1),
			Rules: []*teov20220901.RuleEngineItem{
				{
					RuleId:       ptrString("rule-import-1"),
					RuleName:     ptrString("imported-rule"),
					Status:       ptrString("enable"),
					RulePriority: ptrInt64(1),
				},
			},
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoL7AccRules()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
	d.SetId("zone-import-test")

	// Simulate import
	importData, errs := res.Importer.State(d, meta)
	assert.Empty(t, errs)
	assert.Len(t, importData, 1)
	assert.Equal(t, "zone-import-test", importData[0].Id())

	err := res.Read(importData[0], meta)
	assert.NoError(t, err)

	rules := importData[0].Get("rules").([]interface{})
	assert.Len(t, rules, 1)
}

// Note: ptrString, ptrInt64, mockMeta, and newMockMeta are defined in other test files in this package
var _ tccommon.ProviderMeta = &mockMeta{}
var _ = connectivity.TencentCloudClient{}
