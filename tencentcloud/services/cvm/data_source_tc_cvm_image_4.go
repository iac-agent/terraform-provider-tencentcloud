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
				Description: "List of image IDs to query. Mutually exclusive with `filters`.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"filters": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Filter conditions for the query. Mutually exclusive with `image_ids`.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Filter field name.",
						},
						"values": {
							Type:        schema.TypeSet,
							Required:    true,
							Description: "Filter field values.",
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
				Description: "Instance type for compatibility check, e.g., `SA5.MEDIUM2`.",
			},
			"image_set": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of images.",
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
							Description: "OS name of the image.",
						},
						"image_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Image type, e.g., `PUBLIC_IMAGE`, `PRIVATE_IMAGE`, `SHARED_IMAGE`.",
						},
						"created_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Image creation time.",
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
							Description: "Architecture, e.g., `x86_64`, `arm`, `i386`.",
						},
						"image_state": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Image state, e.g., `CREATING`, `NORMAL`, `CREATEFAILED`, `SYNCING`, `IMPORTING`, `IMPORTFAILED`.",
						},
						"platform": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Source platform.",
						},
						"image_creator": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Image creator.",
						},
						"image_source": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Image source, e.g., `OFFICIAL`, `CREATE_IMAGE`, `EXTERNAL_IMPORT`.",
						},
						"sync_percent": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Sync percentage.",
						},
						"is_support_cloudinit": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether cloud-init is supported.",
						},
						"snapshot_set": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Snapshot list of the image.",
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
										Description: "Disk usage.",
									},
									"disk_size": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Disk size in GiB.",
									},
								},
							},
						},
						"tags": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Tag list of the image.",
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
							Description: "License type, e.g., `TencentCloud`, `BYOL`.",
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

	var (
		logId   = tccommon.GetLogId(nil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = CvmService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("image_ids"); ok {
		imageIdsSet := v.([]interface{})
		tmpSet := make([]*string, 0, len(imageIdsSet))
		for _, item := range imageIdsSet {
			tmpSet = append(tmpSet, helper.String(item.(string)))
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
			if v, ok := filtersMap["values"]; ok {
				valueSet := v.(*schema.Set).List()
				for i := range valueSet {
					filter.Values = append(filter.Values, helper.String(valueSet[i].(string)))
				}
			}
			tmpSet = append(tmpSet, &filter)
		}
		paramMap["Filters"] = tmpSet
	}

	if v, ok := d.GetOk("instance_type"); ok {
		paramMap["InstanceType"] = v.(string)
	}

	var respData []*cvm.Image
	reqErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeCvmImage4ByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		respData = result
		return nil
	})

	if reqErr != nil {
		return reqErr
	}

	imageSetList := make([]map[string]interface{}, 0, len(respData))
	if respData != nil {
		for _, image := range respData {
			imageMap := map[string]interface{}{}
			if image.ImageId != nil {
				imageMap["image_id"] = image.ImageId
			}
			if image.OsName != nil {
				imageMap["os_name"] = image.OsName
			}
			if image.ImageType != nil {
				imageMap["image_type"] = image.ImageType
			}
			if image.CreatedTime != nil {
				imageMap["created_time"] = image.CreatedTime
			}
			if image.ImageName != nil {
				imageMap["image_name"] = image.ImageName
			}
			if image.ImageDescription != nil {
				imageMap["image_description"] = image.ImageDescription
			}
			if image.ImageSize != nil {
				imageMap["image_size"] = image.ImageSize
			}
			if image.Architecture != nil {
				imageMap["architecture"] = image.Architecture
			}
			if image.ImageState != nil {
				imageMap["image_state"] = image.ImageState
			}
			if image.Platform != nil {
				imageMap["platform"] = image.Platform
			}
			if image.ImageCreator != nil {
				imageMap["image_creator"] = image.ImageCreator
			}
			if image.ImageSource != nil {
				imageMap["image_source"] = image.ImageSource
			}
			if image.SyncPercent != nil {
				imageMap["sync_percent"] = image.SyncPercent
			}
			if image.IsSupportCloudinit != nil {
				imageMap["is_support_cloudinit"] = image.IsSupportCloudinit
			}
			if image.SnapshotSet != nil {
				snapshotList := make([]map[string]interface{}, 0, len(image.SnapshotSet))
				for _, snapshot := range image.SnapshotSet {
					snapshotMap := map[string]interface{}{}
					if snapshot.SnapshotId != nil {
						snapshotMap["snapshot_id"] = snapshot.SnapshotId
					}
					if snapshot.DiskUsage != nil {
						snapshotMap["disk_usage"] = snapshot.DiskUsage
					}
					if snapshot.DiskSize != nil {
						snapshotMap["disk_size"] = snapshot.DiskSize
					}
					snapshotList = append(snapshotList, snapshotMap)
				}
				imageMap["snapshot_set"] = snapshotList
			}
			if image.Tags != nil {
				tagList := make([]map[string]interface{}, 0, len(image.Tags))
				for _, tag := range image.Tags {
					tagMap := map[string]interface{}{}
					if tag.Key != nil {
						tagMap["key"] = tag.Key
					}
					if tag.Value != nil {
						tagMap["value"] = tag.Value
					}
					tagList = append(tagList, tagMap)
				}
				imageMap["tags"] = tagList
			}
			if image.LicenseType != nil {
				imageMap["license_type"] = image.LicenseType
			}
			if image.ImageFamily != nil {
				imageMap["image_family"] = image.ImageFamily
			}
			if image.ImageDeprecated != nil {
				imageMap["image_deprecated"] = image.ImageDeprecated
			}
			if image.CdcCacheStatus != nil {
				imageMap["cdc_cache_status"] = image.CdcCacheStatus
			}
			imageSetList = append(imageSetList, imageMap)
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
