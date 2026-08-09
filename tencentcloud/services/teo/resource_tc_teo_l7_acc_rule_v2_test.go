package teo_test

import (
	"fmt"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/teo"
)

func TestAccTencentCloudTeoL7AccRuleV2Resource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PRIVATE) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTeoL7V2AccRule,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_teo_l7_acc_rule_v2.teo_l7_acc_rule_v2", "id"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule_v2.teo_l7_acc_rule_v2", "zone_id", "zone-3fkff38fyw8s"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule_v2.teo_l7_acc_rule_v2", "description.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule_v2.teo_l7_acc_rule_v2", "rule_name", "Web Acceleration 1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule_v2.teo_l7_acc_rule_v2", "status", "enable"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_l7_acc_rule_v2.teo_l7_acc_rule_v2", "rule_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_l7_acc_rule_v2.teo_l7_acc_rule_v2", "rule_priority"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule_v2.teo_l7_acc_rule_v2", "branches.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule_v2.teo_l7_acc_rule_v2", "branches.0.condition", "${http.request.host} in ['aaa.makn.cn']"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule_v2.teo_l7_acc_rule_v2", "branches.0.actions.#", "2"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule_v2.teo_l7_acc_rule_v2", "branches.0.actions.0.name", "Cache"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule_v2.teo_l7_acc_rule_v2", "branches.0.actions.0.cache_parameters.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule_v2.teo_l7_acc_rule_v2", "branches.0.actions.0.cache_parameters.0.custom_time.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule_v2.teo_l7_acc_rule_v2", "branches.0.actions.0.cache_parameters.0.custom_time.0.cache_time", "2592000"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule_v2.teo_l7_acc_rule_v2", "branches.0.actions.0.cache_parameters.0.custom_time.0.ignore_cache_control", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule_v2.teo_l7_acc_rule_v2", "branches.0.actions.0.cache_parameters.0.custom_time.0.switch", "on"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule_v2.teo_l7_acc_rule_v2", "branches.0.actions.1.name", "CacheKey"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule_v2.teo_l7_acc_rule_v2", "branches.0.actions.1.cache_key_parameters.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule_v2.teo_l7_acc_rule_v2", "branches.0.actions.1.cache_key_parameters.0.full_url_cache", "on"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule_v2.teo_l7_acc_rule_v2", "branches.0.actions.1.cache_key_parameters.0.ignore_case", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule_v2.teo_l7_acc_rule_v2", "branches.0.actions.1.cache_key_parameters.0.query_string.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule_v2.teo_l7_acc_rule_v2", "branches.0.actions.1.cache_key_parameters.0.query_string.0.switch", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule_v2.teo_l7_acc_rule_v2", "branches.0.actions.1.cache_key_parameters.0.query_string.0.values.#", "0"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule_v2.teo_l7_acc_rule_v2", "branches.0.sub_rules.#", "2"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule_v2.teo_l7_acc_rule_v2", "branches.0.sub_rules.0.description.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule_v2.teo_l7_acc_rule_v2", "branches.0.sub_rules.0.branches.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule_v2.teo_l7_acc_rule_v2", "branches.0.sub_rules.0.branches.0.condition", "lower(${http.request.file_extension}) in ['php', 'jsp', 'asp', 'aspx']"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule_v2.teo_l7_acc_rule_v2", "branches.0.sub_rules.0.branches.0.actions.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule_v2.teo_l7_acc_rule_v2", "branches.0.sub_rules.0.branches.0.actions.0.name", "Cache"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule_v2.teo_l7_acc_rule_v2", "branches.0.sub_rules.0.branches.0.actions.0.cache_parameters.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule_v2.teo_l7_acc_rule_v2", "branches.0.sub_rules.0.branches.0.actions.0.cache_parameters.0.no_cache.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule_v2.teo_l7_acc_rule_v2", "branches.0.sub_rules.0.branches.0.actions.0.cache_parameters.0.no_cache.0.switch", "on"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule_v2.teo_l7_acc_rule_v2", "branches.0.sub_rules.1.description.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule_v2.teo_l7_acc_rule_v2", "branches.0.sub_rules.1.branches.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule_v2.teo_l7_acc_rule_v2", "branches.0.sub_rules.1.branches.0.condition", "${http.request.file_extension} in ['jpg', 'png', 'gif', 'bmp', 'svg', 'webp']"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule_v2.teo_l7_acc_rule_v2", "branches.0.sub_rules.1.branches.0.actions.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule_v2.teo_l7_acc_rule_v2", "branches.0.sub_rules.1.branches.0.actions.0.name", "MaxAge"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule_v2.teo_l7_acc_rule_v2", "branches.0.sub_rules.1.branches.0.actions.0.max_age_parameters.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule_v2.teo_l7_acc_rule_v2", "branches.0.sub_rules.1.branches.0.actions.0.max_age_parameters.0.cache_time", "3600"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule_v2.teo_l7_acc_rule_v2", "branches.0.sub_rules.1.branches.0.actions.0.max_age_parameters.0.follow_origin", "off"),
				),
			},
			{
				Config: testAccTeoL7V2AccRuleUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule_v2.teo_l7_acc_rule_v2", "description.0", "2"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule_v2.teo_l7_acc_rule_v2", "rule_name", "Web Acceleration 2"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule_v2.teo_l7_acc_rule_v2", "branches.0.actions.#", "2"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule_v2.teo_l7_acc_rule_v2", "branches.0.sub_rules.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule_v2.teo_l7_acc_rule_v2", "branches.0.actions.1.name", "OriginPullProtocol"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule_v2.teo_l7_acc_rule_v2", "branches.0.actions.1.origin_pull_protocol_parameters.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l7_acc_rule_v2.teo_l7_acc_rule_v2", "branches.0.actions.1.origin_pull_protocol_parameters.0.protocol", "https"),
				),
			},
			{
				ResourceName:      "tencentcloud_teo_l7_acc_rule_v2.teo_l7_acc_rule_v2",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTeoL7V2AccRule = `
resource "tencentcloud_teo_l7_acc_rule_v2" "teo_l7_acc_rule_v2" {
  zone_id     = "zone-3fkff38fyw8s"
  description = ["1"]
  rule_name   = "Web Acceleration 1"
  status = "enable"
  branches {
    condition = "$${http.request.host} in ['aaa.makn.cn']"
    actions {
      name = "Cache"
      cache_parameters {
        custom_time {
          cache_time           = 2592000
          ignore_cache_control = "off"
          switch               = "on"
        }
      }
    }

    actions {
      name = "CacheKey"
      cache_key_parameters {
        full_url_cache = "on"
        ignore_case    = "off"
        query_string {
          switch = "off"
          values = []
        }
      }
    }

    sub_rules {
      description = ["1-1"]
      branches {
        condition = "lower($${http.request.file_extension}) in ['php', 'jsp', 'asp', 'aspx']"
        actions {
          name = "Cache"
          cache_parameters {
            no_cache {
              switch = "on"
            }
          }
        }
      }
    }

    sub_rules {
      description = ["1-2"]
      branches {
        condition = "$${http.request.file_extension} in ['jpg', 'png', 'gif', 'bmp', 'svg', 'webp']"
        actions {
          name = "MaxAge"
          max_age_parameters {
            cache_time    = 3600
            follow_origin = "off"
          }
        }
      }
    }
  }
}
`

const testAccTeoL7V2AccRuleUpdate = `
resource "tencentcloud_teo_l7_acc_rule_v2" "teo_l7_acc_rule_v2" {
  zone_id     = "zone-3fkff38fyw8s"
  description = ["2"]
  rule_name   = "Web Acceleration 2"
  status = "enable"
  branches {
    condition = "$${http.request.host} in ['aaa.makn.cn']"
    actions {
      name = "Cache"
      cache_parameters {
        custom_time {
          cache_time           = 2592000
          ignore_cache_control = "off"
          switch               = "on"
        }
      }
    }
    actions {
      name = "OriginPullProtocol"
      origin_pull_protocol_parameters {
          protocol = "https"
      }
    }

    sub_rules {
      description = ["01-1"]
      branches {
        condition = "lower($${http.request.file_extension}) in ['php', 'jsp', 'asp', 'aspx']"
        actions {
          name = "Cache"
          cache_parameters {
            no_cache {
              switch = "on"
            }
          }
        }
      }
    }

  }
}
`

// mockMetaL7AccRuleV2 implements tccommon.ProviderMeta for unit tests
type mockMetaL7AccRuleV2 struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaL7AccRuleV2) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaL7AccRuleV2{}

func newMockMetaL7AccRuleV2() *mockMetaL7AccRuleV2 {
	return &mockMetaL7AccRuleV2{client: &connectivity.TencentCloudClient{}}
}

func ptrStringL7(s string) *string {
	return &s
}

func ptrInt64L7(i int64) *int64 {
	return &i
}

// TestTeoL7AccRuleV2_ReadWithFilters tests the Read function with filters specified
// go test ./tencentcloud/services/teo/ -run "TestTeoL7AccRuleV2_ReadWithFilters" -v -count=1 -gcflags="all=-l"
func TestTeoL7AccRuleV2_ReadWithFilters(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaL7AccRuleV2().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeL7AccRules", func(request *teov20220901.DescribeL7AccRulesRequest) (*teov20220901.DescribeL7AccRulesResponse, error) {
		resp := teov20220901.NewDescribeL7AccRulesResponse()
		resp.Response = &teov20220901.DescribeL7AccRulesResponseParams{
			TotalCount: ptrInt64L7(2),
			Rules: []*teov20220901.RuleEngineItem{
				{
					Status:       ptrStringL7("enable"),
					RuleId:       ptrStringL7("rule-abc123"),
					RuleName:     ptrStringL7("Test Rule 1"),
					Description:  []*string{ptrStringL7("desc1")},
					RulePriority: ptrInt64L7(10),
				},
				{
					Status:       ptrStringL7("disable"),
					RuleId:       ptrStringL7("rule-def456"),
					RuleName:     ptrStringL7("Test Rule 2"),
					Description:  []*string{ptrStringL7("desc2")},
					RulePriority: ptrInt64L7(20),
				},
			},
			RequestId: ptrStringL7("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaL7AccRuleV2()
	res := teo.ResourceTencentCloudTeoL7AccRuleV2()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-test123",
		"filters": []interface{}{
			map[string]interface{}{
				"name":   "rule-id",
				"values": []interface{}{"rule-abc123"},
			},
		},
	})
	d.SetId("zone-test123#rule-abc123")

	err := res.Read(d, meta)
	assert.NoError(t, err)

	// Verify rules computed field was populated
	rules := d.Get("rules").([]interface{})
	assert.Len(t, rules, 2)

	rule0 := rules[0].(map[string]interface{})
	assert.Equal(t, "enable", rule0["status"])
	assert.Equal(t, "rule-abc123", rule0["rule_id"])
	assert.Equal(t, "Test Rule 1", rule0["rule_name"])
	assert.Equal(t, 10, rule0["rule_priority"])

	rule1 := rules[1].(map[string]interface{})
	assert.Equal(t, "disable", rule1["status"])
	assert.Equal(t, "rule-def456", rule1["rule_id"])
	assert.Equal(t, "Test Rule 2", rule1["rule_name"])
	assert.Equal(t, 20, rule1["rule_priority"])
}

// TestTeoL7AccRuleV2_ReadWithoutFilters tests the Read function without filters (backward compatibility)
// go test ./tencentcloud/services/teo/ -run "TestTeoL7AccRuleV2_ReadWithoutFilters" -v -count=1 -gcflags="all=-l"
func TestTeoL7AccRuleV2_ReadWithoutFilters(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaL7AccRuleV2().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeL7AccRules", func(request *teov20220901.DescribeL7AccRulesRequest) (*teov20220901.DescribeL7AccRulesResponse, error) {
		resp := teov20220901.NewDescribeL7AccRulesResponse()
		resp.Response = &teov20220901.DescribeL7AccRulesResponseParams{
			TotalCount: ptrInt64L7(1),
			Rules: []*teov20220901.RuleEngineItem{
				{
					Status:       ptrStringL7("enable"),
					RuleId:       ptrStringL7("rule-abc123"),
					RuleName:     ptrStringL7("Test Rule 1"),
					Description:  []*string{ptrStringL7("desc1")},
					RulePriority: ptrInt64L7(10),
				},
			},
			RequestId: ptrStringL7("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaL7AccRuleV2()
	res := teo.ResourceTencentCloudTeoL7AccRuleV2()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-test123",
	})
	d.SetId("zone-test123#rule-abc123")

	err := res.Read(d, meta)
	assert.NoError(t, err)

	// Verify existing behavior still works (single rule read)
	assert.Equal(t, "enable", d.Get("status"))
	assert.Equal(t, "Test Rule 1", d.Get("rule_name"))
	assert.Equal(t, 10, d.Get("rule_priority"))

	// Verify rules computed field is not populated when filters not specified
	rules := d.Get("rules").([]interface{})
	assert.Len(t, rules, 0)
}

// TestTeoL7AccRuleV2_ReadWithFiltersEmptyResponse tests Read with filters returning empty rules
// go test ./tencentcloud/services/teo/ -run "TestTeoL7AccRuleV2_ReadWithFiltersEmptyResponse" -v -count=1 -gcflags="all=-l"
func TestTeoL7AccRuleV2_ReadWithFiltersEmptyResponse(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaL7AccRuleV2().client, "UseTeoV20220901Client", teoClient)

	callCount := 0
	patches.ApplyMethodFunc(teoClient, "DescribeL7AccRules", func(request *teov20220901.DescribeL7AccRulesRequest) (*teov20220901.DescribeL7AccRulesResponse, error) {
		callCount++
		resp := teov20220901.NewDescribeL7AccRulesResponse()
		if callCount == 1 {
			resp.Response = &teov20220901.DescribeL7AccRulesResponseParams{
				TotalCount: ptrInt64L7(1),
				Rules: []*teov20220901.RuleEngineItem{
					{
						Status:       ptrStringL7("enable"),
						RuleId:       ptrStringL7("rule-abc123"),
						RuleName:     ptrStringL7("Test Rule"),
						RulePriority: ptrInt64L7(10),
					},
				},
				RequestId: ptrStringL7("fake-request-id"),
			}
		} else {
			resp.Response = &teov20220901.DescribeL7AccRulesResponseParams{
				TotalCount: ptrInt64L7(0),
				Rules:      []*teov20220901.RuleEngineItem{},
				RequestId:  ptrStringL7("fake-request-id"),
			}
		}
		return resp, nil
	})

	meta := newMockMetaL7AccRuleV2()
	res := teo.ResourceTencentCloudTeoL7AccRuleV2()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-test123",
		"filters": []interface{}{
			map[string]interface{}{
				"name":   "rule-id",
				"values": []interface{}{"rule-nonexistent"},
			},
		},
	})
	d.SetId("zone-test123#rule-abc123")

	err := res.Read(d, meta)
	assert.NoError(t, err)

	// Verify rules computed field is empty list
	rules := d.Get("rules").([]interface{})
	assert.Len(t, rules, 0)
}

// TestTeoL7AccRuleV2_ReadWithFiltersAPIError tests Read with API error on filter query
// go test ./tencentcloud/services/teo/ -run "TestTeoL7AccRuleV2_ReadWithFiltersAPIError" -v -count=1 -gcflags="all=-l"
func TestTeoL7AccRuleV2_ReadWithFiltersAPIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaL7AccRuleV2().client, "UseTeoV20220901Client", teoClient)

	callCount := 0
	patches.ApplyMethodFunc(teoClient, "DescribeL7AccRules", func(request *teov20220901.DescribeL7AccRulesRequest) (*teov20220901.DescribeL7AccRulesResponse, error) {
		callCount++
		if callCount == 1 {
			resp := teov20220901.NewDescribeL7AccRulesResponse()
			resp.Response = &teov20220901.DescribeL7AccRulesResponseParams{
				TotalCount: ptrInt64L7(1),
				Rules: []*teov20220901.RuleEngineItem{
					{
						Status:       ptrStringL7("enable"),
						RuleId:       ptrStringL7("rule-abc123"),
						RuleName:     ptrStringL7("Test Rule"),
						RulePriority: ptrInt64L7(10),
					},
				},
				RequestId: ptrStringL7("fake-request-id"),
			}
			return resp, nil
		}
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=InternalError, Message=internal error")
	})

	meta := newMockMetaL7AccRuleV2()
	res := teo.ResourceTencentCloudTeoL7AccRuleV2()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-test123",
		"filters": []interface{}{
			map[string]interface{}{
				"name":   "rule-id",
				"values": []interface{}{"rule-abc123"},
			},
		},
	})
	d.SetId("zone-test123#rule-abc123")

	err := res.Read(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "InternalError")
}

// TestTeoL7AccRuleV2_Schema validates the new filters and rules schema fields
// go test ./tencentcloud/services/teo/ -run "TestTeoL7AccRuleV2_Schema" -v -count=1 -gcflags="all=-l"
func TestTeoL7AccRuleV2_Schema(t *testing.T) {
	res := teo.ResourceTencentCloudTeoL7AccRuleV2()

	assert.NotNil(t, res)

	// Check filters schema
	filtersSchema, ok := res.Schema["filters"]
	assert.True(t, ok, "filters field should exist in schema")
	assert.Equal(t, schema.TypeList, filtersSchema.Type)
	assert.True(t, filtersSchema.Optional)
	assert.False(t, filtersSchema.Computed)

	// Check rules schema
	rulesSchema, ok := res.Schema["rules"]
	assert.True(t, ok, "rules field should exist in schema")
	assert.Equal(t, schema.TypeList, rulesSchema.Type)
	assert.True(t, rulesSchema.Computed)
	assert.False(t, rulesSchema.Optional)

	// Check existing fields still exist
	assert.Contains(t, res.Schema, "zone_id")
	assert.Contains(t, res.Schema, "status")
	assert.Contains(t, res.Schema, "rule_name")
	assert.Contains(t, res.Schema, "description")
	assert.Contains(t, res.Schema, "branches")
	assert.Contains(t, res.Schema, "rule_id")
	assert.Contains(t, res.Schema, "rule_priority")
}
