## 1. Schema & Create Function

- [x] 1.1 Add `enable_jumbo_frame` parameter to the resource schema in `tencentcloud/services/cvm/resource_tc_cvm_launch_template.go` (Optional, ForceNew, TypeBool, with description)
- [x] 1.2 Add `enable_jumbo_frame` handling in the `resourceTencentCloudCvmLaunchTemplateCreate` function: set `request.EnableJumboFrame` when the parameter is provided, following the pattern used for `disable_api_termination`

## 2. Read Function

- [x] 2.1 Since `DescribeLaunchTemplateVersions` response (`LaunchTemplateVersionData`) does not include `EnableJumboFrame`, do NOT add read-back logic in the `resourceTencentCloudCvmLaunchTemplateRead` function. The value will be preserved in Terraform state from the configuration (write-only behavior, consistent with `disable_api_termination`)

## 3. Unit Tests

- [x] 3.1 Add unit test in `tencentcloud/services/cvm/resource_tc_cvm_launch_template_test.go` to verify the `enable_jumbo_frame` parameter is correctly passed to the `CreateLaunchTemplate` API request using gomock/gomonkey
- [x] 3.2 Run the unit test with `go test -gcflags=all=-l` to ensure it passes

## 4. Documentation

- [x] 4.1 Update `tencentcloud/services/cvm/resource_tc_cvm_launch_template.md` to add the `enable_jumbo_frame` parameter in the example usage and description sections
