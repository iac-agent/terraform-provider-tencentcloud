package tcaplusdb

import (
	"context"
	"errors"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tcaplusdbv20190823 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tcaplusdb/v20190823"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

type tcaplusTableGroupTestMeta struct{}

func (m *tcaplusTableGroupTestMeta) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return nil
}

var _ tccommon.ProviderMeta = &tcaplusTableGroupTestMeta{}

func newTcaplusTableGroupTestClient(t *testing.T) (*connectivity.TencentCloudClient, *tcaplusdbv20190823.Client, *gomonkey.Patches) {
	t.Helper()

	patches := gomonkey.NewPatches()
	client := &connectivity.TencentCloudClient{}
	sdkClient := &tcaplusdbv20190823.Client{}
	patches.ApplyMethodReturn(client, "UseTcaplusClient", sdkClient)

	return client, sdkClient, patches
}

func TestParseTcaplusTableGroupId(t *testing.T) {
	clusterId, tableGroupId, err := parseTcaplusTableGroupId("cluster-abcd" + tccommon.FILED_SP + "109")
	if err != nil {
		t.Fatalf("expect no error, got %v", err)
	}
	if clusterId != "cluster-abcd" {
		t.Fatalf("expect cluster id cluster-abcd, got %s", clusterId)
	}
	if tableGroupId != "109" {
		t.Fatalf("expect table group id 109, got %s", tableGroupId)
	}

	if _, _, err = parseTcaplusTableGroupId("cluster-abcd"); err == nil {
		t.Fatal("expect error for broken id, got nil")
	}
}

func TestResourceTencentCloudTcaplusTableGroupSchema(t *testing.T) {
	r := ResourceTencentCloudTcaplusTableGroup()
	if r.Create == nil || r.Read == nil || r.Update == nil || r.Delete == nil {
		t.Fatal("expect CRUD handlers to be set")
	}
	if r.Importer == nil || r.Importer.State == nil {
		t.Fatal("expect importer to be set")
	}

	for _, name := range []string{"cluster_id", "table_group_name", "table_group_id", "resource_tags"} {
		if _, ok := r.Schema[name]; !ok {
			t.Fatalf("expect schema field %s", name)
		}
	}

	if !r.Schema["cluster_id"].ForceNew {
		t.Fatal("expect cluster_id to be ForceNew")
	}
	if !r.Schema["table_group_id"].ForceNew {
		t.Fatal("expect table_group_id to be ForceNew")
	}
	if !r.Schema["table_group_id"].Optional || !r.Schema["table_group_id"].Computed {
		t.Fatal("expect table_group_id to be Optional and Computed")
	}
	if r.Schema["table_group_name"].ForceNew {
		t.Fatal("expect table_group_name to be mutable")
	}
}

func TestResourceTencentCloudTcaplusTableGroupCreateSuccess(t *testing.T) {
	client, sdkClient, patches := newTcaplusTableGroupTestClient(t)
	defer patches.Reset()

	var gotRequest *tcaplusdbv20190823.CreateTableGroupRequest
	patches.ApplyMethodFunc(sdkClient, "CreateTableGroupWithContext", func(_ context.Context, request *tcaplusdbv20190823.CreateTableGroupRequest) (*tcaplusdbv20190823.CreateTableGroupResponse, error) {
		gotRequest = request
		response := tcaplusdbv20190823.NewCreateTableGroupResponse()
		response.Response = &tcaplusdbv20190823.CreateTableGroupResponseParams{
			TableGroupId: helper.String("109"),
		}
		return response, nil
	})

	patches.ApplyMethodFunc(client, "UseTcaplusClient", func() *tcaplusdbv20190823.Client {
		return sdkClient
	})

	describeCalls := 0
	patches.ApplyMethodFunc(sdkClient, "DescribeTableGroupsWithContext", func(_ context.Context, request *tcaplusdbv20190823.DescribeTableGroupsRequest) (*tcaplusdbv20190823.DescribeTableGroupsResponse, error) {
		describeCalls++
		response := tcaplusdbv20190823.NewDescribeTableGroupsResponse()
		response.Response = &tcaplusdbv20190823.DescribeTableGroupsResponseParams{
			TableGroups: []*tcaplusdbv20190823.TableGroupInfo{
				{
					TableGroupId:       helper.String("109"),
					TableGroupName:     helper.String("tf-example"),
					CreatedTime:        helper.String("2026-01-01 00:00:00"),
					TableCount:         helper.Uint64(0),
					TotalSize:          helper.Uint64(0),
					TxhBackupExpireDay: helper.Uint64(7),
					EnableMysql:        helper.Uint64(0),
				},
			},
		}
		return response, nil
	})

	resourceData := schema.TestResourceDataRaw(t, ResourceTencentCloudTcaplusTableGroup().Schema, map[string]interface{}{
		"cluster_id":       "cluster-abcd",
		"table_group_name": "tf-example",
	})

	if err := resourceTencentCloudTcaplusTableGroupCreate(resourceData, &tcaplusTableGroupTestMeta{}); err != nil {
		t.Fatalf("expect no error, got %v", err)
	}

	if gotRequest == nil || gotRequest.ClusterId == nil || *gotRequest.ClusterId != "cluster-abcd" {
		t.Fatal("expect CreateTableGroup request to carry cluster id cluster-abcd")
	}
	if gotRequest.TableGroupName == nil || *gotRequest.TableGroupName != "tf-example" {
		t.Fatal("expect CreateTableGroup request to carry table group name tf-example")
	}
	if gotRequest.TableGroupId != nil {
		t.Fatal("expect CreateTableGroup request to leave table group id unset")
	}

	expectId := "cluster-abcd" + tccommon.FILED_SP + "109"
	if resourceData.Id() != expectId {
		t.Fatalf("expect id %s, got %s", expectId, resourceData.Id())
	}

	if describeCalls == 0 {
		t.Fatal("expect DescribeTableGroups to be called after create")
	}

	if v := resourceData.Get("table_group_name").(string); v != "tf-example" {
		t.Fatalf("expect table_group_name tf-example, got %s", v)
	}
}

func TestResourceTencentCloudTcaplusTableGroupCreateEmptyResponse(t *testing.T) {
	_, sdkClient, patches := newTcaplusTableGroupTestClient(t)
	defer patches.Reset()

	patches.ApplyMethodFunc(sdkClient, "CreateTableGroupWithContext", func(_ context.Context, request *tcaplusdbv20190823.CreateTableGroupRequest) (*tcaplusdbv20190823.CreateTableGroupResponse, error) {
		response := tcaplusdbv20190823.NewCreateTableGroupResponse()
		response.Response = &tcaplusdbv20190823.CreateTableGroupResponseParams{}
		return response, nil
	})

	resourceData := schema.TestResourceDataRaw(t, ResourceTencentCloudTcaplusTableGroup().Schema, map[string]interface{}{
		"cluster_id":       "cluster-abcd",
		"table_group_name": "tf-example",
	})

	if err := resourceTencentCloudTcaplusTableGroupCreate(resourceData, &tcaplusTableGroupTestMeta{}); err == nil {
		t.Fatal("expect error when TableGroupId is empty, got nil")
	}
	if resourceData.Id() != "" {
		t.Fatalf("expect empty id, got %s", resourceData.Id())
	}
}

func TestResourceTencentCloudTcaplusTableGroupCreateApiError(t *testing.T) {
	_, sdkClient, patches := newTcaplusTableGroupTestClient(t)
	defer patches.Reset()

	patches.ApplyMethodFunc(sdkClient, "CreateTableGroupWithContext", func(_ context.Context, request *tcaplusdbv20190823.CreateTableGroupRequest) (*tcaplusdbv20190823.CreateTableGroupResponse, error) {
		return nil, errors.New("mock api error")
	})

	resourceData := schema.TestResourceDataRaw(t, ResourceTencentCloudTcaplusTableGroup().Schema, map[string]interface{}{
		"cluster_id":       "cluster-abcd",
		"table_group_name": "tf-example",
	})

	if err := resourceTencentCloudTcaplusTableGroupCreate(resourceData, &tcaplusTableGroupTestMeta{}); err == nil {
		t.Fatal("expect error when api returns error, got nil")
	}
	if resourceData.Id() != "" {
		t.Fatalf("expect empty id, got %s", resourceData.Id())
	}
}

func TestResourceTencentCloudTcaplusTableGroupUpdateImmutable(t *testing.T) {
	_, _, patches := newTcaplusTableGroupTestClient(t)
	defer patches.Reset()

	resourceData := schema.TestResourceDataRaw(t, ResourceTencentCloudTcaplusTableGroup().Schema, map[string]interface{}{
		"cluster_id":       "cluster-abcd",
		"table_group_name": "tf-example",
		"table_group_id":   "109",
	})
	resourceData.SetId("cluster-abcd" + tccommon.FILED_SP + "109")

	hasChangePatches := gomonkey.NewPatches()
	hasChangePatches.ApplyMethodReturn(resourceData, "HasChange", true)
	defer hasChangePatches.Reset()

	err := resourceTencentCloudTcaplusTableGroupUpdate(resourceData, &tcaplusTableGroupTestMeta{})
	if err == nil {
		t.Fatal("expect error when immutable argument changed, got nil")
	}
}

func TestResourceTencentCloudTcaplusTableGroupDeleteSuccess(t *testing.T) {
	_, sdkClient, patches := newTcaplusTableGroupTestClient(t)
	defer patches.Reset()

	var gotRequest *tcaplusdbv20190823.DeleteTableGroupRequest
	patches.ApplyMethodFunc(sdkClient, "DeleteTableGroupWithContext", func(_ context.Context, request *tcaplusdbv20190823.DeleteTableGroupRequest) (*tcaplusdbv20190823.DeleteTableGroupResponse, error) {
		gotRequest = request
		response := tcaplusdbv20190823.NewDeleteTableGroupResponse()
		response.Response = &tcaplusdbv20190823.DeleteTableGroupResponseParams{
			TaskId: helper.String("task-1"),
		}
		return response, nil
	})

	resourceData := schema.TestResourceDataRaw(t, ResourceTencentCloudTcaplusTableGroup().Schema, map[string]interface{}{
		"cluster_id":       "cluster-abcd",
		"table_group_name": "tf-example",
		"table_group_id":   "109",
	})
	resourceData.SetId("cluster-abcd" + tccommon.FILED_SP + "109")

	if err := resourceTencentCloudTcaplusTableGroupDelete(resourceData, &tcaplusTableGroupTestMeta{}); err != nil {
		t.Fatalf("expect no error, got %v", err)
	}

	if gotRequest == nil || gotRequest.ClusterId == nil || *gotRequest.ClusterId != "cluster-abcd" {
		t.Fatal("expect DeleteTableGroup request to carry cluster id cluster-abcd")
	}
	if gotRequest.TableGroupId == nil || *gotRequest.TableGroupId != "109" {
		t.Fatal("expect DeleteTableGroup request to carry table group id 109")
	}
}

func TestResourceTencentCloudTcaplusTableGroupDeleteBrokenId(t *testing.T) {
	resourceData := schema.TestResourceDataRaw(t, ResourceTencentCloudTcaplusTableGroup().Schema, map[string]interface{}{
		"cluster_id":       "cluster-abcd",
		"table_group_name": "tf-example",
	})
	resourceData.SetId("cluster-abcd")

	if err := resourceTencentCloudTcaplusTableGroupDelete(resourceData, &tcaplusTableGroupTestMeta{}); err == nil {
		t.Fatal("expect error for broken id, got nil")
	}
}
