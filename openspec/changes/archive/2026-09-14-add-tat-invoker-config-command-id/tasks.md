## 1. Schema & Resource Code Update

- [x] 1.1 Add `command_id` field (TypeString, Computed: true) to the resource Schema in `resource_tc_tat_invoker_config.go`
- [x] 1.2 Update `resourceTencentCloudTatInvokerConfigRead` to read `command_id` from `invokerConfig.CommandId` and set it via `d.Set("command_id", ...)`, with nil check
- [ ] 1.3 Run `gofmt` to ensure code formatting is correct (skipped, performed in finalize phase)

## 2. Unit Test Update

- [x] 2.1 Update `resource_tc_tat_invoker_config_test.go` to add test validation for the new `command_id` attribute

## 3. Documentation Update

- [x] 3.1 Update `resource_tc_tat_invoker_config.md` to reflect the new `command_id` computed attribute in the example/description
- [ ] 3.2 Run `make doc` to regenerate website documentation (skipped, performed in finalize phase)