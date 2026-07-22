package ci

import (
	"context"
	"fmt"
	"log"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/pkg/errors"
	"github.com/tencentyun/cos-go-sdk-v5"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudCiMediaSnapshotTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCiMediaSnapshotTemplateCreate,
		Read:   resourceTencentCloudCiMediaSnapshotTemplateRead,
		Update: resourceTencentCloudCiMediaSnapshotTemplateUpdate,
		Delete: resourceTencentCloudCiMediaSnapshotTemplateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"bucket": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "存储桶名称",
			},
			"name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "模板名称 仅 支持 `Chinese`，`English`，`numbers`，`_`，`-` 和 `*`。",
			},

			"snapshot": {
				Required:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "screenshot。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"mode": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Screenshot 模式，取值范围：{Interval，Average，KeyFrame}- Interval 表示 间隔 模式 Average 表示 average 模式- KeyFrame 表示 键 frame 模式- Interval 模式: Start，TimeInterval， Count 参数 takes effect. 当 Count 是 集合 和 TimeInterval 是 不 集合，表示to capture all frames， 总数 的 Count pictures- Average 模式: Start， Count 参数 takes effect. express。",
						},
						"start": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Starting 时间，[0 视频 时长] （秒）， Support float 格式， execution accuracy 是 accurate 到 milliseconds。",
						},
						"time_interval": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Screenshot 时间间隔，(0 3600]，（秒）， Support float 格式， execution accuracy 是 accurate 到 milliseconds。",
						},
						"count": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "数量 screenshots，范围 (0 10000]。",
						},
						"width": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "wide，取值范围：[128，4096]，单位：像素，如果 仅 宽度 是 集合，高度 是 calculated according 到 original ratio 的 视频。",
						},
						"height": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "high，取值范围：[128，4096]，单位：像素，如果 仅 高度 是 集合，宽度 是 calculated according 到 original ratio 的 视频。",
						},
						"ci_param": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Screenshot 镜像 processing 参数，对于 示例: imageMogr2/格式/png。",
						},
						"is_check_count": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "是否check 数量 screenshots forcibly，当 使用 自定义 间隔 模式 到 take screenshots， 视频 时间 是 不 long enough 到 capture Count screenshots，您 可以 switch 到 average screenshot 模式 到 capture Count screenshots。",
						},
						"is_check_black": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "是否enable black screen detection true/false。",
						},
						"black_level": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Screenshot black screen detection 参数，有效 当 IsCheckBlack=true，值 reference 范围 [30，100]，indicating proportion 的 black pixels， smaller 值， smaller proportion 的 black pixels，Start&gt;0， 参数 setting 是 无效，无 过滤器 black screen，Start =0 参数 是 有效， 开始时间 的 frame capture 是 first frame non-black screen start。",
						},
						"pixel_black_threshold": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Screenshot black screen detection 参数，有效 当 IsCheckBlack=true， 阈值 对于 judging whether pixel 是 black point，取值范围：[0，255]。",
						},
						"snapshot_out_mode": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Screenshot output 模式 参数，取值范围：{OnlySnapshot，OnlySprite，SnapshotAndSprite}，OnlySnapshot 表示 output 仅 screenshot 模式 OnlySprite 表示 仅 output sprite 模式 SnapshotAndSprite 表示 output screenshot 和 sprite 模式",
						},
						"sprite_snapshot_config": {
							Type:        schema.TypeList,
							Optional:    true,
							Computed:    true,
							MaxItems:    1,
							Description: "Screenshot output 配置。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"cell_width": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "Single 镜像 宽度 取值范围：[8，4096]，单位：像素。",
									},
									"cell_height": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "Single 镜像 高度 取值范围：[8，4096]，单位：像素。",
									},
									"padding": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "screenshot padding 大小，取值范围：[8，4096]，单位：像素。",
									},
									"margin": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "screenshot margin 大小，取值范围：[8，4096]，单位：像素。",
									},
									"color": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "See `https://www.ffmpeg.org/ffmpeg-utils.html#color-syntax` 对于 details 在 支持 colors。",
									},
									"columns": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "数量 screenshot columns，取值范围：[1，10000]。",
									},
									"lines": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "数量 screenshot lines，取值范围：[1，10000]。",
									},
								},
							},
						},
					},
				},
			},

			"template_id": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "模板 ID",
			},

			"update_time": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "更新时间。",
			},

			"create_time": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "创建时间。",
			},
		},
	}
}

func resourceTencentCloudCiMediaSnapshotTemplateCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ci_media_snapshot_template.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	var templateId string
	var bucket string
	if v, ok := d.GetOk("bucket"); ok {
		bucket = v.(string)
	} else {
		return errors.New("get bucket failed!")
	}
	request := cos.CreateMediaSnapshotTemplateOptions{
		Tag: "Snapshot",
	}

	if v, ok := d.GetOk("name"); ok {
		request.Name = v.(string)
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "snapshot"); ok {
		snapshot := cos.Snapshot{}
		if v, ok := dMap["mode"]; ok {
			snapshot.Mode = v.(string)
		}
		if v, ok := dMap["start"]; ok {
			snapshot.Start = v.(string)
		}
		if v, ok := dMap["time_interval"]; ok {
			snapshot.TimeInterval = v.(string)
		}
		if v, ok := dMap["count"]; ok {
			snapshot.Count = v.(string)
		}
		if v, ok := dMap["width"]; ok {
			snapshot.Width = v.(string)
		}
		if v, ok := dMap["height"]; ok {
			snapshot.Height = v.(string)
		}
		if v, ok := dMap["ci_param"]; ok {
			snapshot.CIParam = v.(string)
		}
		if v, ok := dMap["is_check_count"]; ok {
			if v.(string) == "true" {
				snapshot.IsCheckCount = true
			} else {
				snapshot.IsCheckCount = false
			}
		}
		if v, ok := dMap["is_check_black"]; ok {
			if v.(string) == "true" {
				snapshot.IsCheckBlack = true
			} else {
				snapshot.IsCheckBlack = false
			}
		}
		if v, ok := dMap["black_level"]; ok {
			snapshot.BlackLevel = v.(string)
		}
		if v, ok := dMap["pixel_black_threshold"]; ok {
			snapshot.PixelBlackThreshold = v.(string)
		}
		if v, ok := dMap["snapshot_out_mode"]; ok {
			snapshot.SnapshotOutMode = v.(string)
		}
		if spriteSnapshotConfigMap, ok := helper.InterfaceToMap(dMap, "sprite_snapshot_config"); ok {
			spriteSnapshotConfig := cos.SpriteSnapshotConfig{}
			if v, ok := spriteSnapshotConfigMap["cell_width"]; ok {
				spriteSnapshotConfig.CellWidth = v.(string)
			}
			if v, ok := spriteSnapshotConfigMap["cell_height"]; ok {
				spriteSnapshotConfig.CellHeight = v.(string)
			}
			if v, ok := spriteSnapshotConfigMap["padding"]; ok {
				spriteSnapshotConfig.Padding = v.(string)
			}
			if v, ok := spriteSnapshotConfigMap["margin"]; ok {
				spriteSnapshotConfig.Margin = v.(string)
			}
			if v, ok := spriteSnapshotConfigMap["color"]; ok {
				spriteSnapshotConfig.Color = v.(string)
			}
			if v, ok := spriteSnapshotConfigMap["columns"]; ok {
				spriteSnapshotConfig.Columns = v.(string)
			}
			if v, ok := spriteSnapshotConfigMap["lines"]; ok {
				spriteSnapshotConfig.Lines = v.(string)
			}
			snapshot.SpriteSnapshotConfig = &spriteSnapshotConfig
		}
		request.Snapshot = &snapshot
	}

	var response *cos.CreateMediaTemplateResult
	ciClient := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCiClient(bucket)
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, _, e := ciClient.CI.CreateMediaSnapshotTemplate(ctx, &request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%v], response body [%+v]\n", logId, "CreateMediaSnapshotTemplate", request, result)
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create ci mediaSnapshotTemplate failed, reason:%+v", logId, err)
		return err
	}

	templateId = response.Template.TemplateId
	d.SetId(bucket + tccommon.FILED_SP + templateId)

	return resourceTencentCloudCiMediaSnapshotTemplateRead(d, meta)
}

func resourceTencentCloudCiMediaSnapshotTemplateRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ci_media_snapshot_template.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := CiService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	bucket := idSplit[0]
	templateId := idSplit[1]

	mediaSnapshotTemplate, err := service.DescribeCiMediaTemplateById(ctx, bucket, templateId)
	if err != nil {
		return err
	}

	if mediaSnapshotTemplate == nil {
		d.SetId("")
		return fmt.Errorf("resource `track` %s does not exist", d.Id())
	}

	_ = d.Set("bucket", bucket)
	if mediaSnapshotTemplate.Name != "" {
		_ = d.Set("name", mediaSnapshotTemplate.Name)
	}

	log.Printf("[DEBUG]Snapshot api[%+v]", mediaSnapshotTemplate.Snapshot)
	if mediaSnapshotTemplate.Snapshot != nil {
		snapshotMap := map[string]interface{}{}

		if mediaSnapshotTemplate.Snapshot.Mode != "" {
			snapshotMap["mode"] = mediaSnapshotTemplate.Snapshot.Mode
		}

		if mediaSnapshotTemplate.Snapshot.Start != "" {
			snapshotMap["start"] = mediaSnapshotTemplate.Snapshot.Start
		}

		if mediaSnapshotTemplate.Snapshot.TimeInterval != "" {
			snapshotMap["time_interval"] = mediaSnapshotTemplate.Snapshot.TimeInterval
		}

		if mediaSnapshotTemplate.Snapshot.Count != "" {
			snapshotMap["count"] = mediaSnapshotTemplate.Snapshot.Count
		}

		if mediaSnapshotTemplate.Snapshot.Width != "" {
			snapshotMap["width"] = mediaSnapshotTemplate.Snapshot.Width
		}

		if mediaSnapshotTemplate.Snapshot.Height != "" {
			snapshotMap["height"] = mediaSnapshotTemplate.Snapshot.Height
		}

		if mediaSnapshotTemplate.Snapshot.CIParam != "" {
			snapshotMap["ci_param"] = mediaSnapshotTemplate.Snapshot.CIParam
		}

		snapshotMap["is_check_count"] = fmt.Sprintf("%t", mediaSnapshotTemplate.Snapshot.IsCheckCount)
		snapshotMap["is_check_black"] = fmt.Sprintf("%t", mediaSnapshotTemplate.Snapshot.IsCheckBlack)

		if mediaSnapshotTemplate.Snapshot.BlackLevel != "" {
			snapshotMap["black_level"] = mediaSnapshotTemplate.Snapshot.BlackLevel
		}

		if mediaSnapshotTemplate.Snapshot.PixelBlackThreshold != "" {
			snapshotMap["pixel_black_threshold"] = mediaSnapshotTemplate.Snapshot.PixelBlackThreshold
		}

		if mediaSnapshotTemplate.Snapshot.SnapshotOutMode != "" {
			snapshotMap["snapshot_out_mode"] = mediaSnapshotTemplate.Snapshot.SnapshotOutMode
		}

		if mediaSnapshotTemplate.Snapshot.SpriteSnapshotConfig != nil {
			spriteSnapshotConfigMap := map[string]interface{}{}

			if mediaSnapshotTemplate.Snapshot.SpriteSnapshotConfig.CellWidth != "" {
				spriteSnapshotConfigMap["cell_width"] = mediaSnapshotTemplate.Snapshot.SpriteSnapshotConfig.CellWidth
			}

			if mediaSnapshotTemplate.Snapshot.SpriteSnapshotConfig.CellHeight != "" {
				spriteSnapshotConfigMap["cell_height"] = mediaSnapshotTemplate.Snapshot.SpriteSnapshotConfig.CellHeight
			}

			if mediaSnapshotTemplate.Snapshot.SpriteSnapshotConfig.Padding != "" {
				spriteSnapshotConfigMap["padding"] = mediaSnapshotTemplate.Snapshot.SpriteSnapshotConfig.Padding
			}

			if mediaSnapshotTemplate.Snapshot.SpriteSnapshotConfig.Margin != "" {
				spriteSnapshotConfigMap["margin"] = mediaSnapshotTemplate.Snapshot.SpriteSnapshotConfig.Margin
			}

			if mediaSnapshotTemplate.Snapshot.SpriteSnapshotConfig.Color != "" {
				spriteSnapshotConfigMap["color"] = mediaSnapshotTemplate.Snapshot.SpriteSnapshotConfig.Color
			}

			if mediaSnapshotTemplate.Snapshot.SpriteSnapshotConfig.Columns != "" {
				spriteSnapshotConfigMap["columns"] = mediaSnapshotTemplate.Snapshot.SpriteSnapshotConfig.Columns
			}

			if mediaSnapshotTemplate.Snapshot.SpriteSnapshotConfig.Lines != "" {
				spriteSnapshotConfigMap["lines"] = mediaSnapshotTemplate.Snapshot.SpriteSnapshotConfig.Lines
			}

			snapshotMap["sprite_snapshot_config"] = []interface{}{spriteSnapshotConfigMap}
		}

		err = d.Set("snapshot", []interface{}{snapshotMap})
		if err != nil {
			return err
		}
	}

	if mediaSnapshotTemplate.TemplateId != "" {
		_ = d.Set("template_id", mediaSnapshotTemplate.TemplateId)
	}

	if mediaSnapshotTemplate.UpdateTime != "" {
		_ = d.Set("update_time", mediaSnapshotTemplate.UpdateTime)
	}

	if mediaSnapshotTemplate.CreateTime != "" {
		_ = d.Set("create_time", mediaSnapshotTemplate.CreateTime)
	}

	return nil
}

func resourceTencentCloudCiMediaSnapshotTemplateUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ci_media_snapshot_template.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	bucket := idSplit[0]
	templateId := idSplit[1]

	request := cos.CreateMediaSnapshotTemplateOptions{
		Tag: "Snapshot",
	}

	if v, ok := d.GetOk("name"); ok {
		request.Name = v.(string)
	}
	if d.HasChange("snapshot") {
		if dMap, ok := helper.InterfacesHeadMap(d, "snapshot"); ok {
			snapshot := cos.Snapshot{}
			if v, ok := dMap["mode"]; ok {
				snapshot.Mode = v.(string)
			}
			if v, ok := dMap["start"]; ok {
				snapshot.Start = v.(string)
			}
			if v, ok := dMap["time_interval"]; ok {
				snapshot.TimeInterval = v.(string)
			}
			if v, ok := dMap["count"]; ok {
				snapshot.Count = v.(string)
			}
			if v, ok := dMap["width"]; ok {
				snapshot.Width = v.(string)
			}
			if v, ok := dMap["height"]; ok {
				snapshot.Height = v.(string)
			}
			if v, ok := dMap["ci_param"]; ok {
				snapshot.CIParam = v.(string)
			}
			if v, ok := dMap["is_check_count"]; ok {
				snapshot.IsCheckCount = v.(bool)
			}
			if v, ok := dMap["is_check_black"]; ok {
				snapshot.IsCheckBlack = v.(bool)
			}
			if v, ok := dMap["black_level"]; ok {
				snapshot.BlackLevel = v.(string)
			}
			if v, ok := dMap["pixel_black_threshold"]; ok {
				snapshot.PixelBlackThreshold = v.(string)
			}
			if v, ok := dMap["snapshot_out_mode"]; ok {
				snapshot.SnapshotOutMode = v.(string)
			}
			if spriteSnapshotConfigMap, ok := helper.InterfacesHeadMap(d, "sprite_snapshot_config"); ok {
				spriteSnapshotConfig := cos.SpriteSnapshotConfig{}
				if v, ok := spriteSnapshotConfigMap["cell_width"]; ok {
					spriteSnapshotConfig.CellWidth = v.(string)
				}
				if v, ok := spriteSnapshotConfigMap["cell_height"]; ok {
					spriteSnapshotConfig.CellHeight = v.(string)
				}
				if v, ok := spriteSnapshotConfigMap["padding"]; ok {
					spriteSnapshotConfig.Padding = v.(string)
				}
				if v, ok := spriteSnapshotConfigMap["margin"]; ok {
					spriteSnapshotConfig.Margin = v.(string)
				}
				if v, ok := spriteSnapshotConfigMap["color"]; ok {
					spriteSnapshotConfig.Color = v.(string)
				}
				if v, ok := spriteSnapshotConfigMap["columns"]; ok {
					spriteSnapshotConfig.Columns = v.(string)
				}
				if v, ok := spriteSnapshotConfigMap["lines"]; ok {
					spriteSnapshotConfig.Lines = v.(string)
				}
				snapshot.SpriteSnapshotConfig = &spriteSnapshotConfig
			}
			request.Snapshot = &snapshot
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, _, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCiClient(bucket).CI.UpdateMediaSnapshotTemplate(ctx, &request, templateId)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%v], response body [%v]\n", logId, "UpdateMediaSnapshotTemplate", request, result)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create ci mediaSnapshotTemplate failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudCiMediaSnapshotTemplateRead(d, meta)
}

func resourceTencentCloudCiMediaSnapshotTemplateDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ci_media_snapshot_template.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := CiService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	bucket := idSplit[0]
	templateId := idSplit[1]

	if err := service.DeleteCiMediaTemplateById(ctx, bucket, templateId); err != nil {
		return err
	}

	return nil
}
