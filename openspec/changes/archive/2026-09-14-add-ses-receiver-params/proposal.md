## Why

The `tencentcloud_ses_receiver` resource currently supports creating recipient groups with basic fields (`receivers_name`, `desc`) and recipient details (`data` with `email` and `template_data`). According to the requirements analysis, all specified parameters are already implemented in the existing resource code. The `InvalidCount` parameter mentioned in the requirements is not available in the vendor SDK (`ReceiverData` struct has no `InvalidCount` field), so it cannot be added to the Terraform provider.

## What Changes

After validating against the vendor SDK (`github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ses/v20201002`):

- **No new parameters to add**: All required input parameters (`receivers_name`, `desc`, `data`, `email`, `template_data`) are already implemented in `resource_tc_ses_receiver.go`
- **`InvalidCount` cannot be added**: The `ReceiverData` struct in the vendor SDK does not contain an `InvalidCount` field, making it infeasible to implement

## Capabilities

### New Capabilities
- (none - all parameters already implemented)

### Modified Capabilities
- (none - no spec-level behavior changes needed)

## Impact

- No code changes required for the resource, as all specified parameters are already implemented
- The existing `resource_tc_ses_receiver.go`, `service_tencentcloud_ses.go`, and related test files already cover the full CRUD lifecycle