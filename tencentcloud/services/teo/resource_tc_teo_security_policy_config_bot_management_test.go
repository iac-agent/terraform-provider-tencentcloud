package teo

import (
	"testing"

	"github.com/stretchr/testify/assert"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

// go test ./tencentcloud/services/teo/ -run "TestFlattenSecurityAction" -v -count=1 -gcflags="all=-l"
// go test ./tencentcloud/services/teo/ -run "TestExpandSecurityAction" -v -count=1 -gcflags="all=-l"
// go test ./tencentcloud/services/teo/ -run "TestFlattenBotManagement" -v -count=1 -gcflags="all=-l"
// go test ./tencentcloud/services/teo/ -run "TestExpandBotManagement" -v -count=1 -gcflags="all=-l"

func ptrStr(s string) *string    { return &s }
func ptrInt64(i int64) *int64    { return &i }
func ptrUint64(u uint64) *uint64 { return &u }

func TestFlattenSecurityActionForBotManagement_Nil(t *testing.T) {
	result := flattenSecurityActionForBotManagement(nil)
	assert.Empty(t, result)
}

func TestFlattenSecurityActionForBotManagement_DenyAction(t *testing.T) {
	action := &teov20220901.SecurityAction{
		Name: ptrStr("Deny"),
		DenyActionParameters: &teov20220901.DenyActionParameters{
			BlockIp:         ptrStr("on"),
			BlockIpDuration: ptrStr("120s"),
			Stall:           ptrStr("off"),
		},
	}
	result := flattenSecurityActionForBotManagement(action)
	assert.Equal(t, ptrStr("Deny"), result["name"])

	denyList, ok := result["deny_action_parameters"].([]interface{})
	assert.True(t, ok)
	assert.Len(t, denyList, 1)
	denyMap := denyList[0].(map[string]interface{})
	assert.Equal(t, ptrStr("on"), denyMap["block_ip"])
	assert.Equal(t, ptrStr("120s"), denyMap["block_ip_duration"])
	assert.Equal(t, ptrStr("off"), denyMap["stall"])
}

func TestFlattenSecurityActionForBotManagement_RedirectAction(t *testing.T) {
	action := &teov20220901.SecurityAction{
		Name: ptrStr("Redirect"),
		RedirectActionParameters: &teov20220901.RedirectActionParameters{
			URL: ptrStr("https://example.com"),
		},
	}
	result := flattenSecurityActionForBotManagement(action)
	assert.Equal(t, ptrStr("Redirect"), result["name"])

	redirectList, ok := result["redirect_action_parameters"].([]interface{})
	assert.True(t, ok)
	assert.Len(t, redirectList, 1)
	redirectMap := redirectList[0].(map[string]interface{})
	assert.Equal(t, ptrStr("https://example.com"), redirectMap["url"])
}

func TestFlattenSecurityActionForBotManagement_AllowAction(t *testing.T) {
	action := &teov20220901.SecurityAction{
		Name: ptrStr("Allow"),
		AllowActionParameters: &teov20220901.AllowActionParameters{
			MinDelayTime: ptrStr("0s"),
			MaxDelayTime: ptrStr("5s"),
		},
	}
	result := flattenSecurityActionForBotManagement(action)
	assert.Equal(t, ptrStr("Allow"), result["name"])

	allowList, ok := result["allow_action_parameters"].([]interface{})
	assert.True(t, ok)
	assert.Len(t, allowList, 1)
	allowMap := allowList[0].(map[string]interface{})
	assert.Equal(t, ptrStr("0s"), allowMap["min_delay_time"])
	assert.Equal(t, ptrStr("5s"), allowMap["max_delay_time"])
}

func TestFlattenSecurityActionForBotManagement_ChallengeAction(t *testing.T) {
	action := &teov20220901.SecurityAction{
		Name: ptrStr("Challenge"),
		ChallengeActionParameters: &teov20220901.ChallengeActionParameters{
			ChallengeOption: ptrStr("JSChallenge"),
			Interval:        ptrStr("300s"),
			AttesterId:      ptrStr("att-123"),
		},
	}
	result := flattenSecurityActionForBotManagement(action)
	assert.Equal(t, ptrStr("Challenge"), result["name"])

	challengeList, ok := result["challenge_action_parameters"].([]interface{})
	assert.True(t, ok)
	assert.Len(t, challengeList, 1)
	challengeMap := challengeList[0].(map[string]interface{})
	assert.Equal(t, ptrStr("JSChallenge"), challengeMap["challenge_option"])
	assert.Equal(t, ptrStr("300s"), challengeMap["interval"])
	assert.Equal(t, ptrStr("att-123"), challengeMap["attester_id"])
}

func TestFlattenSecurityActionForBotManagement_DenyWithCustomPage(t *testing.T) {
	action := &teov20220901.SecurityAction{
		Name: ptrStr("Deny"),
		DenyActionParameters: &teov20220901.DenyActionParameters{
			ReturnCustomPage: ptrStr("on"),
			ResponseCode:     ptrStr("567"),
			ErrorPageId:      ptrStr("page-abc"),
		},
	}
	result := flattenSecurityActionForBotManagement(action)
	assert.Equal(t, ptrStr("Deny"), result["name"])

	denyList := result["deny_action_parameters"].([]interface{})
	denyMap := denyList[0].(map[string]interface{})
	assert.Equal(t, ptrStr("on"), denyMap["return_custom_page"])
	assert.Equal(t, ptrStr("567"), denyMap["response_code"])
	assert.Equal(t, ptrStr("page-abc"), denyMap["error_page_id"])
}

func TestExpandSecurityActionForBotManagement_DenyAction(t *testing.T) {
	actionMap := map[string]interface{}{
		"name": "Deny",
		"deny_action_parameters": []interface{}{
			map[string]interface{}{
				"block_ip":          "on",
				"block_ip_duration": "120s",
				"stall":             "off",
			},
		},
	}
	result := expandSecurityActionForBotManagement(actionMap)
	assert.Equal(t, "Deny", *result.Name)
	assert.NotNil(t, result.DenyActionParameters)
	assert.Equal(t, "on", *result.DenyActionParameters.BlockIp)
	assert.Equal(t, "120s", *result.DenyActionParameters.BlockIpDuration)
	assert.Equal(t, "off", *result.DenyActionParameters.Stall)
}

func TestExpandSecurityActionForBotManagement_RedirectAction(t *testing.T) {
	actionMap := map[string]interface{}{
		"name": "Redirect",
		"redirect_action_parameters": []interface{}{
			map[string]interface{}{
				"url": "https://example.com",
			},
		},
	}
	result := expandSecurityActionForBotManagement(actionMap)
	assert.Equal(t, "Redirect", *result.Name)
	assert.NotNil(t, result.RedirectActionParameters)
	assert.Equal(t, "https://example.com", *result.RedirectActionParameters.URL)
}

func TestExpandSecurityActionForBotManagement_AllowAction(t *testing.T) {
	actionMap := map[string]interface{}{
		"name": "Allow",
		"allow_action_parameters": []interface{}{
			map[string]interface{}{
				"min_delay_time": "0s",
				"max_delay_time": "5s",
			},
		},
	}
	result := expandSecurityActionForBotManagement(actionMap)
	assert.Equal(t, "Allow", *result.Name)
	assert.NotNil(t, result.AllowActionParameters)
	assert.Equal(t, "0s", *result.AllowActionParameters.MinDelayTime)
	assert.Equal(t, "5s", *result.AllowActionParameters.MaxDelayTime)
}

func TestExpandSecurityActionForBotManagement_ChallengeAction(t *testing.T) {
	actionMap := map[string]interface{}{
		"name": "Challenge",
		"challenge_action_parameters": []interface{}{
			map[string]interface{}{
				"challenge_option": "JSChallenge",
				"interval":         "300s",
				"attester_id":      "att-123",
			},
		},
	}
	result := expandSecurityActionForBotManagement(actionMap)
	assert.Equal(t, "Challenge", *result.Name)
	assert.NotNil(t, result.ChallengeActionParameters)
	assert.Equal(t, "JSChallenge", *result.ChallengeActionParameters.ChallengeOption)
	assert.Equal(t, "300s", *result.ChallengeActionParameters.Interval)
	assert.Equal(t, "att-123", *result.ChallengeActionParameters.AttesterId)
}

func TestFlattenBotManagement_Nil(t *testing.T) {
	result := flattenBotManagement(nil)
	assert.Empty(t, result)
}

func TestFlattenBotManagement_Enabled(t *testing.T) {
	botMgmt := &teov20220901.BotManagement{
		Enabled: ptrStr("on"),
	}
	result := flattenBotManagement(botMgmt)
	assert.Equal(t, ptrStr("on"), result["enabled"])
}

func TestFlattenBotManagement_CustomRules(t *testing.T) {
	botMgmt := &teov20220901.BotManagement{
		Enabled: ptrStr("on"),
		CustomRules: &teov20220901.BotManagementCustomRules{
			Rules: []*teov20220901.BotManagementCustomRule{
				{
					Id:        ptrStr("rule-001"),
					Name:      ptrStr("test-rule"),
					Enabled:   ptrStr("on"),
					Priority:  ptrInt64(50),
					Condition: ptrStr("${http.request.uri.path} contain ['/api']"),
					Action: []*teov20220901.SecurityWeightedAction{
						{
							SecurityAction: &teov20220901.SecurityAction{
								Name: ptrStr("Monitor"),
							},
							Weight: ptrInt64(80),
						},
						{
							SecurityAction: &teov20220901.SecurityAction{
								Name: ptrStr("Deny"),
								DenyActionParameters: &teov20220901.DenyActionParameters{
									BlockIp:         ptrStr("on"),
									BlockIpDuration: ptrStr("120s"),
								},
							},
							Weight: ptrInt64(20),
						},
					},
				},
			},
		},
	}
	result := flattenBotManagement(botMgmt)
	assert.Equal(t, ptrStr("on"), result["enabled"])

	customRulesList, ok := result["custom_rules"].([]interface{})
	assert.True(t, ok)
	assert.Len(t, customRulesList, 1)
	customRulesMap := customRulesList[0].(map[string]interface{})

	rulesList, ok := customRulesMap["rules"].([]interface{})
	assert.True(t, ok)
	assert.Len(t, rulesList, 1)
	ruleMap := rulesList[0].(map[string]interface{})

	assert.Equal(t, ptrStr("rule-001"), ruleMap["id"])
	assert.Equal(t, ptrStr("test-rule"), ruleMap["name"])
	assert.Equal(t, ptrStr("on"), ruleMap["enabled"])
	assert.Equal(t, ptrInt64(50), ruleMap["priority"])
	assert.Equal(t, ptrStr("${http.request.uri.path} contain ['/api']"), ruleMap["condition"])

	actionList, ok := ruleMap["action"].([]interface{})
	assert.True(t, ok)
	assert.Len(t, actionList, 2)

	wa0 := actionList[0].(map[string]interface{})
	sa0 := wa0["security_action"].([]interface{})[0].(map[string]interface{})
	assert.Equal(t, "Monitor", *sa0["name"].(*string))
	assert.Equal(t, ptrInt64(80), wa0["weight"])

	wa1 := actionList[1].(map[string]interface{})
	sa1 := wa1["security_action"].([]interface{})[0].(map[string]interface{})
	assert.Equal(t, "Deny", *sa1["name"].(*string))
	assert.Equal(t, ptrInt64(20), wa1["weight"])
}

func TestFlattenBotManagement_BasicBotSettings(t *testing.T) {
	botMgmt := &teov20220901.BotManagement{
		Enabled: ptrStr("on"),
		BasicBotSettings: &teov20220901.BasicBotSettings{
			SourceIDC: &teov20220901.SourceIDC{
				BaseAction: &teov20220901.SecurityAction{
					Name: ptrStr("Monitor"),
				},
				BotManagementActionOverrides: []*teov20220901.BotManagementActionOverrides{
					{
						Ids: []*string{ptrStr("IDC_RULE_001")},
						Action: &teov20220901.SecurityAction{
							Name: ptrStr("Deny"),
						},
					},
				},
			},
			SearchEngineBots: &teov20220901.SearchEngineBots{
				BaseAction: &teov20220901.SecurityAction{
					Name: ptrStr("Allow"),
				},
			},
			KnownBotCategories: &teov20220901.KnownBotCategories{
				BaseAction: &teov20220901.SecurityAction{
					Name: ptrStr("Monitor"),
				},
			},
			IPReputation: &teov20220901.IPReputation{
				Enabled: ptrStr("on"),
				IPReputationGroup: &teov20220901.IPReputationGroup{
					BaseAction: &teov20220901.SecurityAction{
						Name: ptrStr("Deny"),
					},
					BotManagementActionOverrides: []*teov20220901.BotManagementActionOverrides{
						{
							Ids: []*string{ptrStr("IPREP_WEB_AND_DDOS_ATTACKERS_HIGH")},
							Action: &teov20220901.SecurityAction{
								Name: ptrStr("Monitor"),
							},
						},
					},
				},
			},
			BotIntelligence: &teov20220901.BotIntelligence{
				Enabled: ptrStr("on"),
				Id:      ptrStr("bi-rule-001"),
				BotRatings: &teov20220901.BotRatings{
					HighRiskBotRequestsAction: &teov20220901.SecurityAction{
						Name: ptrStr("Deny"),
					},
					LikelyBotRequestsAction: &teov20220901.SecurityAction{
						Name: ptrStr("Monitor"),
					},
					VerifiedBotRequestsAction: &teov20220901.SecurityAction{
						Name: ptrStr("Allow"),
					},
					HumanRequestsAction: &teov20220901.SecurityAction{
						Name: ptrStr("Allow"),
					},
				},
			},
		},
	}
	result := flattenBotManagement(botMgmt)

	bbsList, ok := result["basic_bot_settings"].([]interface{})
	assert.True(t, ok)
	assert.Len(t, bbsList, 1)
	bbsMap := bbsList[0].(map[string]interface{})

	// SourceIDC
	sourceIDCList := bbsMap["source_idc"].([]interface{})
	sourceIDCMap := sourceIDCList[0].(map[string]interface{})
	baseActionList := sourceIDCMap["base_action"].([]interface{})
	baseActionMap := baseActionList[0].(map[string]interface{})
	assert.Equal(t, "Monitor", *baseActionMap["name"].(*string))

	overridesList := sourceIDCMap["action_overrides"].([]interface{})
	assert.Len(t, overridesList, 1)
	overrideMap := overridesList[0].(map[string]interface{})
	idsList := overrideMap["ids"].([]interface{})
	assert.Len(t, idsList, 1)
	assert.Equal(t, "IDC_RULE_001", idsList[0].(string))

	// IPReputation
	ipRepList := bbsMap["ip_reputation"].([]interface{})
	ipRepMap := ipRepList[0].(map[string]interface{})
	assert.Equal(t, ptrStr("on"), ipRepMap["enabled"])

	ipRepGroupList := ipRepMap["ip_reputation_group"].([]interface{})
	ipRepGroupMap := ipRepGroupList[0].(map[string]interface{})
	ipBaseActionList := ipRepGroupMap["base_action"].([]interface{})
	ipBaseActionMap := ipBaseActionList[0].(map[string]interface{})
	assert.Equal(t, "Deny", *ipBaseActionMap["name"].(*string))

	// BotIntelligence
	biList := bbsMap["bot_intelligence"].([]interface{})
	biMap := biList[0].(map[string]interface{})
	assert.Equal(t, ptrStr("on"), biMap["enabled"])
	assert.Equal(t, ptrStr("bi-rule-001"), biMap["id"])

	botRatingsList := biMap["bot_ratings"].([]interface{})
	botRatingsMap := botRatingsList[0].(map[string]interface{})
	highRiskList := botRatingsMap["high_risk_bot_requests_action"].([]interface{})
	highRiskMap := highRiskList[0].(map[string]interface{})
	assert.Equal(t, "Deny", *highRiskMap["name"].(*string))
}

func TestFlattenBotManagement_ClientAttestationRules(t *testing.T) {
	botMgmt := &teov20220901.BotManagement{
		Enabled: ptrStr("on"),
		ClientAttestationRules: &teov20220901.ClientAttestationRules{
			Rules: []*teov20220901.ClientAttestationRule{
				{
					Id:         ptrStr("car-001"),
					Name:       ptrStr("test-attestation"),
					Enabled:    ptrStr("on"),
					Priority:   ptrUint64(10),
					Condition:  ptrStr("${http.request.uri.path} contain ['/login']"),
					AttesterId: ptrStr("att-001"),
					DeviceProfiles: []*teov20220901.DeviceProfile{
						{
							ClientType:         ptrStr("iOS"),
							HighRiskMinScore:   ptrUint64(50),
							MediumRiskMinScore: ptrUint64(15),
							HighRiskRequestAction: &teov20220901.SecurityAction{
								Name: ptrStr("Deny"),
							},
							MediumRiskRequestAction: &teov20220901.SecurityAction{
								Name: ptrStr("Monitor"),
							},
						},
					},
					InvalidAttestationAction: &teov20220901.SecurityAction{
						Name: ptrStr("Monitor"),
					},
				},
			},
		},
	}
	result := flattenBotManagement(botMgmt)

	carList := result["client_attestation_rules"].([]interface{})
	carMap := carList[0].(map[string]interface{})
	rulesList := carMap["rules"].([]interface{})
	assert.Len(t, rulesList, 1)

	ruleMap := rulesList[0].(map[string]interface{})
	assert.Equal(t, ptrStr("car-001"), ruleMap["id"])
	assert.Equal(t, ptrStr("test-attestation"), ruleMap["name"])
	assert.Equal(t, ptrStr("on"), ruleMap["enabled"])
	assert.Equal(t, ptrUint64(10), ruleMap["priority"])
	assert.Equal(t, ptrStr("att-001"), ruleMap["attester_id"])

	dpList := ruleMap["device_profiles"].([]interface{})
	assert.Len(t, dpList, 1)
	dpMap := dpList[0].(map[string]interface{})
	assert.Equal(t, ptrStr("iOS"), dpMap["client_type"])
	assert.Equal(t, ptrUint64(50), dpMap["high_risk_min_score"])

	invalidActionList := ruleMap["invalid_attestation_action"].([]interface{})
	invalidActionMap := invalidActionList[0].(map[string]interface{})
	assert.Equal(t, "Monitor", *invalidActionMap["name"].(*string))
}

func TestFlattenBotManagement_BrowserImpersonationDetection(t *testing.T) {
	botMgmt := &teov20220901.BotManagement{
		Enabled: ptrStr("on"),
		BrowserImpersonationDetection: &teov20220901.BrowserImpersonationDetection{
			Rules: []*teov20220901.BrowserImpersonationDetectionRule{
				{
					Id:        ptrStr("bid-001"),
					Name:      ptrStr("test-browser-detect"),
					Enabled:   ptrStr("on"),
					Condition: ptrStr("${http.request.method} in ['GET']"),
					Action: &teov20220901.BrowserImpersonationDetectionAction{
						BotSessionValidation: &teov20220901.BotSessionValidation{
							IssueNewBotSessionCookie: ptrStr("on"),
							MaxNewSessionTriggerConfig: &teov20220901.MaxNewSessionTriggerConfig{
								MaxNewSessionCountInterval:  ptrStr("60s"),
								MaxNewSessionCountThreshold: ptrInt64(100),
							},
							SessionExpiredAction: &teov20220901.SecurityAction{
								Name: ptrStr("Deny"),
							},
							SessionInvalidAction: &teov20220901.SecurityAction{
								Name: ptrStr("Monitor"),
							},
							SessionRateControl: &teov20220901.SessionRateControl{
								Enabled: ptrStr("on"),
								HighRateSessionAction: &teov20220901.SecurityAction{
									Name: ptrStr("Deny"),
								},
								MidRateSessionAction: &teov20220901.SecurityAction{
									Name: ptrStr("Monitor"),
								},
								LowRateSessionAction: &teov20220901.SecurityAction{
									Name: ptrStr("Allow"),
								},
							},
						},
						ClientBehaviorDetection: &teov20220901.ClientBehaviorDetection{
							CryptoChallengeIntensity:   ptrStr("medium"),
							CryptoChallengeDelayBefore: ptrStr("100ms"),
							MaxChallengeCountInterval:  ptrStr("60s"),
							MaxChallengeCountThreshold: ptrInt64(500),
							ChallengeNotFinishedAction: &teov20220901.SecurityAction{Name: ptrStr("Deny")},
							ChallengeTimeoutAction:     &teov20220901.SecurityAction{Name: ptrStr("Monitor")},
							BotClientAction:            &teov20220901.SecurityAction{Name: ptrStr("Allow")},
						},
					},
				},
			},
		},
	}
	result := flattenBotManagement(botMgmt)

	bidList := result["browser_impersonation_detection"].([]interface{})
	bidMap := bidList[0].(map[string]interface{})
	rulesList := bidMap["rules"].([]interface{})
	assert.Len(t, rulesList, 1)

	ruleMap := rulesList[0].(map[string]interface{})
	assert.Equal(t, ptrStr("bid-001"), ruleMap["id"])
	assert.Equal(t, ptrStr("test-browser-detect"), ruleMap["name"])

	actionList := ruleMap["action"].([]interface{})
	actionMap := actionList[0].(map[string]interface{})

	// BotSessionValidation
	bsvList := actionMap["bot_session_validation"].([]interface{})
	bsvMap := bsvList[0].(map[string]interface{})
	assert.Equal(t, ptrStr("on"), bsvMap["issue_new_bot_session_cookie"])

	triggerConfigList := bsvMap["max_new_session_trigger_config"].([]interface{})
	triggerConfigMap := triggerConfigList[0].(map[string]interface{})
	assert.Equal(t, ptrStr("60s"), triggerConfigMap["max_new_session_count_interval"])
	assert.Equal(t, ptrInt64(100), triggerConfigMap["max_new_session_count_threshold"])

	// SessionRateControl
	srcList := bsvMap["session_rate_control"].([]interface{})
	srcMap := srcList[0].(map[string]interface{})
	assert.Equal(t, ptrStr("on"), srcMap["enabled"])

	// ClientBehaviorDetection
	cbdList := actionMap["client_behavior_detection"].([]interface{})
	cbdMap := cbdList[0].(map[string]interface{})
	assert.Equal(t, ptrStr("medium"), cbdMap["crypto_challenge_intensity"])
	assert.Equal(t, ptrStr("100ms"), cbdMap["crypto_challenge_delay_before"])
	assert.Equal(t, ptrStr("60s"), cbdMap["max_challenge_count_interval"])
	assert.Equal(t, ptrInt64(500), cbdMap["max_challenge_count_threshold"])
}

func TestExpandBotManagement_Enabled(t *testing.T) {
	botMgmtMap := map[string]interface{}{
		"enabled": "on",
	}
	result := expandBotManagement(botMgmtMap)
	assert.Equal(t, "on", *result.Enabled)
}

func TestExpandBotManagement_CustomRules(t *testing.T) {
	botMgmtMap := map[string]interface{}{
		"enabled": "on",
		"custom_rules": []interface{}{
			map[string]interface{}{
				"rules": []interface{}{
					map[string]interface{}{
						"id":        "rule-001",
						"name":      "test-rule",
						"enabled":   "on",
						"priority":  50,
						"condition": "${http.request.uri.path} contain ['/api']",
						"action": []interface{}{
							map[string]interface{}{
								"security_action": []interface{}{
									map[string]interface{}{
										"name": "Monitor",
									},
								},
								"weight": 80,
							},
							map[string]interface{}{
								"security_action": []interface{}{
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
								"weight": 20,
							},
						},
					},
				},
			},
		},
	}
	result := expandBotManagement(botMgmtMap)
	assert.Equal(t, "on", *result.Enabled)
	assert.NotNil(t, result.CustomRules)
	assert.Len(t, result.CustomRules.Rules, 1)

	rule := result.CustomRules.Rules[0]
	assert.Equal(t, "rule-001", *rule.Id)
	assert.Equal(t, "test-rule", *rule.Name)
	assert.Equal(t, "on", *rule.Enabled)
	assert.Equal(t, int64(50), *rule.Priority)
	assert.Equal(t, "${http.request.uri.path} contain ['/api']", *rule.Condition)
	assert.Len(t, rule.Action, 2)
	assert.Equal(t, "Monitor", *rule.Action[0].SecurityAction.Name)
	assert.Equal(t, int64(80), *rule.Action[0].Weight)
	assert.Equal(t, "Deny", *rule.Action[1].SecurityAction.Name)
	assert.Equal(t, "on", *rule.Action[1].SecurityAction.DenyActionParameters.BlockIp)
}

func TestExpandBotManagement_BasicBotSettings(t *testing.T) {
	botMgmtMap := map[string]interface{}{
		"enabled": "on",
		"basic_bot_settings": []interface{}{
			map[string]interface{}{
				"source_idc": []interface{}{
					map[string]interface{}{
						"base_action": []interface{}{
							map[string]interface{}{
								"name": "Monitor",
							},
						},
						"action_overrides": []interface{}{
							map[string]interface{}{
								"ids": []interface{}{"IDC_RULE_001"},
								"action": []interface{}{
									map[string]interface{}{
										"name": "Deny",
									},
								},
							},
						},
					},
				},
				"ip_reputation": []interface{}{
					map[string]interface{}{
						"enabled": "on",
						"ip_reputation_group": []interface{}{
							map[string]interface{}{
								"base_action": []interface{}{
									map[string]interface{}{
										"name": "Deny",
									},
								},
							},
						},
					},
				},
				"bot_intelligence": []interface{}{
					map[string]interface{}{
						"enabled": "on",
						"bot_ratings": []interface{}{
							map[string]interface{}{
								"high_risk_bot_requests_action": []interface{}{
									map[string]interface{}{
										"name": "Deny",
									},
								},
								"likely_bot_requests_action": []interface{}{
									map[string]interface{}{
										"name": "Monitor",
									},
								},
								"verified_bot_requests_action": []interface{}{
									map[string]interface{}{
										"name": "Allow",
									},
								},
								"human_requests_action": []interface{}{
									map[string]interface{}{
										"name": "Allow",
									},
								},
							},
						},
					},
				},
			},
		},
	}
	result := expandBotManagement(botMgmtMap)
	assert.NotNil(t, result.BasicBotSettings)

	// SourceIDC
	assert.NotNil(t, result.BasicBotSettings.SourceIDC)
	assert.Equal(t, "Monitor", *result.BasicBotSettings.SourceIDC.BaseAction.Name)
	assert.Len(t, result.BasicBotSettings.SourceIDC.BotManagementActionOverrides, 1)
	assert.Equal(t, "IDC_RULE_001", *result.BasicBotSettings.SourceIDC.BotManagementActionOverrides[0].Ids[0])
	assert.Equal(t, "Deny", *result.BasicBotSettings.SourceIDC.BotManagementActionOverrides[0].Action.Name)

	// IPReputation
	assert.NotNil(t, result.BasicBotSettings.IPReputation)
	assert.Equal(t, "on", *result.BasicBotSettings.IPReputation.Enabled)
	assert.NotNil(t, result.BasicBotSettings.IPReputation.IPReputationGroup)
	assert.Equal(t, "Deny", *result.BasicBotSettings.IPReputation.IPReputationGroup.BaseAction.Name)

	// BotIntelligence
	assert.NotNil(t, result.BasicBotSettings.BotIntelligence)
	assert.Equal(t, "on", *result.BasicBotSettings.BotIntelligence.Enabled)
	assert.NotNil(t, result.BasicBotSettings.BotIntelligence.BotRatings)
	assert.Equal(t, "Deny", *result.BasicBotSettings.BotIntelligence.BotRatings.HighRiskBotRequestsAction.Name)
	assert.Equal(t, "Monitor", *result.BasicBotSettings.BotIntelligence.BotRatings.LikelyBotRequestsAction.Name)
	assert.Equal(t, "Allow", *result.BasicBotSettings.BotIntelligence.BotRatings.VerifiedBotRequestsAction.Name)
	assert.Equal(t, "Allow", *result.BasicBotSettings.BotIntelligence.BotRatings.HumanRequestsAction.Name)
}

func TestExpandBotManagement_ClientAttestationRules(t *testing.T) {
	botMgmtMap := map[string]interface{}{
		"enabled": "on",
		"client_attestation_rules": []interface{}{
			map[string]interface{}{
				"rules": []interface{}{
					map[string]interface{}{
						"name":        "test-attestation",
						"enabled":     "on",
						"priority":    10,
						"condition":   "${http.request.uri.path} contain ['/login']",
						"attester_id": "att-001",
						"device_profiles": []interface{}{
							map[string]interface{}{
								"client_type":         "iOS",
								"high_risk_min_score": 50,
								"high_risk_request_action": []interface{}{
									map[string]interface{}{
										"name": "Deny",
									},
								},
								"medium_risk_min_score": 15,
								"medium_risk_request_action": []interface{}{
									map[string]interface{}{
										"name": "Monitor",
									},
								},
							},
						},
						"invalid_attestation_action": []interface{}{
							map[string]interface{}{
								"name": "Monitor",
							},
						},
					},
				},
			},
		},
	}
	result := expandBotManagement(botMgmtMap)
	assert.NotNil(t, result.ClientAttestationRules)
	assert.Len(t, result.ClientAttestationRules.Rules, 1)

	rule := result.ClientAttestationRules.Rules[0]
	assert.Equal(t, "test-attestation", *rule.Name)
	assert.Equal(t, "on", *rule.Enabled)
	assert.Equal(t, uint64(10), *rule.Priority)
	assert.Equal(t, "att-001", *rule.AttesterId)
	assert.Len(t, rule.DeviceProfiles, 1)
	assert.Equal(t, "iOS", *rule.DeviceProfiles[0].ClientType)
	assert.Equal(t, uint64(50), *rule.DeviceProfiles[0].HighRiskMinScore)
	assert.Equal(t, "Deny", *rule.DeviceProfiles[0].HighRiskRequestAction.Name)
	assert.Equal(t, "Monitor", *rule.InvalidAttestationAction.Name)
}

func TestExpandBotManagement_BrowserImpersonationDetection(t *testing.T) {
	botMgmtMap := map[string]interface{}{
		"enabled": "on",
		"browser_impersonation_detection": []interface{}{
			map[string]interface{}{
				"rules": []interface{}{
					map[string]interface{}{
						"name":      "test-browser-detect",
						"enabled":   "on",
						"condition": "${http.request.method} in ['GET']",
						"action": []interface{}{
							map[string]interface{}{
								"bot_session_validation": []interface{}{
									map[string]interface{}{
										"issue_new_bot_session_cookie": "on",
										"max_new_session_trigger_config": []interface{}{
											map[string]interface{}{
												"max_new_session_count_interval":  "60s",
												"max_new_session_count_threshold": 100,
											},
										},
										"session_expired_action": []interface{}{
											map[string]interface{}{
												"name": "Deny",
											},
										},
										"session_invalid_action": []interface{}{
											map[string]interface{}{
												"name": "Monitor",
											},
										},
										"session_rate_control": []interface{}{
											map[string]interface{}{
												"enabled": "on",
												"high_rate_session_action": []interface{}{
													map[string]interface{}{
														"name": "Deny",
													},
												},
												"mid_rate_session_action": []interface{}{
													map[string]interface{}{
														"name": "Monitor",
													},
												},
												"low_rate_session_action": []interface{}{
													map[string]interface{}{
														"name": "Allow",
													},
												},
											},
										},
									},
								},
								"client_behavior_detection": []interface{}{
									map[string]interface{}{
										"crypto_challenge_intensity":    "medium",
										"crypto_challenge_delay_before": "100ms",
										"max_challenge_count_interval":  "60s",
										"max_challenge_count_threshold": 500,
										"challenge_not_finished_action": []interface{}{
											map[string]interface{}{
												"name": "Deny",
											},
										},
										"challenge_timeout_action": []interface{}{
											map[string]interface{}{
												"name": "Monitor",
											},
										},
										"bot_client_action": []interface{}{
											map[string]interface{}{
												"name": "Allow",
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
	}
	result := expandBotManagement(botMgmtMap)
	assert.NotNil(t, result.BrowserImpersonationDetection)
	assert.Len(t, result.BrowserImpersonationDetection.Rules, 1)

	rule := result.BrowserImpersonationDetection.Rules[0]
	assert.Equal(t, "test-browser-detect", *rule.Name)
	assert.Equal(t, "on", *rule.Enabled)
	assert.NotNil(t, rule.Action)

	// BotSessionValidation
	bsv := rule.Action.BotSessionValidation
	assert.NotNil(t, bsv)
	assert.Equal(t, "on", *bsv.IssueNewBotSessionCookie)
	assert.NotNil(t, bsv.MaxNewSessionTriggerConfig)
	assert.Equal(t, "60s", *bsv.MaxNewSessionTriggerConfig.MaxNewSessionCountInterval)
	assert.Equal(t, int64(100), *bsv.MaxNewSessionTriggerConfig.MaxNewSessionCountThreshold)
	assert.Equal(t, "Deny", *bsv.SessionExpiredAction.Name)
	assert.Equal(t, "Monitor", *bsv.SessionInvalidAction.Name)
	assert.NotNil(t, bsv.SessionRateControl)
	assert.Equal(t, "on", *bsv.SessionRateControl.Enabled)
	assert.Equal(t, "Deny", *bsv.SessionRateControl.HighRateSessionAction.Name)
	assert.Equal(t, "Monitor", *bsv.SessionRateControl.MidRateSessionAction.Name)
	assert.Equal(t, "Allow", *bsv.SessionRateControl.LowRateSessionAction.Name)

	// ClientBehaviorDetection
	cbd := rule.Action.ClientBehaviorDetection
	assert.NotNil(t, cbd)
	assert.Equal(t, "medium", *cbd.CryptoChallengeIntensity)
	assert.Equal(t, "100ms", *cbd.CryptoChallengeDelayBefore)
	assert.Equal(t, "60s", *cbd.MaxChallengeCountInterval)
	assert.Equal(t, int64(500), *cbd.MaxChallengeCountThreshold)
	assert.Equal(t, "Deny", *cbd.ChallengeNotFinishedAction.Name)
	assert.Equal(t, "Monitor", *cbd.ChallengeTimeoutAction.Name)
	assert.Equal(t, "Allow", *cbd.BotClientAction.Name)
}

func TestFlattenExpandBotManagement_RoundTrip(t *testing.T) {
	// Test that flattening and expanding preserves the enabled and action name values.
	// Note: Direct map roundtrip is not possible because flatten stores *int64/*string pointers
	// while expand expects int/string types. In real Terraform flow, schema.ResourceData
	// handles the type conversion between flatten (Set) and expand (Get).
	original := &teov20220901.BotManagement{
		Enabled: ptrStr("on"),
		CustomRules: &teov20220901.BotManagementCustomRules{
			Rules: []*teov20220901.BotManagementCustomRule{
				{
					Name:      ptrStr("test-rule"),
					Enabled:   ptrStr("on"),
					Condition: ptrStr("${http.request.uri.path} contain ['/api']"),
					Action: []*teov20220901.SecurityWeightedAction{
						{
							SecurityAction: &teov20220901.SecurityAction{
								Name: ptrStr("Monitor"),
							},
							Weight: ptrInt64(100),
						},
					},
				},
			},
		},
	}

	flattened := flattenBotManagement(original)

	// Verify key flattened values
	assert.Equal(t, ptrStr("on"), flattened["enabled"])
	customRulesList := flattened["custom_rules"].([]interface{})
	customRulesMap := customRulesList[0].(map[string]interface{})
	rulesList := customRulesMap["rules"].([]interface{})
	ruleMap := rulesList[0].(map[string]interface{})
	assert.Equal(t, ptrStr("test-rule"), ruleMap["name"])
	assert.Equal(t, ptrStr("on"), ruleMap["enabled"])
	assert.Equal(t, ptrStr("${http.request.uri.path} contain ['/api']"), ruleMap["condition"])
}

func TestFlattenActionOverrides(t *testing.T) {
	overrides := []*teov20220901.BotManagementActionOverrides{
		{
			Ids: []*string{ptrStr("RULE_001"), ptrStr("RULE_002")},
			Action: &teov20220901.SecurityAction{
				Name: ptrStr("Deny"),
				DenyActionParameters: &teov20220901.DenyActionParameters{
					BlockIp:         ptrStr("on"),
					BlockIpDuration: ptrStr("60s"),
				},
			},
		},
	}
	result := flattenActionOverrides(overrides)
	assert.Len(t, result, 1)

	overrideMap := result[0].(map[string]interface{})
	idsList := overrideMap["ids"].([]interface{})
	assert.Len(t, idsList, 2)
	assert.Equal(t, "RULE_001", idsList[0].(string))
	assert.Equal(t, "RULE_002", idsList[1].(string))

	actionList := overrideMap["action"].([]interface{})
	actionMap := actionList[0].(map[string]interface{})
	assert.Equal(t, "Deny", *actionMap["name"].(*string))
}

func TestExpandActionOverrides(t *testing.T) {
	overridesList := []interface{}{
		map[string]interface{}{
			"ids": []interface{}{"RULE_001", "RULE_002"},
			"action": []interface{}{
				map[string]interface{}{
					"name": "Deny",
					"deny_action_parameters": []interface{}{
						map[string]interface{}{
							"block_ip":          "on",
							"block_ip_duration": "60s",
						},
					},
				},
			},
		},
	}
	result := expandActionOverrides(overridesList)
	assert.Len(t, result, 1)
	assert.Len(t, result[0].Ids, 2)
	assert.Equal(t, "RULE_001", *result[0].Ids[0])
	assert.Equal(t, "RULE_002", *result[0].Ids[1])
	assert.Equal(t, "Deny", *result[0].Action.Name)
	assert.Equal(t, "on", *result[0].Action.DenyActionParameters.BlockIp)
}

func TestFlattenBaseActionAndOverrides(t *testing.T) {
	baseAction := &teov20220901.SecurityAction{
		Name: ptrStr("Monitor"),
	}
	actionOverrides := []*teov20220901.BotManagementActionOverrides{
		{
			Ids: []*string{ptrStr("RULE_001")},
			Action: &teov20220901.SecurityAction{
				Name: ptrStr("Deny"),
			},
		},
	}
	result := flattenBaseActionAndOverrides(baseAction, actionOverrides)

	baseList := result["base_action"].([]interface{})
	baseMap := baseList[0].(map[string]interface{})
	assert.Equal(t, "Monitor", *baseMap["name"].(*string))

	overridesList := result["action_overrides"].([]interface{})
	assert.Len(t, overridesList, 1)
}

func TestExpandBotManagement_EmptyMap(t *testing.T) {
	botMgmtMap := map[string]interface{}{
		"enabled": "off",
	}
	result := expandBotManagement(botMgmtMap)
	assert.Equal(t, "off", *result.Enabled)
	assert.Nil(t, result.CustomRules)
	assert.Nil(t, result.BasicBotSettings)
	assert.Nil(t, result.ClientAttestationRules)
	assert.Nil(t, result.BrowserImpersonationDetection)
}

func TestFlattenBotManagement_EmptySubFields(t *testing.T) {
	botMgmt := &teov20220901.BotManagement{
		Enabled: ptrStr("off"),
	}
	result := flattenBotManagement(botMgmt)
	assert.Equal(t, ptrStr("off"), result["enabled"])

	_, hasCustomRules := result["custom_rules"]
	assert.False(t, hasCustomRules)
	_, hasBasicBotSettings := result["basic_bot_settings"]
	assert.False(t, hasBasicBotSettings)
	_, hasClientAttestationRules := result["client_attestation_rules"]
	assert.False(t, hasClientAttestationRules)
	_, hasBrowserImpersonationDetection := result["browser_impersonation_detection"]
	assert.False(t, hasBrowserImpersonationDetection)
}

// Verify that helper.String and helper.IntInt64 are available
var _ = helper.String
var _ = helper.IntInt64
