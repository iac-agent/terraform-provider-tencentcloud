package fwm_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	fwmv20250611 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/fwm/v20250611"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/fwm"
)

// mockMeta implements tccommon.ProviderMeta
type mockMeta struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMeta) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMeta{}

func newMockMeta() *mockMeta {
	return &mockMeta{client: &connectivity.TencentCloudClient{}}
}

func ptrString(s string) *string {
	return &s
}

func ptrInt64(i int64) *int64 {
	return &i
}

// go test ./tencentcloud/services/fwm/ -run "TestFwmEdgeAclRule" -v -count=1

const edgeAclRuleId = "cfw-group-xxxx#cfw-rule-xxxx#1"

func ptrUint64(i uint64) *uint64 {
	return &i
}

// buildEdgeAclRuleResp builds a response that fills every leaf field of
// EdgeAclRuleResp so that Read coverage of the read-only fields is complete.
func buildEdgeAclRuleResp() *fwmv20250611.EdgeAclRuleResp {
	return &fwmv20250611.EdgeAclRuleResp{
		RuleId:            ptrString("cfw-rule-xxxx"),
		Sequence:          ptrInt64(1),
		Direction:         ptrInt64(1),
		SourceContent:     ptrString("1.1.1.1"),
		SourceType:        ptrString("ip"),
		SourceName:        ptrString("source-name"),
		TargetContent:     ptrString("2.2.2.2"),
		TargetType:        ptrString("ip"),
		TargetName:        ptrString("target-name"),
		Port:              ptrString("80"),
		Protocol:          ptrString("TCP"),
		RuleAction:        ptrString("accept-pass"),
		Description:       ptrString("remark."),
		Scope:             ptrString("all"),
		CountryName:       ptrString("china"),
		CityName:          ptrString("shenzhen"),
		ParamTemplateId:   ptrString("pt-xxxx"),
		ParamTemplateName: ptrString("pt-name"),
		Invalid:           ptrInt64(0),
		BelongMember: &fwmv20250611.MemberInfo{
			AppId:    ptrString("1250000000"),
			Uin:      ptrString("100000000001"),
			Nickname: ptrString("member-nick"),
			MemberId: ptrString("mem-xxxx"),
			NodeName: ptrString("dept"),
		},
		CountryCode:  ptrUint64(1),
		CityCode:     ptrUint64(2),
		IsRegion:     ptrUint64(1),
		CloudCode:    ptrString("tencent"),
		IsCloud:      ptrUint64(1),
		InstanceName: ptrString("ins-xxxx"),
		CountryKey:   ptrString("CN"),
		CityKey:      ptrString("SZ"),
		CreateTime:   ptrString("2025-01-01 00:00:00"),
		UpdateTime:   ptrString("2025-01-02 00:00:00"),
		DnsParseCnt:  ptrInt64(3),
		HitCount:     ptrInt64(9),
		LastHitTime:  ptrString("2025-01-03 00:00:00"),
	}
}

func mockDescribeEdgeAclRules(t *testing.T, patches *gomonkey.Patches, client *fwmv20250611.Client, rules ...*fwmv20250611.EdgeAclRuleResp) {
	patches.ApplyMethodFunc(client, "DescribeEdgeAclRulesWithContext", func(_ context.Context, request *fwmv20250611.DescribeEdgeAclRulesRequest) (*fwmv20250611.DescribeEdgeAclRulesResponse, error) {
		if request != nil && request.GroupId != nil {
			assert.Equal(t, "cfw-group-xxxx", *request.GroupId)
		}
		if request != nil && request.Direction != nil {
			assert.Equal(t, int64(1), *request.Direction)
		}
		if request != nil && request.Limit != nil {
			assert.Equal(t, int64(1000), *request.Limit)
		}

		resp := fwmv20250611.NewDescribeEdgeAclRulesResponse()
		total := int64(len(rules))
		resp.Response = &fwmv20250611.DescribeEdgeAclRulesResponseParams{
			TotalCount: &total,
			Rules:      rules,
			RequestId:  ptrString("fake-request-id"),
		}
		return resp, nil
	})
}

func TestFwmEdgeAclRule_Schema(t *testing.T) {
	res := fwm.ResourceTencentCloudFwmEdgeAclRule()

	assert.NotNil(t, res)
	assert.NotNil(t, res.Create)
	assert.NotNil(t, res.Read)
	assert.NotNil(t, res.Update)
	assert.NotNil(t, res.Delete)
	assert.NotNil(t, res.Importer)

	groupId := res.Schema["group_id"]
	assert.NotNil(t, groupId)
	assert.True(t, groupId.Required)
	assert.False(t, groupId.ForceNew, "group_id replacement is decided in CustomizeDiff")

	direction := res.Schema["direction"]
	assert.NotNil(t, direction)
	assert.Equal(t, schema.TypeInt, direction.Type)
	assert.True(t, direction.Required)
	assert.False(t, direction.ForceNew, "direction replacement is decided in CustomizeDiff")

	// source_type / target_type are mutable through ModifyEdgeAclRule.
	sourceType := res.Schema["source_type"]
	assert.NotNil(t, sourceType)
	assert.True(t, sourceType.Optional)
	assert.False(t, sourceType.ForceNew)

	targetType := res.Schema["target_type"]
	assert.NotNil(t, targetType)
	assert.True(t, targetType.Optional)
	assert.False(t, targetType.ForceNew)

	ruleId := res.Schema["rule_id"]
	assert.NotNil(t, ruleId)
	assert.True(t, ruleId.Computed)

	sequence := res.Schema["sequence"]
	assert.NotNil(t, sequence)
	assert.True(t, sequence.Computed)

	// read-only fields returned by DescribeEdgeAclRules
	for _, name := range []string{
		"source_name", "target_name", "param_template_name", "country_name", "city_name",
		"country_code", "city_code", "country_key", "city_key", "is_region", "cloud_code",
		"is_cloud", "instance_name", "invalid", "dns_parse_cnt", "hit_count", "last_hit_time",
		"create_time", "update_time",
	} {
		field, ok := res.Schema[name]
		assert.True(t, ok, "schema %s is missing", name)
		if ok {
			assert.True(t, field.Computed, "%s must be computed", name)
			assert.False(t, field.Optional, "%s must not be configurable", name)
		}
	}

	belongMember := res.Schema["belong_member"]
	assert.NotNil(t, belongMember)
	assert.Equal(t, schema.TypeList, belongMember.Type)
	assert.True(t, belongMember.Computed)
	elem, ok := belongMember.Elem.(*schema.Resource)
	assert.True(t, ok)
	for _, name := range []string{"app_id", "uin", "nickname", "member_id", "node_name"} {
		field, ok := elem.Schema[name]
		assert.True(t, ok, "belong_member.%s is missing", name)
		if ok {
			assert.True(t, field.Computed)
		}
	}

	// order_index is not part of the schema; sequence is read-only.
	_, ok = res.Schema["order_index"]
	assert.False(t, ok)
}

// TestFwmEdgeAclRule_Create verifies the composite id now carries direction.
func TestFwmEdgeAclRule_Create(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseFwmV20250611Client", fwmClient)

	patches.ApplyMethodFunc(fwmClient, "CreateEdgeAclRuleWithContext", func(_ context.Context, request *fwmv20250611.CreateEdgeAclRuleRequest) (*fwmv20250611.CreateEdgeAclRuleResponse, error) {
		assert.NotNil(t, request.GroupId)
		assert.Equal(t, "cfw-group-xxxx", *request.GroupId)
		assert.Len(t, request.Rules, 1)

		rule := request.Rules[0]
		assert.Equal(t, int64(1), *rule.OrderIndex)
		assert.Equal(t, int64(1), *rule.Direction)
		assert.Equal(t, "ip", *rule.SourceType)
		assert.Equal(t, "accept-pass", *rule.RuleAction)
		assert.Equal(t, "all", *rule.Scope)

		resp := fwmv20250611.NewCreateEdgeAclRuleResponse()
		resp.Response = &fwmv20250611.CreateEdgeAclRuleResponseParams{
			RuleIds:   []*string{ptrString("cfw-rule-xxxx")},
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	mockDescribeEdgeAclRules(t, patches, fwmClient, buildEdgeAclRuleResp())

	meta := newMockMeta()
	res := fwm.ResourceTencentCloudFwmEdgeAclRule()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"group_id":       "cfw-group-xxxx",
		"direction":      1,
		"source_content": "1.1.1.1",
		"source_type":    "ip",
		"target_content": "2.2.2.2",
		"target_type":    "ip",
		"port":           "80",
		"protocol":       "TCP",
		"rule_action":    "accept-pass",
		"scope":          "all",
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	// groupId#ruleId#direction
	assert.Equal(t, edgeAclRuleId, d.Id())
	assert.Equal(t, "cfw-rule-xxxx", d.Get("rule_id"))
	assert.Equal(t, 1, d.Get("sequence"))
}

// TestFwmEdgeAclRule_Create_EmptyRuleIds verifies an empty RuleIds fails Create.
func TestFwmEdgeAclRule_Create_EmptyRuleIds(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseFwmV20250611Client", fwmClient)

	patches.ApplyMethodFunc(fwmClient, "CreateEdgeAclRuleWithContext", func(_ context.Context, _ *fwmv20250611.CreateEdgeAclRuleRequest) (*fwmv20250611.CreateEdgeAclRuleResponse, error) {
		resp := fwmv20250611.NewCreateEdgeAclRuleResponse()
		resp.Response = &fwmv20250611.CreateEdgeAclRuleResponseParams{
			RuleIds:   []*string{},
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := fwm.ResourceTencentCloudFwmEdgeAclRule()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"group_id":  "cfw-group-xxxx",
		"direction": 1,
	})

	err := res.Create(d, meta)
	assert.Error(t, err)
	assert.Empty(t, d.Id())
}

// TestFwmEdgeAclRule_Create_APIError verifies Create surfaces API errors.
func TestFwmEdgeAclRule_Create_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseFwmV20250611Client", fwmClient)

	patches.ApplyMethodFunc(fwmClient, "CreateEdgeAclRuleWithContext", func(_ context.Context, _ *fwmv20250611.CreateEdgeAclRuleRequest) (*fwmv20250611.CreateEdgeAclRuleResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=InvalidParameter, Message=Invalid group_id")
	})

	meta := newMockMeta()
	res := fwm.ResourceTencentCloudFwmEdgeAclRule()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"group_id":  "cfw-group-invalid",
		"direction": 1,
	})

	err := res.Create(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "InvalidParameter")
}

// TestFwmEdgeAclRule_Read verifies Read parses the 3-segment id, queries with
// direction and maps the response back into state.
func TestFwmEdgeAclRule_Read(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseFwmV20250611Client", fwmClient)
	mockDescribeEdgeAclRules(t, patches, fwmClient, buildEdgeAclRuleResp())

	meta := newMockMeta()
	res := fwm.ResourceTencentCloudFwmEdgeAclRule()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"group_id":  "cfw-group-xxxx",
		"direction": 1,
	})
	d.SetId(edgeAclRuleId)

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, edgeAclRuleId, d.Id())
	assert.Equal(t, "cfw-rule-xxxx", d.Get("rule_id"))
	assert.Equal(t, 1, d.Get("sequence"))
	assert.Equal(t, 1, d.Get("direction"))
	assert.Equal(t, "mem-xxxx", d.Get("belong_member_id"))
	assert.Equal(t, "cfw-group-xxxx", d.Get("group_id"))
	assert.Equal(t, "1.1.1.1", d.Get("source_content"))
	assert.Equal(t, "ip", d.Get("source_type"))
	assert.Equal(t, "source-name", d.Get("source_name"))
	assert.Equal(t, "2.2.2.2", d.Get("target_content"))
	assert.Equal(t, "ip", d.Get("target_type"))
	assert.Equal(t, "target-name", d.Get("target_name"))
	assert.Equal(t, "80", d.Get("port"))
	assert.Equal(t, "TCP", d.Get("protocol"))
	assert.Equal(t, "accept-pass", d.Get("rule_action"))
	assert.Equal(t, "remark.", d.Get("description"))
	assert.Equal(t, "all", d.Get("scope"))
	assert.Equal(t, "pt-xxxx", d.Get("param_template_id"))
	assert.Equal(t, "pt-name", d.Get("param_template_name"))
	assert.Equal(t, "china", d.Get("country_name"))
	assert.Equal(t, "shenzhen", d.Get("city_name"))
	assert.Equal(t, 1, d.Get("country_code"))
	assert.Equal(t, 2, d.Get("city_code"))
	assert.Equal(t, "CN", d.Get("country_key"))
	assert.Equal(t, "SZ", d.Get("city_key"))
	assert.Equal(t, 1, d.Get("is_region"))
	assert.Equal(t, "tencent", d.Get("cloud_code"))
	assert.Equal(t, 1, d.Get("is_cloud"))
	assert.Equal(t, "ins-xxxx", d.Get("instance_name"))
	assert.Equal(t, 0, d.Get("invalid"))
	assert.Equal(t, 3, d.Get("dns_parse_cnt"))
	assert.Equal(t, 9, d.Get("hit_count"))
	assert.Equal(t, "2025-01-03 00:00:00", d.Get("last_hit_time"))
	assert.Equal(t, "2025-01-01 00:00:00", d.Get("create_time"))
	assert.Equal(t, "2025-01-02 00:00:00", d.Get("update_time"))

	belongMember, ok := d.Get("belong_member").([]interface{})
	assert.True(t, ok)
	assert.Len(t, belongMember, 1)
	member := belongMember[0].(map[string]interface{})
	assert.Equal(t, "1250000000", member["app_id"])
	assert.Equal(t, "100000000001", member["uin"])
	assert.Equal(t, "member-nick", member["nickname"])
	assert.Equal(t, "mem-xxxx", member["member_id"])
	assert.Equal(t, "dept", member["node_name"])
}

// TestFwmEdgeAclRule_Read_MinimalResponse verifies Read does not crash and keeps
// the state untouched when the response only carries the identity fields.
func TestFwmEdgeAclRule_Read_MinimalResponse(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseFwmV20250611Client", fwmClient)
	mockDescribeEdgeAclRules(t, patches, fwmClient, &fwmv20250611.EdgeAclRuleResp{
		RuleId: ptrString("cfw-rule-xxxx"),
	})

	meta := newMockMeta()
	res := fwm.ResourceTencentCloudFwmEdgeAclRule()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"group_id":  "cfw-group-xxxx",
		"direction": 1,
	})
	d.SetId(edgeAclRuleId)

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, edgeAclRuleId, d.Id())
	assert.Equal(t, "cfw-rule-xxxx", d.Get("rule_id"))
	assert.Len(t, d.Get("belong_member").([]interface{}), 0)
}

// TestFwmEdgeAclRule_Read_APIError verifies a non ResourceNotFound error is
// surfaced instead of clearing the id.
func TestFwmEdgeAclRule_Read_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseFwmV20250611Client", fwmClient)
	patches.ApplyMethodFunc(fwmClient, "DescribeEdgeAclRulesWithContext", func(_ context.Context, _ *fwmv20250611.DescribeEdgeAclRulesRequest) (*fwmv20250611.DescribeEdgeAclRulesResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=InternalError, Message=internal error")
	})

	meta := newMockMeta()
	res := fwm.ResourceTencentCloudFwmEdgeAclRule()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"group_id":  "cfw-group-xxxx",
		"direction": 1,
	})
	d.SetId(edgeAclRuleId)

	err := res.Read(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "InternalError")
	assert.Equal(t, edgeAclRuleId, d.Id())
}

// TestFwmEdgeAclRule_Read_ResourceNotFound verifies a ResourceNotFound SDK error
// clears the id so that the rule is recreated.
func TestFwmEdgeAclRule_Read_ResourceNotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseFwmV20250611Client", fwmClient)
	patches.ApplyMethodFunc(fwmClient, "DescribeEdgeAclRulesWithContext", func(_ context.Context, _ *fwmv20250611.DescribeEdgeAclRulesRequest) (*fwmv20250611.DescribeEdgeAclRulesResponse, error) {
		return nil, &sdkErrors.TencentCloudSDKError{Code: "ResourceNotFound"}
	})

	meta := newMockMeta()
	res := fwm.ResourceTencentCloudFwmEdgeAclRule()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"group_id":  "cfw-group-xxxx",
		"direction": 1,
	})
	d.SetId(edgeAclRuleId)

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}

// TestFwmEdgeAclRule_Read_UnexpectedRule verifies Read clears the id when the
// describe result does not match the rule id from the resource id.
func TestFwmEdgeAclRule_Read_UnexpectedRule(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseFwmV20250611Client", fwmClient)
	other := buildEdgeAclRuleResp()
	other.RuleId = ptrString("cfw-rule-other")
	mockDescribeEdgeAclRules(t, patches, fwmClient, other)

	meta := newMockMeta()
	res := fwm.ResourceTencentCloudFwmEdgeAclRule()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"group_id":  "cfw-group-xxxx",
		"direction": 1,
	})
	d.SetId(edgeAclRuleId)

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}

// TestFwmEdgeAclRule_Read_NotFound verifies Read clears id when the rule is gone.
func TestFwmEdgeAclRule_Read_NotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseFwmV20250611Client", fwmClient)
	mockDescribeEdgeAclRules(t, patches, fwmClient)

	meta := newMockMeta()
	res := fwm.ResourceTencentCloudFwmEdgeAclRule()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"group_id":  "cfw-group-xxxx",
		"direction": 1,
	})
	d.SetId(edgeAclRuleId)

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}

// TestFwmEdgeAclRule_Read_LegacyId verifies the old 2-segment id is rejected.
func TestFwmEdgeAclRule_Read_LegacyId(t *testing.T) {
	meta := newMockMeta()
	res := fwm.ResourceTencentCloudFwmEdgeAclRule()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"group_id":  "cfw-group-xxxx",
		"direction": 1,
	})
	d.SetId("cfw-group-xxxx#cfw-rule-xxxx")

	err := res.Read(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "id is broken")
}

// TestFwmEdgeAclRule_Read_InvalidDirection verifies a non-numeric direction
// segment is reported instead of being silently ignored.
func TestFwmEdgeAclRule_Read_InvalidDirection(t *testing.T) {
	meta := newMockMeta()
	res := fwm.ResourceTencentCloudFwmEdgeAclRule()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"group_id":  "cfw-group-xxxx",
		"direction": 1,
	})
	d.SetId("cfw-group-xxxx#cfw-rule-xxxx#ingress")

	err := res.Read(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "id is broken")
}

// TestFwmEdgeAclRule_Update verifies Update parses the 3-segment id and sends
// the mutable fields.
func TestFwmEdgeAclRule_Update(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseFwmV20250611Client", fwmClient)

	var capturedRequest *fwmv20250611.ModifyEdgeAclRuleRequest
	patches.ApplyMethodFunc(fwmClient, "ModifyEdgeAclRuleWithContext", func(_ context.Context, request *fwmv20250611.ModifyEdgeAclRuleRequest) (*fwmv20250611.ModifyEdgeAclRuleResponse, error) {
		capturedRequest = request
		resp := fwmv20250611.NewModifyEdgeAclRuleResponse()
		resp.Response = &fwmv20250611.ModifyEdgeAclRuleResponseParams{
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	mockDescribeEdgeAclRules(t, patches, fwmClient, func() *fwmv20250611.EdgeAclRuleResp {
		rule := buildEdgeAclRuleResp()
		rule.Port = ptrString("8080")
		rule.RuleAction = ptrString("drop-block")
		rule.Description = ptrString("updated")
		return rule
	}())

	meta := newMockMeta()
	res := fwm.ResourceTencentCloudFwmEdgeAclRule()

	oldData := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"group_id":       "cfw-group-xxxx",
		"direction":      1,
		"source_content": "1.1.1.1",
		"source_type":    "ip",
		"target_content": "2.2.2.2",
		"target_type":    "ip",
		"port":           "80",
		"protocol":       "TCP",
		"rule_action":    "accept-pass",
		"scope":          "all",
	})
	oldData.SetId(edgeAclRuleId)
	state := oldData.State()

	rawConfig := terraform.NewResourceConfigRaw(map[string]interface{}{
		"group_id":       "cfw-group-xxxx",
		"direction":      1,
		"source_content": "1.1.1.1",
		"source_type":    "ip",
		"target_content": "2.2.2.2",
		"target_type":    "ip",
		"port":           "8080",
		"protocol":       "TCP",
		"rule_action":    "drop-block",
		"description":    "updated",
		"scope":          "all",
	})
	diff, err := res.Diff(nil, state, rawConfig, nil)
	assert.NoError(t, err)
	d, err := schema.InternalMap(res.Schema).Data(state, diff)
	assert.NoError(t, err)

	err = res.Update(d, meta)
	assert.NoError(t, err)

	assert.NotNil(t, capturedRequest)
	assert.Equal(t, "cfw-group-xxxx", *capturedRequest.GroupId)
	assert.NotNil(t, capturedRequest.Rule)
	assert.Equal(t, "cfw-rule-xxxx", *capturedRequest.Rule.RuleId)
	assert.Equal(t, "8080", *capturedRequest.Rule.Port)
	assert.Equal(t, "drop-block", *capturedRequest.Rule.RuleAction)
	assert.Equal(t, "updated", *capturedRequest.Rule.Description)
}

// TestFwmEdgeAclRule_Delete verifies Delete parses the 3-segment id.
func TestFwmEdgeAclRule_Delete(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseFwmV20250611Client", fwmClient)

	patches.ApplyMethodFunc(fwmClient, "DeleteEdgeAclRuleWithContext", func(_ context.Context, request *fwmv20250611.DeleteEdgeAclRuleRequest) (*fwmv20250611.DeleteEdgeAclRuleResponse, error) {
		assert.Equal(t, "cfw-group-xxxx", *request.GroupId)
		assert.Len(t, request.RuleIds, 1)
		assert.Equal(t, "cfw-rule-xxxx", *request.RuleIds[0])

		resp := fwmv20250611.NewDeleteEdgeAclRuleResponse()
		resp.Response = &fwmv20250611.DeleteEdgeAclRuleResponseParams{
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := fwm.ResourceTencentCloudFwmEdgeAclRule()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
	d.SetId(edgeAclRuleId)

	err := res.Delete(d, meta)
	assert.NoError(t, err)
}

// TestFwmEdgeAclRule_Delete_NotFound verifies Delete is idempotent on ResourceNotFound.
func TestFwmEdgeAclRule_Delete_NotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseFwmV20250611Client", fwmClient)

	patches.ApplyMethodFunc(fwmClient, "DeleteEdgeAclRuleWithContext", func(_ context.Context, _ *fwmv20250611.DeleteEdgeAclRuleRequest) (*fwmv20250611.DeleteEdgeAclRuleResponse, error) {
		return nil, &sdkErrors.TencentCloudSDKError{Code: "ResourceNotFound", Message: "rule not found"}
	})

	meta := newMockMeta()
	res := fwm.ResourceTencentCloudFwmEdgeAclRule()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
	d.SetId(edgeAclRuleId)

	err := res.Delete(d, meta)
	assert.NoError(t, err)
}

// TestFwmEdgeAclRule_Update_SourceType verifies source_type is now sent through
// ModifyEdgeAclRule instead of forcing a replacement.
func TestFwmEdgeAclRule_Update_SourceType(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseFwmV20250611Client", fwmClient)

	var capturedRequest *fwmv20250611.ModifyEdgeAclRuleRequest
	patches.ApplyMethodFunc(fwmClient, "ModifyEdgeAclRuleWithContext", func(_ context.Context, request *fwmv20250611.ModifyEdgeAclRuleRequest) (*fwmv20250611.ModifyEdgeAclRuleResponse, error) {
		capturedRequest = request
		resp := fwmv20250611.NewModifyEdgeAclRuleResponse()
		resp.Response = &fwmv20250611.ModifyEdgeAclRuleResponseParams{RequestId: ptrString("fake-request-id")}
		return resp, nil
	})

	updatedRule := buildEdgeAclRuleResp()
	updatedRule.SourceType = ptrString("domain")
	updatedRule.SourceContent = ptrString("a.example.com")
	mockDescribeEdgeAclRules(t, patches, fwmClient, updatedRule)

	meta := newMockMeta()
	res := fwm.ResourceTencentCloudFwmEdgeAclRule()

	oldData := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"group_id":       "cfw-group-xxxx",
		"direction":      1,
		"source_content": "1.1.1.1",
		"source_type":    "ip",
		"target_content": "2.2.2.2",
		"target_type":    "ip",
	})
	oldData.SetId(edgeAclRuleId)
	state := oldData.State()

	rawConfig := terraform.NewResourceConfigRaw(map[string]interface{}{
		"group_id":       "cfw-group-xxxx",
		"direction":      1,
		"source_content": "a.example.com",
		"source_type":    "domain",
		"target_content": "2.2.2.2",
		"target_type":    "ip",
	})
	diff, err := res.Diff(nil, state, rawConfig, nil)
	assert.NoError(t, err)
	assert.False(t, diff.RequiresNew(), "source_type change must be applied in place")
	d, err := schema.InternalMap(res.Schema).Data(state, diff)
	assert.NoError(t, err)

	err = res.Update(d, meta)
	assert.NoError(t, err)

	assert.NotNil(t, capturedRequest)
	assert.Equal(t, "domain", *capturedRequest.Rule.SourceType)
	assert.Equal(t, "a.example.com", *capturedRequest.Rule.SourceContent)
	assert.Equal(t, "domain", d.Get("source_type"))
}

// TestFwmEdgeAclRule_Update_DirectionRequiresNew verifies a direction change is
// planned as a replacement instead of an in-place update.
func TestFwmEdgeAclRule_Update_DirectionRequiresNew(t *testing.T) {
	res := fwm.ResourceTencentCloudFwmEdgeAclRule()

	oldData := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"group_id":  "cfw-group-xxxx",
		"direction": 1,
	})
	oldData.SetId(edgeAclRuleId)
	state := oldData.State()

	rawConfig := terraform.NewResourceConfigRaw(map[string]interface{}{
		"group_id":  "cfw-group-xxxx",
		"direction": 0,
	})
	diff, err := res.Diff(nil, state, rawConfig, nil)
	assert.NoError(t, err)
	assert.NotNil(t, diff)
	assert.True(t, diff.RequiresNew(), "direction is part of the resource id and must force a new resource")
}

// TestFwmEdgeAclRule_Update_TypeOnlyWithMemberRequiresNew verifies an isolated
// address type change with belong_member_id configured forces a replacement.
func TestFwmEdgeAclRule_Update_TypeOnlyWithMemberRequiresNew(t *testing.T) {
	res := fwm.ResourceTencentCloudFwmEdgeAclRule()

	oldData := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"group_id":         "cfw-group-xxxx",
		"direction":        1,
		"source_content":   "1.1.1.1",
		"source_type":      "ip",
		"target_content":   "2.2.2.2",
		"target_type":      "ip",
		"belong_member_id": "mem-xxxx",
	})
	oldData.SetId(edgeAclRuleId)
	state := oldData.State()

	rawConfig := terraform.NewResourceConfigRaw(map[string]interface{}{
		"group_id":         "cfw-group-xxxx",
		"direction":        1,
		"source_content":   "1.1.1.1",
		"source_type":      "instance",
		"target_content":   "2.2.2.2",
		"target_type":      "ip",
		"belong_member_id": "mem-xxxx",
	})
	diff, err := res.Diff(nil, state, rawConfig, nil)
	assert.NoError(t, err)
	assert.NotNil(t, diff)
	assert.True(t, diff.RequiresNew())
}

// TestFwmEdgeAclRule_Update_ImmutableArgs verifies Update refuses to modify the
// composite id parts.
func TestFwmEdgeAclRule_Update_ImmutableArgs(t *testing.T) {
	meta := newMockMeta()
	res := fwm.ResourceTencentCloudFwmEdgeAclRule()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"group_id":  "cfw-group-yyyy",
		"direction": 1,
	})
	d.SetId(edgeAclRuleId)

	err := res.Update(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "immutable")
}

// TestFwmEdgeAclRule_ImportState verifies the imported composite id is kept
// as-is and can be resolved by Read without touching the API.
func TestFwmEdgeAclRule_ImportState(t *testing.T) {
	res := fwm.ResourceTencentCloudFwmEdgeAclRule()
	assert.NotNil(t, res.Importer)
	assert.NotNil(t, res.Importer.State)

	d := res.TestResourceData()
	d.SetId(edgeAclRuleId)

	imported, err := res.Importer.State(d, nil)
	assert.NoError(t, err)
	assert.Len(t, imported, 1)
	assert.Equal(t, edgeAclRuleId, imported[0].Id())
}

// TestFwmEdgeAclRule_ParseId covers the composite id split failures that an
// import must report instead of silently accepting a wrong id.
// Only failing ids are asserted here, a well formed id would reach the API.
func TestFwmEdgeAclRule_ParseId(t *testing.T) {
	res := fwm.ResourceTencentCloudFwmEdgeAclRule()

	cases := []struct {
		name string
		id   string
	}{
		{name: "missing direction", id: "cfw-group-xxxx#cfw-rule-xxxx"},
		{name: "too many segments", id: "cfw-group-xxxx#cfw-rule-xxxx#1#extra"},
		{name: "non numeric direction", id: "cfw-group-xxxx#cfw-rule-xxxx#ingress"},
	}

	for _, c := range cases {
		d := res.TestResourceData()
		d.SetId(c.id)
		err := res.Read(d, newMockMeta())
		if err == nil || !strings.Contains(err.Error(), "id is broken") {
			t.Errorf("%s: expected an `id is broken` error, got %v", c.name, err)
		}
	}
}
