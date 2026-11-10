package ags

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	agssdk "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ags/v20250920"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/fw"
	fwhelper "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

// Compile-time interface assertions.
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
// tencentcloud_ags_sandbox_tool. The cloud addresses the tool with a single
// ToolId returned by CreateSandboxTool, which is used as the Terraform id.
type AgsSandboxToolResource struct {
	fw.ResourceWithConfigure
}

// AgsSandboxToolResourceModel maps the top-level schema attributes.
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

// agsNetworkConfigurationModel maps network_configuration.
type agsNetworkConfigurationModel struct {
	NetworkMode types.String `tfsdk:"network_mode"`
	VpcConfig   types.List   `tfsdk:"vpc_config"`
}

// agsVpcConfigModel maps network_configuration.vpc_config.
type agsVpcConfigModel struct {
	SubnetIds        types.Set `tfsdk:"subnet_ids"`
	SecurityGroupIds types.Set `tfsdk:"security_group_ids"`
}

// agsTagModel maps tags.
type agsTagModel struct {
	Key   types.String `tfsdk:"key"`
	Value types.String `tfsdk:"value"`
}

// agsStorageMountModel maps storage_mounts.
type agsStorageMountModel struct {
	Name           types.String `tfsdk:"name"`
	StorageSource  types.List   `tfsdk:"storage_source"`
	MountPath      types.String `tfsdk:"mount_path"`
	ReadOnly       types.Bool   `tfsdk:"read_only"`
	Volume         types.List   `tfsdk:"volume"`
	VolumeTemplate types.List   `tfsdk:"volume_template"`
}

// agsStorageSourceModel maps storage_mounts.storage_source.
type agsStorageSourceModel struct {
	Cos         types.List `tfsdk:"cos"`
	Image       types.List `tfsdk:"image"`
	Cfs         types.List `tfsdk:"cfs"`
	AgentBucket types.List `tfsdk:"agent_bucket"`
}

// agsCosStorageSourceModel maps storage_mounts.storage_source.cos.
type agsCosStorageSourceModel struct {
	Endpoint   types.String `tfsdk:"endpoint"`
	BucketName types.String `tfsdk:"bucket_name"`
	BucketPath types.String `tfsdk:"bucket_path"`
}

// agsImageStorageSourceModel maps storage_mounts.storage_source.image.
type agsImageStorageSourceModel struct {
	Reference         types.String `tfsdk:"reference"`
	ImageRegistryType types.String `tfsdk:"image_registry_type"`
	SubPath           types.String `tfsdk:"sub_path"`
	Digest            types.String `tfsdk:"digest"`
}

// agsCfsStorageSourceModel maps storage_mounts.storage_source.cfs.
type agsCfsStorageSourceModel struct {
	FileSystemId types.String `tfsdk:"file_system_id"`
	Path         types.String `tfsdk:"path"`
}

// agsAgentBucketStorageSourceModel maps storage_mounts.storage_source.agent_bucket.
type agsAgentBucketStorageSourceModel struct {
	LibraryId    types.String `tfsdk:"library_id"`
	SpaceId      types.String `tfsdk:"space_id"`
	AccessDomain types.String `tfsdk:"access_domain"`
}

// agsVolumeRefModel maps storage_mounts.volume.
type agsVolumeRefModel struct {
	VolumeName types.String `tfsdk:"volume_name"`
	VolumeId   types.String `tfsdk:"volume_id"`
}

// agsVolumeTemplateRefModel maps storage_mounts.volume_template.
type agsVolumeTemplateRefModel struct {
	VolumeTemplateName types.String `tfsdk:"volume_template_name"`
	ReuseKey           types.String `tfsdk:"reuse_key"`
}

// agsCustomConfigurationModel maps custom_configuration.
type agsCustomConfigurationModel struct {
	Image             types.String `tfsdk:"image"`
	ImageRegistryType types.String `tfsdk:"image_registry_type"`
	Command           types.Set    `tfsdk:"command"`
	Args              types.Set    `tfsdk:"args"`
	Env               types.List   `tfsdk:"env"`
	Ports             types.List   `tfsdk:"ports"`
	Resources         types.List   `tfsdk:"resources"`
	Probe             types.List   `tfsdk:"probe"`
	DNSConfig         types.List   `tfsdk:"dns_config"`
	SecurityContext   types.List   `tfsdk:"security_context"`
}

// agsEnvVarModel maps custom_configuration.env.
type agsEnvVarModel struct {
	Name  types.String `tfsdk:"name"`
	Value types.String `tfsdk:"value"`
}

// agsPortConfigurationModel maps custom_configuration.ports.
type agsPortConfigurationModel struct {
	Name     types.String `tfsdk:"name"`
	Port     types.Int64  `tfsdk:"port"`
	Protocol types.String `tfsdk:"protocol"`
}

// agsResourceConfigurationModel maps custom_configuration.resources.
type agsResourceConfigurationModel struct {
	Cpu                types.String `tfsdk:"cpu"`
	Memory             types.String `tfsdk:"memory"`
	Storage            types.String `tfsdk:"storage"`
	GpuResourceGroupId types.String `tfsdk:"gpu_resource_group_id"`
}

// agsProbeConfigurationModel maps custom_configuration.probe.
type agsProbeConfigurationModel struct {
	HttpGet          types.List  `tfsdk:"http_get"`
	ReadyTimeoutMs   types.Int64 `tfsdk:"ready_timeout_ms"`
	ProbeTimeoutMs   types.Int64 `tfsdk:"probe_timeout_ms"`
	ProbePeriodMs    types.Int64 `tfsdk:"probe_period_ms"`
	SuccessThreshold types.Int64 `tfsdk:"success_threshold"`
	FailureThreshold types.Int64 `tfsdk:"failure_threshold"`
}

// agsHttpGetActionModel maps custom_configuration.probe.http_get.
type agsHttpGetActionModel struct {
	Path   types.String `tfsdk:"path"`
	Port   types.Int64  `tfsdk:"port"`
	Scheme types.String `tfsdk:"scheme"`
}

// agsDNSConfigModel maps custom_configuration.dns_config.
type agsDNSConfigModel struct {
	Servers  types.Set `tfsdk:"servers"`
	Searches types.Set `tfsdk:"searches"`
	Options  types.Set `tfsdk:"options"`
}

// agsSecurityContextModel maps custom_configuration.security_context.
type agsSecurityContextModel struct {
	Privileged   types.Bool `tfsdk:"privileged"`
	Capabilities types.List `tfsdk:"capabilities"`
}

// agsCapabilitiesModel maps custom_configuration.security_context.capabilities.
type agsCapabilitiesModel struct {
	Add  types.Set `tfsdk:"add"`
	Drop types.Set `tfsdk:"drop"`
}

// agsComputerConfigurationModel maps computer_configuration.
type agsComputerConfigurationModel struct {
	WaaConfiguration     types.List `tfsdk:"waa_configuration"`
	OsWorldConfiguration types.List `tfsdk:"os_world_configuration"`
}

// agsWaaConfigurationModel maps computer_configuration.waa_configuration.
type agsWaaConfigurationModel struct {
	ImageId types.String `tfsdk:"image_id"`
}

// agsOsWorldConfigurationModel maps computer_configuration.os_world_configuration.
type agsOsWorldConfigurationModel struct {
	Version types.String `tfsdk:"version"`
}

// agsLogConfigurationModel maps log_configuration.
type agsLogConfigurationModel struct {
	ClsConfig  types.List `tfsdk:"cls_config"`
	LogSources types.List `tfsdk:"log_sources"`
}

// agsClsConfigModel maps log_configuration.cls_config.
type agsClsConfigModel struct {
	TopicId types.String `tfsdk:"topic_id"`
}

// agsLogSourcesModel maps log_configuration.log_sources.
type agsLogSourcesModel struct {
	Files types.Set `tfsdk:"files"`
}

func (r *AgsSandboxToolResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "tencentcloud_ags_sandbox_tool"
}

func (r *AgsSandboxToolResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a resource to create a sandbox tool of AGS (Agent Development Platform).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Resource identifier, the sandbox tool ID.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"tool_id": schema.StringAttribute{
				Computed:    true,
				Description: "The created sandbox tool ID.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"tool_name": schema.StringAttribute{
				Required: true,
				Description: "Sandbox tool name, 1-50 characters, supports English letters, digits, underscores and hyphens. " +
					"The name must be unique under the same AppId. Immutable: the Update API does not accept it, so changing it forces replacement.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"tool_type": schema.StringAttribute{
				Required: true,
				Description: "Sandbox tool type. Valid values: `browser`, `code-interpreter`, `mobile`, `osworld`, `custom`, " +
					"`swebench`, `aio`, `android-world`, `waa`. Immutable: the Update API does not accept it, so changing it forces replacement.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"network_configuration": networkConfigurationAttribute(true),
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "Sandbox tool description, up to 200 characters.",
			},
			"default_timeout": schema.StringAttribute{
				Optional: true,
				Description: "Default timeout, supports formats such as `5m`, `300s`, `1h`. Defaults to the system value (5 minutes) " +
					"when unset, maximum 24 hours. Immutable: the Update API does not accept it, so changing it forces replacement. " +
					"Note: the Describe API returns the value in seconds (default_timeout_seconds), so drift detection on this field is limited.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"tags": schema.ListNestedAttribute{
				Optional:    true,
				Description: "Tag specification, binds tags to the sandbox tool.",
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
				Optional:    true,
				Description: "Idempotency token for the create call, up to 64 characters. Immutable after creation.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"role_arn": schema.StringAttribute{
				Optional:    true,
				Description: "Role ARN. Immutable: the Update API does not accept it, so changing it forces replacement.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"storage_mounts": schema.ListNestedAttribute{
				Optional: true,
				Description: "Storage mount configurations of the sandbox tool. Immutable: the Update API does not accept it, " +
					"so changing it forces replacement.",
				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplace(),
				},
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Optional:    true,
							Description: "Storage mount configuration name.",
						},
						"storage_source": schema.ListNestedAttribute{
							Optional:    true,
							Description: "Storage configuration. This field may return null.",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"cos": schema.ListNestedAttribute{
										Optional:    true,
										Description: "Object storage (COS) bucket configuration.",
										NestedObject: schema.NestedAttributeObject{
											Attributes: map[string]schema.Attribute{
												"endpoint": schema.StringAttribute{
													Optional:    true,
													Description: "Object storage access domain.",
												},
												"bucket_name": schema.StringAttribute{
													Optional:    true,
													Description: "Object storage bucket name.",
												},
												"bucket_path": schema.StringAttribute{
													Optional:    true,
													Description: "Object storage bucket path, must be an absolute path starting with `/`.",
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
												"digest": schema.StringAttribute{
													Optional:    true,
													Description: "Image digest. Not required on requests.",
												},
											},
										},
									},
									"cfs": schema.ListNestedAttribute{
										Optional:    true,
										Description: "File storage (CFS) configuration.",
										NestedObject: schema.NestedAttributeObject{
											Attributes: map[string]schema.Attribute{
												"file_system_id": schema.StringAttribute{
													Optional:    true,
													Description: "CFS resource ID.",
												},
												"path": schema.StringAttribute{
													Optional:    true,
													Description: "CFS mount path.",
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
													Description: "AgentBucket LibraryID.",
												},
												"space_id": schema.StringAttribute{
													Optional:    true,
													Description: "AgentBucket spaceId.",
												},
												"access_domain": schema.StringAttribute{
													Optional:    true,
													Description: "AgentBucket AccessDomain.",
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
							Description: "Whether the mount is read-only, defaults to false.",
						},
						"volume": schema.ListNestedAttribute{
							Optional:    true,
							Description: "Volume reference. This field may return null.",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"volume_name": schema.StringAttribute{
										Optional:    true,
										Description: "Volume name; resolved server-side to a stable VolumeId under the current AppId.",
									},
									"volume_id": schema.StringAttribute{
										Optional:    true,
										Description: "Stable Volume ID, format `vol-xxxxxxxx`.",
									},
								},
							},
						},
						"volume_template": schema.ListNestedAttribute{
							Optional:    true,
							Description: "VolumeTemplate reference and derivation rules. This field may return null.",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"volume_template_name": schema.StringAttribute{
										Required:    true,
										Description: "VolumeTemplate name.",
									},
									"reuse_key": schema.StringAttribute{
										Optional:    true,
										Description: "Reuse key, a literal or `${metadata.xxx}`.",
									},
								},
							},
						},
					},
				},
			},
			"custom_configuration": schema.ListNestedAttribute{
				Optional:    true,
				Description: "Custom sandbox tool configuration.",
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
							Description: "Start command.",
						},
						"args": schema.SetAttribute{
							Optional:    true,
							ElementType: types.StringType,
							Description: "Start arguments.",
						},
						"env": schema.ListNestedAttribute{
							Optional:    true,
							Description: "Environment variables.",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"name": schema.StringAttribute{
										Optional:    true,
										Description: "Environment variable name.",
									},
									"value": schema.StringAttribute{
										Optional:    true,
										Description: "Environment variable value.",
									},
								},
							},
						},
						"ports": schema.ListNestedAttribute{
							Optional:    true,
							Description: "Port configurations.",
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
									"gpu_resource_group_id": schema.StringAttribute{
										Optional:    true,
										Description: "Associated GPU resource group ID.",
									},
								},
							},
						},
						"probe": schema.ListNestedAttribute{
							Optional:    true,
							Description: "Probe configuration.",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"http_get": schema.ListNestedAttribute{
										Optional:    true,
										Description: "HTTP GET probe configuration.",
										NestedObject: schema.NestedAttributeObject{
											Attributes: map[string]schema.Attribute{
												"path": schema.StringAttribute{
													Optional:    true,
													Description: "Path.",
												},
												"port": schema.Int64Attribute{
													Optional:    true,
													Description: "Port number.",
												},
												"scheme": schema.StringAttribute{
													Optional:    true,
													Description: "Scheme.",
												},
											},
										},
									},
									"ready_timeout_ms": schema.Int64Attribute{
										Optional:    true,
										Description: "Health check readiness timeout in milliseconds.",
									},
									"probe_timeout_ms": schema.Int64Attribute{
										Optional:    true,
										Description: "Health check single probe timeout in milliseconds.",
									},
									"probe_period_ms": schema.Int64Attribute{
										Optional:    true,
										Description: "Health check interval in milliseconds.",
									},
									"success_threshold": schema.Int64Attribute{
										Optional:    true,
										Description: "Health check success threshold.",
									},
									"failure_threshold": schema.Int64Attribute{
										Optional:    true,
										Description: "Health check failure threshold.",
									},
								},
							},
						},
						"dns_config": schema.ListNestedAttribute{
							Optional:    true,
							Description: "Sandbox DNS configuration.",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"servers": schema.SetAttribute{
										Required:    true,
										ElementType: types.StringType,
										Description: "DNS server addresses, valid IP addresses. Defaults to `10.0.0.1`.",
									},
									"searches": schema.SetAttribute{
										Optional:    true,
										ElementType: types.StringType,
										Description: "Search domains (the `search` directive of resolv.conf).",
									},
									"options": schema.SetAttribute{
										Optional:    true,
										ElementType: types.StringType,
										Description: "Options (the `options` directive of resolv.conf).",
									},
								},
							},
						},
						"security_context": schema.ListNestedAttribute{
							Optional:    true,
							Description: "Sandbox security context.",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"privileged": schema.BoolAttribute{
										Optional:    true,
										Description: "Whether privileged mode is enabled. Defaults to true.",
									},
									"capabilities": schema.ListNestedAttribute{
										Optional:    true,
										Description: "POSIX capabilities added to or removed from the running container.",
										NestedObject: schema.NestedAttributeObject{
											Attributes: map[string]schema.Attribute{
												"add": schema.SetAttribute{
													Optional:    true,
													ElementType: types.StringType,
													Description: "Capabilities to add, e.g. `[\"SYS_ADMIN\"]`.",
												},
												"drop": schema.SetAttribute{
													Optional:    true,
													ElementType: types.StringType,
													Description: "Capabilities to drop, e.g. `[\"ALL\"]`.",
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
				Description: "Desktop computer environment sandbox configuration.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"waa_configuration": schema.ListNestedAttribute{
							Optional:    true,
							Description: "WAA sandbox tool configuration.",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"image_id": schema.StringAttribute{
										Optional:    true,
										Description: "Custom WAA image ID.",
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
										Optional:    true,
										Description: "Built-in OSWorld version. Valid values: `osworld1`, `osworld2`. Defaults to `osworld1`.",
									},
								},
							},
						},
					},
				},
			},
			"log_configuration": schema.ListNestedAttribute{
				Optional: true,
				Description: "Sandbox tool log pushing configuration. Immutable: the Update API does not accept it, " +
					"so changing it forces replacement.",
				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplace(),
				},
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"cls_config": schema.ListNestedAttribute{
							Optional:    true,
							Description: "CLS log pushing configuration.",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"topic_id": schema.StringAttribute{
										Optional:    true,
										Description: "CLS log topic ID used by the sandbox tool log pushing.",
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
										Description: "Log file paths to collect, must be files under the `/logs/` directory, " +
											"subdirectories are not supported, up to 10 files.",
									},
								},
							},
						},
					},
				},
			},
			"persistent": schema.BoolAttribute{
				Optional:    true,
				Description: "Whether the sandbox is resident (persistent). Immutable: the Update API does not accept it, so changing it forces replacement.",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
			},
			"status": schema.StringAttribute{
				Computed:    true,
				Description: "Sandbox tool status. Valid values: `CREATING`, `ACTIVE`, `DELETING`, `FAILED`.",
			},
			"status_reason": schema.StringAttribute{
				Computed:    true,
				Description: "Reason why the sandbox tool is in its current status.",
			},
			"create_time": schema.StringAttribute{
				Computed:    true,
				Description: "Creation time of the sandbox tool, ISO8601 format.",
			},
			"update_time": schema.StringAttribute{
				Computed:    true,
				Description: "Last update time of the sandbox tool, ISO8601 format.",
			},
		},
	}
}

// Create handles `terraform apply` for new resources.
func (r *AgsSandboxToolResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan AgsSandboxToolResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, ok := r.agsClient(&resp.Diagnostics)
	if !ok {
		return
	}

	request := agssdk.NewCreateSandboxToolRequest()
	request.ToolName = fwhelper.StringPointerFromValue(plan.ToolName)
	request.ToolType = fwhelper.StringPointerFromValue(plan.ToolType)
	request.Description = fwhelper.StringPointerFromValue(plan.Description)
	request.DefaultTimeout = fwhelper.StringPointerFromValue(plan.DefaultTimeout)
	request.ClientToken = fwhelper.StringPointerFromValue(plan.ClientToken)
	request.RoleArn = fwhelper.StringPointerFromValue(plan.RoleArn)
	request.Persistent = fwhelper.BoolPointerFromValue(plan.Persistent)

	networkConfiguration, networkDiags := expandAgsNetworkConfiguration(ctx, plan.NetworkConfiguration)
	resp.Diagnostics.Append(networkDiags...)
	if resp.Diagnostics.HasError() {
		return
	}
	request.NetworkConfiguration = networkConfiguration

	tags, tagsDiags := expandAgsTags(ctx, plan.Tags)
	resp.Diagnostics.Append(tagsDiags...)
	if resp.Diagnostics.HasError() {
		return
	}
	request.Tags = tags

	storageMounts, storageDiags := expandAgsStorageMounts(ctx, plan.StorageMounts)
	resp.Diagnostics.Append(storageDiags...)
	if resp.Diagnostics.HasError() {
		return
	}
	request.StorageMounts = storageMounts

	customConfiguration, customDiags := expandAgsCustomConfiguration(ctx, plan.CustomConfiguration)
	resp.Diagnostics.Append(customDiags...)
	if resp.Diagnostics.HasError() {
		return
	}
	request.CustomConfiguration = customConfiguration

	computerConfiguration, computerDiags := expandAgsComputerConfiguration(ctx, plan.ComputerConfiguration)
	resp.Diagnostics.Append(computerDiags...)
	if resp.Diagnostics.HasError() {
		return
	}
	request.ComputerConfiguration = computerConfiguration

	logConfiguration, logDiags := expandAgsLogConfiguration(ctx, plan.LogConfiguration)
	resp.Diagnostics.Append(logDiags...)
	if resp.Diagnostics.HasError() {
		return
	}
	request.LogConfiguration = logConfiguration

	response, err := client.CreateSandboxToolWithContext(ctx, request)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating AGS sandbox tool",
			fmt.Sprintf("CreateSandboxTool failed (tool_name=%q): %s", plan.ToolName.ValueString(), err),
		)
		return
	}
	if response == nil || response.Response == nil || response.Response.ToolId == nil || *response.Response.ToolId == "" {
		resp.Diagnostics.AddError(
			"Error creating AGS sandbox tool",
			fmt.Sprintf("CreateSandboxTool returned no ToolId (tool_name=%q).", plan.ToolName.ValueString()),
		)
		return
	}

	toolId := *response.Response.ToolId
	// Persist the identity before any further read-back so a later failure
	// still leaves a destroyable resource behind.
	plan.Id = types.StringValue(toolId)
	plan.ToolId = types.StringValue(toolId)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	item, err := describeAgsSandboxTool(ctx, client, toolId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading AGS sandbox tool after create",
			fmt.Sprintf("CreateSandboxTool succeeded but the tool %s cannot be read back: %s", toolId, err),
		)
		return
	}
	if item == nil {
		resp.Diagnostics.AddWarning(
			"AGS sandbox tool not readable yet",
			fmt.Sprintf("The sandbox tool %s was created but is not returned by DescribeSandboxToolList yet; "+
				"state keeps the planned values until the next refresh.", toolId),
		)
		return
	}

	resp.Diagnostics.Append(refreshAgsSandboxToolModel(ctx, &plan, item)...)
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

	client, ok := r.agsClient(&resp.Diagnostics)
	if !ok {
		return
	}

	toolId := state.Id.ValueString()
	item, err := describeAgsSandboxTool(ctx, client, toolId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading AGS sandbox tool",
			fmt.Sprintf("DescribeSandboxToolList failed (tool_id=%q): %s", toolId, err),
		)
		return
	}
	if item == nil {
		// A confirmed "absent from the list" means the tool was removed
		// outside Terraform, so it is dropped from state.
		resp.Diagnostics.AddWarning(
			"Resource not found",
			fmt.Sprintf("AGS sandbox tool %s no longer exists, removing it from state.", toolId),
		)
		resp.State.RemoveResource(ctx)
		return
	}

	resp.Diagnostics.Append(refreshAgsSandboxToolModel(ctx, &state, item)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update handles in-place updates. Only the fields accepted by
// UpdateSandboxTool are sent; the remaining attributes are marked immutable
// with RequiresReplace in the schema.
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

	client, ok := r.agsClient(&resp.Diagnostics)
	if !ok {
		return
	}

	toolId := state.Id.ValueString()
	request := agssdk.NewUpdateSandboxToolRequest()
	request.ToolId = helper.String(toolId)
	request.Description = fwhelper.StringPointerFromValue(plan.Description)

	networkConfiguration, networkDiags := expandAgsNetworkConfiguration(ctx, plan.NetworkConfiguration)
	resp.Diagnostics.Append(networkDiags...)
	if resp.Diagnostics.HasError() {
		return
	}
	request.NetworkConfiguration = networkConfiguration

	// Tags are passed through as configured. Whether the cloud treats an
	// empty or missing tag list as "clear all tags" is not documented, so the
	// plan value is sent verbatim.
	tags, tagsDiags := expandAgsTags(ctx, plan.Tags)
	resp.Diagnostics.Append(tagsDiags...)
	if resp.Diagnostics.HasError() {
		return
	}
	request.Tags = tags

	customConfiguration, customDiags := expandAgsCustomConfiguration(ctx, plan.CustomConfiguration)
	resp.Diagnostics.Append(customDiags...)
	if resp.Diagnostics.HasError() {
		return
	}
	request.CustomConfiguration = customConfiguration

	computerConfiguration, computerDiags := expandAgsComputerConfiguration(ctx, plan.ComputerConfiguration)
	resp.Diagnostics.Append(computerDiags...)
	if resp.Diagnostics.HasError() {
		return
	}
	request.ComputerConfiguration = computerConfiguration

	if _, err := client.UpdateSandboxToolWithContext(ctx, request); err != nil {
		resp.Diagnostics.AddError(
			"Error updating AGS sandbox tool",
			fmt.Sprintf("UpdateSandboxTool failed (tool_id=%q): %s", toolId, err),
		)
		return
	}

	item, err := describeAgsSandboxTool(ctx, client, toolId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading AGS sandbox tool after update",
			fmt.Sprintf("UpdateSandboxTool succeeded but the tool %s cannot be read back: %s", toolId, err),
		)
		return
	}
	if item == nil {
		resp.Diagnostics.AddWarning(
			"AGS sandbox tool not readable yet",
			fmt.Sprintf("The sandbox tool %s was updated but is not returned by DescribeSandboxToolList yet; "+
				"state keeps the planned values until the next refresh.", toolId),
		)
		resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
		return
	}

	resp.Diagnostics.Append(refreshAgsSandboxToolModel(ctx, &plan, item)...)
	if resp.Diagnostics.HasError() {
		return
	}
	// Fields the Describe API does not echo back keep their planned values.
	plan.DefaultTimeout = state.DefaultTimeout
	plan.ClientToken = state.ClientToken
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete handles resource deletion.
func (r *AgsSandboxToolResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state AgsSandboxToolResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, ok := r.agsClient(&resp.Diagnostics)
	if !ok {
		return
	}

	toolId := state.Id.ValueString()
	request := agssdk.NewDeleteSandboxToolRequest()
	request.ToolId = helper.String(toolId)

	if _, err := client.DeleteSandboxToolWithContext(ctx, request); err != nil {
		resp.Diagnostics.AddError(
			"Error deleting AGS sandbox tool",
			fmt.Sprintf("DeleteSandboxTool failed (tool_id=%q): %s", toolId, err),
		)
	}
}

// ImportState supports `terraform import` using the ToolId. The remaining
// attributes are filled in by the Read that follows.
func (r *AgsSandboxToolResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if req.ID == "" {
		resp.Diagnostics.AddError(
			"Invalid import ID",
			fmt.Sprintf("Import id %q is invalid, expected the sandbox tool ID.", req.ID),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
}

// agsClient returns the AGS SDK client bound to the shared provider client.
func (r *AgsSandboxToolResource) agsClient(diags *diag.Diagnostics) (*agssdk.Client, bool) {
	if r.Client() == nil {
		diags.AddError(
			"Provider not configured",
			"The provider client is not available. This is unexpected; please report it.",
		)
		return nil, false
	}
	return r.Client().UseAgsV20250920Client(), true
}

// describeAgsSandboxTool reads one sandbox tool by ID. It returns nil (without
// error) when the tool is absent from the list, so callers can distinguish a
// confirmed deletion from an API failure.
func describeAgsSandboxTool(ctx context.Context, client *agssdk.Client, toolId string) (*agssdk.SandboxTool, error) {
	request := agssdk.NewDescribeSandboxToolListRequest()
	request.ToolIds = []*string{helper.String(toolId)}
	request.Limit = helper.Int64(100)
	request.Offset = helper.Int64(0)

	response, err := client.DescribeSandboxToolListWithContext(ctx, request)
	if err != nil {
		return nil, err
	}
	if response == nil || response.Response == nil {
		return nil, fmt.Errorf("DescribeSandboxToolList returned an empty response (tool_id=%q)", toolId)
	}
	for _, item := range response.Response.SandboxToolSet {
		if item != nil && item.ToolId != nil && *item.ToolId == toolId {
			return item, nil
		}
	}
	return nil, nil
}

// refreshAgsSandboxToolModel copies the cloud state into the model. Attributes
// the Describe API may omit keep their previous model value, so a value
// configured by the practitioner is never nulled out after apply.
func refreshAgsSandboxToolModel(ctx context.Context, m *AgsSandboxToolResourceModel, item *agssdk.SandboxTool) diag.Diagnostics {
	var diags diag.Diagnostics

	if item == nil {
		diags.AddError("Malformed AGS sandbox tool response",
			"DescribeSandboxToolList returned an empty tool entry.")
		return diags
	}

	m.ToolId = fwhelper.StringValueOrPrev(item.ToolId, m.ToolId)
	if m.Id.IsNull() || m.Id.IsUnknown() || m.Id.ValueString() == "" {
		m.Id = m.ToolId
	}
	m.ToolName = fwhelper.StringValueOrPrev(item.ToolName, m.ToolName)
	m.ToolType = fwhelper.StringValueOrPrev(item.ToolType, m.ToolType)
	m.Description = fwhelper.StringValueOrPrev(item.Description, m.Description)
	m.RoleArn = fwhelper.StringValueOrPrev(item.RoleArn, m.RoleArn)
	m.Persistent = agsBoolValueOrPrev(item.Persistent, m.Persistent)
	m.Status = types.StringPointerValue(item.Status)
	m.StatusReason = types.StringPointerValue(item.StatusReason)
	m.CreateTime = types.StringPointerValue(item.CreateTime)
	m.UpdateTime = types.StringPointerValue(item.UpdateTime)
	// default_timeout and client_token are not echoed by DescribeSandboxToolList;
	// the caller decides whether to keep the planned/state value.

	networkConfiguration, networkDiags := flattenAgsNetworkConfiguration(ctx, item.NetworkConfiguration, m.NetworkConfiguration)
	diags.Append(networkDiags...)
	m.NetworkConfiguration = networkConfiguration

	tags, tagsDiags := flattenAgsTags(ctx, item.Tags, m.Tags)
	diags.Append(tagsDiags...)
	m.Tags = tags

	storageMounts, storageDiags := flattenAgsStorageMounts(ctx, item.StorageMounts, m.StorageMounts)
	diags.Append(storageDiags...)
	m.StorageMounts = storageMounts

	customConfiguration, customDiags := flattenAgsCustomConfiguration(ctx, item.CustomConfiguration, m.CustomConfiguration)
	diags.Append(customDiags...)
	m.CustomConfiguration = customConfiguration

	computerConfiguration, computerDiags := flattenAgsComputerConfiguration(ctx, item.ComputerConfiguration, m.ComputerConfiguration)
	diags.Append(computerDiags...)
	m.ComputerConfiguration = computerConfiguration

	logConfiguration, logDiags := flattenAgsLogConfiguration(ctx, item.LogConfiguration, m.LogConfiguration)
	diags.Append(logDiags...)
	m.LogConfiguration = logConfiguration

	if m.ToolName.IsNull() || m.NetworkConfiguration.IsNull() {
		diags.AddError(
			"Malformed AGS sandbox tool response",
			fmt.Sprintf("DescribeSandboxToolList returned no tool_name or network_configuration for tool %q.", m.ToolId.ValueString()),
		)
	}
	return diags
}

// ---- expand helpers (Terraform model -> SDK request) ----

func expandAgsNetworkConfiguration(ctx context.Context, list types.List) (*agssdk.NetworkConfiguration, diag.Diagnostics) {
	var diags diag.Diagnostics
	if list.IsNull() || list.IsUnknown() {
		return nil, diags
	}
	var models []agsNetworkConfigurationModel
	diags.Append(list.ElementsAs(ctx, &models, false)...)
	if diags.HasError() || len(models) == 0 {
		return nil, diags
	}
	model := models[0]

	out := &agssdk.NetworkConfiguration{
		NetworkMode: fwhelper.StringPointerFromValue(model.NetworkMode),
	}
	if !model.VpcConfig.IsNull() && !model.VpcConfig.IsUnknown() {
		var vpcModels []agsVpcConfigModel
		diags.Append(model.VpcConfig.ElementsAs(ctx, &vpcModels, false)...)
		if diags.HasError() {
			return nil, diags
		}
		if len(vpcModels) > 0 {
			subnetIds, subnetDiags := agsStringSetToSlice(ctx, vpcModels[0].SubnetIds)
			diags.Append(subnetDiags...)
			securityGroupIds, sgDiags := agsStringSetToSlice(ctx, vpcModels[0].SecurityGroupIds)
			diags.Append(sgDiags...)
			if diags.HasError() {
				return nil, diags
			}
			out.VpcConfig = &agssdk.VPCConfig{
				SubnetIds:        subnetIds,
				SecurityGroupIds: securityGroupIds,
			}
		}
	}
	return out, diags
}

func expandAgsTags(ctx context.Context, list types.List) ([]*agssdk.Tag, diag.Diagnostics) {
	var diags diag.Diagnostics
	if list.IsNull() || list.IsUnknown() {
		return nil, diags
	}
	var models []agsTagModel
	diags.Append(list.ElementsAs(ctx, &models, false)...)
	if diags.HasError() {
		return nil, diags
	}
	out := make([]*agssdk.Tag, 0, len(models))
	for _, model := range models {
		out = append(out, &agssdk.Tag{
			Key:   fwhelper.StringPointerFromValue(model.Key),
			Value: fwhelper.StringPointerFromValue(model.Value),
		})
	}
	return out, diags
}

func expandAgsStorageMounts(ctx context.Context, list types.List) ([]*agssdk.StorageMount, diag.Diagnostics) {
	var diags diag.Diagnostics
	if list.IsNull() || list.IsUnknown() {
		return nil, diags
	}
	var models []agsStorageMountModel
	diags.Append(list.ElementsAs(ctx, &models, false)...)
	if diags.HasError() {
		return nil, diags
	}
	out := make([]*agssdk.StorageMount, 0, len(models))
	for _, model := range models {
		mount := &agssdk.StorageMount{
			Name:      fwhelper.StringPointerFromValue(model.Name),
			MountPath: fwhelper.StringPointerFromValue(model.MountPath),
			ReadOnly:  fwhelper.BoolPointerFromValue(model.ReadOnly),
		}
		if !model.StorageSource.IsNull() && !model.StorageSource.IsUnknown() {
			var sourceModels []agsStorageSourceModel
			diags.Append(model.StorageSource.ElementsAs(ctx, &sourceModels, false)...)
			if diags.HasError() {
				return nil, diags
			}
			if len(sourceModels) > 0 {
				source := &agssdk.StorageSource{}
				cos, cosDiags := expandAgsCosStorageSource(ctx, sourceModels[0].Cos)
				diags.Append(cosDiags...)
				source.Cos = cos
				image, imageDiags := expandAgsImageStorageSource(ctx, sourceModels[0].Image)
				diags.Append(imageDiags...)
				source.Image = image
				cfs, cfsDiags := expandAgsCfsStorageSource(ctx, sourceModels[0].Cfs)
				diags.Append(cfsDiags...)
				source.Cfs = cfs
				agentBucket, agentDiags := expandAgsAgentBucketStorageSource(ctx, sourceModels[0].AgentBucket)
				diags.Append(agentDiags...)
				source.AgentBucket = agentBucket
				if diags.HasError() {
					return nil, diags
				}
				mount.StorageSource = source
			}
		}
		if !model.Volume.IsNull() && !model.Volume.IsUnknown() {
			var volumeModels []agsVolumeRefModel
			diags.Append(model.Volume.ElementsAs(ctx, &volumeModels, false)...)
			if diags.HasError() {
				return nil, diags
			}
			if len(volumeModels) > 0 {
				mount.Volume = &agssdk.VolumeRef{
					VolumeName: fwhelper.StringPointerFromValue(volumeModels[0].VolumeName),
					VolumeId:   fwhelper.StringPointerFromValue(volumeModels[0].VolumeId),
				}
			}
		}
		if !model.VolumeTemplate.IsNull() && !model.VolumeTemplate.IsUnknown() {
			var templateModels []agsVolumeTemplateRefModel
			diags.Append(model.VolumeTemplate.ElementsAs(ctx, &templateModels, false)...)
			if diags.HasError() {
				return nil, diags
			}
			if len(templateModels) > 0 {
				mount.VolumeTemplate = &agssdk.VolumeTemplateRef{
					VolumeTemplateName: fwhelper.StringPointerFromValue(templateModels[0].VolumeTemplateName),
					ReuseKey:           fwhelper.StringPointerFromValue(templateModels[0].ReuseKey),
				}
			}
		}
		out = append(out, mount)
	}
	return out, diags
}

func expandAgsCosStorageSource(ctx context.Context, list types.List) (*agssdk.CosStorageSource, diag.Diagnostics) {
	var diags diag.Diagnostics
	if list.IsNull() || list.IsUnknown() {
		return nil, diags
	}
	var models []agsCosStorageSourceModel
	diags.Append(list.ElementsAs(ctx, &models, false)...)
	if diags.HasError() || len(models) == 0 {
		return nil, diags
	}
	return &agssdk.CosStorageSource{
		Endpoint:   fwhelper.StringPointerFromValue(models[0].Endpoint),
		BucketName: fwhelper.StringPointerFromValue(models[0].BucketName),
		BucketPath: fwhelper.StringPointerFromValue(models[0].BucketPath),
	}, diags
}

func expandAgsImageStorageSource(ctx context.Context, list types.List) (*agssdk.ImageStorageSource, diag.Diagnostics) {
	var diags diag.Diagnostics
	if list.IsNull() || list.IsUnknown() {
		return nil, diags
	}
	var models []agsImageStorageSourceModel
	diags.Append(list.ElementsAs(ctx, &models, false)...)
	if diags.HasError() || len(models) == 0 {
		return nil, diags
	}
	return &agssdk.ImageStorageSource{
		Reference:         fwhelper.StringPointerFromValue(models[0].Reference),
		ImageRegistryType: fwhelper.StringPointerFromValue(models[0].ImageRegistryType),
		SubPath:           fwhelper.StringPointerFromValue(models[0].SubPath),
		Digest:            fwhelper.StringPointerFromValue(models[0].Digest),
	}, diags
}

func expandAgsCfsStorageSource(ctx context.Context, list types.List) (*agssdk.CfsStorageSource, diag.Diagnostics) {
	var diags diag.Diagnostics
	if list.IsNull() || list.IsUnknown() {
		return nil, diags
	}
	var models []agsCfsStorageSourceModel
	diags.Append(list.ElementsAs(ctx, &models, false)...)
	if diags.HasError() || len(models) == 0 {
		return nil, diags
	}
	return &agssdk.CfsStorageSource{
		FileSystemId: fwhelper.StringPointerFromValue(models[0].FileSystemId),
		Path:         fwhelper.StringPointerFromValue(models[0].Path),
	}, diags
}

func expandAgsAgentBucketStorageSource(ctx context.Context, list types.List) (*agssdk.AgentBucketStorageSource, diag.Diagnostics) {
	var diags diag.Diagnostics
	if list.IsNull() || list.IsUnknown() {
		return nil, diags
	}
	var models []agsAgentBucketStorageSourceModel
	diags.Append(list.ElementsAs(ctx, &models, false)...)
	if diags.HasError() || len(models) == 0 {
		return nil, diags
	}
	return &agssdk.AgentBucketStorageSource{
		LibraryId:    fwhelper.StringPointerFromValue(models[0].LibraryId),
		SpaceId:      fwhelper.StringPointerFromValue(models[0].SpaceId),
		AccessDomain: fwhelper.StringPointerFromValue(models[0].AccessDomain),
	}, diags
}

func expandAgsCustomConfiguration(ctx context.Context, list types.List) (*agssdk.CustomConfiguration, diag.Diagnostics) {
	var diags diag.Diagnostics
	if list.IsNull() || list.IsUnknown() {
		return nil, diags
	}
	var models []agsCustomConfigurationModel
	diags.Append(list.ElementsAs(ctx, &models, false)...)
	if diags.HasError() || len(models) == 0 {
		return nil, diags
	}
	model := models[0]

	command, commandDiags := agsStringSetToSlice(ctx, model.Command)
	diags.Append(commandDiags...)
	args, argsDiags := agsStringSetToSlice(ctx, model.Args)
	diags.Append(argsDiags...)
	if diags.HasError() {
		return nil, diags
	}

	out := &agssdk.CustomConfiguration{
		Image:             fwhelper.StringPointerFromValue(model.Image),
		ImageRegistryType: fwhelper.StringPointerFromValue(model.ImageRegistryType),
		Command:           command,
		Args:              args,
	}
	if !model.Env.IsNull() && !model.Env.IsUnknown() {
		var envModels []agsEnvVarModel
		diags.Append(model.Env.ElementsAs(ctx, &envModels, false)...)
		if diags.HasError() {
			return nil, diags
		}
		for _, env := range envModels {
			out.Env = append(out.Env, &agssdk.EnvVar{
				Name:  fwhelper.StringPointerFromValue(env.Name),
				Value: fwhelper.StringPointerFromValue(env.Value),
			})
		}
	}
	if !model.Ports.IsNull() && !model.Ports.IsUnknown() {
		var portModels []agsPortConfigurationModel
		diags.Append(model.Ports.ElementsAs(ctx, &portModels, false)...)
		if diags.HasError() {
			return nil, diags
		}
		for _, port := range portModels {
			out.Ports = append(out.Ports, &agssdk.PortConfiguration{
				Name:     fwhelper.StringPointerFromValue(port.Name),
				Port:     fwhelper.Int64PointerFromValue(port.Port),
				Protocol: fwhelper.StringPointerFromValue(port.Protocol),
			})
		}
	}
	if !model.Resources.IsNull() && !model.Resources.IsUnknown() {
		var resourceModels []agsResourceConfigurationModel
		diags.Append(model.Resources.ElementsAs(ctx, &resourceModels, false)...)
		if diags.HasError() {
			return nil, diags
		}
		if len(resourceModels) > 0 {
			out.Resources = &agssdk.ResourceConfiguration{
				CPU:                fwhelper.StringPointerFromValue(resourceModels[0].Cpu),
				Memory:             fwhelper.StringPointerFromValue(resourceModels[0].Memory),
				Storage:            fwhelper.StringPointerFromValue(resourceModels[0].Storage),
				GPUResourceGroupId: fwhelper.StringPointerFromValue(resourceModels[0].GpuResourceGroupId),
			}
		}
	}
	if !model.Probe.IsNull() && !model.Probe.IsUnknown() {
		var probeModels []agsProbeConfigurationModel
		diags.Append(model.Probe.ElementsAs(ctx, &probeModels, false)...)
		if diags.HasError() {
			return nil, diags
		}
		if len(probeModels) > 0 {
			probe := &agssdk.ProbeConfiguration{
				ReadyTimeoutMs:   fwhelper.Int64PointerFromValue(probeModels[0].ReadyTimeoutMs),
				ProbeTimeoutMs:   fwhelper.Int64PointerFromValue(probeModels[0].ProbeTimeoutMs),
				ProbePeriodMs:    fwhelper.Int64PointerFromValue(probeModels[0].ProbePeriodMs),
				SuccessThreshold: fwhelper.Int64PointerFromValue(probeModels[0].SuccessThreshold),
				FailureThreshold: fwhelper.Int64PointerFromValue(probeModels[0].FailureThreshold),
			}
			if !probeModels[0].HttpGet.IsNull() && !probeModels[0].HttpGet.IsUnknown() {
				var httpModels []agsHttpGetActionModel
				diags.Append(probeModels[0].HttpGet.ElementsAs(ctx, &httpModels, false)...)
				if diags.HasError() {
					return nil, diags
				}
				if len(httpModels) > 0 {
					probe.HttpGet = &agssdk.HttpGetAction{
						Path:   fwhelper.StringPointerFromValue(httpModels[0].Path),
						Port:   fwhelper.Int64PointerFromValue(httpModels[0].Port),
						Scheme: fwhelper.StringPointerFromValue(httpModels[0].Scheme),
					}
				}
			}
			out.Probe = probe
		}
	}
	if !model.DNSConfig.IsNull() && !model.DNSConfig.IsUnknown() {
		var dnsModels []agsDNSConfigModel
		diags.Append(model.DNSConfig.ElementsAs(ctx, &dnsModels, false)...)
		if diags.HasError() {
			return nil, diags
		}
		if len(dnsModels) > 0 {
			servers, serversDiags := agsStringSetToSlice(ctx, dnsModels[0].Servers)
			diags.Append(serversDiags...)
			searches, searchesDiags := agsStringSetToSlice(ctx, dnsModels[0].Searches)
			diags.Append(searchesDiags...)
			options, optionsDiags := agsStringSetToSlice(ctx, dnsModels[0].Options)
			diags.Append(optionsDiags...)
			if diags.HasError() {
				return nil, diags
			}
			out.DNSConfig = &agssdk.DNSConfig{
				Servers:  servers,
				Searches: searches,
				Options:  options,
			}
		}
	}
	if !model.SecurityContext.IsNull() && !model.SecurityContext.IsUnknown() {
		var securityModels []agsSecurityContextModel
		diags.Append(model.SecurityContext.ElementsAs(ctx, &securityModels, false)...)
		if diags.HasError() {
			return nil, diags
		}
		if len(securityModels) > 0 {
			security := &agssdk.SecurityContext{
				Privileged: fwhelper.BoolPointerFromValue(securityModels[0].Privileged),
			}
			if !securityModels[0].Capabilities.IsNull() && !securityModels[0].Capabilities.IsUnknown() {
				var capModels []agsCapabilitiesModel
				diags.Append(securityModels[0].Capabilities.ElementsAs(ctx, &capModels, false)...)
				if diags.HasError() {
					return nil, diags
				}
				if len(capModels) > 0 {
					add, addDiags := agsStringSetToSlice(ctx, capModels[0].Add)
					diags.Append(addDiags...)
					drop, dropDiags := agsStringSetToSlice(ctx, capModels[0].Drop)
					diags.Append(dropDiags...)
					if diags.HasError() {
						return nil, diags
					}
					security.Capabilities = &agssdk.Capabilities{
						Add:  add,
						Drop: drop,
					}
				}
			}
			out.SecurityContext = security
		}
	}
	return out, diags
}

func expandAgsComputerConfiguration(ctx context.Context, list types.List) (*agssdk.ComputerConfiguration, diag.Diagnostics) {
	var diags diag.Diagnostics
	if list.IsNull() || list.IsUnknown() {
		return nil, diags
	}
	var models []agsComputerConfigurationModel
	diags.Append(list.ElementsAs(ctx, &models, false)...)
	if diags.HasError() || len(models) == 0 {
		return nil, diags
	}
	model := models[0]

	out := &agssdk.ComputerConfiguration{}
	if !model.WaaConfiguration.IsNull() && !model.WaaConfiguration.IsUnknown() {
		var waaModels []agsWaaConfigurationModel
		diags.Append(model.WaaConfiguration.ElementsAs(ctx, &waaModels, false)...)
		if diags.HasError() {
			return nil, diags
		}
		if len(waaModels) > 0 {
			out.WAAConfiguration = &agssdk.WAAConfiguration{
				ImageId: fwhelper.StringPointerFromValue(waaModels[0].ImageId),
			}
		}
	}
	if !model.OsWorldConfiguration.IsNull() && !model.OsWorldConfiguration.IsUnknown() {
		var osModels []agsOsWorldConfigurationModel
		diags.Append(model.OsWorldConfiguration.ElementsAs(ctx, &osModels, false)...)
		if diags.HasError() {
			return nil, diags
		}
		if len(osModels) > 0 {
			out.OSWorldConfiguration = &agssdk.OSWorldConfiguration{
				Version: fwhelper.StringPointerFromValue(osModels[0].Version),
			}
		}
	}
	return out, diags
}

func expandAgsLogConfiguration(ctx context.Context, list types.List) (*agssdk.LogConfiguration, diag.Diagnostics) {
	var diags diag.Diagnostics
	if list.IsNull() || list.IsUnknown() {
		return nil, diags
	}
	var models []agsLogConfigurationModel
	diags.Append(list.ElementsAs(ctx, &models, false)...)
	if diags.HasError() || len(models) == 0 {
		return nil, diags
	}
	model := models[0]

	out := &agssdk.LogConfiguration{}
	if !model.ClsConfig.IsNull() && !model.ClsConfig.IsUnknown() {
		var clsModels []agsClsConfigModel
		diags.Append(model.ClsConfig.ElementsAs(ctx, &clsModels, false)...)
		if diags.HasError() {
			return nil, diags
		}
		if len(clsModels) > 0 {
			out.CLSConfig = &agssdk.CLSConfig{
				TopicId: fwhelper.StringPointerFromValue(clsModels[0].TopicId),
			}
		}
	}
	if !model.LogSources.IsNull() && !model.LogSources.IsUnknown() {
		var sourceModels []agsLogSourcesModel
		diags.Append(model.LogSources.ElementsAs(ctx, &sourceModels, false)...)
		if diags.HasError() {
			return nil, diags
		}
		if len(sourceModels) > 0 {
			files, filesDiags := agsStringSetToSlice(ctx, sourceModels[0].Files)
			diags.Append(filesDiags...)
			if diags.HasError() {
				return nil, diags
			}
			out.LogSources = &agssdk.LogSources{
				Files: files,
			}
		}
	}
	return out, diags
}

// agsStringSetToSlice converts a Terraform string set into a *[]string slice;
// a null or unknown set yields nil.
func agsStringSetToSlice(ctx context.Context, set types.Set) ([]*string, diag.Diagnostics) {
	var diags diag.Diagnostics
	if set.IsNull() || set.IsUnknown() {
		return nil, diags
	}
	var values []string
	diags.Append(set.ElementsAs(ctx, &values, false)...)
	if diags.HasError() {
		return nil, diags
	}
	out := make([]*string, 0, len(values))
	for _, value := range values {
		out = append(out, helper.String(value))
	}
	return out, diags
}

// ---- flatten helpers (SDK response -> Terraform model) ----

func flattenAgsNetworkConfiguration(ctx context.Context, config *agssdk.NetworkConfiguration, prevList types.List) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics
	if config == nil {
		return prevList, diags
	}

	vpcConfig := types.ListNull(agsVpcConfigObjectType())
	if config.VpcConfig != nil {
		var prevModels []agsVpcConfigModel
		if !prevList.IsNull() && !prevList.IsUnknown() {
			var prevNetworkModels []agsNetworkConfigurationModel
			diags.Append(prevList.ElementsAs(ctx, &prevNetworkModels, false)...)
			if !diags.HasError() && len(prevNetworkModels) > 0 && !prevNetworkModels[0].VpcConfig.IsNull() && !prevNetworkModels[0].VpcConfig.IsUnknown() {
				diags.Append(prevNetworkModels[0].VpcConfig.ElementsAs(ctx, &prevModels, false)...)
			}
		}
		var prev agsVpcConfigModel
		if len(prevModels) > 0 {
			prev = prevModels[0]
		}

		subnetIds, subnetDiags := agsFlattenStringSet(ctx, config.VpcConfig.SubnetIds, prev.SubnetIds)
		diags.Append(subnetDiags...)
		securityGroupIds, sgDiags := agsFlattenStringSet(ctx, config.VpcConfig.SecurityGroupIds, prev.SecurityGroupIds)
		diags.Append(sgDiags...)
		if diags.HasError() {
			return types.ListNull(agsNetworkConfigurationObjectType()), diags
		}
		vpcConfig, _ = types.ListValueFrom(ctx, agsVpcConfigObjectType(), []agsVpcConfigModel{{
			SubnetIds:        subnetIds,
			SecurityGroupIds: securityGroupIds,
		}})
	}

	var prevNetworkMode types.String
	if !prevList.IsNull() && !prevList.IsUnknown() {
		var prevNetworkModels []agsNetworkConfigurationModel
		diags.Append(prevList.ElementsAs(ctx, &prevNetworkModels, false)...)
		if diags.HasError() {
			return types.ListNull(agsNetworkConfigurationObjectType()), diags
		}
		if len(prevNetworkModels) > 0 {
			prevNetworkMode = prevNetworkModels[0].NetworkMode
			if config.VpcConfig == nil {
				vpcConfig = prevNetworkModels[0].VpcConfig
			}
		}
	}

	out := []agsNetworkConfigurationModel{{
		NetworkMode: fwhelper.StringValueOrPrev(config.NetworkMode, prevNetworkMode),
		VpcConfig:   vpcConfig,
	}}
	list, listDiags := types.ListValueFrom(ctx, agsNetworkConfigurationObjectType(), out)
	diags.Append(listDiags...)
	return list, diags
}

func flattenAgsTags(ctx context.Context, tags []*agssdk.Tag, prevList types.List) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics
	if len(tags) == 0 {
		// The Describe API omits tags entirely when none are bound; keep the
		// previous value so a configured tag list is not cleared on refresh.
		if prevList.IsNull() || prevList.IsUnknown() {
			return types.ListNull(agsTagObjectType()), diags
		}
		return prevList, diags
	}
	out := make([]agsTagModel, 0, len(tags))
	for _, tag := range tags {
		if tag == nil {
			continue
		}
		out = append(out, agsTagModel{
			Key:   types.StringPointerValue(tag.Key),
			Value: types.StringPointerValue(tag.Value),
		})
	}
	list, listDiags := types.ListValueFrom(ctx, agsTagObjectType(), out)
	diags.Append(listDiags...)
	return list, diags
}

func flattenAgsStorageMounts(ctx context.Context, mounts []*agssdk.StorageMount, prevList types.List) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics
	if len(mounts) == 0 {
		if prevList.IsNull() || prevList.IsUnknown() {
			return types.ListNull(agsStorageMountObjectType()), diags
		}
		return prevList, diags
	}

	var prevModels []agsStorageMountModel
	if !prevList.IsNull() && !prevList.IsUnknown() {
		diags.Append(prevList.ElementsAs(ctx, &prevModels, false)...)
		if diags.HasError() {
			return types.ListNull(agsStorageMountObjectType()), diags
		}
	}

	out := make([]agsStorageMountModel, 0, len(mounts))
	for i, mount := range mounts {
		if mount == nil {
			continue
		}
		var prev agsStorageMountModel
		if i < len(prevModels) {
			prev = prevModels[i]
		}

		storageSource := types.ListNull(agsStorageSourceObjectType())
		if mount.StorageSource != nil {
			cos, cosDiags := flattenAgsCosStorageSource(ctx, mount.StorageSource.Cos, prev.StorageSource)
			diags.Append(cosDiags...)
			image, imageDiags := flattenAgsImageStorageSource(ctx, mount.StorageSource.Image, prev.StorageSource)
			diags.Append(imageDiags...)
			cfs, cfsDiags := flattenAgsCfsStorageSource(ctx, mount.StorageSource.Cfs, prev.StorageSource)
			diags.Append(cfsDiags...)
			agentBucket, agentDiags := flattenAgsAgentBucketStorageSource(ctx, mount.StorageSource.AgentBucket, prev.StorageSource)
			diags.Append(agentDiags...)
			if diags.HasError() {
				return types.ListNull(agsStorageMountObjectType()), diags
			}
			sourceOut := []agsStorageSourceModel{{
				Cos:         cos,
				Image:       image,
				Cfs:         cfs,
				AgentBucket: agentBucket,
			}}
			storageSource, _ = types.ListValueFrom(ctx, agsStorageSourceObjectType(), sourceOut)
		}

		volume := types.ListNull(agsVolumeRefObjectType())
		if mount.Volume != nil {
			volume, _ = types.ListValueFrom(ctx, agsVolumeRefObjectType(), []agsVolumeRefModel{{
				VolumeName: types.StringPointerValue(mount.Volume.VolumeName),
				VolumeId:   types.StringPointerValue(mount.Volume.VolumeId),
			}})
		}

		volumeTemplate := types.ListNull(agsVolumeTemplateRefObjectType())
		if mount.VolumeTemplate != nil {
			volumeTemplate, _ = types.ListValueFrom(ctx, agsVolumeTemplateRefObjectType(), []agsVolumeTemplateRefModel{{
				VolumeTemplateName: fwhelper.StringValueOrPrev(mount.VolumeTemplate.VolumeTemplateName, types.StringNull()),
				ReuseKey:           types.StringPointerValue(mount.VolumeTemplate.ReuseKey),
			}})
		}

		out = append(out, agsStorageMountModel{
			Name:           fwhelper.StringValueOrPrev(mount.Name, prev.Name),
			StorageSource:  storageSource,
			MountPath:      fwhelper.StringValueOrPrev(mount.MountPath, prev.MountPath),
			ReadOnly:       agsBoolValueOrPrev(mount.ReadOnly, prev.ReadOnly),
			Volume:         volume,
			VolumeTemplate: volumeTemplate,
		})
	}
	list, listDiags := types.ListValueFrom(ctx, agsStorageMountObjectType(), out)
	diags.Append(listDiags...)
	return list, diags
}

func flattenAgsCosStorageSource(ctx context.Context, cos *agssdk.CosStorageSource, prevSourceList types.List) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics
	var prev agsCosStorageSourceModel
	if cos == nil {
		return types.ListNull(agsCosStorageSourceObjectType()), diags
	}
	if !prevSourceList.IsNull() && !prevSourceList.IsUnknown() {
		var sourceModels []agsStorageSourceModel
		if diags.Append(prevSourceList.ElementsAs(ctx, &sourceModels, false)...); !diags.HasError() && len(sourceModels) > 0 && !sourceModels[0].Cos.IsNull() && !sourceModels[0].Cos.IsUnknown() {
			var cosModels []agsCosStorageSourceModel
			diags.Append(sourceModels[0].Cos.ElementsAs(ctx, &cosModels, false)...)
			if !diags.HasError() && len(cosModels) > 0 {
				prev = cosModels[0]
			}
		}
	}
	out := []agsCosStorageSourceModel{{
		Endpoint:   fwhelper.StringValueOrPrev(cos.Endpoint, prev.Endpoint),
		BucketName: fwhelper.StringValueOrPrev(cos.BucketName, prev.BucketName),
		BucketPath: fwhelper.StringValueOrPrev(cos.BucketPath, prev.BucketPath),
	}}
	list, listDiags := types.ListValueFrom(ctx, agsCosStorageSourceObjectType(), out)
	diags.Append(listDiags...)
	return list, diags
}

func flattenAgsImageStorageSource(ctx context.Context, image *agssdk.ImageStorageSource, prevSourceList types.List) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics
	var prev agsImageStorageSourceModel
	if image == nil {
		return types.ListNull(agsImageStorageSourceObjectType()), diags
	}
	if !prevSourceList.IsNull() && !prevSourceList.IsUnknown() {
		var sourceModels []agsStorageSourceModel
		if diags.Append(prevSourceList.ElementsAs(ctx, &sourceModels, false)...); !diags.HasError() && len(sourceModels) > 0 && !sourceModels[0].Image.IsNull() && !sourceModels[0].Image.IsUnknown() {
			var imageModels []agsImageStorageSourceModel
			diags.Append(sourceModels[0].Image.ElementsAs(ctx, &imageModels, false)...)
			if !diags.HasError() && len(imageModels) > 0 {
				prev = imageModels[0]
			}
		}
	}
	out := []agsImageStorageSourceModel{{
		Reference:         fwhelper.StringValueOrPrev(image.Reference, prev.Reference),
		ImageRegistryType: fwhelper.StringValueOrPrev(image.ImageRegistryType, prev.ImageRegistryType),
		SubPath:           fwhelper.StringValueOrPrev(image.SubPath, prev.SubPath),
		Digest:            fwhelper.StringValueOrPrev(image.Digest, prev.Digest),
	}}
	list, listDiags := types.ListValueFrom(ctx, agsImageStorageSourceObjectType(), out)
	diags.Append(listDiags...)
	return list, diags
}

func flattenAgsCfsStorageSource(ctx context.Context, cfs *agssdk.CfsStorageSource, prevSourceList types.List) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics
	var prev agsCfsStorageSourceModel
	if cfs == nil {
		return types.ListNull(agsCfsStorageSourceObjectType()), diags
	}
	if !prevSourceList.IsNull() && !prevSourceList.IsUnknown() {
		var sourceModels []agsStorageSourceModel
		if diags.Append(prevSourceList.ElementsAs(ctx, &sourceModels, false)...); !diags.HasError() && len(sourceModels) > 0 && !sourceModels[0].Cfs.IsNull() && !sourceModels[0].Cfs.IsUnknown() {
			var cfsModels []agsCfsStorageSourceModel
			diags.Append(sourceModels[0].Cfs.ElementsAs(ctx, &cfsModels, false)...)
			if !diags.HasError() && len(cfsModels) > 0 {
				prev = cfsModels[0]
			}
		}
	}
	out := []agsCfsStorageSourceModel{{
		FileSystemId: fwhelper.StringValueOrPrev(cfs.FileSystemId, prev.FileSystemId),
		Path:         fwhelper.StringValueOrPrev(cfs.Path, prev.Path),
	}}
	list, listDiags := types.ListValueFrom(ctx, agsCfsStorageSourceObjectType(), out)
	diags.Append(listDiags...)
	return list, diags
}

func flattenAgsAgentBucketStorageSource(ctx context.Context, agentBucket *agssdk.AgentBucketStorageSource, prevSourceList types.List) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics
	var prev agsAgentBucketStorageSourceModel
	if agentBucket == nil {
		return types.ListNull(agsAgentBucketStorageSourceObjectType()), diags
	}
	if !prevSourceList.IsNull() && !prevSourceList.IsUnknown() {
		var sourceModels []agsStorageSourceModel
		if diags.Append(prevSourceList.ElementsAs(ctx, &sourceModels, false)...); !diags.HasError() && len(sourceModels) > 0 && !sourceModels[0].AgentBucket.IsNull() && !sourceModels[0].AgentBucket.IsUnknown() {
			var agentModels []agsAgentBucketStorageSourceModel
			diags.Append(sourceModels[0].AgentBucket.ElementsAs(ctx, &agentModels, false)...)
			if !diags.HasError() && len(agentModels) > 0 {
				prev = agentModels[0]
			}
		}
	}
	out := []agsAgentBucketStorageSourceModel{{
		LibraryId:    fwhelper.StringValueOrPrev(agentBucket.LibraryId, prev.LibraryId),
		SpaceId:      fwhelper.StringValueOrPrev(agentBucket.SpaceId, prev.SpaceId),
		AccessDomain: fwhelper.StringValueOrPrev(agentBucket.AccessDomain, prev.AccessDomain),
	}}
	list, listDiags := types.ListValueFrom(ctx, agsAgentBucketStorageSourceObjectType(), out)
	diags.Append(listDiags...)
	return list, diags
}

func flattenAgsCustomConfiguration(ctx context.Context, config *agssdk.CustomConfigurationDetail, prevList types.List) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics
	if config == nil {
		if prevList.IsNull() || prevList.IsUnknown() {
			return types.ListNull(agsCustomConfigurationObjectType()), diags
		}
		return prevList, diags
	}

	var prev agsCustomConfigurationModel
	if !prevList.IsNull() && !prevList.IsUnknown() {
		var prevModels []agsCustomConfigurationModel
		diags.Append(prevList.ElementsAs(ctx, &prevModels, false)...)
		if diags.HasError() {
			return types.ListNull(agsCustomConfigurationObjectType()), diags
		}
		if len(prevModels) > 0 {
			prev = prevModels[0]
		}
	}

	command, commandDiags := agsFlattenStringSet(ctx, config.Command, prev.Command)
	diags.Append(commandDiags...)
	args, argsDiags := agsFlattenStringSet(ctx, config.Args, prev.Args)
	diags.Append(argsDiags...)
	if diags.HasError() {
		return types.ListNull(agsCustomConfigurationObjectType()), diags
	}

	env := types.ListNull(agsEnvVarObjectType())
	if len(config.Env) > 0 {
		envOut := make([]agsEnvVarModel, 0, len(config.Env))
		for _, item := range config.Env {
			if item == nil {
				continue
			}
			envOut = append(envOut, agsEnvVarModel{
				Name:  types.StringPointerValue(item.Name),
				Value: types.StringPointerValue(item.Value),
			})
		}
		env, _ = types.ListValueFrom(ctx, agsEnvVarObjectType(), envOut)
	} else if !prev.Env.IsNull() && !prev.Env.IsUnknown() {
		env = prev.Env
	}

	ports := types.ListNull(agsPortConfigurationObjectType())
	if len(config.Ports) > 0 {
		portOut := make([]agsPortConfigurationModel, 0, len(config.Ports))
		for _, item := range config.Ports {
			if item == nil {
				continue
			}
			portOut = append(portOut, agsPortConfigurationModel{
				Name:     types.StringPointerValue(item.Name),
				Port:     agsInt64ValueOrPrev(item.Port, types.Int64Null()),
				Protocol: types.StringPointerValue(item.Protocol),
			})
		}
		ports, _ = types.ListValueFrom(ctx, agsPortConfigurationObjectType(), portOut)
	} else if !prev.Ports.IsNull() && !prev.Ports.IsUnknown() {
		ports = prev.Ports
	}

	resources := types.ListNull(agsResourceConfigurationObjectType())
	if config.Resources != nil {
		resources, _ = types.ListValueFrom(ctx, agsResourceConfigurationObjectType(), []agsResourceConfigurationModel{{
			Cpu:                fwhelper.StringValueOrPrev(config.Resources.CPU, types.StringNull()),
			Memory:             fwhelper.StringValueOrPrev(config.Resources.Memory, types.StringNull()),
			Storage:            fwhelper.StringValueOrPrev(config.Resources.Storage, types.StringNull()),
			GpuResourceGroupId: fwhelper.StringValueOrPrev(config.Resources.GPUResourceGroupId, types.StringNull()),
		}})
	} else if !prev.Resources.IsNull() && !prev.Resources.IsUnknown() {
		resources = prev.Resources
	}

	probe := types.ListNull(agsProbeConfigurationObjectType())
	if config.Probe != nil {
		httpGet := types.ListNull(agsHttpGetActionObjectType())
		if config.Probe.HttpGet != nil {
			httpGet, _ = types.ListValueFrom(ctx, agsHttpGetActionObjectType(), []agsHttpGetActionModel{{
				Path:   types.StringPointerValue(config.Probe.HttpGet.Path),
				Port:   agsInt64ValueOrPrev(config.Probe.HttpGet.Port, types.Int64Null()),
				Scheme: types.StringPointerValue(config.Probe.HttpGet.Scheme),
			}})
		} else if !prev.Probe.IsNull() && !prev.Probe.IsUnknown() {
			var prevProbeModels []agsProbeConfigurationModel
			if diags.Append(prev.Probe.ElementsAs(ctx, &prevProbeModels, false)...); !diags.HasError() && len(prevProbeModels) > 0 {
				httpGet = prevProbeModels[0].HttpGet
			}
		}
		probe, _ = types.ListValueFrom(ctx, agsProbeConfigurationObjectType(), []agsProbeConfigurationModel{{
			HttpGet:          httpGet,
			ReadyTimeoutMs:   agsInt64ValueOrPrev(config.Probe.ReadyTimeoutMs, types.Int64Null()),
			ProbeTimeoutMs:   agsInt64ValueOrPrev(config.Probe.ProbeTimeoutMs, types.Int64Null()),
			ProbePeriodMs:    agsInt64ValueOrPrev(config.Probe.ProbePeriodMs, types.Int64Null()),
			SuccessThreshold: agsInt64ValueOrPrev(config.Probe.SuccessThreshold, types.Int64Null()),
			FailureThreshold: agsInt64ValueOrPrev(config.Probe.FailureThreshold, types.Int64Null()),
		}})
	} else if !prev.Probe.IsNull() && !prev.Probe.IsUnknown() {
		probe = prev.Probe
	}

	dnsConfig := types.ListNull(agsDNSConfigObjectType())
	if config.DNSConfig != nil {
		var prevDns agsDNSConfigModel
		if !prev.DNSConfig.IsNull() && !prev.DNSConfig.IsUnknown() {
			var prevDnsModels []agsDNSConfigModel
			if diags.Append(prev.DNSConfig.ElementsAs(ctx, &prevDnsModels, false)...); !diags.HasError() && len(prevDnsModels) > 0 {
				prevDns = prevDnsModels[0]
			}
		}
		servers, serversDiags := agsFlattenStringSet(ctx, config.DNSConfig.Servers, prevDns.Servers)
		diags.Append(serversDiags...)
		searches, searchesDiags := agsFlattenStringSet(ctx, config.DNSConfig.Searches, prevDns.Searches)
		diags.Append(searchesDiags...)
		options, optionsDiags := agsFlattenStringSet(ctx, config.DNSConfig.Options, prevDns.Options)
		diags.Append(optionsDiags...)
		if diags.HasError() {
			return types.ListNull(agsCustomConfigurationObjectType()), diags
		}
		dnsConfig, _ = types.ListValueFrom(ctx, agsDNSConfigObjectType(), []agsDNSConfigModel{{
			Servers:  servers,
			Searches: searches,
			Options:  options,
		}})
	} else if !prev.DNSConfig.IsNull() && !prev.DNSConfig.IsUnknown() {
		dnsConfig = prev.DNSConfig
	}

	securityContext := types.ListNull(agsSecurityContextObjectType())
	if config.SecurityContext != nil {
		var prevSecurity agsSecurityContextModel
		if !prev.SecurityContext.IsNull() && !prev.SecurityContext.IsUnknown() {
			var prevSecurityModels []agsSecurityContextModel
			if diags.Append(prev.SecurityContext.ElementsAs(ctx, &prevSecurityModels, false)...); !diags.HasError() && len(prevSecurityModels) > 0 {
				prevSecurity = prevSecurityModels[0]
			}
		}
		capabilities := types.ListNull(agsCapabilitiesObjectType())
		if config.SecurityContext.Capabilities != nil {
			var prevCaps agsCapabilitiesModel
			if !prevSecurity.Capabilities.IsNull() && !prevSecurity.Capabilities.IsUnknown() {
				var prevCapModels []agsCapabilitiesModel
				if diags.Append(prevSecurity.Capabilities.ElementsAs(ctx, &prevCapModels, false)...); !diags.HasError() && len(prevCapModels) > 0 {
					prevCaps = prevCapModels[0]
				}
			}
			add, addDiags := agsFlattenStringSet(ctx, config.SecurityContext.Capabilities.Add, prevCaps.Add)
			diags.Append(addDiags...)
			drop, dropDiags := agsFlattenStringSet(ctx, config.SecurityContext.Capabilities.Drop, prevCaps.Drop)
			diags.Append(dropDiags...)
			if diags.HasError() {
				return types.ListNull(agsCustomConfigurationObjectType()), diags
			}
			capabilities, _ = types.ListValueFrom(ctx, agsCapabilitiesObjectType(), []agsCapabilitiesModel{{
				Add:  add,
				Drop: drop,
			}})
		} else if !prevSecurity.Capabilities.IsNull() && !prevSecurity.Capabilities.IsUnknown() {
			capabilities = prevSecurity.Capabilities
		}
		securityContext, _ = types.ListValueFrom(ctx, agsSecurityContextObjectType(), []agsSecurityContextModel{{
			Privileged:   agsBoolValueOrPrev(config.SecurityContext.Privileged, prevSecurity.Privileged),
			Capabilities: capabilities,
		}})
	} else if !prev.SecurityContext.IsNull() && !prev.SecurityContext.IsUnknown() {
		securityContext = prev.SecurityContext
	}

	out := []agsCustomConfigurationModel{{
		Image:             fwhelper.StringValueOrPrev(config.Image, prev.Image),
		ImageRegistryType: fwhelper.StringValueOrPrev(config.ImageRegistryType, prev.ImageRegistryType),
		Command:           command,
		Args:              args,
		Env:               env,
		Ports:             ports,
		Resources:         resources,
		Probe:             probe,
		DNSConfig:         dnsConfig,
		SecurityContext:   securityContext,
	}}
	list, listDiags := types.ListValueFrom(ctx, agsCustomConfigurationObjectType(), out)
	diags.Append(listDiags...)
	return list, diags
}

func flattenAgsComputerConfiguration(ctx context.Context, config *agssdk.ComputerConfiguration, prevList types.List) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics
	if config == nil {
		if prevList.IsNull() || prevList.IsUnknown() {
			return types.ListNull(agsComputerConfigurationObjectType()), diags
		}
		return prevList, diags
	}

	var prev agsComputerConfigurationModel
	if !prevList.IsNull() && !prevList.IsUnknown() {
		var prevModels []agsComputerConfigurationModel
		diags.Append(prevList.ElementsAs(ctx, &prevModels, false)...)
		if diags.HasError() {
			return types.ListNull(agsComputerConfigurationObjectType()), diags
		}
		if len(prevModels) > 0 {
			prev = prevModels[0]
		}
	}

	waaConfiguration := types.ListNull(agsWaaConfigurationObjectType())
	if config.WAAConfiguration != nil {
		waaConfiguration, _ = types.ListValueFrom(ctx, agsWaaConfigurationObjectType(), []agsWaaConfigurationModel{{
			ImageId: fwhelper.StringValueOrPrev(config.WAAConfiguration.ImageId, types.StringNull()),
		}})
	} else if !prev.WaaConfiguration.IsNull() && !prev.WaaConfiguration.IsUnknown() {
		waaConfiguration = prev.WaaConfiguration
	}

	osWorldConfiguration := types.ListNull(agsOsWorldConfigurationObjectType())
	if config.OSWorldConfiguration != nil {
		osWorldConfiguration, _ = types.ListValueFrom(ctx, agsOsWorldConfigurationObjectType(), []agsOsWorldConfigurationModel{{
			Version: fwhelper.StringValueOrPrev(config.OSWorldConfiguration.Version, types.StringNull()),
		}})
	} else if !prev.OsWorldConfiguration.IsNull() && !prev.OsWorldConfiguration.IsUnknown() {
		osWorldConfiguration = prev.OsWorldConfiguration
	}

	out := []agsComputerConfigurationModel{{
		WaaConfiguration:     waaConfiguration,
		OsWorldConfiguration: osWorldConfiguration,
	}}
	list, listDiags := types.ListValueFrom(ctx, agsComputerConfigurationObjectType(), out)
	diags.Append(listDiags...)
	return list, diags
}

func flattenAgsLogConfiguration(ctx context.Context, config *agssdk.LogConfiguration, prevList types.List) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics
	if config == nil {
		if prevList.IsNull() || prevList.IsUnknown() {
			return types.ListNull(agsLogConfigurationObjectType()), diags
		}
		return prevList, diags
	}

	var prev agsLogConfigurationModel
	if !prevList.IsNull() && !prevList.IsUnknown() {
		var prevModels []agsLogConfigurationModel
		diags.Append(prevList.ElementsAs(ctx, &prevModels, false)...)
		if diags.HasError() {
			return types.ListNull(agsLogConfigurationObjectType()), diags
		}
		if len(prevModels) > 0 {
			prev = prevModels[0]
		}
	}

	clsConfig := types.ListNull(agsClsConfigObjectType())
	if config.CLSConfig != nil {
		clsConfig, _ = types.ListValueFrom(ctx, agsClsConfigObjectType(), []agsClsConfigModel{{
			TopicId: fwhelper.StringValueOrPrev(config.CLSConfig.TopicId, types.StringNull()),
		}})
	} else if !prev.ClsConfig.IsNull() && !prev.ClsConfig.IsUnknown() {
		clsConfig = prev.ClsConfig
	}

	logSources := types.ListNull(agsLogSourcesObjectType())
	if config.LogSources != nil {
		var prevSources agsLogSourcesModel
		if !prev.LogSources.IsNull() && !prev.LogSources.IsUnknown() {
			var prevSourceModels []agsLogSourcesModel
			if diags.Append(prev.LogSources.ElementsAs(ctx, &prevSourceModels, false)...); !diags.HasError() && len(prevSourceModels) > 0 {
				prevSources = prevSourceModels[0]
			}
		}
		files, filesDiags := agsFlattenStringSet(ctx, config.LogSources.Files, prevSources.Files)
		diags.Append(filesDiags...)
		if diags.HasError() {
			return types.ListNull(agsLogConfigurationObjectType()), diags
		}
		logSources, _ = types.ListValueFrom(ctx, agsLogSourcesObjectType(), []agsLogSourcesModel{{
			Files: files,
		}})
	} else if !prev.LogSources.IsNull() && !prev.LogSources.IsUnknown() {
		logSources = prev.LogSources
	}

	out := []agsLogConfigurationModel{{
		ClsConfig:  clsConfig,
		LogSources: logSources,
	}}
	list, listDiags := types.ListValueFrom(ctx, agsLogConfigurationObjectType(), out)
	diags.Append(listDiags...)
	return list, diags
}

// agsFlattenStringSet converts a *[]string into a Terraform string set; when
// the cloud value is nil the previous state set is kept so a configured value
// is not cleared on refresh.
func agsFlattenStringSet(ctx context.Context, values []*string, prev types.Set) (types.Set, diag.Diagnostics) {
	var diags diag.Diagnostics
	if values == nil {
		if prev.IsNull() || prev.IsUnknown() {
			return types.SetNull(types.StringType), diags
		}
		return prev, diags
	}
	set, setDiags := types.SetValueFrom(ctx, types.StringType, values)
	diags.Append(setDiags...)
	return set, diags
}

func agsBoolValueOrPrev(value *bool, prev types.Bool) types.Bool {
	if value == nil {
		return prev
	}
	return types.BoolValue(*value)
}

func agsInt64ValueOrPrev(value *int64, prev types.Int64) types.Int64 {
	if value == nil {
		return prev
	}
	return types.Int64Value(*value)
}

// ---- attribute type helpers used by the flatten functions ----

func agsNetworkConfigurationObjectType() attr.Type {
	return types.ObjectType{AttrTypes: agsNetworkConfigurationAttrTypes()}
}

func agsNetworkConfigurationAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"network_mode": types.StringType,
		"vpc_config":   types.ListType{ElemType: agsVpcConfigObjectType()},
	}
}

func agsVpcConfigObjectType() attr.Type {
	return types.ObjectType{AttrTypes: agsVpcConfigAttrTypes()}
}

func agsVpcConfigAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"subnet_ids":         types.SetType{ElemType: types.StringType},
		"security_group_ids": types.SetType{ElemType: types.StringType},
	}
}

func agsTagObjectType() attr.Type {
	return types.ObjectType{AttrTypes: agsTagAttrTypes()}
}

func agsTagAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"key":   types.StringType,
		"value": types.StringType,
	}
}

func agsStorageMountObjectType() attr.Type {
	return types.ObjectType{AttrTypes: agsStorageMountAttrTypes()}
}

func agsStorageMountAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"name":            types.StringType,
		"storage_source":  types.ListType{ElemType: agsStorageSourceObjectType()},
		"mount_path":      types.StringType,
		"read_only":       types.BoolType,
		"volume":          types.ListType{ElemType: agsVolumeRefObjectType()},
		"volume_template": types.ListType{ElemType: agsVolumeTemplateRefObjectType()},
	}
}

func agsStorageSourceObjectType() attr.Type {
	return types.ObjectType{AttrTypes: agsStorageSourceAttrTypes()}
}

func agsStorageSourceAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"cos":          types.ListType{ElemType: agsCosStorageSourceObjectType()},
		"image":        types.ListType{ElemType: agsImageStorageSourceObjectType()},
		"cfs":          types.ListType{ElemType: agsCfsStorageSourceObjectType()},
		"agent_bucket": types.ListType{ElemType: agsAgentBucketStorageSourceObjectType()},
	}
}

func agsCosStorageSourceObjectType() attr.Type {
	return types.ObjectType{AttrTypes: map[string]attr.Type{
		"endpoint":    types.StringType,
		"bucket_name": types.StringType,
		"bucket_path": types.StringType,
	}}
}

func agsImageStorageSourceObjectType() attr.Type {
	return types.ObjectType{AttrTypes: map[string]attr.Type{
		"reference":           types.StringType,
		"image_registry_type": types.StringType,
		"sub_path":            types.StringType,
		"digest":              types.StringType,
	}}
}

func agsCfsStorageSourceObjectType() attr.Type {
	return types.ObjectType{AttrTypes: map[string]attr.Type{
		"file_system_id": types.StringType,
		"path":           types.StringType,
	}}
}

func agsAgentBucketStorageSourceObjectType() attr.Type {
	return types.ObjectType{AttrTypes: map[string]attr.Type{
		"library_id":    types.StringType,
		"space_id":      types.StringType,
		"access_domain": types.StringType,
	}}
}

func agsVolumeRefObjectType() attr.Type {
	return types.ObjectType{AttrTypes: map[string]attr.Type{
		"volume_name": types.StringType,
		"volume_id":   types.StringType,
	}}
}

func agsVolumeTemplateRefObjectType() attr.Type {
	return types.ObjectType{AttrTypes: map[string]attr.Type{
		"volume_template_name": types.StringType,
		"reuse_key":            types.StringType,
	}}
}

func agsCustomConfigurationObjectType() attr.Type {
	return types.ObjectType{AttrTypes: agsCustomConfigurationAttrTypes()}
}

func agsCustomConfigurationAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"image":               types.StringType,
		"image_registry_type": types.StringType,
		"command":             types.SetType{ElemType: types.StringType},
		"args":                types.SetType{ElemType: types.StringType},
		"env":                 types.ListType{ElemType: agsEnvVarObjectType()},
		"ports":               types.ListType{ElemType: agsPortConfigurationObjectType()},
		"resources":           types.ListType{ElemType: agsResourceConfigurationObjectType()},
		"probe":               types.ListType{ElemType: agsProbeConfigurationObjectType()},
		"dns_config":          types.ListType{ElemType: agsDNSConfigObjectType()},
		"security_context":    types.ListType{ElemType: agsSecurityContextObjectType()},
	}
}

func agsEnvVarObjectType() attr.Type {
	return types.ObjectType{AttrTypes: map[string]attr.Type{
		"name":  types.StringType,
		"value": types.StringType,
	}}
}

func agsPortConfigurationObjectType() attr.Type {
	return types.ObjectType{AttrTypes: map[string]attr.Type{
		"name":     types.StringType,
		"port":     types.Int64Type,
		"protocol": types.StringType,
	}}
}

func agsResourceConfigurationObjectType() attr.Type {
	return types.ObjectType{AttrTypes: map[string]attr.Type{
		"cpu":                   types.StringType,
		"memory":                types.StringType,
		"storage":               types.StringType,
		"gpu_resource_group_id": types.StringType,
	}}
}

func agsProbeConfigurationObjectType() attr.Type {
	return types.ObjectType{AttrTypes: agsProbeConfigurationAttrTypes()}
}

func agsProbeConfigurationAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"http_get":          types.ListType{ElemType: agsHttpGetActionObjectType()},
		"ready_timeout_ms":  types.Int64Type,
		"probe_timeout_ms":  types.Int64Type,
		"probe_period_ms":   types.Int64Type,
		"success_threshold": types.Int64Type,
		"failure_threshold": types.Int64Type,
	}
}

func agsHttpGetActionObjectType() attr.Type {
	return types.ObjectType{AttrTypes: map[string]attr.Type{
		"path":   types.StringType,
		"port":   types.Int64Type,
		"scheme": types.StringType,
	}}
}

func agsDNSConfigObjectType() attr.Type {
	return types.ObjectType{AttrTypes: map[string]attr.Type{
		"servers":  types.SetType{ElemType: types.StringType},
		"searches": types.SetType{ElemType: types.StringType},
		"options":  types.SetType{ElemType: types.StringType},
	}}
}

func agsSecurityContextObjectType() attr.Type {
	return types.ObjectType{AttrTypes: map[string]attr.Type{
		"privileged":   types.BoolType,
		"capabilities": types.ListType{ElemType: agsCapabilitiesObjectType()},
	}}
}

func agsCapabilitiesObjectType() attr.Type {
	return types.ObjectType{AttrTypes: map[string]attr.Type{
		"add":  types.SetType{ElemType: types.StringType},
		"drop": types.SetType{ElemType: types.StringType},
	}}
}

func agsComputerConfigurationObjectType() attr.Type {
	return types.ObjectType{AttrTypes: agsComputerConfigurationAttrTypes()}
}

func agsComputerConfigurationAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"waa_configuration":      types.ListType{ElemType: agsWaaConfigurationObjectType()},
		"os_world_configuration": types.ListType{ElemType: agsOsWorldConfigurationObjectType()},
	}
}

func agsWaaConfigurationObjectType() attr.Type {
	return types.ObjectType{AttrTypes: map[string]attr.Type{
		"image_id": types.StringType,
	}}
}

func agsOsWorldConfigurationObjectType() attr.Type {
	return types.ObjectType{AttrTypes: map[string]attr.Type{
		"version": types.StringType,
	}}
}

func agsLogConfigurationObjectType() attr.Type {
	return types.ObjectType{AttrTypes: agsLogConfigurationAttrTypes()}
}

func agsLogConfigurationAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"cls_config":  types.ListType{ElemType: agsClsConfigObjectType()},
		"log_sources": types.ListType{ElemType: agsLogSourcesObjectType()},
	}
}

func agsClsConfigObjectType() attr.Type {
	return types.ObjectType{AttrTypes: map[string]attr.Type{
		"topic_id": types.StringType,
	}}
}

func agsLogSourcesObjectType() attr.Type {
	return types.ObjectType{AttrTypes: map[string]attr.Type{
		"files": types.SetType{ElemType: types.StringType},
	}}
}

// networkConfigurationAttribute builds the network_configuration attribute;
// required is true on the top-level schema (Create requires it).
func networkConfigurationAttribute(required bool) schema.ListNestedAttribute {
	attr := schema.ListNestedAttribute{
		Optional:    !required,
		Required:    required,
		Description: "Network configuration.",
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"network_mode": schema.StringAttribute{
					Required:    true,
					Description: "Network mode. Valid values: `PUBLIC`, `VPC`, `SANDBOX`.",
				},
				"vpc_config": schema.ListNestedAttribute{
					Optional:    true,
					Description: "VPC network configuration.",
					NestedObject: schema.NestedAttributeObject{
						Attributes: map[string]schema.Attribute{
							"subnet_ids": schema.SetAttribute{
								Optional:    true,
								ElementType: types.StringType,
								Description: "VPC subnet ID list.",
							},
							"security_group_ids": schema.SetAttribute{
								Optional:    true,
								ElementType: types.StringType,
								Description: "Security group ID list.",
							},
						},
					},
				},
			},
		},
	}
	return attr
}
