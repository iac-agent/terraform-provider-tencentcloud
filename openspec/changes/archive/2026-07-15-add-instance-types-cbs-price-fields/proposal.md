## Why

The `tencentcloud_instance_types` data source's `cbs_configs` nested block currently maps most fields from the CBS `DescribeDiskConfigQuota` API's `DiskConfig` response, but omits the pricing fields available in `DiskConfig.Price`. Users querying instance types with CBS filter cannot see disk pricing information (unit price, original price, discount price, and their high-precision variants), which limits their ability to perform cost analysis when selecting cloud disk configurations.

## What Changes

- Add 9 new computed fields to the `cbs_configs` nested schema within the `instance_types` data source, sourced from `DiskConfig.Price` in the CBS `DescribeDiskConfigQuota` API response:
  - `charge_unit` (from `Price.ChargeUnit`) - Billing unit for postpaid disks
  - `unit_price` (from `Price.UnitPrice`) - Original unit price for postpaid disks
  - `unit_price_discount` (from `Price.UnitPriceDiscount`) - Discounted unit price for postpaid disks
  - `unit_price_high` (from `Price.UnitPriceHigh`) - High-precision original unit price
  - `unit_price_discount_high` (from `Price.UnitPriceDiscountHigh`) - High-precision discounted unit price
  - `original_price` (from `Price.OriginalPrice`) - Original price for prepaid disks
  - `original_price_high` (from `Price.OriginalPriceHigh`) - High-precision original price for prepaid disks
  - `discount_price` (from `Price.DiscountPrice`) - Discount price for prepaid disks
  - `discount_price_high` (from `Price.DiscountPriceHigh`) - High-precision discount price for prepaid disks

## Capabilities

### New Capabilities
- `instance-types-cbs-pricing`: Adds pricing fields to the cbs_configs block of the instance_types data source, exposing disk pricing information from the DescribeDiskConfigQuota API

### Modified Capabilities

## Impact

- `tencentcloud/services/cvm/data_source_tc_instance_types.go` - Add new computed fields to `cbs_configs` nested schema and mapping logic
- `tencentcloud/services/cvm/data_source_tc_instance_types_test.go` - Add test coverage for new pricing fields
- `tencentcloud/services/cvm/data_source_tc_instance_types.md` - Update documentation with new fields
- No breaking changes - all new fields are optional computed fields
