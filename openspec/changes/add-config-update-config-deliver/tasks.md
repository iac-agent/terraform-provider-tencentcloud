## 1. Service Layer

- [x] 1.1 In `tencentcloud/services/config/service_tencentcloud_config.go`, append
      `func (me *ConfigService) DescribeConfigDeliver(ctx context.Context) (*configv20220802.DescribeConfigDeliverResponseParams, error)`:
      build `configv20220802.NewDescribeConfigDeliverRequest()`, call
      `me.client.UseConfigV20220802Client().DescribeConfigDeliverWithContext(ctx, request)` inside a retry loop
      bounded by `tccommon.ReadRetryTimeout`, return `response.Response`; on failure wrap the error with the API
      action name and log with `tccommon.GetLogId(ctx)`. Return an error when `response.Response` is nil.
- [x] 1.2 In the same file, append
      `func (me *ConfigService) UpdateConfigDeliver(ctx context.Context, request *configv20220802.UpdateConfigDeliverRequest) error`:
      call `UpdateConfigDeliverWithContext` inside a retry loop bounded by `tccommon.WriteRetryTimeout`; treat
      a nil `response.Response` as an error; log the log ID and both request/response bodies on both paths.
- [x] 1.3 Confirm neither new method restructures existing `ConfigService` members and that both reuse
      `UseConfigV20220802Client()` (no new SDK client, no `go.mod` / `vendor/` change).

## 2. Resource Schema and Scaffolding

- [x] 2.1 Create `tencentcloud/services/config/resource_tc_config_update_config_deliver_config.go` with package
      `config`, the sentinel constant `configUpdateConfigDeliverId = "config_update_config_deliver"`, and the
      three compile-time assertions for `resource.Resource`, `resource.ResourceWithConfigure`,
      `resource.ResourceWithImportState`.
- [x] 2.2 Define `ConfigUpdateConfigDeliverResource{ client ConfigService }` and
      `NewConfigUpdateConfigDeliverResource() resource.Resource`. Do not import or embed
      `tencentcloud/framework/fw`.
- [x] 2.3 Define `ConfigUpdateConfigDeliverResourceModel` with `id`, `status`, `deliver_name`, `target_arn`,
      `deliver_prefix`, `deliver_type`, `deliver_content_type`, `create_time` (`types.String` / `types.Int64`).
- [x] 2.4 Implement `Metadata` setting `resp.TypeName = "tencentcloud_config_update_config_deliver"`.
- [x] 2.5 Implement `Schema` exactly as the design table: `id` Computed + `stringplanmodifier.UseStateForUnknown()`;
      `status` Required int64; `deliver_name` / `target_arn` / `deliver_prefix` Optional+Computed string;
      `deliver_type` Optional+Computed string with `stringvalidator.OneOf("COS", "CLS")`;
      `deliver_content_type` Optional+Computed int64 with `int64validator.Between(1, 3)`; `create_time` Computed string.
      No `RequiresReplace`, no `Default`. Description of the schema states that destroying the resource disables
      Config log delivery and never deletes the Config service, COS bucket or CLS topic.

## 3. Resource Lifecycle

- [x] 3.1 Implement `Configure` by hand: return when `req.ProviderData == nil`; assert `*sharedmeta.ProviderMeta`;
      on failure add an "Unexpected ProviderData type" diagnostic including `%T`; on success set
      `r.client = NewConfigService(meta.Client)`.
- [x] 3.2 Implement the shared private `modifyConfig(ctx, *Model) error`: build
      `configv20220802.NewUpdateConfigDeliverRequest()`, always set `Status` via
      `tccommonhelper.Uint64(uint64(plan.Status.ValueInt64()))` (import
      `"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"` aliased as
      `tccommonhelper` — this is the SDKv2-side helper package, whose `Uint64(uint64) *uint64` matches the SDK
      field type; the framework-side `framework/internal/helper` has no `*uint64` producer), and set each
      remaining field only when `!IsNull() && !IsUnknown()`; call `r.client.UpdateConfigDeliver`; wrap errors with
      the resource ID.
- [x] 3.3 Implement `Create`: `req.Plan.Get` → `modifyConfig` → set `plan.Id = types.StringValue(configUpdateConfigDeliverId)`
      → `resp.State.Set` (do not call `d.SetId`; this is Framework).
- [x] 3.4 Implement `Read`: `req.State.Get`; restore the sentinel ID into state if the ID is empty (post-import);
      call `r.client.DescribeConfigDeliver`; error out when `response == nil` (keep the resource in state, do not
      `RemoveResource`); map each field with `types.StringPointerValue` / `helper.Int64ValueFromUint` so nil pointers
      become null; `log.Printf` including the resource ID on failure; `resp.State.Set`.
- [x] 3.5 Implement `Update`: `req.Plan.Get` → `modifyConfig` → preserve the ID from state → `resp.State.Set`.
- [x] 3.6 Implement `Delete` (turn-off semantic): `req.State.Get` to obtain the ID for logging, build
      `UpdateConfigDeliverRequest{Status: <*uint64 pointing to 0>}` as the sole field (use the project's
      `tencentcloud/internal/helper.Uint64(0)`, which returns `*uint64`; `framework/internal/helper` only
      offers `Int64PointerFromValue(v types.Int64) *int64`, so it cannot provide a `*uint64`), call
      `r.client.UpdateConfigDeliver`; on failure add a diagnostic; on success return without calling
      `RemoveResource` (the Framework drops the state entry after Delete returns nil). Never delete the host
      resource.
- [x] 3.7 Implement `ImportState`: reject any ID other than the sentinel with a diagnostic naming the valid ID;
      otherwise `resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)`.

## 4. Framework Registration

- [x] 4.1 Add the `config` service package import and `config.NewConfigUpdateConfigDeliverResource` to
      `resourceFactories` in `tencentcloud/framework/registry.go`, keeping the one-entry-per-line style.
- [x] 4.2 Verify `tencentcloud/provider.go` is **not** modified for this resource.
- [x] 4.3 Add `tencentcloud_config_update_config_deliver` under the Config product's Resources list in
      `tencentcloud/provider.md`, adjacent to `tencentcloud_config_deliver_config` and
      `tencentcloud_config_recorder_config`.

## 5. Documentation

- [x] 5.1 Create `tencentcloud/services/config/resource_tc_config_update_config_deliver_config.md` in English:
      a one-paragraph `Provides a resource to ...` description naming the Tencent Cloud Config product and
      stating the singleton nature and the turn-off-on-destroy semantic; an `Example Usage` HCL block showing
      COS and CLS variants; an `Import` block documenting the sentinel ID
      (`terraform import tencentcloud_config_update_config_deliver.example config_update_config_deliver`).
      No `Argument Reference` / `Attribute Reference` sections (generated by tooling). Do not hand-write
      anything under `website/`.

## 6. Unit Tests

- [x] 6.1 Create `tencentcloud/services/config/resource_tc_config_update_config_deliver_config_test.go` using
      gomonkey to mock `ConfigService.DescribeConfigDeliver` / `ConfigService.UpdateConfigDeliver` (no Terraform
      acceptance-test harness, no real cloud calls). Do not execute `go test`.
- [x] 6.2 Cover: `Metadata` TypeName; `Schema` attribute modes (`status` Required, `create_time` Computed,
      Optional+Computed on the rest, `id` has `UseStateForUnknown`).
- [x] 6.3 Cover `modifyConfig` request building: `status` always sent; optional fields omitted when null/unknown;
      `deliver_type` / `deliver_content_type` mapped correctly; API error propagated.
- [x] 6.4 Cover `Read` nil-handling: nil response pointer → diagnostic and state retained; nil field pointers →
      null state values; non-nil `uint64` pointers → correct `types.Int64`.
- [x] 6.5 Cover `Delete`: asserts the request carries `Status = 0` and no other field, and that a service error
      surfaces as a diagnostic.

## 7. Validation Handoff

- [x] 7.1 Self-review the diff with `git diff` against the spec scenarios; do not run `go build`, `go vet`,
      `golangci-lint` or `go test` in this change phase.
- [x] 7.2 Hand off to the finalize phase for `gofmt`, `make doc` and the `.changelog` entry.

## 8. Follow-up Hardening

- [x] 8.1 Guard `modifyConfig` against a null or unknown planned `status`: `ValueInt64()` would coerce it to
      `0` and silently disable log delivery, so return an error before any `UpdateConfigDeliver` call.
      Add the matching unit test (service must not be invoked in either the null or the unknown case).
- [x] 8.2 Log the resource ID on the successful `Read` path as well, so a refresh can be traced end to end.
- [x] 8.3 Confirm the framework-side `tencentcloud/framework/internal/helper` package cannot be imported from
      `tencentcloud/services/config` (Go's `internal` visibility rule restricts it to the
      `tencentcloud/framework/...` subtree), so the `*uint64` -> `types.Int64` conversion stays local.
- [x] 8.4 Re-verify the vendored `config/v20220802` structs still expose every field used by the resource:
      `UpdateConfigDeliverRequest.Status/DeliverName/TargetArn/DeliverPrefix/DeliverType/DeliverContentType`
      and `DescribeConfigDeliverResponseParams` plus `CreateTime`.

