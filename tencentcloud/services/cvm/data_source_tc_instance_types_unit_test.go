package cvm_test

import (
	"context"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	cbs "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cbs/v20170312"
	cvmSDK "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"

	svccbs "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/cbs"
	svcvm "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/cvm"
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

func ptrString(s string) *string { return &s }
func ptrInt64(i int64) *int64 { return &i }
func ptrBool(b bool) *bool { return &b }

// go test ./tencentcloud/services/cvm/ -run "TestInstanceTypesInquiryType" -v -count=1 -gcflags="all=-l"
// TestInstanceTypesInquiryType verifies InquiryType parameter is correctly passed to the CBS API request
func TestInstanceTypesInquiryType(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	// Mock CvmService.DescribeInstancesSellTypeByFilter
	patches.ApplyMethodFunc(&svcvm.CvmService{}, "DescribeInstancesSellTypeByFilter",
		func(ctx context.Context, filters map[string][]string) ([]*cvmSDK.InstanceTypeQuotaItem, error) {
			return []*cvmSDK.InstanceTypeQuotaItem{
				{
					Zone:            ptrString("ap-guangzhou-6"),
					Cpu:             ptrInt64(4),
					Memory:          ptrInt64(8),
					InstanceFamily:  ptrString("S5"),
					InstanceType:    ptrString("S5.MEDIUM4"),
					Status:          ptrString("SELL"),
					InstanceChargeType: ptrString("POSTPAID_BY_HOUR"),
				},
			}, nil
		},
	)

	// Mock CbsService.DescribeDiskConfigQuota to capture inquiry_type
	var capturedInquiryType string
	patches.ApplyMethodFunc(&svccbs.CbsService{}, "DescribeDiskConfigQuota",
		func(ctx context.Context, cvmInfo map[string]interface{}) ([]*cbs.DiskConfig, error) {
			if v, ok := cvmInfo["inquiry_type"].(string); ok {
				capturedInquiryType = v
			}
			return []*cbs.DiskConfig{
				{
					Available:       ptrBool(true),
					DiskChargeType:  ptrString("PREPAID"),
					Zone:            ptrString("ap-guangzhou-6"),
					InstanceFamily:  ptrString("S5"),
					DiskType:        ptrString("CLOUD_SSD"),
				},
			}, nil
		},
	)

	res := svcvm.DataSourceTencentCloudInstanceTypes()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"cpu_core_count":    4,
		"memory_size":       8,
		"availability_zone": "ap-guangzhou-6",
		"inquiry_type":      "INQUIRY_CBS_CONFIG",
		"disk_charge_type":  "PREPAID",
		"exclude_sold_out":  false,
	})

	meta := newMockMeta()
	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "INQUIRY_CBS_CONFIG", capturedInquiryType)
}

// go test ./tencentcloud/services/cvm/ -run "TestInstanceTypesDiskChargeTypePrecedence" -v -count=1 -gcflags="all=-l"
// TestInstanceTypesDiskChargeTypePrecedence verifies top-level DiskChargeType takes precedence over cbs_filter.disk_charge_type
func TestInstanceTypesDiskChargeTypePrecedence(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	patches.ApplyMethodFunc(&svcvm.CvmService{}, "DescribeInstancesSellTypeByFilter",
		func(ctx context.Context, filters map[string][]string) ([]*cvmSDK.InstanceTypeQuotaItem, error) {
			return []*cvmSDK.InstanceTypeQuotaItem{
				{
					Zone:            ptrString("ap-guangzhou-6"),
					Cpu:             ptrInt64(4),
					Memory:          ptrInt64(8),
					InstanceFamily:  ptrString("S5"),
					InstanceType:    ptrString("S5.MEDIUM4"),
					Status:          ptrString("SELL"),
					InstanceChargeType: ptrString("POSTPAID_BY_HOUR"),
				},
			}, nil
		},
	)

	var capturedDiskChargeType string
	patches.ApplyMethodFunc(&svccbs.CbsService{}, "DescribeDiskConfigQuota",
		func(ctx context.Context, cvmInfo map[string]interface{}) ([]*cbs.DiskConfig, error) {
			if v, ok := cvmInfo["disk_charge_type"].(string); ok {
				capturedDiskChargeType = v
			}
			return []*cbs.DiskConfig{
				{
					Available:       ptrBool(true),
					DiskChargeType:  ptrString("POSTPAID_BY_HOUR"),
					Zone:            ptrString("ap-guangzhou-6"),
					InstanceFamily:  ptrString("S5"),
					DiskType:        ptrString("CLOUD_SSD"),
				},
			}, nil
		},
	)

	res := svcvm.DataSourceTencentCloudInstanceTypes()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"cpu_core_count":    4,
		"memory_size":       8,
		"availability_zone": "ap-guangzhou-6",
		"disk_charge_type":  "POSTPAID_BY_HOUR",
		"exclude_sold_out":  false,
		"cbs_filter": []interface{}{
			map[string]interface{}{
				"disk_types":       []interface{}{"CLOUD_SSD"},
				"disk_charge_type": "PREPAID",
				"disk_usage":       "SYSTEM_DISK",
			},
		},
	})

	meta := newMockMeta()
	err := res.Read(d, meta)
	assert.NoError(t, err)
	// Top-level disk_charge_type should take precedence
	assert.Equal(t, "POSTPAID_BY_HOUR", capturedDiskChargeType)
}

// go test ./tencentcloud/services/cvm/ -run "TestInstanceTypesInstanceFamiliesOverride" -v -count=1 -gcflags="all=-l"
// TestInstanceTypesInstanceFamiliesOverride verifies top-level InstanceFamilies overrides auto-populated family
func TestInstanceTypesInstanceFamiliesOverride(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	patches.ApplyMethodFunc(&svcvm.CvmService{}, "DescribeInstancesSellTypeByFilter",
		func(ctx context.Context, filters map[string][]string) ([]*cvmSDK.InstanceTypeQuotaItem, error) {
			return []*cvmSDK.InstanceTypeQuotaItem{
				{
					Zone:            ptrString("ap-guangzhou-6"),
					Cpu:             ptrInt64(4),
					Memory:          ptrInt64(8),
					InstanceFamily:  ptrString("S5"),
					InstanceType:    ptrString("S5.MEDIUM4"),
					Status:          ptrString("SELL"),
					InstanceChargeType: ptrString("POSTPAID_BY_HOUR"),
				},
			}, nil
		},
	)

	var capturedInstanceFamilies []string
	var hasInstanceFamilies bool
	patches.ApplyMethodFunc(&svccbs.CbsService{}, "DescribeDiskConfigQuota",
		func(ctx context.Context, cvmInfo map[string]interface{}) ([]*cbs.DiskConfig, error) {
			if v, ok := cvmInfo["instance_families"].([]string); ok {
				capturedInstanceFamilies = v
				hasInstanceFamilies = true
			}
			return []*cbs.DiskConfig{
				{
					Available:       ptrBool(true),
					DiskChargeType:  ptrString("PREPAID"),
					Zone:            ptrString("ap-guangzhou-6"),
					InstanceFamily:  ptrString("S5"),
					DiskType:        ptrString("CLOUD_SSD"),
				},
			}, nil
		},
	)

	res := svcvm.DataSourceTencentCloudInstanceTypes()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"cpu_core_count":    4,
		"memory_size":       8,
		"availability_zone": "ap-guangzhou-6",
		"instance_families": []interface{}{"S5", "M5"},
		"disk_charge_type":  "PREPAID",
		"exclude_sold_out":  false,
	})

	meta := newMockMeta()
	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.True(t, hasInstanceFamilies)
	assert.Equal(t, []string{"S5", "M5"}, capturedInstanceFamilies)
}

// go test ./tencentcloud/services/cvm/ -run "TestInstanceTypesAvailableAttribute" -v -count=1 -gcflags="all=-l"
// TestInstanceTypesAvailableAttribute verifies Available computed attribute is correctly set based on DiskConfig availability
func TestInstanceTypesAvailableAttribute(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	patches.ApplyMethodFunc(&svcvm.CvmService{}, "DescribeInstancesSellTypeByFilter",
		func(ctx context.Context, filters map[string][]string) ([]*cvmSDK.InstanceTypeQuotaItem, error) {
			return []*cvmSDK.InstanceTypeQuotaItem{
				{
					Zone:            ptrString("ap-guangzhou-6"),
					Cpu:             ptrInt64(4),
					Memory:          ptrInt64(8),
					InstanceFamily:  ptrString("S5"),
					InstanceType:    ptrString("S5.MEDIUM4"),
					Status:          ptrString("SELL"),
					InstanceChargeType: ptrString("POSTPAID_BY_HOUR"),
				},
			}, nil
		},
	)

	// Test case: some disk configs are available
	patches.ApplyMethodFunc(&svccbs.CbsService{}, "DescribeDiskConfigQuota",
		func(ctx context.Context, cvmInfo map[string]interface{}) ([]*cbs.DiskConfig, error) {
			return []*cbs.DiskConfig{
				{
					Available:       ptrBool(true),
					DiskChargeType:  ptrString("PREPAID"),
					Zone:            ptrString("ap-guangzhou-6"),
					InstanceFamily:  ptrString("S5"),
					DiskType:        ptrString("CLOUD_SSD"),
				},
				{
					Available:       ptrBool(false),
					DiskChargeType:  ptrString("PREPAID"),
					Zone:            ptrString("ap-guangzhou-6"),
					InstanceFamily:  ptrString("S5"),
					DiskType:        ptrString("CLOUD_PREMIUM"),
				},
			}, nil
		},
	)

	res := svcvm.DataSourceTencentCloudInstanceTypes()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"cpu_core_count":    4,
		"memory_size":       8,
		"availability_zone": "ap-guangzhou-6",
		"disk_charge_type":  "PREPAID",
		"exclude_sold_out":  false,
	})

	meta := newMockMeta()
	err := res.Read(d, meta)
	assert.NoError(t, err)
	// available should be true since at least one DiskConfig has Available=true
	assert.Equal(t, true, d.Get("available"))
}

// go test ./tencentcloud/services/cvm/ -run "TestInstanceTypesAvailableAllFalse" -v -count=1 -gcflags="all=-l"
// TestInstanceTypesAvailableAllFalse verifies Available is false when all DiskConfigs have Available=false
func TestInstanceTypesAvailableAllFalse(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	patches.ApplyMethodFunc(&svcvm.CvmService{}, "DescribeInstancesSellTypeByFilter",
		func(ctx context.Context, filters map[string][]string) ([]*cvmSDK.InstanceTypeQuotaItem, error) {
			return []*cvmSDK.InstanceTypeQuotaItem{
				{
					Zone:            ptrString("ap-guangzhou-6"),
					Cpu:             ptrInt64(4),
					Memory:          ptrInt64(8),
					InstanceFamily:  ptrString("S5"),
					InstanceType:    ptrString("S5.MEDIUM4"),
					Status:          ptrString("SELL"),
					InstanceChargeType: ptrString("POSTPAID_BY_HOUR"),
				},
			}, nil
		},
	)

	patches.ApplyMethodFunc(&svccbs.CbsService{}, "DescribeDiskConfigQuota",
		func(ctx context.Context, cvmInfo map[string]interface{}) ([]*cbs.DiskConfig, error) {
			return []*cbs.DiskConfig{
				{
					Available:       ptrBool(false),
					DiskChargeType:  ptrString("PREPAID"),
					Zone:            ptrString("ap-guangzhou-6"),
					InstanceFamily:  ptrString("S5"),
					DiskType:        ptrString("CLOUD_SSD"),
				},
			}, nil
		},
	)

	res := svcvm.DataSourceTencentCloudInstanceTypes()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"cpu_core_count":    4,
		"memory_size":       8,
		"availability_zone": "ap-guangzhou-6",
		"disk_charge_type":  "PREPAID",
		"exclude_sold_out":  false,
	})

	meta := newMockMeta()
	err := res.Read(d, meta)
	assert.NoError(t, err)
	// available should be false since all DiskConfigs have Available=false
	assert.Equal(t, false, d.Get("available"))
}

// go test ./tencentcloud/services/cvm/ -run "TestInstanceTypesBackwardCompatibility" -v -count=1 -gcflags="all=-l"
// TestInstanceTypesBackwardCompatibility verifies backward compatibility when new parameters are not provided
func TestInstanceTypesBackwardCompatibility(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	patches.ApplyMethodFunc(&svcvm.CvmService{}, "DescribeInstancesSellTypeByFilter",
		func(ctx context.Context, filters map[string][]string) ([]*cvmSDK.InstanceTypeQuotaItem, error) {
			return []*cvmSDK.InstanceTypeQuotaItem{
				{
					Zone:            ptrString("ap-guangzhou-6"),
					Cpu:             ptrInt64(4),
					Memory:          ptrInt64(8),
					InstanceFamily:  ptrString("S5"),
					InstanceType:    ptrString("S5.MEDIUM4"),
					Status:          ptrString("SELL"),
					InstanceChargeType: ptrString("POSTPAID_BY_HOUR"),
				},
			}, nil
		},
	)

	var capturedInquiryType string
	var capturedFamily string
	var capturedDiskChargeType string
	var cbsQueryCalled bool
	patches.ApplyMethodFunc(&svccbs.CbsService{}, "DescribeDiskConfigQuota",
		func(ctx context.Context, cvmInfo map[string]interface{}) ([]*cbs.DiskConfig, error) {
			cbsQueryCalled = true
			if v, ok := cvmInfo["inquiry_type"].(string); ok {
				capturedInquiryType = v
			}
			if v, ok := cvmInfo["family"].(string); ok {
				capturedFamily = v
			}
			if v, ok := cvmInfo["disk_charge_type"].(string); ok {
				capturedDiskChargeType = v
			}
			return []*cbs.DiskConfig{
				{
					Available:       ptrBool(true),
					DiskChargeType:  ptrString("PREPAID"),
					Zone:            ptrString("ap-guangzhou-6"),
					InstanceFamily:  ptrString("S5"),
					DiskType:        ptrString("CLOUD_SSD"),
				},
			}, nil
		},
	)

	res := svcvm.DataSourceTencentCloudInstanceTypes()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"cpu_core_count":    4,
		"memory_size":       8,
		"availability_zone": "ap-guangzhou-6",
		"exclude_sold_out":  false,
		"cbs_filter": []interface{}{
			map[string]interface{}{
				"disk_types":       []interface{}{"CLOUD_SSD"},
				"disk_charge_type": "PREPAID",
				"disk_usage":       "SYSTEM_DISK",
			},
		},
	})

	meta := newMockMeta()
	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.True(t, cbsQueryCalled)
	// Default inquiry_type should be INQUIRY_CVM_CONFIG
	assert.Equal(t, "INQUIRY_CVM_CONFIG", capturedInquiryType)
	// Family should be auto-populated from instance type result
	assert.Equal(t, "S5", capturedFamily)
	// disk_charge_type should come from cbs_filter
	assert.Equal(t, "PREPAID", capturedDiskChargeType)
}

// go test ./tencentcloud/services/cvm/ -run "TestInstanceTypesNoCbsQueryAvailableFalse" -v -count=1 -gcflags="all=-l"
// TestInstanceTypesNoCbsQueryAvailableFalse verifies Available defaults to false when no CBS query is triggered
func TestInstanceTypesNoCbsQueryAvailableFalse(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	patches.ApplyMethodFunc(&svcvm.CvmService{}, "DescribeInstancesSellTypeByFilter",
		func(ctx context.Context, filters map[string][]string) ([]*cvmSDK.InstanceTypeQuotaItem, error) {
			return []*cvmSDK.InstanceTypeQuotaItem{
				{
					Zone:            ptrString("ap-guangzhou-6"),
					Cpu:             ptrInt64(4),
					Memory:          ptrInt64(8),
					InstanceFamily:  ptrString("S5"),
					InstanceType:    ptrString("S5.MEDIUM4"),
					Status:          ptrString("SELL"),
					InstanceChargeType: ptrString("POSTPAID_BY_HOUR"),
				},
			}, nil
		},
	)

	var cbsQueryCalled bool
	patches.ApplyMethodFunc(&svccbs.CbsService{}, "DescribeDiskConfigQuota",
		func(ctx context.Context, cvmInfo map[string]interface{}) ([]*cbs.DiskConfig, error) {
			cbsQueryCalled = true
			return []*cbs.DiskConfig{}, nil
		},
	)

	res := svcvm.DataSourceTencentCloudInstanceTypes()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"cpu_core_count":    4,
		"memory_size":       8,
		"availability_zone": "ap-guangzhou-6",
		"exclude_sold_out":  false,
	})

	meta := newMockMeta()
	err := res.Read(d, meta)
	assert.NoError(t, err)
	// CBS query should NOT be called when no cbs_filter or top-level CBS params are provided
	assert.False(t, cbsQueryCalled)
	// Available should default to false when no CBS query
	assert.Equal(t, false, d.Get("available"))
}
