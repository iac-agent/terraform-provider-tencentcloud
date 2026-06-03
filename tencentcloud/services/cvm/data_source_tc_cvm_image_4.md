Use this data source to query detailed information of CVM images

Example Usage

Query by image ids

```hcl
data "tencentcloud_cvm_image_4" "example" {
  image_ids = ["img-8toqc6s3"]
}
```

Query by filters

```hcl
data "tencentcloud_cvm_image_4" "example" {
  filters {
    name   = "image-type"
    values = ["PUBLIC_IMAGE"]
  }
}
```

Query by instance type

```hcl
data "tencentcloud_cvm_image_4" "example" {
  instance_type = "SA5.MEDIUM2"
  filters {
    name   = "image-type"
    values = ["PUBLIC_IMAGE"]
  }
}
```
