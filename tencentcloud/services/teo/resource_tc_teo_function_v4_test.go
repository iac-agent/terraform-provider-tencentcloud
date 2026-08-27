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
	svcteo "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/teo"
)

// go test ./tencentcloud/services/teo/ -run "TestTeoFunctionV4_" -v -count=1 -gcflags="all=-l"

// mockMetaFunctionV4 implements tccommon.ProviderMeta
type mockMetaFunctionV4 struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaFunctionV4) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaFunctionV4{}

func newMockMetaFunctionV4() *mockMetaFunctionV4 {
	return &mockMetaFunctionV4{client: &connectivity.TencentCloudClient{}}
}

func ptrStringFunctionV4(s string) *string {
	return &s
}

func ptrInt64FunctionV4(i int64) *int64 {
	return &i
}

// mockDescribeFunctionResponse builds a DescribeFunctions response with a single function.
func mockDescribeFunctionResponse(function *teov20220901.Function) (*teov20220901.DescribeFunctionsResponse, error) {
	resp := teov20220901.NewDescribeFunctionsResponse()
	resp.Response = &teov20220901.DescribeFunctionsResponseParams{
		TotalCount: ptrInt64FunctionV4(1),
		Functions:  []*teov20220901.Function{function},
		RequestId:  ptrStringFunctionV4("fake-request-id"),
	}
	return resp, nil
}

func newFunctionV4(id, zoneId, name, remark, content, domain, createTime, updateTime string) *teov20220901.Function {
	return &teov20220901.Function{
		FunctionId: ptrStringFunctionV4(id),
		ZoneId:     ptrStringFunctionV4(zoneId),
		Name:       ptrStringFunctionV4(name),
		Remark:     ptrStringFunctionV4(remark),
		Content:    ptrStringFunctionV4(content),
		Domain:     ptrStringFunctionV4(domain),
		CreateTime: ptrStringFunctionV4(createTime),
		UpdateTime: ptrStringFunctionV4(updateTime),
	}
}

// TestTeoFunctionV4_Create tests creating a teo_function_v4 resource successfully.
func TestTeoFunctionV4_Create(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaFunctionV4().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "CreateFunctionWithContext", func(_ context.Context, request *teov20220901.CreateFunctionRequest) (*teov20220901.CreateFunctionResponse, error) {
		assert.Equal(t, "zone-test123", *request.ZoneId)
		assert.Equal(t, "my-func", *request.Name)
		assert.Equal(t, "test remark", *request.Remark)
		assert.Equal(t, "addEventListener('fetch', e => {});", *request.Content)

		resp := teov20220901.NewCreateFunctionResponse()
		resp.Response = &teov20220901.CreateFunctionResponseParams{
			FunctionId: ptrStringFunctionV4("ef-test456"),
			RequestId:  ptrStringFunctionV4("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(teoClient, "DescribeFunctions", func(request *teov20220901.DescribeFunctionsRequest) (*teov20220901.DescribeFunctionsResponse, error) {
		assert.Equal(t, "zone-test123", *request.ZoneId)
		assert.Equal(t, 1, len(request.FunctionIds))
		assert.Equal(t, "ef-test456", *request.FunctionIds[0])
		return mockDescribeFunctionResponse(newFunctionV4(
			"ef-test456", "zone-test123", "my-func-zone-test123",
			"test remark", "addEventListener('fetch', e => {});",
			"my-func.example.com", "2024-01-01T00:00:00Z", "2024-01-01T00:00:00Z",
		))
	})

	meta := newMockMetaFunctionV4()
	res := svcteo.ResourceTencentCloudTeoFunctionV4()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-test123",
		"name":    "my-func",
		"remark":  "test remark",
		"content": "addEventListener('fetch', e => {});",
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-test123#ef-test456", d.Id())
	assert.Equal(t, "ef-test456", d.Get("function_id"))
	assert.Equal(t, "my-func", d.Get("name"))
	assert.Equal(t, "my-func.example.com", d.Get("domain"))
}

// TestTeoFunctionV4_CreateEmptyFunctionId tests that create returns an error when FunctionId is empty.
func TestTeoFunctionV4_CreateEmptyFunctionId(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaFunctionV4().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "CreateFunctionWithContext", func(_ context.Context, request *teov20220901.CreateFunctionRequest) (*teov20220901.CreateFunctionResponse, error) {
		resp := teov20220901.NewCreateFunctionResponse()
		resp.Response = &teov20220901.CreateFunctionResponseParams{
			FunctionId: ptrStringFunctionV4(""),
			RequestId:  ptrStringFunctionV4("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaFunctionV4()
	res := svcteo.ResourceTencentCloudTeoFunctionV4()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-test123",
		"name":    "my-func",
		"content": "addEventListener('fetch', e => {});",
	})

	err := res.Create(d, meta)
	assert.Error(t, err)
	assert.Equal(t, "", d.Id())
}

// TestTeoFunctionV4_CreateNilResponse tests that create returns an error when the response is nil.
func TestTeoFunctionV4_CreateNilResponse(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaFunctionV4().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "CreateFunctionWithContext", func(_ context.Context, request *teov20220901.CreateFunctionRequest) (*teov20220901.CreateFunctionResponse, error) {
		resp := teov20220901.NewCreateFunctionResponse()
		resp.Response = nil
		return resp, nil
	})

	meta := newMockMetaFunctionV4()
	res := svcteo.ResourceTencentCloudTeoFunctionV4()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-test123",
		"name":    "my-func",
		"content": "addEventListener('fetch', e => {});",
	})

	err := res.Create(d, meta)
	assert.Error(t, err)
	assert.Equal(t, "", d.Id())
}

// TestTeoFunctionV4_Read tests reading an existing teo_function_v4 resource (name is split from zone suffix).
func TestTeoFunctionV4_Read(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaFunctionV4().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeFunctions", func(request *teov20220901.DescribeFunctionsRequest) (*teov20220901.DescribeFunctionsResponse, error) {
		assert.Equal(t, "zone-test123", *request.ZoneId)
		return mockDescribeFunctionResponse(newFunctionV4(
			"ef-test456", "zone-test123", "my-func-zone-test123",
			"test remark", "addEventListener('fetch', e => {});",
			"my-func.example.com", "2024-01-01T00:00:00Z", "2024-01-01T00:00:00Z",
		))
	})

	meta := newMockMetaFunctionV4()
	res := svcteo.ResourceTencentCloudTeoFunctionV4()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-test123",
		"name":    "placeholder",
		"content": "placeholder",
	})
	d.SetId("zone-test123#ef-test456")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-test123#ef-test456", d.Id())
	assert.Equal(t, "zone-test123", d.Get("zone_id"))
	assert.Equal(t, "ef-test456", d.Get("function_id"))
	assert.Equal(t, "my-func", d.Get("name"))
	assert.Equal(t, "test remark", d.Get("remark"))
	assert.Equal(t, "addEventListener('fetch', e => {});", d.Get("content"))
	assert.Equal(t, "my-func.example.com", d.Get("domain"))
	assert.Equal(t, "2024-01-01T00:00:00Z", d.Get("create_time"))
	assert.Equal(t, "2024-01-01T00:00:00Z", d.Get("update_time"))
}

// TestTeoFunctionV4_ReadNotFound tests that read clears the id when the resource is externally deleted.
func TestTeoFunctionV4_ReadNotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaFunctionV4().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeFunctions", func(request *teov20220901.DescribeFunctionsRequest) (*teov20220901.DescribeFunctionsResponse, error) {
		resp := teov20220901.NewDescribeFunctionsResponse()
		resp.Response = &teov20220901.DescribeFunctionsResponseParams{
			TotalCount: ptrInt64FunctionV4(0),
			Functions:  []*teov20220901.Function{},
			RequestId:  ptrStringFunctionV4("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaFunctionV4()
	res := svcteo.ResourceTencentCloudTeoFunctionV4()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-test123",
		"name":    "my-func",
		"content": "placeholder",
	})
	d.SetId("zone-test123#ef-test456")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}

// TestTeoFunctionV4_ReadBrokenId tests that read returns an error when the composite id is broken.
func TestTeoFunctionV4_ReadBrokenId(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaFunctionV4().client, "UseTeoV20220901Client", teoClient)

	meta := newMockMetaFunctionV4()
	res := svcteo.ResourceTencentCloudTeoFunctionV4()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-test123",
		"name":    "my-func",
		"content": "placeholder",
	})
	d.SetId("broken-id")

	err := res.Read(d, meta)
	assert.Error(t, err)
}

// TestTeoFunctionV4_UpdateImmutable tests that updating an immutable field (name) returns an error.
func TestTeoFunctionV4_UpdateImmutable(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaFunctionV4().client, "UseTeoV20220901Client", teoClient)

	meta := newMockMetaFunctionV4()
	res := svcteo.ResourceTencentCloudTeoFunctionV4()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-test123",
		"name":    "new-func",
		"content": "placeholder",
	})
	d.SetId("zone-test123#ef-test456")

	// Force only the immutable field "name" to be detected as changed.
	patches.ApplyMethodFunc(d, "HasChange", func(key string) bool {
		return key == "name"
	})

	err := res.Update(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "name")
}

// TestTeoFunctionV4_UpdateMutable tests that updating mutable fields (remark/content) calls ModifyFunction.
func TestTeoFunctionV4_UpdateMutable(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaFunctionV4().client, "UseTeoV20220901Client", teoClient)

	modifyCalled := false
	patches.ApplyMethodFunc(teoClient, "ModifyFunctionWithContext", func(_ context.Context, request *teov20220901.ModifyFunctionRequest) (*teov20220901.ModifyFunctionResponse, error) {
		modifyCalled = true
		assert.Equal(t, "zone-test123", *request.ZoneId)
		assert.Equal(t, "ef-test456", *request.FunctionId)
		assert.Equal(t, "updated remark", *request.Remark)
		assert.Equal(t, "addEventListener('fetch', e => { return 1; });", *request.Content)

		resp := teov20220901.NewModifyFunctionResponse()
		resp.Response = &teov20220901.ModifyFunctionResponseParams{
			RequestId: ptrStringFunctionV4("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(teoClient, "DescribeFunctions", func(request *teov20220901.DescribeFunctionsRequest) (*teov20220901.DescribeFunctionsResponse, error) {
		return mockDescribeFunctionResponse(newFunctionV4(
			"ef-test456", "zone-test123", "my-func-zone-test123",
			"updated remark", "addEventListener('fetch', e => { return 1; });",
			"my-func.example.com", "2024-01-01T00:00:00Z", "2024-01-02T00:00:00Z",
		))
	})

	meta := newMockMetaFunctionV4()
	res := svcteo.ResourceTencentCloudTeoFunctionV4()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-test123",
		"name":    "my-func",
		"remark":  "updated remark",
		"content": "addEventListener('fetch', e => { return 1; });",
	})
	d.SetId("zone-test123#ef-test456")

	// Force only mutable fields (remark, content) to be detected as changed; immutable field
	// "name" returns false so the immutable-args check does not error out.
	patches.ApplyMethodFunc(d, "HasChange", func(key string) bool {
		return key == "remark" || key == "content"
	})

	err := res.Update(d, meta)
	assert.NoError(t, err)
	assert.True(t, modifyCalled)
	assert.Equal(t, "updated remark", d.Get("remark"))
	assert.Equal(t, "addEventListener('fetch', e => { return 1; });", d.Get("content"))
}

// TestTeoFunctionV4_UpdateNoChange tests that update skips ModifyFunction when no mutable field changed.
func TestTeoFunctionV4_UpdateNoChange(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaFunctionV4().client, "UseTeoV20220901Client", teoClient)

	modifyCalled := false
	patches.ApplyMethodFunc(teoClient, "ModifyFunctionWithContext", func(_ context.Context, request *teov20220901.ModifyFunctionRequest) (*teov20220901.ModifyFunctionResponse, error) {
		modifyCalled = true
		resp := teov20220901.NewModifyFunctionResponse()
		resp.Response = &teov20220901.ModifyFunctionResponseParams{
			RequestId: ptrStringFunctionV4("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(teoClient, "DescribeFunctions", func(request *teov20220901.DescribeFunctionsRequest) (*teov20220901.DescribeFunctionsResponse, error) {
		return mockDescribeFunctionResponse(newFunctionV4(
			"ef-test456", "zone-test123", "my-func-zone-test123",
			"test remark", "placeholder",
			"my-func.example.com", "2024-01-01T00:00:00Z", "2024-01-01T00:00:00Z",
		))
	})

	meta := newMockMetaFunctionV4()
	res := svcteo.ResourceTencentCloudTeoFunctionV4()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-test123",
		"name":    "my-func",
		"remark":  "test remark",
		"content": "placeholder",
	})
	d.SetId("zone-test123#ef-test456")

	// No field is detected as changed.
	patches.ApplyMethodFunc(d, "HasChange", func(key string) bool {
		return false
	})

	err := res.Update(d, meta)
	assert.NoError(t, err)
	assert.False(t, modifyCalled)
}

// TestTeoFunctionV4_UpdateBrokenId tests that update returns an error when the composite id is broken.
func TestTeoFunctionV4_UpdateBrokenId(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaFunctionV4().client, "UseTeoV20220901Client", teoClient)

	meta := newMockMetaFunctionV4()
	res := svcteo.ResourceTencentCloudTeoFunctionV4()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-test123",
		"name":    "my-func",
		"content": "placeholder",
	})
	d.SetId("broken-id")

	// No changes detected, so we reach the id parsing logic first.
	patches.ApplyMethodFunc(d, "HasChange", func(key string) bool {
		return false
	})

	err := res.Update(d, meta)
	assert.Error(t, err)
}

// TestTeoFunctionV4_Delete tests deleting a teo_function_v4 resource successfully.
func TestTeoFunctionV4_Delete(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaFunctionV4().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DeleteFunctionWithContext", func(_ context.Context, request *teov20220901.DeleteFunctionRequest) (*teov20220901.DeleteFunctionResponse, error) {
		assert.Equal(t, "zone-test123", *request.ZoneId)
		assert.Equal(t, "ef-test456", *request.FunctionId)

		resp := teov20220901.NewDeleteFunctionResponse()
		resp.Response = &teov20220901.DeleteFunctionResponseParams{
			RequestId: ptrStringFunctionV4("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaFunctionV4()
	res := svcteo.ResourceTencentCloudTeoFunctionV4()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-test123",
		"name":    "my-func",
		"content": "placeholder",
	})
	d.SetId("zone-test123#ef-test456")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
}

// TestTeoFunctionV4_DeleteBrokenId tests that delete returns an error when the composite id is broken.
func TestTeoFunctionV4_DeleteBrokenId(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaFunctionV4().client, "UseTeoV20220901Client", teoClient)

	meta := newMockMetaFunctionV4()
	res := svcteo.ResourceTencentCloudTeoFunctionV4()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-test123",
		"name":    "my-func",
		"content": "placeholder",
	})
	d.SetId("broken-id")

	err := res.Delete(d, meta)
	assert.Error(t, err)
}

// TestTeoFunctionV4_ParseOriginalName covers parseTeoFunctionV4OriginalName via the Read path by
// observing how the returned Name is trimmed of the "-<zoneId>" suffix. Each case feeds a Name
// through the mocked DescribeFunctions and checks the resulting "name" state value.
func TestTeoFunctionV4_ParseOriginalName(t *testing.T) {
	tests := []struct {
		name     string
		fullName string
		zoneId   string
		expected string
	}{
		{
			name:     "standard concatenated name",
			fullName: "my-func-zone-test123",
			zoneId:   "zone-test123",
			expected: "my-func",
		},
		{
			name:     "original name without hyphens",
			fullName: "myfunc-zone-test123",
			zoneId:   "zone-test123",
			expected: "myfunc",
		},
		{
			name:     "original name with multiple hyphens",
			fullName: "my-test-func-v2-zone-test123",
			zoneId:   "zone-test123",
			expected: "my-test-func-v2",
		},
		{
			name:     "name without -zone substring returns original",
			fullName: "myfunc",
			zoneId:   "zone-test123",
			expected: "myfunc",
		},
		{
			name:     "empty string returns empty",
			fullName: "",
			zoneId:   "zone-test123",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			patches := gomonkey.NewPatches()
			defer patches.Reset()

			teoClient := &teov20220901.Client{}
			patches.ApplyMethodReturn(newMockMetaFunctionV4().client, "UseTeoV20220901Client", teoClient)

			patches.ApplyMethodFunc(teoClient, "DescribeFunctions", func(request *teov20220901.DescribeFunctionsRequest) (*teov20220901.DescribeFunctionsResponse, error) {
				fn := newFunctionV4("ef-test456", tt.zoneId, tt.fullName, "", "", "", "", "")
				return mockDescribeFunctionResponse(fn)
			})

			meta := newMockMetaFunctionV4()
			res := svcteo.ResourceTencentCloudTeoFunctionV4()
			d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
				"zone_id": tt.zoneId,
				"name":    "placeholder",
				"content": "placeholder",
			})
			d.SetId(tt.zoneId + "#ef-test456")

			err := res.Read(d, meta)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, d.Get("name"))
		})
	}
}
