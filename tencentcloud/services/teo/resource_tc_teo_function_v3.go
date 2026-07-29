package teo

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudTeoFunctionV3() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTeoFunctionV3Create,
		Read:   resourceTencentCloudTeoFunctionV3Read,
		Update: resourceTencentCloudTeoFunctionV3Update,
		Delete: resourceTencentCloudTeoFunctionV3Delete,
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

			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Function name. It can only contain lowercase letters, numbers, hyphens, must start and end with a letter or number, and can have a maximum length of 30 characters.",
			},

			"content": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Function content, currently only supports JavaScript code, with a maximum size of 5MB.",
			},

			"remark": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Function description, maximum support of 60 characters.",
			},

			"function_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the Function.",
			},

			"domain": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The default domain name for the function.",
			},

			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation time. The time is in Coordinated Universal Time (UTC) and follows the date and time format specified by the ISO 8601 standard.",
			},

			"update_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Modification time. The time is in Coordinated Universal Time (UTC) and follows the date and time format specified by the ISO 8601 standard.",
			},
		},
	}
}

func resourceTencentCloudTeoFunctionV3Create(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_function_v3.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request    = teov20220901.NewCreateFunctionRequest()
		response   = teov20220901.NewCreateFunctionResponse()
		zoneId     string
		functionId string
	)

	if v, ok := d.GetOk("zone_id"); ok {
		zoneId = v.(string)
		request.ZoneId = helper.String(zoneId)
	}

	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}

	if v, ok := d.GetOk("content"); ok {
		request.Content = helper.String(v.(string))
	}

	if v, ok := d.GetOk("remark"); ok {
		request.Remark = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().CreateFunctionWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create teo_function_v3 failed, reason:%+v", logId, err)
		return err
	}

	if response == nil || response.Response == nil {
		return fmt.Errorf("create teo_function_v3 failed, Response is nil")
	}

	log.Printf("[DEBUG]%s create teo_function_v3, logId=%s, d.Id()=%s", logId, logId, d.Id())

	if response.Response.FunctionId == nil || *response.Response.FunctionId == "" {
		return fmt.Errorf("create teo_function_v3 failed, FunctionId is empty")
	}

	functionId = *response.Response.FunctionId

	if _, err := (&resource.StateChangeConf{
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
		Pending:    []string{"false"},
		Refresh:    resourceTeoFunctionV3CreateStateRefreshFunc(ctx, zoneId, functionId),
		Target:     []string{"true"},
		Timeout:    600 * time.Second,
	}).WaitForStateContext(ctx); err != nil {
		return err
	}

	d.SetId(strings.Join([]string{zoneId, functionId}, tccommon.FILED_SP))

	return resourceTencentCloudTeoFunctionV3Read(d, meta)
}

func resourceTencentCloudTeoFunctionV3Read(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_function_v3.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = TeoService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	zoneId := idSplit[0]
	functionId := idSplit[1]

	_ = d.Set("zone_id", zoneId)

	respData, err := service.DescribeTeoFunctionV3ById(ctx, zoneId, functionId)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[CRUD] teo_function_v3 id=%s", d.Id())
		d.SetId("")
		return nil
	}

	if respData.FunctionId != nil {
		_ = d.Set("function_id", respData.FunctionId)
	}

	if respData.Name != nil {
		_ = d.Set("name", respData.Name)
	}

	if respData.Remark != nil {
		_ = d.Set("remark", respData.Remark)
	}

	if respData.Content != nil {
		_ = d.Set("content", respData.Content)
	}

	if respData.Domain != nil {
		_ = d.Set("domain", respData.Domain)
	}

	if respData.CreateTime != nil {
		_ = d.Set("create_time", respData.CreateTime)
	}

	if respData.UpdateTime != nil {
		_ = d.Set("update_time", respData.UpdateTime)
	}

	return nil
}

func resourceTencentCloudTeoFunctionV3Update(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_function_v3.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId = tccommon.GetLogId(tccommon.ContextNil)
		ctx   = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
	)

	immutableArgs := []string{"name"}
	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	zoneId := idSplit[0]
	functionId := idSplit[1]

	needChange := false
	mutableArgs := []string{"remark", "content"}
	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}

	if needChange {
		request := teov20220901.NewModifyFunctionRequest()

		request.ZoneId = helper.String(zoneId)
		request.FunctionId = helper.String(functionId)

		if v, ok := d.GetOk("remark"); ok {
			request.Remark = helper.String(v.(string))
		}

		if v, ok := d.GetOk("content"); ok {
			request.Content = helper.String(v.(string))
		}

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().ModifyFunctionWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update teo_function_v3 failed, reason:%+v", logId, err)
			return err
		}
	}

	return resourceTencentCloudTeoFunctionV3Read(d, meta)
}

func resourceTencentCloudTeoFunctionV3Delete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_function_v3.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId = tccommon.GetLogId(tccommon.ContextNil)
		ctx   = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	zoneId := idSplit[0]
	functionId := idSplit[1]

	request := teov20220901.NewDeleteFunctionRequest()
	request.ZoneId = helper.String(zoneId)
	request.FunctionId = helper.String(functionId)

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().DeleteFunctionWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s delete teo_function_v3 failed, reason:%+v", logId, err)
		return err
	}

	return nil
}

func resourceTeoFunctionV3CreateStateRefreshFunc(ctx context.Context, zoneId string, functionId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		meta := tccommon.ProviderMetaFromContext(ctx)
		if meta == nil {
			return nil, "", fmt.Errorf("resource data can not be nil")
		}

		req := teov20220901.NewDescribeFunctionsRequest()
		req.ZoneId = helper.String(zoneId)
		req.FunctionIds = []*string{helper.String(functionId)}

		resp, err := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().DescribeFunctionsWithContext(ctx, req)
		if err != nil {
			return nil, "", err
		}

		if resp == nil || resp.Response == nil {
			return nil, "false", nil
		}

		if len(resp.Response.Functions) < 1 {
			return nil, "false", nil
		}

		function := resp.Response.Functions[0]
		if function.Domain != nil && *function.Domain != "" {
			return resp.Response, "true", nil
		}

		return resp.Response, "false", nil
	}
}
