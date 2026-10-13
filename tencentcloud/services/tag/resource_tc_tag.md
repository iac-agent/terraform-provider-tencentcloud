Provides a resource to create a Tag

Example Usage

```hcl
resource "tencentcloud_tag" "example" {
  tag_key   = "tagKey"
  tag_value = "tagValue"
}
```

Argument Reference

The following arguments are supported:

* `tag_key` - (Required, ForceNew) Tag key.
* `tag_value` - (Required, ForceNew) Tag value.

Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `Category` - Tag category, valid values: `Custom` (custom tag), `System` (system tag), `All` (all tags). This value is returned by the cloud platform and cannot be set by users.

Import

Tag can be imported using the tagKey#tagValue, e.g.

```
terraform import tencentcloud_tag.example tagKey#tagValue
```