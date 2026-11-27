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
	svcat "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/alb"
)

func albListenerCertificatesStr(s string) *string    { return &s }
func albListenerCertificatesBool(b bool) *bool       { return &b }
func albListenerCertificatesUint64(v uint64) *uint64 { return &v }

func albListenerCertificatesSchema(t *testing.T) schema.Schema {
	d := svcat.NewAlbListenerCertificatesDataSource()
	resp := &datasource.SchemaResponse{}
	d.Schema(context.Background(), datasource.SchemaRequest{}, resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("failed to build data source schema: %v", resp.Diagnostics)
	}

	return resp.Schema
}

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

func buildAlbListenerCertificatesRawValue(certificateType, listenerId, loadBalancerId string) tftypes.Value {
	attrs := map[string]tftypes.Value{
		"id":               tftypes.NewValue(tftypes.String, nil),
		"certificate_type": tftypes.NewValue(tftypes.String, certificateType),
		"listener_id":      tftypes.NewValue(tftypes.String, listenerId),
		"load_balancer_id": tftypes.NewValue(tftypes.String, loadBalancerId),
		"certificates":     tftypes.NewValue(tftypes.List{ElementType: albListenerCertificatesItemType()}, nil),
	}

	return tftypes.NewValue(albListenerCertificatesRootType(), attrs)
}

func setAlbListenerCertificatesClient(d *svcat.AlbListenerCertificatesDataSource, client *connectivity.TencentCloudClient) {
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

func TestAlbListenerCertificates_MetadataAndSchema(t *testing.T) {
	d := svcat.NewAlbListenerCertificatesDataSource()

	metaResp := &datasource.MetadataResponse{}
	d.Metadata(context.Background(), datasource.MetadataRequest{}, metaResp)
	assert.Equal(t, "tencentcloud_alb_listener_certificates", metaResp.TypeName)

	attrs := albListenerCertificatesSchema(t).Attributes
	assert.True(t, attrs["certificate_type"].(schema.StringAttribute).Required)
	assert.True(t, attrs["listener_id"].(schema.StringAttribute).Required)
	assert.True(t, attrs["load_balancer_id"].(schema.StringAttribute).Required)
	assert.True(t, attrs["certificates"].(schema.ListNestedAttribute).Computed)

	assert.NotContains(t, attrs, "limit")
	assert.NotContains(t, attrs, "offset")
	assert.NotContains(t, attrs, "next_token")
	assert.NotContains(t, attrs, "total_count")
}

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
		if request.ListenerId == nil || *request.ListenerId != "lst-d9p3k7wa" {
			t.Errorf("expected listener_id to be passed to the API request")
		}
		if request.LoadBalancerId == nil || *request.LoadBalancerId != "alb-f8q2xk9m" {
			t.Errorf("expected load_balancer_id to be passed to the API request")
		}

		return &albv20251030.DescribeListenerCertificatesResponse{Response: &albv20251030.DescribeListenerCertificatesResponseParams{
			Certificates: []*albv20251030.CertificateInfo{
				{
					AssociatedTime:  albListenerCertificatesStr("2024-06-12T10:34:56Z"),
					CertificateId:   albListenerCertificatesStr("cert-12345"),
					CertificateType: albListenerCertificatesStr("SVR"),
					IsDefault:       albListenerCertificatesBool(true),
					Status:          albListenerCertificatesStr("Associated"),
				},
			},
			MaxResults: albListenerCertificatesUint64(100),
			RequestId:  albListenerCertificatesStr("fake-request-id"),
		}}, nil
	})

	ds := svcat.NewAlbListenerCertificatesDataSource().(*svcat.AlbListenerCertificatesDataSource)
	setAlbListenerCertificatesClient(ds, client)

	sch := albListenerCertificatesSchema(t)
	raw := buildAlbListenerCertificatesRawValue("SVR", "lst-d9p3k7wa", "alb-f8q2xk9m")
	req := newAlbListenerCertificatesReadRequest(sch, raw)
	resp := newAlbListenerCertificatesReadResponse(sch)

	ds.Read(context.Background(), req, resp)

	assert.False(t, resp.Diagnostics.HasError(), "expected no error diagnostics, got: %v", resp.Diagnostics)

	var result svcat.AlbListenerCertificatesDataSourceModel
	resp.Diagnostics.Append(resp.State.Get(context.Background(), &result)...)
	assert.False(t, resp.Diagnostics.HasError())
	assert.NotEmpty(t, result.Id.ValueString())
	assert.False(t, result.Certificates.IsNull())

	var items []svcat.AlbListenerCertificateItemModel
	resp.Diagnostics.Append(result.Certificates.ElementsAs(context.Background(), &items, false)...)
	assert.False(t, resp.Diagnostics.HasError())
	assert.Len(t, items, 1)
	assert.Equal(t, "2024-06-12T10:34:56Z", items[0].AssociatedTime.ValueString())
	assert.Equal(t, "cert-12345", items[0].CertificateId.ValueString())
	assert.Equal(t, "SVR", items[0].CertificateType.ValueString())
	assert.True(t, items[0].IsDefault.ValueBool())
	assert.Equal(t, "Associated", items[0].Status.ValueString())
}

func TestAlbListenerCertificates_ReadEmpty(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	albClient := &albv20251030.Client{}
	patches.ApplyMethodReturn(client, "UseAlbV20251030Client", albClient)
	patches.ApplyMethodFunc(albClient, "DescribeListenerCertificates", func(request *albv20251030.DescribeListenerCertificatesRequest) (*albv20251030.DescribeListenerCertificatesResponse, error) {
		return &albv20251030.DescribeListenerCertificatesResponse{Response: &albv20251030.DescribeListenerCertificatesResponseParams{}}, nil
	})

	ds := svcat.NewAlbListenerCertificatesDataSource().(*svcat.AlbListenerCertificatesDataSource)
	setAlbListenerCertificatesClient(ds, client)

	sch := albListenerCertificatesSchema(t)
	raw := buildAlbListenerCertificatesRawValue("SVR", "lst-d9p3k7wa", "alb-f8q2xk9m")
	req := newAlbListenerCertificatesReadRequest(sch, raw)
	resp := newAlbListenerCertificatesReadResponse(sch)

	ds.Read(context.Background(), req, resp)

	assert.False(t, resp.Diagnostics.HasError(), "expected no error diagnostics, got: %v", resp.Diagnostics)

	var result svcat.AlbListenerCertificatesDataSourceModel
	resp.Diagnostics.Append(resp.State.Get(context.Background(), &result)...)
	assert.False(t, resp.Diagnostics.HasError())
	assert.False(t, result.Certificates.IsNull())

	var items []svcat.AlbListenerCertificateItemModel
	resp.Diagnostics.Append(result.Certificates.ElementsAs(context.Background(), &items, false)...)
	assert.False(t, resp.Diagnostics.HasError())
	assert.Len(t, items, 0)
}

func TestAlbListenerCertificates_ReadApiError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	albClient := &albv20251030.Client{}
	patches.ApplyMethodReturn(client, "UseAlbV20251030Client", albClient)
	patches.ApplyMethodFunc(albClient, "DescribeListenerCertificates", func(request *albv20251030.DescribeListenerCertificatesRequest) (*albv20251030.DescribeListenerCertificatesResponse, error) {
		return nil, fmt.Errorf("mocked API error")
	})

	ds := svcat.NewAlbListenerCertificatesDataSource().(*svcat.AlbListenerCertificatesDataSource)
	setAlbListenerCertificatesClient(ds, client)

	sch := albListenerCertificatesSchema(t)
	raw := buildAlbListenerCertificatesRawValue("SVR", "lst-d9p3k7wa", "alb-f8q2xk9m")
	req := newAlbListenerCertificatesReadRequest(sch, raw)
	resp := newAlbListenerCertificatesReadResponse(sch)

	ds.Read(context.Background(), req, resp)

	assert.True(t, resp.Diagnostics.HasError(), "expected error diagnostics for API failure")
}

func TestAlbListenerCertificatesService_ReadAllPages(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	albClient := &albv20251030.Client{}
	patches.ApplyMethodReturn(client, "UseAlbV20251030Client", albClient)

	callCount := 0
	patches.ApplyMethodFunc(albClient, "DescribeListenerCertificates", func(request *albv20251030.DescribeListenerCertificatesRequest) (*albv20251030.DescribeListenerCertificatesResponse, error) {
		callCount++
		if request.MaxResults == nil || *request.MaxResults != 100 {
			t.Errorf("expected MaxResults 100 on call %d, got %v", callCount, request.MaxResults)
		}
		switch callCount {
		case 1:
			if request.NextToken != nil {
				t.Errorf("expected no NextToken on first call, got %v", *request.NextToken)
			}
			return &albv20251030.DescribeListenerCertificatesResponse{Response: &albv20251030.DescribeListenerCertificatesResponseParams{
				Certificates: []*albv20251030.CertificateInfo{
					{CertificateId: albListenerCertificatesStr("cert-page1-0"), Status: albListenerCertificatesStr("Associated")},
				},
				NextToken: albListenerCertificatesStr("next-token-page2"),
			}}, nil
		case 2:
			if request.NextToken == nil || *request.NextToken != "next-token-page2" {
				t.Errorf("expected NextToken from first page on second call")
			}
			return &albv20251030.DescribeListenerCertificatesResponse{Response: &albv20251030.DescribeListenerCertificatesResponseParams{
				Certificates: []*albv20251030.CertificateInfo{
					{CertificateId: albListenerCertificatesStr("cert-page2-0"), Status: albListenerCertificatesStr("Associated")},
				},
			}}, nil
		default:
			t.Errorf("unexpected extra API call %d", callCount)
			return &albv20251030.DescribeListenerCertificatesResponse{}, nil
		}
	})

	service := svcat.NewAlbService(client)
	describeListenerCertificatesRequest := albv20251030.NewDescribeListenerCertificatesRequest()
	describeListenerCertificatesRequest.CertificateType = albListenerCertificatesStr("SVR")
	describeListenerCertificatesRequest.ListenerId = albListenerCertificatesStr("lst-d9p3k7wa")
	describeListenerCertificatesRequest.LoadBalancerId = albListenerCertificatesStr("alb-f8q2xk9m")

	certificates, err := service.DescribeListenerCertificates(context.Background(), describeListenerCertificatesRequest)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if callCount != 2 {
		t.Errorf("expected 2 API calls, got %d", callCount)
	}
	if len(certificates) != 2 {
		t.Fatalf("expected 2 certificates across pages, got %d", len(certificates))
	}
	if certificates[0].CertificateId == nil || *certificates[0].CertificateId != "cert-page1-0" {
		t.Errorf("unexpected first certificate id: %v", certificates[0].CertificateId)
	}
	if certificates[1].CertificateId == nil || *certificates[1].CertificateId != "cert-page2-0" {
		t.Errorf("unexpected last certificate id: %v", certificates[1].CertificateId)
	}
}
