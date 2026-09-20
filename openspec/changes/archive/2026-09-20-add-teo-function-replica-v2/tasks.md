## 1. 资源实现（tencentcloud/services/teo/resource_tc_teo_function_replica_v2.go）

- [x] 1.1 创建 `resource_tc_teo_function_replica_v2.go`：定义 `ResourceTencentCloudTeoFunctionReplicaV2()`，schema 含 `zone_id`/`function_id`/`replica_name`（Required+ForceNew）、`content`（Required）、`remark`（Optional）、`created_on`/`modified_on`（Computed），带 `Importer: &schema.ResourceImporter{State: schema.ImportStatePassthrough}`；无文件头注释、无 `_extension.go`
- [x] 1.2 实现 `resourceTencentCloudTeoFunctionReplicaV2Create`：`d.GetOk` 填充 CreateFunctionReplicaRequest（ZoneId/FunctionId/ReplicaName/Content/Remark），`resource.Retry(tccommon.WriteRetryTimeout)` 内调 `CreateFunctionReplicaWithContext` 并做 Response nil 检查（nil 返回 NonRetryableError），retry 成功后 `d.SetId(strings.Join([]string{zoneId, functionId, replicaName}, tccommon.FILED_SP))`，最后调 Read
- [x] 1.3 实现 `resourceTencentCloudTeoFunctionReplicaV2Read`：split 复合 ID（长度≠3 报 "id is broken"），`resource.Retry(tccommon.ReadRetryTimeout)` 内调 `DescribeFunctionReplicasWithContext`（`Filters` 用 `replica-name` 精确值、`Limit=helper.Int64(200)`）；retry 外若 response/Response/列表为空或未匹配到同名副本，先 `log.Printf` 保留 id 现场再 `d.SetId("")` 返回 nil；匹配成功后仅对非 nil 字段 `d.Set`（zone_id/function_id/replica_name/content/remark/created_on/modified_on）
- [x] 1.4 实现 `resourceTencentCloudTeoFunctionReplicaV2Update`：split 复合 ID；`immutableArgs := []string{"zone_id", "function_id", "replica_name"}` 任一 HasChange 即返回 error；`content`/`remark` 任一 HasChange 时构造 ModifyFunctionReplicaRequest（ZoneId/FunctionId/ReplicaName 从 ID 解析 + Content/Remark）并 `resource.Retry(tccommon.WriteRetryTimeout)` 调 `ModifyFunctionReplicaWithContext`；完成后调 Read
- [x] 1.5 实现 `resourceTencentCloudTeoFunctionReplicaV2Delete`：split 复合 ID，构造 DeleteFunctionReplicaRequest（ZoneId/FunctionId/ReplicaNames=单元素切片），`resource.Retry(tccommon.WriteRetryTimeout)` 内调 `DeleteFunctionReplicaWithContext`
- [x] 1.6 全文件遵循规范：每个 CRUD 函数带 `defer tccommon.LogElapsed("resource.tencentcloud_teo_function_replica_v2.<op>")()` 与 `defer tccommon.InconsistentCheck(d, meta)()`；错误用 `tccommon.RetryError(e)` 包装；日志/错误信息统一使用资源名 `teo function replica v2` 措辞；检查所有 error 返回值

## 2. Provider 注册

- [x] 2.1 在 `tencentcloud/provider.go` 的 teo 资源 map 中新增 `"tencentcloud_teo_function_replica_v2": teo.ResourceTencentCloudTeoFunctionReplicaV2()`（按字典序插入）
- [x] 2.2 在 `tencentcloud/provider.md` Resource 列表中追加 `tencentcloud_teo_function_replica_v2`（紧跟 `tencentcloud_teo_function_replica` 之后）

## 3. 单元测试（tencentcloud/services/teo/resource_tc_teo_function_replica_v2_test.go）

- [x] 3.1 创建 gomonkey mock 测试文件（不使用 terraform 测试套件），复用 V1 测试的 mock 模式（mockMeta 实现 tccommon.ProviderMeta、`ApplyMethodReturn(client, "UseTeoV20220901Client", ...)`、`ApplyMethodFunc` mock 各 WithContext 方法、`schema.TestResourceDataRaw` 构造 ResourceData）
- [x] 3.2 编写用例：TestTeoFunctionReplicaV2_Create（断言请求参数与复合 ID）、_Read（断言 Limit/Filter 与 state 字段）、_ReadNotFound（断言 SetId 清空）、_Update（断言 Modify 请求参数）、_UpdateImmutableError（修改 replica_name 时返回 error 且不调 Modify）、_Delete（断言 ReplicaNames 单元素）
- [x] 3.3 核对测试代码可编译：import 无未使用项、gomonkey/testify 依赖与仓库现有用法一致、断言 helper 指针构造正确

## 4. 文档（tencentcloud/services/teo/resource_tc_teo_function_replica_v2.md）

- [x] 4.1 创建 md 文件：一句话描述（"Provides a resource to create a TEO function replica" 格式，含 TEO 产品名）+ Example Usage（hcl 示例引用 zone/function 资源）+ Import 部分（说明使用 `zone_id#function_id#replica_name` 复合 ID 导入）；不添加 Argument Reference / Attribute Reference

## 5. 代码正确性检查与收尾

- [x] 5.1 对照 design.md 中的云 API struct 定义，核对 Create/Update/Delete 请求字段均在对应云 API 入参中存在，Read 的输出字段均在 FunctionReplica 响应中存在；确认无字段遗漏或越界（CRUD 参数一致性检查）
- [x] 5.2 确认未直接新增/修改 website/ 目录文件、未在收尾外执行 gofmt/make doc/创建 .changelog 文件；不执行 go build/go vet/go test（由后续流程验证）
- [x] 5.3 最终检查：确认 V1 资源 `resource_tc_teo_function_replica.go` 及其测试/文档未被修改，变更仅包含 V2 新增文件与 provider.go/provider.md 两处注册改动
