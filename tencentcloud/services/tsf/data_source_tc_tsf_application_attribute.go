package tsf

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tsf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tsf/v20180326"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudTsfApplicationAttribute() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTsfApplicationAttributeRead,
		Schema: map[string]*schema.Schema{
			"application_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "应用 ID。",
			},

			"result": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "应用 列表 other attribute。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Total 数量 实例.注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"run_instance_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "数量 running 实例.注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"group_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "数量 部署 groups under 应用.注意：此字段可能返回 null，表示无法获取有效值。",
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

func dataSourceTencentCloudTsfApplicationAttributeRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_tsf_application_attribute.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	ids := ""

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("application_id"); ok {
		ids = v.(string)
		paramMap["ApplicationId"] = helper.String(v.(string))
	}

	service := TsfService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var attribute *tsf.ApplicationAttribute

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeTsfApplicationAttributeByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		attribute = result
		return nil
	})
	if err != nil {
		return err
	}

	applicationAttributeMap := map[string]interface{}{}
	if attribute != nil {
		if attribute.InstanceCount != nil {
			applicationAttributeMap["instance_count"] = attribute.InstanceCount
		}

		if attribute.RunInstanceCount != nil {
			applicationAttributeMap["run_instance_count"] = attribute.RunInstanceCount
		}

		if attribute.GroupCount != nil {
			applicationAttributeMap["group_count"] = attribute.GroupCount
		}

		_ = d.Set("result", []interface{}{applicationAttributeMap})
	}

	d.SetId(ids)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), applicationAttributeMap); e != nil {
			return e
		}
	}
	return nil
}
