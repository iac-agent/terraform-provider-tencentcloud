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

func ResourceTencentCloudTeoFunctionReplicaV1() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTeoFunctionReplicaV1Create,
		Read:   resourceTencentCloudTeoFunctionReplicaV1Read,
		Update: resourceTencentCloudTeoFunctionReplicaV1Update,
		Delete: resourceTencentCloudTeoFunctionReplicaV1Delete,
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
				Description: "ID of the edge function.",
			},

			"replica_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the edge function replica. It can only contain lowercase letters, numbers, hyphens, must start and end with a letter or number, with a maximum length of 50 characters. The replica name must be unique under the same FunctionId.",
			},

			"content": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Content of the edge function replica. Currently only supports JavaScript code, with a maximum size of 5MB.",
			},

			"remark": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Description of the edge function replica. Maximum support of 50 characters.",
			},

			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation time of the edge function replica.",
			},

			"update_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Update time of the edge function replica.",
			},
		},
	}
}

func resourceTencentCloudTeoFunctionReplicaV1Create(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_function_replica_v1.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	var (
		zoneId      string
		functionId  string
		replicaName string
	)
	var (
		request = teov20220901.NewCreateFunctionReplicaRequest()
	)

	if v, ok := d.GetOk("zone_id"); ok {
		zoneId = v.(string)
	}
	request.ZoneId = helper.String(zoneId)

	if v, ok := d.GetOk("function_id"); ok {
		functionId = v.(string)
	}
	request.FunctionId = helper.String(functionId)

	if v, ok := d.GetOk("replica_name"); ok {
		replicaName = v.(string)
	}
	request.ReplicaName = helper.String(replicaName)

	if v, ok := d.GetOk("content"); ok {
		request.Content = helper.String(v.(string))
	}

	if v, ok := d.GetOk("remark"); ok {
		request.Remark = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().CreateFunctionReplicaWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("create function_replica_v1 response is nil, logId=%s", logId))
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create function_replica_v1 failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(strings.Join([]string{zoneId, functionId, replicaName}, tccommon.FILED_SP))

	return resourceTencentCloudTeoFunctionReplicaV1Read(d, meta)
}

func resourceTencentCloudTeoFunctionReplicaV1Read(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_function_replica_v1.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	service := TeoService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	zoneId := idSplit[0]
	functionId := idSplit[1]
	replicaName := idSplit[2]

	_ = d.Set("zone_id", zoneId)
	_ = d.Set("function_id", functionId)

	respData, err := service.DescribeTeoFunctionReplicaV1ByFilter(ctx, zoneId, functionId, replicaName)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[CRUD] function_replica_v1 id=%s", d.Id())
		d.SetId("")
		return nil
	}

	if respData.FunctionId != nil {
		_ = d.Set("function_id", respData.FunctionId)
	}

	if respData.ReplicaName != nil {
		_ = d.Set("replica_name", respData.ReplicaName)
	}

	if respData.Content != nil {
		_ = d.Set("content", respData.Content)
	}

	if respData.Remark != nil {
		_ = d.Set("remark", respData.Remark)
	}

	if respData.CreatedOn != nil {
		_ = d.Set("create_time", respData.CreatedOn)
	}

	if respData.ModifiedOn != nil {
		_ = d.Set("update_time", respData.ModifiedOn)
	}

	return nil
}

func resourceTencentCloudTeoFunctionReplicaV1Update(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_function_replica_v1.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	immutableArgs := []string{"zone_id", "function_id", "replica_name"}
	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

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

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().ModifyFunctionReplicaWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update function_replica_v1 failed, reason:%+v", logId, err)
			return err
		}
	}

	return resourceTencentCloudTeoFunctionReplicaV1Read(d, meta)
}

func resourceTencentCloudTeoFunctionReplicaV1Delete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_function_replica_v1.delete")()
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
		request = teov20220901.NewDeleteFunctionReplicaRequest()
	)

	request.ZoneId = helper.String(zoneId)
	request.FunctionId = helper.String(functionId)
	request.ReplicaNames = []*string{helper.String(replicaName)}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().DeleteFunctionReplicaWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s delete function_replica_v1 failed, reason:%+v", logId, err)
		return err
	}

	return nil
}
