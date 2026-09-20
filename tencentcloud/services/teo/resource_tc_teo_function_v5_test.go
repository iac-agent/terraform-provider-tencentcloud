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

// mockMetaFunctionV5 implements tccommon.ProviderMeta
type mockMetaFunctionV5 struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaFunctionV5) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaFunctionV5{}

func newMockMetaFunctionV5() *mockMetaFunctionV5 {
	return &mockMetaFunctionV5{client: &connectivity.TencentCloudClient{}}
}

func ptrStringFunctionV5(s string) *string {
	return &s
}

func ptrInt64FunctionV5(i int64) *int64 {
	return &i
}

// go test ./tencentcloud/services/teo/ -run "TestTeoFunctionV5_Create" -v -count=1 -gcflags="all=-l"
func TestTeoFunctionV5_Create(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaFunctionV5().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "CreateFunctionWithContext", func(_ context.Context, request *teov20220901.CreateFunctionRequest) (*teov20220901.CreateFunctionResponse, error) {
		assert.Equal(t, "zone-test123", *request.ZoneId)
		assert.Equal(t, "test-func", *request.Name)
		assert.Equal(t, "console.log(123)", *request.Content)
		assert.Equal(t, "test remark", *request.Remark)

		resp := teov20220901.NewCreateFunctionResponse()
		resp.Response = &teov20220901.CreateFunctionResponseParams{
			FunctionId: ptrStringFunctionV5("ef-test456"),
			RequestId:  ptrStringFunctionV5("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(teoClient, "DescribeFunctionsWithContext", func(_ context.Context, request *teov20220901.DescribeFunctionsRequest) (*teov20220901.DescribeFunctionsResponse, error) {
		resp := teov20220901.NewDescribeFunctionsResponse()
		resp.Response = &teov20220901.DescribeFunctionsResponseParams{
			TotalCount: ptrInt64FunctionV5(1),
			Functions: []*teov20220901.Function{
				{
					FunctionId: ptrStringFunctionV5("ef-test456"),
					ZoneId:     ptrStringFunctionV5("zone-test123"),
					Name:       ptrStringFunctionV5("test-func"),
					Remark:     ptrStringFunctionV5("test remark"),
					Content:    ptrStringFunctionV5("console.log(123)"),
					Domain:     ptrStringFunctionV5("test.edgeone.app"),
					DomainComplianceRestrictions: []*teov20220901.ComplianceRestriction{
						{
							Reason: ptrStringFunctionV5("ICP_RECORD_REQUIRED"),
							Region: ptrStringFunctionV5("CN"),
						},
					},
					CreateTime: ptrStringFunctionV5("2024-01-01T00:00:00Z"),
					UpdateTime: ptrStringFunctionV5("2024-01-01T00:00:00Z"),
				},
			},
			RequestId: ptrStringFunctionV5("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaFunctionV5()
	res := teo.ResourceTencentCloudTeoFunctionV5()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-test123",
		"name":    "test-func",
		"content": "console.log(123)",
		"remark":  "test remark",
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-test123#ef-test456", d.Id())
	assert.Equal(t, "test.edgeone.app", d.Get("domain"))
	assert.Equal(t, "2024-01-01T00:00:00Z", d.Get("create_time"))
}

// go test ./tencentcloud/services/teo/ -run "TestTeoFunctionV5_Read" -v -count=1 -gcflags="all=-l"
func TestTeoFunctionV5_Read(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaFunctionV5().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeFunctions", func(request *teov20220901.DescribeFunctionsRequest) (*teov20220901.DescribeFunctionsResponse, error) {
		assert.Equal(t, "zone-test123", *request.ZoneId)
		assert.Equal(t, "ef-test456", *request.FunctionIds[0])

		resp := teov20220901.NewDescribeFunctionsResponse()
		resp.Response = &teov20220901.DescribeFunctionsResponseParams{
			TotalCount: ptrInt64FunctionV5(1),
			Functions: []*teov20220901.Function{
				{
					FunctionId: ptrStringFunctionV5("ef-test456"),
					ZoneId:     ptrStringFunctionV5("zone-test123"),
					Name:       ptrStringFunctionV5("test-func"),
					Remark:     ptrStringFunctionV5("test remark"),
					Content:    ptrStringFunctionV5("console.log(123)"),
					Domain:     ptrStringFunctionV5("test.edgeone.app"),
					DomainComplianceRestrictions: []*teov20220901.ComplianceRestriction{
						{
							Reason: ptrStringFunctionV5("ICP_RECORD_REQUIRED"),
							Region: ptrStringFunctionV5("CN"),
						},
					},
					CreateTime: ptrStringFunctionV5("2024-01-01T00:00:00Z"),
					UpdateTime: ptrStringFunctionV5("2024-01-02T00:00:00Z"),
				},
			},
			RequestId: ptrStringFunctionV5("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaFunctionV5()
	res := teo.ResourceTencentCloudTeoFunctionV5()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-test123",
		"name":    "test-func",
		"content": "console.log(123)",
		"remark":  "test remark",
	})
	d.SetId("zone-test123#ef-test456")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "ef-test456", d.Get("function_id"))
	assert.Equal(t, "test.edgeone.app", d.Get("domain"))
	assert.Equal(t, "2024-01-01T00:00:00Z", d.Get("create_time"))
	assert.Equal(t, "2024-01-02T00:00:00Z", d.Get("update_time"))

	complianceRestrictions := d.Get("domain_compliance_restrictions").([]interface{})
	assert.Equal(t, 1, len(complianceRestrictions))
	restriction := complianceRestrictions[0].(map[string]interface{})
	assert.Equal(t, "ICP_RECORD_REQUIRED", restriction["reason"])
	assert.Equal(t, "CN", restriction["region"])
}

// go test ./tencentcloud/services/teo/ -run "TestTeoFunctionV5_ReadNotFound" -v -count=1 -gcflags="all=-l"
func TestTeoFunctionV5_ReadNotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaFunctionV5().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeFunctions", func(request *teov20220901.DescribeFunctionsRequest) (*teov20220901.DescribeFunctionsResponse, error) {
		resp := teov20220901.NewDescribeFunctionsResponse()
		resp.Response = &teov20220901.DescribeFunctionsResponseParams{
			TotalCount: ptrInt64FunctionV5(0),
			Functions:  []*teov20220901.Function{},
			RequestId:  ptrStringFunctionV5("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaFunctionV5()
	res := teo.ResourceTencentCloudTeoFunctionV5()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-test123",
		"name":    "test-func",
		"content": "console.log(123)",
		"remark":  "test remark",
	})
	d.SetId("zone-test123#ef-test456")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}

// go test ./tencentcloud/services/teo/ -run "TestTeoFunctionV5_Update" -v -count=1 -gcflags="all=-l"
func TestTeoFunctionV5_Update(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaFunctionV5().client, "UseTeoV20220901Client", teoClient)

	modifyCalled := false
	patches.ApplyMethodFunc(teoClient, "ModifyFunctionWithContext", func(_ context.Context, request *teov20220901.ModifyFunctionRequest) (*teov20220901.ModifyFunctionResponse, error) {
		modifyCalled = true
		assert.Equal(t, "zone-test123", *request.ZoneId)
		assert.Equal(t, "ef-test456", *request.FunctionId)
		assert.Equal(t, "console.log(456)", *request.Content)
		assert.Equal(t, "updated remark", *request.Remark)

		resp := teov20220901.NewModifyFunctionResponse()
		resp.Response = &teov20220901.ModifyFunctionResponseParams{
			RequestId: ptrStringFunctionV5("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(teoClient, "DescribeFunctions", func(request *teov20220901.DescribeFunctionsRequest) (*teov20220901.DescribeFunctionsResponse, error) {
		resp := teov20220901.NewDescribeFunctionsResponse()
		resp.Response = &teov20220901.DescribeFunctionsResponseParams{
			TotalCount: ptrInt64FunctionV5(1),
			Functions: []*teov20220901.Function{
				{
					FunctionId: ptrStringFunctionV5("ef-test456"),
					ZoneId:     ptrStringFunctionV5("zone-test123"),
					Name:       ptrStringFunctionV5("test-func"),
					Remark:     ptrStringFunctionV5("updated remark"),
					Content:    ptrStringFunctionV5("console.log(456)"),
					Domain:     ptrStringFunctionV5("test.edgeone.app"),
					DomainComplianceRestrictions: []*teov20220901.ComplianceRestriction{
						{
							Reason: ptrStringFunctionV5("ICP_RECORD_REQUIRED"),
							Region: ptrStringFunctionV5("CN"),
						},
					},
					CreateTime: ptrStringFunctionV5("2024-01-01T00:00:00Z"),
					UpdateTime: ptrStringFunctionV5("2024-01-02T00:00:00Z"),
				},
			},
			RequestId: ptrStringFunctionV5("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaFunctionV5()
	res := teo.ResourceTencentCloudTeoFunctionV5()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-test123",
		"name":    "test-func",
		"content": "console.log(456)",
		"remark":  "updated remark",
	})
	d.SetId("zone-test123#ef-test456")

	err := res.Update(d, meta)
	assert.NoError(t, err)
	assert.True(t, modifyCalled, "ModifyFunctionWithContext should be called")
	assert.Equal(t, "console.log(456)", d.Get("content"))
	assert.Equal(t, "updated remark", d.Get("remark"))
}

// go test ./tencentcloud/services/teo/ -run "TestTeoFunctionV5_Delete" -v -count=1 -gcflags="all=-l"
func TestTeoFunctionV5_Delete(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaFunctionV5().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DeleteFunctionWithContext", func(_ context.Context, request *teov20220901.DeleteFunctionRequest) (*teov20220901.DeleteFunctionResponse, error) {
		assert.Equal(t, "zone-test123", *request.ZoneId)
		assert.Equal(t, "ef-test456", *request.FunctionId)

		resp := teov20220901.NewDeleteFunctionResponse()
		resp.Response = &teov20220901.DeleteFunctionResponseParams{
			RequestId: ptrStringFunctionV5("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaFunctionV5()
	res := teo.ResourceTencentCloudTeoFunctionV5()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-test123",
		"name":    "test-func",
		"content": "console.log(123)",
		"remark":  "test remark",
	})
	d.SetId("zone-test123#ef-test456")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
}
