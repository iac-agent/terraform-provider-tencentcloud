package eb

import (
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	eb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/eb/v20210416"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudEbPutEvents() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudEbPutEventsCreate,
		Read:   resourceTencentCloudEbPutEventsRead,
		Delete: resourceTencentCloudEbPutEventsDelete,

		Schema: map[string]*schema.Schema{
			"event_list": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				Description: "事件 列表。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"source": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Event 来源 信息，new product 报告 必须 comply 使用 EB specifications。",
						},
						"data": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Event 数据， 内容 是 controlled 通过 系统 该 创建 事件， 当前 datacontenttype 仅 支持 应用/json;字符集=utf-8，so 此 字段 是 json 字符串。",
						},
						"type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Event 类型，customizable，可选 云 服务 writes COS:Created:PostObject 通过 默认值，使用: 到 separate 类型 字段。",
						},
						"subject": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Detailed 描述 事件 来源，customizable，可选 云 服务 默认为 standard qcs 资源 representation syntax: qcs::dts:ap-guangzhou:appid/uin:xxx。",
						},
						"time": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "时间戳 在 milliseconds 当 事件 occurred,时间.Now().UnixNano()/1e6。",
						},
					},
				},
			},

			"event_bus_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "事件 bus ID。",
			},
		},
	}
}

func resourceTencentCloudEbPutEventsCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_eb_put_events.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request    = eb.NewPutEventsRequest()
		eventBusId string
	)
	if v, ok := d.GetOk("event_list"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			event := eb.Event{}
			if v, ok := dMap["source"]; ok {
				event.Source = helper.String(v.(string))
			}
			if v, ok := dMap["data"]; ok {
				event.Data = helper.String(v.(string))
			}
			if v, ok := dMap["type"]; ok {
				event.Type = helper.String(v.(string))
			}
			if v, ok := dMap["subject"]; ok {
				event.Subject = helper.String(v.(string))
			}
			if v, ok := dMap["time"]; ok {
				event.Time = helper.IntInt64(v.(int))
			}
			request.EventList = append(request.EventList, &event)
		}
	}

	if v, ok := d.GetOk("event_bus_id"); ok {
		eventBusId = v.(string)
		request.EventBusId = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseEbClient().PutEvents(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate eb putEvents failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(eventBusId)

	return resourceTencentCloudEbPutEventsRead(d, meta)
}

func resourceTencentCloudEbPutEventsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_eb_put_events.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudEbPutEventsDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_eb_put_events.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
