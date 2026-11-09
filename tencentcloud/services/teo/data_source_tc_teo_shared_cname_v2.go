package teo

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/fw"
)

// _ datasource.DataSource is intentionally omitted: the factory function
// NewTeoSharedCnameV2DataSource returns datasource.DataSource, which already
// performs the implicit compile-time check at the return statement. Only the
// extended-interface assertion is kept here to verify the Configure method
// signature (promoted from fw.DataSourceWithConfigure) at compile time.
var _ datasource.DataSourceWithConfigure = &TeoSharedCnameV2DataSource{}

// NewTeoSharedCnameV2DataSource is the factory referenced by
// tencentcloud/framework/registry.go to register this data source.
func NewTeoSharedCnameV2DataSource() datasource.DataSource {
	return &TeoSharedCnameV2DataSource{}
}

// TeoSharedCnameV2DataSource implements datasource.DataSource for
// tencentcloud_teo_shared_cname_v2.
type TeoSharedCnameV2DataSource struct {
	fw.DataSourceWithConfigure
}

// TeoSharedCnameV2DataSourceModel maps the schema attributes.
type TeoSharedCnameV2DataSourceModel struct {
	ZoneId          types.String `tfsdk:"zone_id"`
	Filters         types.List   `tfsdk:"filters"`
	Direction       types.String `tfsdk:"direction"`
	Match           types.String `tfsdk:"match"`
	Order           types.String `tfsdk:"order"`
	SharedCnameInfo types.List   `tfsdk:"shared_cname_info"`
}

// teoSharedCnameV2FilterModel is the element type of filters.
type teoSharedCnameV2FilterModel struct {
	Name   types.String `tfsdk:"name"`
	Values types.Set    `tfsdk:"values"`
	Fuzzy  types.Bool   `tfsdk:"fuzzy"`
}

// TeoSharedCnameV2InfoModel is the element type of shared_cname_info.
type TeoSharedCnameV2InfoModel struct {
	Type                types.String `tfsdk:"type"`
	SharedCname         types.String `tfsdk:"shared_cname"`
	Description         types.String `tfsdk:"description"`
	IpsslConfig         types.Object `tfsdk:"ipssl_config"`
	BindDomainCount     types.Int64  `tfsdk:"bind_domain_count"`
	AccelerationDomains types.List   `tfsdk:"acceleration_domains"`
}

// TeoSharedCnameV2IpsslConfigModel is the object type of ipssl_config.
type TeoSharedCnameV2IpsslConfigModel struct {
	ZoneId           types.String `tfsdk:"zone_id"`
	AssociatedDomain types.String `tfsdk:"associated_domain"`
	Status           types.String `tfsdk:"status"`
}

// TeoSharedCnameV2AccelerationDomainModel is the element type of acceleration_domains.
type TeoSharedCnameV2AccelerationDomainModel struct {
	ZoneId   types.String `tfsdk:"zone_id"`
	Type     types.String `tfsdk:"type"`
	Instance types.String `tfsdk:"instance"`
}

// teoSharedCnameV2IpsslConfigAttrTypes returns the attribute types of ipssl_config.
func teoSharedCnameV2IpsslConfigAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"zone_id":           types.StringType,
		"associated_domain": types.StringType,
		"status":            types.StringType,
	}
}

// teoSharedCnameV2AccelerationDomainAttrTypes returns the attribute types of
// acceleration_domains elements.
func teoSharedCnameV2AccelerationDomainAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"zone_id":  types.StringType,
		"type":     types.StringType,
		"instance": types.StringType,
	}
}

// teoSharedCnameV2InfoAttrTypes returns the attribute types of
// shared_cname_info elements.
func teoSharedCnameV2InfoAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"type":                 types.StringType,
		"shared_cname":         types.StringType,
		"description":          types.StringType,
		"ipssl_config":         types.ObjectType{AttrTypes: teoSharedCnameV2IpsslConfigAttrTypes()},
		"bind_domain_count":    types.Int64Type,
		"acceleration_domains": types.ListType{ElemType: types.ObjectType{AttrTypes: teoSharedCnameV2AccelerationDomainAttrTypes()}},
	}
}

func (d *TeoSharedCnameV2DataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "tencentcloud_teo_shared_cname_v2"
}

func (d *TeoSharedCnameV2DataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this data source to query the shared CNAME list of TEO (TencentCloud EdgeOne) within a zone.",
		Attributes: map[string]schema.Attribute{
			"zone_id": schema.StringAttribute{
				Required:    true,
				Description: "Zone ID to which the shared CNAMEs belong.",
			},
			"filters": schema.ListNestedAttribute{
				Optional: true,
				Description: "Filter conditions of the shared CNAME list. The upper limit of Filters.Values is 20. " +
					"Valid filter names: `shared-cname` (filter by shared CNAME), `type` (filter by shared CNAME type), `description` (filter by description).",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Required:    true,
							Description: "Name of the field to filter by.",
						},
						"values": schema.SetAttribute{
							Required:    true,
							ElementType: types.StringType,
							Description: "Values of the filtered field.",
						},
						"fuzzy": schema.BoolAttribute{
							Optional:    true,
							Description: "Whether to enable fuzzy query.",
						},
					},
				},
			},
			"direction": schema.StringAttribute{
				Optional:    true,
				Description: "Sorting direction of the list. Valid values: `asc`, `desc`. Default value: `asc`.",
			},
			"match": schema.StringAttribute{
				Optional: true,
				Description: "Matching mode of multiple query conditions. Valid values: `all` (return shared CNAMEs matching all query conditions), " +
					"`any` (return shared CNAMEs matching any one of the query conditions). Default value: `all`.",
			},
			"order": schema.StringAttribute{
				Optional:    true,
				Description: "Sorting basis of the list. Valid values: `create-time` (creation time), `shared-cname` (shared CNAME). Default sorting is by the shared-cname attribute.",
			},
			"shared_cname_info": schema.ListNestedAttribute{
				Computed:    true,
				Description: "List of shared CNAMEs.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"type": schema.StringAttribute{
							Computed: true,
							Description: "Type of the shared CNAME. Valid values: `custom` (custom shared CNAME created by user), `ip-ssl` (IP SSL shared CNAME), " +
								"`zero-rating` (zero-rating shared CNAME), `preset` (preset resource shared CNAME).",
						},
						"shared_cname": schema.StringAttribute{
							Computed:    true,
							Description: "Name of the shared CNAME.",
						},
						"description": schema.StringAttribute{
							Computed:    true,
							Description: "Description of the shared CNAME.",
						},
						"ipssl_config": schema.SingleNestedAttribute{
							Computed: true,
							Description: "IP SSL configuration associated with the shared CNAME, returned only when type is `ip-ssl`. " +
								"When Status is `unbound`, zone_id and associated_domain are null.",
							Attributes: map[string]schema.Attribute{
								"zone_id": schema.StringAttribute{
									Computed:    true,
									Description: "Zone ID of the domain associated with the IP SSL configuration. Null when Status is `unbound`.",
								},
								"associated_domain": schema.StringAttribute{
									Computed:    true,
									Description: "Domain associated with the IP SSL configuration. Null when Status is `unbound`.",
								},
								"status": schema.StringAttribute{
									Computed: true,
									Description: "Association status of the IP SSL configuration. Valid values: `bound` (bound), `binding` (binding), " +
										"`unbinding` (unbinding), `unbound` (not bound).",
								},
							},
						},
						"bind_domain_count": schema.Int64Attribute{
							Computed:    true,
							Description: "Number of acceleration domains bound to the shared CNAME.",
						},
						"acceleration_domains": schema.ListNestedAttribute{
							Computed:    true,
							Description: "List of acceleration domains joined to the shared CNAME. Only the first 100 acceleration domains are returned when more than 100 domains are joined.",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"zone_id": schema.StringAttribute{
										Computed:    true,
										Description: "Zone ID.",
									},
									"type": schema.StringAttribute{
										Computed:    true,
										Description: "Instance type. Valid value: `acceleration-domain` (acceleration domain).",
									},
									"instance": schema.StringAttribute{
										Computed:    true,
										Description: "Information of the referenced/referencing instance.",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

// Read queries the shared CNAME list and writes the result into resp.State.
func (d *TeoSharedCnameV2DataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data TeoSharedCnameV2DataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if d.Client() == nil {
		resp.Diagnostics.AddError(
			"Provider not configured",
			"The provider client is not available. This is unexpected; please report it.",
		)
		return
	}

	request := teov20220901.NewDescribeSharedCNAMERequest()
	request.ZoneId = data.ZoneId.ValueStringPointer()

	if !data.Filters.IsNull() && !data.Filters.IsUnknown() {
		var filters []teoSharedCnameV2FilterModel
		resp.Diagnostics.Append(data.Filters.ElementsAs(ctx, &filters, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		for _, filter := range filters {
			values := make([]*string, 0, len(filter.Values.Elements()))
			for _, v := range filter.Values.Elements() {
				values = append(values, v.(types.String).ValueStringPointer())
			}
			advancedFilter := &teov20220901.AdvancedFilter{
				Name:   filter.Name.ValueStringPointer(),
				Values: values,
			}
			if !filter.Fuzzy.IsNull() && !filter.Fuzzy.IsUnknown() {
				advancedFilter.Fuzzy = filter.Fuzzy.ValueBoolPointer()
			}
			request.Filters = append(request.Filters, advancedFilter)
		}
	}

	if !data.Direction.IsNull() && !data.Direction.IsUnknown() {
		request.Direction = data.Direction.ValueStringPointer()
	}
	if !data.Match.IsNull() && !data.Match.IsUnknown() {
		request.Match = data.Match.ValueStringPointer()
	}
	if !data.Order.IsNull() && !data.Order.IsUnknown() {
		request.Order = data.Order.ValueStringPointer()
	}

	// Paginate through the full shared CNAME list so that the data source
	// always returns the complete result set regardless of page size.
	var offset int64 = 0
	var limit int64 = 200
	request.Limit = &limit
	request.Offset = &offset

	response, err := d.Client().UseTeoClient().DescribeSharedCNAMEWithContext(ctx, request)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error querying TEO shared CNAMEs",
			fmt.Sprintf("DescribeSharedCNAME failed: %v", err),
		)
		return
	}

	// Two different "empty" outcomes must stay distinguishable:
	//   - nil response  -> API failure -> error
	//   - len == 0      -> genuinely nothing matched -> empty list, normal return
	if response == nil || response.Response == nil {
		resp.Diagnostics.AddError(
			"Error querying TEO shared CNAMEs",
			"DescribeSharedCNAME returned an empty response.",
		)
		return
	}

	infos := make([]*teov20220901.SharedCNAMEInfo, 0, len(response.Response.SharedCNAMEInfo))
	infos = append(infos, response.Response.SharedCNAMEInfo...)
	for len(response.Response.SharedCNAMEInfo) >= int(limit) {
		offset += limit
		request.Offset = &offset
		response, err = d.Client().UseTeoClient().DescribeSharedCNAMEWithContext(ctx, request)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error querying TEO shared CNAMEs",
				fmt.Sprintf("DescribeSharedCNAME failed at offset %d: %v", offset, err),
			)
			return
		}
		if response == nil || response.Response == nil {
			resp.Diagnostics.AddError(
				"Error querying TEO shared CNAMEs",
				fmt.Sprintf("DescribeSharedCNAME returned an empty response at offset %d.", offset),
			)
			return
		}
		infos = append(infos, response.Response.SharedCNAMEInfo...)
	}

	items := make([]TeoSharedCnameV2InfoModel, 0, len(infos))
	for _, info := range infos {
		if info == nil {
			continue
		}
		item := TeoSharedCnameV2InfoModel{
			Type:            types.StringPointerValue(info.Type),
			SharedCname:     types.StringPointerValue(info.SharedCNAME),
			Description:     types.StringPointerValue(info.Description),
			BindDomainCount: types.Int64PointerValue(info.BindDomainCount),
		}

		if info.IPSSLConfig != nil {
			ipsslObj, diags := types.ObjectValueFrom(ctx, teoSharedCnameV2IpsslConfigAttrTypes(), TeoSharedCnameV2IpsslConfigModel{
				ZoneId:           types.StringPointerValue(info.IPSSLConfig.ZoneId),
				AssociatedDomain: types.StringPointerValue(info.IPSSLConfig.AssociatedDomain),
				Status:           types.StringPointerValue(info.IPSSLConfig.Status),
			})
			resp.Diagnostics.Append(diags...)
			if resp.Diagnostics.HasError() {
				return
			}
			item.IpsslConfig = ipsslObj
		} else {
			item.IpsslConfig = types.ObjectNull(teoSharedCnameV2IpsslConfigAttrTypes())
		}

		domains := make([]TeoSharedCnameV2AccelerationDomainModel, 0, len(info.AccelerationDomains))
		for _, domain := range info.AccelerationDomains {
			if domain == nil {
				continue
			}
			domains = append(domains, TeoSharedCnameV2AccelerationDomainModel{
				ZoneId:   types.StringPointerValue(domain.ZoneId),
				Type:     types.StringPointerValue(domain.Type),
				Instance: types.StringPointerValue(domain.Instance),
			})
		}
		domainList, diags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: teoSharedCnameV2AccelerationDomainAttrTypes()}, domains)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		item.AccelerationDomains = domainList

		items = append(items, item)
	}

	infoList, diags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: teoSharedCnameV2InfoAttrTypes()}, items)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.SharedCnameInfo = infoList
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
