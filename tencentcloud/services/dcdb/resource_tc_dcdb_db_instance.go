package dcdb

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dcdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dcdb/v20180411"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudDcdbDbInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDcdbDbInstanceCreate,
		Read:   resourceTencentCloudDcdbDbInstanceRead,
		Update: resourceTencentCloudDcdbDbInstanceUpdate,
		Delete: resourceTencentCloudDcdbDbInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"zones": {
				Required: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "availability 可用区 distribution 的 分片 nodes 可以 是 filled 使用 up 到 two availability zones. 当 分片 规格 是 一个 master 和 two slaves，two 的 nodes 是 在 first availability 可用区Note 该 当前 availability 可用区 该 可以 是 sold needs 到 是 pulled through DescribeDCDBSaleInfo interface。",
			},

			"period": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "长度 的 时间 您 want 到 buy，单位: month。",
			},

			"shard_memory": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Shard 内存 大小，单位: GB，可以 pass DescribeShardSpec Query 实例 规格 到 obtain。",
			},

			"shard_storage": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Shard 存储 大小，单位: GB，可以 pass DescribeShardSpec Query 实例 规格 到 obtain。",
			},

			"shard_node_count": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "数量 单个 分片 nodes，可以 pass DescribeShardSpec Query 实例 规格 到 obtain。",
			},

			"shard_count": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "数量 实例 fragments， 可选 范围 是 2-8，和 new fragments 可以 是 added 到 最大 的 64 fragments 通过 upgrading 实例。",
			},

			// "count": {
			// 	Optional:    true,
			// 	Type:        schema.TypeInt,
			// 	Description: "数量 的 实例 到 是 purchased.",
			// },

			"project_id": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "项目 ID，其中 可以 是 获取 通过 viewing 项目 列表，如果未传入，它 将 是 associated 使用 默认值 项目。",
			},

			"vpc_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Virtual 私有 网络 ID，如果未传入 或 passed 空，它 表示 该 它 是 创建 作为 basic 网络。",
			},

			"subnet_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Virtual 私有 网络 子网 ID，必填 当 VpcId 是 不 空。",
			},

			"db_version_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Database 引擎 版本，currently 可用: 8.0.18，10.1.9，5.7.17. 8.0.18 - MySQL 8.0.18; 10.1.9 - Mariadb 10.1.9; 5.7.17 - Percona 5.7.17 如果未填写， 默认为 5.7.17，其中 表示 Percona 5.7.17。",
			},

			"auto_voucher": {
				Optional:    true,
				Type:        schema.TypeBool,
				Description: "是否automatically 使用 vouchers 对于 payment，不 使用 通过 默认值。",
			},

			"voucher_ids": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Voucher ID 列表，currently 仅 支持 specifying 一个 voucher。",
			},

			"instance_name": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "实例名称，您 可以 集合 名称 实例 independently through 此 字段。",
			},

			"ipv6_flag": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "是否support IPv6。",
			},

			"extranet_access": {
				Optional:    true,
				Type:        schema.TypeBool,
				Description: "是否open extranet 访问。",
			},

			"vip": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeString,
				Description: "字段 为必填项 到 指定VIP。",
			},

			"vipv6": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeString,
				Description: "字段 为必填项 到 指定VIPv6。",
			},

			"vport": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Intranet 端口",
			},

			"resource_tags": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "数组 标签键-值 pairs。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tag_key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "键 的 标签",
						},
						"tag_value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "值 的 标签",
						},
					},
				},
			},

			"init_params": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "参数 列表. 可选 值 的 此 interface 是: character_set_server (character 集合，必须 是 passed)， lower_case_table_names (表 名称 是 case sensitive，必须 是 passed，0 - sensitive; 1 - insensitive)， innodb_page_size (innodb 数据 页面，默认值 16K)， sync_mode ( Synchronous 模式: 0 - asynchronous; 1 - strong synchronous; 2 - strong synchronous degenerate. 默认为 strong synchronous degenerate) 。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"param": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "名称 参数。",
						},
						"value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "值 的 参数。",
						},
					},
				},
			},

			"dcn_region": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "DCN 来源 地域",
			},

			"dcn_instance_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "DCN 来源 实例 ID。",
			},

			"auto_renew_flag": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Automatic renewal flag，0 表示 默认值 state ( 用户 has 不 集合 它，该 是， initial state 是 manual renewal，和 用户 has activated prepaid non-stop privilege 和 将 also perform automatic renewal). 1 表示 automatic renewal，2 表示 无 automatic renewal (用户 setting). 如果 business has 无 concept 的 renewal 或 automatic renewal 不是必填项，它 needs 到 是 集合 到 0。",
			},

			"security_group_ids": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Security 组 ids， 安全 组 可以 是 passed 在 form 的 数组，compatible 使用 previous SecurityGroupId 参数。",
			},
		},
	}
}

func resourceTencentCloudDcdbDbInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dcdb_db_instance.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	var (
		request       = dcdb.NewCreateDCDBInstanceRequest()
		response      = dcdb.NewCreateDCDBInstanceResponse()
		instanceId    string
		dcnInstanceId string
		vpcId         string
		subnetId      string
		ipv6Flag      int
		service       = DcdbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)
	if v, ok := d.GetOk("zones"); ok {
		zonesSet := v.(*schema.Set).List()
		request.Zones = helper.InterfacesStringsPoint(zonesSet)
	}

	if v, _ := d.GetOk("period"); v != nil {
		request.Period = helper.IntInt64(v.(int))
	}

	if v, _ := d.GetOk("shard_memory"); v != nil {
		request.ShardMemory = helper.IntInt64(v.(int))
	}

	if v, _ := d.GetOk("shard_storage"); v != nil {
		request.ShardStorage = helper.IntInt64(v.(int))
	}

	if v, _ := d.GetOk("shard_node_count"); v != nil {
		request.ShardNodeCount = helper.IntInt64(v.(int))
	}

	if v, _ := d.GetOk("shard_count"); v != nil {
		request.ShardCount = helper.IntInt64(v.(int))
	}

	if v, _ := d.GetOk("project_id"); v != nil {
		request.ProjectId = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("vpc_id"); ok {
		request.VpcId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("subnet_id"); ok {
		request.SubnetId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("db_version_id"); ok {
		request.DbVersionId = helper.String(v.(string))
	}

	if v, _ := d.GetOk("auto_voucher"); v != nil {
		request.AutoVoucher = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOk("voucher_ids"); ok {
		voucherIdsSet := v.(*schema.Set).List()
		request.VoucherIds = helper.InterfacesStringsPoint(voucherIdsSet)
	}

	if v, ok := d.GetOk("instance_name"); ok {
		request.InstanceName = helper.String(v.(string))
	}

	if v, _ := d.GetOk("ipv6_flag"); v != nil {
		request.Ipv6Flag = helper.IntInt64(v.(int))
		ipv6Flag = v.(int)
	}

	if v, ok := d.GetOk("resource_tags"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			resourceTag := dcdb.ResourceTag{}
			if v, ok := dMap["tag_key"]; ok {
				resourceTag.TagKey = helper.String(v.(string))
			}
			if v, ok := dMap["tag_value"]; ok {
				resourceTag.TagValue = helper.String(v.(string))
			}
			request.ResourceTags = append(request.ResourceTags, &resourceTag)
		}
	}

	if v, ok := d.GetOk("init_params"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			dBParamValue := dcdb.DBParamValue{}
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
		dcnInstanceId = v.(string)
	}

	if v, _ := d.GetOk("auto_renew_flag"); v != nil {
		request.AutoRenewFlag = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("security_group_ids"); ok {
		securityGroupIdsSet := v.(*schema.Set).List()
		request.SecurityGroupIds = helper.InterfacesStringsPoint(securityGroupIdsSet)
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDcdbClient().CreateDCDBInstance(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create dcdb dbInstance failed, reason:%+v", logId, err)
		return err
	}

	if response == nil || len(response.Response.InstanceIds) < 1 {
		d.SetId("")
		return fmt.Errorf("[CRITAL]%s create dcdb dbInstance failed.", d.Id())
	}

	instanceId = *response.Response.InstanceIds[0]
	d.SetId(instanceId)

	if len(request.InitParams) < 1 {
		defaultInitParams := []*dcdb.DBParamValue{
			{
				Param: helper.String("character_set_server"),
				Value: helper.String("utf8mb4"),
			},
			{
				Param: helper.String("lower_case_table_names"),
				Value: helper.String("1"),
			},
			{
				Param: helper.String("sync_mode"),
				Value: helper.String("2"),
			},
			{
				Param: helper.String("innodb_page_size"),
				Value: helper.String("16384"),
			},
		}
		request.InitParams = defaultInitParams
	}

	initRet, flowId, e := service.InitDcdbDbInstance(ctx, instanceId, request.InitParams)
	if e != nil {
		return e
	}
	if !initRet {
		return fmt.Errorf("db instance init failed")
	}

	if flowId != nil {
		// need to wait init operation success
		// 0:success; 1:failed, 2:running
		conf := tccommon.BuildStateChangeConf([]string{}, []string{"0"}, 3*tccommon.ReadRetryTimeout, time.Second, service.DcdbDbInstanceStateRefreshFunc(helper.UInt64Int64(*flowId), []string{}))
		if _, e := conf.WaitForState(); e != nil {
			return e
		}
	}

	if dcnInstanceId != "" {
		// need to wait dcn init processing complete
		// 0:none; 1:creating, 2:running
		conf := tccommon.BuildStateChangeConf([]string{}, []string{"2"}, 3*tccommon.ReadRetryTimeout, time.Second, service.DcdbDcnStateRefreshFunc(instanceId, []string{}))
		if _, e := conf.WaitForState(); e != nil {
			return e
		}
	}

	if v, ok := d.GetOkExists("extranet_access"); ok && v != nil {
		flag := v.(bool)
		err := service.SetDcdbExtranetAccess(ctx, instanceId, ipv6Flag, flag)
		if err != nil {
			return err
		}
	}

	var (
		vip   string
		vipv6 string
	)

	if v, ok := d.GetOk("vip"); ok {
		vip = v.(string)
	}
	if v, ok := d.GetOk("vipv6"); ok {
		vipv6 = v.(string)
	}

	if vip != "" || vipv6 != "" {
		if vpcId == "" || subnetId == "" {
			return fmt.Errorf("`vpc_id` and `subnet_id` cannot be empty when setting `vip` or `vipv6` fields!")
		}

		err := service.SetNetworkVip(ctx, instanceId, vpcId, subnetId, vip, vipv6)
		if err != nil {
			return err
		}
	}

	return resourceTencentCloudDcdbDbInstanceRead(d, meta)
}

func resourceTencentCloudDcdbDbInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dcdb_db_instance.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := DcdbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	instanceId := d.Id()

	ret, err := service.DescribeDcdbDbInstance(ctx, instanceId)
	if err != nil {
		return err
	}

	if ret == nil || len(ret.Instances) < 1 {
		d.SetId("")
		return fmt.Errorf("resource `DcdbDbInstance` %s does not exist", d.Id())
	}

	dbInstance := ret.Instances[0]

	// if dbInstance.Period != nil {
	// 	_ = d.Set("period", dbInstance.Period)
	// }

	if dbInstance.ShardDetail[0] != nil { // Memory and Storage is params for one shard
		shard := dbInstance.ShardDetail[0]

		if shard.Memory != nil {
			_ = d.Set("shard_memory", shard.Memory)
		}

		if shard.Storage != nil {
			_ = d.Set("shard_storage", shard.Storage)
		}
	}

	if dbInstance.NodeCount != nil {
		_ = d.Set("shard_node_count", dbInstance.NodeCount)
	}

	if dbInstance.ShardCount != nil {
		_ = d.Set("shard_count", dbInstance.ShardCount)
	}

	if dbInstance.ProjectId != nil {
		_ = d.Set("project_id", dbInstance.ProjectId)
	}

	if dbInstance.UniqueVpcId != nil {
		_ = d.Set("vpc_id", dbInstance.UniqueVpcId)
	}

	if dbInstance.UniqueSubnetId != nil {
		_ = d.Set("subnet_id", dbInstance.UniqueSubnetId)
	}

	if dbInstance.DbVersionId != nil {
		_ = d.Set("db_version_id", dbInstance.DbVersionId)
	}

	// if dbInstance.AutoVoucher != nil {
	// 	_ = d.Set("auto_voucher", dbInstance.AutoVoucher)
	// }

	// if dbInstance.VoucherIds != nil {
	// 	_ = d.Set("voucher_ids", dbInstance.VoucherIds)
	// }

	if dbInstance.InstanceName != nil {
		_ = d.Set("instance_name", dbInstance.InstanceName)
	}

	if dbInstance.Ipv6Flag != nil {
		_ = d.Set("ipv6_flag", dbInstance.Ipv6Flag)
	}

	if dbInstance.WanStatus != nil {
		//0-未开通；1-已开通；2-关闭；3-开通中
		if *dbInstance.WanStatus == DCDB_WAN_STATUS_UNOPEN || *dbInstance.WanStatus == DCDB_WAN_STATUS_CLOSED {
			_ = d.Set("extranet_access", false)
		}

		if *dbInstance.WanStatus == DCDB_WAN_STATUS_OPENED {
			_ = d.Set("extranet_access", true)
		}
	}

	if dbInstance.ResourceTags != nil {
		resourceTagsList := []interface{}{}
		for _, resourceTags := range dbInstance.ResourceTags {
			resourceTagsMap := map[string]interface{}{}

			if resourceTags.TagKey != nil {
				resourceTagsMap["tag_key"] = resourceTags.TagKey
			}

			if resourceTags.TagValue != nil {
				resourceTagsMap["tag_value"] = resourceTags.TagValue
			}

			resourceTagsList = append(resourceTagsList, resourceTagsMap)
		}

		_ = d.Set("resource_tags", resourceTagsList)

	}

	// if dbInstance.InitParams != nil {
	// 	initParamsList := []interface{}{}
	// 	for _, initParams := range dbInstance.InitParams {
	// 		initParamsMap := map[string]interface{}{}

	// 		if dbInstance.InitParams.Param != nil {
	// 			initParamsMap["param"] = dbInstance.InitParams.Param
	// 		}

	// 		if dbInstance.InitParams.Value != nil {
	// 			initParamsMap["value"] = dbInstance.InitParams.Value
	// 		}

	// 		initParamsList = append(initParamsList, initParamsMap)
	// 	}

	// 	_ = d.Set("init_params", initParamsList)

	// }

	if dbInstance.AutoRenewFlag != nil {
		_ = d.Set("auto_renew_flag", dbInstance.AutoRenewFlag)
	}

	if sg, err := service.DescribeDcdbSecurityGroup(ctx, instanceId); err == nil {
		sgIds := []*string{}
		for _, sg := range sg.Groups {
			sgIds = append(sgIds, sg.SecurityGroupId)
		}

		// fake sg
		var tmpSet []interface{}
		if v, ok := d.GetOk("security_group_ids"); ok {
			tmpSet = v.(*schema.Set).List()
			sgIds = helper.InterfacesStringsPoint(tmpSet)
		}
		// end

		_ = d.Set("security_group_ids", sgIds)
	} else {
		return err
	}

	// set dcn id and region
	if dcns, err := service.DescribeDcnDetailById(ctx, instanceId); err == nil {
		for _, dcn := range dcns {
			var master *dcdb.DcnDetailItem
			if *dcn.DcnFlag == DCDB_DCN_FLAG_MASTER {
				master = dcn
				_ = d.Set("dcn_region", master.Region)
				_ = d.Set("dcn_instance_id", master.InstanceId)
			}
		}
	} else {
		return err
	}

	// set vip, vipv6 and vport
	if detail, err := service.DescribeDcdbDbInstanceDetailById(ctx, instanceId); err == nil {
		if detail != nil {
			_ = d.Set("vip", detail.Vip)
			_ = d.Set("vipv6", detail.Vip6)
			_ = d.Set("vport", detail.Vport)

			if detail.MasterZone != nil {
				zones := []*string{detail.MasterZone}
				if detail.SlaveZones != nil {
					zones = append(zones, detail.SlaveZones...)
				}
				_ = d.Set("zones", zones)
			}
		}
	} else {
		return err
	}

	return nil
}

func resourceTencentCloudDcdbDbInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dcdb_db_instance.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	var (
		request = dcdb.NewModifyDBInstanceNameRequest()
		service = DcdbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	instanceId := d.Id()

	request.InstanceId = helper.String(instanceId)
	if d.HasChange("zones") {
		return fmt.Errorf("`zones` do not support change now.")
	}

	if d.HasChange("period") || d.HasChange("auto_voucher") || d.HasChange("voucher_ids") {
		if period, ok := d.GetOk("period"); ok {
			request := dcdb.NewRenewDCDBInstanceRequest()

			request.InstanceId = &instanceId
			request.Period = helper.IntInt64(period.(int))
			if v, _ := d.GetOk("auto_voucher"); v != nil {
				request.AutoVoucher = helper.Bool(v.(bool))
			}
			if v, ok := d.GetOk("voucher_ids"); ok {
				voucherIdsSet := v.(*schema.Set).List()
				for i := range voucherIdsSet {
					if voucherIdsSet[i] != nil {
						voucherIds := voucherIdsSet[i].(string)
						request.VoucherIds = append(request.VoucherIds, &voucherIds)
					}
				}
			}

			err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
				result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDcdbClient().RenewDCDBInstance(request)
				if e != nil {
					return tccommon.RetryError(e)
				} else {
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
				}
				return nil
			})
			if err != nil {
				log.Printf("[CRITAL]%s operate dcdb renewDCDBInstanceOperation failed, reason:%+v", logId, err)
				return err
			}
			_ = d.Set("period", period)
		}
		time.Sleep(2 * time.Second)
	}
	if d.HasChange("shard_memory") {
		return fmt.Errorf("`shard_memory` do not support change now.")
	}
	if d.HasChange("shard_storage") {
		return fmt.Errorf("`shard_storage` do not support change now.")
	}
	if d.HasChange("shard_node_count") {
		return fmt.Errorf("`shard_node_count` do not support change now.")
	}
	if d.HasChange("shard_count") {
		return fmt.Errorf("`shard_count` do not support change now.")
	}

	if v, ok := d.GetOkExists("extranet_access"); ok && v != nil {
		flag := v.(bool)
		var ipv6Flag int
		if v, _ := d.GetOk("ipv6_flag"); v != nil {
			ipv6Flag = v.(int)
		}
		err := service.SetDcdbExtranetAccess(ctx, instanceId, ipv6Flag, flag)
		if err != nil {
			return err
		}
		time.Sleep(2 * time.Second)
	}

	if d.HasChange("project_id") {
		if projectId, ok := d.GetOk("project_id"); ok {
			request := dcdb.NewModifyDBInstancesProjectRequest()

			request.InstanceIds = []*string{&instanceId}
			request.ProjectId = helper.IntInt64(projectId.(int))

			err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
				result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDcdbClient().ModifyDBInstancesProject(request)
				if e != nil {
					return tccommon.RetryError(e)
				} else {
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
				}
				return nil
			})
			if err != nil {
				log.Printf("[CRITAL]%s operate dcdb modifyInstanceProjectOperation failed, reason:%+v", logId, err)
				return err
			}
		}
		time.Sleep(2 * time.Second)
	}

	if d.HasChange("vpc_id") || d.HasChange("subnet_id") || d.HasChange("vip") || d.HasChange("vipv6") {
		var (
			vip      string
			vipv6    string
			vpcId    string
			subnetId string
		)
		if v, ok := d.GetOk("vip"); ok {
			vip = v.(string)
		}
		if v, ok := d.GetOk("vipv6"); ok {
			vipv6 = v.(string)
		}
		if v, ok := d.GetOk("vpc_id"); ok {
			vpcId = v.(string)
		}
		if v, ok := d.GetOk("subnet_id"); ok {
			subnetId = v.(string)
		}

		if vpcId == "" || subnetId == "" {
			return fmt.Errorf("`vpc_id` and `subnet_id` cannot be empty when updating network configs!")
		}

		err := service.SetNetworkVip(ctx, instanceId, vpcId, subnetId, vip, vipv6)
		if err != nil {
			return err
		}
	}

	if d.HasChange("db_version_id") {
		return fmt.Errorf("`db_version_id` do not support change now.")
	}

	if d.HasChange("ipv6_flag") {
		return fmt.Errorf("`ipv6_flag` do not support change now.")
	}
	if d.HasChange("resource_tags") {
		return fmt.Errorf("`resource_tags` do not support change now.")
	}
	if d.HasChange("init_params") {
		return fmt.Errorf("`init_params` do not support change now.")
	}
	if d.HasChange("dcn_region") {
		return fmt.Errorf("`dcn_region` do not support change now.")
	}
	if d.HasChange("dcn_instance_id") {
		return fmt.Errorf("`dcn_instance_id` do not support change now.")
	}
	if d.HasChange("auto_renew_flag") {
		return fmt.Errorf("`auto_renew_flag` do not support change now.")
	}
	if d.HasChange("security_group_ids") {
		return fmt.Errorf("`security_group_ids` do not support change now.")
	}
	if d.HasChange("instance_name") {
		if v, ok := d.GetOk("instance_name"); ok {
			request.InstanceName = helper.String(v.(string))
		}
		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDcdbClient().ModifyDBInstanceName(request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update dcdb dbInstance failed, reason:%+v", logId, err)
			return err
		}
	}
	return resourceTencentCloudDcdbDbInstanceRead(d, meta)
}

func resourceTencentCloudDcdbDbInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dcdb_db_instance.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := DcdbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	instanceId := d.Id()

	if err := service.DeleteDcdbDbInstanceById(ctx, instanceId); err != nil {
		return err
	}

	return nil
}
