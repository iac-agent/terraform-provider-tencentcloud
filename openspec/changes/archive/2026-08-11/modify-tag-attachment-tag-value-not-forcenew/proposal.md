# Modify Tag Attachment Tag Value (Not ForceNew)

## Why

The `tencentcloud_tag_attachment` resource currently requires `tag_value` to be
ForceNew. This means that if a user sets `tag_value` to a value and then changes
it, Terraform will force a recreation of the entire resource.

## What Changes

We will modify the `tencentcloud_tag_attachment` resource so that `tag_value` is
no longer ForceNew. This allows users to update the tag value without recreating
the resource.

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

2. In the Read function, call `d.Set("tag_value", ...)` if `tag_value` has a
   value.

3. Add unit tests to verify the behavior.

## References

- Similar resource implementation: `resource_tc_tag.go`
- Data source implementation: `data_source_tc_tag_keys.go`
