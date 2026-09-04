package vpc_test

import (
	"context"
	"fmt"
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svcvpc "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/vpc"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccTencentCloudProtocolTemplate_basic_and_update(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckProtocolTemplateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccProtocolTemplate_basic,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_protocol_template.template", "name", "test"),
					resource.TestCheckResourceAttr("tencentcloud_protocol_template.template", "protocols.#", "1"),
				),
			},
			{
				ResourceName:      "tencentcloud_protocol_template.template",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccProtocolTemplate_basic_update_remark,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckProtocolTemplateExists("tencentcloud_protocol_template.template"),
					resource.TestCheckResourceAttr("tencentcloud_protocol_template.template", "name", "test_update"),
					resource.TestCheckResourceAttr("tencentcloud_protocol_template.template", "protocols.#", "2"),
				),
			},
		},
	})
}

func TestAccTencentCloudProtocolTemplate_tags(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckProtocolTemplateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccProtocolTemplate_tags,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckProtocolTemplateExists("tencentcloud_protocol_template.template"),
					resource.TestCheckResourceAttr("tencentcloud_protocol_template.template", "name", "tf-tags-test"),
					resource.TestCheckResourceAttr("tencentcloud_protocol_template.template", "tags.env", "prod"),
					resource.TestCheckResourceAttr("tencentcloud_protocol_template.template", "tags.team", "infra"),
				),
			},
			{
				ResourceName:      "tencentcloud_protocol_template.template",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccProtocolTemplate_tags_update,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckProtocolTemplateExists("tencentcloud_protocol_template.template"),
					resource.TestCheckResourceAttr("tencentcloud_protocol_template.template", "tags.env", "dev"),
					resource.TestCheckResourceAttr("tencentcloud_protocol_template.template", "tags.owner", "tf"),
				),
			},
		},
	})
}

func testAccCheckProtocolTemplateDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	vpcService := svcvpc.NewVpcService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_protocol_template" {
			continue
		}

		_, has, err := vpcService.DescribeServiceTemplateById(ctx, rs.Primary.ID)
		if has {
			return fmt.Errorf("protocol template still exists")
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func testAccCheckProtocolTemplateExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Service template %s is not found", n)
		}

		vpcService := svcvpc.NewVpcService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
		_, has, err := vpcService.DescribeServiceTemplateById(ctx, rs.Primary.ID)
		if !has {
			return fmt.Errorf("Service template %s is not found", rs.Primary.ID)
		}
		if err != nil {
			return err
		}

		return nil
	}
}

const testAccProtocolTemplate_basic = `
resource "tencentcloud_protocol_template" "template" {
  name = "test"
  protocols = ["tcp:80"]
}`

const testAccProtocolTemplate_basic_update_remark = `
resource "tencentcloud_protocol_template" "template" {
  name = "test_update"
  protocols = ["udp:all", "tcp:80,90"]
}`

const testAccProtocolTemplate_tags = `
resource "tencentcloud_protocol_template" "template" {
  name = "tf-tags-test"
  protocols = ["tcp:80"]
  tags = {
    "env" = "prod"
    "team" = "infra"
  }
}`

const testAccProtocolTemplate_tags_update = `
resource "tencentcloud_protocol_template" "template" {
  name = "tf-tags-test"
  protocols = ["tcp:80"]
  tags = {
    "env" = "dev"
    "owner" = "tf"
  }
}`
