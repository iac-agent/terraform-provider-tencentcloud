## ADDED Requirements

### Requirement: cbs_configs includes disk pricing fields from DescribeDiskConfigQuota API
The `cbs_configs` nested block in the `tencentcloud_instance_types` data source SHALL include the following 9 computed fields sourced from `DiskConfig.Price` in the CBS `DescribeDiskConfigQuota` API response:
- `charge_unit` (TypeString) - Billing unit for postpaid cloud disks (e.g., HOUR)
- `unit_price` (TypeFloat) - Original unit price for postpaid cloud disks
- `unit_price_discount` (TypeFloat) - Discounted unit price for postpaid cloud disks
- `unit_price_high` (TypeString) - High-precision original unit price
- `unit_price_discount_high` (TypeString) - High-precision discounted unit price
- `original_price` (TypeFloat) - Original price for prepaid cloud disks
- `original_price_high` (TypeString) - High-precision original price for prepaid cloud disks
- `discount_price` (TypeFloat) - Discount price for prepaid cloud disks
- `discount_price_high` (TypeString) - High-precision discount price for prepaid cloud disks

#### Scenario: cbs_configs populated with pricing data when Price is available
- **WHEN** a user provides `cbs_filter` in the `instance_types` data source and the `DescribeDiskConfigQuota` API returns `DiskConfig` items with a non-nil `Price` field
- **THEN** each `cbs_configs` block SHALL contain all 9 pricing fields mapped from the corresponding `DiskConfig.Price` structure

#### Scenario: cbs_configs pricing fields omitted when Price is nil
- **WHEN** a user provides `cbs_filter` and the `DescribeDiskConfigQuota` API returns `DiskConfig` items with a nil `Price` field
- **THEN** the pricing fields SHALL be omitted from the `cbs_configs` block without causing errors

#### Scenario: backward compatibility with existing configurations
- **WHEN** a user uses the `instance_types` data source without `cbs_filter` or with existing configurations
- **THEN** the data source SHALL behave identically to before, with no breaking changes to existing state or configuration
