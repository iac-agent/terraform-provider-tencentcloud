---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_origin_group_v2"
sidebar_current: "docs-tencentcloud-resource-teo_origin_group_v2"
description: |-
  Provides a resource to create a TEO origin group.
---

# tencentcloud_teo_origin_group_v2

Provides a resource to create a TEO origin group.

~> **NOTE:** This is the framework implementation of the TEO origin group resource. `zone_id` forces a resource replacement when changed; `name`, `type`, `host_header` and `records` are updated in place.

## Example Usage

### IP domain origin group

```hcl
resource "tencentcloud_teo_origin_group_v2" "example" {
  zone_id = "zone-297z8rf93cfw"
  name    = "tf-test-origin-group"
  type    = "GENERAL"

  records {
    record = "tf-teo.xyz"
    type   = "IP_DOMAIN"
    weight = 100
  }
  records {
    record = "21.1.1.1"
    type   = "IP_DOMAIN"
    weight = 100
  }
}
```

### Object storage origin group with private authentication

```hcl
resource "tencentcloud_teo_origin_group_v2" "example" {
  zone_id     = "zone-297z8rf93cfw"
  name        = "tf-test-origin-group"
  type        = "HTTP"
  host_header = "origin.example.com"

  records {
    record  = "bucket-1250000000.cos.ap-guangzhou.myqcloud.com"
    type    = "COS"
    private = true

    private_parameters {
      name  = "AccessKeyId"
      value = "AKIDexample"
    }
    private_parameters {
      name  = "SecretAccessKey"
      value = "secret-example"
    }
    private_parameters {
      name  = "SignatureVersion"
      value = "v4"
    }
    private_parameters {
      name  = "Region"
      value = "ap-guangzhou"
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `records` - (Required, List) Origin records of the origin group. The whole list is sent to the API on update, so a change replaces the whole record set.
* `type` - (Required, String) Origin group type. Valid values: `GENERAL`: universal origin group, which only supports IP/domain origins and can be referenced by domain service, rule engine, layer-4 proxy, general load balancing and HTTP-specific load balancing; `HTTP`: HTTP-specific origin group, which supports IP/domain and object storage origins and cannot be referenced by the layer-4 proxy.
* `zone_id` - (Required, String) Site ID.
* `host_header` - (Optional, String) Back-to-origin Host Header, which only takes effect when type is HTTP. The rule engine Host Header configuration has a higher priority than the origin group Host Header.
* `name` - (Optional, String) Origin group name. It contains 1 to 200 characters, and only supports a - z, A - Z, 0 - 9, _, -. When omitted in the configuration, keep the original name: a null value in the configuration is sent as an explicit empty string, which keeps the group name unchanged; remove the attribute to keep the current name.

The `records` object supports the following:

* `record` - (Required, String) Origin record value without the port. It can be an IPv4, IPv6 or domain name.
* `private_parameters` - (Optional, List) Private authentication parameters, which take effect only when private is true.
* `private` - (Optional, Bool) Whether to use private authentication. It takes effect when the origin type is COS or AWS_S3. Valid values: `true`: use private authentication; `false`: not use.
* `record_id` - (Optional, String) Origin record ID, which is server-assigned and back-filled by Read.
* `type` - (Optional, String) Origin type. Valid values: `IP_DOMAIN`: IPv4, IPv6 or domain name origin; `COS`: COS origin; `AWS_S3`: AWS S3 object storage origin.
* `weight` - (Optional, Int) Origin weight for traffic scheduling, ranging from 0 to 100 (integer). Leave it unset and the system schedules by the default policy; set it to 0 and no traffic is scheduled to this origin (at least one origin must have a weight greater than 0). A null value in the configuration means "not set".

The `private_parameters` object of `records` supports the following:

* `name` - (Required, String) Private authentication parameter name. Valid values: `AccessKeyId`: Access Key ID; `SecretAccessKey`: Secret Access Key; `SignatureVersion`: signature version, v2 or v4; `Region`: bucket region.
* `value` - (Required, String) Private authentication parameter value.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

  * `alias_zone_name` - Zone alias of the referenced instance.
  * `instance_id` - Instance ID of the reference type.
  * `instance_name` - Instance name of the reference type.
  * `instance_type` - Reference service type. Valid values: `acceleration-domain`; `rule-engine`; `load-balancer`; `application-proxy`.
  * `zone_id` - Zone ID of the referenced instance.
  * `zone_name` - Zone name of the referenced instance.
* `create_time` - Origin group creation time.
* `id` - Resource identifier in the form `<zone_id>#<group_id>`.
* `references` - List of instances referencing the origin group.
* `update_time` - Origin group update time.

## API Reference

This resource interacts with the following TencentCloud APIs:

* Create:
  - [CreateOriginGroup](https://www.tencentcloud.com/document/api/1145/50974)
  - [DescribeOriginGroup](https://www.tencentcloud.com/document/api/1145/50970)
* Read: [DescribeOriginGroup](https://www.tencentcloud.com/document/api/1145/50970)
* Update:
  - [DescribeOriginGroup](https://www.tencentcloud.com/document/api/1145/50970)
  - [ModifyOriginGroup](https://www.tencentcloud.com/document/api/1145/50967)
* Delete: [DeleteOriginGroup](https://www.tencentcloud.com/document/api/1145/50972)

## Import

TEO origin groups can be imported using the `zone_id#group_id`, e.g.

```
terraform import tencentcloud_teo_origin_group_v2.example zone-297z8rf93cfw#origin-4f8a30b2-3720-11ed-b66b-525400dceb86
```

