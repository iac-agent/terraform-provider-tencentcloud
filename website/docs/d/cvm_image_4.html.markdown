---
subcategory: "Cloud Virtual Machine(CVM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cvm_image_4"
sidebar_current: "docs-tencentcloud-datasource-cvm_image_4"
description: |-
  Use this data source to query detailed information of CVM images
---

# tencentcloud_cvm_image_4

Use this data source to query detailed information of CVM images

## Example Usage

### Query by image ids

```hcl
data "tencentcloud_cvm_image_4" "example" {
  image_ids = ["img-8toqc6s3"]
}
```

### Query by filters

```hcl
data "tencentcloud_cvm_image_4" "example" {
  filters {
    name   = "image-type"
    values = ["PUBLIC_IMAGE"]
  }
}
```

### Query by instance type

```hcl
data "tencentcloud_cvm_image_4" "example" {
  instance_type = "SA5.MEDIUM2"
  filters {
    name   = "image-type"
    values = ["PUBLIC_IMAGE"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `filters` - (Optional, List) Filter conditions for the query. Mutually exclusive with `image_ids`.
* `image_ids` - (Optional, List: [`String`]) List of image IDs to query. Mutually exclusive with `filters`.
* `instance_type` - (Optional, String) Instance type for compatibility check, e.g., `SA5.MEDIUM2`.
* `result_output_file` - (Optional, String) Used to save results.

The `filters` object supports the following:

* `name` - (Required, String) Filter field name.
* `values` - (Required, Set) Filter field values.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `image_set` - List of images.
  * `architecture` - Architecture, e.g., `x86_64`, `arm`, `i386`.
  * `cdc_cache_status` - CDC image cache status.
  * `created_time` - Image creation time.
  * `image_creator` - Image creator.
  * `image_deprecated` - Whether the image is deprecated.
  * `image_description` - Image description.
  * `image_family` - Image family.
  * `image_id` - Image ID.
  * `image_name` - Image name.
  * `image_size` - Image size in GiB.
  * `image_source` - Image source, e.g., `OFFICIAL`, `CREATE_IMAGE`, `EXTERNAL_IMPORT`.
  * `image_state` - Image state, e.g., `CREATING`, `NORMAL`, `CREATEFAILED`, `SYNCING`, `IMPORTING`, `IMPORTFAILED`.
  * `image_type` - Image type, e.g., `PUBLIC_IMAGE`, `PRIVATE_IMAGE`, `SHARED_IMAGE`.
  * `is_support_cloudinit` - Whether cloud-init is supported.
  * `license_type` - License type, e.g., `TencentCloud`, `BYOL`.
  * `os_name` - OS name of the image.
  * `platform` - Source platform.
  * `snapshot_set` - Snapshot list of the image.
    * `disk_size` - Disk size in GiB.
    * `disk_usage` - Disk usage.
    * `snapshot_id` - Snapshot ID.
  * `sync_percent` - Sync percentage.
  * `tags` - Tag list of the image.
    * `key` - Tag key.
    * `value` - Tag value.


