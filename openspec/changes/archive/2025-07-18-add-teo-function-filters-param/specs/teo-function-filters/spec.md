## ADDED Requirements

### Requirement: Teo Function Filters Parameter

The `tencentcloud_teo_function` resource SHALL support a `filters` parameter that maps to the `Filters` field in the `DescribeFunctions` API request, enabling users to filter functions by name and remark during the READ operation.

#### Scenario: filters schema definition

- **WHEN** the `tencentcloud_teo_function` resource schema is defined
- **THEN** the schema SHALL include a `filters` field of type `schema.TypeList`
- **AND** the `filters` field SHALL be Optional
- **AND** the `filters` field SHALL contain nested blocks with `name` (TypeString, Required) and `values` (TypeList of TypeString, Required)
- **AND** the `filters` field description SHALL indicate supported filter names: `name` (function name fuzzy match) and `remark` (function description fuzzy match)

#### Scenario: filters passed to DescribeFunctions API during READ

- **WHEN** the READ operation is executed and the user has specified `filters` in the resource configuration
- **THEN** the system SHALL pass the `filters` parameter to the `DescribeFunctions` API request
- **AND** each filter entry SHALL be converted to the SDK's `Filter` struct with `Name` and `Values` fields populated

#### Scenario: filters not specified during READ

- **WHEN** the READ operation is executed and the user has NOT specified `filters` in the resource configuration
- **THEN** the system SHALL NOT set the `Filters` field in the `DescribeFunctions` API request
- **AND** the READ operation SHALL behave identically to the current implementation

#### Scenario: filters treated as immutable

- **WHEN** the user changes the `filters` parameter in an existing resource configuration
- **THEN** the UPDATE function SHALL return an error indicating that `filters` is an immutable argument
- **AND** the `filters` parameter SHALL be included in the `immutableArgs` list in the Update function

#### Scenario: filters not used in CREATE operation

- **WHEN** the CREATE operation is executed
- **THEN** the `filters` parameter SHALL NOT be passed to the `CreateFunction` API request
- **AND** the `filters` parameter SHALL NOT affect the resource creation logic

#### Scenario: filters not used in DELETE operation

- **WHEN** the DELETE operation is executed
- **THEN** the `filters` parameter SHALL NOT be passed to the `DeleteFunction` API request
- **AND** the `filters` parameter SHALL NOT affect the resource deletion logic

#### Scenario: filters not used in UPDATE API call

- **WHEN** the UPDATE operation detects changes in mutable arguments and calls `ModifyFunction`
- **THEN** the `filters` parameter SHALL NOT be passed to the `ModifyFunction` API request
- **AND** the `filters` parameter is only used in the subsequent READ operation after UPDATE
