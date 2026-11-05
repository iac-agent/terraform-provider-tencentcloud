package teo_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	datasource_schema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/assert"
	teosdk "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/sharedmeta"
	svcteo "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/teo"
)

// ptrStrTeoTokens returns a pointer to the given string.
func ptrStrTeoTokens(s string) *string {
	return &s
}

// ptrInt64TeoTokens returns a pointer to the given int64.
func ptrInt64TeoTokens(n int64) *int64 {
	return &n
}

// teoInferenceAPITokensSchema returns the framework data source schema used
// to build the raw config value for the ReadRequest.
func teoInferenceAPITokensSchema(t *testing.T) datasource_schema.Schema {
	ds := svcteo.NewTeoInferenceAPITokensDataSource()
	schemaResp := &datasource.SchemaResponse{}
	ds.Schema(context.Background(), datasource.SchemaRequest{}, schemaResp)
	if schemaResp.Diagnostics.HasError() {
		t.Fatalf("failed to build data source schema: %v", schemaResp.Diagnostics)
	}
	return schemaResp.Schema
}

// teoInferenceAPITokenObjType builds the tftypes.Object for a tokens element.
func teoInferenceAPITokenObjType() tftypes.Object {
	return tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"token_id":    tftypes.String,
			"name":        tftypes.String,
			"content":     tftypes.String,
			"create_time": tftypes.String,
		},
	}
}

// newTeoInferenceAPITokensReadRequest builds a datasource.ReadRequest with the
// given zone_id as the raw config value.
func newTeoInferenceAPITokensReadRequest(t *testing.T, zoneId string) datasource.ReadRequest {
	s := teoInferenceAPITokensSchema(t)

	raw := tftypes.NewValue(tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"zone_id": tftypes.String,
			"tokens":  tftypes.List{ElementType: teoInferenceAPITokenObjType()},
		},
	}, map[string]tftypes.Value{
		"zone_id": tftypes.NewValue(tftypes.String, zoneId),
		"tokens":  tftypes.NewValue(tftypes.List{ElementType: teoInferenceAPITokenObjType()}, nil),
	})

	req := datasource.ReadRequest{}
	req.Config.Raw = raw
	req.Config.Schema = s
	return req
}

// setTeoInferenceAPITokensClient injects a mock connectivity client into the
// data source instance via the Configure method.
func setTeoInferenceAPITokensClient(ds *svcteo.TeoInferenceAPITokensDataSource, client *connectivity.TencentCloudClient) {
	meta := &sharedmeta.ProviderMeta{Client: client}
	configureReq := datasource.ConfigureRequest{ProviderData: meta}
	configureResp := &datasource.ConfigureResponse{}
	ds.Configure(context.Background(), configureReq, configureResp)
}

// newTeoInferenceAPITokensReadResponse builds a datasource.ReadResponse with
// the state pre-populated from the config raw value (mirrors fwserver).
func newTeoInferenceAPITokensReadResponse(t *testing.T) *datasource.ReadResponse {
	s := teoInferenceAPITokensSchema(t)

	raw := tftypes.NewValue(tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"zone_id": tftypes.String,
			"tokens":  tftypes.List{ElementType: teoInferenceAPITokenObjType()},
		},
	}, map[string]tftypes.Value{
		"zone_id": tftypes.NewValue(tftypes.String, "zone-2qtuhspy7cr6"),
		"tokens":  tftypes.NewValue(tftypes.List{ElementType: teoInferenceAPITokenObjType()}, nil),
	})

	return &datasource.ReadResponse{
		State: tfsdk.State{
			Schema: s,
			Raw:    raw,
		},
	}
}

// go test ./tencentcloud/services/teo/ -run "TestTeoInferenceAPITokensDS" -v -count=1 -gcflags="all=-l"

// TestTeoInferenceAPITokensDS_ReadBasic covers the normal read path with two
// tokens returned by the API, asserting every leaf field of the collection.
func TestTeoInferenceAPITokensDS_ReadBasic(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	teoClient := &teosdk.Client{}
	patches.ApplyMethodReturn(client, "UseTeoV20220901Client", teoClient)

	var gotRequest *teosdk.DescribeInferenceAPITokensRequest
	patches.ApplyMethodFunc(teoClient, "DescribeInferenceAPITokensWithContext", func(_ context.Context, request *teosdk.DescribeInferenceAPITokensRequest) (*teosdk.DescribeInferenceAPITokensResponse, error) {
		gotRequest = request
		assert.Equal(t, "zone-2qtuhspy7cr6", *request.ZoneId)
		resp := teosdk.NewDescribeInferenceAPITokensResponse()
		resp.Response = &teosdk.DescribeInferenceAPITokensResponseParams{
			TotalCount: ptrInt64TeoTokens(2),
			Tokens: []*teosdk.InferenceAPIToken{
				{
					TokenId:    ptrStrTeoTokens("token-abc123"),
					Name:       ptrStrTeoTokens("demo-token"),
					Content:    ptrStrTeoTokens("tok-secret-content-1"),
					CreateTime: ptrStrTeoTokens("2025-01-01T00:00:00Z"),
				},
				{
					TokenId:    ptrStrTeoTokens("token-def456"),
					Name:       ptrStrTeoTokens("prod-token"),
					Content:    ptrStrTeoTokens("tok-secret-content-2"),
					CreateTime: ptrStrTeoTokens("2025-06-01T00:00:00Z"),
				},
			},
			RequestId: ptrStrTeoTokens("fake-request-id"),
		}
		return resp, nil
	})

	ds := svcteo.NewTeoInferenceAPITokensDataSource().(*svcteo.TeoInferenceAPITokensDataSource)
	setTeoInferenceAPITokensClient(ds, client)

	req := newTeoInferenceAPITokensReadRequest(t, "zone-2qtuhspy7cr6")
	resp := newTeoInferenceAPITokensReadResponse(t)
	ds.Read(context.Background(), req, resp)

	assert.False(t, resp.Diagnostics.HasError(), "expected no error diagnostics, got: %v", resp.Diagnostics)
	assert.NotNil(t, gotRequest)

	var state svcteo.TeoInferenceAPITokensDataSourceModel
	diags := resp.State.Get(context.Background(), &state)
	assert.False(t, diags.HasError(), "failed to read state: %v", diags)

	assert.False(t, state.Tokens.IsNull())
	assert.Equal(t, 2, len(state.Tokens.Elements()))

	var tokens []svcteo.TeoInferenceAPITokenModel
	if !state.Tokens.IsNull() {
		elemDiags := state.Tokens.ElementsAs(context.Background(), &tokens, false)
		assert.False(t, elemDiags.HasError(), "failed to convert tokens: %v", elemDiags)
	}

	assert.Len(t, tokens, 2)
	assert.Equal(t, "token-abc123", tokens[0].TokenId.ValueString())
	assert.Equal(t, "demo-token", tokens[0].Name.ValueString())
	assert.Equal(t, "tok-secret-content-1", tokens[0].Content.ValueString())
	assert.Equal(t, "2025-01-01T00:00:00Z", tokens[0].CreateTime.ValueString())

	assert.Equal(t, "token-def456", tokens[1].TokenId.ValueString())
	assert.Equal(t, "prod-token", tokens[1].Name.ValueString())
	assert.Equal(t, "tok-secret-content-2", tokens[1].Content.ValueString())
	assert.Equal(t, "2025-06-01T00:00:00Z", tokens[1].CreateTime.ValueString())
}

// TestTeoInferenceAPITokensDS_ReadEmptyResult covers a legitimately empty
// token list: an empty (but non-null) list must be written to state with no
// error.
func TestTeoInferenceAPITokensDS_ReadEmptyResult(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	teoClient := &teosdk.Client{}
	patches.ApplyMethodReturn(client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeInferenceAPITokensWithContext", func(_ context.Context, request *teosdk.DescribeInferenceAPITokensRequest) (*teosdk.DescribeInferenceAPITokensResponse, error) {
		resp := teosdk.NewDescribeInferenceAPITokensResponse()
		resp.Response = &teosdk.DescribeInferenceAPITokensResponseParams{
			TotalCount: ptrInt64TeoTokens(0),
			Tokens:     []*teosdk.InferenceAPIToken{},
			RequestId:  ptrStrTeoTokens("fake-request-id"),
		}
		return resp, nil
	})

	ds := svcteo.NewTeoInferenceAPITokensDataSource().(*svcteo.TeoInferenceAPITokensDataSource)
	setTeoInferenceAPITokensClient(ds, client)

	req := newTeoInferenceAPITokensReadRequest(t, "zone-2qtuhspy7cr6")
	resp := newTeoInferenceAPITokensReadResponse(t)
	ds.Read(context.Background(), req, resp)

	assert.False(t, resp.Diagnostics.HasError(), "expected no error diagnostics, got: %v", resp.Diagnostics)

	var state svcteo.TeoInferenceAPITokensDataSourceModel
	diags := resp.State.Get(context.Background(), &state)
	assert.False(t, diags.HasError(), "failed to read state: %v", diags)

	assert.False(t, state.Tokens.IsNull())
	assert.Equal(t, 0, len(state.Tokens.Elements()))
}

// TestTeoInferenceAPITokensDS_ReadNilLeafFields covers nil pointer fields in
// the API response: the corresponding state attributes must become Null.
func TestTeoInferenceAPITokensDS_ReadNilLeafFields(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	teoClient := &teosdk.Client{}
	patches.ApplyMethodReturn(client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeInferenceAPITokensWithContext", func(_ context.Context, request *teosdk.DescribeInferenceAPITokensRequest) (*teosdk.DescribeInferenceAPITokensResponse, error) {
		resp := teosdk.NewDescribeInferenceAPITokensResponse()
		resp.Response = &teosdk.DescribeInferenceAPITokensResponseParams{
			TotalCount: ptrInt64TeoTokens(1),
			Tokens: []*teosdk.InferenceAPIToken{
				{
					TokenId: ptrStrTeoTokens("token-abc123"),
				},
			},
			RequestId: ptrStrTeoTokens("fake-request-id"),
		}
		return resp, nil
	})

	ds := svcteo.NewTeoInferenceAPITokensDataSource().(*svcteo.TeoInferenceAPITokensDataSource)
	setTeoInferenceAPITokensClient(ds, client)

	req := newTeoInferenceAPITokensReadRequest(t, "zone-2qtuhspy7cr6")
	resp := newTeoInferenceAPITokensReadResponse(t)
	ds.Read(context.Background(), req, resp)

	assert.False(t, resp.Diagnostics.HasError(), "expected no error diagnostics, got: %v", resp.Diagnostics)

	var state svcteo.TeoInferenceAPITokensDataSourceModel
	diags := resp.State.Get(context.Background(), &state)
	assert.False(t, diags.HasError(), "failed to read state: %v", diags)

	var tokens []svcteo.TeoInferenceAPITokenModel
	elemDiags := state.Tokens.ElementsAs(context.Background(), &tokens, false)
	assert.False(t, elemDiags.HasError(), "failed to convert tokens: %v", elemDiags)

	assert.Len(t, tokens, 1)
	assert.Equal(t, "token-abc123", tokens[0].TokenId.ValueString())
	assert.True(t, tokens[0].Name.IsNull())
	assert.True(t, tokens[0].Content.IsNull())
	assert.True(t, tokens[0].CreateTime.IsNull())
}

// TestTeoInferenceAPITokensDS_ReadPagination covers offset pagination across
// two API pages: all tokens from both pages must be aggregated in state.
func TestTeoInferenceAPITokensDS_ReadPagination(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	teoClient := &teosdk.Client{}
	patches.ApplyMethodReturn(client, "UseTeoV20220901Client", teoClient)

	offsets := []int64{}
	patches.ApplyMethodFunc(teoClient, "DescribeInferenceAPITokensWithContext", func(_ context.Context, request *teosdk.DescribeInferenceAPITokensRequest) (*teosdk.DescribeInferenceAPITokensResponse, error) {
		offsets = append(offsets, *request.Offset)
		resp := teosdk.NewDescribeInferenceAPITokensResponse()
		if *request.Offset == 0 {
			resp.Response = &teosdk.DescribeInferenceAPITokensResponseParams{
				TotalCount: ptrInt64TeoTokens(3),
				Tokens: []*teosdk.InferenceAPIToken{
					{TokenId: ptrStrTeoTokens("token-1"), Name: ptrStrTeoTokens("one"), Content: ptrStrTeoTokens("content-1"), CreateTime: ptrStrTeoTokens("2025-01-01T00:00:00Z")},
					{TokenId: ptrStrTeoTokens("token-2"), Name: ptrStrTeoTokens("two"), Content: ptrStrTeoTokens("content-2"), CreateTime: ptrStrTeoTokens("2025-01-02T00:00:00Z")},
				},
			}
			return resp, nil
		}
		resp.Response = &teosdk.DescribeInferenceAPITokensResponseParams{
			TotalCount: ptrInt64TeoTokens(3),
			Tokens: []*teosdk.InferenceAPIToken{
				{TokenId: ptrStrTeoTokens("token-3"), Name: ptrStrTeoTokens("three"), Content: ptrStrTeoTokens("content-3"), CreateTime: ptrStrTeoTokens("2025-01-03T00:00:00Z")},
			},
		}
		return resp, nil
	})

	ds := svcteo.NewTeoInferenceAPITokensDataSource().(*svcteo.TeoInferenceAPITokensDataSource)
	setTeoInferenceAPITokensClient(ds, client)

	req := newTeoInferenceAPITokensReadRequest(t, "zone-2qtuhspy7cr6")
	resp := newTeoInferenceAPITokensReadResponse(t)
	ds.Read(context.Background(), req, resp)

	assert.False(t, resp.Diagnostics.HasError(), "expected no error diagnostics, got: %v", resp.Diagnostics)
	assert.Equal(t, []int64{0, 2}, offsets, "expected two pages with offsets 0 then 2")

	var state svcteo.TeoInferenceAPITokensDataSourceModel
	diags := resp.State.Get(context.Background(), &state)
	assert.False(t, diags.HasError(), "failed to read state: %v", diags)

	var tokens []svcteo.TeoInferenceAPITokenModel
	elemDiags := state.Tokens.ElementsAs(context.Background(), &tokens, false)
	assert.False(t, elemDiags.HasError(), "failed to convert tokens: %v", elemDiags)

	assert.Len(t, tokens, 3)
	assert.Equal(t, "token-1", tokens[0].TokenId.ValueString())
	assert.Equal(t, "token-2", tokens[1].TokenId.ValueString())
	assert.Equal(t, "token-3", tokens[2].TokenId.ValueString())
}

// TestTeoInferenceAPITokensDS_ReadAPIError covers the API error path: an
// error diagnostic must be returned.
func TestTeoInferenceAPITokensDS_ReadAPIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	teoClient := &teosdk.Client{}
	patches.ApplyMethodReturn(client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeInferenceAPITokensWithContext", func(_ context.Context, request *teosdk.DescribeInferenceAPITokensRequest) (*teosdk.DescribeInferenceAPITokensResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=AuthFailure, Message=authentication failed")
	})

	ds := svcteo.NewTeoInferenceAPITokensDataSource().(*svcteo.TeoInferenceAPITokensDataSource)
	setTeoInferenceAPITokensClient(ds, client)

	req := newTeoInferenceAPITokensReadRequest(t, "zone-2qtuhspy7cr6")
	resp := newTeoInferenceAPITokensReadResponse(t)
	ds.Read(context.Background(), req, resp)

	assert.True(t, resp.Diagnostics.HasError())
	assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "querying TEO inference API tokens")
	assert.Contains(t, resp.Diagnostics.Errors()[0].Detail(), "AuthFailure")
}

// TestTeoInferenceAPITokensDS_ReadNilResponse covers the nil response path:
// a nil response is an API failure, not an empty result.
func TestTeoInferenceAPITokensDS_ReadNilResponse(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	teoClient := &teosdk.Client{}
	patches.ApplyMethodReturn(client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeInferenceAPITokensWithContext", func(_ context.Context, request *teosdk.DescribeInferenceAPITokensRequest) (*teosdk.DescribeInferenceAPITokensResponse, error) {
		resp := teosdk.NewDescribeInferenceAPITokensResponse()
		resp.Response = nil
		return resp, nil
	})

	ds := svcteo.NewTeoInferenceAPITokensDataSource().(*svcteo.TeoInferenceAPITokensDataSource)
	setTeoInferenceAPITokensClient(ds, client)

	req := newTeoInferenceAPITokensReadRequest(t, "zone-2qtuhspy7cr6")
	resp := newTeoInferenceAPITokensReadResponse(t)
	ds.Read(context.Background(), req, resp)

	assert.True(t, resp.Diagnostics.HasError())
	assert.Contains(t, resp.Diagnostics.Errors()[0].Detail(), "response is nil")
}

// TestTeoInferenceAPITokensDS_ReadMissingZoneId covers the empty zone_id
// guard: an error must be raised before any API call.
func TestTeoInferenceAPITokensDS_ReadMissingZoneId(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	teoClient := &teosdk.Client{}
	patches.ApplyMethodReturn(client, "UseTeoV20220901Client", teoClient)

	apiCalled := false
	patches.ApplyMethodFunc(teoClient, "DescribeInferenceAPITokensWithContext", func(_ context.Context, request *teosdk.DescribeInferenceAPITokensRequest) (*teosdk.DescribeInferenceAPITokensResponse, error) {
		apiCalled = true
		return teosdk.NewDescribeInferenceAPITokensResponse(), nil
	})

	ds := svcteo.NewTeoInferenceAPITokensDataSource().(*svcteo.TeoInferenceAPITokensDataSource)
	setTeoInferenceAPITokensClient(ds, client)

	req := newTeoInferenceAPITokensReadRequest(t, "")
	resp := newTeoInferenceAPITokensReadResponse(t)
	ds.Read(context.Background(), req, resp)

	assert.True(t, resp.Diagnostics.HasError())
	assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Missing zone_id")
	assert.False(t, apiCalled)
}

// TestTeoInferenceAPITokensDS_ReadClientNotConfigured covers the nil-client
// guard when the provider is not configured.
func TestTeoInferenceAPITokensDS_ReadClientNotConfigured(t *testing.T) {
	ds := svcteo.NewTeoInferenceAPITokensDataSource().(*svcteo.TeoInferenceAPITokensDataSource)

	req := newTeoInferenceAPITokensReadRequest(t, "zone-2qtuhspy7cr6")
	resp := newTeoInferenceAPITokensReadResponse(t)
	ds.Read(context.Background(), req, resp)

	assert.True(t, resp.Diagnostics.HasError())
	assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Provider not configured")
}

// TestTeoInferenceAPITokensDS_MetadataAndSchema validates the metadata type
// name and schema flags of the data source.
func TestTeoInferenceAPITokensDS_MetadataAndSchema(t *testing.T) {
	ds := svcteo.NewTeoInferenceAPITokensDataSource().(*svcteo.TeoInferenceAPITokensDataSource)

	metaResp := &datasource.MetadataResponse{}
	ds.Metadata(context.Background(), datasource.MetadataRequest{}, metaResp)
	assert.Equal(t, "tencentcloud_teo_inference_api_tokens", metaResp.TypeName)

	schemaResp := &datasource.SchemaResponse{}
	ds.Schema(context.Background(), datasource.SchemaRequest{}, schemaResp)
	assert.False(t, schemaResp.Diagnostics.HasError())

	attrs := schemaResp.Schema.Attributes
	assert.Contains(t, attrs, "zone_id")
	assert.Contains(t, attrs, "tokens")

	zoneIdAttr, ok := attrs["zone_id"].(datasource_schema.StringAttribute)
	assert.True(t, ok)
	assert.True(t, zoneIdAttr.IsRequired())
	assert.False(t, zoneIdAttr.IsComputed())

	tokensAttr, ok := attrs["tokens"].(datasource_schema.ListNestedAttribute)
	assert.True(t, ok)
	assert.True(t, tokensAttr.IsComputed())
	assert.False(t, tokensAttr.IsSensitive())

	nestedAttrs := tokensAttr.NestedObject.Attributes
	assert.Contains(t, nestedAttrs, "token_id")
	assert.Contains(t, nestedAttrs, "name")
	assert.Contains(t, nestedAttrs, "content")
	assert.Contains(t, nestedAttrs, "create_time")

	tokenIdAttr, ok := nestedAttrs["token_id"].(datasource_schema.StringAttribute)
	assert.True(t, ok)
	assert.True(t, tokenIdAttr.IsComputed())
	assert.False(t, tokenIdAttr.IsSensitive())

	contentAttr, ok := nestedAttrs["content"].(datasource_schema.StringAttribute)
	assert.True(t, ok)
	assert.True(t, contentAttr.IsComputed())
	assert.True(t, contentAttr.IsSensitive())
}

// TestTeoInferenceAPITokensDS_FactoryAndInterfaces validates the factory
// return type and the implemented interfaces.
func TestTeoInferenceAPITokensDS_FactoryAndInterfaces(t *testing.T) {
	got := svcteo.NewTeoInferenceAPITokensDataSource()
	if got == nil {
		t.Fatal("NewTeoInferenceAPITokensDataSource returned nil")
	}
	if _, ok := got.(*svcteo.TeoInferenceAPITokensDataSource); !ok {
		t.Fatalf("factory returned %T, want *TeoInferenceAPITokensDataSource", got)
	}

	var _ datasource.DataSourceWithConfigure = &svcteo.TeoInferenceAPITokensDataSource{}
}

// TestTeoInferenceAPITokensDS_ConfigureNilProviderData ensures Configure with
// nil ProviderData is a silent no-op.
func TestTeoInferenceAPITokensDS_ConfigureNilProviderData(t *testing.T) {
	ds := svcteo.NewTeoInferenceAPITokensDataSource().(*svcteo.TeoInferenceAPITokensDataSource)

	req := datasource.ConfigureRequest{}
	resp := &datasource.ConfigureResponse{}
	ds.Configure(context.Background(), req, resp)

	assert.False(t, resp.Diagnostics.HasError())
}

// TestTeoInferenceAPITokensDS_ModelTypes validates the model struct fields
// are accessible (compile-time guard against accidental renames).
func TestTeoInferenceAPITokensDS_ModelTypes(t *testing.T) {
	model := svcteo.TeoInferenceAPITokensDataSourceModel{
		ZoneId: types.StringValue("zone-2qtuhspy7cr6"),
	}
	_ = model.ZoneId
	_ = model.Tokens

	item := svcteo.TeoInferenceAPITokenModel{}
	_ = item.TokenId
	_ = item.Name
	_ = item.Content
	_ = item.CreateTime
}
