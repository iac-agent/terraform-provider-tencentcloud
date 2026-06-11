## 1. Service Layer

- [x] 1.1 Add `DescribeTeoFunctionReplicaById` method to `TeoService` in `tencentcloud/services/teo/service_tencentcloud_teo.go` - call `DescribeFunctionReplicas` with `AdvancedFilter{Name: "replica-name", Values: [replicaName], Fuzzy: false}`, `Limit=200`, return first matching `FunctionReplica` or nil

## 2. Resource Implementation

- [x] 2.1 Create `tencentcloud/services/teo/resource_tc_teo_function_replica.go` with schema definition (zone_id, function_id, replica_name as Required+ForceNew; content as Required; remark as Optional; created_on/modified_on as Computed)
- [x] 2.2 Implement `resourceTencentCloudTeoFunctionReplicaCreate` - call `CreateFunctionReplica`, set composite ID `zone_id#function_id#replica_name`, call Read
- [x] 2.3 Implement `resourceTencentCloudTeoFunctionReplicaRead` - split composite ID, call service `DescribeTeoFunctionReplicaById`, set schema fields with nil checks, handle not-found with log + `d.SetId("")`
- [x] 2.4 Implement `resourceTencentCloudTeoFunctionReplicaUpdate` - check `d.HasChange` for `content` and `remark`, call `ModifyFunctionReplica` with ZoneId/FunctionId/ReplicaName/changed fields, call Read
- [x] 2.5 Implement `resourceTencentCloudTeoFunctionReplicaDelete` - split composite ID, call `DeleteFunctionReplica` with `ReplicaNames: [replica_name]`
- [x] 2.6 Add `Importer` support with `schema.ImportStatePassthrough` for 3-part composite ID import

## 3. Provider Registration

- [x] 3.1 Register `tencentcloud_teo_function_replica` resource in `tencentcloud/provider.go` ResourcesMap
- [x] 3.2 Add resource entry in `tencentcloud/provider.md`

## 4. Documentation

- [x] 4.1 Create `tencentcloud/services/teo/resource_tc_teo_function_replica.md` with one-line description mentioning TEO, example usage, and import section (with composite ID format note)

## 5. Unit Tests

- [x] 5.1 Create `tencentcloud/services/teo/resource_tc_teo_function_replica_test.go` with gomonkey mock tests for Create, Read, Update, Delete functions
- [x] 5.2 Run `go test -gcflags=all=-l` on the test file to verify all tests pass
