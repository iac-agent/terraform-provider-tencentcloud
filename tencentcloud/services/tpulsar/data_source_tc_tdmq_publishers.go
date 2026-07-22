package tpulsar

import (
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svctdmq "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tdmq"

	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tdmq "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tdmq/v20200217"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudTdmqPublishers() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTdmqPublishersRead,
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "集群 ID",
			},
			"namespace": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "命名空间 名称",
			},
			"topic": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "主题 名称",
			},
			"filters": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "Parameter 过滤器，support ProducerName，地址 字段。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "名称 过滤器 参数。",
						},
						"values": {
							Type:        schema.TypeSet,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Optional:    true,
							Description: "值",
						},
					},
				},
			},
			"sort": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "sorter。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "sorter。",
						},
						"order": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Ascending ASC，descending DESC。",
						},
					},
				},
			},
			// computed
			"publishers": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Producer Information List注意：此字段可能返回 null，表示无法获取有效值。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"producer_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "producer id注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"producer_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "producer name注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"address": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "producer address注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"client_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "客户端 version注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"msg_rate_in": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "消息 production 速率 (articles/second)注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"msg_throughput_in": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "消息 production 吞吐量 速率 (bytes/second)注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"average_msg_size": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "Average 消息 大小 (bytes)注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"connected_since": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "连接 time注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"partition": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "主题 分区 数量 producer connection注意：此字段可能返回 null，表示无法获取有效值。",
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

func dataSourceTencentCloudTdmqPublishersRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_tdmq_publishers.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service    = svctdmq.NewTdmqService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
		publishers []*tdmq.Publisher
		clusterId  string
		Namespace  string
		Topic      string
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("cluster_id"); ok {
		paramMap["ClusterId"] = helper.String(v.(string))
		clusterId = v.(string)
	}

	if v, ok := d.GetOk("namespace"); ok {
		paramMap["Namespace"] = helper.String(v.(string))
		Namespace = v.(string)
	}

	if v, ok := d.GetOk("topic"); ok {
		paramMap["Topic"] = helper.String(v.(string))
		Topic = v.(string)
	}

	if v, ok := d.GetOk("filters"); ok {
		filtersSet := v.([]interface{})
		tmpSet := make([]*tdmq.Filter, 0, len(filtersSet))
		for _, item := range filtersSet {
			filter := tdmq.Filter{}
			filterMap := item.(map[string]interface{})

			if v, ok := filterMap["name"]; ok {
				filter.Name = helper.String(v.(string))
			}
			if v, ok := filterMap["values"]; ok {
				valuesSet := v.(*schema.Set).List()
				filter.Values = helper.InterfacesStringsPoint(valuesSet)
			}
			tmpSet = append(tmpSet, &filter)
		}
		paramMap["filters"] = tmpSet
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "sort"); ok {
		sort := tdmq.Sort{}
		if v, ok := dMap["name"]; ok {
			sort.Name = helper.String(v.(string))
		}
		if v, ok := dMap["order"]; ok {
			sort.Order = helper.String(v.(string))
		}
		paramMap["sort"] = &sort
	}

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeTdmqPublishersByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		publishers = result
		return nil
	})

	if err != nil {
		return err
	}

	ids := make([]string, 0)
	tmpList := make([]map[string]interface{}, 0, len(publishers))

	if publishers != nil {
		for _, publisher := range publishers {
			publisherMap := map[string]interface{}{}

			if publisher.ProducerId != nil {
				publisherMap["producer_id"] = publisher.ProducerId
			}

			if publisher.ProducerName != nil {
				publisherMap["producer_name"] = publisher.ProducerName
			}

			if publisher.Address != nil {
				publisherMap["address"] = publisher.Address
			}

			if publisher.ClientVersion != nil {
				publisherMap["client_version"] = publisher.ClientVersion
			}

			if publisher.MsgRateIn != nil {
				publisherMap["msg_rate_in"] = publisher.MsgRateIn
			}

			if publisher.MsgThroughputIn != nil {
				publisherMap["msg_throughput_in"] = publisher.MsgThroughputIn
			}

			if publisher.AverageMsgSize != nil {
				publisherMap["average_msg_size"] = publisher.AverageMsgSize
			}

			if publisher.ConnectedSince != nil {
				publisherMap["connected_since"] = publisher.ConnectedSince
			}

			if publisher.Partition != nil {
				publisherMap["partition"] = publisher.Partition
			}

			tmpList = append(tmpList, publisherMap)
		}

		_ = d.Set("publishers", tmpList)
	}

	ids = append(ids, clusterId)
	ids = append(ids, Namespace)
	ids = append(ids, Topic)
	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), tmpList); e != nil {
			return e
		}
	}

	return nil
}
