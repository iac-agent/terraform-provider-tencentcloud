package tsf

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tsf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tsf/v20180326"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudTsfMicroservice() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTsfMicroserviceRead,
		Schema: map[string]*schema.Schema{
			"namespace_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "namespace id。",
			},

			"status": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "状态 filter，online，offline，single_online。",
			},

			"microservice_id_list": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "microservice id list。",
			},

			"microservice_name_list": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "列表 service names for search。",
			},

			"result": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Microservice paging list information. 注意：此字段可能返回 null，表示无法获取有效值。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"total_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Microservice paging list information. 注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"content": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Microservice list information. 注意：此字段可能返回 null，表示无法获取有效值。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"microservice_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Microservice Id. 注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"microservice_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Microservice 名称 注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"microservice_desc": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Microservice 描述 注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"create_time": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "CreationTime. 注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"update_time": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "last 更新时间.  注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"namespace_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Namespace Id.  注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"run_instance_count": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "run instance count in namespace.  注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"critical_instance_count": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "offline instance count.  注意：此字段可能返回 null，表示无法获取有效值。",
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

func dataSourceTencentCloudTsfMicroserviceRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_tsf_microservice.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("namespace_id"); ok {
		paramMap["NamespaceId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("status"); ok {
		statusSet := v.(*schema.Set).List()
		paramMap["Status"] = helper.InterfacesStringsPoint(statusSet)
	}

	if v, ok := d.GetOk("microservice_id_list"); ok {
		microserviceIdListSet := v.(*schema.Set).List()
		paramMap["MicroserviceIdList"] = helper.InterfacesStringsPoint(microserviceIdListSet)
	}

	if v, ok := d.GetOk("microservice_name_list"); ok {
		microserviceNameListSet := v.(*schema.Set).List()
		paramMap["MicroserviceNameList"] = helper.InterfacesStringsPoint(microserviceNameListSet)
	}

	service := TsfService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var microservice *tsf.TsfPageMicroservice
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeTsfMicroserviceByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		microservice = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(microservice.Content))
	tsfPageMicroserviceMap := map[string]interface{}{}
	if microservice != nil {
		if microservice.TotalCount != nil {
			tsfPageMicroserviceMap["total_count"] = microservice.TotalCount
		}

		if microservice.Content != nil {
			contentList := []interface{}{}
			for _, content := range microservice.Content {
				contentMap := map[string]interface{}{}

				if content.MicroserviceId != nil {
					contentMap["microservice_id"] = content.MicroserviceId
				}

				if content.MicroserviceName != nil {
					contentMap["microservice_name"] = content.MicroserviceName
				}

				if content.MicroserviceDesc != nil {
					contentMap["microservice_desc"] = content.MicroserviceDesc
				}

				if content.CreateTime != nil {
					contentMap["create_time"] = content.CreateTime
				}

				if content.UpdateTime != nil {
					contentMap["update_time"] = content.UpdateTime
				}

				if content.NamespaceId != nil {
					contentMap["namespace_id"] = content.NamespaceId
				}

				if content.RunInstanceCount != nil {
					contentMap["run_instance_count"] = content.RunInstanceCount
				}

				if content.CriticalInstanceCount != nil {
					contentMap["critical_instance_count"] = content.CriticalInstanceCount
				}

				contentList = append(contentList, contentMap)
				ids = append(ids, *content.NamespaceId)
			}

			tsfPageMicroserviceMap["content"] = contentList
		}

		err = d.Set("result", []interface{}{tsfPageMicroserviceMap})
		if err != nil {
			return err
		}
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), tsfPageMicroserviceMap); e != nil {
			return e
		}
	}
	return nil
}
