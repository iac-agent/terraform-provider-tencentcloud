package cdb

import (
	"fmt"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdb/v20170320"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudMysqlInstance() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudMysqlInstanceRead,
		Schema: map[string]*schema.Schema{
			"mysql_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "实例 ID，例如“cdb-c1nl9rpv”。它与数据库控制台页面中显示的实例 ID 相同。",
			},
			"instance_role": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue([]string{"master", "ro", "dr"}),
				Description: "实例类型。支持的值包括：“master”- 主实例、“dr”- 灾难恢复实例和“ro”- 只读实例。",
			},
			"status": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: tccommon.ValidateAllowedIntValue([]int{0, 1, 4, 5}),
				Description: "实例状态。可用值：“0”-创建； `1` - 运行； `4`-隔离； `5`-隔离。",
			},
			"security_group_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "实例的安全组ID。",
			},
			"pay_type": {
				Type:         schema.TypeInt,
				Optional:     true,
				Deprecated:   "It has been deprecated from version 1.36.0. Please use `charge_type` instead.",
				ValidateFunc: tccommon.ValidateAllowedIntValue([]int{0, 1}),
				Description: "实例的付费类型，‘0’：预付费，‘1’：后付费。",
			},
			"charge_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue([]string{MYSQL_CHARGE_TYPE_PREPAID, MYSQL_CHARGE_TYPE_POSTPAID}),
				Description: "实例的付费类型，有效值为`PREPAID`和`POSTPAID`。",
			},
			"instance_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "mysql 实例的名称。",
			},
			"engine_version": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue([]string{"5.1", "5.5", "5.6", "5.7", "8.0"}),
				Description: "要使用的数据库引擎的版本号。支持的版本包括5.5/5.6/5.7/8.0。",
			},
			"init_flag": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: tccommon.ValidateAllowedIntValue([]int{0, 1}),
				Description: "初始化标记。可用值： `0` - 未初始化； `1` - 已初始化。",
			},
			"with_dr": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: tccommon.ValidateAllowedIntValue([]int{0, 1}),
				Description: "是否查询容灾实例。",
			},
			"with_ro": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: tccommon.ValidateAllowedIntValue([]int{0, 1}),
				Description: "是否查询只读实例。",
			},
			"with_master": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: tccommon.ValidateAllowedIntValue([]int{0, 1}),
				Description: "是否查询master实例。",
			},
			"offset": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      0,
				ValidateFunc: tccommon.ValidateIntegerInRange(0, 1000),
				Description: "记录偏移量。默认值为 0。",
			},
			"limit": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      20,
				ValidateFunc: tccommon.ValidateIntegerInRange(1, 2000),
				Description: "单个请求返回的结果数。默认值为“20”，最大值为 2000。",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "用于存储结果。",
			},
			// Computed values
			"instance_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "实例列表。每个元素包含以下属性：",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"mysql_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "实例 ID，例如“cdb-c1nl9rpv”。它与数据库控制台页面中显示的实例 ID 相同。",
						},
						"instance_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "mysql 实例的名称。",
						},
						"instance_role": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "实例类型。支持的值包括：“master”- 主实例、“dr”- 灾难恢复实例和“ro”- 只读实例。",
						},
						"init_flag": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "初始化标记。可用值： `0` - 未初始化； `1` - 已初始化。",
						},
						"status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "实例状态。可用值：“0”-创建； `1` - 运行； `4`-隔离； `5`-隔离。",
						},
						"zone": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "可用区域的信息。",
						},
						"auto_renew_flag": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "自动更新标志。注意：仅支持预付费实例。",
						},
						"engine_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "要使用的数据库引擎的版本号。支持的版本包括`5.5`/`5.6`/`5.7`/`8.0`。",
						},
						"cpu_core_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "CPU 计数。",
						},
						"memory_size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "内存大小（以 MB 为单位）。",
						},
						"volume_size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "磁盘容量（以 GB 为单位）。",
						},
						"internet_status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "公共网络的状态。",
						},
						"internet_host": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "公网域名。",
						},
						"internet_port": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "公网端口。",
						},
						"intranet_ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "用于内部访问的实例IP。",
						},
						"intranet_port": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "用于内部用途的传输层端口号。",
						},
						"project_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "当前实例所属的项目ID。",
						},
						"vpc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "虚拟私有云ID。",
						},
						"subnet_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "当前实例所属子网ID。",
						},
						"slave_sync_mode": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "数据复制模式。 `0` - 异步复制； `1` - 半同步复制； `2` - 强同步复制。",
						},
						"device_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "支持的实例模型。 `HA`-高可用版本； `Basic` - 基本版本。",
						},
						"pay_type": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "实例的付费类型，‘0’：预付费，‘1’：后付费。",
						},
						"charge_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "实例付费类型。",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "创建实例的时间。",
						},
						"dead_line_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "实例的到期日期。注意：仅支持预付费实例。",
						},
						"master_instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "恢复实例的主实例ID。",
						},
						"ro_instance_ids": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "与当前实例关联的只读类型的ID列表。",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"ro_groups": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "只读实例组。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"group_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "组 ID，例如“cdbrg-pz7vg37p”。",
									},
									"instance_ids": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "与当前实例关联的只读类型的ID列表。",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},
						"dr_instance_ids": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "与当前实例关联的灾难恢复类型的ID列表。",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudMysqlInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_mysql_instance.read")()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	request := cdb.NewDescribeDBInstancesRequest()

	if mysqlId, ok := d.GetOk("mysql_id"); ok {
		mysqlIdValue := mysqlId.(string)
		request.InstanceIds = []*string{&mysqlIdValue}
	}
	if instanceType, ok := d.GetOk("instance_role"); ok {
		instanceTypeValue := instanceType.(string)
		var instanceRole uint64 = 1
		for k, v := range MYSQL_ROLE_MAP {
			if instanceTypeValue == v {
				instanceRole = uint64(k)
			}
		}
		request.InstanceTypes = []*uint64{&instanceRole}
	}
	if status, ok := d.GetOk("status"); ok {
		statusValue := uint64(status.(int))
		request.Status = []*uint64{&statusValue}
	}
	if securityGroupId, ok := d.GetOk("security_group_id"); ok {
		securityGroupIdValue := securityGroupId.(string)
		request.SecurityGroupId = &securityGroupIdValue
	}
	if payType, ok := d.GetOk("pay_type"); ok {
		payTypeValue := uint64(payType.(int))
		request.PayTypes = []*uint64{&payTypeValue}
	}
	if chargeType, ok := d.GetOk("charge_type"); ok {
		var payType int
		if chargeType == MYSQL_CHARGE_TYPE_PREPAID {
			payType = MysqlPayByMonth
		} else {
			payType = MysqlPayByUse
		}
		payTypeValue := uint64(payType)
		request.PayTypes = []*uint64{&payTypeValue}
	}
	if instanceName, ok := d.GetOk("instance_name"); ok {
		instanceNameValue := instanceName.(string)
		request.InstanceNames = []*string{&instanceNameValue}
	}
	if taskStatus, ok := d.GetOk("task_status"); ok {
		taskStatusValue := uint64(taskStatus.(int))
		request.TaskStatus = []*uint64{&taskStatusValue}
	}
	if engineVersion, ok := d.GetOk("engine_version"); ok {
		engineVersionValue := engineVersion.(string)
		request.EngineVersions = []*string{&engineVersionValue}
	}
	if initFlag, ok := d.GetOk("init_flag"); ok {
		initFlagValue := int64(initFlag.(int))
		request.InitFlag = &initFlagValue
	}
	if withDr, ok := d.GetOk("with_dr"); ok {
		withDrValue := int64(withDr.(int))
		request.WithDr = &withDrValue
	}
	if withRo, ok := d.GetOk("with_ro"); ok {
		withRoValue := int64(withRo.(int))
		request.WithRo = &withRoValue
	}
	if withMaster, ok := d.GetOk("with_master"); ok {
		withMasterValue := int64(withMaster.(int))
		request.WithMaster = &withMasterValue
	}
	offset := d.Get("offset")
	offsetValue := uint64(offset.(int))
	request.Offset = &offsetValue
	limit := d.Get("limit")
	limitValue := uint64(limit.(int))
	request.Limit = &limitValue

	client := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
	response, err := client.UseMysqlClient().DescribeDBInstances(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return fmt.Errorf("api[DescribeDBInstances]fail, return %s", err.Error())
	}

	instanceDetails := response.Response.Items
	instanceList := make([]map[string]interface{}, 0, len(instanceDetails))
	ids := make([]string, 0, len(instanceDetails))
	for _, item := range instanceDetails {
		mapping := map[string]interface{}{
			"mysql_id":        item.InstanceId,
			"instance_name":   item.InstanceName,
			"instance_role":   MYSQL_ROLE_MAP[*item.InstanceType],
			"init_flag":       item.InitFlag,
			"status":          item.Status,
			"zone":            item.Zone,
			"auto_renew_flag": item.AutoRenew,
			"engine_version":  item.EngineVersion,
			"cpu_core_count":  item.Cpu,
			"memory_size":     item.Memory,
			"volume_size":     item.Volume,
			"internet_status": item.WanStatus,
			"internet_host":   item.WanDomain,
			"internet_port":   item.WanPort,
			"intranet_ip":     item.Vip,
			"intranet_port":   item.Vport,
			"project_id":      item.ProjectId,
			"vpc_id":          item.UniqVpcId,
			"subnet_id":       item.UniqSubnetId,
			"slave_sync_mode": item.ProtectMode,
			"device_type":     item.DeviceType,
			"pay_type":        item.PayType,
			"create_time":     item.CreateTime,
			"dead_line_time":  item.DeadlineTime,
			"charge_type":     MYSQL_CHARGE_TYPE[int(*item.PayType)],
		}
		if item.MasterInfo != nil {
			mapping["master_instance_id"] = item.MasterInfo.InstanceId
		} else {
			mapping["master_instance_id"] = ""
		}
		if len(item.RoGroups) > 0 {
			roInstanceIds := make([]string, 0)
			roGroupList := make([]map[string]interface{}, 0, len(item.RoGroups))
			for _, roGroupInfo := range item.RoGroups {
				roInstanceId := make([]string, 0)
				for _, roInfo := range roGroupInfo.RoInstances {
					roInstanceId = append(roInstanceId, *roInfo.InstanceId)
				}
				roGroup := map[string]interface{}{
					"group_id":     *roGroupInfo.RoGroupId,
					"instance_ids": roInstanceId,
				}

				roInstanceIds = append(roInstanceIds, roInstanceId...)
				roGroupList = append(roGroupList, roGroup)
			}
			mapping["ro_instance_ids"] = roInstanceIds
			mapping["ro_groups"] = roGroupList
		}
		if len(item.DrInfo) > 0 {
			drInstanceIds := make([]string, 0)
			for _, drInfo := range item.DrInfo {
				drInstanceIds = append(drInstanceIds, *drInfo.InstanceId)
			}
			mapping["dr_instance_ids"] = drInstanceIds
		}

		ids = append(ids, *item.InstanceId)
		instanceList = append(instanceList, mapping)
	}
	d.SetId(helper.DataResourceIdsHash(ids))
	err = d.Set("instance_list", instanceList)
	if err != nil {
		log.Printf("[CRITAL]%s provider set instance list fail, reason:%s\n ", logId, err.Error())
	}
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		err = tccommon.WriteToFile(output.(string), instanceList)
		if err != nil {
			return err
		}
	}
	return nil
}
