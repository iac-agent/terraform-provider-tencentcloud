## 1. Service Layer

- [x] 1.1 Add `CreateSsmService(ctx context.Context, payMode, clientToken *string) (resp *ssm.CreateServiceResponseParams, errRet error)` to `tencentcloud/services/ssm/service_tencentcloud_ssm.go`, following the in-file SDKv2 conventions verbatim: `logId = tccommon.GetLogId(ctx)`, `request = ssm.NewCreateServiceRequest()`, the `[CRITAL]` failure `defer` log, `ratelimit.Check(request.GetAction())`, `me.client.UseSsmClient().CreateService(request)` and the `[DEBUG]` success log
- [x] 1.2 Set `request.PayMode` / `request.ClientToken` from the arguments without mutating them, return `response.Response` nil-safely (return `nil` when `response == nil`), and do not add any retry wrapper inside this method
- [x] 1.3 Reuse the existing `DescribeSsmServiceStatusByFilter` unchanged as the `GetServiceStatus` read (verify it already applies the three log conventions; do not duplicate it)

## 2. Action Implementation

- [x] 2.1 Create `tencentcloud/services/ssm/action_tc_ssm_service.go` with package `ssm`, `var _ action.ActionWithConfigure = &SsmServiceAction{}`, factory `func NewSsmServiceAction() action.Action`, struct `SsmServiceAction struct { fw.ActionWithConfigure }` and model `SsmServiceActionModel { PayMode types.String; ClientToken types.String; Timeouts actiontimeouts.Value }` (English comments only)
- [x] 2.2 Add package-level `defaultSsmServiceInvokeTimeout = 10 * time.Minute` and `var ssmServicePollInterval = 5 * time.Second` (the interval must be a `var` so tests can shorten it)
- [x] 2.3 Implement `Metadata` returning `resp.TypeName = "tencentcloud_ssm_service"`
- [x] 2.4 Implement `Schema`: `pay_mode` (Required String), `client_token` (Optional + Computed String), `timeouts` block via `actiontimeouts.BlockWithOpts` exposing only `invoke`, and a `Description` stating one-time semantics, repeated-invoke behaviour, billing order placement and that destroy does not deactivate SSM
- [x] 2.5 Implement `Invoke` input validation in order: `req.Config.Get` + `HasError` guard, `pay_mode` missing/empty check, `pay_mode != "PostPay"` check, `client_token` charset/length check (`^[\x21-\x7E]{1,64}$`), `a.Client() == nil` check — every failure uses `resp.Diagnostics.AddError` and returns before any API call
- [x] 2.6 Resolve `client_token` once before the retry block: use the user value when set, otherwise `fmt.Sprintf("ssm-%s", helper.BuildToken())`; log the resolved token
- [x] 2.7 Read `timeouts` via `data.Timeouts.Invoke(ctx, defaultSsmServiceInvokeTimeout)` and return on diagnostics error
- [x] 2.8 Build the service with `NewSsmService(a.Client())`, call `CreateSsmService` inside `resource.Retry(tccommon.WriteRetryTimeout, ...)` using `tccommon.RetryError(err)`, and log the returned `ResourceId` / `OrderId` nil-safely without writing state
- [x] 2.9 Add `func isSsmServiceAlreadyEnabled(err error) bool` (package-level, unit-testable) detecting `ssm.RESOURCEINUSE_SERVICEALREADYENABLED` via `*tcerr.TencentCloudSDKError` code and a string fallback; treat a hit as success by logging it, calling `resp.SendProgress` and continuing to convergence
- [x] 2.10 Implement the convergence loop: a `progress := func(msg string) { if resp.SendProgress != nil { resp.SendProgress(action.InvokeProgressEvent{Message: msg}) } }` helper, repeated `DescribeSsmServiceStatusByFilter` calls with `ssmServicePollInterval` sleeps, success on `ServiceEnabled == true`, fail-fast on `InvalidType` 2 or 3, nil-response tolerance, and no nested retry
- [x] 2.11 On timeout emit `resp.Diagnostics.AddError` whose summary says SSM did not converge and whose detail carries the last observed `service_enabled`, `invalid_type` and `resource_id`
- [x] 2.12 Add `resp.Diagnostics.AddWarning` when the `ResourceId` from `CreateService` and from `GetServiceStatus` are both non-empty and differ; never fail on that mismatch
- [x] 2.13 Self-check the file: only `Metadata` / `Schema` / `Invoke` implemented, no id, no state writes, no output attributes, no `Importer`, every error path returns through `Diagnostics`

## 3. Registration

- [x] 3.1 Append `ssm.NewSsmServiceAction,` to the `actionFactories` slice in `tencentcloud/framework/registry.go` (the `ssm` import already exists; do not edit `tencentcloud/provider.go`)
- [x] 3.2 Add an `Action` node listing `tencentcloud_ssm_service` under the `Secrets Manager(SSM)` node in `tencentcloud/provider.md`, placed after the `Ephemeral Resource` block and before the next product node

## 4. Documentation

- [x] 4.1 Create `tencentcloud/services/ssm/action_tc_ssm_service.md`: one-line English description mentioning SSM (Secrets Manager), a Terraform 1.14+ NOTE, an `Example Usage` section with an `action "tencentcloud_ssm_service" "example" { config { pay_mode = "PostPay" } }` block plus the `terraform apply -invoke=action.tencentcloud_ssm_service.example` command, and a note that repeated invokes re-trigger the activation order and that destroy does not deactivate SSM; no `Import` section and no Argument/Attribute Reference sections
- [x] 4.2 Verify no hand-written content is added under `website/` (the `website/docs/a/ssm_service.html.markdown` file must come from `make doc` in the finalize stage)

## 5. Unit Tests

- [x] 5.1 Create `tencentcloud/services/ssm/action_tc_ssm_service_test.go` in package `ssm_test` following `action_tc_bdrc_run_copy_pair_tasks_test.go`: obtain the schema via `a.Schema(...)`, build `action.InvokeRequest` with `tftypes.NewValue(tftypes.Object{...})`, and inject the client through `sharedmeta.ProviderMeta` + `a.Configure(...)`
- [x] 5.2 Add helpers to build invoke requests (valid, missing/empty `pay_mode`, non-`PostPay` `pay_mode`, malformed `client_token`, and a `timeouts` object carrying `invoke`) so that every case exercises the real schema
- [x] 5.3 Mock with gomonkey: `patches.ApplyMethodReturn(client, "UseSsmClient", ssmClient)` plus `ApplyMethodFunc` on `CreateService` and `GetServiceStatus`; shorten the poll interval via the package-level `var` (`ssmServicePollInterval`) or a short `timeouts.invoke` value so the timeout case finishes fast
- [x] 5.4 Cover the scenarios: successful invoke, already-enabled invoke, create API error, validation failures, provider-client-nil, and convergence timeout containing the last snapshot; test names prefixed `TestSsmServiceAction_`
- [x] 5.5 Ensure the test file compiles against the current dependency set (gomonkey, testify, tftypes already vendored) without executing it

## 6. Verification (code-level only; do not run builds or tests in this stage)

- [x] 6.1 Cross-check every SDK field used against `design.md`'s vendor excerpt: `CreateServiceRequest.PayMode` / `ClientToken`, `CreateServiceResponseParams.ResourceId` / `OrderId`, `GetServiceStatusResponseParams.ServiceEnabled` / `InvalidType` / `ResourceId`
- [x] 6.2 Confirm the implementation follows the design's parameter-mapping constraint: only `PayMode` and `ClientToken` are sent on create (they are the only fields the API accepts), and `GetServiceStatus` is sent with an empty request
- [x] 6.3 Confirm no vendor or `go.mod` change is introduced by this change
- [x] 6.4 Confirm the generated artifacts contain no non-English text in Go code, comments, `.md` files or tests

## 7. Finalize (executed by the tfpacer-finalize stage only)

- [x] 7.1 Run `gofmt` on the new/changed Go files
- [x] 7.2 Run `make doc` so that `website/docs/a/ssm_service.html.markdown` is generated
- [x] 7.3 Add the changelog file under `.changelog/` with `release-note:new-action`
