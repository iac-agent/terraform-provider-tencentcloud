## ADDED Requirements

### Requirement: FWM strategy resource CRUD
The provider SHALL expose a Terraform Plugin Framework resource `tencentcloud_fwm_strategy` that manages TencentCloud FWM strategies (下发规则) with Create, Read, Update, Delete and Import operations backed by the `fwm/v20250611` APIs `CreateStrategy`, `DescribeStrategy`, `ModifyStrategy` and `DeleteStrategy`.

#### Scenario: Create a strategy
- **WHEN** a user applies a configuration with `product` set (and optional `receive_account`, `receive_group`, `pre_strategy`, `post_strategy`)
- **THEN** the provider calls `CreateStrategy`, takes the first id from `Response.StrategyIds` as the resource id, and returns an error if the response or the id list is empty

#### Scenario: Read refreshes state from the cloud
- **WHEN** the provider refreshes state
- **THEN** it calls `DescribeStrategy` with the stored strategy id and updates all computed attributes from `Strategy`; if the strategy no longer exists it removes the resource from state with a warning

#### Scenario: Update applies in-place changes
- **WHEN** a user changes `receive_account`, `receive_group`, `group_id` or `sequence`
- **THEN** the provider calls `ModifyStrategy` with `StrategyId` taken from the resource id and the planned values, without replacement; changing `product`, `pre_strategy` or `post_strategy` forces replacement

#### Scenario: Delete removes the strategy
- **WHEN** the resource is destroyed
- **THEN** the provider calls `DeleteStrategy` with `StrategyIds` containing the resource id

### Requirement: Schema completeness and types
The resource schema SHALL include all fields mapped from the cloud APIs: `product` (Required), `receive_account` (List of String), `receive_group` (List of String), `pre_strategy` and `post_strategy` (SingleNested with `group_id` and `sequence`), plus Computed attributes `strategy_id`, `group_id`, `group_name`, `rule_count`, `rule_status`, `receive_account_detail` (nested `uin`/`nickname`/`receiver_type`/`members`), `sequence`, `err_msg`, `error_type`, `update_time`, `create_by`, `update_by`, `exec_area`, `create_name`, `update_name`, `create_time`.

#### Scenario: Immutable create-only fields
- **WHEN** a user modifies `product`, `pre_strategy` or `post_strategy` after creation
- **THEN** Terraform plans a replacement because those fields carry `RequiresReplace` plan modifiers

#### Scenario: Import by id
- **WHEN** a user runs `terraform import tencentcloud_fwm_strategy.example <strategy_id>`
- **THEN** the provider imports the resource by id passthrough and fills state via Read

### Requirement: Registration and documentation
The resource factory `NewFwmStrategyResource` SHALL be registered in `tencentcloud/framework/registry.go` and listed in `tencentcloud/provider.md`; a `UseFwmClient` accessor SHALL be added to `tencentcloud/connectivity/client.go`; an English `.md` document with description, Example Usage and Import sections SHALL accompany the resource.

#### Scenario: Provider builds and registers the resource
- **WHEN** the provider is built
- **THEN** `tencentcloud_fwm_strategy` is available and the connectivity client exposes `UseFwmClient()`

### Requirement: Unit tests with mocked cloud APIs
The implementation SHALL be covered by unit tests using gomonkey to mock the `fwm/v20250611` client methods, without the Terraform acceptance test suite, and SHALL compile in the repository environment.

#### Scenario: Unit test coverage
- **WHEN** unit tests are compiled and reviewed
- **THEN** Create/Read/Update/Delete model conversion and empty-response error handling paths are exercised via mocked SDK responses
