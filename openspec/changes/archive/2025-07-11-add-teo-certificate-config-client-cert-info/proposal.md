## Why

The `tencentcloud_teo_certificate_config` resource currently lacks support for the `client_cert_info` parameter in the `ModifyHostsCertificate` API. This parameter configures the client-side CA certificate for edge mutual TLS (mTLS) authentication in EdgeOne (TEO). Without this parameter, users cannot manage client certificate authentication for their TEO domains through Terraform, which is a critical security feature for edge-side mutual TLS scenarios.

## What Changes

- Add a new `client_cert_info` parameter (TypeList, MaxItems: 1) to the `tencentcloud_teo_certificate_config` resource schema, corresponding to the `ClientCertInfo` field in the `ModifyHostsCertificate` API request.
- The `client_cert_info` parameter maps to the `MutualTLS` cloud API struct, which contains:
  - `switch` (Required, string): Mutual TLS configuration switch, values: `on` / `off`.
  - `cert_infos` (Optional, TypeList): Mutual TLS certificate list, containing `cert_id` (Required, string) sub-fields.
- Update the resource's update (create) logic in `resource_tc_teo_certificate_config_extension.go` to populate `request.ClientCertInfo` when the parameter is set.
- Update the resource's read logic to handle the `ClientCertInfo` response from the `DescribeAccelerationDomains` API.
- Add unit tests for the new parameter handling.

## Capabilities

### New Capabilities
- `teo-certificate-config-client-cert-info`: Adds `client_cert_info` parameter to the `tencentcloud_teo_certificate_config` resource, enabling edge-side mutual TLS client certificate configuration for TEO domains.

### Modified Capabilities

## Impact

- **Affected files**:
  - `tencentcloud/services/teo/resource_tc_teo_certificate_config.go` — Schema definition
  - `tencentcloud/services/teo/resource_tc_teo_certificate_config_extension.go` — Read/Update logic
  - `tencentcloud/services/teo/resource_tc_teo_certificate_config_test.go` — Unit tests
  - `tencentcloud/services/teo/resource_tc_teo_certificate_config.md` — Documentation
- **Cloud API**: `ModifyHostsCertificate` (teo v20220901) — `ClientCertInfo` field (type: `MutualTLS`)
- **Backward compatibility**: Fully backward compatible — the new parameter is Optional/Computed, existing configurations remain unchanged.
