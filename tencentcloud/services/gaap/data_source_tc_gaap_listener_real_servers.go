package gaap

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	gaap "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/gaap/v20180529"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudGaapListenerRealServers() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudGaapListenerRealServersRead,
		Schema: map[string]*schema.Schema{
			"listener_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "listener ID。",
			},

			"real_server_set": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Real Server Set。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"real_server_i_p": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Real Server IP。",
						},
						"real_server_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Real Server Id。",
						},
						"real_server_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Real Server 名称",
						},
						"project_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "项目 ID",
						},
						"in_ban_blacklist": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Is it on the banned blacklist? 0 表示not on the blacklist，and 1 表示on the blacklist。",
						},
					},
				},
			},

			"bind_real_server_set": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Bound real server Information List。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"real_server_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Real Server Id。",
						},
						"real_server_i_p": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Real Server IP or 域名",
						},
						"real_server_weight": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The 权重 of this real server。",
						},
						"real_server_status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "real server health check 状态，where:0 表示normal;1 表示an exception.When the health check 状态 is not 已启用，it is always normal.注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"real_server_port": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The 端口 数量 the real server注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"down_i_p_list": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Computed:    true,
							Description: "When the real server is a 域名 名称，the 域名 名称 is resolved to one or more IPs，and this field represents the 列表 abnormal IPs. When the 状态 is abnormal，but the field is empty，it 表示that the 域名 名称 resolution is abnormal。",
						},
						"real_server_failover_role": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The primary and secondary roles of the real server，&#39;master&#39; represents primary，&#39;slave&#39; represents secondary，and this parameter must be in the 活跃 and standby 模式 of the real server when the listener is turned on。",
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

func dataSourceTencentCloudGaapListenerRealServersRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_gaap_listener_real_servers.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	var listenerId string

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("listener_id"); ok {
		listenerId = v.(string)
		paramMap["ListenerId"] = helper.String(listenerId)
	}

	service := GaapService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var (
		realServerSet     []*gaap.RealServer
		bindRealServerSet []*gaap.BindRealServer
	)

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		resultRealServerSet, resultBindRealServerSet, e := service.DescribeGaapListenerRealServersByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		realServerSet = resultRealServerSet
		bindRealServerSet = resultBindRealServerSet
		return nil
	})
	if err != nil {
		return err
	}

	tmpRealServerList := make([]map[string]interface{}, 0, len(realServerSet))
	tmpBindRealServerList := make([]map[string]interface{}, 0, len(bindRealServerSet))

	if realServerSet != nil {
		for _, realServer := range realServerSet {
			realServerMap := map[string]interface{}{}

			if realServer.RealServerIP != nil {
				realServerMap["real_server_i_p"] = realServer.RealServerIP
			}

			if realServer.RealServerId != nil {
				realServerMap["real_server_id"] = realServer.RealServerId
			}

			if realServer.RealServerName != nil {
				realServerMap["real_server_name"] = realServer.RealServerName
			}

			if realServer.ProjectId != nil {
				realServerMap["project_id"] = realServer.ProjectId
			}

			if realServer.InBanBlacklist != nil {
				realServerMap["in_ban_blacklist"] = realServer.InBanBlacklist
			}

			tmpRealServerList = append(tmpRealServerList, realServerMap)
		}

		_ = d.Set("real_server_set", tmpRealServerList)
	}

	if bindRealServerSet != nil {
		for _, bindRealServer := range bindRealServerSet {
			bindRealServerMap := map[string]interface{}{}

			if bindRealServer.RealServerId != nil {
				bindRealServerMap["real_server_id"] = bindRealServer.RealServerId
			}

			if bindRealServer.RealServerIP != nil {
				bindRealServerMap["real_server_i_p"] = bindRealServer.RealServerIP
			}

			if bindRealServer.RealServerWeight != nil {
				bindRealServerMap["real_server_weight"] = bindRealServer.RealServerWeight
			}

			if bindRealServer.RealServerStatus != nil {
				bindRealServerMap["real_server_status"] = bindRealServer.RealServerStatus
			}

			if bindRealServer.RealServerPort != nil {
				bindRealServerMap["real_server_port"] = bindRealServer.RealServerPort
			}

			if bindRealServer.DownIPList != nil {
				bindRealServerMap["down_i_p_list"] = bindRealServer.DownIPList
			}

			if bindRealServer.RealServerFailoverRole != nil {
				bindRealServerMap["real_server_failover_role"] = bindRealServer.RealServerFailoverRole
			}

			tmpBindRealServerList = append(tmpBindRealServerList, bindRealServerMap)
		}

		_ = d.Set("bind_real_server_set", tmpBindRealServerList)
	}

	d.SetId(listenerId)
	result := map[string]interface{}{
		"real_server_set":      tmpRealServerList,
		"bind_real_server_set": tmpBindRealServerList,
	}
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), result); e != nil {
			return e
		}
	}
	return nil
}
