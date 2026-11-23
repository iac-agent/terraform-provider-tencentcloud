package alb_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

var testAlbListenerLoadBalancerId = os.Getenv("TENCENTCLOUD_ALB_LOAD_BALANCER_ID")
var testAlbListenerTargetGroupId = os.Getenv("TENCENTCLOUD_ALB_TARGET_GROUP_ID")

func TestAccTencentCloudAlbListenerResource_basic(t *testing.T) {
	t.Parallel()

	if testAlbListenerLoadBalancerId == "" || testAlbListenerTargetGroupId == "" {
		t.Skip("TENCENTCLOUD_ALB_LOAD_BALANCER_ID and TENCENTCLOUD_ALB_TARGET_GROUP_ID are required")
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccAlbListenerBasic(testAlbListenerLoadBalancerId, testAlbListenerTargetGroupId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_alb_listener.example", "id"),
					resource.TestCheckResourceAttr("tencentcloud_alb_listener.example", "listener_port", "80"),
					resource.TestCheckResourceAttr("tencentcloud_alb_listener.example", "listener_protocol", "HTTP"),
					resource.TestCheckResourceAttr("tencentcloud_alb_listener.example", "listener_name", "tf-example"),
				),
			},
			{
				Config: testAccAlbListenerUpdate(testAlbListenerLoadBalancerId, testAlbListenerTargetGroupId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_alb_listener.example", "id"),
					resource.TestCheckResourceAttr("tencentcloud_alb_listener.example", "listener_name", "tf-example-update"),
				),
			},
		},
	})
}

func testAccAlbListenerBasic(loadBalancerId, targetGroupId string) string {
	return fmt.Sprintf(`
resource "tencentcloud_alb_listener" "example" {
  load_balancer_id  = "%s"
  listener_port     = 80
  listener_protocol = "HTTP"
  listener_name     = "tf-example"

  default_actions {
    type = "Forward"
    target_group_config {
      target_groups {
        target_group_id = "%s"
        weight          = 10
      }
    }
  }
}
`, loadBalancerId, targetGroupId)
}

func testAccAlbListenerUpdate(loadBalancerId, targetGroupId string) string {
	return fmt.Sprintf(`
resource "tencentcloud_alb_listener" "example" {
  load_balancer_id  = "%s"
  listener_port     = 80
  listener_protocol = "HTTP"
  listener_name     = "tf-example-update"
  request_timeout   = 120

  default_actions {
    type = "Forward"
    target_group_config {
      target_groups {
        target_group_id = "%s"
        weight          = 20
      }
    }
  }
}
`, loadBalancerId, targetGroupId)
}
