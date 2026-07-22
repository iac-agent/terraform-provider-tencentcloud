package ssl

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ssl "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ssl/v20191205"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudSslDescribeHostDeployRecordDetail() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudSslDescribeHostDeployRecordDetailRead,
		Schema: map[string]*schema.Schema{
			"deploy_record_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Deployment 记录 ID。",
			},

			"deploy_record_detail_list": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Certificate 部署 记录 listNote: 此 字段 可能 返回 NULL，indicating 该 有效 值 不能 是 获取。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Deployment 记录 details ID。",
						},
						"cert_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Deployment 证书 ID",
						},
						"old_cert_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Original binding 证书 IDNote: 此 字段 可能 返回 NULL，indicating 该 有效 值 不能 是 获取。",
						},
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Deployment 实例 ID。",
						},
						"instance_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Deployment 示例 名称",
						},
						"listener_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Deployment 监控 IDNote: 此 字段 可能 返回 NULL，indicating 该 有效 值 不能 是 获取。",
						},
						"domains": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Computed:    true,
							Description: "列表 部署 域名",
						},
						"protocol": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Deployment 监控 protocolNote: 此 字段 可能 返回 NULL，indicating 该 有效 值 不能 是 获取。",
						},
						"status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Deployment state。",
						},
						"error_msg": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Deployment 错误 messageNote: 此 字段 可能 返回 NULL，indicating 该 有效 值 不能 是 获取。",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Deployment 记录 details 创建时间。",
						},
						"update_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Deployment 记录 details last 更新时间。",
						},
						"listener_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Delicate 监控 名称",
						},
						"sni_switch": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "是否turn 在 SNI。",
						},
						"bucket": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "COS 存储 barrel nameNote: 此 字段 可能 返回 NULL，indicating 该 有效 值 不能 是 获取。",
						},
						"namespace": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Named space nameNote: 此 字段 可能 返回 NULL，indicating 该 有效 值 不能 是 获取。",
						},
						"secret_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Secret nameNote: 此 字段 可能 返回 NULL，indicating 该 有效 值 不能 是 获取。",
						},
						"port": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "portNote: 此 字段 可能 返回 NULL，indicating 该 有效 值 不能 是 获取。",
						},
						"env_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "TCB 环境 IDNote: 此 字段 可能 返回 NULL，indicating 该 有效 值 不能 是 获取。",
						},
						"tcb_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Deployed TCB typeNote: 此 字段 可能 返回 NULL，indicating 该 有效 值 不能 是 获取。",
						},
						"region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Deployed TCB regionNote: 此 字段 可能 返回 NULL，indicating 该 有效 值 不能 是 获取。",
						},
					},
				},
			},

			"success_total_count": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Total successNote: 此 字段 可能 返回 NULL，indicating 该 有效 值 不能 是 获取。",
			},

			"failed_total_count": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Total 数量 failuresNote: 此 字段 可能 返回 NULL，indicating 该 有效 值 不能 是 获取。",
			},

			"running_total_count": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Total 数量 deploymentNote: 此 字段 可能 返回 NULL，indicating 该 有效 值 不能 是 获取。",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "用于保存结果。",
			},
		},
	}
}

func dataSourceTencentCloudSslDescribeHostDeployRecordDetailRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_ssl_describe_host_deploy_record_detail.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("deploy_record_id"); ok {
		paramMap["DeployRecordId"] = helper.String(v.(string))
	}

	service := SslService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var deployRecordDetailList []*ssl.DeployRecordDetail
	var successTotalCount, failedTotalCount, runningTotalCount *int64

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, successTotal, failedTotal, runningTotal, e := service.DescribeSslDescribeHostDeployRecordDetailByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		deployRecordDetailList = result
		successTotalCount, failedTotalCount, runningTotalCount = successTotal, failedTotal, runningTotal
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(deployRecordDetailList))
	tmpList := make([]map[string]interface{}, 0, len(deployRecordDetailList))

	if deployRecordDetailList != nil {
		for _, deployRecordDetail := range deployRecordDetailList {
			deployRecordDetailMap := map[string]interface{}{}

			if deployRecordDetail.Id != nil {
				deployRecordDetailMap["id"] = deployRecordDetail.Id
			}

			if deployRecordDetail.CertId != nil {
				deployRecordDetailMap["cert_id"] = deployRecordDetail.CertId
			}

			if deployRecordDetail.OldCertId != nil {
				deployRecordDetailMap["old_cert_id"] = deployRecordDetail.OldCertId
			}

			if deployRecordDetail.InstanceId != nil {
				deployRecordDetailMap["instance_id"] = deployRecordDetail.InstanceId
			}

			if deployRecordDetail.InstanceName != nil {
				deployRecordDetailMap["instance_name"] = deployRecordDetail.InstanceName
			}

			if deployRecordDetail.ListenerId != nil {
				deployRecordDetailMap["listener_id"] = deployRecordDetail.ListenerId
			}

			if deployRecordDetail.Domains != nil {
				deployRecordDetailMap["domains"] = deployRecordDetail.Domains
			}

			if deployRecordDetail.Protocol != nil {
				deployRecordDetailMap["protocol"] = deployRecordDetail.Protocol
			}

			if deployRecordDetail.Status != nil {
				deployRecordDetailMap["status"] = deployRecordDetail.Status
			}

			if deployRecordDetail.ErrorMsg != nil {
				deployRecordDetailMap["error_msg"] = deployRecordDetail.ErrorMsg
			}

			if deployRecordDetail.CreateTime != nil {
				deployRecordDetailMap["create_time"] = deployRecordDetail.CreateTime
			}

			if deployRecordDetail.UpdateTime != nil {
				deployRecordDetailMap["update_time"] = deployRecordDetail.UpdateTime
			}

			if deployRecordDetail.ListenerName != nil {
				deployRecordDetailMap["listener_name"] = deployRecordDetail.ListenerName
			}

			if deployRecordDetail.SniSwitch != nil {
				deployRecordDetailMap["sni_switch"] = deployRecordDetail.SniSwitch
			}

			if deployRecordDetail.Bucket != nil {
				deployRecordDetailMap["bucket"] = deployRecordDetail.Bucket
			}

			if deployRecordDetail.Namespace != nil {
				deployRecordDetailMap["namespace"] = deployRecordDetail.Namespace
			}

			if deployRecordDetail.SecretName != nil {
				deployRecordDetailMap["secret_name"] = deployRecordDetail.SecretName
			}

			if deployRecordDetail.Port != nil {
				deployRecordDetailMap["port"] = deployRecordDetail.Port
			}

			if deployRecordDetail.EnvId != nil {
				deployRecordDetailMap["env_id"] = deployRecordDetail.EnvId
			}

			if deployRecordDetail.TCBType != nil {
				deployRecordDetailMap["tcb_type"] = deployRecordDetail.TCBType
			}

			if deployRecordDetail.Region != nil {
				deployRecordDetailMap["region"] = deployRecordDetail.Region
			}

			ids = append(ids, *deployRecordDetail.InstanceId)
			tmpList = append(tmpList, deployRecordDetailMap)
		}

		_ = d.Set("deploy_record_detail_list", tmpList)
	}

	if successTotalCount != nil {
		_ = d.Set("success_total_count", successTotalCount)
	}

	if failedTotalCount != nil {
		_ = d.Set("failed_total_count", failedTotalCount)
	}

	if runningTotalCount != nil {
		_ = d.Set("running_total_count", runningTotalCount)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
