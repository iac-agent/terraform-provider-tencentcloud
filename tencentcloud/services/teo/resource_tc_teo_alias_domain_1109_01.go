package teo

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	teo "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/fw"
)

// teoAliasDomain110901StateWriter is implemented by the per-response adapter
// types so that refreshFromRead can write state and diagnostics without
// generics (the repo targets Go 1.21).
type teoAliasDomain110901StateWriter interface {
	// created reports whether this writer belongs to a Create call; only
	// then is "not readable right after create" an error instead of a
	// plain "not found" removal.
	created() bool
	addError(summary, detail string)
	set(ctx context.Context, m *TeoAliasDomain110901ResourceModel)
	removeResource(ctx context.Context)
}

// _ resource.Resource is intentionally omitted: the factory function
// NewTeoAliasDomain110901Resource returns resource.Resource, which already
// performs the implicit compile-time check at the return statement. Only the
// extended-interface assertions are kept here.
var (
	_ resource.ResourceWithConfigure     = &TeoAliasDomain110901Resource{}
	_ resource.ResourceWithImportState = &TeoAliasDomain110901Resource{}
)

// NewTeoAliasDomain110901Resource is the factory referenced by
// tencentcloud/framework/registry.go to register this resource.
func NewTeoAliasDomain110901Resource() resource.Resource {
	return &TeoAliasDomain110901Resource{}
}

// TeoAliasDomain110901Resource implements resource.Resource for
// tencentcloud_teo_alias_domain_1109_01.
type TeoAliasDomain110901Resource struct {
	fw.ResourceWithConfigure
	fw.WithImportByID
}

// TeoAliasDomain110901ResourceModel maps the schema attributes.
type TeoAliasDomain110901ResourceModel struct {
	Id         types.String `tfsdk:"id"`
	ZoneId     types.String `tfsdk:"zone_id"`
	AliasName  types.String `tfsdk:"alias_name"`
	TargetName types.String `tfsdk:"target_name"`
	CertType   types.String `tfsdk:"cert_type"`
	CertId     types.Set    `tfsdk:"cert_id"`
	Domains    types.List   `tfsdk:"domains"`
}

// teoAliasDomain110901DomainModel maps one element of the computed
// `domains` list, mirroring the DescribeInferenceDomains output.
type teoAliasDomain110901DomainModel struct {
	Domain                 types.String `tfsdk:"domain"`
	Status                 types.String `tfsdk:"status"`
	AuthSwitch             types.String `tfsdk:"auth_switch"`
	Cname                  types.String `tfsdk:"cname"`
	OwnershipVerification types.Object `tfsdk:"ownership_verification"`
	Certificate            types.Object `tfsdk:"certificate"`
	CreateTime             types.String `tfsdk:"create_time"`
	UpdateTime             types.String `tfsdk:"update_time"`
}

// teoAliasDomain110901OwnershipVerificationModel mirrors the
// OwnershipVerification object of the DescribeInferenceDomains output.
type teoAliasDomain110901OwnershipVerificationModel struct {
	DnsVerification  types.Object `tfsdk:"dns_verification"`
	FileVerification types.Object `tfsdk:"file_verification"`
	NsVerification   types.Object `tfsdk:"ns_verification"`
}

// teoAliasDomain110901DnsVerificationModel mirrors the DnsVerification object.
type teoAliasDomain110901DnsVerificationModel struct {
	Subdomain   types.String `tfsdk:"subdomain"`
	RecordType  types.String `tfsdk:"record_type"`
	RecordValue types.String `tfsdk:"record_value"`
}

// teoAliasDomain110901FileVerificationModel mirrors the FileVerification object.
type teoAliasDomain110901FileVerificationModel struct {
	Path    types.String `tfsdk:"path"`
	Content types.String `tfsdk:"content"`
}

// teoAliasDomain110901NsVerificationModel mirrors the NsVerification object.
type teoAliasDomain110901NsVerificationModel struct {
	NameServers types.Set `tfsdk:"name_servers"`
}

// teoAliasDomain110901CertificateModel mirrors the Certificate object.
type teoAliasDomain110901CertificateModel struct {
	Host         types.String `tfsdk:"host"`
	Mode         types.String `tfsdk:"mode"`
	HostCertInfo types.List   `tfsdk:"host_cert_info"`
	ApplyType    types.String `tfsdk:"apply_type"`
}

// teoAliasDomain110901HostCertInfoModel mirrors one HostCertInfo element.
type teoAliasDomain110901HostCertInfoModel struct {
	CertId     types.String `tfsdk:"cert_id"`
	Alias      types.String `tfsdk:"alias"`
	Type       types.String `tfsdk:"type"`
	ExpireTime types.String `tfsdk:"expire_time"`
	DeployTime types.String `tfsdk:"deploy_time"`
	SignAlgo   types.String `tfsdk:"sign_algo"`
	Status     types.String `tfsdk:"status"`
}

func (r *TeoAliasDomain110901Resource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "tencentcloud_teo_alias_domain_1109_01"
}

func (r *TeoAliasDomain110901Resource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a resource to create a TEO alias domain " +
			"(terraform-plugin-framework implementation).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Resource identifier in the form `<zone_id>#<alias_name>`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"zone_id": schema.StringAttribute{
				Required:    true,
				Description: "Site ID.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"alias_name": schema.StringAttribute{
				Required:    true,
				Description: "Alias domain name.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"target_name": schema.StringAttribute{
				Required:    true,
				Description: "Target domain name.",
			},
			"cert_type": schema.StringAttribute{
				Optional: true,
				Description: "Certificate configuration. Valid values: `none` (no configuration), " +
					"`hosting` (SSL hosted certificate). Default value: `none`.",
			},
			"cert_id": schema.SetAttribute{
				Optional:    true,
				ElementType: types.StringType,
				Description: "Certificate ID list. Required when `cert_type` is `hosting`.",
			},
			"domains": schema.ListNestedAttribute{
				Computed:    true,
				Description: "Inference service domain list returned by DescribeInferenceDomains.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"domain": schema.StringAttribute{
							Computed:    true,
							Description: "Domain name.",
						},
						"status": schema.StringAttribute{
							Computed:    true,
							Description: "Domain status.",
						},
						"auth_switch": schema.StringAttribute{
							Computed:    true,
							Description: "Inference task request authentication switch.",
						},
						"cname": schema.StringAttribute{
							Computed:    true,
							Description: "CNAME address.",
						},
						"ownership_verification": schema.SingleNestedAttribute{
							Computed:    true,
							Description: "Ownership verification information.",
							Attributes: map[string]schema.Attribute{
								"dns_verification": schema.SingleNestedAttribute{
									Computed:    true,
									Description: "DNS resolution verification information.",
									Attributes: map[string]schema.Attribute{
										"subdomain": schema.StringAttribute{
											Computed:    true,
											Description: "Host record.",
										},
										"record_type": schema.StringAttribute{
											Computed:    true,
											Description: "Record type.",
										},
										"record_value": schema.StringAttribute{
											Computed:    true,
											Description: "Record value.",
										},
									},
								},
								"file_verification": schema.SingleNestedAttribute{
									Computed:    true,
									Description: "File verification information.",
									Attributes: map[string]schema.Attribute{
										"path": schema.StringAttribute{
											Computed:    true,
											Description: "URL path of the verification file.",
										},
										"content": schema.StringAttribute{
											Computed:    true,
											Description: "Content of the verification file.",
										},
									},
								},
								"ns_verification": schema.SingleNestedAttribute{
									Computed:    true,
									Description: "NS access verification information.",
									Attributes: map[string]schema.Attribute{
										"name_servers": schema.SetAttribute{
											Computed:    true,
											ElementType: types.StringType,
											Description: "DNS server addresses assigned to the user.",
										},
									},
								},
							},
						},
						"certificate": schema.SingleNestedAttribute{
							Computed:    true,
							Description: "Domain certificate information.",
							Attributes: map[string]schema.Attribute{
								"host": schema.StringAttribute{
									Computed:    true,
									Description: "Domain name.",
								},
								"mode": schema.StringAttribute{
									Computed:    true,
									Description: "Certificate configuration mode.",
								},
								"host_cert_info": schema.ListNestedAttribute{
									Computed:    true,
									Description: "Server certificate configuration.",
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{
											"cert_id": schema.StringAttribute{
												Computed:    true,
												Description: "Server certificate ID.",
											},
											"alias": schema.StringAttribute{
												Computed:    true,
												Description: "Certificate remark name.",
											},
											"type": schema.StringAttribute{
												Computed:    true,
												Description: "Certificate type.",
											},
											"expire_time": schema.StringAttribute{
												Computed:    true,
												Description: "Certificate expiration time.",
											},
											"deploy_time": schema.StringAttribute{
												Computed:    true,
												Description: "Certificate deployment time.",
											},
											"sign_algo": schema.StringAttribute{
												Computed:    true,
												Description: "Signature algorithm.",
											},
											"status": schema.StringAttribute{
												Computed:    true,
												Description: "Certificate status.",
											},
										},
									},
								},
								"apply_type": schema.StringAttribute{
									Computed:    true,
									Description: "Application type.",
								},
							},
						},
						"create_time": schema.StringAttribute{
							Computed:    true,
							Description: "Creation time.",
						},
						"update_time": schema.StringAttribute{
							Computed:    true,
							Description: "Modification time.",
						},
					},
				},
			},
		},
	}
}

func (r *TeoAliasDomain110901Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan TeoAliasDomain110901ResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if r.Client() == nil {
		resp.Diagnostics.AddError(
			"Provider not configured",
			"The provider client is not available. This is unexpected; please report it.",
		)
		return
	}

	zoneId := plan.ZoneId.ValueString()
	aliasName := plan.AliasName.ValueString()
	targetName := plan.TargetName.ValueString()

	request := teo.NewCreateAliasDomainRequest()
	request.ZoneId = &zoneId
	request.AliasName = &aliasName
	request.TargetName = &targetName
	if !plan.CertType.IsNull() && !plan.CertType.IsUnknown() {
		request.CertType = plan.CertType.ValueStringPointer()
	}
	if !plan.CertId.IsNull() && !plan.CertId.IsUnknown() {
		certIds := make([]*string, 0, len(plan.CertId.Elements()))
		for _, v := range plan.CertId.Elements() {
			s, ok := v.(types.String)
			if !ok {
				continue
			}
			certIds = append(certIds, s.ValueStringPointer())
		}
		// Only send the field when the user configured at least one id; an
		// empty non-nil slice would be serialized as an explicit empty list.
		if len(certIds) > 0 {
			request.CertId = certIds
		}
	}

	if _, err := r.Client().UseTeoV20220901Client().CreateAliasDomainWithContext(ctx, request); err != nil {
		resp.Diagnostics.AddError(
			"Error creating TEO alias domain",
			fmt.Sprintf("CreateAliasDomain failed: %v", err),
		)
		return
	}

	plan.Id = types.StringValue(strings.Join([]string{zoneId, aliasName}, tccommon.FILED_SP))

	// Persist the identity BEFORE any read-back: if the read-back fails, the
	// resource is still in state and `terraform destroy` can clean it up.
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	r.refreshFromRead(ctx, &plan, newTeoAliasDomain110901CreateWriter(resp))
}

func (r *TeoAliasDomain110901Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state TeoAliasDomain110901ResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if r.Client() == nil {
		resp.Diagnostics.AddError(
			"Provider not configured",
			"The provider client is not available. This is unexpected; please report it.",
		)
		return
	}

	r.refreshFromRead(ctx, &state, newTeoAliasDomain110901ReadWriter(resp))
}

// Update handles in-place modifications. The only cloud-side update API in
// scope is ModifyAliasDomainStatus, which flips a `paused` switch that has
// no HCL representation in this resource; no writable attribute can change
// in place (target_name/cert_type/cert_id changes are not supported by the
// mapped Update API), so Update only refreshes the state from the cloud.
func (r *TeoAliasDomain110901Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan TeoAliasDomain110901ResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if r.Client() == nil {
		resp.Diagnostics.AddError(
			"Provider not configured",
			"The provider client is not available. This is unexpected; please report it.",
		)
		return
	}

	r.refreshFromRead(ctx, &plan, newTeoAliasDomain110901UpdateWriter(resp))
}

func (r *TeoAliasDomain110901Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state TeoAliasDomain110901ResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if r.Client() == nil {
		resp.Diagnostics.AddError(
			"Provider not configured",
			"The provider client is not available. This is unexpected; please report it.",
		)
		return
	}

	zoneId, aliasName, ok := teoSplitAliasDomain110901Id(state.Id.ValueString())
	if !ok {
		resp.Diagnostics.AddError(
			"Invalid resource ID",
			fmt.Sprintf("TEO alias domain id is broken, id is %q (expected `<zone_id>%s<alias_name>`)", state.Id.ValueString(), tccommon.FILED_SP),
		)
		return
	}

	request := teo.NewDeleteAliasDomainRequest()
	request.ZoneId = &zoneId
	request.AliasNames = []*string{&aliasName}

	if _, err := r.Client().UseTeoV20220901Client().DeleteAliasDomainWithContext(ctx, request); err != nil {
		resp.Diagnostics.AddError(
			"Error deleting TEO alias domain",
			fmt.Sprintf("DeleteAliasDomain failed: %v", err),
		)
		return
	}
}

// refreshFromRead re-reads the alias domain (DescribeAliasDomains) and the
// inference domain list (DescribeInferenceDomains), then writes the merged
// model back through the writer. When the alias domain no longer exists, the
// resource is removed from state.
func (r *TeoAliasDomain110901Resource) refreshFromRead(ctx context.Context, model *TeoAliasDomain110901ResourceModel, writer teoAliasDomain110901StateWriter) {
	zoneId, aliasName, ok := teoSplitAliasDomain110901Id(model.Id.ValueString())
	if !ok {
		writer.addError(
			"Invalid resource ID",
			fmt.Sprintf("TEO alias domain id is broken, id is %q (expected `<zone_id>%s<alias_name>`)", model.Id.ValueString(), tccommon.FILED_SP),
		)
		return
	}

	service := NewTeoService(r.Client())
	aliasDomain, err := service.DescribeTeoAliasDomainById(ctx, zoneId, aliasName)
	if err != nil {
		writer.addError(
			"Error reading TEO alias domain",
			fmt.Sprintf("DescribeAliasDomains failed for zone %q alias %q: %v", zoneId, aliasName, err),
		)
		return
	}

	// The alias domain no longer exists on the cloud side. During a
	// refresh/update this removes the resource from state so Terraform plans
	// a recreation; right after Create it means the created alias domain is
	// not readable back, which is reported as an error instead (the identity
	// is already in state so `terraform destroy` can clean it up).
	if aliasDomain == nil {
		if writer.created() {
			writer.addError(
				"TEO alias domain not readable after create",
				fmt.Sprintf("the alias domain %q was created but cannot be read back; the identity is kept in state so it can be destroyed", model.Id.ValueString()),
			)
			return
		}
		writer.removeResource(ctx)
		return
	}

	model.ZoneId = types.StringValue(zoneId)
	model.AliasName = types.StringValue(aliasName)
	if aliasDomain.TargetName != nil {
		model.TargetName = types.StringPointerValue(aliasDomain.TargetName)
	}
	// cert_type and cert_id are not returned by DescribeAliasDomains; keep the
	// user-configured values from the plan/state to avoid spurious diffs.

	domains, err := service.DescribeTeoInferenceDomains(ctx, zoneId)
	if err != nil {
		writer.addError(
			"Error reading TEO inference domains",
			fmt.Sprintf("DescribeInferenceDomains failed for zone %q: %v", zoneId, err),
		)
		return
	}

	domainModels, diags := teoAliasDomain110901DomainModelsFrom(ctx, domains)
	if diags.HasError() {
		writer.addError(
			"Error reading TEO inference domains",
			"converting the DescribeInferenceDomains response failed",
		)
		return
	}

	list, diags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: teoAliasDomain110901DomainAttrTypes()}, domainModels)
	if diags.HasError() {
		writer.addError(
			"Error reading TEO inference domains",
			"building the domains list value failed",
		)
		return
	}
	model.Domains = list

	writer.set(ctx, model)
}

// --- writer adapters -------------------------------------------------------

type teoAliasDomain110901CreateWriter struct{ resp *resource.CreateResponse }

func newTeoAliasDomain110901CreateWriter(resp *resource.CreateResponse) *teoAliasDomain110901CreateWriter {
	return &teoAliasDomain110901CreateWriter{resp: resp}
}

func (w *teoAliasDomain110901CreateWriter) created() bool { return true }

func (w *teoAliasDomain110901CreateWriter) addError(summary, detail string) {
	w.resp.Diagnostics.AddError(summary, detail)
}

func (w *teoAliasDomain110901CreateWriter) set(ctx context.Context, m *TeoAliasDomain110901ResourceModel) {
	w.resp.Diagnostics.Append(w.resp.State.Set(ctx, m)...)
}

func (w *teoAliasDomain110901CreateWriter) removeResource(ctx context.Context) {
	w.resp.State.RemoveResource(ctx)
}

type teoAliasDomain110901ReadWriter struct{ resp *resource.ReadResponse }

func newTeoAliasDomain110901ReadWriter(resp *resource.ReadResponse) *teoAliasDomain110901ReadWriter {
	return &teoAliasDomain110901ReadWriter{resp: resp}
}

func (w *teoAliasDomain110901ReadWriter) created() bool { return false }

func (w *teoAliasDomain110901ReadWriter) addError(summary, detail string) {
	w.resp.Diagnostics.AddError(summary, detail)
}

func (w *teoAliasDomain110901ReadWriter) set(ctx context.Context, m *TeoAliasDomain110901ResourceModel) {
	w.resp.Diagnostics.Append(w.resp.State.Set(ctx, m)...)
}

func (w *teoAliasDomain110901ReadWriter) removeResource(ctx context.Context) {
	w.resp.State.RemoveResource(ctx)
}

type teoAliasDomain110901UpdateWriter struct{ resp *resource.UpdateResponse }

func newTeoAliasDomain110901UpdateWriter(resp *resource.UpdateResponse) *teoAliasDomain110901UpdateWriter {
	return &teoAliasDomain110901UpdateWriter{resp: resp}
}

func (w *teoAliasDomain110901UpdateWriter) created() bool { return false }

func (w *teoAliasDomain110901UpdateWriter) addError(summary, detail string) {
	w.resp.Diagnostics.AddError(summary, detail)
}

func (w *teoAliasDomain110901UpdateWriter) set(ctx context.Context, m *TeoAliasDomain110901ResourceModel) {
	w.resp.Diagnostics.Append(w.resp.State.Set(ctx, m)...)
}

func (w *teoAliasDomain110901UpdateWriter) removeResource(ctx context.Context) {
	w.resp.State.RemoveResource(ctx)
}

// --- conversion helpers ----------------------------------------------------

// SplitTeoAliasDomain110901IdForTest exposes the id parsing used by the
// resource for unit tests in the external test package.
func SplitTeoAliasDomain110901IdForTest(id string) (string, string, bool) {
	return teoSplitAliasDomain110901Id(id)
}

// teoSplitAliasDomain110901Id parses `<zone_id>#<alias_name>`.
func teoSplitAliasDomain110901Id(id string) (string, string, bool) {
	parts := strings.Split(id, tccommon.FILED_SP)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", false
	}
	return parts[0], parts[1], true
}

func teoAliasDomain110901DomainAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"domain":                  types.StringType,
		"status":                  types.StringType,
		"auth_switch":             types.StringType,
		"cname":                   types.StringType,
		"ownership_verification":  types.ObjectType{AttrTypes: teoAliasDomain110901OwnershipVerificationAttrTypes()},
		"certificate":             types.ObjectType{AttrTypes: teoAliasDomain110901CertificateAttrTypes()},
		"create_time":             types.StringType,
		"update_time":             types.StringType,
	}
}

func teoAliasDomain110901OwnershipVerificationAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"dns_verification":  types.ObjectType{AttrTypes: teoAliasDomain110901DnsVerificationAttrTypes()},
		"file_verification": types.ObjectType{AttrTypes: teoAliasDomain110901FileVerificationAttrTypes()},
		"ns_verification":   types.ObjectType{AttrTypes: teoAliasDomain110901NsVerificationAttrTypes()},
	}
}

func teoAliasDomain110901DnsVerificationAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"subdomain":   types.StringType,
		"record_type": types.StringType,
		"record_value": types.StringType,
	}
}

func teoAliasDomain110901FileVerificationAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"path":    types.StringType,
		"content": types.StringType,
	}
}

func teoAliasDomain110901NsVerificationAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"name_servers": types.SetType{ElemType: types.StringType},
	}
}

func teoAliasDomain110901CertificateAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"host":          types.StringType,
		"mode":          types.StringType,
		"host_cert_info": types.ListType{ElemType: types.ObjectType{AttrTypes: teoAliasDomain110901HostCertInfoAttrTypes()}},
		"apply_type":    types.StringType,
	}
}

func teoAliasDomain110901HostCertInfoAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"cert_id":     types.StringType,
		"alias":       types.StringType,
		"type":        types.StringType,
		"expire_time": types.StringType,
		"deploy_time": types.StringType,
		"sign_algo":   types.StringType,
		"status":      types.StringType,
	}
}

// teoAliasDomain110901DomainModelsFrom converts the DescribeInferenceDomains
// response items into framework models. A nil slice becomes an empty list.
func teoAliasDomain110901DomainModelsFrom(ctx context.Context, domains []*teo.InferenceDomain) ([]teoAliasDomain110901DomainModel, diag.Diagnostics) {
	var diags diag.Diagnostics
	items := make([]teoAliasDomain110901DomainModel, 0, len(domains))
	for _, d := range domains {
		if d == nil {
			continue
		}
		item := teoAliasDomain110901DomainModel{
			Domain:     types.StringPointerValue(d.Domain),
			Status:     types.StringPointerValue(d.Status),
			AuthSwitch: types.StringPointerValue(d.AuthSwitch),
			Cname:      types.StringPointerValue(d.Cname),
			CreateTime: types.StringPointerValue(d.CreateTime),
			UpdateTime: types.StringPointerValue(d.UpdateTime),
		}
		item.OwnershipVerification = teoAliasDomain110901OwnershipVerificationObject(ctx, d.OwnershipVerification, &diags)
		item.Certificate = teoAliasDomain110901CertificateObject(ctx, d.Certificate, &diags)
		items = append(items, item)
	}
	return items, diags
}

func teoAliasDomain110901OwnershipVerificationObject(ctx context.Context, v *teo.OwnershipVerification, diags *diag.Diagnostics) types.Object {
	if v == nil {
		return types.ObjectNull(teoAliasDomain110901OwnershipVerificationAttrTypes())
	}
	model := teoAliasDomain110901OwnershipVerificationModel{
		DnsVerification:  teoAliasDomain110901DnsVerificationObject(ctx, v.DnsVerification, diags),
		FileVerification: teoAliasDomain110901FileVerificationObject(ctx, v.FileVerification, diags),
		NsVerification:   teoAliasDomain110901NsVerificationObject(ctx, v.NsVerification, diags),
	}
	obj, d := types.ObjectValueFrom(ctx, teoAliasDomain110901OwnershipVerificationAttrTypes(), model)
	diags.Append(d...)
	return obj
}

func teoAliasDomain110901DnsVerificationObject(ctx context.Context, v *teo.DnsVerification, diags *diag.Diagnostics) types.Object {
	if v == nil {
		return types.ObjectNull(teoAliasDomain110901DnsVerificationAttrTypes())
	}
	model := teoAliasDomain110901DnsVerificationModel{
		Subdomain:   types.StringPointerValue(v.Subdomain),
		RecordType:  types.StringPointerValue(v.RecordType),
		RecordValue: types.StringPointerValue(v.RecordValue),
	}
	obj, d := types.ObjectValueFrom(ctx, teoAliasDomain110901DnsVerificationAttrTypes(), model)
	diags.Append(d...)
	return obj
}

func teoAliasDomain110901FileVerificationObject(ctx context.Context, v *teo.FileVerification, diags *diag.Diagnostics) types.Object {
	if v == nil {
		return types.ObjectNull(teoAliasDomain110901FileVerificationAttrTypes())
	}
	model := teoAliasDomain110901FileVerificationModel{
		Path:    types.StringPointerValue(v.Path),
		Content: types.StringPointerValue(v.Content),
	}
	obj, d := types.ObjectValueFrom(ctx, teoAliasDomain110901FileVerificationAttrTypes(), model)
	diags.Append(d...)
	return obj
}

func teoAliasDomain110901NsVerificationObject(ctx context.Context, v *teo.NsVerification, diags *diag.Diagnostics) types.Object {
	if v == nil {
		return types.ObjectNull(teoAliasDomain110901NsVerificationAttrTypes())
	}
	nameServers := make([]types.String, 0, len(v.NameServers))
	for _, ns := range v.NameServers {
		nameServers = append(nameServers, types.StringPointerValue(ns))
	}
	setVal, d := types.SetValueFrom(ctx, types.StringType, nameServers)
	diags.Append(d...)
	if d.HasError() {
		return types.ObjectNull(teoAliasDomain110901NsVerificationAttrTypes())
	}
	model := teoAliasDomain110901NsVerificationModel{
		NameServers: setVal,
	}
	obj, d := types.ObjectValueFrom(ctx, teoAliasDomain110901NsVerificationAttrTypes(), model)
	diags.Append(d...)
	return obj
}

func teoAliasDomain110901CertificateObject(ctx context.Context, v *teo.HostsCertificate, diags *diag.Diagnostics) types.Object {
	if v == nil {
		return types.ObjectNull(teoAliasDomain110901CertificateAttrTypes())
	}
	hostCertInfos := make([]teoAliasDomain110901HostCertInfoModel, 0, len(v.HostCertInfo))
	for _, h := range v.HostCertInfo {
		if h == nil {
			continue
		}
		hostCertInfos = append(hostCertInfos, teoAliasDomain110901HostCertInfoModel{
			CertId:     types.StringPointerValue(h.CertId),
			Alias:      types.StringPointerValue(h.Alias),
			Type:       types.StringPointerValue(h.Type),
			ExpireTime: types.StringPointerValue(h.ExpireTime),
			DeployTime: types.StringPointerValue(h.DeployTime),
			SignAlgo:   types.StringPointerValue(h.SignAlgo),
			Status:     types.StringPointerValue(h.Status),
		})
	}
	hostCertInfoList, d := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: teoAliasDomain110901HostCertInfoAttrTypes()}, hostCertInfos)
	diags.Append(d...)
	if d.HasError() {
		return types.ObjectNull(teoAliasDomain110901CertificateAttrTypes())
	}
	model := teoAliasDomain110901CertificateModel{
		Host:         types.StringPointerValue(v.Host),
		Mode:         types.StringPointerValue(v.Mode),
		HostCertInfo: hostCertInfoList,
		ApplyType:    types.StringPointerValue(v.ApplyType),
	}
	obj, d := types.ObjectValueFrom(ctx, teoAliasDomain110901CertificateAttrTypes(), model)
	diags.Append(d...)
	return obj
}
