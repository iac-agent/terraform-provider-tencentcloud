---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_plans"
sidebar_current: "docs-tencentcloud-datasource-teo_plans"
description: |-
  Use this data source to query detailed information of TEO plans
---

# tencentcloud_teo_plans

Use this data source to query detailed information of TEO plans

## Example Usage

### Query all plans

```hcl
data "tencentcloud_teo_plans" "example" {}
```

### Query plans by filters

```hcl
data "tencentcloud_teo_plans" "example" {
  order     = "expire-time"
  direction = "desc"
  filters {
    name = "plan-id"
    values = [
      "edgeone-2o1xvpmq7nn",
      "edgeone-2mezmk9s2xdx"
    ]
  }

  filters {
    name = "plan-type"
    values = [
      "plan-trial",
      "plan-personal",
      "plan-basic",
      "plan-standard",
      "plan-enterprise"
    ]
  }

  filters {
    name = "area"
    values = [
      "overseas",
      "mainland",
      "global"
    ]
  }

  filters {
    name = "status"
    values = [
      "normal",
      "expiring-soon",
      "expired",
      "isolated"
    ]
  }
}
```

## Argument Reference

The following arguments are supported:

* `direction` - (Optional, String) Sorting direction, the possible values are: asc: sort from small to large; desc: sort from large to small. If not filled in, the default value desc will be used.
* `filters` - (Optional, List) Filter conditions, the upper limit of Filters.Values is 20. The detailed filtering conditions are as follows: plan-type filter by plan type (plan-trial, plan-personal, plan-basic, plan-standard, plan-enterprise); plan-id filter by plan ID (form: edgeone-268z103ob0sx); area filter by plan acceleration region (mainland, overseas, global); status filter by plan status (normal, expiring-soon, expired, isolated).
* `order` - (Optional, String) Sorting field, the values are: enable-time: effective time; expire-time: expiration time. If not filled in, the default value enable-time will be used.
* `result_output_file` - (Optional, String) Used to save results.

The `filters` object supports the following:

* `name` - (Required, String) Filter name.
* `values` - (Required, Set) Filter value.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `plans` - Plan list.
  * `acc_traffic_capacity` - The content acceleration traffic specifications in the package, unit: byte.
  * `area` - Service area, the values are: mainland: Mainland China; overseas: Worldwide (excluding Mainland China); global: Worldwide (including Mainland China).
  * `bindable` - Whether the package allows binding of new sites, the values are: true: allows binding of new sites; false: does not allow binding of new sites.
  * `cross_mlc_traffic_capacity` - The optimized traffic specifications of the Chinese mainland network in the package, unit: bytes.
  * `ddos_traffic_capacity` - DDoS protection traffic specifications within the package, unit: bytes.
  * `enabled_time` - The package effective time.
  * `expired_time` - The expiration date of the package.
  * `features` - The functions supported by the package, such as ContentAcceleration, SmartAcceleration, L4, Waf, QUIC, CrossMLC, ProcessMedia, L4DDoS and L7DDoS specifications.
  * `l4_traffic_capacity` - Layer 4 acceleration traffic specifications within the package, unit: byte.
  * `pay_mode` - Payment type, possible values: 0: post-payment; 1: pre-payment.
  * `plan_id` - Plan ID.
  * `plan_type` - Plan type. Possible values are: plan-trial: Trial plan; plan-personal: Personal plan; plan-basic: Basic plan; plan-standard: Standard plan; plan-enterprise-v2: Enterprise plan; plan-enterprise-model-a: Enterprise Model A plan; plan-enterprise: Old Enterprise plan.
  * `sec_request_capacity` - The number of secure requests in the package, unit: times.
  * `sec_traffic_capacity` - The security flow specification within the package, unit: byte.
  * `smart_request_capacity` - The number of intelligent acceleration requests in the package, unit: times.
  * `smart_traffic_capacity` - Smart acceleration traffic specifications within the package, unit: byte.
  * `status` - Package status, the values are: normal: normal status; expiring-soon: about to expire; expired: expired; isolated: isolated; overdue-isolated: overdue isolated.
  * `vau_capacity` - VAU specifications in the package, unit: piece.
  * `zones_info` - Site information bound to the package, including site ID, site name, and site status.
    * `paused` - Whether the site is disabled. The possible values are: false: not disabled; true: disabled.
    * `zone_id` - Zone ID.
    * `zone_name` - Zone name.

## API Reference

This data source queries the following TencentCloud APIs:

* Read: [DescribePlans](https://www.tencentcloud.com/document/api/1145/72310)

