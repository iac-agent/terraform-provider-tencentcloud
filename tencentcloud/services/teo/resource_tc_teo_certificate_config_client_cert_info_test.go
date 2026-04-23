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

type certConfigMockMeta struct {
	client *connectivity.TencentCloudClient
}

func (m *certConfigMockMeta) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &certConfigMockMeta{}

func newCertConfigMockMeta() *certConfigMockMeta {
	return &certConfigMockMeta{client: &connectivity.TencentCloudClient{}}
}

func ptrStringCertConfig(s string) *string {
	return &s
}

// go test ./tencentcloud/services/teo/ -run "TestTeoCertificateConfigClientCertInfo" -v -count=1 -gcflags="all=-l"

// TestTeoCertificateConfigClientCertInfo_CreateWithClientCertInfo tests creating a resource with client_cert_info
func TestTeoCertificateConfigClientCertInfo_CreateWithClientCertInfo(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newCertConfigMockMeta().client, "UseTeoClient", teoClient)

	// Mock ModifyHostsCertificate to verify ClientCertInfo is sent
	patches.ApplyMethodFunc(teoClient, "ModifyHostsCertificate", func(request *teov20220901.ModifyHostsCertificateRequest) (*teov20220901.ModifyHostsCertificateResponse, error) {
		assert.Equal(t, "zone-test123", *request.ZoneId)
		assert.NotNil(t, request.ClientCertInfo)
		assert.Equal(t, "on", *request.ClientCertInfo.Switch)
		assert.Len(t, request.ClientCertInfo.CertInfos, 1)
		assert.Equal(t, "cert-abc123", *request.ClientCertInfo.CertInfos[0].CertId)

		resp := teov20220901.NewModifyHostsCertificateResponse()
		resp.Response = &teov20220901.ModifyHostsCertificateResponseParams{
			RequestId: ptrStringCertConfig("fake-request-id"),
		}
		return resp, nil
	})

	// Mock DescribeAccelerationDomains for the read after update
	patches.ApplyMethodFunc(teoClient, "DescribeAccelerationDomains", func(request *teov20220901.DescribeAccelerationDomainsRequest) (*teov20220901.DescribeAccelerationDomainsResponse, error) {
		resp := teov20220901.NewDescribeAccelerationDomainsResponse()
		resp.Response = &teov20220901.DescribeAccelerationDomainsResponseParams{
			AccelerationDomains: []*teov20220901.AccelerationDomain{
				{
					DomainName:   ptrStringCertConfig("test.example.com"),
					DomainStatus: ptrStringCertConfig("online"),
					Certificate: &teov20220901.AccelerationDomainCertificate{
						Mode: ptrStringCertConfig("sslcert"),
						List: []*teov20220901.CertificateInfo{
							{
								CertId:   ptrStringCertConfig("cert-abc123"),
								Alias:    ptrStringCertConfig("test-cert"),
								Type:     ptrStringCertConfig("managed"),
								SignAlgo: ptrStringCertConfig("RSA 2048"),
							},
						},
						ClientCertInfo: &teov20220901.MutualTLS{
							Switch: ptrStringCertConfig("on"),
							CertInfos: []*teov20220901.CertificateInfo{
								{
									CertId:     ptrStringCertConfig("cert-abc123"),
									Alias:      ptrStringCertConfig("client-ca-cert"),
									Type:       ptrStringCertConfig("upload"),
									ExpireTime: ptrStringCertConfig("2026-01-01T00:00:00Z"),
									DeployTime: ptrStringCertConfig("2025-01-01T00:00:00Z"),
									SignAlgo:   ptrStringCertConfig("ECDSA P256"),
								},
							},
						},
					},
				},
			},
			RequestId: ptrStringCertConfig("fake-request-id"),
		}
		return resp, nil
	})

	// Mock DescribeTeoZone for the read flow (used to get ZoneName for common_name)
	patches.ApplyMethodFunc(teoClient, "DescribeZones", func(request *teov20220901.DescribeZonesRequest) (*teov20220901.DescribeZonesResponse, error) {
		resp := teov20220901.NewDescribeZonesResponse()
		resp.Response = &teov20220901.DescribeZonesResponseParams{
			Zones: []*teov20220901.Zone{
				{
					ZoneId:   ptrStringCertConfig("zone-test123"),
					ZoneName: ptrStringCertConfig("example.com"),
				},
			},
			RequestId: ptrStringCertConfig("fake-request-id"),
		}
		return resp, nil
	})

	meta := newCertConfigMockMeta()
	res := teo.ResourceTencentCloudTeoCertificateConfig()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-test123",
		"host":    "test.example.com",
		"mode":    "sslcert",
		"server_cert_info": []interface{}{
			map[string]interface{}{
				"cert_id": "cert-abc123",
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
	assert.Equal(t, "zone-test123#test.example.com", d.Id())

	// Verify client_cert_info was read into state
	clientCertInfo := d.Get("client_cert_info").([]interface{})
	assert.Len(t, clientCertInfo, 1)
	clientCertInfoMap := clientCertInfo[0].(map[string]interface{})
	assert.Equal(t, "on", clientCertInfoMap["switch"])

	certInfos := clientCertInfoMap["cert_infos"].([]interface{})
	assert.Len(t, certInfos, 1)
	certInfoMap := certInfos[0].(map[string]interface{})
	assert.Equal(t, "cert-abc123", certInfoMap["cert_id"])
	assert.Equal(t, "client-ca-cert", certInfoMap["alias"])
	assert.Equal(t, "upload", certInfoMap["type"])
	assert.Equal(t, "2026-01-01T00:00:00Z", certInfoMap["expire_time"])
	assert.Equal(t, "2025-01-01T00:00:00Z", certInfoMap["deploy_time"])
	assert.Equal(t, "ECDSA P256", certInfoMap["sign_algo"])
}

// TestTeoCertificateConfigClientCertInfo_CreateWithoutClientCertInfo tests creating a resource without client_cert_info
func TestTeoCertificateConfigClientCertInfo_CreateWithoutClientCertInfo(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newCertConfigMockMeta().client, "UseTeoClient", teoClient)

	// Mock ModifyHostsCertificate - ClientCertInfo should not be set
	patches.ApplyMethodFunc(teoClient, "ModifyHostsCertificate", func(request *teov20220901.ModifyHostsCertificateRequest) (*teov20220901.ModifyHostsCertificateResponse, error) {
		assert.Nil(t, request.ClientCertInfo)

		resp := teov20220901.NewModifyHostsCertificateResponse()
		resp.Response = &teov20220901.ModifyHostsCertificateResponseParams{
			RequestId: ptrStringCertConfig("fake-request-id"),
		}
		return resp, nil
	})

	// Mock DescribeAccelerationDomains for the read after update
	patches.ApplyMethodFunc(teoClient, "DescribeAccelerationDomains", func(request *teov20220901.DescribeAccelerationDomainsRequest) (*teov20220901.DescribeAccelerationDomainsResponse, error) {
		resp := teov20220901.NewDescribeAccelerationDomainsResponse()
		resp.Response = &teov20220901.DescribeAccelerationDomainsResponseParams{
			AccelerationDomains: []*teov20220901.AccelerationDomain{
				{
					DomainName:   ptrStringCertConfig("test.example.com"),
					DomainStatus: ptrStringCertConfig("online"),
					Certificate: &teov20220901.AccelerationDomainCertificate{
						Mode: ptrStringCertConfig("eofreecert"),
						List: []*teov20220901.CertificateInfo{},
					},
				},
			},
			RequestId: ptrStringCertConfig("fake-request-id"),
		}
		return resp, nil
	})

	// Mock DescribeZones for the read flow (used to get ZoneName for common_name)
	patches.ApplyMethodFunc(teoClient, "DescribeZones", func(request *teov20220901.DescribeZonesRequest) (*teov20220901.DescribeZonesResponse, error) {
		resp := teov20220901.NewDescribeZonesResponse()
		resp.Response = &teov20220901.DescribeZonesResponseParams{
			Zones: []*teov20220901.Zone{
				{
					ZoneId:   ptrStringCertConfig("zone-test123"),
					ZoneName: ptrStringCertConfig("example.com"),
				},
			},
			RequestId: ptrStringCertConfig("fake-request-id"),
		}
		return resp, nil
	})

	meta := newCertConfigMockMeta()
	res := teo.ResourceTencentCloudTeoCertificateConfig()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-test123",
		"host":    "test.example.com",
		"mode":    "eofreecert",
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-test123#test.example.com", d.Id())
}

// TestTeoCertificateConfigClientCertInfo_ReadWithClientCertInfo tests reading client_cert_info from API response
func TestTeoCertificateConfigClientCertInfo_ReadWithClientCertInfo(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newCertConfigMockMeta().client, "UseTeoClient", teoClient)

	// Mock DescribeAccelerationDomains for the read flow
	patches.ApplyMethodFunc(teoClient, "DescribeAccelerationDomains", func(request *teov20220901.DescribeAccelerationDomainsRequest) (*teov20220901.DescribeAccelerationDomainsResponse, error) {
		resp := teov20220901.NewDescribeAccelerationDomainsResponse()
		resp.Response = &teov20220901.DescribeAccelerationDomainsResponseParams{
			AccelerationDomains: []*teov20220901.AccelerationDomain{
				{
					DomainName:   ptrStringCertConfig("test.example.com"),
					DomainStatus: ptrStringCertConfig("online"),
					Certificate: &teov20220901.AccelerationDomainCertificate{
						Mode: ptrStringCertConfig("sslcert"),
						List: []*teov20220901.CertificateInfo{
							{
								CertId:   ptrStringCertConfig("server-cert-001"),
								Alias:    ptrStringCertConfig("server-cert"),
								Type:     ptrStringCertConfig("managed"),
								SignAlgo: ptrStringCertConfig("RSA 2048"),
							},
						},
						ClientCertInfo: &teov20220901.MutualTLS{
							Switch: ptrStringCertConfig("on"),
							CertInfos: []*teov20220901.CertificateInfo{
								{
									CertId:     ptrStringCertConfig("client-cert-001"),
									Alias:      ptrStringCertConfig("client-ca-cert"),
									Type:       ptrStringCertConfig("upload"),
									ExpireTime: ptrStringCertConfig("2026-06-01T00:00:00Z"),
									DeployTime: ptrStringCertConfig("2025-06-01T00:00:00Z"),
									SignAlgo:   ptrStringCertConfig("RSA 4096"),
								},
							},
						},
					},
				},
			},
			RequestId: ptrStringCertConfig("fake-request-id"),
		}
		return resp, nil
	})

	// Mock DescribeZones for the read flow
	patches.ApplyMethodFunc(teoClient, "DescribeZones", func(request *teov20220901.DescribeZonesRequest) (*teov20220901.DescribeZonesResponse, error) {
		resp := teov20220901.NewDescribeZonesResponse()
		resp.Response = &teov20220901.DescribeZonesResponseParams{
			Zones: []*teov20220901.Zone{
				{
					ZoneId:   ptrStringCertConfig("zone-test123"),
					ZoneName: ptrStringCertConfig("example.com"),
				},
			},
			RequestId: ptrStringCertConfig("fake-request-id"),
		}
		return resp, nil
	})

	meta := newCertConfigMockMeta()
	res := teo.ResourceTencentCloudTeoCertificateConfig()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-test123",
		"host":    "test.example.com",
	})
	d.SetId("zone-test123#test.example.com")

	err := res.Read(d, meta)
	assert.NoError(t, err)

	// Verify client_cert_info was read into state
	clientCertInfo := d.Get("client_cert_info").([]interface{})
	assert.Len(t, clientCertInfo, 1)
	clientCertInfoMap := clientCertInfo[0].(map[string]interface{})
	assert.Equal(t, "on", clientCertInfoMap["switch"])

	certInfos := clientCertInfoMap["cert_infos"].([]interface{})
	assert.Len(t, certInfos, 1)
	certInfoMap := certInfos[0].(map[string]interface{})
	assert.Equal(t, "client-cert-001", certInfoMap["cert_id"])
	assert.Equal(t, "client-ca-cert", certInfoMap["alias"])
	assert.Equal(t, "upload", certInfoMap["type"])
	assert.Equal(t, "2026-06-01T00:00:00Z", certInfoMap["expire_time"])
	assert.Equal(t, "2025-06-01T00:00:00Z", certInfoMap["deploy_time"])
	assert.Equal(t, "RSA 4096", certInfoMap["sign_algo"])

	// Verify mode was also read
	assert.Equal(t, "sslcert", d.Get("mode"))
}

// TestTeoCertificateConfigClientCertInfo_UpdateClientCertInfo tests updating client_cert_info from off to on
func TestTeoCertificateConfigClientCertInfo_UpdateClientCertInfo(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newCertConfigMockMeta().client, "UseTeoClient", teoClient)

	// Mock ModifyHostsCertificate to verify ClientCertInfo is updated
	patches.ApplyMethodFunc(teoClient, "ModifyHostsCertificate", func(request *teov20220901.ModifyHostsCertificateRequest) (*teov20220901.ModifyHostsCertificateResponse, error) {
		assert.NotNil(t, request.ClientCertInfo)
		assert.Equal(t, "on", *request.ClientCertInfo.Switch)
		assert.Len(t, request.ClientCertInfo.CertInfos, 2)
		assert.Equal(t, "cert-new-001", *request.ClientCertInfo.CertInfos[0].CertId)
		assert.Equal(t, "cert-new-002", *request.ClientCertInfo.CertInfos[1].CertId)

		resp := teov20220901.NewModifyHostsCertificateResponse()
		resp.Response = &teov20220901.ModifyHostsCertificateResponseParams{
			RequestId: ptrStringCertConfig("fake-request-id"),
		}
		return resp, nil
	})

	// Mock DescribeAccelerationDomains for the read after update
	patches.ApplyMethodFunc(teoClient, "DescribeAccelerationDomains", func(request *teov20220901.DescribeAccelerationDomainsRequest) (*teov20220901.DescribeAccelerationDomainsResponse, error) {
		resp := teov20220901.NewDescribeAccelerationDomainsResponse()
		resp.Response = &teov20220901.DescribeAccelerationDomainsResponseParams{
			AccelerationDomains: []*teov20220901.AccelerationDomain{
				{
					DomainName:   ptrStringCertConfig("test.example.com"),
					DomainStatus: ptrStringCertConfig("online"),
					Certificate: &teov20220901.AccelerationDomainCertificate{
						Mode: ptrStringCertConfig("sslcert"),
						List: []*teov20220901.CertificateInfo{
							{
								CertId:   ptrStringCertConfig("server-cert-001"),
								Alias:    ptrStringCertConfig("server-cert"),
								Type:     ptrStringCertConfig("managed"),
								SignAlgo: ptrStringCertConfig("RSA 2048"),
							},
						},
						ClientCertInfo: &teov20220901.MutualTLS{
							Switch: ptrStringCertConfig("on"),
							CertInfos: []*teov20220901.CertificateInfo{
								{
									CertId:   ptrStringCertConfig("cert-new-001"),
									Alias:    ptrStringCertConfig("new-cert-1"),
									Type:     ptrStringCertConfig("upload"),
									SignAlgo: ptrStringCertConfig("ECDSA P256"),
								},
								{
									CertId:   ptrStringCertConfig("cert-new-002"),
									Alias:    ptrStringCertConfig("new-cert-2"),
									Type:     ptrStringCertConfig("default"),
									SignAlgo: ptrStringCertConfig("RSA 2048"),
								},
							},
						},
					},
				},
			},
			RequestId: ptrStringCertConfig("fake-request-id"),
		}
		return resp, nil
	})

	// Mock DescribeZones for the read flow
	patches.ApplyMethodFunc(teoClient, "DescribeZones", func(request *teov20220901.DescribeZonesRequest) (*teov20220901.DescribeZonesResponse, error) {
		resp := teov20220901.NewDescribeZonesResponse()
		resp.Response = &teov20220901.DescribeZonesResponseParams{
			Zones: []*teov20220901.Zone{
				{
					ZoneId:   ptrStringCertConfig("zone-test123"),
					ZoneName: ptrStringCertConfig("example.com"),
				},
			},
			RequestId: ptrStringCertConfig("fake-request-id"),
		}
		return resp, nil
	})

	meta := newCertConfigMockMeta()
	res := teo.ResourceTencentCloudTeoCertificateConfig()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-test123",
		"host":    "test.example.com",
		"mode":    "sslcert",
		"server_cert_info": []interface{}{
			map[string]interface{}{
				"cert_id": "server-cert-001",
			},
		},
		"client_cert_info": []interface{}{
			map[string]interface{}{
				"switch": "on",
				"cert_infos": []interface{}{
					map[string]interface{}{
						"cert_id": "cert-new-001",
					},
					map[string]interface{}{
						"cert_id": "cert-new-002",
					},
				},
			},
		},
	})
	d.SetId("zone-test123#test.example.com")

	err := res.Update(d, meta)
	assert.NoError(t, err)

	// Verify client_cert_info was read into state after update
	clientCertInfo := d.Get("client_cert_info").([]interface{})
	assert.Len(t, clientCertInfo, 1)
	clientCertInfoMap := clientCertInfo[0].(map[string]interface{})
	assert.Equal(t, "on", clientCertInfoMap["switch"])

	certInfos := clientCertInfoMap["cert_infos"].([]interface{})
	assert.Len(t, certInfos, 2)
	assert.Equal(t, "cert-new-001", certInfos[0].(map[string]interface{})["cert_id"])
	assert.Equal(t, "cert-new-002", certInfos[1].(map[string]interface{})["cert_id"])
}

// TestTeoCertificateConfigClientCertInfo_APIError tests API error handling during update
func TestTeoCertificateConfigClientCertInfo_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newCertConfigMockMeta().client, "UseTeoClient", teoClient)

	// Mock ModifyHostsCertificate to return an error
	patches.ApplyMethodFunc(teoClient, "ModifyHostsCertificate", func(request *teov20220901.ModifyHostsCertificateRequest) (*teov20220901.ModifyHostsCertificateResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=ResourceNotFound, Message=Zone not found")
	})

	meta := newCertConfigMockMeta()
	res := teo.ResourceTencentCloudTeoCertificateConfig()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-invalid",
		"host":    "test.example.com",
		"mode":    "sslcert",
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
	d.SetId("zone-invalid#test.example.com")

	err := res.Update(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ResourceNotFound")
}

// TestTeoCertificateConfigClientCertInfo_Schema validates the client_cert_info schema definition
func TestTeoCertificateConfigClientCertInfo_Schema(t *testing.T) {
	res := teo.ResourceTencentCloudTeoCertificateConfig()

	assert.NotNil(t, res)
	assert.Contains(t, res.Schema, "client_cert_info")

	clientCertInfo := res.Schema["client_cert_info"]
	assert.Equal(t, schema.TypeList, clientCertInfo.Type)
	assert.True(t, clientCertInfo.Optional)
	assert.True(t, clientCertInfo.Computed)
	assert.Equal(t, 1, clientCertInfo.MaxItems)

	// Check nested schema
	elem := clientCertInfo.Elem.(*schema.Resource)
	assert.Contains(t, elem.Schema, "switch")
	assert.Contains(t, elem.Schema, "cert_infos")

	// Check switch field
	switchField := elem.Schema["switch"]
	assert.Equal(t, schema.TypeString, switchField.Type)
	assert.True(t, switchField.Required)

	// Check cert_infos field
	certInfosField := elem.Schema["cert_infos"]
	assert.Equal(t, schema.TypeList, certInfosField.Type)
	assert.True(t, certInfosField.Optional)
	assert.True(t, certInfosField.Computed)

	// Check cert_infos nested schema
	certInfosElem := certInfosField.Elem.(*schema.Resource)
	assert.Contains(t, certInfosElem.Schema, "cert_id")
	assert.Contains(t, certInfosElem.Schema, "alias")
	assert.Contains(t, certInfosElem.Schema, "type")
	assert.Contains(t, certInfosElem.Schema, "expire_time")
	assert.Contains(t, certInfosElem.Schema, "deploy_time")
	assert.Contains(t, certInfosElem.Schema, "sign_algo")

	// cert_id is Required
	assert.True(t, certInfosElem.Schema["cert_id"].Required)
	// Other fields are Computed
	assert.True(t, certInfosElem.Schema["alias"].Computed)
	assert.True(t, certInfosElem.Schema["type"].Computed)
	assert.True(t, certInfosElem.Schema["expire_time"].Computed)
	assert.True(t, certInfosElem.Schema["deploy_time"].Computed)
	assert.True(t, certInfosElem.Schema["sign_algo"].Computed)
}
