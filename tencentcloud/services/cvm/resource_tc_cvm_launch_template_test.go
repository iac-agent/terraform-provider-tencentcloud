package cvm_test

import (
	"fmt"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	cvmv20170312 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/cvm"
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

func ptrString(s string) *string {
	return &s
}

func ptrBool(b bool) *bool {
	return &b
}

func ptrUint64(u uint64) *uint64 {
	return &u
}

func ptrInt64(i int64) *int64 {
	return &i
}

// mockReadAPIs mocks the DescribeLaunchTemplates and DescribeLaunchTemplateVersions APIs
// so that the Read function called after Create can succeed.
func mockReadAPIs(patches *gomonkey.Patches, cvmClient *cvmv20170312.Client) {
	patches.ApplyMethodFunc(cvmClient, "DescribeLaunchTemplates", func(request *cvmv20170312.DescribeLaunchTemplatesRequest) (*cvmv20170312.DescribeLaunchTemplatesResponse, error) {
		resp := cvmv20170312.NewDescribeLaunchTemplatesResponse()
		resp.Response = &cvmv20170312.DescribeLaunchTemplatesResponseParams{
			TotalCount: ptrInt64(1),
			LaunchTemplateSet: []*cvmv20170312.LaunchTemplateInfo{
				{
					LaunchTemplateId:   ptrString("lt-abcdefghij"),
					LaunchTemplateName: ptrString("test-launch-template"),
				},
			},
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(cvmClient, "DescribeLaunchTemplateVersions", func(request *cvmv20170312.DescribeLaunchTemplateVersionsRequest) (*cvmv20170312.DescribeLaunchTemplateVersionsResponse, error) {
		resp := cvmv20170312.NewDescribeLaunchTemplateVersionsResponse()
		resp.Response = &cvmv20170312.DescribeLaunchTemplateVersionsResponseParams{
			TotalCount: ptrUint64(1),
			LaunchTemplateVersionSet: []*cvmv20170312.LaunchTemplateVersionInfo{
				{
					LaunchTemplateId:      ptrString("lt-abcdefghij"),
					LaunchTemplateVersion: ptrUint64(1),
					LaunchTemplateVersionData: &cvmv20170312.LaunchTemplateVersionData{
						Placement: &cvmv20170312.Placement{
							Zone:      ptrString("ap-guangzhou-6"),
							ProjectId: ptrInt64(0),
						},
						ImageId:       ptrString("img-9qrfy1xt"),
						InstanceCount: ptrUint64(1),
					},
				},
			},
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})
}

// go test ./tencentcloud/services/cvm/ -run "TestCvmLaunchTemplateEnableJumboFrame" -v -count=1 -gcflags="all=-l"

// TestCvmLaunchTemplateEnableJumboFrame_CreateWithTrue tests Create sets EnableJumboFrame=true in request
func TestCvmLaunchTemplateEnableJumboFrame_CreateWithTrue(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	cvmClient := &cvmv20170312.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseCvmClient", cvmClient)

	patches.ApplyMethodFunc(cvmClient, "CreateLaunchTemplate", func(request *cvmv20170312.CreateLaunchTemplateRequest) (*cvmv20170312.CreateLaunchTemplateResponse, error) {
		assert.NotNil(t, request.EnableJumboFrame)
		assert.Equal(t, true, *request.EnableJumboFrame)
		resp := cvmv20170312.NewCreateLaunchTemplateResponse()
		resp.Response = &cvmv20170312.CreateLaunchTemplateResponseParams{
			LaunchTemplateId: ptrString("lt-abcdefghij"),
			RequestId:        ptrString("fake-request-id"),
		}
		return resp, nil
	})

	mockReadAPIs(patches, cvmClient)

	meta := newMockMeta()
	res := cvm.ResourceTencentCloudCvmLaunchTemplate()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"launch_template_name": "test-launch-template",
		"placement": []interface{}{
			map[string]interface{}{
				"zone":       "ap-guangzhou-6",
				"project_id": 0,
				"host_ids":   []interface{}{},
				"host_ips":   []interface{}{},
			},
		},
		"image_id":           "img-9qrfy1xt",
		"enable_jumbo_frame": true,
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "lt-abcdefghij", d.Id())
}

// TestCvmLaunchTemplateEnableJumboFrame_CreateWithFalse tests Create sets EnableJumboFrame=false in request
func TestCvmLaunchTemplateEnableJumboFrame_CreateWithFalse(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	cvmClient := &cvmv20170312.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseCvmClient", cvmClient)

	patches.ApplyMethodFunc(cvmClient, "CreateLaunchTemplate", func(request *cvmv20170312.CreateLaunchTemplateRequest) (*cvmv20170312.CreateLaunchTemplateResponse, error) {
		assert.NotNil(t, request.EnableJumboFrame)
		assert.Equal(t, false, *request.EnableJumboFrame)
		resp := cvmv20170312.NewCreateLaunchTemplateResponse()
		resp.Response = &cvmv20170312.CreateLaunchTemplateResponseParams{
			LaunchTemplateId: ptrString("lt-abcdefghij"),
			RequestId:        ptrString("fake-request-id"),
		}
		return resp, nil
	})

	mockReadAPIs(patches, cvmClient)

	meta := newMockMeta()
	res := cvm.ResourceTencentCloudCvmLaunchTemplate()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"launch_template_name": "test-launch-template",
		"placement": []interface{}{
			map[string]interface{}{
				"zone":       "ap-guangzhou-6",
				"project_id": 0,
				"host_ids":   []interface{}{},
				"host_ips":   []interface{}{},
			},
		},
		"image_id":           "img-9qrfy1xt",
		"enable_jumbo_frame": false,
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "lt-abcdefghij", d.Id())
}

// TestCvmLaunchTemplateEnableJumboFrame_CreateWithout tests Create when enable_jumbo_frame is not specified.
// Since d.GetOk returns the zero value (false) for TypeBool when unset, and v != nil for interface{}{false},
// the code sets EnableJumboFrame=false, which is consistent with disable_api_termination and dry_run.
func TestCvmLaunchTemplateEnableJumboFrame_CreateWithout(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	cvmClient := &cvmv20170312.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseCvmClient", cvmClient)

	patches.ApplyMethodFunc(cvmClient, "CreateLaunchTemplate", func(request *cvmv20170312.CreateLaunchTemplateRequest) (*cvmv20170312.CreateLaunchTemplateResponse, error) {
		// When enable_jumbo_frame is not set, d.GetOk returns (false, false) for TypeBool.
		// The code uses `if v, _ := d.GetOk("enable_jumbo_frame"); v != nil` which is true for interface{}{false},
		// so EnableJumboFrame is set to false (consistent with disable_api_termination and dry_run behavior).
		assert.NotNil(t, request.EnableJumboFrame)
		assert.Equal(t, false, *request.EnableJumboFrame)
		resp := cvmv20170312.NewCreateLaunchTemplateResponse()
		resp.Response = &cvmv20170312.CreateLaunchTemplateResponseParams{
			LaunchTemplateId: ptrString("lt-abcdefghij"),
			RequestId:        ptrString("fake-request-id"),
		}
		return resp, nil
	})

	mockReadAPIs(patches, cvmClient)

	meta := newMockMeta()
	res := cvm.ResourceTencentCloudCvmLaunchTemplate()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"launch_template_name": "test-launch-template",
		"placement": []interface{}{
			map[string]interface{}{
				"zone":       "ap-guangzhou-6",
				"project_id": 0,
				"host_ids":   []interface{}{},
				"host_ips":   []interface{}{},
			},
		},
		"image_id": "img-9qrfy1xt",
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "lt-abcdefghij", d.Id())
}

// TestCvmLaunchTemplateEnableJumboFrame_Create_APIError tests Create handles API error
func TestCvmLaunchTemplateEnableJumboFrame_Create_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	cvmClient := &cvmv20170312.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseCvmClient", cvmClient)

	patches.ApplyMethodFunc(cvmClient, "CreateLaunchTemplate", func(request *cvmv20170312.CreateLaunchTemplateRequest) (*cvmv20170312.CreateLaunchTemplateResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=InvalidParameter, Message=Invalid parameter")
	})

	meta := newMockMeta()
	res := cvm.ResourceTencentCloudCvmLaunchTemplate()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"launch_template_name": "test-launch-template",
		"placement": []interface{}{
			map[string]interface{}{
				"zone":       "ap-guangzhou-6",
				"project_id": 0,
				"host_ids":   []interface{}{},
				"host_ips":   []interface{}{},
			},
		},
		"image_id":           "img-9qrfy1xt",
		"enable_jumbo_frame": true,
	})

	err := res.Create(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "InvalidParameter")
}

// TestCvmLaunchTemplateEnableJumboFrame_Schema validates schema definition
func TestCvmLaunchTemplateEnableJumboFrame_Schema(t *testing.T) {
	res := cvm.ResourceTencentCloudCvmLaunchTemplate()

	assert.NotNil(t, res)
	assert.Contains(t, res.Schema, "enable_jumbo_frame")

	enableJumboFrame := res.Schema["enable_jumbo_frame"]
	assert.Equal(t, schema.TypeBool, enableJumboFrame.Type)
	assert.True(t, enableJumboFrame.Optional)
	assert.True(t, enableJumboFrame.ForceNew)
	assert.False(t, enableJumboFrame.Required)
	assert.False(t, enableJumboFrame.Computed)
}

func TestAccTencentCloudCvmLaunchTemplateResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCvmLaunchTemplate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_cvm_launch_template.launch_template", "id"),
					resource.TestCheckResourceAttr("tencentcloud_cvm_launch_template.launch_template", "image_id", "img-9qrfy1xt"),
				),
			},
		},
	})
}

const testAccCvmLaunchTemplate = `
resource "tencentcloud_cvm_launch_template" "launch_template" {
	launch_template_name = "test_launch_template"
	placement {
	  zone = "ap-guangzhou-6"
	  project_id = 0
	  host_ids = []
	  host_ips = []
	}
	image_id = "img-9qrfy1xt"
  }
`
