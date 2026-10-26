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
		Description: "Use this data source to query TKE addon parameter values.",
		Attributes: map[string]schema.Attribute{
			"cluster_id": schema.StringAttribute{
				Required:    true,
				Description: "ID of the TKE cluster.",
			},
			"addon_name": schema.StringAttribute{
				Required:    true,
				Description: "Name of the addon, which can be obtained from the component management page of the cluster.",
			},
			"values": schema.StringAttribute{
				Computed: true,
				Description: "Parameter JSON string of the addon. When the addon is installed, the user-set parameters " +
					"are rendered together with the chart default parameters; when the addon is not installed, the value is empty.",
			},
			"default_values": schema.StringAttribute{
				Computed:    true,
				Description: "Parameter JSON string supported by the addon, with values taken from the chart defaults.",
			},
		},
	}
}

// Configure receives the *sharedmeta.ProviderMeta populated by
// framework.Provider.Configure (which reuses the SDKv2-built
// *connectivity.TencentCloudClient) and builds the SDKv2-style TkeService
// from the shared client, mirroring how SDKv2 references in this package do it.
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

// Read queries the addon parameter values through DescribeAddonValues.
func (d *TkeAddonValuesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config TkeAddonValuesDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	request := tke.NewDescribeAddonValuesRequest()
	request.ClusterId = config.ClusterId.ValueStringPointer()
	request.AddonName = config.AddonName.ValueStringPointer()

	response, err := d.client.DescribeAddonValues(ctx, request)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error querying TKE addon values",
			fmt.Sprintf("DescribeAddonValues failed for cluster %s addon %s: %v",
				config.ClusterId.ValueString(), config.AddonName.ValueString(), err),
		)
		return
	}

	if response == nil || response.Response == nil {
		resp.Diagnostics.AddError(
			"Error querying TKE addon values",
			fmt.Sprintf("DescribeAddonValues returned an empty response for cluster %s addon %s",
				config.ClusterId.ValueString(), config.AddonName.ValueString()),
		)
		return
	}

	config.Values = types.StringPointerValue(response.Response.Values)
	config.DefaultValues = types.StringPointerValue(response.Response.DefaultValues)

	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}
