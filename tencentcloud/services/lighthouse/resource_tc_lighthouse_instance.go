package lighthouse

import (
	"context"
	"fmt"
	"log"
	"time"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	lighthouse "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/lighthouse/v20200324"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudLighthouseInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudLighthouseInstanceCreate,
		Read:   resourceTencentCloudLighthouseInstanceRead,
		Delete: resourceTencentCloudLighthouseInstanceDelete,
		Update: resourceTencentCloudLighthouseInstanceUpdate,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, i interface{}) ([]*schema.ResourceData, error) {
				_ = d.Set("is_update_bundle_id_auto_voucher", false)
				_ = d.Set("isolate_data_disk", true)

				return []*schema.ResourceData{d}, nil
			},
		},
		Schema: map[string]*schema.Schema{
			"bundle_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID Lighthouse 包。",
			},
			"is_update_bundle_id_auto_voucher": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "是否voucher 是 deducted automatically 当 update bundle ID. 取值范围：`true`: 表示automatic deduction 的 vouchers，`false`: does 不 automatically deduct vouchers. 默认值：`false`。",
			},
			"blueprint_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID Lighthouse 镜像。",
			},
			"period": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Subscription 周期 在 months. 有效值：1，2，3，4，5，6，7，8，9，10，11，12，24，36，48，60。",
			},
			"renew_flag": {
				Type:     schema.TypeString,
				Required: true,
				Description: "Auto-Renewal flag. 有效 值: NOTIFY_AND_AUTO_RENEW: notify upon expiration 和 renew automatically; NOTIFY_AND_MANUAL_RENEW: notify upon expiration 但 do 不 renew automatically. You need 到 manually renew DISABLE_NOTIFY_AND_AUTO_RENEW: neither notify upon expiration nor renew automatically." +
					"Default value: NOTIFY_AND_MANUAL_RENEW.",
			},
			"instance_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "display 名称 Lighthouse 实例。",
			},
			"zone": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "列表 availability zones. A random AZ 是 selected 通过 默认值。",
			},
			"dry_run": {
				Type:     schema.TypeBool,
				Optional: true,
				Description: "Whether 请求 是 dry run 仅." +
					"true: dry run only. The request will not create instance(s). A dry run can check whether all the required parameters are specified, whether the request format is right, whether the request exceeds service limits, and whether the specified CVMs are available. If the dry run fails, the corresponding error code will be returned.If the dry run succeeds, the RequestId will be returned." +
					"false (default value): send a normal request and create instance(s) if all the requirements are met.",
			},
			"client_token": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A 唯一 字符串 supplied 通过 客户端 到 ensure 该 请求 是 idempotent. Its 最大 长度 是 64 ASCII 字符. 如果 此 参数 是 不 指定， idem-potency 的 请求 不能 是 guaranteed。",
			},
			"login_configuration": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "Login 密码 的 实例. It 是 仅 可用 对于 Windows 实例. 如果 它 是 不 指定，它 表示 该 用户 choose 到 集合 login 密码 after 实例 creation。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"auto_generate_password": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "whether auto generate 密码 如果 false，need 集合 密码",
						},
						"password": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Login 密码",
						},
					},
				},
			},
			"permit_default_key_pair_login": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue([]string{"YES", "NO"}),
				Deprecated:   "It has been deprecated from version v1.81.8. Use `tencentcloud_lighthouse_key_pair_attachment` manage key pair.",
				Description:  "是否allow login 使用 默认值 键 pair. `YES`: allow login; `NO`: disable login. 默认值：`YES`。",
			},
			"isolate_data_disk": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "是否return mounted 数据 磁盘. `true`: 返回both 实例 和 mounted 数据 磁盘; `false`: 返回instance 和 无 longer 返回its mounted 数据 磁盘. 默认值：`true`。",
			},
			"containers": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Configuration 的 containers 到 create。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"container_image": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Container 镜像 地址",
						},
						"container_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Container 名称",
						},
						"envs": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "列表 环境 variables。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Environment variable 键",
									},
									"value": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Environment variable 值",
									},
								},
							},
						},
						"publish_ports": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "列表 mappings 的 容器 ports 和 主机 ports。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"host_port": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "主机端口",
									},
									"container_port": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Container 端口",
									},
									"ip": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "External IP. It 默认为 0.0.0.0。",
									},
									"protocol": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "协议 默认为 tcp. 有效值：tcp，udp 和 sctp。",
									},
								},
							},
						},
						"volumes": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "列表 容器 mount volumes。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"container_path": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Container 路径",
									},
									"host_path": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "主机 路径",
									},
								},
							},
						},
						"command": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "command 到 run。",
						},
					},
				},
			},
			"firewall_template_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Firewall 模板 ID 如果 此 参数 是 不 指定， 默认值 firewall 策略 是 使用。",
			},
			"public_addresses": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
				Description: "Public addresses。",
			},
			"private_addresses": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
				Description: "Private addresses。",
			},
		},
	}
}

func resourceTencentCloudLighthouseInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_lighthouse_instance.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request    = lighthouse.NewCreateInstancesRequest()
		instanceId string
	)

	if v, ok := d.GetOk("bundle_id"); ok {
		request.BundleId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("blueprint_id"); ok {
		request.BlueprintId = helper.String(v.(string))
	}

	instanceChargePrepaid := lighthouse.InstanceChargePrepaid{}
	if v, ok := d.GetOk("period"); ok {
		instanceChargePrepaid.Period = helper.IntInt64(v.(int))
	}
	if v, ok := d.GetOk("renew_flag"); ok {
		instanceChargePrepaid.RenewFlag = helper.String(v.(string))
	}
	request.InstanceChargePrepaid = &instanceChargePrepaid

	if v, ok := d.GetOk("instance_name"); ok {
		request.InstanceName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("zone"); ok {
		request.Zones = append(request.Zones, helper.String(v.(string)))
	}

	if v, ok := d.GetOk("dry_run"); ok {
		request.DryRun = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOk("client_token"); ok {
		request.ClientToken = helper.String(v.(string))
	}

	if loginConfigurationMap, ok := helper.InterfacesHeadMap(d, "login_configuration"); ok {
		loginConfiguration := lighthouse.LoginConfiguration{}
		if v, ok := loginConfigurationMap["auto_generate_password"]; ok {
			loginConfiguration.AutoGeneratePassword = helper.String(v.(string))
		}
		if v, ok := loginConfigurationMap["password"]; ok && v.(string) != "" {
			loginConfiguration.Password = helper.String(v.(string))
		}
		request.LoginConfiguration = &loginConfiguration
	}

	if v, ok := d.GetOk("containers"); ok {
		for _, container := range v.([]interface{}) {
			dockerContainerConfiguration := lighthouse.DockerContainerConfiguration{}
			containerMap := container.(map[string]interface{})
			if v, ok := containerMap["container_image"]; ok {
				dockerContainerConfiguration.ContainerImage = helper.String(v.(string))
			}
			if v, ok := containerMap["container_name"]; ok {
				dockerContainerConfiguration.ContainerName = helper.String(v.(string))
			}
			if v, ok := containerMap["envs"]; ok {
				for _, env := range v.([]interface{}) {
					containerEnv := lighthouse.ContainerEnv{}
					envMap := env.(map[string]interface{})
					if v, ok := envMap["key"]; ok {
						containerEnv.Key = helper.String(v.(string))
					}
					if v, ok := envMap["value"]; ok {
						containerEnv.Value = helper.String(v.(string))
					}
					dockerContainerConfiguration.Envs = append(dockerContainerConfiguration.Envs, &containerEnv)
				}
			}
			if v, ok := containerMap["publish_ports"]; ok {
				for _, publishPort := range v.([]interface{}) {
					dockerContainerPublishPort := lighthouse.DockerContainerPublishPort{}
					publishPortMap := publishPort.(map[string]interface{})
					if v, ok := publishPortMap["host_port"]; ok {
						dockerContainerPublishPort.HostPort = helper.IntInt64(v.(int))
					}
					if v, ok := publishPortMap["container_port"]; ok {
						dockerContainerPublishPort.ContainerPort = helper.IntInt64(v.(int))
					}
					if v, ok := publishPortMap["ip"]; ok {
						dockerContainerPublishPort.Ip = helper.String(v.(string))
					}
					if v, ok := publishPortMap["protocol"]; ok {
						dockerContainerPublishPort.Protocol = helper.String(v.(string))
					}
					dockerContainerConfiguration.PublishPorts = append(dockerContainerConfiguration.PublishPorts, &dockerContainerPublishPort)
				}
			}
			if v, ok := containerMap["volumes"]; ok {
				for _, volume := range v.([]interface{}) {
					dockerContainerVolume := lighthouse.DockerContainerVolume{}
					volumeMap := volume.(map[string]interface{})
					if v, ok := volumeMap["container_path"]; ok {
						dockerContainerVolume.ContainerPath = helper.String(v.(string))
					}
					if v, ok := volumeMap["host_path"]; ok {
						dockerContainerVolume.HostPath = helper.String(v.(string))
					}
					dockerContainerConfiguration.Volumes = append(dockerContainerConfiguration.Volumes, &dockerContainerVolume)
				}
			}
			if v, ok := containerMap["command"]; ok {
				dockerContainerConfiguration.Command = helper.String(v.(string))
			}
			request.Containers = append(request.Containers, &dockerContainerConfiguration)
		}
	}

	if v, ok := d.GetOk("firewall_template_id"); ok {
		request.FirewallTemplateId = helper.String(v.(string))
	}

	result, err := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseLighthouseClient().CreateInstances(request)

	if err != nil {
		log.Printf("[CRITAL]%s create lighthouse instance failed, reason:%+v", logId, err)
		return err
	}

	instanceId = *result.Response.InstanceIdSet[0]

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	lighthouseService := LightHouseService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	err = resource.Retry(5*tccommon.ReadRetryTimeout, func() *resource.RetryError {
		instance, errRet := lighthouseService.DescribeLighthouseInstanceById(ctx, instanceId)
		if errRet != nil {
			return tccommon.RetryError(errRet, tccommon.InternalError)
		}
		if instance != nil && (*instance.InstanceState == "RUNNING") {
			return nil
		}
		if instance == nil || instance.InstanceState == nil {
			return resource.RetryableError(fmt.Errorf("lighthouse instance creating..."))
		}
		return resource.RetryableError(fmt.Errorf("lighthouse instance status is %s, retry...", *instance.InstanceState))
	})
	if err != nil {
		return err
	}

	d.SetId(instanceId)

	return resourceTencentCloudLighthouseInstanceRead(d, meta)
}

func resourceTencentCloudLighthouseInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_lighthouse_instance.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	lighthouseService := LightHouseService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	id := d.Id()

	instance, err := lighthouseService.DescribeLighthouseInstanceById(ctx, id)

	if err != nil {
		return err
	}

	if instance == nil {
		d.SetId("")
		return fmt.Errorf("resource `lighthouse instance` %s does not exist", id)
	}

	if instance.BundleId != nil {
		_ = d.Set("bundle_id", instance.BundleId)
	}

	if instance.BlueprintId != nil {
		_ = d.Set("blueprint_id", instance.BlueprintId)
	}

	if instance.InstanceChargeType != nil {
		_ = d.Set("renew_flag", instance.RenewFlag)
	}

	if instance.InstanceName != nil {
		_ = d.Set("instance_name", instance.InstanceName)
	}

	if instance.Zone != nil {
		_ = d.Set("zone", instance.Zone)
	}

	if len(instance.PublicAddresses) > 0 {
		_ = d.Set("public_addresses", instance.PublicAddresses)
	}

	if len(instance.PrivateAddresses) > 0 {
		_ = d.Set("private_addresses", instance.PrivateAddresses)
	}

	return nil
}

func resourceTencentCloudLighthouseInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_lighthouse_instance.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request = lighthouse.NewModifyInstancesAttributeRequest()
	)

	immutableArgs := []string{"firewall_template_id"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	id := d.Id()

	request.InstanceIds = append(request.InstanceIds, helper.String(id))

	if d.HasChange("instance_name") {
		if v, ok := d.GetOk("instance_name"); ok {
			request.InstanceName = helper.String(v.(string))
		}
		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseLighthouseClient().ModifyInstancesAttribute(request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
					logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			return nil
		})

		if err != nil {
			return err
		}

		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service := LightHouseService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		err = resource.Retry(5*tccommon.ReadRetryTimeout, func() *resource.RetryError {
			instance, errRet := service.DescribeLighthouseInstanceById(ctx, id)
			if errRet != nil {
				return tccommon.RetryError(errRet, tccommon.InternalError)
			}
			if instance.LatestOperationState == nil {
				return resource.RetryableError(fmt.Errorf("waiting for instance operation update"))
			}
			if *instance.LatestOperationState == "OPERATING" {
				return resource.RetryableError(fmt.Errorf("waiting for instance %s operation", id))
			}
			if *instance.LatestOperationState == "FAILED" {
				return resource.NonRetryableError(fmt.Errorf("failed operation"))
			}
			return nil
		})

		if err != nil {
			return err
		}
	}

	if d.HasChange("bundle_id") {
		_, new := d.GetChange("bundle_id")
		request := lighthouse.NewModifyInstancesBundleRequest()
		request.InstanceIds = helper.StringsStringsPoint([]string{id})
		request.BundleId = helper.String(new.(string))
		autoVoucher := d.Get("is_update_bundle_id_auto_voucher").(bool)
		request.AutoVoucher = &autoVoucher
		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseLighthouseClient().ModifyInstancesBundle(request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update lighthouse instanceModifyBundle failed, reason:%+v", logId, err)
			return err
		}

		service := LightHouseService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

		conf := tccommon.BuildStateChangeConf([]string{}, []string{"SUCCESS"}, 20*tccommon.ReadRetryTimeout, time.Second, service.LighthouseInstanceStateRefreshFunc(id, []string{}))

		if _, e := conf.WaitForState(); e != nil {
			return e
		}
	}

	if d.HasChange("blueprint_id") {
		_, new := d.GetChange("blueprint_id")
		request := lighthouse.NewResetInstanceRequest()
		request.InstanceId = helper.String(id)
		request.BlueprintId = helper.String(new.(string))

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseLighthouseClient().ResetInstance(request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s operate lighthouse resetInstance failed, reason:%+v", logId, err)
			return err
		}

		service := LightHouseService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

		conf := tccommon.BuildStateChangeConf([]string{}, []string{"SUCCESS"}, 20*tccommon.ReadRetryTimeout, time.Second, service.LighthouseInstanceStateRefreshFunc(id, []string{}))

		if _, e := conf.WaitForState(); e != nil {
			return e
		}
	}

	if d.HasChange("period") {
		old, _ := d.GetChange("period")
		_ = d.Set("period", old)
		return fmt.Errorf("`period` do not support change now.")
	}

	if d.HasChange("renew_flag") {
		_, new := d.GetChange("renew_flag")
		request := lighthouse.NewModifyInstancesRenewFlagRequest()
		request.InstanceIds = helper.StringsStringsPoint([]string{id})
		request.RenewFlag = helper.String(new.(string))
		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseLighthouseClient().ModifyInstancesRenewFlag(request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s operate lighthouse modifyInstanceRenewFlag failed, reason:%+v", logId, err)
			return err
		}
	}

	if d.HasChange("zone") {
		old, _ := d.GetChange("zone")
		_ = d.Set("zone", old)
		return fmt.Errorf("`zone` do not support change now.")
	}

	if d.HasChange("dry_run") {
		old, _ := d.GetChange("dry_run")
		_ = d.Set("dry_run", old)
		return fmt.Errorf("`dry_run` do not support change now.")
	}

	if d.HasChange("client_token") {
		old, _ := d.GetChange("client_token")
		_ = d.Set("client_token", old)
		return fmt.Errorf("`client_token` do not support change now.")
	}

	if d.HasChange("login_configuration.0.auto_generate_password") {
		old, _ := d.GetChange("login_configuration")
		_ = d.Set("login_configuration", old)
		return fmt.Errorf("`auto_generate_password` do not support change now.")
	}
	if d.HasChange("login_configuration.0.password") {
		_, new := d.GetChange("login_configuration")
		var newLoginConfiguration map[string]interface{}
		if len(new.([]interface{})) > 0 {
			newLoginConfiguration = new.([]interface{})[0].(map[string]interface{})
		}
		newPassword := newLoginConfiguration["password"].(string)
		request := lighthouse.NewResetInstancesPasswordRequest()
		request.InstanceIds = helper.StringsStringsPoint([]string{id})
		request.Password = helper.String(newPassword)
		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseLighthouseClient().ResetInstancesPassword(request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s operate lighthouse resetInstancesPassword failed, reason:%+v", logId, err)
			return err
		}
		service := LightHouseService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

		conf := tccommon.BuildStateChangeConf([]string{}, []string{"SUCCESS"}, 20*tccommon.ReadRetryTimeout, time.Second, service.LighthouseInstanceStateRefreshFunc(id, []string{}))

		if _, e := conf.WaitForState(); e != nil {
			return e
		}

		_ = d.Set("login_configuration", new)
	}

	if d.HasChange("containers") {
		old, _ := d.GetChange("containers")
		_ = d.Set("containers", old)
		return fmt.Errorf("`containers` do not support change now.")
	}

	return resourceTencentCloudLighthouseInstanceRead(d, meta)
}

func resourceTencentCloudLighthouseInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_lighthouse_instance.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	service := LightHouseService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	id := d.Id()
	isolateDataDisk := d.Get("isolate_data_disk").(bool)

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		if err := service.IsolateLighthouseInstanceById(ctx, id, isolateDataDisk); err != nil {
			return tccommon.RetryError(err)
		}
		return nil
	})
	if err != nil {
		return err
	}

	err = resource.Retry(5*tccommon.ReadRetryTimeout, func() *resource.RetryError {
		instance, errRet := service.DescribeLighthouseInstanceById(ctx, id)
		if errRet != nil {
			return tccommon.RetryError(errRet, tccommon.InternalError)
		}
		if instance.LatestOperationState == nil {
			return resource.RetryableError(fmt.Errorf("waiting for instance operation update"))
		}
		if *instance.LatestOperationState == "FAILED" {
			return resource.NonRetryableError(fmt.Errorf("failed operation"))
		}
		if *instance.InstanceState == "SHUTDOWN" && *instance.LatestOperationState != "OPERATING" {
			return nil
		}
		return resource.RetryableError(fmt.Errorf("instance status is %s, retry...", *instance.InstanceState))
	})

	if err != nil {
		return err
	}

	if err := service.DeleteLighthouseInstanceById(ctx, id); err != nil {
		return err
	}
	return nil
}
