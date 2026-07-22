package dbbrain

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dbbrain "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dbbrain/v20210527"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudDbbrainDbSpaceStatus() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudDbbrainDbSpaceStatusRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "实例 ID.",
			},

			"range_days": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "数量 的 days 在 时间 周期, deadline 是 当前 day, 和 默认值 是 7 days.",
			},

			"product": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Service product 类型, 支持 值 include: mysql - 云 数据库 MySQL, cynosdb - 云 数据库 CynosDB 对于 MySQL, 默认值 是 mysql.",
			},

			"growth": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Disk growth (MB).",
			},

			"remain": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Disk remaining (MB).",
			},

			"total": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Total 磁盘 大小 (MB).",
			},

			"available_days": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Estimated 数量 的 days 可用.",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used 到 save results.",
			},
		},
	}
}

func dataSourceTencentCloudDbbrainDbSpaceStatusRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_dbbrain_db_space_status.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var instanceId string

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		paramMap["InstanceId"] = helper.String(v.(string))
		instanceId = v.(string)
	}

	if v, _ := d.GetOk("range_days"); v != nil {
		paramMap["RangeDays"] = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("product"); ok {
		paramMap["Product"] = helper.String(v.(string))
	}

	service := DbbrainService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var rows *dbbrain.DescribeDBSpaceStatusResponseParams

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeDbbrainDbSpaceStatusByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		rows = result
		return nil
	})
	if err != nil {
		return err
	}

	tmpList := []map[string]interface{}{}

	if rows != nil {

		if rows.Growth != nil {
			_ = d.Set("growth", rows.Growth)
		}

		if rows.Remain != nil {
			_ = d.Set("remain", rows.Remain)
		}

		if rows.Total != nil {
			_ = d.Set("total", rows.Total)
		}

		if rows.AvailableDays != nil {
			_ = d.Set("available_days", rows.AvailableDays)
		}
		tmpList = append(tmpList, map[string]interface{}{
			"growth":         rows.Growth,
			"remain":         rows.Remain,
			"total":          rows.Total,
			"available_days": rows.AvailableDays,
		})

	}

	d.SetId(helper.DataResourceIdHash(instanceId))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
