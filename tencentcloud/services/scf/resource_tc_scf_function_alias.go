package scf

import (
	"context"
	"fmt"
	"log"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	scf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/scf/v20180416"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudScfFunctionAlias() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudScfFunctionAliasCreate,
		Read:   resourceTencentCloudScfFunctionAliasRead,
		Update: resourceTencentCloudScfFunctionAliasUpdate,
		Delete: resourceTencentCloudScfFunctionAliasDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Alias 名称，其中 必须 是 唯一 在 函数，可以 contain 1 到 64 letters，digits，_，和 -，和 必须 begin 使用 letter。",
			},

			"function_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Function 名称",
			},

			"function_version": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Master 版本 pointed 到 通过 alias。",
			},

			"namespace": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Function 命名空间。",
			},

			"routing_config": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Request routing 配置 的 alias。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"additional_version_weights": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Additional 版本 使用 random 权重-based routing。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"version": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Function 版本 名称",
									},
									"weight": {
										Type:        schema.TypeFloat,
										Required:    true,
										Description: "版本 权重",
									},
								},
							},
						},
						"additional_version_matches": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Additional 版本 使用 规则-based routing。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"version": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Function 版本 名称",
									},
									"key": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Matching 规则 键 当 API 是 called，pass 在 键 到 路由 请求 到 指定 版本 based 在 matching ruleHeader 方法:Enter invoke.headers.用户 对于 键 和 pass 在 RoutingKey:{用户:值} 当 invoking 函数 through invoke 对于 invocation based 在 规则 matching。",
									},
									"method": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Match 方法. 有效 值:范围: Range matchexact: exact 字符串 match。",
									},
									"expression": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Rule requirements 对于 范围 match:It should 是 described 在 open 或 closed 范围，i.e.，(,b) 或 [,b]，其中 both 和 b 是 integersRule requirements 对于 exact match:Exact 字符串 match。",
									},
								},
							},
						},
					},
				},
			},

			"description": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Alias 描述 信息。",
			},
		},
	}
}

func resourceTencentCloudScfFunctionAliasCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_scf_function_alias.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request      = scf.NewCreateAliasRequest()
		namespace    string
		functionName string
		name         string
	)
	if v, ok := d.GetOk("name"); ok {
		name = v.(string)
		request.Name = helper.String(v.(string))
	}

	if v, ok := d.GetOk("function_name"); ok {
		functionName = v.(string)
		request.FunctionName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("function_version"); ok {
		request.FunctionVersion = helper.String(v.(string))
	}

	if v, ok := d.GetOk("namespace"); ok {
		namespace = v.(string)
		request.Namespace = helper.String(v.(string))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "routing_config"); ok {
		routingConfig := scf.RoutingConfig{}
		if v, ok := dMap["additional_version_weights"]; ok {
			for _, item := range v.([]interface{}) {
				additionalVersionWeightsMap := item.(map[string]interface{})
				versionWeight := scf.VersionWeight{}
				if v, ok := additionalVersionWeightsMap["version"]; ok {
					versionWeight.Version = helper.String(v.(string))
				}
				if v, ok := additionalVersionWeightsMap["weight"]; ok {
					versionWeight.Weight = helper.Float64(v.(float64))
				}
				routingConfig.AdditionalVersionWeights = append(routingConfig.AdditionalVersionWeights, &versionWeight)
			}
		}
		if v, ok := dMap["additional_version_matches"]; ok {
			for _, item := range v.([]interface{}) {
				addtionVersionMatchsMap := item.(map[string]interface{})
				versionMatch := scf.VersionMatch{}
				if v, ok := addtionVersionMatchsMap["version"]; ok {
					versionMatch.Version = helper.String(v.(string))
				}
				if v, ok := addtionVersionMatchsMap["key"]; ok {
					versionMatch.Key = helper.String(v.(string))
				}
				if v, ok := addtionVersionMatchsMap["method"]; ok {
					versionMatch.Method = helper.String(v.(string))
				}
				if v, ok := addtionVersionMatchsMap["expression"]; ok {
					versionMatch.Expression = helper.String(v.(string))
				}
				routingConfig.AddtionVersionMatchs = append(routingConfig.AddtionVersionMatchs, &versionMatch)
			}
		}
		request.RoutingConfig = &routingConfig
	}

	if v, ok := d.GetOk("description"); ok {
		request.Description = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseScfClient().CreateAlias(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create scf FunctionAlias failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(namespace + tccommon.FILED_SP + functionName + tccommon.FILED_SP + name)

	return resourceTencentCloudScfFunctionAliasRead(d, meta)
}

func resourceTencentCloudScfFunctionAliasRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_scf_function_alias.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := ScfService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	namespace := idSplit[0]
	functionName := idSplit[1]
	name := idSplit[2]

	functionAlias, err := service.DescribeScfFunctionAliasById(ctx, namespace, functionName, name)
	if err != nil {
		return err
	}

	if functionAlias == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `ScfFunctionAlias` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("name", name)
	_ = d.Set("function_name", functionName)
	_ = d.Set("namespace", namespace)

	if functionAlias.Response.FunctionVersion != nil {
		_ = d.Set("function_version", functionAlias.Response.FunctionVersion)
	}

	if functionAlias.Response.RoutingConfig != nil {
		routingConfigMap := map[string]interface{}{}

		routingConfig := functionAlias.Response.RoutingConfig

		if routingConfig.AdditionalVersionWeights != nil {
			additionalVersionWeightsList := []interface{}{}
			for _, additionalVersionWeights := range routingConfig.AdditionalVersionWeights {
				additionalVersionWeightsMap := map[string]interface{}{}

				if additionalVersionWeights.Version != nil {
					additionalVersionWeightsMap["version"] = additionalVersionWeights.Version
				}

				if additionalVersionWeights.Weight != nil {
					additionalVersionWeightsMap["weight"] = additionalVersionWeights.Weight
				}

				additionalVersionWeightsList = append(additionalVersionWeightsList, additionalVersionWeightsMap)
			}

			routingConfigMap["additional_version_weights"] = additionalVersionWeightsList
		}

		if routingConfig.AddtionVersionMatchs != nil {
			addtionVersionMatchsList := []interface{}{}
			for _, addtionVersionMatchs := range routingConfig.AddtionVersionMatchs {
				addtionVersionMatchsMap := map[string]interface{}{}

				if addtionVersionMatchs.Version != nil {
					addtionVersionMatchsMap["version"] = addtionVersionMatchs.Version
				}

				if addtionVersionMatchs.Key != nil {
					addtionVersionMatchsMap["key"] = addtionVersionMatchs.Key
				}

				if addtionVersionMatchs.Method != nil {
					addtionVersionMatchsMap["method"] = addtionVersionMatchs.Method
				}

				if addtionVersionMatchs.Expression != nil {
					addtionVersionMatchsMap["expression"] = addtionVersionMatchs.Expression
				}

				addtionVersionMatchsList = append(addtionVersionMatchsList, addtionVersionMatchsMap)
			}

			routingConfigMap["additional_version_matches"] = addtionVersionMatchsList
		}

		_ = d.Set("routing_config", []interface{}{routingConfigMap})
	}

	if functionAlias.Response.Description != nil {
		_ = d.Set("description", functionAlias.Response.Description)
	}

	return nil
}

func resourceTencentCloudScfFunctionAliasUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_scf_function_alias.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	request := scf.NewUpdateAliasRequest()

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	namespace := idSplit[0]
	functionName := idSplit[1]
	name := idSplit[2]

	request.Namespace = &namespace
	request.FunctionName = &functionName
	request.Name = &name
	request.FunctionVersion = helper.String(d.Get("function_version").(string))

	mutableArgs := []string{"routing_config", "description"}

	needChange := false
	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}

	if needChange {

		if dMap, ok := helper.InterfacesHeadMap(d, "routing_config"); ok {
			routingConfig := scf.RoutingConfig{}
			if v, ok := dMap["additional_version_weights"]; ok {
				for _, item := range v.([]interface{}) {
					additionalVersionWeightsMap := item.(map[string]interface{})
					versionWeight := scf.VersionWeight{}
					if v, ok := additionalVersionWeightsMap["version"]; ok {
						versionWeight.Version = helper.String(v.(string))
					}
					if v, ok := additionalVersionWeightsMap["weight"]; ok {
						versionWeight.Weight = helper.Float64(v.(float64))
					}
					routingConfig.AdditionalVersionWeights = append(routingConfig.AdditionalVersionWeights, &versionWeight)
				}
			}
			if v, ok := dMap["additional_version_matches"]; ok {
				for _, item := range v.([]interface{}) {
					addtionVersionMatchsMap := item.(map[string]interface{})
					versionMatch := scf.VersionMatch{}
					if v, ok := addtionVersionMatchsMap["version"]; ok {
						versionMatch.Version = helper.String(v.(string))
					}
					if v, ok := addtionVersionMatchsMap["key"]; ok {
						versionMatch.Key = helper.String(v.(string))
					}
					if v, ok := addtionVersionMatchsMap["method"]; ok {
						versionMatch.Method = helper.String(v.(string))
					}
					if v, ok := addtionVersionMatchsMap["expression"]; ok {
						versionMatch.Expression = helper.String(v.(string))
					}
					routingConfig.AddtionVersionMatchs = append(routingConfig.AddtionVersionMatchs, &versionMatch)
				}
			}
			request.RoutingConfig = &routingConfig
		}

		if v, ok := d.GetOk("description"); ok {
			request.Description = helper.String(v.(string))
		}

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseScfClient().UpdateAlias(request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update scf FunctionAlias failed, reason:%+v", logId, err)
			return err
		}
	}

	return resourceTencentCloudScfFunctionAliasRead(d, meta)
}

func resourceTencentCloudScfFunctionAliasDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_scf_function_alias.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := ScfService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	namespace := idSplit[0]
	functionName := idSplit[1]
	name := idSplit[2]

	if err := service.DeleteScfFunctionAliasById(ctx, namespace, functionName, name); err != nil {
		return err
	}

	return nil
}
