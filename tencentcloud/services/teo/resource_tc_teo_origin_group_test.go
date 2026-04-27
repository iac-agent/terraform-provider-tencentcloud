package teo_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svcteo "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/teo"

	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/teo"
)

// go test -i; go test -test.run TestAccTencentCloudTeoOriginGroup_basic -v
func TestAccTencentCloudTeoOriginGroup_basic(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PRIVATE) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckOriginGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTeoOriginGroup,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOriginGroupExists("tencentcloud_teo_origin_group.basic"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_origin_group.basic", "zone_id"),
					resource.TestCheckResourceAttr("tencentcloud_teo_origin_group.basic", "name", "keep-group-1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_origin_group.basic", "type", "GENERAL"),
					resource.TestCheckResourceAttr("tencentcloud_teo_origin_group.basic", "records.#", "3"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_origin_group.basic", "records.0.record"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_origin_group.basic", "records.0.type"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_origin_group.basic", "records.0.weight"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_origin_group.basic", "records.0.private"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_origin_group.basic", "records.1.record"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_origin_group.basic", "records.1.type"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_origin_group.basic", "records.1.weight"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_origin_group.basic", "records.1.private"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_origin_group.basic", "records.2.record"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_origin_group.basic", "records.2.type"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_origin_group.basic", "records.2.weight"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_origin_group.basic", "records.2.private"),
				),
			},
			{
				ResourceName:      "tencentcloud_teo_origin_group.basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccTeoOriginGroupUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOriginGroupExists("tencentcloud_teo_origin_group.basic"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_origin_group.basic", "zone_id"),
					resource.TestCheckResourceAttr("tencentcloud_teo_origin_group.basic", "records.0.private", "true"),
					resource.TestCheckResourceAttr("tencentcloud_teo_origin_group.basic", "records.0.private_parameters.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_origin_group.basic", "records.0.private_parameters.0.name", "SecretAccessKey"),
					resource.TestCheckResourceAttr("tencentcloud_teo_origin_group.basic", "records.0.private_parameters.0.value", "test"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_origin_group.basic", "create_time"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_origin_group.basic", "update_time"),
				),
			},
		},
	})
}

func testAccCheckOriginGroupDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	service := svcteo.NewTeoService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_teo_origin_group" {
			continue
		}
		idSplit := strings.Split(rs.Primary.ID, tccommon.FILED_SP)
		if len(idSplit) != 2 {
			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
		}
		zoneId := idSplit[0]
		originGroupId := idSplit[1]

		originGroup, err := service.DescribeTeoOriginGroup(ctx, zoneId, originGroupId)
		if originGroup != nil {
			return fmt.Errorf("zone originGroup %s still exists", rs.Primary.ID)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func testAccCheckOriginGroupExists(r string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}
		idSplit := strings.Split(rs.Primary.ID, tccommon.FILED_SP)
		if len(idSplit) != 2 {
			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
		}
		zoneId := idSplit[0]
		originGroupId := idSplit[1]

		service := svcteo.NewTeoService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
		originGroup, err := service.DescribeTeoOriginGroup(ctx, zoneId, originGroupId)
		if originGroup == nil {
			return fmt.Errorf("zone originGroup %s is not found", rs.Primary.ID)
		}
		if err != nil {
			return err
		}

		return nil
	}
}

const testAccTeoOriginGroup = testAccTeoZone + `

resource "tencentcloud_teo_origin_group" "basic" {
  name    = "keep-group-1"
  type    = "GENERAL"
  zone_id = tencentcloud_teo_zone.basic.id

  records {
    record  = var.zone_name
    type    = "IP_DOMAIN"
    weight  = 100
    private = false
  }
  records {
    private   = false
    record    = "21.1.1.1"
    type      = "IP_DOMAIN"
    weight    = 100
  }
  records {
    private   = false
    record    = "21.1.1.2"
    type      = "IP_DOMAIN"
    weight    = 11
  }
}

`

const testAccTeoOriginGroupUpdate = testAccTeoZone + `

resource "tencentcloud_teo_origin_group" "basic" {
  name    = "keep-group-1"
  type    = "GENERAL"
  zone_id = tencentcloud_teo_zone.basic.id

  records {
    record  = var.zone_name
    type    = "IP_DOMAIN"
    weight  = 100
    private = true

    private_parameters {
      name = "SecretAccessKey"
      value = "test"
    }
  }
}

`

// Unit tests using gomonkey for filters parameter

// mockMetaForOriginGroup implements tccommon.ProviderMeta
type mockMetaForOriginGroup struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaForOriginGroup) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

func newMockMetaForOriginGroup() *mockMetaForOriginGroup {
	return &mockMetaForOriginGroup{client: &connectivity.TencentCloudClient{}}
}

func ptrStringOriginGroup(s string) *string { return &s }
func ptrUint64OriginGroup(u uint64) *uint64  { return &u }

// go test ./tencentcloud/services/teo/ -run "TestTeoOriginGroupRead_NoFilters" -v -count=1 -gcflags="all=-l"
func TestTeoOriginGroupRead_NoFilters(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForOriginGroup().client, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeOriginGroup", func(request *teov20220901.DescribeOriginGroupRequest) (*teov20220901.DescribeOriginGroupResponse, error) {
		resp := teov20220901.NewDescribeOriginGroupResponse()
		resp.Response = &teov20220901.DescribeOriginGroupResponseParams{
			TotalCount: ptrUint64OriginGroup(1),
			OriginGroups: []*teov20220901.OriginGroup{
				{
					GroupId:   ptrStringOriginGroup("origin-test-group-id"),
					Name:      ptrStringOriginGroup("test-origin-group"),
					Type:      ptrStringOriginGroup("GENERAL"),
					HostHeader: ptrStringOriginGroup(""),
					Records: []*teov20220901.OriginRecord{
						{
							Record: ptrStringOriginGroup("1.1.1.1"),
							Type:   ptrStringOriginGroup("IP_DOMAIN"),
							Weight: ptrUint64OriginGroup(100),
						},
					},
				},
			},
			RequestId: ptrStringOriginGroup("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForOriginGroup()
	res := teo.ResourceTencentCloudTeoOriginGroup()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-test123",
		"type":    "GENERAL",
		"records": []interface{}{
			map[string]interface{}{
				"record":  "1.1.1.1",
				"type":    "IP_DOMAIN",
				"weight":  100,
				"private": false,
			},
		},
	})
	d.SetId("zone-test123#origin-test-group-id")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-test123#origin-test-group-id", d.Id())
	assert.Equal(t, "origin-test-group-id", d.Get("origin_group_id"))
	assert.Equal(t, "test-origin-group", d.Get("name"))
	assert.Equal(t, "GENERAL", d.Get("type"))
}

// go test ./tencentcloud/services/teo/ -run "TestTeoOriginGroupRead_WithFilters" -v -count=1 -gcflags="all=-l"
func TestTeoOriginGroupRead_WithFilters(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForOriginGroup().client, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeOriginGroup", func(request *teov20220901.DescribeOriginGroupRequest) (*teov20220901.DescribeOriginGroupResponse, error) {
		resp := teov20220901.NewDescribeOriginGroupResponse()
		resp.Response = &teov20220901.DescribeOriginGroupResponseParams{
			TotalCount: ptrUint64OriginGroup(1),
			OriginGroups: []*teov20220901.OriginGroup{
				{
					GroupId:    ptrStringOriginGroup("origin-test-group-id"),
					Name:       ptrStringOriginGroup("my-origin-group"),
					Type:       ptrStringOriginGroup("GENERAL"),
					HostHeader: ptrStringOriginGroup(""),
					Records: []*teov20220901.OriginRecord{
						{
							Record: ptrStringOriginGroup("2.2.2.2"),
							Type:   ptrStringOriginGroup("IP_DOMAIN"),
							Weight: ptrUint64OriginGroup(80),
						},
					},
				},
			},
			RequestId: ptrStringOriginGroup("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForOriginGroup()
	res := teo.ResourceTencentCloudTeoOriginGroup()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-test456",
		"type":    "GENERAL",
		"records": []interface{}{
			map[string]interface{}{
				"record":  "2.2.2.2",
				"type":    "IP_DOMAIN",
				"weight":  80,
				"private": false,
			},
		},
		"filters": []interface{}{
			map[string]interface{}{
				"name":   "origin-group-name",
				"values": []interface{}{"my-origin-group"},
				"fuzzy":  true,
			},
		},
	})
	d.SetId("zone-test456#origin-test-group-id")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-test456#origin-test-group-id", d.Id())
	assert.Equal(t, "my-origin-group", d.Get("name"))
	assert.Equal(t, "GENERAL", d.Get("type"))
}

// go test ./tencentcloud/services/teo/ -run "TestTeoOriginGroupRead_NotFound" -v -count=1 -gcflags="all=-l"
func TestTeoOriginGroupRead_NotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForOriginGroup().client, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeOriginGroup", func(request *teov20220901.DescribeOriginGroupRequest) (*teov20220901.DescribeOriginGroupResponse, error) {
		resp := teov20220901.NewDescribeOriginGroupResponse()
		resp.Response = &teov20220901.DescribeOriginGroupResponseParams{
			TotalCount:  ptrUint64OriginGroup(0),
			OriginGroups: []*teov20220901.OriginGroup{},
			RequestId:   ptrStringOriginGroup("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForOriginGroup()
	res := teo.ResourceTencentCloudTeoOriginGroup()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-test789",
		"type":    "GENERAL",
		"records": []interface{}{
			map[string]interface{}{
				"record":  "3.3.3.3",
				"type":    "IP_DOMAIN",
				"weight":  50,
				"private": false,
			},
		},
	})
	d.SetId("zone-test789#origin-nonexistent-id")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}
