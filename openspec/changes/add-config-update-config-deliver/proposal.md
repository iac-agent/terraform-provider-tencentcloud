## Why

Tencent Cloud Config (云审计/配置审计) delivers compliance and resource-snapshot logs to a COS bucket or a
CLS topic through a per-account, per-region "delivery setting" (投递设置). That setting can only be read
(`DescribeConfigDeliver`) and edited (`UpdateConfigDeliver`); there is no create or delete API for it — it
always exists as a global singleton. The provider today has no way to declare this setting as code, so
teams managing Config with Terraform must configure delivery manually in the console, which breaks
reproducibility and drift detection for a compliance-critical setting.

## What Changes

- Add a new resource `tencentcloud_config_update_config_deliver` (RESOURCE_KIND_CONFIG) that manages the
  Config delivery setting: read it with `DescribeConfigDeliver`, set it with `UpdateConfigDeliver`.
- Implement it with **Terraform Plugin Framework** (not SDKv2/helper/schema), registered through
  `tencentcloud/framework/registry.go`.
- Expose the updatable fields `status`, `deliver_name`, `target_arn`, `deliver_prefix`, `deliver_type`,
  `deliver_content_type`; expose the read-only field `create_time`.
- Delete semantics are **turn-off** (开关启用型): destroy sets `Status = 0` on the cloud side, because
  leaving delivery enabled after `terraform destroy` would silently keep shipping audit logs.
- The resource is a singleton per account/region. Its state ID is a fixed sentinel token so that
  `terraform import` can be supported without inventing a fake cloud-side identifier.
- No change to any existing resource; this is purely additive.

Not a breaking change.

## Capabilities

### New Capabilities

- `config-update-config-deliver`: the `tencentcloud_config_update_config_deliver` resource — its schema
  (which API field each attribute maps to, requiredness, computed-ness), its read/update lifecycle,
  its singleton state-ID and import behaviour, and its turn-off-on-destroy semantics.

### Modified Capabilities

None. No existing capability's requirements change.

## Impact

- **New files**
  - `tencentcloud/services/config/resource_tc_config_update_config_deliver_config.go`
  - `tencentcloud/services/config/resource_tc_config_update_config_deliver_config.md`
  - `tencentcloud/services/config/resource_tc_config_update_config_deliver_config_test.go`
- **Modified files**
  - `tencentcloud/services/config/service_tencentcloud_config.go` — append framework-style
    `DescribeConfigDeliver` / `UpdateConfigDeliver` service methods.
  - `tencentcloud/framework/registry.go` — add the factory to `resourceFactories` and import the
    `config` service package.
  - `tencentcloud/provider.md` — register `tencentcloud_config_update_config_deliver` under the Config
    product nodes so `make doc` emits a website page.
  - `.changelog/` — changelog entry (produced in the finalize phase).
- **APIs**: `config:DescribeConfigDeliver`, `config:UpdateConfigDeliver`
  (`github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/config/v20220802`). Already vendored —
  no `go.mod` / `vendor/` change needed.
- **Dependencies**: reuses `tencentcloud/internal/sharedmeta.ProviderMeta` and
  `tencentcloud/framework/internal/helper`; no new external dependency.
- **Name collision check**: `tencentcloud_config_update_config_deliver` exists nowhere in the repo
  (`provider.go` count = 0). The pre-existing SDKv2 resource `tencentcloud_config_deliver_config`
  (same two APIs) is left untouched; see `design.md` for the relationship between the two.
