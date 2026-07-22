package cdb

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdb/v20170320"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudMysqlDbFeatures() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudMysqlDbFeaturesRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "实例ID，格式为：cdb-c1nl9rpv或cdbro-c1nl9rpv，与云数据库控制台页面显示的实例ID相同。",
			},

			"is_support_audit": {
				Computed:    true,
				Type:        schema.TypeBool,
				Description: "是否支持数据库审计功能。",
			},

			"audit_need_upgrade": {
				Computed:    true,
				Type:        schema.TypeBool,
				Description: "是否开启审计需要升级内核版本。",
			},

			"is_support_encryption": {
				Computed:    true,
				Type:        schema.TypeBool,
				Description: "是否支持数据库加密功能。",
			},

			"encryption_need_upgrade": {
				Computed:    true,
				Type:        schema.TypeBool,
				Description: "是否启用加密需要升级内核版本。",
			},

			"is_remote_ro": {
				Computed:    true,
				Type:        schema.TypeBool,
				Description: "是否为远程只读实例。",
			},

			"master_region": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "主实例所在地域。",
			},

			"is_support_update_sub_version": {
				Computed:    true,
				Type:        schema.TypeBool,
				Description: "是否支持小版本升级。",
			},

			"current_sub_version": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "当前内核版本。",
			},

			"target_sub_version": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "可用于升级的内核版本。",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "用于保存结果。",
			},
		},
	}
}

func dataSourceTencentCloudMysqlDbFeaturesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_mysql_db_features.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	var instanceId string
	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		paramMap["InstanceId"] = helper.String(v.(string))
	}

	service := MysqlService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	var dbFeatures *cdb.DescribeDBFeaturesResponseParams
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeMysqlDbFeaturesByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		dbFeatures = result
		return nil
	})
	if err != nil {
		return err
	}

	if dbFeatures.IsSupportAudit != nil {
		_ = d.Set("is_support_audit", dbFeatures.IsSupportAudit)
	}

	if dbFeatures.AuditNeedUpgrade != nil {
		_ = d.Set("audit_need_upgrade", dbFeatures.AuditNeedUpgrade)
	}

	if dbFeatures.IsSupportEncryption != nil {
		_ = d.Set("is_support_encryption", dbFeatures.IsSupportEncryption)
	}

	if dbFeatures.EncryptionNeedUpgrade != nil {
		_ = d.Set("encryption_need_upgrade", dbFeatures.EncryptionNeedUpgrade)
	}

	if dbFeatures.IsRemoteRo != nil {
		_ = d.Set("is_remote_ro", dbFeatures.IsRemoteRo)
	}

	if dbFeatures.MasterRegion != nil {
		_ = d.Set("master_region", dbFeatures.MasterRegion)
	}

	if dbFeatures.IsSupportUpdateSubVersion != nil {
		_ = d.Set("is_support_update_sub_version", dbFeatures.IsSupportUpdateSubVersion)
	}

	if dbFeatures.CurrentSubVersion != nil {
		_ = d.Set("current_sub_version", dbFeatures.CurrentSubVersion)
	}

	if dbFeatures.TargetSubVersion != nil {
		_ = d.Set("target_sub_version", dbFeatures.TargetSubVersion)
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
