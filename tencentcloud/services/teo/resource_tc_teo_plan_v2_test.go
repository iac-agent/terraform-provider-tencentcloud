package teo_test

import (
	"fmt"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/teo"
)

// mockMeta implements tccommon.ProviderMeta
type mockMetaPlanV2 struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaPlanV2) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaPlanV2{}

func newMockMetaPlanV2() *mockMetaPlanV2 {
	return &mockMetaPlanV2{client: &connectivity.TencentCloudClient{}}
}

// go test ./tencentcloud/services/teo/ -run "TestTeoPlanV2" -v -count=1 -gcflags="all=-l"

// TestTeoPlanV2_Create_Success tests Create calls CreatePlanWithContext and sets ID
func TestTeoPlanV2_Create_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaPlanV2().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "CreatePlanWithContext", func(ctx interface{}, request *teov20220901.CreatePlanRequest) (*teov20220901.CreatePlanResponse, error) {
		assert.Equal(t, "personal", *request.PlanType)
		assert.Equal(t, "true", *request.AutoUseVoucher)
		assert.NotNil(t, request.PrepaidPlanParam)
		assert.Equal(t, int64(1), *request.PrepaidPlanParam.Period)
		assert.Equal(t, "on", *request.PrepaidPlanParam.RenewFlag)

		resp := teov20220901.NewCreatePlanResponse()
		resp.Response = &teov20220901.CreatePlanResponseParams{
			PlanId:    ptrStringPlanV2("edgeone-2unuvzjmmn2q"),
			DealName:  ptrStringPlanV2("20250101123456789000"),
			RequestId: ptrStringPlanV2("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(teoClient, "DescribePlans", func(request *teov20220901.DescribePlansRequest) (*teov20220901.DescribePlansResponse, error) {
		payMode := int64(1)
		resp := teov20220901.NewDescribePlansResponse()
		resp.Response = &teov20220901.DescribePlansResponseParams{
			Plans: []*teov20220901.Plan{
				{
					PlanId:      ptrStringPlanV2("edgeone-2unuvzjmmn2q"),
					PlanType:    ptrStringPlanV2("personal"),
					Area:        ptrStringPlanV2("mainland"),
					Status:      ptrStringPlanV2("normal"),
					PayMode:     &payMode,
					EnabledTime: ptrStringPlanV2("2025-01-01T00:00:00+08:00"),
					ExpiredTime: ptrStringPlanV2("2025-02-01T00:00:00+08:00"),
				},
			},
			RequestId: ptrStringPlanV2("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaPlanV2()
	res := teo.ResourceTencentCloudTeoPlanV2()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"plan_type":        "personal",
		"auto_use_voucher": "true",
		"prepaid_plan_param": []interface{}{
			map[string]interface{}{
				"period":     1,
				"renew_flag": "on",
			},
		},
		"renew_flag": "on",
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "edgeone-2unuvzjmmn2q", d.Id())
	assert.Equal(t, "edgeone-2unuvzjmmn2q", d.Get("plan_id"))
	assert.Equal(t, "personal", d.Get("plan_type"))
}

// TestTeoPlanV2_Create_WithoutOptionalParams tests Create with only required params
func TestTeoPlanV2_Create_WithoutOptionalParams(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaPlanV2().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "CreatePlanWithContext", func(ctx interface{}, request *teov20220901.CreatePlanRequest) (*teov20220901.CreatePlanResponse, error) {
		assert.Equal(t, "enterprise", *request.PlanType)
		assert.Nil(t, request.AutoUseVoucher)
		assert.Nil(t, request.PrepaidPlanParam)

		resp := teov20220901.NewCreatePlanResponse()
		resp.Response = &teov20220901.CreatePlanResponseParams{
			PlanId:    ptrStringPlanV2("edgeone-enterprise-abc"),
			DealName:  ptrStringPlanV2("20250101123456789001"),
			RequestId: ptrStringPlanV2("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(teoClient, "DescribePlans", func(request *teov20220901.DescribePlansRequest) (*teov20220901.DescribePlansResponse, error) {
		payMode := int64(0)
		resp := teov20220901.NewDescribePlansResponse()
		resp.Response = &teov20220901.DescribePlansResponseParams{
			Plans: []*teov20220901.Plan{
				{
					PlanId:      ptrStringPlanV2("edgeone-enterprise-abc"),
					PlanType:    ptrStringPlanV2("enterprise"),
					Area:        ptrStringPlanV2("global"),
					Status:      ptrStringPlanV2("normal"),
					PayMode:     &payMode,
					EnabledTime: ptrStringPlanV2("2025-01-01T00:00:00+08:00"),
					ExpiredTime: ptrStringPlanV2(""),
				},
			},
			RequestId: ptrStringPlanV2("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaPlanV2()
	res := teo.ResourceTencentCloudTeoPlanV2()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"plan_type": "enterprise",
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "edgeone-enterprise-abc", d.Id())
}

// TestTeoPlanV2_Create_APIError tests Create handles API error
func TestTeoPlanV2_Create_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaPlanV2().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "CreatePlanWithContext", func(ctx interface{}, request *teov20220901.CreatePlanRequest) (*teov20220901.CreatePlanResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=InvalidParameter, Message=Invalid plan_type")
	})

	meta := newMockMetaPlanV2()
	res := teo.ResourceTencentCloudTeoPlanV2()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"plan_type": "personal",
	})

	err := res.Create(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "InvalidParameter")
}

// TestTeoPlanV2_Create_NilPlanId tests Create handles nil PlanId in response
func TestTeoPlanV2_Create_NilPlanId(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaPlanV2().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "CreatePlanWithContext", func(ctx interface{}, request *teov20220901.CreatePlanRequest) (*teov20220901.CreatePlanResponse, error) {
		resp := teov20220901.NewCreatePlanResponse()
		resp.Response = &teov20220901.CreatePlanResponseParams{
			PlanId:    nil,
			DealName:  ptrStringPlanV2("20250101123456789000"),
			RequestId: ptrStringPlanV2("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaPlanV2()
	res := teo.ResourceTencentCloudTeoPlanV2()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"plan_type": "personal",
	})

	err := res.Create(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "PlanId is nil")
}

// TestTeoPlanV2_Read_Success tests Read retrieves plan data
func TestTeoPlanV2_Read_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaPlanV2().client, "UseTeoV20220901Client", teoClient)

	payMode := int64(1)
	patches.ApplyMethodFunc(teoClient, "DescribePlans", func(request *teov20220901.DescribePlansRequest) (*teov20220901.DescribePlansResponse, error) {
		resp := teov20220901.NewDescribePlansResponse()
		resp.Response = &teov20220901.DescribePlansResponseParams{
			Plans: []*teov20220901.Plan{
				{
					PlanId:      ptrStringPlanV2("edgeone-2unuvzjmmn2q"),
					PlanType:    ptrStringPlanV2("personal"),
					Area:        ptrStringPlanV2("mainland"),
					Status:      ptrStringPlanV2("normal"),
					PayMode:     &payMode,
					EnabledTime: ptrStringPlanV2("2025-01-01T00:00:00+08:00"),
					ExpiredTime: ptrStringPlanV2("2025-02-01T00:00:00+08:00"),
				},
			},
			RequestId: ptrStringPlanV2("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaPlanV2()
	res := teo.ResourceTencentCloudTeoPlanV2()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"plan_type": "personal",
	})
	d.SetId("edgeone-2unuvzjmmn2q")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "edgeone-2unuvzjmmn2q", d.Get("plan_id"))
	assert.Equal(t, "personal", d.Get("plan_type"))
	assert.Equal(t, "mainland", d.Get("area"))
	assert.Equal(t, "normal", d.Get("status"))
	assert.Equal(t, "1", d.Get("pay_mode"))
	assert.Equal(t, "2025-01-01T00:00:00+08:00", d.Get("enabled_time"))
	assert.Equal(t, "2025-02-01T00:00:00+08:00", d.Get("expired_time"))
}

// TestTeoPlanV2_Read_NotFound tests Read handles plan not found
func TestTeoPlanV2_Read_NotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaPlanV2().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribePlans", func(request *teov20220901.DescribePlansRequest) (*teov20220901.DescribePlansResponse, error) {
		resp := teov20220901.NewDescribePlansResponse()
		resp.Response = &teov20220901.DescribePlansResponseParams{
			Plans:     []*teov20220901.Plan{},
			RequestId: ptrStringPlanV2("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaPlanV2()
	res := teo.ResourceTencentCloudTeoPlanV2()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"plan_type": "personal",
	})
	d.SetId("edgeone-not-exist")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}

// TestTeoPlanV2_Update_Success tests Update calls ModifyPlanWithContext
func TestTeoPlanV2_Update_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaPlanV2().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "ModifyPlanWithContext", func(ctx interface{}, request *teov20220901.ModifyPlanRequest) (*teov20220901.ModifyPlanResponse, error) {
		assert.Equal(t, "edgeone-2unuvzjmmn2q", *request.PlanId)
		assert.NotNil(t, request.RenewFlag)
		assert.Equal(t, "off", *request.RenewFlag.Switch)

		resp := teov20220901.NewModifyPlanResponse()
		resp.Response = &teov20220901.ModifyPlanResponseParams{
			RequestId: ptrStringPlanV2("fake-request-id"),
		}
		return resp, nil
	})

	payMode := int64(1)
	patches.ApplyMethodFunc(teoClient, "DescribePlans", func(request *teov20220901.DescribePlansRequest) (*teov20220901.DescribePlansResponse, error) {
		resp := teov20220901.NewDescribePlansResponse()
		resp.Response = &teov20220901.DescribePlansResponseParams{
			Plans: []*teov20220901.Plan{
				{
					PlanId:      ptrStringPlanV2("edgeone-2unuvzjmmn2q"),
					PlanType:    ptrStringPlanV2("personal"),
					Area:        ptrStringPlanV2("mainland"),
					Status:      ptrStringPlanV2("normal"),
					PayMode:     &payMode,
					EnabledTime: ptrStringPlanV2("2025-01-01T00:00:00+08:00"),
					ExpiredTime: ptrStringPlanV2("2025-02-01T00:00:00+08:00"),
				},
			},
			RequestId: ptrStringPlanV2("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaPlanV2()
	res := teo.ResourceTencentCloudTeoPlanV2()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"plan_type":        "personal",
		"auto_use_voucher": "true",
		"prepaid_plan_param": []interface{}{
			map[string]interface{}{
				"period":     1,
				"renew_flag": "on",
			},
		},
		"renew_flag": "off",
	})
	d.SetId("edgeone-2unuvzjmmn2q")

	err := res.Update(d, meta)
	assert.NoError(t, err)
}

// TestTeoPlanV2_Update_APIError tests Update handles API error
func TestTeoPlanV2_Update_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaPlanV2().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "ModifyPlanWithContext", func(ctx interface{}, request *teov20220901.ModifyPlanRequest) (*teov20220901.ModifyPlanResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=ResourceNotFound, Message=Plan not found")
	})

	meta := newMockMetaPlanV2()
	res := teo.ResourceTencentCloudTeoPlanV2()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"plan_type":  "personal",
		"renew_flag": "off",
	})
	d.SetId("edgeone-not-exist")

	err := res.Update(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ResourceNotFound")
}

// TestTeoPlanV2_Delete_Success tests Delete calls DestroyPlanWithContext
func TestTeoPlanV2_Delete_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaPlanV2().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DestroyPlanWithContext", func(ctx interface{}, request *teov20220901.DestroyPlanRequest) (*teov20220901.DestroyPlanResponse, error) {
		assert.Equal(t, "edgeone-2unuvzjmmn2q", *request.PlanId)

		resp := teov20220901.NewDestroyPlanResponse()
		resp.Response = &teov20220901.DestroyPlanResponseParams{
			RequestId: ptrStringPlanV2("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaPlanV2()
	res := teo.ResourceTencentCloudTeoPlanV2()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"plan_type": "personal",
	})
	d.SetId("edgeone-2unuvzjmmn2q")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
}

// TestTeoPlanV2_Delete_APIError tests Delete handles API error
func TestTeoPlanV2_Delete_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaPlanV2().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DestroyPlanWithContext", func(ctx interface{}, request *teov20220901.DestroyPlanRequest) (*teov20220901.DestroyPlanResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=ResourceNotFound, Message=Plan not found")
	})

	meta := newMockMetaPlanV2()
	res := teo.ResourceTencentCloudTeoPlanV2()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"plan_type": "personal",
	})
	d.SetId("edgeone-not-exist")

	err := res.Delete(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ResourceNotFound")
}

// TestTeoPlanV2_Schema validates schema definition
func TestTeoPlanV2_Schema(t *testing.T) {
	res := teo.ResourceTencentCloudTeoPlanV2()

	assert.NotNil(t, res)
	assert.NotNil(t, res.Create)
	assert.NotNil(t, res.Read)
	assert.NotNil(t, res.Update)
	assert.NotNil(t, res.Delete)
	assert.NotNil(t, res.Importer)

	// Check required fields with ForceNew
	assert.Contains(t, res.Schema, "plan_type")
	planType := res.Schema["plan_type"]
	assert.Equal(t, schema.TypeString, planType.Type)
	assert.True(t, planType.Required)
	assert.True(t, planType.ForceNew)

	// Check optional fields with ForceNew
	assert.Contains(t, res.Schema, "auto_use_voucher")
	autoUseVoucher := res.Schema["auto_use_voucher"]
	assert.Equal(t, schema.TypeString, autoUseVoucher.Type)
	assert.True(t, autoUseVoucher.Optional)
	assert.True(t, autoUseVoucher.ForceNew)

	assert.Contains(t, res.Schema, "prepaid_plan_param")
	prepaidPlanParam := res.Schema["prepaid_plan_param"]
	assert.Equal(t, schema.TypeList, prepaidPlanParam.Type)
	assert.True(t, prepaidPlanParam.Optional)
	assert.True(t, prepaidPlanParam.ForceNew)
	assert.Equal(t, 1, prepaidPlanParam.MaxItems)

	// Check optional field without ForceNew (top-level renew_flag for ModifyPlan)
	assert.Contains(t, res.Schema, "renew_flag")
	renewFlag := res.Schema["renew_flag"]
	assert.Equal(t, schema.TypeString, renewFlag.Type)
	assert.True(t, renewFlag.Optional)
	assert.False(t, renewFlag.ForceNew)

	// Check computed fields
	assert.Contains(t, res.Schema, "plan_id")
	planId := res.Schema["plan_id"]
	assert.Equal(t, schema.TypeString, planId.Type)
	assert.True(t, planId.Computed)

	assert.Contains(t, res.Schema, "deal_name")
	dealName := res.Schema["deal_name"]
	assert.Equal(t, schema.TypeString, dealName.Type)
	assert.True(t, dealName.Computed)

	assert.Contains(t, res.Schema, "area")
	area := res.Schema["area"]
	assert.Equal(t, schema.TypeString, area.Type)
	assert.True(t, area.Computed)

	assert.Contains(t, res.Schema, "status")
	status := res.Schema["status"]
	assert.Equal(t, schema.TypeString, status.Type)
	assert.True(t, status.Computed)

	assert.Contains(t, res.Schema, "pay_mode")
	payMode := res.Schema["pay_mode"]
	assert.Equal(t, schema.TypeString, payMode.Type)
	assert.True(t, payMode.Computed)

	assert.Contains(t, res.Schema, "enabled_time")
	enabledTime := res.Schema["enabled_time"]
	assert.Equal(t, schema.TypeString, enabledTime.Type)
	assert.True(t, enabledTime.Computed)

	assert.Contains(t, res.Schema, "expired_time")
	expiredTime := res.Schema["expired_time"]
	assert.Equal(t, schema.TypeString, expiredTime.Type)
	assert.True(t, expiredTime.Computed)
}

func ptrStringPlanV2(s string) *string {
	return &s
}
