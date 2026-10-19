## Context

The `tencentcloud_ses_receiver` resource manages SES (Simple Email Service) recipient groups in TencentCloud. The resource already implements Create, Read, and Delete operations with `receivers_name`, `desc`, and `data` (containing `email` and `template_data`) parameters. The `DescribeSesReceiverById` service method calls the `ListReceivers` API and returns a `*ses.ReceiverData` struct which contains a `Count` field (`*uint64`) representing the total number of recipient email addresses. This `Count` field is not currently exposed in the Terraform resource schema.

## Goals / Non-Goals

**Goals:**
- Add a computed `count` attribute to the `tencentcloud_ses_receiver` resource schema, mapping to `ReceiverData.Count` from the `ListReceivers` API response
- Populate the `count` field in the resource Read method
- Update unit tests to cover the new attribute
- Update the resource documentation (.md file)

**Non-Goals:**
- Do not modify existing schema fields (`receivers_name`, `desc`, `data`) or their behavior
- Do not add Update support (resource remains ForceNew for all fields)
- Do not modify the `tencentcloud_ses_receivers` datasource

## Decisions

1. **Schema attribute type for `count`**: Use `schema.TypeInt` with `Computed: true` since `ReceiverData.Count` is `*uint64` in the SDK. This follows the existing pattern in the codebase for integer computed attributes.

2. **Where to read `count`**: The `count` value is already available in the `DescribeSesReceiverById` service method result (`receiver.Count`). No additional API call is needed — just set it in the existing Read method alongside `receivers_name` and `desc`.

3. **Attribute naming**: Use `count` (lowercase) as the Terraform schema name, matching the `ReceiverData.Count` JSON field. This follows the Terraform naming convention for snake_case attributes.

## Risks / Trade-offs

- [Backward compatibility] → Adding a computed attribute is backward compatible; existing Terraform configurations and states will not be affected.
- [API response nil check] → The `Count` field in `ReceiverData` could be nil (the SDK marks it as `omitnil`). The Read method must check for nil before setting the value.
