package alb

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	albv20251030 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/alb/v20251030"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/fw"
	fwHelper "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/helper"
)

var _ datasource.DataSourceWithConfigure = &AlbListenerCertificatesDataSource{}

// NewAlbListenerCertificatesDataSource is the factory referenced by
// tencentcloud/framework/registry.go to register this data source.
func NewAlbListenerCertificatesDataSource() datasource.DataSource {
	return &AlbListenerCertificatesDataSource{}
}

// AlbListenerCertificatesDataSource implements datasource.DataSource for
// tencentcloud_alb_listener_certificates.
type AlbListenerCertificatesDataSource struct {
	fw.DataSourceWithConfigure
}

// AlbListenerCertificatesDataSourceModel maps the schema attributes.
type AlbListenerCertificatesDataSourceModel struct {
	Id              types.String `tfsdk:"id"`
	CertificateType types.String `tfsdk:"certificate_type"`
	ListenerId      types.String `tfsdk:"listener_id"`
	LoadBalancerId  types.String `tfsdk:"load_balancer_id"`
	Certificates    types.List   `tfsdk:"certificates"`
}

// AlbListenerCertificatesItemModel is the element type of certificates.
type AlbListenerCertificatesItemModel struct {
	AssociatedTime  types.String `tfsdk:"associated_time"`
	CertificateId   types.String `tfsdk:"certificate_id"`
	CertificateType types.String `tfsdk:"certificate_type"`
	IsDefault       types.Bool   `tfsdk:"is_default"`
	Status          types.String `tfsdk:"status"`
}

// AlbListenerCertificatesItemAttrTypes is the object type of the certificates elements.
var AlbListenerCertificatesItemAttrTypes = map[string]attr.Type{
	"associated_time":  types.StringType,
	"certificate_id":   types.StringType,
	"certificate_type": types.StringType,
	"is_default":       types.BoolType,
	"status":           types.StringType,
}

func (d *AlbListenerCertificatesDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "tencentcloud_alb_listener_certificates"
}

func (d *AlbListenerCertificatesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this data source to query the list of certificates bound to an ALB listener.",
		Attributes: map[string]schema.Attribute{
			// The id carries no cloud meaning — it is a random token
			// generated at read time (see the end of Read).
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "ID of the data source, a random token.",
			},
			// ---- Query inputs ----
			"certificate_type": schema.StringAttribute{
				Required:    true,
				Description: "Certificate type. Values: `CA` or `SVR` (server certificate).",
			},
			"listener_id": schema.StringAttribute{
				Required:    true,
				Description: "ID of the listener, in the format of lst- followed by 8 alphanumeric characters.",
			},
			"load_balancer_id": schema.StringAttribute{
				Required:    true,
				Description: "ID of the ALB instance, in the format of alb- followed by 8 alphanumeric characters.",
			},
			// ---- Query outputs: Computed ----
			"certificates": schema.ListNestedAttribute{
				Computed:    true,
				Description: "List of certificates bound to the listener.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"associated_time": schema.StringAttribute{
							Computed:    true,
							Description: "Time when the certificate was bound.",
						},
						"certificate_id": schema.StringAttribute{
							Computed:    true,
							Description: "ID of the certificate.",
						},
						"certificate_type": schema.StringAttribute{
							Computed:    true,
							Description: "Certificate type. Values: `CA` or `SVR` (server certificate).",
						},
						"is_default": schema.BoolAttribute{
							Computed:    true,
							Description: "Whether the certificate is the default certificate of the listener. Values: `true` (default certificate), `false` (additional certificate).",
						},
						"status": schema.StringAttribute{
							Computed:    true,
							Description: "Binding status of the certificate and the listener. Values: `Associated` (associated), `Associating` (associating), `Disassociating` (disassociating), `Error` (error).",
						},
					},
				},
			},
		},
	}
}

// Read queries the cloud API and writes the result into resp.State.
func (d *AlbListenerCertificatesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config AlbListenerCertificatesDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	describeListenerCertificatesRequest := albv20251030.NewDescribeListenerCertificatesRequest()
	describeListenerCertificatesRequest.CertificateType = config.CertificateType.ValueStringPointer()
	describeListenerCertificatesRequest.ListenerId = config.ListenerId.ValueStringPointer()
	describeListenerCertificatesRequest.LoadBalancerId = config.LoadBalancerId.ValueStringPointer()

	service := NewAlbService(d.Client())

	// The service method internally pages with MaxResults/NextToken and
	// returns the full list — this file never loops by itself and never
	// surfaces the total count.
	certificates, err := service.DescribeListenerCertificates(ctx, describeListenerCertificatesRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error querying ALB listener certificates",
			fmt.Sprintf("DescribeListenerCertificates failed: %v", err),
		)

		return
	}

	// An empty result is a genuinely empty list, not a failure.
	items := make([]AlbListenerCertificatesItemModel, 0, len(certificates))
	for _, item := range certificates {
		if item == nil {
			continue
		}

		items = append(items, AlbListenerCertificatesItemModel{
			AssociatedTime:  types.StringPointerValue(item.AssociatedTime),
			CertificateId:   types.StringPointerValue(item.CertificateId),
			CertificateType: types.StringPointerValue(item.CertificateType),
			IsDefault:       types.BoolPointerValue(item.IsDefault),
			Status:          types.StringPointerValue(item.Status),
		})
	}

	list, diags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: AlbListenerCertificatesItemAttrTypes}, items)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	config.Certificates = list

	// The state id is a random token (JSON resource_id = fwHelper.BuildToken()).
	config.Id = types.StringValue(fwHelper.BuildToken())

	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}
