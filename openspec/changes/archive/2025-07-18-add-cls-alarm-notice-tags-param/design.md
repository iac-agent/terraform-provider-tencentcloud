## Context

`tencentcloud_cls_alarm_notice` 资源当前通过独立的 Tag 服务（`svctag.ModifyTags` 和 `tagService.DescribeResourceTags`）来管理标签。这与其他 CLS 资源（如 `cls_dashboard`、`cls_machine_group`、`cls_topic`）的实现方式不一致，后者直接通过 CLS API 的 `Tags` 参数传递标签。

当前实现的问题：
- Create: 先创建 AlarmNotice，再通过 tagService.ModifyTags 单独调用设置标签（非原子操作）
- Read: 通过 tagService.DescribeResourceTags 读取标签，而非从 CLS API 响应中获取
- Update: 通过 tagService.ModifyTags 更新标签，而非通过 ModifyAlarmNotice API 的 Tags 参数

CLS API 的 `CreateAlarmNoticeRequest` 和 `ModifyAlarmNoticeRequest` 均已支持 `Tags []*Tag` 参数，`AlarmNotice` 响应也包含 `Tags []*Tag` 字段。

## Goals / Non-Goals

**Goals:**
- 在 CreateAlarmNotice 请求中直接传递 Tags 参数，与 CLS API 原生支持对齐
- 在 Read 中从 CLS API 响应的 AlarmNotice.Tags 字段读取标签
- 在 Update 的 ModifyAlarmNotice 请求中直接传递 Tags 参数
- 移除对独立 Tag 服务的依赖（svctag 包）
- 保持与 `cls_dashboard` 等 CLS 资源一致的标签管理模式
- 向后兼容：现有 TF 配置和 state 不受影响

**Non-Goals:**
- 不修改 `tags` schema 字段的定义（保持 TypeMap, Optional）
- 不涉及其他 CLS 资源的修改
- 不涉及 tag 服务的其他使用者

## Decisions

### 1. Tags 传递方式

**决策**: 在 CreateAlarmNotice 和 ModifyAlarmNotice 请求中直接传递 Tags 参数，参考 `cls_dashboard` 的实现模式。

**实现方式**:
- Create: 使用 `helper.GetTags(d, "tags")` 获取标签 map，遍历转换为 `[]*cls.Tag`，设置到 `request.Tags`
- Update: 在 `d.HasChange("tags")` 时，同上方式设置 `request.Tags` 到 ModifyAlarmNotice 请求
- 移除 Create 后的 `tagService.ModifyTags` 调用
- 移除 Update 中的 `tagService.ModifyTags` 调用

**理由**: 与其他 CLS 资源保持一致，减少 API 调用次数，创建/更新操作更原子化。

**备选方案**: 保留 tag 服务调用但同时传递 CLS API Tags 参数。但这会导致标签被设置两次，可能产生冲突。

### 2. Tags 读取方式

**决策**: 从 CLS API 响应的 `AlarmNotice.Tags []*Tag` 字段读取标签，转换为 `map[string]string` 后通过 `d.Set("tags", ...)` 写入 state。

**实现方式**: 参考 `cls_dashboard` 的 Read 实现：
```go
if alarmNotice.Tags != nil {
    tags := make(map[string]string)
    for _, tag := range alarmNotice.Tags {
        if tag.Key != nil && tag.Value != nil {
            tags[*tag.Key] = *tag.Value
        }
    }
    _ = d.Set("tags", tags)
}
```

**理由**: 直接从资源 API 响应读取标签，无需额外 API 调用，与创建/更新流程一致。

**备选方案**: 保留 tag 服务读取。但这增加了不必要的 API 调用，且与 Create/Update 流程不一致。

### 3. 移除 svctag 依赖

**决策**: 移除 `resource_tc_cls_alarm_notice.go` 中对 `svctag` 包的导入，以及所有相关的 tag 服务调用代码。

**理由**: 所有标签操作均通过 CLS API 完成，不再需要独立 tag 服务。

### 4. Update 流程中 Tags 的处理位置

**决策**: 将 Tags 加入 `mutableArgs` 数组，使其在 `d.HasChange` 检测范围内。在 ModifyAlarmNotice 请求构建中，当 `d.HasChange("tags")` 时，将 Tags 参数设置到请求中。

**理由**: 与其他可变参数（name, type, notice_receivers 等）的处理方式一致。

## Risks / Trade-offs

### Risk 1: CLS API Tags 参数限制差异
**问题**: CreateAlarmNotice 的 Tags 限制为最多 50 个键值对，ModifyAlarmNotice 的 Tags 限制为最多 10 个键值对。而通过 tag 服务设置标签的限制可能不同。

**缓解措施**: Terraform schema 中 `tags` 字段已定义为 TypeMap，无数量限制。用户需遵守 CLS API 的限制。API 会返回错误提示。

### Risk 2: 旧资源 state 迁移
**问题**: 已存在的 Terraform state 中，tags 通过 tag 服务读取并存储。切换到 CLS API 读取后，标签值格式应保持一致（均为 map[string]string）。

**缓解措施**: 两种读取方式产生的 map 格式相同（key: string, value: string），无需 state 迁移。下次 `terraform refresh` 或 `terraform plan` 时会自动从 CLS API 读取最新值。

### Risk 3: 向后兼容性
**问题**: 移除 tag 服务调用后，对已有配置的行为是否一致。

**缓解措施**: `tags` schema 字段定义不变，TF 配置语法不变。标签的增删改查功能通过 CLS API 实现，与之前通过 tag 服务实现效果等价。
