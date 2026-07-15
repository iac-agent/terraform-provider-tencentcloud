## Context

The `tencentcloud_cvm_launch_template` resource is a RESOURCE_KIND_GENERAL resource in the CVM service that only supports Create, Read, and Delete operations (no Update). The current resource schema does not include the `EnableJumboFrame` parameter, which was recently added to the CVM `CreateLaunchTemplate` API.

The CVM SDK (`github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312`) already includes `EnableJumboFrame` as a `*bool` field on `CreateLaunchTemplateRequest`. However, the `LaunchTemplateVersionData` response struct (returned by `DescribeLaunchTemplateVersions`) does not include this field, meaning it cannot be read back from the API.

This is similar to the existing `disable_api_termination` parameter in the same resource, which is also set during creation but not returned by the Read API (currently commented out in the Read function).

## Goals / Non-Goals

**Goals:**
- Add the `enable_jumbo_frame` parameter to the `tencentcloud_cvm_launch_template` resource schema
- Pass the parameter value to the `CreateLaunchTemplate` API request during resource creation
- Handle the write-only nature of the parameter correctly (stored in state but not read from API)
- Add unit test coverage for the new parameter

**Non-Goals:**
- Adding the `Metadata` (Key/Value) parameter - these are nested inside a complex struct and the Read API also doesn't return them; these should be addressed in a separate change if needed
- Adding Update functionality to this resource (it remains CRD-only)
- Modifying any other CVM resources

## Decisions

1. **Parameter schema definition**: `enable_jumbo_frame` will be defined as `Optional`, `ForceNew`, `TypeBool` with a default of `false`. This follows the pattern of other bool parameters in this resource (e.g., `disable_api_termination`, `dry_run`).

   - *Alternative considered*: Making it `Computed` with a default - rejected because the API doesn't return this value, so Computed would cause perpetual drift.

2. **Read behavior**: Since `DescribeLaunchTemplateVersions` does not return `EnableJumboFrame` in `LaunchTemplateVersionData`, the parameter will NOT be read back from the API. Terraform will preserve the value from the configuration in state. This is the same approach used for `disable_api_termination` in this resource.

3. **ForceNew**: Since this resource has no Update operation (CRD-only), all parameters must be ForceNew. Changing `enable_jumbo_frame` will trigger resource recreation, which is consistent with existing behavior.

## Risks / Trade-offs

- [Write-only parameter] → The `enable_jumbo_frame` value stored in Terraform state may not reflect the actual cloud configuration if it's modified outside Terraform. Users should be aware of this limitation. This is acceptable because the same pattern is already used for `disable_api_termination`.
- [No API read support] → If the CVM API later adds `EnableJumboFrame` to the `LaunchTemplateVersionData` response, the Read function should be updated to read it back. This can be done in a future change without breaking existing configurations.
