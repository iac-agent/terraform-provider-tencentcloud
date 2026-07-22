package dnspod

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	// dnspod "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dnspod/v20210323"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudDnspodDomainLogList() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudDnspodDomainLogListRead,
		Schema: map[string]*schema.Schema{
			"domain": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "域名",
			},

			"domain_id": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "域名 ID. 参数 DomainId has higher 优先级 比 参数 域名 如果 参数 DomainId 是 passed， 参数 域名 将 是 ignored. You 可以 find all Domains 和 DomainIds through DescribeDomainList interface。",
			},

			"log_list": {
				Computed: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "域名 Operation Log List. 注意：此字段可能返回 null，表示无法获取有效值。",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "用于保存结果。",
			},
		},
	}
}

func dataSourceTencentCloudDnspodDomainLogListRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_dnspod_domain_log_list.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	var domain string

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("domain"); ok {
		domain = v.(string)
		paramMap["Domain"] = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("domain_id"); ok {
		paramMap["DomainId"] = helper.IntUint64(v.(int))
	}

	service := DnspodService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var logList []*string

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeDnspodDomainLogListByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		logList = result
		return nil
	})
	if err != nil {
		return err
	}

	// ids := make([]string, 0, len(logList))
	if logList != nil {
		_ = d.Set("log_list", logList)
	}

	d.SetId(helper.DataResourceIdHash(domain))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), logList); e != nil {
			return e
		}
	}
	return nil
}
