package cfs

import (
	"context"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cfs "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cfs/v20190719"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudCfsAccessRules() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCfsAccessRulesRead,

		Schema: map[string]*schema.Schema{
			"access_group_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "A 指定 访问 组 ID 用于query。",
			},
			"access_rule_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A 指定 访问 规则 ID 用于query。",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "用于保存结果。",
			},

			"access_rule_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "An 信息 列表 CFS 访问 规则. Each element 包含following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"access_rule_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID 访问 规则。",
						},
						"auth_client_ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Allowed IP 的 访问 规则。",
						},
						"rw_permission": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Read 和 write permissions。",
						},
						"user_permission": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "permissions 的 accessing users。",
						},
						"priority": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "优先级 级别 的 访问 规则。",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudCfsAccessRulesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_cfs_access_rules.read")()
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	cfsService := CfsService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}

	var accessRuleId string
	accessGroupId := d.Get("access_group_id").(string)
	if v, ok := d.GetOk("access_rule_id"); ok {
		accessRuleId = v.(string)
	}

	var accessRules []*cfs.PGroupRuleInfo
	var errRet error
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		accessRules, errRet = cfsService.DescribeAccessRule(ctx, accessGroupId, accessRuleId)
		if errRet != nil {
			return tccommon.RetryError(errRet, tccommon.InternalError)
		}
		return nil
	})
	if err != nil {
		return err
	}

	accessRuleList := make([]map[string]interface{}, 0, len(accessRules))
	ids := make([]string, 0, len(accessRules))
	for _, accessRule := range accessRules {
		mapping := map[string]interface{}{
			"access_rule_id":  accessRule.RuleId,
			"auth_client_ip":  accessRule.AuthClientIp,
			"rw_permission":   accessRule.RWPermission,
			"user_permission": accessRule.UserPermission,
			"priority":        accessRule.Priority,
		}
		accessRuleList = append(accessRuleList, mapping)
		ids = append(ids, *accessRule.RuleId)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	err = d.Set("access_rule_list", accessRuleList)
	if err != nil {
		log.Printf("[CRITAL]%s provider set cfs access rule list fail, reason:%s\n ", logId, err.Error())
		return err
	}
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if err := tccommon.WriteToFile(output.(string), accessRuleList); err != nil {
			return err
		}
	}
	return nil
}
