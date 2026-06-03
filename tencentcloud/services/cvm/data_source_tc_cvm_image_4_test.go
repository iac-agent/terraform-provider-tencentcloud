package cvm_test

import (
	"fmt"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	cvmservice "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/cvm"
)

// go test ./tencentcloud/services/cvm/ -run "TestCvmImage4" -v -count=1 -gcflags="all=-l"

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

func ptrInt64(i int64) *int64 {
	return &i
}

func ptrBool(b bool) *bool {
	return &b
}

// TestCvmImage4_Read_ByImageIds tests Read with image_ids parameter
func TestCvmImage4_Read_ByImageIds(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	cvmClient := &cvm.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseCvmClient", cvmClient)

	patches.ApplyMethodFunc(cvmClient, "DescribeImages", func(request *cvm.DescribeImagesRequest) (*cvm.DescribeImagesResponse, error) {
		resp := cvm.NewDescribeImagesResponse()
		resp.Response = &cvm.DescribeImagesResponseParams{
			TotalCount: ptrInt64(1),
			ImageSet: []*cvm.Image{
				{
					ImageId:            ptrString("img-8toqc6s3"),
					OsName:             ptrString("TencentOS Server 3.1"),
					ImageType:          ptrString("PUBLIC_IMAGE"),
					CreatedTime:        ptrString("2024-01-01T00:00:00Z"),
					ImageName:          ptrString("TencentOS Server 3.1"),
					ImageDescription:   ptrString("TencentOS Server 3.1 64bit"),
					ImageSize:          ptrInt64(50),
					Architecture:       ptrString("x86_64"),
					ImageState:         ptrString("NORMAL"),
					Platform:           ptrString("TencentOS"),
					ImageCreator:       ptrString(""),
					ImageSource:        ptrString("OFFICIAL"),
					SyncPercent:        ptrInt64(100),
					IsSupportCloudinit: ptrBool(true),
					SnapshotSet: []*cvm.Snapshot{
						{
							SnapshotId: ptrString("snap-xxxxxxxx"),
							DiskUsage:  ptrString("SYSTEM_DISK"),
							DiskSize:   ptrInt64(50),
						},
					},
					Tags: []*cvm.Tag{
						{
							Key:   ptrString("test-key"),
							Value: ptrString("test-value"),
						},
					},
					LicenseType:     ptrString("TencentCloud"),
					ImageFamily:     ptrString("TencentOS"),
					ImageDeprecated: ptrBool(false),
					CdcCacheStatus:  ptrString(""),
				},
			},
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := cvmservice.DataSourceTencentCloudCvmImage4()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"image_ids": []interface{}{"img-8toqc6s3"},
	})

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.NotEmpty(t, d.Id())

	imageSet := d.Get("image_set").([]interface{})
	assert.Equal(t, 1, len(imageSet))
	image := imageSet[0].(map[string]interface{})
	assert.Equal(t, "img-8toqc6s3", image["image_id"])
	assert.Equal(t, "TencentOS Server 3.1", image["os_name"])
	assert.Equal(t, "PUBLIC_IMAGE", image["image_type"])
	assert.Equal(t, "NORMAL", image["image_state"])
	assert.Equal(t, 50, image["image_size"])
	assert.Equal(t, "x86_64", image["architecture"])
	assert.Equal(t, true, image["is_support_cloudinit"])

	snapshotSet := image["snapshot_set"].([]interface{})
	assert.Equal(t, 1, len(snapshotSet))
	snapshot := snapshotSet[0].(map[string]interface{})
	assert.Equal(t, "snap-xxxxxxxx", snapshot["snapshot_id"])
	assert.Equal(t, "SYSTEM_DISK", snapshot["disk_usage"])
	assert.Equal(t, 50, snapshot["disk_size"])

	tags := image["tags"].([]interface{})
	assert.Equal(t, 1, len(tags))
	tag := tags[0].(map[string]interface{})
	assert.Equal(t, "test-key", tag["key"])
	assert.Equal(t, "test-value", tag["value"])

	assert.Equal(t, "TencentCloud", image["license_type"])
	assert.Equal(t, "TencentOS", image["image_family"])
	assert.Equal(t, false, image["image_deprecated"])
}

// TestCvmImage4_Read_ByFilters tests Read with filters parameter
func TestCvmImage4_Read_ByFilters(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	cvmClient := &cvm.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseCvmClient", cvmClient)

	patches.ApplyMethodFunc(cvmClient, "DescribeImages", func(request *cvm.DescribeImagesRequest) (*cvm.DescribeImagesResponse, error) {
		resp := cvm.NewDescribeImagesResponse()
		resp.Response = &cvm.DescribeImagesResponseParams{
			TotalCount: ptrInt64(1),
			ImageSet: []*cvm.Image{
				{
					ImageId:   ptrString("img-xxxxxxxx"),
					OsName:    ptrString("Ubuntu 22.04"),
					ImageType: ptrString("PUBLIC_IMAGE"),
					ImageName: ptrString("Ubuntu Server 22.04 LTS"),
					ImageSize: ptrInt64(40),
				},
			},
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := cvmservice.DataSourceTencentCloudCvmImage4()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"filters": []interface{}{
			map[string]interface{}{
				"name":   "image-type",
				"values": []interface{}{"PUBLIC_IMAGE"},
			},
		},
	})

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.NotEmpty(t, d.Id())

	imageSet := d.Get("image_set").([]interface{})
	assert.Equal(t, 1, len(imageSet))
	image := imageSet[0].(map[string]interface{})
	assert.Equal(t, "img-xxxxxxxx", image["image_id"])
	assert.Equal(t, "Ubuntu 22.04", image["os_name"])
}

// TestCvmImage4_Read_WithInstanceType tests Read with instance_type parameter
func TestCvmImage4_Read_WithInstanceType(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	cvmClient := &cvm.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseCvmClient", cvmClient)

	patches.ApplyMethodFunc(cvmClient, "DescribeImages", func(request *cvm.DescribeImagesRequest) (*cvm.DescribeImagesResponse, error) {
		resp := cvm.NewDescribeImagesResponse()
		resp.Response = &cvm.DescribeImagesResponseParams{
			TotalCount: ptrInt64(1),
			ImageSet: []*cvm.Image{
				{
					ImageId:   ptrString("img-yyyyyyyy"),
					OsName:    ptrString("CentOS 8.0"),
					ImageType: ptrString("PUBLIC_IMAGE"),
					ImageName: ptrString("CentOS 8.0 64bit"),
					ImageSize: ptrInt64(50),
				},
			},
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := cvmservice.DataSourceTencentCloudCvmImage4()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"instance_type": "SA5.MEDIUM2",
	})

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.NotEmpty(t, d.Id())

	imageSet := d.Get("image_set").([]interface{})
	assert.Equal(t, 1, len(imageSet))
	image := imageSet[0].(map[string]interface{})
	assert.Equal(t, "img-yyyyyyyy", image["image_id"])
}

// TestCvmImage4_Read_NilFields tests Read handles nil fields in response
func TestCvmImage4_Read_NilFields(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	cvmClient := &cvm.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseCvmClient", cvmClient)

	patches.ApplyMethodFunc(cvmClient, "DescribeImages", func(request *cvm.DescribeImagesRequest) (*cvm.DescribeImagesResponse, error) {
		resp := cvm.NewDescribeImagesResponse()
		resp.Response = &cvm.DescribeImagesResponseParams{
			TotalCount: ptrInt64(1),
			ImageSet: []*cvm.Image{
				{
					ImageId:   ptrString("img-zzzzzzzz"),
					ImageName: ptrString("Minimal Image"),
				},
			},
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := cvmservice.DataSourceTencentCloudCvmImage4()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.NotEmpty(t, d.Id())

	imageSet := d.Get("image_set").([]interface{})
	assert.Equal(t, 1, len(imageSet))
	image := imageSet[0].(map[string]interface{})
	assert.Equal(t, "img-zzzzzzzz", image["image_id"])
	assert.Equal(t, "Minimal Image", image["image_name"])
}

// TestCvmImage4_Read_APIError tests Read handles API error
func TestCvmImage4_Read_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	cvmClient := &cvm.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseCvmClient", cvmClient)

	patches.ApplyMethodFunc(cvmClient, "DescribeImages", func(request *cvm.DescribeImagesRequest) (*cvm.DescribeImagesResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=InvalidParameter, Message=Invalid image id")
	})

	meta := newMockMeta()
	res := cvmservice.DataSourceTencentCloudCvmImage4()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"image_ids": []interface{}{"img-invalid"},
	})

	err := res.Read(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "InvalidParameter")
}

// TestCvmImage4_Schema validates schema definition
func TestCvmImage4_Schema(t *testing.T) {
	res := cvmservice.DataSourceTencentCloudCvmImage4()

	assert.NotNil(t, res)
	assert.NotNil(t, res.Read)

	assert.Contains(t, res.Schema, "image_ids")
	imageIds := res.Schema["image_ids"]
	assert.Equal(t, schema.TypeList, imageIds.Type)
	assert.True(t, imageIds.Optional)

	assert.Contains(t, res.Schema, "filters")
	filters := res.Schema["filters"]
	assert.Equal(t, schema.TypeList, filters.Type)
	assert.True(t, filters.Optional)

	assert.Contains(t, res.Schema, "instance_type")
	instanceType := res.Schema["instance_type"]
	assert.Equal(t, schema.TypeString, instanceType.Type)
	assert.True(t, instanceType.Optional)

	assert.Contains(t, res.Schema, "image_set")
	imageSet := res.Schema["image_set"]
	assert.Equal(t, schema.TypeList, imageSet.Type)
	assert.True(t, imageSet.Computed)

	assert.Contains(t, res.Schema, "result_output_file")
	resultOutputFile := res.Schema["result_output_file"]
	assert.Equal(t, schema.TypeString, resultOutputFile.Type)
	assert.True(t, resultOutputFile.Optional)
}
