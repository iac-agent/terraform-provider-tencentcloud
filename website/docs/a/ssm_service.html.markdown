---
subcategory: "Secrets Manager(SSM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ssm_service"
sidebar_current: "docs-tencentcloud-action-ssm_service"
description: |-
  Provides an action to subscribe the current account to Secrets Manager (SSM).
---

# tencentcloud_ssm_service

Provides an action to subscribe the current account to Secrets Manager (SSM).

This action performs a one-shot, account-level operation: it keeps no Terraform state, every invoke
performs a real `CreateService` call against the cloud API, and no cloud API can unsubscribe, so
removing the action from the configuration (`terraform destroy` / `terraform apply`) does not roll the
subscription back.

The action is only executed when it is explicitly invoked, for example with
`terraform apply -invoke=action.tencentcloud_ssm_service.example`. Invoking an already subscribed
account is safe: the `ResourceInUse.ServiceAlreadyEnabled` response is treated as the desired end
state instead of an error.

## Example Usage

```hcl
resource "random_uuid" "token" {}

action "tencentcloud_ssm_service" "example" {
  pay_mode     = "PostPay"
  client_token = "ssm-${random_uuid.token.result}"
}
```

### Invoke the action:

```hcl
shell
terraform apply - invoke = action.tencentcloud_ssm_service.example
```

## Argument Reference

The following arguments are supported:

* `client_token` - (Required, String) Client-generated idempotency token of the subscription request, 1 to 64 visible ASCII characters. Reuse the same value when retrying the same logical creation, otherwise a duplicated billing order may be created.
* `pay_mode` - (Required, String) Billing mode of the SSM subscription. Only `PostPay` is supported (case-sensitive).


