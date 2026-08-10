package fl_test

import (
	"fmt"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/fl"
)

// mockMeta implements tccommon.ProviderMeta
type mockMeta struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMeta) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMeta{}

func newMockMeta() *mockMeta {
	return &mockMeta{client: &connectivity.TencentCloudClient{}}
}

func ptrString(s string) *string {
	return &s
}

// TestFlowLogConfigCreate_EnableTrue tests Create with enable=true
func TestFlowLogConfigCreate_EnableTrue(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	vpcClient := &vpc.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseVpcClient", vpcClient)

	patches.ApplyMethodFunc(vpcClient, "EnableFlowLogs", func(_ *vpc.Client, request *vpc.EnableFlowLogsRequest) (*vpc.EnableFlowLogsResponse, error) {
		assert.Equal(t, 1, len(request.FlowLogIds))
		assert.Equal(t, "fl-test123", *request.FlowLogIds[0])
		resp := vpc.NewEnableFlowLogsResponse()
		resp.Response = &vpc.EnableFlowLogsResponseParams{
			RequestId: ptrString("fake-req-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(vpcClient, "DescribeFlowLogs", func(_ *vpc.Client, request *vpc.DescribeFlowLogsRequest) (*vpc.DescribeFlowLogsResponse, error) {
		assert.Equal(t, "fl-test123", *request.FlowLogId)
		resp := vpc.NewDescribeFlowLogsResponse()
		enable := true
		resp.Response = &vpc.DescribeFlowLogsResponseParams{
			FlowLog: []*vpc.FlowLog{
				{
					FlowLogId: ptrString("fl-test123"),
					Enable:    &enable,
				},
			},
			RequestId: ptrString("fake-req-id"),
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := fl.ResourceTencentCloudVpcFlowLogConfig()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"flow_log_id": "fl-test123",
		"enable":      true,
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "fl-test123", d.Id())
	assert.True(t, d.Get("enable").(bool))
}

// TestFlowLogConfigCreate_EnableFalse tests Create with enable=false
func TestFlowLogConfigCreate_EnableFalse(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	vpcClient := &vpc.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseVpcClient", vpcClient)

	patches.ApplyMethodFunc(vpcClient, "DisableFlowLogs", func(_ *vpc.Client, request *vpc.DisableFlowLogsRequest) (*vpc.DisableFlowLogsResponse, error) {
		assert.Equal(t, 1, len(request.FlowLogIds))
		assert.Equal(t, "fl-test456", *request.FlowLogIds[0])
		resp := vpc.NewDisableFlowLogsResponse()
		resp.Response = &vpc.DisableFlowLogsResponseParams{
			RequestId: ptrString("fake-req-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(vpcClient, "DescribeFlowLogs", func(_ *vpc.Client, request *vpc.DescribeFlowLogsRequest) (*vpc.DescribeFlowLogsResponse, error) {
		assert.Equal(t, "fl-test456", *request.FlowLogId)
		resp := vpc.NewDescribeFlowLogsResponse()
		enable := false
		resp.Response = &vpc.DescribeFlowLogsResponseParams{
			FlowLog: []*vpc.FlowLog{
				{
					FlowLogId: ptrString("fl-test456"),
					Enable:    &enable,
				},
			},
			RequestId: ptrString("fake-req-id"),
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := fl.ResourceTencentCloudVpcFlowLogConfig()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"flow_log_id": "fl-test456",
		"enable":      false,
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "fl-test456", d.Id())
	assert.False(t, d.Get("enable").(bool))
}

// TestFlowLogConfigRead_EnablePresent tests Read with Enable field present
func TestFlowLogConfigRead_EnablePresent(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	vpcClient := &vpc.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseVpcClient", vpcClient)

	patches.ApplyMethodFunc(vpcClient, "DescribeFlowLogs", func(_ *vpc.Client, request *vpc.DescribeFlowLogsRequest) (*vpc.DescribeFlowLogsResponse, error) {
		assert.Equal(t, "fl-read-test", *request.FlowLogId)
		resp := vpc.NewDescribeFlowLogsResponse()
		enable := true
		resp.Response = &vpc.DescribeFlowLogsResponseParams{
			FlowLog: []*vpc.FlowLog{
				{
					FlowLogId: ptrString("fl-read-test"),
					Enable:    &enable,
				},
			},
			RequestId: ptrString("fake-req-id"),
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := fl.ResourceTencentCloudVpcFlowLogConfig()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"flow_log_id": "fl-read-test",
		"enable":      false,
	})
	d.SetId("fl-read-test")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.True(t, d.Get("enable").(bool))
}

// TestFlowLogConfigRead_EmptyResult tests Read with empty result (flow log not found)
func TestFlowLogConfigRead_EmptyResult(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	vpcClient := &vpc.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseVpcClient", vpcClient)

	patches.ApplyMethodFunc(vpcClient, "DescribeFlowLogs", func(_ *vpc.Client, request *vpc.DescribeFlowLogsRequest) (*vpc.DescribeFlowLogsResponse, error) {
		assert.Equal(t, "fl-not-found", *request.FlowLogId)
		resp := vpc.NewDescribeFlowLogsResponse()
		resp.Response = &vpc.DescribeFlowLogsResponseParams{
			FlowLog:   []*vpc.FlowLog{},
			RequestId: ptrString("fake-req-id"),
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := fl.ResourceTencentCloudVpcFlowLogConfig()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"flow_log_id": "fl-not-found",
		"enable":      true,
	})
	d.SetId("fl-not-found")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}

// TestFlowLogConfigUpdate_DisableToEnable tests Update from disabled to enabled
func TestFlowLogConfigUpdate_DisableToEnable(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	vpcClient := &vpc.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseVpcClient", vpcClient)

	enableCalled := false
	patches.ApplyMethodFunc(vpcClient, "EnableFlowLogs", func(_ *vpc.Client, request *vpc.EnableFlowLogsRequest) (*vpc.EnableFlowLogsResponse, error) {
		assert.Equal(t, 1, len(request.FlowLogIds))
		assert.Equal(t, "fl-update", *request.FlowLogIds[0])
		enableCalled = true
		resp := vpc.NewEnableFlowLogsResponse()
		resp.Response = &vpc.EnableFlowLogsResponseParams{
			RequestId: ptrString("fake-req-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(vpcClient, "DescribeFlowLogs", func(_ *vpc.Client, request *vpc.DescribeFlowLogsRequest) (*vpc.DescribeFlowLogsResponse, error) {
		resp := vpc.NewDescribeFlowLogsResponse()
		enable := true
		resp.Response = &vpc.DescribeFlowLogsResponseParams{
			FlowLog: []*vpc.FlowLog{
				{
					FlowLogId: ptrString("fl-update"),
					Enable:    &enable,
				},
			},
			RequestId: ptrString("fake-req-id"),
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := fl.ResourceTencentCloudVpcFlowLogConfig()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"flow_log_id": "fl-update",
		"enable":      true,
	})
	d.SetId("fl-update")

	err := res.Update(d, meta)
	assert.NoError(t, err)
	assert.True(t, enableCalled)
	assert.True(t, d.Get("enable").(bool))
}

// TestFlowLogConfigUpdate_EnableToDisable tests Update from enabled to disabled
func TestFlowLogConfigUpdate_EnableToDisable(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	vpcClient := &vpc.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseVpcClient", vpcClient)

	disableCalled := false
	patches.ApplyMethodFunc(vpcClient, "DisableFlowLogs", func(_ *vpc.Client, request *vpc.DisableFlowLogsRequest) (*vpc.DisableFlowLogsResponse, error) {
		assert.Equal(t, 1, len(request.FlowLogIds))
		assert.Equal(t, "fl-update2", *request.FlowLogIds[0])
		disableCalled = true
		resp := vpc.NewDisableFlowLogsResponse()
		resp.Response = &vpc.DisableFlowLogsResponseParams{
			RequestId: ptrString("fake-req-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(vpcClient, "DescribeFlowLogs", func(_ *vpc.Client, request *vpc.DescribeFlowLogsRequest) (*vpc.DescribeFlowLogsResponse, error) {
		resp := vpc.NewDescribeFlowLogsResponse()
		enable := false
		resp.Response = &vpc.DescribeFlowLogsResponseParams{
			FlowLog: []*vpc.FlowLog{
				{
					FlowLogId: ptrString("fl-update2"),
					Enable:    &enable,
				},
			},
			RequestId: ptrString("fake-req-id"),
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := fl.ResourceTencentCloudVpcFlowLogConfig()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"flow_log_id": "fl-update2",
		"enable":      false,
	})
	d.SetId("fl-update2")

	err := res.Update(d, meta)
	assert.NoError(t, err)
	assert.True(t, disableCalled)
	assert.False(t, d.Get("enable").(bool))
}

// TestFlowLogConfigDelete_NoOp tests Delete is a no-op
func TestFlowLogConfigDelete_NoOp(t *testing.T) {
	res := fl.ResourceTencentCloudVpcFlowLogConfig()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"flow_log_id": "fl-delete",
		"enable":      true,
	})
	d.SetId("fl-delete")

	err := res.Delete(d, nil)
	assert.NoError(t, err)
}

// TestFlowLogConfigUpdate_APIError tests Update handles API errors
func TestFlowLogConfigUpdate_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	vpcClient := &vpc.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseVpcClient", vpcClient)

	patches.ApplyMethodFunc(vpcClient, "EnableFlowLogs", func(_ *vpc.Client, request *vpc.EnableFlowLogsRequest) (*vpc.EnableFlowLogsResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=InternalError, Message=Internal error")
	})

	meta := newMockMeta()
	res := fl.ResourceTencentCloudVpcFlowLogConfig()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"flow_log_id": "fl-error",
		"enable":      true,
	})
	d.SetId("fl-error")

	err := res.Update(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "InternalError")
}

// TestFlowLogConfigSchema tests schema definition
func TestFlowLogConfigSchema(t *testing.T) {
	res := fl.ResourceTencentCloudVpcFlowLogConfig()

	assert.NotNil(t, res)
	assert.NotNil(t, res.Create)
	assert.NotNil(t, res.Read)
	assert.NotNil(t, res.Update)
	assert.NotNil(t, res.Delete)
	assert.NotNil(t, res.Importer)

	// Verify flow_log_id schema
	flowLogId := res.Schema["flow_log_id"]
	assert.Equal(t, schema.TypeString, flowLogId.Type)
	assert.True(t, flowLogId.Required)
	assert.True(t, flowLogId.ForceNew)

	// Verify enable schema
	enable := res.Schema["enable"]
	assert.Equal(t, schema.TypeBool, enable.Type)
	assert.True(t, enable.Required)
}
