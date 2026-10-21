## Why

The TencentCloud provider has no support for the FWM (Firewall Management, 云防火墙下发规则) strategy resource. Users managing FWM下发规则 (strategy) with Terraform cannot create, read, update, or delete them today. The FWM SDK package `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/fwm/v20250611` is already vendored, so the provider should expose this resource.

## What Changes

- Add a new Terraform Plugin Framework resource `tencentcloud_fwm_strategy` (RESOURCE_KIND_GENERAL, full CRUD + import).
- New file `tencentcloud/services/fwm/resource_tc_fwm_strategy.go` with factory `NewFwmStrategyResource()`.
- Cloud APIs used: `CreateStrategy` (create), `DescribeStrategy` (read), `ModifyStrategy` (update), `DeleteStrategy` (delete) from the `fwm/v20250611` SDK package.
- Register the factory in `tencentcloud/framework/registry.go` (`resourceFactories`) and add the resource to the `fwm` node in `tencentcloud/provider.md`.
- Add `UseFwmClient()` accessor in `tencentcloud/connectivity/client.go` (no FWM client exists yet).
- Add sibling docs file `resource_tc_fwm_strategy.md` and unit test `resource_tc_fwm_strategy_test.go` (gomonkey mock, no TF test suite, not executed).

## Capabilities

### New Capabilities
- `fwm-strategy-resource`: Full lifecycle management (Create/Read/Update/Delete/Import) of a TencentCloud FWM strategy (下发规则) via the Terraform Plugin Framework resource `tencentcloud_fwm_strategy`.

### Modified Capabilities
<!-- none -->

## Impact

- Code: `tencentcloud/services/fwm/` (new package), `tencentcloud/framework/registry.go`, `tencentcloud/provider.md`, `tencentcloud/connectivity/client.go`.
- APIs: `fwm/v20250611` `CreateStrategy` / `DescribeStrategy` / `ModifyStrategy` / `DeleteStrategy`.
- No breaking change; purely additive.
