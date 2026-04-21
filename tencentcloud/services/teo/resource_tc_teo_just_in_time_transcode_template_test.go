package teo_test

import (
	"fmt"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/teo"
)

// go test ./tencentcloud/services/teo/ -run "TestJustInTimeTranscodeTemplate" -v -count=1 -gcflags="all=-l"

func TestJustInTimeTranscodeTemplate_Schema(t *testing.T) {
	res := teo.ResourceTencentCloudTeoJustInTimeTranscodeTemplate()

	assert.NotNil(t, res)
	assert.NotNil(t, res.Create)
	assert.NotNil(t, res.Read)
	assert.Nil(t, res.Update)
	assert.NotNil(t, res.Delete)
	assert.NotNil(t, res.Importer)

	// Verify required fields
	assert.Contains(t, res.Schema, "zone_id")
	assert.Contains(t, res.Schema, "template_name")

	// Verify ForceNew on all non-computed fields
	forceNewFields := []string{"zone_id", "template_name", "comment", "video_stream_switch", "audio_stream_switch", "video_template", "audio_template"}
	for _, field := range forceNewFields {
		s, ok := res.Schema[field]
		assert.True(t, ok, "field %s should exist in schema", field)
		if ok {
			assert.True(t, s.ForceNew, "field %s should have ForceNew=true", field)
		}
	}

	// Verify computed fields
	computedFields := []string{"template_id", "create_time", "update_time", "type"}
	for _, field := range computedFields {
		s, ok := res.Schema[field]
		assert.True(t, ok, "field %s should exist in schema", field)
		if ok {
			assert.True(t, s.Computed, "field %s should be Computed", field)
		}
	}

	// Verify nested structures
	vt, ok := res.Schema["video_template"]
	assert.True(t, ok)
	if ok {
		assert.Equal(t, schema.TypeList, vt.Type)
		assert.Equal(t, 1, vt.MaxItems)
	}

	at, ok := res.Schema["audio_template"]
	assert.True(t, ok)
	if ok {
		assert.Equal(t, schema.TypeList, at.Type)
		assert.Equal(t, 1, at.MaxItems)
	}
}

func TestJustInTimeTranscodeTemplate_CreateSuccess(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(&connectivity.TencentCloudClient{}, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "CreateJustInTimeTranscodeTemplateWithContext", func(ctx interface{}, request *teov20220901.CreateJustInTimeTranscodeTemplateRequest) (*teov20220901.CreateJustInTimeTranscodeTemplateResponse, error) {
		resp := teov20220901.NewCreateJustInTimeTranscodeTemplateResponse()
		resp.Response = &teov20220901.CreateJustInTimeTranscodeTemplateResponseParams{
			TemplateId: ptrString("C1LZ7982VgTpYhJ7M"),
			RequestId:  ptrString("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(teoClient, "DescribeJustInTimeTranscodeTemplates", func(request *teov20220901.DescribeJustInTimeTranscodeTemplatesRequest) (*teov20220901.DescribeJustInTimeTranscodeTemplatesResponse, error) {
		resp := teov20220901.NewDescribeJustInTimeTranscodeTemplatesResponse()
		resp.Response = &teov20220901.DescribeJustInTimeTranscodeTemplatesResponseParams{
			TotalCount: ptrUint64(1),
			TemplateSet: []*teov20220901.JustInTimeTranscodeTemplate{
				{
					TemplateId:        ptrString("C1LZ7982VgTpYhJ7M"),
					TemplateName:      ptrString("test-template"),
					Comment:           ptrString("test comment"),
					Type:              ptrString("custom"),
					VideoStreamSwitch: ptrString("on"),
					AudioStreamSwitch: ptrString("on"),
					VideoTemplate: &teov20220901.VideoTemplateInfo{
						Codec:              ptrString("H.264"),
						Fps:                ptrFloat64(30),
						Bitrate:            ptrUint64(2000),
						ResolutionAdaptive: ptrString("open"),
						Width:              ptrUint64(1920),
						Height:             ptrUint64(1080),
						FillType:           ptrString("black"),
					},
					AudioTemplate: &teov20220901.AudioTemplateInfo{
						Codec:        ptrString("libfdk_aac"),
						AudioChannel: ptrUint64(2),
					},
					CreateTime: ptrString("2024-01-01T00:00:00Z"),
					UpdateTime: ptrString("2024-01-01T00:00:00Z"),
				},
			},
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoJustInTimeTranscodeTemplate()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":             "zone-3edjdliiw3he",
		"template_name":       "test-template",
		"comment":             "test comment",
		"video_stream_switch": "on",
		"audio_stream_switch": "on",
		"video_template": []interface{}{
			map[string]interface{}{
				"codec":               "H.264",
				"fps":                 30.0,
				"bitrate":             2000,
				"resolution_adaptive": "open",
				"width":               1920,
				"height":              1080,
				"fill_type":           "black",
			},
		},
		"audio_template": []interface{}{
			map[string]interface{}{
				"codec":         "libfdk_aac",
				"audio_channel": 2,
			},
		},
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-3edjdliiw3he#C1LZ7982VgTpYhJ7M", d.Id())
	assert.Equal(t, "C1LZ7982VgTpYhJ7M", d.Get("template_id"))
	assert.Equal(t, "test-template", d.Get("template_name"))
	assert.Equal(t, "custom", d.Get("type"))
}

func TestJustInTimeTranscodeTemplate_CreateAPIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(&connectivity.TencentCloudClient{}, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "CreateJustInTimeTranscodeTemplateWithContext", func(ctx interface{}, request *teov20220901.CreateJustInTimeTranscodeTemplateRequest) (*teov20220901.CreateJustInTimeTranscodeTemplateResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=InvalidParameter, Message=Invalid zone_id")
	})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoJustInTimeTranscodeTemplate()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":       "zone-invalid",
		"template_name": "test-template",
	})

	err := res.Create(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "InvalidParameter")
}

func TestJustInTimeTranscodeTemplate_ReadSuccess(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(&connectivity.TencentCloudClient{}, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeJustInTimeTranscodeTemplates", func(request *teov20220901.DescribeJustInTimeTranscodeTemplatesRequest) (*teov20220901.DescribeJustInTimeTranscodeTemplatesResponse, error) {
		resp := teov20220901.NewDescribeJustInTimeTranscodeTemplatesResponse()
		resp.Response = &teov20220901.DescribeJustInTimeTranscodeTemplatesResponseParams{
			TotalCount: ptrUint64(1),
			TemplateSet: []*teov20220901.JustInTimeTranscodeTemplate{
				{
					TemplateId:        ptrString("C1LZ7982VgTpYhJ7M"),
					TemplateName:      ptrString("test-template"),
					Comment:           ptrString("test comment"),
					Type:              ptrString("custom"),
					VideoStreamSwitch: ptrString("on"),
					AudioStreamSwitch: ptrString("on"),
					VideoTemplate: &teov20220901.VideoTemplateInfo{
						Codec:              ptrString("H.264"),
						Fps:                ptrFloat64(30),
						Bitrate:            ptrUint64(2000),
						ResolutionAdaptive: ptrString("open"),
						Width:              ptrUint64(1920),
						Height:             ptrUint64(1080),
						FillType:           ptrString("black"),
					},
					AudioTemplate: &teov20220901.AudioTemplateInfo{
						Codec:        ptrString("libfdk_aac"),
						AudioChannel: ptrUint64(2),
					},
					CreateTime: ptrString("2024-01-01T00:00:00Z"),
					UpdateTime: ptrString("2024-01-01T00:00:00Z"),
				},
			},
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoJustInTimeTranscodeTemplate()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":       "zone-3edjdliiw3he",
		"template_name": "test-template",
	})
	d.SetId("zone-3edjdliiw3he#C1LZ7982VgTpYhJ7M")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-3edjdliiw3he#C1LZ7982VgTpYhJ7M", d.Id())
	assert.Equal(t, "C1LZ7982VgTpYhJ7M", d.Get("template_id"))
	assert.Equal(t, "test-template", d.Get("template_name"))
	assert.Equal(t, "custom", d.Get("type"))
	assert.Equal(t, "on", d.Get("video_stream_switch"))
	assert.Equal(t, "on", d.Get("audio_stream_switch"))
	assert.Equal(t, "test comment", d.Get("comment"))
	assert.Equal(t, "2024-01-01T00:00:00Z", d.Get("create_time"))
	assert.Equal(t, "2024-01-01T00:00:00Z", d.Get("update_time"))
}

func TestJustInTimeTranscodeTemplate_ReadNotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(&connectivity.TencentCloudClient{}, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeJustInTimeTranscodeTemplates", func(request *teov20220901.DescribeJustInTimeTranscodeTemplatesRequest) (*teov20220901.DescribeJustInTimeTranscodeTemplatesResponse, error) {
		resp := teov20220901.NewDescribeJustInTimeTranscodeTemplatesResponse()
		resp.Response = &teov20220901.DescribeJustInTimeTranscodeTemplatesResponseParams{
			TotalCount:  ptrUint64(0),
			TemplateSet: []*teov20220901.JustInTimeTranscodeTemplate{},
			RequestId:   ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoJustInTimeTranscodeTemplate()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":       "zone-3edjdliiw3he",
		"template_name": "test-template",
	})
	d.SetId("zone-3edjdliiw3he#C1LZ7982VgTpYhJ7M")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}

func TestJustInTimeTranscodeTemplate_DeleteSuccess(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(&connectivity.TencentCloudClient{}, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DeleteJustInTimeTranscodeTemplatesWithContext", func(ctx interface{}, request *teov20220901.DeleteJustInTimeTranscodeTemplatesRequest) (*teov20220901.DeleteJustInTimeTranscodeTemplatesResponse, error) {
		assert.Equal(t, "zone-3edjdliiw3he", *request.ZoneId)
		assert.Len(t, request.TemplateIds, 1)
		assert.Equal(t, "C1LZ7982VgTpYhJ7M", *request.TemplateIds[0])

		resp := teov20220901.NewDeleteJustInTimeTranscodeTemplatesResponse()
		resp.Response = &teov20220901.DeleteJustInTimeTranscodeTemplatesResponseParams{
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoJustInTimeTranscodeTemplate()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":       "zone-3edjdliiw3he",
		"template_name": "test-template",
	})
	d.SetId("zone-3edjdliiw3he#C1LZ7982VgTpYhJ7M")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
}

func TestJustInTimeTranscodeTemplate_DeleteAPIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(&connectivity.TencentCloudClient{}, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DeleteJustInTimeTranscodeTemplatesWithContext", func(ctx interface{}, request *teov20220901.DeleteJustInTimeTranscodeTemplatesRequest) (*teov20220901.DeleteJustInTimeTranscodeTemplatesResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=ResourceNotFound, Message=Template not found")
	})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoJustInTimeTranscodeTemplate()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":       "zone-3edjdliiw3he",
		"template_name": "test-template",
	})
	d.SetId("zone-3edjdliiw3he#C1LZ7982VgTpYhJ7M")

	err := res.Delete(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ResourceNotFound")
}

func ptrUint64(i uint64) *uint64 {
	return &i
}

func ptrFloat64(f float64) *float64 {
	return &f
}
