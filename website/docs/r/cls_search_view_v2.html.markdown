---
subcategory: "Cloud Log Service(CLS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cls_search_view_v2"
sidebar_current: "docs-tencentcloud-resource-cls_search_view_v2"
description: |-
  Provides a resource to create a CLS (Cloud Log Service) search view
---

# tencentcloud_cls_search_view_v2

Provides a resource to create a CLS (Cloud Log Service) search view

## Example Usage

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

  description    = "terraform example search view"
  view_id_prefix = "tf-example"
}
```

## Argument Reference

The following arguments are supported:

* `logset_id` - (Required, String, ForceNew) Logset ID to which the search view belongs.
* `logset_region` - (Required, String, ForceNew) Region of the logset, e.g., ap-guangzhou.
* `topics` - (Required, List) Topics included in the search view, max 10 topics.
* `view_name` - (Required, String) Search view name, max 255 characters.
* `view_type` - (Required, String) Search view type. Valid values: log, metric.
* `description` - (Optional, String) Description of the search view.
* `view_id_prefix` - (Optional, String, ForceNew) Custom search view ID prefix.

The `topics` object supports the following:

* `logset_id` - (Required, String) Logset ID of the topic.
* `region` - (Required, String) Region of the topic.
* `topic_id` - (Required, String) Topic ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - Creation time.
* `update_time` - Last update time.
* `view_id` - Search view ID.


## Import

CLS search view can be imported using the id, e.g.

```
terraform import tencentcloud_cls_search_view_v2.example view-abc123-view
```

