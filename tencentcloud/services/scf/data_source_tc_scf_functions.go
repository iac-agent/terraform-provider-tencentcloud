package scf

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/pkg/errors"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudScfFunctions() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudScfFunctionsRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "名称 SCF 函数 到 是 queried。",
			},
			"namespace": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Namespace 的 SCF 函数 到 是 queried。",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "描述 SCF 函数 到 是 queried。",
			},
			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "标签 的 SCF 函数 到 是 queried，可以 使用 up 到 10 标签",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "用于保存结果。",
			},

			"functions": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "An 信息 列表 functions. Each element 包含following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "名称 SCF 函数。",
						},
						"handler": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Handler 的 SCF 函数。",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "描述 SCF 函数。",
						},
						"mem_size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Memory 大小 的 SCF 函数 runtime，单位 是 M。",
						},
						"timeout": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Timeout 的 SCF 函数 最大 执行时间，单位 是 second。",
						},
						"environment": {
							Type:        schema.TypeMap,
							Computed:    true,
							Description: "Environment variable 的 SCF 函数。",
						},
						"runtime": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Runtime 的 SCF 函数。",
						},
						"vpc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "私有网络 ID SCF 函数。",
						},
						"subnet_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "子网 ID SCF 函数。",
						},
						"namespace": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Namespace 的 SCF 函数。",
						},
						"role": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "CAM 角色 的 SCF 函数。",
						},
						"cls_logset_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "CLS logset ID SCF 函数。",
						},
						"cls_topic_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "CLS 主题 ID SCF 函数。",
						},
						"l5_enable": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "是否enable L5。",
						},
						"enable_public_net": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "是否public net 已启用",
						},
						"enable_eip_config": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "是否EIP 已启用",
						},
						"tags": {
							Type:        schema.TypeMap,
							Computed:    true,
							Description: "标签 的 SCF 函数。",
						},
						"async_run_enable": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Whether asynchronous attribute 是 已启用",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "创建时间 的 SCF 函数。",
						},
						"modify_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "修改时间 的 SCF 函数。",
						},
						"code_size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "代码 大小 的 SCF 函数。",
						},
						"code_result": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "代码 结果 的 SCF 函数。",
						},
						"code_error": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "代码 错误 的 SCF 函数。",
						},
						"err_no": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Errno 的 SCF 函数。",
						},
						"install_dependency": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "是否automatically install dependencies。",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "状态 SCF 函数。",
						},
						"status_desc": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "状态 描述 SCF 函数。",
						},
						"eip_fixed": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether EIP 是 fixed IP。",
						},
						"eips": {
							Type:        schema.TypeList,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Computed:    true,
							Description: "EIP 列表 SCF 函数。",
						},
						"host": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "主机 的 SCF 函数。",
						},
						"vip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Vip 的 SCF 函数。",
						},
						"trigger_info": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Trigger details 列表 SCF 函数. Each element 包含following attributes:",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "名称 SCF 函数 触发器。",
									},
									"type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "类型 SCF 函数 触发器。",
									},
									"trigger_desc": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "TriggerDesc 的 SCF 函数 触发器。",
									},
									"enable": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "是否enable SCF 函数 触发器。",
									},
									"create_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "创建时间 的 SCF 函数 触发器。",
									},
									"modify_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "修改时间 的 SCF 函数 触发器。",
									},
									"custom_argument": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "用户-defined 参数 的 SCF 函数 触发器。",
									},
								},
							},
						},
						"image_config": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Image 的 SCF 函数，conflict 使用 `cos_bucket_name`，`cos_object_name`，`cos_bucket_region`，`zip_file`。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"image_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "镜像 类型 personal 或 enterprise。",
									},
									"image_uri": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "uri 的 镜像。",
									},
									"registry_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "registry ID TCR. 当 镜像 类型 是 enterprise，它 必须 是 集合。",
									},
									"entry_point": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "entrypoint 的 app。",
									},
									"command": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "command 的 entrypoint。",
									},
									"args": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "参数 的 command。",
									},
									"container_image_accelerate": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Image accelerate switch。",
									},
									"image_port": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Image 函数 端口 setting. 默认为 `9000`，-1 表示no 端口 mirroring 函数. Other 值 ranges 0 ~ 65535。",
									},
								},
							},
						},
						"dns_cache": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "是否enable Dns caching capability，仅 EVENT 函数 是 支持. 默认为 false。",
						},
						"intranet_config": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Intranet 访问 配置。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ip_fixed": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "是否enable fixed intranet IP，ENABLE 是 已启用，DISABLE 是 已禁用",
									},
									"ip_address": {
										Type: schema.TypeList,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Computed:    true,
										Description: "如果 fixed intranet IP 是 已启用，此 字段 返回IP 列表 使用。",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudScfFunctionsRead(d *schema.ResourceData, m interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_scf_functions.read")()
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := ScfService{client: m.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var (
		name      *string
		namespace *string
		desc      *string
	)

	if raw, ok := d.GetOk("name"); ok {
		name = helper.String(raw.(string))
	}
	if raw, ok := d.GetOk("namespace"); ok {
		namespace = helper.String(raw.(string))
	}
	if raw, ok := d.GetOk("description"); ok {
		desc = helper.String(raw.(string))
	}

	tags := helper.GetTags(d, "tags")
	if len(tags) > 10 {
		return errors.Errorf("can't set more than 10 tags")
	}

	respFunctions, err := service.DescribeFunctions(ctx, name, namespace, desc, tags)
	if err != nil {
		log.Printf("[CRITAL]%s get function list failed: %+v", logId, err)
		return err
	}

	functions := make([]map[string]interface{}, 0, len(respFunctions))
	ids := make([]string, 0, len(respFunctions))

	for _, fn := range respFunctions {
		ids = append(ids, fmt.Sprintf("%s+%s", *fn.Namespace, *fn.FunctionName))

		m := map[string]interface{}{
			"name":        fn.FunctionName,
			"description": fn.Description,
			"runtime":     fn.Runtime,
			"namespace":   fn.Namespace,
			"create_time": fn.AddTime,
			"modify_time": fn.ModTime,
			"status":      fn.Status,
			"status_desc": fn.StatusDesc,
		}

		rawResp, err := service.DescribeFunction(ctx, *fn.FunctionName, *fn.Namespace)
		if err != nil {
			log.Printf("[CRITAL]%s read function detail failed: %+v", logId, err)
			return err
		}

		// if function deleted, ignore it
		if rawResp == nil {
			continue
		}

		resp := rawResp.Response

		m["handler"] = resp.Handler
		m["mem_size"] = resp.MemorySize
		m["timeout"] = resp.Timeout

		env := make(map[string]string, len(resp.Environment.Variables))
		for _, v := range resp.Environment.Variables {
			env[*v.Key] = *v.Value
		}
		m["environment"] = env

		m["vpc_id"] = resp.VpcConfig.VpcId
		m["subnet_id"] = resp.VpcConfig.SubnetId
		m["role"] = resp.Role
		m["cls_logset_id"] = resp.ClsLogsetId
		m["cls_topic_id"] = resp.ClsTopicId
		m["code_size"] = resp.CodeSize
		m["code_result"] = resp.CodeResult
		m["code_error"] = resp.CodeError
		m["err_no"] = resp.ErrNo
		m["install_dependency"] = *resp.InstallDependency == "TRUE"
		m["eip_fixed"] = *resp.EipConfig.EipFixed == "TRUE"
		m["eips"] = resp.EipConfig.Eips
		m["host"] = resp.AccessInfo.Host
		m["vip"] = resp.AccessInfo.Vip
		m["l5_enable"] = *resp.L5Enable == "TRUE"
		if resp.PublicNetConfig != nil {
			m["enable_public_net"] = *resp.PublicNetConfig.PublicNetStatus == "ENABLE"
			m["enable_eip_config"] = *resp.PublicNetConfig.EipConfig.EipStatus == "ENABLE"
		}

		triggers := make([]map[string]interface{}, 0, len(resp.Triggers))
		for _, trigger := range resp.Triggers {
			switch *trigger.Type {
			case SCF_TRIGGER_TYPE_TIMER:
				data := struct {
					Cron string `json:"cron"`
				}{}
				if err := json.Unmarshal([]byte(*trigger.TriggerDesc), &data); err != nil {
					log.Printf("[WARN]%s unmarshal timer trigger trigger_desc failed: %+v", logId, errors.WithStack(err))
					continue
				}
				*trigger.TriggerDesc = data.Cron
			}

			triggers = append(triggers, map[string]interface{}{
				"name":            trigger.TriggerName,
				"type":            trigger.Type,
				"trigger_desc":    trigger.TriggerDesc,
				"enable":          *trigger.Enable == 1,
				"create_time":     trigger.AddTime,
				"modify_time":     trigger.ModTime,
				"custom_argument": trigger.CustomArgument,
			})
		}
		m["trigger_info"] = triggers

		fnTags := make(map[string]string, len(resp.Tags))
		for _, tag := range resp.Tags {
			fnTags[*tag.Key] = *tag.Value
		}
		m["tags"] = fnTags
		m["async_run_enable"] = resp.AsyncRunEnable

		imageConfigs := make([]map[string]interface{}, 0, 1)
		if resp.ImageConfig != nil {
			imageConfigResp := resp.ImageConfig

			imageConfig := map[string]interface{}{
				"image_type": imageConfigResp.ImageType,
				"image_uri":  imageConfigResp.ImageUri,
			}
			if imageConfigResp.RegistryId != nil {
				imageConfig["registry_id"] = imageConfigResp.RegistryId
			}
			if imageConfigResp.EntryPoint != nil {
				imageConfig["entry_point"] = imageConfigResp.EntryPoint
			}
			if imageConfigResp.Command != nil {
				imageConfig["command"] = imageConfigResp.Command
			}
			if imageConfigResp.Args != nil {
				imageConfig["args"] = imageConfigResp.Args
			}
			if imageConfigResp.ContainerImageAccelerate != nil {
				imageConfig["container_image_accelerate"] = imageConfigResp.ContainerImageAccelerate
			}
			if imageConfigResp.ImagePort != nil {
				imageConfig["image_port"] = imageConfigResp.ImagePort
			}
			imageConfigs = append(imageConfigs, imageConfig)
		}
		m["image_config"] = imageConfigs

		if resp.DnsCache != nil {
			m["dns_cache"] = *resp.DnsCache == "TRUE"
		}

		intranetConfigs := make([]map[string]interface{}, 0, 1)
		if resp.IntranetConfig != nil {
			intranetConfigResp := resp.IntranetConfig

			intranetConfig := map[string]interface{}{
				"ip_fixed": intranetConfigResp.IpFixed,
			}
			if intranetConfigResp.IpAddress != nil {
				intranetConfig["ip_address"] = intranetConfigResp.IpAddress
			}
			intranetConfigs = append(intranetConfigs, intranetConfig)
		}
		m["intranet_config"] = intranetConfigs

		functions = append(functions, m)
	}

	if err := d.Set("functions", functions); err != nil {
		return err
	}
	d.SetId(helper.DataResourceIdsHash(ids))

	if output, ok := d.GetOk("result_output_file"); ok && output.(string) != "" {
		if err := tccommon.WriteToFile(output.(string), functions); err != nil {
			err = errors.WithStack(err)
			log.Printf("[CRITAL]%s output file[%s] fail, reason: %+v", logId, output.(string), err)
			return err
		}
	}

	return nil
}
