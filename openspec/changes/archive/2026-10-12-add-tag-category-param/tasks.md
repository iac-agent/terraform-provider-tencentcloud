# Tasks: 为 tencentcloud_tag 资源新增 Category 参数

## 1. Schema 与 Read 逻辑修改

- [x] 1.1 在 `tencentcloud/services/tag/resource_tc_tag.go` 的 `ResourceTencentCloudTag()` schema map 中（`tag_value` 定义之后）新增 `Category` 参数：`Computed: true, Type: schema.TypeString`，Description 说明取值为 Custom（自定义标签）/ System（系统标签）/ All（全部标签）
- [x] 1.2 在 `resourceTencentCloudTagResourceRead()` 中（`tag_value` 回填逻辑之后）新增 `Category` 回填：`tagRes.Category != nil` 时 `_ = d.Set("Category", tagRes.Category)`，为 nil 时跳过；不修改 Create/Delete 逻辑（`CreateTag` 无 `Category` 入参，已对照 vendor SDK 核实）

## 2. 测试用例补充

- [x] 2.1 在 `tencentcloud/services/tag/resource_tc_tag_test.go` 中沿用现有 terraform 测试套件（`resource.Test` / `TestCheckResourceAttr`）补充 `Category` 断言：创建后检查 `Category` 属性为 `Custom`（用户创建的标签均为自定义标签）；确保生成的测试代码可正确构建（禁止执行 `go test`）

## 3. 文档更新

- [x] 3.1 更新 `tencentcloud/services/tag/resource_tc_tag.md`：补充 `Category` 参数说明（Computed，标签类型：Custom/System/All，由云平台返回）及示例中的属性输出；`website/docs/` 目录由收尾阶段 `make doc` 生成，禁止手动修改

## 4. 验证

- [x] 4.1 使用 `git diff` 核对全部变更：schema 新增 `Category`（Computed）、Read 判空回填、测试与文档同步更新，确认未引入 `go build`/`go vet` 等验证命令执行，确认未改动 `service_tencentcloud_tag.go`、`provider.go`、`provider.md` 及 vendor 目录
