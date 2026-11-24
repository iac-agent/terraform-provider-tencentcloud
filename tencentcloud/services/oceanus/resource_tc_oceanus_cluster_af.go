package oceanus

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	resourcetimeouts "github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"

	oceanusv20190422 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/oceanus/v20190422"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/fw"
	fwHelper "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/helper"
)

// Compile-time interface assertions.
var (
	_ resource.ResourceWithConfigure   = &OceanusClusterAfResource{}
	_ resource.ResourceWithImportState = &OceanusClusterAfResource{}
)

// NewOceanusClusterAfResource is the factory referenced by
// tencentcloud/framework/registry.go to register this resource.
func NewOceanusClusterAfResource() resource.Resource {
	return &OceanusClusterAfResource{}
}

// OceanusClusterAfResource implements resource.Resource for
// tencentcloud_oceanus_cluster_af.
//
// The cloud side addresses a cluster by a single id (cluster-xxx), so the
// Terraform id is that value and import passes it straight into the state.
type OceanusClusterAfResource struct {
	fw.ResourceWithConfigure
	fw.WithImportByID
}

// OceanusClusterAfResourceModel maps the schema attributes; each attribute
// appears exactly once.
type OceanusClusterAfResourceModel struct {
	Id                          types.String           `tfsdk:"id"`
	Timeouts                    resourcetimeouts.Value `tfsdk:"timeouts"`
	ClusterName                 types.String           `tfsdk:"cluster_name"`
	RegionId                    types.Int64            `tfsdk:"region_id"`
	ZoneId                      types.Int64            `tfsdk:"zone_id"`
	LoginPassword               types.String           `tfsdk:"login_password"`
	VpcDescriptions             types.List             `tfsdk:"vpc_descriptions"`
	DefaultCosBucket            types.String           `tfsdk:"default_cos_bucket"`
	Cu                          types.Int64            `tfsdk:"cu"`
	Remark                      types.String           `tfsdk:"remark"`
	Period                      types.Int64            `tfsdk:"period"`
	InstanceChargeType          types.String           `tfsdk:"instance_charge_type"`
	ClusterType                 types.String           `tfsdk:"cluster_type"`
	RenewFlag                   types.String           `tfsdk:"renew_flag"`
	FlinkWebUiNetworkAccessType types.String           `tfsdk:"flink_web_ui_network_access_type"`
	SlaveVpcDescriptions        types.List             `tfsdk:"slave_vpc_descriptions"`
	CuMemory                    types.Int64            `tfsdk:"cu_memory"`
}

// oceanusClusterAfVpcDescriptionModel maps one element of vpc_descriptions and
// slave_vpc_descriptions; both API structures carry the same fields.
type oceanusClusterAfVpcDescriptionModel struct {
	VpcId    types.String `tfsdk:"vpc_id"`
	SubnetId types.String `tfsdk:"subnet_id"`
	AppId    types.Int64  `tfsdk:"app_id"`
	OwnerUin types.String `tfsdk:"owner_uin"`
}

func (r *OceanusClusterAfResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "tencentcloud_oceanus_cluster_af"
}

func (r *OceanusClusterAfResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	vpcDescriptionAttributes := map[string]schema.Attribute{
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
			Description: "User AppId. Note: This field may return null, indicating that no valid value can be obtained.",
		},
		"owner_uin": schema.StringAttribute{
			Optional:    true,
			Description: "User UIN. Note: This field may return null, indicating that no valid value can be obtained.",
		},
	}

	resp.Schema = schema.Schema{
		Description: "Provides a resource to manage an Oceanus cluster.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "ID of the cluster.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"cluster_name": schema.StringAttribute{
				Required: true,
				Description: "Cluster name. Input restrictions: 1-50 English letters, Chinese characters, digits, " +
					"hyphens (-) or underscores (_) are supported.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"region_id": schema.Int64Attribute{
				Required: true,
				Description: "Region ID. It can be queried through DescribeRegions of the region management system " +
					"with the Product parameter set to oceanus.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"zone_id": schema.Int64Attribute{
				Required: true,
				Description: "Availability zone ID. It can be queried through DescribeZones of the region management " +
					"system with the Product parameter set to oceanus.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"login_password": schema.StringAttribute{
				Required:  true,
				Sensitive: true,
				Description: "Flink UI access password. The username and password are used to log in to the Flink UI " +
					"of the job; the cluster username is admin.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"vpc_descriptions": schema.ListNestedAttribute{
				Required: true,
				Description: "Stream compute accesses other cloud product resources in the same region through VPC " +
					"and elastic network interfaces, and occupies a certain number of subnet IPs. Please ensure that " +
					"the selected subnet has enough available IPs. If the existing network does not meet your " +
					"requirements, create a new VPC or subnet in the VPC console.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: vpcDescriptionAttributes,
				},
				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplace(),
				},
			},
			"default_cos_bucket": schema.StringAttribute{
				Required: true,
				Description: "Stream compute uses Cloud Object Storage (COS) to store job checkpoints, jar packages or " +
					"delivered logs. If there is no available bucket in the local region, create one in the COS " +
					"console. To ensure normal usage, refer to the document for the lifecycle configuration of the " +
					"corresponding COS bucket.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"cu": schema.Int64Attribute{
				Optional: true,
				Computed: true,
				Description: "Number of cluster CUs. 12 CU is the minimum computing resource and billing unit of " +
					"stream compute; 1 CU contains 1 CPU and 4 GB of memory. The management node fee is waived when " +
					"the number of CUs is greater than or equal to 48.",
			},
			"remark": schema.StringAttribute{
				Optional: true,
				Description: "Cluster description. Input restrictions: 1-50 English letters, Chinese characters, " +
					"digits, hyphens (-) or underscores (_) are supported.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"period": schema.Int64Attribute{
				Optional:    true,
				Description: "Purchase duration in months. Value range: [1, 48].",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"instance_charge_type": schema.StringAttribute{
				Optional: true,
				Description: "Cluster billing type. Valid values: PREPAID (monthly subscription) and " +
					"POSTPAID_BY_SECOND (pay-as-you-go). Default value: POSTPAID_BY_SECOND.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"cluster_type": schema.StringAttribute{
				Optional:    true,
				Description: "Cluster type. Valid values: MULTI_AZ_CLUSTER (multi availability zone cluster).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"renew_flag": schema.StringAttribute{
				Optional: true,
				Description: "Auto-renewal flag. Valid values: NOTIFY_AND_MANUAL_RENEW (notify and renew manually), " +
					"NOTIFY_AND_AUTO_RENEW (notify and renew automatically) and DISABLE_NOTIFY_AND_MANUAL_RENEW " +
					"(neither notify nor renew automatically). Default value: NOTIFY_AND_MANUAL_RENEW. The flag is " +
					"required for monthly subscription clusters (InstanceChargeType set to PREPAID) and is not " +
					"required for pay-as-you-go clusters.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"flink_web_ui_network_access_type": schema.StringAttribute{
				Optional: true,
				Description: "Flink UI access type. Valid values: NetworkAccess_INTERNAL (intranet access) and " +
					"NetworkAccess_EXTERNAL (public network access). Default value: NetworkAccess_EXTERNAL.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"slave_vpc_descriptions": schema.ListNestedAttribute{
				Optional:    true,
				Description: "Multi availability zone VPC.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: vpcDescriptionAttributes,
				},
				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplace(),
				},
			},
			// TODO(human): CreateOceanusCluster.CUMemory is the "core memory ratio"
			// ([0, 2, 4, 8]) while DescribeClusters only returns CuMem ("CU memory
			// spec"); confirm whether the two describe the same setting before
			// mapping CuMem back into cu_memory for drift detection. Until then the
			// attribute is create-only (RequiresReplace) and is not refreshed.
			"cu_memory": schema.Int64Attribute{
				Optional:    true,
				Description: "Core memory ratio. Only [0, 2, 4, 8] are supported.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
		},
		Blocks: map[string]schema.Block{
			"timeouts": oceanusClusterAfTimeoutsBlock(ctx),
		},
	}
}

// defaultOceanusClusterAfAsyncTimeout is the built-in maximum duration to wait
// for an asynchronous cluster creation when the user does not override it
// through the `timeouts` block, e.g. `timeouts { create = "30m" }`.
const defaultOceanusClusterAfAsyncTimeout = 30 * time.Minute

// Values accepted by ScaleOceanusCluster.ScaleMode.
const (
	oceanusClusterAfScaleUp   = "ScaleUp"
	oceanusClusterAfScaleDown = "ScaleDown"
)

// oceanusClusterAfTimeoutsBlock builds the `timeouts` block from the official
// terraform-plugin-framework-timeouts module. Only the create stage waits: the
// scale and destroy APIs are synchronous.
func oceanusClusterAfTimeoutsBlock(ctx context.Context) schema.Block {
	block := resourcetimeouts.Block(ctx, resourcetimeouts.Opts{
		Create: true,
		CreateDescription: "A string that can be parsed as a duration (https://pkg.go.dev/time#ParseDuration), " +
			"e.g. `30m`. It limits how long to wait for the asynchronous cluster creation task. Default is `30m`.",
	})

	if nested, ok := block.(schema.SingleNestedBlock); ok {
		nested.Description = "The timeouts block allows you to specify the timeout for the create operation."
		return nested
	}

	return block
}

// Create handles `terraform apply` for new resources.
func (r *OceanusClusterAfResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan OceanusClusterAfResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createOceanusClusterRequest := oceanusv20190422.NewCreateOceanusClusterRequest()
	createOceanusClusterRequest.ClusterName = plan.ClusterName.ValueStringPointer()
	createOceanusClusterRequest.RegionId = fwHelper.Int64PointerFromValue(plan.RegionId)
	createOceanusClusterRequest.ZoneId = fwHelper.Int64PointerFromValue(plan.ZoneId)
	createOceanusClusterRequest.LoginPassword = plan.LoginPassword.ValueStringPointer()
	createOceanusClusterRequest.DefaultCOSBucket = plan.DefaultCosBucket.ValueStringPointer()

	vpcDescriptions, vpcDiags := expandOceanusClusterAfVpcDescriptions(ctx, plan.VpcDescriptions)
	resp.Diagnostics.Append(vpcDiags...)
	if resp.Diagnostics.HasError() {
		return
	}
	createOceanusClusterRequest.VpcDescriptions = vpcDescriptions

	slaveVpcDescriptions, slaveVpcDiags := expandOceanusClusterAfSlaveVpcDescriptions(ctx, plan.SlaveVpcDescriptions)
	resp.Diagnostics.Append(slaveVpcDiags...)
	if resp.Diagnostics.HasError() {
		return
	}
	createOceanusClusterRequest.SlaveVpcDescriptions = slaveVpcDescriptions

	createOceanusClusterRequest.CU = fwHelper.Int64PointerFromValue(plan.Cu)
	createOceanusClusterRequest.Remark = fwHelper.StringPointerFromValue(plan.Remark)
	createOceanusClusterRequest.Period = fwHelper.Int64PointerFromValue(plan.Period)
	createOceanusClusterRequest.InstanceChargeType = fwHelper.StringPointerFromValue(plan.InstanceChargeType)
	createOceanusClusterRequest.ClusterType = fwHelper.StringPointerFromValue(plan.ClusterType)
	createOceanusClusterRequest.RenewFlag = fwHelper.StringPointerFromValue(plan.RenewFlag)
	createOceanusClusterRequest.FlinkWebUINetworkAccessType = fwHelper.StringPointerFromValue(plan.FlinkWebUiNetworkAccessType)
	createOceanusClusterRequest.CUMemory = fwHelper.Int64PointerFromValue(plan.CuMemory)

	service := NewOceanusService(r.Client())

	// The timeout is resolved BEFORE the write API, so a malformed duration
	// fails the apply without any cloud-side effect.
	timeout, timeoutDiags := plan.Timeouts.Create(ctx, defaultOceanusClusterAfAsyncTimeout)
	resp.Diagnostics.Append(timeoutDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	clusterId, err := service.CreateOceanusCluster(ctx, createOceanusClusterRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Oceanus cluster",
			fmt.Sprintf("CreateOceanusCluster failed: %v", err),
		)

		return
	}

	if clusterId == "" {
		resp.Diagnostics.AddError(
			"Error creating Oceanus cluster",
			"CreateOceanusCluster returned an empty cluster id.",
		)

		return
	}

	plan.Id = types.StringValue(clusterId)

	normalizeOceanusClusterAfComputedUnknowns(&plan)

	// Persist the identity before waiting so a partial failure still leaves a
	// destroyable resource in state.
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// CreateOceanusCluster is asynchronous and returns no task id: the cluster
	// state itself is polled through DescribeClusters (the interface named by
	// the create entry's async_info).
	if err := service.WaitOceanusClusterReady(ctx, clusterId, timeout); err != nil {
		resp.Diagnostics.AddError(
			"Error waiting for Oceanus cluster creation",
			fmt.Sprintf("cluster %s was created but never became ready: %v", clusterId, err),
		)

		return
	}

	describeClustersRequest := oceanusv20190422.NewDescribeClustersRequest()
	describeClustersRequest.ClusterIds = []*string{&clusterId}
	clusters, err := service.DescribeOceanusClusters(ctx, describeClustersRequest)
	if err != nil {
		resp.Diagnostics.AddError("Error reading Oceanus cluster after create", err.Error())
		return
	}

	item := findOceanusClusterAf(clusters, clusterId)
	if item == nil {
		resp.Diagnostics.AddError(
			"Oceanus cluster not readable after create",
			fmt.Sprintf("cluster %s was created but cannot be read back", clusterId),
		)

		return
	}

	resp.Diagnostics.Append(refreshOceanusClusterAfModel(&plan, item)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Read handles `terraform refresh` and post-apply state refresh.
func (r *OceanusClusterAfResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state OceanusClusterAfResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	clusterId := state.Id.ValueString()

	describeClustersRequest := oceanusv20190422.NewDescribeClustersRequest()
	describeClustersRequest.ClusterIds = []*string{&clusterId}

	service := NewOceanusService(r.Client())

	clusters, err := service.DescribeOceanusClusters(ctx, describeClustersRequest)
	if err != nil {
		resp.Diagnostics.AddError("Error reading Oceanus cluster", err.Error())
		return
	}

	item := findOceanusClusterAf(clusters, clusterId)

	// Confirmed "not found" on the cloud side => drop it from state so
	// Terraform plans a recreation. The judgement stays inline: a later
	// unconditional State.Set would overwrite the removal.
	if item == nil {
		resp.Diagnostics.AddWarning(
			"Resource not found",
			fmt.Sprintf("Oceanus cluster %s not found, removing from state.", clusterId),
		)

		resp.State.RemoveResource(ctx)
		return
	}

	resp.Diagnostics.Append(refreshOceanusClusterAfModel(&state, item)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update handles in-place modifications. Only the CU number is updatable
// (ScaleOceanusCluster); every other attribute carries RequiresReplace.
func (r *OceanusClusterAfResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan OceanusClusterAfResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state OceanusClusterAfResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	clusterId := state.Id.ValueString()

	service := NewOceanusService(r.Client())

	// A plan differing from the state in the `timeouts` block alone also
	// reaches Update; the scale API must not be called then.
	if oceanusClusterAfCuChanged(plan.Cu, state.Cu) {
		scaleOceanusClusterRequest := oceanusv20190422.NewScaleOceanusClusterRequest()
		scaleOceanusClusterRequest.ClusterId = &clusterId
		scaleOceanusClusterRequest.NewCU = fwHelper.Int64PointerFromValue(plan.Cu)
		scaleOceanusClusterRequest.ScaleMode = oceanusClusterAfScaleMode(state.Cu, plan.Cu)

		if err := service.ScaleOceanusCluster(ctx, scaleOceanusClusterRequest); err != nil {
			resp.Diagnostics.AddError(
				"Error scaling Oceanus cluster",
				fmt.Sprintf("ScaleOceanusCluster failed for cluster %s: %v", clusterId, err),
			)

			return
		}
	}

	describeClustersRequest := oceanusv20190422.NewDescribeClustersRequest()
	describeClustersRequest.ClusterIds = []*string{&clusterId}
	clusters, err := service.DescribeOceanusClusters(ctx, describeClustersRequest)
	if err != nil {
		resp.Diagnostics.AddError("Error reading Oceanus cluster after update", err.Error())
		return
	}

	item := findOceanusClusterAf(clusters, clusterId)
	if item == nil {
		resp.Diagnostics.AddError(
			"Oceanus cluster not readable after update",
			fmt.Sprintf("cluster %s was updated but cannot be read back", clusterId),
		)

		return
	}

	plan.Id = state.Id
	resp.Diagnostics.Append(refreshOceanusClusterAfModel(&plan, item)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete handles resource deletion.
func (r *OceanusClusterAfResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state OceanusClusterAfResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	clusterId := state.Id.ValueString()

	deleteOceanusClusterRequest := oceanusv20190422.NewDeleteOceanusClusterRequest()
	deleteOceanusClusterRequest.ClusterId = &clusterId

	service := NewOceanusService(r.Client())

	if err := service.DeleteOceanusCluster(ctx, deleteOceanusClusterRequest); err != nil {
		resp.Diagnostics.AddError(
			"Error deleting Oceanus cluster",
			fmt.Sprintf("DeleteOceanusCluster failed for cluster %s: %v", clusterId, err),
		)

		return
	}
}

// refreshOceanusClusterAfModel copies the cloud-owned fields into the model.
//
// User-configurable attributes keep their planned value when the Describe API
// omits them (*OrPrev), so a configured value is never nulled out after apply.
// Attributes the Describe API does not echo back in the same form (region_id,
// zone_id, login_password, cu_memory, the VPC descriptions and the purchase
// options) are left untouched so the plan/state values survive the refresh.
func refreshOceanusClusterAfModel(m *OceanusClusterAfResourceModel, item *oceanusv20190422.Cluster) diag.Diagnostics {
	var diags diag.Diagnostics

	if item == nil {
		return diags
	}

	m.ClusterName = fwHelper.StringValueOrPrev(item.Name, m.ClusterName)
	m.Remark = fwHelper.StringValueOrPrev(item.Remark, m.Remark)
	m.DefaultCosBucket = fwHelper.StringValueOrPrev(item.DefaultCOSBucket, m.DefaultCosBucket)
	m.Cu = fwHelper.Int64ValueOrPrev(item.CuNum, m.Cu)

	return diags
}

// normalizeOceanusClusterAfComputedUnknowns rewrites unknown computed
// attributes to null so the model can be persisted into state before the
// read-back fills them with the cloud-side values. Only `cu` is
// Optional+Computed; every other attribute is either required or not computed
// at all, so a plan value is never unknown for them.
func normalizeOceanusClusterAfComputedUnknowns(m *OceanusClusterAfResourceModel) {
	if m.Cu.IsUnknown() {
		m.Cu = types.Int64Null()
	}
}

// oceanusClusterAfCuChanged reports whether the planned CU number is a usable
// value that differs from the one recorded in state.
func oceanusClusterAfCuChanged(plan, state types.Int64) bool {
	if plan.IsNull() || plan.IsUnknown() {
		return false
	}

	return !plan.Equal(state)
}

// oceanusClusterAfScaleMode derives the ScaleOceanusCluster.ScaleMode value
// from the CU numbers: shrinking requests ScaleDown, everything else ScaleUp
// (the API default).
func oceanusClusterAfScaleMode(state, plan types.Int64) *string {
	mode := oceanusClusterAfScaleUp
	if !state.IsNull() && !state.IsUnknown() && plan.ValueInt64() < state.ValueInt64() {
		mode = oceanusClusterAfScaleDown
	}

	return fwHelper.StringPointerFromValue(types.StringValue(mode))
}

// findOceanusClusterAf picks this resource's own item from a describe list.
func findOceanusClusterAf(clusters []*oceanusv20190422.Cluster, clusterId string) *oceanusv20190422.Cluster {
	for _, cluster := range clusters {
		if cluster != nil && cluster.ClusterId != nil && *cluster.ClusterId == clusterId {
			return cluster
		}
	}

	return nil
}

// expandOceanusClusterAfVpcDescriptions converts the planned VPC descriptions
// into the SDK slice.
func expandOceanusClusterAfVpcDescriptions(ctx context.Context, list types.List) ([]*oceanusv20190422.VPCDescription, diag.Diagnostics) {
	var diags diag.Diagnostics
	if list.IsNull() || list.IsUnknown() {
		return nil, diags
	}

	var descriptions []oceanusClusterAfVpcDescriptionModel
	diags.Append(list.ElementsAs(ctx, &descriptions, false)...)
	if diags.HasError() {
		return nil, diags
	}

	out := make([]*oceanusv20190422.VPCDescription, 0, len(descriptions))
	for _, description := range descriptions {
		out = append(out, &oceanusv20190422.VPCDescription{
			VpcId:    fwHelper.StringPointerFromValue(description.VpcId),
			SubnetId: fwHelper.StringPointerFromValue(description.SubnetId),
			AppId:    fwHelper.Int64PointerFromValue(description.AppId),
			OwnerUin: fwHelper.StringPointerFromValue(description.OwnerUin),
		})
	}

	return out, diags
}

// expandOceanusClusterAfSlaveVpcDescriptions converts the planned multi-AZ VPC
// descriptions into the SDK slice.
func expandOceanusClusterAfSlaveVpcDescriptions(ctx context.Context, list types.List) ([]*oceanusv20190422.SlaveVpcDescriptions, diag.Diagnostics) {
	var diags diag.Diagnostics
	if list.IsNull() || list.IsUnknown() {
		return nil, diags
	}

	var descriptions []oceanusClusterAfVpcDescriptionModel
	diags.Append(list.ElementsAs(ctx, &descriptions, false)...)
	if diags.HasError() {
		return nil, diags
	}

	out := make([]*oceanusv20190422.SlaveVpcDescriptions, 0, len(descriptions))
	for _, description := range descriptions {
		out = append(out, &oceanusv20190422.SlaveVpcDescriptions{
			VpcId:    fwHelper.StringPointerFromValue(description.VpcId),
			SubnetId: fwHelper.StringPointerFromValue(description.SubnetId),
			AppId:    fwHelper.Int64PointerFromValue(description.AppId),
			OwnerUin: fwHelper.StringPointerFromValue(description.OwnerUin),
		})
	}

	return out, diags
}
