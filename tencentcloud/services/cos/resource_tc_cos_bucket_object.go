package cos

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mitchellh/go-homedir"
)

func ResourceTencentCloudCosBucketObject() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCosBucketObjectCreate,
		Read:   resourceTencentCloudCosBucketObjectRead,
		Update: resourceTencentCloudCosBucketObjectUpdate,
		Delete: resourceTencentCloudCosBucketObjectDelete,

		Schema: map[string]*schema.Schema{
			"bucket": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "名称 存储桶 到 使用. 存储桶 格式 should 是 [自定义 名称]-[appid]，对于 示例 `mycos-1258798060`。",
			},
			"key": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "名称 对象 once 它 是 在 存储桶",
			},
			"source": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"content"},
				Description:   "路径 到 来源 文件 being uploaded 到 存储桶",
			},
			"content": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"source"},
				Description:   "Literal 字符串 值 到 使用 作为 对象 内容，其中 将 是 uploaded 作为 UTF-8-encoded text。",
			},
			"acl": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  s3.ObjectCannedACLPrivate,
				ValidateFunc: tccommon.ValidateAllowedStringValue([]string{
					s3.ObjectCannedACLPrivate,
					s3.ObjectCannedACLPublicRead,
					s3.ObjectCannedACLPublicReadWrite,
				}),
				Description: "canned ACL 到 apply. Available 值 include `私有`，`公有-read`，和 `公有-read-write`. 默认为 `私有`。",
			},
			"cache_control": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "指定caching behavior along 请求/reply chain. For further details，RFC2616 可以 是 referred。",
			},
			"content_disposition": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "指定presentational 信息 对于 对象。",
			},
			"content_encoding": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "指定what 内容 encodings have been applied 到 对象 和 thus what decoding mechanisms 必须 是 applied 到 obtain media-类型 referenced 通过 内容-类型 头部 字段。",
			},
			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "标签 的 对象。",
			},
			"content_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "A standard MIME 类型 describing 格式 的 对象 数据。",
			},
			"storage_class": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Object 存储 类型，Available 值 include `STANDARD_IA`，`MAZ_STANDARD_IA`，`INTELLIGENT_TIERING`，`MAZ_INTELLIGENT_TIERING`，`ARCHIVE`，`DEEP_ARCHIVE`. For more 信息，please refer 到: https://云.tencent.com/document/product/436/33417。",
			},
			"etag": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "ETag generated 对于 对象 ( MD5 sum 的 对象 内容)。",
			},
		},
	}
}

func resourceTencentCloudCosBucketObjectCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cos_bucket_object.create")()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	bucket := d.Get("bucket").(string)
	key := d.Get("key").(string)
	var body io.ReadSeeker
	if v, ok := d.GetOk("source"); ok {
		source := v.(string)
		path, err := homedir.Expand(source)
		if err != nil {
			return fmt.Errorf("cos object source (%s) homedir expand error: %s", source, err.Error())
		}
		file, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("cos object source (%s) open error: %s", source, err.Error())
		}
		body = file
		defer func() {
			err := file.Close()
			if err != nil {
				log.Printf("closing cos object source (%s) error: %s", path, err.Error())
			}
		}()
	} else if v, ok := d.GetOk("content"); ok {
		content := v.(string)
		body = bytes.NewReader([]byte(content))
	} else {
		return fmt.Errorf("must specify \"source\" or \"content\" field")
	}

	request := &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   body,
	}

	if v, ok := d.GetOk("acl"); ok {
		request.ACL = aws.String(v.(string))
	}
	if v, ok := d.GetOk("cache_control"); ok {
		request.CacheControl = aws.String(v.(string))
	}
	if v, ok := d.GetOk("content_disposition"); ok {
		request.ContentDisposition = aws.String(v.(string))
	}
	if v, ok := d.GetOk("content_encoding"); ok {
		request.ContentEncoding = aws.String(v.(string))
	}
	if v, ok := d.GetOk("content_type"); ok {
		request.ContentType = aws.String(v.(string))
	}
	if v, ok := d.GetOk("storage_class"); ok {
		request.StorageClass = aws.String(v.(string))
	}

	response, err := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCosClient().PutObject(request)
	if err != nil {
		return fmt.Errorf("putting object (%s) in cos bucket (%s) error: %s", key, bucket, err.Error())
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, "put object", request.String(), response.String())

	if v, ok := d.GetOk("tags"); ok {
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service := CosService{
			client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
		}
		tags := make(map[string]string)

		for key, val := range v.(map[string]interface{}) {
			tags[key] = val.(string)
		}

		if err := service.SetObjectTags(ctx, bucket, key, tags); err != nil {
			log.Printf("[WARN] set object tags error, skip processing")
		}
	}

	d.SetId(bucket + key)
	return resourceTencentCloudCosBucketObjectRead(d, meta)
}

func resourceTencentCloudCosBucketObjectRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cos_bucket_object.read")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	bucket := d.Get("bucket").(string)
	key := d.Get("key").(string)

	cosService := CosService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}
	response, err := cosService.HeadObject(ctx, bucket, key)
	if err != nil {
		if awsError, ok := err.(awserr.RequestFailure); ok && awsError.StatusCode() == 404 {
			log.Printf("[WARN]%s object (%s) in bucket (%s) not found, error code (404)", logId, key, bucket)
			d.SetId("")
			return nil
		}
		return err
	}

	_ = d.Set("cache_control", response.CacheControl)
	_ = d.Set("content_disposition", response.ContentDisposition)
	_ = d.Set("content_encoding", response.ContentEncoding)
	_ = d.Set("content_type", response.ContentType)
	_ = d.Set("etag", strings.Trim(*response.ETag, `"`))
	_ = d.Set("storage_class", s3.StorageClassStandard)
	if response.StorageClass != nil {
		_ = d.Set("storage_class", response.StorageClass)
	}

	_, aclResponse, aclErr := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTencentCosClient(bucket).Object.GetACL(ctx, key)
	if aclErr != nil {
		return fmt.Errorf("cos [GetACL] error: %s, bucket: %s, object: %s", aclErr.Error(), bucket, key)
	}
	if aclResponse.StatusCode == 404 {
		log.Printf("[WARN] [GetACL] returns %d, %s", 404, err)
		return nil
	}

	_ = d.Set("acl", aclResponse.Header.Get("x-cos-acl"))

	var tags map[string]string
	tags, err = cosService.GetObjectTags(ctx, bucket, key)
	if err != nil {
		if awsError, ok := err.(awserr.RequestFailure); ok && awsError.StatusCode() == 404 {
			log.Printf("[WARN]%s tags in object (%s) of bucket (%s) not found, error code (404)", logId, key, bucket)
			d.SetId("")
			return nil
		}
		return err
	}
	_ = d.Set("tags", tags)

	return nil
}

func resourceTencentCloudCosBucketObjectUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cos_bucket_object.update")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	fields := []string{
		"cache_control",
		"content_disposition",
		"content_encoding",
		"content_type",
		"source",
		"content",
		"storage_class",
		"etag",
	}
	for _, key := range fields {
		if d.HasChange(key) {
			return resourceTencentCloudCosBucketObjectCreate(d, meta)
		}
	}

	bucket := d.Get("bucket").(string)
	key := d.Get("key").(string)
	cosService := CosService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}

	if d.HasChange("acl") {
		acl := d.Get("acl").(string)
		err := cosService.PutObjectAcl(ctx, bucket, key, acl)
		if err != nil {
			return err
		}
	}

	if d.HasChange("tags") {
		v := d.Get("tags").(map[string]interface{})
		tags := make(map[string]string)
		for key, val := range v {
			tags[key] = val.(string)
		}
		if err := cosService.SetObjectTags(ctx, bucket, key, tags); err != nil {
			return err
		}
	}

	return nil
}

func resourceTencentCloudCosBucketObjectDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cos_bucket_object.delete")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	bucket := d.Get("bucket").(string)
	key := d.Get("key").(string)

	cosService := CosService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}
	err := cosService.DeleteObject(ctx, bucket, key)
	if err != nil {
		return err
	}

	return nil
}
