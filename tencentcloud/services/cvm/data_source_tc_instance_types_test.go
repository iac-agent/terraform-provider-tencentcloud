package cvm_test

import (
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	cbs "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cbs/v20170312"
	cvmSdk "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	svccbs "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/cbs"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/cvm"
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

func ptrInt64(i int64) *int64 {
	return &i
}

func ptrUint64(u uint64) *uint64 {
	return &u
}

func ptrBool(b bool) *bool {
	return &b
}

func ptrFloat64(f float64) *float64 {
	return &f
}

// TestInstanceTypesDataSource_Schema_NewParams validates schema definition for new parameters
func TestInstanceTypesDataSource_Schema_NewParams(t *testing.T) {
	res := cvm.DataSourceTencentCloudInstanceTypes()

	assert.NotNil(t, res)
	assert.Contains(t, res.Schema, "disk_types")
	assert.Contains(t, res.Schema, "zones")
	assert.Contains(t, res.Schema, "memory")

	// Check disk_types schema
	diskTypes := res.Schema["disk_types"]
	assert.Equal(t, schema.TypeList, diskTypes.Type)
	assert.True(t, diskTypes.Optional)
	assert.False(t, diskTypes.Required)
	assert.False(t, diskTypes.Computed)

	// Check zones schema
	zones := res.Schema["zones"]
	assert.Equal(t, schema.TypeList, zones.Type)
	assert.True(t, zones.Optional)
	assert.False(t, zones.Required)
	assert.False(t, zones.Computed)

	// Check memory schema
	memory := res.Schema["memory"]
	assert.Equal(t, schema.TypeInt, memory.Type)
	assert.True(t, memory.Optional)
	assert.False(t, memory.Required)
	assert.False(t, memory.Computed)
}

// TestInstanceTypesDataSource_Read_DiskTypesOverride tests that top-level disk_types overrides cbs_filter.disk_types
func TestInstanceTypesDataSource_Read_DiskTypesOverride(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	cvmClient := &cvmSdk.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseCvmClient", cvmClient)

	// Mock DescribeZoneInstanceConfigInfos to return instance types
	patches.ApplyMethodFunc(cvmClient, "DescribeZoneInstanceConfigInfos", func(request *cvmSdk.DescribeZoneInstanceConfigInfosRequest) (*cvmSdk.DescribeZoneInstanceConfigInfosResponse, error) {
		resp := cvmSdk.NewDescribeZoneInstanceConfigInfosResponse()
		resp.Response = &cvmSdk.DescribeZoneInstanceConfigInfosResponseParams{
			InstanceTypeQuotaSet: []*cvmSdk.InstanceTypeQuotaItem{
				{
					Zone:               ptrString("ap-guangzhou-3"),
					Cpu:                ptrInt64(4),
					Memory:             ptrInt64(8),
					InstanceFamily:     ptrString("S5"),
					InstanceType:       ptrString("S5.MEDIUM4"),
					Status:             ptrString("Sell"),
					InstanceChargeType: ptrString("POSTPAID_BY_HOUR"),
				},
			},
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	cbsClient := &cbs.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseCbsClient", cbsClient)

	// Mock DescribeDiskConfigQuota - verify that disk_types_override is used
	patches.ApplyMethodFunc(cbsClient, "DescribeDiskConfigQuota", func(request *cbs.DescribeDiskConfigQuotaRequest) (*cbs.DescribeDiskConfigQuotaResponse, error) {
		// Verify that DiskTypes comes from the override, not cbs_filter
		assert.NotNil(t, request.DiskTypes)
		diskTypes := make([]string, 0)
		for _, dt := range request.DiskTypes {
			diskTypes = append(diskTypes, *dt)
		}
		assert.Contains(t, diskTypes, "CLOUD_HSSD")
		assert.Contains(t, diskTypes, "CLOUD_SSD")

		resp := cbs.NewDescribeDiskConfigQuotaResponse()
		resp.Response = &cbs.DescribeDiskConfigQuotaResponseParams{
			DiskConfigSet: []*cbs.DiskConfig{
				{
					Available:      ptrBool(true),
					DiskChargeType: ptrString("POSTPAID_BY_HOUR"),
					Zone:           ptrString("ap-guangzhou-3"),
					InstanceFamily: ptrString("S5"),
					DiskType:       ptrString("CLOUD_SSD"),
					StepSize:       ptrUint64(10),
					DiskUsage:      ptrString("DATA_DISK"),
					MinDiskSize:    ptrUint64(20),
					MaxDiskSize:    ptrUint64(500),
				},
			},
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := cvm.DataSourceTencentCloudInstanceTypes()

	cbsFilter := []interface{}{
		map[string]interface{}{
			"disk_types":       []interface{}{"CLOUD_BASIC"},
			"disk_charge_type": "POSTPAID_BY_HOUR",
			"disk_usage":       "DATA_DISK",
		},
	}

	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"availability_zone": "ap-guangzhou-3",
		"cpu_core_count":    4,
		"memory_size":       8,
		"cbs_filter":        cbsFilter,
		"disk_types":        []interface{}{"CLOUD_SSD", "CLOUD_HSSD"},
	})

	err := res.Read(d, meta)
	assert.NoError(t, err)
}

// TestInstanceTypesDataSource_Read_ZonesOverride tests that top-level zones overrides instance type's zone
func TestInstanceTypesDataSource_Read_ZonesOverride(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	cvmClient := &cvmSdk.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseCvmClient", cvmClient)

	// Mock DescribeZoneInstanceConfigInfos to return instance types
	patches.ApplyMethodFunc(cvmClient, "DescribeZoneInstanceConfigInfos", func(request *cvmSdk.DescribeZoneInstanceConfigInfosRequest) (*cvmSdk.DescribeZoneInstanceConfigInfosResponse, error) {
		resp := cvmSdk.NewDescribeZoneInstanceConfigInfosResponse()
		resp.Response = &cvmSdk.DescribeZoneInstanceConfigInfosResponseParams{
			InstanceTypeQuotaSet: []*cvmSdk.InstanceTypeQuotaItem{
				{
					Zone:               ptrString("ap-guangzhou-3"),
					Cpu:                ptrInt64(4),
					Memory:             ptrInt64(8),
					InstanceFamily:     ptrString("S5"),
					InstanceType:       ptrString("S5.MEDIUM4"),
					Status:             ptrString("Sell"),
					InstanceChargeType: ptrString("POSTPAID_BY_HOUR"),
				},
			},
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	cbsClient := &cbs.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseCbsClient", cbsClient)

	// Mock DescribeDiskConfigQuota - verify that zones_override is used instead of instance type's zone
	patches.ApplyMethodFunc(cbsClient, "DescribeDiskConfigQuota", func(request *cbs.DescribeDiskConfigQuotaRequest) (*cbs.DescribeDiskConfigQuotaResponse, error) {
		// Verify that Zones comes from the override, not instance type's availability_zone
		assert.NotNil(t, request.Zones)
		zones := make([]string, 0)
		for _, z := range request.Zones {
			zones = append(zones, *z)
		}
		assert.Contains(t, zones, "ap-guangzhou-6")
		assert.Contains(t, zones, "ap-guangzhou-4")
		assert.NotContains(t, zones, "ap-guangzhou-3") // The instance type's zone should not be used

		resp := cbs.NewDescribeDiskConfigQuotaResponse()
		resp.Response = &cbs.DescribeDiskConfigQuotaResponseParams{
			DiskConfigSet: []*cbs.DiskConfig{
				{
					Available:      ptrBool(true),
					DiskChargeType: ptrString("POSTPAID_BY_HOUR"),
					Zone:           ptrString("ap-guangzhou-6"),
					InstanceFamily: ptrString("S5"),
					DiskType:       ptrString("CLOUD_SSD"),
					StepSize:       ptrUint64(10),
					DiskUsage:      ptrString("DATA_DISK"),
					MinDiskSize:    ptrUint64(20),
					MaxDiskSize:    ptrUint64(500),
				},
			},
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := cvm.DataSourceTencentCloudInstanceTypes()

	cbsFilter := []interface{}{
		map[string]interface{}{
			"disk_charge_type": "POSTPAID_BY_HOUR",
			"disk_usage":       "DATA_DISK",
		},
	}

	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"availability_zone": "ap-guangzhou-3",
		"cpu_core_count":    4,
		"memory_size":       8,
		"cbs_filter":        cbsFilter,
		"zones":             []interface{}{"ap-guangzhou-6", "ap-guangzhou-4"},
	})

	err := res.Read(d, meta)
	assert.NoError(t, err)
}

// TestInstanceTypesDataSource_Read_MemoryOverride tests that top-level memory overrides instance type's memory_size
func TestInstanceTypesDataSource_Read_MemoryOverride(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	cvmClient := &cvmSdk.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseCvmClient", cvmClient)

	// Mock DescribeZoneInstanceConfigInfos to return instance types
	patches.ApplyMethodFunc(cvmClient, "DescribeZoneInstanceConfigInfos", func(request *cvmSdk.DescribeZoneInstanceConfigInfosRequest) (*cvmSdk.DescribeZoneInstanceConfigInfosResponse, error) {
		resp := cvmSdk.NewDescribeZoneInstanceConfigInfosResponse()
		resp.Response = &cvmSdk.DescribeZoneInstanceConfigInfosResponseParams{
			InstanceTypeQuotaSet: []*cvmSdk.InstanceTypeQuotaItem{
				{
					Zone:               ptrString("ap-guangzhou-3"),
					Cpu:                ptrInt64(4),
					Memory:             ptrInt64(8), // instance type has 8GB memory
					InstanceFamily:     ptrString("S5"),
					InstanceType:       ptrString("S5.MEDIUM4"),
					Status:             ptrString("Sell"),
					InstanceChargeType: ptrString("POSTPAID_BY_HOUR"),
				},
			},
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	cbsClient := &cbs.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseCbsClient", cbsClient)

	// Mock DescribeDiskConfigQuota - verify that memory_override is used instead of instance type's memory_size
	patches.ApplyMethodFunc(cbsClient, "DescribeDiskConfigQuota", func(request *cbs.DescribeDiskConfigQuotaRequest) (*cbs.DescribeDiskConfigQuotaResponse, error) {
		// Verify that Memory comes from the override (16), not instance type's memory_size (8)
		assert.NotNil(t, request.Memory)
		assert.Equal(t, uint64(16), *request.Memory)

		resp := cbs.NewDescribeDiskConfigQuotaResponse()
		resp.Response = &cbs.DescribeDiskConfigQuotaResponseParams{
			DiskConfigSet: []*cbs.DiskConfig{
				{
					Available:      ptrBool(true),
					DiskChargeType: ptrString("POSTPAID_BY_HOUR"),
					Zone:           ptrString("ap-guangzhou-3"),
					InstanceFamily: ptrString("S5"),
					DiskType:       ptrString("CLOUD_SSD"),
					StepSize:       ptrUint64(10),
					DiskUsage:      ptrString("DATA_DISK"),
					MinDiskSize:    ptrUint64(20),
					MaxDiskSize:    ptrUint64(500),
				},
			},
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := cvm.DataSourceTencentCloudInstanceTypes()

	cbsFilter := []interface{}{
		map[string]interface{}{
			"disk_charge_type": "POSTPAID_BY_HOUR",
			"disk_usage":       "DATA_DISK",
		},
	}

	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"availability_zone": "ap-guangzhou-3",
		"cpu_core_count":    4,
		"memory_size":       8,
		"cbs_filter":        cbsFilter,
		"memory":            16,
	})

	err := res.Read(d, meta)
	assert.NoError(t, err)
}

// TestInstanceTypesDataSource_Read_BackwardCompatibility tests that when new params are not provided, behavior is unchanged
func TestInstanceTypesDataSource_Read_BackwardCompatibility(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	cvmClient := &cvmSdk.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseCvmClient", cvmClient)

	// Mock DescribeZoneInstanceConfigInfos to return instance types
	patches.ApplyMethodFunc(cvmClient, "DescribeZoneInstanceConfigInfos", func(request *cvmSdk.DescribeZoneInstanceConfigInfosRequest) (*cvmSdk.DescribeZoneInstanceConfigInfosResponse, error) {
		resp := cvmSdk.NewDescribeZoneInstanceConfigInfosResponse()
		resp.Response = &cvmSdk.DescribeZoneInstanceConfigInfosResponseParams{
			InstanceTypeQuotaSet: []*cvmSdk.InstanceTypeQuotaItem{
				{
					Zone:               ptrString("ap-guangzhou-3"),
					Cpu:                ptrInt64(4),
					Memory:             ptrInt64(8),
					InstanceFamily:     ptrString("S5"),
					InstanceType:       ptrString("S5.MEDIUM4"),
					Status:             ptrString("Sell"),
					InstanceChargeType: ptrString("POSTPAID_BY_HOUR"),
				},
			},
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	cbsClient := &cbs.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseCbsClient", cbsClient)

	// Mock DescribeDiskConfigQuota - verify that default behavior (no overrides) works correctly
	patches.ApplyMethodFunc(cbsClient, "DescribeDiskConfigQuota", func(request *cbs.DescribeDiskConfigQuotaRequest) (*cbs.DescribeDiskConfigQuotaResponse, error) {
		// Verify Zones uses availability_zone from instance type
		assert.NotNil(t, request.Zones)
		assert.Equal(t, "ap-guangzhou-3", *request.Zones[0])
		assert.Len(t, request.Zones, 1)

		// Verify Memory uses memory_size from instance type
		assert.NotNil(t, request.Memory)
		assert.Equal(t, uint64(8), *request.Memory)

		// Verify DiskTypes uses cbs_filter.disk_types
		assert.NotNil(t, request.DiskTypes)
		assert.Equal(t, "CLOUD_SSD", *request.DiskTypes[0])

		resp := cbs.NewDescribeDiskConfigQuotaResponse()
		resp.Response = &cbs.DescribeDiskConfigQuotaResponseParams{
			DiskConfigSet: []*cbs.DiskConfig{
				{
					Available:      ptrBool(true),
					DiskChargeType: ptrString("PREPAID"),
					Zone:           ptrString("ap-guangzhou-3"),
					InstanceFamily: ptrString("S5"),
					DiskType:       ptrString("CLOUD_SSD"),
					StepSize:       ptrUint64(10),
					DiskUsage:      ptrString("SYSTEM_DISK"),
					MinDiskSize:    ptrUint64(50),
					MaxDiskSize:    ptrUint64(500),
				},
			},
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := cvm.DataSourceTencentCloudInstanceTypes()

	cbsFilter := []interface{}{
		map[string]interface{}{
			"disk_types":       []interface{}{"CLOUD_SSD"},
			"disk_charge_type": "PREPAID",
			"disk_usage":       "SYSTEM_DISK",
		},
	}

	// No override parameters - just cbs_filter with availability_zone
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"availability_zone": "ap-guangzhou-3",
		"cpu_core_count":    4,
		"memory_size":       8,
		"cbs_filter":        cbsFilter,
	})

	err := res.Read(d, meta)
	assert.NoError(t, err)

	// Verify instance_types were set
	instanceTypes := d.Get("instance_types").([]interface{})
	assert.Len(t, instanceTypes, 1)

	instanceTypeMap := instanceTypes[0].(map[string]interface{})
	cbsConfigs := instanceTypeMap["cbs_configs"].([]interface{})
	assert.Len(t, cbsConfigs, 1)

	cbsConfigMap := cbsConfigs[0].(map[string]interface{})
	assert.Equal(t, true, cbsConfigMap["available"])
	assert.Equal(t, "PREPAID", cbsConfigMap["disk_charge_type"])
	assert.Equal(t, "CLOUD_SSD", cbsConfigMap["disk_type"])
}

// TestDescribeDiskConfigQuota_DiskTypesOverride tests CBS service handling of disk_types_override
func TestDescribeDiskConfigQuota_DiskTypesOverride(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	cbsClient := &cbs.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseCbsClient", cbsClient)

	patches.ApplyMethodFunc(cbsClient, "DescribeDiskConfigQuota", func(request *cbs.DescribeDiskConfigQuotaRequest) (*cbs.DescribeDiskConfigQuotaResponse, error) {
		// Verify that DiskTypes comes from override
		assert.NotNil(t, request.DiskTypes)
		diskTypes := make([]string, 0)
		for _, dt := range request.DiskTypes {
			diskTypes = append(diskTypes, *dt)
		}
		assert.Contains(t, diskTypes, "CLOUD_HSSD")

		resp := cbs.NewDescribeDiskConfigQuotaResponse()
		resp.Response = &cbs.DescribeDiskConfigQuotaResponseParams{
			DiskConfigSet: []*cbs.DiskConfig{},
			RequestId:     ptrString("fake-request-id"),
		}
		return resp, nil
	})

	client := newMockMeta().client
	cbsService := svccbs.NewCbsService(client)

	cvmInfo := map[string]interface{}{
		"availability_zone":   "ap-guangzhou-3",
		"cpu_core_count":      int64(4),
		"memory_size":         int64(8),
		"family":              "S5",
		"disk_types_override": []string{"CLOUD_HSSD"},
		"disk_charge_type":    "PREPAID",
		"disk_usage":          "SYSTEM_DISK",
	}

	diskConfigSet, err := cbsService.DescribeDiskConfigQuota(nil, cvmInfo)
	assert.NoError(t, err)
	assert.NotNil(t, diskConfigSet)
}

// TestDescribeDiskConfigQuota_ZonesOverride tests CBS service handling of zones_override
func TestDescribeDiskConfigQuota_ZonesOverride(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	cbsClient := &cbs.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseCbsClient", cbsClient)

	patches.ApplyMethodFunc(cbsClient, "DescribeDiskConfigQuota", func(request *cbs.DescribeDiskConfigQuotaRequest) (*cbs.DescribeDiskConfigQuotaResponse, error) {
		// Verify that Zones comes from override (multiple zones)
		assert.NotNil(t, request.Zones)
		zones := make([]string, 0)
		for _, z := range request.Zones {
			zones = append(zones, *z)
		}
		assert.Len(t, zones, 2)
		assert.Contains(t, zones, "ap-guangzhou-3")
		assert.Contains(t, zones, "ap-guangzhou-6")

		resp := cbs.NewDescribeDiskConfigQuotaResponse()
		resp.Response = &cbs.DescribeDiskConfigQuotaResponseParams{
			DiskConfigSet: []*cbs.DiskConfig{},
			RequestId:     ptrString("fake-request-id"),
		}
		return resp, nil
	})

	client := newMockMeta().client
	cbsService := svccbs.NewCbsService(client)

	cvmInfo := map[string]interface{}{
		"availability_zone": "ap-guangzhou-3", // This should be ignored when zones_override is present
		"cpu_core_count":    int64(4),
		"memory_size":       int64(8),
		"family":            "S5",
		"zones_override":    []string{"ap-guangzhou-3", "ap-guangzhou-6"},
		"disk_types":        []string{"CLOUD_SSD"},
		"disk_charge_type":  "PREPAID",
		"disk_usage":        "SYSTEM_DISK",
	}

	diskConfigSet, err := cbsService.DescribeDiskConfigQuota(nil, cvmInfo)
	assert.NoError(t, err)
	assert.NotNil(t, diskConfigSet)
}

// TestDescribeDiskConfigQuota_MemoryOverride tests CBS service handling of memory_override
func TestDescribeDiskConfigQuota_MemoryOverride(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	cbsClient := &cbs.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseCbsClient", cbsClient)

	patches.ApplyMethodFunc(cbsClient, "DescribeDiskConfigQuota", func(request *cbs.DescribeDiskConfigQuotaRequest) (*cbs.DescribeDiskConfigQuotaResponse, error) {
		// Verify that Memory comes from override
		assert.NotNil(t, request.Memory)
		assert.Equal(t, uint64(16), *request.Memory)

		resp := cbs.NewDescribeDiskConfigQuotaResponse()
		resp.Response = &cbs.DescribeDiskConfigQuotaResponseParams{
			DiskConfigSet: []*cbs.DiskConfig{},
			RequestId:     ptrString("fake-request-id"),
		}
		return resp, nil
	})

	client := newMockMeta().client
	cbsService := svccbs.NewCbsService(client)

	cvmInfo := map[string]interface{}{
		"availability_zone": "ap-guangzhou-3",
		"cpu_core_count":    int64(4),
		"memory_size":       int64(8), // This should be ignored when memory_override is present
		"family":            "S5",
		"memory_override":   int64(16),
		"disk_types":        []string{"CLOUD_SSD"},
		"disk_charge_type":  "PREPAID",
		"disk_usage":        "SYSTEM_DISK",
	}

	diskConfigSet, err := cbsService.DescribeDiskConfigQuota(nil, cvmInfo)
	assert.NoError(t, err)
	assert.NotNil(t, diskConfigSet)
}

// TestDescribeDiskConfigQuota_NoOverride tests CBS service backward compatibility without overrides
func TestDescribeDiskConfigQuota_NoOverride(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	cbsClient := &cbs.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseCbsClient", cbsClient)

	patches.ApplyMethodFunc(cbsClient, "DescribeDiskConfigQuota", func(request *cbs.DescribeDiskConfigQuotaRequest) (*cbs.DescribeDiskConfigQuotaResponse, error) {
		// Verify default behavior without overrides
		assert.NotNil(t, request.Zones)
		assert.Len(t, request.Zones, 1)
		assert.Equal(t, "ap-guangzhou-3", *request.Zones[0])

		assert.NotNil(t, request.Memory)
		assert.Equal(t, uint64(8), *request.Memory)

		assert.NotNil(t, request.DiskTypes)
		assert.Equal(t, "CLOUD_SSD", *request.DiskTypes[0])

		resp := cbs.NewDescribeDiskConfigQuotaResponse()
		resp.Response = &cbs.DescribeDiskConfigQuotaResponseParams{
			DiskConfigSet: []*cbs.DiskConfig{
				{
					Available:      ptrBool(true),
					DiskChargeType: ptrString("PREPAID"),
					Zone:           ptrString("ap-guangzhou-3"),
					InstanceFamily: ptrString("S5"),
					DiskType:       ptrString("CLOUD_SSD"),
					StepSize:       ptrUint64(10),
					DiskUsage:      ptrString("SYSTEM_DISK"),
					MinDiskSize:    ptrUint64(50),
					MaxDiskSize:    ptrUint64(500),
				},
			},
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	client := newMockMeta().client
	cbsService := svccbs.NewCbsService(client)

	cvmInfo := map[string]interface{}{
		"availability_zone": "ap-guangzhou-3",
		"cpu_core_count":    int64(4),
		"memory_size":       int64(8),
		"family":            "S5",
		"disk_types":        []string{"CLOUD_SSD"},
		"disk_charge_type":  "PREPAID",
		"disk_usage":        "SYSTEM_DISK",
	}

	diskConfigSet, err := cbsService.DescribeDiskConfigQuota(nil, cvmInfo)
	assert.NoError(t, err)
	assert.Len(t, diskConfigSet, 1)
}
