## Why

The `tencentcloud_ses_receiver` resource is missing the `count` field, which represents the total number of recipient email addresses in a receiver list. This field is returned by the `ListReceivers` API (`response.Response.Data.Count`) but is not currently exposed in the Terraform resource schema. Users cannot read or reference the recipient count of a receiver list through Terraform.

## What Changes

- Add a new computed attribute `count` (type: integer) to the `tencentcloud_ses_receiver` resource schema
- Set the `count` value in the resource Read function when the receiver data is fetched from the `ListReceivers` API response

## Capabilities

### New Capabilities
- `ses-receiver-count`: Expose the recipient count field from the ListReceivers API response as a computed attribute in the tencentcloud_ses_receiver resource

### Modified Capabilities

## Impact

- `tencentcloud/services/ses/resource_tc_ses_receiver.go`: Add `count` to schema and set it in Read function
- `tencentcloud/services/ses/resource_tc_ses_receiver_test.go`: Add test coverage for the new `count` attribute
- `tencentcloud/services/ses/resource_tc_ses_receiver.md`: Update documentation to reflect the new `count` attribute
