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

func ResourceTencentCloudTeoFunctionReplicaV5() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTeoFunctionReplicaV5Create,
		Read:   resourceTencentCloudTeoFunctionReplicaV5Read,
		Update: resourceTencentCloudTeoFunctionReplicaV5Update,
		Delete: resourceTencentCloudTeoFunctionReplicaV5Delete,
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
				Description: "Sort basis of the replica list. Valid values: `created-on`: sort by creation time. Default sorted by the `created-on` attribute.",
			},
			"sort_order": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Sort order of the replica list. Valid values: `asc`: sort in ascending order; `desc`: sort in descending order. Default value: `asc`.",
			},
			"filters": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Filter conditions of the replica list. The maximum value of Filters.Values is 20. If this parameter is not filled in, all function replicas under the function ID will be returned. Valid values: `replica-name`: filter by function replica name, which supports fuzzy query.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Field to be filtered. Valid value: `replica-name`.",
						},
						"values": {
							Type:        schema.TypeSet,
							Required:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "Value of the filtered field. The maximum number of values is 20.",
						},
						"fuzzy": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether to enable fuzzy query.",
						},
					},
				},
			},
			"replica_names": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Names of the function replicas to be deleted when destroying the resource. If not configured, only the replica corresponding to this resource (the `replica_name` in the resource ID) will be deleted by default.",
			},
			"created_on": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation time of the edge function replica.",
			},
			"modified_on": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Last modified time of the edge function replica.",
			},
		},
	}
}

func resourceTencentCloudTeoFunctionReplicaV5Create(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_function_replica_v5.create")()
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
			return resource.NonRetryableError(fmt.Errorf("Create teo function replica v5 failed, Response is nil."))
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create teo function replica v5 failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	log.Printf("[DEBUG]%s create teo function replica v5 done, logId: %s, zoneId: %s, functionId: %s, replicaName: %s", logId, logId, zoneId, functionId, replicaName)
	d.SetId(strings.Join([]string{zoneId, functionId, replicaName}, tccommon.FILED_SP))
	return resourceTencentCloudTeoFunctionReplicaV5Read(d, meta)
}

func resourceTencentCloudTeoFunctionReplicaV5Read(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_function_replica_v5.read")()
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

	if v, ok := d.GetOk("sort_by"); ok {
		request.SortBy = helper.String(v.(string))
	}

	if v, ok := d.GetOk("sort_order"); ok {
		request.SortOrder = helper.String(v.(string))
	}

	if v, ok := d.GetOk("filters"); ok {
		filtersSet := v.([]interface{})
		tmpSet := make([]*teov20220901.AdvancedFilter, 0, len(filtersSet))
		for _, item := range filtersSet {
			filtersMap := item.(map[string]interface{})
			advancedFilter := teov20220901.AdvancedFilter{}
			if v, ok := filtersMap["name"].(string); ok && v != "" {
				advancedFilter.Name = helper.String(v)
			}
			if v, ok := filtersMap["values"]; ok {
				valuesSet := v.(*schema.Set).List()
				for i := range valuesSet {
					advancedFilter.Values = append(advancedFilter.Values, helper.String(valuesSet[i].(string)))
				}
			}
			if v, ok := filtersMap["fuzzy"].(bool); ok {
				advancedFilter.Fuzzy = helper.Bool(v)
			}
			tmpSet = append(tmpSet, &advancedFilter)
		}
		request.Filters = append(request.Filters, tmpSet...)
	}

	request.Filters = append(request.Filters, &teov20220901.AdvancedFilter{
		Name:   helper.String("replica-name"),
		Values: []*string{helper.String(replicaName)},
	})
	request.Limit = helper.Int64(200)

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
		log.Printf("[CRITAL]%s read teo function replica v5 failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	if response == nil || response.Response == nil {
		log.Printf("[CRUD] teo_function_replica_v5 id=%s", d.Id())
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
		log.Printf("[CRUD] teo_function_replica_v5 id=%s", d.Id())
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

	return nil
}

func resourceTencentCloudTeoFunctionReplicaV5Update(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_function_replica_v5.update")()
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

			if result == nil || result.Response == nil {
				return resource.NonRetryableError(fmt.Errorf("Modify teo function replica v5 failed, Response is nil."))
			}

			return nil
		})

		if reqErr != nil {
			log.Printf("[CRITAL]%s update teo function replica v5 failed, reason:%+v", logId, reqErr)
			return reqErr
		}
	}

	return resourceTencentCloudTeoFunctionReplicaV5Read(d, meta)
}

func resourceTencentCloudTeoFunctionReplicaV5Delete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_function_replica_v5.delete")()
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

	if v, ok := d.GetOk("replica_names"); ok {
		replicaNamesSet := v.([]interface{})
		replicaNames := make([]*string, 0, len(replicaNamesSet))
		for _, item := range replicaNamesSet {
			replicaNames = append(replicaNames, helper.String(item.(string)))
		}
		request.ReplicaNames = replicaNames
	} else {
		request.ReplicaNames = []*string{helper.String(replicaName)}
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().DeleteFunctionReplicaWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Delete teo function replica v5 failed, Response is nil."))
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete teo function replica v5 failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return nil
}
