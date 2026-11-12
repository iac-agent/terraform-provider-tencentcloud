package ags_test

import (
	"context"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
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
	agsTestPreCacheImageId     = "cache-001"
	agsTestPreCacheImage       = "ccr.ccs.tencentyun.com/tf-example/demo:v1"
	agsTestPreCacheRegistry    = "personal"
	agsTestPreCacheDigest      = "sha256:0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"
	agsTestPreCacheStatus      = "READY"
	agsTestPreCacheMessage     = "pre cache finished"
	agsTestPreCacheCreateTime  = "2026-01-01T00:00:00Z"
	agsTestPreCacheReadyTime   = "2026-01-01T00:10:00Z"
	agsTestPreCacheLastUsed    = "2026-01-02T00:00:00Z"
	agsTestPreCacheSourceType  = "EXPLICIT"
	agsTestPreCacheCachedBytes = int64(1048576)
)

// agsPreCacheImageTestContext wires a resource instance to a patched SDK client.
type agsPreCacheImageTestContext struct {
	resource resource.Resource
	client   *agsv20250920.Client
	schema   schema.Schema
}

func newAgsPreCacheImageTestContext(t *testing.T, patches *gomonkey.Patches) *agsPreCacheImageTestContext {
	t.Helper()

	connectivityClient := &connectivity.TencentCloudClient{}
	agsClient := &agsv20250920.Client{}
	patches.ApplyMethodReturn(connectivityClient, "UseAgsClient", agsClient)

	res := svcags.NewAgsPreCacheImageTaskResource()
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

	return &agsPreCacheImageTestContext{resource: res, client: agsClient, schema: schemaResp.Schema}
}

func (c *agsPreCacheImageTestContext) plan(t *testing.T, model svcags.AgsPreCacheImageTaskResourceModel) tfsdk.Plan {
	t.Helper()

	plan := tfsdk.Plan{Schema: c.schema}
	if diags := plan.Set(context.Background(), model); diags.HasError() {
		t.Fatalf("failed to build plan: %v", diags.Errors())
	}
	return plan
}

func (c *agsPreCacheImageTestContext) state(t *testing.T, model svcags.AgsPreCacheImageTaskResourceModel) tfsdk.State {
	t.Helper()

	state := tfsdk.State{Schema: c.schema}
	if diags := state.Set(context.Background(), model); diags.HasError() {
		t.Fatalf("failed to build state: %v", diags.Errors())
	}
	return state
}

// sampleModel returns the model of a resource that is not created yet, so the
// computed attributes stay unknown.
func (c *agsPreCacheImageTestContext) sampleModel() svcags.AgsPreCacheImageTaskResourceModel {
	return svcags.AgsPreCacheImageTaskResourceModel{
		Id:                   types.StringUnknown(),
		PreCacheImageId:      types.StringUnknown(),
		Image:                types.StringValue(agsTestPreCacheImage),
		ImageRegistryType:    types.StringValue(agsTestPreCacheRegistry),
		ImageDigest:          types.StringUnknown(),
		Status:               types.StringUnknown(),
		Message:              types.StringUnknown(),
		CreateTime:           types.StringUnknown(),
		ReadyTime:            types.StringUnknown(),
		LastUsedTime:         types.StringUnknown(),
		SourceType:           types.StringUnknown(),
		CachedImageSizeBytes: types.Int64Unknown(),
	}
}

// createdModel returns the model of a resource that already exists.
func (c *agsPreCacheImageTestContext) createdModel() svcags.AgsPreCacheImageTaskResourceModel {
	return svcags.AgsPreCacheImageTaskResourceModel{
		Id:                   types.StringValue(agsTestPreCacheImageId),
		PreCacheImageId:      types.StringValue(agsTestPreCacheImageId),
		Image:                types.StringValue(agsTestPreCacheImage),
		ImageRegistryType:    types.StringValue(agsTestPreCacheRegistry),
		ImageDigest:          types.StringValue(agsTestPreCacheDigest),
		Status:               types.StringValue(agsTestPreCacheStatus),
		Message:              types.StringValue(agsTestPreCacheMessage),
		CreateTime:           types.StringValue(agsTestPreCacheCreateTime),
		ReadyTime:            types.StringValue(agsTestPreCacheReadyTime),
		LastUsedTime:         types.StringValue(agsTestPreCacheLastUsed),
		SourceType:           types.StringValue(agsTestPreCacheSourceType),
		CachedImageSizeBytes: types.Int64Value(agsTestPreCacheCachedBytes),
	}
}

func describePreCacheImageResponse() *agsv20250920.DescribePreCacheImageTaskResponse {
	return &agsv20250920.DescribePreCacheImageTaskResponse{
		Response: &agsv20250920.DescribePreCacheImageTaskResponseParams{
			PreCacheImageId:      agsString(agsTestPreCacheImageId),
			Image:                agsString(agsTestPreCacheImage),
			ImageDigest:          agsString(agsTestPreCacheDigest),
			ImageRegistryType:    agsString(agsTestPreCacheRegistry),
			Status:               agsString(agsTestPreCacheStatus),
			Message:              agsString(agsTestPreCacheMessage),
			CreateTime:           agsString(agsTestPreCacheCreateTime),
			ReadyTime:            agsString(agsTestPreCacheReadyTime),
			LastUsedTime:         agsString(agsTestPreCacheLastUsed),
			SourceType:           agsString(agsTestPreCacheSourceType),
			CachedImageSizeBytes: agsInt64(agsTestPreCacheCachedBytes),
		},
	}
}

func assertPreCacheImageState(t *testing.T, state svcags.AgsPreCacheImageTaskResourceModel) {
	t.Helper()

	if state.Id.ValueString() != agsTestPreCacheImageId {
		t.Fatalf("id is %q, want %q", state.Id.ValueString(), agsTestPreCacheImageId)
	}
	if state.PreCacheImageId.ValueString() != agsTestPreCacheImageId {
		t.Fatalf("pre_cache_image_id is %q, want %q", state.PreCacheImageId.ValueString(), agsTestPreCacheImageId)
	}
	if state.Image.ValueString() != agsTestPreCacheImage {
		t.Fatalf("image is %q, want %q", state.Image.ValueString(), agsTestPreCacheImage)
	}
	if state.ImageRegistryType.ValueString() != agsTestPreCacheRegistry {
		t.Fatalf("image_registry_type is %q, want %q", state.ImageRegistryType.ValueString(), agsTestPreCacheRegistry)
	}
	if state.ImageDigest.ValueString() != agsTestPreCacheDigest {
		t.Fatalf("image_digest is %q, want %q", state.ImageDigest.ValueString(), agsTestPreCacheDigest)
	}
	if state.Status.ValueString() != agsTestPreCacheStatus {
		t.Fatalf("status is %q, want %q", state.Status.ValueString(), agsTestPreCacheStatus)
	}
	if state.Message.ValueString() != agsTestPreCacheMessage {
		t.Fatalf("message is %q, want %q", state.Message.ValueString(), agsTestPreCacheMessage)
	}
	if state.CreateTime.ValueString() != agsTestPreCacheCreateTime {
		t.Fatalf("create_time is %q, want %q", state.CreateTime.ValueString(), agsTestPreCacheCreateTime)
	}
	if state.ReadyTime.ValueString() != agsTestPreCacheReadyTime {
		t.Fatalf("ready_time is %q, want %q", state.ReadyTime.ValueString(), agsTestPreCacheReadyTime)
	}
	if state.LastUsedTime.ValueString() != agsTestPreCacheLastUsed {
		t.Fatalf("last_used_time is %q, want %q", state.LastUsedTime.ValueString(), agsTestPreCacheLastUsed)
	}
	if state.SourceType.ValueString() != agsTestPreCacheSourceType {
		t.Fatalf("source_type is %q, want %q", state.SourceType.ValueString(), agsTestPreCacheSourceType)
	}
	if state.CachedImageSizeBytes.IsNull() || state.CachedImageSizeBytes.IsUnknown() {
		t.Fatalf("cached_image_size_bytes has no value")
	}
	if state.CachedImageSizeBytes.ValueInt64() != agsTestPreCacheCachedBytes {
		t.Fatalf("cached_image_size_bytes is %d, want %d", state.CachedImageSizeBytes.ValueInt64(), agsTestPreCacheCachedBytes)
	}
}

// TestAgsPreCacheImageTaskSchema verifies the attribute set and its flags.
func TestAgsPreCacheImageTaskSchema(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	testContext := newAgsPreCacheImageTestContext(t, patches)
	top := testContext.schema.Attributes

	if len(top) != 12 {
		t.Fatalf("schema has %d attributes, want 12", len(top))
	}

	required := []string{"image", "image_registry_type"}
	for _, name := range required {
		attribute, ok := top[name]
		if !ok {
			t.Fatalf("attribute %q is missing", name)
		}
		if !attribute.IsRequired() {
			t.Fatalf("attribute %q must be required", name)
		}
	}

	computed := []string{
		"id", "pre_cache_image_id", "image_digest", "status", "message",
		"create_time", "ready_time", "last_used_time", "source_type", "cached_image_size_bytes",
	}
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

	if _, ok := top["cached_image_size_bytes"].(schema.Int64Attribute); !ok {
		t.Fatalf("attribute cached_image_size_bytes is not an int64 attribute")
	}
}

// TestAgsPreCacheImageTaskCreate verifies the create request and the read back
// state, including that the identity is stored before the wait.
func TestAgsPreCacheImageTaskCreate(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	testContext := newAgsPreCacheImageTestContext(t, patches)

	var capturedRequest *agsv20250920.CreatePreCacheImageTaskRequest
	createCalls := 0
	patches.ApplyMethodFunc(
		testContext.client,
		"CreatePreCacheImageTaskWithContext",
		func(_ context.Context, request *agsv20250920.CreatePreCacheImageTaskRequest) (*agsv20250920.CreatePreCacheImageTaskResponse, error) {
			createCalls++
			capturedRequest = request
			return &agsv20250920.CreatePreCacheImageTaskResponse{
				Response: &agsv20250920.CreatePreCacheImageTaskResponseParams{
					PreCacheImageId:   agsString(agsTestPreCacheImageId),
					Image:             agsString(agsTestPreCacheImage),
					ImageRegistryType: agsString(agsTestPreCacheRegistry),
				},
			}, nil
		},
	)
	describeCalls := 0
	patches.ApplyMethodFunc(
		testContext.client,
		"DescribePreCacheImageTaskWithContext",
		func(_ context.Context, request *agsv20250920.DescribePreCacheImageTaskRequest) (*agsv20250920.DescribePreCacheImageTaskResponse, error) {
			describeCalls++
			if request.PreCacheImageId == nil || *request.PreCacheImageId != agsTestPreCacheImageId {
				t.Fatalf("describe request PreCacheImageId is not %q", agsTestPreCacheImageId)
			}
			if request.PreCacheType == nil || *request.PreCacheType != "image" {
				t.Fatalf("describe request PreCacheType is not image")
			}
			return describePreCacheImageResponse(), nil
		},
	)

	createResp := &resource.CreateResponse{State: tfsdk.State{Schema: testContext.schema}}
	testContext.resource.Create(
		context.Background(),
		resource.CreateRequest{Plan: testContext.plan(t, testContext.sampleModel())},
		createResp,
	)
	if createResp.Diagnostics.HasError() {
		t.Fatalf("unexpected error diagnostics: %v", createResp.Diagnostics.Errors())
	}

	if createCalls != 1 {
		t.Fatalf("CreatePreCacheImageTask was called %d times, want 1", createCalls)
	}
	if describeCalls == 0 {
		t.Fatalf("DescribePreCacheImageTask was never called after create")
	}
	if capturedRequest == nil {
		t.Fatalf("CreatePreCacheImageTask was not called")
	}
	if capturedRequest.Image == nil || *capturedRequest.Image != agsTestPreCacheImage {
		t.Fatalf("CreatePreCacheImageTask Image is not %q", agsTestPreCacheImage)
	}
	if capturedRequest.ImageRegistryType == nil || *capturedRequest.ImageRegistryType != agsTestPreCacheRegistry {
		t.Fatalf("CreatePreCacheImageTask ImageRegistryType is not %q", agsTestPreCacheRegistry)
	}

	var state svcags.AgsPreCacheImageTaskResourceModel
	if diags := createResp.State.Get(context.Background(), &state); diags.HasError() {
		t.Fatalf("failed to read created state: %v", diags.Errors())
	}
	assertPreCacheImageState(t, state)
}

// TestAgsPreCacheImageTaskCreateKeepsStateWhenWaitFails verifies the identity is
// stored even when the created task never becomes ready.
func TestAgsPreCacheImageTaskCreateKeepsStateWhenWaitFails(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	testContext := newAgsPreCacheImageTestContext(t, patches)

	patches.ApplyMethodFunc(
		testContext.client,
		"CreatePreCacheImageTaskWithContext",
		func(_ context.Context, request *agsv20250920.CreatePreCacheImageTaskRequest) (*agsv20250920.CreatePreCacheImageTaskResponse, error) {
			return &agsv20250920.CreatePreCacheImageTaskResponse{
				Response: &agsv20250920.CreatePreCacheImageTaskResponseParams{
					PreCacheImageId: agsString(agsTestPreCacheImageId),
				},
			}, nil
		},
	)
	patches.ApplyMethodFunc(
		testContext.client,
		"DescribePreCacheImageTaskWithContext",
		func(_ context.Context, request *agsv20250920.DescribePreCacheImageTaskRequest) (*agsv20250920.DescribePreCacheImageTaskResponse, error) {
			return nil, &mockSdkError{message: "describe failed"}
		},
	)

	createResp := &resource.CreateResponse{State: tfsdk.State{Schema: testContext.schema}}
	testContext.resource.Create(
		context.Background(),
		resource.CreateRequest{Plan: testContext.plan(t, testContext.sampleModel())},
		createResp,
	)
	if !createResp.Diagnostics.HasError() {
		t.Fatalf("expected an error diagnostic when the wait fails")
	}

	// The identity must be in state so `terraform destroy` can clean up.
	var state svcags.AgsPreCacheImageTaskResourceModel
	if diags := createResp.State.Get(context.Background(), &state); diags.HasError() {
		t.Fatalf("failed to read state after a failed wait: %v", diags.Errors())
	}
	if state.Id.ValueString() != agsTestPreCacheImageId || state.PreCacheImageId.ValueString() != agsTestPreCacheImageId {
		t.Fatalf("identity was not stored before waiting")
	}
}

// TestAgsPreCacheImageTaskCreateEmptyId verifies an empty id is reported.
func TestAgsPreCacheImageTaskCreateEmptyId(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	testContext := newAgsPreCacheImageTestContext(t, patches)

	patches.ApplyMethodFunc(
		testContext.client,
		"CreatePreCacheImageTaskWithContext",
		func(_ context.Context, request *agsv20250920.CreatePreCacheImageTaskRequest) (*agsv20250920.CreatePreCacheImageTaskResponse, error) {
			return &agsv20250920.CreatePreCacheImageTaskResponse{
				Response: &agsv20250920.CreatePreCacheImageTaskResponseParams{},
			}, nil
		},
	)

	createResp := &resource.CreateResponse{State: tfsdk.State{Schema: testContext.schema}}
	testContext.resource.Create(
		context.Background(),
		resource.CreateRequest{Plan: testContext.plan(t, testContext.sampleModel())},
		createResp,
	)
	if !createResp.Diagnostics.HasError() {
		t.Fatalf("expected an error diagnostic when the created id is empty")
	}
}

// TestAgsPreCacheImageTaskCreateApiError verifies a create API failure is
// reported and no state is written.
func TestAgsPreCacheImageTaskCreateApiError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	testContext := newAgsPreCacheImageTestContext(t, patches)

	patches.ApplyMethodFunc(
		testContext.client,
		"CreatePreCacheImageTaskWithContext",
		func(_ context.Context, request *agsv20250920.CreatePreCacheImageTaskRequest) (*agsv20250920.CreatePreCacheImageTaskResponse, error) {
			return nil, &mockSdkError{message: "create failed"}
		},
	)

	createResp := &resource.CreateResponse{State: tfsdk.State{Schema: testContext.schema}}
	testContext.resource.Create(
		context.Background(),
		resource.CreateRequest{Plan: testContext.plan(t, testContext.sampleModel())},
		createResp,
	)
	if !createResp.Diagnostics.HasError() {
		t.Fatalf("expected an error diagnostic when the create API fails")
	}
	if !createResp.State.Raw.IsNull() {
		t.Fatalf("state must stay empty when creation fails")
	}
}

// TestAgsPreCacheImageTaskRead verifies every leaf attribute is refreshed.
func TestAgsPreCacheImageTaskRead(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	testContext := newAgsPreCacheImageTestContext(t, patches)

	patches.ApplyMethodFunc(
		testContext.client,
		"DescribePreCacheImageTaskWithContext",
		func(_ context.Context, request *agsv20250920.DescribePreCacheImageTaskRequest) (*agsv20250920.DescribePreCacheImageTaskResponse, error) {
			return describePreCacheImageResponse(), nil
		},
	)

	model := testContext.createdModel()
	model.Status = types.StringNull()
	model.ImageDigest = types.StringNull()

	readResp := &resource.ReadResponse{State: tfsdk.State{Schema: testContext.schema}}
	testContext.resource.Read(
		context.Background(),
		resource.ReadRequest{State: testContext.state(t, model)},
		readResp,
	)
	if readResp.Diagnostics.HasError() {
		t.Fatalf("unexpected error diagnostics: %v", readResp.Diagnostics.Errors())
	}

	var state svcags.AgsPreCacheImageTaskResourceModel
	if diags := readResp.State.Get(context.Background(), &state); diags.HasError() {
		t.Fatalf("failed to read refreshed state: %v", diags.Errors())
	}
	assertPreCacheImageState(t, state)
}

// TestAgsPreCacheImageTaskReadNullableFields verifies empty cloud fields become
// null instead of empty strings.
func TestAgsPreCacheImageTaskReadNullableFields(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	testContext := newAgsPreCacheImageTestContext(t, patches)

	patches.ApplyMethodFunc(
		testContext.client,
		"DescribePreCacheImageTaskWithContext",
		func(_ context.Context, request *agsv20250920.DescribePreCacheImageTaskRequest) (*agsv20250920.DescribePreCacheImageTaskResponse, error) {
			return &agsv20250920.DescribePreCacheImageTaskResponse{
				Response: &agsv20250920.DescribePreCacheImageTaskResponseParams{
					PreCacheImageId: agsString(agsTestPreCacheImageId),
					Image:           agsString(agsTestPreCacheImage),
				},
			}, nil
		},
	)

	readResp := &resource.ReadResponse{State: tfsdk.State{Schema: testContext.schema}}
	testContext.resource.Read(
		context.Background(),
		resource.ReadRequest{State: testContext.state(t, testContext.createdModel())},
		readResp,
	)
	if readResp.Diagnostics.HasError() {
		t.Fatalf("unexpected error diagnostics: %v", readResp.Diagnostics.Errors())
	}

	var state svcags.AgsPreCacheImageTaskResourceModel
	if diags := readResp.State.Get(context.Background(), &state); diags.HasError() {
		t.Fatalf("failed to read refreshed state: %v", diags.Errors())
	}

	nullable := map[string]types.String{
		"image_registry_type": state.ImageRegistryType,
		"image_digest":        state.ImageDigest,
		"status":              state.Status,
		"message":             state.Message,
		"create_time":         state.CreateTime,
		"ready_time":          state.ReadyTime,
		"last_used_time":      state.LastUsedTime,
		"source_type":         state.SourceType,
	}
	for name, value := range nullable {
		if !value.IsNull() {
			t.Fatalf("attribute %q is %q, want null", name, value.ValueString())
		}
	}
	if !state.CachedImageSizeBytes.IsNull() {
		t.Fatalf("cached_image_size_bytes is %d, want null", state.CachedImageSizeBytes.ValueInt64())
	}
}

// TestAgsPreCacheImageTaskReadNotFound verifies a deleted task is removed from
// state.
func TestAgsPreCacheImageTaskReadNotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	testContext := newAgsPreCacheImageTestContext(t, patches)

	patches.ApplyMethodFunc(
		testContext.client,
		"DescribePreCacheImageTaskWithContext",
		func(_ context.Context, request *agsv20250920.DescribePreCacheImageTaskRequest) (*agsv20250920.DescribePreCacheImageTaskResponse, error) {
			return &agsv20250920.DescribePreCacheImageTaskResponse{
				Response: &agsv20250920.DescribePreCacheImageTaskResponseParams{},
			}, nil
		},
	)

	readResp := &resource.ReadResponse{State: tfsdk.State{Schema: testContext.schema}}
	testContext.resource.Read(
		context.Background(),
		resource.ReadRequest{State: testContext.state(t, testContext.createdModel())},
		readResp,
	)
	if readResp.Diagnostics.HasError() {
		t.Fatalf("unexpected error diagnostics: %v", readResp.Diagnostics.Errors())
	}
	if !readResp.State.Raw.IsNull() {
		t.Fatalf("state must be removed when the pre-cache image task does not exist")
	}
}

// TestAgsPreCacheImageTaskReadApiError verifies an API failure never drops
// state.
func TestAgsPreCacheImageTaskReadApiError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	testContext := newAgsPreCacheImageTestContext(t, patches)

	patches.ApplyMethodFunc(
		testContext.client,
		"DescribePreCacheImageTaskWithContext",
		func(_ context.Context, request *agsv20250920.DescribePreCacheImageTaskRequest) (*agsv20250920.DescribePreCacheImageTaskResponse, error) {
			return nil, &mockSdkError{message: "describe failed"}
		},
	)

	readResp := &resource.ReadResponse{State: tfsdk.State{Schema: testContext.schema}}
	testContext.resource.Read(
		context.Background(),
		resource.ReadRequest{State: testContext.state(t, testContext.createdModel())},
		readResp,
	)
	if !readResp.Diagnostics.HasError() {
		t.Fatalf("expected an error diagnostic when the describe API fails")
	}
	if readResp.State.Raw.IsNull() {
		t.Fatalf("state must not be removed when the describe API fails")
	}
}

// TestAgsPreCacheImageTaskReadEmptyResponse verifies an empty API response is
// an error and never removes state.
func TestAgsPreCacheImageTaskReadEmptyResponse(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	testContext := newAgsPreCacheImageTestContext(t, patches)

	patches.ApplyMethodFunc(
		testContext.client,
		"DescribePreCacheImageTaskWithContext",
		func(_ context.Context, request *agsv20250920.DescribePreCacheImageTaskRequest) (*agsv20250920.DescribePreCacheImageTaskResponse, error) {
			return nil, nil
		},
	)

	readResp := &resource.ReadResponse{State: tfsdk.State{Schema: testContext.schema}}
	testContext.resource.Read(
		context.Background(),
		resource.ReadRequest{State: testContext.state(t, testContext.createdModel())},
		readResp,
	)
	if !readResp.Diagnostics.HasError() {
		t.Fatalf("expected an error diagnostic for an empty describe response")
	}
	if readResp.State.Raw.IsNull() {
		t.Fatalf("state must not be removed for an empty describe response")
	}
}

// TestAgsPreCacheImageTaskDelete verifies the delete request carries the
// identity and the descriptive fields read from the cloud.
func TestAgsPreCacheImageTaskDelete(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	testContext := newAgsPreCacheImageTestContext(t, patches)

	var capturedRequest *agsv20250920.DeletePreCacheImageRequest
	deleteCalls := 0
	patches.ApplyMethodFunc(
		testContext.client,
		"DeletePreCacheImageWithContext",
		func(_ context.Context, request *agsv20250920.DeletePreCacheImageRequest) (*agsv20250920.DeletePreCacheImageResponse, error) {
			deleteCalls++
			capturedRequest = request
			return &agsv20250920.DeletePreCacheImageResponse{
				Response: &agsv20250920.DeletePreCacheImageResponseParams{RequestId: agsString("request-id")},
			}, nil
		},
	)
	describeCalls := 0
	patches.ApplyMethodFunc(
		testContext.client,
		"DescribePreCacheImageTaskWithContext",
		func(_ context.Context, request *agsv20250920.DescribePreCacheImageTaskRequest) (*agsv20250920.DescribePreCacheImageTaskResponse, error) {
			describeCalls++
			// The first describe locates the item, the following ones confirm it
			// is gone.
			if describeCalls == 1 {
				return describePreCacheImageResponse(), nil
			}
			return &agsv20250920.DescribePreCacheImageTaskResponse{
				Response: &agsv20250920.DescribePreCacheImageTaskResponseParams{},
			}, nil
		},
	)

	deleteResp := &resource.DeleteResponse{State: tfsdk.State{Schema: testContext.schema}}
	testContext.resource.Delete(
		context.Background(),
		resource.DeleteRequest{State: testContext.state(t, testContext.createdModel())},
		deleteResp,
	)
	if deleteResp.Diagnostics.HasError() {
		t.Fatalf("unexpected error diagnostics: %v", deleteResp.Diagnostics.Errors())
	}
	if deleteCalls != 1 {
		t.Fatalf("DeletePreCacheImage was called %d times, want 1", deleteCalls)
	}
	if capturedRequest == nil {
		t.Fatalf("DeletePreCacheImage was not called")
	}
	if capturedRequest.PreCacheImageId == nil || *capturedRequest.PreCacheImageId != agsTestPreCacheImageId {
		t.Fatalf("DeletePreCacheImage PreCacheImageId is not %q", agsTestPreCacheImageId)
	}
	if capturedRequest.Image == nil || *capturedRequest.Image != agsTestPreCacheImage {
		t.Fatalf("DeletePreCacheImage Image is not %q", agsTestPreCacheImage)
	}
	if capturedRequest.ImageDigest == nil || *capturedRequest.ImageDigest != agsTestPreCacheDigest {
		t.Fatalf("DeletePreCacheImage ImageDigest is not %q", agsTestPreCacheDigest)
	}
	if capturedRequest.ImageRegistryType == nil || *capturedRequest.ImageRegistryType != agsTestPreCacheRegistry {
		t.Fatalf("DeletePreCacheImage ImageRegistryType is not %q", agsTestPreCacheRegistry)
	}
}

// TestAgsPreCacheImageTaskDeleteAlreadyGone verifies a repeated destroy is
// idempotent when the item no longer exists.
func TestAgsPreCacheImageTaskDeleteAlreadyGone(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	testContext := newAgsPreCacheImageTestContext(t, patches)

	patches.ApplyMethodFunc(
		testContext.client,
		"DescribePreCacheImageTaskWithContext",
		func(_ context.Context, request *agsv20250920.DescribePreCacheImageTaskRequest) (*agsv20250920.DescribePreCacheImageTaskResponse, error) {
			return &agsv20250920.DescribePreCacheImageTaskResponse{
				Response: &agsv20250920.DescribePreCacheImageTaskResponseParams{},
			}, nil
		},
	)
	deleteCalls := 0
	patches.ApplyMethodFunc(
		testContext.client,
		"DeletePreCacheImageWithContext",
		func(_ context.Context, request *agsv20250920.DeletePreCacheImageRequest) (*agsv20250920.DeletePreCacheImageResponse, error) {
			deleteCalls++
			return nil, &mockSdkError{message: "should not be called"}
		},
	)

	deleteResp := &resource.DeleteResponse{State: tfsdk.State{Schema: testContext.schema}}
	testContext.resource.Delete(
		context.Background(),
		resource.DeleteRequest{State: testContext.state(t, testContext.createdModel())},
		deleteResp,
	)
	if deleteResp.Diagnostics.HasError() {
		t.Fatalf("unexpected error diagnostics: %v", deleteResp.Diagnostics.Errors())
	}
	if deleteCalls != 0 {
		t.Fatalf("DeletePreCacheImage was called %d times for an item that is already gone", deleteCalls)
	}
}

// TestAgsPreCacheImageTaskDeleteApiError verifies a delete failure keeps an
// error when the item still exists.
func TestAgsPreCacheImageTaskDeleteApiError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	testContext := newAgsPreCacheImageTestContext(t, patches)

	patches.ApplyMethodFunc(
		testContext.client,
		"DescribePreCacheImageTaskWithContext",
		func(_ context.Context, request *agsv20250920.DescribePreCacheImageTaskRequest) (*agsv20250920.DescribePreCacheImageTaskResponse, error) {
			return describePreCacheImageResponse(), nil
		},
	)
	patches.ApplyMethodFunc(
		testContext.client,
		"DeletePreCacheImageWithContext",
		func(_ context.Context, request *agsv20250920.DeletePreCacheImageRequest) (*agsv20250920.DeletePreCacheImageResponse, error) {
			return nil, &mockSdkError{message: "delete failed"}
		},
	)

	deleteResp := &resource.DeleteResponse{State: tfsdk.State{Schema: testContext.schema}}
	testContext.resource.Delete(
		context.Background(),
		resource.DeleteRequest{State: testContext.state(t, testContext.createdModel())},
		deleteResp,
	)
	if !deleteResp.Diagnostics.HasError() {
		t.Fatalf("expected an error diagnostic when the delete API fails")
	}
}

// TestAgsPreCacheImageTaskDeleteToleratesAlreadyDeleted verifies a delete
// failure is tolerated when the item is confirmed gone afterwards.
func TestAgsPreCacheImageTaskDeleteToleratesAlreadyDeleted(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	testContext := newAgsPreCacheImageTestContext(t, patches)

	describeCalls := 0
	patches.ApplyMethodFunc(
		testContext.client,
		"DescribePreCacheImageTaskWithContext",
		func(_ context.Context, request *agsv20250920.DescribePreCacheImageTaskRequest) (*agsv20250920.DescribePreCacheImageTaskResponse, error) {
			describeCalls++
			if describeCalls == 1 {
				return describePreCacheImageResponse(), nil
			}
			return &agsv20250920.DescribePreCacheImageTaskResponse{
				Response: &agsv20250920.DescribePreCacheImageTaskResponseParams{},
			}, nil
		},
	)
	patches.ApplyMethodFunc(
		testContext.client,
		"DeletePreCacheImageWithContext",
		func(_ context.Context, request *agsv20250920.DeletePreCacheImageRequest) (*agsv20250920.DeletePreCacheImageResponse, error) {
			return nil, &mockSdkError{message: "already deleted"}
		},
	)

	deleteResp := &resource.DeleteResponse{State: tfsdk.State{Schema: testContext.schema}}
	testContext.resource.Delete(
		context.Background(),
		resource.DeleteRequest{State: testContext.state(t, testContext.createdModel())},
		deleteResp,
	)
	if deleteResp.Diagnostics.HasError() {
		t.Fatalf("unexpected error diagnostics for an already deleted item: %v", deleteResp.Diagnostics.Errors())
	}
}

// TestAgsPreCacheImageTaskUpdateNotSupported verifies the resource has no
// in-place update: every change is expressed as a replacement.
func TestAgsPreCacheImageTaskUpdateNotSupported(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	testContext := newAgsPreCacheImageTestContext(t, patches)

	// The resource still exposes Update (embedded WithNoUpdate), so a plan that
	// reaches it reports an explicit error instead of silently ignoring the
	// change. Every configurable attribute is RequiresReplace anyway.
	updateResp := &resource.UpdateResponse{State: tfsdk.State{Schema: testContext.schema}}
	testContext.resource.Update(
		context.Background(),
		resource.UpdateRequest{
			Plan:  testContext.plan(t, testContext.createdModel()),
			State: testContext.state(t, testContext.createdModel()),
		},
		updateResp,
	)
	if !updateResp.Diagnostics.HasError() {
		t.Fatalf("expected an error diagnostic because in-place update is not supported")
	}
}

// TestAgsPreCacheImageTaskImportState verifies the id only import path.
func TestAgsPreCacheImageTaskImportState(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	testContext := newAgsPreCacheImageTestContext(t, patches)

	importer, ok := testContext.resource.(resource.ResourceWithImportState)
	if !ok {
		t.Fatalf("resource does not implement resource.ResourceWithImportState")
	}

	importResp := &resource.ImportStateResponse{State: tfsdk.State{Schema: testContext.schema}}
	importer.ImportState(
		context.Background(),
		resource.ImportStateRequest{ID: agsTestPreCacheImageId},
		importResp,
	)
	if importResp.Diagnostics.HasError() {
		t.Fatalf("unexpected error diagnostics: %v", importResp.Diagnostics.Errors())
	}

	var state svcags.AgsPreCacheImageTaskResourceModel
	if diags := importResp.State.Get(context.Background(), &state); diags.HasError() {
		t.Fatalf("failed to read imported state: %v", diags.Errors())
	}
	if state.Id.ValueString() != agsTestPreCacheImageId {
		t.Fatalf("imported id is %q, want %q", state.Id.ValueString(), agsTestPreCacheImageId)
	}
	if state.PreCacheImageId.ValueString() != agsTestPreCacheImageId {
		t.Fatalf("imported pre_cache_image_id is %q, want %q", state.PreCacheImageId.ValueString(), agsTestPreCacheImageId)
	}
}

// TestAgsPreCacheImageTaskImportStateEmptyId verifies an empty import id is
// rejected.
func TestAgsPreCacheImageTaskImportStateEmptyId(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	testContext := newAgsPreCacheImageTestContext(t, patches)

	importer := testContext.resource.(resource.ResourceWithImportState)
	importResp := &resource.ImportStateResponse{State: tfsdk.State{Schema: testContext.schema}}
	importer.ImportState(
		context.Background(),
		resource.ImportStateRequest{ID: ""},
		importResp,
	)
	if !importResp.Diagnostics.HasError() {
		t.Fatalf("expected an error diagnostic for an empty import id")
	}
}
