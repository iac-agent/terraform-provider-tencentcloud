Provides a resource to create and manage an edge function of TEO (Tencent Cloud EdgeOne).

The resource wraps the `CreateFunction`, `DescribeFunctions`, `ModifyFunction` and `DeleteFunction`
APIs. Only `content` and `remark` can be updated in place through `ModifyFunction`; changing
`zone_id` or `name` forces replacement. Note that the cloud keeps the existing remark when the
field is not sent to `ModifyFunction`, so removing `remark` from the configuration does not clear
the remote value.

Deleting the function only removes the edge function itself; other zone-level configurations are
not affected.

Example Usage

```hcl
resource "tencentcloud_teo_function_1111_01" "example" {
  zone_id  = "zone-2p1ghslb"
  name     = "edge-func-example"
  content  = file("edge_function.js")
  remark   = "managed by terraform"
}
```

Import

TEO edge functions can be imported using the `<zone_id>#<function_id>` format, e.g.

```
terraform import tencentcloud_teo_function_1111_01.example zone-2p1ghslb#ef-2p1ghslb-1234
```
