package tke

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	tke "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tke/v20180525"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/sharedmeta"
)

var _ datasource.DataSource = &TkeAddonValuesDataSource{}
var _ datasource.DataSourceWithConfigure = &TkeAddonValuesDataSource{}

// NewTkeAddonValuesDataSource is the factory referenced by
// tencentcloud/framework/registry.go to register this data source.
func NewTkeAddonValuesDataSource() datasource.DataSource {
	return &TkeAddonValuesDataSource{}
}

// TkeAddonValuesDataSource implements datasource.DataSource for
// tencentcloud_tke_addon_values.
type TkeAddonValuesDataSource struct {
	client TkeService
}

// TkeAddonValuesDataSourceModel maps the schema attributes.
type TkeAddonValuesDataSourceModel struct {
	ClusterId     types.String `tfsdk:"cluster_id"`
	AddonName     types.String `tfsdk:"addon_name"`
	Values        types.String `tfsdk:"values"`
	DefaultValues types.String `tfsdk:"default_values"`
}

func (d *TkeAddonValuesDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "tencentcloud_tke_addon_values"
}

func (d *TkeAddonValuesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this data source to query the parameter values of a TKE cluster addon. " +
			"It returns the parameter values currently rendered for the addon and the default " +
			"parameter values supported by the addon chart, both as JSON strings.",
		Attributes: map[string]schema.Attribute{
			"cluster_id": schema.StringAttribute{
				Required:    true,
				Description: "ID of the TKE cluster the addon belongs to.",
			},
			"addon_name": schema.StringAttribute{
				Required:    true,
				Description: "Name of the addon, as listed in the cluster's addon management page.",
			},
			"values": schema.StringAttribute{
				Computed:    true,
				Description: "Parameter values of the addon rendered with the values already set, as a JSON string. Empty or null when the addon is not installed.",
			},
			"default_values": schema.StringAttribute{
				Computed:    true,
				Description: "Default parameter values supported by the addon, as a JSON string.",
			},
		},
	}
}

// Configure receives the *sharedmeta.ProviderMeta populated by
// framework.Provider.Configure (which itself reuses the SDKv2-built
// *connectivity.TencentCloudClient). The SDKv2-style TkeService is then
// constructed from the shared client, mirroring how SDKv2 resources in this
// package construct it.
func (d *TkeAddonValuesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	meta, ok := req.ProviderData.(*sharedmeta.ProviderMeta)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected ProviderData type",
			fmt.Sprintf("Expected *sharedmeta.ProviderMeta, got: %T. "+
				"This is a bug in the provider; please report it.", req.ProviderData),
		)
		return
	}

	d.client = NewTkeService(meta.Client)
}

// Read queries the addon parameter values and writes them back to state.
func (d *TkeAddonValuesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config TkeAddonValuesDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	clusterId := config.ClusterId.ValueString()
	addonName := config.AddonName.ValueString()

	if clusterId == "" {
		resp.Diagnostics.AddError("Missing cluster_id", "cluster_id is required and cannot be empty.")
		return
	}
	if addonName == "" {
		resp.Diagnostics.AddError("Missing addon_name", "addon_name is required and cannot be empty.")
		return
	}

	values, defaultValues, err := d.client.DescribeAddonValues(ctx, clusterId, addonName)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error querying TKE addon values",
			fmt.Sprintf("DescribeAddonValues failed for cluster %q addon %q: %v", clusterId, addonName, err),
		)
		return
	}

	fillTkeAddonValues(&config, &tke.DescribeAddonValuesResponseParams{
		Values:        values,
		DefaultValues: defaultValues,
	})

	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}

// fillTkeAddonValues maps an API response into the data source model. A nil
// SDK pointer becomes null rather than an empty string, so that an addon which
// is not installed (empty Values) can be told apart from an addon installed
// with empty parameters.
func fillTkeAddonValues(config *TkeAddonValuesDataSourceModel, in *tke.DescribeAddonValuesResponseParams) {
	if in == nil {
		config.Values = types.StringNull()
		config.DefaultValues = types.StringNull()
		return
	}

	config.Values = types.StringPointerValue(in.Values)
	config.DefaultValues = types.StringPointerValue(in.DefaultValues)
}
