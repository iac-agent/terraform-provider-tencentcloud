## ADDED Requirements

### Requirement: client_cert_info parameter for edge mutual TLS configuration
The `tencentcloud_teo_certificate_config` resource SHALL support a `client_cert_info` parameter (TypeList, MaxItems: 1, Optional, Computed) that configures edge-side mutual TLS (client CA certificate authentication) for the acceleration domain. The `client_cert_info` parameter SHALL contain a `switch` sub-field (Required, TypeString, values: "on"/"off") and a `cert_infos` sub-field (Optional, Computed, TypeList) representing the client CA certificate list.

#### Scenario: Create resource with client_cert_info enabled
- **WHEN** a user creates a `tencentcloud_teo_certificate_config` resource with `client_cert_info` containing `switch = "on"` and `cert_infos` with a `cert_id`
- **THEN** the resource SHALL send the `ClientCertInfo` field in the `ModifyHostsCertificate` API request with the specified switch and certificate IDs

#### Scenario: Create resource without client_cert_info
- **WHEN** a user creates a `tencentcloud_teo_certificate_config` resource without specifying `client_cert_info`
- **THEN** the resource SHALL NOT send the `ClientCertInfo` field in the `ModifyHostsCertificate` API request, preserving backward compatibility

#### Scenario: Read client_cert_info from API response
- **WHEN** the resource read function calls `DescribeAccelerationDomains` and the response contains `ClientCertInfo` in `AccelerationDomainCertificate`
- **THEN** the resource SHALL parse and set the `client_cert_info` field in the Terraform state, including `switch` and all `cert_infos` sub-fields (`cert_id`, `alias`, `type`, `expire_time`, `deploy_time`, `sign_algo`, `status`)

#### Scenario: Read with nil ClientCertInfo in API response
- **WHEN** the resource read function calls `DescribeAccelerationDomains` and the response has `ClientCertInfo` as nil
- **THEN** the resource SHALL NOT set the `client_cert_info` field in the Terraform state

### Requirement: cert_infos sub-field structure in client_cert_info
Each element in the `client_cert_info.cert_infos` list SHALL contain `cert_id` (Required, TypeString) as the certificate ID, and `alias` (Computed, TypeString), `type` (Computed, TypeString), `expire_time` (Computed, TypeString), `deploy_time` (Computed, TypeString), `sign_algo` (Computed, TypeString), `status` (Computed, TypeString) as read-only computed fields.

#### Scenario: Write cert_infos with only cert_id
- **WHEN** a user specifies `client_cert_info.cert_infos` with only `cert_id`
- **THEN** the resource SHALL construct the `MutualTLS.CertInfos` in the `ModifyHostsCertificate` request with only the `CertId` field populated

#### Scenario: Read cert_infos with all computed fields
- **WHEN** the API response returns `ClientCertInfo.CertInfos` with all fields populated
- **THEN** the resource SHALL set all computed fields (`alias`, `type`, `expire_time`, `deploy_time`, `sign_algo`, `status`) in the Terraform state for each certificate entry
