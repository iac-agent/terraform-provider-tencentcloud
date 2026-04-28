## 1. Schema 修改

- [x] 1.1 确认 `tags` schema 字段定义不变（TypeMap, Optional），无需修改

## 2. Create 逻辑修改

- [x] 2.1 在 `resourceTencentCloudClsAlarmNoticeCreate` 函数中，在调用 CreateAlarmNotice API 之前，添加 Tags 参数处理逻辑：使用 `helper.GetTags(d, "tags")` 获取标签 map，遍历转换为 `[]*cls.Tag`，设置到 `request.Tags`
- [x] 2.2 移除 Create 函数中创建成功后的 `tagService.ModifyTags` 调用及相关代码（tagService 初始化、resourceName 构建、ModifyTags 调用）

## 3. Read 逻辑修改

- [x] 3.1 在 `resourceTencentCloudClsAlarmNoticeRead` 函数中，添加从 `alarmNotice.Tags` 读取标签的逻辑：将 `[]*cls.Tag` 转换为 `map[string]string`，通过 `d.Set("tags", tagsMap)` 写入 state
- [x] 3.2 移除 Read 函数中的 `tagService.DescribeResourceTags` 调用及相关代码（tcClient 初始化、tagService 初始化、DescribeResourceTags 调用、d.Set("tags", tags)）

## 4. Update 逻辑修改

- [x] 4.1 将 `tags` 加入 `mutableArgs` 数组，使其在 `d.HasChange` 检测范围内
- [x] 4.2 在 ModifyAlarmNotice 请求构建中，当 `d.HasChange("tags")` 时，添加 Tags 参数处理逻辑：使用 `helper.GetTags(d, "tags")` 获取标签 map，遍历转换为 `[]*cls.Tag`，设置到 `request.Tags`
- [x] 4.3 移除 Update 函数中 `d.HasChange("tags")` 块内的 tagService 相关代码（tagService 初始化、DiffTags 调用、BuildTagResourceName 调用、ModifyTags 调用）

## 5. 导入清理

- [x] 5.1 移除 `resource_tc_cls_alarm_notice.go` 中对 `svctag` 包的导入（`svctag "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tag"`），确认无其他引用
- [x] 5.2 移除 Create 函数中不再使用的 `context` 和 `fmt` 相关代码（如无其他使用）

## 6. 测试

- [x] 6.1 在 `resource_tc_cls_alarm_notice_test.go` 中补充单元测试用例，验证 Tags 通过 CLS API 传递的正确性（使用 gomonkey mock CLS API）
- [x] 6.2 使用 `go test -gcflags=all=-l` 运行单元测试验证通过

## 7. 文档

- [x] 7.1 更新 `tencentcloud/services/cls/resource_tc_cls_alarm_notice.md` 示例文件（如需调整 tags 用法说明）
