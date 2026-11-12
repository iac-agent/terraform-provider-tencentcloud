Provides a resource to create an AGS pre-cache image task.

-> **Note:** `image` and `image_registry_type` cannot be changed after creation: the cloud API has no update interface for this resource, so changing either of them destroys the existing pre-cache image task and creates a new one.

-> **Note:** Creating this resource waits until the pre-cache image task reports a ready time. The wait stops with an error when the timeout is reached; the task is kept in state so that `terraform destroy` can still remove it.

-> **Note:** Destroying this resource removes the pre-cache image entry itself only. It does not terminate running sandbox instances and it does not touch the source image registry.

Example Usage

```hcl
resource "tencentcloud_ags_pre_cache_image_task" "example" {
  image               = "ccr.ccs.tencentyun.com/tf-example/demo:v1"
  image_registry_type = "personal"
}
```

Import

AGS pre-cache image task can be imported using the pre-cache image id, e.g.

```
$ terraform import tencentcloud_ags_pre_cache_image_task.example cache-xxxxxxxx
```