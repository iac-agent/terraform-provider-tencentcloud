## 1. Resource Implementation

- [x] 1.1 Create `tencentcloud/services/teo/resource_tc_teo_function_replica_v4.go` with `ResourceTencentCloudTeoFunctionReplicaV4()` schema definition: `zone_id`/`function_id`/`replica_name` (Required+ForceNew, TypeString), `content` (Required, TypeString), `remark` (Optional, TypeString), `created_on`/`modified_on` (Computed, TypeString), plus Importer (ImportStatePassthrough)
- [x] 1.2 Implement `resourceTencentCloudTeoFunctionReplicaV4Create`: fill `CreateFunctionReplicaRequest` via `d.GetOk`, wrap `CreateFunctionReplicaWithContext` in `resource.Retry(tccommon.WriteRetryTimeout, ...)` with `tccommon.RetryError(e)` wrapping, return NonRetryableError on nil result/Response, then `d.SetId(strings.Join([]string{zoneId, functionId, replicaName}, tccommon.FILED_SP))` outside the retry block, finally invoke Read
- [x] 1.3 Implement `resourceTencentCloudTeoFunctionReplicaV4Read`: split 3-segment ID (`len != 3` → error), build `DescribeFunctionReplicasRequest` with Filters `replica-name=[replicaName]` and `Limit=helper.Int64(200)`, wrap in `resource.Retry(tccommon.ReadRetryTimeout, ...)`, exact-match ReplicaName in response, set zone_id/function_id/replica_name and nil-guarded content/remark/created_on/modified_on; on not-found log `[CRUD]` line with `d.Id()` BEFORE `d.SetId("")`
- [x] 1.4 Implement `resourceTencentCloudTeoFunctionReplicaV4Update`: split ID, `mutableArgs := []string{"content", "remark"}` needChange check, call `ModifyFunctionReplicaWithContext` (ZoneId/FunctionId/ReplicaName from ID + Content/Remark when set) in `resource.Retry(tccommon.WriteRetryTimeout, ...)`, then invoke Read
- [x] 1.5 Implement `resourceTencentCloudTeoFunctionReplicaV4Delete`: split ID, `request.ReplicaNames = []*string{helper.String(replicaName)}` (single element), call `DeleteFunctionReplicaWithContext` in `resource.Retry(tccommon.WriteRetryTimeout, ...)`

## 2. Provider Registration

- [x] 2.1 Register `"tencentcloud_teo_function_replica_v4": teo.ResourceTencentCloudTeoFunctionReplicaV4()` in `tencentcloud/provider.go` (adjacent to the existing `tencentcloud_teo_function_replica` line)
- [x] 2.2 Add `tencentcloud_teo_function_replica_v4` entry to the TEO Resources list in `tencentcloud/provider.md` (adjacent to the v1 entry)

## 3. Unit Tests

- [x] 3.1 Create `tencentcloud/services/teo/resource_tc_teo_function_replica_v4_test.go` (package teo_test) with gomonkey helpers using v4-suffixed names (e.g. `mockMetaFunctionReplicaV4`, `ptrStringFunctionReplicaV4`) to avoid conflicts with the existing v1 test file in the same package
- [x] 3.2 Add `TestTeoFunctionReplicaV4_Create` covering creation with all parameters (assert request fields, composite ID result) and creation without remark
- [x] 3.3 Add `TestTeoFunctionReplicaV4_Read` asserting content/remark/created_on/modified_on population, and `TestTeoFunctionReplicaV4_ReadNotFound` asserting `d.SetId("")` on empty response
- [x] 3.4 Add `TestTeoFunctionReplicaV4_Update` asserting `ModifyFunctionReplicaWithContext` receives updated content/remark, and a no-change case asserting Modify is skipped
- [x] 3.5 Add `TestTeoFunctionReplicaV4_Delete` asserting `DeleteFunctionReplicaWithContext` receives a single-element ReplicaNames list

## 4. Documentation

- [x] 4.1 Create `tencentcloud/services/teo/resource_tc_teo_function_replica_v4.md` with one-line description mentioning TEO, Example Usage (zone_id/function_id/replica_name/content/remark), NOTE about not mixing with the v1 resource, and Import section documenting the `zone_id#function_id#replica_name` composite ID; no manual Argument/Attribute Reference sections

## 5. Verification

- [x] 5.1 Verify the new code compiles (no build errors) and check all function error returns are handled per project conventions
- [x] 5.2 Verify schema field ↔ cloud API parameter consistency: Create params exist in CreateFunctionReplicaRequest, Update params exist in ModifyFunctionReplicaRequest, Delete params exist in DeleteFunctionReplicaRequest, Read params exist in DescribeFunctionReplicasRequest
- [x] 5.3 Verify existing v1 resource files (`resource_tc_teo_function_replica.go` / `_test.go` / `.md`) and their provider registrations are untouched

## 6. Finalization (tfpacer-finalize skill 阶段执行)

- [ ] 6.1 Run `gofmt` formatting on the changed Go files
- [ ] 6.2 Run `make doc` to generate `website/docs/r/teo_function_replica_v4.html.markdown` and update provider docs
- [ ] 6.3 Create the `.changelog/` entry for the new resource
