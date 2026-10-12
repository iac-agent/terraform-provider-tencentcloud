# Design: 为 tencentcloud_tag 资源新增 Category 参数

## Context

`tencentcloud_tag` 是 RESOURCE_KIND_GENERAL 通用资源（文件：`tencentcloud/services/tag/resource_tc_tag.go`），管理标签键值的完整生命周期（Create via `CreateTag` / Read via `GetTags` / Delete via `DeleteTag`）。资源 ID 为复合 ID：`tagKey + tccommon.FILED_SP + tagValue`。

当前 schema 仅有两个参数，均来自 `CreateTag` 入参与 `GetTags` 出参：

```go
// tencentcloud/services/tag/resource_tc_tag.go 现有 schema（第 27-41 行）
Schema: map[string]*schema.Schema{
    "tag_key": {
        Required:    true,
        ForceNew:    true,
        Type:        schema.TypeString,
        Description: "Tag key.",
    },
    "tag_value": {
        Required:    true,
        ForceNew:    true,
        Type:        schema.TypeString,
        Description: "Tag value.",
    },
},
```

Read 流程现状（`resourceTencentCloudTagResourceRead`，第 92-130 行）：
1. 从 `d.Id()` 按 `tccommon.FILED_SP` 拆分出 `tagKey`/`tagValue`
2. 调用 `service.DescribeTagResourceById(ctx, tagKey, tagValue)`（服务层方法，内部调用 `GetTags` API，已在 `resource.Retry(tccommon.ReadRetryTimeout, ...)` 块中处理重试）
3. `tagRes == nil` 时打日志并 `d.SetId("")`
4. `tagRes.TagKey` / `tagRes.TagValue` 非 nil 时分别 `d.Set("tag_key", ...)` / `d.Set("tag_value", ...)`

服务层 `DescribeTagResourceById`（`tencentcloud/services/tag/service_tencentcloud_tag.go` 第 300-345 行）已返回完整的 `*tag.Tag` 结构体（含 `Category` 字段），**无需修改服务层**。

## Goals / Non-Goals

**Goals:**

- 在 `tencentcloud_tag` 资源 schema 中新增 `Category` 参数（Computed, String），暴露标签类型（Custom/System/All）
- 资源 Read 时从 `GetTags` 响应的 `Tag.Category` 回填 `Category` 到 Terraform state
- 同步更新资源文档 `resource_tc_tag.md` 与测试用例 `resource_tc_tag_test.go`

**Non-Goals:**

- 不修改 `CreateTag` / `DeleteTag` 调用逻辑（`CreateTagRequest` 无 `Category` 入参，标签类型由云平台侧决定）
- 不修改服务层 `service_tencentcloud_tag.go`（现有返回结构已包含所需字段）
- 不新增/修改 provider.go、provider.md（资源已注册）
- 不处理 `tag_key` / `tag_value`（已存在于现有 schema 中，需求映射中的这两个参数为既有参数，非本次变更范围）

## Decisions

### D1: 参数命名为 `Category`（保留需求指定大小写），Computed 属性

**决策**: schema key 直接使用需求指定的 SchemaName `Category`（非小写蛇形），属性为 Computed。

**理由**:
- 需求中 JsonPath `response.Tags.Category` → SchemaName `Category` 的映射已明确指定参数名为 `Category`
- `Category` 为 `GetTags` 出参，`CreateTag` 入参中不存在该字段，用户不可设置 → Computed
- Computed 字段对现有配置/state 完全向后兼容，refresh 时自动填充

**备选方案**: 命名为小写 `category`。放弃原因：需求 SchemaName 明确指定为 `Category`，以需求为准。

### D2: Read 中判空后回填，不引入 Update 流程

**决策**: 在 `resourceTencentCloudTagResourceRead` 中，`tagRes.Category != nil` 时 `_ = d.Set("Category", tagRes.Category)`；为 nil 时跳过（SDK 注释标注该字段可能返回 null）。

**理由**:
- 遵循现有 Read 代码对 `TagKey`/`TagValue` 的判空处理模式
- SDK 中 `Tag.Category` 带"注意：此字段可能返回 null"注释，直接 Set nil 指针会触发 panic 或写入空值，判空是必须的
- 资源无 Update 方法（仅 Create/Read/Delete），Computed 参数天然不需要更新逻辑

### D3: 沿用现有 terraform 测试套件补充测试

**决策**: 在 `resource_tc_tag_test.go` 中沿用现有 `resource.Test` + `TestCheckResourceAttr` 风格补充 `Category` 相关断言与用例（修改现有资源需求，按项目规则使用 terraform 测试套件）。

## 云 API 参考（摘自 vendor，实施时无需回读）

以下为 vendored SDK（`vendor/github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tag/v20180813/models.go`）中与本次变更相关的 struct 定义摘录。

### CreateTagRequest（models.go 第 256-264 行）— Create 接口入参

```go
type CreateTagRequest struct {
	*tchttp.BaseRequest

	// 标签键
	TagKey *string `json:"TagKey,omitnil" name:"TagKey"`

	// 标签值
	TagValue *string `json:"TagValue,omitnil" name:"TagValue"`
}
```

说明：`CreateTag` 入参仅有 `TagKey` / `TagValue`，**无 `Category` 入参**（已在 vendor 中核对），故 `Category` 只能作为 Computed 参数，Create 流程不传该字段。

### GetTagsRequest（models.go 第 2161-2179 行）— Read 接口入参

```go
type GetTagsRequest struct {
	*tchttp.BaseRequest

	// 从上一页的响应中获取的下一页的Token值。
	// 如果是第一次请求，设置为空。
	PaginationToken *string `json:"PaginationToken,omitnil" name:"PaginationToken"`

	// 每一页返回的数据最大条数，最大1000。
	// 缺省值：50。
	MaxResults *uint64 `json:"MaxResults,omitnil" name:"MaxResults"`

	// 标签键。
	// 返回所有标签键列表对应的标签。
	// 最大长度：20
	TagKeys []*string `json:"TagKeys,omitnil" name:"TagKeys"`

	// 标签类型。取值： Custom：自定义标签。 System：系统标签。 All：全部标签。 默认值：All。
	Category *string `json:"Category,omitnil" name:"Category"`
}
```

说明：`GetTags` 请求支持按 `Category` 过滤，但服务层 `DescribeTagResourceById` 现有实现未传该过滤参数（默认 All，返回全部），无需修改即可获取到目标标签的 `Category` 值。

### GetTagsResponseParams / GetTagsResponse（models.go 第 2204-2218 行）— Read 接口出参

```go
type GetTagsResponseParams struct {
	// 获取的下一页的Token值
	PaginationToken *string `json:"PaginationToken,omitnil" name:"PaginationToken"`

	// 标签列表。
	Tags []*Tag `json:"Tags,omitnil" name:"Tags"`

	// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
	RequestId *string `json:"RequestId,omitnil" name:"RequestId"`
}

type GetTagsResponse struct {
	*tchttp.BaseResponse
	Response *GetTagsResponseParams `json:"Response"`
}
```

### Tag（models.go 第 2446-2456 行）— Tags 列表元素结构体（**本次新增参数的数据来源**）

```go
type Tag struct {
	// 标签键
	TagKey *string `json:"TagKey,omitnil" name:"TagKey"`

	// 标签值
	TagValue *string `json:"TagValue,omitnil" name:"TagValue"`

	// 标签类型。取值： Custom：自定义标签。 System：系统标签。 All：全部标签。 默认值：All。
	// 注意：此字段可能返回 null，表示取不到有效值。
	Category *string `json:"Category,omitnil" name:"Category"`
}
```

JsonPath 映射：`response.Response.Tags[].Category` → Terraform schema `Category`（服务层 `DescribeTagResourceById` 返回单个 `*tag.Tag`，即列表中匹配 tagKey+tagValue 的元素）。

### 服务层现有方法（tencentcloud/services/tag/service_tencentcloud_tag.go 第 300-345 行，无需修改）

```go
func (me *TagService) DescribeTagResourceById(ctx context.Context, tagKey string, tagValue string) (tagRes *tag.Tag, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := tag.NewGetTagsRequest()
	response := tag.NewGetTagsResponse()
	request.TagKeys = []*string{&tagKey}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseTagClient().GetTags(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil || result.Response.Tags == nil || len(result.Response.Tags) == 0 {
			return resource.NonRetryableError(fmt.Errorf("Describe tags failed, Response is nil."))
		}

		response = result
		return nil
	})

	if err != nil {
		errRet = err
		return
	}

	for _, v := range response.Response.Tags {
		if v != nil && v.TagKey != nil && v.TagValue != nil {
			if *v.TagKey == tagKey && *v.TagValue == tagValue {
				tagRes = v
				break
			}
		}
	}

	return
}
```

注：`ratelimit` 为该文件头部引入的包（`tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"` 同文件已有相关 import，实施时无需重复处理）。返回的 `tagRes` 即 `*tag.Tag`，天然携带 `Category` 字段。

## 实施要点（代码级）

### schema 新增（resource_tc_tag.go）

在 `Schema` map 中 `tag_value` 之后追加：

```go
"Category": {
	Computed:    true,
	Type:        schema.TypeString,
	Description: "Tag category, valid values: `Custom` (custom tag), `System` (system tag), `All` (all tags).",
},
```

### Read 回填（resourceTencentCloudTagResourceRead）

在现有 `tag_value` 回填逻辑（第 125-127 行）之后追加：

```go
if tagRes.Category != nil {
	_ = d.Set("Category", tagRes.Category)
}
```

### 文档更新（resource_tc_tag.md）

在示例中补充 `Category` 的输出说明（Computed 属性，值由云平台返回：Custom/System）。

注意：`website/docs` 目录由收尾阶段 `make doc` 生成，本变更在实施阶段仅更新 `tencentcloud/services/tag/resource_tc_tag.md`。

### 测试补充（resource_tc_tag_test.go）

沿用现有测试套件（`resource.Test` / `TestCheckResourceAttr`），在 check 函数中补充 `Category` 断言（期望值为 `Custom`，用户通过 `tencentcloud_tag` 创建的标签均为自定义标签）。禁止执行 `go test`，仅保证代码可正确构建。

## Risks / Trade-offs

- [SDK 返回 `Category == nil`] → Read 中判空后再 Set，为 nil 时跳过，不写入 state、不报错（与现有 `TagKey`/`TagValue` 处理模式一致）
- [参数名为大写驼峰 `Category`，与仓库内其他小写蛇形参数风格不一致] → 以需求指定的 SchemaName 为准；Computed 参数不影响 HCL 语法，用户以 `Category` 引用
- [现有用户 state 中无 `Category` 字段] → Computed 字段在下次 plan/apply 时自动填充，向后兼容，无迁移成本
- [服务层返回的匹配元素可能不携带 Category（旧 SDK 版本/云 API 未返回）] → 已核对 vendored SDK `Tag` 结构体包含 `Category` 字段；即使值为 nil 也仅跳过回填，不产生错误

## Migration Plan

1. 实施 schema + Read 回填 + 文档 + 测试变更（单 PR/单变更）
2. 回滚策略：直接 revert 变更即可，Computed 参数不落盘任何用户输入，无 state 迁移需求

## Open Questions

无。
