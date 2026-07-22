package scf

import (
	"context"
	"log"
	"time"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/pkg/errors"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudScfLogs() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudScfLogsRead,
		Schema: map[string]*schema.Schema{
			"function_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "名称 SCF 函数 到 是 queried。",
			},
			"offset": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "Log 偏移量，默认为 `0`，偏移量+限制 不能 是 greater 比 10000。",
			},
			"limit": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     10000,
				Description: "数量 logs， 默认为 `10000`，偏移量+限制 不能 是 greater 比 10000。",
			},
			"order": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue(SCF_LOGS_ORDERS),
				Default:      SCF_LOGS_ORDER_DESC,
				Description:  "顺序 到 sort 日志，可选 值 `desc` 和 `asc`，默认值 `desc`。",
			},
			"order_by": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue(SCF_LOGS_ORDER_BY),
				Default:      SCF_LOGS_ORDER_BY_START_TIME,
				Description:  "Sort logs according 到 following 字段: `function_name`，`时长`，`mem_usage`，`start_time`，默认值 `start_time`。",
			},
			"ret_code": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue(SCF_LOGS_RET_CODES),
				Description:  "Use 到 过滤器 日志，可选 值: `not0` 仅 返回error 日志. `is0` 仅 返回correct 日志. `TimeLimitExceeded` 返回log 的 函数 call 超时. `ResourceLimitExceeded` 返回function call generation 资源 overrun 日志. `UserCodeException` 返回logs 的 用户 代码 错误 该 occurred 在 函数 call. Not passing 参数 表示 returning all logs。",
			},
			"namespace": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "default",
				Description: "Namespace 的 SCF 函数 到 是 queried。",
			},
			"invoke_request_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Corresponding requestId 当 executing 函数。",
			},
			"start_time": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: tccommon.ValidateTime(SCF_LOGS_DESCRIBE_TIME_FORMAT),
				Description:  "开始时间 的 查询， 格式 是 `2017-05-16 20:00:00`，其中 可以 仅 是 within 一个 day 从 `end_time`。",
			},
			"end_time": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: tccommon.ValidateTime(SCF_LOGS_DESCRIBE_TIME_FORMAT),
				Description:  "结束时间 的 查询， 格式 是 `2017-05-16 20:00:00`，其中 可以 仅 是 within 一个 day 从 `start_time`。",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "用于保存结果。",
			},

			// computed
			"logs": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "An 信息 列表 logs. Each element 包含following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"function_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "名称 SCF 函数。",
						},
						"ret_msg": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "返回值 after 函数 execution 是 completed。",
						},
						"request_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Execute requestId corresponding 到 函数。",
						},
						"start_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Point 在 时间 在 其中 函数 begins execution。",
						},
						"ret_code": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Execution 结果 的 函数，`0` 表示 execution 是 successful，other 值 indicate failure。",
						},
						"invoke_finished": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "是否function call 结束，`1` 表示 execution 结束，other 值 indicate call exception。",
						},
						"duration": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "Function 执行时间-consuming，单位 是 ms。",
						},
						"bill_duration": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Function billing 时间，according 到 时长 up 到 last 100ms，单位 是 ms。",
						},
						"mem_usage": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "actual 内存 大小 consumed 在 execution 的 函数，单位 是 Byte。",
						},
						"log": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Log output during 函数 execution。",
						},
						"level": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Log 级别",
						},
						"source": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Log 来源",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudScfLogsRead(d *schema.ResourceData, m interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_scf_logs.read")()
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := ScfService{client: m.(tccommon.ProviderMeta).GetAPIV3Conn()}

	functionName := d.Get("function_name").(string)
	namespace := d.Get("namespace").(string)

	offset := d.Get("offset").(int)
	limit := d.Get("limit").(int)
	if offset+limit > 10000 {
		return errors.New("offset + limit can't greater than 10000")
	}

	order := d.Get("order").(string)
	orderBy := d.Get("order_by").(string)

	var (
		retCode         *string
		invokeRequestId *string
		startTime       *string
		endTime         *string
	)

	if raw, ok := d.GetOk("ret_code"); ok {
		retCode = helper.String(raw.(string))
	}
	if raw, ok := d.GetOk("invoke_request_id"); ok {
		invokeRequestId = helper.String(raw.(string))
	}

	if raw, ok := d.GetOk("start_time"); ok {
		startTime = helper.String(raw.(string))
	}
	if raw, ok := d.GetOk("end_time"); ok {
		endTime = helper.String(raw.(string))
	}
	if err := helper.CheckIfSetTogether(d, "start_time", "end_time"); err != nil {
		return err
	}

	if startTime != nil && endTime != nil {
		startTime, _ := time.Parse(SCF_LOGS_DESCRIBE_TIME_FORMAT, *startTime)
		endTime, _ := time.Parse(SCF_LOGS_DESCRIBE_TIME_FORMAT, *endTime)

		if endTime.Sub(startTime) > 24*time.Hour {
			return errors.New("end_time - start_time can't greater then 1 day")
		}
	}

	respLogs, err := service.DescribeLogs(ctx,
		functionName, namespace, order, orderBy,
		offset, limit,
		retCode, invokeRequestId, startTime, endTime,
	)
	if err != nil {
		log.Printf("[CRITAL]%s read function logs failed: %+v", logId, err)
		return err
	}

	logs := make([]map[string]interface{}, 0, len(respLogs))
	ids := make([]string, 0, len(respLogs))
	for _, l := range respLogs {
		ids = append(ids, *l.RequestId)

		logs = append(logs, map[string]interface{}{
			"function_name":   l.FunctionName,
			"ret_msg":         l.RetMsg,
			"request_id":      l.RequestId,
			"start_time":      l.StartTime,
			"ret_code":        l.RetCode,
			"invoke_finished": l.InvokeFinished,
			"duration":        l.Duration,
			"bill_duration":   l.BillDuration,
			"mem_usage":       l.MemUsage,
			"log":             l.Log,
			"level":           l.Level,
			"source":          l.Source,
		})
	}

	_ = d.Set("logs", logs)
	d.SetId(helper.DataResourceIdsHash(ids))

	if output, ok := d.GetOk("result_output_file"); ok && output.(string) != "" {
		if err := tccommon.WriteToFile(output.(string), logs); err != nil {
			err = errors.WithStack(err)
			log.Printf("[CRITAL]%s output file[%s] fail, reason: %+v", logId, output.(string), err)
			return err
		}
	}

	return nil
}
