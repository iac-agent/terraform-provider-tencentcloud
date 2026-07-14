## ADDED Requirements

### Requirement: InstanceFamilies Input Parameter
The data source SHALL accept an `InstanceFamilies` top-level input parameter to filter instance types by their instance family names.

#### Scenario: Filter by instance families
- **WHEN** user provides `InstanceFamilies` parameter with one or more instance family names (e.g., ["S5", "M5"])
- **THEN** the data source SHALL only return instance types belonging to the specified families
- **AND** `InstanceFamilies` SHALL be an Optional parameter of type TypeList with TypeString elements
- **AND** the values SHALL be translated to `instance-family` filter values when calling `DescribeZoneInstanceConfigInfos` API

#### Scenario: InstanceFamilies combined with existing filter block
- **WHEN** user provides both `InstanceFamilies` parameter and a `filter` block with `instance-family` filter name
- **THEN** the data source SHALL merge values from both sources into a single `instance-family` filter
- **AND** the combined filter SHALL contain all unique instance family names from both sources

#### Scenario: InstanceFamilies without filter block
- **WHEN** user provides `InstanceFamilies` parameter without any `filter` block
- **THEN** the data source SHALL create an `instance-family` filter using the `InstanceFamilies` values
- **AND** the data source SHALL return only instance types matching the specified families

#### Scenario: InstanceFamilies used with DescribeDiskConfigQuota
- **WHEN** user provides `InstanceFamilies` parameter and `cbs_filter` or `DiskTypes` is also provided (triggering CBS config query)
- **THEN** the `InstanceFamilies` values SHALL be passed directly to the `DescribeDiskConfigQuota` API's `InstanceFamilies` request field
- **AND** the CBS config results SHALL be filtered by the specified instance families

#### Scenario: InstanceFamilies not provided
- **WHEN** user does not provide `InstanceFamilies` parameter
- **THEN** the data source SHALL behave exactly as before, with no change to filtering behavior
- **AND** for CBS config queries, instance families SHALL continue to be auto-populated from each instance type's family field

### Requirement: DiskTypes Input Parameter
The data source SHALL accept a `DiskTypes` top-level input parameter to filter disk configurations by disk media type.

#### Scenario: Filter by disk types
- **WHEN** user provides `DiskTypes` parameter with one or more disk type values (e.g., ["CLOUD_SSD", "CLOUD_PREMIUM"])
- **THEN** the data source SHALL pass the `DiskTypes` values directly to the `DescribeDiskConfigQuota` API's `DiskTypes` request field
- **AND** the CBS config results SHALL be filtered by the specified disk types
- **AND** `DiskTypes` SHALL be an Optional parameter of type TypeList with TypeString elements

#### Scenario: DiskTypes triggers CBS config query
- **WHEN** user provides `DiskTypes` parameter without `cbs_filter` block
- **THEN** the data source SHALL trigger CBS config queries for each instance type
- **AND** `DiskTypes` values SHALL be passed to each `DescribeDiskConfigQuota` API call

#### Scenario: DiskTypes combined with cbs_filter
- **WHEN** user provides both `DiskTypes` top-level parameter and `cbs_filter.disk_types`
- **THEN** the top-level `DiskTypes` parameter SHALL take precedence
- **AND** the `DiskTypes` values SHALL be used instead of `cbs_filter.disk_types` in the `DescribeDiskConfigQuota` API call

#### Scenario: DiskTypes not provided
- **WHEN** user does not provide `DiskTypes` parameter and does not provide `cbs_filter`
- **THEN** the data source SHALL behave exactly as before, with no CBS config query triggered
- **AND** backward compatibility SHALL be maintained

### Requirement: Backward Compatibility
The data source MUST maintain backward compatibility with existing Terraform configurations.

#### Scenario: Existing configurations remain functional
- **WHEN** user has existing Terraform configuration using instance_types data source without `InstanceFamilies` or `DiskTypes` parameters
- **THEN** configuration SHALL continue to work without modification
- **AND** terraform plan SHALL not show any changes
- **AND** data source output SHALL contain same values as before for existing fields

#### Scenario: New parameters do not break existing outputs
- **WHEN** user queries instance_types data source with existing configuration
- **THEN** omission of `InstanceFamilies` and `DiskTypes` parameters SHALL not cause errors
- **AND** existing `filter` and `cbs_filter` mechanisms SHALL continue to work as before
