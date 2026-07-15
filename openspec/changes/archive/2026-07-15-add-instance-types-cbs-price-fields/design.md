## Context

The `tencentcloud_instance_types` data source queries CVM instance type configurations. When the `cbs_filter` attribute is provided, it also calls the CBS `DescribeDiskConfigQuota` API to populate the `cbs_configs` nested attribute with disk configuration details.

Currently, the `cbs_configs` schema maps 11 out of 12 fields from the CBS `DiskConfig` struct. The `Price` field, which contains disk pricing information (unit price, discount price, original price, and high-precision variants), is not mapped. This prevents users from accessing pricing data for disk configurations.

The CBS SDK `Price` struct (defined in `cbs/v20170312/models.go`) contains 9 fields:
- `ChargeUnit` (*string) - Post-paid billing unit (e.g., "HOUR")
- `DiscountPrice` (*float64) - Pre-paid discount price
- `DiscountPriceHigh` (*string) - High-precision pre-paid discount price
- `OriginalPrice` (*float64) - Pre-paid original price
- `OriginalPriceHigh` (*string) - High-precision pre-paid original price
- `UnitPrice` (*float64) - Post-paid unit price
- `UnitPriceDiscount` (*float64) - Post-paid discount unit price
- `UnitPriceDiscountHigh` (*string) - High-precision post-paid discount unit price
- `UnitPriceHigh` (*string) - High-precision post-paid unit price

## Goals / Non-Goals

**Goals:**
- Expose all 9 pricing fields from the `Price` struct within the `cbs_configs` nested schema
- Maintain backward compatibility — existing configurations and state must not break
- Follow existing code patterns in the data source for schema definition and data mapping

**Non-Goals:**
- Do not add pricing fields from the `DescribeZoneInstanceConfigInfos` API (that is covered by the existing `price` attribute within `instance_types`)
- Do not modify the CBS service layer function (`DescribeDiskConfigQuota`) — it already returns the full `DiskConfig` set including `Price`
- Do not add new filter or query parameters to the data source

## Decisions

### Decision 1: Flat fields vs. nested Price schema
**Choice**: Add pricing fields as flat computed attributes directly within the `cbs_configs` nested schema.

**Rationale**: The existing `cbs_configs` schema already uses flat fields. Adding a nested `price` block within `cbs_configs` would create unnecessary nesting depth. Flat fields are consistent with how the other `DiskConfig` fields (like `disk_charge_type`, `zone`) are already exposed in `cbs_configs`.

**Alternative considered**: Create a nested `price` block within `cbs_configs` — rejected because it adds unnecessary schema nesting and deviates from the existing flat pattern.

### Decision 2: High-precision price fields as TypeString
**Choice**: Use `TypeString` for high-precision price fields (`discount_price_high`, `original_price_high`, `unit_price_discount_high`, `unit_price_high`).

**Rationale**: The SDK returns these as `*string` type to preserve precision (avoiding floating-point truncation). Using `TypeString` in Terraform preserves this precision. This is consistent with how other high-precision price fields are handled in the provider.

### Decision 3: Standard price fields as TypeFloat
**Choice**: Use `TypeFloat` for standard price fields (`discount_price`, `original_price`, `unit_price`, `unit_price_discount`).

**Rationale**: The SDK returns these as `*float64`. Using `TypeFloat` matches the SDK type and is consistent with how pricing fields are handled in the existing `price` attribute within `instance_types`.

## Risks / Trade-offs

- **[Risk] Nil pointer dereference when Price is nil** → Mitigation: Check `diskConfig.Price != nil` before accessing Price fields, following the same nil-check pattern used for `instanceType.Price` in the existing code.
- **[Risk] High-precision fields may be empty strings** → Mitigation: These are computed fields; empty/nil values will simply not be set. No user impact since they are optional computed fields.
