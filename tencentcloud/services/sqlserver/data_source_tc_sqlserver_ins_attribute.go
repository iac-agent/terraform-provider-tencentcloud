package sqlserver

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sqlserver "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sqlserver/v20180328"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudSqlserverInsAttribute() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudSqlserverInsAttributeRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "实例 ID.",
			},
			"regular_backup_enable": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Archive 备份 状态. 有效 值: 启用 (已启用), disable (已禁用).",
			},
			"regular_backup_save_days": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Archive 备份 retention 周期: [90-3650] days.",
			},
			"regular_backup_strategy": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Archive 备份 策略. 有效 值: years (yearly); quarters (quarterly);months` (monthly).",
			},
			"regular_backup_counts": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "数量 的 retained archive backups.",
			},
			"regular_backup_start_time": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Archive 备份 start date 在 YYYY-MM-DD 格式, 其中 是 当前 时间 通过 默认值.",
			},
			"blocked_threshold": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Block process 阈值 在 milliseconds.",
			},
			"event_save_days": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Retention 周期 对于 files 的 slow SQL, blocking, deadlock, 和 extended events.",
			},
			"tde_config": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "TDE Transparent Data Encryption Configuration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"certificate_attribution": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Certificate ownership. Self - indicates 使用 account's own 证书, others - indicates referencing certificates 从 other accounts, 和 none - indicates 无 证书.",
						},
						"encryption": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "TDE 加密, '启用' - 已启用, 'disable' - 不 已启用.",
						},
						"quote_uin": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Other primary account IDs referenced 当 activating TDE 加密\nNote: 此 字段 可能 返回 null, indicating 该 有效 值 不能 是 获取.",
						},
					},
				},
			},
			"ssl_config": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "SSL 加密.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"encryption": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "SSL 加密 状态, 启用 - turned 在, disable-不 turned 在, enable_doing - enabling, disable_doing-closing, renew_doing-updating, wait_doing-wait 对于 execution within maintenance 时间 注意: 此 字段 可能 返回 null, indicating 该 无 有效 值 可以 是 获取.",
						},
						"ssl_validity_period": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "SSL 证书 validity 周期, 时间 格式 YYYY-MM-DD HH:MM:SS 注意: 此 字段 可能 返回 null, indicating 该 无 有效 值 可以 是 获取.",
						},
						"ssl_validity": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "SSL 证书 validity, 0-无效, 1-有效 注意: 此 字段 可能 返回 null, indicating 该 无 有效 值 可以 是 获取.",
						},
					},
				},
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used 到 save results.",
			},
		},
	}
}

func dataSourceTencentCloudSqlserverInsAttributeRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_sqlserver_ins_attribute.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId        = tccommon.GetLogId(tccommon.ContextNil)
		ctx          = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service      = SqlserverService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		insAttribute *sqlserver.DescribeDBInstancesAttributeResponseParams
		instanceId   string
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		paramMap["InstanceId"] = helper.String(v.(string))
		instanceId = v.(string)
	}

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeSqlserverInsAttributeByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		insAttribute = result
		return nil
	})

	if err != nil {
		return err
	}

	if insAttribute.InstanceId != nil {
		_ = d.Set("instance_id", instanceId)
	}

	if insAttribute.RegularBackupEnable != nil {
		_ = d.Set("regular_backup_enable", insAttribute.RegularBackupEnable)
	}

	if insAttribute.RegularBackupSaveDays != nil {
		_ = d.Set("regular_backup_save_days", insAttribute.RegularBackupSaveDays)
	}

	if insAttribute.RegularBackupStrategy != nil {
		_ = d.Set("regular_backup_strategy", insAttribute.RegularBackupStrategy)
	}

	if insAttribute.RegularBackupCounts != nil {
		_ = d.Set("regular_backup_counts", insAttribute.RegularBackupCounts)
	}

	if insAttribute.RegularBackupStartTime != nil {
		_ = d.Set("regular_backup_start_time", insAttribute.RegularBackupStartTime)
	}

	if insAttribute.BlockedThreshold != nil {
		_ = d.Set("blocked_threshold", insAttribute.BlockedThreshold)
	}

	if insAttribute.EventSaveDays != nil {
		_ = d.Set("event_save_days", insAttribute.EventSaveDays)
	}

	if insAttribute.TDEConfig != nil {
		tmpList := make([]map[string]interface{}, 0)
		configMap := map[string]interface{}{}
		if insAttribute.TDEConfig.CertificateAttribution != nil {
			configMap["certificate_attribution"] = insAttribute.TDEConfig.CertificateAttribution
		}

		if insAttribute.TDEConfig.Encryption != nil {
			configMap["encryption"] = insAttribute.TDEConfig.Encryption
		}

		if insAttribute.TDEConfig.QuoteUin != nil {
			configMap["quote_uin"] = insAttribute.TDEConfig.QuoteUin
		}

		tmpList = append(tmpList, configMap)

		_ = d.Set("tde_config", tmpList)
	}

	if insAttribute.SSLConfig != nil {
		tmpList := make([]map[string]interface{}, 0)
		configMap := map[string]interface{}{}
		if insAttribute.SSLConfig.Encryption != nil {
			configMap["encryption"] = insAttribute.SSLConfig.Encryption
		}

		if insAttribute.SSLConfig.SSLValidityPeriod != nil {
			configMap["ssl_validity_period"] = insAttribute.SSLConfig.SSLValidityPeriod
		}

		if insAttribute.SSLConfig.SSLValidity != nil {
			configMap["ssl_validity"] = insAttribute.SSLConfig.SSLValidity
		}

		tmpList = append(tmpList, configMap)

		_ = d.Set("ssl_config", tmpList)
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
