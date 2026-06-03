## 1. Service Layer Implementation

- [x] 1.1 Add `DescribeLighthouseDockerActivitiesByFilter` method to `LightHouseService` in `tencentcloud/services/lighthouse/service_tencentcloud_lighthouse.go`, implementing pagination with Offset/Limit (pageSize=100), mapping InstanceId, ActivityIds, CreatedTimeBegin, CreatedTimeEnd from param map to request, and returning `[]*lighthouse.DockerActivity`

## 2. Data Source Schema and Read Function

- [x] 2.1 Create `tencentcloud/services/lighthouse/data_source_tc_lighthouse_docker_activitie.go` with `DataSourceTencentCloudLighthouseDockerActivitie()` function defining schema: instance_id (Optional, TypeString), activity_ids (Optional, TypeSet of TypeString), created_time_begin (Optional, TypeInt), created_time_end (Optional, TypeInt), docker_activity_set (Computed, TypeList with nested schema), result_output_file (Optional, TypeString)
- [x] 2.2 Implement `dataSourceTencentCloudLighthouseDockerActivitieRead()` function: parse input params, call service layer with retry, map response DockerActivity fields (activity_id, activity_name, activity_state, activity_command_output, container_ids, created_time, end_time) with nil checks, set resource ID, handle result_output_file

## 3. Provider Registration

- [x] 3.1 Register `tencentcloud_lighthouse_docker_activitie` data source in `tencentcloud/provider.go` by adding to the dataSources map
- [x] 3.2 Add data source entry to `tencentcloud/provider.md`

## 4. Documentation

- [x] 4.1 Create `tencentcloud/services/lighthouse/data_source_tc_lighthouse_docker_activitie.md` with description, example usage, and import section per documentation guidelines

## 5. Unit Tests

- [x] 5.1 Create `tencentcloud/services/lighthouse/data_source_tc_lighthouse_docker_activitie_test.go` with gomonkey mock-based unit tests for the Read function, verifying parameter mapping and response field handling
- [x] 5.2 Run unit tests with `go test -gcflags=all=-l` to verify they pass

## 6. Verification

- [x] 6.1 Verify all new/modified files compile correctly (go vet check)
- [x] 6.2 Verify provider.go and provider.md are correctly updated
