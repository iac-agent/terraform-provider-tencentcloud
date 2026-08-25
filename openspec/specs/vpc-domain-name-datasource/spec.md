# vpc-domain-name-datasource Specification

## Purpose
TBD - created by syncing change add-vpc-domain-name-datasource. Update Purpose after archive.
## Requirements
### Requirement: tencentcloud_vpc 数据源暴露 DomainName 出参

`tencentcloud_vpc` 数据源 SHALL 在其 schema 中提供 `domain_name` 字段（Computed，类型为 string），用于暴露腾讯云 VPC API `DescribeVpcs` 响应中 `VpcSet[].DomainName`（DHCP 域名选项值）。

数据源的 Read 操作 SHALL 通过 `VpcService.DescribeVpcs` 获取 `VpcBasicInfo`，并将 `VpcBasicInfo.domainName` 写入 `domain_name` 字段。服务层 `DescribeVpcs` SHALL 在解析 `VpcSet` 元素时，仅当 `item.DomainName != nil` 才将其值赋给 `VpcBasicInfo.domainName`，否则保留空字符串零值。

#### Scenario: 查询到的 VPC 配置了 DHCP 域名

- **WHEN** 用户通过 `tencentcloud_vpc` 数据源查询一个已配置 DHCP 域名（例如 `"example.com"`）的 VPC
- **THEN** 数据源 Read 操作调用 `DescribeVpcs` 成功返回，且 Terraform state 中 `domain_name` 字段的值等于云 API 返回的 `VpcSet[0].DomainName`（即 `"example.com"`）

#### Scenario: 查询到的 VPC 未配置 DHCP 域名

- **WHEN** 用户通过 `tencentcloud_vpc` 数据源查询一个未配置 DHCP 域名的 VPC，云 API 返回的 `VpcSet[0].DomainName` 为空字符串或 nil
- **THEN** 数据源 Read 操作调用 `DescribeVpcs` 成功返回，服务层对 nil 指针判空跳过赋值，Terraform state 中 `domain_name` 字段的值为空字符串，且整个 Read 操作不发生 panic

#### Scenario: 数据源保留向后兼容

- **WHEN** 用户使用已有的 `tencentcloud_vpc` 数据源配置（未引用 `domain_name`）执行 plan/apply
- **THEN** 现有字段（`id`、`name`、`cidr_block`、`is_default`、`is_multicast`）的行为和值保持不变，新增的 `domain_name` 作为 Computed 字段被填充，不产生破坏性变更
