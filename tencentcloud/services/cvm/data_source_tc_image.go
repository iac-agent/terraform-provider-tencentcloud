package cvm

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudImage() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudImageRead,

		Schema: map[string]*schema.Schema{
			"filter": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "One 或 more 名称/值 pairs 到 过滤器。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "键 的 过滤器，有效 keys: `镜像-ID`，`镜像-类型`，`镜像-名称`。",
						},
						"values": {
							Type:        schema.TypeList,
							Required:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "Values 的 过滤器。",
						},
					},
				},
			},
			"image_name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: tccommon.ValidateNameRegex,
				Description: "应用于腾讯云返回的图像列表的正则表达式字符串。 **注意**：它不是通配符，应该类似于 `image_name_regex = \"^CentOS\\s+6\\.8\\s+64\\w*\"`。",
			},
			"os_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: tccommon.ValidateNotEmpty,
				Description:  "A 字符串 到 apply 使用 fuzzy match 到 os_name attribute 在 镜像 列表 返回 通过 TencentCloud. **NOTE**: 当 os_name 是 提供，highest 优先级 是 applied 在 此 字段 instead 的 `image_name_regex`。",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "用于保存结果。",
			},
			"image_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "An 镜像 ID indicate uniqueness 的 certain 镜像， 其中 可以 是 用于instance creation 或 resetting。",
			},
			"image_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "名称 此 镜像。",
			},
		},
	}
}

func dataSourceTencentCloudImageRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_image.read")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	cvmService := CvmService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}

	filter := make(map[string][]string)
	filters, ok := d.GetOk("filter")
	if ok {
		for _, v := range filters.(*schema.Set).List() {
			vv := v.(map[string]interface{})
			name := vv["name"].(string)
			filter[name] = []string{}
			for _, vvv := range vv["values"].([]interface{}) {
				filter[name] = append(filter[name], vvv.(string))
			}
		}
	}

	var images []*cvm.Image
	var errRet error
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		images, errRet = cvmService.DescribeImagesByFilter(ctx, filter, "")
		if errRet != nil {
			return tccommon.RetryError(errRet, tccommon.InternalError)
		}
		return nil
	})
	if err != nil {
		return err
	}

	if len(images) == 0 {
		return errors.New("No image found")
	}

	var osName string
	if v, ok := d.GetOk("os_name"); ok {
		osName = v.(string)
	}

	var regImageName string
	var imageNameRegex *regexp.Regexp
	if v, ok := d.GetOk("image_name_regex"); ok {
		regImageName = v.(string)
		imageNameRegex, err = regexp.Compile(regImageName)
		if err != nil {
			return fmt.Errorf("image_name_regex format error,%s", err.Error())
		}
	}

	var resultImageId string
	images = sortImages(images)
	for _, image := range images {
		if osName != "" {
			if strings.Contains(strings.ToLower(*image.OsName), strings.ToLower(osName)) {
				resultImageId = *image.ImageId
				_ = d.Set("image_name", *image.ImageName)
				break
			}
			continue
		}

		if imageNameRegex != nil {
			if imageNameRegex.MatchString(*image.ImageName) {
				resultImageId = *image.ImageId
				_ = d.Set("image_name", *image.ImageName)
				break
			}
			continue
		}

		resultImageId = *image.ImageId
		_ = d.Set("image_name", *image.ImageName)
		break
	}

	if resultImageId == "" {
		return errors.New("No image found")
	}

	d.SetId(helper.DataResourceIdHash(resultImageId))
	if err := d.Set("image_id", resultImageId); err != nil {
		return err
	}

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if err = tccommon.WriteToFile(output.(string), resultImageId); err != nil {
			return err
		}
	}

	return nil
}
