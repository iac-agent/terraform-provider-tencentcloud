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
				Description: "Cloud 数据库 账号 信息。",
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
							Description: "From 其中 主机 用户 可以 日志 在 (corresponding 到 主机 字段 的 MySQL users，UserName + 主机 uniquely identifies 用户，在 form 的 IP， IP segment 结束 使用 %; 支持 filling 在 %; 如果 它 是 空，它 默认为 %)。",
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
							Description: "Read-仅 flag，0: No，1: SQL 请求 的 此 账号 是 preferentially executed 在 standby machine，和 主机 是 selected 对于 execution 当 standby machine 是 unavailable. 2: standby machine 是 preferentially selected 对于 execution，和 operation fails 当 standby machine 是 unavailable。",
						},
						"delay_thresh": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "如果 standby machine 延迟 exceeds setting 值 的 此 参数， 系统 将 consider 该 standby machine 是 faulty 和 recommend 该 参数 值 是 greater 比 10. 此 参数 takes effect 当 ReadOnly selects 1 和 2。",
						},
						"slave_const": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "For read-仅 accounts，集合 策略 是否fix standby machine，0: 不 fix standby machine，该 是， standby machine 将 不 disconnect 从 客户端 如果 它 does 不 meet conditions， Proxy selects other 可用 standby machines，1: standby machine 将 是 disconnected 如果 conditions 是 不 met，Make sure 连接 是 secured 到 standby machine。",
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
