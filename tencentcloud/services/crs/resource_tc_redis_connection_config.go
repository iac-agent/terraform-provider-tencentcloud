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

func ResourceTencentCloudRedisConnectionConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudRedisConnectionConfigCreate,
		Read:   resourceTencentCloudRedisConnectionConfigRead,
		Update: resourceTencentCloudRedisConnectionConfigUpdate,
		Delete: resourceTencentCloudRedisConnectionConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "ID 实例。",
			},

			"client_limit": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "总数 数量 connections per 分片.如果 read-仅 replicas 是 不 已启用， lower 限制 是 10,000 和 upper 限制 是 40,000.当 您 启用 read-仅 replicas， 最小 限制 是 10,000 和 upper 限制 是 10,000 * ( 数量 read replicas +3)。",
			},

			"total_bandwidth": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Total 带宽 的 实例 = additional 带宽 * 数量 shards + standard 带宽 * 数量 shards * (数量 primary nodes + 数量 read-仅 副本 nodes)， 数量 shards 的 standard architecture = 1，在 Mb/s。",
			},

			"base_bandwidth": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "standard 带宽. Refers 到 带宽 allocated 通过 系统 到 each 节点 当 实例 是 purchased。",
			},

			"add_bandwidth": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Refers 到 additional 带宽 的 实例. 当 standard 带宽 does 不 meet demand， 用户 可以 increase 带宽 通过 himself. 当 read-仅 copy 是 已启用， 总数 带宽 的 实例 = additional 带宽 * 数量 fragments + standard 带宽 * 数量 fragments * Max ([数量 read-仅 replicas，1] )， 数量 shards 在 standard architecture = 1，和 当 read-仅 replicas 是 不 已启用， 总数 带宽 的 实例 = additional 带宽 * 数量 shards + standard 带宽 * 数量 shards，和 数量 shards 在 standard architecture = 1。",
			},

			"min_add_bandwidth": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Additional 带宽 sets lower 限制",
			},

			"max_add_bandwidth": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Additional 带宽 是 capped。",
			},
		},
	}
}

func resourceTencentCloudRedisConnectionConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_redis_connection_config.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		instanceId string
	)
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}

	d.SetId(instanceId)

	return resourceTencentCloudRedisConnectionConfigUpdate(d, meta)
}

func resourceTencentCloudRedisConnectionConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_redis_connection_config.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := RedisService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	instanceId := d.Id()

	connectionConfig, err := service.DescribeRedisInstanceById(ctx, instanceId)
	if err != nil {
		return err
	}

	if connectionConfig == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `RedisConnectionConfig` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if connectionConfig.InstanceId != nil {
		_ = d.Set("instance_id", connectionConfig.InstanceId)
	}

	if connectionConfig.ClientLimit != nil {
		_ = d.Set("client_limit", connectionConfig.ClientLimit)
	}

	if connectionConfig.NetLimit != nil && connectionConfig.RedisShardNum != nil {
		netLimt := *connectionConfig.NetLimit
		shardNum := *connectionConfig.RedisShardNum
		_ = d.Set("total_bandwidth", netLimt*shardNum*8)
	}

	bandwidthRange, err := service.DescribeBandwidthRangeById(ctx, instanceId)
	if err != nil {
		return err
	}

	if connectionConfig == nil {
		log.Printf("[WARN]%s resource `DescribeBandwidthRangeById` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if bandwidthRange.BaseBandwidth != nil {
		_ = d.Set("base_bandwidth", bandwidthRange.BaseBandwidth)
	}
	if bandwidthRange.AddBandwidth != nil {
		_ = d.Set("add_bandwidth", bandwidthRange.AddBandwidth)
	}
	if bandwidthRange.MinAddBandwidth != nil {
		_ = d.Set("min_add_bandwidth", bandwidthRange.MinAddBandwidth)
	}
	if bandwidthRange.MaxAddBandwidth != nil {
		_ = d.Set("max_add_bandwidth", bandwidthRange.MaxAddBandwidth)
	}

	return nil
}

func resourceTencentCloudRedisConnectionConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_redis_connection_config.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	request := redis.NewModifyConnectionConfigRequest()
	response := redis.NewModifyConnectionConfigResponse()

	instanceId := d.Id()
	request.InstanceId = &instanceId

	if v, ok := d.GetOkExists("client_limit"); ok {
		request.ClientLimit = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("add_bandwidth"); ok {
		request.Bandwidth = helper.IntInt64(v.(int))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseRedisClient().ModifyConnectionConfig(request)
		if e != nil {
			if ee, ok := e.(*sdkErrors.TencentCloudSDKError); ok {
				if ee.Code == "FailedOperation.SystemError" {
					return resource.NonRetryableError(e)
				}
			}
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update redis param failed, reason:%+v", logId, err)
		return err
	}

	service := RedisService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	taskId := *response.Response.TaskId
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
			return resource.RetryableError(fmt.Errorf("change account is processing"))
		}
	})

	if err != nil {
		log.Printf("[CRITAL]%s redis change connection fail, reason:%s\n", logId, err.Error())
		return err
	}

	return resourceTencentCloudRedisConnectionConfigRead(d, meta)
}

func resourceTencentCloudRedisConnectionConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_redis_connection_config.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
