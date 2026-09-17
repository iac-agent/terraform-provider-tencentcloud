## 1. Schema 变更

- [x] 1.1 在 `tencentcloud/services/teo/resource_tc_teo_plan.go` 的 `ResourceTencentCloudTeoPlan()` Schema 中新增 `auto_use_voucher` 字段（TypeString、Optional），使用 `tccommon.ValidateAllowedStringValue([]string{"true", "false"})` 校验取值，Description 说明：是否自动使用代金券，取值 `true`/`false`，默认 `false`，仅在 `plan_type` 为 `personal`/`basic`/`standard`（预付费套餐）时有效，创建后不可修改
- [x] 1.2 确认无需在 Schema 中添加 Timeouts 块（现有 Create/Update/Delete 均无轮询等待逻辑，维持现状）

## 2. Create 函数变更

- [x] 2.1 在 `ResourceTencentCloudTeoPlanCreate` 中新增参数读取逻辑：`if v, ok := d.GetOk("auto_use_voucher"); ok { request.AutoUseVoucher = helper.String(v.(string)) }`，仅在用户显式配置时传入（未配置时 API 默认 `false`，保持向后兼容）
- [x] 2.2 复核 Create 函数中现有的返回值空值检查（`result == nil || result.Response == nil`、`response.Response.PlanId == nil`）保持不变，并确保在检查 id 前打印 logId 与 d.Id() 便于排障

## 3. Read 函数确认

- [x] 3.1 确认 `ResourceTencentCloudTeoPlanRead` 不设置 `auto_use_voucher`（`DescribePlans` 返回的 `Plan` 结构体无该字段，无法回填），保持现有字段回填逻辑不变

## 4. Update 函数变更

- [x] 4.1 在 `ResourceTencentCloudTeoPlanUpdate` 中新增 `immutableArgs` 数组（包含 `auto_use_voucher`），遍历数组并在 `d.HasChange(...)` 时返回 `fmt.Errorf("argument '%s' cannot be changed", ...)` 类型的错误，阻止创建后修改该参数
- [x] 4.2 确认 Update 函数中现有的 `plan_type`（UpgradePlan）、`prepaid_plan_param.0.period`（RenewPlan）、`prepaid_plan_param.0.renew_flag`（ModifyPlan）分支逻辑不受影响

## 5. 单元测试

- [x] 5.1 新建 `tencentcloud/services/teo/resource_tc_teo_plan_test.go`，使用 gomonkey mock 云API（`CreatePlanWithContext`、`DescribePlans`、`UpgradePlanWithContext`、`RenewPlanWithContext`、`ModifyPlanWithContext`、`DestroyPlanWithContext`），不使用 terraform 测试套件
- [x] 5.2 编写用例：Create 时配置 `auto_use_voucher = "true"`，断言 `CreatePlan` 请求中包含 `AutoUseVoucher`
- [x] 5.3 编写用例：Create 时未配置 `auto_use_voucher`，断言 `CreatePlan` 请求中 `AutoUseVoucher` 为 nil
- [x] 5.4 编写用例：Read 正常回填 Computed 字段（mock `DescribePlans` 返回 Plan 数据）
- [x] 5.5 编写用例：Update 中修改 `auto_use_voucher` 时返回 immutable 错误且不调用任何 Update API
- [x] 5.6 编写用例：Update 中仅修改 `plan_type` 时正常调用 `UpgradePlan`（auto_use_voucher 未变化不报错）
- [x] 5.7 编写用例：Delete 正常调用 `DestroyPlan` 并清理资源

## 6. 文档同步

- [x] 6.1 新建 `tencentcloud/services/teo/resource_tc_teo_plan.md`：一句话描述（带上 TEO 产品名，格式 "Provides a resource to ..."）、Example Usage（示例中包含 `auto_use_voucher` 参数）、Import 部分（说明使用 plan id 导入）；不手动添加 `Argument Reference` 和 `Attribute Reference`
- [x] 6.2 确认 `tencentcloud/provider.go` 和 `tencentcloud/provider.md` 中 `tencentcloud_teo_plan` 已注册（provider.go 注册代码已存在；provider.md 中 TEO Resource 索引列表缺失该项，已补充）

## 7. 验证（由后续流程执行，不在本阶段运行）

- [x] 7.1 代码正确性检查：确认新增参数在 `CreatePlan` 云API创建接口中存在（`request.AutoUseVoucher`），CRUD 各接口参数匹配无误
- [ ] 7.2 编译验证（`go build` 等）与 `gofmt`、`make doc`、changelog 生成统一由收尾阶段 tfpacer-finalize skill 执行
