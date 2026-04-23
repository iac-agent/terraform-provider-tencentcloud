## ADDED Requirements

### Requirement: apply_type parameter for EO hosting type
The `tencentcloud_teo_certificate_config` resource SHALL support an `apply_type` parameter of type `schema.TypeString` with `Optional: true` and `Computed: true`. The parameter represents the EO hosting type, with valid values `none` (no EO hosting) and `apply` (host EO).

#### Scenario: Set apply_type to apply
- **WHEN** a user configures `apply_type = "apply"` in the resource
- **THEN** the `ModifyHostsCertificate` API SHALL be called with `ApplyType` set to `"apply"`

#### Scenario: Set apply_type to none
- **WHEN** a user configures `apply_type = "none"` in the resource
- **THEN** the `ModifyHostsCertificate` API SHALL be called with `ApplyType` set to `"none"`

#### Scenario: Read apply_type from API
- **WHEN** the resource read operation is performed
- **THEN** the system SHALL call `DescribeHostsSetting` API to retrieve the `ApplyType` value and set it in the Terraform state

#### Scenario: apply_type not configured
- **WHEN** the user does not specify `apply_type` in the configuration
- **THEN** the resource SHALL use the computed value from the API response, and `ModifyHostsCertificate` SHALL NOT include `ApplyType` in the request

### Requirement: client_cert_info parameter for client mutual TLS
The `tencentcloud_teo_certificate_config` resource SHALL support a `client_cert_info` parameter of type `schema.TypeList` with `MaxItems: 1`, `Optional: true`, and `Computed: true`. The parameter contains `switch` (TypeString, Required) and `cert_infos` (TypeList, Optional+Computed) sub-fields. The `cert_infos` list contains items with `cert_id` (TypeString, Required).

#### Scenario: Enable client mutual TLS with certificates
- **WHEN** a user configures `client_cert_info` with `switch = "on"` and `cert_infos` containing certificates
- **THEN** the `ModifyHostsCertificate` API SHALL be called with `ClientCertInfo` containing the configured switch and certificate IDs

#### Scenario: Disable client mutual TLS
- **WHEN** a user configures `client_cert_info` with `switch = "off"`
- **THEN** the `ModifyHostsCertificate` API SHALL be called with `ClientCertInfo.Switch` set to `"off"`

#### Scenario: Read client_cert_info from API
- **WHEN** the resource read operation is performed and the `DescribeAccelerationDomains` response contains `Certificate.ClientCertInfo`
- **THEN** the system SHALL populate `client_cert_info` in the Terraform state from the `AccelerationDomainCertificate.ClientCertInfo` field

#### Scenario: client_cert_info not configured
- **WHEN** the user does not specify `client_cert_info` in the configuration
- **THEN** the resource SHALL use the computed value from the API response, and `ModifyHostsCertificate` SHALL NOT include `ClientCertInfo` in the request

### Requirement: Update method sends new parameters to API
The resource update method SHALL include `apply_type` and `client_cert_info` in the `ModifyHostsCertificate` API request when these parameters are configured in the Terraform resource.

#### Scenario: Update apply_type
- **WHEN** a user changes the `apply_type` value in an existing resource
- **THEN** the `ModifyHostsCertificate` API SHALL be called with the updated `ApplyType` value

#### Scenario: Update client_cert_info
- **WHEN** a user changes the `client_cert_info` value in an existing resource
- **THEN** the `ModifyHostsCertificate` API SHALL be called with the updated `ClientCertInfo` value
