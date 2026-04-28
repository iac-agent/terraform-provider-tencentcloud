## Context

`tencentcloud_cls_alarm_notice` 资源当前已实现 `name`、`type`、`notice_receivers`、`web_callbacks`、`tags` 参数的 schema 定义和基本 CRUD 逻辑。其中 `tags` 通过 tag 服务（`svctag.TagService`）间接管理，而非通过 CLS API 的 `CreateAlarmNoticeRequest.Tags` 和 `ModifyAlarmNoticeRequest.Tags` 字段直接传递。

CLS API 的 `CreateAlarmNoticeRequest` 和 `ModifyAlarmNoticeRequest` 均支持 `Tags` 字段（类型 `[]*cls.Tag`，含 `Key`/`Value`），允许在创建和修改告警通知渠道组时直接传递标签。当前实现未利用此字段，而是先创建资源，再通过 tag 服务的 `ModifyTags` 接口设置标签。同样，`AlarmNotice` 响应中的 `Tags` 字段未被读取，标签读取完全依赖 tag 服务的 `DescribeResourceTags`。

此外，`ModifyAlarmNoticeRequest` 已支持 `Name`、`Type`、`NoticeReceivers`、`WebCallbacks` 字段，当前代码已正确映射这些参数。

## Goals / Non-Goals

**Goals:**
- 在 `resourceTencentCloudClsAlarmNoticeCreate` 中通过 `CreateAlarmNoticeRequest.Tags` 直接传递标签参数
- 在 `resourceTencentCloudClsAlarmNoticeUpdate` 中通过 `ModifyAlarmNoticeRequest.Tags` 直接传递标签参数
- 在 `resourceTencentCloudClsAlarmNoticeRead` 中从 `AlarmNotice.Tags` 读取标签信息
- 保持向后兼容：现有使用 tag 服务的逻辑保留作为兜底

**Non-Goals:**
- 不修改 `name`、`type`、`notice_receivers`、`web_callbacks` 的现有实现（已正确映射）
- 不修改 `tags` 的 schema 定义（保持 TypeMap）
- 不移除现有 tag 服务的标签管理逻辑

## Decisions

### 1. Create 中 Tags 直接传递

**决策**: 在 `resourceTencentCloudClsAlarmNoticeCreate` 中，将 `tags` 参数从 schema 的 `TypeMap` 格式转换为 `[]*cls.Tag` 并赋值给 `CreateAlarmNoticeRequest.Tags`，同时在创建成功后保留 tag 服务逻辑作为兜底。

**理由**:
- CLS API 原生支持在创建时传递 Tags，减少一次额外的 API 调用
- 保留 tag 服务逻辑可确保即使 CLS API 的 Tags 字段行为变更，标签仍能正确设置
- `TypeMap` 到 `[]*cls.Tag` 的转换逻辑简单：遍历 map 的 key-value 对，构建 `cls.Tag{Key: &k, Value: &v}`

### 2. Update 中 Tags 直接传递

**决策**: 在 `resourceTencentCloudClsAlarmNoticeUpdate` 中，当 `tags` 变更时，将新标签转换为 `[]*cls.Tag` 并赋值给 `ModifyAlarmNoticeRequest.Tags`。

**理由**:
- 与 Create 保持一致，通过 CLS API 直接传递标签
- `ModifyAlarmNoticeRequest.Tags` 可同时设置新增和修改的标签

### 3. Read 中 Tags 读取

**决策**: 在 `resourceTencentCloudClsAlarmNoticeRead` 中，从 `AlarmNotice.Tags` 字段读取标签，将 `[]*cls.Tag` 转换为 `map[string]interface{}` 格式后设置到 state。

**理由**:
- 直接从 CLS API 响应读取标签，减少对 tag 服务的依赖
- 转换逻辑：遍历 `[]*cls.Tag`，将每个 `Tag.Key`/`Tag.Value` 添加到 map 中

### 4. Tag 格式转换

**决策**: 封装两个辅助函数：
- `clsTagsToMap(tags []*cls.Tag) map[string]interface{}`: 将 `[]*cls.Tag` 转换为 `map[string]interface{}`
- `mapToClsTags(tagsMap map[string]interface{}) []*cls.Tag`: 将 `map[string]interface{}` 转换为 `[]*cls.Tag`

**理由**: 复用转换逻辑，避免在 Create/Read/Update 中重复代码。

## Risks / Trade-offs

### Risk 1: Create 中 Tags 和 tag 服务双重设置
**问题**: 在 Create 中同时通过 `request.Tags` 和 tag 服务设置标签，可能导致重复设置。

**缓解措施**: 优先使用 `request.Tags`，tag 服务逻辑作为兜底。如果 CLS API 的 Tags 字段已正确处理标签，tag 服务的 `ModifyTags` 调用将是幂等的（设置相同的标签不会产生副作用）。

### Risk 2: ModifyAlarmNoticeRequest.Tags 与 tag 服务的行为差异
**问题**: `ModifyAlarmNoticeRequest.Tags` 可能是全量替换语义，而 tag 服务的 `ModifyTags` 支持增量更新（添加/删除指定标签）。

**缓解措施**: 在 Update 中，计算全量标签列表后通过 `ModifyAlarmNoticeRequest.Tags` 传递，确保标签状态与 Terraform 配置一致。保留 tag 服务用于处理标签的增量差异计算。
