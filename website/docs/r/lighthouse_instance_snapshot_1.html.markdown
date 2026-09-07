---
subcategory: "TencentCloud Lighthouse(Lighthouse)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_lighthouse_instance_snapshot_1"
sidebar_current: "docs-tencentcloud-resource-lighthouse_instance_snapshot_1"
description: |-
  Provides a resource to create a Lighthouse instance snapshot, with tags support and comprehensive computed snapshot attributes.
---

# tencentcloud_lighthouse_instance_snapshot_1

Provides a resource to create a Lighthouse instance snapshot, with tags support and comprehensive computed snapshot attributes.

## Example Usage

```hcl
resource "tencentcloud_lighthouse_instance_snapshot_1" "instance_snapshot" {
  instance_id   = "lhins-xxxxxx"
  snapshot_name = "snap_20240101"

  tags {
    key   = "env"
    value = "production"
  }
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, ForceNew) ID of the instance for which to create a snapshot.
* `snapshot_name` - (Optional, String) Snapshot name, which can contain up to 60 characters.
* `tags` - (Optional, List) Tags to associate with the snapshot.

The `tags` object supports the following:

* `key` - (Required, String) Tag key.
* `value` - (Required, String) Tag value.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `created_time` - Snapshot creation time.
* `disk_id` - ID of the disk for which the snapshot was created.
* `disk_size` - Size of the disk for which the snapshot was created, in GB.
* `disk_usage` - Disk type for which the snapshot was created. Values: `SYSTEM_DISK` (system disk).
* `latest_operation_request_id` - The unique request ID for the latest operation of the snapshot. Only recorded when creating or rolling back the snapshot.
* `latest_operation_state` - Latest operation state of the snapshot. Only recorded when creating or rolling back the snapshot. Values: `SUCCESS`, `OPERATING`, `FAILED`.
* `latest_operation` - Latest operation of the snapshot. Only recorded when creating or rolling back the snapshot. Values: `CreateInstanceSnapshot`, `RollbackInstanceSnapshot`.
* `percent` - Snapshot creation progress percentage. After the snapshot is created successfully, the value of this field will be 100.
* `snapshot_id` - Snapshot ID.
* `snapshot_state` - Snapshot state. Values: `NORMAL`, `CREATING`, `ROLLBACKING`.


## Import

Lighthouse instance snapshot can be imported using the id, e.g.

```
terraform import tencentcloud_lighthouse_instance_snapshot_1.instance_snapshot lhsnap-xxxxxx
```

