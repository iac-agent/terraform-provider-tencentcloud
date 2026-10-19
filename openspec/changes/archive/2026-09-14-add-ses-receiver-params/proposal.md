## Why

The `tencentcloud_ses_receiver` resource currently already defines `receivers_name`, `desc`, and `data` (with nested `email` and `template_data`) parameters in its schema and create logic, but the Read method does not populate the `data` block from the `ListReceiverDetails` API response for entries that only have `Email` (without `TemplateData`). Additionally, the `ListReceivers` response includes a `Count` field (`ReceiverData.Count`) that is not exposed as a computed attribute in the resource schema. This change adds the missing `count` computed attribute to fully reflect the cloud API response data in the Terraform resource state.

## What Changes

- Add a new computed attribute `count` (type: int) to the `tencentcloud_ses_receiver` resource schema, mapping to `ReceiverData.Count` from the `ListReceivers` API response
- Populate the `count` field in the Read method by reading it from the `ListReceivers` API response

## Capabilities

### New Capabilities

- `ses-receiver-count-attr`: Expose the `count` (recipient count) computed attribute from the ListReceivers API response in the tencentcloud_ses_receiver resource

### Modified Capabilities

(None - no existing spec-level behavior changes)

## Impact

- `tencentcloud/services/ses/resource_tc_ses_receiver.go`: Schema definition and Read method
- `tencentcloud/services/ses/resource_tc_ses_receiver.md`: Documentation update
- `tencentcloud/services/ses/resource_tc_ses_receiver_test.go`: Unit test update
