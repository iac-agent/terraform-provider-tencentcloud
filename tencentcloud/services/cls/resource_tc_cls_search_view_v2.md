Provides a resource to create a CLS search view.

Example Usage

Create a CLS search view with log type

```hcl
resource "tencentcloud_cls_search_view_v2" "example" {
  logset_id     = "xxxxxxxx-xxxx-xxxx-xxxx"
  logset_region = "ap-guangzhou"
  view_name     = "my-search-view"
  view_type     = "log"

  topics {
    region    = "ap-guangzhou"
    logset_id = "xxxxxxxx-xxxx-xxxx-xxxx"
    topic_id  = "xxxxxxxx-xxxx-xxxx-xxxx"
  }

  description    = "example search view"
  view_id_prefix = "my-prefix"
}
```

Create a CLS search view with metric type

```hcl
resource "tencentcloud_cls_search_view_v2" "metric_example" {
  logset_id     = "xxxxxxxx-xxxx-xxxx-xxxx"
  logset_region = "ap-guangzhou"
  view_name     = "my-metric-view"
  view_type     = "metric"

  topics {
    region    = "ap-guangzhou"
    logset_id = "xxxxxxxx-xxxx-xxxx-xxxx"
    topic_id  = "xxxxxxxx-xxxx-xxxx-xxxx"
  }

  description = "example metric search view"
}
```

Import

CLS search view can be imported using the id, e.g.

```
terraform import tencentcloud_cls_search_view_v2.example my-prefix-view
```