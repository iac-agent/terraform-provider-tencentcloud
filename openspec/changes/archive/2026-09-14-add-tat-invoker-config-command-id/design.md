## Context

The `tencentcloud_tat_invoker_config` resource currently manages the on/off state of a TAT invoker. It uses:
- `EnableInvoker` / `DisableInvoker` APIs for state changes (Update)
- `DescribeInvokers` API for reading the current state (Read)
- The invoker ID is set as the resource ID

The `Invoker` struct returned by `DescribeInvokers` contains a `CommandId` field (`*string`), which is already available in the API response but not exposed to Terraform users. This field is read-only (computed) and provides the command ID associated with the invoker.

## Goals / Non-Goals

**Goals:**
- Expose `command_id` as a computed attribute on the `tencentcloud_tat_invoker_config` resource
- Populate `command_id` from `DescribeInvokers` response's `InvokerSet[0].CommandId` in the Read method
- Maintain full backward compatibility with existing Terraform configurations and state

**Non-Goals:**
- Making `command_id` a configurable (user-set) attribute — it is read-only from the API
- Modifying the Create/Update/Delete business logic
- Adding any other new parameters beyond `command_id`

## Decisions

1. **Schema type**: `schema.TypeString` with `Computed: true` — a computed string attribute that is populated from the API response, never set by the user.
2. **Read logic**: In `resourceTencentCloudTatInvokerConfigRead`, after reading `invokerConfig.CommandId`, check for nil and set via `d.Set("command_id", invokerConfig.CommandId)`. This follows the existing pattern used for `invoker_id`.
3. **Documentation**: Update the `.md` file to mention the new `command_id` attribute in the example output.
4. **No schema changes to other fields**: The existing `invoker_id` (Required) and `invoker_status` (Required) remain unchanged.

## Risks / Trade-offs

- **Minimal risk**: This is a purely additive change — a new computed field that is read from the API response. No existing behavior is modified.
- **Backward compatibility**: Existing `terraform plan/apply` workflows are unaffected. Existing state files will not contain `command_id` until the next `terraform refresh` or `terraform apply` that triggers a Read.
- **SDK availability**: The `CommandId` field already exists in the vendor SDK's `Invoker` struct, so no dependency update is needed.