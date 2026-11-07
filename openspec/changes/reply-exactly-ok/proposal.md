## Why

任务给出的是需求描述 `reply exactly: OK`。它描述的期望行为是：当收到该指令时，回复内容必须**精确等于** `OK`（两个大写字母，不带任何前后缀、解释、标点或额外文本），而不是一个可以落到腾讯云云 API 上的基础设施能力。

在把它转成 OpenSpec 变更之前，先做了一次「提案正确性 / 可行性」校验，结论如下：

- 在 `vendor/github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/`（99 个服务包）中检索 `reply exactly` 等关键字，**没有任何云 API 提供"按指令精确回复文本"的能力**。
- Terraform Provider 的职责是把云资源状态映射为 IaC（Create / Read / Update / Delete / DataSource），而"精确回复 OK"既不对应任何资源、数据源或 Action 接口，也不对应任何可持久化、可读回的状态字段。
- 因此本需求**不应**被降级实现为"新增某个资源/数据源字段"来假装满足；那会引入无意义的 schema 与永久 `Diff`。

本变更的目的因此被收敛为：**在仓库层面显式登记这条"精确回复"约定**，把它作为一条可被后续自动化流程校验的行为契约记录下来，并**明确声明不做任何云 API 适配、不改动任何资源 schema**，避免后续执行者误以为需要新增资源或参数。

## What Changes

- 新增一条 OpenSpec capability `exact-reply-contract`，用可被机器/人复核的 Requirement + Scenario 形式描述"收到 `reply exactly: OK` 时必须且仅输出 `OK`"。
- 明确 **Non-Goal（不做的事）**，防止越界实现：
  - 不在 `vendor/` 中新增或修改任何云 API SDK 代码；
  - 不在 `tencentcloud/services/` 下新增任何资源 / 数据源 / service 层代码；
  - 不修改 `tencentcloud/provider.go`、`tencentcloud/provider.md` 的资源注册列表；
  - 不新增 `.changelog/` 条目，不做 `make doc` 文档生成（无 schema 变更可生成）；
  - 不新增 HCL 示例、`_test.go` 验收用例或 gomonkey mock（无代码可测）。
- 变更产物仅为 OpenSpec 提案文件本身：`proposal.md`、`design.md`、`tasks.md`、`specs/exact-reply-contract/spec.md`。

### 云 API 可行性校验结论

| 校验项 | 结论 | 依据 |
| --- | --- | --- |
| 是否存在对应云 API | 否 | `vendor/` 下 99 个 SDK 服务包全量检索无匹配；不存在 "ReplyExactly" 类接口 |
| 是否可映射为资源 | 否 | 无 CRUD 语义，无可持久化 ID |
| 是否可映射为数据源 | 否 | 无可查询的远端集合，无入参/出参模型 |
| 是否需要改 Schema | 否 | 无任何资源/数据源字段涉及 |
| 是否需要改 provider 注册 | 否 | 无新增资源/数据源 |
| 是否需要 changelog / 文档 | 否 | 无面向用户的 resource/data source 变更 |

**结论：本提案的"正确性、可行性"体现在它只登记行为契约、显式拒绝 API 适配，不产出任何无法由云 API 支撑的 Terraform 能力。**

## Capabilities

### New Capabilities

- `exact-reply-contract`: 登记"收到 `reply exactly: OK` 时必须精确输出 `OK`"这一行为契约，并固化其边界（不涉及云 API、不涉及 Terraform schema）。

### Modified Capabilities
<!-- openspec/specs/ 下不存在与本需求相关的既有 capability，故无 Modified Capabilities。 -->

## Impact

- **代码**：无。`tencentcloud/`、`vendor/`、`website/` 均零改动。
- **依赖**：无。`go.mod` / `go.sum` / `vendor/modules.txt` 不变。
- **云 API**：无调用。不依赖任何腾讯云接口的可用性。
- **文档生成**：无。不触发 `make doc`。
- **向后兼容性**：完全兼容。现有 Terraform 配置、state、Provider schema 均不受影响。
- **风险**：唯一风险是后续执行者把该契约误解为"需要实现某个资源"，从而产出无 API 支撑的伪实现。本提案的 capability 规范中已用显式 Scenario 禁止该行为。
