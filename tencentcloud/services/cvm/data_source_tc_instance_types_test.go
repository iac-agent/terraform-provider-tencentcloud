package cvm_test

import (
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	cbsSDK "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cbs/v20170312"
	cvmSDK "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"

	acctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/cvm"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

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

// instanceTypesMockMeta implements tccommon.ProviderMeta
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

func ptrString(v string) *string {
	return &v
}

func ptrInt64(v int64) *int64 {
	return &v
}

func ptrUint64(v uint64) *uint64 {
	return &v
}

func ptrBool(v bool) *bool {
	return &v
}

func ptrFloat64(v float64) *float64 {
	return &v
}

// go test ./tencentcloud/services/cvm/ -run "TestCvmInstanceTypes_InstanceFamilies" -v -count=1 -gcflags="all=-l"

// TestCvmInstanceTypes_Read_WithInstanceFamiliesInCbsFilter tests that instance_families in cbs_filter is correctly passed to DescribeDiskConfigQuota
func TestCvmInstanceTypes_Read_WithInstanceFamiliesInCbsFilter(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	cvmClient := &cvmSDK.Client{}
	cbsClient := &cbsSDK.Client{}

	patches.ApplyMethodReturn(newInstanceTypesMockMeta().client, "UseCvmClient", cvmClient)
	patches.ApplyMethodReturn(newInstanceTypesMockMeta().client, "UseCbsClient", cbsClient)

	// Mock DescribeZoneInstanceConfigInfos to return instance types
	patches.ApplyMethodFunc(cvmClient, "DescribeZoneInstanceConfigInfos", func(request *cvmSDK.DescribeZoneInstanceConfigInfosRequest) (*cvmSDK.DescribeZoneInstanceConfigInfosResponse, error) {
		resp := cvmSDK.NewDescribeZoneInstanceConfigInfosResponse()
		resp.Response = &cvmSDK.DescribeZoneInstanceConfigInfosResponseParams{
			InstanceTypeQuotaSet: []*cvmSDK.InstanceTypeQuotaItem{
				{
					Zone:               ptrString("ap-guangzhou-6"),
					InstanceType:       ptrString("S5.MEDIUM4"),
					Cpu:                ptrInt64(2),
					Memory:             ptrInt64(4),
					InstanceFamily:     ptrString("S5"),
					InstanceChargeType: ptrString("POSTPAID_BY_HOUR"),
					Status:             ptrString("SELL"),
				},
			},
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	// Mock DescribeDiskConfigQuota - verify that InstanceFamilies is set correctly
	var capturedInstanceFamilies []*string
	patches.ApplyMethodFunc(cbsClient, "DescribeDiskConfigQuota", func(request *cbsSDK.DescribeDiskConfigQuotaRequest) (*cbsSDK.DescribeDiskConfigQuotaResponse, error) {
		// Capture InstanceFamilies for verification
		capturedInstanceFamilies = request.InstanceFamilies

		resp := cbsSDK.NewDescribeDiskConfigQuotaResponse()
		resp.Response = &cbsSDK.DescribeDiskConfigQuotaResponseParams{
			DiskConfigSet: []*cbsSDK.DiskConfig{
				{
					Available:      ptrBool(true),
					DiskChargeType: ptrString("PREPAID"),
					Zone:           ptrString("ap-guangzhou-6"),
					InstanceFamily: ptrString("S5"),
					DiskType:       ptrString("CLOUD_SSD"),
					StepSize:       ptrUint64(10),
					DiskUsage:      ptrString("SYSTEM_DISK"),
					MinDiskSize:    ptrUint64(50),
					MaxDiskSize:    ptrUint64(500),
					DeviceClass:    ptrString("default"),
				},
			},
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newInstanceTypesMockMeta()
	res := cvm.DataSourceTencentCloudInstanceTypes()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"availability_zone": "ap-guangzhou-6",
		"cpu_core_count":    2,
		"memory_size":       4,
		"cbs_filter": []interface{}{
			map[string]interface{}{
				"disk_types":        []interface{}{"CLOUD_SSD"},
				"disk_charge_type":  "PREPAID",
				"disk_usage":        "SYSTEM_DISK",
				"instance_families": []interface{}{"S5", "M5"},
			},
		},
	})

	err := res.Read(d, meta)
	assert.NoError(t, err)

	// Verify InstanceFamilies was set with the user-provided values, not the derived single family
	assert.NotNil(t, capturedInstanceFamilies)
	assert.Equal(t, 2, len(capturedInstanceFamilies))
	assert.Equal(t, "S5", *capturedInstanceFamilies[0])
	assert.Equal(t, "M5", *capturedInstanceFamilies[1])
}

// TestCvmInstanceTypes_Read_WithoutInstanceFamiliesInCbsFilter tests that when instance_families is not provided, the single family from instance type is used
func TestCvmInstanceTypes_Read_WithoutInstanceFamiliesInCbsFilter(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	cvmClient := &cvmSDK.Client{}
	cbsClient := &cbsSDK.Client{}

	patches.ApplyMethodReturn(newInstanceTypesMockMeta().client, "UseCvmClient", cvmClient)
	patches.ApplyMethodReturn(newInstanceTypesMockMeta().client, "UseCbsClient", cbsClient)

	// Mock DescribeZoneInstanceConfigInfos to return instance types
	patches.ApplyMethodFunc(cvmClient, "DescribeZoneInstanceConfigInfos", func(request *cvmSDK.DescribeZoneInstanceConfigInfosRequest) (*cvmSDK.DescribeZoneInstanceConfigInfosResponse, error) {
		resp := cvmSDK.NewDescribeZoneInstanceConfigInfosResponse()
		resp.Response = &cvmSDK.DescribeZoneInstanceConfigInfosResponseParams{
			InstanceTypeQuotaSet: []*cvmSDK.InstanceTypeQuotaItem{
				{
					Zone:               ptrString("ap-guangzhou-6"),
					InstanceType:       ptrString("S5.MEDIUM4"),
					Cpu:                ptrInt64(2),
					Memory:             ptrInt64(4),
					InstanceFamily:     ptrString("S5"),
					InstanceChargeType: ptrString("POSTPAID_BY_HOUR"),
					Status:             ptrString("SELL"),
				},
			},
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	// Mock DescribeDiskConfigQuota - verify that InstanceFamilies uses the derived single family
	var capturedInstanceFamilies []*string
	patches.ApplyMethodFunc(cbsClient, "DescribeDiskConfigQuota", func(request *cbsSDK.DescribeDiskConfigQuotaRequest) (*cbsSDK.DescribeDiskConfigQuotaResponse, error) {
		// Capture InstanceFamilies for verification
		capturedInstanceFamilies = request.InstanceFamilies

		resp := cbsSDK.NewDescribeDiskConfigQuotaResponse()
		resp.Response = &cbsSDK.DescribeDiskConfigQuotaResponseParams{
			DiskConfigSet: []*cbsSDK.DiskConfig{
				{
					Available:      ptrBool(true),
					DiskChargeType: ptrString("PREPAID"),
					Zone:           ptrString("ap-guangzhou-6"),
					InstanceFamily: ptrString("S5"),
					DiskType:       ptrString("CLOUD_SSD"),
					StepSize:       ptrUint64(10),
					DiskUsage:      ptrString("SYSTEM_DISK"),
					MinDiskSize:    ptrUint64(50),
					MaxDiskSize:    ptrUint64(500),
					DeviceClass:    ptrString("default"),
				},
			},
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newInstanceTypesMockMeta()
	res := cvm.DataSourceTencentCloudInstanceTypes()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"availability_zone": "ap-guangzhou-6",
		"cpu_core_count":    2,
		"memory_size":       4,
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

	// Verify InstanceFamilies was set with the single family derived from instance type result
	assert.NotNil(t, capturedInstanceFamilies)
	assert.Equal(t, 1, len(capturedInstanceFamilies))
	assert.Equal(t, "S5", *capturedInstanceFamilies[0])
}

// TestCvmInstanceTypes_Schema_InstanceFamilies validates the instance_families schema definition in cbs_filter
func TestCvmInstanceTypes_Schema_InstanceFamilies(t *testing.T) {
	res := cvm.DataSourceTencentCloudInstanceTypes()

	assert.NotNil(t, res)
	assert.NotNil(t, res.Schema)

	// Check cbs_filter schema exists
	assert.Contains(t, res.Schema, "cbs_filter")
	cbsFilter := res.Schema["cbs_filter"]
	assert.Equal(t, schema.TypeList, cbsFilter.Type)
	assert.True(t, cbsFilter.Optional)

	// Check instance_families exists within cbs_filter
	cbsFilterElem := cbsFilter.Elem.(*schema.Resource)
	assert.Contains(t, cbsFilterElem.Schema, "instance_families")
	instanceFamilies := cbsFilterElem.Schema["instance_families"]
	assert.Equal(t, schema.TypeList, instanceFamilies.Type)
	assert.True(t, instanceFamilies.Optional)
	assert.NotNil(t, instanceFamilies.Elem)
}
