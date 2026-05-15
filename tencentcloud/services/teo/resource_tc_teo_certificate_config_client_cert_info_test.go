package teo_test

import (
	"fmt"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/teo"
)

// go test ./tencentcloud/services/teo/ -run "TestTeoCertificateConfigClientCertInfo" -v -count=1 -gcflags="all=-l"

// TestTeoCertificateConfigClientCertInfo_Schema validates client_cert_info schema definition
func TestTeoCertificateConfigClientCertInfo_Schema(t *testing.T) {
	res := teo.ResourceTencentCloudTeoCertificateConfig()

	assert.NotNil(t, res)
	assert.Contains(t, res.Schema, "client_cert_info")

	clientCertInfo := res.Schema["client_cert_info"]
	assert.Equal(t, schema.TypeList, clientCertInfo.Type)
	assert.True(t, clientCertInfo.Optional)
	assert.True(t, clientCertInfo.Computed)
	assert.Equal(t, 1, clientCertInfo.MaxItems)

	elem := clientCertInfo.Elem.(*schema.Resource)
	assert.Contains(t, elem.Schema, "switch")
	switchField := elem.Schema["switch"]
	assert.Equal(t, schema.TypeString, switchField.Type)
	assert.True(t, switchField.Required)

	assert.Contains(t, elem.Schema, "cert_infos")
	certInfos := elem.Schema["cert_infos"]
	assert.Equal(t, schema.TypeList, certInfos.Type)
	assert.True(t, certInfos.Optional)
	assert.True(t, certInfos.Computed)

	certInfosElem := certInfos.Elem.(*schema.Resource)
	assert.Contains(t, certInfosElem.Schema, "cert_id")
	assert.Equal(t, schema.TypeString, certInfosElem.Schema["cert_id"].Type)
	assert.True(t, certInfosElem.Schema["cert_id"].Required)

	assert.Contains(t, certInfosElem.Schema, "alias")
	assert.Equal(t, schema.TypeString, certInfosElem.Schema["alias"].Type)
	assert.True(t, certInfosElem.Schema["alias"].Computed)

	assert.Contains(t, certInfosElem.Schema, "type")
	assert.Equal(t, schema.TypeString, certInfosElem.Schema["type"].Type)
	assert.True(t, certInfosElem.Schema["type"].Computed)

	assert.Contains(t, certInfosElem.Schema, "expire_time")
	assert.Equal(t, schema.TypeString, certInfosElem.Schema["expire_time"].Type)
	assert.True(t, certInfosElem.Schema["expire_time"].Computed)

	assert.Contains(t, certInfosElem.Schema, "deploy_time")
	assert.Equal(t, schema.TypeString, certInfosElem.Schema["deploy_time"].Type)
	assert.True(t, certInfosElem.Schema["deploy_time"].Computed)

	assert.Contains(t, certInfosElem.Schema, "sign_algo")
	assert.Equal(t, schema.TypeString, certInfosElem.Schema["sign_algo"].Type)
	assert.True(t, certInfosElem.Schema["sign_algo"].Computed)

	assert.Contains(t, certInfosElem.Schema, "status")
	assert.Equal(t, schema.TypeString, certInfosElem.Schema["status"].Type)
	assert.True(t, certInfosElem.Schema["status"].Computed)
}

// TestTeoCertificateConfigClientCertInfo_ReadWithClientCertInfo tests that read path populates
// client_cert_info correctly when ClientCertInfo is present in the API response
func TestTeoCertificateConfigClientCertInfo_ReadWithClientCertInfo(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeAccelerationDomains",
		func(request *teov20220901.DescribeAccelerationDomainsRequest) (*teov20220901.DescribeAccelerationDomainsResponse, error) {
			resp := teov20220901.NewDescribeAccelerationDomainsResponse()
			resp.Response = &teov20220901.DescribeAccelerationDomainsResponseParams{
				AccelerationDomains: []*teov20220901.AccelerationDomain{
					{
						DomainName: ptrString("test.example.com"),
						Certificate: &teov20220901.AccelerationDomainCertificate{
							Mode: ptrString("sslcert"),
							List: []*teov20220901.CertificateInfo{
								{
									CertId: ptrString("server-cert-id-1"),
								},
							},
							ClientCertInfo: &teov20220901.MutualTLS{
								Switch: ptrString("on"),
								CertInfos: []*teov20220901.CertificateInfo{
									{
										CertId:     ptrString("client-cert-id-1"),
										Alias:      ptrString("client-cert-alias"),
										Type:       ptrString("upload"),
										ExpireTime: ptrString("2025-12-31T23:59:59Z"),
										DeployTime: ptrString("2024-01-01T00:00:00Z"),
										SignAlgo:   ptrString("RSA 2048"),
										Status:     ptrString("deployed"),
									},
								},
							},
						},
					},
				},
				RequestId: ptrString("fake-request-id"),
			}
			return resp, nil
		})

	patches.ApplyMethodFunc(teoClient, "DescribeZones",
		func(request *teov20220901.DescribeZonesRequest) (*teov20220901.DescribeZonesResponse, error) {
			resp := teov20220901.NewDescribeZonesResponse()
			resp.Response = &teov20220901.DescribeZonesResponseParams{
				Zones: []*teov20220901.Zone{
					{
						ZoneId:   ptrString("zone-1234567890"),
						ZoneName: ptrString("example.com"),
					},
				},
				RequestId: ptrString("fake-request-id"),
			}
			return resp, nil
		})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoCertificateConfig()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-1234567890",
		"host":    "test.example.com",
		"mode":    "sslcert",
	})
	d.SetId("zone-1234567890#test.example.com")

	err := res.Read(d, meta)
	assert.NoError(t, err)

	clientCertInfoList := d.Get("client_cert_info").([]interface{})
	assert.Equal(t, 1, len(clientCertInfoList))

	clientCertInfoMap := clientCertInfoList[0].(map[string]interface{})
	assert.Equal(t, "on", clientCertInfoMap["switch"])

	certInfosList := clientCertInfoMap["cert_infos"].([]interface{})
	assert.Equal(t, 1, len(certInfosList))

	certInfoMap := certInfosList[0].(map[string]interface{})
	assert.Equal(t, "client-cert-id-1", certInfoMap["cert_id"])
	assert.Equal(t, "client-cert-alias", certInfoMap["alias"])
	assert.Equal(t, "upload", certInfoMap["type"])
	assert.Equal(t, "2025-12-31T23:59:59Z", certInfoMap["expire_time"])
	assert.Equal(t, "2024-01-01T00:00:00Z", certInfoMap["deploy_time"])
	assert.Equal(t, "RSA 2048", certInfoMap["sign_algo"])
	assert.Equal(t, "deployed", certInfoMap["status"])
}

// TestTeoCertificateConfigClientCertInfo_ReadWithNilClientCertInfo tests that read path
// does not set client_cert_info when ClientCertInfo is nil in the API response
func TestTeoCertificateConfigClientCertInfo_ReadWithNilClientCertInfo(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeAccelerationDomains",
		func(request *teov20220901.DescribeAccelerationDomainsRequest) (*teov20220901.DescribeAccelerationDomainsResponse, error) {
			resp := teov20220901.NewDescribeAccelerationDomainsResponse()
			resp.Response = &teov20220901.DescribeAccelerationDomainsResponseParams{
				AccelerationDomains: []*teov20220901.AccelerationDomain{
					{
						DomainName: ptrString("test.example.com"),
						Certificate: &teov20220901.AccelerationDomainCertificate{
							Mode: ptrString("sslcert"),
							List: []*teov20220901.CertificateInfo{
								{
									CertId: ptrString("server-cert-id-1"),
								},
							},
						},
					},
				},
				RequestId: ptrString("fake-request-id"),
			}
			return resp, nil
		})

	patches.ApplyMethodFunc(teoClient, "DescribeZones",
		func(request *teov20220901.DescribeZonesRequest) (*teov20220901.DescribeZonesResponse, error) {
			resp := teov20220901.NewDescribeZonesResponse()
			resp.Response = &teov20220901.DescribeZonesResponseParams{
				Zones: []*teov20220901.Zone{
					{
						ZoneId:   ptrString("zone-1234567890"),
						ZoneName: ptrString("example.com"),
					},
				},
				RequestId: ptrString("fake-request-id"),
			}
			return resp, nil
		})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoCertificateConfig()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-1234567890",
		"host":    "test.example.com",
		"mode":    "sslcert",
	})
	d.SetId("zone-1234567890#test.example.com")

	err := res.Read(d, meta)
	assert.NoError(t, err)

	clientCertInfoList := d.Get("client_cert_info").([]interface{})
	assert.Equal(t, 0, len(clientCertInfoList))
}

// TestTeoCertificateConfigClientCertInfo_UpdateWithClientCertInfo tests that update path
// correctly sets ClientCertInfo in the ModifyHostsCertificate request
func TestTeoCertificateConfigClientCertInfo_UpdateWithClientCertInfo(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoClient", teoClient)

	var capturedRequest *teov20220901.ModifyHostsCertificateRequest
	patches.ApplyMethodFunc(teoClient, "ModifyHostsCertificate",
		func(request *teov20220901.ModifyHostsCertificateRequest) (*teov20220901.ModifyHostsCertificateResponse, error) {
			capturedRequest = request
			resp := teov20220901.NewModifyHostsCertificateResponse()
			resp.Response = &teov20220901.ModifyHostsCertificateResponseParams{
				RequestId: ptrString("fake-request-id"),
			}
			return resp, nil
		})

	// Mock DescribeAccelerationDomains for the read after update and CheckAccelerationDomainStatus
	patches.ApplyMethodFunc(teoClient, "DescribeAccelerationDomains",
		func(request *teov20220901.DescribeAccelerationDomainsRequest) (*teov20220901.DescribeAccelerationDomainsResponse, error) {
			resp := teov20220901.NewDescribeAccelerationDomainsResponse()
			resp.Response = &teov20220901.DescribeAccelerationDomainsResponseParams{
				AccelerationDomains: []*teov20220901.AccelerationDomain{
					{
						DomainName:   ptrString("test.example.com"),
						DomainStatus: ptrString("online"),
						Certificate: &teov20220901.AccelerationDomainCertificate{
							Mode: ptrString("sslcert"),
							List: []*teov20220901.CertificateInfo{
								{
									CertId: ptrString("server-cert-id-1"),
								},
							},
							ClientCertInfo: &teov20220901.MutualTLS{
								Switch: ptrString("on"),
								CertInfos: []*teov20220901.CertificateInfo{
									{
										CertId: ptrString("client-cert-id-1"),
									},
								},
							},
						},
					},
				},
				RequestId: ptrString("fake-request-id"),
			}
			return resp, nil
		})

	// Mock DescribeZones for the read after update
	patches.ApplyMethodFunc(teoClient, "DescribeZones",
		func(request *teov20220901.DescribeZonesRequest) (*teov20220901.DescribeZonesResponse, error) {
			resp := teov20220901.NewDescribeZonesResponse()
			resp.Response = &teov20220901.DescribeZonesResponseParams{
				Zones: []*teov20220901.Zone{
					{
						ZoneId:   ptrString("zone-1234567890"),
						ZoneName: ptrString("example.com"),
					},
				},
				RequestId: ptrString("fake-request-id"),
			}
			return resp, nil
		})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoCertificateConfig()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-1234567890",
		"host":    "test.example.com",
		"mode":    "sslcert",
		"client_cert_info": []interface{}{
			map[string]interface{}{
				"switch": "on",
				"cert_infos": []interface{}{
					map[string]interface{}{
						"cert_id": "client-cert-id-1",
					},
				},
			},
		},
	})
	d.SetId("zone-1234567890#test.example.com")

	err := res.Update(d, meta)
	assert.NoError(t, err)

	// Verify the ModifyHostsCertificate request was captured and has ClientCertInfo set
	assert.NotNil(t, capturedRequest)
	assert.NotNil(t, capturedRequest.ClientCertInfo)
	assert.Equal(t, "on", *capturedRequest.ClientCertInfo.Switch)
	assert.Equal(t, 1, len(capturedRequest.ClientCertInfo.CertInfos))
	assert.Equal(t, "client-cert-id-1", *capturedRequest.ClientCertInfo.CertInfos[0].CertId)
}

// TestTeoCertificateConfigClientCertInfo_UpdateAPIError tests that update path handles API errors
func TestTeoCertificateConfigClientCertInfo_UpdateAPIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "ModifyHostsCertificate",
		func(request *teov20220901.ModifyHostsCertificateRequest) (*teov20220901.ModifyHostsCertificateResponse, error) {
			return nil, fmt.Errorf("[TencentCloudSDKError] Code=InvalidParameter, Message=Invalid cert_id")
		})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoCertificateConfig()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-1234567890",
		"host":    "test.example.com",
		"mode":    "sslcert",
		"client_cert_info": []interface{}{
			map[string]interface{}{
				"switch": "on",
				"cert_infos": []interface{}{
					map[string]interface{}{
						"cert_id": "invalid-cert-id",
					},
				},
			},
		},
	})
	d.SetId("zone-1234567890#test.example.com")

	err := res.Update(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "InvalidParameter")
}
