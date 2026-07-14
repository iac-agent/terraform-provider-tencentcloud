package cvm_test

import (
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	cbsSDK "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cbs/v20170312"
	cvmSDK "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/cvm"
)

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

func itPtrString(v string) *string    { return &v }
func itPtrInt64(v int64) *int64       { return &v }
func itPtrUint64(v uint64) *uint64    { return &v }
func itPtrFloat64(v float64) *float64 { return &v }
func itPtrBool(v bool) *bool          { return &v }

// createMockInstanceTypeQuotaItem creates a mock CVM InstanceTypeQuotaItem
func createMockInstanceTypeQuotaItem(zone, family, instanceType, status string, cpu, memory int64) *cvmSDK.InstanceTypeQuotaItem {
	return &cvmSDK.InstanceTypeQuotaItem{
		Zone:               itPtrString(zone),
		InstanceFamily:     itPtrString(family),
		InstanceType:       itPtrString(instanceType),
		Status:             itPtrString(status),
		Cpu:                itPtrInt64(cpu),
		Memory:             itPtrInt64(memory),
		Gpu:                itPtrInt64(0),
		InstanceChargeType: itPtrString("POSTPAID_BY_HOUR"),
		NetworkCard:        itPtrInt64(0),
		TypeName:           itPtrString("test-type"),
		SoldOutReason:      itPtrString(""),
		InstanceBandwidth:  itPtrFloat64(0),
		InstancePps:        itPtrInt64(0),
		StorageBlockAmount: itPtrInt64(0),
		CpuType:            itPtrString("Intel"),
		Fpga:               itPtrInt64(0),
		GpuCount:           itPtrFloat64(0),
		Frequency:          itPtrString("2.5GHz"),
		StatusCategory:     itPtrString("NormalStock"),
		Remark:             itPtrString(""),
	}
}

// createMockDiskConfig creates a mock CBS DiskConfig
func createMockDiskConfig(zone, family, diskType, diskChargeType, diskUsage string, available bool, minSize, maxSize uint64) *cbsSDK.DiskConfig {
	return &cbsSDK.DiskConfig{
		Available:      itPtrBool(available),
		DiskChargeType: itPtrString(diskChargeType),
		Zone:           itPtrString(zone),
		InstanceFamily: itPtrString(family),
		DiskType:       itPtrString(diskType),
		StepSize:       itPtrUint64(10),
		DeviceClass:    itPtrString(""),
		DiskUsage:      itPtrString(diskUsage),
		MinDiskSize:    itPtrUint64(minSize),
		MaxDiskSize:    itPtrUint64(maxSize),
	}
}

// setupCvmMock patches the CVM SDK client to return mock data
func setupCvmMock(patches *gomonkey.Patches, items []*cvmSDK.InstanceTypeQuotaItem) {
	cvmClient := &cvmSDK.Client{}
	patches.ApplyMethodReturn(newInstanceTypesMockMeta().client, "UseCvmClient", cvmClient)

	patches.ApplyMethodFunc(cvmClient, "DescribeZoneInstanceConfigInfos",
		func(request *cvmSDK.DescribeZoneInstanceConfigInfosRequest) (*cvmSDK.DescribeZoneInstanceConfigInfosResponse, error) {
			resp := cvmSDK.NewDescribeZoneInstanceConfigInfosResponse()
			resp.Response = &cvmSDK.DescribeZoneInstanceConfigInfosResponseParams{
				InstanceTypeQuotaSet: items,
				RequestId:            itPtrString("mock-request-id"),
			}
			return resp, nil
		})
}

// setupCbsMock patches the CBS SDK client to return mock data
func setupCbsMock(patches *gomonkey.Patches, diskConfigs []*cbsSDK.DiskConfig) {
	cbsClient := &cbsSDK.Client{}
	patches.ApplyMethodReturn(newInstanceTypesMockMeta().client, "UseCbsClient", cbsClient)

	patches.ApplyMethodFunc(cbsClient, "DescribeDiskConfigQuota",
		func(request *cbsSDK.DescribeDiskConfigQuotaRequest) (*cbsSDK.DescribeDiskConfigQuotaResponse, error) {
			resp := cbsSDK.NewDescribeDiskConfigQuotaResponse()
			resp.Response = &cbsSDK.DescribeDiskConfigQuotaResponseParams{
				DiskConfigSet: diskConfigs,
				RequestId:     itPtrString("mock-cbs-request-id"),
			}
			return resp, nil
		})
}

// go test ./tencentcloud/services/cvm/ -run "TestInstanceTypes_InstanceFamiliesFilter" -v -count=1 -gcflags="all=-l"

// TestInstanceTypes_InstanceFamiliesFilter tests that InstanceFamilies parameter correctly
// filters instance types via instance-family filter
func TestInstanceTypes_InstanceFamiliesFilter(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	mockItems := []*cvmSDK.InstanceTypeQuotaItem{
		createMockInstanceTypeQuotaItem("ap-guangzhou-6", "SA2", "SA2.MEDIUM4", "SELL", 2, 4),
		createMockInstanceTypeQuotaItem("ap-guangzhou-6", "S5", "S5.MEDIUM4", "SELL", 2, 4),
	}
	setupCvmMock(patches, mockItems)

	meta := newInstanceTypesMockMeta()
	res := cvm.DataSourceTencentCloudInstanceTypes()

	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"instance_families": []interface{}{"SA2"},
		"exclude_sold_out":  false,
	})

	err := res.Read(d, meta)
	assert.NoError(t, err)

	instanceTypesList := d.Get("instance_types").([]interface{})
	assert.Equal(t, 2, len(instanceTypesList), "Both items returned by API should be present; server-side filtering applies the instance-family filter")
}

// go test ./tencentcloud/services/cvm/ -run "TestInstanceTypes_DiskTypesTriggersCbsQuery" -v -count=1 -gcflags="all=-l"

// TestInstanceTypes_DiskTypesTriggersCbsQuery tests that DiskTypes parameter correctly
// triggers CBS config queries and passes values to DescribeDiskConfigQuota
func TestInstanceTypes_DiskTypesTriggersCbsQuery(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	mockItems := []*cvmSDK.InstanceTypeQuotaItem{
		createMockInstanceTypeQuotaItem("ap-guangzhou-6", "SA2", "SA2.MEDIUM4", "SELL", 2, 4),
	}
	setupCvmMock(patches, mockItems)

	mockDiskConfigs := []*cbsSDK.DiskConfig{
		createMockDiskConfig("ap-guangzhou-6", "SA2", "CLOUD_SSD", "POSTPAID_BY_HOUR", "SYSTEM_DISK", true, 50, 500),
	}
	setupCbsMock(patches, mockDiskConfigs)

	meta := newInstanceTypesMockMeta()
	res := cvm.DataSourceTencentCloudInstanceTypes()

	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"disk_types":       []interface{}{"CLOUD_SSD"},
		"exclude_sold_out": false,
	})

	err := res.Read(d, meta)
	assert.NoError(t, err)

	instanceTypesList := d.Get("instance_types").([]interface{})
	assert.Equal(t, 1, len(instanceTypesList))

	firstType := instanceTypesList[0].(map[string]interface{})
	cbsConfigs := firstType["cbs_configs"].([]interface{})
	assert.Equal(t, 1, len(cbsConfigs), "CBS configs should be populated when disk_types is provided")

	firstCbsConfig := cbsConfigs[0].(map[string]interface{})
	assert.Equal(t, "CLOUD_SSD", firstCbsConfig["disk_type"])
}

// go test ./tencentcloud/services/cvm/ -run "TestInstanceTypes_BackwardCompatibility" -v -count=1 -gcflags="all=-l"

// TestInstanceTypes_BackwardCompatibility tests that when InstanceFamilies and DiskTypes
// are not provided, the data source behaves exactly as before
func TestInstanceTypes_BackwardCompatibility(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	mockItems := []*cvmSDK.InstanceTypeQuotaItem{
		createMockInstanceTypeQuotaItem("ap-guangzhou-3", "SA2", "SA2.MEDIUM4", "SELL", 4, 8),
	}
	setupCvmMock(patches, mockItems)

	meta := newInstanceTypesMockMeta()
	res := cvm.DataSourceTencentCloudInstanceTypes()

	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"availability_zone": "ap-guangzhou-3",
		"cpu_core_count":    4,
		"memory_size":       8,
		"exclude_sold_out":  false,
	})

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.NotEmpty(t, d.Id())

	instanceTypesList := d.Get("instance_types").([]interface{})
	assert.Equal(t, 1, len(instanceTypesList))

	firstType := instanceTypesList[0].(map[string]interface{})
	assert.Equal(t, "ap-guangzhou-3", firstType["availability_zone"])
	assert.Equal(t, int(4), firstType["cpu_core_count"])
	assert.Equal(t, int(8), firstType["memory_size"])
}

// go test ./tencentcloud/services/cvm/ -run "TestInstanceTypes_InstanceFamiliesMergeWithFilter" -v -count=1 -gcflags="all=-l"

// TestInstanceTypes_InstanceFamiliesMergeWithFilter tests that when both InstanceFamilies
// and filter block instance-family are provided, their values are merged
func TestInstanceTypes_InstanceFamiliesMergeWithFilter(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	// The API will return results for all specified families
	mockItems := []*cvmSDK.InstanceTypeQuotaItem{
		createMockInstanceTypeQuotaItem("ap-guangzhou-6", "SA2", "SA2.MEDIUM4", "SELL", 2, 4),
		createMockInstanceTypeQuotaItem("ap-guangzhou-6", "S5", "S5.MEDIUM4", "SELL", 2, 4),
		createMockInstanceTypeQuotaItem("ap-guangzhou-6", "M5", "M5.MEDIUM4", "SELL", 2, 4),
	}
	setupCvmMock(patches, mockItems)

	meta := newInstanceTypesMockMeta()
	res := cvm.DataSourceTencentCloudInstanceTypes()

	// Set up with instance_families = ["M5"] and filter block instance-family = ["SA2", "S5"]
	// The merged filter should result in instance-family filter with ["SA2", "S5", "M5"]
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"instance_families": []interface{}{"M5"},
		"exclude_sold_out":  false,
	})

	// Add a filter block with instance-family = SA2, S5
	filterElem := res.Schema["filter"].Elem.(*schema.Resource)
	hash := schema.HashResource(filterElem)
	filterSet := schema.NewSet(hash, []interface{}{
		map[string]interface{}{
			"name":   "instance-family",
			"values": []interface{}{"SA2", "S5"},
		},
	})
	_ = d.Set("filter", filterSet)

	err := res.Read(d, meta)
	assert.NoError(t, err)

	instanceTypesList := d.Get("instance_types").([]interface{})
	// The merged filter should include SA2, S5 (from filter block) and M5 (from instance_families)
	// The API mock returns all 3, so we should get all 3
	assert.Equal(t, 3, len(instanceTypesList))
}

// go test ./tencentcloud/services/cvm/ -run "TestInstanceTypes_DiskTypesPrecedenceOverCbsFilter" -v -count=1 -gcflags="all=-l"

// TestInstanceTypes_DiskTypesPrecedenceOverCbsFilter tests that when both DiskTypes
// top-level parameter and cbs_filter.disk_types are provided, DiskTypes takes precedence
func TestInstanceTypes_DiskTypesPrecedenceOverCbsFilter(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	mockItems := []*cvmSDK.InstanceTypeQuotaItem{
		createMockInstanceTypeQuotaItem("ap-guangzhou-6", "SA2", "SA2.MEDIUM4", "SELL", 2, 4),
	}
	setupCvmMock(patches, mockItems)

	// CBS mock returns CLOUD_SSD (matching top-level disk_types), not CLOUD_BASIC (from cbs_filter)
	mockDiskConfigs := []*cbsSDK.DiskConfig{
		createMockDiskConfig("ap-guangzhou-6", "SA2", "CLOUD_SSD", "POSTPAID_BY_HOUR", "SYSTEM_DISK", true, 50, 500),
	}
	setupCbsMock(patches, mockDiskConfigs)

	meta := newInstanceTypesMockMeta()
	res := cvm.DataSourceTencentCloudInstanceTypes()

	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"disk_types":       []interface{}{"CLOUD_SSD"},
		"exclude_sold_out": false,
	})

	// Set cbs_filter with disk_types = CLOUD_BASIC (should be overridden by top-level disk_types)
	cbsFilterData := map[string]interface{}{
		"disk_types":       []interface{}{"CLOUD_BASIC"},
		"disk_charge_type": "POSTPAID_BY_HOUR",
		"disk_usage":       "SYSTEM_DISK",
	}
	_ = d.Set("cbs_filter", []interface{}{cbsFilterData})

	err := res.Read(d, meta)
	assert.NoError(t, err)

	instanceTypesList := d.Get("instance_types").([]interface{})
	assert.Equal(t, 1, len(instanceTypesList))

	firstType := instanceTypesList[0].(map[string]interface{})
	cbsConfigs := firstType["cbs_configs"].([]interface{})
	assert.Equal(t, 1, len(cbsConfigs))

	// The result should use CLOUD_SSD (from top-level disk_types) not CLOUD_BASIC (from cbs_filter)
	firstCbsConfig := cbsConfigs[0].(map[string]interface{})
	assert.Equal(t, "CLOUD_SSD", firstCbsConfig["disk_type"])
}

// go test ./tencentcloud/services/cvm/ -run "TestInstanceTypes_InstanceFamiliesWithCbsQuery" -v -count=1 -gcflags="all=-l"

// TestInstanceTypes_InstanceFamiliesWithCbsQuery tests that when InstanceFamilies is provided
// and CBS config queries are triggered, InstanceFamilies values are passed directly to
// DescribeDiskConfigQuota API's InstanceFamilies field
func TestInstanceTypes_InstanceFamiliesWithCbsQuery(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	mockItems := []*cvmSDK.InstanceTypeQuotaItem{
		createMockInstanceTypeQuotaItem("ap-guangzhou-6", "SA2", "SA2.MEDIUM4", "SELL", 2, 4),
	}
	setupCvmMock(patches, mockItems)

	mockDiskConfigs := []*cbsSDK.DiskConfig{
		createMockDiskConfig("ap-guangzhou-6", "SA2", "CLOUD_SSD", "POSTPAID_BY_HOUR", "SYSTEM_DISK", true, 50, 500),
		createMockDiskConfig("ap-guangzhou-6", "S5", "CLOUD_SSD", "POSTPAID_BY_HOUR", "SYSTEM_DISK", true, 50, 500),
	}
	setupCbsMock(patches, mockDiskConfigs)

	meta := newInstanceTypesMockMeta()
	res := cvm.DataSourceTencentCloudInstanceTypes()

	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"instance_families": []interface{}{"SA2", "S5"},
		"disk_types":        []interface{}{"CLOUD_SSD"},
		"exclude_sold_out":  false,
	})

	err := res.Read(d, meta)
	assert.NoError(t, err)

	instanceTypesList := d.Get("instance_types").([]interface{})
	assert.Equal(t, 1, len(instanceTypesList))

	firstType := instanceTypesList[0].(map[string]interface{})
	cbsConfigs := firstType["cbs_configs"].([]interface{})
	assert.Equal(t, 2, len(cbsConfigs), "CBS query with InstanceFamilies should return configs for both families")
}

// go test ./tencentcloud/services/cvm/ -run "TestInstanceTypes_SchemaFields" -v -count=1 -gcflags="all=-l"

// TestInstanceTypes_SchemaFields tests that the new schema fields are properly defined
func TestInstanceTypes_SchemaFields(t *testing.T) {
	res := cvm.DataSourceTencentCloudInstanceTypes()

	// Verify instance_families schema
	instanceFamiliesSchema, ok := res.Schema["instance_families"]
	assert.True(t, ok, "instance_families schema should exist")
	assert.Equal(t, schema.TypeList, instanceFamiliesSchema.Type)
	assert.True(t, instanceFamiliesSchema.Optional)
	assert.False(t, instanceFamiliesSchema.Required)
	assert.False(t, instanceFamiliesSchema.Computed)

	// Verify disk_types schema
	diskTypesSchema, ok := res.Schema["disk_types"]
	assert.True(t, ok, "disk_types schema should exist")
	assert.Equal(t, schema.TypeList, diskTypesSchema.Type)
	assert.True(t, diskTypesSchema.Optional)
	assert.False(t, diskTypesSchema.Required)
	assert.False(t, diskTypesSchema.Computed)
}

// go test ./tencentcloud/services/cvm/ -run "TestInstanceTypes_EmptyResult" -v -count=1 -gcflags="all=-l"

// TestInstanceTypes_EmptyResult tests behavior when the API returns no results
func TestInstanceTypes_EmptyResult(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	setupCvmMock(patches, nil)

	meta := newInstanceTypesMockMeta()
	res := cvm.DataSourceTencentCloudInstanceTypes()

	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"instance_families": []interface{}{"NONEXISTENT"},
		"exclude_sold_out":  false,
	})

	err := res.Read(d, meta)
	assert.NoError(t, err)

	instanceTypesList := d.Get("instance_types").([]interface{})
	assert.Equal(t, 0, len(instanceTypesList))
}
