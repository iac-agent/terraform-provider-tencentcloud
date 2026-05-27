## Context

The `tencentcloud_mqtt_instance` resource manages MQTT broker instances on TencentCloud. Currently, the resource supports parameters like `instance_type`, `name`, `sku_code`, `remark`, `vpc_list`, `renew_flag`, `time_span`, `pay_mode`, `automatic_activation`, and `authorization_policy`.

The cloud API's ModifyInstance interface supports a `MessageRate` parameter (int64) that controls per-client message send/receive rate limits (unit: messages/second). This parameter is also returned in the DescribeInstance response. However, it is not available in the CreateInstance interface, meaning it can only be configured after the instance is created.

The existing code already has a pattern for handling parameters that only exist in ModifyInstance (not in CreateInstance): `automatic_activation` and `authorization_policy` are set via a ModifyInstance call after the instance is created and reaches RUNNING status.

## Goals / Non-Goals

**Goals:**
- Add `message_rate` parameter to the `tencentcloud_mqtt_instance` resource schema
- Enable users to configure per-client message rate limits via Terraform
- Follow the existing pattern for post-creation parameter configuration (similar to `automatic_activation` and `authorization_policy`)
- Maintain full backward compatibility

**Non-Goals:**
- Adding other missing parameters (EnablePublic, Bandwidth, IpRules, X509Mode, UseDefaultServerCert) - these are out of scope for this change
- Changing the existing behavior of any current parameters
- Modifying the `tencentcloud_mqtt_instance_public_endpoint` resource

## Decisions

### Decision 1: Schema definition for `message_rate`
- **Choice**: TypeInt, Optional+Computed
- **Rationale**: The parameter is optional (not required for creation) and should be Computed to allow the cloud API to return its value. This matches the pattern of other mutable parameters like `automatic_activation`.
- **Alternatives**: Making it Required would break existing configurations. Making it only Optional without Computed would cause state drift if the cloud API returns a value that wasn't explicitly set.

### Decision 2: Setting `message_rate` during Create
- **Choice**: Add `message_rate` to the existing post-creation ModifyInstance call in the Create function (alongside `automatic_activation` and `authorization_policy`)
- **Rationale**: The CreateInstance API does not support `MessageRate`, so it must be set via ModifyInstance after the instance reaches RUNNING status. The existing code already has a post-creation ModifyInstance call for `automatic_activation` and `authorization_policy` - we extend this to include `message_rate`.
- **Alternatives**: A separate ModifyInstance call just for `message_rate` would add unnecessary API calls and complexity.

### Decision 3: Setting `message_rate` during Update
- **Choice**: Add `message_rate` to the `mutableArgs` list and include it in the ModifyInstance request during update
- **Rationale**: The ModifyInstance API supports `MessageRate`, so it can be updated. Adding it to `mutableArgs` ensures it triggers the update flow, and setting it in the request ensures the API call includes it.

### Decision 4: Reading `message_rate` during Read
- **Choice**: Read `MessageRate` from the DescribeInstance response in the Read function
- **Rationale**: The DescribeInstance response includes `MessageRate`, so it should be read back and stored in state. Following the existing pattern, check for nil before setting.

## Risks / Trade-offs

- **Risk**: The ModifyInstance call after creation may fail, leaving the instance in a partially configured state → Mitigation: This is already the existing pattern for `automatic_activation` and `authorization_policy`, and the Terraform state will reflect the actual value on the next read/plan.
- **Risk**: If the cloud API returns 0 for `MessageRate` when not explicitly set, it may be unclear whether the user set it to 0 or it's a default → Mitigation: Using Optional+Computed ensures the schema handles this correctly; the cloud API documentation clarifies the default behavior.
