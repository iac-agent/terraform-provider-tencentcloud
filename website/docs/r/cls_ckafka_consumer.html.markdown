---
subcategory: "Cloud Log Service(CLS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cls_ckafka_consumer"
sidebar_current: "docs-tencentcloud-resource-cls_ckafka_consumer"
description: |-
  Provides a resource to create a CLS ckafka consumer
---

# tencentcloud_cls_ckafka_consumer

Provides a resource to create a CLS ckafka consumer

## Example Usage

```hcl
resource "tencentcloud_cls_ckafka_consumer" "ckafka_consumer" {
  compression  = 1
  need_content = true
  topic_id     = "7e34a3a7-635e-4da8-9005-88106c1fde69"

  ckafka {
    instance_id   = "ckafka-qzoeaqx8"
    instance_name = "ckafka-instance"
    topic_id      = "topic-c6tm4kpm"
    topic_name    = "name"
    vip           = "172.16.112.23"
    vport         = "9092"
  }

  content {
    enable_tag = true
    meta_fields = [
      "__FILENAME__",
      "__HOSTNAME__",
      "__PKGID__",
      "__SOURCE__",
      "__TIMESTAMP__",
    ]
    tag_json_not_tiled = true
    timestamp_accuracy = 2
  }

  effective   = true
  role_arn    = "qcs::cam::uin/123456789:roleName/MyRole"
  external_id = "my-external-id"

  advanced_config {
    partition_hash_status = true
    partition_fields = [
      "__SOURCE__",
      "__HOSTNAME__",
    ]
  }
}
```

## Argument Reference

The following arguments are supported:

* `topic_id` - (Required, String, ForceNew) topic id.
* `advanced_config` - (Optional, List) advanced configuration for ckafka consumer.
* `ckafka` - (Optional, List) ckafka info.
* `compression` - (Optional, Int) compression method. 0 for NONE, 2 for SNAPPY, 3 for LZ4.
* `content` - (Optional, List) metadata information.
* `effective` - (Optional, Bool) whether the delivery task is effective.
* `external_id` - (Optional, String) external ID for role assumption.
* `need_content` - (Optional, Bool) whether to deliver the metadata information of the log.
* `role_arn` - (Optional, String) role access description name for cross-account access.

The `advanced_config` object supports the following:

* `partition_fields` - (Optional, Set) fields for partition hash, max 5 fields.
* `partition_hash_status` - (Optional, Bool) whether to enable partition hash. true: enable; false: disable.

The `ckafka` object supports the following:

* `instance_id` - (Required, String) instance id.
* `instance_name` - (Required, String) instance name.
* `topic_id` - (Required, String) topic id of ckafka.
* `topic_name` - (Required, String) topic name of ckafka.
* `vip` - (Required, String) vip.
* `vport` - (Required, String) vport.

The `content` object supports the following:

* `enable_tag` - (Required, Bool) whether to deliver the TAG info.
* `meta_fields` - (Required, Set) metadata info list.
* `tag_json_not_tiled` - (Optional, Bool) whether to tiling tag json.
* `timestamp_accuracy` - (Optional, Int) delivery timestamp precision,1 for second, 2 for millisecond.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

CLS ckafka consumer can be imported using the id, e.g.

```
terraform import tencentcloud_cls_ckafka_consumer.ckafka_consumer topic_id
```

