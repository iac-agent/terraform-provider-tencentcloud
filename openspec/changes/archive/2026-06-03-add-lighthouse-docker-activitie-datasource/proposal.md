## Why

Lighthouse 产品当前缺少对 Docker 活动查询的数据源支持。用户在管理轻量应用服务器上的 Docker 容器时，需要查看 Docker 活动列表（如容器创建、启动、停止等操作的执行状态和输出），以便监控和排查容器操作的结果。通过新增 `tencentcloud_lighthouse_docker_activitie` 数据源，用户可以在 Terraform 中查询 Docker 活动信息。

## What Changes

- 新增数据源 `tencentcloud_lighthouse_docker_activitie`，支持通过 `DescribeDockerActivities` 云 API 接口查询 Docker 活动列表
- 支持按实例 ID、活动 ID 列表、创建时间范围等条件过滤查询
- 在 provider.go 和 provider.md 中注册该数据源
- 新增数据源文档 data_source_tc_lighthouse_docker_activitie.md

## Capabilities

### New Capabilities
- `lighthouse-docker-activitie-datasource`: 新增 Lighthouse Docker 活动数据源，支持查询 Docker 活动列表，包含活动 ID、活动名称、活动状态、命令输出、容器 ID 列表、创建时间和结束时间等信息

### Modified Capabilities

## Impact

- 新增文件: `tencentcloud/services/lighthouse/data_source_tc_lighthouse_docker_activitie.go`
- 新增测试文件: `tencentcloud/services/lighthouse/data_source_tc_lighthouse_docker_activitie_test.go`
- 新增文档文件: `tencentcloud/services/lighthouse/data_source_tc_lighthouse_docker_activitie.md`
- 修改文件: `tencentcloud/provider.go`（注册数据源）、`tencentcloud/provider.md`（添加数据源条目）
- 依赖云 API: `DescribeDockerActivities`（lighthouse v20200324）
