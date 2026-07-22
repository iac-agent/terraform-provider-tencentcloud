package cbs

import (
	"context"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cbs "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cbs/v20170312"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudCbsSnapshotPolicies() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCbsSnapshotPoliciesRead,

		Schema: map[string]*schema.Schema{
			"snapshot_policy_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID 快照 策略 到 是 queried。",
			},
			"snapshot_policy_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "名称 快照 策略 到 是 queried。",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "用于保存结果。",
			},
			"snapshot_policy_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A 列表 快照 策略. Each element 包含following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"snapshot_policy_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID 快照 策略。",
						},
						"snapshot_policy_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "名称 快照 策略。",
						},
						"repeat_weekdays": {
							Type:        schema.TypeList,
							Elem:        &schema.Schema{Type: schema.TypeInt},
							Computed:    true,
							Description: "Trigger days 的 periodic 快照。",
						},
						"repeat_hours": {
							Type:        schema.TypeList,
							Elem:        &schema.Schema{Type: schema.TypeInt},
							Computed:    true,
							Description: "Trigger hours 的 periodic 快照。",
						},
						"retention_days": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Retention days 的 快照。",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "状态 快照 策略。",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "创建时间 的 快照 策略。",
						},
						"attached_storage_ids": {
							Type:        schema.TypeList,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Computed:    true,
							Description: "Storage IDs 该 快照 策略 attached。",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudCbsSnapshotPoliciesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_cbs_snapshot_policies.read")()
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	var policyId string
	var policyName string
	if v, ok := d.GetOk("snapshot_policy_id"); ok {
		policyId = v.(string)
	}
	if v, ok := d.GetOk("snapshot_policy_name"); ok {
		policyName = v.(string)
	}
	cbsService := CbsService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}
	var policies []*cbs.AutoSnapshotPolicy
	var errRet error
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		policies, errRet = cbsService.DescribeSnapshotPolicy(ctx, policyId, policyName)
		if errRet != nil {
			return tccommon.RetryError(errRet, tccommon.InternalError)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read cbs snapshot policies failed, reason:%s\n ", logId, err.Error())
		return err
	}

	ids := make([]string, 0, len(policies))
	policyList := make([]map[string]interface{}, 0, len(policies))
	for _, policy := range policies {
		mapping := map[string]interface{}{
			"snapshot_policy_id":   policy.AutoSnapshotPolicyId,
			"snapshot_policy_name": policy.AutoSnapshotPolicyName,
			"retention_days":       policy.RetentionDays,
			"status":               policy.AutoSnapshotPolicyState,
			"create_time":          policy.CreateTime,
			"attached_storage_ids": helper.StringsInterfaces(policy.DiskIdSet),
		}
		if len(policy.Policy) < 1 {
			continue
		}
		mapping["repeat_weekdays"] = helper.Uint64sInterfaces(policy.Policy[0].DayOfWeek)
		mapping["repeat_hours"] = helper.Uint64sInterfaces(policy.Policy[0].Hour)
		policyList = append(policyList, mapping)
		ids = append(ids, *policy.AutoSnapshotPolicyId)
	}
	d.SetId(helper.DataResourceIdsHash(ids))
	if err = d.Set("snapshot_policy_list", policyList); err != nil {
		log.Printf("[CRITAL]%s provider set snapshot policy list fail, reason:%s\n ", logId, err.Error())
		return err
	}

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if err = tccommon.WriteToFile(output.(string), policyList); err != nil {
			return err
		}
	}
	return nil
}
