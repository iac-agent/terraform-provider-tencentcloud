## 1. Resource Implementation

- [x] 1.1 Create `tencentcloud/services/vpc/resource_tc_vpc_traffic_mirror_filter_rules_v2.go` with resource schema definition (traffic_mirror_id Required/ForceNew, ingress_filter_rules Optional TypeList, egress_filter_rules Optional TypeList with nested fields src_net/dst_net/protocol/src_port/dst_port/traffic_mirror_filter_rule_id Computed/priority TypeInt/action/description/created_time Computed)
- [x] 1.2 Implement `resourceTencentCloudVpcTrafficMirrorFilterRulesV2Create` function calling `CreateTrafficMirrorFilterRules` API with `tccommon.ReadRetryTimeout` retry, validating non-empty response and setting `d.SetId()` to traffic_mirror_id
- [x] 1.3 Implement `resourceTencentCloudVpcTrafficMirrorFilterRulesV2Read` function calling `DescribeTrafficMirrorFilterRules` API, mapping response rules back to schema, handling nil/empty response with `d.SetId("")`
- [x] 1.4 Implement `resourceTencentCloudVpcTrafficMirrorFilterRulesV2Update` function calling `ModifyTrafficMirrorFilterRules` API with full rule set from configuration
- [x] 1.5 Implement `resourceTencentCloudVpcTrafficMirrorFilterRulesV2Delete` function calling `DeleteTrafficMirrorFilterRules` API with current rule IDs from state
- [x] 1.6 Add helper functions to convert between Terraform schema rule lists and SDK `TrafficMirrorFilter` structs

## 2. Provider Registration

- [x] 2.1 Register `tencentcloud_vpc_traffic_mirror_filter_rules_v2` resource in `tencentcloud/provider.go`
- [x] 2.2 Register `tencentcloud_vpc_traffic_mirror_filter_rules_v2` resource in `tencentcloud/provider.md`

## 3. Documentation

- [x] 3.1 Create `tencentcloud/services/vpc/resource_tc_vpc_traffic_mirror_filter_rules_v2.md` with one-line description, Example Usage, and Import section following gendoc format

## 4. Tests

- [x] 4.1 Create `tencentcloud/services/vpc/resource_tc_vpc_traffic_mirror_filter_rules_v2_test.go` with unit tests using gomonkey to mock cloud API calls (Create/Read/Update/Delete)
- [x] 4.2 Add test cases covering create with ingress+egress rules, read existing rules, update rules, delete rules, and empty response handling