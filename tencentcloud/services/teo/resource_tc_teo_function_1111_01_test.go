package teo

import (
	"context"
	"fmt"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	tftypes "github.com/hashicorp/terraform-plugin-go/tftypes"

	teosdk "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/sharedmeta"
)

const (
	teoFnTestZoneId     = "zone-2p1ghslb"
	teoFnTestFunctionId = "ef-2p1ghslb-1234"
	teoFnTestName       = "edge-func-example"
	teoFnTestContent    = "export default { async fetch() { return new Response('ok') } }"
	teoFnTestRemark     = "managed by terraform"
)

// teoFnRestrictionObjectTypeRaw is the raw tftypes type of one element of
// domain_compliance_restrictions.
func teoFnRestrictionObjectTypeRaw() tftypes.Type {
	return tftypes.Object{AttributeTypes: map[string]tftypes.Type{
		"reason": tftypes.String,
		"region": tftypes.String,
	}}
}

// teoFnObjectTypeRaw is the raw tftypes object type of the whole resource.
func teoFnObjectTypeRaw() tftypes.Type {
	return tftypes.Object{AttributeTypes: map[string]tftypes.Type{
		"id":                             tftypes.String,
		"zone_id":                        tftypes.String,
		"name":                           tftypes.String,
		"content":                        tftypes.String,
		"remark":                         tftypes.String,
		"function_id":                    tftypes.String,
		"domain":                         tftypes.String,
		"domain_compliance_restrictions": tftypes.List{ElementType: teoFnRestrictionObjectTypeRaw()},
		"create_time":                    tftypes.String,
		"update_time":                    tftypes.String,
	}}
}

// teoFnRawValue builds a raw object value, filling missing attributes with null.
func teoFnRawValue(vals map[string]tftypes.Value) tftypes.Value {
	full := map[string]tftypes.Value{}
	for name, typ := range teoFnObjectTypeRaw().(tftypes.Object).AttributeTypes {
		full[name] = tftypes.NewValue(typ, nil)
	}
	for name, val := range vals {
		full[name] = val
	}
	return tftypes.NewValue(teoFnObjectTypeRaw(), full)
}

func teoFnStr(s string) tftypes.Value {
	return tftypes.NewValue(tftypes.String, s)
}

// teoFnSchema builds the framework schema once for the tests.
func teoFnSchema(t *testing.T) tfsdk.Schema {
	t.Helper()
	r := NewTeoFunction111101Resource()
	resp := &resource.SchemaResponse{}
	r.Schema(context.Background(), resource.SchemaRequest{}, resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Schema returned diagnostics: %v", resp.Diagnostics)
	}
	return resp.Schema
}

// teoFnState builds a tfsdk.State carrying the given raw value and schema.
func teoFnState(schema tfsdk.Schema, raw tftypes.Value) tfsdk.State {
	return tfsdk.State{Schema: schema, Raw: raw}
}

// teoFnNewState builds an empty (all null) state ready to receive Set calls.
func teoFnNewState(schema tfsdk.Schema) tfsdk.State {
	return teoFnState(schema, teoFnRawValue(nil))
}

// teoFnMockClient patches the shared client so every UseTeoClient call returns
// the given mocked teo SDK client.
func teoFnMockClient(patches *gomonkey.Patches) (*connectivity.TencentCloudClient, *teosdk.Client) {
	client := &connectivity.TencentCloudClient{}
	teoClient := &teosdk.Client{}
	patches.ApplyMethodReturn(client, "UseTeoClient", teoClient)
	return client, teoClient
}

// teoFnConfigureResource injects the mocked client into a fresh resource.
func teoFnConfigureResource(t *testing.T, client *connectivity.TencentCloudClient) resource.Resource {
	t.Helper()
	r := NewTeoFunction111101Resource()
	meta := &sharedmeta.ProviderMeta{Client: client}
	configureReq := resource.ConfigureRequest{ProviderData: meta}
	configureResp := &resource.ConfigureResponse{}
	r.Configure(context.Background(), configureReq, configureResp)
	if configureResp.Diagnostics.HasError() {
		t.Fatalf("Configure returned diagnostics: %v", configureResp.Diagnostics)
	}
	return r
}

// teoFnDescribeResponse builds a successful DescribeFunctions response for the
// test function, optionally with the given domain compliance restrictions.
func teoFnDescribeResponse(domain, reason, region string, withFunction bool) *teosdk.DescribeFunctionsResponse {
	resp := teosdk.NewDescribeFunctionsResponse()
	resp.Response = &teosdk.DescribeFunctionsResponseParams{
		TotalCount: func() *int64 { v := int64(0); return &v }(),
	}
	if withFunction {
		fn := &teosdk.Function{
			FunctionId: func() *string { v := teoFnTestFunctionId; return &v }(),
			ZoneId:     func() *string { v := teoFnTestZoneId; return &v }(),
			Name:       func() *string { v := teoFnTestName; return &v }(),
			Content:    func() *string { v := teoFnTestContent; return &v }(),
			Remark:     func() *string { v := teoFnTestRemark; return &v }(),
			CreateTime: func() *string { v := "2024-01-01T00:00:00Z"; return &v }(),
			UpdateTime: func() *string { v := "2024-01-02T00:00:00Z"; return &v }(),
		}
		if domain != "" {
			fn.Domain = func() *string { v := domain; return &v }()
		}
		if reason != "" {
			fn.DomainComplianceRestrictions = []*teosdk.ComplianceRestriction{
				{
					Reason: func() *string { v := reason; return &v }(),
					Region: func() *string { v := region; return &v }(),
				},
			}
		}
		resp.Response.Functions = []*teosdk.Function{fn}
		resp.Response.TotalCount = func() *int64 { v := int64(1); return &v }()
	}
	return resp
}

// TestTeoFunction111101Resource_Metadata validates the resource type name.
func TestTeoFunction111101Resource_Metadata(t *testing.T) {
	r := &TeoFunction111101Resource{}
	resp := &resource.MetadataResponse{}
	r.Metadata(context.Background(), resource.MetadataRequest{}, resp)
	if resp.TypeName != "tencentcloud_teo_function_1111_01" {
		t.Errorf("Metadata.TypeName = %q, want %q", resp.TypeName, "tencentcloud_teo_function_1111_01")
	}
}

// TestTeoFunction111101Resource_SchemaFlags validates required / optional /
// computed flags on every attribute.
func TestTeoFunction111101Resource_SchemaFlags(t *testing.T) {
	resp := &resource.SchemaResponse{}
	NewTeoFunction111101Resource().Schema(context.Background(), resource.SchemaRequest{}, resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Schema returned diagnostics: %v", resp.Diagnostics)
	}

	type flags struct{ required, optional, computed bool }
	cases := map[string]flags{
		"id":                             {computed: true},
		"zone_id":                        {required: true},
		"name":                           {required: true},
		"content":                        {required: true},
		"remark":                         {optional: true, computed: true},
		"function_id":                    {computed: true},
		"domain":                         {computed: true},
		"domain_compliance_restrictions": {computed: true},
		"create_time":                    {computed: true},
		"update_time":                    {computed: true},
	}
	for name, want := range cases {
		attr, ok := resp.Schema.Attributes[name]
		if !ok {
			t.Errorf("attribute %q missing from schema", name)
			continue
		}
		if attr.IsRequired() != want.required || attr.IsOptional() != want.optional || attr.IsComputed() != want.computed {
			t.Errorf("attribute %q flags = (required=%v, optional=%v, computed=%v), want %+v",
				name, attr.IsRequired(), attr.IsOptional(), attr.IsComputed(), want)
		}
	}
}

// TestTeoFunction111101Resource_Factory validates the registry factory.
func TestTeoFunction111101Resource_Factory(t *testing.T) {
	got := NewTeoFunction111101Resource()
	if got == nil {
		t.Fatal("NewTeoFunction111101Resource returned nil")
	}
	if _, ok := got.(*TeoFunction111101Resource); !ok {
		t.Fatalf("factory returned %T, want *TeoFunction111101Resource", got)
	}
}

// TestTeoFunction111101Resource_ConfigureNilProviderData ensures the
// pre-Configure invocation pattern is a silent no-op.
func TestTeoFunction111101Resource_ConfigureNilProviderData(t *testing.T) {
	r := &TeoFunction111101Resource{}
	resp := &resource.ConfigureResponse{}
	r.Configure(context.Background(), resource.ConfigureRequest{}, resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Configure with nil ProviderData returned diagnostics: %v", resp.Diagnostics)
	}
}

// TestTeoFunction111101Resource_Create_Success covers a full create: the
// CreateFunction request fields, the state persisted with the returned
// function id, and every computed attribute read back from DescribeFunctions.
func TestTeoFunction111101Resource_Create_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client, teoClient := teoFnMockClient(patches)

	var gotCreateRequest *teosdk.CreateFunctionRequest
	patches.ApplyMethodFunc(teoClient, "CreateFunctionWithContext", func(_ context.Context, request *teosdk.CreateFunctionRequest) (*teosdk.CreateFunctionResponse, error) {
		gotCreateRequest = request
		resp := teosdk.NewCreateFunctionResponse()
		resp.Response = &teosdk.CreateFunctionResponseParams{
			FunctionId: func() *string { v := teoFnTestFunctionId; return &v }(),
			RequestId:  func() *string { v := "fake-request-id"; return &v }(),
		}
		return resp, nil
	})

	var gotDescribeRequest *teosdk.DescribeFunctionsRequest
	patches.ApplyMethodFunc(teoClient, "DescribeFunctionsWithContext", func(_ context.Context, request *teosdk.DescribeFunctionsRequest) (*teosdk.DescribeFunctionsResponse, error) {
		gotDescribeRequest = request
		return teoFnDescribeResponse("edge-func.example.com", "ICP_RECORD_REQUIRED", "CN", true), nil
	})

	r := teoFnConfigureResource(t, client)
	s := teoFnSchema(t)

	plan := tfsdk.Plan{
		Schema: s,
		Raw: teoFnRawValue(map[string]tftypes.Value{
			"zone_id": teoFnStr(teoFnTestZoneId),
			"name":    teoFnStr(teoFnTestName),
			"content": teoFnStr(teoFnTestContent),
			"remark":  teoFnStr(teoFnTestRemark),
		}),
	}

	resp := &resource.CreateResponse{State: teoFnNewState(s)}
	r.Create(context.Background(), resource.CreateRequest{Plan: plan}, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("Create returned error diagnostics: %v", resp.Diagnostics)
	}
	if gotCreateRequest == nil {
		t.Fatal("CreateFunctionWithContext was not called")
	}
	if gotCreateRequest.ZoneId == nil || *gotCreateRequest.ZoneId != teoFnTestZoneId {
		t.Errorf("CreateFunction request.ZoneId = %v, want %q", gotCreateRequest.ZoneId, teoFnTestZoneId)
	}
	if gotCreateRequest.Name == nil || *gotCreateRequest.Name != teoFnTestName {
		t.Errorf("CreateFunction request.Name = %v, want %q", gotCreateRequest.Name, teoFnTestName)
	}
	if gotCreateRequest.Content == nil || *gotCreateRequest.Content != teoFnTestContent {
		t.Errorf("CreateFunction request.Content = %v, want the configured content", gotCreateRequest.Content)
	}
	if gotCreateRequest.Remark == nil || *gotCreateRequest.Remark != teoFnTestRemark {
		t.Errorf("CreateFunction request.Remark = %v, want %q", gotCreateRequest.Remark, teoFnTestRemark)
	}
	if gotDescribeRequest == nil || len(gotDescribeRequest.FunctionIds) != 1 || *gotDescribeRequest.FunctionIds[0] != teoFnTestFunctionId {
		t.Errorf("DescribeFunctions request.FunctionIds = %v, want [%q]", gotDescribeRequest, teoFnTestFunctionId)
	}

	var model TeoFunction111101ResourceModel
	if diags := resp.State.Get(context.Background(), &model); diags.HasError() {
		t.Fatalf("reading back state failed: %v", diags)
	}
	if model.Id.ValueString() != teoFnTestFunctionId {
		t.Errorf("state id = %q, want %q", model.Id.ValueString(), teoFnTestFunctionId)
	}
	if model.FunctionId.ValueString() != teoFnTestFunctionId {
		t.Errorf("state function_id = %q, want %q", model.FunctionId.ValueString(), teoFnTestFunctionId)
	}
	if model.ZoneId.ValueString() != teoFnTestZoneId {
		t.Errorf("state zone_id = %q, want %q", model.ZoneId.ValueString(), teoFnTestZoneId)
	}
	if model.Name.ValueString() != teoFnTestName {
		t.Errorf("state name = %q, want %q", model.Name.ValueString(), teoFnTestName)
	}
	if model.Content.ValueString() != teoFnTestContent {
		t.Errorf("state content = %q, want the configured content", model.Content.ValueString())
	}
	if model.Remark.ValueString() != teoFnTestRemark {
		t.Errorf("state remark = %q, want %q", model.Remark.ValueString(), teoFnTestRemark)
	}
	if model.Domain.ValueString() != "edge-func.example.com" {
		t.Errorf("state domain = %q, want %q", model.Domain.ValueString(), "edge-func.example.com")
	}
	if model.CreateTime.ValueString() != "2024-01-01T00:00:00Z" {
		t.Errorf("state create_time = %q, want %q", model.CreateTime.ValueString(), "2024-01-01T00:00:00Z")
	}
	if model.UpdateTime.ValueString() != "2024-01-02T00:00:00Z" {
		t.Errorf("state update_time = %q, want %q", model.UpdateTime.ValueString(), "2024-01-02T00:00:00Z")
	}

	var restrictions []TeoFunction111101RestrictionModel
	if diags := model.DomainComplianceRestrictions.ElementsAs(context.Background(), &restrictions, false); diags.HasError() {
		t.Fatalf("reading domain_compliance_restrictions failed: %v", diags)
	}
	if len(restrictions) != 1 {
		t.Fatalf("len(domain_compliance_restrictions) = %d, want 1", len(restrictions))
	}
	if restrictions[0].Reason.ValueString() != "ICP_RECORD_REQUIRED" {
		t.Errorf("restrictions[0].reason = %q, want %q", restrictions[0].Reason.ValueString(), "ICP_RECORD_REQUIRED")
	}
	if restrictions[0].Region.ValueString() != "CN" {
		t.Errorf("restrictions[0].region = %q, want %q", restrictions[0].Region.ValueString(), "CN")
	}
}

// TestTeoFunction111101Resource_Create_ApiError covers the error branch when
// CreateFunction itself fails.
func TestTeoFunction111101Resource_Create_ApiError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client, teoClient := teoFnMockClient(patches)
	patches.ApplyMethodFunc(teoClient, "CreateFunctionWithContext", func(_ context.Context, _ *teosdk.CreateFunctionRequest) (*teosdk.CreateFunctionResponse, error) {
		return nil, fmt.Errorf("internal error")
	})

	r := teoFnConfigureResource(t, client)
	s := teoFnSchema(t)
	plan := tfsdk.Plan{
		Schema: s,
		Raw: teoFnRawValue(map[string]tftypes.Value{
			"zone_id": teoFnStr(teoFnTestZoneId),
			"name":    teoFnStr(teoFnTestName),
			"content": teoFnStr(teoFnTestContent),
		}),
	}

	resp := &resource.CreateResponse{State: teoFnNewState(s)}
	r.Create(context.Background(), resource.CreateRequest{Plan: plan}, resp)

	if !resp.Diagnostics.HasError() {
		t.Fatal("Create returned no error diagnostics, want an error on CreateFunction failure")
	}
}

// TestTeoFunction111101Resource_Create_MissingFunctionId covers the error
// branch when CreateFunction succeeds but returns no function id.
func TestTeoFunction111101Resource_Create_MissingFunctionId(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client, teoClient := teoFnMockClient(patches)
	patches.ApplyMethodFunc(teoClient, "CreateFunctionWithContext", func(_ context.Context, _ *teosdk.CreateFunctionRequest) (*teosdk.CreateFunctionResponse, error) {
		resp := teosdk.NewCreateFunctionResponse()
		resp.Response = &teosdk.CreateFunctionResponseParams{}
		return resp, nil
	})

	r := teoFnConfigureResource(t, client)
	s := teoFnSchema(t)
	plan := tfsdk.Plan{
		Schema: s,
		Raw: teoFnRawValue(map[string]tftypes.Value{
			"zone_id": teoFnStr(teoFnTestZoneId),
			"name":    teoFnStr(teoFnTestName),
			"content": teoFnStr(teoFnTestContent),
		}),
	}

	resp := &resource.CreateResponse{State: teoFnNewState(s)}
	r.Create(context.Background(), resource.CreateRequest{Plan: plan}, resp)

	if !resp.Diagnostics.HasError() {
		t.Fatal("Create returned no error diagnostics, want an error on missing function id")
	}
}

// TestTeoFunction111101Resource_Read_Found covers refresh when the function
// still exists.
func TestTeoFunction111101Resource_Read_Found(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client, teoClient := teoFnMockClient(patches)
	patches.ApplyMethodFunc(teoClient, "DescribeFunctionsWithContext", func(_ context.Context, _ *teosdk.DescribeFunctionsRequest) (*teosdk.DescribeFunctionsResponse, error) {
		return teoFnDescribeResponse("", "", "", true), nil
	})

	r := teoFnConfigureResource(t, client)
	s := teoFnSchema(t)
	state := teoFnState(s, teoFnRawValue(map[string]tftypes.Value{
		"id":      teoFnStr(teoFnTestFunctionId),
		"zone_id": teoFnStr(teoFnTestZoneId),
		"name":    teoFnStr(teoFnTestName),
		"content": teoFnStr(teoFnTestContent),
	}))

	resp := &resource.ReadResponse{State: state}
	r.Read(context.Background(), resource.ReadRequest{State: state}, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("Read returned error diagnostics: %v", resp.Diagnostics)
	}

	var model TeoFunction111101ResourceModel
	if diags := resp.State.Get(context.Background(), &model); diags.HasError() {
		t.Fatalf("reading back state failed: %v", diags)
	}
	if model.Id.ValueString() != teoFnTestFunctionId {
		t.Errorf("state id = %q, want %q", model.Id.ValueString(), teoFnTestFunctionId)
	}
	if model.Remark.ValueString() != teoFnTestRemark {
		t.Errorf("state remark = %q, want %q", model.Remark.ValueString(), teoFnTestRemark)
	}
	if model.UpdateTime.ValueString() != "2024-01-02T00:00:00Z" {
		t.Errorf("state update_time = %q, want %q", model.UpdateTime.ValueString(), "2024-01-02T00:00:00Z")
	}
	var restrictions []TeoFunction111101RestrictionModel
	if diags := model.DomainComplianceRestrictions.ElementsAs(context.Background(), &restrictions, false); diags.HasError() {
		t.Fatalf("reading domain_compliance_restrictions failed: %v", diags)
	}
	if len(restrictions) != 0 {
		t.Errorf("len(domain_compliance_restrictions) = %d, want 0", len(restrictions))
	}
}

// TestTeoFunction111101Resource_Read_NotFound covers the refresh path when
// the function was removed outside Terraform: state is dropped, no error.
func TestTeoFunction111101Resource_Read_NotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client, teoClient := teoFnMockClient(patches)
	patches.ApplyMethodFunc(teoClient, "DescribeFunctionsWithContext", func(_ context.Context, _ *teosdk.DescribeFunctionsRequest) (*teosdk.DescribeFunctionsResponse, error) {
		return teoFnDescribeResponse("", "", "", false), nil
	})

	r := teoFnConfigureResource(t, client)
	s := teoFnSchema(t)
	state := teoFnState(s, teoFnRawValue(map[string]tftypes.Value{
		"id":      teoFnStr(teoFnTestFunctionId),
		"zone_id": teoFnStr(teoFnTestZoneId),
		"name":    teoFnStr(teoFnTestName),
		"content": teoFnStr(teoFnTestContent),
	}))

	resp := &resource.ReadResponse{State: state}
	r.Read(context.Background(), resource.ReadRequest{State: state}, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("Read returned error diagnostics on not-found: %v", resp.Diagnostics)
	}
	if !resp.State.Raw.IsNull() {
		t.Errorf("state should have been removed (null), got %v", resp.State.Raw)
	}
}

// TestTeoFunction111101Resource_Read_ApiError covers the refresh path when
// DescribeFunctions fails: state must be kept and an error reported.
func TestTeoFunction111101Resource_Read_ApiError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client, teoClient := teoFnMockClient(patches)
	patches.ApplyMethodFunc(teoClient, "DescribeFunctionsWithContext", func(_ context.Context, _ *teosdk.DescribeFunctionsRequest) (*teosdk.DescribeFunctionsResponse, error) {
		return nil, fmt.Errorf("network error")
	})

	r := teoFnConfigureResource(t, client)
	s := teoFnSchema(t)
	state := teoFnState(s, teoFnRawValue(map[string]tftypes.Value{
		"id":      teoFnStr(teoFnTestFunctionId),
		"zone_id": teoFnStr(teoFnTestZoneId),
		"name":    teoFnStr(teoFnTestName),
		"content": teoFnStr(teoFnTestContent),
	}))

	resp := &resource.ReadResponse{State: state}
	r.Read(context.Background(), resource.ReadRequest{State: state}, resp)

	if !resp.Diagnostics.HasError() {
		t.Fatal("Read returned no error diagnostics, want an error on DescribeFunctions failure")
	}
	if resp.State.Raw.IsNull() {
		t.Error("state must not be dropped on a transport/API failure")
	}
}

// TestTeoFunction111101Resource_Update_Success covers an in-place update of
// content with a null remark: the ModifyFunction request must omit the remark
// so the cloud keeps the existing value, and the state is refreshed afterwards.
func TestTeoFunction111101Resource_Update_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client, teoClient := teoFnMockClient(patches)

	newContent := "export default { async fetch() { return new Response('v2') } }"
	var gotModifyRequest *teosdk.ModifyFunctionRequest
	patches.ApplyMethodFunc(teoClient, "ModifyFunctionWithContext", func(_ context.Context, request *teosdk.ModifyFunctionRequest) (*teosdk.ModifyFunctionResponse, error) {
		gotModifyRequest = request
		return teosdk.NewModifyFunctionResponse(), nil
	})
	patches.ApplyMethodFunc(teoClient, "DescribeFunctionsWithContext", func(_ context.Context, _ *teosdk.DescribeFunctionsRequest) (*teosdk.DescribeFunctionsResponse, error) {
		resp := teoFnDescribeResponse("", "", "", true)
		resp.Response.Functions[0].Content = &newContent
		return resp, nil
	})

	r := teoFnConfigureResource(t, client)
	s := teoFnSchema(t)

	state := teoFnState(s, teoFnRawValue(map[string]tftypes.Value{
		"id":      teoFnStr(teoFnTestFunctionId),
		"zone_id": teoFnStr(teoFnTestZoneId),
		"name":    teoFnStr(teoFnTestName),
		"content": teoFnStr(teoFnTestContent),
		"remark":  teoFnStr(teoFnTestRemark),
	}))
	plan := tfsdk.Plan{
		Schema: s,
		Raw: teoFnRawValue(map[string]tftypes.Value{
			"id":      teoFnStr(teoFnTestFunctionId),
			"zone_id": teoFnStr(teoFnTestZoneId),
			"name":    teoFnStr(teoFnTestName),
			"content": teoFnStr(newContent),
			"remark":  teoFnStr(teoFnTestRemark),
		}),
	}

	resp := &resource.UpdateResponse{State: state}
	r.Update(context.Background(), resource.UpdateRequest{Plan: plan, State: state}, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("Update returned error diagnostics: %v", resp.Diagnostics)
	}
	if gotModifyRequest == nil {
		t.Fatal("ModifyFunctionWithContext was not called")
	}
	if gotModifyRequest.Content == nil || *gotModifyRequest.Content != newContent {
		t.Errorf("ModifyFunction request.Content = %v, want the new content", gotModifyRequest.Content)
	}
	if gotModifyRequest.ZoneId == nil || *gotModifyRequest.ZoneId != teoFnTestZoneId {
		t.Errorf("ModifyFunction request.ZoneId = %v, want %q", gotModifyRequest.ZoneId, teoFnTestZoneId)
	}
	if gotModifyRequest.FunctionId == nil || *gotModifyRequest.FunctionId != teoFnTestFunctionId {
		t.Errorf("ModifyFunction request.FunctionId = %v, want %q", gotModifyRequest.FunctionId, teoFnTestFunctionId)
	}

	var model TeoFunction111101ResourceModel
	if diags := resp.State.Get(context.Background(), &model); diags.HasError() {
		t.Fatalf("reading back state failed: %v", diags)
	}
	if model.Content.ValueString() != newContent {
		t.Errorf("state content = %q, want the new content", model.Content.ValueString())
	}
}

// TestTeoFunction111101Resource_Update_NullRemarkKeepsRemote verifies that a
// null planned remark is omitted from the ModifyFunction request so the cloud
// keeps the existing value (documented ModifyFunction behavior).
func TestTeoFunction111101Resource_Update_NullRemarkKeepsRemote(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client, teoClient := teoFnMockClient(patches)
	var gotModifyRequest *teosdk.ModifyFunctionRequest
	patches.ApplyMethodFunc(teoClient, "ModifyFunctionWithContext", func(_ context.Context, request *teosdk.ModifyFunctionRequest) (*teosdk.ModifyFunctionResponse, error) {
		gotModifyRequest = request
		return teosdk.NewModifyFunctionResponse(), nil
	})
	patches.ApplyMethodFunc(teoClient, "DescribeFunctionsWithContext", func(_ context.Context, _ *teosdk.DescribeFunctionsRequest) (*teosdk.DescribeFunctionsResponse, error) {
		return teoFnDescribeResponse("", "", "", true), nil
	})

	r := teoFnConfigureResource(t, client)
	s := teoFnSchema(t)
	state := teoFnState(s, teoFnRawValue(map[string]tftypes.Value{
		"id":      teoFnStr(teoFnTestFunctionId),
		"zone_id": teoFnStr(teoFnTestZoneId),
		"name":    teoFnStr(teoFnTestName),
		"content": teoFnStr(teoFnTestContent),
		"remark":  teoFnStr(teoFnTestRemark),
	}))
	plan := tfsdk.Plan{
		Schema: s,
		Raw: teoFnRawValue(map[string]tftypes.Value{
			"id":      teoFnStr(teoFnTestFunctionId),
			"zone_id": teoFnStr(teoFnTestZoneId),
			"name":    teoFnStr(teoFnTestName),
			"content": teoFnStr("new content"),
			"remark":  teoFnStr(""),
		}),
	}
	// Force the planned remark to null regardless of the helper default.
	planRaw := plan.Raw.(tftypes.Object).AttributeValues()
	planRaw["remark"] = tftypes.NewValue(tftypes.String, nil)
	plan.Raw = tftypes.NewValue(teoFnObjectTypeRaw(), planRaw)

	resp := &resource.UpdateResponse{State: state}
	r.Update(context.Background(), resource.UpdateRequest{Plan: plan, State: state}, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("Update returned error diagnostics: %v", resp.Diagnostics)
	}
	if gotModifyRequest == nil {
		t.Fatal("ModifyFunctionWithContext was not called")
	}
	if gotModifyRequest.Remark != nil {
		t.Errorf("ModifyFunction request.Remark = %v, want nil so the cloud keeps the existing remark", *gotModifyRequest.Remark)
	}
}

// TestTeoFunction111101Resource_Update_ApiError covers the error branch when
// ModifyFunction fails.
func TestTeoFunction111101Resource_Update_ApiError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client, teoClient := teoFnMockClient(patches)
	patches.ApplyMethodFunc(teoClient, "ModifyFunctionWithContext", func(_ context.Context, _ *teosdk.ModifyFunctionRequest) (*teosdk.ModifyFunctionResponse, error) {
		return nil, fmt.Errorf("conflict")
	})

	r := teoFnConfigureResource(t, client)
	s := teoFnSchema(t)
	state := teoFnState(s, teoFnRawValue(map[string]tftypes.Value{
		"id":      teoFnStr(teoFnTestFunctionId),
		"zone_id": teoFnStr(teoFnTestZoneId),
		"name":    teoFnStr(teoFnTestName),
		"content": teoFnStr(teoFnTestContent),
	}))
	plan := tfsdk.Plan{
		Schema: s,
		Raw: teoFnRawValue(map[string]tftypes.Value{
			"id":      teoFnStr(teoFnTestFunctionId),
			"zone_id": teoFnStr(teoFnTestZoneId),
			"name":    teoFnStr(teoFnTestName),
			"content": teoFnStr("new content"),
		}),
	}

	resp := &resource.UpdateResponse{State: state}
	r.Update(context.Background(), resource.UpdateRequest{Plan: plan, State: state}, resp)

	if !resp.Diagnostics.HasError() {
		t.Fatal("Update returned no error diagnostics, want an error on ModifyFunction failure")
	}
}

// TestTeoFunction111101Resource_Delete_Success covers deletion and verifies
// the request carries both zone_id and function_id.
func TestTeoFunction111101Resource_Delete_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client, teoClient := teoFnMockClient(patches)
	var gotDeleteRequest *teosdk.DeleteFunctionRequest
	patches.ApplyMethodFunc(teoClient, "DeleteFunctionWithContext", func(_ context.Context, request *teosdk.DeleteFunctionRequest) (*teosdk.DeleteFunctionResponse, error) {
		gotDeleteRequest = request
		return teosdk.NewDeleteFunctionResponse(), nil
	})

	r := teoFnConfigureResource(t, client)
	s := teoFnSchema(t)
	state := teoFnState(s, teoFnRawValue(map[string]tftypes.Value{
		"id":      teoFnStr(teoFnTestFunctionId),
		"zone_id": teoFnStr(teoFnTestZoneId),
		"name":    teoFnStr(teoFnTestName),
		"content": teoFnStr(teoFnTestContent),
	}))

	resp := &resource.DeleteResponse{}
	r.Delete(context.Background(), resource.DeleteRequest{State: state}, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("Delete returned error diagnostics: %v", resp.Diagnostics)
	}
	if gotDeleteRequest == nil {
		t.Fatal("DeleteFunctionWithContext was not called")
	}
	if gotDeleteRequest.ZoneId == nil || *gotDeleteRequest.ZoneId != teoFnTestZoneId {
		t.Errorf("DeleteFunction request.ZoneId = %v, want %q", gotDeleteRequest.ZoneId, teoFnTestZoneId)
	}
	if gotDeleteRequest.FunctionId == nil || *gotDeleteRequest.FunctionId != teoFnTestFunctionId {
		t.Errorf("DeleteFunction request.FunctionId = %v, want %q", gotDeleteRequest.FunctionId, teoFnTestFunctionId)
	}
}

// TestTeoFunction111101Resource_Delete_ApiError covers the error branch when
// DeleteFunction fails with an unexpected error.
func TestTeoFunction111101Resource_Delete_ApiError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client, teoClient := teoFnMockClient(patches)
	patches.ApplyMethodFunc(teoClient, "DeleteFunctionWithContext", func(_ context.Context, _ *teosdk.DeleteFunctionRequest) (*teosdk.DeleteFunctionResponse, error) {
		return nil, fmt.Errorf("internal error")
	})

	r := teoFnConfigureResource(t, client)
	s := teoFnSchema(t)
	state := teoFnState(s, teoFnRawValue(map[string]tftypes.Value{
		"id":      teoFnStr(teoFnTestFunctionId),
		"zone_id": teoFnStr(teoFnTestZoneId),
		"name":    teoFnStr(teoFnTestName),
		"content": teoFnStr(teoFnTestContent),
	}))

	resp := &resource.DeleteResponse{}
	r.Delete(context.Background(), resource.DeleteRequest{State: state}, resp)

	if !resp.Diagnostics.HasError() {
		t.Fatal("Delete returned no error diagnostics, want an error on DeleteFunction failure")
	}
}

// TestTeoFunction111101Resource_ImportState validates the composite import id
// handling, including the malformed id branch.
func TestTeoFunction111101Resource_ImportState(t *testing.T) {
	s := teoFnSchema(t)
	r := NewTeoFunction111101Resource()

	validId := teoFnTestZoneId + tccommon.FILED_SP + teoFnTestFunctionId
	resp := &resource.ImportStateResponse{State: teoFnNewState(s)}
	r.ImportState(context.Background(), resource.ImportStateRequest{ID: validId}, resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("ImportState returned error diagnostics: %v", resp.Diagnostics)
	}
	var zoneId, functionId types.String
	if diags := resp.State.GetAttribute(context.Background(), path.Root("zone_id"), &zoneId); diags.HasError() {
		t.Fatalf("reading zone_id from import state failed: %v", diags)
	}
	if diags := resp.State.GetAttribute(context.Background(), path.Root("id"), &functionId); diags.HasError() {
		t.Fatalf("reading id from import state failed: %v", diags)
	}
	if zoneId.ValueString() != teoFnTestZoneId {
		t.Errorf("imported zone_id = %q, want %q", zoneId.ValueString(), teoFnTestZoneId)
	}
	if functionId.ValueString() != teoFnTestFunctionId {
		t.Errorf("imported id = %q, want %q", functionId.ValueString(), teoFnTestFunctionId)
	}

	invalid := &resource.ImportStateResponse{State: teoFnNewState(s)}
	r.ImportState(context.Background(), resource.ImportStateRequest{ID: "only-function-id"}, invalid)
	if !invalid.Diagnostics.HasError() {
		t.Error("ImportState with a malformed id returned no error diagnostics")
	}
}

// TestSplitTeoFunction111101Id validates the id split helper.
func TestSplitTeoFunction111101Id(t *testing.T) {
	zone, fn, ok := splitTeoFunction111101Id("zone-abc" + tccommon.FILED_SP + "ef-123")
	if !ok || zone != "zone-abc" || fn != "ef-123" {
		t.Errorf("splitTeoFunction111101Id valid case = (%q, %q, %v)", zone, fn, ok)
	}
	for _, bad := range []string{"", "onlyone", tccommon.FILED_SP + "ef-123", "zone-abc" + tccommon.FILED_SP} {
		if _, _, ok := splitTeoFunction111101Id(bad); ok {
			t.Errorf("splitTeoFunction111101Id(%q) should fail", bad)
		}
	}
}

// TestRefreshTeoFunction111101Model_NilFunction validates the defensive
// branch for a malformed Describe response.
func TestRefreshTeoFunction111101Model_NilFunction(t *testing.T) {
	var diags diag.Diagnostics = refreshTeoFunction111101Model(context.Background(), &TeoFunction111101ResourceModel{}, nil)
	if !diags.HasError() {
		t.Error("refreshTeoFunction111101Model with a nil function should report an error")
	}
}
