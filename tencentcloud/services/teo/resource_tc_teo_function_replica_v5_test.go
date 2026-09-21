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

// mockMetaFunctionReplicaV5 implements tccommon.ProviderMeta
type mockMetaFunctionReplicaV5 struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaFunctionReplicaV5) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaFunctionReplicaV5{}

func newMockMetaFunctionReplicaV5() *mockMetaFunctionReplicaV5 {
	return &mockMetaFunctionReplicaV5{client: &connectivity.TencentCloudClient{}}
}

func ptrStringFunctionReplicaV5(s string) *string {
	return &s
}

func ptrInt64FunctionReplicaV5(i int64) *int64 {
	return &i
}

// go test ./tencentcloud/services/teo/ -run "TestTeoFunctionReplicaV5_Create" -v -count=1 -gcflags="all=-l"
func TestTeoFunctionReplicaV5_Create(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaFunctionReplicaV5().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "CreateFunctionReplicaWithContext", func(_ context.Context, request *teov20220901.CreateFunctionReplicaRequest) (*teov20220901.CreateFunctionReplicaResponse, error) {
		assert.Equal(t, "zone-test123", *request.ZoneId)
		assert.Equal(t, "ef-test456", *request.FunctionId)
		assert.Equal(t, "test-replica", *request.ReplicaName)
		assert.Equal(t, "console.log(123)", *request.Content)
		assert.Equal(t, "test remark", *request.Remark)

		resp := teov20220901.NewCreateFunctionReplicaResponse()
		resp.Response = &teov20220901.CreateFunctionReplicaResponseParams{
			RequestId: ptrStringFunctionReplicaV5("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(teoClient, "DescribeFunctionReplicasWithContext", func(_ context.Context, request *teov20220901.DescribeFunctionReplicasRequest) (*teov20220901.DescribeFunctionReplicasResponse, error) {
		assert.Equal(t, "zone-test123", *request.ZoneId)
		assert.Equal(t, "ef-test456", *request.FunctionId)
		assert.Equal(t, int64(200), *request.Limit)

		resp := teov20220901.NewDescribeFunctionReplicasResponse()
		resp.Response = &teov20220901.DescribeFunctionReplicasResponseParams{
			TotalCount: ptrInt64FunctionReplicaV5(1),
			FunctionReplicas: []*teov20220901.FunctionReplica{
				{
					FunctionId:  ptrStringFunctionReplicaV5("ef-test456"),
					ReplicaName: ptrStringFunctionReplicaV5("test-replica"),
					Content:     ptrStringFunctionReplicaV5("console.log(123)"),
					Remark:      ptrStringFunctionReplicaV5("test remark"),
					CreatedOn:   ptrStringFunctionReplicaV5("2024-01-01T00:00:00Z"),
					ModifiedOn:  ptrStringFunctionReplicaV5("2024-01-01T00:00:00Z"),
				},
			},
			RequestId: ptrStringFunctionReplicaV5("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaFunctionReplicaV5()
	res := teo.ResourceTencentCloudTeoFunctionReplicaV5()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":      "zone-test123",
		"function_id":  "ef-test456",
		"replica_name": "test-replica",
		"content":      "console.log(123)",
		"remark":       "test remark",
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-test123#ef-test456#test-replica", d.Id())
	assert.Equal(t, "console.log(123)", d.Get("content"))
	assert.Equal(t, "test remark", d.Get("remark"))
	assert.Equal(t, "2024-01-01T00:00:00Z", d.Get("created_on"))
	assert.Equal(t, "2024-01-01T00:00:00Z", d.Get("modified_on"))
}

// go test ./tencentcloud/services/teo/ -run "TestTeoFunctionReplicaV5_Read" -v -count=1 -gcflags="all=-l"
func TestTeoFunctionReplicaV5_Read(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaFunctionReplicaV5().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeFunctionReplicasWithContext", func(_ context.Context, request *teov20220901.DescribeFunctionReplicasRequest) (*teov20220901.DescribeFunctionReplicasResponse, error) {
		assert.Equal(t, "zone-test123", *request.ZoneId)
		assert.Equal(t, "ef-test456", *request.FunctionId)
		assert.Equal(t, int64(200), *request.Limit)
		assert.Equal(t, "created-on", *request.SortBy)
		assert.Equal(t, "desc", *request.SortOrder)

		// user-configured filter + internal replica-name filter
		assert.Len(t, request.Filters, 2)
		assert.Equal(t, "replica-name", *request.Filters[0].Name)
		assert.Equal(t, 1, len(request.Filters[0].Values))
		assert.Equal(t, "test-replica", *request.Filters[0].Values[0])
		assert.True(t, *request.Filters[0].Fuzzy)
		assert.Equal(t, "replica-name", *request.Filters[1].Name)
		assert.Equal(t, 1, len(request.Filters[1].Values))
		assert.Equal(t, "test-replica", *request.Filters[1].Values[0])

		resp := teov20220901.NewDescribeFunctionReplicasResponse()
		resp.Response = &teov20220901.DescribeFunctionReplicasResponseParams{
			TotalCount: ptrInt64FunctionReplicaV5(1),
			FunctionReplicas: []*teov20220901.FunctionReplica{
				{
					FunctionId:  ptrStringFunctionReplicaV5("ef-test456"),
					ReplicaName: ptrStringFunctionReplicaV5("test-replica"),
					Content:     ptrStringFunctionReplicaV5("console.log(123)"),
					Remark:      ptrStringFunctionReplicaV5("test remark"),
					CreatedOn:   ptrStringFunctionReplicaV5("2024-01-01T00:00:00Z"),
					ModifiedOn:  ptrStringFunctionReplicaV5("2024-01-02T00:00:00Z"),
				},
			},
			RequestId: ptrStringFunctionReplicaV5("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaFunctionReplicaV5()
	res := teo.ResourceTencentCloudTeoFunctionReplicaV5()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":      "zone-test123",
		"function_id":  "ef-test456",
		"replica_name": "test-replica",
		"content":      "console.log(123)",
		"remark":       "test remark",
		"sort_by":      "created-on",
		"sort_order":   "desc",
		"filters": []interface{}{
			map[string]interface{}{
				"name":   "replica-name",
				"values": []interface{}{"test-replica"},
				"fuzzy":  true,
			},
		},
	})
	d.SetId("zone-test123#ef-test456#test-replica")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-test123", d.Get("zone_id"))
	assert.Equal(t, "ef-test456", d.Get("function_id"))
	assert.Equal(t, "test-replica", d.Get("replica_name"))
	assert.Equal(t, "console.log(123)", d.Get("content"))
	assert.Equal(t, "test remark", d.Get("remark"))
	assert.Equal(t, "2024-01-01T00:00:00Z", d.Get("created_on"))
	assert.Equal(t, "2024-01-02T00:00:00Z", d.Get("modified_on"))
}

// go test ./tencentcloud/services/teo/ -run "TestTeoFunctionReplicaV5_ReadFuzzyMatch" -v -count=1 -gcflags="all=-l"
func TestTeoFunctionReplicaV5_ReadFuzzyMatch(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaFunctionReplicaV5().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeFunctionReplicasWithContext", func(_ context.Context, request *teov20220901.DescribeFunctionReplicasRequest) (*teov20220901.DescribeFunctionReplicasResponse, error) {
		// fuzzy filter returns multiple replicas, exact match must pick the right one
		resp := teov20220901.NewDescribeFunctionReplicasResponse()
		resp.Response = &teov20220901.DescribeFunctionReplicasResponseParams{
			TotalCount: ptrInt64FunctionReplicaV5(2),
			FunctionReplicas: []*teov20220901.FunctionReplica{
				{
					FunctionId:  ptrStringFunctionReplicaV5("ef-test456"),
					ReplicaName: ptrStringFunctionReplicaV5("test-replica-2"),
					Content:     ptrStringFunctionReplicaV5("console.log(999)"),
					Remark:      ptrStringFunctionReplicaV5("other replica"),
					CreatedOn:   ptrStringFunctionReplicaV5("2024-01-01T00:00:00Z"),
					ModifiedOn:  ptrStringFunctionReplicaV5("2024-01-01T00:00:00Z"),
				},
				{
					FunctionId:  ptrStringFunctionReplicaV5("ef-test456"),
					ReplicaName: ptrStringFunctionReplicaV5("test-replica"),
					Content:     ptrStringFunctionReplicaV5("console.log(123)"),
					Remark:      ptrStringFunctionReplicaV5("test remark"),
					CreatedOn:   ptrStringFunctionReplicaV5("2024-01-01T00:00:00Z"),
					ModifiedOn:  ptrStringFunctionReplicaV5("2024-01-02T00:00:00Z"),
				},
			},
			RequestId: ptrStringFunctionReplicaV5("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaFunctionReplicaV5()
	res := teo.ResourceTencentCloudTeoFunctionReplicaV5()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":      "zone-test123",
		"function_id":  "ef-test456",
		"replica_name": "test-replica",
		"content":      "console.log(123)",
	})
	d.SetId("zone-test123#ef-test456#test-replica")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "console.log(123)", d.Get("content"))
	assert.Equal(t, "test-replica", d.Get("replica_name"))
}

// go test ./tencentcloud/services/teo/ -run "TestTeoFunctionReplicaV5_ReadNotFound" -v -count=1 -gcflags="all=-l"
func TestTeoFunctionReplicaV5_ReadNotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaFunctionReplicaV5().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeFunctionReplicasWithContext", func(_ context.Context, request *teov20220901.DescribeFunctionReplicasRequest) (*teov20220901.DescribeFunctionReplicasResponse, error) {
		resp := teov20220901.NewDescribeFunctionReplicasResponse()
		resp.Response = &teov20220901.DescribeFunctionReplicasResponseParams{
			TotalCount:       ptrInt64FunctionReplicaV5(0),
			FunctionReplicas: []*teov20220901.FunctionReplica{},
			RequestId:        ptrStringFunctionReplicaV5("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaFunctionReplicaV5()
	res := teo.ResourceTencentCloudTeoFunctionReplicaV5()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":      "zone-test123",
		"function_id":  "ef-test456",
		"replica_name": "test-replica",
		"content":      "console.log(123)",
		"remark":       "test remark",
	})
	d.SetId("zone-test123#ef-test456#test-replica")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}

// go test ./tencentcloud/services/teo/ -run "TestTeoFunctionReplicaV5_Update" -v -count=1 -gcflags="all=-l"
func TestTeoFunctionReplicaV5_Update(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaFunctionReplicaV5().client, "UseTeoV20220901Client", teoClient)

	modifyCalled := false
	patches.ApplyMethodFunc(teoClient, "ModifyFunctionReplicaWithContext", func(_ context.Context, request *teov20220901.ModifyFunctionReplicaRequest) (*teov20220901.ModifyFunctionReplicaResponse, error) {
		modifyCalled = true
		assert.Equal(t, "zone-test123", *request.ZoneId)
		assert.Equal(t, "ef-test456", *request.FunctionId)
		assert.Equal(t, "test-replica", *request.ReplicaName)
		assert.Equal(t, "console.log(456)", *request.Content)
		assert.Equal(t, "updated remark", *request.Remark)

		resp := teov20220901.NewModifyFunctionReplicaResponse()
		resp.Response = &teov20220901.ModifyFunctionReplicaResponseParams{
			RequestId: ptrStringFunctionReplicaV5("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(teoClient, "DescribeFunctionReplicasWithContext", func(_ context.Context, request *teov20220901.DescribeFunctionReplicasRequest) (*teov20220901.DescribeFunctionReplicasResponse, error) {
		resp := teov20220901.NewDescribeFunctionReplicasResponse()
		resp.Response = &teov20220901.DescribeFunctionReplicasResponseParams{
			TotalCount: ptrInt64FunctionReplicaV5(1),
			FunctionReplicas: []*teov20220901.FunctionReplica{
				{
					FunctionId:  ptrStringFunctionReplicaV5("ef-test456"),
					ReplicaName: ptrStringFunctionReplicaV5("test-replica"),
					Content:     ptrStringFunctionReplicaV5("console.log(456)"),
					Remark:      ptrStringFunctionReplicaV5("updated remark"),
					CreatedOn:   ptrStringFunctionReplicaV5("2024-01-01T00:00:00Z"),
					ModifiedOn:  ptrStringFunctionReplicaV5("2024-01-02T00:00:00Z"),
				},
			},
			RequestId: ptrStringFunctionReplicaV5("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaFunctionReplicaV5()
	res := teo.ResourceTencentCloudTeoFunctionReplicaV5()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":      "zone-test123",
		"function_id":  "ef-test456",
		"replica_name": "test-replica",
		"content":      "console.log(456)",
		"remark":       "updated remark",
	})
	d.SetId("zone-test123#ef-test456#test-replica")

	err := res.Update(d, meta)
	assert.NoError(t, err)
	assert.True(t, modifyCalled)
	assert.Equal(t, "console.log(456)", d.Get("content"))
	assert.Equal(t, "updated remark", d.Get("remark"))
}

// go test ./tencentcloud/services/teo/ -run "TestTeoFunctionReplicaV5_UpdateNoChange" -v -count=1 -gcflags="all=-l"
func TestTeoFunctionReplicaV5_UpdateNoChange(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaFunctionReplicaV5().client, "UseTeoV20220901Client", teoClient)

	modifyCalled := false
	patches.ApplyMethodFunc(teoClient, "ModifyFunctionReplicaWithContext", func(_ context.Context, request *teov20220901.ModifyFunctionReplicaRequest) (*teov20220901.ModifyFunctionReplicaResponse, error) {
		modifyCalled = true

		resp := teov20220901.NewModifyFunctionReplicaResponse()
		resp.Response = &teov20220901.ModifyFunctionReplicaResponseParams{
			RequestId: ptrStringFunctionReplicaV5("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(teoClient, "DescribeFunctionReplicasWithContext", func(_ context.Context, request *teov20220901.DescribeFunctionReplicasRequest) (*teov20220901.DescribeFunctionReplicasResponse, error) {
		resp := teov20220901.NewDescribeFunctionReplicasResponse()
		resp.Response = &teov20220901.DescribeFunctionReplicasResponseParams{
			TotalCount: ptrInt64FunctionReplicaV5(1),
			FunctionReplicas: []*teov20220901.FunctionReplica{
				{
					FunctionId:  ptrStringFunctionReplicaV5("ef-test456"),
					ReplicaName: ptrStringFunctionReplicaV5("test-replica"),
					Content:     ptrStringFunctionReplicaV5("console.log(123)"),
					Remark:      ptrStringFunctionReplicaV5("test remark"),
					CreatedOn:   ptrStringFunctionReplicaV5("2024-01-01T00:00:00Z"),
					ModifiedOn:  ptrStringFunctionReplicaV5("2024-01-01T00:00:00Z"),
				},
			},
			RequestId: ptrStringFunctionReplicaV5("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaFunctionReplicaV5()
	res := teo.ResourceTencentCloudTeoFunctionReplicaV5()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":      "zone-test123",
		"function_id":  "ef-test456",
		"replica_name": "test-replica",
		"content":      "console.log(123)",
		"remark":       "test remark",
	})
	d.SetId("zone-test123#ef-test456#test-replica")

	// no changes detected
	patches.ApplyMethodFunc(d, "HasChange", func(key string) bool {
		return false
	})

	err := res.Update(d, meta)
	assert.NoError(t, err)
	assert.False(t, modifyCalled)
	assert.Equal(t, "console.log(123)", d.Get("content"))
}

// go test ./tencentcloud/services/teo/ -run "TestTeoFunctionReplicaV5_Delete" -v -count=1 -gcflags="all=-l"
func TestTeoFunctionReplicaV5_Delete(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaFunctionReplicaV5().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DeleteFunctionReplicaWithContext", func(_ context.Context, request *teov20220901.DeleteFunctionReplicaRequest) (*teov20220901.DeleteFunctionReplicaResponse, error) {
		assert.Equal(t, "zone-test123", *request.ZoneId)
		assert.Equal(t, "ef-test456", *request.FunctionId)
		// replica_names not configured, falls back to single replica name from composite id
		assert.Equal(t, 1, len(request.ReplicaNames))
		assert.Equal(t, "test-replica", *request.ReplicaNames[0])

		resp := teov20220901.NewDeleteFunctionReplicaResponse()
		resp.Response = &teov20220901.DeleteFunctionReplicaResponseParams{
			RequestId: ptrStringFunctionReplicaV5("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaFunctionReplicaV5()
	res := teo.ResourceTencentCloudTeoFunctionReplicaV5()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":      "zone-test123",
		"function_id":  "ef-test456",
		"replica_name": "test-replica",
		"content":      "console.log(123)",
		"remark":       "test remark",
	})
	d.SetId("zone-test123#ef-test456#test-replica")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
}

// go test ./tencentcloud/services/teo/ -run "TestTeoFunctionReplicaV5_DeleteWithReplicaNames" -v -count=1 -gcflags="all=-l"
func TestTeoFunctionReplicaV5_DeleteWithReplicaNames(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaFunctionReplicaV5().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DeleteFunctionReplicaWithContext", func(_ context.Context, request *teov20220901.DeleteFunctionReplicaRequest) (*teov20220901.DeleteFunctionReplicaResponse, error) {
		assert.Equal(t, "zone-test123", *request.ZoneId)
		assert.Equal(t, "ef-test456", *request.FunctionId)
		// replica_names configured, takes precedence over the composite id replica name
		assert.Equal(t, 2, len(request.ReplicaNames))
		assert.Equal(t, "test-replica", *request.ReplicaNames[0])
		assert.Equal(t, "test-replica-2", *request.ReplicaNames[1])

		resp := teov20220901.NewDeleteFunctionReplicaResponse()
		resp.Response = &teov20220901.DeleteFunctionReplicaResponseParams{
			RequestId: ptrStringFunctionReplicaV5("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaFunctionReplicaV5()
	res := teo.ResourceTencentCloudTeoFunctionReplicaV5()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":       "zone-test123",
		"function_id":   "ef-test456",
		"replica_name":  "test-replica",
		"content":       "console.log(123)",
		"remark":        "test remark",
		"replica_names": []interface{}{"test-replica", "test-replica-2"},
	})
	d.SetId("zone-test123#ef-test456#test-replica")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
}

// go test ./tencentcloud/services/teo/ -run "TestTeoFunctionReplicaV5_ReadBrokenId" -v -count=1 -gcflags="all=-l"
func TestTeoFunctionReplicaV5_ReadBrokenId(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaFunctionReplicaV5().client, "UseTeoV20220901Client", teoClient)

	meta := newMockMetaFunctionReplicaV5()
	res := teo.ResourceTencentCloudTeoFunctionReplicaV5()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":      "zone-test123",
		"function_id":  "ef-test456",
		"replica_name": "test-replica",
		"content":      "console.log(123)",
	})
	d.SetId("broken-id")

	err := res.Read(d, meta)
	assert.Error(t, err)
}
