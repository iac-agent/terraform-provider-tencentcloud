## 1. Service 层实现

- [x] 1.1 在 `tencentcloud/services/teo/service_tencentcloud_teo.go` 新增 `DescribeTeoFunctionV5ById(ctx, zoneId, functionId) (*teo.Function, error)` 方法，调用 `DescribeFunctions`（入参 ZoneId + FunctionIds），返回 `Functions[0]`，空响应时返回 nil 不报错。

## 2. 资源主体代码

- [x] 2.1 创建 `tencentcloud/services/teo/resource_tc_teo_function_v5.go`，定义 `ResourceTencentCloudTeoFunctionV5()` 返回 `*schema.Resource`，含 Create/Read/Update/Delete/Importer。
- [x] 2.2 定义 Schema：入参 `zone_id`(Required,ForceNew)、`name`(Required)、`content`(Required)、`remark`(Optional)；计算字段 `function_id`、`domain`、`domain_compliance_restrictions`(TypeList 含 reason/region)、`create_time`、`update_time`。
- [x] 2.3 实现 `resourceTencentCloudTeoFunctionV5Create`：构建 `CreateFunctionRequest`，retry 调用，校验返回值非空（Response/FunctionId），先 `d.SetId(zoneId#functionId)`，再轮询 `DescribeFunctions` 直到 `Domain` 返回，最后调 Read。
- [x] 2.4 实现 `resourceTencentCloudTeoFunctionV5Read`：拆分复合 ID，调用 `DescribeTeoFunctionV5ById`，nil 时先 `log.Printf("[CRUD] teo_function_v5 id=%s", d.Id())` 再 `d.SetId("")`；逐字段 nil 检查后 `d.Set`，`domain_compliance_restrictions` 遍历构造 map 列表。
- [x] 2.5 实现 `resourceTencentCloudTeoFunctionV5Update`：拆分 ID，`immutableArgs=["zone_id","name"]` 拦截，`mutableArgs=["remark","content"]` 检测 HasChange，构建 `ModifyFunctionRequest` retry 调用，再调 Read。
- [x] 2.6 实现 `resourceTencentCloudTeoFunctionV5Delete`：拆分 ID，构建 `DeleteFunctionRequest` retry 调用。
- [x] 2.7 实现 `resourceTeoFunctionV5CreateStateRefreshFunc` 异步轮询函数：调用 `DescribeFunctions`，通过 go-template 判断 `Functions[0].Domain` 是否存在。

## 3. Provider 注册

- [x] 3.1 在 `tencentcloud/provider.go` 注册 `"tencentcloud_teo_function_v5": teo.ResourceTencentCloudTeoFunctionV5()`（按字母序插入 `tencentcloud_teo_function_component_binding` 之后）。
- [x] 3.2 在 `tencentcloud/provider.md` 新增 `tencentcloud_teo_function_v5` 行。

## 4. 文档

- [x] 4.1 创建 `tencentcloud/services/teo/resource_tc_teo_function_v5.md`：一句话描述（带 TEO 名称）、Example Usage、Import（说明使用联合 id `zone_id#function_id`）。不添加 Argument Reference / Attribute Reference。

## 5. 单元测试

- [x] 5.1 创建 `tencentcloud/services/teo/resource_tc_teo_function_v5_test.go`，使用 gomonkey mock 云 API（CreateFunctionWithContext、DescribeFunctions、ModifyFunctionWithContext、DeleteFunctionWithContext），覆盖 Create（含轮询）、Read（含 domain_compliance_restrictions）、Update、Delete 流程。

## 6. 验证

- [x] 6.1 检查所有函数返回的 error 均被处理，无需处理的用 `_ =` 赋值。
- [x] 6.2 通过 `git diff` 核对变更范围，确认未误改现有 `tencentcloud_teo_function` 资源。
