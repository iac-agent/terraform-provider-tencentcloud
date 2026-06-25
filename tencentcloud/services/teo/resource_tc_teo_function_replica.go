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

func ResourceTencentCloudTeoFunctionReplica() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTeoFunctionReplicaCreate,
		Read:   resourceTencentCloudTeoFunctionReplicaRead,
		Update: resourceTencentCloudTeoFunctionReplicaUpdate,
		Delete: resourceTencentCloudTeoFunctionReplicaDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the site.",
			},

			"function_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the Function.",
			},

			"replica_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the function replica. It can only contain lowercase letters, numbers, hyphens, must start and end with a letter or number, and can have a maximum length of 50 characters.",
			},

			"content": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Function replica content, currently only supports JavaScript code, with a maximum size of 5MB.",
			},

			"remark": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Function replica description, maximum support of 50 characters.",
			},

			"created_on": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation time. The time is in Coordinated Universal Time (UTC) and follows the date and time format specified by the ISO 8601 standard.",
			},

			"modified_on": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Modification time. The time is in Coordinated Universal Time (UTC) and follows the date and time format specified by the ISO 8601 standard.",
			},
		},
	}
}

func resourceTencentCloudTeoFunctionReplicaCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_function_replica.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	var (
		zoneId      string
		functionId  string
		replicaName string
	)
	var (
		request  = teov20220901.NewCreateFunctionReplicaRequest()
		response = teov20220901.NewCreateFunctionReplicaResponse()
	)

	if v, ok := d.GetOk("zone_id"); ok {
		zoneId = v.(string)
	}
	if v, ok := d.GetOk("function_id"); ok {
		functionId = v.(string)
	}
	if v, ok := d.GetOk("replica_name"); ok {
		replicaName = v.(string)
	}

	request.ZoneId = helper.String(zoneId)
	request.FunctionId = helper.String(functionId)
	request.ReplicaName = helper.String(replicaName)

	if v, ok := d.GetOk("content"); ok {
		request.Content = helper.String(v.(string))
	}

	if v, ok := d.GetOk("remark"); ok {
		request.Remark = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().CreateFunctionReplicaWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create teo function_replica failed, reason:%+v", logId, err)
		return err
	}

	log.Printf("[DEBUG]%s create teo function_replica response: %s, current d.Id()=%s", logId, response.ToJsonString(), d.Id())
	if response == nil || response.Response == nil {
		return fmt.Errorf("create teo function_replica failed, empty response")
	}

	d.SetId(strings.Join([]string{zoneId, functionId, replicaName}, tccommon.FILED_SP))

	return resourceTencentCloudTeoFunctionReplicaRead(d, meta)
}

func resourceTencentCloudTeoFunctionReplicaRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_function_replica.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	zoneId := idSplit[0]
	functionId := idSplit[1]
	replicaName := idSplit[2]

	_ = d.Set("zone_id", zoneId)
	_ = d.Set("function_id", functionId)

	var (
		offset int64 = 0
		limit  int64 = 200
	)

	var replicas []*teov20220901.FunctionReplica

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		request := teov20220901.NewDescribeFunctionReplicasRequest()
		request.ZoneId = helper.String(zoneId)
		request.FunctionId = helper.String(functionId)
		request.Offset = helper.Int64(offset)
		request.Limit = helper.Int64(limit)

		request.Filters = []*teov20220901.AdvancedFilter{
			{
				Name:   helper.String("replica-name"),
				Values: []*string{helper.String(replicaName)},
			},
		}

		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().DescribeFunctionReplicasWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("describe teo function_replica failed, empty response"))
		}
		replicas = result.Response.FunctionReplicas
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read teo function_replica failed, reason:%+v", logId, err)
		return err
	}

	if len(replicas) == 0 {
		log.Printf("[WARN]%s resource `teo_function_replica` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	// exact match on replica_name since the API does fuzzy matching
	var matchedReplica *teov20220901.FunctionReplica
	for _, replica := range replicas {
		if replica.ReplicaName != nil && *replica.ReplicaName == replicaName {
			matchedReplica = replica
			break
		}
	}

	if matchedReplica == nil {
		log.Printf("[WARN]%s resource `teo_function_replica` [%s] not found (exact match failed), please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	if matchedReplica.ReplicaName != nil {
		_ = d.Set("replica_name", matchedReplica.ReplicaName)
	}

	if matchedReplica.Content != nil {
		_ = d.Set("content", matchedReplica.Content)
	}

	if matchedReplica.Remark != nil {
		_ = d.Set("remark", matchedReplica.Remark)
	}

	if matchedReplica.CreatedOn != nil {
		_ = d.Set("created_on", matchedReplica.CreatedOn)
	}

	if matchedReplica.ModifiedOn != nil {
		_ = d.Set("modified_on", matchedReplica.ModifiedOn)
	}

	return nil
}

func resourceTencentCloudTeoFunctionReplicaUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_function_replica.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	zoneId := idSplit[0]
	functionId := idSplit[1]
	replicaName := idSplit[2]

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

		err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().ModifyFunctionReplicaWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update teo function_replica failed, reason:%+v", logId, err)
			return err
		}
	}

	return resourceTencentCloudTeoFunctionReplicaRead(d, meta)
}

func resourceTencentCloudTeoFunctionReplicaDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_function_replica.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	zoneId := idSplit[0]
	functionId := idSplit[1]
	replicaName := idSplit[2]

	var (
		request  = teov20220901.NewDeleteFunctionReplicaRequest()
		response = teov20220901.NewDeleteFunctionReplicaResponse()
	)

	request.ZoneId = helper.String(zoneId)
	request.FunctionId = helper.String(functionId)
	request.ReplicaNames = []*string{helper.String(replicaName)}

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().DeleteFunctionReplicaWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s delete teo function_replica failed, reason:%+v", logId, err)
		return err
	}

	_ = response
	return nil
}
