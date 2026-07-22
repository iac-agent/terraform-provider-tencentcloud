package cfs

import (
	"context"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cfs "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cfs/v20190719"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudCfsFileSystems() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCfsFileSystemsRead,

		Schema: map[string]*schema.Schema{
			"file_system_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A 指定 文件 系统 ID 用于query。",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A 文件 系统 名称 用于query。",
			},
			"availability_zone": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "可用 可用区 该 文件 系统 locates 在。",
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID vpc 到 是 queried。",
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID vpc 子网。",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "用于保存结果。",
			},

			"file_system_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "An 信息 列表 云 文件 系统. Each element 包含following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"file_system_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID 文件 系统。",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "名称 文件 系统。",
						},
						"availability_zone": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "可用 可用区 该 文件 系统 locates 在。",
						},
						"protocol": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "协议 的 文件 系统。",
						},
						"access_group_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID 访问 组。",
						},
						"storage_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Storage 类型 文件 系统。",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "状态 文件 系统。",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "创建时间 的 文件 系统。",
						},
						"size_limit": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Size 限制 的 文件 系统。",
						},
						"size_used": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Size 使用 的 文件 系统。",
						},
						"mount_ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "IP 的 文件 系统。",
						},
						"fs_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Mount root-directory。",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudCfsFileSystemsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_cfs_file_systems.read")()
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	cfsService := CfsService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}

	var fileSystemId string
	var vpcId string
	var subnetId string
	var name string
	var zone string
	if v, ok := d.GetOk("file_system_id"); ok {
		fileSystemId = v.(string)
	}
	if v, ok := d.GetOk("vpc_id"); ok {
		vpcId = v.(string)
	}
	if v, ok := d.GetOk("subnet_id"); ok {
		subnetId = v.(string)
	}
	if v, ok := d.GetOk("name"); ok {
		name = v.(string)
	}
	if v, ok := d.GetOk("availability_zone"); ok {
		zone = v.(string)
	}

	var fileSystems []*cfs.FileSystemInfo
	var errRet error
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		fileSystems, errRet = cfsService.DescribeFileSystem(ctx, fileSystemId, vpcId, subnetId)
		if errRet != nil {
			return tccommon.RetryError(errRet, tccommon.InternalError)
		}
		return nil
	})
	if err != nil {
		return err
	}

	fileSystemList := make([]map[string]interface{}, 0, len(fileSystems))
	ids := make([]string, 0, len(fileSystems))
	for _, fileSystem := range fileSystems {
		if name != "" && name != *fileSystem.FsName {
			continue
		}
		if zone != "" && zone != *fileSystem.Zone {
			continue
		}
		mapping := map[string]interface{}{
			"file_system_id":    fileSystem.FileSystemId,
			"name":              fileSystem.FsName,
			"availability_zone": fileSystem.Zone,
			"protocol":          fileSystem.Protocol,
			"access_group_id":   fileSystem.PGroup.PGroupId,
			"storage_type":      fileSystem.StorageType,
			"status":            fileSystem.LifeCycleState,
			"create_time":       fileSystem.CreationTime,
			"size_limit":        fileSystem.SizeLimit,
			"size_used":         fileSystem.SizeByte,
		}
		targets, err := cfsService.DescribeMountTargets(ctx, *fileSystem.FileSystemId)
		if err != nil {
			return err
		}
		var mountTarget *cfs.MountInfo
		if len(targets) > 0 {
			mountTarget = targets[0]
		}
		if mountTarget != nil {
			mapping["mount_ip"] = mountTarget.IpAddress
			mapping["fs_id"] = mountTarget.FSID
		}
		fileSystemList = append(fileSystemList, mapping)
		ids = append(ids, *fileSystem.FileSystemId)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	err = d.Set("file_system_list", fileSystemList)
	if err != nil {
		log.Printf("[CRITAL]%s provider set cfs file system list fail, reason:%s\n ", logId, err.Error())
		return err
	}
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if err := tccommon.WriteToFile(output.(string), fileSystemList); err != nil {
			return err
		}
	}
	return nil
}
