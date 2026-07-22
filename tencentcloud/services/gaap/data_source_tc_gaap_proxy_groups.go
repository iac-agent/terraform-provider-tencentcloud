package gaap

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	gaap "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/gaap/v20180529"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudGaapProxyGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudGaapProxyGroupsRead,
		Schema: map[string]*schema.Schema{
			"project_id": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "项目 ID 值 范围:-1，All projects under 此 user0，默认值 projectOther 值，指定 items。",
			},

			"filters": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "过滤器 conditions, upper 限制 的 过滤器.Values per 请求 是 5.RealServerRegion - String - 必填: No - (filtering criteria) 过滤器 通过 real 服务器 地域，refer 到 RegionId 在 返回 results 的 DescribeDestRegions interface.PackageType - String - 必填: No - (过滤器 condition) proxy 组 类型，其中 &amp;#39;Thunder&amp;#39; 表示 standard proxy 组 和 &amp;#39;Accelerator&amp;#39; 表示 silver acceleration proxy 组。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "过滤器 conditions。",
						},
						"values": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Required:    true,
							Description: "filtering 值",
						},
					},
				},
			},

			"tag_set": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "标签列表，当 此 字段 exists，pulls 资源 列表 under corresponding 标签Supports 最大 的 5 labels. 当 there 是 two 或 more labels 和 any 一个 的 them 是 met， proxy 组 将 是 pulled out。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tag_key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "标签键",
						},
						"tag_value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "标签值",
						},
					},
				},
			},

			"proxy_group_list": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "列表 proxy groups.注意：此字段可能返回 null，表示无法获取有效值。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"group_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "proxy 组 ID。",
						},
						"domain": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "proxy 组 域名 name注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"group_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "proxy Group Name注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"project_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "项目 ID",
						},
						"real_server_region_info": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Real Server 地域 Info。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"region_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "地域 ID",
									},
									"region_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "地域 名称",
									},
									"region_area": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "地域 的 computer room。",
									},
									"region_area_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "地域名称 的 computer room。",
									},
									"idc_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "类型 computer room，其中 &#39;dc&#39; 表示 DataCenter 数据 center 和 &#39;ec&#39; 表示 EdgeComputing edge 节点。",
									},
									"feature_bitmap": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Property bitmap，其中 each bit 表示 属性，其中:0，表示that 功能 是 不 支持;1，表示support 对于 此 功能. meaning 的 功能 bitmap 是 作为 follows (从 right 到 left): first bit 支持 4-layer acceleration; second bit 支持 7-layer acceleration; third bit 支持 Http3 访问; fourth bit 支持 IPv6; fifth bit 支持 high-quality BGP 访问; 6th bit 支持 three 网络 访问; 7th bit 支持 QoS acceleration 在 访问 segment.注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"support_feature": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Ability 到 访问 regional support注意：此字段可能返回 null，表示无法获取有效值。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"network_type": {
													Type: schema.TypeSet,
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
													Computed:    true,
													Description: "A 列表 网络 types 支持 通过 访问 area，使用 &#39;normal&#39; indicating support 对于 regular BGP，&#39;cn2&#39; indicating premium BGP，&#39;triple&#39; indicating three networks，和 &#39;secure_EIP&#39; 表示 自定义 secure EIP。",
												},
											},
										},
									},
								},
							},
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "proxy 组 状态Among them,&#39;RUNNING&#39; 表示running;&#39;CREATING&#39; 表示being 创建;&#39;DESTROYING&#39; 表示being destroyed;&#39;MOVING&#39; 表示that proxy 是 being migrated;&#39;CHANGING&#39; 表示partial 部署。",
						},
						"tag_set": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "标签 Set。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"tag_key": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "标签键",
									},
									"tag_value": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "标签值",
									},
								},
							},
						},
						"version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "proxy Group Version注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"create_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Create Time注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"proxy_type": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Does proxy 组 include Microsoft proxys注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"http3_supported": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Supports identification 的 Http3 features，其中:0 表示shutdown;1 表示enabled.注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"feature_bitmap": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Property bitmap，其中 each bit 表示 属性，其中:0，表示that 功能 是 不 支持;1，表示support 对于 此 功能. meaning 的 功能 bitmap 是 作为 follows (从 right 到 left): first bit 支持 4-layer acceleration; second bit 支持 7-layer acceleration; third bit 支持 Http3 访问; fourth bit 支持 IPv6; fifth bit 支持 high-quality BGP 访问; 6th bit 支持 three 网络 访问; 7th bit 支持 QoS acceleration 在 访问 segment.注意：此字段可能返回 null，表示无法获取有效值。",
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

func dataSourceTencentCloudGaapProxyGroupsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_gaap_proxy_groups.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, _ := d.GetOk("project_id"); v != nil {
		paramMap["ProjectId"] = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("filters"); ok {
		filtersSet := v.([]interface{})
		tmpSet := make([]*gaap.Filter, 0, len(filtersSet))

		for _, item := range filtersSet {
			filter := gaap.Filter{}
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

	if v, ok := d.GetOk("tag_set"); ok {
		tagSetSet := v.([]interface{})
		tmpSet := make([]*gaap.TagPair, 0, len(tagSetSet))

		for _, item := range tagSetSet {
			tagPair := gaap.TagPair{}
			tagPairMap := item.(map[string]interface{})

			if v, ok := tagPairMap["tag_key"]; ok {
				tagPair.TagKey = helper.String(v.(string))
			}
			if v, ok := tagPairMap["tag_value"]; ok {
				tagPair.TagValue = helper.String(v.(string))
			}
			tmpSet = append(tmpSet, &tagPair)
		}
		paramMap["TagSet"] = tmpSet
	}

	service := GaapService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var proxyGroupList []*gaap.ProxyGroupInfo

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeGaapProxyGroupsByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		proxyGroupList = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(proxyGroupList))
	tmpList := make([]map[string]interface{}, 0, len(proxyGroupList))

	if proxyGroupList != nil {
		for _, proxyGroupInfo := range proxyGroupList {
			proxyGroupInfoMap := map[string]interface{}{}

			if proxyGroupInfo.GroupId != nil {
				proxyGroupInfoMap["group_id"] = proxyGroupInfo.GroupId
			}

			if proxyGroupInfo.Domain != nil {
				proxyGroupInfoMap["domain"] = proxyGroupInfo.Domain
			}

			if proxyGroupInfo.GroupName != nil {
				proxyGroupInfoMap["group_name"] = proxyGroupInfo.GroupName
			}

			if proxyGroupInfo.ProjectId != nil {
				proxyGroupInfoMap["project_id"] = proxyGroupInfo.ProjectId
			}

			if proxyGroupInfo.RealServerRegionInfo != nil {
				realServerRegionInfoMap := map[string]interface{}{}

				if proxyGroupInfo.RealServerRegionInfo.RegionId != nil {
					realServerRegionInfoMap["region_id"] = proxyGroupInfo.RealServerRegionInfo.RegionId
				}

				if proxyGroupInfo.RealServerRegionInfo.RegionName != nil {
					realServerRegionInfoMap["region_name"] = proxyGroupInfo.RealServerRegionInfo.RegionName
				}

				if proxyGroupInfo.RealServerRegionInfo.RegionArea != nil {
					realServerRegionInfoMap["region_area"] = proxyGroupInfo.RealServerRegionInfo.RegionArea
				}

				if proxyGroupInfo.RealServerRegionInfo.RegionAreaName != nil {
					realServerRegionInfoMap["region_area_name"] = proxyGroupInfo.RealServerRegionInfo.RegionAreaName
				}

				if proxyGroupInfo.RealServerRegionInfo.IDCType != nil {
					realServerRegionInfoMap["idc_type"] = proxyGroupInfo.RealServerRegionInfo.IDCType
				}

				if proxyGroupInfo.RealServerRegionInfo.FeatureBitmap != nil {
					realServerRegionInfoMap["feature_bitmap"] = proxyGroupInfo.RealServerRegionInfo.FeatureBitmap
				}

				if proxyGroupInfo.RealServerRegionInfo.SupportFeature != nil {
					supportFeatureMap := map[string]interface{}{}

					if proxyGroupInfo.RealServerRegionInfo.SupportFeature.NetworkType != nil {
						supportFeatureMap["network_type"] = proxyGroupInfo.RealServerRegionInfo.SupportFeature.NetworkType
					}

					realServerRegionInfoMap["support_feature"] = []interface{}{supportFeatureMap}
				}

				proxyGroupInfoMap["real_server_region_info"] = []interface{}{realServerRegionInfoMap}
			}

			if proxyGroupInfo.Status != nil {
				proxyGroupInfoMap["status"] = proxyGroupInfo.Status
			}

			if proxyGroupInfo.TagSet != nil {
				tagSetList := []interface{}{}
				for _, tagSet := range proxyGroupInfo.TagSet {
					tagSetMap := map[string]interface{}{}

					if tagSet.TagKey != nil {
						tagSetMap["tag_key"] = tagSet.TagKey
					}

					if tagSet.TagValue != nil {
						tagSetMap["tag_value"] = tagSet.TagValue
					}

					tagSetList = append(tagSetList, tagSetMap)
				}

				proxyGroupInfoMap["tag_set"] = tagSetList
			}

			if proxyGroupInfo.Version != nil {
				proxyGroupInfoMap["version"] = proxyGroupInfo.Version
			}

			if proxyGroupInfo.CreateTime != nil {
				proxyGroupInfoMap["create_time"] = proxyGroupInfo.CreateTime
			}

			if proxyGroupInfo.ProxyType != nil {
				proxyGroupInfoMap["proxy_type"] = proxyGroupInfo.ProxyType
			}

			if proxyGroupInfo.Http3Supported != nil {
				proxyGroupInfoMap["http3_supported"] = proxyGroupInfo.Http3Supported
			}

			if proxyGroupInfo.FeatureBitmap != nil {
				proxyGroupInfoMap["feature_bitmap"] = proxyGroupInfo.FeatureBitmap
			}

			ids = append(ids, *proxyGroupInfo.GroupId)
			tmpList = append(tmpList, proxyGroupInfoMap)
		}

		_ = d.Set("proxy_group_list", tmpList)
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
