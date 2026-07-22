package vpc

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudVpcClassicLinkInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudVpcClassicLinkInstancesRead,
		Schema: map[string]*schema.Schema{
			"filters": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "过滤器 conditions.`vpc-ID` - String - (过滤器 condition) VPC 实例 ID `vm-ip` - String - (过滤器 condition) IP 地址 的 CVM 在 basic 网络。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "attribute 名称 如果 more 比 一个 过滤器 exists， logical relation between these Filters 是 `AND`。",
						},
						"values": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Required:    true,
							Description: "attribute 值 如果 there 是 多个 Values 对于 一个 过滤器， logical relation between these Values under same 过滤器 是 `OR`。",
						},
					},
				},
			},

			"classic_link_instance_set": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Classiclink 实例。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vpc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "VPC 实例 ID",
						},
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "唯一 ID CVM 实例。",
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

func dataSourceTencentCloudVpcClassicLinkInstancesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_vpc_classic_link_instances.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("filters"); ok {
		filtersSet := v.([]interface{})
		tmpSet := make([]*vpc.FilterObject, 0, len(filtersSet))

		for _, item := range filtersSet {
			filterObject := vpc.FilterObject{}
			filterObjectMap := item.(map[string]interface{})

			if v, ok := filterObjectMap["name"]; ok {
				filterObject.Name = helper.String(v.(string))
			}
			if v, ok := filterObjectMap["values"]; ok {
				valuesSet := v.(*schema.Set).List()
				filterObject.Values = helper.InterfacesStringsPoint(valuesSet)
			}
			tmpSet = append(tmpSet, &filterObject)
		}
		paramMap["Filters"] = tmpSet
	}

	service := VpcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var classicLinkInstanceSet []*vpc.ClassicLinkInstance

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeVpcClassicLinkInstancesByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		classicLinkInstanceSet = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(classicLinkInstanceSet))
	tmpList := make([]map[string]interface{}, 0, len(classicLinkInstanceSet))

	if classicLinkInstanceSet != nil {
		for _, classicLinkInstance := range classicLinkInstanceSet {
			classicLinkInstanceMap := map[string]interface{}{}

			if classicLinkInstance.VpcId != nil {
				classicLinkInstanceMap["vpc_id"] = classicLinkInstance.VpcId
			}

			if classicLinkInstance.InstanceId != nil {
				classicLinkInstanceMap["instance_id"] = classicLinkInstance.InstanceId
			}

			ids = append(ids, *classicLinkInstance.InstanceId)
			tmpList = append(tmpList, classicLinkInstanceMap)
		}

		_ = d.Set("classic_link_instance_set", tmpList)
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
