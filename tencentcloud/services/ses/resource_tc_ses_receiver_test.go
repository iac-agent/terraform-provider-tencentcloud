package ses_test

import (
	"context"
	"fmt"
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svcses "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/ses"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"
	ses "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ses/v20201002"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	localses "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/ses"
)

// go test -test.run TestAccTencentCloudSesReceiverResource_basic -v
func TestAccTencentCloudSesReceiverResource_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccStepSetRegion(t, "ap-hongkong")
			tcacctest.AccPreCheckBusiness(t, tcacctest.ACCOUNT_TYPE_SES)
		},
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckSesReceiverDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSesReceiver,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSesReceiverExists("tencentcloud_ses_receiver.receiver"),
					resource.TestCheckResourceAttrSet("tencentcloud_ses_receiver.receiver", "id"),
					resource.TestCheckResourceAttr("tencentcloud_ses_receiver.receiver", "receivers_name", "terraform_test"),
					resource.TestCheckResourceAttr("tencentcloud_ses_receiver.receiver", "desc", "description"),
					resource.TestCheckResourceAttr("tencentcloud_ses_receiver.receiver", "data.#", "2"),
					resource.TestCheckResourceAttr("tencentcloud_ses_receiver.receiver", "data.0.email", "abc@abc.com"),
					resource.TestCheckResourceAttr("tencentcloud_ses_receiver.receiver", "data.0.template_data", "{\"name\":\"xxx\",\"age\":\"xx\"}"),
					resource.TestCheckResourceAttrSet("tencentcloud_ses_receiver.receiver", "count"),
				),
			},
			{
				ResourceName:      "tencentcloud_ses_receiver.receiver",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckSesReceiverDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	service := svcses.NewSesService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_ses_receiver" {
			continue
		}

		res, err := service.DescribeSesReceiverById(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}

		if res != nil {
			return fmt.Errorf("ses receiver %s still exists", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckSesReceiverExists(r string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}

		service := svcses.NewSesService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
		res, err := service.DescribeSesReceiverById(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}

		if res == nil {
			return fmt.Errorf("ses receiver %s is not found", rs.Primary.ID)
		}

		return nil
	}
}

const testAccSesReceiver = `

resource "tencentcloud_ses_receiver" "receiver" {
  receivers_name = "terraform_test"
  desc = "description"

  data {
    email = "abc@abc.com"
    template_data = "{\"name\":\"xxx\",\"age\":\"xx\"}"
  }

  data {
    email = "abcd@abcd.com"
    template_data = "{\"name\":\"xxx\",\"age\":\"xx\"}"
  }
}

`

// --- gomonkey mock unit tests for count attribute ---

type mockMetaForSesReceiver struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaForSesReceiver) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaForSesReceiver{}

func newMockMetaForSesReceiver() *mockMetaForSesReceiver {
	return &mockMetaForSesReceiver{client: &connectivity.TencentCloudClient{}}
}

func ptrUint64SR(v uint64) *uint64 { return &v }
func ptrStringSR(s string) *string { return &s }

// TestSesReceiver_Read_SetsCount verifies that the Read function sets count from API response
// go test ./tencentcloud/services/ses/ -run "TestSesReceiver_Read_SetsCount" -v -count=1 -gcflags="all=-l"
func TestSesReceiver_Read_SetsCount(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	receiverData := &ses.ReceiverData{
		ReceiverId:    ptrUint64SR(1001),
		ReceiversName: ptrStringSR("test-receiver"),
		Desc:          ptrStringSR("test-desc"),
		Count:         ptrUint64SR(42),
	}

	patches.ApplyMethodFunc(&localses.SesService{}, "DescribeSesReceiverById", func(_ context.Context, _ string) (*ses.ReceiverData, error) {
		return receiverData, nil
	})

	patches.ApplyMethodFunc(&localses.SesService{}, "DescribeSesReceiverDetailById", func(_ context.Context, _ string) ([]*ses.ReceiverDetail, error) {
		return nil, nil
	})

	meta := newMockMetaForSesReceiver()
	res := localses.ResourceTencentCloudSesReceiver()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"receivers_name": "test-receiver",
		"data":           []interface{}{},
	})
	d.SetId("1001")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, 42, d.Get("count").(int))
	assert.Equal(t, "test-receiver", d.Get("receivers_name").(string))
	assert.Equal(t, "test-desc", d.Get("desc").(string))
}

// TestSesReceiver_Read_CountNil verifies that when Count is nil in API response, count is not set
// go test ./tencentcloud/services/ses/ -run "TestSesReceiver_Read_CountNil" -v -count=1 -gcflags="all=-l"
func TestSesReceiver_Read_CountNil(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	receiverData := &ses.ReceiverData{
		ReceiverId:    ptrUint64SR(1002),
		ReceiversName: ptrStringSR("test-receiver-nil"),
		Desc:          ptrStringSR("test-desc-nil"),
		Count:         nil,
	}

	patches.ApplyMethodFunc(&localses.SesService{}, "DescribeSesReceiverById", func(_ context.Context, _ string) (*ses.ReceiverData, error) {
		return receiverData, nil
	})

	patches.ApplyMethodFunc(&localses.SesService{}, "DescribeSesReceiverDetailById", func(_ context.Context, _ string) ([]*ses.ReceiverDetail, error) {
		return nil, nil
	})

	meta := newMockMetaForSesReceiver()
	res := localses.ResourceTencentCloudSesReceiver()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"receivers_name": "test-receiver-nil",
		"data":           []interface{}{},
	})
	d.SetId("1002")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, 0, d.Get("count").(int))
}

// TestSesReceiver_Schema_CountComputed verifies that count is a computed attribute in the schema
// go test ./tencentcloud/services/ses/ -run "TestSesReceiver_Schema_CountComputed" -v -count=1 -gcflags="all=-l"
func TestSesReceiver_Schema_CountComputed(t *testing.T) {
	res := localses.ResourceTencentCloudSesReceiver()

	assert.Contains(t, res.Schema, "count")
	countSchema := res.Schema["count"]
	assert.Equal(t, schema.TypeInt, countSchema.Type)
	assert.True(t, countSchema.Computed)
	assert.False(t, countSchema.Optional)
	assert.False(t, countSchema.Required)
}
