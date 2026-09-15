package teo_test

import (
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	svcteo "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/teo"
)

func ptrStringL4Proxy1(s string) *string {
	return &s
}

func ptrInt64L4Proxy1(i int64) *int64 {
	return &i
}

// go test ./tencentcloud/services/teo/ -run "TestTeoL4Proxy1_" -v -count=1 -gcflags="all=-l"

// TestTeoL4Proxy1_Schema tests schema definition
func TestTeoL4Proxy1_Schema(t *testing.T) {
	res := svcteo.ResourceTencentCloudTeoL4Proxy1()

	assert.NotNil(t, res)

	// proxy_id is computed only
	assert.Contains(t, res.Schema, "proxy_id")
	proxyIdSchema := res.Schema["proxy_id"]
	assert.Equal(t, schema.TypeString, proxyIdSchema.Type)
	assert.True(t, proxyIdSchema.Computed)
	assert.False(t, proxyIdSchema.Optional)
	assert.False(t, proxyIdSchema.Required)

	// zone_id is required and ForceNew
	zoneIdSchema := res.Schema["zone_id"]
	assert.True(t, zoneIdSchema.Required)
	assert.True(t, zoneIdSchema.ForceNew)

	// proxy_name is Required and ForceNew
	proxyNameSchema := res.Schema["proxy_name"]
	assert.True(t, proxyNameSchema.Required)
	assert.True(t, proxyNameSchema.ForceNew)

	// area is Required and ForceNew
	areaSchema := res.Schema["area"]
	assert.True(t, areaSchema.Required)
	assert.True(t, areaSchema.ForceNew)

	// static_ip is Optional and ForceNew
	staticIpSchema := res.Schema["static_ip"]
	assert.True(t, staticIpSchema.ForceNew)
	assert.True(t, staticIpSchema.Optional)

	// ipv6 is Optional and NOT ForceNew (mutable)
	ipv6Schema := res.Schema["ipv6"]
	assert.False(t, ipv6Schema.ForceNew)
	assert.True(t, ipv6Schema.Optional)

	// accelerate_mainland is Optional and NOT ForceNew (mutable)
	accelSchema := res.Schema["accelerate_mainland"]
	assert.False(t, accelSchema.ForceNew)
	assert.True(t, accelSchema.Optional)

	// cname is computed only
	cnameSchema := res.Schema["cname"]
	assert.True(t, cnameSchema.Computed)
	assert.False(t, cnameSchema.Optional)

	// ips is computed only
	ipsSchema := res.Schema["ips"]
	assert.True(t, ipsSchema.Computed)

	// status is computed only
	statusSchema := res.Schema["status"]
	assert.True(t, statusSchema.Computed)

	// l4proxy_rule_count is computed only
	countSchema := res.Schema["l4proxy_rule_count"]
	assert.True(t, countSchema.Computed)

	// update_time is computed only
	updateTimeSchema := res.Schema["update_time"]
	assert.True(t, updateTimeSchema.Computed)

	// ddos_protection_config is deprecated and ForceNew
	ddosSchema := res.Schema["ddos_protection_config"]
	assert.True(t, ddosSchema.Deprecated != "")
	assert.True(t, ddosSchema.ForceNew)

	// Importer exists
	assert.NotNil(t, res.Importer)

	// CRUD functions exist
	assert.NotNil(t, res.Create)
	assert.NotNil(t, res.Read)
	assert.NotNil(t, res.Update)
	assert.NotNil(t, res.Delete)
}

// TestTeoL4Proxy1_Read tests that read populates fields correctly
func TestTeoL4Proxy1_Read(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	patches.ApplyMethodFunc(&svcteo.TeoService{}, "DescribeTeoL4ProxyById", func(_ interface{}, zoneId string, proxyId string) (*teov20220901.L4Proxy, error) {
		assert.Equal(t, "zone-test1234", zoneId)
		assert.Equal(t, "proxy-87654321", proxyId)
		return &teov20220901.L4Proxy{
			ZoneId:             ptrStringL4Proxy1("zone-test1234"),
			ProxyId:            ptrStringL4Proxy1("proxy-87654321"),
			ProxyName:          ptrStringL4Proxy1("proxy-test"),
			Area:               ptrStringL4Proxy1("overseas"),
			Cname:              ptrStringL4Proxy1("proxy-test.edgeone.com"),
			Ips:                []*string{ptrStringL4Proxy1("1.2.3.4"), ptrStringL4Proxy1("5.6.7.8")},
			Status:             ptrStringL4Proxy1("online"),
			Ipv6:               ptrStringL4Proxy1("off"),
			StaticIp:           ptrStringL4Proxy1("off"),
			AccelerateMainland: ptrStringL4Proxy1("off"),
			L4ProxyRuleCount:   ptrInt64L4Proxy1(3),
			UpdateTime:         ptrStringL4Proxy1("2024-01-01T00:00:00Z"),
		}, nil
	})

	res := svcteo.ResourceTencentCloudTeoL4Proxy1()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":    "zone-test1234",
		"proxy_name": "proxy-test",
	})
	d.SetId("zone-test1234#proxy-87654321")

	meta := newMockMetaL4Proxy()
	err := res.Read(d, meta)
	assert.NoError(t, err)

	proxyId := d.Get("proxy_id").(string)
	assert.Equal(t, "proxy-87654321", proxyId)
	assert.Equal(t, "zone-test1234", d.Get("zone_id").(string))
	assert.Equal(t, "proxy-test", d.Get("proxy_name").(string))
	assert.Equal(t, "overseas", d.Get("area").(string))
	assert.Equal(t, "proxy-test.edgeone.com", d.Get("cname").(string))
	assert.Equal(t, "online", d.Get("status").(string))
	assert.Equal(t, "off", d.Get("ipv6").(string))
	assert.Equal(t, int64(3), d.Get("l4proxy_rule_count").(int64))
}

// TestTeoL4Proxy1_ReadNotFound tests read when resource is not found
func TestTeoL4Proxy1_ReadNotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	patches.ApplyMethodFunc(&svcteo.TeoService{}, "DescribeTeoL4ProxyById", func(_ interface{}, zoneId string, proxyId string) (*teov20220901.L4Proxy, error) {
		return nil, nil
	})

	res := svcteo.ResourceTencentCloudTeoL4Proxy1()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":    "zone-test123",
		"proxy_name": "test-proxy",
	})
	d.SetId("zone-test123#proxy-abc123")

	meta := newMockMetaL4Proxy()
	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}

// TestTeoL4Proxy1_Import tests import state
func TestTeoL4Proxy1_Import(t *testing.T) {
	res := svcteo.ResourceTencentCloudTeoL4Proxy1()

	assert.NotNil(t, res.Importer)
	d := &schema.ResourceData{}
	d.SetId("zone-test#proxy-test")

	states, err := res.Importer.State(d, nil)
	assert.NoError(t, err)
	assert.NotNil(t, states)
}

// TestTeoL4Proxy1_Create tests the create function with gomonkey mocks
func TestTeoL4Proxy1_Create(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodFunc(teoClient, "CreateL4ProxyWithContext", func(_ interface{}, _ *teov20220901.CreateL4ProxyRequest) (*teov20220901.CreateL4ProxyResponse, error) {
		resp := teov20220901.NewCreateL4ProxyResponse()
		resp.Response = &teov20220901.CreateL4ProxyResponseParams{
			ProxyId:   ptrStringL4Proxy1("proxy-12345678"),
			RequestId: ptrStringL4Proxy1("fake-request-id"),
		}
		return resp, nil
	})

	// Mock TeoService.DescribeTeoL4ProxyById for state refresh and Read
	patches.ApplyMethodFunc(&svcteo.TeoService{}, "DescribeTeoL4ProxyById", func(_ interface{}, zoneId string, proxyId string) (*teov20220901.L4Proxy, error) {
		return &teov20220901.L4Proxy{
			ZoneId:             ptrStringL4Proxy1("zone-test1234"),
			ProxyId:            ptrStringL4Proxy1("proxy-12345678"),
			ProxyName:          ptrStringL4Proxy1("proxy-test"),
			Area:               ptrStringL4Proxy1("overseas"),
			Ipv6:               ptrStringL4Proxy1("on"),
			StaticIp:           ptrStringL4Proxy1("off"),
			AccelerateMainland: ptrStringL4Proxy1("off"),
			Status:             ptrStringL4Proxy1("online"),
		}, nil
	})

	res := svcteo.ResourceTencentCloudTeoL4Proxy1()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":             "zone-test1234",
		"proxy_name":          "proxy-test",
		"area":                "overseas",
		"ipv6":                "on",
		"static_ip":           "off",
		"accelerate_mainland": "off",
	})

	meta := newMockMetaL4Proxy()
	err := res.Create(d, meta)
	assert.NoError(t, err)

	proxyId := d.Get("proxy_id").(string)
	assert.Equal(t, "proxy-12345678", proxyId)
	assert.Equal(t, "zone-test1234#proxy-12345678", d.Id())
}

// TestTeoL4Proxy1_CreateEmptyResponse tests error handling for empty response
func TestTeoL4Proxy1_CreateEmptyResponse(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodFunc(teoClient, "CreateL4ProxyWithContext", func(_ interface{}, _ *teov20220901.CreateL4ProxyRequest) (*teov20220901.CreateL4ProxyResponse, error) {
		return nil, nil
	})

	res := svcteo.ResourceTencentCloudTeoL4Proxy1()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":    "zone-test1234",
		"proxy_name": "proxy-test",
		"area":       "overseas",
	})

	meta := newMockMetaL4Proxy()
	err := res.Create(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "empty response")
}

// TestTeoL4Proxy1_CreateNilProxyId tests error handling for nil proxy_id
func TestTeoL4Proxy1_CreateNilProxyId(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodFunc(teoClient, "CreateL4ProxyWithContext", func(_ interface{}, _ *teov20220901.CreateL4ProxyRequest) (*teov20220901.CreateL4ProxyResponse, error) {
		resp := teov20220901.NewCreateL4ProxyResponse()
		resp.Response = &teov20220901.CreateL4ProxyResponseParams{
			ProxyId:   nil,
			RequestId: ptrStringL4Proxy1("fake-request-id"),
		}
		return resp, nil
	})

	res := svcteo.ResourceTencentCloudTeoL4Proxy1()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":    "zone-test1234",
		"proxy_name": "proxy-test",
		"area":       "overseas",
	})

	meta := newMockMetaL4Proxy()
	err := res.Create(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "proxy_id is nil")
}

// TestTeoL4Proxy1_Delete tests the delete function
func TestTeoL4Proxy1_Delete(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}

	// Mock DescribeTeoL4ProxyById to return nil (not found) - skip stop
	patches.ApplyMethodFunc(&svcteo.TeoService{}, "DescribeTeoL4ProxyById", func(_ interface{}, zoneId string, proxyId string) (*teov20220901.L4Proxy, error) {
		return nil, nil
	})

	// Mock DeleteL4ProxyWithContext
	patches.ApplyMethodFunc(teoClient, "DeleteL4ProxyWithContext", func(_ interface{}, _ *teov20220901.DeleteL4ProxyRequest) (*teov20220901.DeleteL4ProxyResponse, error) {
		resp := teov20220901.NewDeleteL4ProxyResponse()
		resp.Response = &teov20220901.DeleteL4ProxyResponseParams{
			RequestId: ptrStringL4Proxy1("fake-request-id"),
		}
		return resp, nil
	})

	res := svcteo.ResourceTencentCloudTeoL4Proxy1()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":    "zone-test1234",
		"proxy_name": "proxy-test",
	})
	d.SetId("zone-test1234#proxy-98765432")

	meta := newMockMetaL4Proxy()
	err := res.Delete(d, meta)
	assert.NoError(t, err)
}

// TestTeoL4Proxy1_DeleteOnline tests delete when instance is online (needs stop first)
func TestTeoL4Proxy1_DeleteOnline(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}

	// Mock DescribeTeoL4ProxyById to return online status
	patches.ApplyMethodFunc(&svcteo.TeoService{}, "DescribeTeoL4ProxyById", func(_ interface{}, zoneId string, proxyId string) (*teov20220901.L4Proxy, error) {
		return &teov20220901.L4Proxy{
			ZoneId:             ptrStringL4Proxy1("zone-test1234"),
			ProxyId:            ptrStringL4Proxy1("proxy-98765432"),
			ProxyName:          ptrStringL4Proxy1("proxy-test"),
			Status:             ptrStringL4Proxy1("online"),
			AccelerateMainland: ptrStringL4Proxy1("off"),
		}, nil
	})

	// Mock ModifyL4ProxyStatusWithContext for stopping
	patches.ApplyMethodFunc(teoClient, "ModifyL4ProxyStatusWithContext", func(_ interface{}, _ *teov20220901.ModifyL4ProxyStatusRequest) (*teov20220901.ModifyL4ProxyStatusResponse, error) {
		resp := teov20220901.NewModifyL4ProxyStatusResponse()
		resp.Response = &teov20220901.ModifyL4ProxyStatusResponseParams{
			RequestId: ptrStringL4Proxy1("fake-request-id"),
		}
		return resp, nil
	})

	// Mock DeleteL4ProxyWithContext
	patches.ApplyMethodFunc(teoClient, "DeleteL4ProxyWithContext", func(_ interface{}, _ *teov20220901.DeleteL4ProxyRequest) (*teov20220901.DeleteL4ProxyResponse, error) {
		resp := teov20220901.NewDeleteL4ProxyResponse()
		resp.Response = &teov20220901.DeleteL4ProxyResponseParams{
			RequestId: ptrStringL4Proxy1("fake-request-id"),
		}
		return resp, nil
	})

	res := svcteo.ResourceTencentCloudTeoL4Proxy1()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":    "zone-test1234",
		"proxy_name": "proxy-test",
	})
	d.SetId("zone-test1234#proxy-98765432")

	meta := newMockMetaL4Proxy()
	err := res.Delete(d, meta)
	assert.NoError(t, err)
}

// TestTeoL4Proxy1_Update tests the update function
func TestTeoL4Proxy1_Update(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	// Mock DescribeTeoL4ProxyById for the read after update
	patches.ApplyMethodFunc(&svcteo.TeoService{}, "DescribeTeoL4ProxyById", func(_ interface{}, zoneId string, proxyId string) (*teov20220901.L4Proxy, error) {
		return &teov20220901.L4Proxy{
			ZoneId:             ptrStringL4Proxy1("zone-test1234"),
			ProxyId:            ptrStringL4Proxy1("proxy-12345678"),
			ProxyName:          ptrStringL4Proxy1("proxy-test"),
			Area:               ptrStringL4Proxy1("overseas"),
			Ipv6:               ptrStringL4Proxy1("off"),
			StaticIp:           ptrStringL4Proxy1("off"),
			AccelerateMainland: ptrStringL4Proxy1("off"),
			Status:             ptrStringL4Proxy1("online"),
		}, nil
	})

	res := svcteo.ResourceTencentCloudTeoL4Proxy1()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":             "zone-test1234",
		"proxy_name":          "proxy-test",
		"area":                "overseas",
		"ipv6":                "off",
		"static_ip":           "off",
		"accelerate_mainland": "off",
	})
	d.SetId("zone-test1234#proxy-12345678")

	meta := newMockMetaL4Proxy()
	err := res.Update(d, meta)
	assert.NoError(t, err)
}
