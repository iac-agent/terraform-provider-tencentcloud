## ADDED Requirements

### Requirement: ReviewReply computed attribute
The `tencentcloud_sms_template` resource SHALL expose a computed attribute `review_reply` of type string that reflects the `ReviewReply` field from the `DescribeSmsTemplateList` API response (`response.DescribeTemplateStatusSet[].ReviewReply`).

The attribute SHALL be read-only (Computed: true) and SHALL NOT be settable by the user via Terraform configuration.

The Read method SHALL check if `template.ReviewReply != nil` before calling `d.Set("review_reply", template.ReviewReply)`, consistent with the existing nil-check pattern for other fields.

#### Scenario: ReviewReply is returned by the API
- **WHEN** the `DescribeSmsTemplateList` API returns a non-nil `ReviewReply` value for a template
- **THEN** the `review_reply` attribute in the Terraform state SHALL be set to that value

#### Scenario: ReviewReply is nil
- **WHEN** the `DescribeSmsTemplateList` API returns nil for `ReviewReply`
- **THEN** the `review_reply` attribute SHALL NOT be set (the nil check skips the Set call)

#### Scenario: ReviewReply is empty string
- **WHEN** the `DescribeSmsTemplateList` API returns an empty string for `ReviewReply`
- **THEN** the `review_reply` attribute SHALL be set to an empty string
