package cam

import (
	"context"
	"fmt"
	"log"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cam "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cam/v20190116"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudCamPolicyVersion() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCamPolicyVersionCreate,
		Read:   resourceTencentCloudCamPolicyVersionRead,
		Update: resourceTencentCloudCamPolicyVersionUpdate,
		Delete: resourceTencentCloudCamPolicyVersionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"policy_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "Strategy ID。",
			},

			"policy_document": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Strategic text 信息。",
			},

			"set_as_default": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeBool,
				Description: "是否set 作为 版本 的 当前 strategy。",
			},

			"policy_version": {
				Computed:    true,
				Optional:    true,
				Type:        schema.TypeList,
				Description: "Strategic 版本 detailsNote: 此 字段 可能 返回 NULL，indicating 该 有效 值 不能 是 获取。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"version_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Strategic 版本 numberNote: 此 字段 可能 返回 NULL，indicating 该 有效 值 不能 是 获取。",
						},
						"create_date": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Strategic 版本 creation timeNote: 此 字段 可能 返回 NULL，indicating 该 有效 值 不能 是 获取。",
						},
						"is_default_version": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "是否为an effective 版本0 表示 不，1 表示 yesNote: 此 字段 可能 返回 NULL，indicating 该 有效 值 不能 是 获取。",
						},
						"document": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Strategic grammar textNote: 此 字段 可能 返回 NULL，indicating 该 有效 值 不能 是 获取。",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudCamPolicyVersionCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cam_policy_version.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request   = cam.NewCreatePolicyVersionRequest()
		response  = cam.NewCreatePolicyVersionResponse()
		policyId  string
		versionId string
	)
	if v, ok := d.GetOkExists("policy_id"); ok {
		policyId = helper.IntToStr(v.(int))
		request.PolicyId = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("policy_document"); ok {
		request.PolicyDocument = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("set_as_default"); ok {
		request.SetAsDefault = helper.Bool(v.(bool))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCamClient().CreatePolicyVersion(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create cam policyVersion failed, reason:%+v", logId, err)
		return err
	}
	if response == nil || response.Response == nil || response.Response.VersionId == nil {
		return fmt.Errorf("CAM policy version is null")
	}
	versionId = helper.UInt64ToStr(*response.Response.VersionId)
	d.SetId(policyId + tccommon.FILED_SP + versionId)

	return resourceTencentCloudCamPolicyVersionRead(d, meta)
}

func resourceTencentCloudCamPolicyVersionRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cam_policy_version.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := CamService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	policyId := idSplit[0]
	versionId := idSplit[1]

	policyVersion, err := service.DescribeCamPolicyVersionById(ctx, helper.StrToUInt64(policyId), helper.StrToUInt64(versionId))
	if err != nil {
		return err
	}

	if policyVersion == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `CamPolicyVersion` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("policy_id", helper.StrToInt64(policyId))

	if policyVersion != nil {
		policyVersionMap := map[string]interface{}{}

		if policyVersion.VersionId != nil {
			policyVersionMap["version_id"] = policyVersion.VersionId
		}

		if policyVersion.CreateDate != nil {
			policyVersionMap["create_date"] = policyVersion.CreateDate
		}

		if policyVersion.IsDefaultVersion != nil {
			policyVersionMap["is_default_version"] = policyVersion.IsDefaultVersion
		}

		if policyVersion.Document != nil {
			policyVersionMap["document"] = policyVersion.Document
			_ = d.Set("policy_document", policyVersion.Document)
		}

		_ = d.Set("policy_version", []interface{}{policyVersionMap})

		if policyVersion.IsDefaultVersion != nil {
			if *policyVersion.IsDefaultVersion == 0 {
				_ = d.Set("set_as_default", false)
			} else {
				_ = d.Set("set_as_default", true)
			}
		}
	}

	return nil
}

func resourceTencentCloudCamPolicyVersionDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cam_policy_version.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := CamService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	policyId := idSplit[0]
	versionId := idSplit[1]

	if err := service.DeleteCamPolicyVersionById(ctx, helper.StrToUInt64(policyId), helper.StrToUInt64(versionId)); err != nil {
		return err
	}

	return nil
}
func resourceTencentCloudCamPolicyVersionUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cam_policy_version.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	return resourceTencentCloudCamPolicyVersionRead(d, meta)
}
