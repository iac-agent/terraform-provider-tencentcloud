Provides a resource to create a TEO multi-path gateway line.

Example Usage

```hcl
resource "tencentcloud_teo_multi_path_gateway_line" "example" {
  zone_id      = "zone-2qtuhspy7cr6"
  gateway_id   = "gw-abcdefgh"
  line_type    = "custom"
  line_address = "1.2.3.4:8080"
}
```

```hcl
resource "tencentcloud_teo_multi_path_gateway_line" "proxy_example" {
  zone_id      = "zone-2qtuhspy7cr6"
  gateway_id   = "gw-abcdefgh"
  line_type    = "proxy"
  line_address = "5.6.7.8:443"
  proxy_id     = "sid-xxxxxxxx"
  rule_id      = "rule-xxxxxxxx"
}
```

Import

TEO multi-path gateway line can be imported using the joint id "zone_id#gateway_id#line_id", e.g.

```
terraform import tencentcloud_teo_multi_path_gateway_line.example zone-2qtuhspy7cr6#gw-abcdefgh#line-2
```
