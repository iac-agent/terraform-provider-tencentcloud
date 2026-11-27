Use this data source to query the list of AGS sandbox instances.

Example Usage

Query all sandbox instances:

```hcl
data "tencentcloud_ags_sandbox_instances" "all" {
}
```

Query sandbox instances by ID:

```hcl
data "tencentcloud_ags_sandbox_instances" "filtered" {
  instance_ids = ["sandbox-8f0b1c2d", "sandbox-9a1b2c3d"]
}
```
