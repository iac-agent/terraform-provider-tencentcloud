## 1. Schema Changes

- [x] 1.1 Add `allow_duplicates` (TypeBool, Optional, ForceNew) to `tencentcloud_teo_zone` resource schema in `resource_tc_teo_zone.go`
- [x] 1.2 Add `jump_start` (TypeBool, Optional, ForceNew) to `tencentcloud_teo_zone` resource schema in `resource_tc_teo_zone.go`
- [x] 1.3 Add `file_verification` sub-block (TypeList, Computed, MaxItems: 1) with `path` and `content` fields under `ownership_verification` in schema
- [x] 1.4 Add `ns_verification` sub-block (TypeList, Computed, MaxItems: 1) with `name_servers` (TypeList of TypeString) field under `ownership_verification` in schema

## 2. Create Method Updates

- [x] 2.1 Read `allow_duplicates` from schema and set `request.AllowDuplicates` in `resourceTencentCloudTeoZoneCreate`
- [x] 2.2 Read `jump_start` from schema and set `request.JumpStart` in `resourceTencentCloudTeoZoneCreate`

## 3. Read Method Updates

- [x] 3.1 Populate `file_verification` sub-block from `respData.OwnershipVerification.FileVerification` in `resourceTencentCloudTeoZoneRead`
- [x] 3.2 Populate `ns_verification` sub-block from `respData.OwnershipVerification.NsVerification` in `resourceTencentCloudTeoZoneRead`
- [x] 3.3 Handle nil cases for `FileVerification` and `NsVerification` (set empty list when nil)

## 4. Tests

- [x] 4.1 Add unit test cases for `allow_duplicates` and `jump_start` schema fields in `resource_tc_teo_zone_test.go`
- [x] 4.2 Add unit test cases for `file_verification` and `ns_verification` in `ownership_verification` block in tests

## 5. Documentation

- [x] 5.1 Update `resource_tc_teo_zone.md` with new parameter descriptions and usage examples
- [x] 5.2 Run `make doc` to regenerate `website/docs/` documentation (handled by tfpacer-finalize skill)

## 6. Verification

- [x] 6.1 Verify code compiles without errors
- [x] 6.2 Verify all new schema fields are properly set in state during Read
- [x] 6.3 Verify backward compatibility (existing configs continue to work without the new fields)