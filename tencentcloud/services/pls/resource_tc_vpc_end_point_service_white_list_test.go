package pls_test

import (
	"context"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/pls"
	svcvpc "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/vpc"

	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
)

// mockMeta implements tccommon.ProviderMeta
type mockMeta struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMeta) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMeta{}

func newMockMeta() *mockMeta {
	return &mockMeta{client: &connectivity.TencentCloudClient{}}
}

func ptrString(s string) *string {
	return &s
}

func ptrUint64(v uint64) *uint64 {
	return &v
}

func TestAccTencentCloudVpcEndPointServiceWhiteListResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcEndPointServiceWhiteList,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_vpc_end_point_service_white_list.end_point_service_white_list", "id")),
			},
			{
				ResourceName:      "tencentcloud_vpc_end_point_service_white_list.end_point_service_white_list",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccVpcEndPointServiceWhiteList = `

resource "tencentcloud_vpc_end_point_service_white_list" "end_point_service_white_list" {
  user_uin = "100020512675"
  end_point_service_id = "vpcsvc-98jddhcz"
  description = "terraform for test"
}

`

// go test ./tencentcloud/services/pls/ -run "TestVpcEndPointServiceWhiteListRead_TotalCountSet" -v -count=1 -gcflags="all=-l"
// TestVpcEndPointServiceWhiteListRead_TotalCountSet verifies total_count is set when the mocked API response has a non-nil TotalCount.
func TestVpcEndPointServiceWhiteListRead_TotalCountSet(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	expectedUserUin := "100020512675"
	expectedEndPointServiceId := "vpcsvc-98jddhcz"
	expectedTotalCount := uint64(3)

	whiteListRecord := &vpc.VpcEndPointServiceUser{
		UserUin:           ptrString(expectedUserUin),
		EndPointServiceId: ptrString(expectedEndPointServiceId),
		Description:       ptrString("terraform for test"),
		Owner:             ptrUint64(1234567890),
		CreateTime:        ptrString("2024-01-01 00:00:00"),
	}

	patches.ApplyMethodFunc(&svcvpc.VpcService{}, "DescribeVpcEndPointServiceWhiteListById",
		func(_ context.Context, userUin string, endPointServiceId string) (*vpc.VpcEndPointServiceUser, *uint64, error) {
			assert.Equal(t, expectedUserUin, userUin)
			assert.Equal(t, expectedEndPointServiceId, endPointServiceId)
			return whiteListRecord, &expectedTotalCount, nil
		})

	res := pls.ResourceTencentCloudVpcEndPointServiceWhiteList()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"user_uin":             expectedUserUin,
		"end_point_service_id": expectedEndPointServiceId,
	})
	d.SetId(expectedUserUin + tccommon.FILED_SP + expectedEndPointServiceId)

	err := res.Read(d, newMockMeta())
	assert.NoError(t, err)

	assert.Equal(t, expectedUserUin, d.Get("user_uin"))
	assert.Equal(t, expectedEndPointServiceId, d.Get("end_point_service_id"))
	assert.Equal(t, int(expectedTotalCount), d.Get("total_count"))
}

// go test ./tencentcloud/services/pls/ -run "TestVpcEndPointServiceWhiteListRead_TotalCountNil" -v -count=1 -gcflags="all=-l"
// TestVpcEndPointServiceWhiteListRead_TotalCountNil verifies Read does not panic and leaves total_count unset when TotalCount is nil.
func TestVpcEndPointServiceWhiteListRead_TotalCountNil(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	expectedUserUin := "100020512675"
	expectedEndPointServiceId := "vpcsvc-98jddhcz"

	whiteListRecord := &vpc.VpcEndPointServiceUser{
		UserUin:           ptrString(expectedUserUin),
		EndPointServiceId: ptrString(expectedEndPointServiceId),
		Description:       ptrString("terraform for test"),
	}

	patches.ApplyMethodFunc(&svcvpc.VpcService{}, "DescribeVpcEndPointServiceWhiteListById",
		func(_ context.Context, userUin string, endPointServiceId string) (*vpc.VpcEndPointServiceUser, *uint64, error) {
			return whiteListRecord, nil, nil
		})

	res := pls.ResourceTencentCloudVpcEndPointServiceWhiteList()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"user_uin":             expectedUserUin,
		"end_point_service_id": expectedEndPointServiceId,
	})
	d.SetId(expectedUserUin + tccommon.FILED_SP + expectedEndPointServiceId)

	err := res.Read(d, newMockMeta())
	assert.NoError(t, err)

	// total_count should remain 0 (unset) since the mocked response returned nil totalCount
	assert.Equal(t, 0, d.Get("total_count"))
}
