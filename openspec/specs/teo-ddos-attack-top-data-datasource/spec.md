## Requirements

### Requirement: Query DDoS Attack Top Data
The system SHALL provide a Terraform data source `tencentcloud_teo_d_do_s_attack_top_data` that queries DDoS attack Top data via the `DescribeDDoSAttackTopData` API.

The data source SHALL accept the following input parameters:
- `start_time` (Required, TypeString): The start time of the query range.
- `end_time` (Required, TypeString): The end time of the query range.
- `metric_name` (Required, TypeString): The statistical metric to query.
- `zone_ids` (Optional, TypeSet of TypeString): The set of zone IDs to query.
- `policy_ids` (Optional, TypeSet of TypeInt): The set of DDoS policy IDs to query.
- `attack_type` (Optional, TypeString): The attack type filter.
- `protocol_type` (Optional, TypeString): The protocol type filter.
- `port` (Optional, TypeInt): The port number filter.
- `area` (Optional, TypeString): The data area (mainland or overseas).

The data source SHALL return the following computed attributes:
- `data` (Computed, TypeList): The list of TopEntry items, each containing:
  - `key` (Computed, TypeString): The dimension value of the Top query.
  - `value` (Computed, TypeList): The list of TopEntryValue items, each containing:
    - `name` (Computed, TypeString): The ranking entity name.
    - `count` (Computed, TypeInt): The ranking entity count.

#### Scenario: Query with required parameters only
- **WHEN** user provides `start_time`, `end_time`, and `metric_name`
- **THEN** the system calls `DescribeDDoSAttackTopData` API with the provided parameters and returns the Top data list

#### Scenario: Query with all optional parameters
- **WHEN** user provides all parameters including `zone_ids`, `policy_ids`, `attack_type`, `protocol_type`, `port`, and `area`
- **THEN** the system calls `DescribeDDoSAttackTopData` API with all provided parameters and returns the filtered Top data list

#### Scenario: API returns empty data
- **WHEN** the API response contains no Top data entries (Data is nil or empty)
- **THEN** the system SHALL return an empty `data` list without error

#### Scenario: API returns error
- **WHEN** the `DescribeDDoSAttackTopData` API call fails
- **THEN** the system SHALL retry with `tccommon.ReadRetryTimeout` and return a NonRetryableError if retries are exhausted