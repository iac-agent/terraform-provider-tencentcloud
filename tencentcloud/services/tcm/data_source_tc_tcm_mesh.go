package tcm

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tcm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tcm/v20210413"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudTcmMesh() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTcmMeshRead,
		Schema: map[string]*schema.Schema{
			"mesh_id": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Mesh 实例 ID。",
			},

			"mesh_name": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Display 名称",
			},

			"tags": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "标签键",
			},

			"mesh_cluster": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Mesh 名称",
			},

			"mesh_list": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "mesh 信息 是 queried注意：此字段可能返回 null，表示有效值不可用。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"mesh_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Mesh 实例 ID。",
						},
						"display_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Mesh 名称",
						},
						"version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Mesh 版本",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Mesh 类型 值 范围:- `STANDALONE`: Standalone mesh- `HOSTED`: hosted mesh。",
						},
						"config": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Mesh 配置。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"istio": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Istio 配置。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"outbound_traffic_policy": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Outbound 流量 策略。",
												},
												"disable_policy_checks": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Disable 策略 checks。",
												},
												"enable_pilot_http": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Enable HTTP/1.0 support。",
												},
												"disable_http_retry": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Disable http retry。",
												},
												"smart_dns": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "SmartDNS 配置。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"istio_meta_dns_capture": {
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Enable dns proxy。",
															},
															"istio_meta_dns_auto_allocate": {
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Enable auto allocate 地址",
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
						"tag_list": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "A 列表 associated 标签",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "标签键",
									},
									"value": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "标签值",
									},
									"passthrough": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Passthrough 到 other related product。",
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

func dataSourceTencentCloudTcmMeshRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_tcm_mesh.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string][]*string)
	if v, ok := d.GetOk("mesh_id"); ok {
		meshIdSet := v.(*schema.Set).List()
		paramMap["MeshId"] = helper.InterfacesStringsPoint(meshIdSet)
	}

	if v, ok := d.GetOk("mesh_name"); ok {
		meshName := v.(*schema.Set).List()
		paramMap["MeshName"] = helper.InterfacesStringsPoint(meshName)
	}

	if v, ok := d.GetOk("tags"); ok {
		tagsSet := v.(*schema.Set).List()
		paramMap["Tags"] = helper.InterfacesStringsPoint(tagsSet)
	}

	if v, ok := d.GetOk("mesh_cluster"); ok {
		meshClusterSet := v.(*schema.Set).List()
		paramMap["MeshCluster"] = helper.InterfacesStringsPoint(meshClusterSet)
	}

	service := TcmService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var meshList []*tcm.Mesh

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeTcmMeshByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		meshList = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(meshList))
	tmpList := make([]map[string]interface{}, 0, len(meshList))

	if meshList != nil {
		for _, mesh := range meshList {
			meshMap := map[string]interface{}{}

			if mesh.MeshId != nil {
				meshMap["mesh_id"] = mesh.MeshId
			}

			if mesh.DisplayName != nil {
				meshMap["display_name"] = mesh.DisplayName
			}

			if mesh.Version != nil {
				meshMap["version"] = mesh.Version
			}

			if mesh.Type != nil {
				meshMap["type"] = mesh.Type
			}

			if mesh.Config != nil {
				configMap := map[string]interface{}{}

				if mesh.Config.Istio != nil {
					istioMap := map[string]interface{}{}

					if mesh.Config.Istio.OutboundTrafficPolicy != nil {
						istioMap["outbound_traffic_policy"] = mesh.Config.Istio.OutboundTrafficPolicy
					}

					if mesh.Config.Istio.DisablePolicyChecks != nil {
						istioMap["disable_policy_checks"] = mesh.Config.Istio.DisablePolicyChecks
					}

					if mesh.Config.Istio.EnablePilotHTTP != nil {
						istioMap["enable_pilot_http"] = mesh.Config.Istio.EnablePilotHTTP
					}

					if mesh.Config.Istio.DisableHTTPRetry != nil {
						istioMap["disable_http_retry"] = mesh.Config.Istio.DisableHTTPRetry
					}

					if mesh.Config.Istio.SmartDNS != nil {
						smartDNSMap := map[string]interface{}{}

						if mesh.Config.Istio.SmartDNS.IstioMetaDNSCapture != nil {
							smartDNSMap["istio_meta_dns_capture"] = mesh.Config.Istio.SmartDNS.IstioMetaDNSCapture
						}

						if mesh.Config.Istio.SmartDNS.IstioMetaDNSAutoAllocate != nil {
							smartDNSMap["istio_meta_dns_auto_allocate"] = mesh.Config.Istio.SmartDNS.IstioMetaDNSAutoAllocate
						}

						istioMap["smart_dns"] = []interface{}{smartDNSMap}
					}

					configMap["istio"] = []interface{}{istioMap}
				}

				meshMap["config"] = []interface{}{configMap}
			}

			if mesh.TagList != nil {
				tagListList := []interface{}{}
				for _, tagList := range mesh.TagList {
					tagListMap := map[string]interface{}{}

					if tagList.Key != nil {
						tagListMap["key"] = tagList.Key
					}

					if tagList.Value != nil {
						tagListMap["value"] = tagList.Value
					}

					if tagList.Passthrough != nil {
						tagListMap["passthrough"] = tagList.Passthrough
					}

					tagListList = append(tagListList, tagListMap)
				}

				meshMap["tag_list"] = tagListList
			}

			ids = append(ids, *mesh.MeshId)
			tmpList = append(tmpList, meshMap)
		}

		err := d.Set("mesh_list", tmpList)
		if err != nil {
			return err
		}
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
