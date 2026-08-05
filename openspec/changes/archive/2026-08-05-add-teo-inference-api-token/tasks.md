## 1. 资源实现

- [x] 1.1 创建文件 `tencentcloud/services/teo/resource_tc_teo_inference_api_token_v7.go`
- [x] 1.2 实现 `ResourceTencentCloudTeoInferenceApiTokenV7()` 函数，返回 `*schema.Resource`，包含 Create/Read/Delete 方法和 Importer
- [x] 1.3 在 Resource Schema 中定义 `zone_id` 字段（TypeString, Required, ForceNew）
- [x] 1.4 在 Resource Schema 中定义 `name` 字段（TypeString, Required, ForceNew）
- [x] 1.5 在 Resource Schema 中定义 `token_id` 字段（TypeString, Computed）
- [x] 1.6 在 Resource Schema 中定义 `content` 字段（TypeString, Computed, Sensitive）
- [x] 1.7 添加 Importer 配置: `State: schema.ImportStatePassthrough`

### Create 方法
- [x] 1.8 实现 `resourceTencentCloudTeoInferenceApiTokenV7Create` 函数
- [x] 1.9 在 Create 函数开始处添加 `defer tccommon.LogElapsed()` 和 `defer tccommon.InconsistentCheck()`
- [x] 1.10 从 `d.GetOk()` 获取 `zone_id` 和 `name`，构造 `CreateInferenceAPITokenRequest`
- [x] 1.11 使用 `resource.Retry(tccommon.ReadRetryTimeout, ...)` 包装 API 调用
- [x] 1.12 在 retry 中调用 `meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoClient().CreateInferenceAPITokenWithContext()`
- [x] 1.13 检查 API 返回值是否为空（Response 为 nil 或 TokenId 为空），若为空返回 `NonRetryableError`
- [x] 1.14 设置 `d.SetId(*response.Response.TokenId)`
- [x] 1.15 调用 Read 函数更新 state

### Read 方法
- [x] 1.16 实现 `resourceTencentCloudTeoInferenceApiTokenV7Read` 函数
- [x] 1.17 在 Read 函数开始处添加 defer 语句
- [x] 1.18 使用 `d.Id()` 获取 TokenId，从 `d.GetOk("zone_id")` 获取 ZoneId
- [x] 1.19 构造 `DescribeInferenceAPITokensRequest`，设置 `ZoneId`、`Limit: 100`
- [x] 1.20 使用 `resource.Retry(tccommon.ReadRetryTimeout, ...)` 包装 API 调用
- [x] 1.21 在 Describe 响应中遍历 `Tokens` 列表，找到匹配的 `TokenId`
- [x] 1.22 若找到匹配 Token，设置 `zone_id`、`name`、`token_id`（非 nil 检查后设置）
- [x] 1.23 若 `Content` 非 nil，设置 `content`
- [x] 1.24 若未找到匹配 Token，打印日志 `log.Printf("[WARN] resource tencentcloud_teo_inference_api_token_v7 [%s] not found", d.Id())` 后调用 `d.SetId("")`

### Delete 方法
- [x] 1.25 实现 `resourceTencentCloudTeoInferenceApiTokenV7Delete` 函数
- [x] 1.26 在 Delete 函数开始处添加 defer 语句
- [x] 1.27 使用 `d.Id()` 获取 TokenId，从 `d.GetOk("zone_id")` 获取 ZoneId
- [x] 1.28 构造 `DeleteInferenceAPITokenRequest`，设置 `ZoneId` 和 `TokenId`
- [x] 1.29 使用 `resource.Retry(tccommon.ReadRetryTimeout, ...)` 包装 API 调用
- [x] 1.30 调用 `meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoClient().DeleteInferenceAPITokenWithContext()`

## 2. Provider 注册

- [x] 2.1 在 `tencentcloud/provider.go` 的 ResourcesMap 中添加资源注册: `"tencentcloud_teo_inference_api_token_v7": teo.ResourceTencentCloudTeoInferenceApiTokenV7()`
- [x] 2.2 在 `tencentcloud/provider.md` 中添加资源注册条目

## 3. 单元测试

- [x] 3.1 创建文件 `tencentcloud/services/teo/resource_tc_teo_inference_api_token_v7_test.go`
- [x] 3.2 使用 gomonkey 对云 API 进行 mock 处理
- [x] 3.3 测试 Create 成功场景：验证资源 ID 被正确设置
- [x] 3.4 测试 Create 失败场景：API 返回空 TokenId 时返回错误
- [x] 3.5 测试 Read 成功场景：验证 state 字段被正确设置
- [x] 3.6 测试 Read 资源不存在场景：验证 ID 被清空
- [x] 3.7 测试 Delete 成功场景

## 4. 文档

- [x] 4.1 创建文件 `tencentcloud/services/teo/resource_tc_teo_inference_api_token_v7.md`
- [x] 4.2 添加资源描述: "Provides a resource to create a TEO inference API token"
- [x] 4.3 添加 Example Usage 部分，包含完整的 Terraform 配置示例
- [x] 4.4 添加 Import 部分，说明导入格式: `terraform import tencentcloud_teo_inference_api_token_v7.foo <token_id>`
- [x] 4.5 不添加 `Argument Reference` 和 `Attribute Reference` 部分（由工具自动生成）

## 5. 收尾

- [ ] 5.1 使用 `gofmt` 格式化所有修改/新增的 Go 文件
- [ ] 5.2 运行 `make doc` 生成 website 文档
- [ ] 5.3 在 `.changelog/` 下创建 changelog 文件