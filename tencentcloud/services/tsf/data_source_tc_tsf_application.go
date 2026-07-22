package tsf

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tsf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tsf/v20180326"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudTsfApplication() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTsfApplicationRead,
		Schema: map[string]*schema.Schema{
			"application_type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "应用 类型 V OR C，V 表示 VM，C 表示 容器。",
			},

			"microservice_type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "microservice 类型 应用。",
			},

			"application_resource_type_list": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "An 数组 应用 资源 types。",
			},

			"application_id_list": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "ID 列表。",
			},

			"result": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "应用 paging 列表 信息。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"total_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "总数 数量 applications。",
						},
						"content": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "列表 应用 信息。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"application_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "ID 应用。",
									},
									"application_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "名称 应用。",
									},
									"application_desc": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "描述 应用。",
									},
									"application_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "类型 应用。",
									},
									"microservice_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "microservice 类型 应用。",
									},
									"prog_lang": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Programming 语言",
									},
									"create_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "创建时间。",
									},
									"update_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "更新时间。",
									},
									"application_resource_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "应用 资源类型",
									},
									"application_runtime_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "应用 runtime 类型",
									},
									"apigateway_service_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "网关 服务 ID。",
									},
									"application_remark_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "备注 名称",
									},
									"service_config_list": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "服务 配置 列表。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "serviceName。",
												},
												"ports": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "端口 列表。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"target_port": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "服务 端口",
															},
															"protocol": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "协议",
															},
														},
													},
												},
												"health_check": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "health check setting。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"path": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "health check 路径",
															},
														},
													},
												},
											},
										},
									},
									"ignore_create_image_repository": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "whether ignore create 镜像 repository。",
									},
								},
							},
						},
					},
				},
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "用于保存结果。",
			},
		},
	}
}

func dataSourceTencentCloudTsfApplicationRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_tsf_application.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("application_type"); ok {
		paramMap["ApplicationType"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("microservice_type"); ok {
		paramMap["MicroserviceType"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("application_resource_type_list"); ok {
		applicationResourceTypeListSet := v.(*schema.Set).List()
		paramMap["ApplicationResourceTypeList"] = helper.InterfacesStringsPoint(applicationResourceTypeListSet)
	}

	if v, ok := d.GetOk("application_id_list"); ok {
		applicationIdListSet := v.(*schema.Set).List()
		paramMap["ApplicationIdList"] = helper.InterfacesStringsPoint(applicationIdListSet)
	}

	service := TsfService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var application *tsf.TsfPageApplication
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeTsfApplicationByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		application = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(application.Content))
	tsfPageApplicationMap := map[string]interface{}{}
	if application != nil {
		if application.TotalCount != nil {
			tsfPageApplicationMap["total_count"] = application.TotalCount
		}

		if application.Content != nil {
			contentList := []interface{}{}
			for _, content := range application.Content {
				contentMap := map[string]interface{}{}

				if content.ApplicationId != nil {
					contentMap["application_id"] = content.ApplicationId
				}

				if content.ApplicationName != nil {
					contentMap["application_name"] = content.ApplicationName
				}

				if content.ApplicationDesc != nil {
					contentMap["application_desc"] = content.ApplicationDesc
				}

				if content.ApplicationType != nil {
					contentMap["application_type"] = content.ApplicationType
				}

				if content.MicroserviceType != nil {
					contentMap["microservice_type"] = content.MicroserviceType
				}

				if content.ProgLang != nil {
					contentMap["prog_lang"] = content.ProgLang
				}

				if content.CreateTime != nil {
					contentMap["create_time"] = content.CreateTime
				}

				if content.UpdateTime != nil {
					contentMap["update_time"] = content.UpdateTime
				}

				if content.ApplicationResourceType != nil {
					contentMap["application_resource_type"] = content.ApplicationResourceType
				}

				if content.ApplicationRuntimeType != nil {
					contentMap["application_runtime_type"] = content.ApplicationRuntimeType
				}

				if content.ApigatewayServiceId != nil {
					contentMap["apigateway_service_id"] = content.ApigatewayServiceId
				}

				if content.ApplicationRemarkName != nil {
					contentMap["application_remark_name"] = content.ApplicationRemarkName
				}

				if content.ServiceConfigList != nil {
					serviceConfigListList := []interface{}{}
					for _, serviceConfigList := range content.ServiceConfigList {
						serviceConfigListMap := map[string]interface{}{}

						if serviceConfigList.Name != nil {
							serviceConfigListMap["name"] = serviceConfigList.Name
						}

						if serviceConfigList.Ports != nil {
							portsList := []interface{}{}
							for _, ports := range serviceConfigList.Ports {
								portsMap := map[string]interface{}{}

								if ports.TargetPort != nil {
									portsMap["target_port"] = ports.TargetPort
								}

								if ports.Protocol != nil {
									portsMap["protocol"] = ports.Protocol
								}

								portsList = append(portsList, portsMap)
							}

							serviceConfigListMap["ports"] = portsList
						}

						if serviceConfigList.HealthCheck != nil {
							healthCheckMap := map[string]interface{}{}

							if serviceConfigList.HealthCheck.Path != nil {
								healthCheckMap["path"] = serviceConfigList.HealthCheck.Path
							}

							serviceConfigListMap["health_check"] = []interface{}{healthCheckMap}
						}

						serviceConfigListList = append(serviceConfigListList, serviceConfigListMap)
					}

					contentMap["service_config_list"] = serviceConfigListList
				}

				if content.IgnoreCreateImageRepository != nil {
					contentMap["ignore_create_image_repository"] = content.IgnoreCreateImageRepository
				}

				contentList = append(contentList, contentMap)
				ids = append(ids, *content.ApplicationId)
			}

			tsfPageApplicationMap["content"] = contentList
		}

		_ = d.Set("result", []interface{}{tsfPageApplicationMap})
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), tsfPageApplicationMap); e != nil {
			return e
		}
	}
	return nil
}
