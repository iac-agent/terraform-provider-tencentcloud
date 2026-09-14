package ses_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/stretchr/testify/assert"
	ses "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ses/v20201002"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	svcses "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/ses"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// mockMetaSesReceiver implements tccommon.ProviderMeta for unit testing
type mockMetaSesReceiver struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaSesReceiver) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaSesReceiver{}

func newMockMetaSesReceiver() *mockMetaSesReceiver {
	return &mockMetaSesReceiver{client: &connectivity.TencentCloudClient{}}
}

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

// TestSesReceiverReadCountFromAPI verifies that the Read method populates the
// computed `count` attribute from the DescribeSesReceiverById API response.
func TestSesReceiverReadCountFromAPI(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	// Mock DescribeSesReceiverById to return a ReceiverData with Count set
	var receiverId uint64 = 1001
	var count uint64 = 50
	receiverData := &ses.ReceiverData{
		ReceiverId:    &receiverId,
		ReceiversName: strPtrSesReceiver("test-receiver"),
		Desc:          strPtrSesReceiver("test description"),
		Count:         &count,
	}

	patches.ApplyMethodFunc(&svcses.SesService{}, "DescribeSesReceiverById", func(_ context.Context, _ string) (*ses.ReceiverData, error) {
		return receiverData, nil
	})
	patches.ApplyMethodFunc(&svcses.SesService{}, "DescribeSesReceiverDetailById", func(_ context.Context, _ string) ([]*ses.ReceiverDetail, error) {
		return []*ses.ReceiverDetail{}, nil
	})

	// Create a resource data instance and call Read
	res := svcses.ResourceTencentCloudSesReceiver()
	d := res.TestResourceData()
	d.SetId("1001")

	meta := newMockMetaSesReceiver()
	err := res.Read(d, meta)
	assert.NoError(t, err)

	// Verify count attribute is set correctly
	assert.Equal(t, 50, d.Get("count"))
}

// TestSesReceiverReadCountNilFromAPI verifies that when the API returns nil Count,
// the Read method skips setting the count attribute without error.
func TestSesReceiverReadCountNilFromAPI(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	// Mock DescribeSesReceiverById to return a ReceiverData with Count nil
	var receiverId uint64 = 1001
	receiverData := &ses.ReceiverData{
		ReceiverId:    &receiverId,
		ReceiversName: strPtrSesReceiver("test-receiver"),
		Desc:          strPtrSesReceiver("test description"),
		Count:         nil,
	}

	patches.ApplyMethodFunc(&svcses.SesService{}, "DescribeSesReceiverById", func(_ context.Context, _ string) (*ses.ReceiverData, error) {
		return receiverData, nil
	})
	patches.ApplyMethodFunc(&svcses.SesService{}, "DescribeSesReceiverDetailById", func(_ context.Context, _ string) ([]*ses.ReceiverDetail, error) {
		return []*ses.ReceiverDetail{}, nil
	})

	// Create a resource data instance and call Read
	res := svcses.ResourceTencentCloudSesReceiver()
	d := res.TestResourceData()
	d.SetId("1001")

	meta := newMockMetaSesReceiver()
	err := res.Read(d, meta)
	assert.NoError(t, err)

	// Verify count attribute is 0 (default zero value, not set since Count was nil)
	assert.Equal(t, 0, d.Get("count"))
}

func strPtrSesReceiver(s string) *string {
	return &s
}
