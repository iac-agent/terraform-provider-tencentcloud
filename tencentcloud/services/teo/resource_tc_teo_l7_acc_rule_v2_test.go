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

// mockMetaL7AccRule implements tccommon.ProviderMeta
type mockMetaL7AccRule struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaL7AccRule) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaL7AccRule{}

func newMockMetaL7AccRule() *mockMetaL7AccRule {
	return &mockMetaL7AccRule{client: &connectivity.TencentCloudClient{}}
}

func ptrStringL7(s string) *string { return &s }
func ptrInt64L7(i int64) *int64    { return &i }

// go test ./tencentcloud/services/teo/ -run "TestTeoL7AccRuleV2_Read_WithFilters" -v -count=1 -gcflags="all=-l"
func TestTeoL7AccRuleV2_Read_WithFilters(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaL7AccRule().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeL7AccRules", func(request *teov20220901.DescribeL7AccRulesRequest) (*teov20220901.DescribeL7AccRulesResponse, error) {
		// Verify filters are passed correctly
		if len(request.Filters) > 0 {
			assert.Equal(t, "rule-id", *request.Filters[0].Name)
			assert.Equal(t, []*string{ptrStringL7("rule-test123")}, request.Filters[0].Values)
		}
		resp := teov20220901.NewDescribeL7AccRulesResponse()
		resp.Response = &teov20220901.DescribeL7AccRulesResponseParams{
			TotalCount: ptrInt64L7(1),
			Rules: []*teov20220901.RuleEngineItem{
				{
					Status:       ptrStringL7("enable"),
					RuleId:       ptrStringL7("rule-test123"),
					RuleName:     ptrStringL7("Test Rule"),
					Description:  []*string{ptrStringL7("test description")},
					RulePriority: ptrInt64L7(10),
					Branches:     []*teov20220901.RuleBranch{},
				},
			},
			RequestId: ptrStringL7("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaL7AccRule()
	res := teo.ResourceTencentCloudTeoL7AccRuleV2()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-3fkff38fyw8s",
		"filters": []interface{}{
			map[string]interface{}{
				"name":   "rule-id",
				"values": []interface{}{"rule-test123"},
			},
		},
	})
	d.SetId("zone-3fkff38fyw8s#rule-test123")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "enable", d.Get("status"))
	assert.Equal(t, "Test Rule", d.Get("rule_name"))

	// Verify rules computed parameter is populated
	rules := d.Get("rules").([]interface{})
	assert.Equal(t, 1, len(rules))
	ruleMap := rules[0].(map[string]interface{})
	assert.Equal(t, "enable", ruleMap["status"])
	assert.Equal(t, "rule-test123", ruleMap["rule_id"])
	assert.Equal(t, "Test Rule", ruleMap["rule_name"])
}

// go test ./tencentcloud/services/teo/ -run "TestTeoL7AccRuleV2_Read_WithoutFilters" -v -count=1 -gcflags="all=-l"
func TestTeoL7AccRuleV2_Read_WithoutFilters(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaL7AccRule().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeL7AccRules", func(request *teov20220901.DescribeL7AccRulesRequest) (*teov20220901.DescribeL7AccRulesResponse, error) {
		// Verify default filter by rule-id is used when no filters specified
		assert.Equal(t, 1, len(request.Filters))
		assert.Equal(t, "rule-id", *request.Filters[0].Name)
		assert.Equal(t, []*string{ptrStringL7("rule-abc")}, request.Filters[0].Values)

		resp := teov20220901.NewDescribeL7AccRulesResponse()
		resp.Response = &teov20220901.DescribeL7AccRulesResponseParams{
			TotalCount: ptrInt64L7(1),
			Rules: []*teov20220901.RuleEngineItem{
				{
					Status:       ptrStringL7("enable"),
					RuleId:       ptrStringL7("rule-abc"),
					RuleName:     ptrStringL7("Default Rule"),
					RulePriority: ptrInt64L7(5),
					Branches:     []*teov20220901.RuleBranch{},
				},
			},
			RequestId: ptrStringL7("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaL7AccRule()
	res := teo.ResourceTencentCloudTeoL7AccRuleV2()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-3fkff38fyw8s",
	})
	d.SetId("zone-3fkff38fyw8s#rule-abc")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "enable", d.Get("status"))
	assert.Equal(t, "Default Rule", d.Get("rule_name"))

	// Verify rules computed parameter is populated
	rules := d.Get("rules").([]interface{})
	assert.Equal(t, 1, len(rules))
}

// go test ./tencentcloud/services/teo/ -run "TestTeoL7AccRuleV2_Read_MultipleRulesOutput" -v -count=1 -gcflags="all=-l"
func TestTeoL7AccRuleV2_Read_MultipleRulesOutput(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaL7AccRule().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeL7AccRules", func(request *teov20220901.DescribeL7AccRulesRequest) (*teov20220901.DescribeL7AccRulesResponse, error) {
		resp := teov20220901.NewDescribeL7AccRulesResponse()
		resp.Response = &teov20220901.DescribeL7AccRulesResponseParams{
			TotalCount: ptrInt64L7(2),
			Rules: []*teov20220901.RuleEngineItem{
				{
					Status:       ptrStringL7("enable"),
					RuleId:       ptrStringL7("rule-001"),
					RuleName:     ptrStringL7("Rule One"),
					Description:  []*string{ptrStringL7("first rule")},
					RulePriority: ptrInt64L7(10),
					Branches:     []*teov20220901.RuleBranch{},
				},
				{
					Status:       ptrStringL7("disable"),
					RuleId:       ptrStringL7("rule-002"),
					RuleName:     ptrStringL7("Rule Two"),
					Description:  []*string{ptrStringL7("second rule")},
					RulePriority: ptrInt64L7(20),
					Branches:     []*teov20220901.RuleBranch{},
				},
			},
			RequestId: ptrStringL7("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaL7AccRule()
	res := teo.ResourceTencentCloudTeoL7AccRuleV2()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-3fkff38fyw8s",
	})
	d.SetId("zone-3fkff38fyw8s#rule-001")

	err := res.Read(d, meta)
	assert.NoError(t, err)

	// Verify rules computed parameter contains both rules
	rules := d.Get("rules").([]interface{})
	assert.Equal(t, 2, len(rules))

	rule0 := rules[0].(map[string]interface{})
	assert.Equal(t, "enable", rule0["status"])
	assert.Equal(t, "rule-001", rule0["rule_id"])
	assert.Equal(t, "Rule One", rule0["rule_name"])

	rule1 := rules[1].(map[string]interface{})
	assert.Equal(t, "disable", rule1["status"])
	assert.Equal(t, "rule-002", rule1["rule_id"])
	assert.Equal(t, "Rule Two", rule1["rule_name"])
}

// go test ./tencentcloud/services/teo/ -run "TestTeoL7AccRuleV2_Read_EmptyRulesOutput" -v -count=1 -gcflags="all=-l"
func TestTeoL7AccRuleV2_Read_EmptyRulesOutput(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaL7AccRule().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeL7AccRules", func(request *teov20220901.DescribeL7AccRulesRequest) (*teov20220901.DescribeL7AccRulesResponse, error) {
		resp := teov20220901.NewDescribeL7AccRulesResponse()
		resp.Response = &teov20220901.DescribeL7AccRulesResponseParams{
			TotalCount: ptrInt64L7(0),
			Rules:      []*teov20220901.RuleEngineItem{},
			RequestId:  ptrStringL7("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaL7AccRule()
	res := teo.ResourceTencentCloudTeoL7AccRuleV2()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-3fkff38fyw8s",
	})
	d.SetId("zone-3fkff38fyw8s#rule-nonexistent")

	err := res.Read(d, meta)
	assert.NoError(t, err)

	// Resource should be marked as removed when no rules found
	assert.Equal(t, "", d.Id())
}

// go test ./tencentcloud/services/teo/ -run "TestTeoL7AccRuleV2_Read_APIError" -v -count=1 -gcflags="all=-l"
func TestTeoL7AccRuleV2_Read_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaL7AccRule().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeL7AccRules", func(request *teov20220901.DescribeL7AccRulesRequest) (*teov20220901.DescribeL7AccRulesResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=ResourceNotFound, Message=Rule not found")
	})

	meta := newMockMetaL7AccRule()
	res := teo.ResourceTencentCloudTeoL7AccRuleV2()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-3fkff38fyw8s",
	})
	d.SetId("zone-3fkff38fyw8s#rule-nonexistent")

	err := res.Read(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ResourceNotFound")
}

// go test ./tencentcloud/services/teo/ -run "TestTeoL7AccRuleV2_MultipleFilters" -v -count=1 -gcflags="all=-l"
func TestTeoL7AccRuleV2_MultipleFilters(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaL7AccRule().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeL7AccRules", func(request *teov20220901.DescribeL7AccRulesRequest) (*teov20220901.DescribeL7AccRulesResponse, error) {
		// Verify multiple filters are passed
		assert.Equal(t, 2, len(request.Filters))
		assert.Equal(t, "rule-id", *request.Filters[0].Name)

		resp := teov20220901.NewDescribeL7AccRulesResponse()
		resp.Response = &teov20220901.DescribeL7AccRulesResponseParams{
			TotalCount: ptrInt64L7(1),
			Rules: []*teov20220901.RuleEngineItem{
				{
					Status:       ptrStringL7("enable"),
					RuleId:       ptrStringL7("rule-multi"),
					RuleName:     ptrStringL7("Multi Filter Rule"),
					RulePriority: ptrInt64L7(15),
					Branches:     []*teov20220901.RuleBranch{},
				},
			},
			RequestId: ptrStringL7("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaL7AccRule()
	res := teo.ResourceTencentCloudTeoL7AccRuleV2()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-3fkff38fyw8s",
		"filters": []interface{}{
			map[string]interface{}{
				"name":   "rule-id",
				"values": []interface{}{"rule-multi"},
			},
			map[string]interface{}{
				"name":   "rule-id",
				"values": []interface{}{"rule-other"},
			},
		},
	})
	d.SetId("zone-3fkff38fyw8s#rule-multi")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "enable", d.Get("status"))
	assert.Equal(t, "Multi Filter Rule", d.Get("rule_name"))
}

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
