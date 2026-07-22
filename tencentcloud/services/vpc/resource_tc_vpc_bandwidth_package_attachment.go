package vpc

import (
	"context"
	"fmt"
	"log"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudVpcBandwidthPackageAttachment() *schema.Resource {
	return &schema.Resource{
		Read:   resourceTencentCloudVpcBandwidthPackageAttachmentRead,
		Create: resourceTencentCloudVpcBandwidthPackageAttachmentCreate,
		Delete: resourceTencentCloudVpcBandwidthPackageAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"resource_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "唯一 ID 资源，currently 支持 EIP resources 和 LB resources，such 作为 `eip-xxxx`，`lb-xxxx`。",
			},

			"bandwidth_package_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Bandwidth 包 唯一 ID，在 form 的 `bwp-xxxx`。",
			},

			"network_type": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Bandwidth packet 类型，currently 支持 `BGP` 类型，indicating 该 内部 资源 是 BGP IP。",
			},

			"resource_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "Resource types，包括 `地址`，`LoadBalance`。",
			},

			"protocol": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Bandwidth packet 协议 类型 Currently `ipv4` 和 `ipv6` 协议 types 是 支持。",
			},
		},
	}
}

func resourceTencentCloudVpcBandwidthPackageAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vpc_bandwidth_package_attachment.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request            = vpc.NewAddBandwidthPackageResourcesRequest()
		bandwidthPackageId string
		resourceId         string
	)

	if v, ok := d.GetOk("resource_id"); ok {
		resourceId = v.(string)
		request.ResourceIds = append(request.ResourceIds, helper.String(v.(string)))
	}

	if v, ok := d.GetOk("bandwidth_package_id"); ok {
		bandwidthPackageId = v.(string)
		request.BandwidthPackageId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("network_type"); ok {
		request.NetworkType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("resource_type"); ok {
		request.ResourceType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("protocol"); ok {
		request.Protocol = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().AddBandwidthPackageResources(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create vpc bandwidthPackageAttachment failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(bandwidthPackageId + tccommon.FILED_SP + resourceId)
	return resourceTencentCloudVpcBandwidthPackageAttachmentRead(d, meta)
}

func resourceTencentCloudVpcBandwidthPackageAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vpc_bandwidth_package_attachment.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := VpcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	bandwidthPackageId := idSplit[0]
	resourceId := idSplit[1]

	bandwidthPackageAttachment, err := service.DescribeVpcBandwidthPackageAttachment(ctx, bandwidthPackageId, resourceId)

	if err != nil {
		return err
	}

	if bandwidthPackageAttachment == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `tencentcloud_vpc_bandwidth_package_attachment` [%s] not found, please check if it has been deleted.",
			logId, bandwidthPackageId,
		)
		return nil
	}

	_ = d.Set("bandwidth_package_id", bandwidthPackageId)

	if bandwidthPackageAttachment.ResourceId != nil {
		_ = d.Set("resource_id", bandwidthPackageAttachment.ResourceId)
	}

	//if bandwidthPackageAttachment.NetworkType != nil {
	//	_ = d.Set("network_type", bandwidthPackageAttachment.NetworkType)
	//}

	if bandwidthPackageAttachment.ResourceType != nil {
		_ = d.Set("resource_type", bandwidthPackageAttachment.ResourceType)
	}

	//if bandwidthPackageAttachment.Protocol != nil {
	//	_ = d.Set("protocol", bandwidthPackageAttachment.Protocol)
	//}

	return nil
}

func resourceTencentCloudVpcBandwidthPackageAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vpc_bandwidth_package_attachment.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := VpcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	bandwidthPackageId := idSplit[0]
	resourceId := idSplit[1]

	if err := service.DeleteVpcBandwidthPackageAttachmentById(ctx, bandwidthPackageId, resourceId); err != nil {
		return err
	}

	return nil
}
