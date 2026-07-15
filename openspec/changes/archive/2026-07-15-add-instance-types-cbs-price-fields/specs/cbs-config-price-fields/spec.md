## ADDED Requirements

### Requirement: CBS Config Pricing Fields
The `cbs_configs` nested attribute within the `instance_types` data source SHALL expose all pricing fields from the `Price` structure returned by the `DescribeDiskConfigQuota` API.

#### Scenario: Query CBS config charge unit
- **WHEN** user queries instance types data source with `cbs_filter` parameter
- **THEN** each CBS config entry SHALL include `charge_unit` field indicating the post-paid billing unit
- **AND** `charge_unit` field SHALL be optional computed field of type string
- **AND** `charge_unit` field SHALL be null if not provided by API

#### Scenario: Query CBS config discount price
- **WHEN** user queries instance types data source with `cbs_filter` parameter
- **THEN** each CBS config entry SHALL include `discount_price` field indicating the pre-paid discount price in CNY
- **AND** `discount_price` field SHALL be optional computed field of type float
- **AND** `discount_price` field SHALL be null if not provided by API

#### Scenario: Query CBS config discount price high precision
- **WHEN** user queries instance types data source with `cbs_filter` parameter
- **THEN** each CBS config entry SHALL include `discount_price_high` field indicating the high-precision pre-paid discount price in CNY
- **AND** `discount_price_high` field SHALL be optional computed field of type string
- **AND** `discount_price_high` field SHALL be null if not provided by API

#### Scenario: Query CBS config original price
- **WHEN** user queries instance types data source with `cbs_filter` parameter
- **THEN** each CBS config entry SHALL include `original_price` field indicating the pre-paid original price in CNY
- **AND** `original_price` field SHALL be optional computed field of type float
- **AND** `original_price` field SHALL be null if not provided by API

#### Scenario: Query CBS config original price high precision
- **WHEN** user queries instance types data source with `cbs_filter` parameter
- **THEN** each CBS config entry SHALL include `original_price_high` field indicating the high-precision pre-paid original price in CNY
- **AND** `original_price_high` field SHALL be optional computed field of type string
- **AND** `original_price_high` field SHALL be null if not provided by API

#### Scenario: Query CBS config unit price
- **WHEN** user queries instance types data source with `cbs_filter` parameter
- **THEN** each CBS config entry SHALL include `unit_price` field indicating the post-paid unit price in CNY
- **AND** `unit_price` field SHALL be optional computed field of type float
- **AND** `unit_price` field SHALL be null if not provided by API

#### Scenario: Query CBS config unit price discount
- **WHEN** user queries instance types data source with `cbs_filter` parameter
- **THEN** each CBS config entry SHALL include `unit_price_discount` field indicating the post-paid discount unit price in CNY
- **AND** `unit_price_discount` field SHALL be optional computed field of type float
- **AND** `unit_price_discount` field SHALL be null if not provided by API

#### Scenario: Query CBS config unit price discount high precision
- **WHEN** user queries instance types data source with `cbs_filter` parameter
- **THEN** each CBS config entry SHALL include `unit_price_discount_high` field indicating the high-precision post-paid discount unit price in CNY
- **AND** `unit_price_discount_high` field SHALL be optional computed field of type string
- **AND** `unit_price_discount_high` field SHALL be null if not provided by API

#### Scenario: Query CBS config unit price high precision
- **WHEN** user queries instance types data source with `cbs_filter` parameter
- **THEN** each CBS config entry SHALL include `unit_price_high` field indicating the high-precision post-paid unit price in CNY
- **AND** `unit_price_high` field SHALL be optional computed field of type string
- **AND** `unit_price_high` field SHALL be null if not provided by API

### Requirement: Backward Compatibility
The data source MUST maintain backward compatibility with existing Terraform configurations.

#### Scenario: Existing configurations remain functional
- **WHEN** user has existing Terraform configuration using instance_types data source
- **AND** configuration references only previously supported `cbs_configs` fields (available, disk_charge_type, zone, instance_family, disk_type, step_size, extra_performance_range, device_class, disk_usage, min_disk_size, max_disk_size)
- **THEN** configuration SHALL continue to work without modification
- **AND** terraform plan SHALL not show any changes

#### Scenario: New pricing fields are populated when Price is available
- **WHEN** user queries instance_types data source with `cbs_filter` parameter
- **AND** the DescribeDiskConfigQuota API returns Price data
- **THEN** the 9 new pricing fields SHALL be populated within `cbs_configs`
- **AND** if Price data is nil, the pricing fields SHALL be null or empty
