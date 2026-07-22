package crs

import (
	"context"
	"fmt"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	redis "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/redis/v20180412"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudRedisUpgradeMultiZoneOperation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudRedisUpgradeMultiZoneOperationCreate,
		Read:   resourceTencentCloudRedisUpgradeMultiZoneOperationRead,
		Delete: resourceTencentCloudRedisUpgradeMultiZoneOperationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "ID 实例。",
			},

			"upgrade_proxy_and_redis_server": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeBool,
				Description: "After 您 upgrade Multi-AZ，是否nearby 访问 功能 是 支持.true: Supports nearby 访问. upgrade process，其中 requires upgrading both proxy 版本 和 Redis kernel minor 版本，involves 数据 迁移 和 可以 take several hours.false: No need 到 support nearby 访问.Upgrading Multi-AZ 仅 involves managing metadata 迁移，使用 无 服务 impact，和 upgrade process typically completes within 3 minutes。",
			},
		},
	}
}

func resourceTencentCloudRedisUpgradeMultiZoneOperationCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_redis_upgrade_multi_zone_operation.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	var (
		request    = redis.NewUpgradeVersionToMultiAvailabilityZonesRequest()
		response   = redis.NewUpgradeVersionToMultiAvailabilityZonesResponse()
		instanceId string
	)
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		request.InstanceId = helper.String(v.(string))
	}

	if v, _ := d.GetOk("upgrade_proxy_and_redis_server"); v != nil {
		request.UpgradeProxyAndRedisServer = helper.Bool(v.(bool))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseRedisClient().UpgradeVersionToMultiAvailabilityZones(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate redis upgradeMultiZoneOperation failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(instanceId)

	service := RedisService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	taskId := *response.Response.FlowId
	err = resource.Retry(6*tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ok, err := service.DescribeTaskInfo(ctx, instanceId, taskId)
		if err != nil {
			if _, ok := err.(*sdkErrors.TencentCloudSDKError); !ok {
				return resource.RetryableError(err)
			} else {
				return resource.NonRetryableError(err)
			}
		}
		if ok {
			return nil
		} else {
			return resource.RetryableError(fmt.Errorf("upgrade multi zone is processing"))
		}
	})

	if err != nil {
		log.Printf("[CRITAL]%s redis upgrade multi zone fail, reason:%s\n", logId, err.Error())
		return err
	}

	return resourceTencentCloudRedisUpgradeMultiZoneOperationRead(d, meta)
}

func resourceTencentCloudRedisUpgradeMultiZoneOperationRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_redis_upgrade_multi_zone_operation.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudRedisUpgradeMultiZoneOperationDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_redis_upgrade_multi_zone_operation.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
