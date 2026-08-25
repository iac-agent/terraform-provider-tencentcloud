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

type mockMetaInferenceApiToken struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaInferenceApiToken) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaInferenceApiToken{}

func newMockMetaInferenceApiToken() *mockMetaInferenceApiToken {
	return &mockMetaInferenceApiToken{client: &connectivity.TencentCloudClient{}}
}

func ptrStringIAT(s string) *string {
	return &s
}

func ptrInt64IAT(i int64) *int64 {
	return &i
}

func TestTeoInferenceApiToken_Create_Normal(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaInferenceApiToken().client, "UseTeoV20220901Client", teoClient)

	// Mock Create API
	patches.ApplyMethodFunc(teoClient, "CreateInferenceAPITokenWithContext", func(_ context.Context, request *teov20220901.CreateInferenceAPITokenRequest) (*teov20220901.CreateInferenceAPITokenResponse, error) {
		assert.NotNil(t, request.ZoneId)
		assert.Equal(t, "zone-test123", *request.ZoneId)
		assert.NotNil(t, request.Name)
		assert.Equal(t, "tf-example", *request.Name)

		resp := teov20220901.NewCreateInferenceAPITokenResponse()
		resp.Response = &teov20220901.CreateInferenceAPITokenResponseParams{
			TokenId:   ptrStringIAT("token-abcd1234"),
			Content:   ptrStringIAT("eo-token-content"),
			RequestId: ptrStringIAT("fake-request-id-create"),
		}
		return resp, nil
	})

	// Mock Read API (called after Create)
	patches.ApplyMethodFunc(teoClient, "DescribeInferenceAPITokensWithContext", func(_ context.Context, request *teov20220901.DescribeInferenceAPITokensRequest) (*teov20220901.DescribeInferenceAPITokensResponse, error) {
		assert.NotNil(t, request.ZoneId)
		assert.Equal(t, "zone-test123", *request.ZoneId)
		assert.NotNil(t, request.Limit)
		assert.Equal(t, int64(100), *request.Limit)

		resp := teov20220901.NewDescribeInferenceAPITokensResponse()
		resp.Response = &teov20220901.DescribeInferenceAPITokensResponseParams{
			TotalCount: ptrInt64IAT(1),
			Tokens: []*teov20220901.InferenceAPIToken{
				{
					TokenId:    ptrStringIAT("token-abcd1234"),
					Name:       ptrStringIAT("tf-example"),
					Content:    ptrStringIAT("eo-token-content"),
					CreateTime: ptrStringIAT("2024-01-01T00:00:00Z"),
				},
			},
			RequestId: ptrStringIAT("fake-request-id-read"),
		}
		return resp, nil
	})

	meta := newMockMetaInferenceApiToken()
	res := teo.ResourceTencentCloudTeoInferenceApiToken()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-test123",
		"name":    "tf-example",
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-test123#token-abcd1234", d.Id())
	assert.Equal(t, "zone-test123", d.Get("zone_id"))
	assert.Equal(t, "tf-example", d.Get("name"))
	assert.Equal(t, "token-abcd1234", d.Get("token_id"))
	assert.Equal(t, "eo-token-content", d.Get("content"))
	assert.Equal(t, "2024-01-01T00:00:00Z", d.Get("create_time"))
}

func TestTeoInferenceApiToken_Create_EmptyTokenId(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaInferenceApiToken().client, "UseTeoV20220901Client", teoClient)

	// Mock Create API returning an empty TokenId
	patches.ApplyMethodFunc(teoClient, "CreateInferenceAPITokenWithContext", func(_ context.Context, request *teov20220901.CreateInferenceAPITokenRequest) (*teov20220901.CreateInferenceAPITokenResponse, error) {
		resp := teov20220901.NewCreateInferenceAPITokenResponse()
		resp.Response = &teov20220901.CreateInferenceAPITokenResponseParams{
			TokenId:   ptrStringIAT(""),
			Content:   ptrStringIAT(""),
			RequestId: ptrStringIAT("fake-request-id-create"),
		}
		return resp, nil
	})

	meta := newMockMetaInferenceApiToken()
	res := teo.ResourceTencentCloudTeoInferenceApiToken()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-test123",
		"name":    "tf-example",
	})

	err := res.Create(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "TokenId is empty")
	// ensure no empty id was written
	assert.Equal(t, "", d.Id())
}

func TestTeoInferenceApiToken_Read_Normal(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaInferenceApiToken().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeInferenceAPITokensWithContext", func(_ context.Context, request *teov20220901.DescribeInferenceAPITokensRequest) (*teov20220901.DescribeInferenceAPITokensResponse, error) {
		resp := teov20220901.NewDescribeInferenceAPITokensResponse()
		resp.Response = &teov20220901.DescribeInferenceAPITokensResponseParams{
			TotalCount: ptrInt64IAT(1),
			Tokens: []*teov20220901.InferenceAPIToken{
				{
					TokenId:    ptrStringIAT("token-abcd1234"),
					Name:       ptrStringIAT("tf-example"),
					Content:    ptrStringIAT("eo-token-content"),
					CreateTime: ptrStringIAT("2024-01-01T00:00:00Z"),
				},
			},
			RequestId: ptrStringIAT("fake-request-id-read"),
		}
		return resp, nil
	})

	meta := newMockMetaInferenceApiToken()
	res := teo.ResourceTencentCloudTeoInferenceApiToken()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-test123",
		"name":    "tf-example",
	})
	d.SetId("zone-test123#token-abcd1234")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-test123#token-abcd1234", d.Id())
	assert.Equal(t, "zone-test123", d.Get("zone_id"))
	assert.Equal(t, "tf-example", d.Get("name"))
	assert.Equal(t, "token-abcd1234", d.Get("token_id"))
	assert.Equal(t, "eo-token-content", d.Get("content"))
	assert.Equal(t, "2024-01-01T00:00:00Z", d.Get("create_time"))
}

func TestTeoInferenceApiToken_Read_NotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaInferenceApiToken().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeInferenceAPITokensWithContext", func(_ context.Context, request *teov20220901.DescribeInferenceAPITokensRequest) (*teov20220901.DescribeInferenceAPITokensResponse, error) {
		resp := teov20220901.NewDescribeInferenceAPITokensResponse()
		resp.Response = &teov20220901.DescribeInferenceAPITokensResponseParams{
			TotalCount: ptrInt64IAT(0),
			Tokens:     []*teov20220901.InferenceAPIToken{},
			RequestId:  ptrStringIAT("fake-request-id-read"),
		}
		return resp, nil
	})

	meta := newMockMetaInferenceApiToken()
	res := teo.ResourceTencentCloudTeoInferenceApiToken()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-test123",
		"name":    "tf-example",
	})
	d.SetId("zone-test123#token-notexist")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}

func TestTeoInferenceApiToken_Delete_Normal(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaInferenceApiToken().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DeleteInferenceAPITokenWithContext", func(_ context.Context, request *teov20220901.DeleteInferenceAPITokenRequest) (*teov20220901.DeleteInferenceAPITokenResponse, error) {
		assert.NotNil(t, request.ZoneId)
		assert.Equal(t, "zone-test123", *request.ZoneId)
		assert.NotNil(t, request.TokenId)
		assert.Equal(t, "token-abcd1234", *request.TokenId)

		resp := teov20220901.NewDeleteInferenceAPITokenResponse()
		resp.Response = &teov20220901.DeleteInferenceAPITokenResponseParams{
			RequestId: ptrStringIAT("fake-request-id-delete"),
		}
		return resp, nil
	})

	meta := newMockMetaInferenceApiToken()
	res := teo.ResourceTencentCloudTeoInferenceApiToken()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-test123",
		"name":    "tf-example",
	})
	d.SetId("zone-test123#token-abcd1234")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
}

func TestTeoInferenceApiToken_Update_ImmutableArgsChanged(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaInferenceApiToken().client, "UseTeoV20220901Client", teoClient)

	meta := newMockMetaInferenceApiToken()
	res := teo.ResourceTencentCloudTeoInferenceApiToken()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-test123",
		"name":    "tf-example",
	})
	d.SetId("zone-test123#token-abcd1234")

	// Force name to be detected as changed (immutable field)
	patches.ApplyMethodFunc(d, "HasChange", func(key string) bool {
		return key == "name"
	})

	err := res.Update(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "name")
	assert.Contains(t, err.Error(), "cannot be changed")
}
