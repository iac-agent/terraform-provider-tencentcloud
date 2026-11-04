package teo_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	fw_datasource_schema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/assert"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/sharedmeta"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/teo"
)

func ptrStringInferenceDomains(s string) *string { return &s }
func ptrInt64InferenceDomains(n int64) *int64    { return &n }

// inferenceDomainsSchema builds the framework data source schema used to
// construct the raw config value for the ReadRequest.
func inferenceDomainsSchema(t *testing.T) tftypes.Object {
	ds := teo.NewTeoInferenceDomainsV2DataSource()
	schemaResp := &datasource.SchemaResponse{}
	ds.Schema(context.Background(), datasource.SchemaRequest{}, schemaResp)
	if schemaResp.Diagnostics.HasError() {
		t.Fatalf("failed to build data source schema: %v", schemaResp.Diagnostics)
	}
	return schemaResp.Schema.Type().TerraformType(context.Background()).(tftypes.Object)
}

// newInferenceDomainsReadRequest builds a datasource.ReadRequest carrying the
// given raw config values. Optional query fields left empty become null.
func newInferenceDomainsReadRequest(t *testing.T, zoneId, serviceId, sortBy, sortOrder string) datasource.ReadRequest {
	schemaType := inferenceDomainsSchema(t)

	strVal := func(v string) tftypes.Value {
		if v == "" {
			return tftypes.NewValue(tftypes.String, nil)
		}
		return tftypes.NewValue(tftypes.String, v)
	}

	raw := tftypes.NewValue(schemaType, map[string]tftypes.Value{
		"zone_id":    strVal(zoneId),
		"service_id": strVal(serviceId),
		"sort_by":    strVal(sortBy),
		"sort_order": strVal(sortOrder),
		"domains":    tftypes.NewValue(tftypes.List{ElementType: inferenceDomainElemType()}, nil),
	})

	req := datasource.ReadRequest{}
	req.Config.Raw = raw
	req.Config.Schema = mustInferenceDomainsFrameworkSchema(t)
	return req
}

// mustInferenceDomainsFrameworkSchema returns the framework data source schema.
func mustInferenceDomainsFrameworkSchema(t *testing.T) fw_datasource_schema.Schema {
	ds := teo.NewTeoInferenceDomainsV2DataSource()
	schemaResp := &datasource.SchemaResponse{}
	ds.Schema(context.Background(), datasource.SchemaRequest{}, schemaResp)
	if schemaResp.Diagnostics.HasError() {
		t.Fatalf("failed to build data source schema: %v", schemaResp.Diagnostics)
	}
	return schemaResp.Schema
}

// inferenceDomainElemType returns the tftypes.Object type of a domains element.
func inferenceDomainElemType() tftypes.Object {
	return tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"domain":                 tftypes.String,
			"status":                 tftypes.String,
			"auth_switch":            tftypes.String,
			"cname":                  tftypes.String,
			"ownership_verification": ownershipVerificationType(),
			"certificate":            certificateType(),
			"create_time":            tftypes.String,
			"update_time":            tftypes.String,
		},
	}
}

func ownershipVerificationType() tftypes.Object {
	return tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"dns_verification":  dnsVerificationType(),
			"file_verification": fileVerificationType(),
			"ns_verification":   nsVerificationType(),
		},
	}
}

func dnsVerificationType() tftypes.Object {
	return tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"subdomain":    tftypes.String,
			"record_type":  tftypes.String,
			"record_value": tftypes.String,
		},
	}
}

func fileVerificationType() tftypes.Object {
	return tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"path":    tftypes.String,
			"content": tftypes.String,
		},
	}
}

func nsVerificationType() tftypes.Object {
	return tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"name_servers": tftypes.Set{ElementType: tftypes.String},
		},
	}
}

func certificateType() tftypes.Object {
	return tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"host": tftypes.String,
			"mode": tftypes.String,
			"host_cert_info": tftypes.List{
				ElementType: hostCertInfoType(),
			},
			"apply_type": tftypes.String,
		},
	}
}

func hostCertInfoType() tftypes.Object {
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

// runInferenceDomainsRead executes Read with a mocked client and returns the
// response together with the state model decoded from resp.State.
func runInferenceDomainsRead(t *testing.T, req datasource.ReadRequest) (*datasource.ReadResponse, teo.TeoInferenceDomainsV2DataSourceModel) {
	ds := teo.NewTeoInferenceDomainsV2DataSource()
	configureInferenceDomainsDS(ds.(*teo.TeoInferenceDomainsV2DataSource), &connectivity.TencentCloudClient{})

	resp := &datasource.ReadResponse{
		State: tfsdk.State{
			Schema: mustInferenceDomainsFrameworkSchema(t),
		},
	}
	ds.Read(context.Background(), req, resp)

	var model teo.TeoInferenceDomainsV2DataSourceModel
	if !resp.State.Raw.IsNull() {
		diags := resp.State.Get(context.Background(), &model)
		assert.False(t, diags.HasError(), "state get should not error: %v", diags)
	}
	return resp, model
}

// configureInferenceDomainsDS injects a mock connectivity client into the
// data source via the Configure method.
func configureInferenceDomainsDS(ds *teo.TeoInferenceDomainsV2DataSource, client *connectivity.TencentCloudClient) {
	meta := &sharedmeta.ProviderMeta{Client: client}
	configureReq := datasource.ConfigureRequest{ProviderData: meta}
	configureResp := &datasource.ConfigureResponse{}
	ds.Configure(context.Background(), configureReq, configureResp)
}

// go test ./tencentcloud/services/teo/ -run "TestTeoInferenceDomainsV2" -v -count=1 -gcflags="all=-l"

// TestTeoInferenceDomainsV2DataSource_ReadFull verifies a full Read with one
// domain carrying every nested object, and asserts every leaf field.
func TestTeoInferenceDomainsV2DataSource_ReadFull(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(client, "UseTeoClient", teoClient)

	var gotRequest *teov20220901.DescribeInferenceDomainsRequest
	patches.ApplyMethodFunc(teoClient, "DescribeInferenceDomainsWithContext", func(_ context.Context, request *teov20220901.DescribeInferenceDomainsRequest) (*teov20220901.DescribeInferenceDomainsResponse, error) {
		gotRequest = request
		resp := teov20220901.NewDescribeInferenceDomainsResponse()
		resp.Response = &teov20220901.DescribeInferenceDomainsResponseParams{
			TotalCount: ptrInt64InferenceDomains(1),
			Domains: []*teov20220901.InferenceDomain{
				{
					Domain:     ptrStringInferenceDomains("infer.example.com"),
					Status:     ptrStringInferenceDomains("Online"),
					AuthSwitch: ptrStringInferenceDomains("On"),
					Cname:      ptrStringInferenceDomains("infer.example.com.cdn.dnsv1.com"),
					OwnershipVerification: &teov20220901.OwnershipVerification{
						DnsVerification: &teov20220901.DnsVerification{
							Subdomain:   ptrStringInferenceDomains("_teo_verify"),
							RecordType:  ptrStringInferenceDomains("TXT"),
							RecordValue: ptrStringInferenceDomains("verify-token"),
						},
						FileVerification: &teov20220901.FileVerification{
							Path:    ptrStringInferenceDomains("/.well-known/teo-verification/z12h416twn.txt"),
							Content: ptrStringInferenceDomains("file-content"),
						},
						NsVerification: &teov20220901.NsVerification{
							NameServers: []*string{
								ptrStringInferenceDomains("ns1.teo.dnsv1.com"),
								ptrStringInferenceDomains("ns2.teo.dnsv1.com"),
							},
						},
					},
					Certificate: &teov20220901.HostsCertificate{
						Host: ptrStringInferenceDomains("infer.example.com"),
						Mode: ptrStringInferenceDomains("eofreecert"),
						HostCertInfo: []*teov20220901.HostCertInfo{
							{
								CertId:     ptrStringInferenceDomains("cert-abc123"),
								Alias:      ptrStringInferenceDomains("my cert"),
								Type:       ptrStringInferenceDomains("managed"),
								ExpireTime: ptrStringInferenceDomains("2026-01-01 00:00:00"),
								DeployTime: ptrStringInferenceDomains("2025-01-01 00:00:00"),
								SignAlgo:   ptrStringInferenceDomains("RSA"),
								Status:     ptrStringInferenceDomains("deployed"),
							},
						},
						ApplyType: ptrStringInferenceDomains("apply"),
					},
					CreateTime: ptrStringInferenceDomains("2025-06-01 10:00:00"),
					UpdateTime: ptrStringInferenceDomains("2025-06-02 10:00:00"),
				},
			},
			RequestId: ptrStringInferenceDomains("fake-request-id"),
		}
		return resp, nil
	})

	req := newInferenceDomainsReadRequest(t, "zone-2qtuhspy7cr6", "sid-123456", "UpdateTime", "Asc")
	resp, model := runInferenceDomainsRead(t, req)

	assert.False(t, resp.Diagnostics.HasError(), "Read should have no error diagnostics: %v", resp.Diagnostics)

	// request wiring
	assert.NotNil(t, gotRequest)
	assert.Equal(t, "zone-2qtuhspy7cr6", *gotRequest.ZoneId)
	assert.Equal(t, "sid-123456", *gotRequest.ServiceId)
	assert.Equal(t, "UpdateTime", *gotRequest.SortBy)
	assert.Equal(t, "Asc", *gotRequest.SortOrder)

	// state echo of query inputs
	assert.Equal(t, "zone-2qtuhspy7cr6", model.ZoneId.ValueString())
	assert.Equal(t, "sid-123456", model.ServiceId.ValueString())
	assert.Equal(t, "UpdateTime", model.SortBy.ValueString())
	assert.Equal(t, "Asc", model.SortOrder.ValueString())

	// domains list
	assert.False(t, model.Domains.IsNull())
	items := unmarshalInferenceDomainItems(t, model.Domains)
	assert.Len(t, items, 1)

	item := items[0].Attributes()
	assert.Equal(t, "infer.example.com", item["domain"].(types.String).ValueString())
	assert.Equal(t, "Online", item["status"].(types.String).ValueString())
	assert.Equal(t, "On", item["auth_switch"].(types.String).ValueString())
	assert.Equal(t, "infer.example.com.cdn.dnsv1.com", item["cname"].(types.String).ValueString())
	assert.Equal(t, "2025-06-01 10:00:00", item["create_time"].(types.String).ValueString())
	assert.Equal(t, "2025-06-02 10:00:00", item["update_time"].(types.String).ValueString())

	// ownership verification
	ov := item["ownership_verification"].(types.Object)
	assert.False(t, ov.IsNull())
	ovAttrs := ov.Attributes()
	dns := ovAttrs["dns_verification"].(types.Object).Attributes()
	assert.Equal(t, "_teo_verify", dns["subdomain"].(types.String).ValueString())
	assert.Equal(t, "TXT", dns["record_type"].(types.String).ValueString())
	assert.Equal(t, "verify-token", dns["record_value"].(types.String).ValueString())

	file := ovAttrs["file_verification"].(types.Object).Attributes()
	assert.Equal(t, "/.well-known/teo-verification/z12h416twn.txt", file["path"].(types.String).ValueString())
	assert.Equal(t, "file-content", file["content"].(types.String).ValueString())

	ns := ovAttrs["ns_verification"].(types.Object).Attributes()
	nsServers, ok := ns["name_servers"].(types.Set)
	assert.True(t, ok, "name_servers should be a set")
	var nsValues []types.String
	assert.False(t, nsServers.ElementsAs(context.Background(), &nsValues, false).HasError(),
		"name_servers ElementsAs should not error")
	assert.Len(t, nsValues, 2)
	got := []string{nsValues[0].ValueString(), nsValues[1].ValueString()}
	assert.ElementsMatch(t, []string{"ns1.teo.dnsv1.com", "ns2.teo.dnsv1.com"}, got)

	// certificate
	cert := item["certificate"].(types.Object)
	assert.False(t, cert.IsNull())
	certAttrs := cert.Attributes()
	assert.Equal(t, "infer.example.com", certAttrs["host"].(types.String).ValueString())
	assert.Equal(t, "eofreecert", certAttrs["mode"].(types.String).ValueString())
	assert.Equal(t, "apply", certAttrs["apply_type"].(types.String).ValueString())

	var hciItems []types.Object
	hciList, ok := certAttrs["host_cert_info"].(types.List)
	assert.True(t, ok, "host_cert_info should be a list")
	assert.False(t, hciList.ElementsAs(context.Background(), &hciItems, false).HasError(),
		"host_cert_info ElementsAs should not error")
	assert.Len(t, hciItems, 1)
	hciAttrs := hciItems[0].Attributes()
	assert.Equal(t, "cert-abc123", hciAttrs["cert_id"].(types.String).ValueString())
	assert.Equal(t, "my cert", hciAttrs["alias"].(types.String).ValueString())
	assert.Equal(t, "managed", hciAttrs["type"].(types.String).ValueString())
	assert.Equal(t, "2026-01-01 00:00:00", hciAttrs["expire_time"].(types.String).ValueString())
	assert.Equal(t, "2025-01-01 00:00:00", hciAttrs["deploy_time"].(types.String).ValueString())
	assert.Equal(t, "RSA", hciAttrs["sign_algo"].(types.String).ValueString())
	assert.Equal(t, "deployed", hciAttrs["status"].(types.String).ValueString())
}

// unmarshalInferenceDomainItems decodes a domains types.List into typed
// object elements whose attributes can be inspected leaf by leaf.
func unmarshalInferenceDomainItems(t *testing.T, list types.List) []types.Object {
	var items []types.Object
	assert.False(t, list.ElementsAs(context.Background(), &items, false).HasError(),
		"domains ElementsAs should not error")
	return items
}

// TestTeoInferenceDomainsV2DataSource_ReadEmptyList verifies that an empty
// result list is a valid business state and lands in state as an empty list.
func TestTeoInferenceDomainsV2DataSource_ReadEmptyList(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(client, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeInferenceDomainsWithContext", func(_ context.Context, request *teov20220901.DescribeInferenceDomainsRequest) (*teov20220901.DescribeInferenceDomainsResponse, error) {
		resp := teov20220901.NewDescribeInferenceDomainsResponse()
		resp.Response = &teov20220901.DescribeInferenceDomainsResponseParams{
			TotalCount: ptrInt64InferenceDomains(0),
			Domains:    []*teov20220901.InferenceDomain{},
			RequestId:  ptrStringInferenceDomains("fake-request-id"),
		}
		return resp, nil
	})

	req := newInferenceDomainsReadRequest(t, "zone-2qtuhspy7cr6", "sid-123456", "", "")
	resp, model := runInferenceDomainsRead(t, req)

	assert.False(t, resp.Diagnostics.HasError(), "Read should have no error diagnostics: %v", resp.Diagnostics)
	assert.False(t, model.Domains.IsNull())
	items := unmarshalInferenceDomainItems(t, model.Domains)
	assert.Len(t, items, 0)
}

// TestTeoInferenceDomainsV2DataSource_ReadNilSubObjects verifies that nil
// nested objects (ownership verification / certificate) land as null objects
// instead of crashing or being dropped.
func TestTeoInferenceDomainsV2DataSource_ReadNilSubObjects(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(client, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeInferenceDomainsWithContext", func(_ context.Context, request *teov20220901.DescribeInferenceDomainsRequest) (*teov20220901.DescribeInferenceDomainsResponse, error) {
		resp := teov20220901.NewDescribeInferenceDomainsResponse()
		resp.Response = &teov20220901.DescribeInferenceDomainsResponseParams{
			TotalCount: ptrInt64InferenceDomains(1),
			Domains: []*teov20220901.InferenceDomain{
				{
					Domain: ptrStringInferenceDomains("infer.example.com"),
					Status: ptrStringInferenceDomains("Init"),
					// OwnershipVerification: nil
					// Certificate: nil
				},
			},
			RequestId: ptrStringInferenceDomains("fake-request-id"),
		}
		return resp, nil
	})

	req := newInferenceDomainsReadRequest(t, "zone-2qtuhspy7cr6", "sid-123456", "", "")
	resp, model := runInferenceDomainsRead(t, req)

	assert.False(t, resp.Diagnostics.HasError(), "Read should have no error diagnostics: %v", resp.Diagnostics)
	items := unmarshalInferenceDomainItems(t, model.Domains)
	assert.Len(t, items, 1)

	ov := items[0].Attributes()["ownership_verification"].(types.Object)
	assert.True(t, ov.IsNull(), "ownership_verification should be null when cloud returns nil")

	cert := items[0].Attributes()["certificate"].(types.Object)
	assert.True(t, cert.IsNull(), "certificate should be null when cloud returns nil")
}

// TestTeoInferenceDomainsV2DataSource_ReadPagination verifies that Read walks
// all pages when TotalCount exceeds a single page.
func TestTeoInferenceDomainsV2DataSource_ReadPagination(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(client, "UseTeoClient", teoClient)

	callCount := 0
	patches.ApplyMethodFunc(teoClient, "DescribeInferenceDomainsWithContext", func(_ context.Context, request *teov20220901.DescribeInferenceDomainsRequest) (*teov20220901.DescribeInferenceDomainsResponse, error) {
		callCount++
		resp := teov20220901.NewDescribeInferenceDomainsResponse()
		var domains []*teov20220901.InferenceDomain
		if callCount == 1 {
			for i := 0; i < 2; i++ {
				domains = append(domains, &teov20220901.InferenceDomain{
					Domain: ptrStringInferenceDomains(fmt.Sprintf("infer-%d.example.com", i)),
					Status: ptrStringInferenceDomains("Online"),
				})
			}
		} else if callCount == 2 {
			domains = append(domains, &teov20220901.InferenceDomain{
				Domain: ptrStringInferenceDomains("infer-2.example.com"),
				Status: ptrStringInferenceDomains("Online"),
			})
		}
		resp.Response = &teov20220901.DescribeInferenceDomainsResponseParams{
			TotalCount: ptrInt64InferenceDomains(3),
			Domains:    domains,
			RequestId:  ptrStringInferenceDomains("fake-request-id"),
		}
		return resp, nil
	})

	req := newInferenceDomainsReadRequest(t, "zone-2qtuhspy7cr6", "sid-123456", "", "")
	resp, model := runInferenceDomainsRead(t, req)

	assert.False(t, resp.Diagnostics.HasError(), "Read should have no error diagnostics: %v", resp.Diagnostics)
	assert.Equal(t, 2, callCount, "two pages should be fetched")
	items := unmarshalInferenceDomainItems(t, model.Domains)
	assert.Len(t, items, 3)
}

// TestTeoInferenceDomainsV2DataSource_ReadApiError verifies that an API error
// surfaces as an error diagnostic.
func TestTeoInferenceDomainsV2DataSource_ReadApiError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(client, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeInferenceDomainsWithContext", func(_ context.Context, request *teov20220901.DescribeInferenceDomainsRequest) (*teov20220901.DescribeInferenceDomainsResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=AuthFailure, Message=authentication failed")
	})

	req := newInferenceDomainsReadRequest(t, "zone-2qtuhspy7cr6", "sid-123456", "", "")
	resp, _ := runInferenceDomainsRead(t, req)

	assert.True(t, resp.Diagnostics.HasError(), "API error should surface as error diagnostics")
}

// TestTeoInferenceDomainsV2DataSource_Metadata validates the data source type name.
func TestTeoInferenceDomainsV2DataSource_Metadata(t *testing.T) {
	ds := teo.NewTeoInferenceDomainsV2DataSource()
	req := datasource.MetadataRequest{}
	resp := &datasource.MetadataResponse{}

	ds.Metadata(context.Background(), req, resp)

	if resp.TypeName != "tencentcloud_teo_inference_domains_v2" {
		t.Errorf("Metadata.TypeName = %q, want %q", resp.TypeName, "tencentcloud_teo_inference_domains_v2")
	}
}

// TestTeoInferenceDomainsV2DataSource_Factory validates the registry factory.
func TestTeoInferenceDomainsV2DataSource_Factory(t *testing.T) {
	got := teo.NewTeoInferenceDomainsV2DataSource()
	if got == nil {
		t.Fatal("NewTeoInferenceDomainsV2DataSource returned nil")
	}
	if _, ok := got.(*teo.TeoInferenceDomainsV2DataSource); !ok {
		t.Fatalf("factory returned %T, want *TeoInferenceDomainsV2DataSource", got)
	}
}

// TestTeoInferenceDomainsV2DataSource_ImplementsInterface validates that the
// data source implements datasource.DataSourceWithConfigure.
func TestTeoInferenceDomainsV2DataSource_ImplementsInterface(t *testing.T) {
	var _ datasource.DataSourceWithConfigure = &teo.TeoInferenceDomainsV2DataSource{}
}

// TestTeoInferenceDomainsV2DataSource_ConfigureNilProviderData ensures the
// pre-Configure invocation pattern (nil ProviderData) is silently no-op.
func TestTeoInferenceDomainsV2DataSource_ConfigureNilProviderData(t *testing.T) {
	ds := &teo.TeoInferenceDomainsV2DataSource{}
	req := datasource.ConfigureRequest{}
	resp := &datasource.ConfigureResponse{}

	ds.Configure(context.Background(), req, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("Configure with nil ProviderData returned diagnostics: %v", resp.Diagnostics)
	}
}
