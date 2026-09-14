## Context

The `tencentcloud_ses_receiver` resource resides at `tencentcloud/services/ses/resource_tc_ses_receiver.go`. It currently supports full CRUD lifecycle for SES recipient groups using the following cloud APIs:

- **CreateReceiver**: Creates a recipient group with `ReceiversName` and `Desc`
- **CreateReceiverDetailWithData**: Adds recipients with template data (`ReceiverInputData` containing `Email` and `TemplateData`)
- **CreateReceiverDetail**: Adds recipients by email list (`Emails []*string`)
- **ListReceivers**: Queries recipient groups, returning `ReceiverData` with fields: `ReceiverId`, `ReceiversName`, `Count`, `Desc`, `ReceiversStatus`, `CreateTime`
- **ListReceiverDetails**: Queries recipient details, returning `ReceiverDetail` with fields: `Email`, `CreateTime`, `TemplateData`
- **DeleteReceiver**: Deletes a recipient group

The vendor SDK (`github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ses/v20201002`) has been verified - the `ReceiverData` struct does **not** contain an `InvalidCount` field.

## Goals / Non-Goals

**Goals:**
- Validate whether new parameters need to be added to `tencentcloud_ses_receiver`

**Non-Goals:**
- No actual code changes required since all specified parameters are already implemented

## Decisions

| Decision | Rationale |
|---|---|
| **No new parameters needed** | All specified parameters (`receivers_name`, `desc`, `data`, `email`, `template_data`) are already present in the existing resource schema |
| **`InvalidCount` cannot be added** | The `ReceiverData` struct in the vendor SDK has no `InvalidCount` field; the cloud API `ListReceivers` response does not expose this field |

## Risks / Trade-offs

- **No risk**: Since no code changes are needed, there are no deployment, migration, or rollback concerns