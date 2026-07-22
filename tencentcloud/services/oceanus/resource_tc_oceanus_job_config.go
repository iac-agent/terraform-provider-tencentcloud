package oceanus

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	oceanus "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/oceanus/v20190422"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudOceanusJobConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudOceanusJobConfigCreate,
		Read:   resourceTencentCloudOceanusJobConfigRead,
		Update: resourceTencentCloudOceanusJobConfigUpdate,
		Delete: resourceTencentCloudOceanusJobConfigDelete,

		Schema: map[string]*schema.Schema{
			"job_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "作业 ID",
			},
			"entrypoint_class": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Main class。",
			},
			"program_args": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Main class 参数。",
			},
			"remark": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "备注",
			},
			"resource_refs": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "Resource reference 数组。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"resource_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "资源 ID",
						},
						"version": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Resource 版本 ID，-1 表示latest 版本",
						},
						"type": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Reference 资源类型，对于 示例，setting main 资源 到 1 表示 jar 包 其中 main class 是 located。",
						},
					},
				},
			},
			"default_parallelism": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Job 默认值 parallelism。",
			},
			"properties": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "System 参数。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "System 配置 键",
						},
						"value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "System 配置 值",
						},
					},
				},
			},
			"cos_bucket": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeString,
				Description: "COS 存储 存储桶名称 使用 通过 作业。",
			},
			"log_collect": {
				Optional:    true,
				Type:        schema.TypeBool,
				Description: "是否collect 作业 logs。",
			},
			"job_manager_spec": {
				Optional:    true,
				Type:        schema.TypeFloat,
				Description: "JobManager 规格。",
			},
			"task_manager_spec": {
				Optional:    true,
				Type:        schema.TypeFloat,
				Description: "TaskManager 规格。",
			},
			"cls_logset_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "CLS logset ID。",
			},
			"cls_topic_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "CLS 日志主题 ID",
			},
			"log_collect_type": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Log collection 类型 2:CLS; 3:COS。",
			},
			"python_version": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Python 版本 使用 通过 pyflink 作业 在 runtime。",
			},
			"work_space_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Workspace SerialId。",
			},
			"log_level": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Log 级别",
			},
			"auto_recover": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Oceanus 平台 作业 recovery switch 1: 在 -1: 关闭。",
			},
			"clazz_levels": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "Class 日志 级别",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"clazz": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Java class full path注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"level": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Log 级别 TRACE，DEBUG，INFO，WARN，ERROR注意：此字段可能返回 null，表示无法获取有效值。",
						},
					},
				},
			},
			"expert_mode_on": {
				Optional:    true,
				Type:        schema.TypeBool,
				Description: "是否enable expert 模式",
			},
			"expert_mode_configuration": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Expert 模式 配置。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"job_graph": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Job graph注意：此字段可能返回 null，表示无法获取有效值。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"nodes": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Point 集合 的 running graph注意：此字段可能返回 null，表示无法获取有效值。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"id": {
													Type:        schema.TypeInt,
													Required:    true,
													Description: "Node ID注意：此字段可能返回 null，表示无法获取有效值。",
												},
												"description": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Node description注意：此字段可能返回 null，表示无法获取有效值。",
												},
												"name": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Node name注意：此字段可能返回 null，表示无法获取有效值。",
												},
												"parallelism": {
													Type:        schema.TypeInt,
													Required:    true,
													Description: "Node parallelism注意：此字段可能返回 null，表示无法获取有效值。",
												},
											},
										},
									},
									"edges": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Edge 集合 的 running graph注意：此字段可能返回 null，表示无法获取有效值。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"source": {
													Type:        schema.TypeInt,
													Required:    true,
													Description: "Starting 节点 ID edge注意：此字段可能返回 null，表示无法获取有效值。",
												},
												"target": {
													Type:        schema.TypeInt,
													Required:    true,
													Description: "Target 节点 ID edge注意：此字段可能返回 null，表示无法获取有效值。",
												},
											},
										},
									},
								},
							},
						},
						"node_config": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Node configuration注意：此字段可能返回 null，表示无法获取有效值。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Node ID注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"parallelism": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Node parallelism注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"slot_sharing_group": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Slot sharing group注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"configuration": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Configuration properties注意：此字段可能返回 null，表示无法获取有效值。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"key": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "System 配置 键",
												},
												"value": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "System 配置 值",
												},
											},
										},
									},
									"state_ttl": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "State TTL 配置 的 节点，separated 通过 semicolons注意：此字段可能返回 null，表示无法获取有效值。",
									},
								},
							},
						},
						"slot_sharing_groups": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Slot sharing groups注意：此字段可能返回 null，表示无法获取有效值。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "名称 SlotSharingGroup注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"spec": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Required:    true,
										Description: "Specification 的 SlotSharingGroup注意：此字段可能返回 null，表示无法获取有效值。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"cpu": {
													Type:        schema.TypeFloat,
													Required:    true,
													Description: "Applicable CPU注意：此字段可能返回 null，表示无法获取有效值。",
												},
												"heap_memory": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "默认为 b，supporting units 是 b，kb，mb，gb注意：此字段可能返回 null，表示无法获取有效值。",
												},
												"off_heap_memory": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "默认为 b，supporting units 是 b，kb，mb，gb注意：此字段可能返回 null，表示无法获取有效值。",
												},
												"managed_memory": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "默认为 b，supporting units 是 b，kb，mb，gb注意：此字段可能返回 null，表示无法获取有效值。",
												},
											},
										},
									},
									"description": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "描述 SlotSharingGroup注意：此字段可能返回 null，表示无法获取有效值。",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudOceanusJobConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_oceanus_job_config.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId    = tccommon.GetLogId(tccommon.ContextNil)
		request  = oceanus.NewCreateJobConfigRequest()
		response = oceanus.NewCreateJobConfigResponse()
		jobId    string
		version  string
	)
	if v, ok := d.GetOk("job_id"); ok {
		request.JobId = helper.String(v.(string))
		jobId = v.(string)
	}

	if v, ok := d.GetOk("entrypoint_class"); ok {
		request.EntrypointClass = helper.String(v.(string))
	}

	if v, ok := d.GetOk("program_args"); ok {
		request.ProgramArgs = helper.String(v.(string))
	}

	if v, ok := d.GetOk("remark"); ok {
		request.Remark = helper.String(v.(string))
	}

	if v, ok := d.GetOk("resource_refs"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			resourceRef := oceanus.ResourceRef{}
			if v, ok := dMap["resource_id"]; ok {
				resourceRef.ResourceId = helper.String(v.(string))
			}

			if v, ok := dMap["version"]; ok {
				resourceRef.Version = helper.IntInt64(v.(int))
			}

			if v, ok := dMap["type"]; ok {
				resourceRef.Type = helper.IntInt64(v.(int))
			}

			request.ResourceRefs = append(request.ResourceRefs, &resourceRef)
		}
	}

	if v, ok := d.GetOkExists("default_parallelism"); ok {
		request.DefaultParallelism = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("properties"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			property := oceanus.Property{}
			if v, ok := dMap["key"]; ok {
				property.Key = helper.String(v.(string))
			}

			if v, ok := dMap["value"]; ok {
				property.Value = helper.String(v.(string))
			}

			request.Properties = append(request.Properties, &property)
		}
	}

	if v, ok := d.GetOk("cos_bucket"); ok {
		request.COSBucket = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("log_collect"); ok {
		request.LogCollect = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOkExists("job_manager_spec"); ok {
		request.JobManagerSpec = helper.Float64(v.(float64))
	}

	if v, ok := d.GetOkExists("task_manager_spec"); ok {
		request.TaskManagerSpec = helper.Float64(v.(float64))
	}

	if v, ok := d.GetOk("cls_logset_id"); ok {
		request.ClsLogsetId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("cls_topic_id"); ok {
		request.ClsTopicId = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("log_collect_type"); ok {
		request.LogCollectType = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("python_version"); ok {
		request.PythonVersion = helper.String(v.(string))
	}

	if v, ok := d.GetOk("work_space_id"); ok {
		request.WorkSpaceId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("log_level"); ok {
		request.LogLevel = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("auto_recover"); ok {
		request.AutoRecover = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("clazz_levels"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			clazzLevel := oceanus.ClazzLevel{}
			if v, ok := dMap["clazz"]; ok {
				clazzLevel.Clazz = helper.String(v.(string))
			}

			if v, ok := dMap["level"]; ok {
				clazzLevel.Level = helper.String(v.(string))
			}

			request.ClazzLevels = append(request.ClazzLevels, &clazzLevel)
		}
	}

	if v, ok := d.GetOkExists("expert_mode_on"); ok {
		request.ExpertModeOn = helper.Bool(v.(bool))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "expert_mode_configuration"); ok {
		expertModeConfiguration := oceanus.ExpertModeConfiguration{}
		if jobGraphMap, ok := helper.InterfaceToMap(dMap, "job_graph"); ok {
			jobGraph := oceanus.JobGraph{}
			if v, ok := jobGraphMap["nodes"]; ok {
				for _, item := range v.([]interface{}) {
					nodesMap := item.(map[string]interface{})
					jobGraphNode := oceanus.JobGraphNode{}
					if v, ok := nodesMap["id"]; ok {
						jobGraphNode.Id = helper.IntInt64(v.(int))
					}

					if v, ok := nodesMap["description"]; ok {
						jobGraphNode.Description = helper.String(v.(string))
					}

					if v, ok := nodesMap["name"]; ok {
						jobGraphNode.Name = helper.String(v.(string))
					}

					if v, ok := nodesMap["parallelism"]; ok {
						jobGraphNode.Parallelism = helper.IntInt64(v.(int))
					}

					jobGraph.Nodes = append(jobGraph.Nodes, &jobGraphNode)
				}
			}

			if v, ok := jobGraphMap["edges"]; ok {
				for _, item := range v.([]interface{}) {
					edgesMap := item.(map[string]interface{})
					jobGraphEdge := oceanus.JobGraphEdge{}
					if v, ok := edgesMap["source"]; ok {
						jobGraphEdge.Source = helper.IntInt64(v.(int))
					}

					if v, ok := edgesMap["target"]; ok {
						jobGraphEdge.Target = helper.IntInt64(v.(int))
					}

					jobGraph.Edges = append(jobGraph.Edges, &jobGraphEdge)
				}
			}

			expertModeConfiguration.JobGraph = &jobGraph
		}

		if v, ok := dMap["node_config"]; ok {
			for _, item := range v.([]interface{}) {
				nodeConfigMap := item.(map[string]interface{})
				nodeConfig := oceanus.NodeConfig{}
				if v, ok := nodeConfigMap["id"]; ok {
					nodeConfig.Id = helper.IntInt64(v.(int))
				}

				if v, ok := nodeConfigMap["parallelism"]; ok {
					nodeConfig.Parallelism = helper.IntInt64(v.(int))
				}

				if v, ok := nodeConfigMap["slot_sharing_group"]; ok {
					nodeConfig.SlotSharingGroup = helper.String(v.(string))
				}

				if v, ok := nodeConfigMap["configuration"]; ok {
					for _, item := range v.([]interface{}) {
						configurationMap := item.(map[string]interface{})
						property := oceanus.Property{}
						if v, ok := configurationMap["key"]; ok {
							property.Key = helper.String(v.(string))
						}

						if v, ok := configurationMap["value"]; ok {
							property.Value = helper.String(v.(string))
						}

						nodeConfig.Configuration = append(nodeConfig.Configuration, &property)
					}
				}

				if v, ok := nodeConfigMap["state_ttl"]; ok {
					nodeConfig.StateTTL = helper.String(v.(string))
				}

				expertModeConfiguration.NodeConfig = append(expertModeConfiguration.NodeConfig, &nodeConfig)
			}
		}

		if v, ok := dMap["slot_sharing_groups"]; ok {
			for _, item := range v.([]interface{}) {
				slotSharingGroupsMap := item.(map[string]interface{})
				slotSharingGroup := oceanus.SlotSharingGroup{}
				if v, ok := slotSharingGroupsMap["name"]; ok {
					slotSharingGroup.Name = helper.String(v.(string))
				}

				if specMap, ok := helper.InterfaceToMap(slotSharingGroupsMap, "spec"); ok {
					slotSharingGroupSpec := oceanus.SlotSharingGroupSpec{}
					if v, ok := specMap["cpu"]; ok {
						slotSharingGroupSpec.CPU = helper.Float64(v.(float64))
					}

					if v, ok := specMap["heap_memory"]; ok {
						slotSharingGroupSpec.HeapMemory = helper.String(v.(string))
					}

					if v, ok := specMap["off_heap_memory"]; ok {
						slotSharingGroupSpec.OffHeapMemory = helper.String(v.(string))
					}

					if v, ok := specMap["managed_memory"]; ok {
						slotSharingGroupSpec.ManagedMemory = helper.String(v.(string))
					}

					slotSharingGroup.Spec = &slotSharingGroupSpec
				}

				if v, ok := slotSharingGroupsMap["description"]; ok {
					slotSharingGroup.Description = helper.String(v.(string))
				}

				expertModeConfiguration.SlotSharingGroups = append(expertModeConfiguration.SlotSharingGroups, &slotSharingGroup)
			}
		}

		request.ExpertModeConfiguration = &expertModeConfiguration
	}

	request.AutoDelete = helper.IntInt64(0)

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseOceanusClient().CreateJobConfig(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil {
			e = fmt.Errorf("oceanus JobConfig not exists")
			return resource.NonRetryableError(e)
		}

		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create oceanus JobConfig failed, reason:%+v", logId, err)
		return err
	}

	versionInt := *response.Response.Version
	version = strconv.FormatUint(versionInt, 10)
	d.SetId(strings.Join([]string{jobId, version}, tccommon.FILED_SP))

	return resourceTencentCloudOceanusJobConfigRead(d, meta)
}

func resourceTencentCloudOceanusJobConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_oceanus_job_config.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = OceanusService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	jobId := idSplit[0]
	version := idSplit[1]

	JobConfig, err := service.DescribeOceanusJobConfigById(ctx, jobId, version)
	if err != nil {
		return err
	}

	if JobConfig == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `OceanusJobConfig` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if JobConfig.JobId != nil {
		_ = d.Set("job_id", JobConfig.JobId)
	}

	if JobConfig.EntrypointClass != nil {
		_ = d.Set("entrypoint_class", JobConfig.EntrypointClass)
	}

	if JobConfig.ProgramArgs != nil {
		_ = d.Set("program_args", JobConfig.ProgramArgs)
	}

	if JobConfig.Remark != nil {
		_ = d.Set("remark", JobConfig.Remark)
	}

	if JobConfig.ResourceRefDetails != nil {
		resourceRefsList := []interface{}{}
		for _, resourceRefs := range JobConfig.ResourceRefDetails {
			resourceRefsMap := map[string]interface{}{}

			if resourceRefs.ResourceId != nil {
				resourceRefsMap["resource_id"] = resourceRefs.ResourceId
			}

			if resourceRefs.Version != nil {
				resourceRefsMap["version"] = resourceRefs.Version
			}

			if resourceRefs.Type != nil {
				resourceRefsMap["type"] = resourceRefs.Type
			}

			resourceRefsList = append(resourceRefsList, resourceRefsMap)
		}

		_ = d.Set("resource_refs", resourceRefsList)
	}

	if JobConfig.DefaultParallelism != nil {
		_ = d.Set("default_parallelism", JobConfig.DefaultParallelism)
	}

	if JobConfig.Properties != nil {
		propertiesList := []interface{}{}
		for _, properties := range JobConfig.Properties {
			propertiesMap := map[string]interface{}{}

			if properties.Key != nil {
				propertiesMap["key"] = properties.Key
			}

			if properties.Value != nil {
				propertiesMap["value"] = properties.Value
			}

			propertiesList = append(propertiesList, propertiesMap)
		}

		_ = d.Set("properties", propertiesList)
	}

	if JobConfig.COSBucket != nil {
		_ = d.Set("cos_bucket", JobConfig.COSBucket)
	}

	if JobConfig.LogCollect != nil {
		if *JobConfig.LogCollect == 1 || *JobConfig.LogCollect == 3 || *JobConfig.LogCollect == 4 {
			_ = d.Set("log_collect", true)
		} else {
			_ = d.Set("log_collect", false)
		}
	}

	if JobConfig.JobManagerSpec != nil {
		_ = d.Set("job_manager_spec", JobConfig.JobManagerSpec)
	}

	if JobConfig.TaskManagerSpec != nil {
		_ = d.Set("task_manager_spec", JobConfig.TaskManagerSpec)
	}

	if JobConfig.ClsLogsetId != nil {
		_ = d.Set("cls_logset_id", JobConfig.ClsLogsetId)
	}

	if JobConfig.ClsTopicId != nil {
		_ = d.Set("cls_topic_id", JobConfig.ClsTopicId)
	}

	if JobConfig.PythonVersion != nil {
		_ = d.Set("python_version", JobConfig.PythonVersion)
	}

	if JobConfig.LogLevel != nil {
		_ = d.Set("log_level", JobConfig.LogLevel)
	}

	if JobConfig.AutoRecover != nil {
		_ = d.Set("auto_recover", JobConfig.AutoRecover)
	}

	if JobConfig.ClazzLevels != nil {
		clazzLevelsList := []interface{}{}
		for _, clazzLevels := range JobConfig.ClazzLevels {
			clazzLevelsMap := map[string]interface{}{}

			if clazzLevels.Clazz != nil {
				clazzLevelsMap["clazz"] = clazzLevels.Clazz
			}

			if clazzLevels.Level != nil {
				clazzLevelsMap["level"] = clazzLevels.Level
			}

			clazzLevelsList = append(clazzLevelsList, clazzLevelsMap)
		}

		_ = d.Set("clazz_levels", clazzLevelsList)
	}

	if JobConfig.ExpertModeOn != nil {
		_ = d.Set("expert_mode_on", JobConfig.ExpertModeOn)
	}

	if JobConfig.ExpertModeConfiguration != nil {
		expertModeConfigurationMap := map[string]interface{}{}

		if JobConfig.ExpertModeConfiguration.JobGraph != nil {
			jobGraphMap := map[string]interface{}{}

			if JobConfig.ExpertModeConfiguration.JobGraph.Nodes != nil {
				nodesList := []interface{}{}
				for _, nodes := range JobConfig.ExpertModeConfiguration.JobGraph.Nodes {
					nodesMap := map[string]interface{}{}

					if nodes.Id != nil {
						nodesMap["id"] = nodes.Id
					}

					if nodes.Description != nil {
						nodesMap["description"] = nodes.Description
					}

					if nodes.Name != nil {
						nodesMap["name"] = nodes.Name
					}

					if nodes.Parallelism != nil {
						nodesMap["parallelism"] = nodes.Parallelism
					}

					nodesList = append(nodesList, nodesMap)
				}

				jobGraphMap["nodes"] = nodesList
			}

			if JobConfig.ExpertModeConfiguration.JobGraph.Edges != nil {
				edgesList := []interface{}{}
				for _, edges := range JobConfig.ExpertModeConfiguration.JobGraph.Edges {
					edgesMap := map[string]interface{}{}

					if edges.Source != nil {
						edgesMap["source"] = edges.Source
					}

					if edges.Target != nil {
						edgesMap["target"] = edges.Target
					}

					edgesList = append(edgesList, edgesMap)
				}

				jobGraphMap["edges"] = edgesList
			}

			expertModeConfigurationMap["job_graph"] = []interface{}{jobGraphMap}
		}

		if JobConfig.ExpertModeConfiguration.NodeConfig != nil {
			nodeConfigList := []interface{}{}
			for _, nodeConfig := range JobConfig.ExpertModeConfiguration.NodeConfig {
				nodeConfigMap := map[string]interface{}{}

				if nodeConfig.Id != nil {
					nodeConfigMap["id"] = nodeConfig.Id
				}

				if nodeConfig.Parallelism != nil {
					nodeConfigMap["parallelism"] = nodeConfig.Parallelism
				}

				if nodeConfig.SlotSharingGroup != nil {
					nodeConfigMap["slot_sharing_group"] = nodeConfig.SlotSharingGroup
				}

				if nodeConfig.Configuration != nil {
					configurationList := []interface{}{}
					for _, configuration := range nodeConfig.Configuration {
						configurationMap := map[string]interface{}{}

						if configuration.Key != nil {
							configurationMap["key"] = configuration.Key
						}

						if configuration.Value != nil {
							configurationMap["value"] = configuration.Value
						}

						configurationList = append(configurationList, configurationMap)
					}

					nodeConfigMap["configuration"] = configurationList
				}

				if nodeConfig.StateTTL != nil {
					nodeConfigMap["state_ttl"] = nodeConfig.StateTTL
				}

				nodeConfigList = append(nodeConfigList, nodeConfigMap)
			}

			expertModeConfigurationMap["node_config"] = nodeConfigList
		}

		if JobConfig.ExpertModeConfiguration.SlotSharingGroups != nil {
			slotSharingGroupsList := []interface{}{}
			for _, slotSharingGroups := range JobConfig.ExpertModeConfiguration.SlotSharingGroups {
				slotSharingGroupsMap := map[string]interface{}{}

				if slotSharingGroups.Name != nil {
					slotSharingGroupsMap["name"] = slotSharingGroups.Name
				}

				if slotSharingGroups.Spec != nil {
					specMap := map[string]interface{}{}

					if slotSharingGroups.Spec.CPU != nil {
						specMap["cpu"] = slotSharingGroups.Spec.CPU
					}

					if slotSharingGroups.Spec.HeapMemory != nil {
						specMap["heap_memory"] = slotSharingGroups.Spec.HeapMemory
					}

					if slotSharingGroups.Spec.OffHeapMemory != nil {
						specMap["off_heap_memory"] = slotSharingGroups.Spec.OffHeapMemory
					}

					if slotSharingGroups.Spec.ManagedMemory != nil {
						specMap["managed_memory"] = slotSharingGroups.Spec.ManagedMemory
					}

					slotSharingGroupsMap["spec"] = []interface{}{specMap}
				}

				if slotSharingGroups.Description != nil {
					slotSharingGroupsMap["description"] = slotSharingGroups.Description
				}

				slotSharingGroupsList = append(slotSharingGroupsList, slotSharingGroupsMap)
			}

			expertModeConfigurationMap["slot_sharing_groups"] = slotSharingGroupsList
		}

		_ = d.Set("expert_mode_configuration", []interface{}{expertModeConfigurationMap})
	}

	return nil
}

func resourceTencentCloudOceanusJobConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_oceanus_job_config.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	immutableArgs := []string{"job_id", "entrypoint_class", "program_args", "remark", "resource_refs", "default_parallelism", "properties", "c_o_s_bucket", "log_collect", "job_manager_spec", "task_manager_spec", "cls_logset_id", "cls_topic_id", "log_collect_type", "python_version", "work_space_id", "log_level", "auto_recover", "clazz_levels", "expert_mode_on", "expert_mode_configuration"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	return resourceTencentCloudOceanusJobConfigRead(d, meta)
}

func resourceTencentCloudOceanusJobConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_oceanus_job_config.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = OceanusService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	jobId := idSplit[0]
	version := idSplit[1]

	if err := service.DeleteOceanusJobConfigById(ctx, jobId, version); err != nil {
		return err
	}

	return nil
}
