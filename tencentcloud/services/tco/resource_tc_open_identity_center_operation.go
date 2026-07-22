package tco

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	organization "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/organization/v20210331"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudOpenIdentityCenterOperation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudInviteOpenIdentityCenterOperationCreate,
		Read:   resourceTencentCloudInviteOpenIdentityCenterOperationRead,
		Delete: resourceTencentCloudInviteOpenIdentityCenterOperationDelete,
		Schema: map[string]*schema.Schema{
			"zone_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Space 名称，其中 必须 是 globally 唯一 和 contain 2-64 字符 包括 lowercase letters，digits，和 hyphens (-). It 可以 neither start 或 end 使用 hyphen (-) nor contain two consecutive hyphens (-)。",
			},

			"zone_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Space ID. z-Prefix starts 使用 12 random numbers/lowercase letters followed 通过。",
			},
		},
	}
}

func resourceTencentCloudInviteOpenIdentityCenterOperationCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_open_identity_center_operation.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request  = organization.NewOpenIdentityCenterRequest()
		response = organization.NewOpenIdentityCenterResponse()
	)

	if v, ok := d.GetOk("zone_name"); ok {
		request.ZoneName = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseOrganizationClient().OpenIdentityCenter(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s open identity center operation failed, reason:%+v", logId, err)
		return err
	}

	if response.Response != nil && response.Response.ZoneId != nil {
		d.SetId(*response.Response.ZoneId)
		_ = d.Set("zone_id", *response.Response.ZoneId)
	}

	_ = response

	return resourceTencentCloudInviteOpenIdentityCenterOperationRead(d, meta)
}

func resourceTencentCloudInviteOpenIdentityCenterOperationRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_open_identity_center_operation.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudInviteOpenIdentityCenterOperationDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_open_identity_center_operation.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
