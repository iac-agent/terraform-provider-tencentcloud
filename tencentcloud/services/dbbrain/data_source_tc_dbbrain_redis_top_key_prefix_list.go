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

func DataSourceTencentCloudDbbrainRedisTopKeyPrefixList() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudDbbrainRedisTopKeyPrefixListRead,
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

			"items": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "列表 的 top 键 prefixes.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ave_element_size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Average element 长度.",
						},
						"length": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Total occupied 内存 (Byte).",
						},
						"key_pre_index": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "键 prefix.",
						},
						"item_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "数量 的 elements.",
						},
						"count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "数量 的 keys.",
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

func dataSourceTencentCloudDbbrainRedisTopKeyPrefixListRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_dbbrain_redis_top_key_prefix_list.read")()
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

	service := DbbrainService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var items []*dbbrain.RedisPreKeySpaceData

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeDbbrainRedisTopKeyPrefixListByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		items = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(items))
	tmpList := make([]map[string]interface{}, 0, len(items))

	if items != nil {
		for _, redisPreKeySpaceData := range items {
			redisPreKeySpaceDataMap := map[string]interface{}{}

			if redisPreKeySpaceData.AveElementSize != nil {
				redisPreKeySpaceDataMap["ave_element_size"] = redisPreKeySpaceData.AveElementSize
			}

			if redisPreKeySpaceData.Length != nil {
				redisPreKeySpaceDataMap["length"] = redisPreKeySpaceData.Length
			}

			if redisPreKeySpaceData.KeyPreIndex != nil {
				redisPreKeySpaceDataMap["key_pre_index"] = redisPreKeySpaceData.KeyPreIndex
			}

			if redisPreKeySpaceData.ItemCount != nil {
				redisPreKeySpaceDataMap["item_count"] = redisPreKeySpaceData.ItemCount
			}

			if redisPreKeySpaceData.Count != nil {
				redisPreKeySpaceDataMap["count"] = redisPreKeySpaceData.Count
			}

			if redisPreKeySpaceData.MaxElementSize != nil {
				redisPreKeySpaceDataMap["max_element_size"] = redisPreKeySpaceData.MaxElementSize
			}

			ids = append(ids, strings.Join([]string{instanceId, *redisPreKeySpaceData.KeyPreIndex}, tccommon.FILED_SP))
			tmpList = append(tmpList, redisPreKeySpaceDataMap)
		}

		_ = d.Set("items", tmpList)
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
