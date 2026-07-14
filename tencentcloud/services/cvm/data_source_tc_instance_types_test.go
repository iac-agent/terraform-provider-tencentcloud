package cvm_test

import (
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	cbs "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cbs/v20170312"
	cvmSDK "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"

	acctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/cvm"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// mockMeta implements tccommon.ProviderMeta
type instanceTypesMockMeta struct {
	client *connectivity.TencentCloudClient
}

func (m *instanceTypesMockMeta) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &instanceTypesMockMeta{}

func newInstanceTypesMockMeta() *instanceTypesMockMeta {
	return &instanceTypesMockMeta{client: &connectivity.TencentCloudClient{}}
}

func ptrString(v string) *string { return &v }
func ptrInt64(v int64) *int64    { return &v }

// TestInstanceTypesInquiryType_Specified tests that inquiry_type is correctly passed to DescribeDiskConfigQuota when specified
func TestInstanceTypesInquiryType_Specified(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	// Patch UseCvmClient to return a mock CVM client
	cvmClient := &cvmSDK.Client{}
	patches.ApplyMethodReturn(newInstanceTypesMockMeta().client, "UseCvmClient", cvmClient)

	// Patch DescribeZoneInstanceConfigInfos to return mock instance types
	patches.ApplyMethodFunc(cvmClient, "DescribeZoneInstanceConfigInfos", func(request *cvmSDK.DescribeZoneInstanceConfigInfosRequest) (*cvmSDK.DescribeZoneInstanceConfigInfosResponse, error) {
		resp := cvmSDK.NewDescribeZoneInstanceConfigInfosResponse()
		resp.Response = &cvmSDK.DescribeZoneInstanceConfigInfosResponseParams{
			InstanceTypeQuotaSet: []*cvmSDK.InstanceTypeQuotaItem{
				{
					Zone:           ptrString("ap-guangzhou-6"),
					Cpu:            ptrInt64(2),
					Memory:         ptrInt64(2),
					InstanceFamily: ptrString("S6"),
					InstanceType:   ptrString("S6.SMALL1"),
					Status:         ptrString("Sell"),
				},
			},
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	// Patch UseCbsClient to return a mock CBS client
	cbsClient := &cbs.Client{}
	patches.ApplyMethodReturn(newInstanceTypesMockMeta().client, "UseCbsClient", cbsClient)

	// Patch DescribeDiskConfigQuota to verify inquiry_type is correctly passed
	patches.ApplyMethodFunc(cbsClient, "DescribeDiskConfigQuota", func(request *cbs.DescribeDiskConfigQuotaRequest) (*cbs.DescribeDiskConfigQuotaResponse, error) {
		// Verify that InquiryType is set to INQUIRY_CBS_CONFIG
		assert.NotNil(t, request.InquiryType)
		assert.Equal(t, "INQUIRY_CBS_CONFIG", *request.InquiryType)

		resp := cbs.NewDescribeDiskConfigQuotaResponse()
		resp.Response = &cbs.DescribeDiskConfigQuotaResponseParams{
			DiskConfigSet: []*cbs.DiskConfig{
				{
					Available:      ptrBool(true),
					DiskChargeType: ptrString("PREPAID"),
					Zone:           ptrString("ap-guangzhou-6"),
					InstanceFamily: ptrString("S6"),
					DiskType:       ptrString("CLOUD_SSD"),
					DiskUsage:      ptrString("SYSTEM_DISK"),
					MinDiskSize:    ptrUint64(50),
					MaxDiskSize:    ptrUint64(500),
				},
			},
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newInstanceTypesMockMeta()
	res := cvm.DataSourceTencentCloudInstanceTypes()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"cpu_core_count":    2,
		"memory_size":       2,
		"exclude_sold_out":  false,
		"availability_zone": "ap-guangzhou-6",
		"cbs_filter": []interface{}{
			map[string]interface{}{
				"disk_types":       []interface{}{"CLOUD_SSD"},
				"disk_charge_type": "PREPAID",
				"disk_usage":       "SYSTEM_DISK",
				"inquiry_type":     "INQUIRY_CBS_CONFIG",
			},
		},
	})

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.NotEmpty(t, d.Id())
}

// TestInstanceTypesInquiryType_Default tests that inquiry_type defaults to INQUIRY_CVM_CONFIG when not specified
func TestInstanceTypesInquiryType_Default(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	// Patch UseCvmClient to return a mock CVM client
	cvmClient := &cvmSDK.Client{}
	patches.ApplyMethodReturn(newInstanceTypesMockMeta().client, "UseCvmClient", cvmClient)

	// Patch DescribeZoneInstanceConfigInfos to return mock instance types
	patches.ApplyMethodFunc(cvmClient, "DescribeZoneInstanceConfigInfos", func(request *cvmSDK.DescribeZoneInstanceConfigInfosRequest) (*cvmSDK.DescribeZoneInstanceConfigInfosResponse, error) {
		resp := cvmSDK.NewDescribeZoneInstanceConfigInfosResponse()
		resp.Response = &cvmSDK.DescribeZoneInstanceConfigInfosResponseParams{
			InstanceTypeQuotaSet: []*cvmSDK.InstanceTypeQuotaItem{
				{
					Zone:           ptrString("ap-guangzhou-6"),
					Cpu:            ptrInt64(2),
					Memory:         ptrInt64(2),
					InstanceFamily: ptrString("S6"),
					InstanceType:   ptrString("S6.SMALL1"),
					Status:         ptrString("Sell"),
				},
			},
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	// Patch UseCbsClient to return a mock CBS client
	cbsClient := &cbs.Client{}
	patches.ApplyMethodReturn(newInstanceTypesMockMeta().client, "UseCbsClient", cbsClient)

	// Patch DescribeDiskConfigQuota to verify inquiry_type defaults to INQUIRY_CVM_CONFIG
	patches.ApplyMethodFunc(cbsClient, "DescribeDiskConfigQuota", func(request *cbs.DescribeDiskConfigQuotaRequest) (*cbs.DescribeDiskConfigQuotaResponse, error) {
		// Verify that InquiryType defaults to INQUIRY_CVM_CONFIG when not specified
		assert.NotNil(t, request.InquiryType)
		assert.Equal(t, "INQUIRY_CVM_CONFIG", *request.InquiryType)

		resp := cbs.NewDescribeDiskConfigQuotaResponse()
		resp.Response = &cbs.DescribeDiskConfigQuotaResponseParams{
			DiskConfigSet: []*cbs.DiskConfig{
				{
					Available:      ptrBool(true),
					DiskChargeType: ptrString("PREPAID"),
					Zone:           ptrString("ap-guangzhou-6"),
					InstanceFamily: ptrString("S6"),
					DiskType:       ptrString("CLOUD_SSD"),
					DiskUsage:      ptrString("SYSTEM_DISK"),
					MinDiskSize:    ptrUint64(50),
					MaxDiskSize:    ptrUint64(500),
				},
			},
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newInstanceTypesMockMeta()
	res := cvm.DataSourceTencentCloudInstanceTypes()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"cpu_core_count":    2,
		"memory_size":       2,
		"exclude_sold_out":  false,
		"availability_zone": "ap-guangzhou-6",
		"cbs_filter": []interface{}{
			map[string]interface{}{
				"disk_types":       []interface{}{"CLOUD_SSD"},
				"disk_charge_type": "PREPAID",
				"disk_usage":       "SYSTEM_DISK",
			},
		},
	})

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.NotEmpty(t, d.Id())
}

// TestInstanceTypesInquiryType_Schema validates inquiry_type is in the schema
func TestInstanceTypesInquiryType_Schema(t *testing.T) {
	res := cvm.DataSourceTencentCloudInstanceTypes()

	assert.NotNil(t, res)
	assert.Contains(t, res.Schema, "cbs_filter")

	cbsFilter := res.Schema["cbs_filter"]
	assert.Equal(t, schema.TypeList, cbsFilter.Type)

	cbsFilterElem := cbsFilter.Elem.(*schema.Resource)
	assert.Contains(t, cbsFilterElem.Schema, "inquiry_type")

	inquiryType := cbsFilterElem.Schema["inquiry_type"]
	assert.Equal(t, schema.TypeString, inquiryType.Type)
	assert.True(t, inquiryType.Optional)
	assert.False(t, inquiryType.Required)
}

func ptrBool(v bool) *bool       { return &v }
func ptrUint64(v uint64) *uint64 { return &v }

// go test -i; go test -test.run TestAccTencentCloudInstanceTypesDataSource_basic -v
func TestAccTencentCloudCvmInstanceTypesDataSource_Basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.AccPreCheck(t)
		},
		Providers: acctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCvmInstanceTypesDataSource_BasicCreate,
				Check:  resource.ComposeTestCheckFunc(acctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_instance_types.example"), resource.TestCheckResourceAttr("data.tencentcloud_instance_types.example", "instance_types.0.cpu_core_count", "4"), resource.TestCheckResourceAttr("data.tencentcloud_instance_types.example", "instance_types.0.memory_size", "8"), resource.TestCheckResourceAttr("data.tencentcloud_instance_types.example", "instance_types.0.availability_zone", "ap-guangzhou-3")),
			},
		},
	})
}

const testAccCvmInstanceTypesDataSource_BasicCreate = `

data "tencentcloud_instance_types" "example" {
    availability_zone = "ap-guangzhou-3"
    cpu_core_count = 4
    memory_size = 8
}

`

func TestAccTencentCloudCvmInstanceTypesDataSource_Sell(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.AccPreCheck(t)
		},
		Providers: acctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCvmInstanceTypesDataSource_SellCreate,
				Check:  resource.ComposeTestCheckFunc(acctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_instance_types.example"), resource.TestCheckResourceAttr("data.tencentcloud_instance_types.example", "instance_types.0.cpu_core_count", "2"), resource.TestCheckResourceAttr("data.tencentcloud_instance_types.example", "instance_types.0.memory_size", "2"), resource.TestCheckResourceAttr("data.tencentcloud_instance_types.example", "instance_types.0.availability_zone", "ap-guangzhou-3"), resource.TestCheckResourceAttr("data.tencentcloud_instance_types.example", "instance_types.0.family", "SA2")),
			},
		},
	})
}
func TestAccTencentCloudCvmInstanceTypesDataSource_WithCbsFilter(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.AccPreCheck(t)
		},
		Providers: acctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCvmInstanceTypesDataSource_WithCbsFilter,
				Check: resource.ComposeTestCheckFunc(
					acctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_instance_types.with_cbs_filter"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_instance_types.with_cbs_filter", "instance_types.0.cbs_configs.#"),
					resource.TestCheckResourceAttr("data.tencentcloud_instance_types.with_cbs_filter", "instance_types.0.cbs_configs.0.disk_type", "CLOUD_SSD"),
					resource.TestCheckResourceAttr("data.tencentcloud_instance_types.with_cbs_filter", "instance_types.0.cbs_configs.0.disk_charge_type", "PREPAID"),
					resource.TestCheckResourceAttr("data.tencentcloud_instance_types.with_cbs_filter", "instance_types.0.cbs_configs.0.disk_usage", "SYSTEM_DISK"),
				),
			},
		},
	})
}

const testAccCvmInstanceTypesDataSource_SellCreate = `

data "tencentcloud_instance_types" "example" {
    cpu_core_count = 2
    memory_size = 2
    exclude_sold_out = true
    
    filter {
        name = "instance-family"
        values = ["SA2"]
    }
    filter {
        name = "zone"
        values = ["ap-guangzhou-3"]
    }
}

`

const testAccCvmInstanceTypesDataSource_WithCbsFilter = `

data "tencentcloud_instance_types" "with_cbs_filter" {
    cpu_core_count = 2
    memory_size = 2
    exclude_sold_out = true
    
    filter {
        name = "instance-family"
        values = ["S6"]
    }
    filter {
        name = "zone"
        values = ["ap-guangzhou-6"]
    }
	cbs_filter {
        disk_types = ["CLOUD_SSD"]
        disk_charge_type = "PREPAID"
        disk_usage = "SYSTEM_DISK"
    }
}
`
