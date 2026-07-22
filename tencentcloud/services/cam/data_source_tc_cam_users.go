package cam

import (
	"context"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cam "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cam/v20190116"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudCamUsers() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCamUsersRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "名称 CAM 用户 to be queried。",
			},
			"remark": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "备注 of the CAM 用户 to be queried。",
			},
			"phone_num": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Phone num of the CAM 用户 to be queried。",
			},
			"country_code": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Country 代码 of the CAM 用户 to be queried。",
			},
			"email": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Email of the CAM 用户 to be queried。",
			},
			"uin": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Uin of the CAM 用户 to be queried。",
			},
			"uid": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Uid of the CAM 用户 to be queried。",
			},
			"console_login": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Indicate 是否user can login in。",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "用于保存结果。",
			},
			"user_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A 列表 CAM users. Each element 包含following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"user_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID CAM 用户 Its 值 equals to `名称` argument。",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "名称 CAM 用户",
						},
						"remark": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "备注 of the CAM 用户",
						},
						"phone_num": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Phone num of the CAM 用户",
						},
						"country_code": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Country 代码 of the CAM 用户",
						},
						"email": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Email of the CAM 用户",
						},
						"uin": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Uin of the CAM 用户",
						},
						"uid": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Uid of the CAM 用户",
						},
						"console_login": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Indicate 是否user can login in。",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudCamUsersRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_cam_users.read")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	params := make(map[string]interface{})
	if v, ok := d.GetOk("name"); ok {
		params["name"] = v.(string)
	}
	if v, ok := d.GetOk("uin"); ok {
		params["uin"] = v.(int)
	}
	if v, ok := d.GetOk("remark"); ok {
		params["remark"] = v.(string)
	}
	if v, ok := d.GetOk("uid"); ok {
		params["uid"] = v.(int)
	}
	if v, ok := d.GetOk("phone_num"); ok {
		params["phone_num"] = v.(string)
	}
	if v, ok := d.GetOk("country_code"); ok {
		params["country_code"] = v.(string)
	}
	if v, ok := d.GetOk("email"); ok {
		params["email"] = v.(string)
	}
	if v, ok := d.GetOkExists("console_login"); ok {
		consoleLogin := v.(bool)
		if consoleLogin {
			params["console_login"] = 1
		} else {
			params["console_login"] = 0
		}
	}

	camService := CamService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}
	var users []*cam.SubAccountInfo
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		results, e := camService.DescribeUsersByFilter(ctx, params)
		if e != nil {
			return tccommon.RetryError(e)
		}
		users = results
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read CAM users failed, reason:%s\n", logId, err.Error())
		return err
	}
	userList := make([]map[string]interface{}, 0, len(users))
	ids := make([]string, 0, len(users))
	for _, user := range users {
		mapping := map[string]interface{}{
			"uin":          int(*user.Uin),
			"uid":          int(*user.Uid),
			"name":         *user.Name,
			"remark":       *user.Remark,
			"phone_num":    *user.PhoneNum,
			"country_code": *user.CountryCode,
			"email":        *user.Email,
			"user_id":      *user.Name,
		}
		if int(*user.ConsoleLogin) == 1 {
			mapping["console_login"] = true
		} else {
			mapping["console_login"] = false
		}
		userList = append(userList, mapping)
		ids = append(ids, *user.Name)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	if e := d.Set("user_list", userList); e != nil {
		log.Printf("[CRITAL]%s provider set CAM user list fail, reason:%s\n", logId, e.Error())
		return e
	}

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), userList); e != nil {
			return e
		}
	}

	return nil
}
