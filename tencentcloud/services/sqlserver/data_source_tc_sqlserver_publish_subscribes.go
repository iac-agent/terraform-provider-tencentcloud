package sqlserver

import (
	"context"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudSqlserverPublishSubscribes() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentSqlserverPublishSubscribesRead,
		Schema: map[string]*schema.Schema{
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used 到 store results.",
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID 的 SQL Server 实例.",
			},
			"pub_or_sub_instance_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "subscribe/publish 实例 ID. It 是 related 到 whether `instance_id` 是 publish 实例 或 subscribe 实例. 当 `instance_id` 是 publish 实例, 此 字段 是 filtered according 到 subscribe 实例 ID; 当 `instance_id` 是 subscribe 实例, 此 字段 是 filtering according 到 publish 实例 ID.",
			},
			"pub_or_sub_instance_ip": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "intranet IP 的 subscribe/publish 实例. It 是 related 到 whether `instance_id` 是 publish 实例 或 subscribe 实例. 当 `instance_id` 是 publish 实例, 此 字段 是 filtered according 到 intranet IP 的 subscribe 实例; 当 `instance_id` 是 subscribe 实例, 此 字段 是 based 在 publish 实例 intranet IP 过滤器.",
			},
			"publish_subscribe_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "ID 的 Publish 和 Subscribe.",
			},
			"publish_subscribe_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "名称 的 Publish 和 Subscribe.",
			},
			"publish_database": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name 的 publish 数据库.",
			},
			"subscribe_database": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name 的 subscribe 数据库.",
			},
			"publish_subscribe_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Publish 和 subscribe 列表. Each element contains following attributes.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"publish_subscribe_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "ID 的 Publish 和 Subscribe.",
						},
						"publish_subscribe_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "名称 的 Publish 和 Subscribe.",
						},
						"publish_instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID 的 SQL Server 实例 其中 publish.",
						},
						"publish_instance_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name 的 SQL Server 实例 其中 publish.",
						},
						"publish_instance_ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "IP 的 SQL Server 实例 其中 publish.",
						},
						"subscribe_instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID 的 SQL Server 实例 其中 subscribe.",
						},
						"subscribe_instance_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name 的 SQL Server 实例 其中 subscribe.",
						},
						"subscribe_instance_ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "IP 的 SQL Server 实例 其中 subscribe.",
						},
						"database_tuples": {
							Type:        schema.TypeSet,
							Computed:    true,
							Description: "Database Publish 和 Publish relationship 列表.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"publish_database": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Name 的 publish SQL Server 实例.",
									},
									"subscribe_database": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Name 的 subscribe SQL Server 实例.",
									},
									"last_sync_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Last sync 时间.",
									},
									"status": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Publish 和 subscribe 状态 between databases, 有效 值 是 `running`, `success`, `fail`, `unknow`.",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentSqlserverPublishSubscribesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_sqlserver_publish_subscribes.read")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	sqlserverService := SqlserverService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	paramMap := make(map[string]interface{})
	paramMap["instanceId"] = d.Get("instance_id").(string)
	if v, ok := d.GetOk("pub_or_sub_instance_id"); ok {
		paramMap["pubOrSubInstanceId"] = v.(string)
	}
	if v, ok := d.GetOk("pub_or_sub_instance_ip"); ok {
		paramMap["pubOrSubInstanceIp"] = v.(string)
	}
	if v, ok := d.GetOk("publish_subscribe_name"); ok {
		paramMap["publishSubscribeName"] = v.(string)
	}
	if v, ok := d.GetOk("publish_subscribe_id"); ok {
		id := v.(int)
		paramMap["publishSubscribeId"] = *helper.IntUint64(id)
	} else {
		paramMap["publishSubscribeId"] = *helper.IntUint64(0)
	}
	if v, ok := d.GetOk("publish_database"); ok {
		paramMap["publishDBName"] = v.(string)
		paramMap["subscribeDBName"] = v.(string)
	}

	publishSubscribes, err := sqlserverService.DescribeSqlserverPublishSubscribes(ctx, paramMap)
	if err != nil {
		return err
	}
	instanceList := make([]map[string]interface{}, 0, len(publishSubscribes))
	ids := make([]string, 0, len(publishSubscribes))

	for _, publishSubscribe := range publishSubscribes {
		var databaseTupleStatus []map[string]interface{}
		for _, inst := range publishSubscribe.DatabaseTupleSet {
			databaseTuple := map[string]interface{}{
				"publish_database":   inst.PublishDatabase,
				"subscribe_database": inst.SubscribeDatabase,
				"last_sync_time":     inst.LastSyncTime,
				"status":             inst.Status,
			}
			databaseTupleStatus = append(databaseTupleStatus, databaseTuple)
		}
		instance := map[string]interface{}{
			"publish_subscribe_id":    publishSubscribe.Id,
			"publish_subscribe_name":  publishSubscribe.Name,
			"publish_instance_id":     publishSubscribe.PublishInstanceId,
			"publish_instance_ip":     publishSubscribe.PublishInstanceIp,
			"publish_instance_name":   publishSubscribe.PublishInstanceName,
			"subscribe_instance_id":   publishSubscribe.SubscribeInstanceId,
			"subscribe_instance_ip":   publishSubscribe.SubscribeInstanceIp,
			"subscribe_instance_name": publishSubscribe.SubscribeInstanceName,
			"database_tuples":         databaseTupleStatus,
		}
		resourceId := *publishSubscribe.PublishInstanceId + tccommon.FILED_SP + *publishSubscribe.SubscribeInstanceId
		instanceList = append(instanceList, instance)
		ids = append(ids, resourceId)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	if err = d.Set("publish_subscribe_list", instanceList); err != nil {
		log.Printf("[CRITAL]%s provider set sql server publish and subscribe list fail, reason:%s ", logId, err.Error())
		return err
	}

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if err := tccommon.WriteToFile(output.(string), instanceList); err != nil {
			return err
		}
	}
	return nil
}
