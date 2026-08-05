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

// go test ./tencentcloud/services/teo/ -run "TestInferenceApiTokenV7" -v -count=1 -gcflags="all=-l"
// TestInferenceApiTokenV7_CreateSuccess tests Create calls API and sets ID
func TestInferenceApiTokenV7_CreateSuccess(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "CreateInferenceAPITokenWithContext", func(_ interface{}, request *teov20220901.CreateInferenceAPITokenRequest) (*teov20220901.CreateInferenceAPITokenResponse, error) {
		assert.Equal(t, "zone-12345678", *request.ZoneId)
		assert.Equal(t, "my-token", *request.Name)
		resp := teov20220901.NewCreateInferenceAPITokenResponse()
		resp.Response = &teov20220901.CreateInferenceAPITokenResponseParams{
			TokenId:   ptrString("token-abc123"),
			Content:   ptrString("secret-content-xyz"),
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	// Mock DescribeInferenceAPITokensWithContext for Read
	patches.ApplyMethodFunc(teoClient, "DescribeInferenceAPITokensWithContext", func(_ interface{}, request *teov20220901.DescribeInferenceAPITokensRequest) (*teov20220901.DescribeInferenceAPITokensResponse, error) {
		resp := teov20220901.NewDescribeInferenceAPITokensResponse()
		resp.Response = &teov20220901.DescribeInferenceAPITokensResponseParams{
			TotalCount: ptrInt64(1),
			Tokens: []*teov20220901.InferenceAPIToken{
				{
					TokenId: ptrString("token-abc123"),
					Name:    ptrString("my-token"),
					Content: ptrString("secret-content-xyz"),
				},
			},
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoInferenceApiTokenV7()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-12345678",
		"name":    "my-token",
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "token-abc123", d.Id())
	assert.Equal(t, "my-token", d.Get("name"))
	assert.Equal(t, "token-abc123", d.Get("token_id"))
	assert.Equal(t, "secret-content-xyz", d.Get("content"))
}

// TestInferenceApiTokenV7_CreateEmptyResponse tests Create handles empty API response
func TestInferenceApiTokenV7_CreateEmptyResponse(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "CreateInferenceAPITokenWithContext", func(_ interface{}, request *teov20220901.CreateInferenceAPITokenRequest) (*teov20220901.CreateInferenceAPITokenResponse, error) {
		resp := teov20220901.NewCreateInferenceAPITokenResponse()
		resp.Response = &teov20220901.CreateInferenceAPITokenResponseParams{
			TokenId:   nil,
			Content:   nil,
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoInferenceApiTokenV7()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-12345678",
		"name":    "my-token",
	})

	err := res.Create(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "empty response")
}

// TestInferenceApiTokenV7_CreateAPIError tests Create handles API error
func TestInferenceApiTokenV7_CreateAPIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "CreateInferenceAPITokenWithContext", func(_ interface{}, request *teov20220901.CreateInferenceAPITokenRequest) (*teov20220901.CreateInferenceAPITokenResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=InvalidParameter, Message=Invalid zone_id")
	})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoInferenceApiTokenV7()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-invalid",
		"name":    "my-token",
	})

	err := res.Create(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "InvalidParameter")
}

// TestInferenceApiTokenV7_ReadSuccess tests Read sets state fields correctly
func TestInferenceApiTokenV7_ReadSuccess(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeInferenceAPITokensWithContext", func(_ interface{}, request *teov20220901.DescribeInferenceAPITokensRequest) (*teov20220901.DescribeInferenceAPITokensResponse, error) {
		resp := teov20220901.NewDescribeInferenceAPITokensResponse()
		resp.Response = &teov20220901.DescribeInferenceAPITokensResponseParams{
			TotalCount: ptrInt64(1),
			Tokens: []*teov20220901.InferenceAPIToken{
				{
					TokenId: ptrString("token-abc123"),
					Name:    ptrString("my-token"),
					Content: ptrString("secret-content-xyz"),
				},
			},
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoInferenceApiTokenV7()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-12345678",
		"name":    "my-token",
	})
	d.SetId("token-abc123")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "token-abc123", d.Id())
	assert.Equal(t, "my-token", d.Get("name"))
	assert.Equal(t, "token-abc123", d.Get("token_id"))
	assert.Equal(t, "secret-content-xyz", d.Get("content"))
}

// TestInferenceApiTokenV7_ReadNotFound tests Read clears ID when token not found
func TestInferenceApiTokenV7_ReadNotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeInferenceAPITokensWithContext", func(_ interface{}, request *teov20220901.DescribeInferenceAPITokensRequest) (*teov20220901.DescribeInferenceAPITokensResponse, error) {
		resp := teov20220901.NewDescribeInferenceAPITokensResponse()
		resp.Response = &teov20220901.DescribeInferenceAPITokensResponseParams{
			TotalCount: ptrInt64(0),
			Tokens:     []*teov20220901.InferenceAPIToken{},
			RequestId:  ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoInferenceApiTokenV7()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-12345678",
		"name":    "my-token",
	})
	d.SetId("token-notfound")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}

// TestInferenceApiTokenV7_DeleteSuccess tests Delete calls API successfully
func TestInferenceApiTokenV7_DeleteSuccess(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "DeleteInferenceAPITokenWithContext", func(_ interface{}, request *teov20220901.DeleteInferenceAPITokenRequest) (*teov20220901.DeleteInferenceAPITokenResponse, error) {
		assert.Equal(t, "zone-12345678", *request.ZoneId)
		assert.Equal(t, "token-abc123", *request.TokenId)
		resp := teov20220901.NewDeleteInferenceAPITokenResponse()
		resp.Response = &teov20220901.DeleteInferenceAPITokenResponseParams{
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoInferenceApiTokenV7()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-12345678",
		"name":    "my-token",
	})
	d.SetId("token-abc123")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
}

// TestInferenceApiTokenV7_DeleteAPIError tests Delete handles API error
func TestInferenceApiTokenV7_DeleteAPIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "DeleteInferenceAPITokenWithContext", func(_ interface{}, request *teov20220901.DeleteInferenceAPITokenRequest) (*teov20220901.DeleteInferenceAPITokenResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=ResourceNotFound, Message=Token not found")
	})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoInferenceApiTokenV7()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-12345678",
		"name":    "my-token",
	})
	d.SetId("token-notfound")

	err := res.Delete(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ResourceNotFound")
}

// TestInferenceApiTokenV7_Schema validates schema definition
func TestInferenceApiTokenV7_Schema(t *testing.T) {
	res := teo.ResourceTencentCloudTeoInferenceApiTokenV7()

	assert.NotNil(t, res)
	assert.NotNil(t, res.Create)
	assert.NotNil(t, res.Read)
	assert.Nil(t, res.Update)
	assert.NotNil(t, res.Delete)
	assert.NotNil(t, res.Importer)

	assert.Contains(t, res.Schema, "zone_id")
	assert.Contains(t, res.Schema, "name")
	assert.Contains(t, res.Schema, "token_id")
	assert.Contains(t, res.Schema, "content")

	zoneId := res.Schema["zone_id"]
	assert.Equal(t, schema.TypeString, zoneId.Type)
	assert.True(t, zoneId.Required)
	assert.True(t, zoneId.ForceNew)

	name := res.Schema["name"]
	assert.Equal(t, schema.TypeString, name.Type)
	assert.True(t, name.Required)
	assert.True(t, name.ForceNew)

	tokenId := res.Schema["token_id"]
	assert.Equal(t, schema.TypeString, tokenId.Type)
	assert.True(t, tokenId.Computed)

	content := res.Schema["content"]
	assert.Equal(t, schema.TypeString, content.Type)
	assert.True(t, content.Computed)
	assert.True(t, content.Sensitive)
}
