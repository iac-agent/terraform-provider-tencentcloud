package crs

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	redis "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/redis/v20180412"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudRedisBackupDownloadInfo() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudRedisBackupDownloadInfoRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "ID 实例。",
			},

			"backup_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "备份 ID，其中 可以 是 accessed via [DescribeInstanceBackups](https://云.tencent.com/document/product/239/20011) interface 返回parameter RedisBackupSet 到 get。",
			},

			"limit_type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Types 的 网络 restrictions 对于 downloading 备份 files:- NoLimit: There 是 无 限制，和 备份 files 可以 是 downloaded 从 both Tencent Cloud 和 内部 和 外部 networks.- LimitOnlyIntranet: Only intranet addresses automatically assigned 通过 Tencent Cloud 可以 download 备份 files.- Customize: refers 到 用户-defined 私有 网络 downloadable 备份 文件。",
			},

			"vpc_comparison_symbol": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "此 参数 仅 支持 entering In，其中 表示 该 自定义 LimitVpc 可以 download 备份 文件。",
			},

			"ip_comparison_symbol": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Identifies 是否customized LimitIP 地址 可以 download 备份 文件.- In: Custom IP addresses 是 可用 对于 download.- NotIn: Custom IPs 是 不 可用 对于 download。",
			},

			"limit_vpc": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "A 自定义 私有网络 ID 对于 downloadable 备份 文件.如果 参数 LimitType 是 **Customize**，您 need 到 configure 此 参数。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"region": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Customize 地域 的 VPC 到 其中 备份 文件 是 downloaded。",
						},
						"vpc_list": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Required:    true,
							Description: "Customize 列表 VPCs 到 download 备份 files。",
						},
					},
				},
			},

			"limit_ip": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "A 自定义 VPC IP 地址 对于 downloadable 备份 files.如果 参数 LimitType 是 **Customize**，您 need 到 configure 此 参数。",
			},

			"backup_infos": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "A 列表 备份 文件 信息。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"file_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Backup 文件 名称",
						},
						"file_size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "备份 文件 大小 是 在 单位 B，如果 它 是 0，它 是 无效。",
						},
						"download_url": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Backup 文件 download 地址 在 Internet (6 hours)。",
						},
						"inner_download_url": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Backup 文件 intranet download 地址 (6 hours)。",
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

func dataSourceTencentCloudRedisBackupDownloadInfoRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_redis_backup_download_info.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		paramMap["instance_id"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("backup_id"); ok {
		paramMap["backup_id"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("limit_type"); ok {
		paramMap["limit_type"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("vpc_comparison_symbol"); ok {
		paramMap["vpc_comparison_symbol"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("ip_comparison_symbol"); ok {
		paramMap["ip_comparison_symbol"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("limit_vpc"); ok {
		limitVpcSet := v.([]interface{})
		tmpSet := make([]*redis.BackupLimitVpcItem, 0, len(limitVpcSet))

		for _, item := range limitVpcSet {
			backupLimitVpcItem := redis.BackupLimitVpcItem{}
			backupLimitVpcItemMap := item.(map[string]interface{})

			if v, ok := backupLimitVpcItemMap["region"]; ok {
				backupLimitVpcItem.Region = helper.String(v.(string))
			}
			if v, ok := backupLimitVpcItemMap["vpc_list"]; ok {
				vpcListSet := v.(*schema.Set).List()
				backupLimitVpcItem.VpcList = helper.InterfacesStringsPoint(vpcListSet)
			}
			tmpSet = append(tmpSet, &backupLimitVpcItem)
		}
		paramMap["limit_vpc"] = tmpSet
	}

	if v, ok := d.GetOk("limit_ip"); ok {
		limitIpSet := v.(*schema.Set).List()
		paramMap["limit_ip"] = helper.InterfacesStringsPoint(limitIpSet)
	}

	service := RedisService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var backupInfos []*redis.BackupDownloadInfo

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeRedisBackupDownloadInfoByFilter(ctx, paramMap)
		if e != nil {
			if ee, ok := e.(*sdkErrors.TencentCloudSDKError); ok {
				if ee.Code == "FailedOperation.SystemError" {
					return resource.NonRetryableError(e)
				}
			}
			return tccommon.RetryError(e)
		}
		backupInfos = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(backupInfos))
	tmpList := make([]map[string]interface{}, 0, len(backupInfos))

	if backupInfos != nil {
		for _, backupDownloadInfo := range backupInfos {
			backupDownloadInfoMap := map[string]interface{}{}

			if backupDownloadInfo.FileName != nil {
				backupDownloadInfoMap["file_name"] = backupDownloadInfo.FileName
			}

			if backupDownloadInfo.FileSize != nil {
				backupDownloadInfoMap["file_size"] = backupDownloadInfo.FileSize
			}

			if backupDownloadInfo.DownloadUrl != nil {
				backupDownloadInfoMap["download_url"] = backupDownloadInfo.DownloadUrl
			}

			if backupDownloadInfo.InnerDownloadUrl != nil {
				backupDownloadInfoMap["inner_download_url"] = backupDownloadInfo.InnerDownloadUrl
			}

			ids = append(ids, *backupDownloadInfo.FileName)
			tmpList = append(tmpList, backupDownloadInfoMap)
		}

		_ = d.Set("backup_infos", tmpList)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
