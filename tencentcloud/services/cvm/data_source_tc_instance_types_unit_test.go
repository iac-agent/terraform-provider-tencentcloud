package cvm_test

import (
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"

	cvmSDK "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
	cbsSDK "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cbs/v20170312"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/cvm"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
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

func ptrString(v string) *string    { return &v }
func ptrInt64(v int64) *int64       { return &v }
func ptrFloat64(v float64) *float64 { return &v }
func ptrInt(v int) *int             { return &v }
func ptrBool(v bool) *bool          { return &v }
func ptrUint64(v uint64) *uint64    { return &v }

func makeMockInstanceType() *cvmSDK.InstanceTypeQuotaItem {
	return &cvmSDK.InstanceTypeQuotaItem{
		Zone:           ptrString("ap-guangzhou-6"),
		Cpu:            ptrInt64(2),
		Memory:         ptrInt64(2),
		InstanceFamily: ptrString("S6"),
		InstanceType:   ptrString("S6.MEDIUM2"),
		Status:         ptrString("SELL"),
	}
}

func makeMockDiskConfig() *cbsSDK.DiskConfig {
	return &cbsSDK.DiskConfig{
		Available:      ptrBool(true),
		DiskChargeType: ptrString("PREPAID"),
		Zone:           ptrString("ap-guangzhou-6"),
		InstanceFamily: ptrString("S6"),
		DiskType:       ptrString("CLOUD_SSD"),
		DiskUsage:      ptrString("SYSTEM_DISK"),
		MinDiskSize:    ptrUint64(50),
		MaxDiskSize:    ptrUint64(500),
		StepSize:       ptrUint64(10),
	}
}

// setupInstanceTypesPatches sets up common gomonkey patches for CVM and CBS API clients
func setupInstanceTypesPatches(patches *gomonkey.Patches, cvmClient *cvmSDK.Client, cbsClient *cbsSDK.Client) {
	// Patch ratelimit.Check to avoid rate limiting in tests
	patches.ApplyFunc(ratelimit.Check, func(action string) {})

	// Patch TencentCloudClient.UseCvmClient to return mock cvm client
	patches.ApplyMethodReturn(&connectivity.TencentCloudClient{}, "UseCvmClient", cvmClient)

	// Patch TencentCloudClient.UseCbsClient to return mock cbs client
	patches.ApplyMethodReturn(&connectivity.TencentCloudClient{}, "UseCbsClient", cbsClient)
}

// go test ./tencentcloud/services/cvm/ -run "TestInstanceTypes" -v -count=1 -gcflags="all=-l"

// TestInstanceTypes_Schema_TopLevelParams validates the new top-level parameters exist in schema
func TestInstanceTypes_Schema_TopLevelParams(t *testing.T) {
	res := cvm.DataSourceTencentCloudInstanceTypes()

	assert.NotNil(t, res)
	assert.NotNil(t, res.Read)

	// Check disk_types parameter
	assert.Contains(t, res.Schema, "disk_types")
	diskTypes := res.Schema["disk_types"]
	assert.Equal(t, schema.TypeList, diskTypes.Type)
	assert.True(t, diskTypes.Optional)
	assert.Contains(t, diskTypes.Description, "CLOUD_SSD")

	// Check zones parameter
	assert.Contains(t, res.Schema, "zones")
	zones := res.Schema["zones"]
	assert.Equal(t, schema.TypeList, zones.Type)
	assert.True(t, zones.Optional)

	// Check memory parameter
	assert.Contains(t, res.Schema, "memory")
	memory := res.Schema["memory"]
	assert.Equal(t, schema.TypeInt, memory.Type)
	assert.True(t, memory.Optional)
}

// TestInstanceTypes_DiskTypesOverride tests top-level disk_types overrides cbs_filter.disk_types
func TestInstanceTypes_DiskTypesOverride(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	cvmClient := &cvmSDK.Client{}
	cbsClient := &cbsSDK.Client{}
	setupInstanceTypesPatches(patches, cvmClient, cbsClient)

	// Mock CVM DescribeZoneInstanceConfigInfos
	patches.ApplyMethodFunc(cvmClient, "DescribeZoneInstanceConfigInfos",
		func(request *cvmSDK.DescribeZoneInstanceConfigInfosRequest) (*cvmSDK.DescribeZoneInstanceConfigInfosResponse, error) {
			resp := cvmSDK.NewDescribeZoneInstanceConfigInfosResponse()
			resp.Response = &cvmSDK.DescribeZoneInstanceConfigInfosResponseParams{
				InstanceTypeQuotaSet: []*cvmSDK.InstanceTypeQuotaItem{makeMockInstanceType()},
				RequestId:           ptrString("fake-request-id"),
			}
			return resp, nil
		},
	)

	// Capture CBS DescribeDiskConfigQuota request params
	var capturedDiskQuotaRequest *cbsSDK.DescribeDiskConfigQuotaRequest
	patches.ApplyMethodFunc(cbsClient, "DescribeDiskConfigQuota",
		func(request *cbsSDK.DescribeDiskConfigQuotaRequest) (*cbsSDK.DescribeDiskConfigQuotaResponse, error) {
			capturedDiskQuotaRequest = request
			resp := cbsSDK.NewDescribeDiskConfigQuotaResponse()
			resp.Response = &cbsSDK.DescribeDiskConfigQuotaResponseParams{
				DiskConfigSet: []*cbsSDK.DiskConfig{makeMockDiskConfig()},
				RequestId:     ptrString("fake-request-id"),
			}
			return resp, nil
		},
	)

	meta := newInstanceTypesMockMeta()
	res := cvm.DataSourceTencentCloudInstanceTypes()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"availability_zone": "ap-guangzhou-6",
		"cpu_core_count":    2,
		"memory_size":       2,
		"exclude_sold_out":  false,
		"cbs_filter": []interface{}{
			map[string]interface{}{
				"disk_types":       []interface{}{"CLOUD_SSD"},
				"disk_charge_type": "PREPAID",
				"disk_usage":       "SYSTEM_DISK",
			},
		},
		// Top-level disk_types overrides cbs_filter.disk_types
		"disk_types": []interface{}{"CLOUD_HSSD"},
	})

	err := res.Read(d, meta)
	assert.NoError(t, err)

	// Verify disk_types was overridden to CLOUD_HSSD
	assert.NotNil(t, capturedDiskQuotaRequest)
	assert.NotNil(t, capturedDiskQuotaRequest.DiskTypes)
	assert.Equal(t, 1, len(capturedDiskQuotaRequest.DiskTypes))
	assert.Equal(t, "CLOUD_HSSD", *capturedDiskQuotaRequest.DiskTypes[0])

	// Verify cbs_configs is populated in result
	instanceTypes := d.Get("instance_types").([]interface{})
	assert.GreaterOrEqual(t, len(instanceTypes), 1)
}

// TestInstanceTypes_ZonesOverride tests top-level zones overrides derived availability_zone
func TestInstanceTypes_ZonesOverride(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	cvmClient := &cvmSDK.Client{}
	cbsClient := &cbsSDK.Client{}
	setupInstanceTypesPatches(patches, cvmClient, cbsClient)

	// Mock CVM DescribeZoneInstanceConfigInfos
	patches.ApplyMethodFunc(cvmClient, "DescribeZoneInstanceConfigInfos",
		func(request *cvmSDK.DescribeZoneInstanceConfigInfosRequest) (*cvmSDK.DescribeZoneInstanceConfigInfosResponse, error) {
			resp := cvmSDK.NewDescribeZoneInstanceConfigInfosResponse()
			resp.Response = &cvmSDK.DescribeZoneInstanceConfigInfosResponseParams{
				InstanceTypeQuotaSet: []*cvmSDK.InstanceTypeQuotaItem{makeMockInstanceType()},
				RequestId:           ptrString("fake-request-id"),
			}
			return resp, nil
		},
	)

	// Capture CBS DescribeDiskConfigQuota request params
	var capturedDiskQuotaRequest *cbsSDK.DescribeDiskConfigQuotaRequest
	patches.ApplyMethodFunc(cbsClient, "DescribeDiskConfigQuota",
		func(request *cbsSDK.DescribeDiskConfigQuotaRequest) (*cbsSDK.DescribeDiskConfigQuotaResponse, error) {
			capturedDiskQuotaRequest = request
			resp := cbsSDK.NewDescribeDiskConfigQuotaResponse()
			resp.Response = &cbsSDK.DescribeDiskConfigQuotaResponseParams{
				DiskConfigSet: []*cbsSDK.DiskConfig{makeMockDiskConfig()},
				RequestId:     ptrString("fake-request-id"),
			}
			return resp, nil
		},
	)

	meta := newInstanceTypesMockMeta()
	res := cvm.DataSourceTencentCloudInstanceTypes()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"availability_zone": "ap-guangzhou-6",
		"cpu_core_count":    2,
		"memory_size":       2,
		"exclude_sold_out":  false,
		"cbs_filter": []interface{}{
			map[string]interface{}{
				"disk_types":       []interface{}{"CLOUD_SSD"},
				"disk_charge_type": "PREPAID",
				"disk_usage":       "SYSTEM_DISK",
			},
		},
		// Top-level zones overrides derived availability_zone
		"zones": []interface{}{"ap-guangzhou-6", "ap-guangzhou-3"},
	})

	err := res.Read(d, meta)
	assert.NoError(t, err)

	// Verify zones was used instead of derived availability_zone
	assert.NotNil(t, capturedDiskQuotaRequest)
	assert.NotNil(t, capturedDiskQuotaRequest.Zones)
	assert.Equal(t, 2, len(capturedDiskQuotaRequest.Zones))
	assert.Equal(t, "ap-guangzhou-6", *capturedDiskQuotaRequest.Zones[0])
	assert.Equal(t, "ap-guangzhou-3", *capturedDiskQuotaRequest.Zones[1])
}

// TestInstanceTypes_MemoryOverride tests top-level memory overrides derived memory_size
func TestInstanceTypes_MemoryOverride(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	cvmClient := &cvmSDK.Client{}
	cbsClient := &cbsSDK.Client{}
	setupInstanceTypesPatches(patches, cvmClient, cbsClient)

	// Mock CVM DescribeZoneInstanceConfigInfos
	patches.ApplyMethodFunc(cvmClient, "DescribeZoneInstanceConfigInfos",
		func(request *cvmSDK.DescribeZoneInstanceConfigInfosRequest) (*cvmSDK.DescribeZoneInstanceConfigInfosResponse, error) {
			resp := cvmSDK.NewDescribeZoneInstanceConfigInfosResponse()
			resp.Response = &cvmSDK.DescribeZoneInstanceConfigInfosResponseParams{
				InstanceTypeQuotaSet: []*cvmSDK.InstanceTypeQuotaItem{makeMockInstanceType()},
				RequestId:           ptrString("fake-request-id"),
			}
			return resp, nil
		},
	)

	// Capture CBS DescribeDiskConfigQuota request params
	var capturedDiskQuotaRequest *cbsSDK.DescribeDiskConfigQuotaRequest
	patches.ApplyMethodFunc(cbsClient, "DescribeDiskConfigQuota",
		func(request *cbsSDK.DescribeDiskConfigQuotaRequest) (*cbsSDK.DescribeDiskConfigQuotaResponse, error) {
			capturedDiskQuotaRequest = request
			resp := cbsSDK.NewDescribeDiskConfigQuotaResponse()
			resp.Response = &cbsSDK.DescribeDiskConfigQuotaResponseParams{
				DiskConfigSet: []*cbsSDK.DiskConfig{makeMockDiskConfig()},
				RequestId:     ptrString("fake-request-id"),
			}
			return resp, nil
		},
	)

	meta := newInstanceTypesMockMeta()
	res := cvm.DataSourceTencentCloudInstanceTypes()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"availability_zone": "ap-guangzhou-6",
		"cpu_core_count":    2,
		"memory_size":       2,
		"exclude_sold_out":  false,
		"cbs_filter": []interface{}{
			map[string]interface{}{
				"disk_types":       []interface{}{"CLOUD_SSD"},
				"disk_charge_type": "PREPAID",
				"disk_usage":       "SYSTEM_DISK",
			},
		},
		// Top-level memory overrides derived memory_size (2 -> 8)
		"memory": 8,
	})

	err := res.Read(d, meta)
	assert.NoError(t, err)

	// Verify memory was used (8 GB) instead of derived memory_size (2 GB)
	assert.NotNil(t, capturedDiskQuotaRequest)
	assert.NotNil(t, capturedDiskQuotaRequest.Memory)
	assert.Equal(t, uint64(8), *capturedDiskQuotaRequest.Memory)
}

// TestInstanceTypes_BackwardCompatibility tests backward compatibility when new parameters are not specified
func TestInstanceTypes_BackwardCompatibility(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	cvmClient := &cvmSDK.Client{}
	cbsClient := &cbsSDK.Client{}
	setupInstanceTypesPatches(patches, cvmClient, cbsClient)

	// Mock CVM DescribeZoneInstanceConfigInfos
	patches.ApplyMethodFunc(cvmClient, "DescribeZoneInstanceConfigInfos",
		func(request *cvmSDK.DescribeZoneInstanceConfigInfosRequest) (*cvmSDK.DescribeZoneInstanceConfigInfosResponse, error) {
			resp := cvmSDK.NewDescribeZoneInstanceConfigInfosResponse()
			resp.Response = &cvmSDK.DescribeZoneInstanceConfigInfosResponseParams{
				InstanceTypeQuotaSet: []*cvmSDK.InstanceTypeQuotaItem{makeMockInstanceType()},
				RequestId:           ptrString("fake-request-id"),
			}
			return resp, nil
		},
	)

	// Capture CBS DescribeDiskConfigQuota request params
	var capturedDiskQuotaRequest *cbsSDK.DescribeDiskConfigQuotaRequest
	patches.ApplyMethodFunc(cbsClient, "DescribeDiskConfigQuota",
		func(request *cbsSDK.DescribeDiskConfigQuotaRequest) (*cbsSDK.DescribeDiskConfigQuotaResponse, error) {
			capturedDiskQuotaRequest = request
			resp := cbsSDK.NewDescribeDiskConfigQuotaResponse()
			resp.Response = &cbsSDK.DescribeDiskConfigQuotaResponseParams{
				DiskConfigSet: []*cbsSDK.DiskConfig{makeMockDiskConfig()},
				RequestId:     ptrString("fake-request-id"),
			}
			return resp, nil
		},
	)

	meta := newInstanceTypesMockMeta()
	res := cvm.DataSourceTencentCloudInstanceTypes()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"availability_zone": "ap-guangzhou-6",
		"cpu_core_count":    2,
		"memory_size":       2,
		"exclude_sold_out":  false,
		"cbs_filter": []interface{}{
			map[string]interface{}{
				"disk_types":       []interface{}{"CLOUD_SSD"},
				"disk_charge_type": "PREPAID",
				"disk_usage":       "SYSTEM_DISK",
			},
		},
		// No new top-level parameters - backward compatibility test
	})

	err := res.Read(d, meta)
	assert.NoError(t, err)

	// Verify derived values are used (backward compatibility)
	// Zones should be derived from availability_zone
	assert.NotNil(t, capturedDiskQuotaRequest)
	assert.NotNil(t, capturedDiskQuotaRequest.Zones)
	assert.Equal(t, 1, len(capturedDiskQuotaRequest.Zones))
	assert.Equal(t, "ap-guangzhou-6", *capturedDiskQuotaRequest.Zones[0])

	// Memory should be derived from memory_size
	assert.NotNil(t, capturedDiskQuotaRequest.Memory)
	assert.Equal(t, uint64(2), *capturedDiskQuotaRequest.Memory)

	// disk_types should come from cbs_filter
	assert.NotNil(t, capturedDiskQuotaRequest.DiskTypes)
	assert.Equal(t, 1, len(capturedDiskQuotaRequest.DiskTypes))
	assert.Equal(t, "CLOUD_SSD", *capturedDiskQuotaRequest.DiskTypes[0])

	// Verify cbs_configs is populated
	instanceTypes := d.Get("instance_types").([]interface{})
	assert.GreaterOrEqual(t, len(instanceTypes), 1)
}
