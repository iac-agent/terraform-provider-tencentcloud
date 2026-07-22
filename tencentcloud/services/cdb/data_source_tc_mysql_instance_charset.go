package cdb

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdb/v20170320"
)

func DataSourceTencentCloudMysqlInstanceCharset() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudMysqlInstanceCharsetRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "实例ID，格式为：cdb-c1nl9rpv，与云数据库控制台页面显示的实例ID相同，可以通过【查询实例列表】（https://云.tencent.com/document/api/236/15872）接口获取输出参数中InstanceId字段的值。",
			},

			"charset": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "实例的默认字符集，如“latin1”、“utf8”等。",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "用于保存结果。",
			},
		},
	}
}

func dataSourceTencentCloudMysqlInstanceCharsetRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_mysql_instance_charset.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	var instanceId string
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}

	service := MysqlService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	var instanceCharset *cdb.DescribeDBInstanceCharsetResponseParams
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeMysqlInstanceCharsetByFilter(ctx, instanceId)
		if e != nil {
			return tccommon.RetryError(e)
		}
		instanceCharset = result
		return nil
	})
	if err != nil {
		return err
	}

	if instanceCharset != nil {
		_ = d.Set("charset", instanceCharset.Charset)
	}

	d.SetId(instanceId)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), d); e != nil {
			return e
		}
	}
	return nil
}
