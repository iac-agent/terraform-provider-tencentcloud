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

func ResourceTencentCloudTeoDnsRecordV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTeoDnsRecordV2Create,
		Read:   resourceTencentCloudTeoDnsRecordV2Read,
		Update: resourceTencentCloudTeoDnsRecordV2Update,
		Delete: resourceTencentCloudTeoDnsRecordV2Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Zone ID.",
			},

			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "DNS record name. If the domain name is in Chinese, Korean, or Japanese, it needs to be converted to punycode before input.",
			},

			"type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "DNS record type. Valid values: A, AAAA, MX, CNAME, TXT, NS, CAA, SRV.",
			},

			"content": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "DNS record content. Fill in the corresponding content according to the type value.",
			},

			"location": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "DNS record resolution route. If not specified, the default is Default.",
			},

			"ttl": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Cache time in seconds. Range: 60-86400. Default: 300.",
			},

			"weight": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "DNS record weight. Range: -1 to 100. Default: -1, which means no weight is set.",
			},

			"priority": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "MX record priority. Range: 0-50. Default: 0. Only takes effect when type is MX.",
			},

			"record_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "DNS record ID.",
			},

			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "DNS record resolution status. Valid values: enable, disable.",
			},

			"created_on": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation time.",
			},

			"modified_on": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Modification time.",
			},
		},
	}
}

func resourceTencentCloudTeoDnsRecordV2Create(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_dns_record_v2.create")()
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
		request.TTL = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("weight"); ok {
		request.Weight = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("priority"); ok {
		request.Priority = helper.IntInt64(v.(int))
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().CreateDnsRecordWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create teo dns record v2 failed, Response is nil."))
		}

		response = result
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create teo dns record v2 failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	if response.Response.RecordId == nil {
		return fmt.Errorf("RecordId is nil.")
	}

	recordId = *response.Response.RecordId
	_ = d.Set("record_id", recordId)
	d.SetId(strings.Join([]string{zoneId, recordId}, tccommon.FILED_SP))

	return resourceTencentCloudTeoDnsRecordV2Read(d, meta)
}

func resourceTencentCloudTeoDnsRecordV2Read(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_dns_record_v2.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = TeoService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	zoneId := d.Get("zone_id").(string)
	recordId := d.Get("record_id").(string)

	if zoneId == "" || recordId == "" {
		idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
		if len(idSplit) != 2 {
			return fmt.Errorf("id is broken,%s", d.Id())
		}
		zoneId = idSplit[0]
		recordId = idSplit[1]
	}

	respData, err := service.DescribeTeoDnsRecordV2ById(ctx, zoneId, recordId)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[WARN]%s resource `tencentcloud_teo_dns_record_v2` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	if respData.ZoneId != nil {
		_ = d.Set("zone_id", respData.ZoneId)
	}

	if respData.Name != nil {
		_ = d.Set("name", respData.Name)
	}

	if respData.Type != nil {
		_ = d.Set("type", respData.Type)
	}

	if respData.Content != nil {
		_ = d.Set("content", respData.Content)
	}

	if respData.Location != nil {
		_ = d.Set("location", respData.Location)
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

	if respData.RecordId != nil {
		_ = d.Set("record_id", respData.RecordId)
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

	return nil
}

func resourceTencentCloudTeoDnsRecordV2Update(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_dns_record_v2.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId = tccommon.GetLogId(tccommon.ContextNil)
		ctx   = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
	)

	zoneId := d.Get("zone_id").(string)
	recordId := d.Get("record_id").(string)

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

		reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().ModifyDnsRecordsWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			return nil
		})

		if reqErr != nil {
			log.Printf("[CRITAL]%s update teo dns record v2 failed, reason:%+v", logId, reqErr)
			return reqErr
		}
	}

	return resourceTencentCloudTeoDnsRecordV2Read(d, meta)
}

func resourceTencentCloudTeoDnsRecordV2Delete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_dns_record_v2.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = teov20220901.NewDeleteDnsRecordsRequest()
	)

	zoneId := d.Get("zone_id").(string)
	recordId := d.Get("record_id").(string)

	request.ZoneId = helper.String(zoneId)
	request.RecordIds = helper.Strings([]string{recordId})

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().DeleteDnsRecordsWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete teo dns record v2 failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return nil
}
