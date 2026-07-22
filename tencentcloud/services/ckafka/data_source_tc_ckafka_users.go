package ckafka

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudCkafkaUsers() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCkafkaUsersRead,

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID 的 ckafka 实例。",
			},
			"account_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "账号 名称 使用 当 查询 ckafka users' infos. Could 是 substr 的 用户 名称",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "用于保存结果。",
			},
			"user_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A 列表 ckafka users. Each element 包含following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"account_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "账号 名称 用户",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "创建时间 的 账号",
						},
						"update_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "last 更新时间 的 账号",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudCkafkaUsersRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_ckafka_users.read")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	params := make(map[string]interface{})
	params["instance_id"] = d.Get("instance_id").(string)
	if v, ok := d.GetOk("account_name"); ok {
		params["account_name"] = v.(string)
	}

	ckafkaService := CkafkaService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}
	userInfos, err := ckafkaService.DescribeUserByFilter(ctx, params)
	if err != nil {
		return err
	}
	userList := make([]map[string]interface{}, 0, len(userInfos))
	ids := make([]string, 0, len(userInfos))
	for _, user := range userInfos {
		userList = append(userList, map[string]interface{}{
			"account_name": *user.Name,
			"create_time":  *user.CreateTime,
			"update_time":  *user.UpdateTime,
		})

		ids = append(ids, params["instance_id"].(string)+tccommon.FILED_SP+*user.Name)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	_ = d.Set("user_list", userList)

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), userList); e != nil {
			return e
		}
	}

	return nil
}
