Provides a resource to create a CLS (Cloud Log Service) search view

Example Usage

```hcl
resource "tencentcloud_cls_search_view_v2" "example" {
  logset_id     = "dac3e1a9-d22c-403b-a129-f94f666a33af"
  logset_region = "ap-guangzhou"
  view_name     = "tf-example-search-view"
  view_type     = "log"

  topics {
    region    = "ap-guangzhou"
    logset_id = "dac3e1a9-d22c-403b-a129-f94f666a33af"
    topic_id  = "775c0bc2-2246-43a0-8eb2-f5bc248be183"
  }

  description     = "terraform example search view"
  view_id_prefix  = "tf-example"
}
```

Import

CLS search view can be imported using the id, e.g.

```
terraform import tencentcloud_cls_search_view_v2.example view-abc123-view
```
