## Why

The `tencentcloud_sms_template` resource currently does not expose the `ReviewReply` field from the `DescribeSmsTemplateList` API response. This field contains the reviewer's reply (typically the reason for rejection), which is essential for users to understand why their SMS template was not approved and take corrective action.

## What Changes

- Add a new computed (read-only) attribute `ReviewReply` to the `tencentcloud_sms_template` resource schema
- Read and set the `ReviewReply` value from `DescribeSmsTemplateList` API response in the resource Read method

## Capabilities

### New Capabilities
- `sms-template-review-reply`: Expose the ReviewReply field from DescribeSmsTemplateList API as a computed attribute on tencentcloud_sms_template resource

### Modified Capabilities

## Impact

- Resource file: `tencentcloud/services/sms/resource_tc_sms_template.go` — add schema field and Read logic
- Service file: No changes needed (DescribeSmsTemplate already returns the full struct including ReviewReply)
- Documentation: `tencentcloud/services/sms/resource_tc_sms_template.md` — add ReviewReply attribute
- Backward compatible: Adding a computed attribute does not break existing configurations or state
