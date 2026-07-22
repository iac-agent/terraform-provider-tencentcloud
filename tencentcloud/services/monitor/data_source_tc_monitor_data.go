package monitor

import (
	"crypto/md5"
	"fmt"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	monitor "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/monitor/v20180724"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

func DataSourceTencentCloudMonitorData() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentMonitorDataRead,
		Schema: map[string]*schema.Schema{
			"namespace": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Namespace 的 each 云 product 在 监控 系统，refer 到 `数据.tencentcloud_monitor_product_namespace`。",
			},
			"metric_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "指标名称，please refer 到 documentation 的 监控 interface 的 each product。",
			},
			"dimensions": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "实例 dimension 名称，eg: `实例 ID` 对于 cvm。",
						},
						"value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "实例 dimension 值，eg: `ins-j0hk02zo` 对于 cvm。",
						},
					},
				},
				Description: "Dimensional composition 的 实例 objects。",
			},
			"period": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     300,
				Description: "Statistical 周期",
			},
			"start_time": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "开始时间 对于 此 查询，eg:`2018-09-22T19:51:23+08:00`。",
			},
			"end_time": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "结束时间 对于 此 查询，eg:`2018-09-22T20:00:00+08:00`。",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "用于存储结果。",
			},
			// Computed values
			"list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A 列表 数据 point. Each element 包含following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"timestamp": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Statistical 时间戳。",
						},
						"value": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "Statistical 值",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentMonitorDataRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_monitor_data.read")()

	var (
		monitorService = MonitorService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		request        = monitor.NewGetMonitorDataRequest()
		response       *monitor.GetMonitorDataResponse
		err            error
		dimensions     = d.Get("dimensions").([]interface{})
		instance       monitor.Instance
		list           []interface{}
	)
	request.Namespace = helper.String(d.Get("namespace").(string))
	request.MetricName = helper.String(d.Get("metric_name").(string))
	request.Period = helper.IntUint64(d.Get("period").(int))
	request.StartTime = helper.String(d.Get("start_time").(string))
	request.EndTime = helper.String(d.Get("end_time").(string))
	request.Instances = []*monitor.Instance{&instance}

	for _, dimension := range dimensions {
		kv := dimension.(map[string]interface{})
		instance.Dimensions = append(instance.Dimensions, &monitor.Dimension{
			Name:  helper.String(kv["name"].(string)),
			Value: helper.String(kv["value"].(string)),
		})
	}

	if err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		if response, err = monitorService.client.UseMonitorClient().GetMonitorData(request); err != nil {
			return tccommon.RetryError(err, tccommon.InternalError)
		}
		return nil
	}); err != nil {
		return err
	}

	if len(response.Response.DataPoints) > 0 {
		data := response.Response.DataPoints[0]
		min := len(data.Values)
		if min > len(data.Timestamps) {
			min = len(data.Timestamps)
		}
		for i := 0; i < min; i++ {
			kv := make(map[string]interface{})
			kv["timestamp"] = int64(*data.Timestamps[i])
			kv["value"] = data.Values[i]
			list = append(list, kv)
		}
	}

	if err = d.Set("list", list); err != nil {
		return err
	}

	md := md5.New()
	_, _ = md.Write([]byte(request.ToJsonString()))
	id := fmt.Sprintf("%x", md.Sum(nil))
	d.SetId(id)
	if output, ok := d.GetOk("result_output_file"); ok {
		return tccommon.WriteToFile(output.(string), list)
	}
	return nil
}
