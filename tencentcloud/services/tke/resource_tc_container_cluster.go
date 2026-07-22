package tke

import (
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svccvm "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/cvm"

	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudContainerCluster() *schema.Resource {
	return &schema.Resource{
		DeprecationMessage: "This resource has been deprecated in Terraform TencentCloud provider version 1.16.0. Please use 'tencentcloud_kubernetes_cluster' instead.",
		Create:             resourceTencentCloudContainerClusterCreate,
		Read:               resourceTencentCloudContainerClusterRead,
		Update:             resourceTencentCloudContainerClusterUpdate,
		Delete:             resourceTencentCloudContainerClusterDelete,

		Schema: map[string]*schema.Schema{
			"cluster_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The 名称 cluster。",
			},
			"cpu": {
				Type:        schema.TypeInt,
				Optional:    true,
				Deprecated:  "It has been deprecated from version 1.16.0. Set 'instance_type' instead.",
				Description: "The cpu of the node。",
			},
			"mem": {
				Type:        schema.TypeInt,
				Optional:    true,
				Deprecated:  "It has been deprecated from version 1.16.0. Set 'instance_type' instead.",
				Description: "The memory of the node。",
			},
			"os_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The system os 名称 node。",
			},
			"bandwidth": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The network bandwidth of the node。",
			},
			"bandwidth_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The network 类型 node。",
			},
			"require_wan_ip": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Indicate whether wan ip is needed。",
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The 子网 ID which the node stays in。",
			},
			"is_vpc_gateway": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Describe 是否node enable the gateway capability。",
			},
			"storage_size": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The size of the data volume。",
			},
			"storage_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The 类型 data volume. see more from CVM。",
			},
			"root_size": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The size of the root volume。",
			},
			"root_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The 类型 root volume. see more from CVM。",
			},
			"goods_num": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The node number is going to create in the cluster。",
			},
			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The 密码 of each node。",
			},
			"key_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The key_id of each node(if using 键 pair to access)。",
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "指定vpc which the node(s) stay in。",
			},
			"cluster_cidr": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The CIDR which the cluster is going to use。",
			},
			"cluster_desc": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The 描述 cluster。",
			},
			"cvm_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "类型 node needed by cvm。",
			},
			"period": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The puchase duration of the node needed by cvm。",
			},
			"zone_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The 可用区 which the node stays in。",
			},
			"instance_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The instance 类型 node needed by cvm。",
			},
			"sg_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The 安全组 ID",
			},
			"mount_target": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The 路径 which volume is going to be mounted。",
			},
			"docker_graph_path": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The docker graph 路径 is going to mounted。",
			},
			"user_script": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "用户 defined script in a base64-格式 The script runs after the kubernetes component is ready on node. see more from CCS api documents。",
			},
			"unschedulable": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Determine 是否node will be schedulable. 0 is the default meaning node will be schedulable. 1 for unschedulable。",
			},
			"instance_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The 名称 ot node。",
			},
			"cluster_version": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The kubernetes 版本 of the cluster。",
			},
			"kubernetes_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The kubernetes 版本 of the cluster。",
			},
			"nodes_num": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The node 数量 the cluster。",
			},
			"nodes_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The node 状态 cluster。",
			},
			"total_cpu": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The total cpu of the cluster。",
			},
			"total_mem": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The total memory of the cluster。",
			},
		},
	}
}

func resourceTencentCloudContainerClusterUpdate(d *schema.ResourceData, m interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_container_cluster.update")()

	return fmt.Errorf("the container cluster resource doesn't support update")
}

func resourceTencentCloudContainerClusterCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_container_cluster.create")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	service := TkeService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	cvmService := svccvm.NewCvmService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())

	runInstancesPara := cvm.NewRunInstancesRequest()

	var place cvm.Placement
	if v, ok := d.GetOkExists("zone_id"); ok {
		var zones []*cvm.ZoneInfo
		err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			var e error
			zones, e = cvmService.DescribeZones(ctx)
			if e != nil {
				return tccommon.RetryError(e, tccommon.InternalError)
			}
			return nil
		})
		if err != nil {
			return err
		}
		zone := ""
		for _, z := range zones {
			if *z.ZoneId == v.(string) {
				zone = *z.Zone
				break
			}
		}
		place.Zone = helper.String(zone)
	}
	runInstancesPara.Placement = &place

	var basic ClusterBasicSetting
	var cAdvanced ClusterAdvancedSettings
	var cvms RunInstancesForNode
	var iAdvanced InstanceAdvancedSettings
	var cidrSet ClusterCidrSettings

	basic.ClusterName = d.Get("cluster_name").(string)

	if clusterDescRaw, ok := d.GetOkExists("cluster_desc"); ok {
		if len(clusterDescRaw.(string)) > 0 {
			basic.ClusterDescription = clusterDescRaw.(string)
		}
	}

	if osNameRaw, ok := d.GetOkExists("os_name"); ok {
		if len(osNameRaw.(string)) > 0 {
			basic.ClusterOs = osNameRaw.(string)
		}
	}

	if clusterVersionRaw, ok := d.GetOkExists("cluster_version"); ok {
		if len(clusterVersionRaw.(string)) > 0 {
			basic.ClusterVersion = clusterVersionRaw.(string)
		}
	}

	if v, ok := d.GetOkExists("vpc_id"); ok {
		if len(v.(string)) > 0 {
			vpcId := v.(string)
			subnetId := ""
			asVpcGateway := false
			if subnetIdRaw, ok := d.GetOkExists("subnet_id"); ok {
				subnetId = subnetIdRaw.(string)
			}
			if isVpcGatewayRaw, ok := d.GetOkExists("is_vpc_gateway"); ok {
				if isVpcGatewayRaw.(int) == 1 {
					asVpcGateway = true
				}
			}
			runInstancesPara.VirtualPrivateCloud = &cvm.VirtualPrivateCloud{
				VpcId:        common.StringPtr(vpcId),
				SubnetId:     common.StringPtr(subnetId),
				AsVpcGateway: common.BoolPtr(asVpcGateway),
			}
			basic.VpcId = vpcId
		}
	}

	if instanceTypeRaw, ok := d.GetOkExists("instance_type"); ok {
		if len(instanceTypeRaw.(string)) > 0 {
			runInstancesPara.InstanceType = common.StringPtr(instanceTypeRaw.(string))
		}
	}

	if v, ok := d.GetOkExists("require_wan_ip"); ok {
		publicIpAssigned := false
		internetMaxBandwidthOut := int64(0)
		internetChargeType := ""
		if v.(int) == 1 {
			publicIpAssigned = true
			if v, ok := d.GetOkExists("bandwidth"); ok {
				internetMaxBandwidthOut = int64(v.(int))
			}
			if v, ok := d.GetOkExists("bandwidth_type"); ok {
				bandwidthTypes := map[string]string{
					"PayByMonth":   "BANDWIDTH_PREPAID",
					"PayByTraffic": "TRAFFIC_POSTPAID_BY_HOUR",
					"PayByHour":    "TRAFFIC_POSTPAID_BY_HOUR",
				}
				if v, ok := bandwidthTypes[v.(string)]; ok {
					internetChargeType = v
				}
			}
		}
		runInstancesPara.InternetAccessible = &cvm.InternetAccessible{
			PublicIpAssigned:        common.BoolPtr(publicIpAssigned),
			InternetMaxBandwidthOut: common.Int64Ptr(internetMaxBandwidthOut),
			InternetChargeType:      common.StringPtr(internetChargeType),
		}
	}

	if v, ok := d.GetOkExists("instance_name"); ok {
		runInstancesPara.InstanceName = common.StringPtr(v.(string))
	}

	if v, ok := d.GetOkExists("goods_num"); ok {
		runInstancesPara.InstanceCount = common.Int64Ptr(int64(v.(int)))
	}

	if v, ok := d.GetOkExists("sg_id"); ok {
		runInstancesPara.SecurityGroupIds = []*string{common.StringPtr(v.(string))}
	}

	if v, ok := d.GetOkExists("password"); ok {
		runInstancesPara.LoginSettings = &cvm.LoginSettings{
			Password: common.StringPtr(v.(string)),
		}
	}

	if v, ok := d.GetOkExists("key_id"); ok {
		runInstancesPara.LoginSettings = &cvm.LoginSettings{
			KeyIds: []*string{common.StringPtr(v.(string))},
		}
	}

	runInstancesPara.SystemDisk = &cvm.SystemDisk{}
	if v, ok := d.GetOkExists("root_size"); ok {
		runInstancesPara.SystemDisk.DiskSize = common.Int64Ptr(int64(v.(int)))
	}

	if v, ok := d.GetOkExists("root_type"); ok {
		runInstancesPara.SystemDisk.DiskType = common.StringPtr(v.(string))
	}

	if v, ok := d.GetOkExists("storage_size"); ok {
		if v.(int) > 0 {
			dataDisk := &cvm.DataDisk{
				DiskSize: common.Int64Ptr(int64(v.(int))),
				DiskType: common.StringPtr("CLOUD_PREMIUM"),
			}
			if v, ok := d.GetOkExists("storage_type"); ok {
				if len(v.(string)) > 0 {
					dataDisk.DiskType = common.StringPtr(v.(string))
				}
			}
			runInstancesPara.DataDisks = []*cvm.DataDisk{dataDisk}
		}
	}

	if v, ok := d.GetOkExists("cvm_type"); ok {
		cvmTypes := map[string]string{
			"PayByHour":  "POSTPAID_BY_HOUR",
			"PayByMonth": "PREPAID",
		}
		if vv, ok := cvmTypes[v.(string)]; ok {
			runInstancesPara.InstanceChargeType = common.StringPtr(vv)
		}
	}

	if v, ok := d.GetOkExists("period"); ok {
		runInstancesPara.InstanceChargePrepaid = &cvm.InstanceChargePrepaid{
			Period: common.Int64Ptr(int64(v.(int))),
		}
	}

	if clusterCidrRaw, ok := d.GetOkExists("cluster_cidr"); ok {
		cidrSet.ClusterCidr = clusterCidrRaw.(string)
	}

	if v, ok := d.GetOkExists("mount_target"); ok {
		iAdvanced.MountTarget = v.(string)
	}

	if v, ok := d.GetOkExists("docker_graph_path"); ok {
		iAdvanced.DockerGraphPath = v.(string)
	}

	if v, ok := d.GetOkExists("user_script"); ok {
		iAdvanced.UserScript = v.(string)
	}

	if v, ok := d.GetOkExists("unschedulable"); ok {
		iAdvanced.Unschedulable = int64(v.(int))
	}

	runInstancesParas := runInstancesPara.ToJsonString()
	cvms.Work = []string{runInstancesParas}

	id, err := service.CreateCluster(ctx, basic, cAdvanced, cvms, iAdvanced, cidrSet, map[string]string{}, nil, nil, nil, nil)
	if err != nil {
		return err
	}

	d.SetId(id)
	err = resource.Retry(6*tccommon.WriteRetryTimeout, func() *resource.RetryError {
		_, _, err = service.DescribeClusterInstances(ctx, d.Id())
		if e, ok := err.(*errors.TencentCloudSDKError); ok {
			if e.GetCode() == "InternalError.ClusterNotFound" {
				return nil
			}
		}
		if err != nil {
			return resource.RetryableError(err)
		}
		return nil
	})
	if err != nil {
		return err
	}

	if err = resourceTencentCloudContainerClusterRead(d, meta); err != nil {
		log.Printf("[WARN]%s resource.tencentcloud_container_cluster.read after create fail , %s", logId, err.Error())
	}

	return nil
}

func resourceTencentCloudContainerClusterRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_container_cluster.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	service := TkeService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var info ClusterInfo
	var has bool

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		var e error
		info, has, e = service.DescribeCluster(ctx, d.Id())
		if e != nil {
			return tccommon.RetryError(e)
		}
		return nil
	})
	if err != nil {
		return nil
	}

	if !has {
		d.SetId("")
		return nil
	}

	_ = d.Set("cluster_name", info.ClusterName)
	_ = d.Set("cluster_desc", info.ClusterDescription)
	_ = d.Set("kubernetes_version", info.ClusterVersion)
	_ = d.Set("vpc_id", info.VpcId)
	_ = d.Set("cluster_cidr", info.ClusterCidr)
	_ = d.Set("nodes_num", info.ClusterNodeNum)
	_ = d.Set("nodes_status", info.ClusterStatus)
	_ = d.Set("total_cpu", 0)
	_ = d.Set("total_mem", 0)

	var workers []InstanceInfo
	err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		var e error
		_, workers, e = service.DescribeClusterInstances(ctx, info.ClusterId)
		if e, ok := e.(*errors.TencentCloudSDKError); ok {
			if e.GetCode() == "InternalError.ClusterNotFound" {
				d.SetId("")
				return nil
			}
		}
		if e != nil {
			return resource.RetryableError(e)
		}
		return nil
	})
	if err != nil {
		return err
	}

	instanceIds := []*string{}
	for _, v := range workers {
		instanceIds = append(instanceIds, common.StringPtr(v.InstanceId))
	}

	if len(instanceIds) > 0 {
		describeInstancesreq := cvm.NewDescribeInstancesRequest()
		describeInstancesreq.InstanceIds = instanceIds
		var describeInstancesResponse *cvm.DescribeInstancesResponse
		err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCvmClient().DescribeInstances(describeInstancesreq)
			if e != nil {
				log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
					logId, describeInstancesreq.GetAction(), describeInstancesreq.ToJsonString(), e.Error())
				return tccommon.RetryError(e)
			}
			describeInstancesResponse = result
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s DescribeInstances failed, reason:%s\n ", logId, err.Error())
			return err
		}

		total_cpu := int64(0)
		total_mem := int64(0)
		for _, v := range describeInstancesResponse.Response.InstanceSet {
			total_cpu += *v.CPU
			total_mem += *v.Memory
		}

		_ = d.Set("total_cpu", total_cpu)
		_ = d.Set("total_mem", total_mem)
	}

	return nil
}

func resourceTencentCloudContainerClusterDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_container_cluster.delete")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	service := TkeService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	return resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		err := service.DeleteCluster(ctx, d.Id())
		if err != nil {
			return tccommon.RetryError(err)
		}
		return nil
	})
}
