package cynosdb

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cynosdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cynosdb/v20190107"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudCynosdbClusterDetailDatabases() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCynosdbClusterDetailDatabasesRead,
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "集群 ID。",
			},
			"db_name": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "数据库名称。",
			},
			"db_infos": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "数据库信息说明：该字段可能返回null，表示无法获取有效值。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"db_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "数据库名称。",
						},
						"character_set": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "字符集类型。",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "数据库状态。",
						},
						"collate_rule": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "捕获规则。",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "数据库说明：该字段可能返回null，表示无法获取到有效值。",
						},
						"user_host_privileges": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "用户权限说明：该字段可能返回null，表示无法获取到有效值。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"db_user_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "数据库用户名。",
									},
									"db_host": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "数据库主机。",
									},
									"db_privilege": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "用户权限说明：该字段可能返回null，表示无法获取到有效值。",
									},
								},
							},
						},
						"db_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "数据库ID 注意：该字段可能返回null，表示无法获取有效值。",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "创建时注意：该字段可能返回null，表示无法获取有效值。",
						},
						"update_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "更新时间注意：该字段可能返回null，表示无法获取到有效值。",
						},
						"app_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "用户appid注意：该字段可能返回null，表示无法获取有效值。",
						},
						"uin": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "用户 Uin 注意：该字段可能返回null，表示无法获取有效值。",
						},
						"cluster_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Cluster Id 注意：该字段可能返回null，表示无法获取有效值。",
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

func dataSourceTencentCloudCynosdbClusterDetailDatabasesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_cynosdb_cluster_detail_databases.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId     = tccommon.GetLogId(tccommon.ContextNil)
		ctx       = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service   = CynosdbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		dbInfos   []*cynosdb.DbInfo
		clusterId string
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("cluster_id"); ok {
		paramMap["ClusterId"] = helper.String(v.(string))
		clusterId = v.(string)
	}

	if v, ok := d.GetOk("db_name"); ok {
		paramMap["DbName"] = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeCynosdbClusterDetailDatabasesByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		dbInfos = result
		return nil
	})

	if err != nil {
		return err
	}

	tmpList := make([]map[string]interface{}, 0, len(dbInfos))

	if dbInfos != nil {
		for _, dbInfo := range dbInfos {
			dbInfoMap := map[string]interface{}{}

			if dbInfo.DbName != nil {
				dbInfoMap["db_name"] = dbInfo.DbName
			}

			if dbInfo.CharacterSet != nil {
				dbInfoMap["character_set"] = dbInfo.CharacterSet
			}

			if dbInfo.Status != nil {
				dbInfoMap["status"] = dbInfo.Status
			}

			if dbInfo.CollateRule != nil {
				dbInfoMap["collate_rule"] = dbInfo.CollateRule
			}

			if dbInfo.Description != nil {
				dbInfoMap["description"] = dbInfo.Description
			}

			if dbInfo.UserHostPrivileges != nil {
				userHostPrivilegesList := []interface{}{}
				for _, userHostPrivileges := range dbInfo.UserHostPrivileges {
					userHostPrivilegesMap := map[string]interface{}{}

					if userHostPrivileges.DbUserName != nil {
						userHostPrivilegesMap["db_user_name"] = userHostPrivileges.DbUserName
					}

					if userHostPrivileges.DbHost != nil {
						userHostPrivilegesMap["db_host"] = userHostPrivileges.DbHost
					}

					if userHostPrivileges.DbPrivilege != nil {
						userHostPrivilegesMap["db_privilege"] = userHostPrivileges.DbPrivilege
					}

					userHostPrivilegesList = append(userHostPrivilegesList, userHostPrivilegesMap)
				}

				dbInfoMap["user_host_privileges"] = userHostPrivilegesList
			}

			if dbInfo.DbId != nil {
				dbInfoMap["db_id"] = dbInfo.DbId
			}

			if dbInfo.CreateTime != nil {
				dbInfoMap["create_time"] = dbInfo.CreateTime
			}

			if dbInfo.UpdateTime != nil {
				dbInfoMap["update_time"] = dbInfo.UpdateTime
			}

			if dbInfo.AppId != nil {
				dbInfoMap["app_id"] = dbInfo.AppId
			}

			if dbInfo.Uin != nil {
				dbInfoMap["uin"] = dbInfo.Uin
			}

			if dbInfo.ClusterId != nil {
				dbInfoMap["cluster_id"] = dbInfo.ClusterId
			}

			tmpList = append(tmpList, dbInfoMap)
		}

		_ = d.Set("db_infos", tmpList)
	}

	d.SetId(clusterId)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), tmpList); e != nil {
			return e
		}
	}

	return nil
}
