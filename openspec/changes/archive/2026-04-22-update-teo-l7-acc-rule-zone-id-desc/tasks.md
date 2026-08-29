# Tasks: Update teo_l7_acc_rule zone_id Description

## Task 1: 更新 zone_id 字段描述

**文件**: `tencentcloud/services/teo/resource_tc_teo_l7_acc_rule.go`

### 实施步骤

#### 1.1 定位 zone_id 字段

找到 Schema 中的 `zone_id` 字段定义（约第 27-32 行）：
```go
"zone_id": {
    Type:        schema.TypeString,
    Required:    true,
    ForceNew:    true,
    Description: "Zone id, required field.",
},
```

#### 1.2 替换 Description 内容

将 Description 从 `"Zone id, required field."` 替换为 `"Zone id, which must be a valid value and cannot be null or empty string."`

### 验收标准

- [x] `zone_id` 字段的 Description 已更新
- [x] 无语法错误

---

## Task 2: 更新文档描述

**文件**: `website/docs/r/teo_l7_acc_rule.html.markdown`

### 实施步骤

#### 2.1 定位 zone_id 文档描述

找到 Argument Reference 中的 `zone_id` 行：
```
* `zone_id` - (Required, String, ForceNew) Zone id, required field.
```

#### 2.2 替换描述内容

将描述从 `Zone id, required field.` 替换为 `Zone id, which must be a valid value and cannot be null or empty string.`

### 验收标准

- [x] 文档中 `zone_id` 的描述已更新
- [x] 描述与 Go 源代码中的 Description 一致

---

## Task 3: 更新 .md 示例文件

**文件**: `tencentcloud/services/teo/resource_tc_teo_l7_acc_rule.md`

### 实施步骤

#### 3.1 检查 .md 示例文件是否存在，如果存在则更新对应描述

### 验收标准

- [x] .md 示例文件中的描述已同步更新（如果文件存在）

---

## 验收清单

### 描述准确性
- [x] Description 明确说明 zone_id 必须输入有效值
- [x] Description 说明不能为 null
- [x] Description 说明不能为空字符串 ""

### 代码质量
- [x] 无语法错误
- [x] 遵循现有代码风格

### 完整性
- [x] Go 源代码描述已更新
- [x] 文档描述已更新
- [x] .md 示例文件已同步（如存在）
