## Context

The `tencentcloud_teo_certificate_config` resource manages SSL/TLS certificate configurations for EdgeOne (TEO) acceleration domains. Currently, the resource supports configuring server certificates (`server_cert_info`), mode selection (`mode`), and upstream certificate info (`upstream_cert_info` with `upstream_mutual_tls`). However, it lacks support for the `client_cert_info` parameter, which configures client-side CA certificates for edge mutual TLS (mTLS) authentication.

The `ClientCertInfo` field is part of the `ModifyHostsCertificate` API request (package: `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901`), using the `MutualTLS` struct type. This is the same `MutualTLS` struct used by the existing `upstream_cert_info.upstream_mutual_tls` parameter, but in a different context — here it represents the client CA certificate deployed at the EO edge node for client authentication, rather than the certificate carried during origin-pull.

The current resource already uses `ModifyHostsCertificate` in its update/create flow (`resourceTencentCloudTeoCertificateConfigUpdateOnStart`), so adding `ClientCertInfo` to the request is straightforward.

## Goals / Non-Goals

**Goals:**
- Add the `client_cert_info` parameter to the `tencentcloud_teo_certificate_config` resource schema
- Support reading `ClientCertInfo` from the API response in the read flow
- Support writing `ClientCertInfo` in the update/create flow via `ModifyHostsCertificate`
- Maintain full backward compatibility — the new parameter is Optional/Computed

**Non-Goals:**
- Adding `apply_type` parameter (deprecated in the cloud API)
- Modifying any other existing parameters or their behavior
- Adding `upstream_certificate_verify` sub-field to `upstream_cert_info` (out of scope for this change)
- Changing the resource ID format or any ForceNew behavior

## Decisions

### 1. Schema design for `client_cert_info`

**Decision**: Use `TypeList` with `MaxItems: 1` containing a nested `schema.Resource` with `switch` (Required, string) and `cert_infos` (Optional, TypeList of cert_id strings).

**Rationale**: The cloud API `ClientCertInfo` is of type `*MutualTLS`, which contains `Switch` (string) and `CertInfos` ([]*CertificateInfo). This follows the same pattern as the existing `upstream_cert_info.upstream_mutual_tls` parameter in the same resource. Using `TypeList` with `MaxItems: 1` is the standard pattern for nested object parameters in this codebase.

### 2. Reading `ClientCertInfo` from API response

**Decision**: Read `ClientCertInfo` from the `AccelerationDomain.Certificate.ClientCertInfo` field in the read post-handle function.

**Rationale**: The `DescribeAccelerationDomains` API returns certificate configuration under `AccelerationDomain.Certificate`. The `ClientCertInfo` field follows the same `MutualTLS` structure. This is consistent with how the existing `upstream_cert_info` is read from `AccelerationDomain.Certificate.UpstreamCertInfo`.

### 3. Parameter optionality

**Decision**: Make `client_cert_info` Optional and Computed.

**Rationale**: The `ClientCertInfo` field in the API is optional (marked `omitnil`). When not specified, the API keeps the original configuration. Setting it as Computed ensures the read flow can populate it from the API response. This maintains backward compatibility — existing configurations without `client_cert_info` will continue to work.

## Risks / Trade-offs

- **[Risk] API response structure uncertainty**: The `ClientCertInfo` field in the read response may not be nested under `Certificate` the same way as `UpstreamCertInfo`. → **Mitigation**: Verify the API response structure by examining the `Certificate` struct in the SDK models. The `ClientCertInfo` field is already present in the `Certificate` struct as `*MutualTLS`, following the same pattern.

- **[Risk] Backward compatibility**: Adding a new parameter could potentially affect existing state files. → **Mitigation**: The parameter is Optional/Computed, so existing configurations without it will continue to work. Terraform will automatically populate it from the API response on the next read.
