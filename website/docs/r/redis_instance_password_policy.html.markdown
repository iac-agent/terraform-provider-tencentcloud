---
subcategory: "TencentDB for Redis(crs)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_redis_instance_password_policy"
sidebar_current: "docs-tencentcloud-resource-redis_instance_password_policy"
description: |-
  Provides a resource to manage the password complexity policy of a TencentCloud Redis instance.
---

# tencentcloud_redis_instance_password_policy

Provides a resource to manage the password complexity policy of a TencentCloud Redis instance.

## Example Usage

```hcl
data "tencentcloud_redis_zone_config" "zone" {
  type_id = 7
}

resource "tencentcloud_vpc" "vpc" {
  cidr_block = "10.0.0.0/16"
  name       = "tf_redis_vpc"
}

resource "tencentcloud_subnet" "subnet" {
  vpc_id            = tencentcloud_vpc.vpc.id
  availability_zone = data.tencentcloud_redis_zone_config.zone.list[0].zone
  name              = "tf_redis_subnet"
  cidr_block        = "10.0.1.0/24"
}

resource "tencentcloud_redis_instance" "example" {
  availability_zone  = data.tencentcloud_redis_zone_config.zone.list[0].zone
  type_id            = data.tencentcloud_redis_zone_config.zone.list[0].type_id
  password           = "Password@123"
  mem_size           = 8192
  redis_shard_num    = data.tencentcloud_redis_zone_config.zone.list[0].redis_shard_nums[0]
  redis_replicas_num = data.tencentcloud_redis_zone_config.zone.list[0].redis_replicas_nums[0]
  name               = "tf_example"
  port               = 6379
  vpc_id             = tencentcloud_vpc.vpc.id
  subnet_id          = tencentcloud_subnet.subnet.id
}

resource "tencentcloud_redis_instance_password_policy" "example" {
  instance_id       = tencentcloud_redis_instance.example.id
  enabled           = true
  min_letter_count  = 2
  min_digit_count   = 1
  min_special_count = 1
  min_length        = 8
}
```

## Argument Reference

The following arguments are supported:

* `enabled` - (Required, Bool) Whether the instance-level password complexity policy is enabled.
* `instance_id` - (Required, String) The ID of instance.
* `min_digit_count` - (Optional, Int) Minimum number of digit characters. Range [1, 16], default 1.
* `min_length` - (Optional, Int) Minimum total password length. Range [8, 64], default 8.
* `min_letter_count` - (Optional, Int) Minimum number of letter (upper/lower case) characters. Range [1, 16], default 1.
* `min_special_count` - (Optional, Int) Minimum number of special characters. Range [1, 16], default 1.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

Redis instance password policy can be imported using the instance_id, e.g.

```
terraform import tencentcloud_redis_instance_password_policy.example crs-cqdfdzvt
```

