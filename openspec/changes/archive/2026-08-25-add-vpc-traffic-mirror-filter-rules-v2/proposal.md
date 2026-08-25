## Why

TencentCloud VPC traffic mirror now supports a new "NO-DIRECTION" mode that requires the use of `CreateTrafficMirrorFilterRules` / `ModifyTrafficMirrorFilterRules` / `DeleteTrafficMirrorFilterRules` / `DescribeTrafficMirrorFilterRules` APIs to manage filter rules with direction, priority, and individual editing capabilities. The existing Terraform provider lacks resource support for these new filtering rule APIs, preventing users from managing traffic mirror filter rules via IaC.

## What Changes

- Add a new Terraform resource `tencentcloud_vpc_traffic_mirror_filter_rules_v2` to manage traffic mirror filter rules (ingress and egress) under a given traffic mirror instance.
- The resource supports full CRUD lifecycle: create, read, update, and delete filter rules.
- Supports both ingress and egress filter rules with fields: src_net, dst_net, protocol, src_port, dst_port, traffic_mirror_filter_rule_id, priority, action, description, created_time.
- Resource ID is composed of `traffic_mirror_id` concatenated with all rule IDs for tracking.

## Capabilities

### New Capabilities

- `vpc-traffic-mirror-filter-rules-v2`: Manage VPC traffic mirror filter rules (both ingress and egress) with full CRUD operations via the new traffic mirror filter rules APIs.

### Modified Capabilities

None.

## Impact

- New file: `tencentcloud/services/vpc/resource_tc_vpc_traffic_mirror_filter_rules_v2.go`
- New test file: `tencentcloud/services/vpc/resource_tc_vpc_traffic_mirror_filter_rules_v2_test.go`
- New doc file: `tencentcloud/services/vpc/resource_tc_vpc_traffic_mirror_filter_rules_v2.md`
- Modified: `tencentcloud/provider.go` (register new resource)
- Modified: `tencentcloud/provider.md` (register new resource in docs)
- Depends on existing SDK: `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312`