package teo_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudTeoInferenceDomainResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTeoInferenceDomain,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_teo_inference_domain.example", "id"),
					resource.TestCheckResourceAttr("tencentcloud_teo_inference_domain.example", "zone_id", "zone-3v0yhgt4m1z9"),
					resource.TestCheckResourceAttr("tencentcloud_teo_inference_domain.example", "service_id", "inf-vs06pb1cdz61"),
					resource.TestCheckResourceAttr("tencentcloud_teo_inference_domain.example", "domain", "test-ai.eo.cn"),
					resource.TestCheckResourceAttr("tencentcloud_teo_inference_domain.example", "auth_switch", "On"),
				),
			},
		},
	})
}

const testAccTeoInferenceDomain = `
resource "tencentcloud_teo_inference_domain" "example" {
  zone_id     = "zone-3v0yhgt4m1z9"
  service_id  = "inf-vs06pb1cdz61"
  domain      = "test-ai.eo.cn"
  auth_switch = "On"
}
`
