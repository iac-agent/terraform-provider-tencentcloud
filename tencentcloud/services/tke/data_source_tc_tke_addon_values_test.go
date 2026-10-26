package tke_test

import (
	"context"
	"errors"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/assert"
	tkev20180525 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tke/v20180525"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/sharedmeta"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tke"
)

func ptrStringAddonValues(s string) *string {
	return &s
}

// buildAddonValuesDataSource returns a configured data source together with its schema,
// mirroring what the framework does before invoking Read.
func buildAddonValuesDataSource(t *testing.T) (datasource.DataSource, datasource.SchemaResponse) {
	t.Helper()

	ds := tke.NewTkeAddonValuesDataSource()

	var schemaResp datasource.SchemaResponse
	ds.Schema(context.Background(), datasource.SchemaRequest{}, &schemaResp)
	assert.False(t, schemaResp.Diagnostics.HasError())

	return ds, schemaResp
}

// configureAddonValuesDataSource injects a provider client into the data source.
func configureAddonValuesDataSource(t *testing.T, ds datasource.DataSource) {
	t.Helper()

	var configureResp datasource.ConfigureResponse
	ds.Configure(context.Background(), datasource.ConfigureRequest{
		ProviderData: &sharedmeta.ProviderMeta{Client: &connectivity.TencentCloudClient{}},
	}, &configureResp)
	assert.False(t, configureResp.Diagnostics.HasError())
}

// buildConfigValue builds the data source configuration value: required inputs set,
// computed outputs unknown.
func buildConfigValue(clusterId, addonName string) tftypes.Value {
	return tftypes.NewValue(
		tftypes.Object{
			AttributeTypes: map[string]tftypes.Type{
				"cluster_id":     tftypes.String,
				"addon_name":     tftypes.String,
				"values":         tftypes.String,
				"default_values": tftypes.String,
			},
		},
		map[string]tftypes.Value{
			"cluster_id":     tftypes.NewValue(tftypes.String, clusterId),
			"addon_name":     tftypes.NewValue(tftypes.String, addonName),
			"values":         tftypes.NewValue(tftypes.String, tftypes.UnknownValue),
			"default_values": tftypes.NewValue(tftypes.String, tftypes.UnknownValue),
		},
	)
}

// readAddonValues invokes Read and returns the resulting model.
func readAddonValues(t *testing.T, ds datasource.DataSource, schemaResp datasource.SchemaResponse, clusterId, addonName string) (datasource.ReadResponse, tke.TkeAddonValuesDataSourceModel) {
	t.Helper()

	readReq := datasource.ReadRequest{Config: tfsdk.Config{Schema: schemaResp.Schema}}
	readReq.Config.Raw = buildConfigValue(clusterId, addonName)

	readResp := datasource.ReadResponse{}
	readResp.State = tfsdk.State{Schema: schemaResp.Schema}

	ds.Read(context.Background(), readReq, &readResp)

	var model tke.TkeAddonValuesDataSourceModel
	if !readResp.Diagnostics.HasError() {
		readResp.Diagnostics.Append(readResp.State.Get(context.Background(), &model)...)
	}

	return readResp, model
}

// TestAddonValues_Metadata verifies the registered Terraform type name.
func TestAddonValues_Metadata(t *testing.T) {
	ds := tke.NewTkeAddonValuesDataSource()

	var resp datasource.MetadataResponse
	ds.Metadata(context.Background(), datasource.MetadataRequest{ProviderTypeName: "tencentcloud"}, &resp)

	assert.Equal(t, "tencentcloud_tke_addon_values", resp.TypeName)
}

// TestAddonValues_Schema verifies the required inputs and computed outputs.
func TestAddonValues_Schema(t *testing.T) {
	_, schemaResp := buildAddonValuesDataSource(t)

	assert.Contains(t, schemaResp.Schema.Description, "Use this data source to query TKE addon parameter values")

	assert.True(t, schemaResp.Schema.Attributes["cluster_id"].IsRequired())
	assert.True(t, schemaResp.Schema.Attributes["addon_name"].IsRequired())
	assert.False(t, schemaResp.Schema.Attributes["cluster_id"].IsComputed())
	assert.False(t, schemaResp.Schema.Attributes["addon_name"].IsComputed())

	assert.True(t, schemaResp.Schema.Attributes["values"].IsComputed())
	assert.True(t, schemaResp.Schema.Attributes["default_values"].IsComputed())
	assert.False(t, schemaResp.Schema.Attributes["values"].IsRequired())
	assert.False(t, schemaResp.Schema.Attributes["default_values"].IsRequired())
}

// TestAddonValues_Configure verifies that the provider client is injected.
func TestAddonValues_Configure(t *testing.T) {
	ds := tke.NewTkeAddonValuesDataSource()

	var configureResp datasource.ConfigureResponse
	ds.Configure(context.Background(), datasource.ConfigureRequest{ProviderData: nil}, &configureResp)
	assert.False(t, configureResp.Diagnostics.HasError())

	configureResp = datasource.ConfigureResponse{}
	ds.Configure(context.Background(), datasource.ConfigureRequest{
		ProviderData: &sharedmeta.ProviderMeta{Client: &connectivity.TencentCloudClient{}},
	}, &configureResp)
	assert.False(t, configureResp.Diagnostics.HasError())

	configureResp = datasource.ConfigureResponse{}
	ds.Configure(context.Background(), datasource.ConfigureRequest{ProviderData: "unexpected"}, &configureResp)
	assert.True(t, configureResp.Diagnostics.HasError())
}

// TestAddonValues_Read_Installed covers an installed addon: both outputs are populated.
func TestAddonValues_Read_Installed(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	sdkClient := &tkev20180525.Client{}
	patches.ApplyMethodReturn(&connectivity.TencentCloudClient{}, "UseTkeClient", sdkClient)
	patches.ApplyMethodFunc(sdkClient, "DescribeAddonValuesWithContext", func(ctx context.Context, request *tkev20180525.DescribeAddonValuesRequest) (*tkev20180525.DescribeAddonValuesResponse, error) {
		assert.Equal(t, "cls-5yezvaxo", *request.ClusterId)
		assert.Equal(t, "ip-masq-agent", *request.AddonName)

		resp := tkev20180525.NewDescribeAddonValuesResponse()
		resp.Response = &tkev20180525.DescribeAddonValuesResponseParams{
			Values:        ptrStringAddonValues(`{"nodeSelector":{"agent":"true"}}`),
			DefaultValues: ptrStringAddonValues(`{"nodeSelector":{}}`),
			RequestId:     ptrStringAddonValues("fake-request-id"),
		}
		return resp, nil
	})

	ds, schemaResp := buildAddonValuesDataSource(t)
	configureAddonValuesDataSource(t, ds)

	readResp, model := readAddonValues(t, ds, schemaResp, "cls-5yezvaxo", "ip-masq-agent")
	assert.False(t, readResp.Diagnostics.HasError())
	assert.Equal(t, "cls-5yezvaxo", model.ClusterId.ValueString())
	assert.Equal(t, "ip-masq-agent", model.AddonName.ValueString())
	assert.Equal(t, `{"nodeSelector":{"agent":"true"}}`, model.Values.ValueString())
	assert.Equal(t, `{"nodeSelector":{}}`, model.DefaultValues.ValueString())
}

// TestAddonValues_Read_NotInstalled covers an addon that is not installed: values is null
// while default_values is still returned, and no diagnostic error is raised.
func TestAddonValues_Read_NotInstalled(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	sdkClient := &tkev20180525.Client{}
	patches.ApplyMethodReturn(&connectivity.TencentCloudClient{}, "UseTkeClient", sdkClient)
	patches.ApplyMethodFunc(sdkClient, "DescribeAddonValuesWithContext", func(ctx context.Context, request *tkev20180525.DescribeAddonValuesRequest) (*tkev20180525.DescribeAddonValuesResponse, error) {
		resp := tkev20180525.NewDescribeAddonValuesResponse()
		resp.Response = &tkev20180525.DescribeAddonValuesResponseParams{
			DefaultValues: ptrStringAddonValues(`{"nodeSelector":{}}`),
			RequestId:     ptrStringAddonValues("fake-request-id"),
		}
		return resp, nil
	})

	ds, schemaResp := buildAddonValuesDataSource(t)
	configureAddonValuesDataSource(t, ds)

	readResp, model := readAddonValues(t, ds, schemaResp, "cls-5yezvaxo", "not-installed-addon")
	assert.False(t, readResp.Diagnostics.HasError())
	assert.True(t, model.Values.IsNull())
	assert.Equal(t, `{"nodeSelector":{}}`, model.DefaultValues.ValueString())
}

// TestAddonValues_Read_ApiError covers a failing API call.
func TestAddonValues_Read_ApiError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	sdkClient := &tkev20180525.Client{}
	patches.ApplyMethodReturn(&connectivity.TencentCloudClient{}, "UseTkeClient", sdkClient)
	patches.ApplyMethodFunc(sdkClient, "DescribeAddonValuesWithContext", func(ctx context.Context, request *tkev20180525.DescribeAddonValuesRequest) (*tkev20180525.DescribeAddonValuesResponse, error) {
		return nil, errors.New("ResourceUnavailable")
	})

	ds, schemaResp := buildAddonValuesDataSource(t)
	configureAddonValuesDataSource(t, ds)

	readResp, _ := readAddonValues(t, ds, schemaResp, "cls-5yezvaxo", "ip-masq-agent")
	assert.True(t, readResp.Diagnostics.HasError())
	assert.Contains(t, readResp.Diagnostics.Errors()[0].Detail(), "DescribeAddonValues failed")
	assert.Contains(t, readResp.Diagnostics.Errors()[0].Detail(), "ResourceUnavailable")
}

// TestAddonValues_Read_NilResponse covers a nil response payload.
func TestAddonValues_Read_NilResponse(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	sdkClient := &tkev20180525.Client{}
	patches.ApplyMethodReturn(&connectivity.TencentCloudClient{}, "UseTkeClient", sdkClient)
	patches.ApplyMethodFunc(sdkClient, "DescribeAddonValuesWithContext", func(ctx context.Context, request *tkev20180525.DescribeAddonValuesRequest) (*tkev20180525.DescribeAddonValuesResponse, error) {
		return nil, nil
	})

	ds, schemaResp := buildAddonValuesDataSource(t)
	configureAddonValuesDataSource(t, ds)

	readResp, _ := readAddonValues(t, ds, schemaResp, "cls-5yezvaxo", "ip-masq-agent")
	assert.True(t, readResp.Diagnostics.HasError())
	assert.Contains(t, readResp.Diagnostics.Errors()[0].Detail(), "empty response")
}
