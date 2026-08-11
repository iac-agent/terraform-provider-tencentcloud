# Proposal: Remove ForceNew from tag_value

## Why

The `tag_value` field in `tencentcloud_tag_attachment` currently has `ForceNew: true`, which means any change to `tag_value` forces resource recreation. This is problematic because:

1. If a user sets `tag_value` to a value and then changes it, Terraform will delete the old tag and create a new one, potentially causing the second request to fail.
2. Users should be able to update the tag value without recreating the resource.

## What Changes

We will modify the `tencentcloud_tag_attachment` resource so that `tag_value` is no longer ForceNew. This allows users to update the tag value without recreating the resource.

## Implementation

1. Change the schema for `tag_value` from:
   ```go
   Required: true,
   ForceNew: true,
   ```
   to:
   ```go
   Optional: true,
   Computed: true,
   ```

2. In the Read function, add logic to set `tag_value` from the API response if it exists.

3. Add unit tests to verify the behavior.

## References

- Similar resource implementation: `resource_tc_tag.go`
- Data source implementation: `data_source_tc_tag_keys.go`