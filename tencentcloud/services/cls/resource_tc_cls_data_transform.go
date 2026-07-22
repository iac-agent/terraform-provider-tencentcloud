package cls

import (
	"context"
	"fmt"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cls "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cls/v20201016"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudClsDataTransform() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudClsDataTransformCreate,
		Read:   resourceTencentCloudClsDataTransformRead,
		Update: resourceTencentCloudClsDataTransformUpdate,
		Delete: resourceTencentCloudClsDataTransformDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"func_type": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "任务类型。 `1`：指定主题； `2`：动态创建。",
			},

			"src_topic_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "源主题 ID。",
			},

			"name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "任务名称。",
			},

			"etl_content": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "数据转换内容。如果“func_type”为“2”，则必须使用“log_auto_output”。",
			},

			"task_type": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "任务类型。 `1`：使用源日志主题中的随机数据进行处理预览； `2`：使用用户自定义的测试数据进行处理预览； `3`：创建真实的加工任务。",
			},

			"enable_flag": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "任务启用标志。 `1`：启用，`2`：禁用，默认为`1`。",
			},

			"dst_resources": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "数据转换资源。如果“func_type”为“1”，则此参数是必需的。如果“func_type”为“2”，则该参数无需填写。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"topic_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "目标主题 ID。",
						},
						"alias": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "别名。",
						},
					},
				},
			},

			"backup_give_up_data": {
				Optional:    true,
				Type:        schema.TypeBool,
				Description: "当 func_type 为 2 时，动态创建的日志集和主题数量超过产品规格限制时是否丢弃数据。默认为“假”。 `false`：创建备份日志集和主题，并将日志写入备份主题； `true`：丢弃日志数据。",
			},

			"has_services_log": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "是否开启服务日志下发。 `1`：禁用； `2`：启用。",
			},

			"data_transform_type": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "数据转换类型。 `0`：标准数据转换任务； `1`：预处理数据转换任务（在写入日志主题之前处理收集的日志）。",
			},

			"keep_failure_log": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "保留故障日志状态。 `1`：不保留（默认）； `2`：保留。",
			},

			"failure_log_key": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeString,
				Description: "失败日志的字段名称。",
			},

			"process_from_timestamp": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "指定处理数据的开始时间，单位为秒级时间戳。日志主题生命周期内的任意时间范围。如果超过生命周期，则只处理生命周期内数据的部分。",
			},

			"process_to_timestamp": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "指定处理数据的结束时间，单位为秒级时间戳。无法指定未来时间。如果不填，则表示继续执行。",
			},

			"data_transform_sql_data_sources": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "关联数据源信息。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"data_source": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "数据源类型。 `1`：MySQL； `2`：自建MySQL； `3`：PostgreSQL。",
						},
						"region": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "实例 ID 区域。例如：ap-广州。",
						},
						"instance_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "实例ID。当DataSource为‘1’时，代表云数据库MySQL实例ID，如：cdb-zxcvbnm。",
						},
						"user": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "MySQL 访问用户名。",
						},
						"alias_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "别名。用于数据转换语句。",
						},
						"password": {
							Type:        schema.TypeString,
							Required:    true,
							Sensitive:   true,
							Description: "MySQL访问密码。",
						},
					},
				},
			},

			"env_infos": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "设置环境变量。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "环境变量名称。",
						},
						"value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "环境变量值。",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudClsDataTransformCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cls_data_transform.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId    = tccommon.GetLogId(tccommon.ContextNil)
		request  = cls.NewCreateDataTransformRequest()
		response = cls.NewCreateDataTransformResponse()
		taskId   string
	)

	if v, ok := d.GetOkExists("func_type"); ok {
		request.FuncType = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("src_topic_id"); ok {
		request.SrcTopicId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}

	if v, ok := d.GetOk("etl_content"); ok {
		request.EtlContent = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("task_type"); ok {
		request.TaskType = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("enable_flag"); ok {
		request.EnableFlag = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("dst_resources"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			dataTransformResouceInfo := cls.DataTransformResouceInfo{}
			if v, ok := dMap["topic_id"]; ok {
				dataTransformResouceInfo.TopicId = helper.String(v.(string))
			}

			if v, ok := dMap["alias"]; ok {
				dataTransformResouceInfo.Alias = helper.String(v.(string))
			}

			request.DstResources = append(request.DstResources, &dataTransformResouceInfo)
		}
	}

	if v, ok := d.GetOkExists("backup_give_up_data"); ok {
		request.BackupGiveUpData = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOkExists("has_services_log"); ok {
		request.HasServicesLog = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("data_transform_type"); ok {
		request.DataTransformType = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("keep_failure_log"); ok {
		request.KeepFailureLog = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("failure_log_key"); ok {
		request.FailureLogKey = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("process_from_timestamp"); ok {
		request.ProcessFromTimestamp = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("process_to_timestamp"); ok {
		request.ProcessToTimestamp = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("data_transform_sql_data_sources"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			dataTransformSqlDataSource := cls.DataTransformSqlDataSource{}
			if v, ok := dMap["data_source"]; ok {
				dataTransformSqlDataSource.DataSource = helper.IntUint64(v.(int))
			}

			if v, ok := dMap["region"]; ok {
				dataTransformSqlDataSource.Region = helper.String(v.(string))
			}

			if v, ok := dMap["instance_id"]; ok {
				dataTransformSqlDataSource.InstanceId = helper.String(v.(string))
			}

			if v, ok := dMap["user"]; ok {
				dataTransformSqlDataSource.User = helper.String(v.(string))
			}

			if v, ok := dMap["alias_name"]; ok {
				dataTransformSqlDataSource.AliasName = helper.String(v.(string))
			}

			if v, ok := dMap["password"]; ok {
				dataTransformSqlDataSource.Password = helper.String(v.(string))
			}

			request.DataTransformSqlDataSources = append(request.DataTransformSqlDataSources, &dataTransformSqlDataSource)
		}
	}

	if v, ok := d.GetOk("env_infos"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			envInfo := cls.EnvInfo{}
			if v, ok := dMap["key"]; ok {
				envInfo.Key = helper.String(v.(string))
			}

			if v, ok := dMap["value"]; ok {
				envInfo.Value = helper.String(v.(string))
			}

			request.EnvInfos = append(request.EnvInfos, &envInfo)
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClsClient().CreateDataTransform(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create cls dataTransform failed, reason:%+v", logId, err)
		return err
	}

	taskId = *response.Response.TaskId
	d.SetId(taskId)

	return resourceTencentCloudClsDataTransformRead(d, meta)
}

func resourceTencentCloudClsDataTransformRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cls_data_transform.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId               = tccommon.GetLogId(tccommon.ContextNil)
		ctx                 = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service             = ClsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		dataTransformTaskId = d.Id()
	)

	dataTransform, err := service.DescribeClsDataTransformById(ctx, dataTransformTaskId)
	if err != nil {
		return err
	}

	if dataTransform == nil {
		log.Printf("[WARN]%s resource `tencentcloud_cls_data_transform` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	if dataTransform.SrcTopicId != nil {
		_ = d.Set("src_topic_id", dataTransform.SrcTopicId)
	}

	if dataTransform.Name != nil {
		_ = d.Set("name", dataTransform.Name)
	}

	if dataTransform.EtlContent != nil {
		_ = d.Set("etl_content", dataTransform.EtlContent)
	}

	if dataTransform.EnableFlag != nil {
		_ = d.Set("enable_flag", dataTransform.EnableFlag)
	}

	if dataTransform.DstResources != nil {
		var dstResourcesList []interface{}
		for _, dstResources := range dataTransform.DstResources {
			dstResourcesMap := map[string]interface{}{}

			if dstResources.TopicId != nil {
				dstResourcesMap["topic_id"] = dstResources.TopicId
			}

			if dstResources.Alias != nil {
				dstResourcesMap["alias"] = dstResources.Alias
			}

			dstResourcesList = append(dstResourcesList, dstResourcesMap)
		}

		_ = d.Set("dst_resources", dstResourcesList)
	}

	if dataTransform.BackupGiveUpData != nil {
		_ = d.Set("backup_give_up_data", dataTransform.BackupGiveUpData)
	}

	if dataTransform.HasServicesLog != nil {
		_ = d.Set("has_services_log", dataTransform.HasServicesLog)
	}

	if dataTransform.DataTransformType != nil {
		_ = d.Set("data_transform_type", dataTransform.DataTransformType)
	}

	if dataTransform.KeepFailureLog != nil {
		_ = d.Set("keep_failure_log", dataTransform.KeepFailureLog)
	}

	if dataTransform.FailureLogKey != nil {
		_ = d.Set("failure_log_key", dataTransform.FailureLogKey)
	}

	if dataTransform.ProcessFromTimestamp != nil {
		_ = d.Set("process_from_timestamp", dataTransform.ProcessFromTimestamp)
	}

	if dataTransform.ProcessToTimestamp != nil {
		_ = d.Set("process_to_timestamp", dataTransform.ProcessToTimestamp)
	}

	if dataTransform.DataTransformSqlDataSources != nil {
		var dataSourcesList []interface{}
		for _, dataSource := range dataTransform.DataTransformSqlDataSources {
			dataSourceMap := map[string]interface{}{}

			if dataSource.DataSource != nil {
				dataSourceMap["data_source"] = dataSource.DataSource
			}

			if dataSource.Region != nil {
				dataSourceMap["region"] = dataSource.Region
			}

			if dataSource.InstanceId != nil {
				dataSourceMap["instance_id"] = dataSource.InstanceId
			}

			if dataSource.User != nil {
				dataSourceMap["user"] = dataSource.User
			}

			if dataSource.AliasName != nil {
				dataSourceMap["alias_name"] = dataSource.AliasName
			}

			if dataSource.Password != nil {
				dataSourceMap["password"] = dataSource.Password
			}

			dataSourcesList = append(dataSourcesList, dataSourceMap)
		}

		_ = d.Set("data_transform_sql_data_sources", dataSourcesList)
	}

	if dataTransform.EnvInfos != nil {
		var envInfosList []interface{}
		for _, envInfo := range dataTransform.EnvInfos {
			envInfoMap := map[string]interface{}{}

			if envInfo.Key != nil {
				envInfoMap["key"] = envInfo.Key
			}

			if envInfo.Value != nil {
				envInfoMap["value"] = envInfo.Value
			}

			envInfosList = append(envInfosList, envInfoMap)
		}

		_ = d.Set("env_infos", envInfosList)
	}

	return nil
}

func resourceTencentCloudClsDataTransformUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cls_data_transform.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId               = tccommon.GetLogId(tccommon.ContextNil)
		request             = cls.NewModifyDataTransformRequest()
		dataTransformTaskId = d.Id()
	)

	immutableArgs := []string{"src_topic_id", "preview_log_statistics", "data_transform_type", "process_from_timestamp", "process_to_timestamp"}
	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	request.TaskId = &dataTransformTaskId

	if d.HasChange("name") {
		if v, ok := d.GetOk("name"); ok {
			request.Name = helper.String(v.(string))
		}
	}

	if d.HasChange("etl_content") {
		if v, ok := d.GetOk("etl_content"); ok {
			request.EtlContent = helper.String(v.(string))
		}
	}

	if d.HasChange("enable_flag") {
		if v, ok := d.GetOkExists("enable_flag"); ok {
			request.EnableFlag = helper.IntInt64(v.(int))
		}
	}

	if d.HasChange("dst_resources") {
		if v, ok := d.GetOk("dst_resources"); ok {
			for _, item := range v.([]interface{}) {
				dataTransformResouceInfo := cls.DataTransformResouceInfo{}
				dMap := item.(map[string]interface{})
				if v, ok := dMap["topic_id"]; ok {
					dataTransformResouceInfo.TopicId = helper.String(v.(string))
				}

				if v, ok := dMap["alias"]; ok {
					dataTransformResouceInfo.Alias = helper.String(v.(string))
				}

				request.DstResources = append(request.DstResources, &dataTransformResouceInfo)
			}
		}
	}

	if d.HasChange("backup_give_up_data") {
		if v, ok := d.GetOkExists("backup_give_up_data"); ok {
			request.BackupGiveUpData = helper.Bool(v.(bool))
		}
	}

	if d.HasChange("has_services_log") {
		if v, ok := d.GetOkExists("has_services_log"); ok {
			request.HasServicesLog = helper.IntUint64(v.(int))
		}
	}

	if d.HasChange("keep_failure_log") {
		if v, ok := d.GetOkExists("keep_failure_log"); ok {
			request.KeepFailureLog = helper.IntUint64(v.(int))
		}
	}

	if d.HasChange("failure_log_key") {
		if v, ok := d.GetOk("failure_log_key"); ok {
			request.FailureLogKey = helper.String(v.(string))
		}
	}

	if d.HasChange("data_transform_sql_data_sources") {
		if v, ok := d.GetOk("data_transform_sql_data_sources"); ok {
			for _, item := range v.([]interface{}) {
				dMap := item.(map[string]interface{})
				dataTransformSqlDataSource := cls.DataTransformSqlDataSource{}
				if v, ok := dMap["data_source"]; ok {
					dataTransformSqlDataSource.DataSource = helper.IntUint64(v.(int))
				}

				if v, ok := dMap["region"]; ok {
					dataTransformSqlDataSource.Region = helper.String(v.(string))
				}

				if v, ok := dMap["instance_id"]; ok {
					dataTransformSqlDataSource.InstanceId = helper.String(v.(string))
				}

				if v, ok := dMap["user"]; ok {
					dataTransformSqlDataSource.User = helper.String(v.(string))
				}

				if v, ok := dMap["alias_name"]; ok {
					dataTransformSqlDataSource.AliasName = helper.String(v.(string))
				}

				if v, ok := dMap["password"]; ok {
					dataTransformSqlDataSource.Password = helper.String(v.(string))
				}

				request.DataTransformSqlDataSources = append(request.DataTransformSqlDataSources, &dataTransformSqlDataSource)
			}
		}
	}

	if d.HasChange("env_infos") {
		if v, ok := d.GetOk("env_infos"); ok {
			for _, item := range v.([]interface{}) {
				dMap := item.(map[string]interface{})
				envInfo := cls.EnvInfo{}
				if v, ok := dMap["key"]; ok {
					envInfo.Key = helper.String(v.(string))
				}

				if v, ok := dMap["value"]; ok {
					envInfo.Value = helper.String(v.(string))
				}

				request.EnvInfos = append(request.EnvInfos, &envInfo)
			}
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClsClient().ModifyDataTransform(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s update cls dataTransform failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudClsDataTransformRead(d, meta)
}

func resourceTencentCloudClsDataTransformDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cls_data_transform.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId               = tccommon.GetLogId(tccommon.ContextNil)
		ctx                 = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service             = ClsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		dataTransformTaskId = d.Id()
	)

	if err := service.DeleteClsDataTransformById(ctx, dataTransformTaskId); err != nil {
		return err
	}

	return nil
}
