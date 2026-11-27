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

// AlbListenerCertificateItemModel is the element type of certificates.
type AlbListenerCertificateItemModel struct {
	AssociatedTime  types.String `tfsdk:"associated_time"`
	CertificateId   types.String `tfsdk:"certificate_id"`
	CertificateType types.String `tfsdk:"certificate_type"`
	IsDefault       types.Bool   `tfsdk:"is_default"`
	Status          types.String `tfsdk:"status"`
}

// AlbListenerCertificateItemAttrTypes is the object type of the certificates elements.
var AlbListenerCertificateItemAttrTypes = map[string]attr.Type{
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
		Description: "Query the list of certificates bound to a listener.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "ID of the data source, a random token.",
			},
			"certificate_type": schema.StringAttribute{
				Required:    true,
				Description: "Certificate type. Value: CA or SVR (server certificate).",
			},
			"listener_id": schema.StringAttribute{
				Required:    true,
				Description: "Listener ID, in the format of lst- followed by 8 alphanumeric characters.",
			},
			"load_balancer_id": schema.StringAttribute{
				Required:    true,
				Description: "Load balancer instance ID, in the format of alb- followed by 8 alphanumeric characters.",
			},
			"certificates": schema.ListNestedAttribute{
				Computed:    true,
				Description: "List of certificates bound to the listener.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"associated_time": schema.StringAttribute{
							Computed:    true,
							Description: "Certificate binding time.",
						},
						"certificate_id": schema.StringAttribute{
							Computed:    true,
							Description: "Certificate ID.",
						},
						"certificate_type": schema.StringAttribute{
							Computed:    true,
							Description: "Certificate type. Value: CA or SVR (server certificate).",
						},
						"is_default": schema.BoolAttribute{
							Computed:    true,
							Description: "Whether it is the default certificate of the listener. Value: true: default certificate; false: extended certificate.",
						},
						"status": schema.StringAttribute{
							Computed:    true,
							Description: "Certificate-listener binding status. Value: Associated (associated), Associating (associating), Disassociating (disassociating), Error (error).",
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
	certificates, err := service.DescribeListenerCertificates(ctx, describeListenerCertificatesRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error querying ALB listener certificates",
			fmt.Sprintf("DescribeListenerCertificates failed: %v", err),
		)

		return
	}

	items := make([]AlbListenerCertificateItemModel, 0, len(certificates))
	for _, item := range certificates {
		if item == nil {
			continue
		}

		items = append(items, AlbListenerCertificateItemModel{
			AssociatedTime:  types.StringPointerValue(item.AssociatedTime),
			CertificateId:   types.StringPointerValue(item.CertificateId),
			CertificateType: types.StringPointerValue(item.CertificateType),
			IsDefault:       types.BoolPointerValue(item.IsDefault),
			Status:          types.StringPointerValue(item.Status),
		})
	}

	list, diags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: AlbListenerCertificateItemAttrTypes}, items)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	config.Certificates = list
	config.Id = types.StringValue(fwHelper.BuildToken())

	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}
