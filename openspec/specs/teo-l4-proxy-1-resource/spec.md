## ADDED Requirements

### Requirement: Resource registration
The system SHALL register `tencentcloud_teo_l4_proxy_1` as a Terraform resource in the TEO provider, exposing all CRUD operations via the TEO v20220901 SDK.

#### Scenario: Resource is registered in provider
- **WHEN** Terraform initializes the tencentcloud provider
- **THEN** `tencentcloud_teo_l4_proxy_1` is available as a managed resource type

### Requirement: Create L4 proxy instance
The system SHALL support creating a TEO L4 proxy instance via `CreateL4Proxy` API with required parameters `zone_id`, `proxy_name`, `area` and optional parameters `ipv6`, `static_ip`, `accelerate_mainland`, and `ddos_protection_config`.

#### Scenario: Successful creation with required parameters
- **WHEN** user applies a configuration with valid `zone_id`, `proxy_name`, and `area`
- **THEN** the system calls `CreateL4Proxy` API and sets the resource ID to `zone_id#proxy_id`

#### Scenario: Creation with optional IPv6 parameter
- **WHEN** user applies a configuration with `ipv6` set to `on`
- **THEN** the `CreateL4Proxy` request includes `Ipv6=on`

#### Scenario: Creation with DDoS protection config
- **WHEN** user applies a configuration with `ddos_protection_config` block containing `level_mainland`, `max_bandwidth_mainland`, and `level_overseas`
- **THEN** the `CreateL4Proxy` request includes the `DDosProtectionConfig` nested struct

#### Scenario: Creation returns empty response
- **WHEN** `CreateL4Proxy` API returns nil response or nil Response field
- **THEN** the system returns an error indicating "empty response" and does NOT set an empty resource ID

#### Scenario: Creation returns empty proxy_id
- **WHEN** `CreateL4Proxy` API returns response with nil ProxyId
- **THEN** the system returns a NonRetryableError

### Requirement: Read L4 proxy instance
The system SHALL support reading a TEO L4 proxy instance via `DescribeL4Proxy` API with proxy-id filter, populating all computed and configured fields from the API response.

#### Scenario: Successful read of existing instance
- **WHEN** the resource ID is `zone_id#proxy_id` and the instance exists
- **THEN** all schema fields (`zone_id`, `proxy_id`, `proxy_name`, `area`, `cname`, `ips`, `status`, `ipv6`, `static_ip`, `accelerate_mainland`, `ddos_protection_config`, `l4proxy_rule_count`, `update_time`) are populated from the API response

#### Scenario: Read when instance no longer exists
- **WHEN** `DescribeL4Proxy` returns empty L4Proxies list for the given proxy-id
- **THEN** the system calls `d.SetId("")` to remove the resource from state and logs a warning

#### Scenario: Read with nil fields in response
- **WHEN** the API response has nil values for optional fields (e.g., `Ips`, `DDosProtectionConfig`)
- **THEN** the system skips setting those fields without error

### Requirement: Update L4 proxy instance
The system SHALL support updating mutable parameters (`ipv6`, `accelerate_mainland`) via `ModifyL4Proxy` API, and SHALL reject changes to immutable parameters (`proxy_name`, `area`, `static_ip`, `ddos_protection_config`, `zone_id`).

#### Scenario: Successful update of IPv6
- **WHEN** user changes `ipv6` from `off` to `on`
- **THEN** the system calls `ModifyL4Proxy` with updated `Ipv6` value

#### Scenario: Rejected update of immutable parameter
- **WHEN** user attempts to change `proxy_name`, `area`, `static_ip`, or `ddos_protection_config`
- **THEN** the system returns an error indicating the argument cannot be changed

#### Scenario: No API call when no mutable parameter changed
- **WHEN** no mutable parameter (`ipv6`, `accelerate_mainland`) has changed
- **THEN** the system skips calling `ModifyL4Proxy` API

### Requirement: Delete L4 proxy instance
The system SHALL support deleting a TEO L4 proxy instance via `DeleteL4Proxy` API, first stopping the instance if it is online via `ModifyL4ProxyStatus`, then deleting it.

#### Scenario: Delete online instance
- **WHEN** the instance status is `online`
- **THEN** the system first calls `ModifyL4ProxyStatus` to set status to `offline`, waits for the status change, then calls `DeleteL4Proxy`

#### Scenario: Delete already stopped instance
- **WHEN** the instance status is `offline`
- **THEN** the system directly calls `DeleteL4Proxy` without calling `ModifyL4ProxyStatus`

### Requirement: Import support
The system SHALL support importing existing L4 proxy instances using the composite ID format `zone_id#proxy_id`.

#### Scenario: Import existing instance
- **WHEN** user runs `terraform import tencentcloud_teo_l4_proxy_1.example zone-xxx#sid-xxx`
- **THEN** the resource is imported and subsequent `terraform plan` shows no changes

### Requirement: Asynchronous state waiting
The system SHALL use `StateChangeConf` to wait for the L4 proxy instance to reach the target state after create and update operations, with a timeout of 10x `ReadRetryTimeout`.

#### Scenario: Wait for instance to become online after creation
- **WHEN** `CreateL4Proxy` returns successfully
- **THEN** the system polls `DescribeL4Proxy` until the instance status is `online` or timeout is reached

#### Scenario: Timeout waiting for instance
- **WHEN** the instance does not reach the target state within the timeout
- **THEN** the system returns a timeout error

### Requirement: Documentation
The system SHALL provide a Markdown documentation file for `tencentcloud_teo_l4_proxy_1` with example usage, argument reference, and import instructions.

#### Scenario: Documentation includes import example
- **WHEN** generating documentation for `tencentcloud_teo_l4_proxy_1`
- **THEN** the documentation includes an import section showing `terraform import tencentcloud_teo_l4_proxy_1.example zone_id#proxy_id`