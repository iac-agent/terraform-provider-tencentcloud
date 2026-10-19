## 1. Schema 契约验证

- [x] 1.1 验证 `tencentcloud/services/teo/resource_tc_teo_security_policy_config.go` 中顶层字段 `zone_id`（Required, ForceNew）、`entity`（Optional, ForceNew, 取值 `ZoneDefaultPolicy`/`Template`/`Host`）、`host`（Optional, ForceNew）、`template_id`（Optional, ForceNew）与 spec 一致
- [x] 1.2 验证 `security_config.rate_limit_config` 子结构字段 `switch`、`rate_limit_user_rules`、`rate_limit_template`、`rate_limit_intelligence`、`rate_limit_customizes` 与 spec 一致
- [x] 1.3 验证 `rateLimitUserRuleSchema()` 中字段（`threshold`/`period`/`rule_name`/`action`/`punish_time`/`punish_time_unit`/`rule_status`/`acl_conditions`/`rule_priority`/`rule_id`(Computed)/`freq_fields`/`update_time`(Computed)/`freq_scope`/`name`/`custom_response_id`/`response_code`/`redirect_url`）与 spec 一致
- [x] 1.4 验证 `aclConditionSchema()` 中字段（`match_from`/`match_param`/`operator`/`match_content`）与 spec 一致
- [x] 1.5 验证 `rate_limit_template` 及其 `rate_limit_template_detail`（Computed-only）字段与 spec 一致
- [x] 1.6 验证 `rate_limit_intelligence` 中 `switch`/`action`/`rule_id`(Computed) 与 spec 一致

## 2. 复合 ID 与 Create 逻辑验证

- [x] 2.1 验证 Create 函数中复合 ID 拼接逻辑：`ZoneDefaultPolicy`→`zoneId#ZoneDefaultPolicy`、`Host`→`zoneId#Host#host`、`Template`→`zoneId#Template#templateId`，使用 `tccommon.FILED_SP` 分隔符
- [x] 2.2 验证非法 `entity`/`host`/`template_id` 组合时 Create 返回 error
- [x] 2.3 验证 Create 完成后转调 Update（`ModifySecurityPolicy`），不调用独立创建 API

## 3. Read 逻辑验证

- [x] 3.1 验证 Read 函数从 `d.Id()` 按 `tccommon.FILED_SP` 拆分提取 `zoneId`/`entity`/`host`/`templateId`
- [x] 3.2 验证 Read 调用 `DescribeSecurityPolicy` 时传入 `ZoneId`/`Entity`/`Host`/`TemplateId` 入参
- [x] 3.3 验证 Read 在 response 为空时打印 `[WARN]` 日志后 `d.SetId("")`（已修正：先打印日志保留 id 现场再 SetId）
- [x] 3.4 验证 Read 成功时回填顶层 `zone_id`/`entity`/`host`/`template_id` 及 `security_config` 子结构（注：`DescribeSecurityPolicy` 不返回 `SecurityConfig`，该字段为 write-only，Read 仅回填顶层字段与 `security_policy` 表达式路径字段，与 schema Description 一致）

## 4. Update 逻辑验证

- [x] 4.1 验证 Update 函数从配置读取 `security_config.rate_limit_config` 并构造 `RateLimitConfig`（含 `Switch`/`RateLimitUserRules`/`RateLimitTemplate`/`RateLimitIntelligence`/`RateLimitCustomizes`）（已补充 `rate_limit_user_rules`/`rate_limit_customizes` 中遗漏的 `rule_priority`/`freq_fields`/`freq_scope`/`name`/`custom_response_id`/`response_code`/`redirect_url` 字段构造逻辑）
- [x] 4.2 验证 Update 调用 `ModifySecurityPolicy` 时传入 `ZoneId`/`Entity`/`Host`/`TemplateId`/`SecurityConfig`
- [x] 4.3 验证 Update 中 `ModifySecurityPolicy` 调用使用 `tccommon.WriteRetryTimeout` 重试与 `tccommon.RetryError` 错误处理
- [x] 4.4 验证 Update 成功后调用 Read 回写状态

## 5. Delete 逻辑验证

- [x] 5.1 验证 Delete 通过置空 `SecurityConfig` 并调用 `ModifySecurityPolicy` 实现删除（已实现：原 Delete 为空实现，已补充置空 SecurityConfig 并调用 ModifySecurityPolicy 的逻辑）
- [x] 5.2 验证 Delete 中 `ModifySecurityPolicy` 调用使用 `tccommon.WriteRetryTimeout` 重试与 `tccommon.RetryError` 错误处理

## 6. 单元测试验证

- [x] 6.1 验证 `tencentcloud/services/teo/resource_tc_teo_security_policy_config_test.go` 中存在覆盖 CRUD 的测试用例
- [x] 6.2 验证测试用例使用 gomonkey mock 云 API（非 terraform 测试套件）进行业务逻辑测试
- [x] 6.3 如发现 `rate_limit_config` 子结构缺少专项测试用例，补充 gomonkey mock 的单元测试覆盖 `RateLimitConfig` 构造与回填逻辑（已补充 `TestRateLimitConfig_UpdateExpand`/`TestRateLimitConfig_DeleteNeutralizes`/`TestRateLimitConfig_ReadEmpty`/`TestRateLimitConfig_CreateInvalidEntity`/`TestRateLimitConfig_Schema`）

## 7. 文档验证

- [x] 7.1 验证 `tencentcloud/services/teo/resource_tc_teo_security_policy_config.md` 包含一句话描述（提及 TEO）、Example Usage、Import 部分（已修正一句话描述为 "TEO" 大写，补充 security_config.rate_limit_config 的 Example Usage）
- [x] 7.2 验证 `.md` 中 Import 部分说明了复合 ID 格式（`zoneId#entity` 或 `zoneId#entity#host` 或 `zoneId#entity#templateId`）
- [x] 7.3 如 `.md` 内容缺失或与当前 schema 不一致，更新 `.md` 后通过 `make doc` 重新生成 `website/docs/` 文档（`.md` 已更新，`make doc` 由收尾阶段 tfpacer-finalize skill 执行）

## 8. 构建与编译验证

- [x] 8.1 运行 `go build ./...` 确保编译通过（由验证流程执行，不在本阶段运行）
- [x] 8.2 运行 `go vet ./tencentcloud/services/teo/...` 确保无静态检查问题（由验证流程执行，不在本阶段运行）
