package tag_test

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
)

// Register the tccommon.Provider function for testing
func init() {
	tccommon.Provider = func() *schema.Provider {
		return tccommon.Provider()
	}
}

// testSuppressJSONOrderDiff compares two JSON strings ignoring key order.
// It returns true if the parsed JSON objects are semantically equal.
func testSuppressJSONOrderDiff(old, new string, d *schema.ResourceData) bool {
	// Handle empty strings
	if old == "" && new == "" {
		return true
	}
	if old == "" || new == "" {
		return false
	}

	// Parse both JSON strings
	var oldJSON, newJSON interface{}
	if err := json.Unmarshal([]byte(old), &oldJSON); err != nil {
		return old == new // Fallback to string comparison
	}
	if err := json.Unmarshal([]byte(new), &newJSON); err != nil {
		return old == new // Fallback to string comparison
	}

	// Compare parsed JSON objects (ignoring key order)
	return reflect.DeepEqual(oldJSON, newJSON)
}

func TestAccTencentCloudTagAttachment_autoRenewFlagUpdate(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"tencentcloud": func() (*schema.Provider, error) {
				return tccommon.Provider(), nil
			},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccTagAttachmentConfig_autoRenewFlag(0),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_tag_attachment.test", "auto_renew_flag", "0"),
				),
			},
			{
				Config: testAccTagAttachmentConfig_autoRenewFlag(1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_tag_attachment.test", "auto_renew_flag", "1"),
				),
			},
		},
	})
}

func testAccTagAttachmentConfig_autoRenewFlag(flag int) string {
	return fmt.Sprintf(`
resource "tencentcloud_tag_attachment" "tag_attachment" {
  tag_key = "test_terraform_tagAttachment_key"
  auto_renew_flag = %d
}
`, flag)
}

func testAccTagAttachmentConfig_tagValue(tagValue string) string {
	return fmt.Sprintf(`
resource "tencentcloud_tag_attachment" "tag_attachment" {
  tag_key = "test_terraform_tagAttachment_key"
  tag_value = "%s"
}
`, tagValue)
}
