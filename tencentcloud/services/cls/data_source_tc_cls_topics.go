package cls

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	clsv20201016 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cls/v20201016"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudClsTopics() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudClsTopicsRead,
		Schema: map[string]*schema.Schema{
			"filters": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "<li>topicName：按**日志主题名称**过滤。默认实现模糊匹配。您可以使用“PreciseSearch”参数来设置精确匹配。类型：字符串。必需的。否。<br><li>logsetName：按**日志集名称**过滤。默认实现模糊匹配。您可以使用“PreciseSearch”参数来设置精确匹配。类型：字符串。必填：否。<br><li>topicId：按**日志主题 ID** 过滤。类型：字符串。必需：否。<br><li>logsetId：按 **日志集 ID** 过滤。您可以调用“DescribeLogsets”查询已创建的日志集列表或登录控制台查看。您还可以调用 CreateLogset 创建日志集。类型：字符串。必填：否。 <br><li>tagKey：按**标签键**过滤。类型：字符串。必填：否。<br><li>标签:tagKey：按**标签键值对**过滤。 `tagKey` 应替换为指定的标签键，例如 `标签:exampleKey`。类型：字符串。必需：否。<br><li>存储类型：按**日志主题存储类型**过滤。有效值：“热”（标准存储）和“冷”（IA 存储）。类型：字符串。必需：否。每个请求最多可以有 10 个“Filters”和 100 个“过滤器.Values”。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "要过滤的字段。",
						},
						"values": {
							Type:        schema.TypeSet,
							Required:    true,
							Description: "要过滤的值。",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},

			"precise_search": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "`Filters` 字段的匹配模式。\n- 0: `topicName` 和 `logsetName` 的模糊匹配。这是默认值。\n- 1：“topicName”精确匹配。\n- 2：“logsetName”精确匹配。\n- 3：“topicName”和“logsetName”精确匹配。",
			},

			"biz_type": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "主题类型\n- 0（默认）：日志主题。\n- 1：指标主题。",
			},

			"topics": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "记录主题列表。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"logset_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "日志集 ID。",
						},
						"topic_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "主题 ID。",
						},
						"topic_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "主题名称。",
						},
						"partition_count": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "主题分区的数量。",
						},
						"index": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "主题是否启用索引（主题类型必须是日志主题）。",
						},
						"assumer_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "云产品标识符。当主题是由其他云产品创建时，该字段显示云产品的名称，例如CDN、TKE。注意：该字段可能返回null，表示取不到有效值。",
						},
						"create_time": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "创作时间。",
						},
						"status": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "主题是否启用日志收集。 true：启用收集； false：禁用采集。创建日志主题时默认开启日志采集，可以通过SDK调用ModifyTopic修改该字段。控制台目前不支持修改该参数。",
						},
						"tags": {
							Type:        schema.TypeList,
							Required:    true,
							Description: "主题绑定的标签信息注意：该字段可能返回null，表示取不到有效值。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "标签键。\n注意：该字段可能返回null，表示取不到有效值。",
									},
									"value": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "标签值。\n注意：该字段可能返回null，表示取不到有效值。",
									},
								},
							},
						},
						"auto_split": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "该主题是否启用自动拆分\n注意：该字段可能返回“null”，表示取不到有效值。",
						},
						"max_split_partitions": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "如果启用自动拆分，该主题拆分的最大分区数\n注意：该字段可能返回“null”，表示无法获取有效值。",
						},
						"storage_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "topic的存储类型注意：该字段可能返回null，表示取不到有效值。",
						},
						"period": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "生命周期以天为单位。取值范围：1-3600（3640表示永久保留）\n注意：该字段可能返回“null”，表示未找到有效值。",
						},
						"sub_assumer_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "云产品子标识符。如果日志主题是由其他云产品创建的，则此字段返回云产品的名称及其日志类型，例如“TKE-Audit”或“TKE-Event”。部分产品仅返回云产品标识（`AssumerName`），没有此字段。\n注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"describes": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "主题描述\n注意：该字段可能返回null，表示取不到有效值。",
						},
						"hot_period": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "启用日志下沉，采用标准存储的生命周期，其中hotPeriod < 周期。对于标准存储，使用hotPeriod，对于不频繁访问的存储，使用Period-hotPeriod。 （主题类型必须是日志主题）HotPeriod=0表示不开启日志下沉。\n注意：该字段可能返回null，表示取不到有效值。",
						},
						"biz_type": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Topic类型。\n- 0：日志 Topic \n- 1：Metric Topic\n注意：该字段可能返回null，表示取不到有效值。",
						},
						"is_web_tracking": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "免费认证开关。 false：禁用； true：启用。启用后，指定操作将支持匿名访问日志主题。具体请参见日志主题（https://intl.云.tencent.com/document/product/614/41035?from_cn_redirect=1）。注意：该字段可能返回null，表示取不到有效值。",
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

func dataSourceTencentCloudClsTopicsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_cls_topics.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = ClsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("filters"); ok {
		filtersSet := v.([]interface{})
		tmpSet := make([]*clsv20201016.Filter, 0, len(filtersSet))
		for _, item := range filtersSet {
			filtersMap := item.(map[string]interface{})
			filter := clsv20201016.Filter{}
			if v, ok := filtersMap["key"].(string); ok && v != "" {
				filter.Key = helper.String(v)
			}

			if v, ok := filtersMap["values"]; ok {
				valuesSet := v.(*schema.Set).List()
				for i := range valuesSet {
					values := valuesSet[i].(string)
					filter.Values = append(filter.Values, helper.String(values))
				}
			}

			tmpSet = append(tmpSet, &filter)
		}

		paramMap["Filters"] = tmpSet
	}

	if v, ok := d.GetOkExists("precise_search"); ok {
		paramMap["PreciseSearch"] = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("biz_type"); ok {
		paramMap["BizType"] = helper.IntUint64(v.(int))
	}

	var respData []*clsv20201016.TopicInfo
	reqErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeClsTopicsByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		respData = result
		return nil
	})

	if reqErr != nil {
		return reqErr
	}

	ids := make([]string, 0, len(respData))
	topicsList := make([]map[string]interface{}, 0, len(respData))
	if respData != nil {
		for _, topics := range respData {
			topicsMap := map[string]interface{}{}
			if topics.LogsetId != nil {
				topicsMap["logset_id"] = topics.LogsetId
			}

			if topics.TopicId != nil {
				topicsMap["topic_id"] = topics.TopicId
				ids = append(ids, *topics.TopicId)
			}

			if topics.TopicName != nil {
				topicsMap["topic_name"] = topics.TopicName
			}

			if topics.PartitionCount != nil {
				topicsMap["partition_count"] = topics.PartitionCount
			}

			if topics.Index != nil {
				topicsMap["index"] = topics.Index
			}

			if topics.AssumerName != nil {
				topicsMap["assumer_name"] = topics.AssumerName
			}

			if topics.CreateTime != nil {
				topicsMap["create_time"] = topics.CreateTime
			}

			if topics.Status != nil {
				topicsMap["status"] = topics.Status
			}

			tagsList := make([]map[string]interface{}, 0, len(topics.Tags))
			if topics.Tags != nil {
				for _, tags := range topics.Tags {
					tagsMap := map[string]interface{}{}
					if tags.Key != nil {
						tagsMap["key"] = tags.Key
					}

					if tags.Value != nil {
						tagsMap["value"] = tags.Value
					}

					tagsList = append(tagsList, tagsMap)
				}

				topicsMap["tags"] = tagsList
			}

			if topics.AutoSplit != nil {
				topicsMap["auto_split"] = topics.AutoSplit
			}

			if topics.MaxSplitPartitions != nil {
				topicsMap["max_split_partitions"] = topics.MaxSplitPartitions
			}

			if topics.StorageType != nil {
				topicsMap["storage_type"] = topics.StorageType
			}

			if topics.Period != nil {
				topicsMap["period"] = topics.Period
			}

			if topics.SubAssumerName != nil {
				topicsMap["sub_assumer_name"] = topics.SubAssumerName
			}

			if topics.Describes != nil {
				topicsMap["describes"] = topics.Describes
			}

			if topics.HotPeriod != nil {
				topicsMap["hot_period"] = topics.HotPeriod
			}

			if topics.BizType != nil {
				topicsMap["biz_type"] = topics.BizType
			}

			if topics.IsWebTracking != nil {
				topicsMap["is_web_tracking"] = topics.IsWebTracking
			}

			topicsList = append(topicsList, topicsMap)
		}

		_ = d.Set("topics", topicsList)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), topicsList); e != nil {
			return e
		}
	}

	return nil
}
