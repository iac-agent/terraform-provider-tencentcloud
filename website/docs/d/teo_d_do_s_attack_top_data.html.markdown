---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_d_do_s_attack_top_data"
sidebar_current: "docs-tencentcloud-datasource-teo_d_do_s_attack_top_data"
description: |-
  Use this data source to query detailed information of TEO DDoS attack top data
---

# tencentcloud_teo_d_do_s_attack_top_data

Use this data source to query detailed information of TEO DDoS attack top data

## Example Usage

```hcl
data "tencentcloud_teo_d_do_s_attack_top_data" "example" {
  start_time  = "2024-01-01T00:00:00Z"
  end_time    = "2024-01-02T00:00:00Z"
  metric_name = "ddos_attackFlux_protocol"
  zone_ids    = ["zone-2qtuhspy7cr6"]
  area        = "overseas"
}
```

## Argument Reference

The following arguments are supported:

* `end_time` - (Required, String) End time of the query range. The query time range (EndTime - StartTime) must be less than or equal to 31 days.
* `metric_name` - (Required, String) The statistical metric to query. Valid values: `ddos_attackFlux_protocol`, `ddos_attackPackageNum_protocol`, `ddos_attackNum_attackType`, `ddos_attackNum_sregion`, `ddos_attackFlux_sip`, `ddos_attackFlux_sregion`.
* `start_time` - (Required, String) Start time of the query range.
* `area` - (Optional, String) Data area. Valid values: `overseas`, `mainland`. If not specified, the area is intelligently selected based on user location.
* `attack_type` - (Optional, String) Attack type filter. Valid values: `flood`, `icmpFlood`, `all`. Default is `all`.
* `policy_ids` - (Optional, Set: [`Int`]) Set of DDoS policy IDs to query. If not specified, all policy IDs are selected by default.
* `port` - (Optional, Int) Port number filter.
* `protocol_type` - (Optional, String) Protocol type filter. Valid values: `tcp`, `udp`, `all`. Default is `all`.
* `result_output_file` - (Optional, String) Used to save results.
* `zone_ids` - (Optional, Set: [`String`]) Set of zone IDs to query. Up to 100 zone IDs. Use `*` to query all zones under the master account.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `data` - DDoS attack Top data list.
  * `key` - The dimension value of the Top query.
  * `value` - The list of TopEntryValue items.
    * `count` - The ranking entity count.
    * `name` - The ranking entity name.


