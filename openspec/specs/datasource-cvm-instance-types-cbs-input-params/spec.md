## ADDED Requirements

### Requirement: Inquiry Type Input Parameter
The `tencentcloud_instance_types` data source SHALL support an `inquiry_type` input parameter within the `cbs_filter` block to specify the query category for the `DescribeDiskConfigQuota` API call.

#### Scenario: User specifies inquiry_type as INQUIRY_CBS_CONFIG
- **WHEN** user provides `inquiry_type` with value `"INQUIRY_CBS_CONFIG"` in the `cbs_filter` block
- **THEN** the `DescribeDiskConfigQuota` API call SHALL set `InquiryType` to `"INQUIRY_CBS_CONFIG"`
- **AND** the API SHALL query cloud disk configuration list independently

#### Scenario: User specifies inquiry_type as INQUIRY_CVM_CONFIG
- **WHEN** user provides `inquiry_type` with value `"INQUIRY_CVM_CONFIG"` in the `cbs_filter` block
- **THEN** the `DescribeDiskConfigQuota` API call SHALL set `InquiryType` to `"INQUIRY_CVM_CONFIG"`
- **AND** the API SHALL query cloud disk configurations combined with instance types

#### Scenario: User does not specify inquiry_type
- **WHEN** user does not provide `inquiry_type` in the `cbs_filter` block
- **THEN** the `DescribeDiskConfigQuota` API call SHALL default `InquiryType` to `"INQUIRY_CVM_CONFIG"`
- **AND** behavior SHALL be identical to the current hardcoded behavior

#### Scenario: inquiry_type schema definition
- **WHEN** the data source schema is defined
- **THEN** `inquiry_type` SHALL be an optional string field within the `cbs_filter` nested schema
- **AND** `inquiry_type` SHALL have description indicating valid values: `INQUIRY_CBS_CONFIG` and `INQUIRY_CVM_CONFIG`

### Requirement: Instance Families Input Parameter
The `tencentcloud_instance_types` data source SHALL support an `instance_families` input parameter within the `cbs_filter` block to specify instance family names for filtering in the `DescribeDiskConfigQuota` API call.

#### Scenario: User specifies instance_families
- **WHEN** user provides `instance_families` as a list of strings (e.g., `["S5", "M5"]`) in the `cbs_filter` block
- **THEN** the `DescribeDiskConfigQuota` API call SHALL set `InstanceFamilies` to the user-provided list
- **AND** the user-provided list SHALL override the instance type's own family field

#### Scenario: User does not specify instance_families
- **WHEN** user does not provide `instance_families` in the `cbs_filter` block
- **THEN** the `DescribeDiskConfigQuota` API call SHALL derive `InstanceFamilies` from each instance type's `family` field
- **AND** behavior SHALL be identical to the current behavior

#### Scenario: instance_families schema definition
- **WHEN** the data source schema is defined
- **THEN** `instance_families` SHALL be an optional list of strings field within the `cbs_filter` nested schema
- **AND** `instance_families` SHALL have description indicating it overrides the instance type's family for CBS filtering

### Requirement: Disk Charge Type Input Parameter Behavior
The `tencentcloud_instance_types` data source SHALL ensure the existing `disk_charge_type` parameter in the `cbs_filter` block properly passes its value to the `DescribeDiskConfigQuota` API call's `DiskChargeType` field.

#### Scenario: User specifies disk_charge_type in cbs_filter
- **WHEN** user provides `disk_charge_type` with value `"PREPAID"` or `"POSTPAID_BY_HOUR"` in the `cbs_filter` block
- **THEN** the `DescribeDiskConfigQuota` API call SHALL set `DiskChargeType` to the user-provided value
- **AND** this behavior SHALL remain unchanged from the current implementation

#### Scenario: User does not specify disk_charge_type in cbs_filter
- **WHEN** user does not provide `disk_charge_type` in the `cbs_filter` block
- **THEN** the `DescribeDiskConfigQuota` API call SHALL NOT set `DiskChargeType`
- **AND** the API SHALL use its default behavior

### Requirement: Backward Compatibility
The `tencentcloud_instance_types` data source MUST maintain backward compatibility with existing Terraform configurations.

#### Scenario: Existing configurations without new parameters remain functional
- **WHEN** user has existing Terraform configuration using `cbs_filter` with only `disk_types`, `disk_charge_type`, and `disk_usage`
- **THEN** configuration SHALL continue to work without modification
- **AND** `InquiryType` SHALL default to `"INQUIRY_CVM_CONFIG"` as currently hardcoded
- **AND** `InstanceFamilies` SHALL be derived from instance type's family field as currently implemented
- **AND** `DiskChargeType` SHALL be passed from `disk_charge_type` input as currently implemented

#### Scenario: New parameters do not break existing outputs
- **WHEN** user queries instance_types data source with existing configuration
- **THEN** new optional parameters in `cbs_filter` SHALL NOT affect existing computed fields
- **AND** omission of new parameters SHALL NOT cause errors
