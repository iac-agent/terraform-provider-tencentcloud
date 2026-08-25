package crs

import (
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	redis "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/redis/v20180412"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudRedisInstancePasswordPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudRedisInstancePasswordPolicyCreate,
		Read:   resourceTencentCloudRedisInstancePasswordPolicyRead,
		Update: resourceTencentCloudRedisInstancePasswordPolicyUpdate,
		Delete: resourceTencentCloudRedisInstancePasswordPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "The ID of instance.",
			},

			"enabled": {
				Required:    true,
				Type:        schema.TypeBool,
				Description: "Whether the instance-level password complexity policy is enabled.",
			},

			"min_letter_count": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Minimum number of letter (upper/lower case) characters. Range [1, 16], default 1.",
			},

			"min_digit_count": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Minimum number of digit characters. Range [1, 16], default 1.",
			},

			"min_special_count": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Minimum number of special characters. Range [1, 16], default 1.",
			},

			"min_length": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Minimum total password length. Range [8, 64], default 8.",
			},
		},
	}
}

func resourceTencentCloudRedisInstancePasswordPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_redis_instance_password_policy.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		instanceId string
	)
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}

	d.SetId(instanceId)

	return resourceTencentCloudRedisInstancePasswordPolicyUpdate(d, meta)
}

func resourceTencentCloudRedisInstancePasswordPolicyRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_redis_instance_password_policy.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	instanceId := d.Id()

	request := redis.NewDescribeInstancePasswordPolicyRequest()
	request.InstanceId = &instanceId

	var passwordPolicy *redis.PasswordPolicy
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseRedisClient().DescribeInstancePasswordPolicy(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		if result == nil || result.Response == nil || result.Response.PasswordPolicy == nil {
			log.Printf("[CRUD] redis_instance_password_policy id=%s", d.Id())
			d.SetId("")
			return nil
		}
		passwordPolicy = result.Response.PasswordPolicy
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read redis_instance_password_policy failed, reason:%+v", logId, err)
		return err
	}

	if passwordPolicy == nil {
		return nil
	}

	_ = d.Set("instance_id", instanceId)

	if passwordPolicy.Enabled != nil {
		_ = d.Set("enabled", passwordPolicy.Enabled)
	}

	if passwordPolicy.MinLetterCount != nil {
		_ = d.Set("min_letter_count", passwordPolicy.MinLetterCount)
	}

	if passwordPolicy.MinDigitCount != nil {
		_ = d.Set("min_digit_count", passwordPolicy.MinDigitCount)
	}

	if passwordPolicy.MinSpecialCount != nil {
		_ = d.Set("min_special_count", passwordPolicy.MinSpecialCount)
	}

	if passwordPolicy.MinLength != nil {
		_ = d.Set("min_length", passwordPolicy.MinLength)
	}

	return nil
}

func resourceTencentCloudRedisInstancePasswordPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_redis_instance_password_policy.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	request := redis.NewModifyInstancePasswordPolicyRequest()

	instanceId := d.Id()
	request.InstanceId = &instanceId

	passwordPolicy := redis.PasswordPolicy{}

	enabled := d.Get("enabled").(bool)
	passwordPolicy.Enabled = &enabled

	if v, ok := d.GetOk("min_letter_count"); ok {
		passwordPolicy.MinLetterCount = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("min_digit_count"); ok {
		passwordPolicy.MinDigitCount = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("min_special_count"); ok {
		passwordPolicy.MinSpecialCount = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("min_length"); ok {
		passwordPolicy.MinLength = helper.IntInt64(v.(int))
	}

	request.PasswordPolicy = &passwordPolicy

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseRedisClient().ModifyInstancePasswordPolicy(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update redis_instance_password_policy failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudRedisInstancePasswordPolicyRead(d, meta)
}

func resourceTencentCloudRedisInstancePasswordPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_redis_instance_password_policy.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
