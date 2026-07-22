package dlc

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dlc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dlc/v20210125"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudDlcDescribeUpdatableDataEngines() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudDlcDescribeUpdatableDataEnginesRead,
		Schema: map[string]*schema.Schema{
			"data_engine_config_command": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Operation commands 的 引擎 配置. UpdateSparkSQLLakefsPath updates 路径 的 managed tables，和 UpdateSparkSQLResultPath updates 路径 的 结果 buckets。",
			},

			"data_engine_basic_infos": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Basic 集群 信息。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"data_engine_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "DataEngine 名称",
						},
						"state": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "EData 引擎 状态: -2: 删除; -1: failed; 0: initializing; 1: suspended; 2: running; 3: ready 到 delete; 4: deleting。",
						},
						"create_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "创建时间。",
						},
						"update_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "更新时间。",
						},
						"message": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Returned 信息。",
						},
						"data_engine_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Engine ID。",
						},
						"data_engine_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Engine types，和 有效 值 是 PrestoSQL，SparkSQL，和 SparkBatch。",
						},
						"app_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "用户 ID。",
						},
						"user_uin": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "账号 uin。",
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

func dataSourceTencentCloudDlcDescribeUpdatableDataEnginesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_dlc_describe_updatable_data_engines.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("data_engine_config_command"); ok {
		paramMap["DataEngineConfigCommand"] = helper.String(v.(string))
	}

	service := DlcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var dataEngineBasicInfos []*dlc.DataEngineBasicInfo

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeDlcDescribeUpdatableDataEnginesByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		dataEngineBasicInfos = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(dataEngineBasicInfos))
	tmpList := make([]map[string]interface{}, 0, len(dataEngineBasicInfos))

	if dataEngineBasicInfos != nil {
		for _, dataEngineBasicInfo := range dataEngineBasicInfos {
			dataEngineBasicInfoMap := map[string]interface{}{}

			if dataEngineBasicInfo.DataEngineName != nil {
				dataEngineBasicInfoMap["data_engine_name"] = dataEngineBasicInfo.DataEngineName
			}

			if dataEngineBasicInfo.State != nil {
				dataEngineBasicInfoMap["state"] = dataEngineBasicInfo.State
			}

			if dataEngineBasicInfo.CreateTime != nil {
				dataEngineBasicInfoMap["create_time"] = dataEngineBasicInfo.CreateTime
			}

			if dataEngineBasicInfo.UpdateTime != nil {
				dataEngineBasicInfoMap["update_time"] = dataEngineBasicInfo.UpdateTime
			}

			if dataEngineBasicInfo.Message != nil {
				dataEngineBasicInfoMap["message"] = dataEngineBasicInfo.Message
			}

			if dataEngineBasicInfo.DataEngineId != nil {
				dataEngineBasicInfoMap["data_engine_id"] = dataEngineBasicInfo.DataEngineId
			}

			if dataEngineBasicInfo.DataEngineType != nil {
				dataEngineBasicInfoMap["data_engine_type"] = dataEngineBasicInfo.DataEngineType
			}

			if dataEngineBasicInfo.AppId != nil {
				dataEngineBasicInfoMap["app_id"] = dataEngineBasicInfo.AppId
			}

			if dataEngineBasicInfo.UserUin != nil {
				dataEngineBasicInfoMap["user_uin"] = dataEngineBasicInfo.UserUin
			}

			ids = append(ids, *dataEngineBasicInfo.DataEngineId)
			tmpList = append(tmpList, dataEngineBasicInfoMap)
		}

		_ = d.Set("data_engine_basic_infos", tmpList)
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
