package teo

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"strings"
	"text/template"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudTeoFunctionV6() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTeoFunctionV6Create,
		Read:   resourceTencentCloudTeoFunctionV6Read,
		Update: resourceTencentCloudTeoFunctionV6Update,
		Delete: resourceTencentCloudTeoFunctionV6Delete,
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
				ForceNew:    true,
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

			"domain_compliance_restrictions": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Domain compliance restrictions list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"reason": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Restriction reason.",
						},
						"region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Restricted region.",
						},
					},
				},
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

func resourceTencentCloudTeoFunctionV6Create(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_function_v6.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	var (
		zoneId     string
		functionId string
	)
	var (
		request  = teov20220901.NewCreateFunctionRequest()
		response = teov20220901.NewCreateFunctionResponse()
	)

	if v, ok := d.GetOk("zone_id"); ok {
		zoneId = v.(string)
	}

	request.ZoneId = helper.String(zoneId)

	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}

	if v, ok := d.GetOk("remark"); ok {
		request.Remark = helper.String(v.(string))
	}

	if v, ok := d.GetOk("content"); ok {
		request.Content = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().CreateFunctionWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create teo_function_v6 failed, Response is nil."))
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create teo_function_v6 failed, reason:%+v", logId, err)
		return err
	}

	log.Printf("[INFO]%s create teo_function_v6, id=%s", logId, *response.Response.FunctionId)

	if response.Response.FunctionId == nil || *response.Response.FunctionId == "" {
		return fmt.Errorf("FunctionId is nil or empty.")
	}

	functionId = *response.Response.FunctionId

	d.SetId(strings.Join([]string{zoneId, functionId}, tccommon.FILED_SP))

	if _, err := (&resource.StateChangeConf{
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
		Pending:    []string{"false"},
		Refresh:    resourceTeoFunctionV6CreateStateRefreshFunc(ctx, zoneId, functionId),
		Target:     []string{"true"},
		Timeout:    600 * time.Second,
	}).WaitForStateContext(ctx); err != nil {
		return err
	}

	return resourceTencentCloudTeoFunctionV6Read(d, meta)
}

func resourceTencentCloudTeoFunctionV6Read(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_function_v6.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	service := TeoService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	zoneId := idSplit[0]
	functionId := idSplit[1]

	respData, err := service.DescribeTeoFunctionV6ById(ctx, zoneId, functionId)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[CRUD] teo_function_v6 id=%s", d.Id())
		d.SetId("")
		return nil
	}

	_ = d.Set("zone_id", zoneId)

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

	if respData.DomainComplianceRestrictions != nil {
		restrictionsList := make([]map[string]interface{}, 0, len(respData.DomainComplianceRestrictions))
		for _, restriction := range respData.DomainComplianceRestrictions {
			restrictionMap := map[string]interface{}{}
			if restriction.Reason != nil {
				restrictionMap["reason"] = restriction.Reason
			}
			if restriction.Region != nil {
				restrictionMap["region"] = restriction.Region
			}
			restrictionsList = append(restrictionsList, restrictionMap)
		}
		_ = d.Set("domain_compliance_restrictions", restrictionsList)
	}

	if respData.CreateTime != nil {
		_ = d.Set("create_time", respData.CreateTime)
	}

	if respData.UpdateTime != nil {
		_ = d.Set("update_time", respData.UpdateTime)
	}

	return nil
}

func resourceTencentCloudTeoFunctionV6Update(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_function_v6.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	immutableArgs := []string{"name", "zone_id"}
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
			log.Printf("[CRITAL]%s update teo_function_v6 failed, reason:%+v", logId, err)
			return err
		}
	}

	return resourceTencentCloudTeoFunctionV6Read(d, meta)
}

func resourceTencentCloudTeoFunctionV6Delete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_function_v6.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	zoneId := idSplit[0]
	functionId := idSplit[1]

	var (
		request  = teov20220901.NewDeleteFunctionRequest()
		response = teov20220901.NewDeleteFunctionResponse()
	)

	request.ZoneId = helper.String(zoneId)

	request.FunctionId = helper.String(functionId)

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().DeleteFunctionWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s delete teo_function_v6 failed, reason:%+v", logId, err)
		return err
	}

	_ = response
	return nil
}

func resourceTeoFunctionV6CreateStateRefreshFunc(ctx context.Context, zoneId string, functionId string) resource.StateRefreshFunc {
	var req *teov20220901.DescribeFunctionsRequest
	t := template.New("gotpl")
	var tplObj *template.Template
	return func() (interface{}, string, error) {
		meta := tccommon.ProviderMetaFromContext(ctx)
		if meta == nil {
			return nil, "", fmt.Errorf("resource data can not be nil")
		}
		if req == nil {
			d := tccommon.ResourceDataFromContext(ctx)
			if d == nil {
				return nil, "", fmt.Errorf("resource data can not be nil")
			}
			_ = d
			req = teov20220901.NewDescribeFunctionsRequest()
			req.ZoneId = helper.String(zoneId)

			req.FunctionIds = []*string{helper.String(functionId)}

		}
		resp, err := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().DescribeFunctionsWithContext(ctx, req)
		if err != nil {
			return nil, "", err
		}
		if resp == nil || resp.Response == nil {
			return nil, "", nil
		}
		if tplObj == nil {
			tplObj, err = t.Parse("{{ if .Functions }}{{ $firstFunction := index .Functions 0 }}{{ if $firstFunction.Domain }}{{ true }}{{ else }}{{ false }}{{ end }}{{ end }}")
			if err != nil {
				return resp.Response, "", fmt.Errorf("parse state go-template error: %w", err)
			}
		}
		stream := new(bytes.Buffer)
		if err := tplObj.Execute(stream, resp.Response); err != nil {
			return resp.Response, "", err
		}
		stateBytes, err := io.ReadAll(stream)
		if err != nil {
			return resp.Response, "", err
		}
		state := string(stateBytes)
		return resp.Response, state, nil
	}
}
