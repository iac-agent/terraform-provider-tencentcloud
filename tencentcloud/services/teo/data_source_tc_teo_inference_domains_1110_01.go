package teo

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/fw"
)

var _ datasource.DataSourceWithConfigure = &TeoInferenceDomainsDataSource{}

// NewTeoInferenceDomainsDataSource is the factory referenced by
// tencentcloud/framework/registry.go to register this data source.
func NewTeoInferenceDomainsDataSource() datasource.DataSource {
	return &TeoInferenceDomainsDataSource{}
}

// TeoInferenceDomainsDataSource implements datasource.DataSource for
// the tencentcloud_teo_inference_domains_1110_01 data source.
type TeoInferenceDomainsDataSource struct {
	fw.DataSourceWithConfigure
}

// TeoInferenceDomainsDataSourceModel maps the schema attributes.
type TeoInferenceDomainsDataSourceModel struct {
	ZoneId    types.String `tfsdk:"zone_id"`
	ServiceId types.String `tfsdk:"service_id"`
	SortBy    types.String `tfsdk:"sort_by"`
	SortOrder types.String `tfsdk:"sort_order"`
	Domains   types.List   `tfsdk:"domains"`
}

// TeoInferenceDomainsDnsVerificationModel is the DNS verification element model.
type TeoInferenceDomainsDnsVerificationModel struct {
	Subdomain   types.String `tfsdk:"subdomain"`
	RecordType  types.String `tfsdk:"record_type"`
	RecordValue types.String `tfsdk:"record_value"`
}

var TeoInferenceDomainsDnsVerificationAttrTypes = map[string]attr.Type{
	"subdomain":    types.StringType,
	"record_type":  types.StringType,
	"record_value": types.StringType,
}

// TeoInferenceDomainsFileVerificationModel is the file verification element model.
type TeoInferenceDomainsFileVerificationModel struct {
	Path    types.String `tfsdk:"path"`
	Content types.String `tfsdk:"content"`
}

var TeoInferenceDomainsFileVerificationAttrTypes = map[string]attr.Type{
	"path":    types.StringType,
	"content": types.StringType,
}

// TeoInferenceDomainsNsVerificationModel is the NS verification element model.
type TeoInferenceDomainsNsVerificationModel struct {
	NameServers types.Set `tfsdk:"name_servers"`
}

var TeoInferenceDomainsNsVerificationAttrTypes = map[string]attr.Type{
	"name_servers": types.SetType{ElemType: types.StringType},
}

// TeoInferenceDomainsOwnershipVerificationModel is the ownership verification element model.
type TeoInferenceDomainsOwnershipVerificationModel struct {
	DnsVerification  types.List `tfsdk:"dns_verification"`
	FileVerification types.List `tfsdk:"file_verification"`
	NsVerification   types.List `tfsdk:"ns_verification"`
}

var TeoInferenceDomainsOwnershipVerificationAttrTypes = map[string]attr.Type{
	"dns_verification":  types.ListType{ElemType: types.ObjectType{AttrTypes: TeoInferenceDomainsDnsVerificationAttrTypes}},
	"file_verification": types.ListType{ElemType: types.ObjectType{AttrTypes: TeoInferenceDomainsFileVerificationAttrTypes}},
	"ns_verification":   types.ListType{ElemType: types.ObjectType{AttrTypes: TeoInferenceDomainsNsVerificationAttrTypes}},
}

// TeoInferenceDomainsHostCertInfoModel is the host certificate info element model.
type TeoInferenceDomainsHostCertInfoModel struct {
	CertId     types.String `tfsdk:"cert_id"`
	Alias      types.String `tfsdk:"alias"`
	Type       types.String `tfsdk:"type"`
	ExpireTime types.String `tfsdk:"expire_time"`
	DeployTime types.String `tfsdk:"deploy_time"`
	SignAlgo   types.String `tfsdk:"sign_algo"`
	Status     types.String `tfsdk:"status"`
}

var TeoInferenceDomainsHostCertInfoAttrTypes = map[string]attr.Type{
	"cert_id":     types.StringType,
	"alias":       types.StringType,
	"type":        types.StringType,
	"expire_time": types.StringType,
	"deploy_time": types.StringType,
	"sign_algo":   types.StringType,
	"status":      types.StringType,
}

// TeoInferenceDomainsCertificateModel is the certificate element model.
type TeoInferenceDomainsCertificateModel struct {
	Host         types.String `tfsdk:"host"`
	Mode         types.String `tfsdk:"mode"`
	HostCertInfo types.List   `tfsdk:"host_cert_info"`
	ApplyType    types.String `tfsdk:"apply_type"`
}

var TeoInferenceDomainsCertificateAttrTypes = map[string]attr.Type{
	"host":           types.StringType,
	"mode":           types.StringType,
	"host_cert_info": types.ListType{ElemType: types.ObjectType{AttrTypes: TeoInferenceDomainsHostCertInfoAttrTypes}},
	"apply_type":     types.StringType,
}

// TeoInferenceDomainsDomainModel is the domain element model.
type TeoInferenceDomainsDomainModel struct {
	Domain                types.String `tfsdk:"domain"`
	Status                types.String `tfsdk:"status"`
	AuthSwitch            types.String `tfsdk:"auth_switch"`
	Cname                 types.String `tfsdk:"cname"`
	OwnershipVerification types.List   `tfsdk:"ownership_verification"`
	Certificate           types.List   `tfsdk:"certificate"`
	CreateTime            types.String `tfsdk:"create_time"`
	UpdateTime            types.String `tfsdk:"update_time"`
}

var TeoInferenceDomainsDomainAttrTypes = map[string]attr.Type{
	"domain":                 types.StringType,
	"status":                 types.StringType,
	"auth_switch":            types.StringType,
	"cname":                  types.StringType,
	"ownership_verification": types.ListType{ElemType: types.ObjectType{AttrTypes: TeoInferenceDomainsOwnershipVerificationAttrTypes}},
	"certificate":            types.ListType{ElemType: types.ObjectType{AttrTypes: TeoInferenceDomainsCertificateAttrTypes}},
	"create_time":            types.StringType,
	"update_time":            types.StringType,
}

func (d *TeoInferenceDomainsDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "tencentcloud_teo_inference_domains_1110_01"
}

func (d *TeoInferenceDomainsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this data source to query the list of TEO (EdgeOne) inference service domains.",
		Attributes: map[string]schema.Attribute{
			"zone_id": schema.StringAttribute{
				Required:    true,
				Description: "Site ID.",
			},
			"service_id": schema.StringAttribute{
				Required:    true,
				Description: "Inference service ID.",
			},
			"sort_by": schema.StringAttribute{
				Optional:    true,
				Description: "Sort field. Valid values: CreateTime (domain creation time), UpdateTime (domain modification time). Default value: CreateTime.",
			},
			"sort_order": schema.StringAttribute{
				Optional:    true,
				Description: "Sort order. Valid values: Asc (ascending), Desc (descending). Default value: Desc.",
			},
			"domains": schema.ListNestedAttribute{
				Computed:    true,
				Description: "List of inference service domains.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"domain": schema.StringAttribute{
							Computed:    true,
							Description: "Domain name.",
						},
						"status": schema.StringAttribute{
							Computed:    true,
							Description: "Domain status. Valid values: Online (effective), Process (deploying), Offline (deactivated), Init (not effective, waiting for site activation), Forbidden (domain disabled).",
						},
						"auth_switch": schema.StringAttribute{
							Computed:    true,
							Description: "Authentication switch of inference task requests. Valid values: Off (disabled), On (enabled). Default value: On.",
						},
						"cname": schema.StringAttribute{
							Computed:    true,
							Description: "CNAME address. To verify the CNAME configuration status of the domain, refer to the CheckCnameStatus API.",
						},
						"ownership_verification": schema.ListNestedAttribute{
							Computed:    true,
							Description: "Information required for ownership verification when the domain needs ownership verification to continue serving. Refer to the VerifyOwnership API for ownership verification.",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"dns_verification": schema.ListNestedAttribute{
										Computed:    true,
										Description: "Information required when using DNS resolution verification for CNAME or no-domain access. Note: this field may return null, indicating that no valid value can be obtained.",
										NestedObject: schema.NestedAttributeObject{
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
									},
									"file_verification": schema.ListNestedAttribute{
										Computed:    true,
										Description: "Information required when using file verification for CNAME or no-domain access. Note: this field may return null, indicating that no valid value can be obtained.",
										NestedObject: schema.NestedAttributeObject{
											Attributes: map[string]schema.Attribute{
												"path": schema.StringAttribute{
													Computed:    true,
													Description: "The URL path part of the verification file. EdgeOne fetches it via http://{Host}{URL Path}.",
												},
												"content": schema.StringAttribute{
													Computed:    true,
													Description: "Content of the verification file, which needs to be written into the txt file returned by the path field.",
												},
											},
										},
									},
									"ns_verification": schema.ListNestedAttribute{
										Computed:    true,
										Description: "Information required for switching DNS servers during NS access. Note: this field may return null, indicating that no valid value can be obtained.",
										NestedObject: schema.NestedAttributeObject{
											Attributes: map[string]schema.Attribute{
												"name_servers": schema.SetAttribute{
													Computed:    true,
													ElementType: types.StringType,
													Description: "DNS server addresses assigned to the user during NS access; the domain NameServer needs to be switched to these addresses.",
												},
											},
										},
									},
								},
							},
						},
						"certificate": schema.ListNestedAttribute{
							Computed:    true,
							Description: "Domain certificate information. Refer to ApplyFreeCertificate, CheckFreeCertificateVerification and ModifyHostsCertificate APIs.",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"host": schema.StringAttribute{
										Computed:    true,
										Description: "Domain. Note: this field may return null, indicating that no valid value can be obtained.",
									},
									"mode": schema.StringAttribute{
										Computed:    true,
										Description: "Certificate configuration mode. Valid values: disable (no certificate), eofreecert (EdgeOne free certificate), sslcert (SSL certificate).",
									},
									"host_cert_info": schema.ListNestedAttribute{
										Computed:    true,
										Description: "Server certificate configuration. Note: this field may return null, indicating that no valid value can be obtained.",
										NestedObject: schema.NestedAttributeObject{
											Attributes: map[string]schema.Attribute{
												"cert_id": schema.StringAttribute{
													Computed:    true,
													Description: "Server certificate ID. Note: this field may return null, indicating that no valid value can be obtained.",
												},
												"alias": schema.StringAttribute{
													Computed:    true,
													Description: "Certificate alias. Note: this field may return null, indicating that no valid value can be obtained.",
												},
												"type": schema.StringAttribute{
													Computed:    true,
													Description: "Certificate type. Valid values: default (default certificate), upload (user uploaded), managed (Tencent Cloud managed). Note: this field may return null, indicating that no valid value can be obtained.",
												},
												"expire_time": schema.StringAttribute{
													Computed:    true,
													Description: "Certificate expiration time. Note: this field may return null, indicating that no valid value can be obtained.",
												},
												"deploy_time": schema.StringAttribute{
													Computed:    true,
													Description: "Certificate deployment time. Note: this field may return null, indicating that no valid value can be obtained.",
												},
												"sign_algo": schema.StringAttribute{
													Computed:    true,
													Description: "Signature algorithm. Note: this field may return null, indicating that no valid value can be obtained.",
												},
												"status": schema.StringAttribute{
													Computed:    true,
													Description: "Certificate status. Valid values: deployed (deployed), processing (deploying), applying (applying), failed (application failed), issued (binding failed). Note: this field may return null, indicating that no valid value can be obtained.",
												},
											},
										},
									},
									"apply_type": schema.StringAttribute{
										Computed:    true,
										Description: "Application type. Valid values: apply (managed by EdgeOne), none (not managed by EdgeOne). Note: this field may return null, indicating that no valid value can be obtained.",
									},
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

// buildDnsVerificationList converts a DnsVerification pointer into a single-element
// list (or a null list when the API returns nil).
func (d *TeoInferenceDomainsDataSource) buildDnsVerificationList(ctx context.Context, dns *teov20220901.DnsVerification) types.List {
	if dns == nil {
		return types.ListNull(types.ObjectType{AttrTypes: TeoInferenceDomainsDnsVerificationAttrTypes})
	}
	model := TeoInferenceDomainsDnsVerificationModel{
		Subdomain:   types.StringPointerValue(dns.Subdomain),
		RecordType:  types.StringPointerValue(dns.RecordType),
		RecordValue: types.StringPointerValue(dns.RecordValue),
	}
	list, diags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: TeoInferenceDomainsDnsVerificationAttrTypes}, []TeoInferenceDomainsDnsVerificationModel{model})
	if diags.HasError() {
		return types.ListNull(types.ObjectType{AttrTypes: TeoInferenceDomainsDnsVerificationAttrTypes})
	}
	return list
}

// buildFileVerificationList converts a FileVerification pointer into a single-element
// list (or a null list when the API returns nil).
func (d *TeoInferenceDomainsDataSource) buildFileVerificationList(ctx context.Context, file *teov20220901.FileVerification) types.List {
	if file == nil {
		return types.ListNull(types.ObjectType{AttrTypes: TeoInferenceDomainsFileVerificationAttrTypes})
	}
	model := TeoInferenceDomainsFileVerificationModel{
		Path:    types.StringPointerValue(file.Path),
		Content: types.StringPointerValue(file.Content),
	}
	list, diags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: TeoInferenceDomainsFileVerificationAttrTypes}, []TeoInferenceDomainsFileVerificationModel{model})
	if diags.HasError() {
		return types.ListNull(types.ObjectType{AttrTypes: TeoInferenceDomainsFileVerificationAttrTypes})
	}
	return list
}

// buildNsVerificationList converts a NsVerification pointer into a single-element
// list (or a null list when the API returns nil).
func (d *TeoInferenceDomainsDataSource) buildNsVerificationList(ctx context.Context, ns *teov20220901.NsVerification) types.List {
	if ns == nil {
		return types.ListNull(types.ObjectType{AttrTypes: TeoInferenceDomainsNsVerificationAttrTypes})
	}
	nameServers := types.SetNull(types.StringType)
	if ns.NameServers != nil {
		values := make([]attr.Value, 0, len(ns.NameServers))
		for _, v := range ns.NameServers {
			values = append(values, types.StringPointerValue(v))
		}
		set, diags := types.SetValue(types.StringType, values)
		if diags.HasError() {
			return types.ListNull(types.ObjectType{AttrTypes: TeoInferenceDomainsNsVerificationAttrTypes})
		}
		nameServers = set
	}
	model := TeoInferenceDomainsNsVerificationModel{NameServers: nameServers}
	list, diags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: TeoInferenceDomainsNsVerificationAttrTypes}, []TeoInferenceDomainsNsVerificationModel{model})
	if diags.HasError() {
		return types.ListNull(types.ObjectType{AttrTypes: TeoInferenceDomainsNsVerificationAttrTypes})
	}
	return list
}

// buildOwnershipVerificationList converts an OwnershipVerification pointer into a
// single-element list (or a null list when the API returns nil).
func (d *TeoInferenceDomainsDataSource) buildOwnershipVerificationList(ctx context.Context, ov *teov20220901.OwnershipVerification) types.List {
	if ov == nil {
		return types.ListNull(types.ObjectType{AttrTypes: TeoInferenceDomainsOwnershipVerificationAttrTypes})
	}
	model := TeoInferenceDomainsOwnershipVerificationModel{
		DnsVerification:  d.buildDnsVerificationList(ctx, ov.DnsVerification),
		FileVerification: d.buildFileVerificationList(ctx, ov.FileVerification),
		NsVerification:   d.buildNsVerificationList(ctx, ov.NsVerification),
	}
	list, diags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: TeoInferenceDomainsOwnershipVerificationAttrTypes}, []TeoInferenceDomainsOwnershipVerificationModel{model})
	if diags.HasError() {
		return types.ListNull(types.ObjectType{AttrTypes: TeoInferenceDomainsOwnershipVerificationAttrTypes})
	}
	return list
}

// buildHostCertInfoList converts HostCertInfo slice into a list (or a null list
// when the API returns nil).
func (d *TeoInferenceDomainsDataSource) buildHostCertInfoList(ctx context.Context, certInfos []*teov20220901.HostCertInfo) types.List {
	if certInfos == nil {
		return types.ListNull(types.ObjectType{AttrTypes: TeoInferenceDomainsHostCertInfoAttrTypes})
	}
	models := make([]TeoInferenceDomainsHostCertInfoModel, 0, len(certInfos))
	for _, info := range certInfos {
		if info == nil {
			continue
		}
		models = append(models, TeoInferenceDomainsHostCertInfoModel{
			CertId:     types.StringPointerValue(info.CertId),
			Alias:      types.StringPointerValue(info.Alias),
			Type:       types.StringPointerValue(info.Type),
			ExpireTime: types.StringPointerValue(info.ExpireTime),
			DeployTime: types.StringPointerValue(info.DeployTime),
			SignAlgo:   types.StringPointerValue(info.SignAlgo),
			Status:     types.StringPointerValue(info.Status),
		})
	}
	list, diags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: TeoInferenceDomainsHostCertInfoAttrTypes}, models)
	if diags.HasError() {
		return types.ListNull(types.ObjectType{AttrTypes: TeoInferenceDomainsHostCertInfoAttrTypes})
	}
	return list
}

// buildCertificateList converts a HostsCertificate pointer into a single-element
// list (or a null list when the API returns nil).
func (d *TeoInferenceDomainsDataSource) buildCertificateList(ctx context.Context, cert *teov20220901.HostsCertificate) types.List {
	if cert == nil {
		return types.ListNull(types.ObjectType{AttrTypes: TeoInferenceDomainsCertificateAttrTypes})
	}
	model := TeoInferenceDomainsCertificateModel{
		Host:         types.StringPointerValue(cert.Host),
		Mode:         types.StringPointerValue(cert.Mode),
		HostCertInfo: d.buildHostCertInfoList(ctx, cert.HostCertInfo),
		ApplyType:    types.StringPointerValue(cert.ApplyType),
	}
	list, diags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: TeoInferenceDomainsCertificateAttrTypes}, []TeoInferenceDomainsCertificateModel{model})
	if diags.HasError() {
		return types.ListNull(types.ObjectType{AttrTypes: TeoInferenceDomainsCertificateAttrTypes})
	}
	return list
}

// Read queries the DescribeInferenceDomains API and writes the result into resp.State.
func (d *TeoInferenceDomainsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config TeoInferenceDomainsDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	request := teov20220901.NewDescribeInferenceDomainsRequest()
	request.ZoneId = config.ZoneId.ValueStringPointer()
	request.ServiceId = config.ServiceId.ValueStringPointer()
	if !config.SortBy.IsNull() && !config.SortBy.IsUnknown() {
		request.SortBy = config.SortBy.ValueStringPointer()
	}
	if !config.SortOrder.IsNull() && !config.SortOrder.IsUnknown() {
		request.SortOrder = config.SortOrder.ValueStringPointer()
	}

	teoClient := d.Client().UseTeoClient()
	response, err := teoClient.DescribeInferenceDomainsWithContext(ctx, request)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error querying TEO inference service domains",
			fmt.Sprintf("DescribeInferenceDomains failed: %v", err),
		)
		return
	}
	if response == nil || response.Response == nil {
		resp.Diagnostics.AddError(
			"Error querying TEO inference service domains",
			"DescribeInferenceDomains returned an empty response",
		)
		return
	}

	// An empty Domains list is a genuinely empty result and is returned as an
	// empty list; a nil response is reported as an error above.
	domains := make([]TeoInferenceDomainsDomainModel, 0, len(response.Response.Domains))
	for _, domain := range response.Response.Domains {
		if domain == nil {
			continue
		}
		domains = append(domains, TeoInferenceDomainsDomainModel{
			Domain:                types.StringPointerValue(domain.Domain),
			Status:                types.StringPointerValue(domain.Status),
			AuthSwitch:            types.StringPointerValue(domain.AuthSwitch),
			Cname:                 types.StringPointerValue(domain.Cname),
			OwnershipVerification: d.buildOwnershipVerificationList(ctx, domain.OwnershipVerification),
			Certificate:           d.buildCertificateList(ctx, domain.Certificate),
			CreateTime:            types.StringPointerValue(domain.CreateTime),
			UpdateTime:            types.StringPointerValue(domain.UpdateTime),
		})
	}

	domainsList, diags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: TeoInferenceDomainsDomainAttrTypes}, domains)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	config.Domains = domainsList
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}
