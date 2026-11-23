package alb_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudAlbListenerResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccAlbListener,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_alb_listener.listener", "id"),
					resource.TestCheckResourceAttr("tencentcloud_alb_listener.listener", "listener_port", "80"),
					resource.TestCheckResourceAttr("tencentcloud_alb_listener.listener", "listener_protocol", "HTTP"),
				),
			},
		},
	})
}

const testAccAlbListener = `
variable "load_balancer_id" {
  default = "alb-xxxxxxxx"
}
variable "target_group_id" {
  default = "lbtg-xxxxxxxx"
}
resource "tencentcloud_alb_listener" "listener" {
  load_balancer_id  = var.load_balancer_id
  listener_port     = 80
  listener_protocol = "HTTP"
  listener_name     = "tf-example-listener"
  default_actions {
    type = "TargetGroup"
    target_group_config {
      target_groups {
        target_group_id = var.target_group_id
        weight          = 10
      }
    }
  }
}
`
