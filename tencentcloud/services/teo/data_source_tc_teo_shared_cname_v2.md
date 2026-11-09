Use this data source to query the shared CNAME list of TEO (TencentCloud EdgeOne) within a zone.

Example Usage

```hcl
data "tencentcloud_teo_shared_cname_v2" "example" {
  zone_id = "zone-2qtuhspy7cr6"
}

# Query with filters and sorting.
data "tencentcloud_teo_shared_cname_v2" "filtered" {
  zone_id = "zone-2qtuhspy7cr6"

  filters {
    name   = "shared-cname"
    values = ["prefix-1.2qtuhspy7cr6.share.dnse1.com"]
  }

  filters {
    name   = "description"
    values = ["demo"]
    fuzzy  = true
  }

  match     = "any"
  direction = "desc"
  order     = "create-time"
}
```
