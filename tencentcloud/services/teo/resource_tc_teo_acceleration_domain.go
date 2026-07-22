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
				Description: "ID site related 使用 accelerated 域名 名称",
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
				Description: "Details 的 源站。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"origin_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Origin 服务器 类型，使用 值: IP_DOMAIN: IPv4，IPv6，或 域名 名称 类型 源站 服务器; COS: Tencent Cloud COS 源站 服务器; AWS_S3: AWS S3 源站 服务器; ORIGIN_GROUP: 源站 服务器 组 类型 源站 服务器; VOD: Video 在 Demand; SPACE: 源站 服务器 uninstallation. Currently 仅 可用 到 allowlist; LB: load balancing. Currently 仅 可用 到 allowlist。",
						},
						"origin": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Origin 服务器 地址，其中 varies according 到 值 的 OriginType: 当 OriginType = IP_DOMAIN，fill 在 IPv4 地址， IPv6 地址，或 域名 名称; 当 OriginType = COS，fill 在 访问 域名 名称 COS 存储桶; 当 OriginType = AWS_S3，fill 在 访问 域名 名称 S3 存储桶; 当 OriginType = ORIGIN_GROUP，fill 在 源站 服务器 组 ID; 当 OriginType = VOD，fill 在 VOD 应用 ID; 当 OriginType = LB，fill 在 Cloud Load Balancer 实例 ID. 此 功能 是 currently 仅 可用 到 allowlist; 当 OriginType = SPACE，fill 在 源站 服务器 uninstallation space ID. 此 功能 是 currently 仅 可用 到 allowlist。",
						},
						"backup_origin": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "ID secondary 源站 组. 此 参数 是 有效 仅 当 OriginType 是 ORIGIN_GROUP. 此 字段 表示old 版本 capability，其中 不能 是 已配置 或 modified 在 control panel after being called. Please 提交 ticket 如果 必填",
						},
						"private_access": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Whether 访问 到 私有 Cloud Object Storage 源站 服务器 是 allowed. 此 参数 是 有效 仅 当 OriginType 是 COS 或 AWS_S3. 有效值：在: Enable 私有 authentication; 关闭: Disable 私有 authentication. 如果 它 是 不 指定， 默认值为 关闭。",
						},
						"private_parameters": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Private authentication 参数. 此 参数 是 有效 仅 当 `private_access` 是 在。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "参数 名称 有效值：`AccessKeyId`: Access 键 ID; `SecretAccessKey`: Secret Access 键; `SignatureVersion`: authentication 版本，v2 或 v4; `地域`: 存储桶 地域",
									},
									"value": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "参数 值",
									},
								},
							},
						},
						"host_header": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Custom 源站 服务器 HOST 头部. 此 参数 是 有效 仅 当 OriginType=IP_DOMAIN.如果 OriginType 是 another 类型 源站，此 参数 does 不 need 到 是 passed 在，otherwise 错误 将 是 reported. 如果 OriginType 是 COS 或 AWS_S3， HOST 头部 对于 源站-pull 将 remain consistent 使用 源站 服务器 域名 名称 如果 OriginType 是 ORIGIN_GROUP， HOST 头部 follows ORIGIN site GROUP 配置. 如果 不 已配置，它 默认为 acceleration 域名 名称 如果 OriginType 是 VOD 或 SPACE，无 配置 为必填项 对于 此 头部，和 域名 名称 takes effect based 在 corresponding 源站。",
						},
						"vod_origin_scope": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "范围 的 云 在-demand back-到-来源 此 参数 是 effective 当 OriginType = VOD. possible 值 是: all: all files 在 云 在-demand 应用 corresponding 到 当前 源站 station. 默认值为 all; 存储桶: files 在 指定 存储桶 under 云 在-demand 应用 corresponding 到 当前 源站 station. 存储桶 是 指定 通过 参数 VodBucketId。",
						},
						"vod_bucket_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "VOD 存储桶 ID. 此 参数 为必填项 当 OriginType = VOD 和 VodOriginScope = 存储桶 Data 来源: 存储 ID 存储桶 在 Cloud VOD Professional Edition 应用。",
						},
					},
				},
			},

			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue([]string{"online", "offline"}),
				Description:  "Accelerated 域名 名称 状态， 值 是: `online`: 已启用; `offline`: 已禁用 默认为 `online`。",
			},

			"origin_protocol": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Origin 返回 协议，possible 值 是: `FOLLOW`: 协议 follow; `HTTP`: HTTP 协议 back 到 来源; `HTTPS`: HTTPS 协议 back 到 来源 如果未填写 在， 默认值 是: `FOLLOW`。",
			},

			"http_origin_port": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "HTTP back-到-源站 端口， 值 是 1-65535，effective 当 OriginProtocol=FOLLOW/HTTP，如果未填写 在， 默认值为 80。",
			},

			"https_origin_port": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "HTTPS back-到-源站 端口 值 范围 是 1-65535. It takes effect 当 OriginProtocol=FOLLOW/HTTPS. 如果 它 是 不 filled 在， 默认值为 443。",
			},

			"ipv6_status": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "IPv6 状态， 值 是: `follow`: follow site IPv6 配置; `在`: 在; `关闭`: 关闭. 如果未填写 在， 默认值 是: `follow`。",
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
