package cam

import (
	"context"
	"log"
	"strconv"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cam "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cam/v20190116"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudCamUserPolicyAttachments() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCamUserPolicyAttachmentsRead,

		Schema: map[string]*schema.Schema{
			"user_id": {
				Type:         schema.TypeString,
				Optional:     true,
				AtLeastOneOf: []string{"user_id", "user_name"},
				Deprecated:   "It has been deprecated from version 1.59.6. Use `user_name` instead.",
				Description:  "ID attached CAM 用户 to be queried。",
			},
			"user_name": {
				Type:         schema.TypeString,
				Optional:     true,
				AtLeastOneOf: []string{"user_id", "user_name"},
				Description:  "名称 attached CAM 用户 as unique 键 to be queried。",
			},
			"policy_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID CAM policy to be queried。",
			},
			"create_mode": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: tccommon.ValidateAllowedIntValue([]int{1, 2}),
				Description:  "模式 of Creation of the CAM 用户 policy attachment. `1` means the CAM policy attachment is created by production，and the others indicate syntax strategy ways。",
			},
			"policy_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue(CAM_POLICY_CREATE_STRATEGY),
				Description:  "类型 policy strategy. '用户' means customer strategy and 'QCS' means preset strategy。",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "用于保存结果。",
			},
			"user_policy_attachment_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A 列表 CAM 用户 policy attachments. Each element 包含following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"user_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Deprecated:  "It has been deprecated from version 1.59.6. Use `user_name` instead.",
							Description: "ID CAM 用户",
						},
						"user_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "名称 CAM 用户 as unique 键",
						},
						"policy_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "名称 CAM 用户",
						},
						"create_mode": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "模式 of Creation of the CAM 用户 policy attachment. `1` means the cam policy attachment is created by production，and the others indicate syntax strategy ways。",
						},
						"policy_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "类型 policy strategy. '用户' means customer strategy and 'QCS' means preset strategy。",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The 创建时间 of the CAM 用户 policy attachment。",
						},
						"policy_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The 名称 policy。",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudCamUserPolicyAttachmentsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_cam_user_policy_attachments.read")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	params := make(map[string]interface{})
	userId, _, err := getUserId(d)
	if err != nil {
		return err
	}
	params["user_id"] = userId
	if v, ok := d.GetOk("policy_id"); ok {
		policyId, err := strconv.Atoi(v.(string))
		if err != nil {
			return err
		}
		params["policy_id"] = uint64(policyId)
	}
	if v, ok := d.GetOk("policy_type"); ok {
		params["policy_type"] = v.(string)
	}
	if v, ok := d.GetOk("create_mode"); ok {
		params["create_mode"] = v.(int)
	}

	camService := CamService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}
	var policyOfUsers []*cam.AttachPolicyInfo
	err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		results, e := camService.DescribeUserPolicyAttachmentsByFilter(ctx, params)
		if e != nil {
			return tccommon.RetryError(e)
		}
		policyOfUsers = results
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read CAM user policy attachments failed, reason:%s\n", logId, err.Error())
		return err
	}
	policyOfUserList := make([]map[string]interface{}, 0, len(policyOfUsers))
	ids := make([]string, 0, len(policyOfUsers))
	for _, policy := range policyOfUsers {
		mapping := map[string]interface{}{
			"user_id":     userId,
			"user_name":   userId,
			"policy_id":   strconv.Itoa(int(*policy.PolicyId)),
			"create_time": *policy.AddTime,
			"create_mode": *policy.CreateMode,
			"policy_type": *policy.PolicyType,
			"policy_name": *policy.PolicyName,
		}
		policyOfUserList = append(policyOfUserList, mapping)
		ids = append(ids, userId+"#"+strconv.Itoa(int(*policy.PolicyId)))
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	if e := d.Set("user_policy_attachment_list", policyOfUserList); e != nil {
		log.Printf("[CRITAL]%s provider set CAM user polilcy attachment list fail, reason:%s\n", logId, e.Error())
		return e
	}

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), policyOfUserList); e != nil {
			return e
		}
	}

	return nil
}
