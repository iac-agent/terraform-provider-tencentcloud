Provides a resource to switch CVM CHC physical server network mode.

Example Usage

```hcl
resource "tencentcloud_cvm_chc_network_mode" "example" {
  chc_ids     = ["chc-1a2b3c4d"]
  network_mode = "DEPLOY"
}
```
