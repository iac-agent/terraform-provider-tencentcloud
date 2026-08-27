## 1. Schema Changes

- [x] 1.1 Add `offset` (Optional, TypeInt) to `tencentcloud_teo_zone` resource schema in `resource_tc_teo_zone.go`
- [x] 1.2 Add `limit` (Optional, TypeInt) to `tencentcloud_teo_zone` resource schema in `resource_tc_teo_zone.go`
- [x] 1.3 Add `total_count` (Computed, TypeInt) to `tencentcloud_teo_zone` resource schema in `resource_tc_teo_zone.go`

## 2. Read Method Update

- [x] 2.1 Modify `resourceTencentCloudTeoZoneRead` to set `total_count` from the DescribeZones API response `TotalCount` field passed through the service layer

## 3. Service Layer Update

- [x] 3.1 Modify `DescribeTeoZoneById` to return `TotalCount` from the API response alongside the zone data

## 4. Unit Tests

- [x] 4.1 Add test cases in `resource_tc_teo_zone_test.go` to validate `offset`, `limit`, and `total_count` parameters

## 5. Documentation

- [x] 5.1 Update `resource_tc_teo_zone.md` with new parameters in the example usage section

## 6. Verification

- [x] 6.1 Run `gofmt` on all modified Go files (deferred to tfpacer-finalize skill)
- [x] 6.2 Run `make doc` to regenerate website documentation (deferred to tfpacer-finalize skill)