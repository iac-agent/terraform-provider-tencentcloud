## 1. Schema & Schema Registration

- [x] 1.1 Add `total_count` computed field to the `tencentcloud_teo_zone` resource schema in `tencentcloud/services/teo/resource_tc_teo_zone.go`

## 2. Service Layer

- [x] 2.1 Modify `DescribeTeoZoneById` in `tencentcloud/services/teo/service_tencentcloud_teo.go` to also return `TotalCount` (`*int64`) from the API response

## 3. Read Function

- [x] 3.1 Update `resourceTencentCloudTeoZoneRead` in `tencentcloud/services/teo/resource_tc_teo_zone.go` to capture and set `total_count` from the modified service call

## 4. Documentation

- [x] 4.1 Update `tencentcloud/services/teo/resource_tc_teo_zone.md` to include `total_count` in the example usage

## 5. Testing

- [x] 5.1 Add unit test case for `total_count` field in `tencentcloud/services/teo/resource_tc_teo_zone_test.go`

## 6. Verification

- [x] 6.1 Verify all Go files compile successfully in the `teo` package (code correctness check)
