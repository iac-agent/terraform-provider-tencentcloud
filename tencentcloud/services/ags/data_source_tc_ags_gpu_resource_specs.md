Use this data source to query the sellable GPU resource specs and the available quota of the current tenant in the current region.

-> **Note:** `gpu_type`, `cpu`, `memory` and `gpu_num` must be provided together or omitted together; when they are provided, filtering matches all four items exactly.

-> **Note:** `available_count` is measured in instances, not GPU cards. 0 means no quota and null means the quota is unknown, so it is mapped to null instead of 0. The quota is verified again at purchase time.

Example Usage

Query all sellable GPU specs

```hcl
data "tencentcloud_ags_gpu_resource_specs" "example" {}
```

Query a specific GPU spec precisely

```hcl
data "tencentcloud_ags_gpu_resource_specs" "example" {
  gpu_type = "5090d"
  cpu      = 64
  memory   = 192
  gpu_num  = 2
}
```

Query with pagination

```hcl
data "tencentcloud_ags_gpu_resource_specs" "example" {
  offset = 0
  limit  = 20
}
```
