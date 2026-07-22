package dnspod

import (
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dnspod "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dnspod/v20210323"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudDnspodDomainInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudDnspodDomainInstancesRead,
		Schema: map[string]*schema.Schema{
			"domain": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "域名",
			},

			"instance_list": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "域名 列表。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"domain": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "域名",
						},
						"group_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Group ID 的 域名",
						},
						"is_mark": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "是否Mark 域名",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "状态 域名",
						},
						"remark": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "备注 的 域名",
						},
						"id": {
							Computed:    true,
							Type:        schema.TypeString,
							Description: "ID 域名",
						},
						"domain_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "ID 域名",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "创建时间 的 域名",
						},
						"slave_dns": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Is secondary DNS 已启用",
						},
						"record_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "数量 DNS records under 此 域名",
						},
						"grade": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "DNS plan/包 grade 的 域名 (e.g.，DP_Free，DP_Plus)。",
						},
						"updated_on": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Last 修改时间 的 域名",
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

func dataSourceTencentCloudDnspodDomainInstancesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_dnspod_domain_instances.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	domain := d.Get("domain").(string)

	request := dnspod.NewDescribeDomainRequest()
	request.Domain = helper.String(domain)

	var response *dnspod.DescribeDomainResponse

	tmpList := make([]map[string]interface{}, 0)

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDnsPodClient().DescribeDomain(request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		response = result
		info := response.Response.DomainInfo
		domainMap := make(map[string]interface{})

		domainMap["id"] = domain
		domainMap["domain_id"] = info.DomainId
		domainMap["domain"] = domain
		domainMap["create_time"] = info.CreatedOn
		domainMap["is_mark"] = info.IsMark
		domainMap["slave_dns"] = info.SlaveDNS
		domainMap["record_count"] = info.RecordCount
		domainMap["grade"] = info.Grade
		domainMap["updated_on"] = info.UpdatedOn

		if info.Status != nil {
			if *info.Status == "pause" {
				domainMap["status"] = DNSPOD_DOMAIN_STATUS_DISABLE
			} else {
				domainMap["status"] = info.Status
			}
		}

		if info.Remark != nil {
			domainMap["remark"] = info.Remark
		}

		if info.GroupId != nil {
			domainMap["group_id"] = info.GroupId
		}

		tmpList = append(tmpList, domainMap)

		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read DnsPod Domain failed, reason:%s\n", logId, err.Error())
		return err
	}

	d.SetId(helper.DataResourceIdsHash([]string{domain}))
	_ = d.Set("instance_list", tmpList)

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
