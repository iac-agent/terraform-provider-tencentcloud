package monitor

import (
	"context"
	"fmt"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	monitor "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/monitor/v20180724"
)

func DataSourceTencentCloudMonitorBindingObjects() *schema.Resource {

	return &schema.Resource{
		Read: dataSourceTencentMonitorBindingObjectRead,
		Schema: map[string]*schema.Schema{
			"group_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Policy 组 ID 对于 查询。",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "用于存储结果。",
			},
			// Computed values
			"list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A 列表 objects. Each element 包含following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"unique_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Object 唯一 ID。",
						},
						"dimensions_json": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Represents collection 的 dimensions 的 对象 实例，json 格式",
						},
						"is_shielded": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "是否object 是 shielded 或 不，`0` 表示 unshielded 和 `1` 表示 shielded。",
						},
						"region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "地域 其中 对象 是 located。",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentMonitorBindingObjectRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_monitor_binding_objects.read")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	var (
		monitorService = MonitorService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		objects        []*monitor.DescribeBindingPolicyObjectListInstance
		err            error

		list = make([]interface{}, 0, len(objects))
	)

	id := int64(d.Get("group_id").(int))

	objects, err = monitorService.DescribeBindingPolicyObjectList(ctx, id)
	if err != nil {
		return err
	}

	for _, event := range objects {
		var listItem = map[string]interface{}{}
		listItem["region"] = event.Region
		listItem["unique_id"] = event.UniqueId
		listItem["dimensions_json"] = event.Dimensions
		listItem["is_shielded"] = event.IsShielded
		listItem["region"] = event.Region
		list = append(list, listItem)
	}
	if err = d.Set("list", list); err != nil {
		return err
	}
	d.SetId(fmt.Sprintf("%d", id))
	if output, ok := d.GetOk("result_output_file"); ok {
		return tccommon.WriteToFile(output.(string), list)
	}
	return nil
}
