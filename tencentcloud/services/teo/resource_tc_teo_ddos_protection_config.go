package teo

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudTeoDdosProtectionConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTeoDdosProtectionConfigCreate,
		Read:   resourceTencentCloudTeoDdosProtectionConfigRead,
		Update: resourceTencentCloudTeoDdosProtectionConfigUpdate,
		Delete: resourceTencentCloudTeoDdosProtectionConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "可用区 ID",
			},

			"ddos_protection": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Description: "指定exclusive Anti-DDoS 配置。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"protection_option": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "指定protection 范围 的 standalone DDoS. 有效 值:.\n<li>protect_all_domains: 指定exclusive Anti-DDoS protection 对于 all 域名 names 在 site. newly added 域名 names automatically 启用 exclusive Anti-DDoS protection. 当 此 参数 是 指定，DomainDDoSProtections 将 不 是 processed.</li>.\n<li>protect_specified_domains: 仅 applicable 到 指定 domains. 特定 范围 可以 是 集合 via DomainDDoSProtection 参数.</li>。",
						},
						"domain_ddos_protections": {
							Type:        schema.TypeSet,
							Optional:    true,
							Computed:    true,
							Description: "Anti-DDoS 配置 的 域名 指定exclusive ddos protection settings 对于 域名 在 请求 参数.\n<li>当 ProtectionOption remains protect_specified_domains， 域名 names 不 filled 在 keep their exclusive Anti-DDoS protection 配置 unchanged，while explicitly 指定 域名 names 是 更新 according 到 input 参数.</li>.\n<li>当 ProtectionOption switches 从 protect_all_domains 到 protect_specified_domains: 如果 DomainDDoSProtections 是 空，disable exclusive DDoS protection 对于 all domains under site; 如果 DomainDDoSProtections 是 不 空，disable 或 maintain exclusive DDoS protection 对于 域名 names 指定 在 参数，和 disable exclusive DDoS protection 对于 other unlisted 域名 names.</li>。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"domain": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "域名 名称",
									},
									"switch": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Standalone DDoS switch 的 域名 有效 值:.\n<li>在: 已启用;</li>.\n<li>关闭: closed.</li>。",
									},
								},
							},
						},
						"shared_cname_ddos_protections": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "指定exclusive DDoS protection 配置 的 shared CNAME. 使用 作为 output 参数。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"domain": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "域名 名称",
									},
									"switch": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Standalone DDoS switch 的 域名 有效 值:.\n<li>在: 已启用;</li>.\n<li>关闭: closed.</li>。",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudTeoDdosProtectionConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_ddos_protection_config.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var zoneId string
	if v, ok := d.GetOk("zone_id"); ok {
		zoneId = v.(string)
	}

	d.SetId(zoneId)
	return resourceTencentCloudTeoDdosProtectionConfigUpdate(d, meta)
}

func resourceTencentCloudTeoDdosProtectionConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_ddos_protection_config.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = TeoService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		zoneId  = d.Id()
	)

	respData, err := service.DescribeTeoDdosProtectionConfigById(ctx, zoneId)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[WARN]%s resource `tencentcloud_teo_ddos_protection_config` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	_ = d.Set("zone_id", zoneId)

	dMap := make(map[string]interface{}, 0)
	if respData.ProtectionOption != nil {
		dMap["protection_option"] = respData.ProtectionOption
	}

	if respData.DomainDDoSProtections != nil {
		domainDDoSProtectionsList := make([]map[string]interface{}, 0, len(respData.DomainDDoSProtections))
		for _, domainDDoSProtections := range respData.DomainDDoSProtections {
			domainDDoSProtectionsMap := map[string]interface{}{}
			if domainDDoSProtections.Domain != nil {
				domainDDoSProtectionsMap["domain"] = domainDDoSProtections.Domain
			}

			if domainDDoSProtections.Switch != nil {
				domainDDoSProtectionsMap["switch"] = domainDDoSProtections.Switch
			}

			domainDDoSProtectionsList = append(domainDDoSProtectionsList, domainDDoSProtectionsMap)
		}

		dMap["domain_ddos_protections"] = domainDDoSProtectionsList
	}

	if respData.SharedCNAMEDDoSProtections != nil {
		sharedCNAMEDDoSProtectionsList := make([]map[string]interface{}, 0, len(respData.SharedCNAMEDDoSProtections))
		for _, sharedCNAMEDDoSProtections := range respData.SharedCNAMEDDoSProtections {
			sharedCNAMEDDoSProtectionsMap := map[string]interface{}{}
			if sharedCNAMEDDoSProtections.Domain != nil {
				sharedCNAMEDDoSProtectionsMap["domain"] = sharedCNAMEDDoSProtections.Domain
			}

			if sharedCNAMEDDoSProtections.Switch != nil {
				sharedCNAMEDDoSProtectionsMap["switch"] = sharedCNAMEDDoSProtections.Switch
			}

			sharedCNAMEDDoSProtectionsList = append(sharedCNAMEDDoSProtectionsList, sharedCNAMEDDoSProtectionsMap)
		}

		dMap["shared_cname_ddos_protections"] = sharedCNAMEDDoSProtectionsList
	}

	_ = d.Set("ddos_protection", []interface{}{dMap})
	return nil
}

func resourceTencentCloudTeoDdosProtectionConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_ddos_protection_config.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = teov20220901.NewModifyDDoSProtectionRequest()
		zoneId  = d.Id()
	)

	if dDoSProtectionMap, ok := helper.InterfacesHeadMap(d, "ddos_protection"); ok {
		dDoSProtection := teov20220901.DDoSProtection{}
		if v, ok := dDoSProtectionMap["protection_option"].(string); ok && v != "" {
			dDoSProtection.ProtectionOption = helper.String(v)
		}

		if v, ok := dDoSProtectionMap["domain_ddos_protections"]; ok {
			for _, item := range v.(*schema.Set).List() {
				domainDDoSProtectionsMap := item.(map[string]interface{})
				domainDDoSProtection := teov20220901.DomainDDoSProtection{}
				if v, ok := domainDDoSProtectionsMap["domain"].(string); ok && v != "" {
					domainDDoSProtection.Domain = helper.String(v)
				}

				if v, ok := domainDDoSProtectionsMap["switch"].(string); ok && v != "" {
					domainDDoSProtection.Switch = helper.String(v)
				}

				dDoSProtection.DomainDDoSProtections = append(dDoSProtection.DomainDDoSProtections, &domainDDoSProtection)
			}
		}

		request.DDoSProtection = &dDoSProtection
	}

	request.ZoneId = &zoneId
	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().ModifyDDoSProtectionWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s update teo ddos protection config failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return resourceTencentCloudTeoDdosProtectionConfigRead(d, meta)
}

func resourceTencentCloudTeoDdosProtectionConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_ddos_protection_config.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
