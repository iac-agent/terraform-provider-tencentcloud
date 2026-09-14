## Why

The `tencentcloud_sms_sign` resource currently does not expose the `QualificationStatusCode` field from the `DescribeSmsSignList` API response. Users cannot observe the qualification status of their SMS sign through Terraform, which is important for monitoring whether the domestic SMS qualification is pending review, approved, rejected, or needs supplemental submission.

## What Changes

- Add a new computed attribute `qualification_status_code` (TypeInt, Computed) to the `tencentcloud_sms_sign` resource schema
- Read and set the `QualificationStatusCode` field from `DescribeSignListStatus` in the resource Read method

## Capabilities

### New Capabilities
- `sms-sign-qualification-status-code`: Expose the QualificationStatusCode field from the DescribeSmsSignList API response as a computed attribute on the tencentcloud_sms_sign resource, allowing users to track the domestic SMS qualification review status (0=pending, 1=approved, 2=rejected, 3=supplement required, 4=modified pending, 5=modified rejected)

### Modified Capabilities

## Impact

- `tencentcloud/services/sms/resource_tc_sms_sign.go`: Add schema field and read logic
- `tencentcloud/services/sms/resource_tc_sms_sign_test.go`: Add unit test for the new field
- `tencentcloud/services/sms/resource_tc_sms_sign.md`: Update documentation
