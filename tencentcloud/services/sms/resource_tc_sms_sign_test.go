package sms_test

import (
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/sms"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// mockMetaForSmsSign implements tccommon.ProviderMeta
type mockMetaForSmsSign struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaForSmsSign) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaForSmsSign{}

func newMockMetaForSmsSign() *mockMetaForSmsSign {
	return &mockMetaForSmsSign{client: &connectivity.TencentCloudClient{}}
}

func ptrUint64SmsSign(v uint64) *uint64 {
	return &v
}

func ptrInt64SmsSign(v int64) *int64 {
	return &v
}

func ptrStrSmsSign(v string) *string {
	return &v
}

// TestSmsSign_Read_QualificationStatusCode_NonNil tests that resourceTencentCloudSmsSignRead
// correctly sets qualification_status_code when the API returns a non-nil value
// go test ./tencentcloud/services/sms/ -run "TestSmsSign_Read_QualificationStatusCode_NonNil" -v -count=1 -gcflags="all=-l"
func TestSmsSign_Read_QualificationStatusCode_NonNil(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	smsClient := &sms.Client{}
	patches.ApplyMethodReturn(newMockMetaForSmsSign().client, "UseSmsClient", smsClient)

	patches.ApplyMethodFunc(smsClient, "DescribeSmsSignList", func(request *sms.DescribeSmsSignListRequest) (*sms.DescribeSmsSignListResponse, error) {
		resp := sms.NewDescribeSmsSignListResponse()
		resp.Response = &sms.DescribeSmsSignListResponseParams{
			DescribeSignListStatusSet: []*sms.DescribeSignListStatus{
				{
					SignId:                  ptrUint64SmsSign(12345),
					SignName:                ptrStrSmsSign("test-sign"),
					International:           ptrUint64SmsSign(0),
					QualificationStatusCode: ptrInt64SmsSign(1),
				},
			},
			RequestId: ptrStrSmsSign("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForSmsSign()
	res := sms.ResourceTencentCloudSmsSign()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"sign_name":     "test-sign",
		"sign_type":     1,
		"document_type": 4,
		"international": 0,
		"sign_purpose":  0,
		"proof_image":   "dGhpcyBpcyBhIGV4YW1wbGU=",
	})
	d.SetId(helper.UInt64ToStr(12345) + tccommon.FILED_SP + "0")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, 1, d.Get("qualification_status_code"))
}

// TestSmsSign_Read_QualificationStatusCode_Nil tests that resourceTencentCloudSmsSignRead
// does not error when QualificationStatusCode is nil
// go test ./tencentcloud/services/sms/ -run "TestSmsSign_Read_QualificationStatusCode_Nil" -v -count=1 -gcflags="all=-l"
func TestSmsSign_Read_QualificationStatusCode_Nil(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	smsClient := &sms.Client{}
	patches.ApplyMethodReturn(newMockMetaForSmsSign().client, "UseSmsClient", smsClient)

	patches.ApplyMethodFunc(smsClient, "DescribeSmsSignList", func(request *sms.DescribeSmsSignListRequest) (*sms.DescribeSmsSignListResponse, error) {
		resp := sms.NewDescribeSmsSignListResponse()
		resp.Response = &sms.DescribeSmsSignListResponseParams{
			DescribeSignListStatusSet: []*sms.DescribeSignListStatus{
				{
					SignId:                  ptrUint64SmsSign(12345),
					SignName:                ptrStrSmsSign("test-sign"),
					International:           ptrUint64SmsSign(0),
					QualificationStatusCode: nil,
				},
			},
			RequestId: ptrStrSmsSign("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForSmsSign()
	res := sms.ResourceTencentCloudSmsSign()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"sign_name":     "test-sign",
		"sign_type":     1,
		"document_type": 4,
		"international": 0,
		"sign_purpose":  0,
		"proof_image":   "dGhpcyBpcyBhIGV4YW1wbGU=",
	})
	d.SetId(helper.UInt64ToStr(12345) + tccommon.FILED_SP + "0")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, 0, d.Get("qualification_status_code"))
}

func TestAccTencentCloudSmsSign_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_SMS) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSmsSign,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_sms_sign.sign", "id"),
					resource.TestCheckResourceAttr("tencentcloud_sms_sign.sign", "sign_name", "terraform"),
				),
			},
		},
	})
}

const testAccSmsSign = `

resource "tencentcloud_sms_sign" "sign" {
  sign_name     = "terraform"
  sign_type     = 1
  document_type = 4
  international = 0
  sign_purpose  = 0
  proof_image = "dGhpcyBpcyBhIGV4YW1wbGU="
}

`
