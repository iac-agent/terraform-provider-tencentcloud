package teo_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svcteo "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/teo"

	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// go test -i; go test -test.run TestAccTencentCloudTeoRuleEngine_basic -v
func TestAccTencentCloudTeoRuleEngine_basic(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PRIVATE) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckRuleEngineDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTeoRuleEngine,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRuleEngineExists("tencentcloud_teo_rule_engine.basic"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_rule_engine.basic", "zone_id"),
					resource.TestCheckResourceAttr("tencentcloud_teo_rule_engine.basic", "rule_name", "rule-1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_rule_engine.basic", "status", "enable"),
					resource.TestCheckResourceAttr("tencentcloud_teo_rule_engine.basic", "rules.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_rule_engine.basic", "rules.0.actions.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_rule_engine.basic", "rules.0.or.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_rule_engine.basic", "rules.0.sub_rules.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_rule_engine.basic", "rules.0.sub_rules.0.tags.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_rule_engine.basic", "rules.0.sub_rules.0.rules.0.or.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_rule_engine.basic", "rules.0.sub_rules.0.rules.0.or.0.and.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_rule_engine.basic", "rules.0.sub_rules.0.rules.0.or.0.and.0.operator", "equal"),
					resource.TestCheckResourceAttr("tencentcloud_teo_rule_engine.basic", "rules.0.sub_rules.0.rules.0.or.0.and.0.target", "filename"),
					resource.TestCheckResourceAttr("tencentcloud_teo_rule_engine.basic", "rules.0.sub_rules.0.rules.0.or.0.and.0.ignore_case", "false"),
					resource.TestCheckResourceAttr("tencentcloud_teo_rule_engine.basic", "rules.0.sub_rules.0.rules.0.or.0.and.0.values.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_rule_engine.basic", "rules.0.sub_rules.0.rules.0.actions.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_rule_engine.basic", "rules.0.sub_rules.0.rules.0.actions.0.normal_action.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_rule_engine.basic", "rules.0.sub_rules.0.rules.0.actions.0.normal_action.0.action", "HostHeader"),
					resource.TestCheckResourceAttr("tencentcloud_teo_rule_engine.basic", "rules.0.sub_rules.0.rules.0.actions.0.normal_action.0.parameters.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_rule_engine.basic", "rules.0.sub_rules.0.rules.0.actions.0.normal_action.0.parameters.0.name", "ServerName"),
					resource.TestCheckResourceAttr("tencentcloud_teo_rule_engine.basic", "rules.0.sub_rules.0.rules.0.actions.0.normal_action.0.parameters.0.values.#", "1"),
				),
			},
			{
				ResourceName:      "tencentcloud_teo_rule_engine.basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccTeoRuleEngineUp,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRuleEngineExists("tencentcloud_teo_rule_engine.basic"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_rule_engine.basic", "zone_id"),
					resource.TestCheckResourceAttr("tencentcloud_teo_rule_engine.basic", "rule_name", "rule-up"),
					resource.TestCheckResourceAttr("tencentcloud_teo_rule_engine.basic", "status", "enable"),
					resource.TestCheckResourceAttr("tencentcloud_teo_rule_engine.basic", "rules.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_rule_engine.basic", "rules.0.actions.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_rule_engine.basic", "rules.0.or.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_rule_engine.basic", "rules.0.sub_rules.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_rule_engine.basic", "rules.0.sub_rules.0.tags.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_rule_engine.basic", "rules.0.sub_rules.0.rules.0.or.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_rule_engine.basic", "rules.0.sub_rules.0.rules.0.or.0.and.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_rule_engine.basic", "rules.0.sub_rules.0.rules.0.or.0.and.0.operator", "equal"),
					resource.TestCheckResourceAttr("tencentcloud_teo_rule_engine.basic", "rules.0.sub_rules.0.rules.0.or.0.and.0.target", "filename"),
					resource.TestCheckResourceAttr("tencentcloud_teo_rule_engine.basic", "rules.0.sub_rules.0.rules.0.or.0.and.0.ignore_case", "false"),
					resource.TestCheckResourceAttr("tencentcloud_teo_rule_engine.basic", "rules.0.sub_rules.0.rules.0.or.0.and.0.values.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_rule_engine.basic", "rules.0.sub_rules.0.rules.0.actions.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_rule_engine.basic", "rules.0.sub_rules.0.rules.0.actions.0.normal_action.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_rule_engine.basic", "rules.0.sub_rules.0.rules.0.actions.0.normal_action.0.action", "HostHeader"),
					resource.TestCheckResourceAttr("tencentcloud_teo_rule_engine.basic", "rules.0.sub_rules.0.rules.0.actions.0.normal_action.0.parameters.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_rule_engine.basic", "rules.0.sub_rules.0.rules.0.actions.0.normal_action.0.parameters.0.name", "ServerName"),
					resource.TestCheckResourceAttr("tencentcloud_teo_rule_engine.basic", "rules.0.sub_rules.0.rules.0.actions.0.normal_action.0.parameters.0.values.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_rule_engine.basic", "tags.#", "2"),
				),
			},
			{
				Config: testAccTeoRuleEngineActionUp,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRuleEngineExists("tencentcloud_teo_rule_engine.basic"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_rule_engine.basic", "zone_id"),
					resource.TestCheckResourceAttr("tencentcloud_teo_rule_engine.basic", "rule_name", "rule-up"),
					resource.TestCheckResourceAttr("tencentcloud_teo_rule_engine.basic", "status", "enable"),
					resource.TestCheckResourceAttr("tencentcloud_teo_rule_engine.basic", "rules.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_rule_engine.basic", "rules.0.or.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_rule_engine.basic", "rules.0.sub_rules.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_rule_engine.basic", "rules.0.sub_rules.0.tags.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_rule_engine.basic", "rules.0.sub_rules.0.rules.0.or.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_rule_engine.basic", "rules.0.sub_rules.0.rules.0.or.0.and.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_rule_engine.basic", "rules.0.sub_rules.0.rules.0.or.0.and.0.operator", "equal"),
					resource.TestCheckResourceAttr("tencentcloud_teo_rule_engine.basic", "rules.0.sub_rules.0.rules.0.or.0.and.0.target", "filename"),
					resource.TestCheckResourceAttr("tencentcloud_teo_rule_engine.basic", "rules.0.sub_rules.0.rules.0.or.0.and.0.ignore_case", "false"),
					resource.TestCheckResourceAttr("tencentcloud_teo_rule_engine.basic", "rules.0.sub_rules.0.rules.0.or.0.and.0.values.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_rule_engine.basic", "tags.#", "2"),
				),
			},
		},
	})
}

func testAccCheckRuleEngineDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	service := svcteo.NewTeoService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_teo_rule_engine" {
			continue
		}
		idSplit := strings.Split(rs.Primary.ID, tccommon.FILED_SP)
		if len(idSplit) != 2 {
			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
		}
		zoneId := idSplit[0]
		ruleId := idSplit[1]

		originGroup, err := service.DescribeTeoRuleEngine(ctx, zoneId, ruleId)
		if originGroup != nil {
			return fmt.Errorf("zone ruleEngine %s still exists", rs.Primary.ID)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func testAccCheckRuleEngineExists(r string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}
		idSplit := strings.Split(rs.Primary.ID, tccommon.FILED_SP)
		if len(idSplit) != 2 {
			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
		}
		zoneId := idSplit[0]
		ruleId := idSplit[1]

		service := svcteo.NewTeoService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
		originGroup, err := service.DescribeTeoRuleEngine(ctx, zoneId, ruleId)
		if originGroup == nil {
			return fmt.Errorf("zone ruleEngine %s is not found", rs.Primary.ID)
		}
		if err != nil {
			return err
		}

		return nil
	}
}

const testAccZoneVar = `
variable "zone_id" {
  default = "zone-2qtuhspy7cr6"
}
`

const testAccTeoRuleEngine = testAccZoneVar + `

resource "tencentcloud_teo_rule_engine" "basic" {
	rule_name = "rule-1"
	status    = "enable"
	zone_id   = var.zone_id

	rules {
	  actions {

		rewrite_action {
		  action = "ResponseHeader"

		  parameters {
			action = "set"
			name   = "project"
			values = [
			  "1111",
			]
		  }
		}
	  }

	  or {
		and {
		  operator = "equal"
		  target   = "extension"
		  values   = [
			"11",
		  ]
		}
	  }
	  sub_rules {
		  tags = ["test-tag",]
		  rules {
			or {
			  and {
				operator = "equal"
				target = "filename"
				ignore_case = false
				values = ["test.txt"]
			  }
			}
			actions {
				normal_action {
					action = "HostHeader"
					parameters {
						name = "ServerName"
						values = ["terraform-test.com"]
					}
				}
			}
		  }
	  }
	}
  }
`

const testAccTeoRuleEngineUp = testAccZoneVar + `

resource "tencentcloud_teo_rule_engine" "basic" {
	rule_name = "rule-up"
	status    = "enable"
	zone_id   = var.zone_id

	tags = ["keep-test-np1", "keep-test-np2"]

	rules {
	  actions {

		rewrite_action {
		  action = "ResponseHeader"

		  parameters {
			action = "set"
			name   = "project"
			values = [
			  "1111",
			]
		  }
		}
	  }

	  or {
		and {
		  operator = "equal"
		  target   = "extension"
		  values   = [
			"11",
		  ]
		}
	  }
	  sub_rules {
		  tags = ["test-tag",]
		  rules {
			or {
			  and {
				operator = "equal"
				target = "filename"
				ignore_case = false
				values = ["test.txt"]
			  }
			}
			actions {
				normal_action {
					action = "HostHeader"
					parameters {
						name = "ServerName"
						values = ["terraform-test.com"]
					}
				}
			}
		  }
	  }
	}
  }
`

const testAccTeoRuleEngineActionUp = testAccZoneVar + `

resource "tencentcloud_teo_rule_engine" "basic" {
	rule_name = "rule-up"
	status    = "enable"
	zone_id   = var.zone_id

	tags = ["keep-test-np1", "keep-test-np2"]

	rules {
	  or {
		and {
		  operator = "equal"
		  target   = "extension"
		  values   = [
			"11",
		  ]
		}
	  }
	  sub_rules {
		  tags = ["test-tag",]
		  rules {
			or {
			  and {
				operator = "equal"
				target = "filename"
				ignore_case = false
				values = ["test.txt"]
			  }
			}
			actions {
				normal_action {
					action = "HostHeader"
					parameters {
						name = "ServerName"
						values = ["terraform-test.com"]
					}
				}
			}
		  }
	  }
	}
  }
`

// go test ./tencentcloud/services/teo/ -run "TestAccTeoRuleEngineRead_WithFilters" -v -count=1 -gcflags="all=-l"
// TestAccTeoRuleEngineRead_WithFilters tests read with custom filters parameter
func TestAccTeoRuleEngineRead_WithFilters(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeRules", func(request *teov20220901.DescribeRulesRequest) (*teov20220901.DescribeRulesResponse, error) {
		assert.Equal(t, "zone-test123", *request.ZoneId)
		// Verify filters are passed correctly
		assert.Equal(t, 1, len(request.Filters))
		assert.Equal(t, "rule-id", *request.Filters[0].Name)
		assert.Equal(t, 1, len(request.Filters[0].Values))
		assert.Equal(t, "rule-test456", *request.Filters[0].Values[0])

		resp := teov20220901.NewDescribeRulesResponse()
		resp.Response = &teov20220901.DescribeRulesResponseParams{
			ZoneId: ptrString("zone-test123"),
			RuleItems: []*teov20220901.RuleItem{
				{
					RuleId:       ptrString("rule-test456"),
					RuleName:     ptrString("test-rule"),
					Status:       ptrString("enable"),
					RulePriority: ptrInt64(1),
					Tags:         []*string{ptrString("tag1")},
				},
			},
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := svcteo.ResourceTencentCloudTeoRuleEngine()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":   "zone-test123",
		"rule_name": "test-rule",
		"status":    "enable",
		"rules": []interface{}{
			map[string]interface{}{
				"or": []interface{}{
					map[string]interface{}{
						"and": []interface{}{
							map[string]interface{}{
								"operator": "equal",
								"target":   "host",
								"values":   []interface{}{"example.com"},
							},
						},
					},
				},
			},
		},
		"filters": []interface{}{
			map[string]interface{}{
				"name":   "rule-id",
				"values": []interface{}{"rule-test456"},
			},
		},
	})
	d.SetId("zone-test123#rule-test456")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-test123", d.Get("zone_id"))
	assert.Equal(t, "rule-test456", d.Get("rule_id"))
	assert.Equal(t, "test-rule", d.Get("rule_name"))
	assert.Equal(t, "enable", d.Get("status"))

	ruleItems := d.Get("rule_items").([]interface{})
	assert.Equal(t, 1, len(ruleItems))
	firstItem := ruleItems[0].(map[string]interface{})
	assert.Equal(t, "rule-test456", firstItem["rule_id"])
	assert.Equal(t, "test-rule", firstItem["rule_name"])
	assert.Equal(t, "enable", firstItem["status"])
}

// go test ./tencentcloud/services/teo/ -run "TestAccTeoRuleEngineRead_RuleItemsFlatten" -v -count=1 -gcflags="all=-l"
// TestAccTeoRuleEngineRead_RuleItemsFlatten tests that rule_items computed output is correctly flattened
func TestAccTeoRuleEngineRead_RuleItemsFlatten(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeRules", func(request *teov20220901.DescribeRulesRequest) (*teov20220901.DescribeRulesResponse, error) {
		resp := teov20220901.NewDescribeRulesResponse()
		resp.Response = &teov20220901.DescribeRulesResponseParams{
			ZoneId: ptrString("zone-test123"),
			RuleItems: []*teov20220901.RuleItem{
				{
					RuleId:       ptrString("rule-001"),
					RuleName:     ptrString("rule-one"),
					Status:       ptrString("enable"),
					RulePriority: ptrInt64(2),
					Tags:         []*string{ptrString("tag1"), ptrString("tag2")},
				},
				{
					RuleId:       ptrString("rule-002"),
					RuleName:     ptrString("rule-two"),
					Status:       ptrString("disable"),
					RulePriority: ptrInt64(1),
					Tags:         []*string{ptrString("tag3")},
				},
			},
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := svcteo.ResourceTencentCloudTeoRuleEngine()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":   "zone-test123",
		"rule_name": "rule-one",
		"status":    "enable",
		"rules": []interface{}{
			map[string]interface{}{
				"or": []interface{}{
					map[string]interface{}{
						"and": []interface{}{
							map[string]interface{}{
								"operator": "equal",
								"target":   "host",
								"values":   []interface{}{"example.com"},
							},
						},
					},
				},
			},
		},
	})
	d.SetId("zone-test123#rule-001")

	err := res.Read(d, meta)
	assert.NoError(t, err)

	ruleItems := d.Get("rule_items").([]interface{})
	assert.Equal(t, 2, len(ruleItems))

	firstItem := ruleItems[0].(map[string]interface{})
	assert.Equal(t, "rule-001", firstItem["rule_id"])
	assert.Equal(t, "rule-one", firstItem["rule_name"])
	assert.Equal(t, "enable", firstItem["status"])

	secondItem := ruleItems[1].(map[string]interface{})
	assert.Equal(t, "rule-002", secondItem["rule_id"])
	assert.Equal(t, "rule-two", secondItem["rule_name"])
	assert.Equal(t, "disable", secondItem["status"])
}

// go test ./tencentcloud/services/teo/ -run "TestAccTeoRuleEngineRead_BackwardCompatible" -v -count=1 -gcflags="all=-l"
// TestAccTeoRuleEngineRead_BackwardCompatible tests that read without filters uses the default behavior
func TestAccTeoRuleEngineRead_BackwardCompatible(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoClient", teoClient)

	callCount := 0
	patches.ApplyMethodFunc(teoClient, "DescribeRules", func(request *teov20220901.DescribeRulesRequest) (*teov20220901.DescribeRulesResponse, error) {
		callCount++
		resp := teov20220901.NewDescribeRulesResponse()
		resp.Response = &teov20220901.DescribeRulesResponseParams{
			ZoneId: ptrString("zone-test123"),
			RuleItems: []*teov20220901.RuleItem{
				{
					RuleId:       ptrString("rule-test456"),
					RuleName:     ptrString("test-rule"),
					Status:       ptrString("enable"),
					RulePriority: ptrInt64(1),
				},
			},
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := svcteo.ResourceTencentCloudTeoRuleEngine()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":   "zone-test123",
		"rule_name": "test-rule",
		"status":    "enable",
		"rules": []interface{}{
			map[string]interface{}{
				"or": []interface{}{
					map[string]interface{}{
						"and": []interface{}{
							map[string]interface{}{
								"operator": "equal",
								"target":   "host",
								"values":   []interface{}{"example.com"},
							},
						},
					},
				},
			},
		},
	})
	d.SetId("zone-test123#rule-test456")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-test123", d.Get("zone_id"))
	assert.Equal(t, "rule-test456", d.Get("rule_id"))
	assert.Equal(t, "test-rule", d.Get("rule_name"))
	assert.Equal(t, "enable", d.Get("status"))
	// Should have been called twice: once for DescribeTeoRuleEngineById, once for DescribeTeoRuleEngineByFilters
	assert.Equal(t, 2, callCount)
}

// go test ./tencentcloud/services/teo/ -run "TestAccTeoRuleEngineRead_EmptyRuleItems" -v -count=1 -gcflags="all=-l"
// TestAccTeoRuleEngineRead_EmptyRuleItems tests that empty RuleItems results in empty rule_items
func TestAccTeoRuleEngineRead_EmptyRuleItems(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeRules", func(request *teov20220901.DescribeRulesRequest) (*teov20220901.DescribeRulesResponse, error) {
		resp := teov20220901.NewDescribeRulesResponse()
		resp.Response = &teov20220901.DescribeRulesResponseParams{
			ZoneId:    ptrString("zone-test123"),
			RuleItems: []*teov20220901.RuleItem{},
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := svcteo.ResourceTencentCloudTeoRuleEngine()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":   "zone-test123",
		"rule_name": "test-rule",
		"status":    "enable",
		"rules": []interface{}{
			map[string]interface{}{
				"or": []interface{}{
					map[string]interface{}{
						"and": []interface{}{
							map[string]interface{}{
								"operator": "equal",
								"target":   "host",
								"values":   []interface{}{"example.com"},
							},
						},
					},
				},
			},
		},
	})
	d.SetId("zone-test123#rule-test456")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	// Resource should be marked as removed since no matching rule found
	assert.Equal(t, "", d.Id())
}

// go test ./tencentcloud/services/teo/ -run "TestAccTeoRuleEngineRead_SchemaValidation" -v -count=1
// TestAccTeoRuleEngineRead_SchemaValidation tests that schema contains filters and rule_items
func TestAccTeoRuleEngineRead_SchemaValidation(t *testing.T) {
	res := svcteo.ResourceTencentCloudTeoRuleEngine()

	assert.NotNil(t, res)
	assert.Contains(t, res.Schema, "filters")
	assert.Contains(t, res.Schema, "rule_items")

	filters := res.Schema["filters"]
	assert.Equal(t, schema.TypeList, filters.Type)
	assert.True(t, filters.Optional)
	assert.False(t, filters.Required)
	assert.False(t, filters.Computed)

	ruleItems := res.Schema["rule_items"]
	assert.Equal(t, schema.TypeList, ruleItems.Type)
	assert.True(t, ruleItems.Computed)
	assert.False(t, ruleItems.Required)
	assert.False(t, ruleItems.Optional)
}

func ptrInt64(v int64) *int64 {
	return &v
}
