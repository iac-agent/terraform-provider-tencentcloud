Provides a resource to create a TEO zone.

~> **NOTE:** A freshly created zone only serves traffic after a plan is bound to it (`plan_id`) and after the ownership verification returned by `ownership_verification` is finished. When `plan_id` is omitted the zone is created in the `init` state and it does not show up in the console until a plan is bound to it.

~> **NOTE:** `allow_duplicates` and `jump_start` are create-only switches that the TEO `DescribeZones` API never returns. Their configured values are kept as-is in state, so out-of-band changes of these two attributes cannot be detected. Changing either of them replaces the zone.

Example Usage

```hcl
resource "tencentcloud_teo_zone_1110_01" "example" {
  type            = "partial"
  zone_name       = "example.com"
  area            = "overseas"
  plan_id         = "edgeone-2wdo315m2y4c"
  alias_zone_name = "example-alias"
  allow_duplicates = false
  jump_start       = false

  tags {
    tag_key   = "env"
    tag_value = "test"
  }
}
```

### Zone without a domain

```hcl
resource "tencentcloud_teo_zone_1110_01" "no_domain" {
  type      = "noDomainAccess"
  plan_id   = "edgeone-2wdo315m2y4c"
  jump_start = true
}
```

Import

TEO zone can be imported using the zone ID, e.g.

```
terraform import tencentcloud_teo_zone_1110_01.example zone-2noz78a8ev6k
```

~> **NOTE:** `ownership_verification` is a creation-time output, so it stays empty on an imported zone.
