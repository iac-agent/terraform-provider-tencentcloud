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
				Description: "Site ID.",
			},

			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "DNS record name. If the domain name is in Chinese, Korean, or Japanese, it needs to be converted to punycode before input.",
			},

			"type": {
				Type:     schema.TypeString,
				Required: true,
				Description: "DNS record type. Valid values are:\n" +
					"	- A: Points the domain name to an external IPv4 address, such as 8.8.8.8;\n" +
					"	- AAAA: Points the domain name to an external IPv6 address;\n" +
					"	- MX: Used for email servers. When there are multiple MX records, the lower the priority value, the higher the priority;\n" +
					"	- CNAME: Points the domain name to another domain name, which then resolves to the final IP address;\n" +
					"	- TXT: Identifies and describes the domain name, commonly used for domain verification and SPF records (anti-spam);\n" +
					"	- NS: If you need to delegate the subdomain to another DNS service provider for resolution, you need to add an NS record. The root domain cannot add NS records;\n" +
					"	- CAA: Specifies the CA that can issue certificates for this site;\n" +
					"	- SRV: Identifies a server using a service, commonly used in Microsoft's directory management.\n" +
					"Different record types, such as SRV and CAA records, have different requirements for host record names and record value formats. For detailed descriptions and format examples of each record type, please refer to: [Introduction to DNS Record Types](https://intl.cloud.tencent.com/document/product/1552/90453?from_cn_redirect=1#2f681022-91ab-4a9e-ac3d-0a6c454d954e).",
			},

			"content": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "DNS record content. Fill in the corresponding content according to the Type value. If the domain name is in Chinese, Korean, or Japanese, it needs to be converted to punycode before input.",
			},

			"location": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "DNS record resolution route. If not specified, the default is DEFAULT, which means the default resolution route and is effective in all regions.\n\n- Resolution route configuration is only applicable when Type (DNS record type) is A, AAAA, or CNAME.\n- Resolution route configuration is only applicable to Standard Edition and Enterprise Edition packages. For valid values, please refer to: [Resolution Routes and Corresponding Code Enumeration](https://intl.cloud.tencent.com/document/product/1552/112542?from_cn_redirect=1).",
			},

			"ttl": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Cache time. Users can specify a value range of 60-86400. The smaller the value, the faster the modification records will take effect in all regions. Default value: 300. Unit: seconds.",
			},

			"weight": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "DNS record weight. Users can specify a value range of -1 to 100. A value of 0 means no resolution. If not specified, the default is -1, which means no weight is set. Weight configuration is only applicable when Type (DNS record type) is A, AAAA, or CNAME. Note: For the same subdomain, different DNS records with the same resolution route should either all have weights set or none have weights set.",
			},

			"priority": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "MX record priority, which takes effect only when Type (DNS record type) is MX. The smaller the value, the higher the priority. Users can specify a value range of 0-50. The default value is 0 if not specified.",
			},

			"record_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "DNS record id.",
			},

			"status": {
				Type:     schema.TypeString,
				Computed: true,
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

			"dns_records": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "DNS record list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"zone_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Site ID. Note: ZoneId is only used as an output parameter and cannot be used as an input parameter in ModifyDnsRecords. If this parameter is passed, it will be ignored.",
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
							Type:     schema.TypeString,
							Computed: true,
							Description: "DNS record type. Valid values are:\n" +
								"	- A: Points the domain name to an external IPv4 address, such as 8.8.8.8;\n" +
								"	- AAAA: Points the domain name to an external IPv6 address;\n" +
								"	- MX: Used for email servers. When there are multiple MX records, the lower the priority value, the higher the priority;\n" +
								"	- CNAME: Points the domain name to another domain name, which then resolves to the final IP address;\n" +
								"	- TXT: Identifies and describes the domain name, commonly used for domain verification and SPF records (anti-spam);\n" +
								"	- NS: If you need to delegate the subdomain to another DNS service provider for resolution, you need to add an NS record. The root domain cannot add NS records;\n" +
								"	- CAA: Specifies the CA that can issue certificates for this site;\n" +
								"	- SRV: Identifies a server using a service, commonly used in Microsoft's directory management.",
						},
						"location": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "DNS record resolution route. If not specified, the default is DEFAULT, which means the default resolution route and is effective in all regions. Resolution route configuration is only applicable when Type (DNS record type) is A, AAAA, or CNAME. For valid values, please refer to: [Resolution Routes and Corresponding Code Enumeration](https://intl.cloud.tencent.com/document/product/1552/112542?from_cn_redirect=1).",
						},
						"content": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "DNS record content. Fill in the corresponding content according to the Type value.",
						},
						"ttl": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Cache time. Value range: 60-86400. The smaller the value, the faster the modification records will take effect in all regions. Unit: seconds.",
						},
						"weight": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "DNS record weight. Value range: -1 to 100. A value of -1 means no weight is assigned, and a value of 0 means no resolution. Weight configuration is only applicable when Type (DNS record type) is A, AAAA, or CNAME.",
						},
						"priority": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "MX record priority. Value range: 0-50. The smaller the value, the higher the priority.",
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
							Description: "DNS record resolution status. Valid values:\n" +
								"	- enable: has taken effect;\n" +
								"	- disable: has been disabled.\n" +
								"Note: Status is only used as an output parameter and cannot be used as an input parameter in ModifyDnsRecords. If this parameter is passed, it will be ignored.",
						},
						"created_on": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Creation time. Note: CreatedOn is only used as an output parameter and cannot be used as an input parameter in ModifyDnsRecords. If this parameter is passed, it will be ignored.",
						},
						"modified_on": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Modification time. Note: ModifiedOn is only used as an output parameter and cannot be used as an input parameter in ModifyDnsRecords. If this parameter is passed, it will be ignored.",
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
