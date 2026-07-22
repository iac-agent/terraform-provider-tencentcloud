package apigateway

import (
	"context"
	"fmt"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudAPIGatewayCustomDomain() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudAPIGatewayCustomDomainCreate,
		Read:   resourceTencentCloudAPIGatewayCustomDomainRead,
		Update: resourceTencentCloudAPIGatewayCustomDomainUpdate,
		Delete: resourceTencentCloudAPIGatewayCustomDomainDelete,

		Schema: map[string]*schema.Schema{
			"service_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: tccommon.ValidateNotEmpty,
				ForceNew:     true,
				Description:  "Unique 服务 ID",
			},
			"sub_domain": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Custom 域名 名称 to be bound。",
			},
			"protocol": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: tccommon.ValidateNotEmpty,
				Description:  "协议 supported by service. 有效值：`http`，`https`，`http&https`。",
			},
			"net_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: tccommon.ValidateNotEmpty,
				Description:  "Network 类型 有效值：`OUTER`，`INNER`。",
			},
			"is_default_mapping": {
				Type:        schema.TypeBool,
				Default:     true,
				Optional:    true,
				Description: "是否default 路径 mapping is used. The 默认值为 `true`. When it is `false`，it means custom 路径 mapping. In this case，the `path_mappings` attribute 为必填项。",
			},
			"default_domain": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: tccommon.ValidateNotEmpty,
				Description:  "Default 域名 名称",
			},
			"certificate_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Unique 证书 ID custom 域名 名称 to be bound. You can choose to upload for the `协议` attribute 值 `https` or `http&https`。",
			},
			"path_mappings": {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Custom 域名 名称 路径 mapping. The data 格式 is: `路径#environment`. 可选 values for the environment are `test`，`prepub`，and `release`。",
			},
			"is_forced_https": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "是否force HTTP requests to jump to HTTPS，默认为 false. When the parameter is true，the API gateway will redirect all HTTP 协议 requests using the custom 域名 名称 to the HTTPS 协议 for forwarding。",
			},
			//compute
			"status": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "域名 名称 resolution 状态 `1` means normal analysis，`0` means parsing failed。",
			},
		},
	}
}

func resourceTencentCloudAPIGatewayCustomDomainCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_api_gateway_custom_domain.create")()

	var (
		logId             = tccommon.GetLogId(tccommon.ContextNil)
		ctx               = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		apiGatewayService = APIGatewayService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		serviceId         = d.Get("service_id").(string)
		subDomain         = d.Get("sub_domain").(string)
		protocol          = d.Get("protocol").(string)
		netType           = d.Get("net_type").(string)
		defaultDomain     = d.Get("default_domain").(string)
		isDefaultMapping  = d.Get("is_default_mapping").(bool)
		isForcedHttps     = d.Get("is_forced_https").(bool)
		certificateId     string
		pathMappings      []string
		err               error
	)

	if v, ok := d.GetOk("certificate_id"); ok {
		certificateId = v.(string)
	}
	if v, ok := d.GetOk("path_mappings"); ok {
		pathMappings = helper.InterfacesStrings(v.(*schema.Set).List())
	}

	err = apiGatewayService.BindSubDomainService(ctx, serviceId, subDomain, protocol, netType, defaultDomain, isDefaultMapping, certificateId, pathMappings, isForcedHttps)
	if err != nil {
		return err
	}

	d.SetId(strings.Join([]string{serviceId, subDomain}, tccommon.FILED_SP))

	return resourceTencentCloudAPIGatewayCustomDomainRead(d, meta)
}

func resourceTencentCloudAPIGatewayCustomDomainRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_api_gateway_custom_domain.read")()

	var (
		logId = tccommon.GetLogId(tccommon.ContextNil)
		ctx   = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		apiGatewayService = APIGatewayService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		id                = d.Id()
		err               error
	)

	results := strings.Split(id, tccommon.FILED_SP)
	if len(results) != 2 {
		return fmt.Errorf("ids param is error. id:  %s", id)
	}
	serviceId := results[0]
	subDomain := results[1]
	resultList, err := apiGatewayService.DescribeServiceSubDomainsService(ctx, serviceId, subDomain)
	if err != nil {
		return err
	}

	if len(resultList) == 0 {
		d.SetId("")
		return nil
	}

	resultInfo := resultList[0]
	info, err := apiGatewayService.DescribeServiceSubDomainMappings(ctx, serviceId, *resultInfo.DomainName)
	if err != nil {
		return fmt.Errorf("DescribeServiceSubDomainMappings err: %s", err.Error())
	}
	pathMap := make([]string, 0, len(info.PathMappingSet))
	for _, v := range info.PathMappingSet {
		pathMap = append(pathMap, strings.Join([]string{*v.Path, *v.Environment}, tccommon.FILED_SP))
	}

	_ = d.Set("path_mappings", pathMap)
	_ = d.Set("status", resultInfo.Status)
	_ = d.Set("certificate_id", resultInfo.CertificateId)
	_ = d.Set("is_default_mapping", resultInfo.IsDefaultMapping)
	_ = d.Set("protocol", resultInfo.Protocol)
	_ = d.Set("net_type", resultInfo.NetType)
	_ = d.Set("service_id", serviceId)
	_ = d.Set("is_forced_https", resultInfo.IsForcedHttps)

	return nil
}

func resourceTencentCloudAPIGatewayCustomDomainUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_api_gateway_custom_domain.update")()

	var (
		logId             = tccommon.GetLogId(tccommon.ContextNil)
		ctx               = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		apiGatewayService = APIGatewayService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		id                = d.Id()
		subDomain         string
		isDefaultMapping  bool
		certificateId     string
		protocol          string
		netType           string
		pathMappings      []string
		isForcedHttps     bool
		hasChange         bool
	)

	results := strings.Split(id, tccommon.FILED_SP)
	if len(results) != 2 {
		return fmt.Errorf("ids param is error. setId:  %s", id)
	}
	serviceId := results[0]

	subDomain = d.Get("sub_domain").(string)
	if d.HasChange("sub_domain") {
		hasChange = true
	}

	isDefaultMapping = d.Get("is_default_mapping").(bool)
	if d.HasChange("is_default_mapping") {
		hasChange = true
	}

	if v, ok := d.GetOk("certificate_id"); ok {
		certificateId = v.(string)
	}
	if d.HasChange("certificate_id") {
		hasChange = true
	}

	netType = d.Get("net_type").(string)
	if d.HasChange("net_type") {
		hasChange = true
	}

	protocol = d.Get("protocol").(string)
	if d.HasChange("protocol") {
		hasChange = true
	}

	if v, ok := d.GetOk("path_mappings"); ok {
		pathMappings = helper.InterfacesStrings(v.(*schema.Set).List())
	}

	if d.HasChange("path_mappings") {
		hasChange = true
	}

	isForcedHttps = d.Get("is_forced_https").(bool)
	if d.HasChange("is_forced_https") {
		hasChange = true
	}

	if hasChange {
		err := apiGatewayService.ModifySubDomainService(ctx, serviceId, subDomain, isDefaultMapping, certificateId, protocol, netType, pathMappings, isForcedHttps)
		if err != nil {
			return err
		}
	}

	return resourceTencentCloudAPIGatewayCustomDomainRead(d, meta)
}

func resourceTencentCloudAPIGatewayCustomDomainDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_api_gateway_custom_domain.delete")()

	var (
		logId             = tccommon.GetLogId(tccommon.ContextNil)
		ctx               = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		id                = d.Id()
		apigatewayService = APIGatewayService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	results := strings.Split(id, tccommon.FILED_SP)
	if len(results) != 2 {
		return fmt.Errorf("ids param is error. setId:  %s", id)
	}
	serviceId := results[0]
	subDomain := results[1]

	return apigatewayService.UnBindSubDomainService(ctx, serviceId, subDomain)
}
