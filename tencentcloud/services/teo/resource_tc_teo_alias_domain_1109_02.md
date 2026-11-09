Provides a resource to create a TEO alias domain. An alias domain points to a target domain within the same site and shares its acceleration configuration.

~> **NOTE:** This resource is implemented with terraform-plugin-framework. `cert_type` and `cert_id` are not returned by the DescribeAliasDomains API, so their values cannot be drift-detected after apply: they are kept from the configuration, and out-of-band changes on the cloud side are not detected.

Example Usage

Basic

```hcl
resource "tencentcloud_teo_alias_domain_1109_02" "example" {
  zone_id     = "zone-3fkff38fyw8s"
  alias_name  = "alias.demo.com"
  target_name = "target.demo.com"
}
```

With SSL hosted certificate

```hcl
resource "tencentcloud_teo_alias_domain_1109_02" "example" {
  zone_id     = "zone-3fkff38fyw8s"
  alias_name  = "alias.demo.com"
  target_name = "target.demo.com"
  cert_type   = "hosting"
  cert_id     = ["your-cert-id"]
}
```

Switch the certificate configuration back to none explicitly

```hcl
resource "tencentcloud_teo_alias_domain_1109_02" "example" {
  zone_id     = "zone-3fkff38fyw8s"
  alias_name  = "alias.demo.com"
  target_name = "target.demo.com"
  cert_type   = "none"
}
```

Import

TEO alias domains can be imported using the composite id `<zone_id>#<alias_name>`, e.g.

```
$ terraform import tencentcloud_teo_alias_domain_1109_02.example zone-3fkff38fyw8s#alias.demo.com
```

When destroying this resource, only the alias domain itself is deleted; the target domain and any hosted certificates are left untouched. Removing `cert_type` or `cert_id` from the configuration keeps the previous certificate configuration on the cloud side (the documented update semantics); set `cert_type = "none"` explicitly to switch to no certificate.
