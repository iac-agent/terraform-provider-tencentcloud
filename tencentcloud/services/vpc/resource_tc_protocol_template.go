package vpc

import (
	"context"
	"fmt"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
	svctag "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceTencentCloudProtocolTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudProtocolTemplateCreate,
		Read:   resourceTencentCloudProtocolTemplateRead,
		Update: resourceTencentCloudProtocolTemplateUpdate,
		Delete: resourceTencentCloudProtocolTemplateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the protocol template.",
			},
			"protocols": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: tccommon.ValidateLowCase,
				},
				Required:    true,
				Description: "Protocol list. Valid protocols are  `tcp`, `udp`, `icmp`, `gre`. Single port(tcp:80), multi-port(tcp:80,443), port range(tcp:3306-20000), all(tcp:all) format are support. Protocol `icmp` and `gre` cannot specify port.",
			},
			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Tags of the protocol template.",
			},
		},
	}
}

func resourceTencentCloudProtocolTemplateCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_protocol_template.create")()
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	name := d.Get("name").(string)
	protocols := d.Get("protocols").(*schema.Set).List()

	var tags []*vpc.Tag
	if v, ok := d.GetOk("tags"); ok {
		tagMap := v.(map[string]interface{})
		tags = make([]*vpc.Tag, 0, len(tagMap))
		for k, t := range tagMap {
			key := k
			value := t.(string)
			tag := vpc.Tag{
				Key:   &key,
				Value: &value,
			}
			tags = append(tags, &tag)
		}
	}

	vpcProtocol := VpcService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}
	var outErr, inErr error
	var templateId string

	outErr = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		templateId, inErr = vpcProtocol.CreateServiceTemplate(ctx, name, protocols, tags)
		if inErr != nil {
			return tccommon.RetryError(inErr)
		}
		return nil
	})
	if outErr != nil {
		return outErr
	}

	d.SetId(templateId)

	return resourceTencentCloudProtocolTemplateRead(d, meta)
}

func resourceTencentCloudProtocolTemplateRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_protocol_template.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	templateId := d.Id()
	var outErr, inErr error
	vpcProtocol := VpcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	template, has, outErr := vpcProtocol.DescribeServiceTemplateById(ctx, templateId)
	if outErr != nil {
		outErr = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			template, has, inErr = vpcProtocol.DescribeServiceTemplateById(ctx, templateId)
			if inErr != nil {
				return tccommon.RetryError(inErr)
			}
			return nil
		})
	}
	if outErr != nil {
		return outErr
	}
	if !has {
		d.SetId("")
		return nil
	}

	_ = d.Set("name", template.ServiceTemplateName)
	_ = d.Set("protocols", template.ServiceSet)

	if len(template.TagSet) > 0 {
		tags := make(map[string]string)
		for _, tag := range template.TagSet {
			if tag.Key != nil && tag.Value != nil {
				tags[*tag.Key] = *tag.Value
			}
		}
		_ = d.Set("tags", tags)
	}

	return nil
}

func resourceTencentCloudProtocolTemplateUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_protocol_template.update")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	templateId := d.Id()

	if d.HasChange("name") || d.HasChange("protocols") {
		var outErr, inErr error
		name := d.Get("name").(string)
		protocols := d.Get("protocols").(*schema.Set).List()
		vpcProtocol := VpcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		outErr = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			inErr = vpcProtocol.ModifyServiceTemplate(ctx, templateId, name, protocols)
			if inErr != nil {
				return tccommon.RetryError(inErr, "UnsupportedOperation.MutexOperationTaskRunning")
			}
			return nil
		})
		if outErr != nil {
			return outErr
		}
	}

	if d.HasChange("tags") {
		client := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
		tagService := svctag.NewTagService(client)
		oldValue, newValue := d.GetChange("tags")
		replaceTags, deleteTags := svctag.DiffTags(oldValue.(map[string]interface{}), newValue.(map[string]interface{}))
		region := client.Region

		resourceName := tccommon.BuildTagResourceName("vpc", "service", region, d.Id())
		err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags)
		if err != nil {
			return err
		}
	}

	return resourceTencentCloudProtocolTemplateRead(d, meta)
}

func resourceTencentCloudProtocolTemplateDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_protocol_template.delete")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	templateId := d.Id()
	vpcProtocol := VpcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	var outErr, inErr error

	outErr = vpcProtocol.DeleteServiceTemplate(ctx, templateId)
	if outErr != nil {
		outErr = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			inErr = vpcProtocol.DeleteServiceTemplate(ctx, templateId)
			if inErr != nil {
				return tccommon.RetryError(inErr, "UnsupportedOperation.MutexOperationTaskRunning")
			}
			return nil
		})
	}

	if outErr != nil {
		return outErr
	}
	//check not exist
	outErr = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		_, has, inErr := vpcProtocol.DescribeServiceTemplateById(ctx, templateId)
		if inErr != nil {
			return tccommon.RetryError(inErr)
		}
		if has {
			return resource.RetryableError(fmt.Errorf("protocol template %s is still exists, retry...", templateId))
		} else {
			return nil
		}
	})

	return outErr
}
