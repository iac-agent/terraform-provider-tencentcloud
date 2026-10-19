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

func TestAccTencentCloudTeoSecurityPolicyResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{{
			Config: testAccTeoSecurityPolicy,
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttrSet("tencentcloud_teo_security_policy_config_config.example", "id"),
			),
		},
			{
				ResourceName:      "tencentcloud_teo_security_policy_config.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTeoSecurityPolicy = `
resource "tencentcloud_teo_security_policy_config" "example" {
  zone_id = "zone-37u62pwxfo8s"
  entity  = "ZoneDefaultPolicy"
  security_policy {
    custom_rules {
      rules {
        name      = "rule1"
        condition = "$${http.request.host} contain ['abc']"
        enabled   = "on"
        rule_type = "PreciseMatchRule"
        priority  = 50
        action {
          name = "BlockIP"
          block_ip_action_parameters {
            duration = "120s"
          }
        }
      }

      rules {
        name      = "rule2"
        condition = "$${http.request.ip} in ['119.28.103.58']"
        enabled   = "off"
        id        = "2182252647"
        rule_type = "BasicAccessRule"
        action {
          name = "Deny"
        }
      }
    }

    managed_rules {
      enabled           = "on"
      detection_only    = "off"
      semantic_analysis = "off"
      auto_update {
        auto_update_to_latest_version = "off"
      }

      managed_rule_groups {
        group_id          = "wafgroup-webshell-attacks"
        sensitivity_level = "strict"
        action {
          name = "Deny"
        }
      }

      managed_rule_groups {
        group_id          = "wafgroup-xxe-attacks"
        sensitivity_level = "strict"
        action {
          name = "Deny"
        }
      }

      managed_rule_groups {
        group_id          = "wafgroup-non-compliant-protocol-usages"
        sensitivity_level = "strict"
        action {
          name = "Deny"
        }
      }

      managed_rule_groups {
        group_id          = "wafgroup-file-upload-attacks"
        sensitivity_level = "strict"
        action {
          name = "Deny"
        }
      }

      managed_rule_groups {
        group_id          = "wafgroup-command-and-code-injections"
        sensitivity_level = "strict"
        action {
          name = "Deny"
        }
      }

      managed_rule_groups {
        group_id          = "wafgroup-ldap-injections"
        sensitivity_level = "strict"
        action {
          name = "Deny"
        }
      }

      managed_rule_groups {
        group_id          = "wafgroup-ssrf-attacks"
        sensitivity_level = "strict"
        action {
          name = "Deny"
        }
      }

      managed_rule_groups {
        group_id          = "wafgroup-unauthorized-accesses"
        sensitivity_level = "strict"
        action {
          name = "Deny"
        }
      }

      managed_rule_groups {
        group_id          = "wafgroup-xss-attacks"
        sensitivity_level = "strict"
        action {
          name = "Deny"
        }
      }

      managed_rule_groups {
        group_id          = "wafgroup-vulnerability-scanners"
        sensitivity_level = "strict"
        action {
          name = "Deny"
        }
      }

      managed_rule_groups {
        group_id          = "wafgroup-cms-vulnerabilities"
        sensitivity_level = "strict"
        action {
          name = "Deny"
        }
      }

      managed_rule_groups {
        group_id          = "wafgroup-other-vulnerabilities"
        sensitivity_level = "strict"
        action {
          name = "Deny"
        }
      }

      managed_rule_groups {
        group_id          = "wafgroup-sql-injections"
        sensitivity_level = "strict"
        action {
          name = "Deny"
        }
      }

      managed_rule_groups {
        group_id          = "wafgroup-unauthorized-file-accesses"
        sensitivity_level = "strict"
        action {
          name = "Deny"
        }
      }

      managed_rule_groups {
        group_id          = "wafgroup-oa-vulnerabilities"
        sensitivity_level = "strict"
        action {
          name = "Deny"
        }
      }

      managed_rule_groups {
        group_id          = "wafgroup-ssti-attacks"
        sensitivity_level = "strict"
        action {
          name = "Deny"
        }
      }

      managed_rule_groups {
        group_id          = "wafgroup-shiro-vulnerabilities"
        sensitivity_level = "strict"
        action {
          name = "Deny"
        }
      }
    }
  }
}
`

// mockMetaSecurityPolicy implements tccommon.ProviderMeta
type mockMetaSecurityPolicy struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaSecurityPolicy) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaSecurityPolicy{}

func newMockMetaSecurityPolicy() *mockMetaSecurityPolicy {
	return &mockMetaSecurityPolicy{client: &connectivity.TencentCloudClient{}}
}

func ptrStringSecurityPolicy(s string) *string {
	return &s
}

// go test ./tencentcloud/services/teo/ -run "TestBotManagementLite_ReadWithBotManagementLite" -v -count=1 -gcflags="all=-l"
// TestBotManagementLite_ReadWithBotManagementLite tests Read flattens BotManagementLite from API response
func TestBotManagementLite_ReadWithBotManagementLite(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaSecurityPolicy().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeSecurityPolicy", func(request *teov20220901.DescribeSecurityPolicyRequest) (*teov20220901.DescribeSecurityPolicyResponse, error) {
		resp := teov20220901.NewDescribeSecurityPolicyResponse()
		resp.Response = &teov20220901.DescribeSecurityPolicyResponseParams{
			SecurityPolicy: &teov20220901.SecurityPolicy{
				BotManagementLite: &teov20220901.BotManagementLite{
					CAPTCHAPageChallenge: &teov20220901.CAPTCHAPageChallenge{
						Enabled: ptrStringSecurityPolicy("on"),
					},
					AICrawlerDetection: &teov20220901.AICrawlerDetection{
						Enabled: ptrStringSecurityPolicy("on"),
						Action: &teov20220901.SecurityAction{
							Name: ptrStringSecurityPolicy("Deny"),
							DenyActionParameters: &teov20220901.DenyActionParameters{
								BlockIp:         ptrStringSecurityPolicy("on"),
								BlockIpDuration: ptrStringSecurityPolicy("120s"),
							},
						},
					},
				},
			},
			RequestId: ptrStringSecurityPolicy("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaSecurityPolicy()
	res := teo.ResourceTencentCloudTeoSecurityPolicyConfig()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-12345678",
		"entity":  "ZoneDefaultPolicy",
	})
	d.SetId("zone-12345678#ZoneDefaultPolicy")

	err := res.Read(d, meta)
	assert.NoError(t, err)

	securityPolicy := d.Get("security_policy").([]interface{})
	assert.Len(t, securityPolicy, 1)
	spMap := securityPolicy[0].(map[string]interface{})

	botMgmtLite := spMap["bot_management_lite"].([]interface{})
	assert.Len(t, botMgmtLite, 1)
	bmlMap := botMgmtLite[0].(map[string]interface{})

	captchaPageChallenge := bmlMap["captcha_page_challenge"].([]interface{})
	assert.Len(t, captchaPageChallenge, 1)
	cpcMap := captchaPageChallenge[0].(map[string]interface{})
	assert.Equal(t, "on", cpcMap["enabled"])

	aiCrawlerDetection := bmlMap["ai_crawler_detection"].([]interface{})
	assert.Len(t, aiCrawlerDetection, 1)
	acdMap := aiCrawlerDetection[0].(map[string]interface{})
	assert.Equal(t, "on", acdMap["enabled"])

	action := acdMap["action"].([]interface{})
	assert.Len(t, action, 1)
	actionMap := action[0].(map[string]interface{})
	assert.Equal(t, "Deny", actionMap["name"])

	denyParams := actionMap["deny_action_parameters"].([]interface{})
	assert.Len(t, denyParams, 1)
	denyMap := denyParams[0].(map[string]interface{})
	assert.Equal(t, "on", denyMap["block_ip"])
	assert.Equal(t, "120s", denyMap["block_ip_duration"])
}

// TestBotManagementLite_ReadWithNilBotManagementLite tests Read when BotManagementLite is nil
func TestBotManagementLite_ReadWithNilBotManagementLite(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaSecurityPolicy().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeSecurityPolicy", func(request *teov20220901.DescribeSecurityPolicyRequest) (*teov20220901.DescribeSecurityPolicyResponse, error) {
		resp := teov20220901.NewDescribeSecurityPolicyResponse()
		resp.Response = &teov20220901.DescribeSecurityPolicyResponseParams{
			SecurityPolicy: &teov20220901.SecurityPolicy{},
			RequestId:      ptrStringSecurityPolicy("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaSecurityPolicy()
	res := teo.ResourceTencentCloudTeoSecurityPolicyConfig()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-12345678",
		"entity":  "ZoneDefaultPolicy",
	})
	d.SetId("zone-12345678#ZoneDefaultPolicy")

	err := res.Read(d, meta)
	assert.NoError(t, err)

	securityPolicy := d.Get("security_policy").([]interface{})
	if len(securityPolicy) > 0 && securityPolicy[0] != nil {
		spMap := securityPolicy[0].(map[string]interface{})
		botMgmtLite := spMap["bot_management_lite"].([]interface{})
		assert.Len(t, botMgmtLite, 0)
	}
}

// TestBotManagementLite_ReadWithPartialBotManagementLite tests Read when only CAPTCHAPageChallenge is set
func TestBotManagementLite_ReadWithPartialBotManagementLite(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaSecurityPolicy().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeSecurityPolicy", func(request *teov20220901.DescribeSecurityPolicyRequest) (*teov20220901.DescribeSecurityPolicyResponse, error) {
		resp := teov20220901.NewDescribeSecurityPolicyResponse()
		resp.Response = &teov20220901.DescribeSecurityPolicyResponseParams{
			SecurityPolicy: &teov20220901.SecurityPolicy{
				BotManagementLite: &teov20220901.BotManagementLite{
					CAPTCHAPageChallenge: &teov20220901.CAPTCHAPageChallenge{
						Enabled: ptrStringSecurityPolicy("on"),
					},
				},
			},
			RequestId: ptrStringSecurityPolicy("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaSecurityPolicy()
	res := teo.ResourceTencentCloudTeoSecurityPolicyConfig()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-12345678",
		"entity":  "ZoneDefaultPolicy",
	})
	d.SetId("zone-12345678#ZoneDefaultPolicy")

	err := res.Read(d, meta)
	assert.NoError(t, err)

	securityPolicy := d.Get("security_policy").([]interface{})
	assert.Len(t, securityPolicy, 1)
	spMap := securityPolicy[0].(map[string]interface{})

	botMgmtLite := spMap["bot_management_lite"].([]interface{})
	assert.Len(t, botMgmtLite, 1)
	bmlMap := botMgmtLite[0].(map[string]interface{})

	captchaPageChallenge := bmlMap["captcha_page_challenge"].([]interface{})
	assert.Len(t, captchaPageChallenge, 1)
	cpcMap := captchaPageChallenge[0].(map[string]interface{})
	assert.Equal(t, "on", cpcMap["enabled"])

	aiCrawlerDetection := bmlMap["ai_crawler_detection"].([]interface{})
	assert.Len(t, aiCrawlerDetection, 0)
}

// TestBotManagementLite_ReadWithAllowAction tests Read with Allow action parameters
func TestBotManagementLite_ReadWithAllowAction(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaSecurityPolicy().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeSecurityPolicy", func(request *teov20220901.DescribeSecurityPolicyRequest) (*teov20220901.DescribeSecurityPolicyResponse, error) {
		resp := teov20220901.NewDescribeSecurityPolicyResponse()
		resp.Response = &teov20220901.DescribeSecurityPolicyResponseParams{
			SecurityPolicy: &teov20220901.SecurityPolicy{
				BotManagementLite: &teov20220901.BotManagementLite{
					AICrawlerDetection: &teov20220901.AICrawlerDetection{
						Enabled: ptrStringSecurityPolicy("on"),
						Action: &teov20220901.SecurityAction{
							Name: ptrStringSecurityPolicy("Allow"),
							AllowActionParameters: &teov20220901.AllowActionParameters{
								MinDelayTime: ptrStringSecurityPolicy("0s"),
								MaxDelayTime: ptrStringSecurityPolicy("5s"),
							},
						},
					},
				},
			},
			RequestId: ptrStringSecurityPolicy("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaSecurityPolicy()
	res := teo.ResourceTencentCloudTeoSecurityPolicyConfig()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-12345678",
		"entity":  "ZoneDefaultPolicy",
	})
	d.SetId("zone-12345678#ZoneDefaultPolicy")

	err := res.Read(d, meta)
	assert.NoError(t, err)

	securityPolicy := d.Get("security_policy").([]interface{})
	assert.Len(t, securityPolicy, 1)
	spMap := securityPolicy[0].(map[string]interface{})

	botMgmtLite := spMap["bot_management_lite"].([]interface{})
	assert.Len(t, botMgmtLite, 1)
	bmlMap := botMgmtLite[0].(map[string]interface{})

	aiCrawlerDetection := bmlMap["ai_crawler_detection"].([]interface{})
	assert.Len(t, aiCrawlerDetection, 1)
	acdMap := aiCrawlerDetection[0].(map[string]interface{})
	assert.Equal(t, "on", acdMap["enabled"])

	action := acdMap["action"].([]interface{})
	assert.Len(t, action, 1)
	actionMap := action[0].(map[string]interface{})
	assert.Equal(t, "Allow", actionMap["name"])

	allowParams := actionMap["allow_action_parameters"].([]interface{})
	assert.Len(t, allowParams, 1)
	allowMap := allowParams[0].(map[string]interface{})
	assert.Equal(t, "0s", allowMap["min_delay_time"])
	assert.Equal(t, "5s", allowMap["max_delay_time"])
}

// TestBotManagementLite_ReadWithChallengeAction tests Read with Challenge action parameters
func TestBotManagementLite_ReadWithChallengeAction(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaSecurityPolicy().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeSecurityPolicy", func(request *teov20220901.DescribeSecurityPolicyRequest) (*teov20220901.DescribeSecurityPolicyResponse, error) {
		resp := teov20220901.NewDescribeSecurityPolicyResponse()
		resp.Response = &teov20220901.DescribeSecurityPolicyResponseParams{
			SecurityPolicy: &teov20220901.SecurityPolicy{
				BotManagementLite: &teov20220901.BotManagementLite{
					AICrawlerDetection: &teov20220901.AICrawlerDetection{
						Enabled: ptrStringSecurityPolicy("on"),
						Action: &teov20220901.SecurityAction{
							Name: ptrStringSecurityPolicy("Challenge"),
							ChallengeActionParameters: &teov20220901.ChallengeActionParameters{
								ChallengeOption: ptrStringSecurityPolicy("JSChallenge"),
								Interval:        ptrStringSecurityPolicy("300s"),
							},
						},
					},
				},
			},
			RequestId: ptrStringSecurityPolicy("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaSecurityPolicy()
	res := teo.ResourceTencentCloudTeoSecurityPolicyConfig()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-12345678",
		"entity":  "ZoneDefaultPolicy",
	})
	d.SetId("zone-12345678#ZoneDefaultPolicy")

	err := res.Read(d, meta)
	assert.NoError(t, err)

	securityPolicy := d.Get("security_policy").([]interface{})
	assert.Len(t, securityPolicy, 1)
	spMap := securityPolicy[0].(map[string]interface{})

	botMgmtLite := spMap["bot_management_lite"].([]interface{})
	assert.Len(t, botMgmtLite, 1)
	bmlMap := botMgmtLite[0].(map[string]interface{})

	aiCrawlerDetection := bmlMap["ai_crawler_detection"].([]interface{})
	assert.Len(t, aiCrawlerDetection, 1)
	acdMap := aiCrawlerDetection[0].(map[string]interface{})
	assert.Equal(t, "on", acdMap["enabled"])

	action := acdMap["action"].([]interface{})
	assert.Len(t, action, 1)
	actionMap := action[0].(map[string]interface{})
	assert.Equal(t, "Challenge", actionMap["name"])

	challengeParams := actionMap["challenge_action_parameters"].([]interface{})
	assert.Len(t, challengeParams, 1)
	challengeMap := challengeParams[0].(map[string]interface{})
	assert.Equal(t, "JSChallenge", challengeMap["challenge_option"])
	assert.Equal(t, "300s", challengeMap["interval"])
}

// TestBotManagementLite_UpdateExpand tests Update expands bot_management_lite into ModifySecurityPolicy request
func TestBotManagementLite_UpdateExpand(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaSecurityPolicy().client, "UseTeoV20220901Client", teoClient)

	var capturedRequest *teov20220901.ModifySecurityPolicyRequest
	patches.ApplyMethodFunc(teoClient, "ModifySecurityPolicyWithContext", func(_ context.Context, request *teov20220901.ModifySecurityPolicyRequest) (*teov20220901.ModifySecurityPolicyResponse, error) {
		capturedRequest = request
		resp := teov20220901.NewModifySecurityPolicyResponse()
		resp.Response = &teov20220901.ModifySecurityPolicyResponseParams{
			RequestId: ptrStringSecurityPolicy("fake-request-id"),
		}
		return resp, nil
	})

	// Also mock DescribeSecurityPolicy for the Read call after Update
	patches.ApplyMethodFunc(teoClient, "DescribeSecurityPolicy", func(request *teov20220901.DescribeSecurityPolicyRequest) (*teov20220901.DescribeSecurityPolicyResponse, error) {
		resp := teov20220901.NewDescribeSecurityPolicyResponse()
		resp.Response = &teov20220901.DescribeSecurityPolicyResponseParams{
			SecurityPolicy: &teov20220901.SecurityPolicy{
				BotManagementLite: &teov20220901.BotManagementLite{
					CAPTCHAPageChallenge: &teov20220901.CAPTCHAPageChallenge{
						Enabled: ptrStringSecurityPolicy("on"),
					},
					AICrawlerDetection: &teov20220901.AICrawlerDetection{
						Enabled: ptrStringSecurityPolicy("on"),
						Action: &teov20220901.SecurityAction{
							Name: ptrStringSecurityPolicy("Deny"),
							DenyActionParameters: &teov20220901.DenyActionParameters{
								BlockIp:         ptrStringSecurityPolicy("on"),
								BlockIpDuration: ptrStringSecurityPolicy("120s"),
							},
						},
					},
				},
			},
			RequestId: ptrStringSecurityPolicy("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaSecurityPolicy()
	res := teo.ResourceTencentCloudTeoSecurityPolicyConfig()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-12345678",
		"entity":  "ZoneDefaultPolicy",
		"security_policy": []interface{}{
			map[string]interface{}{
				"bot_management_lite": []interface{}{
					map[string]interface{}{
						"captcha_page_challenge": []interface{}{
							map[string]interface{}{
								"enabled": "on",
							},
						},
						"ai_crawler_detection": []interface{}{
							map[string]interface{}{
								"enabled": "on",
								"action": []interface{}{
									map[string]interface{}{
										"name": "Deny",
										"deny_action_parameters": []interface{}{
											map[string]interface{}{
												"block_ip":          "on",
												"block_ip_duration": "120s",
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	})
	d.SetId("zone-12345678#ZoneDefaultPolicy")

	err := res.Update(d, meta)
	assert.NoError(t, err)
	assert.NotNil(t, capturedRequest)
	assert.NotNil(t, capturedRequest.SecurityPolicy)
	assert.NotNil(t, capturedRequest.SecurityPolicy.BotManagementLite)
	assert.NotNil(t, capturedRequest.SecurityPolicy.BotManagementLite.CAPTCHAPageChallenge)
	assert.Equal(t, "on", *capturedRequest.SecurityPolicy.BotManagementLite.CAPTCHAPageChallenge.Enabled)
	assert.NotNil(t, capturedRequest.SecurityPolicy.BotManagementLite.AICrawlerDetection)
	assert.Equal(t, "on", *capturedRequest.SecurityPolicy.BotManagementLite.AICrawlerDetection.Enabled)
	assert.NotNil(t, capturedRequest.SecurityPolicy.BotManagementLite.AICrawlerDetection.Action)
	assert.Equal(t, "Deny", *capturedRequest.SecurityPolicy.BotManagementLite.AICrawlerDetection.Action.Name)
	assert.NotNil(t, capturedRequest.SecurityPolicy.BotManagementLite.AICrawlerDetection.Action.DenyActionParameters)
	assert.Equal(t, "on", *capturedRequest.SecurityPolicy.BotManagementLite.AICrawlerDetection.Action.DenyActionParameters.BlockIp)
	assert.Equal(t, "120s", *capturedRequest.SecurityPolicy.BotManagementLite.AICrawlerDetection.Action.DenyActionParameters.BlockIpDuration)
}

// TestRateLimitConfig_UpdateExpand tests Update expands rate_limit_config into ModifySecurityPolicy request
func TestRateLimitConfig_UpdateExpand(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaSecurityPolicy().client, "UseTeoV20220901Client", teoClient)

	var capturedRequest *teov20220901.ModifySecurityPolicyRequest
	patches.ApplyMethodFunc(teoClient, "ModifySecurityPolicyWithContext", func(_ context.Context, request *teov20220901.ModifySecurityPolicyRequest) (*teov20220901.ModifySecurityPolicyResponse, error) {
		capturedRequest = request
		resp := teov20220901.NewModifySecurityPolicyResponse()
		resp.Response = &teov20220901.ModifySecurityPolicyResponseParams{
			RequestId: ptrStringSecurityPolicy("fake-request-id"),
		}
		return resp, nil
	})

	// Also mock DescribeSecurityPolicy for the Read call after Update
	patches.ApplyMethodFunc(teoClient, "DescribeSecurityPolicy", func(request *teov20220901.DescribeSecurityPolicyRequest) (*teov20220901.DescribeSecurityPolicyResponse, error) {
		resp := teov20220901.NewDescribeSecurityPolicyResponse()
		resp.Response = &teov20220901.DescribeSecurityPolicyResponseParams{
			SecurityPolicy: &teov20220901.SecurityPolicy{},
			RequestId:      ptrStringSecurityPolicy("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaSecurityPolicy()
	res := teo.ResourceTencentCloudTeoSecurityPolicyConfig()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-12345678",
		"entity":  "ZoneDefaultPolicy",
		"security_config": []interface{}{
			map[string]interface{}{
				"rate_limit_config": []interface{}{
					map[string]interface{}{
						"switch": "on",
						"rate_limit_user_rules": []interface{}{
							map[string]interface{}{
								"threshold":          100,
								"period":             60,
								"rule_name":          "test-rule",
								"action":             "drop",
								"punish_time":        120,
								"punish_time_unit":   "second",
								"rule_status":        "on",
								"rule_priority":      50,
								"freq_fields":        []interface{}{"sip"},
								"freq_scope":         []interface{}{"client_to_eo"},
								"name":               "custom-page",
								"custom_response_id": "default",
								"response_code":      567,
								"redirect_url":       "https://example.com/redirect",
								"acl_conditions": []interface{}{
									map[string]interface{}{
										"match_from":    "url",
										"match_param":   "",
										"operator":      "equal",
										"match_content": "/test",
									},
								},
							},
						},
						"rate_limit_template": []interface{}{
							map[string]interface{}{
								"mode":   "normal",
								"action": "alg",
							},
						},
						"rate_limit_intelligence": []interface{}{
							map[string]interface{}{
								"switch": "on",
								"action": "monitor",
							},
						},
						"rate_limit_customizes": []interface{}{
							map[string]interface{}{
								"threshold":        200,
								"period":           30,
								"rule_name":        "customize-rule",
								"action":           "monitor",
								"punish_time":      60,
								"punish_time_unit": "minutes",
								"rule_status":      "on",
								"rule_priority":    80,
								"freq_fields":      []interface{}{"sip"},
								"freq_scope":       []interface{}{"source_to_eo"},
								"acl_conditions": []interface{}{
									map[string]interface{}{
										"match_from":    "host",
										"match_param":   "",
										"operator":      "include",
										"match_content": "example.com",
									},
								},
							},
						},
					},
				},
			},
		},
	})
	d.SetId("zone-12345678#ZoneDefaultPolicy")

	err := res.Update(d, meta)
	assert.NoError(t, err)
	assert.NotNil(t, capturedRequest)
	assert.NotNil(t, capturedRequest.SecurityConfig)
	assert.NotNil(t, capturedRequest.SecurityConfig.RateLimitConfig)

	rlConfig := capturedRequest.SecurityConfig.RateLimitConfig
	assert.Equal(t, "on", *rlConfig.Switch)

	// Verify rate_limit_user_rules
	assert.Len(t, rlConfig.RateLimitUserRules, 1)
	userRule := rlConfig.RateLimitUserRules[0]
	assert.Equal(t, int64(100), *userRule.Threshold)
	assert.Equal(t, int64(60), *userRule.Period)
	assert.Equal(t, "test-rule", *userRule.RuleName)
	assert.Equal(t, "drop", *userRule.Action)
	assert.Equal(t, int64(120), *userRule.PunishTime)
	assert.Equal(t, "second", *userRule.PunishTimeUnit)
	assert.Equal(t, "on", *userRule.RuleStatus)
	assert.Equal(t, int64(50), *userRule.RulePriority)
	assert.Len(t, userRule.FreqFields, 1)
	assert.Equal(t, "sip", *userRule.FreqFields[0])
	assert.Len(t, userRule.FreqScope, 1)
	assert.Equal(t, "client_to_eo", *userRule.FreqScope[0])
	assert.Equal(t, "custom-page", *userRule.Name)
	assert.Equal(t, "default", *userRule.CustomResponseId)
	assert.Equal(t, int64(567), *userRule.ResponseCode)
	assert.Equal(t, "https://example.com/redirect", *userRule.RedirectUrl)
	assert.Len(t, userRule.AclConditions, 1)
	assert.Equal(t, "url", *userRule.AclConditions[0].MatchFrom)
	assert.Equal(t, "equal", *userRule.AclConditions[0].Operator)
	assert.Equal(t, "/test", *userRule.AclConditions[0].MatchContent)

	// Verify rate_limit_template
	assert.NotNil(t, rlConfig.RateLimitTemplate)
	assert.Equal(t, "normal", *rlConfig.RateLimitTemplate.Mode)
	assert.Equal(t, "alg", *rlConfig.RateLimitTemplate.Action)

	// Verify rate_limit_intelligence
	assert.NotNil(t, rlConfig.RateLimitIntelligence)
	assert.Equal(t, "on", *rlConfig.RateLimitIntelligence.Switch)
	assert.Equal(t, "monitor", *rlConfig.RateLimitIntelligence.Action)

	// Verify rate_limit_customizes
	assert.Len(t, rlConfig.RateLimitCustomizes, 1)
	customRule := rlConfig.RateLimitCustomizes[0]
	assert.Equal(t, int64(200), *customRule.Threshold)
	assert.Equal(t, int64(30), *customRule.Period)
	assert.Equal(t, "customize-rule", *customRule.RuleName)
	assert.Equal(t, "monitor", *customRule.Action)
	assert.Equal(t, int64(60), *customRule.PunishTime)
	assert.Equal(t, "minutes", *customRule.PunishTimeUnit)
	assert.Equal(t, "on", *customRule.RuleStatus)
	assert.Equal(t, int64(80), *customRule.RulePriority)
	assert.Len(t, customRule.FreqFields, 1)
	assert.Equal(t, "sip", *customRule.FreqFields[0])
	assert.Len(t, customRule.FreqScope, 1)
	assert.Equal(t, "source_to_eo", *customRule.FreqScope[0])
	assert.Len(t, customRule.AclConditions, 1)
	assert.Equal(t, "host", *customRule.AclConditions[0].MatchFrom)
	assert.Equal(t, "include", *customRule.AclConditions[0].Operator)
	assert.Equal(t, "example.com", *customRule.AclConditions[0].MatchContent)
}

// TestRateLimitConfig_DeleteNeutralizes tests Delete sets empty SecurityConfig
func TestRateLimitConfig_DeleteNeutralizes(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaSecurityPolicy().client, "UseTeoV20220901Client", teoClient)

	var capturedRequest *teov20220901.ModifySecurityPolicyRequest
	patches.ApplyMethodFunc(teoClient, "ModifySecurityPolicyWithContext", func(_ context.Context, request *teov20220901.ModifySecurityPolicyRequest) (*teov20220901.ModifySecurityPolicyResponse, error) {
		capturedRequest = request
		resp := teov20220901.NewModifySecurityPolicyResponse()
		resp.Response = &teov20220901.ModifySecurityPolicyResponseParams{
			RequestId: ptrStringSecurityPolicy("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaSecurityPolicy()
	res := teo.ResourceTencentCloudTeoSecurityPolicyConfig()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-12345678",
		"entity":  "ZoneDefaultPolicy",
	})
	d.SetId("zone-12345678#ZoneDefaultPolicy")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
	assert.NotNil(t, capturedRequest)
	assert.NotNil(t, capturedRequest.SecurityConfig)
	assert.Nil(t, capturedRequest.SecurityConfig.RateLimitConfig)
	assert.Equal(t, "zone-12345678", *capturedRequest.ZoneId)
	assert.Equal(t, "ZoneDefaultPolicy", *capturedRequest.Entity)
}

// TestRateLimitConfig_ReadEmpty tests Read when response is nil (resource not found)
func TestRateLimitConfig_ReadEmpty(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaSecurityPolicy().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeSecurityPolicy", func(request *teov20220901.DescribeSecurityPolicyRequest) (*teov20220901.DescribeSecurityPolicyResponse, error) {
		resp := teov20220901.NewDescribeSecurityPolicyResponse()
		resp.Response = &teov20220901.DescribeSecurityPolicyResponseParams{
			SecurityPolicy: nil,
			RequestId:      ptrStringSecurityPolicy("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaSecurityPolicy()
	res := teo.ResourceTencentCloudTeoSecurityPolicyConfig()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-12345678",
		"entity":  "ZoneDefaultPolicy",
	})
	d.SetId("zone-12345678#ZoneDefaultPolicy")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}

// TestRateLimitConfig_CreateInvalidEntity tests Create with illegal entity/host/template_id combination
func TestRateLimitConfig_CreateInvalidEntity(t *testing.T) {
	meta := newMockMetaSecurityPolicy()
	res := teo.ResourceTencentCloudTeoSecurityPolicyConfig()

	// ZoneDefaultPolicy with host set - illegal
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-12345678",
		"entity":  "ZoneDefaultPolicy",
		"host":    "www.example.com",
	})
	err := res.Create(d, meta)
	assert.Error(t, err)

	// Host without host - illegal
	d = schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-12345678",
		"entity":  "Host",
	})
	err = res.Create(d, meta)
	assert.Error(t, err)

	// Template without template_id - illegal
	d = schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-12345678",
		"entity":  "Template",
	})
	err = res.Create(d, meta)
	assert.Error(t, err)
}

// TestRateLimitConfig_Schema tests the rate_limit_config schema definition
func TestRateLimitConfig_Schema(t *testing.T) {
	res := teo.ResourceTencentCloudTeoSecurityPolicyConfig()

	assert.NotNil(t, res)
	assert.Contains(t, res.Schema, "security_config")

	scSchema := res.Schema["security_config"]
	assert.NotNil(t, scSchema.Elem)
	scRes := scSchema.Elem.(*schema.Resource)
	assert.Contains(t, scRes.Schema, "rate_limit_config")

	rlSchema := scRes.Schema["rate_limit_config"]
	assert.Equal(t, schema.TypeList, rlSchema.Type)
	assert.True(t, rlSchema.Optional)
	assert.True(t, rlSchema.Computed)
	assert.Equal(t, 1, rlSchema.MaxItems)

	rlRes := rlSchema.Elem.(*schema.Resource)
	assert.Contains(t, rlRes.Schema, "switch")
	assert.Contains(t, rlRes.Schema, "rate_limit_user_rules")
	assert.Contains(t, rlRes.Schema, "rate_limit_template")
	assert.Contains(t, rlRes.Schema, "rate_limit_intelligence")
	assert.Contains(t, rlRes.Schema, "rate_limit_customizes")

	// Verify rate_limit_user_rules element schema
	ulrSchema := rlRes.Schema["rate_limit_user_rules"]
	ulrRes := ulrSchema.Elem.(*schema.Resource)
	assert.Contains(t, ulrRes.Schema, "threshold")
	assert.Contains(t, ulrRes.Schema, "period")
	assert.Contains(t, ulrRes.Schema, "rule_name")
	assert.Contains(t, ulrRes.Schema, "action")
	assert.Contains(t, ulrRes.Schema, "punish_time")
	assert.Contains(t, ulrRes.Schema, "punish_time_unit")
	assert.Contains(t, ulrRes.Schema, "rule_status")
	assert.Contains(t, ulrRes.Schema, "acl_conditions")
	assert.Contains(t, ulrRes.Schema, "rule_priority")
	assert.Contains(t, ulrRes.Schema, "rule_id")
	assert.True(t, ulrRes.Schema["rule_id"].Computed)
	assert.False(t, ulrRes.Schema["rule_id"].Optional)
	assert.Contains(t, ulrRes.Schema, "freq_fields")
	assert.Contains(t, ulrRes.Schema, "update_time")
	assert.True(t, ulrRes.Schema["update_time"].Computed)
	assert.False(t, ulrRes.Schema["update_time"].Optional)
	assert.Contains(t, ulrRes.Schema, "freq_scope")
	assert.Contains(t, ulrRes.Schema, "name")
	assert.Contains(t, ulrRes.Schema, "custom_response_id")
	assert.Contains(t, ulrRes.Schema, "response_code")
	assert.Contains(t, ulrRes.Schema, "redirect_url")

	// Verify acl_conditions element schema
	acSchema := ulrRes.Schema["acl_conditions"]
	acRes := acSchema.Elem.(*schema.Resource)
	assert.Contains(t, acRes.Schema, "match_from")
	assert.Contains(t, acRes.Schema, "match_param")
	assert.Contains(t, acRes.Schema, "operator")
	assert.Contains(t, acRes.Schema, "match_content")

	// Verify rate_limit_template schema
	rltSchema := rlRes.Schema["rate_limit_template"]
	rltRes := rltSchema.Elem.(*schema.Resource)
	assert.Contains(t, rltRes.Schema, "mode")
	assert.Contains(t, rltRes.Schema, "action")
	assert.Contains(t, rltRes.Schema, "rate_limit_template_detail")
	assert.True(t, rltRes.Schema["rate_limit_template_detail"].Computed)
	assert.False(t, rltRes.Schema["rate_limit_template_detail"].Optional)

	// Verify rate_limit_template_detail is Computed-only
	rltdSchema := rltRes.Schema["rate_limit_template_detail"]
	rltdRes := rltdSchema.Elem.(*schema.Resource)
	assert.Contains(t, rltdRes.Schema, "mode")
	assert.Contains(t, rltdRes.Schema, "id")
	assert.Contains(t, rltdRes.Schema, "action")
	assert.Contains(t, rltdRes.Schema, "punish_time")
	assert.Contains(t, rltdRes.Schema, "threshold")
	assert.Contains(t, rltdRes.Schema, "period")
	for _, s := range rltdRes.Schema {
		assert.True(t, s.Computed)
		assert.False(t, s.Optional)
	}

	// Verify rate_limit_intelligence schema
	rliSchema := rlRes.Schema["rate_limit_intelligence"]
	rliRes := rliSchema.Elem.(*schema.Resource)
	assert.Contains(t, rliRes.Schema, "switch")
	assert.Contains(t, rliRes.Schema, "action")
	assert.Contains(t, rliRes.Schema, "rule_id")
	assert.True(t, rliRes.Schema["rule_id"].Computed)
	assert.False(t, rliRes.Schema["rule_id"].Optional)

	// Verify rate_limit_customizes reuses rate_limit_user_rules structure
	rlcSchema := rlRes.Schema["rate_limit_customizes"]
	rlcRes := rlcSchema.Elem.(*schema.Resource)
	assert.Contains(t, rlcRes.Schema, "threshold")
	assert.Contains(t, rlcRes.Schema, "period")
	assert.Contains(t, rlcRes.Schema, "rule_name")
	assert.Contains(t, rlcRes.Schema, "action")
	assert.Contains(t, rlcRes.Schema, "acl_conditions")
	assert.Contains(t, rlcRes.Schema, "rule_priority")
}

func TestBotManagementLite_Schema(t *testing.T) {
	res := teo.ResourceTencentCloudTeoSecurityPolicyConfig()

	assert.NotNil(t, res)
	assert.Contains(t, res.Schema, "security_policy")

	spSchema := res.Schema["security_policy"]
	assert.NotNil(t, spSchema.Elem)
	spRes := spSchema.Elem.(*schema.Resource)
	assert.Contains(t, spRes.Schema, "bot_management_lite")

	bmlSchema := spRes.Schema["bot_management_lite"]
	assert.Equal(t, schema.TypeList, bmlSchema.Type)
	assert.True(t, bmlSchema.Optional)
	assert.True(t, bmlSchema.Computed)
	assert.Equal(t, 1, bmlSchema.MaxItems)

	bmlRes := bmlSchema.Elem.(*schema.Resource)
	assert.Contains(t, bmlRes.Schema, "captcha_page_challenge")
	assert.Contains(t, bmlRes.Schema, "ai_crawler_detection")

	cpcSchema := bmlRes.Schema["captcha_page_challenge"]
	assert.Equal(t, schema.TypeList, cpcSchema.Type)
	assert.True(t, cpcSchema.Optional)
	assert.Equal(t, 1, cpcSchema.MaxItems)

	cpcRes := cpcSchema.Elem.(*schema.Resource)
	assert.Contains(t, cpcRes.Schema, "enabled")
	assert.Equal(t, schema.TypeString, cpcRes.Schema["enabled"].Type)
	assert.True(t, cpcRes.Schema["enabled"].Required)

	acdSchema := bmlRes.Schema["ai_crawler_detection"]
	assert.Equal(t, schema.TypeList, acdSchema.Type)
	assert.True(t, acdSchema.Optional)
	assert.Equal(t, 1, acdSchema.MaxItems)

	acdRes := acdSchema.Elem.(*schema.Resource)
	assert.Contains(t, acdRes.Schema, "enabled")
	assert.Contains(t, acdRes.Schema, "action")
	assert.Equal(t, schema.TypeString, acdRes.Schema["enabled"].Type)
	assert.True(t, acdRes.Schema["enabled"].Required)

	actionSchema := acdRes.Schema["action"]
	assert.Equal(t, schema.TypeList, actionSchema.Type)
	assert.True(t, actionSchema.Optional)
	assert.Equal(t, 1, actionSchema.MaxItems)

	actionRes := actionSchema.Elem.(*schema.Resource)
	assert.Contains(t, actionRes.Schema, "name")
	assert.Contains(t, actionRes.Schema, "deny_action_parameters")
	assert.Contains(t, actionRes.Schema, "allow_action_parameters")
	assert.Contains(t, actionRes.Schema, "challenge_action_parameters")

	denySchema := actionRes.Schema["deny_action_parameters"]
	denyRes := denySchema.Elem.(*schema.Resource)
	assert.Contains(t, denyRes.Schema, "block_ip")
	assert.Contains(t, denyRes.Schema, "block_ip_duration")
	assert.Contains(t, denyRes.Schema, "return_custom_page")
	assert.Contains(t, denyRes.Schema, "response_code")
	assert.Contains(t, denyRes.Schema, "error_page_id")
	assert.Contains(t, denyRes.Schema, "stall")

	allowSchema := actionRes.Schema["allow_action_parameters"]
	allowRes := allowSchema.Elem.(*schema.Resource)
	assert.Contains(t, allowRes.Schema, "min_delay_time")
	assert.Contains(t, allowRes.Schema, "max_delay_time")

	challengeSchema := actionRes.Schema["challenge_action_parameters"]
	challengeRes := challengeSchema.Elem.(*schema.Resource)
	assert.Contains(t, challengeRes.Schema, "challenge_option")
	assert.Contains(t, challengeRes.Schema, "interval")
	assert.Contains(t, challengeRes.Schema, "attester_id")
}
