package teo_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/assert"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/sharedmeta"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/teo"
)

// go test ./tencentcloud/services/teo/ -run "TestTeoZone110901Resource" -v -count=1 -gcflags="all=-l"

func ptrStringZone1109(s string) *string {
	return &s
}

func ptrBoolZone1109(b bool) *bool {
	return &b
}

func ptrInt64Zone1109(i int64) *int64 {
	return &i
}

// teoZone110901ResourceSchema builds the resource schema for tests.
func teoZone110901ResourceSchema(t *testing.T) schema.Schema {
	r := teo.NewTeoZone110901Resource()
	resp := &resource.SchemaResponse{}
	r.Schema(context.Background(), resource.SchemaRequest{}, resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("failed to build resource schema: %v", resp.Diagnostics)
	}
	return resp.Schema
}

// teoZone110901AttrTypes builds the tftypes attribute map of the schema.
func teoZone110901AttrTypes(t *testing.T) map[string]tftypes.Type {
	s := teoZone110901ResourceSchema(t)
	out := map[string]tftypes.Type{}
	for name, attr := range s.GetAttributes() {
		out[name] = attr.GetType().TerraformType(context.Background())
	}
	return out
}

// teoZone110901StringSet builds a tftypes set of strings.
func teoZone110901StringSet(elems []string) tftypes.Value {
	vals := make([]tftypes.Value, 0, len(elems))
	for _, e := range elems {
		vals = append(vals, tftypes.NewValue(tftypes.String, e))
	}
	return tftypes.NewValue(tftypes.Set{ElementType: tftypes.String}, vals)
}

// teoZone110901TagObject builds a tftypes tag object.
func teoZone110901TagObject(key, value string) tftypes.Value {
	return tftypes.NewValue(tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"tag_key":   tftypes.String,
			"tag_value": tftypes.String,
		},
	}, map[string]tftypes.Value{
		"tag_key":   tftypes.NewValue(tftypes.String, key),
		"tag_value": tftypes.NewValue(tftypes.String, value),
	})
}

// teoZone110901VanityObject builds a tftypes vanity_name_servers object.
func teoZone110901VanityObject(switchVal string, servers []string) tftypes.Value {
	return tftypes.NewValue(tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"switch":  tftypes.String,
			"servers": tftypes.Set{ElementType: tftypes.String},
		},
	}, map[string]tftypes.Value{
		"switch":  tftypes.NewValue(tftypes.String, switchVal),
		"servers": teoZone110901StringSet(servers),
	})
}

// teoZone110901FullValues builds the full raw object value used as plan.
func teoZone110901FullValues(t *testing.T) map[string]tftypes.Value {
	attrTypes := teoZone110901AttrTypes(t)
	vals := map[string]tftypes.Value{
		"id":               tftypes.NewValue(tftypes.String, nil),
		"zone_id":          tftypes.NewValue(tftypes.String, nil),
		"type":             tftypes.NewValue(tftypes.String, "partial"),
		"zone_name":        tftypes.NewValue(tftypes.String, "example.com"),
		"area":             tftypes.NewValue(tftypes.String, "overseas"),
		"plan_id":          tftypes.NewValue(tftypes.String, "edgeone-2kfv1h391n6w"),
		"alias_zone_name":  tftypes.NewValue(tftypes.String, "my-alias"),
		"allow_duplicates": tftypes.NewValue(tftypes.Bool, true),
		"jump_start":       tftypes.NewValue(tftypes.Bool, false),
		"tags": tftypes.NewValue(tftypes.Set{
			ElementType: tftypes.Object{
				AttributeTypes: map[string]tftypes.Type{
					"tag_key":   tftypes.String,
					"tag_value": tftypes.String,
				},
			},
		}, []tftypes.Value{teoZone110901TagObject("createdBy", "terraform")}),
	}
	// fill the remaining computed attributes with nulls
	for name, typ := range attrTypes {
		if _, ok := vals[name]; !ok {
			vals[name] = tftypes.NewValue(typ, nil)
		}
	}
	return vals
}

// newTeoZone110901CreateRequest builds a CreateRequest with a full plan.
func newTeoZone110901CreateRequest(t *testing.T) resource.CreateRequest {
	s := teoZone110901ResourceSchema(t)
	raw := tftypes.NewValue(tftypes.Object{AttributeTypes: teoZone110901AttrTypes(t)}, teoZone110901FullValues(t))
	req := resource.CreateRequest{}
	req.Plan.Raw = raw
	req.Plan.Schema = s
	req.Config.Raw = raw
	req.Config.Schema = s
	return req
}

// setTeoZone110901Client injects a mock connectivity client into the resource
// via the Configure method (same path as the framework provider).
func setTeoZone110901Client(r *teo.TeoZone110901Resource, client *connectivity.TencentCloudClient) {
	meta := &sharedmeta.ProviderMeta{Client: client}
	configureReq := resource.ConfigureRequest{ProviderData: meta}
	configureResp := &resource.ConfigureResponse{}
	r.Configure(context.Background(), configureReq, configureResp)
}

// mockTeoZone110901Zone builds a full zone object as returned by DescribeZones.
func mockTeoZone110901Zone(zoneId string) *teov20220901.Zone {
	return &teov20220901.Zone{
		ZoneId:        ptrStringZone1109(zoneId),
		ZoneName:      ptrStringZone1109("example.com"),
		AliasZoneName: ptrStringZone1109("my-alias"),
		Area:          ptrStringZone1109("overseas"),
		Type:          ptrStringZone1109("partial"),
		Tags: []*teov20220901.Tag{
			{TagKey: ptrStringZone1109("createdBy"), TagValue: ptrStringZone1109("terraform")},
		},
		Resources: []*teov20220901.Resource{
			{
				Id:            ptrStringZone1109("edgeone-2kfv1h391n6w"),
				PayMode:       ptrInt64Zone1109(0),
				PlanId:        ptrStringZone1109("edgeone-2kfv1h391n6w"),
				Status:        ptrStringZone1109("normal"),
				Area:          ptrStringZone1109("overseas"),
				Group:         ptrStringZone1109("plan"),
				CreateTime:    ptrStringZone1109("2025-01-01T00:00:00+08:00"),
				EnableTime:    ptrStringZone1109("2025-01-01T00:00:00+08:00"),
				ExpireTime:    ptrStringZone1109("2026-01-01T00:00:00+08:00"),
				AutoRenewFlag: ptrInt64Zone1109(1),
				ZoneNumber:    ptrInt64Zone1109(1),
			},
		},
		Status:              ptrStringZone1109("pending"),
		CnameStatus:         ptrStringZone1109("pending"),
		ActiveStatus:        ptrStringZone1109("inactive"),
		LockStatus:          ptrStringZone1109("enable"),
		Paused:              ptrBoolZone1109(false),
		CreatedOn:           ptrStringZone1109("2025-01-01T00:00:00+08:00"),
		ModifiedOn:          ptrStringZone1109("2025-01-01T00:00:00+08:00"),
		NameServers:         []*string{ptrStringZone1109("ns1.teo.example.com")},
		OriginalNameServers: []*string{ptrStringZone1109("ns1.old.example.com")},
		OwnershipVerification: &teov20220901.OwnershipVerification{
			DnsVerification: &teov20220901.DnsVerification{
				Subdomain:   ptrStringZone1109("_teo.example.com"),
				RecordType:  ptrStringZone1109("TXT"),
				RecordValue: ptrStringZone1109("teo-site-verification=abc123"),
			},
			FileVerification: &teov20220901.FileVerification{
				Path:    ptrStringZone1109("/.well-known/teo-verification/z12h416twn.txt"),
				Content: ptrStringZone1109("content-123"),
			},
			NsVerification: &teov20220901.NsVerification{
				NameServers: []*string{ptrStringZone1109("ns1.teo.example.com"), ptrStringZone1109("ns2.teo.example.com")},
			},
		},
		VanityNameServers: &teov20220901.VanityNameServers{
			Switch:  ptrStringZone1109("on"),
			Servers: []*string{ptrStringZone1109("ns1.example.com")},
		},
	}
}

// TestTeoZone110901Resource_Metadata validates the resource type name.
func TestTeoZone110901Resource_Metadata(t *testing.T) {
	r := teo.NewTeoZone110901Resource()
	req := resource.MetadataRequest{}
	resp := &resource.MetadataResponse{}

	r.Metadata(context.Background(), req, resp)

	if resp.TypeName != "tencentcloud_teo_zone_1109_01" {
		t.Errorf("Metadata.TypeName = %q, want %q", resp.TypeName, "tencentcloud_teo_zone_1109_01")
	}
}

// TestTeoZone110901Resource_Schema validates required / optional / computed
// flags on every attribute.
func TestTeoZone110901Resource_Schema(t *testing.T) {
	r := teo.NewTeoZone110901Resource()
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
		"id":                     {computed: true},
		"zone_id":                {computed: true},
		"type":                   {optional: true, computed: true},
		"zone_name":              {optional: true, computed: true},
		"area":                   {optional: true, computed: true},
		"plan_id":                {optional: true, computed: true},
		"alias_zone_name":        {optional: true, computed: true},
		"tags":                   {optional: true, computed: true},
		"allow_duplicates":       {optional: true, computed: true},
		"jump_start":             {optional: true, computed: true},
		"vanity_name_servers":    {optional: true, computed: true},
		"status":                 {computed: true},
		"cname_status":           {computed: true},
		"active_status":          {computed: true},
		"lock_status":            {computed: true},
		"paused":                 {computed: true},
		"created_on":             {computed: true},
		"modified_on":            {computed: true},
		"name_servers":           {computed: true},
		"original_name_servers":  {computed: true},
		"ownership_verification": {computed: true},
	}

	for name, want := range cases {
		attr, ok := resp.Schema.GetAttributes()[name]
		if !ok {
			t.Errorf("attribute %q missing from schema", name)
			continue
		}
		var gotRequired, gotOptional, gotComputed bool
		switch a := attr.(type) {
		case schema.StringAttribute:
			gotRequired, gotOptional, gotComputed = a.Required, a.Optional, a.Computed
		case schema.BoolAttribute:
			gotRequired, gotOptional, gotComputed = a.Required, a.Optional, a.Computed
		case schema.SetNestedAttribute:
			gotRequired, gotOptional, gotComputed = a.Required, a.Optional, a.Computed
		case schema.SetAttribute:
			gotRequired, gotOptional, gotComputed = a.Required, a.Optional, a.Computed
		case schema.SingleNestedAttribute:
			gotRequired, gotOptional, gotComputed = a.Required, a.Optional, a.Computed
		default:
			t.Errorf("attribute %q has unexpected type %T", name, attr)
			continue
		}
		if gotRequired != want.required {
			t.Errorf("attribute %q Required=%v, want %v", name, gotRequired, want.required)
		}
		if gotOptional != want.optional {
			t.Errorf("attribute %q Optional=%v, want %v", name, gotOptional, want.optional)
		}
		if gotComputed != want.computed {
			t.Errorf("attribute %q Computed=%v, want %v", name, gotComputed, want.computed)
		}
	}

	// nested attribute checks
	tags, ok := resp.Schema.GetAttributes()["tags"].(schema.SetNestedAttribute)
	if !ok {
		t.Fatalf("tags is not a SetNestedAttribute")
	}
	if tagKey, ok := tags.GetNestedObject().GetAttributes()["tag_key"].(schema.StringAttribute); !ok {
		t.Errorf("tags.tag_key missing or not a string attribute")
	} else if !tagKey.Required {
		t.Errorf("tags.tag_key must be Required")
	}
	if tagValue, ok := tags.GetNestedObject().GetAttributes()["tag_value"].(schema.StringAttribute); !ok {
		t.Errorf("tags.tag_value missing or not a string attribute")
	} else if !tagValue.Required {
		t.Errorf("tags.tag_value must be Required")
	}

	vns, ok := resp.Schema.GetAttributes()["vanity_name_servers"].(schema.SingleNestedAttribute)
	if !ok {
		t.Fatalf("vanity_name_servers is not a SingleNestedAttribute")
	}
	if sw, ok := vns.GetAttributes()["switch"].(schema.StringAttribute); !ok {
		t.Errorf("vanity_name_servers.switch missing or not a string attribute")
	} else if !sw.Required {
		t.Errorf("vanity_name_servers.switch must be Required")
	}
}

// TestTeoZone110901Resource_ImplementsInterfaces verifies the resource
// satisfies the framework interfaces.
func TestTeoZone110901Resource_ImplementsInterfaces(t *testing.T) {
	var _ resource.ResourceWithConfigure = &teo.TeoZone110901Resource{}
	var _ resource.ResourceWithImportState = &teo.TeoZone110901Resource{}
}

// TestTeoZone110901Resource_Factory validates the registry factory.
func TestTeoZone110901Resource_Factory(t *testing.T) {
	got := teo.NewTeoZone110901Resource()
	if got == nil {
		t.Fatal("NewTeoZone110901Resource returned nil")
	}
	if _, ok := got.(*teo.TeoZone110901Resource); !ok {
		t.Fatalf("factory returned %T, want *TeoZone110901Resource", got)
	}
}

// TestTeoZone110901Resource_ImportState validates the embedded passthrough
// import behavior (importing the zone id writes it to the id attribute).
func TestTeoZone110901Resource_ImportState(t *testing.T) {
	r := &teo.TeoZone110901Resource{}
	s := teoZone110901ResourceSchema(t)

	req := resource.ImportStateRequest{ID: "zone-2noz78a8ev6k"}
	resp := &resource.ImportStateResponse{
		State: tfsdk.State{
			Raw:    tftypes.NewValue(tftypes.Object{AttributeTypes: teoZone110901AttrTypes(t)}, nil),
			Schema: s,
		},
	}
	r.ImportState(context.Background(), req, resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("ImportState returned diagnostics: %v", resp.Diagnostics)
	}

	var model teo.TeoZone110901ResourceModel
	if diags := resp.State.Get(context.Background(), &model); diags.HasError() {
		t.Fatalf("failed to read imported state: %v", diags)
	}
	if model.Id.ValueString() != "zone-2noz78a8ev6k" {
		t.Errorf("imported id = %q, want zone-2noz78a8ev6k", model.Id.ValueString())
	}
}

// TestTeoZone110901Resource_Create_Success covers a successful create with a
// full plan and read-back, asserting every leaf field written to state.
func TestTeoZone110901Resource_Create_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(client, "UseTeoClient", teoClient)

	var gotCreate *teov20220901.CreateZoneRequest
	patches.ApplyMethodFunc(teoClient, "CreateZone", func(request *teov20220901.CreateZoneRequest) (*teov20220901.CreateZoneResponse, error) {
		gotCreate = request
		resp := teov20220901.NewCreateZoneResponse()
		resp.Response = &teov20220901.CreateZoneResponseParams{
			ZoneId: ptrStringZone1109("zone-2noz78a8ev6k"),
			OwnershipVerification: &teov20220901.OwnershipVerification{
				DnsVerification: &teov20220901.DnsVerification{
					Subdomain:   ptrStringZone1109("_teo.example.com"),
					RecordType:  ptrStringZone1109("TXT"),
					RecordValue: ptrStringZone1109("teo-site-verification=abc123"),
				},
			},
			RequestId: ptrStringZone1109("fake-request-id"),
		}
		return resp, nil
	})

	var gotDescribe *teov20220901.DescribeZonesRequest
	patches.ApplyMethodFunc(teoClient, "DescribeZones", func(request *teov20220901.DescribeZonesRequest) (*teov20220901.DescribeZonesResponse, error) {
		gotDescribe = request
		resp := teov20220901.NewDescribeZonesResponse()
		resp.Response = &teov20220901.DescribeZonesResponseParams{
			TotalCount: ptrInt64Zone1109(1),
			Zones:      []*teov20220901.Zone{mockTeoZone110901Zone("zone-2noz78a8ev6k")},
			RequestId:  ptrStringZone1109("fake-request-id"),
		}
		return resp, nil
	})

	r := &teo.TeoZone110901Resource{}
	setTeoZone110901Client(r, client)

	req := newTeoZone110901CreateRequest(t)
	resp := &resource.CreateResponse{
		State: tfsdk.State{Raw: req.Plan.Raw, Schema: req.Plan.Schema},
	}

	r.Create(context.Background(), req, resp)
	assert.False(t, resp.Diagnostics.HasError(), "expected no error diagnostics, got: %v", resp.Diagnostics)

	// request assertions
	assert.NotNil(t, gotCreate)
	assert.Equal(t, "partial", *gotCreate.Type)
	assert.Equal(t, "example.com", *gotCreate.ZoneName)
	assert.Equal(t, "overseas", *gotCreate.Area)
	assert.Equal(t, "edgeone-2kfv1h391n6w", *gotCreate.PlanId)
	assert.Equal(t, "my-alias", *gotCreate.AliasZoneName)
	assert.True(t, *gotCreate.AllowDuplicates)
	assert.False(t, *gotCreate.JumpStart)
	assert.Len(t, gotCreate.Tags, 1)
	assert.Equal(t, "createdBy", *gotCreate.Tags[0].TagKey)
	assert.Equal(t, "terraform", *gotCreate.Tags[0].TagValue)
	// vanity_name_servers is ModifyZone-only; CreateZone has no such field.

	// describe must filter by zone-id
	assert.NotNil(t, gotDescribe)
	assert.Len(t, gotDescribe.Filters, 1)
	assert.Equal(t, "zone-id", *gotDescribe.Filters[0].Name)
	assert.Equal(t, "zone-2noz78a8ev6k", *gotDescribe.Filters[0].Values[0])

	// state assertions
	var model teo.TeoZone110901ResourceModel
	diags := resp.State.Get(context.Background(), &model)
	assert.False(t, diags.HasError(), "failed to read state: %v", diags)

	assert.Equal(t, "zone-2noz78a8ev6k", model.Id.ValueString())
	assert.Equal(t, "zone-2noz78a8ev6k", model.ZoneId.ValueString())
	assert.Equal(t, "partial", model.Type.ValueString())
	assert.Equal(t, "example.com", model.ZoneName.ValueString())
	assert.Equal(t, "overseas", model.Area.ValueString())
	assert.Equal(t, "edgeone-2kfv1h391n6w", model.PlanId.ValueString())
	assert.Equal(t, "my-alias", model.AliasZoneName.ValueString())
	assert.Equal(t, "pending", model.Status.ValueString())
	assert.Equal(t, "pending", model.CnameStatus.ValueString())
	assert.Equal(t, "inactive", model.ActiveStatus.ValueString())
	assert.Equal(t, "enable", model.LockStatus.ValueString())
	assert.False(t, model.Paused.ValueBool())
	assert.Equal(t, "2025-01-01T00:00:00+08:00", model.CreatedOn.ValueString())
	assert.Equal(t, "2025-01-01T00:00:00+08:00", model.ModifiedOn.ValueString())

	// tags set read back
	assert.False(t, model.Tags.IsNull())
	var tagElems []types.Object
	diags = model.Tags.ElementsAs(context.Background(), &tagElems, false)
	assert.False(t, diags.HasError())
	assert.Len(t, tagElems, 1)
	assert.Equal(t, "createdBy", tagElems[0].Attributes()["tag_key"].(types.String).ValueString())
	assert.Equal(t, "terraform", tagElems[0].Attributes()["tag_value"].(types.String).ValueString())

	// name servers
	assert.False(t, model.NameServers.IsNull())
	var nsVals []types.String
	diags = model.NameServers.ElementsAs(context.Background(), &nsVals, false)
	assert.False(t, diags.HasError())
	assert.Len(t, nsVals, 1)
	assert.Equal(t, "ns1.teo.example.com", nsVals[0].ValueString())

	assert.False(t, model.OriginalNameServers.IsNull())
	var onsVals []types.String
	diags = model.OriginalNameServers.ElementsAs(context.Background(), &onsVals, false)
	assert.False(t, diags.HasError())
	assert.Len(t, onsVals, 1)
	assert.Equal(t, "ns1.old.example.com", onsVals[0].ValueString())

	// ownership verification full structure (types.Object attribute access)
	assert.False(t, model.OwnershipVerification.IsNull())
	ovAttrs := model.OwnershipVerification.Attributes()

	assert.False(t, ovAttrs["dns_verification"].IsNull())
	dnsAttrs := ovAttrs["dns_verification"].(types.Object).Attributes()
	assert.Equal(t, "_teo.example.com", dnsAttrs["subdomain"].(types.String).ValueString())
	assert.Equal(t, "TXT", dnsAttrs["record_type"].(types.String).ValueString())
	assert.Equal(t, "teo-site-verification=abc123", dnsAttrs["record_value"].(types.String).ValueString())

	assert.False(t, ovAttrs["file_verification"].IsNull())
	fileAttrs := ovAttrs["file_verification"].(types.Object).Attributes()
	assert.Equal(t, "/.well-known/teo-verification/z12h416twn.txt", fileAttrs["path"].(types.String).ValueString())
	assert.Equal(t, "content-123", fileAttrs["content"].(types.String).ValueString())

	assert.False(t, ovAttrs["ns_verification"].IsNull())
	nsAttrs := ovAttrs["ns_verification"].(types.Object).Attributes()
	nsServers := nsAttrs["name_servers"].(types.Set)
	assert.False(t, nsServers.IsNull())
	var nsvVals []types.String
	diags = nsServers.ElementsAs(context.Background(), &nsvVals, false)
	assert.False(t, diags.HasError())
	assert.Len(t, nsvVals, 2)
	assert.Equal(t, "ns1.teo.example.com", nsvVals[0].ValueString())
	assert.Equal(t, "ns2.teo.example.com", nsvVals[1].ValueString())

	// vanity name servers read back (types.Object attribute access)
	assert.False(t, model.VanityNameServers.IsNull())
	vnsAttrs := model.VanityNameServers.Attributes()
	assert.Equal(t, "on", vnsAttrs["switch"].(types.String).ValueString())
	vnsServers := vnsAttrs["servers"].(types.Set)
	assert.False(t, vnsServers.IsNull())
	var vnsServerVals []types.String
	diags = vnsServers.ElementsAs(context.Background(), &vnsServerVals, false)
	assert.False(t, diags.HasError())
	assert.Len(t, vnsServerVals, 1)
	assert.Equal(t, "ns1.example.com", vnsServerVals[0].ValueString())

	// allow_duplicates / jump_start are not returned by DescribeZones: the
	// plan value is kept in state (true / false as planned).
	assert.True(t, model.AllowDuplicates.ValueBool(), "allow_duplicates should keep the planned value")
	assert.False(t, model.JumpStart.ValueBool(), "jump_start should keep the planned value")
}

// TestTeoZone110901Resource_Create_APIError covers the create API failure path.
func TestTeoZone110901Resource_Create_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(client, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "CreateZone", func(request *teov20220901.CreateZoneRequest) (*teov20220901.CreateZoneResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=FailedOperation.ZoneExisted, Message=zone existed")
	})

	r := &teo.TeoZone110901Resource{}
	setTeoZone110901Client(r, client)

	req := newTeoZone110901CreateRequest(t)
	resp := &resource.CreateResponse{State: tfsdk.State{Raw: req.Plan.Raw, Schema: req.Plan.Schema}}

	r.Create(context.Background(), req, resp)

	assert.True(t, resp.Diagnostics.HasError())
	assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Error creating TEO zone")
}

// TestTeoZone110901Resource_Create_EmptyResponse covers the empty create
// response path (nil response / empty zone id must not be written to state).
func TestTeoZone110901Resource_Create_EmptyResponse(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(client, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "CreateZone", func(request *teov20220901.CreateZoneRequest) (*teov20220901.CreateZoneResponse, error) {
		resp := teov20220901.NewCreateZoneResponse()
		resp.Response = &teov20220901.CreateZoneResponseParams{}
		return resp, nil
	})

	r := &teo.TeoZone110901Resource{}
	setTeoZone110901Client(r, client)

	req := newTeoZone110901CreateRequest(t)
	resp := &resource.CreateResponse{State: tfsdk.State{Raw: req.Plan.Raw, Schema: req.Plan.Schema}}

	r.Create(context.Background(), req, resp)

	assert.True(t, resp.Diagnostics.HasError())
	assert.Contains(t, resp.Diagnostics.Errors()[0].Detail(), "empty zone id")
}

// TestTeoZone110901Resource_Create_NotReadableAfterCreate covers the case
// where the zone is created but the read-back finds nothing.
func TestTeoZone110901Resource_Create_NotReadableAfterCreate(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(client, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "CreateZone", func(request *teov20220901.CreateZoneRequest) (*teov20220901.CreateZoneResponse, error) {
		resp := teov20220901.NewCreateZoneResponse()
		resp.Response = &teov20220901.CreateZoneResponseParams{
			ZoneId: ptrStringZone1109("zone-2noz78a8ev6k"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(teoClient, "DescribeZones", func(request *teov20220901.DescribeZonesRequest) (*teov20220901.DescribeZonesResponse, error) {
		resp := teov20220901.NewDescribeZonesResponse()
		resp.Response = &teov20220901.DescribeZonesResponseParams{
			TotalCount: ptrInt64Zone1109(0),
			Zones:      []*teov20220901.Zone{},
		}
		return resp, nil
	})

	r := &teo.TeoZone110901Resource{}
	setTeoZone110901Client(r, client)

	req := newTeoZone110901CreateRequest(t)
	resp := &resource.CreateResponse{State: tfsdk.State{Raw: req.Plan.Raw, Schema: req.Plan.Schema}}

	r.Create(context.Background(), req, resp)

	assert.True(t, resp.Diagnostics.HasError())
	assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "not readable after create")
}

// newTeoZone110901ReadRequest builds a ReadRequest with the zone id in state.
func newTeoZone110901ReadRequest(t *testing.T, zoneId string) resource.ReadRequest {
	s := teoZone110901ResourceSchema(t)
	vals := map[string]tftypes.Value{}
	for name, typ := range teoZone110901AttrTypes(t) {
		vals[name] = tftypes.NewValue(typ, nil)
	}
	vals["id"] = tftypes.NewValue(tftypes.String, zoneId)
	vals["zone_id"] = tftypes.NewValue(tftypes.String, zoneId)
	raw := tftypes.NewValue(tftypes.Object{AttributeTypes: teoZone110901AttrTypes(t)}, vals)
	req := resource.ReadRequest{}
	req.State.Raw = raw
	req.State.Schema = s
	return req
}

// TestTeoZone110901Resource_Read_Success covers a successful read.
func TestTeoZone110901Resource_Read_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(client, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeZones", func(request *teov20220901.DescribeZonesRequest) (*teov20220901.DescribeZonesResponse, error) {
		resp := teov20220901.NewDescribeZonesResponse()
		resp.Response = &teov20220901.DescribeZonesResponseParams{
			TotalCount: ptrInt64Zone1109(1),
			Zones:      []*teov20220901.Zone{mockTeoZone110901Zone("zone-2noz78a8ev6k")},
		}
		return resp, nil
	})

	r := &teo.TeoZone110901Resource{}
	setTeoZone110901Client(r, client)

	req := newTeoZone110901ReadRequest(t, "zone-2noz78a8ev6k")
	resp := &resource.ReadResponse{
		State: tfsdk.State{Raw: req.State.Raw, Schema: req.State.Schema},
	}

	r.Read(context.Background(), req, resp)
	assert.False(t, resp.Diagnostics.HasError(), "expected no error diagnostics, got: %v", resp.Diagnostics)

	var model teo.TeoZone110901ResourceModel
	diags := resp.State.Get(context.Background(), &model)
	assert.False(t, diags.HasError())
	assert.Equal(t, "zone-2noz78a8ev6k", model.ZoneId.ValueString())
	assert.Equal(t, "example.com", model.ZoneName.ValueString())
	assert.Equal(t, "partial", model.Type.ValueString())
	assert.Equal(t, "edgeone-2kfv1h391n6w", model.PlanId.ValueString())
	assert.Equal(t, "pending", model.Status.ValueString())
}

// TestTeoZone110901Resource_Read_NotFound covers the not-found path: a
// successful query with no matching zone removes the resource from state.
func TestTeoZone110901Resource_Read_NotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(client, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeZones", func(request *teov20220901.DescribeZonesRequest) (*teov20220901.DescribeZonesResponse, error) {
		resp := teov20220901.NewDescribeZonesResponse()
		resp.Response = &teov20220901.DescribeZonesResponseParams{
			TotalCount: ptrInt64Zone1109(0),
			Zones:      []*teov20220901.Zone{},
		}
		return resp, nil
	})

	r := &teo.TeoZone110901Resource{}
	setTeoZone110901Client(r, client)

	req := newTeoZone110901ReadRequest(t, "zone-2noz78a8ev6k")
	resp := &resource.ReadResponse{
		State: tfsdk.State{Raw: req.State.Raw, Schema: req.State.Schema},
	}

	r.Read(context.Background(), req, resp)
	assert.False(t, resp.Diagnostics.HasError())
	assert.True(t, resp.State.Raw.IsNull(), "state should have been removed")
}

// TestTeoZone110901Resource_Read_MismatchedId covers the case where the
// response contains a zone with a different id: it must not be treated as
// this resource (falls through to not-found removal).
func TestTeoZone110901Resource_Read_MismatchedId(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(client, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeZones", func(request *teov20220901.DescribeZonesRequest) (*teov20220901.DescribeZonesResponse, error) {
		other := mockTeoZone110901Zone("zone-9999999999999")
		resp := teov20220901.NewDescribeZonesResponse()
		resp.Response = &teov20220901.DescribeZonesResponseParams{
			TotalCount: ptrInt64Zone1109(1),
			Zones:      []*teov20220901.Zone{other},
		}
		return resp, nil
	})

	r := &teo.TeoZone110901Resource{}
	setTeoZone110901Client(r, client)

	req := newTeoZone110901ReadRequest(t, "zone-2noz78a8ev6k")
	resp := &resource.ReadResponse{
		State: tfsdk.State{Raw: req.State.Raw, Schema: req.State.Schema},
	}

	r.Read(context.Background(), req, resp)
	assert.False(t, resp.Diagnostics.HasError())
	assert.True(t, resp.State.Raw.IsNull(), "state should have been removed when no zone matches the id")
}

// TestTeoZone110901Resource_Read_APIError covers the API failure path: the
// error is surfaced and the state is kept.
func TestTeoZone110901Resource_Read_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(client, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeZones", func(request *teov20220901.DescribeZonesRequest) (*teov20220901.DescribeZonesResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=InternalError, Message=boom")
	})

	r := &teo.TeoZone110901Resource{}
	setTeoZone110901Client(r, client)

	req := newTeoZone110901ReadRequest(t, "zone-2noz78a8ev6k")
	resp := &resource.ReadResponse{
		State: tfsdk.State{Raw: req.State.Raw, Schema: req.State.Schema},
	}

	r.Read(context.Background(), req, resp)

	assert.True(t, resp.Diagnostics.HasError())
	assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Error reading TEO zone")
	assert.False(t, resp.State.Raw.IsNull(), "state must be kept on API error")
}

// TestTeoZone110901Resource_Read_NilClient covers the unconfigured provider
// path.
func TestTeoZone110901Resource_Read_NilClient(t *testing.T) {
	r := &teo.TeoZone110901Resource{}
	req := newTeoZone110901ReadRequest(t, "zone-2noz78a8ev6k")
	resp := &resource.ReadResponse{
		State: tfsdk.State{Raw: req.State.Raw, Schema: req.State.Schema},
	}

	r.Read(context.Background(), req, resp)
	assert.True(t, resp.Diagnostics.HasError())
	assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Provider not configured")
}

// newTeoZone110901UpdateRequest builds an UpdateRequest with changed
// type, area, alias_zone_name and vanity_name_servers in the plan.
func newTeoZone110901UpdateRequest(t *testing.T, zoneId string) resource.UpdateRequest {
	s := teoZone110901ResourceSchema(t)
	attrTypes := teoZone110901AttrTypes(t)
	stateVals := map[string]tftypes.Value{}
	planVals := map[string]tftypes.Value{}
	for name, typ := range attrTypes {
		stateVals[name] = tftypes.NewValue(typ, nil)
		planVals[name] = tftypes.NewValue(typ, nil)
	}
	stateVals["id"] = tftypes.NewValue(tftypes.String, zoneId)
	stateVals["zone_id"] = tftypes.NewValue(tftypes.String, zoneId)
	stateVals["type"] = tftypes.NewValue(tftypes.String, "partial")
	stateVals["zone_name"] = tftypes.NewValue(tftypes.String, "example.com")
	stateVals["area"] = tftypes.NewValue(tftypes.String, "overseas")
	stateVals["alias_zone_name"] = tftypes.NewValue(tftypes.String, "my-alias")

	planVals["id"] = tftypes.NewValue(tftypes.String, zoneId)
	planVals["zone_id"] = tftypes.NewValue(tftypes.String, zoneId)
	planVals["type"] = tftypes.NewValue(tftypes.String, "full")
	planVals["zone_name"] = tftypes.NewValue(tftypes.String, "example.com")
	planVals["area"] = tftypes.NewValue(tftypes.String, "global")
	planVals["alias_zone_name"] = tftypes.NewValue(tftypes.String, "new-alias")
	planVals["vanity_name_servers"] = teoZone110901VanityObject("on", []string{"ns1.example.com"})

	req := resource.UpdateRequest{}
	req.State.Raw = tftypes.NewValue(tftypes.Object{AttributeTypes: attrTypes}, stateVals)
	req.State.Schema = s
	req.Plan.Raw = tftypes.NewValue(tftypes.Object{AttributeTypes: attrTypes}, planVals)
	req.Plan.Schema = s
	req.Config.Raw = req.Plan.Raw
	req.Config.Schema = s
	return req
}

// TestTeoZone110901Resource_Update_Success covers a successful update.
func TestTeoZone110901Resource_Update_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(client, "UseTeoClient", teoClient)

	var gotModify *teov20220901.ModifyZoneRequest
	patches.ApplyMethodFunc(teoClient, "ModifyZone", func(request *teov20220901.ModifyZoneRequest) (*teov20220901.ModifyZoneResponse, error) {
		gotModify = request
		resp := teov20220901.NewModifyZoneResponse()
		resp.Response = &teov20220901.ModifyZoneResponseParams{RequestId: ptrStringZone1109("fake-request-id")}
		return resp, nil
	})

	patches.ApplyMethodFunc(teoClient, "DescribeZones", func(request *teov20220901.DescribeZonesRequest) (*teov20220901.DescribeZonesResponse, error) {
		resp := teov20220901.NewDescribeZonesResponse()
		resp.Response = &teov20220901.DescribeZonesResponseParams{
			TotalCount: ptrInt64Zone1109(1),
			Zones:      []*teov20220901.Zone{mockTeoZone110901Zone("zone-2noz78a8ev6k")},
		}
		return resp, nil
	})

	r := &teo.TeoZone110901Resource{}
	setTeoZone110901Client(r, client)

	req := newTeoZone110901UpdateRequest(t, "zone-2noz78a8ev6k")
	resp := &resource.UpdateResponse{
		State: tfsdk.State{Raw: req.Plan.Raw, Schema: req.Plan.Schema},
	}

	r.Update(context.Background(), req, resp)
	assert.False(t, resp.Diagnostics.HasError(), "expected no error diagnostics, got: %v", resp.Diagnostics)

	// request assertions
	assert.NotNil(t, gotModify)
	assert.Equal(t, "zone-2noz78a8ev6k", *gotModify.ZoneId)
	assert.Equal(t, "full", *gotModify.Type)
	assert.Equal(t, "global", *gotModify.Area)
	assert.Equal(t, "new-alias", *gotModify.AliasZoneName)
	assert.NotNil(t, gotModify.VanityNameServers)
	assert.Equal(t, "on", *gotModify.VanityNameServers.Switch)
	assert.Len(t, gotModify.VanityNameServers.Servers, 1)
	assert.Equal(t, "ns1.example.com", *gotModify.VanityNameServers.Servers[0])

	// state assertions after read-back
	var model teo.TeoZone110901ResourceModel
	diags := resp.State.Get(context.Background(), &model)
	assert.False(t, diags.HasError())
	assert.Equal(t, "zone-2noz78a8ev6k", model.Id.ValueString())
	assert.Equal(t, "my-alias", model.AliasZoneName.ValueString())
}

// TestTeoZone110901Resource_Update_NullKeepsOmitted covers the keep-semantics:
// plan fields that are null are not sent to ModifyZone.
func TestTeoZone110901Resource_Update_NullKeepsOmitted(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(client, "UseTeoClient", teoClient)

	var gotModify *teov20220901.ModifyZoneRequest
	patches.ApplyMethodFunc(teoClient, "ModifyZone", func(request *teov20220901.ModifyZoneRequest) (*teov20220901.ModifyZoneResponse, error) {
		gotModify = request
		resp := teov20220901.NewModifyZoneResponse()
		resp.Response = &teov20220901.ModifyZoneResponseParams{}
		return resp, nil
	})

	patches.ApplyMethodFunc(teoClient, "DescribeZones", func(request *teov20220901.DescribeZonesRequest) (*teov20220901.DescribeZonesResponse, error) {
		resp := teov20220901.NewDescribeZonesResponse()
		resp.Response = &teov20220901.DescribeZonesResponseParams{
			TotalCount: ptrInt64Zone1109(1),
			Zones:      []*teov20220901.Zone{mockTeoZone110901Zone("zone-2noz78a8ev6k")},
		}
		return resp, nil
	})

	r := &teo.TeoZone110901Resource{}
	setTeoZone110901Client(r, client)

	// plan with only the id set: everything else null
	s := teoZone110901ResourceSchema(t)
	attrTypes := teoZone110901AttrTypes(t)
	vals := map[string]tftypes.Value{}
	for name, typ := range attrTypes {
		vals[name] = tftypes.NewValue(typ, nil)
	}
	vals["id"] = tftypes.NewValue(tftypes.String, "zone-2noz78a8ev6k")
	vals["zone_id"] = tftypes.NewValue(tftypes.String, "zone-2noz78a8ev6k")

	req := resource.UpdateRequest{}
	req.State.Raw = tftypes.NewValue(tftypes.Object{AttributeTypes: attrTypes}, vals)
	req.State.Schema = s
	req.Plan.Raw = tftypes.NewValue(tftypes.Object{AttributeTypes: attrTypes}, vals)
	req.Plan.Schema = s
	req.Config.Raw = req.Plan.Raw
	req.Config.Schema = s

	resp := &resource.UpdateResponse{
		State: tfsdk.State{Raw: req.Plan.Raw, Schema: req.Plan.Schema},
	}

	r.Update(context.Background(), req, resp)
	assert.False(t, resp.Diagnostics.HasError(), "expected no error diagnostics, got: %v", resp.Diagnostics)

	// all updatable fields must be omitted (keep cloud-side values)
	assert.NotNil(t, gotModify)
	assert.Nil(t, gotModify.Type)
	assert.Nil(t, gotModify.ZoneName)
	assert.Nil(t, gotModify.Area)
	assert.Nil(t, gotModify.AliasZoneName)
	assert.Nil(t, gotModify.VanityNameServers)
}

// TestTeoZone110901Resource_Update_APIError covers the update API failure.
func TestTeoZone110901Resource_Update_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(client, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "ModifyZone", func(request *teov20220901.ModifyZoneRequest) (*teov20220901.ModifyZoneResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=ResourceInUse, Message=busy")
	})

	r := &teo.TeoZone110901Resource{}
	setTeoZone110901Client(r, client)

	req := newTeoZone110901UpdateRequest(t, "zone-2noz78a8ev6k")
	resp := &resource.UpdateResponse{
		State: tfsdk.State{Raw: req.Plan.Raw, Schema: req.Plan.Schema},
	}

	r.Update(context.Background(), req, resp)

	assert.True(t, resp.Diagnostics.HasError())
	assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Error updating TEO zone")
}

// TestTeoZone110901Resource_Delete_Success covers a successful delete.
func TestTeoZone110901Resource_Delete_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(client, "UseTeoClient", teoClient)

	var gotDelete *teov20220901.DeleteZoneRequest
	patches.ApplyMethodFunc(teoClient, "DeleteZone", func(request *teov20220901.DeleteZoneRequest) (*teov20220901.DeleteZoneResponse, error) {
		gotDelete = request
		resp := teov20220901.NewDeleteZoneResponse()
		resp.Response = &teov20220901.DeleteZoneResponseParams{}
		return resp, nil
	})

	r := &teo.TeoZone110901Resource{}
	setTeoZone110901Client(r, client)

	readReq := newTeoZone110901ReadRequest(t, "zone-2noz78a8ev6k")
	delReq := resource.DeleteRequest{
		State: tfsdk.State{Raw: readReq.State.Raw, Schema: readReq.State.Schema},
	}
	resp := &resource.DeleteResponse{}

	r.Delete(context.Background(), delReq, resp)
	assert.False(t, resp.Diagnostics.HasError(), "expected no error diagnostics, got: %v", resp.Diagnostics)

	assert.NotNil(t, gotDelete)
	assert.Equal(t, "zone-2noz78a8ev6k", *gotDelete.ZoneId)
}

// TestTeoZone110901Resource_Delete_APIError covers the delete API failure:
// the error must surface (no documented "already gone" code exists).
func TestTeoZone110901Resource_Delete_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(client, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "DeleteZone", func(request *teov20220901.DeleteZoneRequest) (*teov20220901.DeleteZoneResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=ResourceNotFound, Message=zone not found")
	})

	r := &teo.TeoZone110901Resource{}
	setTeoZone110901Client(r, client)

	readReq := newTeoZone110901ReadRequest(t, "zone-2noz78a8ev6k")
	delReq := resource.DeleteRequest{
		State: tfsdk.State{Raw: readReq.State.Raw, Schema: readReq.State.Schema},
	}
	resp := &resource.DeleteResponse{}

	r.Delete(context.Background(), delReq, resp)

	assert.True(t, resp.Diagnostics.HasError())
	assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Error deleting TEO zone")
}
