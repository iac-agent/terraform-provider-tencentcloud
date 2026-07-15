## Proposal: Add `enable_jumbo_frame` parameter to `tencentcloud_cvm_launch_template`

### Background

The CVM `CreateLaunchTemplate` API now supports the `EnableJumboFrame` parameter, which allows instances created from the launch template to enable jumbo frames (only instance types that support jumbo frames can use this feature). The current Terraform resource `tencentcloud_cvm_launch_template` does not expose this parameter.

### What

Add the `enable_jumbo_frame` parameter to the `tencentcloud_cvm_launch_template` resource schema, and pass it to the `CreateLaunchTemplate` API request during resource creation.

### Why

Users need to configure jumbo frame support when creating CVM launch templates through Terraform. Without this parameter, they cannot set this important networking feature via Terraform, forcing them to use the console or API directly.

### Key Constraints

- `EnableJumboFrame` is a write-only parameter: the CVM `DescribeLaunchTemplateVersions` API (used for Read) does NOT return `EnableJumboFrame` in `LaunchTemplateVersionData`. This means the value will be preserved in Terraform state from the configuration (consistent with `disable_api_termination` behavior in this resource).
- Since this resource is CRD-only (no Update operation), the parameter must be `ForceNew`.
- The parameter type is `bool` and is `Optional`.
