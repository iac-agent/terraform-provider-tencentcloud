# Design: Update teo_l7_acc_rule zone_id Description

## Architecture

This is a documentation-only change that updates the Schema description of the `zone_id` field in the `tencentcloud_teo_l7_acc_rule` resource. No functional logic or data structures will be modified.

### File Structure

```
tencentcloud/services/teo/resource_tc_teo_l7_acc_rule.go  # Schema definition (修改)
website/docs/r/teo_l7_acc_rule.html.markdown              # 文档 (修改)
```

## Detailed Design

### 1. Current Description

```go
"zone_id": {
    Type:        schema.TypeString,
    Required:    true,
    ForceNew:    true,
    Description: "Zone id, required field.",
},
```

### 2. Updated Description

```go
"zone_id": {
    Type:        schema.TypeString,
    Required:    true,
    ForceNew:    true,
    Description: "Zone id, which must be a valid value and cannot be null or empty string.",
},
```

### 3. Documentation Update

The website/docs/r/teo_l7_acc_rule.html.markdown file's `zone_id` argument reference also needs to be updated to match the new description.

**Current**:
```
* `zone_id` - (Required, String, ForceNew) Zone id, required field.
```

**Updated**:
```
* `zone_id` - (Required, String, ForceNew) Zone id, which must be a valid value and cannot be null or empty string.
```

## Validation

### Success Criteria
- [x] `zone_id` field Description in Go source code has been updated
- [x] Documentation in website/docs/r/teo_l7_acc_rule.html.markdown has been updated
- [x] No functional logic changes
- [x] No syntax errors

### Code Quality
- [x] Follows existing code style
- [x] Description is clear and accurate
