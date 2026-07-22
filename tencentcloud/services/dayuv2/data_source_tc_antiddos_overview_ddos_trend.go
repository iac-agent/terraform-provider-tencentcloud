package dayuv2

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svcantiddos "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/antiddos"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	antiddos "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/antiddos/v20200309"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudAntiddosOverviewDdosTrend() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudAntiddosOverviewDdosTrendRead,
		Schema: map[string]*schema.Schema{
			"period": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Statistical granularity，值 [300 (5 minutes)，3600 (hours)，86400 (days)]。",
			},

			"start_time": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "StartTime。",
			},

			"end_time": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "EndTime。",
			},

			"metric_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Indicator，值 [bps (attack 流量 带宽，pps (attack packet 速率)]。",
			},

			"business": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Dayu sub product 代码 (bgpip 表示 advanced defense IP; net 表示 professional 版本 的 advanced defense IP)。",
			},

			"ip_list": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "实例 IpList。",
			},

			"data": {
				Computed: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Description: "Array，attack 流量 带宽 在 Mbps，packet 速率 在 pps。",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "用于保存结果。",
			},
		},
	}
}

func dataSourceTencentCloudAntiddosOverviewDdosTrendRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_antiddos_overview_ddos_trend.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, _ := d.GetOk("period"); v != nil {
		paramMap["Period"] = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("start_time"); ok {
		paramMap["StartTime"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("end_time"); ok {
		paramMap["EndTime"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("metric_name"); ok {
		paramMap["MetricName"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("business"); ok {
		paramMap["Business"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("ip_list"); ok {
		ipListSet := v.(*schema.Set).List()
		paramMap["IpList"] = helper.InterfacesStringsPoint(ipListSet)
	}

	service := svcantiddos.NewAntiddosService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())

	var describeOverviewDDoSTrendResponseParams *antiddos.DescribeOverviewDDoSTrendResponseParams
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeAntiddosOverviewDdosTrendByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		describeOverviewDDoSTrendResponseParams = result
		return nil
	})
	if err != nil {
		return err
	}

	resultMap := make(map[string]interface{})
	if describeOverviewDDoSTrendResponseParams.Data != nil {
		resultMap["data"] = describeOverviewDDoSTrendResponseParams.Data
		_ = d.Set("data", describeOverviewDDoSTrendResponseParams.Data)
	}

	d.SetId(helper.BuildToken())
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), resultMap); e != nil {
			return e
		}
	}
	return nil
}
