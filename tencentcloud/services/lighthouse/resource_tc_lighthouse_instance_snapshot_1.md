Provides a resource to create a Lighthouse instance snapshot, with tags support and comprehensive computed snapshot attributes.

Example Usage

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

Import

Lighthouse instance snapshot can be imported using the id, e.g.

```
terraform import tencentcloud_lighthouse_instance_snapshot_1.instance_snapshot lhsnap-xxxxxx
```