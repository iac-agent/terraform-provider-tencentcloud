Use this data source to query detailed information of TEO origin acl

The `origin_acl_info` block contains the following computed fields:
- `status` - Origin ACL status. Valid values: `online`, `offline`, `updating`.

Example Usage

Query origin acl by zone Id

```hcl
data "tencentcloud_teo_origin_acl" "example" {
  zone_id = "zone-3fkff38fyw8s"
}
```
