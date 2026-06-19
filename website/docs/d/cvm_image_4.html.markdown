---
subcategory: "Cloud Virtual Machine(CVM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cvm_image_4"
sidebar_current: "docs-tencentcloud-datasource-cvm_image_4"
description: |-
  Use this data source to query detailed information of CVM image
---

# tencentcloud_cvm_image_4

Use this data source to query detailed information of CVM image

## Example Usage

```hcl
data "tencentcloud_cvm_image_4" "example" {
  image_ids = ["img-0elsru2u"]

  filters {
    name   = "image-type"
    values = ["PRIVATE_IMAGE"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `filters` - (Optional, List) Filter conditions.
* `image_ids` - (Optional, List: [`String`]) Image ID list.
* `instance_type` - (Optional, String) Instance type, such as `SA5.MEDIUM2`.
* `result_output_file` - (Optional, String) Used to save results.

The `filters` object supports the following:

* `name` - (Required, String) Filter field name.
* `values` - (Required, List) Filter field value.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `image_set` - A structure about image details, including the main status and attributes of the image.
  * `architecture` - Image architecture. Valid values include: `x86_64`, `arm`, `i386`.
  * `cdc_cache_status` - CDC image cache status.
  * `created_time` - Image creation time. Format: YYYY-MM-DDThh:mm:ssZ (ISO8601 standard, UTC time).
  * `image_creator` - Image creator.
  * `image_deprecated` - Whether the image is deprecated.
  * `image_description` - Image description.
  * `image_family` - Image family.
  * `image_id` - Image ID.
  * `image_name` - Image name.
  * `image_size` - Image size in GiB.
  * `image_source` - Image source. Valid values include: `OFFICIAL`, `CREATE_IMAGE`, `EXTERNAL_IMPORT`.
  * `image_state` - Image state. Valid values: CREATING, NORMAL, CREATEFAILED, SYNCING, IMPORTING, IMPORTFAILED.
  * `image_type` - Image type. Valid values include: `PUBLIC_IMAGE`, `PRIVATE_IMAGE`, `SHARED_IMAGE`.
  * `is_support_cloudinit` - Whether the image supports cloud-init.
  * `license_type` - Image license type. Valid values include: `TencentCloud`, `BYOL`.
  * `os_name` - Image OS name.
  * `platform` - Image source platform, including TencentOS, CentOS, Windows, Ubuntu, Debian, Fedora, etc.
  * `snapshot_set` - Snapshot information associated with the image.
    * `disk_size` - Cloud disk size for this snapshot in GiB.
    * `disk_usage` - Cloud disk type for this snapshot. Valid values: SYSTEM_DISK, DATA_DISK.
    * `snapshot_id` - Snapshot ID.
  * `sync_percent` - Sync percentage. Note: This field may return null, indicating that no valid value can be obtained.
  * `tags` - List of tags associated with the image.
    * `key` - Tag key.
    * `value` - Tag value.


