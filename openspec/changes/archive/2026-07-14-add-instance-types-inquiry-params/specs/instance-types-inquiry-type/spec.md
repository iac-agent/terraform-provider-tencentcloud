## ADDED Requirements

### Requirement: Configurable Inquiry Type for CBS Disk Config Query
The data source SHALL allow users to specify the `inquiry_type` parameter within the `cbs_filter` block to control the query category of the `DescribeDiskConfigQuota` CBS API call.

#### Scenario: User specifies inquiry_type as INQUIRY_CBS_CONFIG
- **WHEN** user provides `cbs_filter` with `inquiry_type` set to `"INQUIRY_CBS_CONFIG"`
- **THEN** the `DescribeDiskConfigQuota` API SHALL be called with `InquiryType` parameter set to `"INQUIRY_CBS_CONFIG"`
- **AND** the data source SHALL return disk configuration results without instance pairing constraints

#### Scenario: User specifies inquiry_type as INQUIRY_CVM_CONFIG
- **WHEN** user provides `cbs_filter` with `inquiry_type` set to `"INQUIRY_CVM_CONFIG"`
- **THEN** the `DescribeDiskConfigQuota` API SHALL be called with `InquiryType` parameter set to `"INQUIRY_CVM_CONFIG"`
- **AND** the data source SHALL return disk configuration results with instance pairing constraints

#### Scenario: User does not specify inquiry_type in cbs_filter
- **WHEN** user provides `cbs_filter` without specifying `inquiry_type`
- **THEN** the `DescribeDiskConfigQuota` API SHALL be called with `InquiryType` parameter defaulted to `"INQUIRY_CVM_CONFIG"`
- **AND** the data source behavior SHALL be identical to the current hardcoded behavior
- **AND** existing Terraform configurations SHALL continue to work without modification

### Requirement: Inquiry Type Schema Definition
The `inquiry_type` field SHALL be defined as an optional string field within the `cbs_filter` nested schema block.

#### Scenario: inquiry_type schema field properties
- **WHEN** the data source schema is defined
- **THEN** `inquiry_type` SHALL be a field of type `schema.TypeString`
- **AND** `inquiry_type` SHALL be `Optional` with no `Default` value set in schema (default applied in code)
- **AND** `inquiry_type` SHALL have a description listing valid values: `"INQUIRY_CBS_CONFIG"` (query cloud disk configuration list) and `"INQUIRY_CVM_CONFIG"` (query cloud disk and instance pairing configuration list)
- **AND** `inquiry_type` SHALL NOT be `Required` to maintain backward compatibility

### Requirement: Backward Compatibility
The addition of `inquiry_type` SHALL NOT break any existing Terraform configuration using the `tencentcloud_instance_types` data source.

#### Scenario: Existing cbs_filter configurations without inquiry_type
- **WHEN** user has an existing configuration using `cbs_filter` without `inquiry_type`
- **THEN** the data source SHALL continue to work as before
- **AND** `terraform plan` SHALL NOT show any changes for existing configurations
- **AND** `DescribeDiskConfigQuota` SHALL continue to use `"INQUIRY_CVM_CONFIG"` as default
