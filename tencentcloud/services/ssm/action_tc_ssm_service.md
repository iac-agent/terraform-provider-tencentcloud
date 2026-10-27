Provides an action to subscribe the current account to Secrets Manager (SSM).

This action performs a one-shot, account-level operation: it keeps no Terraform state, every invoke
performs a real `CreateService` call against the cloud API, and no cloud API can unsubscribe, so
removing the action from the configuration (`terraform destroy` / `terraform apply`) does not roll the
subscription back.

The action is only executed when it is explicitly invoked, for example with
`terraform apply -invoke=action.tencentcloud_ssm_service.example`. Invoking an already subscribed
account is safe: the `ResourceInUse.ServiceAlreadyEnabled` response is treated as the desired end
state instead of an error.

Example Usage

```hcl
resource "random_uuid" "token" {}

action "tencentcloud_ssm_service" "example" {
  pay_mode     = "PostPay"
  client_token = "ssm-${random_uuid.token.result}"
}
```

Invoke the action:

```shell
terraform apply -invoke=action.tencentcloud_ssm_service.example
```
