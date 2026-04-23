Provides a resource to import TEO (EdgeOne) zone configuration.

Example Usage

```hcl
resource "tencentcloud_teo_import_zone_config_operation" "example" {
  zone_id = "zone-2qtuhspy7cr6"
  content = jsonencode({
    ZoneId = "zone-2qtuhspy7cr6"
    CacheConfig = {
      CustomTime = {
        Switch = "on"
        CacheTime = 86400
      }
    }
  })
}
```

Argument Reference

The following arguments are supported:

* `zone_id` - (Required, ForceNew) Zone ID.
* `content` - (Required, ForceNew) The configuration content to import. It must be in JSON format and UTF-8 encoded. The content can be obtained via the ExportZoneConfig API.

Attributes Reference

The following attributes are exported:

* `task_id` - The task ID returned by the API, used for polling the import result.
