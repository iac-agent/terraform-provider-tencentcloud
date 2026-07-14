## ADDED Requirements

### Requirement: inquiry_type parameter for tencentcloud_instance_types datasource

The `tencentcloud_instance_types` datasource SHALL accept an optional top-level parameter `inquiry_type` (type: string) that specifies the query category for the `DescribeDiskConfigQuota` API call.

Valid values SHALL be:
- `INQUIRY_CBS_CONFIG`: Query cloud disk configuration list only
- `INQUIRY_CVM_CONFIG`: Query cloud disk and instance combination configuration list

When `inquiry_type` is not specified, it SHALL default to `INQUIRY_CVM_CONFIG` to maintain backward compatibility with existing configurations.

The `inquiry_type` parameter SHALL only be used when `cbs_filter` is present, as it controls the behavior of the `DescribeDiskConfigQuota` API call.

#### Scenario: User specifies inquiry_type as INQUIRY_CBS_CONFIG
- **WHEN** user sets `inquiry_type = "INQUIRY_CBS_CONFIG"` and provides `cbs_filter`
- **THEN** the `DescribeDiskConfigQuota` API request SHALL have `InquiryType` set to `INQUIRY_CBS_CONFIG`

#### Scenario: User does not specify inquiry_type
- **WHEN** user provides `cbs_filter` without specifying `inquiry_type`
- **THEN** the `DescribeDiskConfigQuota` API request SHALL have `InquiryType` set to `INQUIRY_CVM_CONFIG` (default)

#### Scenario: User specifies inquiry_type without cbs_filter
- **WHEN** user sets `inquiry_type` but does not provide `cbs_filter`
- **THEN** the `inquiry_type` parameter SHALL be ignored since `DescribeDiskConfigQuota` is not called

### Requirement: disk_charge_type parameter for tencentcloud_instance_types datasource

The `tencentcloud_instance_types` datasource SHALL accept an optional top-level parameter `disk_charge_type` (type: string) that specifies the payment model for the `DescribeDiskConfigQuota` API call.

Valid values SHALL be:
- `PREPAID`: Prepaid mode
- `POSTPAID_BY_HOUR`: Postpaid by hour mode

When both the top-level `disk_charge_type` and `cbs_filter.disk_charge_type` are specified, the top-level `disk_charge_type` SHALL take precedence.

When only `cbs_filter.disk_charge_type` is specified (no top-level `disk_charge_type`), the `cbs_filter.disk_charge_type` SHALL continue to be used for backward compatibility.

When neither is specified, `DiskChargeType` SHALL not be set in the API request.

#### Scenario: User specifies top-level disk_charge_type
- **WHEN** user sets `disk_charge_type = "PREPAID"` and provides `cbs_filter`
- **THEN** the `DescribeDiskConfigQuota` API request SHALL have `DiskChargeType` set to `PREPAID`

#### Scenario: User specifies both top-level disk_charge_type and cbs_filter.disk_charge_type
- **WHEN** user sets top-level `disk_charge_type = "PREPAID"` and also sets `cbs_filter.disk_charge_type = "POSTPAID_BY_HOUR"`
- **THEN** the `DescribeDiskConfigQuota` API request SHALL have `DiskChargeType` set to `PREPAID` (top-level takes precedence)

#### Scenario: User specifies only cbs_filter.disk_charge_type (backward compatibility)
- **WHEN** user does not set top-level `disk_charge_type` but sets `cbs_filter.disk_charge_type = "POSTPAID_BY_HOUR"`
- **THEN** the `DescribeDiskConfigQuota` API request SHALL have `DiskChargeType` set to `POSTPAID_BY_HOUR` (existing behavior preserved)

#### Scenario: User specifies neither disk_charge_type
- **WHEN** user does not set top-level `disk_charge_type` and does not set `cbs_filter.disk_charge_type`
- **THEN** the `DescribeDiskConfigQuota` API request SHALL NOT have `DiskChargeType` set
