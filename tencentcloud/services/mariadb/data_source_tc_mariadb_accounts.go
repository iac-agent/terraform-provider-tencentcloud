package mariadb

import (
	"context"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mariadb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mariadb/v20170312"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudMariadbAccounts() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudMariadbAccountsRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "实例 ID",
			},

			"list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "账号 列表。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"user_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "用户名",
						},
						"host": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "主机 从 其中 用户 可以 日志 在 (corresponding 到 主机 字段 的 MySQL users，UserName + 主机 uniquely identifies 用户，在 form 的 IP，和 IP segment 结束 使用 %; 支持 filling 在 %; 如果 它 是 空，它 默认为 %)。",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "用户 备注",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "创建时间。",
						},
						"update_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "更新时间。",
						},
						"read_only": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Read-仅 flag，`0`: No，`1`: SQL 请求 的 此 账号 是 preferentially executed 在 standby machine，和 主机 machine 是 selected 对于 execution 当 standby machine 是 unavailable，`2`: standby machine 是 preferentially selected 对于 execution，和 operation fails 当 standby machine 是 unavailable。",
						},
						"delay_thresh": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "此 字段 是 meaningful 对于 read-仅 accounts，indicating 该 standby machine 使用 活跃-standby 延迟 less 比 此 值 是 selected。",
						},
						"slave_const": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "For read-仅 accounts，集合 是否policy 是 到 fix standby machine，`0`: standby machine 是 不 fixed，该 是， standby machine does 不 meet conditions 和 将 不 disconnect 从 客户端，和 Proxy selects other 可用 standby machines，`1`: standby machine does 不 meet conditions Disconnect，make sure 一个 连接 secures standby。",
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

func dataSourceTencentCloudMariadbAccountsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_mariadb_accounts.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	var instanceId string

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		paramMap["instance_id"] = helper.String(v.(string))
	}

	mariadbService := MariadbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var users []*mariadb.DBAccount
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		results, e := mariadbService.DescribeMariadbAccountsByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		users = results
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read Mariadb users failed, reason:%+v", logId, err)
		return err
	}

	userList := []interface{}{}
	if users != nil {
		for _, user := range users {
			userMap := map[string]interface{}{}
			if user.UserName != nil {
				userMap["user_name"] = user.UserName
			}
			if user.Host != nil {
				userMap["host"] = user.Host
			}
			if user.Description != nil {
				userMap["description"] = user.Description
			}
			if user.CreateTime != nil {
				userMap["create_time"] = user.CreateTime
			}
			if user.UpdateTime != nil {
				userMap["update_time"] = user.UpdateTime
			}
			if user.ReadOnly != nil {
				userMap["read_only"] = user.ReadOnly
			}
			if user.DelayThresh != nil {
				userMap["delay_thresh"] = user.DelayThresh
			}
			if user.SlaveConst != nil {
				userMap["slave_const"] = user.SlaveConst
			}

			userList = append(userList, userMap)
		}
		err = d.Set("list", userList)
		if err != nil {
			log.Printf("[CRITAL]%s provider set instances list fail, reason:%s\n ", logId, err.Error())
			return err
		}
	}
	d.SetId(instanceId)

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), userList); e != nil {
			return e
		}
	}

	return nil
}
