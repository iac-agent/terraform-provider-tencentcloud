## 1. Schema Extension

- [x] 1.1 Add `charge_unit` computed field (TypeString) to the `cbs_configs` nested schema in `data_source_tc_instance_types.go`
- [x] 1.2 Add `discount_price` computed field (TypeFloat) to the `cbs_configs` nested schema in `data_source_tc_instance_types.go`
- [x] 1.3 Add `discount_price_high` computed field (TypeString) to the `cbs_configs` nested schema in `data_source_tc_instance_types.go`
- [x] 1.4 Add `original_price` computed field (TypeFloat) to the `cbs_configs` nested schema in `data_source_tc_instance_types.go`
- [x] 1.5 Add `original_price_high` computed field (TypeString) to the `cbs_configs` nested schema in `data_source_tc_instance_types.go`
- [x] 1.6 Add `unit_price` computed field (TypeFloat) to the `cbs_configs` nested schema in `data_source_tc_instance_types.go`
- [x] 1.7 Add `unit_price_discount` computed field (TypeFloat) to the `cbs_configs` nested schema in `data_source_tc_instance_types.go`
- [x] 1.8 Add `unit_price_discount_high` computed field (TypeString) to the `cbs_configs` nested schema in `data_source_tc_instance_types.go`
- [x] 1.9 Add `unit_price_high` computed field (TypeString) to the `cbs_configs` nested schema in `data_source_tc_instance_types.go`

## 2. Data Mapping Implementation

- [x] 2.1 Update the `cbs_configs` mapping block in `dataSourceTencentCloudInstanceTypesRead` to check `diskConfig.Price != nil` before accessing Price fields
- [x] 2.2 Map `charge_unit` from `diskConfig.Price.ChargeUnit`
- [x] 2.3 Map `discount_price` from `diskConfig.Price.DiscountPrice`
- [x] 2.4 Map `discount_price_high` from `diskConfig.Price.DiscountPriceHigh`
- [x] 2.5 Map `original_price` from `diskConfig.Price.OriginalPrice`
- [x] 2.6 Map `original_price_high` from `diskConfig.Price.OriginalPriceHigh`
- [x] 2.7 Map `unit_price` from `diskConfig.Price.UnitPrice`
- [x] 2.8 Map `unit_price_discount` from `diskConfig.Price.UnitPriceDiscount`
- [x] 2.9 Map `unit_price_discount_high` from `diskConfig.Price.UnitPriceDiscountHigh`
- [x] 2.10 Map `unit_price_high` from `diskConfig.Price.UnitPriceHigh`

## 3. Documentation Updates

- [x] 3.1 Update `data_source_tc_instance_types.md` to add descriptions for the 9 new pricing fields in the `cbs_configs` section

## 4. Testing

- [x] 4.1 Add unit test cases in `data_source_tc_instance_types_cbs_price_test.go` to verify the new pricing fields are correctly mapped when Price is populated
- [x] 4.2 Add unit test case to verify Price fields are null/empty when Price is nil
- [x] 4.3 Run unit tests with `go test -gcflags=all=-l` to verify all tests pass

## Dependencies
- Step 2 (Data Mapping) depends on Step 1 (Schema Extension) completion
- Step 3 (Documentation) can be done in parallel with Step 2
- Step 4 (Testing) depends on Steps 1-2 completion
