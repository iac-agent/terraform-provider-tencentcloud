---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_just_in_time_transcode_template"
sidebar_current: "docs-tencentcloud-resource-teo_just_in_time_transcode_template"
description: |-
  Provides a resource to create a TEO just in time transcode template
---

# tencentcloud_teo_just_in_time_transcode_template

Provides a resource to create a TEO just in time transcode template

## Example Usage

```hcl
resource "tencentcloud_teo_just_in_time_transcode_template" "example" {
  zone_id             = "zone-3edjdliiw3he"
  template_name       = "tf-example"
  comment             = "test comment"
  video_stream_switch = "on"
  audio_stream_switch = "on"
  video_template {
    codec               = "H.264"
    fps                 = 30
    bitrate             = 2000
    resolution_adaptive = "open"
    width               = 1920
    height              = 1080
    fill_type           = "black"
  }
  audio_template {
    codec         = "libfdk_aac"
    audio_channel = 2
  }
}
```

## Argument Reference

The following arguments are supported:

* `template_name` - (Required, String, ForceNew) Just-in-time transcode template name.
* `zone_id` - (Required, String, ForceNew) Zone ID.
* `audio_stream_switch` - (Optional, String, ForceNew) Enable audio stream switch, valid values: on, off.
* `audio_template` - (Optional, List, ForceNew) Audio stream configuration parameters.
* `comment` - (Optional, String, ForceNew) Template description.
* `video_stream_switch` - (Optional, String, ForceNew) Enable video stream switch, valid values: on, off.
* `video_template` - (Optional, List, ForceNew) Video stream configuration parameters.

The `audio_template` object supports the following:

* `audio_channel` - (Optional, Int, ForceNew) Audio channel count, valid value: 2.
* `codec` - (Optional, String, ForceNew) Audio stream encoding format, valid value: libfdk_aac.

The `video_template` object supports the following:

* `bitrate` - (Optional, Int, ForceNew) Video stream bitrate, range: 0 and [128, 10000].
* `codec` - (Optional, String, ForceNew) Video stream encoding format, valid values: H.264, H.265.
* `fill_type` - (Optional, String, ForceNew) Fill type, valid values: stretch, black, white, gauss.
* `fps` - (Optional, Float64, ForceNew) Video frame rate, range: [0, 30].
* `height` - (Optional, Int, ForceNew) Maximum height of video stream, range: 0 and [128, 1080].
* `resolution_adaptive` - (Optional, String, ForceNew) Resolution adaptive, valid values: open, close.
* `width` - (Optional, Int, ForceNew) Maximum width of video stream, range: 0 and [128, 1920].

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - Template creation time.
* `template_id` - Just-in-time transcode template unique identifier.
* `type` - Template type, valid values: preset, custom.
* `update_time` - Template last modified time.


## Import

TEO just in time transcode template can be imported using the id, e.g.

```
terraform import tencentcloud_teo_just_in_time_transcode_template.example zone-3edjdliiw3he#C1LZ7982VgTpYhJ7M
```

