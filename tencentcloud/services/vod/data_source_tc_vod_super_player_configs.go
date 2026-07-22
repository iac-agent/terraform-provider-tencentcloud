package vod

import (
	"context"
	"log"
	"strconv"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudVodSuperPlayerConfigs() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudVodSuperPlayerConfigsRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "名称 super player 配置",
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "配置 类型 filter. 有效值：`Preset`，`Custom`. `Preset`: preset template; `Custom`: custom template。",
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
			"config_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A 列表 super player configs. Each element 包含following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Template 类型 filter. 有效值：`Preset`，`Custom`. `Preset`: preset template; `Custom`: custom template。",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Player configuration 名称，which can contain up to 64 letters，digits，underscores，and hyphens (such as test_ABC-123) and must be unique under a 用户",
						},
						"drm_switch": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Switch of DRM-protected adaptive bitstream playback: `true`: 已启用，indicating to play back only output adaptive bitstreams protected by DRM; `false`: 已禁用，indicating to play back unencrypted output adaptive bitstreams。",
						},
						"adaptive_dynamic_streaming_definition": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID unencrypted adaptive bitrate streaming template that allows output，which 为必填项 if `drm_switch` is `false`。",
						},
						"drm_streaming_info": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "内容 of the DRM-protected adaptive bitrate streaming template that allows output，which 为必填项 if `drm_switch` is `true`。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"simple_aes_definition": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "ID adaptive dynamic streaming template whose protection 类型 is `SimpleAES`。",
									},
								},
							},
						},
						"image_sprite_definition": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID image sprite template that allows output。",
						},
						"resolution_names": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Display 名称 player for substreams with different resolutions. 如果此参数为空 or an empty array，the default configuration will be used: `min_edge_length: 240，名称: LD`; `min_edge_length: 480，名称: SD`; `min_edge_length: 720，名称: HD`; `min_edge_length: 1080，名称: FHD`; `min_edge_length: 1440，名称: 2K`; `min_edge_length: 2160，名称: 4K`; `min_edge_length: 4320，名称: 8K`。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"min_edge_length": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Length of video short side （像素）。",
									},
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Display 名称",
									},
								},
							},
						},
						"domain": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "域名 名称 用于playback. If it is left empty or set to `Default`，the 域名 名称 configured in [Default Distribution Configuration](https://cloud.tencent.com/document/product/266/33373) will be used。",
						},
						"scheme": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Scheme 用于playback. If it is left empty or set to `Default`，the scheme configured in [Default Distribution Configuration](https://cloud.tencent.com/document/product/266/33373) will be used. Other 有效值：`HTTP`; `HTTPS`。",
						},
						"comment": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "模板描述",
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

func dataSourceTencentCloudVodSuperPlayerConfigsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_vod_super_player_configs.read")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	filter := make(map[string]interface{})
	if v, ok := d.GetOk("name"); ok {
		filter["names"] = []string{v.(string)}
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
	configs, err := vodService.DescribeSuperPlayerConfigsByFilter(ctx, filter)
	if err != nil {
		return err
	}

	configsList := make([]map[string]interface{}, 0, len(configs))
	ids := make([]string, 0, len(configs))
	for _, item := range configs {
		configsList = append(configsList, func() map[string]interface{} {
			mapping := map[string]interface{}{
				"type":        item.Type,
				"name":        item.Name,
				"drm_switch":  *item.DrmSwitch == "ON",
				"domain":      item.Domain,
				"scheme":      item.Scheme,
				"comment":     item.Comment,
				"create_time": item.CreateTime,
				"update_time": item.UpdateTime,
			}
			// workaround for AdaptiveDynamicStreamingDefinition para cuz it's dirty data.
			if *item.DrmSwitch == "OFF" {
				mapping["adaptive_dynamic_streaming_definition"] = strconv.FormatUint(*item.AdaptiveDynamicStreamingDefinition, 10)
			}
			if item.DrmStreamingsInfo != nil && item.DrmStreamingsInfo.SimpleAesDefinition != nil {
				mapping["drm_streaming_info"] = []map[string]interface{}{
					{
						"simple_aes_definition": strconv.FormatUint(*item.DrmStreamingsInfo.SimpleAesDefinition, 10),
					},
				}
			}
			if item.ImageSpriteDefinition != nil {
				mapping["image_sprite_definition"] = strconv.FormatUint(*item.ImageSpriteDefinition, 10)
			}
			mapping["resolution_names"] = func() []map[string]interface{} {
				namesMap := make([]map[string]interface{}, 0, len(item.ResolutionNameSet))
				for _, v := range item.ResolutionNameSet {
					namesMap = append(namesMap, map[string]interface{}{
						"min_edge_length": v.MinEdgeLength,
						"name":            v.Name,
					})
				}
				return namesMap
			}()
			ids = append(ids, *item.Name)
			return mapping
		}())
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	if e := d.Set("config_list", configsList); e != nil {
		log.Printf("[CRITAL]%s provider set vod super player config list fail, reason:%s ", logId, e.Error())
	}

	if output, ok := d.GetOk("result_output_file"); ok && output.(string) != "" {
		if err := tccommon.WriteToFile(output.(string), configsList); err != nil {
			log.Printf("[CRITAL]%s output file[%s] fail, reason[%s]", logId, output.(string), err.Error())
			return err
		}
	}

	return nil
}
