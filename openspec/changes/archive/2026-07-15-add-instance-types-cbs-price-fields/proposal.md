## Why

The `tencentcloud_instance_types` data source currently exposes CBS disk configuration information through the `cbs_configs` computed attribute, which is populated by the `DescribeDiskConfigQuota` API. However, the `Price` field from the `DiskConfig` response structure is not mapped, preventing users from accessing disk pricing information (such as unit price, discount price, original price, and high-precision price variants) when querying instance types with CBS filters. This information is essential for cost planning and optimization.

## What Changes

- Add 9 new computed fields to the `cbs_configs` nested schema in the `instance_types` data source, mapping all fields from the `Price` structure returned by the `DescribeDiskConfigQuota` API
- New fields: `charge_unit`, `discount_price`, `discount_price_high`, `original_price`, `original_price_high`, `unit_price`, `unit_price_discount`, `unit_price_discount_high`, `unit_price_high`
- Update data source mapping logic to populate the new fields from `DiskConfig.Price`
- All new fields are optional computed fields — no breaking changes

## Capabilities

### New Capabilities
- `cbs-config-price-fields`: Expose pricing fields from the DescribeDiskConfigQuota API's Price structure within the cbs_configs attribute of the instance_types data source

### Modified Capabilities

## Impact

### Affected Code
- `tencentcloud/services/cvm/data_source_tc_instance_types.go` - Add new Price fields to `cbs_configs` schema and mapping logic
- `tencentcloud/services/cvm/data_source_tc_instance_types.md` - Update documentation with new fields

### Affected APIs
- `DescribeDiskConfigQuota` (CBS API) - No API changes; only mapping additional response fields

### Breaking Changes
None - all changes are additive (new computed fields within existing `cbs_configs` nested schema)

### Dependencies
None - uses existing `DescribeDiskConfigQuota` API call; the `Price` field is already returned by the API but not mapped
