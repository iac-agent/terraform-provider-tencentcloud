package ags

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	v20250920 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ags/v20250920"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/fw"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

// AGS_SANDBOX_TOOL_STATE_POLL_INTERVAL is the delay between two consecutive
// DescribeSandboxToolList calls while waiting for a sandbox tool to reach a
// stable status. It is a package level variable so that unit tests can shrink it.
var AGS_SANDBOX_TOOL_STATE_POLL_INTERVAL = 5 * time.Second

// _ resource.Resource is intentionally omitted: the factory function
// NewAgsSandboxToolResource returns resource.Resource, which already performs
// the implicit compile-time check at the return statement.
var (
	_ resource.ResourceWithConfigure   = &AgsSandboxToolResource{}
	_ resource.ResourceWithImportState = &AgsSandboxToolResource{}
)

// NewAgsSandboxToolResource is the factory referenced by
// tencentcloud/framework/registry.go to register this resource.
func NewAgsSandboxToolResource() resource.Resource {
	return &AgsSandboxToolResource{}
}

// AgsSandboxToolResource implements resource.Resource for
// tencentcloud_ags_sandbox_tool.
type AgsSandboxToolResource struct {
	fw.ResourceWithConfigure
	fw.WithImportByID
}

// AgsSandboxToolResourceModel maps the schema attributes.
type AgsSandboxToolResourceModel struct {
	Id                    types.String `tfsdk:"id"`
	ToolId                types.String `tfsdk:"tool_id"`
	ToolName              types.String `tfsdk:"tool_name"`
	ToolType              types.String `tfsdk:"tool_type"`
	NetworkConfiguration  types.List   `tfsdk:"network_configuration"`
	Description           types.String `tfsdk:"description"`
	DefaultTimeout        types.String `tfsdk:"default_timeout"`
	Tags                  types.List   `tfsdk:"tags"`
	ClientToken           types.String `tfsdk:"client_token"`
	RoleArn               types.String `tfsdk:"role_arn"`
	StorageMounts         types.List   `tfsdk:"storage_mounts"`
	CustomConfiguration   types.List   `tfsdk:"custom_configuration"`
	ComputerConfiguration types.List   `tfsdk:"computer_configuration"`
	LogConfiguration      types.List   `tfsdk:"log_configuration"`
	Persistent            types.Bool   `tfsdk:"persistent"`
	Status                types.String `tfsdk:"status"`
	StatusReason          types.String `tfsdk:"status_reason"`
	CreateTime            types.String `tfsdk:"create_time"`
	UpdateTime            types.String `tfsdk:"update_time"`
}

// Nested models. Block attributes are modelled with their own structs so that
// ElementsAs / ListValueFrom can convert them directly.

type agsSandboxToolVpcConfigModel struct {
	SubnetIds        types.Set `tfsdk:"subnet_ids"`
	SecurityGroupIds types.Set `tfsdk:"security_group_ids"`
}

type agsSandboxToolNetworkConfigurationModel struct {
	NetworkMode types.String                   `tfsdk:"network_mode"`
	VpcConfig   []agsSandboxToolVpcConfigModel `tfsdk:"vpc_config"`
}

type agsSandboxToolTagModel struct {
	Key   types.String `tfsdk:"key"`
	Value types.String `tfsdk:"value"`
}

type agsSandboxToolCosModel struct {
	Endpoint   types.String `tfsdk:"endpoint"`
	BucketName types.String `tfsdk:"bucket_name"`
	BucketPath types.String `tfsdk:"bucket_path"`
}

type agsSandboxToolImageModel struct {
	Reference         types.String `tfsdk:"reference"`
	ImageRegistryType types.String `tfsdk:"image_registry_type"`
	SubPath           types.String `tfsdk:"sub_path"`
}

type agsSandboxToolCfsModel struct {
	FileSystemId types.String `tfsdk:"file_system_id"`
	Path         types.String `tfsdk:"path"`
}

type agsSandboxToolAgentBucketModel struct {
	LibraryId    types.String `tfsdk:"library_id"`
	SpaceId      types.String `tfsdk:"space_id"`
	AccessDomain types.String `tfsdk:"access_domain"`
}

type agsSandboxToolStorageSourceModel struct {
	Cos         []agsSandboxToolCosModel         `tfsdk:"cos"`
	Image       []agsSandboxToolImageModel       `tfsdk:"image"`
	Cfs         []agsSandboxToolCfsModel         `tfsdk:"cfs"`
	AgentBucket []agsSandboxToolAgentBucketModel `tfsdk:"agent_bucket"`
}

type agsSandboxToolVolumeModel struct {
	VolumeName types.String `tfsdk:"volume_name"`
	VolumeId   types.String `tfsdk:"volume_id"`
}

type agsSandboxToolVolumeTemplateModel struct {
	VolumeTemplateName types.String `tfsdk:"volume_template_name"`
	ReuseKey           types.String `tfsdk:"reuse_key"`
}

type agsSandboxToolStorageMountModel struct {
	Name           types.String                        `tfsdk:"name"`
	StorageSource  []agsSandboxToolStorageSourceModel  `tfsdk:"storage_source"`
	MountPath      types.String                        `tfsdk:"mount_path"`
	ReadOnly       types.Bool                          `tfsdk:"read_only"`
	Volume         []agsSandboxToolVolumeModel         `tfsdk:"volume"`
	VolumeTemplate []agsSandboxToolVolumeTemplateModel `tfsdk:"volume_template"`
}

type agsSandboxToolEnvModel struct {
	Name  types.String `tfsdk:"name"`
	Value types.String `tfsdk:"value"`
}

type agsSandboxToolPortModel struct {
	Name     types.String `tfsdk:"name"`
	Port     types.Int64  `tfsdk:"port"`
	Protocol types.String `tfsdk:"protocol"`
}

type agsSandboxToolResourcesModel struct {
	Cpu     types.String `tfsdk:"cpu"`
	Memory  types.String `tfsdk:"memory"`
	Storage types.String `tfsdk:"storage"`
}

type agsSandboxToolHttpGetModel struct {
	Path   types.String `tfsdk:"path"`
	Port   types.Int64  `tfsdk:"port"`
	Scheme types.String `tfsdk:"scheme"`
}

type agsSandboxToolProbeModel struct {
	HttpGet          []agsSandboxToolHttpGetModel `tfsdk:"http_get"`
	ReadyTimeoutMs   types.Int64                  `tfsdk:"ready_timeout_ms"`
	ProbeTimeoutMs   types.Int64                  `tfsdk:"probe_timeout_ms"`
	ProbePeriodMs    types.Int64                  `tfsdk:"probe_period_ms"`
	SuccessThreshold types.Int64                  `tfsdk:"success_threshold"`
	FailureThreshold types.Int64                  `tfsdk:"failure_threshold"`
}

type agsSandboxToolDnsConfigModel struct {
	Servers  types.Set `tfsdk:"servers"`
	Searches types.Set `tfsdk:"searches"`
	Options  types.Set `tfsdk:"options"`
}

type agsSandboxToolCapabilitiesModel struct {
	Add  types.Set `tfsdk:"add"`
	Drop types.Set `tfsdk:"drop"`
}

type agsSandboxToolSecurityContextModel struct {
	Privileged   types.Bool                        `tfsdk:"privileged"`
	Capabilities []agsSandboxToolCapabilitiesModel `tfsdk:"capabilities"`
}

type agsSandboxToolCustomConfigurationModel struct {
	Image             types.String                         `tfsdk:"image"`
	ImageRegistryType types.String                         `tfsdk:"image_registry_type"`
	Command           types.Set                            `tfsdk:"command"`
	Args              types.Set                            `tfsdk:"args"`
	Env               []agsSandboxToolEnvModel             `tfsdk:"env"`
	Ports             []agsSandboxToolPortModel            `tfsdk:"ports"`
	Resources         []agsSandboxToolResourcesModel       `tfsdk:"resources"`
	Probe             []agsSandboxToolProbeModel           `tfsdk:"probe"`
	DnsConfig         []agsSandboxToolDnsConfigModel       `tfsdk:"dns_config"`
	SecurityContext   []agsSandboxToolSecurityContextModel `tfsdk:"security_context"`
}

type agsSandboxToolWaaConfigurationModel struct {
	ImageId types.String `tfsdk:"image_id"`
}

type agsSandboxToolOsWorldConfigurationModel struct {
	Version types.String `tfsdk:"version"`
}

type agsSandboxToolComputerConfigurationModel struct {
	WaaConfiguration     []agsSandboxToolWaaConfigurationModel     `tfsdk:"waa_configuration"`
	OsWorldConfiguration []agsSandboxToolOsWorldConfigurationModel `tfsdk:"os_world_configuration"`
}

type agsSandboxToolClsConfigModel struct {
	TopicId types.String `tfsdk:"topic_id"`
}

type agsSandboxToolLogSourcesModel struct {
	Files types.Set `tfsdk:"files"`
}

type agsSandboxToolLogConfigurationModel struct {
	ClsConfig  []agsSandboxToolClsConfigModel  `tfsdk:"cls_config"`
	LogSources []agsSandboxToolLogSourcesModel `tfsdk:"log_sources"`
}

func (r *AgsSandboxToolResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "tencentcloud_ags_sandbox_tool"
}

func (r *AgsSandboxToolResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a resource to manage an AGS (Agent Sandbox) sandbox tool. " +
			"A sandbox tool is the runtime template (tool type, network configuration, storage mounts, custom container configuration) " +
			"used to create AGS sandbox instances. Deleting this resource only removes the sandbox tool definition itself; " +
			"the COS buckets, CFS file systems, agent buckets and volumes referenced by the tool are never deleted.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "ID of the sandbox tool, same as `tool_id`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"tool_id": schema.StringAttribute{
				Computed:    true,
				Description: "ID of the sandbox tool assigned by the cloud side. It is the resource identity.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"tool_name": schema.StringAttribute{
				Required: true,
				Description: "Name of the sandbox tool. Length 1-50 characters, supporting English letters, digits, underscores and hyphens. " +
					"The name must be unique under the same AppId. Changing it forces a new resource.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"tool_type": schema.StringAttribute{
				Required: true,
				Description: "Type of the sandbox tool. Valid values: `browser`, `code-interpreter`, `mobile`, `osworld`, " +
					"`custom`, `swebench`, `aio`, `android-world`, `waa`. Changing it forces a new resource.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"network_configuration": schema.ListNestedAttribute{
				Required:    true,
				Description: "Network configuration of sandbox instances created from this tool. Changing it forces a new resource.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"network_mode": schema.StringAttribute{
							Required:    true,
							Description: "Network mode. Valid values: `PUBLIC`, `VPC`, `SANDBOX`.",
						},
						"vpc_config": schema.ListNestedAttribute{
							Optional:    true,
							Description: "VPC network configuration, used when `network_mode` is `VPC`.",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"subnet_ids": schema.SetAttribute{
										Optional:    true,
										ElementType: types.StringType,
										Description: "List of VPC subnet IDs.",
									},
									"security_group_ids": schema.SetAttribute{
										Optional:    true,
										ElementType: types.StringType,
										Description: "List of security group IDs.",
									},
								},
							},
						},
					},
				},
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "Description of the sandbox tool. Maximum length 200 characters.",
			},
			"default_timeout": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Description: "Default timeout of a sandbox instance, e.g. `5m`, `300s`, `1h`. The system default (5 minutes) is used when omitted, " +
					"maximum 24 hours. The cloud side only returns this as a number of seconds, so the refreshed value is normalised to the `Ns` form. " +
					"Changing it forces a new resource.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
			},
			"tags": schema.ListNestedAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Tags bound to the sandbox tool.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"key": schema.StringAttribute{
							Optional:    true,
							Description: "Tag key.",
						},
						"value": schema.StringAttribute{
							Optional:    true,
							Description: "Tag value.",
						},
					},
				},
			},
			"client_token": schema.StringAttribute{
				Optional: true,
				Description: "Idempotency token of the create request, at most 64 characters. It is only sent on creation and is never read back " +
					"from the cloud side. Changing it forces a new resource.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"role_arn": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "ARN of the role bound to the sandbox tool. Changing it forces a new resource.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
			},
			"storage_mounts": schema.ListNestedAttribute{
				Optional: true,
				Description: "Storage mount configuration of sandbox instances created from this tool. " +
					"`UpdateSandboxTool` has no storage mount parameter, so changing it forces a new resource.",
				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplace(),
				},
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Optional:    true,
							Description: "Name of the storage mount configuration.",
						},
						"storage_source": schema.ListNestedAttribute{
							Optional:    true,
							Description: "Storage source configuration.",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"cos": schema.ListNestedAttribute{
										Optional:    true,
										Description: "COS bucket configuration.",
										NestedObject: schema.NestedAttributeObject{
											Attributes: map[string]schema.Attribute{
												"endpoint": schema.StringAttribute{
													Optional:    true,
													Description: "Access domain of the object storage.",
												},
												"bucket_name": schema.StringAttribute{
													Optional:    true,
													Description: "Name of the object storage bucket.",
												},
												"bucket_path": schema.StringAttribute{
													Optional:    true,
													Description: "Path inside the bucket. It must be an absolute path starting with `/`.",
												},
											},
										},
									},
									"image": schema.ListNestedAttribute{
										Optional:    true,
										Description: "Image volume configuration.",
										NestedObject: schema.NestedAttributeObject{
											Attributes: map[string]schema.Attribute{
												"reference": schema.StringAttribute{
													Optional:    true,
													Description: "Image address.",
												},
												"image_registry_type": schema.StringAttribute{
													Optional:    true,
													Description: "Image registry type. Valid values: `enterprise`, `personal`.",
												},
												"sub_path": schema.StringAttribute{
													Optional:    true,
													Description: "Path inside the image.",
												},
											},
										},
									},
									"cfs": schema.ListNestedAttribute{
										Optional:    true,
										Description: "CFS file system configuration.",
										NestedObject: schema.NestedAttributeObject{
											Attributes: map[string]schema.Attribute{
												"file_system_id": schema.StringAttribute{
													Optional:    true,
													Description: "ID of the CFS file system.",
												},
												"path": schema.StringAttribute{
													Optional:    true,
													Description: "Mount path of the CFS file system.",
												},
											},
										},
									},
									"agent_bucket": schema.ListNestedAttribute{
										Optional:    true,
										Description: "AgentBucket storage configuration.",
										NestedObject: schema.NestedAttributeObject{
											Attributes: map[string]schema.Attribute{
												"library_id": schema.StringAttribute{
													Optional:    true,
													Description: "Library ID passed to the AgentBucket.",
												},
												"space_id": schema.StringAttribute{
													Optional:    true,
													Description: "Space ID passed to the AgentBucket.",
												},
												"access_domain": schema.StringAttribute{
													Optional:    true,
													Description: "Access domain passed to the AgentBucket.",
												},
											},
										},
									},
								},
							},
						},
						"mount_path": schema.StringAttribute{
							Optional:    true,
							Description: "Local mount path inside the sandbox instance.",
						},
						"read_only": schema.BoolAttribute{
							Optional:    true,
							Description: "Whether the storage mount is read-only. Defaults to `false` on the cloud side.",
						},
						"volume": schema.ListNestedAttribute{
							Optional:    true,
							Description: "Reference of a volume.",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"volume_name": schema.StringAttribute{
										Optional:    true,
										Description: "Volume name, resolved to a stable volume id by the cloud side under the current AppId.",
									},
									"volume_id": schema.StringAttribute{
										Optional:    true,
										Description: "Stable volume id, in the form of `vol-xxxxxxxx`.",
									},
								},
							},
						},
						"volume_template": schema.ListNestedAttribute{
							Optional:    true,
							Description: "Reference of a volume template together with its derivation rules.",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"volume_template_name": schema.StringAttribute{
										Optional:    true,
										Description: "Name of the volume template.",
									},
									"reuse_key": schema.StringAttribute{
										Optional:    true,
										Description: "Reuse key. It can be a literal value or `${metadata.xxx}`.",
									},
								},
							},
						},
					},
				},
			},
			"custom_configuration": schema.ListNestedAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Custom configuration of the sandbox tool, used when `tool_type` is `custom`.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"image": schema.StringAttribute{
							Optional:    true,
							Description: "Image address.",
						},
						"image_registry_type": schema.StringAttribute{
							Optional:    true,
							Description: "Image registry type. Valid values: `enterprise`, `personal`, `custom`.",
						},
						"command": schema.SetAttribute{
							Optional:    true,
							ElementType: types.StringType,
							Description: "Startup command.",
						},
						"args": schema.SetAttribute{
							Optional:    true,
							ElementType: types.StringType,
							Description: "Startup arguments.",
						},
						"env": schema.ListNestedAttribute{
							Optional:    true,
							Description: "Environment variables.",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"name": schema.StringAttribute{
										Optional:    true,
										Description: "Name of the environment variable.",
									},
									"value": schema.StringAttribute{
										Optional:    true,
										Description: "Value of the environment variable.",
									},
								},
							},
						},
						"ports": schema.ListNestedAttribute{
							Optional:    true,
							Description: "Port configuration.",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"name": schema.StringAttribute{
										Optional:    true,
										Description: "Port name.",
									},
									"port": schema.Int64Attribute{
										Optional:    true,
										Description: "Port number.",
									},
									"protocol": schema.StringAttribute{
										Optional:    true,
										Description: "Protocol.",
									},
								},
							},
						},
						"resources": schema.ListNestedAttribute{
							Optional:    true,
							Description: "Resource configuration.",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"cpu": schema.StringAttribute{
										Optional:    true,
										Description: "CPU resource amount.",
									},
									"memory": schema.StringAttribute{
										Optional:    true,
										Description: "Memory resource amount.",
									},
									"storage": schema.StringAttribute{
										Optional:    true,
										Description: "Custom disk size. Valid values: `1Gi`, `5Gi`, `10Gi`, `20Gi`.",
									},
								},
							},
						},
						"probe": schema.ListNestedAttribute{
							Optional:    true,
							Description: "Health probe configuration.",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"http_get": schema.ListNestedAttribute{
										Optional:    true,
										Description: "HTTP GET probe configuration.",
										NestedObject: schema.NestedAttributeObject{
											Attributes: map[string]schema.Attribute{
												"path": schema.StringAttribute{
													Optional:    true,
													Description: "Request path.",
												},
												"port": schema.Int64Attribute{
													Optional:    true,
													Description: "Request port.",
												},
												"scheme": schema.StringAttribute{
													Optional:    true,
													Description: "Request scheme.",
												},
											},
										},
									},
									"ready_timeout_ms": schema.Int64Attribute{
										Optional:    true,
										Description: "Ready timeout of the health check, in milliseconds.",
									},
									"probe_timeout_ms": schema.Int64Attribute{
										Optional:    true,
										Description: "Timeout of a single probe, in milliseconds.",
									},
									"probe_period_ms": schema.Int64Attribute{
										Optional:    true,
										Description: "Interval between two probes, in milliseconds.",
									},
									"success_threshold": schema.Int64Attribute{
										Optional:    true,
										Description: "Number of consecutive successes required to be considered healthy.",
									},
									"failure_threshold": schema.Int64Attribute{
										Optional:    true,
										Description: "Number of consecutive failures required to be considered unhealthy.",
									},
								},
							},
						},
						"dns_config": schema.ListNestedAttribute{
							Optional:    true,
							Description: "DNS configuration of the sandbox.",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"servers": schema.SetAttribute{
										Optional:    true,
										ElementType: types.StringType,
										Description: "DNS server addresses. Valid IP addresses are required. The cloud side defaults to `10.0.0.1`.",
									},
									"searches": schema.SetAttribute{
										Optional:    true,
										ElementType: types.StringType,
										Description: "Search domains, mapping to the `search` directive of resolv.conf.",
									},
									"options": schema.SetAttribute{
										Optional:    true,
										ElementType: types.StringType,
										Description: "Options, mapping to the `options` directive of resolv.conf.",
									},
								},
							},
						},
						"security_context": schema.ListNestedAttribute{
							Optional:    true,
							Description: "Security context of the sandbox.",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"privileged": schema.BoolAttribute{
										Optional:    true,
										Description: "Whether to enable privileged mode. The cloud side defaults to `true`.",
									},
									"capabilities": schema.ListNestedAttribute{
										Optional:    true,
										Description: "POSIX capabilities that can be added to or removed from the running container.",
										NestedObject: schema.NestedAttributeObject{
											Attributes: map[string]schema.Attribute{
												"add": schema.SetAttribute{
													Optional:    true,
													ElementType: types.StringType,
													Description: "Capabilities to add, e.g. `SYS_ADMIN`.",
												},
												"drop": schema.SetAttribute{
													Optional:    true,
													ElementType: types.StringType,
													Description: "Capabilities to remove, e.g. `ALL`.",
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			"computer_configuration": schema.ListNestedAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Desktop computer environment configuration of the sandbox tool.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"waa_configuration": schema.ListNestedAttribute{
							Optional:    true,
							Description: "WAA sandbox tool configuration.",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"image_id": schema.StringAttribute{
										Optional:    true,
										Description: "ID of the custom WAA image.",
									},
								},
							},
						},
						"os_world_configuration": schema.ListNestedAttribute{
							Optional:    true,
							Description: "Built-in OSWorld configuration.",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"version": schema.StringAttribute{
										Optional: true,
										Description: "Built-in OSWorld version. Valid values: `osworld1`, `osworld2`. " +
											"The cloud side defaults to `osworld1`.",
									},
								},
							},
						},
					},
				},
			},
			"log_configuration": schema.ListNestedAttribute{
				Optional: true,
				Description: "Log shipping configuration of the sandbox tool. " +
					"`UpdateSandboxTool` has no log configuration parameter, so changing it forces a new resource.",
				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplace(),
				},
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"cls_config": schema.ListNestedAttribute{
							Optional:    true,
							Description: "CLS log shipping configuration.",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"topic_id": schema.StringAttribute{
										Optional:    true,
										Description: "ID of the CLS log topic used for shipping sandbox tool logs.",
									},
								},
							},
						},
						"log_sources": schema.ListNestedAttribute{
							Optional:    true,
							Description: "Log source configuration.",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"files": schema.SetAttribute{
										Optional:    true,
										ElementType: types.StringType,
										Description: "Log file paths to collect. They must be files under `/logs/`, sub directories are not supported, at most 10 files.",
									},
								},
							},
						},
					},
				},
			},
			"persistent": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Whether the sandbox is persistent. Changing it forces a new resource.",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
					boolplanmodifier.RequiresReplace(),
				},
			},
			"status": schema.StringAttribute{
				Computed:    true,
				Description: "Status of the sandbox tool. Valid values: `CREATING`, `ACTIVE`, `DELETING`, `FAILED`.",
			},
			"status_reason": schema.StringAttribute{
				Computed:    true,
				Description: "Reason describing why the sandbox tool is in its current status.",
			},
			"create_time": schema.StringAttribute{
				Computed:    true,
				Description: "Creation time of the sandbox tool, in ISO8601 format.",
			},
			"update_time": schema.StringAttribute{
				Computed:    true,
				Description: "Last update time of the sandbox tool, in ISO8601 format.",
			},
		},
	}
}

// Create handles `terraform apply` for new sandbox tools.
func (r *AgsSandboxToolResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan AgsSandboxToolResourceModel
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

	request := v20250920.NewCreateSandboxToolRequest()
	request.ToolName = plan.ToolName.ValueStringPointer()
	request.ToolType = plan.ToolType.ValueStringPointer()

	networkConfiguration, diags := expandNetworkConfiguration(ctx, plan.NetworkConfiguration)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	request.NetworkConfiguration = networkConfiguration

	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		request.Description = plan.Description.ValueStringPointer()
	}
	if !plan.DefaultTimeout.IsNull() && !plan.DefaultTimeout.IsUnknown() {
		request.DefaultTimeout = plan.DefaultTimeout.ValueStringPointer()
	}
	if !plan.ClientToken.IsNull() && !plan.ClientToken.IsUnknown() {
		request.ClientToken = plan.ClientToken.ValueStringPointer()
	}
	if !plan.RoleArn.IsNull() && !plan.RoleArn.IsUnknown() {
		request.RoleArn = plan.RoleArn.ValueStringPointer()
	}
	if !plan.Persistent.IsNull() && !plan.Persistent.IsUnknown() {
		request.Persistent = plan.Persistent.ValueBoolPointer()
	}

	tags, diags := expandTags(ctx, plan.Tags)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	request.Tags = tags

	storageMounts, diags := expandStorageMounts(ctx, plan.StorageMounts)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	request.StorageMounts = storageMounts

	customConfiguration, diags := expandCustomConfiguration(ctx, plan.CustomConfiguration)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	request.CustomConfiguration = customConfiguration

	computerConfiguration, diags := expandComputerConfiguration(ctx, plan.ComputerConfiguration)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	request.ComputerConfiguration = computerConfiguration

	logConfiguration, diags := expandLogConfiguration(ctx, plan.LogConfiguration)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	request.LogConfiguration = logConfiguration

	service := NewAgsService(r.Client())
	toolId, err := service.CreateSandboxTool(ctx, request)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating AGS sandbox tool",
			fmt.Sprintf("CreateSandboxTool failed: %v", err),
		)
		return
	}

	// The tool is created asynchronously (status CREATING). Persist the identity
	// BEFORE waiting so that a failed wait can still be cleaned up by
	// `terraform destroy`.
	plan.Id = types.StringValue(toolId)
	plan.ToolId = types.StringValue(toolId)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	item, err := r.waitForToolStatus(ctx, toolId, true)
	if err != nil {
		resp.Diagnostics.AddError("Error waiting for AGS sandbox tool to become active", err.Error())
		return
	}

	resp.Diagnostics.Append(r.refreshModel(ctx, &plan, item)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Read handles `terraform refresh` and post-apply state refresh.
func (r *AgsSandboxToolResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state AgsSandboxToolResourceModel
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

	toolId := state.ToolId.ValueString()
	if toolId == "" {
		toolId = state.Id.ValueString()
	}

	service := NewAgsService(r.Client())
	item, err := service.DescribeSandboxToolById(ctx, toolId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading AGS sandbox tool",
			fmt.Sprintf("DescribeSandboxToolList failed for tool %s: %v", toolId, err),
		)
		return
	}

	// A confirmed miss on the cloud side removes the resource from state so that
	// Terraform plans a recreation. A nil response is an API failure and has
	// already been turned into an error above.
	if item == nil {
		resp.Diagnostics.AddWarning(
			"Sandbox tool not found",
			fmt.Sprintf("AGS sandbox tool %s not found, removing it from state.", toolId),
		)
		resp.State.RemoveResource(ctx)
		return
	}

	if toolId != state.Id.ValueString() {
		state.Id = types.StringValue(toolId)
	}
	resp.Diagnostics.Append(r.refreshModel(ctx, &state, item)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update handles in-place updates.
func (r *AgsSandboxToolResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan AgsSandboxToolResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state AgsSandboxToolResourceModel
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

	toolId := state.ToolId.ValueString()

	request := v20250920.NewUpdateSandboxToolRequest()
	request.ToolId = helper.String(toolId)

	// UpdateSandboxTool does not document whether an omitted argument means
	// "keep the previous value" or "clear it". Only attributes whose new value
	// is known are sent, and a removal is reported instead of being guessed.
	if plan.Description.IsUnknown() {
		// nothing to send
	} else if plan.Description.IsNull() {
		if !state.Description.IsNull() {
			resp.Diagnostics.AddError(
				"Clear semantics of description is not confirmed",
				"UpdateSandboxTool does not document whether an omitted Description keeps or clears the previous value, "+
					"so removing this attribute from the configuration is refused instead of guessing a direction.",
			)
			return
		}
	} else {
		request.Description = plan.Description.ValueStringPointer()
	}

	if !plan.Tags.IsNull() && !plan.Tags.IsUnknown() {
		tags, diags := expandTags(ctx, plan.Tags)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		request.Tags = tags
	}

	if !plan.NetworkConfiguration.IsNull() && !plan.NetworkConfiguration.IsUnknown() {
		networkConfiguration, diags := expandNetworkConfiguration(ctx, plan.NetworkConfiguration)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		request.NetworkConfiguration = networkConfiguration
	}

	if !plan.CustomConfiguration.IsNull() && !plan.CustomConfiguration.IsUnknown() {
		customConfiguration, diags := expandCustomConfiguration(ctx, plan.CustomConfiguration)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		request.CustomConfiguration = customConfiguration
	}

	if !plan.ComputerConfiguration.IsNull() && !plan.ComputerConfiguration.IsUnknown() {
		computerConfiguration, diags := expandComputerConfiguration(ctx, plan.ComputerConfiguration)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		request.ComputerConfiguration = computerConfiguration
	}

	service := NewAgsService(r.Client())
	if err := service.UpdateSandboxTool(ctx, request); err != nil {
		resp.Diagnostics.AddError(
			"Error updating AGS sandbox tool",
			fmt.Sprintf("UpdateSandboxTool failed for tool %s: %v", toolId, err),
		)
		return
	}

	item, err := r.waitForToolStatus(ctx, toolId, true)
	if err != nil {
		resp.Diagnostics.AddError("Error waiting for AGS sandbox tool to become active after update", err.Error())
		return
	}

	plan.Id = state.Id
	plan.ToolId = state.ToolId
	resp.Diagnostics.Append(r.refreshModel(ctx, &plan, item)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete handles sandbox tool deletion.
func (r *AgsSandboxToolResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state AgsSandboxToolResourceModel
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

	toolId := state.ToolId.ValueString()
	service := NewAgsService(r.Client())

	if err := service.DeleteSandboxTool(ctx, toolId); err != nil {
		// ResourceNotFound(.SandboxTool) means the tool is already gone, which is
		// the desired end state of a delete.
		if !isAgsSandboxToolNotFound(err) {
			resp.Diagnostics.AddError(
				"Error deleting AGS sandbox tool",
				fmt.Sprintf("DeleteSandboxTool failed for tool %s: %v", toolId, err),
			)
			return
		}
		// nothing left to wait for
		return
	}

	if _, err := r.waitForToolStatus(ctx, toolId, false); err != nil {
		resp.Diagnostics.AddError("Error waiting for AGS sandbox tool to be deleted", err.Error())
		return
	}
}

// ImportState is embedded via fw.WithImportByID: `terraform import` passes the
// tool id directly and the following Read fills the remaining attributes.

// waitForToolStatus polls DescribeSandboxToolList until the sandbox tool
// reaches ACTIVE (wantActive) or disappears (wantActive == false).
func (r *AgsSandboxToolResource) waitForToolStatus(ctx context.Context, toolId string, wantActive bool) (*v20250920.SandboxTool, error) {
	service := NewAgsService(r.Client())
	deadline := time.Now().Add(tccommon.WriteRetryTimeout)

	for {
		item, err := service.DescribeSandboxToolById(ctx, toolId)
		if err != nil {
			return nil, err
		}

		if item == nil {
			if !wantActive {
				return nil, nil
			}
			return nil, fmt.Errorf("sandbox tool %s disappeared while waiting for status %s", toolId, AGS_SANDBOX_TOOL_STATUS_ACTIVE)
		}

		status := ""
		if item.Status != nil {
			status = *item.Status
		}

		if wantActive && status == AGS_SANDBOX_TOOL_STATUS_ACTIVE {
			return item, nil
		}

		if status == AGS_SANDBOX_TOOL_STATUS_FAILED {
			reason := ""
			if item.StatusReason != nil {
				reason = *item.StatusReason
			}
			return nil, fmt.Errorf("sandbox tool %s reached status %s: %s", toolId, AGS_SANDBOX_TOOL_STATUS_FAILED, reason)
		}

		if time.Now().After(deadline) {
			return nil, fmt.Errorf("timed out waiting for sandbox tool %s, last observed status: %s", toolId, status)
		}

		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("context cancelled while waiting for sandbox tool %s, last observed status: %s", toolId, status)
		case <-time.After(AGS_SANDBOX_TOOL_STATE_POLL_INTERVAL):
		}
	}
}

// isAgsSandboxToolNotFound reports whether the error means the tool is already gone.
func isAgsSandboxToolNotFound(err error) bool {
	return tccommon.IsExpectError(err, []string{
		"ResourceNotFound",
		"ResourceNotFound.SandboxTool",
	})
}

// refreshModel copies cloud owned attributes into the model.
//
// The cloud side only returns a part of what was sent on creation: it never
// echoes the raw duration string of `default_timeout`, and it may normalise
// (or echo differently) the values nested inside the blocks. Blocks whose
// attributes are user owned therefore keep the configured value whenever the
// configuration already carries one; the cloud value is only used to fill in
// what the configuration cannot have (i.e. right after an import).
func (r *AgsSandboxToolResource) refreshModel(ctx context.Context, m *AgsSandboxToolResourceModel, item *v20250920.SandboxTool) diag.Diagnostics {
	var diags diag.Diagnostics
	if item == nil {
		return diags
	}

	m.ToolId = types.StringPointerValue(item.ToolId)
	m.Id = types.StringPointerValue(item.ToolId)

	// tool_name and tool_type are described differently by the create and the
	// describe documentation, so the configured value wins whenever it is known.
	if m.ToolName.IsNull() || m.ToolName.IsUnknown() || m.ToolName.ValueString() == "" {
		m.ToolName = types.StringPointerValue(item.ToolName)
	}
	if m.ToolType.IsNull() || m.ToolType.IsUnknown() || m.ToolType.ValueString() == "" {
		m.ToolType = types.StringPointerValue(item.ToolType)
	}

	// Cloud owned outputs.
	m.Status = types.StringPointerValue(item.Status)
	m.StatusReason = types.StringPointerValue(item.StatusReason)
	m.CreateTime = types.StringPointerValue(item.CreateTime)
	m.UpdateTime = types.StringPointerValue(item.UpdateTime)
	m.Persistent = types.BoolPointerValue(item.Persistent)

	// Attributes that can only be set on creation: the configuration is the only
	// thing that can express them (the cloud side either never returns them or
	// returns them in a different shape), so they are only filled in when state
	// does not carry a value yet.
	if isAgsEmptyString(m.Description) {
		m.Description = types.StringPointerValue(item.Description)
	}
	if isAgsEmptyString(m.RoleArn) {
		m.RoleArn = types.StringPointerValue(item.RoleArn)
	}
	if isAgsEmptyString(m.ClientToken) {
		// ClientToken is never returned by the describe API, so it stays null
		// after an import and is never part of the drift detection.
		m.ClientToken = types.StringNull()
	}
	if isAgsEmptyString(m.DefaultTimeout) {
		if item.DefaultTimeoutSeconds != nil {
			m.DefaultTimeout = types.StringValue(fmt.Sprintf("%ds", *item.DefaultTimeoutSeconds))
		} else {
			m.DefaultTimeout = types.StringNull()
		}
	}

	if m.NetworkConfiguration.IsNull() || m.NetworkConfiguration.IsUnknown() {
		networkConfiguration, d := flattenNetworkConfiguration(ctx, item.NetworkConfiguration)
		diags.Append(d...)
		if diags.HasError() {
			return diags
		}
		m.NetworkConfiguration = networkConfiguration
	}

	// Optional + Computed blocks: the configured value wins whenever it is known,
	// because the cloud side may echo it in a normalised form (which would show
	// up as a permanent diff). The cloud value is only used when the block is not
	// managed by the configuration, i.e. right after an import or after the user
	// removed the block from HCL.
	if m.Tags.IsNull() || m.Tags.IsUnknown() {
		tags, d := flattenTags(ctx, item.Tags)
		diags.Append(d...)
		if diags.HasError() {
			return diags
		}
		m.Tags = tags
	}

	if m.StorageMounts.IsNull() || m.StorageMounts.IsUnknown() {
		storageMounts, d := flattenStorageMounts(ctx, item.StorageMounts)
		diags.Append(d...)
		if diags.HasError() {
			return diags
		}
		m.StorageMounts = storageMounts
	}

	if m.CustomConfiguration.IsNull() || m.CustomConfiguration.IsUnknown() {
		customConfiguration, d := flattenCustomConfiguration(ctx, item.CustomConfiguration)
		diags.Append(d...)
		if diags.HasError() {
			return diags
		}
		m.CustomConfiguration = customConfiguration
	}

	if m.ComputerConfiguration.IsNull() || m.ComputerConfiguration.IsUnknown() {
		computerConfiguration, d := flattenComputerConfiguration(ctx, item.ComputerConfiguration)
		diags.Append(d...)
		if diags.HasError() {
			return diags
		}
		m.ComputerConfiguration = computerConfiguration
	}

	if m.LogConfiguration.IsNull() || m.LogConfiguration.IsUnknown() {
		logConfiguration, d := flattenLogConfiguration(ctx, item.LogConfiguration)
		diags.Append(d...)
		if diags.HasError() {
			return diags
		}
		m.LogConfiguration = logConfiguration
	}

	// The sandbox tool cannot be created without a network configuration, so a
	// null value here means the configuration was lost; report it instead of
	// writing an incomplete state.
	if m.NetworkConfiguration.IsNull() {
		diags.AddError(
			"Ags sandbox tool network configuration is missing",
			fmt.Sprintf("sandbox tool %s has no network configuration in state and the cloud side returned none", m.Id.ValueString()),
		)
	}

	return diags
}

// isAgsEmptyString reports whether the value is null, unknown or an empty string.
func isAgsEmptyString(v types.String) bool {
	return v.IsNull() || v.IsUnknown() || v.ValueString() == ""
}

// objectType builds a types.ObjectType from the given attribute types.
func objectType(attrTypes map[string]attr.Type) types.ObjectType {
	return types.ObjectType{AttrTypes: attrTypes}
}

// stringSliceFromSet converts a types.Set of strings into []*string.
func stringSliceFromSet(ctx context.Context, set types.Set, diags *diag.Diagnostics) []*string {
	if set.IsNull() || set.IsUnknown() {
		return nil
	}

	var values []string
	diags.Append(set.ElementsAs(ctx, &values, false)...)
	if diags.HasError() || len(values) == 0 {
		return nil
	}
	return helper.Strings(values)
}

// stringSetFromSlice converts []*string into a types.Set of strings.
func stringSetFromSlice(ctx context.Context, values []*string, diags *diag.Diagnostics) types.Set {
	if values == nil {
		return types.SetNull(types.StringType)
	}

	elems := make([]attr.Value, 0, len(values))
	for _, v := range values {
		if v == nil {
			continue
		}
		elems = append(elems, types.StringValue(*v))
	}

	set, d := types.SetValue(types.StringType, elems)
	diags.Append(d...)
	return set
}

// expandNetworkConfiguration converts the Terraform model into the API struct.
func expandNetworkConfiguration(ctx context.Context, list types.List) (*v20250920.NetworkConfiguration, diag.Diagnostics) {
	var diags diag.Diagnostics
	if list.IsNull() || list.IsUnknown() {
		return nil, diags
	}

	var items []agsSandboxToolNetworkConfigurationModel
	diags.Append(list.ElementsAs(ctx, &items, false)...)
	if diags.HasError() || len(items) == 0 {
		return nil, diags
	}

	item := items[0]
	out := &v20250920.NetworkConfiguration{
		NetworkMode: item.NetworkMode.ValueStringPointer(),
	}
	if len(item.VpcConfig) > 0 {
		vpc := item.VpcConfig[0]
		out.VpcConfig = &v20250920.VPCConfig{
			SubnetIds:        stringSliceFromSet(ctx, vpc.SubnetIds, &diags),
			SecurityGroupIds: stringSliceFromSet(ctx, vpc.SecurityGroupIds, &diags),
		}
	}
	return out, diags
}

// flattenNetworkConfiguration converts the API struct into the Terraform model.
func flattenNetworkConfiguration(ctx context.Context, in *v20250920.NetworkConfiguration) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics

	vpcConfigTypes := map[string]attr.Type{
		"subnet_ids":         types.SetType{ElemType: types.StringType},
		"security_group_ids": types.SetType{ElemType: types.StringType},
	}
	networkConfigurationTypes := map[string]attr.Type{
		"network_mode": types.StringType,
		"vpc_config":   types.ListType{ElemType: objectType(vpcConfigTypes)},
	}
	listType := objectType(networkConfigurationTypes)

	if in == nil {
		return types.ListNull(listType), diags
	}

	item := agsSandboxToolNetworkConfigurationModel{
		NetworkMode: types.StringPointerValue(in.NetworkMode),
		VpcConfig:   []agsSandboxToolVpcConfigModel{},
	}
	if in.VpcConfig != nil {
		item.VpcConfig = append(item.VpcConfig, agsSandboxToolVpcConfigModel{
			SubnetIds:        stringSetFromSlice(ctx, in.VpcConfig.SubnetIds, &diags),
			SecurityGroupIds: stringSetFromSlice(ctx, in.VpcConfig.SecurityGroupIds, &diags),
		})
	}

	list, d := types.ListValueFrom(ctx, listType, []agsSandboxToolNetworkConfigurationModel{item})
	diags.Append(d...)
	return list, diags
}

// expandTags converts the Terraform model into the API struct.
func expandTags(ctx context.Context, list types.List) ([]*v20250920.Tag, diag.Diagnostics) {
	var diags diag.Diagnostics
	if list.IsNull() || list.IsUnknown() {
		return nil, diags
	}

	var items []agsSandboxToolTagModel
	diags.Append(list.ElementsAs(ctx, &items, false)...)
	if diags.HasError() || len(items) == 0 {
		return nil, diags
	}

	out := make([]*v20250920.Tag, 0, len(items))
	for _, item := range items {
		out = append(out, &v20250920.Tag{
			Key:   item.Key.ValueStringPointer(),
			Value: item.Value.ValueStringPointer(),
		})
	}
	return out, diags
}

// flattenTags converts the API struct into the Terraform model.
func flattenTags(ctx context.Context, in []*v20250920.Tag) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics

	tagTypes := map[string]attr.Type{
		"key":   types.StringType,
		"value": types.StringType,
	}
	listType := objectType(tagTypes)

	if in == nil {
		return types.ListNull(listType), diags
	}

	items := make([]agsSandboxToolTagModel, 0, len(in))
	for _, tag := range in {
		if tag == nil {
			continue
		}
		items = append(items, agsSandboxToolTagModel{
			Key:   types.StringPointerValue(tag.Key),
			Value: types.StringPointerValue(tag.Value),
		})
	}

	list, d := types.ListValueFrom(ctx, listType, items)
	diags.Append(d...)
	return list, diags
}

// expandStorageMounts converts the Terraform model into the API struct.
func expandStorageMounts(ctx context.Context, list types.List) ([]*v20250920.StorageMount, diag.Diagnostics) {
	var diags diag.Diagnostics
	if list.IsNull() || list.IsUnknown() {
		return nil, diags
	}

	var items []agsSandboxToolStorageMountModel
	diags.Append(list.ElementsAs(ctx, &items, false)...)
	if diags.HasError() || len(items) == 0 {
		return nil, diags
	}

	out := make([]*v20250920.StorageMount, 0, len(items))
	for _, item := range items {
		mount := &v20250920.StorageMount{
			Name:      item.Name.ValueStringPointer(),
			MountPath: item.MountPath.ValueStringPointer(),
		}
		if !item.ReadOnly.IsNull() && !item.ReadOnly.IsUnknown() {
			mount.ReadOnly = item.ReadOnly.ValueBoolPointer()
		}
		if len(item.StorageSource) > 0 {
			source := item.StorageSource[0]
			storageSource := &v20250920.StorageSource{}
			if len(source.Cos) > 0 {
				cos := source.Cos[0]
				storageSource.Cos = &v20250920.CosStorageSource{
					Endpoint:   cos.Endpoint.ValueStringPointer(),
					BucketName: cos.BucketName.ValueStringPointer(),
					BucketPath: cos.BucketPath.ValueStringPointer(),
				}
			}
			if len(source.Image) > 0 {
				image := source.Image[0]
				storageSource.Image = &v20250920.ImageStorageSource{
					Reference:         image.Reference.ValueStringPointer(),
					ImageRegistryType: image.ImageRegistryType.ValueStringPointer(),
					SubPath:           image.SubPath.ValueStringPointer(),
				}
			}
			if len(source.Cfs) > 0 {
				cfs := source.Cfs[0]
				storageSource.Cfs = &v20250920.CfsStorageSource{
					FileSystemId: cfs.FileSystemId.ValueStringPointer(),
					Path:         cfs.Path.ValueStringPointer(),
				}
			}
			if len(source.AgentBucket) > 0 {
				bucket := source.AgentBucket[0]
				storageSource.AgentBucket = &v20250920.AgentBucketStorageSource{
					LibraryId:    bucket.LibraryId.ValueStringPointer(),
					SpaceId:      bucket.SpaceId.ValueStringPointer(),
					AccessDomain: bucket.AccessDomain.ValueStringPointer(),
				}
			}
			mount.StorageSource = storageSource
		}
		if len(item.Volume) > 0 {
			volume := item.Volume[0]
			mount.Volume = &v20250920.VolumeRef{
				VolumeName: volume.VolumeName.ValueStringPointer(),
				VolumeId:   volume.VolumeId.ValueStringPointer(),
			}
		}
		if len(item.VolumeTemplate) > 0 {
			template := item.VolumeTemplate[0]
			mount.VolumeTemplate = &v20250920.VolumeTemplateRef{
				VolumeTemplateName: template.VolumeTemplateName.ValueStringPointer(),
				ReuseKey:           template.ReuseKey.ValueStringPointer(),
			}
		}
		out = append(out, mount)
	}
	return out, diags
}

// flattenStorageMounts converts the API struct into the Terraform model.
func flattenStorageMounts(ctx context.Context, in []*v20250920.StorageMount) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics

	cosTypes := map[string]attr.Type{
		"endpoint":    types.StringType,
		"bucket_name": types.StringType,
		"bucket_path": types.StringType,
	}
	imageTypes := map[string]attr.Type{
		"reference":           types.StringType,
		"image_registry_type": types.StringType,
		"sub_path":            types.StringType,
	}
	cfsTypes := map[string]attr.Type{
		"file_system_id": types.StringType,
		"path":           types.StringType,
	}
	agentBucketTypes := map[string]attr.Type{
		"library_id":    types.StringType,
		"space_id":      types.StringType,
		"access_domain": types.StringType,
	}
	storageSourceTypes := map[string]attr.Type{
		"cos":          types.ListType{ElemType: objectType(cosTypes)},
		"image":        types.ListType{ElemType: objectType(imageTypes)},
		"cfs":          types.ListType{ElemType: objectType(cfsTypes)},
		"agent_bucket": types.ListType{ElemType: objectType(agentBucketTypes)},
	}
	volumeTypes := map[string]attr.Type{
		"volume_name": types.StringType,
		"volume_id":   types.StringType,
	}
	volumeTemplateTypes := map[string]attr.Type{
		"volume_template_name": types.StringType,
		"reuse_key":            types.StringType,
	}
	storageMountTypes := map[string]attr.Type{
		"name":            types.StringType,
		"storage_source":  types.ListType{ElemType: objectType(storageSourceTypes)},
		"mount_path":      types.StringType,
		"read_only":       types.BoolType,
		"volume":          types.ListType{ElemType: objectType(volumeTypes)},
		"volume_template": types.ListType{ElemType: objectType(volumeTemplateTypes)},
	}
	listType := objectType(storageMountTypes)

	if in == nil {
		return types.ListNull(listType), diags
	}

	items := make([]agsSandboxToolStorageMountModel, 0, len(in))
	for _, mount := range in {
		if mount == nil {
			continue
		}
		item := agsSandboxToolStorageMountModel{
			Name:           types.StringPointerValue(mount.Name),
			MountPath:      types.StringPointerValue(mount.MountPath),
			ReadOnly:       types.BoolPointerValue(mount.ReadOnly),
			StorageSource:  []agsSandboxToolStorageSourceModel{},
			Volume:         []agsSandboxToolVolumeModel{},
			VolumeTemplate: []agsSandboxToolVolumeTemplateModel{},
		}
		if mount.StorageSource != nil {
			source := agsSandboxToolStorageSourceModel{}
			if mount.StorageSource.Cos != nil {
				source.Cos = []agsSandboxToolCosModel{{
					Endpoint:   types.StringPointerValue(mount.StorageSource.Cos.Endpoint),
					BucketName: types.StringPointerValue(mount.StorageSource.Cos.BucketName),
					BucketPath: types.StringPointerValue(mount.StorageSource.Cos.BucketPath),
				}}
			}
			if mount.StorageSource.Image != nil {
				source.Image = []agsSandboxToolImageModel{{
					Reference:         types.StringPointerValue(mount.StorageSource.Image.Reference),
					ImageRegistryType: types.StringPointerValue(mount.StorageSource.Image.ImageRegistryType),
					SubPath:           types.StringPointerValue(mount.StorageSource.Image.SubPath),
				}}
			}
			if mount.StorageSource.Cfs != nil {
				source.Cfs = []agsSandboxToolCfsModel{{
					FileSystemId: types.StringPointerValue(mount.StorageSource.Cfs.FileSystemId),
					Path:         types.StringPointerValue(mount.StorageSource.Cfs.Path),
				}}
			}
			if mount.StorageSource.AgentBucket != nil {
				source.AgentBucket = []agsSandboxToolAgentBucketModel{{
					LibraryId:    types.StringPointerValue(mount.StorageSource.AgentBucket.LibraryId),
					SpaceId:      types.StringPointerValue(mount.StorageSource.AgentBucket.SpaceId),
					AccessDomain: types.StringPointerValue(mount.StorageSource.AgentBucket.AccessDomain),
				}}
			}
			item.StorageSource = append(item.StorageSource, source)
		}
		if mount.Volume != nil {
			item.Volume = append(item.Volume, agsSandboxToolVolumeModel{
				VolumeName: types.StringPointerValue(mount.Volume.VolumeName),
				VolumeId:   types.StringPointerValue(mount.Volume.VolumeId),
			})
		}
		if mount.VolumeTemplate != nil {
			item.VolumeTemplate = append(item.VolumeTemplate, agsSandboxToolVolumeTemplateModel{
				VolumeTemplateName: types.StringPointerValue(mount.VolumeTemplate.VolumeTemplateName),
				ReuseKey:           types.StringPointerValue(mount.VolumeTemplate.ReuseKey),
			})
		}
		items = append(items, item)
	}

	list, d := types.ListValueFrom(ctx, listType, items)
	diags.Append(d...)
	return list, diags
}

// expandCustomConfiguration converts the Terraform model into the API struct.
func expandCustomConfiguration(ctx context.Context, list types.List) (*v20250920.CustomConfiguration, diag.Diagnostics) {
	var diags diag.Diagnostics
	if list.IsNull() || list.IsUnknown() {
		return nil, diags
	}

	var items []agsSandboxToolCustomConfigurationModel
	diags.Append(list.ElementsAs(ctx, &items, false)...)
	if diags.HasError() || len(items) == 0 {
		return nil, diags
	}

	item := items[0]
	out := &v20250920.CustomConfiguration{
		Image:             item.Image.ValueStringPointer(),
		ImageRegistryType: item.ImageRegistryType.ValueStringPointer(),
		Command:           stringSliceFromSet(ctx, item.Command, &diags),
		Args:              stringSliceFromSet(ctx, item.Args, &diags),
	}

	for _, env := range item.Env {
		out.Env = append(out.Env, &v20250920.EnvVar{
			Name:  env.Name.ValueStringPointer(),
			Value: env.Value.ValueStringPointer(),
		})
	}

	for _, port := range item.Ports {
		out.Ports = append(out.Ports, &v20250920.PortConfiguration{
			Name:     port.Name.ValueStringPointer(),
			Port:     port.Port.ValueInt64Pointer(),
			Protocol: port.Protocol.ValueStringPointer(),
		})
	}

	if len(item.Resources) > 0 {
		resources := item.Resources[0]
		out.Resources = &v20250920.ResourceConfiguration{
			CPU:     resources.Cpu.ValueStringPointer(),
			Memory:  resources.Memory.ValueStringPointer(),
			Storage: resources.Storage.ValueStringPointer(),
		}
	}

	if len(item.Probe) > 0 {
		probe := item.Probe[0]
		out.Probe = &v20250920.ProbeConfiguration{
			ReadyTimeoutMs:   probe.ReadyTimeoutMs.ValueInt64Pointer(),
			ProbeTimeoutMs:   probe.ProbeTimeoutMs.ValueInt64Pointer(),
			ProbePeriodMs:    probe.ProbePeriodMs.ValueInt64Pointer(),
			SuccessThreshold: probe.SuccessThreshold.ValueInt64Pointer(),
			FailureThreshold: probe.FailureThreshold.ValueInt64Pointer(),
		}
		if len(probe.HttpGet) > 0 {
			httpGet := probe.HttpGet[0]
			out.Probe.HttpGet = &v20250920.HttpGetAction{
				Path:   httpGet.Path.ValueStringPointer(),
				Port:   httpGet.Port.ValueInt64Pointer(),
				Scheme: httpGet.Scheme.ValueStringPointer(),
			}
		}
	}

	if len(item.DnsConfig) > 0 {
		dnsConfig := item.DnsConfig[0]
		out.DNSConfig = &v20250920.DNSConfig{
			Servers:  stringSliceFromSet(ctx, dnsConfig.Servers, &diags),
			Searches: stringSliceFromSet(ctx, dnsConfig.Searches, &diags),
			Options:  stringSliceFromSet(ctx, dnsConfig.Options, &diags),
		}
	}

	if len(item.SecurityContext) > 0 {
		securityContext := item.SecurityContext[0]
		out.SecurityContext = &v20250920.SecurityContext{}
		if !securityContext.Privileged.IsNull() && !securityContext.Privileged.IsUnknown() {
			out.SecurityContext.Privileged = securityContext.Privileged.ValueBoolPointer()
		}
		if len(securityContext.Capabilities) > 0 {
			capabilities := securityContext.Capabilities[0]
			out.SecurityContext.Capabilities = &v20250920.Capabilities{
				Add:  stringSliceFromSet(ctx, capabilities.Add, &diags),
				Drop: stringSliceFromSet(ctx, capabilities.Drop, &diags),
			}
		}
	}

	return out, diags
}

// flattenCustomConfiguration converts the API struct into the Terraform model.
func flattenCustomConfiguration(ctx context.Context, in *v20250920.CustomConfigurationDetail) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics

	envTypes := map[string]attr.Type{
		"name":  types.StringType,
		"value": types.StringType,
	}
	portTypes := map[string]attr.Type{
		"name":     types.StringType,
		"port":     types.Int64Type,
		"protocol": types.StringType,
	}
	resourcesTypes := map[string]attr.Type{
		"cpu":     types.StringType,
		"memory":  types.StringType,
		"storage": types.StringType,
	}
	httpGetTypes := map[string]attr.Type{
		"path":   types.StringType,
		"port":   types.Int64Type,
		"scheme": types.StringType,
	}
	probeTypes := map[string]attr.Type{
		"http_get":          types.ListType{ElemType: objectType(httpGetTypes)},
		"ready_timeout_ms":  types.Int64Type,
		"probe_timeout_ms":  types.Int64Type,
		"probe_period_ms":   types.Int64Type,
		"success_threshold": types.Int64Type,
		"failure_threshold": types.Int64Type,
	}
	dnsConfigTypes := map[string]attr.Type{
		"servers":  types.SetType{ElemType: types.StringType},
		"searches": types.SetType{ElemType: types.StringType},
		"options":  types.SetType{ElemType: types.StringType},
	}
	capabilitiesTypes := map[string]attr.Type{
		"add":  types.SetType{ElemType: types.StringType},
		"drop": types.SetType{ElemType: types.StringType},
	}
	securityContextTypes := map[string]attr.Type{
		"privileged":   types.BoolType,
		"capabilities": types.ListType{ElemType: objectType(capabilitiesTypes)},
	}
	customConfigurationTypes := map[string]attr.Type{
		"image":               types.StringType,
		"image_registry_type": types.StringType,
		"command":             types.SetType{ElemType: types.StringType},
		"args":                types.SetType{ElemType: types.StringType},
		"env":                 types.ListType{ElemType: objectType(envTypes)},
		"ports":               types.ListType{ElemType: objectType(portTypes)},
		"resources":           types.ListType{ElemType: objectType(resourcesTypes)},
		"probe":               types.ListType{ElemType: objectType(probeTypes)},
		"dns_config":          types.ListType{ElemType: objectType(dnsConfigTypes)},
		"security_context":    types.ListType{ElemType: objectType(securityContextTypes)},
	}
	listType := objectType(customConfigurationTypes)

	if in == nil {
		return types.ListNull(listType), diags
	}

	item := agsSandboxToolCustomConfigurationModel{
		Image:             types.StringPointerValue(in.Image),
		ImageRegistryType: types.StringPointerValue(in.ImageRegistryType),
		Command:           stringSetFromSlice(ctx, in.Command, &diags),
		Args:              stringSetFromSlice(ctx, in.Args, &diags),
		Env:               []agsSandboxToolEnvModel{},
		Ports:             []agsSandboxToolPortModel{},
	}

	for _, env := range in.Env {
		if env == nil {
			continue
		}
		item.Env = append(item.Env, agsSandboxToolEnvModel{
			Name:  types.StringPointerValue(env.Name),
			Value: types.StringPointerValue(env.Value),
		})
	}

	for _, port := range in.Ports {
		if port == nil {
			continue
		}
		item.Ports = append(item.Ports, agsSandboxToolPortModel{
			Name:     types.StringPointerValue(port.Name),
			Port:     types.Int64PointerValue(port.Port),
			Protocol: types.StringPointerValue(port.Protocol),
		})
	}

	if in.Resources != nil {
		item.Resources = []agsSandboxToolResourcesModel{{
			Cpu:     types.StringPointerValue(in.Resources.CPU),
			Memory:  types.StringPointerValue(in.Resources.Memory),
			Storage: types.StringPointerValue(in.Resources.Storage),
		}}
	}

	if in.Probe != nil {
		probe := agsSandboxToolProbeModel{
			ReadyTimeoutMs:   types.Int64PointerValue(in.Probe.ReadyTimeoutMs),
			ProbeTimeoutMs:   types.Int64PointerValue(in.Probe.ProbeTimeoutMs),
			ProbePeriodMs:    types.Int64PointerValue(in.Probe.ProbePeriodMs),
			SuccessThreshold: types.Int64PointerValue(in.Probe.SuccessThreshold),
			FailureThreshold: types.Int64PointerValue(in.Probe.FailureThreshold),
		}
		if in.Probe.HttpGet != nil {
			probe.HttpGet = []agsSandboxToolHttpGetModel{{
				Path:   types.StringPointerValue(in.Probe.HttpGet.Path),
				Port:   types.Int64PointerValue(in.Probe.HttpGet.Port),
				Scheme: types.StringPointerValue(in.Probe.HttpGet.Scheme),
			}}
		}
		item.Probe = []agsSandboxToolProbeModel{probe}
	}

	if in.DNSConfig != nil {
		item.DnsConfig = []agsSandboxToolDnsConfigModel{{
			Servers:  stringSetFromSlice(ctx, in.DNSConfig.Servers, &diags),
			Searches: stringSetFromSlice(ctx, in.DNSConfig.Searches, &diags),
			Options:  stringSetFromSlice(ctx, in.DNSConfig.Options, &diags),
		}}
	}

	if in.SecurityContext != nil {
		securityContext := agsSandboxToolSecurityContextModel{
			Privileged: types.BoolPointerValue(in.SecurityContext.Privileged),
		}
		if in.SecurityContext.Capabilities != nil {
			securityContext.Capabilities = []agsSandboxToolCapabilitiesModel{{
				Add:  stringSetFromSlice(ctx, in.SecurityContext.Capabilities.Add, &diags),
				Drop: stringSetFromSlice(ctx, in.SecurityContext.Capabilities.Drop, &diags),
			}}
		}
		item.SecurityContext = []agsSandboxToolSecurityContextModel{securityContext}
	}

	list, d := types.ListValueFrom(ctx, listType, []agsSandboxToolCustomConfigurationModel{item})
	diags.Append(d...)
	return list, diags
}

// expandComputerConfiguration converts the Terraform model into the API struct.
func expandComputerConfiguration(ctx context.Context, list types.List) (*v20250920.ComputerConfiguration, diag.Diagnostics) {
	var diags diag.Diagnostics
	if list.IsNull() || list.IsUnknown() {
		return nil, diags
	}

	var items []agsSandboxToolComputerConfigurationModel
	diags.Append(list.ElementsAs(ctx, &items, false)...)
	if diags.HasError() || len(items) == 0 {
		return nil, diags
	}

	item := items[0]
	out := &v20250920.ComputerConfiguration{}
	if len(item.WaaConfiguration) > 0 {
		out.WAAConfiguration = &v20250920.WAAConfiguration{
			ImageId: item.WaaConfiguration[0].ImageId.ValueStringPointer(),
		}
	}
	if len(item.OsWorldConfiguration) > 0 {
		out.OSWorldConfiguration = &v20250920.OSWorldConfiguration{
			Version: item.OsWorldConfiguration[0].Version.ValueStringPointer(),
		}
	}
	return out, diags
}

// flattenComputerConfiguration converts the API struct into the Terraform model.
func flattenComputerConfiguration(ctx context.Context, in *v20250920.ComputerConfiguration) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics

	waaTypes := map[string]attr.Type{
		"image_id": types.StringType,
	}
	osWorldTypes := map[string]attr.Type{
		"version": types.StringType,
	}
	computerConfigurationTypes := map[string]attr.Type{
		"waa_configuration":      types.ListType{ElemType: objectType(waaTypes)},
		"os_world_configuration": types.ListType{ElemType: objectType(osWorldTypes)},
	}
	listType := objectType(computerConfigurationTypes)

	if in == nil {
		return types.ListNull(listType), diags
	}

	item := agsSandboxToolComputerConfigurationModel{}
	if in.WAAConfiguration != nil {
		item.WaaConfiguration = []agsSandboxToolWaaConfigurationModel{{
			ImageId: types.StringPointerValue(in.WAAConfiguration.ImageId),
		}}
	}
	if in.OSWorldConfiguration != nil {
		item.OsWorldConfiguration = []agsSandboxToolOsWorldConfigurationModel{{
			Version: types.StringPointerValue(in.OSWorldConfiguration.Version),
		}}
	}

	list, d := types.ListValueFrom(ctx, listType, []agsSandboxToolComputerConfigurationModel{item})
	diags.Append(d...)
	return list, diags
}

// expandLogConfiguration converts the Terraform model into the API struct.
func expandLogConfiguration(ctx context.Context, list types.List) (*v20250920.LogConfiguration, diag.Diagnostics) {
	var diags diag.Diagnostics
	if list.IsNull() || list.IsUnknown() {
		return nil, diags
	}

	var items []agsSandboxToolLogConfigurationModel
	diags.Append(list.ElementsAs(ctx, &items, false)...)
	if diags.HasError() || len(items) == 0 {
		return nil, diags
	}

	item := items[0]
	out := &v20250920.LogConfiguration{}
	if len(item.ClsConfig) > 0 {
		out.CLSConfig = &v20250920.CLSConfig{
			TopicId: item.ClsConfig[0].TopicId.ValueStringPointer(),
		}
	}
	if len(item.LogSources) > 0 {
		out.LogSources = &v20250920.LogSources{
			Files: stringSliceFromSet(ctx, item.LogSources[0].Files, &diags),
		}
	}
	return out, diags
}

// flattenLogConfiguration converts the API struct into the Terraform model.
func flattenLogConfiguration(ctx context.Context, in *v20250920.LogConfiguration) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics

	clsConfigTypes := map[string]attr.Type{
		"topic_id": types.StringType,
	}
	logSourcesTypes := map[string]attr.Type{
		"files": types.SetType{ElemType: types.StringType},
	}
	logConfigurationTypes := map[string]attr.Type{
		"cls_config":  types.ListType{ElemType: objectType(clsConfigTypes)},
		"log_sources": types.ListType{ElemType: objectType(logSourcesTypes)},
	}
	listType := objectType(logConfigurationTypes)

	if in == nil {
		return types.ListNull(listType), diags
	}

	item := agsSandboxToolLogConfigurationModel{}
	if in.CLSConfig != nil {
		item.ClsConfig = []agsSandboxToolClsConfigModel{{
			TopicId: types.StringPointerValue(in.CLSConfig.TopicId),
		}}
	}
	if in.LogSources != nil {
		item.LogSources = []agsSandboxToolLogSourcesModel{{
			Files: stringSetFromSlice(ctx, in.LogSources.Files, &diags),
		}}
	}

	list, d := types.ListValueFrom(ctx, listType, []agsSandboxToolLogConfigurationModel{item})
	diags.Append(d...)
	return list, diags
}
