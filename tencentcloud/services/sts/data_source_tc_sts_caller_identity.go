package sts

import (
	"context"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sts "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sts/v20180813"
)

func DataSourceTencentCloudStsCallerIdentity() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudStsCallerIdentityRead,
		Schema: map[string]*schema.Schema{
			"arn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Current caller ARN。",
			},

			"account_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The primary 账号 Uin to which the current caller belongs。",
			},

			"user_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Identity:- When the caller is a cloud 账号，the current 账号 `Uin` is returned.- When the caller is a 角色，it 返回`roleId:roleSessionName`- When the caller is a federated identity，it 返回`uin:federatedUserName`。",
			},

			"principal_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "账号 Uin to which the 键 belongs:- The caller is a cloud 账号，and the returned current 账号 Uin- The caller is a 角色，and the returned 账号 Uin that applies for the 角色 键",
			},

			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Identity 类型",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "用于保存结果。",
			},
		},
	}
}

func dataSourceTencentCloudStsCallerIdentityRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_sts_caller_identity.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	stsService := StsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var callerIdentity *sts.GetCallerIdentityResponseParams
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		results, e := stsService.DescribeStsCallerIdentityByFilter(ctx)
		if e != nil {
			return tccommon.RetryError(e)
		}
		callerIdentity = results
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read Sts instances failed, reason:%+v", logId, err)
		return err
	}

	if callerIdentity.Arn != nil {
		_ = d.Set("arn", callerIdentity.Arn)
	}

	if callerIdentity.AccountId != nil {
		_ = d.Set("account_id", callerIdentity.AccountId)
	}

	if callerIdentity.UserId != nil {
		_ = d.Set("user_id", callerIdentity.UserId)
	}

	if callerIdentity.PrincipalId != nil {
		_ = d.Set("principal_id", callerIdentity.PrincipalId)
	}

	if callerIdentity.Type != nil {
		_ = d.Set("type", callerIdentity.Type)
	}

	d.SetId(*callerIdentity.Arn)

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), map[string]interface{}{
			"arn":          callerIdentity.Arn,
			"account_id":   callerIdentity.AccountId,
			"user_id":      callerIdentity.UserId,
			"principal_id": callerIdentity.PrincipalId,
			"type":         callerIdentity.Type,
		}); e != nil {
			return e
		}
	}

	return nil
}
