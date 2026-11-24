package oceanus

import (
	"context"
	"fmt"
	"log"
	"time"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	oceanus "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/oceanus/v20190422"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudOceanusOceanusCluster() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudOceanusOceanusClusterCreate,
		Read:   resourceTencentCloudOceanusOceanusClusterRead,
		Update: resourceTencentCloudOceanusOceanusClusterUpdate,
		Delete: resourceTencentCloudOceanusOceanusClusterDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"cluster_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Cluster name. Supports 1-50 english letters, chinese characters, numbers, hyphens or underscores.",
			},
			"region_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "Region ID, which can be queried via DescribeRegions with product oceanus.",
			},
			"zone_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "Zone ID, which can be queried via DescribeZones with product oceanus.",
			},
			"login_password": {
				Required:    true,
				Type:        schema.TypeString,
				Sensitive:   true,
				Description: "Flink UI access password. The cluster user name is admin.",
			},
			"vpc_descriptions": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeList,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vpc_id": {
							Required:    true,
							Type:        schema.TypeString,
							Description: "VPC ID.",
						},
						"subnet_id": {
							Required:    true,
							Type:        schema.TypeString,
							Description: "Subnet ID.",
						},
						"app_id": {
							Optional:    true,
							Computed:    true,
							Type:        schema.TypeInt,
							Description: "User AppId.",
						},
						"owner_uin": {
							Optional:    true,
							Computed:    true,
							Type:        schema.TypeString,
							Description: "User UIN.",
						},
					},
				},
				Description: "VPC descriptions. Stream compute accesses other cloud resources in the same region via VPC and ENI and occupies subnet IPs.",
			},
			"default_cos_bucket": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Default COS bucket used by stream compute to save checkpoints, jars or logs.",
			},
			"cu": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Cluster CU number. 12 CU is the minimum compute resource and billing unit, 1 CU contains 1 CPU and 4GB memory. Management node fee is waived when CU >= 48.",
			},
			"remark": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Cluster description. Supports 1-50 english letters, chinese characters, numbers, hyphens or underscores.",
			},
			"period": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Purchase period in months, range [1, 48]. Only valid for PREPAID clusters.",
			},
			"instance_charge_type": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Cluster charge type. Valid values: PREPAID, POSTPAID_BY_SECOND. Default: POSTPAID_BY_SECOND.",
			},
			"cluster_type": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Cluster type. Valid value: MULTI_AZ_CLUSTER.",
			},
			"renew_flag": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Auto renewal flag. Valid values: NOTIFY_AND_MANUAL_RENEW, NOTIFY_AND_AUTO_RENEW, DISABLE_NOTIFY_AND_MANUAL_RENEW. Default: NOTIFY_AND_MANUAL_RENEW. Only valid for PREPAID clusters.",
			},
			"flink_web_ui_network_access_type": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Flink UI network access type. Valid values: NetworkAccess_INTERNAL, NetworkAccess_EXTERNAL. Default: NetworkAccess_EXTERNAL.",
			},
			"slave_vpc_descriptions": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeList,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vpc_id": {
							Required:    true,
							Type:        schema.TypeString,
							Description: "VPC ID.",
						},
						"subnet_id": {
							Required:    true,
							Type:        schema.TypeString,
							Description: "Subnet ID.",
						},
						"app_id": {
							Optional:    true,
							Computed:    true,
							Type:        schema.TypeInt,
							Description: "User AppId.",
						},
						"owner_uin": {
							Optional:    true,
							Computed:    true,
							Type:        schema.TypeString,
							Description: "User UIN.",
						},
					},
				},
				Description: "Multi-AZ VPC descriptions. Only valid for MULTI_AZ_CLUSTER.",
			},
			"cu_memory": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "CU memory ratio, only supports [0, 2, 4, 8].",
			},
			// computed
			"status": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Cluster status. 1: uninitialized, 2: running, 3: initializing.",
			},
			"create_time": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Cluster create time.",
			},
			"update_time": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Cluster last operation time.",
			},
		},
	}
}

func resourceTencentCloudOceanusOceanusClusterCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_oceanus_oceanus_cluster.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId     = tccommon.GetLogId(tccommon.ContextNil)
		request   = oceanus.NewCreateOceanusClusterRequest()
		response  = oceanus.NewCreateOceanusClusterResponse()
		clusterId string
	)

	if v, ok := d.GetOk("cluster_name"); ok {
		request.ClusterName = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("region_id"); ok {
		request.RegionId = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("zone_id"); ok {
		request.ZoneId = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("login_password"); ok {
		request.LoginPassword = helper.String(v.(string))
	}

	if v, ok := d.GetOk("vpc_descriptions"); ok {
		vpcDescriptions := v.([]interface{})
		if len(vpcDescriptions) > 0 {
			vpcMap := vpcDescriptions[0].(map[string]interface{})
			vpcDescription := oceanus.VPCDescription{}
			if vpcId, ok := vpcMap["vpc_id"]; ok && vpcId.(string) != "" {
				vpcDescription.VpcId = helper.String(vpcId.(string))
			}
			if subnetId, ok := vpcMap["subnet_id"]; ok && subnetId.(string) != "" {
				vpcDescription.SubnetId = helper.String(subnetId.(string))
			}
			if appId, ok := vpcMap["app_id"]; ok {
				vpcDescription.AppId = helper.IntInt64(appId.(int))
			}
			if ownerUin, ok := vpcMap["owner_uin"]; ok && ownerUin.(string) != "" {
				vpcDescription.OwnerUin = helper.String(ownerUin.(string))
			}
			request.VpcDescriptions = []*oceanus.VPCDescription{&vpcDescription}
		}
	}

	if v, ok := d.GetOk("default_cos_bucket"); ok {
		request.DefaultCOSBucket = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("cu"); ok {
		request.CU = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("remark"); ok {
		request.Remark = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("period"); ok {
		request.Period = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("instance_charge_type"); ok {
		request.InstanceChargeType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("cluster_type"); ok {
		request.ClusterType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("renew_flag"); ok {
		request.RenewFlag = helper.String(v.(string))
	}

	if v, ok := d.GetOk("flink_web_ui_network_access_type"); ok {
		request.FlinkWebUINetworkAccessType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("slave_vpc_descriptions"); ok {
		slaveVpcDescriptions := v.([]interface{})
		if len(slaveVpcDescriptions) > 0 {
			slaveVpcMap := slaveVpcDescriptions[0].(map[string]interface{})
			slaveVpcDescription := oceanus.SlaveVpcDescriptions{}
			if vpcId, ok := slaveVpcMap["vpc_id"]; ok && vpcId.(string) != "" {
				slaveVpcDescription.VpcId = helper.String(vpcId.(string))
			}
			if subnetId, ok := slaveVpcMap["subnet_id"]; ok && subnetId.(string) != "" {
				slaveVpcDescription.SubnetId = helper.String(subnetId.(string))
			}
			if appId, ok := slaveVpcMap["app_id"]; ok {
				slaveVpcDescription.AppId = helper.IntInt64(appId.(int))
			}
			if ownerUin, ok := slaveVpcMap["owner_uin"]; ok && ownerUin.(string) != "" {
				slaveVpcDescription.OwnerUin = helper.String(ownerUin.(string))
			}
			request.SlaveVpcDescriptions = []*oceanus.SlaveVpcDescriptions{&slaveVpcDescription}
		}
	}

	if v, ok := d.GetOkExists("cu_memory"); ok {
		request.CUMemory = helper.IntInt64(v.(int))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseOceanusClient().CreateOceanusCluster(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil || result.Response.ClusterId == nil {
			e = fmt.Errorf("oceanus OceanusCluster not exists")
			return resource.NonRetryableError(e)
		}

		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create oceanus OceanusCluster failed, reason:%+v", logId, err)
		return err
	}

	clusterId = *response.Response.ClusterId
	d.SetId(clusterId)

	err = resource.RetryContext(context.TODO(), 15*time.Minute, func() *resource.RetryError {
		service := OceanusService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		cluster, e := service.DescribeOceanusClusterById(context.WithValue(context.TODO(), tccommon.LogIdKey, logId), clusterId)
		if e != nil {
			return resource.NonRetryableError(e)
		}
		if cluster == nil {
			return resource.NonRetryableError(fmt.Errorf("oceanus OceanusCluster %s not exists", clusterId))
		}
		if cluster.Status != nil && *cluster.Status == 2 {
			return nil
		}
		return resource.RetryableError(fmt.Errorf("oceanus OceanusCluster %s status is %d, retrying", clusterId, *cluster.Status))
	})

	if err != nil {
		log.Printf("[CRITAL]%s waiting for oceanus OceanusCluster [%s] running failed, reason:%+v", logId, clusterId, err)
		return err
	}

	return resourceTencentCloudOceanusOceanusClusterRead(d, meta)
}

func resourceTencentCloudOceanusOceanusClusterRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_oceanus_oceanus_cluster.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = OceanusService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	clusterId := d.Id()

	cluster, err := service.DescribeOceanusClusterById(ctx, clusterId)
	if err != nil {
		return err
	}

	if cluster == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `OceanusOceanusCluster` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if cluster.ClusterId != nil {
		_ = d.Set("cluster_id", cluster.ClusterId)
	}

	if cluster.Name != nil {
		_ = d.Set("cluster_name", cluster.Name)
	}

	if cluster.Remark != nil {
		_ = d.Set("remark", cluster.Remark)
	}

	if cluster.CuNum != nil {
		_ = d.Set("cu", int(*cluster.CuNum))
	}

	if cluster.CuMem != nil {
		_ = d.Set("cu_memory", int(*cluster.CuMem))
	}

	if cluster.DefaultCOSBucket != nil {
		_ = d.Set("default_cos_bucket", cluster.DefaultCOSBucket)
	}

	if cluster.Status != nil {
		_ = d.Set("status", int(*cluster.Status))
	}

	if cluster.CreateTime != nil {
		_ = d.Set("create_time", cluster.CreateTime)
	}

	if cluster.UpdateTime != nil {
		_ = d.Set("update_time", cluster.UpdateTime)
	}

	return nil
}

func resourceTencentCloudOceanusOceanusClusterUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_oceanus_oceanus_cluster.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	return resourceTencentCloudOceanusOceanusClusterRead(d, meta)
}

func resourceTencentCloudOceanusOceanusClusterDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_oceanus_oceanus_cluster.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = OceanusService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	clusterId := d.Id()

	if err := service.DeleteOceanusOceanusClusterById(ctx, clusterId); err != nil {
		return err
	}

	return nil
}
