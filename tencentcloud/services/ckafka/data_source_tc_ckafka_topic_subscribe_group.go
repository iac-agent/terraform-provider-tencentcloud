package ckafka

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ckafka "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ckafka/v20190819"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudCkafkaTopicSubscribeGroup() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCkafkaTopicSubscribeGroupRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "实例 ID",
			},

			"topic_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "TopicName。",
			},

			"groups_info": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Consumer 组 信息。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"error_code": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "错误码，normally 0。",
						},
						"state": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Group state 描述 (commonly Empty，Stable，和 Dead states): Dead: consumption 组 does 不 exist Empty: consumption 组 does 不 currently have any 消费者 subscriptions PreparingRebalance: consumption 组 是 在 rebalance state CompletingRebalance: consumption 组 是 在 rebalance state Stable: Each 消费者 在 consumption 组 has joined 和 是 在 stable state。",
						},
						"protocol_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "协议 类型 selected 通过 consumption 组 是 normally 消费者，但 some systems 使用 their own 协议，such 作为 kafka-connect，其中 uses connect. Only standard 消费者 协议，此 interface knows 格式 的 特定 allocation 方法，和 可以 analyze 特定 分区 allocation。",
						},
						"protocol": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Common 消费者 分区 allocation algorithms 是 作为 follows ( 默认值 选项 对于 Kafka 消费者 SDK 是 范围) 范围|roundrobin| sticky。",
						},
						"members": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "此 数组 包含information 仅 如果 state 是 Stable 和 protocol_type 是 消费者。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"member_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "ID 该 coordinator generated 对于 消费者。",
									},
									"client_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "客户端.ID 信息 集合 通过 客户端 消费者 SDK itself。",
									},
									"client_host": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Generally store customer&#39;s IP 地址",
									},
									"assignment": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Stores 分区 信息 assigned 到 消费者。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"version": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "assignment 版本 信息。",
												},
												"topics": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "主题 列表。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"topic": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "主题 名称",
															},
															"partitions": {
																Type: schema.TypeSet,
																Elem: &schema.Schema{
																	Type: schema.TypeInt,
																},
																Computed:    true,
																Description: "分区 列表。",
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
						"group": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Kafka 消费者 组。",
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

func dataSourceTencentCloudCkafkaTopicSubscribeGroupRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_ckafka_topic_subscribe_group.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		paramMap["InstanceId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("topic_name"); ok {
		paramMap["TopicName"] = helper.String(v.(string))
	}

	service := CkafkaService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var result []*ckafka.GroupInfoResponse

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		groupInfos, e := service.DescribeCkafkaTopicSubscribeGroupByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		result = groupInfos
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(result))

	groupsInfoList := []interface{}{}
	for _, groupsInfo := range result {
		groupsInfoMap := map[string]interface{}{}

		if groupsInfo.ErrorCode != nil {
			groupsInfoMap["error_code"] = groupsInfo.ErrorCode
		}

		if groupsInfo.State != nil {
			groupsInfoMap["state"] = groupsInfo.State
		}

		if groupsInfo.ProtocolType != nil {
			groupsInfoMap["protocol_type"] = groupsInfo.ProtocolType
		}

		if groupsInfo.Protocol != nil {
			groupsInfoMap["protocol"] = groupsInfo.Protocol
		}

		if groupsInfo.Members != nil {
			membersList := []interface{}{}
			for _, members := range groupsInfo.Members {
				membersMap := map[string]interface{}{}

				if members.MemberId != nil {
					membersMap["member_id"] = members.MemberId
				}

				if members.ClientId != nil {
					membersMap["client_id"] = members.ClientId
				}

				if members.ClientHost != nil {
					membersMap["client_host"] = members.ClientHost
				}

				if members.Assignment != nil {
					assignmentMap := map[string]interface{}{}

					if members.Assignment.Version != nil {
						assignmentMap["version"] = members.Assignment.Version
					}

					if members.Assignment.Topics != nil {
						topicsList := []interface{}{}
						for _, topics := range members.Assignment.Topics {
							topicsMap := map[string]interface{}{}

							if topics.Topic != nil {
								topicsMap["topic"] = topics.Topic
							}

							if topics.Partitions != nil {
								topicsMap["partitions"] = topics.Partitions
							}

							topicsList = append(topicsList, topicsMap)
						}

						assignmentMap["topics"] = []interface{}{topicsList}
					}

					membersMap["assignment"] = []interface{}{assignmentMap}
				}

				membersList = append(membersList, membersMap)
			}

			groupsInfoMap["members"] = membersList
		}

		if groupsInfo.Group != nil {
			ids = append(ids, *groupsInfo.Group)

			groupsInfoMap["group"] = groupsInfo.Group
		}

		groupsInfoList = append(groupsInfoList, groupsInfoMap)
	}

	_ = d.Set("groups_info", groupsInfoList)

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), groupsInfoList); e != nil {
			return e
		}
	}
	return nil
}
