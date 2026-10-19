## ADDED Requirements

### Requirement: Computed command_id attribute
The `tencentcloud_tat_invoker_config` resource SHALL expose the invoker's associated command ID as a computed string attribute named `command_id`.

#### Scenario: Read returns command_id from DescribeInvokers
- **WHEN** the resource Read method is called
- **THEN** it SHALL query the `DescribeInvokers` API and set `command_id` from `response.InvokerSet[0].CommandId`

#### Scenario: command_id is populated after create
- **WHEN** a `tencentcloud_tat_invoker_config` resource is created (enabled/disabled)
- **THEN** the `command_id` attribute SHALL be readable in the Terraform state after the subsequent Read

#### Scenario: command_id is computed only
- **WHEN** a user defines a `tencentcloud_tat_invoker_config` resource in HCL
- **THEN** the `command_id` attribute SHALL NOT be configurable (no user input), SHALL be populated from the API response