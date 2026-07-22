package chdfs

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	chdfs "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/chdfs/v20201112"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudChdfsAccessGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudChdfsAccessGroupsRead,
		Schema: map[string]*schema.Schema{
			"vpc_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "get groups belongs 到 私有网络 ID，必须 集合 但 仅 可以 使用 一个 的 VpcId 和 OwnerUin 到 get groups。",
			},

			"owner_uin": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "get groups belongs 到 所有者 uin，必须 集合 但 仅 可以 使用 一个 的 VpcId 和 OwnerUin 到 get groups。",
			},

			"access_groups": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "访问 组 列表。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"access_group_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "访问 组 ID。",
						},
						"access_group_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "访问 组名称",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "访问 组 描述",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "创建时间。",
						},
						"vpc_type": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "vpc 网络 类型(1:CVM，2:BM 1.0)。",
						},
						"vpc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "私有网络 ID",
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

func dataSourceTencentCloudChdfsAccessGroupsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_chdfs_access_groups.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("vpc_id"); ok {
		paramMap["vpc_id"] = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("owner_uin"); ok {
		paramMap["owner_uin"] = helper.IntUint64(v.(int))
	}

	service := ChdfsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var accessGroups []*chdfs.AccessGroup

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeChdfsAccessGroupsByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		accessGroups = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(accessGroups))
	tmpList := make([]map[string]interface{}, 0, len(accessGroups))

	if accessGroups != nil {
		for _, accessGroup := range accessGroups {
			accessGroupMap := map[string]interface{}{}

			if accessGroup.AccessGroupId != nil {
				accessGroupMap["access_group_id"] = accessGroup.AccessGroupId
			}

			if accessGroup.AccessGroupName != nil {
				accessGroupMap["access_group_name"] = accessGroup.AccessGroupName
			}

			if accessGroup.Description != nil {
				accessGroupMap["description"] = accessGroup.Description
			}

			if accessGroup.CreateTime != nil {
				accessGroupMap["create_time"] = accessGroup.CreateTime
			}

			if accessGroup.VpcType != nil {
				accessGroupMap["vpc_type"] = accessGroup.VpcType
			}

			if accessGroup.VpcId != nil {
				accessGroupMap["vpc_id"] = accessGroup.VpcId
			}

			ids = append(ids, *accessGroup.AccessGroupId)
			tmpList = append(tmpList, accessGroupMap)
		}

		_ = d.Set("access_groups", tmpList)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
