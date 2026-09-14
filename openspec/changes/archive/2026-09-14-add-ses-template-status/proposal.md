## Why

The `tencentcloud_ses_template` resource currently does not expose the template review status (`TemplateStatus`) returned by the `GetEmailTemplate` API. Users need this information to know whether their email template has been approved, is pending review, or has been rejected. Adding this computed attribute provides better visibility into the template lifecycle.

## What Changes

- Add a new computed output attribute `template_status` (TypeInt) to the `tencentcloud_ses_template` resource schema
- Read `template_status` from `TemplateStatus` field in the `GetEmailTemplate` API response during the Read operation
- Update the resource documentation (`resource_tc_ses_template.md`) to reflect the new attribute

## Capabilities

### New Capabilities
- `ses-template-status`: Expose the email template review status as a computed output attribute

### Modified Capabilities

*(No existing capabilities are being modified)*

## Impact

- **Code**: Only `tencentcloud/services/ses/resource_tc_ses_template.go` needs modification - add schema field and populate it in the Read function
- **Documentation**: Update `tencentcloud/services/ses/resource_tc_ses_template.md` with the new attribute
- **Test**: Update unit tests in `resource_tc_ses_template_test.go` to verify the new field
- **Backward Compatibility**: Fully backward compatible - adding only a computed (read-only) field