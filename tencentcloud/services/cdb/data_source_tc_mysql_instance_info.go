package cdb

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdb/v20170320"
)

func DataSourceTencentCloudMysqlInstanceInfo() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudMysqlInstanceInfoRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "实例 ID。",
			},

			"instance_name": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "实例名称。",
			},

			"encryption": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "是否启用加密，YES 启用，NO 不启用。",
			},

			"key_id": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "用于加密的密钥 ID。",
			},

			"key_region": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "密钥所在的区域。",
			},

			"default_kms_region": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "当前CDB后端服务使用的KMS服务的默认区域。",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "用于保存结果。",
			},
		},
	}
}

func dataSourceTencentCloudMysqlInstanceInfoRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_mysql_instance_info.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	var instanceId string
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}

	service := MysqlService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	var instanceInfo *cdb.DescribeDBInstanceInfoResponseParams
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeMysqlInstanceInfoById(ctx, instanceId)
		if e != nil {
			return tccommon.RetryError(e)
		}
		instanceInfo = result
		return nil
	})
	if err != nil {
		return err
	}

	if instanceInfo.InstanceName != nil {
		_ = d.Set("instance_name", instanceInfo.InstanceName)
	}

	if instanceInfo.Encryption != nil {
		_ = d.Set("encryption", instanceInfo.Encryption)
	}

	if instanceInfo.KeyId != nil {
		_ = d.Set("key_id", instanceInfo.KeyId)
	}

	if instanceInfo.KeyRegion != nil {
		_ = d.Set("key_region", instanceInfo.KeyRegion)
	}

	if instanceInfo.DefaultKmsRegion != nil {
		_ = d.Set("default_kms_region", instanceInfo.DefaultKmsRegion)
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
