package clb

import (
	"context"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	clb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/clb/v20180317"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudClbServerAttachments() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudClbServerAttachmentsRead,

		Schema: map[string]*schema.Schema{
			"clb_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "需要查询的CLB ID。",
			},
			"listener_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "需要查询的CLB监听ID。",
			},
			"rule_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "CLB监听规则ID。如果监听协议是`HTTP`/`HTTPS`，则此段为必填项。",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "用于保存结果。",
			},
			"attachment_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "云负载均衡器附件配置列表。每个元素包含以下属性：",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"clb_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "CLB的ID。",
						},
						"listener_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "CLB监听器ID。",
						},
						"protocol_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "侦听器中的协议类型，可用值包括“TCP”、“UDP”、“HTTP”、“HTTPS”和“TCP_SSL”。注意：`TCP_SSL`正在内部测试，如需使用请申请。",
						},
						"rule_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "CLB监听规则ID。",
						},
						"targets": {
							Type:        schema.TypeSet,
							Computed:    true,
							Description: "要附加的后端信息。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"instance_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "后端服务器的id。",
									},
									"eni_ip": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "弹性网卡唯一ID。",
									},
									"port": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "后端服务器的端口。",
									},
									"weight": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "后端服务的转发权重，范围[0，100]，默认为`10`。",
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

func dataSourceTencentCloudClbServerAttachmentsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_clb_attachments.read")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	params := make(map[string]string)
	clbId := d.Get("clb_id").(string)
	listenerId := d.Get("listener_id").(string)
	checkErr := ListenerIdCheck(listenerId)
	if checkErr != nil {
		return checkErr
	}
	locationId := ""
	if v, ok := d.GetOk("rule_id"); ok {
		locationId = v.(string)
		checkErr := RuleIdCheck(locationId)
		if checkErr != nil {
			return checkErr
		}
		params["rule_id"] = locationId
	}
	params["clb_id"] = clbId
	params["listener_id"] = listenerId

	clbService := ClbService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}
	var attachments []*clb.ListenerBackend
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		results, e := clbService.DescribeAttachmentsByFilter(ctx, params)
		if e != nil {
			return tccommon.RetryError(e)
		}
		attachments = results
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read CLB attachments failed, reason:%+v", logId, err)
		return err
	}
	attachmentList := make([]map[string]interface{}, 0, len(attachments))
	ids := make([]string, 0, len(attachments))
	for _, attachment := range attachments {
		mapping := map[string]interface{}{
			"clb_id":        clbId,
			"listener_id":   listenerId,
			"protocol_type": attachment.Protocol,
		}
		if locationId != "" {
			mapping["rule_id"] = locationId
		}
		if *attachment.Protocol == CLB_LISTENER_PROTOCOL_HTTP || *attachment.Protocol == CLB_LISTENER_PROTOCOL_HTTPS {
			if len(attachment.Rules) > 0 {
				for _, loc := range attachment.Rules {
					if locationId == "" || locationId == *loc.LocationId {
						mapping["targets"] = flattenBackendList(loc.Targets)
					}
				}
			}
		} else if *attachment.Protocol == CLB_LISTENER_PROTOCOL_TCP || *attachment.Protocol == CLB_LISTENER_PROTOCOL_UDP ||
			*attachment.Protocol == CLB_LISTENER_PROTOCOL_TCPSSL || *attachment.Protocol == CLB_LISTENER_PROTOCOL_QUIC {
			mapping["targets"] = flattenBackendList(attachment.Targets)
		}
		attachmentList = append(attachmentList, mapping)
		ids = append(ids, locationId+"#"+listenerId+"#"+clbId)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	if e := d.Set("attachment_list", attachmentList); e != nil {
		log.Printf("[CRITAL]%s provider set attachment list fail, reason:%+v", logId, e)
		return e
	}

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), attachmentList); e != nil {
			return e
		}
	}

	return nil
}
