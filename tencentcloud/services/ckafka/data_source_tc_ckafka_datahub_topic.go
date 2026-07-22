package ckafka

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ckafka "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ckafka/v20190819"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudCkafkaDatahubTopic() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCkafkaDatahubTopicRead,
		Schema: map[string]*schema.Schema{
			"search_word": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "query 键 word。",
			},

			"offset": {
				Optional:    true,
				Default:     0,
				Type:        schema.TypeInt,
				Description: "The 偏移量 position of this query，the 默认为 0。",
			},

			"limit": {
				Optional:    true,
				Default:     50,
				Type:        schema.TypeInt,
				Description: "The 最大results returned this time，the 默认为 50，and the maximum 值 is 50。",
			},

			"topic_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Topic list。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "名称",
						},
						"topic_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Topic 名称",
						},
						"topic_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Topic Id。",
						},
						"partition_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "数量 partitions。",
						},
						"retention_ms": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Expiration。",
						},
						"note": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "备注",
						},
						"status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "状态，1 in use，2 in deletion。",
						},
					},
				},
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "用于保存结果。",
			},
		},
	}
}

func dataSourceTencentCloudCkafkaDatahubTopicRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_ckafka_datahub_topic.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("search_word"); ok {
		paramMap["search_word"] = v.(string)
	}

	if v, _ := d.GetOk("offset"); v != nil {
		paramMap["offset"] = v.(int)
	}

	if v, _ := d.GetOk("limit"); v != nil {
		paramMap["limit"] = v.(int)
	}

	service := CkafkaService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var describeDatahubTopicsResp *ckafka.DescribeDatahubTopicsResp

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		response, e := service.DescribeCkafkaDatahubTopicByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		describeDatahubTopicsResp = response
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0)
	topicList := make([]map[string]interface{}, 0)

	if describeDatahubTopicsResp != nil {

		if len(describeDatahubTopicsResp.TopicList) != 0 {
			for _, topic := range describeDatahubTopicsResp.TopicList {
				topicMap := map[string]interface{}{}

				if topic.Name != nil {
					topicMap["name"] = topic.Name
				}

				if topic.TopicName != nil {
					topicMap["topic_name"] = topic.TopicName
					ids = append(ids, *topic.TopicName)
				}

				if topic.TopicId != nil {
					topicMap["topic_id"] = topic.TopicId
				}

				if topic.PartitionNum != nil {
					topicMap["partition_num"] = topic.PartitionNum
				}

				if topic.RetentionMs != nil {
					topicMap["retention_ms"] = topic.RetentionMs
				}

				if topic.Note != nil {
					topicMap["note"] = topic.Note
				}

				if topic.Status != nil {
					topicMap["status"] = topic.Status
				}

				topicList = append(topicList, topicMap)
			}

		}
	}
	d.SetId(helper.DataResourceIdsHash(ids))
	_ = d.Set("topic_list", topicList)

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), topicList); e != nil {
			return e
		}
	}
	return nil
}
