package oceanus

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	oceanus "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/oceanus/v20190422"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudOceanusWorkSpaces() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudOceanusWorkSpacesRead,
		Schema: map[string]*schema.Schema{
			"order_type": {
				Optional:     true,
				Type:         schema.TypeInt,
				Default:      WORK_SPACES_ORDER_TYPE_0,
				ValidateFunc: tccommon.ValidateAllowedIntValue(WORK_SPACES_ORDER_TYPE),
				Description:  "1:排序方式 创建时间 在 降序 (默认值); 2:排序方式 创建时间 在 升序; 3:排序方式 状态 在 降序; 4:排序方式 状态 在 升序; 默认为 0。",
			},
			"filters": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "过滤器 规则。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "待过滤字段",
						},
						"values": {
							Type:        schema.TypeSet,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Required:    true,
							Description: "过滤器 值 对于 字段。",
						},
					},
				},
			},
			"work_space_set_item": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "列表 workspace details注意：此字段可能返回 null，表示无法获取有效值。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"serial_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Workspace SerialId。",
						},
						"app_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "用户 APPID。",
						},
						"owner_uin": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Main 账号 UIN。",
						},
						"creator_uin": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "创建者 UIN。",
						},
						"work_space_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Workspace 名称",
						},
						"region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "地域",
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
						"status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "1:uninitialized; 2:可用; -1:删除。",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Workspace 描述",
						},
						"cluster_group_set_item": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Workspace 集群 信息。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"cluster_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "SerialId 的 clusterGroup。",
									},
									"name": {
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
									"app_id": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "账号 APPID。",
									},
									"owner_uin": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Main 账号 UIN。",
									},
									"creator_uin": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "创建者 账号 UIN。",
									},
									"cu_num": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "CU quantity。",
									},
									"cu_mem": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "CU 内存 规格。",
									},
									"status": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Cluster 状态，1:uninitialized，3:initializing，2:running。",
									},
									"status_desc": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "状态 描述",
									},
									"create_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Cluster 创建时间。",
									},
									"update_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Last operation 时间 在 集群。",
									},
									"remark": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "描述",
									},
									"net_environment_type": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Network。",
									},
									"free_cu_num": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Free CU。",
									},
									"free_cu": {
										Type:        schema.TypeFloat,
										Computed:    true,
										Description: "Free CU under fine-grained resources。",
									},
									"running_cu": {
										Type:        schema.TypeFloat,
										Computed:    true,
										Description: "Running CU。",
									},
									"pay_mode": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Payment 模式",
									},
								},
							},
						},
						"role_auth": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Workspace 角色 信息。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"app_id": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "用户 AppID。",
									},
									"work_space_serial_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Workspace SerialId。",
									},
									"owner_uin": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Main 账号 UIN。",
									},
									"creator_uin": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "创建者 UIN。",
									},
									"auth_sub_account_uin": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Bound authorized UIN。",
									},
									"permission": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Corresponding 到 ID 在 角色 表。",
									},
									"create_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "创建时间。",
									},
									"update_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Last operation 时间。",
									},
									"status": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "2:已启用，1:已禁用",
									},
									"id": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "ID注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"work_space_id": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Workspace ID注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"role_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Permission name注意：此字段可能返回 null，表示无法获取有效值。",
									},
								},
							},
						},
						"role_auth_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Workspace member count。",
						},
						"work_space_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Workspace SerialId。",
						},
						"jobs_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "注意：此字段可能返回 null，表示无法获取有效值。",
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

func dataSourceTencentCloudOceanusWorkSpacesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_oceanus_work_spaces.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId         = tccommon.GetLogId(tccommon.ContextNil)
		ctx           = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service       = OceanusService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		workSpaceList []*oceanus.WorkSpaceSetItem
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOkExists("order_type"); ok {
		paramMap["OrderType"] = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("filters"); ok {
		filtersSet := v.([]interface{})
		tmpSet := make([]*oceanus.Filter, 0, len(filtersSet))

		for _, item := range filtersSet {
			filter := oceanus.Filter{}
			filterMap := item.(map[string]interface{})

			if v, ok := filterMap["name"]; ok {
				filter.Name = helper.String(v.(string))
			}

			if v, ok := filterMap["values"]; ok {
				valuesSet := v.(*schema.Set).List()
				filter.Values = helper.InterfacesStringsPoint(valuesSet)
			}

			tmpSet = append(tmpSet, &filter)
		}

		paramMap["Filters"] = tmpSet
	}

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeOceanusWorkSpacesByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		workSpaceList = result
		return nil
	})

	if err != nil {
		return err
	}

	ids := make([]string, 0, len(workSpaceList))
	tmpList := make([]map[string]interface{}, 0, len(workSpaceList))
	if workSpaceList != nil {
		for _, workSpaceSetItem := range workSpaceList {
			workSpaceSetItemMap := map[string]interface{}{}

			if workSpaceSetItem.SerialId != nil {
				workSpaceSetItemMap["serial_id"] = workSpaceSetItem.SerialId
			}

			if workSpaceSetItem.AppId != nil {
				workSpaceSetItemMap["app_id"] = workSpaceSetItem.AppId
			}

			if workSpaceSetItem.OwnerUin != nil {
				workSpaceSetItemMap["owner_uin"] = workSpaceSetItem.OwnerUin
			}

			if workSpaceSetItem.CreatorUin != nil {
				workSpaceSetItemMap["creator_uin"] = workSpaceSetItem.CreatorUin
			}

			if workSpaceSetItem.WorkSpaceName != nil {
				workSpaceSetItemMap["work_space_name"] = workSpaceSetItem.WorkSpaceName
			}

			if workSpaceSetItem.Region != nil {
				workSpaceSetItemMap["region"] = workSpaceSetItem.Region
			}

			if workSpaceSetItem.CreateTime != nil {
				workSpaceSetItemMap["create_time"] = workSpaceSetItem.CreateTime
			}

			if workSpaceSetItem.UpdateTime != nil {
				workSpaceSetItemMap["update_time"] = workSpaceSetItem.UpdateTime
			}

			if workSpaceSetItem.Status != nil {
				workSpaceSetItemMap["status"] = workSpaceSetItem.Status
			}

			if workSpaceSetItem.Description != nil {
				workSpaceSetItemMap["description"] = workSpaceSetItem.Description
			}

			if workSpaceSetItem.ClusterGroupSetItem != nil {
				clusterGroupSetItemList := []interface{}{}
				for _, clusterGroupSetItem := range workSpaceSetItem.ClusterGroupSetItem {
					clusterGroupSetItemMap := map[string]interface{}{}

					if clusterGroupSetItem.ClusterId != nil {
						clusterGroupSetItemMap["cluster_id"] = clusterGroupSetItem.ClusterId
					}

					if clusterGroupSetItem.Name != nil {
						clusterGroupSetItemMap["name"] = clusterGroupSetItem.Name
					}

					if clusterGroupSetItem.Region != nil {
						clusterGroupSetItemMap["region"] = clusterGroupSetItem.Region
					}

					if clusterGroupSetItem.Zone != nil {
						clusterGroupSetItemMap["zone"] = clusterGroupSetItem.Zone
					}

					if clusterGroupSetItem.AppId != nil {
						clusterGroupSetItemMap["app_id"] = clusterGroupSetItem.AppId
					}

					if clusterGroupSetItem.OwnerUin != nil {
						clusterGroupSetItemMap["owner_uin"] = clusterGroupSetItem.OwnerUin
					}

					if clusterGroupSetItem.CreatorUin != nil {
						clusterGroupSetItemMap["creator_uin"] = clusterGroupSetItem.CreatorUin
					}

					if clusterGroupSetItem.CuNum != nil {
						clusterGroupSetItemMap["cu_num"] = clusterGroupSetItem.CuNum
					}

					if clusterGroupSetItem.CuMem != nil {
						clusterGroupSetItemMap["cu_mem"] = clusterGroupSetItem.CuMem
					}

					if clusterGroupSetItem.Status != nil {
						clusterGroupSetItemMap["status"] = clusterGroupSetItem.Status
					}

					if clusterGroupSetItem.StatusDesc != nil {
						clusterGroupSetItemMap["status_desc"] = clusterGroupSetItem.StatusDesc
					}

					if clusterGroupSetItem.CreateTime != nil {
						clusterGroupSetItemMap["create_time"] = clusterGroupSetItem.CreateTime
					}

					if clusterGroupSetItem.UpdateTime != nil {
						clusterGroupSetItemMap["update_time"] = clusterGroupSetItem.UpdateTime
					}

					if clusterGroupSetItem.Remark != nil {
						clusterGroupSetItemMap["remark"] = clusterGroupSetItem.Remark
					}

					if clusterGroupSetItem.NetEnvironmentType != nil {
						clusterGroupSetItemMap["net_environment_type"] = clusterGroupSetItem.NetEnvironmentType
					}

					if clusterGroupSetItem.FreeCuNum != nil {
						clusterGroupSetItemMap["free_cu_num"] = clusterGroupSetItem.FreeCuNum
					}

					if clusterGroupSetItem.FreeCu != nil {
						clusterGroupSetItemMap["free_cu"] = clusterGroupSetItem.FreeCu
					}

					if clusterGroupSetItem.RunningCu != nil {
						clusterGroupSetItemMap["running_cu"] = clusterGroupSetItem.RunningCu
					}

					if clusterGroupSetItem.PayMode != nil {
						clusterGroupSetItemMap["pay_mode"] = clusterGroupSetItem.PayMode
					}

					clusterGroupSetItemList = append(clusterGroupSetItemList, clusterGroupSetItemMap)
				}

				workSpaceSetItemMap["cluster_group_set_item"] = clusterGroupSetItemList
			}

			if workSpaceSetItem.RoleAuth != nil {
				roleAuthList := []interface{}{}
				for _, roleAuth := range workSpaceSetItem.RoleAuth {
					roleAuthMap := map[string]interface{}{}

					if roleAuth.AppId != nil {
						roleAuthMap["app_id"] = roleAuth.AppId
					}

					if roleAuth.WorkSpaceSerialId != nil {
						roleAuthMap["work_space_serial_id"] = roleAuth.WorkSpaceSerialId
					}

					if roleAuth.OwnerUin != nil {
						roleAuthMap["owner_uin"] = roleAuth.OwnerUin
					}

					if roleAuth.CreatorUin != nil {
						roleAuthMap["creator_uin"] = roleAuth.CreatorUin
					}

					if roleAuth.AuthSubAccountUin != nil {
						roleAuthMap["auth_sub_account_uin"] = roleAuth.AuthSubAccountUin
					}

					if roleAuth.Permission != nil {
						roleAuthMap["permission"] = roleAuth.Permission
					}

					if roleAuth.CreateTime != nil {
						roleAuthMap["create_time"] = roleAuth.CreateTime
					}

					if roleAuth.UpdateTime != nil {
						roleAuthMap["update_time"] = roleAuth.UpdateTime
					}

					if roleAuth.Status != nil {
						roleAuthMap["status"] = roleAuth.Status
					}

					if roleAuth.Id != nil {
						roleAuthMap["id"] = roleAuth.Id
					}

					if roleAuth.WorkSpaceId != nil {
						roleAuthMap["work_space_id"] = roleAuth.WorkSpaceId
					}

					if roleAuth.RoleName != nil {
						roleAuthMap["role_name"] = roleAuth.RoleName
					}

					roleAuthList = append(roleAuthList, roleAuthMap)
				}

				workSpaceSetItemMap["role_auth"] = roleAuthList
			}

			if workSpaceSetItem.RoleAuthCount != nil {
				workSpaceSetItemMap["role_auth_count"] = workSpaceSetItem.RoleAuthCount
			}

			if workSpaceSetItem.WorkSpaceId != nil {
				workSpaceSetItemMap["work_space_id"] = workSpaceSetItem.WorkSpaceId
			}

			if workSpaceSetItem.JobsCount != nil {
				workSpaceSetItemMap["jobs_count"] = workSpaceSetItem.JobsCount
			}

			ids = append(ids, *workSpaceSetItem.WorkSpaceId)
			tmpList = append(tmpList, workSpaceSetItemMap)
		}

		_ = d.Set("work_space_set_item", tmpList)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), d); e != nil {
			return e
		}
	}

	return nil
}
