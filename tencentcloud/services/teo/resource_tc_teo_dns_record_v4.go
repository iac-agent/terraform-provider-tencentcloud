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

func ResourceTencentCloudTeoDnsRecordV4() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTeoDnsRecordV4Create,
		Read:   resourceTencentCloudTeoDnsRecordV4Read,
		Update: resourceTencentCloudTeoDnsRecordV4Update,
		Delete: resourceTencentCloudTeoDnsRecordV4Delete,
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
				Description: "DNS record resolution route. If not specified, the default is Default, which means the default resolution route and is effective in all regions.",
			},

			"ttl": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Cache time in seconds. Value range: 60-86400. Default: 300.",
			},

			"weight": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "DNS record weight. Value range: -1 to 100. -1 means no weight is set. 0 means no resolution.",
			},

			"priority": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "MX record priority, which takes effect only when type is MX. Value range: 0-50. Default: 0.",
			},

			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "DNS record resolution status. Valid values: enable (active), disable (disabled).",
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

func resourceTencentCloudTeoDnsRecordV4Create(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_dns_record_v4.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	var (
		zoneId   string
		recordId string
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
		log.Printf("[CRITAL]%s create teo dns record v4 failed, reason:%+v", logId, err)
		return err
	}

	if response.Response == nil || response.Response.RecordId == nil {
		return fmt.Errorf("Create teo dns record v4 failed, RecordId is nil.")
	}

	recordId = *response.Response.RecordId

	d.SetId(strings.Join([]string{zoneId, recordId}, tccommon.FILED_SP))

	return resourceTencentCloudTeoDnsRecordV4Read(d, meta)
}

func resourceTencentCloudTeoDnsRecordV4Read(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_dns_record_v4.read")()
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

	respData, err := service.DescribeTeoDnsRecordV4ById(ctx, zoneId, recordId)
	if err != nil {
		return err
	}

	if respData == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `teo_dns_record_v4` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
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

	_ = recordId
	return nil
}

func resourceTencentCloudTeoDnsRecordV4Update(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_dns_record_v4.update")()
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
			log.Printf("[CRITAL]%s update teo dns record v4 failed, reason:%+v", logId, err)
			return err
		}
	}

	needChange1 := false
	mutableArgs1 := []string{"status"}
	for _, v := range mutableArgs1 {
		if d.HasChange(v) {
			needChange1 = true
			break
		}
	}

	if needChange1 {
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
			log.Printf("[CRITAL]%s update teo dns record v4 status failed, reason:%+v", logId, err)
			return err
		}
	}

	_ = zoneId
	_ = recordId
	return resourceTencentCloudTeoDnsRecordV4Read(d, meta)
}

func resourceTencentCloudTeoDnsRecordV4Delete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_dns_record_v4.delete")()
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
		log.Printf("[CRITAL]%s delete teo dns record v4 failed, reason:%+v", logId, err)
		return err
	}

	_ = response
	_ = recordId
	return nil
}
