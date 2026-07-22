package sqlserver

import (
	"context"
	"fmt"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceTencentCloudSqlserverZoneConfig() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentSqlserverZoneConfigRead,
		Schema: map[string]*schema.Schema{
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used 到 store results.",
			},
			"zone_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A 列表 的 availability zones. Each element contains following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"availability_zone": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Alphabet ID 的 availability zone.",
						},
						"zone_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Number ID 的 availability zone.",
						},
						"specinfo_list": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "A 列表 的 specinfo configurations 对于 特定 availability zone. Each element contains following attributes:",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"spec_id": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "实例 规格 ID.",
									},
									"machine_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Model ID.",
									},
									"db_version": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Database 版本 信息. 有效 值: `2008R2 (SQL Server 2008 Enterprise)`, `2012SP3 (SQL Server 2012 Enterprise)`, `2016SP1 (SQL Server 2016 Enterprise)`, `201602 (SQL Server 2016 Standard)`, `2017 (SQL Server 2017 Enterprise)`.",
									},
									"db_version_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Version 名称 corresponding 到 `db_version` 字段.",
									},
									"memory": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Memory 大小 在 GB.",
									},
									"cpu": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Number 的 CPU cores.",
									},
									"min_storage_size": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Minimum 磁盘 大小 under 此 规格 在 GB.",
									},
									"max_storage_size": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Maximum 磁盘 大小 under 此 规格 在 GB.",
									},
									"qps": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "QPS 的 此 规格.",
									},
									"charge_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Billing 模式 under 此 规格. 有效 值 是 `POSTPAID_BY_HOUR`, `PREPAID` 和 `ALL`. `ALL` 表示 both POSTPAID_BY_HOUR 和 PREPAID.",
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

func dataSourceTencentSqlserverZoneConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencent_sqlserver_zone_config.read")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	sqlserverService := SqlserverService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	// get zoneinfo
	zoneInfoList, err := sqlserverService.DescribeZones(ctx)
	if err != nil {
		return fmt.Errorf("api[DescribeZones]fail, return %s", err.Error())
	}
	zoneSet := make(map[string]map[string]interface{})
	for _, zoneInfo := range zoneInfoList {
		zoneSetInfo := make(map[string]interface{}, 1)
		zoneSetInfo["id"] = zoneInfo.ZoneId
		zoneSet[*zoneInfo.Zone] = zoneSetInfo
	}

	var zoneList []interface{}
	for k, v := range zoneSet {
		var zoneListItem = make(map[string]interface{})
		zoneListItem["availability_zone"] = k
		zoneListItem["zone_id"] = v["id"]

		// get specinfo for each zone
		specinfoList, err := sqlserverService.DescribeProductConfig(ctx, k)
		if err != nil {
			return fmt.Errorf("api[DescribeProductConfig]fail, return %s", err.Error())
		}
		var specinfoConfigs []interface{}
		for _, specinfoItem := range specinfoList {
			var specinfoConfig = make(map[string]interface{})
			specinfoConfig["spec_id"] = specinfoItem.SpecId
			specinfoConfig["machine_type"] = specinfoItem.MachineType
			specinfoConfig["db_version"] = specinfoItem.Version
			specinfoConfig["db_version_name"] = specinfoItem.VersionName
			specinfoConfig["memory"] = specinfoItem.Memory
			specinfoConfig["cpu"] = specinfoItem.CPU
			specinfoConfig["min_storage_size"] = specinfoItem.MinStorage
			specinfoConfig["max_storage_size"] = specinfoItem.MaxStorage
			specinfoConfig["qps"] = specinfoItem.QPS
			specinfoConfig["charge_type"] = SQLSERVER_CHARGE_TYPE_NAME[*specinfoItem.PayModeStatus]

			specinfoConfigs = append(specinfoConfigs, specinfoConfig)
		}
		zoneListItem["specinfo_list"] = specinfoConfigs
		zoneList = append(zoneList, zoneListItem)
	}

	// set zone_list
	if err := d.Set("zone_list", zoneList); err != nil {
		return fmt.Errorf("[CRITAL]%s provider set zone_list fail, reason:%s\n ", logId, err.Error())
	}

	d.SetId("zone_config")

	if output, ok := d.GetOk("result_output_file"); ok && output.(string) != "" {
		if err := tccommon.WriteToFile(output.(string), zoneList); err != nil {
			log.Printf("[CRITAL]%s output file[%s] fail, reason[%s]\n",
				logId, output.(string), err.Error())
		}

	}
	return nil
}
