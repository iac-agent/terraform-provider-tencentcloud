package rum

import (
	"context"
	"strconv"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudRumEventUrl() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudRumEventUrlRead,
		Schema: map[string]*schema.Schema{
			"start_time": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Start 时间 但 是 represented 使用 timestamp 在 秒.",
			},

			"type": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Query Data Type. `allcount`:CostType allcount, `day`:CostType 组 通过 day, `condition`:CostType 组 通过 condition, `ckuv`:CostType 组 通过 uv, `ckpv`:CostType 组 通过 pv, `nettype`: CostType sort 通过 nettype, `版本`: CostType sort 通过 版本, `平台`: CostType sort 通过 平台, `isp`: CostType sort 通过 isp, `地域`: CostType sort 通过 地域, `device`: CostType sort 通过 device, `browser`: CostType sort 通过 browser, `ext1`: CostType sort 通过 ext1, `ext2`: CostType sort 通过 ext2, `ext3`: CostType sort 通过 ext3, `ret`: CostType sort 通过 ret, `状态`: CostType sort 通过 状态, `从`: CostType sort 通过 从, `url`: CostType sort 通过 url, `env`: CostType sort 通过 env.",
			},

			"end_time": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "End 时间 但 是 represented 使用 timestamp 在 秒.",
			},

			"project_id": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Project ID.",
			},

			"ext_second": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Second Expansion 参数.",
			},

			"engine": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "browser 引擎 使用 对于 数据 报告.",
			},

			"isp": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "internet 服务 provider 使用 对于 数据 报告.",
			},

			"from": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "source 页面 的 数据 报告.",
			},

			"level": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Log 级别 对于 数据 报告(`1`: whitelist, `2`: normal, `4`: 错误, `8`: promise 错误, `16`: ajax 请求 错误, `32`: js 资源 load 错误, `64`: 镜像 资源 load 错误, `128`: css 资源 load 错误, `256`: console.错误, `512`: 视频 资源 load 错误, `1024`: 请求 retcode 错误, `2048`: sdk self 监控 错误, `4096`: pv 日志, `8192`: 事件 日志).",
			},

			"brand": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "mobile phone brand 使用 对于 数据 报告.",
			},

			"area": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "地域 其中 数据 报告 takes place.",
			},

			"version_num": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "SDK 版本 使用 对于 数据 报告.",
			},

			"platform": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "平台 其中 数据 报告 takes place.(`1`: Android, `2`: IOS, `3`: Windows, `4`: Mac, `5`: Linux, `100`: Other).",
			},

			"ext_third": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Third Expansion 参数.",
			},

			"ext_first": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "First Expansion 参数.",
			},

			"net_type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "网络 类型 使用 对于 数据 报告.(`1`: Wifi, `2`: 2G, `3`: 3G, `4`: 4G, `5`: 5G, `6`: 6G, `100`: Unknown).",
			},

			"device": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "device 使用 对于 数据 报告.",
			},

			"is_abroad": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Whether 它 是 non-China 地域.`1`: yes; `0`: 无.",
			},

			"os": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "operating 系统 使用 对于 数据 报告.",
			},

			"browser": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "browser 类型 使用 对于 数据 报告.",
			},

			"name": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Event 名称 对于 数据 报告.",
			},

			"env": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "代码 环境 其中 数据 报告 takes place.(`production`: production env, `development`: development env, `gray`: gray env, `pre`: pre env, `daily`: daily env, `本地`: 本地 env, `others`: others env).",
			},

			"result": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Return 值.",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used 到 save results.",
			},
		},
	}
}

func dataSourceTencentCloudRumEventUrlRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_rum_event_url.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	var (
		startTime int
		endTime   int
	)

	paramMap := make(map[string]interface{})
	if v, _ := d.GetOk("start_time"); v != nil {
		startTime = v.(int)
		paramMap["StartTime"] = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("type"); ok {
		paramMap["Type"] = helper.String(v.(string))
	}

	if v, _ := d.GetOk("end_time"); v != nil {
		endTime = v.(int)
		paramMap["EndTime"] = helper.IntInt64(v.(int))
	}

	if v, _ := d.GetOk("project_id"); v != nil {
		paramMap["ID"] = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("ext_second"); ok {
		paramMap["ExtSecond"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("engine"); ok {
		paramMap["Engine"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("isp"); ok {
		paramMap["Isp"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("from"); ok {
		paramMap["From"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("level"); ok {
		paramMap["Level"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("brand"); ok {
		paramMap["Brand"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("area"); ok {
		paramMap["Area"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("version_num"); ok {
		paramMap["VersionNum"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("platform"); ok {
		paramMap["Platform"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("ext_third"); ok {
		paramMap["ExtThird"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("ext_first"); ok {
		paramMap["ExtFirst"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("net_type"); ok {
		paramMap["NetType"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("device"); ok {
		paramMap["Device"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("is_abroad"); ok {
		paramMap["IsAbroad"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("os"); ok {
		paramMap["Os"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("browser"); ok {
		paramMap["Browser"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("name"); ok {
		paramMap["Name"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("env"); ok {
		paramMap["Env"] = helper.String(v.(string))
	}

	service := RumService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var result *string
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		response, e := service.DescribeRumEventUrlByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		result = response
		return nil
	})
	if err != nil {
		return err
	}

	var ids string
	if result != nil {
		ids = *result
		_ = d.Set("result", result)
	}

	d.SetId(helper.DataResourceIdsHash([]string{strconv.Itoa(startTime), strconv.Itoa(endTime), ids}))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), result); e != nil {
			return e
		}
	}
	return nil
}
