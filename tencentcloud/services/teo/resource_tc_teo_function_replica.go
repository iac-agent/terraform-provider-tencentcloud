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
				Description: "ID of the edge function.",
			},

			"replica_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Edge function replica name. It can only contain lowercase letters, numbers, and hyphens. It must start and end with a letter or number, cannot have consecutive hyphens, and has a maximum length of 50 characters. The replica name must be unique under the same FunctionId.",
			},

			"content": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Edge function replica content. Currently only JavaScript code is supported, with a maximum size of 5MB.",
			},

			"remark": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Edge function replica description. Maximum support of 50 characters.",
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
	}
}

func resourceTencentCloudTeoFunctionReplicaCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_function_replica.create")()
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

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().CreateFunctionReplicaWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			log.Printf("[CRITAL]%s create teo_function_replica failed, logId=%s, result is nil", logId, logId)
			return resource.NonRetryableError(fmt.Errorf("create teo_function_replica failed, Response is nil"))
		}

		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create teo_function_replica failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(strings.Join([]string{zoneId, functionId, replicaName}, tccommon.FILED_SP))

	return resourceTencentCloudTeoFunctionReplicaRead(d, meta)
}

func resourceTencentCloudTeoFunctionReplicaRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_function_replica.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = TeoService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	zoneId := idSplit[0]
	functionId := idSplit[1]
	replicaName := idSplit[2]

	respData, err := service.DescribeTeoFunctionReplicaById(ctx, zoneId, functionId, replicaName)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[CRUD] teo_function_replica id=%s", d.Id())
		d.SetId("")
		return nil
	}

	_ = d.Set("zone_id", zoneId)
	_ = d.Set("function_id", functionId)

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
		_ = d.Set("created_on", respData.CreatedOn)
	}

	if respData.ModifiedOn != nil {
		_ = d.Set("modified_on", respData.ModifiedOn)
	}

	return nil
}

func resourceTencentCloudTeoFunctionReplicaUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_function_replica.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId = tccommon.GetLogId(tccommon.ContextNil)
		ctx   = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
	)

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
			log.Printf("[CRITAL]%s update teo_function_replica failed, reason:%+v", logId, err)
			return err
		}
	}

	return resourceTencentCloudTeoFunctionReplicaRead(d, meta)
}

func resourceTencentCloudTeoFunctionReplicaDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_function_replica.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = teov20220901.NewDeleteFunctionReplicaRequest()
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	zoneId := idSplit[0]
	functionId := idSplit[1]
	replicaName := idSplit[2]

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
		log.Printf("[CRITAL]%s delete teo_function_replica failed, reason:%+v", logId, err)
		return err
	}

	return nil
}
