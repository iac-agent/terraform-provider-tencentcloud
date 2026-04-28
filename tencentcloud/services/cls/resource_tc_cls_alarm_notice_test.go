package cls_test

import (
	"fmt"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	tccls "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cls/v20201016"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
	localcls "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/cls"
)

// mockMeta implements tccommon.ProviderMeta
type alarmNoticeMockMeta struct {
	client *connectivity.TencentCloudClient
}

func (m *alarmNoticeMockMeta) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &alarmNoticeMockMeta{}

func newAlarmNoticeMockMeta() *alarmNoticeMockMeta {
	return &alarmNoticeMockMeta{client: &connectivity.TencentCloudClient{}}
}

func ptrString(s string) *string { return &s }
func ptrUint64(u uint64) *uint64 { return &u }
func ptrInt64(i int64) *int64    { return &i }
func ptrBool(b bool) *bool       { return &b }

// go test ./tencentcloud/services/cls/ -run "TestClsAlarmNotice" -v -count=1 -gcflags="all=-l"

// TestClsAlarmNotice_Create_WithTags tests Create passes Tags through CLS API
func TestClsAlarmNotice_Create_WithTags(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	clsClient := &tccls.Client{}
	patches.ApplyMethodReturn(newAlarmNoticeMockMeta().client, "UseClsClient", clsClient)

	patches.ApplyMethodFunc(clsClient, "CreateAlarmNotice", func(request *tccls.CreateAlarmNoticeRequest) (*tccls.CreateAlarmNoticeResponse, error) {
		assert.NotNil(t, request.Tags)
		assert.Equal(t, 2, len(request.Tags), "Expected 2 tags in Create request")
		tagMap := make(map[string]string)
		for _, tag := range request.Tags {
			if tag.Key != nil && tag.Value != nil {
				tagMap[*tag.Key] = *tag.Value
			}
		}
		assert.Equal(t, "terraform", tagMap["createdBy"])
		assert.Equal(t, "cls", tagMap["env"])

		resp := tccls.NewCreateAlarmNoticeResponse()
		resp.Response = &tccls.CreateAlarmNoticeResponseParams{
			AlarmNoticeId: ptrString("notice-abcdefghij"),
			RequestId:     ptrString("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(clsClient, "DescribeAlarmNotices", func(request *tccls.DescribeAlarmNoticesRequest) (*tccls.DescribeAlarmNoticesResponse, error) {
		resp := tccls.NewDescribeAlarmNoticesResponse()
		resp.Response = &tccls.DescribeAlarmNoticesResponseParams{
			AlarmNotices: []*tccls.AlarmNotice{
				{
					AlarmNoticeId: ptrString("notice-abcdefghij"),
					Name:          ptrString("test-alarm-notice"),
					Type:          ptrString("All"),
					Tags: []*tccls.Tag{
						{Key: ptrString("createdBy"), Value: ptrString("terraform")},
						{Key: ptrString("env"), Value: ptrString("cls")},
					},
				},
			},
			TotalCount: ptrInt64(1),
			RequestId:  ptrString("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyFuncReturn(ratelimit.Check)

	meta := newAlarmNoticeMockMeta()
	res := localcls.ResourceTencentCloudClsAlarmNotice()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"name": "test-alarm-notice",
		"type": "All",
		"tags": map[string]interface{}{
			"createdBy": "terraform",
			"env":       "cls",
		},
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "notice-abcdefghij", d.Id())

	tags := d.Get("tags").(map[string]interface{})
	assert.Equal(t, "terraform", tags["createdBy"])
	assert.Equal(t, "cls", tags["env"])
}

// TestClsAlarmNotice_Create_WithoutTags tests Create without tags
func TestClsAlarmNotice_Create_WithoutTags(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	clsClient := &tccls.Client{}
	patches.ApplyMethodReturn(newAlarmNoticeMockMeta().client, "UseClsClient", clsClient)

	patches.ApplyMethodFunc(clsClient, "CreateAlarmNotice", func(request *tccls.CreateAlarmNoticeRequest) (*tccls.CreateAlarmNoticeResponse, error) {
		assert.Equal(t, 0, len(request.Tags), "Expected no tags in Create request")
		resp := tccls.NewCreateAlarmNoticeResponse()
		resp.Response = &tccls.CreateAlarmNoticeResponseParams{
			AlarmNoticeId: ptrString("notice-no-tags"),
			RequestId:     ptrString("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(clsClient, "DescribeAlarmNotices", func(request *tccls.DescribeAlarmNoticesRequest) (*tccls.DescribeAlarmNoticesResponse, error) {
		resp := tccls.NewDescribeAlarmNoticesResponse()
		resp.Response = &tccls.DescribeAlarmNoticesResponseParams{
			AlarmNotices: []*tccls.AlarmNotice{
				{
					AlarmNoticeId: ptrString("notice-no-tags"),
					Name:          ptrString("test-alarm-notice"),
					Type:          ptrString("All"),
				},
			},
			TotalCount: ptrInt64(1),
			RequestId:  ptrString("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyFuncReturn(ratelimit.Check)

	meta := newAlarmNoticeMockMeta()
	res := localcls.ResourceTencentCloudClsAlarmNotice()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"name": "test-alarm-notice",
		"type": "All",
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "notice-no-tags", d.Id())
}

// TestClsAlarmNotice_Create_APIError tests Create handles API error
func TestClsAlarmNotice_Create_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	clsClient := &tccls.Client{}
	patches.ApplyMethodReturn(newAlarmNoticeMockMeta().client, "UseClsClient", clsClient)

	patches.ApplyMethodFunc(clsClient, "CreateAlarmNotice", func(request *tccls.CreateAlarmNoticeRequest) (*tccls.CreateAlarmNoticeResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=InvalidParameter, Message=Invalid parameter")
	})

	meta := newAlarmNoticeMockMeta()
	res := localcls.ResourceTencentCloudClsAlarmNotice()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"name": "test-alarm-notice",
	})

	err := res.Create(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "InvalidParameter")
}

// TestClsAlarmNotice_Read_WithTags tests Read retrieves tags from CLS API response
func TestClsAlarmNotice_Read_WithTags(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	clsClient := &tccls.Client{}
	patches.ApplyMethodReturn(newAlarmNoticeMockMeta().client, "UseClsClient", clsClient)

	patches.ApplyMethodFunc(clsClient, "DescribeAlarmNotices", func(request *tccls.DescribeAlarmNoticesRequest) (*tccls.DescribeAlarmNoticesResponse, error) {
		resp := tccls.NewDescribeAlarmNoticesResponse()
		resp.Response = &tccls.DescribeAlarmNoticesResponseParams{
			AlarmNotices: []*tccls.AlarmNotice{
				{
					AlarmNoticeId:     ptrString("notice-abcdefghij"),
					Name:              ptrString("test-alarm-notice"),
					Type:              ptrString("All"),
					AlarmShieldStatus: ptrUint64(2),
					Tags: []*tccls.Tag{
						{Key: ptrString("createdBy"), Value: ptrString("terraform")},
						{Key: ptrString("env"), Value: ptrString("cls")},
					},
				},
			},
			TotalCount: ptrInt64(1),
			RequestId:  ptrString("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyFuncReturn(ratelimit.Check)

	meta := newAlarmNoticeMockMeta()
	res := localcls.ResourceTencentCloudClsAlarmNotice()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"name": "test-alarm-notice",
	})
	d.SetId("notice-abcdefghij")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "test-alarm-notice", d.Get("name"))
	assert.Equal(t, "All", d.Get("type"))

	tags := d.Get("tags").(map[string]interface{})
	assert.Equal(t, "terraform", tags["createdBy"])
	assert.Equal(t, "cls", tags["env"])
}

// TestClsAlarmNotice_Read_WithoutTags tests Read with no tags in response
func TestClsAlarmNotice_Read_WithoutTags(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	clsClient := &tccls.Client{}
	patches.ApplyMethodReturn(newAlarmNoticeMockMeta().client, "UseClsClient", clsClient)

	patches.ApplyMethodFunc(clsClient, "DescribeAlarmNotices", func(request *tccls.DescribeAlarmNoticesRequest) (*tccls.DescribeAlarmNoticesResponse, error) {
		resp := tccls.NewDescribeAlarmNoticesResponse()
		resp.Response = &tccls.DescribeAlarmNoticesResponseParams{
			AlarmNotices: []*tccls.AlarmNotice{
				{
					AlarmNoticeId:     ptrString("notice-no-tags"),
					Name:              ptrString("test-alarm-notice"),
					Type:              ptrString("All"),
					AlarmShieldStatus: ptrUint64(2),
				},
			},
			TotalCount: ptrInt64(1),
			RequestId:  ptrString("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyFuncReturn(ratelimit.Check)

	meta := newAlarmNoticeMockMeta()
	res := localcls.ResourceTencentCloudClsAlarmNotice()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"name": "test-alarm-notice",
	})
	d.SetId("notice-no-tags")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "test-alarm-notice", d.Get("name"))
	tags := d.Get("tags").(map[string]interface{})
	assert.Equal(t, 0, len(tags))
}

// TestClsAlarmNotice_Read_NotFound tests Read handles not found
func TestClsAlarmNotice_Read_NotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	clsClient := &tccls.Client{}
	patches.ApplyMethodReturn(newAlarmNoticeMockMeta().client, "UseClsClient", clsClient)

	patches.ApplyMethodFunc(clsClient, "DescribeAlarmNotices", func(request *tccls.DescribeAlarmNoticesRequest) (*tccls.DescribeAlarmNoticesResponse, error) {
		resp := tccls.NewDescribeAlarmNoticesResponse()
		resp.Response = &tccls.DescribeAlarmNoticesResponseParams{
			AlarmNotices: []*tccls.AlarmNotice{},
			TotalCount:   ptrInt64(0),
			RequestId:    ptrString("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyFuncReturn(ratelimit.Check)

	meta := newAlarmNoticeMockMeta()
	res := localcls.ResourceTencentCloudClsAlarmNotice()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"name": "test-alarm-notice",
	})
	d.SetId("notice-not-exist")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id(), "Resource ID should be cleared when not found")
}

// TestClsAlarmNotice_Update_WithTags tests Update passes Tags through CLS API
func TestClsAlarmNotice_Update_WithTags(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	clsClient := &tccls.Client{}
	patches.ApplyMethodReturn(newAlarmNoticeMockMeta().client, "UseClsClient", clsClient)

	patches.ApplyMethodFunc(clsClient, "ModifyAlarmNotice", func(request *tccls.ModifyAlarmNoticeRequest) (*tccls.ModifyAlarmNoticeResponse, error) {
		assert.NotNil(t, request.Tags, "Tags should be set in Modify request")
		tagMap := make(map[string]string)
		for _, tag := range request.Tags {
			if tag.Key != nil && tag.Value != nil {
				tagMap[*tag.Key] = *tag.Value
			}
		}
		assert.Equal(t, "terraform-updated", tagMap["createdBy"])
		assert.Equal(t, "production", tagMap["env"])

		resp := tccls.NewModifyAlarmNoticeResponse()
		resp.Response = &tccls.ModifyAlarmNoticeResponseParams{
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(clsClient, "DescribeAlarmNotices", func(request *tccls.DescribeAlarmNoticesRequest) (*tccls.DescribeAlarmNoticesResponse, error) {
		resp := tccls.NewDescribeAlarmNoticesResponse()
		resp.Response = &tccls.DescribeAlarmNoticesResponseParams{
			AlarmNotices: []*tccls.AlarmNotice{
				{
					AlarmNoticeId:     ptrString("notice-abcdefghij"),
					Name:              ptrString("test-alarm-notice-updated"),
					Type:              ptrString("All"),
					AlarmShieldStatus: ptrUint64(2),
					Tags: []*tccls.Tag{
						{Key: ptrString("createdBy"), Value: ptrString("terraform-updated")},
						{Key: ptrString("env"), Value: ptrString("production")},
					},
				},
			},
			TotalCount: ptrInt64(1),
			RequestId:  ptrString("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyFuncReturn(ratelimit.Check)

	meta := newAlarmNoticeMockMeta()
	res := localcls.ResourceTencentCloudClsAlarmNotice()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"name": "test-alarm-notice-updated",
		"type": "All",
		"tags": map[string]interface{}{
			"createdBy": "terraform-updated",
			"env":       "production",
		},
	})
	d.SetId("notice-abcdefghij")

	err := res.Update(d, meta)
	assert.NoError(t, err)
}

// TestClsAlarmNotice_Delete_Success tests Delete removes alarm notice
func TestClsAlarmNotice_Delete_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	clsClient := &tccls.Client{}
	patches.ApplyMethodReturn(newAlarmNoticeMockMeta().client, "UseClsClient", clsClient)

	patches.ApplyMethodFunc(clsClient, "DeleteAlarmNotice", func(request *tccls.DeleteAlarmNoticeRequest) (*tccls.DeleteAlarmNoticeResponse, error) {
		assert.NotNil(t, request.AlarmNoticeId)
		assert.Equal(t, "notice-abcdefghij", *request.AlarmNoticeId)

		resp := tccls.NewDeleteAlarmNoticeResponse()
		resp.Response = &tccls.DeleteAlarmNoticeResponseParams{
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyFuncReturn(ratelimit.Check)

	meta := newAlarmNoticeMockMeta()
	res := localcls.ResourceTencentCloudClsAlarmNotice()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"name": "test-alarm-notice",
	})
	d.SetId("notice-abcdefghij")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
}

// TestClsAlarmNotice_Delete_APIError tests Delete handles API error
func TestClsAlarmNotice_Delete_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	clsClient := &tccls.Client{}
	patches.ApplyMethodReturn(newAlarmNoticeMockMeta().client, "UseClsClient", clsClient)

	patches.ApplyMethodFunc(clsClient, "DeleteAlarmNotice", func(request *tccls.DeleteAlarmNoticeRequest) (*tccls.DeleteAlarmNoticeResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=ResourceNotFound, Message=AlarmNotice not found")
	})

	patches.ApplyFuncReturn(ratelimit.Check)

	meta := newAlarmNoticeMockMeta()
	res := localcls.ResourceTencentCloudClsAlarmNotice()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"name": "test-alarm-notice",
	})
	d.SetId("notice-not-exist")

	err := res.Delete(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ResourceNotFound")
}

// TestClsAlarmNotice_Schema_Tags tests tags schema is TypeMap Optional
func TestClsAlarmNotice_Schema_Tags(t *testing.T) {
	res := localcls.ResourceTencentCloudClsAlarmNotice()

	assert.NotNil(t, res)
	assert.Contains(t, res.Schema, "tags")

	tagsSchema := res.Schema["tags"]
	assert.Equal(t, schema.TypeMap, tagsSchema.Type)
	assert.True(t, tagsSchema.Optional)
	assert.False(t, tagsSchema.Required)
	assert.False(t, tagsSchema.Computed)
}

// TestClsAlarmNotice_Schema_Basic verifies basic schema properties
func TestClsAlarmNotice_Schema_Basic(t *testing.T) {
	res := localcls.ResourceTencentCloudClsAlarmNotice()

	assert.NotNil(t, res.Create)
	assert.NotNil(t, res.Read)
	assert.NotNil(t, res.Update)
	assert.NotNil(t, res.Delete)
	assert.NotNil(t, res.Importer)

	assert.Contains(t, res.Schema, "name")
	nameSchema := res.Schema["name"]
	assert.Equal(t, schema.TypeString, nameSchema.Type)
	assert.True(t, nameSchema.Required)

	assert.Contains(t, res.Schema, "tags")
	tagsSchema := res.Schema["tags"]
	assert.Equal(t, schema.TypeMap, tagsSchema.Type)
	assert.True(t, tagsSchema.Optional)
}
