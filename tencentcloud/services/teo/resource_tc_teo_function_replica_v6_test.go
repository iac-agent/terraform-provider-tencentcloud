package teo_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/assert"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/sharedmeta"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/teo"
)

// go test ./tencentcloud/services/teo/ -run "TestTeoFunctionReplicaV6" -v -count=1 -gcflags="all=-l"

func ptrStringFrV6(s string) *string { return &s }

func ptrInt64FrV6(i int64) *int64 { return &i }

// newTeoFunctionReplicaV6Resource builds the resource and returns it together
// with its framework schema, both used to assemble raw plan/state values.
func newTeoFunctionReplicaV6Resource(t *testing.T) (*teo.TeoFunctionReplicaV6Resource, schema.Schema) {
	t.Helper()
	r, ok := teo.NewTeoFunctionReplicaV6Resource().(*teo.TeoFunctionReplicaV6Resource)
	if !ok {
		t.Fatalf("NewTeoFunctionReplicaV6Resource returned %T, want *teo.TeoFunctionReplicaV6Resource", r)
	}
	schemaResp := &resource.SchemaResponse{}
	r.Schema(context.Background(), resource.SchemaRequest{}, schemaResp)
	if schemaResp.Diagnostics.HasError() {
		t.Fatalf("failed to build resource schema: %v", schemaResp.Diagnostics)
	}
	return r, schemaResp.Schema
}

// teoFunctionReplicaV6Raw builds a tftypes object value for the resource
// model; attributes not present in values default to null.
func teoFunctionReplicaV6Raw(t *testing.T, s schema.Schema, values map[string]tftypes.Value) tftypes.Value {
	t.Helper()
	tfType := s.Type().TerraformType(context.Background())
	objType, ok := tfType.(tftypes.Object)
	if !ok {
		t.Fatalf("schema type is %T, want tftypes.Object", tfType)
	}
	full := make(map[string]tftypes.Value, len(objType.AttributeTypes))
	for name, typ := range objType.AttributeTypes {
		if v, ok := values[name]; ok {
			full[name] = v
		} else {
			full[name] = tftypes.NewValue(typ, nil)
		}
	}
	return tftypes.NewValue(objType, full)
}

func teoFunctionReplicaV6Null(s schema.Schema) tftypes.Value {
	return tftypes.NewValue(s.Type().TerraformType(context.Background()), nil)
}

func frV6Str(v string) tftypes.Value { return tftypes.NewValue(tftypes.String, v) }

func frV6Unknown() tftypes.Value { return tftypes.NewValue(tftypes.String, tftypes.UnknownValue) }

// configureTeoFunctionReplicaV6 injects the mock client into the resource via
// the Configure method.
func configureTeoFunctionReplicaV6(r *teo.TeoFunctionReplicaV6Resource, client *connectivity.TencentCloudClient) {
	meta := &sharedmeta.ProviderMeta{Client: client}
	configureReq := resource.ConfigureRequest{ProviderData: meta}
	configureResp := &resource.ConfigureResponse{}
	r.Configure(context.Background(), configureReq, configureResp)
}

func frV6DescribeReplica() *teov20220901.FunctionReplica {
	return &teov20220901.FunctionReplica{
		FunctionId:  ptrStringFrV6("ef-test456"),
		ReplicaName: ptrStringFrV6("test-replica"),
		Content:     ptrStringFrV6("console.log(123)"),
		Remark:      ptrStringFrV6("test remark"),
		CreatedOn:   ptrStringFrV6("2024-01-01T00:00:00Z"),
		ModifiedOn:  ptrStringFrV6("2024-01-01T00:00:00Z"),
	}
}

func frV6DescribeResponse(items ...*teov20220901.FunctionReplica) *teov20220901.DescribeFunctionReplicasResponse {
	resp := teov20220901.NewDescribeFunctionReplicasResponse()
	resp.Response = &teov20220901.DescribeFunctionReplicasResponseParams{
		TotalCount:       ptrInt64FrV6(int64(len(items))),
		FunctionReplicas: items,
		RequestId:        ptrStringFrV6("fake-request-id"),
	}
	return resp
}

func frV6CreateResponse() *teov20220901.CreateFunctionReplicaResponse {
	resp := teov20220901.NewCreateFunctionReplicaResponse()
	resp.Response = &teov20220901.CreateFunctionReplicaResponseParams{
		RequestId: ptrStringFrV6("fake-request-id"),
	}
	return resp
}

// TestTeoFunctionReplicaV6Resource_Metadata validates the resource type name.
func TestTeoFunctionReplicaV6Resource_Metadata(t *testing.T) {
	r, _ := newTeoFunctionReplicaV6Resource(t)

	req := resource.MetadataRequest{ProviderTypeName: "tencentcloud"}
	resp := &resource.MetadataResponse{}
	r.Metadata(context.Background(), req, resp)

	assert.Equal(t, "tencentcloud_teo_function_replica_v6", resp.TypeName)
}

// TestTeoFunctionReplicaV6Resource_Schema validates required / optional /
// computed flags on every attribute.
func TestTeoFunctionReplicaV6Resource_Schema(t *testing.T) {
	_, s := newTeoFunctionReplicaV6Resource(t)

	attrs := s.Attributes
	assert.Contains(t, attrs, "id")
	assert.Contains(t, attrs, "zone_id")
	assert.Contains(t, attrs, "function_id")
	assert.Contains(t, attrs, "replica_name")
	assert.Contains(t, attrs, "content")
	assert.Contains(t, attrs, "remark")

	idAttr, ok := attrs["id"].(schema.StringAttribute)
	assert.True(t, ok)
	assert.True(t, idAttr.Computed)
	assert.False(t, idAttr.Required)
	assert.False(t, idAttr.Optional)

	for _, name := range []string{"zone_id", "function_id", "replica_name"} {
		attr, ok := attrs[name].(schema.StringAttribute)
		assert.True(t, ok, "%s should be a StringAttribute", name)
		assert.True(t, attr.Required, "%s should be required", name)
		assert.False(t, attr.Optional, "%s should not be optional", name)
		assert.False(t, attr.Computed, "%s should not be computed", name)
	}

	contentAttr, ok := attrs["content"].(schema.StringAttribute)
	assert.True(t, ok)
	assert.True(t, contentAttr.Required)
	assert.False(t, contentAttr.Optional)

	remarkAttr, ok := attrs["remark"].(schema.StringAttribute)
	assert.True(t, ok)
	assert.False(t, remarkAttr.Required)
	assert.True(t, remarkAttr.Optional)
	assert.False(t, remarkAttr.Computed)
}

// TestTeoFunctionReplicaV6Resource_Factory validates the registry factory.
func TestTeoFunctionReplicaV6Resource_Factory(t *testing.T) {
	got := teo.NewTeoFunctionReplicaV6Resource()
	if got == nil {
		t.Fatal("NewTeoFunctionReplicaV6Resource returned nil")
	}
	if _, ok := got.(*teo.TeoFunctionReplicaV6Resource); !ok {
		t.Fatalf("factory returned %T, want *teo.TeoFunctionReplicaV6Resource", got)
	}
}

// TestTeoFunctionReplicaV6Resource_ImplementsInterfaces verifies the resource
// satisfies the framework interfaces it advertises in package-level asserts.
func TestTeoFunctionReplicaV6Resource_ImplementsInterfaces(t *testing.T) {
	var _ resource.ResourceWithConfigure = &teo.TeoFunctionReplicaV6Resource{}
	var _ resource.ResourceWithImportState = &teo.TeoFunctionReplicaV6Resource{}
}

// TestTeoFunctionReplicaV6Resource_ConfigureNilProviderData ensures a
// Configure invocation with nil ProviderData is silently a no-op.
func TestTeoFunctionReplicaV6Resource_ConfigureNilProviderData(t *testing.T) {
	r, _ := newTeoFunctionReplicaV6Resource(t)

	req := resource.ConfigureRequest{}
	resp := &resource.ConfigureResponse{}
	r.Configure(context.Background(), req, resp)

	assert.False(t, resp.Diagnostics.HasError(), "unexpected diagnostics: %v", resp.Diagnostics)
}

// TestTeoFunctionReplicaV6Resource_Create tests a successful create flow: the
// Create API receives all attributes, and the state carries the composite id
// plus every leaf field read back from DescribeFunctionReplicas.
func TestTeoFunctionReplicaV6Resource_Create(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(client, "UseTeoV20220901Client", teoClient)

	var gotCreate *teov20220901.CreateFunctionReplicaRequest
	patches.ApplyMethodFunc(teoClient, "CreateFunctionReplicaWithContext", func(_ context.Context, request *teov20220901.CreateFunctionReplicaRequest) (*teov20220901.CreateFunctionReplicaResponse, error) {
		gotCreate = request
		return frV6CreateResponse(), nil
	})

	patches.ApplyMethodFunc(teoClient, "DescribeFunctionReplicasWithContext", func(_ context.Context, request *teov20220901.DescribeFunctionReplicasRequest) (*teov20220901.DescribeFunctionReplicasResponse, error) {
		assert.Equal(t, "zone-test123", *request.ZoneId)
		assert.Equal(t, "ef-test456", *request.FunctionId)
		assert.Equal(t, "replica-name", *request.Filters[0].Name)
		assert.Equal(t, "test-replica", *request.Filters[0].Values[0])
		return frV6DescribeResponse(frV6DescribeReplica()), nil
	})

	r, s := newTeoFunctionReplicaV6Resource(t)
	configureTeoFunctionReplicaV6(r, client)

	raw := teoFunctionReplicaV6Raw(t, s, map[string]tftypes.Value{
		"zone_id":      frV6Str("zone-test123"),
		"function_id":  frV6Str("ef-test456"),
		"replica_name": frV6Str("test-replica"),
		"content":      frV6Str("console.log(123)"),
		"remark":       frV6Str("test remark"),
		"id":           frV6Unknown(),
	})

	req := resource.CreateRequest{
		Config: tfsdk.Config{Raw: raw, Schema: s},
		Plan:   tfsdk.Plan{Raw: raw, Schema: s},
	}
	resp := &resource.CreateResponse{State: tfsdk.State{Raw: teoFunctionReplicaV6Null(s), Schema: s}}

	r.Create(context.Background(), req, resp)

	assert.False(t, resp.Diagnostics.HasError(), "unexpected diagnostics: %v", resp.Diagnostics)

	assert.NotNil(t, gotCreate)
	assert.Equal(t, "zone-test123", *gotCreate.ZoneId)
	assert.Equal(t, "ef-test456", *gotCreate.FunctionId)
	assert.Equal(t, "test-replica", *gotCreate.ReplicaName)
	assert.Equal(t, "console.log(123)", *gotCreate.Content)
	assert.Equal(t, "test remark", *gotCreate.Remark)

	var model teo.TeoFunctionReplicaV6ResourceModel
	diags := resp.State.Get(context.Background(), &model)
	assert.False(t, diags.HasError(), "state get returned diagnostics: %v", diags)

	assert.Equal(t, "zone-test123#ef-test456#test-replica", model.Id.ValueString())
	assert.Equal(t, "zone-test123", model.ZoneId.ValueString())
	assert.Equal(t, "ef-test456", model.FunctionId.ValueString())
	assert.Equal(t, "test-replica", model.ReplicaName.ValueString())
	assert.Equal(t, "console.log(123)", model.Content.ValueString())
	assert.Equal(t, "test remark", model.Remark.ValueString())
}

// TestTeoFunctionReplicaV6Resource_CreateWithoutRemark verifies a null
// remark is not sent to the Create API and stays null in state.
func TestTeoFunctionReplicaV6Resource_CreateWithoutRemark(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(client, "UseTeoV20220901Client", teoClient)

	var gotCreate *teov20220901.CreateFunctionReplicaRequest
	patches.ApplyMethodFunc(teoClient, "CreateFunctionReplicaWithContext", func(_ context.Context, request *teov20220901.CreateFunctionReplicaRequest) (*teov20220901.CreateFunctionReplicaResponse, error) {
		gotCreate = request
		return frV6CreateResponse(), nil
	})

	replica := frV6DescribeReplica()
	replica.Remark = nil
	patches.ApplyMethodFunc(teoClient, "DescribeFunctionReplicasWithContext", func(_ context.Context, request *teov20220901.DescribeFunctionReplicasRequest) (*teov20220901.DescribeFunctionReplicasResponse, error) {
		return frV6DescribeResponse(replica), nil
	})

	r, s := newTeoFunctionReplicaV6Resource(t)
	configureTeoFunctionReplicaV6(r, client)

	raw := teoFunctionReplicaV6Raw(t, s, map[string]tftypes.Value{
		"zone_id":      frV6Str("zone-test123"),
		"function_id":  frV6Str("ef-test456"),
		"replica_name": frV6Str("test-replica"),
		"content":      frV6Str("console.log(123)"),
		"id":           frV6Unknown(),
	})

	req := resource.CreateRequest{
		Config: tfsdk.Config{Raw: raw, Schema: s},
		Plan:   tfsdk.Plan{Raw: raw, Schema: s},
	}
	resp := &resource.CreateResponse{State: tfsdk.State{Raw: teoFunctionReplicaV6Null(s), Schema: s}}

	r.Create(context.Background(), req, resp)

	assert.False(t, resp.Diagnostics.HasError(), "unexpected diagnostics: %v", resp.Diagnostics)

	assert.NotNil(t, gotCreate)
	assert.Nil(t, gotCreate.Remark, "remark should not be sent when null")

	var model teo.TeoFunctionReplicaV6ResourceModel
	diags := resp.State.Get(context.Background(), &model)
	assert.False(t, diags.HasError(), "state get returned diagnostics: %v", diags)
	assert.True(t, model.Remark.IsNull(), "remark in state should be null, got %v", model.Remark)
}

// TestTeoFunctionReplicaV6Resource_CreateAPIError verifies the create error
// branch surfaces an error diagnostic.
func TestTeoFunctionReplicaV6Resource_CreateAPIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "CreateFunctionReplicaWithContext", func(_ context.Context, request *teov20220901.CreateFunctionReplicaRequest) (*teov20220901.CreateFunctionReplicaResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=InvalidParameter.ReplicaNameInUse, Message=replica already exists")
	})

	r, s := newTeoFunctionReplicaV6Resource(t)
	configureTeoFunctionReplicaV6(r, client)

	raw := teoFunctionReplicaV6Raw(t, s, map[string]tftypes.Value{
		"zone_id":      frV6Str("zone-test123"),
		"function_id":  frV6Str("ef-test456"),
		"replica_name": frV6Str("test-replica"),
		"content":      frV6Str("console.log(123)"),
		"id":           frV6Unknown(),
	})

	req := resource.CreateRequest{
		Config: tfsdk.Config{Raw: raw, Schema: s},
		Plan:   tfsdk.Plan{Raw: raw, Schema: s},
	}
	resp := &resource.CreateResponse{State: tfsdk.State{Raw: teoFunctionReplicaV6Null(s), Schema: s}}

	r.Create(context.Background(), req, resp)

	assert.True(t, resp.Diagnostics.HasError(), "expected error diagnostics")
	assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Error creating TEO function replica")
	assert.Contains(t, resp.Diagnostics.Errors()[0].Detail(), "ReplicaNameInUse")
}

// TestTeoFunctionReplicaV6Resource_CreateClientNotConfigured verifies the
// nil-client guard fails before any API call.
func TestTeoFunctionReplicaV6Resource_CreateClientNotConfigured(t *testing.T) {
	r, s := newTeoFunctionReplicaV6Resource(t)

	raw := teoFunctionReplicaV6Raw(t, s, map[string]tftypes.Value{
		"zone_id":      frV6Str("zone-test123"),
		"function_id":  frV6Str("ef-test456"),
		"replica_name": frV6Str("test-replica"),
		"content":      frV6Str("console.log(123)"),
		"id":           frV6Unknown(),
	})

	req := resource.CreateRequest{
		Config: tfsdk.Config{Raw: raw, Schema: s},
		Plan:   tfsdk.Plan{Raw: raw, Schema: s},
	}
	resp := &resource.CreateResponse{State: tfsdk.State{Raw: teoFunctionReplicaV6Null(s), Schema: s}}

	r.Create(context.Background(), req, resp)

	assert.True(t, resp.Diagnostics.HasError(), "expected error diagnostics")
	assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Provider not configured")
}

// TestTeoFunctionReplicaV6Resource_Read tests a successful read: all leaf
// fields are refreshed from the Describe response, remark null round-trips.
func TestTeoFunctionReplicaV6Resource_Read(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeFunctionReplicasWithContext", func(_ context.Context, request *teov20220901.DescribeFunctionReplicasRequest) (*teov20220901.DescribeFunctionReplicasResponse, error) {
		return frV6DescribeResponse(frV6DescribeReplica()), nil
	})

	r, s := newTeoFunctionReplicaV6Resource(t)
	configureTeoFunctionReplicaV6(r, client)

	raw := teoFunctionReplicaV6Raw(t, s, map[string]tftypes.Value{
		"id":           frV6Str("zone-test123#ef-test456#test-replica"),
		"zone_id":      frV6Str("zone-test123"),
		"function_id":  frV6Str("ef-test456"),
		"replica_name": frV6Str("test-replica"),
		"content":      frV6Str("old content"),
		"remark":       frV6Str("test remark"),
	})

	req := resource.ReadRequest{State: tfsdk.State{Raw: raw, Schema: s}}
	resp := &resource.ReadResponse{State: tfsdk.State{Raw: raw.Copy(), Schema: s}}

	r.Read(context.Background(), req, resp)

	assert.False(t, resp.Diagnostics.HasError(), "unexpected diagnostics: %v", resp.Diagnostics)

	var model teo.TeoFunctionReplicaV6ResourceModel
	diags := resp.State.Get(context.Background(), &model)
	assert.False(t, diags.HasError(), "state get returned diagnostics: %v", diags)

	assert.Equal(t, "zone-test123#ef-test456#test-replica", model.Id.ValueString())
	assert.Equal(t, "zone-test123", model.ZoneId.ValueString())
	assert.Equal(t, "ef-test456", model.FunctionId.ValueString())
	assert.Equal(t, "test-replica", model.ReplicaName.ValueString())
	assert.Equal(t, "console.log(123)", model.Content.ValueString())
	assert.Equal(t, "test remark", model.Remark.ValueString())
}

// TestTeoFunctionReplicaV6Resource_ReadEmptyList verifies a successful query
// with no matching replica removes the resource from state.
func TestTeoFunctionReplicaV6Resource_ReadEmptyList(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeFunctionReplicasWithContext", func(_ context.Context, request *teov20220901.DescribeFunctionReplicasRequest) (*teov20220901.DescribeFunctionReplicasResponse, error) {
		return frV6DescribeResponse(), nil
	})

	r, s := newTeoFunctionReplicaV6Resource(t)
	configureTeoFunctionReplicaV6(r, client)

	raw := teoFunctionReplicaV6Raw(t, s, map[string]tftypes.Value{
		"id":           frV6Str("zone-test123#ef-test456#test-replica"),
		"zone_id":      frV6Str("zone-test123"),
		"function_id":  frV6Str("ef-test456"),
		"replica_name": frV6Str("test-replica"),
		"content":      frV6Str("console.log(123)"),
	})

	req := resource.ReadRequest{State: tfsdk.State{Raw: raw, Schema: s}}
	resp := &resource.ReadResponse{State: tfsdk.State{Raw: raw.Copy(), Schema: s}}

	r.Read(context.Background(), req, resp)

	assert.False(t, resp.Diagnostics.HasError(), "unexpected diagnostics: %v", resp.Diagnostics)
	assert.True(t, resp.State.Raw.IsNull(), "state should be removed, got %v", resp.State.Raw)
}

// TestTeoFunctionReplicaV6Resource_ReadAPIError verifies a nil/failed
// Describe response surfaces an error instead of removing state.
func TestTeoFunctionReplicaV6Resource_ReadAPIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeFunctionReplicasWithContext", func(_ context.Context, request *teov20220901.DescribeFunctionReplicasRequest) (*teov20220901.DescribeFunctionReplicasResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=InternalError, Message=server error")
	})

	r, s := newTeoFunctionReplicaV6Resource(t)
	configureTeoFunctionReplicaV6(r, client)

	raw := teoFunctionReplicaV6Raw(t, s, map[string]tftypes.Value{
		"id":           frV6Str("zone-test123#ef-test456#test-replica"),
		"zone_id":      frV6Str("zone-test123"),
		"function_id":  frV6Str("ef-test456"),
		"replica_name": frV6Str("test-replica"),
		"content":      frV6Str("console.log(123)"),
	})

	req := resource.ReadRequest{State: tfsdk.State{Raw: raw, Schema: s}}
	resp := &resource.ReadResponse{State: tfsdk.State{Raw: raw.Copy(), Schema: s}}

	r.Read(context.Background(), req, resp)

	assert.True(t, resp.Diagnostics.HasError(), "expected error diagnostics")
	assert.False(t, resp.State.Raw.IsNull(), "state must not be removed on API failure")
}

// TestTeoFunctionReplicaV6Resource_ReadBrokenId verifies a malformed id fails
// with a clear error.
func TestTeoFunctionReplicaV6Resource_ReadBrokenId(t *testing.T) {
	r, s := newTeoFunctionReplicaV6Resource(t)
	configureTeoFunctionReplicaV6(r, &connectivity.TencentCloudClient{})

	raw := teoFunctionReplicaV6Raw(t, s, map[string]tftypes.Value{
		"id":      frV6Str("broken-id"),
		"zone_id": frV6Str("zone-test123"),
	})

	req := resource.ReadRequest{State: tfsdk.State{Raw: raw, Schema: s}}
	resp := &resource.ReadResponse{State: tfsdk.State{Raw: raw.Copy(), Schema: s}}

	r.Read(context.Background(), req, resp)

	assert.True(t, resp.Diagnostics.HasError(), "expected error diagnostics")
	assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Invalid resource ID")
}

// TestTeoFunctionReplicaV6Resource_Update verifies the update flow sends
// content and remark to ModifyFunctionReplica and refreshes state.
func TestTeoFunctionReplicaV6Resource_Update(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(client, "UseTeoV20220901Client", teoClient)

	var gotModify *teov20220901.ModifyFunctionReplicaRequest
	patches.ApplyMethodFunc(teoClient, "ModifyFunctionReplicaWithContext", func(_ context.Context, request *teov20220901.ModifyFunctionReplicaRequest) (*teov20220901.ModifyFunctionReplicaResponse, error) {
		gotModify = request
		resp := teov20220901.NewModifyFunctionReplicaResponse()
		resp.Response = &teov20220901.ModifyFunctionReplicaResponseParams{
			RequestId: ptrStringFrV6("fake-request-id"),
		}
		return resp, nil
	})

	replica := frV6DescribeReplica()
	replica.Content = ptrStringFrV6("console.log(456)")
	replica.Remark = ptrStringFrV6("updated remark")
	patches.ApplyMethodFunc(teoClient, "DescribeFunctionReplicasWithContext", func(_ context.Context, request *teov20220901.DescribeFunctionReplicasRequest) (*teov20220901.DescribeFunctionReplicasResponse, error) {
		return frV6DescribeResponse(replica), nil
	})

	r, s := newTeoFunctionReplicaV6Resource(t)
	configureTeoFunctionReplicaV6(r, client)

	stateRaw := teoFunctionReplicaV6Raw(t, s, map[string]tftypes.Value{
		"id":           frV6Str("zone-test123#ef-test456#test-replica"),
		"zone_id":      frV6Str("zone-test123"),
		"function_id":  frV6Str("ef-test456"),
		"replica_name": frV6Str("test-replica"),
		"content":      frV6Str("console.log(123)"),
		"remark":       frV6Str("test remark"),
	})
	planRaw := teoFunctionReplicaV6Raw(t, s, map[string]tftypes.Value{
		"id":           frV6Str("zone-test123#ef-test456#test-replica"),
		"zone_id":      frV6Str("zone-test123"),
		"function_id":  frV6Str("ef-test456"),
		"replica_name": frV6Str("test-replica"),
		"content":      frV6Str("console.log(456)"),
		"remark":       frV6Str("updated remark"),
	})

	req := resource.UpdateRequest{
		Config: tfsdk.Config{Raw: planRaw, Schema: s},
		Plan:   tfsdk.Plan{Raw: planRaw, Schema: s},
		State:  tfsdk.State{Raw: stateRaw, Schema: s},
	}
	resp := &resource.UpdateResponse{State: tfsdk.State{Raw: planRaw.Copy(), Schema: s}}

	r.Update(context.Background(), req, resp)

	assert.False(t, resp.Diagnostics.HasError(), "unexpected diagnostics: %v", resp.Diagnostics)

	assert.NotNil(t, gotModify)
	assert.Equal(t, "zone-test123", *gotModify.ZoneId)
	assert.Equal(t, "ef-test456", *gotModify.FunctionId)
	assert.Equal(t, "test-replica", *gotModify.ReplicaName)
	assert.Equal(t, "console.log(456)", *gotModify.Content)
	assert.Equal(t, "updated remark", *gotModify.Remark)

	var model teo.TeoFunctionReplicaV6ResourceModel
	diags := resp.State.Get(context.Background(), &model)
	assert.False(t, diags.HasError(), "state get returned diagnostics: %v", diags)

	assert.Equal(t, "zone-test123#ef-test456#test-replica", model.Id.ValueString())
	assert.Equal(t, "console.log(456)", model.Content.ValueString())
	assert.Equal(t, "updated remark", model.Remark.ValueString())
}

// TestTeoFunctionReplicaV6Resource_UpdateRemarkRemoved verifies that removing
// the remark attribute from the configuration is rejected loudly: the clear
// semantics of ModifyFunctionReplica for an omitted Remark is not documented,
// so no direction is guessed.
func TestTeoFunctionReplicaV6Resource_UpdateRemarkRemoved(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(client, "UseTeoV20220901Client", teoClient)

	apiCalled := false
	patches.ApplyMethodFunc(teoClient, "ModifyFunctionReplicaWithContext", func(_ context.Context, request *teov20220901.ModifyFunctionReplicaRequest) (*teov20220901.ModifyFunctionReplicaResponse, error) {
		apiCalled = true
		return teov20220901.NewModifyFunctionReplicaResponse(), nil
	})

	r, s := newTeoFunctionReplicaV6Resource(t)
	configureTeoFunctionReplicaV6(r, client)

	stateRaw := teoFunctionReplicaV6Raw(t, s, map[string]tftypes.Value{
		"id":           frV6Str("zone-test123#ef-test456#test-replica"),
		"zone_id":      frV6Str("zone-test123"),
		"function_id":  frV6Str("ef-test456"),
		"replica_name": frV6Str("test-replica"),
		"content":      frV6Str("console.log(123)"),
		"remark":       frV6Str("test remark"),
	})
	// remark absent from the new configuration.
	planRaw := teoFunctionReplicaV6Raw(t, s, map[string]tftypes.Value{
		"id":           frV6Str("zone-test123#ef-test456#test-replica"),
		"zone_id":      frV6Str("zone-test123"),
		"function_id":  frV6Str("ef-test456"),
		"replica_name": frV6Str("test-replica"),
		"content":      frV6Str("console.log(123)"),
	})

	req := resource.UpdateRequest{
		Config: tfsdk.Config{Raw: planRaw, Schema: s},
		Plan:   tfsdk.Plan{Raw: planRaw, Schema: s},
		State:  tfsdk.State{Raw: stateRaw, Schema: s},
	}
	resp := &resource.UpdateResponse{State: tfsdk.State{Raw: planRaw.Copy(), Schema: s}}

	r.Update(context.Background(), req, resp)

	assert.True(t, resp.Diagnostics.HasError(), "expected error diagnostics")
	assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Clear semantics of remark is not confirmed")
	assert.False(t, apiCalled, "ModifyFunctionReplica must not be called when the clear semantics is unconfirmed")
}

// TestTeoFunctionReplicaV6Resource_UpdateRemarkStaysNull verifies that when
// remark is absent on both sides, the update proceeds without sending Remark.
func TestTeoFunctionReplicaV6Resource_UpdateRemarkStaysNull(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(client, "UseTeoV20220901Client", teoClient)

	var gotModify *teov20220901.ModifyFunctionReplicaRequest
	patches.ApplyMethodFunc(teoClient, "ModifyFunctionReplicaWithContext", func(_ context.Context, request *teov20220901.ModifyFunctionReplicaRequest) (*teov20220901.ModifyFunctionReplicaResponse, error) {
		gotModify = request
		return teov20220901.NewModifyFunctionReplicaResponse(), nil
	})

	replica := frV6DescribeReplica()
	replica.Remark = nil
	patches.ApplyMethodFunc(teoClient, "DescribeFunctionReplicasWithContext", func(_ context.Context, request *teov20220901.DescribeFunctionReplicasRequest) (*teov20220901.DescribeFunctionReplicasResponse, error) {
		return frV6DescribeResponse(replica), nil
	})

	r, s := newTeoFunctionReplicaV6Resource(t)
	configureTeoFunctionReplicaV6(r, client)

	stateRaw := teoFunctionReplicaV6Raw(t, s, map[string]tftypes.Value{
		"id":           frV6Str("zone-test123#ef-test456#test-replica"),
		"zone_id":      frV6Str("zone-test123"),
		"function_id":  frV6Str("ef-test456"),
		"replica_name": frV6Str("test-replica"),
		"content":      frV6Str("console.log(123)"),
	})
	planRaw := teoFunctionReplicaV6Raw(t, s, map[string]tftypes.Value{
		"id":           frV6Str("zone-test123#ef-test456#test-replica"),
		"zone_id":      frV6Str("zone-test123"),
		"function_id":  frV6Str("ef-test456"),
		"replica_name": frV6Str("test-replica"),
		"content":      frV6Str("console.log(456)"),
	})

	req := resource.UpdateRequest{
		Config: tfsdk.Config{Raw: planRaw, Schema: s},
		Plan:   tfsdk.Plan{Raw: planRaw, Schema: s},
		State:  tfsdk.State{Raw: stateRaw, Schema: s},
	}
	resp := &resource.UpdateResponse{State: tfsdk.State{Raw: planRaw.Copy(), Schema: s}}

	r.Update(context.Background(), req, resp)

	assert.False(t, resp.Diagnostics.HasError(), "unexpected diagnostics: %v", resp.Diagnostics)

	assert.NotNil(t, gotModify)
	assert.Nil(t, gotModify.Remark, "remark should not be sent when null on both sides")
	assert.Equal(t, "console.log(456)", *gotModify.Content)

	var model teo.TeoFunctionReplicaV6ResourceModel
	diags := resp.State.Get(context.Background(), &model)
	assert.False(t, diags.HasError(), "state get returned diagnostics: %v", diags)
	assert.True(t, model.Remark.IsNull(), "remark in state should stay null, got %v", model.Remark)
}

// TestTeoFunctionReplicaV6Resource_UpdateAPIError verifies the update error
// branch surfaces an error diagnostic.
func TestTeoFunctionReplicaV6Resource_UpdateAPIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "ModifyFunctionReplicaWithContext", func(_ context.Context, request *teov20220901.ModifyFunctionReplicaRequest) (*teov20220901.ModifyFunctionReplicaResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=InvalidParameter, Message=content too large")
	})

	r, s := newTeoFunctionReplicaV6Resource(t)
	configureTeoFunctionReplicaV6(r, client)

	stateRaw := teoFunctionReplicaV6Raw(t, s, map[string]tftypes.Value{
		"id":           frV6Str("zone-test123#ef-test456#test-replica"),
		"zone_id":      frV6Str("zone-test123"),
		"function_id":  frV6Str("ef-test456"),
		"replica_name": frV6Str("test-replica"),
		"content":      frV6Str("console.log(123)"),
	})
	planRaw := teoFunctionReplicaV6Raw(t, s, map[string]tftypes.Value{
		"id":           frV6Str("zone-test123#ef-test456#test-replica"),
		"zone_id":      frV6Str("zone-test123"),
		"function_id":  frV6Str("ef-test456"),
		"replica_name": frV6Str("test-replica"),
		"content":      frV6Str("console.log(456)"),
	})

	req := resource.UpdateRequest{
		Config: tfsdk.Config{Raw: planRaw, Schema: s},
		Plan:   tfsdk.Plan{Raw: planRaw, Schema: s},
		State:  tfsdk.State{Raw: stateRaw, Schema: s},
	}
	resp := &resource.UpdateResponse{State: tfsdk.State{Raw: planRaw.Copy(), Schema: s}}

	r.Update(context.Background(), req, resp)

	assert.True(t, resp.Diagnostics.HasError(), "expected error diagnostics")
	assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Error updating TEO function replica")
}

// TestTeoFunctionReplicaV6Resource_Delete verifies the delete flow sends the
// single replica name.
func TestTeoFunctionReplicaV6Resource_Delete(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeFunctionReplicasWithContext", func(_ context.Context, request *teov20220901.DescribeFunctionReplicasRequest) (*teov20220901.DescribeFunctionReplicasResponse, error) {
		return frV6DescribeResponse(frV6DescribeReplica()), nil
	})

	var gotDelete *teov20220901.DeleteFunctionReplicaRequest
	patches.ApplyMethodFunc(teoClient, "DeleteFunctionReplicaWithContext", func(_ context.Context, request *teov20220901.DeleteFunctionReplicaRequest) (*teov20220901.DeleteFunctionReplicaResponse, error) {
		gotDelete = request
		resp := teov20220901.NewDeleteFunctionReplicaResponse()
		resp.Response = &teov20220901.DeleteFunctionReplicaResponseParams{
			RequestId: ptrStringFrV6("fake-request-id"),
		}
		return resp, nil
	})

	r, s := newTeoFunctionReplicaV6Resource(t)
	configureTeoFunctionReplicaV6(r, client)

	raw := teoFunctionReplicaV6Raw(t, s, map[string]tftypes.Value{
		"id":           frV6Str("zone-test123#ef-test456#test-replica"),
		"zone_id":      frV6Str("zone-test123"),
		"function_id":  frV6Str("ef-test456"),
		"replica_name": frV6Str("test-replica"),
		"content":      frV6Str("console.log(123)"),
	})

	req := resource.DeleteRequest{State: tfsdk.State{Raw: raw, Schema: s}}
	resp := &resource.DeleteResponse{State: tfsdk.State{Raw: raw.Copy(), Schema: s}}

	r.Delete(context.Background(), req, resp)

	assert.False(t, resp.Diagnostics.HasError(), "unexpected diagnostics: %v", resp.Diagnostics)

	assert.NotNil(t, gotDelete)
	assert.Equal(t, "zone-test123", *gotDelete.ZoneId)
	assert.Equal(t, "ef-test456", *gotDelete.FunctionId)
	assert.Equal(t, []*string{ptrStringFrV6("test-replica")}, gotDelete.ReplicaNames)
}

// TestTeoFunctionReplicaV6Resource_DeleteNotFound verifies destroy
// idempotency: when the replica is confirmed gone, Delete succeeds without
// calling the Delete API.
func TestTeoFunctionReplicaV6Resource_DeleteNotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeFunctionReplicasWithContext", func(_ context.Context, request *teov20220901.DescribeFunctionReplicasRequest) (*teov20220901.DescribeFunctionReplicasResponse, error) {
		return frV6DescribeResponse(), nil
	})

	deleteCalled := false
	patches.ApplyMethodFunc(teoClient, "DeleteFunctionReplicaWithContext", func(_ context.Context, request *teov20220901.DeleteFunctionReplicaRequest) (*teov20220901.DeleteFunctionReplicaResponse, error) {
		deleteCalled = true
		return teov20220901.NewDeleteFunctionReplicaResponse(), nil
	})

	r, s := newTeoFunctionReplicaV6Resource(t)
	configureTeoFunctionReplicaV6(r, client)

	raw := teoFunctionReplicaV6Raw(t, s, map[string]tftypes.Value{
		"id":           frV6Str("zone-test123#ef-test456#test-replica"),
		"zone_id":      frV6Str("zone-test123"),
		"function_id":  frV6Str("ef-test456"),
		"replica_name": frV6Str("test-replica"),
		"content":      frV6Str("console.log(123)"),
	})

	req := resource.DeleteRequest{State: tfsdk.State{Raw: raw, Schema: s}}
	resp := &resource.DeleteResponse{State: tfsdk.State{Raw: raw.Copy(), Schema: s}}

	r.Delete(context.Background(), req, resp)

	assert.False(t, resp.Diagnostics.HasError(), "unexpected diagnostics: %v", resp.Diagnostics)
	assert.False(t, deleteCalled, "DeleteFunctionReplica must not be called when the replica is already gone")
}

// TestTeoFunctionReplicaV6Resource_DeleteAPIError verifies the delete error
// branch surfaces an error diagnostic.
func TestTeoFunctionReplicaV6Resource_DeleteAPIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeFunctionReplicasWithContext", func(_ context.Context, request *teov20220901.DescribeFunctionReplicasRequest) (*teov20220901.DescribeFunctionReplicasResponse, error) {
		return frV6DescribeResponse(frV6DescribeReplica()), nil
	})

	patches.ApplyMethodFunc(teoClient, "DeleteFunctionReplicaWithContext", func(_ context.Context, request *teov20220901.DeleteFunctionReplicaRequest) (*teov20220901.DeleteFunctionReplicaResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=InternalError, Message=server error")
	})

	r, s := newTeoFunctionReplicaV6Resource(t)
	configureTeoFunctionReplicaV6(r, client)

	raw := teoFunctionReplicaV6Raw(t, s, map[string]tftypes.Value{
		"id":           frV6Str("zone-test123#ef-test456#test-replica"),
		"zone_id":      frV6Str("zone-test123"),
		"function_id":  frV6Str("ef-test456"),
		"replica_name": frV6Str("test-replica"),
		"content":      frV6Str("console.log(123)"),
	})

	req := resource.DeleteRequest{State: tfsdk.State{Raw: raw, Schema: s}}
	resp := &resource.DeleteResponse{State: tfsdk.State{Raw: raw.Copy(), Schema: s}}

	r.Delete(context.Background(), req, resp)

	assert.True(t, resp.Diagnostics.HasError(), "expected error diagnostics")
	assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Error deleting TEO function replica")
}

// TestTeoFunctionReplicaV6Resource_ImportState verifies the composite import
// id round-trips into state attributes.
func TestTeoFunctionReplicaV6Resource_ImportState(t *testing.T) {
	r, s := newTeoFunctionReplicaV6Resource(t)

	req := resource.ImportStateRequest{ID: "zone-test123#ef-test456#test-replica"}
	resp := &resource.ImportStateResponse{
		State: tfsdk.State{Raw: teoFunctionReplicaV6Null(s), Schema: s},
	}

	r.ImportState(context.Background(), req, resp)

	assert.False(t, resp.Diagnostics.HasError(), "unexpected diagnostics: %v", resp.Diagnostics)

	var model teo.TeoFunctionReplicaV6ResourceModel
	diags := resp.State.Get(context.Background(), &model)
	assert.False(t, diags.HasError(), "state get returned diagnostics: %v", diags)

	assert.Equal(t, "zone-test123#ef-test456#test-replica", model.Id.ValueString())
	assert.Equal(t, "zone-test123", model.ZoneId.ValueString())
	assert.Equal(t, "ef-test456", model.FunctionId.ValueString())
	assert.Equal(t, "test-replica", model.ReplicaName.ValueString())
}

// TestTeoFunctionReplicaV6Resource_ImportStateBrokenID covers the import
// failure branches of the composite id parser.
func TestTeoFunctionReplicaV6Resource_ImportStateBrokenID(t *testing.T) {
	r, s := newTeoFunctionReplicaV6Resource(t)

	for _, id := range []string{
		"",
		"only-zone",
		"zone#function",
		"zone#function#",
		"#function#replica",
		"a#b#c#d",
	} {
		req := resource.ImportStateRequest{ID: id}
		resp := &resource.ImportStateResponse{
			State: tfsdk.State{Raw: teoFunctionReplicaV6Null(s), Schema: s},
		}
		r.ImportState(context.Background(), req, resp)
		assert.True(t, resp.Diagnostics.HasError(), "expected error diagnostics for import id %q", id)
		assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Invalid import ID")
	}
}
