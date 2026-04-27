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

// go test ./tencentcloud/services/teo/ -run "TestMultiPathGatewayLine" -v -count=1 -gcflags="all=-l"

// TestMultiPathGatewayLine_Create_Success tests Create calls API and sets ID
func TestMultiPathGatewayLine_Create_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "CreateMultiPathGatewayLineWithContext", func(_ context.Context, request *teov20220901.CreateMultiPathGatewayLineRequest) (*teov20220901.CreateMultiPathGatewayLineResponse, error) {
		assert.Equal(t, "zone-12345678", *request.ZoneId)
		assert.Equal(t, "gw-abcdefgh", *request.GatewayId)
		assert.Equal(t, "custom", *request.LineType)
		assert.Equal(t, "1.2.3.4:8080", *request.LineAddress)

		resp := teov20220901.NewCreateMultiPathGatewayLineResponse()
		resp.Response = &teov20220901.CreateMultiPathGatewayLineResponseParams{
			LineId:    ptrString("line-2"),
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(teoClient, "DescribeMultiPathGatewayLine", func(request *teov20220901.DescribeMultiPathGatewayLineRequest) (*teov20220901.DescribeMultiPathGatewayLineResponse, error) {
		resp := teov20220901.NewDescribeMultiPathGatewayLineResponse()
		resp.Response = &teov20220901.DescribeMultiPathGatewayLineResponseParams{
			Line: &teov20220901.MultiPathGatewayLine{
				LineId:      ptrString("line-2"),
				LineType:    ptrString("custom"),
				LineAddress: ptrString("1.2.3.4:8080"),
			},
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoMultiPathGatewayLine()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":      "zone-12345678",
		"gateway_id":   "gw-abcdefgh",
		"line_type":    "custom",
		"line_address": "1.2.3.4:8080",
		"proxy_id":     "",
		"rule_id":      "",
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-12345678#gw-abcdefgh#line-2", d.Id())
	assert.Equal(t, "line-2", d.Get("line_id"))
	assert.Equal(t, "custom", d.Get("line_type"))
	assert.Equal(t, "1.2.3.4:8080", d.Get("line_address"))
}

// TestMultiPathGatewayLine_Create_WithProxy tests Create with proxy type
func TestMultiPathGatewayLine_Create_WithProxy(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "CreateMultiPathGatewayLineWithContext", func(_ context.Context, request *teov20220901.CreateMultiPathGatewayLineRequest) (*teov20220901.CreateMultiPathGatewayLineResponse, error) {
		assert.Equal(t, "proxy", *request.LineType)
		assert.Equal(t, "proxy-12345", *request.ProxyId)
		assert.Equal(t, "rule-67890", *request.RuleId)

		resp := teov20220901.NewCreateMultiPathGatewayLineResponse()
		resp.Response = &teov20220901.CreateMultiPathGatewayLineResponseParams{
			LineId:    ptrString("line-3"),
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(teoClient, "DescribeMultiPathGatewayLine", func(request *teov20220901.DescribeMultiPathGatewayLineRequest) (*teov20220901.DescribeMultiPathGatewayLineResponse, error) {
		resp := teov20220901.NewDescribeMultiPathGatewayLineResponse()
		resp.Response = &teov20220901.DescribeMultiPathGatewayLineResponseParams{
			Line: &teov20220901.MultiPathGatewayLine{
				LineId:      ptrString("line-3"),
				LineType:    ptrString("proxy"),
				LineAddress: ptrString("5.6.7.8:443"),
				ProxyId:     ptrString("proxy-12345"),
				RuleId:      ptrString("rule-67890"),
			},
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoMultiPathGatewayLine()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":      "zone-12345678",
		"gateway_id":   "gw-abcdefgh",
		"line_type":    "proxy",
		"line_address": "5.6.7.8:443",
		"proxy_id":     "proxy-12345",
		"rule_id":      "rule-67890",
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-12345678#gw-abcdefgh#line-3", d.Id())
}

// TestMultiPathGatewayLine_Create_APIError tests Create handles API error
func TestMultiPathGatewayLine_Create_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "CreateMultiPathGatewayLineWithContext", func(_ context.Context, request *teov20220901.CreateMultiPathGatewayLineRequest) (*teov20220901.CreateMultiPathGatewayLineResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=InvalidParameter, Message=Invalid zone_id")
	})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoMultiPathGatewayLine()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":      "zone-invalid",
		"gateway_id":   "gw-invalid",
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
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeMultiPathGatewayLine", func(request *teov20220901.DescribeMultiPathGatewayLineRequest) (*teov20220901.DescribeMultiPathGatewayLineResponse, error) {
		assert.Equal(t, "zone-12345678", *request.ZoneId)
		assert.Equal(t, "gw-abcdefgh", *request.GatewayId)
		assert.Equal(t, "line-2", *request.LineId)

		resp := teov20220901.NewDescribeMultiPathGatewayLineResponse()
		resp.Response = &teov20220901.DescribeMultiPathGatewayLineResponseParams{
			Line: &teov20220901.MultiPathGatewayLine{
				LineId:      ptrString("line-2"),
				LineType:    ptrString("custom"),
				LineAddress: ptrString("1.2.3.4:8080"),
			},
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoMultiPathGatewayLine()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":      "zone-12345678",
		"gateway_id":   "gw-abcdefgh",
		"line_type":    "custom",
		"line_address": "1.2.3.4:8080",
		"line_id":      "line-2",
	})
	d.SetId("zone-12345678#gw-abcdefgh#line-2")

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
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeMultiPathGatewayLine", func(request *teov20220901.DescribeMultiPathGatewayLineRequest) (*teov20220901.DescribeMultiPathGatewayLineResponse, error) {
		resp := teov20220901.NewDescribeMultiPathGatewayLineResponse()
		resp.Response = &teov20220901.DescribeMultiPathGatewayLineResponseParams{
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoMultiPathGatewayLine()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":      "zone-12345678",
		"gateway_id":   "gw-abcdefgh",
		"line_type":    "custom",
		"line_address": "1.2.3.4:8080",
	})
	d.SetId("zone-12345678#gw-abcdefgh#line-99")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}

// TestMultiPathGatewayLine_Update_Success tests Update modifies line
func TestMultiPathGatewayLine_Update_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "ModifyMultiPathGatewayLineWithContext", func(_ context.Context, request *teov20220901.ModifyMultiPathGatewayLineRequest) (*teov20220901.ModifyMultiPathGatewayLineResponse, error) {
		assert.Equal(t, "zone-12345678", *request.ZoneId)
		assert.Equal(t, "gw-abcdefgh", *request.GatewayId)
		assert.Equal(t, "line-2", *request.LineId)
		assert.Equal(t, "5.6.7.8:9090", *request.LineAddress)

		resp := teov20220901.NewModifyMultiPathGatewayLineResponse()
		resp.Response = &teov20220901.ModifyMultiPathGatewayLineResponseParams{
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(teoClient, "DescribeMultiPathGatewayLine", func(request *teov20220901.DescribeMultiPathGatewayLineRequest) (*teov20220901.DescribeMultiPathGatewayLineResponse, error) {
		resp := teov20220901.NewDescribeMultiPathGatewayLineResponse()
		resp.Response = &teov20220901.DescribeMultiPathGatewayLineResponseParams{
			Line: &teov20220901.MultiPathGatewayLine{
				LineId:      ptrString("line-2"),
				LineType:    ptrString("custom"),
				LineAddress: ptrString("5.6.7.8:9090"),
			},
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoMultiPathGatewayLine()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":      "zone-12345678",
		"gateway_id":   "gw-abcdefgh",
		"line_type":    "custom",
		"line_address": "5.6.7.8:9090",
		"line_id":      "line-2",
	})
	d.SetId("zone-12345678#gw-abcdefgh#line-2")

	err := res.Update(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "5.6.7.8:9090", d.Get("line_address"))
}

// TestMultiPathGatewayLine_Delete_Success tests Delete removes line
func TestMultiPathGatewayLine_Delete_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DeleteMultiPathGatewayLineWithContext", func(_ context.Context, request *teov20220901.DeleteMultiPathGatewayLineRequest) (*teov20220901.DeleteMultiPathGatewayLineResponse, error) {
		assert.Equal(t, "zone-12345678", *request.ZoneId)
		assert.Equal(t, "gw-abcdefgh", *request.GatewayId)
		assert.Equal(t, "line-2", *request.LineId)

		resp := teov20220901.NewDeleteMultiPathGatewayLineResponse()
		resp.Response = &teov20220901.DeleteMultiPathGatewayLineResponseParams{
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoMultiPathGatewayLine()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":      "zone-12345678",
		"gateway_id":   "gw-abcdefgh",
		"line_type":    "custom",
		"line_address": "1.2.3.4:8080",
		"line_id":      "line-2",
	})
	d.SetId("zone-12345678#gw-abcdefgh#line-2")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
}

// TestMultiPathGatewayLine_Delete_APIError tests Delete handles API error
func TestMultiPathGatewayLine_Delete_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DeleteMultiPathGatewayLineWithContext", func(_ context.Context, request *teov20220901.DeleteMultiPathGatewayLineRequest) (*teov20220901.DeleteMultiPathGatewayLineResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=ResourceNotFound, Message=Line not found")
	})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoMultiPathGatewayLine()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":      "zone-12345678",
		"gateway_id":   "gw-abcdefgh",
		"line_type":    "custom",
		"line_address": "1.2.3.4:8080",
		"line_id":      "line-2",
	})
	d.SetId("zone-12345678#gw-abcdefgh#line-2")

	err := res.Delete(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ResourceNotFound")
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
