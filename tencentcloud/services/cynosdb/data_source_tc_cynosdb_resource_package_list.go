package cynosdb

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cynosdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cynosdb/v20190107"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudCynosdbResourcePackageList() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCynosdbResourcePackageListRead,
		Schema: map[string]*schema.Schema{
			"package_id": {
				Optional:    true,
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "资源包唯一ID。",
			},
			"package_name": {
				Optional:    true,
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "资源包名称。",
			},
			"package_type": {
				Optional:    true,
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "资源包类型 CCU - 计算资源包，DISK - 存储资源包。",
			},
			"package_region": {
				Optional:    true,
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "资源包使用地区 中国-中国大陆通用，海外-港澳台、海外通用。",
			},
			"status": {
				Optional:    true,
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "资源包状态创建-创建；使用 - 使用中； Expired——已经过期； Normal_ Finish——用完； Apply_ Refund——申请退款；退款 - 费用已退还。",
			},
			"order_by": {
				Optional:    true,
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "支持的排序条件：startTime - 有效时间、expireTime - 过期时间、packageUsedSpec - 使用容量、packageTotalSpec - 总存储容量。按数组顺序排列；。",
			},
			"order_direction": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "排序依据：DESC 降序、ASC 升序。",
			},
			"resource_package_list": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "资源包详情说明：该字段可能返回null，表示无法获取到有效值。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"app_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "AppID注意：该字段可能返回null，表示无法获取有效值。",
						},
						"package_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "资源包唯一ID 注意：该字段可能返回null，表示无法获取有效值。",
						},
						"package_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "资源包名称 注意：该字段可能返回null，表示无法获取有效值。",
						},
						"package_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "资源包类型 CCU - 计算资源包，DISK - 存储资源包 注意：该字段可能返回 null，表示无法获取有效值。",
						},
						"package_region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "资源包国内使用，常用的是中国大陆，海外使用，常用的是港澳台，海外。注意：该字段可能返回null，表示无法获取到有效值。",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "资源包状态创建-创建；使用 - 使用中； Expired——已经过期； Normal_ Finish——用完； Apply_ Refund——申请退款；退款 - 费用已退还。注意：该字段可能返回null，表示无法获取到有效值。",
						},
						"package_total_spec": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "注意资源包总量：该字段可能返回null，表示无法获取有效值。",
						},
						"package_used_spec": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "资源包使用注意：该字段可能返回null，表示无法获取到有效值。",
						},
						"has_quota": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "资源包使用注意：该字段可能返回null，表示无法获取到有效值。",
						},
						"bind_instance_infos": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "绑定实例信息注意：该字段可能返回null，表示无法获取到有效值。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"instance_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "实例ID。",
									},
									"instance_region": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "实例所在区域。",
									},
									"instance_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "实例类型。",
									},
								},
							},
						},
						"start_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "生效时间：2022年7月1日00:00:00 注意：该字段可能返回null，表示无法获取有效值。",
						},
						"expire_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "过期时间：2022年8月1日00:00:00 注意：该字段可能返回null，表示无法获取有效值。",
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

func dataSourceTencentCloudCynosdbResourcePackageListRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_cynosdb_resource_package_list.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service    = CynosdbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		detailList []*cynosdb.Package
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("package_id"); ok {
		packageIdSet := v.(*schema.Set).List()
		paramMap["PackageId"] = helper.InterfacesStringsPoint(packageIdSet)
	}

	if v, ok := d.GetOk("package_name"); ok {
		packageNameSet := v.(*schema.Set).List()
		paramMap["PackageName"] = helper.InterfacesStringsPoint(packageNameSet)
	}

	if v, ok := d.GetOk("package_type"); ok {
		packageTypeSet := v.(*schema.Set).List()
		paramMap["PackageType"] = helper.InterfacesStringsPoint(packageTypeSet)
	}

	if v, ok := d.GetOk("package_region"); ok {
		packageRegionSet := v.(*schema.Set).List()
		paramMap["PackageRegion"] = helper.InterfacesStringsPoint(packageRegionSet)
	}

	if v, ok := d.GetOk("status"); ok {
		statusSet := v.(*schema.Set).List()
		paramMap["Status"] = helper.InterfacesStringsPoint(statusSet)
	}

	if v, ok := d.GetOk("order_by"); ok {
		orderBySet := v.(*schema.Set).List()
		paramMap["OrderBy"] = helper.InterfacesStringsPoint(orderBySet)
	}

	if v, ok := d.GetOk("order_direction"); ok {
		paramMap["OrderDirection"] = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeCynosdbResourcePackageListByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		detailList = result
		return nil
	})

	if err != nil {
		return err
	}

	ids := make([]string, 0, len(detailList))

	if detailList != nil {
		tmpList := []interface{}{}
		for _, detail := range detailList {
			detailMap := map[string]interface{}{}
			if detail.AppId != nil {
				detailMap["app_id"] = detail.AppId
			}

			if detail.PackageId != nil {
				detailMap["package_id"] = detail.PackageId
			}

			if detail.PackageName != nil {
				detailMap["package_name"] = detail.PackageName
			}

			if detail.PackageType != nil {
				detailMap["package_type"] = detail.PackageType
			}

			if detail.PackageRegion != nil {
				detailMap["package_region"] = detail.PackageRegion
			}

			if detail.Status != nil {
				detailMap["status"] = detail.Status
			}

			if detail.PackageTotalSpec != nil {
				detailMap["package_total_spec"] = detail.PackageTotalSpec
			}

			if detail.PackageUsedSpec != nil {
				detailMap["package_used_spec"] = detail.PackageUsedSpec
			}

			if detail.HasQuota != nil {
				detailMap["has_quota"] = detail.HasQuota
			}

			if detail.BindInstanceInfos != nil {
				insList := []interface{}{}
				for _, instanceInfo := range detail.BindInstanceInfos {
					insMap := map[string]interface{}{}
					if instanceInfo.InstanceId != nil {
						insMap["instance_id"] = instanceInfo.InstanceId
					}

					if instanceInfo.InstanceRegion != nil {
						insMap["instance_region"] = instanceInfo.InstanceRegion
					}

					if instanceInfo.InstanceType != nil {
						insMap["instance_type"] = instanceInfo.InstanceType
					}
					insList = append(insList, insMap)
				}

				detailMap["bind_instance_infos"] = insList
			}

			if detail.StartTime != nil {
				detailMap["start_time"] = detail.StartTime
			}

			if detail.ExpireTime != nil {
				detailMap["expire_time"] = detail.ExpireTime
			}

			tmpList = append(tmpList, detailMap)
		}

		_ = d.Set("resource_package_list", tmpList)
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
