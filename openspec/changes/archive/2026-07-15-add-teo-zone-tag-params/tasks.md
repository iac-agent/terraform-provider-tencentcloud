## 1. Schema 定义

- [x] 1.1 在 `tencentcloud/services/teo/resource_tc_teo_zone.go` 的 Schema 中添加 `ResourceRegion` 字段：TypeString, Optional, Computed, Description 为 "资源所在地域，用于标签操作。不区分地域的资源可忽略该参数。默认使用 provider 配置的地域。"
- [x] 1.2 在 `tencentcloud/services/teo/resource_tc_teo_zone.go` 的 Schema 中添加 `ServiceType` 字段：TypeString, Optional, Computed, Description 为 "业务类型，用于标签操作。默认为 teo。"

## 2. CRUD 函数修改

- [x] 2.1 修改 `resourceTencentCloudTeoZoneCreate`：从 schema 读取 `ResourceRegion` 和 `ServiceType`，用于构建 QCS resource name（替换硬编码的 `"teo"` 和 `tcClient.Region`），传递给 `tagService.ModifyTags` 调用
- [x] 2.2 修改 `resourceTencentCloudTeoZoneRead`：从 schema 读取 `ResourceRegion` 和 `ServiceType`，传递给 `tagService.DescribeResourceTags` 调用（替换硬编码的 `"teo"` 和 `tcClient.Region`）；在 Read 函数末尾将 `ResourceRegion` 和 `ServiceType` 的有效值设置到 state 中
- [x] 2.3 修改 `resourceTencentCloudTeoZoneUpdate`：从 schema 读取 `ResourceRegion` 和 `ServiceType`，用于构建 QCS resource name（替换硬编码的 `"teo"` 和 `tcClient.Region`），传递给 `tagService.ModifyTags` 调用

## 3. 单元测试

- [x] 3.1 在 `tencentcloud/services/teo/resource_tc_teo_zone_tag_params_test.go` 中补充 `ResourceRegion` 和 `ServiceType` 参数相关的单元测试用例，使用 mock（gomonkey）方式测试 CRUD 函数中 tag 操作对 `ResourceRegion` 和 `ServiceType` 的处理逻辑

## 4. 文档更新

- [x] 4.1 更新 `tencentcloud/services/teo/resource_tc_teo_zone.md` 文件，在 Example Usage 中补充 `ResourceRegion` 和 `ServiceType` 参数的使用示例
