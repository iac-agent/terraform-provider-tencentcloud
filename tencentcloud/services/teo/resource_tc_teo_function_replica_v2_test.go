package teo_test

import (
	"context"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/teo"
)

// mockMetaFunctionReplicaV2 implements tccommon.ProviderMeta
type mockMetaFunctionReplicaV2 struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaFunctionReplicaV2) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaFunctionReplicaV2{}

func newMockMetaFunctionReplicaV2() *mockMetaFunctionReplicaV2 {
	return &mockMetaFunctionReplicaV2{client: &connectivity.TencentCloudClient{}}
}

func ptrStringFunctionReplicaV2(s string) *string {
	return &s
}

func ptrInt64FunctionReplicaV2(i int64) *int64 {
	return &i
}

// go test ./tencentcloud/services/teo/ -run "TestTeoFunctionReplicaV2_Create" -v -count=1 -gcflags="all=-l"
func TestTeoFunctionReplicaV2_Create(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaFunctionReplicaV2().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "CreateFunctionReplicaWithContext", func(_ context.Context, request *teov20220901.CreateFunctionReplicaRequest) (*teov20220901.CreateFunctionReplicaResponse, error) {
		assert.Equal(t, "zone-test123", *request.ZoneId)
		assert.Equal(t, "ef-test456", *request.FunctionId)
		assert.Equal(t, "test-replica", *request.ReplicaName)
		assert.Equal(t, "console.log(123)", *request.Content)
		assert.Equal(t, "test remark", *request.Remark)

		resp := teov20220901.NewCreateFunctionReplicaResponse()
		resp.Response = &teov20220901.CreateFunctionReplicaResponseParams{
			RequestId: ptrStringFunctionReplicaV2("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(teoClient, "DescribeFunctionReplicasWithContext", func(_ context.Context, request *teov20220901.DescribeFunctionReplicasRequest) (*teov20220901.DescribeFunctionReplicasResponse, error) {
		resp := teov20220901.NewDescribeFunctionReplicasResponse()
		resp.Response = &teov20220901.DescribeFunctionReplicasResponseParams{
			TotalCount: ptrInt64FunctionReplicaV2(1),
			FunctionReplicas: []*teov20220901.FunctionReplica{
				{
					FunctionId:  ptrStringFunctionReplicaV2("ef-test456"),
					ReplicaName: ptrStringFunctionReplicaV2("test-replica"),
					Content:     ptrStringFunctionReplicaV2("console.log(123)"),
					Remark:      ptrStringFunctionReplicaV2("test remark"),
					CreatedOn:   ptrStringFunctionReplicaV2("2024-01-01T00:00:00Z"),
					ModifiedOn:  ptrStringFunctionReplicaV2("2024-01-01T00:00:00Z"),
				},
			},
			RequestId: ptrStringFunctionReplicaV2("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaFunctionReplicaV2()
	res := teo.ResourceTencentCloudTeoFunctionReplicaV2()
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

// go test ./tencentcloud/services/teo/ -run "TestTeoFunctionReplicaV2_Read" -v -count=1 -gcflags="all=-l"
func TestTeoFunctionReplicaV2_Read(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaFunctionReplicaV2().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeFunctionReplicasWithContext", func(_ context.Context, request *teov20220901.DescribeFunctionReplicasRequest) (*teov20220901.DescribeFunctionReplicasResponse, error) {
		assert.Equal(t, "zone-test123", *request.ZoneId)
		assert.Equal(t, "ef-test456", *request.FunctionId)
		assert.Equal(t, int64(200), *request.Limit)
		assert.Equal(t, 1, len(request.Filters))
		assert.Equal(t, "replica-name", *request.Filters[0].Name)
		assert.Equal(t, "test-replica", *request.Filters[0].Values[0])

		resp := teov20220901.NewDescribeFunctionReplicasResponse()
		resp.Response = &teov20220901.DescribeFunctionReplicasResponseParams{
			TotalCount: ptrInt64FunctionReplicaV2(1),
			FunctionReplicas: []*teov20220901.FunctionReplica{
				{
					FunctionId:  ptrStringFunctionReplicaV2("ef-test456"),
					ReplicaName: ptrStringFunctionReplicaV2("test-replica"),
					Content:     ptrStringFunctionReplicaV2("console.log(123)"),
					Remark:      ptrStringFunctionReplicaV2("test remark"),
					CreatedOn:   ptrStringFunctionReplicaV2("2024-01-01T00:00:00Z"),
					ModifiedOn:  ptrStringFunctionReplicaV2("2024-01-02T00:00:00Z"),
				},
			},
			RequestId: ptrStringFunctionReplicaV2("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaFunctionReplicaV2()
	res := teo.ResourceTencentCloudTeoFunctionReplicaV2()
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
	assert.Equal(t, "zone-test123#ef-test456#test-replica", d.Id())
	assert.Equal(t, "console.log(123)", d.Get("content"))
	assert.Equal(t, "test remark", d.Get("remark"))
	assert.Equal(t, "2024-01-01T00:00:00Z", d.Get("created_on"))
	assert.Equal(t, "2024-01-02T00:00:00Z", d.Get("modified_on"))
}

// go test ./tencentcloud/services/teo/ -run "TestTeoFunctionReplicaV2_ReadNotFound" -v -count=1 -gcflags="all=-l"
func TestTeoFunctionReplicaV2_ReadNotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaFunctionReplicaV2().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeFunctionReplicasWithContext", func(_ context.Context, request *teov20220901.DescribeFunctionReplicasRequest) (*teov20220901.DescribeFunctionReplicasResponse, error) {
		resp := teov20220901.NewDescribeFunctionReplicasResponse()
		resp.Response = &teov20220901.DescribeFunctionReplicasResponseParams{
			TotalCount:       ptrInt64FunctionReplicaV2(0),
			FunctionReplicas: []*teov20220901.FunctionReplica{},
			RequestId:        ptrStringFunctionReplicaV2("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaFunctionReplicaV2()
	res := teo.ResourceTencentCloudTeoFunctionReplicaV2()
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

// go test ./tencentcloud/services/teo/ -run "TestTeoFunctionReplicaV2_Update" -v -count=1 -gcflags="all=-l"
func TestTeoFunctionReplicaV2_Update(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaFunctionReplicaV2().client, "UseTeoV20220901Client", teoClient)

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
			RequestId: ptrStringFunctionReplicaV2("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(teoClient, "DescribeFunctionReplicasWithContext", func(_ context.Context, request *teov20220901.DescribeFunctionReplicasRequest) (*teov20220901.DescribeFunctionReplicasResponse, error) {
		resp := teov20220901.NewDescribeFunctionReplicasResponse()
		resp.Response = &teov20220901.DescribeFunctionReplicasResponseParams{
			TotalCount: ptrInt64FunctionReplicaV2(1),
			FunctionReplicas: []*teov20220901.FunctionReplica{
				{
					FunctionId:  ptrStringFunctionReplicaV2("ef-test456"),
					ReplicaName: ptrStringFunctionReplicaV2("test-replica"),
					Content:     ptrStringFunctionReplicaV2("console.log(456)"),
					Remark:      ptrStringFunctionReplicaV2("updated remark"),
					CreatedOn:   ptrStringFunctionReplicaV2("2024-01-01T00:00:00Z"),
					ModifiedOn:  ptrStringFunctionReplicaV2("2024-01-02T00:00:00Z"),
				},
			},
			RequestId: ptrStringFunctionReplicaV2("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaFunctionReplicaV2()
	res := teo.ResourceTencentCloudTeoFunctionReplicaV2()

	// Build a prior state where content/remark are the old values, and a new
	// config where content/remark are updated, so that d.HasChange("content")
	// and d.HasChange("remark") return true inside Update.
	state := &terraform.InstanceState{
		ID: "zone-test123#ef-test456#test-replica",
		Attributes: map[string]string{
			"id":           "zone-test123#ef-test456#test-replica",
			"zone_id":      "zone-test123",
			"function_id":  "ef-test456",
			"replica_name": "test-replica",
			"content":      "console.log(123)",
			"remark":       "test remark",
		},
	}

	rawConfig := terraform.NewResourceConfigRaw(map[string]interface{}{
		"zone_id":      "zone-test123",
		"function_id":  "ef-test456",
		"replica_name": "test-replica",
		"content":      "console.log(456)",
		"remark":       "updated remark",
	})

	diff, err := res.Diff(nil, state, rawConfig, meta)
	assert.NoError(t, err)
	assert.NotNil(t, diff)

	d, err := schema.InternalMap(res.Schema).Data(state, diff)
	assert.NoError(t, err)

	err = res.Update(d, meta)
	assert.NoError(t, err)
	assert.True(t, modifyCalled)
	assert.Equal(t, "console.log(456)", d.Get("content"))
	assert.Equal(t, "updated remark", d.Get("remark"))
}

// go test ./tencentcloud/services/teo/ -run "TestTeoFunctionReplicaV2_UpdateImmutableError" -v -count=1 -gcflags="all=-l"
func TestTeoFunctionReplicaV2_UpdateImmutableError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaFunctionReplicaV2().client, "UseTeoV20220901Client", teoClient)

	modifyCalled := false
	patches.ApplyMethodFunc(teoClient, "ModifyFunctionReplicaWithContext", func(_ context.Context, request *teov20220901.ModifyFunctionReplicaRequest) (*teov20220901.ModifyFunctionReplicaResponse, error) {
		modifyCalled = true
		resp := teov20220901.NewModifyFunctionReplicaResponse()
		resp.Response = &teov20220901.ModifyFunctionReplicaResponseParams{
			RequestId: ptrStringFunctionReplicaV2("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaFunctionReplicaV2()
	res := teo.ResourceTencentCloudTeoFunctionReplicaV2()

	// Build a prior state where replica_name is the old value, and a new
	// config where replica_name is changed, so that d.HasChange("replica_name")
	// returns true inside Update and triggers the immutable-argument error.
	state := &terraform.InstanceState{
		ID: "zone-test123#ef-test456#test-replica",
		Attributes: map[string]string{
			"id":           "zone-test123#ef-test456#test-replica",
			"zone_id":      "zone-test123",
			"function_id":  "ef-test456",
			"replica_name": "test-replica",
			"content":      "console.log(123)",
			"remark":       "test remark",
		},
	}

	rawConfig := terraform.NewResourceConfigRaw(map[string]interface{}{
		"zone_id":      "zone-test123",
		"function_id":  "ef-test456",
		"replica_name": "replica-changed",
		"content":      "console.log(123)",
		"remark":       "test remark",
	})

	diff, err := res.Diff(nil, state, rawConfig, meta)
	assert.NoError(t, err)
	assert.NotNil(t, diff)

	d, err := schema.InternalMap(res.Schema).Data(state, diff)
	assert.NoError(t, err)

	err = res.Update(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "immutable")
	assert.False(t, modifyCalled)
}

// go test ./tencentcloud/services/teo/ -run "TestTeoFunctionReplicaV2_Delete" -v -count=1 -gcflags="all=-l"
func TestTeoFunctionReplicaV2_Delete(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaFunctionReplicaV2().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DeleteFunctionReplicaWithContext", func(_ context.Context, request *teov20220901.DeleteFunctionReplicaRequest) (*teov20220901.DeleteFunctionReplicaResponse, error) {
		assert.Equal(t, "zone-test123", *request.ZoneId)
		assert.Equal(t, "ef-test456", *request.FunctionId)
		assert.Equal(t, 1, len(request.ReplicaNames))
		assert.Equal(t, "test-replica", *request.ReplicaNames[0])

		resp := teov20220901.NewDeleteFunctionReplicaResponse()
		resp.Response = &teov20220901.DeleteFunctionReplicaResponseParams{
			RequestId: ptrStringFunctionReplicaV2("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaFunctionReplicaV2()
	res := teo.ResourceTencentCloudTeoFunctionReplicaV2()
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
