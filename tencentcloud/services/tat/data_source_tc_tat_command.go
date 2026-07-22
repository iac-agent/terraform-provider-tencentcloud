package tat

import (
	"context"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tat "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tat/v20201028"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudTatCommand() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTatCommandRead,
		Schema: map[string]*schema.Schema{
			"command_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "命令 ID",
			},

			"command_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "命令名称",
			},

			"command_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "命令类型，值 is `SHELL` or `POWERSHELL`。",
			},

			"created_by": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Command 创建者 `TAT` 表示a public command and `USER` 表示a personal command。",
			},

			"command_set": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "列表 command details。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"command_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "命令 ID",
						},
						"command_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "命令名称",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "命令描述",
						},
						"content": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "command。",
						},
						"command_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "命令类型",
						},
						"working_directory": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Command execution 路径",
						},
						"timeout": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Command timeout 周期",
						},
						"created_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Command 创建时间。",
						},
						"updated_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Command 更新时间。",
						},
						"enable_parameter": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "是否enable the custom parameter feature。",
						},
						"default_parameters": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Default custom parameter 值",
						},
						"formatted_description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Formatted 描述 command. This parameter is an empty string for 用户 commands and 包含values for public commands。",
						},
						"created_by": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Command 创建者 `TAT` 表示a public command and `USER` 表示a personal command。",
						},
						"tags": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "标签 bound to the command. At most 10 标签 are allowed。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "标签键",
									},
									"value": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "标签值",
									},
								},
							},
						},
						"username": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The 用户 who executes the command on the instance。",
						},
						"output_cos_bucket_url": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The COS 存储桶 URL for uploading logs。",
						},
						"output_cos_key_prefix": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The COS 存储桶 directory where the logs are saved。",
						},
					},
				},
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "用于保存结果。",
			},
		},
	}
}

func dataSourceTencentCloudTatCommandRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_tat_command.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("command_id"); ok {
		paramMap["command_id"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("command_name"); ok {
		paramMap["command_name"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("command_type"); ok {
		paramMap["command_type"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("created_by"); ok {
		paramMap["created_by"] = helper.String(v.(string))
	}

	tatService := TatService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var commandSet []*tat.Command
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		results, e := tatService.DescribeTatCommandByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		commandSet = results
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read Tat commandSet failed, reason:%+v", logId, err)
		return err
	}

	commandSetList := []interface{}{}
	ids := make([]string, 0, len(commandSet))
	if commandSet != nil {
		for _, commandSet := range commandSet {
			commandSetMap := map[string]interface{}{}
			if commandSet.CommandId != nil {
				commandSetMap["command_id"] = commandSet.CommandId
			}
			if commandSet.CommandName != nil {
				commandSetMap["command_name"] = commandSet.CommandName
			}
			if commandSet.Description != nil {
				commandSetMap["description"] = commandSet.Description
			}
			if commandSet.Content != nil {
				content, err := tccommon.Base64ToString(*commandSet.Content)
				if err == nil {
					commandSetMap["content"] = content
				} else {
					log.Printf("[CRITAL]%s base64 to string failed, err:%+v", logId, err)
				}
			}
			if commandSet.CommandType != nil {
				commandSetMap["command_type"] = commandSet.CommandType
			}
			if commandSet.WorkingDirectory != nil {
				commandSetMap["working_directory"] = commandSet.WorkingDirectory
			}
			if commandSet.Timeout != nil {
				commandSetMap["timeout"] = commandSet.Timeout
			}
			if commandSet.CreatedTime != nil {
				commandSetMap["created_time"] = commandSet.CreatedTime
			}
			if commandSet.UpdatedTime != nil {
				commandSetMap["updated_time"] = commandSet.UpdatedTime
			}
			if commandSet.EnableParameter != nil {
				commandSetMap["enable_parameter"] = commandSet.EnableParameter
			}
			if commandSet.DefaultParameters != nil {
				commandSetMap["default_parameters"] = commandSet.DefaultParameters
			}
			if commandSet.FormattedDescription != nil {
				commandSetMap["formatted_description"] = commandSet.FormattedDescription
			}
			if commandSet.CreatedBy != nil {
				commandSetMap["created_by"] = commandSet.CreatedBy
			}
			if commandSet.Tags != nil {
				tagsList := []interface{}{}
				for _, tags := range commandSet.Tags {
					tagsMap := map[string]interface{}{}
					if tags.Key != nil {
						tagsMap["key"] = tags.Key
					}
					if tags.Value != nil {
						tagsMap["value"] = tags.Value
					}

					tagsList = append(tagsList, tagsMap)
				}
				commandSetMap["tags"] = tagsList
			}
			if commandSet.Username != nil {
				commandSetMap["username"] = commandSet.Username
			}
			if commandSet.OutputCOSBucketUrl != nil {
				commandSetMap["output_cos_bucket_url"] = commandSet.OutputCOSBucketUrl
			}
			if commandSet.OutputCOSKeyPrefix != nil {
				commandSetMap["output_cos_key_prefix"] = commandSet.OutputCOSKeyPrefix
			}

			commandSetList = append(commandSetList, commandSetMap)
			ids = append(ids, *commandSet.CommandId)
		}
		_ = d.Set("command_set", commandSetList)
	}

	d.SetId(helper.DataResourceIdsHash(ids))

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), commandSetList); e != nil {
			return e
		}
	}

	return nil
}
