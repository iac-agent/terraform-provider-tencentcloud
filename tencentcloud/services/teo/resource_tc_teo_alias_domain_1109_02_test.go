package teo_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/assert"

	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/sharedmeta"
	providerTeo "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/teo"
)

func ptrStringAliasDomain110902(s string) *string { return &s }
func ptrInt64AliasDomain110902(n int64) *int64    { return &n }

// teoAliasDomain110902Schema builds the resource schema via the resource's own
// Schema method.
func teoAliasDomain110902Schema(t *testing.T) schema.Schema {
	r := providerTeo.NewTeoAliasDomain110902Resource()
	schemaResp := &resource.SchemaResponse{}
	r.Schema(context.Background(), resource.SchemaRequest{}, schemaResp)
	if schemaResp.Diagnostics.HasError() {
		t.Fatalf("failed to build resource schema: %v", schemaResp.Diagnostics)
	}
	return schemaResp.Schema
}

// teoAliasDomain110902AttrTypes returns the tftypes attribute types of the
// resource object.
func teoAliasDomain110902AttrTypes() map[string]tftypes.Type {
	return map[string]tftypes.Type{
		"id":          tftypes.String,
		"zone_id":     tftypes.String,
		"alias_name":  tftypes.String,
		"target_name": tftypes.String,
		"cert_type":   tftypes.String,
		"cert_id":     tftypes.Set{ElementType: tftypes.String},
		"alias_domains": tftypes.List{ElementType: tftypes.Object{
			AttributeTypes: map[string]tftypes.Type{
				"alias_name":  tftypes.String,
				"zone_id":     tftypes.String,
				"target_name": tftypes.String,
				"status":      tftypes.String,
				"forbid_mode": tftypes.Number,
				"created_on":  tftypes.String,
				"modified_on": tftypes.String,
			},
		}},
	}
}

// teoAliasDomain110902ObjValue builds the whole resource object as one
// tftypes.Value from a Go model, mirroring how the provider writes state.
func teoAliasDomain110902ObjValue(t *testing.T, model providerTeo.TeoAliasDomain110902ResourceModel) tftypes.Value {
	var certIdVal tftypes.Value
	if model.CertId.IsNull() {
		certIdVal = tftypes.NewValue(tftypes.Set{ElementType: tftypes.String}, nil)
	} else {
		elems := make([]tftypes.Value, 0, len(model.CertId.Elements()))
		for _, e := range model.CertId.Elements() {
			elems = append(elems, tftypes.NewValue(tftypes.String, e.(types.String).ValueString()))
		}
		certIdVal = tftypes.NewValue(tftypes.Set{ElementType: tftypes.String}, elems)
	}

	var aliasDomainsVal tftypes.Value
	if model.AliasDomains == nil {
		aliasDomainsVal = tftypes.NewValue(tftypes.List{ElementType: teoAliasDomain110902ElemType()}, nil)
	} else {
		elems := make([]tftypes.Value, 0, len(model.AliasDomains))
		for _, item := range model.AliasDomains {
			elems = append(elems, teoAliasDomain110902ElemValue(item))
		}
		aliasDomainsVal = tftypes.NewValue(tftypes.List{ElementType: teoAliasDomain110902ElemType()}, elems)
	}

	str := func(v types.String) tftypes.Value {
		if v.IsNull() || v.IsUnknown() {
			return tftypes.NewValue(tftypes.String, nil)
		}
		return tftypes.NewValue(tftypes.String, v.ValueString())
	}

	return tftypes.NewValue(tftypes.Object{
		AttributeTypes: teoAliasDomain110902AttrTypes(),
	}, map[string]tftypes.Value{
		"id":            str(model.Id),
		"zone_id":       str(model.ZoneId),
		"alias_name":    str(model.AliasName),
		"target_name":   str(model.TargetName),
		"cert_type":     str(model.CertType),
		"cert_id":       certIdVal,
		"alias_domains": aliasDomainsVal,
	})
}

// teoAliasDomain110902ElemType returns the tftypes type of one
// alias_domains list element.
func teoAliasDomain110902ElemType() tftypes.Object {
	return tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"alias_name":  tftypes.String,
			"zone_id":     tftypes.String,
			"target_name": tftypes.String,
			"status":      tftypes.String,
			"forbid_mode": tftypes.Number,
			"created_on":  tftypes.String,
			"modified_on": tftypes.String,
		},
	}
}

// teoAliasDomain110902ElemValue builds one alias_domains element value.
func teoAliasDomain110902ElemValue(item providerTeo.TeoAliasDomainItem) tftypes.Value {
	str := func(v types.String) tftypes.Value {
		if v.IsNull() || v.IsUnknown() {
			return tftypes.NewValue(tftypes.String, nil)
		}
		return tftypes.NewValue(tftypes.String, v.ValueString())
	}
	num := func(v types.Int64) tftypes.Value {
		if v.IsNull() || v.IsUnknown() {
			return tftypes.NewValue(tftypes.Number, nil)
		}
		return tftypes.NewValue(tftypes.Number, v.ValueInt64())
	}
	return tftypes.NewValue(teoAliasDomain110902ElemType(), map[string]tftypes.Value{
		"alias_name":  str(item.AliasName),
		"zone_id":     str(item.ZoneId),
		"target_name": str(item.TargetName),
		"status":      str(item.Status),
		"forbid_mode": num(item.ForbidMode),
		"created_on":  str(item.CreatedOn),
		"modified_on": str(item.ModifiedOn),
	})
}

// setTeoAliasDomain110902Client injects a mock connectivity client into the
// resource instance via the Configure method.
func setTeoAliasDomain110902Client(r *providerTeo.TeoAliasDomain110902Resource, client *connectivity.TencentCloudClient) {
	meta := &sharedmeta.ProviderMeta{Client: client}
	configureReq := resource.ConfigureRequest{ProviderData: meta}
	configureResp := &resource.ConfigureResponse{}
	r.Configure(context.Background(), configureReq, configureResp)
}

// mockTeoAliasDomain110902SDK mocks the teo SDK client methods used by the
// resource: create / modify / delete return canned responses, describe returns
// the given alias domain list (nil item means "not found").
func mockTeoAliasDomain110902SDK(patches *gomonkey.Patches, describeItems []*teov20220901.AliasDomain) *teov20220901.Client {
	client := &connectivity.TencentCloudClient{}
	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "CreateAliasDomainWithContext", func(_ context.Context, request *teov20220901.CreateAliasDomainRequest) (*teov20220901.CreateAliasDomainResponse, error) {
		resp := teov20220901.NewCreateAliasDomainResponse()
		resp.Response = &teov20220901.CreateAliasDomainResponseParams{
			RequestId: ptrStringAliasDomain110902("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(teoClient, "ModifyAliasDomainWithContext", func(_ context.Context, request *teov20220901.ModifyAliasDomainRequest) (*teov20220901.ModifyAliasDomainResponse, error) {
		resp := teov20220901.NewModifyAliasDomainResponse()
		resp.Response = &teov20220901.ModifyAliasDomainResponseParams{
			RequestId: ptrStringAliasDomain110902("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(teoClient, "DeleteAliasDomainWithContext", func(_ context.Context, request *teov20220901.DeleteAliasDomainRequest) (*teov20220901.DeleteAliasDomainResponse, error) {
		resp := teov20220901.NewDeleteAliasDomainResponse()
		resp.Response = &teov20220901.DeleteAliasDomainResponseParams{
			RequestId: ptrStringAliasDomain110902("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(teoClient, "DescribeAliasDomainsWithContext", func(_ context.Context, request *teov20220901.DescribeAliasDomainsRequest) (*teov20220901.DescribeAliasDomainsResponse, error) {
		resp := teov20220901.NewDescribeAliasDomainsResponse()
		resp.Response = &teov20220901.DescribeAliasDomainsResponseParams{
			TotalCount:   ptrInt64AliasDomain110902(int64(len(describeItems))),
			AliasDomains: describeItems,
			RequestId:    ptrStringAliasDomain110902("fake-request-id"),
		}
		return resp, nil
	})

	return teoClient
}

// buildTeoAliasDomain110902 returns a configured resource with the mock
// client already injected.
func buildTeoAliasDomain110902(patches *gomonkey.Patches, describeItems []*teov20220901.AliasDomain) *providerTeo.TeoAliasDomain110902Resource {
	mockTeoAliasDomain110902SDK(patches, describeItems)
	r := providerTeo.NewTeoAliasDomain110902Resource().(*providerTeo.TeoAliasDomain110902Resource)
	setTeoAliasDomain110902Client(r, &connectivity.TencentCloudClient{})
	return r
}

// buildTeoAliasDomain110902State builds a tfsdk.State carrying the model.
func buildTeoAliasDomain110902State(t *testing.T, model providerTeo.TeoAliasDomain110902ResourceModel) tfsdk.State {
	return tfsdk.State{
		Raw:    teoAliasDomain110902ObjValue(t, model),
		Schema: teoAliasDomain110902Schema(t),
	}
}

// buildTeoAliasDomain110902Plan builds a tfsdk.Plan carrying the model.
func buildTeoAliasDomain110902Plan(t *testing.T, model providerTeo.TeoAliasDomain110902ResourceModel) tfsdk.Plan {
	return tfsdk.Plan{
		Raw:    teoAliasDomain110902ObjValue(t, model),
		Schema: teoAliasDomain110902Schema(t),
	}
}

// readTeoAliasDomain110902State reads the final state model out of the
// response state.
func readTeoAliasDomain110902State(t *testing.T, state *tfsdk.State) providerTeo.TeoAliasDomain110902ResourceModel {
	var model providerTeo.TeoAliasDomain110902ResourceModel
	if diags := state.Get(context.Background(), &model); diags.HasError() {
		t.Fatalf("failed to read back state: %v", diags)
	}
	return model
}

// newTeoAliasDomain110902CertIdSet builds a types.Set of certificate ids.
func newTeoAliasDomain110902CertIdSet(t *testing.T, ids []string) types.Set {
	if ids == nil {
		return types.SetNull(types.StringType)
	}
	elems := make([]attr.Value, 0, len(ids))
	for _, id := range ids {
		elems = append(elems, types.StringValue(id))
	}
	set, diags := types.SetValue(types.StringType, elems)
	if diags.HasError() {
		t.Fatalf("failed to build cert id set: %v", diags)
	}
	return set
}

// buildTeoAliasDomain110902SDKItem builds the canned cloud-side AliasDomain
// returned by DescribeAliasDomains.
func buildTeoAliasDomain110902SDKItem() *teov20220901.AliasDomain {
	return &teov20220901.AliasDomain{
		AliasName:  ptrStringAliasDomain110902("alias.example.com"),
		ZoneId:     ptrStringAliasDomain110902("zone-2rdhxffj"),
		TargetName: ptrStringAliasDomain110902("target.example.com"),
		Status:     ptrStringAliasDomain110902("active"),
		ForbidMode: ptrInt64AliasDomain110902(0),
		CreatedOn:  ptrStringAliasDomain110902("2025-01-01T00:00:00+08:00"),
		ModifiedOn: ptrStringAliasDomain110902("2025-01-02T00:00:00+08:00"),
	}
}

// go test ./tencentcloud/services/teo/ -run "TestTeoAliasDomain110902" -v -count=1 -gcflags="all=-l"

// TestTeoAliasDomain110902_MetadataAndSchema validates metadata and schema flags.
func TestTeoAliasDomain110902_MetadataAndSchema(t *testing.T) {
	r := providerTeo.NewTeoAliasDomain110902Resource()

	metaResp := &resource.MetadataResponse{}
	r.Metadata(context.Background(), resource.MetadataRequest{}, metaResp)
	assert.Equal(t, "tencentcloud_teo_alias_domain_1109_02", metaResp.TypeName)

	schemaResp := &resource.SchemaResponse{}
	r.Schema(context.Background(), resource.SchemaRequest{}, schemaResp)
	assert.False(t, schemaResp.Diagnostics.HasError(), "schema diagnostics: %v", schemaResp.Diagnostics)

	type expect struct {
		required, optional, computed bool
	}
	cases := map[string]expect{
		"id":            {computed: true},
		"zone_id":       {required: true},
		"alias_name":    {required: true},
		"target_name":   {required: true},
		"cert_type":     {optional: true},
		"cert_id":       {optional: true},
		"alias_domains": {computed: true},
	}
	for name, want := range cases {
		attr, ok := schemaResp.Schema.Attributes[name]
		if !ok {
			t.Errorf("attribute %q missing from schema", name)
			continue
		}
		if attr.IsRequired() != want.required {
			t.Errorf("attribute %q Required=%v, want %v", name, attr.IsRequired(), want.required)
		}
		if attr.IsOptional() != want.optional {
			t.Errorf("attribute %q Optional=%v, want %v", name, attr.IsOptional(), want.optional)
		}
		if attr.IsComputed() != want.computed {
			t.Errorf("attribute %q Computed=%v, want %v", name, attr.IsComputed(), want.computed)
		}
	}
}

// TestTeoAliasDomain110902_Create_Success covers a successful create with
// certificate configuration and full read-back of every alias_domains leaf.
func TestTeoAliasDomain110902_Create_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	sdkItem := buildTeoAliasDomain110902SDKItem()
	r := buildTeoAliasDomain110902(patches, []*teov20220901.AliasDomain{sdkItem})

	certIds, diags := types.SetValue(types.StringType, []attr.Value{types.StringValue("cert-xxxx")})
	if diags.HasError() {
		t.Fatalf("failed to build cert id set: %v", diags)
	}

	plan := providerTeo.TeoAliasDomain110902ResourceModel{
		ZoneId:     types.StringValue("zone-2rdhxffj"),
		AliasName:  types.StringValue("alias.example.com"),
		TargetName: types.StringValue("target.example.com"),
		CertType:   types.StringValue("hosting"),
		CertId:     certIds,
	}

	req := resource.CreateRequest{
		Plan: buildTeoAliasDomain110902Plan(t, plan),
	}
	resp := &resource.CreateResponse{
		State: tfsdk.State{
			Schema: teoAliasDomain110902Schema(t),
			Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: teoAliasDomain110902AttrTypes()}, nil),
		},
	}
	r.Create(context.Background(), req, resp)

	assert.False(t, resp.Diagnostics.HasError(), "expected no error diagnostics, got: %v", resp.Diagnostics)

	got := readTeoAliasDomain110902State(t, &resp.State)
	assert.Equal(t, "zone-2rdhxffj#alias.example.com", got.Id.ValueString())
	assert.Equal(t, "zone-2rdhxffj", got.ZoneId.ValueString())
	assert.Equal(t, "alias.example.com", got.AliasName.ValueString())
	assert.Equal(t, "target.example.com", got.TargetName.ValueString())
	assert.Equal(t, "hosting", got.CertType.ValueString())
	// cert_id is not echoed by DescribeAliasDomains: state keeps the plan value.
	assert.Len(t, got.CertId.Elements(), 1)
	assert.Equal(t, "cert-xxxx", got.CertId.Elements()[0].(types.String).ValueString())

	assert.Len(t, got.AliasDomains, 1)
	item := got.AliasDomains[0]
	assert.Equal(t, "alias.example.com", item.AliasName.ValueString())
	assert.Equal(t, "zone-2rdhxffj", item.ZoneId.ValueString())
	assert.Equal(t, "target.example.com", item.TargetName.ValueString())
	assert.Equal(t, "active", item.Status.ValueString())
	assert.Equal(t, int64(0), item.ForbidMode.ValueInt64())
	assert.Equal(t, "2025-01-01T00:00:00+08:00", item.CreatedOn.ValueString())
	assert.Equal(t, "2025-01-02T00:00:00+08:00", item.ModifiedOn.ValueString())
}

// TestTeoAliasDomain110902_Create_APIError verifies create errors surface as
// error diagnostics and no state is written.
func TestTeoAliasDomain110902_Create_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(client, "UseTeoV20220901Client", teoClient)
	patches.ApplyMethodFunc(teoClient, "CreateAliasDomainWithContext", func(_ context.Context, request *teov20220901.CreateAliasDomainRequest) (*teov20220901.CreateAliasDomainResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=FailedOperation.DomainInUsing, Message=domain in using")
	})

	r := providerTeo.NewTeoAliasDomain110902Resource().(*providerTeo.TeoAliasDomain110902Resource)
	setTeoAliasDomain110902Client(r, client)

	plan := providerTeo.TeoAliasDomain110902ResourceModel{
		ZoneId:     types.StringValue("zone-2rdhxffj"),
		AliasName:  types.StringValue("alias.example.com"),
		TargetName: types.StringValue("target.example.com"),
	}
	req := resource.CreateRequest{
		Plan: buildTeoAliasDomain110902Plan(t, plan),
	}
	resp := &resource.CreateResponse{
		State: tfsdk.State{
			Schema: teoAliasDomain110902Schema(t),
			Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: teoAliasDomain110902AttrTypes()}, nil),
		},
	}
	r.Create(context.Background(), req, resp)

	assert.True(t, resp.Diagnostics.HasError())
	assert.Contains(t, resp.Diagnostics.Errors()[0].Detail(), "FailedOperation.DomainInUsing")
}

// TestTeoAliasDomain110902_Read_Success covers a successful refresh,
// including the "keep state value" behavior of cert_type / cert_id which the
// API never echoes back.
func TestTeoAliasDomain110902_Read_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	sdkItem := buildTeoAliasDomain110902SDKItem()
	r := buildTeoAliasDomain110902(patches, []*teov20220901.AliasDomain{sdkItem})

	state := providerTeo.TeoAliasDomain110902ResourceModel{
		Id:           types.StringValue("zone-2rdhxffj#alias.example.com"),
		ZoneId:       types.StringValue("zone-2rdhxffj"),
		AliasName:    types.StringValue("alias.example.com"),
		TargetName:   types.StringValue("target.example.com"),
		CertType:     types.StringValue("hosting"),
		CertId:       newTeoAliasDomain110902CertIdSet(t, []string{"cert-xxxx"}),
		AliasDomains: []providerTeo.TeoAliasDomainItem{},
	}

	req := resource.ReadRequest{
		State: buildTeoAliasDomain110902State(t, state),
	}
	resp := &resource.ReadResponse{
		State: buildTeoAliasDomain110902State(t, state),
	}
	r.Read(context.Background(), req, resp)

	assert.False(t, resp.Diagnostics.HasError(), "expected no error diagnostics, got: %v", resp.Diagnostics)

	got := readTeoAliasDomain110902State(t, &resp.State)
	assert.Equal(t, "zone-2rdhxffj#alias.example.com", got.Id.ValueString())
	assert.Equal(t, "hosting", got.CertType.ValueString())
	assert.Len(t, got.CertId.Elements(), 1)
	assert.Equal(t, "cert-xxxx", got.CertId.Elements()[0].(types.String).ValueString())

	assert.Len(t, got.AliasDomains, 1)
	item := got.AliasDomains[0]
	assert.Equal(t, "alias.example.com", item.AliasName.ValueString())
	assert.Equal(t, "zone-2rdhxffj", item.ZoneId.ValueString())
	assert.Equal(t, "target.example.com", item.TargetName.ValueString())
	assert.Equal(t, "active", item.Status.ValueString())
	assert.Equal(t, int64(0), item.ForbidMode.ValueInt64())
	assert.Equal(t, "2025-01-01T00:00:00+08:00", item.CreatedOn.ValueString())
	assert.Equal(t, "2025-01-02T00:00:00+08:00", item.ModifiedOn.ValueString())
}

// TestTeoAliasDomain110902_Read_NotFound verifies a confirmed-empty describe
// result removes the resource from state.
func TestTeoAliasDomain110902_Read_NotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	r := buildTeoAliasDomain110902(patches, nil)

	state := providerTeo.TeoAliasDomain110902ResourceModel{
		Id:           types.StringValue("zone-2rdhxffj#alias.example.com"),
		ZoneId:       types.StringValue("zone-2rdhxffj"),
		AliasName:    types.StringValue("alias.example.com"),
		TargetName:   types.StringValue("target.example.com"),
		AliasDomains: []providerTeo.TeoAliasDomainItem{},
	}

	req := resource.ReadRequest{
		State: buildTeoAliasDomain110902State(t, state),
	}
	resp := &resource.ReadResponse{
		State: buildTeoAliasDomain110902State(t, state),
	}
	r.Read(context.Background(), req, resp)

	assert.False(t, resp.Diagnostics.HasError(), "expected no error diagnostics, got: %v", resp.Diagnostics)

	var model providerTeo.TeoAliasDomain110902ResourceModel
	if diags := resp.State.Get(context.Background(), &model); diags.HasError() {
		t.Fatalf("failed to read back state: %v", diags)
	}
	assert.True(t, model.Id.IsNull(), "expected the whole state to be nulled out after removal, got id %q", model.Id.ValueString())
}

// TestTeoAliasDomain110902_Read_APIError verifies a describe API failure is
// reported as an error and keeps the state untouched.
func TestTeoAliasDomain110902_Read_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(client, "UseTeoV20220901Client", teoClient)
	patches.ApplyMethodFunc(teoClient, "DescribeAliasDomainsWithContext", func(_ context.Context, request *teov20220901.DescribeAliasDomainsRequest) (*teov20220901.DescribeAliasDomainsResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=UnauthorizedOperation.ZoneNotFound, Message=zone not found")
	})

	r := providerTeo.NewTeoAliasDomain110902Resource().(*providerTeo.TeoAliasDomain110902Resource)
	setTeoAliasDomain110902Client(r, client)

	state := providerTeo.TeoAliasDomain110902ResourceModel{
		Id:           types.StringValue("zone-2rdhxffj#alias.example.com"),
		ZoneId:       types.StringValue("zone-2rdhxffj"),
		AliasName:    types.StringValue("alias.example.com"),
		TargetName:   types.StringValue("target.example.com"),
		AliasDomains: []providerTeo.TeoAliasDomainItem{},
	}

	req := resource.ReadRequest{
		State: buildTeoAliasDomain110902State(t, state),
	}
	resp := &resource.ReadResponse{
		State: buildTeoAliasDomain110902State(t, state),
	}
	r.Read(context.Background(), req, resp)

	assert.True(t, resp.Diagnostics.HasError())
	assert.Contains(t, resp.Diagnostics.Errors()[0].Detail(), "UnauthorizedOperation.ZoneNotFound")
}

// TestTeoAliasDomain110902_Update_Success covers the update path: target_name
// change plus the confirmed "omit cert_type keeps the previous cloud-side
// configuration" contract (a null plan value is not sent and the final state
// takes the plan value, so removing the attribute neither clears the cloud-side
// certificate nor produces a permanent diff).
func TestTeoAliasDomain110902_Update_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	sdkItem := buildTeoAliasDomain110902SDKItem()
	r := buildTeoAliasDomain110902(patches, []*teov20220901.AliasDomain{sdkItem})

	state := providerTeo.TeoAliasDomain110902ResourceModel{
		Id:           types.StringValue("zone-2rdhxffj#alias.example.com"),
		ZoneId:       types.StringValue("zone-2rdhxffj"),
		AliasName:    types.StringValue("alias.example.com"),
		TargetName:   types.StringValue("old-target.example.com"),
		CertType:     types.StringValue("hosting"),
		CertId:       newTeoAliasDomain110902CertIdSet(t, []string{"cert-xxxx"}),
		AliasDomains: []providerTeo.TeoAliasDomainItem{},
	}

	// cert_type and cert_id removed from the configuration: the confirmed
	// "keep" semantics must not send them, and the final state must take the
	// plan values (null) so no permanent diff is produced.
	plan := providerTeo.TeoAliasDomain110902ResourceModel{
		Id:           types.StringValue("zone-2rdhxffj#alias.example.com"),
		ZoneId:       types.StringValue("zone-2rdhxffj"),
		AliasName:    types.StringValue("alias.example.com"),
		TargetName:   types.StringValue("target.example.com"),
		CertType:     types.StringNull(),
		CertId:       types.SetNull(types.StringType),
		AliasDomains: []providerTeo.TeoAliasDomainItem{},
	}

	req := resource.UpdateRequest{
		Plan:  buildTeoAliasDomain110902Plan(t, plan),
		State: buildTeoAliasDomain110902State(t, state),
	}
	resp := &resource.UpdateResponse{
		State: tfsdk.State{
			Schema: teoAliasDomain110902Schema(t),
			Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: teoAliasDomain110902AttrTypes()}, nil),
		},
	}
	r.Update(context.Background(), req, resp)

	assert.False(t, resp.Diagnostics.HasError(), "expected no error diagnostics, got: %v", resp.Diagnostics)

	got := readTeoAliasDomain110902State(t, &resp.State)
	assert.Equal(t, "zone-2rdhxffj#alias.example.com", got.Id.ValueString())
	assert.Equal(t, "target.example.com", got.TargetName.ValueString())
	// Keep semantics: the omitted fields are not sent; the final state takes
	// the plan values (null), so the next plan shows no diff.
	assert.True(t, got.CertType.IsNull(), "expected cert_type to follow the plan (null), got %q", got.CertType.ValueString())
	assert.True(t, got.CertId.IsNull(), "expected cert_id to follow the plan (null)")
	assert.Len(t, got.AliasDomains, 1)
	assert.Equal(t, "target.example.com", got.AliasDomains[0].TargetName.ValueString())
}

// TestTeoAliasDomain110902_Update_WireRequest verifies the ModifyAliasDomain
// wire request: identity fields always sent, omitted cert fields never sent,
// explicit cert fields sent verbatim.
func TestTeoAliasDomain110902_Update_WireRequest(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	sdkItem := buildTeoAliasDomain110902SDKItem()
	client := &connectivity.TencentCloudClient{}
	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(client, "UseTeoV20220901Client", teoClient)

	var gotRequest *teov20220901.ModifyAliasDomainRequest
	patches.ApplyMethodFunc(teoClient, "ModifyAliasDomainWithContext", func(_ context.Context, request *teov20220901.ModifyAliasDomainRequest) (*teov20220901.ModifyAliasDomainResponse, error) {
		gotRequest = request
		resp := teov20220901.NewModifyAliasDomainResponse()
		resp.Response = &teov20220901.ModifyAliasDomainResponseParams{
			RequestId: ptrStringAliasDomain110902("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(teoClient, "DescribeAliasDomainsWithContext", func(_ context.Context, request *teov20220901.DescribeAliasDomainsRequest) (*teov20220901.DescribeAliasDomainsResponse, error) {
		resp := teov20220901.NewDescribeAliasDomainsResponse()
		resp.Response = &teov20220901.DescribeAliasDomainsResponseParams{
			TotalCount:   ptrInt64AliasDomain110902(1),
			AliasDomains: []*teov20220901.AliasDomain{sdkItem},
			RequestId:    ptrStringAliasDomain110902("fake-request-id"),
		}
		return resp, nil
	})

	r := providerTeo.NewTeoAliasDomain110902Resource().(*providerTeo.TeoAliasDomain110902Resource)
	setTeoAliasDomain110902Client(r, client)

	state := providerTeo.TeoAliasDomain110902ResourceModel{
		Id:           types.StringValue("zone-2rdhxffj#alias.example.com"),
		ZoneId:       types.StringValue("zone-2rdhxffj"),
		AliasName:    types.StringValue("alias.example.com"),
		TargetName:   types.StringValue("old-target.example.com"),
		CertType:     types.StringValue("hosting"),
		CertId:       newTeoAliasDomain110902CertIdSet(t, []string{"cert-xxxx"}),
		AliasDomains: []providerTeo.TeoAliasDomainItem{},
	}
	plan := providerTeo.TeoAliasDomain110902ResourceModel{
		Id:           types.StringValue("zone-2rdhxffj#alias.example.com"),
		ZoneId:       types.StringValue("zone-2rdhxffj"),
		AliasName:    types.StringValue("alias.example.com"),
		TargetName:   types.StringValue("target.example.com"),
		CertType:     types.StringNull(),
		CertId:       types.SetNull(types.StringType),
		AliasDomains: []providerTeo.TeoAliasDomainItem{},
	}

	req := resource.UpdateRequest{
		Plan:  buildTeoAliasDomain110902Plan(t, plan),
		State: buildTeoAliasDomain110902State(t, state),
	}
	resp := &resource.UpdateResponse{
		State: tfsdk.State{
			Schema: teoAliasDomain110902Schema(t),
			Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: teoAliasDomain110902AttrTypes()}, nil),
		},
	}
	r.Update(context.Background(), req, resp)

	assert.False(t, resp.Diagnostics.HasError(), "expected no error diagnostics, got: %v", resp.Diagnostics)
	assert.NotNil(t, gotRequest)
	assert.Equal(t, "zone-2rdhxffj", *gotRequest.ZoneId)
	assert.Equal(t, "alias.example.com", *gotRequest.AliasName)
	assert.Equal(t, "target.example.com", *gotRequest.TargetName)
	// Keep semantics: omitted cert fields must not be sent at all.
	assert.Nil(t, gotRequest.CertType)
	assert.Empty(t, gotRequest.CertId)
}

// TestTeoAliasDomain110902_Update_ExplicitValues covers an update that
// switches the certificate configuration explicitly.
func TestTeoAliasDomain110902_Update_ExplicitValues(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	sdkItem := buildTeoAliasDomain110902SDKItem()
	r := buildTeoAliasDomain110902(patches, []*teov20220901.AliasDomain{sdkItem})

	state := providerTeo.TeoAliasDomain110902ResourceModel{
		Id:           types.StringValue("zone-2rdhxffj#alias.example.com"),
		ZoneId:       types.StringValue("zone-2rdhxffj"),
		AliasName:    types.StringValue("alias.example.com"),
		TargetName:   types.StringValue("target.example.com"),
		CertType:     types.StringValue("hosting"),
		CertId:       newTeoAliasDomain110902CertIdSet(t, []string{"cert-old"}),
		AliasDomains: []providerTeo.TeoAliasDomainItem{},
	}

	plan := providerTeo.TeoAliasDomain110902ResourceModel{
		Id:           types.StringValue("zone-2rdhxffj#alias.example.com"),
		ZoneId:       types.StringValue("zone-2rdhxffj"),
		AliasName:    types.StringValue("alias.example.com"),
		TargetName:   types.StringValue("target.example.com"),
		CertType:     types.StringValue("none"),
		CertId:       types.SetNull(types.StringType),
		AliasDomains: []providerTeo.TeoAliasDomainItem{},
	}

	req := resource.UpdateRequest{
		Plan:  buildTeoAliasDomain110902Plan(t, plan),
		State: buildTeoAliasDomain110902State(t, state),
	}
	resp := &resource.UpdateResponse{
		State: tfsdk.State{
			Schema: teoAliasDomain110902Schema(t),
			Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: teoAliasDomain110902AttrTypes()}, nil),
		},
	}
	r.Update(context.Background(), req, resp)

	assert.False(t, resp.Diagnostics.HasError(), "expected no error diagnostics, got: %v", resp.Diagnostics)

	got := readTeoAliasDomain110902State(t, &resp.State)
	assert.Equal(t, "none", got.CertType.ValueString())
	assert.True(t, got.CertId.IsNull())
}

// TestTeoAliasDomain110902_Update_APIError verifies modify failures surface
// as error diagnostics.
func TestTeoAliasDomain110902_Update_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(client, "UseTeoV20220901Client", teoClient)
	patches.ApplyMethodFunc(teoClient, "ModifyAliasDomainWithContext", func(_ context.Context, request *teov20220901.ModifyAliasDomainRequest) (*teov20220901.ModifyAliasDomainResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=FailedOperation.DomainNotExists, Message=domain not exists")
	})

	r := providerTeo.NewTeoAliasDomain110902Resource().(*providerTeo.TeoAliasDomain110902Resource)
	setTeoAliasDomain110902Client(r, client)

	state := providerTeo.TeoAliasDomain110902ResourceModel{
		Id:           types.StringValue("zone-2rdhxffj#alias.example.com"),
		ZoneId:       types.StringValue("zone-2rdhxffj"),
		AliasName:    types.StringValue("alias.example.com"),
		TargetName:   types.StringValue("target.example.com"),
		AliasDomains: []providerTeo.TeoAliasDomainItem{},
	}
	plan := state

	req := resource.UpdateRequest{
		Plan:  buildTeoAliasDomain110902Plan(t, plan),
		State: buildTeoAliasDomain110902State(t, state),
	}
	resp := &resource.UpdateResponse{
		State: tfsdk.State{
			Schema: teoAliasDomain110902Schema(t),
			Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: teoAliasDomain110902AttrTypes()}, nil),
		},
	}
	r.Update(context.Background(), req, resp)

	assert.True(t, resp.Diagnostics.HasError())
	assert.Contains(t, resp.Diagnostics.Errors()[0].Detail(), "FailedOperation.DomainNotExists")
}

// TestTeoAliasDomain110902_Delete_Success covers deleting an existing alias
// domain and verifies the single-element AliasNames list on the wire.
func TestTeoAliasDomain110902_Delete_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	sdkItem := buildTeoAliasDomain110902SDKItem()
	client := &connectivity.TencentCloudClient{}
	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeAliasDomainsWithContext", func(_ context.Context, request *teov20220901.DescribeAliasDomainsRequest) (*teov20220901.DescribeAliasDomainsResponse, error) {
		resp := teov20220901.NewDescribeAliasDomainsResponse()
		resp.Response = &teov20220901.DescribeAliasDomainsResponseParams{
			TotalCount:   ptrInt64AliasDomain110902(1),
			AliasDomains: []*teov20220901.AliasDomain{sdkItem},
			RequestId:    ptrStringAliasDomain110902("fake-request-id"),
		}
		return resp, nil
	})

	var gotDeleteRequest *teov20220901.DeleteAliasDomainRequest
	patches.ApplyMethodFunc(teoClient, "DeleteAliasDomainWithContext", func(_ context.Context, request *teov20220901.DeleteAliasDomainRequest) (*teov20220901.DeleteAliasDomainResponse, error) {
		gotDeleteRequest = request
		resp := teov20220901.NewDeleteAliasDomainResponse()
		resp.Response = &teov20220901.DeleteAliasDomainResponseParams{
			RequestId: ptrStringAliasDomain110902("fake-request-id"),
		}
		return resp, nil
	})

	r := providerTeo.NewTeoAliasDomain110902Resource().(*providerTeo.TeoAliasDomain110902Resource)
	setTeoAliasDomain110902Client(r, client)

	state := providerTeo.TeoAliasDomain110902ResourceModel{
		Id:           types.StringValue("zone-2rdhxffj#alias.example.com"),
		ZoneId:       types.StringValue("zone-2rdhxffj"),
		AliasName:    types.StringValue("alias.example.com"),
		TargetName:   types.StringValue("target.example.com"),
		AliasDomains: []providerTeo.TeoAliasDomainItem{},
	}

	req := resource.DeleteRequest{
		State: buildTeoAliasDomain110902State(t, state),
	}
	resp := &resource.DeleteResponse{}
	r.Delete(context.Background(), req, resp)

	assert.False(t, resp.Diagnostics.HasError(), "expected no error diagnostics, got: %v", resp.Diagnostics)
	assert.NotNil(t, gotDeleteRequest)
	assert.Equal(t, "zone-2rdhxffj", *gotDeleteRequest.ZoneId)
	assert.Len(t, gotDeleteRequest.AliasNames, 1)
	assert.Equal(t, "alias.example.com", *gotDeleteRequest.AliasNames[0])
}

// TestTeoAliasDomain110902_Delete_AlreadyGone verifies that deleting a
// confirmed-missing alias domain succeeds without calling DeleteAliasDomain.
func TestTeoAliasDomain110902_Delete_AlreadyGone(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeAliasDomainsWithContext", func(_ context.Context, request *teov20220901.DescribeAliasDomainsRequest) (*teov20220901.DescribeAliasDomainsResponse, error) {
		resp := teov20220901.NewDescribeAliasDomainsResponse()
		resp.Response = &teov20220901.DescribeAliasDomainsResponseParams{
			TotalCount:   ptrInt64AliasDomain110902(0),
			AliasDomains: []*teov20220901.AliasDomain{},
			RequestId:    ptrStringAliasDomain110902("fake-request-id"),
		}
		return resp, nil
	})

	deleteCalled := false
	patches.ApplyMethodFunc(teoClient, "DeleteAliasDomainWithContext", func(_ context.Context, request *teov20220901.DeleteAliasDomainRequest) (*teov20220901.DeleteAliasDomainResponse, error) {
		deleteCalled = true
		resp := teov20220901.NewDeleteAliasDomainResponse()
		resp.Response = &teov20220901.DeleteAliasDomainResponseParams{
			RequestId: ptrStringAliasDomain110902("fake-request-id"),
		}
		return resp, nil
	})

	r := providerTeo.NewTeoAliasDomain110902Resource().(*providerTeo.TeoAliasDomain110902Resource)
	setTeoAliasDomain110902Client(r, client)

	state := providerTeo.TeoAliasDomain110902ResourceModel{
		Id:           types.StringValue("zone-2rdhxffj#alias.example.com"),
		ZoneId:       types.StringValue("zone-2rdhxffj"),
		AliasName:    types.StringValue("alias.example.com"),
		TargetName:   types.StringValue("target.example.com"),
		AliasDomains: []providerTeo.TeoAliasDomainItem{},
	}

	req := resource.DeleteRequest{
		State: buildTeoAliasDomain110902State(t, state),
	}
	resp := &resource.DeleteResponse{}
	r.Delete(context.Background(), req, resp)

	assert.False(t, resp.Diagnostics.HasError(), "expected no error diagnostics, got: %v", resp.Diagnostics)
	assert.False(t, deleteCalled, "DeleteAliasDomain must not be called when the alias domain is confirmed missing")
}

// TestTeoAliasDomain110902_Delete_APIError verifies delete failures surface
// as error diagnostics so the resource stays in state for a retry.
func TestTeoAliasDomain110902_Delete_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	sdkItem := buildTeoAliasDomain110902SDKItem()
	client := &connectivity.TencentCloudClient{}
	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeAliasDomainsWithContext", func(_ context.Context, request *teov20220901.DescribeAliasDomainsRequest) (*teov20220901.DescribeAliasDomainsResponse, error) {
		resp := teov20220901.NewDescribeAliasDomainsResponse()
		resp.Response = &teov20220901.DescribeAliasDomainsResponseParams{
			TotalCount:   ptrInt64AliasDomain110902(1),
			AliasDomains: []*teov20220901.AliasDomain{sdkItem},
			RequestId:    ptrStringAliasDomain110902("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(teoClient, "DeleteAliasDomainWithContext", func(_ context.Context, request *teov20220901.DeleteAliasDomainRequest) (*teov20220901.DeleteAliasDomainResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=FailedOperation.DomainNotExists, Message=domain not exists")
	})

	r := providerTeo.NewTeoAliasDomain110902Resource().(*providerTeo.TeoAliasDomain110902Resource)
	setTeoAliasDomain110902Client(r, client)

	state := providerTeo.TeoAliasDomain110902ResourceModel{
		Id:           types.StringValue("zone-2rdhxffj#alias.example.com"),
		ZoneId:       types.StringValue("zone-2rdhxffj"),
		AliasName:    types.StringValue("alias.example.com"),
		TargetName:   types.StringValue("target.example.com"),
		AliasDomains: []providerTeo.TeoAliasDomainItem{},
	}

	req := resource.DeleteRequest{
		State: buildTeoAliasDomain110902State(t, state),
	}
	resp := &resource.DeleteResponse{}
	r.Delete(context.Background(), req, resp)

	assert.True(t, resp.Diagnostics.HasError())
	assert.Contains(t, resp.Diagnostics.Errors()[0].Detail(), "FailedOperation.DomainNotExists")
}

// TestTeoAliasDomain110902_ImportState covers import id passthrough and the
// composite id format validation.
func TestTeoAliasDomain110902_ImportState(t *testing.T) {
	r := providerTeo.NewTeoAliasDomain110902Resource().(*providerTeo.TeoAliasDomain110902Resource)

	// Valid composite id: passthrough into the id attribute.
	resp := &resource.ImportStateResponse{
		State: tfsdk.State{
			Schema: teoAliasDomain110902Schema(t),
			Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: teoAliasDomain110902AttrTypes()}, nil),
		},
	}
	r.ImportState(context.Background(), resource.ImportStateRequest{ID: "zone-2rdhxffj#alias.example.com"}, resp)
	assert.False(t, resp.Diagnostics.HasError(), "expected no error diagnostics, got: %v", resp.Diagnostics)

	var id types.String
	if diags := resp.State.GetAttribute(context.Background(), path.Root("id"), &id); diags.HasError() {
		t.Fatalf("failed to read back id: %v", diags)
	}
	assert.Equal(t, "zone-2rdhxffj#alias.example.com", id.ValueString())

	// Broken ids must be rejected, never silently half-imported.
	for _, bad := range []string{"", "zone-2rdhxffj", "#alias.example.com", "zone-2rdhxffj#", "a#b#c"} {
		badResp := &resource.ImportStateResponse{
			State: tfsdk.State{
				Schema: teoAliasDomain110902Schema(t),
				Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: teoAliasDomain110902AttrTypes()}, nil),
			},
		}
		r.ImportState(context.Background(), resource.ImportStateRequest{ID: bad}, badResp)
		assert.True(t, badResp.Diagnostics.HasError(), "expected error for import id %q", bad)
	}
}

// TestTeoAliasDomain110902_SplitId covers the composite id parser.
func TestTeoAliasDomain110902_SplitId(t *testing.T) {
	good := "zone-2rdhxffj#alias.example.com"
	zone, alias, ok := providerTeo.SplitTeoAliasDomain110902Id(good)
	if !ok || zone != "zone-2rdhxffj" || alias != "alias.example.com" {
		t.Errorf("SplitTeoAliasDomain110902Id(%q) = (%q, %q, %v)", good, zone, alias, ok)
	}
	for _, bad := range []string{"", "zone-only", "#alias", "zone#", "a#b#c"} {
		if _, _, ok := providerTeo.SplitTeoAliasDomain110902Id(bad); ok {
			t.Errorf("SplitTeoAliasDomain110902Id(%q) ok=true, want false", bad)
		}
	}
}
