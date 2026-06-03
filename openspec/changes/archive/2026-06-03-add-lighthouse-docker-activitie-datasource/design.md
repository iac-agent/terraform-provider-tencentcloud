## Context

Terraform Provider for TencentCloud 需要新增 Lighthouse Docker 活动数据源。Lighthouse（轻量应用服务器）产品提供了 `DescribeDockerActivities` 云 API 接口用于查询 Docker 活动列表，但当前 Terraform Provider 中尚无对应的数据源支持。

当前状态：
- Lighthouse 服务已存在多个数据源（blueprints、disks、bundle 等），遵循统一的代码模式
- 服务层 `LightHouseService` 已在 `service_tencentcloud_lighthouse.go` 中实现
- `DescribeDockerActivities` API 支持 InstanceId、ActivityIds、CreatedTimeBegin、CreatedTimeEnd 过滤条件，并支持分页

## Goals / Non-Goals

**Goals:**
- 新增 `tencentcloud_lighthouse_docker_activitie` 数据源，支持通过 `DescribeDockerActivities` API 查询 Docker 活动列表
- 支持按实例 ID、活动 ID 列表、创建时间范围等条件过滤查询
- 返回完整的 Docker 活动信息（活动 ID、名称、状态、命令输出、容器 ID 列表、创建/结束时间）
- 自动处理分页，对用户透明
- 在 provider.go 和 provider.md 中注册数据源
- 新增数据源文档和单元测试

**Non-Goals:**
- 不修改现有数据源或资源的 schema
- 不新增 Docker 容器管理相关的 CRUD 资源
- 不处理异步操作轮询（本数据源仅执行只读查询）

## Decisions

### 1. 数据源命名
**Decision**: 使用 `tencentcloud_lighthouse_docker_activitie` 作为数据源名称，与云 API 接口 `DescribeDockerActivities` 保持语义一致。

**Rationale**: 遵循需求中指定的命名，与云 API 保持一致。

### 2. 代码组织模式
**Decision**: 采用与 `tencentcloud_lighthouse_blueprints` 数据源相同的模式：
- 数据源文件: `data_source_tc_lighthouse_docker_activitie.go`
- 测试文件: `data_source_tc_lighthouse_docker_activitie_test.go`（使用 gomonkey mock）
- 文档文件: `data_source_tc_lighthouse_docker_activitie.md`
- 服务层方法: `DescribeLighthouseDockerActivitiesByFilter()`

**Rationale**: 遵循项目现有 lighthouse 数据源的统一模式。

### 3. 输入参数设计
**Decision**: 将 API 的 InstanceId、ActivityIds、CreatedTimeBegin、CreatedTimeEnd 映射为数据源 schema 的 Optional 参数。Offset 和 Limit 参数不在 schema 中暴露，由服务层自动处理分页。

**Rationale**:
- InstanceId 为核心过滤条件
- ActivityIds 用于精确查询指定活动
- CreatedTimeBegin/CreatedTimeEnd 用于时间范围过滤（注意：API 中类型为 int64 时间戳秒数，schema 中也使用 TypeInt）
- 分页参数对用户隐藏，自动循环获取所有结果

### 4. 输出结构设计
**Decision**: `docker_activity_set` 作为 Computed TypeList，每个元素包含 DockerActivity 结构体的所有字段：
- activity_id (string)
- activity_name (string)
- activity_state (string)
- activity_command_output (string)
- container_ids ([]string)
- created_time (string)
- end_time (string)

**Rationale**: 完整映射 API 响应中 DockerActivity 结构体的所有字段。

### 5. CreatedTimeBegin/CreatedTimeEnd 参数类型
**Decision**: 使用 `schema.TypeInt` 类型，与云 API 的 int64 时间戳秒数保持一致。

**Rationale**: 云 API 中这两个字段为 int64 类型（时间戳秒数），直接映射为 TypeInt 避免类型转换问题。

## Risks / Trade-offs

- **[API 字段可空性]** → 所有输出字段进行 nil 检查，避免空指针异常
- **[API 返回空结果]** → 正常处理，返回空列表
- **[分页边界情况]** → 服务层使用 pageSize=100（API 最大值），逐页获取直到结果不足一页
