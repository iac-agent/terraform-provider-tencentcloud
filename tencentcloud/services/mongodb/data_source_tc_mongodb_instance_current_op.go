package mongodb

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mongodb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mongodb/v20190725"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudMongodbInstanceCurrentOp() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudMongodbInstanceCurrentOpRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "实例 ID， 格式 是: cmgo-9d0p6umb.Same 作为 实例 ID displayed 在 云 数据库 console 页面。",
			},

			"ns": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "过滤器 condition， 命名空间 命名空间 到 其中 operation belongs，在 格式 的 db.collection。",
			},

			"millisecond_running": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "过滤器 condition， 时间 该 operation has been executed (单位: millisecond), 结果 将 返回 operation 该 exceeds 集合 时间， 默认值为 0,和 值 范围 是 [0，3600000]。",
			},

			"op": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "过滤器 condition，操作类型，possible 值: none，update，insert，查询，command，getmore,remove 和 killcursors。",
			},

			"replica_set_name": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "过滤器 condition，分片 名称",
			},

			"state": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "过滤器 condition，节点 状态，possible 值: primary，secondary。",
			},

			"order_by": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "返回sorted 字段 的 结果 集合，currently 支持: MicrosecsRunning/microsecsrunning, 默认为 ascending sort。",
			},

			"order_by_type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "返回sorting 方法 的 结果 集合，possible 值: ASC/asc 或 DESC/desc。",
			},

			"current_ops": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "当前 operation 列表。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"op_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "operation ID。",
						},
						"ns": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "operation 命名空间。",
						},
						"query": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "operation 查询。",
						},
						"op": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "operation 值",
						},
						"replica_set_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Replication 名称",
						},
						"state": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "operation state。",
						},
						"operation": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "operation info。",
						},
						"node_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "节点名称",
						},
						"microsecs_running": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "running 时间(ms)。",
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

func dataSourceTencentCloudMongodbInstanceCurrentOpRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_mongodb_instance_current_op.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		paramMap["instance_id"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("ns"); ok {
		paramMap["ns"] = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("millisecond_running"); ok {
		paramMap["millisecond_running"] = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("op"); ok {
		paramMap["op"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("replica_set_name"); ok {
		paramMap["replica_set_name"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("state"); ok {
		paramMap["state"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("order_by"); ok {
		paramMap["order_by"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("order_by_type"); ok {
		paramMap["order_by_type"] = helper.String(v.(string))
	}

	service := MongodbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var currentOps []*mongodb.CurrentOp

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeMongodbInstanceCurrentOpByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		currentOps = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(currentOps))
	tmpList := make([]map[string]interface{}, 0, len(currentOps))

	if currentOps != nil {
		for _, currentOp := range currentOps {
			currentOpMap := map[string]interface{}{}

			if currentOp.OpId != nil {
				currentOpMap["op_id"] = currentOp.OpId
			}

			if currentOp.Ns != nil {
				currentOpMap["ns"] = currentOp.Ns
			}

			if currentOp.Query != nil {
				currentOpMap["query"] = currentOp.Query
			}

			if currentOp.Op != nil {
				currentOpMap["op"] = currentOp.Op
			}

			if currentOp.ReplicaSetName != nil {
				currentOpMap["replica_set_name"] = currentOp.ReplicaSetName
			}

			if currentOp.State != nil {
				currentOpMap["state"] = currentOp.State
			}

			if currentOp.Operation != nil {
				currentOpMap["operation"] = currentOp.Operation
			}

			if currentOp.NodeName != nil {
				currentOpMap["node_name"] = currentOp.NodeName
			}

			if currentOp.MicrosecsRunning != nil {
				currentOpMap["microsecs_running"] = currentOp.MicrosecsRunning
			}

			ids = append(ids, helper.Int64ToStr(*currentOp.OpId))
			tmpList = append(tmpList, currentOpMap)
		}

		_ = d.Set("current_ops", tmpList)
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
