## ADDED Requirements

### Requirement: Create teo function v5 resource
The system SHALL provide a `tencentcloud_teo_function_v5` resource that creates an EdgeOne edge function via the `CreateFunction` API with `zone_id`, `name`, `content`, `remark` parameters.

#### Scenario: Successful creation
- **WHEN** a user applies a `tencentcloud_teo_function_v5` resource with required `zone_id`, `name`, `content` and optional `remark`
- **THEN** the system SHALL call `CreateFunction`, set the resource id to `zone_id#function_id`, poll `DescribeFunctions` until the `Domain` field is returned, and read back all computed fields.

#### Scenario: Create return value validation
- **WHEN** `CreateFunction` returns a nil response or an empty `FunctionId`
- **THEN** the system SHALL return a `NonRetryableError` and not set the resource id.

### Requirement: Read teo function v5 resource
The system SHALL read the edge function via `DescribeFunctions` using the composite id `zone_id#function_id` and populate computed fields including `function_id`, `domain`, `domain_compliance_restrictions`, `create_time`, `update_time`.

#### Scenario: Resource exists
- **WHEN** the system reads a `tencentcloud_teo_function_v5` resource whose id maps to an existing function
- **THEN** the system SHALL set all computed fields, flattening `DomainComplianceRestrictions` into a list of maps with `reason` and `region` keys.

#### Scenario: Resource not found
- **WHEN** `DescribeFunctions` returns no matching function
- **THEN** the system SHALL log `[CRUD] teo_function_v5 id=<id>` and call `d.SetId("")` to remove the resource from state.

### Requirement: Update teo function v5 resource
The system SHALL update the edge function via `ModifyFunction` when `remark` or `content` changes, and SHALL reject changes to `name`.

#### Scenario: Update mutable fields
- **WHEN** a user changes `remark` or `content`
- **THEN** the system SHALL call `ModifyFunction` with `zone_id`, `function_id`, and the changed fields, then read back the resource.

#### Scenario: Reject immutable field change
- **WHEN** a user changes `name`
- **THEN** the system SHALL return an error indicating the argument cannot be changed.

### Requirement: Delete teo function v5 resource
The system SHALL delete the edge function via `DeleteFunction` using `zone_id` and `function_id`.

#### Scenario: Successful deletion
- **WHEN** a user destroys a `tencentcloud_teo_function_v5` resource
- **THEN** the system SHALL call `DeleteFunction` and return nil on success.

### Requirement: Import teo function v5 resource
The system SHALL support importing a `tencentcloud_teo_function_v5` resource using the composite id `zone_id#function_id`.

#### Scenario: Import existing function
- **WHEN** a user runs `terraform import tencentcloud_teo_function_v5.foo zone_id#function_id`
- **THEN** the system SHALL populate the resource state by reading the function and mapping the composite id.

### Requirement: Provider registration for teo function v5
The system SHALL register `tencentcloud_teo_function_v5` in `provider.go` and document it in `provider.md`.

#### Scenario: Resource available in provider
- **WHEN** the provider is initialized
- **THEN** `tencentcloud_teo_function_v5` SHALL be a registered resource type backed by `teo.ResourceTencentCloudTeoFunctionV5()`.