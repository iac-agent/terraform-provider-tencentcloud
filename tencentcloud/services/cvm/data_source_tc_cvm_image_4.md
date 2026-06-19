Use this data source to query detailed information of CVM image

Example Usage

```hcl
data "tencentcloud_cvm_image_4" "example" {
  image_ids = ["img-0elsru2u"]

  filters {
    name   = "image-type"
    values = ["PRIVATE_IMAGE"]
  }
}
```
