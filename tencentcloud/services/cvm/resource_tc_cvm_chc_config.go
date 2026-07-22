package cvm

import (
	"context"
	"log"
	"time"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudCvmChcConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCvmChcConfigCreate,
		Update: resourceTencentCloudCvmChcConfigUpdate,
		Read:   resourceTencentCloudCvmChcConfigRead,
		Delete: resourceTencentCloudCvmChcConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"chc_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "CHC 主机 ID。",
			},

			"instance_name": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeString,
				Description: "CHC 主机名",
			},

			"device_type": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Server 类型",
			},

			"bmc_user": {
				Optional:     true,
				RequiredWith: []string{"password"},
				Type:         schema.TypeString,
				Description:  "有效 字符: Letters，numbers，hyphens 和 underscores. Only 集合 当 update 密码",
			},

			"password": {
				Optional:     true,
				Sensitive:    true,
				RequiredWith: []string{"bmc_user"},
				Type:         schema.TypeString,

				Description: "密码 可以 contain 8 到 16 字符，包括 letters，numbers 和 special symbols (()`~!@#$%^&amp;amp;*-+=_|{})。",
			},

			"bmc_virtual_private_cloud": {
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Out-的-band 网络 信息。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vpc_id": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "私有网络 ID 在 格式 的 vpc-xxx. To obtain 有效 VPC IDs，您 可以 日志 在 到 [console](https://console.tencentcloud.com/vpc/vpc?rid=1) 或 call DescribeVpcEx API 和 look 对于 unVpcId 字段 在 response. 如果 您 指定DEFAULT 对于 both VpcId 和 SubnetId 当 creating 实例， 默认值 VPC 将 是 使用。",
						},
						"subnet_id": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "VPC 子网 ID 在 格式 子网-xxx. To obtain 有效 子网 IDs，您 可以 日志 在 到 [console](https://console.tencentcloud.com/vpc/vpc?rid=1) 或 call DescribeSubnets 和 look 对于 unSubnetId 字段 在 response. 如果 您 指定DEFAULT 对于 both SubnetId 和 VpcId 当 creating 实例， 默认值 VPC 将 是 使用。",
						},
						"as_vpc_gateway": {
							Type:        schema.TypeBool,
							Optional:    true,
							ForceNew:    true,
							Description: "是否use CVM 实例 作为 公有 网关. 公有 网关 是 仅 可用 当 实例 has 公有 IP 和 resides 在 VPC. 有效 值:&lt;br&gt;&lt;li&gt;TRUE: yes;&lt;br&gt;&lt;li&gt;FALSE: 无&lt;br&gt;&lt;br&gt;默认值：FALSE。",
						},
						"private_ip_addresses": {
							Type: schema.TypeList,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Optional:    true,
							Computed:    true,
							ForceNew:    true,
							Description: "数组 VPC 子网 IPs. You 可以 使用 此 参数 当 creating 实例 或 modifying VPC attributes 的 实例. Currently 您 可以 指定multiple IPs 在 一个 子网 仅 当 creating 多个 实例 在 same 时间。",
						},
						"ipv6_address_count": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							ForceNew:    true,
							Description: "数量 IPv6 addresses randomly generated 对于 ENI。",
						},
					},
				},
			},

			"bmc_security_group_ids": {
				Optional: true,
				Computed: true,
				ForceNew: true,
				Type:     schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				RequiredWith: []string{"bmc_virtual_private_cloud"},
				Description:  "Out-的-band 网络 安全 组 列表。",
			},

			"deploy_virtual_private_cloud": {
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Deployment 网络 信息。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vpc_id": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "私有网络 ID 在 格式 的 vpc-xxx. To obtain 有效 VPC IDs，您 可以 日志 在 到 [console](https://console.tencentcloud.com/vpc/vpc?rid=1) 或 call DescribeVpcEx API 和 look 对于 unVpcId 字段 在 response. 如果 您 指定DEFAULT 对于 both VpcId 和 SubnetId 当 creating 实例， 默认值 VPC 将 是 使用。",
						},
						"subnet_id": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "VPC 子网 ID 在 格式 子网-xxx. To obtain 有效 子网 IDs，您 可以 日志 在 到 [console](https://console.tencentcloud.com/vpc/vpc?rid=1) 或 call DescribeSubnets 和 look 对于 unSubnetId 字段 在 response. 如果 您 指定DEFAULT 对于 both SubnetId 和 VpcId 当 creating 实例， 默认值 VPC 将 是 使用。",
						},
						"as_vpc_gateway": {
							Type:        schema.TypeBool,
							Optional:    true,
							ForceNew:    true,
							Description: "是否use CVM 实例 作为 公有 网关. 公有 网关 是 仅 可用 当 实例 has 公有 IP 和 resides 在 VPC. 有效 值:&lt;br&gt;&lt;li&gt;TRUE: yes;&lt;br&gt;&lt;li&gt;FALSE: 无&lt;br&gt;&lt;br&gt;默认值：FALSE。",
						},
						"private_ip_addresses": {
							Type: schema.TypeList,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Optional:    true,
							Computed:    true,
							ForceNew:    true,
							Description: "数组 VPC 子网 IPs. You 可以 使用 此 参数 当 creating 实例 或 modifying VPC attributes 的 实例. Currently 您 可以 指定multiple IPs 在 一个 子网 仅 当 creating 多个 实例 在 same 时间。",
						},
						"ipv6_address_count": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							ForceNew:    true,
							Description: "数量 IPv6 addresses randomly generated 对于 ENI。",
						},
					},
				},
			},

			"deploy_security_group_ids": {
				Optional: true,
				Computed: true,
				ForceNew: true,
				Type:     schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				RequiredWith: []string{"deploy_virtual_private_cloud"},
				Description:  "Deployment 网络 安全 组 列表。",
			},
		},
	}
}

func resourceTencentCloudCvmChcConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cvm_chc_config.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		assistChange  bool
		deployChange  bool
		chcId         string
		vpcId         string
		assistRequest = cvm.NewConfigureChcAssistVpcRequest()
		deployRequest = cvm.NewConfigureChcDeployVpcRequest()
	)
	service := CvmService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	if v, ok := d.GetOk("chc_id"); ok {
		chcId = v.(string)
	}

	if v, ok := d.GetOk("instance_name"); ok {
		attributeRequest := cvm.NewModifyChcAttributeRequest()
		attributeRequest.InstanceName = helper.String(v.(string))
		attributeRequest.ChcIds = []*string{&chcId}
		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCvmClient().ModifyChcAttribute(attributeRequest)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, attributeRequest.GetAction(), attributeRequest.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s operate cvm chcAttribute failed, reason:%+v", logId, err)
			return err
		}
	}

	if v, ok := d.GetOk("device_type"); ok {
		attributeRequest := cvm.NewModifyChcAttributeRequest()
		attributeRequest.DeviceType = helper.String(v.(string))
		attributeRequest.ChcIds = []*string{&chcId}
		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCvmClient().ModifyChcAttribute(attributeRequest)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, attributeRequest.GetAction(), attributeRequest.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s operate cvm chcAttribute failed, reason:%+v", logId, err)
			return err
		}
	}
	bmcUser, bmcUserok := d.GetOk("bmc_user")
	password, passwordOk := d.GetOk("password")
	if bmcUserok && passwordOk {
		attributeRequest := cvm.NewModifyChcAttributeRequest()
		attributeRequest.BmcUser = helper.String(bmcUser.(string))
		attributeRequest.Password = helper.String(password.(string))
		attributeRequest.ChcIds = []*string{&chcId}
		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCvmClient().ModifyChcAttribute(attributeRequest)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, attributeRequest.GetAction(), attributeRequest.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s operate cvm chcAttribute failed, reason:%+v", logId, err)
			return err
		}
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "bmc_virtual_private_cloud"); ok {
		virtualPrivateCloud := cvm.VirtualPrivateCloud{}
		if v, ok := dMap["vpc_id"]; ok {
			virtualPrivateCloud.VpcId = helper.String(v.(string))
		}
		if v, ok := dMap["subnet_id"]; ok {
			virtualPrivateCloud.SubnetId = helper.String(v.(string))
		}
		if v, ok := dMap["as_vpc_gateway"]; ok {
			virtualPrivateCloud.AsVpcGateway = helper.Bool(v.(bool))
		}
		if v, ok := dMap["private_ip_addresses"]; ok {
			privateIpAddresses := v.([]interface{})
			for i := range privateIpAddresses {
				privateIpAddresses := privateIpAddresses[i].(string)
				virtualPrivateCloud.PrivateIpAddresses = append(virtualPrivateCloud.PrivateIpAddresses, &privateIpAddresses)
			}
		}
		if v, ok := dMap["ipv6_address_count"]; ok {
			virtualPrivateCloud.Ipv6AddressCount = helper.IntUint64(v.(int))
		}
		assistChange = true
		assistRequest.BmcVirtualPrivateCloud = &virtualPrivateCloud
	}

	if v, ok := d.GetOk("bmc_security_group_ids"); ok {
		bmcSecurityGroupIds := v.([]interface{})
		for i := range bmcSecurityGroupIds {
			bmcSecurityGroupIds := bmcSecurityGroupIds[i].(string)
			assistRequest.BmcSecurityGroupIds = append(assistRequest.BmcSecurityGroupIds, &bmcSecurityGroupIds)
		}
		assistChange = true
	}

	if assistChange {
		assistRequest.ChcIds = []*string{&chcId}
		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCvmClient().ConfigureChcAssistVpc(assistRequest)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, assistRequest.GetAction(), assistRequest.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s create cvm chcAssistVpc failed, reason:%+v", logId, err)
			return err
		}
		conf := tccommon.BuildStateChangeConf([]string{}, []string{"READY"}, 20*tccommon.ReadRetryTimeout, time.Second, service.CvmChcInstanceStateRefreshFunc(chcId, []string{}))

		if _, e := conf.WaitForState(); e != nil {
			return e
		}
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "deploy_virtual_private_cloud"); ok {
		virtualPrivateCloud := cvm.VirtualPrivateCloud{}
		if v, ok := dMap["vpc_id"]; ok {
			vpcId = v.(string)
			virtualPrivateCloud.VpcId = helper.String(vpcId)
		}
		if v, ok := dMap["subnet_id"]; ok {
			virtualPrivateCloud.SubnetId = helper.String(v.(string))
		}
		if v, ok := dMap["as_vpc_gateway"]; ok {
			virtualPrivateCloud.AsVpcGateway = helper.Bool(v.(bool))
		}
		if v, ok := dMap["private_ip_addresses"]; ok {
			privateIpAddresses := v.([]interface{})
			for i := range privateIpAddresses {
				privateIpAddresses := privateIpAddresses[i].(string)
				virtualPrivateCloud.PrivateIpAddresses = append(virtualPrivateCloud.PrivateIpAddresses, &privateIpAddresses)
			}
		}
		if v, ok := dMap["ipv6_address_count"]; ok {
			virtualPrivateCloud.Ipv6AddressCount = helper.IntUint64(v.(int))
		}
		deployRequest.DeployVirtualPrivateCloud = &virtualPrivateCloud
		deployChange = true
	}

	if v, ok := d.GetOk("deploy_security_group_ids"); ok {
		deploySecurityGroupIds := v.([]interface{})
		for i := range deploySecurityGroupIds {
			deploySecurityGroupIds := deploySecurityGroupIds[i].(string)
			deployRequest.DeploySecurityGroupIds = append(deployRequest.DeploySecurityGroupIds, &deploySecurityGroupIds)
		}
		deployChange = true
	}

	if deployChange {
		deployRequest.ChcIds = []*string{&chcId}
		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCvmClient().ConfigureChcDeployVpc(deployRequest)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, deployRequest.GetAction(), deployRequest.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s create cvm chcDeployVpc failed, reason:%+v", logId, err)
			return err
		}

		conf := tccommon.BuildStateChangeConf([]string{}, []string{vpcId}, 10*tccommon.ReadRetryTimeout, time.Second, service.CvmChcInstanceDeployVpcStateRefreshFunc(chcId, []string{}))

		if _, e := conf.WaitForState(); e != nil {
			return e
		}
	}

	d.SetId(chcId)

	return resourceTencentCloudCvmChcConfigRead(d, meta)
}

func resourceTencentCloudCvmChcConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cvm_chc_config.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	chcId := d.Id()

	if d.HasChange("instance_name") {
		attributeRequest := cvm.NewModifyChcAttributeRequest()
		attributeRequest.ChcIds = []*string{&chcId}
		if v, ok := d.GetOk("instance_name"); ok {
			attributeRequest.InstanceName = helper.String(v.(string))
		}
		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCvmClient().ModifyChcAttribute(attributeRequest)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, attributeRequest.GetAction(), attributeRequest.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s operate cvm chcAttribute failed, reason:%+v", logId, err)
			return err
		}
	}
	if d.HasChange("device_type") {
		attributeRequest := cvm.NewModifyChcAttributeRequest()
		attributeRequest.ChcIds = []*string{&chcId}
		if v, ok := d.GetOk("device_type"); ok {
			attributeRequest.DeviceType = helper.String(v.(string))
		}
		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCvmClient().ModifyChcAttribute(attributeRequest)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, attributeRequest.GetAction(), attributeRequest.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s operate cvm chcAttribute failed, reason:%+v", logId, err)
			return err
		}
	}
	if d.HasChange("bmc_user") || d.HasChange("password") {
		attributeRequest := cvm.NewModifyChcAttributeRequest()
		attributeRequest.ChcIds = []*string{&chcId}
		if v, ok := d.GetOk("bmc_user"); ok {
			attributeRequest.BmcUser = helper.String(v.(string))
		}

		if v, ok := d.GetOk("password"); ok {
			attributeRequest.Password = helper.String(v.(string))
		}
		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCvmClient().ModifyChcAttribute(attributeRequest)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, attributeRequest.GetAction(), attributeRequest.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s operate cvm chcAttribute failed, reason:%+v", logId, err)
			return err
		}
	}

	return resourceTencentCloudCvmChcConfigRead(d, meta)
}
func resourceTencentCloudCvmChcConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cvm_chc_config.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := CvmService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	chcId := d.Id()

	params := map[string]interface{}{
		"chc_ids": []string{chcId},
	}
	chcHosts, err := service.DescribeCvmChcHostsByFilter(ctx, params)
	if err != nil {
		return err
	}

	if len(chcHosts) < 1 {
		d.SetId("")
		log.Printf("[WARN]%s resource `CvmChcAssistVpc` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	chcHost := chcHosts[0]
	if chcHost.ChcId != nil {
		_ = d.Set("chc_id", chcHost.ChcId)
	}
	_ = d.Set("instance_name", chcHost.InstanceName)
	_ = d.Set("device_type", chcHost.DeviceType)
	if chcHost.BmcVirtualPrivateCloud != nil {
		bmcVirtualPrivateCloudMap := map[string]interface{}{}

		if chcHost.BmcVirtualPrivateCloud.VpcId != nil {
			bmcVirtualPrivateCloudMap["vpc_id"] = chcHost.BmcVirtualPrivateCloud.VpcId
		}

		if chcHost.BmcVirtualPrivateCloud.SubnetId != nil {
			bmcVirtualPrivateCloudMap["subnet_id"] = chcHost.BmcVirtualPrivateCloud.SubnetId
		}

		if chcHost.BmcVirtualPrivateCloud.AsVpcGateway != nil {
			bmcVirtualPrivateCloudMap["as_vpc_gateway"] = chcHost.BmcVirtualPrivateCloud.AsVpcGateway
		}

		if chcHost.BmcVirtualPrivateCloud.PrivateIpAddresses != nil {
			privateIpAddresses := make([]string, 0)
			for _, p := range chcHost.BmcVirtualPrivateCloud.PrivateIpAddresses {
				privateIpAddresses = append(privateIpAddresses, *p)
			}
			bmcVirtualPrivateCloudMap["private_ip_addresses"] = privateIpAddresses
		}

		if chcHost.BmcVirtualPrivateCloud.Ipv6AddressCount != nil {
			bmcVirtualPrivateCloudMap["ipv6_address_count"] = chcHost.BmcVirtualPrivateCloud.Ipv6AddressCount
		}

		_ = d.Set("bmc_virtual_private_cloud", []interface{}{bmcVirtualPrivateCloudMap})
	}

	if chcHost.BmcSecurityGroupIds != nil {
		bmcSecurityGroupIds := make([]string, 0)
		for _, sgId := range chcHost.BmcSecurityGroupIds {
			bmcSecurityGroupIds = append(bmcSecurityGroupIds, *sgId)
		}
		_ = d.Set("bmc_security_group_ids", bmcSecurityGroupIds)
	}

	if chcHost.DeployVirtualPrivateCloud != nil {
		deployVirtualPrivateCloudMap := map[string]interface{}{}

		if chcHost.DeployVirtualPrivateCloud.VpcId != nil {
			deployVirtualPrivateCloudMap["vpc_id"] = chcHost.DeployVirtualPrivateCloud.VpcId
		}

		if chcHost.DeployVirtualPrivateCloud.SubnetId != nil {
			deployVirtualPrivateCloudMap["subnet_id"] = chcHost.DeployVirtualPrivateCloud.SubnetId
		}

		if chcHost.DeployVirtualPrivateCloud.AsVpcGateway != nil {
			deployVirtualPrivateCloudMap["as_vpc_gateway"] = chcHost.DeployVirtualPrivateCloud.AsVpcGateway
		}

		if chcHost.DeployVirtualPrivateCloud.PrivateIpAddresses != nil {
			privateIpAddresses := make([]string, 0)
			for _, p := range chcHost.DeployVirtualPrivateCloud.PrivateIpAddresses {
				privateIpAddresses = append(privateIpAddresses, *p)
			}
			deployVirtualPrivateCloudMap["private_ip_addresses"] = privateIpAddresses
		}

		if chcHost.DeployVirtualPrivateCloud.Ipv6AddressCount != nil {
			deployVirtualPrivateCloudMap["ipv6_address_count"] = chcHost.DeployVirtualPrivateCloud.Ipv6AddressCount
		}

		_ = d.Set("deploy_virtual_private_cloud", []interface{}{deployVirtualPrivateCloudMap})
	}

	if chcHost.DeploySecurityGroupIds != nil {
		deploySecurityGroupIds := make([]string, 0)
		for _, sgId := range chcHost.DeploySecurityGroupIds {
			deploySecurityGroupIds = append(deploySecurityGroupIds, *sgId)
		}
		_ = d.Set("deploy_security_group_ids", deploySecurityGroupIds)
	}

	return nil
}

func resourceTencentCloudCvmChcConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cvm_chc_config.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := CvmService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	chcId := d.Id()

	request := cvm.NewRemoveChcDeployVpcRequest()
	request.ChcIds = []*string{&chcId}
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCvmClient().RemoveChcDeployVpc(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s remove Chc deploy vpc failed, reason:%+v", logId, err)
		return err
	}

	conf := tccommon.BuildStateChangeConf([]string{}, []string{""}, 5*tccommon.ReadRetryTimeout, time.Second, service.CvmChcInstanceDeployVpcStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	params := map[string]interface{}{
		"chc_ids": []string{chcId},
	}
	chcHosts, err := service.DescribeCvmChcHostsByFilter(ctx, params)
	if err != nil {
		return err
	}
	if len(chcHosts) > 0 && *chcHosts[0].InstanceState == "INIT" {
		return nil
	}

	if err := service.DeleteCvmChcAssistVpcById(ctx, chcId); err != nil {
		return err
	}

	conf = tccommon.BuildStateChangeConf([]string{}, []string{"INIT"}, 10*tccommon.ReadRetryTimeout, time.Second, service.CvmChcInstanceStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return nil
}
