package dnspod

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dnspod "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dnspod/v20210323"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudSubdomainValidateTxtValueOperation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudSubdomainValidateTxtValueOperationCreate,
		Read:   resourceTencentCloudSubdomainValidateTxtValueOperationRead,
		Delete: resourceTencentCloudSubdomainValidateTxtValueOperationDelete,
		Schema: map[string]*schema.Schema{
			"domain_zone": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "subdomain 到 add 可用区 域名",
			},
			"domain": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "域名 名称 对于 其中 TXT records need 到 是 added。",
			},
			"subdomain": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "主机 records 该 need 到 是 added 到 TXT records。",
			},
			"record_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Record types need 到 是 added。",
			},
			"value": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "记录 值 的 TXT 记录 needs 到 是 added。",
			},
		},
	}
}

func resourceTencentCloudSubdomainValidateTxtValueOperationCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_subdomain_validate_txt_value_operation.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	var (
		domainZone string
	)
	var (
		request  = dnspod.NewCreateSubdomainValidateTXTValueRequest()
		response = dnspod.NewCreateSubdomainValidateTXTValueResponse()
	)

	if v, ok := d.GetOk("domain_zone"); ok {
		domainZone = v.(string)
	}

	request.DomainZone = helper.String(domainZone)

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDnsPodClient().CreateSubdomainValidateTXTValueWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create subdomain validate txt value operation failed, reason:%+v", logId, err)
		return err
	}

	if response.Response != nil {
		if response.Response.Domain != nil {
			_ = d.Set("domain", response.Response.Domain)
		}
		if response.Response.Domain != nil {
			_ = d.Set("subdomain", response.Response.Subdomain)
		}
		if response.Response.Domain != nil {
			_ = d.Set("record_type", response.Response.RecordType)
		}
		if response.Response.Domain != nil {
			_ = d.Set("value", response.Response.Value)
		}
	}

	_ = response

	d.SetId(domainZone)

	return resourceTencentCloudSubdomainValidateTxtValueOperationRead(d, meta)
}

func resourceTencentCloudSubdomainValidateTxtValueOperationRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_subdomain_validate_txt_value_operation.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudSubdomainValidateTxtValueOperationDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_subdomain_validate_txt_value_operation.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
