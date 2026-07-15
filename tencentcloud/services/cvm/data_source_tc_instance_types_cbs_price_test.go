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

// instanceTypesCbsPriceMockMeta implements tccommon.ProviderMeta
type instanceTypesCbsPriceMockMeta struct {
	client *connectivity.TencentCloudClient
}

func (m *instanceTypesCbsPriceMockMeta) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &instanceTypesCbsPriceMockMeta{}

func newInstanceTypesCbsPriceMockMeta() *instanceTypesCbsPriceMockMeta {
	return &instanceTypesCbsPriceMockMeta{client: &connectivity.TencentCloudClient{}}
}

func ptrString(v string) *string    { return &v }
func ptrInt64(v int64) *int64       { return &v }
func ptrUint64(v uint64) *uint64    { return &v }
func ptrFloat64(v float64) *float64 { return &v }
func ptrBool(v bool) *bool          { return &v }

// go test ./tencentcloud/services/cvm/ -run "TestInstanceTypesCbsConfigPriceFields" -v -count=1 -gcflags="all=-l"

// TestInstanceTypesCbsConfigPriceFields_Read_Success tests that CBS config pricing fields are correctly mapped when Price is populated
func TestInstanceTypesCbsConfigPriceFields_Read_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	// Mock CVM SDK client
	cvmClient := &cvmSDK.Client{}
	patches.ApplyMethodReturn(newInstanceTypesCbsPriceMockMeta().client, "UseCvmClient", cvmClient)

	patches.ApplyMethodFunc(cvmClient, "DescribeZoneInstanceConfigInfos",
		func(request *cvmSDK.DescribeZoneInstanceConfigInfosRequest) (*cvmSDK.DescribeZoneInstanceConfigInfosResponse, error) {
			resp := cvmSDK.NewDescribeZoneInstanceConfigInfosResponse()
			resp.Response = &cvmSDK.DescribeZoneInstanceConfigInfosResponseParams{
				InstanceTypeQuotaSet: []*cvmSDK.InstanceTypeQuotaItem{
					{
						Zone:               ptrString("ap-guangzhou-6"),
						Cpu:                ptrInt64(2),
						Memory:             ptrInt64(2),
						InstanceFamily:     ptrString("S6"),
						InstanceType:       ptrString("S6.MEDIUM2"),
						InstanceChargeType: ptrString("POSTPAID_BY_HOUR"),
						Status:             ptrString("SELL"),
					},
				},
				RequestId: ptrString("fake-request-id"),
			}
			return resp, nil
		})

	// Mock CBS SDK client
	cbsClient := &cbsSDK.Client{}
	patches.ApplyMethodReturn(newInstanceTypesCbsPriceMockMeta().client, "UseCbsClient", cbsClient)

	patches.ApplyMethodFunc(cbsClient, "DescribeDiskConfigQuota",
		func(request *cbsSDK.DescribeDiskConfigQuotaRequest) (*cbsSDK.DescribeDiskConfigQuotaResponse, error) {
			resp := cbsSDK.NewDescribeDiskConfigQuotaResponse()
			resp.Response = &cbsSDK.DescribeDiskConfigQuotaResponseParams{
				DiskConfigSet: []*cbsSDK.DiskConfig{
					{
						Available:      ptrBool(true),
						DiskChargeType: ptrString("PREPAID"),
						Zone:           ptrString("ap-guangzhou-6"),
						InstanceFamily: ptrString("S6"),
						DiskType:       ptrString("CLOUD_SSD"),
						StepSize:       ptrUint64(10),
						DiskUsage:      ptrString("SYSTEM_DISK"),
						MinDiskSize:    ptrUint64(50),
						MaxDiskSize:    ptrUint64(500),
						Price: &cbsSDK.Price{
							ChargeUnit:            ptrString("HOUR"),
							DiscountPrice:         ptrFloat64(100.0),
							DiscountPriceHigh:     ptrString("100.00"),
							OriginalPrice:         ptrFloat64(200.0),
							OriginalPriceHigh:     ptrString("200.00"),
							UnitPrice:             ptrFloat64(0.5),
							UnitPriceDiscount:     ptrFloat64(0.3),
							UnitPriceDiscountHigh: ptrString("0.30"),
							UnitPriceHigh:         ptrString("0.50"),
						},
					},
				},
				RequestId: ptrString("fake-request-id"),
			}
			return resp, nil
		})

	meta := newInstanceTypesCbsPriceMockMeta()
	res := cvm.DataSourceTencentCloudInstanceTypes()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"cpu_core_count":   2,
		"memory_size":      2,
		"exclude_sold_out": false,
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

	// Verify the instance_types list was populated
	instanceTypes := d.Get("instance_types").([]interface{})
	assert.Equal(t, 1, len(instanceTypes))

	// Get the first instance type
	instanceType := instanceTypes[0].(map[string]interface{})
	cbsConfigs := instanceType["cbs_configs"].([]interface{})
	assert.Equal(t, 1, len(cbsConfigs))

	// Get the first CBS config and check pricing fields
	cbsConfig := cbsConfigs[0].(map[string]interface{})
	assert.Equal(t, "HOUR", cbsConfig["charge_unit"])
	assert.Equal(t, 100.0, cbsConfig["discount_price"])
	assert.Equal(t, "100.00", cbsConfig["discount_price_high"])
	assert.Equal(t, 200.0, cbsConfig["original_price"])
	assert.Equal(t, "200.00", cbsConfig["original_price_high"])
	assert.Equal(t, 0.5, cbsConfig["unit_price"])
	assert.Equal(t, 0.3, cbsConfig["unit_price_discount"])
	assert.Equal(t, "0.30", cbsConfig["unit_price_discount_high"])
	assert.Equal(t, "0.50", cbsConfig["unit_price_high"])
}

// TestInstanceTypesCbsConfigPriceFields_Read_PriceNil tests that CBS config pricing fields are null/empty when Price is nil
func TestInstanceTypesCbsConfigPriceFields_Read_PriceNil(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	// Mock CVM SDK client
	cvmClient := &cvmSDK.Client{}
	patches.ApplyMethodReturn(newInstanceTypesCbsPriceMockMeta().client, "UseCvmClient", cvmClient)

	patches.ApplyMethodFunc(cvmClient, "DescribeZoneInstanceConfigInfos",
		func(request *cvmSDK.DescribeZoneInstanceConfigInfosRequest) (*cvmSDK.DescribeZoneInstanceConfigInfosResponse, error) {
			resp := cvmSDK.NewDescribeZoneInstanceConfigInfosResponse()
			resp.Response = &cvmSDK.DescribeZoneInstanceConfigInfosResponseParams{
				InstanceTypeQuotaSet: []*cvmSDK.InstanceTypeQuotaItem{
					{
						Zone:               ptrString("ap-guangzhou-6"),
						Cpu:                ptrInt64(2),
						Memory:             ptrInt64(2),
						InstanceFamily:     ptrString("S6"),
						InstanceType:       ptrString("S6.MEDIUM2"),
						InstanceChargeType: ptrString("POSTPAID_BY_HOUR"),
						Status:             ptrString("SELL"),
					},
				},
				RequestId: ptrString("fake-request-id"),
			}
			return resp, nil
		})

	// Mock CBS SDK client
	cbsClient := &cbsSDK.Client{}
	patches.ApplyMethodReturn(newInstanceTypesCbsPriceMockMeta().client, "UseCbsClient", cbsClient)

	patches.ApplyMethodFunc(cbsClient, "DescribeDiskConfigQuota",
		func(request *cbsSDK.DescribeDiskConfigQuotaRequest) (*cbsSDK.DescribeDiskConfigQuotaResponse, error) {
			resp := cbsSDK.NewDescribeDiskConfigQuotaResponse()
			resp.Response = &cbsSDK.DescribeDiskConfigQuotaResponseParams{
				DiskConfigSet: []*cbsSDK.DiskConfig{
					{
						Available:      ptrBool(true),
						DiskChargeType: ptrString("PREPAID"),
						Zone:           ptrString("ap-guangzhou-6"),
						InstanceFamily: ptrString("S6"),
						DiskType:       ptrString("CLOUD_SSD"),
						StepSize:       ptrUint64(10),
						DiskUsage:      ptrString("SYSTEM_DISK"),
						MinDiskSize:    ptrUint64(50),
						MaxDiskSize:    ptrUint64(500),
						Price:          nil,
					},
				},
				RequestId: ptrString("fake-request-id"),
			}
			return resp, nil
		})

	meta := newInstanceTypesCbsPriceMockMeta()
	res := cvm.DataSourceTencentCloudInstanceTypes()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"cpu_core_count":   2,
		"memory_size":      2,
		"exclude_sold_out": false,
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

	// Verify the instance_types list was populated
	instanceTypes := d.Get("instance_types").([]interface{})
	assert.Equal(t, 1, len(instanceTypes))

	// Get the first instance type
	instanceType := instanceTypes[0].(map[string]interface{})
	cbsConfigs := instanceType["cbs_configs"].([]interface{})
	assert.Equal(t, 1, len(cbsConfigs))

	// Get the first CBS config and verify pricing fields are nil/zero when Price is nil
	cbsConfig := cbsConfigs[0].(map[string]interface{})

	// When Price is nil, the pricing fields should have nil/zero values
	assert.Equal(t, "", cbsConfig["charge_unit"], "charge_unit should be empty when Price is nil")
	assert.Equal(t, 0.0, cbsConfig["discount_price"], "discount_price should be 0 when Price is nil")
	assert.Equal(t, "", cbsConfig["discount_price_high"], "discount_price_high should be empty when Price is nil")
	assert.Equal(t, 0.0, cbsConfig["original_price"], "original_price should be 0 when Price is nil")
	assert.Equal(t, "", cbsConfig["original_price_high"], "original_price_high should be empty when Price is nil")
	assert.Equal(t, 0.0, cbsConfig["unit_price"], "unit_price should be 0 when Price is nil")
	assert.Equal(t, 0.0, cbsConfig["unit_price_discount"], "unit_price_discount should be 0 when Price is nil")
	assert.Equal(t, "", cbsConfig["unit_price_discount_high"], "unit_price_discount_high should be empty when Price is nil")
	assert.Equal(t, "", cbsConfig["unit_price_high"], "unit_price_high should be empty when Price is nil")
}

// TestInstanceTypesCbsConfigPriceFields_Schema validates that the new pricing fields exist in the cbs_configs schema
func TestInstanceTypesCbsConfigPriceFields_Schema(t *testing.T) {
	res := cvm.DataSourceTencentCloudInstanceTypes()

	assert.NotNil(t, res)
	assert.NotNil(t, res.Schema)

	// Check instance_types computed field exists
	instanceTypesSchema, ok := res.Schema["instance_types"]
	assert.True(t, ok, "instance_types should exist in schema")
	assert.True(t, instanceTypesSchema.Computed)

	// Get cbs_configs nested schema
	instanceTypesElem := instanceTypesSchema.Elem.(*schema.Resource)
	cbsConfigsSchema, ok := instanceTypesElem.Schema["cbs_configs"]
	assert.True(t, ok, "cbs_configs should exist in instance_types schema")
	assert.True(t, cbsConfigsSchema.Computed)

	// Get cbs_configs element schema
	cbsConfigsElem := cbsConfigsSchema.Elem.(*schema.Resource)
	cbsConfigsFields := cbsConfigsElem.Schema

	// Verify all 9 new pricing fields exist in cbs_configs schema
	expectedFields := map[string]schema.ValueType{
		"charge_unit":              schema.TypeString,
		"discount_price":           schema.TypeFloat,
		"discount_price_high":      schema.TypeString,
		"original_price":           schema.TypeFloat,
		"original_price_high":      schema.TypeString,
		"unit_price":               schema.TypeFloat,
		"unit_price_discount":      schema.TypeFloat,
		"unit_price_discount_high": schema.TypeString,
		"unit_price_high":          schema.TypeString,
	}

	for fieldName, expectedType := range expectedFields {
		field, exists := cbsConfigsFields[fieldName]
		assert.True(t, exists, "field %s should exist in cbs_configs schema", fieldName)
		if exists {
			assert.Equal(t, expectedType, field.Type, "field %s should have type %v", fieldName, expectedType)
			assert.True(t, field.Computed, "field %s should be computed", fieldName)
		}
	}
}
