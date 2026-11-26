Use this data source to query the list of AGS sandbox instances.

Example Usage

```hcl
data "tencentcloud_ags_sandbox_instances" "example" {
  instance_ids = ["ags-fake-instance-id"]
}
```

Query all sandbox instances

```hcl
data "tencentcloud_ags_sandbox_instances" "all" {
}
```
