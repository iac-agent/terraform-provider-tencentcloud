package cfs

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cfs "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cfs/v20190719"
)

func DataSourceTencentCloudCfsMountTargets() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCfsMountTargetsRead,
		Schema: map[string]*schema.Schema{
			"file_system_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "File 系统 ID。",
			},

			"mount_targets": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "挂载目标 details。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"file_system_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "File 系统 ID。",
						},
						"mount_target_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "挂载目标 ID。",
						},
						"ip_address": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "挂载目标 IP。",
						},
						"fs_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Mount root-directory。",
						},
						"life_cycle_state": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "挂载目标 状态",
						},
						"network_interface": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Network 类型",
						},
						"vpc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "私有网络 ID",
						},
						"vpc_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "VPC 名称",
						},
						"subnet_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "子网 ID",
						},
						"subnet_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Subnet 名称",
						},
						"ccn_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "CCN 实例 ID 使用 通过 CFS Turbo。",
						},
						"cidr_block": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "CCN IP 范围 使用 通过 CFS Turbo。",
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

func dataSourceTencentCloudCfsMountTargetsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_cfs_mount_targets.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := CfsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var mountTargets []*cfs.MountInfo

	fsId := d.Get("file_system_id").(string)
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeCfsMountTargetsById(ctx, fsId)
		if e != nil {
			return tccommon.RetryError(e)
		}
		mountTargets = result
		return nil
	})
	if err != nil {
		return err
	}

	tmpList := make([]map[string]interface{}, 0, len(mountTargets))

	if mountTargets != nil {
		for _, mountInfo := range mountTargets {
			mountInfoMap := map[string]interface{}{}

			if mountInfo.FileSystemId != nil {
				mountInfoMap["file_system_id"] = mountInfo.FileSystemId
			}

			if mountInfo.MountTargetId != nil {
				mountInfoMap["mount_target_id"] = mountInfo.MountTargetId
			}

			if mountInfo.IpAddress != nil {
				mountInfoMap["ip_address"] = mountInfo.IpAddress
			}

			if mountInfo.FSID != nil {
				mountInfoMap["fs_id"] = mountInfo.FSID
			}

			if mountInfo.LifeCycleState != nil {
				mountInfoMap["life_cycle_state"] = mountInfo.LifeCycleState
			}

			if mountInfo.NetworkInterface != nil {
				mountInfoMap["network_interface"] = mountInfo.NetworkInterface
			}

			if mountInfo.VpcId != nil {
				mountInfoMap["vpc_id"] = mountInfo.VpcId
			}

			if mountInfo.VpcName != nil {
				mountInfoMap["vpc_name"] = mountInfo.VpcName
			}

			if mountInfo.SubnetId != nil {
				mountInfoMap["subnet_id"] = mountInfo.SubnetId
			}

			if mountInfo.SubnetName != nil {
				mountInfoMap["subnet_name"] = mountInfo.SubnetName
			}

			if mountInfo.CcnID != nil {
				mountInfoMap["ccn_id"] = mountInfo.CcnID
			}

			if mountInfo.CidrBlock != nil {
				mountInfoMap["cidr_block"] = mountInfo.CidrBlock
			}

			tmpList = append(tmpList, mountInfoMap)
		}

		_ = d.Set("mount_targets", tmpList)
	}

	d.SetId(fsId)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
