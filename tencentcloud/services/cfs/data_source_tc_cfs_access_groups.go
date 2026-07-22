package cfs

import (
	"context"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cfs "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cfs/v20190719"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudCfsAccessGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCfsAccessGroupsRead,

		Schema: map[string]*schema.Schema{
			"access_group_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A 指定 访问 组 ID 用于query。",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A 访问 组 名称 用于query。",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "用于保存结果。",
			},

			"access_group_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "An 信息 列表 CFS 访问 组. Each element 包含following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"access_group_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID 访问 组。",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "名称 访问 组。",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "描述 访问 组。",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "创建时间 的 访问 组。",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudCfsAccessGroupsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_cfs_access_groups.read")()
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	cfsService := CfsService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}

	var accessGroupId string
	var name string
	if v, ok := d.GetOk("access_group_id"); ok {
		accessGroupId = v.(string)
	}
	if v, ok := d.GetOk("name"); ok {
		name = v.(string)
	}

	var accessGroups []*cfs.PGroupInfo
	var errRet error
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		accessGroups, errRet = cfsService.DescribeAccessGroup(ctx, accessGroupId, name)
		if errRet != nil {
			return tccommon.RetryError(errRet, tccommon.InternalError)
		}
		return nil
	})
	if err != nil {
		return err
	}

	accessGroupList := make([]map[string]interface{}, 0, len(accessGroups))
	ids := make([]string, 0, len(accessGroups))
	for _, accessGroup := range accessGroups {
		mapping := map[string]interface{}{
			"access_group_id": accessGroup.PGroupId,
			"name":            accessGroup.Name,
			"description":     accessGroup.DescInfo,
			"create_time":     accessGroup.CDate,
		}
		accessGroupList = append(accessGroupList, mapping)
		ids = append(ids, *accessGroup.PGroupId)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	err = d.Set("access_group_list", accessGroupList)
	if err != nil {
		log.Printf("[CRITAL]%s provider set cfs access group list fail, reason:%s\n ", logId, err.Error())
		return err
	}
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if err := tccommon.WriteToFile(output.(string), accessGroupList); err != nil {
			return err
		}
	}
	return nil
}
