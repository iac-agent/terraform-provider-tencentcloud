package tag_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	tag "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tag/v20180813"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	svctag "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// mockMetaTagAttachment implements tccommon.ProviderMeta
type mockMetaTagAttachment struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaTagAttachment) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaTagAttachment{}

func ptrStringTagAttachment(s string) *string {
	return &s
}

// go test -i; go test -test.run TestAccTencentCloudTagAttachmentResource_basic -v
func TestAccTencentCloudTagAttachmentResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckTagAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTagResourceTag,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTagAttachmentExists("tencentcloud_tag_attachment.tag_attachment"),
					resource.TestCheckResourceAttr("tencentcloud_tag_attachment.tag_attachment", "tag_key", "test_terraform_tagAttachment_key"),
					resource.TestCheckResourceAttr("tencentcloud_tag_attachment.tag_attachment", "tag_value", "Terraform_tagAttachment_value"),
					resource.TestCheckResourceAttrSet("tencentcloud_tag_attachment.tag_attachment", "resource")),
			},
			{
				ResourceName:      "tencentcloud_tag_attachment.tag_attachment",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
func testAccCheckTagAttachmentDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_tag_attachment" {
			continue
		}
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service := svctag.NewTagService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())

		tags, err := service.DescribeTagTagAttachmentById(ctx, rs.Primary.Attributes["tag_key"],
			rs.Primary.Attributes["tag_value"], rs.Primary.Attributes["resource"])
		if err != nil {
			return err
		}
		if tags == nil {
			return nil
		}
		return fmt.Errorf("delete tagAttachment key %s fail, still on server", rs.Primary.Attributes["tag_key"])
	}
	return nil
}

func testAccCheckTagAttachmentExists(r string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}

		service := svctag.NewTagService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
		res, err := service.DescribeTagTagAttachmentById(ctx, rs.Primary.Attributes["tag_key"],
			rs.Primary.Attributes["tag_value"], rs.Primary.Attributes["resource"])
		if err != nil {
			return err
		}
		if res != nil && res.Resource != nil && res.Tags != nil {
			return nil
		}

		return fmt.Errorf("tagAttachment %s not found on server", rs.Primary.Attributes["tag_key"])
	}
}

const testAccTagResourceTag = tcacctest.DefaultCvmModificationVariable + `
data "tencentcloud_user_info" "info" {}

locals {
  uin = data.tencentcloud_user_info.info.uin
}

resource "tencentcloud_tag_attachment" "tag_attachment" {
  tag_key = "test_terraform_tagAttachment_key"
  tag_value = "Terraform_tagAttachment_value"
  resource = "qcs::cvm:ap-guangzhou:uin/${local.uin}:instance/${var.cvm_id}"
}

`

const tagAttachmentTestResourceSixSegment = "qcs::cvm:ap-guangzhou:uin/123456789:instance/ins-test1234"

// TestTagAttachmentUpdate_TagValueChanged tests that when tag_value changes,
// ModifyTags is called with the full six-segment resource string and the
// ReplaceTags contains the correct key/value, and the ID is updated afterwards.
//
// go test ./tencentcloud/services/tag/ -run "TestTagAttachmentUpdate_TagValueChanged" -v -count=1 -gcflags="all=-l"
func TestTagAttachmentUpdate_TagValueChanged(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	var (
		capturedResource    string
		capturedReplaceTags map[string]string
		capturedDeleteKeys  []string
	)
	patches.ApplyMethodFunc(&svctag.TagService{}, "ModifyTags", func(_ context.Context, resourceName string, replaceTags map[string]string, deleteKeys []string) error {
		capturedResource = resourceName
		capturedReplaceTags = replaceTags
		capturedDeleteKeys = deleteKeys
		return nil
	})

	resourcePtr := tagAttachmentTestResourceSixSegment
	patches.ApplyMethodFunc(&svctag.TagService{}, "DescribeTagTagAttachmentById", func(_ context.Context, tagKey string, tagValue string, resourceId string) (*tag.ResourceTagMapping, error) {
		return &tag.ResourceTagMapping{
			Resource: &resourcePtr,
			Tags: []*tag.Tag{
				{
					TagKey:   ptrStringTagAttachment(tagKey),
					TagValue: ptrStringTagAttachment(tagValue),
				},
			},
		}, nil
	})

	res := svctag.ResourceTencentCloudTagAttachment()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"tag_key":   "运营产品",
		"tag_value": "B",
		"resource":  tagAttachmentTestResourceSixSegment,
	})
	d.SetId("运营产品" + tccommon.FILED_SP + "A" + tccommon.FILED_SP + tagAttachmentTestResourceSixSegment)

	// Simulate that only tag_value changed.
	patches.ApplyMethodFunc(d, "HasChange", func(key string) bool {
		return key == "tag_value"
	})

	meta := &mockMetaTagAttachment{}
	err := res.Update(d, meta)
	assert.NoError(t, err)

	assert.Equal(t, tagAttachmentTestResourceSixSegment, capturedResource)
	assert.Equal(t, map[string]string{"运营产品": "B"}, capturedReplaceTags)
	assert.Nil(t, capturedDeleteKeys)

	assert.Equal(t, "运营产品"+tccommon.FILED_SP+"B"+tccommon.FILED_SP+tagAttachmentTestResourceSixSegment, d.Id())
}

// TestTagAttachmentUpdate_TagValueUnchanged tests that when tag_value is unchanged,
// ModifyTags is NOT called.
//
// go test ./tencentcloud/services/tag/ -run "TestTagAttachmentUpdate_TagValueUnchanged" -v -count=1 -gcflags="all=-l"
func TestTagAttachmentUpdate_TagValueUnchanged(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	var modifyCalled bool
	patches.ApplyMethodFunc(&svctag.TagService{}, "ModifyTags", func(_ context.Context, resourceName string, replaceTags map[string]string, deleteKeys []string) error {
		modifyCalled = true
		return nil
	})

	resourcePtr := tagAttachmentTestResourceSixSegment
	patches.ApplyMethodFunc(&svctag.TagService{}, "DescribeTagTagAttachmentById", func(_ context.Context, tagKey string, tagValue string, resourceId string) (*tag.ResourceTagMapping, error) {
		return &tag.ResourceTagMapping{
			Resource: &resourcePtr,
			Tags: []*tag.Tag{
				{
					TagKey:   ptrStringTagAttachment(tagKey),
					TagValue: ptrStringTagAttachment(tagValue),
				},
			},
		}, nil
	})

	res := svctag.ResourceTencentCloudTagAttachment()
	originalId := "运营产品" + tccommon.FILED_SP + "A" + tccommon.FILED_SP + tagAttachmentTestResourceSixSegment
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"tag_key":   "运营产品",
		"tag_value": "A",
		"resource":  tagAttachmentTestResourceSixSegment,
	})
	d.SetId(originalId)

	// No field changed.
	patches.ApplyMethodFunc(d, "HasChange", func(key string) bool {
		return false
	})

	meta := &mockMetaTagAttachment{}
	err := res.Update(d, meta)
	assert.NoError(t, err)

	assert.False(t, modifyCalled)
	assert.Equal(t, originalId, d.Id())
}

// TestTagAttachmentUpdate_BrokenId tests that when the ID is not a three-segment
// format, Update returns an error without calling any cloud API.
//
// go test ./tencentcloud/services/tag/ -run "TestTagAttachmentUpdate_BrokenId" -v -count=1 -gcflags="all=-l"
func TestTagAttachmentUpdate_BrokenId(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	var modifyCalled bool
	patches.ApplyMethodFunc(&svctag.TagService{}, "ModifyTags", func(_ context.Context, resourceName string, replaceTags map[string]string, deleteKeys []string) error {
		modifyCalled = true
		return nil
	})

	res := svctag.ResourceTencentCloudTagAttachment()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"tag_key":   "运营产品",
		"tag_value": "B",
		"resource":  tagAttachmentTestResourceSixSegment,
	})
	d.SetId("broken-id-without-separator")

	patches.ApplyMethodFunc(d, "HasChange", func(key string) bool {
		return key == "tag_value"
	})

	meta := &mockMetaTagAttachment{}
	err := res.Update(d, meta)
	assert.Error(t, err)

	assert.False(t, modifyCalled)
}

// TestTagAttachmentUpdate_ModifyTagsError tests that when ModifyTags returns an
// error, Update propagates the error upward.
//
// go test ./tencentcloud/services/tag/ -run "TestTagAttachmentUpdate_ModifyTagsError" -v -count=1 -gcflags="all=-l"
func TestTagAttachmentUpdate_ModifyTagsError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	expectedErr := fmt.Errorf("mock ModifyTags error")
	patches.ApplyMethodFunc(&svctag.TagService{}, "ModifyTags", func(_ context.Context, resourceName string, replaceTags map[string]string, deleteKeys []string) error {
		return expectedErr
	})

	res := svctag.ResourceTencentCloudTagAttachment()
	originalId := "运营产品" + tccommon.FILED_SP + "A" + tccommon.FILED_SP + tagAttachmentTestResourceSixSegment
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"tag_key":   "运营产品",
		"tag_value": "B",
		"resource":  tagAttachmentTestResourceSixSegment,
	})
	d.SetId(originalId)

	patches.ApplyMethodFunc(d, "HasChange", func(key string) bool {
		return key == "tag_value"
	})

	meta := &mockMetaTagAttachment{}
	err := res.Update(d, meta)
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)

	// ID should not have been updated since ModifyTags failed.
	assert.Equal(t, originalId, d.Id())
}
