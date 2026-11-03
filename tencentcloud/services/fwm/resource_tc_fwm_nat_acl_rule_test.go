package fwm_test

import (
	"context"
	"strconv"
	"strings"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	fwmv20250611 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/fwm/v20250611"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/fwm"
)

const (
	testFwmNatGroupId   = "fwmrg-abc"
	testFwmNatRuleId    = "fwmrule-xyz"
	testFwmNatDirection = int64(0)
)

type mockMetaFwmNatAclRule struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaFwmNatAclRule) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaFwmNatAclRule{}

func newMockMetaFwmNatAclRule() *mockMetaFwmNatAclRule {
	return &mockMetaFwmNatAclRule{client: &connectivity.TencentCloudClient{}}
}

func ptrStringFwmNatAcl(s string) *string {
	return &s
}

func ptrInt64FwmNatAcl(i int64) *int64 {
	return &i
}

// fwmNatAclRuleStateId builds the expected resource id: <groupId>#<ruleId>#<direction>.
func fwmNatAclRuleStateId(groupId, ruleId string, direction int64) string {
	return strings.Join([]string{groupId, ruleId, strconv.FormatInt(direction, 10)}, tccommon.FILED_SP)
}

func baseFwmNatAclRuleInput() map[string]interface{} {
	return map[string]interface{}{
		"group_id":       testFwmNatGroupId,
		"direction":      int(testFwmNatDirection),
		"source_content": "1.1.1.1",
		"source_type":    "ip",
		"target_content": "0.0.0.0/0",
		"target_type":    "ip",
		"protocol":       "TCP",
		"rule_action":    "accept",
		"scope":          "ALL",
		"order_index":    1,
		"port":           "80",
		"description":    "tf example nat acl rule",
	}
}

func newFwmNatAclRuleRespFull() *fwmv20250611.NatAclRuleResp {
	resp := newFwmNatAclRuleResp()
	resp.SourceName = ptrStringFwmNatAcl("src-asset")
	resp.TargetName = ptrStringFwmNatAcl("dst-asset")
	resp.ScopeDesc = ptrStringFwmNatAcl("global")
	resp.FwInsId = ptrStringFwmNatAcl("cfwnat-abc")
	resp.CountryName = ptrStringFwmNatAcl("China")
	resp.CityName = ptrStringFwmNatAcl("Shenzhen")
	resp.CountryCode = ptrInt64FwmNatAcl(86)
	resp.CityCode = ptrInt64FwmNatAcl(755)
	resp.CountryKey = ptrStringFwmNatAcl("CN")
	resp.CityKey = ptrStringFwmNatAcl("SZ")
	resp.IsRegion = ptrInt64FwmNatAcl(0)
	resp.CloudCode = ptrStringFwmNatAcl("tencent")
	resp.IsCloud = ptrInt64FwmNatAcl(0)
	resp.ParamTemplateId = ptrStringFwmNatAcl("pt-abc")
	resp.ParamTemplateName = ptrStringFwmNatAcl("web-ports")
	resp.Invalid = ptrInt64FwmNatAcl(0)
	resp.CreateTime = ptrStringFwmNatAcl("2025-01-01 00:00:00")
	resp.UpdateTime = ptrStringFwmNatAcl("2025-01-02 00:00:00")
	resp.DnsParseCnt = ptrInt64FwmNatAcl(3)
	resp.HitCount = ptrInt64FwmNatAcl(12)
	resp.LastHitTime = ptrStringFwmNatAcl("2025-01-02 10:00:00")
	resp.BelongMember = &fwmv20250611.MemberInfo{
		AppId:    ptrStringFwmNatAcl("1300000000"),
		Uin:      ptrStringFwmNatAcl("100000000000"),
		Nickname: ptrStringFwmNatAcl("member-nick"),
		MemberId: ptrStringFwmNatAcl("mem-abc"),
		NodeName: ptrStringFwmNatAcl("dept-a"),
	}

	return resp
}

func newFwmNatAclRuleResp() *fwmv20250611.NatAclRuleResp {
	return &fwmv20250611.NatAclRuleResp{
		RuleId:        ptrStringFwmNatAcl(testFwmNatRuleId),
		Sequence:      ptrInt64FwmNatAcl(1),
		Direction:     ptrInt64FwmNatAcl(testFwmNatDirection),
		SourceContent: ptrStringFwmNatAcl("1.1.1.1"),
		SourceType:    ptrStringFwmNatAcl("ip"),
		TargetContent: ptrStringFwmNatAcl("0.0.0.0/0"),
		TargetType:    ptrStringFwmNatAcl("ip"),
		Port:          ptrStringFwmNatAcl("80"),
		Protocol:      ptrStringFwmNatAcl("TCP"),
		RuleAction:    ptrStringFwmNatAcl("accept"),
		Description:   ptrStringFwmNatAcl("tf example nat acl rule"),
		Scope:         ptrStringFwmNatAcl("ALL"),
	}
}

func TestFwmNatAclRule_Create(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaFwmNatAclRule().client, "UseFwmV20250611Client", fwmClient)

	patches.ApplyMethodFunc(fwmClient, "CreateNatAclRuleWithContext", func(_ context.Context, request *fwmv20250611.CreateNatAclRuleRequest) (*fwmv20250611.CreateNatAclRuleResponse, error) {
		assert.NotNil(t, request.GroupId)
		assert.Equal(t, testFwmNatGroupId, *request.GroupId)
		assert.Equal(t, 1, len(request.Rules))
		rule := request.Rules[0]
		assert.Equal(t, "1.1.1.1", *rule.SourceContent)
		assert.Equal(t, "ip", *rule.SourceType)
		assert.Equal(t, "ip", *rule.TargetType)
		assert.Equal(t, int64(1), *rule.OrderIndex)
		assert.Equal(t, testFwmNatDirection, *rule.Direction)
		assert.Equal(t, "80", *rule.Port)
		assert.Nil(t, rule.RuleId)

		resp := fwmv20250611.NewCreateNatAclRuleResponse()
		resp.Response = &fwmv20250611.CreateNatAclRuleResponseParams{
			RuleIds:   []*string{ptrStringFwmNatAcl(testFwmNatRuleId)},
			RequestId: ptrStringFwmNatAcl("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(fwmClient, "DescribeNatAclRulesWithContext", func(_ context.Context, request *fwmv20250611.DescribeNatAclRulesRequest) (*fwmv20250611.DescribeNatAclRulesResponse, error) {
		assert.NotNil(t, request.GroupId)
		assert.Equal(t, testFwmNatGroupId, *request.GroupId)
		assert.NotNil(t, request.Direction)
		assert.Equal(t, testFwmNatDirection, *request.Direction)
		assert.NotNil(t, request.Filters)
		assert.Equal(t, 1, len(request.Filters))

		resp := fwmv20250611.NewDescribeNatAclRulesResponse()
		resp.Response = &fwmv20250611.DescribeNatAclRulesResponseParams{
			Rules: []*fwmv20250611.NatAclRuleResp{newFwmNatAclRuleResp()},
		}
		return resp, nil
	})

	meta := newMockMetaFwmNatAclRule()
	res := fwm.ResourceTencentCloudFwmNatAclRule()
	d := schema.TestResourceDataRaw(t, res.Schema, baseFwmNatAclRuleInput())

	err := res.Create(d, meta)
	assert.NoError(t, err)
	// resource id is <groupId>#<ruleId>#<direction>
	assert.Equal(t, fwmNatAclRuleStateId(testFwmNatGroupId, testFwmNatRuleId, testFwmNatDirection), d.Id())
	assert.Equal(t, testFwmNatRuleId, d.Get("rule_id"))
	assert.Equal(t, 1, d.Get("sequence"))
}

func TestFwmNatAclRule_Create_InboundDirection(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaFwmNatAclRule().client, "UseFwmV20250611Client", fwmClient)

	patches.ApplyMethodFunc(fwmClient, "CreateNatAclRuleWithContext", func(_ context.Context, request *fwmv20250611.CreateNatAclRuleRequest) (*fwmv20250611.CreateNatAclRuleResponse, error) {
		assert.Equal(t, int64(1), *request.Rules[0].Direction)

		resp := fwmv20250611.NewCreateNatAclRuleResponse()
		resp.Response = &fwmv20250611.CreateNatAclRuleResponseParams{
			RuleIds:   []*string{ptrStringFwmNatAcl(testFwmNatRuleId)},
			RequestId: ptrStringFwmNatAcl("fake-request-id-inbound"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(fwmClient, "DescribeNatAclRulesWithContext", func(_ context.Context, request *fwmv20250611.DescribeNatAclRulesRequest) (*fwmv20250611.DescribeNatAclRulesResponse, error) {
		assert.Equal(t, int64(1), *request.Direction)

		respRule := newFwmNatAclRuleResp()
		respRule.Direction = ptrInt64FwmNatAcl(1)
		resp := fwmv20250611.NewDescribeNatAclRulesResponse()
		resp.Response = &fwmv20250611.DescribeNatAclRulesResponseParams{
			Rules: []*fwmv20250611.NatAclRuleResp{respRule},
		}
		return resp, nil
	})

	meta := newMockMetaFwmNatAclRule()
	res := fwm.ResourceTencentCloudFwmNatAclRule()
	input := baseFwmNatAclRuleInput()
	input["direction"] = 1
	d := schema.TestResourceDataRaw(t, res.Schema, input)

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, testFwmNatGroupId+"#"+testFwmNatRuleId+"#1", d.Id())
}

func TestFwmNatAclRule_Read(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaFwmNatAclRule().client, "UseFwmV20250611Client", fwmClient)

	patches.ApplyMethodFunc(fwmClient, "DescribeNatAclRulesWithContext", func(_ context.Context, request *fwmv20250611.DescribeNatAclRulesRequest) (*fwmv20250611.DescribeNatAclRulesResponse, error) {
		assert.Equal(t, testFwmNatGroupId, *request.GroupId)
		// direction comes from the resource id, not from state.
		assert.Equal(t, int64(1), *request.Direction)
		assert.Equal(t, int64(0), *request.Offset)
		assert.Equal(t, int64(100), *request.Limit)
		assert.Equal(t, 1, len(request.Filters))
		assert.Equal(t, "RuleId", *request.Filters[0].Name)
		assert.Equal(t, testFwmNatRuleId, *request.Filters[0].Values[0])

		respRule := newFwmNatAclRuleResp()
		respRule.Direction = ptrInt64FwmNatAcl(1)
		resp := fwmv20250611.NewDescribeNatAclRulesResponse()
		resp.Response = &fwmv20250611.DescribeNatAclRulesResponseParams{
			Rules: []*fwmv20250611.NatAclRuleResp{respRule},
		}
		return resp, nil
	})

	meta := newMockMetaFwmNatAclRule()
	res := fwm.ResourceTencentCloudFwmNatAclRule()
	d := schema.TestResourceDataRaw(t, res.Schema, baseFwmNatAclRuleInput())
	d.SetId(testFwmNatGroupId + "#" + testFwmNatRuleId + "#1")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, testFwmNatGroupId+"#"+testFwmNatRuleId+"#1", d.Id())
	assert.Equal(t, testFwmNatRuleId, d.Get("rule_id"))
	assert.Equal(t, 1, d.Get("direction"))
	assert.Equal(t, 1, d.Get("sequence"))
}

func TestFwmNatAclRule_ReadBrokenId(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaFwmNatAclRule().client, "UseFwmV20250611Client", fwmClient)

	called := false
	patches.ApplyMethodFunc(fwmClient, "DescribeNatAclRulesWithContext", func(_ context.Context, request *fwmv20250611.DescribeNatAclRulesRequest) (*fwmv20250611.DescribeNatAclRulesResponse, error) {
		called = true
		return fwmv20250611.NewDescribeNatAclRulesResponse(), nil
	})

	meta := newMockMetaFwmNatAclRule()
	res := fwm.ResourceTencentCloudFwmNatAclRule()

	// two segments only.
	d := schema.TestResourceDataRaw(t, res.Schema, baseFwmNatAclRuleInput())
	d.SetId(testFwmNatGroupId + tccommon.FILED_SP + testFwmNatRuleId)
	assert.Error(t, res.Read(d, meta))

	// direction is not a number.
	d2 := schema.TestResourceDataRaw(t, res.Schema, baseFwmNatAclRuleInput())
	d2.SetId(testFwmNatGroupId + tccommon.FILED_SP + testFwmNatRuleId + tccommon.FILED_SP + "inbound")
	assert.Error(t, res.Read(d2, meta))

	assert.False(t, called)
}

// TestFwmNatAclRule_ReadRetryableError verifies the describe call is retried on a retryable SDK error.
func TestFwmNatAclRule_ReadRetryableError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaFwmNatAclRule().client, "UseFwmV20250611Client", fwmClient)

	callCount := 0
	patches.ApplyMethodFunc(fwmClient, "DescribeNatAclRulesWithContext", func(_ context.Context, request *fwmv20250611.DescribeNatAclRulesRequest) (*fwmv20250611.DescribeNatAclRulesResponse, error) {
		callCount++
		if callCount == 1 {
			return nil, &sdkErrors.TencentCloudSDKError{Code: "RequestLimitExceeded", Message: "too many requests"}
		}

		resp := fwmv20250611.NewDescribeNatAclRulesResponse()
		resp.Response = &fwmv20250611.DescribeNatAclRulesResponseParams{
			Rules: []*fwmv20250611.NatAclRuleResp{newFwmNatAclRuleResp()},
		}
		return resp, nil
	})

	meta := newMockMetaFwmNatAclRule()
	res := fwm.ResourceTencentCloudFwmNatAclRule()
	d := schema.TestResourceDataRaw(t, res.Schema, baseFwmNatAclRuleInput())
	d.SetId(fwmNatAclRuleStateId(testFwmNatGroupId, testFwmNatRuleId, testFwmNatDirection))

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, 2, callCount)
	assert.Equal(t, testFwmNatRuleId, d.Get("rule_id"))
}

// TestFwmNatAclRule_ReadNilResponse verifies a nil describe response surfaces an error instead of being treated as deleted.
func TestFwmNatAclRule_ReadNilResponse(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaFwmNatAclRule().client, "UseFwmV20250611Client", fwmClient)

	patches.ApplyMethodFunc(fwmClient, "DescribeNatAclRulesWithContext", func(_ context.Context, request *fwmv20250611.DescribeNatAclRulesRequest) (*fwmv20250611.DescribeNatAclRulesResponse, error) {
		return &fwmv20250611.DescribeNatAclRulesResponse{}, nil
	})

	meta := newMockMetaFwmNatAclRule()
	res := fwm.ResourceTencentCloudFwmNatAclRule()
	d := schema.TestResourceDataRaw(t, res.Schema, baseFwmNatAclRuleInput())
	d.SetId(fwmNatAclRuleStateId(testFwmNatGroupId, testFwmNatRuleId, testFwmNatDirection))

	err := res.Read(d, meta)
	assert.Error(t, err)
	assert.Equal(t, fwmNatAclRuleStateId(testFwmNatGroupId, testFwmNatRuleId, testFwmNatDirection), d.Id())
}

func TestFwmNatAclRule_ReadNotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaFwmNatAclRule().client, "UseFwmV20250611Client", fwmClient)

	patches.ApplyMethodFunc(fwmClient, "DescribeNatAclRulesWithContext", func(_ context.Context, request *fwmv20250611.DescribeNatAclRulesRequest) (*fwmv20250611.DescribeNatAclRulesResponse, error) {
		resp := fwmv20250611.NewDescribeNatAclRulesResponse()
		resp.Response = &fwmv20250611.DescribeNatAclRulesResponseParams{
			Rules: []*fwmv20250611.NatAclRuleResp{},
		}
		return resp, nil
	})

	meta := newMockMetaFwmNatAclRule()
	res := fwm.ResourceTencentCloudFwmNatAclRule()
	d := schema.TestResourceDataRaw(t, res.Schema, baseFwmNatAclRuleInput())
	d.SetId(fwmNatAclRuleStateId(testFwmNatGroupId, testFwmNatRuleId, testFwmNatDirection))

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}

func TestFwmNatAclRule_Update(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaFwmNatAclRule().client, "UseFwmV20250611Client", fwmClient)

	patches.ApplyMethodFunc(fwmClient, "ModifyNatAclRuleWithContext", func(_ context.Context, request *fwmv20250611.ModifyNatAclRuleRequest) (*fwmv20250611.ModifyNatAclRuleResponse, error) {
		assert.NotNil(t, request.GroupId)
		assert.Equal(t, testFwmNatGroupId, *request.GroupId)
		assert.NotNil(t, request.Rule)
		assert.NotNil(t, request.Rule.RuleId)
		assert.Equal(t, testFwmNatRuleId, *request.Rule.RuleId)
		assert.Equal(t, "ap-guangzhou", *request.Rule.Scope)
		// direction is taken from the resource id.
		assert.Equal(t, int64(1), *request.Rule.Direction)

		resp := fwmv20250611.NewModifyNatAclRuleResponse()
		resp.Response = &fwmv20250611.ModifyNatAclRuleResponseParams{
			RequestId: ptrStringFwmNatAcl("fake-request-id-update"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(fwmClient, "DescribeNatAclRulesWithContext", func(_ context.Context, request *fwmv20250611.DescribeNatAclRulesRequest) (*fwmv20250611.DescribeNatAclRulesResponse, error) {
		respRule := newFwmNatAclRuleResp()
		respRule.Scope = ptrStringFwmNatAcl("ap-guangzhou")
		respRule.Direction = ptrInt64FwmNatAcl(1)
		resp := fwmv20250611.NewDescribeNatAclRulesResponse()
		resp.Response = &fwmv20250611.DescribeNatAclRulesResponseParams{
			Rules: []*fwmv20250611.NatAclRuleResp{respRule},
		}
		return resp, nil
	})

	meta := newMockMetaFwmNatAclRule()
	res := fwm.ResourceTencentCloudFwmNatAclRule()
	input := baseFwmNatAclRuleInput()
	input["scope"] = "ap-guangzhou"
	d := schema.TestResourceDataRaw(t, res.Schema, input)
	d.SetId(testFwmNatGroupId + "#" + testFwmNatRuleId + "#1")

	err := res.Update(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "ap-guangzhou", d.Get("scope"))
}

func TestFwmNatAclRule_UpdateBrokenId(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaFwmNatAclRule().client, "UseFwmV20250611Client", fwmClient)

	res := fwm.ResourceTencentCloudFwmNatAclRule()
	d := schema.TestResourceDataRaw(t, res.Schema, baseFwmNatAclRuleInput())
	d.SetId(testFwmNatGroupId + tccommon.FILED_SP + testFwmNatRuleId)

	assert.Error(t, res.Update(d, newMockMetaFwmNatAclRule()))
}

func TestFwmNatAclRule_Delete(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaFwmNatAclRule().client, "UseFwmV20250611Client", fwmClient)

	patches.ApplyMethodFunc(fwmClient, "DeleteNatAclRuleWithContext", func(_ context.Context, request *fwmv20250611.DeleteNatAclRuleRequest) (*fwmv20250611.DeleteNatAclRuleResponse, error) {
		assert.NotNil(t, request.GroupId)
		assert.Equal(t, testFwmNatGroupId, *request.GroupId)
		assert.NotNil(t, request.RuleIds)
		assert.Equal(t, 1, len(request.RuleIds))
		assert.Equal(t, testFwmNatRuleId, *request.RuleIds[0])

		resp := fwmv20250611.NewDeleteNatAclRuleResponse()
		resp.Response = &fwmv20250611.DeleteNatAclRuleResponseParams{
			RequestId: ptrStringFwmNatAcl("fake-request-id-delete"),
		}
		return resp, nil
	})

	meta := newMockMetaFwmNatAclRule()
	res := fwm.ResourceTencentCloudFwmNatAclRule()
	d := schema.TestResourceDataRaw(t, res.Schema, baseFwmNatAclRuleInput())
	d.SetId(fwmNatAclRuleStateId(testFwmNatGroupId, testFwmNatRuleId, testFwmNatDirection))

	err := res.Delete(d, meta)
	assert.NoError(t, err)
}

func TestFwmNatAclRule_DeleteBrokenId(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaFwmNatAclRule().client, "UseFwmV20250611Client", fwmClient)

	res := fwm.ResourceTencentCloudFwmNatAclRule()
	d := schema.TestResourceDataRaw(t, res.Schema, baseFwmNatAclRuleInput())
	d.SetId(testFwmNatGroupId + tccommon.FILED_SP + testFwmNatRuleId)

	assert.Error(t, res.Delete(d, newMockMetaFwmNatAclRule()))
}

func TestFwmNatAclRule_DeleteResourceNotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaFwmNatAclRule().client, "UseFwmV20250611Client", fwmClient)

	patches.ApplyMethodFunc(fwmClient, "DeleteNatAclRuleWithContext", func(_ context.Context, request *fwmv20250611.DeleteNatAclRuleRequest) (*fwmv20250611.DeleteNatAclRuleResponse, error) {
		return nil, &sdkErrors.TencentCloudSDKError{Code: "ResourceNotFound"}
	})

	meta := newMockMetaFwmNatAclRule()
	res := fwm.ResourceTencentCloudFwmNatAclRule()
	d := schema.TestResourceDataRaw(t, res.Schema, baseFwmNatAclRuleInput())
	d.SetId(fwmNatAclRuleStateId(testFwmNatGroupId, testFwmNatRuleId, testFwmNatDirection))

	err := res.Delete(d, meta)
	assert.NoError(t, err)
}

// TestFwmNatAclRule_CreateOrderIndexAndUserRuleId verifies order_index is sent as the rule priority
// and a user supplied rule_id is passed through on create.
func TestFwmNatAclRule_CreateOrderIndexAndUserRuleId(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaFwmNatAclRule().client, "UseFwmV20250611Client", fwmClient)

	patches.ApplyMethodFunc(fwmClient, "CreateNatAclRuleWithContext", func(_ context.Context, request *fwmv20250611.CreateNatAclRuleRequest) (*fwmv20250611.CreateNatAclRuleResponse, error) {
		assert.Equal(t, 1, len(request.Rules))
		assert.Equal(t, int64(7), *request.Rules[0].OrderIndex)
		assert.NotNil(t, request.Rules[0].RuleId)
		assert.Equal(t, "user-rule-id", *request.Rules[0].RuleId)
		assert.Equal(t, "pt-abc", *request.Rules[0].ParamTemplateId)
		assert.Equal(t, "mem-abc", *request.Rules[0].BelongMemberId)

		resp := fwmv20250611.NewCreateNatAclRuleResponse()
		resp.Response = &fwmv20250611.CreateNatAclRuleResponseParams{
			RuleIds:   []*string{ptrStringFwmNatAcl("user-rule-id")},
			RequestId: ptrStringFwmNatAcl("fake-request-id-order"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(fwmClient, "DescribeNatAclRulesWithContext", func(_ context.Context, request *fwmv20250611.DescribeNatAclRulesRequest) (*fwmv20250611.DescribeNatAclRulesResponse, error) {
		respRule := newFwmNatAclRuleResp()
		respRule.RuleId = ptrStringFwmNatAcl("user-rule-id")
		respRule.Sequence = ptrInt64FwmNatAcl(7)
		resp := fwmv20250611.NewDescribeNatAclRulesResponse()
		resp.Response = &fwmv20250611.DescribeNatAclRulesResponseParams{
			Rules: []*fwmv20250611.NatAclRuleResp{respRule},
		}
		return resp, nil
	})

	meta := newMockMetaFwmNatAclRule()
	res := fwm.ResourceTencentCloudFwmNatAclRule()
	input := baseFwmNatAclRuleInput()
	input["order_index"] = 7
	input["rule_id"] = "user-rule-id"
	input["param_template_id"] = "pt-abc"
	input["belong_member_id"] = "mem-abc"
	d := schema.TestResourceDataRaw(t, res.Schema, input)

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, fwmNatAclRuleStateId(testFwmNatGroupId, "user-rule-id", testFwmNatDirection), d.Id())
	assert.Equal(t, "user-rule-id", d.Get("rule_id"))
	assert.Equal(t, 7, d.Get("order_index"))
	assert.Equal(t, 7, d.Get("sequence"))
}

// TestFwmNatAclRule_CreateEmptyRuleIds verifies an empty RuleIds response never produces an empty id.
func TestFwmNatAclRule_CreateEmptyRuleIds(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaFwmNatAclRule().client, "UseFwmV20250611Client", fwmClient)

	patches.ApplyMethodFunc(fwmClient, "CreateNatAclRuleWithContext", func(_ context.Context, request *fwmv20250611.CreateNatAclRuleRequest) (*fwmv20250611.CreateNatAclRuleResponse, error) {
		resp := fwmv20250611.NewCreateNatAclRuleResponse()
		resp.Response = &fwmv20250611.CreateNatAclRuleResponseParams{
			RuleIds: []*string{},
		}
		return resp, nil
	})

	meta := newMockMetaFwmNatAclRule()
	res := fwm.ResourceTencentCloudFwmNatAclRule()
	d := schema.TestResourceDataRaw(t, res.Schema, baseFwmNatAclRuleInput())

	err := res.Create(d, meta)
	assert.Error(t, err)
	assert.Equal(t, "", d.Id())
}

// TestFwmNatAclRule_CreateMultipleRuleIds verifies a create that returns more than one rule id is rejected
// instead of silently picking the first one.
func TestFwmNatAclRule_CreateMultipleRuleIds(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaFwmNatAclRule().client, "UseFwmV20250611Client", fwmClient)

	patches.ApplyMethodFunc(fwmClient, "CreateNatAclRuleWithContext", func(_ context.Context, request *fwmv20250611.CreateNatAclRuleRequest) (*fwmv20250611.CreateNatAclRuleResponse, error) {
		resp := fwmv20250611.NewCreateNatAclRuleResponse()
		resp.Response = &fwmv20250611.CreateNatAclRuleResponseParams{
			RuleIds: []*string{ptrStringFwmNatAcl("rule-1"), ptrStringFwmNatAcl("rule-2")},
		}
		return resp, nil
	})

	meta := newMockMetaFwmNatAclRule()
	res := fwm.ResourceTencentCloudFwmNatAclRule()
	d := schema.TestResourceDataRaw(t, res.Schema, baseFwmNatAclRuleInput())

	err := res.Create(d, meta)
	assert.Error(t, err)
	assert.Equal(t, "", d.Id())
}

// TestFwmNatAclRule_CreateEmptyRuleIdValue verifies an empty rule id string is rejected.
func TestFwmNatAclRule_CreateEmptyRuleIdValue(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaFwmNatAclRule().client, "UseFwmV20250611Client", fwmClient)

	patches.ApplyMethodFunc(fwmClient, "CreateNatAclRuleWithContext", func(_ context.Context, request *fwmv20250611.CreateNatAclRuleRequest) (*fwmv20250611.CreateNatAclRuleResponse, error) {
		resp := fwmv20250611.NewCreateNatAclRuleResponse()
		resp.Response = &fwmv20250611.CreateNatAclRuleResponseParams{
			RuleIds: []*string{ptrStringFwmNatAcl("")},
		}
		return resp, nil
	})

	meta := newMockMetaFwmNatAclRule()
	res := fwm.ResourceTencentCloudFwmNatAclRule()
	d := schema.TestResourceDataRaw(t, res.Schema, baseFwmNatAclRuleInput())

	err := res.Create(d, meta)
	assert.Error(t, err)
	assert.Equal(t, "", d.Id())
}

// TestFwmNatAclRule_CreateNilResponse verifies a nil create response surfaces an error.
func TestFwmNatAclRule_CreateNilResponse(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaFwmNatAclRule().client, "UseFwmV20250611Client", fwmClient)

	patches.ApplyMethodFunc(fwmClient, "CreateNatAclRuleWithContext", func(_ context.Context, request *fwmv20250611.CreateNatAclRuleRequest) (*fwmv20250611.CreateNatAclRuleResponse, error) {
		return &fwmv20250611.CreateNatAclRuleResponse{}, nil
	})

	meta := newMockMetaFwmNatAclRule()
	res := fwm.ResourceTencentCloudFwmNatAclRule()
	d := schema.TestResourceDataRaw(t, res.Schema, baseFwmNatAclRuleInput())

	assert.Error(t, res.Create(d, meta))
	assert.Equal(t, "", d.Id())
}

// TestFwmNatAclRule_ReadAllFields verifies every readable leaf field coming from DescribeNatAclRules,
// including the nested belong_member block, is written into state.
func TestFwmNatAclRule_ReadAllFields(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaFwmNatAclRule().client, "UseFwmV20250611Client", fwmClient)

	patches.ApplyMethodFunc(fwmClient, "DescribeNatAclRulesWithContext", func(_ context.Context, request *fwmv20250611.DescribeNatAclRulesRequest) (*fwmv20250611.DescribeNatAclRulesResponse, error) {
		resp := fwmv20250611.NewDescribeNatAclRulesResponse()
		resp.Response = &fwmv20250611.DescribeNatAclRulesResponseParams{
			Rules:         []*fwmv20250611.NatAclRuleResp{newFwmNatAclRuleRespFull()},
			TotalCount:    ptrInt64FwmNatAcl(1),
			AllTotalCount: ptrInt64FwmNatAcl(5),
			RequestId:     ptrStringFwmNatAcl("fake-request-id-read-all"),
		}
		return resp, nil
	})

	meta := newMockMetaFwmNatAclRule()
	res := fwm.ResourceTencentCloudFwmNatAclRule()
	d := schema.TestResourceDataRaw(t, res.Schema, baseFwmNatAclRuleInput())
	d.SetId(fwmNatAclRuleStateId(testFwmNatGroupId, testFwmNatRuleId, testFwmNatDirection))

	err := res.Read(d, meta)
	assert.NoError(t, err)

	expects := map[string]interface{}{
		"group_id":            testFwmNatGroupId,
		"rule_id":             testFwmNatRuleId,
		"direction":           int(testFwmNatDirection),
		"source_content":      "1.1.1.1",
		"source_type":         "ip",
		"source_name":         "src-asset",
		"target_content":      "0.0.0.0/0",
		"target_type":         "ip",
		"target_name":         "dst-asset",
		"protocol":            "TCP",
		"rule_action":         "accept",
		"scope":               "ALL",
		"scope_desc":          "global",
		"sequence":            1,
		"order_index":         1,
		"port":                "80",
		"description":         "tf example nat acl rule",
		"param_template_id":   "pt-abc",
		"param_template_name": "web-ports",
		"belong_member_id":    "mem-abc",
		"fw_ins_id":           "cfwnat-abc",
		"country_name":        "China",
		"city_name":           "Shenzhen",
		"country_code":        86,
		"city_code":           755,
		"country_key":         "CN",
		"city_key":            "SZ",
		"is_region":           0,
		"cloud_code":          "tencent",
		"is_cloud":            0,
		"invalid":             0,
		"create_time":         "2025-01-01 00:00:00",
		"update_time":         "2025-01-02 00:00:00",
		"dns_parse_cnt":       3,
		"hit_count":           12,
		"last_hit_time":       "2025-01-02 10:00:00",
	}
	for k, v := range expects {
		assert.Equal(t, v, d.Get(k), "field %s", k)
	}

	raw := d.Get("belong_member")
	members, ok := raw.([]interface{})
	assert.True(t, ok, "belong_member should be a list")
	assert.Equal(t, 1, len(members))
	member, ok := members[0].(map[string]interface{})
	assert.True(t, ok, "belong_member element should be a map")
	assert.Equal(t, "1300000000", member["app_id"])
	assert.Equal(t, "100000000000", member["uin"])
	assert.Equal(t, "member-nick", member["nickname"])
	assert.Equal(t, "mem-abc", member["member_id"])
	assert.Equal(t, "dept-a", member["node_name"])
}

// TestFwmNatAclRule_ReadClearedOptionalFields verifies clearable optional fields are reset to ""
// when the cloud side returns them as null, to avoid a permanent diff.
func TestFwmNatAclRule_ReadClearedOptionalFields(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaFwmNatAclRule().client, "UseFwmV20250611Client", fwmClient)

	patches.ApplyMethodFunc(fwmClient, "DescribeNatAclRulesWithContext", func(_ context.Context, request *fwmv20250611.DescribeNatAclRulesRequest) (*fwmv20250611.DescribeNatAclRulesResponse, error) {
		respRule := newFwmNatAclRuleResp()
		respRule.Port = nil
		respRule.Description = nil
		respRule.ParamTemplateId = nil
		respRule.BelongMember = nil
		respRule.Sequence = nil
		resp := fwmv20250611.NewDescribeNatAclRulesResponse()
		resp.Response = &fwmv20250611.DescribeNatAclRulesResponseParams{
			Rules: []*fwmv20250611.NatAclRuleResp{respRule},
		}
		return resp, nil
	})

	meta := newMockMetaFwmNatAclRule()
	res := fwm.ResourceTencentCloudFwmNatAclRule()
	d := schema.TestResourceDataRaw(t, res.Schema, baseFwmNatAclRuleInput())
	d.SetId(fwmNatAclRuleStateId(testFwmNatGroupId, testFwmNatRuleId, testFwmNatDirection))

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Get("port"))
	assert.Equal(t, "", d.Get("description"))
	assert.Equal(t, "", d.Get("param_template_id"))
	assert.Equal(t, "", d.Get("belong_member_id"))
	// order_index keeps the configured value when the server returns no sequence.
	assert.Equal(t, 1, d.Get("order_index"))

	members, ok := d.Get("belong_member").([]interface{})
	assert.True(t, ok)
	assert.Equal(t, 0, len(members))
}

// TestFwmNatAclRule_ReadResourceNotFoundSDKError verifies a ResourceNotFound SDK error clears the id.
func TestFwmNatAclRule_ReadResourceNotFoundSDKError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaFwmNatAclRule().client, "UseFwmV20250611Client", fwmClient)

	patches.ApplyMethodFunc(fwmClient, "DescribeNatAclRulesWithContext", func(_ context.Context, request *fwmv20250611.DescribeNatAclRulesRequest) (*fwmv20250611.DescribeNatAclRulesResponse, error) {
		return nil, &sdkErrors.TencentCloudSDKError{Code: "ResourceNotFound"}
	})

	meta := newMockMetaFwmNatAclRule()
	res := fwm.ResourceTencentCloudFwmNatAclRule()
	d := schema.TestResourceDataRaw(t, res.Schema, baseFwmNatAclRuleInput())
	d.SetId(fwmNatAclRuleStateId(testFwmNatGroupId, testFwmNatRuleId, testFwmNatDirection))

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}

// TestFwmNatAclRule_ReadOtherSDKError verifies a non ResourceNotFound SDK error is returned as is.
func TestFwmNatAclRule_ReadOtherSDKError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaFwmNatAclRule().client, "UseFwmV20250611Client", fwmClient)

	patches.ApplyMethodFunc(fwmClient, "DescribeNatAclRulesWithContext", func(_ context.Context, request *fwmv20250611.DescribeNatAclRulesRequest) (*fwmv20250611.DescribeNatAclRulesResponse, error) {
		return nil, &sdkErrors.TencentCloudSDKError{Code: "UnauthorizedOperation"}
	})

	meta := newMockMetaFwmNatAclRule()
	res := fwm.ResourceTencentCloudFwmNatAclRule()
	d := schema.TestResourceDataRaw(t, res.Schema, baseFwmNatAclRuleInput())
	d.SetId(fwmNatAclRuleStateId(testFwmNatGroupId, testFwmNatRuleId, testFwmNatDirection))

	err := res.Read(d, meta)
	assert.Error(t, err)
	assert.Equal(t, fwmNatAclRuleStateId(testFwmNatGroupId, testFwmNatRuleId, testFwmNatDirection), d.Id())
}

// TestFwmNatAclRule_UpdateSendsAllWritableFields verifies modify always sends the whole rule,
// including cleared optional fields (empty string) and the configured order_index.
func TestFwmNatAclRule_UpdateSendsAllWritableFields(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaFwmNatAclRule().client, "UseFwmV20250611Client", fwmClient)

	patches.ApplyMethodFunc(fwmClient, "ModifyNatAclRuleWithContext", func(_ context.Context, request *fwmv20250611.ModifyNatAclRuleRequest) (*fwmv20250611.ModifyNatAclRuleResponse, error) {
		assert.Equal(t, testFwmNatGroupId, *request.GroupId)
		assert.NotNil(t, request.Rule)
		assert.Equal(t, testFwmNatRuleId, *request.Rule.RuleId)
		assert.Equal(t, testFwmNatDirection, *request.Rule.Direction)
		assert.Equal(t, int64(3), *request.Rule.OrderIndex)
		assert.Equal(t, "1.1.1.2", *request.Rule.SourceContent)
		assert.Equal(t, "ip", *request.Rule.SourceType)
		assert.Equal(t, "10.0.0.0/8", *request.Rule.TargetContent)
		assert.Equal(t, "ip", *request.Rule.TargetType)
		assert.Equal(t, "UDP", *request.Rule.Protocol)
		assert.Equal(t, "drop", *request.Rule.RuleAction)
		assert.Equal(t, "ALL", *request.Rule.Scope)

		resp := fwmv20250611.NewModifyNatAclRuleResponse()
		resp.Response = &fwmv20250611.ModifyNatAclRuleResponseParams{
			RequestId: ptrStringFwmNatAcl("fake-request-id-update-all"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(fwmClient, "DescribeNatAclRulesWithContext", func(_ context.Context, request *fwmv20250611.DescribeNatAclRulesRequest) (*fwmv20250611.DescribeNatAclRulesResponse, error) {
		respRule := newFwmNatAclRuleResp()
		respRule.SourceContent = ptrStringFwmNatAcl("1.1.1.2")
		respRule.TargetContent = ptrStringFwmNatAcl("10.0.0.0/8")
		respRule.Protocol = ptrStringFwmNatAcl("UDP")
		respRule.RuleAction = ptrStringFwmNatAcl("drop")
		respRule.Sequence = ptrInt64FwmNatAcl(3)
		respRule.Port = nil
		respRule.Description = nil
		resp := fwmv20250611.NewDescribeNatAclRulesResponse()
		resp.Response = &fwmv20250611.DescribeNatAclRulesResponseParams{
			Rules: []*fwmv20250611.NatAclRuleResp{respRule},
		}
		return resp, nil
	})

	meta := newMockMetaFwmNatAclRule()
	res := fwm.ResourceTencentCloudFwmNatAclRule()
	input := map[string]interface{}{
		"group_id":       testFwmNatGroupId,
		"direction":      int(testFwmNatDirection),
		"source_content": "1.1.1.2",
		"source_type":    "ip",
		"target_content": "10.0.0.0/8",
		"target_type":    "ip",
		"protocol":       "UDP",
		"rule_action":    "drop",
		"scope":          "ALL",
		"order_index":    3,
	}
	d := schema.TestResourceDataRaw(t, res.Schema, input)
	d.SetId(fwmNatAclRuleStateId(testFwmNatGroupId, testFwmNatRuleId, testFwmNatDirection))

	err := res.Update(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Get("port"))
	assert.Equal(t, "", d.Get("description"))
	assert.Equal(t, 3, d.Get("order_index"))
}

// TestFwmNatAclRule_UpdateNilResponse verifies a nil modify response surfaces an error.
func TestFwmNatAclRule_UpdateNilResponse(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaFwmNatAclRule().client, "UseFwmV20250611Client", fwmClient)

	patches.ApplyMethodFunc(fwmClient, "ModifyNatAclRuleWithContext", func(_ context.Context, request *fwmv20250611.ModifyNatAclRuleRequest) (*fwmv20250611.ModifyNatAclRuleResponse, error) {
		return &fwmv20250611.ModifyNatAclRuleResponse{}, nil
	})

	meta := newMockMetaFwmNatAclRule()
	res := fwm.ResourceTencentCloudFwmNatAclRule()
	d := schema.TestResourceDataRaw(t, res.Schema, baseFwmNatAclRuleInput())
	d.SetId(fwmNatAclRuleStateId(testFwmNatGroupId, testFwmNatRuleId, testFwmNatDirection))

	assert.Error(t, res.Update(d, meta))
}

// TestFwmNatAclRule_DeleteNilResponse verifies a nil delete response surfaces an error.
func TestFwmNatAclRule_DeleteNilResponse(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	fwmClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(newMockMetaFwmNatAclRule().client, "UseFwmV20250611Client", fwmClient)

	patches.ApplyMethodFunc(fwmClient, "DeleteNatAclRuleWithContext", func(_ context.Context, request *fwmv20250611.DeleteNatAclRuleRequest) (*fwmv20250611.DeleteNatAclRuleResponse, error) {
		return &fwmv20250611.DeleteNatAclRuleResponse{}, nil
	})

	meta := newMockMetaFwmNatAclRule()
	res := fwm.ResourceTencentCloudFwmNatAclRule()
	d := schema.TestResourceDataRaw(t, res.Schema, baseFwmNatAclRuleInput())
	d.SetId(fwmNatAclRuleStateId(testFwmNatGroupId, testFwmNatRuleId, testFwmNatDirection))

	assert.Error(t, res.Delete(d, meta))
}

// TestFwmNatAclRule_SchemaConsistency guards the schema flags that carry the resource contract.
func TestFwmNatAclRule_SchemaConsistency(t *testing.T) {
	res := fwm.ResourceTencentCloudFwmNatAclRule()

	required := []string{
		"group_id", "direction", "source_content", "source_type", "target_content",
		"target_type", "protocol", "rule_action", "order_index", "scope",
	}
	for _, name := range required {
		s := res.Schema[name]
		if s == nil {
			t.Fatalf("schema %s is missing", name)
		}
		assert.True(t, s.Required, "%s should be required", name)
	}

	optional := []string{"port", "description", "param_template_id", "belong_member_id"}
	for _, name := range optional {
		s := res.Schema[name]
		if s == nil {
			t.Fatalf("schema %s is missing", name)
		}
		assert.True(t, s.Optional, "%s should be optional", name)
	}

	for _, name := range []string{"group_id", "direction", "rule_id"} {
		assert.True(t, res.Schema[name].ForceNew, "%s should force a new rule", name)
	}

	computedOnly := []string{
		"sequence", "source_name", "target_name", "scope_desc", "fw_ins_id",
		"country_name", "city_name", "country_code", "city_code", "country_key",
		"city_key", "is_region", "cloud_code", "is_cloud", "param_template_name",
		"invalid", "belong_member", "create_time", "update_time", "dns_parse_cnt",
		"hit_count", "last_hit_time",
	}
	for _, name := range computedOnly {
		s := res.Schema[name]
		if s == nil {
			t.Fatalf("schema %s is missing", name)
		}
		assert.True(t, s.Computed, "%s should be computed", name)
		assert.False(t, s.Optional, "%s should not be configurable", name)
	}

	assert.True(t, res.Schema["rule_id"].Computed)
	assert.True(t, res.Schema["rule_id"].Optional)
}
