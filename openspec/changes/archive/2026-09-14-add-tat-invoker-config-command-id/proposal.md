## Why

The `tencentcloud_tat_invoker_config` resource currently manages the enable/disable state of a TAT invoker but does not expose the associated command ID (`command_id`) as a readable attribute. Users need to call the `DescribeInvokers` API separately or reference the `tencentcloud_tat_invoker` resource to obtain the invoker's command ID. Adding a computed `command_id` field to this resource provides users with a more complete view of the invoker configuration directly from the resource, improving usability and consistency.

## What Changes

- Add a new computed field `command_id` (TypeString, Computed: true) to the `tencentcloud_tat_invoker_config` resource schema, populated from the `DescribeInvokers` API response (`InvokerSet.CommandId`).
- Update the Read method to read and set `command_id` from the queried invoker data.
- Update documentation and tests accordingly.

## Capabilities

### New Capabilities

- `tat-invoker-config-command-id`: Expose the invoker's associated command ID as a computed attribute on the `tencentcloud_tat_invoker_config` resource, readable after creation or import.

### Modified Capabilities

<!-- No existing capabilities are modified; this is purely additive. -->

## Impact

- **Code**: `tencentcloud/services/tat/resource_tc_tat_invoker_config.go` – add `command_id` to Schema and Read logic.
- **Documentation**: `tencentcloud/services/tat/resource_tc_tat_invoker_config.md` – update example and attribute reference.
- **Tests**: `resource_tc_tat_invoker_config_test.go` – add validation for the new `command_id` attribute.
- **Provider registration**: No change needed (resource already registered).
- **SDK**: No change needed (`CommandId` already exists in `Invoker` struct).