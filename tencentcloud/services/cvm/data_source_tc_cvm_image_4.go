package cvm

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudCvmImage4() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCvmImage4Read,
		Schema: map[string]*schema.Schema{
			"image_ids": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Image ID list.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"filters": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Filter conditions.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Filter field name.",
						},
						"values": {
							Type:        schema.TypeList,
							Required:    true,
							Description: "Filter field value.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},

			"instance_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Instance type, such as `SA5.MEDIUM2`.",
			},

			"image_set": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A structure about image details, including the main status and attributes of the image.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"image_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Image ID.",
						},
						"os_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Image OS name.",
						},
						"image_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Image type. Valid values include: `PUBLIC_IMAGE`, `PRIVATE_IMAGE`, `SHARED_IMAGE`.",
						},
						"created_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Image creation time. Format: YYYY-MM-DDThh:mm:ssZ (ISO8601 standard, UTC time).",
						},
						"image_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Image name.",
						},
						"image_description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Image description.",
						},
						"image_size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Image size in GiB.",
						},
						"architecture": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Image architecture. Valid values include: `x86_64`, `arm`, `i386`.",
						},
						"image_state": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Image state. Valid values: CREATING, NORMAL, CREATEFAILED, SYNCING, IMPORTING, IMPORTFAILED.",
						},
						"platform": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Image source platform, including TencentOS, CentOS, Windows, Ubuntu, Debian, Fedora, etc.",
						},
						"image_creator": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Image creator.",
						},
						"image_source": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Image source. Valid values include: `OFFICIAL`, `CREATE_IMAGE`, `EXTERNAL_IMPORT`.",
						},
						"sync_percent": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Sync percentage. Note: This field may return null, indicating that no valid value can be obtained.",
						},
						"is_support_cloudinit": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether the image supports cloud-init.",
						},
						"snapshot_set": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Snapshot information associated with the image.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"snapshot_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Snapshot ID.",
									},
									"disk_usage": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Cloud disk type for this snapshot. Valid values: SYSTEM_DISK, DATA_DISK.",
									},
									"disk_size": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Cloud disk size for this snapshot in GiB.",
									},
								},
							},
						},
						"tags": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "List of tags associated with the image.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Tag key.",
									},
									"value": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Tag value.",
									},
								},
							},
						},
						"license_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Image license type. Valid values include: `TencentCloud`, `BYOL`.",
						},
						"image_family": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Image family.",
						},
						"image_deprecated": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether the image is deprecated.",
						},
						"cdc_cache_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "CDC image cache status.",
						},
					},
				},
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudCvmImage4Read(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_cvm_image_4.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := CvmService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("image_ids"); ok {
		imageIdsSet := v.([]interface{})
		tmpSet := make([]*string, 0, len(imageIdsSet))
		for _, item := range imageIdsSet {
			if item != nil {
				tmpSet = append(tmpSet, helper.String(item.(string)))
			}
		}
		paramMap["ImageIds"] = tmpSet
	}

	if v, ok := d.GetOk("filters"); ok {
		filtersSet := v.([]interface{})
		tmpSet := make([]*cvm.Filter, 0, len(filtersSet))
		for _, item := range filtersSet {
			filtersMap := item.(map[string]interface{})
			filter := cvm.Filter{}
			if v, ok := filtersMap["name"].(string); ok && v != "" {
				filter.Name = helper.String(v)
			}
			if v, ok := filtersMap["values"].([]interface{}); ok {
				for _, item := range v {
					filter.Values = append(filter.Values, helper.String(item.(string)))
				}
			}
			tmpSet = append(tmpSet, &filter)
		}
		paramMap["Filters"] = tmpSet
	}

	if v, ok := d.GetOk("instance_type"); ok {
		paramMap["InstanceType"] = helper.String(v.(string))
	}

	var images []*cvm.Image
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeCvmImage4ByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		images = result
		return nil
	})
	if err != nil {
		return err
	}

	imageSetList := make([]map[string]interface{}, 0, len(images))
	if images != nil {
		for _, image := range images {
			imageSetMap := map[string]interface{}{}

			if image.ImageId != nil {
				imageSetMap["image_id"] = image.ImageId
			}

			if image.OsName != nil {
				imageSetMap["os_name"] = image.OsName
			}

			if image.ImageType != nil {
				imageSetMap["image_type"] = image.ImageType
			}

			if image.CreatedTime != nil {
				imageSetMap["created_time"] = image.CreatedTime
			}

			if image.ImageName != nil {
				imageSetMap["image_name"] = image.ImageName
			}

			if image.ImageDescription != nil {
				imageSetMap["image_description"] = image.ImageDescription
			}

			if image.ImageSize != nil {
				imageSetMap["image_size"] = image.ImageSize
			}

			if image.Architecture != nil {
				imageSetMap["architecture"] = image.Architecture
			}

			if image.ImageState != nil {
				imageSetMap["image_state"] = image.ImageState
			}

			if image.Platform != nil {
				imageSetMap["platform"] = image.Platform
			}

			if image.ImageCreator != nil {
				imageSetMap["image_creator"] = image.ImageCreator
			}

			if image.ImageSource != nil {
				imageSetMap["image_source"] = image.ImageSource
			}

			if image.SyncPercent != nil {
				imageSetMap["sync_percent"] = image.SyncPercent
			}

			if image.IsSupportCloudinit != nil {
				imageSetMap["is_support_cloudinit"] = image.IsSupportCloudinit
			}

			if image.SnapshotSet != nil {
				snapshotSetList := make([]map[string]interface{}, 0, len(image.SnapshotSet))
				for _, snapshot := range image.SnapshotSet {
					snapshotSetMap := map[string]interface{}{}
					if snapshot.SnapshotId != nil {
						snapshotSetMap["snapshot_id"] = snapshot.SnapshotId
					}
					if snapshot.DiskUsage != nil {
						snapshotSetMap["disk_usage"] = snapshot.DiskUsage
					}
					if snapshot.DiskSize != nil {
						snapshotSetMap["disk_size"] = snapshot.DiskSize
					}
					snapshotSetList = append(snapshotSetList, snapshotSetMap)
				}
				imageSetMap["snapshot_set"] = snapshotSetList
			}

			if image.Tags != nil {
				tagsList := make([]map[string]interface{}, 0, len(image.Tags))
				for _, tag := range image.Tags {
					tagMap := map[string]interface{}{}
					if tag.Key != nil {
						tagMap["key"] = tag.Key
					}
					if tag.Value != nil {
						tagMap["value"] = tag.Value
					}
					tagsList = append(tagsList, tagMap)
				}
				imageSetMap["tags"] = tagsList
			}

			if image.LicenseType != nil {
				imageSetMap["license_type"] = image.LicenseType
			}

			if image.ImageFamily != nil {
				imageSetMap["image_family"] = image.ImageFamily
			}

			if image.ImageDeprecated != nil {
				imageSetMap["image_deprecated"] = image.ImageDeprecated
			}

			if image.CdcCacheStatus != nil {
				imageSetMap["cdc_cache_status"] = image.CdcCacheStatus
			}

			imageSetList = append(imageSetList, imageSetMap)
		}

		_ = d.Set("image_set", imageSetList)
	}

	d.SetId(helper.BuildToken())
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), d); e != nil {
			return e
		}
	}

	return nil
}
