Provides a resource to create a TEO just in time transcode template

Example Usage

```hcl
resource "tencentcloud_teo_just_in_time_transcode_template" "example" {
  zone_id             = "zone-3edjdliiw3he"
  template_name       = "tf-example"
  comment             = "test comment"
  video_stream_switch = "on"
  audio_stream_switch = "on"
  video_template {
    codec              = "H.264"
    fps                = 30
    bitrate            = 2000
    resolution_adaptive = "open"
    width              = 1920
    height             = 1080
    fill_type          = "black"
  }
  audio_template {
    codec        = "libfdk_aac"
    audio_channel = 2
  }
}
```

Import

TEO just in time transcode template can be imported using the id, e.g.

```
terraform import tencentcloud_teo_just_in_time_transcode_template.example zone-3edjdliiw3he#C1LZ7982VgTpYhJ7M
```
