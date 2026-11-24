package ags_test

import (
	"context"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"

	agsv20250920 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ags/v20250920"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/sharedmeta"
	svcags "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/ags"
)

const agsTestInstanceId = "inst-001"

// agsSandboxInstanceTestContext wires a resource instance to a patched SDK client.
type agsSandboxInstanceTestContext struct {
	resource resource.Resource
	client   *agsv20250920.Client
	schema   schema.Schema
}

func newAgsSandboxInstanceTestContext(t *testing.T, patches *gomonkey.Patches) *agsSandboxInstanceTestContext {
	t.Helper()

	connectivityClient := &connectivity.TencentCloudClient{}
	agsClient := &agsv20250920.Client{}
	patches.ApplyMethodReturn(connectivityClient, "UseAgsClient", agsClient)

	res := svcags.NewAgsSandboxInstanceResource()
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

	return &agsSandboxInstanceTestContext{resource: res, client: agsClient, schema: schemaResp.Schema}
}

func (c *agsSandboxInstanceTestContext) plan(t *testing.T, model svcags.AgsSandboxInstanceResourceModel) tfsdk.Plan {
	t.Helper()

	plan := tfsdk.Plan{Schema: c.schema}
	if diags := plan.Set(context.Background(), model); diags.HasError() {
		t.Fatalf("failed to build plan: %v", diags.Errors())
	}
	return plan
}

func (c *agsSandboxInstanceTestContext) state(t *testing.T, model svcags.AgsSandboxInstanceResourceModel) tfsdk.State {
	t.Helper()

	state := tfsdk.State{Schema: c.schema}
	if diags := state.Set(context.Background(), model); diags.HasError() {
		t.Fatalf("failed to build state: %v", diags.Errors())
	}
	return state
}

func objectAttributeTypesOfInstance(t *testing.T, attributes map[string]schema.Attribute, name string) map[string]attr.Type {
	t.Helper()

	nested, ok := attributes[name].(schema.SingleNestedAttribute)
	if !ok {
		t.Fatalf("attribute %q is not a single nested attribute", name)
	}
	objectType, ok := nested.GetType().(types.ObjectType)
	if !ok {
		t.Fatalf("attribute %q does not have an object type", name)
	}
	return objectType.AttrTypes
}

func listElementAttributeTypesOfInstance(t *testing.T, attributes map[string]schema.Attribute, name string) map[string]attr.Type {
	t.Helper()

	nested, ok := attributes[name].(schema.ListNestedAttribute)
	if !ok {
		t.Fatalf("attribute %q is not a list nested attribute", name)
	}
	listType, ok := nested.GetType().(types.ListType)
	if !ok {
		t.Fatalf("attribute %q does not have a list type", name)
	}
	objectType, ok := listType.ElemType.(types.ObjectType)
	if !ok {
		t.Fatalf("attribute %q does not have an object element type", name)
	}
	return objectType.AttrTypes
}

// sampleInstanceModel returns a plan model with every configurable leaf populated.
func (c *agsSandboxInstanceTestContext) sampleInstanceModel(t *testing.T) svcags.AgsSandboxInstanceResourceModel {
	t.Helper()

	top := c.schema.Attributes

	mountOptionTypes := listElementAttributeTypesOfInstance(t, top, "mount_options")
	metadataTypes := listElementAttributeTypesOfInstance(t, top, "metadata")

	customAttributes := top["custom_configuration"].(schema.SingleNestedAttribute).Attributes
	probeAttributes := customAttributes["probe"].(schema.SingleNestedAttribute).Attributes
	securityContextAttributes := customAttributes["security_context"].(schema.SingleNestedAttribute).Attributes
	resourcesTypes := objectAttributeTypesOfInstance(t, customAttributes, "resources")
	httpGetTypes := objectAttributeTypesOfInstance(t, probeAttributes, "http_get")
	dnsConfigTypes := objectAttributeTypesOfInstance(t, customAttributes, "dns_config")
	capabilitiesTypes := objectAttributeTypesOfInstance(t, securityContextAttributes, "capabilities")
	securityContextTypes := objectAttributeTypesOfInstance(t, customAttributes, "security_context")
	probeTypes := objectAttributeTypesOfInstance(t, customAttributes, "probe")
	customConfigurationTypes := objectAttributeTypesOfInstance(t, top, "custom_configuration")
	envTypes := listElementAttributeTypesOfInstance(t, customAttributes, "env")
	portTypes := listElementAttributeTypesOfInstance(t, customAttributes, "ports")

	return svcags.AgsSandboxInstanceResourceModel{
		Id:               types.StringNull(),
		InstanceId:       types.StringNull(),
		ToolId:           types.StringValue(agsTestToolId),
		ToolName:         types.StringNull(),
		SourceInstanceId: types.StringNull(),
		SnapshotId:       types.StringNull(),
		Timeout:          types.StringValue("10m"),
		ClientToken:      types.StringNull(),
		MountOptions: mustList(t, mountOptionTypes, map[string]attr.Value{
			"name":       types.StringValue("workspace"),
			"mount_path": types.StringValue("/workspace"),
			"sub_path":   types.StringValue("/sub"),
			"read_only":  types.BoolValue(true),
		}),
		CustomConfiguration: mustObject(t, customConfigurationTypes, map[string]attr.Value{
			"image":               types.StringValue("img:1"),
			"image_registry_type": types.StringValue("personal"),
			"command":             mustStringSet(t, "/bin/sh"),
			"args":                mustStringSet(t, "-c"),
			"env": mustList(t, envTypes, map[string]attr.Value{
				"name":  types.StringValue("K"),
				"value": types.StringValue("V"),
			}),
			"ports": mustList(t, portTypes, map[string]attr.Value{
				"name":     types.StringValue("http"),
				"port":     types.Int64Value(8080),
				"protocol": types.StringValue("TCP"),
			}),
			"resources": mustObject(t, resourcesTypes, map[string]attr.Value{
				"cpu":     types.StringValue("1"),
				"memory":  types.StringValue("2Gi"),
				"storage": types.StringValue("5Gi"),
			}),
			"probe": mustObject(t, probeTypes, map[string]attr.Value{
				"http_get": mustObject(t, httpGetTypes, map[string]attr.Value{
					"path":   types.StringValue("/healthz"),
					"port":   types.Int64Value(8080),
					"scheme": types.StringValue("HTTP"),
				}),
				"ready_timeout_ms":  types.Int64Value(1000),
				"probe_timeout_ms":  types.Int64Value(2000),
				"probe_period_ms":   types.Int64Value(3000),
				"success_threshold": types.Int64Value(1),
				"failure_threshold": types.Int64Value(3),
			}),
			"dns_config": mustObject(t, dnsConfigTypes, map[string]attr.Value{
				"servers":  mustStringSet(t, "10.0.0.1"),
				"searches": mustStringSet(t, "example.com"),
				"options":  mustStringSet(t, "ndots:2"),
			}),
			"security_context": mustObject(t, securityContextTypes, map[string]attr.Value{
				"privileged": types.BoolValue(true),
				"capabilities": mustObject(t, capabilitiesTypes, map[string]attr.Value{
					"add":  mustStringSet(t, "SYS_ADMIN"),
					"drop": mustStringSet(t, "ALL"),
				}),
			}),
		}),
		AuthMode: types.StringValue("TOKEN"),
		Metadata: mustList(t, metadataTypes, map[string]attr.Value{
			"name":  types.StringValue("k1"),
			"value": types.StringValue("v1"),
		}),
		Status:         types.StringNull(),
		Persistent:     types.BoolNull(),
		TimeoutSeconds: types.Int64Null(),
		ExpiresAt:      types.StringNull(),
		StopReason:     types.StringNull(),
		NetworkMode:    types.StringNull(),
		CreateTime:     types.StringNull(),
		UpdateTime:     types.StringNull(),
	}
}

// sampleInstance returns the cloud side representation used by the describe mock.
func sampleInstance() *agsv20250920.SandboxInstance {
	return &agsv20250920.SandboxInstance{
		InstanceId:     agsString(agsTestInstanceId),
		ToolId:         agsString(agsTestToolId),
		ToolName:       agsString(agsTestToolName),
		Status:         agsString("RUNNING"),
		Persistent:     agsBool(false),
		TimeoutSeconds: helperUint64(600),
		ExpiresAt:      agsString("2026-01-01T01:00:00Z"),
		MountOptions: []*agsv20250920.MountOption{
			{
				Name:      agsString("workspace"),
				MountPath: agsString("/workspace"),
				SubPath:   agsString("/sub"),
				ReadOnly:  agsBool(true),
			},
		},
		CustomConfiguration: &agsv20250920.CustomConfigurationDetail{
			Image:             agsString("img:1"),
			ImageRegistryType: agsString("personal"),
			Env: []*agsv20250920.EnvVar{
				{Name: agsString("K"), Value: agsString("V")},
			},
			Resources: &agsv20250920.ResourceConfiguration{
				CPU:    agsString("1"),
				Memory: agsString("2Gi"),
			},
		},
		NetworkMode: agsString("PUBLIC"),
		Metadata: []*agsv20250920.MetadataVar{
			{Name: agsString("k1"), Value: agsString("v1")},
		},
		AuthMode:   agsString("TOKEN"),
		CreateTime: agsString(agsTestCreateTime),
		UpdateTime: agsString(agsTestUpdateTime),
	}
}

func describeSandboxInstancesResponse(instances ...*agsv20250920.SandboxInstance) *agsv20250920.DescribeSandboxInstancesResponse {
	return &agsv20250920.DescribeSandboxInstancesResponse{
		Response: &agsv20250920.DescribeSandboxInstancesResponseParams{
			InstanceSet: instances,
		},
	}
}

// TestAgsSandboxInstanceSchema verifies the attribute set and its required flags.
func TestAgsSandboxInstanceSchema(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	testContext := newAgsSandboxInstanceTestContext(t, patches)
	top := testContext.schema.Attributes

	if len(top) != 20 {
		t.Fatalf("schema has %d attributes, want 20", len(top))
	}

	computed := []string{"id", "instance_id", "status", "persistent", "timeout_seconds", "expires_at", "stop_reason", "network_mode", "create_time", "update_time"}
	for _, name := range computed {
		attribute, ok := top[name]
		if !ok {
			t.Fatalf("attribute %q is missing", name)
		}
		if !attribute.IsComputed() {
			t.Fatalf("attribute %q must be computed", name)
		}
		if attribute.IsRequired() || attribute.IsOptional() {
			t.Fatalf("attribute %q must not be configurable", name)
		}
	}

	optional := []string{"tool_id", "tool_name", "source_instance_id", "snapshot_id", "timeout", "client_token", "mount_options", "custom_configuration", "auth_mode", "metadata"}
	for _, name := range optional {
		attribute, ok := top[name]
		if !ok {
			t.Fatalf("attribute %q is missing", name)
		}
		if !attribute.IsOptional() {
			t.Fatalf("attribute %q must be optional", name)
		}
	}

	authMode := top["auth_mode"].(schema.StringAttribute)
	if !authMode.IsComputed() {
		t.Fatalf("auth_mode must be computed as well")
	}
}

// TestAgsSandboxInstanceCreate verifies the create request and the read back state.
func TestAgsSandboxInstanceCreate(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	testContext := newAgsSandboxInstanceTestContext(t, patches)

	var capturedStartRequest *agsv20250920.StartSandboxInstanceRequest
	patches.ApplyMethodFunc(
		testContext.client,
		"StartSandboxInstanceWithContext",
		func(_ context.Context, request *agsv20250920.StartSandboxInstanceRequest) (*agsv20250920.StartSandboxInstanceResponse, error) {
			capturedStartRequest = request
			return &agsv20250920.StartSandboxInstanceResponse{
				Response: &agsv20250920.StartSandboxInstanceResponseParams{Instance: sampleInstance()},
			}, nil
		},
	)
	patches.ApplyMethodFunc(
		testContext.client,
		"DescribeSandboxInstancesWithContext",
		func(_ context.Context, request *agsv20250920.DescribeSandboxInstancesRequest) (*agsv20250920.DescribeSandboxInstancesResponse, error) {
			assertStringPointerList(t, "describe request InstanceIds", request.InstanceIds, agsTestInstanceId)
			return describeSandboxInstancesResponse(sampleInstance()), nil
		},
	)

	model := testContext.sampleInstanceModel(t)
	createResp := &resource.CreateResponse{State: tfsdk.State{Schema: testContext.schema}}
	testContext.resource.Create(context.Background(), resource.CreateRequest{Plan: testContext.plan(t, model)}, createResp)
	assertNoErrorDiagnostics(t, createResp.Diagnostics)

	if capturedStartRequest == nil {
		t.Fatalf("StartSandboxInstance was not called")
	}
	if capturedStartRequest.ToolId == nil || *capturedStartRequest.ToolId != agsTestToolId {
		t.Fatalf("StartSandboxInstance ToolId is not %q", agsTestToolId)
	}
	if capturedStartRequest.Timeout == nil || *capturedStartRequest.Timeout != "10m" {
		t.Fatalf("StartSandboxInstance Timeout is not 10m")
	}
	// client_token is generated internally as an idempotency key.
	if capturedStartRequest.ClientToken == nil || *capturedStartRequest.ClientToken == "" {
		t.Fatalf("StartSandboxInstance ClientToken is empty")
	}
	if capturedStartRequest.AuthMode == nil || *capturedStartRequest.AuthMode != "TOKEN" {
		t.Fatalf("StartSandboxInstance AuthMode is not TOKEN")
	}
	if len(capturedStartRequest.MountOptions) != 1 {
		t.Fatalf("StartSandboxInstance MountOptions has %d elements, want 1", len(capturedStartRequest.MountOptions))
	}
	if capturedStartRequest.MountOptions[0].Name == nil || *capturedStartRequest.MountOptions[0].Name != "workspace" {
		t.Fatalf("MountOption.Name is not workspace")
	}
	if capturedStartRequest.MountOptions[0].ReadOnly == nil || !*capturedStartRequest.MountOptions[0].ReadOnly {
		t.Fatalf("MountOption.ReadOnly is not true")
	}
	if capturedStartRequest.CustomConfiguration == nil || capturedStartRequest.CustomConfiguration.Image == nil ||
		*capturedStartRequest.CustomConfiguration.Image != "img:1" {
		t.Fatalf("StartSandboxInstance CustomConfiguration.Image is not img:1")
	}
	if len(capturedStartRequest.Metadata) != 1 || capturedStartRequest.Metadata[0].Name == nil || *capturedStartRequest.Metadata[0].Name != "k1" {
		t.Fatalf("StartSandboxInstance Metadata is not [k1=v1]")
	}

	var state svcags.AgsSandboxInstanceResourceModel
	if diags := createResp.State.Get(context.Background(), &state); diags.HasError() {
		t.Fatalf("failed to read created state: %v", diags.Errors())
	}
	if state.Id.ValueString() != agsTestInstanceId || state.InstanceId.ValueString() != agsTestInstanceId {
		t.Fatalf("created state identity is not %q", agsTestInstanceId)
	}
	if state.Status.ValueString() != "RUNNING" {
		t.Fatalf("status is %q, want RUNNING", state.Status.ValueString())
	}
	if state.TimeoutSeconds.ValueInt64() != 600 {
		t.Fatalf("timeout_seconds is %d, want 600", state.TimeoutSeconds.ValueInt64())
	}
	if state.NetworkMode.ValueString() != "PUBLIC" {
		t.Fatalf("network_mode is %q, want PUBLIC", state.NetworkMode.ValueString())
	}
	if state.AuthMode.ValueString() != "TOKEN" {
		t.Fatalf("auth_mode is %q, want TOKEN", state.AuthMode.ValueString())
	}
	// timeout is create-only and is kept as it was configured.
	if state.Timeout.ValueString() != "10m" {
		t.Fatalf("timeout is %q, want 10m", state.Timeout.ValueString())
	}
}

// TestAgsSandboxInstanceRead verifies the refresh of an existing instance.
func TestAgsSandboxInstanceRead(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	testContext := newAgsSandboxInstanceTestContext(t, patches)

	patches.ApplyMethodFunc(
		testContext.client,
		"DescribeSandboxInstancesWithContext",
		func(_ context.Context, request *agsv20250920.DescribeSandboxInstancesRequest) (*agsv20250920.DescribeSandboxInstancesResponse, error) {
			return describeSandboxInstancesResponse(sampleInstance()), nil
		},
	)

	model := testContext.sampleInstanceModel(t)
	model.Id = types.StringValue(agsTestInstanceId)
	model.InstanceId = types.StringValue(agsTestInstanceId)

	readResp := &resource.ReadResponse{State: tfsdk.State{Schema: testContext.schema}}
	testContext.resource.Read(context.Background(), resource.ReadRequest{State: testContext.state(t, model)}, readResp)
	assertNoErrorDiagnostics(t, readResp.Diagnostics)

	var state svcags.AgsSandboxInstanceResourceModel
	if diags := readResp.State.Get(context.Background(), &state); diags.HasError() {
		t.Fatalf("failed to read refreshed state: %v", diags.Errors())
	}
	if state.Status.ValueString() != "RUNNING" {
		t.Fatalf("status is %q, want RUNNING", state.Status.ValueString())
	}
	// timeout is create-only and is kept as it is stored in state.
	if state.Timeout.ValueString() != "10m" {
		t.Fatalf("timeout is %q, want 10m", state.Timeout.ValueString())
	}
}

// TestAgsSandboxInstanceReadNotFound verifies a deleted instance is removed from state.
func TestAgsSandboxInstanceReadNotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	testContext := newAgsSandboxInstanceTestContext(t, patches)

	patches.ApplyMethodFunc(
		testContext.client,
		"DescribeSandboxInstancesWithContext",
		func(_ context.Context, request *agsv20250920.DescribeSandboxInstancesRequest) (*agsv20250920.DescribeSandboxInstancesResponse, error) {
			return describeSandboxInstancesResponse(), nil
		},
	)

	model := testContext.sampleInstanceModel(t)
	model.Id = types.StringValue(agsTestInstanceId)
	model.InstanceId = types.StringValue(agsTestInstanceId)

	readResp := &resource.ReadResponse{State: tfsdk.State{Schema: testContext.schema}}
	testContext.resource.Read(context.Background(), resource.ReadRequest{State: testContext.state(t, model)}, readResp)
	assertNoErrorDiagnostics(t, readResp.Diagnostics)

	if !readResp.State.Raw.IsNull() {
		t.Fatalf("state must be removed when the sandbox instance does not exist")
	}
}

// TestAgsSandboxInstanceDelete verifies the stop request of destroy.
func TestAgsSandboxInstanceDelete(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	testContext := newAgsSandboxInstanceTestContext(t, patches)

	var capturedStopRequest *agsv20250920.StopSandboxInstanceRequest
	patches.ApplyMethodFunc(
		testContext.client,
		"StopSandboxInstanceWithContext",
		func(_ context.Context, request *agsv20250920.StopSandboxInstanceRequest) (*agsv20250920.StopSandboxInstanceResponse, error) {
			capturedStopRequest = request
			return &agsv20250920.StopSandboxInstanceResponse{
				Response: &agsv20250920.StopSandboxInstanceResponseParams{},
			}, nil
		},
	)

	model := testContext.sampleInstanceModel(t)
	model.Id = types.StringValue(agsTestInstanceId)
	model.InstanceId = types.StringValue(agsTestInstanceId)

	deleteResp := &resource.DeleteResponse{State: tfsdk.State{Schema: testContext.schema}}
	testContext.resource.Delete(context.Background(), resource.DeleteRequest{State: testContext.state(t, model)}, deleteResp)
	assertNoErrorDiagnostics(t, deleteResp.Diagnostics)

	if capturedStopRequest == nil {
		t.Fatalf("StopSandboxInstance was not called")
	}
	if capturedStopRequest.InstanceId == nil || *capturedStopRequest.InstanceId != agsTestInstanceId {
		t.Fatalf("StopSandboxInstance InstanceId is not %q", agsTestInstanceId)
	}
}
