package mdl

import (
	"context"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mdl "github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/mdl/v20200326"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudMdlStreamLiveInput() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMdlStreamLiveInputCreate,
		Read:   resourceTencentCloudMdlStreamLiveInputRead,
		Update: resourceTencentCloudMdlStreamLiveInputUpdate,
		Delete: resourceTencentCloudMdlStreamLiveInputDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Input 名称，其中 可以 contain 1-32 case-sensitive letters，digits，和 underscores 和 必须 是 唯一 在 地域 级别",
			},

			"type": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Input type有效值：`RTMP_PUSH`，`RTP_PUSH`，`UDP_PUSH`，`RTMP_PULL`，`HLS_PULL`，`MP4_PULL`。",
			},

			"security_group_ids": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "ID input 安全 组 到 attachYou 可以 attach 仅 一个 安全 组 到 input。",
			},

			"input_settings": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Input settings. For 类型 `RTMP_PUSH`，`RTMP_PULL`，`HLS_PULL`，或 `MP4_PULL`，1 或 2 inputs 的 corresponding 类型 可以 是 已配置。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"app_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Application 名称，其中 是 有效 如果 `类型` 是 `RTMP_PUSH` 和 可以 contain 1-32 letters 和 digitsNote: 此 字段 可能 返回 `null`，indicating 该 无 有效 值 是 found。",
						},
						"stream_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Stream 名称，其中 是 有效 如果 `类型` 是 `RTMP_PUSH` 和 可以 contain 1-32 letters 和 digitsNote: 此 字段 可能 返回 `null`，indicating 该 无 有效 值 是 found。",
						},
						"source_url": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "来源 URL，其中 是 有效 如果 `类型` 是 `RTMP_PULL`，`HLS_PULL`，或 `MP4_PULL` 和 可以 contain 1-512 charactersNote: 此 字段 可能 返回 `null`，indicating 该 无 有效 值 是 found。",
						},
						"input_address": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "RTP/UDP input 地址，其中 does 不 need 到 是 entered 对于 input 参数.注意: 此 字段 可能 返回 null，indicating 该 无 有效 值 可以 是 获取。",
						},
						"source_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "来源 类型 对于 流 pulling 和 relaying. To pull 内容 从 私有-read COS buckets under 当前 账号，集合 此 参数 到 `TencentCOS`; otherwise，leave 它 空.注意: 此 字段 可能 返回 `null`，indicating 该 无 有效 值 是 found。",
						},
						"delay_time": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Delayed 时间 (ms) 对于 playback，其中 是 有效 如果 `类型` 是 `RTMP_PUSH`取值范围：0 (默认值) 或 10000-600000The 值 必须 是 多个 的 1,000.注意: 此 字段 可能 返回 `null`，indicating 该 无 有效 值 是 found。",
						},
						"input_domain": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "域名 的 SRT_PUSH 地址 如果 此 是 请求 参数，您 do 不 need 到 指定it.注意: 此 字段 可能 返回 `null`，indicating 该 无 有效 值 是 found。",
						},
						"user_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "用户名，其中 是 用于authentication.注意: 此 字段 可能 返回 `null`，indicating 该 无 有效 值 是 found。",
						},
						"password": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "密码，其中 是 用于authentication.注意: 此 字段 可能 返回 `null`，indicating 该 无 有效 值 是 found。",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudMdlStreamLiveInputCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mdl_stream_live_input.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request  = mdl.NewCreateStreamLiveInputRequest()
		response = mdl.NewCreateStreamLiveInputResponse()
		id       string
	)
	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}

	if v, ok := d.GetOk("type"); ok {
		request.Type = helper.String(v.(string))
	}

	if v, ok := d.GetOk("security_group_ids"); ok {
		securityGroupIdsSet := v.(*schema.Set).List()
		for i := range securityGroupIdsSet {
			securityGroupIds := securityGroupIdsSet[i].(string)
			request.SecurityGroupIds = append(request.SecurityGroupIds, &securityGroupIds)
		}
	}

	if v, ok := d.GetOk("input_settings"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			inputSettingInfo := mdl.InputSettingInfo{}
			if v, ok := dMap["app_name"]; ok {
				inputSettingInfo.AppName = helper.String(v.(string))
			}
			if v, ok := dMap["stream_name"]; ok {
				inputSettingInfo.StreamName = helper.String(v.(string))
			}
			if v, ok := dMap["source_url"]; ok {
				inputSettingInfo.SourceUrl = helper.String(v.(string))
			}
			if v, ok := dMap["input_address"]; ok {
				inputSettingInfo.InputAddress = helper.String(v.(string))
			}
			if v, ok := dMap["source_type"]; ok {
				inputSettingInfo.SourceType = helper.String(v.(string))
			}
			if v, ok := dMap["delay_time"]; ok {
				inputSettingInfo.DelayTime = helper.IntInt64(v.(int))
			}
			if v, ok := dMap["input_domain"]; ok {
				inputSettingInfo.InputDomain = helper.String(v.(string))
			}
			if v, ok := dMap["user_name"]; ok {
				inputSettingInfo.UserName = helper.String(v.(string))
			}
			if v, ok := dMap["password"]; ok {
				inputSettingInfo.Password = helper.String(v.(string))
			}
			request.InputSettings = append(request.InputSettings, &inputSettingInfo)
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMdlClient().CreateStreamLiveInput(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create mdl streamliveInput failed, reason:%+v", logId, err)
		return err
	}

	id = *response.Response.Id
	d.SetId(id)

	return resourceTencentCloudMdlStreamLiveInputRead(d, meta)
}

func resourceTencentCloudMdlStreamLiveInputRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mdl_stream_live_input.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := MdlService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	id := d.Id()

	streamliveInput, err := service.DescribeMdlStreamLiveInputById(ctx, id)
	if err != nil {
		return err
	}

	if streamliveInput == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `MdlStreamliveInput` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if streamliveInput.Name != nil {
		_ = d.Set("name", streamliveInput.Name)
	}

	if streamliveInput.Type != nil {
		_ = d.Set("type", streamliveInput.Type)
	}

	if streamliveInput.SecurityGroupIds != nil {
		_ = d.Set("security_group_ids", streamliveInput.SecurityGroupIds)
	}

	if streamliveInput.InputSettings != nil {
		inputSettingsList := []interface{}{}
		for _, inputSettings := range streamliveInput.InputSettings {
			inputSettingsMap := map[string]interface{}{}

			if inputSettings.AppName != nil {
				inputSettingsMap["app_name"] = inputSettings.AppName
			}

			if inputSettings.StreamName != nil {
				inputSettingsMap["stream_name"] = inputSettings.StreamName
			}

			if inputSettings.SourceUrl != nil {
				inputSettingsMap["source_url"] = inputSettings.SourceUrl
			}

			if inputSettings.InputAddress != nil {
				inputSettingsMap["input_address"] = inputSettings.InputAddress
			}

			if inputSettings.SourceType != nil {
				inputSettingsMap["source_type"] = inputSettings.SourceType
			}

			if inputSettings.DelayTime != nil {
				inputSettingsMap["delay_time"] = inputSettings.DelayTime
			}

			if inputSettings.InputDomain != nil {
				inputSettingsMap["input_domain"] = inputSettings.InputDomain
			}

			if inputSettings.UserName != nil {
				inputSettingsMap["user_name"] = inputSettings.UserName
			}

			if inputSettings.Password != nil {
				inputSettingsMap["password"] = inputSettings.Password
			}

			inputSettingsList = append(inputSettingsList, inputSettingsMap)
		}

		_ = d.Set("input_settings", inputSettingsList)

	}

	return nil
}

func resourceTencentCloudMdlStreamLiveInputUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mdl_streamlive_input.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	request := mdl.NewModifyStreamLiveInputRequest()

	id := d.Id()

	request.Id = &id

	needChange := false
	mutableArgs := []string{"name", "security_group_ids", "input_settings"}

	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}

	if needChange {

		if v, ok := d.GetOk("name"); ok {
			request.Name = helper.String(v.(string))
		}

		if v, ok := d.GetOk("security_group_ids"); ok {
			securityGroupIdsSet := v.(*schema.Set).List()
			for i := range securityGroupIdsSet {
				securityGroupIds := securityGroupIdsSet[i].(string)
				request.SecurityGroupIds = append(request.SecurityGroupIds, &securityGroupIds)
			}
		}

		if v, ok := d.GetOk("input_settings"); ok {
			for _, item := range v.([]interface{}) {
				dMap := item.(map[string]interface{})
				inputSettingInfo := mdl.InputSettingInfo{}
				if v, ok := dMap["app_name"]; ok {
					inputSettingInfo.AppName = helper.String(v.(string))
				}
				if v, ok := dMap["stream_name"]; ok {
					inputSettingInfo.StreamName = helper.String(v.(string))
				}
				if v, ok := dMap["source_url"]; ok {
					inputSettingInfo.SourceUrl = helper.String(v.(string))
				}
				if v, ok := dMap["input_address"]; ok {
					inputSettingInfo.InputAddress = helper.String(v.(string))
				}
				if v, ok := dMap["source_type"]; ok {
					inputSettingInfo.SourceType = helper.String(v.(string))
				}
				if v, ok := dMap["delay_time"]; ok {
					inputSettingInfo.DelayTime = helper.IntInt64(v.(int))
				}
				if v, ok := dMap["input_domain"]; ok {
					inputSettingInfo.InputDomain = helper.String(v.(string))
				}
				if v, ok := dMap["user_name"]; ok {
					inputSettingInfo.UserName = helper.String(v.(string))
				}
				if v, ok := dMap["password"]; ok {
					inputSettingInfo.Password = helper.String(v.(string))
				}
				request.InputSettings = append(request.InputSettings, &inputSettingInfo)
			}
		}

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMdlClient().ModifyStreamLiveInput(request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update mdl streamliveInput failed, reason:%+v", logId, err)
			return err
		}
	}

	return resourceTencentCloudMdlStreamLiveInputRead(d, meta)
}

func resourceTencentCloudMdlStreamLiveInputDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mdl_stream_live_input.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := MdlService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	id := d.Id()

	if err := service.DeleteMdlStreamLiveInputById(ctx, id); err != nil {
		return err
	}

	return nil
}
