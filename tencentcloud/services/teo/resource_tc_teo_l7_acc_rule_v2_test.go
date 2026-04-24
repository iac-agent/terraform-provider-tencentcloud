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

// go test ./tencentcloud/services/teo/ -run "TestTeoL7AccRuleV2" -v -count=1 -gcflags="all=-l"

// TestTeoL7AccRuleV2_Create_WithRuleBlock tests Create calls API with rule block and sets ID
func TestTeoL7AccRuleV2_Create_WithRuleBlock(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "CreateL7AccRules", func(request *teov20220901.CreateL7AccRulesRequest) (*teov20220901.CreateL7AccRulesResponse, error) {
		assert.Equal(t, "zone-1234567890", *request.ZoneId)
		assert.Equal(t, 1, len(request.Rules))
		rule := request.Rules[0]
		assert.Equal(t, "enable", *rule.Status)
		assert.Equal(t, "test-rule", *rule.RuleName)

		resp := teov20220901.NewCreateL7AccRulesResponse()
		resp.Response = &teov20220901.CreateL7AccRulesResponseParams{
			RuleIds:   []*string{ptrString("rule-abcdefghij")},
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
					Status:       ptrString("enable"),
					RuleId:       ptrString("rule-abcdefghij"),
					RuleName:     ptrString("test-rule"),
					Description:  []*string{ptrString("test-desc")},
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
	assert.Equal(t, "zone-1234567890#rule-abcdefghij", d.Id())
}

// TestTeoL7AccRuleV2_Create_WithTopLevelFields tests Create with top-level fields (backward compatibility)
func TestTeoL7AccRuleV2_Create_WithTopLevelFields(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "CreateL7AccRules", func(request *teov20220901.CreateL7AccRulesRequest) (*teov20220901.CreateL7AccRulesResponse, error) {
		assert.Equal(t, "zone-1234567890", *request.ZoneId)
		assert.Equal(t, 1, len(request.Rules))
		rule := request.Rules[0]
		assert.Equal(t, "enable", *rule.Status)
		assert.Equal(t, "test-rule", *rule.RuleName)

		resp := teov20220901.NewCreateL7AccRulesResponse()
		resp.Response = &teov20220901.CreateL7AccRulesResponseParams{
			RuleIds:   []*string{ptrString("rule-abcdefghij")},
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
					Status:       ptrString("enable"),
					RuleId:       ptrString("rule-abcdefghij"),
					RuleName:     ptrString("test-rule"),
					Description:  []*string{ptrString("test-desc")},
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
		"zone_id":     "zone-1234567890",
		"status":      "enable",
		"rule_name":   "test-rule",
		"description": []interface{}{"test-desc"},
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-1234567890#rule-abcdefghij", d.Id())
}

// TestTeoL7AccRuleV2_Create_APIError tests Create handles API error
func TestTeoL7AccRuleV2_Create_APIError(t *testing.T) {
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

// TestTeoL7AccRuleV2_Update_WithRuleBlock tests Update with rule block changes
func TestTeoL7AccRuleV2_Update_WithRuleBlock(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "ModifyL7AccRule", func(request *teov20220901.ModifyL7AccRuleRequest) (*teov20220901.ModifyL7AccRuleResponse, error) {
		assert.Equal(t, "zone-1234567890", *request.ZoneId)
		assert.NotNil(t, request.Rule)
		assert.Equal(t, "rule-abcdefghij", *request.Rule.RuleId)
		assert.Equal(t, "disable", *request.Rule.Status)
		assert.Equal(t, "updated-rule", *request.Rule.RuleName)

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
					Status:       ptrString("disable"),
					RuleId:       ptrString("rule-abcdefghij"),
					RuleName:     ptrString("updated-rule"),
					Description:  []*string{ptrString("updated-desc")},
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
		"rule": []interface{}{
			map[string]interface{}{
				"status":      "disable",
				"rule_name":   "updated-rule",
				"description": []interface{}{"updated-desc"},
			},
		},
		"rule_id": "rule-abcdefghij",
	})
	d.SetId("zone-1234567890#rule-abcdefghij")

	err := res.Update(d, meta)
	assert.NoError(t, err)
}

// TestTeoL7AccRuleV2_Update_WithTopLevelFields tests Update with top-level fields (backward compatibility)
func TestTeoL7AccRuleV2_Update_WithTopLevelFields(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "ModifyL7AccRule", func(request *teov20220901.ModifyL7AccRuleRequest) (*teov20220901.ModifyL7AccRuleResponse, error) {
		assert.Equal(t, "zone-1234567890", *request.ZoneId)
		assert.NotNil(t, request.Rule)
		assert.Equal(t, "rule-abcdefghij", *request.Rule.RuleId)
		assert.Equal(t, "disable", *request.Rule.Status)
		assert.Equal(t, "updated-rule", *request.Rule.RuleName)

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
					Status:       ptrString("disable"),
					RuleId:       ptrString("rule-abcdefghij"),
					RuleName:     ptrString("updated-rule"),
					Description:  []*string{ptrString("updated-desc")},
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
		"zone_id":     "zone-1234567890",
		"status":      "disable",
		"rule_name":   "updated-rule",
		"description": []interface{}{"updated-desc"},
		"rule_id":     "rule-abcdefghij",
	})
	d.SetId("zone-1234567890#rule-abcdefghij")

	err := res.Update(d, meta)
	assert.NoError(t, err)
}

// TestTeoL7AccRuleV2_Read_WithRuleBlock tests Read populates rule block
func TestTeoL7AccRuleV2_Read_WithRuleBlock(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeL7AccRules", func(request *teov20220901.DescribeL7AccRulesRequest) (*teov20220901.DescribeL7AccRulesResponse, error) {
		resp := teov20220901.NewDescribeL7AccRulesResponse()
		resp.Response = &teov20220901.DescribeL7AccRulesResponseParams{
			TotalCount: ptrInt64(1),
			Rules: []*teov20220901.RuleEngineItem{
				{
					Status:       ptrString("enable"),
					RuleId:       ptrString("rule-abcdefghij"),
					RuleName:     ptrString("test-rule"),
					Description:  []*string{ptrString("test-desc")},
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
		"rule": []interface{}{
			map[string]interface{}{
				"status":    "enable",
				"rule_name": "test-rule",
			},
		},
	})
	d.SetId("zone-1234567890#rule-abcdefghij")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "enable", d.Get("status"))
	assert.Equal(t, "test-rule", d.Get("rule_name"))

	// Verify rule block is populated
	ruleBlock := d.Get("rule").([]interface{})
	assert.Equal(t, 1, len(ruleBlock))
	ruleMap := ruleBlock[0].(map[string]interface{})
	assert.Equal(t, "enable", ruleMap["status"])
	assert.Equal(t, "test-rule", ruleMap["rule_name"])
}

// TestTeoL7AccRuleV2_Schema validates the schema definition including the new rule parameter
func TestTeoL7AccRuleV2_Schema(t *testing.T) {
	res := teo.ResourceTencentCloudTeoL7AccRuleV2()

	assert.NotNil(t, res)
	assert.NotNil(t, res.Create)
	assert.NotNil(t, res.Read)
	assert.NotNil(t, res.Update)
	assert.NotNil(t, res.Delete)
	assert.NotNil(t, res.Importer)

	// Check required fields
	assert.Contains(t, res.Schema, "zone_id")
	zoneId := res.Schema["zone_id"]
	assert.Equal(t, schema.TypeString, zoneId.Type)
	assert.True(t, zoneId.Required)
	assert.True(t, zoneId.ForceNew)

	// Check optional top-level fields
	assert.Contains(t, res.Schema, "status")
	status := res.Schema["status"]
	assert.Equal(t, schema.TypeString, status.Type)
	assert.True(t, status.Optional)

	assert.Contains(t, res.Schema, "rule_name")
	ruleName := res.Schema["rule_name"]
	assert.Equal(t, schema.TypeString, ruleName.Type)
	assert.True(t, ruleName.Optional)

	assert.Contains(t, res.Schema, "description")
	description := res.Schema["description"]
	assert.Equal(t, schema.TypeList, description.Type)
	assert.True(t, description.Optional)

	assert.Contains(t, res.Schema, "branches")
	branches := res.Schema["branches"]
	assert.Equal(t, schema.TypeList, branches.Type)
	assert.True(t, branches.Optional)

	// Check new rule parameter
	assert.Contains(t, res.Schema, "rule")
	ruleParam := res.Schema["rule"]
	assert.Equal(t, schema.TypeList, ruleParam.Type)
	assert.True(t, ruleParam.Optional)
	assert.Equal(t, 1, ruleParam.MaxItems)

	// Check rule block nested schema
	ruleElem := ruleParam.Elem.(*schema.Resource)
	assert.Contains(t, ruleElem.Schema, "status")
	assert.Contains(t, ruleElem.Schema, "rule_name")
	assert.Contains(t, ruleElem.Schema, "description")
	assert.Contains(t, ruleElem.Schema, "branches")

	ruleStatus := ruleElem.Schema["status"]
	assert.Equal(t, schema.TypeString, ruleStatus.Type)
	assert.True(t, ruleStatus.Optional)

	ruleRuleName := ruleElem.Schema["rule_name"]
	assert.Equal(t, schema.TypeString, ruleRuleName.Type)
	assert.True(t, ruleRuleName.Optional)

	ruleDescription := ruleElem.Schema["description"]
	assert.Equal(t, schema.TypeList, ruleDescription.Type)
	assert.True(t, ruleDescription.Optional)

	ruleBranches := ruleElem.Schema["branches"]
	assert.Equal(t, schema.TypeList, ruleBranches.Type)
	assert.True(t, ruleBranches.Optional)

	// Check computed fields
	assert.Contains(t, res.Schema, "rule_id")
	ruleId := res.Schema["rule_id"]
	assert.Equal(t, schema.TypeString, ruleId.Type)
	assert.True(t, ruleId.Computed)

	assert.Contains(t, res.Schema, "rule_priority")
	rulePriority := res.Schema["rule_priority"]
	assert.Equal(t, schema.TypeInt, rulePriority.Type)
	assert.True(t, rulePriority.Computed)
}

// TestTeoL7AccRuleV2_ExpandRuleBlock tests the expandTeoL7AccRuleV2RuleBlock helper function
func TestTeoL7AccRuleV2_ExpandRuleBlock(t *testing.T) {
	t.Run("empty list returns empty rule", func(t *testing.T) {
		rule := teo.ExpandTeoL7AccRuleV2RuleBlock([]interface{}{})
		assert.NotNil(t, rule)
		assert.Nil(t, rule.Status)
		assert.Nil(t, rule.RuleName)
		assert.Nil(t, rule.Description)
		assert.Nil(t, rule.Branches)
	})

	t.Run("nil list returns empty rule", func(t *testing.T) {
		rule := teo.ExpandTeoL7AccRuleV2RuleBlock(nil)
		assert.NotNil(t, rule)
		assert.Nil(t, rule.Status)
		assert.Nil(t, rule.RuleName)
	})

	t.Run("full rule block", func(t *testing.T) {
		ruleList := []interface{}{
			map[string]interface{}{
				"status":      "enable",
				"rule_name":   "test-rule",
				"description": []interface{}{"desc1", "desc2"},
			},
		}
		rule := teo.ExpandTeoL7AccRuleV2RuleBlock(ruleList)
		assert.NotNil(t, rule)
		assert.Equal(t, "enable", *rule.Status)
		assert.Equal(t, "test-rule", *rule.RuleName)
		assert.Equal(t, 2, len(rule.Description))
		assert.Equal(t, "desc1", *rule.Description[0])
		assert.Equal(t, "desc2", *rule.Description[1])
	})

	t.Run("partial rule block", func(t *testing.T) {
		ruleList := []interface{}{
			map[string]interface{}{
				"status":    "disable",
				"rule_name": "",
			},
		}
		rule := teo.ExpandTeoL7AccRuleV2RuleBlock(ruleList)
		assert.NotNil(t, rule)
		assert.Equal(t, "disable", *rule.Status)
		assert.Nil(t, rule.RuleName)
		assert.Nil(t, rule.Description)
		assert.Nil(t, rule.Branches)
	})
}

// TestTeoL7AccRuleV2_FlattenRuleBlock tests the flattenTeoL7AccRuleV2RuleBlock helper function
func TestTeoL7AccRuleV2_FlattenRuleBlock(t *testing.T) {
	t.Run("nil rule returns empty list", func(t *testing.T) {
		result := teo.FlattenTeoL7AccRuleV2RuleBlock(nil)
		assert.Equal(t, 0, len(result))
	})

	t.Run("full rule", func(t *testing.T) {
		rule := &teov20220901.RuleEngineItem{
			Status:       ptrString("enable"),
			RuleName:     ptrString("test-rule"),
			Description:  []*string{ptrString("desc1"), ptrString("desc2")},
			RulePriority: ptrInt64(10),
		}
		result := teo.FlattenTeoL7AccRuleV2RuleBlock(rule)
		assert.Equal(t, 1, len(result))
		ruleMap := result[0].(map[string]interface{})
		assert.Equal(t, "enable", ruleMap["status"])
		assert.Equal(t, "test-rule", ruleMap["rule_name"])

		descList := ruleMap["description"].([]string)
		assert.Equal(t, 2, len(descList))
		assert.Equal(t, "desc1", descList[0])
		assert.Equal(t, "desc2", descList[1])
	})

	t.Run("rule with nil fields", func(t *testing.T) {
		rule := &teov20220901.RuleEngineItem{
			Status: ptrString("disable"),
		}
		result := teo.FlattenTeoL7AccRuleV2RuleBlock(rule)
		assert.Equal(t, 1, len(result))
		ruleMap := result[0].(map[string]interface{})
		assert.Equal(t, "disable", ruleMap["status"])
		_, hasRuleName := ruleMap["rule_name"]
		assert.False(t, hasRuleName)
	})
}

func ptrInt64(i int64) *int64 {
	return &i
}
