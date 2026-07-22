## ADDED Requirements

### Requirement: Resource schema definition
The `tencentcloud_cvm_chc_network_mode` resource SHALL define the following schema parameters:
- `chc_ids`: Required, TypeList of TypeString, ForceNew - CHC物理服务器ID列表
- `network_mode`: Required, TypeString - 所要切换的网络模式，枚举值：DEPLOY（部署网络模式）、BUSINESS（业务网络模式）

#### Scenario: Valid schema with required fields
- **WHEN** user provides both `chc_ids` and `network_mode` in the resource configuration
- **THEN** the resource SHALL accept the configuration and proceed with Create

#### Scenario: Missing required field
- **WHEN** user omits either `chc_ids` or `network_mode`
- **THEN** Terraform SHALL reject the configuration with a required field error

### Requirement: Resource Create operation
The resource Create operation SHALL call `ModifyChcNetworkMode` API with the provided `chc_ids` and `network_mode` parameters. The resource ID SHALL be constructed by joining `chc_ids` with `tccommon.FILED_SP` as separator.

#### Scenario: Successful network mode switch on create
- **WHEN** user creates the resource with `chc_ids = ["chc-xxx"]` and `network_mode = "DEPLOY"`
- **THEN** the resource SHALL call `ModifyChcNetworkMode` API and set the resource ID to the chc_ids value

#### Scenario: API returns empty response
- **WHEN** `ModifyChcNetworkMode` API returns nil response
- **THEN** the resource SHALL return `NonRetryableError` to prevent writing empty ID

#### Scenario: Multiple CHC IDs
- **WHEN** user creates the resource with multiple `chc_ids` (e.g., `["chc-1a2b3c4d", "chc-5e6f7g8h"]`)
- **THEN** the resource SHALL join the IDs with `tccommon.FILED_SP` as the resource ID and pass all IDs to the API

### Requirement: Resource Read operation
The resource Read operation SHALL call `DescribeChcHosts` API to verify the CHC hosts still exist. The `chc_id` field SHALL be set from the API response. The `network_mode` field SHALL be preserved from the Terraform state (since the `ChcHost` structure does not contain a `NetworkMode` field).

#### Scenario: CHC hosts exist
- **WHEN** Read is called and `DescribeChcHosts` returns the CHC hosts
- **THEN** the resource SHALL set `chc_id` from the API response and preserve `network_mode` from state

#### Scenario: CHC hosts not found
- **WHEN** Read is called and `DescribeChcHosts` returns empty results
- **THEN** the resource SHALL log `[CRUD] chc_network_mode id=<id>` and then call `d.SetId("")` to remove from state

### Requirement: Resource Update operation
The resource Update operation SHALL call `ModifyChcNetworkMode` API when `network_mode` changes. When `chc_ids` changes (which is ForceNew), Terraform SHALL trigger a destroy and recreate cycle.

#### Scenario: Network mode change
- **WHEN** user updates `network_mode` from "DEPLOY" to "BUSINESS"
- **THEN** the resource SHALL call `ModifyChcNetworkMode` API with the new network_mode value

#### Scenario: No change
- **WHEN** no parameters change
- **THEN** the resource SHALL not call any API and proceed to Read

### Requirement: Resource Delete operation
The resource Delete operation SHALL only remove the resource from Terraform state by calling `d.SetId("")`. It SHALL NOT call any cloud API since there is no corresponding delete/revert operation for network mode switching.

#### Scenario: Resource deletion
- **WHEN** user destroys the resource
- **THEN** the resource SHALL remove itself from state without calling any cloud API

### Requirement: Provider registration
The `tencentcloud_cvm_chc_network_mode` resource SHALL be registered in `tencentcloud/provider.go` with the resource name `"tencentcloud_cvm_chc_network_mode"` and mapped to `cvm.ResourceTencentCloudCvmChcNetworkMode()`. An entry SHALL also be added to `tencentcloud/provider.md`.

#### Scenario: Provider initialization
- **WHEN** Terraform initializes the provider
- **THEN** the `tencentcloud_cvm_chc_network_mode` resource SHALL be available for use

### Requirement: Resource documentation
A markdown documentation file SHALL be created at `tencentcloud/services/cvm/resource_tc_cvm_chc_network_mode.md` following the gendoc format, including:
- One-sentence description mentioning CVM product
- Example Usage section with HCL configuration
- No Import section (since this resource does not support import)
- No Argument Reference or Attribute Reference sections (auto-generated)

#### Scenario: Documentation exists
- **WHEN** the resource is implemented
- **THEN** the .md documentation file SHALL exist with proper format and content

### Requirement: Unit tests
Unit test file SHALL be created at `tencentcloud/services/cvm/resource_tc_cvm_chc_network_mode_test.go` using mock (gomonkey) approach for unit testing business logic, not Terraform test suite.

#### Scenario: Unit tests cover CRUD logic
- **WHEN** unit tests are run with `go test -gcflags=all=-l`
- **THEN** all tests SHALL pass covering Create, Read, Update, and Delete operations
