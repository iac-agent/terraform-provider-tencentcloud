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

func ResourceTencentCloudTeoDnsRecordV5() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTeoDnsRecordV5Create,
		Read:   resourceTencentCloudTeoDnsRecordV5Read,
		Update: resourceTencentCloudTeoDnsRecordV5Update,
		Delete: resourceTencentCloudTeoDnsRecordV5Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Zone id.",
			},

			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "DNS record name. if the domain name is in chinese, korean, or japanese, it needs to be converted to punycode before input.",
			},

			"type": {
				Type:     schema.TypeString,
				Required: true,
				Description: "DNS record type. valid values are:\n" +
					"	- A: points the domain name to an external ipv4 address, such as 8.8.8.8;\n" +
					"	- AAAA: points the domain name to an external ipv6 address;\n" +
					"	- MX: used for email servers. when there are multiple mx records, the lower the priority value, the higher the priority;\n" +
					"	- CNAME: points the domain name to another domain name, which then resolves to the final ip address;\n" +
					"	- TXT: identifies and describes the domain name, commonly used for domain verification and spf records (anti-spam);\n" +
					"	- NS: if you need to delegate the subdomain to another dns service provider for resolution, you need to add an ns record. the root domain cannot add ns records;\n" +
					"	- CAA: specifies the ca that can issue certificates for this site;\n" +
					"	- SRV: identifies a server using a service, commonly used in microsoft's directory management.\n" +
					"Different record types, such as SRV and CAA records, have different requirements for host record names and record value formats.",
			},

			"content": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "DNS record content. fill in the corresponding content according to the type value. if the domain name is in chinese, korean, or japanese, it needs to be converted to punycode before input.",
			},

			"location": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "DNS record resolution route. if not specified, the default is DEFAULT, which means the default resolution route and is effective in all regions.",
			},

			"ttl": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Cache time. users can specify a value range of 60-86400. the smaller the value, the faster the modification records will take effect in all regions. default value: 300. unit: seconds.",
			},

			"weight": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "DNS record weight. users can specify a value range of -1 to 100. a value of 0 means no resolution. if not specified, the default is -1, which means no weight is set.",
			},

			"priority": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "MX record priority, which takes effect only when type (dns record type) is MX. the smaller the value, the higher the priority. users can specify a value range of 0-50. the default value is 0 if not specified.",
			},

			"status": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
				Description: "DNS record resolution status, the following values:\n" +
					"	- enable: has taken effect;\n" +
					"	- disable: has been disabled.",
			},

			"record_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "DNS record ID.",
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

func resourceTencentCloudTeoDnsRecordV5Create(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_dns_record_v5.create")()
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

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().CreateDnsRecordWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create teo dns record v5 failed, Response is nil."))
		}

		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create teo dns record v5 failed, reason:%+v", logId, err)
		return err
	}

	if response.Response.RecordId == nil {
		return fmt.Errorf("RecordId is nil.")
	}

	recordId = *response.Response.RecordId
	d.SetId(strings.Join([]string{zoneId, recordId}, tccommon.FILED_SP))

	return resourceTencentCloudTeoDnsRecordV5Read(d, meta)
}

func resourceTencentCloudTeoDnsRecordV5Read(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_dns_record_v5.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = TeoService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	zoneId := d.Get("zone_id").(string)
	if zoneId == "" {
		zoneId = idSplit[0]
	}
	recordId := idSplit[1]

	respData, err := service.DescribeTeoDnsRecordById(ctx, zoneId, recordId)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[WARN]%s resource `tencentcloud_teo_dns_record_v5` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
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

	if respData.RecordId != nil {
		_ = d.Set("record_id", respData.RecordId)
	}

	if respData.CreatedOn != nil {
		_ = d.Set("created_on", respData.CreatedOn)
	}

	if respData.ModifiedOn != nil {
		_ = d.Set("modified_on", respData.ModifiedOn)
	}

	return nil
}

func resourceTencentCloudTeoDnsRecordV5Update(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_dns_record_v5.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId = tccommon.GetLogId(tccommon.ContextNil)
		ctx   = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	zoneId := d.Get("zone_id").(string)
	if zoneId == "" {
		zoneId = idSplit[0]
	}
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
			log.Printf("[CRITAL]%s update teo dns record v5 failed, reason:%+v", logId, err)
			return err
		}
	}

	if d.HasChange("status") {
		request1 := teov20220901.NewModifyDnsRecordsStatusRequest()
		request1.ZoneId = helper.String(zoneId)

		if v, ok := d.GetOk("status"); ok {
			status := v.(string)
			if status == "enable" {
				request1.RecordsToEnable = helper.Strings([]string{recordId})
			}
			if status == "disable" {
				request1.RecordsToDisable = helper.Strings([]string{recordId})
			}
		}

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().ModifyDnsRecordsStatusWithContext(ctx, request1)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request1.GetAction(), request1.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update teo dns record v5 status failed, reason:%+v", logId, err)
			return err
		}
	}

	return resourceTencentCloudTeoDnsRecordV5Read(d, meta)
}

func resourceTencentCloudTeoDnsRecordV5Delete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_dns_record_v5.delete")()
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

	zoneId := d.Get("zone_id").(string)
	if zoneId == "" {
		zoneId = idSplit[0]
	}
	recordId := idSplit[1]

	request.ZoneId = helper.String(zoneId)
	request.RecordIds = helper.Strings([]string{recordId})

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
		log.Printf("[CRITAL]%s delete teo dns record v5 failed, reason:%+v", logId, err)
		return err
	}

	return nil
}
