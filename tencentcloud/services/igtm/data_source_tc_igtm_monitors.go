package igtm

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	igtmv20231024 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/igtm/v20231024"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudIgtmMonitors() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudIgtmMonitorsRead,
		Schema: map[string]*schema.Schema{
			"filters": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Query 过滤器 conditions。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "过滤字段名称，支持 MonitorName: 监控 名称; MonitorId: 监控 ID。",
						},
						"value": {
							Type:        schema.TypeSet,
							Required:    true,
							Description: "过滤器 字段 值。",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"fuzzy": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "是否enable fuzzy 查询，仅 支持 过滤字段名称 作为 域名\nWhen fuzzy 查询 是 已启用，值 最大 长度 是 1，otherwise 值 最大 长度 是 5. (Reserved 字段，currently unused)。",
						},
					},
				},
			},

			"is_detect_num": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "是否query detection count，0 对于 无，1 对于 yes。",
			},

			"monitor_data_set": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Monitor 列表。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"monitor_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Detection 规则 ID。",
						},
						"monitor_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Monitor 名称",
						},
						"uin": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "所有者 用户",
						},
						"detector_group_ids": {
							Type:        schema.TypeSet,
							Computed:    true,
							Description: "Monitoring 节点 ID 组。",
							Elem: &schema.Schema{
								Type: schema.TypeInt,
							},
						},
						"check_protocol": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Detection 协议 PING TCP HTTP HTTPS。",
						},
						"check_interval": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Detection 周期",
						},
						"ping_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Packet count。",
						},
						"tcp_port": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "TCP 端口",
						},
						"host": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Detection 主机",
						},
						"path": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Detection 路径",
						},
						"return_code_threshold": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "返回值 阈值。",
						},
						"enable_redirect": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "是否enable 3xx redirect following ENABLED DISABLED。",
						},
						"enable_sni": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "是否enable SNI.\nENABLED DISABLED。",
						},
						"packet_loss_rate": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Packet loss 速率 upper 限制",
						},
						"timeout": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Detection 超时。",
						},
						"fail_times": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Failure count。",
						},
						"fail_rate": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Failure 速率 upper 限制 100。",
						},
						"created_on": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "创建时间。",
						},
						"updated_on": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "更新时间。",
						},
						"detector_style": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Monitoring 节点 类型\nAUTO INTERNAL OVERSEAS IPV6 ALL。",
						},
						"detect_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Detection count。",
						},
						"continue_period": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Continuous 周期 count。",
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

func dataSourceTencentCloudIgtmMonitorsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_igtm_monitors.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(nil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = IgtmService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("filters"); ok {
		filtersSet := v.([]interface{})
		tmpSet := make([]*igtmv20231024.ResourceFilter, 0, len(filtersSet))
		for _, item := range filtersSet {
			filtersMap := item.(map[string]interface{})
			resourceFilter := igtmv20231024.ResourceFilter{}
			if v, ok := filtersMap["name"].(string); ok && v != "" {
				resourceFilter.Name = helper.String(v)
			}

			if v, ok := filtersMap["value"]; ok {
				valueSet := v.(*schema.Set).List()
				for i := range valueSet {
					value := valueSet[i].(string)
					resourceFilter.Value = append(resourceFilter.Value, helper.String(value))
				}
			}

			if v, ok := filtersMap["fuzzy"].(bool); ok {
				resourceFilter.Fuzzy = helper.Bool(v)
			}

			tmpSet = append(tmpSet, &resourceFilter)
		}

		paramMap["Filters"] = tmpSet
	}

	if v, ok := d.GetOkExists("is_detect_num"); ok {
		paramMap["IsDetectNum"] = helper.IntUint64(v.(int))
	}

	var respData []*igtmv20231024.MonitorDetail
	reqErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeIgtmMonitorsByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		respData = result
		return nil
	})

	if reqErr != nil {
		return reqErr
	}

	monitorDataSetList := make([]map[string]interface{}, 0, len(respData))
	if respData != nil {
		for _, monitorDataSet := range respData {
			monitorDataSetMap := map[string]interface{}{}
			if monitorDataSet.MonitorId != nil {
				monitorDataSetMap["monitor_id"] = monitorDataSet.MonitorId
			}

			if monitorDataSet.MonitorName != nil {
				monitorDataSetMap["monitor_name"] = monitorDataSet.MonitorName
			}

			if monitorDataSet.Uin != nil {
				monitorDataSetMap["uin"] = monitorDataSet.Uin
			}

			if monitorDataSet.DetectorGroupIds != nil {
				monitorDataSetMap["detector_group_ids"] = monitorDataSet.DetectorGroupIds
			}

			if monitorDataSet.CheckProtocol != nil {
				monitorDataSetMap["check_protocol"] = monitorDataSet.CheckProtocol
			}

			if monitorDataSet.CheckInterval != nil {
				monitorDataSetMap["check_interval"] = monitorDataSet.CheckInterval
			}

			if monitorDataSet.PingNum != nil {
				monitorDataSetMap["ping_num"] = monitorDataSet.PingNum
			}

			if monitorDataSet.TcpPort != nil {
				monitorDataSetMap["tcp_port"] = monitorDataSet.TcpPort
			}

			if monitorDataSet.Host != nil {
				monitorDataSetMap["host"] = monitorDataSet.Host
			}

			if monitorDataSet.Path != nil {
				monitorDataSetMap["path"] = monitorDataSet.Path
			}

			if monitorDataSet.ReturnCodeThreshold != nil {
				monitorDataSetMap["return_code_threshold"] = monitorDataSet.ReturnCodeThreshold
			}

			if monitorDataSet.EnableRedirect != nil {
				monitorDataSetMap["enable_redirect"] = monitorDataSet.EnableRedirect
			}

			if monitorDataSet.EnableSni != nil {
				monitorDataSetMap["enable_sni"] = monitorDataSet.EnableSni
			}

			if monitorDataSet.PacketLossRate != nil {
				monitorDataSetMap["packet_loss_rate"] = monitorDataSet.PacketLossRate
			}

			if monitorDataSet.Timeout != nil {
				monitorDataSetMap["timeout"] = monitorDataSet.Timeout
			}

			if monitorDataSet.FailTimes != nil {
				monitorDataSetMap["fail_times"] = monitorDataSet.FailTimes
			}

			if monitorDataSet.FailRate != nil {
				monitorDataSetMap["fail_rate"] = monitorDataSet.FailRate
			}

			if monitorDataSet.CreatedOn != nil {
				monitorDataSetMap["created_on"] = monitorDataSet.CreatedOn
			}

			if monitorDataSet.UpdatedOn != nil {
				monitorDataSetMap["updated_on"] = monitorDataSet.UpdatedOn
			}

			if monitorDataSet.DetectorStyle != nil {
				monitorDataSetMap["detector_style"] = monitorDataSet.DetectorStyle
			}

			if monitorDataSet.DetectNum != nil {
				monitorDataSetMap["detect_num"] = monitorDataSet.DetectNum
			}

			if monitorDataSet.ContinuePeriod != nil {
				monitorDataSetMap["continue_period"] = monitorDataSet.ContinuePeriod
			}

			monitorDataSetList = append(monitorDataSetList, monitorDataSetMap)
		}

		_ = d.Set("monitor_data_set", monitorDataSetList)
	}

	d.SetId(helper.BuildToken())
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), monitorDataSetList); e != nil {
			return e
		}
	}

	return nil
}
