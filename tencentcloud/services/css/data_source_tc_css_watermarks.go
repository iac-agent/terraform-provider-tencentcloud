package css

import (
	"context"
	"strconv"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	css "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/live/v20180801"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudCssWatermarks() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCssWatermarksRead,
		Schema: map[string]*schema.Schema{
			"watermark_list": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Watermark 信息 列表。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"watermark_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Watermark ID。",
						},
						"picture_url": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Watermark 镜像 URL",
						},
						"x_position": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Display position: X-axis 偏移量",
						},
						"y_position": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Display position: Y-axis 偏移量",
						},
						"watermark_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Watermark 名称",
						},
						"status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Current 状态 0: 不 使用. 1: 在 使用。",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "时间 当 水印 是 added.注意: Beijing 时间 (UTC+8) 是 使用。",
						},
						"width": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Watermark 宽度。",
						},
						"height": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Watermark 高度。",
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

func dataSourceTencentCloudCssWatermarksRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_css_watermarks.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	service := CssService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var watermarks []*css.WatermarkInfo
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeCssWatermarksByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		watermarks = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(watermarks))
	tmpList := make([]map[string]interface{}, 0)
	if watermarks != nil {
		for _, watermarkInfo := range watermarks {
			watermarkInfoMap := map[string]interface{}{}

			if watermarkInfo.WatermarkId != nil {
				watermarkInfoMap["watermark_id"] = watermarkInfo.WatermarkId
			}

			if watermarkInfo.PictureUrl != nil {
				watermarkInfoMap["picture_url"] = watermarkInfo.PictureUrl
			}

			if watermarkInfo.XPosition != nil {
				watermarkInfoMap["x_position"] = watermarkInfo.XPosition
			}

			if watermarkInfo.YPosition != nil {
				watermarkInfoMap["y_position"] = watermarkInfo.YPosition
			}

			if watermarkInfo.WatermarkName != nil {
				watermarkInfoMap["watermark_name"] = watermarkInfo.WatermarkName
			}

			if watermarkInfo.Status != nil {
				watermarkInfoMap["status"] = watermarkInfo.Status
			}

			if watermarkInfo.CreateTime != nil {
				watermarkInfoMap["create_time"] = watermarkInfo.CreateTime
			}

			if watermarkInfo.Width != nil {
				watermarkInfoMap["width"] = watermarkInfo.Width
			}

			if watermarkInfo.Height != nil {
				watermarkInfoMap["height"] = watermarkInfo.Height
			}

			ids = append(ids, strconv.FormatInt(*watermarkInfo.WatermarkId, 10))
			tmpList = append(tmpList, watermarkInfoMap)
		}

		_ = d.Set("watermark_list", tmpList)
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
