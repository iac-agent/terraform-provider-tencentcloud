package cam

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cam "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cam/v20190116"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudCamListEntitiesForPolicy() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCamListEntitiesForPolicyRead,
		Schema: map[string]*schema.Schema{
			"policy_id": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Policy ID。",
			},

			"rp": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Per 每页数量，默认值为 20。",
			},

			"entity_filter": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Can take 值 的 &amp;amp;#39;All&amp;amp;#39;，&amp;amp;#39;用户&amp;amp;#39;，&amp;amp;#39;Group&amp;amp;#39;，和 &amp;amp;#39;角色&amp;amp;#39;. &amp;amp;#39;All&amp;amp;#39; 表示 obtaining all entity types，&amp;amp;#39;用户&amp;amp;#39; 表示 仅 obtaining sub accounts，&amp;amp;#39;Group&amp;amp;#39; 表示 仅 obtaining 用户 groups，和 &amp;amp;#39;角色&amp;amp;#39; 表示 仅 obtaining roles. 默认值 是&amp;amp;#39; All &amp;amp;#39;。",
			},

			"list": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Entity List注意：此字段可能返回 null，表示无法获取有效值。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Entity ID。",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Entity Name注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"uin": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Entity Uin注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"related_type": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Association 类型 1. 用户 association; 2 用户 Group Association。",
						},
						"attachment_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Policy association time注意：此字段可能返回 null，表示无法获取有效值。",
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

func dataSourceTencentCloudCamListEntitiesForPolicyRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_cam_list_entities_for_policy.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, _ := d.GetOkExists("policy_id"); v != nil {
		paramMap["PolicyId"] = helper.IntUint64(v.(int))
	}

	if v, _ := d.GetOkExists("rp"); v != nil {
		paramMap["Rp"] = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("entity_filter"); ok {
		paramMap["EntityFilter"] = helper.String(v.(string))
	}

	service := CamService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var listEntitiesForPolicy []*cam.AttachEntityOfPolicy
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeCamListEntitiesForPolicyByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		listEntitiesForPolicy = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(listEntitiesForPolicy))
	tmpList := make([]map[string]interface{}, 0)

	if listEntitiesForPolicy != nil {
		for _, attachEntityOfPolicy := range listEntitiesForPolicy {
			attachEntityOfPolicyMap := map[string]interface{}{}

			if attachEntityOfPolicy.Id != nil {
				attachEntityOfPolicyMap["id"] = attachEntityOfPolicy.Id
			}

			if attachEntityOfPolicy.Name != nil {
				attachEntityOfPolicyMap["name"] = attachEntityOfPolicy.Name
			}

			if attachEntityOfPolicy.Uin != nil {
				attachEntityOfPolicyMap["uin"] = attachEntityOfPolicy.Uin
			}

			if attachEntityOfPolicy.RelatedType != nil {
				attachEntityOfPolicyMap["related_type"] = attachEntityOfPolicy.RelatedType
			}

			if attachEntityOfPolicy.AttachmentTime != nil {
				attachEntityOfPolicyMap["attachment_time"] = attachEntityOfPolicy.AttachmentTime
			}

			ids = append(ids, helper.UInt64ToStr(*attachEntityOfPolicy.Uin))
			tmpList = append(tmpList, attachEntityOfPolicyMap)
		}

		_ = d.Set("list", tmpList)
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
