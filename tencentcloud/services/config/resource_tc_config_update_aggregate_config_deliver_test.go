package config_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudConfigUpdateAggregateConfigDeliverResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfigUpdateAggregateConfigDeliver,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_config_update_aggregate_config_deliver.example", "id"),
					resource.TestCheckResourceAttr("tencentcloud_config_update_aggregate_config_deliver.example", "status", "1"),
					resource.TestCheckResourceAttr("tencentcloud_config_update_aggregate_config_deliver.example", "deliver_type", "COS"),
					resource.TestCheckResourceAttr("tencentcloud_config_update_aggregate_config_deliver.example", "deliver_content_type", "3"),
				),
			},
			{
				Config: testAccConfigUpdateAggregateConfigDeliverUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_config_update_aggregate_config_deliver.example", "status", "0"),
				),
			},
		},
	})
}

// NOTE: Replace account_group_id and target_arn with real values in your environment.
const testAccConfigUpdateAggregateConfigDeliver = `
resource "tencentcloud_config_update_aggregate_config_deliver" "example" {
  account_group_id     = "ca-xxxxxxxx"
  status               = 1
  deliver_name         = "tf-example-deliver"
  target_arn           = "qcs::cos:ap-guangzhou:uin/100000005287:prefix/1307050748/my-config-bucket"
  deliver_prefix       = "config"
  deliver_type         = "COS"
  deliver_uin          = 0
  deliver_content_type = 3
}
`

const testAccConfigUpdateAggregateConfigDeliverUpdate = `
resource "tencentcloud_config_update_aggregate_config_deliver" "example" {
  account_group_id = "ca-xxxxxxxx"
  status           = 0
}
`
