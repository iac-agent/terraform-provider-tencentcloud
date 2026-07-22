package dcdb

import (
	"context"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dcdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dcdb/v20180411"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudDcdbAccounts() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudDcdbAccountsRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "实例 ID",
			},
			"list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Cloud database 账号 information。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"user_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "用户 名称",
						},
						"host": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "From which 主机 the 用户 can log in (corresponding to the 主机 field of MySQL users，UserName + 主机 uniquely identifies a 用户，in the form of IP，the IP segment ends with %; supports filling in %; if it is empty，it 默认为 %)。",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "用户 备注 info。",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "创建时间。",
						},
						"update_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Last 更新时间。",
						},
						"read_only": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Read-only flag，0: No，1: The SQL request of this 账号 is preferentially executed on the standby machine，and the 主机 is selected for execution when the standby machine is unavailable. 2: The standby machine is preferentially selected for execution，and the operation fails when the standby machine is unavailable。",
						},
						"delay_thresh": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "If the standby machine 延迟 exceeds the setting 值 of this parameter，the system will consider that the standby machine is faulty and recommend that the parameter 值 be greater than 10. This parameter takes effect when ReadOnly selects 1 and 2。",
						},
						"slave_const": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "For read-only accounts，set the policy 是否fix the standby machine，0: not fix the standby machine，that is，the standby machine will not disconnect from the client if it does not meet the conditions，the Proxy selects other available standby machines，1: the standby machine will be disconnected if the conditions are not met，Make sure a connection is secured to the standby machine。",
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

func dataSourceTencentCloudDcdbAccountsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_dcdb_accounts.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		paramMap["instance_id"] = helper.String(v.(string))
	}

	dcdbService := DcdbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var dbAccountList []*dcdb.DBAccount
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		results, e := dcdbService.DescribeDcdbAccountsByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		dbAccountList = results
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read Dcdb list failed, reason:%+v", logId, err)
		return err
	}

	retList := []interface{}{}
	if dbAccountList != nil {
		ids := make([]string, 0, len(dbAccountList))
		for _, dbA := range dbAccountList {
			listMap := map[string]interface{}{}
			if dbA.UserName != nil {
				listMap["user_name"] = dbA.UserName
			}
			if dbA.Host != nil {
				listMap["host"] = dbA.Host
			}
			if dbA.Description != nil {
				listMap["description"] = dbA.Description
			}
			if dbA.CreateTime != nil {
				listMap["create_time"] = dbA.CreateTime
			}
			if dbA.UpdateTime != nil {
				listMap["update_time"] = dbA.UpdateTime
			}
			if dbA.ReadOnly != nil {
				listMap["read_only"] = dbA.ReadOnly
			}
			if dbA.DelayThresh != nil {
				listMap["delay_thresh"] = dbA.DelayThresh
			}
			if dbA.SlaveConst != nil {
				listMap["slave_const"] = dbA.SlaveConst
			}
			ids = append(ids, *dbA.UserName+tccommon.FILED_SP+*dbA.Host)
			retList = append(retList, listMap)
		}

		d.SetId(helper.DataResourceIdsHash(ids))
		_ = d.Set("list", retList)
	}

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), retList); e != nil {
			return e
		}
	}

	return nil
}
