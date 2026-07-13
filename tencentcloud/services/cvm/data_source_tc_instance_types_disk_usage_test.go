package cvm_test

import (
	"context"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"

	cbsSDK "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cbs/v20170312"
	cvmSDK "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	cbsSvc "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/cbs"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/cvm"
)

// instanceTypesDiskUsageMockMeta implements tccommon.ProviderMeta
type instanceTypesDiskUsageMockMeta struct {
	client *connectivity.TencentCloudClient
}

func (m *instanceTypesDiskUsageMockMeta) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &instanceTypesDiskUsageMockMeta{}

func newInstanceTypesDiskUsageMockMeta() *instanceTypesDiskUsageMockMeta {
	return &instanceTypesDiskUsageMockMeta{client: &connectivity.TencentCloudClient{}}
}

// go test ./tencentcloud/services/cvm/ -run "TestInstanceTypesDiskUsage" -v -count=1 -gcflags="all=-l"

// Helper functions for creating pointers
func ptrString(v string) *string { return &v }
func ptrInt64(v int64) *int64    { return &v }
func ptrBool(v bool) *bool       { return &v }
func ptrUint64(v uint64) *uint64 { return &v }

// TestInstanceTypesDiskUsage_WithCbsFilter tests that disk_usage is correctly populated at the top level when cbs_filter with disk_usage is provided
func TestInstanceTypesDiskUsage_WithCbsFilter(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	// Mock DescribeInstancesSellTypeByFilter to return instance type data
	patches.ApplyMethodFunc(&cvm.CvmService{}, "DescribeInstancesSellTypeByFilter", func(ctx context.Context, filters map[string][]string) ([]*cvmSDK.InstanceTypeQuotaItem, error) {
		return []*cvmSDK.InstanceTypeQuotaItem{
			{
				Zone:               ptrString("ap-guangzhou-6"),
				InstanceType:       ptrString("S6.MEDIUM2"),
				Cpu:                ptrInt64(2),
				Memory:             ptrInt64(2),
				InstanceFamily:     ptrString("S6"),
				Status:             ptrString("SELL"),
				InstanceChargeType: ptrString("POSTPAID_BY_HOUR"),
				NetworkCard:        ptrInt64(0),
				TypeName:           ptrString("Standard S6.MEDIUM2"),
			},
		}, nil
	})

	// Mock DescribeDiskConfigQuota to return disk configs with DiskUsage
	patches.ApplyMethodFunc(&cbsSvc.CbsService{}, "DescribeDiskConfigQuota", func(ctx context.Context, cvmInfo map[string]interface{}) ([]*cbsSDK.DiskConfig, error) {
		return []*cbsSDK.DiskConfig{
			{
				Available:      ptrBool(true),
				DiskChargeType: ptrString("PREPAID"),
				Zone:           ptrString("ap-guangzhou-6"),
				InstanceFamily: ptrString("S6"),
				DiskType:       ptrString("CLOUD_SSD"),
				StepSize:       ptrUint64(10),
				DeviceClass:    ptrString("S6"),
				DiskUsage:      ptrString("SYSTEM_DISK"),
				MinDiskSize:    ptrUint64(50),
				MaxDiskSize:    ptrUint64(500),
			},
		}, nil
	})

	meta := newInstanceTypesDiskUsageMockMeta()
	res := cvm.DataSourceTencentCloudInstanceTypes()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"cpu_core_count":   2,
		"memory_size":      2,
		"exclude_sold_out": false,
		"filter": []interface{}{
			map[string]interface{}{
				"name":   "instance-family",
				"values": []interface{}{"S6"},
			},
			map[string]interface{}{
				"name":   "zone",
				"values": []interface{}{"ap-guangzhou-6"},
			},
		},
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

	instanceTypes := d.Get("instance_types").([]interface{})
	assert.GreaterOrEqual(t, len(instanceTypes), 1)

	firstInstance := instanceTypes[0].(map[string]interface{})
	// Check top-level disk_usage is populated from cbs_filter input
	assert.Equal(t, "SYSTEM_DISK", firstInstance["disk_usage"])

	// Check cbs_configs also has disk_usage
	cbsConfigs := firstInstance["cbs_configs"].([]interface{})
	assert.GreaterOrEqual(t, len(cbsConfigs), 1)
	firstCbsConfig := cbsConfigs[0].(map[string]interface{})
	assert.Equal(t, "SYSTEM_DISK", firstCbsConfig["disk_usage"])
}

// TestInstanceTypesDiskUsage_WithCbsFilterNoDiskUsageInput tests that disk_usage is populated from the API response when cbs_filter is provided but without disk_usage input
func TestInstanceTypesDiskUsage_WithCbsFilterNoDiskUsageInput(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	// Mock DescribeInstancesSellTypeByFilter to return instance type data
	patches.ApplyMethodFunc(&cvm.CvmService{}, "DescribeInstancesSellTypeByFilter", func(ctx context.Context, filters map[string][]string) ([]*cvmSDK.InstanceTypeQuotaItem, error) {
		return []*cvmSDK.InstanceTypeQuotaItem{
			{
				Zone:               ptrString("ap-guangzhou-6"),
				InstanceType:       ptrString("S6.MEDIUM2"),
				Cpu:                ptrInt64(2),
				Memory:             ptrInt64(2),
				InstanceFamily:     ptrString("S6"),
				Status:             ptrString("SELL"),
				InstanceChargeType: ptrString("POSTPAID_BY_HOUR"),
			},
		}, nil
	})

	// Mock DescribeDiskConfigQuota to return disk configs with DiskUsage from API response
	patches.ApplyMethodFunc(&cbsSvc.CbsService{}, "DescribeDiskConfigQuota", func(ctx context.Context, cvmInfo map[string]interface{}) ([]*cbsSDK.DiskConfig, error) {
		return []*cbsSDK.DiskConfig{
			{
				Available:      ptrBool(true),
				DiskChargeType: ptrString("POSTPAID_BY_HOUR"),
				Zone:           ptrString("ap-guangzhou-6"),
				InstanceFamily: ptrString("S6"),
				DiskType:       ptrString("CLOUD_SSD"),
				DiskUsage:      ptrString("DATA_DISK"),
				MinDiskSize:    ptrUint64(50),
				MaxDiskSize:    ptrUint64(500),
			},
		}, nil
	})

	meta := newInstanceTypesDiskUsageMockMeta()
	res := cvm.DataSourceTencentCloudInstanceTypes()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"cpu_core_count":   2,
		"memory_size":      2,
		"exclude_sold_out": false,
		"filter": []interface{}{
			map[string]interface{}{
				"name":   "instance-family",
				"values": []interface{}{"S6"},
			},
			map[string]interface{}{
				"name":   "zone",
				"values": []interface{}{"ap-guangzhou-6"},
			},
		},
		"cbs_filter": []interface{}{
			map[string]interface{}{
				"disk_types":       []interface{}{"CLOUD_SSD"},
				"disk_charge_type": "POSTPAID_BY_HOUR",
			},
		},
	})

	err := res.Read(d, meta)
	assert.NoError(t, err)

	instanceTypes := d.Get("instance_types").([]interface{})
	assert.GreaterOrEqual(t, len(instanceTypes), 1)

	firstInstance := instanceTypes[0].(map[string]interface{})
	// Check top-level disk_usage is populated from API response (DATA_DISK)
	assert.Equal(t, "DATA_DISK", firstInstance["disk_usage"])
}

// TestInstanceTypesDiskUsage_WithoutCbsFilter tests that disk_usage is empty/null when cbs_filter is not provided
func TestInstanceTypesDiskUsage_WithoutCbsFilter(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	// Mock DescribeInstancesSellTypeByFilter to return instance type data
	patches.ApplyMethodFunc(&cvm.CvmService{}, "DescribeInstancesSellTypeByFilter", func(ctx context.Context, filters map[string][]string) ([]*cvmSDK.InstanceTypeQuotaItem, error) {
		return []*cvmSDK.InstanceTypeQuotaItem{
			{
				Zone:               ptrString("ap-guangzhou-6"),
				InstanceType:       ptrString("S6.MEDIUM2"),
				Cpu:                ptrInt64(2),
				Memory:             ptrInt64(2),
				InstanceFamily:     ptrString("S6"),
				Status:             ptrString("SELL"),
				InstanceChargeType: ptrString("POSTPAID_BY_HOUR"),
			},
		}, nil
	})

	meta := newInstanceTypesDiskUsageMockMeta()
	res := cvm.DataSourceTencentCloudInstanceTypes()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"availability_zone": "ap-guangzhou-6",
		"cpu_core_count":    2,
		"memory_size":       2,
		"exclude_sold_out":  false,
	})

	err := res.Read(d, meta)
	assert.NoError(t, err)

	instanceTypes := d.Get("instance_types").([]interface{})
	assert.GreaterOrEqual(t, len(instanceTypes), 1)

	firstInstance := instanceTypes[0].(map[string]interface{})
	// Check top-level disk_usage is empty/null when cbs_filter is not provided
	assert.Equal(t, "", firstInstance["disk_usage"])
}

// TestInstanceTypesDiskUsage_Schema tests that disk_usage is defined in the schema
func TestInstanceTypesDiskUsage_Schema(t *testing.T) {
	res := cvm.DataSourceTencentCloudInstanceTypes()

	assert.NotNil(t, res)

	// Check instance_types schema
	instanceTypesSchema := res.Schema["instance_types"]
	assert.NotNil(t, instanceTypesSchema)
	assert.Equal(t, schema.TypeList, instanceTypesSchema.Type)
	assert.True(t, instanceTypesSchema.Computed)

	// Check disk_usage is inside instance_types element schema
	elemResource := instanceTypesSchema.Elem.(*schema.Resource)
	assert.Contains(t, elemResource.Schema, "disk_usage")

	diskUsageSchema := elemResource.Schema["disk_usage"]
	assert.Equal(t, schema.TypeString, diskUsageSchema.Type)
	assert.True(t, diskUsageSchema.Computed)
	assert.Contains(t, diskUsageSchema.Description, "SYSTEM_DISK")
	assert.Contains(t, diskUsageSchema.Description, "DATA_DISK")
}
