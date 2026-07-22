package cls

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cls "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cls/v20201016"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudClsMachines() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudClsMachinesRead,
		Schema: map[string]*schema.Schema{
			"group_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "组 ID。",
			},

			"machines": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "机器信息。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "机器的ip。",
						},
						"status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "机器的状态。",
						},
						"offline_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "机器离线时间。",
						},
						"auto_update": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "如果打开自动更新标志。",
						},
						"version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "当前机器版本。",
						},
						"update_status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "机器更新状态。",
						},
						"err_code": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "更新操作的代码。",
						},
						"err_msg": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "更新操作的消息。",
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

func dataSourceTencentCloudClsMachinesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_cls_machines.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("group_id"); ok {
		paramMap["GroupId"] = helper.String(v.(string))
	}

	service := ClsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var machines []*cls.MachineInfo

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeClsMachinesByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		machines = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(machines))
	tmpList := make([]map[string]interface{}, 0, len(machines))

	if machines != nil {
		for _, machineInfo := range machines {
			machineInfoMap := map[string]interface{}{}

			if machineInfo.Ip != nil {
				machineInfoMap["ip"] = machineInfo.Ip
			}

			if machineInfo.Status != nil {
				machineInfoMap["status"] = machineInfo.Status
			}

			if machineInfo.OfflineTime != nil {
				machineInfoMap["offline_time"] = machineInfo.OfflineTime
			}

			if machineInfo.AutoUpdate != nil {
				machineInfoMap["auto_update"] = machineInfo.AutoUpdate
			}

			if machineInfo.Version != nil {
				machineInfoMap["version"] = machineInfo.Version
			}

			if machineInfo.UpdateStatus != nil {
				machineInfoMap["update_status"] = machineInfo.UpdateStatus
			}

			if machineInfo.ErrCode != nil {
				machineInfoMap["err_code"] = machineInfo.ErrCode
			}

			if machineInfo.ErrMsg != nil {
				machineInfoMap["err_msg"] = machineInfo.ErrMsg
			}

			ids = append(ids, *machineInfo.Ip)
			tmpList = append(tmpList, machineInfoMap)
		}

		_ = d.Set("machines", tmpList)
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
