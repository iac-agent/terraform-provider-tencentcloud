## Context

The `tencentcloud_teo_certificate_config` resource manages certificate configurations for EdgeOne (TEO) acceleration domains. It currently supports `zone_id`, `host`, `mode`, `server_cert_info`, and `upstream_cert_info` parameters. The underlying `ModifyHostsCertificate` API also supports `ApplyType` and `ClientCertInfo` parameters, which are not yet exposed in the Terraform resource.

Current state:
- **Write API**: `ModifyHostsCertificate` supports `ZoneId`, `Hosts`, `Mode`, `ServerCertInfo`, `ApplyType`, `ClientCertInfo`, `UpstreamCertInfo`
- **Read API**: `DescribeAccelerationDomains` returns `AccelerationDomainCertificate` with `Mode`, `List`, `ClientCertInfo`, `UpstreamCertInfo` — but does NOT return `ApplyType`
- The `ApplyType` field IS available in `DescribeHostsSetting` API response (via `DetailHost.Https.ApplyType`)

## Goals / Non-Goals

**Goals:**
- Add `apply_type` parameter (Optional + Computed, TypeString) to support EO hosting type configuration
- Add `client_cert_info` parameter (Optional + Computed, TypeList/MaxItems:1) to support client-side mutual TLS certificate configuration
- Ensure `apply_type` can be read back by calling the `DescribeHostsSetting` API when `ApplyType` is not available from the primary read API
- Maintain backward compatibility — existing configurations must continue to work

**Non-Goals:**
- Changing the existing `host` parameter to `hosts` (list type) — this would be a breaking change
- Modifying the read service method signature
- Adding `UpstreamCertificateVerify` sub-fields to `upstream_cert_info` (out of scope for this change)

## Decisions

### Decision 1: How to read `apply_type`

The primary read API (`DescribeAccelerationDomains`) returns `AccelerationDomainCertificate` which does not contain `ApplyType`. The `DescribeHostsSetting` API returns `DetailHost.Https.ApplyType`.

**Choice**: Call `DescribeHostsSetting` in the read post-handler to retrieve `ApplyType`. This adds an extra API call but ensures the field is accurately read.

**Alternative considered**: Mark `apply_type` as write-only (no read-back). Rejected because this is inconsistent with Terraform best practices and would cause perpetual diffs.

### Decision 2: Schema design for `client_cert_info`

The `ClientCertInfo` in the cloud API is of type `*MutualTLS`, which contains `Switch` (string) and `CertInfos` ([]*CertificateInfo).

**Choice**: Use `TypeList` with `MaxItems: 1` containing `switch` (TypeString, Required) and `cert_infos` (TypeList, Optional+Computed) sub-fields. The `cert_infos` list contains items with `cert_id` (TypeString, Required). This mirrors the existing `upstream_cert_info.upstream_mutual_tls` pattern already in the resource.

### Decision 3: Parameter mutability

Both new parameters are Optional + Computed and can be updated without recreating the resource (not ForceNew). This is consistent with the existing `mode`, `server_cert_info`, and `upstream_cert_info` parameters.

## Risks / Trade-offs

- **Extra API call for read**: Adding a `DescribeHostsSetting` call increases read latency. Mitigation: The call is only made when the certificate config exists, and it provides accurate state.
- **`ApplyType` deprecated in some API versions**: The `ModifyHostsCertificateRequest` has a deprecation note on `ApplyType`. However, it is still a functional parameter. If it is fully removed in the future, we may need to remove the Terraform parameter as well.
- **Backward compatibility**: Both new parameters are Optional + Computed, so existing configurations are not affected.
