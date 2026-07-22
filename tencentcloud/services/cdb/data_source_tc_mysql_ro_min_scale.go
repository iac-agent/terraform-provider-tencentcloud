package cdb

import (
	"context"
	"strconv"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdb/v20170320"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudMysqlRoMinScale() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudMysqlRoMinScaleRead,
		Schema: map[string]*schema.Schema{
			"ro_instance_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "只读实例ID，格式为cdbro-c1nl9rpv，与云数据库控制台页面显示的实例ID一致。该参数与MasterInstanceId参数不能同时为空。",
			},

			"master_instance_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "主实例ID与云数据库控制台页面显示的实例ID一致，格式为：cdb-c1nl9rpv。该参数与RoInstanceId参数不能同时为空。注意，当入参包含RoInstanceId时，返回值为只读实例升级时的最小规格；当入参仅包含MasterInstanceId时，返回值为购买只读实例时的最低规格。",
			},

			"memory": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "内存规格大小，单位：MB。",
			},

			"volume": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "磁盘规格大小，单位：GB。",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "用于保存结果。",
			},
		},
	}
}

func dataSourceTencentCloudMysqlRoMinScaleRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_mysql_ro_min_scale.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("ro_instance_id"); ok {
		paramMap["RoInstanceId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("master_instance_id"); ok {
		paramMap["MasterInstanceId"] = helper.String(v.(string))
	}

	service := MysqlService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	var minScale *cdb.DescribeRoMinScaleResponseParams
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeMysqlRoMinScaleByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		minScale = result
		return nil
	})
	if err != nil {
		return err
	}

	if minScale.Memory != nil {
		_ = d.Set("memory", minScale.Memory)
	}

	if minScale.Volume != nil {
		_ = d.Set("volume", minScale.Volume)
	}

	d.SetId(helper.DataResourceIdsHash([]string{strconv.FormatInt(*minScale.Memory, 10), strconv.FormatInt(*minScale.Volume, 10)}))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), map[string]interface{}{
			"memory": minScale.Memory,
			"volume": minScale.Volume,
		}); e != nil {
			return e
		}
	}
	return nil
}
