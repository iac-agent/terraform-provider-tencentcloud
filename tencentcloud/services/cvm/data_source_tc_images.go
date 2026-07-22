package cvm

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"

	svccbs "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/cbs"
)

func DataSourceTencentCloudImages() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudImagesRead,

		Schema: map[string]*schema.Schema{
			"image_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID 镜像 到 是 queried。",
			},
			"image_type": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "A 列表 镜像 类型 到 是 queried. 有效值：'PUBLIC_IMAGE'，'PRIVATE_IMAGE'，'SHARED_IMAGE'，'MARKET_IMAGE'。",
			},
			"image_name_regex": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"os_name"},
				ValidateFunc:  tccommon.ValidateNameRegex,
				Description: "应用于腾讯云返回的图像列表的正则表达式字符串，与“os_name”冲突。 **注意**：它不是通配符，应该类似于 `image_name_regex = \"^CentOS\\s+6\\.8\\s+64\\w*\"`。",
			},
			"os_name": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"image_name_regex"},
				ValidateFunc:  tccommon.ValidateNotEmpty,
				Description:   "A 字符串 到 apply 使用 fuzzy match 到 os_name attribute 在 镜像 列表 返回 通过 TencentCloud，conflict 使用 'image_name_regex'。",
			},
			"instance_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "实例类型，such 作为 `S1.SMALL1`。",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "用于保存结果。",
			},
			"images": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "An 信息 列表 镜像. Each element 包含following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"image_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID 镜像。",
						},
						"os_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "OS 名称 镜像。",
						},
						"image_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "类型 镜像。",
						},
						"created_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Created 时间 的 镜像。",
						},
						"image_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "名称 镜像。",
						},
						"image_description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "描述 镜像。",
						},
						"image_size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Size 的 镜像。",
						},
						"architecture": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Architecture 的 镜像。",
						},
						"image_state": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "State 的 镜像。",
						},
						"platform": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Platform 的 镜像。",
						},
						"image_creator": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Image 创建者 的 镜像。",
						},
						"image_source": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Image 来源 的 镜像。",
						},
						"sync_percent": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Sync percent 的 镜像。",
						},
						"support_cloud_init": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether support 云-init。",
						},
						"snapshots": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "列表 快照 details。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"snapshot_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Snapshot ID。",
									},
									"snapshot_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Snapshot 名称， 用户-defined 快照 alias。",
									},
									"disk_usage": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "类型 云 磁盘 用于create 快照。",
									},
									"disk_size": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Size 的 云 磁盘 用于create 快照; 单位: GB。",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudImagesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_images.read")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	cvmService := CvmService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}

	cbsService := svccbs.NewCbsService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())

	var (
		imageId        string
		imageType      []string
		imageName      string
		osName         string
		imageNameRegex *regexp.Regexp
		err            error
	)

	filter := make(map[string][]string)

	if v, ok := d.GetOk("image_id"); ok {
		imageId = v.(string)
		if imageId != "" {
			filter["image-id"] = []string{imageId}
		}
	}

	if v, ok := d.GetOk("image_type"); ok {
		for _, vv := range v.([]interface{}) {
			if vv, ok := vv.(string); ok && vv != "" {
				imageType = append(imageType, vv)
			}
		}
		if len(imageType) > 0 {
			filter["image-type"] = imageType
		}
	}

	if v, ok := d.GetOk("image_name_regex"); ok {
		imageName = v.(string)
		if imageName != "" {
			imageNameRegex, err = regexp.Compile(imageName)
			if err != nil {
				return fmt.Errorf("image_name_regex format error,%s", err.Error())
			}
		}
	}

	if v, ok := d.GetOk("os_name"); ok {
		osName = v.(string)
	}

	var instanceType string
	if v, ok := d.GetOk("instance_type"); ok {
		instanceType = v.(string)
	}

	var images []*cvm.Image
	err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		var e error
		images, e = cvmService.DescribeImagesByFilter(ctx, filter, instanceType)
		if e != nil {
			return tccommon.RetryError(e, tccommon.InternalError)
		}
		return nil
	})
	if err != nil {
		return err
	}

	var results []*cvm.Image
	images = sortImages(images)

	if osName == "" && imageName == "" {
		results = images
	} else {
		for _, image := range images {
			if osName != "" {
				if strings.Contains(strings.ToLower(*image.OsName), strings.ToLower(osName)) {
					results = append(results, image)
					continue
				}
			}
			if imageNameRegex != nil {
				if imageNameRegex.MatchString(*image.ImageName) {
					results = append(results, image)
					continue
				}
			}
		}
	}

	imageList := make([]map[string]interface{}, 0, len(results))
	ids := make([]string, 0, len(results))
	for _, image := range results {
		snapshots, err := imagesReadSnapshotByIds(ctx, cbsService, image)
		if err != nil {
			return err
		}

		mapping := map[string]interface{}{
			"image_id":           image.ImageId,
			"os_name":            image.OsName,
			"image_type":         image.ImageType,
			"created_time":       image.CreatedTime,
			"image_name":         image.ImageName,
			"image_description":  image.ImageDescription,
			"image_size":         image.ImageSize,
			"architecture":       image.Architecture,
			"image_state":        image.ImageState,
			"platform":           image.Platform,
			"image_creator":      image.ImageCreator,
			"image_source":       image.ImageSource,
			"sync_percent":       image.SyncPercent,
			"support_cloud_init": image.IsSupportCloudinit,
			"snapshots":          snapshots,
		}
		imageList = append(imageList, mapping)
		ids = append(ids, *image.ImageId)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	err = d.Set("images", imageList)
	if err != nil {
		log.Printf("[CRITAL]%s provider set image list fail, reason:%s\n ", logId, err.Error())
		return err
	}

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if err := tccommon.WriteToFile(output.(string), imageList); err != nil {
			return err
		}
	}

	return nil
}

func imagesReadSnapshotByIds(ctx context.Context, cbsService svccbs.CbsService, image *cvm.Image) (snapshotResults []map[string]interface{}, errRet error) {
	if len(image.SnapshotSet) == 0 {
		return
	}

	snapshotByIds := make([]*string, 0, len(image.SnapshotSet))
	for _, snapshot := range image.SnapshotSet {
		snapshotByIds = append(snapshotByIds, snapshot.SnapshotId)
	}

	snapshots, errRet := cbsService.DescribeSnapshotByIds(ctx, snapshotByIds)
	if errRet != nil {
		return
	}

	snapshotResults = make([]map[string]interface{}, 0, len(snapshots))
	for _, snapshot := range snapshots {
		snapshotMap := make(map[string]interface{}, 4)
		snapshotMap["snapshot_id"] = snapshot.SnapshotId
		snapshotMap["disk_usage"] = snapshot.DiskUsage
		snapshotMap["disk_size"] = snapshot.DiskSize
		snapshotMap["snapshot_name"] = snapshot.SnapshotName

		snapshotResults = append(snapshotResults, snapshotMap)
	}

	return
}
