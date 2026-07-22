package cdb

import (
	"context"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mysql "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdb/v20170320"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudMysqlRoInstanceIp() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMysqlRoInstanceIpCreate,
		Read:   resourceTencentCloudMysqlRoInstanceIpRead,
		Delete: resourceTencentCloudMysqlRoInstanceIpDelete,

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "只读实例ID，格式为：cdbro-3i70uj0k，与云数据库控制台页面显示的只读实例ID相同。",
			},

			"uniq_subnet_id": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "子网描述符，例如：subnet-1typ0s7d。",
			},

			"uniq_vpc_id": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "vpc描述符，例如：vpc-a23yt67j，如果传递该字段，则必须传递UniqSubnetId。",
			},

			"ro_vip": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "只读实例的内网IP地址。",
			},

			"ro_vport": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "只读实例的内网端口号。",
			},
		},
	}
}

func resourceTencentCloudMysqlRoInstanceIpCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mysql_ro_instance_ip.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request    = mysql.NewCreateRoInstanceIpRequest()
		instanceId string
	)
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		request.InstanceId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("uniq_subnet_id"); ok {
		request.UniqSubnetId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("uniq_vpc_id"); ok {
		request.UniqVpcId = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMysqlClient().CreateRoInstanceIp(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate mysql roInstanceIp failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(instanceId)

	return resourceTencentCloudMysqlRoInstanceIpRead(d, meta)
}

func resourceTencentCloudMysqlRoInstanceIpRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mysql_ro_instance_ip.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := MysqlService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	instanceId := d.Id()

	switchForUpgrade, err := service.DescribeDBInstanceById(ctx, instanceId)
	if err != nil {
		return err
	}

	if switchForUpgrade == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `MysqlSwitchForUpgrade` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if switchForUpgrade.InstanceId != nil {
		_ = d.Set("instance_id", switchForUpgrade.InstanceId)
	}

	if switchForUpgrade.UniqVpcId != nil {
		_ = d.Set("uniq_vpc_id", switchForUpgrade.UniqVpcId)
	}

	if switchForUpgrade.UniqSubnetId != nil {
		_ = d.Set("uniq_subnet_id", switchForUpgrade.UniqSubnetId)
	}

	if switchForUpgrade.RoVipInfo != nil {
		if switchForUpgrade.RoVipInfo.RoVip != nil {
			_ = d.Set("ro_vip", switchForUpgrade.RoVipInfo.RoVip)
		}

		if switchForUpgrade.RoVipInfo.RoVport != nil {
			_ = d.Set("ro_vport", switchForUpgrade.RoVipInfo.RoVport)
		}
	}

	return nil
}

func resourceTencentCloudMysqlRoInstanceIpDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mysql_ro_instance_ip.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
