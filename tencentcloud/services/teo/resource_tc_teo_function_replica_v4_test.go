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

// mockMetaFunctionReplicaV4 implements tccommon.ProviderMeta
type mockMetaFunctionReplicaV4 struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaFunctionReplicaV4) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaFunctionReplicaV4{}

func newMockMetaFunctionReplicaV4() *mockMetaFunctionReplicaV4 {
	return &mockMetaFunctionReplicaV4{client: &connectivity.TencentCloudClient{}}
}

func ptrStringFunctionReplicaV4(s string) *string {
	return &s
}

func ptrInt64FunctionReplicaV4(i int64) *int64 {
	return &i
}

// go test ./tencentcloud/services/teo/ -run "TestTeoFunctionReplicaV4_Create" -v -count=1 -gcflags="all=-l"
func TestTeoFunctionReplicaV4_Create(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaFunctionReplicaV4().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "CreateFunctionReplicaWithContext", func(_ context.Context, request *teov20220901.CreateFunctionReplicaRequest) (*teov20220901.CreateFunctionReplicaResponse, error) {
		assert.Equal(t, "zone-test123", *request.ZoneId)
		assert.Equal(t, "ef-test456", *request.FunctionId)
		assert.Equal(t, "test-replica", *request.ReplicaName)
		assert.Equal(t, "console.log(123)", *request.Content)
		assert.Equal(t, "test remark", *request.Remark)

		resp := teov20220901.NewCreateFunctionReplicaResponse()
		resp.Response = &teov20220901.CreateFunctionReplicaResponseParams{
			RequestId: ptrStringFunctionReplicaV4("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(teoClient, "DescribeFunctionReplicasWithContext", func(_ context.Context, request *teov20220901.DescribeFunctionReplicasRequest) (*teov20220901.DescribeFunctionReplicasResponse, error) {
		assert.Equal(t, int64(200), *request.Limit)
		assert.Equal(t, "replica-name", *request.Filters[0].Name)
		assert.Equal(t, "test-replica", *request.Filters[0].Values[0])

		resp := teov20220901.NewDescribeFunctionReplicasResponse()
		resp.Response = &teov20220901.DescribeFunctionReplicasResponseParams{
			TotalCount: ptrInt64FunctionReplicaV4(1),
			FunctionReplicas: []*teov20220901.FunctionReplica{
				{
					FunctionId:  ptrStringFunctionReplicaV4("ef-test456"),
					ReplicaName: ptrStringFunctionReplicaV4("test-replica"),
					Content:     ptrStringFunctionReplicaV4("console.log(123)"),
					Remark:      ptrStringFunctionReplicaV4("test remark"),
					CreatedOn:   ptrStringFunctionReplicaV4("2024-01-01T00:00:00Z"),
					ModifiedOn:  ptrStringFunctionReplicaV4("2024-01-01T00:00:00Z"),
				},
			},
			RequestId: ptrStringFunctionReplicaV4("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaFunctionReplicaV4()
	res := teo.ResourceTencentCloudTeoFunctionReplicaV4()
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
	assert.Equal(t, "2024-01-01T00:00:00Z", d.Get("created_on"))
	assert.Equal(t, "2024-01-01T00:00:00Z", d.Get("modified_on"))
}

// go test ./tencentcloud/services/teo/ -run "TestTeoFunctionReplicaV4_CreateNoRemark" -v -count=1 -gcflags="all=-l"
func TestTeoFunctionReplicaV4_CreateNoRemark(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaFunctionReplicaV4().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "CreateFunctionReplicaWithContext", func(_ context.Context, request *teov20220901.CreateFunctionReplicaRequest) (*teov20220901.CreateFunctionReplicaResponse, error) {
		assert.Equal(t, "zone-test123", *request.ZoneId)
		assert.Equal(t, "ef-test456", *request.FunctionId)
		assert.Equal(t, "test-replica", *request.ReplicaName)
		assert.Equal(t, "console.log(123)", *request.Content)
		assert.Nil(t, request.Remark)

		resp := teov20220901.NewCreateFunctionReplicaResponse()
		resp.Response = &teov20220901.CreateFunctionReplicaResponseParams{
			RequestId: ptrStringFunctionReplicaV4("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(teoClient, "DescribeFunctionReplicasWithContext", func(_ context.Context, request *teov20220901.DescribeFunctionReplicasRequest) (*teov20220901.DescribeFunctionReplicasResponse, error) {
		resp := teov20220901.NewDescribeFunctionReplicasResponse()
		resp.Response = &teov20220901.DescribeFunctionReplicasResponseParams{
			TotalCount: ptrInt64FunctionReplicaV4(1),
			FunctionReplicas: []*teov20220901.FunctionReplica{
				{
					FunctionId:  ptrStringFunctionReplicaV4("ef-test456"),
					ReplicaName: ptrStringFunctionReplicaV4("test-replica"),
					Content:     ptrStringFunctionReplicaV4("console.log(123)"),
					CreatedOn:   ptrStringFunctionReplicaV4("2024-01-01T00:00:00Z"),
					ModifiedOn:  ptrStringFunctionReplicaV4("2024-01-01T00:00:00Z"),
				},
			},
			RequestId: ptrStringFunctionReplicaV4("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaFunctionReplicaV4()
	res := teo.ResourceTencentCloudTeoFunctionReplicaV4()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":      "zone-test123",
		"function_id":  "ef-test456",
		"replica_name": "test-replica",
		"content":      "console.log(123)",
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-test123#ef-test456#test-replica", d.Id())
}

// go test ./tencentcloud/services/teo/ -run "TestTeoFunctionReplicaV4_Read" -v -count=1 -gcflags="all=-l"
func TestTeoFunctionReplicaV4_Read(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaFunctionReplicaV4().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeFunctionReplicasWithContext", func(_ context.Context, request *teov20220901.DescribeFunctionReplicasRequest) (*teov20220901.DescribeFunctionReplicasResponse, error) {
		assert.Equal(t, "zone-test123", *request.ZoneId)
		assert.Equal(t, "ef-test456", *request.FunctionId)
		assert.Equal(t, int64(200), *request.Limit)
		assert.Equal(t, "replica-name", *request.Filters[0].Name)
		assert.Equal(t, "test-replica", *request.Filters[0].Values[0])

		resp := teov20220901.NewDescribeFunctionReplicasResponse()
		resp.Response = &teov20220901.DescribeFunctionReplicasResponseParams{
			TotalCount: ptrInt64FunctionReplicaV4(1),
			FunctionReplicas: []*teov20220901.FunctionReplica{
				{
					FunctionId:  ptrStringFunctionReplicaV4("ef-test456"),
					ReplicaName: ptrStringFunctionReplicaV4("test-replica"),
					Content:     ptrStringFunctionReplicaV4("console.log(123)"),
					Remark:      ptrStringFunctionReplicaV4("test remark"),
					CreatedOn:   ptrStringFunctionReplicaV4("2024-01-01T00:00:00Z"),
					ModifiedOn:  ptrStringFunctionReplicaV4("2024-01-02T00:00:00Z"),
				},
			},
			RequestId: ptrStringFunctionReplicaV4("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaFunctionReplicaV4()
	res := teo.ResourceTencentCloudTeoFunctionReplicaV4()
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
	assert.Equal(t, "zone-test123", d.Get("zone_id"))
	assert.Equal(t, "ef-test456", d.Get("function_id"))
	assert.Equal(t, "test-replica", d.Get("replica_name"))
	assert.Equal(t, "console.log(123)", d.Get("content"))
	assert.Equal(t, "test remark", d.Get("remark"))
	assert.Equal(t, "2024-01-01T00:00:00Z", d.Get("created_on"))
	assert.Equal(t, "2024-01-02T00:00:00Z", d.Get("modified_on"))
}

// go test ./tencentcloud/services/teo/ -run "TestTeoFunctionReplicaV4_ReadNotFound" -v -count=1 -gcflags="all=-l"
func TestTeoFunctionReplicaV4_ReadNotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaFunctionReplicaV4().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeFunctionReplicasWithContext", func(_ context.Context, request *teov20220901.DescribeFunctionReplicasRequest) (*teov20220901.DescribeFunctionReplicasResponse, error) {
		resp := teov20220901.NewDescribeFunctionReplicasResponse()
		resp.Response = &teov20220901.DescribeFunctionReplicasResponseParams{
			TotalCount:       ptrInt64FunctionReplicaV4(0),
			FunctionReplicas: []*teov20220901.FunctionReplica{},
			RequestId:        ptrStringFunctionReplicaV4("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaFunctionReplicaV4()
	res := teo.ResourceTencentCloudTeoFunctionReplicaV4()
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

// go test ./tencentcloud/services/teo/ -run "TestTeoFunctionReplicaV4_Update" -v -count=1 -gcflags="all=-l"
func TestTeoFunctionReplicaV4_Update(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaFunctionReplicaV4().client, "UseTeoV20220901Client", teoClient)

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
			RequestId: ptrStringFunctionReplicaV4("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(teoClient, "DescribeFunctionReplicasWithContext", func(_ context.Context, request *teov20220901.DescribeFunctionReplicasRequest) (*teov20220901.DescribeFunctionReplicasResponse, error) {
		resp := teov20220901.NewDescribeFunctionReplicasResponse()
		resp.Response = &teov20220901.DescribeFunctionReplicasResponseParams{
			TotalCount: ptrInt64FunctionReplicaV4(1),
			FunctionReplicas: []*teov20220901.FunctionReplica{
				{
					FunctionId:  ptrStringFunctionReplicaV4("ef-test456"),
					ReplicaName: ptrStringFunctionReplicaV4("test-replica"),
					Content:     ptrStringFunctionReplicaV4("console.log(456)"),
					Remark:      ptrStringFunctionReplicaV4("updated remark"),
					CreatedOn:   ptrStringFunctionReplicaV4("2024-01-01T00:00:00Z"),
					ModifiedOn:  ptrStringFunctionReplicaV4("2024-01-02T00:00:00Z"),
				},
			},
			RequestId: ptrStringFunctionReplicaV4("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaFunctionReplicaV4()
	res := teo.ResourceTencentCloudTeoFunctionReplicaV4()
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
	assert.Equal(t, "2024-01-02T00:00:00Z", d.Get("modified_on"))
}

// go test ./tencentcloud/services/teo/ -run "TestTeoFunctionReplicaV4_UpdateNoChange" -v -count=1 -gcflags="all=-l"
func TestTeoFunctionReplicaV4_UpdateNoChange(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaFunctionReplicaV4().client, "UseTeoV20220901Client", teoClient)

	modifyCalled := false
	patches.ApplyMethodFunc(teoClient, "ModifyFunctionReplicaWithContext", func(_ context.Context, request *teov20220901.ModifyFunctionReplicaRequest) (*teov20220901.ModifyFunctionReplicaResponse, error) {
		modifyCalled = true

		resp := teov20220901.NewModifyFunctionReplicaResponse()
		resp.Response = &teov20220901.ModifyFunctionReplicaResponseParams{
			RequestId: ptrStringFunctionReplicaV4("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(teoClient, "DescribeFunctionReplicasWithContext", func(_ context.Context, request *teov20220901.DescribeFunctionReplicasRequest) (*teov20220901.DescribeFunctionReplicasResponse, error) {
		resp := teov20220901.NewDescribeFunctionReplicasResponse()
		resp.Response = &teov20220901.DescribeFunctionReplicasResponseParams{
			TotalCount: ptrInt64FunctionReplicaV4(1),
			FunctionReplicas: []*teov20220901.FunctionReplica{
				{
					FunctionId:  ptrStringFunctionReplicaV4("ef-test456"),
					ReplicaName: ptrStringFunctionReplicaV4("test-replica"),
					Content:     ptrStringFunctionReplicaV4("console.log(123)"),
					Remark:      ptrStringFunctionReplicaV4("test remark"),
					CreatedOn:   ptrStringFunctionReplicaV4("2024-01-01T00:00:00Z"),
					ModifiedOn:  ptrStringFunctionReplicaV4("2024-01-01T00:00:00Z"),
				},
			},
			RequestId: ptrStringFunctionReplicaV4("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaFunctionReplicaV4()
	res := teo.ResourceTencentCloudTeoFunctionReplicaV4()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":      "zone-test123",
		"function_id":  "ef-test456",
		"replica_name": "test-replica",
		"content":      "console.log(123)",
		"remark":       "test remark",
	})
	d.SetId("zone-test123#ef-test456#test-replica")

	err := res.Update(d, meta)
	assert.NoError(t, err)
	assert.False(t, modifyCalled)
}

// go test ./tencentcloud/services/teo/ -run "TestTeoFunctionReplicaV4_Delete" -v -count=1 -gcflags="all=-l"
func TestTeoFunctionReplicaV4_Delete(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaFunctionReplicaV4().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DeleteFunctionReplicaWithContext", func(_ context.Context, request *teov20220901.DeleteFunctionReplicaRequest) (*teov20220901.DeleteFunctionReplicaResponse, error) {
		assert.Equal(t, "zone-test123", *request.ZoneId)
		assert.Equal(t, "ef-test456", *request.FunctionId)
		assert.Equal(t, 1, len(request.ReplicaNames))
		assert.Equal(t, "test-replica", *request.ReplicaNames[0])

		resp := teov20220901.NewDeleteFunctionReplicaResponse()
		resp.Response = &teov20220901.DeleteFunctionReplicaResponseParams{
			RequestId: ptrStringFunctionReplicaV4("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaFunctionReplicaV4()
	res := teo.ResourceTencentCloudTeoFunctionReplicaV4()
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
