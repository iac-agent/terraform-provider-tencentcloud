package tat

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tat "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tat/v20201028"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudTatAgent() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTatAgentRead,
		Schema: map[string]*schema.Schema{
			"instance_ids": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "列表 实例 IDs 对于 查询。",
			},

			"filters": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "过滤器 conditions. agent-状态 - String - 必填: No - (过滤器 condition) 过滤器 通过 agent 状态 有效值：Online，Offline. 环境 - String - 必填: No - (过滤器 condition) 过滤器 通过 agent 环境. 有效 值: Linux. 实例-ID - String - 必填: No - (过滤器 condition) 过滤器 通过 实例 ID. Up 到 10 Filters allowed 在 一个 请求. For each 过滤器，five 过滤器.Values 可以 是 指定. InstanceIds 和 Filters 不能 是 指定 在 same 时间。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "待过滤字段",
						},
						"values": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Required:    true,
							Description: "过滤器 值 的 字段。",
						},
					},
				},
			},

			"automation_agent_set": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "列表 agent 消息",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "实例 ID",
						},
						"version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Agent 版本",
						},
						"last_heartbeat_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Time 的 last heartbeat。",
						},
						"agent_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Agent 状态Ranges:&lt;li&gt; Online:Online&lt;li&gt; Offline:Offline。",
						},
						"environment": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Environment 对于 Agent.Ranges:&lt;li&gt; Linux:Linux 实例&lt;li&gt; Windows:Windows 实例。",
						},
						"support_features": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Computed:    true,
							Description: "列表 功能 Agent support。",
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

func dataSourceTencentCloudTatAgentRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_tat_agent.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_ids"); ok {
		instanceIdsSet := v.(*schema.Set).List()
		paramMap["InstanceIds"] = helper.InterfacesStringsPoint(instanceIdsSet)
	}

	if v, ok := d.GetOk("filters"); ok {
		filtersSet := v.([]interface{})
		tmpSet := make([]*tat.Filter, 0, len(filtersSet))

		for _, item := range filtersSet {
			filter := tat.Filter{}
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

	service := TatService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var automationAgentSet []*tat.AutomationAgentInfo

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeTatAgentByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		automationAgentSet = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(automationAgentSet))
	tmpList := make([]map[string]interface{}, 0, len(automationAgentSet))

	if automationAgentSet != nil {
		for _, automationAgentInfo := range automationAgentSet {
			automationAgentInfoMap := map[string]interface{}{}

			if automationAgentInfo.InstanceId != nil {
				automationAgentInfoMap["instance_id"] = automationAgentInfo.InstanceId
			}

			if automationAgentInfo.Version != nil {
				automationAgentInfoMap["version"] = automationAgentInfo.Version
			}

			if automationAgentInfo.LastHeartbeatTime != nil {
				automationAgentInfoMap["last_heartbeat_time"] = automationAgentInfo.LastHeartbeatTime
			}

			if automationAgentInfo.AgentStatus != nil {
				automationAgentInfoMap["agent_status"] = automationAgentInfo.AgentStatus
			}

			if automationAgentInfo.Environment != nil {
				automationAgentInfoMap["environment"] = automationAgentInfo.Environment
			}

			if automationAgentInfo.SupportFeatures != nil {
				automationAgentInfoMap["support_features"] = automationAgentInfo.SupportFeatures
			}

			ids = append(ids, *automationAgentInfo.InstanceId)
			tmpList = append(tmpList, automationAgentInfoMap)
		}

		_ = d.Set("automation_agent_set", tmpList)
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
