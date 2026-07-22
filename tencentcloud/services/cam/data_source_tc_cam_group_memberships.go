package cam

import (
	"context"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudCamGroupMemberships() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCamGroupMembershipsRead,

		Schema: map[string]*schema.Schema{
			"group_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID CAM 组 到 是 queried。",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "用于保存结果。",
			},
			"membership_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A 列表 CAM 组 membership. Each element 包含following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"group_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID CAM 组。",
						},
						"user_ids": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Deprecated:  "It has been deprecated from version 1.59.5. Use `user_names` instead.",
							Description: "ID 集合 的 CAM 组 members。",
						},
						"user_names": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Description: "ID 集合 的 CAM 组 members。",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudCamGroupMembershipsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_cam_group_memberships.read")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	groupId := d.Get("group_id").(string)
	camService := CamService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}
	var memberships []*string
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		results, e := camService.DescribeGroupMembershipById(ctx, groupId)
		if e != nil {
			return tccommon.RetryError(e)
		}
		memberships = results
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read CAM group memberships failed, reason:%s\n", logId, err.Error())
		return err
	}
	groupList := make([]map[string]interface{}, 0, 1)
	ids := make([]string, 0, 1)
	mapping := map[string]interface{}{
		"group_id":   groupId,
		"user_ids":   memberships,
		"user_names": memberships,
	}
	groupList = append(groupList, mapping)
	ids = append(ids, groupId)

	d.SetId(helper.DataResourceIdsHash(ids))
	if e := d.Set("membership_list", groupList); e != nil {
		log.Printf("[CRITAL]%s provider set membership list fail, reason:%s\n", logId, e.Error())
		return e
	}

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), groupList); e != nil {
			return e
		}
	}

	return nil
}
