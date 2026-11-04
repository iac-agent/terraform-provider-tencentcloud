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

// _ datasource.DataSource is intentionally omitted: the factory function
// NewTeoInferenceDomainsV2DataSource returns datasource.DataSource, which
// already performs the implicit compile-time check at the return statement.
// Only the extended-interface assertion is kept here to verify the Configure
// method signature (promoted from fw.DataSourceWithConfigure) at compile time.
var _ datasource.DataSourceWithConfigure = &TeoInferenceDomainsV2DataSource{}

// NewTeoInferenceDomainsV2DataSource is the factory referenced by
// tencentcloud/framework/registry.go to register this data source.
func NewTeoInferenceDomainsV2DataSource() datasource.DataSource {
	return &TeoInferenceDomainsV2DataSource{}
}

// TeoInferenceDomainsV2DataSource implements datasource.DataSource for
// tencentcloud_teo_inference_domains_v2.
type TeoInferenceDomainsV2DataSource struct {
	fw.DataSourceWithConfigure
}

// TeoInferenceDomainsV2DataSourceModel maps the schema attributes.
type TeoInferenceDomainsV2DataSourceModel struct {
	ZoneId    types.String `tfsdk:"zone_id"`
	ServiceId types.String `tfsdk:"service_id"`
	SortBy    types.String `tfsdk:"sort_by"`
	SortOrder types.String `tfsdk:"sort_order"`
	Domains   types.List   `tfsdk:"domains"`
}

// inferenceDomainItemModel is the element type of domains.
type inferenceDomainItemModel struct {
	Domain                types.String `tfsdk:"domain"`
	Status                types.String `tfsdk:"status"`
	AuthSwitch            types.String `tfsdk:"auth_switch"`
	Cname                 types.String `tfsdk:"cname"`
	OwnershipVerification types.Object `tfsdk:"ownership_verification"`
	Certificate           types.Object `tfsdk:"certificate"`
	CreateTime            types.String `tfsdk:"create_time"`
	UpdateTime            types.String `tfsdk:"update_time"`
}

// ownershipVerificationModel maps the ownership_verification attribute.
type ownershipVerificationModel struct {
	DnsVerification  types.Object `tfsdk:"dns_verification"`
	FileVerification types.Object `tfsdk:"file_verification"`
	NsVerification   types.Object `tfsdk:"ns_verification"`
}

// dnsVerificationModel maps the dns_verification attribute.
type dnsVerificationModel struct {
	Subdomain   types.String `tfsdk:"subdomain"`
	RecordType  types.String `tfsdk:"record_type"`
	RecordValue types.String `tfsdk:"record_value"`
}

// fileVerificationModel maps the file_verification attribute.
type fileVerificationModel struct {
	Path    types.String `tfsdk:"path"`
	Content types.String `tfsdk:"content"`
}

// nsVerificationModel maps the ns_verification attribute.
type nsVerificationModel struct {
	NameServers types.Set `tfsdk:"name_servers"`
}

// certificateModel maps the certificate attribute.
type certificateModel struct {
	Host         types.String `tfsdk:"host"`
	Mode         types.String `tfsdk:"mode"`
	HostCertInfo types.List   `tfsdk:"host_cert_info"`
	ApplyType    types.String `tfsdk:"apply_type"`
}

// hostCertInfoItemModel is the element type of host_cert_info.
type hostCertInfoItemModel struct {
	CertId     types.String `tfsdk:"cert_id"`
	Alias      types.String `tfsdk:"alias"`
	Type       types.String `tfsdk:"type"`
	ExpireTime types.String `tfsdk:"expire_time"`
	DeployTime types.String `tfsdk:"deploy_time"`
	SignAlgo   types.String `tfsdk:"sign_algo"`
	Status     types.String `tfsdk:"status"`
}

// Attribute types for the nested objects written into state.
var (
	inferenceDomainAttrTypes = map[string]attr.Type{
		"domain":                 types.StringType,
		"status":                 types.StringType,
		"auth_switch":            types.StringType,
		"cname":                  types.StringType,
		"ownership_verification": types.ObjectType{AttrTypes: ownershipVerificationAttrTypes},
		"certificate":            types.ObjectType{AttrTypes: certificateAttrTypes},
		"create_time":            types.StringType,
		"update_time":            types.StringType,
	}
	ownershipVerificationAttrTypes = map[string]attr.Type{
		"dns_verification":  types.ObjectType{AttrTypes: dnsVerificationAttrTypes},
		"file_verification": types.ObjectType{AttrTypes: fileVerificationAttrTypes},
		"ns_verification":   types.ObjectType{AttrTypes: nsVerificationAttrTypes},
	}
	dnsVerificationAttrTypes = map[string]attr.Type{
		"subdomain":    types.StringType,
		"record_type":  types.StringType,
		"record_value": types.StringType,
	}
	fileVerificationAttrTypes = map[string]attr.Type{
		"path":    types.StringType,
		"content": types.StringType,
	}
	nsVerificationAttrTypes = map[string]attr.Type{
		"name_servers": types.SetType{ElemType: types.StringType},
	}
	certificateAttrTypes = map[string]attr.Type{
		"host":           types.StringType,
		"mode":           types.StringType,
		"host_cert_info": types.ListType{ElemType: types.ObjectType{AttrTypes: hostCertInfoAttrTypes}},
		"apply_type":     types.StringType,
	}
	hostCertInfoAttrTypes = map[string]attr.Type{
		"cert_id":     types.StringType,
		"alias":       types.StringType,
		"type":        types.StringType,
		"expire_time": types.StringType,
		"deploy_time": types.StringType,
		"sign_algo":   types.StringType,
		"status":      types.StringType,
	}
)

func (d *TeoInferenceDomainsV2DataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "tencentcloud_teo_inference_domains_v2"
}

func (d *TeoInferenceDomainsV2DataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this data source to query the list of TEO inference service domains under a zone and inference service.",
		Attributes: map[string]schema.Attribute{
			// ---- Query inputs: Required / Optional only ----
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
				Description: "Sort field. Valid values: `CreateTime` (domain creation time), `UpdateTime` (domain modification time). Default value: `CreateTime`.",
			},
			"sort_order": schema.StringAttribute{
				Optional:    true,
				Description: "Sort order. Valid values: `Asc` (ascending), `Desc` (descending). Default value: `Desc`.",
			},
			// ---- Query outputs: Computed ----
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
							Description: "Domain status. Valid values: `Online`, `Process`, `Offline`, `Init`, `Forbidden`.",
						},
						"auth_switch": schema.StringAttribute{
							Computed:    true,
							Description: "Inference task request authentication switch. Valid values: `Off`, `On`.",
						},
						"cname": schema.StringAttribute{
							Computed:    true,
							Description: "CNAME address.",
						},
						"ownership_verification": schema.SingleNestedAttribute{
							Computed:    true,
							Description: "Ownership verification information. Only present when the domain needs ownership verification to continue serving.",
							Attributes: map[string]schema.Attribute{
								"dns_verification": schema.SingleNestedAttribute{
									Computed:    true,
									Description: "DNS resolution verification information, used for CNAME and non-domain-name access.",
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
									Description: "File verification information, used for CNAME and non-domain-name access.",
									Attributes: map[string]schema.Attribute{
										"path": schema.StringAttribute{
											Computed:    true,
											Description: "URL path portion of the file verification endpoint.",
										},
										"content": schema.StringAttribute{
											Computed:    true,
											Description: "Content of the verification file.",
										},
									},
								},
								"ns_verification": schema.SingleNestedAttribute{
									Computed:    true,
									Description: "NS verification information, used for NS access to switch DNS servers.",
									Attributes: map[string]schema.Attribute{
										"name_servers": schema.SetAttribute{
											Computed:    true,
											ElementType: types.StringType,
											Description: "DNS server addresses assigned to the user for NS access.",
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
									Description: "Certificate mode. Valid values: `disable`, `eofreecert`, `sslcert`.",
								},
								"host_cert_info": schema.ListNestedAttribute{
									Computed:    true,
									Description: "Server certificate configuration list.",
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
												Description: "Certificate type. Valid values: `default`, `upload`, `managed`.",
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
												Description: "Certificate status. Valid values: `deployed`, `processing`, `applying`, `failed`, `issued`.",
											},
										},
									},
								},
								"apply_type": schema.StringAttribute{
									Computed:    true,
									Description: "Application type. Valid values: `apply`, `none`.",
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

// Read queries the cloud API and writes the result into resp.State.
func (d *TeoInferenceDomainsV2DataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config TeoInferenceDomainsV2DataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if d.Client() == nil {
		resp.Diagnostics.AddError(
			"Provider not configured",
			"The provider client is not available. This is unexpected; please report it.",
		)
		return
	}

	zoneId := config.ZoneId.ValueString()
	serviceId := config.ServiceId.ValueString()
	if zoneId == "" {
		resp.Diagnostics.AddError(
			"Missing zone_id",
			"zone_id is required and cannot be empty.",
		)
		return
	}
	if serviceId == "" {
		resp.Diagnostics.AddError(
			"Missing service_id",
			"service_id is required and cannot be empty.",
		)
		return
	}

	service := NewTeoService(d.Client())
	domains, err := service.DescribeTeoInferenceDomains(ctx, zoneId, serviceId, config.SortBy.ValueString(), config.SortOrder.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error querying TEO inference domains for zone (%s) and service (%s)", zoneId, serviceId),
			err.Error(),
		)
		return
	}

	items := make([]inferenceDomainItemModel, 0, len(domains))
	for _, item := range domains {
		items = append(items, teoInferenceDomainToModel(item))
	}

	domainsList, diags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: inferenceDomainAttrTypes}, items)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	config.Domains = domainsList
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}

// teoInferenceDomainToModel converts a cloud SDK InferenceDomain into the
// list element model. Cloud-nullable sub-objects become null objects, and
// empty collections stay empty (never silently dropped).
func teoInferenceDomainToModel(item *teov20220901.InferenceDomain) inferenceDomainItemModel {
	if item == nil {
		return inferenceDomainItemModel{
			OwnershipVerification: types.ObjectNull(ownershipVerificationAttrTypes),
			Certificate:           types.ObjectNull(certificateAttrTypes),
		}
	}

	var ownershipVerification types.Object
	if item.OwnershipVerification != nil {
		var dnsVerification types.Object
		if item.OwnershipVerification.DnsVerification != nil {
			dnsVerification = dnsVerificationToModel(item.OwnershipVerification.DnsVerification)
		} else {
			dnsVerification = types.ObjectNull(dnsVerificationAttrTypes)
		}

		var fileVerification types.Object
		if item.OwnershipVerification.FileVerification != nil {
			fileVerification = fileVerificationToModel(item.OwnershipVerification.FileVerification)
		} else {
			fileVerification = types.ObjectNull(fileVerificationAttrTypes)
		}

		var nsVerification types.Object
		if item.OwnershipVerification.NsVerification != nil {
			nsVerification = nsVerificationToModel(item.OwnershipVerification.NsVerification)
		} else {
			nsVerification = types.ObjectNull(nsVerificationAttrTypes)
		}

		ownershipVerification = types.ObjectValueMust(ownershipVerificationAttrTypes, map[string]attr.Value{
			"dns_verification":  dnsVerification,
			"file_verification": fileVerification,
			"ns_verification":   nsVerification,
		})
	} else {
		ownershipVerification = types.ObjectNull(ownershipVerificationAttrTypes)
	}

	var certificate types.Object
	if item.Certificate != nil {
		certificate = certificateToModel(item.Certificate)
	} else {
		certificate = types.ObjectNull(certificateAttrTypes)
	}

	return inferenceDomainItemModel{
		Domain:                types.StringPointerValue(item.Domain),
		Status:                types.StringPointerValue(item.Status),
		AuthSwitch:            types.StringPointerValue(item.AuthSwitch),
		Cname:                 types.StringPointerValue(item.Cname),
		OwnershipVerification: ownershipVerification,
		Certificate:           certificate,
		CreateTime:            types.StringPointerValue(item.CreateTime),
		UpdateTime:            types.StringPointerValue(item.UpdateTime),
	}
}

func dnsVerificationToModel(item *teov20220901.DnsVerification) types.Object {
	return types.ObjectValueMust(dnsVerificationAttrTypes, map[string]attr.Value{
		"subdomain":    types.StringPointerValue(item.Subdomain),
		"record_type":  types.StringPointerValue(item.RecordType),
		"record_value": types.StringPointerValue(item.RecordValue),
	})
}

func fileVerificationToModel(item *teov20220901.FileVerification) types.Object {
	return types.ObjectValueMust(fileVerificationAttrTypes, map[string]attr.Value{
		"path":    types.StringPointerValue(item.Path),
		"content": types.StringPointerValue(item.Content),
	})
}

func nsVerificationToModel(item *teov20220901.NsVerification) types.Object {
	nameServers := make([]attr.Value, 0, len(item.NameServers))
	for _, ns := range item.NameServers {
		nameServers = append(nameServers, types.StringPointerValue(ns))
	}
	nameServersSet, diags := types.SetValue(types.StringType, nameServers)
	if diags.HasError() {
		return types.ObjectNull(nsVerificationAttrTypes)
	}
	return types.ObjectValueMust(nsVerificationAttrTypes, map[string]attr.Value{
		"name_servers": nameServersSet,
	})
}

func certificateToModel(item *teov20220901.HostsCertificate) types.Object {
	hostCertInfo := make([]hostCertInfoItemModel, 0, len(item.HostCertInfo))
	for _, cert := range item.HostCertInfo {
		if cert == nil {
			continue
		}
		hostCertInfo = append(hostCertInfo, hostCertInfoItemModel{
			CertId:     types.StringPointerValue(cert.CertId),
			Alias:      types.StringPointerValue(cert.Alias),
			Type:       types.StringPointerValue(cert.Type),
			ExpireTime: types.StringPointerValue(cert.ExpireTime),
			DeployTime: types.StringPointerValue(cert.DeployTime),
			SignAlgo:   types.StringPointerValue(cert.SignAlgo),
			Status:     types.StringPointerValue(cert.Status),
		})
	}

	hostCertInfoList, diags := types.ListValueFrom(context.Background(), types.ObjectType{AttrTypes: hostCertInfoAttrTypes}, hostCertInfo)
	if diags.HasError() {
		return types.ObjectNull(certificateAttrTypes)
	}

	return types.ObjectValueMust(certificateAttrTypes, map[string]attr.Value{
		"host":           types.StringPointerValue(item.Host),
		"mode":           types.StringPointerValue(item.Mode),
		"host_cert_info": hostCertInfoList,
		"apply_type":     types.StringPointerValue(item.ApplyType),
	})
}
