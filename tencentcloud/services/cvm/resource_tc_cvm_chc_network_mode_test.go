package cvm_test

import (
	"fmt"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	cvmSDK "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/cvm"
)

// go test ./tencentcloud/services/cvm/ -run "TestChcNetworkMode" -v -count=1 -gcflags="all=-l"

type chcNetworkModeMockMeta struct {
	client *connectivity.TencentCloudClient
}

func (m *chcNetworkModeMockMeta) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &chcNetworkModeMockMeta{}

func newChcNetworkModeMockMeta() *chcNetworkModeMockMeta {
	return &chcNetworkModeMockMeta{client: &connectivity.TencentCloudClient{}}
}

func ptrString(s string) *string {
	return &s
}

// TestChcNetworkMode_Schema validates schema definition
func TestChcNetworkMode_Schema(t *testing.T) {
	res := cvm.ResourceTencentCloudCvmChcNetworkMode()

	assert.NotNil(t, res)
	assert.NotNil(t, res.Create)
	assert.NotNil(t, res.Read)
	assert.NotNil(t, res.Update)
	assert.NotNil(t, res.Delete)

	// Check chc_ids field
	assert.Contains(t, res.Schema, "chc_ids")
	chcIds := res.Schema["chc_ids"]
	assert.Equal(t, schema.TypeList, chcIds.Type)
	assert.True(t, chcIds.Required)
	assert.True(t, chcIds.ForceNew)

	// Check network_mode field
	assert.Contains(t, res.Schema, "network_mode")
	networkMode := res.Schema["network_mode"]
	assert.Equal(t, schema.TypeString, networkMode.Type)
	assert.True(t, networkMode.Required)
	assert.False(t, networkMode.ForceNew)
}

// TestChcNetworkMode_Create_Success tests Create calls API and sets ID
func TestChcNetworkMode_Create_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	cvmClient := &cvmSDK.Client{}
	patches.ApplyMethodReturn(newChcNetworkModeMockMeta().client, "UseCvmClient", cvmClient)

	patches.ApplyMethodFunc(cvmClient, "ModifyChcNetworkMode", func(request *cvmSDK.ModifyChcNetworkModeRequest) (*cvmSDK.ModifyChcNetworkModeResponse, error) {
		resp := cvmSDK.NewModifyChcNetworkModeResponse()
		resp.Response = &cvmSDK.ModifyChcNetworkModeResponseParams{
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(cvmClient, "DescribeChcHosts", func(request *cvmSDK.DescribeChcHostsRequest) (*cvmSDK.DescribeChcHostsResponse, error) {
		resp := cvmSDK.NewDescribeChcHostsResponse()
		resp.Response = &cvmSDK.DescribeChcHostsResponseParams{
			ChcHostSet: []*cvmSDK.ChcHost{
				{
					ChcId:        ptrString("chc-1a2b3c4d"),
					InstanceName: ptrString("test-chc"),
				},
			},
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newChcNetworkModeMockMeta()
	res := cvm.ResourceTencentCloudCvmChcNetworkMode()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"chc_ids":      []interface{}{"chc-1a2b3c4d"},
		"network_mode": "DEPLOY",
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "chc-1a2b3c4d", d.Id())
}

// TestChcNetworkMode_Create_APIError tests Create handles API error
func TestChcNetworkMode_Create_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	cvmClient := &cvmSDK.Client{}
	patches.ApplyMethodReturn(newChcNetworkModeMockMeta().client, "UseCvmClient", cvmClient)

	patches.ApplyMethodFunc(cvmClient, "ModifyChcNetworkMode", func(request *cvmSDK.ModifyChcNetworkModeRequest) (*cvmSDK.ModifyChcNetworkModeResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=InvalidParameterValue.ChcHostsNotFound, Message=CHC hosts not found")
	})

	meta := newChcNetworkModeMockMeta()
	res := cvm.ResourceTencentCloudCvmChcNetworkMode()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"chc_ids":      []interface{}{"chc-invalid"},
		"network_mode": "DEPLOY",
	})

	err := res.Create(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ChcHostsNotFound")
}

// TestChcNetworkMode_Create_NilResponse tests Create handles nil response
func TestChcNetworkMode_Create_NilResponse(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	cvmClient := &cvmSDK.Client{}
	patches.ApplyMethodReturn(newChcNetworkModeMockMeta().client, "UseCvmClient", cvmClient)

	patches.ApplyMethodFunc(cvmClient, "ModifyChcNetworkMode", func(request *cvmSDK.ModifyChcNetworkModeRequest) (*cvmSDK.ModifyChcNetworkModeResponse, error) {
		resp := cvmSDK.NewModifyChcNetworkModeResponse()
		resp.Response = nil
		return resp, nil
	})

	meta := newChcNetworkModeMockMeta()
	res := cvm.ResourceTencentCloudCvmChcNetworkMode()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"chc_ids":      []interface{}{"chc-1a2b3c4d"},
		"network_mode": "DEPLOY",
	})

	err := res.Create(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Response is nil")
}

// TestChcNetworkMode_Read_Success tests Read retrieves CHC host data
func TestChcNetworkMode_Read_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	cvmClient := &cvmSDK.Client{}
	patches.ApplyMethodReturn(newChcNetworkModeMockMeta().client, "UseCvmClient", cvmClient)

	patches.ApplyMethodFunc(cvmClient, "DescribeChcHosts", func(request *cvmSDK.DescribeChcHostsRequest) (*cvmSDK.DescribeChcHostsResponse, error) {
		resp := cvmSDK.NewDescribeChcHostsResponse()
		resp.Response = &cvmSDK.DescribeChcHostsResponseParams{
			ChcHostSet: []*cvmSDK.ChcHost{
				{
					ChcId:        ptrString("chc-1a2b3c4d"),
					InstanceName: ptrString("test-chc"),
				},
			},
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newChcNetworkModeMockMeta()
	res := cvm.ResourceTencentCloudCvmChcNetworkMode()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"chc_ids":      []interface{}{"chc-1a2b3c4d"},
		"network_mode": "DEPLOY",
	})
	d.SetId("chc-1a2b3c4d")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "chc-1a2b3c4d", d.Id())
	// network_mode should be preserved from state since ChcHost doesn't have NetworkMode field
	assert.Equal(t, "DEPLOY", d.Get("network_mode"))
}

// TestChcNetworkMode_Read_NotFound tests Read handles CHC host not found
func TestChcNetworkMode_Read_NotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	cvmClient := &cvmSDK.Client{}
	patches.ApplyMethodReturn(newChcNetworkModeMockMeta().client, "UseCvmClient", cvmClient)

	patches.ApplyMethodFunc(cvmClient, "DescribeChcHosts", func(request *cvmSDK.DescribeChcHostsRequest) (*cvmSDK.DescribeChcHostsResponse, error) {
		resp := cvmSDK.NewDescribeChcHostsResponse()
		resp.Response = &cvmSDK.DescribeChcHostsResponseParams{
			ChcHostSet: []*cvmSDK.ChcHost{},
			RequestId:  ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newChcNetworkModeMockMeta()
	res := cvm.ResourceTencentCloudCvmChcNetworkMode()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"chc_ids":      []interface{}{"chc-1a2b3c4d"},
		"network_mode": "DEPLOY",
	})
	d.SetId("chc-1a2b3c4d")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}

// TestChcNetworkMode_Update_Success tests Update calls API when network_mode changes
func TestChcNetworkMode_Update_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	cvmClient := &cvmSDK.Client{}
	patches.ApplyMethodReturn(newChcNetworkModeMockMeta().client, "UseCvmClient", cvmClient)

	patches.ApplyMethodFunc(cvmClient, "ModifyChcNetworkMode", func(request *cvmSDK.ModifyChcNetworkModeRequest) (*cvmSDK.ModifyChcNetworkModeResponse, error) {
		resp := cvmSDK.NewModifyChcNetworkModeResponse()
		resp.Response = &cvmSDK.ModifyChcNetworkModeResponseParams{
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(cvmClient, "DescribeChcHosts", func(request *cvmSDK.DescribeChcHostsRequest) (*cvmSDK.DescribeChcHostsResponse, error) {
		resp := cvmSDK.NewDescribeChcHostsResponse()
		resp.Response = &cvmSDK.DescribeChcHostsResponseParams{
			ChcHostSet: []*cvmSDK.ChcHost{
				{
					ChcId:        ptrString("chc-1a2b3c4d"),
					InstanceName: ptrString("test-chc"),
				},
			},
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newChcNetworkModeMockMeta()
	res := cvm.ResourceTencentCloudCvmChcNetworkMode()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"chc_ids":      []interface{}{"chc-1a2b3c4d"},
		"network_mode": "BUSINESS",
	})
	d.SetId("chc-1a2b3c4d")

	err := res.Update(d, meta)
	assert.NoError(t, err)
}

// TestChcNetworkMode_Update_APIError tests Update handles API error
func TestChcNetworkMode_Update_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	cvmClient := &cvmSDK.Client{}
	patches.ApplyMethodReturn(newChcNetworkModeMockMeta().client, "UseCvmClient", cvmClient)

	patches.ApplyMethodFunc(cvmClient, "ModifyChcNetworkMode", func(request *cvmSDK.ModifyChcNetworkModeRequest) (*cvmSDK.ModifyChcNetworkModeResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=OperationDenied.ChcHostStateNotSupported, Message=CHC host state not supported")
	})

	meta := newChcNetworkModeMockMeta()
	res := cvm.ResourceTencentCloudCvmChcNetworkMode()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"chc_ids":      []interface{}{"chc-1a2b3c4d"},
		"network_mode": "BUSINESS",
	})
	d.SetId("chc-1a2b3c4d")

	err := res.Update(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ChcHostStateNotSupported")
}

// TestChcNetworkMode_Delete_Success tests Delete only removes from state
func TestChcNetworkMode_Delete_Success(t *testing.T) {
	meta := newChcNetworkModeMockMeta()
	res := cvm.ResourceTencentCloudCvmChcNetworkMode()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"chc_ids":      []interface{}{"chc-1a2b3c4d"},
		"network_mode": "DEPLOY",
	})
	d.SetId("chc-1a2b3c4d")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}

// TestChcNetworkMode_Create_MultipleChcIds tests Create with multiple CHC IDs
func TestChcNetworkMode_Create_MultipleChcIds(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	cvmClient := &cvmSDK.Client{}
	patches.ApplyMethodReturn(newChcNetworkModeMockMeta().client, "UseCvmClient", cvmClient)

	patches.ApplyMethodFunc(cvmClient, "ModifyChcNetworkMode", func(request *cvmSDK.ModifyChcNetworkModeRequest) (*cvmSDK.ModifyChcNetworkModeResponse, error) {
		resp := cvmSDK.NewModifyChcNetworkModeResponse()
		resp.Response = &cvmSDK.ModifyChcNetworkModeResponseParams{
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(cvmClient, "DescribeChcHosts", func(request *cvmSDK.DescribeChcHostsRequest) (*cvmSDK.DescribeChcHostsResponse, error) {
		resp := cvmSDK.NewDescribeChcHostsResponse()
		resp.Response = &cvmSDK.DescribeChcHostsResponseParams{
			ChcHostSet: []*cvmSDK.ChcHost{
				{
					ChcId:        ptrString("chc-1a2b3c4d"),
					InstanceName: ptrString("test-chc-1"),
				},
				{
					ChcId:        ptrString("chc-5e6f7g8h"),
					InstanceName: ptrString("test-chc-2"),
				},
			},
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newChcNetworkModeMockMeta()
	res := cvm.ResourceTencentCloudCvmChcNetworkMode()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"chc_ids":      []interface{}{"chc-1a2b3c4d", "chc-5e6f7g8h"},
		"network_mode": "BUSINESS",
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "chc-1a2b3c4d#chc-5e6f7g8h", d.Id())
}
