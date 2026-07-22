package teo

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudTeoDnsRecord24() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTeoDnsRecord24Create,
		Read:   resourceTencentCloudTeoDnsRecord24Read,
		Update: resourceTencentCloudTeoDnsRecord24Update,
		Delete: resourceTencentCloudTeoDnsRecord24Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "可用区 ID",
			},

			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "DNS 记录 名称 如果 域名 名称 是 在 chinese，korean，或 japanese，它 needs 到 是 converted 到 punycode before input。",
			},

			"type": {
				Type:     schema.TypeString,
				Required: true,
				Description: "DNS 记录 类型. 有效 值 是:\n" +
					"	- A: points the domain name to an external ipv4 address, such as 8.8.8.8;\n" +
					"	- AAAA: points the domain name to an external ipv6 address;\n" +
					"	- MX: used for email servers. when there are multiple mx records, the lower the priority value, the higher the priority;\n" +
					"	- CNAME: points the domain name to another domain name, which then resolves to the final ip address;\n" +
					"	- TXT: identifies and describes the domain name, commonly used for domain verification and spf records (anti-spam);\n" +
					"	- NS: if you need to delegate the subdomain to another dns service provider for resolution, you need to add an ns record. the root domain cannot add ns records;\n" +
					"	- CAA: specifies the ca that can issue certificates for this site;\n" +
					"	- SRV: identifies a server using a service, commonly used in microsoft's directory management.\n" +
					"Different record types, such as SRV and CAA records, have different requirements for host record names and record value formats. for detailed descriptions and format examples of each record type, please refer to: [introduction to dns record types](https://intl.cloud.tencent.com/document/product/1552/90453?from_cn_redirect=1#2f681022-91ab-4a9e-ac3d-0a6c454d954e).",
			},

			"content": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "DNS 记录 内容 fill 在 corresponding 内容 according 到 类型 值 如果 域名 名称 是 在 chinese，korean，或 japanese，它 needs 到 是 converted 到 punycode before input。",
			},

			"location": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "DNS 记录 resolution 路由. 如果未指定， 默认为 DEFAULT，其中 表示 默认值 resolution 路由 和 是 effective 在 all regions.\n\n- resolution 路由 配置 是 仅 applicable 当 类型 (dns 记录 类型) 是 A，AAAA，或 CNAME.\n- resolution 路由 配置 是 仅 applicable 到 standard 版本 和 enterprise edition packages. 对于 有效 值，please refer 到: [resolution routes 和 corresponding 代码 enumeration](https://intl.云.tencent.com/document/product/1552/112542?from_cn_redirect=1)。",
			},

			"ttl": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Cache 时间. users 可以 指定a 值 范围 的 60-86400. smaller 值， faster modification records 将 take effect 在 all regions. 默认值：300. 单位: 秒。",
			},

			"weight": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "DNS 记录 权重 users 可以 指定a 值 范围 的 -1 到 100. 值 的 0 表示 无 resolution. 如果未指定， 默认为 -1，其中 表示 无 权重 是 集合. 权重 配置 是 仅 applicable 当 类型 (dns 记录 类型) 是 A，AAAA，或 CNAME. note: 对于 same subdomain，different dns records 使用 same resolution 路由 should either all have weights 集合 或 none have weights 集合。",
			},

			"priority": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "MX 记录 优先级，其中 takes effect 仅 当 类型 (dns 记录 类型) 是 MX. smaller 值， higher 优先级 users 可以 指定a 值 范围 的 0-50. 默认值为 0 如果未指定。",
			},

			"record_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "DNS 记录 ID。",
			},

			"status": {
				Type:     schema.TypeString,
				Computed: true,
				Description: "DNS 记录 resolution 状态, following 值:\n" +
					"	- enable: has taken effect;\n" +
					"	- disable: has been disabled.",
			},

			"created_on": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "创建时间。",
			},

			"modified_on": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "修改时间。",
			},

			"dns_records": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "DNS 记录 列表。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"zone_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "可用区 ID note: ZoneId 是 仅 使用 作为 output 参数，和 不能 是 使用 作为 input 参数 在 ModifyDnsRecords. 如果 此 参数 是 passed，它 将 是 ignored。",
						},
						"record_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "DNS 记录 ID。",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "DNS 记录 名称",
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
							Description: "DNS 记录 类型. 有效 值 是:\n" +
								"	- A: points the domain name to an external ipv4 address, such as 8.8.8.8;\n" +
								"	- AAAA: points the domain name to an external ipv6 address;\n" +
								"	- MX: used for email servers. when there are multiple mx records, the lower the priority value, the higher the priority;\n" +
								"	- CNAME: points the domain name to another domain name, which then resolves to the final ip address;\n" +
								"	- TXT: identifies and describes the domain name, commonly used for domain verification and spf records (anti-spam);\n" +
								"	- NS: if you need to delegate the subdomain to another dns service provider for resolution, you need to add an ns record. the root domain cannot add ns records;\n" +
								"	- CAA: specifies the ca that can issue certificates for this site;\n" +
								"	- SRV: identifies a server using a service, commonly used in microsoft's directory management.",
						},
						"location": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "DNS 记录 resolution 路由. 如果未指定， 默认为 DEFAULT，其中 表示 默认值 resolution 路由 和 是 effective 在 all regions. resolution 路由 配置 是 仅 applicable 当 类型 (dns 记录 类型) 是 A，AAAA，或 CNAME. 对于 有效 值，please refer 到: [resolution routes 和 corresponding 代码 enumeration](https://intl.云.tencent.com/document/product/1552/112542?from_cn_redirect=1)。",
						},
						"content": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "DNS 记录 内容 fill 在 corresponding 内容 according 到 类型 值",
						},
						"ttl": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Cache 时间. 值 范围 60-86400. smaller 值， faster modification records 将 take effect 在 all regions. 单位: 秒。",
						},
						"weight": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "DNS 记录 权重 值 范围 -1 到 100. -1 表示 无 权重 是 assigned，0 表示 无 resolution. 权重 配置 是 仅 applicable 当 类型 (dns 记录 类型) 是 A，AAAA，或 CNAME。",
						},
						"priority": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "MX 记录 优先级 值 范围 0-50. smaller 值， higher 优先级",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "DNS 记录 resolution 状态 有效值：启用: has taken effect; disable: has been 已禁用 note: 状态 是 仅 使用 作为 output 参数，和 不能 是 使用 作为 input 参数 在 ModifyDnsRecords. 如果 此 参数 是 passed，它 将 是 ignored。",
						},
						"created_on": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "创建时间. note: CreatedOn 是 仅 使用 作为 output 参数，和 不能 是 使用 作为 input 参数 在 ModifyDnsRecords. 如果 此 参数 是 passed，它 将 是 ignored。",
						},
						"modified_on": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "修改时间. note: ModifiedOn 是 仅 使用 作为 output 参数，和 不能 是 使用 作为 input 参数 在 ModifyDnsRecords. 如果 此 参数 是 passed，它 将 是 ignored。",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudTeoDnsRecord24Create(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_dns_record_24.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	var (
		zoneId   string
		recordId string
	)
	var (
		request  = teov20220901.NewCreateDnsRecordRequest()
		response = teov20220901.NewCreateDnsRecordResponse()
	)

	if v, ok := d.GetOk("zone_id"); ok {
		zoneId = v.(string)
		request.ZoneId = helper.String(zoneId)
	}

	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}

	if v, ok := d.GetOk("type"); ok {
		request.Type = helper.String(v.(string))
	}

	if v, ok := d.GetOk("content"); ok {
		request.Content = helper.String(v.(string))
	}

	if v, ok := d.GetOk("location"); ok {
		request.Location = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("ttl"); ok {
		request.TTL = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("weight"); ok {
		request.Weight = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("priority"); ok {
		request.Priority = helper.IntInt64(v.(int))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().CreateDnsRecordWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create teo dns record 24 failed, reason:%+v", logId, err)
		return err
	}

	if response.Response == nil || response.Response.RecordId == nil {
		return fmt.Errorf("create teo dns record 24 failed, response is nil")
	}

	recordId = *response.Response.RecordId

	d.SetId(strings.Join([]string{zoneId, recordId}, tccommon.FILED_SP))

	return resourceTencentCloudTeoDnsRecord24Read(d, meta)
}

func resourceTencentCloudTeoDnsRecord24Read(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_dns_record_24.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	service := TeoService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	zoneId := idSplit[0]
	recordId := idSplit[1]

	respData, err := service.DescribeTeoDnsRecordById(ctx, zoneId, recordId)
	if err != nil {
		return err
	}
	if respData == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `teo_dns_record_24` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if respData.ZoneId != nil {
		_ = d.Set("zone_id", respData.ZoneId)
	}

	if respData.RecordId != nil {
		_ = d.Set("record_id", respData.RecordId)
	}

	if respData.Name != nil {
		_ = d.Set("name", respData.Name)
	}

	if respData.Type != nil {
		_ = d.Set("type", respData.Type)
	}

	if respData.Location != nil {
		_ = d.Set("location", respData.Location)
	}

	if respData.Content != nil {
		_ = d.Set("content", respData.Content)
	}

	if respData.TTL != nil {
		_ = d.Set("ttl", respData.TTL)
	}

	if respData.Weight != nil {
		_ = d.Set("weight", respData.Weight)
	}

	if respData.Priority != nil {
		_ = d.Set("priority", respData.Priority)
	}

	if respData.Status != nil {
		_ = d.Set("status", respData.Status)
	}

	if respData.CreatedOn != nil {
		_ = d.Set("created_on", respData.CreatedOn)
	}

	if respData.ModifiedOn != nil {
		_ = d.Set("modified_on", respData.ModifiedOn)
	}

	_ = recordId
	return nil
}

func resourceTencentCloudTeoDnsRecord24Update(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_dns_record_24.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	zoneId := idSplit[0]
	recordId := idSplit[1]

	needChange := false
	mutableArgs := []string{"name", "type", "content", "location", "ttl", "weight", "priority"}
	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}

	if needChange {
		request := teov20220901.NewModifyDnsRecordsRequest()
		request.ZoneId = helper.String(zoneId)
		dnsRecord := &teov20220901.DnsRecord{
			RecordId: helper.String(recordId),
		}
		if v, ok := d.GetOk("name"); ok {
			dnsRecord.Name = helper.String(v.(string))
		}

		if v, ok := d.GetOk("type"); ok {
			dnsRecord.Type = helper.String(v.(string))
		}

		if v, ok := d.GetOk("content"); ok {
			dnsRecord.Content = helper.String(v.(string))
		}

		if v, ok := d.GetOk("location"); ok {
			dnsRecord.Location = helper.String(v.(string))
		}

		if v, ok := d.GetOkExists("ttl"); ok {
			dnsRecord.TTL = helper.IntInt64(v.(int))
		}

		if v, ok := d.GetOkExists("weight"); ok {
			dnsRecord.Weight = helper.IntInt64(v.(int))
		}

		if v, ok := d.GetOkExists("priority"); ok {
			dnsRecord.Priority = helper.IntInt64(v.(int))
		}
		request.DnsRecords = []*teov20220901.DnsRecord{dnsRecord}

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().ModifyDnsRecordsWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update teo dns record 24 failed, reason:%+v", logId, err)
			return err
		}
	}

	_ = zoneId
	_ = recordId
	return resourceTencentCloudTeoDnsRecord24Read(d, meta)
}

func resourceTencentCloudTeoDnsRecord24Delete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_dns_record_24.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	zoneId := idSplit[0]
	recordId := idSplit[1]

	var (
		request  = teov20220901.NewDeleteDnsRecordsRequest()
		response = teov20220901.NewDeleteDnsRecordsResponse()
	)

	request.ZoneId = helper.String(zoneId)
	request.RecordIds = helper.Strings([]string{recordId})

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().DeleteDnsRecordsWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s delete teo dns record 24 failed, reason:%+v", logId, err)
		return err
	}

	_ = response
	_ = recordId
	return nil
}
