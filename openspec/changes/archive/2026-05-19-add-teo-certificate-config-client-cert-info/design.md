## Context

The `tencentcloud_teo_certificate_config` resource manages SSL/TLS certificate configuration for EdgeOne (TEO) acceleration domains. Currently, it supports server certificates (`server_cert_info`), upstream mutual TLS (`upstream_cert_info`), and the certificate mode (`mode`). However, the resource does not support configuring edge-side mutual TLS (client CA certificate authentication), which is available in the cloud API as the `ClientCertInfo` field.

The cloud API `ModifyHostsCertificate` accepts `ClientCertInfo` (type `MutualTLS`) in its request, and `DescribeAccelerationDomains` returns it within the `AccelerationDomainCertificate` structure. The `MutualTLS` type contains a `Switch` field and a `CertInfos` list of `CertificateInfo` objects.

The existing `upstream_cert_info.upstream_mutual_tls` parameter already uses the same `MutualTLS` pattern, so this design follows the established approach.

## Goals / Non-Goals

**Goals:**
- Add `client_cert_info` parameter to the `tencentcloud_teo_certificate_config` resource schema
- Support reading `ClientCertInfo` from `DescribeAccelerationDomains` response
- Support writing `ClientCertInfo` via `ModifyHostsCertificate` request
- Maintain backward compatibility (new parameter is Optional + Computed)

**Non-Goals:**
- Modifying existing `server_cert_info`, `upstream_cert_info`, or `mode` parameters
- Adding `UpstreamCertificateVerify` support (out of scope for this change)
- Adding `Status` field to existing `server_cert_info` or `upstream_cert_info.cert_infos` sub-schemas

## Decisions

### Decision 1: Schema structure for `client_cert_info`

**Choice**: Use `TypeList` with `MaxItems: 1` containing `switch` and `cert_infos` sub-fields, mirroring the existing `upstream_cert_info` pattern.

**Rationale**: The cloud API's `ClientCertInfo` is a `MutualTLS` struct, which is the same type used for `UpstreamCertInfo.UpstreamMutualTLS`. Following the existing pattern of `upstream_cert_info` (which wraps `upstream_mutual_tls`) ensures consistency and familiarity for users.

**Alternative considered**: Use `TypeSet` for `cert_infos` to avoid ordering issues — rejected because the existing `upstream_cert_info.upstream_mutual_tls.cert_infos` uses `TypeList`, and consistency is more important.

### Decision 2: `cert_infos` sub-field structure

**Choice**: Include `cert_id` (Required), `alias`, `type`, `expire_time`, `deploy_time`, `sign_algo`, `status` (all Computed) in each `cert_infos` element.

**Rationale**: The cloud API's `CertificateInfo` struct (used in `MutualTLS.CertInfos`) contains all these fields. When writing via `ModifyHostsCertificate`, only `cert_id` is needed as input. The rest are read-only metadata returned by `DescribeAccelerationDomains`. The `status` field is new compared to the existing `cert_infos` in `upstream_mutual_tls`, but it exists in the `CertificateInfo` struct and provides useful deployment status information.

### Decision 3: Read path for `ClientCertInfo`

**Choice**: Parse `ClientCertInfo` from `AccelerationDomainCertificate.ClientCertInfo` in the existing read post-handle function.

**Rationale**: The `DescribeAccelerationDomains` response already returns `ClientCertInfo` within `AccelerationDomainCertificate`. The current read function (`resourceTencentCloudTeoCertificateConfigReadPostHandleResponse0`) already processes `AccelerationDomain.Certificate`, so adding `ClientCertInfo` parsing there is the natural extension.

### Decision 4: Write path for `ClientCertInfo`

**Choice**: Set `request.ClientCertInfo` in the existing update-on-start function (`resourceTencentCloudTeoCertificateConfigUpdateOnStart`).

**Rationale**: The `ModifyHostsCertificate` request already supports `ClientCertInfo`. The current update function already constructs and sends the request, so adding `ClientCertInfo` there follows the established pattern.

## Risks / Trade-offs

- **[Backward compatibility]** → Mitigation: `client_cert_info` is Optional + Computed, existing configurations remain unaffected. Users who don't specify this field will see no change in behavior.
- **[API field naming consistency]** → The cloud API uses `ClientCertInfo` at the top level of both the request and response, but in the response it's nested inside `AccelerationDomainCertificate`. The terraform schema mirrors the API structure, which is consistent with how `upstream_cert_info` is handled.
