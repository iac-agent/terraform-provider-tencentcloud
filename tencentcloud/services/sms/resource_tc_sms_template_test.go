package sms_test

import (
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	svcsms "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/sms"
)

type mockMetaSmsTemplate struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaSmsTemplate) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaSmsTemplate{}

func newMockMetaSmsTemplate() *mockMetaSmsTemplate {
	return &mockMetaSmsTemplate{client: &connectivity.TencentCloudClient{}}
}

func ptrStrSmsTpl(s string) *string {
	return &s
}

func ptrUint64SmsTpl(u uint64) *uint64 {
	return &u
}

// TestSmsTemplate_Read_ReviewReplyNotNil verifies that when the API returns a
// non-nil ReviewReply, the Read method sets the review_reply attribute correctly.
func TestSmsTemplate_Read_ReviewReplyNotNil(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	smsClient := &sms.Client{}
	patches.ApplyMethodReturn(newMockMetaSmsTemplate().client, "UseSmsClient", smsClient)

	patches.ApplyMethodFunc(smsClient, "DescribeSmsTemplateList", func(request *sms.DescribeSmsTemplateListRequest) (*sms.DescribeSmsTemplateListResponse, error) {
		resp := sms.NewDescribeSmsTemplateListResponse()
		resp.Response = &sms.DescribeSmsTemplateListResponseParams{
			RequestId: ptrStrSmsTpl("fake-request-id-read"),
			DescribeTemplateStatusSet: []*sms.DescribeTemplateListStatus{
				{
					TemplateName:    ptrStrSmsTpl("test_template"),
					TemplateContent: ptrStrSmsTpl("test content"),
					International:   ptrUint64SmsTpl(0),
					ReviewReply:     ptrStrSmsTpl("template content is not compliant"),
				},
			},
		}
		return resp, nil
	})

	meta := newMockMetaSmsTemplate()
	res := svcsms.ResourceTencentCloudSmsTemplate()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"template_name":    "test_template",
		"template_content": "test content",
		"international":    0,
		"sms_type":         0,
		"remark":           "test",
	})
	d.SetId("12345" + tccommon.FILED_SP + "0")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "test_template", d.Get("template_name"))
	assert.Equal(t, "test content", d.Get("template_content"))
	assert.Equal(t, "template content is not compliant", d.Get("review_reply"))
}

// TestSmsTemplate_Read_ReviewReplyNil verifies that when the API returns nil for
// ReviewReply, the Read method does not set the review_reply attribute (nil check skips Set).
func TestSmsTemplate_Read_ReviewReplyNil(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	smsClient := &sms.Client{}
	patches.ApplyMethodReturn(newMockMetaSmsTemplate().client, "UseSmsClient", smsClient)

	patches.ApplyMethodFunc(smsClient, "DescribeSmsTemplateList", func(request *sms.DescribeSmsTemplateListRequest) (*sms.DescribeSmsTemplateListResponse, error) {
		resp := sms.NewDescribeSmsTemplateListResponse()
		resp.Response = &sms.DescribeSmsTemplateListResponseParams{
			RequestId: ptrStrSmsTpl("fake-request-id-read"),
			DescribeTemplateStatusSet: []*sms.DescribeTemplateListStatus{
				{
					TemplateName:    ptrStrSmsTpl("test_template"),
					TemplateContent: ptrStrSmsTpl("test content"),
					International:   ptrUint64SmsTpl(0),
					ReviewReply:     nil,
				},
			},
		}
		return resp, nil
	})

	meta := newMockMetaSmsTemplate()
	res := svcsms.ResourceTencentCloudSmsTemplate()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"template_name":    "test_template",
		"template_content": "test content",
		"international":    0,
		"sms_type":         0,
		"remark":           "test",
	})
	d.SetId("12345" + tccommon.FILED_SP + "0")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "test_template", d.Get("template_name"))
	// review_reply should remain empty string (default for TypeString) since nil check skips Set
	assert.Equal(t, "", d.Get("review_reply"))
}

// TestSmsTemplate_Read_ReviewReplyEmptyString verifies that when the API returns
// an empty string for ReviewReply, the Read method sets review_reply to empty string.
func TestSmsTemplate_Read_ReviewReplyEmptyString(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	smsClient := &sms.Client{}
	patches.ApplyMethodReturn(newMockMetaSmsTemplate().client, "UseSmsClient", smsClient)

	patches.ApplyMethodFunc(smsClient, "DescribeSmsTemplateList", func(request *sms.DescribeSmsTemplateListRequest) (*sms.DescribeSmsTemplateListResponse, error) {
		resp := sms.NewDescribeSmsTemplateListResponse()
		resp.Response = &sms.DescribeSmsTemplateListResponseParams{
			RequestId: ptrStrSmsTpl("fake-request-id-read"),
			DescribeTemplateStatusSet: []*sms.DescribeTemplateListStatus{
				{
					TemplateName:    ptrStrSmsTpl("test_template"),
					TemplateContent: ptrStrSmsTpl("test content"),
					International:   ptrUint64SmsTpl(0),
					ReviewReply:     ptrStrSmsTpl(""),
				},
			},
		}
		return resp, nil
	})

	meta := newMockMetaSmsTemplate()
	res := svcsms.ResourceTencentCloudSmsTemplate()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"template_name":    "test_template",
		"template_content": "test content",
		"international":    0,
		"sms_type":         0,
		"remark":           "test",
	})
	d.SetId("12345" + tccommon.FILED_SP + "0")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Get("review_reply"))
}
