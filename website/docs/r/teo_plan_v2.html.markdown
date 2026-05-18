---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_plan_v2"
sidebar_current: "docs-tencentcloud-resource-teo_plan_v2"
description: |-
  Provides a resource to create a TEO plan.
---

# tencentcloud_teo_plan_v2

Provides a resource to create a TEO plan.

## Example Usage

### Create a personal plan with prepaid parameters

```hcl
resource "tencentcloud_teo_plan_v2" "example" {
  plan_type        = "personal"
  auto_use_voucher = "true"
  prepaid_plan_param {
    period     = 1
    renew_flag = "on"
  }
  renew_flag = "on"
}
```

### Create an enterprise plan

```hcl
resource "tencentcloud_teo_plan_v2" "example" {
  plan_type = "enterprise"
}
```

## Argument Reference

The following arguments are supported:

* `plan_type` - (Required, String, ForceNew) The subscription package type, the possible values are: `personal`: personal package, prepaid package; `basic`: basic package, prepaid package; `standard`: standard package, prepaid package; `enterprise`: enterprise package, postpaid package.
* `auto_use_voucher` - (Optional, String, ForceNew) Whether to auto-use vouchers, the possible values are: `true`: yes; `false`: no. This parameter is only effective when PlanType is personal, basic, or standard. If not filled, the default value is false.
* `prepaid_plan_param` - (Optional, List, ForceNew) Subscription prepaid package parameters. When PlanType is personal, basic, or standard, this parameter is optional and is used to enter the subscription duration of the package and whether to enable automatic renewal. If this parameter is not filled in, the default subscription duration is 1 month and automatic renewal is not enabled.
* `renew_flag` - (Optional, String) Auto-renewal flag for ModifyPlan, the values are: `on`: turn on automatic renewal; `off`: do not turn on automatic renewal.

The `prepaid_plan_param` object supports the following:

* `period` - (Optional, Int) The subscription period of the prepaid package, in months, with possible values: 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 24, 36. If not filled in, the default value 1 is used.
* `renew_flag` - (Optional, String) The automatic renewal flag of the prepaid package, the values are: `on`: turn on automatic renewal; `off`: do not turn on automatic renewal. If not filled in, the default value off is used. When automatic renewal occurs, the default renewal period is 1 month.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `area` - Service area, possible values are: mainland: Mainland China; overseas: Worldwide (excluding Mainland China); global: Worldwide (including Mainland China).
* `deal_name` - Order number.
* `enabled_time` - The time when the package takes effect.
* `expired_time` - The expiration date of the package.
* `pay_mode` - Payment type, possible values: 0: post-payment; 1: pre-payment.
* `plan_id` - Plan ID.
* `status` - Package status, the values are: normal: normal status; expiring-soon: about to expire; expired: expired; isolated: isolated; overdue-isolated: overdue isolated.


## Import

TEO plan can be imported using the id, e.g.

```
terraform import tencentcloud_teo_plan_v2.example edgeone-2unuvzjmmn2q
```

