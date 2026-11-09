Provides a resource to create a TEO alias domain.

~> **NOTE:** This feature is only supported by the Enterprise edition plan and is currently in beta testing. Please [contact us](https://cloud.tencent.com/online-service?from=connect-us) if you need to use it.

~> **NOTE:** `cert_type` and `cert_id` are not returned by the `DescribeAliasDomains` API, so Terraform cannot detect drift on these two fields. Removing either of them from the configuration keeps the previous cloud-side configuration; to remove a hosted certificate binding, set `cert_type` to `none` explicitly.

Example Usage

```hcl
resource "tencentcloud_teo_alias_domain_1109_03" "example" {
  zone_id     = "zone-2qtuhspy7cr6"
  alias_name  = "alias.demo.com"
  target_name = "target.demo.com"
  cert_type   = "none"
}
```

With SSL hosted certificate

```hcl
resource "tencentcloud_teo_alias_domain_1109_03" "example" {
  zone_id     = "zone-2qtuhspy7cr6"
  alias_name  = "alias.demo.com"
  target_name = "target.demo.com"
  cert_type   = "hosting"
  cert_id     = ["your-cert-id"]
}
```

Import

TEO alias domain can be imported using the composite id `<zone_id>#<alias_name>`, e.g.

```
terraform import tencentcloud_teo_alias_domain_1109_03.example zone-2qtuhspy7cr6#alias.demo.com
```
