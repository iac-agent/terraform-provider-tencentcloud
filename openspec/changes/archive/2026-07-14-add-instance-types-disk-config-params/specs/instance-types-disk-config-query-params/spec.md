## ADDED Requirements

### Requirement: Disk types parameter for DescribeDiskConfigQuota
The data source SHALL accept a `disk_types` optional input parameter that maps to the `request.DiskTypes` field of the `DescribeDiskConfigQuota` API.

#### Scenario: User specifies disk_types to filter CBS disk configurations
- **WHEN** user provides `disk_types` parameter in the data source configuration (e.g., `disk_types = ["CLOUD_SSD", "CLOUD_HSSD"]`)
- **AND** `cbs_filter` is also provided
- **THEN** the `DescribeDiskConfigQuota` API request SHALL use the `disk_types` value as the `DiskTypes` parameter
- **AND** the `disk_types` value SHALL override `cbs_filter.disk_types` if both are specified

#### Scenario: User does not specify disk_types
- **WHEN** user does not provide `disk_types` parameter
- **AND** `cbs_filter` is provided with `disk_types` inside it
- **THEN** the `DescribeDiskConfigQuota` API request SHALL use `cbs_filter.disk_types` as the `DiskTypes` parameter (existing behavior preserved)

#### Scenario: Neither top-level disk_types nor cbs_filter.disk_types is specified
- **WHEN** user does not provide `disk_types` parameter
- **AND** `cbs_filter` is provided without `disk_types` inside it
- **THEN** the `DescribeDiskConfigQuota` API request SHALL NOT include `DiskTypes` parameter

### Requirement: Zones parameter for DescribeDiskConfigQuota
The data source SHALL accept a `zones` optional input parameter that maps to the `request.Zones` field of the `DescribeDiskConfigQuota` API.

#### Scenario: User specifies zones to override CBS query zones
- **WHEN** user provides `zones` parameter in the data source configuration (e.g., `zones = ["ap-guangzhou-1", "ap-guangzhou-2"]`)
- **AND** `cbs_filter` is also provided
- **THEN** the `DescribeDiskConfigQuota` API request SHALL use the `zones` value as the `Zones` parameter
- **AND** the `zones` value SHALL override the instance type's `availability_zone` derived value

#### Scenario: User does not specify zones
- **WHEN** user does not provide `zones` parameter
- **AND** `cbs_filter` is provided
- **THEN** the `DescribeDiskConfigQuota` API request SHALL use each instance type's `availability_zone` as the `Zones` parameter (existing behavior preserved)

### Requirement: Memory parameter for DescribeDiskConfigQuota
The data source SHALL accept a `memory` optional input parameter that maps to the `request.Memory` field of the `DescribeDiskConfigQuota` API.

#### Scenario: User specifies memory to override CBS query memory
- **WHEN** user provides `memory` parameter in the data source configuration (e.g., `memory = 8`)
- **AND** `cbs_filter` is also provided
- **THEN** the `DescribeDiskConfigQuota` API request SHALL use the `memory` value as the `Memory` parameter
- **AND** the `memory` value SHALL override the instance type's `memory_size` derived value

#### Scenario: User does not specify memory
- **WHEN** user does not provide `memory` parameter
- **AND** `cbs_filter` is provided
- **THEN** the `DescribeDiskConfigQuota` API request SHALL use each instance type's `memory_size` as the `Memory` parameter (existing behavior preserved)

### Requirement: Backward compatibility
The data source MUST maintain backward compatibility with existing Terraform configurations.

#### Scenario: Existing configurations without new parameters
- **WHEN** user has an existing configuration using `tencentcloud_instance_types` data source without `disk_types`, `zones`, or `memory` parameters
- **THEN** the data source SHALL behave identically to the previous version
- **AND** all existing parameters and output fields SHALL continue to work without modification

#### Scenario: New parameters without cbs_filter
- **WHEN** user provides `disk_types`, `zones`, or `memory` parameters
- **AND** `cbs_filter` is NOT provided
- **THEN** these parameters SHALL have no effect on the data source behavior (since `DescribeDiskConfigQuota` is only called when `cbs_filter` is present)
