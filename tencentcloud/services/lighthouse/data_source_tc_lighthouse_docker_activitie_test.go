package lighthouse_test

import (
	"fmt"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	lighthouseSDK "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/lighthouse/v20200324"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/lighthouse"
)

// dockerActivitieMockMeta implements tccommon.ProviderMeta
type dockerActivitieMockMeta struct {
	client *connectivity.TencentCloudClient
}

func (m *dockerActivitieMockMeta) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &dockerActivitieMockMeta{}

func newDockerActivitieMockMeta() *dockerActivitieMockMeta {
	return &dockerActivitieMockMeta{client: &connectivity.TencentCloudClient{}}
}

func ptrString(v string) *string {
	return &v
}

func ptrInt64(v int64) *int64 {
	return &v
}

// go test ./tencentcloud/services/lighthouse/ -run "TestLighthouseDockerActivitie" -v -count=1 -gcflags="all=-l"

// TestLighthouseDockerActivitie_Read_Success tests Read retrieves Docker activities
func TestLighthouseDockerActivitie_Read_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	lighthouseClient := &lighthouseSDK.Client{}
	patches.ApplyMethodReturn(newDockerActivitieMockMeta().client, "UseLighthouseClient", lighthouseClient)

	patches.ApplyMethodFunc(lighthouseClient, "DescribeDockerActivities", func(request *lighthouseSDK.DescribeDockerActivitiesRequest) (*lighthouseSDK.DescribeDockerActivitiesResponse, error) {
		resp := lighthouseSDK.NewDescribeDockerActivitiesResponse()
		resp.Response = &lighthouseSDK.DescribeDockerActivitiesResponseParams{
			TotalCount: ptrInt64(1),
			DockerActivitySet: []*lighthouseSDK.DockerActivity{
				{
					ActivityId:            ptrString("lhda-12345678"),
					ActivityName:          ptrString("CreateContainer"),
					ActivityState:         ptrString("SUCCESS"),
					ActivityCommandOutput: ptrString("dGVzdCBvdXRwdXQ="),
					ContainerIds:          []*string{ptrString("lhcon-abcdef12")},
					CreatedTime:           ptrString("2024-06-01T10:00:00Z"),
					EndTime:               ptrString("2024-06-01T10:01:00Z"),
				},
			},
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newDockerActivitieMockMeta()
	res := lighthouse.DataSourceTencentCloudLighthouseDockerActivitie()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"instance_id": "lhins-12345678",
	})

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "lhins-12345678", d.Get("instance_id"))
}

// TestLighthouseDockerActivitie_Read_WithActivityIds tests Read with activity_ids filter
func TestLighthouseDockerActivitie_Read_WithActivityIds(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	lighthouseClient := &lighthouseSDK.Client{}
	patches.ApplyMethodReturn(newDockerActivitieMockMeta().client, "UseLighthouseClient", lighthouseClient)

	patches.ApplyMethodFunc(lighthouseClient, "DescribeDockerActivities", func(request *lighthouseSDK.DescribeDockerActivitiesRequest) (*lighthouseSDK.DescribeDockerActivitiesResponse, error) {
		resp := lighthouseSDK.NewDescribeDockerActivitiesResponse()
		resp.Response = &lighthouseSDK.DescribeDockerActivitiesResponseParams{
			TotalCount: ptrInt64(2),
			DockerActivitySet: []*lighthouseSDK.DockerActivity{
				{
					ActivityId:    ptrString("lhda-12345678"),
					ActivityName:  ptrString("CreateContainer"),
					ActivityState: ptrString("SUCCESS"),
					ContainerIds:  []*string{ptrString("lhcon-abcdef12")},
					CreatedTime:   ptrString("2024-06-01T10:00:00Z"),
					EndTime:       ptrString("2024-06-01T10:01:00Z"),
				},
				{
					ActivityId:    ptrString("lhda-87654321"),
					ActivityName:  ptrString("DeleteContainer"),
					ActivityState: ptrString("OPERATING"),
					ContainerIds:  []*string{ptrString("lhcon-fedcba21")},
					CreatedTime:   ptrString("2024-06-02T10:00:00Z"),
				},
			},
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newDockerActivitieMockMeta()
	res := lighthouse.DataSourceTencentCloudLighthouseDockerActivitie()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"instance_id":  "lhins-12345678",
		"activity_ids": []interface{}{"lhda-12345678", "lhda-87654321"},
	})

	err := res.Read(d, meta)
	assert.NoError(t, err)
}

// TestLighthouseDockerActivitie_Read_WithTimeRange tests Read with time range filter
func TestLighthouseDockerActivitie_Read_WithTimeRange(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	lighthouseClient := &lighthouseSDK.Client{}
	patches.ApplyMethodReturn(newDockerActivitieMockMeta().client, "UseLighthouseClient", lighthouseClient)

	patches.ApplyMethodFunc(lighthouseClient, "DescribeDockerActivities", func(request *lighthouseSDK.DescribeDockerActivitiesRequest) (*lighthouseSDK.DescribeDockerActivitiesResponse, error) {
		assert.NotNil(t, request.CreatedTimeBegin)
		assert.NotNil(t, request.CreatedTimeEnd)

		resp := lighthouseSDK.NewDescribeDockerActivitiesResponse()
		resp.Response = &lighthouseSDK.DescribeDockerActivitiesResponseParams{
			TotalCount:        ptrInt64(0),
			DockerActivitySet: []*lighthouseSDK.DockerActivity{},
			RequestId:         ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newDockerActivitieMockMeta()
	res := lighthouse.DataSourceTencentCloudLighthouseDockerActivitie()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"instance_id":        "lhins-12345678",
		"created_time_begin": 1717200000,
		"created_time_end":   1719800000,
	})

	err := res.Read(d, meta)
	assert.NoError(t, err)
}

// TestLighthouseDockerActivitie_Read_APIError tests Read handles API error
func TestLighthouseDockerActivitie_Read_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	lighthouseClient := &lighthouseSDK.Client{}
	patches.ApplyMethodReturn(newDockerActivitieMockMeta().client, "UseLighthouseClient", lighthouseClient)

	patches.ApplyMethodFunc(lighthouseClient, "DescribeDockerActivities", func(request *lighthouseSDK.DescribeDockerActivitiesRequest) (*lighthouseSDK.DescribeDockerActivitiesResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=InvalidParameter, Message=Invalid instance_id")
	})

	meta := newDockerActivitieMockMeta()
	res := lighthouse.DataSourceTencentCloudLighthouseDockerActivitie()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"instance_id": "lhins-invalid",
	})

	err := res.Read(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "InvalidParameter")
}

// TestLighthouseDockerActivitie_Read_NilFields tests Read handles nil fields gracefully
func TestLighthouseDockerActivitie_Read_NilFields(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	lighthouseClient := &lighthouseSDK.Client{}
	patches.ApplyMethodReturn(newDockerActivitieMockMeta().client, "UseLighthouseClient", lighthouseClient)

	patches.ApplyMethodFunc(lighthouseClient, "DescribeDockerActivities", func(request *lighthouseSDK.DescribeDockerActivitiesRequest) (*lighthouseSDK.DescribeDockerActivitiesResponse, error) {
		resp := lighthouseSDK.NewDescribeDockerActivitiesResponse()
		resp.Response = &lighthouseSDK.DescribeDockerActivitiesResponseParams{
			TotalCount: ptrInt64(1),
			DockerActivitySet: []*lighthouseSDK.DockerActivity{
				{
					ActivityId:    ptrString("lhda-12345678"),
					ActivityName:  ptrString("CreateContainer"),
					ActivityState: ptrString("INIT"),
					CreatedTime:   ptrString("2024-06-01T10:00:00Z"),
				},
			},
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newDockerActivitieMockMeta()
	res := lighthouse.DataSourceTencentCloudLighthouseDockerActivitie()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"instance_id": "lhins-12345678",
	})

	err := res.Read(d, meta)
	assert.NoError(t, err)
}

// TestLighthouseDockerActivitie_Schema validates schema definition
func TestLighthouseDockerActivitie_Schema(t *testing.T) {
	res := lighthouse.DataSourceTencentCloudLighthouseDockerActivitie()

	assert.NotNil(t, res)
	assert.NotNil(t, res.Read)

	// Check optional fields
	assert.Contains(t, res.Schema, "instance_id")
	instanceId := res.Schema["instance_id"]
	assert.Equal(t, schema.TypeString, instanceId.Type)
	assert.True(t, instanceId.Optional)

	assert.Contains(t, res.Schema, "activity_ids")
	activityIds := res.Schema["activity_ids"]
	assert.Equal(t, schema.TypeSet, activityIds.Type)
	assert.True(t, activityIds.Optional)

	assert.Contains(t, res.Schema, "created_time_begin")
	createdTimeBegin := res.Schema["created_time_begin"]
	assert.Equal(t, schema.TypeInt, createdTimeBegin.Type)
	assert.True(t, createdTimeBegin.Optional)

	assert.Contains(t, res.Schema, "created_time_end")
	createdTimeEnd := res.Schema["created_time_end"]
	assert.Equal(t, schema.TypeInt, createdTimeEnd.Type)
	assert.True(t, createdTimeEnd.Optional)

	// Check computed fields
	assert.Contains(t, res.Schema, "docker_activity_set")
	dockerActivitySet := res.Schema["docker_activity_set"]
	assert.Equal(t, schema.TypeList, dockerActivitySet.Type)
	assert.True(t, dockerActivitySet.Computed)

	assert.Contains(t, res.Schema, "result_output_file")
	resultOutputFile := res.Schema["result_output_file"]
	assert.Equal(t, schema.TypeString, resultOutputFile.Type)
	assert.True(t, resultOutputFile.Optional)
}
