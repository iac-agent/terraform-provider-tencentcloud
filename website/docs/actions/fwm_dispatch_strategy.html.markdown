---
subcategory: "Firewall Manager(FWM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_fwm_dispatch_strategy"
sidebar_current: "docs-tencentcloud-action-fwm_dispatch_strategy"
description: |-
  Provides an action to dispatch FWM (Firewall Manager) rule configurations to the cloud via the DispatchStrategy API, and waits until the dispatch reaches a terminal status. `Dispatch` pushes the strategy definitions that are already configured (for example through `tencentcloud_fwm_strategy` and `tencentcloud_fwm_rule_group`) so that they take effect.
---

# tencentcloud_fwm_dispatch_strategy

Provides an action to dispatch FWM (Firewall Manager) rule configurations to the cloud via the DispatchStrategy API, and waits until the dispatch reaches a terminal status. `Dispatch` pushes the strategy definitions that are already configured (for example through `tencentcloud_fwm_strategy` and `tencentcloud_fwm_rule_group`) so that they take effect.

~> **NOTE:** Actions are supported in HashiCorp Terraform version 1.14 and later.

~> **NOTE:** This action only dispatches, it does not create, update or destroy any Terraform-managed resource and it exposes no attributes referenceable from other resources. Destroying or removing the block does not revert anything that was already dispatched: the rule configurations stay effective in the cloud.

~> **NOTE:** Every successful invocation triggers a new cloud-side dispatch, so repeating the invocation repeats the side effect. Terraform does not guarantee that the remote operation happens only once, especially when the action is triggered by `lifecycle.action_trigger`, which runs automatically during apply (including on drift-driven `after_update` events), or when `count` expands the same configuration. A failed invocation does not roll back the part that was already dispatched. `on_failure` accepts `halt` (default), `continue` and `taint`; it is **not** a destroy guard, because a failing action during a destroy walk is downgraded to a warning.

## Example Usage

### Dispatch the strategy configuration of a single product

```hcl
action "tencentcloud_fwm_dispatch_strategy" "example" {
  config {
    product = "cfw_edge_acl"
    status  = 1

    timeouts {
      invoke = "10m"
    }
  }
}
```

### Dispatch the strategy configuration of every product

```hcl
action "tencentcloud_fwm_dispatch_strategy" "all_products" {
  config {}
}
```

### Abort an ongoing dispatch

```hcl
action "tencentcloud_fwm_dispatch_strategy" "abort" {
  config {
    product = "cfw_edge_acl"
    status  = 2
  }
}
```

### Trigger the action automatically after a rule group is created

```hcl
resource "tencentcloud_fwm_rule_group" "example" {
  # ...

  lifecycle {
    action_trigger {
      events     = [after_create, after_update]
      actions    = [action.tencentcloud_fwm_dispatch_strategy.example]
      on_failure = halt
    }
  }
}
```

### Trigger the action explicitly

```hcl
shell
terraform apply - invoke = action.tencentcloud_fwm_dispatch_strategy.example
```

## Argument Reference

The following arguments are supported:

* `product` - (Optional, String) Product whose rule configurations are dispatched, for example `cfw_edge_acl`. Omit it to dispatch every product. Only strategies of the given scope are affected; strategies of other products are left untouched.
* `status` - (Optional, Int) Dispatch operation to send: `1` dispatches the configured strategies, `2` aborts an ongoing dispatch. Default is `1`. Aborting is not a rollback: strategies that were already dispatched stay in effect. When `2` is used the action does not wait for a successful dispatch.
* `timeouts` - (Optional, Object) The timeouts block allows you to specify the timeout for the invoke operation.

## API Reference

This action invokes the following TencentCloud APIs:

* Invoke:
  - [DescribeStrategyDispatchStatus](https://www.tencentcloud.com/search/?wd=DescribeStrategyDispatchStatus&type=document&lang=en)
  - [DispatchStrategy](https://www.tencentcloud.com/search/?wd=DispatchStrategy&type=document&lang=en)

