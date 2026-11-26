package alb_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/assert"

	albv20251030 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/alb/v20251030"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/sharedmeta"
	svalb "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/alb"
)

func albListenerCertificatesStr(s string) *string { return &s }
func albListenerCertificatesBool(v bool) *bool    { return &v }
func albListenerCertificatesUint64(v uint64) *uint64 {
	return &v
}

// albListenerCertificatesSchema returns the framework data source schema
// used to build the raw config/state value.
func albListenerCertificatesSchema(t *testing.T) schema.Schema {
	d := svalb.NewAlbListenerCertificatesDataSource()
	resp := &datasource.SchemaResponse{}
	d.Schema(context.Background(), datasource.SchemaRequest{}, resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("failed to build data source schema: %v", resp.Diagnostics)
	}

	return resp.Schema
}

// albListenerCertificatesRootType builds the tftypes.Object type for the
// data source root, mirroring every schema attribute.
func albListenerCertificatesRootType() tftypes.Object {
	return tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"id":               tftypes.String,
			"certificate_type": tftypes.String,
			"listener_id":      tftypes.String,
			"load_balancer_id": tftypes.String,
			"certificates":     tftypes.List{ElementType: albListenerCertificatesItemType()},
		},
	}
}

// albListenerCertificatesItemType builds the tftypes.Object type for a
// certificates element.
func albListenerCertificatesItemType() tftypes.Object {
	return tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"associated_time":  tftypes.String,
			"certificate_id":   tftypes.String,
			"certificate_type": tftypes.String,
			"is_default":       tftypes.Bool,
			"status":           tftypes.String,
		},
	}
}

// buildAlbListenerCertificatesRawValue builds the raw config value with the
// given query inputs; every other attribute is unset (nil).
func buildAlbListenerCertificatesRawValue(certificateType string, listenerId string, loadBalancerId string) tftypes.Value {
	attrs := map[string]tftypes.Value{
		"id":               tftypes.NewValue(tftypes.String, nil),
		"certificate_type": tftypes.NewValue(tftypes.String, certificateType),
		"listener_id":      tftypes.NewValue(tftypes.String, listenerId),
		"load_balancer_id": tftypes.NewValue(tftypes.String, loadBalancerId),
		"certificates":     tftypes.NewValue(tftypes.List{ElementType: albListenerCertificatesItemType()}, nil),
	}

	return tftypes.NewValue(albListenerCertificatesRootType(), attrs)
}

// setAlbListenerCertificatesClient injects a mock connectivity client into
// the data source instance via the Configure method.
func setAlbListenerCertificatesClient(d *svalb.AlbListenerCertificatesDataSource, client *connectivity.TencentCloudClient) {
	meta := &sharedmeta.ProviderMeta{Client: client}
	configureReq := datasource.ConfigureRequest{ProviderData: meta}
	configureResp := &datasource.ConfigureResponse{}
	d.Configure(context.Background(), configureReq, configureResp)
}

func newAlbListenerCertificatesReadRequest(sch schema.Schema, raw tftypes.Value) datasource.ReadRequest {
	req := datasource.ReadRequest{}
	req.Config.Raw = raw
	req.Config.Schema = sch

	return req
}

func newAlbListenerCertificatesReadResponse(sch schema.Schema) *datasource.ReadResponse {
	resp := &datasource.ReadResponse{}
	resp.State.Schema = sch

	return resp
}

// TestAlbListenerCertificates_MetadataAndSchema validates metadata and
// schema definition.
func TestAlbListenerCertificates_MetadataAndSchema(t *testing.T) {
	d := svalb.NewAlbListenerCertificatesDataSource()

	metaResp := &datasource.MetadataResponse{}
	d.Metadata(context.Background(), datasource.MetadataRequest{}, metaResp)
	assert.Equal(t, "tencentcloud_alb_listener_certificates", metaResp.TypeName)

	sch := albListenerCertificatesSchema(t)
	attrs := sch.Attributes

	assert.True(t, attrs["certificate_type"].(schema.StringAttribute).Required)
	assert.True(t, attrs["listener_id"].(schema.StringAttribute).Required)
	assert.True(t, attrs["load_balancer_id"].(schema.StringAttribute).Required)
	assert.True(t, attrs["certificates"].(schema.ListNestedAttribute).Computed)
	assert.True(t, attrs["id"].(schema.StringAttribute).Computed)

	// Paging knobs and total-count aggregates are NEVER exposed.
	assert.NotContains(t, attrs, "limit")
	assert.NotContains(t, attrs, "offset")
	assert.NotContains(t, attrs, "total_count")
	assert.NotContains(t, attrs, "next_token")
	assert.NotContains(t, attrs, "max_results")
}

// TestAlbListenerCertificates_ReadBasic tests a successful read with one
// record and asserts every leaf field is read back.
func TestAlbListenerCertificates_ReadBasic(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	albClient := &albv20251030.Client{}
	patches.ApplyMethodReturn(client, "UseAlbV20251030Client", albClient)
	patches.ApplyMethodFunc(albClient, "DescribeListenerCertificates", func(request *albv20251030.DescribeListenerCertificatesRequest) (*albv20251030.DescribeListenerCertificatesResponse, error) {
		if request.CertificateType == nil || *request.CertificateType != "SVR" {
			t.Errorf("expected certificate_type to be passed to the API request")
		}
		if request.ListenerId == nil || *request.ListenerId != "lst-12345678" {
			t.Errorf("expected listener_id to be passed to the API request")
		}
		if request.LoadBalancerId == nil || *request.LoadBalancerId != "alb-12345678" {
			t.Errorf("expected load_balancer_id to be passed to the API request")
		}

		return &albv20251030.DescribeListenerCertificatesResponse{Response: &albv20251030.DescribeListenerCertificatesResponseParams{
			Certificates: []*albv20251030.CertificateInfo{
				{
					AssociatedTime:  albListenerCertificatesStr("2024-01-01 00:00:00"),
					CertificateId:   albListenerCertificatesStr("cert-12345678"),
					CertificateType: albListenerCertificatesStr("SVR"),
					IsDefault:       albListenerCertificatesBool(true),
					Status:          albListenerCertificatesStr("Associated"),
				},
			},
			MaxResults: albListenerCertificatesUint64(100),
			TotalCount: albListenerCertificatesUint64(1),
			RequestId:  albListenerCertificatesStr("fake-request-id"),
		}}, nil
	})

	ds := svalb.NewAlbListenerCertificatesDataSource().(*svalb.AlbListenerCertificatesDataSource)
	setAlbListenerCertificatesClient(ds, client)

	sch := albListenerCertificatesSchema(t)
	raw := buildAlbListenerCertificatesRawValue("SVR", "lst-12345678", "alb-12345678")
	req := newAlbListenerCertificatesReadRequest(sch, raw)
	resp := newAlbListenerCertificatesReadResponse(sch)

	ds.Read(context.Background(), req, resp)

	assert.False(t, resp.Diagnostics.HasError(), "expected no error diagnostics, got: %v", resp.Diagnostics)

	var result svalb.AlbListenerCertificatesDataSourceModel
	resp.Diagnostics.Append(resp.State.Get(context.Background(), &result)...)
	assert.False(t, resp.Diagnostics.HasError())

	// The state id is a random token.
	assert.NotEmpty(t, result.Id.ValueString())

	assert.False(t, result.Certificates.IsNull())

	var items []svalb.AlbListenerCertificatesItemModel
	resp.Diagnostics.Append(result.Certificates.ElementsAs(context.Background(), &items, false)...)
	assert.False(t, resp.Diagnostics.HasError())
	assert.Len(t, items, 1)

	assert.Equal(t, "2024-01-01 00:00:00", items[0].AssociatedTime.ValueString())
	assert.Equal(t, "cert-12345678", items[0].CertificateId.ValueString())
	assert.Equal(t, "SVR", items[0].CertificateType.ValueString())
	assert.True(t, items[0].IsDefault.ValueBool())
	assert.Equal(t, "Associated", items[0].Status.ValueString())
}

// TestAlbListenerCertificates_ReadEmpty covers an empty result: a legitimate
// answer, not a failure — the data source must succeed with an empty list.
func TestAlbListenerCertificates_ReadEmpty(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	albClient := &albv20251030.Client{}
	patches.ApplyMethodReturn(client, "UseAlbV20251030Client", albClient)
	patches.ApplyMethodFunc(albClient, "DescribeListenerCertificates", func(request *albv20251030.DescribeListenerCertificatesRequest) (*albv20251030.DescribeListenerCertificatesResponse, error) {
		return &albv20251030.DescribeListenerCertificatesResponse{Response: &albv20251030.DescribeListenerCertificatesResponseParams{}}, nil
	})

	ds := svalb.NewAlbListenerCertificatesDataSource().(*svalb.AlbListenerCertificatesDataSource)
	setAlbListenerCertificatesClient(ds, client)

	sch := albListenerCertificatesSchema(t)
	raw := buildAlbListenerCertificatesRawValue("CA", "lst-12345678", "alb-12345678")
	req := newAlbListenerCertificatesReadRequest(sch, raw)
	resp := newAlbListenerCertificatesReadResponse(sch)

	ds.Read(context.Background(), req, resp)

	assert.False(t, resp.Diagnostics.HasError(), "expected no error diagnostics, got: %v", resp.Diagnostics)

	var result svalb.AlbListenerCertificatesDataSourceModel
	resp.Diagnostics.Append(resp.State.Get(context.Background(), &result)...)
	assert.False(t, resp.Diagnostics.HasError())

	assert.False(t, result.Certificates.IsNull())

	var items []svalb.AlbListenerCertificatesItemModel
	resp.Diagnostics.Append(result.Certificates.ElementsAs(context.Background(), &items, false)...)
	assert.False(t, resp.Diagnostics.HasError())
	assert.Len(t, items, 0)
}

// TestAlbListenerCertificates_ReadApiError covers an API-level failure: the
// error must surface through diagnostics instead of silently persisting
// empty state.
func TestAlbListenerCertificates_ReadApiError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	albClient := &albv20251030.Client{}
	patches.ApplyMethodReturn(client, "UseAlbV20251030Client", albClient)
	patches.ApplyMethodFunc(albClient, "DescribeListenerCertificates", func(request *albv20251030.DescribeListenerCertificatesRequest) (*albv20251030.DescribeListenerCertificatesResponse, error) {
		return nil, fmt.Errorf("mocked API error")
	})

	ds := svalb.NewAlbListenerCertificatesDataSource().(*svalb.AlbListenerCertificatesDataSource)
	setAlbListenerCertificatesClient(ds, client)

	sch := albListenerCertificatesSchema(t)
	raw := buildAlbListenerCertificatesRawValue("SVR", "lst-12345678", "alb-12345678")
	req := newAlbListenerCertificatesReadRequest(sch, raw)
	resp := newAlbListenerCertificatesReadResponse(sch)

	ds.Read(context.Background(), req, resp)

	assert.True(t, resp.Diagnostics.HasError(), "expected error diagnostics for API failure")
}

// TestAlbListenerCertificatesService_ReadAllPages verifies that the service
// layer loops through every page with MaxResults/NextToken and returns the
// complete set (the total count stays internal).
func TestAlbListenerCertificatesService_ReadAllPages(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	albClient := &albv20251030.Client{}
	patches.ApplyMethodReturn(client, "UseAlbV20251030Client", albClient)

	callCount := 0
	patches.ApplyMethodFunc(albClient, "DescribeListenerCertificates", func(request *albv20251030.DescribeListenerCertificatesRequest) (*albv20251030.DescribeListenerCertificatesResponse, error) {
		callCount++
		// The service must drive pagination itself with the page size 100.
		if request.MaxResults == nil || *request.MaxResults != 100 {
			t.Errorf("expected MaxResults 100 on call %d, got %v", callCount, request.MaxResults)
		}
		var nextToken string
		if request.NextToken != nil {
			nextToken = *request.NextToken
		}
		switch callCount {
		case 1:
			if nextToken != "" {
				t.Errorf("expected empty next_token on first call, got %q", nextToken)
			}
			certificates := make([]*albv20251030.CertificateInfo, 0, 100)
			for i := 0; i < 100; i++ {
				certificates = append(certificates, &albv20251030.CertificateInfo{
					CertificateId: albListenerCertificatesStr(fmt.Sprintf("cert-page1-%d", i)),
					Status:        albListenerCertificatesStr("Associated"),
				})
			}

			return &albv20251030.DescribeListenerCertificatesResponse{Response: &albv20251030.DescribeListenerCertificatesResponseParams{
				Certificates: certificates,
				NextToken:    albListenerCertificatesStr("next-token-page1"),
				TotalCount:   albListenerCertificatesUint64(101),
				RequestId:    albListenerCertificatesStr("fake-request-id"),
			}}, nil
		case 2:
			if nextToken != "next-token-page1" {
				t.Errorf("expected next_token from page 1 on second call, got %q", nextToken)
			}

			return &albv20251030.DescribeListenerCertificatesResponse{Response: &albv20251030.DescribeListenerCertificatesResponseParams{
				Certificates: []*albv20251030.CertificateInfo{
					{CertificateId: albListenerCertificatesStr("cert-page2-0"), Status: albListenerCertificatesStr("Associated")},
				},
				TotalCount: albListenerCertificatesUint64(101),
				RequestId:  albListenerCertificatesStr("fake-request-id"),
			}}, nil
		default:
			t.Errorf("unexpected extra API call %d", callCount)

			return &albv20251030.DescribeListenerCertificatesResponse{}, nil
		}
	})

	service := svalb.NewAlbService(client)
	describeListenerCertificatesRequest := albv20251030.NewDescribeListenerCertificatesRequest()
	describeListenerCertificatesRequest.CertificateType = albListenerCertificatesStr("SVR")
	describeListenerCertificatesRequest.ListenerId = albListenerCertificatesStr("lst-12345678")
	describeListenerCertificatesRequest.LoadBalancerId = albListenerCertificatesStr("alb-12345678")

	certificates, err := service.DescribeListenerCertificates(context.Background(), describeListenerCertificatesRequest)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if callCount != 2 {
		t.Errorf("expected 2 API calls, got %d", callCount)
	}
	if len(certificates) != 101 {
		t.Fatalf("expected 101 certificates across pages, got %d", len(certificates))
	}
	if certificates[0].CertificateId == nil || *certificates[0].CertificateId != "cert-page1-0" {
		t.Errorf("unexpected first certificate id: %v", certificates[0].CertificateId)
	}
	if certificates[100].CertificateId == nil || *certificates[100].CertificateId != "cert-page2-0" {
		t.Errorf("unexpected last certificate id: %v", certificates[100].CertificateId)
	}
}
