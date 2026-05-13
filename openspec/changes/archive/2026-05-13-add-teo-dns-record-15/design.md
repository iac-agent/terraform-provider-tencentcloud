## Context

Terraform Provider for TencentCloud 需要为云产品 TEO (EdgeOne) 新增 DNS 记录管理资源 `tencentcloud_teo_dns_record_15`。目前已存在 `tencentcloud_teo_dns_record` 资源，本次新增为同一 API 的另一个版本资源，采用相同的设计模式和 API 接口。

该资源使用以下云 API 接口（均来自 `teo v20220901` SDK 包）：
- **CreateDnsRecord**: 创建单条 DNS 记录，返回 RecordId
- **DescribeDnsRecords**: 查询 DNS 记录列表，使用 AdvancedFilter 进行过滤，支持分页
- **ModifyDnsRecords**: 批量修改 DNS 记录，通过 DnsRecord 结构体传入修改内容
- **DeleteDnsRecords**: 批量删除 DNS 记录，通过 RecordIds 指定要删除的记录

所有接口均为同步接口，无需异步轮询。

当前代码库中已有 `tencentcloud_teo_dns_record` 资源作为参考实现，新资源将遵循完全相同的设计模式。

## Goals / Non-Goals

**Goals:**
- 实现 `tencentcloud_teo_dns_record_15` 资源的完整 CRUD 功能
- 遵循项目现有的代码模式（参考 `tencentcloud_teo_dns_record` 和 `tencentcloud_igtm_strategy`）
- 使用 zone_id + record_id 作为复合 ID（以 FILED_SP 分隔）
- 支持 Terraform Import
- 在 provider.go 和 provider.md 中正确注册资源
- 编写使用 gomonkey mock 的单元测试
- 生成 .md 文档文件

**Non-Goals:**
- 不修改现有的 `tencentcloud_teo_dns_record` 资源
- 不实现数据源（datasource）
- 不处理异步操作（所有 API 均为同步）
- 不修改 website/ 目录下的文件（由 make doc 生成）

## Decisions

### 1. 资源 Schema 设计
**决策**: 资源 schema 包含以下字段：
- `zone_id` (Required, ForceNew): 站点 ID
- `name` (Required): DNS 记录名称
- `type` (Required): DNS 记录类型（A/AAAA/MX/CNAME/TXT/NS/CAA/SRV）
- `content` (Required): DNS 记录内容
- `location` (Optional, Computed): 解析线路，默认 "Default"
- `ttl` (Optional, Computed): 缓存时间，默认 300
- `weight` (Optional, Computed): 权重，默认 -1
- `priority` (Optional, Computed): MX 优先级，默认 0
- `record_id` (Computed): DNS 记录 ID（由 CreateDnsRecord 返回）
- `status` (Computed): 记录状态
- `created_on` (Computed): 创建时间
- `modified_on` (Computed): 修改时间

**理由**: 完全对齐 CreateDnsRecord 入参和 DescribeDnsRecords 返回的 DnsRecord 结构体字段。zone_id 设为 ForceNew 因为 CreateDnsRecord 的 ZoneId 是创建时指定、不可更改的。

### 2. 复合 ID 格式
**决策**: 使用 `zone_id + FILED_SP + record_id` 作为复合 ID

**理由**: 与现有 `tencentcloud_teo_dns_record` 资源保持一致，因为 DescribeDnsRecords 需要 ZoneId 作为必要参数来查询记录。

### 3. Read 方法实现
**决策**: 在服务层新增 `DescribeTeoDnsRecord15ById` 方法，使用 DescribeDnsRecords 接口并通过 AdvancedFilter 按 id 过滤查询单条记录

**理由**: 没有单条记录查询接口，只能通过列表接口过滤获取。与现有 `tencentcloud_teo_dns_record` 的实现方式一致。

### 4. Update 方法实现
**决策**: 使用 ModifyDnsRecords 接口更新记录，将可变字段（name、type、content、location、ttl、weight、priority）通过 DnsRecord 结构体传入

**理由**: ModifyDnsRecords 接受 DnsRecord 数组，需要包含 RecordId 以标识要修改的记录。可变字段与 CreateDnsRecord 入参一致（不含 zone_id，因为它是 ForceNew）。

### 5. Delete 方法实现
**决策**: 使用 DeleteDnsRecords 接口，传入 RecordIds 数组删除单条记录

**理由**: DeleteDnsRecords 接受 RecordIds 字符串数组，直接传入从复合 ID 解析出的 recordId 即可。

### 6. 测试策略
**决策**: 使用 gomonkey 进行 mock 测试，不使用 Terraform 测试套件

**理由**: 按照项目规范，新增资源应使用 gomonkey mock 云 API，只测试业务逻辑。

## Risks / Trade-offs

- **[风险] DescribeDnsRecords 返回列表而非单条记录** → 通过 AdvancedFilter 按 id 精确过滤，并取第一条结果，降低误匹配风险
- **[风险] ModifyDnsRecords 是批量接口但只修改一条记录** → 只传入包含单个 DnsRecord 的数组，符合 API 规范（最大支持 100 条）
- **[风险] 与现有 tencentcloud_teo_dns_record 资源功能重复** → 这是需求要求，新资源作为独立版本存在
