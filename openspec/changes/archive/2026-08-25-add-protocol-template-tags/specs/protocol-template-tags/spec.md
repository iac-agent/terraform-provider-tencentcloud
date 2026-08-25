## ADDED Requirements

### Requirement: Protocol template resource SHALL support Key and Value parameters for tags
The `tencentcloud_protocol_template` resource SHALL support optional `Key` and `Value` parameters that map to the TencentCloud VPC API's Tags parameter for ServiceTemplate resources.

#### Scenario: Create protocol template with Key and Value
- **WHEN** user creates a `tencentcloud_protocol_template` resource with `Key` and `Value` parameters specified
- **THEN** the system SHALL call `CreateServiceTemplate` API with the `Tags` parameter containing the specified `Key` and `Value`
- **AND** the system SHALL set the resource ID to the returned `ServiceTemplateId`

#### Scenario: Create protocol template without Key and Value
- **WHEN** user creates a `tencentcloud_protocol_template` resource without `Key` and `Value` parameters
- **THEN** the system SHALL call `CreateServiceTemplate` API without the `Tags` parameter
- **AND** the system SHALL set the resource ID to the returned `ServiceTemplateId`

#### Scenario: Read protocol template with tags
- **WHEN** the system reads an existing `tencentcloud_protocol_template` resource that has tags
- **THEN** the system SHALL call `DescribeServiceTemplates` API to retrieve the resource details
- **AND** the system SHALL extract the `TagSet` from the response
- **AND** the system SHALL set the `Key` and `Value` state values from the first tag in `TagSet`

#### Scenario: Read protocol template without tags
- **WHEN** the system reads an existing `tencentcloud_protocol_template` resource that does not have tags
- **THEN** the system SHALL call `DescribeServiceTemplates` API to retrieve the resource details
- **AND** the system SHALL set `Key` and `Value` state values to empty strings (or nil)

### Requirement: Service layer methods SHALL support Tags parameter
The VpcService methods `CreateServiceTemplate` and `DescribeServiceTemplateById` SHALL be modified to support the Tags parameter.

#### Scenario: CreateServiceTemplate with tags
- **WHEN** `CreateServiceTemplate` method is called with `key` and `value` parameters
- **THEN** the method SHALL construct a `Tags` array with a single `Tag` element containing the specified `Key` and `Value`
- **AND** the method SHALL set this `Tags` array on the `CreateServiceTemplateRequest`

#### Scenario: CreateServiceTemplate without tags
- **WHEN** `CreateServiceTemplate` method is called with empty `key` and `value` parameters
- **THEN** the method SHALL NOT set the `Tags` parameter on the `CreateServiceTemplateRequest`

#### Scenario: DescribeServiceTemplateById returns tags
- **WHEN** `DescribeServiceTemplateById` method retrieves a ServiceTemplate that has TagSet
- **THEN** the method SHALL return the ServiceTemplate object containing the TagSet field
- **AND** the caller SHALL be able to access `template.TagSet` to retrieve tag information
