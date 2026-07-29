## Context

TencentCloud EdgeOne (TEO) 边缘函数（Edge Function）允许用户在边缘节点上运行 JavaScript 代码，实现自定义的请求处理逻辑。当前 provider 中已存在 `tencentcloud_teo_function` 资源，但需要新增一个 `tencentcloud_teo_function_v3` 资源，以提供更规范的实现。

云 API 接口位于 `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901` 包中，包括：
- `CreateFunction`：创建边缘函数，入参 ZoneId/Name/Content/Remark，出参 FunctionId
- `DescribeFunctions`：查询边缘函数列表，入参 ZoneId/FunctionIds/Filters/Offset/Limit，出参 TotalCount/Functions
- `ModifyFunction`：修改边缘函数，入参 ZoneId/FunctionId/Remark/Content
- `DeleteFunction`：删除边缘函数，入参 ZoneId/FunctionId

`DescribeFunctions` 返回的 `Function` 结构体包含：FunctionId、ZoneId、Name、Remark、Content、Domain、CreateTime、UpdateTime。

## Goals / Non-Goals

**Goals:**
- 新增 `tencentcloud_teo_function_v3` 资源，支持 TEO 边缘函数的完整 CRUD 生命周期管理
- 使用联合 ID（zone_id + function_id）作为资源标识，支持 import
- 遵循现有 provider 代码规范（参考 tencentcloud_igtm_strategy 资源）
- 在 `provider.go` 和 `provider.md` 中注册新资源
- 生成对应的 `.md` 文档和单元测试文件

**Non-Goals:**
- 不修改或替换现有的 `tencentcloud_teo_function` 资源
- 不支持边缘函数规则（FunctionRule）管理，该功能由独立资源管理
- 不支持边缘函数副本（FunctionReplica）管理

## Decisions

1. **联合 ID 方案**：使用 `zone_id + function_id` 作为复合 ID（以 `tccommon.FILED_SP` 分隔），与现有 `tencentcloud_teo_function` 资源保持一致。这样在 Read/Update/Delete 中可以从 d.Id() 解析出 zone_id 和 function_id。

2. **name 字段不可变**：ModifyFunction API 不支持修改 name 字段，因此将 name 设为 immutable 参数，在 Update 方法中检查变更时返回错误。

3. **CreateFunction 后轮询**：CreateFunction API 是异步操作，创建后需要轮询 DescribeFunctions 直到函数的 Domain 字段非空，表示函数创建完成。这与现有 `tencentcloud_teo_function` 资源的行为一致。

4. **DescribeFunctions 查询**：使用 ZoneId + FunctionIds 进行查询，取返回列表中的第一个元素作为当前资源状态。

5. **计算属性**：function_id、domain、create_time、update_time 为 Computed 属性，仅在 Read 时从 API 响应中设置。

6. **测试策略**：使用 gomonkey 进行 mock 单元测试，不使用 terraform 测试套件。

## Risks / Trade-offs

- [CreateFunction 异步延迟] → 通过 StateChangeConf 轮询 DescribeFunctions 直到 Domain 字段非空，超时时间 600 秒
- [ModifyFunction 不支持修改 name] → 在 Update 方法中将 name 加入 immutableArgs，变更时返回错误
- [DescribeFunctions 返回列表而非单个资源] → 通过 FunctionIds 精确过滤，取第一个结果
- [联合 ID 解析失败] → 在 Read/Update/Delete 中检查 idSplit 长度，不等于 2 时返回错误
