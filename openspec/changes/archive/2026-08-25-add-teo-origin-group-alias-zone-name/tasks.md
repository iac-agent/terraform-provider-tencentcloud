## 1. Schema and Read Logic

- [x] 1.1 Add `alias_zone_name` computed string field to the `references` nested block schema in `ResourceTencentCloudTeoOriginGroup()` in `tencentcloud/services/teo/resource_tc_teo_origin_group.go`
- [x] 1.2 Add the nil-check and field-setting logic for `alias_zone_name` inside the `references` loop within `resourceTencentCloudTeoOriginGroupRead()` in `tencentcloud/services/teo/resource_tc_teo_origin_group.go` (`if references.AliasZoneName != nil { referencesMap["alias_zone_name"] = references.AliasZoneName }`)

## 2. Unit Tests

- [x] 2.1 Add gomonkey-based mock unit tests for the new `alias_zone_name` computed field in `tencentcloud/services/teo/resource_tc_teo_origin_group_test.go`, covering the case where `AliasZoneName` is returned by the API and the case where it is nil

## 3. Documentation

- [x] 3.1 Update `tencentcloud/services/teo/resource_tc_teo_origin_group.md` to reflect the new `alias_zone_name` computed field in the `references` block (the `website/docs/` markdown is generated via `make doc` during finalization)

## 4. Verification

- [x] 4.1 Verify the generated Go code compiles and the resource schema, read logic, and tests are consistent with the cloud API `OriginGroupReference.AliasZoneName` field