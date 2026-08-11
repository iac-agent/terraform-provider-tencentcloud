# Tasks: Remove ForceNew from tag_value

## Phase 1: Schema Modification
- [ ] Change `tag_value` from `Required: true, ForceNew: true` to `Optional: true, Computed: true`
- [ ] Update description to reflect new behavior

## Phase 2: Read Function Update
- [ ] Add logic to set `tag_value` from API response
- [ ] Add unit tests for Read function

## Phase 3: Testing
- [ ] Add acceptance tests for the resource