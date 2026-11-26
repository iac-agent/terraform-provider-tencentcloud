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

var _ datasource.DataSourceWithConfigure = &AlbListenersDataSource{}

// NewAlbListenersDataSource is the factory referenced by
// tencentcloud/framework/registry.go to register this data source.
func NewAlbListenersDataSource() datasource.DataSource {
	return &AlbListenersDataSource{}
}

// AlbListenersDataSource implements datasource.DataSource for
// tencentcloud_alb_listeners.
type AlbListenersDataSource struct {
	fw.DataSourceWithConfigure
}

// AlbListenersDataSourceModel maps the schema attributes.
type AlbListenersDataSourceModel struct {
	Id             types.String `tfsdk:"id"`
	LoadBalancerId types.String `tfsdk:"load_balancer_id"`
	ListenerIds    types.List   `tfsdk:"listener_ids"`
	Filters        types.List   `tfsdk:"filters"`
	Listeners      types.List   `tfsdk:"listeners"`
}

// AlbListenersFilterModel is the element type of filters.
type AlbListenersFilterModel struct {
	Name   types.String `tfsdk:"name"`
	Values types.List   `tfsdk:"values"`
}

// AlbListenersTagItemModel is the element type of listeners.tags.
type AlbListenersTagItemModel struct {
	TagKey   types.String `tfsdk:"tag_key"`
	TagValue types.String `tfsdk:"tag_value"`
}

// AlbListenersTagItemAttrTypes is the object type of the tags elements.
var AlbListenersTagItemAttrTypes = map[string]attr.Type{
	"tag_key":   types.StringType,
	"tag_value": types.StringType,
}

// AlbListenersXForwardedForModel is the model of listeners.x_forwarded_for_config.
type AlbListenersXForwardedForModel struct {
	XForwardedForAlbIdEnabled         types.Bool   `tfsdk:"x_forwarded_for_alb_id_enabled"`
	XForwardedForClientSrcPortEnabled types.Bool   `tfsdk:"x_forwarded_for_client_src_port_enabled"`
	XForwardedForHostEnabled          types.Bool   `tfsdk:"x_forwarded_for_host_enabled"`
	XForwardedForMode                 types.String `tfsdk:"x_forwarded_for_mode"`
	XForwardedForPortEnabled          types.Bool   `tfsdk:"x_forwarded_for_port_enabled"`
	XForwardedForProtoEnabled         types.Bool   `tfsdk:"x_forwarded_for_proto_enabled"`
	XTencentClientIDNEnabled          types.Bool   `tfsdk:"x_tencent_client_idn_enabled"`
	XTencentClientSDNEnabled          types.Bool   `tfsdk:"x_tencent_client_sdn_enabled"`
	XTencentClientSerialEnabled       types.Bool   `tfsdk:"x_tencent_client_serial_enabled"`
	XTencentClientVerifyEnabled       types.Bool   `tfsdk:"x_tencent_client_verify_enabled"`
}

// AlbListenersXForwardedForAttrTypes is the object type of x_forwarded_for_config.
var AlbListenersXForwardedForAttrTypes = map[string]attr.Type{
	"x_forwarded_for_alb_id_enabled":          types.BoolType,
	"x_forwarded_for_client_src_port_enabled": types.BoolType,
	"x_forwarded_for_host_enabled":            types.BoolType,
	"x_forwarded_for_mode":                    types.StringType,
	"x_forwarded_for_port_enabled":            types.BoolType,
	"x_forwarded_for_proto_enabled":           types.BoolType,
	"x_tencent_client_idn_enabled":            types.BoolType,
	"x_tencent_client_sdn_enabled":            types.BoolType,
	"x_tencent_client_serial_enabled":         types.BoolType,
	"x_tencent_client_verify_enabled":         types.BoolType,
}

// AlbListenersItemModel is the element type of listeners.
type AlbListenersItemModel struct {
	CaEnable            types.Bool   `tfsdk:"ca_enable"`
	CreateTime          types.String `tfsdk:"create_time"`
	GzipEnabled         types.Bool   `tfsdk:"gzip_enabled"`
	Http2Enable         types.Bool   `tfsdk:"http2_enable"`
	IdleTimeout         types.Int64  `tfsdk:"idle_timeout"`
	ListenerId          types.String `tfsdk:"listener_id"`
	ListenerName        types.String `tfsdk:"listener_name"`
	ListenerPort        types.Int64  `tfsdk:"listener_port"`
	ListenerProtocol    types.String `tfsdk:"listener_protocol"`
	ListenerStatus      types.String `tfsdk:"listener_status"`
	ModifyTime          types.String `tfsdk:"modify_time"`
	RequestTimeout      types.Int64  `tfsdk:"request_timeout"`
	Tags                types.List   `tfsdk:"tags"`
	TlsSecurityPolicyId types.String `tfsdk:"tls_security_policy_id"`
	XForwardedForConfig types.Object `tfsdk:"x_forwarded_for_config"`
}

// AlbListenersItemAttrTypes is the object type of the listeners elements.
var AlbListenersItemAttrTypes = map[string]attr.Type{
	"ca_enable":              types.BoolType,
	"create_time":            types.StringType,
	"gzip_enabled":           types.BoolType,
	"http2_enable":           types.BoolType,
	"idle_timeout":           types.Int64Type,
	"listener_id":            types.StringType,
	"listener_name":          types.StringType,
	"listener_port":          types.Int64Type,
	"listener_protocol":      types.StringType,
	"listener_status":        types.StringType,
	"modify_time":            types.StringType,
	"request_timeout":        types.Int64Type,
	"tags":                   types.ListType{ElemType: types.ObjectType{AttrTypes: AlbListenersTagItemAttrTypes}},
	"tls_security_policy_id": types.StringType,
	"x_forwarded_for_config": types.ObjectType{AttrTypes: AlbListenersXForwardedForAttrTypes},
}

func (d *AlbListenersDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "tencentcloud_alb_listeners"
}

func (d *AlbListenersDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this data source to query the list of ALB listeners.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "ID of the data source, a random token.",
			},
			// ---- Query inputs ----
			"load_balancer_id": schema.StringAttribute{
				Required:    true,
				Description: "Load balancer instance ID, in the format of alb- followed by 8 alphanumeric characters.",
			},
			"listener_ids": schema.ListAttribute{
				Optional:    true,
				ElementType: types.StringType,
				Description: "List of listener IDs. An ID is in the format of lst- followed by 8 alphanumeric characters.",
			},
			"filters": schema.ListNestedAttribute{
				Optional:    true,
				Description: "Filter conditions, with a maximum of 20. The following fields are supported: Protocol (protocol type), Tags (tags).",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Required:    true,
							Description: "Name of the filter.",
						},
						"values": schema.ListAttribute{
							Required:    true,
							ElementType: types.StringType,
							Description: "Values of the filter.",
						},
					},
				},
			},
			// ---- Query outputs: Computed ----
			"listeners": schema.ListNestedAttribute{
				Computed:    true,
				Description: "Listener information.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"ca_enable": schema.BoolAttribute{
							Computed:    true,
							Description: "Whether mutual authentication is enabled.",
						},
						"create_time": schema.StringAttribute{
							Computed:    true,
							Description: "Creation time of the listener instance, in ISO 8601 format (for example, 2025-01-01T08:30:00+08:00).",
						},
						"gzip_enabled": schema.BoolAttribute{
							Computed:    true,
							Description: "Whether Gzip compression is enabled.",
						},
						"http2_enable": schema.BoolAttribute{
							Computed:    true,
							Description: "Whether HTTP/2 is enabled.",
						},
						"idle_timeout": schema.Int64Attribute{
							Computed:    true,
							Description: "Idle timeout period.",
						},
						"listener_id": schema.StringAttribute{
							Computed:    true,
							Description: "Listener ID, in the format of lst- followed by 8 alphanumeric characters.",
						},
						"listener_name": schema.StringAttribute{
							Computed:    true,
							Description: "Listener name.",
						},
						"listener_port": schema.Int64Attribute{
							Computed:    true,
							Description: "Listener port.",
						},
						"listener_protocol": schema.StringAttribute{
							Computed:    true,
							Description: "Listener protocol.",
						},
						"listener_status": schema.StringAttribute{
							Computed:    true,
							Description: "Listener status. Valid values: Active (running), Provisioning (creating), Configuring (modifying configuration), ProvisionFailed (creation failed).",
						},
						"modify_time": schema.StringAttribute{
							Computed:    true,
							Description: "Last modification time of the listener instance, in ISO 8601 format (for example, 2025-01-01T08:30:00+08:00).",
						},
						"request_timeout": schema.Int64Attribute{
							Computed:    true,
							Description: "Connection request timeout period.",
						},
						"tags": schema.ListNestedAttribute{
							Computed:    true,
							Description: "Tags.",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"tag_key": schema.StringAttribute{
										Computed:    true,
										Description: "Tag key.",
									},
									"tag_value": schema.StringAttribute{
										Computed:    true,
										Description: "Tag value.",
									},
								},
							},
						},
						"tls_security_policy_id": schema.StringAttribute{
							Computed:    true,
							Description: "Security policy ID.",
						},
						"x_forwarded_for_config": schema.SingleNestedAttribute{
							Computed:    true,
							Description: "XForwardedFor configuration.",
							Attributes: map[string]schema.Attribute{
								"x_forwarded_for_alb_id_enabled": schema.BoolAttribute{
									Computed:    true,
									Description: "Whether to obtain the load balancer instance ID through the ALB-ID header field. true: yes; false: no.",
								},
								"x_forwarded_for_client_src_port_enabled": schema.BoolAttribute{
									Computed:    true,
									Description: "Whether to obtain the client port accessing the load balancer through the X-Forwarded-Client-srcport header field. true: yes; false: no.",
								},
								"x_forwarded_for_host_enabled": schema.BoolAttribute{
									Computed:    true,
									Description: "Whether to enable obtaining the domain name accessing the load balancer through the X-Forwarded-Host header field. true: yes; false: no.",
								},
								"x_forwarded_for_mode": schema.StringAttribute{
									Computed:    true,
									Description: "How to handle the X-Forwarded-For (XFF) HTTP header field. append: append mode (default), appending the client real IP address to the end of the X-Forwarded-For header while preserving the original XFF chain information; remove: removal mode, removing the X-Forwarded-For header field and not passing it to the backend server; passthrough: passthrough mode, keeping X-Forwarded-For unchanged and passing it to the backend server without any modification.",
								},
								"x_forwarded_for_port_enabled": schema.BoolAttribute{
									Computed:    true,
									Description: "Whether to obtain the listener port of the load balancer through the X-Forwarded-Port header field. true: yes; false: no.",
								},
								"x_forwarded_for_proto_enabled": schema.BoolAttribute{
									Computed:    true,
									Description: "Whether to obtain the listener protocol of the load balancer through the X-Forwarded-Proto header field. true: yes; false: no.",
								},
								"x_tencent_client_idn_enabled": schema.BoolAttribute{
									Computed:    true,
									Description: "Whether to access the issuer of the client certificate ($ssl_client_i_dn) through the X-Tencent-Client-IDN header. true: yes; false: no.",
								},
								"x_tencent_client_sdn_enabled": schema.BoolAttribute{
									Computed:    true,
									Description: "Whether to access the subject of the client certificate ($ssl_client_s_dn) through the X-Tencent-Client-SDN header. true: yes; false: no.",
								},
								"x_tencent_client_serial_enabled": schema.BoolAttribute{
									Computed:    true,
									Description: "Whether to access the serial number of the client certificate ($ssl_client_serial) through the X-Tencent-Client-Serial header. true: yes; false: no.",
								},
								"x_tencent_client_verify_enabled": schema.BoolAttribute{
									Computed:    true,
									Description: "Whether to access the verification result of the client certificate ($ssl_client_verify) through the X-Tencent-Client-Verify header. true: yes; false: no.",
								},
							},
						},
					},
				},
			},
		},
	}
}

// Read queries the cloud API and writes the result into resp.State.
func (d *AlbListenersDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config AlbListenersDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	describeListenersRequest := albv20251030.NewDescribeListenersRequest()
	describeListenersRequest.LoadBalancerId = config.LoadBalancerId.ValueStringPointer()

	describeListenersRequest.ListenerIds = fwHelper.PointersFromList[string](ctx, config.ListenerIds, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	if !config.Filters.IsNull() && !config.Filters.IsUnknown() {
		var filterModels []AlbListenersFilterModel
		resp.Diagnostics.Append(config.Filters.ElementsAs(ctx, &filterModels, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		filters := make([]*albv20251030.Filter, 0, len(filterModels))
		for _, fm := range filterModels {
			filter := &albv20251030.Filter{
				Name:   fm.Name.ValueStringPointer(),
				Values: fwHelper.PointersFromList[string](ctx, fm.Values, &resp.Diagnostics),
			}
			filters = append(filters, filter)
		}
		if resp.Diagnostics.HasError() {
			return
		}
		describeListenersRequest.Filters = filters
	}

	service := NewAlbService(d.Client())
	listeners, err := service.DescribeListeners(ctx, describeListenersRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error querying ALB listeners",
			fmt.Sprintf("DescribeListeners failed: %v", err),
		)
		return
	}

	// An empty result is a genuinely empty list, not a failure.
	items := make([]AlbListenersItemModel, 0, len(listeners))
	for _, item := range listeners {
		if item == nil {
			continue
		}

		tags := make([]AlbListenersTagItemModel, 0, len(item.Tags))
		for _, tag := range item.Tags {
			if tag == nil {
				continue
			}
			tags = append(tags, AlbListenersTagItemModel{
				TagKey:   fwHelper.StringPointerValueOrNull(tag.TagKey),
				TagValue: fwHelper.StringPointerValueOrNull(tag.TagValue),
			})
		}
		tagsList, diags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: AlbListenersTagItemAttrTypes}, tags)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		xForwardedFor := types.ObjectNull(AlbListenersXForwardedForAttrTypes)
		if item.XForwardedForConfig != nil {
			xForwardedForValue, diags := types.ObjectValueFrom(ctx, AlbListenersXForwardedForAttrTypes, AlbListenersXForwardedForModel{
				XForwardedForAlbIdEnabled:         fwHelper.BoolValueOrNull(item.XForwardedForConfig.XForwardedForAlbIdEnabled),
				XForwardedForClientSrcPortEnabled: fwHelper.BoolValueOrNull(item.XForwardedForConfig.XForwardedForClientSrcPortEnabled),
				XForwardedForHostEnabled:          fwHelper.BoolValueOrNull(item.XForwardedForConfig.XForwardedForHostEnabled),
				XForwardedForMode:                 fwHelper.StringPointerValueOrNull(item.XForwardedForConfig.XForwardedForMode),
				XForwardedForPortEnabled:          fwHelper.BoolValueOrNull(item.XForwardedForConfig.XForwardedForPortEnabled),
				XForwardedForProtoEnabled:         fwHelper.BoolValueOrNull(item.XForwardedForConfig.XForwardedForProtoEnabled),
				XTencentClientIDNEnabled:          fwHelper.BoolValueOrNull(item.XForwardedForConfig.XTencentClientIDNEnabled),
				XTencentClientSDNEnabled:          fwHelper.BoolValueOrNull(item.XForwardedForConfig.XTencentClientSDNEnabled),
				XTencentClientSerialEnabled:       fwHelper.BoolValueOrNull(item.XForwardedForConfig.XTencentClientSerialEnabled),
				XTencentClientVerifyEnabled:       fwHelper.BoolValueOrNull(item.XForwardedForConfig.XTencentClientVerifyEnabled),
			})
			resp.Diagnostics.Append(diags...)
			if resp.Diagnostics.HasError() {
				return
			}
			xForwardedFor = xForwardedForValue
		}

		items = append(items, AlbListenersItemModel{
			CaEnable:            fwHelper.BoolValueOrNull(item.CaEnable),
			CreateTime:          fwHelper.StringPointerValueOrNull(item.CreateTime),
			GzipEnabled:         fwHelper.BoolValueOrNull(item.GzipEnabled),
			Http2Enable:         fwHelper.BoolValueOrNull(item.Http2Enable),
			IdleTimeout:         fwHelper.Int64ValueFromUint(item.IdleTimeout),
			ListenerId:          fwHelper.StringPointerValueOrNull(item.ListenerId),
			ListenerName:        fwHelper.StringPointerValueOrNull(item.ListenerName),
			ListenerPort:        fwHelper.Int64ValueFromUint(item.ListenerPort),
			ListenerProtocol:    fwHelper.StringPointerValueOrNull(item.ListenerProtocol),
			ListenerStatus:      fwHelper.StringPointerValueOrNull(item.ListenerStatus),
			ModifyTime:          fwHelper.StringPointerValueOrNull(item.ModifyTime),
			RequestTimeout:      fwHelper.Int64ValueFromUint(item.RequestTimeout),
			Tags:                tagsList,
			TlsSecurityPolicyId: fwHelper.StringPointerValueOrNull(item.TlsSecurityPolicyId),
			XForwardedForConfig: xForwardedFor,
		})
	}

	list, diags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: AlbListenersItemAttrTypes}, items)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	config.Listeners = list
	config.Id = types.StringValue(fwHelper.BuildToken())

	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}
