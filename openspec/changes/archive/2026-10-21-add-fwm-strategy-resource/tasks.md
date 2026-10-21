## 1. Client & Registration Setup

- [x] 1.1 Add `UseFwmClient()` accessor and `fwm/v20250611` import in `tencentcloud/connectivity/client.go`
- [x] 1.2 Create `tencentcloud/services/fwm/` package

## 2. Resource Implementation

- [x] 2.1 Implement `tencentcloud/services/fwm/resource_tc_fwm_strategy.go`: model, Metadata, Schema (attributes per design.md, with RequiresReplace/UseStateForUnknown plan modifiers)
- [x] 2.2 Implement Configure (sharedmeta.ProviderMeta), Create (CreateStrategy, empty-id NonRetryableError, SetId then read-back), Read (read helper + RemoveResource on not-found)
- [x] 2.3 Implement Update (ModifyStrategy from Plan) and Delete (DeleteStrategy with StrategyIds), plus refreshModel copying computed fields with nil-safe types.XxxPointerValue

## 3. Registration & Docs

- [x] 3.1 Register `fwm.NewFwmStrategyResource` in `tencentcloud/framework/registry.go` resourceFactories
- [x] 3.2 Add `tencentcloud_fwm_strategy` entry under the `fwm` node in `tencentcloud/provider.md`
- [x] 3.3 Write `tencentcloud/services/fwm/resource_tc_fwm_strategy.md` (English description "Provides a resource to ...", Example Usage with jsonencode where applicable, Import section)

## 4. Unit Tests

- [x] 4.1 Write `tencentcloud/services/fwm/resource_tc_fwm_strategy_test.go` with gomonkey-mocked SDK CRUD responses (no TF test suite, not executed)

## 5. Verification

- [x] 5.1 Self-check: request/response field mapping against design.md excerpts, error handling, ID semantics, English-only content; confirm all artifacts complete via `openspec status`
