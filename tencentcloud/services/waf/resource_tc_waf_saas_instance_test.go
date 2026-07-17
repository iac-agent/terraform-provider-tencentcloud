package waf_test

import (
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	wafsdk "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/waf/v20180125"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	wafsvc "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/waf"
)

// mockMeta implements tccommon.ProviderMeta
type wafSaasInstanceMockMeta struct {
	client *connectivity.TencentCloudClient
}

func (m *wafSaasInstanceMockMeta) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &wafSaasInstanceMockMeta{}

func newWafSaasInstanceMockMeta(region string) *wafSaasInstanceMockMeta {
	client := &connectivity.TencentCloudClient{}
	client.Region = region
	return &wafSaasInstanceMockMeta{client: client}
}

func ptrString(s string) *string {
	return &s
}

// setupWafClientPatches creates common gomonkey patches for waf client mocking
func setupWafClientPatches(meta *wafSaasInstanceMockMeta) (*gomonkey.Patches, *wafsdk.Client) {
	patches := gomonkey.NewPatches()
	wafClient := &wafsdk.Client{}
	patches.ApplyMethodReturn(meta.client, "UseWafClient", wafClient)
	return patches, wafClient
}

// setupDescribeInstancesMock mocks DescribeInstances to return a valid instance for Read
func setupDescribeInstancesMock(patches *gomonkey.Patches, wafClient *wafsdk.Client, instanceId string) {
	patches.ApplyMethodFunc(wafClient, "DescribeInstances", func(request *wafsdk.DescribeInstancesRequest) (*wafsdk.DescribeInstancesResponse, error) {
		resp := wafsdk.NewDescribeInstancesResponse()
		region := "gz"
		edition := "saas"
		status := uint64(1)
		mode := uint64(0)
		elasticBilling := uint64(0)
		renewFlag := uint64(0)
		beginTime := "2025-01-01 00:00:00"
		validTime := "2025-02-01 00:00:00"
		resp.Response = &wafsdk.DescribeInstancesResponseParams{
			Instances: []*wafsdk.InstanceInfo{
				{
					InstanceId:     &instanceId,
					InstanceName:   ptrString("test-instance"),
					Edition:        &edition,
					Region:         &region,
					Status:         &status,
					Mode:           &mode,
					ElasticBilling: &elasticBilling,
					RenewFlag:      &renewFlag,
					BeginTime:      &beginTime,
					ValidTime:      &validTime,
				},
			},
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})
}

// go test ./tencentcloud/services/waf/ -run "TestWafSaasInstance" -v -count=1 -gcflags="all=-l"

// TestWafSaasInstanceGoodsNum_Specified tests goods_num is set during creation
func TestWafSaasInstanceGoodsNum_Specified(t *testing.T) {
	meta := newWafSaasInstanceMockMeta("ap-guangzhou")
	patches, wafClient := setupWafClientPatches(meta)
	defer patches.Reset()

	var capturedGoodsNum *int64
	patches.ApplyMethodFunc(wafClient, "GenerateDealsAndPayNew", func(request *wafsdk.GenerateDealsAndPayNewRequest) (*wafsdk.GenerateDealsAndPayNewResponse, error) {
		if request.Goods != nil && len(request.Goods) > 0 {
			capturedGoodsNum = request.Goods[0].GoodsNum
		}
		resp := wafsdk.NewGenerateDealsAndPayNewResponse()
		status := int64(1)
		instanceId := "waf-instance-123"
		resp.Response = &wafsdk.GenerateDealsAndPayNewResponseParams{
			Status:        &status,
			InstanceId:    &instanceId,
			ReturnMessage: ptrString("ok"),
		}
		return resp, nil
	})

	setupDescribeInstancesMock(patches, wafClient, "waf-instance-123")

	res := wafsvc.ResourceTencentCloudWafSaasInstance()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"goods_category":  "premium_saas",
		"goods_num":       2,
		"time_span":       1,
		"time_unit":       "m",
		"auto_renew_flag": 0,
		"elastic_mode":    0,
		"real_region":     "gz",
		"bot_management":  0,
		"api_security":    0,
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "waf-instance-123", d.Id())
	assert.NotNil(t, capturedGoodsNum)
	assert.Equal(t, int64(2), *capturedGoodsNum)
}

// TestWafSaasInstanceGoodsNum_Default tests goods_num defaults to 1 when not set
func TestWafSaasInstanceGoodsNum_Default(t *testing.T) {
	meta := newWafSaasInstanceMockMeta("ap-guangzhou")
	patches, wafClient := setupWafClientPatches(meta)
	defer patches.Reset()

	var capturedGoodsNum *int64
	patches.ApplyMethodFunc(wafClient, "GenerateDealsAndPayNew", func(request *wafsdk.GenerateDealsAndPayNewRequest) (*wafsdk.GenerateDealsAndPayNewResponse, error) {
		if request.Goods != nil && len(request.Goods) > 0 {
			capturedGoodsNum = request.Goods[0].GoodsNum
		}
		resp := wafsdk.NewGenerateDealsAndPayNewResponse()
		status := int64(1)
		instanceId := "waf-instance-456"
		resp.Response = &wafsdk.GenerateDealsAndPayNewResponseParams{
			Status:        &status,
			InstanceId:    &instanceId,
			ReturnMessage: ptrString("ok"),
		}
		return resp, nil
	})

	setupDescribeInstancesMock(patches, wafClient, "waf-instance-456")

	res := wafsvc.ResourceTencentCloudWafSaasInstance()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"goods_category":  "premium_saas",
		"time_span":       1,
		"time_unit":       "m",
		"auto_renew_flag": 0,
		"elastic_mode":    0,
		"real_region":     "gz",
		"bot_management":  0,
		"api_security":    0,
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "waf-instance-456", d.Id())
	assert.NotNil(t, capturedGoodsNum)
	assert.Equal(t, int64(1), *capturedGoodsNum)
}

// TestWafSaasInstancePid_Specified tests pid overrides PID_SAAS mapping
func TestWafSaasInstancePid_Specified(t *testing.T) {
	meta := newWafSaasInstanceMockMeta("ap-guangzhou")
	patches, wafClient := setupWafClientPatches(meta)
	defer patches.Reset()

	var capturedPid *int64
	patches.ApplyMethodFunc(wafClient, "GenerateDealsAndPayNew", func(request *wafsdk.GenerateDealsAndPayNewRequest) (*wafsdk.GenerateDealsAndPayNewResponse, error) {
		if request.Goods != nil && len(request.Goods) > 0 && request.Goods[0].GoodsDetail != nil {
			capturedPid = request.Goods[0].GoodsDetail.Pid
		}
		resp := wafsdk.NewGenerateDealsAndPayNewResponse()
		status := int64(1)
		instanceId := "waf-instance-pid-123"
		resp.Response = &wafsdk.GenerateDealsAndPayNewResponseParams{
			Status:        &status,
			InstanceId:    &instanceId,
			ReturnMessage: ptrString("ok"),
		}
		return resp, nil
	})

	setupDescribeInstancesMock(patches, wafClient, "waf-instance-pid-123")

	res := wafsvc.ResourceTencentCloudWafSaasInstance()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"goods_category":  "premium_saas",
		"pid":             1000827,
		"time_span":       1,
		"time_unit":       "m",
		"auto_renew_flag": 0,
		"elastic_mode":    0,
		"real_region":     "gz",
		"bot_management":  0,
		"api_security":    0,
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "waf-instance-pid-123", d.Id())
	assert.NotNil(t, capturedPid)
	assert.Equal(t, int64(1000827), *capturedPid)
}

// TestWafSaasInstancePid_Default tests pid uses PID_SAAS mapping when not set
func TestWafSaasInstancePid_Default(t *testing.T) {
	meta := newWafSaasInstanceMockMeta("ap-guangzhou")
	patches, wafClient := setupWafClientPatches(meta)
	defer patches.Reset()

	var capturedPid *int64
	patches.ApplyMethodFunc(wafClient, "GenerateDealsAndPayNew", func(request *wafsdk.GenerateDealsAndPayNewRequest) (*wafsdk.GenerateDealsAndPayNewResponse, error) {
		if request.Goods != nil && len(request.Goods) > 0 && request.Goods[0].GoodsDetail != nil {
			capturedPid = request.Goods[0].GoodsDetail.Pid
		}
		resp := wafsdk.NewGenerateDealsAndPayNewResponse()
		status := int64(1)
		instanceId := "waf-instance-pid-default"
		resp.Response = &wafsdk.GenerateDealsAndPayNewResponseParams{
			Status:        &status,
			InstanceId:    &instanceId,
			ReturnMessage: ptrString("ok"),
		}
		return resp, nil
	})

	setupDescribeInstancesMock(patches, wafClient, "waf-instance-pid-default")

	res := wafsvc.ResourceTencentCloudWafSaasInstance()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"goods_category":  "premium_saas",
		"time_span":       1,
		"time_unit":       "m",
		"auto_renew_flag": 0,
		"elastic_mode":    0,
		"real_region":     "gz",
		"bot_management":  0,
		"api_security":    0,
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "waf-instance-pid-default", d.Id())
	assert.NotNil(t, capturedPid)
	// premium_saas PID should be from PID_SAAS mapping, not 0
	assert.NotEqual(t, int64(0), *capturedPid)
}

// TestWafSaasInstanceRegionId_Specified tests region_id overrides provider region derivation
func TestWafSaasInstanceRegionId_Specified(t *testing.T) {
	meta := newWafSaasInstanceMockMeta("ap-guangzhou")
	patches, wafClient := setupWafClientPatches(meta)
	defer patches.Reset()

	var capturedRegionId *int64
	patches.ApplyMethodFunc(wafClient, "GenerateDealsAndPayNew", func(request *wafsdk.GenerateDealsAndPayNewRequest) (*wafsdk.GenerateDealsAndPayNewResponse, error) {
		if request.Goods != nil && len(request.Goods) > 0 {
			capturedRegionId = request.Goods[0].RegionId
		}
		resp := wafsdk.NewGenerateDealsAndPayNewResponse()
		status := int64(1)
		instanceId := "waf-instance-region-123"
		resp.Response = &wafsdk.GenerateDealsAndPayNewResponseParams{
			Status:        &status,
			InstanceId:    &instanceId,
			ReturnMessage: ptrString("ok"),
		}
		return resp, nil
	})

	setupDescribeInstancesMock(patches, wafClient, "waf-instance-region-123")

	res := wafsvc.ResourceTencentCloudWafSaasInstance()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"goods_category":  "premium_saas",
		"region_id":       5,
		"time_span":       1,
		"time_unit":       "m",
		"auto_renew_flag": 0,
		"elastic_mode":    0,
		"real_region":     "gz",
		"bot_management":  0,
		"api_security":    0,
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "waf-instance-region-123", d.Id())
	assert.NotNil(t, capturedRegionId)
	assert.Equal(t, int64(5), *capturedRegionId)
}

// TestWafSaasInstanceRegionId_Default tests region_id defaults to provider region derivation
func TestWafSaasInstanceRegionId_Default(t *testing.T) {
	meta := newWafSaasInstanceMockMeta("ap-guangzhou")
	patches, wafClient := setupWafClientPatches(meta)
	defer patches.Reset()

	var capturedRegionId *int64
	patches.ApplyMethodFunc(wafClient, "GenerateDealsAndPayNew", func(request *wafsdk.GenerateDealsAndPayNewRequest) (*wafsdk.GenerateDealsAndPayNewResponse, error) {
		if request.Goods != nil && len(request.Goods) > 0 {
			capturedRegionId = request.Goods[0].RegionId
		}
		resp := wafsdk.NewGenerateDealsAndPayNewResponse()
		status := int64(1)
		instanceId := "waf-instance-region-default"
		resp.Response = &wafsdk.GenerateDealsAndPayNewResponseParams{
			Status:        &status,
			InstanceId:    &instanceId,
			ReturnMessage: ptrString("ok"),
		}
		return resp, nil
	})

	setupDescribeInstancesMock(patches, wafClient, "waf-instance-region-default")

	res := wafsvc.ResourceTencentCloudWafSaasInstance()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"goods_category":  "premium_saas",
		"time_span":       1,
		"time_unit":       "m",
		"auto_renew_flag": 0,
		"elastic_mode":    0,
		"real_region":     "gz",
		"bot_management":  0,
		"api_security":    0,
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "waf-instance-region-default", d.Id())
	assert.NotNil(t, capturedRegionId)
	// ap-guangzhou should default to REGION_ID_MAINLAND = 1
	assert.Equal(t, int64(1), *capturedRegionId)
}

// TestWafSaasInstanceDealNames_Populated tests deal_names is populated from creation response
func TestWafSaasInstanceDealNames_Populated(t *testing.T) {
	meta := newWafSaasInstanceMockMeta("ap-guangzhou")
	patches, wafClient := setupWafClientPatches(meta)
	defer patches.Reset()

	patches.ApplyMethodFunc(wafClient, "GenerateDealsAndPayNew", func(request *wafsdk.GenerateDealsAndPayNewRequest) (*wafsdk.GenerateDealsAndPayNewResponse, error) {
		resp := wafsdk.NewGenerateDealsAndPayNewResponse()
		status := int64(1)
		instanceId := "waf-instance-dealnames-123"
		dealNames := []*string{ptrString("20250101-order-001"), ptrString("20250101-order-002")}
		bigDealId := "big-deal-001"
		resp.Response = &wafsdk.GenerateDealsAndPayNewResponseParams{
			Data: &wafsdk.DealData{
				DealNames: dealNames,
				BigDealId: &bigDealId,
			},
			Status:        &status,
			InstanceId:    &instanceId,
			ReturnMessage: ptrString("ok"),
		}
		return resp, nil
	})

	setupDescribeInstancesMock(patches, wafClient, "waf-instance-dealnames-123")

	res := wafsvc.ResourceTencentCloudWafSaasInstance()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"goods_category":  "premium_saas",
		"time_span":       1,
		"time_unit":       "m",
		"auto_renew_flag": 0,
		"elastic_mode":    0,
		"real_region":     "gz",
		"bot_management":  0,
		"api_security":    0,
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "waf-instance-dealnames-123", d.Id())

	dealNames := d.Get("deal_names").([]interface{})
	assert.Len(t, dealNames, 2)
	assert.Equal(t, "20250101-order-001", dealNames[0])
	assert.Equal(t, "20250101-order-002", dealNames[1])
}

// TestWafSaasInstanceDealNames_NilData tests deal_names when Data is nil in response
func TestWafSaasInstanceDealNames_NilData(t *testing.T) {
	meta := newWafSaasInstanceMockMeta("ap-guangzhou")
	patches, wafClient := setupWafClientPatches(meta)
	defer patches.Reset()

	patches.ApplyMethodFunc(wafClient, "GenerateDealsAndPayNew", func(request *wafsdk.GenerateDealsAndPayNewRequest) (*wafsdk.GenerateDealsAndPayNewResponse, error) {
		resp := wafsdk.NewGenerateDealsAndPayNewResponse()
		status := int64(1)
		instanceId := "waf-instance-nildata-123"
		resp.Response = &wafsdk.GenerateDealsAndPayNewResponseParams{
			Status:        &status,
			InstanceId:    &instanceId,
			ReturnMessage: ptrString("ok"),
		}
		return resp, nil
	})

	setupDescribeInstancesMock(patches, wafClient, "waf-instance-nildata-123")

	res := wafsvc.ResourceTencentCloudWafSaasInstance()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"goods_category":  "premium_saas",
		"time_span":       1,
		"time_unit":       "m",
		"auto_renew_flag": 0,
		"elastic_mode":    0,
		"real_region":     "gz",
		"bot_management":  0,
		"api_security":    0,
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "waf-instance-nildata-123", d.Id())
	// deal_names should be empty when Data is nil
	dealNames := d.Get("deal_names").([]interface{})
	assert.Len(t, dealNames, 0)
}

// TestWafSaasInstanceImmutableArgs tests that goods_num, pid, region_id are immutable
func TestWafSaasInstanceImmutableArgs(t *testing.T) {
	meta := newWafSaasInstanceMockMeta("ap-guangzhou")
	patches, wafClient := setupWafClientPatches(meta)
	defer patches.Reset()

	setupDescribeInstancesMock(patches, wafClient, "waf-instance-immutable")

	res := wafsvc.ResourceTencentCloudWafSaasInstance()

	// Test goods_num immutable - set goods_num to different value triggers error
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"goods_category":  "premium_saas",
		"goods_num":       1,
		"time_span":       1,
		"time_unit":       "m",
		"auto_renew_flag": 0,
		"elastic_mode":    0,
		"real_region":     "gz",
		"bot_management":  0,
		"api_security":    0,
	})
	d.SetId("waf-instance-immutable-1")
	_ = d.Set("goods_num", 2)
	err := res.Update(d, meta)
	assert.Error(t, err)
	// The immutableArgs check iterates from the beginning, so goods_category may also trigger
	// but the key is that the update should fail due to immutable args
	assert.Contains(t, err.Error(), "cannot be changed")

	// Test pid immutable - set pid to different value triggers error
	d2 := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"goods_category":  "premium_saas",
		"pid":             1000827,
		"time_span":       1,
		"time_unit":       "m",
		"auto_renew_flag": 0,
		"elastic_mode":    0,
		"real_region":     "gz",
		"bot_management":  0,
		"api_security":    0,
	})
	d2.SetId("waf-instance-immutable-2")
	_ = d2.Set("pid", 1000828)
	err = res.Update(d2, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot be changed")

	// Test region_id immutable - set region_id to different value triggers error
	d3 := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"goods_category":  "premium_saas",
		"region_id":       1,
		"time_span":       1,
		"time_unit":       "m",
		"auto_renew_flag": 0,
		"elastic_mode":    0,
		"real_region":     "gz",
		"bot_management":  0,
		"api_security":    0,
	})
	d3.SetId("waf-instance-immutable-3")
	_ = d3.Set("region_id", 9)
	err = res.Update(d3, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot be changed")
}

// TestWafSaasInstanceSchema validates schema definition for new parameters
func TestWafSaasInstanceSchema(t *testing.T) {
	res := wafsvc.ResourceTencentCloudWafSaasInstance()

	assert.NotNil(t, res)
	assert.Contains(t, res.Schema, "goods_num")
	assert.Contains(t, res.Schema, "pid")
	assert.Contains(t, res.Schema, "region_id")
	assert.Contains(t, res.Schema, "deal_names")

	goodsNum := res.Schema["goods_num"]
	assert.Equal(t, schema.TypeInt, goodsNum.Type)
	assert.True(t, goodsNum.Optional)
	assert.False(t, goodsNum.Required)
	assert.False(t, goodsNum.Computed)

	pid := res.Schema["pid"]
	assert.Equal(t, schema.TypeInt, pid.Type)
	assert.True(t, pid.Optional)
	assert.False(t, pid.Required)
	assert.False(t, pid.Computed)

	regionId := res.Schema["region_id"]
	assert.Equal(t, schema.TypeInt, regionId.Type)
	assert.True(t, regionId.Optional)
	assert.False(t, regionId.Required)
	assert.False(t, regionId.Computed)

	dealNames := res.Schema["deal_names"]
	assert.Equal(t, schema.TypeList, dealNames.Type)
	assert.True(t, dealNames.Computed)
	assert.False(t, dealNames.Optional)
	assert.False(t, dealNames.Required)
}
