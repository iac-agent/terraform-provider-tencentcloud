package teo_test

import (
	"fmt"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/teo"
)

// go test ./tencentcloud/services/teo/ -run "TestRealtimeLogDelivery" -v -count=1 -gcflags="all=-l"

// TestRealtimeLogDelivery_Read_WithFilters tests Read with filters specified
func TestRealtimeLogDelivery_Read_WithFilters(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeRealtimeLogDeliveryTasks", func(request *teov20220901.DescribeRealtimeLogDeliveryTasksRequest) (*teov20220901.DescribeRealtimeLogDeliveryTasksResponse, error) {
		resp := teov20220901.NewDescribeRealtimeLogDeliveryTasksResponse()
		resp.Response = &teov20220901.DescribeRealtimeLogDeliveryTasksResponseParams{
			TotalCount: ptrUint64(1),
			RealtimeLogDeliveryTasks: []*teov20220901.RealtimeLogDeliveryTask{
				{
					TaskId:         ptrString("task-abc123"),
					TaskName:       ptrString("test-task"),
					DeliveryStatus: ptrString("enabled"),
					TaskType:       ptrString("cls"),
					EntityList:     []*string{ptrString("domain.example.com")},
					LogType:        ptrString("domain"),
					Area:           ptrString("mainland"),
					Fields:         []*string{ptrString("Field1")},
					Sample:         ptrUint64(1000),
					CreateTime:     ptrString("2024-01-01T00:00:00Z"),
					UpdateTime:     ptrString("2024-01-02T00:00:00Z"),
				},
			},
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoRealtimeLogDelivery()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":     "zone-1234567890",
		"task_name":   "test-task",
		"task_type":   "cls",
		"entity_list": []interface{}{"domain.example.com"},
		"log_type":    "domain",
		"area":        "mainland",
		"fields":      []interface{}{"Field1"},
		"sample":      1000,
		"filters": []interface{}{
			map[string]interface{}{
				"name":   "task-type",
				"values": []interface{}{"cls"},
			},
		},
	})
	d.SetId("zone-1234567890#task-abc123")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "task-abc123", d.Get("task_id"))
	assert.Equal(t, "test-task", d.Get("task_name"))
	assert.Equal(t, "enabled", d.Get("delivery_status"))
	assert.Equal(t, "cls", d.Get("task_type"))

	tasks := d.Get("realtime_log_delivery_tasks").([]interface{})
	assert.Equal(t, 1, len(tasks))
	taskMap := tasks[0].(map[string]interface{})
	assert.Equal(t, "task-abc123", taskMap["task_id"])
	assert.Equal(t, "test-task", taskMap["task_name"])
	assert.Equal(t, "enabled", taskMap["delivery_status"])
	assert.Equal(t, "cls", taskMap["task_type"])
	assert.Equal(t, "domain", taskMap["log_type"])
	assert.Equal(t, "mainland", taskMap["area"])
	assert.Equal(t, "2024-01-01T00:00:00Z", taskMap["create_time"])
	assert.Equal(t, "2024-01-02T00:00:00Z", taskMap["update_time"])
}

// TestRealtimeLogDelivery_Read_WithoutFilters tests Read without filters (default behavior)
func TestRealtimeLogDelivery_Read_WithoutFilters(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeRealtimeLogDeliveryTasks", func(request *teov20220901.DescribeRealtimeLogDeliveryTasksRequest) (*teov20220901.DescribeRealtimeLogDeliveryTasksResponse, error) {
		resp := teov20220901.NewDescribeRealtimeLogDeliveryTasksResponse()
		resp.Response = &teov20220901.DescribeRealtimeLogDeliveryTasksResponseParams{
			TotalCount: ptrUint64(1),
			RealtimeLogDeliveryTasks: []*teov20220901.RealtimeLogDeliveryTask{
				{
					TaskId:         ptrString("task-abc123"),
					TaskName:       ptrString("test-task"),
					DeliveryStatus: ptrString("enabled"),
					TaskType:       ptrString("cls"),
					EntityList:     []*string{ptrString("domain.example.com")},
					LogType:        ptrString("domain"),
					Area:           ptrString("mainland"),
					Fields:         []*string{ptrString("Field1")},
					Sample:         ptrUint64(1000),
				},
			},
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoRealtimeLogDelivery()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":     "zone-1234567890",
		"task_name":   "test-task",
		"task_type":   "cls",
		"entity_list": []interface{}{"domain.example.com"},
		"log_type":    "domain",
		"area":        "mainland",
		"fields":      []interface{}{"Field1"},
		"sample":      1000,
	})
	d.SetId("zone-1234567890#task-abc123")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "task-abc123", d.Get("task_id"))
	assert.Equal(t, "test-task", d.Get("task_name"))
	assert.Equal(t, "enabled", d.Get("delivery_status"))
}

// TestRealtimeLogDelivery_Read_NotFound tests Read when resource is not found
func TestRealtimeLogDelivery_Read_NotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeRealtimeLogDeliveryTasks", func(request *teov20220901.DescribeRealtimeLogDeliveryTasksRequest) (*teov20220901.DescribeRealtimeLogDeliveryTasksResponse, error) {
		resp := teov20220901.NewDescribeRealtimeLogDeliveryTasksResponse()
		resp.Response = &teov20220901.DescribeRealtimeLogDeliveryTasksResponseParams{
			TotalCount:               ptrUint64(0),
			RealtimeLogDeliveryTasks: []*teov20220901.RealtimeLogDeliveryTask{},
			RequestId:                ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoRealtimeLogDelivery()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":     "zone-1234567890",
		"task_name":   "test-task",
		"task_type":   "cls",
		"entity_list": []interface{}{"domain.example.com"},
		"log_type":    "domain",
		"area":        "mainland",
		"fields":      []interface{}{"Field1"},
		"sample":      1000,
	})
	d.SetId("zone-1234567890#task-abc123")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}

// TestRealtimeLogDelivery_Read_APIError tests Read handles API error
func TestRealtimeLogDelivery_Read_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeRealtimeLogDeliveryTasks", func(request *teov20220901.DescribeRealtimeLogDeliveryTasksRequest) (*teov20220901.DescribeRealtimeLogDeliveryTasksResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=InternalError, Message=Internal error")
	})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoRealtimeLogDelivery()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":     "zone-1234567890",
		"task_name":   "test-task",
		"task_type":   "cls",
		"entity_list": []interface{}{"domain.example.com"},
		"log_type":    "domain",
		"area":        "mainland",
		"fields":      []interface{}{"Field1"},
		"sample":      1000,
	})
	d.SetId("zone-1234567890#task-abc123")

	err := res.Read(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "InternalError")
}

// TestRealtimeLogDelivery_Read_WithFiltersAndPagination tests Read with filters and pagination
func TestRealtimeLogDelivery_Read_WithFiltersAndPagination(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoClient", teoClient)

	callCount := 0
	patches.ApplyMethodFunc(teoClient, "DescribeRealtimeLogDeliveryTasks", func(request *teov20220901.DescribeRealtimeLogDeliveryTasksRequest) (*teov20220901.DescribeRealtimeLogDeliveryTasksResponse, error) {
		callCount++
		resp := teov20220901.NewDescribeRealtimeLogDeliveryTasksResponse()
		if callCount == 1 {
			// First page: return 2 tasks (one matching the resource ID), TotalCount=3 to trigger pagination
			resp.Response = &teov20220901.DescribeRealtimeLogDeliveryTasksResponseParams{
				TotalCount: ptrUint64(3),
				RealtimeLogDeliveryTasks: []*teov20220901.RealtimeLogDeliveryTask{
					{
						TaskId:   ptrString("task-abc123"),
						TaskName: ptrString("task-one"),
						TaskType: ptrString("cls"),
					},
					{
						TaskId:   ptrString("task-002"),
						TaskName: ptrString("task-two"),
						TaskType: ptrString("cls"),
					},
				},
				RequestId: ptrString("fake-request-id"),
			}
		} else {
			// Second page: return 1 more task
			resp.Response = &teov20220901.DescribeRealtimeLogDeliveryTasksResponseParams{
				TotalCount: ptrUint64(3),
				RealtimeLogDeliveryTasks: []*teov20220901.RealtimeLogDeliveryTask{
					{
						TaskId:   ptrString("task-003"),
						TaskName: ptrString("task-three"),
						TaskType: ptrString("s3"),
					},
				},
				RequestId: ptrString("fake-request-id-2"),
			}
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoRealtimeLogDelivery()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":     "zone-1234567890",
		"task_name":   "test-task",
		"task_type":   "cls",
		"entity_list": []interface{}{"domain.example.com"},
		"log_type":    "domain",
		"area":        "mainland",
		"fields":      []interface{}{"Field1"},
		"sample":      1000,
		"filters": []interface{}{
			map[string]interface{}{
				"name":   "task-type",
				"values": []interface{}{"cls"},
			},
		},
	})
	d.SetId("zone-1234567890#task-abc123")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, 2, callCount)

	tasks := d.Get("realtime_log_delivery_tasks").([]interface{})
	assert.Equal(t, 3, len(tasks))
}

// TestRealtimeLogDelivery_Schema validates schema definition for new parameters
func TestRealtimeLogDelivery_Schema(t *testing.T) {
	res := teo.ResourceTencentCloudTeoRealtimeLogDelivery()

	assert.NotNil(t, res)
	assert.NotNil(t, res.Create)
	assert.NotNil(t, res.Read)
	assert.NotNil(t, res.Update)
	assert.NotNil(t, res.Delete)
	assert.NotNil(t, res.Importer)

	// Check filters parameter
	assert.Contains(t, res.Schema, "filters")
	filters := res.Schema["filters"]
	assert.Equal(t, schema.TypeList, filters.Type)
	assert.True(t, filters.Optional)
	assert.False(t, filters.Computed)
	assert.NotNil(t, filters.Elem)

	// Check filters nested schema
	filtersRes := filters.Elem.(*schema.Resource)
	assert.Contains(t, filtersRes.Schema, "name")
	nameField := filtersRes.Schema["name"]
	assert.Equal(t, schema.TypeString, nameField.Type)
	assert.True(t, nameField.Required)

	assert.Contains(t, filtersRes.Schema, "values")
	valuesField := filtersRes.Schema["values"]
	assert.Equal(t, schema.TypeList, valuesField.Type)
	assert.True(t, valuesField.Required)

	assert.Contains(t, filtersRes.Schema, "fuzzy")
	fuzzyField := filtersRes.Schema["fuzzy"]
	assert.Equal(t, schema.TypeBool, fuzzyField.Type)
	assert.True(t, fuzzyField.Optional)

	// Check realtime_log_delivery_tasks parameter
	assert.Contains(t, res.Schema, "realtime_log_delivery_tasks")
	tasksField := res.Schema["realtime_log_delivery_tasks"]
	assert.Equal(t, schema.TypeList, tasksField.Type)
	assert.True(t, tasksField.Computed)
	assert.NotNil(t, tasksField.Elem)

	// Check realtime_log_delivery_tasks nested schema
	tasksRes := tasksField.Elem.(*schema.Resource)
	assert.Contains(t, tasksRes.Schema, "task_id")
	assert.True(t, tasksRes.Schema["task_id"].Computed)
	assert.Contains(t, tasksRes.Schema, "task_name")
	assert.True(t, tasksRes.Schema["task_name"].Computed)
	assert.Contains(t, tasksRes.Schema, "delivery_status")
	assert.True(t, tasksRes.Schema["delivery_status"].Computed)
	assert.Contains(t, tasksRes.Schema, "task_type")
	assert.True(t, tasksRes.Schema["task_type"].Computed)
	assert.Contains(t, tasksRes.Schema, "entity_list")
	assert.Contains(t, tasksRes.Schema, "log_type")
	assert.Contains(t, tasksRes.Schema, "area")
	assert.Contains(t, tasksRes.Schema, "fields")
	assert.Contains(t, tasksRes.Schema, "custom_fields")
	assert.Contains(t, tasksRes.Schema, "delivery_conditions")
	assert.Contains(t, tasksRes.Schema, "sample")
	assert.Contains(t, tasksRes.Schema, "log_format")
	assert.Contains(t, tasksRes.Schema, "cls")
	assert.Contains(t, tasksRes.Schema, "custom_endpoint")
	assert.Contains(t, tasksRes.Schema, "s3")
	assert.Contains(t, tasksRes.Schema, "create_time")
	assert.True(t, tasksRes.Schema["create_time"].Computed)
	assert.Contains(t, tasksRes.Schema, "update_time")
	assert.True(t, tasksRes.Schema["update_time"].Computed)
}
