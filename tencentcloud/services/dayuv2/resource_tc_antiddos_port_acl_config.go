package dayuv2

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svcantiddos "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/antiddos"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	antiddos "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/antiddos/v20200309"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudAntiddosPortAclConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudAntiddosPortAclConfigCreate,
		Read:   resourceTencentCloudAntiddosPortAclConfigRead,
		Delete: resourceTencentCloudAntiddosPortAclConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "InstanceIdList。",
			},

			"acl_config": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "端口 ACL Policy。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"forward_protocol": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "协议 类型，可以 take TCP，udp，all 值。",
						},
						"d_port_start": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Starting 从 端口，使用 范围 的 0~65535 值。",
						},
						"d_port_end": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "end 从 端口，使用 范围 的 0~65535 值。",
						},
						"s_port_start": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Starting 从 来源 端口，使用 值 范围 的 0~65535。",
						},
						"s_port_end": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "end 从 来源 端口，使用 值 范围 的 0~65535。",
						},
						"action": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "操作，可以 take 值: drop，transmit，forward。",
						},
						"priority": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "策略 优先级， smaller 数量， higher 级别，和 higher matching 的 规则，使用 值 ranging 从 1 到 1000. 注意：此字段可能返回 null，表示无法获取有效值。",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudAntiddosPortAclConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_antiddos_port_acl_config.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request    = antiddos.NewCreatePortAclConfigRequest()
		instanceId string
	)
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		request.InstanceId = helper.String(instanceId)
	}
	var config string
	if dMap, ok := helper.InterfacesHeadMap(d, "acl_config"); ok {
		aclConfig := antiddos.AclConfig{}
		if v, ok := dMap["forward_protocol"]; ok {
			aclConfig.ForwardProtocol = helper.String(v.(string))
		}
		if v, ok := dMap["d_port_start"]; ok {
			aclConfig.DPortStart = helper.IntUint64(v.(int))
		}
		if v, ok := dMap["d_port_end"]; ok {
			aclConfig.DPortEnd = helper.IntUint64(v.(int))
		}
		if v, ok := dMap["s_port_start"]; ok {
			aclConfig.SPortStart = helper.IntUint64(v.(int))
		}
		if v, ok := dMap["s_port_end"]; ok {
			aclConfig.SPortEnd = helper.IntUint64(v.(int))
		}
		if v, ok := dMap["action"]; ok {
			aclConfig.Action = helper.String(v.(string))
		}
		if v, ok := dMap["priority"]; ok {
			aclConfig.Priority = helper.IntUint64(v.(int))
		}
		request.AclConfig = &aclConfig
		dMapJson, err := json.Marshal(dMap)
		if err != nil {
			return err
		}
		config = string(dMapJson)
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseAntiddosClient().CreatePortAclConfig(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create antiddos portAclConfig failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(instanceId + tccommon.FILED_SP + config)

	return resourceTencentCloudAntiddosPortAclConfigRead(d, meta)
}

func resourceTencentCloudAntiddosPortAclConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_antiddos_port_acl_config.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := svcantiddos.NewAntiddosService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", idSplit)
	}
	instanceId := idSplit[0]
	configJson := idSplit[1]

	portAclConfigs, err := service.DescribeAntiddosPortAclConfigById(ctx, instanceId)
	if err != nil {
		return err
	}

	if portAclConfigs == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `AntiddosPortAclConfig` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	configMap := make(map[string]interface{})
	err = json.Unmarshal([]byte(configJson), &configMap)
	if err != nil {
		return err
	}

	var targetConfig *antiddos.AclConfig
	for _, item := range portAclConfigs {
		portAclConfig := item.AclConfig
		if v, ok := configMap["forward_protocol"]; ok {
			if *portAclConfig.ForwardProtocol != v.(string) {
				continue
			}
		}
		if v, ok := configMap["d_port_start"]; ok {
			if int(*portAclConfig.DPortStart) != int(v.(float64)) {
				continue
			}
		}
		if v, ok := configMap["d_port_end"]; ok {
			if int(*portAclConfig.DPortEnd) != int(v.(float64)) {
				continue
			}
		}
		if v, ok := configMap["s_port_start"]; ok {
			if int(*portAclConfig.SPortStart) != int(v.(float64)) {
				continue
			}
		}
		if v, ok := configMap["s_port_end"]; ok {
			if int(*portAclConfig.SPortEnd) != int(v.(float64)) {
				continue
			}
		}
		if v, ok := configMap["action"]; ok {
			if *portAclConfig.Action != v.(string) {
				continue
			}
		}
		if v, ok := configMap["priority"]; ok {
			if int(*portAclConfig.Priority) != int(v.(float64)) {
				continue
			}
		}
		targetConfig = portAclConfig
	}

	_ = d.Set("instance_id", instanceId)

	if targetConfig != nil {
		aclConfigMap := map[string]interface{}{}

		if targetConfig.ForwardProtocol != nil {
			aclConfigMap["forward_protocol"] = targetConfig.ForwardProtocol
		}

		if targetConfig.DPortStart != nil {
			aclConfigMap["d_port_start"] = targetConfig.DPortStart
		}

		if targetConfig.DPortEnd != nil {
			aclConfigMap["d_port_end"] = targetConfig.DPortEnd
		}

		if targetConfig.SPortStart != nil {
			aclConfigMap["s_port_start"] = targetConfig.SPortStart
		}

		if targetConfig.SPortEnd != nil {
			aclConfigMap["s_port_end"] = targetConfig.SPortEnd
		}

		if targetConfig.Action != nil {
			aclConfigMap["action"] = targetConfig.Action
		}

		if targetConfig.Priority != nil {
			aclConfigMap["priority"] = targetConfig.Priority
		}

		_ = d.Set("acl_config", []interface{}{aclConfigMap})
	}

	return nil
}

func resourceTencentCloudAntiddosPortAclConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_antiddos_port_acl_config.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := svcantiddos.NewAntiddosService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", idSplit)
	}
	instanceId := idSplit[0]
	configJson := idSplit[1]

	configMap := make(map[string]interface{})
	err := json.Unmarshal([]byte(configJson), &configMap)
	if err != nil {
		return err
	}

	aclConfig := antiddos.AclConfig{}
	if v, ok := configMap["forward_protocol"]; ok {
		aclConfig.ForwardProtocol = helper.String(v.(string))
	}
	if v, ok := configMap["d_port_start"]; ok {
		aclConfig.DPortStart = helper.IntUint64(int(v.(float64)))
	}
	if v, ok := configMap["d_port_end"]; ok {
		aclConfig.DPortEnd = helper.IntUint64(int(v.(float64)))
	}
	if v, ok := configMap["s_port_start"]; ok {
		aclConfig.SPortStart = helper.IntUint64(int(v.(float64)))
	}
	if v, ok := configMap["s_port_end"]; ok {
		aclConfig.SPortEnd = helper.IntUint64(int(v.(float64)))
	}
	if v, ok := configMap["action"]; ok {
		aclConfig.Action = helper.String(v.(string))
	}
	if v, ok := configMap["priority"]; ok {
		aclConfig.Priority = helper.IntUint64(int(v.(float64)))
	}

	if err := service.DeleteAntiddosPortAclConfigById(ctx, instanceId, &aclConfig); err != nil {
		return err
	}

	return nil
}
