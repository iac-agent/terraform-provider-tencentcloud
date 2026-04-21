## Why

TEO（边缘安全加速平台）的即时转码模板功能目前缺少 Terraform 资源支持，用户无法通过 Terraform 管理即时转码模板的生命周期（创建、查询、删除）。需要新增 `tencentcloud_teo_just_in_time_transcode_template` 资源，使基础设施即代码的管理覆盖到该功能。

## What Changes

- 新增 Terraform 资源 `tencentcloud_teo_just_in_time_transcode_template`，支持即时转码模板的 Create/Read/Delete 操作
- 由于云 API 不存在 ModifyJustInTimeTranscodeTemplate 接口，该资源不支持 Update 操作，所有可变参数变更将触发 ForceNew（删除重建）
- 新增资源对应的单元测试文件，使用 gomonkey mock 方式测试业务逻辑
- 新增资源的 .md 文档文件
- 在 provider.go 和 provider.md 中注册新资源

## Capabilities

### New Capabilities
- `teo-just-in-time-transcode-template`: 管理 TEO 即时转码模板的 Terraform 资源，支持创建、查询、删除操作，包含视频流和音频流配置参数

### Modified Capabilities

## Impact

- 新增文件：`tencentcloud/services/teo/resource_tc_teo_just_in_time_transcode_template.go`
- 新增文件：`tencentcloud/services/teo/resource_tc_teo_just_in_time_transcode_template_test.go`
- 新增文件：`tencentcloud/services/teo/resource_tc_teo_just_in_time_transcode_template.md`
- 修改文件：`tencentcloud/provider.go`（注册新资源）
- 修改文件：`tencentcloud/provider.md`（列出新资源）
- 依赖云 API：`CreateJustInTimeTranscodeTemplate`、`DescribeJustInTimeTranscodeTemplates`、`DeleteJustInTimeTranscodeTemplates`（teo v20220901）
