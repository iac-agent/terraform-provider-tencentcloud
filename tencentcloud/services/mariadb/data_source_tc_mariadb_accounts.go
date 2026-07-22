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
				Description: "账号 list。",
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
							Description: "The 主机 from which the 用户 can log in (corresponding to the 主机 field of MySQL users，UserName + 主机 uniquely identifies a 用户，in the form of IP，and the IP segment ends with %; supports filling in %; if it is empty，it 默认为 %)。",
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
							Description: "Read-only flag，`0`: No，`1`: The SQL request of this 账号 is preferentially executed on the standby machine，and the 主机 machine is selected for execution when the standby machine is unavailable，`2`: The standby machine is preferentially selected for execution，and the operation fails when the standby machine is unavailable。",
						},
						"delay_thresh": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "This field is meaningful for read-only accounts，indicating that the standby machine with the 活跃-standby 延迟 less than this 值 is selected。",
						},
						"slave_const": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "For read-only accounts，set 是否policy is to fix the standby machine，`0`: The standby machine is not fixed，that is，the standby machine does not meet the conditions and will not disconnect from the client，and the Proxy selects other available standby machines，`1`: The standby machine does not meet the conditions Disconnect，make sure one connection secures the standby。",
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
