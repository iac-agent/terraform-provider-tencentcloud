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
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/teo"
)

func ptrString(s string) *string {
	return &s
}

func TestTeoFunctionV3Schema(t *testing.T) {
	res := teo.ResourceTencentCloudTeoFunctionV3()

	assert.NotNil(t, res)
	assert.NotNil(t, res.Create)
	assert.NotNil(t, res.Read)
	assert.NotNil(t, res.Update)
	assert.NotNil(t, res.Delete)
	assert.NotNil(t, res.Importer)

	expectedFields := []string{"zone_id", "name", "content", "remark", "function_id", "domain", "create_time", "update_time"}
	for _, field := range expectedFields {
		assert.Contains(t, res.Schema, field, "schema should contain field: %s", field)
	}

	// Required fields
	assert.True(t, res.Schema["zone_id"].Required)
	assert.True(t, res.Schema["zone_id"].ForceNew)
	assert.True(t, res.Schema["name"].Required)
	assert.True(t, res.Schema["content"].Required)

	// Optional fields
	assert.True(t, res.Schema["remark"].Optional)

	// Computed fields
	assert.True(t, res.Schema["function_id"].Computed)
	assert.True(t, res.Schema["domain"].Computed)
	assert.True(t, res.Schema["create_time"].Computed)
	assert.True(t, res.Schema["update_time"].Computed)
}

func TestTeoFunctionV3Create_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	// Patch UseTeoV20220901Client to return a non-nil client
	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(&connectivity.TencentCloudClient{}, "UseTeoV20220901Client", teoClient)

	// Patch CreateFunctionWithContext to return success
	patches.ApplyMethodFunc(teoClient, "CreateFunctionWithContext", func(_ interface{}, request *teov20220901.CreateFunctionRequest) (*teov20220901.CreateFunctionResponse, error) {
		resp := teov20220901.NewCreateFunctionResponse()
		resp.Response = &teov20220901.CreateFunctionResponseParams{
			FunctionId: ptrString("ef-test123456"),
			RequestId:  ptrString("fake-request-id"),
		}
		return resp, nil
	})

	// Patch DescribeFunctionsWithContext for state refresh
	patches.ApplyMethodFunc(teoClient, "DescribeFunctionsWithContext", func(_ interface{}, request *teov20220901.DescribeFunctionsRequest) (*teov20220901.DescribeFunctionsResponse, error) {
		resp := teov20220901.NewDescribeFunctionsResponse()
		resp.Response = &teov20220901.DescribeFunctionsResponseParams{
			TotalCount: helper.Int64(1),
			Functions: []*teov20220901.Function{
				{
					FunctionId: ptrString("ef-test123456"),
					ZoneId:     ptrString("zone-2qtuhspy7cr6"),
					Name:       ptrString("test-function"),
					Remark:     ptrString("test remark"),
					Content:    ptrString("addEventListener('fetch', e => { e.respondWith(new Response('Hello')); });"),
					Domain:     ptrString("test-function.zone-2qtuhspy7cr6.edgeone.app"),
					CreateTime: ptrString("2024-01-01T00:00:00Z"),
					UpdateTime: ptrString("2024-01-01T00:00:00Z"),
				},
			},
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := &mockMeta{client: &connectivity.TencentCloudClient{}}
	res := teo.ResourceTencentCloudTeoFunctionV3()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-2qtuhspy7cr6",
		"name":    "test-function",
		"content": "addEventListener('fetch', e => { e.respondWith(new Response('Hello')); });",
		"remark":  "test remark",
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-2qtuhspy7cr6#ef-test123456", d.Id())
	assert.Equal(t, "ef-test123456", d.Get("function_id"))
}

func TestTeoFunctionV3Create_EmptyFunctionId(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(&connectivity.TencentCloudClient{}, "UseTeoV20220901Client", teoClient)

	// Patch CreateFunctionWithContext to return empty FunctionId
	patches.ApplyMethodFunc(teoClient, "CreateFunctionWithContext", func(_ interface{}, request *teov20220901.CreateFunctionRequest) (*teov20220901.CreateFunctionResponse, error) {
		resp := teov20220901.NewCreateFunctionResponse()
		resp.Response = &teov20220901.CreateFunctionResponseParams{
			FunctionId: ptrString(""),
			RequestId:  ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := &mockMeta{client: &connectivity.TencentCloudClient{}}
	res := teo.ResourceTencentCloudTeoFunctionV3()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-2qtuhspy7cr6",
		"name":    "test-function",
		"content": "test content",
	})

	err := res.Create(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "FunctionId is empty")
}

func TestTeoFunctionV3Create_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(&connectivity.TencentCloudClient{}, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "CreateFunctionWithContext", func(_ interface{}, request *teov20220901.CreateFunctionRequest) (*teov20220901.CreateFunctionResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=InvalidParameter, Message=Invalid zone_id")
	})

	meta := &mockMeta{client: &connectivity.TencentCloudClient{}}
	res := teo.ResourceTencentCloudTeoFunctionV3()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-invalid",
		"name":    "test-function",
		"content": "test content",
	})

	err := res.Create(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "InvalidParameter")
}

func TestTeoFunctionV3Read_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(&connectivity.TencentCloudClient{}, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeFunctionsWithContext", func(_ interface{}, request *teov20220901.DescribeFunctionsRequest) (*teov20220901.DescribeFunctionsResponse, error) {
		resp := teov20220901.NewDescribeFunctionsResponse()
		resp.Response = &teov20220901.DescribeFunctionsResponseParams{
			TotalCount: helper.Int64(1),
			Functions: []*teov20220901.Function{
				{
					FunctionId: ptrString("ef-test123456"),
					ZoneId:     ptrString("zone-2qtuhspy7cr6"),
					Name:       ptrString("test-function"),
					Remark:     ptrString("test remark"),
					Content:    ptrString("test content"),
					Domain:     ptrString("test-function.zone-2qtuhspy7cr6.edgeone.app"),
					CreateTime: ptrString("2024-01-01T00:00:00Z"),
					UpdateTime: ptrString("2024-01-01T00:00:00Z"),
				},
			},
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := &mockMeta{client: &connectivity.TencentCloudClient{}}
	res := teo.ResourceTencentCloudTeoFunctionV3()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-2qtuhspy7cr6",
		"name":    "test-function",
		"content": "test content",
	})
	d.SetId("zone-2qtuhspy7cr6#ef-test123456")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-2qtuhspy7cr6#ef-test123456", d.Id())
	assert.Equal(t, "ef-test123456", d.Get("function_id"))
	assert.Equal(t, "test-function", d.Get("name"))
	assert.Equal(t, "test remark", d.Get("remark"))
	assert.Equal(t, "test content", d.Get("content"))
	assert.Equal(t, "test-function.zone-2qtuhspy7cr6.edgeone.app", d.Get("domain"))
}

func TestTeoFunctionV3Read_NotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(&connectivity.TencentCloudClient{}, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeFunctionsWithContext", func(_ interface{}, request *teov20220901.DescribeFunctionsRequest) (*teov20220901.DescribeFunctionsResponse, error) {
		resp := teov20220901.NewDescribeFunctionsResponse()
		resp.Response = &teov20220901.DescribeFunctionsResponseParams{
			TotalCount: helper.Int64(0),
			Functions:  []*teov20220901.Function{},
			RequestId:  ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := &mockMeta{client: &connectivity.TencentCloudClient{}}
	res := teo.ResourceTencentCloudTeoFunctionV3()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-2qtuhspy7cr6",
		"name":    "test-function",
		"content": "test content",
	})
	d.SetId("zone-2qtuhspy7cr6#ef-test123456")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}

func TestTeoFunctionV3Update_ImmutableName(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(&connectivity.TencentCloudClient{}, "UseTeoV20220901Client", teoClient)

	// Patch DescribeFunctionsWithContext for read
	patches.ApplyMethodFunc(teoClient, "DescribeFunctionsWithContext", func(_ interface{}, request *teov20220901.DescribeFunctionsRequest) (*teov20220901.DescribeFunctionsResponse, error) {
		resp := teov20220901.NewDescribeFunctionsResponse()
		resp.Response = &teov20220901.DescribeFunctionsResponseParams{
			TotalCount: helper.Int64(1),
			Functions: []*teov20220901.Function{
				{
					FunctionId: ptrString("ef-test123456"),
					ZoneId:     ptrString("zone-2qtuhspy7cr6"),
					Name:       ptrString("test-function"),
					Remark:     ptrString("updated remark"),
					Content:    ptrString("test content"),
					Domain:     ptrString("test-function.zone-2qtuhspy7cr6.edgeone.app"),
					CreateTime: ptrString("2024-01-01T00:00:00Z"),
					UpdateTime: ptrString("2024-01-01T00:00:00Z"),
				},
			},
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := &mockMeta{client: &connectivity.TencentCloudClient{}}
	res := teo.ResourceTencentCloudTeoFunctionV3()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-2qtuhspy7cr6",
		"name":    "test-function-changed",
		"content": "test content",
		"remark":  "updated remark",
	})
	d.SetId("zone-2qtuhspy7cr6#ef-test123456")

	// Simulate that name has changed
	_ = d.Set("name", "test-function-changed")

	err := res.Update(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "argument `name` cannot be changed")
}

func TestTeoFunctionV3Update_ContentChange(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(&connectivity.TencentCloudClient{}, "UseTeoV20220901Client", teoClient)

	// Patch ModifyFunctionWithContext to return success
	patches.ApplyMethodFunc(teoClient, "ModifyFunctionWithContext", func(_ interface{}, request *teov20220901.ModifyFunctionRequest) (*teov20220901.ModifyFunctionResponse, error) {
		resp := teov20220901.NewModifyFunctionResponse()
		resp.Response = &teov20220901.ModifyFunctionResponseParams{
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	// Patch DescribeFunctionsWithContext for read
	patches.ApplyMethodFunc(teoClient, "DescribeFunctionsWithContext", func(_ interface{}, request *teov20220901.DescribeFunctionsRequest) (*teov20220901.DescribeFunctionsResponse, error) {
		resp := teov20220901.NewDescribeFunctionsResponse()
		resp.Response = &teov20220901.DescribeFunctionsResponseParams{
			TotalCount: helper.Int64(1),
			Functions: []*teov20220901.Function{
				{
					FunctionId: ptrString("ef-test123456"),
					ZoneId:     ptrString("zone-2qtuhspy7cr6"),
					Name:       ptrString("test-function"),
					Remark:     ptrString("test remark"),
					Content:    ptrString("updated content"),
					Domain:     ptrString("test-function.zone-2qtuhspy7cr6.edgeone.app"),
					CreateTime: ptrString("2024-01-01T00:00:00Z"),
					UpdateTime: ptrString("2024-01-01T00:00:00Z"),
				},
			},
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := &mockMeta{client: &connectivity.TencentCloudClient{}}
	res := teo.ResourceTencentCloudTeoFunctionV3()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-2qtuhspy7cr6",
		"name":    "test-function",
		"content": "updated content",
		"remark":  "test remark",
	})
	d.SetId("zone-2qtuhspy7cr6#ef-test123456")

	err := res.Update(d, meta)
	assert.NoError(t, err)
}

func TestTeoFunctionV3Delete_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(&connectivity.TencentCloudClient{}, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DeleteFunctionWithContext", func(_ interface{}, request *teov20220901.DeleteFunctionRequest) (*teov20220901.DeleteFunctionResponse, error) {
		resp := teov20220901.NewDeleteFunctionResponse()
		resp.Response = &teov20220901.DeleteFunctionResponseParams{
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := &mockMeta{client: &connectivity.TencentCloudClient{}}
	res := teo.ResourceTencentCloudTeoFunctionV3()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-2qtuhspy7cr6",
		"name":    "test-function",
		"content": "test content",
	})
	d.SetId("zone-2qtuhspy7cr6#ef-test123456")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
}

func TestTeoFunctionV3Delete_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(&connectivity.TencentCloudClient{}, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DeleteFunctionWithContext", func(_ interface{}, request *teov20220901.DeleteFunctionRequest) (*teov20220901.DeleteFunctionResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=ResourceNotFound, Message=Function not found")
	})

	meta := &mockMeta{client: &connectivity.TencentCloudClient{}}
	res := teo.ResourceTencentCloudTeoFunctionV3()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-2qtuhspy7cr6",
		"name":    "test-function",
		"content": "test content",
	})
	d.SetId("zone-2qtuhspy7cr6#ef-test123456")

	err := res.Delete(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ResourceNotFound")
}
