package oceanus

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

	oceanus "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/oceanus/v20190422"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/sharedmeta"
)

const (
	oceanusClusterAfTestId        = "cluster-12345"
	oceanusClusterAfTestName      = "tf-example-cluster"
	oceanusClusterAfTestBucket    = "tf-test-bucket"
	oceanusClusterAfTestVpcId     = "vpc-12345"
	oceanusClusterAfTestSubnetId  = "subnet-12345"
	oceanusClusterAfTestPassword  = "tf-example-password"
	oceanusClusterAfTestStatusRun = "2"
	oceanusClusterAfTestScaledCU  = 19
	oceanusClusterAfTestBaseCU    = 12
)

// oceanusClusterAfStr is the *string helper for stub data.
func oceanusClusterAfStr(v string) *string { return &v }

// oceanusClusterAfInt64 is the *int64 helper for stub data.
func oceanusClusterAfInt64(v int64) *int64 { return &v }

// oceanusClusterAfUint64 is the *uint64 helper for stub data.
func oceanusClusterAfUint64(v uint64) *uint64 { return &v }

// oceanusClusterAfTestItem builds the DescribeClusters response item.
func oceanusClusterAfTestItem() *oceanus.Cluster {
	return &oceanus.Cluster{
		ClusterId:        oceanusClusterAfStr(oceanusClusterAfTestId),
		Name:             oceanusClusterAfStr(oceanusClusterAfTestName),
		Zone:             oceanusClusterAfStr("ap-guangzhou-1"),
		Status:           oceanusClusterAfInt64(2),
		StatusDesc:       oceanusClusterAfStr("Running"),
		CreateTime:       oceanusClusterAfStr("2024-01-01 00:00:00"),
		UpdateTime:       oceanusClusterAfStr("2024-01-01 00:00:00"),
		DefaultCOSBucket: oceanusClusterAfStr(oceanusClusterAfTestBucket),
		CuNum:            oceanusClusterAfUint64(oceanusClusterAfTestBaseCU),
	}
}

// oceanusClusterAfTestModel builds the terraform model used for plan/state.
// Every unset attribute carries an explicit typed null value — the raw zero
// value (types.List{} etc.) breaks tfsdk Plan/State reflection. Cases
// exercising a wait path pass a second-scale short user value (5s) so the
// code path is verified without hard-waiting the minute-scale default.
func oceanusClusterAfTestModel() OceanusClusterAfResourceModel {
	return OceanusClusterAfResourceModel{
		Id:                          types.StringValue(oceanusClusterAfTestId),
		ClusterName:                 types.StringValue(oceanusClusterAfTestName),
		RegionId:                    types.Int64Value(1),
		ZoneId:                      types.Int64Value(100001),
		LoginPassword:               types.StringValue(oceanusClusterAfTestPassword),
		VpcDescriptions:             oceanusClusterAfTestVpcList(),
		DefaultCOSBucket:            types.StringValue(oceanusClusterAfTestBucket),
		CU:                          types.Int64Value(oceanusClusterAfTestBaseCU),
		Remark:                      types.StringNull(),
		Period:                      types.Int64Null(),
		InstanceChargeType:          types.StringNull(),
		ClusterType:                 types.StringNull(),
		RenewFlag:                   types.StringNull(),
		FlinkWebUINetworkAccessType: types.StringNull(),
		SlaveVpcDescriptions:        types.ListNull(oceanusClusterAfSlaveAttrType()),
		CUMemory:                    types.Int64Null(),
		Status:                      types.StringNull(),
		StatusDesc:                  types.StringNull(),
		CreateTime:                  types.StringNull(),
		UpdateTime:                  types.StringNull(),
		Zone:                        types.StringNull(),
		ClusterReadId:               types.StringNull(),
		Timeouts: resourcetimeouts.Value{Object: types.ObjectValueMust(
			oceanusClusterAfTimeoutsAttrTypes(),
			map[string]attr.Value{
				"create": types.StringValue("5s"),
				"update": types.StringValue("5s"),
				"delete": types.StringValue("5s"),
			},
		)},
	}
}

// oceanusClusterAfTestVpcList builds the required vpc_descriptions list.
func oceanusClusterAfTestVpcList() types.List {
	return types.ListValueMust(
		oceanusClusterAfVpcAttrType(),
		[]attr.Value{types.ObjectValueMust(
			oceanusClusterAfVpcAttrTypes(),
			map[string]attr.Value{
				"vpc_id":    types.StringValue(oceanusClusterAfTestVpcId),
				"subnet_id": types.StringValue(oceanusClusterAfTestSubnetId),
			},
		)},
	)
}

// oceanusClusterAfVpcAttrType is the element type of vpc_descriptions.
func oceanusClusterAfVpcAttrType() attr.Type {
	return types.ObjectType{AttrTypes: oceanusClusterAfVpcAttrTypes()}
}

// oceanusClusterAfVpcAttrTypes mirrors the vpc_descriptions element schema.
func oceanusClusterAfVpcAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"vpc_id":    types.StringType,
		"subnet_id": types.StringType,
	}
}

// oceanusClusterAfSlaveAttrType is the element type of slave_vpc_descriptions.
func oceanusClusterAfSlaveAttrType() attr.Type {
	return types.ObjectType{AttrTypes: oceanusClusterAfVpcAttrTypes()}
}

// oceanusClusterAfTimeoutsAttrTypes mirrors the attribute types of the
// timeouts block declared by the resource schema.
func oceanusClusterAfTimeoutsAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"create": types.StringType,
		"update": types.StringType,
		"delete": types.StringType,
	}
}

// patchOceanusClusterAfSdkClient stubs the client accessor and returns the
// mock SDK client whose methods the test patches next (two-level mode).
func patchOceanusClusterAfSdkClient(patches *gomonkey.Patches, client *connectivity.TencentCloudClient) *oceanus.Client {
	oceanusClient := &oceanus.Client{}
	patches.ApplyMethodReturn(client, "UseOceanusClient", oceanusClient)

	return oceanusClient
}

// stubOceanusClusterAfDescribe returns a describe answer carrying the given
// items for the first call.
func stubOceanusClusterAfDescribe(patches *gomonkey.Patches, oceanusClient *oceanus.Client, items []*oceanus.Cluster) {
	patches.ApplyMethodFunc(oceanusClient, "DescribeClusters", func(req *oceanus.DescribeClustersRequest) (*oceanus.DescribeClustersResponse, error) {
		return &oceanus.DescribeClustersResponse{Response: &oceanus.DescribeClustersResponseParams{
			ClusterSet: items,
			TotalCount: oceanusClusterAfInt64(int64(len(items))),
			RequestId:  oceanusClusterAfStr("fake-request-id"),
		}}, nil
	})
}

// stubOceanusClusterAfDescribeNeverVisible stubs DescribeClusters to answer
// with an empty list — used to exercise wait timeout paths.
func stubOceanusClusterAfDescribeNeverVisible(patches *gomonkey.Patches, oceanusClient *oceanus.Client) {
	patches.ApplyMethodFunc(oceanusClient, "DescribeClusters", func(req *oceanus.DescribeClustersRequest) (*oceanus.DescribeClustersResponse, error) {
		return &oceanus.DescribeClustersResponse{Response: &oceanus.DescribeClustersResponseParams{
			ClusterSet: []*oceanus.Cluster{},
			RequestId:  oceanusClusterAfStr("fake-request-id"),
		}}, nil
	})
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

	res := NewOceanusClusterAfResource()
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

func (c *oceanusClusterAfTestContext) plan(t *testing.T, m OceanusClusterAfResourceModel) tfsdk.Plan {
	t.Helper()

	plan := tfsdk.Plan{Schema: c.schema}
	if diags := plan.Set(context.Background(), m); diags.HasError() {
		t.Fatalf("failed to build plan: %v", diags.Errors())
	}

	return plan
}

func (c *oceanusClusterAfTestContext) state(t *testing.T, m OceanusClusterAfResourceModel) tfsdk.State {
	t.Helper()

	state := tfsdk.State{Schema: c.schema}
	if diags := state.Set(context.Background(), m); diags.HasError() {
		t.Fatalf("failed to build state: %v", diags.Errors())
	}

	return state
}

// assertOceanusClusterAfModel asserts the read-back fields on a fresh model.
func assertOceanusClusterAfModel(t *testing.T, m OceanusClusterAfResourceModel) {
	t.Helper()

	if m.Id.ValueString() != oceanusClusterAfTestId {
		t.Errorf("id = %q, want %q", m.Id.ValueString(), oceanusClusterAfTestId)
	}
	if m.ClusterReadId.ValueString() != oceanusClusterAfTestId {
		t.Errorf("cluster_read_id = %q, want %q", m.ClusterReadId.ValueString(), oceanusClusterAfTestId)
	}
	if m.ClusterName.ValueString() != oceanusClusterAfTestName {
		t.Errorf("cluster_name = %q, want %q", m.ClusterName.ValueString(), oceanusClusterAfTestName)
	}
	if m.Status.ValueString() != oceanusClusterAfTestStatusRun {
		t.Errorf("status = %q, want %q", m.Status.ValueString(), oceanusClusterAfTestStatusRun)
	}
	if m.Zone.ValueString() != "ap-guangzhou-1" {
		t.Errorf("zone = %q, want ap-guangzhou-1", m.Zone.ValueString())
	}
	if m.CU.ValueInt64() != oceanusClusterAfTestBaseCU {
		t.Errorf("cu = %d, want %d", m.CU.ValueInt64(), oceanusClusterAfTestBaseCU)
	}
}

// TestOceanusClusterAfCreateSuccess verifies that Create sends all request
// fields, tracks the returned cluster id and refreshes the state.
func TestOceanusClusterAfCreateSuccess(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ctx := newOceanusClusterAfTestContext(t)
	oceanusClient := patchOceanusClusterAfSdkClient(patches, ctx.client)

	var captured *oceanus.CreateOceanusClusterRequest
	patches.ApplyMethodFunc(oceanusClient, "CreateOceanusCluster", func(req *oceanus.CreateOceanusClusterRequest) (*oceanus.CreateOceanusClusterResponse, error) {
		captured = req

		return &oceanus.CreateOceanusClusterResponse{Response: &oceanus.CreateOceanusClusterResponseParams{
			ClusterId: oceanusClusterAfStr(oceanusClusterAfTestId),
			RequestId: oceanusClusterAfStr("fake-request-id"),
		}}, nil
	})
	stubOceanusClusterAfDescribe(patches, oceanusClient, []*oceanus.Cluster{oceanusClusterAfTestItem()})

	model := oceanusClusterAfTestModel()
	var createResp resource.CreateResponse
	createResp.State = ctx.state(t, model)
	ctx.resource.Create(context.Background(), resource.CreateRequest{Plan: ctx.plan(t, model)}, &createResp)
	if createResp.Diagnostics.HasError() {
		t.Fatalf("Create failed: %v", createResp.Diagnostics.Errors())
	}

	if captured == nil {
		t.Fatal("CreateOceanusCluster request is nil")
	}
	if captured.ClusterName == nil || *captured.ClusterName != oceanusClusterAfTestName {
		t.Errorf("create request cluster_name = %v, want %q", captured.ClusterName, oceanusClusterAfTestName)
	}
	if captured.RegionId == nil || *captured.RegionId != 1 {
		t.Errorf("create request region_id = %v, want 1", captured.RegionId)
	}
	if captured.ZoneId == nil || *captured.ZoneId != 100001 {
		t.Errorf("create request zone_id = %v, want 100001", captured.ZoneId)
	}
	if captured.LoginPassword == nil || *captured.LoginPassword != oceanusClusterAfTestPassword {
		t.Errorf("create request login_password = %v, want %q", captured.LoginPassword, oceanusClusterAfTestPassword)
	}
	if captured.DefaultCOSBucket == nil || *captured.DefaultCOSBucket != oceanusClusterAfTestBucket {
		t.Errorf("create request default_cos_bucket = %v, want %q", captured.DefaultCOSBucket, oceanusClusterAfTestBucket)
	}
	if captured.CU == nil || *captured.CU != oceanusClusterAfTestBaseCU {
		t.Errorf("create request cu = %v, want %d", captured.CU, oceanusClusterAfTestBaseCU)
	}
	if len(captured.VpcDescriptions) != 1 {
		t.Fatalf("create request vpc_descriptions length = %d, want 1", len(captured.VpcDescriptions))
	}
	if captured.VpcDescriptions[0].VpcId == nil || *captured.VpcDescriptions[0].VpcId != oceanusClusterAfTestVpcId {
		t.Errorf("create request vpc_id = %v, want %q", captured.VpcDescriptions[0].VpcId, oceanusClusterAfTestVpcId)
	}
	if captured.VpcDescriptions[0].SubnetId == nil || *captured.VpcDescriptions[0].SubnetId != oceanusClusterAfTestSubnetId {
		t.Errorf("create request subnet_id = %v, want %q", captured.VpcDescriptions[0].SubnetId, oceanusClusterAfTestSubnetId)
	}

	var got OceanusClusterAfResourceModel
	if diags := createResp.State.Get(context.Background(), &got); diags.HasError() {
		t.Fatalf("failed to read back state: %v", diags.Errors())
	}
	assertOceanusClusterAfModel(t, got)
}

// TestOceanusClusterAfCreateEmptyId verifies that Create fails loudly when the
// API returns no cluster id.
func TestOceanusClusterAfCreateEmptyId(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ctx := newOceanusClusterAfTestContext(t)
	oceanusClient := patchOceanusClusterAfSdkClient(patches, ctx.client)
	patches.ApplyMethodFunc(oceanusClient, "CreateOceanusCluster", func(req *oceanus.CreateOceanusClusterRequest) (*oceanus.CreateOceanusClusterResponse, error) {
		return &oceanus.CreateOceanusClusterResponse{Response: &oceanus.CreateOceanusClusterResponseParams{
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

// TestOceanusClusterAfCreateWaitTimeout verifies that a cluster that never
// becomes visible ends the apply with an error once the (short) user timeout
// elapses — this exercises the `timeouts` block end to end.
func TestOceanusClusterAfCreateWaitTimeout(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ctx := newOceanusClusterAfTestContext(t)
	oceanusClient := patchOceanusClusterAfSdkClient(patches, ctx.client)
	patches.ApplyMethodFunc(oceanusClient, "CreateOceanusCluster", func(req *oceanus.CreateOceanusClusterRequest) (*oceanus.CreateOceanusClusterResponse, error) {
		return &oceanus.CreateOceanusClusterResponse{Response: &oceanus.CreateOceanusClusterResponseParams{
			ClusterId: oceanusClusterAfStr(oceanusClusterAfTestId),
			RequestId: oceanusClusterAfStr("fake-request-id"),
		}}, nil
	})
	stubOceanusClusterAfDescribeNeverVisible(patches, oceanusClient)

	var createResp resource.CreateResponse
	createResp.State = ctx.state(t, oceanusClusterAfTestModel())
	ctx.resource.Create(context.Background(), resource.CreateRequest{Plan: ctx.plan(t, oceanusClusterAfTestModel())}, &createResp)
	if !createResp.Diagnostics.HasError() {
		t.Fatal("Create should fail when the cluster never becomes visible")
	}
}

// TestOceanusClusterAfCreateApiError verifies API failures surface as errors.
func TestOceanusClusterAfCreateApiError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ctx := newOceanusClusterAfTestContext(t)
	oceanusClient := patchOceanusClusterAfSdkClient(patches, ctx.client)
	patches.ApplyMethodFunc(oceanusClient, "CreateOceanusCluster", func(req *oceanus.CreateOceanusClusterRequest) (*oceanus.CreateOceanusClusterResponse, error) {
		return nil, fmt.Errorf("create failed")
	})

	var createResp resource.CreateResponse
	createResp.State = ctx.state(t, oceanusClusterAfTestModel())
	ctx.resource.Create(context.Background(), resource.CreateRequest{Plan: ctx.plan(t, oceanusClusterAfTestModel())}, &createResp)
	if !createResp.Diagnostics.HasError() {
		t.Fatal("Create should fail on an API error")
	}
}

// TestOceanusClusterAfReadSuccess verifies Read refreshes all read-back fields.
func TestOceanusClusterAfReadSuccess(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ctx := newOceanusClusterAfTestContext(t)
	oceanusClient := patchOceanusClusterAfSdkClient(patches, ctx.client)

	var gotClusterId string
	patches.ApplyMethodFunc(oceanusClient, "DescribeClusters", func(req *oceanus.DescribeClustersRequest) (*oceanus.DescribeClustersResponse, error) {
		if len(req.ClusterIds) != 1 || req.ClusterIds[0] == nil {
			t.Errorf("expected the cluster id filter to be passed")
		} else {
			gotClusterId = *req.ClusterIds[0]
		}

		return &oceanus.DescribeClustersResponse{Response: &oceanus.DescribeClustersResponseParams{
			ClusterSet: []*oceanus.Cluster{oceanusClusterAfTestItem()},
			RequestId:  oceanusClusterAfStr("fake-request-id"),
		}}, nil
	})

	var readResp resource.ReadResponse
	readResp.State = ctx.state(t, oceanusClusterAfTestModel())
	ctx.resource.Read(context.Background(), resource.ReadRequest{State: ctx.state(t, oceanusClusterAfTestModel())}, &readResp)
	if readResp.Diagnostics.HasError() {
		t.Fatalf("Read failed: %v", readResp.Diagnostics.Errors())
	}

	if gotClusterId != oceanusClusterAfTestId {
		t.Errorf("DescribeClusters called with %q, want %q", gotClusterId, oceanusClusterAfTestId)
	}

	var got OceanusClusterAfResourceModel
	if diags := readResp.State.Get(context.Background(), &got); diags.HasError() {
		t.Fatalf("failed to read back state: %v", diags.Errors())
	}
	assertOceanusClusterAfModel(t, got)
}

// TestOceanusClusterAfReadNotFound verifies that a describe list without the
// tracked cluster drops the resource from state (confirmed not-found).
func TestOceanusClusterAfReadNotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ctx := newOceanusClusterAfTestContext(t)
	oceanusClient := patchOceanusClusterAfSdkClient(patches, ctx.client)
	stubOceanusClusterAfDescribe(patches, oceanusClient, []*oceanus.Cluster{})

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

// TestOceanusClusterAfUpdateSuccess verifies Update sends the PLANNED CU value
// to the scale request, waits for the resize, and records the POST-update
// cloud state from the read-back — not the plan.
func TestOceanusClusterAfUpdateSuccess(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ctx := newOceanusClusterAfTestContext(t)
	oceanusClient := patchOceanusClusterAfSdkClient(patches, ctx.client)

	var captured *oceanus.ScaleOceanusClusterRequest
	patches.ApplyMethodFunc(oceanusClient, "ScaleOceanusCluster", func(req *oceanus.ScaleOceanusClusterRequest) (*oceanus.ScaleOceanusClusterResponse, error) {
		captured = req

		return &oceanus.ScaleOceanusClusterResponse{Response: &oceanus.ScaleOceanusClusterResponseParams{
			TaskExecResult: oceanusClusterAfStr("success"),
			RequestId:      oceanusClusterAfStr("fake-request-id"),
		}}, nil
	})
	patches.ApplyMethodFunc(oceanusClient, "DescribeClusters", func(req *oceanus.DescribeClustersRequest) (*oceanus.DescribeClustersResponse, error) {
		scaled := oceanusClusterAfTestItem()
		scaled.CuNum = oceanusClusterAfUint64(oceanusClusterAfTestScaledCU)

		return &oceanus.DescribeClustersResponse{Response: &oceanus.DescribeClustersResponseParams{
			ClusterSet: []*oceanus.Cluster{scaled},
			RequestId:  oceanusClusterAfStr("fake-request-id"),
		}}, nil
	})

	plannedModel := oceanusClusterAfTestModel()
	plannedModel.CU = types.Int64Value(oceanusClusterAfTestScaledCU)

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
		t.Errorf("update request cluster_id = %v, want %q", captured.ClusterId, oceanusClusterAfTestId)
	}
	if captured.NewCU == nil || *captured.NewCU != oceanusClusterAfTestScaledCU {
		t.Errorf("update request new_cu = %v, want the planned %d", captured.NewCU, oceanusClusterAfTestScaledCU)
	}

	var got OceanusClusterAfResourceModel
	if diags := updateResp.State.Get(context.Background(), &got); diags.HasError() {
		t.Fatalf("failed to read back state: %v", diags.Errors())
	}
	if got.Id.ValueString() != oceanusClusterAfTestId {
		t.Errorf("id = %q, want the unchanged %q", got.Id.ValueString(), oceanusClusterAfTestId)
	}
	if got.CU.ValueInt64() != oceanusClusterAfTestScaledCU {
		t.Errorf("cu = %d, want the read-back cloud value %d", got.CU.ValueInt64(), oceanusClusterAfTestScaledCU)
	}
}

// TestOceanusClusterAfUpdateTimeoutsOnly verifies that a plan differing from
// the state in the `timeouts` block alone calls NO write API: no attribute the
// scale API covers changed, so the call would be a pointless no-change
// request. The read-back still runs and the new timeouts are persisted.
func TestOceanusClusterAfUpdateTimeoutsOnly(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ctx := newOceanusClusterAfTestContext(t)
	oceanusClient := patchOceanusClusterAfSdkClient(patches, ctx.client)

	scaleCalls := 0
	patches.ApplyMethodFunc(oceanusClient, "ScaleOceanusCluster", func(req *oceanus.ScaleOceanusClusterRequest) (*oceanus.ScaleOceanusClusterResponse, error) {
		scaleCalls++

		return nil, fmt.Errorf("the scale API must not be called")
	})
	stubOceanusClusterAfDescribe(patches, oceanusClient, []*oceanus.Cluster{oceanusClusterAfTestItem()})

	plannedModel := oceanusClusterAfTestModel()
	plannedModel.Timeouts = resourcetimeouts.Value{Object: types.ObjectValueMust(
		oceanusClusterAfTimeoutsAttrTypes(),
		map[string]attr.Value{
			"create": types.StringValue("5s"),
			"update": types.StringValue("7s"),
			"delete": types.StringValue("5s"),
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

	var got OceanusClusterAfResourceModel
	if diags := updateResp.State.Get(context.Background(), &got); diags.HasError() {
		t.Fatalf("failed to read back state: %v", diags.Errors())
	}
	assertOceanusClusterAfModel(t, got)
}

// TestOceanusClusterAfDeleteSuccess verifies Delete sends the delete request
// with the tracked id and then waits until the cluster is gone.
func TestOceanusClusterAfDeleteSuccess(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ctx := newOceanusClusterAfTestContext(t)
	oceanusClient := patchOceanusClusterAfSdkClient(patches, ctx.client)

	var captured *oceanus.DeleteOceanusClusterRequest
	patches.ApplyMethodFunc(oceanusClient, "DeleteOceanusCluster", func(req *oceanus.DeleteOceanusClusterRequest) (*oceanus.DeleteOceanusClusterResponse, error) {
		captured = req

		return &oceanus.DeleteOceanusClusterResponse{Response: &oceanus.DeleteOceanusClusterResponseParams{
			TaskExecResult: oceanusClusterAfStr("success"),
			RequestId:      oceanusClusterAfStr("fake-request-id"),
		}}, nil
	})
	// The instance-status wait polls the describe API; an empty answer is the
	// terminal "destroyed" state, reported on the first call.
	stubOceanusClusterAfDescribe(patches, oceanusClient, []*oceanus.Cluster{})

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

// TestOceanusClusterAfDeleteWaitTimeout verifies that a cluster that never
// disappears ends the destroy with an error once the (short) user timeout
// elapses.
func TestOceanusClusterAfDeleteWaitTimeout(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ctx := newOceanusClusterAfTestContext(t)
	oceanusClient := patchOceanusClusterAfSdkClient(patches, ctx.client)

	patches.ApplyMethodFunc(oceanusClient, "DeleteOceanusCluster", func(req *oceanus.DeleteOceanusClusterRequest) (*oceanus.DeleteOceanusClusterResponse, error) {
		return &oceanus.DeleteOceanusClusterResponse{Response: &oceanus.DeleteOceanusClusterResponseParams{
			TaskExecResult: oceanusClusterAfStr("success"),
			RequestId:      oceanusClusterAfStr("fake-request-id"),
		}}, nil
	})
	stubOceanusClusterAfDescribe(patches, oceanusClient, []*oceanus.Cluster{oceanusClusterAfTestItem()})

	var deleteResp resource.DeleteResponse
	ctx.resource.Delete(context.Background(), resource.DeleteRequest{State: ctx.state(t, oceanusClusterAfTestModel())}, &deleteResp)
	if !deleteResp.Diagnostics.HasError() {
		t.Fatal("Delete should fail when the cluster never disappears")
	}
}
