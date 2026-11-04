Provides a resource to create a TEO origin group.

~> **NOTE:** This is the framework implementation of the TEO origin group resource. `zone_id` forces a resource replacement when changed; `name`, `type`, `host_header` and `records` are updated in place.

Example Usage

IP domain origin group

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

Object storage origin group with private authentication

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

Import

TEO origin groups can be imported using the `zone_id#group_id`, e.g.

```
terraform import tencentcloud_teo_origin_group_v2.example zone-297z8rf93cfw#origin-4f8a30b2-3720-11ed-b66b-525400dceb86
```
