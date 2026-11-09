Provides a resource to create a TEO (EdgeOne) zone.

Example Usage

Create a CNAME-access zone bound to an existing plan

```hcl
resource "tencentcloud_teo_zone_1109_01" "example" {
  type            = "partial"
  zone_name       = "example.com"
  area            = "overseas"
  plan_id         = "edgeone-2kfv1h391n6w"
  alias_zone_name = "my-alias"

  tags {
    tag_key   = "createdBy"
    tag_value = "terraform"
  }
}
```

Create an NS-access zone with custom name servers updated afterwards

```hcl
resource "tencentcloud_teo_zone_1109_01" "example" {
  type      = "full"
  zone_name = "example.com"
  area      = "global"
  plan_id   = "edgeone-2kfv1h391n6w"

  vanity_name_servers = {
    switch  = "on"
    servers = ["ns1.example.com", "ns2.example.com"]
  }
}
```

Create a no-domain-access zone (keep `zone_name` and `area` unset)

```hcl
resource "tencentcloud_teo_zone_1109_01" "example" {
  type    = "noDomainAccess"
  plan_id = "edgeone-2kfv1h391n6w"
}
```

Import

TEO zones can be imported using the zone id, e.g.

```
terraform import tencentcloud_teo_zone_1109_01.example zone-2noz78a8ev6k
```

Notes

- `plan_id`, `tags`, `allow_duplicates` and `jump_start` are only accepted at
  creation (ModifyZone has no counterpart for them); changing any of them
  forces the zone to be recreated.
- `type`, `zone_name`, `area`, `alias_zone_name` and `vanity_name_servers` are
  updated in place via ModifyZone. ModifyZone keeps the cloud-side value for
  omitted fields, so removing one of them from the configuration keeps the
  current cloud value instead of clearing it.
- `allow_duplicates` and `jump_start` are not returned by DescribeZones, so
  drift detection for these two fields is limited.
- After creation the zone still needs ownership verification before it serves
  traffic; use the computed `ownership_verification` attribute for the TXT
  record or file content, or switch the domain NS servers to the returned
  `name_servers` for NS access.
