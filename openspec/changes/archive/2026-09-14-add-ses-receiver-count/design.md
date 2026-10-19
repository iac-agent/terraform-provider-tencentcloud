## Context

The `tencentcloud_ses_receiver` resource manages SES recipient groups. It currently exposes `receivers_name`, `desc`, and `data` (with nested `email` and `template_data`) attributes. The `ListReceivers` API returns a `Count` field (`*uint64`) on each `ReceiverData` object, indicating the total number of recipient email addresses in the list. This field is not currently surfaced in the Terraform resource schema.

The `Count` field is already available from the `DescribeSesReceiverById` service method (which calls `ListReceivers` and returns `*ses.ReceiverData`), but it is never set on the `schema.ResourceData`.

## Goals / Non-Goals

**Goals:**
- Expose the `count` attribute as a computed integer field in the `tencentcloud_ses_receiver` resource schema
- Set the `count` value in the Read function from the `ListReceivers` API response

**Non-Goals:**
- Modifying the Create, Update, or Delete flows (count is read-only)
- Adding the `count` field to any data source (not in scope)
- Changing any existing schema fields or behavior

## Decisions

1. **Schema type for `count`**: Use `schema.TypeInt` with `Computed: true`. The API returns `*uint64`, and Terraform's TypeInt is the standard mapping for integer fields. Computed ensures it is read-only and never set by the user.

2. **Where to set `count`**: In `resourceTencentCloudSesReceiverRead`, after calling `service.DescribeSesReceiverById`, set `count` from `receiver.Count` (with nil check). This follows the existing pattern for `receivers_name` and `desc`.

3. **No update to service layer**: The `DescribeSesReceiverById` method already returns the full `*ses.ReceiverData` struct which includes `Count`, so no service-layer changes are needed.

## Risks / Trade-offs

- **[Backward compatibility]** Adding a computed field is backward compatible — existing Terraform configurations and state files are unaffected. → No mitigation needed.
- **[API field may be nil]** The `Count` field may be nil in some edge cases (API documentation notes it can be null for some fields). → Mitigation: check for nil before calling `d.Set()`.
