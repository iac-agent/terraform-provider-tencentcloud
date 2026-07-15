## ADDED Requirements

### Requirement: ResourceRegion parameter for tag operations
The `tencentcloud_teo_zone` resource SHALL accept an optional `ResourceRegion` string parameter that specifies the resource region for tag API operations. When not provided by the user, it SHALL default to the provider's configured region (`tcClient.Region`). This parameter SHALL be used when:
1. Reading tags via `DescribeResourceTagsByResourceIds` — passed as `ResourceRegion` in the request
2. Creating tags via `ModifyResourceTags` — used to construct the QCS resource name via `BuildTagResourceName`
3. Updating tags via `ModifyResourceTags` — used to construct the QCS resource name via `BuildTagResourceName`

#### Scenario: ResourceRegion not specified by user
- **WHEN** a user creates a `tencentcloud_teo_zone` resource without specifying `ResourceRegion`
- **THEN** the resource SHALL use the provider's default region for tag operations, and the Read operation SHALL set `ResourceRegion` to the provider's region in the Terraform state

#### Scenario: ResourceRegion explicitly specified by user
- **WHEN** a user creates a `tencentcloud_teo_zone` resource with `ResourceRegion = "ap-guangzhou"`
- **THEN** the resource SHALL use `"ap-guangzhou"` for all tag API operations (DescribeResourceTagsByResourceIds and ModifyResourceTags), and the Read operation SHALL set `ResourceRegion` to `"ap-guangzhou"` in the Terraform state

#### Scenario: ResourceRegion preserved on Read
- **WHEN** the Read operation runs for a `tencentcloud_teo_zone` resource
- **THEN** `ResourceRegion` SHALL be persisted in the Terraform state, using the value from the resource configuration if available, or the provider's default region if not configured

### Requirement: ServiceType parameter for tag operations
The `tencentcloud_teo_zone` resource SHALL accept an optional `ServiceType` string parameter that specifies the service type for tag API operations. When not provided by the user, it SHALL default to `"teo"`. This parameter SHALL be used when:
1. Reading tags via `DescribeResourceTagsByResourceIds` — passed as `ServiceType` in the request
2. Creating tags via `ModifyResourceTags` — used to construct the QCS resource name via `BuildTagResourceName`
3. Updating tags via `ModifyResourceTags` — used to construct the QCS resource name via `BuildTagResourceName`

#### Scenario: ServiceType not specified by user
- **WHEN** a user creates a `tencentcloud_teo_zone` resource without specifying `ServiceType`
- **THEN** the resource SHALL use `"teo"` as the service type for tag operations, and the Read operation SHALL set `ServiceType` to `"teo"` in the Terraform state

#### Scenario: ServiceType explicitly specified by user
- **WHEN** a user creates a `tencentcloud_teo_zone` resource with `ServiceType = "edgeone"`
- **THEN** the resource SHALL use `"edgeone"` for all tag API operations (DescribeResourceTagsByResourceIds and ModifyResourceTags), and the Read operation SHALL set `ServiceType` to `"edgeone"` in the Terraform state

#### Scenario: ServiceType preserved on Read
- **WHEN** the Read operation runs for a `tencentcloud_teo_zone` resource
- **THEN** `ServiceType` SHALL be persisted in the Terraform state, using the value from the resource configuration if available, or `"teo"` if not configured

### Requirement: Backward compatibility
Adding `ResourceRegion` and `ServiceType` parameters SHALL NOT break existing Terraform configurations that do not specify these parameters. Existing state files SHALL continue to work without modification.

#### Scenario: Existing configuration without new parameters
- **WHEN** a user applies an existing `tencentcloud_teo_zone` configuration that does not include `ResourceRegion` or `ServiceType`
- **THEN** the resource SHALL function identically to the previous version, using the hardcoded defaults (`tcClient.Region` for `ResourceRegion` and `"teo"` for `ServiceType`), and no plan diff SHALL be produced for existing resources
