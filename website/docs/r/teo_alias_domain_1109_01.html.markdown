---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_alias_domain_1109_01"
sidebar_current: "docs-tencentcloud-resource-teo_alias_domain_1109_01"
description: |-
  Provides a resource to create a TEO alias domain.
---

# tencentcloud_teo_alias_domain_1109_01

Provides a resource to create a TEO alias domain.

~> **NOTE:** This feature is only supported by the Enterprise edition plan and is currently in beta testing. Please [contact us](https://cloud.tencent.com/online-service?from=connect-us) if you need to use it.

~> **NOTE:** `cert_type` and `cert_id` are not returned by the alias domain query API, so Terraform cannot detect drift on them after creation. Changing them in the configuration does not update the cloud side; recreate the resource instead.

## Example Usage

```hcl
resource "tencentcloud_teo_alias_domain_1109_01" "example" {
  zone_id     = "zone-3fkff38fyw8s"
  alias_name  = "alias.demo.com"
  target_name = "target.demo.com"
  cert_type   = "none"
}
```

### With SSL hosted certificate

```hcl
resource "tencentcloud_teo_alias_domain_1109_01" "example" {
  zone_id     = "zone-3fkff38fyw8s"
  alias_name  = "alias.demo.com"
  target_name = "target.demo.com"
  cert_type   = "hosting"
  cert_id     = ["your-cert-id"]
}
```

## Argument Reference

The following arguments are supported:

* `alias_name` - (Required, String) Alias domain name.
* `target_name` - (Required, String) Target domain name.
* `zone_id` - (Required, String) Site ID.
* `cert_id` - (Optional, Set) Certificate ID list. Required when `cert_type` is `hosting`.
* `cert_type` - (Optional, String) Certificate configuration. Valid values: `none` (no configuration), `hosting` (SSL hosted certificate). Default value: `none`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

      * `alias` - Certificate remark name.
      * `cert_id` - Server certificate ID.
      * `content` - Content of the verification file.
      * `deploy_time` - Certificate deployment time.
      * `expire_time` - Certificate expiration time.
      * `name_servers` - DNS server addresses assigned to the user.
      * `path` - URL path of the verification file.
      * `record_type` - Record type.
      * `record_value` - Record value.
      * `sign_algo` - Signature algorithm.
      * `status` - Certificate status.
      * `subdomain` - Host record.
      * `type` - Certificate type.
    * `apply_type` - Application type.
    * `dns_verification` - DNS resolution verification information.
    * `file_verification` - File verification information.
    * `host_cert_info` - Server certificate configuration.
    * `host` - Domain name.
    * `mode` - Certificate configuration mode.
    * `ns_verification` - NS access verification information.
  * `auth_switch` - Inference task request authentication switch.
  * `certificate` - Domain certificate information.
  * `cname` - CNAME address.
  * `create_time` - Creation time.
  * `domain` - Domain name.
  * `ownership_verification` - Ownership verification information.
  * `status` - Domain status.
  * `update_time` - Modification time.
* `domains` - Inference service domain list returned by DescribeInferenceDomains.
* `id` - Resource identifier in the form `<zone_id>#<alias_name>`.

## API Reference

This resource interacts with the following TencentCloud APIs:

* Create:
  - [CreateAliasDomain](https://www.tencentcloud.com/document/api/1145/51740)
  - [DescribeAliasDomains](https://www.tencentcloud.com/document/api/1145/51738)
  - [DescribeInferenceDomains](https://www.tencentcloud.com/search/?wd=DescribeInferenceDomains&type=document&lang=en)
* Read:
  - [DescribeAliasDomains](https://www.tencentcloud.com/document/api/1145/51738)
  - [DescribeInferenceDomains](https://www.tencentcloud.com/search/?wd=DescribeInferenceDomains&type=document&lang=en)
* Update:
  - [DescribeAliasDomains](https://www.tencentcloud.com/document/api/1145/51738)
  - [DescribeInferenceDomains](https://www.tencentcloud.com/search/?wd=DescribeInferenceDomains&type=document&lang=en)
* Delete: [DeleteAliasDomain](https://www.tencentcloud.com/document/api/1145/51739)

## Import

TEO alias domains can be imported using the `id`, e.g.

```
$ terraform import tencentcloud_teo_alias_domain_1109_01.example zone-3fkff38fyw8s#alias.demo.com
```

