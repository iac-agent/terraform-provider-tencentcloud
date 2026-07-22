package dayuv2

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svcantiddos "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/antiddos"
	svcdayu "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/dayu"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	antiddos "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/antiddos/v20200309"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudDayuDDosIpAttachmentV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDayuDDosIpAttachmentCreateV2,
		Read:   resourceTencentCloudDayuDDosIpAttachmentReadV2,
		Delete: resourceTencentCloudDayuDDosIpAttachmentDeleteV2,
		Schema: map[string]*schema.Schema{
			"bgp_instance_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Anti-DDoS 实例 ID。",
			},

			"bound_ip_list": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				Description: "数组 IPs 到 bind 到 Anti-DDoS 实例. For Anti-DDoS Pro Single IP 实例， 数组 包含only 一个 IP. 如果 there 是 无 IPs 到 bind，它 是 空; however，either BoundDevList 或 UnBoundDevList 必须 不 是 空。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "IP 地址",
						},
						"biz_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Category 的 product 该 可以 是 bound. 有效值：公有 (CVM 和 CLB)，bm (BM)，eni (ENI)，vpngw (VPN 网关)，natgw (NAT 网关)，waf (WAF)，fpc (financial products)，gaap (GAAP)，和 other (hosted IP). 此 字段 为必填项 当 您 perform binding。",
						},
						"instance_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Anti-DDoS 实例 ID IP. 此 字段 为必填项 仅 当 实例 是 bound 到 IP. For 示例，此 字段 实例 ID 将 是 eni-* 如果 实例 ID 是 bound 到 ENI IP; none 如果 there 是 无 实例 到 bind 到 managed IP。",
						},
						"device_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Sub-product category. 有效值：cvm (CVM)，lb (Load balancer)，eni (ENI)，vpngw (VPN 网关)，natgw (NAT 网关)，waf (WAF)，fpc (financial products)，gaap (GAAP)，eip (BM EIP) 和 other (managed IP). 此 字段 为必填项 当 您 perform binding。",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudDayuDDosIpAttachmentCreateV2(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dayu_ddos_ip_attachment_v2.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := svcantiddos.NewAntiddosService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())

	var bgpInstanceId string
	boundIps := make([]string, 0)
	request := antiddos.NewCreateBoundIPRequest()

	request.Business = helper.String(svcdayu.BUSINESS_BGP_MULTIP)
	if v, ok := d.GetOk("bgp_instance_id"); ok {
		bgpInstanceId = v.(string)
		request.Id = helper.String(bgpInstanceId)
	}

	if v, ok := d.GetOk("bound_ip_list"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			boundIpInfo := antiddos.BoundIpInfo{}
			if v, ok := dMap["ip"]; ok {
				boundIps = append(boundIps, v.(string))
				boundIpInfo.Ip = helper.String(v.(string))
			}
			if v, ok := dMap["biz_type"]; ok {
				boundIpInfo.BizType = helper.String(v.(string))
			}
			if v, ok := dMap["instance_id"]; ok {
				boundIpInfo.InstanceId = helper.String(v.(string))
			}
			if v, ok := dMap["device_type"]; ok {
				boundIpInfo.DeviceType = helper.String(v.(string))
			}
			boundIpInfo.IspCode = helper.IntUint64(svcdayu.ISP_CODE_BGP)
			request.BoundDevList = append(request.BoundDevList, &boundIpInfo)
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseAntiddosClient().CreateBoundIP(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		if *result.Response.Success.Code != svcdayu.RESPONSE_SUCCESS_CODE {
			return resource.RetryableError(errors.New("request failed"))
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create antiddos boundip failed, reason:%+v", logId, err)
		return err
	}
	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		boundip, e := service.DescribeAntiddosBoundipById(ctx, bgpInstanceId)
		if e != nil {
			return tccommon.RetryError(e)
		}
		if *boundip.BoundStatus == svcdayu.DAYU_BOUNDSTATUS_IDLE {
			return nil
		}
		return resource.RetryableError(errors.New("still building."))

	})
	if err != nil {
		return err
	}

	d.SetId(bgpInstanceId + tccommon.FILED_SP + strings.Join(boundIps, tccommon.COMMA_SP))

	return resourceTencentCloudDayuDDosIpAttachmentReadV2(d, meta)
}

func resourceTencentCloudDayuDDosIpAttachmentReadV2(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dayu_ddos_ip_attachment_v2.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := svcantiddos.NewAntiddosService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	bgpInstanceId := idSplit[0]
	boundIps := idSplit[1]
	boundIpMap := make(map[string]bool)
	for _, boundIp := range strings.Split(boundIps, tccommon.COMMA_SP) {
		boundIpMap[boundIp] = true
	}

	boundip, err := service.DescribeAntiddosBoundipById(ctx, bgpInstanceId)
	if err != nil {
		return err
	}

	if boundip == nil {
		d.SetId("")
		log.Printf("resource `AntiddosIp` %s does not exist", d.Id())
		return nil
	}
	_ = d.Set("bgp_instance_id", bgpInstanceId)
	boundIpList := make([]map[string]interface{}, 0)
	if boundip.EipProductInfos != nil {
		for _, item := range boundip.EipProductInfos {
			if _, ok := boundIpMap[*item.Ip]; !ok {
				continue
			}
			boundIpListItem := make(map[string]interface{})
			boundIpListItem["ip"] = *item.Ip
			boundIpListItem["biz_type"] = *item.BizType
			boundIpListItem["instance_id"] = *item.InstanceId
			boundIpListItem["device_type"] = *item.DeviceType
			boundIpList = append(boundIpList, boundIpListItem)
		}
	}
	_ = d.Set("bound_ip_list", boundIpList)

	return nil
}

func resourceTencentCloudDayuDDosIpAttachmentDeleteV2(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dayu_ddos_ip_attachment_v2.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := svcantiddos.NewAntiddosService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	bgpInstanceId := idSplit[0]
	boundIps := idSplit[1]

	request := antiddos.NewCreateBoundIPRequest()
	request.Business = helper.String(svcdayu.BUSINESS_BGP_MULTIP)
	request.Id = helper.String(bgpInstanceId)
	ubBoundDevList := make([]*antiddos.BoundIpInfo, 0)
	for _, boundIp := range strings.Split(boundIps, tccommon.COMMA_SP) {
		boundIpValue := boundIp
		ubBoundDevList = append(
			ubBoundDevList,
			&antiddos.BoundIpInfo{
				Ip: &boundIpValue,
			},
		)
	}
	request.UnBoundDevList = ubBoundDevList

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseAntiddosClient().CreateBoundIP(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		if *result.Response.Success.Code != svcdayu.RESPONSE_SUCCESS_CODE {
			return resource.RetryableError(errors.New("request failed"))
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create antiddos boundip failed, reason:%+v", logId, err)
		return err
	}

	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		boundip, e := service.DescribeAntiddosBoundipById(ctx, bgpInstanceId)
		if e != nil {
			return tccommon.RetryError(e)
		}
		if *boundip.BoundStatus == svcdayu.DAYU_BOUNDSTATUS_IDLE {
			return nil
		}
		return resource.RetryableError(errors.New("still building."))

	})
	if err != nil {
		return err
	}

	return nil
}
