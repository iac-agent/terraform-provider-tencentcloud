## ADDED Requirements

### Requirement: Resource Schema Definition
The `tencentcloud_teo_function_v6` resource SHALL define a schema that maps to the EdgeOne Function cloud API with the following user-configurable input fields and computed output fields.

Input fields:
- `zone_id` (string, required, ForceNew): Site ID.
- `name` (string, required, ForceNew): Function name.
- `content` (string, required): Function content (JavaScript code).
- `remark` (string, optional): Function description.

Computed fields:
- `function_id` (string): Function ID returned by CreateFunction.
- `domain` (string): Function default domain.
- `domain_compliance_restrictions` (list): Domain compliance restrictions list, each element with `reason` (string) and `region` (string).
- `create_time` (string): Creation time.
- `update_time` (string): Update time.

#### Scenario: Schema contains all required input fields
- **WHEN** the resource schema is registered
- **THEN** the schema SHALL include `zone_id`, `name`, `content` as Required and `remark` as Optional

#### Scenario: Computed fields are read-only
- **WHEN** the resource schema is registered
- **THEN** `function_id`, `domain`, `domain_compliance_restrictions`, `create_time`, `update_time` SHALL be Computed and not user-settable

#### Scenario: zone_id and name are ForceNew
- **WHEN** the user changes `zone_id` or `name` on an existing resource
- **THEN** the provider SHALL destroy and recreate the resource rather than calling ModifyFunction

### Requirement: Create Function
The resource Create method SHALL call the `CreateFunction` cloud API with `ZoneId`, `Name`, `Content`, and `Remark` parameters. After receiving the `FunctionId`, the method SHALL set the resource ID to `zoneId#functionId` (using `tccommon.FILED_SP` as separator) BEFORE polling for deployment completion. The method SHALL poll `DescribeFunctions` until the function's `Domain` field is non-empty (indicating deployment success).

#### Scenario: Successful creation
- **WHEN** the user creates a `tencentcloud_teo_function_v6` resource with valid `zone_id`, `name`, `content`, and `remark`
- **THEN** the provider SHALL call `CreateFunction`, set the composite ID `zoneId#functionId`, poll `DescribeFunctions` until `Domain` is populated, and populate all computed fields via Read

#### Scenario: CreateFunction returns empty FunctionId
- **WHEN** `CreateFunction` succeeds but `response.Response.FunctionId` is nil or empty string
- **THEN** the provider SHALL return a NonRetryableError and NOT set the resource ID

#### Scenario: ID set before async polling
- **WHEN** `CreateFunction` returns a valid `FunctionId`
- **THEN** the provider SHALL call `d.SetId()` BEFORE starting the async deployment polling, so that even if polling fails the resource ID is preserved in tfstate for `terraform destroy`

### Requirement: Read Function
The resource Read method SHALL parse the composite ID into `zoneId` and `functionId`, then call `DescribeFunctions` (via the service layer `DescribeTeoFunctionV6ById`) with `ZoneId` and `FunctionIds` to fetch the function. The method SHALL set each field only when the corresponding response field is non-nil.

#### Scenario: Function exists
- **WHEN** Read is called and `DescribeFunctions` returns a matching function
- **THEN** the provider SHALL set `zone_id`, `function_id`, `name`, `remark`, `content`, `domain`, `domain_compliance_restrictions`, `create_time`, `update_time` from the response

#### Scenario: Function not found
- **WHEN** Read is called and `DescribeFunctions` returns an empty function list
- **THEN** the provider SHALL first log `log.Printf("[CRUD] ... id=%s", d.Id())` to preserve the id context, then call `d.SetId("")` to remove the resource from state

### Requirement: Update Function
The resource Update method SHALL detect changes to `remark` and `content`, and call `ModifyFunction` with `ZoneId`, `FunctionId`, `Remark`, and `Content`. Since `ModifyFunction` does not accept a `Name` parameter, changing `name` SHALL trigger resource recreation (ForceNew).

#### Scenario: Update remark and content
- **WHEN** the user changes `remark` or `content` on an existing resource
- **THEN** the provider SHALL call `ModifyFunction` with `ZoneId`, `FunctionId`, and the changed `Remark` and/or `Content`, then call Read to refresh state

#### Scenario: Attempt to update immutable field
- **WHEN** the user changes `name` on an existing resource
- **THEN** the provider SHALL NOT call `ModifyFunction` and SHALL instead destroy and recreate the resource (ForceNew behavior)

### Requirement: Delete Function
The resource Delete method SHALL parse the composite ID into `zoneId` and `functionId`, then call `DeleteFunction` with `ZoneId` and `FunctionId`.

#### Scenario: Successful deletion
- **WHEN** the user destroys a `tencentcloud_teo_function_v6` resource
- **THEN** the provider SHALL call `DeleteFunction` with the parsed `zone_id` and `function_id`

### Requirement: Import Support
The resource SHALL support Terraform import using the composite ID `zoneId#functionId`.

#### Scenario: Import by composite ID
- **WHEN** the user runs `terraform import tencentcloud_teo_function_v6.example zoneId#functionId`
- **THEN** the provider SHALL parse the composite ID, call Read to populate state, and the resource SHALL be manageable thereafter

### Requirement: Retry and Error Handling
All cloud API calls (Create, Read, Update, Delete) SHALL be wrapped in `resource.Retry` with appropriate timeouts (`tccommon.WriteRetryTimeout` for writes, `tccommon.ReadRetryTimeout` for reads). API errors SHALL be wrapped with `tccommon.RetryError()`.

#### Scenario: Transient API failure retries
- **WHEN** a cloud API call fails with a retryable error
- **THEN** the provider SHALL retry within the configured timeout before failing

### Requirement: Provider Registration and Documentation
The resource SHALL be registered in `tencentcloud/provider.go` as `"tencentcloud_teo_function_v6"` and listed in `tencentcloud/provider.md`. A documentation file `tencentcloud/services/teo/resource_tc_teo_function_v6.md` SHALL be created with a one-line description mentioning EdgeOne/TEO, an Example Usage HCL block, and an Import section describing the composite ID.

#### Scenario: Resource registered in provider
- **WHEN** the provider is initialized
- **THEN** `tencentcloud_teo_function_v6` SHALL be available as a managed resource

#### Scenario: Documentation exists
- **WHEN** the documentation is generated via `make doc`
- **THEN** the `tencentcloud_teo_function_v6` resource page SHALL exist with description, example usage, and import instructions using the composite ID