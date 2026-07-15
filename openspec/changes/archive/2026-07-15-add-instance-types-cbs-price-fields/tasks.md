## 1. Schema Extension

- [x] 1.1 Add `charge_unit` computed field (TypeString) to `cbs_configs` nested schema in `data_source_tc_instance_types.go`
- [x] 1.2 Add `unit_price` computed field (TypeFloat) to `cbs_configs` nested schema
- [x] 1.3 Add `unit_price_discount` computed field (TypeFloat) to `cbs_configs` nested schema
- [x] 1.4 Add `unit_price_high` computed field (TypeString) to `cbs_configs` nested schema
- [x] 1.5 Add `unit_price_discount_high` computed field (TypeString) to `cbs_configs` nested schema
- [x] 1.6 Add `original_price` computed field (TypeFloat) to `cbs_configs` nested schema
- [x] 1.7 Add `original_price_high` computed field (TypeString) to `cbs_configs` nested schema
- [x] 1.8 Add `discount_price` computed field (TypeFloat) to `cbs_configs` nested schema
- [x] 1.9 Add `discount_price_high` computed field (TypeString) to `cbs_configs` nested schema

## 2. Data Mapping Implementation

- [x] 2.1 Add nil-check for `diskConfig.Price` before accessing pricing fields in the `cbs_configs` mapping loop
- [x] 2.2 Map `charge_unit` from `diskConfig.Price.ChargeUnit` in the cbs_configs mapping
- [x] 2.3 Map `unit_price` from `diskConfig.Price.UnitPrice` in the cbs_configs mapping
- [x] 2.4 Map `unit_price_discount` from `diskConfig.Price.UnitPriceDiscount` in the cbs_configs mapping
- [x] 2.5 Map `unit_price_high` from `diskConfig.Price.UnitPriceHigh` in the cbs_configs mapping
- [x] 2.6 Map `unit_price_discount_high` from `diskConfig.Price.UnitPriceDiscountHigh` in the cbs_configs mapping
- [x] 2.7 Map `original_price` from `diskConfig.Price.OriginalPrice` in the cbs_configs mapping
- [x] 2.8 Map `original_price_high` from `diskConfig.Price.OriginalPriceHigh` in the cbs_configs mapping
- [x] 2.9 Map `discount_price` from `diskConfig.Price.DiscountPrice` in the cbs_configs mapping
- [x] 2.10 Map `discount_price_high` from `diskConfig.Price.DiscountPriceHigh` in the cbs_configs mapping

## 3. Documentation Updates

- [x] 3.1 Update `data_source_tc_instance_types.md` to document the 9 new pricing fields in the `cbs_configs` section with descriptions and example usage

## 4. Testing

- [x] 4.1 Add unit test case in `data_source_tc_instance_types_test.go` verifying pricing fields are correctly mapped when Price is non-nil
- [x] 4.2 Add unit test case verifying cbs_configs mapping handles nil Price gracefully without errors
