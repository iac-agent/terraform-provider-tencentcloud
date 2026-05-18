Provides a resource to create a TEO plan.

Example Usage

Create a personal plan with prepaid parameters

```hcl
resource "tencentcloud_teo_plan_v2" "example" {
  plan_type        = "personal"
  auto_use_voucher = "true"
  prepaid_plan_param {
    period    = 1
    renew_flag = "on"
  }
  renew_flag = "on"
}
```

Create an enterprise plan

```hcl
resource "tencentcloud_teo_plan_v2" "example" {
  plan_type = "enterprise"
}
```

Import

TEO plan can be imported using the id, e.g.

```
terraform import tencentcloud_teo_plan_v2.example edgeone-2unuvzjmmn2q
```
