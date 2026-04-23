package teo_test

import (
	"fmt"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/teo"
)

func TestAccTencentCloudTeoFunctionResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTeoFunction,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_teo_function.teo_function", "id"),
					resource.TestCheckResourceAttr("tencentcloud_teo_function.teo_function", "name", "aaa-zone-2qtuhspy7cr6-1310708577"),
					resource.TestCheckResourceAttr("tencentcloud_teo_function.teo_function", "remark", "test"),
					resource.TestCheckResourceAttr("tencentcloud_teo_function.teo_function", "content", `addEventListener('fetch', e => {
  const response = new Response('Hello World!!');
  e.respondWith(response);
});
`),
				),
			},
			{
				ResourceName:      "tencentcloud_teo_function.teo_function",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccTeoFunctionUp,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_teo_function.teo_function", "id"),
					resource.TestCheckResourceAttr("tencentcloud_teo_function.teo_function", "name", "aaa-zone-2qtuhspy7cr6-1310708577"),
					resource.TestCheckResourceAttr("tencentcloud_teo_function.teo_function", "remark", "test-update"),
					resource.TestCheckResourceAttr("tencentcloud_teo_function.teo_function", "content", `addEventListener('fetch', e => {
  const response = new Response('Hello World');
  e.respondWith(response);
});
`),
				),
			},
		},
	})
}

const testAccTeoFunction = `

resource "tencentcloud_teo_function" "teo_function" {
    content     = <<-EOT
        addEventListener('fetch', e => {
          const response = new Response('Hello World!!');
          e.respondWith(response);
        });
    EOT
    name        = "aaa-zone-2qtuhspy7cr6-1310708577"
    remark      = "test"
    zone_id     = "zone-2qtuhspy7cr6"
}
`
const testAccTeoFunctionUp = `

resource "tencentcloud_teo_function" "teo_function" {
    content     = <<-EOT
        addEventListener('fetch', e => {
          const response = new Response('Hello World');
          e.respondWith(response);
        });
    EOT
    name        = "aaa-zone-2qtuhspy7cr6-1310708577"
    remark      = "test-update"
    zone_id     = "zone-2qtuhspy7cr6"
}
`

// go test ./tencentcloud/services/teo/ -run "TestTeoFunctionReadWithFunctionIds" -v -count=1 -gcflags="all=-l"
func TestTeoFunctionReadWithFunctionIds(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)

	// Patch DescribeFunctions (called by DescribeTeoFunctionById)
	patches.ApplyMethodFunc(teoClient, "DescribeFunctions", func(request *teov20220901.DescribeFunctionsRequest) (*teov20220901.DescribeFunctionsResponse, error) {
		resp := teov20220901.NewDescribeFunctionsResponse()
		resp.Response = &teov20220901.DescribeFunctionsResponseParams{
			Functions: []*teov20220901.Function{
				{
					FunctionId: ptrString("func-001"),
					ZoneId:     ptrString("zone-12345678"),
					Name:       ptrString("test-func-1"),
					Remark:     ptrString("test remark 1"),
					Content:    ptrString("console.log('hello')"),
					Domain:     ptrString("func-001.example.com"),
					CreateTime: ptrString("2024-01-01T00:00:00Z"),
					UpdateTime: ptrString("2024-01-02T00:00:00Z"),
				},
				{
					FunctionId: ptrString("func-002"),
					ZoneId:     ptrString("zone-12345678"),
					Name:       ptrString("test-func-2"),
					Remark:     ptrString("test remark 2"),
					Content:    ptrString("console.log('world')"),
					Domain:     ptrString("func-002.example.com"),
					CreateTime: ptrString("2024-01-03T00:00:00Z"),
					UpdateTime: ptrString("2024-01-04T00:00:00Z"),
				},
			},
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoFunction()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":      "zone-12345678",
		"name":         "test-func-1",
		"content":      "console.log('hello')",
		"function_ids": []interface{}{"func-001", "func-002"},
	})
	d.SetId("zone-12345678#func-001")

	err := res.Read(d, meta)
	assert.NoError(t, err)

	functions := d.Get("functions").([]interface{})
	assert.Equal(t, 2, len(functions))

	func0 := functions[0].(map[string]interface{})
	assert.Equal(t, "func-001", func0["function_id"].(string))
	assert.Equal(t, "zone-12345678", func0["zone_id"].(string))
	assert.Equal(t, "test-func-1", func0["name"].(string))
	assert.Equal(t, "test remark 1", func0["remark"].(string))
	assert.Equal(t, "console.log('hello')", func0["content"].(string))
	assert.Equal(t, "func-001.example.com", func0["domain"].(string))

	func1 := functions[1].(map[string]interface{})
	assert.Equal(t, "func-002", func1["function_id"].(string))
	assert.Equal(t, "test-func-2", func1["name"].(string))
}

// go test ./tencentcloud/services/teo/ -run "TestTeoFunctionReadWithFilters" -v -count=1 -gcflags="all=-l"
func TestTeoFunctionReadWithFilters(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeFunctions", func(request *teov20220901.DescribeFunctionsRequest) (*teov20220901.DescribeFunctionsResponse, error) {
		// Verify that Filters are passed correctly
		if len(request.Filters) > 0 {
			assert.Equal(t, "name", *request.Filters[0].Name)
			assert.Equal(t, "test-func", *request.Filters[0].Values[0])
		}
		resp := teov20220901.NewDescribeFunctionsResponse()
		resp.Response = &teov20220901.DescribeFunctionsResponseParams{
			Functions: []*teov20220901.Function{
				{
					FunctionId: ptrString("func-001"),
					ZoneId:     ptrString("zone-12345678"),
					Name:       ptrString("test-func-1"),
					Remark:     ptrString("a test function"),
					Content:    ptrString("console.log('hello')"),
					Domain:     ptrString("func-001.example.com"),
					CreateTime: ptrString("2024-01-01T00:00:00Z"),
					UpdateTime: ptrString("2024-01-02T00:00:00Z"),
				},
			},
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoFunction()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-12345678",
		"name":    "test-func-1",
		"content": "console.log('hello')",
		"filters": []interface{}{
			map[string]interface{}{
				"name":   "name",
				"values": []interface{}{"test-func"},
			},
		},
	})
	d.SetId("zone-12345678#func-001")

	err := res.Read(d, meta)
	assert.NoError(t, err)

	functions := d.Get("functions").([]interface{})
	assert.Equal(t, 1, len(functions))

	func0 := functions[0].(map[string]interface{})
	assert.Equal(t, "func-001", func0["function_id"].(string))
	assert.Equal(t, "test-func-1", func0["name"].(string))
}

// go test ./tencentcloud/services/teo/ -run "TestTeoFunctionReadWithoutQueryParams" -v -count=1 -gcflags="all=-l"
func TestTeoFunctionReadWithoutQueryParams(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeFunctions", func(request *teov20220901.DescribeFunctionsRequest) (*teov20220901.DescribeFunctionsResponse, error) {
		resp := teov20220901.NewDescribeFunctionsResponse()
		resp.Response = &teov20220901.DescribeFunctionsResponseParams{
			Functions: []*teov20220901.Function{
				{
					FunctionId: ptrString("func-001"),
					ZoneId:     ptrString("zone-12345678"),
					Name:       ptrString("test-func-1"),
					Remark:     ptrString("test remark"),
					Content:    ptrString("console.log('hello')"),
					Domain:     ptrString("func-001.example.com"),
					CreateTime: ptrString("2024-01-01T00:00:00Z"),
					UpdateTime: ptrString("2024-01-02T00:00:00Z"),
				},
			},
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoFunction()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-12345678",
		"name":    "test-func-1",
		"content": "console.log('hello')",
	})
	d.SetId("zone-12345678#func-001")

	err := res.Read(d, meta)
	assert.NoError(t, err)

	// When no function_ids or filters are specified, functions should be an empty list
	functions := d.Get("functions").([]interface{})
	assert.Equal(t, 0, len(functions))
}

// go test ./tencentcloud/services/teo/ -run "TestTeoFunctionReadAPIError" -v -count=1 -gcflags="all=-l"
func TestTeoFunctionReadAPIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeFunctions", func(request *teov20220901.DescribeFunctionsRequest) (*teov20220901.DescribeFunctionsResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=ResourceNotFound, Message=Zone not found")
	})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoFunction()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-invalid",
		"name":    "test-func",
		"content": "console.log('hello')",
	})
	d.SetId("zone-invalid#func-001")

	err := res.Read(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ResourceNotFound")
}

// go test ./tencentcloud/services/teo/ -run "TestTeoFunctionSchemaNewParams" -v -count=1 -gcflags="all=-l"
func TestTeoFunctionSchemaNewParams(t *testing.T) {
	res := teo.ResourceTencentCloudTeoFunction()

	assert.NotNil(t, res)
	assert.Contains(t, res.Schema, "function_ids")
	assert.Contains(t, res.Schema, "filters")
	assert.Contains(t, res.Schema, "functions")

	functionIds := res.Schema["function_ids"]
	assert.Equal(t, schema.TypeList, functionIds.Type)
	assert.True(t, functionIds.Optional)
	assert.False(t, functionIds.Required)
	assert.False(t, functionIds.Computed)

	filters := res.Schema["filters"]
	assert.Equal(t, schema.TypeList, filters.Type)
	assert.True(t, filters.Optional)
	assert.False(t, filters.Required)
	assert.False(t, filters.Computed)

	functions := res.Schema["functions"]
	assert.Equal(t, schema.TypeList, functions.Type)
	assert.True(t, functions.Computed)
	assert.False(t, functions.Required)
	assert.False(t, functions.Optional)
}
