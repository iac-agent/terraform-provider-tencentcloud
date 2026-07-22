package cdb

import (
	"context"
	"fmt"
	"log"
	"strconv"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdb/v20170320"
)

func TencentMysqlSellType() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"cdb_type": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "实例类型，可能的取值范围有：`UNIVERSAL`（通用类型）、`EXCLUSIVE`（独占类型）、`BASIC`（基本类型）、`BASIC_V2`（基本类型 v2）。",
		},
		"cpu": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "CPU 核心数。",
		},
		"mem_size": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "内存大小（以 MB 为单位）。",
		},
		"min_volume_size": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "最小磁盘大小（以 GB 为单位）。",
		},
		"max_volume_size": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "最大磁盘大小（以 GB 为单位）。",
		},
		"volume_step": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "磁盘增量（以 GB 为单位）。",
		},
		"qps": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "每秒查询数。",
		},
		"info": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "应用场景描述。",
		},
	}
}

func TencentMysqlZoneConfig() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "可用区域的名称，相当于特定数据中心。",
		},
		"is_default": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "指示当前DC是否为该区域的默认DC。可能的返回值：`0` - 否； `1` - 是的。",
		},
		"is_support_disaster_recovery": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "指示是否支持恢复：`0` - 否； `1` - 是的。",
		},
		"is_support_vpc": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "指示是否支持 VPC：`0` - 否； `1` - 是的。",
		},
		"engine_versions": {
			Type:        schema.TypeList,
			Elem:        &schema.Schema{Type: schema.TypeString},
			Computed:    true,
			Description: "要使用的数据库引擎的版本号。支持的版本包括`5.5`/`5.6`/`5.7`。",
		},
		"pay_type": {
			Type:        schema.TypeList,
			Elem:        &schema.Schema{Type: schema.TypeInt},
			Computed:    true,
			Description: "",
		},
		"hour_instance_sale_max_num": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "",
		},
		"support_slave_sync_modes": {
			Type:        schema.TypeList,
			Elem:        &schema.Schema{Type: schema.TypeInt},
			Computed:    true,
			Description: "数据复制模式。 `0` - 异步复制； `1` - 半同步复制； `2` - 强同步复制。",
		},
		"disaster_recovery_zones": {
			Type:        schema.TypeList,
			Elem:        &schema.Schema{Type: schema.TypeString},
			Computed:    true,
			Description: "有关可用恢复区域的信息。",
		},
		"slave_deploy_modes": {
			Type:        schema.TypeList,
			Elem:        &schema.Schema{Type: schema.TypeInt},
			Computed:    true,
			Description: "可用区部署方式。可用值：“0”-单个可用区； `1` - 多个可用区。",
		},
		"first_slave_zones": {
			Type:        schema.TypeList,
			Elem:        &schema.Schema{Type: schema.TypeString},
			Computed:    true,
			Description: "有关第一个从属实例的区域信息。",
		},
		"second_slave_zones": {
			Type:        schema.TypeList,
			Elem:        &schema.Schema{Type: schema.TypeString},
			Computed:    true,
			Description: "有关第二个从属实例的区域信息。",
		},
		"remote_ro_zones": {
			Type:        schema.TypeList,
			Elem:        &schema.Schema{Type: schema.TypeString},
			Computed:    true,
			Description: "有关远程 ro 实例的区域信息。",
		},
		"sells": {Type: schema.TypeList,
			Computed:    true,
			Description: "支持出售的实例类型列表：",
			Elem: &schema.Resource{
				Schema: TencentMysqlSellType(),
			},
		},
	}
}

func DataSourceTencentCloudMysqlZoneConfig() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentMysqlZoneConfigRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "地域 参数，用于标识要处理的数据所属的区域。",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "用于存储结果。",
			},
			// Computed values
			"list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "区域配置列表。每个元素包含以下属性：",
				Elem: &schema.Resource{
					Schema: TencentMysqlZoneConfig(),
				},
			},
		},
	}
}

func dataSourceTencentMysqlZoneConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_mysql_zone_config.read")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	mysqlService := MysqlService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	region := meta.(tccommon.ProviderMeta).GetAPIV3Conn().Region
	if regionInterface, ok := d.GetOk("region"); ok {
		region = regionInterface.(string)
	} else {
		log.Printf("[INFO]%s region is not set,so we use [%s] from env\n ", logId, region)
	}

	sellConfigures, err := mysqlService.DescribeDBZoneConfig(ctx)
	if err != nil {
		return fmt.Errorf("api[DescribeBackups]fail, return %s", err.Error())
	}
	var regionItem *cdb.CdbRegionSellConf

	for _, regionItem = range sellConfigures.Regions {
		if *regionItem.Region == region {
			break
		}
	}
	if regionItem == nil {
		return nil
	}

	var zoneConfigs []interface{}
	var zoneConfig = make(map[string]interface{})

	for _, sellItem := range regionItem.RegionConfig {
		zoneConfig["name"] = regionItem.RegionName
		if sellItem.HourInstanceSaleMaxNum != nil {
			zoneConfig["hour_instance_sale_max_num"] = *sellItem.HourInstanceSaleMaxNum
		}

		if sellItem.IsDefaultZone != nil {
			if *sellItem.IsDefaultZone {
				zoneConfig["is_default"] = 1
			} else {
				zoneConfig["is_default"] = 0
			}
		}

		if sellItem.IsSupportDr != nil {
			if *sellItem.IsSupportDr {
				zoneConfig["is_support_disaster_recovery"] = 1
			} else {
				zoneConfig["is_support_disaster_recovery"] = 0
			}
		}

		if sellItem.IsSupportVpc != nil {
			if *sellItem.IsSupportVpc {
				zoneConfig["is_support_vpc"] = 1
			} else {
				zoneConfig["is_support_vpc"] = 0
			}
		}

		payTypes := make([]int, len(sellItem.PayType))
		for index, strPtr := range sellItem.PayType {
			if tempInt, err := strconv.ParseInt(*strPtr, 10, 64); err != nil {
				errRet := fmt.Errorf("api[DescribeDBZoneConfig]return PayType error,not int")
				log.Printf("[CRITAL]%s %s\n ", logId, errRet.Error())
				return errRet
			} else {
				payTypes[index] = int(tempInt)
			}
		}
		zoneConfig["pay_type"] = payTypes

		supportSlaveSyncModes := make([]string, len(sellItem.ProtectMode))
		for index, intPtr := range sellItem.ProtectMode {
			supportSlaveSyncModes[index] = *intPtr
		}
		zoneConfig["support_slave_sync_modes"] = payTypes

		disasterRecoveryZones := make([]string, len(sellItem.DrZone))
		for index, strPtr := range sellItem.DrZone {
			disasterRecoveryZones[index] = *strPtr
		}
		zoneConfig["disaster_recovery_zones"] = disasterRecoveryZones

		var (
			slaveDeployModes                                 []int
			firstSlaveZones, secondSlaveZones, remoteRoZones []string
		)
		if sellItem.ZoneConf != nil {
			for _, mode := range sellItem.ZoneConf.DeployMode {
				slaveDeployModes = append(slaveDeployModes, int(*mode))
			}
			for _, zoneName := range sellItem.ZoneConf.SlaveZone {
				firstSlaveZones = append(firstSlaveZones, *zoneName)
			}
			for _, zoneName := range sellItem.ZoneConf.BackupZone {
				secondSlaveZones = append(secondSlaveZones, *zoneName)
			}
			for _, zoneName := range sellItem.RemoteRoZone {
				remoteRoZones = append(remoteRoZones, *zoneName)
			}
		}
		zoneConfig["slave_deploy_modes"] = slaveDeployModes
		zoneConfig["first_slave_zones"] = firstSlaveZones
		zoneConfig["second_slave_zones"] = secondSlaveZones
		zoneConfig["remote_ro_zones"] = remoteRoZones
	}

	var (
		engineVersions []string
		sells          []interface{}
	)

	for _, sellItem := range sellConfigures.Configs {
		if *sellItem.Status != ZONE_SELL_STATUS_ONLINE {
			continue
		}
		engineVersions = append(engineVersions, *sellItem.EngineType)

		var showConfigMap = make(map[string]interface{})
		if sellItem.DeviceType != nil {
			showConfigMap["cdb_type"] = *sellItem.DeviceType
		}
		if sellItem.Cpu != nil {
			showConfigMap["cpu"] = int(*sellItem.Cpu)
		}
		if sellItem.Memory != nil {
			showConfigMap["mem_size"] = int(*sellItem.Memory)
		}
		if sellItem.VolumeMax != nil {
			showConfigMap["max_volume_size"] = int(*sellItem.VolumeMax)
		}
		if sellItem.VolumeMin != nil {
			showConfigMap["min_volume_size"] = int(*sellItem.VolumeMin)
		}
		if sellItem.VolumeStep != nil {
			showConfigMap["volume_step"] = int(*sellItem.VolumeStep)
		}
		if sellItem.Iops != nil {
			showConfigMap["qps"] = int(*sellItem.Iops)
		}
		if sellItem.Info != nil {
			showConfigMap["info"] = *sellItem.Info
		}
		sells = append(sells, showConfigMap)

		zoneConfig["engine_versions"] = engineVersions
		zoneConfig["sells"] = sells

		zoneConfigs = append(zoneConfigs, zoneConfig)
	}

	if err := d.Set("list", zoneConfigs); err != nil {
		log.Printf("[CRITAL]%s provider set zoneConfigs fail, reason:%s\n ", logId, err.Error())
	}
	d.SetId("zoneconfig" + region)

	if output, ok := d.GetOk("result_output_file"); ok && output.(string) != "" {
		if err := tccommon.WriteToFile(output.(string), zoneConfigs); err != nil {
			log.Printf("[CRITAL]%s output file[%s] fail, reason[%s]\n",
				logId, output.(string), err.Error())
		}

	}
	return nil
}
