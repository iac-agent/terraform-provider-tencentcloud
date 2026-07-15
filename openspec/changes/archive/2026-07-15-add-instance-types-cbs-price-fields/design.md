## Context

The `tencentcloud_instance_types` data source queries the CVM `DescribeZoneInstanceConfigInfos` API for instance type information, and optionally calls the CBS `DescribeDiskConfigQuota` API (via `cbs_filter`) to populate the `cbs_configs` nested block. The `DiskConfig` response from CBS includes a `Price` field of type `*Price` containing disk pricing information, but the current Terraform schema only maps the non-pricing fields from `DiskConfig` (e.g., `Available`, `DiskChargeType`, `Zone`, `DiskType`, etc.), leaving pricing data inaccessible to Terraform users.

The CBS SDK's `Price` struct (in `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cbs/v20170312`) contains the following fields:
- `UnitPriceDiscount` (*float64) - Discounted unit price for postpaid
- `DiscountPrice` (*float64) - Discount price for prepaid
- `UnitPrice` (*float64) - Original unit price for postpaid
- `UnitPriceHigh` (*string) - High-precision original unit price
- `OriginalPriceHigh` (*string) - High-precision original price for prepaid
- `OriginalPrice` (*float64) - Original price for prepaid
- `DiscountPriceHigh` (*string) - High-precision discount price for prepaid
- `UnitPriceDiscountHigh` (*string) - High-precision discounted unit price
- `ChargeUnit` (*string) - Billing unit (e.g., HOUR)

## Goals / Non-Goals

**Goals:**
- Add all 9 pricing fields from `DiskConfig.Price` to the `cbs_configs` nested schema in the `instance_types` data source
- Map the pricing fields from the CBS API response to the Terraform schema in the read function
- Ensure backward compatibility - all new fields are computed and optional

**Non-Goals:**
- Do not modify the `DescribeZoneInstanceConfigInfos` API mapping (covered by existing change)
- Do not add any new input/filter parameters
- Do not modify the existing non-pricing fields in `cbs_configs`

## Decisions

1. **Flat schema over nested Price block**: Add pricing fields directly to the `cbs_configs` nested schema rather than creating a nested `price` sub-block. Rationale: The `cbs_configs` block already has flat fields from `DiskConfig`; adding a nested `price` block would be inconsistent with the existing pattern. Users can identify pricing fields by their `unit_price*`, `original_price*`, `discount_price*`, and `charge_unit` names.

2. **High-precision fields as TypeString**: The `*High` fields (`UnitPriceHigh`, `OriginalPriceHigh`, `DiscountPriceHigh`, `UnitPriceDiscountHigh`) use `*string` type in the SDK to preserve precision. Map these to `TypeString` in Terraform to avoid float precision loss, consistent with how other high-precision pricing fields are handled in the provider.

3. **Regular pricing fields as TypeFloat**: The non-High pricing fields (`UnitPrice`, `DiscountPrice`, `OriginalPrice`, `UnitPriceDiscount`) use `*float64` in the SDK. Map these to `TypeFloat` in Terraform, consistent with the existing `price` block in `instance_types`.

## Risks / Trade-offs

- [Risk] `Price` field may be nil in some `DiskConfig` responses → Mitigation: Check `diskConfig.Price != nil` before accessing pricing fields, consistent with the existing nil-check pattern for `instanceType.Price`
- [Risk] Individual pricing fields within `Price` may be nil → Mitigation: The Terraform SDK handles nil computed fields gracefully by omitting them from state
