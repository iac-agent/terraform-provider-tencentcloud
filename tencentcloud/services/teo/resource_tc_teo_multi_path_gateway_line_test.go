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
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/teo"
)

// go test ./tencentcloud/services/teo/ -run "TestMultiPathGatewayLine" -v -count=1 -gcflags="all=-l"

// mockMeta implements tccommon.ProviderMeta
type mpglMockMeta struct {
	client *connectivity.TencentCloudClient
}

func (m *mpglMockMeta) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mpglMockMeta{}

func mpglNewMockMeta() *mpglMockMeta {
	return &mpglMockMeta{client: &connectivity.TencentCloudClient{}}
}

func mpglPtrString(s string) *string {
	return &s
}

// TestMultiPathGatewayLine_Create_Success tests Create calls API and sets ID
func TestMultiPathGatewayLine_Create_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(mpglNewMockMeta().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "CreateMultiPathGatewayLine", func(request *teov20220901.CreateMultiPathGatewayLineRequest) (*teov20220901.CreateMultiPathGatewayLineResponse, error) {
		resp := teov20220901.NewCreateMultiPathGatewayLineResponse()
		resp.Response = &teov20220901.CreateMultiPathGatewayLineResponseParams{
			LineId:    mpglPtrString("line-2"),
			RequestId: mpglPtrString("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(teoClient, "DescribeMultiPathGatewayLine", func(request *teov20220901.DescribeMultiPathGatewayLineRequest) (*teov20220901.DescribeMultiPathGatewayLineResponse, error) {
		resp := teov20220901.NewDescribeMultiPathGatewayLineResponse()
		resp.Response = &teov20220901.DescribeMultiPathGatewayLineResponseParams{
			Line: &teov20220901.MultiPathGatewayLine{
				LineId:      mpglPtrString("line-2"),
				LineType:    mpglPtrString("custom"),
				LineAddress: mpglPtrString("1.2.3.4:8080"),
			},
			RequestId: mpglPtrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := mpglNewMockMeta()
	res := teo.ResourceTencentCloudTeoMultiPathGatewayLine()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":      "zone-1234567890",
		"gateway_id":   "gw-abcdefghij",
		"line_type":    "custom",
		"line_address": "1.2.3.4:8080",
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-1234567890#gw-abcdefghij#line-2", d.Id())
	assert.Equal(t, "line-2", d.Get("line_id"))
}

// TestMultiPathGatewayLine_Create_ProxyLine tests Create with proxy line type
func TestMultiPathGatewayLine_Create_ProxyLine(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(mpglNewMockMeta().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "CreateMultiPathGatewayLine", func(request *teov20220901.CreateMultiPathGatewayLineRequest) (*teov20220901.CreateMultiPathGatewayLineResponse, error) {
		resp := teov20220901.NewCreateMultiPathGatewayLineResponse()
		resp.Response = &teov20220901.CreateMultiPathGatewayLineResponseParams{
			LineId:    mpglPtrString("line-1"),
			RequestId: mpglPtrString("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(teoClient, "DescribeMultiPathGatewayLine", func(request *teov20220901.DescribeMultiPathGatewayLineRequest) (*teov20220901.DescribeMultiPathGatewayLineResponse, error) {
		resp := teov20220901.NewDescribeMultiPathGatewayLineResponse()
		resp.Response = &teov20220901.DescribeMultiPathGatewayLineResponseParams{
			Line: &teov20220901.MultiPathGatewayLine{
				LineId:      mpglPtrString("line-1"),
				LineType:    mpglPtrString("proxy"),
				LineAddress: mpglPtrString("1.2.3.4:443"),
				ProxyId:     mpglPtrString("sid-2xzwkzljmm9b"),
				RuleId:      mpglPtrString("rule-2xzwkzljmm9b"),
			},
			RequestId: mpglPtrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := mpglNewMockMeta()
	res := teo.ResourceTencentCloudTeoMultiPathGatewayLine()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":      "zone-1234567890",
		"gateway_id":   "gw-abcdefghij",
		"line_type":    "proxy",
		"line_address": "1.2.3.4:443",
		"proxy_id":     "sid-2xzwkzljmm9b",
		"rule_id":      "rule-2xzwkzljmm9b",
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-1234567890#gw-abcdefghij#line-1", d.Id())
}

// TestMultiPathGatewayLine_Create_APIError tests Create handles API error
func TestMultiPathGatewayLine_Create_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(mpglNewMockMeta().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "CreateMultiPathGatewayLine", func(request *teov20220901.CreateMultiPathGatewayLineRequest) (*teov20220901.CreateMultiPathGatewayLineResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=InvalidParameter, Message=Invalid zone_id")
	})

	meta := mpglNewMockMeta()
	res := teo.ResourceTencentCloudTeoMultiPathGatewayLine()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":      "zone-invalid",
		"gateway_id":   "gw-abcdefghij",
		"line_type":    "custom",
		"line_address": "1.2.3.4:8080",
	})

	err := res.Create(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "InvalidParameter")
}

// TestMultiPathGatewayLine_Read_Success tests Read retrieves line data
func TestMultiPathGatewayLine_Read_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(mpglNewMockMeta().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeMultiPathGatewayLine", func(request *teov20220901.DescribeMultiPathGatewayLineRequest) (*teov20220901.DescribeMultiPathGatewayLineResponse, error) {
		resp := teov20220901.NewDescribeMultiPathGatewayLineResponse()
		resp.Response = &teov20220901.DescribeMultiPathGatewayLineResponseParams{
			Line: &teov20220901.MultiPathGatewayLine{
				LineId:      mpglPtrString("line-2"),
				LineType:    mpglPtrString("custom"),
				LineAddress: mpglPtrString("1.2.3.4:8080"),
			},
			RequestId: mpglPtrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := mpglNewMockMeta()
	res := teo.ResourceTencentCloudTeoMultiPathGatewayLine()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":    "zone-1234567890",
		"gateway_id": "gw-abcdefghij",
		"line_id":    "line-2",
	})
	d.SetId("zone-1234567890#gw-abcdefghij#line-2")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "custom", d.Get("line_type"))
	assert.Equal(t, "1.2.3.4:8080", d.Get("line_address"))
	assert.Equal(t, "line-2", d.Get("line_id"))
}

// TestMultiPathGatewayLine_Read_NotFound tests Read handles not found
func TestMultiPathGatewayLine_Read_NotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(mpglNewMockMeta().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeMultiPathGatewayLine", func(request *teov20220901.DescribeMultiPathGatewayLineRequest) (*teov20220901.DescribeMultiPathGatewayLineResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=ResourceNotFound, Message=line not found")
	})

	meta := mpglNewMockMeta()
	res := teo.ResourceTencentCloudTeoMultiPathGatewayLine()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":    "zone-1234567890",
		"gateway_id": "gw-abcdefghij",
		"line_id":    "line-999",
	})
	d.SetId("zone-1234567890#gw-abcdefghij#line-999")

	err := res.Read(d, meta)
	assert.Error(t, err)
}

// TestMultiPathGatewayLine_Update_Success tests Update calls Modify API
func TestMultiPathGatewayLine_Update_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(mpglNewMockMeta().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "ModifyMultiPathGatewayLine", func(request *teov20220901.ModifyMultiPathGatewayLineRequest) (*teov20220901.ModifyMultiPathGatewayLineResponse, error) {
		resp := teov20220901.NewModifyMultiPathGatewayLineResponse()
		resp.Response = &teov20220901.ModifyMultiPathGatewayLineResponseParams{
			RequestId: mpglPtrString("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(teoClient, "DescribeMultiPathGatewayLine", func(request *teov20220901.DescribeMultiPathGatewayLineRequest) (*teov20220901.DescribeMultiPathGatewayLineResponse, error) {
		resp := teov20220901.NewDescribeMultiPathGatewayLineResponse()
		resp.Response = &teov20220901.DescribeMultiPathGatewayLineResponseParams{
			Line: &teov20220901.MultiPathGatewayLine{
				LineId:      mpglPtrString("line-2"),
				LineType:    mpglPtrString("custom"),
				LineAddress: mpglPtrString("5.6.7.8:9090"),
			},
			RequestId: mpglPtrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := mpglNewMockMeta()
	res := teo.ResourceTencentCloudTeoMultiPathGatewayLine()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":      "zone-1234567890",
		"gateway_id":   "gw-abcdefghij",
		"line_type":    "custom",
		"line_address": "5.6.7.8:9090",
		"line_id":      "line-2",
	})
	d.SetId("zone-1234567890#gw-abcdefghij#line-2")

	err := res.Update(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "5.6.7.8:9090", d.Get("line_address"))
}

// TestMultiPathGatewayLine_Delete_Success tests Delete removes line
func TestMultiPathGatewayLine_Delete_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(mpglNewMockMeta().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DeleteMultiPathGatewayLine", func(request *teov20220901.DeleteMultiPathGatewayLineRequest) (*teov20220901.DeleteMultiPathGatewayLineResponse, error) {
		resp := teov20220901.NewDeleteMultiPathGatewayLineResponse()
		resp.Response = &teov20220901.DeleteMultiPathGatewayLineResponseParams{
			RequestId: mpglPtrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := mpglNewMockMeta()
	res := teo.ResourceTencentCloudTeoMultiPathGatewayLine()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":    "zone-1234567890",
		"gateway_id": "gw-abcdefghij",
		"line_id":    "line-2",
	})
	d.SetId("zone-1234567890#gw-abcdefghij#line-2")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
}

// TestMultiPathGatewayLine_Delete_APIError tests Delete handles API error
func TestMultiPathGatewayLine_Delete_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(mpglNewMockMeta().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DeleteMultiPathGatewayLine", func(request *teov20220901.DeleteMultiPathGatewayLineRequest) (*teov20220901.DeleteMultiPathGatewayLineResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=OperationDenied, Message=proxy line cannot be deleted")
	})

	meta := mpglNewMockMeta()
	res := teo.ResourceTencentCloudTeoMultiPathGatewayLine()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":    "zone-1234567890",
		"gateway_id": "gw-abcdefghij",
		"line_id":    "line-1",
	})
	d.SetId("zone-1234567890#gw-abcdefghij#line-1")

	err := res.Delete(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "OperationDenied")
}

// TestMultiPathGatewayLine_Schema validates schema definition
func TestMultiPathGatewayLine_Schema(t *testing.T) {
	res := teo.ResourceTencentCloudTeoMultiPathGatewayLine()

	assert.NotNil(t, res)
	assert.NotNil(t, res.Create)
	assert.NotNil(t, res.Read)
	assert.NotNil(t, res.Update)
	assert.NotNil(t, res.Delete)
	assert.NotNil(t, res.Importer)

	// Check required fields with ForceNew
	assert.Contains(t, res.Schema, "zone_id")
	zoneId := res.Schema["zone_id"]
	assert.Equal(t, schema.TypeString, zoneId.Type)
	assert.True(t, zoneId.Required)
	assert.True(t, zoneId.ForceNew)

	assert.Contains(t, res.Schema, "gateway_id")
	gatewayId := res.Schema["gateway_id"]
	assert.Equal(t, schema.TypeString, gatewayId.Type)
	assert.True(t, gatewayId.Required)
	assert.True(t, gatewayId.ForceNew)

	// Check required fields without ForceNew
	assert.Contains(t, res.Schema, "line_type")
	lineType := res.Schema["line_type"]
	assert.Equal(t, schema.TypeString, lineType.Type)
	assert.True(t, lineType.Required)
	assert.False(t, lineType.ForceNew)

	assert.Contains(t, res.Schema, "line_address")
	lineAddress := res.Schema["line_address"]
	assert.Equal(t, schema.TypeString, lineAddress.Type)
	assert.True(t, lineAddress.Required)
	assert.False(t, lineAddress.ForceNew)

	// Check optional fields
	assert.Contains(t, res.Schema, "proxy_id")
	proxyId := res.Schema["proxy_id"]
	assert.Equal(t, schema.TypeString, proxyId.Type)
	assert.True(t, proxyId.Optional)
	assert.False(t, proxyId.ForceNew)

	assert.Contains(t, res.Schema, "rule_id")
	ruleId := res.Schema["rule_id"]
	assert.Equal(t, schema.TypeString, ruleId.Type)
	assert.True(t, ruleId.Optional)
	assert.False(t, ruleId.ForceNew)

	// Check computed fields
	assert.Contains(t, res.Schema, "line_id")
	lineId := res.Schema["line_id"]
	assert.Equal(t, schema.TypeString, lineId.Type)
	assert.True(t, lineId.Computed)
}
