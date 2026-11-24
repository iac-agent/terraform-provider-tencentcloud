package ags_test

import (
	"context"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	agsv20250920 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ags/v20250920"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/sharedmeta"
	svcags "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/ags"
)

const (
	agsTestInstanceId   = "sbi-001"
	agsTestInstanceTool = "tool-001"
	agsTestExpiresAt    = "2026-01-01T00:10:00Z"
)

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

// sampleInstanceModel returns a plan model with every configurable leaf populated.
func (c *agsSandboxInstanceTestContext) sampleInstanceModel(t *testing.T) svcags.AgsSandboxInstanceResourceModel {
	t.Helper()

	top := c.schema.Attributes

	mountOptionTypes := listElementAttributeTypes(t, top, "mount_options")
	metadataTypes := listElementAttributeTypes(t, top, "metadata")

	customAttributes := top["custom_configuration"].(schema.SingleNestedAttribute).Attributes
	probeTypes := customNestedTypes(t, customAttributes, "probe")
	probeAttributes := customAttributes["probe"].(schema.SingleNestedAttribute).Attributes
	httpGetTypes := customNestedTypes(t, probeAttributes, "http_get")
	resourcesTypes := objectAttributeTypes(t, customAttributes, "resources")
	customConfigurationTypes := objectAttributeTypes(t, top, "custom_configuration")

	return svcags.AgsSandboxInstanceResourceModel{
		Id:               types.StringNull(),
		InstanceId:       types.StringNull(),
		ToolId:           types.StringNull(),
		ToolName:         types.StringValue(agsTestInstanceTool),
		SourceInstanceId: types.StringNull(),
		SnapshotId:       types.StringNull(),
		Timeout:          types.StringValue("10m"),
		ClientToken:      types.StringNull(),
		MountOptions: mustList(t, mountOptionTypes, map[string]attr.Value{
			"name":       types.StringValue("workspace"),
			"mount_path": types.StringValue("/workspace"),
			"sub_path":   types.StringNull(),
			"read_only":  types.BoolValue(false),
		}),
		CustomConfiguration: mustObject(t, customConfigurationTypes, map[string]attr.Value{
			"image":               types.StringValue("img:1"),
			"image_registry_type": types.StringValue("personal"),
			"command":             mustStringSet(t, "/bin/sh"),
			"args":                mustStringSet(t, "-c"),
			"env":                 types.ListNull(types.ObjectType{AttrTypes: listElementAttributeTypes(t, customAttributes, "env")}),
			"ports":               types.ListNull(types.ObjectType{AttrTypes: listElementAttributeTypes(t, customAttributes, "ports")}),
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
			"dns_config":       types.ObjectNull(customNestedTypes(t, customAttributes, "dns_config")),
			"security_context": types.ObjectNull(customNestedTypes(t, customAttributes, "security_context")),
		}),
		AuthMode: types.StringNull(),
		Metadata: mustList(t, metadataTypes, map[string]attr.Value{
			"name":  types.StringValue("owner"),
			"value": types.StringValue("tf"),
		}),
		NetworkMode:    types.StringNull(),
		Persistent:     types.BoolNull(),
		Status:         types.StringNull(),
		TimeoutSeconds: types.Int64Null(),
		ExpiresAt:      types.StringNull(),
		StopReason:     types.StringNull(),
		CreateTime:     types.StringNull(),
		UpdateTime:     types.StringNull(),
	}
}

// sampleSandboxInstance returns the cloud side representation used by the describe mock.
func sampleSandboxInstance() *agsv20250920.SandboxInstance {
	return &agsv20250920.SandboxInstance{
		InstanceId:     agsString(agsTestInstanceId),
		ToolId:         agsString(agsTestInstanceTool),
		ToolName:       agsString(agsTestInstanceTool),
		Status:         agsString("RUNNING"),
		Persistent:     agsBool(false),
		TimeoutSeconds: helperUint64(600),
		ExpiresAt:      agsString(agsTestExpiresAt),
		CreateTime:     agsString(agsTestCreateTime),
		UpdateTime:     agsString(agsTestUpdateTime),
		MountOptions: []*agsv20250920.MountOption{
			{
				Name:      agsString("workspace"),
				MountPath: agsString("/workspace"),
				ReadOnly:  agsBool(false),
			},
		},
		CustomConfiguration: &agsv20250920.CustomConfigurationDetail{
			Image:             agsString("img:1"),
			ImageRegistryType: agsString("personal"),
			ImageDigest:       agsString("sha256:abc"),
			Command:           []*string{agsString("/bin/sh")},
			Args:              []*string{agsString("-c")},
			Resources: &agsv20250920.ResourceConfiguration{
				CPU:     agsString("1"),
				Memory:  agsString("2Gi"),
				Storage: agsString("5Gi"),
			},
			Probe: &agsv20250920.ProbeConfiguration{
				HttpGet: &agsv20250920.HttpGetAction{
					Path:   agsString("/healthz"),
					Port:   agsInt64(8080),
					Scheme: agsString("HTTP"),
				},
				ReadyTimeoutMs:   agsInt64(1000),
				ProbeTimeoutMs:   agsInt64(2000),
				ProbePeriodMs:    agsInt64(3000),
				SuccessThreshold: agsInt64(1),
				FailureThreshold: agsInt64(3),
			},
		},
		NetworkMode: agsString("PUBLIC"),
		Metadata: []*agsv20250920.MetadataVar{
			{Name: agsString("owner"), Value: agsString("tf")},
		},
		AuthMode: agsString("DEFAULT"),
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

	computed := []string{
		"id", "instance_id", "network_mode", "persistent", "status",
		"timeout_seconds", "expires_at", "stop_reason", "create_time", "update_time",
	}
	for _, name := range computed {
		attribute, ok := top[name]
		if !ok {
			t.Fatalf("attribute %q is missing", name)
		}
		if !attribute.IsComputed() {
			t.Fatalf("attribute %q must be computed", name)
		}
		if attribute.IsRequired() {
			t.Fatalf("attribute %q must not be required", name)
		}
	}

	optional := []string{
		"tool_id", "tool_name", "source_instance_id", "snapshot_id", "timeout",
		"mount_options", "custom_configuration", "auth_mode", "metadata",
	}
	for _, name := range optional {
		attribute, ok := top[name]
		if !ok {
			t.Fatalf("attribute %q is missing", name)
		}
		if !attribute.IsOptional() {
			t.Fatalf("attribute %q must be optional", name)
		}
	}
}

// TestAgsSandboxInstanceCreate verifies the start request and the read back state.
func TestAgsSandboxInstanceCreate(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	testContext := newAgsSandboxInstanceTestContext(t, patches)

	var capturedRequest *agsv20250920.StartSandboxInstanceRequest
	patches.ApplyMethodFunc(
		testContext.client,
		"StartSandboxInstanceWithContext",
		func(_ context.Context, request *agsv20250920.StartSandboxInstanceRequest) (*agsv20250920.StartSandboxInstanceResponse, error) {
			capturedRequest = request
			return &agsv20250920.StartSandboxInstanceResponse{
				Response: &agsv20250920.StartSandboxInstanceResponseParams{
					Instance: &agsv20250920.SandboxInstance{InstanceId: agsString(agsTestInstanceId)},
				},
			}, nil
		},
	)
	patches.ApplyMethodFunc(
		testContext.client,
		"DescribeSandboxInstancesWithContext",
		func(_ context.Context, request *agsv20250920.DescribeSandboxInstancesRequest) (*agsv20250920.DescribeSandboxInstancesResponse, error) {
			assertStringPointerList(t, "describe request InstanceIds", request.InstanceIds, agsTestInstanceId)
			return describeSandboxInstancesResponse(sampleSandboxInstance()), nil
		},
	)

	model := testContext.sampleInstanceModel(t)
	createResp := &resource.CreateResponse{State: tfsdk.State{Schema: testContext.schema}}
	testContext.resource.Create(context.Background(), resource.CreateRequest{Plan: testContext.plan(t, model)}, createResp)
	assertNoErrorDiagnostics(t, createResp.Diagnostics)

	if capturedRequest == nil {
		t.Fatalf("StartSandboxInstance was not called")
	}
	if capturedRequest.ToolName == nil || *capturedRequest.ToolName != agsTestInstanceTool {
		t.Fatalf("StartSandboxInstance ToolName is not %q", agsTestInstanceTool)
	}
	if capturedRequest.Timeout == nil || *capturedRequest.Timeout != "10m" {
		t.Fatalf("StartSandboxInstance Timeout is not 10m")
	}
	// client_token is generated internally as an idempotency key.
	if capturedRequest.ClientToken == nil || *capturedRequest.ClientToken == "" {
		t.Fatalf("StartSandboxInstance ClientToken is empty")
	}
	if len(capturedRequest.MountOptions) != 1 || capturedRequest.MountOptions[0].Name == nil ||
		*capturedRequest.MountOptions[0].Name != "workspace" {
		t.Fatalf("StartSandboxInstance MountOptions is not [workspace]")
	}
	if capturedRequest.MountOptions[0].SubPath != nil {
		t.Fatalf("StartSandboxInstance MountOptions[0].SubPath must not be sent")
	}
	if capturedRequest.CustomConfiguration == nil || capturedRequest.CustomConfiguration.Image == nil ||
		*capturedRequest.CustomConfiguration.Image != "img:1" {
		t.Fatalf("StartSandboxInstance CustomConfiguration.Image is not img:1")
	}
	if len(capturedRequest.Metadata) != 1 || capturedRequest.Metadata[0].Name == nil ||
		*capturedRequest.Metadata[0].Name != "owner" {
		t.Fatalf("StartSandboxInstance Metadata is not [owner]")
	}

	var state svcags.AgsSandboxInstanceResourceModel
	if diags := createResp.State.Get(context.Background(), &state); diags.HasError() {
		t.Fatalf("failed to read created state: %v", diags.Errors())
	}

	if state.Id.ValueString() != agsTestInstanceId {
		t.Fatalf("state id is %q, want %q", state.Id.ValueString(), agsTestInstanceId)
	}
	if state.InstanceId.ValueString() != agsTestInstanceId {
		t.Fatalf("state instance_id is %q, want %q", state.InstanceId.ValueString(), agsTestInstanceId)
	}
	if state.Status.ValueString() != "RUNNING" {
		t.Fatalf("state status is %q, want RUNNING", state.Status.ValueString())
	}
	if state.TimeoutSeconds.ValueInt64() != 600 {
		t.Fatalf("state timeout_seconds is %d, want 600", state.TimeoutSeconds.ValueInt64())
	}
	if state.ExpiresAt.ValueString() != agsTestExpiresAt {
		t.Fatalf("state expires_at is %q", state.ExpiresAt.ValueString())
	}
	if state.NetworkMode.ValueString() != "PUBLIC" {
		t.Fatalf("state network_mode is %q, want PUBLIC", state.NetworkMode.ValueString())
	}
	// create only inputs are not reported by the describe API and must be kept.
	if state.Timeout.ValueString() != "10m" {
		t.Fatalf("state timeout is %q, want 10m", state.Timeout.ValueString())
	}

	top := map[string]attr.Value{
		"id":                   state.Id,
		"instance_id":          state.InstanceId,
		"tool_id":              state.ToolId,
		"tool_name":            state.ToolName,
		"mount_options":        state.MountOptions,
		"custom_configuration": state.CustomConfiguration,
		"metadata":             state.Metadata,
		"network_mode":         state.NetworkMode,
		"persistent":           state.Persistent,
		"status":               state.Status,
		"timeout_seconds":      state.TimeoutSeconds,
		"expires_at":           state.ExpiresAt,
		"stop_reason":          state.StopReason,
		"create_time":          state.CreateTime,
		"update_time":          state.UpdateTime,
	}

	assertString(t, top, "create_time", agsTestCreateTime)
	assertString(t, top, "update_time", agsTestUpdateTime)

	mountOptions := listElements(t, top, "mount_options")
	if len(mountOptions) != 1 {
		t.Fatalf("state mount_options has %d elements, want 1", len(mountOptions))
	}
	mountOption := attributesOf(t, mountOptions[0])
	assertString(t, mountOption, "name", "workspace")
	assertString(t, mountOption, "mount_path", "/workspace")

	customConfiguration := attributesOf(t, top["custom_configuration"])
	assertString(t, customConfiguration, "image", "img:1")
	assertStringSet(t, customConfiguration, "command", "/bin/sh")

	metadata := listElements(t, top, "metadata")
	if len(metadata) != 1 {
		t.Fatalf("state metadata has %d elements, want 1", len(metadata))
	}
	assertString(t, attributesOf(t, metadata[0]), "name", "owner")
}

// TestAgsSandboxInstanceReadNotFound verifies state is dropped for a vanished instance.
func TestAgsSandboxInstanceReadNotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	testContext := newAgsSandboxInstanceTestContext(t, patches)

	patches.ApplyMethodFunc(
		testContext.client,
		"DescribeSandboxInstancesWithContext",
		func(_ context.Context, _ *agsv20250920.DescribeSandboxInstancesRequest) (*agsv20250920.DescribeSandboxInstancesResponse, error) {
			return describeSandboxInstancesResponse(), nil
		},
	)

	stateModel := testContext.sampleInstanceModel(t)
	stateModel.Id = types.StringValue(agsTestInstanceId)
	stateModel.InstanceId = types.StringValue(agsTestInstanceId)

	readResp := &resource.ReadResponse{State: testContext.state(t, stateModel)}
	testContext.resource.Read(context.Background(), resource.ReadRequest{State: testContext.state(t, stateModel)}, readResp)
	assertNoErrorDiagnostics(t, readResp.Diagnostics)

	if !readResp.State.Raw.IsNull() {
		t.Fatalf("state must be removed when the sandbox instance does not exist")
	}
}

// TestAgsSandboxInstanceImportState verifies the import id is stored as identity.
func TestAgsSandboxInstanceImportState(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	testContext := newAgsSandboxInstanceTestContext(t, patches)

	importer, ok := testContext.resource.(resource.ResourceWithImportState)
	if !ok {
		t.Fatalf("resource does not implement resource.ResourceWithImportState")
	}

	importResp := &resource.ImportStateResponse{State: tfsdk.State{Schema: testContext.schema}}
	importer.ImportState(context.Background(), resource.ImportStateRequest{ID: agsTestInstanceId}, importResp)
	assertNoErrorDiagnostics(t, importResp.Diagnostics)

	var importedId types.String
	if diags := importResp.State.GetAttribute(context.Background(), path.Empty().AtName("id"), &importedId); diags.HasError() {
		t.Fatalf("failed to read imported id: %v", diags.Errors())
	}
	if importedId.ValueString() != agsTestInstanceId {
		t.Fatalf("state id is %q, want %q", importedId.ValueString(), agsTestInstanceId)
	}
}

func customNestedTypes(t *testing.T, attributes map[string]schema.Attribute, name string) map[string]attr.Type {
	t.Helper()
	return objectAttributeTypes(t, attributes, name)
}
