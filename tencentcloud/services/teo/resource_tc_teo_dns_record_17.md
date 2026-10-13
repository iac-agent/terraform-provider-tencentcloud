Provides a resource to create a teo teo_dns_record_17

Example Usage

```hcl
resource "tencentcloud_teo_dns_record_17" "teo_dns_record_17" {
  zone_id  = "zone-39quuimqg8r6"
  type     = "A"
  content  = "1.2.3.5"
  location = "Default"
  name     = "a.makn.cn"
  priority = 5
  ttl      = 300
  weight   = -1
  status   = "enable"
  # Optional pagination limit for the DescribeDnsRecords api call used by read,
  # range: 1-1000, defaults to 1000 (the cloud api max value) when not set.
  Limit = 1000
}
```

Import

teo teo_dns_record_17 can be imported using the composite id (joined with `#`, in the form `{zoneId}#{recordId}`), e.g.

```
terraform import tencentcloud_teo_dns_record_17.teo_dns_record_17 {zoneId}#{recordId}
```
