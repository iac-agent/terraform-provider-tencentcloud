package vpc

import (
	"context"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudVpcNetDetect() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudVpcNetDetectCreate,
		Read:   resourceTencentCloudVpcNetDetectRead,
		Update: resourceTencentCloudVpcNetDetectUpdate,
		Delete: resourceTencentCloudVpcNetDetectDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"vpc_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "`VPC` 实例 `ID`. Such 作为:`vpc-12345678`。",
			},

			"subnet_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "子网实例 ID Such 作为:子网-12345678。",
			},

			"net_detect_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Network probe 名称， 最大 长度 不能 exceed 60 bytes。",
			},

			"detect_destination_ip": {
				Required: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "An 数组 probe destination IPv4 addresses. Up 到 two。",
			},

			"next_hop_type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "next hop 类型，currently we support following types: `VPN`: VPN 网关; `DIRECTCONNECT`: 私有 line 网关; `PEERCONNECTION`: peer 连接; `NAT`: NAT 网关; `NORMAL_CVM`: normal 云 服务器; `CCN`: 云 networking 网关; `NONEXTHOP`: 无 next hop。",
			},

			"next_hop_destination": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "destination 网关 的 next hop， 值 是 related 到 next hop 类型 如果 next hop 类型 是 VPN，和 值 是 VPN 网关 ID，such 作为: vpngw-12345678; 如果 next hop 类型 是 DIRECTCONNECT，和 值 是 私有 line 网关 ID，such 作为: dcg-12345678; 如果 next hop 类型 是 PEERCONNECTION，其中 takes 值 的 peer 连接 ID，such 作为: pcx-12345678; 如果 next hop 类型 是 NAT，和 值 是 Nat 网关，such 作为: nat-12345678; 如果 next hop 类型 是 NORMAL_CVM，其中 takes IPv4 地址 的 云 服务器，such 作为: 10.0.0.12; 如果 next hop 类型 是 CCN，和 值 是 云 网络 ID，such 作为: ccn-12345678; 如果 next hop 类型 是 NONEXTHOP，和 指定 网络 probe 是 网络 probe without next hop。",
			},

			"net_detect_description": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Network probe 描述",
			},
		},
	}
}

func resourceTencentCloudVpcNetDetectCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vpc_net_detect.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request     = vpc.NewCreateNetDetectRequest()
		response    = vpc.NewCreateNetDetectResponse()
		netDetectId string
	)
	if v, ok := d.GetOk("vpc_id"); ok {
		request.VpcId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("subnet_id"); ok {
		request.SubnetId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("net_detect_name"); ok {
		request.NetDetectName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("detect_destination_ip"); ok {
		detectDestinationIpSet := v.(*schema.Set).List()
		for i := range detectDestinationIpSet {
			detectDestinationIp := detectDestinationIpSet[i].(string)
			request.DetectDestinationIp = append(request.DetectDestinationIp, &detectDestinationIp)
		}
	}

	if v, ok := d.GetOk("next_hop_type"); ok {
		request.NextHopType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("next_hop_destination"); ok {
		request.NextHopDestination = helper.String(v.(string))
	}

	if v, ok := d.GetOk("net_detect_description"); ok {
		request.NetDetectDescription = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().CreateNetDetect(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create vpc netDetect failed, reason:%+v", logId, err)
		return err
	}

	netDetectId = *response.Response.NetDetect.NetDetectId
	d.SetId(netDetectId)

	return resourceTencentCloudVpcNetDetectRead(d, meta)
}

func resourceTencentCloudVpcNetDetectRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vpc_net_detect.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := VpcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	netDetectId := d.Id()

	netDetect, err := service.DescribeVpcNetDetectById(ctx, netDetectId)
	if err != nil {
		return err
	}

	if netDetect == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `VpcNetDetect` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if netDetect.VpcId != nil {
		_ = d.Set("vpc_id", netDetect.VpcId)
	}

	if netDetect.SubnetId != nil {
		_ = d.Set("subnet_id", netDetect.SubnetId)
	}

	if netDetect.NetDetectName != nil {
		_ = d.Set("net_detect_name", netDetect.NetDetectName)
	}

	if netDetect.DetectDestinationIp != nil {
		_ = d.Set("detect_destination_ip", netDetect.DetectDestinationIp)
	}

	if netDetect.NextHopType != nil {
		_ = d.Set("next_hop_type", netDetect.NextHopType)
	}

	if netDetect.NextHopDestination != nil {
		_ = d.Set("next_hop_destination", netDetect.NextHopDestination)
	}

	if netDetect.NetDetectDescription != nil {
		_ = d.Set("net_detect_description", netDetect.NetDetectDescription)
	}

	return nil
}

func resourceTencentCloudVpcNetDetectUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vpc_net_detect.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	request := vpc.NewModifyNetDetectRequest()

	netDetectId := d.Id()

	request.NetDetectId = &netDetectId

	mutableArgs := []string{
		"net_detect_name", "detect_destination_ip", "next_hop_type",
		"next_hop_destination", "net_detect_description",
	}
	needChange := false
	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}

	if needChange {

		if v, ok := d.GetOk("net_detect_name"); ok {
			request.NetDetectName = helper.String(v.(string))
		}

		if v, ok := d.GetOk("detect_destination_ip"); ok {
			detectDestinationIpSet := v.(*schema.Set).List()
			for i := range detectDestinationIpSet {
				detectDestinationIp := detectDestinationIpSet[i].(string)
				request.DetectDestinationIp = append(request.DetectDestinationIp, &detectDestinationIp)
			}
		}

		if v, ok := d.GetOk("next_hop_type"); ok {
			request.NextHopType = helper.String(v.(string))
		}

		if v, ok := d.GetOk("next_hop_destination"); ok {
			request.NextHopDestination = helper.String(v.(string))
		}

		if v, ok := d.GetOk("net_detect_description"); ok {
			request.NetDetectDescription = helper.String(v.(string))
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().ModifyNetDetect(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update vpc netDetect failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudVpcNetDetectRead(d, meta)
}

func resourceTencentCloudVpcNetDetectDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vpc_net_detect.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := VpcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	netDetectId := d.Id()

	if err := service.DeleteVpcNetDetectById(ctx, netDetectId); err != nil {
		return err
	}

	return nil
}
