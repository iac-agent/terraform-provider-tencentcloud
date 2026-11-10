package vpc

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/fw"
)

var _ datasource.DataSourceWithConfigure = &VpcSecurityGroupLimits111001DataSource{}

// NewVpcSecurityGroupLimits111001DataSource is the factory referenced by
// tencentcloud/framework/registry.go to register this data source.
func NewVpcSecurityGroupLimits111001DataSource() datasource.DataSource {
	return &VpcSecurityGroupLimits111001DataSource{}
}

// VpcSecurityGroupLimits111001DataSource implements datasource.DataSource for
// the tencentcloud_vpc_security_group_limits_1110_01 data source, which reports
// the security group quotas of the current account in the configured region.
type VpcSecurityGroupLimits111001DataSource struct {
	fw.DataSourceWithConfigure
}

// VpcSecurityGroupLimits111001DataSourceModel maps the schema attributes.
type VpcSecurityGroupLimits111001DataSourceModel struct {
	SecurityGroupLimitSet types.List `tfsdk:"security_group_limit_set"`
}

// vpcSecurityGroupLimitItemModel is the element type of security_group_limit_set.
type vpcSecurityGroupLimitItemModel struct {
	SecurityGroupLimit                 types.Int64 `tfsdk:"security_group_limit"`
	SecurityGroupPolicyLimit           types.Int64 `tfsdk:"security_group_policy_limit"`
	ReferedSecurityGroupLimit          types.Int64 `tfsdk:"refered_security_group_limit"`
	SecurityGroupInstanceLimit         types.Int64 `tfsdk:"security_group_instance_limit"`
	InstanceSecurityGroupLimit         types.Int64 `tfsdk:"instance_security_group_limit"`
	SecurityGroupExtendedPolicyLimit   types.Int64 `tfsdk:"security_group_extended_policy_limit"`
	SecurityGroupReferedCvmAndEniLimit types.Int64 `tfsdk:"security_group_refered_cvm_and_eni_limit"`
	SecurityGroupReferedSvcLimit       types.Int64 `tfsdk:"security_group_refered_svc_limit"`
}

// vpcSecurityGroupLimitItemAttrTypes mirrors vpcSecurityGroupLimitItemModel and
// is used to build the typed list written back to state.
var vpcSecurityGroupLimitItemAttrTypes = map[string]attr.Type{
	"security_group_limit":                     types.Int64Type,
	"security_group_policy_limit":              types.Int64Type,
	"refered_security_group_limit":             types.Int64Type,
	"security_group_instance_limit":            types.Int64Type,
	"instance_security_group_limit":            types.Int64Type,
	"security_group_extended_policy_limit":     types.Int64Type,
	"security_group_refered_cvm_and_eni_limit": types.Int64Type,
	"security_group_refered_svc_limit":         types.Int64Type,
}

func (d *VpcSecurityGroupLimits111001DataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "tencentcloud_vpc_security_group_limits_1110_01"
}

func (d *VpcSecurityGroupLimits111001DataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this data source to query the security group quota limits of the current TencentCloud account in the configured region.",
		Attributes: map[string]schema.Attribute{
			"security_group_limit_set": schema.ListNestedAttribute{
				Computed:    true,
				Description: "Security group quota limits of the account.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"security_group_limit": schema.Int64Attribute{
							Computed:    true,
							Description: "Number of security groups that can be created per project per region.",
						},
						"security_group_policy_limit": schema.Int64Attribute{
							Computed:    true,
							Description: "Maximum number of rules in a single security group.",
						},
						"refered_security_group_limit": schema.Int64Attribute{
							Computed:    true,
							Description: "Maximum number of nested security group rules in a security group. The name keeps the server side spelling `Refered`.",
						},
						"security_group_instance_limit": schema.Int64Attribute{
							Computed:    true,
							Description: "Maximum number of instances associated with a single security group.",
						},
						"instance_security_group_limit": schema.Int64Attribute{
							Computed:    true,
							Description: "Maximum number of security groups associated with a single instance.",
						},
						"security_group_extended_policy_limit": schema.Int64Attribute{
							Computed:    true,
							Description: "Limit of the expanded rules of a security group. It differs from `security_group_policy_limit` in that the latter counts the rules as written in the security group while this one counts the rules after referenced security groups are expanded.",
						},
						"security_group_refered_cvm_and_eni_limit": schema.Int64Attribute{
							Computed:    true,
							Description: "Quota of CVM and ENI instances associated with a referenced security group. The name keeps the server side spelling `Refered`.",
						},
						"security_group_refered_svc_limit": schema.Int64Attribute{
							Computed:    true,
							Description: "Quota of database, CLB and other service instances associated with a referenced security group. The name keeps the server side spelling `Refered`.",
						},
					},
				},
			},
		},
	}
}

// Read queries DescribeSecurityGroupLimits and writes the whole result into
// resp.State. A query without filter arguments always returns a single
// quota record; a missing record (nil SecurityGroupLimitSet) means the API
// returned no payload, hence it is reported as an error instead of being
// written as an empty list.
func (d *VpcSecurityGroupLimits111001DataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	// The data source exposes no configurable argument, so there is nothing to
	// read from the config; the result model is built from scratch below.
	var model VpcSecurityGroupLimits111001DataSourceModel

	service := NewVpcService(d.Client())
	limitSet, err := service.DescribeVpcSecurityGroupLimits(ctx, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error querying VPC security group limits",
			fmt.Sprintf("DescribeSecurityGroupLimits failed: %v", err),
		)
		return
	}
	if limitSet == nil {
		resp.Diagnostics.AddError(
			"Error querying VPC security group limits",
			"DescribeSecurityGroupLimits returned an empty response.",
		)
		return
	}

	items := []vpcSecurityGroupLimitItemModel{
		{
			SecurityGroupLimit:                 types.Int64PointerValue(limitSet.SecurityGroupLimit),
			SecurityGroupPolicyLimit:           types.Int64PointerValue(limitSet.SecurityGroupPolicyLimit),
			ReferedSecurityGroupLimit:          types.Int64PointerValue(limitSet.ReferedSecurityGroupLimit),
			SecurityGroupInstanceLimit:         types.Int64PointerValue(limitSet.SecurityGroupInstanceLimit),
			InstanceSecurityGroupLimit:         types.Int64PointerValue(limitSet.InstanceSecurityGroupLimit),
			SecurityGroupExtendedPolicyLimit:   types.Int64PointerValue(limitSet.SecurityGroupExtendedPolicyLimit),
			SecurityGroupReferedCvmAndEniLimit: types.Int64PointerValue(limitSet.SecurityGroupReferedCvmAndEniLimit),
			SecurityGroupReferedSvcLimit:       types.Int64PointerValue(limitSet.SecurityGroupReferedSvcLimit),
		},
	}

	list, diags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: vpcSecurityGroupLimitItemAttrTypes}, items)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	model.SecurityGroupLimitSet = list
	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}
