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

Note: the account group is identified by its single `group_id`. During refresh the whole account group list is scanned page by page (20 entries per page) and the entry whose `group_id` equals the resource id is used; only after every page is exhausted without a match is the resource treated as deleted and removed from state. If a field cannot be written into the state, refresh fails with an error instead of silently leaving stale values. The name can be updated in place, while the description can only be set on creation, so changing it forces a new account group to be created. Deleting the account group only removes the group object itself; member accounts and policies that reference the group are not changed by this resource.