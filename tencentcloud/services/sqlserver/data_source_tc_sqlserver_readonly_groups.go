package sqlserver

import (
	"context"
	"fmt"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudSqlserverReadonlyGroups() *schema.Resource {

	return &schema.Resource{
		Read: dataSourceTencentSqlserverReadonlyGroups,
		Schema: map[string]*schema.Schema{
			"master_instance_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Master SQL Server 实例 ID.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used 到 store results.",
			},
			// Computed values
			"list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A 列表 的 SQL Server readonly 组. Each element contains following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID 的 readonly 组.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name 的 readonly 组.",
						},
						"master_instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Master 实例 ID.",
						},
						"max_delay_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Maximum delay 时间 的 readonly 实例.",
						},
						"is_offline_delay": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Indicate whether 到 offline delayed readonly 实例.",
						},
						"min_instances": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Minimum readonly 实例 该 stays 在 组.",
						},
						"vip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Virtual IP 地址 的 readonly 组.",
						},
						"vport": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Virtual 端口 的 readonly 组.",
						},
						"readonly_instance_set": {
							Type:        schema.TypeList,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Computed:    true,
							Description: "Readonly 实例 ID 集合 的 readonly 组.",
						},
						"status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Status 的 readonly 组. `1` 对于 running, `5` 对于 applying.",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentSqlserverReadonlyGroups(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_sqlserver_readonly_groups.read")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	instanceId := d.Get("master_instance_id").(string)
	sqlserverService := SqlserverService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	groupList, err := sqlserverService.DescribeReadonlyGroupList(ctx, instanceId)

	if err != nil {
		return fmt.Errorf("api[DescribeReadOnlyGroupList]fail, return %s", err.Error())
	}

	var list []map[string]interface{}
	var ids = make([]string, len(groupList))

	for _, item := range groupList {
		roSet := make([]string, 0)
		for _, v := range item.ReadOnlyInstanceSet {
			roSet = append(roSet, *v.InstanceId)
		}
		mapping := map[string]interface{}{
			"name":                  item.ReadOnlyGroupName,
			"vip":                   item.Vip,
			"vport":                 item.Vport,
			"is_offline_delay":      item.IsOfflineDelay,
			"max_delay_time":        item.ReadOnlyMaxDelayTime,
			"min_instances":         item.MinReadOnlyInGroup,
			"status":                item.Status,
			"master_instance_id":    item.MasterInstanceId,
			"id":                    item.ReadOnlyGroupId,
			"readonly_instance_set": roSet,
		}
		list = append(list, mapping)
		ids = append(ids, *item.ReadOnlyGroupId)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	if e := d.Set("list", list); e != nil {
		log.Printf("[CRITAL]%s provider set list fail, reason:%s\n", logId, e.Error())
		return e
	}
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		return tccommon.WriteToFile(output.(string), list)
	}

	return nil
}
