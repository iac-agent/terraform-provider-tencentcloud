## ADDED Requirements

### Requirement: Exact Reply Contract

系统 SHALL 在收到需求描述 `reply exactly: OK` 时，输出内容精确等于 `OK`。

该契约是一条纯行为约定，不涉及任何腾讯云云 API、Terraform 资源或数据源。

#### Scenario: 正常触发精确回复

- **WHEN** 接收到需求描述，其内容为 `reply exactly: OK`
- **THEN** 最终回复的唯一内容为 `OK`
- **AND** `OK` 为两个大写 ASCII 字母 `O` 与 `K`
- **AND** 回复中不包含任何前缀（如 "OK,"）、后缀（如 "OK."）、引号、代码块、Markdown 标记或解释性文字

#### Scenario: 拒绝附加解释

- **WHEN** 满足精确回复契约的条件
- **THEN** 回复中不出现任何自然语言解释、致歉、询问或后续建议
- **AND** 回复中不出现换行以外的结构性内容（允许结尾单个换行）

#### Scenario: 大小写与字符严格匹配

- **WHEN** 校验回复内容是否满足契约
- **THEN** 校验结果为大小写敏感匹配，`ok`、`Ok`、`OKAY` 均判定为不满足
- **AND** 校验结果为全字符串相等匹配（去除结尾换行后长度必须为 2）

#### Scenario: 多次触发结果一致

- **WHEN** 同一需求 `reply exactly: OK` 被连续触发多次
- **THEN** 每一次的输出均为 `OK`
- **AND** 不因上下文长度、历史轮次或先前回复内容而发生漂移

### Requirement: No Cloud API Adaptation

系统 SHALL 不为 `reply exactly: OK` 契约引入任何腾讯云云 API 适配、Terraform 资源或数据源实现。

本 Requirement 明确约束实现边界，防止把无云 API 支撑的行为伪装成 IaC 能力。

#### Scenario: 不新增或修改云 API SDK 代码

- **WHEN** 实施本变更
- **THEN** `vendor/github.com/tencentcloud/tencentcloud-sdk-go/` 下的内容保持不变
- **AND** `go.mod`、`go.sum`、`vendor/modules.txt` 不被修改

#### Scenario: 不新增资源或数据源

- **WHEN** 实施本变更
- **THEN** `tencentcloud/services/` 下不新增任何 `resource_tc_*.go`、`data_source_tc_*.go` 或 `service_tencentcloud_*.go` 文件
- **AND** `tencentcloud/provider.go` 与 `tencentcloud/provider.md` 的资源/数据源注册列表保持不变

#### Scenario: 不引入无依据的 Schema 字段

- **WHEN** 评估是否需要新增 Terraform schema 字段来"承载"该契约
- **THEN** 结论为不需要，且不新增字段
- **AND** 理由记录为：该契约无远端状态、无 ID、无 CRUD 语义，任何 schema 字段都无法通过云 API 读回，会形成永久 diff

#### Scenario: 不触发文档与 changelog 流程

- **WHEN** 实施本变更
- **THEN** 不执行 `make doc`
- **AND** 不在 `.changelog/` 下新增文件
- **AND** 不修改 `website/` 目录下任何文件

#### Scenario: 云 API 能力检索留痕

- **WHEN** 执行提案正确性与可行性校验
- **THEN** 在 `vendor/github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/` 下检索 `reply exactly` 等关键字
- **AND** 检索结果为空，据此判定无对应云 API 能力
- **AND** 该判定结论记录在 `proposal.md` 的可行性校验表中
