## Context

`tencentcloud_vpc` 数据源（`tencentcloud/services/vpc/data_source_tc_vpc.go`）用于读取单个 VPC 的属性。它通过 `VpcService.DescribeVpcs` 调用腾讯云 VPC API `DescribeVpcs`，取返回列表的第一项，将字段写入 Terraform state。

当前数据源已暴露 `name`、`cidr_block`、`is_default`、`is_multicast` 字段，但未暴露 `DomainName`（DHCP 域名选项值）。云 API `DescribeVpcs` 响应中的 `VpcSet[].DomainName` 已经返回了该值（SDK `Vpc` 结构体已包含 `DomainName *string` 字段），但服务层 `VpcBasicInfo` 未携带该字段，数据源 schema 也未映射。

该数据源已在 1.10.0 版本标记为 deprecated（建议改用 `tencentcloud_vpc_instances`），但仍需保持功能完整性。

## Goals / Non-Goals

**Goals:**
- 在 `tencentcloud_vpc` 数据源中新增 `domain_name` 出参（Computed），暴露 `DescribeVpcs` 响应中 `VpcSet[].DomainName` 的值。
- 保持完全向后兼容：不修改任何现有字段的行为，不影响已有 TF 配置和 state。

**Non-Goals:**
- 不新增任何入参（无查询过滤条件变化）。
- 不修改 `tencentcloud_vpc_instances` 数据源或其他 VPC 资源。
- 不涉及 Create/Update/Delete 流程（`DomainName` 仅作为只读出参）。
- 不更新 SDK（现有 vendor 中 `Vpc` 结构体已包含 `DomainName` 字段）。
- 不取消该数据源的 deprecated 状态。

## Decisions

**决策 1：通过 `VpcBasicInfo` 中转字段，而非直接在数据源读取 SDK 响应。**

数据源 `dataSourceTencentCloudVpcRead` 并不直接持有 `*vpc.DescribeVpcsResponse`，而是调用 `VpcService.DescribeVpcs(...)` 得到 `[]VpcBasicInfo` 并取首元素。因此需要在 `VpcBasicInfo` 中新增 `domainName string` 字段，并在 `DescribeVpcs` 的解析循环中填充。

- 备选方案：重构 `DescribeVpcs` 让数据源直接访问原始 SDK 响应 —— 拒绝，因为 `DescribeVpcs` 被多个资源/数据源共用，重构影响面过大，违反最小变更原则。

**决策 2：对 `item.DomainName` 进行判空后再解引用。**

SDK 中 `DomainName` 类型为 `*string`（`omitnil`）。现有 `DescribeVpcs` 解析代码对 `CidrBlock`、`CreatedTime` 等字段直接解引用（`*item.CidrBlock`），依赖云 API 必返回这些字段。但 `DomainName` 仅在 VPC 开启 DHCP 且配置了域名时才可能有值，为稳妥起见，填充时需判空：当 `item.DomainName != nil` 时才赋值，否则保留零值（空字符串），避免指针解引用 panic。

**决策 3：`domain_name` 字段设为 Computed（只读出参），不设 Optional。**

`DomainName` 仅出现在 `DescribeVpcs` 响应中，无对应的写入语义（该数据源本身不含 Create/Update）。因此 schema 中 `domain_name` 仅声明 `Computed: true`，与同类的 `cidr_block`、`is_default`、`is_multicast` 保持一致。

## Risks / Trade-offs

- **[数据源已废弃]** → 用户被引导使用 `tencentcloud_vpc_instances`。本次仅做最小增强，不为已废弃数据源投入额外重构，风险可控。
- **[DomainName 可能为空]** → 部分未配置 DHCP 域名的 VPC 返回空字符串。数据源直接 set 空字符串即可，`d.Set` 对空值是合法操作，无副作用。
- **[指针解引用风险]** → 通过决策 2 的判空处理规避。
