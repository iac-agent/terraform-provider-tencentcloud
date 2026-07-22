package mariadb

import (
	"context"
	"fmt"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svctag "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	mariadb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mariadb/v20170312"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudMariadbInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMariadbInstanceCreate,
		Read:   resourceTencentCloudMariadbInstanceRead,
		Update: resourceTencentCloudMariadbInstanceUpdate,
		Delete: resourceTencentCloudMariadbInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "实例 ID，uniquely identifies TDSQL 实例。",
			},

			"instance_name": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeString,
				Description: "实例名称，您 可以 集合 名称 实例 independently through 此 字段。",
			},

			"zones": {
				Required: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "实例 节点 availability 可用区 distribution，up 到 two availability zones 可以 是 filled. 当 分片 规格 是 一个 master 和 two slaves，two 的 nodes 是 在 first availability 可用区",
			},

			"node_count": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "节点数量，2 是 一个 master 和 一个 slave，3 是 一个 master 和 two slaves。",
			},

			"memory": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Memory 大小，单位: GB，可以 是 获取 通过 querying 实例 specifications through DescribeDBInstanceSpecs。",
			},

			"storage": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Storage 大小，单位: GB. You 可以 查询 实例 specifications through DescribeDBInstanceSpecs 到 obtain lower 和 upper limits 的 磁盘 specifications corresponding 到 different 内存 sizes。",
			},

			"period": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "时长 的 purchase，单位: month。",
			},

			"auto_voucher": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeBool,
				Description: "是否automatically 使用 voucher 对于 payment， 默认为 不 使用。",
			},

			"voucher_ids": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "A 列表 voucher IDs. Currently，仅 一个 voucher 可以 是 指定。",
			},

			"vip": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "内网 IP 地址",
			},

			"vpc_id": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Virtual 私有 网络 ID，如果未传入，它 表示 该 它 是 创建 作为 basic 网络。",
			},

			"subnet_id": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Virtual 私有 网络 子网 ID，必填 当 VpcId 是 不 空。",
			},

			"project_id": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "项目 ID，其中 可以 是 获取 通过 viewing 项目 列表，如果未传入，它 将 是 associated 使用 默认值 项目。",
			},

			"db_version_id": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Database 引擎 版本，currently 可用: 8.0.18，10.1.9，5.7.17. 如果未传入， 默认为 Percona 5.7.17。",
			},

			"security_group_ids": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "安全组 ID 列表。",
			},

			"auto_renew_flag": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Automatic renewal flag，1: automatic renewal，2: 无 automatic renewal。",
			},

			"ipv6_flag": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Whether IPv6 是 支持。",
			},

			"app_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "ID 应用 到 其中 实例 belongs。",
			},

			"region": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "名称 地域 其中 实例 是 located，such 作为 ap-shanghai。",
			},

			"status": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "实例状态: 0 creating，1 process processing，2 running，3 实例 不 initialized，-1 实例 isolated，4 实例 initializing，5 实例 deleting，6 实例 restarting，7 数据 迁移。",
			},

			"vport": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Intranet 端口",
			},
			"wan_domain": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "域名 名称 accessed 从 外部 网络，其中 可以 是 resolved 通过 公有 网络。",
			},
			"wan_vip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Extranet IP 地址，accessible 从 公有 网络。",
			},
			"wan_port": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Internet 端口",
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "实例 创建时间， 格式 是 2006-01-02 15:04:05。",
			},
			"update_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "last 更新时间 的 实例 在 格式 的 2006-01-02 15:04:05。",
			},
			"period_end_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "实例 过期时间， 格式 是 2006-01-02 15:04:05。",
			},
			"uin": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "账号 到 其中 实例 belongs。",
			},
			"tdsql_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "TDSQL 版本 信息。",
			},
			"is_tmp": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "是否为a temporary 实例，0 表示 无，non-zero 表示 yes。",
			},
			"excluster_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Exclusive 集群 ID，如果 它 是 空，它 表示 normal 实例。",
			},
			"pid": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Product 类型 ID。",
			},
			"qps": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Maximum Qps 值",
			},
			"paymode": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Payment 模式",
			},
			"locker": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Asynchronous 任务 process ID 当 实例 是 在 asynchronous 任务。",
			},
			"status_desc": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "描述 当前 running state 的 实例。",
			},
			"wan_status": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "External 网络 状态，0-unopened; 1-opened; 2-closed; 3-opening。",
			},
			"is_audit_supported": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "是否instance 支持 审计. 1-支持; 0-不 支持。",
			},
			"machine": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Machine Model。",
			},
			"is_encrypt_supported": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Whether 数据 加密 是 支持. 1-支持; 0-不 支持。",
			},
			"cpu": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "CPU 核数 的 实例。",
			},
			"vipv6": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Intranet IPv6。",
			},
			"wan_vipv6": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Internet IPv6。",
			},
			"wan_port_ipv6": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Internet IPv6 端口",
			},
			"wan_status_ipv6": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Internet IPv6 状态",
			},
			"db_engine": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Database Engine。",
			},
			"dcn_flag": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "DCN flag，0-none，1-primary 实例，2-disaster 备份 实例。",
			},
			"dcn_status": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "DCN 状态，0-none，1-creating，2-synchronizing，3-disconnected。",
			},

			"dcn_dst_num": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "数量 DCN disaster recovery 实例。",
			},
			"instance_type": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "1: primary 实例 (exclusive)，2: primary 实例，3: disaster recovery 实例，4: disaster recovery 实例 (exclusive 类型)。",
			},

			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "标签列表",
			},

			"init_params": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				Description: "Parameter 列表. 可选 值 的 此 interface 是: character_set_server (character 集合，必填) enum: utf8,latin1,gbk,utf8mb4,gb18030，lower_case_table_names (表 名称 case sensitive，必填，0 - sensitive; 1 - insensitive)，innodb_page_size (innodb 数据 页面，Default 16K)，sync_mode (sync 模式: 0 - asynchronous; 1 - strong synchronous; 2 - strong synchronous 可以 degenerate. 默认为 strong synchronous 可以 degenerate)。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"param": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "参数 名称",
						},
						"value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "参数 值",
						},
					},
				},
			},

			"dcn_region": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "DCN 来源 地域",
			},

			"dcn_instance_id": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "DCN 来源 实例 ID。",
			},
		},
	}
}

func resourceTencentCloudMariadbInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mariadb_instance.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		request    = mariadb.NewCreateDBInstanceRequest()
		response   = mariadb.NewCreateDBInstanceResponse()
		instanceId string
	)

	if v, ok := d.GetOk("zones"); ok {
		zonesSet := v.(*schema.Set).List()
		for i := range zonesSet {
			zones := zonesSet[i].(string)
			request.Zones = append(request.Zones, &zones)
		}
	}

	if v, _ := d.GetOk("node_count"); v != nil {
		request.NodeCount = helper.IntInt64(v.(int))
	}

	if v, _ := d.GetOk("memory"); v != nil {
		request.Memory = helper.IntInt64(v.(int))
	}

	if v, _ := d.GetOk("storage"); v != nil {
		request.Storage = helper.IntInt64(v.(int))
	}

	if v, _ := d.GetOk("period"); v != nil {
		request.Period = helper.IntInt64(v.(int))
	}

	if v, _ := d.GetOk("auto_voucher"); v != nil {
		request.AutoVoucher = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOk("voucher_ids"); ok {
		voucherIdsSet := v.(*schema.Set).List()
		for i := range voucherIdsSet {
			voucherIds := voucherIdsSet[i].(string)
			request.VoucherIds = append(request.VoucherIds, &voucherIds)
		}
	}

	if v, ok := d.GetOk("vpc_id"); ok {
		request.VpcId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("subnet_id"); ok {
		request.SubnetId = helper.String(v.(string))
	}

	if v, _ := d.GetOk("project_id"); v != nil {
		request.ProjectId = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("db_version_id"); ok {
		request.DbVersionId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("instance_name"); ok {
		request.InstanceName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("security_group_ids"); ok {
		securityGroupIdsSet := v.(*schema.Set).List()
		for i := range securityGroupIdsSet {
			securityGroupIds := securityGroupIdsSet[i].(string)
			request.SecurityGroupIds = append(request.SecurityGroupIds, &securityGroupIds)
		}
	}

	if v, _ := d.GetOk("auto_renew_flag"); v != nil {
		request.AutoRenewFlag = helper.IntInt64(v.(int))
	}

	if v, _ := d.GetOk("ipv6_flag"); v != nil {
		request.Ipv6Flag = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("tags"); ok {
		for key, value := range v.(map[string]interface{}) {
			resourceTag := mariadb.ResourceTag{
				TagKey:   helper.String(key),
				TagValue: helper.String(value.(string)),
			}
			request.ResourceTags = append(request.ResourceTags, &resourceTag)
		}
	}

	if v, ok := d.GetOk("init_params"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			dBParamValue := mariadb.DBParamValue{}
			if v, ok := dMap["param"]; ok {
				dBParamValue.Param = helper.String(v.(string))
			}
			if v, ok := dMap["value"]; ok {
				dBParamValue.Value = helper.String(v.(string))
			}
			request.InitParams = append(request.InitParams, &dBParamValue)
		}
	}

	if v, ok := d.GetOk("dcn_region"); ok {
		request.DcnRegion = helper.String(v.(string))
	}

	if v, ok := d.GetOk("dcn_instance_id"); ok {
		request.DcnInstanceId = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMariadbClient().CreateDBInstance(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil || result.Response.InstanceIds == nil || len(result.Response.InstanceIds) == 0 {
			return resource.RetryableError(fmt.Errorf("Create mariadb instance failed, Response is nil."))
		}

		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create mariadb instance failed, reason:%+v", logId, err)
		return err
	}

	instanceId = *response.Response.InstanceIds[0]
	d.SetId(instanceId)

	// wait
	service := MariadbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	err = resource.Retry(7*tccommon.ReadRetryTimeout, func() *resource.RetryError {
		instance, e := service.DescribeMariadbInstanceById(ctx, instanceId)
		if e != nil {
			return resource.NonRetryableError(e)
		}

		if instance == nil {
			err = fmt.Errorf("mariadb %s instance not exists", instanceId)
			return resource.NonRetryableError(err)
		}

		if *instance.Status == 0 || *instance.Status == 1 || *instance.Status == 4 {
			return resource.RetryableError(fmt.Errorf("create mariadb status is %v,start retrying ...", *instance.Status))
		}

		if *instance.Status == 2 {
			return nil
		}

		err = fmt.Errorf("create mariadb status is %v,we won't wait for it finish", *instance.Status)
		return resource.NonRetryableError(err)
	})

	if err != nil {
		log.Printf("[CRITAL]%s create mariadb fail, reason:%s\n ", logId, err.Error())
		return err
	}

	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		tagService := svctag.NewTagService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
		region := meta.(tccommon.ProviderMeta).GetAPIV3Conn().Region
		resourceName := fmt.Sprintf("qcs::mariadb:%s:uin/:instance/%s", region, instanceId)
		if err := tagService.ModifyTags(ctx, resourceName, tags, nil); err != nil {
			return err
		}
	}

	return resourceTencentCloudMariadbInstanceRead(d, meta)
}

func resourceTencentCloudMariadbInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mariadb_instance.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service    = MariadbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		instanceId = d.Id()
	)

	instance, err := service.DescribeMariadbInstanceById(ctx, instanceId)
	if err != nil {
		return err
	}

	if instance == nil {
		log.Printf("[WARN]%s resource `tencentcloud_mariadb_instance` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	if instance.InstanceId != nil {
		_ = d.Set("instance_id", instance.InstanceId)
	}

	if instance.InstanceName != nil {
		_ = d.Set("instance_name", instance.InstanceName)
	}

	if instance.AppId != nil {
		_ = d.Set("app_id", instance.AppId)
	}

	if instance.ProjectId != nil {
		_ = d.Set("project_id", instance.ProjectId)
	}

	if instance.Region != nil {
		_ = d.Set("region", instance.Region)
	}

	if instance.VpcId != nil {
		_ = d.Set("vpc_id", instance.UniqueVpcId)
	}

	if instance.SubnetId != nil {
		_ = d.Set("subnet_id", instance.UniqueSubnetId)
	}

	if instance.Status != nil {
		_ = d.Set("status", instance.Status)
	}

	if instance.Vip != nil {
		_ = d.Set("vip", instance.Vip)
	}

	if instance.Vport != nil {
		_ = d.Set("vport", instance.Vport)
	}

	if instance.WanDomain != nil {
		_ = d.Set("wan_domain", instance.WanDomain)
	}

	if instance.WanVip != nil {
		_ = d.Set("wan_vip", instance.WanVip)
	}

	if instance.WanPort != nil {
		_ = d.Set("wan_port", instance.WanPort)
	}

	if instance.CreateTime != nil {
		_ = d.Set("create_time", instance.CreateTime)
	}

	if instance.UpdateTime != nil {
		_ = d.Set("update_time", instance.UpdateTime)
	}

	if instance.AutoRenewFlag != nil {
		_ = d.Set("auto_renew_flag", instance.AutoRenewFlag)
	}

	if instance.PeriodEndTime != nil {
		_ = d.Set("period_end_time", instance.PeriodEndTime)
	}

	if instance.Uin != nil {
		_ = d.Set("uin", instance.Uin)
	}

	if instance.TdsqlVersion != nil {
		_ = d.Set("tdsql_version", instance.TdsqlVersion)
	}

	if instance.Memory != nil {
		_ = d.Set("memory", instance.Memory)
	}

	if instance.Storage != nil {
		_ = d.Set("storage", instance.Storage)
	}

	if instance.NodeCount != nil {
		_ = d.Set("node_count", instance.NodeCount)
	}

	if instance.IsTmp != nil {
		_ = d.Set("is_tmp", instance.IsTmp)
	}

	if instance.ExclusterId != nil {
		_ = d.Set("excluster_id", instance.ExclusterId)
	}

	if instance.Pid != nil {
		_ = d.Set("pid", instance.Pid)
	}

	if instance.Qps != nil {
		_ = d.Set("qps", instance.Qps)
	}

	if instance.Paymode != nil {
		_ = d.Set("paymode", instance.Paymode)
	}

	if instance.Locker != nil {
		_ = d.Set("locker", instance.Locker)
	}

	if instance.StatusDesc != nil {
		_ = d.Set("status_desc", instance.StatusDesc)
	}

	if instance.WanStatus != nil {
		_ = d.Set("wan_status", instance.WanStatus)
	}

	if instance.IsAuditSupported != nil {
		_ = d.Set("is_audit_supported", instance.IsAuditSupported)
	}

	if instance.Machine != nil {
		_ = d.Set("machine", instance.Machine)
	}

	if instance.IsEncryptSupported != nil {
		_ = d.Set("is_encrypt_supported", instance.IsEncryptSupported)
	}

	if instance.Cpu != nil {
		_ = d.Set("cpu", instance.Cpu)
	}

	if instance.Ipv6Flag != nil {
		_ = d.Set("ipv6_flag", instance.Ipv6Flag)
	}

	if instance.Vipv6 != nil {
		_ = d.Set("vipv6", instance.Vipv6)
	}

	if instance.WanVipv6 != nil {
		_ = d.Set("wan_vipv6", instance.WanVipv6)
	}

	if instance.WanPortIpv6 != nil {
		_ = d.Set("wan_port_ipv6", instance.WanPortIpv6)
	}

	if instance.WanStatusIpv6 != nil {
		_ = d.Set("wan_status_ipv6", instance.WanStatusIpv6)
	}

	if instance.DbEngine != nil {
		_ = d.Set("db_engine", instance.DbEngine)
	}

	if instance.DcnFlag != nil {
		_ = d.Set("dcn_flag", instance.DcnFlag)
	}

	if instance.DcnStatus != nil {
		_ = d.Set("dcn_status", instance.DcnStatus)
	}

	if instance.DcnDstNum != nil {
		_ = d.Set("dcn_dst_num", instance.DcnDstNum)
	}

	if instance.InstanceType != nil {
		_ = d.Set("instance_type", instance.InstanceType)
	}

	if instance.DbVersionId != nil {
		_ = d.Set("db_version_id", instance.DbVersionId)
	}

	client := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
	tagService := svctag.NewTagService(client)
	region := client.Region
	tags, err := tagService.DescribeResourceTags(ctx, "mariadb", "instance", region, instanceId)
	if err != nil {
		return err
	}
	_ = d.Set("tags", tags)

	DbInstance, err := service.DescribeMariadbDbInstanceDetail(ctx, instanceId)
	if err != nil {
		return err
	}

	if DbInstance == nil {
		log.Printf("[WARN]%s resource `tencentcloud_mariadb_instance` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	var zones []*string
	if DbInstance.MasterZone != nil {
		zones = append(zones, DbInstance.MasterZone)
	}

	if DbInstance.SlaveZones != nil {
		zones = append(zones, DbInstance.SlaveZones...)
	}

	_ = d.Set("zones", zones)

	return nil
}

func resourceTencentCloudMariadbInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mariadb_instance.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service    = MariadbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		request    = mariadb.NewModifyDBInstanceNameRequest()
		instanceId = d.Id()
	)

	request.InstanceId = &instanceId

	immutableArgs := []string{"zones", "node_count", "memory", "storage", "period", "count", "auto_voucher", "voucher_ids", "vpc_id", "subnet_id", "db_version_id", "security_group_ids", "auto_renew_flag", "ipv6_flag", "init_params", "dcn_region", "dcn_instance_id", "total_count", "instances"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("instance_name") {
		if v, ok := d.GetOk("instance_name"); ok {
			request.InstanceName = helper.String(v.(string))
		}

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMariadbClient().ModifyDBInstanceName(request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			return nil
		})

		if err != nil {
			log.Printf("[CRITAL]%s update mariadb instance failed, reason:%+v", logId, err)
			return err
		}
	}

	if d.HasChange("tags") {
		tcClient := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
		tagService := svctag.NewTagService(tcClient)
		oldTags, newTags := d.GetChange("tags")
		replaceTags, deleteTags := svctag.DiffTags(oldTags.(map[string]interface{}), newTags.(map[string]interface{}))
		resourceName := tccommon.BuildTagResourceName("mariadb", "instance", tcClient.Region, d.Id())
		if err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags); err != nil {
			return err
		}
	}

	if d.HasChange("project_id") {
		if v, ok := d.GetOkExists("project_id"); ok {
			projectId := int64(v.(int))
			MPRequest := mariadb.NewModifyDBInstancesProjectRequest()
			MPRequest.InstanceIds = common.StringPtrs([]string{instanceId})
			MPRequest.ProjectId = &projectId

			err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
				result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMariadbClient().ModifyDBInstancesProject(MPRequest)
				if e != nil {
					return tccommon.RetryError(e)
				} else {
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
				}

				return nil
			})

			if err != nil {
				log.Printf("[CRITAL]%s operate mariadb modifyInstanceProject failed, reason:%+v", logId, err)
				return err
			}
		}
	}

	if d.HasChange("vip") {
		if v, ok := d.GetOk("vip"); ok {
			Vip := v.(string)
			var VipFlowId int64
			VipRequest := mariadb.NewModifyInstanceNetworkRequest()
			VipRequest.InstanceId = &instanceId
			VipRequest.Vip = &Vip
			if v, ok := d.GetOk("vpc_id"); ok {
				VipRequest.VpcId = helper.String(v.(string))
			}

			if v, ok := d.GetOk("subnet_id"); ok {
				VipRequest.SubnetId = helper.String(v.(string))
			}

			err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
				result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMariadbClient().ModifyInstanceNetwork(VipRequest)
				if e != nil {
					return tccommon.RetryError(e)
				} else {
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
				}

				VipFlowId = *result.Response.FlowId
				return nil
			})

			if err != nil {
				log.Printf("[CRITAL]%s operate mariadb network failed, reason:%+v", logId, err)
				return err
			}

			// wait
			if VipFlowId != NONE_FLOW_TASK {
				err = resource.Retry(10*tccommon.WriteRetryTimeout, func() *resource.RetryError {
					result, e := service.DescribeFlowById(ctx, VipFlowId)
					if e != nil {
						return tccommon.RetryError(e)
					}

					if *result.Status == MARIADB_TASK_SUCCESS {
						return nil
					} else if *result.Status == MARIADB_TASK_RUNNING {
						return resource.RetryableError(fmt.Errorf("operate mariadb network status is running"))
					} else if *result.Status == MARIADB_TASK_FAIL {
						return resource.NonRetryableError(fmt.Errorf("operate mariadb network status is fail"))
					} else {
						e = fmt.Errorf("operate mariadb network status illegal")
						return resource.NonRetryableError(e)
					}
				})

				if err != nil {
					log.Printf("[CRITAL]%s operate mariadb network task failed, reason:%+v", logId, err)
					return err
				}
			}
		}
	}

	return resourceTencentCloudMariadbInstanceRead(d, meta)
}

func resourceTencentCloudMariadbInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mariadb_instance.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := MariadbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	instanceId := d.Id()

	if err := service.IsolateDBInstanceById(ctx, instanceId); err != nil {
		return err
	}

	err := resource.Retry(2*tccommon.ReadRetryTimeout, func() *resource.RetryError {
		instance, e := service.DescribeMariadbInstanceById(ctx, instanceId)
		if e != nil {
			return resource.NonRetryableError(e)
		}
		if instance == nil {
			return nil
		}
		if *instance.Status == 2 {
			return resource.RetryableError(fmt.Errorf("isolate mariadb status is %v,start retrying ...", *instance.Status))
		}
		if *instance.Status == -1 {
			return nil
		}
		err := fmt.Errorf("isolate mariadb status is %v,we won't wait for it finish", *instance.Status)
		return resource.NonRetryableError(err)
	})

	if err != nil {
		log.Printf("[CRITAL]%s isolate mariadb fail, reason:%s\n ", logId, err.Error())
		return err
	}

	if err := service.DeleteMariadbInstanceById(ctx, instanceId); err != nil {
		return err
	}

	err = resource.Retry(2*tccommon.ReadRetryTimeout, func() *resource.RetryError {
		instance, e := service.DescribeMariadbInstanceById(ctx, instanceId)
		if e != nil {
			return resource.NonRetryableError(e)
		}
		if instance == nil {
			return nil
		}

		if *instance.Status == -1 {
			return resource.RetryableError(fmt.Errorf("delete mariadb status is %v,start retrying ...", *instance.Status))
		}

		err := fmt.Errorf("delete mariadb status is %v,we won't wait for it finish", *instance.Status)
		return resource.NonRetryableError(err)
	})

	if err != nil {
		log.Printf("[CRITAL]%s delete mariadb fail, reason:%s\n ", logId, err.Error())
		return err
	}

	return nil
}
