package vod

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
)

const (
	VOD_AUDIO_CHANNEL_MONO   = "mono"
	VOD_AUDIO_CHANNEL_DUAL   = "dual"
	VOD_AUDIO_CHANNEL_STEREO = "stereo"

	VOD_SUB_APPLICATION_RUNNING = "On"
	VOD_SUB_APPLICATION_STOPPED = "Off"
	VOD_SUB_APPLICATION_DESTROY = "Destroyed"

	VOD_DEFAULT_OFFSET = 0
	VOD_MAX_LIMIT      = 100
)

var VOD_SUB_APPLICATION_STATUS = []string{
	VOD_SUB_APPLICATION_RUNNING,
	VOD_SUB_APPLICATION_STOPPED,
	VOD_SUB_APPLICATION_DESTROY,
}

var (
	VOD_AUDIO_CHANNEL_TYPE_TO_INT = map[string]int64{
		VOD_AUDIO_CHANNEL_MONO:   1,
		VOD_AUDIO_CHANNEL_DUAL:   2,
		VOD_AUDIO_CHANNEL_STEREO: 6,
	}
	VOD_AUDIO_CHANNEL_TYPE_TO_STRING = map[int64]string{
		1: VOD_AUDIO_CHANNEL_MONO,
		2: VOD_AUDIO_CHANNEL_DUAL,
		6: VOD_AUDIO_CHANNEL_STEREO,
	}
	DISABLE_HIGHER_VIDEO_BITRATE_TO_UNINT = map[bool]uint64{
		true:  1,
		false: 0,
	}
	DISABLE_HIGHER_VIDEO_RESOLUTION_TO_UNINT = map[bool]uint64{
		true:  1,
		false: 0,
	}
	RESOLUTION_ADAPTIVE_TO_STRING = map[bool]string{
		true:  "open",
		false: "close",
	}
	REMOVE_AUDIO_TO_UNINT = map[bool]uint64{
		true:  1,
		false: 0,
	}
	DRM_SWITCH_TO_STRING = map[bool]string{
		true:  "ON",
		false: "OFF",
	}
)

func VodWatermarkResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"definition": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Watermarking 模板 ID",
			},
			"text_content": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: tccommon.ValidateStringLengthInRange(0, 100),
				Description:  "Text 内容 的 up 到 `100` 字符. 此 needs 到 是 entered 仅 当 水印 类型 是 text. 注意: 此 字段 可能 返回 null，indicating 该 无 有效 值 可以 是 获取。",
			},
			"svg_content": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: tccommon.ValidateStringLengthInRange(0, 2000000),
				Description:  "SVG 内容 的 up 到 `2000000` 字符. 此 needs 到 是 entered 仅 当 水印 类型 是 `SVG`. 注意: 此 字段 可能 返回 null，indicating 该 无 有效 值 可以 是 获取。",
			},
			"start_time_offset": {
				Type:        schema.TypeFloat,
				Optional:    true,
				Description: "开始时间 偏移量 的 水印 （秒）。 如果 此 参数 是 left blank 或 `0` 是 entered， 水印 将 appear upon first 视频 frame. 如果 此 参数 是 left blank 或 `0` 是 entered， 水印 将 appear upon first 视频 frame; 如果 此 值 是 greater 比 `0` (e.g.，n)， 水印 将 appear 在 second n after first 视频 frame; 如果 此 值 是 smaller 比 `0` (e.g.，-n)， 水印 将 appear 在 second n before last 视频 frame。",
			},
			"end_time_offset": {
				Type:        schema.TypeFloat,
				Optional:    true,
				Description: "结束时间 偏移量 的 水印 （秒）。 如果 此 参数 是 left blank 或 `0` 是 entered， 水印 将 exist till last 视频 frame; 如果 此 值 是 greater 比 `0` (e.g.，n)， 水印 将 exist till second n; 如果 此 值 是 smaller 比 `0` (e.g.，-n)， 水印 将 exist till second n before last 视频 frame。",
			},
		},
	}
}
