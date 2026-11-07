## Context

需求描述只有一句：`reply exactly: OK`。它不是云资源建模问题，而是一个**输出格式约束**。

本仓库是 Terraform Provider for TencentCloud，其所有交付物都必须能落到"云 API 调用 → Terraform schema → state 读写"这条链路上。因此需要先判定：这条需求在链路上是否有落点。

判定过程：

1. 在 `vendor/github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/` 下（99 个服务包）检索 `reply exactly` 关键字，无任何命中。
2. 逐一排除三类 Terraform 抽象与云 API 的对应关系：
   - 资源（Resource）：需要 Create/Read/Update/Delete 四类语义与稳定 ID。本需求无。
   - 数据源（DataSource）：需要远端可查询集合与入参/出参模型。本需求无。
   - Action 型接口（如 `StartXxxOperation`）：需要对远端产生可观测状态变更。本需求无。
3. 结论：**该需求没有云 API 落点**，因此不能也不应该被实现为资源、数据源或 schema 字段。

现有仓库中已有同类"非 API 变更"的先例（如 `openspec/changes/add-plugin-framework-muxing`），其做法是登记 capability 规范而不触碰 service 层代码。本设计沿用该模式。

## Goals / Non-Goals

**Goals:**

- 用 OpenSpec 形式登记 `reply exactly: OK` 这一行为契约，使其可被后续流程复核，而不是只存在于自然语言对话里。
- 把契约的判定条件写死为可机械验证的形式：大小写敏感、全字符串相等、去尾换行后长度为 2。
- 显式声明实现边界，避免后续执行者把"精确回复"误读为"需要新增某个资源/参数去承载它"。

**Non-Goals:**

- 不新增、不修改 `vendor/` 下任何云 API SDK 代码或依赖版本。
- 不在 `tencentcloud/services/` 下新增资源、数据源、service 层文件。
- 不修改 `tencentcloud/provider.go` / `tencentcloud/provider.md` 的资源注册列表。
- 不新增 HCL 示例、不新增 `_test.go` 用例（无论验收套件还是 gomonkey mock），因为没有任何 Go 代码被改动。
- 不执行 `make doc`，不新增 `.changelog/` 条目，不触碰 `website/` 目录。
- 不引入任何"为了满足需求而凭空造出来的" schema 字段。

## Decisions

### 决策 1：不做云 API 适配，只登记行为契约

**决策**：本变更不产出任何 Go 代码或 Terraform schema，仅产出 OpenSpec 文档（proposal/design/tasks/spec）。

**理由**：

- 云 API 侧无对应能力（已在 Context 中检索验证）。任何"实现"都只能是对本地逻辑的伪包装，无法通过 Read 从云端读回，会形成永久 diff。
- Terraform Provider 的契约是"state 反映云端真实状态"。引入一个无云端来源的字段会破坏这一契约，属于比不实现更差的结果。

**备选方案与否决理由**：

- 备选 A：新增一个 `tencentcloud_xxx` 资源，其 Create 只返回固定字符串。否决——无云端状态可读，Read 必然清空 ID，形成 create/destroy 循环。
- 备选 B：在某个现有资源上新增 `reply` 字段。否决——该字段无法通过云 API 读回，且与本资源语义无关，属于污染现有 schema。
- 备选 C：在 provider `configure` 中加开关。否决——provider 层字段会进入所有配置的 plan diff，影响面与需求价值严重不匹配。

### 决策 2：capability 命名为 `exact-reply-contract`

**决策**：使用 `exact-reply-contract` 作为 capability 目录名（kebab-case）。

**理由**：命名描述的是"契约"而非"功能"，与 `openspec/config.yaml` 中"新增 spec 用 ADDED Requirements"的规则一致；同时与既有 capability（均以资源/机制命名）不冲突，避免误认为存在一个同名资源。

### 决策 3：把"不做什么"也写成 Requirement

**决策**：`specs/exact-reply-contract/spec.md` 中包含第二个 Requirement `No Cloud API Adaptation`，用 Scenario 固化禁止项（不改 vendor、不新增资源、不改 provider 注册、不跑 doc/changelog）。

**理由**：tasks.md 只约束"做过什么"，而本变更最大的风险是**过度实现**。把禁止项提升为 Requirement + Scenario，可以让后续评审/CI 以同等优先级检查"是否越界"，而不是仅在 tasks 里以勾选框形式存在。

### 决策 4：tasks 只包含文档类任务，且验证阶段只做校验类命令

**决策**：tasks 分组为「提案文件编写」「云 API 可行性校验」「边界确认」「OpenSpec 校验」四组，全部为只读或文档写入操作。

**理由**：与仓库全局约束一致——本阶段禁止执行 `go build` / `go vet` / `go test`，禁止在 `tfpacer-finalize` 之外执行 `gofmt` 与 `make doc`。`openspec validate` 属于允许的轻量校验。

## Risks / Trade-offs

| 风险 | 影响 | 缓解 |
| --- | --- | --- |
| 后续执行者把契约误读为"要实现一个资源" | 产出无云 API 支撑的伪实现，污染 provider schema | spec 中 `No Cloud API Adaptation` 的 5 个 Scenario 逐项禁止；proposal 的可行性校验表留痕 |
| 把 `OK` 写成 `ok`/`OK.`/带解释的长文本 | 契约不满足 | 规范中显式定义大小写敏感 + 全字符串相等 + 去尾换行长度 2 |
| 有人要求"顺便加个 changelog / 文档" | 触发无意义的 `make doc` 与 `.changelog` 条目 | tasks 与 spec 均声明不执行；同时符合仓库"禁止在收尾阶段外新增 changelog"的约束 |
| 提案文件过多但无代码 | 看起来"空转" | 这是刻意选择：显式登记边界比默默跳过更可审计；文件总量控制在 4 个 |

**取舍**：本变更放弃了"产出一个看起来能跑的东西"的短期满足感，换取"不引入无法由云 API 支撑的 Terraform 能力"的长期正确性。这与仓库硬约束「vendor 模式管理依赖、云 API 优先」一致。
