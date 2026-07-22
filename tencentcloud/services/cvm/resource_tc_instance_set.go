package cvm

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"strconv"
	"time"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

func ResourceTencentCloudInstanceSet() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudInstanceSetCreate,
		Read:   resourceTencentCloudInstanceSetRead,
		Update: resourceTencentCloudInstanceSetUpdate,
		Delete: resourceTencentCloudInstanceSetDelete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(600 * time.Second),
			Read:   schema.DefaultTimeout(600 * time.Second),
			Update: schema.DefaultTimeout(600 * time.Second),
			Delete: schema.DefaultTimeout(600 * time.Second),
		},

		Schema: map[string]*schema.Schema{
			"image_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "镜像 到 使用 对于 实例. Changing `image_id` 将 cause 实例 reset。",
			},
			"availability_zone": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "可用 可用区 对于 CVM 实例。",
			},
			"instance_count": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "数量 实例 到 是 purchased. 值 范围:[1,100]; 默认值：1。",
			},
			"exclude_instance_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "实例 ids 列表 到 exclude。",
			},
			"instance_name": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "Terraform-CVM-Instance",
				ValidateFunc: tccommon.ValidateStringLengthInRange(2, 128),
				Description:  "名称 实例. max 长度 的 instance_name 是 128，和 默认值为 `Terraform-CVM-实例`。",
			},
			"instance_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: tccommon.ValidateInstanceType,
				Description:  "类型 实例。",
			},
			"hostname": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "hostname 的 实例. Windows 实例: 名称 should 是 combination 的 2 到 15 字符 comprised 的 letters (case insensitive)，numbers，和 hyphens (-). 周期 (.) 是 不 支持，和 名称 不能 是 字符串 的 pure numbers. Other types (such 作为 Linux) 的 实例: 名称 should 是 combination 的 2 到 60 字符，supporting 多个 periods (.). piece between two periods 是 composed 的 letters (case insensitive)，numbers，和 hyphens (-). Modifications 可能 lead 到 reinstallation 的 实例's operating 系统.。",
			},
			"project_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "项目 实例 belongs 到，默认为 0。",
			},
			"placement_group_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "ID placement 组。",
			},
			// payment
			"instance_charge_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      CVM_CHARGE_TYPE_POSTPAID,
				ValidateFunc: tccommon.ValidateAllowedStringValue(CVM_CHARGE_TYPE),
				Description:  "charge 类型 实例. Only support `POSTPAID_BY_HOUR`。",
			},
			// network
			"internet_charge_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue(CVM_INTERNET_CHARGE_TYPE),
				Description:  "Internet charge 类型 实例，有效 值 是 `BANDWIDTH_PREPAID`，`TRAFFIC_POSTPAID_BY_HOUR`，`BANDWIDTH_POSTPAID_BY_HOUR` 和 `BANDWIDTH_PACKAGE`. 此 值 does 不 need 到 是 集合 当 `allocate_public_ip` 是 false。",
			},
			"bandwidth_package_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "带宽 包 ID. 如果 用户 是 standard 用户，then bandwidth_package_id 是 needed，或 默认值 has bandwidth_package_id。",
			},
			"internet_max_bandwidth_out": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Maximum outgoing 带宽 到 公有 网络，measured 在 Mbps (Mega bits per second). 此 值 does 不 need 到 是 集合 当 `allocate_public_ip` 是 false。",
			},
			"allocate_public_ip": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				ForceNew:    true,
				Description: "Associate 公网 IP 地址 使用 实例 在 VPC 或 Classic. Boolean 值，默认为 false。",
			},
			// vpc
			"vpc_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "ID VPC 网络. 如果 您 want 到 create 实例 在 VPC 网络，此 参数 必须 是 集合。",
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "ID VPC 子网. 如果 您 want 到 create 实例 在 VPC 网络，此 参数 必须 是 集合。",
			},
			"private_ip": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "私有 IP 到 是 assigned 到 此 实例，必须 是 在 提供 子网 和 可用。",
			},
			// security group
			"security_groups": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Computed:    true,
				Description: "A 列表 安全 组 IDs 到 associate 使用。",
			},
			// storage
			"system_disk_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      CVM_DISK_TYPE_CLOUD_PREMIUM,
				ValidateFunc: tccommon.ValidateAllowedStringValue(CVM_DISK_TYPE),
				Description:  "System 磁盘 类型 For more 信息 在 limits 的 系统 磁盘 types，see [Storage Overview](https://intl.云.tencent.com/document/product/213/4952). 有效值：`LOCAL_BASIC`: 本地 磁盘，`LOCAL_SSD`: 本地 SSD 磁盘，`CLOUD_SSD`: SSD，`CLOUD_PREMIUM`: Premium Cloud Storage，`CLOUD_BSSD`: Basic SSD. NOTE: 如果 modified， 实例 可能 强制停止",
			},
			"system_disk_size": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      50,
				ValidateFunc: tccommon.ValidateIntegerInRange(50, 1000),
				Description:  "Size 的 系统 磁盘. 有效 值 ranges: (50~1000). 和 单位 是 GB. 默认为 50GB. 如果 modified， 实例 可能 强制停止",
			},
			"system_disk_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "System 磁盘 快照 ID 用于initialize 系统 磁盘. 当 系统 磁盘 类型 是 `LOCAL_BASIC` 和 `LOCAL_SSD`，磁盘 ID 是 不 支持。",
			},
			// enhance services
			"disable_security_service": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Disable enhance 服务 对于 安全，它 是 已启用 通过 默认值. 当 此 options 是 集合，安全 agent won't 是 installed. Modifications 可能 lead 到 reinstallation 的 实例's operating 系统。",
			},
			"disable_monitor_service": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Disable enhance 服务 对于 监控，它 是 已启用 通过 默认值. 当 此 options 是 集合，监控 agent won't 是 installed. Modifications 可能 lead 到 reinstallation 的 实例's operating 系统。",
			},
			// login
			"key_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "键 pair 到 使用 对于 实例，它 looks like `skey-16jig7tx`. Modifications 可能 lead 到 reinstallation 的 实例's operating 系统。",
			},
			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "密码 对于 实例. In 顺序 对于 new 密码 到 take effect， 实例 将 是 restarted after 密码 change. Modifications 可能 lead 到 reinstallation 的 实例's operating 系统。",
			},
			"keep_image_login": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if new == "false" && old == "" || old == "false" && new == "" {
						return true
					} else {
						return old == new
					}
				},
				ConflictsWith: []string{"key_name", "password"},
				Description:   "是否keep 镜像 login 或 不，默认为 `false`. 当 镜像 类型 是 私有 或 shared 或 imported，此 参数 可以 是 集合 `true`. Modifications 可能 lead 到 reinstallation 的 实例's operating 系统。",
			},
			"user_data": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"user_data_raw"},
				Description:   "用户 数据 到 是 injected into 此 实例. Must 是 base64 encoded 和 up 到 16 KB。",
			},
			"user_data_raw": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"user_data"},
				Description:   "用户 数据 到 是 injected into 此 实例，在 plain text. Conflicts 使用 `user_data`. Up 到 16 KB after base64 encoded。",
			},
			// role
			"cam_role_name": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Optional:    true,
				Description: "被授权访问的 CAM 角色名称",
			},
			// Computed values.
			"instance_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Current 状态 实例。",
			},
			"public_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Public IP 的 实例。",
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "创建时间 的 实例。",
			},
			"expired_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "过期时间 的 实例。",
			},
			"instance_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "实例 ID 列表。",
			},
		},
	}
}

func resourceTencentCloudInstanceSetCreate(d *schema.ResourceData, meta interface{}) error {
	doneChan := make(chan struct{}, 1)
	rspChan := make(chan error, 1)

	timeout := d.Timeout(schema.TimeoutCreate)

	go func(d *schema.ResourceData, meta interface{}) {
		e := doResourceTencentCloudInstanceSetCreate(d, meta)
		doneChan <- struct{}{}
		rspChan <- e
	}(d, meta)

	select {
	case <-doneChan:
		return <-rspChan
	case <-time.After(timeout):
		return fmt.Errorf("Do cvm instance set create action timeout, current timeout :[%.3f]s", timeout.Seconds())
	}
}

func resourceTencentCloudInstanceSetRead(d *schema.ResourceData, meta interface{}) error {
	doneChan := make(chan struct{}, 1)
	rspChan := make(chan error, 1)

	timeout := d.Timeout(schema.TimeoutRead)

	go func(d *schema.ResourceData, meta interface{}) {
		e := doResourceTencentCloudInstanceSetRead(d, meta)
		doneChan <- struct{}{}
		rspChan <- e
	}(d, meta)

	select {
	case <-doneChan:
		return <-rspChan
	case <-time.After(timeout):
		return fmt.Errorf("Do cvm instance set read action timeout, current timeout :[%.3f]s", timeout.Seconds())
	}
}

func resourceTencentCloudInstanceSetUpdate(d *schema.ResourceData, meta interface{}) error {
	doneChan := make(chan struct{}, 1)
	rspChan := make(chan error, 1)

	timeout := d.Timeout(schema.TimeoutUpdate)

	go func(d *schema.ResourceData, meta interface{}) {
		e := doResourceTencentCloudInstanceSetUpdate(d, meta)
		doneChan <- struct{}{}
		rspChan <- e
	}(d, meta)

	select {
	case <-doneChan:
		return <-rspChan
	case <-time.After(timeout):
		return fmt.Errorf("Do cvm instance set update action timeout, current timeout :[%.3f]s", timeout.Seconds())
	}
}

func resourceTencentCloudInstanceSetDelete(d *schema.ResourceData, meta interface{}) error {
	doneChan := make(chan struct{}, 1)
	rspChan := make(chan error, 1)

	timeout := d.Timeout(schema.TimeoutDelete)

	go func(d *schema.ResourceData, meta interface{}) {
		e := doResourceTencentCloudInstanceSetDelete(d, meta)
		doneChan <- struct{}{}
		rspChan <- e
	}(d, meta)

	select {
	case <-doneChan:
		return <-rspChan
	case <-time.After(timeout):
		return fmt.Errorf("Do cvm instance set delete action timeout, current timeout :[%.3f]s", timeout.Seconds())
	}
}

func doResourceTencentCloudInstanceSetCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_instance_set.create")()
	logId := tccommon.GetLogId(tccommon.ContextNil)

	var instanceCount int

	request := cvm.NewRunInstancesRequest()
	request.ImageId = helper.String(d.Get("image_id").(string))
	request.Placement = &cvm.Placement{
		Zone: helper.String(d.Get("availability_zone").(string)),
	}
	if v, ok := d.GetOk("project_id"); ok {
		projectId := int64(v.(int))
		request.Placement.ProjectId = &projectId
	}
	if v, ok := d.GetOk("instance_name"); ok {
		request.InstanceName = helper.String(v.(string))
	}
	if v, ok := d.GetOk("instance_count"); ok {
		instanceCount = v.(int)
		request.InstanceCount = helper.Int64(int64(instanceCount))
	}
	if v, ok := d.GetOk("instance_type"); ok {
		request.InstanceType = helper.String(v.(string))
	}
	if v, ok := d.GetOk("hostname"); ok {
		request.HostName = helper.String(v.(string))
	}
	if v, ok := d.GetOk("cam_role_name"); ok {
		request.CamRoleName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("instance_charge_type"); ok {
		instanceChargeType := v.(string)
		request.InstanceChargeType = &instanceChargeType
	}
	if v, ok := d.GetOk("placement_group_id"); ok {
		request.DisasterRecoverGroupIds = []*string{helper.String(v.(string))}
	}

	// network
	request.InternetAccessible = &cvm.InternetAccessible{}
	if v, ok := d.GetOk("internet_charge_type"); ok {
		request.InternetAccessible.InternetChargeType = helper.String(v.(string))
	}
	if v, ok := d.GetOk("internet_max_bandwidth_out"); ok {
		maxBandwidthOut := int64(v.(int))
		request.InternetAccessible.InternetMaxBandwidthOut = &maxBandwidthOut
	}
	if v, ok := d.GetOk("bandwidth_package_id"); ok {
		request.InternetAccessible.BandwidthPackageId = helper.String(v.(string))
	}
	if v, ok := d.GetOkExists("allocate_public_ip"); ok {
		allocatePublicIp := v.(bool)
		request.InternetAccessible.PublicIpAssigned = &allocatePublicIp
	}

	// vpc
	if v, ok := d.GetOk("vpc_id"); ok {
		request.VirtualPrivateCloud = &cvm.VirtualPrivateCloud{}
		request.VirtualPrivateCloud.VpcId = helper.String(v.(string))

		if v, ok = d.GetOk("subnet_id"); ok {
			request.VirtualPrivateCloud.SubnetId = helper.String(v.(string))
		}

		if v, ok = d.GetOk("private_ip"); ok {
			request.VirtualPrivateCloud.PrivateIpAddresses = []*string{helper.String(v.(string))}
		}
	}

	if v, ok := d.GetOk("security_groups"); ok {
		securityGroups := v.(*schema.Set).List()
		request.SecurityGroupIds = make([]*string, 0, len(securityGroups))
		for _, securityGroup := range securityGroups {
			request.SecurityGroupIds = append(request.SecurityGroupIds, helper.String(securityGroup.(string)))
		}
	}

	// storage
	request.SystemDisk = &cvm.SystemDisk{}
	if v, ok := d.GetOk("system_disk_type"); ok {
		request.SystemDisk.DiskType = helper.String(v.(string))
	}
	if v, ok := d.GetOk("system_disk_size"); ok {
		diskSize := int64(v.(int))
		request.SystemDisk.DiskSize = &diskSize
	}
	if v, ok := d.GetOk("system_disk_id"); ok {
		request.SystemDisk.DiskId = helper.String(v.(string))
	}

	// enhanced service
	request.EnhancedService = &cvm.EnhancedService{}
	if v, ok := d.GetOkExists("disable_security_service"); ok {
		securityService := !(v.(bool))
		request.EnhancedService.SecurityService = &cvm.RunSecurityServiceEnabled{
			Enabled: &securityService,
		}
	}
	if v, ok := d.GetOkExists("disable_monitor_service"); ok {
		monitorService := !(v.(bool))
		request.EnhancedService.MonitorService = &cvm.RunMonitorServiceEnabled{
			Enabled: &monitorService,
		}
	}

	// login
	request.LoginSettings = &cvm.LoginSettings{}
	if v, ok := d.GetOk("key_name"); ok {
		request.LoginSettings.KeyIds = []*string{helper.String(v.(string))}
	}
	if v, ok := d.GetOk("password"); ok {
		request.LoginSettings.Password = helper.String(v.(string))
	}
	v := d.Get("keep_image_login").(bool)
	if v {
		request.LoginSettings.KeepImageLogin = helper.String(CVM_IMAGE_LOGIN)
	} else {
		request.LoginSettings.KeepImageLogin = helper.String(CVM_IMAGE_LOGIN_NOT)
	}

	if v, ok := d.GetOk("user_data"); ok {
		request.UserData = helper.String(v.(string))
	}
	if v, ok := d.GetOk("user_data_raw"); ok {
		userData := base64.StdEncoding.EncodeToString([]byte(v.(string)))
		request.UserData = &userData
	}

	instanceIds := make([]*string, 0)

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check("create")
		response, err := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCvmClient().RunInstances(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), err.Error())
			e, ok := err.(*sdkErrors.TencentCloudSDKError)
			if ok && tccommon.IsContains(CVM_RETRYABLE_ERROR, e.Code) {
				time.Sleep(1 * time.Second) // 需要重试的话，等待1s进行重试
				return resource.RetryableError(fmt.Errorf("cvm create error: %s, retrying", e.Error()))
			}
			return resource.NonRetryableError(err)
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
		if len(response.Response.InstanceIdSet) < instanceCount {
			err = fmt.Errorf("number of instances is less than %s", strconv.Itoa(instanceCount))
			return resource.NonRetryableError(err)
		}
		instanceIds = response.Response.InstanceIdSet

		return nil
	})
	if err != nil {
		return err
	}

	_ = d.Set("instance_ids", instanceIds)
	d.SetId(helper.StrListToStr(instanceIds))

	return nil
}

func doResourceTencentCloudInstanceSetRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_instance_set.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	var instanceSetIds []*string
	if v, ok := d.GetOk("instance_ids"); ok {
		instanceSetIds = helper.InterfacesStringsPoint(v.([]interface{}))
	}

	client := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
	cvmService := CvmService{
		client: client,
	}
	var instanceSet []*cvm.Instance
	var errRet error
	instanceSet, errRet = cvmService.DescribeInstanceSetByIds(ctx, helper.StrListToStr(instanceSetIds))
	if errRet != nil {
		return errRet
	}

	if instanceSet == nil {
		d.SetId("")
		return nil
	}

	instance := instanceSet[0]

	_ = d.Set("image_id", instance.ImageId)
	_ = d.Set("availability_zone", instance.Placement.Zone)
	_ = d.Set("instance_name", d.Get("instance_name"))
	_ = d.Set("instance_type", instance.InstanceType)
	_ = d.Set("project_id", instance.Placement.ProjectId)
	_ = d.Set("instance_charge_type", instance.InstanceChargeType)
	_ = d.Set("internet_charge_type", instance.InternetAccessible.InternetChargeType)
	_ = d.Set("internet_max_bandwidth_out", instance.InternetAccessible.InternetMaxBandwidthOut)
	_ = d.Set("vpc_id", instance.VirtualPrivateCloud.VpcId)
	_ = d.Set("subnet_id", instance.VirtualPrivateCloud.SubnetId)
	_ = d.Set("security_groups", helper.StringsInterfaces(instance.SecurityGroupIds))
	_ = d.Set("system_disk_type", instance.SystemDisk.DiskType)
	_ = d.Set("system_disk_size", instance.SystemDisk.DiskSize)
	_ = d.Set("system_disk_id", instance.SystemDisk.DiskId)
	_ = d.Set("instance_status", instance.InstanceState)
	_ = d.Set("create_time", instance.CreatedTime)
	_ = d.Set("expired_time", instance.ExpiredTime)
	_ = d.Set("cam_role_name", instance.CamRoleName)

	if _, ok := d.GetOkExists("allocate_public_ip"); !ok {
		_ = d.Set("allocate_public_ip", len(instance.PublicIpAddresses) > 0)
	}

	if len(instance.PrivateIpAddresses) > 0 {
		_ = d.Set("private_ip", instance.PrivateIpAddresses[0])
	}
	if len(instance.PublicIpAddresses) > 0 {
		_ = d.Set("public_ip", instance.PublicIpAddresses[0])
	}
	if len(instance.LoginSettings.KeyIds) > 0 {
		_ = d.Set("key_name", instance.LoginSettings.KeyIds[0])
	} else {
		_ = d.Set("key_name", "")
	}
	if instance.LoginSettings.KeepImageLogin != nil {
		_ = d.Set("keep_image_login", *instance.LoginSettings.KeepImageLogin == CVM_IMAGE_LOGIN)
	}

	return nil
}

func doResourceTencentCloudInstanceSetUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	defer tccommon.LogElapsed("resource.tencentcloud_instance_set.update")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	cvmService := CvmService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}
	if d.HasChange("instance_count") {
		return fmt.Errorf("`instance_count` do not support change now, please use `resource_tc_instace` instead.")
	}

	if d.HasChange("exclude_instance_ids") {
		old, new := d.GetChange("exclude_instance_ids")
		olds := old.(*schema.Set)
		news := new.(*schema.Set)
		needExclude := news.Difference(olds).List()
		needCreate := olds.Difference(news).List()

		// need delete instance
		if len(needExclude) > 0 {
			instanceSetIds := helper.StrListToStr(helper.InterfacesStringsPoint(needExclude))
			err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
				errRet := cvmService.DeleteInstanceSetByIds(ctx, instanceSetIds)
				if errRet != nil {
					log.Printf("[CRITAL][first delete]%s api[%s] fail, reason[%s]\n",
						logId, "delete", errRet.Error())
					e, ok := errRet.(*sdkErrors.TencentCloudSDKError)
					if ok && tccommon.IsContains(CVM_RETRYABLE_ERROR, e.Code) {
						time.Sleep(1 * time.Second) // 需要重试的话，等待1s进行重试
						return resource.RetryableError(fmt.Errorf("[first delete]cvm delete error: %s, retrying", e.Error()))
					}
					return resource.NonRetryableError(errRet)
				}
				return nil
			})
			if err != nil {
				return err
			}
		}

		createInstanceIds := make([]*string, 0)
		// need create instance
		if len(needCreate) > 0 {
			instanceCount := len(needCreate)
			request := cvm.NewRunInstancesRequest()
			request.ImageId = helper.String(d.Get("image_id").(string))
			request.InstanceCount = helper.Int64(int64(instanceCount))
			request.Placement = &cvm.Placement{
				Zone: helper.String(d.Get("availability_zone").(string)),
			}
			if v, ok := d.GetOk("project_id"); ok {
				projectId := int64(v.(int))
				request.Placement.ProjectId = &projectId
			}
			if v, ok := d.GetOk("instance_name"); ok {
				request.InstanceName = helper.String(v.(string))
			}

			if v, ok := d.GetOk("instance_type"); ok {
				request.InstanceType = helper.String(v.(string))
			}
			if v, ok := d.GetOk("hostname"); ok {
				request.HostName = helper.String(v.(string))
			}
			if v, ok := d.GetOk("cam_role_name"); ok {
				request.CamRoleName = helper.String(v.(string))
			}

			if v, ok := d.GetOk("instance_charge_type"); ok {
				instanceChargeType := v.(string)
				request.InstanceChargeType = &instanceChargeType
			}
			if v, ok := d.GetOk("placement_group_id"); ok {
				request.DisasterRecoverGroupIds = []*string{helper.String(v.(string))}
			}

			// network
			request.InternetAccessible = &cvm.InternetAccessible{}
			if v, ok := d.GetOk("internet_charge_type"); ok {
				request.InternetAccessible.InternetChargeType = helper.String(v.(string))
			}
			if v, ok := d.GetOk("internet_max_bandwidth_out"); ok {
				maxBandwidthOut := int64(v.(int))
				request.InternetAccessible.InternetMaxBandwidthOut = &maxBandwidthOut
			}
			if v, ok := d.GetOk("bandwidth_package_id"); ok {
				request.InternetAccessible.BandwidthPackageId = helper.String(v.(string))
			}
			if v, ok := d.GetOkExists("allocate_public_ip"); ok {
				allocatePublicIp := v.(bool)
				request.InternetAccessible.PublicIpAssigned = &allocatePublicIp
			}

			// vpc
			if v, ok := d.GetOk("vpc_id"); ok {
				request.VirtualPrivateCloud = &cvm.VirtualPrivateCloud{}
				request.VirtualPrivateCloud.VpcId = helper.String(v.(string))

				if v, ok = d.GetOk("subnet_id"); ok {
					request.VirtualPrivateCloud.SubnetId = helper.String(v.(string))
				}
			}

			if v, ok := d.GetOk("security_groups"); ok {
				securityGroups := v.(*schema.Set).List()
				request.SecurityGroupIds = make([]*string, 0, len(securityGroups))
				for _, securityGroup := range securityGroups {
					request.SecurityGroupIds = append(request.SecurityGroupIds, helper.String(securityGroup.(string)))
				}
			}

			// storage
			request.SystemDisk = &cvm.SystemDisk{}
			if v, ok := d.GetOk("system_disk_type"); ok {
				request.SystemDisk.DiskType = helper.String(v.(string))
			}
			if v, ok := d.GetOk("system_disk_size"); ok {
				diskSize := int64(v.(int))
				request.SystemDisk.DiskSize = &diskSize
			}

			// enhanced service
			request.EnhancedService = &cvm.EnhancedService{}
			if v, ok := d.GetOkExists("disable_security_service"); ok {
				securityService := !(v.(bool))
				request.EnhancedService.SecurityService = &cvm.RunSecurityServiceEnabled{
					Enabled: &securityService,
				}
			}
			if v, ok := d.GetOkExists("disable_monitor_service"); ok {
				monitorService := !(v.(bool))
				request.EnhancedService.MonitorService = &cvm.RunMonitorServiceEnabled{
					Enabled: &monitorService,
				}
			}

			// login
			request.LoginSettings = &cvm.LoginSettings{}
			if v, ok := d.GetOk("key_name"); ok {
				request.LoginSettings.KeyIds = []*string{helper.String(v.(string))}
			}
			if v, ok := d.GetOk("password"); ok {
				request.LoginSettings.Password = helper.String(v.(string))
			}
			v := d.Get("keep_image_login").(bool)
			if v {
				request.LoginSettings.KeepImageLogin = helper.String(CVM_IMAGE_LOGIN)
			} else {
				request.LoginSettings.KeepImageLogin = helper.String(CVM_IMAGE_LOGIN_NOT)
			}

			if v, ok := d.GetOk("user_data"); ok {
				request.UserData = helper.String(v.(string))
			}
			if v, ok := d.GetOk("user_data_raw"); ok {
				userData := base64.StdEncoding.EncodeToString([]byte(v.(string)))
				request.UserData = &userData
			}

			err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
				ratelimit.Check("create")
				response, err := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCvmClient().RunInstances(request)
				if err != nil {
					log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
						logId, request.GetAction(), request.ToJsonString(), err.Error())
					e, ok := err.(*sdkErrors.TencentCloudSDKError)
					if ok && tccommon.IsContains(CVM_RETRYABLE_ERROR, e.Code) {
						time.Sleep(1 * time.Second) // 需要重试的话，等待1s进行重试
						return resource.RetryableError(fmt.Errorf("cvm create error: %s, retrying", e.Error()))
					}
					return resource.NonRetryableError(err)
				}
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
					logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
				if len(response.Response.InstanceIdSet) < instanceCount {
					err = fmt.Errorf("number of instances is less than %s", strconv.Itoa(instanceCount))
					return resource.NonRetryableError(err)
				}
				createInstanceIds = response.Response.InstanceIdSet

				return nil
			})
			if err != nil {
				return err
			}

		}

		var instanceSetIds []*string
		if v, ok := d.GetOk("instance_ids"); ok {
			instanceSetIds = helper.InterfacesStringsPoint(v.([]interface{}))
		}

		var newInstanceSetIds []*string
		for _, v := range instanceSetIds {
			ins := v
			noEqual := true
			for _, u := range needExclude {
				if *ins == u.(string) {
					noEqual = false
				}
			}
			if noEqual {
				newInstanceSetIds = append(newInstanceSetIds, ins)
			}
		}

		if len(needCreate) > 0 {
			for _, v := range createInstanceIds {
				ins := v
				newInstanceSetIds = append(newInstanceSetIds, ins)
			}
		}
		_ = d.Set("instance_ids", newInstanceSetIds)

	}

	return nil
}

func doResourceTencentCloudInstanceSetDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_instance_set.delete")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	//instanceSetIds := d.Id()

	cvmService := CvmService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}

	var instanceSetIds []*string
	if v, ok := d.GetOk("instance_ids"); ok {
		instanceSetIds = helper.InterfacesStringsPoint(v.([]interface{}))
	}

	// delete
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		errRet := cvmService.DeleteInstanceSetByIds(ctx, helper.StrListToStr(instanceSetIds))
		if errRet != nil {
			log.Printf("[CRITAL][first delete]%s api[%s] fail, reason[%s]\n",
				logId, "delete", errRet.Error())
			e, ok := errRet.(*sdkErrors.TencentCloudSDKError)
			if ok && tccommon.IsContains(CVM_RETRYABLE_ERROR, e.Code) {
				time.Sleep(1 * time.Second) // 需要重试的话，等待1s进行重试
				return resource.RetryableError(fmt.Errorf("[first delete]cvm delete error: %s, retrying", e.Error()))
			}
			return resource.NonRetryableError(errRet)
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
