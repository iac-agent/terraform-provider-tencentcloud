## Context

TEO（边缘安全加速平台）提供了即时转码模板功能，允许用户配置视频流和音频流的转码参数。当前 Terraform Provider 中尚未覆盖该资源的生命周期管理。

云 API 提供了 3 个接口：
- `CreateJustInTimeTranscodeTemplate`：创建即时转码模板，返回 TemplateId
- `DescribeJustInTimeTranscodeTemplates`：查询即时转码模板列表（支持过滤和分页）
- `DeleteJustInTimeTranscodeTemplates`：删除即时转码模板（支持批量删除）

**关键发现**：云 API 不存在 `ModifyJustInTimeTranscodeTemplate` 接口，因此该资源不支持 Update 操作。所有参数变更需通过 ForceNew（删除重建）实现。

现有 TEO 资源模式参考：`tencentcloud_igtm_strategy` 资源（RESOURCE_KIND_GENERAL），使用复合 ID 格式 `zoneId#templateId`。

## Goals / Non-Goals

**Goals:**
- 新增 `tencentcloud_teo_just_in_time_transcode_template` Terraform 资源，支持 Create/Read/Delete
- 支持视频流配置（VideoTemplateInfo）和音频流配置（AudioTemplateInfo）的嵌套结构
- 支持资源导入（Import）
- 使用 gomonkey 编写单元测试验证业务逻辑
- 生成 .md 文档文件供 make doc 使用

**Non-Goals:**
- 不实现 Update 操作（云 API 不支持）
- 不实现 datasource（本需求仅为 RESOURCE_KIND_GENERAL）
- 不修改现有 TEO 资源的任何代码

## Decisions

### 1. 资源 ID 格式：复合 ID `zoneId#templateId`
- **选择**：使用 `zoneId#templateId` 作为资源 ID
- **理由**：DeleteJustInTimeTranscodeTemplates 需要 ZoneId 和 TemplateIds 两个参数，DescribeJustInTimeTranscodeTemplates 也需要 ZoneId 作为查询条件。复合 ID 可以在 Read/Delete 时拆分获取这两个参数，与现有 TEO 资源（如 igtm_strategy）模式一致。
- **替代方案**：仅使用 TemplateId 作为 ID，但需要在 schema 中将 zone_id 标记为 Required 且不使用 ForceNew，这会导致 Delete 和 Read 操作需要从 state 中额外读取 zone_id，增加复杂度。

### 2. 所有可变参数设置 ForceNew
- **选择**：所有 Create 接口参数在变更时触发 ForceNew
- **理由**：云 API 不存在 ModifyJustInTimeTranscodeTemplate 接口，无法原地更新资源。当参数变更时，只能先删除旧资源再创建新资源。
- **替代方案**：在 Update 函数中实现删除+创建的逻辑，但这种方式不符合 Terraform 最佳实践，且 ForceNew 机制已经内置了这种行为。

### 3. Read 操作使用 DescribeJustInTimeTranscodeTemplates + Filter
- **选择**：通过 Filter 过滤 template-id 来查询单个模板
- **理由**：云 API 没有单条查询接口，只有列表查询接口 DescribeJustInTimeTranscodeTemplates。使用 `Filters` 中 `template-id` 过滤条件可以精确查询指定模板。
- **替代方案**：获取全部模板后在客户端代码中过滤，但效率较低且不必要。

### 4. 嵌套结构使用 TypeList + MaxItems: 1
- **选择**：video_template 和 audio_template 使用 TypeList 配合 MaxItems: 1
- **理由**：与代码库中其他资源的嵌套结构模式一致，TypeList 支持嵌套 schema，MaxItems: 1 限制为单条记录，同时保持与 Terraform 状态管理的兼容性。

### 5. 单元测试使用 gomonkey mock
- **选择**：使用 gomonkey 对云 API 进行 mock，只测试业务逻辑
- **理由**：根据项目规范，新增 terraform 资源不使用 terraform 测试套件，而是使用 mock 方式进行单元测试，避免依赖真实的云 API 环境。

## Risks / Trade-offs

- **[ForceNew 导致资源重建]** → 所有参数变更都会触发删除+创建，可能导致服务短暂中断。在文档中明确说明该行为。
- **[Describe 返回列表而非单条]** → Read 操作依赖列表查询接口，如果模板被删除，列表中无匹配项时应正确设置 `d.SetId("")` 标记资源已消失。
- **[无 Update 接口]** → 用户的任何参数修改都需要重建资源，与 Update 行为的用户期望不同。通过 ForceNew 机制使 Terraform 能自动处理重建。
