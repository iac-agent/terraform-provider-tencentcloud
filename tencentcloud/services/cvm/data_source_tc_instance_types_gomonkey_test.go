package cvm_test

import (
	"fmt"
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
func itPtrFloat64(v float64) *float64 { return &v }
func itPtrUint64(v uint64) *uint64    { return &v }
func itPtrBool(v bool) *bool          { return &v }

// go test ./tencentcloud/services/cvm/ -run "TestInstanceTypes" -v -count=1 -gcflags="all=-l"

// TestInstanceTypes_Read_WithInquiryType tests that inquiry_type is correctly passed to CBS API
func TestInstanceTypes_Read_WithInquiryType(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	cvmClient := &cvmSDK.Client{}
	cbsClient := &cbsSDK.Client{}

	patches.ApplyMethodReturn(newInstanceTypesMockMeta().client, "UseCvmClient", cvmClient)
	patches.ApplyMethodReturn(newInstanceTypesMockMeta().client, "UseCbsClient", cbsClient)

	// Mock CVM DescribeZoneInstanceConfigInfos
	patches.ApplyMethodFunc(cvmClient, "DescribeZoneInstanceConfigInfos",
		func(request *cvmSDK.DescribeZoneInstanceConfigInfosRequest) (*cvmSDK.DescribeZoneInstanceConfigInfosResponse, error) {
			resp := cvmSDK.NewDescribeZoneInstanceConfigInfosResponse()
			resp.Response = &cvmSDK.DescribeZoneInstanceConfigInfosResponseParams{
				InstanceTypeQuotaSet: []*cvmSDK.InstanceTypeQuotaItem{
					{
						Zone:           itPtrString("ap-guangzhou-6"),
						InstanceType:   itPtrString("S6.MEDIUM4"),
						Cpu:            itPtrInt64(2),
						Memory:         itPtrInt64(4),
						InstanceFamily: itPtrString("S6"),
						Status:         itPtrString("Sell"),
					},
				},
				RequestId: itPtrString("fake-cvm-request-id"),
			}
			return resp, nil
		},
	)

	// Mock CBS DescribeDiskConfigQuota - verify InquiryType is set correctly
	patches.ApplyMethodFunc(cbsClient, "DescribeDiskConfigQuota",
		func(request *cbsSDK.DescribeDiskConfigQuotaRequest) (*cbsSDK.DescribeDiskConfigQuotaResponse, error) {
			// Verify InquiryType is set to INQUIRY_CBS_CONFIG
			assert.NotNil(t, request.InquiryType)
			assert.Equal(t, "INQUIRY_CBS_CONFIG", *request.InquiryType)

			resp := cbsSDK.NewDescribeDiskConfigQuotaResponse()
			resp.Response = &cbsSDK.DescribeDiskConfigQuotaResponseParams{
				DiskConfigSet: []*cbsSDK.DiskConfig{
					{
						Available:      itPtrBool(true),
						DiskChargeType: itPtrString("PREPAID"),
						Zone:           itPtrString("ap-guangzhou-6"),
						InstanceFamily: itPtrString("S6"),
						DiskType:       itPtrString("CLOUD_SSD"),
						DiskUsage:      itPtrString("SYSTEM_DISK"),
						MinDiskSize:    itPtrUint64(50),
						MaxDiskSize:    itPtrUint64(500),
					},
				},
				RequestId: itPtrString("fake-cbs-request-id"),
			}
			return resp, nil
		},
	)

	meta := newInstanceTypesMockMeta()
	res := cvm.DataSourceTencentCloudInstanceTypes()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"availability_zone": "ap-guangzhou-6",
		"cpu_core_count":    2,
		"memory_size":       4,
		"inquiry_type":      "INQUIRY_CBS_CONFIG",
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
}

// TestInstanceTypes_Read_DefaultInquiryType tests that inquiry_type defaults to INQUIRY_CVM_CONFIG
func TestInstanceTypes_Read_DefaultInquiryType(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	cvmClient := &cvmSDK.Client{}
	cbsClient := &cbsSDK.Client{}

	patches.ApplyMethodReturn(newInstanceTypesMockMeta().client, "UseCvmClient", cvmClient)
	patches.ApplyMethodReturn(newInstanceTypesMockMeta().client, "UseCbsClient", cbsClient)

	patches.ApplyMethodFunc(cvmClient, "DescribeZoneInstanceConfigInfos",
		func(request *cvmSDK.DescribeZoneInstanceConfigInfosRequest) (*cvmSDK.DescribeZoneInstanceConfigInfosResponse, error) {
			resp := cvmSDK.NewDescribeZoneInstanceConfigInfosResponse()
			resp.Response = &cvmSDK.DescribeZoneInstanceConfigInfosResponseParams{
				InstanceTypeQuotaSet: []*cvmSDK.InstanceTypeQuotaItem{
					{
						Zone:           itPtrString("ap-guangzhou-6"),
						InstanceType:   itPtrString("S6.MEDIUM4"),
						Cpu:            itPtrInt64(2),
						Memory:         itPtrInt64(4),
						InstanceFamily: itPtrString("S6"),
						Status:         itPtrString("Sell"),
					},
				},
				RequestId: itPtrString("fake-cvm-request-id"),
			}
			return resp, nil
		},
	)

	patches.ApplyMethodFunc(cbsClient, "DescribeDiskConfigQuota",
		func(request *cbsSDK.DescribeDiskConfigQuotaRequest) (*cbsSDK.DescribeDiskConfigQuotaResponse, error) {
			// Verify InquiryType defaults to INQUIRY_CVM_CONFIG
			assert.NotNil(t, request.InquiryType)
			assert.Equal(t, "INQUIRY_CVM_CONFIG", *request.InquiryType)

			resp := cbsSDK.NewDescribeDiskConfigQuotaResponse()
			resp.Response = &cbsSDK.DescribeDiskConfigQuotaResponseParams{
				DiskConfigSet: []*cbsSDK.DiskConfig{},
				RequestId:     itPtrString("fake-cbs-request-id"),
			}
			return resp, nil
		},
	)

	meta := newInstanceTypesMockMeta()
	res := cvm.DataSourceTencentCloudInstanceTypes()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"availability_zone": "ap-guangzhou-6",
		"cpu_core_count":    2,
		"memory_size":       4,
		"cbs_filter": []interface{}{
			map[string]interface{}{
				"disk_types": []interface{}{"CLOUD_SSD"},
				"disk_usage": "SYSTEM_DISK",
			},
		},
	})

	err := res.Read(d, meta)
	assert.NoError(t, err)
}

// TestInstanceTypes_Read_TopLevelDiskChargeTypePrecedence tests top-level disk_charge_type takes precedence
func TestInstanceTypes_Read_TopLevelDiskChargeTypePrecedence(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	cvmClient := &cvmSDK.Client{}
	cbsClient := &cbsSDK.Client{}

	patches.ApplyMethodReturn(newInstanceTypesMockMeta().client, "UseCvmClient", cvmClient)
	patches.ApplyMethodReturn(newInstanceTypesMockMeta().client, "UseCbsClient", cbsClient)

	patches.ApplyMethodFunc(cvmClient, "DescribeZoneInstanceConfigInfos",
		func(request *cvmSDK.DescribeZoneInstanceConfigInfosRequest) (*cvmSDK.DescribeZoneInstanceConfigInfosResponse, error) {
			resp := cvmSDK.NewDescribeZoneInstanceConfigInfosResponse()
			resp.Response = &cvmSDK.DescribeZoneInstanceConfigInfosResponseParams{
				InstanceTypeQuotaSet: []*cvmSDK.InstanceTypeQuotaItem{
					{
						Zone:           itPtrString("ap-guangzhou-6"),
						InstanceType:   itPtrString("S6.MEDIUM4"),
						Cpu:            itPtrInt64(2),
						Memory:         itPtrInt64(4),
						InstanceFamily: itPtrString("S6"),
						Status:         itPtrString("Sell"),
					},
				},
				RequestId: itPtrString("fake-cvm-request-id"),
			}
			return resp, nil
		},
	)

	patches.ApplyMethodFunc(cbsClient, "DescribeDiskConfigQuota",
		func(request *cbsSDK.DescribeDiskConfigQuotaRequest) (*cbsSDK.DescribeDiskConfigQuotaResponse, error) {
			// Verify top-level disk_charge_type (POSTPAID_BY_HOUR) takes precedence over cbs_filter.disk_charge_type (PREPAID)
			assert.NotNil(t, request.DiskChargeType)
			assert.Equal(t, "POSTPAID_BY_HOUR", *request.DiskChargeType)

			resp := cbsSDK.NewDescribeDiskConfigQuotaResponse()
			resp.Response = &cbsSDK.DescribeDiskConfigQuotaResponseParams{
				DiskConfigSet: []*cbsSDK.DiskConfig{},
				RequestId:     itPtrString("fake-cbs-request-id"),
			}
			return resp, nil
		},
	)

	meta := newInstanceTypesMockMeta()
	res := cvm.DataSourceTencentCloudInstanceTypes()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"availability_zone": "ap-guangzhou-6",
		"cpu_core_count":    2,
		"memory_size":       4,
		"disk_charge_type":  "POSTPAID_BY_HOUR",
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
}

// TestInstanceTypes_Read_CbsFilterDiskChargeTypeBackwardCompat tests cbs_filter.disk_charge_type works when no top-level param
func TestInstanceTypes_Read_CbsFilterDiskChargeTypeBackwardCompat(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	cvmClient := &cvmSDK.Client{}
	cbsClient := &cbsSDK.Client{}

	patches.ApplyMethodReturn(newInstanceTypesMockMeta().client, "UseCvmClient", cvmClient)
	patches.ApplyMethodReturn(newInstanceTypesMockMeta().client, "UseCbsClient", cbsClient)

	patches.ApplyMethodFunc(cvmClient, "DescribeZoneInstanceConfigInfos",
		func(request *cvmSDK.DescribeZoneInstanceConfigInfosRequest) (*cvmSDK.DescribeZoneInstanceConfigInfosResponse, error) {
			resp := cvmSDK.NewDescribeZoneInstanceConfigInfosResponse()
			resp.Response = &cvmSDK.DescribeZoneInstanceConfigInfosResponseParams{
				InstanceTypeQuotaSet: []*cvmSDK.InstanceTypeQuotaItem{
					{
						Zone:           itPtrString("ap-guangzhou-6"),
						InstanceType:   itPtrString("S6.MEDIUM4"),
						Cpu:            itPtrInt64(2),
						Memory:         itPtrInt64(4),
						InstanceFamily: itPtrString("S6"),
						Status:         itPtrString("Sell"),
					},
				},
				RequestId: itPtrString("fake-cvm-request-id"),
			}
			return resp, nil
		},
	)

	patches.ApplyMethodFunc(cbsClient, "DescribeDiskConfigQuota",
		func(request *cbsSDK.DescribeDiskConfigQuotaRequest) (*cbsSDK.DescribeDiskConfigQuotaResponse, error) {
			// Verify cbs_filter.disk_charge_type is used when no top-level disk_charge_type
			assert.NotNil(t, request.DiskChargeType)
			assert.Equal(t, "PREPAID", *request.DiskChargeType)

			resp := cbsSDK.NewDescribeDiskConfigQuotaResponse()
			resp.Response = &cbsSDK.DescribeDiskConfigQuotaResponseParams{
				DiskConfigSet: []*cbsSDK.DiskConfig{},
				RequestId:     itPtrString("fake-cbs-request-id"),
			}
			return resp, nil
		},
	)

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
}

// TestInstanceTypes_Read_NoDiskChargeType tests that DiskChargeType is not set when neither is specified
func TestInstanceTypes_Read_NoDiskChargeType(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	cvmClient := &cvmSDK.Client{}
	cbsClient := &cbsSDK.Client{}

	patches.ApplyMethodReturn(newInstanceTypesMockMeta().client, "UseCvmClient", cvmClient)
	patches.ApplyMethodReturn(newInstanceTypesMockMeta().client, "UseCbsClient", cbsClient)

	patches.ApplyMethodFunc(cvmClient, "DescribeZoneInstanceConfigInfos",
		func(request *cvmSDK.DescribeZoneInstanceConfigInfosRequest) (*cvmSDK.DescribeZoneInstanceConfigInfosResponse, error) {
			resp := cvmSDK.NewDescribeZoneInstanceConfigInfosResponse()
			resp.Response = &cvmSDK.DescribeZoneInstanceConfigInfosResponseParams{
				InstanceTypeQuotaSet: []*cvmSDK.InstanceTypeQuotaItem{
					{
						Zone:           itPtrString("ap-guangzhou-6"),
						InstanceType:   itPtrString("S6.MEDIUM4"),
						Cpu:            itPtrInt64(2),
						Memory:         itPtrInt64(4),
						InstanceFamily: itPtrString("S6"),
						Status:         itPtrString("Sell"),
					},
				},
				RequestId: itPtrString("fake-cvm-request-id"),
			}
			return resp, nil
		},
	)

	patches.ApplyMethodFunc(cbsClient, "DescribeDiskConfigQuota",
		func(request *cbsSDK.DescribeDiskConfigQuotaRequest) (*cbsSDK.DescribeDiskConfigQuotaResponse, error) {
			// Verify DiskChargeType is nil when neither top-level nor cbs_filter disk_charge_type is set
			assert.Nil(t, request.DiskChargeType)

			resp := cbsSDK.NewDescribeDiskConfigQuotaResponse()
			resp.Response = &cbsSDK.DescribeDiskConfigQuotaResponseParams{
				DiskConfigSet: []*cbsSDK.DiskConfig{},
				RequestId:     itPtrString("fake-cbs-request-id"),
			}
			return resp, nil
		},
	)

	meta := newInstanceTypesMockMeta()
	res := cvm.DataSourceTencentCloudInstanceTypes()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"availability_zone": "ap-guangzhou-6",
		"cpu_core_count":    2,
		"memory_size":       4,
		"cbs_filter": []interface{}{
			map[string]interface{}{
				"disk_types": []interface{}{"CLOUD_SSD"},
				"disk_usage": "SYSTEM_DISK",
			},
		},
	})

	err := res.Read(d, meta)
	assert.NoError(t, err)
}

// TestInstanceTypes_Read_APIError tests Read handles API error
func TestInstanceTypes_Read_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	cvmClient := &cvmSDK.Client{}
	patches.ApplyMethodReturn(newInstanceTypesMockMeta().client, "UseCvmClient", cvmClient)

	patches.ApplyMethodFunc(cvmClient, "DescribeZoneInstanceConfigInfos",
		func(request *cvmSDK.DescribeZoneInstanceConfigInfosRequest) (*cvmSDK.DescribeZoneInstanceConfigInfosResponse, error) {
			return nil, fmt.Errorf("[TencentCloudSDKError] Code=InvalidParameter, Message=Invalid parameter")
		},
	)

	meta := newInstanceTypesMockMeta()
	res := cvm.DataSourceTencentCloudInstanceTypes()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"availability_zone": "ap-guangzhou-6",
		"cpu_core_count":    2,
		"memory_size":       4,
	})

	err := res.Read(d, meta)
	assert.Error(t, err)
}

// TestInstanceTypes_Schema validates schema definition for new parameters
func TestInstanceTypes_Schema(t *testing.T) {
	res := cvm.DataSourceTencentCloudInstanceTypes()

	assert.NotNil(t, res)
	assert.NotNil(t, res.Read)

	// Check inquiry_type
	assert.Contains(t, res.Schema, "inquiry_type")
	inquiryType := res.Schema["inquiry_type"]
	assert.Equal(t, schema.TypeString, inquiryType.Type)
	assert.True(t, inquiryType.Optional)
	assert.Equal(t, "INQUIRY_CVM_CONFIG", inquiryType.Default)

	// Check disk_charge_type
	assert.Contains(t, res.Schema, "disk_charge_type")
	diskChargeType := res.Schema["disk_charge_type"]
	assert.Equal(t, schema.TypeString, diskChargeType.Type)
	assert.True(t, diskChargeType.Optional)
	assert.Nil(t, diskChargeType.Default)
}
