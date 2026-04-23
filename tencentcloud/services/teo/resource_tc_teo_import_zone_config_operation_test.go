package teo_test

import (
	"fmt"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/teo"
)

// go test ./tencentcloud/services/teo/ -run "TestImportZoneConfigOperation" -v -count=1 -gcflags="all=-l"

// TestImportZoneConfigOperation_Success tests Create calls ImportZoneConfig API, polls until success, and sets ID
func TestImportZoneConfigOperation_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	// Patch UseTeoV20220901Client to return a non-nil client
	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)

	// Patch ImportZoneConfig to return success with TaskId
	patches.ApplyMethodFunc(teoClient, "ImportZoneConfig", func(request *teov20220901.ImportZoneConfigRequest) (*teov20220901.ImportZoneConfigResponse, error) {
		assert.Equal(t, "zone-12345678", *request.ZoneId)
		assert.NotNil(t, request.Content)
		resp := teov20220901.NewImportZoneConfigResponse()
		resp.Response = &teov20220901.ImportZoneConfigResponseParams{
			TaskId:    ptrString("task-87654321"),
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	// Patch DescribeZoneConfigImportResult to return success status immediately
	patches.ApplyMethodFunc(teoClient, "DescribeZoneConfigImportResult", func(request *teov20220901.DescribeZoneConfigImportResultRequest) (*teov20220901.DescribeZoneConfigImportResultResponse, error) {
		assert.Equal(t, "zone-12345678", *request.ZoneId)
		assert.Equal(t, "task-87654321", *request.TaskId)
		resp := teov20220901.NewDescribeZoneConfigImportResultResponse()
		resp.Response = &teov20220901.DescribeZoneConfigImportResultResponseParams{
			Status:    ptrString("success"),
			Message:   ptrString(""),
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoImportZoneConfigOperation()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-12345678",
		"content": `{"ZoneId":"zone-12345678"}`,
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-12345678"+tccommon.FILED_SP+"task-87654321", d.Id())
	assert.Equal(t, "task-87654321", d.Get("task_id").(string))
}

// TestImportZoneConfigOperation_TaskFailure tests Create handles async task failure
func TestImportZoneConfigOperation_TaskFailure(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)

	// Patch ImportZoneConfig to return success with TaskId
	patches.ApplyMethodFunc(teoClient, "ImportZoneConfig", func(request *teov20220901.ImportZoneConfigRequest) (*teov20220901.ImportZoneConfigResponse, error) {
		resp := teov20220901.NewImportZoneConfigResponse()
		resp.Response = &teov20220901.ImportZoneConfigResponseParams{
			TaskId:    ptrString("task-failed"),
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	// Patch DescribeZoneConfigImportResult to return failure status
	patches.ApplyMethodFunc(teoClient, "DescribeZoneConfigImportResult", func(request *teov20220901.DescribeZoneConfigImportResultRequest) (*teov20220901.DescribeZoneConfigImportResultResponse, error) {
		resp := teov20220901.NewDescribeZoneConfigImportResultResponse()
		resp.Response = &teov20220901.DescribeZoneConfigImportResultResponseParams{
			Status:    ptrString("failure"),
			Message:   ptrString("invalid configuration format"),
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoImportZoneConfigOperation()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-12345678",
		"content": `{"invalid": true}`,
	})

	err := res.Create(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failure")
}

// TestImportZoneConfigOperation_APIError tests Create handles ImportZoneConfig API error
func TestImportZoneConfigOperation_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)

	// Patch ImportZoneConfig to return API error
	patches.ApplyMethodFunc(teoClient, "ImportZoneConfig", func(request *teov20220901.ImportZoneConfigRequest) (*teov20220901.ImportZoneConfigResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=InvalidParameter, Message=Invalid zone_id")
	})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoImportZoneConfigOperation()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-invalid",
		"content": `{"ZoneId":"zone-invalid"}`,
	})

	err := res.Create(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "InvalidParameter")
}

// TestImportZoneConfigOperation_Read tests Read is no-op
func TestImportZoneConfigOperation_Read(t *testing.T) {
	res := teo.ResourceTencentCloudTeoImportZoneConfigOperation()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-12345678",
		"content": `{"ZoneId":"zone-12345678"}`,
	})
	d.SetId("zone-12345678" + tccommon.FILED_SP + "task-87654321")

	err := res.Read(d, nil)
	assert.NoError(t, err)
}

// TestImportZoneConfigOperation_Delete tests Delete is no-op
func TestImportZoneConfigOperation_Delete(t *testing.T) {
	res := teo.ResourceTencentCloudTeoImportZoneConfigOperation()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-12345678",
		"content": `{"ZoneId":"zone-12345678"}`,
	})
	d.SetId("zone-12345678" + tccommon.FILED_SP + "task-87654321")

	err := res.Delete(d, nil)
	assert.NoError(t, err)
}

// TestImportZoneConfigOperation_Schema validates schema definition
func TestImportZoneConfigOperation_Schema(t *testing.T) {
	res := teo.ResourceTencentCloudTeoImportZoneConfigOperation()

	assert.NotNil(t, res)
	assert.NotNil(t, res.Create)
	assert.NotNil(t, res.Read)
	assert.Nil(t, res.Update)
	assert.NotNil(t, res.Delete)

	assert.Contains(t, res.Schema, "zone_id")
	assert.Contains(t, res.Schema, "content")
	assert.Contains(t, res.Schema, "task_id")

	zoneId := res.Schema["zone_id"]
	assert.Equal(t, schema.TypeString, zoneId.Type)
	assert.True(t, zoneId.Required)
	assert.True(t, zoneId.ForceNew)

	content := res.Schema["content"]
	assert.Equal(t, schema.TypeString, content.Type)
	assert.True(t, content.Required)
	assert.True(t, content.ForceNew)

	taskId := res.Schema["task_id"]
	assert.Equal(t, schema.TypeString, taskId.Type)
	assert.True(t, taskId.Computed)
}
