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

func ResourceTencentCloudTeoDnsRecord50() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTeoDnsRecord50Create,
		Read:   resourceTencentCloudTeoDnsRecord50Read,
		Update: resourceTencentCloudTeoDnsRecord50Update,
		Delete: resourceTencentCloudTeoDnsRecord50Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Site ID.",
			},

			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "DNS record name. For domains in Chinese, Japanese or Korean, convert it to punycode before input.",
			},

			"type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "DNS record type. Valid values: `A`, `AAAA`, `MX`, `CNAME`, `TXT`, `NS`, `CAA`, `SRV`.",
			},

			"content": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "DNS record content, which should match the value of `type`.",
			},

			"location": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Resolution route, defaults to `Default`. Only applicable when `type` is `A`, `AAAA` or `CNAME`.",
			},

			"ttl": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Cache time in seconds, range 60-86400, defaults to 300.",
			},

			"weight": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Record weight, range -1-100. `-1` means no weight, `0` means no resolution. Only applicable when `type` is `A`, `AAAA` or `CNAME`.",
			},

			"priority": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "MX record priority, range 0-50. A smaller value indicates a higher priority. Only applicable when `type` is `MX`.",
			},

			"filters": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Query filter conditions of `DescribeDnsRecords` used on each read. The precise `id` filter of the resource itself is always sent first, and user filters are appended after it.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Filter field. Valid values: `id`, `name`, `content`, `type`, `ttl`.",
						},
						"values": {
							Type:        schema.TypeList,
							Required:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "Filter values, up to 20.",
						},
						"fuzzy": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether to enable fuzzy matching.",
						},
					},
				},
			},

			"sort_by": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Sort key of `DescribeDnsRecords`. Valid values: `content`, `created-on`, `name`, `ttl`, `type`. Defaults to a combined sort by `type` and `name`.",
			},

			"sort_order": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Sort order of `DescribeDnsRecords`. Valid values: `asc`, `desc`. Defaults to `asc`.",
			},

			"match": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Match mode of `DescribeDnsRecords`. Valid values: `all`, `any`. Defaults to `all`.",
			},

			"record_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "DNS record ID.",
			},

			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "DNS record resolution status. Valid values: `enable`, `disable`.",
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
		},
	}
}

func resourceTencentCloudTeoDnsRecord50Create(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_dns_record_50.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId    = tccommon.GetLogId(tccommon.ContextNil)
		ctx      = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request  = teov20220901.NewCreateDnsRecordRequest()
		response = teov20220901.NewCreateDnsRecordResponse()
		zoneId   string
		recordId string
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
		request.TTL = helper.Int64(int64(v.(int)))
	}

	if v, ok := d.GetOkExists("weight"); ok {
		request.Weight = helper.Int64(int64(v.(int)))
	}

	if v, ok := d.GetOkExists("priority"); ok {
		request.Priority = helper.Int64(int64(v.(int)))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().CreateDnsRecordWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil || result.Response.RecordId == nil || *result.Response.RecordId == "" {
			return resource.NonRetryableError(fmt.Errorf("[CRITAL]%s create teo dns_record_50 failed, record id is empty, request body [%s]", logId, request.ToJsonString()))
		}

		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create teo dns_record_50 failed, reason:%+v", logId, err)
		return err
	}

	recordId = *response.Response.RecordId
	d.SetId(strings.Join([]string{zoneId, recordId}, tccommon.FILED_SP))

	return resourceTencentCloudTeoDnsRecord50Read(d, meta)
}

func resourceTencentCloudTeoDnsRecord50Read(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_dns_record_50.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = teov20220901.NewDescribeDnsRecordsRequest()
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	zoneId := idSplit[0]
	recordId := idSplit[1]

	request.ZoneId = helper.String(zoneId)
	request.Limit = helper.Int64(1000)

	// Always locate the resource itself with a precise id filter first.
	request.Filters = append(request.Filters, &teov20220901.AdvancedFilter{
		Name:   helper.String("id"),
		Values: []*string{helper.String(recordId)},
	})

	// Append user-configured extra filters if any.
	if v, ok := d.GetOk("filters"); ok {
		for _, item := range v.([]interface{}) {
			filterMap := item.(map[string]interface{})
			advancedFilter := teov20220901.AdvancedFilter{}
			if name, ok := filterMap["name"].(string); ok && name != "" {
				advancedFilter.Name = helper.String(name)
			}
			if values, ok := filterMap["values"].([]interface{}); ok {
				valueList := make([]*string, 0, len(values))
				for _, value := range values {
					if s, ok := value.(string); ok {
						valueList = append(valueList, helper.String(s))
					}
				}
				advancedFilter.Values = valueList
			}
			if fuzzy, ok := filterMap["fuzzy"].(bool); ok {
				advancedFilter.Fuzzy = helper.Bool(fuzzy)
			}
			request.Filters = append(request.Filters, &advancedFilter)
		}
	}

	if v, ok := d.GetOk("sort_by"); ok {
		request.SortBy = helper.String(v.(string))
	}

	if v, ok := d.GetOk("sort_order"); ok {
		request.SortOrder = helper.String(v.(string))
	}

	if v, ok := d.GetOk("match"); ok {
		request.Match = helper.String(v.(string))
	}

	var dnsRecord *teov20220901.DnsRecord
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().DescribeDnsRecordsWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("[CRITAL]%s read teo dns_record_50 failed, response is nil, request body [%s]", logId, request.ToJsonString()))
		}

		if len(result.Response.DnsRecords) == 0 {
			log.Printf("[CRUD] read teo dns_record_50 id=%s", d.Id())
			return nil
		}

		dnsRecord = result.Response.DnsRecords[0]
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s read teo dns_record_50 failed, reason:%+v", logId, err)
		return err
	}

	if dnsRecord == nil {
		d.SetId("")
		return nil
	}

	_ = d.Set("zone_id", zoneId)

	if dnsRecord.RecordId != nil {
		_ = d.Set("record_id", dnsRecord.RecordId)
	}

	if dnsRecord.Name != nil {
		_ = d.Set("name", dnsRecord.Name)
	}

	if dnsRecord.Type != nil {
		_ = d.Set("type", dnsRecord.Type)
	}

	if dnsRecord.Location != nil {
		_ = d.Set("location", dnsRecord.Location)
	}

	if dnsRecord.Content != nil {
		_ = d.Set("content", dnsRecord.Content)
	}

	if dnsRecord.TTL != nil {
		_ = d.Set("ttl", int(*dnsRecord.TTL))
	}

	if dnsRecord.Weight != nil {
		_ = d.Set("weight", int(*dnsRecord.Weight))
	}

	if dnsRecord.Priority != nil {
		_ = d.Set("priority", int(*dnsRecord.Priority))
	}

	if dnsRecord.Status != nil {
		_ = d.Set("status", dnsRecord.Status)
	}

	if dnsRecord.CreatedOn != nil {
		_ = d.Set("created_on", dnsRecord.CreatedOn)
	}

	if dnsRecord.ModifiedOn != nil {
		_ = d.Set("modified_on", dnsRecord.ModifiedOn)
	}

	return nil
}

func resourceTencentCloudTeoDnsRecord50Update(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_dns_record_50.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId = tccommon.GetLogId(tccommon.ContextNil)
		ctx   = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
	)

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
			dnsRecord.TTL = helper.Int64(int64(v.(int)))
		}

		if v, ok := d.GetOkExists("weight"); ok {
			dnsRecord.Weight = helper.Int64(int64(v.(int)))
		}

		if v, ok := d.GetOkExists("priority"); ok {
			dnsRecord.Priority = helper.Int64(int64(v.(int)))
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
			log.Printf("[CRITAL]%s update teo dns_record_50 failed, reason:%+v", logId, err)
			return err
		}
	}

	return resourceTencentCloudTeoDnsRecord50Read(d, meta)
}

func resourceTencentCloudTeoDnsRecord50Delete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_dns_record_50.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = teov20220901.NewDeleteDnsRecordsRequest()
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	zoneId := idSplit[0]
	recordId := idSplit[1]

	request.ZoneId = helper.String(zoneId)
	request.RecordIds = []*string{helper.String(recordId)}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().DeleteDnsRecordsWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s delete teo dns_record_50 failed, reason:%+v", logId, err)
		return err
	}

	return nil
}
