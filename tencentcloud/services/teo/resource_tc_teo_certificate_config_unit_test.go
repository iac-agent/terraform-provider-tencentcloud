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

// go test ./tencentcloud/services/teo/ -run "TestTeoCertificateConfig" -v -count=1 -gcflags="all=-l"

// TestTeoCertificateConfig_Schema validates that apply_type and client_cert_info are present in schema
func TestTeoCertificateConfig_Schema(t *testing.T) {
	res := teo.ResourceTencentCloudTeoCertificateConfig()

	assert.NotNil(t, res)
	assert.NotNil(t, res.Create)
	assert.NotNil(t, res.Read)
	assert.NotNil(t, res.Update)
	assert.NotNil(t, res.Delete)
	assert.NotNil(t, res.Importer)

	// Check apply_type field
	assert.Contains(t, res.Schema, "apply_type")
	applyType := res.Schema["apply_type"]
	assert.Equal(t, schema.TypeString, applyType.Type)
	assert.True(t, applyType.Optional)
	assert.True(t, applyType.Computed)
	assert.False(t, applyType.ForceNew)

	// Check client_cert_info field
	assert.Contains(t, res.Schema, "client_cert_info")
	clientCertInfo := res.Schema["client_cert_info"]
	assert.Equal(t, schema.TypeList, clientCertInfo.Type)
	assert.True(t, clientCertInfo.Optional)
	assert.True(t, clientCertInfo.Computed)
	assert.Equal(t, 1, clientCertInfo.MaxItems)
	assert.False(t, clientCertInfo.ForceNew)
}

// TestTeoCertificateConfig_Create_WithApplyType tests Create with apply_type parameter
func TestTeoCertificateConfig_Create_WithApplyType(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoClient", teoClient)
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)

	// Mock ModifyHostsCertificate
	patches.ApplyMethodFunc(teoClient, "ModifyHostsCertificate", func(request *teov20220901.ModifyHostsCertificateRequest) (*teov20220901.ModifyHostsCertificateResponse, error) {
		resp := teov20220901.NewModifyHostsCertificateResponse()
		resp.Response = &teov20220901.ModifyHostsCertificateResponseParams{
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	// Mock DescribeAccelerationDomains for read and status check
	patches.ApplyMethodFunc(teoClient, "DescribeAccelerationDomains", func(request *teov20220901.DescribeAccelerationDomainsRequest) (*teov20220901.DescribeAccelerationDomainsResponse, error) {
		resp := teov20220901.NewDescribeAccelerationDomainsResponse()
		resp.Response = &teov20220901.DescribeAccelerationDomainsResponseParams{
			AccelerationDomains: []*teov20220901.AccelerationDomain{
				{
					DomainName:   ptrString("test.example.com"),
					DomainStatus: ptrString("online"),
					Certificate: &teov20220901.AccelerationDomainCertificate{
						Mode: ptrString("sslcert"),
					},
				},
			},
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	// Mock DescribeZones for zone name lookup
	patches.ApplyMethodFunc(teoClient, "DescribeZones", func(request *teov20220901.DescribeZonesRequest) (*teov20220901.DescribeZonesResponse, error) {
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

	// Mock DescribeHostsSetting for apply_type read-back
	patches.ApplyMethodFunc(teoClient, "DescribeHostsSetting", func(request *teov20220901.DescribeHostsSettingRequest) (*teov20220901.DescribeHostsSettingResponse, error) {
		resp := teov20220901.NewDescribeHostsSettingResponse()
		resp.Response = &teov20220901.DescribeHostsSettingResponseParams{
			DetailHosts: []*teov20220901.DetailHost{
				{
					Host: ptrString("test.example.com"),
					Https: &teov20220901.Https{
						ApplyType: ptrString("apply"),
					},
				},
			},
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoCertificateConfig()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":    "zone-1234567890",
		"host":       "test.example.com",
		"mode":       "sslcert",
		"apply_type": "apply",
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-1234567890#test.example.com", d.Id())
}

// TestTeoCertificateConfig_Create_WithClientCertInfo tests Create with client_cert_info parameter
func TestTeoCertificateConfig_Create_WithClientCertInfo(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoClient", teoClient)
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)

	// Mock ModifyHostsCertificate
	patches.ApplyMethodFunc(teoClient, "ModifyHostsCertificate", func(request *teov20220901.ModifyHostsCertificateRequest) (*teov20220901.ModifyHostsCertificateResponse, error) {
		resp := teov20220901.NewModifyHostsCertificateResponse()
		resp.Response = &teov20220901.ModifyHostsCertificateResponseParams{
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	// Mock DescribeAccelerationDomains for read
	patches.ApplyMethodFunc(teoClient, "DescribeAccelerationDomains", func(request *teov20220901.DescribeAccelerationDomainsRequest) (*teov20220901.DescribeAccelerationDomainsResponse, error) {
		resp := teov20220901.NewDescribeAccelerationDomainsResponse()
		resp.Response = &teov20220901.DescribeAccelerationDomainsResponseParams{
			AccelerationDomains: []*teov20220901.AccelerationDomain{
				{
					DomainName:   ptrString("test.example.com"),
					DomainStatus: ptrString("online"),
					Certificate: &teov20220901.AccelerationDomainCertificate{
						Mode: ptrString("sslcert"),
						ClientCertInfo: &teov20220901.MutualTLS{
							Switch: ptrString("on"),
							CertInfos: []*teov20220901.CertificateInfo{
								{
									CertId: ptrString("cert-abc123"),
									Alias:  ptrString("test-cert"),
									Type:   ptrString("upload"),
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

	// Mock DescribeZones for zone name lookup
	patches.ApplyMethodFunc(teoClient, "DescribeZones", func(request *teov20220901.DescribeZonesRequest) (*teov20220901.DescribeZonesResponse, error) {
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

	// Mock DescribeHostsSetting for apply_type read-back
	patches.ApplyMethodFunc(teoClient, "DescribeHostsSetting", func(request *teov20220901.DescribeHostsSettingRequest) (*teov20220901.DescribeHostsSettingResponse, error) {
		resp := teov20220901.NewDescribeHostsSettingResponse()
		resp.Response = &teov20220901.DescribeHostsSettingResponseParams{
			DetailHosts: []*teov20220901.DetailHost{
				{
					Host: ptrString("test.example.com"),
					Https: &teov20220901.Https{
						ApplyType: ptrString("none"),
					},
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

// TestTeoCertificateConfig_Read_WithApplyType tests Read populates apply_type from DescribeHostsSetting
func TestTeoCertificateConfig_Read_WithApplyType(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoClient", teoClient)
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)

	// Mock DescribeAccelerationDomains
	patches.ApplyMethodFunc(teoClient, "DescribeAccelerationDomains", func(request *teov20220901.DescribeAccelerationDomainsRequest) (*teov20220901.DescribeAccelerationDomainsResponse, error) {
		resp := teov20220901.NewDescribeAccelerationDomainsResponse()
		resp.Response = &teov20220901.DescribeAccelerationDomainsResponseParams{
			AccelerationDomains: []*teov20220901.AccelerationDomain{
				{
					DomainName: ptrString("test.example.com"),
					Certificate: &teov20220901.AccelerationDomainCertificate{
						Mode: ptrString("sslcert"),
					},
				},
			},
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	// Mock DescribeZones
	patches.ApplyMethodFunc(teoClient, "DescribeZones", func(request *teov20220901.DescribeZonesRequest) (*teov20220901.DescribeZonesResponse, error) {
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

	// Mock DescribeHostsSetting to return apply_type
	patches.ApplyMethodFunc(teoClient, "DescribeHostsSetting", func(request *teov20220901.DescribeHostsSettingRequest) (*teov20220901.DescribeHostsSettingResponse, error) {
		resp := teov20220901.NewDescribeHostsSettingResponse()
		resp.Response = &teov20220901.DescribeHostsSettingResponseParams{
			DetailHosts: []*teov20220901.DetailHost{
				{
					Host: ptrString("test.example.com"),
					Https: &teov20220901.Https{
						ApplyType: ptrString("apply"),
					},
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
	})
	d.SetId("zone-1234567890#test.example.com")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "apply", d.Get("apply_type"))
	assert.Equal(t, "sslcert", d.Get("mode"))
}

// TestTeoCertificateConfig_Read_WithClientCertInfo tests Read populates client_cert_info from API response
func TestTeoCertificateConfig_Read_WithClientCertInfo(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoClient", teoClient)
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)

	// Mock DescribeAccelerationDomains with ClientCertInfo
	patches.ApplyMethodFunc(teoClient, "DescribeAccelerationDomains", func(request *teov20220901.DescribeAccelerationDomainsRequest) (*teov20220901.DescribeAccelerationDomainsResponse, error) {
		resp := teov20220901.NewDescribeAccelerationDomainsResponse()
		resp.Response = &teov20220901.DescribeAccelerationDomainsResponseParams{
			AccelerationDomains: []*teov20220901.AccelerationDomain{
				{
					DomainName: ptrString("test.example.com"),
					Certificate: &teov20220901.AccelerationDomainCertificate{
						Mode: ptrString("sslcert"),
						ClientCertInfo: &teov20220901.MutualTLS{
							Switch: ptrString("on"),
							CertInfos: []*teov20220901.CertificateInfo{
								{
									CertId: ptrString("cert-abc123"),
									Alias:  ptrString("test-cert"),
									Type:   ptrString("upload"),
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

	// Mock DescribeZones
	patches.ApplyMethodFunc(teoClient, "DescribeZones", func(request *teov20220901.DescribeZonesRequest) (*teov20220901.DescribeZonesResponse, error) {
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

	// Mock DescribeHostsSetting
	patches.ApplyMethodFunc(teoClient, "DescribeHostsSetting", func(request *teov20220901.DescribeHostsSettingRequest) (*teov20220901.DescribeHostsSettingResponse, error) {
		resp := teov20220901.NewDescribeHostsSettingResponse()
		resp.Response = &teov20220901.DescribeHostsSettingResponseParams{
			DetailHosts: []*teov20220901.DetailHost{
				{
					Host: ptrString("test.example.com"),
					Https: &teov20220901.Https{
						ApplyType: ptrString("none"),
					},
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
	})
	d.SetId("zone-1234567890#test.example.com")

	err := res.Read(d, meta)
	assert.NoError(t, err)

	// Verify client_cert_info is populated
	clientCertInfo := d.Get("client_cert_info").([]interface{})
	assert.Len(t, clientCertInfo, 1)
	clientCertInfoMap := clientCertInfo[0].(map[string]interface{})
	assert.Equal(t, "on", clientCertInfoMap["switch"])

	certInfos := clientCertInfoMap["cert_infos"].([]interface{})
	assert.Len(t, certInfos, 1)
	certInfoMap := certInfos[0].(map[string]interface{})
	assert.Equal(t, "cert-abc123", certInfoMap["cert_id"])
}

// TestTeoCertificateConfig_Update_WithApplyType tests Update with apply_type parameter
func TestTeoCertificateConfig_Update_WithApplyType(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoClient", teoClient)
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)

	// Mock ModifyHostsCertificate
	patches.ApplyMethodFunc(teoClient, "ModifyHostsCertificate", func(request *teov20220901.ModifyHostsCertificateRequest) (*teov20220901.ModifyHostsCertificateResponse, error) {
		// Verify ApplyType is set in the request
		assert.NotNil(t, request.ApplyType)
		assert.Equal(t, "apply", *request.ApplyType)

		resp := teov20220901.NewModifyHostsCertificateResponse()
		resp.Response = &teov20220901.ModifyHostsCertificateResponseParams{
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	// Mock DescribeAccelerationDomains for read after update
	patches.ApplyMethodFunc(teoClient, "DescribeAccelerationDomains", func(request *teov20220901.DescribeAccelerationDomainsRequest) (*teov20220901.DescribeAccelerationDomainsResponse, error) {
		resp := teov20220901.NewDescribeAccelerationDomainsResponse()
		resp.Response = &teov20220901.DescribeAccelerationDomainsResponseParams{
			AccelerationDomains: []*teov20220901.AccelerationDomain{
				{
					DomainName:   ptrString("test.example.com"),
					DomainStatus: ptrString("online"),
					Certificate: &teov20220901.AccelerationDomainCertificate{
						Mode: ptrString("sslcert"),
					},
				},
			},
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	// Mock DescribeZones
	patches.ApplyMethodFunc(teoClient, "DescribeZones", func(request *teov20220901.DescribeZonesRequest) (*teov20220901.DescribeZonesResponse, error) {
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

	// Mock DescribeHostsSetting
	patches.ApplyMethodFunc(teoClient, "DescribeHostsSetting", func(request *teov20220901.DescribeHostsSettingRequest) (*teov20220901.DescribeHostsSettingResponse, error) {
		resp := teov20220901.NewDescribeHostsSettingResponse()
		resp.Response = &teov20220901.DescribeHostsSettingResponseParams{
			DetailHosts: []*teov20220901.DetailHost{
				{
					Host: ptrString("test.example.com"),
					Https: &teov20220901.Https{
						ApplyType: ptrString("apply"),
					},
				},
			},
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoCertificateConfig()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":    "zone-1234567890",
		"host":       "test.example.com",
		"mode":       "sslcert",
		"apply_type": "apply",
	})
	d.SetId("zone-1234567890#test.example.com")

	err := res.Update(d, meta)
	assert.NoError(t, err)
}

// TestTeoCertificateConfig_Create_APIError tests Create handles API error
func TestTeoCertificateConfig_Create_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoClient", teoClient)
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)

	// Mock ModifyHostsCertificate to return error
	patches.ApplyMethodFunc(teoClient, "ModifyHostsCertificate", func(request *teov20220901.ModifyHostsCertificateRequest) (*teov20220901.ModifyHostsCertificateResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=InvalidParameter, Message=Invalid zone_id")
	})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoCertificateConfig()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":    "zone-invalid",
		"host":       "test.example.com",
		"mode":       "sslcert",
		"apply_type": "apply",
	})

	err := res.Create(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "InvalidParameter")
}
