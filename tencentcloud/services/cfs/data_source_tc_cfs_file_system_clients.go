package cfs

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cfs "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cfs/v20190719"
)

func DataSourceTencentCloudCfsFileSystemClients() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCfsFileSystemClientsRead,
		Schema: map[string]*schema.Schema{
			"file_system_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "File 系统 ID。",
			},

			"client_list": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Client 列表。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cfs_vip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "IP 地址 的 文件 系统。",
						},
						"client_ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "客户端 IP",
						},
						"vpc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "File 系统 VPCID。",
						},
						"zone": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "名称 availability 可用区，e.g. ap-beijing-1. For more 信息，see regions 和 availability zones 在 Overview document。",
						},
						"zone_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "AZ 名称",
						},
						"mount_directory": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "路径 在 其中 文件 系统 是 mounted 到 客户端。",
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

func dataSourceTencentCloudCfsFileSystemClientsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_cfs_file_system_clients.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	fsId := d.Get("file_system_id").(string)

	service := CfsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var clientList []*cfs.FileSystemClient

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeCfsFileSystemClientsById(ctx, fsId)
		if e != nil {
			return tccommon.RetryError(e)
		}
		clientList = result
		return nil
	})
	if err != nil {
		return err
	}

	tmpList := make([]map[string]interface{}, 0, len(clientList))

	if clientList != nil {
		for _, fileSystemClient := range clientList {
			fileSystemClientMap := map[string]interface{}{}

			if fileSystemClient.CfsVip != nil {
				fileSystemClientMap["cfs_vip"] = fileSystemClient.CfsVip
			}

			if fileSystemClient.ClientIp != nil {
				fileSystemClientMap["client_ip"] = fileSystemClient.ClientIp
			}

			if fileSystemClient.VpcId != nil {
				fileSystemClientMap["vpc_id"] = fileSystemClient.VpcId
			}

			if fileSystemClient.Zone != nil {
				fileSystemClientMap["zone"] = fileSystemClient.Zone
			}

			if fileSystemClient.ZoneName != nil {
				fileSystemClientMap["zone_name"] = fileSystemClient.ZoneName
			}

			if fileSystemClient.MountDirectory != nil {
				fileSystemClientMap["mount_directory"] = fileSystemClient.MountDirectory
			}

			tmpList = append(tmpList, fileSystemClientMap)
		}

		_ = d.Set("client_list", tmpList)
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
