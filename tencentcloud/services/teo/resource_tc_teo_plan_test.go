package teo_test

import (
	"context"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/teo"
)

// mockMetaTeoPlan implements tccommon.ProviderMeta
type mockMetaTeoPlan struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaTeoPlan) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaTeoPlan{}

func newMockMetaTeoPlan() *mockMetaTeoPlan {
	return &mockMetaTeoPlan{client: &connectivity.TencentCloudClient{}}
}

func ptrStringTeoPlan(s string) *string {
	return &s
}

func ptrInt64TeoPlan(i int64) *int64 {
	return &i
}

// go test ./tencentcloud/services/teo/ -run "TestTeoPlan" -v -count=1 -gcflags="all=-l"

// mockDescribePlansTeoPlan mocks the DescribePlans API used by ResourceTencentCloudTeoPlanRead
// through the TeoService DescribeTeoPlansById helper.
func mockDescribePlansTeoPlan(patches *gomonkey.Patches, teoClient *teov20220901.Client, planId string) {
	patches.ApplyMethodFunc(teoClient, "DescribePlans", func(request *teov20220901.DescribePlansRequest) (*teov20220901.DescribePlansResponse, error) {
		resp := teov20220901.NewDescribePlansResponse()
		resp.Response = &teov20220901.DescribePlansResponseParams{
			TotalCount: ptrInt64TeoPlan(1),
			Plans: []*teov20220901.Plan{
				{
					PlanType:    ptrStringTeoPlan("plan-personal"),
					PlanId:      ptrStringTeoPlan(planId),
					Area:        ptrStringTeoPlan("global"),
					Status:      ptrStringTeoPlan("normal"),
					PayMode:     ptrInt64TeoPlan(1),
					EnabledTime: ptrStringTeoPlan("2026-01-01 00:00:00"),
					ExpiredTime: ptrStringTeoPlan("2026-02-01 00:00:00"),
				},
			},
			RequestId: ptrStringTeoPlan("fake-request-id"),
		}
		return resp, nil
	})
}

// TestTeoPlan_CreateWithAutoUseVoucher verifies auto_use_voucher is passed to the CreatePlan request.
func TestTeoPlan_CreateWithAutoUseVoucher(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	meta := newMockMetaTeoPlan()
	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(meta.client, "UseTeoV20220901Client", teoClient)

	var capturedRequest *teov20220901.CreatePlanRequest
	patches.ApplyMethodFunc(teoClient, "CreatePlanWithContext", func(_ context.Context, request *teov20220901.CreatePlanRequest) (*teov20220901.CreatePlanResponse, error) {
		capturedRequest = request
		resp := teov20220901.NewCreatePlanResponse()
		resp.Response = &teov20220901.CreatePlanResponseParams{
			PlanId:    ptrStringTeoPlan("edgeone-2unuvzjmmn2q"),
			RequestId: ptrStringTeoPlan("fake-request-id"),
		}
		return resp, nil
	})

	mockDescribePlansTeoPlan(patches, teoClient, "edgeone-2unuvzjmmn2q")

	res := teo.ResourceTencentCloudTeoPlan()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"plan_type":          "personal",
		"auto_use_voucher":   "true",
		"prepaid_plan_param": []interface{}{map[string]interface{}{"period": 1, "renew_flag": "off"}},
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "edgeone-2unuvzjmmn2q", d.Id())

	// Verify AutoUseVoucher was passed to the CreatePlan request.
	assert.NotNil(t, capturedRequest)
	assert.NotNil(t, capturedRequest.AutoUseVoucher)
	assert.Equal(t, "true", *capturedRequest.AutoUseVoucher)
	assert.Equal(t, "personal", *capturedRequest.PlanType)
	assert.NotNil(t, capturedRequest.PrepaidPlanParam)
	assert.Equal(t, int64(1), *capturedRequest.PrepaidPlanParam.Period)
	assert.Equal(t, "off", *capturedRequest.PrepaidPlanParam.RenewFlag)

	// Read after create refreshes computed fields.
	assert.Equal(t, "global", d.Get("area").(string))
	assert.Equal(t, "normal", d.Get("status").(string))
	assert.Equal(t, "edgeone-2unuvzjmmn2q", d.Get("plan_id").(string))
}

// TestTeoPlan_CreateWithoutAutoUseVoucher verifies AutoUseVoucher is omitted when not configured.
func TestTeoPlan_CreateWithoutAutoUseVoucher(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	meta := newMockMetaTeoPlan()
	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(meta.client, "UseTeoV20220901Client", teoClient)

	var capturedRequest *teov20220901.CreatePlanRequest
	patches.ApplyMethodFunc(teoClient, "CreatePlanWithContext", func(_ context.Context, request *teov20220901.CreatePlanRequest) (*teov20220901.CreatePlanResponse, error) {
		capturedRequest = request
		resp := teov20220901.NewCreatePlanResponse()
		resp.Response = &teov20220901.CreatePlanResponseParams{
			PlanId:    ptrStringTeoPlan("edgeone-2unuvzjmmn2q"),
			RequestId: ptrStringTeoPlan("fake-request-id"),
		}
		return resp, nil
	})

	mockDescribePlansTeoPlan(patches, teoClient, "edgeone-2unuvzjmmn2q")

	res := teo.ResourceTencentCloudTeoPlan()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"plan_type": "personal",
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "edgeone-2unuvzjmmn2q", d.Id())

	// Verify AutoUseVoucher was NOT passed to the CreatePlan request.
	assert.NotNil(t, capturedRequest)
	assert.Nil(t, capturedRequest.AutoUseVoucher)
}

// TestTeoPlan_CreateWithAutoUseVoucherFalse verifies auto_use_voucher="false" is passed as-is.
func TestTeoPlan_CreateWithAutoUseVoucherFalse(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	meta := newMockMetaTeoPlan()
	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(meta.client, "UseTeoV20220901Client", teoClient)

	var capturedRequest *teov20220901.CreatePlanRequest
	patches.ApplyMethodFunc(teoClient, "CreatePlanWithContext", func(_ context.Context, request *teov20220901.CreatePlanRequest) (*teov20220901.CreatePlanResponse, error) {
		capturedRequest = request
		resp := teov20220901.NewCreatePlanResponse()
		resp.Response = &teov20220901.CreatePlanResponseParams{
			PlanId:    ptrStringTeoPlan("edgeone-2unuvzjmmn2q"),
			RequestId: ptrStringTeoPlan("fake-request-id"),
		}
		return resp, nil
	})

	mockDescribePlansTeoPlan(patches, teoClient, "edgeone-2unuvzjmmn2q")

	res := teo.ResourceTencentCloudTeoPlan()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"plan_type":        "standard",
		"auto_use_voucher": "false",
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)

	// Verify AutoUseVoucher="false" was passed to the CreatePlan request.
	assert.NotNil(t, capturedRequest)
	assert.NotNil(t, capturedRequest.AutoUseVoucher)
	assert.Equal(t, "false", *capturedRequest.AutoUseVoucher)
}

// TestTeoPlan_Read verifies Read populates computed fields and does not overwrite auto_use_voucher.
func TestTeoPlan_Read(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	meta := newMockMetaTeoPlan()
	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(meta.client, "UseTeoV20220901Client", teoClient)

	mockDescribePlansTeoPlan(patches, teoClient, "edgeone-2unuvzjmmn2q")

	res := teo.ResourceTencentCloudTeoPlan()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"plan_type":        "personal",
		"auto_use_voucher": "true",
	})
	d.SetId("edgeone-2unuvzjmmn2q")

	err := res.Read(d, meta)
	assert.NoError(t, err)

	// Computed fields are refreshed from the DescribePlans response.
	assert.Equal(t, "edgeone-2unuvzjmmn2q", d.Id())
	assert.Equal(t, "edgeone-2unuvzjmmn2q", d.Get("plan_id").(string))
	assert.Equal(t, "global", d.Get("area").(string))
	assert.Equal(t, "normal", d.Get("status").(string))
	assert.Equal(t, "2026-01-01 00:00:00", d.Get("enabled_time").(string))
	assert.Equal(t, "2026-02-01 00:00:00", d.Get("expired_time").(string))

	// auto_use_voucher is NOT refreshed from the API (Plan struct has no such field),
	// so the locally configured value stays untouched.
	assert.Equal(t, "true", d.Get("auto_use_voucher").(string))
}

// TestTeoPlan_UpdateAutoUseVoucherImmutable verifies changing auto_use_voucher returns an
// immutable error and no update API is invoked.
func TestTeoPlan_UpdateAutoUseVoucherImmutable(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	meta := newMockMetaTeoPlan()
	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(meta.client, "UseTeoV20220901Client", teoClient)

	upgradeCalled := false
	patches.ApplyMethodFunc(teoClient, "UpgradePlanWithContext", func(_ context.Context, request *teov20220901.UpgradePlanRequest) (*teov20220901.UpgradePlanResponse, error) {
		upgradeCalled = true
		resp := teov20220901.NewUpgradePlanResponse()
		resp.Response = &teov20220901.UpgradePlanResponseParams{
			RequestId: ptrStringTeoPlan("fake-request-id"),
		}
		return resp, nil
	})

	renewCalled := false
	patches.ApplyMethodFunc(teoClient, "RenewPlanWithContext", func(_ context.Context, request *teov20220901.RenewPlanRequest) (*teov20220901.RenewPlanResponse, error) {
		renewCalled = true
		resp := teov20220901.NewRenewPlanResponse()
		resp.Response = &teov20220901.RenewPlanResponseParams{
			RequestId: ptrStringTeoPlan("fake-request-id"),
		}
		return resp, nil
	})

	modifyCalled := false
	patches.ApplyMethodFunc(teoClient, "ModifyPlanWithContext", func(_ context.Context, request *teov20220901.ModifyPlanRequest) (*teov20220901.ModifyPlanResponse, error) {
		modifyCalled = true
		resp := teov20220901.NewModifyPlanResponse()
		resp.Response = &teov20220901.ModifyPlanResponseParams{
			RequestId: ptrStringTeoPlan("fake-request-id"),
		}
		return resp, nil
	})

	res := teo.ResourceTencentCloudTeoPlan()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"plan_type":        "personal",
		"auto_use_voucher": "false",
	})
	d.SetId("edgeone-2unuvzjmmn2q")

	// Patch HasChange to simulate a change in auto_use_voucher only.
	patches.ApplyMethodFunc(d, "HasChange", func(key string) bool {
		return key == "auto_use_voucher"
	})

	err := res.Update(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "auto_use_voucher")
	assert.Contains(t, err.Error(), "cannot be changed")

	// No update API should have been invoked.
	assert.False(t, upgradeCalled)
	assert.False(t, renewCalled)
	assert.False(t, modifyCalled)
}

// TestTeoPlan_UpdatePlanType verifies changing plan_type calls UpgradePlan and
// auto_use_voucher unchanged does not trigger the immutable error.
func TestTeoPlan_UpdatePlanType(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	meta := newMockMetaTeoPlan()
	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(meta.client, "UseTeoV20220901Client", teoClient)

	var capturedRequest *teov20220901.UpgradePlanRequest
	upgradeCalled := false
	patches.ApplyMethodFunc(teoClient, "UpgradePlanWithContext", func(_ context.Context, request *teov20220901.UpgradePlanRequest) (*teov20220901.UpgradePlanResponse, error) {
		upgradeCalled = true
		capturedRequest = request
		resp := teov20220901.NewUpgradePlanResponse()
		resp.Response = &teov20220901.UpgradePlanResponseParams{
			RequestId: ptrStringTeoPlan("fake-request-id"),
		}
		return resp, nil
	})

	mockDescribePlansTeoPlan(patches, teoClient, "edgeone-2unuvzjmmn2q")

	res := teo.ResourceTencentCloudTeoPlan()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"plan_type":        "standard",
		"auto_use_voucher": "true",
	})
	d.SetId("edgeone-2unuvzjmmn2q")

	// Patch HasChange to simulate a change in plan_type only (auto_use_voucher unchanged).
	patches.ApplyMethodFunc(d, "HasChange", func(key string) bool {
		return key == "plan_type"
	})

	err := res.Update(d, meta)
	assert.NoError(t, err)

	assert.True(t, upgradeCalled)
	assert.NotNil(t, capturedRequest)
	assert.Equal(t, "edgeone-2unuvzjmmn2q", *capturedRequest.PlanId)
	assert.Equal(t, "standard", *capturedRequest.PlanType)
}

// TestTeoPlan_Delete verifies Delete calls DestroyPlan and succeeds.
func TestTeoPlan_Delete(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	meta := newMockMetaTeoPlan()
	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(meta.client, "UseTeoV20220901Client", teoClient)

	var capturedRequest *teov20220901.DestroyPlanRequest
	destroyCalled := false
	patches.ApplyMethodFunc(teoClient, "DestroyPlanWithContext", func(_ context.Context, request *teov20220901.DestroyPlanRequest) (*teov20220901.DestroyPlanResponse, error) {
		destroyCalled = true
		capturedRequest = request
		resp := teov20220901.NewDestroyPlanResponse()
		resp.Response = &teov20220901.DestroyPlanResponseParams{
			RequestId: ptrStringTeoPlan("fake-request-id"),
		}
		return resp, nil
	})

	res := teo.ResourceTencentCloudTeoPlan()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"plan_type": "personal",
	})
	d.SetId("edgeone-2unuvzjmmn2q")

	err := res.Delete(d, meta)
	assert.NoError(t, err)

	assert.True(t, destroyCalled)
	assert.NotNil(t, capturedRequest)
	assert.Equal(t, "edgeone-2unuvzjmmn2q", *capturedRequest.PlanId)
}

// TestTeoPlan_Schema validates the resource schema definition including auto_use_voucher.
func TestTeoPlan_Schema(t *testing.T) {
	res := teo.ResourceTencentCloudTeoPlan()

	assert.NotNil(t, res)
	assert.NotNil(t, res.Create)
	assert.NotNil(t, res.Read)
	assert.NotNil(t, res.Update)
	assert.NotNil(t, res.Delete)

	assert.Contains(t, res.Schema, "plan_type")
	assert.Contains(t, res.Schema, "auto_use_voucher")
	assert.Contains(t, res.Schema, "prepaid_plan_param")
	assert.Contains(t, res.Schema, "plan_id")

	autoUseVoucher := res.Schema["auto_use_voucher"]
	assert.Equal(t, schema.TypeString, autoUseVoucher.Type)
	assert.True(t, autoUseVoucher.Optional)
	assert.NotNil(t, autoUseVoucher.ValidateFunc)
}
