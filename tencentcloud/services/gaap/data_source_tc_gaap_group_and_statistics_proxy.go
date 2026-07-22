package gaap

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	gaap "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/gaap/v20180529"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudGaapGroupAndStatisticsProxy() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudGaapGroupAndStatisticsProxyRead,
		Schema: map[string]*schema.Schema{
			"project_id": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "项目 ID",
			},

			"group_set": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Channel group information that can be counted。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"group_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Channel 组 ID",
						},
						"group_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Channel 组名称",
						},
						"proxy_set": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Channel list in the proxy group。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"proxy_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Channel Id。",
									},
									"proxy_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Channel 名称",
									},
									"listener_list": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "listener list。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"listener_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "listener Id。",
												},
												"listener_name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "listener 名称",
												},
												"port": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "listened 端口",
												},
												"protocol": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Listener 协议 类型",
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

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "用于保存结果。",
			},
		},
	}
}

func dataSourceTencentCloudGaapGroupAndStatisticsProxyRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_gaap_group_and_statistics_proxy.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, _ := d.GetOk("project_id"); v != nil {
		paramMap["ProjectId"] = helper.IntUint64(v.(int))
	}

	service := GaapService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var groupSet []*gaap.GroupStatisticsInfo

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeGaapGroupAndStatisticsProxyByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		groupSet = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(groupSet))
	tmpList := make([]map[string]interface{}, 0, len(groupSet))

	if groupSet != nil {
		for _, groupStatisticsInfo := range groupSet {
			groupStatisticsInfoMap := map[string]interface{}{}

			if groupStatisticsInfo.GroupId != nil {
				groupStatisticsInfoMap["group_id"] = groupStatisticsInfo.GroupId
				ids = append(ids, *groupStatisticsInfo.GroupId)
			}

			if groupStatisticsInfo.GroupName != nil {
				groupStatisticsInfoMap["group_name"] = groupStatisticsInfo.GroupName
			}

			if groupStatisticsInfo.ProxySet != nil {
				proxySetList := []interface{}{}
				for _, proxySet := range groupStatisticsInfo.ProxySet {
					proxySetMap := map[string]interface{}{}

					if proxySet.ProxyId != nil {
						proxySetMap["proxy_id"] = proxySet.ProxyId
					}

					if proxySet.ProxyName != nil {
						proxySetMap["proxy_name"] = proxySet.ProxyName
					}

					if proxySet.ListenerList != nil {
						listenerListList := []interface{}{}
						for _, listenerList := range proxySet.ListenerList {
							listenerListMap := map[string]interface{}{}

							if listenerList.ListenerId != nil {
								listenerListMap["listener_id"] = listenerList.ListenerId
							}

							if listenerList.ListenerName != nil {
								listenerListMap["listener_name"] = listenerList.ListenerName
							}

							if listenerList.Port != nil {
								listenerListMap["port"] = listenerList.Port
							}

							if listenerList.Protocol != nil {
								listenerListMap["protocol"] = listenerList.Protocol
							}

							listenerListList = append(listenerListList, listenerListMap)
						}

						proxySetMap["listener_list"] = listenerListList
					}

					proxySetList = append(proxySetList, proxySetMap)
				}

				groupStatisticsInfoMap["proxy_set"] = proxySetList
			}

			tmpList = append(tmpList, groupStatisticsInfoMap)
		}

		_ = d.Set("group_set", tmpList)
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
