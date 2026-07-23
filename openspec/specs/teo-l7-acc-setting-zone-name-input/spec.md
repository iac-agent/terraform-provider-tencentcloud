## Requirements

### Requirement: Zone name is optionally inputable
The `zone_name` field in the `tencentcloud_teo_l7_acc_setting` resource SHALL be `Optional: true, Computed: true`, allowing users to optionally specify the zone name in their Terraform configuration while the authoritative value continues to be read from the cloud API.

#### Scenario: User specifies zone_name in configuration
- **WHEN** a user specifies `zone_name` in their `tencentcloud_teo_l7_acc_setting` resource configuration
- **THEN** the Terraform provider SHALL accept the value without error
- **AND** the provider SHALL read the authoritative `zone_name` from the cloud API response during refresh

#### Scenario: User does not specify zone_name
- **WHEN** a user does not specify `zone_name` in their `tencentcloud_teo_l7_acc_setting` resource configuration
- **THEN** the Terraform provider SHALL populate `zone_name` from the cloud API response as before
- **AND** existing configurations SHALL continue to work unchanged

#### Scenario: User specifies a zone_name that differs from API value
- **WHEN** a user specifies a `zone_name` that differs from the actual zone name returned by the API
- **THEN** the Terraform provider SHALL use the API-returned value as the authoritative `zone_name` during state refresh