## 1. 服务层改造（service_tencentcloud_vpc.go）

- [x] 1.1 在 `VpcBasicInfo` 结构体中新增私有字段 `domainName string`，用于中转 `DescribeVpcs` 响应中的 `DomainName` 值。
- [x] 1.2 在 `DescribeVpcs` 方法解析 `response.Response.VpcSet` 的循环中，新增对 `item.DomainName` 的判空读取：仅当 `item.DomainName != nil` 时将 `*item.DomainName` 赋给 `basicInfo.domainName`，否则保留空字符串零值。

## 2. 数据源改造（data_source_tc_vpc.go）

- [x] 2.1 在 `DataSourceTencentCloudVpc()` 的 `Schema` map 中新增 `domain_name` 字段：`Type: schema.TypeString`、`Computed: true`，并补充 `Description`（说明为 DHCP 域名选项值）。保持与同类出参（`cidr_block`、`is_default`、`is_multicast`）风格一致。
- [x] 2.2 在 `dataSourceTencentCloudVpcRead` 中，于取首元素 `vpc := vpcInfos[0]` 之后，新增 `_ = d.Set("domain_name", vpc.domainName)`（注意 `VpcBasicInfo` 字段为私有，需在同包内直接访问；若需暴露则新增 getter，默认同包直接访问）。

## 3. 单元测试（data_source_tc_vpc_test.go）

- [x] 3.1 使用 mock（gomonkey）方式对 `VpcService.DescribeVpcs` 进行 mock，补充单元测试用例，覆盖：查询到配置了 DHCP 域名的 VPC（`domain_name` 被正确填充）以及未配置域名（`domain_name` 为空字符串）两种场景，断言 `domain_name` 字段值符合预期。禁止使用 terraform 测试套件。

## 4. 文档与收尾

- [x] 4.1 更新数据源示例文档 `tencentcloud/services/vpc/data_source_tc_vpc.md`，在 Example Usage 中体现 `domain_name` 出参（无需手动添加 Argument/Attribute Reference，由 `make doc` 自动生成）。
- [x] 4.2 在收尾阶段通过 `make doc` 生成 `website/docs/` 下的文档，并通过 `gofmt` 格式化变更的 Go 文件、生成 changelog，统一由 tfpacer-finalize skill 执行。
