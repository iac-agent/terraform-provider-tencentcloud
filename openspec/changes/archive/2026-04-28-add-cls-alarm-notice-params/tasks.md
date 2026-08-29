## 1. Tag Format Conversion Helpers

- [x] 1.1 Add `mapToClsTags` helper function to convert `map[string]interface{}` (Terraform TypeMap) to `[]*cls.Tag` (CLS API format) in `resource_tc_cls_alarm_notice.go`
- [x] 1.2 Add `clsTagsToMap` helper function to convert `[]*cls.Tag` (CLS API format) to `map[string]interface{}` (Terraform TypeMap) in `resource_tc_cls_alarm_notice.go`

## 2. Create Logic - Tags Direct API Support

- [x] 2.1 In `resourceTencentCloudClsAlarmNoticeCreate`, before the API call, convert `tags` from TypeMap to `[]*cls.Tag` using `mapToClsTags` and set `CreateAlarmNoticeRequest.Tags`
- [x] 2.2 Verify that the existing tag service `ModifyTags` call after create is kept as fallback (idempotent)

## 3. Read Logic - Tags from AlarmNotice Response

- [x] 3.1 In `resourceTencentCloudClsAlarmNoticeRead`, add logic to read tags from `AlarmNotice.Tags` response field, converting using `clsTagsToMap` and setting via `d.Set("tags", ...)`
- [x] 3.2 Ensure that when `AlarmNotice.Tags` is nil, the existing tag service `DescribeResourceTags` is used as fallback

## 4. Update Logic - Tags Direct API Support

- [x] 4.1 In `resourceTencentCloudClsAlarmNoticeUpdate`, when `tags` has changed, convert the new tags from TypeMap to `[]*cls.Tag` using `mapToClsTags` and set `ModifyAlarmNoticeRequest.Tags`
- [x] 4.2 Add `tags` to the `mutableArgs` list if not already present
- [x] 4.3 Ensure the existing tag service `ModifyTags` call in update is kept as fallback

## 5. Unit Tests

- [x] 5.1 Add unit test for `mapToClsTags` function to verify TypeMap to `[]*cls.Tag` conversion
- [x] 5.2 Add unit test for `clsTagsToMap` function to verify `[]*cls.Tag` to TypeMap conversion
- [x] 5.3 Add unit test for Create with tags to verify `CreateAlarmNoticeRequest.Tags` is set correctly (integration test with mapToClsTags)
- [x] 5.4 Add unit test for Read with tags from AlarmNotice response to verify `clsTagsToMap` conversion is correct
- [x] 5.5 Add unit test for Update with tags to verify `ModifyAlarmNoticeRequest.Tags` is set correctly (integration test with mapToClsTags)
- [x] 5.6 Run unit tests with `go test -v -gcflags=all=-l ./tencentcloud/services/cls/...` to verify all tests pass

## 6. Documentation

- [x] 6.1 Update `tencentcloud/services/cls/resource_tc_cls_alarm_notice.md` if needed to reflect tags parameter behavior (already contains tags example, no changes needed)
- [ ] 6.2 Run `make doc` to generate website documentation (only in finalize phase)
