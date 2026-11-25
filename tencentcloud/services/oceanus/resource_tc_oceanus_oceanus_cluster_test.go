package oceanus_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"

	oceanusv20190422 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/oceanus/v20190422"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/sharedmeta"
	svoceanus "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/oceanus"
)

const (
	oceanusClusterTestId          = "cluster-12345678"
	oceanusClusterTestName        = "tf-example-oceanus"
	oceanusClusterTestVpcId       = "vpc-12345678"
	oceanusClusterTestSubnetId    = "subnet-12345678"
	oceanusClusterTestCosBucket   = "tf-example-bucket-1250000000"
	oceanusClusterTestPassword    = "Tf-example-12345"
	oceanusClusterTestRemark      = "created by terraform"
	oceanusClusterTestRegionId    = int64(1)
	oceanusClusterTestZoneId      = int64(100003)
	oceanusClusterTestCu          = int64(12)
	oceanusClusterTestScaledCu    = int64(24)
	oceanusClusterTestCloudCu     = uint64(26)
	oceanusClusterTestChargeType  = "POSTPAID_BY_SECOND"
	oceanusClusterTestWebUiAccess = "NetworkAccess_EXTERNAL"
)

// oceanusClusterTestStr returns a pointer to the given string.
func oceanusClusterTestStr(v string) *string { return &v }

// oceanusClusterTestVpcObjectType mirrors the element type of the
// vpc_descriptions / slave_vpc_descriptions nested attribute declared by the
// resource schema.
func oceanusClusterTestVpcObjectType() types.ObjectType {
	return types.ObjectType{AttrTypes: map[string]attr.Type{
		"vpc_id":    types.StringType,
		"subnet_id": types.StringType,
		"app_id":    types.Int64Type,
		"owner_uin": types.StringType,
	}}
}

// oceanusClusterTestVpcList builds the vpc_descriptions value used by the
// plan/state fixtures.
func oceanusClusterTestVpcList() types.List {
	objectType := oceanusClusterTestVpcObjectType()

	return types.ListValueMust(objectType, []attr.Value{
		types.ObjectValueMust(objectType.AttrTypes, map[string]attr.Value{
			"vpc_id":    types.StringValue(oceanusClusterTestVpcId),
			"subnet_id": types.StringValue(oceanusClusterTestSubnetId),
			"app_id":    types.Int64Null(),
			"owner_uin": types.StringNull(),
		}),
	})
}

// oceanusClusterTestModel builds the terraform model used for plan/state.
// Every attribute the fixture does not configure carries an explicit typed
// null value: the raw zero value (types.List{} etc.) breaks tfsdk Plan/State
// reflection with "Value Conversion Error".
func oceanusClusterTestModel() svoceanus.OceanusOceanusClusterResourceModel {
	return svoceanus.OceanusOceanusClusterResourceModel{
		Id:                          types.StringValue(oceanusClusterTestId),
		ClusterName:                 types.StringValue(oceanusClusterTestName),
		RegionId:                    types.Int64Value(oceanusClusterTestRegionId),
		ZoneId:                      types.Int64Value(oceanusClusterTestZoneId),
		LoginPassword:               types.StringValue(oceanusClusterTestPassword),
		VpcDescriptions:             oceanusClusterTestVpcList(),
		DefaultCosBucket:            types.StringValue(oceanusClusterTestCosBucket),
		Cu:                          types.Int64Value(oceanusClusterTestCu),
		Remark:                      types.StringValue(oceanusClusterTestRemark),
		Period:                      types.Int64Null(),
		InstanceChargeType:          types.StringValue(oceanusClusterTestChargeType),
		ClusterType:                 types.StringNull(),
		RenewFlag:                   types.StringNull(),
		FlinkWebUiNetworkAccessType: types.StringValue(oceanusClusterTestWebUiAccess),
		SlaveVpcDescriptions:        types.ListNull(oceanusClusterTestVpcObjectType()),
		CuMemory:                    types.Int64Null(),
	}
}

// oceanusClusterTestItem builds the DescribeClusters response item.
func oceanusClusterTestItem() *oceanusv20190422.Cluster {
	cuNum := oceanusClusterTestCloudCu

	return &oceanusv20190422.Cluster{
		ClusterId:        oceanusClusterTestStr(oceanusClusterTestId),
		Name:             oceanusClusterTestStr(oceanusClusterTestName),
		Remark:           oceanusClusterTestStr(oceanusClusterTestRemark),
		DefaultCOSBucket: oceanusClusterTestStr(oceanusClusterTestCosBucket),
		CuNum:            &cuNum,
		Status:           oceanusClusterTestInt64(2),
		StatusDesc:       oceanusClusterTestStr("running"),
	}
}

// oceanusClusterTestInt64 returns a pointer to the given int64.
func oceanusClusterTestInt64(v int64) *int64 { return &v }

// oceanusClusterTestContext wires the resource to a configured client; the SDK
// client methods are stubbed per test case.
type oceanusClusterTestContext struct {
	resource resource.Resource
	schema   schema.Schema
	client   *connectivity.TencentCloudClient
}

func newOceanusClusterTestContext(t *testing.T) *oceanusClusterTestContext {
	t.Helper()

	connectivityClient := &connectivity.TencentCloudClient{}

	res := svoceanus.NewOceanusOceanusClusterResource()
	configured, ok := res.(resource.ResourceWithConfigure)
	if !ok {
		t.Fatalf("resource does not implement resource.ResourceWithConfigure")
	}

	var configureResp resource.ConfigureResponse
	configured.Configure(
		context.Background(),
		resource.ConfigureRequest{ProviderData: &sharedmeta.ProviderMeta{Client: connectivityClient}},
		&configureResp,
	)
	if configureResp.Diagnostics.HasError() {
		t.Fatalf("failed to configure resource: %v", configureResp.Diagnostics.Errors())
	}

	var schemaResp resource.SchemaResponse
	res.Schema(context.Background(), resource.SchemaRequest{}, &schemaResp)
	if schemaResp.Diagnostics.HasError() {
		t.Fatalf("failed to build schema: %v", schemaResp.Diagnostics.Errors())
	}

	return &oceanusClusterTestContext{resource: res, schema: schemaResp.Schema, client: connectivityClient}
}

func (c *oceanusClusterTestContext) plan(t *testing.T, m svoceanus.OceanusOceanusClusterResourceModel) tfsdk.Plan {
	t.Helper()

	plan := tfsdk.Plan{Schema: c.schema}
	if diags := plan.Set(context.Background(), m); diags.HasError() {
		t.Fatalf("failed to build plan: %v", diags.Errors())
	}

	return plan
}

func (c *oceanusClusterTestContext) state(t *testing.T, m svoceanus.OceanusOceanusClusterResourceModel) tfsdk.State {
	t.Helper()

	state := tfsdk.State{Schema: c.schema}
	if diags := state.Set(context.Background(), m); diags.HasError() {
		t.Fatalf("failed to build state: %v", diags.Errors())
	}

	return state
}

// patchOceanusClusterSdkClient stubs the Oceanus client accessor and returns
// the mock SDK client whose methods each test case patches next.
func patchOceanusClusterSdkClient(patches *gomonkey.Patches, client *connectivity.TencentCloudClient) *oceanusv20190422.Client {
	oceanusClient := &oceanusv20190422.Client{}
	patches.ApplyMethodReturn(client, "UseOceanusClient", oceanusClient)

	return oceanusClient
}

// TestOceanusClusterMetadata verifies the registered type name.
func TestOceanusClusterMetadata(t *testing.T) {
	ctx := newOceanusClusterTestContext(t)

	var metadataResp resource.MetadataResponse
	ctx.resource.Metadata(context.Background(), resource.MetadataRequest{}, &metadataResp)
	if metadataResp.TypeName != "tencentcloud_oceanus_oceanus_cluster" {
		t.Errorf("type name = %q, want tencentcloud_oceanus_oceanus_cluster", metadataResp.TypeName)
	}
}

// TestOceanusClusterCreateSuccess verifies that Create sends every configured
// field, tracks the returned cluster id and refreshes the state.
func TestOceanusClusterCreateSuccess(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ctx := newOceanusClusterTestContext(t)
	oceanusClient := patchOceanusClusterSdkClient(patches, ctx.client)

	var captured *oceanusv20190422.CreateOceanusClusterRequest
	patches.ApplyMethodFunc(oceanusClient, "CreateOceanusCluster", func(req *oceanusv20190422.CreateOceanusClusterRequest) (*oceanusv20190422.CreateOceanusClusterResponse, error) {
		captured = req

		return &oceanusv20190422.CreateOceanusClusterResponse{Response: &oceanusv20190422.CreateOceanusClusterResponseParams{
			ClusterId: oceanusClusterTestStr(oceanusClusterTestId),
			RequestId: oceanusClusterTestStr("fake-request-id"),
		}}, nil
	})
	patches.ApplyMethodFunc(oceanusClient, "DescribeClusters", func(req *oceanusv20190422.DescribeClustersRequest) (*oceanusv20190422.DescribeClustersResponse, error) {
		return &oceanusv20190422.DescribeClustersResponse{Response: &oceanusv20190422.DescribeClustersResponseParams{
			ClusterSet: []*oceanusv20190422.Cluster{oceanusClusterTestItem()},
			RequestId:  oceanusClusterTestStr("fake-request-id"),
		}}, nil
	})

	model := oceanusClusterTestModel()
	var createResp resource.CreateResponse
	createResp.State = ctx.state(t, model)
	ctx.resource.Create(context.Background(), resource.CreateRequest{Plan: ctx.plan(t, model)}, &createResp)
	if createResp.Diagnostics.HasError() {
		t.Fatalf("Create failed: %v", createResp.Diagnostics.Errors())
	}

	if captured == nil {
		t.Fatal("CreateOceanusCluster request is nil")
	}
	if captured.ClusterName == nil || *captured.ClusterName != oceanusClusterTestName {
		t.Errorf("create request cluster_name = %v, want %q", captured.ClusterName, oceanusClusterTestName)
	}
	if captured.RegionId == nil || *captured.RegionId != oceanusClusterTestRegionId {
		t.Errorf("create request region_id = %v, want %d", captured.RegionId, oceanusClusterTestRegionId)
	}
	if captured.ZoneId == nil || *captured.ZoneId != oceanusClusterTestZoneId {
		t.Errorf("create request zone_id = %v, want %d", captured.ZoneId, oceanusClusterTestZoneId)
	}
	if captured.LoginPassword == nil || *captured.LoginPassword != oceanusClusterTestPassword {
		t.Errorf("create request login_password = %v, want %q", captured.LoginPassword, oceanusClusterTestPassword)
	}
	if captured.DefaultCOSBucket == nil || *captured.DefaultCOSBucket != oceanusClusterTestCosBucket {
		t.Errorf("create request default_cos_bucket = %v, want %q", captured.DefaultCOSBucket, oceanusClusterTestCosBucket)
	}
	if len(captured.VpcDescriptions) != 1 || captured.VpcDescriptions[0] == nil ||
		captured.VpcDescriptions[0].VpcId == nil || *captured.VpcDescriptions[0].VpcId != oceanusClusterTestVpcId ||
		captured.VpcDescriptions[0].SubnetId == nil || *captured.VpcDescriptions[0].SubnetId != oceanusClusterTestSubnetId {
		t.Errorf("create request vpc_descriptions = %v, want one entry %q/%q", captured.VpcDescriptions, oceanusClusterTestVpcId, oceanusClusterTestSubnetId)
	}
	if captured.CU == nil || *captured.CU != oceanusClusterTestCu {
		t.Errorf("create request cu = %v, want %d", captured.CU, oceanusClusterTestCu)
	}
	if captured.Remark == nil || *captured.Remark != oceanusClusterTestRemark {
		t.Errorf("create request remark = %v, want %q", captured.Remark, oceanusClusterTestRemark)
	}
	if captured.InstanceChargeType == nil || *captured.InstanceChargeType != oceanusClusterTestChargeType {
		t.Errorf("create request instance_charge_type = %v, want %q", captured.InstanceChargeType, oceanusClusterTestChargeType)
	}
	if captured.FlinkWebUINetworkAccessType == nil || *captured.FlinkWebUINetworkAccessType != oceanusClusterTestWebUiAccess {
		t.Errorf("create request flink_web_ui_network_access_type = %v, want %q", captured.FlinkWebUINetworkAccessType, oceanusClusterTestWebUiAccess)
	}
	if captured.Period != nil {
		t.Errorf("create request period = %v, want it unset", captured.Period)
	}
	if captured.CUMemory != nil {
		t.Errorf("create request cu_memory = %v, want it unset", captured.CUMemory)
	}

	var got svoceanus.OceanusOceanusClusterResourceModel
	if diags := createResp.State.Get(context.Background(), &got); diags.HasError() {
		t.Fatalf("failed to read back state: %v", diags.Errors())
	}
	if got.Id.ValueString() != oceanusClusterTestId {
		t.Errorf("id = %q, want %q", got.Id.ValueString(), oceanusClusterTestId)
	}
	if got.Cu.ValueInt64() != int64(oceanusClusterTestCloudCu) {
		t.Errorf("cu = %d, want the cloud value %d", got.Cu.ValueInt64(), oceanusClusterTestCloudCu)
	}
	if got.ClusterName.ValueString() != oceanusClusterTestName {
		t.Errorf("cluster_name = %q, want %q", got.ClusterName.ValueString(), oceanusClusterTestName)
	}
}

// TestOceanusClusterCreateEmptyId verifies that Create fails loudly when the
// API returns no cluster id.
func TestOceanusClusterCreateEmptyId(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ctx := newOceanusClusterTestContext(t)
	oceanusClient := patchOceanusClusterSdkClient(patches, ctx.client)
	patches.ApplyMethodFunc(oceanusClient, "CreateOceanusCluster", func(req *oceanusv20190422.CreateOceanusClusterRequest) (*oceanusv20190422.CreateOceanusClusterResponse, error) {
		return &oceanusv20190422.CreateOceanusClusterResponse{Response: &oceanusv20190422.CreateOceanusClusterResponseParams{
			RequestId: oceanusClusterTestStr("fake-request-id"),
		}}, nil
	})

	model := oceanusClusterTestModel()
	var createResp resource.CreateResponse
	createResp.State = ctx.state(t, model)
	ctx.resource.Create(context.Background(), resource.CreateRequest{Plan: ctx.plan(t, model)}, &createResp)
	if !createResp.Diagnostics.HasError() {
		t.Fatal("Create should fail when the API returns an empty cluster id")
	}
}

// TestOceanusClusterCreateApiError verifies API failures surface as errors.
func TestOceanusClusterCreateApiError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ctx := newOceanusClusterTestContext(t)
	oceanusClient := patchOceanusClusterSdkClient(patches, ctx.client)
	patches.ApplyMethodFunc(oceanusClient, "CreateOceanusCluster", func(req *oceanusv20190422.CreateOceanusClusterRequest) (*oceanusv20190422.CreateOceanusClusterResponse, error) {
		return nil, fmt.Errorf("create failed")
	})

	model := oceanusClusterTestModel()
	var createResp resource.CreateResponse
	createResp.State = ctx.state(t, model)
	ctx.resource.Create(context.Background(), resource.CreateRequest{Plan: ctx.plan(t, model)}, &createResp)
	if !createResp.Diagnostics.HasError() {
		t.Fatal("Create should fail on an API error")
	}
}

// TestOceanusClusterReadSuccess verifies Read refreshes the cloud-owned fields.
func TestOceanusClusterReadSuccess(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ctx := newOceanusClusterTestContext(t)
	oceanusClient := patchOceanusClusterSdkClient(patches, ctx.client)

	var gotClusterId string
	patches.ApplyMethodFunc(oceanusClient, "DescribeClusters", func(req *oceanusv20190422.DescribeClustersRequest) (*oceanusv20190422.DescribeClustersResponse, error) {
		if len(req.ClusterIds) != 1 || req.ClusterIds[0] == nil {
			t.Errorf("expected the cluster id filter to be passed")
		} else {
			gotClusterId = *req.ClusterIds[0]
		}

		return &oceanusv20190422.DescribeClustersResponse{Response: &oceanusv20190422.DescribeClustersResponseParams{
			ClusterSet: []*oceanusv20190422.Cluster{oceanusClusterTestItem()},
			RequestId:  oceanusClusterTestStr("fake-request-id"),
		}}, nil
	})

	model := oceanusClusterTestModel()
	var readResp resource.ReadResponse
	readResp.State = ctx.state(t, model)
	ctx.resource.Read(context.Background(), resource.ReadRequest{State: ctx.state(t, model)}, &readResp)
	if readResp.Diagnostics.HasError() {
		t.Fatalf("Read failed: %v", readResp.Diagnostics.Errors())
	}

	if gotClusterId != oceanusClusterTestId {
		t.Errorf("DescribeClusters called with %q, want %q", gotClusterId, oceanusClusterTestId)
	}

	var got svoceanus.OceanusOceanusClusterResourceModel
	if diags := readResp.State.Get(context.Background(), &got); diags.HasError() {
		t.Fatalf("failed to read back state: %v", diags.Errors())
	}
	if got.Cu.ValueInt64() != int64(oceanusClusterTestCloudCu) {
		t.Errorf("cu = %d, want the cloud value %d", got.Cu.ValueInt64(), oceanusClusterTestCloudCu)
	}
}

// TestOceanusClusterReadNotFound verifies that a describe result without the
// tracked cluster drops the resource from state (confirmed not-found).
func TestOceanusClusterReadNotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ctx := newOceanusClusterTestContext(t)
	oceanusClient := patchOceanusClusterSdkClient(patches, ctx.client)
	patches.ApplyMethodFunc(oceanusClient, "DescribeClusters", func(req *oceanusv20190422.DescribeClustersRequest) (*oceanusv20190422.DescribeClustersResponse, error) {
		return &oceanusv20190422.DescribeClustersResponse{Response: &oceanusv20190422.DescribeClustersResponseParams{
			RequestId: oceanusClusterTestStr("fake-request-id"),
		}}, nil
	})

	model := oceanusClusterTestModel()
	var readResp resource.ReadResponse
	readResp.State = ctx.state(t, model)
	ctx.resource.Read(context.Background(), resource.ReadRequest{State: ctx.state(t, model)}, &readResp)
	if readResp.Diagnostics.HasError() {
		t.Fatalf("Read should not error on confirmed not-found: %v", readResp.Diagnostics.Errors())
	}
	if !readResp.State.Raw.IsNull() {
		t.Error("state should be removed after the cluster is confirmed gone")
	}
}

// TestOceanusClusterReadApiError verifies describe failures surface as errors.
func TestOceanusClusterReadApiError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ctx := newOceanusClusterTestContext(t)
	oceanusClient := patchOceanusClusterSdkClient(patches, ctx.client)
	patches.ApplyMethodFunc(oceanusClient, "DescribeClusters", func(req *oceanusv20190422.DescribeClustersRequest) (*oceanusv20190422.DescribeClustersResponse, error) {
		return nil, fmt.Errorf("describe failed")
	})

	model := oceanusClusterTestModel()
	var readResp resource.ReadResponse
	readResp.State = ctx.state(t, model)
	ctx.resource.Read(context.Background(), resource.ReadRequest{State: ctx.state(t, model)}, &readResp)
	if !readResp.Diagnostics.HasError() {
		t.Fatal("Read should fail on an API error")
	}
}

// TestOceanusClusterUpdateScaleUp verifies that a `cu` change calls
// ScaleOceanusCluster with the planned target and records the POST-update
// cloud value from the read-back rather than the plan.
func TestOceanusClusterUpdateScaleUp(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ctx := newOceanusClusterTestContext(t)
	oceanusClient := patchOceanusClusterSdkClient(patches, ctx.client)

	var captured *oceanusv20190422.ScaleOceanusClusterRequest
	patches.ApplyMethodFunc(oceanusClient, "ScaleOceanusCluster", func(req *oceanusv20190422.ScaleOceanusClusterRequest) (*oceanusv20190422.ScaleOceanusClusterResponse, error) {
		captured = req

		return &oceanusv20190422.ScaleOceanusClusterResponse{Response: &oceanusv20190422.ScaleOceanusClusterResponseParams{
			TaskExecResult: oceanusClusterTestStr("success"),
			RequestId:      oceanusClusterTestStr("fake-request-id"),
		}}, nil
	})
	patches.ApplyMethodFunc(oceanusClient, "DescribeClusters", func(req *oceanusv20190422.DescribeClustersRequest) (*oceanusv20190422.DescribeClustersResponse, error) {
		return &oceanusv20190422.DescribeClustersResponse{Response: &oceanusv20190422.DescribeClustersResponseParams{
			ClusterSet: []*oceanusv20190422.Cluster{oceanusClusterTestItem()},
			RequestId:  oceanusClusterTestStr("fake-request-id"),
		}}, nil
	})

	stateModel := oceanusClusterTestModel()
	plannedModel := oceanusClusterTestModel()
	plannedModel.Cu = types.Int64Value(oceanusClusterTestScaledCu)

	var updateResp resource.UpdateResponse
	updateResp.State = ctx.state(t, stateModel)
	ctx.resource.Update(context.Background(),
		resource.UpdateRequest{
			Plan:  ctx.plan(t, plannedModel),
			State: ctx.state(t, stateModel),
		}, &updateResp)
	if updateResp.Diagnostics.HasError() {
		t.Fatalf("Update failed: %v", updateResp.Diagnostics.Errors())
	}

	if captured == nil {
		t.Fatal("ScaleOceanusCluster request is nil")
	}
	if captured.ClusterId == nil || *captured.ClusterId != oceanusClusterTestId {
		t.Errorf("scale request cluster_id = %v, want %q", captured.ClusterId, oceanusClusterTestId)
	}
	if captured.NewCU == nil || *captured.NewCU != oceanusClusterTestScaledCu {
		t.Errorf("scale request new_cu = %v, want the planned %d", captured.NewCU, oceanusClusterTestScaledCu)
	}
	if captured.ScaleMode == nil || *captured.ScaleMode != "ScaleUp" {
		t.Errorf("scale request scale_mode = %v, want ScaleUp", captured.ScaleMode)
	}

	var got svoceanus.OceanusOceanusClusterResourceModel
	if diags := updateResp.State.Get(context.Background(), &got); diags.HasError() {
		t.Fatalf("failed to read back state: %v", diags.Errors())
	}
	if got.Cu.ValueInt64() != int64(oceanusClusterTestCloudCu) {
		t.Errorf("cu = %d, want the read-back cloud value %d", got.Cu.ValueInt64(), oceanusClusterTestCloudCu)
	}
}

// TestOceanusClusterUpdateScaleDown verifies that a decreasing `cu` change is
// sent with the ScaleDown direction.
func TestOceanusClusterUpdateScaleDown(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ctx := newOceanusClusterTestContext(t)
	oceanusClient := patchOceanusClusterSdkClient(patches, ctx.client)

	var captured *oceanusv20190422.ScaleOceanusClusterRequest
	patches.ApplyMethodFunc(oceanusClient, "ScaleOceanusCluster", func(req *oceanusv20190422.ScaleOceanusClusterRequest) (*oceanusv20190422.ScaleOceanusClusterResponse, error) {
		captured = req

		return &oceanusv20190422.ScaleOceanusClusterResponse{Response: &oceanusv20190422.ScaleOceanusClusterResponseParams{
			TaskExecResult: oceanusClusterTestStr("success"),
			RequestId:      oceanusClusterTestStr("fake-request-id"),
		}}, nil
	})
	patches.ApplyMethodFunc(oceanusClient, "DescribeClusters", func(req *oceanusv20190422.DescribeClustersRequest) (*oceanusv20190422.DescribeClustersResponse, error) {
		return &oceanusv20190422.DescribeClustersResponse{Response: &oceanusv20190422.DescribeClustersResponseParams{
			ClusterSet: []*oceanusv20190422.Cluster{oceanusClusterTestItem()},
			RequestId:  oceanusClusterTestStr("fake-request-id"),
		}}, nil
	})

	stateModel := oceanusClusterTestModel()
	stateModel.Cu = types.Int64Value(oceanusClusterTestScaledCu)
	plannedModel := oceanusClusterTestModel()

	var updateResp resource.UpdateResponse
	updateResp.State = ctx.state(t, stateModel)
	ctx.resource.Update(context.Background(),
		resource.UpdateRequest{
			Plan:  ctx.plan(t, plannedModel),
			State: ctx.state(t, stateModel),
		}, &updateResp)
	if updateResp.Diagnostics.HasError() {
		t.Fatalf("Update failed: %v", updateResp.Diagnostics.Errors())
	}

	if captured == nil {
		t.Fatal("ScaleOceanusCluster request is nil")
	}
	if captured.ScaleMode == nil || *captured.ScaleMode != "ScaleDown" {
		t.Errorf("scale request scale_mode = %v, want ScaleDown", captured.ScaleMode)
	}
}

// TestOceanusClusterUpdateNoCuChange verifies that an update without a `cu`
// change calls NO write API: no attribute the scale API covers changed, so the
// call would be a pointless no-change request. The read-back still runs.
func TestOceanusClusterUpdateNoCuChange(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ctx := newOceanusClusterTestContext(t)
	oceanusClient := patchOceanusClusterSdkClient(patches, ctx.client)

	scaleCalls := 0
	patches.ApplyMethodFunc(oceanusClient, "ScaleOceanusCluster", func(req *oceanusv20190422.ScaleOceanusClusterRequest) (*oceanusv20190422.ScaleOceanusClusterResponse, error) {
		scaleCalls++

		return nil, fmt.Errorf("the scale API must not be called")
	})
	patches.ApplyMethodFunc(oceanusClient, "DescribeClusters", func(req *oceanusv20190422.DescribeClustersRequest) (*oceanusv20190422.DescribeClustersResponse, error) {
		return &oceanusv20190422.DescribeClustersResponse{Response: &oceanusv20190422.DescribeClustersResponseParams{
			ClusterSet: []*oceanusv20190422.Cluster{oceanusClusterTestItem()},
			RequestId:  oceanusClusterTestStr("fake-request-id"),
		}}, nil
	})

	model := oceanusClusterTestModel()

	var updateResp resource.UpdateResponse
	updateResp.State = ctx.state(t, model)
	ctx.resource.Update(context.Background(),
		resource.UpdateRequest{
			Plan:  ctx.plan(t, model),
			State: ctx.state(t, model),
		}, &updateResp)
	if updateResp.Diagnostics.HasError() {
		t.Fatalf("Update failed: %v", updateResp.Diagnostics.Errors())
	}
	if scaleCalls != 0 {
		t.Errorf("ScaleOceanusCluster was called %d time(s) without a cu change", scaleCalls)
	}
}

// TestOceanusClusterUpdateApiError verifies scale failures surface as errors.
func TestOceanusClusterUpdateApiError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ctx := newOceanusClusterTestContext(t)
	oceanusClient := patchOceanusClusterSdkClient(patches, ctx.client)
	patches.ApplyMethodFunc(oceanusClient, "ScaleOceanusCluster", func(req *oceanusv20190422.ScaleOceanusClusterRequest) (*oceanusv20190422.ScaleOceanusClusterResponse, error) {
		return nil, fmt.Errorf("scale failed")
	})

	stateModel := oceanusClusterTestModel()
	plannedModel := oceanusClusterTestModel()
	plannedModel.Cu = types.Int64Value(oceanusClusterTestScaledCu)

	var updateResp resource.UpdateResponse
	updateResp.State = ctx.state(t, stateModel)
	ctx.resource.Update(context.Background(),
		resource.UpdateRequest{
			Plan:  ctx.plan(t, plannedModel),
			State: ctx.state(t, stateModel),
		}, &updateResp)
	if !updateResp.Diagnostics.HasError() {
		t.Fatal("Update should fail on an API error")
	}
}

// TestOceanusClusterDeleteSuccess verifies Delete sends the tracked cluster id.
func TestOceanusClusterDeleteSuccess(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ctx := newOceanusClusterTestContext(t)
	oceanusClient := patchOceanusClusterSdkClient(patches, ctx.client)

	var captured *oceanusv20190422.DeleteOceanusClusterRequest
	patches.ApplyMethodFunc(oceanusClient, "DeleteOceanusCluster", func(req *oceanusv20190422.DeleteOceanusClusterRequest) (*oceanusv20190422.DeleteOceanusClusterResponse, error) {
		captured = req

		return &oceanusv20190422.DeleteOceanusClusterResponse{Response: &oceanusv20190422.DeleteOceanusClusterResponseParams{
			TaskExecResult: oceanusClusterTestStr("success"),
			RequestId:      oceanusClusterTestStr("fake-request-id"),
		}}, nil
	})

	model := oceanusClusterTestModel()
	var deleteResp resource.DeleteResponse
	ctx.resource.Delete(context.Background(), resource.DeleteRequest{State: ctx.state(t, model)}, &deleteResp)
	if deleteResp.Diagnostics.HasError() {
		t.Fatalf("Delete failed: %v", deleteResp.Diagnostics.Errors())
	}

	if captured == nil {
		t.Fatal("DeleteOceanusCluster request is nil")
	}
	if captured.ClusterId == nil || *captured.ClusterId != oceanusClusterTestId {
		t.Errorf("delete request cluster_id = %v, want %q", captured.ClusterId, oceanusClusterTestId)
	}
}

// TestOceanusClusterDeleteApiError verifies delete failures surface as errors.
func TestOceanusClusterDeleteApiError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ctx := newOceanusClusterTestContext(t)
	oceanusClient := patchOceanusClusterSdkClient(patches, ctx.client)
	patches.ApplyMethodFunc(oceanusClient, "DeleteOceanusCluster", func(req *oceanusv20190422.DeleteOceanusClusterRequest) (*oceanusv20190422.DeleteOceanusClusterResponse, error) {
		return nil, fmt.Errorf("delete failed")
	})

	model := oceanusClusterTestModel()
	var deleteResp resource.DeleteResponse
	ctx.resource.Delete(context.Background(), resource.DeleteRequest{State: ctx.state(t, model)}, &deleteResp)
	if !deleteResp.Diagnostics.HasError() {
		t.Fatal("Delete should fail on an API error")
	}
}
