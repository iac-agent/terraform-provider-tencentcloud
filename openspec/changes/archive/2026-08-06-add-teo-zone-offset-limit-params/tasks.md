## 1. Schema Changes

- [x] 1.1 Add `offset` (Optional, TypeInt, Description: "Pagination offset for the DescribeZones API call during the Read operation. Default: 0.") to the `tencentcloud_teo_zone` resource schema in `tencentcloud/services/teo/resource_tc_teo_zone.go`
- [x] 1.2 Add `limit` (Optional, TypeInt, Description: "Pagination limit for the DescribeZones API call during the Read operation. Default: 20, maximum: 100.") to the `tencentcloud_teo_zone` resource schema in `tencentcloud/services/teo/resource_tc_teo_zone.go`

## 2. Service Layer Changes

- [x] 2.1 Modify `DescribeTeoZoneById` method signature in `tencentcloud/services/teo/service_tencentcloud_teo.go` to accept `offset` and `limit` parameters (`int64`)
- [x] 2.2 Update the method body to use the provided offset/limit values when they are non-zero, falling back to the existing defaults (`offset=0`, `limit=20`) when zero

## 3. Read Method Update

- [x] 3.1 Update the `resourceTencentCloudTeoZoneRead` function in `tencentcloud/services/teo/resource_tc_teo_zone.go` to read `offset` and `limit` from the schema and pass them to `DescribeTeoZoneById`

## 4. Documentation

- [x] 4.1 Update `tencentcloud/services/teo/resource_tc_teo_zone.md` with example usage for the new `offset` and `limit` parameters

## 5. Unit Tests

- [x] 5.1 Add unit test cases in `tencentcloud/services/teo/resource_tc_teo_zone_offset_limit_test.go` for the new `offset` and `limit` parameters