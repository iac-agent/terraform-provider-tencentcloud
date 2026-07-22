package vpc

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	svctag "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tag"
)

func ResourceTencentCloudElasticPublicIpv6() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudElasticPublicIpv6Create,
		Read:   resourceTencentCloudElasticPublicIpv6Read,
		Update: resourceTencentCloudElasticPublicIpv6Update,
		Delete: resourceTencentCloudElasticPublicIpv6Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"address_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "EIP 名称，用于customize personalized 名称 EIP 当 applying 对于 EIP. 默认值：unnamed。",
			},

			"address_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Elastic IPv6 类型，可选 值:\n\t- EIPv6: Ordinary IPv6\n\t- HighQualityEIPv6: Premium IPv6\nNote: You need 到 contact product 到 open premium IPv6 white 列表，和 仅 some regions support premium IPv6\n默认值：EIPv6。",
			},

			"internet_charge_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Elastic IPv6 charging 方法，可选 值:\n\t- BANDWIDTH_PACKAGE: Payment 对于 Shared Bandwidth Package\n\t- TRAFFIC_POSTPAID_BY_HOUR: Traffic 是 paid 通过 hour\n默认值：TRAFFIC_POSTPAID_BY_HOUR。",
			},

			"internet_service_provider": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Elastic IPv6 line 类型，默认值：BGP.\nFor users who have activated 静态 单个-line IP whitelist，selectable 值:\n\t- CMCC: China Mobile\n\t- CTCC: China Telecom\n\t- CUCC: China Unicom\nNote: Static 单个-wire IP 是 仅 支持 在 some regions。",
			},

			"internet_max_bandwidth_out": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Elastic IPv6 带宽 限制 在 Mbps.\nThe 范围 的 selectable 值 depends 在 EIP billing 方法:\n\t- BANDWIDTH_PACKAGE: 1 Mbps 到 2000 Mbps\n\t- TRAFFIC_POSTPAID_BY_HOUR: 1 Mbps 到 100 Mbps\n默认值：1 Mbps。",
			},

			"bandwidth_package_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Bandwidth packet 唯一 ID 参数. 如果 此 参数 是 集合 和 InternetChargeType 是 BANDWIDTH_PACKAGE，它 表示 该 EIP 创建 是 added 到 BGP 带宽 packet 和 带宽 packet 是 charged。",
			},

			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "标签",
			},

			"egress": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Elastic IPv6 网络 exit，可选 值:\n\t- CENTER_EGRESS_1: Center Exit 1\n\t- CENTER_EGRESS_2: Center Exit 2\n\t- CENTER_EGRESS_3: Center Exit 3\nNote: Network exports corresponding 到 different operators 或 资源 types need 到 contact product 对于 clarification\n默认值：CENTER_EGRESS_1。",
			},

			"address_ip": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "External 网络 IP 地址",
			},
		},
	}
}

func resourceTencentCloudElasticPublicIpv6Create(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_elastic_public_ipv6.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	var (
		ipId   string
		taskId string
	)
	var (
		request  = vpc.NewAllocateIPv6AddressesRequest()
		response = vpc.NewAllocateIPv6AddressesResponse()
	)

	if v, ok := d.GetOk("address_name"); ok {
		request.AddressName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("address_type"); ok {
		request.AddressType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("internet_charge_type"); ok {
		request.InternetChargeType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("internet_service_provider"); ok {
		request.InternetServiceProvider = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("internet_max_bandwidth_out"); ok {
		request.InternetMaxBandwidthOut = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("bandwidth_package_id"); ok {
		request.BandwidthPackageId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("egress"); ok {
		request.Egress = helper.String(v.(string))
	}

	if v := helper.GetTags(d, "tags"); len(v) > 0 {
		for tagKey, tagValue := range v {
			tag := vpc.Tag{
				Key:   helper.String(tagKey),
				Value: helper.String(tagValue),
			}
			request.Tags = append(request.Tags, &tag)
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().AllocateIPv6AddressesWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create elastic public ipv6 failed, reason:%+v", logId, err)
		return err
	}

	if len(response.Response.AddressSet) < 1 {
		return fmt.Errorf("resource `tencentcloud_elastic_public_ipv6` create failed.")
	}

	ipId = *response.Response.AddressSet[0]
	taskId = *response.Response.TaskId

	if _, err := (&resource.StateChangeConf{
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
		Pending:    []string{},
		Refresh:    resourceElasticPublicIpv6CreateStateRefreshFunc_0_0(ctx, taskId),
		Target:     []string{"SUCCESS"},
		Timeout:    600 * time.Second,
	}).WaitForStateContext(ctx); err != nil {
		return err
	}
	d.SetId(ipId)

	return resourceTencentCloudElasticPublicIpv6Read(d, meta)
}

func resourceTencentCloudElasticPublicIpv6Read(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_elastic_public_ipv6.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	service := VpcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	ipId := d.Id()

	respData, err := service.DescribeElasticPublicIpv6ById(ctx, ipId)
	if err != nil {
		return err
	}

	if respData == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `elastic_public_ipv6` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}
	if len(respData.AddressSet) > 0 {
		address := respData.AddressSet[0]
		_ = d.Set("address_name", address.AddressName)
		_ = d.Set("address_type", address.AddressType)
		_ = d.Set("internet_charge_type", address.InternetChargeType)
		_ = d.Set("internet_service_provider", address.InternetServiceProvider)
		_ = d.Set("internet_max_bandwidth_out", address.Bandwidth)
		_ = d.Set("bandwidth_package_id", address.BandwidthPackageId)
		_ = d.Set("egress", address.Egress)
		_ = d.Set("address_ip", address.AddressIp)
	}
	tcClient := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
	tagService := svctag.NewTagService(tcClient)
	region := tcClient.Region
	tags, err := tagService.DescribeResourceTags(ctx, "vpc", "eipv6", region, ipId)
	if err != nil {
		log.Printf("[CRITAL]%s describe eipv6 tags failed: %+v", logId, err)
		return err
	}
	_ = d.Set("tags", tags)
	return nil
}

func resourceTencentCloudElasticPublicIpv6Update(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_elastic_public_ipv6.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	immutableArgs := []string{"address_type", "internet_charge_type", "internet_service_provider", "bandwidth_package_id", "egress"}
	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}
	ipId := d.Id()

	needChange := false
	mutableArgs := []string{"address_name"}
	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}

	if needChange {
		request := vpc.NewModifyIPv6AddressesAttributesRequest()
		request.IPv6AddressIds = []*string{helper.String(ipId)}

		if v, ok := d.GetOk("address_name"); ok {
			request.IPv6AddressName = helper.String(v.(string))
		}

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().ModifyIPv6AddressesAttributesWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update elastic public ipv6 failed, reason:%+v", logId, err)
			return err
		}
	}

	needChange1 := false
	mutableArgs1 := []string{"internet_max_bandwidth_out"}
	for _, v := range mutableArgs1 {
		if d.HasChange(v) {
			needChange1 = true
			break
		}
	}

	if needChange1 {
		request1 := vpc.NewModifyIPv6AddressesBandwidthRequest()
		request1.IPv6AddressIds = []*string{helper.String(ipId)}

		if v, ok := d.GetOkExists("internet_max_bandwidth_out"); ok {
			request1.InternetMaxBandwidthOut = helper.IntInt64(v.(int))
		}

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().ModifyIPv6AddressesBandwidthWithContext(ctx, request1)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request1.GetAction(), request1.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update elastic public ipv6 failed, reason:%+v", logId, err)
			return err
		}
	}

	if d.HasChange("tags") {
		oldTags, newTags := d.GetChange("tags")
		replaceTags, deleteTags := svctag.DiffTags(oldTags.(map[string]interface{}), newTags.(map[string]interface{}))
		tcClient := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
		tagService := svctag.NewTagService(tcClient)
		region := tcClient.Region
		resourceName := tccommon.BuildTagResourceName("vpc", "eipv6", region, ipId)
		if err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags); err != nil {
			log.Printf("[CRITAL]%s update eipv6 tags failed: %+v", logId, err)
			return err
		}

	}

	return resourceTencentCloudElasticPublicIpv6Read(d, meta)
}

func resourceTencentCloudElasticPublicIpv6Delete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_elastic_public_ipv6.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	ipId := d.Id()

	var (
		request  = vpc.NewReleaseIPv6AddressesRequest()
		response = vpc.NewReleaseIPv6AddressesResponse()
	)

	request.IPv6AddressIds = []*string{helper.String(ipId)}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().ReleaseIPv6AddressesWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s delete elastic public ipv6 failed, reason:%+v", logId, err)
		return err
	}

	_ = response
	return nil
}

func resourceElasticPublicIpv6CreateStateRefreshFunc_0_0(ctx context.Context, taskId string) resource.StateRefreshFunc {
	var req *vpc.DescribeTaskResultRequest
	return func() (interface{}, string, error) {
		meta := tccommon.ProviderMetaFromContext(ctx)
		if meta == nil {
			return nil, "", fmt.Errorf("resource data can not be nil")
		}
		if req == nil {
			d := tccommon.ResourceDataFromContext(ctx)
			if d == nil {
				return nil, "", fmt.Errorf("resource data can not be nil")
			}
			_ = d
			req = vpc.NewDescribeTaskResultRequest()
			req.TaskId = helper.StrToUint64Point(taskId)

		}
		resp, err := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().DescribeTaskResultWithContext(ctx, req)
		if err != nil {
			return nil, "", err
		}
		if resp == nil || resp.Response == nil {
			return nil, "", nil
		}
		state := fmt.Sprintf("%v", *resp.Response.Result)
		return resp.Response, state, nil
	}
}
