package vpc

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vpcv20170312 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudVpcReplaceRoutesWithRoutePolicyConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudVpcReplaceRoutesWithRoutePolicyConfigCreate,
		Read:   resourceTencentCloudVpcReplaceRoutesWithRoutePolicyConfigRead,
		Update: resourceTencentCloudVpcReplaceRoutesWithRoutePolicyConfigUpdate,
		Delete: resourceTencentCloudVpcReplaceRoutesWithRoutePolicyConfigDelete,
		Schema: map[string]*schema.Schema{
			"route_table_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Route Table Instance ID.",
			},

			"routes": {
				Type:        schema.TypeSet,
				Required:    true,
				Description: "Routing policy object. requires specifying the unique ID of routing policy (RouteItemId).",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"route_item_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Route unique policy ID.",
						},
						"force_match_policy": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Match the route reception policy tag.",
						},
					},
				},
			},

			"need_router_info": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Indicates whether to obtain route policy information. Maps to the `NeedRouterInfo` parameter of the `DescribeRouteTables` API.",
			},

			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filter name, maps to `Filters.Name` of the `DescribeRouteTables` API. Must be used together with `values`.",
			},

			"values": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Filter values, maps to `Filters.Values` of the `DescribeRouteTables` API. Must be used together with `name`.",
			},

			"total_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total number of matching instances, returned by the `DescribeRouteTables` API.",
			},
		},
	}
}

func resourceTencentCloudVpcReplaceRoutesWithRoutePolicyConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vpc_replace_routes_with_route_policy_config.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		routeTableId string
	)
	if v, ok := d.GetOk("route_table_id"); ok {
		routeTableId = v.(string)
	}

	d.SetId(routeTableId)
	return resourceTencentCloudVpcReplaceRoutesWithRoutePolicyConfigUpdate(d, meta)
}

func resourceTencentCloudVpcReplaceRoutesWithRoutePolicyConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vpc_replace_routes_with_route_policy_config.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = VpcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	routeTableId := d.Id()

	var (
		needRouterInfo *bool
		name           string
		values         []*string
		totalCount     uint64
		respData       []VpcRouteTableBasicInfo
	)

	if v, ok := d.GetOkExists("need_router_info"); ok {
		needRouterInfo = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOk("name"); ok {
		name = v.(string)
	}

	if v, ok := d.GetOk("values"); ok {
		valuesList := v.([]interface{})
		if len(valuesList) > 0 {
			values = make([]*string, 0, len(valuesList))
			for _, item := range valuesList {
				value := item.(string)
				values = append(values, helper.String(value))
			}
		}
	}

	var errRet error
	e := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		respData, totalCount, errRet = service.DescribeRouteTables(ctx, routeTableId, "", "", map[string]string{}, nil, "", needRouterInfo, name, values)
		if errRet != nil {
			return tccommon.RetryError(errRet, tccommon.InternalError)
		}
		return nil
	})
	if e != nil {
		return e
	}

	if len(respData) == 0 {
		log.Printf("[CRUD]%s api[%s] vpc replace_routes_with_route_policy_config id=%s", logId, "DescribeRouteTables", d.Id())
		d.SetId("")
		return nil
	}

	_ = d.Set("route_table_id", routeTableId)
	_ = d.Set("total_count", totalCount)

	return nil
}

func resourceTencentCloudVpcReplaceRoutesWithRoutePolicyConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vpc_replace_routes_with_route_policy_config.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId        = tccommon.GetLogId(tccommon.ContextNil)
		ctx          = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request      = vpcv20170312.NewReplaceRoutesWithRoutePolicyRequest()
		routeTableId = d.Id()
	)

	if v, ok := d.GetOk("routes"); ok {
		for _, item := range v.(*schema.Set).List() {
			routesMap := item.(map[string]interface{})
			replaceRoutesWithRoutePolicyRoute := vpcv20170312.ReplaceRoutesWithRoutePolicyRoute{}
			if v, ok := routesMap["route_item_id"].(string); ok && v != "" {
				replaceRoutesWithRoutePolicyRoute.RouteItemId = helper.String(v)
			}

			if v, ok := routesMap["force_match_policy"].(bool); ok {
				replaceRoutesWithRoutePolicyRoute.ForceMatchPolicy = helper.Bool(v)
			}

			request.Routes = append(request.Routes, &replaceRoutesWithRoutePolicyRoute)
		}
	}

	request.RouteTableId = &routeTableId
	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().ReplaceRoutesWithRoutePolicyWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s update vpc replace routes with route policy config failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return resourceTencentCloudVpcReplaceRoutesWithRoutePolicyConfigRead(d, meta)
}

func resourceTencentCloudVpcReplaceRoutesWithRoutePolicyConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vpc_replace_routes_with_route_policy_config.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
