package dlc

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dlcv20210125 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dlc/v20210125"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudDlcStandardEngineResourceGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDlcStandardEngineResourceGroupCreate,
		Read:   resourceTencentCloudDlcStandardEngineResourceGroupRead,
		Update: resourceTencentCloudDlcStandardEngineResourceGroupUpdate,
		Delete: resourceTencentCloudDlcStandardEngineResourceGroupDelete,
		Schema: map[string]*schema.Schema{
			"engine_resource_group_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Standard 引擎 资源 组名称",
			},

			"data_engine_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Standard 引擎 名称",
			},

			"auto_launch": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Automatic start (任务 submission automatically pulls up 资源 组) 0-automatic start，1-不 automatic start。",
			},

			"auto_pause": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Automatically suspend 资源 groups. 0 - Automatically suspend，1 - Not automatically suspend。",
			},

			"driver_cu_spec": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Driver CU specifications: Currently 支持: small (默认值，1 CU)，medium (2 CU)，large (4 CU)，xlarge (8 CU). Memory CUs 是 CPUs 使用 ratio 的 1:8，m.small (1 CU 内存)，m.medium (2 CU 内存)，m.large (4 CU 内存)，和 m.xlarge (8 CU 内存)。",
			},

			"executor_cu_spec": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Executor CU specifications: Currently 支持: small (默认值，1 CU)，medium (2 CU)，large (4 CU)，xlarge (8 CU). Memory CUs 是 CPUs 使用 ratio 的 1:8，m.small (1 CU 内存)，m.medium (2 CU 内存)，m.large (4 CU 内存)，和 m.xlarge (8 CU 内存)。",
			},

			"min_executor_nums": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "最小executors。",
			},

			"max_executor_nums": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "最大executors。",
			},

			"auto_pause_time": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Automatic suspension 时间，在 minutes，使用 值 范围 的 1-999 (after 无 tasks have reached AutoPauseTime， 资源 组 将 automatically suspend)。",
			},

			"static_config_pairs": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Description: "Static 参数 的 资源 组，其中 require restarting 资源 组 到 take effect。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"config_item": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Configuration items。",
						},
						"config_value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Configuration 值。",
						},
					},
				},
			},

			"dynamic_config_pairs": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Description: "Dynamic 参数 的 资源 组，effective 在 next 任务。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"config_item": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Configuration items。",
						},
						"config_value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Configuration 值。",
						},
					},
				},
			},

			"max_concurrency": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "数量 concurrent tasks 是 5 通过 默认值。",
			},

			"network_config_names": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Network 配置 名称",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"public_domain": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Customized mirror 域名 名称",
			},

			"registry_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Custom 镜像 实例 ID。",
			},

			"frame_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "框架 类型 AI 类型 资源 组，machine-learning，python，spark-ml，如果未填写 在， 默认为 machine-learning。",
			},

			"image_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Image 类型，build-在: built-在，自定义: 自定义，如果未填写 在， 默认为 build-在。",
			},

			"image_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Image 名称 \nExample 值: 镜像-xxx. 如果 使用 built-在 镜像 (ImageType 是 built-在)， ImageName 对于 different frameworks 是: machine-learning: pytorch-v2.5.1，scikit-learn-v1.6.0，tensorflow-v2.18.0，python: python-v3.10，spark-m: Standard-S 1.1。",
			},

			"image_version": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Image ID。",
			},

			"size": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "AI 资源 组 是 有效，和 upper 限制 的 可用 resources 在 资源 组 必须 是 less 比 upper 限制 的 引擎 resources。",
			},

			"resource_group_scene": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Resource 组 scenario。",
			},

			"region_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Custom 镜像 location。",
			},

			"python_cu_spec": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "资源 限制 对于 Python stand-alone 节点 在 Python 资源 组 必须 是 smaller 比 资源 限制 对于 资源 组. Small: 1cu Medium: 2cu Large: 4cu Xlarge: 8cu 4xlarge: 16cu 8xlarge: 32cu 16xlarge: 64cu. 如果 资源类型 是 high 内存，add m before 类型",
			},

			"spark_spec_mode": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Only SQL 资源 组 资源 配置 模式，fast: fast 模式，自定义: 自定义 模式",
			},

			"spark_size": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Only SQL 资源 组 资源 限制，仅 用于the express 模块",
			},

			"running_state": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "state 的 资源 组. true: launch standard 引擎 资源 组; false: pause standard 引擎 资源 组. 默认为 true。",
			},

			// computed
			"engine_resource_group_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Standard 引擎 资源 组 ID",
			},
		},
	}
}

func resourceTencentCloudDlcStandardEngineResourceGroupCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dlc_standard_engine_resource_group.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId                   = tccommon.GetLogId(tccommon.ContextNil)
		ctx                     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request                 = dlcv20210125.NewCreateStandardEngineResourceGroupRequest()
		engineResourceGroupName string
	)

	if v, ok := d.GetOk("engine_resource_group_name"); ok {
		request.EngineResourceGroupName = helper.String(v.(string))
		engineResourceGroupName = v.(string)
	}

	if v, ok := d.GetOk("data_engine_name"); ok {
		request.DataEngineName = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("auto_launch"); ok {
		request.AutoLaunch = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("auto_pause"); ok {
		request.AutoPause = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("driver_cu_spec"); ok {
		request.DriverCuSpec = helper.String(v.(string))
	}

	if v, ok := d.GetOk("executor_cu_spec"); ok {
		request.ExecutorCuSpec = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("min_executor_nums"); ok {
		request.MinExecutorNums = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("max_executor_nums"); ok {
		request.MaxExecutorNums = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("auto_pause_time"); ok {
		request.AutoPauseTime = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("static_config_pairs"); ok {
		for _, item := range v.([]interface{}) {
			staticConfigPairsMap := item.(map[string]interface{})
			engineResourceGroupConfigPair := dlcv20210125.EngineResourceGroupConfigPair{}
			if v, ok := staticConfigPairsMap["config_item"].(string); ok && v != "" {
				engineResourceGroupConfigPair.ConfigItem = helper.String(v)
			}

			if v, ok := staticConfigPairsMap["config_value"].(string); ok && v != "" {
				engineResourceGroupConfigPair.ConfigValue = helper.String(v)
			}

			request.StaticConfigPairs = append(request.StaticConfigPairs, &engineResourceGroupConfigPair)
		}
	}

	if v, ok := d.GetOk("dynamic_config_pairs"); ok {
		for _, item := range v.([]interface{}) {
			dynamicConfigPairsMap := item.(map[string]interface{})
			engineResourceGroupConfigPair := dlcv20210125.EngineResourceGroupConfigPair{}
			if v, ok := dynamicConfigPairsMap["config_item"].(string); ok && v != "" {
				engineResourceGroupConfigPair.ConfigItem = helper.String(v)
			}

			if v, ok := dynamicConfigPairsMap["config_value"].(string); ok && v != "" {
				engineResourceGroupConfigPair.ConfigValue = helper.String(v)
			}

			request.DynamicConfigPairs = append(request.DynamicConfigPairs, &engineResourceGroupConfigPair)
		}
	}

	if v, ok := d.GetOkExists("max_concurrency"); ok {
		request.MaxConcurrency = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("network_config_names"); ok {
		networkConfigNamesSet := v.(*schema.Set).List()
		for i := range networkConfigNamesSet {
			if networkConfigNames, ok := networkConfigNamesSet[i].(string); ok && networkConfigNames != "" {
				request.NetworkConfigNames = append(request.NetworkConfigNames, helper.String(networkConfigNames))
			}
		}
	}

	if v, ok := d.GetOk("public_domain"); ok {
		request.PublicDomain = helper.String(v.(string))
	}

	if v, ok := d.GetOk("registry_id"); ok {
		request.RegistryId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("frame_type"); ok {
		request.FrameType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("image_type"); ok {
		request.ImageType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("image_name"); ok {
		request.ImageName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("image_version"); ok {
		request.ImageVersion = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("size"); ok {
		request.Size = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("resource_group_scene"); ok {
		request.ResourceGroupScene = helper.String(v.(string))
	}

	if v, ok := d.GetOk("region_name"); ok {
		request.RegionName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("python_cu_spec"); ok {
		request.PythonCuSpec = helper.String(v.(string))
	}

	if v, ok := d.GetOk("spark_spec_mode"); ok {
		request.SparkSpecMode = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("spark_size"); ok {
		request.SparkSize = helper.IntInt64(v.(int))
	}

	request.IsLaunchNow = helper.IntInt64(0)
	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDlcClient().CreateStandardEngineResourceGroupWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create dlc standard engine resource group failed, Response is nil."))
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create dlc standard engine resource group failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	d.SetId(engineResourceGroupName)

	// wait
	waitErr := resource.Retry(tccommon.WriteRetryTimeout*4, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDlcClient().DescribeStandardEngineResourceGroupsWithContext(ctx, &dlcv20210125.DescribeStandardEngineResourceGroupsRequest{
			Filters: []*dlcv20210125.Filter{
				{
					Name:   helper.String("engine-resource-group-name-unique"),
					Values: helper.Strings([]string{engineResourceGroupName}),
				},
			},
		})
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe dlc standard engine resource groups failed, Response is nil."))
		}

		if result.Response.UserEngineResourceGroupInfos == nil || len(result.Response.UserEngineResourceGroupInfos) == 0 {
			return resource.NonRetryableError(fmt.Errorf("UserEngineResourceGroupInfos is nil."))
		}

		if len(result.Response.UserEngineResourceGroupInfos) != 1 {
			return resource.NonRetryableError(fmt.Errorf("UserEngineResourceGroupInfos is not 1."))
		}

		state := result.Response.UserEngineResourceGroupInfos[0].ResourceGroupState
		if state != nil {
			if *state == 2 {
				return nil
			}
		} else {
			return resource.NonRetryableError(fmt.Errorf("ResourceGroupState is nil."))
		}

		return resource.RetryableError(fmt.Errorf("UserEngineResourceGroupInfos is not ready, state:%d", *state))
	})

	if waitErr != nil {
		log.Printf("[CRITAL]%s wait for dlc standard engine resource group failed, reason:%+v", logId, waitErr)
		return waitErr
	}

	if v, ok := d.GetOkExists("running_state"); ok {
		// pause resource group
		if !v.(bool) {
			request := dlcv20210125.NewPauseStandardEngineResourceGroupsRequest()
			request.EngineResourceGroupNames = helper.Strings([]string{engineResourceGroupName})
			reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
				result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDlcClient().PauseStandardEngineResourceGroupsWithContext(ctx, request)
				if e != nil {
					return tccommon.RetryError(e)
				} else {
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
				}

				if result == nil || result.Response == nil {
					return resource.NonRetryableError(fmt.Errorf("Pause dlc standard engine resource group failed, Response is nil."))
				}

				return nil
			})

			if reqErr != nil {
				log.Printf("[CRITAL]%s pause dlc standard engine resource group failed, reason:%+v", logId, reqErr)
				return reqErr
			}

			// wait
			waitErr := resource.Retry(tccommon.WriteRetryTimeout*4, func() *resource.RetryError {
				result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDlcClient().DescribeStandardEngineResourceGroupsWithContext(ctx, &dlcv20210125.DescribeStandardEngineResourceGroupsRequest{
					Filters: []*dlcv20210125.Filter{
						{
							Name:   helper.String("engine-resource-group-name-unique"),
							Values: helper.Strings([]string{engineResourceGroupName}),
						},
					},
				})
				if e != nil {
					return tccommon.RetryError(e)
				} else {
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
				}

				if result == nil || result.Response == nil {
					return resource.NonRetryableError(fmt.Errorf("Describe dlc standard engine resource groups failed, Response is nil."))
				}

				if result.Response.UserEngineResourceGroupInfos == nil || len(result.Response.UserEngineResourceGroupInfos) == 0 {
					return resource.NonRetryableError(fmt.Errorf("UserEngineResourceGroupInfos is nil."))
				}

				if len(result.Response.UserEngineResourceGroupInfos) != 1 {
					return resource.NonRetryableError(fmt.Errorf("UserEngineResourceGroupInfos is not 1."))
				}

				state := result.Response.UserEngineResourceGroupInfos[0].ResourceGroupState
				if state != nil {
					if *state == 3 {
						return nil
					}
				} else {
					return resource.NonRetryableError(fmt.Errorf("ResourceGroupState is nil."))
				}

				return resource.RetryableError(fmt.Errorf("UserEngineResourceGroupInfos is not ready, state:%d", *state))
			})

			if waitErr != nil {
				log.Printf("[CRITAL]%s wait for dlc standard engine resource group failed, reason:%+v", logId, waitErr)
				return waitErr
			}
		}
	}

	return resourceTencentCloudDlcStandardEngineResourceGroupRead(d, meta)
}

func resourceTencentCloudDlcStandardEngineResourceGroupRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dlc_standard_engine_resource_group.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId                   = tccommon.GetLogId(tccommon.ContextNil)
		ctx                     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service                 = DlcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		engineResourceGroupName = d.Id()
	)

	respData, err := service.DescribeDlcStandardEngineResourceGroupById(ctx, engineResourceGroupName)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[WARN]%s resource `tencentcloud_dlc_standard_engine_resource_group` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	if respData.EngineResourceGroupName != nil {
		_ = d.Set("engine_resource_group_name", respData.EngineResourceGroupName)
	}

	if respData.DataEngineName != nil {
		_ = d.Set("data_engine_name", respData.DataEngineName)
	}

	if respData.AutoLaunch != nil {
		_ = d.Set("auto_launch", respData.AutoLaunch)
	}

	if respData.AutoPause != nil {
		_ = d.Set("auto_pause", respData.AutoPause)
	}

	if respData.DriverCuSpec != nil {
		_ = d.Set("driver_cu_spec", respData.DriverCuSpec)
	}

	if respData.ExecutorCuSpec != nil {
		_ = d.Set("executor_cu_spec", respData.ExecutorCuSpec)
	}

	if respData.MinExecutorNums != nil {
		_ = d.Set("min_executor_nums", respData.MinExecutorNums)
	}

	if respData.MaxExecutorNums != nil {
		_ = d.Set("max_executor_nums", respData.MaxExecutorNums)
	}

	if respData.AutoPauseTime != nil {
		_ = d.Set("auto_pause_time", respData.AutoPauseTime)
	}

	if respData.MaxConcurrency != nil {
		_ = d.Set("max_concurrency", respData.MaxConcurrency)
	}

	if respData.NetworkConfigNames != nil {
		_ = d.Set("network_config_names", respData.NetworkConfigNames)
	}

	if respData.PublicDomain != nil {
		_ = d.Set("public_domain", respData.PublicDomain)
	}

	if respData.RegistryId != nil {
		_ = d.Set("registry_id", respData.RegistryId)
	}

	if respData.FrameType != nil {
		_ = d.Set("frame_type", respData.FrameType)
	}

	if respData.ImageType != nil {
		_ = d.Set("image_type", respData.ImageType)
	}

	if respData.ImageName != nil {
		_ = d.Set("image_name", respData.ImageName)
	}

	if respData.ImageVersion != nil {
		_ = d.Set("image_version", respData.ImageVersion)
	}

	if respData.Size != nil {
		_ = d.Set("size", respData.Size)
	}

	if respData.ResourceGroupScene != nil {
		_ = d.Set("resource_group_scene", respData.ResourceGroupScene)
	}

	if respData.RegionName != nil {
		_ = d.Set("region_name", respData.RegionName)
	}

	if respData.PythonCuSpec != nil {
		_ = d.Set("python_cu_spec", respData.PythonCuSpec)
	}

	if respData.SparkSpecMode != nil {
		_ = d.Set("spark_spec_mode", respData.SparkSpecMode)
	}

	if respData.SparkSize != nil {
		_ = d.Set("spark_size", respData.SparkSize)
	}

	if respData.EngineResourceGroupId != nil {
		_ = d.Set("engine_resource_group_id", respData.EngineResourceGroupId)
	}

	if respData.ResourceGroupState != nil {
		if *respData.ResourceGroupState == 2 {
			_ = d.Set("running_state", true)
		} else if *respData.ResourceGroupState == 3 {
			_ = d.Set("running_state", false)
		}
	}

	return nil
}

func resourceTencentCloudDlcStandardEngineResourceGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dlc_standard_engine_resource_group.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId                   = tccommon.GetLogId(tccommon.ContextNil)
		ctx                     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		engineResourceGroupName = d.Id()
	)

	immutableArgs := []string{"engine_resource_group_name", "data_engine_name", "static_config_pairs", "dynamic_config_pairs", "resource_group_scene"}
	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	needChange := false
	mutableArgs := []string{"driver_cu_spec", "executor_cu_spec", "min_executor_nums", "max_executor_nums", "size", "image_type", "image_name", "image_version", "frame_type", "public_domain", "registry_id", "region_name", "python_cu_spec", "spark_spec_mode", "spark_size"}
	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}

	if needChange {
		request := dlcv20210125.NewUpdateStandardEngineResourceGroupResourceInfoRequest()
		if v, ok := d.GetOk("driver_cu_spec"); ok {
			request.DriverCuSpec = helper.String(v.(string))
		}

		if v, ok := d.GetOk("executor_cu_spec"); ok {
			request.ExecutorCuSpec = helper.String(v.(string))
		}

		if v, ok := d.GetOkExists("min_executor_nums"); ok {
			request.MinExecutorNums = helper.IntInt64(v.(int))
		}

		if v, ok := d.GetOkExists("max_executor_nums"); ok {
			request.MaxExecutorNums = helper.IntInt64(v.(int))
		}

		if v, ok := d.GetOkExists("is_effective_now"); ok {
			request.IsEffectiveNow = helper.IntInt64(v.(int))
		}

		if v, ok := d.GetOkExists("size"); ok {
			request.Size = helper.IntInt64(v.(int))
		}

		if v, ok := d.GetOk("image_type"); ok {
			request.ImageType = helper.String(v.(string))
		}

		if v, ok := d.GetOk("image_name"); ok {
			request.ImageName = helper.String(v.(string))
		}

		if v, ok := d.GetOk("image_version"); ok {
			request.ImageVersion = helper.String(v.(string))
		}

		if v, ok := d.GetOk("frame_type"); ok {
			request.FrameType = helper.String(v.(string))
		}

		if v, ok := d.GetOk("public_domain"); ok {
			request.PublicDomain = helper.String(v.(string))
		}

		if v, ok := d.GetOk("registry_id"); ok {
			request.RegistryId = helper.String(v.(string))
		}

		if v, ok := d.GetOk("region_name"); ok {
			request.RegionName = helper.String(v.(string))
		}

		if v, ok := d.GetOk("python_cu_spec"); ok {
			request.PythonCuSpec = helper.String(v.(string))
		}

		if v, ok := d.GetOk("spark_spec_mode"); ok {
			request.SparkSpecMode = helper.String(v.(string))
		}

		if v, ok := d.GetOkExists("spark_size"); ok {
			request.SparkSize = helper.IntInt64(v.(int))
		}

		request.EngineResourceGroupName = &engineResourceGroupName
		request.IsEffectiveNow = helper.Int64(0)
		reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDlcClient().UpdateStandardEngineResourceGroupResourceInfoWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			return nil
		})

		if reqErr != nil {
			log.Printf("[CRITAL]%s update dlc standard engine resource group failed, reason:%+v", logId, reqErr)
			return reqErr
		}

		// wait
		waitErr := resource.Retry(tccommon.WriteRetryTimeout*4, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDlcClient().DescribeStandardEngineResourceGroupsWithContext(ctx, &dlcv20210125.DescribeStandardEngineResourceGroupsRequest{
				Filters: []*dlcv20210125.Filter{
					{
						Name:   helper.String("engine-resource-group-name-unique"),
						Values: helper.Strings([]string{engineResourceGroupName}),
					},
				},
			})
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			if result == nil || result.Response == nil {
				return resource.NonRetryableError(fmt.Errorf("Describe dlc standard engine resource groups failed, Response is nil."))
			}

			if result.Response.UserEngineResourceGroupInfos == nil || len(result.Response.UserEngineResourceGroupInfos) == 0 {
				return resource.NonRetryableError(fmt.Errorf("UserEngineResourceGroupInfos is nil."))
			}

			if len(result.Response.UserEngineResourceGroupInfos) != 1 {
				return resource.NonRetryableError(fmt.Errorf("UserEngineResourceGroupInfos is not 1."))
			}

			state := result.Response.UserEngineResourceGroupInfos[0].ResourceGroupState
			if state != nil {
				if *state == 2 {
					return nil
				}
			} else {
				return resource.NonRetryableError(fmt.Errorf("ResourceGroupState is nil."))
			}

			return resource.RetryableError(fmt.Errorf("UserEngineResourceGroupInfos is not ready, state:%d", *state))
		})

		if waitErr != nil {
			log.Printf("[CRITAL]%s wait for dlc standard engine resource group failed, reason:%+v", logId, waitErr)
			return waitErr
		}
	}

	if d.HasChange("network_config_names") {
		request := dlcv20210125.NewUpdateEngineResourceGroupNetworkConfigInfoRequest()
		if v, ok := d.GetOk("network_config_names"); ok {
			networkConfigNamesSet := v.(*schema.Set).List()
			for i := range networkConfigNamesSet {
				if networkConfigNames, ok := networkConfigNamesSet[i].(string); ok && networkConfigNames != "" {
					request.NetworkConfigNames = append(request.NetworkConfigNames, helper.String(networkConfigNames))
				}
			}
		}

		engineResourceGroupId := d.Get("engine_resource_group_id").(string)
		request.EngineResourceGroupId = &engineResourceGroupId
		request.IsEffectiveNow = helper.Int64(0)
		reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDlcClient().UpdateEngineResourceGroupNetworkConfigInfoWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			return nil
		})

		if reqErr != nil {
			log.Printf("[CRITAL]%s update dlc standard engine resource group network config names failed, reason:%+v", logId, reqErr)
			return reqErr
		}

		// wait
		waitErr := resource.Retry(tccommon.WriteRetryTimeout*4, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDlcClient().DescribeStandardEngineResourceGroupsWithContext(ctx, &dlcv20210125.DescribeStandardEngineResourceGroupsRequest{
				Filters: []*dlcv20210125.Filter{
					{
						Name:   helper.String("engine-resource-group-name-unique"),
						Values: helper.Strings([]string{engineResourceGroupName}),
					},
				},
			})
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			if result == nil || result.Response == nil {
				return resource.NonRetryableError(fmt.Errorf("Describe dlc standard engine resource groups failed, Response is nil."))
			}

			if result.Response.UserEngineResourceGroupInfos == nil || len(result.Response.UserEngineResourceGroupInfos) == 0 {
				return resource.NonRetryableError(fmt.Errorf("UserEngineResourceGroupInfos is nil."))
			}

			if len(result.Response.UserEngineResourceGroupInfos) != 1 {
				return resource.NonRetryableError(fmt.Errorf("UserEngineResourceGroupInfos is not 1."))
			}

			state := result.Response.UserEngineResourceGroupInfos[0].ResourceGroupState
			if state != nil {
				if *state == 2 {
					return nil
				}
			} else {
				return resource.NonRetryableError(fmt.Errorf("ResourceGroupState is nil."))
			}

			return resource.RetryableError(fmt.Errorf("UserEngineResourceGroupInfos is not ready, state:%d", *state))
		})

		if waitErr != nil {
			log.Printf("[CRITAL]%s wait for dlc standard engine resource group failed, reason:%+v", logId, waitErr)
			return waitErr
		}
	}

	if d.HasChange("auto_pause") || d.HasChange("auto_launch") || d.HasChange("auto_pause_time") || d.HasChange("max_concurrency") {
		var (
			autoPause  int
			autoLaunch int
		)

		request := dlcv20210125.NewUpdateStandardEngineResourceGroupBaseInfoRequest()
		if v, ok := d.GetOkExists("auto_pause"); ok {
			autoPause = v.(int)
		}

		if v, ok := d.GetOkExists("auto_launch"); ok {
			autoLaunch = v.(int)
		}

		if v, ok := d.GetOkExists("auto_pause_time"); ok {
			request.AutoPauseTime = helper.IntInt64(v.(int))
		}

		if v, ok := d.GetOkExists("max_concurrency"); ok {
			request.MaxConcurrency = helper.IntInt64(v.(int))
		}

		request.AutoPause = helper.IntInt64(autoPause)
		request.AutoLaunch = helper.IntInt64(autoLaunch)
		request.EngineResourceGroupName = &engineResourceGroupName
		reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDlcClient().UpdateStandardEngineResourceGroupBaseInfoWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			return nil
		})

		if reqErr != nil {
			log.Printf("[CRITAL]%s update dlc standard engine resource group base info failed, reason:%+v", logId, reqErr)
			return reqErr
		}

		// wait
		waitErr := resource.Retry(tccommon.WriteRetryTimeout*4, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDlcClient().DescribeStandardEngineResourceGroupsWithContext(ctx, &dlcv20210125.DescribeStandardEngineResourceGroupsRequest{
				Filters: []*dlcv20210125.Filter{
					{
						Name:   helper.String("engine-resource-group-name-unique"),
						Values: helper.Strings([]string{engineResourceGroupName}),
					},
				},
			})
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			if result == nil || result.Response == nil {
				return resource.NonRetryableError(fmt.Errorf("Describe dlc standard engine resource groups failed, Response is nil."))
			}

			if result.Response.UserEngineResourceGroupInfos == nil || len(result.Response.UserEngineResourceGroupInfos) == 0 {
				return resource.NonRetryableError(fmt.Errorf("UserEngineResourceGroupInfos is nil."))
			}

			if len(result.Response.UserEngineResourceGroupInfos) != 1 {
				return resource.NonRetryableError(fmt.Errorf("UserEngineResourceGroupInfos is not 1."))
			}

			state := result.Response.UserEngineResourceGroupInfos[0].ResourceGroupState
			if state != nil {
				if *state == 2 {
					return nil
				}
			} else {
				return resource.NonRetryableError(fmt.Errorf("ResourceGroupState is nil."))
			}

			return resource.RetryableError(fmt.Errorf("UserEngineResourceGroupInfos is not ready, state:%d", *state))
		})

		if waitErr != nil {
			log.Printf("[CRITAL]%s wait for dlc standard engine resource group failed, reason:%+v", logId, waitErr)
			return waitErr
		}
	}

	if d.HasChange("running_state") {
		if v, ok := d.GetOkExists("running_state"); ok {
			if v.(bool) {
				request := dlcv20210125.NewLaunchStandardEngineResourceGroupsRequest()
				request.EngineResourceGroupNames = helper.Strings([]string{engineResourceGroupName})
				reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
					result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDlcClient().LaunchStandardEngineResourceGroupsWithContext(ctx, request)
					if e != nil {
						return tccommon.RetryError(e)
					} else {
						log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
					}

					if result == nil || result.Response == nil {
						return resource.NonRetryableError(fmt.Errorf("Launch dlc standard engine resource group failed, Response is nil."))
					}

					return nil
				})

				if reqErr != nil {
					log.Printf("[CRITAL]%s launch dlc standard engine resource group failed, reason:%+v", logId, reqErr)
					return reqErr
				}

				// wait
				waitErr := resource.Retry(tccommon.WriteRetryTimeout*4, func() *resource.RetryError {
					result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDlcClient().DescribeStandardEngineResourceGroupsWithContext(ctx, &dlcv20210125.DescribeStandardEngineResourceGroupsRequest{
						Filters: []*dlcv20210125.Filter{
							{
								Name:   helper.String("engine-resource-group-name-unique"),
								Values: helper.Strings([]string{engineResourceGroupName}),
							},
						},
					})
					if e != nil {
						return tccommon.RetryError(e)
					} else {
						log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
					}

					if result == nil || result.Response == nil {
						return resource.NonRetryableError(fmt.Errorf("Describe dlc standard engine resource groups failed, Response is nil."))
					}

					if result.Response.UserEngineResourceGroupInfos == nil || len(result.Response.UserEngineResourceGroupInfos) == 0 {
						return resource.NonRetryableError(fmt.Errorf("UserEngineResourceGroupInfos is nil."))
					}

					if len(result.Response.UserEngineResourceGroupInfos) != 1 {
						return resource.NonRetryableError(fmt.Errorf("UserEngineResourceGroupInfos is not 1."))
					}

					state := result.Response.UserEngineResourceGroupInfos[0].ResourceGroupState
					if state != nil {
						if *state == 2 {
							return nil
						}
					} else {
						return resource.NonRetryableError(fmt.Errorf("ResourceGroupState is nil."))
					}

					return resource.RetryableError(fmt.Errorf("UserEngineResourceGroupInfos is not ready, state:%d", *state))
				})

				if waitErr != nil {
					log.Printf("[CRITAL]%s wait for dlc standard engine resource group failed, reason:%+v", logId, waitErr)
					return waitErr
				}
			} else {
				request := dlcv20210125.NewPauseStandardEngineResourceGroupsRequest()
				request.EngineResourceGroupNames = helper.Strings([]string{engineResourceGroupName})
				reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
					result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDlcClient().PauseStandardEngineResourceGroupsWithContext(ctx, request)
					if e != nil {
						return tccommon.RetryError(e)
					} else {
						log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
					}

					if result == nil || result.Response == nil {
						return resource.NonRetryableError(fmt.Errorf("Pause dlc standard engine resource group failed, Response is nil."))
					}

					return nil
				})

				if reqErr != nil {
					log.Printf("[CRITAL]%s pause dlc standard engine resource group failed, reason:%+v", logId, reqErr)
					return reqErr
				}

				// wait
				waitErr := resource.Retry(tccommon.WriteRetryTimeout*4, func() *resource.RetryError {
					result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDlcClient().DescribeStandardEngineResourceGroupsWithContext(ctx, &dlcv20210125.DescribeStandardEngineResourceGroupsRequest{
						Filters: []*dlcv20210125.Filter{
							{
								Name:   helper.String("engine-resource-group-name-unique"),
								Values: helper.Strings([]string{engineResourceGroupName}),
							},
						},
					})
					if e != nil {
						return tccommon.RetryError(e)
					} else {
						log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
					}

					if result == nil || result.Response == nil {
						return resource.NonRetryableError(fmt.Errorf("Describe dlc standard engine resource groups failed, Response is nil."))
					}

					if result.Response.UserEngineResourceGroupInfos == nil || len(result.Response.UserEngineResourceGroupInfos) == 0 {
						return resource.NonRetryableError(fmt.Errorf("UserEngineResourceGroupInfos is nil."))
					}

					if len(result.Response.UserEngineResourceGroupInfos) != 1 {
						return resource.NonRetryableError(fmt.Errorf("UserEngineResourceGroupInfos is not 1."))
					}

					state := result.Response.UserEngineResourceGroupInfos[0].ResourceGroupState
					if state != nil {
						if *state == 3 {
							return nil
						}
					} else {
						return resource.NonRetryableError(fmt.Errorf("ResourceGroupState is nil."))
					}

					return resource.RetryableError(fmt.Errorf("UserEngineResourceGroupInfos is not ready, state:%d", *state))
				})

				if waitErr != nil {
					log.Printf("[CRITAL]%s wait for dlc standard engine resource group failed, reason:%+v", logId, waitErr)
					return waitErr
				}
			}
		}
	}

	return resourceTencentCloudDlcStandardEngineResourceGroupRead(d, meta)
}

func resourceTencentCloudDlcStandardEngineResourceGroupDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dlc_standard_engine_resource_group.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId                   = tccommon.GetLogId(tccommon.ContextNil)
		ctx                     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request                 = dlcv20210125.NewDeleteStandardEngineResourceGroupRequest()
		engineResourceGroupName = d.Id()
	)

	// pause first
	pauseRequest := dlcv20210125.NewPauseStandardEngineResourceGroupsRequest()
	pauseRequest.EngineResourceGroupNames = helper.Strings([]string{engineResourceGroupName})
	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDlcClient().PauseStandardEngineResourceGroupsWithContext(ctx, pauseRequest)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s pause dlc standard engine resource group failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	// wait
	waitErr := resource.Retry(tccommon.WriteRetryTimeout*4, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDlcClient().DescribeStandardEngineResourceGroupsWithContext(ctx, &dlcv20210125.DescribeStandardEngineResourceGroupsRequest{
			Filters: []*dlcv20210125.Filter{
				{
					Name:   helper.String("engine-resource-group-name-unique"),
					Values: helper.Strings([]string{engineResourceGroupName}),
				},
			},
		})
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe dlc standard engine resource groups failed, Response is nil."))
		}

		if result.Response.UserEngineResourceGroupInfos == nil || len(result.Response.UserEngineResourceGroupInfos) == 0 {
			return resource.NonRetryableError(fmt.Errorf("UserEngineResourceGroupInfos is nil."))
		}

		if len(result.Response.UserEngineResourceGroupInfos) != 1 {
			return resource.NonRetryableError(fmt.Errorf("UserEngineResourceGroupInfos is not 1."))
		}

		state := result.Response.UserEngineResourceGroupInfos[0].ResourceGroupState
		if state != nil {
			if *state == 3 {
				return nil
			}
		} else {
			return resource.NonRetryableError(fmt.Errorf("ResourceGroupState is nil."))
		}

		return resource.RetryableError(fmt.Errorf("UserEngineResourceGroupInfos is not pause, state:%d", *state))
	})

	if waitErr != nil {
		log.Printf("[CRITAL]%s wait for dlc standard engine resource group pause failed, reason:%+v", logId, waitErr)
		return waitErr
	}

	// delete end
	request.EngineResourceGroupName = &engineResourceGroupName
	reqErr = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDlcClient().DeleteStandardEngineResourceGroupWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete dlc standard engine resource group failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return nil
}
