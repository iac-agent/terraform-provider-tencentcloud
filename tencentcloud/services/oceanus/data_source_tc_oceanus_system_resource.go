package oceanus

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	oceanus "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/oceanus/v20190422"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudOceanusSystemResource() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudOceanusSystemResourceRead,
		Schema: map[string]*schema.Schema{
			"resource_ids": {
				Optional:    true,
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "数组 资源 IDs 到 是 queried。",
			},
			"filters": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "Query 资源 配置 列表. 如果未指定，返回 all 作业 配置 lists under ResourceIds.N。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "待过滤字段",
						},
						"values": {
							Type:        schema.TypeSet,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Required:    true,
							Description: "过滤器 值 对于 字段。",
						},
					},
				},
			},
			"cluster_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "集群 ID",
			},
			"flink_version": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Query built-在 connectors 对于 corresponding Flink 版本",
			},
			"resource_set": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Collection 的 资源 details。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"resource_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "资源 ID",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "资源名称",
						},
						"resource_type": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "资源类型 1 表示JAR 包，其中 是 currently 仅 支持 值",
						},
						"remark": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Resource 备注",
						},
						"region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "地域 到 其中 资源 belongs。",
						},
						"latest_resource_config_version": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Latest 版本 的 资源。",
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

func dataSourceTencentCloudOceanusSystemResourceRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_oceanus_system_resource.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId              = tccommon.GetLogId(tccommon.ContextNil)
		ctx                = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service            = OceanusService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		systemResourceList []*oceanus.SystemResourceItem
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("resource_ids"); ok {
		resourceIdsSet := v.(*schema.Set).List()
		paramMap["ResourceIds"] = helper.InterfacesStringsPoint(resourceIdsSet)
	}

	if v, ok := d.GetOk("filters"); ok {
		filtersSet := v.([]interface{})
		tmpSet := make([]*oceanus.Filter, 0, len(filtersSet))

		for _, item := range filtersSet {
			filter := oceanus.Filter{}
			filterMap := item.(map[string]interface{})

			if v, ok := filterMap["name"]; ok {
				filter.Name = helper.String(v.(string))
			}

			if v, ok := filterMap["values"]; ok {
				valuesSet := v.(*schema.Set).List()
				filter.Values = helper.InterfacesStringsPoint(valuesSet)
			}

			tmpSet = append(tmpSet, &filter)
		}

		paramMap["Filters"] = tmpSet
	}

	if v, ok := d.GetOk("cluster_id"); ok {
		paramMap["ClusterId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("flink_version"); ok {
		paramMap["FlinkVersion"] = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeOceanusSystemResourceByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		systemResourceList = result
		return nil
	})

	if err != nil {
		return err
	}

	ids := make([]string, 0, len(systemResourceList))
	tmpList := make([]map[string]interface{}, 0, len(systemResourceList))
	if systemResourceList != nil {
		for _, systemResourceItem := range systemResourceList {
			systemResourceItemMap := map[string]interface{}{}

			if systemResourceItem.ResourceId != nil {
				systemResourceItemMap["resource_id"] = systemResourceItem.ResourceId
			}

			if systemResourceItem.Name != nil {
				systemResourceItemMap["name"] = systemResourceItem.Name
			}

			if systemResourceItem.ResourceType != nil {
				systemResourceItemMap["resource_type"] = systemResourceItem.ResourceType
			}

			if systemResourceItem.Remark != nil {
				systemResourceItemMap["remark"] = systemResourceItem.Remark
			}

			if systemResourceItem.Region != nil {
				systemResourceItemMap["region"] = systemResourceItem.Region
			}

			if systemResourceItem.LatestResourceConfigVersion != nil {
				systemResourceItemMap["latest_resource_config_version"] = systemResourceItem.LatestResourceConfigVersion
			}

			ids = append(ids, *systemResourceItem.ResourceId)
			tmpList = append(tmpList, systemResourceItemMap)
		}

		_ = d.Set("resource_set", tmpList)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), d); e != nil {
			return e
		}
	}

	return nil
}
