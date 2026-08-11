# Design: Remove ForceNew from tag_value

## Problem

Currently, `tag_value` is `Required: true` and `ForceNew: true` in the schema. This means any change to `tag_value` forces resource recreation, which is not ideal.

## Solution

Change `tag_value` from `Required: true, ForceNew: true` to `Optional: true, Computed: true`. This makes `tag_value` updatable without resource recreation.

## Changes

1. Schema: Change `tag_value` from Required+ForceNew to Optional+Computed
2. Read function: Set `tag_value` from API if available
3. Tests: Add unit tests to verify Read function behavior