package oceanus_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudOceanusOceanusClusterResource_basic -v
func TestAccTencentCloudOceanusOceanusClusterResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOceanusOceanusCluster,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_oceanus_oceanus_cluster.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_oceanus_oceanus_cluster.example", "cluster_name"),
				),
			},
			{
				ResourceName:      "tencentcloud_oceanus_oceanus_cluster.example",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"login_password",
				},
			},
		},
	})
}

const testAccOceanusOceanusCluster = `
resource "tencentcloud_oceanus_oceanus_cluster" "example" {
  cluster_name         = "tf-example-cluster"
  region_id            = 1
  zone_id              = 100003
  login_password       = "YourPassword123"
  default_cos_bucket   = "example-bucket-1250000000"
  cu                   = 12
  cu_memory            = 4
  instance_charge_type = "POSTPAID_BY_SECOND"

  vpc_descriptions {
    vpc_id    = "vpc-xxx"
    subnet_id = "subnet-xxx"
  }
}
`
