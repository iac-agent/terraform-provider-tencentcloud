Provides a resource to create a CLS ckafka consumer

Example Usage

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
    enable_tag         = true
    meta_fields        = [
      "__FILENAME__",
      "__HOSTNAME__",
      "__PKGID__",
      "__SOURCE__",
      "__TIMESTAMP__",
    ]
    tag_json_not_tiled = true
    timestamp_accuracy = 2
  }

  effective    = true
  role_arn     = "qcs::cam::uin/123456789:roleName/MyRole"
  external_id  = "my-external-id"

  advanced_config {
    partition_hash_status = true
    partition_fields      = [
      "__SOURCE__",
      "__HOSTNAME__",
    ]
  }
}
```

Import

CLS ckafka consumer can be imported using the id, e.g.

```
terraform import tencentcloud_cls_ckafka_consumer.ckafka_consumer topic_id
```
