package teo_test

import (
	"fmt"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/teo"
)

// mockMetaCertCfg implements tccommon.ProviderMeta
type mockMetaCertCfg struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaCertCfg) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaCertCfg{}

func newMockMetaCertCfg() *mockMetaCertCfg {
	return &mockMetaCertCfg{client: &connectivity.TencentCloudClient{}}
}

func ptrStrCfg(s string) *string {
	return &s
}

func ptrInt64Cfg(i int64) *int64 {
	return &i
}

// go test ./tencentcloud/services/teo/ -run "TestTeoCertificateConfigClientCertInfo" -v -count=1 -gcflags="all=-l"

// buildAccelerationDomainCertCfgResp builds a DescribeAccelerationDomains response with the given ClientCertInfo
func buildAccelerationDomainCertCfgResp(clientCertInfo *teov20220901.MutualTLS, domainStatus *string) *teov20220901.DescribeAccelerationDomainsResponse {
	resp := teov20220901.NewDescribeAccelerationDomainsResponse()
	resp.Response = &teov20220901.DescribeAccelerationDomainsResponseParams{
		TotalCount: ptrInt64Cfg(1),
		AccelerationDomains: []*teov20220901.AccelerationDomain{
			{
				ZoneId:       ptrStrCfg("zone-1234567890"),
				DomainName:   ptrStrCfg("test.example.com"),
				DomainStatus: domainStatus,
				Certificate: &teov20220901.AccelerationDomainCertificate{
					Mode: ptrStrCfg("sslcert"),
					List: []*teov20220901.CertificateInfo{
						{
							CertId: ptrStrCfg("8xiUJIJd"),
						},
					},
					ClientCertInfo: clientCertInfo,
				},
			},
		},
		RequestId: ptrStrCfg("fake-request-id"),
	}
	return resp
}

// buildDescribeZonesCertCfgResp builds a DescribeZones response
func buildDescribeZonesCertCfgResp() *teov20220901.DescribeZonesResponse {
	resp := teov20220901.NewDescribeZonesResponse()
	resp.Response = &teov20220901.DescribeZonesResponseParams{
		TotalCount: ptrInt64Cfg(1),
		Zones: []*teov20220901.Zone{
			{
				ZoneId:   ptrStrCfg("zone-1234567890"),
				ZoneName: ptrStrCfg("example.com"),
			},
		},
		RequestId: ptrStrCfg("fake-request-id"),
	}
	return resp
}

// TestTeoCertificateConfigClientCertInfo_CreateWithClientCert tests Create with client_cert_info specified
func TestTeoCertificateConfigClientCertInfo_CreateWithClientCert(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaCertCfg().client, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "ModifyHostsCertificate", func(request *teov20220901.ModifyHostsCertificateRequest) (*teov20220901.ModifyHostsCertificateResponse, error) {
		assert.NotNil(t, request.ClientCertInfo)
		assert.Equal(t, "on", *request.ClientCertInfo.Switch)
		assert.Len(t, request.ClientCertInfo.CertInfos, 1)
		assert.Equal(t, "cert-abc123", *request.ClientCertInfo.CertInfos[0].CertId)

		resp := teov20220901.NewModifyHostsCertificateResponse()
		resp.Response = &teov20220901.ModifyHostsCertificateResponseParams{
			RequestId: ptrStrCfg("fake-request-id"),
		}
		return resp, nil
	})

	clientCertInfo := &teov20220901.MutualTLS{
		Switch: ptrStrCfg("on"),
		CertInfos: []*teov20220901.CertificateInfo{
			{
				CertId:     ptrStrCfg("cert-abc123"),
				Alias:      ptrStrCfg("test-cert"),
				Type:       ptrStrCfg("upload"),
				ExpireTime: ptrStrCfg("2025-12-31T23:59:59Z"),
				DeployTime: ptrStrCfg("2024-06-01T10:00:00Z"),
				SignAlgo:   ptrStrCfg("RSA 2048"),
				Status:     ptrStrCfg("deployed"),
			},
		},
	}

	patches.ApplyMethodFunc(teoClient, "DescribeAccelerationDomains", func(request *teov20220901.DescribeAccelerationDomainsRequest) (*teov20220901.DescribeAccelerationDomainsResponse, error) {
		return buildAccelerationDomainCertCfgResp(clientCertInfo, ptrStrCfg("online")), nil
	})

	patches.ApplyMethodFunc(teoClient, "DescribeZones", func(request *teov20220901.DescribeZonesRequest) (*teov20220901.DescribeZonesResponse, error) {
		return buildDescribeZonesCertCfgResp(), nil
	})

	meta := newMockMetaCertCfg()
	res := teo.ResourceTencentCloudTeoCertificateConfig()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-1234567890",
		"host":    "test.example.com",
		"mode":    "sslcert",
		"server_cert_info": []interface{}{
			map[string]interface{}{
				"cert_id": "8xiUJIJd",
			},
		},
		"client_cert_info": []interface{}{
			map[string]interface{}{
				"switch": "on",
				"cert_infos": []interface{}{
					map[string]interface{}{
						"cert_id": "cert-abc123",
					},
				},
			},
		},
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-1234567890#test.example.com", d.Id())
}

// TestTeoCertificateConfigClientCertInfo_ReadWithClientCert tests Read with ClientCertInfo in response
func TestTeoCertificateConfigClientCertInfo_ReadWithClientCert(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaCertCfg().client, "UseTeoClient", teoClient)

	clientCertInfo := &teov20220901.MutualTLS{
		Switch: ptrStrCfg("on"),
		CertInfos: []*teov20220901.CertificateInfo{
			{
				CertId:     ptrStrCfg("cert-abc123"),
				Alias:      ptrStrCfg("test-cert"),
				Type:       ptrStrCfg("upload"),
				ExpireTime: ptrStrCfg("2025-12-31T23:59:59Z"),
				DeployTime: ptrStrCfg("2024-06-01T10:00:00Z"),
				SignAlgo:   ptrStrCfg("RSA 2048"),
				Status:     ptrStrCfg("deployed"),
			},
		},
	}

	patches.ApplyMethodFunc(teoClient, "DescribeAccelerationDomains", func(request *teov20220901.DescribeAccelerationDomainsRequest) (*teov20220901.DescribeAccelerationDomainsResponse, error) {
		return buildAccelerationDomainCertCfgResp(clientCertInfo, nil), nil
	})

	patches.ApplyMethodFunc(teoClient, "DescribeZones", func(request *teov20220901.DescribeZonesRequest) (*teov20220901.DescribeZonesResponse, error) {
		return buildDescribeZonesCertCfgResp(), nil
	})

	meta := newMockMetaCertCfg()
	res := teo.ResourceTencentCloudTeoCertificateConfig()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-1234567890",
		"host":    "test.example.com",
	})
	d.SetId("zone-1234567890#test.example.com")

	err := res.Read(d, meta)
	assert.NoError(t, err)

	clientCertInfoResult := d.Get("client_cert_info").([]interface{})
	assert.Len(t, clientCertInfoResult, 1)

	clientCertInfoMap := clientCertInfoResult[0].(map[string]interface{})
	assert.Equal(t, "on", clientCertInfoMap["switch"])

	certInfos := clientCertInfoMap["cert_infos"].([]interface{})
	assert.Len(t, certInfos, 1)

	certInfosMap := certInfos[0].(map[string]interface{})
	assert.Equal(t, "cert-abc123", certInfosMap["cert_id"])
	assert.Equal(t, "test-cert", certInfosMap["alias"])
	assert.Equal(t, "upload", certInfosMap["type"])
	assert.Equal(t, "2025-12-31T23:59:59Z", certInfosMap["expire_time"])
	assert.Equal(t, "2024-06-01T10:00:00Z", certInfosMap["deploy_time"])
	assert.Equal(t, "RSA 2048", certInfosMap["sign_algo"])
	assert.Equal(t, "deployed", certInfosMap["status"])
}

// TestTeoCertificateConfigClientCertInfo_ReadNilClientCert tests Read with nil ClientCertInfo
func TestTeoCertificateConfigClientCertInfo_ReadNilClientCert(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaCertCfg().client, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeAccelerationDomains", func(request *teov20220901.DescribeAccelerationDomainsRequest) (*teov20220901.DescribeAccelerationDomainsResponse, error) {
		return buildAccelerationDomainCertCfgResp(nil, nil), nil
	})

	patches.ApplyMethodFunc(teoClient, "DescribeZones", func(request *teov20220901.DescribeZonesRequest) (*teov20220901.DescribeZonesResponse, error) {
		return buildDescribeZonesCertCfgResp(), nil
	})

	meta := newMockMetaCertCfg()
	res := teo.ResourceTencentCloudTeoCertificateConfig()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-1234567890",
		"host":    "test.example.com",
	})
	d.SetId("zone-1234567890#test.example.com")

	err := res.Read(d, meta)
	assert.NoError(t, err)

	clientCertInfo := d.Get("client_cert_info").([]interface{})
	assert.Len(t, clientCertInfo, 0)
}

// TestTeoCertificateConfigClientCertInfo_UpdateWithClientCert tests Update with client_cert_info
func TestTeoCertificateConfigClientCertInfo_UpdateWithClientCert(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaCertCfg().client, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "ModifyHostsCertificate", func(request *teov20220901.ModifyHostsCertificateRequest) (*teov20220901.ModifyHostsCertificateResponse, error) {
		assert.NotNil(t, request.ClientCertInfo)
		assert.Equal(t, "off", *request.ClientCertInfo.Switch)

		resp := teov20220901.NewModifyHostsCertificateResponse()
		resp.Response = &teov20220901.ModifyHostsCertificateResponseParams{
			RequestId: ptrStrCfg("fake-request-id"),
		}
		return resp, nil
	})

	clientCertInfoOff := &teov20220901.MutualTLS{
		Switch: ptrStrCfg("off"),
	}

	patches.ApplyMethodFunc(teoClient, "DescribeAccelerationDomains", func(request *teov20220901.DescribeAccelerationDomainsRequest) (*teov20220901.DescribeAccelerationDomainsResponse, error) {
		return buildAccelerationDomainCertCfgResp(clientCertInfoOff, ptrStrCfg("online")), nil
	})

	patches.ApplyMethodFunc(teoClient, "DescribeZones", func(request *teov20220901.DescribeZonesRequest) (*teov20220901.DescribeZonesResponse, error) {
		return buildDescribeZonesCertCfgResp(), nil
	})

	meta := newMockMetaCertCfg()
	res := teo.ResourceTencentCloudTeoCertificateConfig()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-1234567890",
		"host":    "test.example.com",
		"mode":    "sslcert",
		"server_cert_info": []interface{}{
			map[string]interface{}{
				"cert_id": "8xiUJIJd",
			},
		},
		"client_cert_info": []interface{}{
			map[string]interface{}{
				"switch": "off",
			},
		},
	})
	d.SetId("zone-1234567890#test.example.com")

	err := res.Update(d, meta)
	assert.NoError(t, err)
}

// TestTeoCertificateConfigClientCertInfo_CreateAPIError tests Create handles API error
func TestTeoCertificateConfigClientCertInfo_CreateAPIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaCertCfg().client, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "ModifyHostsCertificate", func(request *teov20220901.ModifyHostsCertificateRequest) (*teov20220901.ModifyHostsCertificateResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=InvalidParameter, Message=Invalid cert_id")
	})

	meta := newMockMetaCertCfg()
	res := teo.ResourceTencentCloudTeoCertificateConfig()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-1234567890",
		"host":    "test.example.com",
		"mode":    "sslcert",
		"server_cert_info": []interface{}{
			map[string]interface{}{
				"cert_id": "invalid-cert-id",
			},
		},
		"client_cert_info": []interface{}{
			map[string]interface{}{
				"switch": "on",
				"cert_infos": []interface{}{
					map[string]interface{}{
						"cert_id": "cert-abc123",
					},
				},
			},
		},
	})

	err := res.Create(d, meta)
	assert.Error(t, err)
}

func TestAccTencentCloudTeoCertificateConfigResource_basic(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTeoCertificateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_teo_certificate_config.certificate", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_certificate_config.certificate", "zone_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_certificate_config.certificate", "host"),
					resource.TestCheckResourceAttr("tencentcloud_teo_certificate_config.certificate", "mode", "sslcert"),
					resource.TestCheckResourceAttr("tencentcloud_teo_certificate_config.certificate", "server_cert_info.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_certificate_config.certificate", "server_cert_info.0.alias", "terraform_test"),
					resource.TestCheckResourceAttr("tencentcloud_teo_certificate_config.certificate", "server_cert_info.0.cert_id", "EEIqXrZt"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_certificate_config.certificate", "server_cert_info.0.common_name"),
					resource.TestCheckResourceAttr("tencentcloud_teo_certificate_config.certificate", "server_cert_info.0.sign_algo", "RSA 2048"),
					resource.TestCheckResourceAttr("tencentcloud_teo_certificate_config.certificate", "server_cert_info.0.type", "managed"),
				),
			},
			{
				ResourceName:      "tencentcloud_teo_certificate_config.certificate",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTeoCertificateConfig = testAccTeoZone + `

resource "tencentcloud_teo_ownership_verify" "ownership_verify" {
  domain = var.zone_name

  depends_on = [ tencentcloud_teo_zone.basic ]
}

resource "tencentcloud_teo_acceleration_domain" "acceleration_domain" {
  zone_id     = tencentcloud_teo_zone.basic.id
  domain_name = "test.tf-teo.xyz"

  origin_info {
    origin      = "150.109.8.1"
    origin_type = "IP_DOMAIN"
  }

  depends_on = [ tencentcloud_teo_ownership_verify.ownership_verify ]
}

resource "tencentcloud_teo_certificate_config" "certificate" {
  host    = format("test.%s", var.zone_name)
  mode    = "sslcert"
  zone_id = tencentcloud_teo_zone.basic.id

  server_cert_info {
    alias       = "terraform_test"
    cert_id     = "EEIqXrZt"
    common_name = var.zone_name
    //deploy_time = "2024-04-22T10:34:13Z"
    //expire_time = "2025-04-22T23:59:59Z"
    sign_algo   = "RSA 2048"
    type        = "managed"
  }

  depends_on = [tencentcloud_teo_acceleration_domain.acceleration_domain]
}

`
