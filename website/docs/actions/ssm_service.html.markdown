---
subcategory: "Secrets Manager(SSM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ssm_service"
sidebar_current: "docs-tencentcloud-action-ssm_service"
description: |-
  Provides an action to activate TencentCloud Secrets Manager (SSM) for the account and region of the current provider configuration via the CreateService API, and then wait until GetServiceStatus reports that the service is enabled. SSM is an account-level, per-region service, so every other SSM resource fails until it is activated. This is a one-time trigger; no Terraform state is persisted after the action completes.
---

# tencentcloud_ssm_service

Provides an action to activate TencentCloud Secrets Manager (SSM) for the account and region of the current provider configuration via the CreateService API, and then wait until GetServiceStatus reports that the service is enabled. SSM is an account-level, per-region service, so every other SSM resource fails until it is activated. This is a one-time trigger; no Terraform state is persisted after the action completes.

~> **NOTE:** Actions are supported in HashiCorp Terraform version 1.14 and later, and the action must be run explicitly with `terraform apply -invoke=action.tencentcloud_ssm_service.example`.

~> **NOTE:** Each invoke places a pay-as-you-go activation order, so repeated invokes re-trigger the cloud call. Pin `client_token` to make retries replayable. Destroying the configuration does not deactivate SSM.

## Example Usage

```hcl
action "tencentcloud_ssm_service" "example" {
  config {
    pay_mode = "PostPay"

    timeouts {
      invoke = "10m"
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `pay_mode` - (Required, String) Billing mode of the SSM activation order. Only `PostPay` (pay-as-you-go) is supported by the CreateService API and the value is case sensitive.
* `client_token` - (Optional, String) Client-generated identifier of the activation request, used by the billing side to deduplicate repeated creation requests. Must be 1 to 64 visible ASCII characters ('!'..'~'), with no space, non-ASCII or control character. When omitted, the action generates `ssm-<random token>` on every invoke; provide an explicit value to replay the same logical activation after a partial failure.
* `timeouts` - (Optional, Object) The timeouts block allows you to specify the timeout for the invoke operation.

## API Reference

This action invokes the following TencentCloud APIs:

* Invoke:
  - [CreateService](https://www.tencentcloud.com/search/?wd=CreateService&type=document&lang=en)
  - [GetServiceStatus](https://www.tencentcloud.com/document/api/1078/38648)

