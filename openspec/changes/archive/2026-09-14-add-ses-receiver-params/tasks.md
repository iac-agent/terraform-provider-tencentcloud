## 1. Schema Definition

- [x] 1.1 Add computed `count` attribute (TypeInt, Computed) to `tencentcloud_ses_receiver` resource schema in `tencentcloud/services/ses/resource_tc_ses_receiver.go`

## 2. Read Method Update

- [x] 2.1 In `resourceTencentCloudSesReceiverRead`, add nil check and set `count` from `receiver.Count` (the `DescribeSesReceiverById` result)

## 3. Unit Tests

- [x] 3.1 Update `tencentcloud/services/ses/resource_tc_ses_receiver_test.go` to add unit test covering the `count` computed attribute being read from the API response

## 4. Documentation

- [x] 4.1 Update `tencentcloud/services/ses/resource_tc_ses_receiver.md` to reflect the new `count` attribute in the example
