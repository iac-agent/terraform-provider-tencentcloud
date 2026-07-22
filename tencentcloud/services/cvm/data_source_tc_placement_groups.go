package cvm

import (
	"context"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudPlacementGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudPlacementGroupsRead,

		Schema: map[string]*schema.Schema{
			"placement_group_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID placement group to be queried。",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "名称 placement group to be queried。",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "用于保存结果。",
			},

			"placement_group_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "An information 列表 placement group. Each element 包含following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"placement_group_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID placement group。",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "名称 placement group。",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "类型 placement group。",
						},
						"cvm_quota_total": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "最大hosts in the placement group。",
						},
						"current_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "数量 hosts in the placement group。",
						},
						"instance_ids": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "主机 IDs in the placement group。",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "创建时间 of the placement group。",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudPlacementGroupsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_placement_groups.read")()
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	cvmService := CvmService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}

	var placementGroupId string
	var name string
	if v, ok := d.GetOk("placement_group_id"); ok {
		placementGroupId = v.(string)
	}
	if v, ok := d.GetOk("name"); ok {
		name = v.(string)
	}

	var placementGroups []*cvm.DisasterRecoverGroup
	var errRet error
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		placementGroups, errRet = cvmService.DescribePlacementGroupByFilter(ctx, placementGroupId, name)
		if errRet != nil {
			return tccommon.RetryError(errRet, tccommon.InternalError)
		}
		return nil
	})
	if err != nil {
		return err
	}

	placementGroupList := make([]map[string]interface{}, 0, len(placementGroups))
	ids := make([]string, 0, len(placementGroups))
	for _, placement := range placementGroups {
		mapping := map[string]interface{}{
			"placement_group_id": placement.DisasterRecoverGroupId,
			"name":               placement.Name,
			"type":               placement.Type,
			"cvm_quota_total":    placement.CvmQuotaTotal,
			"current_num":        placement.CurrentNum,
			"instance_ids":       helper.StringsInterfaces(placement.InstanceIds),
			"create_time":        placement.CreateTime,
		}
		placementGroupList = append(placementGroupList, mapping)
		ids = append(ids, *placement.DisasterRecoverGroupId)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	err = d.Set("placement_group_list", placementGroupList)
	if err != nil {
		log.Printf("[CRITAL]%s provider set placement group list fail, reason:%s\n ", logId, err.Error())
		return err
	}

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if err := tccommon.WriteToFile(output.(string), placementGroupList); err != nil {
			return err
		}
	}
	return nil
}
