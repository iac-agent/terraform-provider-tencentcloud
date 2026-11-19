package ags_test

import (
	"context"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	agsv20250920 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ags/v20250920"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/sharedmeta"
	svcags "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/ags"
)

const (
	agsTestGroupId   = "rg-001"
	agsTestGroupName = "tf-example"
)

func agsUint64(value uint64) *uint64 { return &value }

// agsGpuResourceGroupTestContext wires the GPU resource group resource to a
// patched SDK client.
type agsGpuResourceGroupTestContext struct {
	resource resource.Resource
	client   *agsv20250920.Client
	schema   schema.Schema
}

func newAgsGpuResourceGroupTestContext(t *testing.T, patches *gomonkey.Patches) *agsGpuResourceGroupTestContext {
	t.Helper()

	connectivityClient := &connectivity.TencentCloudClient{}
	agsClient := &agsv20250920.Client{}
	patches.ApplyMethodReturn(connectivityClient, "UseAgsClient", agsClient)

	res := svcags.NewAgsGPUResourceGroupResource()
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

	return &agsGpuResourceGroupTestContext{resource: res, client: agsClient, schema: schemaResp.Schema}
}

func (c *agsGpuResourceGroupTestContext) state(t *testing.T, model svcags.AgsGPUResourceGroupResourceModel) tfsdk.State {
	t.Helper()

	state := tfsdk.State{Schema: c.schema, Raw: tftypes.NewValue(c.schema.Type().TerraformType(context.Background()), nil)}
	if diags := state.Set(context.Background(), model); diags.HasError() {
		t.Fatalf("failed to build state: %v", diags.Errors())
	}
	return state
}

func (c *agsGpuResourceGroupTestContext) plan(t *testing.T, model svcags.AgsGPUResourceGroupResourceModel) tfsdk.Plan {
	t.Helper()

	plan := tfsdk.Plan{Schema: c.schema, Raw: tftypes.NewValue(c.schema.Type().TerraformType(context.Background()), nil)}
	if diags := plan.Set(context.Background(), model); diags.HasError() {
		t.Fatalf("failed to build plan: %v", diags.Errors())
	}
	return plan
}

func sampleGpuResourceGroup() *agsv20250920.GPUResourceGroup {
	return &agsv20250920.GPUResourceGroup{
		ResourceGroupId: agsString(agsTestGroupId),
		Name:            agsString(agsTestGroupName),
		Spec: &agsv20250920.GPUResourceGroupSpec{
			GPUType:      agsString("5090d"),
			GPUTypeAlias: agsString("RTX 5090D"),
			GPUMemoryGiB: agsUint64(32),
			CPU:          agsUint64(16),
			Memory:       agsUint64(64),
			GPUNum:       agsUint64(1),
		},
		Statistics: &agsv20250920.GPUResourceGroupStatistics{
			TotalCount:     agsInt64(4),
			AvailableCount: agsInt64(2),
			AllocatedCount: agsInt64(2),
		},
	}
}

func describeGpuResourceGroupsResponse(groups ...*agsv20250920.GPUResourceGroup) *agsv20250920.DescribeGPUResourceGroupsResponse {
	return &agsv20250920.DescribeGPUResourceGroupsResponse{
		Response: &agsv20250920.DescribeGPUResourceGroupsResponseParams{
			ResourceGroups: groups,
			TotalCount:     agsInt64(int64(len(groups))),
		},
	}
}

func assertGpuGroupString(t *testing.T, value types.String, want string) {
	t.Helper()
	if value.ValueString() != want {
		t.Fatalf("value is %q, want %q", value.ValueString(), want)
	}
}

func assertGpuGroupInt64(t *testing.T, value types.Int64, want int64) {
	t.Helper()
	if value.ValueInt64() != want {
		t.Fatalf("value is %d, want %d", value.ValueInt64(), want)
	}
}

func gpuGroupStateModel() svcags.AgsGPUResourceGroupResourceModel {
	return svcags.AgsGPUResourceGroupResourceModel{
		Id:              types.StringValue(agsTestGroupId),
		ResourceGroupId: types.StringValue(agsTestGroupId),
		Name:            types.StringValue(agsTestGroupName),
	}
}

func TestAgsGpuResourceGroupSchema(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	testContext := newAgsGpuResourceGroupTestContext(t, patches)

	if _, ok := testContext.resource.(resource.ResourceWithConfigure); !ok {
		t.Fatalf("resource does not implement resource.ResourceWithConfigure")
	}
	if _, ok := testContext.resource.(resource.ResourceWithImportState); !ok {
		t.Fatalf("resource does not implement resource.ResourceWithImportState")
	}
	for _, name := range []string{"id", "resource_group_id", "name", "gpu_type", "gpu_type_alias", "gpu_memory_gib", "cpu", "memory", "gpu_num", "total_count", "available_count", "allocated_count"} {
		if _, ok := testContext.schema.Attributes[name]; !ok {
			t.Fatalf("schema attribute %q is missing", name)
		}
	}
}

func TestAgsGpuResourceGroupCreateRejected(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	testContext := newAgsGpuResourceGroupTestContext(t, patches)

	createResp := &resource.CreateResponse{State: tfsdk.State{Schema: testContext.schema}}
	testContext.resource.Create(
		context.Background(),
		resource.CreateRequest{Plan: testContext.plan(t, gpuGroupStateModel())},
		createResp,
	)
	if !createResp.Diagnostics.HasError() {
		t.Fatalf("expected create to be rejected for an import only resource")
	}
}

func TestAgsGpuResourceGroupDeleteRejected(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	testContext := newAgsGpuResourceGroupTestContext(t, patches)

	deleteResp := &resource.DeleteResponse{}
	testContext.resource.Delete(
		context.Background(),
		resource.DeleteRequest{State: testContext.state(t, gpuGroupStateModel())},
		deleteResp,
	)
	if !deleteResp.Diagnostics.HasError() {
		t.Fatalf("expected delete to be rejected for an import only resource")
	}
}

func TestAgsGpuResourceGroupRead(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	testContext := newAgsGpuResourceGroupTestContext(t, patches)

	patches.ApplyMethodFunc(
		testContext.client,
		"DescribeGPUResourceGroupsWithContext",
		func(_ context.Context, request *agsv20250920.DescribeGPUResourceGroupsRequest) (*agsv20250920.DescribeGPUResourceGroupsResponse, error) {
			if len(request.Filters) != 1 || request.Filters[0].Name == nil || *request.Filters[0].Name != "resource-group-id" {
				t.Fatalf("expected a resource-group-id filter, got %v", request.Filters)
			}
			return describeGpuResourceGroupsResponse(sampleGpuResourceGroup()), nil
		},
	)

	readResp := &resource.ReadResponse{State: tfsdk.State{Schema: testContext.schema}}
	testContext.resource.Read(
		context.Background(),
		resource.ReadRequest{State: testContext.state(t, gpuGroupStateModel())},
		readResp,
	)
	assertNoErrorDiagnostics(t, readResp.Diagnostics)

	var state svcags.AgsGPUResourceGroupResourceModel
	if diags := readResp.State.Get(context.Background(), &state); diags.HasError() {
		t.Fatalf("failed to read refreshed state: %v", diags.Errors())
	}

	assertGpuGroupString(t, state.Id, agsTestGroupId)
	assertGpuGroupString(t, state.ResourceGroupId, agsTestGroupId)
	assertGpuGroupString(t, state.Name, agsTestGroupName)
	assertGpuGroupString(t, state.GPUType, "5090d")
	assertGpuGroupString(t, state.GPUTypeAlias, "RTX 5090D")
	assertGpuGroupInt64(t, state.GPUMemoryGiB, 32)
	assertGpuGroupInt64(t, state.CPU, 16)
	assertGpuGroupInt64(t, state.Memory, 64)
	assertGpuGroupInt64(t, state.GPUNum, 1)
	assertGpuGroupInt64(t, state.TotalCount, 4)
	assertGpuGroupInt64(t, state.AvailableCount, 2)
	assertGpuGroupInt64(t, state.AllocatedCount, 2)
}

func TestAgsGpuResourceGroupReadNotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	testContext := newAgsGpuResourceGroupTestContext(t, patches)

	patches.ApplyMethodFunc(
		testContext.client,
		"DescribeGPUResourceGroupsWithContext",
		func(_ context.Context, _ *agsv20250920.DescribeGPUResourceGroupsRequest) (*agsv20250920.DescribeGPUResourceGroupsResponse, error) {
			return describeGpuResourceGroupsResponse(), nil
		},
	)

	readResp := &resource.ReadResponse{State: tfsdk.State{Schema: testContext.schema}}
	testContext.resource.Read(
		context.Background(),
		resource.ReadRequest{State: testContext.state(t, gpuGroupStateModel())},
		readResp,
	)
	assertNoErrorDiagnostics(t, readResp.Diagnostics)

	if !readResp.State.Raw.IsNull() {
		t.Fatalf("expected state to be removed for a missing GPU resource group")
	}
}

func TestAgsGpuResourceGroupReadApiError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	testContext := newAgsGpuResourceGroupTestContext(t, patches)

	patches.ApplyMethodFunc(
		testContext.client,
		"DescribeGPUResourceGroupsWithContext",
		func(ctx context.Context, _ *agsv20250920.DescribeGPUResourceGroupsRequest) (*agsv20250920.DescribeGPUResourceGroupsResponse, error) {
			return nil, &mockSdkError{message: "api failure"}
		},
	)

	readResp := &resource.ReadResponse{State: testContext.state(t, gpuGroupStateModel())}
	testContext.resource.Read(
		context.Background(),
		resource.ReadRequest{State: testContext.state(t, gpuGroupStateModel())},
		readResp,
	)
	if !readResp.Diagnostics.HasError() {
		t.Fatalf("expected an error diagnostic for a failed read")
	}
	if readResp.State.Raw.IsNull() {
		t.Fatalf("state must be kept when the read request fails")
	}
}

func TestAgsGpuResourceGroupUpdate(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	testContext := newAgsGpuResourceGroupTestContext(t, patches)

	var modifiedName string
	patches.ApplyMethodFunc(
		testContext.client,
		"ModifyGPUResourceGroupWithContext",
		func(_ context.Context, request *agsv20250920.ModifyGPUResourceGroupRequest) (*agsv20250920.ModifyGPUResourceGroupResponse, error) {
			if request.ResourceGroupId == nil || *request.ResourceGroupId != agsTestGroupId {
				t.Fatalf("resource group id is %v, want %q", request.ResourceGroupId, agsTestGroupId)
			}
			modifiedName = *request.Name
			return &agsv20250920.ModifyGPUResourceGroupResponse{Response: &agsv20250920.ModifyGPUResourceGroupResponseParams{}}, nil
		},
	)
	patches.ApplyMethodFunc(
		testContext.client,
		"DescribeGPUResourceGroupsWithContext",
		func(_ context.Context, _ *agsv20250920.DescribeGPUResourceGroupsRequest) (*agsv20250920.DescribeGPUResourceGroupsResponse, error) {
			group := sampleGpuResourceGroup()
			group.Name = agsString("renamed")
			return describeGpuResourceGroupsResponse(group), nil
		},
	)

	planModel := gpuGroupStateModel()
	planModel.Name = types.StringValue("renamed")

	updateResp := &resource.UpdateResponse{State: tfsdk.State{Schema: testContext.schema}}
	testContext.resource.Update(
		context.Background(),
		resource.UpdateRequest{
			Plan:   testContext.plan(t, planModel),
			State:  testContext.state(t, gpuGroupStateModel()),
			Config: tfsdk.Config{Schema: testContext.schema},
		},
		updateResp,
	)
	assertNoErrorDiagnostics(t, updateResp.Diagnostics)

	if modifiedName != "renamed" {
		t.Fatalf("ModifyGPUResourceGroup received name %q, want %q", modifiedName, "renamed")
	}

	var state svcags.AgsGPUResourceGroupResourceModel
	if diags := updateResp.State.Get(context.Background(), &state); diags.HasError() {
		t.Fatalf("failed to read updated state: %v", diags.Errors())
	}
	assertGpuGroupString(t, state.Name, "renamed")
	assertGpuGroupString(t, state.ResourceGroupId, agsTestGroupId)
}

func TestAgsGpuResourceGroupUpdateApiError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	testContext := newAgsGpuResourceGroupTestContext(t, patches)

	patches.ApplyMethodFunc(
		testContext.client,
		"ModifyGPUResourceGroupWithContext",
		func(_ context.Context, _ *agsv20250920.ModifyGPUResourceGroupRequest) (*agsv20250920.ModifyGPUResourceGroupResponse, error) {
			return nil, &mockSdkError{message: "api failure"}
		},
	)

	planModel := gpuGroupStateModel()
	planModel.Name = types.StringValue("renamed")

	updateResp := &resource.UpdateResponse{State: tfsdk.State{Schema: testContext.schema}}
	testContext.resource.Update(
		context.Background(),
		resource.UpdateRequest{
			Plan:   testContext.plan(t, planModel),
			State:  testContext.state(t, gpuGroupStateModel()),
			Config: tfsdk.Config{Schema: testContext.schema},
		},
		updateResp,
	)
	if !updateResp.Diagnostics.HasError() {
		t.Fatalf("expected an error diagnostic for a failed update")
	}
}

func TestAgsGpuResourceGroupImportState(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	testContext := newAgsGpuResourceGroupTestContext(t, patches)

	importer, ok := testContext.resource.(resource.ResourceWithImportState)
	if !ok {
		t.Fatalf("resource does not implement resource.ResourceWithImportState")
	}

	importResp := &resource.ImportStateResponse{State: tfsdk.State{Schema: testContext.schema, Raw: tftypes.NewValue(testContext.schema.Type().TerraformType(context.Background()), nil)}}
	importer.ImportState(
		context.Background(),
		resource.ImportStateRequest{ID: agsTestGroupId},
		importResp,
	)
	assertNoErrorDiagnostics(t, importResp.Diagnostics)

	var state svcags.AgsGPUResourceGroupResourceModel
	if diags := importResp.State.Get(context.Background(), &state); diags.HasError() {
		t.Fatalf("failed to read imported state: %v", diags.Errors())
	}
	if state.Id.ValueString() != agsTestGroupId {
		t.Fatalf("imported id is %q, want %q", state.Id.ValueString(), agsTestGroupId)
	}
	if state.ResourceGroupId.ValueString() != agsTestGroupId {
		t.Fatalf("imported resource_group_id is %q, want %q", state.ResourceGroupId.ValueString(), agsTestGroupId)
	}
}

func TestAgsGpuResourceGroupImportStateEmptyId(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	testContext := newAgsGpuResourceGroupTestContext(t, patches)

	importer := testContext.resource.(resource.ResourceWithImportState)
	importResp := &resource.ImportStateResponse{State: tfsdk.State{Schema: testContext.schema, Raw: tftypes.NewValue(testContext.schema.Type().TerraformType(context.Background()), nil)}}
	importer.ImportState(
		context.Background(),
		resource.ImportStateRequest{ID: ""},
		importResp,
	)
	if !importResp.Diagnostics.HasError() {
		t.Fatalf("expected an error diagnostic for an empty import id")
	}
}
