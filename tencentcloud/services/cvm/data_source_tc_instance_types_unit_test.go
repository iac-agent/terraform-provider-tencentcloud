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
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
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

// go test ./tencentcloud/services/cvm/ -run "TestInstanceTypesCbsPricing" -v -count=1 -gcflags="all=-l"

// TestInstanceTypesCbsPricing_Read_WithPriceFields tests that pricing fields from DiskConfig.Price are correctly mapped
func TestInstanceTypesCbsPricing_Read_WithPriceFields(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	cvmClient := &cvmSDK.Client{}
	patches.ApplyMethodReturn(newInstanceTypesMockMeta().client, "UseCvmClient", cvmClient)

	patches.ApplyMethodFunc(cvmClient, "DescribeZoneInstanceConfigInfos", func(request *cvmSDK.DescribeZoneInstanceConfigInfosRequest) (*cvmSDK.DescribeZoneInstanceConfigInfosResponse, error) {
		resp := cvmSDK.NewDescribeZoneInstanceConfigInfosResponse()
		resp.Response = &cvmSDK.DescribeZoneInstanceConfigInfosResponseParams{
			InstanceTypeQuotaSet: []*cvmSDK.InstanceTypeQuotaItem{
				{
					Zone:               helper.String("ap-guangzhou-6"),
					InstanceFamily:     helper.String("S6"),
					InstanceType:       helper.String("S6.MEDIUM4"),
					Cpu:                helper.Int64(2),
					Memory:             helper.Int64(4),
					Status:             helper.String("Sell"),
					InstanceChargeType: helper.String("POSTPAID_BY_HOUR"),
				},
			},
			RequestId: helper.String("fake-request-id"),
		}
		return resp, nil
	})

	cbsClient := &cbsSDK.Client{}
	patches.ApplyMethodReturn(newInstanceTypesMockMeta().client, "UseCbsClient", cbsClient)

	patches.ApplyMethodFunc(cbsClient, "DescribeDiskConfigQuota", func(request *cbsSDK.DescribeDiskConfigQuotaRequest) (*cbsSDK.DescribeDiskConfigQuotaResponse, error) {
		resp := cbsSDK.NewDescribeDiskConfigQuotaResponse()
		resp.Response = &cbsSDK.DescribeDiskConfigQuotaResponseParams{
			DiskConfigSet: []*cbsSDK.DiskConfig{
				{
					Available:      helper.Bool(true),
					DiskChargeType: helper.String("POSTPAID_BY_HOUR"),
					Zone:           helper.String("ap-guangzhou-6"),
					InstanceFamily: helper.String("S6"),
					DiskType:       helper.String("CLOUD_SSD"),
					StepSize:       helper.Uint64(10),
					DiskUsage:      helper.String("DATA_DISK"),
					MinDiskSize:    helper.Uint64(20),
					MaxDiskSize:    helper.Uint64(500),
					Price: &cbsSDK.Price{
						ChargeUnit:            helper.String("HOUR"),
						UnitPrice:             helper.Float64(0.5),
						UnitPriceDiscount:     helper.Float64(0.3),
						UnitPriceHigh:         helper.String("0.5000"),
						UnitPriceDiscountHigh: helper.String("0.3000"),
						OriginalPrice:         helper.Float64(100.0),
						OriginalPriceHigh:     helper.String("100.0000"),
						DiscountPrice:         helper.Float64(80.0),
						DiscountPriceHigh:     helper.String("80.0000"),
					},
				},
			},
			RequestId: helper.String("fake-request-id"),
		}
		return resp, nil
	})

	meta := newInstanceTypesMockMeta()
	res := cvm.DataSourceTencentCloudInstanceTypes()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"cpu_core_count":   2,
		"memory_size":      4,
		"exclude_sold_out": true,
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
				"disk_usage":       "DATA_DISK",
			},
		},
	})

	err := res.Read(d, meta)
	assert.NoError(t, err)

	instanceTypes := d.Get("instance_types").([]interface{})
	assert.NotEmpty(t, instanceTypes)

	firstInstance := instanceTypes[0].(map[string]interface{})
	cbsConfigs := firstInstance["cbs_configs"].([]interface{})
	assert.NotEmpty(t, cbsConfigs)

	firstCbsConfig := cbsConfigs[0].(map[string]interface{})
	assert.Equal(t, "HOUR", firstCbsConfig["charge_unit"])
	assert.Equal(t, 0.5, firstCbsConfig["unit_price"])
	assert.Equal(t, 0.3, firstCbsConfig["unit_price_discount"])
	assert.Equal(t, "0.5000", firstCbsConfig["unit_price_high"])
	assert.Equal(t, "0.3000", firstCbsConfig["unit_price_discount_high"])
	assert.Equal(t, 100.0, firstCbsConfig["original_price"])
	assert.Equal(t, "100.0000", firstCbsConfig["original_price_high"])
	assert.Equal(t, 80.0, firstCbsConfig["discount_price"])
	assert.Equal(t, "80.0000", firstCbsConfig["discount_price_high"])
}

// TestInstanceTypesCbsPricing_Read_NilPrice tests that nil Price in DiskConfig is handled gracefully
func TestInstanceTypesCbsPricing_Read_NilPrice(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	cvmClient := &cvmSDK.Client{}
	patches.ApplyMethodReturn(newInstanceTypesMockMeta().client, "UseCvmClient", cvmClient)

	patches.ApplyMethodFunc(cvmClient, "DescribeZoneInstanceConfigInfos", func(request *cvmSDK.DescribeZoneInstanceConfigInfosRequest) (*cvmSDK.DescribeZoneInstanceConfigInfosResponse, error) {
		resp := cvmSDK.NewDescribeZoneInstanceConfigInfosResponse()
		resp.Response = &cvmSDK.DescribeZoneInstanceConfigInfosResponseParams{
			InstanceTypeQuotaSet: []*cvmSDK.InstanceTypeQuotaItem{
				{
					Zone:               helper.String("ap-guangzhou-6"),
					InstanceFamily:     helper.String("S6"),
					InstanceType:       helper.String("S6.MEDIUM4"),
					Cpu:                helper.Int64(2),
					Memory:             helper.Int64(4),
					Status:             helper.String("Sell"),
					InstanceChargeType: helper.String("POSTPAID_BY_HOUR"),
				},
			},
			RequestId: helper.String("fake-request-id"),
		}
		return resp, nil
	})

	cbsClient := &cbsSDK.Client{}
	patches.ApplyMethodReturn(newInstanceTypesMockMeta().client, "UseCbsClient", cbsClient)

	patches.ApplyMethodFunc(cbsClient, "DescribeDiskConfigQuota", func(request *cbsSDK.DescribeDiskConfigQuotaRequest) (*cbsSDK.DescribeDiskConfigQuotaResponse, error) {
		resp := cbsSDK.NewDescribeDiskConfigQuotaResponse()
		resp.Response = &cbsSDK.DescribeDiskConfigQuotaResponseParams{
			DiskConfigSet: []*cbsSDK.DiskConfig{
				{
					Available:      helper.Bool(true),
					DiskChargeType: helper.String("POSTPAID_BY_HOUR"),
					Zone:           helper.String("ap-guangzhou-6"),
					InstanceFamily: helper.String("S6"),
					DiskType:       helper.String("CLOUD_SSD"),
					StepSize:       helper.Uint64(10),
					DiskUsage:      helper.String("DATA_DISK"),
					MinDiskSize:    helper.Uint64(20),
					MaxDiskSize:    helper.Uint64(500),
					Price:          nil,
				},
			},
			RequestId: helper.String("fake-request-id"),
		}
		return resp, nil
	})

	meta := newInstanceTypesMockMeta()
	res := cvm.DataSourceTencentCloudInstanceTypes()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"cpu_core_count":   2,
		"memory_size":      4,
		"exclude_sold_out": true,
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
				"disk_usage":       "DATA_DISK",
			},
		},
	})

	err := res.Read(d, meta)
	assert.NoError(t, err)

	instanceTypes := d.Get("instance_types").([]interface{})
	assert.NotEmpty(t, instanceTypes)

	firstInstance := instanceTypes[0].(map[string]interface{})
	cbsConfigs := firstInstance["cbs_configs"].([]interface{})
	assert.NotEmpty(t, cbsConfigs)

	firstCbsConfig := cbsConfigs[0].(map[string]interface{})
	// When Price is nil, pricing fields should be zero-value (Terraform SDK returns "" for TypeString and 0 for TypeFloat)
	assert.Equal(t, "", firstCbsConfig["charge_unit"])
	assert.Equal(t, 0.0, firstCbsConfig["unit_price"])
	assert.Equal(t, 0.0, firstCbsConfig["unit_price_discount"])
	assert.Equal(t, "", firstCbsConfig["unit_price_high"])
	assert.Equal(t, "", firstCbsConfig["unit_price_discount_high"])
	assert.Equal(t, 0.0, firstCbsConfig["original_price"])
	assert.Equal(t, "", firstCbsConfig["original_price_high"])
	assert.Equal(t, 0.0, firstCbsConfig["discount_price"])
	assert.Equal(t, "", firstCbsConfig["discount_price_high"])

	// Non-pricing fields should still be present
	assert.Equal(t, "CLOUD_SSD", firstCbsConfig["disk_type"])
	assert.Equal(t, true, firstCbsConfig["available"])
}

// TestInstanceTypesCbsPricing_Schema validates the 9 new pricing fields exist in cbs_configs schema
func TestInstanceTypesCbsPricing_Schema(t *testing.T) {
	res := cvm.DataSourceTencentCloudInstanceTypes()

	assert.NotNil(t, res)
	assert.NotNil(t, res.Read)

	instanceTypesSchema := res.Schema["instance_types"]
	assert.NotNil(t, instanceTypesSchema)
	instanceTypesElem := instanceTypesSchema.Elem.(*schema.Resource)

	cbsConfigsSchema := instanceTypesElem.Schema["cbs_configs"]
	assert.NotNil(t, cbsConfigsSchema)
	cbsConfigsElem := cbsConfigsSchema.Elem.(*schema.Resource)

	// Verify all 9 pricing fields are present and computed
	pricingFields := []struct {
		name     string
		typeName schema.ValueType
	}{
		{"charge_unit", schema.TypeString},
		{"unit_price", schema.TypeFloat},
		{"unit_price_discount", schema.TypeFloat},
		{"unit_price_high", schema.TypeString},
		{"unit_price_discount_high", schema.TypeString},
		{"original_price", schema.TypeFloat},
		{"original_price_high", schema.TypeString},
		{"discount_price", schema.TypeFloat},
		{"discount_price_high", schema.TypeString},
	}

	for _, pf := range pricingFields {
		fieldSchema, exists := cbsConfigsElem.Schema[pf.name]
		assert.True(t, exists, "pricing field %s should exist in cbs_configs schema", pf.name)
		assert.Equal(t, pf.typeName, fieldSchema.Type, "pricing field %s should have type %v", pf.name, pf.typeName)
		assert.True(t, fieldSchema.Computed, "pricing field %s should be computed", pf.name)
	}
}
