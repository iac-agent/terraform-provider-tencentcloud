package tcaplusdb

import (
	"context"
	"fmt"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceTencentCloudTcaplusTableGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTcaplusTableGroupsRead,
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID 的 TcaplusDB 集群 到 是 查询。",
			},
			"tablegroup_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID 的 表 组 到 是 查询。",
			},
			"tablegroup_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "名称 表 组 到 是 查询。",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "File 对于 saving results。",
			},
			"list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A 列表 表 组. Each element 包含following attributes。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tablegroup_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "名称 表 组。",
						},
						"tablegroup_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID 的 表 组。",
						},
						"table_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "数量 tables。",
						},
						"total_size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Total 存储 大小 (MB)。",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "创建时间 的 表 组.。",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudTcaplusTableGroupsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_tcaplus_tablegroups.read")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := TcaplusService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}

	clusterId := d.Get("cluster_id").(string)
	groupId := d.Get("tablegroup_id").(string)
	groupName := d.Get("tablegroup_name").(string)

	groups, err := service.DescribeGroups(ctx, clusterId, groupId, groupName)
	if err != nil {
		groups, err = service.DescribeGroups(ctx, clusterId, groupId, groupName)
	}

	if err != nil {
		return err
	}

	list := make([]map[string]interface{}, 0, len(groups))

	for _, group := range groups {
		listItem := make(map[string]interface{})
		listItem["tablegroup_name"] = group.TableGroupName
		listItem["tablegroup_id"] = fmt.Sprintf("%s:%s", clusterId, *group.TableGroupId)
		listItem["table_count"] = group.TableCount
		listItem["total_size"] = group.TotalSize
		listItem["create_time"] = group.CreatedTime
		list = append(list, listItem)
	}

	d.SetId("group." + clusterId + "." + groupId + "." + groupName)
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
