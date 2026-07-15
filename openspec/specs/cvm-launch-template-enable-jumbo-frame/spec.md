## ADDED Requirements

### Requirement: Enable jumbo frame parameter in launch template schema
The `tencentcloud_cvm_launch_template` resource SHALL include an `enable_jumbo_frame` parameter of type bool, Optional and ForceNew, that allows users to configure jumbo frame support for CVM instances created from the launch template.

#### Scenario: Create launch template with enable_jumbo_frame set to true
- **WHEN** a user creates a `tencentcloud_cvm_launch_template` resource with `enable_jumbo_frame = true`
- **THEN** the Create function SHALL set `request.EnableJumboFrame` to `true` in the `CreateLaunchTemplate` API request

#### Scenario: Create launch template with enable_jumbo_frame set to false
- **WHEN** a user creates a `tencentcloud_cvm_launch_template` resource with `enable_jumbo_frame = false`
- **THEN** the Create function SHALL set `request.EnableJumboFrame` to `false` in the `CreateLaunchTemplate` API request

#### Scenario: Create launch template without enable_jumbo_frame specified
- **WHEN** a user creates a `tencentcloud_cvm_launch_template` resource without specifying `enable_jumbo_frame`
- **THEN** the Create function SHALL NOT set `request.EnableJumboFrame` in the `CreateLaunchTemplate` API request

### Requirement: Write-only parameter handling
Since the `DescribeLaunchTemplateVersions` API response does not include `EnableJumboFrame` in `LaunchTemplateVersionData`, the `enable_jumbo_frame` parameter SHALL be treated as write-only. The Read function SHALL NOT attempt to read this field from the API response, and Terraform SHALL preserve the value from the configuration in state.

#### Scenario: Read operation preserves configured value
- **WHEN** a user runs `terraform refresh` or `terraform plan` on a launch template with `enable_jumbo_frame` configured
- **THEN** the Read function SHALL NOT overwrite the `enable_jumbo_frame` value in state, and the value from the configuration SHALL be preserved

### Requirement: Unit test coverage for enable_jumbo_frame
The resource SHALL include unit test coverage that verifies the `enable_jumbo_frame` parameter is correctly passed to the `CreateLaunchTemplate` API request.

#### Scenario: Unit test verifies enable_jumbo_frame parameter is set in create request
- **WHEN** the unit test creates a launch template with `enable_jumbo_frame = true`
- **THEN** the test SHALL verify that `request.EnableJumboFrame` is set to `true` in the mocked API call
