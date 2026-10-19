## 1. Schema 修改

- [x] 1.1 在 `tencentcloud/services/tag/resource_tc_tag_attachment.go` 中移除 `tag_value` 字段的 `ForceNew: true` 配置，使其支持原地更新
- [x] 1.2 在 `ResourceTencentCloudTagAttachment()` 返回的 `schema.Resource` 中注册 `Update: resourceTencentCloudTagAttachmentUpdate`

## 2. Update 函数实现

- [x] 2.1 在 `tencentcloud/services/tag/resource_tc_tag_attachment.go` 中新增 `resourceTencentCloudTagAttachmentUpdate` 函数
- [x] 2.2 添加标准 defer 处理（`tccommon.LogElapsed("resource.tencentcloud_tag_attachment.update")()` 与 `tccommon.InconsistentCheck(d, meta)()`）
- [x] 2.3 从 `d.Id()` 解析复合 ID（`tag_key + FILED_SP + tag_value + FILED_SP + resource`），校验三段格式，格式异常时返回错误且不调用云 API
- [x] 2.4 通过 `d.HasChange("tag_value")` 检测变更，未变更时不调用任何云 API
- [x] 2.5 当 `tag_value` 变更时，读取新值，构造 `replaceTags = {oldTagKey: newTagValue}`，调用 `TagService.ModifyTags(ctx, resourceSixSegment, replaceTags, nil)`（复用已有服务层方法，传入完整六段式资源字符串）
- [x] 2.6 调用成功后用新 tag_value 重新生成复合 ID 并 `d.SetId(oldTagKey + FILED_SP + newTagValue + FILED_SP + resourceSixSegment)`
- [x] 2.7 Update 函数末尾调用 `resourceTencentCloudTagAttachmentRead(d, meta)` 同步云端状态
- [x] 2.8 在错误路径上打印 logId 与资源名称（tag_attachment）便于日志检索

## 3. 服务层

- [x] 3.1 确认 `service_tencentcloud_tag.go` 中已存在 `ModifyTags(ctx, resourceName, replaceTags, deleteKeys)` 方法且底层对应 `ModifyResourceTags` API，无需新增封装

## 4. 单元测试

- [x] 4.1 在 `tencentcloud/services/tag/resource_tc_tag_attachment_test.go` 中新增针对 Update 的单元测试（使用 gomonkey mock `TagService.ModifyTags`，仅测试业务逻辑）
- [x] 4.2 测试用例：tag_value 变更时调用 ModifyTags 且传入完整六段式、ReplaceTags 包含正确的 key/value、调用成功后 ID 更新为新 tag_value
- [x] 4.3 测试用例：tag_value 未变更时不调用 ModifyTags
- [x] 4.4 测试用例：ID 格式异常（非三段）时 Update 返回错误且不调用云 API
- [x] 4.5 测试用例：ModifyTags 返回错误时 Update 将错误向上返回

## 5. 文档与变更日志（收尾阶段执行）

- [ ] 5.1 更新 `tencentcloud/services/tag/resource_tc_tag_attachment.md` 示例文档，补充修改 tag_value 的说明（由 tfpacer-finalize 阶段处理）
- [ ] 5.2 在收尾阶段通过 `make doc` 重新生成 `website/docs/` 下的 markdown 文档（禁止手动编写 website/ 文件）
- [ ] 5.3 在收尾阶段创建 `.changelog/<PR_NUMBER>.txt` 变更日志条目，格式：`bugfix: resource/tencentcloud_tag_attachment: tag_value no longer forces resource recreation, supports in-place update via ModifyResourceTags`

## 6. 代码质量检查（收尾阶段执行）

- [ ] 6.1 在收尾阶段运行 `gofmt` 格式化代码
- [ ] 6.2 确认所有函数返回的 error 均被检查；必定不出错的用 `_ = func()` 处理
- [ ] 6.3 编译验证（由后续流程执行，不在本阶段运行 go build）

## 依赖关系
- 任务 2 依赖任务 1 完成（需先注册 Update 函数）
- 任务 4 依赖任务 1、2 完成（需新增的 Update 函数存在）
- 任务 5、6 为收尾阶段任务，依赖所有代码修改完成
