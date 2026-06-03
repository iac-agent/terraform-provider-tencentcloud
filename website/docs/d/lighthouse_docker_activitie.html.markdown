---
subcategory: "TencentCloud Lighthouse(Lighthouse)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_lighthouse_docker_activitie"
sidebar_current: "docs-tencentcloud-datasource-lighthouse_docker_activitie"
description: |-
  Use this data source to query Docker activities for a Lighthouse instance.
---

# tencentcloud_lighthouse_docker_activitie

Use this data source to query Docker activities for a Lighthouse instance.

## Example Usage

### Query Docker activities by instance ID

```hcl
data "tencentcloud_lighthouse_docker_activitie" "example" {
  instance_id = "lhins-12345678"
}
```

### Query Docker activities by instance ID and activity IDs

```hcl
data "tencentcloud_lighthouse_docker_activitie" "example" {
  instance_id  = "lhins-12345678"
  activity_ids = ["lhda-12345678", "lhda-87654321"]
}
```

### Query Docker activities by time range

```hcl
data "tencentcloud_lighthouse_docker_activitie" "example" {
  instance_id        = "lhins-12345678"
  created_time_begin = 1717200000
  created_time_end   = 1719800000
}
```

## Argument Reference

The following arguments are supported:

* `activity_ids` - (Optional, Set: [`String`]) Docker activity ID list. Can be obtained from the ActivityId field returned by the DescribeDockerActivities interface.
* `created_time_begin` - (Optional, Int) The start value of the activity creation time, timestamp in seconds.
* `created_time_end` - (Optional, Int) The end value of the activity creation time, timestamp in seconds.
* `instance_id` - (Optional, String) Instance ID. Can be obtained from the InstanceId field returned by the DescribeInstances interface.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `docker_activity_set` - Docker activity list.
  * `activity_command_output` - Activity command output, base64 encoded.
  * `activity_id` - Activity ID.
  * `activity_name` - Activity name.
  * `activity_state` - Activity state. Valid values: INIT, OPERATING, SUCCESS, FAILED.
  * `container_ids` - Container ID list.
  * `created_time` - Creation time according to ISO8601 standard. UTC time is used. Format is YYYY-MM-DDThh:mm:ssZ.
  * `end_time` - End time according to ISO8601 standard. UTC time is used. Format is YYYY-MM-DDThh:mm:ssZ.


