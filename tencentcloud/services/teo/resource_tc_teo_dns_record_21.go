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

func ResourceTencentCloudTeoDnsRecord21() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTeoDnsRecord21Create,
		Read:   resourceTencentCloudTeoDnsRecord21Read,
		Update: resourceTencentCloudTeoDnsRecord21Update,
		Delete: resourceTencentCloudTeoDnsRecord21Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Site ID. This field is required.",
			},

			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "DNS record name. If the domain is in Chinese, Korean, or Japanese, convert it to punycode before input.",
			},

			"type": {
				Type:     schema.TypeString,
				Required: true,
				Description: "DNS record type. Valid values: A, AAAA, MX, CNAME, TXT, NS, CAA, SRV." +
					" Different record types have different requirements for the host record name and record value format. For details, see: [DNS Record Types](https://cloud.tencent.com/document/product/1552/90453#2f681022-91ab-4a9e-ac3d-0a6c454d954e).",
			},

			"content": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "DNS record content. Fill in the corresponding content according to the Type value. If the domain is in Chinese, Korean, or Japanese, convert it to punycode before input.",
			},

			"location": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "DNS record resolution line. Default is Default, representing all regions. Resolution line configuration only applies when Type is A, AAAA, or CNAME. For standard and enterprise editions only. For valid values, see: [Resolution Lines and Codes](https://cloud.tencent.com/document/product/1552/112542).",
			},

			"ttl": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Cache time in seconds. Valid range: 60-86400. Smaller values mean faster propagation. Default: 300.",
			},

			"weight": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "DNS record weight. Valid range: -1 to 100. 0 means no resolution, -1 means no weight. Weight configuration only applies when Type is A, AAAA, or CNAME. Note: DNS records under the same subdomain with the same resolution line should either all have weights set or all have weights unset.",
			},

			"priority": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "MX record priority. Only effective when Type is MX. Lower values indicate higher priority. Valid range: 0-50. Default: 0.",
			},

			"record_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "DNS record ID.",
			},

			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "DNS record resolution status. Valid values: enable (has taken effect), disable (has been disabled).",
			},

			"created_on": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation time.",
			},

			"modified_on": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Last modification time.",
			},

			"dns_records": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of DNS records.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"zone_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Site ID. Note: ZoneId is output-only in ModifyDnsRecords and will be ignored if passed as input.",
						},
						"record_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "DNS record ID.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "DNS record name.",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "DNS record type. Valid values: A, AAAA, MX, CNAME, TXT, NS, CAA, SRV.",
						},
						"location": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "DNS record resolution line. Default is Default, representing all regions. Resolution line configuration only applies when Type is A, AAAA, or CNAME. For valid values, see: [Resolution Lines and Codes](https://cloud.tencent.com/document/product/1552/112542).",
						},
						"content": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "DNS record content. Fill in the corresponding content according to the Type value.",
						},
						"ttl": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Cache time in seconds. Valid range: 60-86400. Smaller values mean faster propagation.",
						},
						"weight": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "DNS record weight. Valid range: -1 to 100. -1 means no weight, 0 means no resolution. Weight configuration only applies when Type is A, AAAA, or CNAME.",
						},
						"priority": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "MX record priority. Valid range: 0-50. Lower values indicate higher priority.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "DNS record resolution status. Valid values: enable (has taken effect), disable (has been disabled). Note: Status is output-only in ModifyDnsRecords and will be ignored if passed as input.",
						},
						"created_on": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Creation time. Note: CreatedOn is output-only in ModifyDnsRecords and will be ignored if passed as input.",
						},
						"modified_on": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Last modification time. Note: ModifiedOn is output-only in ModifyDnsRecords and will be ignored if passed as input.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudTeoDnsRecord21Create(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_dns_record_21.create")()
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
		log.Printf("[CRITAL]%s create teo dns record 21 failed, reason:%+v", logId, err)
		return err
	}

	if response.Response == nil || response.Response.RecordId == nil {
		return fmt.Errorf("create teo dns record 21 failed, response is nil")
	}

	recordId = *response.Response.RecordId

	d.SetId(strings.Join([]string{zoneId, recordId}, tccommon.FILED_SP))

	return resourceTencentCloudTeoDnsRecord21Read(d, meta)
}

func resourceTencentCloudTeoDnsRecord21Read(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_dns_record_21.read")()
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
		log.Printf("[WARN]%s resource `teo_dns_record_21` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
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

func resourceTencentCloudTeoDnsRecord21Update(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_dns_record_21.update")()
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
			log.Printf("[CRITAL]%s update teo dns record 21 failed, reason:%+v", logId, err)
			return err
		}
	}

	_ = zoneId
	_ = recordId
	return resourceTencentCloudTeoDnsRecord21Read(d, meta)
}

func resourceTencentCloudTeoDnsRecord21Delete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_dns_record_21.delete")()
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
		log.Printf("[CRITAL]%s delete teo dns record 21 failed, reason:%+v", logId, err)
		return err
	}

	_ = response
	_ = recordId
	return nil
}
