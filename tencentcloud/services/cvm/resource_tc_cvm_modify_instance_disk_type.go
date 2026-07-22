package cvm

import (
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudCvmModifyInstanceDiskType() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCvmModifyInstanceDiskTypeCreate,
		Read:   resourceTencentCloudCvmModifyInstanceDiskTypeRead,
		Delete: resourceTencentCloudCvmModifyInstanceDiskTypeDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "实例 ID To obtain 实例 IDs，您 可以 call DescribeInstances 和 look 对于 实例 ID 在 response。",
			},

			"data_disks": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				Description: "For 实例 数据 磁盘 配置 信息，您 仅 need 到 指定media 类型 目标 云 磁盘 到 是 converted，和 指定value 的 DiskType. Currently，仅 一个 数据 磁盘 conversion 是 支持. CdcId 参数 是 仅 支持 对于 实例 的 CDHPAID 类型",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"disk_size": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Data 磁盘 大小 (在 GB). 最小 adjustment increment 是 10 GB. 值 范围 varies 通过 数据 磁盘 类型 默认值为 0，indicating 该 无 数据 磁盘 是 purchased. For more 信息，see product documentation。",
						},
						"disk_type": {
							Type:     schema.TypeString,
							Optional: true,
							Description: "数据盘类型。有效值:\n" +
								"- LOCAL_BASIC: local hard disk;\n" +
								"- LOCAL_SSD: local SSD hard disk;\n" +
								"- LOCAL_NVME: local NVME hard disk, which is strongly related to InstanceType and cannot be specified;\n" +
								"- LOCAL_PRO: local HDD hard disk, which is strongly related to InstanceType and cannot be specified;\n" +
								"- CLOUD_BASIC: ordinary cloud disk;\n" +
								"- CLOUD_PREMIUM: high-performance cloud disk;\n" +
								"- CLOUD_SSD:SSD cloud disk;\n" +
								"- CLOUD_HSSD: enhanced SSD cloud disk;\n" +
								"- CLOUD_TSSD: extremely fast SSD cloud disk;\n" +
								"- CLOUD_BSSD: general-purpose SSD cloud disk;\n" +
								"Default value: LOCAL_BASIC.",
						},
						"disk_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Data 磁盘 ID. 注意 该 它's 不 可用 对于 LOCAL_BASIC 和 LOCAL_SSD disks。",
						},
						"delete_with_instance": {
							Type:     schema.TypeBool,
							Optional: true,
							Description: "CVM 终止时是否终止数据盘。有效值:\n" +
								"- TRUE: terminate the data disk when its CVM is terminated. This value only supports pay-as-you-go cloud disks billed on an hourly basis.\n" +
								"- FALSE: retain the data disk when its CVM is terminated.\n" +
								"Default value: TRUE.",
						},
						"snapshot_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Data 磁盘 快照 ID. 大小 的 selected 数据 磁盘 快照 必须 是 smaller 比 该 的 数据 磁盘。",
						},
						"encrypt": {
							Type:     schema.TypeBool,
							Optional: true,
							Description: "数据盘是否加密。有效值:\n" +
								"- TRUE: encrypted\n" +
								"- FALSE: not encrypted\n" +
								"Default value: FALSE.",
						},
						"kms_key_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "ID 自定义 CMK 在 格式 的 UUID 或 “kms-abcd1234”. 此 参数 是 用于encrypt 云 disks。",
						},
						"throughput_performance": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Cloud 磁盘 performance，在 MB/s。",
						},
						"cdc_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "ID dedicated 集群 到 其中 实例 belongs。",
						},
					},
				},
			},

			"system_disk": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "For 实例 系统 磁盘 配置 信息，您 仅 need 到 指定nature 类型 目标 云 磁盘 到 是 converted，和 指定value 的 DiskType. Only CDHPAID 类型 实例 是 支持 到 指定Cd。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"disk_type": {
							Type:     schema.TypeString,
							Optional: true,
							Description: "系统盘类型。有效值：" +
								"- LOCAL_BASIC: local disk\n" +
								"- LOCAL_SSD: local SSD disk\n" +
								"- CLOUD_BASIC: ordinary cloud disk\n" +
								"- CLOUD_SSD: SSD cloud disk\n" +
								"- CLOUD_PREMIUM: Premium cloud storage\n" +
								"- CLOUD_BSSD: Balanced SSD\n" +
								"The disk currently in stock will be used by default.",
						},
						"disk_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "System 磁盘 ID. System disks whose 类型 是 LOCAL_BASIC 或 LOCAL_SSD do 不 have ID 和 do 不 support 此 参数。",
						},
						"disk_size": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "System 磁盘 大小; 单位: GB; 默认值：50 GB。",
						},
						"cdc_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "ID dedicated 集群 到 其中 实例 belongs。",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudCvmModifyInstanceDiskTypeCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cvm_modify_instance_disk_type.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request    = cvm.NewModifyInstanceDiskTypeRequest()
		instanceId string
	)
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId := v.(string)
		request.InstanceId = helper.String(instanceId)
	}

	if v, ok := d.GetOk("data_disks"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			dataDisk := cvm.DataDisk{}
			if v, ok := dMap["disk_size"]; ok {
				dataDisk.DiskSize = helper.IntInt64(v.(int))
			}
			if v, ok := dMap["disk_type"]; ok {
				dataDisk.DiskType = helper.String(v.(string))
			}
			if v, ok := dMap["disk_id"]; ok {
				dataDisk.DiskId = helper.String(v.(string))
			}
			if v, ok := dMap["delete_with_instance"]; ok {
				dataDisk.DeleteWithInstance = helper.Bool(v.(bool))
			}
			if v, ok := dMap["snapshot_id"]; ok {
				dataDisk.SnapshotId = helper.String(v.(string))
			}
			if v, ok := dMap["encrypt"]; ok {
				dataDisk.Encrypt = helper.Bool(v.(bool))
			}
			if v, ok := dMap["kms_key_id"]; ok {
				dataDisk.KmsKeyId = helper.String(v.(string))
			}
			if v, ok := dMap["throughput_performance"]; ok {
				dataDisk.ThroughputPerformance = helper.IntInt64(v.(int))
			}
			if v, ok := dMap["cdc_id"]; ok {
				dataDisk.CdcId = helper.String(v.(string))
			}
			request.DataDisks = append(request.DataDisks, &dataDisk)
		}
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "system_disk"); ok {
		systemDisk := cvm.SystemDisk{}
		if v, ok := dMap["disk_type"]; ok {
			systemDisk.DiskType = helper.String(v.(string))
		}
		if v, ok := dMap["disk_id"]; ok {
			systemDisk.DiskId = helper.String(v.(string))
		}
		if v, ok := dMap["disk_size"]; ok {
			systemDisk.DiskSize = helper.IntInt64(v.(int))
		}
		if v, ok := dMap["cdc_id"]; ok {
			systemDisk.CdcId = helper.String(v.(string))
		}
		request.SystemDisk = &systemDisk
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCvmClient().ModifyInstanceDiskType(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate cvm modifyInstanceDiskType failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(instanceId)

	return resourceTencentCloudCvmModifyInstanceDiskTypeRead(d, meta)
}

func resourceTencentCloudCvmModifyInstanceDiskTypeRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cvm_modify_instance_disk_type.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudCvmModifyInstanceDiskTypeDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cvm_modify_instance_disk_type.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
