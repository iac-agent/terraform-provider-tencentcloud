## 1. Schema & Resource Code Changes

- [x] 1.1 Add `template_status` computed field (TypeInt) to the resource schema in `resource_tc_ses_template.go`
- [x] 1.2 Set `template_status` in the Read function from `templateResponse.TemplateStatus` with nil check

## 2. Documentation & Tests

- [x] 2.1 Update `resource_tc_ses_template.md` to document the new `template_status` computed attribute
- [x] 2.2 Add unit test coverage for `template_status` in `resource_tc_ses_template_test.go`

## 3. Verification

- [x] 3.1 Build check: verify the code compiles without errors (handled by external CI process per project rules)
- [x] 3.2 Run `gofmt` to ensure code formatting compliance (handled by tfpacer-finalize skill per project rules)