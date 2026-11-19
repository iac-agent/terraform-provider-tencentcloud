# resource_tc_ags_gpu_resource_group

Provides an AGS GPU resource group resource. It supports renaming an existing GPU resource group (`ModifyGPUResourceGroup`) and reading its spec and statistics (`DescribeGPUResourceGroups`).

-> **Note:** The AGS product does not provide public APIs to create or release GPU resource groups. Purchase the resource group in the AGS console first, then import it into Terraform. Because of that, the `create` and `delete` operations of this resource are intentionally rejected with an error.

-> **Note:** `DeleteGPUResourceGroup` only succeeds when the group is empty (no valid GPU instances, occupations or bound sandbox tools). It is therefore not wired into `terraform destroy`; release the resource group in the console after it has been emptied.

## Example Usage

```hcl
resource "tencentcloud_ags_gpu_resource_group" "example" {
  name = "tf-example"
}
```

Import

GPU resource group can be imported using the id, e.g.

```
terraform import tencentcloud_ags_gpu_resource_group.example rg-xxxxxxxx
```

After importing, configure `name` to match the current group name to avoid a spurious diff. Renaming the group with `name` is the only supported update.

## Argument Reference

The following arguments are supported:

* `name` - (Required) Name of the GPU resource group. After trimming whitespace it must not be empty and can be at most 128 bytes long.
* `resource_group_id` - (Computed) Id of the GPU resource group. It is set by `terraform import`.
* `gpu_type` - (Computed) GPU type identifier of the fixed spec, for example `5090d`.
* `gpu_type_alias` - (Computed) Display name of the GPU type. When it is empty the cloud side displays `gpu_type`.
* `gpu_memory_gib` - (Computed) Total memory of all GPUs on one instance, in GiB. `0` means the display metadata is not maintained. It does not need to be multiplied by `gpu_num`.
* `cpu` - (Computed) Number of CPU cores on one instance.
* `memory` - (Computed) Host memory of one instance, in GiB. It is not the GPU memory.
* `gpu_num` - (Computed) Number of GPU cards on one instance.
* `total_count` - (Computed) Total number of resources in the group. Destroyed resources and failed resources whose compensation has completed are excluded.
* `available_count` - (Computed) Number of resources that are Active, unoccupied and physically claimable.
* `allocated_count` - (Computed) Number of resources that are still occupied or bound, including those being claimed, in use or whose release has not finished.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Id of the resource, which is the same as `resource_group_id`.
