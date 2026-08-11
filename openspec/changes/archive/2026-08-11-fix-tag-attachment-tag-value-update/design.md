## Context

`tencentcloud_tag_attachment` 资源用于建立云资源与标签 KV 之间的绑定关系。当前实现中（`tencentcloud/services/tag/resource_tc_tag_attachment.go`），三个 schema 字段 `tag_key`、`tag_value`、`resource` 均被设置为 `ForceNew: true`，且资源只实现了 `Create`、`Read`、`Delete` 三个 CRUD 函数，没有 `Update` 函数。

这意味着当用户修改 `tag_value` 时，Terraform 会销毁旧资源（先删除标签绑定）再创建新资源（再绑定新标签值）。在基于标签 KV 做访问权限控制的场景下，删除标签的瞬间资源会因缺少标签而失去访问权限，导致紧接着的「绑定新标签值」请求因鉴权失败而被拒绝。

**当前关键文件：**
- `tencentcloud/services/tag/resource_tc_tag_attachment.go` - 资源定义与 CRUD（无 Update 函数）
- `tencentcloud/services/tag/service_tencentcloud_tag.go` - 服务层
  - 已存在 `ModifyTags(ctx, resourceName, replaceTags, deleteKeys)` 方法，封装云 API `ModifyResourceTags`

**云 API 能力分析（vendor 中已确认）：**

1. `ModifyResourceTags`（已被项目 `TagService.ModifyTags` 封装）：
   - 入参：`Resource`（完整资源六段式）、`ReplaceTags`（`[]*Tag{TagKey, TagValue}`）、`DeleteTags`
   - 语义：若资源未关联该标签键则增加关联；**若已关联，则将该资源关联的键对应的标签值修改为输入值**（即原地修改标签值）。
   - 优点：直接接收完整六段式字符串，与现有 Create/Delete 使用的 `AddResourceTag`/`DeleteResourceTag` 入参形式完全一致；项目中已被 ssl、cat 等多个服务使用。

2. `ModifyResourcesTagValue`（备选）：
   - 入参：`ServiceType`、`ResourceIds`、`TagKey`、`TagValue`、`ResourceRegion`、`ResourcePrefix`（需要将六段式拆解成多个字段）。
   - 缺点：需要额外把完整六段式拆解成 ServiceType/ResourceRegion/ResourcePrefix/ResourceId 等字段，解析逻辑复杂且易错。

**约束：**
- 必须保持向后兼容，现有以 `tag_key#tag_value#resource` 形式存储的 state ID 不需要迁移。
- 资源 ID 由 `tagKey + FILED_SP + tagValue + FILED_SP + resourceId` 组成；`tag_value` 可变后，Update 成功后需要用新 tag_value 重新生成并 SetId。

## Goals / Non-Goals

**Goals:**
- 移除 `tag_value` 的 `ForceNew`，使修改 `tag_value` 时资源不被销毁重建。
- 新增 `Update` 函数，当 `tag_value` 发生变更时通过 `ModifyResourceTags` API 原地修改标签值。
- 在 Update 成功后正确更新复合 ID（用新 tag_value 替换旧 tag_value）。
- 复用项目已有的 `TagService.ModifyTags` 方法，避免新增冗余的 API 封装与六段式拆解逻辑。
- 保持向后兼容，现有配置与 state 不受影响。

**Non-Goals:**
- 不修改 `tag_key` 和 `resource` 字段的行为（仍保持 `ForceNew`），因为修改标签键或目标资源等价于建立一组全新的绑定关系。
- 不引入 `ModifyResourcesTagValue` API（六段式拆解复杂，`ModifyResourceTags` 已能满足需求）。
- 不修改 schema 的字段类型或新增字段。

## Decisions

### Decision 1: 使用 `ModifyResourceTags` 而非 `ModifyResourcesTagValue`
**决策：** Update 时调用已有的 `TagService.ModifyTags`（底层 `ModifyResourceTags`），而非新增 `ModifyResourcesTagValue` 封装。

**理由：**
- `ModifyResourceTags` 直接接收完整资源六段式字符串，与现有 Create（`AddResourceTag`）、Delete（`DeleteResourceTag`）保持入参形式一致，无需额外的六段式拆解逻辑。
- 该 API 的语义明确支持「键已存在则修改其值」，正好匹配本次需求。
- 项目中已有 `TagService.ModifyTags` 的成熟封装，复用可降低实现与维护成本。

**替代方案（已否决）：**
- 新增 `ModifyResourcesTagValue` 封装：需要把六段式 `qcs::cvm:ap-guangzhou:uin/xxx:instance/ins-xxx` 拆解为 ServiceType=`cvm`、ResourceRegion=`ap-guangzhou`、ResourcePrefix=`instance`、ResourceIds=`[ins-xxx]`，解析逻辑复杂、边界情况多（如 cos 存储桶不需要 ResourcePrefix），且收益不大。

### Decision 2: Update 函数仅处理 tag_value 变更
**决策：** 在 `Update` 函数中通过 `d.HasChange("tag_value")` 检测变更；由于 `tag_key` 与 `resource` 仍是 `ForceNew`，它们不会进入 Update 流程，因此 Update 只需处理 `tag_value`。

**理由：**
- 符合 Terraform 的 ForceNew 语义：ForceNew 字段的变更由框架自动触发重建，不会进入 Update。
- Update 逻辑聚焦且简单：构造 `ReplaceTags = [{TagKey, TagValue:new}]`，`DeleteTags = nil`，调用 `ModifyTags`。

**实现：**
```go
func resourceTencentCloudTagAttachmentUpdate(d *schema.ResourceData, meta interface{}) error {
    defer tccommon.LogElapsed("resource.tencentcloud_tag_attachment.update")()
    defer tccommon.InconsistentCheck(d, meta)()

    logId := tccommon.GetLogId(tccommon.ContextNil)
    ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

    idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
    if len(idSplit) != 3 {
        return fmt.Errorf("id is broken,%s", d.Id())
    }
    oldTagKey := idSplit[0]
    oldTagValue := idSplit[1]
    resourceSixSegment := idSplit[2]

    if d.HasChange("tag_value") {
        newTagValue := d.Get("tag_value").(string)
        replaceTags := map[string]string{oldTagKey: newTagValue}
        tagService := TagService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
        if err := tagService.ModifyTags(ctx, resourceSixSegment, replaceTags, nil); err != nil {
            log.Printf("[CRITAL]%s update tag_attachment failed, reason:%+v", logId, err)
            return err
        }
        // 用新 tag_value 重新生成复合 ID
        d.SetId(oldTagKey + tccommon.FILED_SP + newTagValue + tccommon.FILED_SP + resourceSixSegment)
    }

    return resourceTencentCloudTagAttachmentRead(d, meta)
}
```

### Decision 3: Update 成功后更新复合 ID
**决策：** 因为资源 ID 的第二段是 `tag_value`，`tag_value` 变更后必须在 Update 成功后用新值重新 `d.SetId()`，否则后续 Read/Delete 会用错误的旧 tag_value 去查询。

**理由：**
- 复合 ID 是后续 Read/Delete 定位资源的唯一依据。
- 不更新 ID 会导致下次 Read 用旧 tag_value 查询而查不到（资源实际已绑定新 tag_value），从而误判资源已被删除并清空 state，造成数据丢失。

**风险缓解：**
- 在 `SetId` 之前打印 logId 与新旧值，便于排障。

### Decision 4: Update 使用标准重试与超时
**决策：** 复用 `TagService.ModifyTags` 内部已有的 `resource.Retry(tccommon.WriteRetryTimeout, ...)` 重试逻辑，Update 函数内不再额外嵌套 retry。

**理由：**
- 项目规范要求「请不要在retry中再次进行retry，若内层函数有retry，则外层函数无需添加retry逻辑」。
- `ModifyTags` 已封装了瞬态错误重试，符合规范。

## Risks / Trade-offs

### Risk 1: 行为变更（ForceNew → 原地更新）
**风险：** 修改 `tag_value` 的行为从「销毁重建」变为「原地更新」。这是破坏性变更，但符合用户预期。
**缓解措施：** 在变更日志和文档中明确说明此行为变更；现有不修改 `tag_value` 的配置完全不受影响。

### Risk 2: ID 更新失败导致 state 不一致
**风险：** 若 Update 调用 API 成功但后续 `d.SetId()` 未执行（如因 panic），state 中 ID 仍为旧 tag_value。
**缓解措施：** 在 `ModifyTags` 返回后立即执行 `SetId`，二者在同一函数体内紧密相连；并在出错路径上打印完整日志。

### Risk 3: 并发修改标签值
**风险：** 若资源在 Terraform 之外被并发修改标签，Update 可能基于过期状态。
**缓解措施：** 这是所有 Terraform 资源的通用问题，不在本次范围；Update 后会调用 Read 同步云端最新状态。
