## 1. Schema & Resource Code

- [x] 1.1 Add `qualification_status_code` computed attribute (TypeInt, Computed) to the `tencentcloud_sms_sign` resource schema in `tencentcloud/services/sms/resource_tc_sms_sign.go`
- [x] 1.2 Add nil-check and `d.Set("qualification_status_code", sign.QualificationStatusCode)` in `resourceTencentCloudSmsSignRead` after the existing `international` field set, following the existing pattern

## 2. Unit Tests

- [x] 2.1 Add unit test in `tencentcloud/services/sms/resource_tc_sms_sign_test.go` using gomonkey mock to verify that `resourceTencentCloudSmsSignRead` correctly sets `qualification_status_code` when the API returns a non-nil value
- [x] 2.2 Add unit test to verify that `resourceTencentCloudSmsSignRead` does not error when `QualificationStatusCode` is nil (nil-check scenario)

## 3. Documentation

- [x] 3.1 Update `tencentcloud/services/sms/resource_tc_sms_sign.md` to include the new `qualification_status_code` attribute in the example output
