package oceanus

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	common "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	oceanus "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/oceanus/v20190422"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/fw"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

var (
	_ resource.ResourceWithConfigure   = &OceanusOceanusClusterResource{}
	_ resource.ResourceWithImportState = &OceanusOceanusClusterResource{}
)

func NewOceanusOceanusClusterResource() resource.Resource {
	return &OceanusOceanusClusterResource{}
}

type OceanusOceanusClusterResource struct {
	fw.ResourceWithConfigure
	fw.WithNoUpdate
	fw.WithImportByID
}

type OceanusOceanusClusterVpcModel struct {
	VpcId    types.String `tfsdk:"vpc_id"`
	SubnetId types.String `tfsdk:"subnet_id"`
	AppId    types.Int64  `tfsdk:"app_id"`
	OwnerUin types.String `tfsdk:"owner_uin"`
}

type OceanusOceanusClusterResourceModel struct {
	Id                          types.String `tfsdk:"id"`
	ClusterName                 types.String `tfsdk:"cluster_name"`
	RegionId                    types.Int64  `tfsdk:"region_id"`
	ZoneId                      types.Int64  `tfsdk:"zone_id"`
	LoginPassword               types.String `tfsdk:"login_password"`
	VpcDescriptions             types.List   `tfsdk:"vpc_descriptions"`
	DefaultCOSBucket            types.String `tfsdk:"default_cos_bucket"`
	Cu                          types.Int64  `tfsdk:"cu"`
	Remark                      types.String `tfsdk:"remark"`
	Period                      types.Int64  `tfsdk:"period"`
	InstanceChargeType          types.String `tfsdk:"instance_charge_type"`
	ClusterType                 types.String `tfsdk:"cluster_type"`
	RenewFlag                   types.String `tfsdk:"renew_flag"`
	FlinkWebUINetworkAccessType types.String `tfsdk:"flink_web_ui_network_access_type"`
	SlaveVpcDescriptions        types.List   `tfsdk:"slave_vpc_descriptions"`
	CuMemory                    types.Int64  `tfsdk:"cu_memory"`
	Status                      types.Int64  `tfsdk:"status"`
	CreateTime                  types.String `tfsdk:"create_time"`
	UpdateTime                  types.String `tfsdk:"update_time"`
}

func vpcDescriptionsAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"vpc_id":    types.StringType,
		"subnet_id": types.StringType,
		"app_id":    types.Int64Type,
		"owner_uin": types.StringType,
	}
}

func (r *OceanusOceanusClusterResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "tencentcloud_oceanus_oceanus_cluster"
}

func (r *OceanusOceanusClusterResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a resource to create an Oceanus cluster.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Cluster ID in the format `cluster-xxx`.",
			},
			"cluster_name": schema.StringAttribute{
				Required:    true,
				Description: "Cluster name. Supports 1-50 English letters, Chinese characters, digits, hyphens or underscores.",
			},
			"region_id": schema.Int64Attribute{
				Required:    true,
				Description: "Region ID, which can be queried via DescribeRegions with Product set to oceanus.",
			},
			"zone_id": schema.Int64Attribute{
				Required:    true,
				Description: "Zone ID, which can be queried via DescribeZones with Product set to oceanus.",
			},
			"login_password": schema.StringAttribute{
				Required:    true,
				Sensitive:   true,
				Description: "Flink UI access password. The cluster username is `admin`.",
			},
			"vpc_descriptions": schema.ListNestedAttribute{
				Required:    true,
				Description: "VPC descriptions. The stream computing cluster accesses other cloud resources in the same region via VPC and ENI, which occupies subnet IPs.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"vpc_id": schema.StringAttribute{
							Required:    true,
							Description: "VPC ID.",
						},
						"subnet_id": schema.StringAttribute{
							Required:    true,
							Description: "Subnet ID.",
						},
						"app_id": schema.Int64Attribute{
							Optional:    true,
							Computed:    true,
							Description: "User AppId.",
						},
						"owner_uin": schema.StringAttribute{
							Optional:    true,
							Computed:    true,
							Description: "User UIN.",
						},
					},
				},
			},
			"default_cos_bucket": schema.StringAttribute{
				Required:    true,
				Description: "Default COS bucket used by the cluster to save checkpoints, jar packages, or delivery logs.",
			},
			"cu": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "Number of cluster CUs. 12 CU is the minimum computing resource and billing unit. 1 CU contains 1 CPU and 4GB memory. Management node fee is waived when CU is greater than or equal to 48.",
			},
			"remark": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Cluster description. Supports 1-50 English letters, Chinese characters, digits, hyphens or underscores.",
			},
			"period": schema.Int64Attribute{
				Optional:    true,
				Description: "Purchase duration in months. Valid range: [1, 48]. Required when instance_charge_type is PREPAID.",
			},
			"instance_charge_type": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Cluster billing type. Valid values: `PREPAID`, `POSTPAID_BY_SECOND`. Default value: `POSTPAID_BY_SECOND`.",
			},
			"cluster_type": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Cluster type. Valid value: `MULTI_AZ_CLUSTER`.",
			},
			"renew_flag": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Auto renewal flag. Valid values: `NOTIFY_AND_MANUAL_RENEW`, `NOTIFY_AND_AUTO_RENEW`, `DISABLE_NOTIFY_AND_MANUAL_RENEW`. Default value: `NOTIFY_AND_MANUAL_RENEW`. Required for PREPAID clusters.",
			},
			"flink_web_ui_network_access_type": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Flink UI access type. Valid values: `NetworkAccess_INTERNAL`, `NetworkAccess_EXTERNAL`. Default value: `NetworkAccess_EXTERNAL`.",
			},
			"slave_vpc_descriptions": schema.ListNestedAttribute{
				Optional:    true,
				Description: "Multi-AZ VPC descriptions.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"vpc_id": schema.StringAttribute{
							Required:    true,
							Description: "VPC ID.",
						},
						"subnet_id": schema.StringAttribute{
							Required:    true,
							Description: "Subnet ID.",
						},
						"app_id": schema.Int64Attribute{
							Optional:    true,
							Computed:    true,
							Description: "User AppId.",
						},
						"owner_uin": schema.StringAttribute{
							Optional:    true,
							Computed:    true,
							Description: "User UIN.",
						},
					},
				},
			},
			"cu_memory": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "Core memory ratio. Only supports [0, 2, 4, 8].",
			},
			"status": schema.Int64Attribute{
				Computed:    true,
				Description: "Cluster status. `1`: uninitialized, `3`: initializing, `2`: running.",
			},
			"create_time": schema.StringAttribute{
				Computed:    true,
				Description: "Cluster creation time.",
			},
			"update_time": schema.StringAttribute{
				Computed:    true,
				Description: "Time of the last operation on the cluster.",
			},
		},
	}
}

func (r *OceanusOceanusClusterResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan OceanusOceanusClusterResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if r.Client() == nil {
		resp.Diagnostics.AddError("Provider not configured", "The provider client is not available. This is unexpected; please report it.")
		return
	}

	request := oceanus.NewCreateOceanusClusterRequest()
	request.ClusterName = plan.ClusterName.ValueStringPointer()
	request.RegionId = int64Pointer(plan.RegionId.ValueInt64())
	request.ZoneId = int64Pointer(plan.ZoneId.ValueInt64())
	request.LoginPassword = plan.LoginPassword.ValueStringPointer()
	request.DefaultCOSBucket = plan.DefaultCOSBucket.ValueStringPointer()

	if !plan.VpcDescriptions.IsNull() && !plan.VpcDescriptions.IsUnknown() {
		request.VpcDescriptions = expandVpcDescriptions(ctx, plan.VpcDescriptions, &resp.Diagnostics)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	if !plan.Cu.IsNull() && !plan.Cu.IsUnknown() {
		request.CU = int64Pointer(plan.Cu.ValueInt64())
	}
	if !plan.Remark.IsNull() && !plan.Remark.IsUnknown() {
		request.Remark = plan.Remark.ValueStringPointer()
	}
	if !plan.Period.IsNull() && !plan.Period.IsUnknown() {
		request.Period = int64Pointer(plan.Period.ValueInt64())
	}
	if !plan.InstanceChargeType.IsNull() && !plan.InstanceChargeType.IsUnknown() {
		request.InstanceChargeType = plan.InstanceChargeType.ValueStringPointer()
	}
	if !plan.ClusterType.IsNull() && !plan.ClusterType.IsUnknown() {
		request.ClusterType = plan.ClusterType.ValueStringPointer()
	}
	if !plan.RenewFlag.IsNull() && !plan.RenewFlag.IsUnknown() {
		request.RenewFlag = plan.RenewFlag.ValueStringPointer()
	}
	if !plan.FlinkWebUINetworkAccessType.IsNull() && !plan.FlinkWebUINetworkAccessType.IsUnknown() {
		request.FlinkWebUINetworkAccessType = plan.FlinkWebUINetworkAccessType.ValueStringPointer()
	}
	if !plan.SlaveVpcDescriptions.IsNull() && !plan.SlaveVpcDescriptions.IsUnknown() {
		request.SlaveVpcDescriptions = expandSlaveVpcDescriptions(ctx, plan.SlaveVpcDescriptions, &resp.Diagnostics)
		if resp.Diagnostics.HasError() {
			return
		}
	}
	if !plan.CuMemory.IsNull() && !plan.CuMemory.IsUnknown() {
		request.CUMemory = int64Pointer(plan.CuMemory.ValueInt64())
	}

	ratelimit.Check(request.GetAction())

	response, err := r.Client().UseOceanusClient().CreateOceanusCluster(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", "codegen", request.GetAction(), request.ToJsonString(), err.Error())
		resp.Diagnostics.AddError("Error creating Oceanus cluster", err.Error())
		return
	}
	log.Printf("[DEBUG] api[%s] success, request body [%s], response body [%s]\n", request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || response.Response == nil || response.Response.ClusterId == nil {
		resp.Diagnostics.AddError("Error creating Oceanus cluster", "the cloud API did not return a cluster id")
		return
	}

	plan.Id = types.StringValue(*response.Response.ClusterId)
	plan.Read(ctx, r.Client(), &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *OceanusOceanusClusterResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state OceanusOceanusClusterResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if r.Client() == nil {
		resp.Diagnostics.AddError("Provider not configured", "The provider client is not available. This is unexpected; please report it.")
		return
	}

	found := state.Read(ctx, r.Client(), &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	if !found {
		resp.Diagnostics.AddWarning("Oceanus cluster not found", fmt.Sprintf("the Oceanus cluster %s no longer exists and will be removed from state", state.Id.ValueString()))
		resp.State.RemoveResource(ctx)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *OceanusOceanusClusterResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state OceanusOceanusClusterResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if r.Client() == nil {
		resp.Diagnostics.AddError("Provider not configured", "The provider client is not available. This is unexpected; please report it.")
		return
	}

	request := oceanus.NewDeleteOceanusClusterRequest()
	request.ClusterId = state.Id.ValueStringPointer()

	ratelimit.Check(request.GetAction())

	response, err := r.Client().UseOceanusClient().DeleteOceanusCluster(request)
	if err != nil {
		log.Printf("[CRITAL] api[%s] fail, request body [%s], reason[%s]\n", request.GetAction(), request.ToJsonString(), err.Error())
		resp.Diagnostics.AddError("Error deleting Oceanus cluster", err.Error())
		return
	}
	log.Printf("[DEBUG] api[%s] success, request body [%s], response body [%s]\n", request.GetAction(), request.ToJsonString(), response.ToJsonString())
}

func expandVpcDescriptions(ctx context.Context, list types.List, diags *diag.Diagnostics) []*oceanus.VPCDescription {
	var models []OceanusOceanusClusterVpcModel
	diags.Append(list.ElementsAs(ctx, &models, false)...)
	if diags.HasError() {
		return nil
	}
	result := make([]*oceanus.VPCDescription, 0, len(models))
	for _, m := range models {
		item := &oceanus.VPCDescription{
			VpcId:    m.VpcId.ValueStringPointer(),
			SubnetId: m.SubnetId.ValueStringPointer(),
		}
		if !m.AppId.IsNull() && !m.AppId.IsUnknown() {
			item.AppId = int64Pointer(m.AppId.ValueInt64())
		}
		if !m.OwnerUin.IsNull() && !m.OwnerUin.IsUnknown() {
			item.OwnerUin = m.OwnerUin.ValueStringPointer()
		}
		result = append(result, item)
	}
	return result
}

func expandSlaveVpcDescriptions(ctx context.Context, list types.List, diags *diag.Diagnostics) []*oceanus.SlaveVpcDescriptions {
	var models []OceanusOceanusClusterVpcModel
	diags.Append(list.ElementsAs(ctx, &models, false)...)
	if diags.HasError() {
		return nil
	}
	result := make([]*oceanus.SlaveVpcDescriptions, 0, len(models))
	for _, m := range models {
		item := &oceanus.SlaveVpcDescriptions{
			VpcId:    m.VpcId.ValueStringPointer(),
			SubnetId: m.SubnetId.ValueStringPointer(),
		}
		if !m.AppId.IsNull() && !m.AppId.IsUnknown() {
			item.AppId = int64Pointer(m.AppId.ValueInt64())
		}
		if !m.OwnerUin.IsNull() && !m.OwnerUin.IsUnknown() {
			item.OwnerUin = m.OwnerUin.ValueStringPointer()
		}
		result = append(result, item)
	}
	return result
}

func (m *OceanusOceanusClusterResourceModel) Read(ctx context.Context, client interface{ UseOceanusClient() *oceanus.Client }, diags *diag.Diagnostics) bool {
	request := oceanus.NewDescribeClustersRequest()
	request.ClusterIds = common.StringPtrs([]string{m.Id.ValueString()})

	ratelimit.Check(request.GetAction())

	response, err := client.UseOceanusClient().DescribeClusters(request)
	if err != nil {
		diags.AddError("Error reading Oceanus cluster", err.Error())
		return false
	}

	if response == nil || response.Response == nil || len(response.Response.ClusterSet) != 1 {
		return false
	}
	cluster := response.Response.ClusterSet[0]

	m.Id = stringPointerValue(cluster.ClusterId)
	m.ClusterName = stringPointerValue(cluster.Name)
	m.Remark = stringPointerValue(cluster.Remark)
	m.DefaultCOSBucket = stringPointerValue(cluster.DefaultCOSBucket)
	m.Cu = uint64PointerValue(cluster.CuNum)
	m.CuMemory = uint64PointerValue(cluster.CuMem)
	m.Status = int64PointerValue(cluster.Status)
	m.CreateTime = stringPointerValue(cluster.CreateTime)
	m.UpdateTime = stringPointerValue(cluster.UpdateTime)

	if len(cluster.CCNs) > 0 {
		vpcModels := make([]OceanusOceanusClusterVpcModel, 0, len(cluster.CCNs))
		for _, ccn := range cluster.CCNs {
			vpcModels = append(vpcModels, OceanusOceanusClusterVpcModel{
				VpcId:    stringPointerValue(ccn.VpcId),
				SubnetId: stringPointerValue(ccn.SubnetId),
			})
		}
		vpcList, listDiags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: vpcDescriptionsAttrTypes()}, vpcModels)
		diags.Append(listDiags...)
		if diags.HasError() {
			return false
		}
		m.VpcDescriptions = vpcList
	}

	return true
}

func int64Pointer(v int64) *int64 {
	return &v
}

func int64PointerValue(v *int64) types.Int64 {
	if v == nil {
		return types.Int64Null()
	}
	return types.Int64Value(*v)
}

func uint64PointerValue(v *uint64) types.Int64 {
	if v == nil {
		return types.Int64Null()
	}
	return types.Int64Value(int64(*v))
}

func stringPointerValue(v *string) types.String {
	if v == nil {
		return types.StringNull()
	}
	return types.StringValue(*v)
}
