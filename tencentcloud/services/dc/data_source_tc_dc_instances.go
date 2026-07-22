package dc

import (
	"context"
	"crypto/md5"
	"fmt"
	"log"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceTencentCloudDcInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudDcInstancesRead,

		Schema: map[string]*schema.Schema{
			"dc_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID DC 到 是 queried。",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "名称 DC 到 是 queried。",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "用于保存结果。",
			},

			// Computed values
			"instance_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Information 列表 DC。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"dc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID DC。",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "名称 DC。",
						},
						"state": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "State 的 DC，和 可用 值 include `REJECTED`，`TOPAY`，`PAID`，`ALLOCATED`，`AVAILABLE`，`DELETING` 和 `DELETED`。",
						},
						"access_point_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Access point ID tne DC。",
						},
						"line_operator": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "操作者 的 DC，和 可用 值 include `ChinaTelecom`，`ChinaMobile`，`ChinaUnicom`，`In-houseWiring`，`ChinaOther` 和 `InternationalOperator`。",
						},
						"location": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "DC location 其中 连接 是 located。",
						},
						"bandwidth": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Bandwidth 的 DC。",
						},
						"port_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "端口 类型 DC 在 客户端，和 可用 值 include `100Base-T`，`1000Base-T`，`1000Base-LX`，`10GBase-T` 和 `10GBase-LR`. 默认值为 `1000Base-LX`。",
						},
						"circuit_code": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "circuit 代码 提供 通过 操作者 对于 DC。",
						},
						"redundant_dc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID redundant DC。",
						},
						"tencent_address": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Interconnect IP 的 DC within Tencent. 注意: 此 字段 可能 返回 null，indicating 该 无 有效 值 是 taken。",
						},
						"customer_address": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Interconnect IP 的 DC within 客户端. 注意: 此 字段 可能 返回 null，indicating 该 无 有效 值 是 taken。",
						},
						"customer_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Applicant 名称 DC， 默认为 获取 从 账号 注意: 此 字段 可能 返回 null，indicating 该 无 有效 值 是 taken。",
						},
						"customer_email": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Applicant email 的 DC， 默认为 获取 从 账号 注意: 此 字段 可能 返回 null，indicating 该 无 有效 值 是 taken。",
						},
						"customer_phone": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Applicant phone 数量 DC， 默认为 获取 从 账号 注意: 此 字段 可能 返回 null，indicating 该 无 有效 值 是 taken。",
						},
						"fault_report_contact_person": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Contact 的 报告 faulty. 注意: 此 字段 可能 返回 null，indicating 该 无 有效 值 是 taken。",
						},
						"fault_report_contact_phone": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Phone 数量 报告 faulty. 注意: 此 字段 可能 返回 null，indicating 该 无 有效 值 是 taken。",
						},
						"enabled_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Enable 时间 的 资源。",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "创建时间 的 资源。",
						},
						"expired_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Expire date 的 资源。",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudDcInstancesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_dc_instances.read")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := DcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var (
		id   = ""
		name = ""
	)
	if temp, ok := d.GetOk("dc_id"); ok {
		tempStr := temp.(string)
		if tempStr != "" {
			id = tempStr
		}
	}
	if temp, ok := d.GetOk("name"); ok {
		tempStr := temp.(string)
		if tempStr != "" {
			name = tempStr
		}
	}

	var infos, err = service.DescribeDirectConnects(ctx, id, name)

	if err != nil {
		return err
	}
	var instanceList = make([]map[string]interface{}, 0, len(infos))

	for _, item := range infos {

		var infoMap = make(map[string]interface{})
		infoMap["dc_id"] = *item.DirectConnectId
		infoMap["name"] = *item.DirectConnectName

		infoMap["state"] = strings.ToUpper(service.strPt2str(item.State))
		infoMap["access_point_id"] = service.strPt2str(item.AccessPointId)
		infoMap["line_operator"] = service.strPt2str(item.LineOperator)

		infoMap["location"] = service.strPt2str(item.Location)
		infoMap["bandwidth"] = service.int64Pt2int64(item.Bandwidth)
		infoMap["port_type"] = service.strPt2str(item.PortType)

		infoMap["circuit_code"] = service.strPt2str(item.CircuitCode)
		infoMap["redundant_dc_id"] = service.strPt2str(item.RedundantDirectConnectId)
		infoMap["tencent_address"] = service.strPt2str(item.TencentAddress)

		infoMap["customer_address"] = service.strPt2str(item.CustomerAddress)
		infoMap["customer_name"] = service.strPt2str(item.CustomerName)
		infoMap["customer_email"] = service.strPt2str(item.CustomerContactMail)

		infoMap["customer_phone"] = service.strPt2str(item.CustomerContactNumber)
		infoMap["fault_report_contact_person"] = service.strPt2str(item.FaultReportContactPerson)
		infoMap["fault_report_contact_phone"] = service.strPt2str(item.FaultReportContactNumber)

		infoMap["enabled_time"] = service.strPt2str(item.EnabledTime)
		infoMap["create_time"] = service.strPt2str(item.CreatedTime)
		infoMap["expired_time"] = service.strPt2str(item.ExpiredTime)

		instanceList = append(instanceList, infoMap)
	}

	if err := d.Set("instance_list", instanceList); err != nil {
		log.Printf("[CRITAL]%s provider set  dc instances fail, reason:%s\n ", logId, err.Error())
		return err
	}

	m := md5.New()
	_, err = m.Write([]byte(id + "_" + name))
	if err != nil {
		return err
	}
	d.SetId(fmt.Sprintf("%x", m.Sum(nil)))

	if output, ok := d.GetOk("result_output_file"); ok && output.(string) != "" {
		if err := tccommon.WriteToFile(output.(string), instanceList); err != nil {
			log.Printf("[CRITAL]%s output file[%s] fail, reason[%s]\n",
				logId, output.(string), err.Error())
			return err
		}
	}
	return nil
}
