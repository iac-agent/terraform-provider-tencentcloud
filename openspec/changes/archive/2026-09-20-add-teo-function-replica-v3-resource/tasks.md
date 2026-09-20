## 1. Resource Implementation

- [x] 1.1 Create resource file `tencentcloud/services/teo/resource_tc_teo_function_replica_v3.go` defining `ResourceTencentCloudTeoFunctionReplicaV3()` with schema: `zone_id`/`function_id`/`replica_name` (Required, ForceNew, String), `content` (Required, String), `remark` (Optional, String), `sort_by`/`sort_order` (Optional, String), `filters` (Optional, TypeList of Resource with `name` Required String, `values` Required TypeList of String, `fuzzy` Optional Bool), `replica_names` (Required, TypeList of String), `created_on`/`modified_on` (Computed, String), `function_replicas` (Computed, flattened TypeList of Resource with function_id/replica_name/content/remark/created_on/modified_on, no extra wrapper layer), and `Importer: &schema.ResourceImporter{State: schema.ImportStatePassthrough}`
- [x] 1.2 Implement `resourceTencentCloudTeoFunctionReplicaV3Create`: build `CreateFunctionReplicaRequest` from schema (zone_id/function_id/replica_name/content/remark via `d.GetOk` + `helper.String`), call `CreateFunctionReplicaWithContext` inside `resource.Retry(tccommon.WriteRetryTimeout, ...)` with `tccommon.RetryError(e)` wrapping, check `result == nil || result.Response == nil` returning NonRetryableError, log logId before ID handling, then `d.SetId(strings.Join([]string{zoneId, functionId, replicaName}, tccommon.FILED_SP))` outside the retry block and call Read
- [x] 1.3 Implement `resourceTencentCloudTeoFunctionReplicaV3Read`: parse 3-segment ID (`id is broken, id is %s` error otherwise), call `DescribeFunctionReplicasWithContext` in `resource.Retry(tccommon.ReadRetryTimeout, ...)` with `Limit = helper.Int64(200)`, Filters containing `Name=replica-name, Values=[replicaName]`, plus user-configured sort_by/sort_order/filters mapped to SortBy/SortOrder/Filters; after retry, if response empty or no exact replica_name match, log `[CRUD]` line with resource id first then `d.SetId("")`; set zone_id/function_id/replica_name from ID and content/remark/created_on/modified_on (nil-guarded) plus flattened `function_replicas` computed list from the response
- [x] 1.4 Implement `resourceTencentCloudTeoFunctionReplicaV3Update`: parse 3-segment ID, when `d.HasChange` on `content` or `remark` (query-only fields sort_by/sort_order/filters excluded), build `ModifyFunctionReplicaRequest` with ZoneId/FunctionId/ReplicaName from ID and content/remark via `d.GetOk`, call `ModifyFunctionReplicaWithContext` in `resource.Retry(tccommon.WriteRetryTimeout, ...)`, then call Read
- [x] 1.5 Implement `resourceTencentCloudTeoFunctionReplicaV3Delete`: parse 3-segment ID, build `DeleteFunctionReplicaRequest` with ZoneId/FunctionId from ID and `ReplicaNames` from `d.Get("replica_names")` list (fall back to `[replicaName]` from ID when unset), call `DeleteFunctionReplicaWithContext` in `resource.Retry(tccommon.WriteRetryTimeout, ...)`, return nil on success (no re-read)

## 2. Provider Registration

- [x] 2.1 Register `"tencentcloud_teo_function_replica_v3": teo.ResourceTencentCloudTeoFunctionReplicaV3(),` in the TEO resources map of `tencentcloud/provider.go` (after `tencentcloud_teo_function_replica` entry, aligned with existing formatting)
- [x] 2.2 Add `tencentcloud_teo_function_replica_v3` entry to the TEO resources list in `tencentcloud/provider.md` (after `tencentcloud_teo_function_replica`)

## 3. Documentation

- [x] 3.1 Create `tencentcloud/services/teo/resource_tc_teo_function_replica_v3.md` with one-line description ("Provides a resource to create a TEO edge function replica"), Example Usage (hcl block covering required fields and optional sort/filters/replica_names usage), and Import section documenting the composite ID format `zone_id#function_id#replica_name`; do NOT include Argument Reference / Attribute Reference sections

## 4. Unit Tests

- [x] 4.1 Create `tencentcloud/services/teo/resource_tc_teo_function_replica_v3_test.go` (package `teo_test`) using gomonkey: mockMeta implementing `tccommon.ProviderMeta`, `patches.ApplyMethodReturn` on `UseTeoV20220901Client`, `patches.ApplyMethodFunc` mocks for CreateFunctionReplicaWithContext/DescribeFunctionReplicasWithContext/ModifyFunctionReplicaWithContext/DeleteFunctionReplicaWithContext, with tests covering Create (request field assertions + composite ID), Read found (fields set), Read not-found (ID cleared), Update (Modify request assertions), and Delete (ReplicaNames assertions); ensure the code compiles with correct error handling (do not run go test)

## 5. Verification (finalization stage only)

- [x] 5.1 Verify via `git diff` that only the intended files were created/modified (resource .go, _test.go, .md, provider.go, provider.md)
- [x] 5.2 Confirm gofmt/make doc/changelog actions are deferred to the tfpacer-finalize skill (not executed in implementation stage)
