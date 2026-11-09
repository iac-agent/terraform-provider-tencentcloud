package teo_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	datasource_schema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/assert"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/sharedmeta"
	svcteo "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/teo"
)

func ptrStrTeoSharedCnameV2(s string) *string {
	return &s
}

func ptrInt64TeoSharedCnameV2(i int64) *int64 {
	return &i
}

func ptrBoolTeoSharedCnameV2(b bool) *bool {
	return &b
}

// teoSharedCnameV2RawInfoType returns the raw tftypes type of the
// shared_cname_info attribute, matching the data source schema exactly.
func teoSharedCnameV2RawInfoType() tftypes.Type {
	ipsslObjType := tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"zone_id":           tftypes.String,
			"associated_domain": tftypes.String,
			"status":             tftypes.String,
		},
	}
	domainObjType := tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"zone_id":   tftypes.String,
			"type":      tftypes.String,
			"instance":  tftypes.String,
		},
	}
	infoObjType := tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"type":              tftypes.String,
			"shared_cname":      tftypes.String,
			"description":       tftypes.String,
			"ipssl_config":      ipsslObjType,
			"bind_domain_count": tftypes.Number,
			"acceleration_domains": tftypes.List{ElementType: domainObjType},
		},
	}
	return tftypes.List{ElementType: infoObjType}
}

// teoSharedCnameV2FilterObjType returns the raw tftypes type of a filters element.
func teoSharedCnameV2FilterObjType() tftypes.Type {
	return tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"name":   tftypes.String,
			"values": tftypes.Set{ElementType: tftypes.String},
			"fuzzy":  tftypes.Bool,
		},
	}
}

// teoSharedCnameV2TestSchema returns the framework data source schema used to
// build the raw config value for ReadRequest.
func teoSharedCnameV2TestSchema(t *testing.T) datasource_schema.Schema {
	ds := svcteo.NewTeoSharedCnameV2DataSource()
	schemaResp := &datasource.SchemaResponse{}
	ds.Schema(context.Background(), datasource.SchemaRequest{}, schemaResp)
	if schemaResp.Diagnostics.HasError() {
		t.Fatalf("failed to build data source schema: %v", schemaResp.Diagnostics)
	}
	return schemaResp.Schema
}

// newTeoSharedCnameV2ReadRequest builds a datasource.ReadRequest carrying the
// given query inputs as raw tftypes values.
func newTeoSharedCnameV2ReadRequest(t *testing.T, zoneId string, filters []tftypes.Value, direction, match, order *string) datasource.ReadRequest {
	s := teoSharedCnameV2TestSchema(t)

	filterObjType := teoSharedCnameV2FilterObjType()

	var filtersVal tftypes.Value
	if filters == nil {
		filtersVal = tftypes.NewValue(tftypes.List{ElementType: filterObjType}, nil)
	} else {
		filtersVal = tftypes.NewValue(tftypes.List{ElementType: filterObjType}, filters)
	}

	attrValues := map[string]tftypes.Value{
		"zone_id":           tftypes.NewValue(tftypes.String, zoneId),
		"filters":           filtersVal,
		"direction":         tftypes.NewValue(tftypes.String, nil),
		"match":             tftypes.NewValue(tftypes.String, nil),
		"order":             tftypes.NewValue(tftypes.String, nil),
		"shared_cname_info": tftypes.NewValue(teoSharedCnameV2RawInfoType(), nil),
	}
	if direction != nil {
		attrValues["direction"] = tftypes.NewValue(tftypes.String, *direction)
	}
	if match != nil {
		attrValues["match"] = tftypes.NewValue(tftypes.String, *match)
	}
	if order != nil {
		attrValues["order"] = tftypes.NewValue(tftypes.String, *order)
	}

	raw := tftypes.NewValue(tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"zone_id":           tftypes.String,
			"filters":           tftypes.List{ElementType: filterObjType},
			"direction":         tftypes.String,
			"match":             tftypes.String,
			"order":             tftypes.String,
			"shared_cname_info": teoSharedCnameV2RawInfoType(),
		},
	}, attrValues)

	req := datasource.ReadRequest{}
	req.Config.Raw = raw
	req.Config.Schema = s
	return req
}

// newTeoSharedCnameV2ReadResponse builds a datasource.ReadResponse whose
// State carries the data source schema, mirroring what the framework runtime
// provides during a real ReadDataSource RPC.
func newTeoSharedCnameV2ReadResponse(t *testing.T) *datasource.ReadResponse {
	resp := &datasource.ReadResponse{}
	resp.State.Schema = teoSharedCnameV2TestSchema(t)
	return resp
}

// setTeoSharedCnameV2Client injects a mock connectivity client into the data
// source instance via the Configure method.
func setTeoSharedCnameV2Client(d *svcteo.TeoSharedCnameV2DataSource, client *connectivity.TencentCloudClient) {
	meta := &sharedmeta.ProviderMeta{Client: client}
	configureReq := datasource.ConfigureRequest{ProviderData: meta}
	configureResp := &datasource.ConfigureResponse{}
	d.Configure(context.Background(), configureReq, configureResp)
}

// go test ./tencentcloud/services/teo/ -run "TestTeoSharedCnameV2DS" -v -count=1 -gcflags="all=-l"

// TestTeoSharedCnameV2DS_ReadBasic tests a successful read of the shared CNAME list.
func TestTeoSharedCnameV2DS_ReadBasic(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(client, "UseTeoClient", teoClient)

	var gotRequest *teov20220901.DescribeSharedCNAMERequest
	patches.ApplyMethodFunc(teoClient, "DescribeSharedCNAMEWithContext", func(_ context.Context, request *teov20220901.DescribeSharedCNAMERequest) (*teov20220901.DescribeSharedCNAMEResponse, error) {
		gotRequest = request
		resp := teov20220901.NewDescribeSharedCNAMEResponse()
		resp.Response = &teov20220901.DescribeSharedCNAMEResponseParams{
			TotalCount: ptrInt64TeoSharedCnameV2(2),
			SharedCNAMEInfo: []*teov20220901.SharedCNAMEInfo{
				{
					Type:            ptrStrTeoSharedCnameV2("custom"),
					SharedCNAME:     ptrStrTeoSharedCnameV2("prefix-1.2qtuhspy7cr6.share.dnse1.com"),
					Description:     ptrStrTeoSharedCnameV2("first shared cname"),
					BindDomainCount: ptrInt64TeoSharedCnameV2(2),
					AccelerationDomains: []*teov20220901.ReferenceHolder{
						{
							ZoneId:   ptrStrTeoSharedCnameV2("zone-2qtuhspy7cr6"),
							Type:     ptrStrTeoSharedCnameV2("acceleration-domain"),
							Instance: ptrStrTeoSharedCnameV2("www.example.com"),
						},
						{
							ZoneId:   ptrStrTeoSharedCnameV2("zone-2qtuhspy7cr6"),
							Type:     ptrStrTeoSharedCnameV2("acceleration-domain"),
							Instance: ptrStrTeoSharedCnameV2("api.example.com"),
						},
					},
				},
				{
					Type:            ptrStrTeoSharedCnameV2("ip-ssl"),
					SharedCNAME:     ptrStrTeoSharedCnameV2("prefix-2.2qtuhspy7cr6.share.dnse1.com"),
					Description:     ptrStrTeoSharedCnameV2("ip-ssl shared cname"),
					BindDomainCount: ptrInt64TeoSharedCnameV2(0),
					IPSSLConfig: &teov20220901.IPSSLConfig{
						ZoneId:           ptrStrTeoSharedCnameV2("zone-2qtuhspy7cr6"),
						AssociatedDomain: ptrStrTeoSharedCnameV2("ssl.example.com"),
						Status:           ptrStrTeoSharedCnameV2("bound"),
					},
				},
			},
			RequestId: ptrStrTeoSharedCnameV2("fake-request-id"),
		}
		return resp, nil
	})

	ds := svcteo.NewTeoSharedCnameV2DataSource()
	setTeoSharedCnameV2Client(ds.(*svcteo.TeoSharedCnameV2DataSource), client)

	req := newTeoSharedCnameV2ReadRequest(t, "zone-2qtuhspy7cr6", nil, nil, nil, nil)
	resp := newTeoSharedCnameV2ReadResponse(t)
	ds.Read(context.Background(), req, resp)

	assert.False(t, resp.Diagnostics.HasError(), "expected no error diagnostics, got: %v", resp.Diagnostics)

	assert.NotNil(t, gotRequest)
	assert.Equal(t, "zone-2qtuhspy7cr6", *gotRequest.ZoneId)
	assert.Empty(t, gotRequest.Filters)
	assert.Nil(t, gotRequest.Direction)
	assert.Nil(t, gotRequest.Match)
	assert.Nil(t, gotRequest.Order)
	assert.NotNil(t, gotRequest.Limit)
	assert.Equal(t, int64(200), *gotRequest.Limit)
	assert.NotNil(t, gotRequest.Offset)
	assert.Equal(t, int64(0), *gotRequest.Offset)

	var state svcteo.TeoSharedCnameV2DataSourceModel
	assert.False(t, resp.State.Raw.IsNull())
	diags := resp.State.Get(context.Background(), &state)
	assert.False(t, diags.HasError(), "expected no state read diagnostics, got: %v", diags)

	assert.Equal(t, "zone-2qtuhspy7cr6", state.ZoneId.ValueString())
	assert.True(t, state.Filters.IsNull())
	assert.True(t, state.Direction.IsNull())
	assert.True(t, state.Match.IsNull())
	assert.True(t, state.Order.IsNull())

	assert.False(t, state.SharedCnameInfo.IsNull())
	assert.Equal(t, 2, len(state.SharedCnameInfo.Elements()))

	// First element: custom shared CNAME without IP SSL config.
	elem0, ok := state.SharedCnameInfo.Elements()[0].(types.Object)
	assert.True(t, ok, "first element is not a types.Object")
	assert.Equal(t, "custom", elem0.Attributes()["type"].(types.String).ValueString())
	assert.Equal(t, "prefix-1.2qtuhspy7cr6.share.dnse1.com", elem0.Attributes()["shared_cname"].(types.String).ValueString())
	assert.Equal(t, "first shared cname", elem0.Attributes()["description"].(types.String).ValueString())
	assert.Equal(t, int64(2), elem0.Attributes()["bind_domain_count"].(types.Int64).ValueInt64())
	ipssl0 := elem0.Attributes()["ipssl_config"].(types.Object)
	assert.True(t, ipssl0.IsNull(), "expected null ipssl_config for custom type")
	domains0 := elem0.Attributes()["acceleration_domains"].(types.List)
	assert.Equal(t, 2, len(domains0.Elements()))
	domain00 := domains0.Elements()[0].(types.Object)
	assert.Equal(t, "zone-2qtuhspy7cr6", domain00.Attributes()["zone_id"].(types.String).ValueString())
	assert.Equal(t, "acceleration-domain", domain00.Attributes()["type"].(types.String).ValueString())
	assert.Equal(t, "www.example.com", domain00.Attributes()["instance"].(types.String).ValueString())
	domain01 := domains0.Elements()[1].(types.Object)
	assert.Equal(t, "api.example.com", domain01.Attributes()["instance"].(types.String).ValueString())

	// Second element: ip-ssl shared CNAME with IP SSL config and empty domains.
	elem1, ok := state.SharedCnameInfo.Elements()[1].(types.Object)
	assert.True(t, ok, "second element is not a types.Object")
	assert.Equal(t, "ip-ssl", elem1.Attributes()["type"].(types.String).ValueString())
	assert.Equal(t, "prefix-2.2qtuhspy7cr6.share.dnse1.com", elem1.Attributes()["shared_cname"].(types.String).ValueString())
	assert.Equal(t, "ip-ssl shared cname", elem1.Attributes()["description"].(types.String).ValueString())
	assert.Equal(t, int64(0), elem1.Attributes()["bind_domain_count"].(types.Int64).ValueInt64())
	ipssl1 := elem1.Attributes()["ipssl_config"].(types.Object)
	assert.False(t, ipssl1.IsNull())
	assert.Equal(t, "zone-2qtuhspy7cr6", ipssl1.Attributes()["zone_id"].(types.String).ValueString())
	assert.Equal(t, "ssl.example.com", ipssl1.Attributes()["associated_domain"].(types.String).ValueString())
	assert.Equal(t, "bound", ipssl1.Attributes()["status"].(types.String).ValueString())
	domains1 := elem1.Attributes()["acceleration_domains"].(types.List)
	assert.Equal(t, 0, len(domains1.Elements()))
}

// TestTeoSharedCnameV2DS_ReadWithFiltersAndSortParams tests that filters and
// sort parameters are passed through to the API request.
func TestTeoSharedCnameV2DS_ReadWithFiltersAndSortParams(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(client, "UseTeoClient", teoClient)

	var gotRequest *teov20220901.DescribeSharedCNAMERequest
	patches.ApplyMethodFunc(teoClient, "DescribeSharedCNAMEWithContext", func(_ context.Context, request *teov20220901.DescribeSharedCNAMERequest) (*teov20220901.DescribeSharedCNAMEResponse, error) {
		gotRequest = request
		resp := teov20220901.NewDescribeSharedCNAMEResponse()
		resp.Response = &teov20220901.DescribeSharedCNAMEResponseParams{
			TotalCount: ptrInt64TeoSharedCnameV2(0),
			SharedCNAMEInfo: []*teov20220901.SharedCNAMEInfo{},
			RequestId:        ptrStrTeoSharedCnameV2("fake-request-id"),
		}
		return resp, nil
	})

	ds := svcteo.NewTeoSharedCnameV2DataSource()
	setTeoSharedCnameV2Client(ds.(*svcteo.TeoSharedCnameV2DataSource), client)

	filterObjType := teoSharedCnameV2FilterObjType()
	filterVals := []tftypes.Value{
		tftypes.NewValue(filterObjType, map[string]tftypes.Value{
			"name":   tftypes.NewValue(tftypes.String, "shared-cname"),
			"values": tftypes.NewValue(tftypes.Set{ElementType: tftypes.String}, []tftypes.Value{tftypes.NewValue(tftypes.String, "prefix-1.2qtuhspy7cr6.share.dnse1.com")}),
			"fuzzy":  tftypes.NewValue(tftypes.Bool, nil),
		}),
		tftypes.NewValue(filterObjType, map[string]tftypes.Value{
			"name":   tftypes.NewValue(tftypes.String, "description"),
			"values": tftypes.NewValue(tftypes.Set{ElementType: tftypes.String}, []tftypes.Value{tftypes.NewValue(tftypes.String, "desc")}),
			"fuzzy":  tftypes.NewValue(tftypes.Bool, true),
		}),
	}

	req := newTeoSharedCnameV2ReadRequest(t, "zone-2qtuhspy7cr6", filterVals,
		ptrStrTeoSharedCnameV2("desc"), ptrStrTeoSharedCnameV2("any"), ptrStrTeoSharedCnameV2("create-time"))
	resp := newTeoSharedCnameV2ReadResponse(t)
	ds.Read(context.Background(), req, resp)

	assert.False(t, resp.Diagnostics.HasError(), "expected no error diagnostics, got: %v", resp.Diagnostics)

	assert.NotNil(t, gotRequest)
	assert.Len(t, gotRequest.Filters, 2)
	assert.Equal(t, "shared-cname", *gotRequest.Filters[0].Name)
	assert.Equal(t, []*string{ptrStrTeoSharedCnameV2("prefix-1.2qtuhspy7cr6.share.dnse1.com")}, gotRequest.Filters[0].Values)
	assert.Nil(t, gotRequest.Filters[0].Fuzzy)
	assert.Equal(t, "description", *gotRequest.Filters[1].Name)
	assert.Equal(t, []*string{ptrStrTeoSharedCnameV2("desc")}, gotRequest.Filters[1].Values)
	assert.NotNil(t, gotRequest.Filters[1].Fuzzy)
	assert.True(t, *gotRequest.Filters[1].Fuzzy)
	assert.Equal(t, "desc", *gotRequest.Direction)
	assert.Equal(t, "any", *gotRequest.Match)
	assert.Equal(t, "create-time", *gotRequest.Order)

	// Empty result list is a valid state, written to state as an empty list.
	var state svcteo.TeoSharedCnameV2DataSourceModel
	diags := resp.State.Get(context.Background(), &state)
	assert.False(t, diags.HasError(), "expected no state read diagnostics, got: %v", diags)
	assert.False(t, state.SharedCnameInfo.IsNull())
	assert.Equal(t, 0, len(state.SharedCnameInfo.Elements()))
}

// TestTeoSharedCnameV2DS_ReadPaginated tests that the read follows pagination
// until the full list is collected.
func TestTeoSharedCnameV2DS_ReadPaginated(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(client, "UseTeoClient", teoClient)

	var callCount int
	var offsets []int64
	patches.ApplyMethodFunc(teoClient, "DescribeSharedCNAMEWithContext", func(_ context.Context, request *teov20220901.DescribeSharedCNAMERequest) (*teov20220901.DescribeSharedCNAMEResponse, error) {
		callCount++
		if request.Offset != nil {
			offsets = append(offsets, *request.Offset)
		}
		resp := teov20220901.NewDescribeSharedCNAMEResponse()
		// First call returns a full page of 200 entries; the second call
		// returns one extra entry; the loop then stops.
		var infos []*teov20220901.SharedCNAMEInfo
		if callCount == 1 {
			for i := 0; i < 200; i++ {
				infos = append(infos, &teov20220901.SharedCNAMEInfo{
					Type:            ptrStrTeoSharedCnameV2("custom"),
					SharedCNAME:     ptrStrTeoSharedCnameV2(fmt.Sprintf("page1-entry-%d.share.dnse1.com", i)),
					BindDomainCount: ptrInt64TeoSharedCnameV2(0),
				})
			}
		} else {
			infos = append(infos, &teov20220901.SharedCNAMEInfo{
				Type:            ptrStrTeoSharedCnameV2("custom"),
				SharedCNAME:     ptrStrTeoSharedCnameV2("page2-entry-0.share.dnse1.com"),
				BindDomainCount: ptrInt64TeoSharedCnameV2(0),
			})
		}
		resp.Response = &teov20220901.DescribeSharedCNAMEResponseParams{
			TotalCount:      ptrInt64TeoSharedCnameV2(201),
			SharedCNAMEInfo: infos,
			RequestId:       ptrStrTeoSharedCnameV2("fake-request-id"),
		}
		return resp, nil
	})

	ds := svcteo.NewTeoSharedCnameV2DataSource()
	setTeoSharedCnameV2Client(ds.(*svcteo.TeoSharedCnameV2DataSource), client)

	req := newTeoSharedCnameV2ReadRequest(t, "zone-2qtuhspy7cr6", nil, nil, nil, nil)
	resp := newTeoSharedCnameV2ReadResponse(t)
	ds.Read(context.Background(), req, resp)

	assert.False(t, resp.Diagnostics.HasError(), "expected no error diagnostics, got: %v", resp.Diagnostics)
	assert.Equal(t, 2, callCount)
	assert.Equal(t, []int64{200}, offsets)

	var state svcteo.TeoSharedCnameV2DataSourceModel
	diags := resp.State.Get(context.Background(), &state)
	assert.False(t, diags.HasError(), "expected no state read diagnostics, got: %v", diags)
	assert.False(t, state.SharedCnameInfo.IsNull())
	assert.Equal(t, 201, len(state.SharedCnameInfo.Elements()))
}

// TestTeoSharedCnameV2DS_ReadEmptyResponse tests the nil response guard.
func TestTeoSharedCnameV2DS_ReadEmptyResponse(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(client, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeSharedCNAMEWithContext", func(_ context.Context, request *teov20220901.DescribeSharedCNAMERequest) (*teov20220901.DescribeSharedCNAMEResponse, error) {
		return teov20220901.NewDescribeSharedCNAMEResponse(), nil
	})

	ds := svcteo.NewTeoSharedCnameV2DataSource()
	setTeoSharedCnameV2Client(ds.(*svcteo.TeoSharedCnameV2DataSource), client)

	req := newTeoSharedCnameV2ReadRequest(t, "zone-2qtuhspy7cr6", nil, nil, nil, nil)
	resp := newTeoSharedCnameV2ReadResponse(t)
	ds.Read(context.Background(), req, resp)

	assert.True(t, resp.Diagnostics.HasError())
	assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Error querying TEO shared CNAMEs")
	assert.Contains(t, resp.Diagnostics.Errors()[0].Detail(), "empty response")
}

// TestTeoSharedCnameV2DS_ReadAPIError tests API error handling.
func TestTeoSharedCnameV2DS_ReadAPIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(client, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeSharedCNAMEWithContext", func(_ context.Context, request *teov20220901.DescribeSharedCNAMERequest) (*teov20220901.DescribeSharedCNAMEResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=AuthFailure, Message=authentication failed")
	})

	ds := svcteo.NewTeoSharedCnameV2DataSource()
	setTeoSharedCnameV2Client(ds.(*svcteo.TeoSharedCnameV2DataSource), client)

	req := newTeoSharedCnameV2ReadRequest(t, "zone-2qtuhspy7cr6", nil, nil, nil, nil)
	resp := newTeoSharedCnameV2ReadResponse(t)
	ds.Read(context.Background(), req, resp)

	assert.True(t, resp.Diagnostics.HasError())
	assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Error querying TEO shared CNAMEs")
	assert.Contains(t, resp.Diagnostics.Errors()[0].Detail(), "AuthFailure")
}

// TestTeoSharedCnameV2DS_ReadClientNotConfigured tests the nil-client guard.
func TestTeoSharedCnameV2DS_ReadClientNotConfigured(t *testing.T) {
	ds := svcteo.NewTeoSharedCnameV2DataSource()

	req := newTeoSharedCnameV2ReadRequest(t, "zone-2qtuhspy7cr6", nil, nil, nil, nil)
	resp := newTeoSharedCnameV2ReadResponse(t)
	ds.Read(context.Background(), req, resp)

	assert.True(t, resp.Diagnostics.HasError())
	assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Provider not configured")
}

// TestTeoSharedCnameV2DS_MetadataAndSchema validates metadata and schema definition.
func TestTeoSharedCnameV2DS_MetadataAndSchema(t *testing.T) {
	ds := svcteo.NewTeoSharedCnameV2DataSource()

	metaResp := &datasource.MetadataResponse{}
	ds.Metadata(context.Background(), datasource.MetadataRequest{}, metaResp)
	assert.Equal(t, "tencentcloud_teo_shared_cname_v2", metaResp.TypeName)

	schemaResp := &datasource.SchemaResponse{}
	ds.Schema(context.Background(), datasource.SchemaRequest{}, schemaResp)
	assert.False(t, schemaResp.Diagnostics.HasError())

	attrs := schemaResp.Schema.Attributes

	assert.True(t, attrs["zone_id"].IsRequired())
	assert.False(t, attrs["zone_id"].IsComputed())

	assert.True(t, attrs["filters"].IsOptional())
	assert.False(t, attrs["filters"].IsComputed())

	for _, optional := range []string{"direction", "match", "order"} {
		assert.True(t, attrs[optional].IsOptional())
		assert.False(t, attrs[optional].IsComputed())
	}

	assert.True(t, attrs["shared_cname_info"].IsComputed())
	assert.False(t, attrs["shared_cname_info"].IsOptional())
}
