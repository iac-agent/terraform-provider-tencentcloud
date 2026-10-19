## 1. Schema and Resource Code Changes

- [x] 1.1 Add `count` computed attribute (TypeInt, Computed: true) to the `tencentcloud_ses_receiver` resource schema in `tencentcloud/services/ses/resource_tc_ses_receiver.go`
- [x] 1.2 Set `count` value from `receiver.Count` (with nil check) in `resourceTencentCloudSesReceiverRead` in `tencentcloud/services/ses/resource_tc_ses_receiver.go`

## 2. Test Updates

- [x] 2.1 Add unit test coverage for the `count` attribute in `tencentcloud/services/ses/resource_tc_ses_receiver_test.go`, verifying the Read function sets count from API response
- [x] 2.2 Add `resource.TestCheckResourceAttrSet("tencentcloud_ses_receiver.receiver", "count")` to the existing acceptance test check functions in `tencentcloud/services/ses/resource_tc_ses_receiver_test.go`

## 3. Documentation Updates

- [x] 3.1 Update `tencentcloud/services/ses/resource_tc_ses_receiver.md` to reflect the new `count` computed attribute
