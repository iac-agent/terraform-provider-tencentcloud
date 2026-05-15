## Why

The `tencentcloud_teo_certificate_config` resource currently lacks support for configuring client certificate authentication (edge mutual TLS). The cloud API `ModifyHostsCertificate` already supports a `ClientCertInfo` parameter of type `MutualTLS`, and the `DescribeAccelerationDomains` response includes `ClientCertInfo` in `AccelerationDomainCertificate`. Adding this parameter enables Terraform users to manage edge-side client certificate authentication for TEO acceleration domains.

## What Changes

- Add a new `client_cert_info` parameter (TypeList, MaxItems: 1, Optional+Computed) to the `tencentcloud_teo_certificate_config` resource schema, mapping to the cloud API's `ClientCertInfo` field (MutualTLS type).
  - Sub-field `switch` (TypeString, Required): Mutual authentication switch, values: `on` / `off`.
  - Sub-field `cert_infos` (TypeList, Optional+Computed): Client certificate list, with sub-fields: `cert_id` (Required), `alias` (Computed), `type` (Computed), `expire_time` (Computed), `deploy_time` (Computed), `sign_algo` (Computed), `status` (Computed).
- Update the Create/Update logic (`resourceTencentCloudTeoCertificateConfigUpdateOnStart`) to set `ClientCertInfo` in the `ModifyHostsCertificate` request.
- Update the Read logic (`resourceTencentCloudTeoCertificateConfigReadPostHandleResponse0`) to read `ClientCertInfo` from `AccelerationDomainCertificate` and populate the terraform state.
- Update the resource `.md` documentation example to include `client_cert_info`.

## Capabilities

### New Capabilities
- `client-cert-info`: Add `client_cert_info` parameter to `tencentcloud_teo_certificate_config` resource for edge mutual TLS client certificate configuration.

### Modified Capabilities

## Impact

- **Resource file**: `tencentcloud/services/teo/resource_tc_teo_certificate_config.go` (schema definition)
- **Extension file**: `tencentcloud/services/teo/resource_tc_teo_certificate_config_extension.go` (read/write logic)
- **Test file**: `tencentcloud/services/teo/resource_tc_teo_certificate_config_test.go` (unit tests)
- **Documentation**: `tencentcloud/services/teo/resource_tc_teo_certificate_config.md` (example usage)
- **Cloud APIs**: `ModifyHostsCertificate` (write), `DescribeAccelerationDomains` (read), `DescribeZones` (auxiliary read for zone info)
- **SDK package**: `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901`
