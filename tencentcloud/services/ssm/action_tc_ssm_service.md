Provides an action to activate TencentCloud Secrets Manager (SSM) for the account and region of the current provider configuration via the CreateService API, and then wait until GetServiceStatus reports that the service is enabled. SSM is an account-level, per-region service, so every other SSM resource fails until it is activated. This is a one-time trigger; no Terraform state is persisted after the action completes.

~> **NOTE:** Actions are supported in HashiCorp Terraform version 1.14 and later, and the action must be run explicitly with `terraform apply -invoke=action.tencentcloud_ssm_service.example`.

~> **NOTE:** Each invoke places a pay-as-you-go activation order, so repeated invokes re-trigger the cloud call. Pin `client_token` to make retries replayable. Destroying the configuration does not deactivate SSM.

Example Usage

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
