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

// go test ./tencentcloud/services/teo/ -run "TestInferenceServiceV1" -v -count=1 -gcflags="all=-l"

// TestInferenceServiceV1_Create_Success tests Create calls API and sets ID
func TestInferenceServiceV1_Create_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "CreateInferenceService", func(request *teov20220901.CreateInferenceServiceRequest) (*teov20220901.CreateInferenceServiceResponse, error) {
		assert.Equal(t, "zone-1234567890", *request.ZoneId)
		assert.Equal(t, "test-inference-service", *request.Name)
		resp := teov20220901.NewCreateInferenceServiceResponse()
		resp.Response = &teov20220901.CreateInferenceServiceResponseParams{
			ServiceId: ptrString("service-abcdefghij"),
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(teoClient, "DescribeInferenceServices", func(request *teov20220901.DescribeInferenceServicesRequest) (*teov20220901.DescribeInferenceServicesResponse, error) {
		resp := teov20220901.NewDescribeInferenceServicesResponse()
		resp.Response = &teov20220901.DescribeInferenceServicesResponseParams{
			TotalCount: ptrInt64(1),
			Services: []*teov20220901.InferenceService{
				{
					ServiceId:    ptrString("service-abcdefghij"),
					Name:         ptrString("test-inference-service"),
					ListenPort:   ptrInt64(8080),
					Description:  ptrString("test description"),
					Status:       ptrString("Running"),
					InferenceURL: ptrString("https://inference.example.com"),
					CreateTime:   ptrString("2024-01-01T00:00:00Z"),
					UpdateTime:   ptrString("2024-01-01T00:00:00Z"),
					Containers: []*teov20220901.InferenceContainerConfig{
						{
							ImageType:      ptrString("TCR"),
							StartupCommand: ptrString("python app.py"),
							TcrRepositoryConfig: &teov20220901.InferenceTCRRepositoryConfig{
								TCRType:    ptrString("Personal"),
								Image:      ptrString("ccr.ccs.tencentyun.com/ns/img:v1"),
								RegionName: ptrString("ap-guangzhou"),
							},
						},
					},
					ResourceConfig: &teov20220901.InferenceResourceConfig{
						ScalingMode:  ptrString("Auto"),
						HardwareSpec: ptrString("GPU.A10.1"),
						Concurrency:  ptrInt64(10),
					},
					RequestPaths: []*string{ptrString("/predict")},
				},
			},
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoInferenceServiceV1()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":     "zone-1234567890",
		"name":        "test-inference-service",
		"listen_port": 8080,
		"description": "test description",
		"containers": []interface{}{
			map[string]interface{}{
				"image_type":      "TCR",
				"startup_command": "python app.py",
				"tcr_repository_config": []interface{}{
					map[string]interface{}{
						"tcr_type":    "Personal",
						"image":       "ccr.ccs.tencentyun.com/ns/img:v1",
						"region_name": "ap-guangzhou",
					},
				},
			},
		},
		"resource_config": []interface{}{
			map[string]interface{}{
				"scaling_mode":  "Auto",
				"hardware_spec": "GPU.A10.1",
				"concurrency":   10,
			},
		},
		"request_paths": []interface{}{"/predict"},
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-1234567890#service-abcdefghij", d.Id())
}

// TestInferenceServiceV1_Create_APIError tests Create handles API error
func TestInferenceServiceV1_Create_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "CreateInferenceService", func(request *teov20220901.CreateInferenceServiceRequest) (*teov20220901.CreateInferenceServiceResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=InvalidParameter, Message=Invalid zone_id")
	})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoInferenceServiceV1()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":     "zone-invalid",
		"name":        "test-inference-service",
		"listen_port": 8080,
		"containers": []interface{}{
			map[string]interface{}{
				"image_type": "TCR",
				"tcr_repository_config": []interface{}{
					map[string]interface{}{
						"tcr_type":    "Personal",
						"image":       "ccr.ccs.tencentyun.com/ns/img:v1",
						"region_name": "ap-guangzhou",
					},
				},
			},
		},
		"resource_config": []interface{}{
			map[string]interface{}{
				"scaling_mode":  "Auto",
				"hardware_spec": "GPU.A10.1",
			},
		},
	})

	err := res.Create(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "InvalidParameter")
}

// TestInferenceServiceV1_Create_NilResponse tests Create handles nil response
func TestInferenceServiceV1_Create_NilResponse(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "CreateInferenceService", func(request *teov20220901.CreateInferenceServiceRequest) (*teov20220901.CreateInferenceServiceResponse, error) {
		return nil, nil
	})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoInferenceServiceV1()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":     "zone-1234567890",
		"name":        "test-inference-service",
		"listen_port": 8080,
		"containers": []interface{}{
			map[string]interface{}{
				"image_type": "TCR",
				"tcr_repository_config": []interface{}{
					map[string]interface{}{
						"tcr_type":    "Personal",
						"image":       "ccr.ccs.tencentyun.com/ns/img:v1",
						"region_name": "ap-guangzhou",
					},
				},
			},
		},
		"resource_config": []interface{}{
			map[string]interface{}{
				"scaling_mode":  "Auto",
				"hardware_spec": "GPU.A10.1",
			},
		},
	})

	err := res.Create(d, meta)
	assert.Error(t, err)
}

// TestInferenceServiceV1_Read_Success tests Read retrieves inference service data
func TestInferenceServiceV1_Read_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeInferenceServices", func(request *teov20220901.DescribeInferenceServicesRequest) (*teov20220901.DescribeInferenceServicesResponse, error) {
		resp := teov20220901.NewDescribeInferenceServicesResponse()
		resp.Response = &teov20220901.DescribeInferenceServicesResponseParams{
			TotalCount: ptrInt64(1),
			Services: []*teov20220901.InferenceService{
				{
					ServiceId:    ptrString("service-abcdefghij"),
					Name:         ptrString("test-inference-service"),
					ListenPort:   ptrInt64(8080),
					Description:  ptrString("test description"),
					Status:       ptrString("Running"),
					InferenceURL: ptrString("https://inference.example.com"),
					CreateTime:   ptrString("2024-01-01T00:00:00Z"),
					UpdateTime:   ptrString("2024-01-01T00:00:00Z"),
					Containers: []*teov20220901.InferenceContainerConfig{
						{
							ImageType:      ptrString("TCR"),
							StartupCommand: ptrString("python app.py"),
							TcrRepositoryConfig: &teov20220901.InferenceTCRRepositoryConfig{
								TCRType:    ptrString("Personal"),
								Image:      ptrString("ccr.ccs.tencentyun.com/ns/img:v1"),
								RegionName: ptrString("ap-guangzhou"),
							},
						},
					},
					ResourceConfig: &teov20220901.InferenceResourceConfig{
						ScalingMode:  ptrString("Auto"),
						HardwareSpec: ptrString("GPU.A10.1"),
						Concurrency:  ptrInt64(10),
					},
					RequestPaths: []*string{ptrString("/predict")},
				},
			},
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoInferenceServiceV1()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":     "zone-1234567890",
		"name":        "test-inference-service",
		"listen_port": 8080,
		"description": "test description",
		"containers": []interface{}{
			map[string]interface{}{
				"image_type": "TCR",
			},
		},
		"resource_config": []interface{}{
			map[string]interface{}{
				"scaling_mode":  "Auto",
				"hardware_spec": "GPU.A10.1",
			},
		},
	})
	d.SetId("zone-1234567890#service-abcdefghij")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "test-inference-service", d.Get("name"))
	assert.Equal(t, "Running", d.Get("status"))
}

// TestInferenceServiceV1_Read_NotFound tests Read handles service not found
func TestInferenceServiceV1_Read_NotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeInferenceServices", func(request *teov20220901.DescribeInferenceServicesRequest) (*teov20220901.DescribeInferenceServicesResponse, error) {
		resp := teov20220901.NewDescribeInferenceServicesResponse()
		resp.Response = &teov20220901.DescribeInferenceServicesResponseParams{
			TotalCount: ptrInt64(0),
			Services:   []*teov20220901.InferenceService{},
			RequestId:  ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoInferenceServiceV1()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-1234567890",
		"name":    "test-inference-service",
	})
	d.SetId("zone-1234567890#service-abcdefghij")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}

// TestInferenceServiceV1_Update_Success tests Update calls Modify API
func TestInferenceServiceV1_Update_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)

	modifyCalled := false
	patches.ApplyMethodFunc(teoClient, "ModifyInferenceService", func(request *teov20220901.ModifyInferenceServiceRequest) (*teov20220901.ModifyInferenceServiceResponse, error) {
		modifyCalled = true
		resp := teov20220901.NewModifyInferenceServiceResponse()
		resp.Response = &teov20220901.ModifyInferenceServiceResponseParams{
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(teoClient, "DescribeInferenceServices", func(request *teov20220901.DescribeInferenceServicesRequest) (*teov20220901.DescribeInferenceServicesResponse, error) {
		resp := teov20220901.NewDescribeInferenceServicesResponse()
		resp.Response = &teov20220901.DescribeInferenceServicesResponseParams{
			TotalCount: ptrInt64(1),
			Services: []*teov20220901.InferenceService{
				{
					ServiceId:    ptrString("service-abcdefghij"),
					Name:         ptrString("test-inference-service"),
					ListenPort:   ptrInt64(9090),
					Description:  ptrString("updated description"),
					Status:       ptrString("Running"),
					InferenceURL: ptrString("https://inference.example.com"),
					CreateTime:   ptrString("2024-01-01T00:00:00Z"),
					UpdateTime:   ptrString("2024-01-02T00:00:00Z"),
				},
			},
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoInferenceServiceV1()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":     "zone-1234567890",
		"name":        "test-inference-service",
		"listen_port": 9090,
		"description": "updated description",
		"containers": []interface{}{
			map[string]interface{}{
				"image_type": "TCR",
				"tcr_repository_config": []interface{}{
					map[string]interface{}{
						"tcr_type":    "Personal",
						"image":       "ccr.ccs.tencentyun.com/ns/img:v2",
						"region_name": "ap-guangzhou",
					},
				},
			},
		},
		"resource_config": []interface{}{
			map[string]interface{}{
				"scaling_mode":  "Auto",
				"hardware_spec": "GPU.A10.1",
			},
		},
	})
	d.SetId("zone-1234567890#service-abcdefghij")

	err := res.Update(d, meta)
	assert.NoError(t, err)
	assert.True(t, modifyCalled)
}

// TestInferenceServiceV1_Update_Operation tests Update handles Stop/Resume operation
func TestInferenceServiceV1_Update_Operation(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)

	operateCalled := false
	patches.ApplyMethodFunc(teoClient, "OperateInferenceService", func(request *teov20220901.OperateInferenceServiceRequest) (*teov20220901.OperateInferenceServiceResponse, error) {
		operateCalled = true
		assert.Equal(t, "Stop", *request.Operation)
		resp := teov20220901.NewOperateInferenceServiceResponse()
		resp.Response = &teov20220901.OperateInferenceServiceResponseParams{
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(teoClient, "DescribeInferenceServices", func(request *teov20220901.DescribeInferenceServicesRequest) (*teov20220901.DescribeInferenceServicesResponse, error) {
		resp := teov20220901.NewDescribeInferenceServicesResponse()
		resp.Response = &teov20220901.DescribeInferenceServicesResponseParams{
			TotalCount: ptrInt64(1),
			Services: []*teov20220901.InferenceService{
				{
					ServiceId: ptrString("service-abcdefghij"),
					Name:      ptrString("test-inference-service"),
					Status:    ptrString("Stopped"),
				},
			},
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoInferenceServiceV1()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":     "zone-1234567890",
		"name":        "test-inference-service",
		"listen_port": 8080,
		"operation":   "Stop",
		"containers": []interface{}{
			map[string]interface{}{
				"image_type": "TCR",
				"tcr_repository_config": []interface{}{
					map[string]interface{}{
						"tcr_type":    "Personal",
						"image":       "ccr.ccs.tencentyun.com/ns/img:v1",
						"region_name": "ap-guangzhou",
					},
				},
			},
		},
		"resource_config": []interface{}{
			map[string]interface{}{
				"scaling_mode":  "Auto",
				"hardware_spec": "GPU.A10.1",
			},
		},
	})
	d.SetId("zone-1234567890#service-abcdefghij")

	err := res.Update(d, meta)
	assert.NoError(t, err)
	assert.True(t, operateCalled)
}

// TestInferenceServiceV1_Delete_Success tests Delete calls OperateInferenceService
func TestInferenceServiceV1_Delete_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "OperateInferenceService", func(request *teov20220901.OperateInferenceServiceRequest) (*teov20220901.OperateInferenceServiceResponse, error) {
		assert.Equal(t, "Delete", *request.Operation)
		assert.Equal(t, "service-abcdefghij", *request.ServiceId)
		resp := teov20220901.NewOperateInferenceServiceResponse()
		resp.Response = &teov20220901.OperateInferenceServiceResponseParams{
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoInferenceServiceV1()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-1234567890",
		"name":    "test-inference-service",
	})
	d.SetId("zone-1234567890#service-abcdefghij")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
}

// TestInferenceServiceV1_Delete_APIError tests Delete handles API error
func TestInferenceServiceV1_Delete_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "OperateInferenceService", func(request *teov20220901.OperateInferenceServiceRequest) (*teov20220901.OperateInferenceServiceResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=ResourceNotFound.InferenceService, Message=Service not found")
	})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoInferenceServiceV1()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-1234567890",
		"name":    "test-inference-service",
	})
	d.SetId("zone-1234567890#service-abcdefghij")

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

	// Check required fields
	assert.Contains(t, res.Schema, "zone_id")
	zoneId := res.Schema["zone_id"]
	assert.Equal(t, schema.TypeString, zoneId.Type)
	assert.True(t, zoneId.Required)
	assert.True(t, zoneId.ForceNew)

	assert.Contains(t, res.Schema, "name")
	name := res.Schema["name"]
	assert.Equal(t, schema.TypeString, name.Type)
	assert.True(t, name.Required)
	assert.True(t, name.ForceNew)

	assert.Contains(t, res.Schema, "listen_port")
	listenPort := res.Schema["listen_port"]
	assert.Equal(t, schema.TypeInt, listenPort.Type)
	assert.True(t, listenPort.Required)

	assert.Contains(t, res.Schema, "containers")
	containers := res.Schema["containers"]
	assert.Equal(t, schema.TypeList, containers.Type)
	assert.True(t, containers.Required)

	assert.Contains(t, res.Schema, "resource_config")
	resourceConfig := res.Schema["resource_config"]
	assert.Equal(t, schema.TypeList, resourceConfig.Type)
	assert.True(t, resourceConfig.Required)

	// Check optional fields
	assert.Contains(t, res.Schema, "request_paths")
	requestPaths := res.Schema["request_paths"]
	assert.Equal(t, schema.TypeSet, requestPaths.Type)
	assert.True(t, requestPaths.Optional)

	assert.Contains(t, res.Schema, "description")
	description := res.Schema["description"]
	assert.Equal(t, schema.TypeString, description.Type)
	assert.True(t, description.Optional)

	assert.Contains(t, res.Schema, "operation")
	operation := res.Schema["operation"]
	assert.Equal(t, schema.TypeString, operation.Type)
	assert.True(t, operation.Optional)

	// Check computed fields
	assert.Contains(t, res.Schema, "service_id")
	serviceId := res.Schema["service_id"]
	assert.Equal(t, schema.TypeString, serviceId.Type)
	assert.True(t, serviceId.Computed)

	assert.Contains(t, res.Schema, "status")
	status := res.Schema["status"]
	assert.Equal(t, schema.TypeString, status.Type)
	assert.True(t, status.Computed)

	assert.Contains(t, res.Schema, "inference_url")
	assert.Contains(t, res.Schema, "create_time")
	assert.Contains(t, res.Schema, "update_time")

	// Check Timeouts
	assert.NotNil(t, res.Timeouts)
	assert.NotNil(t, res.Timeouts.Create)
	assert.NotNil(t, res.Timeouts.Update)
	assert.NotNil(t, res.Timeouts.Delete)
}
