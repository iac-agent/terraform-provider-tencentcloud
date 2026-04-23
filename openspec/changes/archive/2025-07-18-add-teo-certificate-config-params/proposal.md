## Why

The `tencentcloud_teo_certificate_config` resource currently lacks support for the `apply_type` and `client_cert_info` parameters of the `ModifyHostsCertificate` API. Users cannot configure the EO hosting type (`apply_type`) or the client-side mutual TLS certificate (`client_cert_info`) through Terraform, which limits the ability to fully manage certificate configurations for EdgeOne (TEO) domains.

## What Changes

- Add `apply_type` parameter (Optional + Computed, TypeString) to the `tencentcloud_teo_certificate_config` resource schema, supporting values `none` (no EO hosting) and `apply` (host EO).
- Add `client_cert_info` parameter (Optional + Computed, TypeList, MaxItems: 1) to the `tencentcloud_teo_certificate_config` resource schema, containing `switch` (on/off) and `cert_infos` (list of certificate info with `cert_id`).
- Update the `ModifyHostsCertificate` request construction in the update method to include the new `apply_type` and `client_cert_info` parameters.
- Update the read method to populate `apply_type` and `client_cert_info` from the API response.
- Update the resource documentation (.md file).

## Capabilities

### New Capabilities
- `teo-certificate-config-params`: Adds `apply_type` and `client_cert_info` parameters to the `tencentcloud_teo_certificate_config` resource, enabling full configuration of EO hosting type and client-side mutual TLS certificate.

### Modified Capabilities

## Impact

- Affected files: `tencentcloud/services/teo/resource_tc_teo_certificate_config.go`, `tencentcloud/services/teo/resource_tc_teo_certificate_config_extension.go`, `tencentcloud/services/teo/resource_tc_teo_certificate_config_test.go`, resource documentation `.md` file.
- API: Uses existing `ModifyHostsCertificate` API with new input parameters `ApplyType` and `ClientCertInfo`.
- No breaking changes: both new parameters are Optional + Computed, existing configurations remain compatible.
