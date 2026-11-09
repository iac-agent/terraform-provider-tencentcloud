Provides a resource to create a TEO alias domain.

~> **NOTE:** An alias domain shares the plan of its target domain. Only the alias domain mapping itself is managed here: deleting this resource removes the alias domain only, and does not affect the target domain or the zone.

-> **NOTE:** `cert_type` and `cert_id` are write-only on the cloud side (the query API does not return them), so Terraform cannot detect drift on these two attributes; their configured values are kept in state.

Example Usage

```hcl
resource "tencentcloud_teo_alias_domain_1109_04" "example" {
  zone_id     = "zone-3fkff38fyw8s"
  alias_name  = "alias.demo.com"
  target_name = "target.demo.com"
}
```

With an SSL hosted certificate

```hcl
resource "tencentcloud_teo_alias_domain_1109_04" "example" {
  zone_id     = "zone-3fkff38fyw8s"
  alias_name  = "alias.demo.com"
  target_name = "target.demo.com"
  cert_type   = "hosting"
  cert_id     = ["your-cert-id"]
}
```

Import

TEO alias domains can be imported using the `zone_id#alias_name` composite ID, e.g.
```
$ terraform import tencentcloud_teo_alias_domain_1109_04.example zone-3fkff38fyw8s#alias.demo.com
```
