# Remove ForceNew from tag_value

## Background

The `tencentcloud_tag_attachment` resource currently requires `tag_value` to be ForceNew. This means that if a user sets `tag_value` and then changes it, Terraform will delete the old tag and create a new one, causing the second request to fail.

## Goal

Allow `tag_value` to be updated without recreating the resource.

## Changes

1. Change the schema for `tag_value` from `Required: true, ForceNew: true` to `Optional: true, Computed: true`
2. In the Read function, set `tag_value` from the API if available
3. Add tests to verify the behavior