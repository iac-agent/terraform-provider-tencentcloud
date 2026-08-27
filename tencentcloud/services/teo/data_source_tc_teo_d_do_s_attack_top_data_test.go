package teo_test

import (
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/teo"
)

// go test ./tencentcloud/services/teo/ -run "TestTeoDDoSAttackTopDataDataSource" -v -count=1 -gcflags="all=-l"

type mockMetaDDoSAttackTopData struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaDDoSAttackTopData) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaDDoSAttackTopData{}

func newMockMetaDDoSAttackTopData() *mockMetaDDoSAttackTopData {
	return &mockMetaDDoSAttackTopData{client: &connectivity.TencentCloudClient{}}
}

func ptrStringTopData(s string) *string {
	return &s
}

func ptrInt64TopData(n int64) *int64 {
	return &n
}

func ptrUint64TopData(n uint64) *uint64 {
	return &n
}

// TestTeoDDoSAttackTopDataDataSource_ReadSuccess tests successful read with top data
func TestTeoDDoSAttackTopDataDataSource_ReadSuccess(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	meta := newMockMetaDDoSAttackTopData()
	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(meta.client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeDDoSAttackTopDataWithContext",
		func(_ interface{}, request *teov20220901.DescribeDDoSAttackTopDataRequest) (*teov20220901.DescribeDDoSAttackTopDataResponse, error) {
			resp := teov20220901.NewDescribeDDoSAttackTopDataResponse()
			resp.Response = &teov20220901.DescribeDDoSAttackTopDataResponseParams{
				Data: []*teov20220901.TopEntry{
					{
						Key: ptrStringTopData("TCP"),
						Value: []*teov20220901.TopEntryValue{
							{
								Name:  ptrStringTopData("1.2.3.4"),
								Count: ptrInt64TopData(1000),
							},
							{
								Name:  ptrStringTopData("5.6.7.8"),
								Count: ptrInt64TopData(500),
							},
						},
					},
					{
						Key: ptrStringTopData("UDP"),
						Value: []*teov20220901.TopEntryValue{
							{
								Name:  ptrStringTopData("9.10.11.12"),
								Count: ptrInt64TopData(300),
							},
						},
					},
				},
				TotalCount: ptrUint64TopData(2),
				RequestId:  ptrStringTopData("fake-request-id"),
			}
			return resp, nil
		})

	res := teo.DataSourceTencentCloudTeoDDoSAttackTopData()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"start_time":  "2024-01-01T00:00:00Z",
		"end_time":    "2024-01-02T00:00:00Z",
		"metric_name": "ddos_attackFlux_protocol",
	})

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.NotEmpty(t, d.Id())

	data := d.Get("data").([]interface{})
	assert.Len(t, data, 2)

	// Check first entry
	entry1 := data[0].(map[string]interface{})
	assert.Equal(t, "TCP", entry1["key"])
	values1 := entry1["value"].([]interface{})
	assert.Len(t, values1, 2)
	val1 := values1[0].(map[string]interface{})
	assert.Equal(t, "1.2.3.4", val1["name"])
	assert.Equal(t, int(1000), val1["count"])

	// Check second entry
	entry2 := data[1].(map[string]interface{})
	assert.Equal(t, "UDP", entry2["key"])
	values2 := entry2["value"].([]interface{})
	assert.Len(t, values2, 1)
}

// TestTeoDDoSAttackTopDataDataSource_ReadSuccessWithOptionalParams tests successful read with all optional params
func TestTeoDDoSAttackTopDataDataSource_ReadSuccessWithOptionalParams(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	meta := newMockMetaDDoSAttackTopData()
	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(meta.client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeDDoSAttackTopDataWithContext",
		func(_ interface{}, request *teov20220901.DescribeDDoSAttackTopDataRequest) (*teov20220901.DescribeDDoSAttackTopDataResponse, error) {
			// Verify that optional params are set correctly
			assert.NotNil(t, request.ZoneIds)
			assert.Len(t, request.ZoneIds, 1)
			assert.Equal(t, "zone-2qtuhspy7cr6", *request.ZoneIds[0])
			assert.NotNil(t, request.PolicyIds)
			assert.Len(t, request.PolicyIds, 1)
			assert.Equal(t, int64(100), *request.PolicyIds[0])
			assert.NotNil(t, request.AttackType)
			assert.Equal(t, "flood", *request.AttackType)
			assert.NotNil(t, request.ProtocolType)
			assert.Equal(t, "tcp", *request.ProtocolType)
			assert.NotNil(t, request.Port)
			assert.Equal(t, int64(80), *request.Port)
			assert.NotNil(t, request.Area)
			assert.Equal(t, "overseas", *request.Area)
			assert.NotNil(t, request.Limit)
			assert.Equal(t, int64(100), *request.Limit)

			resp := teov20220901.NewDescribeDDoSAttackTopDataResponse()
			resp.Response = &teov20220901.DescribeDDoSAttackTopDataResponseParams{
				Data: []*teov20220901.TopEntry{
					{
						Key: ptrStringTopData("flood"),
						Value: []*teov20220901.TopEntryValue{
							{
								Name:  ptrStringTopData("US"),
								Count: ptrInt64TopData(500),
							},
						},
					},
				},
				TotalCount: ptrUint64TopData(1),
				RequestId:  ptrStringTopData("fake-request-id"),
			}
			return resp, nil
		})

	res := teo.DataSourceTencentCloudTeoDDoSAttackTopData()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"start_time":    "2024-01-01T00:00:00Z",
		"end_time":      "2024-01-02T00:00:00Z",
		"metric_name":   "ddos_attackNum_attackType",
		"zone_ids":      []interface{}{"zone-2qtuhspy7cr6"},
		"policy_ids":    []interface{}{100},
		"attack_type":   "flood",
		"protocol_type": "tcp",
		"port":          80,
		"area":          "overseas",
	})

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.NotEmpty(t, d.Id())

	data := d.Get("data").([]interface{})
	assert.Len(t, data, 1)
}

// TestTeoDDoSAttackTopDataDataSource_ReadEmptyData tests empty data response
func TestTeoDDoSAttackTopDataDataSource_ReadEmptyData(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	meta := newMockMetaDDoSAttackTopData()
	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(meta.client, "UseTeoV20220901Client", teoClient)

	callCount := 0
	patches.ApplyMethodFunc(teoClient, "DescribeDDoSAttackTopDataWithContext",
		func(_ interface{}, request *teov20220901.DescribeDDoSAttackTopDataRequest) (*teov20220901.DescribeDDoSAttackTopDataResponse, error) {
			callCount++
			resp := teov20220901.NewDescribeDDoSAttackTopDataResponse()
			// Return empty/nil data to trigger retry
			resp.Response = &teov20220901.DescribeDDoSAttackTopDataResponseParams{
				Data:       nil,
				TotalCount: ptrUint64TopData(0),
				RequestId:  ptrStringTopData("fake-request-id"),
			}
			return resp, nil
		})

	res := teo.DataSourceTencentCloudTeoDDoSAttackTopData()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"start_time":  "2024-01-01T00:00:00Z",
		"end_time":    "2024-01-02T00:00:00Z",
		"metric_name": "ddos_attackFlux_protocol",
	})

	err := res.Read(d, meta)
	// Should error because empty data triggers NonRetryableError (retry stops immediately)
	assert.Error(t, err)
	assert.Equal(t, 1, callCount)
}

// TestTeoDDoSAttackTopDataDataSource_ReadNilValue tests null fields in TopEntry
func TestTeoDDoSAttackTopDataDataSource_ReadNilValue(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	meta := newMockMetaDDoSAttackTopData()
	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(meta.client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeDDoSAttackTopDataWithContext",
		func(_ interface{}, request *teov20220901.DescribeDDoSAttackTopDataRequest) (*teov20220901.DescribeDDoSAttackTopDataResponse, error) {
			resp := teov20220901.NewDescribeDDoSAttackTopDataResponse()
			resp.Response = &teov20220901.DescribeDDoSAttackTopDataResponseParams{
				Data: []*teov20220901.TopEntry{
					{
						Key:   nil,
						Value: nil,
					},
				},
				TotalCount: ptrUint64TopData(1),
				RequestId:  ptrStringTopData("fake-request-id"),
			}
			return resp, nil
		})
	res := teo.DataSourceTencentCloudTeoDDoSAttackTopData()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"start_time":  "2024-01-01T00:00:00Z",
		"end_time":    "2024-01-02T00:00:00Z",
		"metric_name": "ddos_attackFlux_protocol",
	})

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.NotEmpty(t, d.Id())

	data := d.Get("data").([]interface{})
	assert.Len(t, data, 1)
	entry := data[0].(map[string]interface{})
	// key should be nil (not set when Key is nil)
	assert.Nil(t, entry["key"])
	// value should be empty list
	values := entry["value"].([]interface{})
	assert.Len(t, values, 0)
}

// TestTeoDDoSAttackTopDataDataSource_Schema tests schema definition
func TestTeoDDoSAttackTopDataDataSource_Schema(t *testing.T) {
	res := teo.DataSourceTencentCloudTeoDDoSAttackTopData()

	assert.NotNil(t, res)
	assert.NotNil(t, res.Read)

	// Required fields
	assert.Contains(t, res.Schema, "start_time")
	assert.Contains(t, res.Schema, "end_time")
	assert.Contains(t, res.Schema, "metric_name")

	// Optional fields
	assert.Contains(t, res.Schema, "zone_ids")
	assert.Contains(t, res.Schema, "policy_ids")
	assert.Contains(t, res.Schema, "attack_type")
	assert.Contains(t, res.Schema, "protocol_type")
	assert.Contains(t, res.Schema, "port")
	assert.Contains(t, res.Schema, "area")

	// Computed fields
	assert.Contains(t, res.Schema, "data")

	// Output
	assert.Contains(t, res.Schema, "result_output_file")

	// start_time is Required
	startTime := res.Schema["start_time"]
	assert.Equal(t, schema.TypeString, startTime.Type)
	assert.True(t, startTime.Required)

	// end_time is Required
	endTime := res.Schema["end_time"]
	assert.Equal(t, schema.TypeString, endTime.Type)
	assert.True(t, endTime.Required)

	// metric_name is Required
	metricName := res.Schema["metric_name"]
	assert.Equal(t, schema.TypeString, metricName.Type)
	assert.True(t, metricName.Required)

	// zone_ids is Optional
	zoneIds := res.Schema["zone_ids"]
	assert.Equal(t, schema.TypeSet, zoneIds.Type)
	assert.True(t, zoneIds.Optional)

	// policy_ids is Optional
	policyIds := res.Schema["policy_ids"]
	assert.Equal(t, schema.TypeSet, policyIds.Type)
	assert.True(t, policyIds.Optional)

	// data is Computed
	data := res.Schema["data"]
	assert.Equal(t, schema.TypeList, data.Type)
	assert.True(t, data.Computed)
}
