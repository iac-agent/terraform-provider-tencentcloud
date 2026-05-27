## ADDED Requirements

### Requirement: message_rate parameter schema
The `tencentcloud_mqtt_instance` resource SHALL support a `message_rate` parameter of TypeInt, Optional and Computed, that controls the per-client message send/receive rate limit in messages/second.

#### Scenario: message_rate is optional
- **WHEN** a user creates an MQTT instance without specifying `message_rate`
- **THEN** the resource SHALL be created successfully and the `message_rate` value SHALL be read from the DescribeInstance response

#### Scenario: message_rate is set during creation
- **WHEN** a user creates an MQTT instance with `message_rate` specified
- **THEN** after the instance reaches RUNNING status, the resource SHALL call ModifyInstance to set the `message_rate` value
- **AND** the `message_rate` SHALL be included in the same ModifyInstance call as `automatic_activation` and `authorization_policy` if they are also set

#### Scenario: message_rate is read from API
- **WHEN** the resource Read function is called
- **THEN** the `message_rate` value SHALL be read from the DescribeInstance response's `MessageRate` field
- **AND** the value SHALL only be set in state if `MessageRate` is not nil in the response

### Requirement: message_rate update support
The `tencentcloud_mqtt_instance` resource SHALL support updating the `message_rate` parameter.

#### Scenario: message_rate is updated
- **WHEN** a user changes the `message_rate` value in the Terraform configuration
- **THEN** the resource SHALL detect the change via the `mutableArgs` check
- **AND** the resource SHALL include `message_rate` in the ModifyInstance request with the new value
- **AND** the ModifyInstance call SHALL use `tccommon.WriteRetryTimeout` with retry logic

#### Scenario: message_rate is unchanged
- **WHEN** a user updates other parameters but not `message_rate`
- **THEN** the `message_rate` SHALL NOT be included in the ModifyInstance request if its value has not changed

### Requirement: Backward compatibility
The addition of `message_rate` SHALL NOT break existing Terraform configurations.

#### Scenario: Existing configuration without message_rate
- **WHEN** a user applies an existing Terraform configuration that does not include `message_rate`
- **THEN** the resource SHALL continue to work without any changes to behavior
- **AND** the `message_rate` value SHALL be computed from the API response

### Requirement: message_rate documentation
The resource .md documentation file SHALL include `message_rate` in the example usage.

#### Scenario: Documentation example includes message_rate
- **WHEN** a user reads the resource documentation
- **THEN** the example SHALL show `message_rate` as an optional parameter
