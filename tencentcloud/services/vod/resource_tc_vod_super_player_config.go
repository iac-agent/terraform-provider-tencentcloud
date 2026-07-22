package vod

import (
	"context"
	"log"
	"strconv"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vod "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vod/v20180717"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

func ResourceTencentCloudVodSuperPlayerConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudVodSuperPlayerConfigCreate,
		Read:   resourceTencentCloudVodSuperPlayerConfigRead,
		Update: resourceTencentCloudVodSuperPlayerConfigUpdate,
		Delete: resourceTencentCloudVodSuperPlayerConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: tccommon.ValidateStringLengthInRange(1, 64),
				Description:  "Player 配置 名称，其中 可以 contain up 到 64 letters，digits，underscores，和 hyphens (such 作为 test_ABC-123) 和 必须 是 唯一 under 用户",
			},
			"drm_switch": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Switch 的 DRM-protected adaptive bitstream playback: `true`: 已启用，indicating 到 play back 仅 output adaptive bitstreams protected 通过 DRM; `false`: 已禁用，indicating 到 play back unencrypted output adaptive bitstreams. 默认值：`false`。",
			},
			"adaptive_dynamic_streaming_definition": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID unencrypted adaptive bitrate streaming template 该 allows output，其中 为必填项 如果 `drm_switch` 是 `false`。",
			},
			"drm_streaming_info": {
				Type:        schema.TypeList,
				MinItems:    1,
				MaxItems:    1,
				Optional:    true,
				Description: "内容 的 DRM-protected adaptive bitrate streaming template 该 allows output，其中 为必填项 如果 `drm_switch` 是 `true`。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"simple_aes_definition": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "ID adaptive 动态 streaming template whose protection 类型 是 `SimpleAES`。",
						},
					},
				},
			},
			"image_sprite_definition": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID 镜像 sprite template 该 allows output。",
			},
			"resolution_names": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Display 名称 player 对于 substreams 使用 different resolutions. 如果此参数为空 或 空 数组， 默认值 配置 将 是 使用: `min_edge_length: 240，名称: LD`; `min_edge_length: 480，名称: SD`; `min_edge_length: 720，名称: HD`; `min_edge_length: 1080，名称: FHD`; `min_edge_length: 1440，名称: 2K`; `min_edge_length: 2160，名称: 4K`; `min_edge_length: 4320，名称: 8K`。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"min_edge_length": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Length 的 视频 short side （像素）。",
						},
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Display 名称",
						},
					},
				},
			},
			"domain": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "Default",
				Description: "域名 名称 用于playback. 如果 它 是 left 空 或 集合 到 `Default`， 域名 名称 已配置 在 [Default Distribution Configuration](https://云.tencent.com/document/product/266/33373) 将 是 使用. `Default` 通过 默认值。",
			},
			"scheme": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "Default",
				Description: "Scheme 用于playback. 如果 它 是 left 空 或 集合 到 `Default`， scheme 已配置 在 [Default Distribution Configuration](https://云.tencent.com/document/product/266/33373) 将 是 使用. Other 有效值：`HTTP`; `HTTPS`。",
			},
			"comment": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: tccommon.ValidateStringLengthInRange(1, 256),
				Description:  "模板描述 Length 限制: 256 字符。",
			},
			"sub_app_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Subapplication ID 在 VOD. 如果 您 need 到 访问 资源 在 subapplication，enter subapplication ID 在 此 字段; otherwise，leave 它 空。",
			},
			// computed
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
	}
}

func resourceTencentCloudVodSuperPlayerConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vod_super_player_config.create")()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		request = vod.NewCreateSuperPlayerConfigRequest()
	)

	request.Name = helper.String(d.Get("name").(string))
	request.DrmSwitch = helper.String(DRM_SWITCH_TO_STRING[d.Get("drm_switch").(bool)])
	if v, ok := d.GetOk("adaptive_dynamic_streaming_definition"); ok {
		idUint, _ := strconv.ParseUint(v.(string), 0, 64)
		request.AdaptiveDynamicStreamingDefinition = &idUint
	}
	if v, ok := d.GetOk("drm_streaming_info"); ok {
		request.DrmStreamingsInfo = &vod.DrmStreamingsInfo{
			SimpleAesDefinition: func(value interface{}) *uint64 {
				vv := value.([]interface{})
				vvv := vv[0].(map[string]interface{})
				idUint, _ := strconv.ParseUint(vvv["simple_aes_definition"].(string), 0, 64)
				return &idUint
			}(v),
		}
	}
	if v, ok := d.GetOk("image_sprite_definition"); ok {
		idUint, _ := strconv.ParseUint(v.(string), 0, 64)
		request.ImageSpriteDefinition = &idUint
	}
	resolutionNames := []*vod.ResolutionNameInfo{
		{
			MinEdgeLength: helper.Uint64(240),
			Name:          helper.String("LD"),
		},
		{
			MinEdgeLength: helper.Uint64(480),
			Name:          helper.String("SD"),
		},
		{
			MinEdgeLength: helper.Uint64(720),
			Name:          helper.String("HD"),
		},
		{
			MinEdgeLength: helper.Uint64(1080),
			Name:          helper.String("FHD"),
		},
		{
			MinEdgeLength: helper.Uint64(1440),
			Name:          helper.String("2K"),
		},
		{
			MinEdgeLength: helper.Uint64(2160),
			Name:          helper.String("4K"),
		},
		{
			MinEdgeLength: helper.Uint64(4320),
			Name:          helper.String("8K"),
		},
	}
	if v, ok := d.GetOk("resolution_names"); ok {
		resolutionNames = make([]*vod.ResolutionNameInfo, 0, len(v.([]interface{})))
		for _, item := range v.([]interface{}) {
			itemV := item.(map[string]interface{})
			resolutionNames = append(resolutionNames, &vod.ResolutionNameInfo{
				MinEdgeLength: helper.IntUint64(itemV["min_edge_length"].(int)),
				Name:          helper.String(itemV["name"].(string)),
			})
		}
	}
	request.ResolutionNames = resolutionNames
	request.Domain = helper.String(d.Get("domain").(string))
	request.Scheme = helper.String(d.Get("scheme").(string))
	if v, ok := d.GetOk("comment"); ok {
		request.Comment = helper.String(v.(string))
	}
	if v, ok := d.GetOk("sub_app_id"); ok {
		request.SubAppId = helper.IntUint64(v.(int))
	}

	var err error
	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		_, err = meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVodClient().CreateSuperPlayerConfig(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, reason:%s", logId, request.GetAction(), err.Error())
			return tccommon.RetryError(err)
		}
		return nil
	})
	if err != nil {
		return err
	}

	d.SetId(d.Get("name").(string))

	return resourceTencentCloudVodSuperPlayerConfigRead(d, meta)
}

func resourceTencentCloudVodSuperPlayerConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vod_super_player_config.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		id         = d.Id()
		client     = meta.(tccommon.ProviderMeta).GetAPIV3Conn()
		vodService = VodService{client: client}
	)
	config, has, err := vodService.DescribeSuperPlayerConfigsById(ctx, id)
	if err != nil {
		return err
	}
	if !has {
		d.SetId("")
		return nil
	}

	_ = d.Set("name", config.Name)
	_ = d.Set("drm_switch", *config.DrmSwitch == "ON")
	// workaround for AdaptiveDynamicStreamingDefinition para cuz it's dirty data.
	if *config.DrmSwitch == "OFF" {
		_ = d.Set("adaptive_dynamic_streaming_definition", strconv.FormatUint(*config.AdaptiveDynamicStreamingDefinition, 10))
	}
	if config.DrmStreamingsInfo != nil && config.DrmStreamingsInfo.SimpleAesDefinition != nil {
		_ = d.Set("drm_streaming_info", []map[string]interface{}{
			{
				"simple_aes_definition": strconv.FormatUint(*config.DrmStreamingsInfo.SimpleAesDefinition, 10),
			},
		})
	}
	if config.ImageSpriteDefinition != nil {
		_ = d.Set("image_sprite_definition", strconv.FormatUint(*config.ImageSpriteDefinition, 10))
	}
	_ = d.Set("resolution_names", func() []map[string]interface{} {
		namesMap := make([]map[string]interface{}, 0, len(config.ResolutionNameSet))
		for _, v := range config.ResolutionNameSet {
			namesMap = append(namesMap, map[string]interface{}{
				"min_edge_length": v.MinEdgeLength,
				"name":            v.Name,
			})
		}
		return namesMap
	}())
	_ = d.Set("domain", config.Domain)
	_ = d.Set("scheme", config.Scheme)
	_ = d.Set("comment", config.Comment)
	_ = d.Set("create_time", config.CreateTime)
	_ = d.Set("update_time", config.UpdateTime)

	return nil
}

func resourceTencentCloudVodSuperPlayerConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vod_super_player_config.update")()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		request    = vod.NewModifySuperPlayerConfigRequest()
		id         = d.Id()
		changeFlag = false
	)

	request.Name = &id
	if d.HasChange("drm_switch") {
		changeFlag = true
		request.DrmSwitch = helper.String(DRM_SWITCH_TO_STRING[d.Get("drm_switch").(bool)])
	}
	if d.HasChange("adaptive_dynamic_streaming_definition") {
		changeFlag = true
		idUint, _ := strconv.ParseUint(d.Get("adaptive_dynamic_streaming_definition").(string), 0, 64)
		request.AdaptiveDynamicStreamingDefinition = &idUint
	}
	if d.HasChange("drm_streaming_info") {
		changeFlag = true
		request.DrmStreamingsInfo = func() *vod.DrmStreamingsInfoForUpdate {
			if v, ok := d.GetOk("drm_streaming_info"); !ok {
				return nil
			} else {
				return &vod.DrmStreamingsInfoForUpdate{
					SimpleAesDefinition: func(value interface{}) *uint64 {
						vv := value.([]interface{})
						vvv := vv[0].(map[string]interface{})
						idUint, _ := strconv.ParseUint(vvv["simple_aes_definition"].(string), 0, 64)
						return &idUint
					}(v),
				}
			}
		}()
	}
	if d.HasChange("image_sprite_definition") {
		changeFlag = true
		idUint, _ := strconv.ParseUint(d.Get("image_sprite_definition").(string), 0, 64)
		request.ImageSpriteDefinition = &idUint
	}
	if d.HasChange("resolution_names") {
		changeFlag = true
		v := d.Get("resolution_names")
		resolutionNames := make([]*vod.ResolutionNameInfo, 0, len(v.([]interface{})))
		for _, item := range v.([]interface{}) {
			itemV := item.(map[string]interface{})
			resolutionNames = append(resolutionNames, &vod.ResolutionNameInfo{
				MinEdgeLength: helper.IntUint64(itemV["min_edge_length"].(int)),
				Name:          helper.String(itemV["name"].(string)),
			})
		}
		request.ResolutionNames = resolutionNames
	}
	if d.HasChange("domain") {
		changeFlag = true
		request.Domain = helper.String(d.Get("domain").(string))
	}
	if d.HasChange("scheme") {
		changeFlag = true
		request.Scheme = helper.String(d.Get("scheme").(string))
	}
	if d.HasChange("comment") {
		changeFlag = true
		request.Comment = helper.String(d.Get("comment").(string))
	}
	if d.HasChange("sub_app_id") {
		changeFlag = true
		request.SubAppId = helper.IntUint64(d.Get("sub_app_id").(int))
	}

	if changeFlag {
		var err error
		err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(request.GetAction())
			_, err = meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVodClient().ModifySuperPlayerConfig(request)
			if err != nil {
				log.Printf("[CRITAL]%s api[%s] fail, reason:%s", logId, request.GetAction(), err.Error())
				return tccommon.RetryError(err)
			}
			return nil
		})
		if err != nil {
			return err
		}

		return resourceTencentCloudVodSuperPlayerConfigRead(d, meta)
	}

	return nil
}

func resourceTencentCloudVodSuperPlayerConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vod_super_player_config.delete")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	id := d.Id()
	vodService := VodService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}

	if err := vodService.DeleteSuperPlayerConfig(ctx, id, uint64(d.Get("sub_app_id").(int))); err != nil {
		return err
	}

	return nil
}
