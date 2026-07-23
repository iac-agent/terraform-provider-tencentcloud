package teo_test

import (
	"fmt"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/teo"
)

// go test ./tencentcloud/services/teo/ -run "TestTeoFunctionReplicaV1" -v -count=1 -gcflags="all=-l"

// TestTeoFunctionReplicaV1_Schema validates schema definition
func TestTeoFunctionReplicaV1_Schema(t *testing.T) {
	res := teo.ResourceTencentCloudTeoFunctionReplicaV1()

	assert.NotNil(t, res)
	assert.NotNil(t, res.Create)
	assert.NotNil(t, res.Read)
	assert.NotNil(t, res.Update)
	assert.NotNil(t, res.Delete)
	assert.Nil(t, res.Importer)

	// Check required fields
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

	assert.Contains(t, res.Schema, "content")
	content := res.Schema["content"]
	assert.Equal(t, schema.TypeString, content.Type)
	assert.True(t, content.Required)

	assert.Contains(t, res.Schema, "remark")
	remark := res.Schema["remark"]
	assert.Equal(t, schema.TypeString, remark.Type)
	assert.True(t, remark.Required)
	assert.False(t, remark.Optional)

	// Check computed fields
	assert.Contains(t, res.Schema, "create_time")
	createTime := res.Schema["create_time"]
	assert.Equal(t, schema.TypeString, createTime.Type)
	assert.True(t, createTime.Computed)

	assert.Contains(t, res.Schema, "update_time")
	updateTime := res.Schema["update_time"]
	assert.Equal(t, schema.TypeString, updateTime.Type)
	assert.True(t, updateTime.Computed)
}

// TestTeoFunctionReplicaV1_Create_Success tests Create calls API and sets composite ID
func TestTeoFunctionReplicaV1_Create_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "CreateFunctionReplicaWithContext", func(ctx interface{}, request *teov20220901.CreateFunctionReplicaRequest) (*teov20220901.CreateFunctionReplicaResponse, error) {
		resp := teov20220901.NewCreateFunctionReplicaResponse()
		resp.Response = &teov20220901.CreateFunctionReplicaResponseParams{
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(teoClient, "DescribeFunctionReplicas", func(request *teov20220901.DescribeFunctionReplicasRequest) (*teov20220901.DescribeFunctionReplicasResponse, error) {
		resp := teov20220901.NewDescribeFunctionReplicasResponse()
		resp.Response = &teov20220901.DescribeFunctionReplicasResponseParams{
			TotalCount: ptrInt64(1),
			FunctionReplicas: []*teov20220901.FunctionReplica{
				{
					FunctionId:  ptrString("func-abcdefghij"),
					ReplicaName: ptrString("replica-test"),
					Content:     ptrString("addEventListener('fetch', e => { e.respondWith(new Response('Hello World')); });"),
					Remark:      ptrString("test replica"),
					CreatedOn:   ptrString("2024-01-01T00:00:00Z"),
					ModifiedOn:  ptrString("2024-01-01T00:00:00Z"),
				},
			},
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoFunctionReplicaV1()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":      "zone-1234567890",
		"function_id":  "func-abcdefghij",
		"replica_name": "replica-test",
		"content":      "addEventListener('fetch', e => { e.respondWith(new Response('Hello World')); });",
		"remark":       "test replica",
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-1234567890#func-abcdefghij#replica-test", d.Id())
}

// TestTeoFunctionReplicaV1_Create_APIError tests Create handles API error
func TestTeoFunctionReplicaV1_Create_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "CreateFunctionReplicaWithContext", func(ctx interface{}, request *teov20220901.CreateFunctionReplicaRequest) (*teov20220901.CreateFunctionReplicaResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=InvalidParameter, Message=Invalid zone_id")
	})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoFunctionReplicaV1()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":      "zone-invalid",
		"function_id":  "func-abcdefghij",
		"replica_name": "replica-test",
		"content":      "code",
		"remark":       "test",
	})

	err := res.Create(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "InvalidParameter")
}

// TestTeoFunctionReplicaV1_Read_Success tests Read retrieves replica data
func TestTeoFunctionReplicaV1_Read_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeFunctionReplicas", func(request *teov20220901.DescribeFunctionReplicasRequest) (*teov20220901.DescribeFunctionReplicasResponse, error) {
		resp := teov20220901.NewDescribeFunctionReplicasResponse()
		resp.Response = &teov20220901.DescribeFunctionReplicasResponseParams{
			TotalCount: ptrInt64(1),
			FunctionReplicas: []*teov20220901.FunctionReplica{
				{
					FunctionId:  ptrString("func-abcdefghij"),
					ReplicaName: ptrString("replica-test"),
					Content:     ptrString("addEventListener('fetch', e => { e.respondWith(new Response('Hello World')); });"),
					Remark:      ptrString("test replica"),
					CreatedOn:   ptrString("2024-01-01T00:00:00Z"),
					ModifiedOn:  ptrString("2024-01-01T00:00:00Z"),
				},
			},
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoFunctionReplicaV1()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":      "zone-1234567890",
		"function_id":  "func-abcdefghij",
		"replica_name": "replica-test",
		"content":      "code",
		"remark":       "test",
	})
	d.SetId("zone-1234567890#func-abcdefghij#replica-test")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-1234567890", d.Get("zone_id"))
	assert.Equal(t, "func-abcdefghij", d.Get("function_id"))
	assert.Equal(t, "replica-test", d.Get("replica_name"))
	assert.Equal(t, "test replica", d.Get("remark"))
	assert.Equal(t, "2024-01-01T00:00:00Z", d.Get("create_time"))
	assert.Equal(t, "2024-01-01T00:00:00Z", d.Get("update_time"))
}

// TestTeoFunctionReplicaV1_Read_NotFound tests Read handles replica not found
func TestTeoFunctionReplicaV1_Read_NotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeFunctionReplicas", func(request *teov20220901.DescribeFunctionReplicasRequest) (*teov20220901.DescribeFunctionReplicasResponse, error) {
		resp := teov20220901.NewDescribeFunctionReplicasResponse()
		resp.Response = &teov20220901.DescribeFunctionReplicasResponseParams{
			TotalCount:       ptrInt64(0),
			FunctionReplicas: []*teov20220901.FunctionReplica{},
			RequestId:        ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoFunctionReplicaV1()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":      "zone-1234567890",
		"function_id":  "func-abcdefghij",
		"replica_name": "replica-test",
		"content":      "code",
		"remark":       "test",
	})
	d.SetId("zone-1234567890#func-abcdefghij#replica-test")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}

// TestTeoFunctionReplicaV1_Update_Success tests Update calls ModifyFunctionReplica API
func TestTeoFunctionReplicaV1_Update_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "ModifyFunctionReplicaWithContext", func(ctx interface{}, request *teov20220901.ModifyFunctionReplicaRequest) (*teov20220901.ModifyFunctionReplicaResponse, error) {
		resp := teov20220901.NewModifyFunctionReplicaResponse()
		resp.Response = &teov20220901.ModifyFunctionReplicaResponseParams{
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(teoClient, "DescribeFunctionReplicas", func(request *teov20220901.DescribeFunctionReplicasRequest) (*teov20220901.DescribeFunctionReplicasResponse, error) {
		resp := teov20220901.NewDescribeFunctionReplicasResponse()
		resp.Response = &teov20220901.DescribeFunctionReplicasResponseParams{
			TotalCount: ptrInt64(1),
			FunctionReplicas: []*teov20220901.FunctionReplica{
				{
					FunctionId:  ptrString("func-abcdefghij"),
					ReplicaName: ptrString("replica-test"),
					Content:     ptrString("addEventListener('fetch', e => { e.respondWith(new Response('Hello World Updated')); });"),
					Remark:      ptrString("test replica updated"),
					CreatedOn:   ptrString("2024-01-01T00:00:00Z"),
					ModifiedOn:  ptrString("2024-01-02T00:00:00Z"),
				},
			},
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	// Patch HasChange on *schema.ResourceData - immutable args return false, mutable args return true
	// zone_id, function_id, replica_name are ForceNew (immutable), content and remark are mutable
	patches.ApplyMethodSeq(&schema.ResourceData{}, "HasChange", []gomonkey.OutputCell{
		{Values: gomonkey.Params{false}, Times: 3}, // zone_id, function_id, replica_name => false
		{Values: gomonkey.Params{true}, Times: 1},  // content => true (triggers needChange)
		{Values: gomonkey.Params{true}, Times: 1},  // remark => true (triggers needChange in loop)
		{Values: gomonkey.Params{false}, Times: 2}, // content and remark in HasChange for mutable args check => both changed already processed
	})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoFunctionReplicaV1()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":      "zone-1234567890",
		"function_id":  "func-abcdefghij",
		"replica_name": "replica-test",
		"content":      "addEventListener('fetch', e => { e.respondWith(new Response('Hello World Updated')); });",
		"remark":       "test replica updated",
	})
	d.SetId("zone-1234567890#func-abcdefghij#replica-test")

	err := res.Update(d, meta)
	assert.NoError(t, err)
}

// TestTeoFunctionReplicaV1_Update_APIError tests Update handles API error
func TestTeoFunctionReplicaV1_Update_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "ModifyFunctionReplicaWithContext", func(ctx interface{}, request *teov20220901.ModifyFunctionReplicaRequest) (*teov20220901.ModifyFunctionReplicaResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=InvalidParameter, Message=Invalid content")
	})

	// Patch HasChange - immutable args return false, content and remark return true
	patches.ApplyMethodSeq(&schema.ResourceData{}, "HasChange", []gomonkey.OutputCell{
		{Values: gomonkey.Params{false}, Times: 3}, // zone_id, function_id, replica_name => false
		{Values: gomonkey.Params{true}, Times: 2},  // content and remark => true (triggers modify)
	})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoFunctionReplicaV1()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":      "zone-1234567890",
		"function_id":  "func-abcdefghij",
		"replica_name": "replica-test",
		"content":      "updated code",
		"remark":       "test replica updated",
	})
	d.SetId("zone-1234567890#func-abcdefghij#replica-test")

	err := res.Update(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "InvalidParameter")
}

// TestTeoFunctionReplicaV1_Delete_Success tests Delete removes replica
func TestTeoFunctionReplicaV1_Delete_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DeleteFunctionReplicaWithContext", func(ctx interface{}, request *teov20220901.DeleteFunctionReplicaRequest) (*teov20220901.DeleteFunctionReplicaResponse, error) {
		resp := teov20220901.NewDeleteFunctionReplicaResponse()
		resp.Response = &teov20220901.DeleteFunctionReplicaResponseParams{
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoFunctionReplicaV1()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":      "zone-1234567890",
		"function_id":  "func-abcdefghij",
		"replica_name": "replica-test",
		"content":      "code",
		"remark":       "test",
	})
	d.SetId("zone-1234567890#func-abcdefghij#replica-test")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
}

// TestTeoFunctionReplicaV1_Delete_APIError tests Delete handles API error
func TestTeoFunctionReplicaV1_Delete_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DeleteFunctionReplicaWithContext", func(ctx interface{}, request *teov20220901.DeleteFunctionReplicaRequest) (*teov20220901.DeleteFunctionReplicaResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=ResourceNotFound, Message=Replica not found")
	})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoFunctionReplicaV1()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":      "zone-1234567890",
		"function_id":  "func-abcdefghij",
		"replica_name": "replica-test",
		"content":      "code",
		"remark":       "test",
	})
	d.SetId("zone-1234567890#func-abcdefghij#replica-test")

	err := res.Delete(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ResourceNotFound")
}
