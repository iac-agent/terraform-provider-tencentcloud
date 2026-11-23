package alb_test

import (
	"context"
	"fmt"
	"testing"

	alb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/alb/v20251030"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccTencentCloudAlbListener_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckAlbListenerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAlbListener_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlbListenerExists("tencentcloud_alb_listener.example"),
					resource.TestCheckResourceAttrSet("tencentcloud_alb_listener.example", "load_balancer_id"),
					resource.TestCheckResourceAttr("tencentcloud_alb_listener.example", "listener_protocol", "HTTP"),
					resource.TestCheckResourceAttr("tencentcloud_alb_listener.example", "listener_port", "80"),
					resource.TestCheckResourceAttr("tencentcloud_alb_listener.example", "listener_name", "tf-example"),
				),
			},
			{
				ResourceName:      "tencentcloud_alb_listener.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckAlbListenerExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("ALB listener %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("ALB listener id is not set")
		}
		conn := tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn().UseAlbClient()
		items := rs.Primary.Attributes["load_balancer_id"]
		listenerId := rs.Primary.Attributes["listener_id"]
		request := alb.NewDescribeListenersRequest()
		request.LoadBalancerId = &items
		request.ListenerIds = []*string{&listenerId}
		result, err := conn.DescribeListenersWithContext(context.TODO(), request)
		if err != nil {
			return err
		}
		if result == nil || result.Response == nil || len(result.Response.Listeners) == 0 {
			return fmt.Errorf("ALB listener %s is not exist", rs.Primary.ID)
		}
		return nil
	}
}

func testAccCheckAlbListenerDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_alb_listener" {
			continue
		}
		conn := tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn().UseAlbClient()
		loadBalancerId := rs.Primary.Attributes["load_balancer_id"]
		listenerId := rs.Primary.Attributes["listener_id"]
		request := alb.NewDescribeListenersRequest()
		request.LoadBalancerId = &loadBalancerId
		request.ListenerIds = []*string{&listenerId}
		result, err := conn.DescribeListenersWithContext(context.TODO(), request)
		if err != nil {
			return err
		}
		if result != nil && result.Response != nil && len(result.Response.Listeners) > 0 {
			return fmt.Errorf("ALB listener %s still exists", rs.Primary.ID)
		}
	}
	return nil
}

const testAccAlbListener_basic = `
resource "tencentcloud_alb_listener" "example" {
  load_balancer_id  = var.load_balancer_id
  listener_protocol = "HTTP"
  listener_port     = 80
  listener_name     = "tf-example"

  default_actions {
    type = "Forward"
    target_group_config {
      target_groups {
        target_group_id = var.target_group_id
        weight          = 10
      }
    }
  }
}

variable "load_balancer_id" {
  default = "alb-xxxxxxxx"
}

variable "target_group_id" {
  default = "lbtg-xxxxxxxx"
}
`
