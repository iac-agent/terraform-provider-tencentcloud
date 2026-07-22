package teo

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	teo "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudTeoAccelerationDomain() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTeoAccelerationDomainCreate,
		Read:   resourceTencentCloudTeoAccelerationDomainRead,
		Update: resourceTencentCloudTeoAccelerationDomainUpdate,
		Delete: resourceTencentCloudTeoAccelerationDomainDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Read:   schema.DefaultTimeout(3 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID site related with the accelerated 域名 名称",
			},

			"domain_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Accelerated 域名 名称",
			},

			"origin_info": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Description: "Details of the origin。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"origin_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Origin server 类型，with values: IP_DOMAIN: IPv4，IPv6，or 域名 名称 类型 origin server; COS: Tencent Cloud COS origin server; AWS_S3: AWS S3 origin server; ORIGIN_GROUP: origin server group 类型 origin server; VOD: Video on Demand; SPACE: origin server uninstallation. Currently only available to the allowlist; LB: load balancing. Currently only available to the allowlist。",
						},
						"origin": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Origin server 地址，which varies according to the 值 of OriginType: When OriginType = IP_DOMAIN，fill in an IPv4 地址，an IPv6 地址，or a 域名 名称; When OriginType = COS，fill in the access 域名 名称 COS 存储桶; When OriginType = AWS_S3，fill in the access 域名 名称 S3 存储桶; When OriginType = ORIGIN_GROUP，fill in the origin server 组 ID; When OriginType = VOD，fill in the VOD application ID; When OriginType = LB，fill in the Cloud Load Balancer instance ID. This feature is currently only available to the allowlist; When OriginType = SPACE，fill in the origin server uninstallation space ID. This feature is currently only available to the allowlist。",
						},
						"backup_origin": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The ID secondary origin group. This parameter is valid only when OriginType is ORIGIN_GROUP. This field 表示old 版本 capability，which cannot be configured or modified on the control panel after being called. Please 提交 a ticket if 必填",
						},
						"private_access": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Whether access to the private Cloud Object Storage origin server is allowed. This parameter is valid only when OriginType is COS or AWS_S3. 有效值：on: Enable private authentication; off: Disable private authentication. If it is not specified，the 默认值为 off。",
						},
						"private_parameters": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Private authentication parameter. This parameter is valid only when `private_access` is on。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "The parameter 名称 有效值：`AccessKeyId`: Access 键 ID; `SecretAccessKey`: Secret Access 键; `SignatureVersion`: authentication 版本，v2 or v4; `地域`: 存储桶 地域",
									},
									"value": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "The parameter 值",
									},
								},
							},
						},
						"host_header": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Custom origin server HOST header. this parameter is valid only when OriginType=IP_DOMAIN.If the OriginType is another 类型 origin，this parameter does not need to be passed in，otherwise an 错误 will be reported. If OriginType is COS or AWS_S3，the HOST header for origin-pull will remain consistent with the origin server 域名 名称 If OriginType is ORIGIN_GROUP，the HOST header follows the ORIGIN site GROUP configuration. if not configured，it 默认为 the acceleration 域名 名称 If OriginType is VOD or SPACE，no configuration 为必填项 for this header，and the 域名 名称 takes effect based on the corresponding origin。",
						},
						"vod_origin_scope": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The 范围 of cloud on-demand back-to-来源 This parameter is effective when OriginType = VOD. The possible values are: all: all files in the cloud on-demand application corresponding to the current origin station. The 默认值为 all; 存储桶: files in a specified 存储桶 under the cloud on-demand application corresponding to the current origin station. The 存储桶 is specified by the parameter VodBucketId。",
						},
						"vod_bucket_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "VOD 存储桶 ID. This parameter 为必填项 when OriginType = VOD and VodOriginScope = 存储桶 Data 来源: the storage ID 存储桶 in the Cloud VOD Professional Edition application。",
						},
					},
				},
			},

			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue([]string{"online", "offline"}),
				Description:  "Accelerated 域名 名称 状态，the values are: `online`: 已启用; `offline`: 已禁用 默认为 `online`。",
			},

			"origin_protocol": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Origin return 协议，possible values are: `FOLLOW`: 协议 follow; `HTTP`: HTTP 协议 back to 来源; `HTTPS`: HTTPS 协议 back to 来源 如果未填写 in，the default is: `FOLLOW`。",
			},

			"http_origin_port": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "HTTP back-to-origin 端口，the 值 is 1-65535，effective when OriginProtocol=FOLLOW/HTTP，如果未填写 in，the 默认值为 80。",
			},

			"https_origin_port": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "HTTPS back-to-origin 端口 The 值 range is 1-65535. It takes effect when OriginProtocol=FOLLOW/HTTPS. If it is not filled in，the 默认值为 443。",
			},

			"ipv6_status": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "IPv6 状态，the 值 is: `follow`: follow the site IPv6 configuration; `on`: on; `off`: off. 如果未填写 in，the default is: `follow`。",
			},

			"cname": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "CNAME 地址",
			},
		},
	}
}

func resourceTencentCloudTeoAccelerationDomainCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_acceleration_domain.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId = tccommon.GetLogId(tccommon.ContextNil)
		ctx   = tccommon.NewResourceLifeCycleHandleFuncContext(
			context.Background(), logId, d, meta)
		request    = teo.NewCreateAccelerationDomainRequest()
		response   = teo.NewCreateAccelerationDomainResponse()
		zoneId     string
		domainName string
	)

	if v, ok := d.GetOk("zone_id"); ok {
		request.ZoneId = helper.String(v.(string))
		zoneId = v.(string)
	}

	if v, ok := d.GetOk("domain_name"); ok {
		request.DomainName = helper.String(v.(string))
		domainName = v.(string)
	}

	if originInfoMap, ok := helper.InterfacesHeadMap(d, "origin_info"); ok {
		originInfo := teo.OriginInfo{}
		var originType string
		if v, ok := originInfoMap["origin_type"]; ok {
			originInfo.OriginType = helper.String(v.(string))
			originType = v.(string)
		}

		if v, ok := originInfoMap["origin"]; ok {
			originInfo.Origin = helper.String(v.(string))
		}

		if v, ok := originInfoMap["backup_origin"]; ok {
			originInfo.BackupOrigin = helper.String(v.(string))
		}

		if v, ok := originInfoMap["private_access"]; ok {
			originInfo.PrivateAccess = helper.String(v.(string))
		}

		if v, ok := originInfoMap["private_parameters"]; ok {
			for _, item := range v.([]interface{}) {
				privateParametersMap := item.(map[string]interface{})
				privateParameter := teo.PrivateParameter{}
				if v, ok := privateParametersMap["name"]; ok {
					privateParameter.Name = helper.String(v.(string))
				}

				if v, ok := privateParametersMap["value"]; ok {
					privateParameter.Value = helper.String(v.(string))
				}

				originInfo.PrivateParameters = append(originInfo.PrivateParameters, &privateParameter)
			}
		}

		if v, ok := originInfoMap["host_header"].(string); ok && v != "" {
			if originType == "IP_DOMAIN" {
				originInfo.HostHeader = helper.String(v)
			} else {
				return fmt.Errorf("Only `origin_type` is `IP_DOMAIN` can set `host_header`.")
			}
		}

		if v, ok := originInfoMap["vod_origin_scope"].(string); ok && v != "" {
			originInfo.VodOriginScope = helper.String(v)
		}

		if v, ok := originInfoMap["vod_bucket_id"].(string); ok && v != "" {
			originInfo.VodBucketId = helper.String(v)
		}

		request.OriginInfo = &originInfo
	}

	if v, ok := d.GetOk("origin_protocol"); ok {
		request.OriginProtocol = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("http_origin_port"); ok {
		request.HttpOriginPort = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("https_origin_port"); ok {
		request.HttpsOriginPort = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("ipv6_status"); ok {
		request.IPv6Status = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoClient().CreateAccelerationDomainWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create teo acceleration domain failed, Response is nil."))
		}

		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create teo acceleration domain failed, reason:%+v", logId, err)
		return err
	}

	// wait
	if err := resourceTencentCloudTeoAccelerationDomainCreatePostHandleResponse0(ctx, response); err != nil {
		return err
	}

	d.SetId(strings.Join([]string{zoneId, domainName}, tccommon.FILED_SP))

	// offline
	if v, ok := d.GetOk("status"); ok {
		if v.(string) == "offline" {
			request := teo.NewModifyAccelerationDomainStatusesRequest()
			request.ZoneId = helper.String(zoneId)
			request.DomainNames = []*string{helper.String(domainName)}
			if v, ok := d.GetOk("status"); ok {
				request.Status = helper.String(v.(string))
			}

			err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
				result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoClient().ModifyAccelerationDomainStatusesWithContext(ctx, request)
				if e != nil {
					return tccommon.RetryError(e)
				} else {
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
				}

				return nil
			})

			if err != nil {
				log.Printf("[CRITAL]%s update teo acceleration domain status failed, reason:%+v", logId, err)
				return err
			}

			// wait
			if err = resourceTencentCloudTeoAccelerationDomainUpdateOnExit(ctx); err != nil {
				return err
			}
		}
	}

	return resourceTencentCloudTeoAccelerationDomainRead(d, meta)
}

func resourceTencentCloudTeoAccelerationDomainRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_acceleration_domain.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId = tccommon.GetLogId(tccommon.ContextNil)
		ctx   = tccommon.NewResourceLifeCycleHandleFuncContext(
			context.Background(), logId, d, meta)
		service = TeoService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	zoneId := idSplit[0]
	domainName := idSplit[1]

	_ = d.Set("zone_id", zoneId)
	_ = d.Set("domain_name", domainName)

	respData, err := service.DescribeTeoAccelerationDomainById(ctx, zoneId, domainName)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[WARN]%s resource `teo_acceleration_domain` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	if respData.ZoneId != nil {
		_ = d.Set("zone_id", respData.ZoneId)
	}

	if respData.DomainName != nil {
		_ = d.Set("domain_name", respData.DomainName)
	}

	if respData.OriginDetail != nil {
		originDetailMap := map[string]interface{}{}
		if respData.OriginDetail.OriginType != nil {
			originDetailMap["origin_type"] = respData.OriginDetail.OriginType
		}

		if respData.OriginDetail.Origin != nil {
			originDetailMap["origin"] = respData.OriginDetail.Origin
		}

		if respData.OriginDetail.BackupOrigin != nil {
			originDetailMap["backup_origin"] = respData.OriginDetail.BackupOrigin
		}

		if respData.OriginDetail.PrivateAccess != nil {
			originDetailMap["private_access"] = respData.OriginDetail.PrivateAccess
		}

		if respData.OriginDetail.PrivateParameters != nil {
			privateParametersList := make([]map[string]interface{}, 0, len(respData.OriginDetail.PrivateParameters))
			for _, privateParameters := range respData.OriginDetail.PrivateParameters {
				privateParametersMap := map[string]interface{}{}
				if privateParameters.Name != nil {
					privateParametersMap["name"] = privateParameters.Name
				}

				if privateParameters.Value != nil {
					privateParametersMap["value"] = privateParameters.Value
				}

				privateParametersList = append(privateParametersList, privateParametersMap)
			}

			originDetailMap["private_parameters"] = privateParametersList
		}

		if respData.OriginDetail.HostHeader != nil {
			originDetailMap["host_header"] = respData.OriginDetail.HostHeader
		}

		if respData.OriginDetail.VodOriginScope != nil {
			originDetailMap["vod_origin_scope"] = respData.OriginDetail.VodOriginScope
		}

		if respData.OriginDetail.VodBucketId != nil {
			originDetailMap["vod_bucket_id"] = respData.OriginDetail.VodBucketId
		}

		_ = d.Set("origin_info", []interface{}{originDetailMap})
	}

	if respData.DomainStatus != nil {
		_ = d.Set("status", respData.DomainStatus)
	}

	if respData.OriginProtocol != nil {
		_ = d.Set("origin_protocol", respData.OriginProtocol)
	}

	if respData.HttpOriginPort != nil {
		_ = d.Set("http_origin_port", respData.HttpOriginPort)
	}

	if respData.HttpsOriginPort != nil {
		_ = d.Set("https_origin_port", respData.HttpsOriginPort)
	}

	if respData.IPv6Status != nil {
		_ = d.Set("ipv6_status", respData.IPv6Status)
	}

	if respData.Cname != nil {
		_ = d.Set("cname", respData.Cname)
	}

	return nil
}

func resourceTencentCloudTeoAccelerationDomainUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_acceleration_domain.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId = tccommon.GetLogId(tccommon.ContextNil)
		ctx   = tccommon.NewResourceLifeCycleHandleFuncContext(
			context.Background(), logId, d, meta)
	)

	immutableArgs := []string{"https_origin_port"}
	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	zoneId := idSplit[0]
	domainName := idSplit[1]

	needChange := false
	mutableArgs := []string{"origin_info", "origin_protocol", "http_origin_port", "https_origin_port", "ipv6_status"}
	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}

	if needChange {
		request := teo.NewModifyAccelerationDomainRequest()
		request.ZoneId = helper.String(zoneId)
		request.DomainName = helper.String(domainName)
		if originInfoMap, ok := helper.InterfacesHeadMap(d, "origin_info"); ok {
			originInfo := teo.OriginInfo{}
			var originType string
			if v, ok := originInfoMap["origin_type"]; ok {
				originInfo.OriginType = helper.String(v.(string))
				originType = v.(string)
			}

			if v, ok := originInfoMap["origin"]; ok {
				originInfo.Origin = helper.String(v.(string))
			}

			if v, ok := originInfoMap["backup_origin"]; ok {
				originInfo.BackupOrigin = helper.String(v.(string))
			}

			if v, ok := originInfoMap["private_access"]; ok {
				originInfo.PrivateAccess = helper.String(v.(string))
			}

			if v, ok := originInfoMap["private_parameters"]; ok {
				for _, item := range v.([]interface{}) {
					privateParametersMap := item.(map[string]interface{})
					privateParameter := teo.PrivateParameter{}
					if v, ok := privateParametersMap["name"]; ok {
						privateParameter.Name = helper.String(v.(string))
					}

					if v, ok := privateParametersMap["value"]; ok {
						privateParameter.Value = helper.String(v.(string))
					}

					originInfo.PrivateParameters = append(originInfo.PrivateParameters, &privateParameter)
				}
			}

			if v, ok := originInfoMap["host_header"].(string); ok && v != "" {
				if originType == "IP_DOMAIN" {
					originInfo.HostHeader = helper.String(v)
				}
			}

			if v, ok := originInfoMap["vod_origin_scope"].(string); ok && v != "" {
				originInfo.VodOriginScope = helper.String(v)
			}

			if v, ok := originInfoMap["vod_bucket_id"].(string); ok && v != "" {
				originInfo.VodBucketId = helper.String(v)
			}

			request.OriginInfo = &originInfo
		}

		if v, ok := d.GetOk("origin_protocol"); ok {
			request.OriginProtocol = helper.String(v.(string))
		}

		if v, ok := d.GetOkExists("http_origin_port"); ok {
			request.HttpOriginPort = helper.IntUint64(v.(int))
		}

		if v, ok := d.GetOkExists("https_origin_port"); ok {
			request.HttpsOriginPort = helper.IntUint64(v.(int))
		}

		if v, ok := d.GetOk("ipv6_status"); ok {
			request.IPv6Status = helper.String(v.(string))
		}

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoClient().ModifyAccelerationDomainWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			return nil
		})

		if err != nil {
			log.Printf("[CRITAL]%s update teo acceleration domain failed, reason:%+v", logId, err)
			return err
		}

		// wait
		if err := resourceTencentCloudTeoAccelerationDomainUpdateOnExit(ctx); err != nil {
			return err
		}
	}

	if d.HasChange("status") {
		request := teo.NewModifyAccelerationDomainStatusesRequest()
		request.ZoneId = helper.String(zoneId)
		request.DomainNames = []*string{helper.String(domainName)}
		if v, ok := d.GetOk("status"); ok {
			request.Status = helper.String(v.(string))
		}

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoClient().ModifyAccelerationDomainStatusesWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			return nil
		})

		if err != nil {
			log.Printf("[CRITAL]%s update teo acceleration domain status failed, reason:%+v", logId, err)
			return err
		}

		// wait
		if err := resourceTencentCloudTeoAccelerationDomainUpdateOnExit(ctx); err != nil {
			return err
		}
	}

	return resourceTencentCloudTeoAccelerationDomainRead(d, meta)
}

func resourceTencentCloudTeoAccelerationDomainDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_acceleration_domain.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId = tccommon.GetLogId(tccommon.ContextNil)
		ctx   = tccommon.NewResourceLifeCycleHandleFuncContext(
			context.Background(), logId, d, meta)
		request  = teo.NewModifyAccelerationDomainStatusesRequest()
		response = teo.NewModifyAccelerationDomainStatusesResponse()
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	zoneId := idSplit[0]
	domainName := idSplit[1]

	// check offline first
	if v, ok := d.GetOk("status"); ok {
		if v.(string) == "online" {
			request.ZoneId = helper.String(zoneId)
			request.DomainNames = []*string{helper.String(domainName)}
			request.Status = helper.String("offline")
			err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
				result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoClient().ModifyAccelerationDomainStatusesWithContext(ctx, request)
				if e != nil {
					return tccommon.RetryError(e)
				} else {
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
				}

				if result == nil || result.Response == nil {
					return resource.NonRetryableError(fmt.Errorf("Modify teo acceleration domain status failed, Response is nil."))
				}

				response = result
				return nil
			})

			if err != nil {
				log.Printf("[CRITAL]%s modify teo acceleration domain status failed, reason:%+v", logId, err)
				return err
			}

			// wait
			if err := resourceTencentCloudTeoAccelerationDomainDeletePostHandleResponse0(ctx, response); err != nil {
				return err
			}
		}
	}

	// delete
	delRequest := teo.NewDeleteAccelerationDomainsRequest()
	delRequest.ZoneId = helper.String(zoneId)
	delRequest.DomainNames = []*string{helper.String(domainName)}
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoClient().DeleteAccelerationDomainsWithContext(ctx, delRequest)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, delRequest.GetAction(), delRequest.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Delete teo acceleration domain failed, Response is nil."))
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s delete teo acceleration domain failed, reason:%+v", logId, err)
		return err
	}

	return nil
}
