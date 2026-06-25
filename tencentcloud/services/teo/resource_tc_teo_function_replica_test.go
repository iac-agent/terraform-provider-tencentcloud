package teo_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/teo"
)

// go test ./tencentcloud/services/teo/ -run "TestFunctionReplica" -v -count=1 -gcflags="all=-l"

// TestFunctionReplica_Create_Success tests Create calls API and sets ID
func TestFunctionReplica_Create_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "CreateFunctionReplicaWithContext", func(ctx context.Context, request *teov20220901.CreateFunctionReplicaRequest) (*teov20220901.CreateFunctionReplicaResponse, error) {
		assert.Equal(t, "zone-12345678", *request.ZoneId)
		assert.Equal(t, "func-abcdefgh", *request.FunctionId)
		assert.Equal(t, "my-replica", *request.ReplicaName)
		resp := teov20220901.NewCreateFunctionReplicaResponse()
		resp.Response = &teov20220901.CreateFunctionReplicaResponseParams{
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(teoClient, "DescribeFunctionReplicasWithContext", func(ctx context.Context, request *teov20220901.DescribeFunctionReplicasRequest) (*teov20220901.DescribeFunctionReplicasResponse, error) {
		resp := teov20220901.NewDescribeFunctionReplicasResponse()
		resp.Response = &teov20220901.DescribeFunctionReplicasResponseParams{
			TotalCount: ptrInt64(1),
			FunctionReplicas: []*teov20220901.FunctionReplica{
				{
					FunctionId:  ptrString("func-abcdefgh"),
					ReplicaName: ptrString("my-replica"),
					Content:     ptrString("function myHandler() {}"),
					Remark:      ptrString("test replica"),
					CreatedOn:   ptrString("2025-01-01T00:00:00Z"),
					ModifiedOn:  ptrString("2025-01-02T00:00:00Z"),
				},
			},
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoFunctionReplica()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":      "zone-12345678",
		"function_id":  "func-abcdefgh",
		"replica_name": "my-replica",
		"content":      "function myHandler() {}",
		"remark":       "test replica",
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-12345678#func-abcdefgh#my-replica", d.Id())
}

// TestFunctionReplica_Create_APIError tests Create handles API error
func TestFunctionReplica_Create_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "CreateFunctionReplicaWithContext", func(ctx context.Context, request *teov20220901.CreateFunctionReplicaRequest) (*teov20220901.CreateFunctionReplicaResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=InvalidParameter, Message=Invalid zone_id")
	})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoFunctionReplica()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":      "zone-invalid",
		"function_id":  "func-abcdefgh",
		"replica_name": "my-replica",
		"content":      "function myHandler() {}",
	})

	err := res.Create(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "InvalidParameter")
}

// TestFunctionReplica_Read_Success tests Read retrieves replica data
func TestFunctionReplica_Read_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeFunctionReplicasWithContext", func(ctx context.Context, request *teov20220901.DescribeFunctionReplicasRequest) (*teov20220901.DescribeFunctionReplicasResponse, error) {
		resp := teov20220901.NewDescribeFunctionReplicasResponse()
		resp.Response = &teov20220901.DescribeFunctionReplicasResponseParams{
			TotalCount: ptrInt64(1),
			FunctionReplicas: []*teov20220901.FunctionReplica{
				{
					FunctionId:  ptrString("func-abcdefgh"),
					ReplicaName: ptrString("my-replica"),
					Content:     ptrString("function myHandler() {}"),
					Remark:      ptrString("test remark"),
					CreatedOn:   ptrString("2025-01-01T00:00:00Z"),
					ModifiedOn:  ptrString("2025-01-02T00:00:00Z"),
				},
			},
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoFunctionReplica()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":      "zone-12345678",
		"function_id":  "func-abcdefgh",
		"replica_name": "my-replica",
	})
	d.SetId("zone-12345678#func-abcdefgh#my-replica")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "my-replica", d.Get("replica_name"))
	assert.Equal(t, "function myHandler() {}", d.Get("content"))
	assert.Equal(t, "test remark", d.Get("remark"))
	assert.Equal(t, "2025-01-01T00:00:00Z", d.Get("created_on"))
	assert.Equal(t, "2025-01-02T00:00:00Z", d.Get("modified_on"))
}

// TestFunctionReplica_Read_NotFound tests Read handles replica not found (empty list)
func TestFunctionReplica_Read_NotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeFunctionReplicasWithContext", func(ctx context.Context, request *teov20220901.DescribeFunctionReplicasRequest) (*teov20220901.DescribeFunctionReplicasResponse, error) {
		resp := teov20220901.NewDescribeFunctionReplicasResponse()
		resp.Response = &teov20220901.DescribeFunctionReplicasResponseParams{
			TotalCount:       ptrInt64(0),
			FunctionReplicas: []*teov20220901.FunctionReplica{},
			RequestId:        ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoFunctionReplica()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":      "zone-12345678",
		"function_id":  "func-abcdefgh",
		"replica_name": "my-replica",
	})
	d.SetId("zone-12345678#func-abcdefgh#my-replica")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}

// TestFunctionReplica_Read_ExactMatchFailed tests Read handles fuzzy match but no exact match
func TestFunctionReplica_Read_ExactMatchFailed(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeFunctionReplicasWithContext", func(ctx context.Context, request *teov20220901.DescribeFunctionReplicasRequest) (*teov20220901.DescribeFunctionReplicasResponse, error) {
		resp := teov20220901.NewDescribeFunctionReplicasResponse()
		resp.Response = &teov20220901.DescribeFunctionReplicasResponseParams{
			TotalCount: ptrInt64(2),
			FunctionReplicas: []*teov20220901.FunctionReplica{
				{
					FunctionId:  ptrString("func-abcdefgh"),
					ReplicaName: ptrString("my-replica-v1"),
					Content:     ptrString("function v1() {}"),
				},
				{
					FunctionId:  ptrString("func-abcdefgh"),
					ReplicaName: ptrString("my-replica-v2"),
					Content:     ptrString("function v2() {}"),
				},
			},
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoFunctionReplica()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":      "zone-12345678",
		"function_id":  "func-abcdefgh",
		"replica_name": "my-replica",
	})
	d.SetId("zone-12345678#func-abcdefgh#my-replica")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}

// TestFunctionReplica_Update_Success tests Update modifies content/remark
func TestFunctionReplica_Update_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)

	modifyCalled := false
	patches.ApplyMethodFunc(teoClient, "ModifyFunctionReplicaWithContext", func(ctx context.Context, request *teov20220901.ModifyFunctionReplicaRequest) (*teov20220901.ModifyFunctionReplicaResponse, error) {
		modifyCalled = true
		assert.Equal(t, "zone-12345678", *request.ZoneId)
		assert.Equal(t, "func-abcdefgh", *request.FunctionId)
		assert.Equal(t, "my-replica", *request.ReplicaName)
		resp := teov20220901.NewModifyFunctionReplicaResponse()
		resp.Response = &teov20220901.ModifyFunctionReplicaResponseParams{
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(teoClient, "DescribeFunctionReplicasWithContext", func(ctx context.Context, request *teov20220901.DescribeFunctionReplicasRequest) (*teov20220901.DescribeFunctionReplicasResponse, error) {
		resp := teov20220901.NewDescribeFunctionReplicasResponse()
		resp.Response = &teov20220901.DescribeFunctionReplicasResponseParams{
			TotalCount: ptrInt64(1),
			FunctionReplicas: []*teov20220901.FunctionReplica{
				{
					FunctionId:  ptrString("func-abcdefgh"),
					ReplicaName: ptrString("my-replica"),
					Content:     ptrString("function updatedHandler() {}"),
					Remark:      ptrString("updated remark"),
					CreatedOn:   ptrString("2025-01-01T00:00:00Z"),
					ModifiedOn:  ptrString("2025-01-03T00:00:00Z"),
				},
			},
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoFunctionReplica()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":      "zone-12345678",
		"function_id":  "func-abcdefgh",
		"replica_name": "my-replica",
		"content":      "function updatedHandler() {}",
		"remark":       "updated remark",
	})
	d.SetId("zone-12345678#func-abcdefgh#my-replica")

	// Simulate a change on content to trigger update
	_ = d.Set("content", "function updatedHandler() {}")

	err := res.Update(d, meta)
	assert.NoError(t, err)
	assert.True(t, modifyCalled)
}

// TestFunctionReplica_Update_NoChange tests Update succeeds with no mutable fields changing
func TestFunctionReplica_Update_NoChange(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "ModifyFunctionReplicaWithContext", func(ctx context.Context, request *teov20220901.ModifyFunctionReplicaRequest) (*teov20220901.ModifyFunctionReplicaResponse, error) {
		resp := teov20220901.NewModifyFunctionReplicaResponse()
		resp.Response = &teov20220901.ModifyFunctionReplicaResponseParams{
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(teoClient, "DescribeFunctionReplicasWithContext", func(ctx context.Context, request *teov20220901.DescribeFunctionReplicasRequest) (*teov20220901.DescribeFunctionReplicasResponse, error) {
		resp := teov20220901.NewDescribeFunctionReplicasResponse()
		resp.Response = &teov20220901.DescribeFunctionReplicasResponseParams{
			TotalCount: ptrInt64(1),
			FunctionReplicas: []*teov20220901.FunctionReplica{
				{
					FunctionId:  ptrString("func-abcdefgh"),
					ReplicaName: ptrString("my-replica"),
					Content:     ptrString("function myHandler() {}"),
					Remark:      ptrString("test remark"),
					CreatedOn:   ptrString("2025-01-01T00:00:00Z"),
					ModifiedOn:  ptrString("2025-01-02T00:00:00Z"),
				},
			},
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoFunctionReplica()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":      "zone-12345678",
		"function_id":  "func-abcdefgh",
		"replica_name": "my-replica",
		"content":      "function myHandler() {}",
		"remark":       "test remark",
	})
	d.SetId("zone-12345678#func-abcdefgh#my-replica")

	err := res.Update(d, meta)
	assert.NoError(t, err)
}

// TestFunctionReplica_Delete_Success tests Delete removes replica
func TestFunctionReplica_Delete_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DeleteFunctionReplicaWithContext", func(ctx context.Context, request *teov20220901.DeleteFunctionReplicaRequest) (*teov20220901.DeleteFunctionReplicaResponse, error) {
		assert.Equal(t, "zone-12345678", *request.ZoneId)
		assert.Equal(t, "func-abcdefgh", *request.FunctionId)
		assert.Equal(t, 1, len(request.ReplicaNames))
		assert.Equal(t, "my-replica", *request.ReplicaNames[0])
		resp := teov20220901.NewDeleteFunctionReplicaResponse()
		resp.Response = &teov20220901.DeleteFunctionReplicaResponseParams{
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoFunctionReplica()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":      "zone-12345678",
		"function_id":  "func-abcdefgh",
		"replica_name": "my-replica",
		"content":      "function myHandler() {}",
	})
	d.SetId("zone-12345678#func-abcdefgh#my-replica")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
}

// TestFunctionReplica_Delete_APIError tests Delete handles API error
func TestFunctionReplica_Delete_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DeleteFunctionReplicaWithContext", func(ctx context.Context, request *teov20220901.DeleteFunctionReplicaRequest) (*teov20220901.DeleteFunctionReplicaResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=ResourceNotFound, Message=Replica not found")
	})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoFunctionReplica()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":      "zone-12345678",
		"function_id":  "func-abcdefgh",
		"replica_name": "my-replica",
		"content":      "function myHandler() {}",
	})
	d.SetId("zone-12345678#func-abcdefgh#my-replica")

	err := res.Delete(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ResourceNotFound")
}

// TestFunctionReplica_Schema validates schema definition
func TestFunctionReplica_Schema(t *testing.T) {
	res := teo.ResourceTencentCloudTeoFunctionReplica()

	assert.NotNil(t, res)
	assert.NotNil(t, res.Create)
	assert.NotNil(t, res.Read)
	assert.NotNil(t, res.Update)
	assert.NotNil(t, res.Delete)
	assert.NotNil(t, res.Importer)

	// Check required + ForceNew fields
	assert.Contains(t, res.Schema, "zone_id")
	zoneId := res.Schema["zone_id"]
	assert.Equal(t, schema.TypeString, zoneId.Type)
	assert.True(t, zoneId.Required)
	assert.True(t, zoneId.ForceNew)

	assert.Contains(t, res.Schema, "function_id")
	functionId := res.Schema["function_id"]
	assert.Equal(t, schema.TypeString, functionId.Type)
	assert.True(t, functionId.Required)
	assert.True(t, functionId.ForceNew)

	assert.Contains(t, res.Schema, "replica_name")
	replicaName := res.Schema["replica_name"]
	assert.Equal(t, schema.TypeString, replicaName.Type)
	assert.True(t, replicaName.Required)
	assert.True(t, replicaName.ForceNew)

	// Check required WITHOUT ForceNew
	assert.Contains(t, res.Schema, "content")
	content := res.Schema["content"]
	assert.Equal(t, schema.TypeString, content.Type)
	assert.True(t, content.Required)
	assert.False(t, content.ForceNew)

	// Check optional fields
	assert.Contains(t, res.Schema, "remark")
	remark := res.Schema["remark"]
	assert.Equal(t, schema.TypeString, remark.Type)
	assert.True(t, remark.Optional)
	assert.False(t, remark.ForceNew)

	// Check computed fields
	assert.Contains(t, res.Schema, "created_on")
	createdOn := res.Schema["created_on"]
	assert.Equal(t, schema.TypeString, createdOn.Type)
	assert.True(t, createdOn.Computed)

	assert.Contains(t, res.Schema, "modified_on")
	modifiedOn := res.Schema["modified_on"]
	assert.Equal(t, schema.TypeString, modifiedOn.Type)
	assert.True(t, modifiedOn.Computed)
}
