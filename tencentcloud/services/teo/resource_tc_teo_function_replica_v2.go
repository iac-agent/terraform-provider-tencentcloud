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

func ResourceTencentCloudTeoFunctionReplicaV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTeoFunctionReplicaV2Create,
		Read:   resourceTencentCloudTeoFunctionReplicaV2Read,
		Update: resourceTencentCloudTeoFunctionReplicaV2Update,
		Delete: resourceTencentCloudTeoFunctionReplicaV2Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Site ID.",
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

			// computed
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
	}
}

func resourceTencentCloudTeoFunctionReplicaV2Create(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_function_replica_v2.create")()
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
			return resource.NonRetryableError(fmt.Errorf("Create teo function replica v2 failed, Response is nil."))
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create teo function replica v2 failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	log.Printf("[DEBUG]%s create teo function replica v2 success, zoneId is %s, functionId is %s, replicaName is %s", logId, zoneId, functionId, replicaName)
	if zoneId == "" || functionId == "" || replicaName == "" {
		return fmt.Errorf("Create teo function replica v2 failed, id params is empty, zoneId is %s, functionId is %s, replicaName is %s.", zoneId, functionId, replicaName)
	}

	d.SetId(strings.Join([]string{zoneId, functionId, replicaName}, tccommon.FILED_SP))
	return resourceTencentCloudTeoFunctionReplicaV2Read(d, meta)
}

func resourceTencentCloudTeoFunctionReplicaV2Read(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_function_replica_v2.read")()
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
	request.Filters = []*teov20220901.AdvancedFilter{
		{
			Name:   helper.String("replica-name"),
			Values: []*string{helper.String(replicaName)},
		},
	}
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
		log.Printf("[CRITAL]%s read teo function replica v2 failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	if response == nil || response.Response == nil || len(response.Response.FunctionReplicas) == 0 {
		log.Printf("[CRUD]%s resource `tencentcloud_teo_function_replica_v2` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
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
		log.Printf("[CRUD]%s resource `tencentcloud_teo_function_replica_v2` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
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

func resourceTencentCloudTeoFunctionReplicaV2Update(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_function_replica_v2.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId = tccommon.GetLogId(tccommon.ContextNil)
		ctx   = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken, id is %s", d.Id())
	}

	zoneId := idSplit[0]
	functionId := idSplit[1]
	replicaName := idSplit[2]

	immutableArgs := []string{"zone_id", "function_id", "replica_name"}
	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("Update teo function replica v2 failed, the arg `%s` is immutable, it will force to recreate the resource.", v)
		}
	}

	needChange := false
	mutableArgs := []string{"content", "remark"}
	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}

	if needChange {
		request := teov20220901.NewModifyFunctionReplicaRequest()
		request.ZoneId = helper.String(zoneId)
		request.FunctionId = helper.String(functionId)
		request.ReplicaName = helper.String(replicaName)

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
				return resource.NonRetryableError(fmt.Errorf("Update teo function replica v2 failed, Response is nil."))
			}

			return nil
		})

		if reqErr != nil {
			log.Printf("[CRITAL]%s update teo function replica v2 failed, reason:%+v", logId, reqErr)
			return reqErr
		}
	}

	return resourceTencentCloudTeoFunctionReplicaV2Read(d, meta)
}

func resourceTencentCloudTeoFunctionReplicaV2Delete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_function_replica_v2.delete")()
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
	request.ReplicaNames = []*string{helper.String(replicaName)}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().DeleteFunctionReplicaWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Delete teo function replica v2 failed, Response is nil."))
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete teo function replica v2 failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return nil
}
