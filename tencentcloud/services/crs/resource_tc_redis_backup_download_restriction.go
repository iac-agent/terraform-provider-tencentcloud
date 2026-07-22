package crs

import (
	"context"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	redis "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/redis/v20180412"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudRedisBackupDownloadRestriction() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudRedisBackupDownloadRestrictionCreate,
		Read:   resourceTencentCloudRedisBackupDownloadRestrictionRead,
		Update: resourceTencentCloudRedisBackupDownloadRestrictionUpdate,
		Delete: resourceTencentCloudRedisBackupDownloadRestrictionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"limit_type": {
				Required:    true,
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
		},
	}
}

func resourceTencentCloudRedisBackupDownloadRestrictionCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_redis_backup_download_restriction.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	region := meta.(tccommon.ProviderMeta).GetAPIV3Conn().Region

	d.SetId(region)

	return resourceTencentCloudRedisBackupDownloadRestrictionUpdate(d, meta)
}

func resourceTencentCloudRedisBackupDownloadRestrictionRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_redis_backup_download_restriction.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := RedisService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	backupDownloadRestriction, err := service.DescribeRedisBackupDownloadRestrictionById(ctx)
	if err != nil {
		return err
	}

	if backupDownloadRestriction == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `RedisBackupDownloadRestriction` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if backupDownloadRestriction.LimitType != nil {
		_ = d.Set("limit_type", backupDownloadRestriction.LimitType)
	}

	if backupDownloadRestriction.VpcComparisonSymbol != nil {
		_ = d.Set("vpc_comparison_symbol", backupDownloadRestriction.VpcComparisonSymbol)
	}

	if backupDownloadRestriction.IpComparisonSymbol != nil {
		_ = d.Set("ip_comparison_symbol", backupDownloadRestriction.IpComparisonSymbol)
	}

	if backupDownloadRestriction.LimitVpc != nil {
		limitVpcList := []interface{}{}
		for _, limitVpc := range backupDownloadRestriction.LimitVpc {
			limitVpcMap := map[string]interface{}{}

			if limitVpc.Region != nil {
				limitVpcMap["region"] = limitVpc.Region
			}

			if limitVpc.VpcList != nil {
				limitVpcMap["vpc_list"] = limitVpc.VpcList
			}

			limitVpcList = append(limitVpcList, limitVpcMap)
		}

		_ = d.Set("limit_vpc", limitVpcList)

	}

	if backupDownloadRestriction.LimitIp != nil {
		_ = d.Set("limit_ip", backupDownloadRestriction.LimitIp)
	}

	return nil
}

func resourceTencentCloudRedisBackupDownloadRestrictionUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_redis_backup_download_restriction.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	request := redis.NewModifyBackupDownloadRestrictionRequest()

	if v, ok := d.GetOk("limit_type"); ok {
		request.LimitType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("vpc_comparison_symbol"); ok {
		request.VpcComparisonSymbol = helper.String(v.(string))
	}

	if v, ok := d.GetOk("ip_comparison_symbol"); ok {
		request.IpComparisonSymbol = helper.String(v.(string))
	}

	if v, ok := d.GetOk("limit_vpc"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			backupLimitVpcItem := redis.BackupLimitVpcItem{}
			if v, ok := dMap["region"]; ok {
				backupLimitVpcItem.Region = helper.String(v.(string))
			}
			if v, ok := dMap["vpc_list"]; ok {
				vpcListSet := v.(*schema.Set).List()
				for i := range vpcListSet {
					vpcList := vpcListSet[i].(string)
					backupLimitVpcItem.VpcList = append(backupLimitVpcItem.VpcList, &vpcList)
				}
			}
			request.LimitVpc = append(request.LimitVpc, &backupLimitVpcItem)
		}
	}

	if d.HasChange("limit_ip") {
		if v, ok := d.GetOk("limit_ip"); ok {
			limitIpSet := v.(*schema.Set).List()
			for i := range limitIpSet {
				limitIp := limitIpSet[i].(string)
				request.LimitIp = append(request.LimitIp, &limitIp)
			}
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseRedisClient().ModifyBackupDownloadRestriction(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update redis backupDownloadRestriction failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudRedisBackupDownloadRestrictionRead(d, meta)
}

func resourceTencentCloudRedisBackupDownloadRestrictionDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_redis_backup_download_restriction.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
