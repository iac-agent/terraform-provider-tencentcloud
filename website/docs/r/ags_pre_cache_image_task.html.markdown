---
subcategory: "Agent Sandbox(ags)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ags_pre_cache_image_task"
sidebar_current: "docs-tencentcloud-resource-ags_pre_cache_image_task"
description: |-
  Provides a resource to create an AGS pre-cache image task.
---

# tencentcloud_ags_pre_cache_image_task

Provides a resource to create an AGS pre-cache image task.

-> **Note:** `image` and `image_registry_type` cannot be changed after creation: the cloud API has no update interface for this resource, so changing either of them destroys the existing pre-cache image task and creates a new one.

-> **Note:** Creating this resource waits until the pre-cache image task reports a ready time. The wait stops with an error when the timeout is reached; the task is kept in state so that `terraform destroy` can still remove it.

-> **Note:** Destroying this resource removes the pre-cache image entry itself only. It does not terminate running sandbox instances and it does not touch the source image registry.

## Example Usage

```hcl
resource "tencentcloud_ags_pre_cache_image_task" "example" {
  image               = "ccr.ccs.tencentyun.com/tf-example/demo:v1"
  image_registry_type = "personal"
}
```

## Argument Reference

The following arguments are supported:

* `image_registry_type` - (Required, String) Image registry type. Valid values: `enterprise`, `personal`, `custom`. `enterprise` is the TCR enterprise image registry and `personal` is the CCR personal image registry. It cannot be changed after creation.
* `image` - (Required, String) Image address. Only `repository:tag`, `repository@sha256:<64 hex digits>` and `repository:tag@sha256:<64 hex digits>` are supported. It cannot be changed after creation.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `cached_image_size_bytes` - Storage size of the pre-cached image, in bytes.
* `create_time` - Creation time of the pre-cache image task, reported by the cloud side as a string.
* `id` - Id of the resource, which is the same as `pre_cache_image_id`.
* `image_digest` - Digest of the image, resolved and reported by the cloud side.
* `last_used_time` - Time when the pre-cached image was last used by a sandbox instance. It changes while the resource is managed, so it is only reported by the cloud side.
* `message` - Description of the pre-cache image task status.
* `pre_cache_image_id` - Id of the pre-cache image task.
* `ready_time` - Time when the pre-cache image task completed, reported by the cloud side as a string.
* `source_type` - Source type of the pre-cache image resource. Valid values: `EXPLICIT` (created manually), `AUTO` (created automatically), `TCR_AUTO` (pre-cached automatically by TCR).
* `status` - Status of the pre-cache image task.

## API Reference

This resource interacts with the following TencentCloud APIs:

* Create:
  - [CreatePreCacheImageTask](https://www.tencentcloud.com/search/?wd=CreatePreCacheImageTask&type=document&lang=en)
  - [DescribePreCacheImageTask](https://www.tencentcloud.com/search/?wd=DescribePreCacheImageTask&type=document&lang=en)
* Read: [DescribePreCacheImageTask](https://www.tencentcloud.com/search/?wd=DescribePreCacheImageTask&type=document&lang=en)
* Update:
  - [DescribeSandboxToolList](https://www.tencentcloud.com/search/?wd=DescribeSandboxToolList&type=document&lang=en)
  - [UpdateSandboxTool](https://www.tencentcloud.com/search/?wd=UpdateSandboxTool&type=document&lang=en)
* Delete:
  - [DeletePreCacheImage](https://www.tencentcloud.com/search/?wd=DeletePreCacheImage&type=document&lang=en)
  - [DescribePreCacheImageTask](https://www.tencentcloud.com/search/?wd=DescribePreCacheImageTask&type=document&lang=en)

## Import

AGS pre-cache image task can be imported using the pre-cache image id, e.g.

```
$ terraform import tencentcloud_ags_pre_cache_image_task.example cache-xxxxxxxx
```

