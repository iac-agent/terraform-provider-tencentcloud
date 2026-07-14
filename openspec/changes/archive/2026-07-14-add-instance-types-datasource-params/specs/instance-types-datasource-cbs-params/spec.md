## ADDED Requirements

### Requirement: InquiryType Parameter for CBS Query
The `tencentcloud_instance_types` data source SHALL support an `InquiryType` parameter within the `cbs_filter` block to control the query mode of the `DescribeDiskConfigQuota` CBS API.

#### Scenario: User specifies InquiryType as INQUIRY_CBS_CONFIG
- **WHEN** user provides `cbs_filter` with `inquiry_type` set to `"INQUIRY_CBS_CONFIG"`
- **THEN** the `DescribeDiskConfigQuota` API call SHALL use `"INQUIRY_CBS_CONFIG"` as the `InquiryType` request parameter
- **AND** the CBS query SHALL return CBS-only disk configurations without CVM instance coupling requirements

#### Scenario: User does not specify InquiryType
- **WHEN** user provides `cbs_filter` without `inquiry_type`
- **THEN** the `DescribeDiskConfigQuota` API call SHALL default to `"INQUIRY_CVM_CONFIG"` as the `InquiryType` request parameter
- **AND** the behavior SHALL be identical to the current hardcoded behavior

#### Scenario: InquiryType field schema definition
- **WHEN** the data source schema is defined
- **THEN** `inquiry_type` SHALL be an optional string field within the `cbs_filter` nested schema block
- **AND** `inquiry_type` SHALL have description indicating valid values are `"INQUIRY_CBS_CONFIG"` and `"INQUIRY_CVM_CONFIG"`
- **AND** `inquiry_type` SHALL NOT be a required field (backward compatibility)

### Requirement: InstanceFamilies Parameter for CBS Query
The `tencentcloud_instance_types` data source SHALL support an `InstanceFamilies` parameter within the `cbs_filter` block to allow users to specify instance families for the `DescribeDiskConfigQuota` CBS API independently of instance type results.

#### Scenario: User specifies InstanceFamilies explicitly
- **WHEN** user provides `cbs_filter` with `instance_families` set to a list of instance family names (e.g., `["S5", "M5"]`)
- **THEN** the `DescribeDiskConfigQuota` API call SHALL use the user-provided instance families list as the `InstanceFamilies` request parameter
- **AND** the CBS query SHALL filter results based on the specified instance families regardless of what instance type results contain

#### Scenario: User does not specify InstanceFamilies
- **WHEN** user provides `cbs_filter` without `instance_families`
- **THEN** the `DescribeDiskConfigQuota` API call SHALL derive `InstanceFamilies` from the instance type results' `family` field for each instance type
- **AND** the behavior SHALL be identical to the current derivation logic

#### Scenario: InstanceFamilies field schema definition
- **WHEN** the data source schema is defined
- **THEN** `instance_families` SHALL be an optional list of strings field within the `cbs_filter` nested schema block
- **AND** `instance_families` SHALL have description indicating it specifies instance family names for CBS configuration filtering
- **AND** `instance_families` SHALL NOT be a required field (backward compatibility)

### Requirement: DiskChargeType Already Exists in cbs_filter
The `DiskChargeType` parameter already exists as `disk_charge_type` within the `cbs_filter` block of the `tencentcloud_instance_types` data source and maps to the `DescribeDiskConfigQuota` API's `DiskChargeType` request parameter. No new schema addition is needed for this parameter.

#### Scenario: DiskChargeType already functional
- **WHEN** user provides `cbs_filter` with `disk_charge_type` set to `"PREPAID"` or `"POSTPAID_BY_HOUR"`
- **THEN** the `DescribeDiskConfigQuota` API call SHALL use the provided value as the `DiskChargeType` request parameter
- **AND** this behavior SHALL be identical to the current existing behavior

#### Scenario: DiskChargeType field is unchanged
- **WHEN** the data source schema is reviewed
- **THEN** `disk_charge_type` SHALL remain as an optional string field within the `cbs_filter` nested schema block
- **AND** no schema modifications SHALL be made to the `disk_charge_type` field

### Requirement: Backward Compatibility
The new parameters MUST maintain backward compatibility with existing Terraform configurations using the `tencentcloud_instance_types` data source.

#### Scenario: Existing cbs_filter configuration without new parameters
- **WHEN** user has existing Terraform configuration using `cbs_filter` with only `disk_types`, `disk_charge_type`, and `disk_usage`
- **THEN** configuration SHALL continue to work without modification
- **AND** `InquiryType` SHALL default to `"INQUIRY_CVM_CONFIG"`
- **AND** `InstanceFamilies` SHALL be derived from instance type results
- **AND** terraform plan SHALL not show any changes to the existing configuration

#### Scenario: Data source without cbs_filter block
- **WHEN** user has existing Terraform configuration without the `cbs_filter` block
- **THEN** configuration SHALL continue to work without modification
- **AND** only the `DescribeZoneInstanceConfigInfos` CVM API SHALL be called
- **AND** no CBS API call SHALL be made
