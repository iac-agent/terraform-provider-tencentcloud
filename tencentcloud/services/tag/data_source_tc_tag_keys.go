package tag

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudTagKeys() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTagKeysRead,
		Schema: map[string]*schema.Schema{
			"create_uin": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "创建者 `Uin`. 如果未指定，`Uin` 是 仅 使用 作为 查询 condition。",
			},

			"show_project": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "是否show 项目. Allow 值: 0: 无，1: yes。",
			},

			"category": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "标签 类型 有效值：Custom: 自定义 标签; System: 系统 标签; All: all 标签 默认值：All。",
			},

			"tags": {
				Type:        schema.TypeSet,
				Computed:    true,
				Description: "标签列表",
				Elem: &schema.Schema{
					Type: schema.TypeString,
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

func dataSourceTencentCloudTagKeysRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_tag_keys.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(nil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = TagService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOkExists("create_uin"); ok {
		paramMap["CreateUin"] = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("show_project"); ok {
		paramMap["ShowProject"] = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("category"); ok {
		paramMap["Category"] = helper.String(v.(string))
	}

	var respData []*string
	reqErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeTagKeysByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		respData = result
		return nil
	})

	if reqErr != nil {
		return reqErr
	}

	if respData != nil {
		_ = d.Set("tags", respData)
	}

	d.SetId(helper.BuildToken())
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), d); e != nil {
			return e
		}
	}

	return nil
}
