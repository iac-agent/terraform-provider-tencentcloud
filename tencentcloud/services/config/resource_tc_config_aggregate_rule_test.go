package config_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudConfigAggregateRuleResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfigAggregateRule,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_config_aggregate_rule.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_config_aggregate_rule.example", "config_rule_id"),
					resource.TestCheckResourceAttr("tencentcloud_config_aggregate_rule.example", "identifier", "cam-user-group-bound"),
					resource.TestCheckResourceAttr("tencentcloud_config_aggregate_rule.example", "risk_level", "3"),
				),
			},
			{
				Config: testAccConfigAggregateRuleUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_config_aggregate_rule.example", "rule_name", "tf-example-aggregate-update"),
					resource.TestCheckResourceAttr("tencentcloud_config_aggregate_rule.example", "risk_level", "2"),
				),
			},
			{
				ResourceName:      "tencentcloud_config_aggregate_rule.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

var testAccConfigAggregateRule = `
data "tencentcloud_config_list_aggregators" "example" {}

resource "tencentcloud_config_aggregate_rule" "example" {
  identifier       = "cam-user-group-bound"
  identifier_type  = "SYSTEM"
  rule_name        = "tf-example-aggregate-rule"
  resource_type    = ["QCS::CAM::User"]
  account_group_id = data.tencentcloud_config_list_aggregators.example.items[0].account_group_id
  risk_level       = 3
  description      = "tf example aggregate config rule"

  trigger_type {
    message_type                = "ScheduledNotification"
    maximum_execution_frequency = "TwentyFour_Hours"
  }
}
`

const testAccConfigAggregateRuleUpdate = `
data "tencentcloud_config_list_aggregators" "example" {}

resource "tencentcloud_config_aggregate_rule" "example" {
  identifier       = "cam-user-group-bound"
  identifier_type  = "SYSTEM"
  rule_name        = "tf-example-aggregate-update"
  resource_type    = ["QCS::CAM::User"]
  account_group_id = data.tencentcloud_config_list_aggregators.example.items[0].account_group_id
  risk_level       = 2
  description      = "tf example aggregate config rule update"

  trigger_type {
    message_type                = "ScheduledNotification"
    maximum_execution_frequency = "TwentyFour_Hours"
  }
}
`
