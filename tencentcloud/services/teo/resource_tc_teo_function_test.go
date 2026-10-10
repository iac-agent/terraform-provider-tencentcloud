package teo_test

import (
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

func newMockTeoFunctionDomainComplianceTeoClient(patches *gomonkey.Patches, restrictions []*teov20220901.ComplianceRestriction) *teov20220901.Client {
	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaTeoFunction().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeFunctions", func(request *teov20220901.DescribeFunctionsRequest) (*teov20220901.DescribeFunctionsResponse, error) {
		assert.Equal(t, "zone-2qtuhspy7cr6", *request.ZoneId)
		assert.Equal(t, "edgefunction-2qtuhspy7cr6", *request.FunctionIds[0])
		resp := teov20220901.NewDescribeFunctionsResponse()
		resp.Response = &teov20220901.DescribeFunctionsResponseParams{
			Functions: []*teov20220901.Function{
				{
					FunctionId:                   ptrStringTeoFunction("edgefunction-2qtuhspy7cr6"),
					ZoneId:                       ptrStringTeoFunction("zone-2qtuhspy7cr6"),
					Name:                         ptrStringTeoFunction("aaa"),
					Remark:                       ptrStringTeoFunction("test"),
					Content:                      ptrStringTeoFunction("addEventListener('fetch', e => { e.respondWith(new Response('Hello World!!')); });"),
					Domain:                       ptrStringTeoFunction("aaa-zone-2qtuhspy7cr6.edgefunction.app"),
					DomainComplianceRestrictions: restrictions,
					CreateTime:                   ptrStringTeoFunction("2025-01-01T00:00:00Z"),
					UpdateTime:                   ptrStringTeoFunction("2025-01-01T00:00:01Z"),
				},
			},
			RequestId: ptrStringTeoFunction("fake-request-id"),
		}
		return resp, nil
	})

	return teoClient
}

// TestTeoFunctionRead_DomainComplianceRestrictions_Present tests Read backfills
// domain_compliance_restrictions when DescribeFunctions returns non-nil restrictions.
// go test ./tencentcloud/services/teo/ -run "TestTeoFunctionRead_DomainComplianceRestrictions" -v -count=1 -gcflags="all=-l"
func TestTeoFunctionRead_DomainComplianceRestrictions_Present(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	newMockTeoFunctionDomainComplianceTeoClient(patches, []*teov20220901.ComplianceRestriction{
		{
			Reason: ptrStringTeoFunction("ICP_RECORD_REQUIRED"),
			Region: ptrStringTeoFunction("CN"),
		},
		{
			Reason: ptrStringTeoFunction("GOVERNMENT_ORDER"),
			Region: ptrStringTeoFunction("US"),
		},
	})

	meta := newMockMetaTeoFunction()
	res := teo.ResourceTencentCloudTeoFunction()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
	d.SetId("zone-2qtuhspy7cr6#edgefunction-2qtuhspy7cr6")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-2qtuhspy7cr6#edgefunction-2qtuhspy7cr6", d.Id())

	restrictions := d.Get("domain_compliance_restrictions").([]interface{})
	assert.Len(t, restrictions, 2)

	first := restrictions[0].(map[string]interface{})
	assert.Equal(t, "ICP_RECORD_REQUIRED", first["reason"].(string))
	assert.Equal(t, "CN", first["region"].(string))

	second := restrictions[1].(map[string]interface{})
	assert.Equal(t, "GOVERNMENT_ORDER", second["reason"].(string))
	assert.Equal(t, "US", second["region"].(string))
}

// TestTeoFunctionRead_DomainComplianceRestrictions_Empty tests Read keeps
// domain_compliance_restrictions as an empty list when the Function has nil restrictions.
func TestTeoFunctionRead_DomainComplianceRestrictions_Empty(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	newMockTeoFunctionDomainComplianceTeoClient(patches, nil)

	meta := newMockMetaTeoFunction()
	res := teo.ResourceTencentCloudTeoFunction()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
	d.SetId("zone-2qtuhspy7cr6#edgefunction-2qtuhspy7cr6")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-2qtuhspy7cr6#edgefunction-2qtuhspy7cr6", d.Id())

	restrictions := d.Get("domain_compliance_restrictions").([]interface{})
	assert.Empty(t, restrictions)
}
