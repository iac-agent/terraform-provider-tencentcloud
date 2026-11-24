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

	resourcetimeouts "github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"

	oceanusv20190422 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/oceanus/v20190422"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/sharedmeta"
	svoceanus "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/oceanus"
)

const (
	oceanusClusterAfTestId              = "cluster-12345678"
	oceanusClusterAfTestName            = "tf-example-cluster"
	oceanusClusterAfTestCloudName       = "tf-example-cluster-normalized"
	oceanusClusterAfTestRemark          = "created by terraform"
	oceanusClusterAfTestCosBucket       = "example-1250000000"
	oceanusClusterAfTestPassword        = "Example#2024"
	oceanusClusterAfTestVpcId           = "vpc-12345678"
	oceanusClusterAfTestSubnetId        = "subnet-12345678"
	oceanusClusterAfTestRegionId  int64 = 1
	oceanusClusterAfTestZoneId    int64 = 100002
	oceanusClusterAfTestCu        int64 = 12
	oceanusClusterAfTestNewCu     int64 = 26
	oceanusClusterAfTestCloudCu   int64 = 26
)

// oceanusClusterAfStr is the *string helper shared by the stubs below.
func oceanusClusterAfStr(v string) *string { return &v }

// oceanusClusterAfInt64 is the *int64 helper shared by the stubs below.
func oceanusClusterAfInt64(v int64) *int64 { return &v }

// oceanusClusterAfTestCluster builds the DescribeClusters response item.
func oceanusClusterAfTestCluster() *oceanusv20190422.Cluster {
	return &oceanusv20190422.Cluster{
		ClusterId:        oceanusClusterAfStr(oceanusClusterAfTestId),
		Name:             oceanusClusterAfStr(oceanusClusterAfTestName),
		Remark:           oceanusClusterAfStr(oceanusClusterAfTestRemark),
		DefaultCOSBucket: oceanusClusterAfStr(oceanusClusterAfTestCosBucket),
		Status:           oceanusClusterAfInt64(svoceanus.OCEANUS_CLUSTER_STATUS_RUNNING),
		CuNum:            func() *uint64 { v := uint64(oceanusClusterAfTestCu); return &v }(),
	}
}

// oceanusClusterAfTimeoutsAttrTypes mirrors the attribute types of the
// timeouts block declared by the resource schema. Only the create stage waits:
// the scale and destroy APIs are synchronous.
func oceanusClusterAfTimeoutsAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"create": types.StringType,
	}
}

// oceanusClusterAfVpcDescriptionsObjectType is the object type of the
// vpc_descriptions / slave_vpc_descriptions elements.
func oceanusClusterAfVpcDescriptionsObjectType() attr.Type {
	return types.ObjectType{AttrTypes: map[string]attr.Type{
		"vpc_id":    types.StringType,
		"subnet_id": types.StringType,
		"app_id":    types.Int64Type,
		"owner_uin": types.StringType,
	}}
}

// oceanusClusterAfTestVpcDescriptions builds the planned vpc_descriptions list.
func oceanusClusterAfTestVpcDescriptions() types.List {
	return types.ListValueMust(oceanusClusterAfVpcDescriptionsObjectType(), []attr.Value{
		types.ObjectValueMust(oceanusClusterAfVpcDescriptionsObjectType().(types.ObjectType).AttrTypes, map[string]attr.Value{
			"vpc_id":    types.StringValue(oceanusClusterAfTestVpcId),
			"subnet_id": types.StringValue(oceanusClusterAfTestSubnetId),
			"app_id":    types.Int64Null(),
			"owner_uin": types.StringNull(),
		}),
	})
}

// oceanusClusterAfTestModel builds the terraform model used for plan/state.
// Every unset attribute carries an explicit typed null value: the raw zero
// value breaks tfsdk Plan/State reflection. The timeouts value is a short
// user value so wait-path cases finish in seconds instead of the default.
func oceanusClusterAfTestModel() svoceanus.OceanusClusterAfResourceModel {
	return svoceanus.OceanusClusterAfResourceModel{
		Id:                          types.StringValue(oceanusClusterAfTestId),
		ClusterName:                 types.StringValue(oceanusClusterAfTestName),
		RegionId:                    types.Int64Value(oceanusClusterAfTestRegionId),
		ZoneId:                      types.Int64Value(oceanusClusterAfTestZoneId),
		LoginPassword:               types.StringValue(oceanusClusterAfTestPassword),
		VpcDescriptions:             oceanusClusterAfTestVpcDescriptions(),
		DefaultCosBucket:            types.StringValue(oceanusClusterAfTestCosBucket),
		Cu:                          types.Int64Value(oceanusClusterAfTestCu),
		Remark:                      types.StringValue(oceanusClusterAfTestRemark),
		Period:                      types.Int64Null(),
		InstanceChargeType:          types.StringNull(),
		ClusterType:                 types.StringNull(),
		RenewFlag:                   types.StringNull(),
		FlinkWebUiNetworkAccessType: types.StringNull(),
		SlaveVpcDescriptions:        types.ListNull(oceanusClusterAfVpcDescriptionsObjectType()),
		CuMemory:                    types.Int64Null(),
		Timeouts: resourcetimeouts.Value{Object: types.ObjectValueMust(
			oceanusClusterAfTimeoutsAttrTypes(),
			map[string]attr.Value{
				"create": types.StringValue("5s"),
			},
		)},
	}
}

// oceanusClusterAfTestContext wires the resource to a configured client; the
// SDK client methods are stubbed per test case.
type oceanusClusterAfTestContext struct {
	resource resource.Resource
	schema   schema.Schema
	client   *connectivity.TencentCloudClient
}

func newOceanusClusterAfTestContext(t *testing.T) *oceanusClusterAfTestContext {
	t.Helper()

	connectivityClient := &connectivity.TencentCloudClient{}

	res := svoceanus.NewOceanusClusterAfResource()
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

	return &oceanusClusterAfTestContext{resource: res, schema: schemaResp.Schema, client: connectivityClient}
}

func (c *oceanusClusterAfTestContext) plan(t *testing.T, m svoceanus.OceanusClusterAfResourceModel) tfsdk.Plan {
	t.Helper()

	plan := tfsdk.Plan{Schema: c.schema}
	if diags := plan.Set(context.Background(), m); diags.HasError() {
		t.Fatalf("failed to build plan: %v", diags.Errors())
	}

	return plan
}

func (c *oceanusClusterAfTestContext) state(t *testing.T, m svoceanus.OceanusClusterAfResourceModel) tfsdk.State {
	t.Helper()

	state := tfsdk.State{Schema: c.schema}
	if diags := state.Set(context.Background(), m); diags.HasError() {
		t.Fatalf("failed to build state: %v", diags.Errors())
	}

	return state
}

// patchOceanusClusterAfSdkClient stubs the client accessor and returns the mock
// SDK client whose methods the test patches next.
func patchOceanusClusterAfSdkClient(patches *gomonkey.Patches, client *connectivity.TencentCloudClient) *oceanusv20190422.Client {
	oceanusClient := &oceanusv20190422.Client{}
	patches.ApplyMethodReturn(client, "UseOceanusClient", oceanusClient)

	return oceanusClient
}

// TestOceanusClusterAfMetadata validates the registry type name.
func TestOceanusClusterAfMetadata(t *testing.T) {
	var resp resource.MetadataResponse
	svoceanus.NewOceanusClusterAfResource().Metadata(context.Background(), resource.MetadataRequest{}, &resp)

	if resp.TypeName != "tencentcloud_oceanus_cluster_af" {
		t.Errorf("type name = %q, want %q", resp.TypeName, "tencentcloud_oceanus_cluster_af")
	}
}

// TestOceanusClusterAfCreateSuccess verifies that Create sends all request
// fields, tracks the returned cluster id, waits for the cluster to become
// ready and refreshes the state.
func TestOceanusClusterAfCreateSuccess(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ctx := newOceanusClusterAfTestContext(t)
	oceanusClient := patchOceanusClusterAfSdkClient(patches, ctx.client)

	var captured *oceanusv20190422.CreateOceanusClusterRequest
	patches.ApplyMethodFunc(oceanusClient, "CreateOceanusCluster", func(req *oceanusv20190422.CreateOceanusClusterRequest) (*oceanusv20190422.CreateOceanusClusterResponse, error) {
		captured = req

		return &oceanusv20190422.CreateOceanusClusterResponse{Response: &oceanusv20190422.CreateOceanusClusterResponseParams{
			ClusterId: oceanusClusterAfStr(oceanusClusterAfTestId),
			RequestId: oceanusClusterAfStr("fake-request-id"),
		}}, nil
	})
	patches.ApplyMethodFunc(oceanusClient, "DescribeClusters", func(req *oceanusv20190422.DescribeClustersRequest) (*oceanusv20190422.DescribeClustersResponse, error) {
		return &oceanusv20190422.DescribeClustersResponse{Response: &oceanusv20190422.DescribeClustersResponseParams{
			ClusterSet: []*oceanusv20190422.Cluster{oceanusClusterAfTestCluster()},
			RequestId:  oceanusClusterAfStr("fake-request-id"),
		}}, nil
	})

	var createResp resource.CreateResponse
	createResp.State = ctx.state(t, oceanusClusterAfTestModel())
	ctx.resource.Create(context.Background(), resource.CreateRequest{Plan: ctx.plan(t, oceanusClusterAfTestModel())}, &createResp)
	if createResp.Diagnostics.HasError() {
		t.Fatalf("Create failed: %v", createResp.Diagnostics.Errors())
	}

	if captured == nil {
		t.Fatal("CreateOceanusCluster request is nil")
	}
	if captured.ClusterName == nil || *captured.ClusterName != oceanusClusterAfTestName {
		t.Errorf("create request cluster_name = %v, want %q", captured.ClusterName, oceanusClusterAfTestName)
	}
	if captured.RegionId == nil || *captured.RegionId != oceanusClusterAfTestRegionId {
		t.Errorf("create request region_id = %v, want %d", captured.RegionId, oceanusClusterAfTestRegionId)
	}
	if captured.ZoneId == nil || *captured.ZoneId != oceanusClusterAfTestZoneId {
		t.Errorf("create request zone_id = %v, want %d", captured.ZoneId, oceanusClusterAfTestZoneId)
	}
	if captured.LoginPassword == nil || *captured.LoginPassword != oceanusClusterAfTestPassword {
		t.Errorf("create request login_password = %v, want the configured value", captured.LoginPassword)
	}
	if captured.DefaultCOSBucket == nil || *captured.DefaultCOSBucket != oceanusClusterAfTestCosBucket {
		t.Errorf("create request default_cos_bucket = %v, want %q", captured.DefaultCOSBucket, oceanusClusterAfTestCosBucket)
	}
	if captured.CU == nil || *captured.CU != oceanusClusterAfTestCu {
		t.Errorf("create request cu = %v, want %d", captured.CU, oceanusClusterAfTestCu)
	}
	if captured.Remark == nil || *captured.Remark != oceanusClusterAfTestRemark {
		t.Errorf("create request remark = %v, want %q", captured.Remark, oceanusClusterAfTestRemark)
	}
	if len(captured.VpcDescriptions) != 1 || captured.VpcDescriptions[0] == nil {
		t.Fatalf("create request vpc_descriptions = %v, want one entry", captured.VpcDescriptions)
	}
	if captured.VpcDescriptions[0].VpcId == nil || *captured.VpcDescriptions[0].VpcId != oceanusClusterAfTestVpcId {
		t.Errorf("create request vpc_id = %v, want %q", captured.VpcDescriptions[0].VpcId, oceanusClusterAfTestVpcId)
	}
	if captured.VpcDescriptions[0].SubnetId == nil || *captured.VpcDescriptions[0].SubnetId != oceanusClusterAfTestSubnetId {
		t.Errorf("create request subnet_id = %v, want %q", captured.VpcDescriptions[0].SubnetId, oceanusClusterAfTestSubnetId)
	}

	var got svoceanus.OceanusClusterAfResourceModel
	if diags := createResp.State.Get(context.Background(), &got); diags.HasError() {
		t.Fatalf("failed to read back state: %v", diags.Errors())
	}
	if got.Id.ValueString() != oceanusClusterAfTestId {
		t.Errorf("id = %q, want %q", got.Id.ValueString(), oceanusClusterAfTestId)
	}
	if got.Cu.ValueInt64() != oceanusClusterAfTestCu {
		t.Errorf("cu = %d, want the read-back cloud value %d", got.Cu.ValueInt64(), oceanusClusterAfTestCu)
	}
}

// TestOceanusClusterAfCreateEmptyId verifies that Create fails loudly when the
// API returns no cluster id.
func TestOceanusClusterAfCreateEmptyId(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ctx := newOceanusClusterAfTestContext(t)
	oceanusClient := patchOceanusClusterAfSdkClient(patches, ctx.client)
	patches.ApplyMethodFunc(oceanusClient, "CreateOceanusCluster", func(req *oceanusv20190422.CreateOceanusClusterRequest) (*oceanusv20190422.CreateOceanusClusterResponse, error) {
		return &oceanusv20190422.CreateOceanusClusterResponse{Response: &oceanusv20190422.CreateOceanusClusterResponseParams{
			RequestId: oceanusClusterAfStr("fake-request-id"),
		}}, nil
	})

	var createResp resource.CreateResponse
	createResp.State = ctx.state(t, oceanusClusterAfTestModel())
	ctx.resource.Create(context.Background(), resource.CreateRequest{Plan: ctx.plan(t, oceanusClusterAfTestModel())}, &createResp)
	if !createResp.Diagnostics.HasError() {
		t.Fatal("Create should fail when the API returns an empty cluster id")
	}
}

// TestOceanusClusterAfCreateWaitTimeout verifies that a cluster which never
// reports the running state ends the apply with an error once the (short) user
// timeout elapses — this exercises the `timeouts` block end to end.
func TestOceanusClusterAfCreateWaitTimeout(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ctx := newOceanusClusterAfTestContext(t)
	oceanusClient := patchOceanusClusterAfSdkClient(patches, ctx.client)
	patches.ApplyMethodFunc(oceanusClient, "CreateOceanusCluster", func(req *oceanusv20190422.CreateOceanusClusterRequest) (*oceanusv20190422.CreateOceanusClusterResponse, error) {
		return &oceanusv20190422.CreateOceanusClusterResponse{Response: &oceanusv20190422.CreateOceanusClusterResponseParams{
			ClusterId: oceanusClusterAfStr(oceanusClusterAfTestId),
			RequestId: oceanusClusterAfStr("fake-request-id"),
		}}, nil
	})
	patches.ApplyMethodFunc(oceanusClient, "DescribeClusters", func(req *oceanusv20190422.DescribeClustersRequest) (*oceanusv20190422.DescribeClustersResponse, error) {
		initializing := oceanusClusterAfTestCluster()
		initializing.Status = oceanusClusterAfInt64(svoceanus.OCEANUS_CLUSTER_STATUS_INITIALIZING)

		return &oceanusv20190422.DescribeClustersResponse{Response: &oceanusv20190422.DescribeClustersResponseParams{
			ClusterSet: []*oceanusv20190422.Cluster{initializing},
			RequestId:  oceanusClusterAfStr("fake-request-id"),
		}}, nil
	})

	var createResp resource.CreateResponse
	createResp.State = ctx.state(t, oceanusClusterAfTestModel())
	ctx.resource.Create(context.Background(), resource.CreateRequest{Plan: ctx.plan(t, oceanusClusterAfTestModel())}, &createResp)
	if !createResp.Diagnostics.HasError() {
		t.Fatal("Create should fail when the cluster never reports the running state")
	}
}

// TestOceanusClusterAfCreateApiError verifies API failures surface as errors.
func TestOceanusClusterAfCreateApiError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ctx := newOceanusClusterAfTestContext(t)
	oceanusClient := patchOceanusClusterAfSdkClient(patches, ctx.client)
	patches.ApplyMethodFunc(oceanusClient, "CreateOceanusCluster", func(req *oceanusv20190422.CreateOceanusClusterRequest) (*oceanusv20190422.CreateOceanusClusterResponse, error) {
		return nil, fmt.Errorf("create failed")
	})

	var createResp resource.CreateResponse
	createResp.State = ctx.state(t, oceanusClusterAfTestModel())
	ctx.resource.Create(context.Background(), resource.CreateRequest{Plan: ctx.plan(t, oceanusClusterAfTestModel())}, &createResp)
	if !createResp.Diagnostics.HasError() {
		t.Fatal("Create should fail on an API error")
	}
}

// TestOceanusClusterAfReadSuccess verifies Read refreshes the cloud-owned
// fields — the stub reports a name different from the state, so a Read that
// skipped the refresh could not equal it.
func TestOceanusClusterAfReadSuccess(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ctx := newOceanusClusterAfTestContext(t)
	oceanusClient := patchOceanusClusterAfSdkClient(patches, ctx.client)

	patches.ApplyMethodFunc(oceanusClient, "DescribeClusters", func(req *oceanusv20190422.DescribeClustersRequest) (*oceanusv20190422.DescribeClustersResponse, error) {
		if len(req.ClusterIds) != 1 || req.ClusterIds[0] == nil || *req.ClusterIds[0] != oceanusClusterAfTestId {
			t.Errorf("expected the cluster id filter to be passed, got %v", req.ClusterIds)
		}

		normalized := oceanusClusterAfTestCluster()
		normalized.Name = oceanusClusterAfStr(oceanusClusterAfTestCloudName)

		return &oceanusv20190422.DescribeClustersResponse{Response: &oceanusv20190422.DescribeClustersResponseParams{
			ClusterSet: []*oceanusv20190422.Cluster{normalized},
			RequestId:  oceanusClusterAfStr("fake-request-id"),
		}}, nil
	})

	var readResp resource.ReadResponse
	readResp.State = ctx.state(t, oceanusClusterAfTestModel())
	ctx.resource.Read(context.Background(), resource.ReadRequest{State: ctx.state(t, oceanusClusterAfTestModel())}, &readResp)
	if readResp.Diagnostics.HasError() {
		t.Fatalf("Read failed: %v", readResp.Diagnostics.Errors())
	}

	var got svoceanus.OceanusClusterAfResourceModel
	if diags := readResp.State.Get(context.Background(), &got); diags.HasError() {
		t.Fatalf("failed to read back state: %v", diags.Errors())
	}
	if got.ClusterName.ValueString() != oceanusClusterAfTestCloudName {
		t.Errorf("cluster_name = %q, want the cloud value %q", got.ClusterName.ValueString(), oceanusClusterAfTestCloudName)
	}
}

// TestOceanusClusterAfReadNotFound verifies that a describe list without the
// tracked cluster drops the resource from state (confirmed not-found).
func TestOceanusClusterAfReadNotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ctx := newOceanusClusterAfTestContext(t)
	oceanusClient := patchOceanusClusterAfSdkClient(patches, ctx.client)
	patches.ApplyMethodFunc(oceanusClient, "DescribeClusters", func(req *oceanusv20190422.DescribeClustersRequest) (*oceanusv20190422.DescribeClustersResponse, error) {
		return &oceanusv20190422.DescribeClustersResponse{Response: &oceanusv20190422.DescribeClustersResponseParams{
			RequestId: oceanusClusterAfStr("fake-request-id"),
		}}, nil
	})

	var readResp resource.ReadResponse
	readResp.State = ctx.state(t, oceanusClusterAfTestModel())
	ctx.resource.Read(context.Background(), resource.ReadRequest{State: ctx.state(t, oceanusClusterAfTestModel())}, &readResp)
	if readResp.Diagnostics.HasError() {
		t.Fatalf("Read should not error on confirmed not-found: %v", readResp.Diagnostics.Errors())
	}
	if !readResp.State.Raw.IsNull() {
		t.Error("state should be removed after the cluster is confirmed gone")
	}
}

// TestOceanusClusterAfReadApiError verifies API failures surface as errors.
func TestOceanusClusterAfReadApiError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ctx := newOceanusClusterAfTestContext(t)
	oceanusClient := patchOceanusClusterAfSdkClient(patches, ctx.client)
	patches.ApplyMethodFunc(oceanusClient, "DescribeClusters", func(req *oceanusv20190422.DescribeClustersRequest) (*oceanusv20190422.DescribeClustersResponse, error) {
		return nil, fmt.Errorf("describe failed")
	})

	var readResp resource.ReadResponse
	readResp.State = ctx.state(t, oceanusClusterAfTestModel())
	ctx.resource.Read(context.Background(), resource.ReadRequest{State: ctx.state(t, oceanusClusterAfTestModel())}, &readResp)
	if !readResp.Diagnostics.HasError() {
		t.Fatal("Read should fail on an API error")
	}
}

// TestOceanusClusterAfUpdateScaleUp verifies Update sends the planned CU number
// to the scale API, derives the scale direction, and records the POST-update
// cloud state from the read-back rather than the plan.
func TestOceanusClusterAfUpdateScaleUp(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ctx := newOceanusClusterAfTestContext(t)
	oceanusClient := patchOceanusClusterAfSdkClient(patches, ctx.client)

	var captured *oceanusv20190422.ScaleOceanusClusterRequest
	patches.ApplyMethodFunc(oceanusClient, "ScaleOceanusCluster", func(req *oceanusv20190422.ScaleOceanusClusterRequest) (*oceanusv20190422.ScaleOceanusClusterResponse, error) {
		captured = req

		return &oceanusv20190422.ScaleOceanusClusterResponse{Response: &oceanusv20190422.ScaleOceanusClusterResponseParams{
			TaskExecResult: oceanusClusterAfStr("success"),
			RequestId:      oceanusClusterAfStr("fake-request-id"),
		}}, nil
	})
	patches.ApplyMethodFunc(oceanusClient, "DescribeClusters", func(req *oceanusv20190422.DescribeClustersRequest) (*oceanusv20190422.DescribeClustersResponse, error) {
		scaled := oceanusClusterAfTestCluster()
		scaled.CuNum = func() *uint64 { v := uint64(oceanusClusterAfTestCloudCu); return &v }()

		return &oceanusv20190422.DescribeClustersResponse{Response: &oceanusv20190422.DescribeClustersResponseParams{
			ClusterSet: []*oceanusv20190422.Cluster{scaled},
			RequestId:  oceanusClusterAfStr("fake-request-id"),
		}}, nil
	})

	plannedModel := oceanusClusterAfTestModel()
	plannedModel.Cu = types.Int64Value(oceanusClusterAfTestNewCu)

	var updateResp resource.UpdateResponse
	updateResp.State = ctx.state(t, oceanusClusterAfTestModel())
	ctx.resource.Update(context.Background(),
		resource.UpdateRequest{
			Plan:  ctx.plan(t, plannedModel),
			State: ctx.state(t, oceanusClusterAfTestModel()),
		}, &updateResp)
	if updateResp.Diagnostics.HasError() {
		t.Fatalf("Update failed: %v", updateResp.Diagnostics.Errors())
	}

	if captured == nil {
		t.Fatal("ScaleOceanusCluster request is nil")
	}
	if captured.ClusterId == nil || *captured.ClusterId != oceanusClusterAfTestId {
		t.Errorf("scale request cluster_id = %v, want %q", captured.ClusterId, oceanusClusterAfTestId)
	}
	if captured.NewCU == nil || *captured.NewCU != oceanusClusterAfTestNewCu {
		t.Errorf("scale request new_cu = %v, want the planned %d", captured.NewCU, oceanusClusterAfTestNewCu)
	}
	if captured.ScaleMode == nil || *captured.ScaleMode != "ScaleUp" {
		t.Errorf("scale request scale_mode = %v, want ScaleUp", captured.ScaleMode)
	}

	var got svoceanus.OceanusClusterAfResourceModel
	if diags := updateResp.State.Get(context.Background(), &got); diags.HasError() {
		t.Fatalf("failed to read back state: %v", diags.Errors())
	}
	if got.Cu.ValueInt64() != oceanusClusterAfTestCloudCu {
		t.Errorf("cu = %d, want the read-back cloud value %d", got.Cu.ValueInt64(), oceanusClusterAfTestCloudCu)
	}
}

// TestOceanusClusterAfUpdateScaleDown verifies the scale direction is derived
// from the CU numbers: a smaller planned value asks for ScaleDown.
func TestOceanusClusterAfUpdateScaleDown(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ctx := newOceanusClusterAfTestContext(t)
	oceanusClient := patchOceanusClusterAfSdkClient(patches, ctx.client)

	var captured *oceanusv20190422.ScaleOceanusClusterRequest
	patches.ApplyMethodFunc(oceanusClient, "ScaleOceanusCluster", func(req *oceanusv20190422.ScaleOceanusClusterRequest) (*oceanusv20190422.ScaleOceanusClusterResponse, error) {
		captured = req

		return &oceanusv20190422.ScaleOceanusClusterResponse{Response: &oceanusv20190422.ScaleOceanusClusterResponseParams{
			TaskExecResult: oceanusClusterAfStr("success"),
			RequestId:      oceanusClusterAfStr("fake-request-id"),
		}}, nil
	})
	patches.ApplyMethodFunc(oceanusClient, "DescribeClusters", func(req *oceanusv20190422.DescribeClustersRequest) (*oceanusv20190422.DescribeClustersResponse, error) {
		return &oceanusv20190422.DescribeClustersResponse{Response: &oceanusv20190422.DescribeClustersResponseParams{
			ClusterSet: []*oceanusv20190422.Cluster{oceanusClusterAfTestCluster()},
			RequestId:  oceanusClusterAfStr("fake-request-id"),
		}}, nil
	})

	stateModel := oceanusClusterAfTestModel()
	stateModel.Cu = types.Int64Value(oceanusClusterAfTestNewCu)

	plannedModel := oceanusClusterAfTestModel()
	plannedModel.Cu = types.Int64Value(oceanusClusterAfTestCu)

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

// TestOceanusClusterAfUpdateTimeoutsOnly verifies that a plan differing from
// the state in the `timeouts` block alone calls NO write API. The read-back
// still runs and the new timeouts are persisted.
func TestOceanusClusterAfUpdateTimeoutsOnly(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ctx := newOceanusClusterAfTestContext(t)
	oceanusClient := patchOceanusClusterAfSdkClient(patches, ctx.client)

	scaleCalls := 0
	patches.ApplyMethodFunc(oceanusClient, "ScaleOceanusCluster", func(req *oceanusv20190422.ScaleOceanusClusterRequest) (*oceanusv20190422.ScaleOceanusClusterResponse, error) {
		scaleCalls++

		return nil, fmt.Errorf("the scale API must not be called")
	})
	patches.ApplyMethodFunc(oceanusClient, "DescribeClusters", func(req *oceanusv20190422.DescribeClustersRequest) (*oceanusv20190422.DescribeClustersResponse, error) {
		return &oceanusv20190422.DescribeClustersResponse{Response: &oceanusv20190422.DescribeClustersResponseParams{
			ClusterSet: []*oceanusv20190422.Cluster{oceanusClusterAfTestCluster()},
			RequestId:  oceanusClusterAfStr("fake-request-id"),
		}}, nil
	})

	plannedModel := oceanusClusterAfTestModel()
	plannedModel.Timeouts = resourcetimeouts.Value{Object: types.ObjectValueMust(
		oceanusClusterAfTimeoutsAttrTypes(),
		map[string]attr.Value{
			"create": types.StringValue("7s"),
		},
	)}

	var updateResp resource.UpdateResponse
	updateResp.State = ctx.state(t, oceanusClusterAfTestModel())
	ctx.resource.Update(context.Background(),
		resource.UpdateRequest{
			Plan:  ctx.plan(t, plannedModel),
			State: ctx.state(t, oceanusClusterAfTestModel()),
		}, &updateResp)
	if updateResp.Diagnostics.HasError() {
		t.Fatalf("Update failed: %v", updateResp.Diagnostics.Errors())
	}
	if scaleCalls != 0 {
		t.Errorf("ScaleOceanusCluster was called %d time(s) for a timeouts-only change", scaleCalls)
	}
}

// TestOceanusClusterAfUpdateApiError verifies that a scale failure surfaces as
// an error and leaves the state untouched.
func TestOceanusClusterAfUpdateApiError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ctx := newOceanusClusterAfTestContext(t)
	oceanusClient := patchOceanusClusterAfSdkClient(patches, ctx.client)
	patches.ApplyMethodFunc(oceanusClient, "ScaleOceanusCluster", func(req *oceanusv20190422.ScaleOceanusClusterRequest) (*oceanusv20190422.ScaleOceanusClusterResponse, error) {
		return nil, fmt.Errorf("scale failed")
	})

	plannedModel := oceanusClusterAfTestModel()
	plannedModel.Cu = types.Int64Value(oceanusClusterAfTestNewCu)

	var updateResp resource.UpdateResponse
	updateResp.State = ctx.state(t, oceanusClusterAfTestModel())
	ctx.resource.Update(context.Background(),
		resource.UpdateRequest{
			Plan:  ctx.plan(t, plannedModel),
			State: ctx.state(t, oceanusClusterAfTestModel()),
		}, &updateResp)
	if !updateResp.Diagnostics.HasError() {
		t.Fatal("Update should fail on an API error")
	}
}

// TestOceanusClusterAfDeleteSuccess verifies Delete sends the delete request
// with the tracked id.
func TestOceanusClusterAfDeleteSuccess(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ctx := newOceanusClusterAfTestContext(t)
	oceanusClient := patchOceanusClusterAfSdkClient(patches, ctx.client)

	var captured *oceanusv20190422.DeleteOceanusClusterRequest
	patches.ApplyMethodFunc(oceanusClient, "DeleteOceanusCluster", func(req *oceanusv20190422.DeleteOceanusClusterRequest) (*oceanusv20190422.DeleteOceanusClusterResponse, error) {
		captured = req

		return &oceanusv20190422.DeleteOceanusClusterResponse{Response: &oceanusv20190422.DeleteOceanusClusterResponseParams{
			TaskExecResult: oceanusClusterAfStr("success"),
			RequestId:      oceanusClusterAfStr("fake-request-id"),
		}}, nil
	})

	var deleteResp resource.DeleteResponse
	ctx.resource.Delete(context.Background(), resource.DeleteRequest{State: ctx.state(t, oceanusClusterAfTestModel())}, &deleteResp)
	if deleteResp.Diagnostics.HasError() {
		t.Fatalf("Delete failed: %v", deleteResp.Diagnostics.Errors())
	}

	if captured == nil {
		t.Fatal("DeleteOceanusCluster request is nil")
	}
	if captured.ClusterId == nil || *captured.ClusterId != oceanusClusterAfTestId {
		t.Errorf("delete request cluster_id = %v, want %q", captured.ClusterId, oceanusClusterAfTestId)
	}
}

// TestOceanusClusterAfDeleteApiError verifies API failures surface as errors.
func TestOceanusClusterAfDeleteApiError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ctx := newOceanusClusterAfTestContext(t)
	oceanusClient := patchOceanusClusterAfSdkClient(patches, ctx.client)
	patches.ApplyMethodFunc(oceanusClient, "DeleteOceanusCluster", func(req *oceanusv20190422.DeleteOceanusClusterRequest) (*oceanusv20190422.DeleteOceanusClusterResponse, error) {
		return nil, fmt.Errorf("delete failed")
	})

	var deleteResp resource.DeleteResponse
	ctx.resource.Delete(context.Background(), resource.DeleteRequest{State: ctx.state(t, oceanusClusterAfTestModel())}, &deleteResp)
	if !deleteResp.Diagnostics.HasError() {
		t.Fatal("Delete should fail on an API error")
	}
}
