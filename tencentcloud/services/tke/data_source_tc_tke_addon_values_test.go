package tke

import (
	"context"
	"fmt"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/assert"
	tke "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tke/v20180525"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/sharedmeta"
)

func strPtrAddonValues(s string) *string {
	return &s
}

func TestFillTkeAddonValues(t *testing.T) {
	cases := []struct {
		name                string
		in                  *tke.DescribeAddonValuesResponseParams
		expectValues        types.String
		expectDefaultValues types.String
	}{
		{
			name: "both fields set",
			in: &tke.DescribeAddonValuesResponseParams{
				Values:        strPtrAddonValues(`{"replicas":1}`),
				DefaultValues: strPtrAddonValues(`{"replicas":2}`),
			},
			expectValues:        types.StringValue(`{"replicas":1}`),
			expectDefaultValues: types.StringValue(`{"replicas":2}`),
		},
		{
			name: "values nil because addon not installed",
			in: &tke.DescribeAddonValuesResponseParams{
				Values:        nil,
				DefaultValues: strPtrAddonValues(`{"replicas":2}`),
			},
			expectValues:        types.StringNull(),
			expectDefaultValues: types.StringValue(`{"replicas":2}`),
		},
		{
			name: "default values nil",
			in: &tke.DescribeAddonValuesResponseParams{
				Values:        strPtrAddonValues(`{"replicas":1}`),
				DefaultValues: nil,
			},
			expectValues:        types.StringValue(`{"replicas":1}`),
			expectDefaultValues: types.StringNull(),
		},
		{
			name: "both fields nil",
			in: &tke.DescribeAddonValuesResponseParams{
				Values:        nil,
				DefaultValues: nil,
			},
			expectValues:        types.StringNull(),
			expectDefaultValues: types.StringNull(),
		},
		{
			name:                "nil response params",
			in:                  nil,
			expectValues:        types.StringNull(),
			expectDefaultValues: types.StringNull(),
		},
		{
			name: "empty strings are preserved",
			in: &tke.DescribeAddonValuesResponseParams{
				Values:        strPtrAddonValues(""),
				DefaultValues: strPtrAddonValues(""),
			},
			expectValues:        types.StringValue(""),
			expectDefaultValues: types.StringValue(""),
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var model TkeAddonValuesDataSourceModel
			fillTkeAddonValues(&model, c.in)

			assert.Equal(t, c.expectValues, model.Values)
			assert.Equal(t, c.expectDefaultValues, model.DefaultValues)
		})
	}
}

func TestTkeAddonValuesDataSourceMetadata(t *testing.T) {
	ds := NewTkeAddonValuesDataSource()
	assert.NotNil(t, ds)

	var resp datasource.MetadataResponse
	ds.Metadata(context.Background(), datasource.MetadataRequest{ProviderTypeName: "tencentcloud"}, &resp)
	assert.Equal(t, "tencentcloud_tke_addon_values", resp.TypeName)
}

func TestTkeAddonValuesDataSourceSchema(t *testing.T) {
	ds := NewTkeAddonValuesDataSource()

	var resp datasource.SchemaResponse
	ds.Schema(context.Background(), datasource.SchemaRequest{}, &resp)

	clusterId, ok := resp.Schema.Attributes["cluster_id"].(schema.StringAttribute)
	assert.True(t, ok)
	assert.True(t, clusterId.Required)
	assert.False(t, clusterId.Computed)

	addonName, ok := resp.Schema.Attributes["addon_name"].(schema.StringAttribute)
	assert.True(t, ok)
	assert.True(t, addonName.Required)
	assert.False(t, addonName.Computed)

	values, ok := resp.Schema.Attributes["values"].(schema.StringAttribute)
	assert.True(t, ok)
	assert.True(t, values.Computed)
	assert.False(t, values.Required)

	defaultValues, ok := resp.Schema.Attributes["default_values"].(schema.StringAttribute)
	assert.True(t, ok)
	assert.True(t, defaultValues.Computed)
	assert.False(t, defaultValues.Required)
}

func TestTkeAddonValuesDataSourceConfigure(t *testing.T) {
	ds := &TkeAddonValuesDataSource{}

	// nil ProviderData is a no-op, not an error.
	var nilResp datasource.ConfigureResponse
	ds.Configure(context.Background(), datasource.ConfigureRequest{ProviderData: nil}, &nilResp)
	assert.False(t, nilResp.Diagnostics.HasError())

	// wrong type produces a diagnostic.
	var wrongResp datasource.ConfigureResponse
	ds.Configure(context.Background(), datasource.ConfigureRequest{ProviderData: "unexpected"}, &wrongResp)
	assert.True(t, wrongResp.Diagnostics.HasError())

	// correct type injects the shared client.
	client := &connectivity.TencentCloudClient{}
	var okResp datasource.ConfigureResponse
	ds.Configure(context.Background(), datasource.ConfigureRequest{ProviderData: &sharedmeta.ProviderMeta{Client: client}}, &okResp)
	assert.False(t, okResp.Diagnostics.HasError())
	assert.NotNil(t, ds.client.client)
}

// addonValuesSchema returns the data source schema for building test requests.
func addonValuesSchema(t *testing.T) schema.Schema {
	t.Helper()

	var resp datasource.SchemaResponse
	NewTkeAddonValuesDataSource().Schema(context.Background(), datasource.SchemaRequest{}, &resp)
	return resp.Schema
}

// buildAddonValuesConfig builds a datasource.ReadRequest with the supplied inputs.
func buildAddonValuesConfig(t *testing.T, clusterId, addonName string) datasource.ReadRequest {
	t.Helper()

	s := addonValuesSchema(t)
	objectType := s.Type().TerraformType(context.Background())

	raw := tftypes.NewValue(objectType, map[string]tftypes.Value{
		"cluster_id":     tftypes.NewValue(tftypes.String, clusterId),
		"addon_name":     tftypes.NewValue(tftypes.String, addonName),
		"values":         tftypes.NewValue(tftypes.String, nil),
		"default_values": tftypes.NewValue(tftypes.String, nil),
	})

	return datasource.ReadRequest{
		Config: tfsdk.Config{
			Raw:    raw,
			Schema: s,
		},
	}
}

// emptyAddonValuesState builds a datasource.ReadResponse holding an empty state.
func emptyAddonValuesState(t *testing.T) datasource.ReadResponse {
	t.Helper()

	s := addonValuesSchema(t)
	objectType := s.Type().TerraformType(context.Background())

	return datasource.ReadResponse{
		State: tfsdk.State{
			Raw:    tftypes.NewValue(objectType, nil),
			Schema: s,
		},
	}
}

// newAddonValuesDataSourceWithMock builds a data source whose TKE client calls
// are delegated to the supplied function.
func newAddonValuesDataSourceWithMock(patches *gomonkey.Patches, fn func(context.Context, *tke.DescribeAddonValuesRequest) (*tke.DescribeAddonValuesResponse, error)) *TkeAddonValuesDataSource {
	client := &connectivity.TencentCloudClient{}
	tkeClient := &tke.Client{}

	patches.ApplyMethodReturn(client, "UseTkeClient", tkeClient)
	patches.ApplyMethodFunc(tkeClient, "DescribeAddonValuesWithContext", fn)

	return &TkeAddonValuesDataSource{client: NewTkeService(client)}
}

func TestTkeAddonValuesDataSourceReadSuccess(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ds := newAddonValuesDataSourceWithMock(patches, func(_ context.Context, request *tke.DescribeAddonValuesRequest) (*tke.DescribeAddonValuesResponse, error) {
		assert.NotNil(t, request.ClusterId)
		assert.Equal(t, "cls-5yezvaxo", *request.ClusterId)
		assert.NotNil(t, request.AddonName)
		assert.Equal(t, "cbs", *request.AddonName)

		resp := tke.NewDescribeAddonValuesResponse()
		resp.Response = &tke.DescribeAddonValuesResponseParams{
			Values:        strPtrAddonValues(`{"replicas":1}`),
			DefaultValues: strPtrAddonValues(`{"replicas":2}`),
			RequestId:     strPtrAddonValues("fake-request-id"),
		}
		return resp, nil
	})

	resp := emptyAddonValuesState(t)
	ds.Read(context.Background(), buildAddonValuesConfig(t, "cls-5yezvaxo", "cbs"), &resp)
	assert.False(t, resp.Diagnostics.HasError())

	var result TkeAddonValuesDataSourceModel
	resp.Diagnostics.Append(resp.State.Get(context.Background(), &result)...)
	assert.False(t, resp.Diagnostics.HasError())
	assert.Equal(t, "cls-5yezvaxo", result.ClusterId.ValueString())
	assert.Equal(t, "cbs", result.AddonName.ValueString())
	assert.Equal(t, `{"replicas":1}`, result.Values.ValueString())
	assert.Equal(t, `{"replicas":2}`, result.DefaultValues.ValueString())
}

func TestTkeAddonValuesDataSourceReadAddonNotInstalled(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ds := newAddonValuesDataSourceWithMock(patches, func(_ context.Context, _ *tke.DescribeAddonValuesRequest) (*tke.DescribeAddonValuesResponse, error) {
		resp := tke.NewDescribeAddonValuesResponse()
		resp.Response = &tke.DescribeAddonValuesResponseParams{
			Values:        nil,
			DefaultValues: strPtrAddonValues(`{"replicas":2}`),
			RequestId:     strPtrAddonValues("fake-request-id"),
		}
		return resp, nil
	})

	resp := emptyAddonValuesState(t)
	ds.Read(context.Background(), buildAddonValuesConfig(t, "cls-5yezvaxo", "cbs"), &resp)
	assert.False(t, resp.Diagnostics.HasError())

	var result TkeAddonValuesDataSourceModel
	resp.Diagnostics.Append(resp.State.Get(context.Background(), &result)...)
	assert.False(t, resp.Diagnostics.HasError())
	assert.True(t, result.Values.IsNull())
	assert.Equal(t, `{"replicas":2}`, result.DefaultValues.ValueString())
}

func TestTkeAddonValuesDataSourceReadEmptyClusterId(t *testing.T) {
	ds := &TkeAddonValuesDataSource{}

	resp := emptyAddonValuesState(t)
	ds.Read(context.Background(), buildAddonValuesConfig(t, "", "cbs"), &resp)
	assert.True(t, resp.Diagnostics.HasError())
	assert.Contains(t, resp.Diagnostics.Errors()[0].Detail(), "cluster_id")
}

func TestTkeAddonValuesDataSourceReadEmptyAddonName(t *testing.T) {
	ds := &TkeAddonValuesDataSource{}

	resp := emptyAddonValuesState(t)
	ds.Read(context.Background(), buildAddonValuesConfig(t, "cls-5yezvaxo", ""), &resp)
	assert.True(t, resp.Diagnostics.HasError())
	assert.Contains(t, resp.Diagnostics.Errors()[0].Detail(), "addon_name")
}

func TestTkeAddonValuesDataSourceReadAPIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ds := newAddonValuesDataSourceWithMock(patches, func(_ context.Context, _ *tke.DescribeAddonValuesRequest) (*tke.DescribeAddonValuesResponse, error) {
		return nil, fmt.Errorf("ResourceUnavailable: cluster not found")
	})

	resp := emptyAddonValuesState(t)
	ds.Read(context.Background(), buildAddonValuesConfig(t, "cls-5yezvaxo", "cbs"), &resp)
	assert.True(t, resp.Diagnostics.HasError())

	detail := resp.Diagnostics.Errors()[0].Detail()
	assert.Contains(t, detail, "DescribeAddonValues")
	assert.Contains(t, detail, "cls-5yezvaxo")
	assert.Contains(t, detail, "cbs")
}

func TestTkeAddonValuesDataSourceReadEmptyEnvelope(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ds := newAddonValuesDataSourceWithMock(patches, func(_ context.Context, _ *tke.DescribeAddonValuesRequest) (*tke.DescribeAddonValuesResponse, error) {
		resp := tke.NewDescribeAddonValuesResponse()
		resp.Response = nil
		return resp, nil
	})

	resp := emptyAddonValuesState(t)
	ds.Read(context.Background(), buildAddonValuesConfig(t, "cls-5yezvaxo", "cbs"), &resp)
	assert.True(t, resp.Diagnostics.HasError())

	detail := resp.Diagnostics.Errors()[0].Detail()
	assert.Contains(t, detail, "DescribeAddonValues")
	assert.Contains(t, detail, "empty response")
}
