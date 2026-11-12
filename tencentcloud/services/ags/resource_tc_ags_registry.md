Provides a resource to create an AGS (Agent Gateway Service) registry.

-> **Note:** All configurable attributes are immutable on the cloud side (there is no Modify API). Changing `name`, `description`, `approval_mode` or `tags` forces the registry to be recreated.

Example Usage

Create a registry with tags

```hcl
resource "tencentcloud_ags_registry" "example" {
  name          = "tf-example-registry"
  description   = "created by terraform"
  approval_mode = "AUTO"

  tags {
    key   = "env"
    value = "prod"
  }

  tags {
    key   = "team"
    value = "devops"
  }
}
```

Import

AGS registry can be imported using the id, e.g.

```sh
terraform import tencentcloud_ags_registry.example reg-abcd1234
```
