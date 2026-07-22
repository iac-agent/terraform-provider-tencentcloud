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
	tke "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tke/v20180525"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudContainerClusterInstance() *schema.Resource {
	return &schema.Resource{
		DeprecationMessage: "This resource has been deprecated in Terraform TencentCloud provider version 1.16.0. Please use 'tencentcloud_kubernetes_scale_worker' instead.",
		Create:             resourceTencentCloudContainerClusterInstancesCreate,
		Read:               resourceTencentCloudContainerClusterInstancesRead,
		Update:             resourceTencentCloudContainerClusterInstancesUpdate,
		Delete:             resourceTencentCloudContainerClusterInstancesDelete,

		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID 集群。",
			},
			"cpu": {
				Type:        schema.TypeInt,
				Optional:    true,
				Deprecated:  "It has been deprecated from version 1.16.0. Set 'instance_type' instead.",
				Description: "cpu 的 节点。",
			},
			"mem": {
				Type:        schema.TypeInt,
				Optional:    true,
				Deprecated:  "It has been deprecated from version 1.16.0. Set 'instance_type' instead.",
				Description: "内存 的 节点。",
			},
			"bandwidth": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "网络 带宽 的 节点。",
			},
			"bandwidth_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "网络 类型 节点。",
			},
			"require_wan_ip": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Indicate whether wan ip 是 needed。",
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "子网 ID 其中 节点 stays 在。",
			},
			"is_vpc_gateway": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Describe 是否node 启用 网关 capability。",
			},
			"storage_size": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "大小 的 数据 卷。",
			},
			"storage_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "类型 数据 卷. see more 从 CVM。",
			},
			"root_size": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "大小 的 root 卷。",
			},
			"root_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "类型 root 卷. see more 从 CVM。",
			},
			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "密码 的 each 节点。",
			},
			"key_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "key_id 的 each 节点(如果 使用 键 pair 到 访问)。",
			},
			"cvm_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "类型 节点 needed 通过 cvm。",
			},
			"period": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "puchase 时长 的 节点 needed 通过 cvm。",
			},
			"zone_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "可用区 其中 节点 stays 在。",
			},
			"instance_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "实例 类型 节点 needed 通过 cvm。",
			},
			"sg_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "安全组 ID",
			},
			"mount_target": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "路径 其中 卷 是 going 到 是 mounted。",
			},
			"docker_graph_path": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "docker graph 路径 是 going 到 mounted。",
			},
			"user_script": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "用户 defined 脚本 在 base64-格式 脚本 runs after kubernetes 组件 是 ready 在 节点. see more 从 CCS api documents。",
			},
			"unschedulable": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Determine 是否node 将 是 schedulable. 0 是 默认值 meaning 节点 将 是 schedulable. 1 对于 unschedulable。",
			},
			"instance_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "名称 ot 节点。",
			},
			"abnormal_reason": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Describe reason 当 节点 是 在 abnormal state(如果 它 是)。",
			},
			"instance_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "An ID identify 节点，提供 通过 cvm。",
			},
			"is_normal": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Describe 是否node 是 normal。",
			},
			"wan_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Describe wan ip 的 节点。",
			},
			"lan_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Describe lan ip 的 节点。",
			},
		},
	}
}

func resourceTencentCloudContainerClusterInstancesUpdate(d *schema.ResourceData, m interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_container_cluster_instance.update")()

	return fmt.Errorf("the container cluster instance resource doesn't support update")
}

func resourceTencentCloudContainerClusterInstancesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_container_cluster_instance.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	service := TkeService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var workers []InstanceInfo
	clusterId := d.Get("cluster_id").(string)

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		var e error
		_, workers, e = service.DescribeClusterInstances(ctx, clusterId)
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

	found := false
	for _, v := range workers {
		if v.InstanceId == d.Id() {
			found = true
			_ = d.Set("instance_id", v.InstanceId)
			_ = d.Set("abnormal_reason", v.FailedReason)
			_ = d.Set("wan_ip", "")
			_ = d.Set("lan_ip", "")
			_ = d.Set("is_normal", true)
			if v.InstanceState == "failed" {
				_ = d.Set("is_normal", false)
			}

			describeInstancesreq := cvm.NewDescribeInstancesRequest()
			describeInstancesreq.InstanceIds = []*string{common.StringPtr(v.InstanceId)}
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

			if len(describeInstancesResponse.Response.InstanceSet) > 0 {
				if len(describeInstancesResponse.Response.InstanceSet[0].PublicIpAddresses) > 0 {
					_ = d.Set("wan_ip", *describeInstancesResponse.Response.InstanceSet[0].PublicIpAddresses[0])
				}
				if len(describeInstancesResponse.Response.InstanceSet[0].PrivateIpAddresses) > 0 {
					_ = d.Set("lan_ip", *describeInstancesResponse.Response.InstanceSet[0].PrivateIpAddresses[0])
				}
			}
		}
	}

	if !found {
		d.SetId("")
	}

	return nil
}

func resourceTencentCloudContainerClusterInstancesCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_container_cluster_instance.create")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	service := TkeService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	cvmService := svccvm.NewCvmService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())

	clusterId := d.Get("cluster_id").(string)
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
	runInstancesPara.InstanceCount = common.Int64Ptr(1)

	var iAdvanced tke.InstanceAdvancedSettings
	var cvms RunInstancesForNode

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

	if v, ok := d.GetOkExists("user_script"); ok {
		runInstancesPara.UserData = common.StringPtr(v.(string))
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

	if v, ok := d.GetOkExists("instance_name"); ok {
		runInstancesPara.InstanceName = common.StringPtr(v.(string))
	}

	if v, ok := d.GetOkExists("mount_target"); ok {
		iAdvanced.MountTarget = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("docker_graph_path"); ok {
		iAdvanced.DockerGraphPath = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("user_script"); ok {
		iAdvanced.UserScript = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("unschedulable"); ok {
		iAdvanced.Unschedulable = helper.IntInt64(v.(int))
	}

	runInstancesParas := runInstancesPara.ToJsonString()
	cvms.Work = []string{runInstancesParas}

	instanceIds, err := service.CreateClusterInstances(ctx, clusterId, cvms.Work[0], iAdvanced)
	if err != nil {
		return err
	}

	if len(instanceIds) != 1 {
		return fmt.Errorf("no instance return")
	}

	err = resource.Retry(6*tccommon.WriteRetryTimeout, func() *resource.RetryError {
		_, workers, e := service.DescribeClusterInstances(ctx, clusterId)
		if ee, ok := e.(*errors.TencentCloudSDKError); ok {
			if ee.GetCode() == "InternalError.ClusterNotFound" {
				return nil
			}
		}
		if err != nil {
			return resource.RetryableError(err)
		}
		for _, v := range workers {
			if v.InstanceId != instanceIds[0] {
				continue
			}
			if v.InstanceState != "running" {
				return resource.RetryableError(fmt.Errorf("instance state not ready"))
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	d.SetId(instanceIds[0])

	return resourceTencentCloudContainerClusterInstancesRead(d, meta)
}

func resourceTencentCloudContainerClusterInstancesDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_container_cluster_instance.delete")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	service := TkeService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var workers []InstanceInfo
	clusterId := d.Get("cluster_id").(string)

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		var e error
		_, workers, e = service.DescribeClusterInstances(ctx, clusterId)
		if ee, ok := e.(*errors.TencentCloudSDKError); ok {
			if ee.GetCode() == "InternalError.ClusterNotFound" {
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

	found := false
	var node InstanceInfo
	for _, v := range workers {
		if v.InstanceId == d.Id() {
			found = true
			node = v
		}
	}
	if !found {
		return nil
	}

	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		e := service.DeleteClusterInstances(ctx, clusterId, []string{node.InstanceId})
		if ee, ok := e.(*errors.TencentCloudSDKError); ok {
			if ee.GetCode() == "InternalError.ClusterNotFound" {
				return nil
			}
		}
		if err != nil {
			return tccommon.RetryError(err)
		}
		return nil
	})

	return err
}
