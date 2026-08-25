## Why

`tencentcloud_vpc` 数据源（`data_source_tc_vpc.go`）用于读取单个 VPC 的属性，但目前未暴露 `DomainName`（DHCP 域名选项值）字段。腾讯云 VPC API `DescribeVpcs` 的响应 `VpcSet[].DomainName` 已经返回了该值，但数据源没有将其映射到 Terraform schema，导致用户无法通过该数据源获取 VPC 的 DHCP 域名配置信息。

## What Changes

- 为数据源 `tencentcloud_vpc` 新增出参 `domain_name`（Computed），用于暴露 `DescribeVpcs` 接口响应中 `VpcSet[].DomainName` 的 DHCP 域名选项值。

## Capabilities

### New Capabilities
- `vpc-domain-name-datasource`: 为 `tencentcloud_vpc` 数据源新增 `domain_name` 出参，映射 `DescribeVpcs` 响应中 `VpcSet[].DomainName` 字段。

### Modified Capabilities
<!-- No existing specs require modification -->

## Impact

- **Affected files:**
  - `tencentcloud/services/vpc/data_source_tc_vpc.go` — 在 `Schema` 中新增 `domain_name`（Computed）字段，并在 `dataSourceTencentCloudVpcRead` 中调用 `d.Set("domain_name", ...)`。
  - `tencentcloud/services/vpc/service_tencentcloud_vpc.go` — 在 `VpcBasicInfo` 结构体中新增 `domainName` 字段，并在 `DescribeVpcs` 解析 `VpcSet` 时从 `item.DomainName` 读取（需判空，避免指针解引用 panic）。
- **SDK dependency:** 现有 `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312` 的 `Vpc` 结构体已包含 `DomainName` 字段，无需更新 SDK。
- **Backward compatibility:** 完全向后兼容，仅为已废弃的数据源新增一个 Computed 出参，不修改任何现有字段行为，不影响已有 TF 配置和 state。
- **API constraints:** `DomainName` 仅出现在 `DescribeVpcs` 响应中（只读出参），无对应的写入操作，本次变更不涉及 Create/Update/Delete 流程。
