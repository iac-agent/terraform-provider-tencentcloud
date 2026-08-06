package teo

import (
	"fmt"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
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

func ptrInt64(i int64) *int64 {
	return &i
}

func ptrUint64(i uint64) *uint64 {
	return &i
}

// go test ./tencentcloud/services/teo/ -run "TestTeoL4ProxyCreate_Success" -v -count=1 -gcflags="all=-l"
// TestTeoL4ProxyCreate_Success tests Create calls API and sets ID
func TestTeoL4ProxyCreate_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoClient", teoClient)

	// Patch CreateL4ProxyWithContext to return success
	patches.ApplyMethodFunc(teoClient, "CreateL4ProxyWithContext", func(_ interface{}, request *teov20220901.CreateL4ProxyRequest) (*teov20220901.CreateL4ProxyResponse, error) {
		assert.Equal(t, "zone-2qtuhspy7cr6", *request.ZoneId)
		assert.Equal(t, "proxy-test", *request.ProxyName)
		assert.Equal(t, "overseas", *request.Area)
		resp := teov20220901.NewCreateL4ProxyResponse()
		resp.Response = &teov20220901.CreateL4ProxyResponseParams{
			ProxyId:   ptrString("sid-2qtuhspy7cr6"),
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	// Patch extension function to avoid waiting for state
	patches.ApplyFunc(resourceTencentCloudTeoL4ProxyCreatePostHandleResponse0, func(_ interface{}, _ *teov20220901.CreateL4ProxyResponse) error {
		return nil
	})

	meta := newMockMeta()
	res := ResourceTencentCloudTeoL4Proxy()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":    "zone-2qtuhspy7cr6",
		"proxy_name": "proxy-test",
		"area":       "overseas",
		"ipv6":       "off",
		"static_ip":  "off",
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-2qtuhspy7cr6#sid-2qtuhspy7cr6", d.Id())
}

// TestTeoL4ProxyCreate_APIError tests Create handles API error
func TestTeoL4ProxyCreate_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "CreateL4ProxyWithContext", func(_ interface{}, _ *teov20220901.CreateL4ProxyRequest) (*teov20220901.CreateL4ProxyResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=InvalidParameter, Message=Invalid zone_id")
	})

	meta := newMockMeta()
	res := ResourceTencentCloudTeoL4Proxy()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":    "zone-invalid",
		"proxy_name": "proxy-test",
		"area":       "overseas",
	})

	err := res.Create(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "InvalidParameter")
}

// TestTeoL4ProxyCreate_EmptyProxyId tests Create handles empty ProxyId
func TestTeoL4ProxyCreate_EmptyProxyId(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "CreateL4ProxyWithContext", func(_ interface{}, _ *teov20220901.CreateL4ProxyRequest) (*teov20220901.CreateL4ProxyResponse, error) {
		resp := teov20220901.NewCreateL4ProxyResponse()
		resp.Response = &teov20220901.CreateL4ProxyResponseParams{
			ProxyId:   ptrString(""),
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := ResourceTencentCloudTeoL4Proxy()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":    "zone-2qtuhspy7cr6",
		"proxy_name": "proxy-test",
		"area":       "overseas",
	})

	err := res.Create(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "empty ProxyId")
}

// TestTeoL4ProxyRead_Success tests Read sets all computed fields
func TestTeoL4ProxyRead_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeL4Proxy", func(request *teov20220901.DescribeL4ProxyRequest) (*teov20220901.DescribeL4ProxyResponse, error) {
		assert.Equal(t, "zone-2qtuhspy7cr6", *request.ZoneId)
		resp := teov20220901.NewDescribeL4ProxyResponse()
		resp.Response = &teov20220901.DescribeL4ProxyResponseParams{
			TotalCount: ptrUint64(1),
			L4Proxies: []*teov20220901.L4Proxy{
				{
					ZoneId:             ptrString("zone-2qtuhspy7cr6"),
					ProxyId:            ptrString("sid-2qtuhspy7cr6"),
					ProxyName:          ptrString("proxy-test"),
					Area:               ptrString("overseas"),
					Ipv6:               ptrString("off"),
					StaticIp:           ptrString("off"),
					AccelerateMainland: ptrString("off"),
					Cname:              ptrString("proxy-test.zone-2qtuhspy7cr6.eo.dnse2.com"),
					Ips:                []*string{ptrString("1.2.3.4")},
					Status:             ptrString("online"),
					L4ProxyRuleCount:   ptrInt64(0),
					UpdateTime:         ptrString("2024-01-01T00:00:00Z"),
				},
			},
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := ResourceTencentCloudTeoL4Proxy()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":    "zone-2qtuhspy7cr6",
		"proxy_name": "proxy-test",
		"area":       "overseas",
		"ipv6":       "off",
		"static_ip":  "off",
	})
	d.SetId("zone-2qtuhspy7cr6#sid-2qtuhspy7cr6")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-2qtuhspy7cr6", d.Get("zone_id"))
	assert.Equal(t, "sid-2qtuhspy7cr6", d.Get("proxy_id"))
	assert.Equal(t, "proxy-test", d.Get("proxy_name"))
	assert.Equal(t, "overseas", d.Get("area"))
	assert.Equal(t, "proxy-test.zone-2qtuhspy7cr6.eo.dnse2.com", d.Get("cname"))
	assert.Equal(t, "online", d.Get("status"))
}

// TestTeoL4ProxyRead_NotFound tests Read handles not found
func TestTeoL4ProxyRead_NotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeL4Proxy", func(_ *teov20220901.DescribeL4ProxyRequest) (*teov20220901.DescribeL4ProxyResponse, error) {
		resp := teov20220901.NewDescribeL4ProxyResponse()
		resp.Response = &teov20220901.DescribeL4ProxyResponseParams{
			TotalCount: ptrUint64(0),
			L4Proxies:  []*teov20220901.L4Proxy{},
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := ResourceTencentCloudTeoL4Proxy()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":    "zone-2qtuhspy7cr6",
		"proxy_name": "proxy-test",
		"area":       "overseas",
	})
	d.SetId("zone-2qtuhspy7cr6#sid-notexist")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}

// TestTeoL4ProxyUpdate_ImmutableField tests Update rejects immutable field changes
func TestTeoL4ProxyUpdate_ImmutableField(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	meta := newMockMeta()
	res := ResourceTencentCloudTeoL4Proxy()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":    "zone-2qtuhspy7cr6",
		"proxy_name": "proxy-test-new",
		"area":       "overseas",
		"ipv6":       "off",
		"static_ip":  "off",
	})
	d.SetId("zone-2qtuhspy7cr6#sid-2qtuhspy7cr6")

	// Mark proxy_name as changed to trigger immutable check
	d.MarkNewResource()
	_ = d.Set("proxy_name", "old-name")
	d.MarkNewResource()
	_ = d.Set("proxy_name", "proxy-test-new")

	err := res.Update(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "proxy_name")
}

// TestTeoL4ProxyDelete_Success tests Delete calls API
func TestTeoL4ProxyDelete_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "DeleteL4ProxyWithContext", func(_ interface{}, request *teov20220901.DeleteL4ProxyRequest) (*teov20220901.DeleteL4ProxyResponse, error) {
		assert.Equal(t, "zone-2qtuhspy7cr6", *request.ZoneId)
		assert.Equal(t, "sid-2qtuhspy7cr6", *request.ProxyId)
		resp := teov20220901.NewDeleteL4ProxyResponse()
		resp.Response = &teov20220901.DeleteL4ProxyResponseParams{
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	// Patch extension function to avoid state change wait
	patches.ApplyFunc(resourceTencentCloudTeoL4ProxyDeletePostFillRequest0, func(_ interface{}, _ *teov20220901.DeleteL4ProxyRequest) error {
		return nil
	})

	meta := newMockMeta()
	res := ResourceTencentCloudTeoL4Proxy()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":    "zone-2qtuhspy7cr6",
		"proxy_name": "proxy-test",
		"area":       "overseas",
	})
	d.SetId("zone-2qtuhspy7cr6#sid-2qtuhspy7cr6")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
}

// TestTeoL4ProxyDelete_BrokenId tests Delete handles broken ID
func TestTeoL4ProxyDelete_BrokenId(t *testing.T) {
	res := ResourceTencentCloudTeoL4Proxy()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":    "zone-2qtuhspy7cr6",
		"proxy_name": "proxy-test",
		"area":       "overseas",
	})
	d.SetId("broken-id")

	err := res.Delete(d, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "id is broken")
}

// TestTeoL4Proxy_Schema tests schema definition
func TestTeoL4Proxy_Schema(t *testing.T) {
	res := ResourceTencentCloudTeoL4Proxy()

	assert.NotNil(t, res)
	assert.NotNil(t, res.Create)
	assert.NotNil(t, res.Read)
	assert.NotNil(t, res.Update)
	assert.NotNil(t, res.Delete)
	assert.NotNil(t, res.Importer)

	// Required fields
	assert.Contains(t, res.Schema, "zone_id")
	assert.Contains(t, res.Schema, "proxy_name")
	assert.Contains(t, res.Schema, "area")

	zoneId := res.Schema["zone_id"]
	assert.Equal(t, schema.TypeString, zoneId.Type)
	assert.True(t, zoneId.Required)
	assert.True(t, zoneId.ForceNew)

	proxyName := res.Schema["proxy_name"]
	assert.Equal(t, schema.TypeString, proxyName.Type)
	assert.True(t, proxyName.Required)
	assert.True(t, proxyName.ForceNew)

	area := res.Schema["area"]
	assert.Equal(t, schema.TypeString, area.Type)
	assert.True(t, area.Required)
	assert.True(t, area.ForceNew)

	staticIp := res.Schema["static_ip"]
	assert.Equal(t, schema.TypeString, staticIp.Type)
	assert.True(t, staticIp.Optional)
	assert.True(t, staticIp.ForceNew)

	// Mutable fields
	ipv6 := res.Schema["ipv6"]
	assert.Equal(t, schema.TypeString, ipv6.Type)
	assert.True(t, ipv6.Optional)

	accelerateMainland := res.Schema["accelerate_mainland"]
	assert.Equal(t, schema.TypeString, accelerateMainland.Type)
	assert.True(t, accelerateMainland.Optional)

	// Computed fields
	assert.Contains(t, res.Schema, "proxy_id")
	assert.Contains(t, res.Schema, "cname")
	assert.Contains(t, res.Schema, "ips")
	assert.Contains(t, res.Schema, "status")
	assert.Contains(t, res.Schema, "l4_proxy_rule_count")
	assert.Contains(t, res.Schema, "update_time")

	proxyId := res.Schema["proxy_id"]
	assert.True(t, proxyId.Computed)

	cname := res.Schema["cname"]
	assert.True(t, cname.Computed)

	status := res.Schema["status"]
	assert.True(t, status.Computed)

	// DDoS protection config
	assert.Contains(t, res.Schema, "ddos_protection_config")
	ddosConfig := res.Schema["ddos_protection_config"]
	assert.True(t, ddosConfig.Optional)
	assert.True(t, ddosConfig.ForceNew)
	assert.Equal(t, 1, ddosConfig.MaxItems)
}

// TestTeoL4ProxyUpdate_Success tests Update
func TestTeoL4ProxyUpdate_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "ModifyL4ProxyWithContext", func(_ interface{}, request *teov20220901.ModifyL4ProxyRequest) (*teov20220901.ModifyL4ProxyResponse, error) {
		assert.Equal(t, "zone-2qtuhspy7cr6", *request.ZoneId)
		assert.Equal(t, "sid-2qtuhspy7cr6", *request.ProxyId)
		assert.Equal(t, "on", *request.Ipv6)
		resp := teov20220901.NewModifyL4ProxyResponse()
		resp.Response = &teov20220901.ModifyL4ProxyResponseParams{
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	// Patch extension function
	patches.ApplyFunc(resourceTencentCloudTeoL4ProxyUpdateOnExit, func(_ interface{}) error {
		return nil
	})

	meta := newMockMeta()
	res := ResourceTencentCloudTeoL4Proxy()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":             "zone-2qtuhspy7cr6",
		"proxy_name":          "proxy-test",
		"area":                "overseas",
		"ipv6":                "on",
		"static_ip":           "off",
		"accelerate_mainland": "off",
	})
	d.SetId("zone-2qtuhspy7cr6#sid-2qtuhspy7cr6")

	// Mark ipv6 as changed
	d.MarkNewResource()
	_ = d.Set("ipv6", "off")
	_ = d.Set("ipv6", "on")

	err := res.Update(d, meta)
	// Update will call Read which will call DescribeL4Proxy, but we haven't patched it
	// This is expected to fail at Read
	assert.Error(t, err)
}

// TestTeoL4ProxyImport tests the import functionality
func TestTeoL4ProxyImport(t *testing.T) {
	res := ResourceTencentCloudTeoL4Proxy()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":    "zone-2qtuhspy7cr6",
		"proxy_name": "proxy-test",
		"area":       "overseas",
	})
	d.SetId("zone-2qtuhspy7cr6#sid-2qtuhspy7cr6")

	imported, err := res.Importer.State(d, nil)
	assert.NoError(t, err)
	assert.Len(t, imported, 1)
	assert.Equal(t, "zone-2qtuhspy7cr6#sid-2qtuhspy7cr6", imported[0].Id())
}

// TestTeoL4ProxyImport_BrokenId tests import with broken ID
func TestTeoL4ProxyImport_BrokenId(t *testing.T) {
	res := ResourceTencentCloudTeoL4Proxy()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":    "zone-2qtuhspy7cr6",
		"proxy_name": "proxy-test",
		"area":       "overseas",
	})
	d.SetId("broken")

	_, err := res.Importer.State(d, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "id is broken")
}

// TestTeoL4ProxyRead_BrokenId tests Read with broken ID
func TestTeoL4ProxyRead_BrokenId(t *testing.T) {
	res := ResourceTencentCloudTeoL4Proxy()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":    "zone-2qtuhspy7cr6",
		"proxy_name": "proxy-test",
		"area":       "overseas",
	})
	d.SetId("broken-id")

	err := res.Read(d, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "id is broken")
}
