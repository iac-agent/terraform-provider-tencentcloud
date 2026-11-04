Provides a resource to create a FWM account group

Example Usage

```hcl
resource "tencentcloud_fwm_account_group" "example" {
  group_name  = "tf-example"
  description = "remark."
}
```

Import

FWM account group can be imported using the GroupId, e.g.

```
terraform import tencentcloud_fwm_account_group.example acg-nyalbmbi
```

Note: deleting the account group only removes the group object itself; member accounts and policies that reference the group are not changed by this resource. The account group name can be updated in place, while the description can only be set on creation and changing it forces a new account group to be created.