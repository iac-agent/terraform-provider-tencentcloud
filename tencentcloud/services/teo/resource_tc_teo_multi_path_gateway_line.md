Provides a resource to create a TEO multi-path gateway line

Example Usage

custom line

```hcl
resource "tencentcloud_teo_multi_path_gateway_line" "custom_line" {
  zone_id      = "zone-2xkazzl8yf6k"
  gateway_id   = "gw-3lchxitnb5pb"
  line_type    = "custom"
  line_address = "1.2.3.4:8080"
}
```

proxy line

```hcl
resource "tencentcloud_teo_multi_path_gateway_line" "proxy_line" {
  zone_id      = "zone-2xkazzl8yf6k"
  gateway_id   = "gw-3lchxitnb5pb"
  line_type    = "proxy"
  line_address = "1.2.3.4:443"
  proxy_id     = "sid-2xzwkzljmm9b"
  rule_id      = "rule-2xzwkzljmm9b"
}
```

Import

TEO multi-path gateway line can be imported using the zoneId#gatewayId#lineId, e.g.

```
terraform import tencentcloud_teo_multi_path_gateway_line.example zone-2xkazzl8yf6k#gw-3lchxitnb5pb#line-2
```
