package teo_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	resource_schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"

	teo "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/sharedmeta"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/teo"
)

func ptrStringTeoAliasDomain110901(s string) *string { return &s }

func ptrInt64TeoAliasDomain110901(v int64) *int64 { return &v }

// teoAliasDomain110901Schema returns the framework resource schema used to
// build raw config/plan/state values.
func teoAliasDomain110901Schema(t *testing.T) resource_schema.Schema {
	r := teo.NewTeoAliasDomain110901Resource()
	schemaResp := &resource.SchemaResponse{}
	r.Schema(context.Background(), resource.SchemaRequest{}, schemaResp)
	if schemaResp.Diagnostics.HasError() {
		t.Fatalf("failed to build resource schema: %v", schemaResp.Diagnostics)
	}
	return schemaResp.Schema
}

// teoAliasDomain110901ObjectType returns the tftypes object type of the
// resource schema.
func teoAliasDomain110901ObjectType(t *testing.T) tftypes.Type {
	return teoAliasDomain110901Schema(t).Type().TerraformType(context.Background())
}

// teoAliasDomain110901RawValue builds a raw tftypes.Value of the resource
// object type from the given attribute values.
func teoAliasDomain110901RawValue(t *testing.T, attrs map[string]tftypes.Value) tftypes.Value {
	return tftypes.NewValue(teoAliasDomain110901ObjectType(t), attrs)
}

// newTeoAliasDomain110901CreateRequest builds a resource.CreateRequest whose
// config and plan carry the given attribute values.
func newTeoAliasDomain110901CreateRequest(t *testing.T, zoneId, aliasName, targetName string, certType *string, certIds []string) resource.CreateRequest {
	s := teoAliasDomain110901Schema(t)
	objType := teoAliasDomain110901ObjectType(t)
	nestedListType := objType.(tftypes.Object).AttributeTypes["domains"]

	var certTypeVal tftypes.Value
	if certType == nil {
		certTypeVal = tftypes.NewValue(tftypes.String, nil)
	} else {
		certTypeVal = tftypes.NewValue(tftypes.String, *certType)
	}

	var certIdVal tftypes.Value
	if certIds == nil {
		certIdVal = tftypes.NewValue(tftypes.Set{ElementType: tftypes.String}, nil)
	} else {
		elems := make([]tftypes.Value, 0, len(certIds))
		for _, c := range certIds {
			elems = append(elems, tftypes.NewValue(tftypes.String, c))
		}
		certIdVal = tftypes.NewValue(tftypes.Set{ElementType: tftypes.String}, elems)
	}

	attrs := map[string]tftypes.Value{
		"id":          tftypes.NewValue(tftypes.String, nil),
		"zone_id":     tftypes.NewValue(tftypes.String, zoneId),
		"alias_name":  tftypes.NewValue(tftypes.String, aliasName),
		"target_name": tftypes.NewValue(tftypes.String, targetName),
		"cert_type":   certTypeVal,
		"cert_id":     certIdVal,
		"domains":     tftypes.NewValue(nestedListType, nil),
	}

	raw := teoAliasDomain110901RawValue(t, attrs)

	req := resource.CreateRequest{}
	req.Config.Raw = raw
	req.Config.Schema = s
	req.Plan.Raw = raw
	req.Plan.Schema = s
	return req
}

// newTeoAliasDomain110901ReadRequest builds a resource.ReadRequest whose
// state carries only the given resource id.
func newTeoAliasDomain110901ReadRequest(t *testing.T, id string) resource.ReadRequest {
	s := teoAliasDomain110901Schema(t)
	objType := teoAliasDomain110901ObjectType(t)
	nestedListType := objType.(tftypes.Object).AttributeTypes["domains"]

	req := resource.ReadRequest{}
	req.State = tfsdk.State{
		Raw: teoAliasDomain110901RawValue(t, map[string]tftypes.Value{
			"id":          tftypes.NewValue(tftypes.String, id),
			"zone_id":     tftypes.NewValue(tftypes.String, nil),
			"alias_name":  tftypes.NewValue(tftypes.String, nil),
			"target_name": tftypes.NewValue(tftypes.String, nil),
			"cert_type":   tftypes.NewValue(tftypes.String, nil),
			"cert_id":     tftypes.NewValue(tftypes.Set{ElementType: tftypes.String}, nil),
			"domains":     tftypes.NewValue(nestedListType, nil),
		}),
		Schema: s,
	}
	return req
}

// newTeoAliasDomain110901UpdateRequest builds an resource.UpdateRequest whose
// plan matches the new config and whose state carries the prior values.
func newTeoAliasDomain110901UpdateRequest(t *testing.T, zoneId, aliasName, targetName string) resource.UpdateRequest {
	s := teoAliasDomain110901Schema(t)
	objType := teoAliasDomain110901ObjectType(t)
	nestedListType := objType.(tftypes.Object).AttributeTypes["domains"]

	stateAttrs := map[string]tftypes.Value{
		"id":          tftypes.NewValue(tftypes.String, zoneId+"#"+aliasName),
		"zone_id":     tftypes.NewValue(tftypes.String, zoneId),
		"alias_name":  tftypes.NewValue(tftypes.String, aliasName),
		"target_name": tftypes.NewValue(tftypes.String, targetName),
		"cert_type":   tftypes.NewValue(tftypes.String, nil),
		"cert_id":     tftypes.NewValue(tftypes.Set{ElementType: tftypes.String}, nil),
		"domains":     tftypes.NewValue(nestedListType, nil),
	}
	planAttrs := map[string]tftypes.Value{
		"id":          tftypes.NewValue(tftypes.String, zoneId+"#"+aliasName),
		"zone_id":     tftypes.NewValue(tftypes.String, zoneId),
		"alias_name":  tftypes.NewValue(tftypes.String, aliasName),
		"target_name": tftypes.NewValue(tftypes.String, targetName),
		"cert_type":   tftypes.NewValue(tftypes.String, nil),
		"cert_id":     tftypes.NewValue(tftypes.Set{ElementType: tftypes.String}, nil),
		"domains":     tftypes.NewValue(nestedListType, nil),
	}

	req := resource.UpdateRequest{}
	req.Config.Raw = teoAliasDomain110901RawValue(t, planAttrs)
	req.Config.Schema = s
	req.Plan.Raw = teoAliasDomain110901RawValue(t, planAttrs)
	req.Plan.Schema = s
	req.State = tfsdk.State{
		Raw:    teoAliasDomain110901RawValue(t, stateAttrs),
		Schema: s,
	}
	return req
}

// newTeoAliasDomain110901DeleteRequest builds a resource.DeleteRequest whose
// state carries only the given resource id.
func newTeoAliasDomain110901DeleteRequest(t *testing.T, id string) resource.DeleteRequest {
	s := teoAliasDomain110901Schema(t)
	objType := teoAliasDomain110901ObjectType(t)
	nestedListType := objType.(tftypes.Object).AttributeTypes["domains"]

	req := resource.DeleteRequest{}
	req.State = tfsdk.State{
		Raw: teoAliasDomain110901RawValue(t, map[string]tftypes.Value{
			"id":          tftypes.NewValue(tftypes.String, id),
			"zone_id":     tftypes.NewValue(tftypes.String, nil),
			"alias_name":  tftypes.NewValue(tftypes.String, nil),
			"target_name": tftypes.NewValue(tftypes.String, nil),
			"cert_type":   tftypes.NewValue(tftypes.String, nil),
			"cert_id":     tftypes.NewValue(tftypes.Set{ElementType: tftypes.String}, nil),
			"domains":     tftypes.NewValue(nestedListType, nil),
		}),
		Schema: s,
	}
	return req
}

// newTeoAliasDomain110901Response builds a response with a state primed from
// the request state/plan so that resp.State.Set can be asserted afterwards.
func newTeoAliasDomain110901CreateResponse(req resource.CreateRequest) *resource.CreateResponse {
	resp := &resource.CreateResponse{}
	resp.State.Schema = req.Plan.Schema
	resp.State.Raw = req.Plan.Raw
	return resp
}

func newTeoAliasDomain110901ReadResponse(req resource.ReadRequest) *resource.ReadResponse {
	resp := &resource.ReadResponse{}
	resp.State = req.State
	return resp
}

func newTeoAliasDomain110901UpdateResponse(req resource.UpdateRequest) *resource.UpdateResponse {
	resp := &resource.UpdateResponse{}
	resp.State.Schema = req.Plan.Schema
	resp.State.Raw = req.Plan.Raw
	return resp
}

// setTeoAliasDomain110901Client injects a mock connectivity client into the
// resource instance via the Configure method.
func setTeoAliasDomain110901Client(r *teo.TeoAliasDomain110901Resource, client *connectivity.TencentCloudClient) {
	meta := &sharedmeta.ProviderMeta{Client: client}
	configureReq := resource.ConfigureRequest{ProviderData: meta}
	configureResp := &resource.ConfigureResponse{}
	r.Configure(context.Background(), configureReq, configureResp)
}

// newTeoAliasDomain110901ResourceUnderTest returns a configured resource
// instance plus the mock teo client used for patching.
func newTeoAliasDomain110901ResourceUnderTest(patches *gomonkey.Patches) (*teo.TeoAliasDomain110901Resource, *teo.Client) {
	client := &connectivity.TencentCloudClient{}
	teoClient := &teo.Client{}
	patches.ApplyMethodReturn(client, "UseTeoV20220901Client", teoClient)

	r := teo.NewTeoAliasDomain110901Resource().(*teo.TeoAliasDomain110901Resource)
	setTeoAliasDomain110901Client(r, client)
	return r, teoClient
}

// mockDescribeAliasDomain110901 patches DescribeAliasDomains to return the
// given alias domain, or none when found is false.
func mockDescribeAliasDomain110901(patches *gomonkey.Patches, teoClient *teo.Client, zoneId, aliasName string, found bool) {
	aliasDomains := []*teo.AliasDomain{}
	if found {
		aliasDomains = append(aliasDomains, &teo.AliasDomain{
			AliasName:  ptrStringTeoAliasDomain110901(aliasName),
			ZoneId:     ptrStringTeoAliasDomain110901(zoneId),
			TargetName: ptrStringTeoAliasDomain110901("target.demo.com"),
			Status:     ptrStringTeoAliasDomain110901("active"),
		})
	}
	resp := teo.NewDescribeAliasDomainsResponse()
	resp.Response = &teo.DescribeAliasDomainsResponseParams{
		TotalCount:   ptrInt64TeoAliasDomain110901(int64(len(aliasDomains))),
		AliasDomains: aliasDomains,
	}
	patches.ApplyMethodFunc(teoClient, "DescribeAliasDomains", func(_ context.Context, request *teo.DescribeAliasDomainsRequest) (*teo.DescribeAliasDomainsResponse, error) {
		return resp, nil
	})
}

// mockDescribeInferenceDomains110901 patches DescribeInferenceDomains to
// return the given inference domain list.
func mockDescribeInferenceDomains110901(patches *gomonkey.Patches, teoClient *teo.Client, domains []*teo.InferenceDomain) {
	resp := &teo.DescribeInferenceDomainsResponse{
		Response: &teo.DescribeInferenceDomainsResponseParams{
			TotalCount: ptrInt64TeoAliasDomain110901(int64(len(domains))),
			Domains:    domains,
		},
	}
	patches.ApplyMethodFunc(teoClient, "DescribeInferenceDomainsWithContext", func(_ context.Context, request *teo.DescribeInferenceDomainsRequest) (*teo.DescribeInferenceDomainsResponse, error) {
		return resp, nil
	})
}

// fullInferenceDomain110901 returns a DescribeInferenceDomains item with
// every leaf field populated.
func fullInferenceDomain110901() *teo.InferenceDomain {
	return &teo.InferenceDomain{
		Domain:     ptrStringTeoAliasDomain110901("inference.example.com"),
		Status:     ptrStringTeoAliasDomain110901("Online"),
		AuthSwitch: ptrStringTeoAliasDomain110901("On"),
		Cname:      ptrStringTeoAliasDomain110901("cn.example.cdn.dnsv1.com"),
		OwnershipVerification: &teo.OwnershipVerification{
			DnsVerification: &teo.DnsVerification{
				Subdomain:   ptrStringTeoAliasDomain110901("_teo_verify"),
				RecordType:  ptrStringTeoAliasDomain110901("TXT"),
				RecordValue: ptrStringTeoAliasDomain110901("20250101120000abc"),
			},
			FileVerification: &teo.FileVerification{
				Path:    ptrStringTeoAliasDomain110901("/.well-known/teo-verification/abc.txt"),
				Content: ptrStringTeoAliasDomain110901("verify-content"),
			},
			NsVerification: &teo.NsVerification{
				NameServers: []*string{
					ptrStringTeoAliasDomain110901("ns1.example.dnsv1.com"),
					ptrStringTeoAliasDomain110901("ns2.example.dnsv1.com"),
				},
			},
		},
		Certificate: &teo.HostsCertificate{
			Host: ptrStringTeoAliasDomain110901("inference.example.com"),
			Mode: ptrStringTeoAliasDomain110901("sslcert"),
			HostCertInfo: []*teo.HostCertInfo{
				{
					CertId:     ptrStringTeoAliasDomain110901("cert-1234"),
					Alias:      ptrStringTeoAliasDomain110901("my cert"),
					Type:       ptrStringTeoAliasDomain110901("upload"),
					ExpireTime: ptrStringTeoAliasDomain110901("2026-01-01"),
					DeployTime: ptrStringTeoAliasDomain110901("2025-01-02"),
					SignAlgo:   ptrStringTeoAliasDomain110901("SHA256"),
					Status:     ptrStringTeoAliasDomain110901("deployed"),
				},
			},
			ApplyType: ptrStringTeoAliasDomain110901("apply"),
		},
		CreateTime: ptrStringTeoAliasDomain110901("2025-01-01 00:00:00"),
		UpdateTime: ptrStringTeoAliasDomain110901("2025-01-02 00:00:00"),
	}
}

// decodeTeoAliasDomain110901State decodes a tfsdk.State into the model.
func decodeTeoAliasDomain110901State(t *testing.T, state tfsdk.State) teo.TeoAliasDomain110901ResourceModel {
	t.Helper()
	var model teo.TeoAliasDomain110901ResourceModel
	if diags := state.Get(context.Background(), &model); diags.HasError() {
		t.Fatalf("failed to decode state: %v", diags)
	}
	return model
}

// assertFullDomainsList110901 asserts every leaf field of the computed
// domains list against the fixture returned by fullInferenceDomain110901.
func assertFullDomainsList110901(t *testing.T, domains types.List) {
	t.Helper()
	if domains.IsNull() {
		t.Fatal("domains list must not be null")
	}
	elems := domains.Elements()
	if len(elems) != 1 {
		t.Fatalf("expected 1 domains element, got %d", len(elems))
	}
	obj, ok := elems[0].(types.Object)
	if !ok {
		t.Fatalf("expected domains[0] to be types.Object, got %T", elems[0])
	}
	attrs := obj.Attributes()

	if got := attrs["domain"].(types.String).ValueString(); got != "inference.example.com" {
		t.Errorf("domains[0].domain mismatch: %s", got)
	}
	if got := attrs["status"].(types.String).ValueString(); got != "Online" {
		t.Errorf("domains[0].status mismatch: %s", got)
	}
	if got := attrs["auth_switch"].(types.String).ValueString(); got != "On" {
		t.Errorf("domains[0].auth_switch mismatch: %s", got)
	}
	if got := attrs["cname"].(types.String).ValueString(); got != "cn.example.cdn.dnsv1.com" {
		t.Errorf("domains[0].cname mismatch: %s", got)
	}
	if got := attrs["create_time"].(types.String).ValueString(); got != "2025-01-01 00:00:00" {
		t.Errorf("domains[0].create_time mismatch: %s", got)
	}
	if got := attrs["update_time"].(types.String).ValueString(); got != "2025-01-02 00:00:00" {
		t.Errorf("domains[0].update_time mismatch: %s", got)
	}

	ov, ok := attrs["ownership_verification"].(types.Object)
	if !ok {
		t.Fatalf("expected ownership_verification to be types.Object, got %T", attrs["ownership_verification"])
	}
	ovAttrs := ov.Attributes()
	dns, ok := ovAttrs["dns_verification"].(types.Object)
	if !ok {
		t.Fatalf("expected dns_verification to be types.Object, got %T", ovAttrs["dns_verification"])
	}
	dnsAttrs := dns.Attributes()
	if got := dnsAttrs["subdomain"].(types.String).ValueString(); got != "_teo_verify" {
		t.Errorf("dns_verification.subdomain mismatch: %s", got)
	}
	if got := dnsAttrs["record_type"].(types.String).ValueString(); got != "TXT" {
		t.Errorf("dns_verification.record_type mismatch: %s", got)
	}
	if got := dnsAttrs["record_value"].(types.String).ValueString(); got != "20250101120000abc" {
		t.Errorf("dns_verification.record_value mismatch: %s", got)
	}
	fileV, ok := ovAttrs["file_verification"].(types.Object)
	if !ok {
		t.Fatalf("expected file_verification to be types.Object, got %T", ovAttrs["file_verification"])
	}
	fileAttrs := fileV.Attributes()
	if got := fileAttrs["path"].(types.String).ValueString(); got != "/.well-known/teo-verification/abc.txt" {
		t.Errorf("file_verification.path mismatch: %s", got)
	}
	if got := fileAttrs["content"].(types.String).ValueString(); got != "verify-content" {
		t.Errorf("file_verification.content mismatch: %s", got)
	}
	ns, ok := ovAttrs["ns_verification"].(types.Object)
	if !ok {
		t.Fatalf("expected ns_verification to be types.Object, got %T", ovAttrs["ns_verification"])
	}
	nsAttrs := ns.Attributes()
	nsSet, ok := nsAttrs["name_servers"].(types.Set)
	if !ok {
		t.Fatalf("expected name_servers to be types.Set, got %T", nsAttrs["name_servers"])
	}
	nsElems := nsSet.Elements()
	if len(nsElems) != 2 {
		t.Fatalf("expected 2 name_servers, got %d", len(nsElems))
	}
	if got := nsElems[0].(types.String).ValueString(); got != "ns1.example.dnsv1.com" {
		t.Errorf("name_servers[0] mismatch: %s", got)
	}
	if got := nsElems[1].(types.String).ValueString(); got != "ns2.example.dnsv1.com" {
		t.Errorf("name_servers[1] mismatch: %s", got)
	}

	cert, ok := attrs["certificate"].(types.Object)
	if !ok {
		t.Fatalf("expected certificate to be types.Object, got %T", attrs["certificate"])
	}
	certAttrs := cert.Attributes()
	if got := certAttrs["host"].(types.String).ValueString(); got != "inference.example.com" {
		t.Errorf("certificate.host mismatch: %s", got)
	}
	if got := certAttrs["mode"].(types.String).ValueString(); got != "sslcert" {
		t.Errorf("certificate.mode mismatch: %s", got)
	}
	if got := certAttrs["apply_type"].(types.String).ValueString(); got != "apply" {
		t.Errorf("certificate.apply_type mismatch: %s", got)
	}
	hciList, ok := certAttrs["host_cert_info"].(types.List)
	if !ok {
		t.Fatalf("expected host_cert_info to be types.List, got %T", certAttrs["host_cert_info"])
	}
	hciElems := hciList.Elements()
	if len(hciElems) != 1 {
		t.Fatalf("expected 1 host_cert_info element, got %d", len(hciElems))
	}
	hci, ok := hciElems[0].(types.Object)
	if !ok {
		t.Fatalf("expected host_cert_info[0] to be types.Object, got %T", hciElems[0])
	}
	hciAttrs := hci.Attributes()
	if got := hciAttrs["cert_id"].(types.String).ValueString(); got != "cert-1234" {
		t.Errorf("host_cert_info[0].cert_id mismatch: %s", got)
	}
	if got := hciAttrs["alias"].(types.String).ValueString(); got != "my cert" {
		t.Errorf("host_cert_info[0].alias mismatch: %s", got)
	}
	if got := hciAttrs["type"].(types.String).ValueString(); got != "upload" {
		t.Errorf("host_cert_info[0].type mismatch: %s", got)
	}
	if got := hciAttrs["expire_time"].(types.String).ValueString(); got != "2026-01-01" {
		t.Errorf("host_cert_info[0].expire_time mismatch: %s", got)
	}
	if got := hciAttrs["deploy_time"].(types.String).ValueString(); got != "2025-01-02" {
		t.Errorf("host_cert_info[0].deploy_time mismatch: %s", got)
	}
	if got := hciAttrs["sign_algo"].(types.String).ValueString(); got != "SHA256" {
		t.Errorf("host_cert_info[0].sign_algo mismatch: %s", got)
	}
	if got := hciAttrs["status"].(types.String).ValueString(); got != "deployed" {
		t.Errorf("host_cert_info[0].status mismatch: %s", got)
	}
}

// go test ./tencentcloud/services/teo/ -run "TestTeoAliasDomain110901" -v -count=1 -gcflags="all=-l"

func TestTeoAliasDomain110901_MetadataAndSchema(t *testing.T) {
	r := teo.NewTeoAliasDomain110901Resource()

	metaResp := &resource.MetadataResponse{}
	r.Metadata(context.Background(), resource.MetadataRequest{}, metaResp)
	if metaResp.TypeName != "tencentcloud_teo_alias_domain_1109_01" {
		t.Fatalf("unexpected TypeName: %s", metaResp.TypeName)
	}

	schemaResp := &resource.SchemaResponse{}
	r.Schema(context.Background(), resource.SchemaRequest{}, schemaResp)
	if schemaResp.Diagnostics.HasError() {
		t.Fatalf("schema error: %v", schemaResp.Diagnostics)
	}

	attrs := schemaResp.Schema.Attributes
	for _, name := range []string{"id", "zone_id", "alias_name", "target_name", "cert_type", "cert_id", "domains"} {
		if _, ok := attrs[name]; !ok {
			t.Errorf("attribute %s missing from schema", name)
		}
	}
	if !attrs["zone_id"].(resource_schema.StringAttribute).Required {
		t.Errorf("zone_id should be required")
	}
	if !attrs["domains"].(resource_schema.ListNestedAttribute).Computed {
		t.Errorf("domains should be computed")
	}
}

func TestTeoAliasDomain110901_Create_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	var gotCreateRequest *teo.CreateAliasDomainRequest
	patchesRef := patches

	r, teoClient := newTeoAliasDomain110901ResourceUnderTest(patchesRef)
	patchesRef.ApplyMethodFunc(teoClient, "CreateAliasDomainWithContext", func(_ context.Context, request *teo.CreateAliasDomainRequest) (*teo.CreateAliasDomainResponse, error) {
		gotCreateRequest = request
		resp := teo.NewCreateAliasDomainResponse()
		resp.Response = &teo.CreateAliasDomainResponseParams{RequestId: ptrStringTeoAliasDomain110901("fake-request-id")}
		return resp, nil
	})
	mockDescribeAliasDomain110901(patchesRef, teoClient, "zone-3fkff38fyw8s", "alias.demo.com", true)
	mockDescribeInferenceDomains110901(patchesRef, teoClient, []*teo.InferenceDomain{fullInferenceDomain110901()})

	req := newTeoAliasDomain110901CreateRequest(t, "zone-3fkff38fyw8s", "alias.demo.com", "target.demo.com", ptrStringTeoAliasDomain110901("hosting"), []string{"cert-1234"})
	resp := newTeoAliasDomain110901CreateResponse(req)
	r.Create(context.Background(), req, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected error diagnostics: %v", resp.Diagnostics)
	}

	if gotCreateRequest == nil {
		t.Fatal("CreateAliasDomainWithContext was not called")
	}
	if gotCreateRequest.ZoneId == nil || *gotCreateRequest.ZoneId != "zone-3fkff38fyw8s" {
		t.Errorf("ZoneId mismatch: %v", gotCreateRequest.ZoneId)
	}
	if gotCreateRequest.AliasName == nil || *gotCreateRequest.AliasName != "alias.demo.com" {
		t.Errorf("AliasName mismatch: %v", gotCreateRequest.AliasName)
	}
	if gotCreateRequest.TargetName == nil || *gotCreateRequest.TargetName != "target.demo.com" {
		t.Errorf("TargetName mismatch: %v", gotCreateRequest.TargetName)
	}
	if gotCreateRequest.CertType == nil || *gotCreateRequest.CertType != "hosting" {
		t.Errorf("CertType mismatch: %v", gotCreateRequest.CertType)
	}
	if len(gotCreateRequest.CertId) != 1 || *gotCreateRequest.CertId[0] != "cert-1234" {
		t.Errorf("CertId mismatch: %v", gotCreateRequest.CertId)
	}

	model := decodeTeoAliasDomain110901State(t, resp.State)
	if model.Id.ValueString() != "zone-3fkff38fyw8s#alias.demo.com" {
		t.Errorf("id mismatch: %s", model.Id.ValueString())
	}
	if model.TargetName.ValueString() != "target.demo.com" {
		t.Errorf("target_name mismatch: %s", model.TargetName.ValueString())
	}
	if model.CertType.ValueString() != "hosting" {
		t.Errorf("cert_type mismatch: %s", model.CertType.ValueString())
	}
	certIdElems := model.CertId.Elements()
	if len(certIdElems) != 1 {
		t.Fatalf("cert_id should have 1 element, got %d", len(certIdElems))
	}
	if got := certIdElems[0].(types.String).ValueString(); got != "cert-1234" {
		t.Errorf("cert_id[0] mismatch: %s", got)
	}
	assertFullDomainsList110901(t, model.Domains)
}

func TestTeoAliasDomain110901_Create_NoOptionalFields(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	var gotCreateRequest *teo.CreateAliasDomainRequest
	r, teoClient := newTeoAliasDomain110901ResourceUnderTest(patches)
	patches.ApplyMethodFunc(teoClient, "CreateAliasDomainWithContext", func(_ context.Context, request *teo.CreateAliasDomainRequest) (*teo.CreateAliasDomainResponse, error) {
		gotCreateRequest = request
		resp := teo.NewCreateAliasDomainResponse()
		resp.Response = &teo.CreateAliasDomainResponseParams{RequestId: ptrStringTeoAliasDomain110901("fake-request-id")}
		return resp, nil
	})
	mockDescribeAliasDomain110901(patches, teoClient, "zone-3fkff38fyw8s", "alias.demo.com", true)
	mockDescribeInferenceDomains110901(patches, teoClient, []*teo.InferenceDomain{})

	req := newTeoAliasDomain110901CreateRequest(t, "zone-3fkff38fyw8s", "alias.demo.com", "target.demo.com", nil, nil)
	resp := newTeoAliasDomain110901CreateResponse(req)
	r.Create(context.Background(), req, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected error diagnostics: %v", resp.Diagnostics)
	}
	if gotCreateRequest.CertType != nil {
		t.Errorf("CertType should be nil, got %v", *gotCreateRequest.CertType)
	}
	if gotCreateRequest.CertId != nil {
		t.Errorf("CertId should be nil, got %v", gotCreateRequest.CertId)
	}

	model := decodeTeoAliasDomain110901State(t, resp.State)
	if model.Id.ValueString() != "zone-3fkff38fyw8s#alias.demo.com" {
		t.Errorf("id mismatch: %s", model.Id.ValueString())
	}
	if !model.Domains.IsNull() && len(model.Domains.Elements()) != 0 {
		t.Errorf("expected empty domains list, got %v", model.Domains)
	}
}

func TestTeoAliasDomain110901_Create_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	r, teoClient := newTeoAliasDomain110901ResourceUnderTest(patches)
	patches.ApplyMethodFunc(teoClient, "CreateAliasDomainWithContext", func(_ context.Context, request *teo.CreateAliasDomainRequest) (*teo.CreateAliasDomainResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=LimitExceeded, Message=quota exceeded")
	})

	req := newTeoAliasDomain110901CreateRequest(t, "zone-3fkff38fyw8s", "alias.demo.com", "target.demo.com", nil, nil)
	resp := newTeoAliasDomain110901CreateResponse(req)
	r.Create(context.Background(), req, resp)

	if !resp.Diagnostics.HasError() {
		t.Fatal("expected error diagnostics")
	}
	if got := resp.Diagnostics.Errors()[0].Summary(); got != "Error creating TEO alias domain" {
		t.Errorf("unexpected summary: %s", got)
	}
}

func TestTeoAliasDomain110901_Create_ClientNotConfigured(t *testing.T) {
	r := teo.NewTeoAliasDomain110901Resource().(*teo.TeoAliasDomain110901Resource)

	req := newTeoAliasDomain110901CreateRequest(t, "zone-3fkff38fyw8s", "alias.demo.com", "target.demo.com", nil, nil)
	resp := newTeoAliasDomain110901CreateResponse(req)
	r.Create(context.Background(), req, resp)

	if !resp.Diagnostics.HasError() {
		t.Fatal("expected error diagnostics")
	}
	if got := resp.Diagnostics.Errors()[0].Summary(); got != "Provider not configured" {
		t.Errorf("unexpected summary: %s", got)
	}
}

func TestTeoAliasDomain110901_Read_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	r, teoClient := newTeoAliasDomain110901ResourceUnderTest(patches)
	mockDescribeAliasDomain110901(patches, teoClient, "zone-3fkff38fyw8s", "alias.demo.com", true)
	mockDescribeInferenceDomains110901(patches, teoClient, []*teo.InferenceDomain{fullInferenceDomain110901()})

	req := newTeoAliasDomain110901ReadRequest(t, "zone-3fkff38fyw8s#alias.demo.com")
	resp := newTeoAliasDomain110901ReadResponse(req)
	r.Read(context.Background(), req, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected error diagnostics: %v", resp.Diagnostics)
	}

	model := decodeTeoAliasDomain110901State(t, resp.State)
	if model.ZoneId.ValueString() != "zone-3fkff38fyw8s" {
		t.Errorf("zone_id mismatch: %s", model.ZoneId.ValueString())
	}
	if model.AliasName.ValueString() != "alias.demo.com" {
		t.Errorf("alias_name mismatch: %s", model.AliasName.ValueString())
	}
	if model.TargetName.ValueString() != "target.demo.com" {
		t.Errorf("target_name mismatch: %s", model.TargetName.ValueString())
	}
	assertFullDomainsList110901(t, model.Domains)
}

func TestTeoAliasDomain110901_Read_NotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	r, teoClient := newTeoAliasDomain110901ResourceUnderTest(patches)
	mockDescribeAliasDomain110901(patches, teoClient, "zone-3fkff38fyw8s", "alias.demo.com", false)
	mockDescribeInferenceDomains110901(patches, teoClient, []*teo.InferenceDomain{})

	req := newTeoAliasDomain110901ReadRequest(t, "zone-3fkff38fyw8s#alias.demo.com")
	resp := newTeoAliasDomain110901ReadResponse(req)
	r.Read(context.Background(), req, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected error diagnostics: %v", resp.Diagnostics)
	}
	if !resp.State.Raw.IsNull() {
		t.Errorf("expected the resource to be removed from state, got %v", resp.State.Raw)
	}
}

func TestTeoAliasDomain110901_Read_EmptyDomainsList(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	r, teoClient := newTeoAliasDomain110901ResourceUnderTest(patches)
	mockDescribeAliasDomain110901(patches, teoClient, "zone-3fkff38fyw8s", "alias.demo.com", true)
	// An empty domains list is a valid business state, not an error.
	mockDescribeInferenceDomains110901(patches, teoClient, []*teo.InferenceDomain{})

	req := newTeoAliasDomain110901ReadRequest(t, "zone-3fkff38fyw8s#alias.demo.com")
	resp := newTeoAliasDomain110901ReadResponse(req)
	r.Read(context.Background(), req, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected error diagnostics: %v", resp.Diagnostics)
	}
	model := decodeTeoAliasDomain110901State(t, resp.State)
	if !model.Domains.IsNull() && len(model.Domains.Elements()) != 0 {
		t.Errorf("expected empty domains list, got %v", model.Domains)
	}
}

func TestTeoAliasDomain110901_Read_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	r, teoClient := newTeoAliasDomain110901ResourceUnderTest(patches)
	patches.ApplyMethodFunc(teoClient, "DescribeAliasDomains", func(_ context.Context, request *teo.DescribeAliasDomainsRequest) (*teo.DescribeAliasDomainsResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=InternalError, Message=server error")
	})

	req := newTeoAliasDomain110901ReadRequest(t, "zone-3fkff38fyw8s#alias.demo.com")
	resp := newTeoAliasDomain110901ReadResponse(req)
	r.Read(context.Background(), req, resp)

	if !resp.Diagnostics.HasError() {
		t.Fatal("expected error diagnostics")
	}
	if !resp.State.Raw.IsNull() {
		t.Error("state must not be removed on an API error")
	}
}

func TestTeoAliasDomain110901_Update_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	r, teoClient := newTeoAliasDomain110901ResourceUnderTest(patches)
	mockDescribeAliasDomain110901(patches, teoClient, "zone-3fkff38fyw8s", "alias.demo.com", true)
	mockDescribeInferenceDomains110901(patches, teoClient, []*teo.InferenceDomain{fullInferenceDomain110901()})

	req := newTeoAliasDomain110901UpdateRequest(t, "zone-3fkff38fyw8s", "alias.demo.com", "target.demo.com")
	resp := newTeoAliasDomain110901UpdateResponse(req)
	r.Update(context.Background(), req, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected error diagnostics: %v", resp.Diagnostics)
	}
	model := decodeTeoAliasDomain110901State(t, resp.State)
	if model.Id.ValueString() != "zone-3fkff38fyw8s#alias.demo.com" {
		t.Errorf("id mismatch: %s", model.Id.ValueString())
	}
	assertFullDomainsList110901(t, model.Domains)
}

func TestTeoAliasDomain110901_Delete_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	var gotDeleteRequest *teo.DeleteAliasDomainRequest
	r, teoClient := newTeoAliasDomain110901ResourceUnderTest(patches)
	patches.ApplyMethodFunc(teoClient, "DeleteAliasDomainWithContext", func(_ context.Context, request *teo.DeleteAliasDomainRequest) (*teo.DeleteAliasDomainResponse, error) {
		gotDeleteRequest = request
		resp := teo.NewDeleteAliasDomainResponse()
		resp.Response = &teo.DeleteAliasDomainResponseParams{RequestId: ptrStringTeoAliasDomain110901("fake-request-id")}
		return resp, nil
	})

	req := newTeoAliasDomain110901DeleteRequest(t, "zone-3fkff38fyw8s#alias.demo.com")
	resp := &resource.DeleteResponse{}
	resp.State = req.State
	r.Delete(context.Background(), req, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected error diagnostics: %v", resp.Diagnostics)
	}
	if gotDeleteRequest == nil {
		t.Fatal("DeleteAliasDomainWithContext was not called")
	}
	if gotDeleteRequest.ZoneId == nil || *gotDeleteRequest.ZoneId != "zone-3fkff38fyw8s" {
		t.Errorf("ZoneId mismatch: %v", gotDeleteRequest.ZoneId)
	}
	if len(gotDeleteRequest.AliasNames) != 1 || *gotDeleteRequest.AliasNames[0] != "alias.demo.com" {
		t.Errorf("AliasNames mismatch: %v", gotDeleteRequest.AliasNames)
	}
}

func TestTeoAliasDomain110901_Delete_BrokenId(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	r, _ := newTeoAliasDomain110901ResourceUnderTest(patches)

	req := newTeoAliasDomain110901DeleteRequest(t, "broken-id-without-separator")
	resp := &resource.DeleteResponse{}
	resp.State = req.State
	r.Delete(context.Background(), req, resp)

	if !resp.Diagnostics.HasError() {
		t.Fatal("expected error diagnostics for a broken id")
	}
}

func TestTeoAliasDomain110901_ImportState(t *testing.T) {
	r := teo.NewTeoAliasDomain110901Resource().(*teo.TeoAliasDomain110901Resource)

	schemaResp := &resource.SchemaResponse{}
	r.Schema(context.Background(), resource.SchemaRequest{}, schemaResp)
	if schemaResp.Diagnostics.HasError() {
		t.Fatalf("schema error: %v", schemaResp.Diagnostics)
	}

	req := resource.ImportStateRequest{ID: "zone-3fkff38fyw8s#alias.demo.com"}
	resp := &resource.ImportStateResponse{}
	resp.State.Schema = schemaResp.Schema
	resp.State.Raw = tftypes.NewValue(schemaResp.Schema.Type().TerraformType(context.Background()), nil)
	r.ImportState(context.Background(), req, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected import diagnostics: %v", resp.Diagnostics)
	}

	var model teo.TeoAliasDomain110901ResourceModel
	if diags := resp.State.Get(context.Background(), &model); diags.HasError() {
		t.Fatalf("failed to decode imported state: %v", diags)
	}
	if model.Id.ValueString() != "zone-3fkff38fyw8s#alias.demo.com" {
		t.Errorf("imported id mismatch: %s", model.Id.ValueString())
	}
}

func TestTeoAliasDomain110901_SplitIdRoundTrip(t *testing.T) {
	cases := []struct {
		id     string
		zone   string
		alias  string
		wantOk bool
	}{
		{"zone-1#alias.example.com", "zone-1", "alias.example.com", true},
		{"zone-1", "", "", false},
		{"#alias.example.com", "", "", false},
		{"zone-1#", "", "", false},
	}
	for _, c := range cases {
		zone, alias, ok := teo.SplitTeoAliasDomain110901IdForTest(c.id)
		if ok != c.wantOk {
			t.Errorf("id %q: ok = %v, want %v", c.id, ok, c.wantOk)
			continue
		}
		if ok && (zone != c.zone || alias != c.alias) {
			t.Errorf("id %q: got (%q, %q), want (%q, %q)", c.id, zone, alias, c.zone, c.alias)
		}
	}
}
