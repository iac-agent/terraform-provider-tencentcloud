package tcr

import (
	"context"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tcr "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tcr/v20190924"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudTCRInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTCRInstancesRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "名称 TCR 实例 到 查询。",
			},
			"instance_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID TCR 实例 到 查询。",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "用于保存结果。",
			},

			// Computed values
			"instance_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Information 列表 dedicated TCR 实例。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID TCR 实例。",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "名称 TCR 实例。",
						},

						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "状态 TCR 实例。",
						},
						"public_domain": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Public 地址 对于 访问 的 TCR 实例。",
						},
						"internal_end_point": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Internal 地址 对于 访问 的 TCR 实例。",
						},
						"instance_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "实例类型",
						},
						"tags": {
							Type:        schema.TypeMap,
							Computed:    true,
							Description: "标签 的 TCR 实例。",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudTCRInstancesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_tcr_instances.read")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	var name, instanceId string
	var filters = make([]*tcr.Filter, 0)
	if v, ok := d.GetOk("name"); ok {
		name = v.(string)
		filters = append(filters, &tcr.Filter{Name: helper.String("RegistryName"), Values: []*string{&name}})
	}

	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}

	tcrService := TCRService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	var outErr, inErr error
	instances, outErr := tcrService.DescribeTCRInstances(ctx, instanceId, filters)
	if outErr != nil {
		outErr = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			instances, inErr = tcrService.DescribeTCRInstances(ctx, instanceId, filters)
			if inErr != nil {
				return tccommon.RetryError(inErr)
			}
			return nil
		})
	}

	if outErr != nil {
		return outErr
	}

	ids := make([]string, 0, len(instances))
	instanceList := make([]map[string]interface{}, 0, len(instances))
	for _, ins := range instances {
		mapping := map[string]interface{}{
			"id":                 ins.RegistryId,
			"name":               ins.RegistryName,
			"status":             ins.Status,
			"public_domain":      ins.PublicDomain,
			"instance_type":      ins.RegistryType,
			"internal_end_point": ins.InternalEndpoint,
		}
		tags := make(map[string]string, len(ins.TagSpecification.Tags))
		for _, tag := range ins.TagSpecification.Tags {
			tags[*tag.Key] = *tag.Value
		}
		mapping["tags"] = tags
		instanceList = append(instanceList, mapping)
		ids = append(ids, *ins.RegistryId)
	}
	d.SetId(helper.DataResourceIdsHash(ids))
	if e := d.Set("instance_list", instanceList); e != nil {
		log.Printf("[CRITAL]%s provider set TCR instance list fail, reason:%s\n", logId, e)
		return e
	}

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), instanceList); e != nil {
			return e
		}
	}

	return nil

}
