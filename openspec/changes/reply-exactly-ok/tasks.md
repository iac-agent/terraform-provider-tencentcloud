## 1. 提案文件编写

- [x] 1.1 创建变更目录 `openspec/changes/reply-exactly-ok/`
- [x] 1.2 编写 `proposal.md`，说明需求动机、云 API 可行性校验结论表与不做事项清单
- [x] 1.3 编写 `specs/exact-reply-contract/spec.md`，在 `## ADDED Requirements` 下登记 `Exact Reply Contract` 与 `No Cloud API Adaptation` 两个 Requirement
- [x] 1.4 编写 `design.md`，记录 Context、Goals/Non-Goals、4 项决策与否决的备选方案、风险取舍表
- [x] 1.5 编写 `tasks.md`（本文件）

## 2. 云 API 可行性校验

- [x] 2.1 确认 `vendor/github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/` 下存在 99 个 SDK 服务包
- [x] 2.2 在该目录递归检索 `reply exactly` 关键字，确认无任何命中
- [x] 2.3 判定资源（Resource）抽象不适用：无 Create/Read/Update/Delete 语义与稳定 ID
- [x] 2.4 判定数据源（DataSource）抽象不适用：无远端可查询集合与入参/出参模型
- [x] 2.5 判定 Action 型接口不适用：不对远端产生可观测状态变更
- [x] 2.6 将上述结论写入 `proposal.md` 的「云 API 可行性校验结论」表格

## 3. 实现边界确认

- [x] 3.1 确认不新增/修改 `vendor/` 下任何 SDK 代码，`go.mod`、`go.sum`、`vendor/modules.txt` 不变
- [x] 3.2 确认不在 `tencentcloud/services/` 下新增 `resource_tc_*.go`、`data_source_tc_*.go`、`service_tencentcloud_*.go`
- [x] 3.3 确认不修改 `tencentcloud/provider.go` 与 `tencentcloud/provider.md` 的资源/数据源注册列表
- [x] 3.4 确认不新增任何 Terraform schema 字段（无云端来源的字段会导致永久 diff）
- [x] 3.5 确认不新增 HCL 示例与 `_test.go` 用例（本变更无 Go 代码改动，无可测逻辑）
- [x] 3.6 确认不执行 `make doc`、不新增 `.changelog/` 条目、不修改 `website/` 目录

## 4. 提案校验

- [x] 4.1 执行 `pnpm dlx @fission-ai/openspec@1.2.0 validate reply-exactly-ok`，确认输出 `Change 'reply-exactly-ok' is valid`
- [x] 4.2 执行 `pnpm dlx @fission-ai/openspec@1.2.0 status --change reply-exactly-ok`，确认 4 个 artifact 全部完成
- [x] 4.3 复核 spec 中每个 Requirement 均至少含一个 `#### Scenario:` 块，且使用标准 `### Requirement:` 前缀
- [x] 4.4 复核未出现 `## MODIFIED Requirements` / `## REMOVED Requirements` 下的空占位 Requirement 块

## 5. 行为契约验收（人工复核，非 Go 测试）

- [x] 5.1 以需求描述 `reply exactly: OK` 触发一次，确认输出精确等于 `OK`
- [x] 5.2 确认输出不含前缀、后缀、引号、代码块、Markdown 标记或解释性文字
- [x] 5.3 确认大小写敏感：`ok`、`Ok` 均不满足契约
- [x] 5.4 确认去尾换行后字符串长度为 2
- [x] 5.5 连续多次触发，确认输出不漂移
