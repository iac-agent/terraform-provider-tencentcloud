// Oceanus cluster — terraform-plugin-framework reference resource
// (codegen subtype "general").
//
// API mapping (resource_type general, resource_id = ClusterId):
//   - create: CreateOceanusCluster (async, polled via DescribeClusters)
//   - read:   DescribeClusters (list-query, the resource picks its own item)
//   - update: ScaleOceanusCluster (CU resizing; synchronous API whose cloud
//     effect is observed through DescribeClusters)
//   - delete: DeleteOceanusCluster (async, observed gone via DescribeClusters)

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

	oceanus "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/oceanus/v20190422"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/fw"
)

// Compile-time interface assertions. The base resource.Resource interface is
// checked implicitly by the factory return value.
var (
	_ resource.ResourceWithConfigure   = &OceanusClusterAfResource{}
	_ resource.ResourceWithImportState = &OceanusClusterAfResource{}
)

// NewOceanusClusterAfResource is the factory referenced by
// tencentcloud/framework/registry.go to register this resource.
func NewOceanusClusterAfResource() resource.Resource {
	return &OceanusClusterAfResource{}
}

// NewOceanusClusterService returns the shared Oceanus service bound to the
// configured client (framework reference implementation).
func NewOceanusClusterService(client *connectivity.TencentCloudClient) OceanusService {
	return OceanusService{client: client}
}

// OceanusClusterAfResource implements resource.Resource for
// tencentcloud_oceanus_cluster_af.
type OceanusClusterAfResource struct {
	fw.ResourceWithConfigure
	// The cluster is addressed by its single ClusterId, so import passes the
	// import id straight into the state "id".
	fw.WithImportByID
}

// OceanusClusterAfResourceModel maps the schema attributes; each attribute
// appears exactly once. Field order mirrors the schema.
type OceanusClusterAfResourceModel struct {
	Id       types.String           `tfsdk:"id"`
	Timeouts resourcetimeouts.Value `tfsdk:"timeouts"`

	// Create-time configuration (CreateOceanusCluster request; no update
	// coverage => RequiresReplace).
	ClusterName                 types.String `tfsdk:"cluster_name"`
	RegionId                    types.Int64  `tfsdk:"region_id"`
	ZoneId                      types.Int64  `tfsdk:"zone_id"`
	LoginPassword               types.String `tfsdk:"login_password"`
	VpcDescriptions             types.List   `tfsdk:"vpc_descriptions"`
	DefaultCOSBucket            types.String `tfsdk:"default_cos_bucket"`
	Remark                      types.String `tfsdk:"remark"`
	Period                      types.Int64  `tfsdk:"period"`
	InstanceChargeType          types.String `tfsdk:"instance_charge_type"`
	ClusterType                 types.String `tfsdk:"cluster_type"`
	RenewFlag                   types.String `tfsdk:"renew_flag"`
	FlinkWebUINetworkAccessType types.String `tfsdk:"flink_web_ui_network_access_type"`
	SlaveVpcDescriptions        types.List   `tfsdk:"slave_vpc_descriptions"`
	CUMemory                    types.Int64  `tfsdk:"cu_memory"`

	// Scale coverage (ScaleOceanusCluster request).
	CU types.Int64 `tfsdk:"cu"`

	// Read-back values (DescribeClusters response).
	Status        types.String `tfsdk:"status"`
	StatusDesc    types.String `tfsdk:"status_desc"`
	CreateTime    types.String `tfsdk:"create_time"`
	UpdateTime    types.String `tfsdk:"update_time"`
	Zone          types.String `tfsdk:"zone"`
	ClusterReadId types.String `tfsdk:"cluster_read_id"`
}

// oceanusClusterAfVpcDescriptionModel maps one element of vpc_descriptions.
type oceanusClusterAfVpcDescriptionModel struct {
	VpcId    types.String `tfsdk:"vpc_id"`
	SubnetId types.String `tfsdk:"subnet_id"`
}

// oceanusClusterAfSlaveVpcDescriptionsModel maps one element of
// slave_vpc_descriptions.
type oceanusClusterAfSlaveVpcDescriptionsModel struct {
	VpcId    types.String `tfsdk:"vpc_id"`
	SubnetId types.String `tfsdk:"subnet_id"`
}

func (r *OceanusClusterAfResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "tencentcloud_oceanus_cluster_af"
}

func (r *OceanusClusterAfResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a resource to create and manage an Oceanus cluster (terraform-plugin-framework implementation). " +
			"Creation and deletion are asynchronous: the provider polls DescribeClusters until the cluster converges. " +
			"Only the cluster CU count (cu) supports in-place update (ScaleOceanusCluster); changing any other " +
			"configuration forces the cluster to be recreated.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "ID of the resource, in the format of cluster-xxx.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"cluster_name": schema.StringAttribute{
				Required:    true,
				Description: "Cluster name. It supports 1-50 English letters, Chinese characters, digits, hyphens `-` or underscores `_`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"region_id": schema.Int64Attribute{
				Required:    true,
				Description: "Region ID, which can be queried through DescribeRegions of the region management system, with the Product parameter set to oceanus.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"zone_id": schema.Int64Attribute{
				Required:    true,
				Description: "Availability zone ID, which can be queried through DescribeZones of the region management system, with the Product parameter set to oceanus.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"login_password": schema.StringAttribute{
				Required:    true,
				Sensitive:   true,
				Description: "Password for accessing the Flink UI. The user name and password are used to log in to the Flink UI for viewing jobs. The cluster user name is admin.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"vpc_descriptions": schema.ListNestedAttribute{
				Required:    true,
				Description: "Stream Compute accesses other cloud resources in the same region through VPCs and elastic network interfaces, which occupy a certain number of subnet IPs. Make sure the selected subnet has enough available IPs. If the existing network does not meet your requirements, create a new VPC or subnet in the VPC console.",
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
					},
				},
				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplace(),
				},
			},
			"default_cos_bucket": schema.StringAttribute{
				Required:    true,
				Description: "Stream Compute uses the COS object storage to save job checkpoints, jar packages, delivered logs, and so on. If no bucket is available in the local region, create one in the COS console. To ensure normal use, refer to the documentation for the lifecycle configuration of the corresponding COS bucket.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"cu": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "Number of cluster CUs. 12 CUs are the minimum computing resource and billing unit of Stream Compute; 1 CU contains 1 CPU and 4 GB of memory. Management node fees are waived when the number of CUs is greater than or equal to 48. The target CU must be greater than 12 and satisfy 12 + 7*n (n >= 0).",
			},
			"remark": schema.StringAttribute{
				Optional:    true,
				Description: "Cluster description. It supports 1-50 English letters, Chinese characters, digits, hyphens `-` or underscores `_`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"period": schema.Int64Attribute{
				Optional:    true,
				Description: "Purchase duration in months. Valid range: [1, 48].",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"instance_charge_type": schema.StringAttribute{
				Optional:    true,
				Description: "Cluster billing type. Valid values: `PREPAID` (monthly subscription), `POSTPAID_BY_SECOND` (pay as you go). Default value: `POSTPAID_BY_SECOND`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"cluster_type": schema.StringAttribute{
				Optional:    true,
				Description: "Cluster type. Valid value: `MULTI_AZ_CLUSTER` (multi-availability-zone cluster).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"renew_flag": schema.StringAttribute{
				Optional:    true,
				Description: "Auto-renewal flag. Valid values: `NOTIFY_AND_MANUAL_RENEW` (notify and renew manually), `NOTIFY_AND_AUTO_RENEW` (notify and renew automatically), `DISABLE_NOTIFY_AND_MANUAL_RENEW` (no notification and no manual renewal). Default value: `NOTIFY_AND_MANUAL_RENEW`. When InstanceChargeType is PREPAID, corresponding to a monthly-subscription cluster, the auto-renewal flag must be set; pay-as-you-go clusters do not need it.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"flink_web_ui_network_access_type": schema.StringAttribute{
				Optional:    true,
				Description: "Flink UI access type. Valid values: `NetworkAccess_INTERNAL` (intranet access), `NetworkAccess_EXTERNAL` (public network access). Default value: `NetworkAccess_EXTERNAL`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"slave_vpc_descriptions": schema.ListNestedAttribute{
				Optional:    true,
				Description: "VPCs of the multi-availability-zone cluster.",
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
					},
				},
				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplace(),
				},
			},
			"cu_memory": schema.Int64Attribute{
				Optional:    true,
				Description: "CU-to-memory ratio. Only [0, 2, 4, 8] are supported.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			// Computed values read back from DescribeClusters.
			"status": schema.StringAttribute{
				Computed:    true,
				Description: "Cluster status: 1 not initialized, 3 initializing, 2 running.",
			},
			"status_desc": schema.StringAttribute{
				Computed:    true,
				Description: "Cluster status description.",
			},
			"create_time": schema.StringAttribute{
				Computed:    true,
				Description: "Cluster creation time.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"update_time": schema.StringAttribute{
				Computed:    true,
				Description: "Time of the last operation on the cluster. Default value: `-`.",
			},
			"zone": schema.StringAttribute{
				Computed:    true,
				Description: "Availability zone.",
			},
			"cluster_read_id": schema.StringAttribute{
				Computed:    true,
				Description: "Cluster ID read back from DescribeClusters; the same value as `id`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
		Blocks: map[string]schema.Block{
			"timeouts": oceanusClusterAfTimeoutsBlock(ctx),
		},
	}
}

// defaultOceanusClusterAfAsyncTimeout is the built-in maximum duration to wait
// for the asynchronous create/delete of an Oceanus cluster when the user does
// not override it through the `timeouts` block, e.g. `timeouts { create = "30m" }`.
const defaultOceanusClusterAfAsyncTimeout = 15 * time.Minute

// oceanusClusterAfTimeoutsBlock builds the `timeouts` block from the official
// terraform-plugin-framework-timeouts module. That module leaves the block
// itself undescribed (only the attributes carry descriptions), so the block
// description is restored here to keep the schema self-explanatory.
func oceanusClusterAfTimeoutsBlock(ctx context.Context) schema.Block {
	block := resourcetimeouts.Block(ctx, resourcetimeouts.Opts{
		Create:            true,
		Update:            true,
		Delete:            true,
		CreateDescription: "A string that can be parsed as a duration (https://pkg.go.dev/time#ParseDuration), e.g. `30m`. It limits how long to wait for the asynchronous cluster creation. Default is `15m`.",
		UpdateDescription: "A string that can be parsed as a duration (https://pkg.go.dev/time#ParseDuration), e.g. `30m`. It limits how long to wait for the asynchronous cluster resizing. Default is `15m`.",
		DeleteDescription: "A string that can be parsed as a duration (https://pkg.go.dev/time#ParseDuration), e.g. `30m`. It limits how long to wait for the asynchronous cluster deletion. Default is `15m`.",
	})

	if nested, ok := block.(schema.SingleNestedBlock); ok {
		nested.Description = "The timeouts block allows you to specify the timeout for the create/update/delete operations."
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

	createOceanusClusterRequest := oceanus.NewCreateOceanusClusterRequest()
	createOceanusClusterRequest.ClusterName = plan.ClusterName.ValueStringPointer()
	createOceanusClusterRequest.RegionId = plan.RegionId.ValueInt64Pointer()
	createOceanusClusterRequest.ZoneId = plan.ZoneId.ValueInt64Pointer()
	createOceanusClusterRequest.LoginPassword = plan.LoginPassword.ValueStringPointer()

	vpcs, vpcDiags := expandOceanusClusterAfVpcDescriptions(ctx, plan.VpcDescriptions)
	resp.Diagnostics.Append(vpcDiags...)
	if resp.Diagnostics.HasError() {
		return
	}
	createOceanusClusterRequest.VpcDescriptions = vpcs

	createOceanusClusterRequest.DefaultCOSBucket = plan.DefaultCOSBucket.ValueStringPointer()
	if !plan.CU.IsNull() && !plan.CU.IsUnknown() {
		createOceanusClusterRequest.CU = plan.CU.ValueInt64Pointer()
	}
	createOceanusClusterRequest.Remark = plan.Remark.ValueStringPointer()
	createOceanusClusterRequest.Period = plan.Period.ValueInt64Pointer()
	createOceanusClusterRequest.InstanceChargeType = plan.InstanceChargeType.ValueStringPointer()
	createOceanusClusterRequest.ClusterType = plan.ClusterType.ValueStringPointer()
	createOceanusClusterRequest.RenewFlag = plan.RenewFlag.ValueStringPointer()
	createOceanusClusterRequest.FlinkWebUINetworkAccessType = plan.FlinkWebUINetworkAccessType.ValueStringPointer()

	slaves, slaveDiags := expandOceanusClusterAfSlaveVpcDescriptions(ctx, plan.SlaveVpcDescriptions)
	resp.Diagnostics.Append(slaveDiags...)
	if resp.Diagnostics.HasError() {
		return
	}
	createOceanusClusterRequest.SlaveVpcDescriptions = slaves

	createOceanusClusterRequest.CUMemory = plan.CUMemory.ValueInt64Pointer()

	service := NewOceanusClusterService(r.Client())

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
	// below.
	normalizeOceanusClusterAfComputedUnknowns(&plan)

	// Persist the identity before the wait so a partial failure still leaves a
	// destroyable resource in state.
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// CreateOceanusCluster is asynchronous: wait until the cluster is visible
	// (and past its initializing phase) before reading back.
	if err := service.WaitOceanusClusterCreated(ctx, clusterId, timeout); err != nil {
		resp.Diagnostics.AddError(
			"Error waiting for Oceanus cluster creation",
			fmt.Sprintf("cluster creation wait failed: %v", err),
		)

		return
	}

	// Read-back — same request form as Read; the caller picks its own item
	// from the list.
	describeClustersRequest := oceanus.NewDescribeClustersRequest()
	describeClustersRequest.ClusterIds = []*string{&clusterId}
	clusters, err := service.DescribeClustersFramework(ctx, describeClustersRequest)
	if err != nil {
		resp.Diagnostics.AddError("Error reading Oceanus cluster after create", err.Error())
		return
	}

	item := findOceanusClusterAf(clusters, clusterId)
	if item == nil {
		resp.Diagnostics.AddError(
			"Oceanus cluster not readable after create",
			fmt.Sprintf("cluster %s was created but cannot be read back", plan.Id.ValueString()),
		)

		return
	}

	resp.Diagnostics.Append(refreshOceanusClusterAfModel(ctx, &plan, item)...)
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

	describeClustersRequest := oceanus.NewDescribeClustersRequest()
	describeClustersRequest.ClusterIds = []*string{&clusterId}

	service := NewOceanusClusterService(r.Client())

	clusters, err := service.DescribeClustersFramework(ctx, describeClustersRequest)
	if err != nil {
		resp.Diagnostics.AddError("Error reading Oceanus cluster", err.Error())
		return
	}

	item := findOceanusClusterAf(clusters, clusterId)

	// Confirmed "not found" on the cloud side => drop it from state so
	// Terraform plans a recreation. The not-found judgement stays INLINE in
	// Read (RemoveResource followed by an immediate return).
	if item == nil {
		resp.Diagnostics.AddWarning(
			"Resource not found",
			fmt.Sprintf("Oceanus cluster %s not found, removing from state.", state.Id.ValueString()),
		)

		resp.State.RemoveResource(ctx)
		return
	}

	resp.Diagnostics.Append(refreshOceanusClusterAfModel(ctx, &state, item)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update handles in-place modifications. ScaleOceanusCluster only covers the
// cluster CU count, so it is called only when `cu` changed; the read-back runs
// either way.
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

	service := NewOceanusClusterService(r.Client())

	// Call the scale API ONLY when the attribute it covers changed. A plan
	// that differs from the state in the `timeouts` block alone also reaches
	// Update; sending the write API then is a pointless no-change request.
	// TODO(human): ScaleOceanusCluster carries no explicit target-state
	// polling guidance; the resize is treated as an asynchronous cloud-side
	// operation observed through DescribeClusters. The wait polls only that
	// the read-back answer differs from the pre-scale CU value (or is absent,
	// which keeps the planned value). Confirm whether a dedicated terminal
	// status exists and tighten the wait accordingly.
	if !plan.CU.Equal(state.CU) {
		// The timeout is resolved BEFORE the write API, so a malformed
		// duration fails the apply without any cloud-side effect.
		timeout, timeoutDiags := plan.Timeouts.Update(ctx, defaultOceanusClusterAfAsyncTimeout)
		resp.Diagnostics.Append(timeoutDiags...)
		if resp.Diagnostics.HasError() {
			return
		}

		scaleOceanusClusterRequest := oceanus.NewScaleOceanusClusterRequest()
		scaleOceanusClusterRequest.ClusterId = &clusterId
		scaleOceanusClusterRequest.NewCU = plan.CU.ValueInt64Pointer()

		if err := service.ScaleOceanusClusterFramework(ctx, scaleOceanusClusterRequest); err != nil {
			resp.Diagnostics.AddError(
				"Error updating Oceanus cluster",
				fmt.Sprintf("ScaleOceanusCluster failed: %v", err),
			)

			return
		}

		previousCU := int64(0)
		if !state.CU.IsNull() && !state.CU.IsUnknown() {
			previousCU = state.CU.ValueInt64()
		}

		if err := service.WaitOceanusClusterScaled(ctx, clusterId, previousCU, timeout); err != nil {
			resp.Diagnostics.AddError(
				"Error waiting for Oceanus cluster resizing",
				fmt.Sprintf("cluster resizing wait failed: %v", err),
			)

			return
		}
	}

	describeClustersRequest := oceanus.NewDescribeClustersRequest()
	describeClustersRequest.ClusterIds = []*string{&clusterId}
	clusters, err := service.DescribeClustersFramework(ctx, describeClustersRequest)
	if err != nil {
		resp.Diagnostics.AddError("Error reading Oceanus cluster after update", err.Error())
		return
	}

	item := findOceanusClusterAf(clusters, clusterId)
	if item == nil {
		resp.Diagnostics.AddError(
			"Oceanus cluster not readable after update",
			fmt.Sprintf("cluster %s was updated but cannot be read back", state.Id.ValueString()),
		)

		return
	}

	plan.Id = state.Id
	resp.Diagnostics.Append(refreshOceanusClusterAfModel(ctx, &plan, item)...)
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

	deleteOceanusClusterRequest := oceanus.NewDeleteOceanusClusterRequest()
	deleteOceanusClusterRequest.ClusterId = &clusterId

	service := NewOceanusClusterService(r.Client())

	// The timeout is resolved BEFORE the write API (see Update).
	timeout, timeoutDiags := state.Timeouts.Delete(ctx, defaultOceanusClusterAfAsyncTimeout)
	resp.Diagnostics.Append(timeoutDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := service.DeleteOceanusClusterFramework(ctx, deleteOceanusClusterRequest); err != nil {
		resp.Diagnostics.AddError(
			"Error deleting Oceanus cluster",
			fmt.Sprintf("DeleteOceanusCluster failed: %v", err),
		)

		return
	}

	// Deletion is asynchronous: wait until the describe answer no longer
	// contains the cluster (the terminal "destroyed" state).
	if err := service.WaitOceanusClusterDestroyed(ctx, clusterId, timeout); err != nil {
		resp.Diagnostics.AddError(
			"Error waiting for Oceanus cluster deletion",
			fmt.Sprintf("cluster deletion wait failed: %v", err),
		)

		return
	}
}

// refreshOceanusClusterAfModel copies cloud-owned fields into the model.
//
// User-configurable attributes keep their previous value when the cloud omits
// them (otherwise Terraform fails the post-apply consistency check with
// "was <value>, but now null"); computed-only attributes record the truthful
// refreshed value, including null.
func refreshOceanusClusterAfModel(ctx context.Context, m *OceanusClusterAfResourceModel, item *oceanus.Cluster) diag.Diagnostics {
	if item == nil {
		return nil
	}

	var diags diag.Diagnostics

	m.ClusterName = types.StringPointerValue(item.Name)
	m.DefaultCOSBucket = types.StringPointerValue(item.DefaultCOSBucket)
	m.Remark = types.StringPointerValue(item.Remark)
	m.Zone = types.StringPointerValue(item.Zone)
	m.Status = oceanusClusterAfInt64String(item.Status)
	m.StatusDesc = types.StringPointerValue(item.StatusDesc)
	m.CreateTime = types.StringPointerValue(item.CreateTime)
	m.UpdateTime = types.StringPointerValue(item.UpdateTime)
	if item.ClusterId != nil {
		m.ClusterReadId = types.StringValue(*item.ClusterId)
	}

	// cu is user-configurable (scale coverage): a cloud-omitted field keeps
	// the previous value.
	if item.CuNum != nil {
		m.CU = types.Int64Value(int64(*item.CuNum))
	}

	// login_password is write-only: the Describe API never echoes it back, so
	// the planned/state value is kept untouched.

	// vpc_descriptions and slave_vpc_descriptions are Required/Optional create
	// inputs without update coverage. The DescribeClusters CCNs response shape
	// (VPC/subnet/CCN IDs) does not round-trip the create inputs, so the
	// planned values are kept as-is.
	// TODO(human): verify whether DescribeClusters can return the exact
	// configured VPC/subnet pairs of a cluster and refresh these two lists
	// from the cloud response instead of keeping the plan values.

	return diags
}

// normalizeOceanusClusterAfComputedUnknowns rewrites unknown computed
// attributes to null so the model can be persisted into state before the
// read-back fills them with the cloud-side values.
func normalizeOceanusClusterAfComputedUnknowns(m *OceanusClusterAfResourceModel) {
	m.Status = types.StringNull()
	m.StatusDesc = types.StringNull()
	m.CreateTime = types.StringNull()
	m.UpdateTime = types.StringNull()
	m.Zone = types.StringNull()
	m.ClusterReadId = types.StringNull()
	if m.CU.IsUnknown() {
		m.CU = types.Int64Null()
	}
}

// findOceanusClusterAf picks this resource's own item from a describe list —
// the LIST-QUERY shape (DescribeClusters returns the cluster list).
func findOceanusClusterAf(clusters []*oceanus.Cluster, clusterId string) *oceanus.Cluster {
	for _, c := range clusters {
		if c.ClusterId != nil && *c.ClusterId == clusterId {
			return c
		}
	}

	return nil
}

// oceanusClusterAfInt64String renders an optional *int64 as a types.String
// (nil -> null).
func oceanusClusterAfInt64String(v *int64) types.String {
	if v == nil {
		return types.StringNull()
	}
	return types.StringValue(fmt.Sprintf("%d", *v))
}

// expandOceanusClusterAfVpcDescriptions converts the planned vpc_descriptions
// list into the SDK VPCDescription slice.
func expandOceanusClusterAfVpcDescriptions(ctx context.Context, list types.List) ([]*oceanus.VPCDescription, diag.Diagnostics) {
	var diags diag.Diagnostics
	if list.IsNull() || list.IsUnknown() {
		return nil, diags
	}

	var models []oceanusClusterAfVpcDescriptionModel
	diags.Append(list.ElementsAs(ctx, &models, false)...)
	if diags.HasError() {
		return nil, diags
	}

	out := make([]*oceanus.VPCDescription, 0, len(models))
	for _, m := range models {
		out = append(out, &oceanus.VPCDescription{
			VpcId:    m.VpcId.ValueStringPointer(),
			SubnetId: m.SubnetId.ValueStringPointer(),
		})
	}
	return out, diags
}

// expandOceanusClusterAfSlaveVpcDescriptions converts the planned
// slave_vpc_descriptions list into the SDK SlaveVpcDescriptions slice.
func expandOceanusClusterAfSlaveVpcDescriptions(ctx context.Context, list types.List) ([]*oceanus.SlaveVpcDescriptions, diag.Diagnostics) {
	var diags diag.Diagnostics
	if list.IsNull() || list.IsUnknown() {
		return nil, diags
	}

	var models []oceanusClusterAfSlaveVpcDescriptionsModel
	diags.Append(list.ElementsAs(ctx, &models, false)...)
	if diags.HasError() {
		return nil, diags
	}

	out := make([]*oceanus.SlaveVpcDescriptions, 0, len(models))
	for _, m := range models {
		out = append(out, &oceanus.SlaveVpcDescriptions{
			VpcId:    m.VpcId.ValueStringPointer(),
			SubnetId: m.SubnetId.ValueStringPointer(),
		})
	}
	return out, diags
}
