package dayu

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sdkError "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	dayu "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dayu/v20180709"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudDayuDdosPolicyCase() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDayuDdosPolicyCaseCreate,
		Read:   resourceTencentCloudDayuDdosPolicyCaseRead,
		Update: resourceTencentCloudDayuDdosPolicyCaseUpdate,
		Delete: resourceTencentCloudDayuDdosPolicyCaseDelete,

		Schema: map[string]*schema.Schema{
			"resource_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue(DAYU_RESOURCE_TYPE_CASE),
				ForceNew:     true,
				Description:  "类型 资源 该 DDoS 策略 case works 对于. 有效值：`bgpip`，`bgp` 和 `bgp-multip`。",
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: tccommon.ValidateStringLengthInRange(1, 64),
				Description:  "名称 DDoS 策略 case. Length should between 1 和 64。",
			},
			"platform_types": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: tccommon.ValidateAllowedStringValue(DAYU_APP_PLATFORM),
					Description:  "Platform 的 DDoS 策略 case. 有效值：`PC`，`MOBILE`，`TV` 和 `SERVER`。",
				},
				Required:    true,
				Description: "Platform 集合 的 DDoS 策略 case。",
			},
			"app_protocols": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: tccommon.ValidateAllowedStringValue(DAYU_PROTOCOL),
					Description:  "App 协议 的 DDoS 策略 case. 有效值：`tcp`，`udp`，`icmp` 和 `all`。",
				},
				Required:    true,
				Description: "App 协议 集合 的 DDoS 策略 case。",
			},
			"app_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue(DAYU_APP_TYPE), //to see the max
				Description:  "App 类型 DDoS 策略 case. 有效值：`WEB`，`GAME`，`APP` 和 `OTHER`。",
			},
			"tcp_start_port": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: tccommon.ValidatePort,
				Description:  "Start 端口 的 TCP 服务. 有效 值 ranges: (0~65535)。",
			},
			"tcp_end_port": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: tccommon.ValidatePort,
				Description:  "End 端口 的 TCP 服务. 有效 值 ranges: (0~65535). It 必须 是 greater 比 `tcp_start_port`。",
			},
			"udp_start_port": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: tccommon.ValidatePort,
				Description:  "Start 端口 的 UDP 服务. 有效 值 ranges: (0~65535)。",
			},
			"udp_end_port": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: tccommon.ValidatePort,
				Description:  "End 端口 的 UDP 服务. 有效 值 ranges: (0~65535). It 必须 是 greater 比 `udp_start_port`。",
			},
			"has_abroad": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue(DAYU_BOOL_FLAG),
				Description:  "Indicate 是否service involves overseas 或 不. 有效值：`无` 和 `yes`。",
			},
			"has_initiate_tcp": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue(DAYU_BOOL_FLAG),
				Description:  "Indicate 是否service actively initiates TCP requests 或 不. 有效值：`无` 和 `yes`。",
			},
			"has_initiate_udp": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue(DAYU_BOOL_FLAG),
				Description:  "Indicate 是否actively initiate UDP requests 或 不. 有效值：`无` 和 `yes`。",
			},
			"has_vpn": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue(DAYU_BOOL_FLAG),
				Description:  "Indicate 是否service involves VPN 服务 或 不. 有效值：`无` 和 `yes`。",
			},
			"peer_tcp_port": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: tccommon.ValidatePort,
				Description:  "端口 该 actively initiates TCP requests. 有效 值 ranges: (1~65535)。",
			},
			"peer_udp_port": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: tccommon.ValidatePort,
				Description:  "端口 该 actively initiates UDP requests. 有效 值 ranges: (1~65535)。",
			},
			"tcp_footprint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: tccommon.ValidateStringLengthInRange(1, 512),
				Description:  "fixed 签名 的 TCP 协议 load，有效 值 长度 是 范围 从 1 到 512。",
			},
			"udp_footprint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: tccommon.ValidateStringLengthInRange(1, 512),
				Description:  "fixed 签名 的 TCP 协议 load，有效 值 长度 是 范围 从 1 到 512。",
			},
			"web_api_urls": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type:        schema.TypeString,
					Description: "Web API URL",
				},
				Required:    true,
				Description: "Web API URL 集合。",
			},
			"min_tcp_package_len": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: tccommon.ValidateStringLengthInRange(1, 1499),
				Description:  "最小长度TCP 消息 包，有效 值 长度 should 是 greater 比 0 和 less 比 1500。",
			},
			"max_tcp_package_len": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: tccommon.ValidateStringLengthInRange(1, 1499),
				Description:  "max 长度 的 TCP 消息 包，有效 值 长度 should 是 greater 比 0 和 less 比 1500. It should 是 greater 比 `min_tcp_package_len`。",
			},
			"min_udp_package_len": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: tccommon.ValidateStringLengthInRange(1, 1499),
				Description:  "最小长度UDP 消息 包，有效 值 长度 should 是 greater 比 0 和 less 比 1500。",
			},
			"max_udp_package_len": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: tccommon.ValidateStringLengthInRange(1, 1499),
				Description:  "max 长度 的 UDP 消息 包，有效 值 长度 should 是 greater 比 0 和 less 比 1500. It should 是 greater 比 `min_udp_package_len`。",
			},
			//computed
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "创建时间 的 DDoS 策略 case。",
			},
			"scene_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID DDoS 策略 case。",
			},
		},
	}
}

func resourceTencentCloudDayuDdosPolicyCaseCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dayu_ddos_policy_case.create")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	request := dayu.NewCreateDDoSPolicyCaseRequest()
	resourceType := d.Get("resource_type").(string)
	request.Business = &resourceType
	request.CaseName = helper.String(d.Get("name").(string))
	platforms := d.Get("platform_types").(*schema.Set).List()
	for _, plat := range platforms {
		request.PlatformTypes = append(request.PlatformTypes, helper.String(plat.(string)))
	}
	protocols := d.Get("app_protocols").(*schema.Set).List()
	for _, protocol := range protocols {
		request.AppProtocols = append(request.AppProtocols, helper.String(protocol.(string)))
	}
	urls := d.Get("web_api_urls").(*schema.Set).List()
	for _, url := range urls {
		request.WebApiUrl = append(request.WebApiUrl, helper.String(url.(string)))
	}
	request.AppType = helper.String(d.Get("app_type").(string))
	request.HasAbroad = helper.String(d.Get("has_abroad").(string))
	request.HasInitiateTcp = helper.String(d.Get("has_initiate_tcp").(string))
	request.HasInitiateUdp = helper.String(d.Get("has_initiate_udp").(string))
	request.HasVPN = helper.String(d.Get("has_vpn").(string))
	request.PeerTcpPort = helper.String(d.Get("peer_tcp_port").(string))
	request.PeerUdpPort = helper.String(d.Get("peer_udp_port").(string))
	request.TcpFootprint = helper.String(d.Get("tcp_footprint").(string))
	request.UdpFootprint = helper.String(d.Get("udp_footprint").(string))

	tcpPortStart := d.Get("tcp_start_port").(string)
	tcpPortEnd := d.Get("tcp_end_port").(string)
	startInt, sErr := strconv.Atoi(tcpPortStart)
	endInt, eErr := strconv.Atoi(tcpPortEnd)
	if sErr != nil {
		return sErr
	}
	if eErr != nil {
		return eErr
	}
	if endInt < startInt {
		return fmt.Errorf("`tcp_start_port`:%s should not be greater than `tcp_end_port`:%s.", tcpPortStart, tcpPortEnd)
	}
	udpPortStart := d.Get("udp_start_port").(string)
	udpPortEnd := d.Get("udp_end_port").(string)
	startInt, sErr = strconv.Atoi(udpPortStart)
	endInt, eErr = strconv.Atoi(udpPortEnd)
	if sErr != nil {
		return sErr
	}
	if eErr != nil {
		return eErr
	}
	if endInt < startInt {
		return fmt.Errorf("`udp_start_port`:%s should not be greater than `udp_end_port`:%s.", udpPortStart, udpPortEnd)
	}
	request.TcpSportStart = &tcpPortStart
	request.TcpSportEnd = &tcpPortEnd
	request.UdpSportStart = &udpPortStart
	request.UdpSportEnd = &udpPortEnd

	minTcpPackageLen := d.Get("min_tcp_package_len").(string)
	maxTcpPackageLen := d.Get("max_tcp_package_len").(string)
	minTcpPackageLenInt, _ := strconv.Atoi(minTcpPackageLen)
	maxTcpPackageLenInt, _ := strconv.Atoi(maxTcpPackageLen)
	if maxTcpPackageLenInt < minTcpPackageLenInt {
		return fmt.Errorf("`min_tcp_package_len`:%s should not be greater than `max_tcp_package_len`:%s.", minTcpPackageLen, maxTcpPackageLen)
	}
	minUdpPackageLen := d.Get("min_udp_package_len").(string)
	maxUdpPackageLen := d.Get("max_udp_package_len").(string)
	minUdpPackageLenInt, _ := strconv.Atoi(minUdpPackageLen)
	maxUdpPackageLenInt, _ := strconv.Atoi(maxUdpPackageLen)
	if maxUdpPackageLenInt < minUdpPackageLenInt {
		return fmt.Errorf("`min_udp_package_len`:%s should not be greater than `max_udp_package_len`:%s.", minUdpPackageLen, maxUdpPackageLen)
	}
	request.MinTcpPackageLen = &minTcpPackageLen
	request.MaxTcpPackageLen = &maxTcpPackageLen
	request.MinUdpPackageLen = &minUdpPackageLen
	request.MaxUdpPackageLen = &maxTcpPackageLen

	dayuService := DayuService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	sceneId := ""
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := dayuService.CreateDdosPolicyCase(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		}
		sceneId = result
		return nil
	})

	if err != nil {
		return err
	}

	d.SetId(resourceType + tccommon.FILED_SP + sceneId)

	return resourceTencentCloudDayuDdosPolicyCaseRead(d, meta)
}

func resourceTencentCloudDayuDdosPolicyCaseRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dayu_ddos_policy_case.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	items := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(items) < 2 {
		return fmt.Errorf("broken ID of DDoS policy case")
	}
	resourceType := items[0]
	sceneId := items[1]

	dayuService := DayuService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	ddosPolicyCase, has, err := dayuService.DescribeDdosPolicyCase(ctx, resourceType, sceneId)
	if err != nil {
		err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			ddosPolicyCase, has, err = dayuService.DescribeDdosPolicyCase(ctx, resourceType, sceneId)
			if err != nil {
				return tccommon.RetryError(err)
			}
			return nil
		})
	}
	if err != nil {
		return err
	}
	if !has {
		d.SetId("")
		return nil
	}

	for _, record := range ddosPolicyCase.Record {
		key := *record.Key
		if key == "CaseName" {
			_ = d.Set("name", record.Value)
		}
		if key == "HasInitiateTcp" {
			_ = d.Set("has_initiate_tcp", record.Value)
		}
		if key == "HasInitiateUdp" {
			_ = d.Set("has_initiate_udp", record.Value)
		}
		if key == "HasVPN" {
			_ = d.Set("has_vpn", record.Value)
		}
		if key == "PeerTcpPort" {
			_ = d.Set("peer_tcp_port", record.Value)
		}
		if key == "PeerUdpPort" {
			_ = d.Set("peer_udp_port", record.Value)
		}
		if key == "TcpFootprint" {
			_ = d.Set("tcp_footprint", record.Value)
		}
		if key == "UdpFootprint" {
			_ = d.Set("udp_footprint", record.Value)
		}
		if key == "HasAbroad" {
			_ = d.Set("has_abroad", record.Value)
		}
		if key == "TcpSportStart" {
			_ = d.Set("tcp_start_port", record.Value)
		}
		if key == "TcpSportEnd" {
			_ = d.Set("tcp_end_port", record.Value)
		}
		if key == "UdpSportStart" {
			_ = d.Set("udp_start_port", record.Value)
		}
		if key == "UdpSportEnd" {
			_ = d.Set("udp_end_port", record.Value)
		}
		if key == "MaxUdpPackageLen" {
			_ = d.Set("max_udp_package_len", record.Value)
		}
		if key == "MinUdpPackageLen" {
			_ = d.Set("min_udp_package_len", record.Value)
		}
		if key == "MaxTcpPackageLen" {
			_ = d.Set("max_tcp_package_len", record.Value)
		}
		if key == "MinTcpPackageLen" {
			_ = d.Set("min_tcp_package_len", record.Value)
		}
		if key == "AppType" {
			_ = d.Set("app_type", record.Value)
		}
		if key == "AppProtocols" {
			_ = d.Set("app_protocols", strings.Split(*record.Value, ";"))
		}
		if key == "WebApiUrl" {
			_ = d.Set("web_api_urls", strings.Split(*record.Value, ";"))
		}
		if key == "PlatformTypes" {
			_ = d.Set("platform_types", strings.Split(*record.Value, ";"))
		}
		if key == "Id" {
			_ = d.Set("scene_id", record.Value)
		}
		if key == "CreateTime" {
			_ = d.Set("create_time", record.Value)
		}
	}
	return nil
}

func resourceTencentCloudDayuDdosPolicyCaseUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dayu_ddos_policy_case.update")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	items := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(items) < 2 {
		return fmt.Errorf("broken ID of DDoS policy case")
	}
	resourceType := items[0]
	sceneId := items[1]

	request := dayu.NewModifyDDoSPolicyCaseRequest()
	request.Business = &resourceType
	request.SceneId = &sceneId

	platforms := d.Get("platform_types").(*schema.Set).List()
	for _, plat := range platforms {
		request.PlatformTypes = append(request.PlatformTypes, helper.String(plat.(string)))
	}
	protocols := d.Get("app_protocols").(*schema.Set).List()
	for _, protocol := range protocols {
		request.AppProtocols = append(request.AppProtocols, helper.String(protocol.(string)))
	}
	urls := d.Get("web_api_urls").(*schema.Set).List()
	for _, url := range urls {
		request.WebApiUrl = append(request.WebApiUrl, helper.String(url.(string)))
	}
	request.AppType = helper.String(d.Get("app_type").(string))
	request.HasAbroad = helper.String(d.Get("has_abroad").(string))
	request.HasInitiateTcp = helper.String(d.Get("has_initiate_tcp").(string))
	request.HasInitiateUdp = helper.String(d.Get("has_initiate_udp").(string))
	request.HasVPN = helper.String(d.Get("has_vpn").(string))
	request.PeerTcpPort = helper.String(d.Get("peer_tcp_port").(string))
	request.PeerUdpPort = helper.String(d.Get("peer_udp_port").(string))
	request.TcpFootprint = helper.String(d.Get("tcp_footprint").(string))
	request.UdpFootprint = helper.String(d.Get("udp_footprint").(string))

	tcpPortStart := d.Get("tcp_start_port").(string)
	tcpPortEnd := d.Get("tcp_end_port").(string)
	startInt, sErr := strconv.Atoi(tcpPortStart)
	endInt, eErr := strconv.Atoi(tcpPortEnd)
	if sErr != nil {
		return sErr
	}
	if eErr != nil {
		return eErr
	}

	if endInt < startInt {
		return fmt.Errorf("`tcp_start_port`:%s should not be greater than `tcp_end_port`:%s.", tcpPortStart, tcpPortEnd)
	}
	udpPortStart := d.Get("udp_start_port").(string)
	udpPortEnd := d.Get("udp_end_port").(string)
	startInt, sErr = strconv.Atoi(udpPortStart)
	endInt, eErr = strconv.Atoi(udpPortEnd)
	if sErr != nil {
		return sErr
	}
	if eErr != nil {
		return eErr
	}

	if endInt < startInt {
		return fmt.Errorf("`udp_start_port`:%s should not be greater than `udp_end_port`:%s.", udpPortStart, udpPortEnd)
	}
	request.TcpSportStart = &tcpPortStart
	request.TcpSportEnd = &tcpPortEnd
	request.UdpSportStart = &udpPortStart
	request.UdpSportEnd = &udpPortEnd

	minTcpPackageLen := d.Get("min_tcp_package_len").(string)
	maxTcpPackageLen := d.Get("max_tcp_package_len").(string)
	minTcpPackageLenInt, _ := strconv.Atoi(minTcpPackageLen)
	maxTcpPackageLenInt, _ := strconv.Atoi(maxTcpPackageLen)
	if maxTcpPackageLenInt < minTcpPackageLenInt {
		return fmt.Errorf("`min_tcp_package_len`:%s should not be greater than `max_tcp_package_len`:%s.", minTcpPackageLen, maxTcpPackageLen)
	}
	minUdpPackageLen := d.Get("min_udp_package_len").(string)
	maxUdpPackageLen := d.Get("max_udp_package_len").(string)
	minUdpPackageLenInt, _ := strconv.Atoi(minUdpPackageLen)
	maxUdpPackageLenInt, _ := strconv.Atoi(maxUdpPackageLen)
	if maxUdpPackageLenInt < minUdpPackageLenInt {
		return fmt.Errorf("`min_udp_package_len`:%s should not be greater than `max_udp_package_len`:%s.", minUdpPackageLen, maxUdpPackageLen)
	}
	request.MinTcpPackageLen = &minTcpPackageLen
	request.MaxTcpPackageLen = &maxTcpPackageLen
	request.MinUdpPackageLen = &minUdpPackageLen
	request.MaxUdpPackageLen = &maxTcpPackageLen

	dayuService := DayuService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		e := dayuService.ModifyDdosPolicyCase(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		}
		return nil
	})
	if err != nil {
		return err
	}

	return resourceTencentCloudDayuDdosPolicyCaseRead(d, meta)
}

func resourceTencentCloudDayuDdosPolicyCaseDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dayu_ddos_policy_case.delete")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	items := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(items) < 2 {
		return fmt.Errorf("broken ID of DDoS policy")
	}
	resourceType := items[0]
	sceneId := items[1]

	dayuService := DayuService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	err := dayuService.DeleteDdosPolicyCase(ctx, resourceType, sceneId)

	if err != nil {
		if sdkErr, ok := err.(*sdkError.TencentCloudSDKError); ok {
			if sdkErr.Code == "ResourceInUse" {
				//if bind automatically, try to unbind policy first
				//get the automatically generated policy
				policies, dErr := dayuService.DescribeDdosPolicies(ctx, resourceType, "")
				if dErr != nil {
					err = dErr
					return err
				}
				bindPolicyId := ""
				bindResourceIds := []string{}
				for _, policy := range policies {
					if *policy.SceneId == sceneId {
						bindPolicyId = *policy.PolicyId
						for _, resources := range (*policy).BoundResources {
							bindResourceIds = append(bindResourceIds, *resources)
						}
					}
				}
				if bindPolicyId == "" || len(bindResourceIds) == 0 {
					return fmt.Errorf("the automatically generated policy of policy case %s can not be find", sceneId)
				}
				//unbind policy and resource
				for _, resourceId := range bindResourceIds {
					bErr := dayuService.UnbindDdosPolicy(ctx, resourceId, resourceType, bindPolicyId)
					if bErr != nil {
						err = bErr
						return err
					}
				}
			}
		}
		err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			e := dayuService.DeleteDdosPolicyCase(ctx, resourceType, sceneId)
			if e != nil {
				return tccommon.RetryError(e)
			}
			return nil
		})
	}

	if err != nil {
		return err
	}

	_, has, err := dayuService.DescribeDdosPolicyCase(ctx, resourceType, sceneId)
	if err != nil || has {
		err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			_, has, err = dayuService.DescribeDdosPolicyCase(ctx, resourceType, sceneId)
			if err != nil {
				return tccommon.RetryError(err)
			}

			if has {
				err = fmt.Errorf("delete DDoS policy case fail, DDoS policy case still exist.")
				return resource.RetryableError(err)
			}

			return nil
		})
	}
	if err != nil {
		return err
	}
	if !has {
		return nil
	} else {
		return errors.New("delete DDoS policy case fail, DDoS policy case still exist.")
	}
}
