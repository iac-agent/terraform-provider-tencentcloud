package teo_test

import (
	"context"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
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
					resource.TestCheckResourceAttr("tencentcloud_teo_function.teo_function", "name", "aaa"),
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
					resource.TestCheckResourceAttr("tencentcloud_teo_function.teo_function", "name", "aaa"),
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
    name        = "aaa"
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
    name        = "aaa"
    remark      = "test-update"
    zone_id     = "zone-2qtuhspy7cr6"
}
`

func TestParseTeoFunctionOriginalName(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    string
		zoneId   string
		expected string
	}{
		{
			name:     "standard concatenated name",
			input:    "my-func-zone-2qtuhspy7cr6-1310708577",
			zoneId:   "zone-2qtuhspy7cr6-1310708577",
			expected: "my-func",
		},
		{
			name:     "original name without hyphens",
			input:    "myfunc-zone-2qtuhspy7cr6-1310708577",
			zoneId:   "zone-2qtuhspy7cr6-1310708577",
			expected: "myfunc",
		},
		{
			name:     "original name with multiple hyphens",
			input:    "my-test-func-v2-zone-2qtuhspy7cr6-1310708577",
			zoneId:   "zone-2qtuhspy7cr6-1310708577",
			expected: "my-test-func-v2",
		},
		{
			name:     "name without -zone- substring",
			input:    "myfunc",
			zoneId:   "zone-2qtuhspy7cr6-1310708577",
			expected: "myfunc",
		},
		{
			name:     "empty string",
			input:    "",
			zoneId:   "zone-2qtuhspy7cr6-1310708577",
			expected: "",
		},
		{
			name:     "name starts with zone",
			input:    "zone-2qtuhspy7cr6-1310708577",
			zoneId:   "zone-2qtuhspy7cr6-1310708577",
			expected: "zone-2qtuhspy7cr6-1310708577",
		},
		{
			name:     "name with zone appearing before -zone-",
			input:    "zone-proxy-zone-2qtuhspy7cr6-1310708577",
			zoneId:   "zone-2qtuhspy7cr6-1310708577",
			expected: "zone-proxy",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := teo.ParseTeoFunctionOriginalName(tt.input, tt.zoneId)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// ---- Unit Tests (gomonkey mock) ----

// go test ./tencentcloud/services/teo/ -run "TestTeoFunction_" -v -count=1 -gcflags="all=-l"

// mockMetaTeoFunction implements tccommon.ProviderMeta
type mockMetaTeoFunction struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaTeoFunction) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaTeoFunction{}

func newMockMetaTeoFunction() *mockMetaTeoFunction {
	return &mockMetaTeoFunction{client: &connectivity.TencentCloudClient{}}
}

func ptrStringTeoFunction(s string) *string {
	return &s
}

// TestTeoFunction_SchemaDomainComplianceRestrictions validates the domain_compliance_restrictions schema definition
func TestTeoFunction_SchemaDomainComplianceRestrictions(t *testing.T) {
	res := teo.ResourceTencentCloudTeoFunction()

	assert.NotNil(t, res)
	assert.Contains(t, res.Schema, "domain_compliance_restrictions")

	block := res.Schema["domain_compliance_restrictions"]
	assert.Equal(t, schema.TypeList, block.Type)
	assert.True(t, block.Computed)
	assert.False(t, block.Required)
	assert.False(t, block.Optional)

	elem, ok := block.Elem.(*schema.Resource)
	assert.True(t, ok)
	assert.Contains(t, elem.Schema, "reason")
	assert.Contains(t, elem.Schema, "region")

	reasonField := elem.Schema["reason"]
	assert.Equal(t, schema.TypeString, reasonField.Type)
	assert.True(t, reasonField.Computed)

	regionField := elem.Schema["region"]
	assert.Equal(t, schema.TypeString, regionField.Type)
	assert.True(t, regionField.Computed)
}

// TestTeoFunction_ReadDomainComplianceRestrictions tests Read flattens DomainComplianceRestrictions into state
func TestTeoFunction_ReadDomainComplianceRestrictions(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaTeoFunction().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeFunctionsWithContext", func(_ context.Context, request *teov20220901.DescribeFunctionsRequest) (*teov20220901.DescribeFunctionsResponse, error) {
		assert.Equal(t, "zone-test123", *request.ZoneId)
		assert.Equal(t, "ef-test456", *request.FunctionIds[0])

		resp := teov20220901.NewDescribeFunctionsResponse()
		resp.Response = &teov20220901.DescribeFunctionsResponseParams{
			TotalCount: func() *int64 {
				i := int64(1)
				return &i
			}(),
			Functions: []*teov20220901.Function{
				{
					FunctionId: ptrStringTeoFunction("ef-test456"),
					Name:       ptrStringTeoFunction("test-function"),
					Content:    ptrStringTeoFunction("console.log(123)"),
					Remark:     ptrStringTeoFunction("test remark"),
					Domain:     ptrStringTeoFunction("test-function.example.com"),
					DomainComplianceRestrictions: []*teov20220901.ComplianceRestriction{
						{
							Reason: ptrStringTeoFunction("ICP_RECORD_REQUIRED"),
							Region: ptrStringTeoFunction("CN"),
						},
						{
							Reason: ptrStringTeoFunction("GOVERNMENT_ORDER"),
							Region: ptrStringTeoFunction("XX"),
						},
					},
					CreateTime: ptrStringTeoFunction("2024-01-01T00:00:00Z"),
					UpdateTime: ptrStringTeoFunction("2024-01-02T00:00:00Z"),
				},
			},
			RequestId: ptrStringTeoFunction("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaTeoFunction()
	res := teo.ResourceTencentCloudTeoFunction()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-test123",
		"name":    "test-function",
		"content": "console.log(123)",
	})
	d.SetId("zone-test123#ef-test456")

	err := res.Read(d, meta)
	assert.NoError(t, err)

	restrictions := d.Get("domain_compliance_restrictions").([]interface{})
	assert.Len(t, restrictions, 2)

	restriction0 := restrictions[0].(map[string]interface{})
	assert.Equal(t, "ICP_RECORD_REQUIRED", restriction0["reason"])
	assert.Equal(t, "CN", restriction0["region"])

	restriction1 := restrictions[1].(map[string]interface{})
	assert.Equal(t, "GOVERNMENT_ORDER", restriction1["reason"])
	assert.Equal(t, "XX", restriction1["region"])
}

// TestTeoFunction_ReadDomainComplianceRestrictionsNil tests Read skips the attribute when API returns nil list
func TestTeoFunction_ReadDomainComplianceRestrictionsNil(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaTeoFunction().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeFunctionsWithContext", func(_ context.Context, request *teov20220901.DescribeFunctionsRequest) (*teov20220901.DescribeFunctionsResponse, error) {
		resp := teov20220901.NewDescribeFunctionsResponse()
		resp.Response = &teov20220901.DescribeFunctionsResponseParams{
			TotalCount: func() *int64 {
				i := int64(1)
				return &i
			}(),
			Functions: []*teov20220901.Function{
				{
					FunctionId:                   ptrStringTeoFunction("ef-test456"),
					Name:                         ptrStringTeoFunction("test-function"),
					Content:                      ptrStringTeoFunction("console.log(123)"),
					Remark:                       ptrStringTeoFunction("test remark"),
					Domain:                       ptrStringTeoFunction("test-function.example.com"),
					DomainComplianceRestrictions: nil,
					CreateTime:                   ptrStringTeoFunction("2024-01-01T00:00:00Z"),
					UpdateTime:                   ptrStringTeoFunction("2024-01-02T00:00:00Z"),
				},
			},
			RequestId: ptrStringTeoFunction("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaTeoFunction()
	res := teo.ResourceTencentCloudTeoFunction()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-test123",
		"name":    "test-function",
		"content": "console.log(123)",
	})
	d.SetId("zone-test123#ef-test456")

	err := res.Read(d, meta)
	assert.NoError(t, err)

	restrictions := d.Get("domain_compliance_restrictions").([]interface{})
	assert.Len(t, restrictions, 0)
}

// TestTeoFunction_ReadDomainComplianceRestrictionsEmptyList tests Read skips the attribute when API returns an empty list
func TestTeoFunction_ReadDomainComplianceRestrictionsEmptyList(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaTeoFunction().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeFunctionsWithContext", func(_ context.Context, request *teov20220901.DescribeFunctionsRequest) (*teov20220901.DescribeFunctionsResponse, error) {
		resp := teov20220901.NewDescribeFunctionsResponse()
		resp.Response = &teov20220901.DescribeFunctionsResponseParams{
			TotalCount: func() *int64 {
				i := int64(1)
				return &i
			}(),
			Functions: []*teov20220901.Function{
				{
					FunctionId:                   ptrStringTeoFunction("ef-test456"),
					Name:                         ptrStringTeoFunction("test-function"),
					Content:                      ptrStringTeoFunction("console.log(123)"),
					Domain:                       ptrStringTeoFunction("test-function.example.com"),
					DomainComplianceRestrictions: []*teov20220901.ComplianceRestriction{},
					CreateTime:                   ptrStringTeoFunction("2024-01-01T00:00:00Z"),
					UpdateTime:                   ptrStringTeoFunction("2024-01-02T00:00:00Z"),
				},
			},
			RequestId: ptrStringTeoFunction("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaTeoFunction()
	res := teo.ResourceTencentCloudTeoFunction()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-test123",
		"name":    "test-function",
		"content": "console.log(123)",
	})
	d.SetId("zone-test123#ef-test456")

	err := res.Read(d, meta)
	assert.NoError(t, err)

	restrictions := d.Get("domain_compliance_restrictions").([]interface{})
	assert.Len(t, restrictions, 0)
}

// TestTeoFunction_ReadDomainComplianceRestrictionsNilSubFields tests Read skips nil sub-fields inside an entry
func TestTeoFunction_ReadDomainComplianceRestrictionsNilSubFields(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaTeoFunction().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeFunctionsWithContext", func(_ context.Context, request *teov20220901.DescribeFunctionsRequest) (*teov20220901.DescribeFunctionsResponse, error) {
		resp := teov20220901.NewDescribeFunctionsResponse()
		resp.Response = &teov20220901.DescribeFunctionsResponseParams{
			TotalCount: func() *int64 {
				i := int64(1)
				return &i
			}(),
			Functions: []*teov20220901.Function{
				{
					FunctionId: ptrStringTeoFunction("ef-test456"),
					Name:       ptrStringTeoFunction("test-function"),
					Content:    ptrStringTeoFunction("console.log(123)"),
					Domain:     ptrStringTeoFunction("test-function.example.com"),
					DomainComplianceRestrictions: []*teov20220901.ComplianceRestriction{
						{
							// Reason is nil, only Region set
							Reason: nil,
							Region: ptrStringTeoFunction("CN"),
						},
						{
							// Region is nil, only Reason set
							Reason: ptrStringTeoFunction("GOVERNMENT_ORDER"),
							Region: nil,
						},
					},
					CreateTime: ptrStringTeoFunction("2024-01-01T00:00:00Z"),
					UpdateTime: ptrStringTeoFunction("2024-01-02T00:00:00Z"),
				},
			},
			RequestId: ptrStringTeoFunction("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaTeoFunction()
	res := teo.ResourceTencentCloudTeoFunction()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-test123",
		"name":    "test-function",
		"content": "console.log(123)",
	})
	d.SetId("zone-test123#ef-test456")

	err := res.Read(d, meta)
	assert.NoError(t, err)

	restrictions := d.Get("domain_compliance_restrictions").([]interface{})
	assert.Len(t, restrictions, 2)

	// nil sub-fields are not written, so reading them back yields the zero value (empty string)
	restriction0 := restrictions[0].(map[string]interface{})
	assert.Equal(t, "", restriction0["reason"])
	assert.Equal(t, "CN", restriction0["region"])

	restriction1 := restrictions[1].(map[string]interface{})
	assert.Equal(t, "GOVERNMENT_ORDER", restriction1["reason"])
	assert.Equal(t, "", restriction1["region"])
}
