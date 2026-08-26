package vpc

import (
	"context"
	"fmt"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

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
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "Tags of the protocol template.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Tag key.",
						},
						"value": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Tag value.",
						},
					},
				},
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

	tags := make(map[string]string)
	if v, ok := d.GetOk("tags"); ok {
		for _, item := range v.([]interface{}) {
			itemMap := item.(map[string]interface{})
			tags[itemMap["key"].(string)] = itemMap["value"].(string)
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

	if template.TagSet != nil {
		tagList := make([]interface{}, 0, len(template.TagSet))
		for _, tag := range template.TagSet {
			tagMap := make(map[string]interface{})
			if tag.Key != nil {
				tagMap["key"] = *tag.Key
			}
			if tag.Value != nil {
				tagMap["value"] = *tag.Value
			}
			tagList = append(tagList, tagMap)
		}
		_ = d.Set("tags", tagList)
	} else {
		_ = d.Set("tags", []interface{}{})
	}

	return nil
}

func resourceTencentCloudProtocolTemplateUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_protocol_template.update")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	templateId := d.Id()

	immutableArgs := []string{"tags"}
	for _, arg := range immutableArgs {
		if d.HasChange(arg) {
			return fmt.Errorf("argument `%s` cannot be changed for protocol_template %s", arg, templateId)
		}
	}

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
