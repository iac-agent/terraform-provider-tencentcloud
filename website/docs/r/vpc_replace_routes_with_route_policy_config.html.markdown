---
subcategory: "Virtual Private Cloud(VPC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_vpc_replace_routes_with_route_policy_config"
sidebar_current: "docs-tencentcloud-resource-vpc_replace_routes_with_route_policy_config"
description: |-
  Provides a resource to create a VPC replace routes with route policy config
---

# tencentcloud_vpc_replace_routes_with_route_policy_config

Provides a resource to create a VPC replace routes with route policy config

## Example Usage

```hcl
resource "tencentcloud_vpc_replace_routes_with_route_policy_config" "example" {
  route_table_id = "rtb-olsbhnyc"
  routes {
    route_item_id      = "rti-araogi5t"
    force_match_policy = true
  }

  routes {
    route_item_id      = "rti-kiyt72op"
    force_match_policy = true
  }
}
```

### Read with custom filter name and values

```hcl
resource "tencentcloud_vpc_replace_routes_with_route_policy_config" "example" {
  route_table_id = "rtb-olsbhnyc"
  routes {
    route_item_id      = "rti-araogi5t"
    force_match_policy = true
  }

  name   = "route-table-name"
  values = ["my-rtb"]
}
```

### Read with explicit route_table_ids (mutually exclusive with filters)

```hcl
resource "tencentcloud_vpc_replace_routes_with_route_policy_config" "example" {
  route_table_id = "rtb-olsbhnyc"
  routes {
    route_item_id      = "rti-araogi5t"
    force_match_policy = true
  }

  route_table_ids = ["rtb-olsbhnyc"]
}
```

### Read with need_router_info and custom limit

```hcl
resource "tencentcloud_vpc_replace_routes_with_route_policy_config" "example" {
  route_table_id = "rtb-olsbhnyc"
  routes {
    route_item_id      = "rti-araogi5t"
    force_match_policy = true
  }

  need_router_info = false
  limit            = 50
}
```

## Argument Reference

The following arguments are supported:

* `route_table_id` - (Required, String, ForceNew) Route Table Instance ID.
* `routes` - (Required, Set) Routing policy object. requires specifying the unique ID of routing policy (RouteItemId).
* `limit` - (Optional, Int) Return quantity of `DescribeRouteTables`, default is 20, max value is 100. When unset, the read helper uses its internal default of 100.
* `name` - (Optional, String) Filter name of the `DescribeRouteTables` read request. Populates the `Name` of a `Filter` entry. Only takes effect when `route_table_ids` is not set.
* `need_router_info` - (Optional, Bool) Indicates whether to obtain route policy info. Maps to `NeedRouterInfo` of `DescribeRouteTables`.
* `route_table_ids` - (Optional, List: [`String`]) Route table instance IDs, e.g. rtb-azd4dt1c. Maps to `RouteTableIds` of `DescribeRouteTables`. Mutually exclusive with `Filters` (including the `route-table-id` filter derived from `route_table_id`); when set, the read path queries by `RouteTableIds` only.
* `values` - (Optional, List: [`String`]) Filter values of the `DescribeRouteTables` read request. Populates the `Values` of the same `Filter` entry as `name`. Only takes effect when `route_table_ids` is not set.

The `routes` object supports the following:

* `force_match_policy` - (Optional, Bool) Match the route reception policy tag.
* `route_item_id` - (Optional, String) Route unique policy ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



