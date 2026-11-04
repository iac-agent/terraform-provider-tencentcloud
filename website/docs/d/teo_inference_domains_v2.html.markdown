---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_inference_domains_v2"
sidebar_current: "docs-tencentcloud-data_source-teo_inference_domains_v2"
description: |-
  Use this data source to query the list of TEO inference service domains under a specified zone and inference service.
---

# tencentcloud_teo_inference_domains_v2

Use this data source to query the list of TEO inference service domains under a specified zone and inference service.

## Example Usage

```hcl
data "tencentcloud_teo_inference_domains_v2" "example" {
  zone_id    = "zone-2qtuhspy7cr6"
  service_id = "sid-123456"
}
```

### Sort by creation time in ascending order:

```hcl
data "tencentcloud_teo_inference_domains_v2" "example" {
  zone_id    = "zone-2qtuhspy7cr6"
  service_id = "sid-123456"
  sort_by    = "CreateTime"
  sort_order = "Asc"
}
```

## Argument Reference

The following arguments are supported:

* `service_id` - (Required, String) Inference service ID.
* `zone_id` - (Required, String) Site ID.
* `sort_by` - (Optional, String) Sort field. Valid values: `CreateTime` (domain creation time), `UpdateTime` (domain modification time). Default value: `CreateTime`.
* `sort_order` - (Optional, String) Sort order. Valid values: `Asc` (ascending), `Desc` (descending). Default value: `Desc`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

      * `alias` - Certificate remark name.
      * `cert_id` - Server certificate ID.
      * `content` - Content of the verification file.
      * `deploy_time` - Certificate deployment time.
      * `expire_time` - Certificate expiration time.
      * `name_servers` - DNS server addresses assigned to the user for NS access.
      * `path` - URL path portion of the file verification endpoint.
      * `record_type` - Record type.
      * `record_value` - Record value.
      * `sign_algo` - Signature algorithm.
      * `status` - Certificate status. Valid values: `deployed`, `processing`, `applying`, `failed`, `issued`.
      * `subdomain` - Host record.
      * `type` - Certificate type. Valid values: `default`, `upload`, `managed`.
    * `apply_type` - Application type. Valid values: `apply`, `none`.
    * `dns_verification` - DNS resolution verification information, used for CNAME and non-domain-name access.
    * `file_verification` - File verification information, used for CNAME and non-domain-name access.
    * `host_cert_info` - Server certificate configuration list.
    * `host` - Domain name.
    * `mode` - Certificate mode. Valid values: `disable`, `eofreecert`, `sslcert`.
    * `ns_verification` - NS verification information, used for NS access to switch DNS servers.
  * `auth_switch` - Inference task request authentication switch. Valid values: `Off`, `On`.
  * `certificate` - Domain certificate information.
  * `cname` - CNAME address.
  * `create_time` - Creation time.
  * `domain` - Domain name.
  * `ownership_verification` - Ownership verification information. Only present when the domain needs ownership verification to continue serving.
  * `status` - Domain status. Valid values: `Online`, `Process`, `Offline`, `Init`, `Forbidden`.
  * `update_time` - Modification time.
* `domains` - List of inference service domains.


