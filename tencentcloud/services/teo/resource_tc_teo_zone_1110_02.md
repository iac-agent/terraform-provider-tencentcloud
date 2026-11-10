Provides a resource to create a zone in TencentCloud EdgeOne (TEO).

Example Usage

```hcl
resource "tencentcloud_teo_zone_1110_02" "example" {
  zone_name = "example.com"
  type      = "full"
  area      = "overseas"
  tags {
    tag_key   = "createdBy"
    tag_value = "terraform"
  }
}

resource "tencentcloud_teo_zone_1110_02" "example_with_plan" {
  zone_name = "example2.com"
  type      = "partial"
  area      = "overseas"
  plan_id   = "plan-xxxxxx"
}
```

Import

teo zone can be imported using the zone id, e.g.

```
terraform import tencentcloud_teo_zone_1110_02.example zone-2noz78a8ev6k
```
