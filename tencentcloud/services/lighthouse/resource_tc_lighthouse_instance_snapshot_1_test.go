package lighthouse_test

import (
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	lighthouse "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/lighthouse/v20200324"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	svclighthouse "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/lighthouse"
)

// mockMetaForSnapshot implements tccommon.ProviderMeta
type mockMetaForSnapshot struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaForSnapshot) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaForSnapshot{}

func newMockMetaForSnapshot() *mockMetaForSnapshot {
	return &mockMetaForSnapshot{client: &connectivity.TencentCloudClient{}}
}

// go test ./tencentcloud/services/lighthouse/ -run "TestLighthouseInstanceSnapshot1" -v -count=1 -gcflags="all=-l"

func TestLighthouseInstanceSnapshot1_Create(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	lighthouseClient := &lighthouse.Client{}
	patches.ApplyMethodReturn(newMockMetaForSnapshot().client, "UseLighthouseClient", lighthouseClient)

	patches.ApplyMethodFunc(lighthouseClient, "CreateInstanceSnapshot", func(request *lighthouse.CreateInstanceSnapshotRequest) (*lighthouse.CreateInstanceSnapshotResponse, error) {
		resp := lighthouse.NewCreateInstanceSnapshotResponse()
		resp.Response = &lighthouse.CreateInstanceSnapshotResponseParams{
			SnapshotId: ptrStr("lhsnap-test123"),
			RequestId:  ptrStr("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(lighthouseClient, "DescribeSnapshots", func(request *lighthouse.DescribeSnapshotsRequest) (*lighthouse.DescribeSnapshotsResponse, error) {
		resp := lighthouse.NewDescribeSnapshotsResponse()
		resp.Response = &lighthouse.DescribeSnapshotsResponseParams{
			SnapshotSet: []*lighthouse.Snapshot{
				{
					SnapshotId:             ptrStr("lhsnap-test123"),
					SnapshotName:           ptrStr("test-snapshot"),
					DiskUsage:              ptrStr("SYSTEM_DISK"),
					DiskId:                 ptrStr("lhdisk-test"),
					DiskSize:               ptrInt64(60),
					SnapshotState:          ptrStr("NORMAL"),
					Percent:                ptrInt64(100),
					LatestOperation:        ptrStr("CreateInstanceSnapshot"),
					LatestOperationState:   ptrStr("SUCCESS"),
					LatestOperationRequestId: ptrStr("op-req-id"),
					CreatedTime:            ptrStr("2024-01-01T00:00:00Z"),
				},
			},
			TotalCount: ptrInt64(1),
			RequestId:  ptrStr("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForSnapshot()
	res := svclighthouse.ResourceTencentCloudLighthouseInstanceSnapshot1()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"instance_id":   "lhins-test123",
		"snapshot_name": "test-snapshot",
		"tags": []interface{}{
			map[string]interface{}{
				"key":   "env",
				"value": "test",
			},
		},
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "lhsnap-test123", d.Id())
}

func TestLighthouseInstanceSnapshot1_Create_NilSnapshotId(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	lighthouseClient := &lighthouse.Client{}
	patches.ApplyMethodReturn(newMockMetaForSnapshot().client, "UseLighthouseClient", lighthouseClient)

	patches.ApplyMethodFunc(lighthouseClient, "CreateInstanceSnapshot", func(request *lighthouse.CreateInstanceSnapshotRequest) (*lighthouse.CreateInstanceSnapshotResponse, error) {
		resp := lighthouse.NewCreateInstanceSnapshotResponse()
		resp.Response = &lighthouse.CreateInstanceSnapshotResponseParams{
			SnapshotId: nil,
			RequestId:  ptrStr("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForSnapshot()
	res := svclighthouse.ResourceTencentCloudLighthouseInstanceSnapshot1()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"instance_id": "lhins-test123",
	})

	err := res.Create(d, meta)
	assert.Error(t, err)
}

func TestLighthouseInstanceSnapshot1_Read(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	lighthouseClient := &lighthouse.Client{}
	patches.ApplyMethodReturn(newMockMetaForSnapshot().client, "UseLighthouseClient", lighthouseClient)

	patches.ApplyMethodFunc(lighthouseClient, "DescribeSnapshots", func(request *lighthouse.DescribeSnapshotsRequest) (*lighthouse.DescribeSnapshotsResponse, error) {
		resp := lighthouse.NewDescribeSnapshotsResponse()
		resp.Response = &lighthouse.DescribeSnapshotsResponseParams{
			SnapshotSet: []*lighthouse.Snapshot{
				{
					SnapshotId:               ptrStr("lhsnap-test123"),
					SnapshotName:             ptrStr("test-snapshot"),
					DiskUsage:                ptrStr("SYSTEM_DISK"),
					DiskId:                   ptrStr("lhdisk-test"),
					DiskSize:                 ptrInt64(60),
					SnapshotState:            ptrStr("NORMAL"),
					Percent:                  ptrInt64(100),
					LatestOperation:          ptrStr("CreateInstanceSnapshot"),
					LatestOperationState:     ptrStr("SUCCESS"),
					LatestOperationRequestId: ptrStr("op-req-id"),
					CreatedTime:              ptrStr("2024-01-01T00:00:00Z"),
					Tags: []*lighthouse.Tag{
						{
							Key:   ptrStr("env"),
							Value: ptrStr("test"),
						},
					},
				},
			},
			TotalCount: ptrInt64(1),
			RequestId:  ptrStr("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForSnapshot()
	res := svclighthouse.ResourceTencentCloudLighthouseInstanceSnapshot1()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"instance_id":   "lhins-test123",
		"snapshot_name": "test-snapshot",
	})
	d.SetId("lhsnap-test123")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "lhsnap-test123", d.Id())
	assert.Equal(t, "lhsnap-test123", d.Get("snapshot_id"))
	assert.Equal(t, "test-snapshot", d.Get("snapshot_name"))
	assert.Equal(t, "SYSTEM_DISK", d.Get("disk_usage"))
	assert.Equal(t, "lhdisk-test", d.Get("disk_id"))
	assert.Equal(t, 60, d.Get("disk_size"))
	assert.Equal(t, "NORMAL", d.Get("snapshot_state"))
	assert.Equal(t, 100, d.Get("percent"))
	assert.Equal(t, "CreateInstanceSnapshot", d.Get("latest_operation"))
	assert.Equal(t, "SUCCESS", d.Get("latest_operation_state"))
	assert.Equal(t, "op-req-id", d.Get("latest_operation_request_id"))
	assert.Equal(t, "2024-01-01T00:00:00Z", d.Get("created_time"))
}

func TestLighthouseInstanceSnapshot1_Read_NotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	lighthouseClient := &lighthouse.Client{}
	patches.ApplyMethodReturn(newMockMetaForSnapshot().client, "UseLighthouseClient", lighthouseClient)

	patches.ApplyMethodFunc(lighthouseClient, "DescribeSnapshots", func(request *lighthouse.DescribeSnapshotsRequest) (*lighthouse.DescribeSnapshotsResponse, error) {
		resp := lighthouse.NewDescribeSnapshotsResponse()
		resp.Response = &lighthouse.DescribeSnapshotsResponseParams{
			SnapshotSet: []*lighthouse.Snapshot{},
			TotalCount:  ptrInt64(0),
			RequestId:   ptrStr("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForSnapshot()
	res := svclighthouse.ResourceTencentCloudLighthouseInstanceSnapshot1()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"instance_id": "lhins-test123",
	})
	d.SetId("lhsnap-test123")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}

func TestLighthouseInstanceSnapshot1_Update(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	lighthouseClient := &lighthouse.Client{}
	patches.ApplyMethodReturn(newMockMetaForSnapshot().client, "UseLighthouseClient", lighthouseClient)

	patches.ApplyMethodFunc(lighthouseClient, "ModifySnapshotAttribute", func(request *lighthouse.ModifySnapshotAttributeRequest) (*lighthouse.ModifySnapshotAttributeResponse, error) {
		resp := lighthouse.NewModifySnapshotAttributeResponse()
		resp.Response = &lighthouse.ModifySnapshotAttributeResponseParams{
			RequestId: ptrStr("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(lighthouseClient, "DescribeSnapshots", func(request *lighthouse.DescribeSnapshotsRequest) (*lighthouse.DescribeSnapshotsResponse, error) {
		resp := lighthouse.NewDescribeSnapshotsResponse()
		resp.Response = &lighthouse.DescribeSnapshotsResponseParams{
			SnapshotSet: []*lighthouse.Snapshot{
				{
					SnapshotId:   ptrStr("lhsnap-test123"),
					SnapshotName: ptrStr("updated-snapshot"),
					SnapshotState: ptrStr("NORMAL"),
				},
			},
			TotalCount: ptrInt64(1),
			RequestId:  ptrStr("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForSnapshot()
	res := svclighthouse.ResourceTencentCloudLighthouseInstanceSnapshot1()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"instance_id":   "lhins-test123",
		"snapshot_name": "updated-snapshot",
	})
	d.SetId("lhsnap-test123")

	err := res.Update(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "lhsnap-test123", d.Id())
}

func TestLighthouseInstanceSnapshot1_Delete(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	lighthouseClient := &lighthouse.Client{}
	patches.ApplyMethodReturn(newMockMetaForSnapshot().client, "UseLighthouseClient", lighthouseClient)

	patches.ApplyMethodFunc(lighthouseClient, "DeleteSnapshots", func(request *lighthouse.DeleteSnapshotsRequest) (*lighthouse.DeleteSnapshotsResponse, error) {
		resp := lighthouse.NewDeleteSnapshotsResponse()
		resp.Response = &lighthouse.DeleteSnapshotsResponseParams{
			RequestId: ptrStr("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForSnapshot()
	res := svclighthouse.ResourceTencentCloudLighthouseInstanceSnapshot1()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"instance_id": "lhins-test123",
	})
	d.SetId("lhsnap-test123")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
}