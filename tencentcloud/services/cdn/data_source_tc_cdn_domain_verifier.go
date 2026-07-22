package cdn

import (
	"context"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/go-multierror"
	sdkError "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cdn "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdn/v20180606"
)

func DataSourceTencentCloudCdnDomainVerifyRecord() *schema.Resource {
	return &schema.Resource{
		Read: resourceTencentCloudCdnDomainVerifyRecordRead,
		Schema: map[string]*schema.Schema{
			"domain": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "指定domain 名称，e.g. `www.examplexxx123.com`。",
			},
			"verify_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "指定verify 类型，值: `dns` (默认值)，`文件`。",
			},
			"auto_verify": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "指定是否keep first create 结果 instead 的 re-create again。",
			},
			"freeze_record": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "指定是否verification 记录 needs 到 是 freeze instead 的 refresh every 8 hours，此 用于domain verification。",
			},
			"verify_result": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Verify 结果",
			},
			"failed_reason": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "表示failed reason 的 verification。",
			},
			"sub_domain": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Sub-域名 resolution。",
			},
			"record": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Resolution 记录 值",
			},
			"record_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "类型 resolution。",
			},
			"file_verify_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "File verify URL guidance。",
			},
			"file_verify_domains": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "列表 文件 verified domains。",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"file_verify_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "名称 文件 verifications。",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "用于save 结果 json。",
			},
		},
	}
}

func resourceTencentCloudCdnDomainVerifyRecordRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cdn_domain_verifier.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := CdnService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	domainName := d.Get("domain").(string)
	freezeRecord := d.Get("freeze_record").(bool)
	autoVerify := d.Get("auto_verify").(bool)

	verified, err := service.VerifyDomainRecord(ctx, domainName)

	if err != nil {
		canContinue, reason := checkCdnDomainVerifyErrReason(err)
		if !canContinue {
			d.SetId("")
			return err
		}
		_ = d.Set("failed_reason", reason)
	}

	_ = d.Set("verify_result", verified)
	d.SetId(domainName)

	if !autoVerify || (verified && freezeRecord) {
		return nil
	}
	response, err := service.CreateVerifyRecord(ctx, domainName)

	if err != nil {
		return err
	}

	var errResults *multierror.Error

	errResults = multierror.Append(errResults, d.Set("record", response.Record))
	errResults = multierror.Append(errResults, d.Set("record_type", response.RecordType))
	errResults = multierror.Append(errResults, d.Set("sub_domain", response.SubDomain))
	if len(response.FileVerifyDomains) > 0 {
		errResults = multierror.Append(errResults, d.Set("file_verify_domains", response.FileVerifyDomains))
	}
	errResults = multierror.Append(errResults, d.Set("file_verify_name", response.FileVerifyName))
	errResults = multierror.Append(errResults, d.Set("file_verify_url", response.FileVerifyUrl))

	if e := errResults.ErrorOrNil(); e != nil {
		return e
	}

	if err != nil {
		return err
	}

	if output, ok := d.GetOk("result_output_file"); ok && output.(string) != "" {
		result := map[string]interface{}{
			"verify_result":       verified,
			"failed_reason":       d.Get("failed_reason"),
			"sub_domain":          response.SubDomain,
			"record":              response.Record,
			"record_type":         response.RecordType,
			"file_verify_url":     response.FileVerifyUrl,
			"file_verify_domains": response.FileVerifyDomains,
			"file_verify_name":    response.FileVerifyName,
		}
		if err := tccommon.WriteToFile(output.(string), result); err != nil {
			log.Printf("[CRITAL]%s output file[%s] fail, reason[%v]",
				logId, output.(string), err)
			return err
		}
	}

	return nil
}

func checkCdnDomainVerifyErrReason(err error) (canContinue bool, code string) {
	sdkErr, ok := err.(*sdkError.TencentCloudSDKError)
	if !ok {
		return
	}
	errCode := sdkErr.Code
	if errCode == cdn.UNAUTHORIZEDOPERATION_CDNDOMAINRECORDNOTVERIFIED ||
		errCode == cdn.UNAUTHORIZEDOPERATION_CDNTXTRECORDVALUENOTMATCH {
		canContinue = true
		code = errCode
	}
	return
}
