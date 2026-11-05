package teo

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/fw"
)

var _ datasource.DataSourceWithConfigure = &TeoInferenceAPITokensDataSource{}

// NewTeoInferenceAPITokensDataSource is the factory referenced by
// tencentcloud/framework/registry.go to register this data source.
func NewTeoInferenceAPITokensDataSource() datasource.DataSource {
	return &TeoInferenceAPITokensDataSource{}
}

// TeoInferenceAPITokensDataSource implements datasource.DataSource for
// tencentcloud_teo_inference_api_tokens.
type TeoInferenceAPITokensDataSource struct {
	fw.DataSourceWithConfigure
}

// TeoInferenceAPITokensDataSourceModel maps the schema attributes.
type TeoInferenceAPITokensDataSourceModel struct {
	ZoneId types.String `tfsdk:"zone_id"`
	Tokens types.List   `tfsdk:"tokens"`
}

// TeoInferenceAPITokenModel is the element type of tokens.
type TeoInferenceAPITokenModel struct {
	TokenId    types.String `tfsdk:"token_id"`
	Name       types.String `tfsdk:"name"`
	Content    types.String `tfsdk:"content"`
	CreateTime types.String `tfsdk:"create_time"`
}

// teoInferenceAPITokenAttrTypes describes the object type of a tokens element.
var teoInferenceAPITokenAttrTypes = map[string]attr.Type{
	"token_id":    types.StringType,
	"name":        types.StringType,
	"content":     types.StringType,
	"create_time": types.StringType,
}

func (d *TeoInferenceAPITokensDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "tencentcloud_teo_inference_api_tokens"
}

func (d *TeoInferenceAPITokensDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this data source to query the inference API tokens of the specified TEO zone.",
		Attributes: map[string]schema.Attribute{
			// ---- Query inputs ----
			"zone_id": schema.StringAttribute{
				Required:    true,
				Description: "Zone ID.",
			},
			// ---- Query outputs ----
			"tokens": schema.ListNestedAttribute{
				Computed:    true,
				Description: "List of inference API tokens.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"token_id": schema.StringAttribute{
							Computed:    true,
							Description: "Inference API Token ID.",
						},
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "Inference API Token name.",
						},
						"content": schema.StringAttribute{
							Computed:    true,
							Sensitive:   true,
							Description: "Inference API Token content.",
						},
						"create_time": schema.StringAttribute{
							Computed:    true,
							Description: "Creation time in ISO date format.",
						},
					},
				},
			},
		},
	}
}

// Read queries the cloud API and writes the result into resp.State.
func (d *TeoInferenceAPITokensDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config TeoInferenceAPITokensDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	zoneId := config.ZoneId.ValueString()
	if zoneId == "" {
		resp.Diagnostics.AddError(
			"Missing zone_id",
			"zone_id is required and cannot be empty.",
		)
		return
	}

	if d.Client() == nil {
		resp.Diagnostics.AddError(
			"Provider not configured",
			"The provider client is not available. This is unexpected; please report it.",
		)
		return
	}

	service := NewTeoService(d.Client())
	tokens, err := service.DescribeTeoInferenceAPITokens(ctx, zoneId)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("querying TEO inference API tokens for zone (%s)", zoneId),
			err.Error(),
		)
		return
	}

	items := make([]TeoInferenceAPITokenModel, 0, len(tokens))
	for _, token := range tokens {
		if token == nil {
			continue
		}
		items = append(items, TeoInferenceAPITokenModel{
			TokenId:    types.StringPointerValue(token.TokenId),
			Name:       types.StringPointerValue(token.Name),
			Content:    types.StringPointerValue(token.Content),
			CreateTime: types.StringPointerValue(token.CreateTime),
		})
	}

	list, diags := types.ListValueFrom(ctx, types.ObjectType{
		AttrTypes: teoInferenceAPITokenAttrTypes,
	}, items)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	config.Tokens = list
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}
