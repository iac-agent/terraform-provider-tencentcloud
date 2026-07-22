package tsf

import (
	"context"
	"fmt"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tsf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tsf/v20180326"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudTsfDeployContainerGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTsfDeployContainerGroupCreate,
		Read:   resourceTencentCloudTsfDeployContainerGroupRead,
		Update: resourceTencentCloudTsfDeployContainerGroupUpdate,
		Delete: resourceTencentCloudTsfDeployContainerGroupDelete,

		Schema: map[string]*schema.Schema{
			"group_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "组 ID。",
			},

			"tag_name": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "镜像 版本 名称，v1。",
			},

			"instance_num": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "实例 数量。",
			},

			"server": {
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "镜像 服务器。",
			},

			"reponame": {
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "old 镜像 名称，eg: /tsf/服务器。",
			},

			"cpu_limit": {
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "最大CPU cores 对于 business 容器，corresponding 到 限制 在 K8S. 如果未指定，它 默认为 twice 请求。",
			},

			"mem_limit": {
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "最大 内存 大小 在 MiB 对于 business 容器，corresponding 到 限制 在 K8S. 如果未指定，它 默认为 twice 请求。",
			},

			"jvm_opts": {
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "jvm options。",
			},

			"cpu_request": {
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "数量 CPU 核数 allocated 到 business 容器，corresponding 到 请求 在 K8S. 默认值为 0.25。",
			},

			"mem_request": {
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "amount 的 内存 在 MiB allocated 到 business 容器，corresponding 到 请求 在 K8S. 默认值为 640 MiB。",
			},

			"do_not_start": {
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeBool,
				Description: "Not start right away。",
			},

			"repo_name": {
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "(优先级 使用) New 镜像 名称，such 作为 /tsf/nginx。",
			},

			"update_type": {
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "Update 方法: 0 对于 fast update，1 对于 rolling update。",
			},

			"update_ivl": {
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "update Interval，为必填项 当 rolling update。",
			},

			"agent_cpu_request": {
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "数量 CPU 核数 allocated 到 agent 容器 corresponds 到 请求 字段 在 Kubernetes。",
			},

			"agent_cpu_limit": {
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "最大CPU cores allocated 到 agent 容器 corresponds 到 限制 字段 在 Kubernetes。",
			},

			"agent_mem_request": {
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "amount 的 内存 在 MiB allocated 到 agent 容器 corresponds 到 请求 字段 在 Kubernetes。",
			},

			"agent_mem_limit": {
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "最大 amount 的 内存 在 MiB allocated 到 agent 容器 corresponds 到 &amp;#39;限制&amp;#39; 字段 在 Kubernetes。",
			},

			"istio_cpu_request": {
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "数量 CPU 核数 allocated 到 istio proxy 容器 corresponds 到 &amp;#39;请求&amp;#39; 字段 在 Kubernetes。",
			},

			"istio_cpu_limit": {
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "最大 amount 的 CPU 核数 allocated 到 istio proxy 容器 corresponds 到 &amp;#39;限制&amp;#39; 字段 在 Kubernetes。",
			},

			"istio_mem_request": {
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "amount 的 内存 在 MiB allocated 到 agent 容器 corresponds 到 请求 字段 在 Kubernetes。",
			},

			"istio_mem_limit": {
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "最大 amount 的 内存 在 MiB allocated 到 agent 容器 corresponds 到 请求 字段 在 Kubernetes。",
			},

			"max_surge": {
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "MaxSurge 参数 在 Kubernetes rolling update strategy。",
			},

			"max_unavailable": {
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "MaxUnavailable 参数 在 Kubernetes rolling update strategy。",
			},

			"health_check_settings": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "配置 信息 对于 health checks. 如果 此 参数 是 不 指定， health check 是 不 集合 通过 默认值。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"liveness_probe": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Liveness probe. 注意：此字段可能返回 null，表示无法获取有效值。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"action_type": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "health check 方法. HTTP: checks through HTTP interface; CMD: checks 通过 executing command; TCP: checks 通过 establishing TCP 连接. 注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"initial_delay_seconds": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "时间 延迟 对于 容器 到 start health check. 注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"timeout_seconds": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "最大 超时 周期 对于 each health check response. 注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"period_seconds": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "时间间隔 对于 performing health checks. 注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"success_threshold": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "数量 consecutive successful health checks 必填 对于 backend 容器 到 transition 从 failure 到 success. 注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"failure_threshold": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "数量 consecutive successful health checks 必填 对于 backend 容器 到 transition 从 success 到 failure. 注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"scheme": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "协议 用于HTTP health checks. HTTP 和 HTTPS 是 支持. 注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"port": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "端口 用于health checks，ranging 从 1 到 65535. 注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"path": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "请求 路径 对于 HTTP health checks. 注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"command": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Optional:    true,
										Description: "command 到 是 executed 对于 command health checks. 注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"type": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "类型 readiness probe. TSF_DEFAULT 表示 默认值 readiness probe 的 TSF，while K8S_NATIVE 表示 native readiness probe 的 Kubernetes. 如果 此 字段 是 不 指定， native readiness probe 的 Kubernetes 是 使用 通过 默认值. 注意：此字段可能返回 null，表示无法获取有效值。",
									},
								},
							},
						},
						"readiness_probe": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Readiness health check. 注意：此字段可能返回 null，表示无法获取有效值。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"action_type": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "health check 方法. HTTP 表示checking through HTTP interface，CMD 表示checking through executing command，和 TCP 表示checking through establishing TCP 连接. 注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"initial_delay_seconds": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "时间 到 延迟 start 的 容器 health check. 注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"timeout_seconds": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "最大 超时 周期 对于 each health check response. 注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"period_seconds": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "时间间隔 对于 performing health checks. 注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"success_threshold": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "数量 consecutive successful health checks 必填 对于 backend 容器 到 transition 从 failure 到 success. 注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"failure_threshold": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "数量 consecutive successful health checks 必填 对于 backend 容器 到 transition 从 success 到 failure. 注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"scheme": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "协议 用于HTTP health checks. HTTP 和 HTTPS 是 支持. 注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"port": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "端口 用于health checks，ranging 从 1 到 65535. 注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"path": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "请求 路径 对于 HTTP health checks. 注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"command": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Optional:    true,
										Description: "command 到 是 executed 对于 command check. 注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"type": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "类型 readiness probe. TSF_DEFAULT 表示 默认值 readiness probe 的 TSF，while K8S_NATIVE 表示 native readiness probe 的 Kubernetes. 如果 此 字段 是 不 指定， native readiness probe 的 Kubernetes 是 使用 通过 默认值. 注意：此字段可能返回 null，表示无法获取有效值。",
									},
								},
							},
						},
					},
				},
			},

			"envs": {
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				Description: "环境 variables 该 应用 runs 在 部署 组. 如果 此 参数 是 不 指定，无 additional 环境 variables 是 集合 通过 默认值。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "env param 名称",
						},
						"value": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "值 的 env。",
						},
						"value_from": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Computed:    true,
							Description: "Kubernetes ValueFrom 配置. 注意：此字段可能返回 null，表示无法获取有效值。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"field_ref": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Computed:    true,
										Description: "FieldRef 配置 的 Kubernetes env. 注意：此字段可能返回 null，表示无法获取有效值。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"field_path": {
													Type:        schema.TypeString,
													Optional:    true,
													Computed:    true,
													Description: "FieldPath 配置 的 Kubernetes. 注意：此字段可能返回 null，表示无法获取有效值。",
												},
											},
										},
									},
									"resource_field_ref": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Computed:    true,
										Description: "ResourceFieldRef 配置 的 Kubernetes env. 注意：此字段可能返回 null，表示无法获取有效值。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"resource": {
													Type:        schema.TypeString,
													Optional:    true,
													Computed:    true,
													Description: "Resource 配置 的 Kubernetes. 注意：此字段可能返回 null，表示无法获取有效值。",
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

			"service_setting": {
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Network settings 对于 容器 部署 groups。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"access_type": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "0: Public 网络，1: Access within 集群，2: NodePort，3: Access within VPC. 注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"protocol_ports": {
							Type:        schema.TypeList,
							Required:    true,
							Description: "Container 端口 mapping. 注意：此字段可能返回 null，表示无法获取有效值。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"protocol": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "TCP 或 UDP。",
									},
									"port": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "端口",
									},
									"target_port": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "容器 端口",
									},
									"node_port": {
										Type:        schema.TypeInt,
										Optional:    true,
										Computed:    true,
										Description: "节点 端口",
									},
								},
							},
						},
						"subnet_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "子网 ID。",
						},
						"disable_service": {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: "是否create Kubernetes 服务. 默认值为 false. 注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"headless_service": {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: "是否service 是 的 headless 类型 注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"allow_delete_service": {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: "当 集合 到 true 和 DisableService 是 also true， previously 创建 服务 将 是 删除. Please 使用 使用 caution. 注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"open_session_affinity": {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: "Enable 会话 affinity. true 表示 已启用，false 表示 已禁用 默认值为 false. 注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"session_affinity_timeout_seconds": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "Session affinity 会话 时间. 默认值为 10800. 注意：此字段可能返回 null，表示无法获取有效值。",
						},
					},
				},
			},

			"deploy_agent": {
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeBool,
				Description: "是否deploy agent 容器. 如果 此 参数 是 不 指定， agent 容器 将 不 是 deployed 通过 默认值。",
			},

			"scheduling_strategy": {
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Node scheduling strategy. 如果 此 参数 是 不 指定， 节点 scheduling strategy 将 不 是 使用 通过 默认值。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "NONE: Do 不 使用 scheduling strategy; CROSS_AZ: Deploy across availability zones. 注意：此字段可能返回 null，表示无法获取有效值。",
						},
					},
				},
			},

			"incremental_deployment": {
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeBool,
				Description: "是否perform incremental 部署. 默认值为 false，其中 表示 full update。",
			},

			"repo_type": {
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "repo 类型，tcr 或 leave 它 blank。",
			},

			"volume_info_list": {
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				Description: "Volume 信息，作为 列表。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"volume_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "卷 类型",
						},
						"volume_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "卷 名称",
						},
						"volume_config": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "卷 配置",
						},
					},
				},
			},

			"volume_mount_info_list": {
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				Description: "Volume mount point 信息，列表 类型",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"volume_mount_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "mount 卷 名称",
						},
						"volume_mount_path": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "mount 路径",
						},
						"volume_mount_sub_path": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "mount subPath。",
						},
						"read_or_write": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Read 和 write 访问 模式 1: Read-仅. 2: Read-write。",
						},
					},
				},
			},

			"volume_clean": {
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeBool,
				Description: "是否clear 卷 信息. 默认为 false。",
			},

			"agent_profile_list": {
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				Description: "javaagent info: SERVICE_AGENT/OT_AGENT。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"agent_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Agent 类型",
						},
						"agent_version": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Agent 版本",
						},
					},
				},
			},

			"warmup_setting": {
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "warmup setting。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: "是否enable preheating。",
						},
						"warmup_time": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "warmup 时间。",
						},
						"curvature": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "Preheating curvature，使用 值 between 1 和 5。",
						},
						"enabled_protection": {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: "是否enable preheating protection. 如果 protection 是 已启用 和 more 比 50% 的 nodes 是 在 preheating state，preheating 将 是 aborted。",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudTsfDeployContainerGroupCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tsf_deploy_container_group.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	var (
		request  = tsf.NewDeployContainerGroupRequest()
		response = tsf.NewDeployContainerGroupResponse()
		groupId  string
	)
	if v, ok := d.GetOk("group_id"); ok {
		groupId = v.(string)
		request.GroupId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("tag_name"); ok {
		request.TagName = helper.String(v.(string))
	}

	if v, _ := d.GetOk("instance_num"); v != nil {
		request.InstanceNum = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("server"); ok {
		request.Server = helper.String(v.(string))
	}

	if v, ok := d.GetOk("reponame"); ok {
		request.Reponame = helper.String(v.(string))
	}

	if v, ok := d.GetOk("cpu_limit"); ok {
		request.CpuLimit = helper.String(v.(string))
	}

	if v, ok := d.GetOk("mem_limit"); ok {
		request.MemLimit = helper.String(v.(string))
	}

	if v, ok := d.GetOk("jvm_opts"); ok {
		request.JvmOpts = helper.String(v.(string))
	}

	if v, ok := d.GetOk("cpu_request"); ok {
		request.CpuRequest = helper.String(v.(string))
	}

	if v, ok := d.GetOk("mem_request"); ok {
		request.MemRequest = helper.String(v.(string))
	}

	if v, _ := d.GetOk("do_not_start"); v != nil {
		request.DoNotStart = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOk("repo_name"); ok {
		request.RepoName = helper.String(v.(string))
	}

	if v, _ := d.GetOk("update_type"); v != nil {
		request.UpdateType = helper.IntInt64(v.(int))
	}

	if v, _ := d.GetOk("update_ivl"); v != nil {
		request.UpdateIvl = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("agent_cpu_request"); ok {
		request.AgentCpuRequest = helper.String(v.(string))
	}

	if v, ok := d.GetOk("agent_cpu_limit"); ok {
		request.AgentCpuLimit = helper.String(v.(string))
	}

	if v, ok := d.GetOk("agent_mem_request"); ok {
		request.AgentMemRequest = helper.String(v.(string))
	}

	if v, ok := d.GetOk("agent_mem_limit"); ok {
		request.AgentMemLimit = helper.String(v.(string))
	}

	if v, ok := d.GetOk("istio_cpu_request"); ok {
		request.IstioCpuRequest = helper.String(v.(string))
	}

	if v, ok := d.GetOk("istio_cpu_limit"); ok {
		request.IstioCpuLimit = helper.String(v.(string))
	}

	if v, ok := d.GetOk("istio_mem_request"); ok {
		request.IstioMemRequest = helper.String(v.(string))
	}

	if v, ok := d.GetOk("istio_mem_limit"); ok {
		request.IstioMemLimit = helper.String(v.(string))
	}

	if v, ok := d.GetOk("max_surge"); ok {
		request.MaxSurge = helper.String(v.(string))
	}

	if v, ok := d.GetOk("max_unavailable"); ok {
		request.MaxUnavailable = helper.String(v.(string))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "health_check_settings"); ok {
		healthCheckSettings := tsf.HealthCheckSettings{}
		if livenessProbeMap, ok := helper.InterfaceToMap(dMap, "liveness_probe"); ok {
			healthCheckSetting := tsf.HealthCheckSetting{}
			if v, ok := livenessProbeMap["action_type"]; ok {
				healthCheckSetting.ActionType = helper.String(v.(string))
			}
			if v, ok := livenessProbeMap["initial_delay_seconds"]; ok {
				healthCheckSetting.InitialDelaySeconds = helper.IntUint64(v.(int))
			}
			if v, ok := livenessProbeMap["timeout_seconds"]; ok {
				healthCheckSetting.TimeoutSeconds = helper.IntUint64(v.(int))
			}
			if v, ok := livenessProbeMap["period_seconds"]; ok {
				healthCheckSetting.PeriodSeconds = helper.IntUint64(v.(int))
			}
			if v, ok := livenessProbeMap["success_threshold"]; ok {
				healthCheckSetting.SuccessThreshold = helper.IntUint64(v.(int))
			}
			if v, ok := livenessProbeMap["failure_threshold"]; ok {
				healthCheckSetting.FailureThreshold = helper.IntUint64(v.(int))
			}
			if v, ok := livenessProbeMap["scheme"]; ok {
				healthCheckSetting.Scheme = helper.String(v.(string))
			}
			if v, ok := livenessProbeMap["port"]; ok {
				healthCheckSetting.Port = helper.IntUint64(v.(int))
			}
			if v, ok := livenessProbeMap["path"]; ok {
				healthCheckSetting.Path = helper.String(v.(string))
			}
			if v, ok := livenessProbeMap["command"]; ok {
				commandSet := v.(*schema.Set).List()
				for i := range commandSet {
					command := commandSet[i].(string)
					healthCheckSetting.Command = append(healthCheckSetting.Command, &command)
				}
			}
			if v, ok := livenessProbeMap["type"]; ok {
				healthCheckSetting.Type = helper.String(v.(string))
			}
			healthCheckSettings.LivenessProbe = &healthCheckSetting
		}
		if readinessProbeMap, ok := helper.InterfaceToMap(dMap, "readiness_probe"); ok {
			healthCheckSetting := tsf.HealthCheckSetting{}
			if v, ok := readinessProbeMap["action_type"]; ok {
				healthCheckSetting.ActionType = helper.String(v.(string))
			}
			if v, ok := readinessProbeMap["initial_delay_seconds"]; ok {
				healthCheckSetting.InitialDelaySeconds = helper.IntUint64(v.(int))
			}
			if v, ok := readinessProbeMap["timeout_seconds"]; ok {
				healthCheckSetting.TimeoutSeconds = helper.IntUint64(v.(int))
			}
			if v, ok := readinessProbeMap["period_seconds"]; ok {
				healthCheckSetting.PeriodSeconds = helper.IntUint64(v.(int))
			}
			if v, ok := readinessProbeMap["success_threshold"]; ok {
				healthCheckSetting.SuccessThreshold = helper.IntUint64(v.(int))
			}
			if v, ok := readinessProbeMap["failure_threshold"]; ok {
				healthCheckSetting.FailureThreshold = helper.IntUint64(v.(int))
			}
			if v, ok := readinessProbeMap["scheme"]; ok {
				healthCheckSetting.Scheme = helper.String(v.(string))
			}
			if v, ok := readinessProbeMap["port"]; ok {
				healthCheckSetting.Port = helper.IntUint64(v.(int))
			}
			if v, ok := readinessProbeMap["path"]; ok {
				healthCheckSetting.Path = helper.String(v.(string))
			}
			if v, ok := readinessProbeMap["command"]; ok {
				commandSet := v.(*schema.Set).List()
				for i := range commandSet {
					command := commandSet[i].(string)
					healthCheckSetting.Command = append(healthCheckSetting.Command, &command)
				}
			}
			if v, ok := readinessProbeMap["type"]; ok {
				healthCheckSetting.Type = helper.String(v.(string))
			}
			healthCheckSettings.ReadinessProbe = &healthCheckSetting
		}
		request.HealthCheckSettings = &healthCheckSettings
	}

	if v, ok := d.GetOk("envs"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			env := tsf.Env{}
			if v, ok := dMap["name"]; ok {
				env.Name = helper.String(v.(string))
			}
			if v, ok := dMap["value"]; ok {
				env.Value = helper.String(v.(string))
			}
			if valueFromMap, ok := helper.InterfaceToMap(dMap, "value_from"); ok {
				valueFrom := tsf.ValueFrom{}
				if fieldRefMap, ok := helper.InterfaceToMap(valueFromMap, "field_ref"); ok {
					fieldRef := tsf.FieldRef{}
					if v, ok := fieldRefMap["field_path"]; ok {
						fieldRef.FieldPath = helper.String(v.(string))
					}
					valueFrom.FieldRef = &fieldRef
				}
				if resourceFieldRefMap, ok := helper.InterfaceToMap(valueFromMap, "resource_field_ref"); ok {
					resourceFieldRef := tsf.ResourceFieldRef{}
					if v, ok := resourceFieldRefMap["resource"]; ok {
						resourceFieldRef.Resource = helper.String(v.(string))
					}
					valueFrom.ResourceFieldRef = &resourceFieldRef
				}
				env.ValueFrom = &valueFrom
			}
			request.Envs = append(request.Envs, &env)
		}
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "service_setting"); ok {
		serviceSetting := tsf.ServiceSetting{}
		if v, ok := dMap["access_type"]; ok {
			serviceSetting.AccessType = helper.IntInt64(v.(int))
		}
		if v, ok := dMap["protocol_ports"]; ok {
			for _, item := range v.([]interface{}) {
				protocolPortsMap := item.(map[string]interface{})
				protocolPort := tsf.ProtocolPort{}
				if v, ok := protocolPortsMap["protocol"]; ok {
					protocolPort.Protocol = helper.String(v.(string))
				}
				if v, ok := protocolPortsMap["port"]; ok {
					protocolPort.Port = helper.IntInt64(v.(int))
				}
				if v, ok := protocolPortsMap["target_port"]; ok {
					protocolPort.TargetPort = helper.IntInt64(v.(int))
				}
				if v, ok := protocolPortsMap["node_port"]; ok {
					protocolPort.NodePort = helper.IntInt64(v.(int))
				}
				serviceSetting.ProtocolPorts = append(serviceSetting.ProtocolPorts, &protocolPort)
			}
		}
		if v, ok := dMap["subnet_id"]; ok {
			serviceSetting.SubnetId = helper.String(v.(string))
		}
		if v, ok := dMap["disable_service"]; ok {
			serviceSetting.DisableService = helper.Bool(v.(bool))
		}
		if v, ok := dMap["headless_service"]; ok {
			serviceSetting.HeadlessService = helper.Bool(v.(bool))
		}
		if v, ok := dMap["allow_delete_service"]; ok {
			serviceSetting.AllowDeleteService = helper.Bool(v.(bool))
		}
		if v, ok := dMap["open_session_affinity"]; ok {
			serviceSetting.OpenSessionAffinity = helper.Bool(v.(bool))
		}
		if v, ok := dMap["session_affinity_timeout_seconds"]; ok {
			serviceSetting.SessionAffinityTimeoutSeconds = helper.IntInt64(v.(int))
		}
		request.ServiceSetting = &serviceSetting
	}

	if v, _ := d.GetOk("deploy_agent"); v != nil {
		request.DeployAgent = helper.Bool(v.(bool))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "scheduling_strategy"); ok {
		schedulingStrategy := tsf.SchedulingStrategy{}
		if v, ok := dMap["type"]; ok {
			schedulingStrategy.Type = helper.String(v.(string))
		}
		request.SchedulingStrategy = &schedulingStrategy
	}

	if v, _ := d.GetOk("incremental_deployment"); v != nil {
		request.IncrementalDeployment = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOk("repo_type"); ok {
		request.RepoType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("volume_info_list"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			volumeInfo := tsf.VolumeInfo{}
			if v, ok := dMap["volume_type"]; ok {
				volumeInfo.VolumeType = helper.String(v.(string))
			}
			if v, ok := dMap["volume_name"]; ok {
				volumeInfo.VolumeName = helper.String(v.(string))
			}
			if v, ok := dMap["volume_config"]; ok {
				volumeInfo.VolumeConfig = helper.String(v.(string))
			}
			request.VolumeInfoList = append(request.VolumeInfoList, &volumeInfo)
		}
	}

	if v, ok := d.GetOk("volume_mount_info_list"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			volumeMountInfo := tsf.VolumeMountInfo{}
			if v, ok := dMap["volume_mount_name"]; ok {
				volumeMountInfo.VolumeMountName = helper.String(v.(string))
			}
			if v, ok := dMap["volume_mount_path"]; ok {
				volumeMountInfo.VolumeMountPath = helper.String(v.(string))
			}
			if v, ok := dMap["volume_mount_sub_path"]; ok {
				volumeMountInfo.VolumeMountSubPath = helper.String(v.(string))
			}
			if v, ok := dMap["read_or_write"]; ok {
				volumeMountInfo.ReadOrWrite = helper.String(v.(string))
			}
			request.VolumeMountInfoList = append(request.VolumeMountInfoList, &volumeMountInfo)
		}
	}

	if v, _ := d.GetOk("volume_clean"); v != nil {
		request.VolumeClean = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOk("agent_profile_list"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			agentProfile := tsf.AgentProfile{}
			if v, ok := dMap["agent_type"]; ok {
				agentProfile.AgentType = helper.String(v.(string))
			}
			if v, ok := dMap["agent_version"]; ok {
				agentProfile.AgentVersion = helper.String(v.(string))
			}
			request.AgentProfileList = append(request.AgentProfileList, &agentProfile)
		}
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "warmup_setting"); ok {
		warmupSetting := tsf.WarmupSetting{}
		if v, ok := dMap["enabled"]; ok {
			warmupSetting.Enabled = helper.Bool(v.(bool))
		}
		if v, ok := dMap["warmup_time"]; ok {
			warmupSetting.WarmupTime = helper.IntUint64(v.(int))
		}
		if v, ok := dMap["curvature"]; ok {
			warmupSetting.Curvature = helper.IntUint64(v.(int))
		}
		if v, ok := dMap["enabled_protection"]; ok {
			warmupSetting.EnabledProtection = helper.Bool(v.(bool))
		}
		request.WarmupSetting = &warmupSetting
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTsfClient().DeployContainerGroup(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s deploy tsf deployContainerGroup failed, reason:%+v", logId, err)
		return err
	}

	if !*response.Response.Result {
		return fmt.Errorf("[CRITAL]%s deploy tsf deployContainerGroup failed", logId)
	}
	d.SetId(groupId)

	service := TsfService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	err = resource.Retry(10*tccommon.WriteRetryTimeout, func() *resource.RetryError {
		groupInfo, err := service.DescribeTsfStartContainerGroupById(ctx, groupId)
		if err != nil {
			return tccommon.RetryError(err)
		}
		if groupInfo == nil {
			err = fmt.Errorf("group %s not exists", groupId)
			return resource.NonRetryableError(err)
		}
		if *groupInfo.Status == "Running" {
			return nil
		}
		if *groupInfo.Status == "Waiting" || *groupInfo.Status == "Updating" {
			return resource.RetryableError(fmt.Errorf("deploy container group status is %s", *groupInfo.Status))
		}
		err = fmt.Errorf("deploy container group status is %v, we won't wait for it finish", *groupInfo.Status)
		return resource.NonRetryableError(err)
	})

	if err != nil {
		log.Printf("[CRITAL]%s deploy container group, reason:%s\n ", logId, err.Error())
		return err
	}

	return resourceTencentCloudTsfDeployContainerGroupRead(d, meta)
}

func resourceTencentCloudTsfDeployContainerGroupRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tsf_deploy_container_group.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := TsfService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	groupId := d.Id()

	deployContainerGroup, err := service.DescribeTsfDeployContainerGroupById(ctx, groupId)
	if err != nil {
		return err
	}

	if deployContainerGroup == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `TsfDeployContainerGroup` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if deployContainerGroup.GroupId != nil {
		_ = d.Set("group_id", deployContainerGroup.GroupId)
	}

	if deployContainerGroup.TagName != nil {
		_ = d.Set("tag_name", deployContainerGroup.TagName)
	}

	if deployContainerGroup.InstanceNum != nil {
		_ = d.Set("instance_num", deployContainerGroup.InstanceNum)
	}

	if deployContainerGroup.Server != nil {
		_ = d.Set("server", deployContainerGroup.Server)
	}

	if deployContainerGroup.Reponame != nil {
		_ = d.Set("reponame", deployContainerGroup.Reponame)
	}

	if deployContainerGroup.CpuLimit != nil {
		_ = d.Set("cpu_limit", deployContainerGroup.CpuLimit)
	}

	if deployContainerGroup.MemLimit != nil {
		_ = d.Set("mem_limit", deployContainerGroup.MemLimit)
	}

	if deployContainerGroup.JvmOpts != nil {
		_ = d.Set("jvm_opts", deployContainerGroup.JvmOpts)
	}

	if deployContainerGroup.CpuRequest != nil {
		_ = d.Set("cpu_request", deployContainerGroup.CpuRequest)
	}

	if deployContainerGroup.MemRequest != nil {
		_ = d.Set("mem_request", deployContainerGroup.MemRequest)
	}

	// if deployContainerGroup.DoNotStart != nil {
	// 	_ = d.Set("do_not_start", deployContainerGroup.DoNotStart)
	// }

	// if deployContainerGroup.RepoName != nil {
	// 	_ = d.Set("repo_name", deployContainerGroup.RepoName)
	// }

	if deployContainerGroup.UpdateType != nil {
		_ = d.Set("update_type", deployContainerGroup.UpdateType)
	}

	if deployContainerGroup.UpdateIvl != nil {
		_ = d.Set("update_ivl", deployContainerGroup.UpdateIvl)
	}

	if deployContainerGroup.AgentCpuRequest != nil {
		_ = d.Set("agent_cpu_request", deployContainerGroup.AgentCpuRequest)
	}

	if deployContainerGroup.AgentCpuLimit != nil {
		_ = d.Set("agent_cpu_limit", deployContainerGroup.AgentCpuLimit)
	}

	if deployContainerGroup.AgentMemRequest != nil {
		_ = d.Set("agent_mem_request", deployContainerGroup.AgentMemRequest)
	}

	if deployContainerGroup.AgentMemLimit != nil {
		_ = d.Set("agent_mem_limit", deployContainerGroup.AgentMemLimit)
	}

	if deployContainerGroup.IstioCpuRequest != nil {
		_ = d.Set("istio_cpu_request", deployContainerGroup.IstioCpuRequest)
	}

	if deployContainerGroup.IstioCpuLimit != nil {
		_ = d.Set("istio_cpu_limit", deployContainerGroup.IstioCpuLimit)
	}

	if deployContainerGroup.IstioMemRequest != nil {
		_ = d.Set("istio_mem_request", deployContainerGroup.IstioMemRequest)
	}

	if deployContainerGroup.IstioMemLimit != nil {
		_ = d.Set("istio_mem_limit", deployContainerGroup.IstioMemLimit)
	}

	// if deployContainerGroup.HealthCheckSettings != nil {
	// 	healthCheckSettingsMap := map[string]interface{}{}

	// 	if deployContainerGroup.HealthCheckSettings.LivenessProbe != nil {
	// 		livenessProbeMap := map[string]interface{}{}

	// 		if deployContainerGroup.HealthCheckSettings.LivenessProbe.ActionType != nil {
	// 			livenessProbeMap["action_type"] = deployContainerGroup.HealthCheckSettings.LivenessProbe.ActionType
	// 		}

	// 		if deployContainerGroup.HealthCheckSettings.LivenessProbe.InitialDelaySeconds != nil {
	// 			livenessProbeMap["initial_delay_seconds"] = deployContainerGroup.HealthCheckSettings.LivenessProbe.InitialDelaySeconds
	// 		}

	// 		if deployContainerGroup.HealthCheckSettings.LivenessProbe.TimeoutSeconds != nil {
	// 			livenessProbeMap["timeout_seconds"] = deployContainerGroup.HealthCheckSettings.LivenessProbe.TimeoutSeconds
	// 		}

	// 		if deployContainerGroup.HealthCheckSettings.LivenessProbe.PeriodSeconds != nil {
	// 			livenessProbeMap["period_seconds"] = deployContainerGroup.HealthCheckSettings.LivenessProbe.PeriodSeconds
	// 		}

	// 		if deployContainerGroup.HealthCheckSettings.LivenessProbe.SuccessThreshold != nil {
	// 			livenessProbeMap["success_threshold"] = deployContainerGroup.HealthCheckSettings.LivenessProbe.SuccessThreshold
	// 		}

	// 		if deployContainerGroup.HealthCheckSettings.LivenessProbe.FailureThreshold != nil {
	// 			livenessProbeMap["failure_threshold"] = deployContainerGroup.HealthCheckSettings.LivenessProbe.FailureThreshold
	// 		}

	// 		if deployContainerGroup.HealthCheckSettings.LivenessProbe.Scheme != nil {
	// 			livenessProbeMap["scheme"] = deployContainerGroup.HealthCheckSettings.LivenessProbe.Scheme
	// 		}

	// 		if deployContainerGroup.HealthCheckSettings.LivenessProbe.Port != nil {
	// 			livenessProbeMap["port"] = deployContainerGroup.HealthCheckSettings.LivenessProbe.Port
	// 		}

	// 		if deployContainerGroup.HealthCheckSettings.LivenessProbe.Path != nil {
	// 			livenessProbeMap["path"] = deployContainerGroup.HealthCheckSettings.LivenessProbe.Path
	// 		}

	// 		if deployContainerGroup.HealthCheckSettings.LivenessProbe.Command != nil {
	// 			livenessProbeMap["command"] = deployContainerGroup.HealthCheckSettings.LivenessProbe.Command
	// 		}

	// 		if deployContainerGroup.HealthCheckSettings.LivenessProbe.Type != nil {
	// 			livenessProbeMap["type"] = deployContainerGroup.HealthCheckSettings.LivenessProbe.Type
	// 		}

	// 		healthCheckSettingsMap["liveness_probe"] = []interface{}{livenessProbeMap}
	// 	}

	// 	if deployContainerGroup.HealthCheckSettings.ReadinessProbe != nil {
	// 		readinessProbeMap := map[string]interface{}{}

	// 		if deployContainerGroup.HealthCheckSettings.ReadinessProbe.ActionType != nil {
	// 			readinessProbeMap["action_type"] = deployContainerGroup.HealthCheckSettings.ReadinessProbe.ActionType
	// 		}

	// 		if deployContainerGroup.HealthCheckSettings.ReadinessProbe.InitialDelaySeconds != nil {
	// 			readinessProbeMap["initial_delay_seconds"] = deployContainerGroup.HealthCheckSettings.ReadinessProbe.InitialDelaySeconds
	// 		}

	// 		if deployContainerGroup.HealthCheckSettings.ReadinessProbe.TimeoutSeconds != nil {
	// 			readinessProbeMap["timeout_seconds"] = deployContainerGroup.HealthCheckSettings.ReadinessProbe.TimeoutSeconds
	// 		}

	// 		if deployContainerGroup.HealthCheckSettings.ReadinessProbe.PeriodSeconds != nil {
	// 			readinessProbeMap["period_seconds"] = deployContainerGroup.HealthCheckSettings.ReadinessProbe.PeriodSeconds
	// 		}

	// 		if deployContainerGroup.HealthCheckSettings.ReadinessProbe.SuccessThreshold != nil {
	// 			readinessProbeMap["success_threshold"] = deployContainerGroup.HealthCheckSettings.ReadinessProbe.SuccessThreshold
	// 		}

	// 		if deployContainerGroup.HealthCheckSettings.ReadinessProbe.FailureThreshold != nil {
	// 			readinessProbeMap["failure_threshold"] = deployContainerGroup.HealthCheckSettings.ReadinessProbe.FailureThreshold
	// 		}

	// 		if deployContainerGroup.HealthCheckSettings.ReadinessProbe.Scheme != nil {
	// 			readinessProbeMap["scheme"] = deployContainerGroup.HealthCheckSettings.ReadinessProbe.Scheme
	// 		}

	// 		if deployContainerGroup.HealthCheckSettings.ReadinessProbe.Port != nil {
	// 			readinessProbeMap["port"] = deployContainerGroup.HealthCheckSettings.ReadinessProbe.Port
	// 		}

	// 		if deployContainerGroup.HealthCheckSettings.ReadinessProbe.Path != nil {
	// 			readinessProbeMap["path"] = deployContainerGroup.HealthCheckSettings.ReadinessProbe.Path
	// 		}

	// 		if deployContainerGroup.HealthCheckSettings.ReadinessProbe.Command != nil {
	// 			readinessProbeMap["command"] = deployContainerGroup.HealthCheckSettings.ReadinessProbe.Command
	// 		}

	// 		if deployContainerGroup.HealthCheckSettings.ReadinessProbe.Type != nil {
	// 			readinessProbeMap["type"] = deployContainerGroup.HealthCheckSettings.ReadinessProbe.Type
	// 		}

	// 		healthCheckSettingsMap["readiness_probe"] = []interface{}{readinessProbeMap}
	// 	}

	// 	_ = d.Set("health_check_settings", []interface{}{healthCheckSettingsMap})
	// }

	if deployContainerGroup.Envs != nil {
		envsList := []interface{}{}
		for _, envs := range deployContainerGroup.Envs {
			envsMap := map[string]interface{}{}

			if envs.Name != nil {
				envsMap["name"] = envs.Name
			}

			if envs.Value != nil {
				envsMap["value"] = envs.Value
			}

			if envs.ValueFrom != nil {
				valueFromMap := map[string]interface{}{}

				if envs.ValueFrom.FieldRef != nil {
					fieldRefMap := map[string]interface{}{}

					if envs.ValueFrom.FieldRef.FieldPath != nil {
						fieldRefMap["field_path"] = envs.ValueFrom.FieldRef.FieldPath
					}

					valueFromMap["field_ref"] = []interface{}{fieldRefMap}
				}

				if envs.ValueFrom.ResourceFieldRef != nil {
					resourceFieldRefMap := map[string]interface{}{}

					if envs.ValueFrom.ResourceFieldRef.Resource != nil {
						resourceFieldRefMap["resource"] = envs.ValueFrom.ResourceFieldRef.Resource
					}

					valueFromMap["resource_field_ref"] = []interface{}{resourceFieldRefMap}
				}

				envsMap["value_from"] = []interface{}{valueFromMap}
			}

			envsList = append(envsList, envsMap)
		}

		_ = d.Set("envs", envsList)

	}

	if deployContainerGroup.DeployAgent != nil {
		_ = d.Set("deploy_agent", deployContainerGroup.DeployAgent)
	}

	if deployContainerGroup.RepoType != nil {
		_ = d.Set("repo_type", deployContainerGroup.RepoType)
	}

	if deployContainerGroup.VolumeInfos != nil {
		volumeInfoListList := []interface{}{}
		for _, volumeInfoList := range deployContainerGroup.VolumeInfos {
			volumeInfoListMap := map[string]interface{}{}

			if volumeInfoList.VolumeType != nil {
				volumeInfoListMap["volume_type"] = volumeInfoList.VolumeType
			}

			if volumeInfoList.VolumeName != nil {
				volumeInfoListMap["volume_name"] = volumeInfoList.VolumeName
			}

			if volumeInfoList.VolumeConfig != nil {
				volumeInfoListMap["volume_config"] = volumeInfoList.VolumeConfig
			}

			volumeInfoListList = append(volumeInfoListList, volumeInfoListMap)
		}

		_ = d.Set("volume_info_list", volumeInfoListList)

	}

	if deployContainerGroup.VolumeMountInfos != nil {
		volumeMountInfoListList := []interface{}{}
		for _, volumeMountInfoList := range deployContainerGroup.VolumeMountInfos {
			volumeMountInfoListMap := map[string]interface{}{}

			if volumeMountInfoList.VolumeMountName != nil {
				volumeMountInfoListMap["volume_mount_name"] = volumeMountInfoList.VolumeMountName
			}

			if volumeMountInfoList.VolumeMountPath != nil {
				volumeMountInfoListMap["volume_mount_path"] = volumeMountInfoList.VolumeMountPath
			}

			if volumeMountInfoList.VolumeMountSubPath != nil {
				volumeMountInfoListMap["volume_mount_sub_path"] = volumeMountInfoList.VolumeMountSubPath
			}

			if volumeMountInfoList.ReadOrWrite != nil {
				volumeMountInfoListMap["read_or_write"] = volumeMountInfoList.ReadOrWrite
			}

			volumeMountInfoListList = append(volumeMountInfoListList, volumeMountInfoListMap)
		}

		_ = d.Set("volume_mount_info_list", volumeMountInfoListList)

	}

	// if deployContainerGroup.VolumeClean != nil {
	// 	_ = d.Set("volume_clean", deployContainerGroup.VolumeClean)
	// }

	// if deployContainerGroup.AgentProfileList != nil {
	// 	agentProfileListList := []interface{}{}
	// 	for _, agentProfileList := range deployContainerGroup.AgentProfileList {
	// 		agentProfileListMap := map[string]interface{}{}

	// 		if deployContainerGroup.AgentProfileList.AgentType != nil {
	// 			agentProfileListMap["agent_type"] = deployContainerGroup.AgentProfileList.AgentType
	// 		}

	// 		if deployContainerGroup.AgentProfileList.AgentVersion != nil {
	// 			agentProfileListMap["agent_version"] = deployContainerGroup.AgentProfileList.AgentVersion
	// 		}

	// 		agentProfileListList = append(agentProfileListList, agentProfileListMap)
	// 	}

	// 	_ = d.Set("agent_profile_list", agentProfileListList)

	// }

	if deployContainerGroup.WarmupSetting != nil {
		warmupSettingMap := map[string]interface{}{}

		if deployContainerGroup.WarmupSetting.Enabled != nil {
			warmupSettingMap["enabled"] = deployContainerGroup.WarmupSetting.Enabled
		}

		if deployContainerGroup.WarmupSetting.WarmupTime != nil {
			warmupSettingMap["warmup_time"] = deployContainerGroup.WarmupSetting.WarmupTime
		}

		if deployContainerGroup.WarmupSetting.Curvature != nil {
			warmupSettingMap["curvature"] = deployContainerGroup.WarmupSetting.Curvature
		}

		if deployContainerGroup.WarmupSetting.EnabledProtection != nil {
			warmupSettingMap["enabled_protection"] = deployContainerGroup.WarmupSetting.EnabledProtection
		}

		_ = d.Set("warmup_setting", []interface{}{warmupSettingMap})
	}

	return nil
}

func resourceTencentCloudTsfDeployContainerGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tsf_deploy_container_group.update")()
	defer tccommon.InconsistentCheck(d, meta)()
	logId := tccommon.GetLogId(tccommon.ContextNil)

	request := tsf.NewModifyContainerReplicasRequest()

	groupId := d.Id()

	request.GroupId = &groupId

	immutableArgs := []string{"group_id", "tag_name", "server", "reponame", "cpu_limit", "mem_limit", "jvm_opts", "cpu_request", "mem_request", "do_not_start", "repo_name", "update_type", "update_ivl", "agent_cpu_request", "agent_cpu_limit", "agent_mem_request", "agent_mem_limit", "istio_cpu_request", "istio_cpu_limit", "istio_mem_request", "istio_mem_limit", "max_surge", "max_unavailable", "health_check_settings", "envs", "service_setting", "deploy_agent", "scheduling_strategy", "incremental_deployment", "repo_type", "volume_infos", "volume_mount_infos", "volume_info_list", "volume_mount_info_list", "volume_clean", "agent_profile_list", "warmup_setting"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("instance_num") {
		if v, ok := d.GetOk("instance_num"); ok {
			request.InstanceNum = helper.IntInt64(v.(int))
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTsfClient().ModifyContainerReplicas(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update tsf unitRule failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudTsfDeployContainerGroupRead(d, meta)
}

func resourceTencentCloudTsfDeployContainerGroupDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tsf_deploy_container_group.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
