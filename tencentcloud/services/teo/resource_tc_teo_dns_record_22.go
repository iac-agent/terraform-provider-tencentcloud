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

func ResourceTencentCloudTeoDnsRecord22() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTeoDnsRecord22Create,
		Read:   resourceTencentCloudTeoDnsRecord22Read,
		Update: resourceTencentCloudTeoDnsRecord22Update,
		Delete: resourceTencentCloudTeoDnsRecord22Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Site ID. Cannot be null or empty string.",
			},

			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "DNS record name. If the domain name is in Chinese, Korean, or Japanese, it needs to be converted to punycode before input.",
			},

			"type": {
				Type:     schema.TypeString,
				Required: true,
				Description: "DNS record type. Valid values: <li>A: Points the domain to an external IPv4 address, e.g., 8.8.8.8;</li><li>AAAA: Points the domain to an external IPv6 address;</li><li>MX: Used for mail servers. Lower priority values are preferred when multiple MX records exist;</li><li>CNAME: Points the domain to another domain, which then resolves to the final IP address;</li><li>TXT: Identifies and describes the domain, commonly used for domain verification and SPF records (anti-spam);</li><li>NS: Required when delegating subdomain resolution to other DNS providers. NS records cannot be added to the root domain;</li><li>CAA: Specifies the CA that can issue certificates for this site;</li><li>SRV: Identifies a server using a specific service, commonly found in Microsoft system directory management.</li>\n" +
					"Different record types (e.g., SRV, CAA) have different requirements for host record names and record value formats. For detailed descriptions and format examples of each record type, please refer to: [DNS Record Type Introduction](https://cloud.tencent.com/document/product/1552/90453#2f681022-91ab-4a9e-ac3d-0a6c454d954e).",
			},

			"content": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "DNS record content. Fill in the corresponding content based on the Type value. If the domain name is in Chinese, Korean, or Japanese, it needs to be converted to punycode before input.",
			},

			"location": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "DNS record resolution line. Defaults to Default, which means the default resolution line that takes effect for all regions.\n\n- Resolution line configuration only applies when Type (DNS record type) is A, AAAA, or CNAME.\n- Resolution line configuration only applies to Standard and Enterprise editions. For valid values, please refer to: [Resolution Line and Corresponding Code Enumeration](https://cloud.tencent.com/document/product/1552/112542).",
			},

			"ttl": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Cache time. The user-specified value range is 60~86400. The smaller the value, the faster the record modification takes effect in each region. Default is 300, in seconds.",
			},

			"weight": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "DNS record weight. The user-specified value range is -1~100. Setting to 0 means no resolution. Default is -1, which means no weight is set. Weight configuration only applies when Type (DNS record type) is A, AAAA, or CNAME.<br>Note: Under the same subdomain, different DNS records with the same resolution line should either all have weights set or all have no weights set.",
			},

			"priority": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "MX record priority. This parameter only takes effect when Type (DNS record type) is MX. The smaller the value, the higher the priority. The user-specified value range is 0~50. Default is 0.",
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
				Description: "DNS record resolution status, the following values:\n" +
					"	- enable: has taken effect;\n" +
					"	- disable: has been disabled.",
			},
			"created_on": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation time.",
			},
			"modified_on": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Modify time.",
			},
			"record_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "DNS record ID.",
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
							Description: "Site ID.<br>Note: ZoneId is only used as an output parameter. It cannot be used as an input parameter in ModifyDnsRecords. If this parameter is provided, it will be ignored.",
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
							Description: "DNS record type. Valid values:\n<li>A: Points the domain to an external IPv4 address, e.g., 8.8.8.8;</li>\n<li>AAAA: Points the domain to an external IPv6 address;</li>\n<li>MX: Used for mail servers. Lower priority values are preferred when multiple MX records exist;</li>\n<li>CNAME: Points the domain to another domain, which then resolves to the final IP address;</li>\n<li>TXT: Identifies and describes the domain, commonly used for domain verification and SPF records (anti-spam);</li>\n<li>NS: Required when delegating subdomain resolution to other DNS providers. NS records cannot be added to the root domain;</li>\n<li>CAA: Specifies the CA that can issue certificates for this site;</li>\n<li>SRV: Identifies a server using a specific service, commonly found in Microsoft system directory management.</li>.",
						},
						"location": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "DNS record resolution line. Defaults to Default, which means the default resolution line that takes effect for all regions.<br>Resolution line configuration only applies when Type (DNS record type) is A, AAAA, or CNAME.<br>For valid values, please refer to: [Resolution Line and Corresponding Code Enumeration](https://cloud.tencent.com/document/product/1552/112542).",
						},
						"content": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "DNS record content. Fill in the corresponding content based on the Type value.",
						},
						"ttl": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Cache time. Value range is 60~86400. The smaller the value, the faster the record modification takes effect in each region, in seconds.",
						},
						"weight": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "DNS record weight. Value range is -1~100. -1 means no weight is assigned, 0 means no resolution. Weight configuration only applies when Type (DNS record type) is A, AAAA, or CNAME.",
						},
						"priority": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "MX record priority. Value range is 0~50. The smaller the value, the higher the priority.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "DNS record resolution status. Valid values: <li>enable: has taken effect;</li><li>disable: has been disabled.</li>Note: Status is only used as an output parameter. It cannot be used as an input parameter in ModifyDnsRecords. If this parameter is provided, it will be ignored.",
						},
						"created_on": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Creation time.<br>Note: CreatedOn is only used as an output parameter. It cannot be used as an input parameter in ModifyDnsRecords. If this parameter is provided, it will be ignored.",
						},
						"modified_on": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Modification time.<br>Note: ModifiedOn is only used as an output parameter. It cannot be used as an input parameter in ModifyDnsRecords. If this parameter is provided, it will be ignored.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudTeoDnsRecord22Create(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_dns_record_22.create")()
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
		log.Printf("[CRITAL]%s create teo dns record 22 failed, reason:%+v", logId, err)
		return err
	}

	if response.Response == nil || response.Response.RecordId == nil {
		return fmt.Errorf("create teo dns record 22 failed, RecordId is nil")
	}

	recordId = *response.Response.RecordId

	d.SetId(strings.Join([]string{zoneId, recordId}, tccommon.FILED_SP))

	return resourceTencentCloudTeoDnsRecord22Read(d, meta)
}

func resourceTencentCloudTeoDnsRecord22Read(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_dns_record_22.read")()
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

	respData, err := service.DescribeTeoDnsRecord22ById(ctx, zoneId, recordId)
	if err != nil {
		return err
	}
	if respData == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `teo_dns_record_22` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
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

	if respData.CreatedOn != nil {
		_ = d.Set("created_on", respData.CreatedOn)
	}

	if respData.ModifiedOn != nil {
		_ = d.Set("modified_on", respData.ModifiedOn)
	}

	if respData.RecordId != nil {
		_ = d.Set("record_id", respData.RecordId)
	}

	dnsRecords := make([]map[string]interface{}, 0)
	dnsRecord := map[string]interface{}{}
	if respData.ZoneId != nil {
		dnsRecord["zone_id"] = respData.ZoneId
	}
	if respData.RecordId != nil {
		dnsRecord["record_id"] = respData.RecordId
	}
	if respData.Name != nil {
		dnsRecord["name"] = respData.Name
	}
	if respData.Type != nil {
		dnsRecord["type"] = respData.Type
	}
	if respData.Location != nil {
		dnsRecord["location"] = respData.Location
	}
	if respData.Content != nil {
		dnsRecord["content"] = respData.Content
	}
	if respData.TTL != nil {
		dnsRecord["ttl"] = respData.TTL
	}
	if respData.Weight != nil {
		dnsRecord["weight"] = respData.Weight
	}
	if respData.Priority != nil {
		dnsRecord["priority"] = respData.Priority
	}
	if respData.Status != nil {
		dnsRecord["status"] = respData.Status
	}
	if respData.CreatedOn != nil {
		dnsRecord["created_on"] = respData.CreatedOn
	}
	if respData.ModifiedOn != nil {
		dnsRecord["modified_on"] = respData.ModifiedOn
	}
	dnsRecords = append(dnsRecords, dnsRecord)
	_ = d.Set("dns_records", dnsRecords)

	_ = recordId
	return nil
}

func resourceTencentCloudTeoDnsRecord22Update(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_dns_record_22.update")()
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
			log.Printf("[CRITAL]%s update teo dns record 22 failed, reason:%+v", logId, err)
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
			log.Printf("[CRITAL]%s update teo dns record 22 failed, reason:%+v", logId, err)
			return err
		}
	}

	_ = zoneId
	_ = recordId
	return resourceTencentCloudTeoDnsRecord22Read(d, meta)
}

func resourceTencentCloudTeoDnsRecord22Delete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_dns_record_22.delete")()
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
		log.Printf("[CRITAL]%s delete teo dns record 22 failed, reason:%+v", logId, err)
		return err
	}

	_ = response
	_ = recordId
	return nil
}
