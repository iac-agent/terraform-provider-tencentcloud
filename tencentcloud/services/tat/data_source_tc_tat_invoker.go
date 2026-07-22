package tat

import (
	"context"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tat "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tat/v20201028"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudTatInvoker() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTatInvokerRead,
		Schema: map[string]*schema.Schema{
			"invoker_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Invoker ID。",
			},

			"command_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "命令 ID",
			},

			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Invoker 类型",
			},

			"invoker_set": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Invoker 信息。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"invoker_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Invoker ID。",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Invoker 名称",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Invoker 类型",
						},
						"command_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "命令 ID",
						},
						"username": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "用户名",
						},
						"parameters": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Custom 参数。",
						},
						"instance_ids": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Computed:    true,
							Description: "实例 ID 列表。",
						},
						"enable": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "是否enable invoker。",
						},
						"schedule_settings": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Execution 调度 的 invoker. 此 字段 是 返回 对于 recurring invokers。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"policy": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Execution 策略: `ONCE`: Execute once; `RECURRENCE`: Execute repeatedly。",
									},
									"recurrence": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Trigger crontab expression. 此 字段 为必填项 如果 `Policy` 是 `RECURRENCE`. crontab expression 是 parsed 在 UTC+8。",
									},
									"invoke_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "next 执行时间 的 invoker. 此 字段 为必填项 如果 Policy 是 ONCE。",
									},
								},
							},
						},
						"created_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "创建时间。",
						},
						"updated_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "修改时间。",
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

func dataSourceTencentCloudTatInvokerRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_tat_invoker.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("invoker_id"); ok {
		paramMap["invoker_id"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("command_id"); ok {
		paramMap["command_id"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("type"); ok {
		paramMap["type"] = helper.String(v.(string))
	}

	tatService := TatService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var invokerSet []*tat.Invoker
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		results, e := tatService.DescribeTatInvokerByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		invokerSet = results
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read Tat invokerSet failed, reason:%+v", logId, err)
		return err
	}

	invokerSetList := []interface{}{}
	ids := make([]string, 0, len(invokerSet))
	if invokerSet != nil {
		for _, invokerSet := range invokerSet {
			invokerSetMap := map[string]interface{}{}
			if invokerSet.InvokerId != nil {
				invokerSetMap["invoker_id"] = invokerSet.InvokerId
			}
			if invokerSet.Name != nil {
				invokerSetMap["name"] = invokerSet.Name
			}
			if invokerSet.Type != nil {
				invokerSetMap["type"] = invokerSet.Type
			}
			if invokerSet.CommandId != nil {
				invokerSetMap["command_id"] = invokerSet.CommandId
			}
			if invokerSet.Username != nil {
				invokerSetMap["username"] = invokerSet.Username
			}
			if invokerSet.Parameters != nil {
				invokerSetMap["parameters"] = invokerSet.Parameters
			}
			if invokerSet.InstanceIds != nil {
				invokerSetMap["instance_ids"] = invokerSet.InstanceIds
			}
			if invokerSet.Enable != nil {
				invokerSetMap["enable"] = invokerSet.Enable
			}
			if invokerSet.ScheduleSettings != nil {
				scheduleSettingsMap := map[string]interface{}{}
				if invokerSet.ScheduleSettings.Policy != nil {
					scheduleSettingsMap["policy"] = invokerSet.ScheduleSettings.Policy
				}
				if invokerSet.ScheduleSettings.Recurrence != nil {
					scheduleSettingsMap["recurrence"] = invokerSet.ScheduleSettings.Recurrence
				}
				if invokerSet.ScheduleSettings.InvokeTime != nil {
					scheduleSettingsMap["invoke_time"] = invokerSet.ScheduleSettings.InvokeTime
				}

				invokerSetMap["schedule_settings"] = []interface{}{scheduleSettingsMap}
			}
			if invokerSet.CreatedTime != nil {
				invokerSetMap["created_time"] = invokerSet.CreatedTime
			}
			if invokerSet.UpdatedTime != nil {
				invokerSetMap["updated_time"] = invokerSet.UpdatedTime
			}

			invokerSetList = append(invokerSetList, invokerSetMap)
			ids = append(ids, *invokerSet.InvokerId)
		}
		_ = d.Set("invoker_set", invokerSetList)
	}

	d.SetId(helper.DataResourceIdsHash(ids))

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), invokerSetList); e != nil {
			return e
		}
	}

	return nil
}
