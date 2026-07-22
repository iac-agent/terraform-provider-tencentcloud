package gaap

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	gaap "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/gaap/v20180529"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudGaapRuleRealServers() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudGaapRuleRealServersRead,
		Schema: map[string]*schema.Schema{
			"rule_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Rule ID。",
			},

			"real_server_set": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Real Server Set。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"real_server_ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Real Server IP 或 域名",
						},
						"real_server_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Real Server ID。",
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
							Description: "Is 它 在 banned blacklist? 0 表示not 在 blacklist，和 1 表示on blacklist。",
						},
					},
				},
			},

			"bind_real_server_set": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Bind Real Server info。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"real_server_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Real Server ID。",
						},
						"real_server_ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Real Server Ip 或 域名",
						},
						"real_server_weight": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Real Server 权重",
						},
						"real_server_status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "RealServerStatus: 0 表示normal;1 表示an exception.当 health check 状态 是 不 已启用，它 是 always normal.注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"real_server_port": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Real Server Port注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"down_ip_list": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Computed:    true,
							Description: "当 real 服务器 是 域名 名称， 域名 名称 是 resolved 到 一个 或 more IPs，和 此 字段 表示 列表 abnormal IPs. 当 状态 是 abnormal，但 字段 是 空，它 表示that 域名 名称 resolution 是 abnormal。",
						},
						"real_server_failover_role": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "primary 和 secondary roles 的 real 服务器:master 表示 primary，slave 表示 secondary，和 此 参数 必须 是 在 活跃 和 standby 模式 的 real 服务器 当 listener 是 turned 在。",
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

func dataSourceTencentCloudGaapRuleRealServersRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_gaap_rule_real_servers.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("rule_id"); ok {
		paramMap["RuleId"] = helper.String(v.(string))
	}

	service := GaapService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var (
		realServerSet     []*gaap.RealServer
		bindRealServerSet []*gaap.BindRealServer
	)

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		resultRuleRealServers, resultBindRealServers, e := service.DescribeGaapRuleRealServersByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		realServerSet = resultRuleRealServers
		bindRealServerSet = resultBindRealServers
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(realServerSet))
	tmpRealServerList := make([]map[string]interface{}, 0, len(realServerSet))
	tmpBindRealServerList := make([]map[string]interface{}, 0, len(bindRealServerSet))

	if realServerSet != nil {
		for _, realServer := range realServerSet {
			realServerMap := map[string]interface{}{}

			if realServer.RealServerIP != nil {
				realServerMap["real_server_ip"] = realServer.RealServerIP
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

			ids = append(ids, *realServer.RealServerIP)
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
				bindRealServerMap["real_server_ip"] = bindRealServer.RealServerIP
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
				bindRealServerMap["down_ip_list"] = bindRealServer.DownIPList
			}

			if bindRealServer.RealServerFailoverRole != nil {
				bindRealServerMap["real_server_failover_role"] = bindRealServer.RealServerFailoverRole
			}

			tmpBindRealServerList = append(tmpBindRealServerList, bindRealServerMap)
		}

		_ = d.Set("bind_real_server_set", tmpBindRealServerList)
	}

	result := map[string]interface{}{
		"real_server_set":      tmpRealServerList,
		"bind_real_server_set": tmpBindRealServerList,
	}
	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), result); e != nil {
			return e
		}
	}
	return nil
}
