## Why

The `tencentcloud_teo_certificate_config` resource currently lacks support for configuring edge-side mutual TLS (client CA certificate authentication). The cloud API's `ModifyHostsCertificate` interface already supports the `ClientCertInfo` field (type `MutualTLS`), and the `DescribeAccelerationDomains` response returns this field in `AccelerationDomainCertificate`. Without this parameter, users cannot manage edge mutual TLS authentication through Terraform, forcing them to use the console or API directly.

## What Changes

- Add a new `client_cert_info` parameter (TypeList, MaxItems: 1) to the `tencentcloud_teo_certificate_config` resource schema, mapping to the cloud API's `ClientCertInfo` field
- The `client_cert_info` parameter contains:
  - `switch` (Required, string): Mutual TLS switch, values: `on` / `off`
  - `cert_infos` (Optional, Computed, TypeList): Client CA certificate list
    - `cert_id` (Required, string): Certificate ID from SSL
    - `alias` (Computed, string): Certificate alias
    - `type` (Computed, string): Certificate type (default/upload/managed)
    - `expire_time` (Computed, string): Certificate expiration time
    - `deploy_time` (Computed, string): Certificate deployment time
    - `sign_algo` (Computed, string): Signature algorithm
    - `status` (Computed, string): Certificate status (deployed/processing/applying/failed/issued)
- Update the resource's update function to send `ClientCertInfo` in `ModifyHostsCertificate` request
- Update the resource's read function to parse `ClientCertInfo` from `DescribeAccelerationDomains` response
- Update the resource's `.md` documentation file with `client_cert_info` example

## Capabilities

### New Capabilities
- `teo-certificate-config-client-cert-info`: Add `client_cert_info` parameter to the `tencentcloud_teo_certificate_config` resource for edge-side mutual TLS (client CA certificate) configuration

### Modified Capabilities
<!-- No existing capability requirements are changing -->

## Impact

- **Code**: `tencentcloud/services/teo/resource_tc_teo_certificate_config.go` (schema definition), `resource_tc_teo_certificate_config_extension.go` (read/update logic)
- **Tests**: `tencentcloud/services/teo/resource_tc_teo_certificate_config_test.go`
- **Docs**: `tencentcloud/services/teo/resource_tc_teo_certificate_config.md`
- **APIs**: `ModifyHostsCertificate` (write), `DescribeAccelerationDomains` (read)
- **Backward Compatibility**: New parameter is Optional + Computed, no breaking changes to existing configurations
