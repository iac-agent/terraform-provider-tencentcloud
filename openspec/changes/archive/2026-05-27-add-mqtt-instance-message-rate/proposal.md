## Why

The `tencentcloud_mqtt_instance` resource currently lacks support for the `MessageRate` parameter, which controls the per-client message send/receive rate limit (unit: messages/second). This parameter is available in the cloud API's ModifyInstance and DescribeInstance interfaces, but is not exposed in the Terraform resource. Users need this parameter to manage client-level rate limiting for their MQTT instances.

## What Changes

- Add `message_rate` (TypeInt, Optional+Computed) parameter to the `tencentcloud_mqtt_instance` resource schema
- In the Create function, set `message_rate` via ModifyInstance call after instance creation (same pattern as `automatic_activation` and `authorization_policy`)
- In the Read function, read `message_rate` from the DescribeInstance response
- In the Update function, include `message_rate` in the mutableArgs list and set it in the ModifyInstance request
- Update the resource documentation (.md file) with the new parameter

## Capabilities

### New Capabilities
- `mqtt-instance-message-rate`: Adds the `message_rate` parameter to the `tencentcloud_mqtt_instance` resource, allowing users to configure per-client message rate limits for MQTT instances via Terraform.

### Modified Capabilities

## Impact

- **Affected code**: `tencentcloud/services/mqtt/resource_tc_mqtt_instance.go` (schema, create, read, update functions)
- **Affected documentation**: Resource .md file for mqtt instance
- **API dependencies**: `ModifyInstance` (v20240516), `DescribeInstance` (v20240516) - both already used by the resource
- **Backward compatibility**: Fully backward compatible - new Optional+Computed parameter, no breaking changes
