package gs

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	gsv20191118 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/gs/v20191118"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudGsAndroidInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudGsAndroidInstancesRead,
		Schema: map[string]*schema.Schema{
			"android_instance_ids": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "列表 实例 IDs 到 查询. Up 到 100 per 请求。",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"android_instance_region": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "实例 地域 Aggregated 查询 across regions 是 不 currently 支持。",
			},

			"android_instance_zone": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "实例 availability 可用区",
			},

			"label_selector": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "实例 标签 selector。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "标签 键",
						},
						"operator": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "操作者 类型 IN: 标签值 必须 match 一个 的 Values; NOT_IN: 必须 不 match any; EXISTS: 标签 键 必须 exist; NOT_EXISTS: 标签 键 必须 不 exist。",
						},
						"values": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "标签值 列表. 必填 对于 IN 和 NOT_IN operators。",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},

			"filters": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Field filters. Supported 过滤器 names: 名称，UserId，HostSerialNumber，HostServerSerialNumber，AndroidInstanceModel。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "过滤字段名称",
						},
						"values": {
							Type:        schema.TypeList,
							Required:    true,
							Description: "过滤器 字段 值。",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},

			// Computed
			"android_instance_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "列表 Android 实例。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"android_instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "实例 ID",
						},
						"android_instance_region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "实例 地域",
						},
						"android_instance_zone": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "实例 availability 可用区",
						},
						"state": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "实例 state: INITIALIZING，NORMAL，PROCESSING。",
						},
						"android_instance_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "实例 规格。",
						},
						"android_instance_image_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "实例 镜像 ID。",
						},
						"width": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Resolution 宽度。",
						},
						"height": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Resolution 高度。",
						},
						"host_serial_number": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "主机 serial 数量。",
						},
						"android_instance_group_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "实例 组 ID",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "实例名称",
						},
						"user_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "用户 ID。",
						},
						"private_ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "内网 IP 地址",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "实例 创建时间。",
						},
						"host_server_serial_number": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Chassis serial 数量。",
						},
						"service_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Service 状态 IDLE: 不 connected; ESTABLISHED: connected。",
						},
						"android_instance_model": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Android 实例 model. YS1: basic; GC0/GC1/GC2: performance。",
						},
						"android_instance_labels": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "实例 标签 列表。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "标签 键",
									},
									"value": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "标签值",
									},
								},
							},
						},
					},
				},
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "用于保存结果。",
			},
		},
	}
}

func dataSourceTencentCloudGsAndroidInstancesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_gs_android_instances.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(nil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = GsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("android_instance_ids"); ok {
		ids := v.([]interface{})
		strIds := make([]*string, 0, len(ids))
		for _, id := range ids {
			strIds = append(strIds, helper.String(id.(string)))
		}
		paramMap["AndroidInstanceIds"] = strIds
	}

	if v, ok := d.GetOk("android_instance_region"); ok {
		paramMap["AndroidInstanceRegion"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("android_instance_zone"); ok {
		paramMap["AndroidInstanceZone"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("label_selector"); ok {
		labelSelectorList := v.([]interface{})
		labelReqs := make([]*gsv20191118.LabelRequirement, 0, len(labelSelectorList))
		for _, item := range labelSelectorList {
			m := item.(map[string]interface{})
			req := &gsv20191118.LabelRequirement{}
			if key, ok := m["key"].(string); ok && key != "" {
				req.Key = helper.String(key)
			}
			if op, ok := m["operator"].(string); ok && op != "" {
				req.Operator = helper.String(op)
			}
			if vals, ok := m["values"]; ok {
				for _, val := range vals.([]interface{}) {
					req.Values = append(req.Values, helper.String(val.(string)))
				}
			}
			labelReqs = append(labelReqs, req)
		}
		paramMap["LabelSelector"] = labelReqs
	}

	if v, ok := d.GetOk("filters"); ok {
		filterList := v.([]interface{})
		filters := make([]*gsv20191118.Filter, 0, len(filterList))
		for _, item := range filterList {
			m := item.(map[string]interface{})
			f := &gsv20191118.Filter{}
			if name, ok := m["name"].(string); ok && name != "" {
				f.Name = helper.String(name)
			}
			if vals, ok := m["values"]; ok {
				for _, val := range vals.([]interface{}) {
					f.Values = append(f.Values, helper.String(val.(string)))
				}
			}
			filters = append(filters, f)
		}
		paramMap["Filters"] = filters
	}

	var respData []*gsv20191118.AndroidInstance
	reqErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeGsAndroidInstancesByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		respData = result
		return nil
	})
	if reqErr != nil {
		return reqErr
	}

	instanceList := make([]map[string]interface{}, 0, len(respData))
	for _, inst := range respData {
		m := map[string]interface{}{}
		if inst.AndroidInstanceId != nil {
			m["android_instance_id"] = inst.AndroidInstanceId
		}
		if inst.AndroidInstanceRegion != nil {
			m["android_instance_region"] = inst.AndroidInstanceRegion
		}
		if inst.AndroidInstanceZone != nil {
			m["android_instance_zone"] = inst.AndroidInstanceZone
		}
		if inst.State != nil {
			m["state"] = inst.State
		}
		if inst.AndroidInstanceType != nil {
			m["android_instance_type"] = inst.AndroidInstanceType
		}
		if inst.AndroidInstanceImageId != nil {
			m["android_instance_image_id"] = inst.AndroidInstanceImageId
		}
		if inst.Width != nil {
			m["width"] = int(*inst.Width)
		}
		if inst.Height != nil {
			m["height"] = int(*inst.Height)
		}
		if inst.HostSerialNumber != nil {
			m["host_serial_number"] = inst.HostSerialNumber
		}
		if inst.Name != nil {
			m["name"] = inst.Name
		}
		if inst.UserId != nil {
			m["user_id"] = inst.UserId
		}
		if inst.PrivateIP != nil {
			m["private_ip"] = inst.PrivateIP
		}
		if inst.CreateTime != nil {
			m["create_time"] = inst.CreateTime
		}
		if inst.HostServerSerialNumber != nil {
			m["host_server_serial_number"] = inst.HostServerSerialNumber
		}
		if inst.ServiceStatus != nil {
			m["service_status"] = inst.ServiceStatus
		}
		if inst.AndroidInstanceModel != nil {
			m["android_instance_model"] = inst.AndroidInstanceModel
		}
		if len(inst.AndroidInstanceLabels) > 0 {
			labelList := make([]map[string]interface{}, 0, len(inst.AndroidInstanceLabels))
			for _, lbl := range inst.AndroidInstanceLabels {
				lblMap := map[string]interface{}{}
				if lbl.Key != nil {
					lblMap["key"] = lbl.Key
				}
				if lbl.Value != nil {
					lblMap["value"] = lbl.Value
				}
				labelList = append(labelList, lblMap)
			}
			m["android_instance_labels"] = labelList
		}
		instanceList = append(instanceList, m)
	}

	_ = d.Set("android_instance_list", instanceList)
	d.SetId(helper.BuildToken())

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), instanceList); e != nil {
			return e
		}
	}

	return nil
}
