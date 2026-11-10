---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_inference_domains_1110_01"
sidebar_current: "docs-tencentcloud-data_source-teo_inference_domains_1110_01"
description: |-
  Use this data source to query the list of TEO (EdgeOne) inference service domains.
---

# tencentcloud_teo_inference_domains_1110_01

Use this data source to query the list of TEO (EdgeOne) inference service domains.

## Example Usage

```hcl
data "tencentcloud_teo_inference_domains_1110_01" "example" {
  zone_id    = "zone-2qtuhspy7cr6"
  service_id = "sid-2rugn89bkla9"
  sort_by    = "CreateTime"
  sort_order = "Desc"
}
```

## Argument Reference

The following arguments are supported:

* `service_id` - (Required, String) Inference service ID.
* `zone_id` - (Required, String) Site ID.
* `sort_by` - (Optional, String) Sort field. Valid values: CreateTime (domain creation time), UpdateTime (domain modification time). Default value: CreateTime.
* `sort_order` - (Optional, String) Sort order. Valid values: Asc (ascending), Desc (descending). Default value: Desc.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

      * `alias` - Certificate alias. Note: this field may return null, indicating that no valid value can be obtained.
      * `cert_id` - Server certificate ID. Note: this field may return null, indicating that no valid value can be obtained.
      * `content` - Content of the verification file, which needs to be written into the txt file returned by the path field.
      * `deploy_time` - Certificate deployment time. Note: this field may return null, indicating that no valid value can be obtained.
      * `expire_time` - Certificate expiration time. Note: this field may return null, indicating that no valid value can be obtained.
      * `name_servers` - DNS server addresses assigned to the user during NS access; the domain NameServer needs to be switched to these addresses.
      * `path` - The URL path part of the verification file. EdgeOne fetches it via http://{Host}{URL Path}.
      * `record_type` - Record type.
      * `record_value` - Record value.
      * `sign_algo` - Signature algorithm. Note: this field may return null, indicating that no valid value can be obtained.
      * `status` - Certificate status. Valid values: deployed (deployed), processing (deploying), applying (applying), failed (application failed), issued (binding failed). Note: this field may return null, indicating that no valid value can be obtained.
      * `subdomain` - Host record.
      * `type` - Certificate type. Valid values: default (default certificate), upload (user uploaded), managed (Tencent Cloud managed). Note: this field may return null, indicating that no valid value can be obtained.
    * `apply_type` - Application type. Valid values: apply (managed by EdgeOne), none (not managed by EdgeOne). Note: this field may return null, indicating that no valid value can be obtained.
    * `dns_verification` - Information required when using DNS resolution verification for CNAME or no-domain access. Note: this field may return null, indicating that no valid value can be obtained.
    * `file_verification` - Information required when using file verification for CNAME or no-domain access. Note: this field may return null, indicating that no valid value can be obtained.
    * `host_cert_info` - Server certificate configuration. Note: this field may return null, indicating that no valid value can be obtained.
    * `host` - Domain. Note: this field may return null, indicating that no valid value can be obtained.
    * `mode` - Certificate configuration mode. Valid values: disable (no certificate), eofreecert (EdgeOne free certificate), sslcert (SSL certificate).
    * `ns_verification` - Information required for switching DNS servers during NS access. Note: this field may return null, indicating that no valid value can be obtained.
  * `auth_switch` - Authentication switch of inference task requests. Valid values: Off (disabled), On (enabled). Default value: On.
  * `certificate` - Domain certificate information. Refer to ApplyFreeCertificate, CheckFreeCertificateVerification and ModifyHostsCertificate APIs.
  * `cname` - CNAME address. To verify the CNAME configuration status of the domain, refer to the CheckCnameStatus API.
  * `create_time` - Creation time.
  * `domain` - Domain name.
  * `ownership_verification` - Information required for ownership verification when the domain needs ownership verification to continue serving. Refer to the VerifyOwnership API for ownership verification.
  * `status` - Domain status. Valid values: Online (effective), Process (deploying), Offline (deactivated), Init (not effective, waiting for site activation), Forbidden (domain disabled).
  * `update_time` - Modification time.
* `domains` - List of inference service domains.


