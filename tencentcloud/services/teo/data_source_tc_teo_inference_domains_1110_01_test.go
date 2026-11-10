package teo_test

import (
	"context"
	"errors"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/assert"

	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/sharedmeta"
	svcteo "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/teo"
)

func ptrStrInfDomains(s string) *string { return &s }

// teoInferenceDomainsSchema returns the framework data source schema.
func teoInferenceDomainsSchema(t *testing.T) schema.Schema {
	d := svcteo.NewTeoInferenceDomainsDataSource()
	resp := &datasource.SchemaResponse{}
	d.Schema(context.Background(), datasource.SchemaRequest{}, resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("failed to build data source schema: %v", resp.Diagnostics)
	}
	return resp.Schema
}

// teoInferenceDomainsRootType builds the tftypes.Object type for the data
// source root, mirroring every schema attribute.
func teoInferenceDomainsRootType() tftypes.Object {
	return tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"zone_id":    tftypes.String,
			"service_id": tftypes.String,
			"sort_by":    tftypes.String,
			"sort_order": tftypes.String,
			"domains":    tftypes.List{ElementType: teoInferenceDomainsDomainType()},
		},
	}
}

func teoInferenceDomainsDomainType() tftypes.Object {
	return tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"domain":                  tftypes.String,
			"status":                  tftypes.String,
			"auth_switch":             tftypes.String,
			"cname":                   tftypes.String,
			"ownership_verification":  tftypes.List{ElementType: teoInferenceDomainsOwnershipVerificationType()},
			"certificate":             tftypes.List{ElementType: teoInferenceDomainsCertificateType()},
			"create_time":             tftypes.String,
			"update_time":             tftypes.String,
		},
	}
}

func teoInferenceDomainsOwnershipVerificationType() tftypes.Object {
	return tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"dns_verification":  tftypes.List{ElementType: teoInferenceDomainsDnsVerificationType()},
			"file_verification": tftypes.List{ElementType: teoInferenceDomainsFileVerificationType()},
			"ns_verification":   tftypes.List{ElementType: teoInferenceDomainsNsVerificationType()},
		},
	}
}

func teoInferenceDomainsDnsVerificationType() tftypes.Object {
	return tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"subdomain":    tftypes.String,
			"record_type":  tftypes.String,
			"record_value": tftypes.String,
		},
	}
}

func teoInferenceDomainsFileVerificationType() tftypes.Object {
	return tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"path":    tftypes.String,
			"content": tftypes.String,
		},
	}
}

func teoInferenceDomainsNsVerificationType() tftypes.Object {
	return tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"name_servers": tftypes.Set{ElementType: tftypes.String},
		},
	}
}

func teoInferenceDomainsCertificateType() tftypes.Object {
	return tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"host":           tftypes.String,
			"mode":           tftypes.String,
			"host_cert_info": tftypes.List{ElementType: teoInferenceDomainsHostCertInfoType()},
			"apply_type":     tftypes.String,
		},
	}
}

func teoInferenceDomainsHostCertInfoType() tftypes.Object {
	return tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"cert_id":     tftypes.String,
			"alias":       tftypes.String,
			"type":        tftypes.String,
			"expire_time": tftypes.String,
			"deploy_time": tftypes.String,
			"sign_algo":   tftypes.String,
			"status":      tftypes.String,
		},
	}
}

// buildTeoInferenceDomainsRawValue builds the raw config value with the given
// query inputs; the computed domains list stays null.
func buildTeoInferenceDomainsRawValue(zoneId, serviceId, sortBy, sortOrder string) tftypes.Value {
	attrs := map[string]tftypes.Value{
		"zone_id":    tftypes.NewValue(tftypes.String, zoneId),
		"service_id": tftypes.NewValue(tftypes.String, serviceId),
		"sort_by":    tftypes.NewValue(tftypes.String, nil),
		"sort_order": tftypes.NewValue(tftypes.String, nil),
		"domains":    tftypes.NewValue(tftypes.List{ElementType: teoInferenceDomainsDomainType()}, nil),
	}
	if sortBy != "" {
		attrs["sort_by"] = tftypes.NewValue(tftypes.String, sortBy)
	}
	if sortOrder != "" {
		attrs["sort_order"] = tftypes.NewValue(tftypes.String, sortOrder)
	}
	return tftypes.NewValue(teoInferenceDomainsRootType(), attrs)
}

// setTeoInferenceDomainsClient injects a mock connectivity client into the
// data source instance via the Configure method.
func setTeoInferenceDomainsClient(d *svcteo.TeoInferenceDomainsDataSource, client *connectivity.TencentCloudClient) {
	meta := &sharedmeta.ProviderMeta{Client: client}
	configureReq := datasource.ConfigureRequest{ProviderData: meta}
	configureResp := &datasource.ConfigureResponse{}
	d.Configure(context.Background(), configureReq, configureResp)
}

func newTeoInferenceDomainsReadRequest(sch schema.Schema, raw tftypes.Value) datasource.ReadRequest {
	req := datasource.ReadRequest{}
	req.Config.Raw = raw
	req.Config.Schema = sch
	return req
}

func newTeoInferenceDomainsReadResponse(sch schema.Schema) *datasource.ReadResponse {
	resp := &datasource.ReadResponse{}
	resp.State.Schema = sch
	return resp
}

// go test ./tencentcloud/services/teo/ -run "TestTeoInferenceDomains" -v -count=1 -gcflags="all=-l"

// TestTeoInferenceDomains_MetadataAndSchema validates metadata and schema definition.
func TestTeoInferenceDomains_MetadataAndSchema(t *testing.T) {
	d := svcteo.NewTeoInferenceDomainsDataSource()

	metaResp := &datasource.MetadataResponse{}
	d.Metadata(context.Background(), datasource.MetadataRequest{}, metaResp)
	assert.Equal(t, "tencentcloud_teo_inference_domains_1110_01", metaResp.TypeName)

	sch := teoInferenceDomainsSchema(t)
	attrs := sch.Attributes

	// Required inputs
	assert.True(t, attrs["zone_id"].(schema.StringAttribute).Required)
	assert.True(t, attrs["service_id"].(schema.StringAttribute).Required)

	// Optional inputs
	assert.True(t, attrs["sort_by"].(schema.StringAttribute).Optional)
	assert.True(t, attrs["sort_order"].(schema.StringAttribute).Optional)

	// Computed output
	domains := attrs["domains"].(schema.ListNestedAttribute)
	assert.True(t, domains.Computed)
	nested := domains.NestedObject.Attributes
	for _, leaf := range []string{"domain", "status", "auth_switch", "cname", "create_time", "update_time"} {
		assert.True(t, nested[leaf].(schema.StringAttribute).Computed, "nested attribute %s should be Computed", leaf)
	}
	ov := nested["ownership_verification"].(schema.ListNestedAttribute)
	assert.True(t, ov.Computed)
	ovNested := ov.NestedObject.Attributes
	assert.True(t, ovNested["dns_verification"].(schema.ListNestedAttribute).Computed)
	assert.True(t, ovNested["file_verification"].(schema.ListNestedAttribute).Computed)
	assert.True(t, ovNested["ns_verification"].(schema.ListNestedAttribute).Computed)
	cert := nested["certificate"].(schema.ListNestedAttribute)
	assert.True(t, cert.Computed)
	certNested := cert.NestedObject.Attributes
	for _, leaf := range []string{"host", "mode", "apply_type"} {
		assert.True(t, certNested[leaf].(schema.StringAttribute).Computed, "nested attribute %s should be Computed", leaf)
	}
	assert.True(t, certNested["host_cert_info"].(schema.ListNestedAttribute).Computed)
}

// TestTeoInferenceDomains_ReadBasic tests a successful read with one domain and
// asserts every leaf field is read back.
func TestTeoInferenceDomains_ReadBasic(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(client, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeInferenceDomainsWithContext", func(_ context.Context, request *teov20220901.DescribeInferenceDomainsRequest) (*teov20220901.DescribeInferenceDomainsResponse, error) {
		assert.Equal(t, "zone-2qtuhspy7cr6", *request.ZoneId)
		assert.Equal(t, "sid-2rugn89bkla9", *request.ServiceId)
		assert.Equal(t, "CreateTime", *request.SortBy)
		assert.Equal(t, "Asc", *request.SortOrder)

		resp := teov20220901.NewDescribeInferenceDomainsResponse()
		resp.Response = &teov20220901.DescribeInferenceDomainsResponseParams{
			TotalCount: func() *int64 { v := int64(1); return &v }(),
			Domains: []*teov20220901.InferenceDomain{
				{
					Domain:     ptrStrInfDomains("infer.example.com"),
					Status:     ptrStrInfDomains("Online"),
					AuthSwitch: ptrStrInfDomains("On"),
					Cname:      ptrStrInfDomains("infer.example.com.dnse1.com"),
					OwnershipVerification: &teov20220901.OwnershipVerification{
						DnsVerification: &teov20220901.DnsVerification{
							Subdomain:   ptrStrInfDomains("_teo-verification"),
							RecordType:  ptrStrInfDomains("CNAME"),
							RecordValue: ptrStrInfDomains("verify.dnse1.com"),
						},
						FileVerification: &teov20220901.FileVerification{
							Path:    ptrStrInfDomains("/.well-known/teo-verification/z12h416twn.txt"),
							Content: ptrStrInfDomains("verification-content"),
						},
						NsVerification: &teov20220901.NsVerification{
							NameServers: []*string{ptrStrInfDomains("ns1.dnse1.com"), ptrStrInfDomains("ns2.dnse1.com")},
						},
					},
					Certificate: &teov20220901.HostsCertificate{
						Host: ptrStrInfDomains("infer.example.com"),
						Mode: ptrStrInfDomains("sslcert"),
						HostCertInfo: []*teov20220901.HostCertInfo{
							{
								CertId:     ptrStrInfDomains("cert-abcdefg"),
								Alias:      ptrStrInfDomains("tf test cert"),
								Type:       ptrStrInfDomains("upload"),
								ExpireTime: ptrStrInfDomains("2026-01-01 00:00:00"),
								DeployTime: ptrStrInfDomains("2024-01-01 00:00:00"),
								SignAlgo:   ptrStrInfDomains("SHA256withRSA"),
								Status:     ptrStrInfDomains("deployed"),
							},
						},
						ApplyType: ptrStrInfDomains("none"),
					},
					CreateTime: ptrStrInfDomains("2024-01-01 00:00:00"),
					UpdateTime: ptrStrInfDomains("2024-06-01 00:00:00"),
				},
			},
			RequestId: ptrStrInfDomains("fake-request-id"),
		}
		return resp, nil
	})

	ds := svcteo.NewTeoInferenceDomainsDataSource().(*svcteo.TeoInferenceDomainsDataSource)
	setTeoInferenceDomainsClient(ds, client)

	sch := teoInferenceDomainsSchema(t)
	raw := buildTeoInferenceDomainsRawValue("zone-2qtuhspy7cr6", "sid-2rugn89bkla9", "CreateTime", "Asc")
	req := newTeoInferenceDomainsReadRequest(sch, raw)
	resp := newTeoInferenceDomainsReadResponse(sch)

	ds.Read(context.Background(), req, resp)

	assert.False(t, resp.Diagnostics.HasError(), "expected no error diagnostics, got: %v", resp.Diagnostics)

	var result svcteo.TeoInferenceDomainsDataSourceModel
	resp.Diagnostics.Append(resp.State.Get(context.Background(), &result)...)
	assert.False(t, resp.Diagnostics.HasError())

	assert.Equal(t, "zone-2qtuhspy7cr6", result.ZoneId.ValueString())
	assert.Equal(t, "sid-2rugn89bkla9", result.ServiceId.ValueString())
	assert.False(t, result.Domains.IsNull())

	var domains []svcteo.TeoInferenceDomainsDomainModel
	resp.Diagnostics.Append(result.Domains.ElementsAs(context.Background(), &domains, false)...)
	assert.False(t, resp.Diagnostics.HasError())
	assert.Len(t, domains, 1)

	// assert every leaf field of the domain
	assert.Equal(t, "infer.example.com", domains[0].Domain.ValueString())
	assert.Equal(t, "Online", domains[0].Status.ValueString())
	assert.Equal(t, "On", domains[0].AuthSwitch.ValueString())
	assert.Equal(t, "infer.example.com.dnse1.com", domains[0].Cname.ValueString())
	assert.Equal(t, "2024-01-01 00:00:00", domains[0].CreateTime.ValueString())
	assert.Equal(t, "2024-06-01 00:00:00", domains[0].UpdateTime.ValueString())

	// ownership_verification nesting
	var ovList []svcteo.TeoInferenceDomainsOwnershipVerificationModel
	resp.Diagnostics.Append(domains[0].OwnershipVerification.ElementsAs(context.Background(), &ovList, false)...)
	assert.False(t, resp.Diagnostics.HasError())
	assert.Len(t, ovList, 1)

	var dnsList []svcteo.TeoInferenceDomainsDnsVerificationModel
	resp.Diagnostics.Append(ovList[0].DnsVerification.ElementsAs(context.Background(), &dnsList, false)...)
	assert.False(t, resp.Diagnostics.HasError())
	assert.Len(t, dnsList, 1)
	assert.Equal(t, "_teo-verification", dnsList[0].Subdomain.ValueString())
	assert.Equal(t, "CNAME", dnsList[0].RecordType.ValueString())
	assert.Equal(t, "verify.dnse1.com", dnsList[0].RecordValue.ValueString())

	var fileList []svcteo.TeoInferenceDomainsFileVerificationModel
	resp.Diagnostics.Append(ovList[0].FileVerification.ElementsAs(context.Background(), &fileList, false)...)
	assert.False(t, resp.Diagnostics.HasError())
	assert.Len(t, fileList, 1)
	assert.Equal(t, "/.well-known/teo-verification/z12h416twn.txt", fileList[0].Path.ValueString())
	assert.Equal(t, "verification-content", fileList[0].Content.ValueString())

	var nsList []svcteo.TeoInferenceDomainsNsVerificationModel
	resp.Diagnostics.Append(ovList[0].NsVerification.ElementsAs(context.Background(), &nsList, false)...)
	assert.False(t, resp.Diagnostics.HasError())
	assert.Len(t, nsList, 1)
	assert.False(t, nsList[0].NameServers.IsNull())
	nameServers := []string{}
	resp.Diagnostics.Append(nsList[0].NameServers.ElementsAs(context.Background(), &nameServers, false)...)
	assert.False(t, resp.Diagnostics.HasError())
	assert.Len(t, nameServers, 2)
	assert.Contains(t, nameServers, "ns1.dnse1.com")
	assert.Contains(t, nameServers, "ns2.dnse1.com")

	// certificate nesting
	var certList []svcteo.TeoInferenceDomainsCertificateModel
	resp.Diagnostics.Append(domains[0].Certificate.ElementsAs(context.Background(), &certList, false)...)
	assert.False(t, resp.Diagnostics.HasError())
	assert.Len(t, certList, 1)
	assert.Equal(t, "infer.example.com", certList[0].Host.ValueString())
	assert.Equal(t, "sslcert", certList[0].Mode.ValueString())
	assert.Equal(t, "none", certList[0].ApplyType.ValueString())

	var hostCertInfos []svcteo.TeoInferenceDomainsHostCertInfoModel
	resp.Diagnostics.Append(certList[0].HostCertInfo.ElementsAs(context.Background(), &hostCertInfos, false)...)
	assert.False(t, resp.Diagnostics.HasError())
	assert.Len(t, hostCertInfos, 1)
	assert.Equal(t, "cert-abcdefg", hostCertInfos[0].CertId.ValueString())
	assert.Equal(t, "tf test cert", hostCertInfos[0].Alias.ValueString())
	assert.Equal(t, "upload", hostCertInfos[0].Type.ValueString())
	assert.Equal(t, "2026-01-01 00:00:00", hostCertInfos[0].ExpireTime.ValueString())
	assert.Equal(t, "2024-01-01 00:00:00", hostCertInfos[0].DeployTime.ValueString())
	assert.Equal(t, "SHA256withRSA", hostCertInfos[0].SignAlgo.ValueString())
	assert.Equal(t, "deployed", hostCertInfos[0].Status.ValueString())
}

// TestTeoInferenceDomains_ReadEmptyResult tests read when the API returns an
// empty Domains list; it must yield an empty list, not an error.
func TestTeoInferenceDomains_ReadEmptyResult(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(client, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeInferenceDomainsWithContext", func(_ context.Context, request *teov20220901.DescribeInferenceDomainsRequest) (*teov20220901.DescribeInferenceDomainsResponse, error) {
		resp := teov20220901.NewDescribeInferenceDomainsResponse()
		resp.Response = &teov20220901.DescribeInferenceDomainsResponseParams{
			Domains:   nil,
			RequestId: ptrStrInfDomains("fake-request-id"),
		}
		return resp, nil
	})

	ds := svcteo.NewTeoInferenceDomainsDataSource().(*svcteo.TeoInferenceDomainsDataSource)
	setTeoInferenceDomainsClient(ds, client)

	sch := teoInferenceDomainsSchema(t)
	raw := buildTeoInferenceDomainsRawValue("zone-2qtuhspy7cr6", "sid-2rugn89bkla9", "", "")
	req := newTeoInferenceDomainsReadRequest(sch, raw)
	resp := newTeoInferenceDomainsReadResponse(sch)

	ds.Read(context.Background(), req, resp)

	assert.False(t, resp.Diagnostics.HasError(), "expected no error diagnostics, got: %v", resp.Diagnostics)

	var result svcteo.TeoInferenceDomainsDataSourceModel
	resp.Diagnostics.Append(resp.State.Get(context.Background(), &result)...)
	assert.False(t, resp.Diagnostics.HasError())
	assert.False(t, result.Domains.IsNull())

	var domains []svcteo.TeoInferenceDomainsDomainModel
	resp.Diagnostics.Append(result.Domains.ElementsAs(context.Background(), &domains, false)...)
	assert.Len(t, domains, 0)
}

// TestTeoInferenceDomains_ReadApiError tests that an API error surfaces as an
// error diagnostic.
func TestTeoInferenceDomains_ReadApiError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(client, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeInferenceDomainsWithContext", func(_ context.Context, request *teov20220901.DescribeInferenceDomainsRequest) (*teov20220901.DescribeInferenceDomainsResponse, error) {
		return nil, errors.New("fake api error")
	})

	ds := svcteo.NewTeoInferenceDomainsDataSource().(*svcteo.TeoInferenceDomainsDataSource)
	setTeoInferenceDomainsClient(ds, client)

	sch := teoInferenceDomainsSchema(t)
	raw := buildTeoInferenceDomainsRawValue("zone-2qtuhspy7cr6", "sid-2rugn89bkla9", "", "")
	req := newTeoInferenceDomainsReadRequest(sch, raw)
	resp := newTeoInferenceDomainsReadResponse(sch)

	ds.Read(context.Background(), req, resp)

	assert.True(t, resp.Diagnostics.HasError(), "expected error diagnostics for API failure")
}

// TestTeoInferenceDomains_ReadNilNestedFields tests a domain whose optional
// nested objects (ownership_verification, certificate) are nil; they must be
// read back as null lists.
func TestTeoInferenceDomains_ReadNilNestedFields(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(client, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeInferenceDomainsWithContext", func(_ context.Context, request *teov20220901.DescribeInferenceDomainsRequest) (*teov20220901.DescribeInferenceDomainsResponse, error) {
		resp := teov20220901.NewDescribeInferenceDomainsResponse()
		resp.Response = &teov20220901.DescribeInferenceDomainsResponseParams{
			Domains: []*teov20220901.InferenceDomain{
				{
					Domain:     ptrStrInfDomains("bare.example.com"),
					Status:     ptrStrInfDomains("Process"),
					AuthSwitch: ptrStrInfDomains("Off"),
				},
			},
			RequestId: ptrStrInfDomains("fake-request-id"),
		}
		return resp, nil
	})

	ds := svcteo.NewTeoInferenceDomainsDataSource().(*svcteo.TeoInferenceDomainsDataSource)
	setTeoInferenceDomainsClient(ds, client)

	sch := teoInferenceDomainsSchema(t)
	raw := buildTeoInferenceDomainsRawValue("zone-2qtuhspy7cr6", "sid-2rugn89bkla9", "", "")
	req := newTeoInferenceDomainsReadRequest(sch, raw)
	resp := newTeoInferenceDomainsReadResponse(sch)

	ds.Read(context.Background(), req, resp)

	assert.False(t, resp.Diagnostics.HasError(), "expected no error diagnostics, got: %v", resp.Diagnostics)

	var result svcteo.TeoInferenceDomainsDataSourceModel
	resp.Diagnostics.Append(resp.State.Get(context.Background(), &result)...)
	assert.False(t, resp.Diagnostics.HasError())

	var domains []svcteo.TeoInferenceDomainsDomainModel
	resp.Diagnostics.Append(result.Domains.ElementsAs(context.Background(), &domains, false)...)
	assert.False(t, resp.Diagnostics.HasError())
	assert.Len(t, domains, 1)
	assert.Equal(t, "bare.example.com", domains[0].Domain.ValueString())
	assert.Equal(t, "Process", domains[0].Status.ValueString())
	assert.Equal(t, "Off", domains[0].AuthSwitch.ValueString())
	assert.True(t, domains[0].OwnershipVerification.IsNull())
	assert.True(t, domains[0].Certificate.IsNull())
	assert.True(t, domains[0].Cname.IsNull())
}
