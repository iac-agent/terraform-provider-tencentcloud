package vod

import (
	"context"
	"log"
	"strconv"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudVodImageSpriteTemplates() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudVodImageSpriteTemplatesRead,

		Schema: map[string]*schema.Schema{
			"definition": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Unique ID filter of image sprite template。",
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Template 类型 filter. 有效值：`Preset`，`Custom`. `Preset`: preset template; `Custom`: custom template。",
			},
			"sub_app_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Subapplication ID in VOD. If you need to access a resource in a subapplication，enter the subapplication ID in this field; otherwise，leave it empty。",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "用于保存结果。",
			},
			"template_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A 列表 image sprite templates. Each element 包含following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"definition": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Unique ID image sprite template。",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Template 类型 filter. 有效值：`Preset`，`Custom`. `Preset`: preset template; `Custom`: custom template。",
						},
						"sample_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Sampling 类型 有效值：`Percent`，`Time`. `Percent`: by percent. `Time`: by 时间间隔。",
						},
						"sample_interval": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Sampling interval. If `sample_type` is `Percent`，sampling will be performed at an interval of the specified percentage. If `sample_type` is `Time`，sampling will be performed at the specified 时间间隔 （秒）。",
						},
						"row_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Subimage row count of an image sprite。",
						},
						"column_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Subimage column count of an image sprite。",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "名称 a time point screen capturing template。",
						},
						"comment": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "模板描述",
						},
						"fill_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Fill refers to the way of processing a screenshot when its aspect ratio is different from that of the 来源 video. The following fill types are supported: `stretch`: stretch. The screenshot will be stretched frame by frame to match the aspect ratio of the 来源 video，which may make the screenshot shorter or longer; `black`: fill with black. This option retains the aspect ratio of the 来源 video for the screenshot and fills the unmatched area with black color blocks。",
						},
						"width": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Maximum 值 of the `width` (or long side) of a screenshot （像素）。 取值范围：0 and [128，4,096]. If both `width` and `height` are `0`，the resolution will be the same as that of the 来源 video; If `width` is `0`，but `height` is not `0`，width will be proportionally scaled; If `width` is not `0`，but `height` is `0`，`height` will be proportionally scaled; If both `width` and `height` are not `0`，the custom resolution will be used。",
						},
						"height": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Maximum 值 of the `height` (or short side) of a screenshot （像素）。 取值范围：0 and [128，4,096]. If both `width` and `height` are `0`，the resolution will be the same as that of the 来源 video; If `width` is `0`，but `height` is not `0`，`width` will be proportionally scaled; If `width` is not `0`，but `height` is `0`，`height` will be proportionally scaled; If both `width` and `height` are not `0`，the custom resolution will be used。",
						},
						"resolution_adaptive": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Resolution adaption. 有效值：`true`: 已启用 In this case，`width` represents the long side of a video，while `height` the short side; `false`: 已禁用 In this case，`width` represents the width of a video，while `height` the height。",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "创建时间 of template in ISO date 格式",
						},
						"update_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "最后修改时间 of template in ISO date 格式",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudVodImageSpriteTemplatesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_vod_image_sprite_templates.read")()

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
	templates, err := vodService.DescribeImageSpriteTemplatesByFilter(ctx, filter)
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
			"sample_type":         item.SampleType,
			"sample_interval":     item.SampleInterval,
			"row_count":           item.RowCount,
			"column_count":        item.ColumnCount,
			"name":                item.Name,
			"comment":             item.Comment,
			"fill_type":           item.FillType,
			"width":               item.Width,
			"height":              item.Height,
			"resolution_adaptive": *item.ResolutionAdaptive == "open",
			"create_time":         item.CreateTime,
			"update_time":         item.UpdateTime,
		})
		ids = append(ids, definitionStr)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	if e := d.Set("template_list", templatesList); e != nil {
		log.Printf("[CRITAL]%s provider set vod image sprite template list fail, reason:%s ", logId, e.Error())
	}

	if output, ok := d.GetOk("result_output_file"); ok && output.(string) != "" {
		if err := tccommon.WriteToFile(output.(string), templatesList); err != nil {
			log.Printf("[CRITAL]%s output file[%s] fail, reason[%s]", logId, output.(string), err.Error())
			return err
		}
	}

	return nil
}
