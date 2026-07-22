package cvm

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svctag "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tag"
	svcvpc "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/vpc"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

func ResourceTencentCloudEip() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudEipCreate,
		Read:   resourceTencentCloudEipRead,
		Update: resourceTencentCloudEipUpdate,
		Delete: resourceTencentCloudEipDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "名称 eip。",
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     svcvpc.EIP_TYPE_EIP,
				ForceNew:    true,
				Description: "类型 eip. 有效 值: `EIP` 和 `AnycastEIP` 和 `HighQualityEIP` 和 `AntiDDoSEIP` 和 `ResidentialEIP`. 默认为 `EIP`。",
			},
			"anycast_zone": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "可用区 的 anycast. 有效 值: `ANYCAST_ZONE_GLOBAL` 和 `ANYCAST_ZONE_OVERSEAS`。",
			},
			"applicable_for_clb": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "表示是否anycast eip 可以 是 associated 到 CLB。",
				Deprecated:  "It has been deprecated from version 1.27.0.",
			},
			"internet_service_provider": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Internet 服务 provider 的 eip. 有效 值: `BGP`，`CMCC`，`CTCC` 和 `CUCC`。",
			},
			"internet_charge_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "charge 类型 eip. 有效值：`BANDWIDTH_PACKAGE`，`BANDWIDTH_POSTPAID_BY_HOUR`，`BANDWIDTH_PREPAID_BY_MONTH` 和 `TRAFFIC_POSTPAID_BY_HOUR`。",
			},
			"prepaid_period": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: tccommon.ValidateAllowedIntValue(svcvpc.EIP_AVAILABLE_PERIOD),
				Description:  "周期 的 实例. 默认值：`1`. 有效 值: `1`，`2`，`3`，`4`，`6`，`7`，`8`，`9`，`12`，`24`，`36`. NOTES: 必须 集合 当 `internet_charge_type` 是 `BANDWIDTH_PREPAID_BY_MONTH`。",
			},

			"auto_renew_flag": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: tccommon.ValidateAllowedIntValue([]int{0, 1, 2}),
				Description:  "Auto 续费标识 0 - 默认值 state (manual renew); 1 - automatic renew; 2 - explicit 无 automatic renew. NOTES: Only 支持 prepaid EIP。",
			},

			"internet_max_bandwidth_out": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "带宽 限制 的 EIP，单位 是 Mbps。",
			},
			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "标签 的 eip。",
			},
			"bandwidth_package_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "ID 带宽 包，它 将 集合 当 `internet_charge_type` 是 `BANDWIDTH_PACKAGE`。",
			},
			"egress": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Network egress. It 默认为 `center_egress1`. 如果 您 want 到 try egress 功能，please [提交 ticket](https://console.云.tencent.com/workorder/category)。",
			},
			"anti_ddos_package_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "ID anti DDos 包，它 必须 集合 当 `类型` 是 `AntiDDoSEIP`。",
			},
			"cdc_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "CDC Unique ID。",
			},

			// computed
			"public_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "elastic IP 地址",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "EIP 当前 状态",
			},
		},
	}
}

func resourceTencentCloudEipCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_eip.create")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	client := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
	vpcService := svcvpc.NewVpcService(client)

	var internetChargeType string

	request := vpc.NewAllocateAddressesRequest()
	if v, ok := d.GetOk("type"); ok {
		request.AddressType = helper.String(v.(string))
	}
	if v, ok := d.GetOk("anycast_zone"); ok {
		request.AnycastZone = helper.String(v.(string))
	}
	if v, ok := d.GetOk("internet_service_provider"); ok {
		request.InternetServiceProvider = helper.String(v.(string))
	}
	if v, ok := d.GetOk("internet_charge_type"); ok {
		internetChargeType = v.(string)
		request.InternetChargeType = helper.String(v.(string))
	}
	if v, ok := d.GetOk("internet_max_bandwidth_out"); ok {
		request.InternetMaxBandwidthOut = helper.IntInt64(v.(int))
	}

	if internetChargeType == "BANDWIDTH_PREPAID_BY_MONTH" {
		addressChargePrepaid := vpc.AddressChargePrepaid{}
		period := d.Get("prepaid_period")
		renewFlag := d.Get("auto_renew_flag")
		addressChargePrepaid.Period = helper.IntInt64(period.(int))
		addressChargePrepaid.AutoRenewFlag = helper.IntInt64(renewFlag.(int))
		request.AddressChargePrepaid = &addressChargePrepaid
	}

	if v := helper.GetTags(d, "tags"); len(v) > 0 {
		for tagKey, tagValue := range v {
			tag := vpc.Tag{
				Key:   helper.String(tagKey),
				Value: helper.String(tagValue),
			}
			request.Tags = append(request.Tags, &tag)
		}
	}
	if v, ok := d.GetOk("bandwidth_package_id"); ok {
		request.BandwidthPackageId = helper.String(v.(string))
	}
	if v, ok := d.GetOk("name"); ok {
		request.AddressName = helper.String(v.(string))
	}
	if v, ok := d.GetOk("egress"); ok {
		request.Egress = helper.String(v.(string))
	}
	if v, ok := d.GetOk("anti_ddos_package_id"); ok {
		request.AntiDDoSPackageId = helper.String(v.(string))
	}
	if v, ok := d.GetOk("cdc_id"); ok {
		request.DedicatedClusterId = helper.String(v.(string))
	}

	eipId := ""
	taskId := ""
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		response, err := client.UseVpcClient().AllocateAddresses(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), err.Error())
			return tccommon.RetryError(err)
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if len(response.Response.AddressSet) < 1 {
			return resource.RetryableError(fmt.Errorf("eip id is nil"))
		}
		eipId = *response.Response.AddressSet[0]
		// taskId maybe nil
		if response.Response.TaskId != nil {
			taskId = *response.Response.TaskId
		}

		return nil
	})
	if err != nil {
		return err
	}
	d.SetId(eipId)

	// wait for status
	if taskId != "" {
		taskIdUint64, err := strconv.ParseUint(taskId, 10, 64)
		if err != nil {
			return err
		}
		taskRequest := vpc.NewDescribeTaskResultRequest()
		taskRequest.TaskId = &taskIdUint64
		err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(taskRequest.GetAction())
			taskResponse, err := client.UseVpcClient().DescribeTaskResult(taskRequest)
			if err != nil {
				return tccommon.RetryError(err)
			}
			if taskResponse.Response.Result != nil && *taskResponse.Response.Result == svcvpc.VPN_TASK_STATUS_RUNNING {
				return resource.RetryableError(errors.New("eip task is running"))
			}
			if taskResponse.Response.Result != nil && *taskResponse.Response.Result == svcvpc.VPN_TASK_STATUS_FAILED {
				return resource.NonRetryableError(errors.New("eip task is failed"))
			}
			return nil
		})
		if err != nil {
			return err
		}
	}

	err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		eip, errRet := vpcService.DescribeEipById(ctx, eipId)
		if errRet != nil {
			return tccommon.RetryError(errRet)
		}
		if eip == nil {
			return resource.NonRetryableError(errors.New("eip is nil"))
		}
		if eip.AddressStatus != nil && *eip.AddressStatus == svcvpc.EIP_STATUS_CREATING {
			return resource.RetryableError(fmt.Errorf("eip is still creating"))
		}
		return nil
	})

	if err != nil {
		return err
	}

	return resourceTencentCloudEipRead(d, meta)
}

func resourceTencentCloudEipRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_eip.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	client := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
	vpcService := svcvpc.NewVpcService(client)
	tagService := svctag.NewTagService(client)
	region := client.Region

	eipId := d.Id()
	var eip *vpc.Address
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		instance, errRet := vpcService.DescribeEipById(ctx, eipId)
		if errRet != nil {
			return tccommon.RetryError(errRet)
		}
		eip = instance
		return nil
	})
	if err != nil {
		return err
	}
	if eip == nil {
		d.SetId("")
		return nil
	}

	tags, err := tagService.DescribeResourceTags(ctx, svcvpc.VPC_SERVICE_TYPE, svcvpc.EIP_RESOURCE_TYPE, region, eipId)
	if err != nil {
		log.Printf("[CRITAL]%s describe eip tags failed: %+v", logId, err)
		return err
	}

	bgp, err := vpcService.DescribeVpcBandwidthPackageByEip(ctx, eipId)
	if err != nil {
		log.Printf("[CRITAL]%s describe eip tags failed: %+v", logId, err)
		return err
	}
	_ = d.Set("name", eip.AddressName)
	_ = d.Set("type", eip.AddressType)
	_ = d.Set("public_ip", eip.AddressIp)
	_ = d.Set("status", eip.AddressStatus)
	_ = d.Set("internet_charge_type", eip.InternetChargeType)
	_ = d.Set("tags", tags)

	if eip.Bandwidth != nil {
		_ = d.Set("internet_max_bandwidth_out", eip.Bandwidth)
	}

	if eip.Egress != nil {
		_ = d.Set("egress", eip.Egress)
	}

	if eip.AntiDDoSPackageId != nil {
		_ = d.Set("anti_ddos_package_id", eip.AntiDDoSPackageId)
	}

	if eip.DedicatedClusterId != nil {
		_ = d.Set("cdc_id", eip.DedicatedClusterId)
	}

	if bgp != nil {
		_ = d.Set("bandwidth_package_id", bgp.BandwidthPackageId)
	}
	return nil
}

func resourceTencentCloudEipUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_eip.update")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	client := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
	vpcService := svcvpc.NewVpcService(client)
	tagService := svctag.NewTagService(client)
	region := client.Region

	eipId := d.Id()

	d.Partial(true)

	unsupportedUpdateFields := []string{
		"bandwidth_package_id",
		"anti_ddos_package_id",
		"egress",
	}
	for _, field := range unsupportedUpdateFields {
		if d.HasChange(field) {
			return fmt.Errorf("tencentcloud_eip update on %s is not support yet", field)
		}
	}

	if d.HasChange("name") {
		name := d.Get("name").(string)
		err := vpcService.ModifyEipName(ctx, eipId, name)
		if err != nil {
			return err
		}
	}

	if d.HasChange("internet_charge_type") {
		var (
			chargeType   string
			bandWidthOut int
		)

		if v, ok := d.GetOk("internet_charge_type"); ok {
			chargeType = v.(string)
		}
		if v, ok := d.GetOk("internet_max_bandwidth_out"); ok {
			bandWidthOut = v.(int)
		}

		period := d.Get("prepaid_period").(int)
		renewFlag := d.Get("auto_renew_flag").(int)

		if chargeType != "" && bandWidthOut != 0 {
			err := vpcService.ModifyEipInternetChargeType(ctx, eipId, chargeType, bandWidthOut, period, renewFlag)
			if err != nil {
				return err
			}
		}
	}

	if d.HasChange("internet_max_bandwidth_out") {
		if v, ok := d.GetOk("internet_max_bandwidth_out"); ok {
			bandwidthOut := v.(int)
			err := vpcService.ModifyEipBandwidthOut(ctx, eipId, bandwidthOut)
			if err != nil {
				return err
			}

		}
	}

	if d.HasChange("prepaid_period") || d.HasChange("auto_renew_flag") {
		period := d.Get("prepaid_period").(int)
		renewFlag := d.Get("auto_renew_flag").(int)
		err := vpcService.RenewAddress(ctx, eipId, period, renewFlag)
		if err != nil {
			return err
		}
	}

	if d.HasChange("tags") {
		oldTags, newTags := d.GetChange("tags")
		replaceTags, deleteTags := svctag.DiffTags(oldTags.(map[string]interface{}), newTags.(map[string]interface{}))
		resourceName := tccommon.BuildTagResourceName(svcvpc.VPC_SERVICE_TYPE, svcvpc.EIP_RESOURCE_TYPE, region, eipId)

		if err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags); err != nil {
			log.Printf("[CRITAL]%s update eip tags failed: %+v", logId, err)
			return err
		}

	}

	d.Partial(false)

	return resourceTencentCloudEipRead(d, meta)
}

func resourceTencentCloudEipDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_eip.delete")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	vpcService := svcvpc.NewVpcService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
	eipId := d.Id()
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		errRet := vpcService.UnattachEip(ctx, eipId)
		if errRet != nil {
			return tccommon.RetryError(errRet, "DesOperation.MutexTaskRunning")
		}
		return nil
	})
	if err != nil {
		return err
	}

	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		errRet := vpcService.DeleteEip(ctx, eipId)
		if errRet != nil {
			return tccommon.RetryError(errRet, "DesOperation.MutexTaskRunning", "OperationDenied.MutexTaskRunning")
		}
		return nil
	})
	if err != nil {
		return err
	}

	var internetChargeType string
	if v, ok := d.GetOk("internet_charge_type"); ok {
		internetChargeType = v.(string)
	}

	if internetChargeType == "BANDWIDTH_PREPAID_BY_MONTH" {
		// isolated
		err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			eip, errRet := vpcService.DescribeEipById(ctx, eipId)
			if errRet != nil {
				return tccommon.RetryError(errRet)
			}
			if !*eip.IsArrears {
				return resource.RetryableError(fmt.Errorf("eip is still isolate"))
			}
			return nil
		})
		if err != nil {
			return err
		}

		// release
		err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			errRet := vpcService.DeleteEip(ctx, eipId)
			if errRet != nil {
				return tccommon.RetryError(errRet, "DesOperation.MutexTaskRunning", "OperationDenied.MutexTaskRunning")
			}
			return nil
		})
		if err != nil {
			return err
		}
	}

	err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		eip, errRet := vpcService.DescribeEipById(ctx, eipId)
		if errRet != nil {
			return tccommon.RetryError(errRet)
		}
		if eip != nil {
			return resource.RetryableError(fmt.Errorf("eip is still deleting"))
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
