# Tasks: Modify Tag Attachment Tag Value (Not ForceNew)

## Phase 1: Schema Modification
- [ ] Change `tag_value` from `Required: true, ForceNew: true` to `Optional: true, Computed: true`
- [ ] Update description to reflect new behavior

## Phase 2: Read Function Update
- [ ] Add logic to set `auto_renew_flag` from API response
- [ ] Add nil check for `instance.AutoRenew`
- [ ] Add type conversion from `*uint64` to `int`

## Phase 3: Testing
- [ ] Add unit tests for Read function
- [ ] Add acceptance tests for the resource
- [ ] Run `terraform plan` to verify no diff
