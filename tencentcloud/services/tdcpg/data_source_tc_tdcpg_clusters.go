package tdcpg

import (
	"context"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tdcpg "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tdcpg/v20211118"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudTdcpgClusters() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTdcpgClustersRead,
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "集群 ID",
			},

			"cluster_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "集群名称",
			},

			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "cluster 状态",
			},

			"pay_mode": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "付费模式",
			},

			"project_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "项目 ID，默认为 0，means default project。",
			},

			"list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "cluster list。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cluster_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "集群 ID",
						},
						"cluster_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "集群名称",
						},
						"region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "地域",
						},
						"zone": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "可用区",
						},
						"db_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "db 版本",
						},
						"project_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "项目 ID",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "状态",
						},
						"status_desc": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "状态 描述",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "创建时间。",
						},
						"storage_used": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "storage used，unit is GB。",
						},
						"storage_limit": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "storage 限制，unit is GB。",
						},
						"pay_mode": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "付费模式",
						},
						"pay_period_end_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "pay 周期 过期时间。",
						},
						"auto_renew_flag": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "auto 续费标识",
						},
						"db_charset": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "db charset。",
						},
						"instance_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "instance count。",
						},
						"endpoint_set": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "endpoint set。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"endpoint_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "endpoint id。",
									},
									"cluster_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "集群 ID",
									},
									"endpoint_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "endpoint 名称",
									},
									"endpoint_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "endpoint 类型",
									},
									"vpc_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "私有网络 ID",
									},
									"subnet_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "子网 ID",
									},
									"private_ip": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "private ip。",
									},
									"private_port": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "private 端口",
									},
									"wan_ip": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "wan ip。",
									},
									"wan_port": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "wan 端口",
									},
									"wan_domain": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "wan 域名",
									},
								},
							},
						},
						"db_major_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "db major 版本",
						},
						"db_kernel_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "db kernel 版本",
						},
						"storage_pay_mode": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "storage 付费模式，可选 值 is PREPAID or POSTPAID_BY_HOUR。",
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

func dataSourceTencentCloudTdcpgClustersRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_tdcpg_clusters.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("cluster_id"); ok {
		paramMap["cluster_id"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("cluster_name"); ok {
		paramMap["cluster_name"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("status"); ok {
		paramMap["status"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("pay_mode"); ok {
		paramMap["pay_mode"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("project_id"); ok {
		paramMap["project_id"] = helper.IntInt64(v.(int))
	}

	tdcpgService := TdcpgService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var clusterSet []*tdcpg.Cluster
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		results, e := tdcpgService.DescribeTdcpgClustersByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		clusterSet = results
		return nil
	})
	if err != nil {
		log.Printf("[CRITICAL]%s read Tdcpg clusterSet failed, reason:%+v", logId, err)
		return err
	}

	ids := make([]string, 0, len(clusterSet))
	clusterList := make([]map[string]interface{}, 0, len(clusterSet))
	if clusterSet != nil {
		for _, cluster := range clusterSet {
			clusterSetMap := map[string]interface{}{}
			if cluster.ClusterId != nil {
				clusterSetMap["cluster_id"] = cluster.ClusterId
			}
			if cluster.ClusterName != nil {
				clusterSetMap["cluster_name"] = cluster.ClusterName
			}
			if cluster.Region != nil {
				clusterSetMap["region"] = cluster.Region
			}
			if cluster.Zone != nil {
				clusterSetMap["zone"] = cluster.Zone
			}
			if cluster.DBVersion != nil {
				clusterSetMap["db_version"] = cluster.DBVersion
			}
			if cluster.ProjectId != nil {
				clusterSetMap["project_id"] = cluster.ProjectId
			}
			if cluster.Status != nil {
				clusterSetMap["status"] = cluster.Status
			}
			if cluster.StatusDesc != nil {
				clusterSetMap["status_desc"] = cluster.StatusDesc
			}
			if cluster.CreateTime != nil {
				clusterSetMap["create_time"] = cluster.CreateTime
			}
			if cluster.StorageUsed != nil {
				clusterSetMap["storage_used"] = cluster.StorageUsed
			}
			if cluster.StorageLimit != nil {
				clusterSetMap["storage_limit"] = cluster.StorageLimit
			}
			if cluster.PayMode != nil {
				clusterSetMap["pay_mode"] = cluster.PayMode
			}
			if cluster.PayPeriodEndTime != nil {
				clusterSetMap["pay_period_end_time"] = cluster.PayPeriodEndTime
			}
			if cluster.AutoRenewFlag != nil {
				clusterSetMap["auto_renew_flag"] = cluster.AutoRenewFlag
			}
			if cluster.DBCharset != nil {
				clusterSetMap["db_charset"] = cluster.DBCharset
			}
			if cluster.InstanceCount != nil {
				clusterSetMap["instance_count"] = cluster.InstanceCount
			}
			if cluster.EndpointSet != nil {
				endpointSetList := []interface{}{}
				for _, endpointSet := range cluster.EndpointSet {
					endpointSetMap := map[string]interface{}{}
					if endpointSet.EndpointId != nil {
						endpointSetMap["endpoint_id"] = endpointSet.EndpointId
					}
					if endpointSet.ClusterId != nil {
						endpointSetMap["cluster_id"] = endpointSet.ClusterId
					}
					if endpointSet.EndpointName != nil {
						endpointSetMap["endpoint_name"] = endpointSet.EndpointName
					}
					if endpointSet.EndpointType != nil {
						endpointSetMap["endpoint_type"] = endpointSet.EndpointType
					}
					if endpointSet.VpcId != nil {
						endpointSetMap["vpc_id"] = endpointSet.VpcId
					}
					if endpointSet.SubnetId != nil {
						endpointSetMap["subnet_id"] = endpointSet.SubnetId
					}
					if endpointSet.PrivateIp != nil {
						endpointSetMap["private_ip"] = endpointSet.PrivateIp
					}
					if endpointSet.PrivatePort != nil {
						endpointSetMap["private_port"] = endpointSet.PrivatePort
					}
					if endpointSet.WanIp != nil {
						endpointSetMap["wan_ip"] = endpointSet.WanIp
					}
					if endpointSet.WanPort != nil {
						endpointSetMap["wan_port"] = endpointSet.WanPort
					}
					if endpointSet.WanDomain != nil {
						endpointSetMap["wan_domain"] = endpointSet.WanDomain
					}

					endpointSetList = append(endpointSetList, endpointSetMap)
				}
				clusterSetMap["endpoint_set"] = endpointSetList
			}
			if cluster.DBMajorVersion != nil {
				clusterSetMap["db_major_version"] = cluster.DBMajorVersion
			}
			if cluster.DBKernelVersion != nil {
				clusterSetMap["db_kernel_version"] = cluster.DBKernelVersion
			}
			if cluster.StoragePayMode != nil {
				clusterSetMap["storage_pay_mode"] = cluster.StoragePayMode
			}
			ids = append(ids, *cluster.ClusterId)
			clusterList = append(clusterList, clusterSetMap)
		}
		d.SetId(helper.DataResourceIdsHash(ids))
		err := d.Set("list", clusterList)
		if err != nil {
			log.Printf("[CRITICAL]%s set tdcpg clusterList failed, reason:%+v", logId, err)
			return err
		}
	}

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), clusterList); e != nil {
			return e
		}
	}

	return nil
}
