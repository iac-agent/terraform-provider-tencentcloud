## Context

The `tencentcloud_teo_certificate_config` resource manages SSL/TLS certificate configuration for TEO (TencentCloud EdgeOne) acceleration domains. Currently, the resource supports server certificates (`server_cert_info`) and upstream mutual TLS (`upstream_cert_info.upstream_mutual_tls`), but does not support the `ClientCertInfo` parameter for edge-side client certificate authentication.

The cloud API already supports `ClientCertInfo` in:
- `ModifyHostsCertificate` request (write path)
- `AccelerationDomainCertificate` response via `DescribeAccelerationDomains` (read path)

The `ClientCertInfo` field is of type `MutualTLS`, the same type used by `upstream_cert_info.upstream_mutual_tls`, providing a proven pattern to follow.

Current resource files:
- Schema: `tencentcloud/services/teo/resource_tc_teo_certificate_config.go`
- Extension (CRUD logic): `tencentcloud/services/teo/resource_tc_teo_certificate_config_extension.go`
- Tests: `tencentcloud/services/teo/resource_tc_teo_certificate_config_test.go`

## Goals / Non-Goals

**Goals:**
- Add `client_cert_info` parameter to enable edge mutual TLS client certificate configuration
- Support both read and write paths for `ClientCertInfo`
- Maintain backward compatibility (new parameter is Optional + Computed)
- Follow existing patterns from `upstream_cert_info.upstream_mutual_tls` for consistency

**Non-Goals:**
- Modifying existing parameters (`server_cert_info`, `upstream_cert_info`, `mode`)
- Adding `ApplyType` parameter (separate concern, not part of this change)
- Adding `UpstreamCertificateVerify` sub-field within `upstream_cert_info` (separate concern)

## Decisions

1. **Schema design**: `client_cert_info` will be a TypeList with MaxItems: 1, Optional + Computed, mirroring the `upstream_cert_info` pattern. This allows the cloud API to return the field even when not explicitly set by the user.

2. **Sub-field `cert_infos`**: Within `client_cert_info`, the `cert_infos` sub-field will contain `cert_id` (Required for input) and computed fields (`alias`, `type`, `expire_time`, `deploy_time`, `sign_algo`, `status`). When used as an input parameter in `ModifyHostsCertificate`, only `cert_id` needs to be provided per the API documentation.

3. **Read path**: `ClientCertInfo` is read from `AccelerationDomainCertificate.ClientCertInfo` in the `DescribeAccelerationDomains` response, which is already used by the existing read function.

4. **Write path**: `ClientCertInfo` is set in the `ModifyHostsCertificate` request within `resourceTencentCloudTeoCertificateConfigUpdateOnStart`.

5. **No new extension file**: All changes go into the existing `_extension.go` file, following the current code organization.

## Risks / Trade-offs

- **[Backward compatibility]** → Adding a new Optional+Computed field ensures existing TF configurations remain valid without changes.
- **[API consistency]** → The `MutualTLS` type is already used for `upstream_cert_info.upstream_mutual_tls`, so the implementation pattern is well-established and tested.
- **[Read consistency]** → The `CertificateInfo` struct in `ClientCertInfo.CertInfos` includes a `Status` field not present in the `ServerCertInfo` mapping. This will be added as a computed-only field in `cert_infos`.
