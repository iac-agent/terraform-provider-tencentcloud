package vpc_test

import (
	"context"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	vpcv20170312 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	vpc "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/vpc"
)

type mockMetaVpc struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaVpc) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaVpc{}

func newMockMetaVpc() *mockMetaVpc {
	return &mockMetaVpc{client: &connectivity.TencentCloudClient{}}
}

func ptrStringVpc(s string) *string {
	return &s
}

func ptrUint64Vpc(i uint64) *uint64 {
	return &i
}

func TestAccVpcTrafficMirrorFilterRulesV2_Create(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	vpcClient := &vpcv20170312.Client{}
	patches.ApplyMethodReturn(newMockMetaVpc().client, "UseVpcClient", vpcClient)

	trafficMirrorId := "imgf-test123"
	ruleId1 := "tmfi-test-rule1"
	ruleId2 := "tmfi-test-rule2"

	patches.ApplyMethodFunc(vpcClient, "CreateTrafficMirrorFilterRulesWithContext", func(_ context.Context, request *vpcv20170312.CreateTrafficMirrorFilterRulesRequest) (*vpcv20170312.CreateTrafficMirrorFilterRulesResponse, error) {
		assert.Equal(t, trafficMirrorId, *request.TrafficMirrorId)
		assert.NotNil(t, request.IngressFilterRules)
		assert.Equal(t, 1, len(request.IngressFilterRules))
		assert.NotNil(t, request.EgressFilterRules)
		assert.Equal(t, 1, len(request.EgressFilterRules))

		resp := vpcv20170312.NewCreateTrafficMirrorFilterRulesResponse()
		resp.Response = &vpcv20170312.CreateTrafficMirrorFilterRulesResponseParams{
			TrafficMirrorId: &trafficMirrorId,
			IngressFilterRules: []*vpcv20170312.TrafficMirrorFilter{
				{
					SrcNet:                    ptrStringVpc("10.0.0.0/24"),
					DstNet:                    ptrStringVpc("10.0.1.0/24"),
					Protocol:                  ptrStringVpc("TCP"),
					SrcPort:                   ptrStringVpc("80"),
					DstPort:                   ptrStringVpc("8080"),
					TrafficMirrorFilterRuleId: &ruleId1,
					Priority:                  ptrUint64Vpc(1),
					Action:                    ptrStringVpc("ACCEPT"),
					Description:               ptrStringVpc("ingress rule"),
					CreatedTime:               ptrStringVpc("2024-01-01 00:00:00"),
				},
			},
			EgressFilterRules: []*vpcv20170312.TrafficMirrorFilter{
				{
					SrcNet:                    ptrStringVpc("10.0.1.0/24"),
					DstNet:                    ptrStringVpc("10.0.0.0/24"),
					Protocol:                  ptrStringVpc("TCP"),
					SrcPort:                   ptrStringVpc("8080"),
					DstPort:                   ptrStringVpc("80"),
					TrafficMirrorFilterRuleId: &ruleId2,
					Priority:                  ptrUint64Vpc(1),
					Action:                    ptrStringVpc("ACCEPT"),
					Description:               ptrStringVpc("egress rule"),
					CreatedTime:               ptrStringVpc("2024-01-01 00:00:00"),
				},
			},
			RequestId: ptrStringVpc("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaVpc()
	res := vpc.ResourceTencentCloudVpcTrafficMirrorFilterRulesV2()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"traffic_mirror_id": trafficMirrorId,
		"ingress_filter_rules": []interface{}{
			map[string]interface{}{
				"src_net":     "10.0.0.0/24",
				"dst_net":     "10.0.1.0/24",
				"protocol":    "TCP",
				"src_port":    "80",
				"dst_port":    "8080",
				"priority":    1,
				"action":      "ACCEPT",
				"description": "ingress rule",
			},
		},
		"egress_filter_rules": []interface{}{
			map[string]interface{}{
				"src_net":     "10.0.1.0/24",
				"dst_net":     "10.0.0.0/24",
				"protocol":    "TCP",
				"src_port":    "8080",
				"dst_port":    "80",
				"priority":    1,
				"action":      "ACCEPT",
				"description": "egress rule",
			},
		},
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, trafficMirrorId, d.Id())
}

func TestAccVpcTrafficMirrorFilterRulesV2_Read(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	vpcClient := &vpcv20170312.Client{}
	patches.ApplyMethodReturn(newMockMetaVpc().client, "UseVpcClient", vpcClient)

	trafficMirrorId := "imgf-test123"
	ruleId1 := "tmfi-test-rule1"

	patches.ApplyMethodFunc(vpcClient, "DescribeTrafficMirrorFilterRulesWithContext", func(_ context.Context, request *vpcv20170312.DescribeTrafficMirrorFilterRulesRequest) (*vpcv20170312.DescribeTrafficMirrorFilterRulesResponse, error) {
		assert.Equal(t, trafficMirrorId, *request.TrafficMirrorId)

		resp := vpcv20170312.NewDescribeTrafficMirrorFilterRulesResponse()
		resp.Response = &vpcv20170312.DescribeTrafficMirrorFilterRulesResponseParams{
			TrafficMirrorId: &trafficMirrorId,
			IngressFilterRules: []*vpcv20170312.TrafficMirrorFilter{
				{
					SrcNet:                    ptrStringVpc("10.0.0.0/24"),
					DstNet:                    ptrStringVpc("10.0.1.0/24"),
					Protocol:                  ptrStringVpc("TCP"),
					SrcPort:                   ptrStringVpc("80"),
					DstPort:                   ptrStringVpc("8080"),
					TrafficMirrorFilterRuleId: &ruleId1,
					Priority:                  ptrUint64Vpc(1),
					Action:                    ptrStringVpc("ACCEPT"),
					Description:               ptrStringVpc("ingress rule"),
					CreatedTime:               ptrStringVpc("2024-01-01 00:00:00"),
				},
			},
			TotalCount: ptrUint64Vpc(1),
			RequestId:  ptrStringVpc("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaVpc()
	res := vpc.ResourceTencentCloudVpcTrafficMirrorFilterRulesV2()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"traffic_mirror_id": trafficMirrorId,
	})
	d.SetId(trafficMirrorId)

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, trafficMirrorId, d.Id())

	ingressRules := d.Get("ingress_filter_rules").([]interface{})
	assert.Equal(t, 1, len(ingressRules))
}

func TestAccVpcTrafficMirrorFilterRulesV2_ReadEmpty(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	vpcClient := &vpcv20170312.Client{}
	patches.ApplyMethodReturn(newMockMetaVpc().client, "UseVpcClient", vpcClient)

	trafficMirrorId := "imgf-test123"

	patches.ApplyMethodFunc(vpcClient, "DescribeTrafficMirrorFilterRulesWithContext", func(_ context.Context, request *vpcv20170312.DescribeTrafficMirrorFilterRulesRequest) (*vpcv20170312.DescribeTrafficMirrorFilterRulesResponse, error) {
		resp := vpcv20170312.NewDescribeTrafficMirrorFilterRulesResponse()
		resp.Response = &vpcv20170312.DescribeTrafficMirrorFilterRulesResponseParams{
			TrafficMirrorId: &trafficMirrorId,
			TotalCount:      ptrUint64Vpc(0),
			RequestId:       ptrStringVpc("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaVpc()
	res := vpc.ResourceTencentCloudVpcTrafficMirrorFilterRulesV2()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"traffic_mirror_id": trafficMirrorId,
	})
	d.SetId(trafficMirrorId)

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}

func TestAccVpcTrafficMirrorFilterRulesV2_Update(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	vpcClient := &vpcv20170312.Client{}
	patches.ApplyMethodReturn(newMockMetaVpc().client, "UseVpcClient", vpcClient)

	trafficMirrorId := "imgf-test123"
	ruleId1 := "tmfi-test-rule1"

	modifyCalled := false
	patches.ApplyMethodFunc(vpcClient, "ModifyTrafficMirrorFilterRulesWithContext", func(_ context.Context, request *vpcv20170312.ModifyTrafficMirrorFilterRulesRequest) (*vpcv20170312.ModifyTrafficMirrorFilterRulesResponse, error) {
		modifyCalled = true
		assert.Equal(t, trafficMirrorId, *request.TrafficMirrorId)
		assert.Equal(t, 1, len(request.IngressFilterRules))

		resp := vpcv20170312.NewModifyTrafficMirrorFilterRulesResponse()
		resp.Response = &vpcv20170312.ModifyTrafficMirrorFilterRulesResponseParams{
			RequestId: ptrStringVpc("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(vpcClient, "DescribeTrafficMirrorFilterRulesWithContext", func(_ context.Context, request *vpcv20170312.DescribeTrafficMirrorFilterRulesRequest) (*vpcv20170312.DescribeTrafficMirrorFilterRulesResponse, error) {
		resp := vpcv20170312.NewDescribeTrafficMirrorFilterRulesResponse()
		resp.Response = &vpcv20170312.DescribeTrafficMirrorFilterRulesResponseParams{
			TrafficMirrorId: &trafficMirrorId,
			IngressFilterRules: []*vpcv20170312.TrafficMirrorFilter{
				{
					SrcNet:                    ptrStringVpc("10.0.0.0/24"),
					DstNet:                    ptrStringVpc("10.0.1.0/24"),
					Protocol:                  ptrStringVpc("TCP"),
					SrcPort:                   ptrStringVpc("80"),
					DstPort:                   ptrStringVpc("8080"),
					TrafficMirrorFilterRuleId: &ruleId1,
					Priority:                  ptrUint64Vpc(1),
					Action:                    ptrStringVpc("ACCEPT"),
					Description:               ptrStringVpc("ingress rule"),
					CreatedTime:               ptrStringVpc("2024-01-01 00:00:00"),
				},
			},
			TotalCount: ptrUint64Vpc(1),
			RequestId:  ptrStringVpc("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaVpc()
	res := vpc.ResourceTencentCloudVpcTrafficMirrorFilterRulesV2()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"traffic_mirror_id": trafficMirrorId,
		"ingress_filter_rules": []interface{}{
			map[string]interface{}{
				"src_net":     "10.0.0.0/24",
				"dst_net":     "10.0.1.0/24",
				"protocol":    "TCP",
				"src_port":    "80",
				"dst_port":    "8080",
				"priority":    1,
				"action":      "ACCEPT",
				"description": "ingress rule",
			},
		},
	})
	d.SetId(trafficMirrorId)

	// Mark HasChange on ingress_filter_rules to trigger update
	patches.ApplyMethodFunc(d, "HasChange", func(key string) bool {
		return key == "ingress_filter_rules"
	})

	err := res.Update(d, meta)
	assert.NoError(t, err)
	assert.True(t, modifyCalled)
}

func TestAccVpcTrafficMirrorFilterRulesV2_Delete(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	vpcClient := &vpcv20170312.Client{}
	patches.ApplyMethodReturn(newMockMetaVpc().client, "UseVpcClient", vpcClient)

	trafficMirrorId := "imgf-test123"
	ruleId1 := "tmfi-test-rule1"

	deleteCalled := false
	patches.ApplyMethodFunc(vpcClient, "DeleteTrafficMirrorFilterRulesWithContext", func(_ context.Context, request *vpcv20170312.DeleteTrafficMirrorFilterRulesRequest) (*vpcv20170312.DeleteTrafficMirrorFilterRulesResponse, error) {
		deleteCalled = true
		assert.Equal(t, trafficMirrorId, *request.TrafficMirrorId)
		assert.Equal(t, 1, len(request.IngressFilterRuleIds))
		assert.Equal(t, ruleId1, *request.IngressFilterRuleIds[0])

		resp := vpcv20170312.NewDeleteTrafficMirrorFilterRulesResponse()
		resp.Response = &vpcv20170312.DeleteTrafficMirrorFilterRulesResponseParams{
			RequestId: ptrStringVpc("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaVpc()
	res := vpc.ResourceTencentCloudVpcTrafficMirrorFilterRulesV2()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"traffic_mirror_id": trafficMirrorId,
		"ingress_filter_rules": []interface{}{
			map[string]interface{}{
				"traffic_mirror_filter_rule_id": ruleId1,
				"src_net":                       "10.0.0.0/24",
				"dst_net":                       "10.0.1.0/24",
				"protocol":                      "TCP",
				"src_port":                      "80",
				"dst_port":                      "8080",
				"priority":                      1,
				"action":                        "ACCEPT",
				"description":                   "ingress rule",
			},
		},
	})
	d.SetId(trafficMirrorId)

	err := res.Delete(d, meta)
	assert.NoError(t, err)
	assert.True(t, deleteCalled)
}

func TestAccVpcTrafficMirrorFilterRulesV2_CreateNoRules(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	vpcClient := &vpcv20170312.Client{}
	patches.ApplyMethodReturn(newMockMetaVpc().client, "UseVpcClient", vpcClient)

	meta := newMockMetaVpc()
	res := vpc.ResourceTencentCloudVpcTrafficMirrorFilterRulesV2()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"traffic_mirror_id": "imgf-test123",
	})

	err := res.Create(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "At least one of")
}
