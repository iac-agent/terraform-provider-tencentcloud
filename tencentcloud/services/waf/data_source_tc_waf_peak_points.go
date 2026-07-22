package waf

import (
	"context"
	"strconv"
	"time"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	waf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/waf/v20180125"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudWafPeakPoints() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudWafPeakPointsRead,
		Schema: map[string]*schema.Schema{
			"from_time": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "开始时间。",
			},
			"to_time": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "结束时间。",
			},
			"domain": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "域名 名称 到 是 queried. 如果 all 域名 名称 数据 是 queried，此 参数 是 不 filled 在。",
			},
			"edition": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Only support sparta-waf 和 clb-waf. 如果未传入，there 将 是 无 filtering。",
			},
			"instance_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "WAF 实例 ID，如果未传入，there 将 是 无 filtering。",
			},
			"metric_name": {
				Optional:     true,
				Type:         schema.TypeString,
				ValidateFunc: tccommon.ValidateAllowedStringValue(MetricNameList),
				Description:  "Twelve 值 是 可用: `访问`-Peak qps trend chart; `botAccess`- bot peak qps trend chart; `down`-Downstream peak 带宽 trend chart; `up`-Upstream peak 带宽 trend chart; `attack`-Trend chart 的 总数 数量 web attacks; `cc`-Trend chart 的 总数 数量 CC attacks; `bw`- Black IP Attack Total Trend Chart; `tamper`- Anti Tamper Attack Total Trend Chart; `leak`- Trend chart 的 总数 数量 anti leakage attacks; `acl`- Trend chart 的 总数 数量 访问 control attacks; `http_status`- Trend chart 的 状态 代码 频率; `wx_access`- WeChat Mini Program Peak QPS Trend Chart。",
			},
			"points": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "point 列表。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Second 级别 时间戳。",
						},
						"access": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "qps。",
						},
						"bot_access": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Bot qps。",
						},
						"attack": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "数量 web attacks。",
						},
						"cc": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "数量 cc attacks。",
						},
						"down": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Peak downlink 带宽，单位 B。",
						},
						"up": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Peak uplink 带宽，单位 B。",
						},
						"status_server_error": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Trend chart 的 数量 状态 codes 返回 通过 WAF 到 服务器。",
						},
						"status_client_error": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Trend chart 的 数量 状态 codes 返回 通过 WAF 到 客户端。",
						},
						"status_redirect": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Trend chart 的 数量 状态 codes 返回 通过 WAF 到 客户端。",
						},
						"status_ok": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Trend chart 的 数量 状态 codes 返回 通过 WAF 到 客户端。",
						},
						"upstream_server_error": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Trend chart 的 数量 状态 codes 返回 到 WAF 通过 源站 site。",
						},
						"upstream_client_error": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Trend chart 的 数量 状态 codes 返回 到 WAF 通过 源站 site。",
						},
						"upstream_redirect": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Trend chart 的 数量 状态 codes 返回 到 WAF 通过 源站 site。",
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

func dataSourceTencentCloudWafPeakPointsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_waf_peak_points.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = WafService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		points  []*waf.PeakPointsItem
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("from_time"); ok {
		paramMap["FromTime"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("to_time"); ok {
		paramMap["ToTime"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("domain"); ok {
		paramMap["Domain"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("edition"); ok {
		paramMap["Edition"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("instance_id"); ok {
		paramMap["InstanceID"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("metric_name"); ok {
		paramMap["MetricName"] = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeWafPeakPointsByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		points = result
		return nil
	})

	if err != nil {
		return err
	}

	tmpList := make([]map[string]interface{}, 0, len(points))

	if points != nil {
		for _, point := range points {
			dMap := map[string]interface{}{}

			if point.Time != nil {
				dMap["time"] = point.Time
			}

			if point.Access != nil {
				dMap["access"] = point.Access
			}

			if point.BotAccess != nil {
				dMap["bot_access"] = point.BotAccess
			}

			if point.Attack != nil {
				dMap["attack"] = point.Attack
			}

			if point.Cc != nil {
				dMap["cc"] = point.Cc
			}

			if point.Down != nil {
				dMap["down"] = point.Down
			}

			if point.Up != nil {
				dMap["up"] = point.Up
			}

			if point.StatusServerError != nil {
				dMap["status_server_error"] = point.StatusServerError
			}

			if point.StatusClientError != nil {
				dMap["status_client_error"] = point.StatusClientError
			}

			if point.StatusRedirect != nil {
				dMap["status_redirect"] = point.StatusRedirect
			}

			if point.StatusOk != nil {
				dMap["status_ok"] = point.StatusOk
			}

			if point.UpstreamServerError != nil {
				dMap["upstream_server_error"] = point.UpstreamServerError
			}

			if point.UpstreamClientError != nil {
				dMap["upstream_client_error"] = point.UpstreamClientError
			}

			if point.UpstreamRedirect != nil {
				dMap["upstream_redirect"] = point.UpstreamRedirect
			}

			tmpList = append(tmpList, dMap)
		}

		_ = d.Set("points", tmpList)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), tmpList); e != nil {
			return e
		}
	}

	return nil
}
