package dcg

import (
	"context"
	"fmt"
	"log"
	"strings"

	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	svctag "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceTencentCloudDcGatewayInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDcGatewayCreate,
		Read:   resourceTencentCloudDcGatewayRead,
		Update: resourceTencentCloudDcGatewayUpdate,
		Delete: resourceTencentCloudDcGatewayDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: tccommon.ValidateStringLengthInRange(1, 60),
				Description:  "名称 DCG。",
			},
			"network_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue(DCG_NETWORK_TYPES),
				Description:  "类型 associated network. Valid 值: `VPC` and `CCN`。",
			},
			"network_instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "If the `network_type` 值 is `VPC`，the available 值 is 私有网络 ID But when the `network_type` 值 is `CCN`，the available 值 is CCN 实例 ID",
			},
			"gateway_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      DCG_GATEWAY_TYPE_NORMAL,
				ValidateFunc: tccommon.ValidateAllowedStringValue(DCG_GATEWAY_TYPES),
				Description:  "类型 gateway. Valid 值: `NORMAL` and `NAT`. 默认为 `NORMAL`. NOTES: CCN only supports `NORMAL` and a VPC can create two DCGs，the one is NAT 类型 and the other is non-NAT 类型",
			},
			"mode_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "CCN route publishing method. 有效值：standard and exquisite. This parameter is only valid for the CCN direct connect gateway。",
			},
			"gateway_asn": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: "Dedicated connection gateway custom ASN，range: 45090，64512-65534 and 4200000000-4294967294。",
			},
			"zone": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: "Availability 可用区 where the direct connect gateway resides。",
			},
			"ha_zone_group_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "ID DC highly available placement group。",
			},
			"cnn_route_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "类型 CCN route. Valid 值: `BGP` and `STATIC`. The property is available when the DCG 类型 is CCN gateway and BGP 已启用",
			},
			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "标签键-值 pairs for the DC gateway. Multiple 标签 can be set。",
			},
			//compute
			"enable_bgp": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "表示是否BGP is 已启用",
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "创建时间 of resource。",
			},
		},
	}
}

func resourceTencentCloudDcGatewayCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dc_gateway.create")()

	var (
		logId             = tccommon.GetLogId(tccommon.ContextNil)
		request           = vpc.NewCreateDirectConnectGatewayRequest()
		response          = vpc.NewCreateDirectConnectGatewayResponse()
		networkType       string
		networkInstanceId string
		gatewayType       string
	)

	if v, ok := d.GetOk("name"); ok {
		request.DirectConnectGatewayName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("network_type"); ok {
		request.NetworkType = helper.String(v.(string))
		networkType = v.(string)
	}

	if v, ok := d.GetOk("network_instance_id"); ok {
		request.NetworkInstanceId = helper.String(v.(string))
		networkInstanceId = v.(string)
	} else {
		request.NetworkInstanceId = helper.String("")
	}

	if v, ok := d.GetOk("gateway_type"); ok {
		request.GatewayType = helper.String(v.(string))
		gatewayType = v.(string)
	}

	if v, ok := d.GetOk("mode_type"); ok {
		request.ModeType = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("gateway_asn"); ok {
		request.GatewayAsn = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("zone"); ok {
		request.Zone = helper.String(v.(string))
	}

	if v, ok := d.GetOk("ha_zone_group_id"); ok {
		request.HaZoneGroupId = helper.String(v.(string))
	}

	// Extract tags from schema
	var tags []*vpc.Tag
	if temp := helper.GetTags(d, "tags"); len(temp) > 0 {
		for k, v := range temp {
			tags = append(tags, &vpc.Tag{
				Key:   helper.String(k),
				Value: helper.String(v),
			})
		}
	}

	// Set tags in request if present
	if len(tags) > 0 {
		request.Tags = tags
	}

	if networkType == DCG_NETWORK_TYPE_VPC && !strings.HasPrefix(networkInstanceId, "vpc") {
		return fmt.Errorf("if `network_type` is '%s', the field `network_instance_id` must be a VPC resource", DCG_NETWORK_TYPE_VPC)
	}

	// if networkType == DCG_NETWORK_TYPE_CCN && !strings.HasPrefix(networkInstanceId, "ccn") {
	// 	return fmt.Errorf("if `network_type` is '%s', the field `network_instance_id` must be a CCN resource", DCG_NETWORK_TYPE_CCN)
	// }

	if networkType == DCG_NETWORK_TYPE_CCN && gatewayType != DCG_GATEWAY_TYPE_NORMAL {
		return fmt.Errorf("if `network_type` is '%s', the field `gateway_type` must be '%s'", DCG_NETWORK_TYPE_CCN, DCG_GATEWAY_TYPE_NORMAL)
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().CreateDirectConnectGateway(request)
		if e != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), e.Error())
			return tccommon.RetryError(e)
		}

		if result == nil || result.Response == nil || result.Response.DirectConnectGateway == nil {
			return resource.NonRetryableError(fmt.Errorf("Create direct connect gateway failed, Response is nil."))
		}

		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create direct connect gateway failed, reason:%s\n", logId, err.Error())
		return err
	}

	if response.Response.DirectConnectGateway.DirectConnectGatewayId == nil {
		return fmt.Errorf("DirectConnectGatewayId is nil.")
	}

	d.SetId(*response.Response.DirectConnectGateway.DirectConnectGatewayId)

	// set ccn route type
	if v, ok := d.GetOk("cnn_route_type"); ok {
		if v.(string) != "" && v.(string) != "STATIC" {
			request := vpc.NewModifyDirectConnectGatewayAttributeRequest()
			request.CcnRouteType = helper.String(v.(string))
			request.DirectConnectGatewayId = helper.String(d.Id())
			err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
				_, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().ModifyDirectConnectGatewayAttribute(request)
				if e != nil {
					log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), e.Error())
					return tccommon.RetryError(e)
				}

				return nil
			})

			if err != nil {
				log.Printf("[CRITAL]%s update direct connect gateway failed, reason:%s\n", logId, err.Error())
				return err
			}
		}
	}

	return resourceTencentCloudDcGatewayRead(d, meta)
}

func resourceTencentCloudDcGatewayRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dc_gateway.read")()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = VpcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		info, has, e := service.DescribeDirectConnectGateway(ctx, d.Id())
		if e != nil {
			return tccommon.RetryError(e)
		}

		if has == 0 {
			d.SetId("")
			return nil
		}

		_ = d.Set("name", info.name)
		_ = d.Set("network_type", info.networkType)
		_ = d.Set("network_instance_id", info.networkInstanceId)
		_ = d.Set("gateway_type", info.gatewayType)
		_ = d.Set("cnn_route_type", info.cnnRouteType)
		_ = d.Set("enable_bgp", info.enableBGP)
		_ = d.Set("create_time", info.createTime)
		if info.modeType != "" {
			_ = d.Set("mode_type", info.modeType)
		}
		if info.gatewayAsn != 0 {
			_ = d.Set("gateway_asn", info.gatewayAsn)
		}
		if info.zone != "" {
			_ = d.Set("zone", info.zone)
		}
		return nil
	})

	if err != nil {
		return err
	}

	// Retrieve tags using tag service
	tcClient := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
	tagService := svctag.NewTagService(tcClient)
	tags, err := tagService.DescribeResourceTags(ctx, "vpc", "dcg", tcClient.Region, d.Id())
	if err != nil {
		return err
	}
	_ = d.Set("tags", tags)

	return nil
}

func resourceTencentCloudDcGatewayUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dc_gateway.update")()

	var (
		logId = tccommon.GetLogId(tccommon.ContextNil)
		ctx   = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		dcgId = d.Id()
	)

	if d.HasChange("name") || d.HasChange("mode_type") || d.HasChange("cnn_route_type") {
		request := vpc.NewModifyDirectConnectGatewayAttributeRequest()
		if d.HasChange("name") {
			if v, ok := d.GetOk("name"); ok {
				request.DirectConnectGatewayName = helper.String(v.(string))
			}
		}

		if d.HasChange("mode_type") {
			if v, ok := d.GetOk("mode_type"); ok {
				request.ModeType = helper.String(v.(string))
			}
		}

		if d.HasChange("cnn_route_type") {
			if v, ok := d.GetOk("cnn_route_type"); ok {
				request.CcnRouteType = helper.String(v.(string))
			}
		}

		request.DirectConnectGatewayId = helper.String(dcgId)
		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			_, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().ModifyDirectConnectGatewayAttribute(request)
			if e != nil {
				log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), e.Error())
				return tccommon.RetryError(e)
			}

			return nil
		})

		if err != nil {
			log.Printf("[CRITAL]%s update direct connect gateway failed, reason:%s\n", logId, err.Error())
			return err
		}
	}

	// Handle tag changes
	if d.HasChange("tags") {
		oldValue, newValue := d.GetChange("tags")
		replaceTags, deleteTags := svctag.DiffTags(oldValue.(map[string]interface{}), newValue.(map[string]interface{}))
		tcClient := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
		tagService := svctag.NewTagService(tcClient)
		resourceName := tccommon.BuildTagResourceName("vpc", "dcg", tcClient.Region, d.Id())
		if err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags); err != nil {
			return err
		}
	}

	return resourceTencentCloudDcGatewayRead(d, meta)
}

func resourceTencentCloudDcGatewayDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dc_gateway.delete")()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		request = vpc.NewDeleteDirectConnectGatewayRequest()
		dcgId   = d.Id()
	)

	request.DirectConnectGatewayId = helper.String(dcgId)
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		_, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().DeleteDirectConnectGateway(request)
		if e != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), e.Error())
			return tccommon.RetryError(e)
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s delete direct connect gateway failed, reason:%s\n", logId, err.Error())
		return err
	}

	return nil
}
