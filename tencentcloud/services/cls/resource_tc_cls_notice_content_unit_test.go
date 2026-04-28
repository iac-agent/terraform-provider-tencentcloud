package cls_test

import (
	"fmt"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	clsv20201016 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cls/v20201016"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/cls"
)

// mockClsMeta implements tccommon.ProviderMeta
type mockClsMeta struct {
	client *connectivity.TencentCloudClient
}

func (m *mockClsMeta) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockClsMeta{}

func newMockClsMeta() *mockClsMeta {
	return &mockClsMeta{client: &connectivity.TencentCloudClient{}}
}

func ptrStr(s string) *string {
	return &s
}

func ptrUint64(v uint64) *uint64 {
	return &v
}

// go test ./tencentcloud/services/cls/ -run "TestClsNoticeContent" -v -count=1 -gcflags="all=-l"

// TestClsNoticeContent_Schema validates the notice_content_id computed attribute exists
func TestClsNoticeContent_Schema(t *testing.T) {
	res := cls.ResourceTencentCloudClsNoticeContent()

	assert.NotNil(t, res)
	assert.Contains(t, res.Schema, "notice_content_id")

	noticeContentIdSchema := res.Schema["notice_content_id"]
	assert.Equal(t, schema.TypeString, noticeContentIdSchema.Type)
	assert.True(t, noticeContentIdSchema.Computed)
	assert.False(t, noticeContentIdSchema.Optional)
	assert.False(t, noticeContentIdSchema.Required)
}

// TestClsNoticeContent_Create_SetsNoticeContentId tests that Create sets notice_content_id
func TestClsNoticeContent_Create_SetsNoticeContentId(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	clsClient := &clsv20201016.Client{}
	patches.ApplyMethodReturn(newMockClsMeta().client, "UseClsV20201016Client", clsClient)

	patches.ApplyMethodFunc(clsClient, "CreateNoticeContentWithContext",
		func(ctx interface{}, request *clsv20201016.CreateNoticeContentRequest) (*clsv20201016.CreateNoticeContentResponse, error) {
			assert.Equal(t, "tf-test-notice", *request.Name)
			resp := clsv20201016.NewCreateNoticeContentResponse()
			resp.Response = &clsv20201016.CreateNoticeContentResponseParams{
				NoticeContentId: ptrStr("noticetemplate-abc123"),
				RequestId:       ptrStr("fake-request-id"),
			}
			return resp, nil
		},
	)

	// Mock the DescribeClsNoticeContentById call that happens during Read after Create
	patches.ApplyMethodFunc(clsClient, "DescribeNoticeContents",
		func(request *clsv20201016.DescribeNoticeContentsRequest) (*clsv20201016.DescribeNoticeContentsResponse, error) {
			resp := clsv20201016.NewDescribeNoticeContentsResponse()
			resp.Response = &clsv20201016.DescribeNoticeContentsResponseParams{
				NoticeContents: []*clsv20201016.NoticeContentTemplate{
					{
						NoticeContentId: ptrStr("noticetemplate-abc123"),
						Name:            ptrStr("tf-test-notice"),
						Type:            ptrUint64(0),
					},
				},
				TotalCount: ptrInt64(1),
				RequestId:  ptrStr("fake-request-id"),
			}
			return resp, nil
		},
	)

	meta := newMockClsMeta()
	res := cls.ResourceTencentCloudClsNoticeContent()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"name": "tf-test-notice",
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "noticetemplate-abc123", d.Id())
	assert.Equal(t, "noticetemplate-abc123", d.Get("notice_content_id"))
}

// TestClsNoticeContent_Read_SetsNoticeContentId tests that Read sets notice_content_id from API response
func TestClsNoticeContent_Read_SetsNoticeContentId(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	clsClient := &clsv20201016.Client{}
	patches.ApplyMethodReturn(newMockClsMeta().client, "UseClsV20201016Client", clsClient)

	patches.ApplyMethodFunc(clsClient, "DescribeNoticeContents",
		func(request *clsv20201016.DescribeNoticeContentsRequest) (*clsv20201016.DescribeNoticeContentsResponse, error) {
			resp := clsv20201016.NewDescribeNoticeContentsResponse()
			resp.Response = &clsv20201016.DescribeNoticeContentsResponseParams{
				NoticeContents: []*clsv20201016.NoticeContentTemplate{
					{
						NoticeContentId: ptrStr("noticetemplate-xyz789"),
						Name:            ptrStr("tf-test-notice-read"),
						Type:            ptrUint64(1),
					},
				},
				TotalCount: ptrInt64(1),
				RequestId:  ptrStr("fake-request-id"),
			}
			return resp, nil
		},
	)

	meta := newMockClsMeta()
	res := cls.ResourceTencentCloudClsNoticeContent()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"name": "tf-test-notice-read",
	})
	d.SetId("noticetemplate-xyz789")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "noticetemplate-xyz789", d.Get("notice_content_id"))
	assert.Equal(t, "tf-test-notice-read", d.Get("name"))
}

// TestClsNoticeContent_Create_APIError tests that Create handles API error
func TestClsNoticeContent_Create_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	clsClient := &clsv20201016.Client{}
	patches.ApplyMethodReturn(newMockClsMeta().client, "UseClsV20201016Client", clsClient)

	patches.ApplyMethodFunc(clsClient, "CreateNoticeContentWithContext",
		func(ctx interface{}, request *clsv20201016.CreateNoticeContentRequest) (*clsv20201016.CreateNoticeContentResponse, error) {
			return nil, fmt.Errorf("[TencentCloudSDKError] Code=InvalidParameter, Message=Invalid name")
		},
	)

	meta := newMockClsMeta()
	res := cls.ResourceTencentCloudClsNoticeContent()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"name": "invalid-name",
	})

	err := res.Create(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "InvalidParameter")
}

// TestClsNoticeContent_Delete_Success tests that Delete calls API properly
func TestClsNoticeContent_Delete_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	clsClient := &clsv20201016.Client{}
	patches.ApplyMethodReturn(newMockClsMeta().client, "UseClsV20201016Client", clsClient)

	patches.ApplyMethodFunc(clsClient, "DeleteNoticeContentWithContext",
		func(ctx interface{}, request *clsv20201016.DeleteNoticeContentRequest) (*clsv20201016.DeleteNoticeContentResponse, error) {
			assert.Equal(t, "noticetemplate-abc123", *request.NoticeContentId)
			resp := clsv20201016.NewDeleteNoticeContentResponse()
			resp.Response = &clsv20201016.DeleteNoticeContentResponseParams{
				RequestId: ptrStr("fake-request-id"),
			}
			return resp, nil
		},
	)

	meta := newMockClsMeta()
	res := cls.ResourceTencentCloudClsNoticeContent()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"name": "tf-test-notice",
	})
	d.SetId("noticetemplate-abc123")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
}

func ptrInt64(v int64) *int64 {
	return &v
}
