package ags

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	agsv20250920 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ags/v20250920"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/fw"
	fwhelper "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

var _ datasource.DataSourceWithConfigure = &AgsGPUResourceSpecsDataSource{}

// NewAgsGPUResourceSpecsDataSource is the factory referenced by
// tencentcloud/framework/registry.go to register this data source.
func NewAgsGPUResourceSpecsDataSource() datasource.DataSource {
	return &AgsGPUResourceSpecsDataSource{}
}

// AgsGPUResourceSpecsDataSource implements datasource.DataSource for
// the tencentcloud_ags_gpu_resource_specs data source.
type AgsGPUResourceSpecsDataSource struct {
	fw.DataSourceWithConfigure
}

// AgsGPUResourceSpecsDataSourceModel maps the schema attributes.
type AgsGPUResourceSpecsDataSourceModel struct {
	GPUType types.String `tfsdk:"gpu_type"`
	CPU     types.Int64  `tfsdk:"cpu"`
	Memory  types.Int64  `tfsdk:"memory"`
	GPUNum  types.Int64  `tfsdk:"gpu_num"`
	Offset  types.Int64  `tfsdk:"offset"`
	Limit   types.Int64  `tfsdk:"limit"`

	TotalCount types.Int64 `tfsdk:"total_count"`
	Specs      types.List  `tfsdk:"specs"`
}

// AgsGPUResourceSpecModel is the element type of specs.
type AgsGPUResourceSpecModel struct {
	GPUType        types.String `tfsdk:"gpu_type"`
	GPUTypeAlias   types.String `tfsdk:"gpu_type_alias"`
	GPUMemoryGiB   types.Int64  `tfsdk:"gpu_memory_gib"`
	CPU            types.Int64  `tfsdk:"cpu"`
	Memory         types.Int64  `tfsdk:"memory"`
	GPUNum         types.Int64  `tfsdk:"gpu_num"`
	AvailableCount types.Int64  `tfsdk:"available_count"`
}

// AgsGPUResourceSpecAttrTypes is the object type of the specs elements.
var AgsGPUResourceSpecAttrTypes = map[string]attr.Type{
	"gpu_type":        types.StringType,
	"gpu_type_alias":  types.StringType,
	"gpu_memory_gib":  types.Int64Type,
	"cpu":             types.Int64Type,
	"memory":          types.Int64Type,
	"gpu_num":         types.Int64Type,
	"available_count": types.Int64Type,
}

func (d *AgsGPUResourceSpecsDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "tencentcloud_ags_gpu_resource_specs"
}

func (d *AgsGPUResourceSpecsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this data source to query the sellable GPU resource specs and the available quota of the current tenant in the current region.",
		Attributes: map[string]schema.Attribute{
			// ---- Query inputs ----
			"gpu_type": schema.StringAttribute{
				Optional:    true,
				Description: "GPU type identifier, for example `5090d`. When it is specified, `cpu`, `memory` and `gpu_num` must be specified together.",
			},
			"cpu": schema.Int64Attribute{
				Optional:    true,
				Description: "Number of CPU cores, must be positive. It must be specified together with `gpu_type`, `memory` and `gpu_num`.",
			},
			"memory": schema.Int64Attribute{
				Optional:    true,
				Description: "Host memory in GiB, must be positive. It must be specified together with `gpu_type`, `cpu` and `gpu_num`.",
			},
			"gpu_num": schema.Int64Attribute{
				Optional:    true,
				Description: "Number of GPU cards per instance, must be a positive integer. It must be specified together with `gpu_type`, `cpu` and `memory`.",
			},
			"offset": schema.Int64Attribute{
				Optional:    true,
				Description: "Pagination offset, default is 0, must be greater than or equal to 0.",
			},
			"limit": schema.Int64Attribute{
				Optional:    true,
				Description: "Page size, 100 is used when it is omitted or 0; a non-zero value must be in the range 1-500.",
			},
			// ---- Query outputs: Computed ----
			"total_count": schema.Int64Attribute{
				Computed:    true,
				Description: "Total number of specs after filtering and before pagination.",
			},
			"specs": schema.ListNestedAttribute{
				Computed:    true,
				Description: "List of sellable GPU resource specs and quota. It is empty when no sellable spec exists.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"gpu_type": schema.StringAttribute{
							Computed:    true,
							Description: "GPU type identifier, for example `5090d`.",
						},
						"gpu_type_alias": schema.StringAttribute{
							Computed:    true,
							Description: "Display name of the GPU type. `gpu_type` is used for display when it is empty.",
						},
						"gpu_memory_gib": schema.Int64Attribute{
							Computed:    true,
							Description: "Total GPU memory of all GPUs on a single instance, in GiB. Do not multiply it by `gpu_num`. 0 means the display metadata is not maintained.",
						},
						"cpu": schema.Int64Attribute{
							Computed:    true,
							Description: "Number of CPU cores of a single instance.",
						},
						"memory": schema.Int64Attribute{
							Computed:    true,
							Description: "Host memory of a single instance in GiB, which is not GPU memory.",
						},
						"gpu_num": schema.Int64Attribute{
							Computed:    true,
							Description: "Number of GPU cards of a single instance.",
						},
						"available_count": schema.Int64Attribute{
							Computed:    true,
							Description: "Currently purchasable count, in instances rather than GPU cards. 0 means no quota and null means the quota is unknown. The quota is verified again at purchase time.",
						},
					},
				},
			},
		},
	}
}

// Read queries the cloud API and writes the result into resp.State.
func (d *AgsGPUResourceSpecsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config AgsGPUResourceSpecsDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	request := agsv20250920.NewDescribeGPUResourceSpecsRequest()
	if !config.GPUType.IsNull() && !config.GPUType.IsUnknown() {
		request.GPUType = config.GPUType.ValueStringPointer()
	}
	if !config.CPU.IsNull() && !config.CPU.IsUnknown() {
		request.CPU = fwhelper.Uint64PointerFromValue(config.CPU)
	}
	if !config.Memory.IsNull() && !config.Memory.IsUnknown() {
		request.Memory = fwhelper.Uint64PointerFromValue(config.Memory)
	}
	if !config.GPUNum.IsNull() && !config.GPUNum.IsUnknown() {
		request.GPUNum = fwhelper.Uint64PointerFromValue(config.GPUNum)
	}
	if !config.Offset.IsNull() && !config.Offset.IsUnknown() {
		request.Offset = helper.Int64(config.Offset.ValueInt64())
	}
	if !config.Limit.IsNull() && !config.Limit.IsUnknown() {
		request.Limit = helper.Int64(config.Limit.ValueInt64())
	}

	service := NewAgsService(d.Client())
	specItems, total, err := service.DescribeGPUResourceSpecs(ctx, request)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error querying AGS GPU resource specs",
			fmt.Sprintf("DescribeGPUResourceSpecs failed: %v", err),
		)
		return
	}

	// Two different "empty" outcomes must stay distinguishable:
	//   - nil response error is reported by the service layer above.
	//   - len == 0 is a genuinely empty result -> empty list, normal return.
	specs := make([]AgsGPUResourceSpecModel, 0, len(specItems))
	for _, item := range specItems {
		if item == nil {
			continue
		}
		specs = append(specs, AgsGPUResourceSpecModel{
			GPUType:        types.StringPointerValue(item.GPUType),
			GPUTypeAlias:   types.StringPointerValue(item.GPUTypeAlias),
			GPUMemoryGiB:   fwhelper.Int64ValueFromUint(item.GPUMemoryGiB),
			CPU:            fwhelper.Int64ValueFromUint(item.CPU),
			Memory:         fwhelper.Int64ValueFromUint(item.Memory),
			GPUNum:         fwhelper.Int64ValueFromUint(item.GPUNum),
			AvailableCount: fwhelper.Int64ValueFromUint(item.AvailableCount),
		})
	}

	list, diags := types.ListValueFrom(ctx, types.ObjectType{
		AttrTypes: AgsGPUResourceSpecAttrTypes,
	}, specs)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	config.Specs = list
	config.TotalCount = types.Int64Value(total)

	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}
