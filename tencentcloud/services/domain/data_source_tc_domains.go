package domain

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	domain "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/domain/v20180808"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudDomains() *schema.Resource {
	return &schema.Resource{
		Read: datasourceTencentCloudDomainsRead,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"offset": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "指定data 偏移量 默认值：0。",
			},
			"limit": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     20,
				Description: "指定data 限制 在 范围 [1，100]. 默认值：20。",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "用于save response 作为 文件 locally。",
			},
			"list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "域名 结果 列表。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"auto_renew": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "是否domain auto renew，0 - manual renew，1 - auto renew。",
						},
						"is_premium": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "是否domain 是 premium。",
						},
						"domain_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "域名 ID。",
						},
						"expiration_date": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "域名 expiration date。",
						},
						"domain_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "域名 名称",
						},
						"code_tld": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "域名 代码 ltd。",
						},
						"creation_date": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "域名 创建时间。",
						},
						"tld": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "域名 ltd。",
						},
						"buy_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "域名 buy 状态",
						},
					},
				},
			},
		},
	}
}

func datasourceTencentCloudDomainsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("datasource.tencentcloud_domains.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	client := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
	service := DomainService{client}
	request := domain.NewDescribeDomainNameListRequest()

	if v, ok := d.GetOk("limit"); ok {
		request.Limit = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("offset"); ok {
		request.Offset = helper.IntUint64(v.(int))
	}

	result, err := service.DescribeDomainNameList(ctx, request)

	if err != nil {
		d.SetId("")
		return err
	}

	list := make([]interface{}, 0, len(result))
	ids := make([]string, 0)

	for i := range result {
		item := result[i]
		ids = append(ids, *item.DomainId)
		list = append(list, map[string]interface{}{
			"auto_renew":      item.AutoRenew,
			"is_premium":      item.IsPremium,
			"domain_id":       item.DomainId,
			"expiration_date": item.ExpirationDate,
			"domain_name":     item.DomainName,
			"code_tld":        item.CodeTld,
			"creation_date":   item.CreationDate,
			"tld":             item.Tld,
			"buy_status":      item.BuyStatus,
		})
	}

	d.SetId("domains-" + helper.DataResourceIdsHash(ids))
	if err := d.Set("list", list); err != nil {
		return err
	}

	if output, ok := d.GetOk("result_output_file"); ok {
		return tccommon.WriteToFile(output.(string), result)
	}

	return nil
}
