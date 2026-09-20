package teo

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudTeoFunctionReplicaV3() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTeoFunctionReplicaV3Create,
		Read:   resourceTencentCloudTeoFunctionReplicaV3Read,
		Update: resourceTencentCloudTeoFunctionReplicaV3Update,
		Delete: resourceTencentCloudTeoFunctionReplicaV3Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Zone ID.",
			},
			"function_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Function ID.",
			},
			"replica_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Edge function replica name. Limited to 1-50 characters, allowed characters are a-z, 0-9, -, and - cannot be used alone or consecutively, nor at the beginning or end. Replica names must be unique under the same FunctionId.",
			},
			"content": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Edge function replica content. Currently only supports JavaScript code, maximum 5MB.",
			},
			"remark": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Edge function replica description. Maximum 50 characters.",
			},
			"sort_by": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Sort by field for querying the function replica list. Valid value: created-on (creation time). Default: created-on.",
			},
			"sort_order": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Sort order for querying the function replica list. Valid values: asc, desc. Default: asc.",
			},
			"filters": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Filter conditions for querying the function replica list. The maximum number of Values is 20. Supported filter key: replica-name (filter by replica name, fuzzy query supported).",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Field to be filtered.",
						},
						"values": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Description: "Filter values of the field.",
						},
						"fuzzy": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Indicates whether fuzzy query is enabled.",
						},
					},
				},
			},
			"replica_names": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "List of edge function replica names to be deleted.",
			},
			"created_on": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation time of the edge function replica.",
			},
			"modified_on": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Last modification time of the edge function replica.",
			},
			"function_replicas": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of the edge function replicas returned by the DescribeFunctionReplicas API.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"function_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Function ID.",
						},
						"replica_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Edge function replica name.",
						},
						"content": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Edge function replica content (JavaScript).",
						},
						"remark": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Edge function replica description.",
						},
						"created_on": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Creation time of the edge function replica.",
						},
						"modified_on": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Last modification time of the edge function replica.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudTeoFunctionReplicaV3Create(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_function_replica_v3.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId       = tccommon.GetLogId(tccommon.ContextNil)
		ctx         = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request     = teov20220901.NewCreateFunctionReplicaRequest()
		zoneId      string
		functionId  string
		replicaName string
	)

	if v, ok := d.GetOk("zone_id"); ok {
		request.ZoneId = helper.String(v.(string))
		zoneId = v.(string)
	}

	if v, ok := d.GetOk("function_id"); ok {
		request.FunctionId = helper.String(v.(string))
		functionId = v.(string)
	}

	if v, ok := d.GetOk("replica_name"); ok {
		request.ReplicaName = helper.String(v.(string))
		replicaName = v.(string)
	}

	if v, ok := d.GetOk("content"); ok {
		request.Content = helper.String(v.(string))
	}

	if v, ok := d.GetOk("remark"); ok {
		request.Remark = helper.String(v.(string))
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().CreateFunctionReplicaWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create teo function replica v3 failed, Response is nil."))
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create teo_function_replica_v3 failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	log.Printf("[DEBUG]%s create teo_function_replica_v3, zone id: %s, function id: %s, replica name: %s", logId, zoneId, functionId, replicaName)
	d.SetId(strings.Join([]string{zoneId, functionId, replicaName}, tccommon.FILED_SP))
	return resourceTencentCloudTeoFunctionReplicaV3Read(d, meta)
}

func resourceTencentCloudTeoFunctionReplicaV3Read(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_function_replica_v3.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = teov20220901.NewDescribeFunctionReplicasRequest()
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken, id is %s", d.Id())
	}

	zoneId := idSplit[0]
	functionId := idSplit[1]
	replicaName := idSplit[2]

	request.ZoneId = helper.String(zoneId)
	request.FunctionId = helper.String(functionId)
	request.Limit = helper.Int64(200)
	request.Filters = []*teov20220901.AdvancedFilter{
		{
			Name:   helper.String("replica-name"),
			Values: []*string{helper.String(replicaName)},
		},
	}

	if v, ok := d.GetOk("sort_by"); ok {
		request.SortBy = helper.String(v.(string))
	}

	if v, ok := d.GetOk("sort_order"); ok {
		request.SortOrder = helper.String(v.(string))
	}

	if filters, ok := d.GetOk("filters"); ok {
		filtersList := filters.([]interface{})
		for _, filter := range filtersList {
			filterMap := filter.(map[string]interface{})
			advancedFilter := &teov20220901.AdvancedFilter{}
			if v, ok := filterMap["name"].(string); ok && v != "" {
				advancedFilter.Name = helper.String(v)
			}
			if values, ok := filterMap["values"].([]interface{}); ok {
				for _, value := range values {
					advancedFilter.Values = append(advancedFilter.Values, helper.String(value.(string)))
				}
			}
			if v, ok := filterMap["fuzzy"].(bool); ok {
				advancedFilter.Fuzzy = helper.Bool(v)
			}
			request.Filters = append(request.Filters, advancedFilter)
		}
	}

	var response *teov20220901.DescribeFunctionReplicasResponse
	reqErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().DescribeFunctionReplicasWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		response = result
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s read teo_function_replica_v3 failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	if response == nil || response.Response == nil {
		log.Printf("[CRUD] teo_function_replica_v3 id=%s", d.Id())
		d.SetId("")
		return nil
	}

	var targetReplica *teov20220901.FunctionReplica
	for _, replica := range response.Response.FunctionReplicas {
		if replica.ReplicaName != nil && *replica.ReplicaName == replicaName {
			targetReplica = replica
			break
		}
	}

	if targetReplica == nil {
		log.Printf("[CRUD] teo_function_replica_v3 id=%s", d.Id())
		d.SetId("")
		return nil
	}

	_ = d.Set("zone_id", zoneId)
	_ = d.Set("function_id", functionId)
	_ = d.Set("replica_name", replicaName)

	if targetReplica.Content != nil {
		_ = d.Set("content", targetReplica.Content)
	}

	if targetReplica.Remark != nil {
		_ = d.Set("remark", targetReplica.Remark)
	}

	if targetReplica.CreatedOn != nil {
		_ = d.Set("created_on", targetReplica.CreatedOn)
	}

	if targetReplica.ModifiedOn != nil {
		_ = d.Set("modified_on", targetReplica.ModifiedOn)
	}

	functionReplicasList := make([]interface{}, 0, len(response.Response.FunctionReplicas))
	for _, replica := range response.Response.FunctionReplicas {
		functionReplicaMap := make(map[string]interface{})
		if replica.FunctionId != nil {
			functionReplicaMap["function_id"] = replica.FunctionId
		}
		if replica.ReplicaName != nil {
			functionReplicaMap["replica_name"] = replica.ReplicaName
		}
		if replica.Content != nil {
			functionReplicaMap["content"] = replica.Content
		}
		if replica.Remark != nil {
			functionReplicaMap["remark"] = replica.Remark
		}
		if replica.CreatedOn != nil {
			functionReplicaMap["created_on"] = replica.CreatedOn
		}
		if replica.ModifiedOn != nil {
			functionReplicaMap["modified_on"] = replica.ModifiedOn
		}
		functionReplicasList = append(functionReplicasList, functionReplicaMap)
	}
	_ = d.Set("function_replicas", functionReplicasList)

	return nil
}

func resourceTencentCloudTeoFunctionReplicaV3Update(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_function_replica_v3.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = teov20220901.NewModifyFunctionReplicaRequest()
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken, id is %s", d.Id())
	}

	zoneId := idSplit[0]
	functionId := idSplit[1]
	replicaName := idSplit[2]

	request.ZoneId = helper.String(zoneId)
	request.FunctionId = helper.String(functionId)
	request.ReplicaName = helper.String(replicaName)

	needChange := false
	mutableArgs := []string{"content", "remark"}
	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}

	if needChange {
		if v, ok := d.GetOk("content"); ok {
			request.Content = helper.String(v.(string))
		}

		if v, ok := d.GetOk("remark"); ok {
			request.Remark = helper.String(v.(string))
		}

		reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().ModifyFunctionReplicaWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			return nil
		})

		if reqErr != nil {
			log.Printf("[CRITAL]%s update teo_function_replica_v3 failed, reason:%+v", logId, reqErr)
			return reqErr
		}
	}

	return resourceTencentCloudTeoFunctionReplicaV3Read(d, meta)
}

func resourceTencentCloudTeoFunctionReplicaV3Delete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_function_replica_v3.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = teov20220901.NewDeleteFunctionReplicaRequest()
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken, id is %s", d.Id())
	}

	zoneId := idSplit[0]
	functionId := idSplit[1]
	replicaName := idSplit[2]

	request.ZoneId = helper.String(zoneId)
	request.FunctionId = helper.String(functionId)

	replicaNamesList := d.Get("replica_names").([]interface{})
	if len(replicaNamesList) == 0 {
		request.ReplicaNames = []*string{helper.String(replicaName)}
	} else {
		for _, v := range replicaNamesList {
			request.ReplicaNames = append(request.ReplicaNames, helper.String(v.(string)))
		}
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().DeleteFunctionReplicaWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete teo_function_replica_v3 failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return nil
}
