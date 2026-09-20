package teo_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	svcteo "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/teo"

	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// mockMetaL4Proxy implements tccommon.ProviderMeta
type mockMetaL4Proxy struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaL4Proxy) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaL4Proxy{}

func newMockMetaL4Proxy() *mockMetaL4Proxy {
	return &mockMetaL4Proxy{client: &connectivity.TencentCloudClient{}}
}

func ptrStringL4Proxy(s string) *string {
	return &s
}

// go test -test.run TestAccTencentCloudTeoL4ProxyResource_basic -v -timeout=0
func TestAccTencentCloudTeoL4ProxyResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PRIVATE) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckL4ProxyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTeoL4Proxy,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckL4ProxyExists("tencentcloud_teo_l4_proxy.teo_l4_proxy"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_l4_proxy.teo_l4_proxy", "zone_id"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l4_proxy.teo_l4_proxy", "accelerate_mainland", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l4_proxy.teo_l4_proxy", "area", "overseas"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l4_proxy.teo_l4_proxy", "ipv6", "on"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l4_proxy.teo_l4_proxy", "proxy_name", "proxy-test"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l4_proxy.teo_l4_proxy", "static_ip", "off"),
				),
			},
			{
				ResourceName:      "tencentcloud_teo_l4_proxy.teo_l4_proxy",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccTeoL4ProxyUp,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckL4ProxyExists("tencentcloud_teo_l4_proxy.teo_l4_proxy"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_l4_proxy.teo_l4_proxy", "zone_id"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l4_proxy.teo_l4_proxy", "accelerate_mainland", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l4_proxy.teo_l4_proxy", "area", "overseas"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l4_proxy.teo_l4_proxy", "ipv6", "off"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l4_proxy.teo_l4_proxy", "proxy_name", "proxy-test"),
					resource.TestCheckResourceAttr("tencentcloud_teo_l4_proxy.teo_l4_proxy", "static_ip", "off"),
				),
			},
		},
	})
}

func testAccCheckL4ProxyDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	service := svcteo.NewTeoService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_teo_l4_proxy" {
			continue
		}
		idSplit := strings.Split(rs.Primary.ID, tccommon.FILED_SP)
		if len(idSplit) != 2 {
			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
		}
		zoneId := idSplit[0]
		proxyId := idSplit[1]

		proxy, _, err := service.DescribeTeoL4ProxyById(ctx, zoneId, proxyId, nil, nil)
		if proxy != nil {
			return fmt.Errorf("zone l4 proxy %s still exists", rs.Primary.ID)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func testAccCheckL4ProxyExists(r string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}
		idSplit := strings.Split(rs.Primary.ID, tccommon.FILED_SP)
		if len(idSplit) != 2 {
			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
		}
		zoneId := idSplit[0]
		proxyId := idSplit[1]

		service := svcteo.NewTeoService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
		proxy, _, err := service.DescribeTeoL4ProxyById(ctx, zoneId, proxyId, nil, nil)
		if proxy == nil {
			return fmt.Errorf("zone l4 proxy %s is not found", rs.Primary.ID)
		}
		if err != nil {
			return err
		}

		return nil
	}
}

const testAccTeoL4Proxy = `

resource "tencentcloud_teo_l4_proxy" "teo_l4_proxy" {
  accelerate_mainland = "off"
  area                = "overseas"
  ipv6                = "on"
  proxy_name          = "proxy-test"
  static_ip           = "off"
  zone_id             = "zone-2qtuhspy7cr6"
}
`

const testAccTeoL4ProxyUp = `

resource "tencentcloud_teo_l4_proxy" "teo_l4_proxy" {
  accelerate_mainland = "off"
  area                = "overseas"
  ipv6                = "off"
  proxy_name          = "proxy-test"
  static_ip           = "off"
  zone_id             = "zone-2qtuhspy7cr6"
}
`

// ---- Unit Tests (gomonkey mock) ----

// go test ./tencentcloud/services/teo/ -run "TestTeoL4Proxy_" -v -count=1 -gcflags="all=-l"

// TestTeoL4Proxy_Create tests that proxy_id is set correctly after Create
func TestTeoL4Proxy_Create(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaL4Proxy().client, "UseTeoClient", teoClient)

	// Mock CreateL4ProxyWithContext to return a response with ProxyId
	patches.ApplyMethodFunc(teoClient, "CreateL4ProxyWithContext", func(_ context.Context, _ *teov20220901.CreateL4ProxyRequest) (*teov20220901.CreateL4ProxyResponse, error) {
		resp := teov20220901.NewCreateL4ProxyResponse()
		resp.Response = &teov20220901.CreateL4ProxyResponseParams{
			ProxyId:   ptrStringL4Proxy("proxy-12345678"),
			RequestId: ptrStringL4Proxy("fake-request-id"),
		}
		return resp, nil
	})

	// Mock TeoService.DescribeTeoL4ProxyById for both the state refresh in CreatePostHandleResponse and the Read call
	patches.ApplyMethodFunc(&svcteo.TeoService{}, "DescribeTeoL4ProxyById", func(_ context.Context, zoneId string, proxyId string, offset *uint64, limit *uint64) (*teov20220901.L4Proxy, *uint64, error) {
		return &teov20220901.L4Proxy{
			ZoneId:             ptrStringL4Proxy("zone-test1234"),
			ProxyId:            ptrStringL4Proxy("proxy-12345678"),
			ProxyName:          ptrStringL4Proxy("proxy-test"),
			Area:               ptrStringL4Proxy("overseas"),
			Ipv6:               ptrStringL4Proxy("on"),
			StaticIp:           ptrStringL4Proxy("off"),
			AccelerateMainland: ptrStringL4Proxy("off"),
			Status:             ptrStringL4Proxy("online"),
		}, nil, nil
	})

	meta := newMockMetaL4Proxy()
	res := svcteo.ResourceTencentCloudTeoL4Proxy()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":             "zone-test1234",
		"proxy_name":          "proxy-test",
		"area":                "overseas",
		"ipv6":                "on",
		"static_ip":           "off",
		"accelerate_mainland": "off",
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)

	// Verify proxy_id is set correctly
	proxyId := d.Get("proxy_id").(string)
	assert.Equal(t, "proxy-12345678", proxyId)

	// Verify composite ID
	assert.Equal(t, "zone-test1234#proxy-12345678", d.Id())
}

// TestTeoL4Proxy_Read tests that proxy_id is set correctly after Read
func TestTeoL4Proxy_Read(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	// Mock TeoService.DescribeTeoL4ProxyById for the Read flow
	patches.ApplyMethodFunc(&svcteo.TeoService{}, "DescribeTeoL4ProxyById", func(_ context.Context, zoneId string, proxyId string, offset *uint64, limit *uint64) (*teov20220901.L4Proxy, *uint64, error) {
		assert.Equal(t, "zone-test1234", zoneId)
		assert.Equal(t, "proxy-87654321", proxyId)
		return &teov20220901.L4Proxy{
			ZoneId:             ptrStringL4Proxy("zone-test1234"),
			ProxyId:            ptrStringL4Proxy("proxy-87654321"),
			ProxyName:          ptrStringL4Proxy("proxy-test"),
			Area:               ptrStringL4Proxy("overseas"),
			Ipv6:               ptrStringL4Proxy("off"),
			StaticIp:           ptrStringL4Proxy("off"),
			AccelerateMainland: ptrStringL4Proxy("off"),
		}, nil, nil
	})

	meta := newMockMetaL4Proxy()
	res := svcteo.ResourceTencentCloudTeoL4Proxy()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":    "zone-test1234",
		"proxy_name": "proxy-test",
	})
	d.SetId("zone-test1234#proxy-87654321")

	err := res.Read(d, meta)
	assert.NoError(t, err)

	// Verify proxy_id is set correctly from Read
	proxyId := d.Get("proxy_id").(string)
	assert.Equal(t, "proxy-87654321", proxyId)
}

// TestTeoL4Proxy_Read_NotFound tests read when resource is not found
func TestTeoL4Proxy_Read_NotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	// Mock TeoService.DescribeTeoL4ProxyById to return nil (not found)
	patches.ApplyMethodFunc(&svcteo.TeoService{}, "DescribeTeoL4ProxyById", func(_ context.Context, zoneId string, proxyId string, offset *uint64, limit *uint64) (*teov20220901.L4Proxy, *uint64, error) {
		return nil, nil, nil
	})

	meta := newMockMetaL4Proxy()
	res := svcteo.ResourceTencentCloudTeoL4Proxy()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":    "zone-test123",
		"proxy_name": "test-proxy",
	})
	d.SetId("zone-test123#proxy-abc123")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}

// TestTeoL4Proxy_Schema tests that proxy_id is defined as computed attribute in schema
func TestTeoL4Proxy_Schema(t *testing.T) {
	res := svcteo.ResourceTencentCloudTeoL4Proxy()

	assert.NotNil(t, res)

	assert.Contains(t, res.Schema, "proxy_id")
	proxyIdSchema := res.Schema["proxy_id"]
	assert.Equal(t, schema.TypeString, proxyIdSchema.Type)
	assert.True(t, proxyIdSchema.Computed)
	assert.False(t, proxyIdSchema.Optional)
	assert.False(t, proxyIdSchema.Required)
}

// TestTeoL4Proxy_Read_OffsetLimit tests that offset and limit are passed to DescribeTeoL4ProxyById and total_count is set
func TestTeoL4Proxy_Read_OffsetLimit(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	var capturedOffset *uint64
	var capturedLimit *uint64

	// Mock TeoService.DescribeTeoL4ProxyById to capture offset/limit and return total_count
	patches.ApplyMethodFunc(&svcteo.TeoService{}, "DescribeTeoL4ProxyById", func(_ context.Context, zoneId string, proxyId string, offset *uint64, limit *uint64) (*teov20220901.L4Proxy, *uint64, error) {
		capturedOffset = offset
		capturedLimit = limit
		totalCount := uint64(10)
		return &teov20220901.L4Proxy{
			ZoneId:             ptrStringL4Proxy("zone-test1234"),
			ProxyId:            ptrStringL4Proxy("proxy-87654321"),
			ProxyName:          ptrStringL4Proxy("proxy-test"),
			Area:               ptrStringL4Proxy("overseas"),
			Ipv6:               ptrStringL4Proxy("off"),
			StaticIp:           ptrStringL4Proxy("off"),
			AccelerateMainland: ptrStringL4Proxy("off"),
		}, &totalCount, nil
	})

	meta := newMockMetaL4Proxy()
	res := svcteo.ResourceTencentCloudTeoL4Proxy()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":    "zone-test1234",
		"proxy_name": "proxy-test",
		"offset":     5,
		"limit":      10,
	})
	d.SetId("zone-test1234#proxy-87654321")

	err := res.Read(d, meta)
	assert.NoError(t, err)

	// Verify offset and limit were passed to the service
	assert.NotNil(t, capturedOffset)
	assert.Equal(t, uint64(5), *capturedOffset)
	assert.NotNil(t, capturedLimit)
	assert.Equal(t, uint64(10), *capturedLimit)

	// Verify total_count is set
	totalCount := d.Get("total_count").(int)
	assert.Equal(t, 10, totalCount)
}

// TestTeoL4Proxy_Read_TotalCount tests that total_count is set when returned from API
func TestTeoL4Proxy_Read_TotalCount(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	totalCount := uint64(5)

	// Mock TeoService.DescribeTeoL4ProxyById to return a specific total_count
	patches.ApplyMethodFunc(&svcteo.TeoService{}, "DescribeTeoL4ProxyById", func(_ context.Context, zoneId string, proxyId string, offset *uint64, limit *uint64) (*teov20220901.L4Proxy, *uint64, error) {
		return &teov20220901.L4Proxy{
			ZoneId:             ptrStringL4Proxy("zone-test1234"),
			ProxyId:            ptrStringL4Proxy("proxy-87654321"),
			ProxyName:          ptrStringL4Proxy("proxy-test"),
			Area:               ptrStringL4Proxy("overseas"),
			Ipv6:               ptrStringL4Proxy("off"),
			StaticIp:           ptrStringL4Proxy("off"),
			AccelerateMainland: ptrStringL4Proxy("off"),
		}, &totalCount, nil
	})

	meta := newMockMetaL4Proxy()
	res := svcteo.ResourceTencentCloudTeoL4Proxy()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":    "zone-test1234",
		"proxy_name": "proxy-test",
	})
	d.SetId("zone-test1234#proxy-87654321")

	err := res.Read(d, meta)
	assert.NoError(t, err)

	// Verify total_count is set correctly
	tc := d.Get("total_count").(int)
	assert.Equal(t, 5, tc)
}

// TestTeoL4Proxy_Schema_NewFields tests that offset, limit, total_count are defined in schema
func TestTeoL4Proxy_Schema_NewFields(t *testing.T) {
	res := svcteo.ResourceTencentCloudTeoL4Proxy()

	assert.NotNil(t, res)

	// offset
	assert.Contains(t, res.Schema, "offset")
	offsetSchema := res.Schema["offset"]
	assert.Equal(t, schema.TypeInt, offsetSchema.Type)
	assert.True(t, offsetSchema.Optional)
	assert.False(t, offsetSchema.Required)
	assert.False(t, offsetSchema.Computed)

	// limit
	assert.Contains(t, res.Schema, "limit")
	limitSchema := res.Schema["limit"]
	assert.Equal(t, schema.TypeInt, limitSchema.Type)
	assert.True(t, limitSchema.Optional)
	assert.False(t, limitSchema.Required)
	assert.False(t, limitSchema.Computed)

	// total_count
	assert.Contains(t, res.Schema, "total_count")
	totalCountSchema := res.Schema["total_count"]
	assert.Equal(t, schema.TypeInt, totalCountSchema.Type)
	assert.True(t, totalCountSchema.Computed)
	assert.False(t, totalCountSchema.Optional)
	assert.False(t, totalCountSchema.Required)
}
