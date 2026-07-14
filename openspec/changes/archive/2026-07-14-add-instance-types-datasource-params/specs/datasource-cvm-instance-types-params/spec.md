## ADDED Requirements

### Requirement: InquiryType parameter
The `tencentcloud_instance_types` data source SHALL accept an optional top-level `InquiryType` parameter (TypeString) that specifies the query category for the `DescribeDiskConfigQuota` API call. Valid values SHALL be "INQUIRY_CBS_CONFIG" and "INQUIRY_CVM_CONFIG". When not specified, the default SHALL be "INQUIRY_CVM_CONFIG" to maintain backward compatibility.

#### Scenario: User specifies InquiryType as INQUIRY_CBS_CONFIG
- **WHEN** user provides `InquiryType = "INQUIRY_CBS_CONFIG"` in the data source configuration
- **THEN** the `DescribeDiskConfigQuota` API request SHALL set `InquiryType` to "INQUIRY_CBS_CONFIG" instead of the default "INQUIRY_CVM_CONFIG"

#### Scenario: User does not specify InquiryType
- **WHEN** user does not provide the `InquiryType` parameter
- **THEN** the `DescribeDiskConfigQuota` API request SHALL use "INQUIRY_CVM_CONFIG" as the default value, preserving existing behavior

### Requirement: DiskChargeType parameter
The `tencentcloud_instance_types` data source SHALL accept an optional top-level `DiskChargeType` parameter (TypeString) that specifies the disk payment model for the `DescribeDiskConfigQuota` API call. Valid values SHALL be "PREPAID" and "POSTPAID_BY_HOUR". When specified at the top level, this value SHALL be passed to the `DescribeDiskConfigQuota` request.

#### Scenario: User specifies top-level DiskChargeType
- **WHEN** user provides `DiskChargeType = "PREPAID"` at the data source top level
- **THEN** the `DescribeDiskConfigQuota` API request SHALL set `DiskChargeType` to "PREPAID"

#### Scenario: User does not specify top-level DiskChargeType
- **WHEN** user does not provide the top-level `DiskChargeType` parameter but provides `cbs_filter.disk_charge_type`
- **THEN** the `cbs_filter.disk_charge_type` value SHALL be used for the `DescribeDiskConfigQuota` request as before

#### Scenario: User specifies both top-level DiskChargeType and cbs_filter disk_charge_type
- **WHEN** user provides both top-level `DiskChargeType` and `cbs_filter.disk_charge_type`
- **THEN** the top-level `DiskChargeType` value SHALL take precedence in the `DescribeDiskConfigQuota` request

### Requirement: InstanceFamilies parameter
The `tencentcloud_instance_types` data source SHALL accept an optional top-level `InstanceFamilies` parameter (TypeList of TypeString) that specifies instance family filters for the `DescribeDiskConfigQuota` API call. When explicitly provided, this value SHALL override the auto-populated instance family derived from instance type query results. When not provided, the instance family SHALL continue to be auto-populated from the `DescribeZoneInstanceConfigInfos` results.

#### Scenario: User specifies InstanceFamilies explicitly
- **WHEN** user provides `InstanceFamilies = ["S5", "M5"]` at the data source top level
- **THEN** the `DescribeDiskConfigQuota` API request SHALL set `InstanceFamilies` to the user-provided list instead of auto-populating from instance type results

#### Scenario: User does not specify InstanceFamilies
- **WHEN** user does not provide the top-level `InstanceFamilies` parameter
- **THEN** the instance family SHALL be auto-populated from each instance type's `family` field as in current behavior

### Requirement: Available computed attribute
The `tencentcloud_instance_types` data source SHALL expose a top-level computed `Available` attribute (TypeBool) that indicates whether disk configurations are available from the `DescribeDiskConfigQuota` API response. When the `DescribeDiskConfigQuota` API is called and returns disk configurations, `Available` SHALL be set to true if at least one disk configuration has `Available = true`. If no disk configurations are available or the API is not called, `Available` SHALL be set to false.

#### Scenario: Disk configurations are available
- **WHEN** `DescribeDiskConfigQuota` is called and at least one returned `DiskConfig` has `Available = true`
- **THEN** the top-level `Available` attribute SHALL be set to true

#### Scenario: No disk configurations are available
- **WHEN** `DescribeDiskConfigQuota` is called and all returned `DiskConfig` entries have `Available = false`
- **THEN** the top-level `Available` attribute SHALL be set to false

#### Scenario: DescribeDiskConfigQuota is not called
- **WHEN** neither `cbs_filter` nor any new top-level CBS parameters are provided
- **THEN** the top-level `Available` attribute SHALL be set to false (no disk config query was performed)
