package vpc

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
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

			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filter name of the `DescribeRouteTables` read request. Populates the `Name` of a `Filter` entry. Only takes effect when `route_table_ids` is not set.",
			},

			"values": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Filter values of the `DescribeRouteTables` read request. Populates the `Values` of the same `Filter` entry as `name`. Only takes effect when `route_table_ids` is not set.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"need_router_info": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Indicates whether to obtain route policy info. Maps to `NeedRouterInfo` of `DescribeRouteTables`.",
			},

			"route_table_ids": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Route table instance IDs, e.g. rtb-azd4dt1c. Maps to `RouteTableIds` of `DescribeRouteTables`. Mutually exclusive with `Filters` (including the `route-table-id` filter derived from `route_table_id`); when set, the read path queries by `RouteTableIds` only.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"limit": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(1, 100),
				Description:  "Return quantity of `DescribeRouteTables`, default is 20, max value is 100. When unset, the read helper uses its internal default of 100.",
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
		logId        = tccommon.GetLogId(tccommon.ContextNil)
		ctx          = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service      = VpcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		routeTableId = d.Id()
	)

	var (
		filterName     string
		filterValues   []*string
		needRouterInfo *bool
		routeTableIds  []*string
		limit          int
	)
	if v, ok := d.GetOk("name"); ok {
		filterName = v.(string)
	}
	if v, ok := d.GetOk("values"); ok {
		for _, item := range v.([]interface{}) {
			val := item.(string)
			filterValues = append(filterValues, helper.String(val))
		}
	}
	if v, ok := d.GetOkExists("need_router_info"); ok {
		needRouterInfo = helper.Bool(v.(bool))
	}
	if v, ok := d.GetOk("route_table_ids"); ok {
		for _, item := range v.([]interface{}) {
			val := item.(string)
			routeTableIds = append(routeTableIds, helper.String(val))
		}
	}
	if v, ok := d.GetOk("limit"); ok {
		limit = v.(int)
	}

	respData, found, err := service.DescribeRouteTablesForReplaceRoutesWithRoutePolicyConfig(ctx, routeTableId, filterName, filterValues, needRouterInfo, routeTableIds, limit)
	if err != nil {
		return err
	}

	if !found {
		log.Printf("[CRUD] tencentcloud_vpc_replace_routes_with_route_policy_config id=%s", d.Id())
		d.SetId("")
		return nil
	}

	_ = respData
	_ = d.Set("route_table_id", routeTableId)

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
