package cvm_test

import (
	"context"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"

	cbs "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cbs/v20170312"
	cvm_sdk "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	svccbs "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/cbs"
	svccvm "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/cvm"
)

type mockInstanceTypesMeta struct {
	client *connectivity.TencentCloudClient
}

func (m *mockInstanceTypesMeta) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockInstanceTypesMeta{}

func newInstanceTypesMockMeta() *mockInstanceTypesMeta {
	return &mockInstanceTypesMeta{client: &connectivity.TencentCloudClient{}}
}

func ptrString(s string) *string {
	return &s
}

func ptrInt64(i int64) *int64 {
	return &i
}

func ptrFloat64(f float64) *float64 {
	return &f
}

func ptrBool(b bool) *bool {
	return &b
}

func ptrUint64(u uint64) *uint64 {
	return &u
}

// TestInstanceTypesInquiryTypePassedToCBS verifies that inquiry_type is correctly extracted from cbs_filter and passed to CBS service
func TestInstanceTypesInquiryTypePassedToCBS(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	mockMeta := newInstanceTypesMockMeta()

	// Mock CvmService.DescribeInstancesSellTypeByFilter
	patches.ApplyMethodFunc(&svccvm.CvmService{}, "DescribeInstancesSellTypeByFilter",
		func(ctx context.Context, filterMap map[string][]string) ([]*cvm_sdk.InstanceTypeQuotaItem, error) {
			return []*cvm_sdk.InstanceTypeQuotaItem{
				{
					Zone:               ptrString("ap-guangzhou-6"),
					Cpu:                ptrInt64(2),
					Gpu:                ptrInt64(0),
					Memory:             ptrInt64(2),
					InstanceFamily:     ptrString("S5"),
					InstanceType:       ptrString("S5.SMALL2"),
					InstanceChargeType: ptrString("POSTPAID_BY_HOUR"),
					Status:             ptrString("SELL"),
					NetworkCard:        ptrInt64(25),
					TypeName:           ptrString("S5.SMALL2"),
					SoldOutReason:      ptrString(""),
					InstanceBandwidth:  ptrFloat64(1.5),
					InstancePps:        ptrInt64(30),
					StorageBlockAmount: ptrInt64(0),
					CpuType:            ptrString("Intel Xeon Cascade Lake"),
					Fpga:               ptrInt64(0),
					GpuCount:           ptrFloat64(0),
					Frequency:          ptrString("2.5GHz"),
					StatusCategory:     ptrString("EnoughStock"),
					Remark:             ptrString(""),
					Price: &cvm_sdk.ItemPrice{
						UnitPrice:         ptrFloat64(0.06),
						ChargeUnit:        ptrString("HOUR"),
						OriginalPrice:     ptrFloat64(0.06),
						DiscountPrice:     ptrFloat64(0.06),
						Discount:          ptrFloat64(100),
						UnitPriceDiscount: ptrFloat64(0.06),
					},
				},
			}, nil
		},
	)

	// Mock CbsService.DescribeDiskConfigQuota - capture the inquiryType parameter
	var capturedInquiryType string
	patches.ApplyMethodFunc(&svccbs.CbsService{}, "DescribeDiskConfigQuota",
		func(ctx context.Context, cvmInfo map[string]interface{}, inquiryType string, instanceFamilies []string) ([]*cbs.DiskConfig, error) {
			capturedInquiryType = inquiryType
			return []*cbs.DiskConfig{
				{
					Available:      ptrBool(true),
					DiskChargeType: ptrString("PREPAID"),
					Zone:           ptrString("ap-guangzhou-6"),
					InstanceFamily: ptrString("S5"),
					DiskType:       ptrString("CLOUD_SSD"),
					StepSize:       ptrUint64(10),
					DiskUsage:      ptrString("SYSTEM_DISK"),
					MinDiskSize:    ptrUint64(20),
					MaxDiskSize:    ptrUint64(500),
					DeviceClass:    ptrString(""),
				},
			}, nil
		},
	)

	res := svccvm.DataSourceTencentCloudInstanceTypes()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"cpu_core_count":    2,
		"memory_size":       2,
		"exclude_sold_out":  false,
		"availability_zone": "ap-guangzhou-6",
		"cbs_filter": []interface{}{
			map[string]interface{}{
				"inquiry_type":     "INQUIRY_CBS_CONFIG",
				"disk_types":       []interface{}{"CLOUD_SSD"},
				"disk_charge_type": "PREPAID",
				"disk_usage":       "SYSTEM_DISK",
			},
		},
	})

	err := res.Read(d, mockMeta)
	assert.NoError(t, err)

	// Verify inquiry_type was correctly passed
	assert.Equal(t, "INQUIRY_CBS_CONFIG", capturedInquiryType)
}

// TestInstanceTypesInstanceFamiliesPassedToCBS verifies that instance_families is correctly extracted from cbs_filter and passed to CBS service
func TestInstanceTypesInstanceFamiliesPassedToCBS(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	mockMeta := newInstanceTypesMockMeta()

	// Mock CvmService.DescribeInstancesSellTypeByFilter
	patches.ApplyMethodFunc(&svccvm.CvmService{}, "DescribeInstancesSellTypeByFilter",
		func(ctx context.Context, filterMap map[string][]string) ([]*cvm_sdk.InstanceTypeQuotaItem, error) {
			return []*cvm_sdk.InstanceTypeQuotaItem{
				{
					Zone:               ptrString("ap-guangzhou-6"),
					Cpu:                ptrInt64(2),
					Gpu:                ptrInt64(0),
					Memory:             ptrInt64(2),
					InstanceFamily:     ptrString("S5"),
					InstanceType:       ptrString("S5.SMALL2"),
					InstanceChargeType: ptrString("POSTPAID_BY_HOUR"),
					Status:             ptrString("SELL"),
					NetworkCard:        ptrInt64(25),
					TypeName:           ptrString("S5.SMALL2"),
					SoldOutReason:      ptrString(""),
					InstanceBandwidth:  ptrFloat64(1.5),
					InstancePps:        ptrInt64(30),
					StorageBlockAmount: ptrInt64(0),
					CpuType:            ptrString("Intel Xeon Cascade Lake"),
					Fpga:               ptrInt64(0),
					GpuCount:           ptrFloat64(0),
					Frequency:          ptrString("2.5GHz"),
					StatusCategory:     ptrString("EnoughStock"),
					Remark:             ptrString(""),
					Price: &cvm_sdk.ItemPrice{
						UnitPrice:         ptrFloat64(0.06),
						ChargeUnit:        ptrString("HOUR"),
						OriginalPrice:     ptrFloat64(0.06),
						DiscountPrice:     ptrFloat64(0.06),
						Discount:          ptrFloat64(100),
						UnitPriceDiscount: ptrFloat64(0.06),
					},
				},
			}, nil
		},
	)

	// Mock CbsService.DescribeDiskConfigQuota - capture the instanceFamilies parameter
	var capturedInquiryType string
	var capturedInstanceFamilies []string
	patches.ApplyMethodFunc(&svccbs.CbsService{}, "DescribeDiskConfigQuota",
		func(ctx context.Context, cvmInfo map[string]interface{}, inquiryType string, instanceFamilies []string) ([]*cbs.DiskConfig, error) {
			capturedInquiryType = inquiryType
			capturedInstanceFamilies = instanceFamilies
			return []*cbs.DiskConfig{
				{
					Available:      ptrBool(true),
					DiskChargeType: ptrString("PREPAID"),
					Zone:           ptrString("ap-guangzhou-6"),
					InstanceFamily: ptrString("S5"),
					DiskType:       ptrString("CLOUD_SSD"),
					StepSize:       ptrUint64(10),
					DiskUsage:      ptrString("SYSTEM_DISK"),
					MinDiskSize:    ptrUint64(20),
					MaxDiskSize:    ptrUint64(500),
					DeviceClass:    ptrString(""),
				},
			}, nil
		},
	)

	res := svccvm.DataSourceTencentCloudInstanceTypes()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"cpu_core_count":    2,
		"memory_size":       2,
		"exclude_sold_out":  false,
		"availability_zone": "ap-guangzhou-6",
		"cbs_filter": []interface{}{
			map[string]interface{}{
				"instance_families": []interface{}{"S5", "SA2"},
				"disk_types":        []interface{}{"CLOUD_SSD"},
				"disk_charge_type":  "PREPAID",
				"disk_usage":        "SYSTEM_DISK",
			},
		},
	})

	err := res.Read(d, mockMeta)
	assert.NoError(t, err)

	// Verify instance_families was correctly passed
	assert.Equal(t, []string{"S5", "SA2"}, capturedInstanceFamilies)
	// Verify inquiry_type defaults to empty (which CBS service will treat as INQUIRY_CVM_CONFIG)
	assert.Equal(t, "", capturedInquiryType)
}

// TestInstanceTypesBackwardCompatibility verifies that when new parameters are not specified, default values are used correctly
func TestInstanceTypesBackwardCompatibility(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	mockMeta := newInstanceTypesMockMeta()

	// Mock CvmService.DescribeInstancesSellTypeByFilter
	patches.ApplyMethodFunc(&svccvm.CvmService{}, "DescribeInstancesSellTypeByFilter",
		func(ctx context.Context, filterMap map[string][]string) ([]*cvm_sdk.InstanceTypeQuotaItem, error) {
			return []*cvm_sdk.InstanceTypeQuotaItem{
				{
					Zone:               ptrString("ap-guangzhou-6"),
					Cpu:                ptrInt64(2),
					Gpu:                ptrInt64(0),
					Memory:             ptrInt64(2),
					InstanceFamily:     ptrString("S5"),
					InstanceType:       ptrString("S5.SMALL2"),
					InstanceChargeType: ptrString("POSTPAID_BY_HOUR"),
					Status:             ptrString("SELL"),
					NetworkCard:        ptrInt64(25),
					TypeName:           ptrString("S5.SMALL2"),
					SoldOutReason:      ptrString(""),
					InstanceBandwidth:  ptrFloat64(1.5),
					InstancePps:        ptrInt64(30),
					StorageBlockAmount: ptrInt64(0),
					CpuType:            ptrString("Intel Xeon Cascade Lake"),
					Fpga:               ptrInt64(0),
					GpuCount:           ptrFloat64(0),
					Frequency:          ptrString("2.5GHz"),
					StatusCategory:     ptrString("EnoughStock"),
					Remark:             ptrString(""),
					Price: &cvm_sdk.ItemPrice{
						UnitPrice:         ptrFloat64(0.06),
						ChargeUnit:        ptrString("HOUR"),
						OriginalPrice:     ptrFloat64(0.06),
						DiscountPrice:     ptrFloat64(0.06),
						Discount:          ptrFloat64(100),
						UnitPriceDiscount: ptrFloat64(0.06),
					},
				},
			}, nil
		},
	)

	// Mock CbsService.DescribeDiskConfigQuota - capture both parameters
	var capturedInquiryType string
	var capturedInstanceFamilies []string
	patches.ApplyMethodFunc(&svccbs.CbsService{}, "DescribeDiskConfigQuota",
		func(ctx context.Context, cvmInfo map[string]interface{}, inquiryType string, instanceFamilies []string) ([]*cbs.DiskConfig, error) {
			capturedInquiryType = inquiryType
			capturedInstanceFamilies = instanceFamilies
			return []*cbs.DiskConfig{
				{
					Available:      ptrBool(true),
					DiskChargeType: ptrString("PREPAID"),
					Zone:           ptrString("ap-guangzhou-6"),
					InstanceFamily: ptrString("S5"),
					DiskType:       ptrString("CLOUD_SSD"),
					StepSize:       ptrUint64(10),
					DiskUsage:      ptrString("SYSTEM_DISK"),
					MinDiskSize:    ptrUint64(20),
					MaxDiskSize:    ptrUint64(500),
					DeviceClass:    ptrString(""),
				},
			}, nil
		},
	)

	res := svccvm.DataSourceTencentCloudInstanceTypes()
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

	err := res.Read(d, mockMeta)
	assert.NoError(t, err)

	// Verify backward compatibility: inquiry_type defaults to empty (CBS service defaults to INQUIRY_CVM_CONFIG)
	assert.Equal(t, "", capturedInquiryType)
	// Verify backward compatibility: instance_families defaults to empty (CBS service derives from cvmInfo["family"])
	assert.Equal(t, []string(nil), capturedInstanceFamilies)
}

// TestInstanceTypesBothNewParams verifies that both inquiry_type and instance_families work together
func TestInstanceTypesBothNewParams(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	mockMeta := newInstanceTypesMockMeta()

	// Mock CvmService.DescribeInstancesSellTypeByFilter
	patches.ApplyMethodFunc(&svccvm.CvmService{}, "DescribeInstancesSellTypeByFilter",
		func(ctx context.Context, filterMap map[string][]string) ([]*cvm_sdk.InstanceTypeQuotaItem, error) {
			return []*cvm_sdk.InstanceTypeQuotaItem{
				{
					Zone:               ptrString("ap-guangzhou-6"),
					Cpu:                ptrInt64(2),
					Gpu:                ptrInt64(0),
					Memory:             ptrInt64(2),
					InstanceFamily:     ptrString("S5"),
					InstanceType:       ptrString("S5.SMALL2"),
					InstanceChargeType: ptrString("POSTPAID_BY_HOUR"),
					Status:             ptrString("SELL"),
					NetworkCard:        ptrInt64(25),
					TypeName:           ptrString("S5.SMALL2"),
					SoldOutReason:      ptrString(""),
					InstanceBandwidth:  ptrFloat64(1.5),
					InstancePps:        ptrInt64(30),
					StorageBlockAmount: ptrInt64(0),
					CpuType:            ptrString("Intel Xeon Cascade Lake"),
					Fpga:               ptrInt64(0),
					GpuCount:           ptrFloat64(0),
					Frequency:          ptrString("2.5GHz"),
					StatusCategory:     ptrString("EnoughStock"),
					Remark:             ptrString(""),
					Price: &cvm_sdk.ItemPrice{
						UnitPrice:         ptrFloat64(0.06),
						ChargeUnit:        ptrString("HOUR"),
						OriginalPrice:     ptrFloat64(0.06),
						DiscountPrice:     ptrFloat64(0.06),
						Discount:          ptrFloat64(100),
						UnitPriceDiscount: ptrFloat64(0.06),
					},
				},
			}, nil
		},
	)

	// Mock CbsService.DescribeDiskConfigQuota
	var capturedInquiryType string
	var capturedInstanceFamilies []string
	patches.ApplyMethodFunc(&svccbs.CbsService{}, "DescribeDiskConfigQuota",
		func(ctx context.Context, cvmInfo map[string]interface{}, inquiryType string, instanceFamilies []string) ([]*cbs.DiskConfig, error) {
			capturedInquiryType = inquiryType
			capturedInstanceFamilies = instanceFamilies
			return []*cbs.DiskConfig{
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
					DeviceClass:    ptrString(""),
				},
			}, nil
		},
	)

	res := svccvm.DataSourceTencentCloudInstanceTypes()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"cpu_core_count":    2,
		"memory_size":       2,
		"exclude_sold_out":  false,
		"availability_zone": "ap-guangzhou-6",
		"cbs_filter": []interface{}{
			map[string]interface{}{
				"inquiry_type":      "INQUIRY_CBS_CONFIG",
				"instance_families": []interface{}{"S5", "SA2", "M5"},
				"disk_types":        []interface{}{"CLOUD_SSD"},
				"disk_charge_type":  "POSTPAID_BY_HOUR",
				"disk_usage":        "DATA_DISK",
			},
		},
	})

	err := res.Read(d, mockMeta)
	assert.NoError(t, err)

	// Verify both new parameters were correctly passed
	assert.Equal(t, "INQUIRY_CBS_CONFIG", capturedInquiryType)
	assert.Equal(t, []string{"S5", "SA2", "M5"}, capturedInstanceFamilies)
}
