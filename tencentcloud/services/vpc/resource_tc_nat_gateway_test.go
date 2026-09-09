package vpc_test

import (
	"context"
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	svctag "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tag"
	svcvpc "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/vpc"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func init() {
	resource.AddTestSweepers("tencentcloud_nat", &resource.Sweeper{
		Name: "tencentcloud_nat",
		F:    testSweepNatInstance,
	})
}

func testSweepNatInstance(region string) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	sharedClient, err := tcacctest.SharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("getting tencentcloud client error: %s", err.Error())
	}
	client := sharedClient.(tccommon.ProviderMeta).GetAPIV3Conn()

	vpcService := svcvpc.NewVpcService(client)

	instances, err := vpcService.DescribeNatGatewayByFilter(ctx, nil)
	if err != nil {
		return fmt.Errorf("get instance list error: %s", err.Error())
	}

	// add scanning resources
	var resources, nonKeepResources []*tccommon.ResourceInstance
	for _, v := range instances {
		if !tccommon.CheckResourcePersist(*v.NatGatewayName, *v.CreatedTime) {
			nonKeepResources = append(nonKeepResources, &tccommon.ResourceInstance{
				Id:   *v.NatGatewayId,
				Name: *v.NatGatewayName,
			})
		}
		resources = append(resources, &tccommon.ResourceInstance{
			Id:         *v.NatGatewayId,
			Name:       *v.NatGatewayName,
			CreateTime: *v.CreatedTime,
		})
	}
	tccommon.ProcessScanCloudResources(client, resources, nonKeepResources, "CreateNatGateway")

	for _, v := range instances {
		instanceId := *v.NatGatewayId
		instanceName := v.NatGatewayName

		now := time.Now()

		createTime := tccommon.StringToTime(*v.CreatedTime)
		interval := now.Sub(createTime).Minutes()
		if instanceName != nil {
			if strings.HasPrefix(*instanceName, tcacctest.KeepResource) || strings.HasPrefix(*instanceName, tcacctest.DefaultResource) {
				continue
			}
		}

		// less than 30 minute, not delete
		if tccommon.NeedProtect == 1 && int64(interval) < 30 {
			continue
		}

		if err = vpcService.DeleteNatGateway(ctx, instanceId); err != nil {
			log.Printf("[ERROR] sweep instance %s error: %s", instanceId, err.Error())
		}
	}
	return nil
}

func TestAccTencentCloudNatGateway_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckNatGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNatGatewayConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNatGatewayExists("tencentcloud_nat_gateway.my_nat"),
					resource.TestCheckResourceAttr("tencentcloud_nat_gateway.my_nat", "name", "terraform_test"),
					resource.TestCheckResourceAttr("tencentcloud_nat_gateway.my_nat", "max_concurrent", "3000000"),
					resource.TestCheckResourceAttr("tencentcloud_nat_gateway.my_nat", "bandwidth", "500"),
					resource.TestCheckResourceAttr("tencentcloud_nat_gateway.my_nat", "assigned_eip_set.#", "2"),
					resource.TestCheckResourceAttr("tencentcloud_nat_gateway.my_nat", "tags.tf", "test"),
				),
			},
			{
				Config: testAccNatGatewayConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNatGatewayExists("tencentcloud_nat_gateway.my_nat"),
					resource.TestCheckResourceAttr("tencentcloud_nat_gateway.my_nat", "name", "new_name"),
					resource.TestCheckResourceAttr("tencentcloud_nat_gateway.my_nat", "max_concurrent", "10000000"),
					resource.TestCheckResourceAttr("tencentcloud_nat_gateway.my_nat", "bandwidth", "1000"),
					resource.TestCheckResourceAttr("tencentcloud_nat_gateway.my_nat", "assigned_eip_set.#", "2"),
					resource.TestCheckResourceAttr("tencentcloud_nat_gateway.my_nat", "tags.tf", "teest"),
				),
			},
		},
	})
}

func testAccCheckNatGatewayDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)

	conn := tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn()
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_nat_gateway" {
			continue
		}
		request := vpc.NewDescribeNatGatewaysRequest()
		request.NatGatewayIds = []*string{&rs.Primary.ID}
		var response *vpc.DescribeNatGatewaysResponse
		err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			result, e := conn.UseVpcClient().DescribeNatGateways(request)
			if e != nil {
				log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
					logId, request.GetAction(), request.ToJsonString(), e.Error())
				return tccommon.RetryError(e)
			}
			response = result
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s read nat gateway failed, reason:%s\n ", logId, err.Error())
			return err
		}
		if len(response.Response.NatGatewaySet) != 0 {
			return fmt.Errorf("nat gateway id is still exists")
		}

	}
	return nil
}

func testAccCheckNatGatewayExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("nat gateway instance %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("nat gateway id is not set")
		}
		conn := tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn()
		request := vpc.NewDescribeNatGatewaysRequest()
		request.NatGatewayIds = []*string{&rs.Primary.ID}
		var response *vpc.DescribeNatGatewaysResponse
		err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			result, e := conn.UseVpcClient().DescribeNatGateways(request)
			if e != nil {
				log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
					logId, request.GetAction(), request.ToJsonString(), e.Error())
				return tccommon.RetryError(e)
			}
			response = result
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s read nat gateway failed, reason:%s\n ", logId, err.Error())
			return err
		}
		if len(response.Response.NatGatewaySet) != 1 {
			return fmt.Errorf("nat gateway id is not found")
		}
		return nil
	}
}

const testAccNatGatewayConfig = `
data "tencentcloud_vpc_instances" "foo" {
  name = "Default-VPC"
}
# Create EIP 
resource "tencentcloud_eip" "eip_dev_dnat" {
  name = "terraform_test"
}
resource "tencentcloud_eip" "eip_test_dnat" {
  name = "terraform_test"
}
resource "tencentcloud_nat_gateway" "my_nat" {
  vpc_id         = data.tencentcloud_vpc_instances.foo.instance_list.0.vpc_id
  name           = "terraform_test"
  max_concurrent = 3000000
  bandwidth      = 500

  assigned_eip_set = [
    tencentcloud_eip.eip_dev_dnat.public_ip,
    tencentcloud_eip.eip_test_dnat.public_ip,
  ]
  tags = {
    tf = "test"
  }
}
`
const testAccNatGatewayConfigUpdate = `
data "tencentcloud_vpc_instances" "foo" {
  name = "Default-VPC"
}
# Create EIP 
resource "tencentcloud_eip" "eip_dev_dnat" {
  name = "terraform_test"
}
resource "tencentcloud_eip" "new_eip" {
  name = "terraform_test"
}

resource "tencentcloud_nat_gateway" "my_nat" {
  vpc_id         = data.tencentcloud_vpc_instances.foo.instance_list.0.vpc_id
  name           = "new_name"
  max_concurrent = 10000000
  bandwidth      = 1000

  assigned_eip_set = [
    tencentcloud_eip.eip_dev_dnat.public_ip,
    tencentcloud_eip.new_eip.public_ip,
  ]
  tags = {
    tf = "teest"
  }
}
`

// ---- gomonkey-based unit tests for deletion_protection_enabled ----
// Run with: go test ./tencentcloud/services/vpc/ -run "TestNatGatewayDeletionProtection" -v -count=1 -gcflags="all=-l"

type mockMetaNatGateway struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaNatGateway) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaNatGateway{}

func newMockMetaNatGateway() *mockMetaNatGateway {
	return &mockMetaNatGateway{
		client: &connectivity.TencentCloudClient{
			Region: "ap-guangzhou",
		},
	}
}

func ptrStringNatGW(s string) *string { return &s }
func ptrUint64NatGW(v uint64) *uint64 { return &v }
func ptrBoolNatGW(b bool) *bool       { return &b }

func buildMockNatGateway(natGatewayId, natGatewayName, vpcId, zone, createdTime string, maxConcurrent, bandwidth uint64, deletionProtection bool) *vpc.NatGateway {
	return &vpc.NatGateway{
		NatGatewayId:              ptrStringNatGW(natGatewayId),
		NatGatewayName:            ptrStringNatGW(natGatewayName),
		VpcId:                     ptrStringNatGW(vpcId),
		Zone:                      ptrStringNatGW(zone),
		CreatedTime:               ptrStringNatGW(createdTime),
		State:                     ptrStringNatGW("AVAILABLE"),
		MaxConcurrentConnection:   ptrUint64NatGW(maxConcurrent),
		InternetMaxBandwidthOut:   ptrUint64NatGW(bandwidth),
		NatProductVersion:         ptrUint64NatGW(1),
		DeletionProtectionEnabled: ptrBoolNatGW(deletionProtection),
		PublicIpAddressSet: []*vpc.NatGatewayAddress{
			{
				AddressId:       ptrStringNatGW("eip-12345678"),
				PublicIpAddress: ptrStringNatGW("1.2.3.4"),
			},
		},
	}
}

// TestNatGatewayDeletionProtection_Create verifies deletion_protection_enabled is passed to CreateNatGateway API during creation.
func TestNatGatewayDeletionProtection_Create(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	vpcClient := &vpc.Client{}
	patches.ApplyMethodReturn(newMockMetaNatGateway().client, "UseVpcClient", vpcClient)

	var capturedRequest *vpc.CreateNatGatewayRequest
	patches.ApplyMethodFunc(vpcClient, "CreateNatGateway", func(request *vpc.CreateNatGatewayRequest) (*vpc.CreateNatGatewayResponse, error) {
		capturedRequest = request
		resp := vpc.NewCreateNatGatewayResponse()
		resp.Response = &vpc.CreateNatGatewayResponseParams{
			NatGatewaySet: []*vpc.NatGateway{
				buildMockNatGateway("nat-df45454", "tf_nat_gateway", "vpc-123", "ap-guangzhou-3", "2024-01-01 00:00:00", 1000000, 100, true),
			},
			RequestId: ptrStringNatGW("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(vpcClient, "DescribeNatGateways", func(request *vpc.DescribeNatGatewaysRequest) (*vpc.DescribeNatGatewaysResponse, error) {
		resp := vpc.NewDescribeNatGatewaysResponse()
		resp.Response = &vpc.DescribeNatGatewaysResponseParams{
			NatGatewaySet: []*vpc.NatGateway{
				buildMockNatGateway("nat-df45454", "tf_nat_gateway", "vpc-123", "ap-guangzhou-3", "2024-01-01 00:00:00", 1000000, 100, true),
			},
			RequestId: ptrStringNatGW("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(vpcClient, "DescribeAddresses", func(request *vpc.DescribeAddressesRequest) (*vpc.DescribeAddressesResponse, error) {
		resp := vpc.NewDescribeAddressesResponse()
		resp.Response = &vpc.DescribeAddressesResponseParams{
			AddressSet: []*vpc.Address{
				{
					AddressId: ptrStringNatGW("eip-12345678"),
					PublicIp:  ptrStringNatGW("1.2.3.4"),
					Bandwidth: ptrUint64NatGW(100),
				},
			},
			RequestId: ptrStringNatGW("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(&svctag.TagService{}, "DescribeResourceTags", func(_ context.Context, _, _, _, _ string) (map[string]string, error) {
		return map[string]string{}, nil
	})

	meta := newMockMetaNatGateway()
	res := svcvpc.ResourceTencentCloudNatGateway()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"vpc_id":                      "vpc-123",
		"name":                        "tf_nat_gateway",
		"assigned_eip_set":            schema.NewSet(schema.HashString, []interface{}{"1.2.3.4"}),
		"deletion_protection_enabled": true,
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "nat-df45454", d.Id())

	assert.NotNil(t, capturedRequest.DeletionProtectionEnabled)
	assert.True(t, *capturedRequest.DeletionProtectionEnabled)

	assert.Equal(t, true, d.Get("deletion_protection_enabled").(bool))
}

// TestNatGatewayDeletionProtection_CreateDefault verifies deletion_protection_enabled is not set on the request when not specified.
func TestNatGatewayDeletionProtection_CreateDefault(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	vpcClient := &vpc.Client{}
	patches.ApplyMethodReturn(newMockMetaNatGateway().client, "UseVpcClient", vpcClient)

	var capturedRequest *vpc.CreateNatGatewayRequest
	patches.ApplyMethodFunc(vpcClient, "CreateNatGateway", func(request *vpc.CreateNatGatewayRequest) (*vpc.CreateNatGatewayResponse, error) {
		capturedRequest = request
		resp := vpc.NewCreateNatGatewayResponse()
		resp.Response = &vpc.CreateNatGatewayResponseParams{
			NatGatewaySet: []*vpc.NatGateway{
				buildMockNatGateway("nat-df45454", "tf_nat_gateway", "vpc-123", "ap-guangzhou-3", "2024-01-01 00:00:00", 1000000, 100, false),
			},
			RequestId: ptrStringNatGW("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(vpcClient, "DescribeNatGateways", func(request *vpc.DescribeNatGatewaysRequest) (*vpc.DescribeNatGatewaysResponse, error) {
		resp := vpc.NewDescribeNatGatewaysResponse()
		resp.Response = &vpc.DescribeNatGatewaysResponseParams{
			NatGatewaySet: []*vpc.NatGateway{
				buildMockNatGateway("nat-df45454", "tf_nat_gateway", "vpc-123", "ap-guangzhou-3", "2024-01-01 00:00:00", 1000000, 100, false),
			},
			RequestId: ptrStringNatGW("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(vpcClient, "DescribeAddresses", func(request *vpc.DescribeAddressesRequest) (*vpc.DescribeAddressesResponse, error) {
		resp := vpc.NewDescribeAddressesResponse()
		resp.Response = &vpc.DescribeAddressesResponseParams{
			AddressSet: []*vpc.Address{
				{
					AddressId: ptrStringNatGW("eip-12345678"),
					PublicIp:  ptrStringNatGW("1.2.3.4"),
					Bandwidth: ptrUint64NatGW(100),
				},
			},
			RequestId: ptrStringNatGW("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(&svctag.TagService{}, "DescribeResourceTags", func(_ context.Context, _, _, _, _ string) (map[string]string, error) {
		return map[string]string{}, nil
	})

	meta := newMockMetaNatGateway()
	res := svcvpc.ResourceTencentCloudNatGateway()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"vpc_id":           "vpc-123",
		"name":             "tf_nat_gateway",
		"assigned_eip_set": schema.NewSet(schema.HashString, []interface{}{"1.2.3.4"}),
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)

	assert.Nil(t, capturedRequest.DeletionProtectionEnabled)
}

// TestNatGatewayDeletionProtection_ReadValue verifies deletion_protection_enabled is read from DescribeNatGateways response.
func TestNatGatewayDeletionProtection_ReadValue(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	vpcClient := &vpc.Client{}
	patches.ApplyMethodReturn(newMockMetaNatGateway().client, "UseVpcClient", vpcClient)

	patches.ApplyMethodFunc(vpcClient, "DescribeNatGateways", func(request *vpc.DescribeNatGatewaysRequest) (*vpc.DescribeNatGatewaysResponse, error) {
		resp := vpc.NewDescribeNatGatewaysResponse()
		resp.Response = &vpc.DescribeNatGatewaysResponseParams{
			NatGatewaySet: []*vpc.NatGateway{
				buildMockNatGateway("nat-df45454", "tf_nat_gateway", "vpc-123", "ap-guangzhou-3", "2024-01-01 00:00:00", 1000000, 100, true),
			},
			RequestId: ptrStringNatGW("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(vpcClient, "DescribeAddresses", func(request *vpc.DescribeAddressesRequest) (*vpc.DescribeAddressesResponse, error) {
		resp := vpc.NewDescribeAddressesResponse()
		resp.Response = &vpc.DescribeAddressesResponseParams{
			AddressSet: []*vpc.Address{
				{
					AddressId: ptrStringNatGW("eip-12345678"),
					PublicIp:  ptrStringNatGW("1.2.3.4"),
					Bandwidth: ptrUint64NatGW(100),
				},
			},
			RequestId: ptrStringNatGW("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(&svctag.TagService{}, "DescribeResourceTags", func(_ context.Context, _, _, _, _ string) (map[string]string, error) {
		return map[string]string{}, nil
	})

	meta := newMockMetaNatGateway()
	res := svcvpc.ResourceTencentCloudNatGateway()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"vpc_id":                      "vpc-123",
		"name":                        "tf_nat_gateway",
		"assigned_eip_set":            schema.NewSet(schema.HashString, []interface{}{"1.2.3.4"}),
		"deletion_protection_enabled": false,
	})
	d.SetId("nat-df45454")

	err := res.Read(d, meta)
	assert.NoError(t, err)

	assert.Equal(t, true, d.Get("deletion_protection_enabled").(bool))
}

// TestNatGatewayDeletionProtection_ReadNil verifies Read skips setting deletion_protection_enabled when API returns nil.
func TestNatGatewayDeletionProtection_ReadNil(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	vpcClient := &vpc.Client{}
	patches.ApplyMethodReturn(newMockMetaNatGateway().client, "UseVpcClient", vpcClient)

	nat := buildMockNatGateway("nat-df45454", "tf_nat_gateway", "vpc-123", "ap-guangzhou-3", "2024-01-01 00:00:00", 1000000, 100, false)
	nat.DeletionProtectionEnabled = nil

	patches.ApplyMethodFunc(vpcClient, "DescribeNatGateways", func(request *vpc.DescribeNatGatewaysRequest) (*vpc.DescribeNatGatewaysResponse, error) {
		resp := vpc.NewDescribeNatGatewaysResponse()
		resp.Response = &vpc.DescribeNatGatewaysResponseParams{
			NatGatewaySet: []*vpc.NatGateway{nat},
			RequestId:     ptrStringNatGW("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(vpcClient, "DescribeAddresses", func(request *vpc.DescribeAddressesRequest) (*vpc.DescribeAddressesResponse, error) {
		resp := vpc.NewDescribeAddressesResponse()
		resp.Response = &vpc.DescribeAddressesResponseParams{
			AddressSet: []*vpc.Address{
				{
					AddressId: ptrStringNatGW("eip-12345678"),
					PublicIp:  ptrStringNatGW("1.2.3.4"),
					Bandwidth: ptrUint64NatGW(100),
				},
			},
			RequestId: ptrStringNatGW("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(&svctag.TagService{}, "DescribeResourceTags", func(_ context.Context, _, _, _, _ string) (map[string]string, error) {
		return map[string]string{}, nil
	})

	meta := newMockMetaNatGateway()
	res := svcvpc.ResourceTencentCloudNatGateway()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"vpc_id":                      "vpc-123",
		"name":                        "tf_nat_gateway",
		"assigned_eip_set":            schema.NewSet(schema.HashString, []interface{}{"1.2.3.4"}),
		"deletion_protection_enabled": true,
	})
	d.SetId("nat-df45454")

	err := res.Read(d, meta)
	assert.NoError(t, err)

	assert.Equal(t, true, d.Get("deletion_protection_enabled").(bool))
}

// TestNatGatewayDeletionProtection_Update verifies deletion_protection_enabled change triggers ModifyNatGatewayAttribute.
func TestNatGatewayDeletionProtection_Update(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	vpcClient := &vpc.Client{}
	patches.ApplyMethodReturn(newMockMetaNatGateway().client, "UseVpcClient", vpcClient)

	var capturedRequest *vpc.ModifyNatGatewayAttributeRequest
	patches.ApplyMethodFunc(vpcClient, "ModifyNatGatewayAttribute", func(request *vpc.ModifyNatGatewayAttributeRequest) (*vpc.ModifyNatGatewayAttributeResponse, error) {
		capturedRequest = request
		resp := vpc.NewModifyNatGatewayAttributeResponse()
		resp.Response = &vpc.ModifyNatGatewayAttributeResponseParams{
			RequestId: ptrStringNatGW("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(vpcClient, "DescribeNatGateways", func(request *vpc.DescribeNatGatewaysRequest) (*vpc.DescribeNatGatewaysResponse, error) {
		resp := vpc.NewDescribeNatGatewaysResponse()
		resp.Response = &vpc.DescribeNatGatewaysResponseParams{
			NatGatewaySet: []*vpc.NatGateway{
				buildMockNatGateway("nat-df45454", "tf_nat_gateway", "vpc-123", "ap-guangzhou-3", "2024-01-01 00:00:00", 1000000, 100, true),
			},
			RequestId: ptrStringNatGW("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(vpcClient, "DescribeAddresses", func(request *vpc.DescribeAddressesRequest) (*vpc.DescribeAddressesResponse, error) {
		resp := vpc.NewDescribeAddressesResponse()
		resp.Response = &vpc.DescribeAddressesResponseParams{
			AddressSet: []*vpc.Address{
				{
					AddressId: ptrStringNatGW("eip-12345678"),
					PublicIp:  ptrStringNatGW("1.2.3.4"),
					Bandwidth: ptrUint64NatGW(100),
				},
			},
			RequestId: ptrStringNatGW("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(&svctag.TagService{}, "DescribeResourceTags", func(_ context.Context, _, _, _, _ string) (map[string]string, error) {
		return map[string]string{}, nil
	})

	meta := newMockMetaNatGateway()
	res := svcvpc.ResourceTencentCloudNatGateway()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"vpc_id":                      "vpc-123",
		"name":                        "tf_nat_gateway",
		"assigned_eip_set":            schema.NewSet(schema.HashString, []interface{}{"1.2.3.4"}),
		"deletion_protection_enabled": true,
	})
	d.SetId("nat-df45454")

	patches.ApplyMethodFunc(d, "HasChange", func(key string) bool {
		return key == "deletion_protection_enabled"
	})

	err := res.Update(d, meta)
	assert.NoError(t, err)

	assert.NotNil(t, capturedRequest)
	assert.NotNil(t, capturedRequest.DeletionProtectionEnabled)
	assert.True(t, *capturedRequest.DeletionProtectionEnabled)

	assert.Equal(t, true, d.Get("deletion_protection_enabled").(bool))
}

// TestNatGatewayDeletionProtection_Schema validates the deletion_protection_enabled schema definition.
func TestNatGatewayDeletionProtection_Schema(t *testing.T) {
	res := svcvpc.ResourceTencentCloudNatGateway()

	field, ok := res.Schema["deletion_protection_enabled"]
	assert.True(t, ok, "deletion_protection_enabled should exist in schema")
	assert.Equal(t, schema.TypeBool, field.Type)
	assert.True(t, field.Optional)
	assert.True(t, field.Computed)
	assert.False(t, field.ForceNew, "deletion_protection_enabled should not be ForceNew")
}
