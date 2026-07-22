package cos

import (
	"context"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cos "github.com/tencentyun/cos-go-sdk-v5"
)

func ResourceTencentCloudCosBucketReferer() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCosBucketRefererCreate,
		Read:   resourceTencentCloudCosBucketRefererRead,
		Update: resourceTencentCloudCosBucketRefererUpdate,
		Delete: resourceTencentCloudCosBucketRefererDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"bucket": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "存储桶 格式 should 是 [自定义 名称]-[appid]，对于 示例 `mycos-1258798060`。",
			},
			"referer_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Hotlink protection 类型 Enumerated 值: `Black-List`，`White-List`。",
			},
			"status": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "是否enable hotlink protection. Enumerated 值: `已启用`，`已禁用`。",
			},
			"domain_list": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Required:    true,
				Description: "A 列表 域名 names 在 blocklist/allowlist。",
			},
			"empty_refer_configuration": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "是否allow 访问 使用 空 referer. Enumerated 值: `Allow`，`Deny` (默认值)。",
			},
		},
	}
}

func resourceTencentCloudCosBucketRefererCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cos_bucket_referer.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var bucket string
	if v, ok := d.GetOk("bucket"); ok {
		bucket = v.(string)
	}

	d.SetId(bucket)

	return resourceTencentCloudCosBucketRefererUpdate(d, meta)
}

func resourceTencentCloudCosBucketRefererRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cos_bucket_referer.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := CosService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	bucket := d.Id()

	bucketReferer, err := service.DescribeCosBucketRefererById(ctx, bucket)
	if err != nil {
		return err
	}

	if bucketReferer == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `CosBucketReferer` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("bucket", bucket)

	if bucketReferer.Status != "" {
		_ = d.Set("status", bucketReferer.Status)
	}

	if bucketReferer.RefererType != "" {
		_ = d.Set("referer_type", bucketReferer.RefererType)
	}

	if len(bucketReferer.DomainList) > 0 {
		_ = d.Set("domain_list", bucketReferer.DomainList)
	}

	if bucketReferer.EmptyReferConfiguration != "" {
		_ = d.Set("empty_refer_configuration", bucketReferer.EmptyReferConfiguration)
	}

	return nil
}

func resourceTencentCloudCosBucketRefererUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cos_bucket_referer.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	bucket := d.Id()

	request := cos.BucketPutRefererOptions{}
	if v, ok := d.GetOk("status"); ok {
		request.Status = v.(string)
	}
	if v, ok := d.GetOk("referer_type"); ok {
		request.RefererType = v.(string)
	}
	if v, ok := d.GetOk("domain_list"); ok {
		domainListSet := v.(*schema.Set).List()
		for i := range domainListSet {
			domainList := domainListSet[i].(string)
			request.DomainList = append(request.DomainList, domainList)
		}
	}
	if v, ok := d.GetOk("empty_refer_configuration"); ok {
		request.EmptyReferConfiguration = v.(string)
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTencentCosClient(bucket).Bucket.PutReferer(ctx, &request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%+v], response status [%s]\n", logId, "PutReferer", request, result.Status)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s cos bucketReferer failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudCosBucketRefererRead(d, meta)
}

func resourceTencentCloudCosBucketRefererDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cos_bucket_referer.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
