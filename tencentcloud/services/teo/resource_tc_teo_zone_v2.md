Provides a resource to create a TEO (EdgeOne) zone

Example Usage

Basic Usage

```hcl
resource "tencentcloud_teo_zone_v2" "zone" {
  type            = "partial"
  zone_name       = "tf-teo.com"
  area            = "global"
  plan_id         = "edgeone-2kfv1h391n6w"
  alias_zone_name = "teo-test"
  allow_duplicates = false
  jump_start      = false

  tags {
    tag_key   = "createdBy"
    tag_value = "terraform"
  }
}
```

Use Custom Name Servers

```hcl
resource "tencentcloud_teo_zone_v2" "zone" {
  type            = "full"
  zone_name       = "tf-teo-vanity.com"
  area            = "global"
  plan_id         = "edgeone-2kfv1h391n6w"
  alias_zone_name = "teo-vanity-test"

  vanity_name_servers {
    switch  = "on"
    servers = ["ns1.example.com", "ns2.example.com"]
  }

  tags {
    tag_key   = "createdBy"
    tag_value = "terraform"
  }
}
```

Import

TEO (EdgeOne) zone can be imported using the zone id, e.g.
```
terraform import tencentcloud_teo_zone_v2.zone zone-2noz78a8ev6k
```
