Provides a resource to create a TEO alias domain.

~> **NOTE:** This feature is only supported by the Enterprise edition plan and is currently in beta testing. Please [contact us](https://cloud.tencent.com/online-service?from=connect-us) if you need to use it.

~> **NOTE:** `cert_type` and `cert_id` are not returned by the alias domain query API, so Terraform cannot detect drift on them after creation. Changing them in the configuration does not update the cloud side; recreate the resource instead.

Example Usage

```hcl
resource "tencentcloud_teo_alias_domain_1109_01" "example" {
  zone_id     = "zone-3fkff38fyw8s"
  alias_name  = "alias.demo.com"
  target_name = "target.demo.com"
  cert_type   = "none"
}
```

With SSL hosted certificate

```hcl
resource "tencentcloud_teo_alias_domain_1109_01" "example" {
  zone_id     = "zone-3fkff38fyw8s"
  alias_name  = "alias.demo.com"
  target_name = "target.demo.com"
  cert_type   = "hosting"
  cert_id     = ["your-cert-id"]
}
```

Import

TEO alias domains can be imported using the `id`, e.g.

```
$ terraform import tencentcloud_teo_alias_domain_1109_01.example zone-3fkff38fyw8s#alias.demo.com
```
