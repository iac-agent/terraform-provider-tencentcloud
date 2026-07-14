# datasource-cvm-instance-types-params Specification

## Purpose
TBD - created by syncing change add-instance-types-datasource-params. Update Purpose after archive.

## Requirements

### Requirement: Disk Types Input Parameter
The data source SHALL accept a `disk_types` top-level optional input parameter that specifies disk media types for the DescribeDiskConfigQuota API call.

#### Scenario: Specify disk types at top level
- **WHEN** user provides `disk_types` parameter at the data source top level with values such as `["CLOUD_SSD", "CLOUD_PREMIUM"]`
- **AND** `cbs_filter` is also provided
- **THEN** the `disk_types` top-level values SHALL be used as `request.DiskTypes` in the DescribeDiskConfigQuota API call
- **AND** top-level `disk_types` SHALL override `cbs_filter.disk_types` if both are specified

#### Scenario: Disk types not specified at top level
- **WHEN** user does NOT provide `disk_types` at the data source top level
- **AND** `cbs_filter.disk_types` is provided
- **THEN** `cbs_filter.disk_types` SHALL continue to be used as `request.DiskTypes` in the DescribeDiskConfigQuota API call

#### Scenario: Disk types parameter schema
- **WHEN** the data source schema is defined
- **THEN** `disk_types` SHALL be an optional field of type TypeList with TypeString elements
- **AND** `disk_types` SHALL have description listing valid values: CLOUD_BASIC, CLOUD_PREMIUM, CLOUD_SSD, CLOUD_HSSD

### Requirement: Zones Input Parameter
The data source SHALL accept a `zones` top-level optional input parameter that specifies availability zones for the DescribeDiskConfigQuota API call.

#### Scenario: Specify zones at top level
- **WHEN** user provides `zones` parameter at the data source top level with values such as `["ap-guangzhou-3", "ap-guangzhou-4"]`
- **AND** `cbs_filter` is also provided
- **THEN** the `zones` top-level values SHALL be used as `request.Zones` in the DescribeDiskConfigQuota API call
- **AND** top-level `zones` SHALL override the derived `availability_zone` from instance type results

#### Scenario: Zones not specified at top level
- **WHEN** user does NOT provide `zones` at the data source top level
- **THEN** the derived `availability_zone` from instance type results SHALL continue to be used as `request.Zones` in the DescribeDiskConfigQuota API call

#### Scenario: Zones parameter schema
- **WHEN** the data source schema is defined
- **THEN** `zones` SHALL be an optional field of type TypeList with TypeString elements
- **AND** `zones` SHALL have description explaining it specifies availability zones for disk config quota query

### Requirement: Memory Input Parameter
The data source SHALL accept a `memory` top-level optional input parameter that specifies instance memory size for the DescribeDiskConfigQuota API call.

#### Scenario: Specify memory at top level
- **WHEN** user provides `memory` parameter at the data source top level with a value such as `8`
- **AND** `cbs_filter` is also provided
- **THEN** the `memory` top-level value SHALL be used as `request.Memory` in the DescribeDiskConfigQuota API call
- **AND** top-level `memory` SHALL override the derived `memory_size` from instance type results

#### Scenario: Memory not specified at top level
- **WHEN** user does NOT provide `memory` at the data source top level
- **THEN** the derived `memory_size` from instance type results SHALL continue to be used as `request.Memory` in the DescribeDiskConfigQuota API call

#### Scenario: Memory parameter schema
- **WHEN** the data source schema is defined
- **THEN** `memory` SHALL be an optional field of type TypeInt
- **AND** `memory` SHALL have description explaining it specifies instance memory size in GB for disk config quota query

### Requirement: Backward Compatibility
The data source MUST maintain backward compatibility with existing Terraform configurations.

#### Scenario: Existing configurations remain functional
- **WHEN** user has existing Terraform configuration using instance_types data source without the new top-level parameters
- **THEN** configuration SHALL continue to work without modification
- **AND** the DescribeDiskConfigQuota call SHALL use derived values as before

#### Scenario: New parameters do not affect DescribeZoneInstanceConfigInfos
- **WHEN** user provides `disk_types`, `zones`, or `memory` at the data source top level
- **THEN** these parameters SHALL NOT affect the DescribeZoneInstanceConfigInfos API call
- **AND** the existing `availability_zone`, `memory_size`, `cpu_core_count`, `gpu_core_count` parameters SHALL continue to filter instance types as before
