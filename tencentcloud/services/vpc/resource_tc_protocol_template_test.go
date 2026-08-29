package vpc_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	tag "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tag/v20180813"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	svctag "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tag"
	svcvpc "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/vpc"
)

// ---- Acceptance tests (existing) ----

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

// ---- gomonkey-based unit tests for tags ----
// Run with: go test ./tencentcloud/services/vpc/ -run "TestProtocolTemplateTags" -v -count=1 -gcflags="all=-l"

type mockMetaProtocolTemplate struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaProtocolTemplate) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaProtocolTemplate{}

func newMockMetaProtocolTemplate() *mockMetaProtocolTemplate {
	return &mockMetaProtocolTemplate{client: &connectivity.TencentCloudClient{Region: "ap-guangzhou"}}
}

func ptrStringPT(s string) *string {
	return &s
}

// TestProtocolTemplateTags_Create verifies tags are passed to CreateServiceTemplate.
func TestProtocolTemplateTags_Create(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	vpcClient := &vpc.Client{}
	patches.ApplyMethodReturn(newMockMetaProtocolTemplate().client, "UseVpcClient", vpcClient)

	templateId := "ppm-test123"
	var capturedRequest *vpc.CreateServiceTemplateRequest
	patches.ApplyMethodFunc(vpcClient, "CreateServiceTemplate", func(request *vpc.CreateServiceTemplateRequest) (*vpc.CreateServiceTemplateResponse, error) {
		capturedRequest = request
		resp := vpc.NewCreateServiceTemplateResponse()
		resp.Response = &vpc.CreateServiceTemplateResponseParams{
			ServiceTemplate: &vpc.ServiceTemplate{
				ServiceTemplateId:   &templateId,
				ServiceTemplateName: ptrStringPT("test-template"),
			},
		}
		return resp, nil
	})

	// Mock DescribeServiceTemplates called by the read after create.
	patches.ApplyMethodFunc(vpcClient, "DescribeServiceTemplates", func(request *vpc.DescribeServiceTemplatesRequest) (*vpc.DescribeServiceTemplatesResponse, error) {
		resp := vpc.NewDescribeServiceTemplatesResponse()
		resp.Response = &vpc.DescribeServiceTemplatesResponseParams{
			ServiceTemplateSet: []*vpc.ServiceTemplate{
				{
					ServiceTemplateId:   &templateId,
					ServiceTemplateName: ptrStringPT("test-template"),
					ServiceSet:          []*string{ptrStringPT("tcp:80")},
					TagSet: []*vpc.Tag{
						{Key: ptrStringPT("env"), Value: ptrStringPT("prod")},
						{Key: ptrStringPT("team"), Value: ptrStringPT("infra")},
					},
				},
			},
		}
		return resp, nil
	})

	meta := newMockMetaProtocolTemplate()
	res := svcvpc.ResourceTencentCloudProtocolTemplate()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"name":      "test-template",
		"protocols": []interface{}{"tcp:80"},
		"tags": map[string]interface{}{
			"env":  "prod",
			"team": "infra",
		},
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, templateId, d.Id())

	// Verify tags were passed to the create request.
	assert.NotNil(t, capturedRequest)
	assert.Len(t, capturedRequest.Tags, 2)

	tagMap := make(map[string]string)
	for _, tag := range capturedRequest.Tags {
		if tag.Key != nil && tag.Value != nil {
			tagMap[*tag.Key] = *tag.Value
		}
	}
	assert.Equal(t, "prod", tagMap["env"])
	assert.Equal(t, "infra", tagMap["team"])

	// Verify tags were read back into state.
	tags := d.Get("tags").(map[string]interface{})
	assert.Equal(t, "prod", tags["env"])
	assert.Equal(t, "infra", tags["team"])
}

// TestProtocolTemplateTags_ReadWithTags verifies tags are read from TagSet in DescribeServiceTemplates response.
func TestProtocolTemplateTags_ReadWithTags(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	vpcClient := &vpc.Client{}
	patches.ApplyMethodReturn(newMockMetaProtocolTemplate().client, "UseVpcClient", vpcClient)

	templateId := "ppm-readtest"
	patches.ApplyMethodFunc(vpcClient, "DescribeServiceTemplates", func(request *vpc.DescribeServiceTemplatesRequest) (*vpc.DescribeServiceTemplatesResponse, error) {
		resp := vpc.NewDescribeServiceTemplatesResponse()
		resp.Response = &vpc.DescribeServiceTemplatesResponseParams{
			ServiceTemplateSet: []*vpc.ServiceTemplate{
				{
					ServiceTemplateId:   &templateId,
					ServiceTemplateName: ptrStringPT("read-template"),
					ServiceSet:          []*string{ptrStringPT("tcp:80")},
					TagSet: []*vpc.Tag{
						{Key: ptrStringPT("env"), Value: ptrStringPT("staging")},
						{Key: ptrStringPT("owner"), Value: ptrStringPT("platform")},
					},
				},
			},
		}
		return resp, nil
	})

	meta := newMockMetaProtocolTemplate()
	res := svcvpc.ResourceTencentCloudProtocolTemplate()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"name":      "read-template",
		"protocols": []interface{}{"tcp:80"},
	})
	d.SetId(templateId)

	err := res.Read(d, meta)
	assert.NoError(t, err)

	// Verify tags were read from TagSet.
	tags := d.Get("tags").(map[string]interface{})
	assert.Equal(t, "staging", tags["env"])
	assert.Equal(t, "platform", tags["owner"])
}

// TestProtocolTemplateTags_ReadWithoutTags verifies no tags set when TagSet is nil.
func TestProtocolTemplateTags_ReadWithoutTags(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	vpcClient := &vpc.Client{}
	patches.ApplyMethodReturn(newMockMetaProtocolTemplate().client, "UseVpcClient", vpcClient)

	templateId := "ppm-notags"
	patches.ApplyMethodFunc(vpcClient, "DescribeServiceTemplates", func(request *vpc.DescribeServiceTemplatesRequest) (*vpc.DescribeServiceTemplatesResponse, error) {
		resp := vpc.NewDescribeServiceTemplatesResponse()
		resp.Response = &vpc.DescribeServiceTemplatesResponseParams{
			ServiceTemplateSet: []*vpc.ServiceTemplate{
				{
					ServiceTemplateId:   &templateId,
					ServiceTemplateName: ptrStringPT("notags-template"),
					ServiceSet:          []*string{ptrStringPT("tcp:80")},
				},
			},
		}
		return resp, nil
	})

	meta := newMockMetaProtocolTemplate()
	res := svcvpc.ResourceTencentCloudProtocolTemplate()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"name":      "notags-template",
		"protocols": []interface{}{"tcp:80"},
	})
	d.SetId(templateId)

	err := res.Read(d, meta)
	assert.NoError(t, err)

	// Verify tags are empty when TagSet is nil.
	tags := d.Get("tags").(map[string]interface{})
	assert.Empty(t, tags)
}

// TestProtocolTemplateTags_Update verifies tag update via tag service when tags change.
func TestProtocolTemplateTags_Update(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	vpcClient := &vpc.Client{}
	patches.ApplyMethodReturn(newMockMetaProtocolTemplate().client, "UseVpcClient", vpcClient)

	tagClient := &tag.Client{}
	patches.ApplyMethodReturn(newMockMetaProtocolTemplate().client, "UseTagClient", tagClient)

	templateId := "ppm-updatetest"

	// Mock ModifyServiceTemplateAttribute (for name/protocols update, not needed for tags-only update)
	patches.ApplyMethodFunc(vpcClient, "ModifyServiceTemplateAttribute", func(request *vpc.ModifyServiceTemplateAttributeRequest) (*vpc.ModifyServiceTemplateAttributeResponse, error) {
		resp := vpc.NewModifyServiceTemplateAttributeResponse()
		resp.Response = &vpc.ModifyServiceTemplateAttributeResponseParams{}
		return resp, nil
	})

	// Mock DescribeServiceTemplates for read after update
	patches.ApplyMethodFunc(vpcClient, "DescribeServiceTemplates", func(request *vpc.DescribeServiceTemplatesRequest) (*vpc.DescribeServiceTemplatesResponse, error) {
		resp := vpc.NewDescribeServiceTemplatesResponse()
		resp.Response = &vpc.DescribeServiceTemplatesResponseParams{
			ServiceTemplateSet: []*vpc.ServiceTemplate{
				{
					ServiceTemplateId:   &templateId,
					ServiceTemplateName: ptrStringPT("update-template"),
					ServiceSet:          []*string{ptrStringPT("tcp:80")},
					TagSet: []*vpc.Tag{
						{Key: ptrStringPT("env"), Value: ptrStringPT("prod")},
						{Key: ptrStringPT("version"), Value: ptrStringPT("v2")},
					},
				},
			},
		}
		return resp, nil
	})

	// Mock tag service ModifyResourceTags
	var capturedModifyRequest *tag.ModifyResourceTagsRequest
	patches.ApplyMethodFunc(tagClient, "ModifyResourceTags", func(request *tag.ModifyResourceTagsRequest) (*tag.ModifyResourceTagsResponse, error) {
		capturedModifyRequest = request
		resp := tag.NewModifyResourceTagsResponse()
		resp.Response = &tag.ModifyResourceTagsResponseParams{}
		return resp, nil
	})

	meta := newMockMetaProtocolTemplate()
	res := svcvpc.ResourceTencentCloudProtocolTemplate()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"name":      "update-template",
		"protocols": []interface{}{"tcp:80"},
	})
	d.SetId(templateId)

	// Set old state: env=staging, team=infra
	_ = d.Set("tags", map[string]interface{}{
		"env":  "staging",
		"team": "infra",
	})

	// Now change to new state: env=prod, version=v2
	// This means: replace env (staging->prod), add version (v2), delete team
	_ = d.Set("tags", map[string]interface{}{
		"env":     "prod",
		"version": "v2",
	})

	err := res.Update(d, meta)
	assert.NoError(t, err)

	// Verify tag updates via tag service.
	// Old: env=staging, team=infra
	// New: env=prod, version=v2
	// Replace: env=prod (changed), version=v2 (new)
	// Delete: team (removed)
	assert.NotNil(t, capturedModifyRequest)
	assert.NotNil(t, capturedModifyRequest.ReplaceTags)

	replaceTagsMap := make(map[string]string)
	for _, t := range capturedModifyRequest.ReplaceTags {
		if t.TagKey != nil && t.TagValue != nil {
			replaceTagsMap[*t.TagKey] = *t.TagValue
		}
	}
	assert.Equal(t, "prod", replaceTagsMap["env"])
	assert.Equal(t, "v2", replaceTagsMap["version"])

	assert.NotNil(t, capturedModifyRequest.DeleteTags)
	deleteKeys := make([]string, 0)
	for _, t := range capturedModifyRequest.DeleteTags {
		if t.TagKey != nil {
			deleteKeys = append(deleteKeys, *t.TagKey)
		}
	}
	assert.Contains(t, deleteKeys, "team")
}
