package teo_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/teo"
)

// mockMetaForInferenceService implements tccommon.ProviderMeta
type mockMetaForInferenceService struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaForInferenceService) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaForInferenceService{}

func newMockMetaForInferenceService() *mockMetaForInferenceService {
	return &mockMetaForInferenceService{client: &connectivity.TencentCloudClient{}}
}

func ptrStrIS(s string) *string {
	return &s
}

func ptrInt64IS(i int64) *int64 {
	return &i
}

// go test ./tencentcloud/services/teo/ -run "TestInferenceServiceV1" -v -count=1 -gcflags="all=-l"

// TestInferenceServiceV1_Create_Success tests Create calls API and sets composite ID
func TestInferenceServiceV1_Create_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForInferenceService().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "CreateInferenceServiceWithContext", func(_ context.Context, request *teov20220901.CreateInferenceServiceRequest) (*teov20220901.CreateInferenceServiceResponse, error) {
		assert.Equal(t, "zone-12345678", *request.ZoneId)
		assert.Equal(t, "test-inference-svc", *request.Name)
		assert.Equal(t, int64(8080), *request.ListenPort)
		resp := teov20220901.NewCreateInferenceServiceResponse()
		resp.Response = &teov20220901.CreateInferenceServiceResponseParams{
			ServiceId: ptrStrIS("svc-abcdefgh"),
			RequestId: ptrStrIS("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(teoClient, "DescribeInferenceServices", func(request *teov20220901.DescribeInferenceServicesRequest) (*teov20220901.DescribeInferenceServicesResponse, error) {
		resp := teov20220901.NewDescribeInferenceServicesResponse()
		resp.Response = &teov20220901.DescribeInferenceServicesResponseParams{
			Services: []*teov20220901.InferenceService{
				{
					ServiceId:  ptrStrIS("svc-abcdefgh"),
					Name:       ptrStrIS("test-inference-svc"),
					ListenPort: ptrInt64IS(8080),
					Status:     ptrStrIS("Deploying"),
				},
			},
			TotalCount: ptrInt64IS(1),
			RequestId:  ptrStrIS("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForInferenceService()
	res := teo.ResourceTencentCloudTeoInferenceServiceV1()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":     "zone-12345678",
		"name":        "test-inference-svc",
		"listen_port": 8080,
		"containers": []interface{}{
			map[string]interface{}{
				"image_type": "TCR",
				"tcr_config": []interface{}{
					map[string]interface{}{
						"tcr_type":    "Personal",
						"image":       "ccr.ccs.tencentyun.com/test/image:v1",
						"region_name": "ap-guangzhou",
					},
				},
			},
		},
		"resource_config": []interface{}{
			map[string]interface{}{
				"scaling_mode": "Manual",
				"manual_instance_config": []interface{}{
					map[string]interface{}{
						"fixed_instance_count": 1,
					},
				},
			},
		},
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-12345678#svc-abcdefgh", d.Id())
}

// TestInferenceServiceV1_Create_WithAffinityConfig tests Create with affinity config
func TestInferenceServiceV1_Create_WithAffinityConfig(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForInferenceService().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "CreateInferenceServiceWithContext", func(_ context.Context, request *teov20220901.CreateInferenceServiceRequest) (*teov20220901.CreateInferenceServiceResponse, error) {
		assert.NotNil(t, request.AffinityConfig)
		assert.Equal(t, "On", *request.AffinityConfig.Switch)
		assert.Equal(t, "SessionId", *request.AffinityConfig.AffinityMode)
		assert.NotNil(t, request.AffinityConfig.SessionIdAffinityConfig)
		assert.Equal(t, "Header", *request.AffinityConfig.SessionIdAffinityConfig.Source)
		assert.Equal(t, "X-Custom-Session-Id", *request.AffinityConfig.SessionIdAffinityConfig.HeaderName)
		resp := teov20220901.NewCreateInferenceServiceResponse()
		resp.Response = &teov20220901.CreateInferenceServiceResponseParams{
			ServiceId: ptrStrIS("svc-affinity-01"),
			RequestId: ptrStrIS("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(teoClient, "DescribeInferenceServices", func(request *teov20220901.DescribeInferenceServicesRequest) (*teov20220901.DescribeInferenceServicesResponse, error) {
		resp := teov20220901.NewDescribeInferenceServicesResponse()
		resp.Response = &teov20220901.DescribeInferenceServicesResponseParams{
			Services: []*teov20220901.InferenceService{
				{
					ServiceId:  ptrStrIS("svc-affinity-01"),
					Name:       ptrStrIS("test-affinity-svc"),
					ListenPort: ptrInt64IS(8080),
					Status:     ptrStrIS("Running"),
				},
			},
			TotalCount: ptrInt64IS(1),
			RequestId:  ptrStrIS("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForInferenceService()
	res := teo.ResourceTencentCloudTeoInferenceServiceV1()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":     "zone-12345678",
		"name":        "test-affinity-svc",
		"listen_port": 8080,
		"containers": []interface{}{
			map[string]interface{}{
				"image_type": "TCR",
				"tcr_config": []interface{}{
					map[string]interface{}{
						"tcr_type":    "Personal",
						"image":       "ccr.ccs.tencentyun.com/test/image:v1",
						"region_name": "ap-guangzhou",
					},
				},
			},
		},
		"resource_config": []interface{}{
			map[string]interface{}{
				"scaling_mode": "Manual",
				"manual_instance_config": []interface{}{
					map[string]interface{}{
						"fixed_instance_count": 1,
					},
				},
			},
		},
		"affinity_config": []interface{}{
			map[string]interface{}{
				"switch":        "On",
				"affinity_mode": "SessionId",
				"source":        "Header",
				"header_name":   "X-Custom-Session-Id",
			},
		},
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-12345678#svc-affinity-01", d.Id())
}

// TestInferenceServiceV1_Create_APIError tests Create handles API error
func TestInferenceServiceV1_Create_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForInferenceService().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "CreateInferenceServiceWithContext", func(_ context.Context, request *teov20220901.CreateInferenceServiceRequest) (*teov20220901.CreateInferenceServiceResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=InvalidParameter, Message=Invalid zone_id")
	})

	meta := newMockMetaForInferenceService()
	res := teo.ResourceTencentCloudTeoInferenceServiceV1()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":     "zone-invalid",
		"name":        "test-inference-svc",
		"listen_port": 8080,
		"containers": []interface{}{
			map[string]interface{}{
				"image_type": "TCR",
				"tcr_config": []interface{}{
					map[string]interface{}{
						"tcr_type":    "Personal",
						"image":       "ccr.ccs.tencentyun.com/test/image:v1",
						"region_name": "ap-guangzhou",
					},
				},
			},
		},
		"resource_config": []interface{}{
			map[string]interface{}{
				"scaling_mode": "Manual",
				"manual_instance_config": []interface{}{
					map[string]interface{}{
						"fixed_instance_count": 1,
					},
				},
			},
		},
	})

	err := res.Create(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "InvalidParameter")
}

// TestInferenceServiceV1_Read_Success tests Read populates state from API
func TestInferenceServiceV1_Read_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForInferenceService().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeInferenceServices", func(request *teov20220901.DescribeInferenceServicesRequest) (*teov20220901.DescribeInferenceServicesResponse, error) {
		assert.Equal(t, "zone-12345678", *request.ZoneId)
		resp := teov20220901.NewDescribeInferenceServicesResponse()
		resp.Response = &teov20220901.DescribeInferenceServicesResponseParams{
			Services: []*teov20220901.InferenceService{
				{
					ServiceId:    ptrStrIS("svc-abcdefgh"),
					Name:         ptrStrIS("test-inference-svc"),
					ListenPort:   ptrInt64IS(8080),
					Status:       ptrStrIS("Running"),
					InferenceURL: ptrStrIS("https://svc-abcdefgh.edgeone.com"),
					CreateTime:   ptrStrIS("2024-01-01T00:00:00Z"),
					UpdateTime:   ptrStrIS("2024-01-01T01:00:00Z"),
				},
			},
			TotalCount: ptrInt64IS(1),
			RequestId:  ptrStrIS("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForInferenceService()
	res := teo.ResourceTencentCloudTeoInferenceServiceV1()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":     "zone-12345678",
		"name":        "test-inference-svc",
		"listen_port": 8080,
	})
	d.SetId("zone-12345678#svc-abcdefgh")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "svc-abcdefgh", d.Get("service_id"))
	assert.Equal(t, "test-inference-svc", d.Get("name"))
	assert.Equal(t, "Running", d.Get("status"))
	assert.Equal(t, "https://svc-abcdefgh.edgeone.com", d.Get("inference_url"))
}

// TestInferenceServiceV1_Read_NotFound tests Read handles resource not found
func TestInferenceServiceV1_Read_NotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForInferenceService().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeInferenceServices", func(request *teov20220901.DescribeInferenceServicesRequest) (*teov20220901.DescribeInferenceServicesResponse, error) {
		resp := teov20220901.NewDescribeInferenceServicesResponse()
		resp.Response = &teov20220901.DescribeInferenceServicesResponseParams{
			Services:   []*teov20220901.InferenceService{},
			TotalCount: ptrInt64IS(0),
			RequestId:  ptrStrIS("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForInferenceService()
	res := teo.ResourceTencentCloudTeoInferenceServiceV1()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":     "zone-12345678",
		"name":        "test-inference-svc",
		"listen_port": 8080,
	})
	d.SetId("zone-12345678#svc-notfound")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}

// TestInferenceServiceV1_Update_Success tests Update calls Modify API
func TestInferenceServiceV1_Update_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForInferenceService().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "ModifyInferenceServiceWithContext", func(_ context.Context, request *teov20220901.ModifyInferenceServiceRequest) (*teov20220901.ModifyInferenceServiceResponse, error) {
		assert.Equal(t, "zone-12345678", *request.ZoneId)
		assert.Equal(t, "svc-abcdefgh", *request.ServiceId)
		assert.Equal(t, int64(9090), *request.ListenPort)
		resp := teov20220901.NewModifyInferenceServiceResponse()
		resp.Response = &teov20220901.ModifyInferenceServiceResponseParams{
			RequestId: ptrStrIS("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(teoClient, "DescribeInferenceServices", func(request *teov20220901.DescribeInferenceServicesRequest) (*teov20220901.DescribeInferenceServicesResponse, error) {
		resp := teov20220901.NewDescribeInferenceServicesResponse()
		resp.Response = &teov20220901.DescribeInferenceServicesResponseParams{
			Services: []*teov20220901.InferenceService{
				{
					ServiceId:  ptrStrIS("svc-abcdefgh"),
					Name:       ptrStrIS("test-inference-svc"),
					ListenPort: ptrInt64IS(9090),
					Status:     ptrStrIS("Running"),
				},
			},
			TotalCount: ptrInt64IS(1),
			RequestId:  ptrStrIS("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForInferenceService()
	res := teo.ResourceTencentCloudTeoInferenceServiceV1()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":     "zone-12345678",
		"name":        "test-inference-svc",
		"listen_port": 9090,
		"containers": []interface{}{
			map[string]interface{}{
				"image_type": "TCR",
				"tcr_config": []interface{}{
					map[string]interface{}{
						"tcr_type":    "Personal",
						"image":       "ccr.ccs.tencentyun.com/test/image:v2",
						"region_name": "ap-guangzhou",
					},
				},
			},
		},
		"resource_config": []interface{}{
			map[string]interface{}{
				"scaling_mode": "Manual",
				"manual_instance_config": []interface{}{
					map[string]interface{}{
						"fixed_instance_count": 2,
					},
				},
			},
		},
	})
	d.SetId("zone-12345678#svc-abcdefgh")

	err := res.Update(d, meta)
	assert.NoError(t, err)
}

// TestInferenceServiceV1_Update_AffinityConfig tests Update with affinity config change
func TestInferenceServiceV1_Update_AffinityConfig(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForInferenceService().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "ModifyInferenceServiceWithContext", func(_ context.Context, request *teov20220901.ModifyInferenceServiceRequest) (*teov20220901.ModifyInferenceServiceResponse, error) {
		assert.NotNil(t, request.AffinityConfig)
		assert.Equal(t, "Off", *request.AffinityConfig.Switch)
		resp := teov20220901.NewModifyInferenceServiceResponse()
		resp.Response = &teov20220901.ModifyInferenceServiceResponseParams{
			RequestId: ptrStrIS("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(teoClient, "DescribeInferenceServices", func(request *teov20220901.DescribeInferenceServicesRequest) (*teov20220901.DescribeInferenceServicesResponse, error) {
		resp := teov20220901.NewDescribeInferenceServicesResponse()
		resp.Response = &teov20220901.DescribeInferenceServicesResponseParams{
			Services: []*teov20220901.InferenceService{
				{
					ServiceId:  ptrStrIS("svc-abcdefgh"),
					Name:       ptrStrIS("test-inference-svc"),
					ListenPort: ptrInt64IS(8080),
					Status:     ptrStrIS("Running"),
				},
			},
			TotalCount: ptrInt64IS(1),
			RequestId:  ptrStrIS("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForInferenceService()
	res := teo.ResourceTencentCloudTeoInferenceServiceV1()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":     "zone-12345678",
		"name":        "test-inference-svc",
		"listen_port": 8080,
		"containers": []interface{}{
			map[string]interface{}{
				"image_type": "TCR",
				"tcr_config": []interface{}{
					map[string]interface{}{
						"tcr_type":    "Personal",
						"image":       "ccr.ccs.tencentyun.com/test/image:v1",
						"region_name": "ap-guangzhou",
					},
				},
			},
		},
		"resource_config": []interface{}{
			map[string]interface{}{
				"scaling_mode": "Manual",
				"manual_instance_config": []interface{}{
					map[string]interface{}{
						"fixed_instance_count": 1,
					},
				},
			},
		},
		"affinity_config": []interface{}{
			map[string]interface{}{
				"switch": "Off",
			},
		},
	})
	d.SetId("zone-12345678#svc-abcdefgh")

	err := res.Update(d, meta)
	assert.NoError(t, err)
}

// TestInferenceServiceV1_Delete_Success tests Delete calls API successfully
func TestInferenceServiceV1_Delete_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForInferenceService().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "OperateInferenceServiceWithContext", func(_ context.Context, request *teov20220901.OperateInferenceServiceRequest) (*teov20220901.OperateInferenceServiceResponse, error) {
		assert.Equal(t, "zone-12345678", *request.ZoneId)
		assert.Equal(t, "svc-abcdefgh", *request.ServiceId)
		assert.Equal(t, "Delete", *request.Operation)
		resp := teov20220901.NewOperateInferenceServiceResponse()
		resp.Response = &teov20220901.OperateInferenceServiceResponseParams{
			RequestId: ptrStrIS("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForInferenceService()
	res := teo.ResourceTencentCloudTeoInferenceServiceV1()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":     "zone-12345678",
		"name":        "test-inference-svc",
		"listen_port": 8080,
	})
	d.SetId("zone-12345678#svc-abcdefgh")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
}

// TestInferenceServiceV1_Delete_APIError tests Delete handles API error
func TestInferenceServiceV1_Delete_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForInferenceService().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "OperateInferenceServiceWithContext", func(_ context.Context, request *teov20220901.OperateInferenceServiceRequest) (*teov20220901.OperateInferenceServiceResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=ResourceNotFound, Message=Service not found")
	})

	meta := newMockMetaForInferenceService()
	res := teo.ResourceTencentCloudTeoInferenceServiceV1()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":     "zone-12345678",
		"name":        "test-inference-svc",
		"listen_port": 8080,
	})
	d.SetId("zone-12345678#svc-abcdefgh")

	err := res.Delete(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ResourceNotFound")
}

// TestInferenceServiceV1_Schema validates schema definition
func TestInferenceServiceV1_Schema(t *testing.T) {
	res := teo.ResourceTencentCloudTeoInferenceServiceV1()

	assert.NotNil(t, res)
	assert.NotNil(t, res.Create)
	assert.NotNil(t, res.Read)
	assert.NotNil(t, res.Update)
	assert.NotNil(t, res.Delete)
	assert.NotNil(t, res.Importer)

	assert.Contains(t, res.Schema, "zone_id")
	assert.Contains(t, res.Schema, "name")
	assert.Contains(t, res.Schema, "listen_port")
	assert.Contains(t, res.Schema, "containers")
	assert.Contains(t, res.Schema, "resource_config")
	assert.Contains(t, res.Schema, "affinity_config")
	assert.Contains(t, res.Schema, "request_paths")
	assert.Contains(t, res.Schema, "description")
	assert.Contains(t, res.Schema, "service_id")
	assert.Contains(t, res.Schema, "status")
	assert.Contains(t, res.Schema, "inference_url")

	zoneId := res.Schema["zone_id"]
	assert.Equal(t, schema.TypeString, zoneId.Type)
	assert.True(t, zoneId.Required)
	assert.True(t, zoneId.ForceNew)

	name := res.Schema["name"]
	assert.Equal(t, schema.TypeString, name.Type)
	assert.True(t, name.Required)
	assert.True(t, name.ForceNew)

	listenPort := res.Schema["listen_port"]
	assert.Equal(t, schema.TypeInt, listenPort.Type)
	assert.True(t, listenPort.Required)

	affinityConfig := res.Schema["affinity_config"]
	assert.Equal(t, schema.TypeList, affinityConfig.Type)
	assert.True(t, affinityConfig.Optional)

	serviceId := res.Schema["service_id"]
	assert.Equal(t, schema.TypeString, serviceId.Type)
	assert.True(t, serviceId.Computed)
	assert.False(t, serviceId.Required)
	assert.False(t, serviceId.Optional)
}
