package teo_test

import (
	"context"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/teo"
)

type mockMetaInferenceAPIToken struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaInferenceAPIToken) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaInferenceAPIToken{}

func newMockMetaInferenceAPIToken() *mockMetaInferenceAPIToken {
	return &mockMetaInferenceAPIToken{client: &connectivity.TencentCloudClient{}}
}

func ptrStringIAT(s string) *string {
	return &s
}

func ptrInt64IAT(i int64) *int64 {
	return &i
}

func TestTeoInferenceAPIToken_Create_Normal(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaInferenceAPIToken().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "CreateInferenceAPITokenWithContext", func(_ context.Context, request *teov20220901.CreateInferenceAPITokenRequest) (*teov20220901.CreateInferenceAPITokenResponse, error) {
		assert.NotNil(t, request.ZoneId)
		assert.Equal(t, "zone-test123", *request.ZoneId)
		assert.NotNil(t, request.Name)
		assert.Equal(t, "test-token", *request.Name)

		resp := teov20220901.NewCreateInferenceAPITokenResponse()
		resp.Response = &teov20220901.CreateInferenceAPITokenResponseParams{
			TokenId:   ptrStringIAT("token-xxxxx"),
			Content:   ptrStringIAT("test-content"),
			RequestId: ptrStringIAT("fake-request-id-create"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(teoClient, "DescribeInferenceAPITokensWithContext", func(_ context.Context, request *teov20220901.DescribeInferenceAPITokensRequest) (*teov20220901.DescribeInferenceAPITokensResponse, error) {
		assert.NotNil(t, request.ZoneId)
		assert.Equal(t, "zone-test123", *request.ZoneId)
		assert.NotNil(t, request.Offset)
		assert.Equal(t, int64(0), *request.Offset)
		assert.NotNil(t, request.Limit)
		assert.Equal(t, int64(20), *request.Limit)

		resp := teov20220901.NewDescribeInferenceAPITokensResponse()
		resp.Response = &teov20220901.DescribeInferenceAPITokensResponseParams{
			TotalCount: ptrInt64IAT(1),
			Tokens: []*teov20220901.InferenceAPIToken{
				{
					TokenId:    ptrStringIAT("token-xxxxx"),
					Name:       ptrStringIAT("test-token"),
					Content:    ptrStringIAT("test-content"),
					CreateTime: ptrStringIAT("2024-01-01T00:00:00Z"),
				},
			},
			RequestId: ptrStringIAT("fake-request-id-read"),
		}
		return resp, nil
	})

	meta := newMockMetaInferenceAPIToken()
	res := teo.ResourceTencentCloudTeoInferenceAPIToken()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-test123",
		"name":    "test-token",
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-test123#token-xxxxx", d.Id())
}

func TestTeoInferenceAPIToken_Read_Normal(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaInferenceAPIToken().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeInferenceAPITokensWithContext", func(_ context.Context, request *teov20220901.DescribeInferenceAPITokensRequest) (*teov20220901.DescribeInferenceAPITokensResponse, error) {
		assert.NotNil(t, request.ZoneId)
		assert.Equal(t, "zone-test123", *request.ZoneId)
		assert.NotNil(t, request.Offset)
		assert.Equal(t, int64(10), *request.Offset)
		assert.NotNil(t, request.Limit)
		assert.Equal(t, int64(50), *request.Limit)

		resp := teov20220901.NewDescribeInferenceAPITokensResponse()
		resp.Response = &teov20220901.DescribeInferenceAPITokensResponseParams{
			TotalCount: ptrInt64IAT(42),
			Tokens: []*teov20220901.InferenceAPIToken{
				{
					TokenId:    ptrStringIAT("token-xxxxx"),
					Name:       ptrStringIAT("test-token"),
					Content:    ptrStringIAT("test-content"),
					CreateTime: ptrStringIAT("2024-01-01T00:00:00Z"),
				},
			},
			RequestId: ptrStringIAT("fake-request-id-read"),
		}
		return resp, nil
	})

	meta := newMockMetaInferenceAPIToken()
	res := teo.ResourceTencentCloudTeoInferenceAPIToken()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-test123",
		"name":    "test-token",
		"offset":  10,
		"limit":   50,
	})
	d.SetId("zone-test123#token-xxxxx")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-test123#token-xxxxx", d.Id())
	assert.Equal(t, "zone-test123", d.Get("zone_id"))
	assert.Equal(t, "test-token", d.Get("name"))
	assert.Equal(t, "token-xxxxx", d.Get("token_id"))
	assert.Equal(t, "test-content", d.Get("content"))
	assert.Equal(t, "2024-01-01T00:00:00Z", d.Get("create_time"))
	assert.Equal(t, 42, d.Get("total_count"))
}

func TestTeoInferenceAPIToken_Read_WithDefaultOffsetLimit(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaInferenceAPIToken().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeInferenceAPITokensWithContext", func(_ context.Context, request *teov20220901.DescribeInferenceAPITokensRequest) (*teov20220901.DescribeInferenceAPITokensResponse, error) {
		assert.NotNil(t, request.Offset)
		assert.Equal(t, int64(0), *request.Offset)
		assert.NotNil(t, request.Limit)
		assert.Equal(t, int64(20), *request.Limit)

		resp := teov20220901.NewDescribeInferenceAPITokensResponse()
		resp.Response = &teov20220901.DescribeInferenceAPITokensResponseParams{
			TotalCount: ptrInt64IAT(1),
			Tokens: []*teov20220901.InferenceAPIToken{
				{
					TokenId:    ptrStringIAT("token-xxxxx"),
					Name:       ptrStringIAT("test-token"),
					Content:    ptrStringIAT("test-content"),
					CreateTime: ptrStringIAT("2024-01-01T00:00:00Z"),
				},
			},
			RequestId: ptrStringIAT("fake-request-id-read"),
		}
		return resp, nil
	})

	meta := newMockMetaInferenceAPIToken()
	res := teo.ResourceTencentCloudTeoInferenceAPIToken()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-test123",
		"name":    "test-token",
	})
	d.SetId("zone-test123#token-xxxxx")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-test123#token-xxxxx", d.Id())
	assert.Equal(t, 1, d.Get("total_count"))
}

func TestTeoInferenceAPIToken_Read_NotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaInferenceAPIToken().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeInferenceAPITokensWithContext", func(_ context.Context, request *teov20220901.DescribeInferenceAPITokensRequest) (*teov20220901.DescribeInferenceAPITokensResponse, error) {
		resp := teov20220901.NewDescribeInferenceAPITokensResponse()
		resp.Response = &teov20220901.DescribeInferenceAPITokensResponseParams{
			TotalCount: ptrInt64IAT(0),
			Tokens:     []*teov20220901.InferenceAPIToken{},
			RequestId:  ptrStringIAT("fake-request-id-read"),
		}
		return resp, nil
	})

	meta := newMockMetaInferenceAPIToken()
	res := teo.ResourceTencentCloudTeoInferenceAPIToken()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-test123",
		"name":    "test-token-notexist",
	})
	d.SetId("zone-test123#token-notexist")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}

func TestTeoInferenceAPIToken_Update_Immutable(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaInferenceAPIToken().client, "UseTeoV20220901Client", teoClient)

	meta := newMockMetaInferenceAPIToken()
	res := teo.ResourceTencentCloudTeoInferenceAPIToken()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-test123",
		"name":    "new-name",
	})
	d.SetId("zone-test123#token-xxxxx")

	// Mark name as changed
	_ = d.Set("name", "old-name")
	d.SetId("zone-test123#token-xxxxx")

	err := res.Update(d, meta)
	assert.NoError(t, err)
}

func TestTeoInferenceAPIToken_Delete_Normal(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaInferenceAPIToken().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DeleteInferenceAPITokenWithContext", func(_ context.Context, request *teov20220901.DeleteInferenceAPITokenRequest) (*teov20220901.DeleteInferenceAPITokenResponse, error) {
		assert.NotNil(t, request.ZoneId)
		assert.Equal(t, "zone-test123", *request.ZoneId)
		assert.NotNil(t, request.TokenId)
		assert.Equal(t, "token-xxxxx", *request.TokenId)

		resp := teov20220901.NewDeleteInferenceAPITokenResponse()
		resp.Response = &teov20220901.DeleteInferenceAPITokenResponseParams{
			RequestId: ptrStringIAT("fake-request-id-delete"),
		}
		return resp, nil
	})

	meta := newMockMetaInferenceAPIToken()
	res := teo.ResourceTencentCloudTeoInferenceAPIToken()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-test123",
		"name":    "test-token",
	})
	d.SetId("zone-test123#token-xxxxx")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
}

func TestTeoInferenceAPIToken_Read_TokenIdMismatch(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaInferenceAPIToken().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeInferenceAPITokensWithContext", func(_ context.Context, request *teov20220901.DescribeInferenceAPITokensRequest) (*teov20220901.DescribeInferenceAPITokensResponse, error) {
		resp := teov20220901.NewDescribeInferenceAPITokensResponse()
		resp.Response = &teov20220901.DescribeInferenceAPITokensResponseParams{
			TotalCount: ptrInt64IAT(1),
			Tokens: []*teov20220901.InferenceAPIToken{
				{
					TokenId:    ptrStringIAT("token-other"),
					Name:       ptrStringIAT("other-token"),
					Content:    ptrStringIAT("other-content"),
					CreateTime: ptrStringIAT("2024-01-01T00:00:00Z"),
				},
			},
			RequestId: ptrStringIAT("fake-request-id-read"),
		}
		return resp, nil
	})

	meta := newMockMetaInferenceAPIToken()
	res := teo.ResourceTencentCloudTeoInferenceAPIToken()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-test123",
		"name":    "test-token",
	})
	d.SetId("zone-test123#token-xxxxx")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}
