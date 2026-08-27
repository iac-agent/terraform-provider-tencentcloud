## 1. 服务层辅助函数

- [x] 1.1 在 `tencentcloud/services/teo/service_tencentcloud_teo.go` 中新增 `DescribeTeoFunctionV4ById(ctx, zoneId, functionId string) (*teo.Function, error)`：构造 `DescribeFunctionsRequest`（`ZoneId`、`FunctionIds=[]*string{functionId}`），在 `resource.Retry(tccommon.ReadRetryTimeout, ...)` 内调用 `DescribeFunctions`，失败用 `tccommon.RetryError(e)` 包装；当 `response.Response == nil || len(response.Response.Functions) < 1` 返回 `(nil, nil)`；严格相等校验 `*function.FunctionId == functionId` 后取第一项；defer 内打印 `[CRITAL]` 失败日志。实现与既有 `DescribeTeoFunctionById` 等价。
- [x] 1.2 确认 `connectivity.UseTeoV20220901Client()` 已就绪、SDK 包导入路径为 `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901`（无需变更）。

## 2. 资源实现

- [x] 2.1 创建 `tencentcloud/services/teo/resource_tc_teo_function_v4.go`（单文件）。顶层布局：package + imports → `ResourceTencentCloudTeoFunctionV4()` schema → `resourceTencentCloudTeoFunctionV4Create/Read/Update/Delete` → `parseTeoFunctionV4OriginalName` 辅助函数。文件开头不添加注释。
- [x] 2.2 按 spec 定义 schema，字段顺序：`zone_id`（Required + ForceNew）、`function_id`（Computed）、`name`（Required + ForceNew）、`remark`（Optional）、`content`（Required）、`domain`（Computed）、`create_time`（Computed）、`update_time`（Computed）。使用 `schema.ImportStatePassthrough` 作为 Importer。
- [x] 2.3 实现 `Create`：从 schema 取 `zone_id`、`name`、`remark`（可选）、`content` 构造 `CreateFunctionRequest`；在 `resource.Retry(tccommon.WriteRetryTimeout, ...)` 内调用 `CreateFunctionWithContext`，失败用 `tccommon.RetryError(e)` 包装并打印 `[CRITAL]` 日志；调用完成后检查 `response == nil || response.Response == nil || response.Response.FunctionId == nil || *response.Response.FunctionId == ""`，任一为空则打印 `logId` 与 `d.Id()` 后返回 `tccommon.NonRetryableError`；retry 块外设置 `d.SetId(strings.Join([]string{zoneId, functionId}, tccommon.FILED_SP))`，最后返回 `Read`。
- [x] 2.4 实现 `Read`：解析复合 ID（`tccommon.FILED_SP` 拆分，长度不为 2 返回 `id is broken`）；`_ = d.Set("zone_id", zoneId)`；调用 `DescribeTeoFunctionV4ById`；返回 `(nil, nil)` 时先 `log.Printf("[CRUD] teo_function_v4 id=%s", d.Id())` 再 `d.SetId("")` 返回无错误；否则对每个响应字段先判 nil 再 `_ = d.Set(...)`，其中 `name` 经 `parseTeoFunctionV4OriginalName(*respData.Name, zoneId)` 拆分后再 set。
- [x] 2.5 实现 `parseTeoFunctionV4OriginalName(name, zoneId string) string`：算法与既有 `ParseTeoFunctionOriginalName` 一致（以 `-`+zoneId 为后缀定位最后一次出现位置作为分割点，找不到则返回原值）。
- [x] 2.6 实现 `Update`：校验 `immutableArgs`（`name`），变更则返回 `argument %s cannot be changed`；解析复合 ID；当 `mutableArgs`（`remark`、`content`）任一变更时构造 `ModifyFunctionRequest`（`ZoneId`、`FunctionId`、`Remark`、`Content` 在用户配置时填入），在 `resource.Retry(tccommon.WriteRetryTimeout, ...)` 内调用 `ModifyFunctionWithContext`，失败用 `tccommon.RetryError(e)` 包装并打印 `[CRITAL]` 日志；retry 块内只调用接口；完成后返回 `Read`。
- [x] 2.7 实现 `Delete`：解析复合 ID；构造 `DeleteFunctionRequest`（`ZoneId`、`FunctionId`），在 `resource.Retry(tccommon.WriteRetryTimeout, ...)` 内调用 `DeleteFunctionWithContext`，失败用 `tccommon.RetryError(e)` 包装并打印 `[CRITAL]` 日志。
- [x] 2.8 每个 CRUD 函数顶部添加 `defer tccommon.LogElapsed("resource.tencentcloud_teo_function_v4.<op>")()` 与 `defer tccommon.InconsistentCheck(d, meta)()`；每次 SDK 调用成功打印 `[DEBUG]` 行（含 action、request body、response body）。日志统一使用资源名 `teo_function_v4`。
- [x] 2.9 检查所有 error 返回；对必定不出错的函数调用用 `_ =` 忽略 err，避免未使用变量错误。

## 3. Provider 注册

- [x] 3.1 在 `tencentcloud/provider.go` 的 `ResourcesMap` 中注册 `"tencentcloud_teo_function_v4": teo.ResourceTencentCloudTeoFunctionV4()`，放置在既有 `tencentcloud_teo_function` 相邻位置以保持 teo 命名空间连续。
- [x] 3.2 在 `tencentcloud/provider.md` 中追加 `tencentcloud_teo_function_v4` 资源注册说明（紧随 `tencentcloud_teo_function` 行之后）。

## 4. 资源文档

- [x] 4.1 创建 `tencentcloud/services/teo/resource_tc_teo_function_v4.md`（参考 `resource_tc_teo_function.md` / `resource_tc_igtm_strategy.md` 风格）。内容：一句话描述（含云产品名称 TEO，格式 "Provides a resource to ..."）；`Example Usage`（HCL 示例展示 `zone_id`、`name`、`content`、`remark`）；`Import` 部分说明需使用联合 ID `zone_id#function_id`。**不手写** `Argument Reference` / `Attribute Reference`（由 `make doc` 自动生成）；**不手写** `website/docs/` 下任何文件。

## 5. 单元测试

- [x] 5.1 创建 `tencentcloud/services/teo/resource_tc_teo_function_v4_test.go`，使用 gomonkey mock 云 API 客户端方法（`CreateFunctionWithContext`、`DescribeFunctions`、`ModifyFunctionWithContext`、`DeleteFunctionWithContext`），**不使用** terraform 测试套件（`resource.TestCase`）。覆盖：Create 成功路径与空 FunctionId 防御、Read 资源存在（含 name 拆分）与外部删除（`(nil,nil)` → `d.SetId("")`）、Update 不可变字段报错与可变字段变更、Delete 成功、`parseTeoFunctionV4OriginalName` 各分支。**禁止通过 `go test` 执行**，但代码须在当前环境下可正确构建执行。

## 6. 代码正确性检查

- [x] 6.1 逐一核对新增参数在云 API 各 CRUD 接口中的对应关系：`CreateFunction`（ZoneId/Name/Content/Remark）、`DescribeFunctions`（ZoneId/FunctionIds）、`ModifyFunction`（ZoneId/FunctionId/Remark/Content）、`DeleteFunction`（ZoneId/FunctionId），确保创建接口参数在创建接口存在、更新接口参数在更新接口存在，以此类推。
- [x] 6.2 核对四个接口均为同步接口、无 TaskId，确认无需异步任务轮询；确认未引入 `_extension.go` 文件（非必须不生成）。

## 7. 收尾阶段（由 tfpacer-finalize skill 统一执行）

- [ ] 7.1 由收尾阶段统一执行 `gofmt`（格式化变更的 Go 文件）、`make doc`（自动生成 `website/docs/` 下文档）、创建 `.changelog/` 下文件并 amend 推送。**禁止**在本阶段手动执行上述操作或新增 `website/`、`.changelog/` 文件。
