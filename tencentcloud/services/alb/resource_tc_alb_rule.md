Provides a resource to create and manage an ALB forwarding rule.

Destroy only deletes the forwarding rule identified by `rule_ids` in `DeleteRules`. The documentation does not state whether it also affects the owning listener, load balancer instance or target groups, so those objects are never deleted by this resource. `dry_run = true` only pre-checks a request and must not be used for a managed rule; a pre-check never creates or changes the cloud object. An Update is rejected when a previously configured optional attribute is removed, because the API does not document whether an omitted value keeps or clears the cloud setting.

Example Usage

```hcl
resource "tencentcloud_alb_rule" "example" {
  load_balancer_id = "alb-12345678"
  listener_id      = "lst-12345678"
  client_token     = "tf-example-rule-token"

  rules = [
    {
      priority  = 10
      rule_name = "tf-example"
      direction = "Request"

      actions = [
        {
          order = 1
          type  = "TargetGroup"

          target_group_config = {
            target_groups = [
              {
                target_group_id = "lbtg-12345678"
                weight          = 30
              },
            ]

            target_group_sticky_session = {
              enabled = true
              timeout = 1000
            }
          }
        },
      ]

      conditions = [
        {
          type        = "Path"
          path_config = ["/api"]

          host_config         = ["example.com"]
          method_config       = ["GET"]
          source_ip_config    = ["10.0.0.1/32"]
          query_string_config = jsonencode([])
        },
      ]

      tags = [
        {
          tag_key   = "env"
          tag_value = "test"
        },
      ]
    },
  ]
}
```

Import

ALB rules can be imported using the composite id `<load_balancer_id>#<listener_id>#<rule_id>`, e.g.

```
terraform import tencentcloud_alb_rule.example alb-12345678#lst-12345678#rule-12345678
```
