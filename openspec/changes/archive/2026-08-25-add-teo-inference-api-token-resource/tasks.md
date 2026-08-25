## 1. Service 层实现

- [x] 1.1 在 `tencentcloud/services/teo/service_tencentcloud_teo.go` 中新增 `AddTeoInferenceApiToken`（调用 `CreateInferenceAPITokenWithContext`）、`DescribeTeoInferenceApiTokenById(zoneId, tokenId)`（调用 `DescribeInferenceAPITokens`，`Offset=0`、`Limit=100`，遍历 `Tokens` 按 `TokenId` 精确匹配返回）、`DeleteTeoInferenceApiToken`（调用 `DeleteInferenceAPITokenWithContext`）三个方法，统一使用 `tccommon.ReadRetryTimeout`/`tccommon.WriteRetryTimeout` 包装 `resource.Retry`，错误经 `tccommon.RetryError` 处理

## 2. Resource Schema 与 CRUD 实现

- [x] 2.1 新建 `tencentcloud/services/teo/resource_tc_teo_inference_api_token.go`，定义 `tencentcloud_teo_inference_api_token` schema：`zone_id`(Required,ForceNew)、`name`(Required,ForceNew)、`token_id`(Computed)、`content`(Computed)、`create_time`(Computed)
- [x] 2.2 实现 `resourceTencentCloudTeoInferenceApiTokenCreate`：填充 `ZoneId`、`Name` 调用 service，调用后校验 `response == nil` / `Response == nil` / `TokenId == nil` / `TokenId == ""`，任一为空则打印 logId 后返回 `NonRetryableError`；成功后设置复合 id（`zoneId # tokenId`），再调用 Read 回填
- [x] 2.3 实现 `resourceTencentCloudTeoInferenceApiTokenRead`：从 `d.Id()` 解析 `zoneId`/`tokenId`，调用 service 查询；对返回字段先判 nil 再 `d.Set`；若返回为空先 `log.Printf("[CRUD] teo_inference_api_token id=%s", d.Id())` 保留现场，再 `d.SetId("")`
- [x] 2.4 实现 `resourceTencentCloudTeoInferenceApiTokenUpdate`：因云 API 无更新接口，将 `name` 加入 `immutableArgs`，命中 `d.HasChange` 即返回 error（不调用任何 Modify 接口）
- [x] 2.5 实现 `resourceTencentCloudTeoInferenceApiTokenDelete`：从 `d.Id()` 解析 `zoneId`/`tokenId`，调用 service 删除

## 3. Provider 注册

- [x] 3.1 在 `tencentcloud/provider.go` 中注册 `tencentcloud_teo_inference_api_token` 资源（参考 `tencentcloud_igtm_strategy` 注册样式）
- [x] 3.2 在 `tencentcloud/provider.md` 中补充资源说明条目

## 4. 单元测试

- [x] 4.1 新建 `tencentcloud/services/teo/resource_tc_teo_inference_api_token_test.go`，使用 gomonkey mock 云 API（不使用 TF 测试套件），新增覆盖 Create 成功、Create 返回空 TokenId（返回 NonRetryableError）、Read 命中、Read 未命中（清空 id）、Delete 成功、Update 检测到 immutableArgs 变更返回 error 等场景的测试用例

## 5. 文档

- [x] 5.1 新建 `tencentcloud/services/teo/resource_tc_teo_inference_api_token.md`：一句话描述带上"TEO"；Example Usage；Import 部分（因使用复合 id，说明需用 `zone_id#token_id`）；不添加 Argument/Attribute Reference

## 6. 收尾验证（由 tfpacer-finalize 执行）

- [ ] 6.1 执行 `gofmt` 格式化变更的 Go 代码
- [ ] 6.2 执行 `make doc` 生成 `website/docs/` 下文档（禁止手改 website/ 目录）
- [ ] 6.3 创建 `.changelog/` 下的 changelog 文件
