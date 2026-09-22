## 1. Schema 定义与 CRUD 实现

- [x] 1.1 在 `tencentcloud/services/sqlserver/resource_tc_sqlserver_basic_instance.go` 中添加 `time_zone` schema 字段（TypeString, Optional, Computed, ForceNew）
- [x] 1.2 在 `tencentcloud/services/sqlserver/resource_tc_sqlserver_basic_instance.go` 中添加 `disk_encrypt_flag` schema 字段（TypeInt, Optional, Computed, ForceNew, ValidateIntegerInRange(0,1)）
- [x] 1.3 在 `resourceTencentCloudSqlserverBasicInstanceUpdate` 的 `immutableArgs` 中添加 `time_zone` 和 `disk_encrypt_flag`
- [x] 1.4 在 `resourceTencentCloudSqlserverBasicInstanceCreate` 中将 `time_zone` 和 `disk_encrypt_flag` 加入 paramMap
- [x] 1.5 在 `resourceTencentCloudSqlserverBasicInstanceRead` 中从 `DescribeDBInstances` 响应读取 `TimeZone` 并 set 到 state（带 nil 检查）
- [x] 1.6 在 `resourceTencentCloudSqlserverBasicInstanceRead` 中调用 `DescribeSqlserverInstanceAttributeById` 读取 `IsDiskEncryptFlag` 并 set 到 state（带双重 nil 检查）

## 2. Service 层实现

- [x] 2.1 在 `tencentcloud/services/sqlserver/service_tencentcloud_sqlserver.go` 中新增 `DescribeSqlserverInstanceAttributeById` 方法，调用 `DescribeDBInstancesAttribute` API
- [x] 2.2 在 `CreateSqlserverBasicInstance` 中处理 paramMap 中的 `time_zone`（`request.TimeZone = helper.String(v.(string))`）
- [x] 2.3 在 `CreateSqlserverBasicInstance` 中处理 paramMap 中的 `disk_encrypt_flag`（`request.DiskEncryptFlag = helper.IntInt64(v.(int))`）

## 3. 单元测试

- [x] 3.1 在 `tencentcloud/services/sqlserver/resource_tc_sqlserver_basic_instance_test.go` 中新增 `TestAccTencentCloudSqlserverBasicInstanceTimezone` 测试用例，验证自定义时区创建与导入
- [x] 3.2 在 `tencentcloud/services/sqlserver/resource_tc_sqlserver_basic_instance_test.go` 中新增 `TestAccTencentCloudSqlserverBasicInstanceDiskEncrypt` 测试用例，验证磁盘加密创建与导入
- [x] 3.3 在 `tencentcloud/services/sqlserver/resource_tc_sqlserver_basic_instance_test.go` 中新增 `TestAccTencentCloudSqlserverBasicInstanceTimezoneAndEncrypt` 测试用例，验证两个参数同时使用

## 4. 文档

- [x] 4.1 在 `tencentcloud/services/sqlserver/resource_tc_sqlserver_basic_instance.md` 中添加 timezone 和 disk encryption 的使用示例
- [x] 4.2 在 `tencentcloud/services/sqlserver/resource_tc_sqlserver_basic_instance.md` 中补充 Argument Reference 参数说明，包含 `time_zone` 和 `disk_encrypt_flag`
- [x] 4.3 在 `tencentcloud/services/sqlserver/resource_tc_sqlserver_basic_instance.md` 中补充 Import 说明，提示导入后会填充 `time_zone` 和 `disk_encrypt_flag`

## 5. 验证

- [x] 5.1 核实 schema 字段、create/read/update 逻辑、service 方法均已到位
- [x] 5.2 核实 `time_zone` 通过 `DescribeDBInstances` 读取，`disk_encrypt_flag` 通过 `DescribeDBInstancesAttribute` 读取
- [x] 5.3 核实 Read 操作中 `DescribeDBInstancesAttribute` 失败时仅记录 warning 不中断整个 Read
