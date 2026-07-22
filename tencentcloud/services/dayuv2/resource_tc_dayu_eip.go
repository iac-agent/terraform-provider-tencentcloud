package dayuv2

import (
	"context"
	"fmt"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svcantiddos "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/antiddos"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceTencentCloudDayuEip() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDayuEipCreate,
		Read:   resourceTencentCloudDayuEipRead,
		Delete: resourceTencentCloudDayuEipDelete,

		Schema: map[string]*schema.Schema{
			"resource_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID 资源。",
			},
			"eip": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Eip 的 资源。",
			},
			"bind_resource_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Resource ID 到 bind。",
			},
			"bind_resource_region": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Resource 地域 到 bind。",
			},
			"bind_resource_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue(DDOS_EIP_BIND_RESOURCE_TYPE),
				Description:  "资源类型 到 bind，值 范围 [`clb`，`cvm`]。",
			},
			"resource_region": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "地域 的 资源 实例。",
			},
			"eip_bound_rsc_ins": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Eip bound rsc ins 的 资源 实例。",
			},
			"eip_bound_rsc_eni": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Eip bound rsc eni 的 资源 实例。",
			},
			"eip_bound_rsc_vip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Eip bound rsc VIP 的 资源 实例。",
			},
			"eip_address_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Eip 地址 状态 资源 实例。",
			},
			"protection_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Protection 状态 资源 实例。",
			},
			"created_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Created 时间 的 资源 实例。",
			},
			"expired_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "过期时间 的 资源 实例。",
			},
			"modify_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "修改时间 的 资源 实例。",
			},
		},
	}
}

func resourceTencentCloudDayuEipCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dayu_l4_rule.create")()
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	antiddosService := svcantiddos.NewAntiddosService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
	bindResourceType := d.Get("bind_resource_type").(string)
	resourceId := d.Get("resource_id").(string)
	eip := d.Get("eip").(string)
	bindResourceId := d.Get("bind_resource_id").(string)
	bindResourceRegion := d.Get("bind_resource_region").(string)
	if bindResourceType == DDOS_EIP_BIND_RESOURCE_TYPE_CLB {
		err := antiddosService.AssociateDDoSEipLoadBalancer(ctx, resourceId, eip, bindResourceId, bindResourceRegion)
		if err != nil {
			return err
		}
	}

	if bindResourceType == DDOS_EIP_BIND_RESOURCE_TYPE_CVM {
		err := antiddosService.AssociateDDoSEipAddress(ctx, resourceId, eip, bindResourceId, bindResourceRegion)
		if err != nil {
			return err
		}
	}

	for {
		bgpIPInstances, err := antiddosService.DescribeListBGPIPInstances(ctx, resourceId, DDOS_EIP_BIND_STATUS, 0, 10)
		if err != nil {
			return err
		}
		if len(bgpIPInstances) != 0 {
			break
		}
	}
	d.SetId(resourceId + tccommon.FILED_SP + eip)
	return resourceTencentCloudDayuEipRead(d, meta)
}

func resourceTencentCloudDayuEipRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dayu_l4_rule.read")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	items := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(items) < 2 {
		return fmt.Errorf("broken ID of dayu eip.")
	}
	resourceId := items[0]
	antiddosService := svcantiddos.NewAntiddosService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
	bgpIPInstances, err := antiddosService.DescribeListBGPIPInstances(ctx, resourceId, DDOS_EIP_BIND_STATUS, 0, 10)
	if err != nil {
		return err
	}
	if len(bgpIPInstances) != 0 {
		posBGPIPInstance := bgpIPInstances[0]
		_ = d.Set("resource_region", *posBGPIPInstance.Region.Region)
		_ = d.Set("eip_bound_rsc_ins", *posBGPIPInstance.EipAddressInfo.EipBoundRscIns)
		_ = d.Set("eip_bound_rsc_eni", *posBGPIPInstance.EipAddressInfo.EipBoundRscEni)
		_ = d.Set("eip_bound_rsc_vip", *posBGPIPInstance.EipAddressInfo.EipBoundRscVip)
		_ = d.Set("eip_address_status", *posBGPIPInstance.EipAddressStatus)
		_ = d.Set("protection_status", *posBGPIPInstance.Status)
		_ = d.Set("created_time", *posBGPIPInstance.CreatedTime)
		_ = d.Set("expired_time", *posBGPIPInstance.ExpiredTime)
		_ = d.Set("modify_time", *posBGPIPInstance.EipAddressInfo.ModifyTime)
	}

	return nil
}

func resourceTencentCloudDayuEipDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dayu_l4_rule.delete")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	items := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(items) < 2 {
		return fmt.Errorf("broken ID of dayu eip.")
	}
	resourceId := items[0]
	eip := items[1]
	antiddosService := svcantiddos.NewAntiddosService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
	err := antiddosService.DisassociateDDoSEipAddress(ctx, resourceId, eip)
	if err != nil {
		return err
	}
	return nil
}
