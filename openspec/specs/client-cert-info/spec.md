# client-cert-info Specification

## Purpose
TBD - created by archiving change add-teo-certificate-config-client-cert-info. Update Purpose after archive.
## Requirements
### Requirement: client_cert_info parameter in resource schema
The `tencentcloud_teo_certificate_config` resource SHALL include a `client_cert_info` parameter of type TypeList with MaxItems: 1, Optional and Computed. The parameter SHALL contain the following sub-fields:
- `switch` (TypeString, Required): Mutual authentication configuration switch, values SHALL be `on` or `off`.
- `cert_infos` (TypeList, Optional and Computed): Client certificate list, with sub-fields:
  - `cert_id` (TypeString, Required): Certificate ID from SSL service.
  - `alias` (TypeString, Computed): Certificate alias.
  - `type` (TypeString, Computed): Certificate type (default/upload/managed).
  - `expire_time` (TypeString, Computed): Certificate expiration time.
  - `deploy_time` (TypeString, Computed): Certificate deployment time.
  - `sign_algo` (TypeString, Computed): Signature algorithm.
  - `status` (TypeString, Computed): Certificate status (deployed/processing/applying/failed/issued).

#### Scenario: User configures client certificate authentication
- **WHEN** user sets `client_cert_info` with `switch = "on"` and provides `cert_infos` with `cert_id`
- **THEN** the resource SHALL call `ModifyHostsCertificate` with `ClientCertInfo` containing the specified switch and certificate IDs

#### Scenario: User disables client certificate authentication
- **WHEN** user sets `client_cert_info` with `switch = "off"`
- **THEN** the resource SHALL call `ModifyHostsCertificate` with `ClientCertInfo.Switch = "off"`

#### Scenario: User does not specify client_cert_info
- **WHEN** user does not set `client_cert_info` in the terraform configuration
- **THEN** the resource SHALL NOT include `ClientCertInfo` in the `ModifyHostsCertificate` request, preserving existing configuration

### Requirement: client_cert_info read from cloud API
The resource read function SHALL populate `client_cert_info` from `AccelerationDomainCertificate.ClientCertInfo` in the `DescribeAccelerationDomains` response.

#### Scenario: ClientCertInfo is present in API response
- **WHEN** the `DescribeAccelerationDomains` response contains `ClientCertInfo` with `Switch = "on"` and `CertInfos`
- **THEN** the resource SHALL set `client_cert_info` in the terraform state with the switch value and certificate information

#### Scenario: ClientCertInfo is nil in API response
- **WHEN** the `DescribeAccelerationDomains` response has `ClientCertInfo` as nil
- **THEN** the resource SHALL NOT set `client_cert_info` in the terraform state

#### Scenario: CertInfos fields are nil
- **WHEN** individual fields within `CertInfos` entries are nil
- **THEN** the resource SHALL skip setting those nil fields (consistent with existing pattern for server_cert_info and upstream_cert_info)

### Requirement: Backward compatibility
Adding the `client_cert_info` parameter SHALL NOT break existing terraform configurations or state files that do not use this parameter.

#### Scenario: Existing configuration without client_cert_info
- **WHEN** an existing terraform configuration for `tencentcloud_teo_certificate_config` does not include `client_cert_info`
- **THEN** the resource SHALL continue to function correctly with no changes to behavior
