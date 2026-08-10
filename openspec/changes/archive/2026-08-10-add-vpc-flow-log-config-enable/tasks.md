## 1. Resource Implementation

- [x] 1.1 Verify resource schema includes `flow_log_id` (Required, TypeString, ForceNew) and `enable` (Required, TypeBool) parameters
- [x] 1.2 Verify Create function sets flow_log_id as resource ID and delegates to Update
- [x] 1.3 Verify Read function calls DescribeFlowLogs API with FlowLogId and syncs Enable field to state
- [x] 1.4 Verify Update function calls EnableFlowLogs (when enable=true) or DisableFlowLogs (when enable=false) with retry
- [x] 1.5 Verify Delete function is a no-op (removes from state without calling cloud API)
- [x] 1.6 Verify ImportStatePassthrough is configured for terraform import support

## 2. Provider Registration

- [x] 2.1 Verify resource is registered in tencentcloud/provider.go as `tencentcloud_vpc_flow_log_config`

## 3. Testing

- [x] 3.1 Verify unit tests cover Create with enable=true and enable=false scenarios
- [x] 3.2 Verify unit tests cover Read with Enable field present and nil
- [x] 3.3 Verify unit tests cover Update from enabled to disabled and vice versa
- [x] 3.4 Verify unit tests cover Delete (no-op)

## 4. Documentation

- [x] 4.1 Verify resource_tc_vpc_flow_log_config.md exists with correct Example Usage and Import sections
- [x] 4.2 Verify website docs are generated via make doc
