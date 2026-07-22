package cos

import (
	"context"
	"encoding/json"
	"log"
	"reflect"
	"time"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	svccam "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/cam"
)

func ResourceTencentCloudCosBucketPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCosBucketPolicyCreate,
		Read:   resourceTencentCloudCosBucketPolicyRead,
		Update: resourceTencentCloudCosBucketPolicyUpdate,
		Delete: resourceTencentCloudCosBucketPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"bucket": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: tccommon.ValidateCosBucketName,
				Description:  "名称 存储桶 到 是 创建. 存储桶 格式 should 是 [自定义 名称]-[appid]，对于 示例 `mycos-1258798060`。",
			},
			"policy": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsJSON,
				DiffSuppressFunc: func(k, olds, news string, d *schema.ResourceData) bool {
					var oldJson interface{}
					err := json.Unmarshal([]byte(olds), &oldJson)
					if err != nil {
						return olds == news
					}
					var newJson interface{}
					err = json.Unmarshal([]byte(news), &newJson)
					if err != nil {
						return olds == news
					}
					flag := reflect.DeepEqual(oldJson, newJson)
					return flag
				},
				Description: "text 的 策略. For more info please refer 到 [Tencent official doc](https://intl.云.tencent.com/document/product/436/18023)， six-segment 资源 scenario 示例 在 document 是 作为 follows: Example 的 specifying 存储桶: `qcs::cos:ap-guangzhou:uid/1250000000:examplebucket-1250000000/*`，In cdc scenario: `qcs::cos:ap-guangzhou:uid/1250000000:cdc_cluster-123456_examplebucket-1250000000/*`; Example 的 specifying 文件夹: `qcs::cos:ap-guangzhou:uid/1250000000:examplebucket-1250000000/文件夹/*`，In cdc scenario: `qcs::cos:ap-guangzhou:uid/1250000000:cdc_cluster-123456_examplebucket-1250000000/文件夹/*`; Specified 对象 示例: `qcs::cos:ap-guangzhou:uid/1250000000:examplebucket-1250000000/文件夹/exampleobject`，In cdc scenario: `qcs::cos:ap-guangzhou:uid/1250000000:cdc_cluster-123456_examplebucket-1250000000/文件夹/exampleobject`。",
			},
		},
	}
}

func resourceTencentCloudCosBucketPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cos_bucket_policy.create")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	bucket := d.Get("bucket").(string)
	policy := d.Get("policy").(string)

	cosService := CosService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}
	camService := svccam.NewCamService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
	policyErr := camService.PolicyDocumentForceCheck(policy)
	if policyErr != nil {
		return policyErr
	}

	err := cosService.PutBucketPolicy(ctx, bucket, policy)
	if err != nil {
		return err
	}
	d.SetId(bucket)

	return resourceTencentCloudCosBucketPolicyRead(d, meta)
}

func resourceTencentCloudCosBucketPolicyRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cos_bucket_policy.read")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	bucket := d.Id()
	cosService := CosService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var result string
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		policy, e := cosService.DescribePolicyByBucket(ctx, bucket)
		if e != nil {
			return tccommon.RetryError(e)
		}
		result = policy
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read cos bucket policy failed, reason:%s\n", logId, err.Error())
		return err
	}
	result, err = removeSid(result)
	if err != nil {
		log.Printf("[CRITAL]%s read cos bucket policy failed, reason:%s\n", logId, err.Error())
		return err
	}
	if result == "" {
		d.SetId("")
		return nil
	}
	_ = d.Set("policy", result)
	_ = d.Set("bucket", bucket)

	return nil
}

func resourceTencentCloudCosBucketPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cos_bucket_policy.update")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	cosService := CosService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	bucket := d.Id()

	if d.HasChange("policy") {
		policy := d.Get("policy").(string)
		camService := svccam.NewCamService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
		policyErr := camService.PolicyDocumentForceCheck(policy)
		if policyErr != nil {
			return policyErr
		}
		err := cosService.PutBucketPolicy(ctx, bucket, policy)
		if err != nil {
			return err
		}
	}

	// wait for update cache
	// if not, the data may be outdated.
	time.Sleep(3 * time.Second)

	return resourceTencentCloudCosBucketPolicyRead(d, meta)
}

func resourceTencentCloudCosBucketPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cos_bucket_policy.delete")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	bucket := d.Id()
	cosService := CosService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}
	err := cosService.DeleteBucketPolicy(ctx, bucket)
	if err != nil {
		return err
	}

	// wait for update cache
	// if not, head bucket may be successful
	time.Sleep(3 * time.Second)

	return nil
}

// In the returned JSON, the SDK automatically adds the Sid, which needs to be removed
func removeSid(v string) (result string, err error) {
	m := make(map[string]interface{})
	err = json.Unmarshal([]byte(v), &m)
	if err != nil {
		return
	}
	var stateMend []interface{}
	if v, ok := m["Statement"]; ok {
		stateMend = v.([]interface{})
	}
	for index, v := range stateMend {
		mp := v.(map[string]interface{})
		delete(mp, "Sid")
		stateMend[index] = mp
	}
	if _, ok := m["Statement"]; ok {
		m["Statement"] = stateMend
	}
	s, err := json.Marshal(m)
	return string(s), err
}
