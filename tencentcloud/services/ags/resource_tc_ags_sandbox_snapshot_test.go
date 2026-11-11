package ags

import (
	"context"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	ags "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ags/v20250920"
	sdkcommon "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/sharedmeta"
)

// --- Metadata / Schema / Factory / Configure ---

// TestAgsSandboxSnapshotResource_Metadata validates the type name.
func TestAgsSandboxSnapshotResource_Metadata(t *testing.T) {
	r := &AgsSandboxSnapshotResource{}
	req := resource.MetadataRequest{}
	resp := &resource.MetadataResponse{}

	r.Metadata(context.Background(), req, resp)

	if resp.TypeName != "tencentcloud_ags_sandbox_snapshot" {
		t.Errorf("Metadata.TypeName = %q, want %q", resp.TypeName, "tencentcloud_ags_sandbox_snapshot")
	}
}

// TestAgsSandboxSnapshotResource_Schema validates required/optional/computed
// flags on every top-level attribute.
func TestAgsSandboxSnapshotResource_Schema(t *testing.T) {
	r := &AgsSandboxSnapshotResource{}
	req := resource.SchemaRequest{}
	resp := &resource.SchemaResponse{}

	r.Schema(context.Background(), req, resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Schema returned diagnostics: %v", resp.Diagnostics)
	}

	type expect struct {
		required, optional, computed bool
	}
	cases := map[string]expect{
		"id":            {computed: true},
		"instance_id":   {optional: true},
		"snapshot_name": {optional: true},
		"memory":        {optional: true},
		"tags":          {optional: true},
		"client_token":  {optional: true},
		"status":        {computed: true},
		"status_reason": {computed: true},
		"create_time":   {computed: true},
		"update_time":   {computed: true},
	}

	for name, want := range cases {
		attr, ok := resp.Schema.Attributes[name]
		if !ok {
			t.Errorf("attribute %q missing from schema", name)
			continue
		}
		var required, optional, computed bool
		switch a := attr.(type) {
		case schema.StringAttribute:
			required, optional, computed = a.Required, a.Optional, a.Computed
		case schema.BoolAttribute:
			required, optional, computed = a.Required, a.Optional, a.Computed
		case schema.ListNestedAttribute:
			required, optional, computed = a.Required, a.Optional, a.Computed
		default:
			t.Errorf("attribute %q unexpected type %T", name, attr)
			continue
		}
		if required != want.required || optional != want.optional || computed != want.computed {
			t.Errorf("attribute %q flags (r=%v,o=%v,c=%v), want (r=%v,o=%v,c=%v)",
				name, required, optional, computed, want.required, want.optional, want.computed)
		}
	}

	// Every writable attribute must force replacement because there is no
	// Modify API for this resource.
	for _, name := range []string{"instance_id", "snapshot_name", "memory", "tags", "client_token"} {
		attr := resp.Schema.Attributes[name]
		if len(attr.GetPlanModifiers()) == 0 {
			t.Errorf("attribute %q should have a RequiresReplace plan modifier", name)
		}
	}
}

// TestNewAgsSandboxSnapshotResource validates the factory.
func TestNewAgsSandboxSnapshotResource(t *testing.T) {
	got := NewAgsSandboxSnapshotResource()
	if got == nil {
		t.Fatal("NewAgsSandboxSnapshotResource returned nil")
	}
	if _, ok := got.(*AgsSandboxSnapshotResource); !ok {
		t.Fatalf("factory returned %T, want *AgsSandboxSnapshotResource", got)
	}
}

// TestAgsSandboxSnapshotResource_ConfigureNilProviderData ensures the
// pre-Configure invocation pattern (nil ProviderData) is silently no-op.
func TestAgsSandboxSnapshotResource_ConfigureNilProviderData(t *testing.T) {
	r := &AgsSandboxSnapshotResource{}
	req := resource.ConfigureRequest{}
	resp := &resource.ConfigureResponse{}

	r.Configure(context.Background(), req, resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Configure with nil ProviderData returned diagnostics: %v", resp.Diagnostics)
	}
}

// --- Tag conversion helpers ---

// TestTagsToSdk verifies the model-to-SDK tag conversion.
func TestTagsToSdk(t *testing.T) {
	tags := []AgsTagModel{
		{Key: types.StringValue("k1"), Value: types.StringValue("v1")},
		{Key: types.StringValue("k2"), Value: types.StringNull()},
	}
	got := tagsToSdk(tags)
	if len(got) != 2 {
		t.Fatalf("tagsToSdk len = %d, want 2", len(got))
	}
	if got[0].Key == nil || *got[0].Key != "k1" || got[0].Value == nil || *got[0].Value != "v1" {
		t.Errorf("tagsToSdk[0] = %+v, want key k1 value v1", *got[0])
	}
	if got[1].Key == nil || *got[1].Key != "k2" || got[1].Value != nil {
		t.Errorf("tagsToSdk[1] should map null value to nil pointer")
	}

	if got := tagsToSdk(nil); len(got) != 0 {
		t.Errorf("tagsToSdk(nil) len = %d, want 0", len(got))
	}
}

// TestTagsFromSdk verifies the SDK-to-model tag conversion.
func TestTagsFromSdk(t *testing.T) {
	sdkTags := []*ags.Tag{
		{Key: sdkcommon.StringPtr("k1"), Value: sdkcommon.StringPtr("v1")},
		nil,
	}
	got := tagsFromSdk(sdkTags)
	if len(got) != 1 {
		t.Fatalf("tagsFromSdk len = %d, want 1 (nil entries skipped)", len(got))
	}
	if got[0].Key.ValueString() != "k1" {
		t.Errorf("tagsFromSdk[0].Key = %v, want k1", got[0].Key)
	}
	if got[0].Value.ValueString() != "v1" {
		t.Errorf("tagsFromSdk[0].Value = %v, want v1", got[0].Value)
	}
	if tagsFromSdk(nil) != nil {
		t.Error("tagsFromSdk(nil) should return nil")
	}
	if tagsFromSdk([]*ags.Tag{}) != nil {
		t.Error("tagsFromSdk(empty) should return nil")
	}
}

// --- read / refreshModel with mocked cloud API ---

// mockAgsClient wires an AgsSandboxSnapshotResource whose underlying client
// returns a patched ags SDK client for the duration of the test.
func mockAgsClient(t *testing.T) (*gomonkey.Patches, *AgsSandboxSnapshotResource, *ags.Client) {
	t.Helper()
	patches := gomonkey.NewPatches()
	client := &connectivity.TencentCloudClient{Region: "ap-guangzhou"}
	agsClient := &ags.Client{}
	patches.ApplyMethodReturn(client, "UseAgsClient", agsClient)
	r := &AgsSandboxSnapshotResource{}
	req := resource.ConfigureRequest{ProviderData: &sharedmeta.ProviderMeta{Client: client}}
	resp := &resource.ConfigureResponse{}
	r.Configure(context.Background(), req, resp)
	if resp.Diagnostics.HasError() {
		patches.Reset()
		t.Fatalf("Configure returned diagnostics: %v", resp.Diagnostics)
	}
	return patches, r, agsClient
}

// TestAgsSandboxSnapshotReadFound verifies read returns the snapshot matching
// the requested id.
func TestAgsSandboxSnapshotReadFound(t *testing.T) {
	patches, r, agsClient := mockAgsClient(t)
	defer patches.Reset()

	var captured *ags.DescribeSandboxSnapshotListRequest
	patches.ApplyMethodFunc(agsClient, "DescribeSandboxSnapshotListWithContext",
		func(ctx context.Context, request *ags.DescribeSandboxSnapshotListRequest) (*ags.DescribeSandboxSnapshotListResponse, error) {
			captured = request
			resp := &ags.DescribeSandboxSnapshotListResponse{}
			resp.Response = &ags.DescribeSandboxSnapshotListResponseParams{
				SnapshotSet: []*ags.SandboxSnapshot{
					{SnapshotId: sdkcommon.StringPtr("snap-other")},
					{
						SnapshotId:   sdkcommon.StringPtr("snap-1"),
						SnapshotName: sdkcommon.StringPtr("demo"),
					},
				},
			}
			return resp, nil
		})

	item, err := r.read(context.Background(), "snap-1")
	if err != nil {
		t.Fatalf("read returned error: %v", err)
	}
	if item == nil || item.SnapshotName == nil || *item.SnapshotName != "demo" {
		t.Fatalf("read returned wrong snapshot: %+v", item)
	}
	if len(captured.SnapshotIds) != 1 || *captured.SnapshotIds[0] != "snap-1" {
		t.Errorf("DescribeSandboxSnapshotList.SnapshotIds = %v, want [snap-1]", captured.SnapshotIds)
	}
}

// TestAgsSandboxSnapshotReadNotFound verifies read returns nil,nil when no
// snapshot matches, which Read maps to "removed from state".
func TestAgsSandboxSnapshotReadNotFound(t *testing.T) {
	patches, r, agsClient := mockAgsClient(t)
	defer patches.Reset()

	patches.ApplyMethodFunc(agsClient, "DescribeSandboxSnapshotListWithContext",
		func(ctx context.Context, request *ags.DescribeSandboxSnapshotListRequest) (*ags.DescribeSandboxSnapshotListResponse, error) {
			resp := &ags.DescribeSandboxSnapshotListResponse{}
			resp.Response = &ags.DescribeSandboxSnapshotListResponseParams{
				SnapshotSet: []*ags.SandboxSnapshot{},
			}
			return resp, nil
		})

	item, err := r.read(context.Background(), "snap-gone")
	if err != nil {
		t.Fatalf("read returned error: %v", err)
	}
	if item != nil {
		t.Errorf("read should return nil for a non-matching snapshot, got %+v", item)
	}
}

// TestAgsSandboxSnapshotReadEmptyResponse verifies a nil Response payload is
// an API failure, not a not-found signal.
func TestAgsSandboxSnapshotReadEmptyResponse(t *testing.T) {
	patches, r, agsClient := mockAgsClient(t)
	defer patches.Reset()

	patches.ApplyMethodFunc(agsClient, "DescribeSandboxSnapshotListWithContext",
		func(ctx context.Context, request *ags.DescribeSandboxSnapshotListRequest) (*ags.DescribeSandboxSnapshotListResponse, error) {
			return &ags.DescribeSandboxSnapshotListResponse{}, nil
		})

	_, err := r.read(context.Background(), "snap-1")
	if err == nil {
		t.Fatal("read should fail on an empty Response payload")
	}
}

// TestAgsSandboxSnapshotReadApiError verifies API errors propagate.
func TestAgsSandboxSnapshotReadApiError(t *testing.T) {
	patches, r, agsClient := mockAgsClient(t)
	defer patches.Reset()

	patches.ApplyMethodFunc(agsClient, "DescribeSandboxSnapshotListWithContext",
		func(ctx context.Context, request *ags.DescribeSandboxSnapshotListRequest) (*ags.DescribeSandboxSnapshotListResponse, error) {
			return nil, context.DeadlineExceeded
		})

	_, err := r.read(context.Background(), "snap-1")
	if err == nil {
		t.Fatal("read should propagate the API error")
	}
}

// TestAgsSandboxSnapshotRefreshModel verifies every computed leaf field is
// refreshed and the non-readable config fields keep their state values.
func TestAgsSandboxSnapshotRefreshModel(t *testing.T) {
	r := &AgsSandboxSnapshotResource{}
	m := AgsSandboxSnapshotResourceModel{
		Id:           types.StringValue("snap-1"),
		InstanceId:   types.StringValue("inst-1"),
		SnapshotName: types.StringValue("old"),
		Memory:       types.BoolValue(true),
		ClientToken:  types.StringValue("token-1"),
	}
	item := &ags.SandboxSnapshot{
		SnapshotId:     sdkcommon.StringPtr("snap-1"),
		SnapshotName:   sdkcommon.StringPtr("demo"),
		SourceInstanceId: sdkcommon.StringPtr("inst-1"),
		MemoryIncluded: sdkcommon.BoolPtr(true),
		Status:         sdkcommon.StringPtr("ACTIVE"),
		StatusReason:   sdkcommon.StringPtr("ok"),
		Tags: []*ags.Tag{
			{Key: sdkcommon.StringPtr("k1"), Value: sdkcommon.StringPtr("v1")},
		},
		CreateTime: sdkcommon.StringPtr("2025-01-01T00:00:00Z"),
		UpdateTime: sdkcommon.StringPtr("2025-01-02T00:00:00Z"),
	}

	r.refreshModel(&m, item)

	if m.Id.ValueString() != "snap-1" {
		t.Errorf("Id = %q, want snap-1", m.Id.ValueString())
	}
	if m.SnapshotName.ValueString() != "demo" {
		t.Errorf("SnapshotName = %q, want demo", m.SnapshotName.ValueString())
	}
	if m.Status.ValueString() != "ACTIVE" {
		t.Errorf("Status = %q, want ACTIVE", m.Status.ValueString())
	}
	if m.StatusReason.ValueString() != "ok" {
		t.Errorf("StatusReason = %q, want ok", m.StatusReason.ValueString())
	}
	if m.CreateTime.ValueString() != "2025-01-01T00:00:00Z" {
		t.Errorf("CreateTime = %q, unexpected", m.CreateTime.ValueString())
	}
	if m.UpdateTime.ValueString() != "2025-01-02T00:00:00Z" {
		t.Errorf("UpdateTime = %q, unexpected", m.UpdateTime.ValueString())
	}
	if len(m.Tags) != 1 || m.Tags[0].Key.ValueString() != "k1" || m.Tags[0].Value.ValueString() != "v1" {
		t.Errorf("Tags = %+v, want one tag k1=v1", m.Tags)
	}
	// Fields without read-back capability keep their state values.
	if m.InstanceId.ValueString() != "inst-1" {
		t.Errorf("InstanceId = %q, want inst-1 preserved", m.InstanceId.ValueString())
	}
	if m.Memory.ValueBool() != true {
		t.Errorf("Memory = %v, want true preserved", m.Memory.ValueBool())
	}
	if m.ClientToken.ValueString() != "token-1" {
		t.Errorf("ClientToken = %q, want token-1 preserved", m.ClientToken.ValueString())
	}
}

// TestAgsSandboxSnapshotRefreshModelNil verifies a nil snapshot item is a no-op.
func TestAgsSandboxSnapshotRefreshModelNil(t *testing.T) {
	r := &AgsSandboxSnapshotResource{}
	m := AgsSandboxSnapshotResourceModel{Id: types.StringValue("snap-1")}
	r.refreshModel(&m, nil)
	if m.Id.ValueString() != "snap-1" {
		t.Errorf("refreshModel(nil) should not modify the model")
	}
}
