---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_import_zone_config_operation"
sidebar_current: "docs-tencentcloud-resource-teo_import_zone_config_operation"
description: |-
  Provides a resource to import TEO (EdgeOne) zone configuration.
---

# tencentcloud_teo_import_zone_config_operation

Provides a resource to import TEO (EdgeOne) zone configuration.

## Example Usage

```hcl
resource "tencentcloud_teo_import_zone_config_operation" "example" {
  zone_id = "zone-2qtuhspy7cr6"
  content = jsonencode({
    ZoneId = "zone-2qtuhspy7cr6"
    CacheConfig = {
      CustomTime = {
        Switch    = "on"
        CacheTime = 86400
      }
    }
  })
}
```

## Argument Reference

The following arguments are supported:

* `content` - (Required, String, ForceNew) The configuration content to import. It must be in JSON format and UTF-8 encoded. The content can be obtained via the ExportZoneConfig API.
* `zone_id` - (Required, String, ForceNew) Zone ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `task_id` - The task ID returned by the API, used for polling the import result.


