package ssl

import (
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ssl "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ssl/v20191205"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudSslDeployCertificateInstanceOperation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudSslDeployCertificateInstanceCreate,
		Read:   resourceTencentCloudSslDeployCertificateInstanceRead,
		Delete: resourceTencentCloudSslDeployCertificateInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"certificate_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "ID 证书 到 是 deployed。",
			},

			"instance_id_list": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Need 到 deploy 实例 列表。",
			},

			"resource_type": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Deployed 云 资源类型",
			},

			"status": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "Deployment 云 资源 状态: Live: -1: 域名 名称 是 不 associated 使用 证书.1: 域名 名称 https 是 已启用0: 域名 名称 https 是 closed。",
			},
		},
	}
}

func resourceTencentCloudSslDeployCertificateInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ssl_deploy_certificate_instance_operation.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request        = ssl.NewDeployCertificateInstanceRequest()
		response       = ssl.NewDeployCertificateInstanceResponse()
		deployRecordId uint64
	)
	if v, ok := d.GetOk("certificate_id"); ok {
		request.CertificateId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("instance_id_list"); ok {
		instanceIdListSet := v.(*schema.Set).List()
		for i := range instanceIdListSet {
			instanceIdList := instanceIdListSet[i].(string)
			request.InstanceIdList = append(request.InstanceIdList, &instanceIdList)
		}
	}

	if v, ok := d.GetOk("resource_type"); ok {
		request.ResourceType = helper.String(v.(string))
	}

	if v, _ := d.GetOk("status"); v != nil {
		request.Status = helper.IntInt64(v.(int))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseSSLCertificateClient().DeployCertificateInstance(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate ssl deployCertificateInstance failed, reason:%+v", logId, err)
		return err
	}

	deployRecordId = *response.Response.DeployRecordId
	d.SetId(helper.UInt64ToStr(deployRecordId))

	return resourceTencentCloudSslDeployCertificateInstanceRead(d, meta)
}

func resourceTencentCloudSslDeployCertificateInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ssl_deploy_certificate_instance_operation.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudSslDeployCertificateInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ssl_deploy_certificate_instance_operation.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
