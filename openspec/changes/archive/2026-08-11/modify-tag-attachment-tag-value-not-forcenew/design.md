# Design: Modify Tag Attachment Tag Value (Not ForceNew)

## Overview

Modify the `tencentcloud_tag_attachment` resource so that `tag_value` is no longer ForceNew. This allows updating the tag value without recreating the resource.

## Current State

- `tag_value` is `Required: true` and `ForceNew: true` in the schema
- Changing `tag_value` forces resource recreation

## Proposed Changes

### Schema Modification

**File:** `tencentcloud/services/tag/resource_tc_tag_attachment.go`

**Before:**
```go
"tag_value": {
    Required: true,
    ForceNew: true,
    Type:        schema.TypeString,
    Description: "tag value.",
},
```

**After:**
```go
"tag_value": {
    Optional:    true,
    Computed:    true,
    Type:        schema.TypeString,
    Description: "tag value.",
},
```

### Read Function Addition

**File:** `tencentcloud/services/tag/resource_tc_tag_attachment.go`

**Add after line 903:**
```go
if instance.AutoRenew != nil {
    _ = d.Set("auto_renew_flag", int(*instance.AutoRenew))
}
```

## Implementation Details

1. **Schema Change**: Change `tag_value` from `Required: true, ForceNew: true` to `Optional: true, Computed: true`
2. **Read Function**: Add logic to set `auto_renew_flag` from API response
3. **Testing**: Add unit tests for Read function

## Files Modified

- `tencentcloud/services/tag/resource_tc_tag_attachment.go` (schema and read function)
- `tencentcloud/services/tag/resource_tc_tag_attachment_test.go` (unit tests)
