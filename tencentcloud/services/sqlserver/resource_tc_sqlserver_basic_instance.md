Provides a SQL Server instance resource to create basic database instances.

Example Usage

```hcl
data "tencentcloud_availability_zones_by_product" "zones" {
  product = "sqlserver"
}

resource "tencentcloud_vpc" "vpc" {
  name       = "vpc-example"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_subnet" "subnet" {
  availability_zone = data.tencentcloud_availability_zones_by_product.zones.zones.4.name
  name              = "subnet-example"
  vpc_id            = tencentcloud_vpc.vpc.id
  cidr_block        = "10.0.0.0/16"
  is_multicast      = false
}

resource "tencentcloud_security_group" "security_group" {
  name        = "sg-example"
  description = "desc."
}

resource "tencentcloud_sqlserver_basic_instance" "example" {
  name                   = "tf-example"
  availability_zone      = data.tencentcloud_availability_zones_by_product.zones.zones.4.name
  charge_type            = "POSTPAID_BY_HOUR"
  vpc_id                 = tencentcloud_vpc.vpc.id
  subnet_id              = tencentcloud_subnet.subnet.id
  project_id             = 0
  memory                 = 4
  storage                = 100
  cpu                    = 2
  machine_type           = "CLOUD_PREMIUM"
  maintenance_week_set   = [1, 2, 3]
  maintenance_start_time = "09:00"
  maintenance_time_span  = 3
  security_groups        = [tencentcloud_security_group.security_group.id]

  tags = {
    "test" = "test"
  }
}
```

Example with custom timezone:

```hcl
resource "tencentcloud_sqlserver_basic_instance" "example_timezone" {
  name                   = "tf-example-utc"
  availability_zone      = data.tencentcloud_availability_zones_by_product.zones.zones.4.name
  charge_type            = "POSTPAID_BY_HOUR"
  vpc_id                 = tencentcloud_vpc.vpc.id
  subnet_id              = tencentcloud_subnet.subnet.id
  memory                 = 4
  storage                = 100
  cpu                    = 2
  machine_type           = "CLOUD_PREMIUM"
  time_zone              = "UTC"
}
```

Example with disk encryption enabled:

```hcl
resource "tencentcloud_sqlserver_basic_instance" "example_encrypted" {
  name                   = "tf-example-encrypted"
  availability_zone      = data.tencentcloud_availability_zones_by_product.zones.zones.4.name
  charge_type            = "POSTPAID_BY_HOUR"
  vpc_id                 = tencentcloud_vpc.vpc.id
  subnet_id              = tencentcloud_subnet.subnet.id
  memory                 = 4
  storage                = 100
  cpu                    = 2
  machine_type           = "CLOUD_SSD"
  disk_encrypt_flag      = 1
  time_zone              = "China Standard Time"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Name of the SQL Server basic instance.
* `cpu` - (Required) The CPU number of the SQL Server basic instance.
* `storage` - (Required) Disk size (in GB). Allowed value must be a multiple of 10.
* `memory` - (Required) Memory size (in GB).
* `machine_type` - (Required) The host type of the purchased instance. Valid values: `CLOUD_PREMIUM`, `CLOUD_SSD`, `CLOUD_HSSD`, `CLOUD_BSSD`.
* `charge_type` - (Optional, ForceNew) Pay type of the SQL Server basic instance. Only `POSTPAID_BY_HOUR` is valid.
* `vpc_id` - (Optional, ForceNew) ID of VPC.
* `subnet_id` - (Optional, ForceNew) ID of subnet.
* `engine_version` - (Optional, ForceNew) Version of the SQL Server basic database engine. Default is `2008R2`.
* `time_zone` - (Optional, ForceNew, Computed) System timezone for the SQL Server instance. Default is `China Standard Time`. Common values: `China Standard Time`, `UTC`, `Eastern Standard Time`. This setting cannot be changed after creation.
* `disk_encrypt_flag` - (Optional, ForceNew, Computed) Disk encryption flag. `0` - Disabled (default), `1` - Enabled. Disk encryption cannot be changed after instance creation.
* `period` - (Optional) Purchase instance period, the default value is 1. The value does not exceed 48.
* `security_groups` - (Optional) Security group bound to the instance.
* `auto_renew` - (Optional) Automatic renewal sign. 0 for normal renewal, 1 for automatic renewal.
* `auto_voucher` - (Optional) Whether to use the voucher automatically. 1 for yes, 0 for no.
* `voucher_ids` - (Optional) An array of voucher IDs.
* `maintenance_week_set` - (Optional, Computed) A list of integer indicates weekly maintenance.
* `maintenance_start_time` - (Optional, Computed) Start time of the maintenance in one day, format like `HH:mm`.
* `maintenance_time_span` - (Optional, Computed) The timespan of maintenance in one day, unit is hour.
* `project_id` - (Optional, Computed) Project ID, default value is 0.
* `availability_zone` - (Optional, ForceNew, Computed) Availability zone.
* `collation` - (Optional) System character set sorting rule, default: Chinese_PRC_CI_AS.
* `tags` - (Optional) The tags of the SQL Server basic instance.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `vip` - IP for private access.
* `vport` - Port for private access.
* `create_time` - Create time of the SQL Server basic instance.
* `status` - Status of the SQL Server basic instance.
* `dns_pod_domain` - Internet address domain name.
* `tgw_wan_vport` - External port number.

## Schema Change Behavior

The following parameters cannot be changed after instance creation and will trigger instance recreation:
* `time_zone`
* `disk_encrypt_flag`
* `engine_version`
* `charge_type`
* `vpc_id`
* `subnet_id`

Import

SQL Server basic instance can be imported using the id, e.g.

```
$ terraform import tencentcloud_sqlserver_basic_instance.example mssql-3cdq7kx5
```

After import, both `time_zone` and `disk_encrypt_flag` will be populated from the API.