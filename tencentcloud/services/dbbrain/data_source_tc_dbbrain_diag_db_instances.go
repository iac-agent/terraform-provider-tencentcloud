package dbbrain

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dbbrain "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dbbrain/v20210527"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudDbbrainDiagDbInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudDbbrainDiagDbInstancesRead,
		Schema: map[string]*schema.Schema{
			"is_supported": {
				Required:    true,
				Type:        schema.TypeBool,
				Description: "whether 它 是 实例 支持 通过 DBbrain, always pass `true`.",
			},

			"product": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "服务 product 类型, 支持 值 include: `mysql` - 云 数据库 MySQL, `cynosdb` - 云 数据库 TDSQL-C 对于 MySQL, 默认值 是 `mysql`.",
			},

			"instance_names": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "查询 based 在 实例 名称 condition.",
			},

			"instance_ids": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "查询 based 在 实例 ID condition.",
			},

			"regions": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "查询 based 在 geographical conditions.",
			},

			"db_scan_status": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "all-实例 inspection 状态. `0`: All-实例 inspection 是 已启用; `1`: All-实例 inspection 是 不 已启用.",
			},

			"items": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "信息 about 实例.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "实例 ID.",
						},
						"instance_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "实例 名称.",
						},
						"region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "地域.",
						},
						"health_score": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "health score.",
						},
						"product": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "belongs 到 product.",
						},
						"event_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "数量 的 abnormal events.",
						},
						"instance_type": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "实例 类型. `1`: MASTER; `2`: DR, `3`: RO, `4`: SDR.",
						},
						"cpu": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "数量 的 cores.",
						},
						"memory": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "内存, 在 MB.",
						},
						"volume": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "hard 磁盘 存储, 在 GB.",
						},
						"engine_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "数据库 版本.",
						},
						"vip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "intranet 地址.",
						},
						"vport": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "intranet 端口.",
						},
						"source": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "访问 source.",
						},
						"group_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "组 ID.",
						},
						"group_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "组 名称.",
						},
						"status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "实例 状态: `0`: Shipping; `1`: Running normally; `4`: Destroying; `5`: Isolating.",
						},
						"uniq_subnet_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "子网 uniform ID.",
						},
						"deploy_mode": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "cdb 类型.",
						},
						"init_flag": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "cdb 实例 initialization flag: `0`: 不 initialized; `1`: initialized.",
						},
						"task_status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "任务 状态.",
						},
						"uniq_vpc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "unified ID 的 私有 网络.",
						},
						"instance_conf": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "状态 的 实例 inspection/overview.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"daily_inspection": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "数据库 inspection switch, Yes/No.",
									},
									"overview_display": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "实例 overview switch, Yes/No.",
									},
									"key_delimiters": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Computed:    true,
										Description: "Custom separator 对于 redis large 键 analysis, 仅 使用 通过 `redis`. 注意: 此 字段 可能 返回 null, indicating 该 无 有效 值 可以 是 获取.",
									},
								},
							},
						},
						"deadline_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "资源 expiration 时间.",
						},
						"is_supported": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "whether 它 是 实例 支持 通过 DBbrain.",
						},
						"sec_audit_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "已启用 状态 的 实例 安全 audit 日志. `ON`: 安全 audit 是 已启用; `OFF`: 安全 audit 是 不 已启用.",
						},
						"audit_policy_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "实例 audit 日志 启用 状态. `ALL_AUDIT`: full audit 是 已启用; `RULE_AUDIT`: 规则 audit 是 已启用; `UNBOUND`: audit 是 不 已启用.",
						},
						"audit_running_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "实例 audit 日志 running 状态. `normal`: running; `paused`: arrears suspended.",
						},
						"internal_vip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Intranet VIPNote: 此 字段 可能 返回 null, indicating 该 无 有效 值 可以 是 获取.",
						},
						"internal_vport": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Intranet portNote: 此 字段 可能 返回 null, indicating 该 无 有效 值 可以 是 获取.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "create 时间.",
						},
					},
				},
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used 到 save results.",
			},
		},
	}
}

func dataSourceTencentCloudDbbrainDiagDbInstancesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_dbbrain_diag_db_instances.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, _ := d.GetOk("is_supported"); v != nil {
		paramMap["IsSupported"] = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOk("product"); ok {
		paramMap["Product"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("instance_names"); ok {
		instanceNamesSet := v.(*schema.Set).List()
		paramMap["InstanceNames"] = helper.InterfacesStringsPoint(instanceNamesSet)
	}

	if v, ok := d.GetOk("instance_ids"); ok {
		instanceIdsSet := v.(*schema.Set).List()
		paramMap["InstanceIds"] = helper.InterfacesStringsPoint(instanceIdsSet)
	}

	if v, ok := d.GetOk("regions"); ok {
		regionsSet := v.(*schema.Set).List()
		paramMap["Regions"] = helper.InterfacesStringsPoint(regionsSet)
	}

	service := DbbrainService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var (
		infos        []*dbbrain.InstanceInfo
		dbScanStatus *int64
		e            error
	)
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		infos, dbScanStatus, e = service.DescribeDbbrainDiagDbInstancesByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(infos))
	tmpList := make([]map[string]interface{}, 0, len(infos))
	if dbScanStatus != nil {
		_ = d.Set("db_scan_status", dbScanStatus)
	}

	if infos != nil {

		for _, instanceInfo := range infos {
			instanceInfoMap := map[string]interface{}{}

			if instanceInfo.InstanceId != nil {
				instanceInfoMap["instance_id"] = instanceInfo.InstanceId
			}

			if instanceInfo.InstanceName != nil {
				instanceInfoMap["instance_name"] = instanceInfo.InstanceName
			}

			if instanceInfo.Region != nil {
				instanceInfoMap["region"] = instanceInfo.Region
			}

			if instanceInfo.HealthScore != nil {
				instanceInfoMap["health_score"] = instanceInfo.HealthScore
			}

			if instanceInfo.Product != nil {
				instanceInfoMap["product"] = instanceInfo.Product
			}

			if instanceInfo.EventCount != nil {
				instanceInfoMap["event_count"] = instanceInfo.EventCount
			}

			if instanceInfo.InstanceType != nil {
				instanceInfoMap["instance_type"] = instanceInfo.InstanceType
			}

			if instanceInfo.Cpu != nil {
				instanceInfoMap["cpu"] = instanceInfo.Cpu
			}

			if instanceInfo.Memory != nil {
				instanceInfoMap["memory"] = instanceInfo.Memory
			}

			if instanceInfo.Volume != nil {
				instanceInfoMap["volume"] = instanceInfo.Volume
			}

			if instanceInfo.EngineVersion != nil {
				instanceInfoMap["engine_version"] = instanceInfo.EngineVersion
			}

			if instanceInfo.Vip != nil {
				instanceInfoMap["vip"] = instanceInfo.Vip
			}

			if instanceInfo.Vport != nil {
				instanceInfoMap["vport"] = instanceInfo.Vport
			}

			if instanceInfo.Source != nil {
				instanceInfoMap["source"] = instanceInfo.Source
			}

			if instanceInfo.GroupId != nil {
				instanceInfoMap["group_id"] = instanceInfo.GroupId
			}

			if instanceInfo.GroupName != nil {
				instanceInfoMap["group_name"] = instanceInfo.GroupName
			}

			if instanceInfo.Status != nil {
				instanceInfoMap["status"] = instanceInfo.Status
			}

			if instanceInfo.UniqSubnetId != nil {
				instanceInfoMap["uniq_subnet_id"] = instanceInfo.UniqSubnetId
			}

			if instanceInfo.DeployMode != nil {
				instanceInfoMap["deploy_mode"] = instanceInfo.DeployMode
			}

			if instanceInfo.InitFlag != nil {
				instanceInfoMap["init_flag"] = instanceInfo.InitFlag
			}

			if instanceInfo.TaskStatus != nil {
				instanceInfoMap["task_status"] = instanceInfo.TaskStatus
			}

			if instanceInfo.UniqVpcId != nil {
				instanceInfoMap["uniq_vpc_id"] = instanceInfo.UniqVpcId
			}

			if instanceInfo.InstanceConf != nil {
				instanceConfMap := map[string]interface{}{}

				if instanceInfo.InstanceConf.DailyInspection != nil {
					instanceConfMap["daily_inspection"] = instanceInfo.InstanceConf.DailyInspection
				}

				if instanceInfo.InstanceConf.OverviewDisplay != nil {
					instanceConfMap["overview_display"] = instanceInfo.InstanceConf.OverviewDisplay
				}

				if instanceInfo.InstanceConf.KeyDelimiters != nil {
					instanceConfMap["key_delimiters"] = instanceInfo.InstanceConf.KeyDelimiters
				}

				instanceInfoMap["instance_conf"] = []interface{}{instanceConfMap}
			}

			if instanceInfo.DeadlineTime != nil {
				instanceInfoMap["deadline_time"] = instanceInfo.DeadlineTime
			}

			if instanceInfo.IsSupported != nil {
				instanceInfoMap["is_supported"] = instanceInfo.IsSupported
			}

			if instanceInfo.SecAuditStatus != nil {
				instanceInfoMap["sec_audit_status"] = instanceInfo.SecAuditStatus
			}

			if instanceInfo.AuditPolicyStatus != nil {
				instanceInfoMap["audit_policy_status"] = instanceInfo.AuditPolicyStatus
			}

			if instanceInfo.AuditRunningStatus != nil {
				instanceInfoMap["audit_running_status"] = instanceInfo.AuditRunningStatus
			}

			if instanceInfo.InternalVip != nil {
				instanceInfoMap["internal_vip"] = instanceInfo.InternalVip
			}

			if instanceInfo.InternalVport != nil {
				instanceInfoMap["internal_vport"] = instanceInfo.InternalVport
			}

			if instanceInfo.CreateTime != nil {
				instanceInfoMap["create_time"] = instanceInfo.CreateTime
			}

			ids = append(ids, *instanceInfo.InstanceId)
			tmpList = append(tmpList, instanceInfoMap)
		}

		_ = d.Set("items", tmpList)
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
