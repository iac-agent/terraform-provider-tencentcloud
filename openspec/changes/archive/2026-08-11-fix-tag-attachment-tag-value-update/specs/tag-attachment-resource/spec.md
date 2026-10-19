## ADDED Requirements

### Requirement: tag_value 参数支持原地修改
`tencentcloud_tag_attachment` 资源 SHALL 支持在不销毁重建资源的情况下修改 `tag_value` 参数。修改 `tag_value` 时 SHALL 通过腾讯云标签服务的 `ModifyResourceTags` API 原地更新资源的标签值，而不是先解绑旧标签再绑定新标签。

#### Scenario: 用户修改已有标签绑定关系的 tag_value
- **GIVEN** 一个已创建的 `tencentcloud_tag_attachment` 资源，tag_key=`运营产品`，tag_value=`A`，resource 指向某个云资源
- **WHEN** 用户将 Terraform 配置中的 `tag_value` 从 `A` 修改为 `B`
- **THEN** Provider SHALL 调用 `ModifyResourceTags` API，将 `ReplaceTags` 设置为 `[{TagKey:"运营产品", TagValue:"B"}]`、`Resource` 设置为资源六段式
- **AND** Provider SHALL NOT 先调用 `DeleteResourceTag` 解绑旧标签
- **AND** 在整个更新过程中标签键 `运营产品` 始终绑定在该资源上，不产生「无标签」的时间窗口

#### Scenario: tag_value 未变更时不触发更新
- **GIVEN** 一个已创建的 `tencentcloud_tag_attachment` 资源
- **WHEN** 用户修改 Terraform 配置，但 `tag_value` 保持不变
- **THEN** Provider SHALL NOT 调用 `ModifyResourceTags` API
- **AND** 资源保持不变

#### Scenario: tag_key 与 resource 仍保持 ForceNew
- **GIVEN** 一个已创建的 `tencentcloud_tag_attachment` 资源
- **WHEN** 用户修改 `tag_key` 或 `resource` 参数
- **THEN** Provider SHALL 销毁旧资源并创建新资源（即触发 ForceNew 重建）
- **AND** SHALL NOT 进入 Update 流程

### Requirement: tag_value 移除 ForceNew 配置
`tencentcloud_tag_attachment` 资源的 `tag_value` 字段 SHALL 不再设置 `ForceNew: true`，从而使 Terraform 框架将其变更路由到 Update 函数而非重建流程。

#### Scenario: tag_value 字段属性校验
- **GIVEN** 资源 Schema 定义
- **WHEN** 检查 `tag_value` 字段属性
- **THEN** 该字段 SHALL 为 `Required`、`TypeString`
- **AND** SHALL NOT 包含 `ForceNew: true`
- **AND** 当 `tag_value` 发生变更时 SHALL 触发 Update 操作而非重建

### Requirement: Update 成功后更新复合 ID
资源 ID 由 `tag_key`、`tag_value`、`resource` 三段以 `tccommon.FILED_SP` 连接而成。当 `tag_value` 变更并成功更新云端标签值后，Provider SHALL 使用新的 `tag_value` 重新生成并设置资源 ID，以保证后续 Read/Delete 能用正确的 tag_value 定位资源。

#### Scenario: 更新 tag_value 后 ID 反映最新值
- **GIVEN** 一个 `tencentcloud_tag_attachment` 资源，原 ID 为 `运营产品#A#<resource六段式>`
- **WHEN** 用户将 `tag_value` 修改为 `B` 且 `ModifyResourceTags` API 调用成功
- **THEN** Provider SHALL 调用 `d.SetId("运营产品" + FILED_SP + "B" + FILED_SP + <resource六段式>)`
- **AND** 后续基于该 ID 的 Read 操作 SHALL 能查到 tag_key=`运营产品`、tag_value=`B` 的绑定关系

#### Scenario: ID 格式异常时拒绝更新
- **GIVEN** 资源的 ID 不符合 `tag_key + FILED_SP + tag_value + FILED_SP + resource` 的三段格式
- **WHEN** 执行 Update 操作解析 ID
- **THEN** Provider SHALL 返回描述 ID 损坏的错误，且不调用任何云 API

### Requirement: 复用 TagService.ModifyTags 服务层方法
Update 操作 SHALL 复用服务层已有的 `TagService.ModifyTags(ctx, resourceName, replaceTags, deleteKeys)` 方法（底层对应 `ModifyResourceTags` API），而不新增针对 `ModifyResourcesTagValue` 的服务层封装。

#### Scenario: Update 调用 ModifyTags 并传入完整六段式
- **GIVEN** `tag_value` 发生变更，资源六段式为 `qcs::cvm:ap-guangzhou:uin/xxx:instance/ins-xxx`
- **WHEN** 执行 Update
- **THEN** Provider SHALL 调用 `tagService.ModifyTags(ctx, "qcs::cvm:ap-guangzhou:uin/xxx:instance/ins-xxx", {tag_key: new_tag_value}, nil)`
- **AND** SHALL NOT 将六段式拆解为 ServiceType/ResourceRegion/ResourcePrefix 等字段

#### Scenario: 云 API 错误处理
- **GIVEN** `ModifyResourceTags` API 返回错误
- **WHEN** 错误为可重试类型（如网络错误、限流）
- **THEN** `ModifyTags` 内部 SHALL 使用 `tccommon.RetryError` 包装并通过 `resource.Retry(tccommon.WriteRetryTimeout, ...)` 重试
- **AND** 当错误不可重试或重试耗尽时 SHALL 返回错误给调用方
