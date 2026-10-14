# Tasks: 修改 tencentcloud_dts_sync_job 资源的 tags 字段为 Set 类型

## 1. Schema 与 CRUD 代码修改

- [x] 1.1 修改 `tencentcloud/services/dts/resource_tc_dts_sync_job.go` 中 `tags` 字段的 schema 定义：`Type` 从 `schema.TypeList` 改为 `schema.TypeSet`，`Elem` 结构（`tag_key`/`tag_value` 子字段）与 `Optional/Computed/ForceNew` 属性保持不变
- [x] 1.2 修改 `resourceTencentCloudDtsSyncJobCreate` 中 tags 参数读取逻辑：`v.([]interface{})` 改为 `v.(*schema.Set).List()`，组装 `dts.TagItem` 的循环体保持不变
- [x] 1.3 核对 `resourceTencentCloudDtsSyncJobRead` 中 tags 写回逻辑：确认 `d.Set("tags", tagsList)` 在 Set 类型下工作正常，`SyncJobInfo.Tags` 为 nil 时不写回的判断保持不变
- [x] 1.4 核对 `resourceTencentCloudDtsSyncJobUpdate`/`Delete`：确认无 tags 相关逻辑残留，tags 保持 ForceNew 语义不受影响

## 2. 测试代码修改

- [x] 2.1 核对 `tencentcloud/services/dts/resource_tc_dts_sync_job_test.go` 中现有断言：`resource.TestCheckResourceAttrSet(..., "tags.#")` 在 Set 类型下仍然有效，测试 HCL 配置中 `tags` 块写法保持不变
- [x] 2.2 在 `resource_tc_dts_sync_job_test.go` 中补充 terraform 测试套件用例：创建带多个 tags 的 tencentcloud_dts_sync_job，校验 tags 数量与 tag_key/tag_value 值（覆盖 Set 类型读取与写回路径）

## 3. 文档更新

- [x] 3.1 核对 `tencentcloud/services/dts/resource_tc_dts_sync_job.md` 示例：`tags` 块 HCL 写法在 Set 类型下语义不变，确认无需修改（Argument Reference 由 `make doc` 自动生成，不手动编写）

## 4. 收尾阶段（由 tfpacer-finalize skill 统一执行）

- [ ] 4.1 对变更的 Go 文件执行 gofmt 格式化
- [ ] 4.2 执行 `make doc` 生成 website/docs/ 下的文档（tags 类型变更后重新生成 Argument Reference）
- [ ] 4.3 在 `.changelog/` 目录下新增变更记录，说明 tags 字段从 List 改为 Set 的破坏性变更及一次性状态迁移影响

## 5. 验证

- [x] 5.1 代码正确性检查：确认 tags 读取使用 `v.(*schema.Set).List()` 标准写法、CreateSyncJob 请求参数在云 API 创建接口中存在（Tags []*TagItem）、Read 写回前有 nil 判断，全部函数返回的 error 已处理
- [x] 5.2 检查无遗漏引用：全局搜索该资源文件中对 `tags` 的其他用法（含 `_test.go`），确认无残留的 List 类型断言（`v.([]interface{})`）
