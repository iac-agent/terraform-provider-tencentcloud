package oceanus

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	oceanus "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/oceanus/v20190422"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/fw"
	fwhelper "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/helper"
)

// Compile-time interface assertions.
var (
	_ resource.ResourceWithConfigure   = &OceanusClusterResource{}
	_ resource.ResourceWithImportState = &OceanusClusterResource{}
)

// oceanusClusterWaitTimeout bounds the eventual-consistency waits used by
// create / scale / delete. Cluster provisioning is a slow, server-side workflow,
// so it gets its own budget instead of the generic 5 minute write retry window.
const oceanusClusterWaitTimeout = 20 * time.Minute

// Cluster status values reported by DescribeClusters: 1 uninitialized,
// 3 initializing, 2 running.
const oceanusClusterStatusRunning = 2

var (
	// errOceanusClusterNotFound marks a DescribeClusters response that no longer
	// contains the cluster, so Read can drop it from state instead of failing.
	errOceanusClusterNotFound = errors.New("oceanus cluster not found")

	// errOceanusClusterPending marks a cluster that exists but has not reached
	// the expected state yet; the wait loops keep polling on it.
	errOceanusClusterPending = errors.New("oceanus cluster state not converged")
)

// NewOceanusClusterResource is the factory referenced by
// tencentcloud/framework/registry.go to register this resource.
func NewOceanusClusterResource() resource.Resource {
	return &OceanusClusterResource{}
}

// OceanusClusterResource implements resource.Resource for
// tencentcloud_oceanus_cluster.
//
// The cloud side has no cluster modification API besides CU scaling, therefore
// every argument except cu is create-only and carries RequiresReplace.
type OceanusClusterResource struct {
	fw.ResourceWithConfigure
	fw.WithImportByID
}

// OceanusClusterResourceModel maps the schema attributes; each attribute
// appears exactly once. Fields follow the schema order.
type OceanusClusterResourceModel struct {
	Id                          types.String  `tfsdk:"id"`
	ClusterName                 types.String  `tfsdk:"cluster_name"`
	RegionId                    types.Int64   `tfsdk:"region_id"`
	ZoneId                      types.Int64   `tfsdk:"zone_id"`
	LoginPassword               types.String  `tfsdk:"login_password"`
	VpcDescriptions             types.List    `tfsdk:"vpc_descriptions"`
	DefaultCosBucket            types.String  `tfsdk:"default_cos_bucket"`
	Cu                          types.Int64   `tfsdk:"cu"`
	Remark                      types.String  `tfsdk:"remark"`
	Period                      types.Int64   `tfsdk:"period"`
	InstanceChargeType          types.String  `tfsdk:"instance_charge_type"`
	ClusterType                 types.String  `tfsdk:"cluster_type"`
	RenewFlag                   types.String  `tfsdk:"renew_flag"`
	FlinkWebUiNetworkAccessType types.String  `tfsdk:"flink_web_ui_network_access_type"`
	SlaveVpcDescriptions        types.List    `tfsdk:"slave_vpc_descriptions"`
	CuMemory                    types.Int64   `tfsdk:"cu_memory"`
	Status                      types.Int64   `tfsdk:"status"`
	StatusDesc                  types.String  `tfsdk:"status_desc"`
	Region                      types.String  `tfsdk:"region"`
	AppId                       types.Int64   `tfsdk:"app_id"`
	OwnerUin                    types.String  `tfsdk:"owner_uin"`
	CreatorUin                  types.String  `tfsdk:"creator_uin"`
	CreateTime                  types.String  `tfsdk:"create_time"`
	UpdateTime                  types.String  `tfsdk:"update_time"`
	Zone                        types.String  `tfsdk:"zone"`
	NetEnvironmentType          types.Int64   `tfsdk:"net_environment_type"`
	FreeCuNum                   types.Int64   `tfsdk:"free_cu_num"`
	FreeCu                      types.Float64 `tfsdk:"free_cu"`
	RunningCu                   types.Float64 `tfsdk:"running_cu"`
	PayMode                     types.Int64   `tfsdk:"pay_mode"`
	AutoRenewFlag               types.Int64   `tfsdk:"auto_renew_flag"`
	ExpireTime                  types.String  `tfsdk:"expire_time"`
	SecondsUntilExpiry          types.String  `tfsdk:"seconds_until_expiry"`
	IsolatedTime                types.String  `tfsdk:"isolated_time"`
	ClsLogSet                   types.String  `tfsdk:"cls_log_set"`
	ClsTopicId                  types.String  `tfsdk:"cls_topic_id"`
	ClsLogName                  types.String  `tfsdk:"cls_log_name"`
	ClsTopicName                types.String  `tfsdk:"cls_topic_name"`
	Version                     types.List    `tfsdk:"version"`
	DefaultLogCollectConf       types.String  `tfsdk:"default_log_collect_conf"`
	CustomizedDnsEnabled        types.Int64   `tfsdk:"customized_dns_enabled"`
	IsNeedManageNode            types.Int64   `tfsdk:"is_need_manage_node"`
	ArchGeneration              types.Int64   `tfsdk:"arch_generation"`
	WebUiType                   types.Int64   `tfsdk:"web_ui_type"`
	Type                        types.Int64   `tfsdk:"type"`
	ResourceType                types.Int64   `tfsdk:"resource_type"`
	BillingResourceMode         types.String  `tfsdk:"billing_resource_mode"`
	MemRatio                    types.Int64   `tfsdk:"mem_ratio"`
	CrossTenantEniMode          types.Int64   `tfsdk:"cross_tenant_eni_mode"`
	TotalCpu                    types.Float64 `tfsdk:"total_cpu"`
	TotalMem                    types.Float64 `tfsdk:"total_mem"`
	RunningCpu                  types.Float64 `tfsdk:"running_cpu"`
	RunningMem                  types.Float64 `tfsdk:"running_mem"`
	DeploymentMode              types.Int64   `tfsdk:"deployment_mode"`
	LogCosBucket                types.String  `tfsdk:"log_cos_bucket"`
	CdcId                       types.String  `tfsdk:"cdc_id"`
	ClusterProcessMsg           types.String  `tfsdk:"cluster_process_msg"`
	MaxCuPerJob                 types.Int64   `tfsdk:"max_cu_per_job"`
	NetEniType                  types.Int64   `tfsdk:"net_eni_type"`
	IsolationPolicyVersion      types.Int64   `tfsdk:"isolation_policy_version"`
	SecurityGroupIds            types.Set     `tfsdk:"security_group_ids"`
	AgentSerialId               types.String  `tfsdk:"agent_serial_id"`
	Tags                        types.List    `tfsdk:"tags"`
	CcNs                        types.List    `tfsdk:"cc_ns"`
	SlaveZones                  types.List    `tfsdk:"slave_zones"`
}

// oceanusClusterVpcDescriptionModel maps one element of vpc_descriptions and
// slave_vpc_descriptions.
type oceanusClusterVpcDescriptionModel struct {
	VpcId    types.String `tfsdk:"vpc_id"`
	SubnetId types.String `tfsdk:"subnet_id"`
	AppId    types.Int64  `tfsdk:"app_id"`
	OwnerUin types.String `tfsdk:"owner_uin"`
}

// oceanusClusterJdkSupportVersionModel maps one element of
// version[*].jdk_support_version.
type oceanusClusterJdkSupportVersionModel struct {
	FlinkVersion types.String `tfsdk:"flink_version"`
	JdkVersions  types.Set    `tfsdk:"jdk_versions"`
}

// oceanusClusterVersionModel maps one element of version.
type oceanusClusterVersionModel struct {
	Flink             types.String `tfsdk:"flink"`
	SupportedFlink    types.Set    `tfsdk:"supported_flink"`
	JdkSupportVersion types.List   `tfsdk:"jdk_support_version"`
}

// oceanusClusterTagModel maps one element of tags.
type oceanusClusterTagModel struct {
	TagKey   types.String `tfsdk:"tag_key"`
	TagValue types.String `tfsdk:"tag_value"`
}

// oceanusClusterCcnModel maps one element of cc_ns.
type oceanusClusterCcnModel struct {
	VpcId    types.String `tfsdk:"vpc_id"`
	SubnetId types.String `tfsdk:"subnet_id"`
	CcnId    types.String `tfsdk:"ccn_id"`
}

// oceanusClusterSlaveZoneModel maps one element of slave_zones.
type oceanusClusterSlaveZoneModel struct {
	VpcId    types.String `tfsdk:"vpc_id"`
	SubnetId types.String `tfsdk:"subnet_id"`
	Zone     types.String `tfsdk:"zone"`
}

var (
	oceanusClusterJdkSupportVersionObjectType = types.ObjectType{AttrTypes: map[string]attr.Type{
		"flink_version": types.StringType,
		"jdk_versions":  types.SetType{ElemType: types.StringType},
	}}
	oceanusClusterVersionObjectType = types.ObjectType{AttrTypes: map[string]attr.Type{
		"flink":               types.StringType,
		"supported_flink":     types.SetType{ElemType: types.StringType},
		"jdk_support_version": types.ListType{ElemType: oceanusClusterJdkSupportVersionObjectType},
	}}
	oceanusClusterTagObjectType = types.ObjectType{AttrTypes: map[string]attr.Type{
		"tag_key":   types.StringType,
		"tag_value": types.StringType,
	}}
	oceanusClusterCcnObjectType = types.ObjectType{AttrTypes: map[string]attr.Type{
		"vpc_id":    types.StringType,
		"subnet_id": types.StringType,
		"ccn_id":    types.StringType,
	}}
	oceanusClusterSlaveZoneObjectType = types.ObjectType{AttrTypes: map[string]attr.Type{
		"vpc_id":    types.StringType,
		"subnet_id": types.StringType,
		"zone":      types.StringType,
	}}
)

func (r *OceanusClusterResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "tencentcloud_oceanus_cluster"
}

func (r *OceanusClusterResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a resource to manage an Oceanus (Stream Compute) cluster.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "ID of the Oceanus cluster, e.g. `cluster-xxxx`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"cluster_name": schema.StringAttribute{
				Required:    true,
				Description: "Cluster name. Supports 1-50 characters of English letters, Chinese characters, digits, `-` and `_`. Immutable: changing it forces resource replacement.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"region_id": schema.Int64Attribute{
				Required:    true,
				Description: "Region ID, which can be queried from the region management system with the product set to `oceanus`. Immutable: changing it forces resource replacement.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"zone_id": schema.Int64Attribute{
				Required:    true,
				Description: "Availability zone ID, which can be queried from the region management system with the product set to `oceanus`. Immutable: changing it forces resource replacement.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"login_password": schema.StringAttribute{
				Required:    true,
				Sensitive:   true,
				Description: "Password of the Flink UI, whose user name is always `admin`. Immutable: changing it forces resource replacement.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"vpc_descriptions": schema.ListNestedAttribute{
				Required:    true,
				Description: "VPC and subnet used by the cluster to reach other cloud resources. Immutable: changing it forces resource replacement.",
				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplace(),
				},
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"vpc_id": schema.StringAttribute{
							Required:    true,
							Description: "Private network (VPC) ID.",
						},
						"subnet_id": schema.StringAttribute{
							Required:    true,
							Description: "Subnet ID.",
						},
						"app_id": schema.Int64Attribute{
							Optional:    true,
							Description: "App ID of the VPC owner.",
						},
						"owner_uin": schema.StringAttribute{
							Optional:    true,
							Description: "UIN of the VPC owner.",
						},
					},
				},
			},
			"default_cos_bucket": schema.StringAttribute{
				Required:    true,
				Description: "Default COS bucket used to store job checkpoints, jar packages and delivered logs. Immutable: changing it forces resource replacement.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"cu": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "Number of CUs of the cluster. 12 CU is the minimum billing unit and one CU contains 1 CPU and 4 GB memory. The value must satisfy 12 + 7*n (n >= 0) and be greater than 12 when scaling.",
			},
			"remark": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Cluster description. Supports 1-50 characters of English letters, Chinese characters, digits, `-` and `_`. Immutable: changing it forces resource replacement.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"period": schema.Int64Attribute{
				Optional:    true,
				Description: "Purchase duration in months. Valid values: [1, 48]. Immutable: changing it forces resource replacement.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"instance_charge_type": schema.StringAttribute{
				Optional:    true,
				Description: "Cluster billing type. Valid values: `PREPAID` (monthly subscription) and `POSTPAID_BY_SECOND` (pay-as-you-go, the default). Immutable: changing it forces resource replacement.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"cluster_type": schema.StringAttribute{
				Optional:    true,
				Description: "Cluster type. Valid value: `MULTI_AZ_CLUSTER` (multi availability zone cluster). Immutable: changing it forces resource replacement.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"renew_flag": schema.StringAttribute{
				Optional:    true,
				Description: "Auto renew flag. Valid values: `NOTIFY_AND_MANUAL_RENEW` (default), `NOTIFY_AND_AUTO_RENEW` and `DISABLE_NOTIFY_AND_MANUAL_RENEW`. Required when `instance_charge_type` is `PREPAID`. Immutable: changing it forces resource replacement.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"flink_web_ui_network_access_type": schema.StringAttribute{
				Optional:    true,
				Description: "Flink UI access type. Valid values: `NetworkAccess_INTERNAL` (private network) and `NetworkAccess_EXTERNAL` (public network, the default). Immutable: changing it forces resource replacement.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"slave_vpc_descriptions": schema.ListNestedAttribute{
				Optional:    true,
				Description: "VPC and subnet of the secondary availability zone, used by multi availability zone clusters. Immutable: changing it forces resource replacement.",
				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplace(),
				},
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"vpc_id": schema.StringAttribute{
							Required:    true,
							Description: "Private network (VPC) ID.",
						},
						"subnet_id": schema.StringAttribute{
							Required:    true,
							Description: "Subnet ID.",
						},
						"app_id": schema.Int64Attribute{
							Optional:    true,
							Description: "App ID of the VPC owner.",
						},
						"owner_uin": schema.StringAttribute{
							Optional:    true,
							Description: "UIN of the VPC owner.",
						},
					},
				},
			},
			"cu_memory": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "Core memory ratio of the cluster. Only `0`, `2`, `4` and `8` are supported. Immutable: changing it forces resource replacement.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"status": schema.Int64Attribute{
				Computed:    true,
				Description: "Cluster status: 1 uninitialized, 3 initializing, 2 running.",
			},
			"status_desc": schema.StringAttribute{
				Computed:    true,
				Description: "Cluster status description.",
			},
			"region": schema.StringAttribute{
				Computed:    true,
				Description: "Region of the cluster.",
			},
			"app_id": schema.Int64Attribute{
				Computed:    true,
				Description: "App ID of the cluster owner.",
			},
			"owner_uin": schema.StringAttribute{
				Computed:    true,
				Description: "UIN of the main account owning the cluster.",
			},
			"creator_uin": schema.StringAttribute{
				Computed:    true,
				Description: "UIN of the cluster creator.",
			},
			"create_time": schema.StringAttribute{
				Computed:    true,
				Description: "Creation time of the cluster.",
			},
			"update_time": schema.StringAttribute{
				Computed:    true,
				Description: "Time of the last operation on the cluster.",
			},
			"zone": schema.StringAttribute{
				Computed:    true,
				Description: "Availability zone of the cluster.",
			},
			"net_environment_type": schema.Int64Attribute{
				Computed:    true,
				Description: "Network environment type of the cluster.",
			},
			"free_cu_num": schema.Int64Attribute{
				Computed:    true,
				Description: "Number of idle CUs.",
			},
			"free_cu": schema.Float64Attribute{
				Computed:    true,
				Description: "Idle CU under fine-grained resources.",
			},
			"running_cu": schema.Float64Attribute{
				Computed:    true,
				Description: "Number of running CUs.",
			},
			"pay_mode": schema.Int64Attribute{
				Computed:    true,
				Description: "Payment mode: 0 pay-as-you-go, 1 monthly subscription.",
			},
			"auto_renew_flag": schema.Int64Attribute{
				Computed:    true,
				Description: "Auto renew flag: 0 default, 1 auto renew, 2 explicitly no auto renew.",
			},
			"expire_time": schema.StringAttribute{
				Computed:    true,
				Description: "Expiration time of the cluster, `-` when the cluster never expires.",
			},
			"seconds_until_expiry": schema.StringAttribute{
				Computed:    true,
				Description: "Seconds until the cluster expires, `-` when the cluster never expires.",
			},
			"isolated_time": schema.StringAttribute{
				Computed:    true,
				Description: "Isolation time of the cluster, `-` when the cluster is not isolated.",
			},
			"cls_log_set": schema.StringAttribute{
				Computed:    true,
				Description: "CLS log set of the cluster.",
			},
			"cls_topic_id": schema.StringAttribute{
				Computed:    true,
				Description: "CLS log topic ID of the cluster.",
			},
			"cls_log_name": schema.StringAttribute{
				Computed:    true,
				Description: "CLS log set name of the cluster.",
			},
			"cls_topic_name": schema.StringAttribute{
				Computed:    true,
				Description: "CLS log topic name of the cluster.",
			},
			"version": schema.ListNestedAttribute{
				Computed:    true,
				Description: "Version information of the cluster.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"flink": schema.StringAttribute{
							Computed:    true,
							Description: "Flink version of the cluster.",
						},
						"supported_flink": schema.SetAttribute{
							Computed:    true,
							ElementType: types.StringType,
							Description: "Flink versions supported by the cluster.",
						},
						"jdk_support_version": schema.ListNestedAttribute{
							Computed:    true,
							Description: "Supported JDK versions per Flink version.",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"flink_version": schema.StringAttribute{
										Computed:    true,
										Description: "Flink version.",
									},
									"jdk_versions": schema.SetAttribute{
										Computed:    true,
										ElementType: types.StringType,
										Description: "JDK versions supported by this Flink version.",
									},
								},
							},
						},
					},
				},
			},
			"default_log_collect_conf": schema.StringAttribute{
				Computed:    true,
				Description: "Default log collection configuration of the cluster.",
			},
			"customized_dns_enabled": schema.Int64Attribute{
				Computed:    true,
				Description: "Whether a customized DNS is configured: 0 not configured, 1 configured, 2 not allowed.",
			},
			"is_need_manage_node": schema.Int64Attribute{
				Computed:    true,
				Description: "Whether the cluster needs the 2 CU management node logic: 1 yes, 0 no.",
			},
			"arch_generation": schema.Int64Attribute{
				Computed:    true,
				Description: "Architecture generation of the cluster; V3 clusters report 2.",
			},
			"web_ui_type": schema.Int64Attribute{
				Computed:    true,
				Description: "Flink UI access type: 0 public network, 1 private network.",
			},
			"type": schema.Int64Attribute{
				Computed:    true,
				Description: "Cluster type; 2 means a dedicated cluster.",
			},
			"resource_type": schema.Int64Attribute{
				Computed:    true,
				Description: "Resource type of the cluster.",
			},
			"billing_resource_mode": schema.StringAttribute{
				Computed:    true,
				Description: "Billing resource mode of the cluster.",
			},
			"mem_ratio": schema.Int64Attribute{
				Computed:    true,
				Description: "CU memory ratio of the cluster.",
			},
			"cross_tenant_eni_mode": schema.Int64Attribute{
				Computed:    true,
				Description: "Whether the cross tenant elastic network interface is enabled.",
			},
			"total_cpu": schema.Float64Attribute{
				Computed:    true,
				Description: "Total CPU of the cluster.",
			},
			"total_mem": schema.Float64Attribute{
				Computed:    true,
				Description: "Total memory of the cluster.",
			},
			"running_cpu": schema.Float64Attribute{
				Computed:    true,
				Description: "Running CPU of the cluster.",
			},
			"running_mem": schema.Float64Attribute{
				Computed:    true,
				Description: "Running memory of the cluster.",
			},
			"deployment_mode": schema.Int64Attribute{
				Computed:    true,
				Description: "Deployment mode: 0 single availability zone, 1 multi availability zone.",
			},
			"log_cos_bucket": schema.StringAttribute{
				Computed:    true,
				Description: "COS bucket storing the logs of the cluster.",
			},
			"cdc_id": schema.StringAttribute{
				Computed:    true,
				Description: "CDC cluster ID.",
			},
			"cluster_process_msg": schema.StringAttribute{
				Computed:    true,
				Description: "Progress message of the cluster.",
			},
			"max_cu_per_job": schema.Int64Attribute{
				Computed:    true,
				Description: "Maximum number of CUs configurable for a single job.",
			},
			"net_eni_type": schema.Int64Attribute{
				Computed:    true,
				Description: "Elastic network interface scheme: 0 POD ENI, 1 Node ENI.",
			},
			"isolation_policy_version": schema.Int64Attribute{
				Computed:    true,
				Description: "Isolation policy of the cluster: 0 seven days, 1 fifteen days.",
			},
			"security_group_ids": schema.SetAttribute{
				Computed:    true,
				ElementType: types.StringType,
				Description: "Security group IDs of the cluster.",
			},
			"agent_serial_id": schema.StringAttribute{
				Computed:    true,
				Description: "Serial ID of the parent cluster.",
			},
			"tags": schema.ListNestedAttribute{
				Computed:    true,
				Description: "Tags bound to the cluster.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"tag_key": schema.StringAttribute{
							Computed:    true,
							Description: "Tag key.",
						},
						"tag_value": schema.StringAttribute{
							Computed:    true,
							Description: "Tag value.",
						},
					},
				},
			},
			"cc_ns": schema.ListNestedAttribute{
				Computed:    true,
				Description: "Networks of the cluster.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"vpc_id": schema.StringAttribute{
							Computed:    true,
							Description: "Private network (VPC) ID.",
						},
						"subnet_id": schema.StringAttribute{
							Computed:    true,
							Description: "Subnet ID.",
						},
						"ccn_id": schema.StringAttribute{
							Computed:    true,
							Description: "Cloud connect network ID, e.g. `ccn-rahigzjd`.",
						},
					},
				},
			},
			"slave_zones": schema.ListNestedAttribute{
				Computed:    true,
				Description: "Secondary availability zones of the cluster.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"vpc_id": schema.StringAttribute{
							Computed:    true,
							Description: "Private network (VPC) ID.",
						},
						"subnet_id": schema.StringAttribute{
							Computed:    true,
							Description: "Subnet ID.",
						},
						"zone": schema.StringAttribute{
							Computed:    true,
							Description: "Availability zone.",
						},
					},
				},
			},
		},
	}
}

// Create handles `terraform apply` for new clusters.
func (r *OceanusClusterResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan OceanusClusterResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	service, ok := r.service(&resp.Diagnostics)
	if !ok {
		return
	}

	vpcDescriptions, vpcDiags := expandOceanusClusterVpcDescriptions(ctx, plan.VpcDescriptions)
	resp.Diagnostics.Append(vpcDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	slaveVpcDescriptions, slaveVpcDiags := expandOceanusClusterSlaveVpcDescriptions(ctx, plan.SlaveVpcDescriptions)
	resp.Diagnostics.Append(slaveVpcDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createOceanusClusterRequest := oceanus.NewCreateOceanusClusterRequest()
	createOceanusClusterRequest.ClusterName = fwhelper.StringPointerFromValue(plan.ClusterName)
	createOceanusClusterRequest.RegionId = fwhelper.Int64PointerFromValue(plan.RegionId)
	createOceanusClusterRequest.ZoneId = fwhelper.Int64PointerFromValue(plan.ZoneId)
	createOceanusClusterRequest.LoginPassword = fwhelper.StringPointerFromValue(plan.LoginPassword)
	createOceanusClusterRequest.VpcDescriptions = vpcDescriptions
	createOceanusClusterRequest.DefaultCOSBucket = fwhelper.StringPointerFromValue(plan.DefaultCosBucket)

	if !plan.Cu.IsNull() && !plan.Cu.IsUnknown() {
		createOceanusClusterRequest.CU = fwhelper.Int64PointerFromValue(plan.Cu)
	}

	if !plan.Remark.IsNull() && !plan.Remark.IsUnknown() {
		createOceanusClusterRequest.Remark = fwhelper.StringPointerFromValue(plan.Remark)
	}

	if !plan.Period.IsNull() && !plan.Period.IsUnknown() {
		createOceanusClusterRequest.Period = fwhelper.Int64PointerFromValue(plan.Period)
	}

	if !plan.InstanceChargeType.IsNull() && !plan.InstanceChargeType.IsUnknown() {
		createOceanusClusterRequest.InstanceChargeType = fwhelper.StringPointerFromValue(plan.InstanceChargeType)
	}

	if !plan.ClusterType.IsNull() && !plan.ClusterType.IsUnknown() {
		createOceanusClusterRequest.ClusterType = fwhelper.StringPointerFromValue(plan.ClusterType)
	}

	if !plan.RenewFlag.IsNull() && !plan.RenewFlag.IsUnknown() {
		createOceanusClusterRequest.RenewFlag = fwhelper.StringPointerFromValue(plan.RenewFlag)
	}

	if !plan.FlinkWebUiNetworkAccessType.IsNull() && !plan.FlinkWebUiNetworkAccessType.IsUnknown() {
		createOceanusClusterRequest.FlinkWebUINetworkAccessType = fwhelper.StringPointerFromValue(plan.FlinkWebUiNetworkAccessType)
	}

	if len(slaveVpcDescriptions) > 0 {
		createOceanusClusterRequest.SlaveVpcDescriptions = slaveVpcDescriptions
	}

	if !plan.CuMemory.IsNull() && !plan.CuMemory.IsUnknown() {
		createOceanusClusterRequest.CUMemory = fwhelper.Int64PointerFromValue(plan.CuMemory)
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
			"CreateOceanusCluster returned an empty cluster ID.",
		)
		return
	}

	plan.Id = types.StringValue(clusterId)
	normalizeOceanusClusterComputedUnknowns(&plan)

	// Persist the identity before waiting so a timeout still leaves a
	// destroyable resource in state.
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := waitOceanusClusterReady(ctx, service, clusterId); err != nil {
		resp.Diagnostics.AddError(
			"Error waiting for Oceanus cluster creation",
			fmt.Sprintf("The cluster %s was created but never became ready: %s", clusterId, err),
		)
		return
	}

	item, err := describeOceanusCluster(ctx, service, clusterId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading Oceanus cluster after create",
			fmt.Sprintf("CreateOceanusCluster succeeded but DescribeClusters failed: %s", err),
		)
		return
	}

	if item == nil {
		resp.Diagnostics.AddError(
			"Oceanus cluster not readable after create",
			fmt.Sprintf("Cluster %s was created but cannot be read back.", clusterId),
		)
		return
	}

	resp.Diagnostics.Append(refreshOceanusClusterModel(ctx, &plan, item)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Read handles `terraform refresh` and post-apply state refresh.
func (r *OceanusClusterResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state OceanusClusterResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	service, ok := r.service(&resp.Diagnostics)
	if !ok {
		return
	}

	item, err := describeOceanusCluster(ctx, service, state.Id.ValueString())
	if err != nil {
		if errors.Is(err, errOceanusClusterNotFound) {
			resp.Diagnostics.AddWarning(
				"Resource not found",
				fmt.Sprintf("Oceanus cluster %s not found, removing from state.", state.Id.ValueString()),
			)
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError("Error reading Oceanus cluster", err.Error())
		return
	}

	resp.Diagnostics.Append(refreshOceanusClusterModel(ctx, &state, item)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update handles in-place modifications. CU scaling is the only supported
// change, every other argument carries RequiresReplace.
func (r *OceanusClusterResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan OceanusClusterResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state OceanusClusterResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	service, ok := r.service(&resp.Diagnostics)
	if !ok {
		return
	}

	clusterId := state.Id.ValueString()

	if !plan.Cu.Equal(state.Cu) && !plan.Cu.IsNull() && !plan.Cu.IsUnknown() {
		scaleMode := "ScaleUp"
		if plan.Cu.ValueInt64() < state.Cu.ValueInt64() {
			scaleMode = "ScaleDown"
		}

		scaleOceanusClusterRequest := oceanus.NewScaleOceanusClusterRequest()
		scaleOceanusClusterRequest.ClusterId = fwhelper.StringPointerFromValue(state.Id)
		scaleOceanusClusterRequest.NewCU = fwhelper.Int64PointerFromValue(plan.Cu)
		scaleOceanusClusterRequest.ScaleMode = &scaleMode

		if err := service.ScaleOceanusCluster(ctx, scaleOceanusClusterRequest); err != nil {
			resp.Diagnostics.AddError(
				"Error scaling Oceanus cluster",
				fmt.Sprintf("ScaleOceanusCluster failed for cluster %s: %s", clusterId, err),
			)
			return
		}

		if err := waitOceanusClusterCu(ctx, service, clusterId, plan.Cu.ValueInt64()); err != nil {
			resp.Diagnostics.AddError(
				"Error waiting for Oceanus cluster scaling",
				fmt.Sprintf("The cluster %s CU was never scaled to %d: %s", clusterId, plan.Cu.ValueInt64(), err),
			)
			return
		}
	}

	item, err := describeOceanusCluster(ctx, service, clusterId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading Oceanus cluster after update",
			fmt.Sprintf("The cluster %s was updated but cannot be read back: %s", clusterId, err),
		)
		return
	}

	if item == nil {
		resp.Diagnostics.AddError(
			"Oceanus cluster not readable after update",
			fmt.Sprintf("Cluster %s was updated but cannot be read back.", clusterId),
		)
		return
	}

	plan.Id = state.Id
	resp.Diagnostics.Append(refreshOceanusClusterModel(ctx, &plan, item)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete handles resource deletion.
func (r *OceanusClusterResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state OceanusClusterResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	service, ok := r.service(&resp.Diagnostics)
	if !ok {
		return
	}

	clusterId := state.Id.ValueString()

	deleteOceanusClusterRequest := oceanus.NewDeleteOceanusClusterRequest()
	deleteOceanusClusterRequest.ClusterId = &clusterId

	if err := service.DeleteOceanusCluster(ctx, deleteOceanusClusterRequest); err != nil {
		resp.Diagnostics.AddError(
			"Error deleting Oceanus cluster",
			fmt.Sprintf("DeleteOceanusCluster failed for cluster %s: %s", clusterId, err),
		)
		return
	}

	if err := waitOceanusClusterDeleted(ctx, service, clusterId); err != nil {
		resp.Diagnostics.AddError(
			"Error waiting for Oceanus cluster deletion",
			fmt.Sprintf("The cluster %s was destroyed but is still readable: %s", clusterId, err),
		)
		return
	}
}

// service returns the Oceanus service bound to the shared client. When the
// provider has not been configured it records a diagnostic and reports false,
// so the CRUD methods never dereference a nil SDK client.
func (r *OceanusClusterResource) service(diags *diag.Diagnostics) (OceanusService, bool) {
	if r.Client() == nil {
		diags.AddError(
			"Provider not configured",
			"The provider client is not available. This is unexpected; please report it.",
		)
		return OceanusService{}, false
	}
	return NewOceanusService(r.Client()), true
}

// describeOceanusCluster reads one cluster by id. A DescribeClusters response
// that does not contain the cluster is reported as errOceanusClusterNotFound.
func describeOceanusCluster(ctx context.Context, service OceanusService, clusterId string) (*oceanus.Cluster, error) {
	if clusterId == "" {
		return nil, errOceanusClusterNotFound
	}

	describeClustersRequest := oceanus.NewDescribeClustersRequest()
	describeClustersRequest.ClusterIds = []*string{&clusterId}

	clusters, err := service.DescribeOceanusClusters(ctx, describeClustersRequest)
	if err != nil {
		return nil, err
	}

	for _, cluster := range clusters {
		if cluster == nil || cluster.ClusterId == nil {
			continue
		}
		if *cluster.ClusterId == clusterId {
			return cluster, nil
		}
	}

	return nil, errOceanusClusterNotFound
}

// waitOceanusClusterReady polls DescribeClusters until the newly created
// cluster reaches the running state.
func waitOceanusClusterReady(ctx context.Context, service OceanusService, clusterId string) error {
	_, err := fwhelper.RetryFramework(ctx, oceanusClusterWaitTimeout, func(err error) bool {
		return errors.Is(err, errOceanusClusterPending)
	}, func() (string, error) {
		item, err := describeOceanusCluster(ctx, service, clusterId)
		if err != nil {
			if errors.Is(err, errOceanusClusterNotFound) {
				return "", errOceanusClusterPending
			}
			return "", err
		}

		if item.Status == nil || *item.Status != oceanusClusterStatusRunning {
			return "", errOceanusClusterPending
		}

		return clusterId, nil
	})

	if err != nil {
		if fwhelper.IsTimeoutError(err) {
			return fmt.Errorf("timed out waiting for cluster %s to be running", clusterId)
		}
		return err
	}

	return nil
}

// waitOceanusClusterCu polls DescribeClusters until the cluster reports the
// requested CU count and is running again.
func waitOceanusClusterCu(ctx context.Context, service OceanusService, clusterId string, expectedCu int64) error {
	_, err := fwhelper.RetryFramework(ctx, oceanusClusterWaitTimeout, func(err error) bool {
		return errors.Is(err, errOceanusClusterPending)
	}, func() (string, error) {
		item, err := describeOceanusCluster(ctx, service, clusterId)
		if err != nil {
			return "", err
		}

		if item.CuNum == nil || int64(*item.CuNum) != expectedCu {
			return "", errOceanusClusterPending
		}

		if item.Status == nil || *item.Status != oceanusClusterStatusRunning {
			return "", errOceanusClusterPending
		}

		return clusterId, nil
	})

	if err != nil {
		if fwhelper.IsTimeoutError(err) {
			return fmt.Errorf("timed out waiting for cluster %s to reach %d CU", clusterId, expectedCu)
		}
		return err
	}

	return nil
}

// waitOceanusClusterDeleted polls DescribeClusters until the destroyed cluster
// is no longer reported.
func waitOceanusClusterDeleted(ctx context.Context, service OceanusService, clusterId string) error {
	_, err := fwhelper.RetryFramework(ctx, oceanusClusterWaitTimeout, func(err error) bool {
		return errors.Is(err, errOceanusClusterPending)
	}, func() (string, error) {
		item, err := describeOceanusCluster(ctx, service, clusterId)
		if err != nil {
			if errors.Is(err, errOceanusClusterNotFound) {
				return clusterId, nil
			}
			return "", err
		}

		if item == nil {
			return clusterId, nil
		}

		return "", errOceanusClusterPending
	})

	if err != nil {
		if fwhelper.IsTimeoutError(err) {
			return fmt.Errorf("timed out waiting for cluster %s to be deleted", clusterId)
		}
		return err
	}

	return nil
}

// expandOceanusClusterVpcDescriptions converts a vpc_descriptions /
// slave_vpc_descriptions list into the CreateOceanusClusterRequest shapes.
func expandOceanusClusterVpcDescriptions(ctx context.Context, list types.List) ([]*oceanus.VPCDescription, diag.Diagnostics) {
	var diags diag.Diagnostics
	if list.IsNull() || list.IsUnknown() {
		return nil, diags
	}

	var models []oceanusClusterVpcDescriptionModel
	diags.Append(list.ElementsAs(ctx, &models, false)...)
	if diags.HasError() {
		return nil, diags
	}

	out := make([]*oceanus.VPCDescription, 0, len(models))
	for _, model := range models {
		description := &oceanus.VPCDescription{
			VpcId:    fwhelper.StringPointerFromValue(model.VpcId),
			SubnetId: fwhelper.StringPointerFromValue(model.SubnetId),
		}

		if !model.AppId.IsNull() && !model.AppId.IsUnknown() {
			description.AppId = fwhelper.Int64PointerFromValue(model.AppId)
		}
		if !model.OwnerUin.IsNull() && !model.OwnerUin.IsUnknown() {
			description.OwnerUin = fwhelper.StringPointerFromValue(model.OwnerUin)
		}

		out = append(out, description)
	}

	return out, diags
}

// expandOceanusClusterSlaveVpcDescriptions converts the
// slave_vpc_descriptions list into the SlaveVpcDescriptions shape expected by
// CreateOceanusCluster.
func expandOceanusClusterSlaveVpcDescriptions(ctx context.Context, list types.List) ([]*oceanus.SlaveVpcDescriptions, diag.Diagnostics) {
	var diags diag.Diagnostics
	if list.IsNull() || list.IsUnknown() {
		return nil, diags
	}

	var models []oceanusClusterVpcDescriptionModel
	diags.Append(list.ElementsAs(ctx, &models, false)...)
	if diags.HasError() {
		return nil, diags
	}

	out := make([]*oceanus.SlaveVpcDescriptions, 0, len(models))
	for _, model := range models {
		description := &oceanus.SlaveVpcDescriptions{
			VpcId:    fwhelper.StringPointerFromValue(model.VpcId),
			SubnetId: fwhelper.StringPointerFromValue(model.SubnetId),
		}

		if !model.AppId.IsNull() && !model.AppId.IsUnknown() {
			description.AppId = fwhelper.Int64PointerFromValue(model.AppId)
		}
		if !model.OwnerUin.IsNull() && !model.OwnerUin.IsUnknown() {
			description.OwnerUin = fwhelper.StringPointerFromValue(model.OwnerUin)
		}

		out = append(out, description)
	}

	return out, diags
}

// refreshOceanusClusterModel copies cloud-owned fields into the model while
// keeping the user-supplied values for fields the cloud API does not echo back
// (region_id, zone_id, login_password, VPC descriptions, billing inputs, Flink
// UI access type). Every echoed field is only overwritten when the cloud
// returned a value, so a configured attribute never becomes null after apply.
func refreshOceanusClusterModel(ctx context.Context, m *OceanusClusterResourceModel, item *oceanus.Cluster) diag.Diagnostics {
	var diags diag.Diagnostics
	if item == nil {
		return diags
	}

	m.ClusterName = fwhelper.StringValueOrPrev(item.Name, m.ClusterName)
	m.DefaultCosBucket = fwhelper.StringValueOrPrev(item.DefaultCOSBucket, m.DefaultCosBucket)
	m.Remark = fwhelper.StringValueOrPrev(item.Remark, m.Remark)
	m.Cu = fwhelper.Int64ValueOrPrev(item.CuNum, m.Cu)
	m.CuMemory = fwhelper.Int64ValueOrPrev(item.CuMem, m.CuMemory)

	m.Status = fwhelper.Int64ValueOrNull(item.Status)
	m.StatusDesc = fwhelper.StringPointerValueOrNull(item.StatusDesc)
	m.Region = fwhelper.StringPointerValueOrNull(item.Region)
	m.AppId = fwhelper.Int64ValueFromUint(item.AppId)
	m.OwnerUin = fwhelper.StringPointerValueOrNull(item.OwnerUin)
	m.CreatorUin = fwhelper.StringPointerValueOrNull(item.CreatorUin)
	m.CreateTime = fwhelper.StringPointerValueOrNull(item.CreateTime)
	m.UpdateTime = fwhelper.StringPointerValueOrNull(item.UpdateTime)
	m.Zone = fwhelper.StringPointerValueOrNull(item.Zone)
	m.NetEnvironmentType = fwhelper.Int64ValueFromUint(item.NetEnvironmentType)
	m.FreeCuNum = fwhelper.Int64ValueOrNull(item.FreeCuNum)
	m.FreeCu = float64ValueOrNull(item.FreeCu)
	m.RunningCu = float64ValueOrNull(item.RunningCu)
	m.PayMode = fwhelper.Int64ValueOrNull(item.PayMode)
	m.AutoRenewFlag = fwhelper.Int64ValueOrNull(item.AutoRenewFlag)
	m.ExpireTime = fwhelper.StringPointerValueOrNull(item.ExpireTime)
	m.SecondsUntilExpiry = fwhelper.StringPointerValueOrNull(item.SecondsUntilExpiry)
	m.IsolatedTime = fwhelper.StringPointerValueOrNull(item.IsolatedTime)
	m.ClsLogSet = fwhelper.StringPointerValueOrNull(item.CLSLogSet)
	m.ClsTopicId = fwhelper.StringPointerValueOrNull(item.CLSTopicId)
	m.ClsLogName = fwhelper.StringPointerValueOrNull(item.CLSLogName)
	m.ClsTopicName = fwhelper.StringPointerValueOrNull(item.CLSTopicName)
	m.DefaultLogCollectConf = fwhelper.StringPointerValueOrNull(item.DefaultLogCollectConf)
	m.CustomizedDnsEnabled = fwhelper.Int64ValueOrNull(item.CustomizedDNSEnabled)
	m.IsNeedManageNode = fwhelper.Int64ValueOrNull(item.IsNeedManageNode)
	m.ArchGeneration = fwhelper.Int64ValueFromUint(item.ArchGeneration)
	m.WebUiType = fwhelper.Int64ValueOrNull(item.WebUIType)
	m.Type = fwhelper.Int64ValueOrNull(item.Type)
	m.ResourceType = fwhelper.Int64ValueOrNull(item.ResourceType)
	m.BillingResourceMode = fwhelper.StringPointerValueOrNull(item.BillingResourceMode)
	m.MemRatio = fwhelper.Int64ValueOrNull(item.MemRatio)
	m.CrossTenantEniMode = fwhelper.Int64ValueOrNull(item.CrossTenantEniMode)
	m.TotalCpu = float64ValueOrNull(item.TotalCpu)
	m.TotalMem = float64ValueOrNull(item.TotalMem)
	m.RunningCpu = float64ValueOrNull(item.RunningCpu)
	m.RunningMem = float64ValueOrNull(item.RunningMem)
	m.DeploymentMode = fwhelper.Int64ValueOrNull(item.DeploymentMode)
	m.LogCosBucket = fwhelper.StringPointerValueOrNull(item.LogCOSBucket)
	m.CdcId = fwhelper.StringPointerValueOrNull(item.CdcId)
	m.ClusterProcessMsg = fwhelper.StringPointerValueOrNull(item.ClusterProcessMsg)
	m.MaxCuPerJob = fwhelper.Int64ValueOrNull(item.MaxCuPerJob)
	m.NetEniType = fwhelper.Int64ValueOrNull(item.NetEniType)
	m.IsolationPolicyVersion = fwhelper.Int64ValueOrNull(item.IsolationPolicyVersion)
	m.AgentSerialId = fwhelper.StringPointerValueOrNull(item.AgentSerialId)

	m.SecurityGroupIds = stringSetOrNull(item.SecurityGroupIds)

	version, versionDiags := flattenOceanusClusterVersion(ctx, item.Version)
	diags.Append(versionDiags...)
	if diags.HasError() {
		return diags
	}
	m.Version = version

	tags, tagDiags := flattenOceanusClusterTags(ctx, item.Tags)
	diags.Append(tagDiags...)
	if diags.HasError() {
		return diags
	}
	m.Tags = tags

	ccNs, ccnDiags := flattenOceanusClusterCcNs(ctx, item.CCNs)
	diags.Append(ccnDiags...)
	if diags.HasError() {
		return diags
	}
	m.CcNs = ccNs

	slaveZones, slaveZoneDiags := flattenOceanusClusterSlaveZones(ctx, item.SlaveZones)
	diags.Append(slaveZoneDiags...)
	if diags.HasError() {
		return diags
	}
	m.SlaveZones = slaveZones

	return diags
}

// flattenOceanusClusterVersion converts the single cluster version object into
// the one element list the schema exposes.
func flattenOceanusClusterVersion(ctx context.Context, version *oceanus.ClusterVersion) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics
	if version == nil {
		return types.ListNull(oceanusClusterVersionObjectType), diags
	}

	jdkModels := make([]oceanusClusterJdkSupportVersionModel, 0, len(version.JdkSupportVersion))
	for _, jdk := range version.JdkSupportVersion {
		if jdk == nil {
			continue
		}

		jdkModels = append(jdkModels, oceanusClusterJdkSupportVersionModel{
			FlinkVersion: fwhelper.StringPointerValueOrNull(jdk.FlinkVersion),
			JdkVersions:  stringSetOrNull(jdk.JdkVersions),
		})
	}

	jdkList, jdkDiags := types.ListValueFrom(ctx, oceanusClusterJdkSupportVersionObjectType, jdkModels)
	diags.Append(jdkDiags...)
	if diags.HasError() {
		return types.ListNull(oceanusClusterVersionObjectType), diags
	}

	models := []oceanusClusterVersionModel{
		{
			Flink:             fwhelper.StringPointerValueOrNull(version.Flink),
			SupportedFlink:    stringSetOrNull(version.SupportedFlink),
			JdkSupportVersion: jdkList,
		},
	}

	list, listDiags := types.ListValueFrom(ctx, oceanusClusterVersionObjectType, models)
	diags.Append(listDiags...)
	if diags.HasError() {
		return types.ListNull(oceanusClusterVersionObjectType), diags
	}

	return list, diags
}

// flattenOceanusClusterTags converts the cluster tags into the schema list.
func flattenOceanusClusterTags(ctx context.Context, items []*oceanus.Tag) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics
	if len(items) == 0 {
		return types.ListNull(oceanusClusterTagObjectType), diags
	}

	models := make([]oceanusClusterTagModel, 0, len(items))
	for _, item := range items {
		if item == nil {
			continue
		}

		models = append(models, oceanusClusterTagModel{
			TagKey:   fwhelper.StringPointerValueOrNull(item.TagKey),
			TagValue: fwhelper.StringPointerValueOrNull(item.TagValue),
		})
	}

	if len(models) == 0 {
		return types.ListNull(oceanusClusterTagObjectType), diags
	}

	list, listDiags := types.ListValueFrom(ctx, oceanusClusterTagObjectType, models)
	diags.Append(listDiags...)
	if diags.HasError() {
		return types.ListNull(oceanusClusterTagObjectType), diags
	}

	return list, diags
}

// flattenOceanusClusterCcNs converts the cluster networks into the schema list.
func flattenOceanusClusterCcNs(ctx context.Context, items []*oceanus.CCN) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics
	if len(items) == 0 {
		return types.ListNull(oceanusClusterCcnObjectType), diags
	}

	models := make([]oceanusClusterCcnModel, 0, len(items))
	for _, item := range items {
		if item == nil {
			continue
		}

		models = append(models, oceanusClusterCcnModel{
			VpcId:    fwhelper.StringPointerValueOrNull(item.VpcId),
			SubnetId: fwhelper.StringPointerValueOrNull(item.SubnetId),
			CcnId:    fwhelper.StringPointerValueOrNull(item.CcnId),
		})
	}

	if len(models) == 0 {
		return types.ListNull(oceanusClusterCcnObjectType), diags
	}

	list, listDiags := types.ListValueFrom(ctx, oceanusClusterCcnObjectType, models)
	diags.Append(listDiags...)
	if diags.HasError() {
		return types.ListNull(oceanusClusterCcnObjectType), diags
	}

	return list, diags
}

// flattenOceanusClusterSlaveZones converts the secondary availability zones
// into the schema list.
func flattenOceanusClusterSlaveZones(ctx context.Context, items []*oceanus.SlaveZone) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics
	if len(items) == 0 {
		return types.ListNull(oceanusClusterSlaveZoneObjectType), diags
	}

	models := make([]oceanusClusterSlaveZoneModel, 0, len(items))
	for _, item := range items {
		if item == nil {
			continue
		}

		models = append(models, oceanusClusterSlaveZoneModel{
			VpcId:    fwhelper.StringPointerValueOrNull(item.VpcId),
			SubnetId: fwhelper.StringPointerValueOrNull(item.SubnetId),
			Zone:     fwhelper.StringPointerValueOrNull(item.Zone),
		})
	}

	if len(models) == 0 {
		return types.ListNull(oceanusClusterSlaveZoneObjectType), diags
	}

	list, listDiags := types.ListValueFrom(ctx, oceanusClusterSlaveZoneObjectType, models)
	diags.Append(listDiags...)
	if diags.HasError() {
		return types.ListNull(oceanusClusterSlaveZoneObjectType), diags
	}

	return list, diags
}

// stringSetOrNull converts a pointer slice into a string set, keeping a null
// set when the cloud returned nothing.
func stringSetOrNull(values []*string) types.Set {
	elements := make([]attr.Value, 0, len(values))
	for _, value := range values {
		if value == nil {
			continue
		}
		elements = append(elements, types.StringValue(*value))
	}

	if len(elements) == 0 {
		return types.SetNull(types.StringType)
	}

	return types.SetValueMust(types.StringType, elements)
}

// float64ValueOrNull converts a *float64 to types.Float64.
func float64ValueOrNull(value *float64) types.Float64 {
	if value == nil {
		return types.Float64Null()
	}
	return types.Float64Value(*value)
}

// normalizeOceanusClusterComputedUnknowns rewrites unknown computed attributes
// to null so the model can be persisted into state before the read-back fills
// it with the cloud-side values.
func normalizeOceanusClusterComputedUnknowns(m *OceanusClusterResourceModel) {
	if m.Cu.IsUnknown() {
		m.Cu = types.Int64Null()
	}
	if m.Remark.IsUnknown() {
		m.Remark = types.StringNull()
	}
	if m.CuMemory.IsUnknown() {
		m.CuMemory = types.Int64Null()
	}
	if m.Status.IsUnknown() {
		m.Status = types.Int64Null()
	}
	if m.StatusDesc.IsUnknown() {
		m.StatusDesc = types.StringNull()
	}
	if m.Region.IsUnknown() {
		m.Region = types.StringNull()
	}
	if m.AppId.IsUnknown() {
		m.AppId = types.Int64Null()
	}
	if m.OwnerUin.IsUnknown() {
		m.OwnerUin = types.StringNull()
	}
	if m.CreatorUin.IsUnknown() {
		m.CreatorUin = types.StringNull()
	}
	if m.CreateTime.IsUnknown() {
		m.CreateTime = types.StringNull()
	}
	if m.UpdateTime.IsUnknown() {
		m.UpdateTime = types.StringNull()
	}
	if m.Zone.IsUnknown() {
		m.Zone = types.StringNull()
	}
	if m.NetEnvironmentType.IsUnknown() {
		m.NetEnvironmentType = types.Int64Null()
	}
	if m.FreeCuNum.IsUnknown() {
		m.FreeCuNum = types.Int64Null()
	}
	if m.FreeCu.IsUnknown() {
		m.FreeCu = types.Float64Null()
	}
	if m.RunningCu.IsUnknown() {
		m.RunningCu = types.Float64Null()
	}
	if m.PayMode.IsUnknown() {
		m.PayMode = types.Int64Null()
	}
	if m.AutoRenewFlag.IsUnknown() {
		m.AutoRenewFlag = types.Int64Null()
	}
	if m.ExpireTime.IsUnknown() {
		m.ExpireTime = types.StringNull()
	}
	if m.SecondsUntilExpiry.IsUnknown() {
		m.SecondsUntilExpiry = types.StringNull()
	}
	if m.IsolatedTime.IsUnknown() {
		m.IsolatedTime = types.StringNull()
	}
	if m.ClsLogSet.IsUnknown() {
		m.ClsLogSet = types.StringNull()
	}
	if m.ClsTopicId.IsUnknown() {
		m.ClsTopicId = types.StringNull()
	}
	if m.ClsLogName.IsUnknown() {
		m.ClsLogName = types.StringNull()
	}
	if m.ClsTopicName.IsUnknown() {
		m.ClsTopicName = types.StringNull()
	}
	if m.Version.IsUnknown() {
		m.Version = types.ListNull(oceanusClusterVersionObjectType)
	}
	if m.DefaultLogCollectConf.IsUnknown() {
		m.DefaultLogCollectConf = types.StringNull()
	}
	if m.CustomizedDnsEnabled.IsUnknown() {
		m.CustomizedDnsEnabled = types.Int64Null()
	}
	if m.IsNeedManageNode.IsUnknown() {
		m.IsNeedManageNode = types.Int64Null()
	}
	if m.ArchGeneration.IsUnknown() {
		m.ArchGeneration = types.Int64Null()
	}
	if m.WebUiType.IsUnknown() {
		m.WebUiType = types.Int64Null()
	}
	if m.Type.IsUnknown() {
		m.Type = types.Int64Null()
	}
	if m.ResourceType.IsUnknown() {
		m.ResourceType = types.Int64Null()
	}
	if m.BillingResourceMode.IsUnknown() {
		m.BillingResourceMode = types.StringNull()
	}
	if m.MemRatio.IsUnknown() {
		m.MemRatio = types.Int64Null()
	}
	if m.CrossTenantEniMode.IsUnknown() {
		m.CrossTenantEniMode = types.Int64Null()
	}
	if m.TotalCpu.IsUnknown() {
		m.TotalCpu = types.Float64Null()
	}
	if m.TotalMem.IsUnknown() {
		m.TotalMem = types.Float64Null()
	}
	if m.RunningCpu.IsUnknown() {
		m.RunningCpu = types.Float64Null()
	}
	if m.RunningMem.IsUnknown() {
		m.RunningMem = types.Float64Null()
	}
	if m.DeploymentMode.IsUnknown() {
		m.DeploymentMode = types.Int64Null()
	}
	if m.LogCosBucket.IsUnknown() {
		m.LogCosBucket = types.StringNull()
	}
	if m.CdcId.IsUnknown() {
		m.CdcId = types.StringNull()
	}
	if m.ClusterProcessMsg.IsUnknown() {
		m.ClusterProcessMsg = types.StringNull()
	}
	if m.MaxCuPerJob.IsUnknown() {
		m.MaxCuPerJob = types.Int64Null()
	}
	if m.NetEniType.IsUnknown() {
		m.NetEniType = types.Int64Null()
	}
	if m.IsolationPolicyVersion.IsUnknown() {
		m.IsolationPolicyVersion = types.Int64Null()
	}
	if m.SecurityGroupIds.IsUnknown() {
		m.SecurityGroupIds = types.SetNull(types.StringType)
	}
	if m.AgentSerialId.IsUnknown() {
		m.AgentSerialId = types.StringNull()
	}
	if m.Tags.IsUnknown() {
		m.Tags = types.ListNull(oceanusClusterTagObjectType)
	}
	if m.CcNs.IsUnknown() {
		m.CcNs = types.ListNull(oceanusClusterCcnObjectType)
	}
	if m.SlaveZones.IsUnknown() {
		m.SlaveZones = types.ListNull(oceanusClusterSlaveZoneObjectType)
	}
}
