---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_function_1111_01"
sidebar_current: "docs-tencentcloud-resource-teo_function_1111_01"
description: |-
  Provides a resource to create and manage an edge function of TEO (Tencent Cloud EdgeOne).
---

# tencentcloud_teo_function_1111_01

Provides a resource to create and manage an edge function of TEO (Tencent Cloud EdgeOne).

The resource wraps the `CreateFunction`, `DescribeFunctions`, `ModifyFunction` and `DeleteFunction`
APIs. Only `content` and `remark` can be updated in place through `ModifyFunction`; changing
`zone_id` or `name` forces replacement. Note that the cloud keeps the existing remark when the
field is not sent to `ModifyFunction`, so removing `remark` from the configuration does not clear
the remote value.

Deleting the function only removes the edge function itself; other zone-level configurations are
not affected.

## Example Usage

```hcl
resource "tencentcloud_teo_function_1111_01" "example" {
  zone_id = "zone-2p1ghslb"
  name    = "edge-func-example"
  content = file("edge_function.js")
  remark  = "managed by terraform"
}
```

## Argument Reference

The following arguments are supported:

* `content` - (Required, String) Function content. Only JavaScript code is supported, up to 5MB.
* `name` - (Required, String) Function name. It can only contain lowercase letters, digits and hyphens, must start and end with a letter or digit, and supports up to 30 characters. Immutable: ModifyFunction cannot change it, so changing it forces resource replacement.
* `zone_id` - (Required, String) ID of the zone the function belongs to. Immutable: changing it forces resource replacement.
* `remark` - (Optional, String) Function description, up to 60 characters. The cloud keeps the existing remark when the field is not sent to ModifyFunction, so removing it from configuration does not clear the remote value.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

  * `reason` - Reason why the access restriction is enforced. Valid values: `ICP_RECORD_REQUIRED`, `GOVERNMENT_ORDER`.
  * `region` - Country/region code of the restricted access area, following the ISO 3166 standard.
* `create_time` - Creation time in UTC, following the ISO 8601 standard.
* `domain_compliance_restrictions` - List of region access restrictions on the default domain caused by compliance issues.
* `domain` - Default domain of the function.
* `function_id` - Function ID.
* `id` - Resource identifier. It is the function id.
* `update_time` - Last modification time in UTC, following the ISO 8601 standard.

## API Reference

This resource interacts with the following TencentCloud APIs:

* Create:
  - [CreateFunction](https://www.tencentcloud.com/document/api/1145/64815)
  - [DescribeFunctions](https://www.tencentcloud.com/document/api/1145/64809)
* Read: [DescribeFunctions](https://www.tencentcloud.com/document/api/1145/64809)
* Update:
  - [DescribeFunctions](https://www.tencentcloud.com/document/api/1145/64809)
  - [ModifyFunction](https://www.tencentcloud.com/document/api/1145/64807)
* Delete: [DeleteFunction](https://www.tencentcloud.com/document/api/1145/64813)

## Import

TEO edge functions can be imported using the `<zone_id>#<function_id>` format, e.g.

```
terraform import tencentcloud_teo_function_1111_01.example zone-2p1ghslb#ef-2p1ghslb-1234
```

