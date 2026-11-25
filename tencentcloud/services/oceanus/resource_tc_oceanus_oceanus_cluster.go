package oceanus

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	oceanusv20190422 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/oceanus/v20190422"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/fw"
	fwHelper "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/helper"
)

// Compile-time interface assertions. The base resource.Resource interface is
// checked implicitly by the factory return value.
var (
	_ resource.ResourceWithConfigure   = &OceanusOceanusClusterResource{}
	_ resource.ResourceWithImportState = &OceanusOceanusClusterResource{}
)

// NewOceanusOceanusClusterResource is the factory referenced by
// tencentcloud/framework/registry.go to register this resource.
func NewOceanusOceanusClusterResource() resource.Resource {
	return &OceanusOceanusClusterResource{}
}

// OceanusOceanusClusterResource implements resource.Resource for
// tencentcloud_oceanus_oceanus_cluster.
type OceanusOceanusClusterResource struct {
	fw.ResourceWithConfigure
	fw.WithImportByID
}

// OceanusOceanusClusterResourceModel maps the schema attributes; each
// attribute appears exactly once, in schema order.
type OceanusOceanusClusterResourceModel struct {
	Id                          types.String `tfsdk:"id"`
	ClusterName                 types.String `tfsdk:"cluster_name"`
	RegionId                    types.Int64  `tfsdk:"region_id"`
	ZoneId                      types.Int64  `tfsdk:"zone_id"`
	LoginPassword               types.String `tfsdk:"login_password"`
	VpcDescriptions             types.List   `tfsdk:"vpc_descriptions"`
	DefaultCosBucket            types.String `tfsdk:"default_cos_bucket"`
	Cu                          types.Int64  `tfsdk:"cu"`
	Remark                      types.String `tfsdk:"remark"`
	Period                      types.Int64  `tfsdk:"period"`
	InstanceChargeType          types.String `tfsdk:"instance_charge_type"`
	ClusterType                 types.String `tfsdk:"cluster_type"`
	RenewFlag                   types.String `tfsdk:"renew_flag"`
	FlinkWebUiNetworkAccessType types.String `tfsdk:"flink_web_ui_network_access_type"`
	SlaveVpcDescriptions        types.List   `tfsdk:"slave_vpc_descriptions"`
	CuMemory                    types.Int64  `tfsdk:"cu_memory"`
}

// OceanusClusterVpcDescriptionModel maps one element of the
// vpc_descriptions / slave_vpc_descriptions nested attribute.
type OceanusClusterVpcDescriptionModel struct {
	VpcId    types.String `tfsdk:"vpc_id"`
	SubnetId types.String `tfsdk:"subnet_id"`
	AppId    types.Int64  `tfsdk:"app_id"`
	OwnerUin types.String `tfsdk:"owner_uin"`
}

func (r *OceanusOceanusClusterResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "tencentcloud_oceanus_oceanus_cluster"
}

func (r *OceanusOceanusClusterResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a resource to manage an Oceanus cluster.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "ID of the resource, the Oceanus cluster ID.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"cluster_name": schema.StringAttribute{
				Required:    true,
				Description: "Cluster name. Limits: 1 to 50 English letters, Chinese characters, digits, hyphens (-) or underscores (_).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"region_id": schema.Int64Attribute{
				Required:    true,
				Description: "Region ID. It can be queried through DescribeRegions of the region management system, with the Product parameter set to oceanus.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"zone_id": schema.Int64Attribute{
				Required:    true,
				Description: "Availability zone ID. It can be queried through DescribeZones of the region management system, with the Product parameter set to oceanus.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"login_password": schema.StringAttribute{
				Required:    true,
				Sensitive:   true,
				Description: "Flink UI access password. The user name and password are used to log in to the Flink UI of the cluster; the cluster user name is admin.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"vpc_descriptions": schema.ListNestedAttribute{
				Required:    true,
				Description: "Stream compute accesses other cloud product resources in the same region through VPC and elastic network interfaces, and occupies a certain number of subnet IPs. Make sure the selected subnet has enough available IPs. If the existing network does not meet your requirements, create a new private network or a new subnet in the VPC console.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: oceanusClusterVpcDescriptionAttributes(),
				},
				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplace(),
				},
			},
			"default_cos_bucket": schema.StringAttribute{
				Required:    true,
				Description: "Stream compute uses Cloud Object Storage (COS) to store job checkpoints, jar packages, or delivered logs. If there is no available bucket in the local region, create one in the COS console. To ensure normal use, refer to the documentation for the lifecycle configuration of the corresponding COS bucket.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			// cu is the cluster CU count on creation, and the target CU count of
			// ScaleOceanusCluster on update (the update API names the same
			// quantity NewCU). The cloud echoes the effective value back as
			// CuNum, so the attribute is Optional+Computed.
			"cu": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "Number of CUs of the cluster. 12 CU is the minimum computing resource and billing unit of stream compute; 1 CU contains 1 CPU and 4 GB of memory. When the number of CUs is greater than or equal to 48, the management node fee is waived.",
			},
			"remark": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Cluster description. Limits: 1 to 50 English letters, Chinese characters, digits, hyphens (-) or underscores (_).",
			},
			"period": schema.Int64Attribute{
				Optional:    true,
				Description: "Purchase duration in months. Value range: [1, 48].",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"instance_charge_type": schema.StringAttribute{
				Optional:    true,
				Description: "Cluster billing type. Valid values: `PREPAID` (monthly subscription), `POSTPAID_BY_SECOND` (pay-as-you-go). Default value: `POSTPAID_BY_SECOND`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"cluster_type": schema.StringAttribute{
				Optional:    true,
				Description: "Cluster type. Valid values: `MULTI_AZ_CLUSTER` (multi-availability-zone cluster).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"renew_flag": schema.StringAttribute{
				Optional:    true,
				Description: "Auto-renewal flag. Valid values: `NOTIFY_AND_MANUAL_RENEW` (notify and renew manually), `NOTIFY_AND_AUTO_RENEW` (notify and renew automatically), `DISABLE_NOTIFY_AND_MANUAL_RENEW` (do not notify and do not renew automatically). Default value: `NOTIFY_AND_MANUAL_RENEW`. When InstanceChargeType is set to `PREPAID`, the cluster is a monthly-subscription cluster and the auto-renewal flag must be set; pay-as-you-go clusters do not need it.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"flink_web_ui_network_access_type": schema.StringAttribute{
				Optional:    true,
				Description: "Flink UI access type. Valid values: `NetworkAccess_INTERNAL` (internal network access), `NetworkAccess_EXTERNAL` (public network access). Default value: `NetworkAccess_EXTERNAL`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"slave_vpc_descriptions": schema.ListNestedAttribute{
				Optional:    true,
				Description: "VPCs of the multi-availability-zone cluster.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: oceanusClusterVpcDescriptionAttributes(),
				},
				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplace(),
				},
			},
			"cu_memory": schema.Int64Attribute{
				Optional:    true,
				Description: "Memory-to-core ratio. Only [0, 2, 4, 8] are supported.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
		},
	}
}

// oceanusClusterVpcDescriptionAttributes is the nested attribute set shared by
// vpc_descriptions and slave_vpc_descriptions; both create-request structures
// (VPCDescription and SlaveVpcDescriptions) carry the same four fields.
func oceanusClusterVpcDescriptionAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"vpc_id": schema.StringAttribute{
			Required:    true,
			Description: "Private network ID.",
		},
		"subnet_id": schema.StringAttribute{
			Required:    true,
			Description: "Subnet ID.",
		},
		"app_id": schema.Int64Attribute{
			Optional:    true,
			Description: "User APPID.",
		},
		"owner_uin": schema.StringAttribute{
			Optional:    true,
			Description: "User UIN.",
		},
	}
}

// Create handles `terraform apply` for new resources.
func (r *OceanusOceanusClusterResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan OceanusOceanusClusterResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createOceanusClusterRequest := oceanusv20190422.NewCreateOceanusClusterRequest()
	createOceanusClusterRequest.ClusterName = fwHelper.StringPointerFromValue(plan.ClusterName)
	createOceanusClusterRequest.RegionId = fwHelper.Int64PointerFromValue(plan.RegionId)
	createOceanusClusterRequest.ZoneId = fwHelper.Int64PointerFromValue(plan.ZoneId)
	createOceanusClusterRequest.LoginPassword = fwHelper.StringPointerFromValue(plan.LoginPassword)
	createOceanusClusterRequest.VpcDescriptions = expandOceanusClusterVpcDescriptions(ctx, plan.VpcDescriptions, &resp.Diagnostics)
	createOceanusClusterRequest.DefaultCOSBucket = fwHelper.StringPointerFromValue(plan.DefaultCosBucket)
	createOceanusClusterRequest.CU = fwHelper.Int64PointerFromValue(plan.Cu)
	if !plan.Remark.IsNull() && !plan.Remark.IsUnknown() {
		createOceanusClusterRequest.Remark = plan.Remark.ValueStringPointer()
	}

	if !plan.Period.IsNull() && !plan.Period.IsUnknown() {
		createOceanusClusterRequest.Period = plan.Period.ValueInt64Pointer()
	}

	if !plan.InstanceChargeType.IsNull() && !plan.InstanceChargeType.IsUnknown() {
		createOceanusClusterRequest.InstanceChargeType = plan.InstanceChargeType.ValueStringPointer()
	}

	if !plan.ClusterType.IsNull() && !plan.ClusterType.IsUnknown() {
		createOceanusClusterRequest.ClusterType = plan.ClusterType.ValueStringPointer()
	}

	if !plan.RenewFlag.IsNull() && !plan.RenewFlag.IsUnknown() {
		createOceanusClusterRequest.RenewFlag = plan.RenewFlag.ValueStringPointer()
	}

	if !plan.FlinkWebUiNetworkAccessType.IsNull() && !plan.FlinkWebUiNetworkAccessType.IsUnknown() {
		createOceanusClusterRequest.FlinkWebUINetworkAccessType = plan.FlinkWebUiNetworkAccessType.ValueStringPointer()
	}

	createOceanusClusterRequest.SlaveVpcDescriptions = expandOceanusClusterSlaveVpcDescriptions(ctx, plan.SlaveVpcDescriptions, &resp.Diagnostics)
	createOceanusClusterRequest.CUMemory = fwHelper.Int64PointerFromValue(plan.CuMemory)

	if resp.Diagnostics.HasError() {
		return
	}

	service := NewOceanusService(r.Client())

	clusterId, err := service.CreateOceanusCluster(ctx, createOceanusClusterRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Oceanus cluster",
			fmt.Sprintf("CreateOceanusCluster failed: %v", err),
		)

		return
	}

	// Never write an empty id into state.
	if clusterId == "" {
		resp.Diagnostics.AddError(
			"Error creating Oceanus cluster",
			"CreateOceanusCluster returned an empty cluster id.",
		)

		return
	}

	plan.Id = types.StringValue(clusterId)

	// Computed attributes are unknown at this point and a state cannot hold
	// unknown values — normalize them to null before the defensive persist
	// below (see normalizeOceanusClusterComputedUnknowns).
	normalizeOceanusClusterComputedUnknowns(&plan)

	// Persist the identity before reading back so a partial failure still
	// leaves a destroyable resource in state.
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	describeClustersRequest := oceanusv20190422.NewDescribeClustersRequest()
	describeClustersRequest.ClusterIds = []*string{fwHelper.String(clusterId)}

	clusters, err := service.DescribeOceanusClusters(ctx, describeClustersRequest)
	if err != nil {
		resp.Diagnostics.AddError("Error reading Oceanus cluster after create", err.Error())
		return
	}

	item := findOceanusCluster(clusters, clusterId)
	if item == nil {
		resp.Diagnostics.AddError(
			"Oceanus cluster not readable after create",
			fmt.Sprintf("cluster %s was created but cannot be read back", plan.Id.ValueString()),
		)

		return
	}

	resp.Diagnostics.Append(refreshOceanusClusterModel(&plan, item)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Read handles `terraform refresh` and post-apply state refresh.
func (r *OceanusOceanusClusterResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state OceanusOceanusClusterResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	clusterId := state.Id.ValueString()

	describeClustersRequest := oceanusv20190422.NewDescribeClustersRequest()
	describeClustersRequest.ClusterIds = []*string{fwHelper.String(clusterId)}

	service := NewOceanusService(r.Client())

	clusters, err := service.DescribeOceanusClusters(ctx, describeClustersRequest)
	if err != nil {
		resp.Diagnostics.AddError("Error reading Oceanus cluster", err.Error())
		return
	}

	item := findOceanusCluster(clusters, clusterId)

	// Confirmed "not found" on the cloud side => drop it from state so
	// Terraform plans a recreation. The judgement stays inline here: a later
	// unconditional State.Set would overwrite the removal.
	if item == nil {
		resp.Diagnostics.AddWarning(
			"Resource not found",
			fmt.Sprintf("Oceanus cluster %s not found, removing from state.", state.Id.ValueString()),
		)

		resp.State.RemoveResource(ctx)
		return
	}

	resp.Diagnostics.Append(refreshOceanusClusterModel(&state, item)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update handles in-place modifications. Only the CU count is mutable: it is
// applied through ScaleOceanusCluster, every other attribute is
// RequiresReplace.
func (r *OceanusOceanusClusterResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan OceanusOceanusClusterResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state OceanusOceanusClusterResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	clusterId := state.Id.ValueString()

	service := NewOceanusService(r.Client())

	if !plan.Cu.Equal(state.Cu) {
		scaleOceanusClusterRequest := oceanusv20190422.NewScaleOceanusClusterRequest()
		scaleOceanusClusterRequest.ClusterId = fwHelper.String(clusterId)
		scaleOceanusClusterRequest.NewCU = fwHelper.Int64PointerFromValue(plan.Cu)
		// The scale API requires the direction explicitly; derive it from the
		// CU delta instead of exposing a separate attribute that could
		// contradict `cu`.
		if !plan.Cu.IsNull() && !plan.Cu.IsUnknown() && !state.Cu.IsNull() && !state.Cu.IsUnknown() && plan.Cu.ValueInt64() < state.Cu.ValueInt64() {
			scaleOceanusClusterRequest.ScaleMode = fwHelper.String("ScaleDown")
		} else {
			scaleOceanusClusterRequest.ScaleMode = fwHelper.String("ScaleUp")
		}

		if err := service.ScaleOceanusCluster(ctx, scaleOceanusClusterRequest); err != nil {
			resp.Diagnostics.AddError(
				"Error updating Oceanus cluster",
				fmt.Sprintf("ScaleOceanusCluster failed: %v", err),
			)

			return
		}
	}

	describeClustersRequest := oceanusv20190422.NewDescribeClustersRequest()
	describeClustersRequest.ClusterIds = []*string{fwHelper.String(clusterId)}

	clusters, err := service.DescribeOceanusClusters(ctx, describeClustersRequest)
	if err != nil {
		resp.Diagnostics.AddError("Error reading Oceanus cluster after update", err.Error())
		return
	}

	item := findOceanusCluster(clusters, clusterId)
	if item == nil {
		resp.Diagnostics.AddError(
			"Oceanus cluster not readable after update",
			fmt.Sprintf("cluster %s was updated but cannot be read back", state.Id.ValueString()),
		)

		return
	}

	plan.Id = state.Id
	resp.Diagnostics.Append(refreshOceanusClusterModel(&plan, item)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete handles resource deletion.
func (r *OceanusOceanusClusterResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state OceanusOceanusClusterResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	deleteOceanusClusterRequest := oceanusv20190422.NewDeleteOceanusClusterRequest()
	deleteOceanusClusterRequest.ClusterId = fwHelper.String(state.Id.ValueString())

	service := NewOceanusService(r.Client())

	if err := service.DeleteOceanusCluster(ctx, deleteOceanusClusterRequest); err != nil {
		resp.Diagnostics.AddError(
			"Error deleting Oceanus cluster",
			fmt.Sprintf("DeleteOceanusCluster failed: %v", err),
		)

		return
	}
}

// expandOceanusClusterVpcDescriptions converts the vpc_descriptions attribute
// into the SDK request structure.
func expandOceanusClusterVpcDescriptions(ctx context.Context, list types.List, diags *diag.Diagnostics) []*oceanusv20190422.VPCDescription {
	if list.IsNull() || list.IsUnknown() {
		return nil
	}

	var models []OceanusClusterVpcDescriptionModel
	diags.Append(list.ElementsAs(ctx, &models, false)...)
	if diags.HasError() {
		return nil
	}

	descriptions := make([]*oceanusv20190422.VPCDescription, 0, len(models))
	for _, model := range models {
		description := &oceanusv20190422.VPCDescription{}
		description.VpcId = fwHelper.StringPointerFromValue(model.VpcId)
		description.SubnetId = fwHelper.StringPointerFromValue(model.SubnetId)
		description.AppId = fwHelper.Int64PointerFromValue(model.AppId)
		description.OwnerUin = fwHelper.StringPointerFromValue(model.OwnerUin)
		descriptions = append(descriptions, description)
	}

	return descriptions
}

// expandOceanusClusterSlaveVpcDescriptions converts the
// slave_vpc_descriptions attribute into the SDK request structure.
func expandOceanusClusterSlaveVpcDescriptions(ctx context.Context, list types.List, diags *diag.Diagnostics) []*oceanusv20190422.SlaveVpcDescriptions {
	if list.IsNull() || list.IsUnknown() {
		return nil
	}

	var models []OceanusClusterVpcDescriptionModel
	diags.Append(list.ElementsAs(ctx, &models, false)...)
	if diags.HasError() {
		return nil
	}

	descriptions := make([]*oceanusv20190422.SlaveVpcDescriptions, 0, len(models))
	for _, model := range models {
		description := &oceanusv20190422.SlaveVpcDescriptions{}
		description.VpcId = fwHelper.StringPointerFromValue(model.VpcId)
		description.SubnetId = fwHelper.StringPointerFromValue(model.SubnetId)
		description.AppId = fwHelper.Int64PointerFromValue(model.AppId)
		description.OwnerUin = fwHelper.StringPointerFromValue(model.OwnerUin)
		descriptions = append(descriptions, description)
	}

	return descriptions
}

// refreshOceanusClusterModel copies cloud-owned fields into the model, keeping
// the user-supplied values for fields the describe API does not echo back.
func refreshOceanusClusterModel(m *OceanusOceanusClusterResourceModel, item *oceanusv20190422.Cluster) diag.Diagnostics {
	if item == nil {
		return nil
	}

	m.ClusterName = fwHelper.StringValueOrPrev(item.Name, m.ClusterName)
	m.DefaultCosBucket = fwHelper.StringValueOrPrev(item.DefaultCOSBucket, m.DefaultCosBucket)
	// CuNum is the effective CU count of the cluster (uint64 in the SDK).
	m.Cu = fwHelper.Int64ValueFromUintOrPrev(item.CuNum, m.Cu)
	m.Remark = fwHelper.StringValueOrPrev(item.Remark, m.Remark)

	return nil
}

// normalizeOceanusClusterComputedUnknowns rewrites unknown computed attributes
// to null so the model can be persisted into state before the read-back fills
// them with the cloud-side values.
func normalizeOceanusClusterComputedUnknowns(m *OceanusOceanusClusterResourceModel) {
	if m.Cu.IsUnknown() {
		m.Cu = types.Int64Null()
	}

	if m.Remark.IsUnknown() {
		m.Remark = types.StringNull()
	}
}

// findOceanusCluster picks this resource's own item from a describe list.
func findOceanusCluster(clusters []*oceanusv20190422.Cluster, clusterId string) *oceanusv20190422.Cluster {
	for _, c := range clusters {
		if c.ClusterId != nil && *c.ClusterId == clusterId {
			return c
		}
	}

	return nil
}
