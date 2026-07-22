package tsf

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tsf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tsf/v20180326"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudTsfPodInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTsfPodInstancesRead,
		Schema: map[string]*schema.Schema{
			"group_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Instance&amp;#39;s 组 ID",
			},

			"pod_name_list": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Filter，pod 名称 list。",
			},

			"result": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "pod instance list。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"total_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Total 数量 records.注意：此字段可能返回 null，表示未找到有效值。",
						},
						"content": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "内容 list.注意：此字段可能返回 null，表示未找到有效值。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"pod_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "实例名称 (corresponding to the pod 名称 in Kubernetes).注意：此字段可能返回 null，表示未找到有效值。",
									},
									"pod_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "实例 ID (corresponding to the pod 实例 ID in Kubernetes).注意：此字段可能返回 null，表示未找到有效值。",
									},
									"status": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "实例状态 Please refer to the definition of instance and container 状态 below. Starting (pod not ready): Starting; Running: Running; Abnormal: Abnormal; Stopped: Stopped;注意：此字段可能返回 null，表示未找到有效值。",
									},
									"reason": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Instance reason for current 状态注意：此字段可能返回 null，表示未找到有效值。",
									},
									"node_ip": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Instance node ip.注意：此字段可能返回 null，表示未找到有效值。",
									},
									"ip": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Instance ip.注意：此字段可能返回 null，表示未找到有效值。",
									},
									"restart_count": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Instance restart count.注意：此字段可能返回 null，表示未找到有效值。",
									},
									"ready_count": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Instance ready count.注意：此字段可能返回 null，表示未找到有效值。",
									},
									"runtime": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Instance run time.注意：此字段可能返回 null，表示未找到有效值。",
									},
									"created_at": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Instance 开始时间.注意：此字段可能返回 null，表示未找到有效值。",
									},
									"service_instance_status": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Instance serve 状态注意：此字段可能返回 null，表示未找到有效值。",
									},
									"instance_available_status": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Instance available 状态注意：此字段可能返回 null，表示未找到有效值。",
									},
									"instance_status": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "实例状态注意：此字段可能返回 null，表示未找到有效值。",
									},
									"node_instance_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Instance node id.注意：此字段可能返回 null，表示未找到有效值。",
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

func dataSourceTencentCloudTsfPodInstancesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_tsf_pod_instances.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("group_id"); ok {
		paramMap["GroupId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("pod_name_list"); ok {
		podNameListSet := v.(*schema.Set).List()
		paramMap["PodNameList"] = helper.InterfacesStringsPoint(podNameListSet)
	}

	service := TsfService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	var groupPodResult *tsf.GroupPodResult
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeTsfDescribePodInstancesByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		groupPodResult = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(groupPodResult.Content))
	groupPodResultMap := map[string]interface{}{}
	if groupPodResult != nil {
		if groupPodResult.TotalCount != nil {
			groupPodResultMap["total_count"] = groupPodResult.TotalCount
		}

		if groupPodResult.Content != nil {
			contentList := []interface{}{}
			for _, content := range groupPodResult.Content {
				contentMap := map[string]interface{}{}

				if content.PodName != nil {
					contentMap["pod_name"] = content.PodName
				}

				if content.PodId != nil {
					contentMap["pod_id"] = content.PodId
				}

				if content.Status != nil {
					contentMap["status"] = content.Status
				}

				if content.Reason != nil {
					contentMap["reason"] = content.Reason
				}

				if content.NodeIp != nil {
					contentMap["node_ip"] = content.NodeIp
				}

				if content.Ip != nil {
					contentMap["ip"] = content.Ip
				}

				if content.RestartCount != nil {
					contentMap["restart_count"] = content.RestartCount
				}

				if content.ReadyCount != nil {
					contentMap["ready_count"] = content.ReadyCount
				}

				if content.Runtime != nil {
					contentMap["runtime"] = content.Runtime
				}

				if content.CreatedAt != nil {
					contentMap["created_at"] = content.CreatedAt
				}

				if content.ServiceInstanceStatus != nil {
					contentMap["service_instance_status"] = content.ServiceInstanceStatus
				}

				if content.InstanceAvailableStatus != nil {
					contentMap["instance_available_status"] = content.InstanceAvailableStatus
				}

				if content.InstanceStatus != nil {
					contentMap["instance_status"] = content.InstanceStatus
				}

				if content.NodeInstanceId != nil {
					contentMap["node_instance_id"] = content.NodeInstanceId
				}

				contentList = append(contentList, contentMap)
				ids = append(ids, *content.PodName)
			}

			groupPodResultMap["content"] = contentList
		}

		_ = d.Set("result", []interface{}{groupPodResultMap})
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), groupPodResultMap); e != nil {
			return e
		}
	}
	return nil
}
