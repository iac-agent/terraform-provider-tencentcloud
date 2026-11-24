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
			},
		},
	})
}

const testAccOceanusOceanusCluster = `
resource "tencentcloud_oceanus_oceanus_cluster" "example" {
  cluster_name       = "tf_example_cluster"
  region_id          = 1
  zone_id            = 100003
  login_password     = "Password@123"
  default_cos_bucket = "tf-example-bucket-1250000000"
  vpc_descriptions {
    vpc_id    = "vpc-xxxxxx"
    subnet_id = "subnet-xxxxxx"
  }
  instance_charge_type             = "POSTPAID_BY_SECOND"
  flink_web_ui_network_access_type = "NetworkAccess_EXTERNAL"
}
`
