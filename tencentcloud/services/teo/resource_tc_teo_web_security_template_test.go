package teo_test

import (
	"context"
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

func TestAccTencentCloudTeoWebSecurityTemplateResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTeoWebSecurityTemplate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_teo_web_security_template.web_security_template", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_web_security_template.web_security_template", "template_id"),
					resource.TestCheckResourceAttr("tencentcloud_teo_web_security_template.web_security_template", "template_name", "tf-test"),
					resource.TestCheckResourceAttr("tencentcloud_teo_web_security_template.web_security_template", "zone_id", "zone-3fkff38fyw8s"),
					resource.TestCheckResourceAttr("tencentcloud_teo_web_security_template.web_security_template", "security_policy.#", "1"),
					// bot_management
					resource.TestCheckResourceAttr("tencentcloud_teo_web_security_template.web_security_template", "security_policy.0.bot_management.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_web_security_template.web_security_template", "security_policy.0.bot_management.0.enabled", "on"),
					resource.TestCheckResourceAttr("tencentcloud_teo_web_security_template.web_security_template", "security_policy.0.bot_management.0.basic_bot_settings.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_web_security_template.web_security_template", "security_policy.0.bot_management.0.basic_bot_settings.0.bot_intelligence.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_web_security_template.web_security_template", "security_policy.0.bot_management.0.basic_bot_settings.0.bot_intelligence.0.enabled", "on"),
					resource.TestCheckResourceAttr("tencentcloud_teo_web_security_template.web_security_template", "security_policy.0.bot_management.0.basic_bot_settings.0.ip_reputation.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_web_security_template.web_security_template", "security_policy.0.bot_management.0.basic_bot_settings.0.ip_reputation.0.enabled", "on"),
					resource.TestCheckResourceAttr("tencentcloud_teo_web_security_template.web_security_template", "security_policy.0.bot_management.0.browser_impersonation_detection.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_web_security_template.web_security_template", "security_policy.0.bot_management.0.custom_rules.#", "1"),
					// custom_rules
					resource.TestCheckResourceAttr("tencentcloud_teo_web_security_template.web_security_template", "security_policy.0.custom_rules.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_web_security_template.web_security_template", "security_policy.0.custom_rules.0.rules.#", "2"),
					// exception_rules
					resource.TestCheckResourceAttr("tencentcloud_teo_web_security_template.web_security_template", "security_policy.0.exception_rules.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_web_security_template.web_security_template", "security_policy.0.exception_rules.0.rules.#", "2"),
					// http_ddos_protection
					resource.TestCheckResourceAttr("tencentcloud_teo_web_security_template.web_security_template", "security_policy.0.http_ddos_protection.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_web_security_template.web_security_template", "security_policy.0.http_ddos_protection.0.adaptive_frequency_control.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_web_security_template.web_security_template", "security_policy.0.http_ddos_protection.0.adaptive_frequency_control.0.enabled", "on"),
					resource.TestCheckResourceAttr("tencentcloud_teo_web_security_template.web_security_template", "security_policy.0.http_ddos_protection.0.adaptive_frequency_control.0.sensitivity", "Loose"),
					resource.TestCheckResourceAttr("tencentcloud_teo_web_security_template.web_security_template", "security_policy.0.http_ddos_protection.0.bandwidth_abuse_defense.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_web_security_template.web_security_template", "security_policy.0.http_ddos_protection.0.bandwidth_abuse_defense.0.enabled", "on"),
					resource.TestCheckResourceAttr("tencentcloud_teo_web_security_template.web_security_template", "security_policy.0.http_ddos_protection.0.client_filtering.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_web_security_template.web_security_template", "security_policy.0.http_ddos_protection.0.client_filtering.0.enabled", "on"),
					resource.TestCheckResourceAttr("tencentcloud_teo_web_security_template.web_security_template", "security_policy.0.http_ddos_protection.0.slow_attack_defense.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_web_security_template.web_security_template", "security_policy.0.http_ddos_protection.0.slow_attack_defense.0.enabled", "on"),
					// rate_limiting_rules
					resource.TestCheckResourceAttr("tencentcloud_teo_web_security_template.web_security_template", "security_policy.0.rate_limiting_rules.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_web_security_template.web_security_template", "security_policy.0.rate_limiting_rules.0.rules.#", "1"),
				),
			},
			{
				ResourceName:      "tencentcloud_teo_web_security_template.web_security_template",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccTeoWebSecurityTemplateUp,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_teo_web_security_template.web_security_template", "id"),
					resource.TestCheckResourceAttr("tencentcloud_teo_web_security_template.web_security_template", "template_name", "tf-test"),
					resource.TestCheckResourceAttr("tencentcloud_teo_web_security_template.web_security_template", "zone_id", "zone-3fkff38fyw8s"),
					resource.TestCheckResourceAttr("tencentcloud_teo_web_security_template.web_security_template", "security_policy.#", "1"),
					// bot_management
					resource.TestCheckResourceAttr("tencentcloud_teo_web_security_template.web_security_template", "security_policy.0.bot_management.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_web_security_template.web_security_template", "security_policy.0.bot_management.0.enabled", "off"),
					// custom_rules
					resource.TestCheckResourceAttr("tencentcloud_teo_web_security_template.web_security_template", "security_policy.0.custom_rules.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_web_security_template.web_security_template", "security_policy.0.custom_rules.0.rules.#", "2"),
					// exception_rules
					resource.TestCheckResourceAttr("tencentcloud_teo_web_security_template.web_security_template", "security_policy.0.exception_rules.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_web_security_template.web_security_template", "security_policy.0.exception_rules.0.rules.#", "2"),
					// http_ddos_protection
					resource.TestCheckResourceAttr("tencentcloud_teo_web_security_template.web_security_template", "security_policy.0.http_ddos_protection.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_web_security_template.web_security_template", "security_policy.0.http_ddos_protection.0.adaptive_frequency_control.0.enabled", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_web_security_template.web_security_template", "security_policy.0.http_ddos_protection.0.bandwidth_abuse_defense.0.enabled", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_web_security_template.web_security_template", "security_policy.0.http_ddos_protection.0.client_filtering.0.enabled", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_web_security_template.web_security_template", "security_policy.0.http_ddos_protection.0.slow_attack_defense.0.enabled", "off"),
					// rate_limiting_rules
					resource.TestCheckResourceAttr("tencentcloud_teo_web_security_template.web_security_template", "security_policy.0.rate_limiting_rules.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_web_security_template.web_security_template", "security_policy.0.rate_limiting_rules.0.rules.#", "1"),
				),
			},
		},
	})
}

const testAccTeoWebSecurityTemplate = `
resource "tencentcloud_teo_web_security_template" "web_security_template" {
  template_name = "tf-test"
  zone_id       = "zone-3fkff38fyw8s"
  security_policy {
    bot_management {
      enabled = "on"
      basic_bot_settings {
        bot_intelligence {
          enabled = "on"
          bot_ratings {
            high_risk_bot_requests_action {
              name = "Monitor"
            }
            human_requests_action {
              name = "Allow"
            }
            likely_bot_requests_action {
              name = "Monitor"
            }
            verified_bot_requests_action {
              name = "Monitor"
            }
          }
        }
        ip_reputation {
          enabled = "on"
          ip_reputation_group {
          }
        }
        known_bot_categories {
          bot_management_action_overrides {
            ids = ["9395241960"]
            action {
              name = "Allow"
            }
          }
        }
        search_engine_bots {
          bot_management_action_overrides {
            ids = ["9126905504"]
            action {
              name = "Deny"
            }
          }
        }
        source_idc {
          bot_management_action_overrides {
            ids = ["8868370049", "8868370048"]
            action {
              name = "Deny"
            }
          }
        }
      }
      browser_impersonation_detection {
        rules {
          condition = "$${http.request.uri.path} like ['/*'] and $${http.request.method} in ['get']"
          enabled   = "on"
          name      = "Block Non-Browser Crawler Access"
          action {
            bot_session_validation {
              issue_new_bot_session_cookie = "on"
              max_new_session_trigger_config {
                max_new_session_count_interval  = "10s"
                max_new_session_count_threshold = 300
              }
              session_expired_action {
                name = "Deny"
              }
              session_invalid_action {
                name = "Deny"
                deny_action_parameters {
                  block_ip           = null
                  block_ip_duration  = null
                  error_page_id      = null
                  response_code      = null
                  return_custom_page = null
                  stall              = "on"
                }
              }
              session_rate_control {
                enabled = "off"
              }
            }
          }
        }
      }
      client_attestation_rules {
      }
      custom_rules {
        rules {
          condition = "$${http.request.ip} in ['222.22.22.0/24'] and $${http.request.headers['user-agent']} contain ['cURL']"
          enabled   = "on"
          name      = "Login API Request Surge Protection"
          priority  = 50
          action {
            weight = 100
            security_action {
              name = "Deny"
              deny_action_parameters {
                block_ip           = null
                block_ip_duration  = null
                error_page_id      = null
                response_code      = null
                return_custom_page = null
                stall              = "on"
              }
            }
          }
        }
      }
    }
    custom_rules {
      rules {
        condition = "$${http.request.headers['user-agent']} contain ['curl/','Wget/','ApacheBench/']"
        enabled   = "on"
        name      = "Malicious User-Agent Blacklist"
        priority  = 50
        rule_type = "PreciseMatchRule"
        action {
          name = "JSChallenge"
        }
      }
      rules {
        condition = "$${http.request.ip} in ['36']"
        enabled   = "on"
        name      = "Custom Rule"
        priority  = 0
        rule_type = "BasicAccessRule"
        action {
          name = "Monitor"
        }
      }
    }
    exception_rules {
      rules {
        condition                          = "$${http.request.method} in ['post'] and $${http.request.uri.path} in ['/api/EventLogUpload']"
        enabled                            = "on"
        managed_rule_groups_for_exception  = []
        managed_rules_for_exception        = []
        name                               = "High Frequency API Skip Rate Limit 1"
        skip_option                        = "SkipOnAllRequestFields"
        skip_scope                         = "WebSecurityModules"
        web_security_modules_for_exception = ["websec-mod-adaptive-control"]
      }
      rules {
        condition                          = "$${http.request.ip} in ['123.123.123.0/24']"
        enabled                            = "on"
        managed_rule_groups_for_exception  = []
        managed_rules_for_exception        = []
        name                               = "IP Whitelist 1"
        skip_option                        = "SkipOnAllRequestFields"
        skip_scope                         = "WebSecurityModules"
        web_security_modules_for_exception = ["websec-mod-adaptive-control", "websec-mod-bot", "websec-mod-custom-rules", "websec-mod-managed-rules", "websec-mod-rate-limiting"]
      }
    }
    http_ddos_protection {
      adaptive_frequency_control {
        enabled     = "on"
        sensitivity = "Loose"
        action {
          name = "Challenge"
          challenge_action_parameters {
            attester_id      = null
            challenge_option = "JSChallenge"
            interval         = null
          }
        }
      }
      bandwidth_abuse_defense {
        enabled = "on"
        action {
          name = "Deny"
        }
      }
      client_filtering {
        enabled = "on"
        action {
          name = "Challenge"
          challenge_action_parameters {
            attester_id      = null
            challenge_option = "JSChallenge"
            interval         = null
          }
        }
      }
      slow_attack_defense {
        enabled = "on"
        action {
          name = "Deny"
        }
        minimal_request_body_transfer_rate {
          counting_period                     = "60s"
          enabled                             = "off"
          minimal_avg_transfer_rate_threshold = "80bps"
        }
        request_body_transfer_timeout {
          enabled      = "off"
          idle_timeout = "5s"
        }
      }
    }
    rate_limiting_rules {
      rules {
        action_duration       = "30m"
        condition             = "$${http.request.uri.path} contain ['/checkout/submit']"
        count_by              = ["http.request.ip"]
        counting_period       = "60s"
        enabled               = "on"
        max_request_threshold = 300
        name                  = "Single IP Request Rate Limit 1"
        priority              = 50
        action {
          name = "Challenge"
          challenge_action_parameters {
            attester_id      = null
            challenge_option = "JSChallenge"
            interval         = null
          }
        }
      }
    }
  }
}
`

const testAccTeoWebSecurityTemplateUp = `
resource "tencentcloud_teo_web_security_template" "web_security_template" {
  template_name = "tf-test"
  zone_id       = "zone-3fkff38fyw8s"
  security_policy {
    bot_management {
      enabled = "off"
      basic_bot_settings {
        bot_intelligence {
          enabled = "on"
          bot_ratings {
            high_risk_bot_requests_action {
              name = "Monitor"
            }
            human_requests_action {
              name = "Allow"
            }
            likely_bot_requests_action {
              name = "Monitor"
            }
            verified_bot_requests_action {
              name = "Monitor"
            }
          }
        }
        ip_reputation {
          enabled = "on"
          ip_reputation_group {
          }
        }
        known_bot_categories {
          bot_management_action_overrides {
            ids = ["9395241960"]
            action {
              name = "Allow"
            }
          }
        }
        search_engine_bots {
          bot_management_action_overrides {
            ids = ["9126905504"]
            action {
              name = "Deny"
            }
          }
        }
        source_idc {
          bot_management_action_overrides {
            ids = ["8868370048", "8868370049"]
            action {
              name = "Deny"
            }
          }
        }
      }
      browser_impersonation_detection {
        rules {
          condition = "$${http.request.uri.path} like ['/*'] and $${http.request.method} in ['get']"
          enabled   = "on"
          name      = "Block Non-Browser Crawler Access"
          action {
            bot_session_validation {
              issue_new_bot_session_cookie = "on"
              max_new_session_trigger_config {
                max_new_session_count_interval  = "10s"
                max_new_session_count_threshold = 300
              }
              session_expired_action {
                name = "Deny"
              }
              session_invalid_action {
                name = "Deny"
                deny_action_parameters {
                  block_ip           = null
                  block_ip_duration  = null
                  error_page_id      = null
                  response_code      = null
                  return_custom_page = null
                  stall              = "on"
                }
              }
              session_rate_control {
                enabled = "off"
              }
            }
          }
        }
      }
      client_attestation_rules {
      }
      custom_rules {
        rules {
          condition = "$${http.request.ip} in ['222.22.22.0/24'] and $${http.request.headers['user-agent']} contain ['cURL']"
          enabled   = "on"
          name      = "Login API Request Surge Protection"
          priority  = 50
          action {
            weight = 100
            security_action {
              name = "Deny"
              deny_action_parameters {
                block_ip           = null
                block_ip_duration  = null
                error_page_id      = null
                response_code      = null
                return_custom_page = null
                stall              = "on"
              }
            }
          }
        }
      }
    }
    custom_rules {
      rules {
        condition = "$${http.request.headers['user-agent']} contain ['curl/','Wget/','ApacheBench/']"
        enabled   = "off"
        name      = "Malicious User-Agent Blacklist"
        priority  = 50
        rule_type = "PreciseMatchRule"
        action {
          name = "JSChallenge"
        }
      }
      rules {
        condition = "$${http.request.ip} in ['36']"
        enabled   = "off"
        name      = "Custom Rule"
        priority  = 0
        rule_type = "BasicAccessRule"
        action {
          name = "Monitor"
        }
      }
    }
    exception_rules {
      rules {
        condition                          = "$${http.request.method} in ['post'] and $${http.request.uri.path} in ['/api/EventLogUpload']"
        enabled                            = "off"
        managed_rule_groups_for_exception  = []
        managed_rules_for_exception        = []
        name                               = "High Frequency API Skip Rate Limit 1"
        skip_option                        = "SkipOnAllRequestFields"
        skip_scope                         = "WebSecurityModules"
        web_security_modules_for_exception = ["websec-mod-adaptive-control"]
      }
      rules {
        condition                          = "$${http.request.ip} in ['123.123.123.0/24']"
        enabled                            = "off"
        managed_rule_groups_for_exception  = []
        managed_rules_for_exception        = []
        name                               = "IP Whitelist 1"
        skip_option                        = "SkipOnAllRequestFields"
        skip_scope                         = "WebSecurityModules"
        web_security_modules_for_exception = ["websec-mod-adaptive-control", "websec-mod-bot", "websec-mod-custom-rules", "websec-mod-managed-rules", "websec-mod-rate-limiting"]
      }
    }
    http_ddos_protection {
      adaptive_frequency_control {
        enabled     = "off"
        sensitivity = null
      }
      bandwidth_abuse_defense {
        enabled = "off"
        action {
          name = "Deny"
        }
      }
      client_filtering {
        enabled = "off"
        action {
          name = "Challenge"
          challenge_action_parameters {
            attester_id      = null
            challenge_option = "JSChallenge"
            interval         = null
          }
        }
      }
      slow_attack_defense {
        enabled = "off"
        action {
          name = "Deny"
        }
        minimal_request_body_transfer_rate {
          counting_period                     = "60s"
          enabled                             = "off"
          minimal_avg_transfer_rate_threshold = "80bps"
        }
        request_body_transfer_timeout {
          enabled      = "off"
          idle_timeout = "5s"
        }
      }
    }
    rate_limiting_rules {
      rules {
        action_duration       = "30m"
        condition             = "$${http.request.uri.path} contain ['/checkout/submit']"
        count_by              = ["http.request.ip"]
        counting_period       = "60s"
        enabled               = "off"
        max_request_threshold = 300
        name                  = "Single IP Request Rate Limit 1"
        priority              = 50
        action {
          name = "Challenge"
          challenge_action_parameters {
            attester_id      = null
            challenge_option = "JSChallenge"
            interval         = null
          }
        }
      }
    }
  }
}
`

// --- Mock-based unit tests for template_id ---

// mockMeta implements tccommon.ProviderMeta
type teoWebSecurityTemplateMockMeta struct {
	client *connectivity.TencentCloudClient
}

func (m *teoWebSecurityTemplateMockMeta) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &teoWebSecurityTemplateMockMeta{}

func newTeoWebSecurityTemplateMockMeta() *teoWebSecurityTemplateMockMeta {
	return &teoWebSecurityTemplateMockMeta{client: &connectivity.TencentCloudClient{}}
}

func ptrInt64(v int64) *int64 {
	return &v
}

// go test ./tencentcloud/services/teo/ -run "TestTeoWebSecurityTemplate_TemplateIdSchema" -v -count=1 -gcflags="all=-l"
// TestTeoWebSecurityTemplate_TemplateIdSchema verifies template_id field exists in schema with correct properties
func TestTeoWebSecurityTemplate_TemplateIdSchema(t *testing.T) {
	res := teo.ResourceTencentCloudTeoWebSecurityTemplate()

	assert.NotNil(t, res)
	assert.Contains(t, res.Schema, "template_id", "template_id should be present in schema")

	templateIdSchema := res.Schema["template_id"]
	assert.Equal(t, schema.TypeString, templateIdSchema.Type, "template_id should be TypeString")
	assert.True(t, templateIdSchema.Optional, "template_id should be Optional")
	assert.True(t, templateIdSchema.Computed, "template_id should be Computed")
}

// go test ./tencentcloud/services/teo/ -run "TestTeoWebSecurityTemplate_ReadSetsTemplateId" -v -count=1 -gcflags="all=-l"
// TestTeoWebSecurityTemplate_ReadSetsTemplateId verifies that Read function sets template_id from composite ID
func TestTeoWebSecurityTemplate_ReadSetsTemplateId(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	// Patch UseTeoV20220901Client to return a non-nil client
	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newTeoWebSecurityTemplateMockMeta().client, "UseTeoV20220901Client", teoClient)

	// Patch DescribeWebSecurityTemplate to return a valid response
	patches.ApplyMethodFunc(teoClient, "DescribeWebSecurityTemplate", func(request *teov20220901.DescribeWebSecurityTemplateRequest) (*teov20220901.DescribeWebSecurityTemplateResponse, error) {
		resp := teov20220901.NewDescribeWebSecurityTemplateResponse()
		resp.Response = &teov20220901.DescribeWebSecurityTemplateResponseParams{
			SecurityPolicy: &teov20220901.SecurityPolicy{},
			RequestId:      ptrString("fake-request-id"),
		}
		return resp, nil
	})

	// Patch DescribeWebSecurityTemplates (list API) used by DescribeTeoWebSecurityTemplateNameById
	patches.ApplyMethodFunc(teoClient, "DescribeWebSecurityTemplates", func(request *teov20220901.DescribeWebSecurityTemplatesRequest) (*teov20220901.DescribeWebSecurityTemplatesResponse, error) {
		resp := teov20220901.NewDescribeWebSecurityTemplatesResponse()
		templateId := "temp-abcdef1234"
		templateName := "test-template"
		resp.Response = &teov20220901.DescribeWebSecurityTemplatesResponseParams{
			TotalCount: ptrInt64(1),
			SecurityPolicyTemplates: []*teov20220901.SecurityPolicyTemplateInfo{
				{
					ZoneId:       ptrString("zone-12345678"),
					TemplateId:   &templateId,
					TemplateName: &templateName,
				},
			},
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	res := teo.ResourceTencentCloudTeoWebSecurityTemplate()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":       "zone-12345678",
		"template_name": "test-template",
	})

	// Set composite ID: zoneId#templateId
	d.SetId("zone-12345678#temp-abcdef1234")

	meta := newTeoWebSecurityTemplateMockMeta()
	err := res.Read(d, meta)
	assert.NoError(t, err)

	// Verify template_id was set from the composite ID
	assert.Equal(t, "temp-abcdef1234", d.Get("template_id").(string), "template_id should be set from composite ID")
	assert.Equal(t, "zone-12345678", d.Get("zone_id").(string), "zone_id should be set from composite ID")
}

// go test ./tencentcloud/services/teo/ -run "TestTeoWebSecurityTemplate_CreateSetsTemplateId" -v -count=1 -gcflags="all=-l"
// TestTeoWebSecurityTemplate_CreateSetsTemplateId verifies that Create function sets template_id from API response
func TestTeoWebSecurityTemplate_CreateSetsTemplateId(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	// Patch UseTeoV20220901Client to return a non-nil client
	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newTeoWebSecurityTemplateMockMeta().client, "UseTeoV20220901Client", teoClient)

	// Patch CreateWebSecurityTemplateWithContext to return a response with TemplateId
	patches.ApplyMethodFunc(teoClient, "CreateWebSecurityTemplateWithContext", func(ctx context.Context, request *teov20220901.CreateWebSecurityTemplateRequest) (*teov20220901.CreateWebSecurityTemplateResponse, error) {
		resp := teov20220901.NewCreateWebSecurityTemplateResponse()
		resp.Response = &teov20220901.CreateWebSecurityTemplateResponseParams{
			TemplateId: ptrString("temp-created123"),
			RequestId:  ptrString("fake-request-id"),
		}
		return resp, nil
	})

	// Patch DescribeWebSecurityTemplate for the Read call after Create
	patches.ApplyMethodFunc(teoClient, "DescribeWebSecurityTemplate", func(request *teov20220901.DescribeWebSecurityTemplateRequest) (*teov20220901.DescribeWebSecurityTemplateResponse, error) {
		resp := teov20220901.NewDescribeWebSecurityTemplateResponse()
		resp.Response = &teov20220901.DescribeWebSecurityTemplateResponseParams{
			SecurityPolicy: &teov20220901.SecurityPolicy{},
			RequestId:      ptrString("fake-request-id"),
		}
		return resp, nil
	})

	// Patch DescribeWebSecurityTemplates (list API) used by DescribeTeoWebSecurityTemplateNameById
	patches.ApplyMethodFunc(teoClient, "DescribeWebSecurityTemplates", func(request *teov20220901.DescribeWebSecurityTemplatesRequest) (*teov20220901.DescribeWebSecurityTemplatesResponse, error) {
		resp := teov20220901.NewDescribeWebSecurityTemplatesResponse()
		templateId := "temp-created123"
		templateName := "test-template"
		resp.Response = &teov20220901.DescribeWebSecurityTemplatesResponseParams{
			TotalCount: ptrInt64(1),
			SecurityPolicyTemplates: []*teov20220901.SecurityPolicyTemplateInfo{
				{
					ZoneId:       ptrString("zone-12345678"),
					TemplateId:   &templateId,
					TemplateName: &templateName,
				},
			},
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	res := teo.ResourceTencentCloudTeoWebSecurityTemplate()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":       "zone-12345678",
		"template_name": "test-template",
	})

	meta := newTeoWebSecurityTemplateMockMeta()
	err := res.Create(d, meta)
	assert.NoError(t, err)

	// Verify composite ID was set
	assert.Equal(t, "zone-12345678#temp-created123", d.Id(), "composite ID should be set after creation")

	// Verify template_id was set from the API response
	assert.Equal(t, "temp-created123", d.Get("template_id").(string), "template_id should be set from Create API response")
}

// go test ./tencentcloud/services/teo/ -run "TestTeoWebSecurityTemplate_TemplateIdImmutable" -v -count=1 -gcflags="all=-l"
// TestTeoWebSecurityTemplate_TemplateIdImmutable verifies that template_id is in the immutable args list
func TestTeoWebSecurityTemplate_TemplateIdImmutable(t *testing.T) {
	res := teo.ResourceTencentCloudTeoWebSecurityTemplate()

	// template_id should not have ForceNew set (it's server-generated, immutability is enforced in Update)
	templateIdSchema := res.Schema["template_id"]
	assert.False(t, templateIdSchema.ForceNew, "template_id should not have ForceNew (immutability is enforced in update function)")
}
