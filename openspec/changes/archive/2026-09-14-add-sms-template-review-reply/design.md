## Context

The `tencentcloud_sms_template` resource manages SMS templates via the TencentCloud SMS API. The resource currently exposes `template_name`, `template_content`, `international`, and `sms_type` as writable fields, and reads back `template_name`, `template_content`, and `international` from the `DescribeSmsTemplateList` API.

The `DescribeSmsTemplateList` API returns a `DescribeTemplateListStatus` struct that includes a `ReviewReply` field (`*string`), which contains the reviewer's reply (typically the reason for template rejection). This field is already present in the SDK vendor code at `vendor/github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111/models.go` but is not currently exposed in the Terraform resource schema.

The service layer `DescribeSmsTemplate` method already returns the full `DescribeTemplateListStatus` struct, so no service-layer changes are needed.

## Goals / Non-Goals

**Goals:**
- Expose `ReviewReply` as a computed (read-only) attribute on `tencentcloud_sms_template`
- Read and set the `ReviewReply` value in the resource Read method
- Update the resource documentation (.md file)

**Non-Goals:**
- Modifying the Create, Update, or Delete operations
- Adding any other fields from `DescribeTemplateListStatus` (e.g., `StatusCode`, `CreateTime`)
- Changing the service layer code

## Decisions

1. **Schema type: Computed string** — `ReviewReply` is a read-only field returned by the API and cannot be set by the user. It SHALL be defined as `TypeString` with `Computed: true` in the resource schema.

2. **Nil check before Set** — Following the existing code pattern, the Read method SHALL check if `template.ReviewReply != nil` before calling `d.Set("review_reply", ...)`. This avoids potential nil pointer dereference and is consistent with the current handling of other fields.

3. **No service layer changes** — The `DescribeSmsTemplate` service method already returns the full `DescribeTemplateListStatus` struct which includes `ReviewReply`. No modifications to `service_tencentcloud_sms.go` are required.

## Risks / Trade-offs

- [Risk] `ReviewReply` may be empty string or nil for approved templates → Mitigation: Nil check ensures no crash; empty string is a valid computed value in Terraform.
- [Risk] Backward compatibility → Mitigation: Adding a computed attribute does not affect existing configurations or state files. Fully backward compatible.
