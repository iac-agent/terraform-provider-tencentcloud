## 1. Schema 与资源定义

- [x] 1.1 创建 `tencentcloud/services/teo/resource_tc_teo_function_replica_v5.go`，定义 `ResourceTencentCloudTeoFunctionReplicaV5()`：Create/Read/Update/Delete 四个函数指针 + `Importer: &schema.ResourceImporter{State: schema.ImportStatePassthrough}`，文件开头不加注释
- [x] 1.2 定义 schema 顶层字段（平铺、不引入 `function_replicas` 包装层）：`zone_id`（Required+ForceNew）、`function_id`（Required+ForceNew）、`replica_name`（Required+ForceNew）、`content`（Required）、`remark`（Optional）、`sort_by`（Optional）、`sort_order`（Optional）、`replica_names`（Optional, TypeList of String）、`created_on`（Computed）、`modified_on`（Computed），Description 使用英文并覆盖取值说明
- [x] 1.3 定义 `filters`（Optional, TypeList，Elem 为 Resource）：`name`（Required, TypeString）、`values`（Required, TypeSet of String）、`fuzzy`（Optional, TypeBool），Description 说明 Filters.Values 上限 20 与 `replica-name` 过滤条件

## 2. Create 实现

- [x] 2.1 实现 `resourceTencentCloudTeoFunctionReplicaV5Create`：`defer tccommon.LogElapsed/InconsistentCheck`，`logId`/`ctx` 初始化，`teov20220901.NewCreateFunctionReplicaRequest()`
- [x] 2.2 用 `d.GetOk` + `helper.String` 填充 ZoneId/FunctionId/ReplicaName/Content/Remark（记录 zoneId/functionId/replicaName 局部变量用于构造复合 ID）
- [x] 2.3 `resource.Retry(tccommon.WriteRetryTimeout, ...)` 包装 `CreateFunctionReplicaWithContext`：API 错误 `tccommon.RetryError(e)`，成功打 `[DEBUG]` 日志，`result == nil || result.Response == nil` 返回 `resource.NonRetryableError`
- [x] 2.4 Retry 块外处理 `reqErr`（`[CRITAL]` 日志 + 返回）；SetId 前打印 `logId` 与入参便于排障，然后 `d.SetId(strings.Join([]string{zoneId, functionId, replicaName}, tccommon.FILED_SP))`，最后 `return resourceTencentCloudTeoFunctionReplicaV5Read(d, meta)`

## 3. Read 实现

- [x] 3.1 拆分复合 ID（`strings.Split(d.Id(), tccommon.FILED_SP)`，len != 3 返回 `"id is broken, id is %s"` 错误），填充 request.ZoneId / request.FunctionId
- [x] 3.2 透传用户查询参数：`d.GetOk("sort_by")` → request.SortBy，`d.GetOk("sort_order")` → request.SortOrder，`d.GetOk("filters")` 按 AdvancedFilter 构造模式填充 request.Filters（name/values/fuzzy）
- [x] 3.3 追加内部过滤条件 `{Name: "replica-name", Values: []*string{helper.String(replicaName)}}` 到 request.Filters，并设置 `request.Limit = helper.Int64(200)`（API 注释标注的最大值）
- [x] 3.4 `resource.Retry(tccommon.ReadRetryTimeout, ...)` 包装 `DescribeFunctionReplicasWithContext`（Retry 块内仅 API 调用与空检查），块外处理 `reqErr`
- [x] 3.5 response / response.Response 为空或未匹配到目标副本时：先 `log.Printf("[CRUD] teo_function_replica_v5 id=%s", d.Id())` 保留现场，再 `d.SetId("")` 返回 nil
- [x] 3.6 按 `*replica.ReplicaName == replicaName` 精确匹配后回填：`zone_id` / `function_id` / `replica_name`（来自 ID 拆分）、`content` / `remark` / `created_on` / `modified_on`（逐字段判 nil 后 `_ = d.Set(...)`）

## 4. Update 实现

- [x] 4.1 实现拆 ID + `needChange` / `mutableArgs := []string{"content", "remark"}` + `d.HasChange` 判断
- [x] 4.2 needChange 时构造 `teov20220901.NewModifyFunctionReplicaRequest()`，填充 ZoneId/FunctionId/ReplicaName（来自 ID）与 `d.GetOk` 的 Content/Remark
- [x] 4.3 `resource.Retry(tccommon.WriteRetryTimeout, ...)` 包装 `ModifyFunctionReplicaWithContext`（错误包装、`[DEBUG]` 日志、Response 空检查），块外处理 `reqErr`，末尾 `return resourceTencentCloudTeoFunctionReplicaV5Read(d, meta)`

## 5. Delete 实现

- [x] 5.1 实现拆 ID，填充 request.ZoneId / request.FunctionId
- [x] 5.2 `d.GetOk("replica_names")` 有值时以其构造 `request.ReplicaNames`（[]*string），否则回退 `[]*string{helper.String(replicaName)}`
- [x] 5.3 `resource.Retry(tccommon.WriteRetryTimeout, ...)` 包装 `DeleteFunctionReplicaWithContext`（错误包装、`[DEBUG]` 日志、Response 空检查），块外处理 `reqErr`，成功 `return nil`

## 6. Provider 注册

- [x] 6.1 在 `tencentcloud/provider.go` 的 ResourcesMap 中（`tencentcloud_teo_function_replica` 相邻位置）新增 `"tencentcloud_teo_function_replica_v5": teo.ResourceTencentCloudTeoFunctionReplicaV5()`，对齐相邻条目的列宽风格
- [x] 6.2 在 `tencentcloud/provider.md` 资源清单中 `tencentcloud_teo_function_replica` 下一行追加 `tencentcloud_teo_function_replica_v5`

## 7. 资源文档

- [x] 7.1 创建 `tencentcloud/services/teo/resource_tc_teo_function_replica_v5.md`：一句话描述（含 TEO 产品名，"Provides a resource to create a TEO edge function replica"）、Example Usage（HCL 示例含必填字段与可选 filters/replica_names 用法）、Import 部分（说明使用 `zone_id#function_id#replica_name` 联合 id）
- [x] 7.2 确认文档不包含 Argument Reference / Attribute Reference 部分（由工具自动生成）；不修改 `website/` 目录下任何文件

## 8. 单元测试（gomonkey mock，不使用 terraform 测试套件）

- [x] 8.1 创建 `tencentcloud/services/teo/resource_tc_teo_function_replica_v5_test.go`：package teo_test，实现 mockMeta（`GetAPIV3Conn()` 返回 `&connectivity.TencentCloudClient{}`）与 `patches.ApplyMethodReturn(..., "UseTeoV20220901Client", teoClient)`
- [x] 8.2 新增 `TestTeoFunctionReplicaV5_Create`：mock CreateFunctionReplicaWithContext（断言五个入参）+ DescribeFunctionReplicasWithContext（回显副本），验证 `d.Id()` 等于 `zone#function#replica` 与 content/remark/created_on/modified_on 回填
- [x] 8.3 新增 `TestTeoFunctionReplicaV5_Read`：mock DescribeFunctionReplicasWithContext（断言 ZoneId/FunctionId/Limit=200/Filters 含内部 replica-name 条目），验证字段回填
- [x] 8.4 新增 `TestTeoFunctionReplicaV5_ReadNotFound`：mock 返回空 FunctionReplicas，验证 `d.Id()` 被清空
- [x] 8.5 新增 `TestTeoFunctionReplicaV5_Update`：mock ModifyFunctionReplicaWithContext（断言入参）+ Describe，验证更新成功
- [x] 8.6 新增 `TestTeoFunctionReplicaV5_Delete`（默认回退单元素 ReplicaNames）与 `TestTeoFunctionReplicaV5_DeleteWithReplicaNames`（显式配置 replica_names），断言 request.ReplicaNames 内容
- [x] 8.7 检查所有 mock 与业务代码的 error 均被处理（必定不出错的用 `_ =` 吸收），保证测试代码可正确编译构建

## 9. 验证与收尾

- [x] 9.1 代码正确性自查：核对 Create/Update/Delete 入参与 vendor 中对应 Request struct 字段一一匹配（CreateFunctionReplicaRequest 五字段、ModifyFunctionReplicaRequest 五字段、DeleteFunctionReplicaRequest 三字段、DescribeFunctionReplicasRequest 八字段），不引入 vendor 不支持的参数
- [x] 9.2 确认未生成 `_extension.go` 文件、未手工改动 `website/` 与 `.changelog/`（收尾阶段由 tfpacer-finalize 统一执行 gofmt / make doc / changelog）
- [x] 9.3 通过 `git diff` 复核全部变更：仅新增 3 个资源文件 + 修改 provider.go / provider.md 两处注册
