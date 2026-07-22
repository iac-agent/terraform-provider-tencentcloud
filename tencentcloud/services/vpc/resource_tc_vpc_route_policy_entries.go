package vpc

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vpcv20170312 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudVpcRoutePolicyEntries() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudVpcRoutePolicyEntriesCreate,
		Read:   resourceTencentCloudVpcRoutePolicyEntriesRead,
		Update: resourceTencentCloudVpcRoutePolicyEntriesUpdate,
		Delete: resourceTencentCloudVpcRoutePolicyEntriesDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"route_policy_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "指定instance ID 路由 reception 策略。",
			},

			"route_policy_entry_set": {
				Type:        schema.TypeSet,
				Required:    true,
				Description: "Route reception 策略 entry 列表。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"route_policy_entry_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "指定unique ID IPv4 routing strategy entry.\n注意：此字段可能返回 null，表示未找到有效值。",
						},
						"cidr_block": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Destination ip 范围.\n注意：此字段可能返回 null，表示未找到有效值。",
						},
						"description": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Describes routing strategy 规则.\n注意：此字段可能返回 null，表示未找到有效值。",
						},
						"route_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Routing 类型\n\n指定USER-customized 数据 类型\nNETD: 指定route 对于 网络 detection.\nCCN: CCN 路由.\n注意：此字段可能返回 null，表示未找到有效值。",
						},
						"gateway_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Next hop 类型 types currently 支持:.\nCVM: 云 virtual machine 使用 公有 网络 网关 类型\nVPN: vpn 网关.\nDIRECTCONNECT: direct connect 网关.\nPEERCONNECTION: peering 连接.\nHAVIP: high availability virtual ip.\nNAT: 指定nat 网关. \nEIP: 指定public ip 地址 的 云 virtual machine.\nLOCAL_GATEWAY: 指定local 网关.\nPVGW: pvgw 网关.\n注意：此字段可能返回 null，表示未找到有效值。",
						},
						"gateway_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Gateway 唯一 ID.\n注意：此字段可能返回 null，表示未找到有效值。",
						},
						"priority": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "优先级 smaller 值 表示a higher 优先级\n注意：此字段可能返回 null，表示未找到有效值。",
						},
						"action": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "操作\nDROP: drop.\nDISABLE: receive 和 disable.\nACCEPT: receive 和 启用.\n注意：此字段可能返回 null，表示未找到有效值。",
						},
						"created_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "创建时间.\n\n注意：此字段可能返回 null，表示未找到有效值。",
						},
						"region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "指定region.\n注意：此字段可能返回 null，表示未找到有效值。",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudVpcRoutePolicyEntriesCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vpc_route_policy_entries.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId         = tccommon.GetLogId(tccommon.ContextNil)
		ctx           = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request       = vpcv20170312.NewCreateRoutePolicyEntriesRequest()
		routePolicyId string
	)

	if v, ok := d.GetOk("route_policy_id"); ok {
		request.RoutePolicyId = helper.String(v.(string))
		routePolicyId = v.(string)
	}

	if v, ok := d.GetOk("route_policy_entry_set"); ok {
		for _, item := range v.(*schema.Set).List() {
			routePolicyEntrySetMap := item.(map[string]interface{})
			routePolicyEntry := vpcv20170312.RoutePolicyEntry{}
			if v, ok := routePolicyEntrySetMap["cidr_block"].(string); ok && v != "" {
				routePolicyEntry.CidrBlock = helper.String(v)
			}

			if v, ok := routePolicyEntrySetMap["description"].(string); ok && v != "" {
				routePolicyEntry.Description = helper.String(v)
			}

			if v, ok := routePolicyEntrySetMap["route_type"].(string); ok && v != "" {
				routePolicyEntry.RouteType = helper.String(v)
			}

			if v, ok := routePolicyEntrySetMap["gateway_type"].(string); ok && v != "" {
				routePolicyEntry.GatewayType = helper.String(v)
			}

			if v, ok := routePolicyEntrySetMap["gateway_id"].(string); ok && v != "" {
				routePolicyEntry.GatewayId = helper.String(v)
			}

			if v, ok := routePolicyEntrySetMap["priority"].(int); ok {
				routePolicyEntry.Priority = helper.IntUint64(v)
			}

			if v, ok := routePolicyEntrySetMap["action"].(string); ok && v != "" {
				routePolicyEntry.Action = helper.String(v)
			}

			request.RoutePolicyEntrySet = append(request.RoutePolicyEntrySet, &routePolicyEntry)
		}
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().CreateRoutePolicyEntriesWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create vpc route policy entries failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	d.SetId(routePolicyId)
	return resourceTencentCloudVpcRoutePolicyEntriesRead(d, meta)
}

func resourceTencentCloudVpcRoutePolicyEntriesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vpc_route_policy_entries.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId         = tccommon.GetLogId(tccommon.ContextNil)
		ctx           = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service       = VpcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		routePolicyId = d.Id()
	)

	respData, err := service.DescribeVpcRoutePolicyEntriesById(ctx, routePolicyId)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[WARN]%s resource `tencentcloud_vpc_route_policy_entries` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	_ = d.Set("route_policy_id", routePolicyId)

	if len(respData) > 0 {
		routePolicyEntrySet := make([]map[string]interface{}, 0, len(respData))
		for _, item := range respData {
			entry := make(map[string]interface{})
			if item.RoutePolicyEntryId != nil {
				entry["route_policy_entry_id"] = *item.RoutePolicyEntryId
			}

			if item.CidrBlock != nil {
				entry["cidr_block"] = *item.CidrBlock
			}

			if item.Description != nil {
				entry["description"] = *item.Description
			}

			if item.RouteType != nil {
				entry["route_type"] = *item.RouteType
			}

			if item.GatewayType != nil {
				entry["gateway_type"] = *item.GatewayType
			}

			if item.GatewayId != nil {
				entry["gateway_id"] = *item.GatewayId
			}

			if item.Priority != nil {
				entry["priority"] = *item.Priority
			}

			if item.Action != nil {
				entry["action"] = *item.Action
			}

			if item.CreatedTime != nil {
				entry["created_time"] = *item.CreatedTime
			}

			if item.Region != nil {
				entry["region"] = *item.Region
			}

			routePolicyEntrySet = append(routePolicyEntrySet, entry)
		}

		_ = d.Set("route_policy_entry_set", routePolicyEntrySet)
	}

	return nil
}

func resourceTencentCloudVpcRoutePolicyEntriesUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vpc_route_policy_entries.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId         = tccommon.GetLogId(tccommon.ContextNil)
		ctx           = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service       = VpcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		request       = vpcv20170312.NewResetRoutePolicyEntriesRequest()
		routePolicyId = d.Id()
	)

	// temp get RoutePolicyDescription and RoutePolicyName
	respData, err := service.DescribeVpcRoutePolicyById(ctx, routePolicyId)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[WARN]%s resource `tencentcloud_vpc_route_policy_entries` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return fmt.Errorf("resource `tencentcloud_vpc_route_policy_entries` [%s] not found, please check if it has been deleted.", d.Id())
	}

	if respData.RoutePolicyDescription != nil {
		request.RoutePolicyDescription = respData.RoutePolicyDescription
	}

	if respData.RoutePolicyName != nil {
		request.RoutePolicyName = respData.RoutePolicyName
	}

	if v, ok := d.GetOk("route_policy_entry_set"); ok {
		for _, item := range v.(*schema.Set).List() {
			routePolicyEntrySetMap := item.(map[string]interface{})
			routePolicyEntry := vpcv20170312.RoutePolicyEntry{}
			if v, ok := routePolicyEntrySetMap["route_policy_entry_id"].(string); ok && v != "" {
				routePolicyEntry.RoutePolicyEntryId = helper.String(v)
			}

			if v, ok := routePolicyEntrySetMap["cidr_block"].(string); ok && v != "" {
				routePolicyEntry.CidrBlock = helper.String(v)
			}

			if v, ok := routePolicyEntrySetMap["description"].(string); ok && v != "" {
				routePolicyEntry.Description = helper.String(v)
			}

			if v, ok := routePolicyEntrySetMap["route_type"].(string); ok && v != "" {
				routePolicyEntry.RouteType = helper.String(v)
			}

			if v, ok := routePolicyEntrySetMap["gateway_type"].(string); ok && v != "" {
				routePolicyEntry.GatewayType = helper.String(v)
			}

			if v, ok := routePolicyEntrySetMap["gateway_id"].(string); ok && v != "" {
				routePolicyEntry.GatewayId = helper.String(v)
			}

			if v, ok := routePolicyEntrySetMap["priority"].(int); ok {
				routePolicyEntry.Priority = helper.IntUint64(v)
			}

			if v, ok := routePolicyEntrySetMap["action"].(string); ok && v != "" {
				routePolicyEntry.Action = helper.String(v)
			}

			request.RoutePolicyEntrySet = append(request.RoutePolicyEntrySet, &routePolicyEntry)
		}
	}

	request.RoutePolicyId = &routePolicyId
	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().ResetRoutePolicyEntriesWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s update vpc route policy entries failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return resourceTencentCloudVpcRoutePolicyEntriesRead(d, meta)
}

func resourceTencentCloudVpcRoutePolicyEntriesDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vpc_route_policy_entries.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId         = tccommon.GetLogId(tccommon.ContextNil)
		ctx           = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service       = VpcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		request       = vpcv20170312.NewDeleteRoutePolicyEntriesRequest()
		routePolicyId = d.Id()
	)

	// get all entryIds first
	respData, err := service.DescribeVpcRoutePolicyEntriesById(ctx, routePolicyId)
	if err != nil {
		return err
	}

	if respData == nil || len(respData) == 0 {
		return nil
	}

	for _, item := range respData {
		routePolicyEntry := vpcv20170312.RoutePolicyEntry{}
		if item.RoutePolicyEntryId != nil {
			routePolicyEntry.RoutePolicyEntryId = item.RoutePolicyEntryId
		}

		request.RoutePolicyEntrySet = append(request.RoutePolicyEntrySet, &routePolicyEntry)
	}

	request.RoutePolicyId = &routePolicyId
	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().DeleteRoutePolicyEntriesWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete vpc route policy entries failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return nil
}
