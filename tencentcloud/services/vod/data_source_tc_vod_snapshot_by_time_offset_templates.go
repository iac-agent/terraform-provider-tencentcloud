package vod

import (
	"context"
	"log"
	"strconv"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudVodSnapshotByTimeOffsetTemplates() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudVodSnapshotByTimeOffsetTemplatesRead,

		Schema: map[string]*schema.Schema{
			"definition": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Unique ID 过滤器 的 快照 通过 时间 偏移量 template。",
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "模板 类型 过滤器. 有效值：`Preset`，`Custom`. `Preset`: preset template; `Custom`: 自定义 template。",
			},
			"sub_app_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Subapplication ID 在 VOD. 如果 您 need 到 访问 资源 在 subapplication，enter subapplication ID 在 此 字段; otherwise，leave 它 空。",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "用于保存结果。",
			},
			"template_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A 列表 快照 通过 时间 偏移量 templates. Each element 包含following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"definition": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Unique ID 快照 通过 时间 偏移量 template。",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "模板 类型 过滤器. 有效值：`Preset`，`Custom`. `Preset`: preset template; `Custom`: 自定义 template。",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "名称 时间 point screen capturing template。",
						},
						"width": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Maximum 值 的 `宽度` (或 long side) 的 screenshot （像素）。 取值范围：0 和 [128，4,096]. 如果 both `宽度` 和 `高度` 是 `0`， resolution 将 是 same 作为 该 的 来源 视频; 如果 `宽度` 是 `0`，但 `高度` 是 不 `0`，宽度 将 是 proportionally scaled; 如果 `宽度` 是 不 `0`，但 `高度` 是 `0`，`高度` 将 是 proportionally scaled; 如果 both `宽度` 和 `高度` 是 不 `0`， 自定义 resolution 将 是 使用。",
						},
						"height": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Maximum 值 的 `高度` (或 short side) 的 screenshot （像素）。 取值范围：0 和 [128，4,096]. 如果 both `宽度` 和 `高度` 是 `0`， resolution 将 是 same 作为 该 的 来源 视频; 如果 `宽度` 是 `0`，但 `高度` 是 不 `0`，`宽度` 将 是 proportionally scaled; 如果 `宽度` 是 不 `0`，但 `高度` 是 `0`，`高度` 将 是 proportionally scaled; 如果 both `宽度` 和 `高度` 是 不 `0`， 自定义 resolution 将 是 使用。",
						},
						"resolution_adaptive": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Resolution adaption. 有效值：`true`，`false`. `true`: 已启用 In 此 case，`宽度` 表示 long side 的 视频，while `高度` short side; `false`: 已禁用 In 此 case，`宽度` 表示 宽度 的 视频，while `高度` 高度。",
						},
						"format": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Image 格式 有效值：`jpg`，`png`。",
						},
						"comment": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "模板描述",
						},
						"fill_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Fill refers 到 way 的 processing screenshot 当 its aspect ratio 是 different 从 该 的 来源 视频. following fill types 是 支持: `stretch`: stretch. screenshot 将 是 stretched frame 通过 frame 到 match aspect ratio 的 来源 视频，其中 可能 make screenshot `shorter` 或 `longer`; `black`: fill 使用 black. 此 选项 retains aspect ratio 的 来源 视频 对于 screenshot 和 fills unmatched area 使用 black color blocks. `white`: fill 使用 white. 此 选项 retains aspect ratio 的 来源 视频 对于 screenshot 和 fills unmatched area 使用 white color blocks. `gauss`: fill 使用 Gaussian blur. 此 选项 retains aspect ratio 的 来源 视频 对于 screenshot 和 fills unmatched area 使用 Gaussian blur。",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "创建时间 的 template 在 ISO date 格式",
						},
						"update_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "最后修改时间 的 template 在 ISO date 格式",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudVodSnapshotByTimeOffsetTemplatesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_vod_snapshot_by_time_offset_templates.read")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	filter := make(map[string]interface{})
	if v, ok := d.GetOk("definition"); ok {
		filter["definitions"] = []string{v.(string)}
	}
	if v, ok := d.GetOk("type"); ok {
		filter["type"] = v.(string)
	}
	if v, ok := d.GetOk("sub_app_id"); ok {
		filter["sub_appid"] = v.(int)
	}

	vodService := VodService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}
	templates, err := vodService.DescribeSnapshotByTimeOffsetTemplatesByFilter(ctx, filter)
	if err != nil {
		return err
	}

	templatesList := make([]map[string]interface{}, 0, len(templates))
	ids := make([]string, 0, len(templates))
	for _, item := range templates {
		definitionStr := strconv.FormatUint(*item.Definition, 10)
		templatesList = append(templatesList, map[string]interface{}{
			"definition":          definitionStr,
			"type":                item.Type,
			"name":                item.Name,
			"width":               item.Width,
			"height":              item.Height,
			"resolution_adaptive": *item.ResolutionAdaptive == "open",
			"format":              item.Format,
			"comment":             item.Comment,
			"fill_type":           item.FillType,
			"create_time":         item.CreateTime,
			"update_time":         item.UpdateTime,
		})
		ids = append(ids, definitionStr)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	if e := d.Set("template_list", templatesList); e != nil {
		log.Printf("[CRITAL]%s provider set vod snapshot by time offset template list fail, reason:%s ", logId, e.Error())
	}

	if output, ok := d.GetOk("result_output_file"); ok && output.(string) != "" {
		if err := tccommon.WriteToFile(output.(string), templatesList); err != nil {
			log.Printf("[CRITAL]%s output file[%s] fail, reason[%s]", logId, output.(string), err.Error())
			return err
		}
	}

	return nil
}
