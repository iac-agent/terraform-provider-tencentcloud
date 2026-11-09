package teo

import (
	"context"
	"fmt"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	fwtypes "github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/assert"
	teo "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/sharedmeta"
)

// ptrStringAliasDomain110904 returns a pointer to the given string.
func ptrStringAliasDomain110904(s string) *string { return &s }

// ptrInt64AliasDomain110904 returns a pointer to the given int64.
func ptrInt64AliasDomain110904(i int64) *int64 { return &i }

// aliasDomain110904Schema returns the framework resource schema, failing the
// test when Schema returns diagnostics.
func aliasDomain110904Schema(t *testing.T) schema.Schema {
	r := &TeoAliasDomain110904Resource{}
	req := resource.SchemaRequest{}
	resp := &resource.SchemaResponse{}
	r.Schema(context.Background(), req, resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Schema returned diagnostics: %v", resp.Diagnostics)
	}
	return resp.Schema
}

// aliasDomain110904Str builds a tftypes string value.
func aliasDomain110904Str(v string) tftypes.Value {
	return tftypes.NewValue(tftypes.String, v)
}

// aliasDomain110904NullStr builds a null tftypes string value.
func aliasDomain110904NullStr() tftypes.Value {
	return tftypes.NewValue(tftypes.String, nil)
}

// aliasDomain110904StrSet builds a tftypes set of strings from the items.
func aliasDomain110904StrSet(items []string) tftypes.Value {
	if items == nil {
		return tftypes.NewValue(tftypes.Set{ElementType: tftypes.String}, nil)
	}
	vals := make([]tftypes.Value, 0, len(items))
	for _, item := range items {
		vals = append(vals, aliasDomain110904Str(item))
	}
	return tftypes.NewValue(tftypes.Set{ElementType: tftypes.String}, vals)
}

// aliasDomain110904UnknownStr builds an unknown tftypes string value.
func aliasDomain110904UnknownStr() tftypes.Value {
	return tftypes.NewValue(tftypes.String, tftypes.UnknownValue)
}

// aliasDomain110904UnknownStrSet builds an unknown tftypes set-of-string value.
func aliasDomain110904UnknownStrSet() tftypes.Value {
	return tftypes.NewValue(tftypes.Set{ElementType: tftypes.String}, tftypes.UnknownValue)
}

// aliasDomain110904RawObject builds the raw tftypes object value of the whole
// resource model; unspecified attributes default to null.
func aliasDomain110904RawObject(s schema.Schema, attrs map[string]tftypes.Value) tftypes.Value {
	all := map[string]tftypes.Value{
		"id":          aliasDomain110904NullStr(),
		"zone_id":     aliasDomain110904NullStr(),
		"alias_name":  aliasDomain110904NullStr(),
		"target_name": aliasDomain110904NullStr(),
		"cert_type":   aliasDomain110904NullStr(),
		"cert_id":     aliasDomain110904StrSet(nil),
		"status":      aliasDomain110904NullStr(),
	}
	for k, v := range attrs {
		all[k] = v
	}
	return tftypes.NewValue(s.Type().TerraformType(context.Background()), all)
}

// tftypesDataAliasDomain110904 builds a tfsdk data wrapper with the given schema
// and raw value.
func tftypesDataAliasDomain110904(s schema.Schema, raw tftypes.Value) tfsdk.State {
	return tfsdk.State{Schema: s, Raw: raw}
}

// setAliasDomain110904Client injects a mock connectivity client into the
// resource instance via the Configure method.
func setAliasDomain110904Client(r *TeoAliasDomain110904Resource, client *connectivity.TencentCloudClient) {
	meta := &sharedmeta.ProviderMeta{Client: client}
	configureReq := resource.ConfigureRequest{ProviderData: meta}
	configureResp := &resource.ConfigureResponse{}
	r.Configure(context.Background(), configureReq, configureResp)
}

// mockTeoClientForAliasDomain110904 patches the connectivity client to return
// the given teo client from UseTeoV20220901Client.
func mockTeoClientForAliasDomain110904(patches *gomonkey.Patches, client *connectivity.TencentCloudClient) *teo.Client {
	teoClient := &teo.Client{}
	patches.ApplyMethodReturn(client, "UseTeoV20220901Client", teoClient)
	return teoClient
}

// describeAliasDomain110904Response builds a DescribeAliasDomains response
// containing the given alias domain.
func describeAliasDomain110904Response(item *teo.AliasDomain) *teo.DescribeAliasDomainsResponse {
	resp := teo.NewDescribeAliasDomainsResponse()
	resp.Response = &teo.DescribeAliasDomainsResponseParams{
		TotalCount:   ptrInt64AliasDomain110904(1),
		AliasDomains: []*teo.AliasDomain{item},
		RequestId:    ptrStringAliasDomain110904("fake-request-id"),
	}
	return resp
}

// emptyAliasDomain110904Response builds a DescribeAliasDomains response that
// contains no alias domains.
func emptyAliasDomain110904Response() *teo.DescribeAliasDomainsResponse {
	resp := teo.NewDescribeAliasDomainsResponse()
	resp.Response = &teo.DescribeAliasDomainsResponseParams{
		TotalCount:   ptrInt64AliasDomain110904(0),
		AliasDomains: []*teo.AliasDomain{},
		RequestId:    ptrStringAliasDomain110904("fake-request-id"),
	}
	return resp
}

// sampleAliasDomain110904Active builds an active alias domain payload.
func sampleAliasDomain110904Active() *teo.AliasDomain {
	return &teo.AliasDomain{
		AliasName:   ptrStringAliasDomain110904("alias.demo.com"),
		ZoneId:      ptrStringAliasDomain110904("zone-3fkff38fyw8s"),
		TargetName:  ptrStringAliasDomain110904("target.demo.com"),
		Status:      ptrStringAliasDomain110904("active"),
		ForbidMode:  ptrInt64AliasDomain110904(0),
		CreatedOn:   ptrStringAliasDomain110904("2024-01-01T00:00:00Z"),
		ModifiedOn:  ptrStringAliasDomain110904("2024-01-01T00:00:00Z"),
	}
}

// TestTeoAliasDomain110904Resource_Metadata validates the resource type name.
func TestTeoAliasDomain110904Resource_Metadata(t *testing.T) {
	r := &TeoAliasDomain110904Resource{}
	req := resource.MetadataRequest{}
	resp := &resource.MetadataResponse{}

	r.Metadata(context.Background(), req, resp)

	if resp.TypeName != "tencentcloud_teo_alias_domain_1109_04" {
		t.Errorf("Metadata.TypeName = %q, want %q", resp.TypeName, "tencentcloud_teo_alias_domain_1109_04")
	}
}

// TestTeoAliasDomain110904Resource_Schema validates required / optional /
// computed flags on every attribute.
func TestTeoAliasDomain110904Resource_Schema(t *testing.T) {
	r := &TeoAliasDomain110904Resource{}
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
		"id":          {computed: true},
		"zone_id":     {required: true},
		"alias_name":  {required: true},
		"target_name": {required: true},
		"cert_type":   {optional: true, computed: true},
		"status":      {computed: true},
	}
	for name, want := range cases {
		attr, ok := resp.Schema.Attributes[name].(schema.StringAttribute)
		if !ok {
			t.Errorf("attribute %q missing or not a StringAttribute", name)
			continue
		}
		if attr.Required != want.required {
			t.Errorf("attribute %q Required=%v, want %v", name, attr.Required, want.required)
		}
		if attr.Optional != want.optional {
			t.Errorf("attribute %q Optional=%v, want %v", name, attr.Optional, want.optional)
		}
		if attr.Computed != want.computed {
			t.Errorf("attribute %q Computed=%v, want %v", name, attr.Computed, want.computed)
		}
	}

	certId, ok := resp.Schema.Attributes["cert_id"].(schema.SetAttribute)
	if !ok {
		t.Fatalf("attribute %q missing or not a SetAttribute", "cert_id")
	}
	if !certId.Optional {
		t.Errorf("attribute %q Optional=%v, want true", "cert_id", certId.Optional)
	}
	if certId.Computed {
		t.Errorf("attribute %q Computed=%v, want false", "cert_id", certId.Computed)
	}
	if !certId.ElementType.Equal(fwtypes.StringType) {
		t.Errorf("attribute %q ElementType=%v, want StringType", "cert_id", certId.ElementType)
	}
}

// TestTeoAliasDomain110904Resource_ImplementsInterfaces verifies the resource
// satisfies the framework interfaces it advertises.
func TestTeoAliasDomain110904Resource_ImplementsInterfaces(t *testing.T) {
	var _ resource.ResourceWithConfigure = &TeoAliasDomain110904Resource{}
	var _ resource.ResourceWithImportState = &TeoAliasDomain110904Resource{}
}

// TestTeoAliasDomain110904Resource_Factory validates the registry factory.
func TestTeoAliasDomain110904Resource_Factory(t *testing.T) {
	got := NewTeoAliasDomain110904Resource()
	if got == nil {
		t.Fatal("NewTeoAliasDomain110904Resource returned nil")
	}
	if _, ok := got.(*TeoAliasDomain110904Resource); !ok {
		t.Fatalf("factory returned %T, want *TeoAliasDomain110904Resource", got)
	}
}

// TestTeoAliasDomain110904Resource_ConfigureNilProviderData ensures the
// pre-Configure invocation pattern (nil ProviderData) is silently no-op.
func TestTeoAliasDomain110904Resource_ConfigureNilProviderData(t *testing.T) {
	r := &TeoAliasDomain110904Resource{}
	req := resource.ConfigureRequest{}
	resp := &resource.ConfigureResponse{}

	r.Configure(context.Background(), req, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("Configure with nil ProviderData returned diagnostics: %v", resp.Diagnostics)
	}
}

// TestTeoAliasDomain110904Resource_ImportState covers the composite-id import:
// a valid `<zone_id>#<alias_name>` sets state.id, and broken ids are rejected.
func TestTeoAliasDomain110904Resource_ImportState(t *testing.T) {
	r := &TeoAliasDomain110904Resource{}
	s := aliasDomain110904Schema(t)
	compositeId := "zone-3fkff38fyw8s" + tccommon.FILED_SP + "alias.demo.com"

	// Valid composite id: id lands in state via the passthrough.
	req := resource.ImportStateRequest{ID: compositeId}
	resp := &resource.ImportStateResponse{
		State: tftypesDataAliasDomain110904(s, aliasDomain110904RawObject(s, nil)),
	}
	r.ImportState(context.Background(), req, resp)
	assert.False(t, resp.Diagnostics.HasError(), "expected no error diagnostics, got: %v", resp.Diagnostics)

	var state TeoAliasDomain110904ResourceModel
	diags := resp.State.Get(context.Background(), &state)
	assert.False(t, diags.HasError(), "State.Get returned diagnostics: %v", diags)
	assert.Equal(t, compositeId, state.Id.ValueString())

	// Broken ids are rejected with an error diagnostic.
	for _, bad := range []string{
		"",
		"only-zone",
		tccommon.FILED_SP + "alias.demo.com",
		"zone-3fkff38fyw8s" + tccommon.FILED_SP,
		"a" + tccommon.FILED_SP + "b" + tccommon.FILED_SP + "c",
	} {
		badReq := resource.ImportStateRequest{ID: bad}
		badResp := &resource.ImportStateResponse{
			State: tftypesDataAliasDomain110904(s, aliasDomain110904RawObject(s, nil)),
		}
		r.ImportState(context.Background(), badReq, badResp)
		assert.True(t, badResp.Diagnostics.HasError(), "import id %q should be rejected", bad)
	}
}

// TestTeoAliasDomain110904Resource_Create_Success covers a successful create:
// identity is persisted before the asynchronous wait, and the refreshed model
// carries every readable leaf field.
func TestTeoAliasDomain110904Resource_Create_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	teoClient := mockTeoClientForAliasDomain110904(patches, client)

	var gotCreate *teo.CreateAliasDomainRequest
	patches.ApplyMethodFunc(teoClient, "CreateAliasDomainWithContext", func(_ context.Context, request *teo.CreateAliasDomainRequest) (*teo.CreateAliasDomainResponse, error) {
		gotCreate = request
		resp := teo.NewCreateAliasDomainResponse()
		resp.Response = &teo.CreateAliasDomainResponseParams{
			RequestId: ptrStringAliasDomain110904("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(teoClient, "DescribeAliasDomainsWithContext", func(_ context.Context, request *teo.DescribeAliasDomainsRequest) (*teo.DescribeAliasDomainsResponse, error) {
		return describeAliasDomain110904Response(sampleAliasDomain110904Active()), nil
	})

	r := &TeoAliasDomain110904Resource{}
	setAliasDomain110904Client(r, client)
	s := aliasDomain110904Schema(t)

	planRaw := aliasDomain110904RawObject(s, map[string]tftypes.Value{
		"zone_id":     aliasDomain110904Str("zone-3fkff38fyw8s"),
		"alias_name":  aliasDomain110904Str("alias.demo.com"),
		"target_name": aliasDomain110904Str("target.demo.com"),
		"cert_type":   aliasDomain110904Str("hosting"),
		"cert_id":     aliasDomain110904StrSet([]string{"cert-xxxx"}),
		"status":      aliasDomain110904NullStr(),
	})

	req := resource.CreateRequest{
		Config: tfsdk.Config{Schema: s, Raw: planRaw},
		Plan:  tfsdk.Plan{Schema: s, Raw: planRaw},
	}
	resp := &resource.CreateResponse{
		State: tftypesDataAliasDomain110904(s, aliasDomain110904RawObject(s, nil)),
	}

	r.Create(context.Background(), req, resp)
	assert.False(t, resp.Diagnostics.HasError(), "expected no error diagnostics, got: %v", resp.Diagnostics)

	// Request mapping: every configurable field lands on the API request.
	assert.NotNil(t, gotCreate)
	assert.Equal(t, "zone-3fkff38fyw8s", *gotCreate.ZoneId)
	assert.Equal(t, "alias.demo.com", *gotCreate.AliasName)
	assert.Equal(t, "target.demo.com", *gotCreate.TargetName)
	assert.Equal(t, "hosting", *gotCreate.CertType)
	assert.Equal(t, []*string{ptrStringAliasDomain110904("cert-xxxx")}, gotCreate.CertId)

	// State mapping: composite id + every readable leaf field.
	var state TeoAliasDomain110904ResourceModel
	diags := resp.State.Get(context.Background(), &state)
	assert.False(t, diags.HasError(), "State.Get returned diagnostics: %v", diags)
	assert.Equal(t, "zone-3fkff38fyw8s"+tccommon.FILED_SP+"alias.demo.com", state.Id.ValueString())
	assert.Equal(t, "zone-3fkff38fyw8s", state.ZoneId.ValueString())
	assert.Equal(t, "alias.demo.com", state.AliasName.ValueString())
	assert.Equal(t, "target.demo.com", state.TargetName.ValueString())
	assert.Equal(t, "hosting", state.CertType.ValueString())
	assert.False(t, state.CertId.IsNull())
	var certIds []fwtypes.String
	certDiags := state.CertId.ElementsAs(context.Background(), &certIds, false)
	assert.False(t, certDiags.HasError(), "CertId.ElementsAs returned diagnostics: %v", certDiags)
	assert.Equal(t, []string{"cert-xxxx"}, stringValuesAliasDomain110904(certIds))
	assert.Equal(t, "active", state.Status.ValueString())
}

// TestTeoAliasDomain110904Resource_Create_NullCert covers a create without
// cert fields: CertType/CertId are not sent.
func TestTeoAliasDomain110904Resource_Create_NullCert(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	teoClient := mockTeoClientForAliasDomain110904(patches, client)

	var gotCreate *teo.CreateAliasDomainRequest
	patches.ApplyMethodFunc(teoClient, "CreateAliasDomainWithContext", func(_ context.Context, request *teo.CreateAliasDomainRequest) (*teo.CreateAliasDomainResponse, error) {
		gotCreate = request
		resp := teo.NewCreateAliasDomainResponse()
		resp.Response = &teo.CreateAliasDomainResponseParams{
			RequestId: ptrStringAliasDomain110904("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(teoClient, "DescribeAliasDomainsWithContext", func(_ context.Context, request *teo.DescribeAliasDomainsRequest) (*teo.DescribeAliasDomainsResponse, error) {
		return describeAliasDomain110904Response(sampleAliasDomain110904Active()), nil
	})

	r := &TeoAliasDomain110904Resource{}
	setAliasDomain110904Client(r, client)
	s := aliasDomain110904Schema(t)

	planRaw := aliasDomain110904RawObject(s, map[string]tftypes.Value{
		"zone_id":     aliasDomain110904Str("zone-3fkff38fyw8s"),
		"alias_name":  aliasDomain110904Str("alias.demo.com"),
		"target_name": aliasDomain110904Str("target.demo.com"),
		"status":      aliasDomain110904NullStr(),
	})

	req := resource.CreateRequest{
		Config: tfsdk.Config{Schema: s, Raw: planRaw},
		Plan:  tfsdk.Plan{Schema: s, Raw: planRaw},
	}
	resp := &resource.CreateResponse{
		State: tftypesDataAliasDomain110904(s, aliasDomain110904RawObject(s, nil)),
	}

	r.Create(context.Background(), req, resp)
	assert.False(t, resp.Diagnostics.HasError(), "expected no error diagnostics, got: %v", resp.Diagnostics)

	assert.NotNil(t, gotCreate)
	assert.Nil(t, gotCreate.CertType)
	assert.Nil(t, gotCreate.CertId)

	var state TeoAliasDomain110904ResourceModel
	diags := resp.State.Get(context.Background(), &state)
	assert.False(t, diags.HasError(), "State.Get returned diagnostics: %v", diags)
	assert.True(t, state.CertType.IsNull(), "cert_type should stay null when not configured")
	assert.True(t, state.CertId.IsNull(), "cert_id should stay null when not configured")
}

// TestTeoAliasDomain110904Resource_Create_APIError covers the create error
// branch: the API failure is surfaced as an error diagnostic.
func TestTeoAliasDomain110904Resource_Create_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	teoClient := mockTeoClientForAliasDomain110904(patches, client)
	patches.ApplyMethodFunc(teoClient, "CreateAliasDomainWithContext", func(_ context.Context, request *teo.CreateAliasDomainRequest) (*teo.CreateAliasDomainResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=FailedOperation, Message=domain already exists")
	})

	r := &TeoAliasDomain110904Resource{}
	setAliasDomain110904Client(r, client)
	s := aliasDomain110904Schema(t)

	planRaw := aliasDomain110904RawObject(s, map[string]tftypes.Value{
		"zone_id":     aliasDomain110904Str("zone-3fkff38fyw8s"),
		"alias_name":  aliasDomain110904Str("alias.demo.com"),
		"target_name": aliasDomain110904Str("target.demo.com"),
		"status":      aliasDomain110904NullStr(),
	})

	req := resource.CreateRequest{
		Config: tfsdk.Config{Schema: s, Raw: planRaw},
		Plan:  tfsdk.Plan{Schema: s, Raw: planRaw},
	}
	resp := &resource.CreateResponse{
		State: tftypesDataAliasDomain110904(s, aliasDomain110904RawObject(s, nil)),
	}

	r.Create(context.Background(), req, resp)
	assert.True(t, resp.Diagnostics.HasError(), "expected error diagnostics for a failed create")
}

// TestTeoAliasDomain110904Resource_Read_Success covers refresh: state is
// populated from the described alias domain and cert fields keep state values.
func TestTeoAliasDomain110904Resource_Read_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	teoClient := mockTeoClientForAliasDomain110904(patches, client)

	var gotDescribe *teo.DescribeAliasDomainsRequest
	patches.ApplyMethodFunc(teoClient, "DescribeAliasDomainsWithContext", func(_ context.Context, request *teo.DescribeAliasDomainsRequest) (*teo.DescribeAliasDomainsResponse, error) {
		gotDescribe = request
		return describeAliasDomain110904Response(sampleAliasDomain110904Active()), nil
	})

	r := &TeoAliasDomain110904Resource{}
	setAliasDomain110904Client(r, client)
	s := aliasDomain110904Schema(t)
	compositeId := "zone-3fkff38fyw8s" + tccommon.FILED_SP + "alias.demo.com"

	stateRaw := aliasDomain110904RawObject(s, map[string]tftypes.Value{
		"id":          aliasDomain110904Str(compositeId),
		"zone_id":     aliasDomain110904Str("zone-3fkff38fyw8s"),
		"alias_name":  aliasDomain110904Str("alias.demo.com"),
		"target_name": aliasDomain110904Str("old-target.demo.com"),
		"cert_type":   aliasDomain110904Str("hosting"),
		"cert_id":     aliasDomain110904StrSet([]string{"cert-xxxx"}),
		"status":      aliasDomain110904Str("pending"),
	})

	req := resource.ReadRequest{
		State: tftypesDataAliasDomain110904(s, stateRaw),
	}
	resp := &resource.ReadResponse{
		State: tftypesDataAliasDomain110904(s, stateRaw),
	}

	r.Read(context.Background(), req, resp)
	assert.False(t, resp.Diagnostics.HasError(), "expected no error diagnostics, got: %v", resp.Diagnostics)

	// Read requests an exact alias-name filter over the zone.
	assert.NotNil(t, gotDescribe)
	assert.Equal(t, "zone-3fkff38fyw8s", *gotDescribe.ZoneId)
	if assert.Len(t, gotDescribe.Filters, 1) {
		assert.Equal(t, "alias-name", *gotDescribe.Filters[0].Name)
		assert.Equal(t, []*string{ptrStringAliasDomain110904("alias.demo.com")}, gotDescribe.Filters[0].Values)
	}

	var state TeoAliasDomain110904ResourceModel
	diags := resp.State.Get(context.Background(), &state)
	assert.False(t, diags.HasError(), "State.Get returned diagnostics: %v", diags)
	assert.Equal(t, compositeId, state.Id.ValueString())
	assert.Equal(t, "target.demo.com", state.TargetName.ValueString())
	assert.Equal(t, "active", state.Status.ValueString())
	// cert_type/cert_id are not returned by the API: state values are kept.
	assert.Equal(t, "hosting", state.CertType.ValueString())
	var certIds []fwtypes.String
	certDiags := state.CertId.ElementsAs(context.Background(), &certIds, false)
	assert.False(t, certDiags.HasError(), "CertId.ElementsAs returned diagnostics: %v", certDiags)
	assert.Equal(t, []string{"cert-xxxx"}, stringValuesAliasDomain110904(certIds))
}

// TestTeoAliasDomain110904Resource_Read_NotFound covers the confirmed-gone
// branch: an empty query result removes the resource from state.
func TestTeoAliasDomain110904Resource_Read_NotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	teoClient := mockTeoClientForAliasDomain110904(patches, client)
	patches.ApplyMethodFunc(teoClient, "DescribeAliasDomainsWithContext", func(_ context.Context, request *eo.DescribeAliasDomainsRequest) (*eo.DescribeAliasDomainsResponse, error) {
		return emptyAliasDomain110904Response(), nil
	})

	r := &TeoAliasDomain110904Resource{}
	setAliasDomain110904Client(r, client)
	s := aliasDomain110904Schema(t)
	compositeId := "zone-3fkff38fyw8s" + tccommon.FILED_SP + "alias.demo.com"

	stateRaw := aliasDomain110904RawObject(s, map[string]tftypes.Value{
		"id":          aliasDomain110904Str(compositeId),
		"zone_id":     aliasDomain110904Str("zone-3fkff38fyw8s"),
		"alias_name":  aliasDomain110904Str("alias.demo.com"),
		"target_name": aliasDomain110904Str("target.demo.com"),
	})

	req := resource.ReadRequest{
		State: tftypesDataAliasDomain110904(s, stateRaw),
	}
	resp := &resource.ReadResponse{
		State: tftypesDataAliasDomain110904(s, stateRaw),
	}

	r.Read(context.Background(), req, resp)
	assert.False(t, resp.Diagnostics.HasError(), "expected no error diagnostics, got: %v", resp.Diagnostics)
	assert.True(t, resp.State.Raw.IsNull(), "expected the resource to be removed from state")
}

// TestTeoAliasDomain110904Resource_Read_BrokenId covers a corrupted state id.
func TestTeoAliasDomain110904Resource_Read_BrokenId(t *testing.T) {
	r := &TeoAliasDomain110904Resource{}
	s := aliasDomain110904Schema(t)

	stateRaw := aliasDomain110904RawObject(s, map[string]tftypes.Value{
		"id": aliasDomain110904Str("broken-id"),
	})

	req := resource.ReadRequest{
		State: tftypesDataAliasDomain110904(s, stateRaw),
	}
	resp := &resource.ReadResponse{
		State: tftypesDataAliasDomain110904(s, stateRaw),
	}

	r.Read(context.Background(), req, resp)
	assert.True(t, resp.Diagnostics.HasError(), "expected error diagnostics for a broken id")
}

// TestTeoAliasDomain110904Resource_Update_Success covers the update branch:
// target_name is always re-sent, cert fields follow the clear contract.
func TestTeoAliasDomain110904Resource_Update_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	teoClient := mockTeoClientForAliasDomain110904(patches, client)

	var gotModify *teo.ModifyAliasDomainRequest
	patches.ApplyMethodFunc(teoClient, "ModifyAliasDomainWithContext", func(_ context.Context, request *teo.ModifyAliasDomainRequest) (*teo.ModifyAliasDomainResponse, error) {
		gotModify = request
		resp := teo.NewModifyAliasDomainResponse()
		resp.Response = &teo.ModifyAliasDomainResponseParams{
			RequestId: ptrStringAliasDomain110904("fake-request-id"),
		}
		return resp, nil
	})

	active := sampleAliasDomain110904Active()
	active.TargetName = ptrStringAliasDomain110904("new-target.demo.com")
	patches.ApplyMethodFunc(teoClient, "DescribeAliasDomainsWithContext", func(_ context.Context, request *teo.DescribeAliasDomainsRequest) (*teo.DescribeAliasDomainsResponse, error) {
		return describeAliasDomain110904Response(active), nil
	})

	r := &TeoAliasDomain110904Resource{}
	setAliasDomain110904Client(r, client)
	s := aliasDomain110904Schema(t)
	compositeId := "zone-3fkff38fyw8s" + tccommon.FILED_SP + "alias.demo.com"

	planRaw := aliasDomain110904RawObject(s, map[string]tftypes.Value{
		"id":          aliasDomain110904Str(compositeId),
		"zone_id":     aliasDomain110904Str("zone-3fkff38fyw8s"),
		"alias_name":  aliasDomain110904Str("alias.demo.com"),
		"target_name": aliasDomain110904Str("new-target.demo.com"),
		"cert_type":   aliasDomain110904Str("hosting"),
		"cert_id":     aliasDomain110904StrSet([]string{"cert-yyyy"}),
		"status":      aliasDomain110904NullStr(),
	})
	stateRaw := aliasDomain110904RawObject(s, map[string]tftypes.Value{
		"id":          aliasDomain110904Str(compositeId),
		"zone_id":     aliasDomain110904Str("zone-3fkff38fyw8s"),
		"alias_name":  aliasDomain110904Str("alias.demo.com"),
		"target_name": aliasDomain110904Str("old-target.demo.com"),
		"cert_type":   aliasDomain110904Str("hosting"),
		"cert_id":     aliasDomain110904StrSet([]string{"cert-xxxx"}),
		"status":      aliasDomain110904Str("active"),
	})

	req := resource.UpdateRequest{
		Config: tfsdk.Config{Schema: s, Raw: planRaw},
		Plan:   tfsdk.Plan{Schema: s, Raw: planRaw},
		State:  tftypesDataAliasDomain110904(s, stateRaw),
	}
	resp := &resource.UpdateResponse{
		State: tftypesDataAliasDomain110904(s, stateRaw),
	}

	r.Update(context.Background(), req, resp)
	assert.False(t, resp.Diagnostics.HasError(), "expected no error diagnostics, got: %v", resp.Diagnostics)

	// Modify request: target_name re-sent, cert_type/cert_id sent with values.
	assert.NotNil(t, gotModify)
	assert.Equal(t, "zone-3fkff38fyw8s", *gotModify.ZoneId)
	assert.Equal(t, "alias.demo.com", *gotModify.AliasName)
	assert.Equal(t, "new-target.demo.com", *gotModify.TargetName)
	assert.Equal(t, "hosting", *gotModify.CertType)
	assert.Equal(t, []*string{ptrStringAliasDomain110904("cert-yyyy")}, gotModify.CertId)

	var state TeoAliasDomain110904ResourceModel
	diags := resp.State.Get(context.Background(), &state)
	assert.False(t, diags.HasError(), "State.Get returned diagnostics: %v", diags)
	assert.Equal(t, compositeId, state.Id.ValueString())
	assert.Equal(t, "new-target.demo.com", state.TargetName.ValueString())
	assert.Equal(t, "active", state.Status.ValueString())
}

// TestTeoAliasDomain110904Resource_Update_ClearCert covers the clear contract:
// explicitly setting cert_type = "none" sends "none" and clears cert_id.
func TestTeoAliasDomain110904Resource_Update_ClearCert(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	teoClient := mockTeoClientForAliasDomain110904(patches, client)

	var gotModify *teo.ModifyAliasDomainRequest
	patches.ApplyMethodFunc(teoClient, "ModifyAliasDomainWithContext", func(_ context.Context, request *teo.ModifyAliasDomainRequest) (*teo.ModifyAliasDomainResponse, error) {
		gotModify = request
		resp := teo.NewModifyAliasDomainResponse()
		resp.Response = &teo.ModifyAliasDomainResponseParams{
			RequestId: ptrStringAliasDomain110904("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(teoClient, "DescribeAliasDomainsWithContext", func(_ context.Context, request *eo.DescribeAliasDomainsRequest) (*teo.DescribeAliasDomainsResponse, error) {
		return describeAliasDomain110904Response(sampleAliasDomain110904Active()), nil
	})

	r := &TeoAliasDomain110904Resource{}
	setAliasDomain110904Client(r, client)
	s := aliasDomain110904Schema(t)
	compositeId := "zone-3fkff38fyw8s" + tccommon.FILED_SP + "alias.demo.com"

	// plan cert_type = "none" (the explicit clear path) vs state "hosting".
	planRaw := aliasDomain110904RawObject(s, map[string]tftypes.Value{
		"id":          aliasDomain110904Str(compositeId),
		"zone_id":     aliasDomain110904Str("zone-3fkff38fyw8s"),
		"alias_name":  aliasDomain110904Str("alias.demo.com"),
		"target_name": aliasDomain110904Str("target.demo.com"),
		"cert_type":   aliasDomain110904Str("none"),
		"cert_id":     aliasDomain110904StrSet(nil),
		"status":      aliasDomain110904NullStr(),
	})
	stateRaw := aliasDomain110904RawObject(s, map[string]tftypes.Value{
		"id":          aliasDomain110904Str(compositeId),
		"zone_id":     aliasDomain110904Str("zone-3fkff38fyw8s"),
		"alias_name":  aliasDomain110904Str("alias.demo.com"),
		"target_name": aliasDomain110904Str("target.demo.com"),
		"cert_type":   aliasDomain110904Str("hosting"),
		"cert_id":     aliasDomain110904StrSet([]string{"cert-xxxx"}),
		"status":      aliasDomain110904Str("active"),
	})

	req := resource.UpdateRequest{
		Config: tfsdk.Config{Schema: s, Raw: planRaw},
		Plan:   tfsdk.Plan{Schema: s, Raw: planRaw},
		State:  tftypesDataAliasDomain110904(s, stateRaw),
	}
	resp := &resource.UpdateResponse{
		State: tftypesDataAliasDomain110904(s, stateRaw),
	}

	r.Update(context.Background(), req, resp)
	assert.False(t, resp.Diagnostics.HasError(), "expected no error diagnostics, got: %v", resp.Diagnostics)

	assert.NotNil(t, gotModify)
	assert.Equal(t, "target.demo.com", *gotModify.TargetName)
	assert.NotNil(t, gotModify.CertType)
	assert.Equal(t, "none", *gotModify.CertType)
	assert.Empty(t, gotModify.CertId)

	var state TeoAliasDomain110904ResourceModel
	diags := resp.State.Get(context.Background(), &state)
	assert.False(t, diags.HasError(), "State.Get returned diagnostics: %v", diags)
	assert.Equal(t, "none", state.CertType.ValueString())
	assert.True(t, state.CertId.IsNull(), "cert_id should be null after clearing, got %v", state.CertId)
}

// TestTeoAliasDomain110904Resource_Update_NeverConfiguredCert covers the
// "never configured" update: an unknown/null plan cert_type against a null
// state value sends no cert fields at all (the cloud keeps its own default).
func TestTeoAliasDomain110904Resource_Update_NeverConfiguredCert(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	teoClient := mockTeoClientForAliasDomain110904(patches, client)

	var gotModify *eo.ModifyAliasDomainRequest
	patches.ApplyMethodFunc(teoClient, "ModifyAliasDomainWithContext", func(_ context.Context, request *eo.ModifyAliasDomainRequest) (*eo.ModifyAliasDomainResponse, error) {
		gotModify = request
		resp := teo.NewModifyAliasDomainResponse()
		resp.Response = &teo.ModifyAliasDomainResponseParams{
			RequestId: ptrStringAliasDomain110904("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(teoClient, "DescribeAliasDomainsWithContext", func(_ context.Context, request *eo.DescribeAliasDomainsRequest) (*teo.DescribeAliasDomainsResponse, error) {
		return describeAliasDomain110904Response(sampleAliasDomain110904Active()), nil
	})

	r := &TeoAliasDomain110904Resource{}
	setAliasDomain110904Client(r, client)
	s := aliasDomain110904Schema(t)
	compositeId := "zone-3fkff38fyw8s" + tccommon.FILED_SP + "alias.demo.com"

	planRaw := aliasDomain110904RawObject(s, map[string]tftypes.Value{
		"id":          aliasDomain110904Str(compositeId),
		"zone_id":     aliasDomain110904Str("zone-3fkff38fyw8s"),
		"alias_name":  aliasDomain110904Str("alias.demo.com"),
		"target_name": aliasDomain110904Str("new-target.demo.com"),
		"cert_type":   aliasDomain110904UnknownStr(),
		"cert_id":     aliasDomain110904UnknownStrSet(),
		"status":      aliasDomain110904NullStr(),
	})
	stateRaw := aliasDomain110904RawObject(s, map[string]tftypes.Value{
		"id":          aliasDomain110904Str(compositeId),
		"zone_id":     aliasDomain110904Str("zone-3fkff38fyw8s"),
		"alias_name":  aliasDomain110904Str("alias.demo.com"),
		"target_name": aliasDomain110904Str("old-target.demo.com"),
		"cert_type":   aliasDomain110904NullStr(),
		"cert_id":     aliasDomain110904StrSet(nil),
		"status":      aliasDomain110904Str("active"),
	})

	req := resource.UpdateRequest{
		Config: tfsdk.Config{Schema: s, Raw: planRaw},
		Plan:   tfsdk.Plan{Schema: s, Raw: planRaw},
		State:  tftypesDataAliasDomain110904(s, stateRaw),
	}
	resp := &resource.UpdateResponse{
		State: tftypesDataAliasDomain110904(s, stateRaw),
	}

	r.Update(context.Background(), req, resp)
	assert.False(t, resp.Diagnostics.HasError(), "expected no error diagnostics, got: %v", resp.Diagnostics)

	assert.NotNil(t, gotModify)
	assert.Equal(t, "new-target.demo.com", *gotModify.TargetName)
	assert.Nil(t, gotModify.CertType, "cert_type should not be sent when never configured")
	assert.Nil(t, gotModify.CertId, "cert_id should not be sent when never configured")
}

// TestTeoAliasDomain110904Resource_Update_UnchangedCert covers the unchanged
// case: when the certificate configuration matches the state, neither
// CertType nor CertId is sent (the cloud keeps the previous configuration).
func TestTeoAliasDomain110904Resource_Update_UnchangedCert(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	teoClient := mockTeoClientForAliasDomain110904(patches, client)

	var gotModify *teo.ModifyAliasDomainRequest
	patches.ApplyMethodFunc(teoClient, "ModifyAliasDomainWithContext", func(_ context.Context, request *eo.ModifyAliasDomainRequest) (*eo.ModifyAliasDomainResponse, error) {
		gotModify = request
		resp := teo.NewModifyAliasDomainResponse()
		resp.Response = &eo.ModifyAliasDomainResponseParams{
			RequestId: ptrStringAliasDomain110904("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(teoClient, "DescribeAliasDomainsWithContext", func(_ context.Context, request *eo.DescribeAliasDomainsRequest) (*eo.DescribeAliasDomainsResponse, error) {
		return describeAliasDomain110904Response(sampleAliasDomain110904Active()), nil
	})

	r := &TeoAliasDomain110904Resource{}
	setAliasDomain110904Client(r, client)
	s := aliasDomain110904Schema(t)
	compositeId := "zone-3fkff38fyw8s" + tccommon.FILED_SP + "alias.demo.com"

	// plan and state carry the same certificate configuration; only the
	// target domain changes.
	planRaw := aliasDomain110904RawObject(s, map[string]tftypes.Value{
		"id":          aliasDomain110904Str(compositeId),
		"zone_id":     aliasDomain110904Str("zone-3fkff38fyw8s"),
		"alias_name":  aliasDomain110904Str("alias.demo.com"),
		"target_name": aliasDomain110904Str("new-target.demo.com"),
		"cert_type":   aliasDomain110904Str("hosting"),
		"cert_id":     aliasDomain110904StrSet([]string{"cert-xxxx"}),
		"status":      aliasDomain110904NullStr(),
	})
	stateRaw := aliasDomain110904RawObject(s, map[string]tftypes.Value{
		"id":          aliasDomain110904Str(compositeId),
		"zone_id":     aliasDomain110904Str("zone-3fkff38fyw8s"),
		"alias_name":  aliasDomain110904Str("alias.demo.com"),
		"target_name": aliasDomain110904Str("old-target.demo.com"),
		"cert_type":   aliasDomain110904Str("hosting"),
		"cert_id":     aliasDomain110904StrSet([]string{"cert-xxxx"}),
		"status":      aliasDomain110904Str("active"),
	})

	req := resource.UpdateRequest{
		Config: tfsdk.Config{Schema: s, Raw: planRaw},
		Plan:   tfsdk.Plan{Schema: s, Raw: planRaw},
		State:  tftypesDataAliasDomain110904(s, stateRaw),
	}
	resp := &resource.UpdateResponse{
		State: tftypesDataAliasDomain110904(s, stateRaw),
	}

	r.Update(context.Background(), req, resp)
	assert.False(t, resp.Diagnostics.HasError(), "expected no error diagnostics, got: %v", resp.Diagnostics)

	assert.NotNil(t, gotModify)
	assert.Equal(t, "new-target.demo.com", *gotModify.TargetName)
	assert.Nil(t, gotModify.CertType, "cert_type should not be sent when unchanged")
	assert.Nil(t, gotModify.CertId, "cert_id should not be sent when unchanged")
}

// TestTeoAliasDomain110904Resource_Update_APIError covers the update error
// branch.
func TestTeoAliasDomain110904Resource_Update_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	teoClient := mockTeoClientForAliasDomain110904(patches, client)
	patches.ApplyMethodFunc(teoClient, "ModifyAliasDomainWithContext", func(_ context.Context, request *eo.ModifyAliasDomainRequest) (*eo.ModifyAliasDomainResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=ResourceNotFound, Message=alias domain not found")
	})

	r := &TeoAliasDomain110904Resource{}
	setAliasDomain110904Client(r, client)
	s := aliasDomain110904Schema(t)
	compositeId := "zone-3fkff38fyw8s" + tccommon.FILED_SP + "alias.demo.com"

	planRaw := aliasDomain110904RawObject(s, map[string]tftypes.Value{
		"id":          aliasDomain110904Str(compositeId),
		"zone_id":     aliasDomain110904Str("zone-3fkff38fyw8s"),
		"alias_name":  aliasDomain110904Str("alias.demo.com"),
		"target_name": aliasDomain110904Str("new-target.demo.com"),
		"status":      aliasDomain110904NullStr(),
	})
	stateRaw := aliasDomain110904RawObject(s, map[string]tftypes.Value{
		"id":          aliasDomain110904Str(compositeId),
		"zone_id":     aliasDomain110904Str("zone-3fkff38fyw8s"),
		"alias_name":  aliasDomain110904Str("alias.demo.com"),
		"target_name": aliasDomain110904Str("target.demo.com"),
		"status":      aliasDomain110904Str("active"),
	})

	req := resource.UpdateRequest{
		Config: tfsdk.Config{Schema: s, Raw: planRaw},
		Plan:   tfsdk.Plan{Schema: s, Raw: planRaw},
		State:  tftypesDataAliasDomain110904(s, stateRaw),
	}
	resp := &resource.UpdateResponse{
		State: tftypesDataAliasDomain110904(s, stateRaw),
	}

	r.Update(context.Background(), req, resp)
	assert.True(t, resp.Diagnostics.HasError(), "expected error diagnostics for a failed update")
}

// TestTeoAliasDomain110904Resource_Delete_Success covers a successful delete:
// the alias domain is deleted with a single-element AliasNames list and the
// delete wait converges.
func TestTeoAliasDomain110904Resource_Delete_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	teoClient := mockTeoClientForAliasDomain110904(patches, client)

	// The first describe (pre-delete existence check) sees the alias domain;
	// every later describe (the delete wait) sees it gone.
	describeCalls := 0
	patches.ApplyMethodFunc(teoClient, "DescribeAliasDomainsWithContext", func(_ context.Context, request *eo.DescribeAliasDomainsRequest) (*eo.DescribeAliasDomainsResponse, error) {
		describeCalls++
		if describeCalls == 1 {
			return describeAliasDomain110904Response(sampleAliasDomain110904Active()), nil
		}
		return emptyAliasDomain110904Response(), nil
	})

	var gotDelete *eo.DeleteAliasDomainRequest
	patches.ApplyMethodFunc(teoClient, "DeleteAliasDomainWithContext", func(_ context.Context, request *eo.DeleteAliasDomainRequest) (*eo.DeleteAliasDomainResponse, error) {
		gotDelete = request
		resp := teo.NewDeleteAliasDomainResponse()
		resp.Response = &eo.DeleteAliasDomainResponseParams{
			RequestId: ptrStringAliasDomain110904("fake-request-id"),
		}
		return resp, nil
	})

	r := &TeoAliasDomain110904Resource{}
	setAliasDomain110904Client(r, client)
	s := aliasDomain110904Schema(t)
	compositeId := "zone-3fkff38fyw8s" + tccommon.FILED_SP + "alias.demo.com"

	stateRaw := aliasDomain110904RawObject(s, map[string]tftypes.Value{
		"id":          aliasDomain110904Str(compositeId),
		"zone_id":     aliasDomain110904Str("zone-3fkff38fyw8s"),
		"alias_name":  aliasDomain110904Str("alias.demo.com"),
		"target_name": aliasDomain110904Str("target.demo.com"),
		"status":      aliasDomain110904Str("active"),
	})

	req := resource.DeleteRequest{
		State: tftypesDataAliasDomain110904(s, stateRaw),
	}
	resp := &resource.DeleteResponse{
		State: tftypesDataAliasDomain110904(s, stateRaw),
	}

	r.Delete(context.Background(), req, resp)
	assert.False(t, resp.Diagnostics.HasError(), "expected no error diagnostics, got: %v", resp.Diagnostics)

	assert.NotNil(t, gotDelete)
	assert.Equal(t, "zone-3fkff38fyw8s", *gotDelete.ZoneId)
	assert.Equal(t, []*string{ptrStringAliasDomain110904("alias.demo.com")}, gotDelete.AliasNames)
}

// TestTeoAliasDomain110904Resource_Delete_AlreadyGone covers the confirmed
// already-gone branch: an empty query result makes Delete a no-op.
func TestTeoAliasDomain110904Resource_Delete_AlreadyGone(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	teoClient := mockTeoClientForAliasDomain110904(patches, client)

	deleteCalled := false
	patches.ApplyMethodFunc(teoClient, "DescribeAliasDomainsWithContext", func(_ context.Context, request *eo.DescribeAliasDomainsRequest) (*eo.DescribeAliasDomainsResponse, error) {
		return emptyAliasDomain110904Response(), nil
	})
	patches.ApplyMethodFunc(teoClient, "DeleteAliasDomainWithContext", func(_ context.Context, request *eo.DeleteAliasDomainRequest) (*eo.DeleteAliasDomainResponse, error) {
		deleteCalled = true
		resp := teo.NewDeleteAliasDomainResponse()
		resp.Response = &eo.DeleteAliasDomainResponseParams{
			RequestId: ptrStringAliasDomain110904("fake-request-id"),
		}
		return resp, nil
	})

	r := &TeoAliasDomain110904Resource{}
	setAliasDomain110904Client(r, client)
	s := aliasDomain110904Schema(t)
	compositeId := "zone-3fkff38fyw8s" + tccommon.FILED_SP + "alias.demo.com"

	stateRaw := aliasDomain110904RawObject(s, map[string]tftypes.Value{
		"id":          aliasDomain110904Str(compositeId),
		"zone_id":     aliasDomain110904Str("zone-3fkff38fyw8s"),
		"alias_name":  aliasDomain110904Str("alias.demo.com"),
		"target_name": aliasDomain110904Str("target.demo.com"),
		"status":      aliasDomain110904Str("active"),
	})

	req := resource.DeleteRequest{
		State: tftypesDataAliasDomain110904(s, stateRaw),
	}
	resp := &resource.DeleteResponse{
		State: tftypesDataAliasDomain110904(s, stateRaw),
	}

	r.Delete(context.Background(), req, resp)
	assert.False(t, resp.Diagnostics.HasError(), "expected no error diagnostics, got: %v", resp.Diagnostics)
	assert.False(t, deleteCalled, "DeleteAliasDomain should not be called when the alias domain is already gone")
}

// TestTeoAliasDomain110904Resource_Delete_APIError covers the delete error
// branch.
func TestTeoAliasDomain110904Resource_Delete_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	teoClient := mockTeoClientForAliasDomain110904(patches, client)

	patches.ApplyMethodFunc(teoClient, "DescribeAliasDomainsWithContext", func(_ context.Context, request *eo.DescribeAliasDomainsRequest) (*eo.DescribeAliasDomainsResponse, error) {
		return describeAliasDomain110904Response(sampleAliasDomain110904Active()), nil
	})
	patches.ApplyMethodFunc(teoClient, "DeleteAliasDomainWithContext", func(_ context.Context, request *eo.DeleteAliasDomainRequest) (*eo.DeleteAliasDomainResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=FailedOperation, Message=delete failed")
	})

	r := &TeoAliasDomain110904Resource{}
	setAliasDomain110904Client(r, client)
	s := aliasDomain110904Schema(t)
	compositeId := "zone-3fkff38fyw8s" + tccommon.FILED_SP + "alias.demo.com"

	stateRaw := aliasDomain110904RawObject(s, map[string]tftypes.Value{
		"id":          aliasDomain110904Str(compositeId),
		"zone_id":     aliasDomain110904Str("zone-3fkff38fyw8s"),
		"alias_name":  aliasDomain110904Str("alias.demo.com"),
		"target_name": aliasDomain110904Str("target.demo.com"),
		"status":      aliasDomain110904Str("active"),
	})

	req := resource.DeleteRequest{
		State: tftypesDataAliasDomain110904(s, stateRaw),
	}
	resp := &resource.DeleteResponse{
		State: tftypesDataAliasDomain110904(s, stateRaw),
	}

	r.Delete(context.Background(), req, resp)
	assert.True(t, resp.Diagnostics.HasError(), "expected error diagnostics for a failed delete")
}

// TestSplitTeoAliasDomain110904Id validates parsing of `<zone_id>#<alias_name>`.
func TestSplitTeoAliasDomain110904Id(t *testing.T) {
	good := "zone-3fkff38fyw8s" + tccommon.FILED_SP + "alias.demo.com"
	zone, alias, ok := splitTeoAliasDomain110904Id(good)
	if !ok || zone != "zone-3fkff38fyw8s" || alias != "alias.demo.com" {
		t.Errorf("splitTeoAliasDomain110904Id(%q) = (%q, %q, %v)", good, zone, alias, ok)
	}

	bad := []string{
		"",
		"only-zone",
		tccommon.FILED_SP + "alias.demo.com",
		"zone-3fkff38fyw8s" + tccommon.FILED_SP,
		"a" + tccommon.FILED_SP + "b" + tccommon.FILED_SP + "c",
	}
	for _, in := range bad {
		if _, _, ok := splitTeoAliasDomain110904Id(in); ok {
			t.Errorf("splitTeoAliasDomain110904Id(%q) ok=true, want false", in)
		}
	}
}

// TestTeoAliasDomain110904CertParams validates the cert_type/cert_id request
// conversion, including the null and empty cases.
func TestTeoAliasDomain110904CertParams(t *testing.T) {
	ctx := context.Background()

	// Null cert_type: nothing to send.
	typ, ids, diags := teoAliasDomain110904CertParams(ctx, fwtypes.StringNull(), fwtypes.SetNull(fwtypes.StringType))
	assert.False(t, diags.HasError())
	assert.Empty(t, typ)
	assert.Nil(t, ids)

	// Known cert_type without cert_id: type only.
	typ, ids, diags = teoAliasDomain110904CertParams(ctx, fwtypes.StringValue("none"), fwtypes.SetNull(fwtypes.StringType))
	assert.False(t, diags.HasError())
	assert.Equal(t, "none", typ)
	assert.Nil(t, ids)

	// Known cert_type with cert_id set: both are sent.
	certSet, setDiags := fwtypes.SetValueFrom(ctx, fwtypes.StringType, []string{"cert-xxxx", "cert-yyyy"})
	assert.False(t, setDiags.HasError())
	typ, ids, diags = teoAliasDomain110904CertParams(ctx, fwtypes.StringValue("hosting"), certSet)
	assert.False(t, diags.HasError())
	assert.Equal(t, "hosting", typ)
	assert.Equal(t, []*string{ptrStringAliasDomain110904("cert-xxxx"), ptrStringAliasDomain110904("cert-yyyy")}, ids)
}

// stringValuesAliasDomain110904 converts []types.String into []string.
func stringValuesAliasDomain110904(values []fwtypes.String) []string {
	out := make([]string, 0, len(values))
	for _, v := range values {
		out = append(out, v.ValueString())
	}
	return out
}
