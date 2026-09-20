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

type mockMetaForTeoFunctionV6 struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaForTeoFunctionV6) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaForTeoFunctionV6{}

func newMockMetaForTeoFunctionV6() *mockMetaForTeoFunctionV6 {
	return &mockMetaForTeoFunctionV6{client: &connectivity.TencentCloudClient{}}
}

func ptrStrFuncV6(s string) *string {
	return &s
}

func ptrInt64FuncV6(v int64) *int64 {
	return &v
}

// go test ./tencentcloud/services/teo/ -run "TestTeoFunctionV6" -v -count=1 -gcflags="all=-l"

// TestTeoFunctionV6_Schema tests the schema definition
func TestTeoFunctionV6_Schema(t *testing.T) {
	res := teo.ResourceTencentCloudTeoFunctionV6()

	assert.NotNil(t, res)
	assert.NotNil(t, res.Importer)

	// Required fields
	assert.Contains(t, res.Schema, "zone_id")
	assert.True(t, res.Schema["zone_id"].Required)
	assert.True(t, res.Schema["zone_id"].ForceNew)

	assert.Contains(t, res.Schema, "name")
	assert.True(t, res.Schema["name"].Required)
	assert.True(t, res.Schema["name"].ForceNew)

	assert.Contains(t, res.Schema, "content")
	assert.True(t, res.Schema["content"].Required)

	// Optional fields
	assert.Contains(t, res.Schema, "remark")
	assert.True(t, res.Schema["remark"].Optional)

	// Computed fields
	assert.Contains(t, res.Schema, "function_id")
	assert.True(t, res.Schema["function_id"].Computed)

	assert.Contains(t, res.Schema, "domain")
	assert.True(t, res.Schema["domain"].Computed)

	assert.Contains(t, res.Schema, "domain_compliance_restrictions")
	assert.True(t, res.Schema["domain_compliance_restrictions"].Computed)

	assert.Contains(t, res.Schema, "create_time")
	assert.True(t, res.Schema["create_time"].Computed)

	assert.Contains(t, res.Schema, "update_time")
	assert.True(t, res.Schema["update_time"].Computed)
}

// TestTeoFunctionV6_Read_Success tests Read populates fields from Function
func TestTeoFunctionV6_Read_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForTeoFunctionV6().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeFunctions", func(request *teov20220901.DescribeFunctionsRequest) (*teov20220901.DescribeFunctionsResponse, error) {
		resp := teov20220901.NewDescribeFunctionsResponse()
		zoneId := "zone-2qtuhspy7cr6"
		resp.Response = &teov20220901.DescribeFunctionsResponseParams{
			TotalCount: ptrInt64FuncV6(1),
			Functions: []*teov20220901.Function{
				{
					FunctionId: ptrStrFuncV6("ef-xxxx"),
					ZoneId:     ptrStrFuncV6(zoneId),
					Name:       ptrStrFuncV6("test"),
					Remark:     ptrStrFuncV6("test remark"),
					Content:    ptrStrFuncV6("addEventListener('fetch', e => {});"),
					Domain:     ptrStrFuncV6("test.example.com"),
					CreateTime: ptrStrFuncV6("2024-01-01T00:00:00Z"),
					UpdateTime: ptrStrFuncV6("2024-01-02T00:00:00Z"),
					DomainComplianceRestrictions: []*teov20220901.ComplianceRestriction{
						{
							Reason: ptrStrFuncV6("icp"),
							Region: ptrStrFuncV6("CN"),
						},
					},
				},
			},
			RequestId: ptrStrFuncV6("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForTeoFunctionV6()
	res := teo.ResourceTencentCloudTeoFunctionV6()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-2qtuhspy7cr6",
		"name":    "test",
		"content": "addEventListener('fetch', e => {});",
		"remark":  "test remark",
	})
	d.SetId("zone-2qtuhspy7cr6#ef-xxxx")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "ef-xxxx", d.Get("function_id"))
	assert.Equal(t, "test", d.Get("name"))
	assert.Equal(t, "test remark", d.Get("remark"))
	assert.Equal(t, "test.example.com", d.Get("domain"))
	assert.Equal(t, 1, d.Get("domain_compliance_restrictions.#"))
	assert.Equal(t, "icp", d.Get("domain_compliance_restrictions.0.reason"))
	assert.Equal(t, "CN", d.Get("domain_compliance_restrictions.0.region"))
	assert.Equal(t, "2024-01-01T00:00:00Z", d.Get("create_time"))
	assert.Equal(t, "2024-01-02T00:00:00Z", d.Get("update_time"))
}

// TestTeoFunctionV6_Read_NotFound tests Read removes resource when not found
func TestTeoFunctionV6_Read_NotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForTeoFunctionV6().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeFunctions", func(request *teov20220901.DescribeFunctionsRequest) (*teov20220901.DescribeFunctionsResponse, error) {
		resp := teov20220901.NewDescribeFunctionsResponse()
		resp.Response = &teov20220901.DescribeFunctionsResponseParams{
			TotalCount: ptrInt64FuncV6(0),
			Functions:  []*teov20220901.Function{},
			RequestId:  ptrStrFuncV6("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForTeoFunctionV6()
	res := teo.ResourceTencentCloudTeoFunctionV6()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-2qtuhspy7cr6",
		"name":    "test",
		"content": "addEventListener('fetch', e => {});",
	})
	d.SetId("zone-2qtuhspy7cr6#ef-xxxx")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}

// TestTeoFunctionV6_Create_Success tests Create sets ID and polls deployment
func TestTeoFunctionV6_Create_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForTeoFunctionV6().client, "UseTeoV20220901Client", teoClient)

	zoneId := "zone-2qtuhspy7cr6"
	functionId := "ef-xxxx"

	patches.ApplyMethodFunc(teoClient, "CreateFunctionWithContext", func(ctx context.Context, request *teov20220901.CreateFunctionRequest) (*teov20220901.CreateFunctionResponse, error) {
		resp := teov20220901.NewCreateFunctionResponse()
		resp.Response = &teov20220901.CreateFunctionResponseParams{
			FunctionId: ptrStrFuncV6(functionId),
			RequestId:  ptrStrFuncV6("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(teoClient, "DescribeFunctionsWithContext", func(ctx context.Context, request *teov20220901.DescribeFunctionsRequest) (*teov20220901.DescribeFunctionsResponse, error) {
		resp := teov20220901.NewDescribeFunctionsResponse()
		resp.Response = &teov20220901.DescribeFunctionsResponseParams{
			TotalCount: ptrInt64FuncV6(1),
			Functions: []*teov20220901.Function{
				{
					FunctionId: ptrStrFuncV6(functionId),
					ZoneId:     ptrStrFuncV6(zoneId),
					Name:       ptrStrFuncV6("test"),
					Remark:     ptrStrFuncV6("test remark"),
					Content:    ptrStrFuncV6("addEventListener('fetch', e => {});"),
					Domain:     ptrStrFuncV6("test.example.com"),
					CreateTime: ptrStrFuncV6("2024-01-01T00:00:00Z"),
					UpdateTime: ptrStrFuncV6("2024-01-02T00:00:00Z"),
				},
			},
			RequestId: ptrStrFuncV6("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForTeoFunctionV6()
	res := teo.ResourceTencentCloudTeoFunctionV6()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": zoneId,
		"name":    "test",
		"content": "addEventListener('fetch', e => {});",
		"remark":  "test remark",
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-2qtuhspy7cr6#ef-xxxx", d.Id())
	assert.Equal(t, "ef-xxxx", d.Get("function_id"))
	assert.Equal(t, "test.example.com", d.Get("domain"))
}

// TestTeoFunctionV6_Update_Success tests Update calls ModifyFunction
func TestTeoFunctionV6_Update_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForTeoFunctionV6().client, "UseTeoV20220901Client", teoClient)

	zoneId := "zone-2qtuhspy7cr6"
	functionId := "ef-xxxx"

	patches.ApplyMethodFunc(teoClient, "ModifyFunctionWithContext", func(ctx context.Context, request *teov20220901.ModifyFunctionRequest) (*teov20220901.ModifyFunctionResponse, error) {
		assert.Equal(t, zoneId, *request.ZoneId)
		assert.Equal(t, functionId, *request.FunctionId)
		assert.Equal(t, "updated remark", *request.Remark)
		resp := teov20220901.NewModifyFunctionResponse()
		resp.Response = &teov20220901.ModifyFunctionResponseParams{
			RequestId: ptrStrFuncV6("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(teoClient, "DescribeFunctions", func(request *teov20220901.DescribeFunctionsRequest) (*teov20220901.DescribeFunctionsResponse, error) {
		resp := teov20220901.NewDescribeFunctionsResponse()
		resp.Response = &teov20220901.DescribeFunctionsResponseParams{
			TotalCount: ptrInt64FuncV6(1),
			Functions: []*teov20220901.Function{
				{
					FunctionId: ptrStrFuncV6(functionId),
					ZoneId:     ptrStrFuncV6(zoneId),
					Name:       ptrStrFuncV6("test"),
					Remark:     ptrStrFuncV6("updated remark"),
					Content:    ptrStrFuncV6("addEventListener('fetch', e => {});"),
					Domain:     ptrStrFuncV6("test.example.com"),
					CreateTime: ptrStrFuncV6("2024-01-01T00:00:00Z"),
					UpdateTime: ptrStrFuncV6("2024-01-02T00:00:00Z"),
				},
			},
			RequestId: ptrStrFuncV6("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForTeoFunctionV6()
	res := teo.ResourceTencentCloudTeoFunctionV6()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": zoneId,
		"name":    "test",
		"content": "addEventListener('fetch', e => {});",
		"remark":  "updated remark",
	})
	d.SetId(zoneId + "#" + functionId)

	// simulate a change on remark
	_ = d.Set("remark", "updated remark")

	err := res.Update(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "updated remark", d.Get("remark"))
}

// TestTeoFunctionV6_Delete_Success tests Delete calls DeleteFunction
func TestTeoFunctionV6_Delete_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForTeoFunctionV6().client, "UseTeoV20220901Client", teoClient)

	zoneId := "zone-2qtuhspy7cr6"
	functionId := "ef-xxxx"

	deleted := false
	patches.ApplyMethodFunc(teoClient, "DeleteFunctionWithContext", func(ctx context.Context, request *teov20220901.DeleteFunctionRequest) (*teov20220901.DeleteFunctionResponse, error) {
		assert.Equal(t, zoneId, *request.ZoneId)
		assert.Equal(t, functionId, *request.FunctionId)
		deleted = true
		resp := teov20220901.NewDeleteFunctionResponse()
		resp.Response = &teov20220901.DeleteFunctionResponseParams{
			RequestId: ptrStrFuncV6("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForTeoFunctionV6()
	res := teo.ResourceTencentCloudTeoFunctionV6()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": zoneId,
		"name":    "test",
		"content": "addEventListener('fetch', e => {});",
	})
	d.SetId(zoneId + "#" + functionId)

	err := res.Delete(d, meta)
	assert.NoError(t, err)
	assert.True(t, deleted)
}
