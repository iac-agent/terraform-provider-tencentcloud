package dbbrain

import (
	"context"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dbbrain "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dbbrain/v20210527"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudDbbrainRedisTopBigKeys() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudDbbrainRedisTopBigKeysRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "实例 ID.",
			},

			"date": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Query date, such 作为 2021-05-27, earliest date 可以 是 previous 30 days.",
			},

			"product": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Service product 类型, 支持 值 include `redis` - 云 数据库 Redis.",
			},

			"sort_by": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Sorting 字段, 值 includes `Capacity` - 内存, `ItemCount` - 数量 的 elements, 默认值 是 `Capacity`.",
			},

			"key_type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Key 类型 过滤器 condition, 默认值 是 无 过滤器, 值 includes `字符串`, `列表`, `集合`, `hash`, `sortedset`, `流`.",
			},

			"top_keys": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "列表 的 top keys.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "键 名称.",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "键 类型.",
						},
						"encoding": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "键 编码 方法.",
						},
						"expire_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Key expiration timestamp (在 milliseconds), 0 表示 无 expiration 时间 是 集合.",
						},
						"length": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Key 内存 大小, 单位 Byte.",
						},
						"item_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "数量 的 elements.",
						},
						"max_element_size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Maximum element 长度.",
						},
					},
				},
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used 到 save results.",
			},
		},
	}
}

func dataSourceTencentCloudDbbrainRedisTopBigKeysRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_dbbrain_redis_top_big_keys.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	var (
		instanceId string
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		paramMap["InstanceId"] = helper.String(v.(string))
		instanceId = v.(string)
	}

	if v, ok := d.GetOk("date"); ok {
		paramMap["Date"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("product"); ok {
		paramMap["Product"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("sort_by"); ok {
		paramMap["SortBy"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("key_type"); ok {
		paramMap["KeyType"] = helper.String(v.(string))
	}

	service := DbbrainService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var topKeys []*dbbrain.RedisKeySpaceData

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeDbbrainRedisTopBigKeysByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		topKeys = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(topKeys))
	tmpList := make([]map[string]interface{}, 0, len(topKeys))

	if topKeys != nil {
		for _, redisKeySpaceData := range topKeys {
			redisKeySpaceDataMap := map[string]interface{}{}

			if redisKeySpaceData.Key != nil {
				redisKeySpaceDataMap["key"] = redisKeySpaceData.Key
			}

			if redisKeySpaceData.Type != nil {
				redisKeySpaceDataMap["type"] = redisKeySpaceData.Type
			}

			if redisKeySpaceData.Encoding != nil {
				redisKeySpaceDataMap["encoding"] = redisKeySpaceData.Encoding
			}

			if redisKeySpaceData.ExpireTime != nil {
				redisKeySpaceDataMap["expire_time"] = redisKeySpaceData.ExpireTime
			}

			if redisKeySpaceData.Length != nil {
				redisKeySpaceDataMap["length"] = redisKeySpaceData.Length
			}

			if redisKeySpaceData.ItemCount != nil {
				redisKeySpaceDataMap["item_count"] = redisKeySpaceData.ItemCount
			}

			if redisKeySpaceData.MaxElementSize != nil {
				redisKeySpaceDataMap["max_element_size"] = redisKeySpaceData.MaxElementSize
			}

			ids = append(ids, strings.Join([]string{instanceId, *redisKeySpaceData.Key}, tccommon.FILED_SP))
			tmpList = append(tmpList, redisKeySpaceDataMap)
		}

		_ = d.Set("top_keys", tmpList)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
