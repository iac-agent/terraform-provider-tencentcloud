Provides a resource to create a TEO plan.

Example Usage

```hcl
resource "tencentcloud_teo_plan" "example" {
  plan_type         = "personal"
  auto_use_voucher = "true"

  prepaid_plan_param {
    period    = 1
    renew_flag = "off"
  }
}
```

Import

TEO plan can be imported using the plan id, e.g.

```
terraform import tencentcloud_teo_plan.example edgeone-2unuvzjmmn2q
```
