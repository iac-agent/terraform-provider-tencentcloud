package emr

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	emr "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/emr/v20190103"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

const (
	EmrInternetStatusCreated int64 = 2
	EmrInternetStatusDeleted int64 = 201
)

const (
	DisplayStrategyIsclusterList = "clusterList"
)

const (
	EMR_MASTER_WAN_TYPE_NEED_MASTER_WAN     = "NEED_MASTER_WAN"
	EMR_MASTER_WAN_TYPE_NOT_NEED_MASTER_WAN = "NOT_NEED_MASTER_WAN"
)

const (
	F_KEY_FLOW_ID  = "FlowId"
	F_KEY_TRACE_ID = "TraceId"
)

var EMR_MASTER_WAN_TYPES = []string{EMR_MASTER_WAN_TYPE_NEED_MASTER_WAN, EMR_MASTER_WAN_TYPE_NOT_NEED_MASTER_WAN}

func buildResourceSpecSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		ForceNew: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"spec": {
					Type:        schema.TypeString,
					Optional:    true,
					ForceNew:    true,
					Description: "Node 规格 描述, such 作为 CVM.SA2.",
				},
				"storage_type": {
					Type:     schema.TypeInt,
					Optional: true,
					ForceNew: true,
					Description: "Storage 类型. Value 范围:\n" +
						"	- 4: Represents cloud SSD;\n" +
						"	- 5: Represents efficient cloud disk;\n" +
						"	- 6: Represents enhanced SSD Cloud Block Storage;\n" +
						"	- 11: Represents throughput Cloud Block Storage;\n" +
						"	- 12: Represents extremely fast SSD Cloud Block Storage.",
				},
				"disk_type": {
					Type:     schema.TypeString,
					Optional: true,
					ForceNew: true,
					Description: "磁盘 types. Value 范围:\n" +
						"	- CLOUD_SSD: Represents cloud SSD;\n" +
						"	- CLOUD_PREMIUM: Represents efficient cloud disk;\n" +
						"	- CLOUD_BASIC: Represents Cloud Block Storage.",
				},
				"mem_size": {
					Type:        schema.TypeInt,
					Optional:    true,
					ForceNew:    true,
					Description: "Memory 大小 在 M.",
				},
				"cpu": {
					Type:        schema.TypeInt,
					Optional:    true,
					ForceNew:    true,
					Description: "Number 的 CPU cores.",
				},
				"disk_size": {
					Type:        schema.TypeInt,
					Optional:    true,
					ForceNew:    true,
					Description: "Data 磁盘 容量.",
				},
				"root_size": {
					Type:        schema.TypeInt,
					Optional:    true,
					ForceNew:    true,
					Description: "Root 磁盘 容量.",
				},
				"multi_disks": {
					Type:        schema.TypeSet,
					Optional:    true,
					Computed:    true,
					ForceNew:    true,
					Description: "Cloud 磁盘 列表. 当 数据 磁盘 是 云 磁盘, 使用 disk_type 和 disk_size 参数 directly, 和 使用 multi_disks 对于 excess parts.",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"disk_type": {
								Type:     schema.TypeString,
								Optional: true,
								ForceNew: true,
								Computed: true,
								Elem: &schema.Schema{
									Type: schema.TypeString,
								},
								Description: "Cloud 磁盘 类型\n" +
									"	- CLOUD_SSD: Represents cloud SSD;\n" +
									"	- CLOUD_PREMIUM: Represents efficient cloud disk;\n" +
									"	- CLOUD_HSSD: Represents enhanced SSD Cloud Block Storage.",
							},
							"volume": {
								Type:        schema.TypeInt,
								Optional:    true,
								ForceNew:    true,
								Computed:    true,
								Description: "Cloud 磁盘 大小.",
							},
							"count": {
								Type:        schema.TypeInt,
								Optional:    true,
								ForceNew:    true,
								Computed:    true,
								Description: "Number 的 云 disks 的 此 类型.",
							},
						},
					},
					Set: func(v interface{}) int {
						m := v.(map[string]interface{})
						return helper.HashString(fmt.Sprintf("%s-%d-%d", m["disk_type"].(string), m["volume"].(int), m["count"].(int)))

					},
				},
			},
		},
		Description: "Resource details.",
	}
}

func ParseMultiDisks(_multiDisks []interface{}) []*emr.MultiDisk {
	multiDisks := make([]*emr.MultiDisk, 0, len(_multiDisks))
	for _, multiDisk := range _multiDisks {
		item := multiDisk.(map[string]interface{})
		var diskType string
		var volume int
		var count int
		for subK, subV := range item {
			if subK == "disk_type" {
				diskType = subV.(string)
			} else if subK == "volume" {
				volume = subV.(int)
			} else if subK == "count" {
				count = subV.(int)
			}
		}
		multiDisks = append(multiDisks,
			&emr.MultiDisk{
				DiskType: helper.String(diskType),
				Volume:   helper.IntInt64(volume),
				Count:    helper.IntInt64(count),
			})
	}

	return multiDisks
}

func ParseTags(_tags []map[string]string) []*emr.Tag {
	tags := make([]*emr.Tag, len(_tags))
	for _, item := range _tags {
		for k, v := range item {
			tags = append(tags, &emr.Tag{TagKey: &k, TagValue: &v})
		}
	}
	return tags
}

func ParseResource(_resource map[string]interface{}) *emr.Resource {
	resultResource := &emr.Resource{}
	for k, v := range _resource {
		if k == "spec" {
			resultResource.Spec = common.StringPtr(v.(string))
		} else if k == "storage_type" {
			resultResource.StorageType = common.Int64Ptr((int64)(v.(int)))
		} else if k == "disk_type" {
			resultResource.DiskType = common.StringPtr(v.(string))
		} else if k == "mem_size" {
			resultResource.MemSize = common.Int64Ptr((int64)(v.(int)))
		} else if k == "cpu" {
			resultResource.Cpu = common.Int64Ptr((int64)(v.(int)))
		} else if k == "disk_size" {
			resultResource.DiskSize = common.Int64Ptr((int64)(v.(int)))
		} else if k == "root_size" {
			resultResource.RootSize = common.Int64Ptr((int64)(v.(int)))
		} else if k == "multi_disks" {
			multiDisks := v.(*schema.Set).List()
			resultResource.MultiDisks = ParseMultiDisks(multiDisks)
		} else if k == "tags" {
			tags := v.([]map[string]string)
			resultResource.Tags = ParseTags(tags)
		} else if k == "instance_type" {
			resultResource.InstanceType = common.StringPtr(v.(string))
		} else if k == "local_disk_num" {
			resultResource.LocalDiskNum = common.Uint64Ptr(v.(uint64))
		} else if k == "DiskNum" {
			resultResource.DiskNum = common.Uint64Ptr(v.(uint64))
		}
	}
	return resultResource
}

func validateMultiDisks(r map[string]interface{}) error {
	if _, ok := r["multi_disks"]; !ok {
		return nil
	}
	multiDiskList := r["multi_disks"].(*schema.Set).List()
	visited := make(map[string]struct{})

	for _, multiDisk := range multiDiskList {
		multiDiskMap := multiDisk.(map[string]interface{})
		key := fmt.Sprintf("%s-%d", multiDiskMap["disk_type"].(string), multiDiskMap["volume"].(int))
		if _, ok := visited[key]; ok {
			return fmt.Errorf("Merge disks of the same specifications")
		} else {
			visited[key] = struct{}{}
		}
	}

	return nil
}

func translateDiskType(diskType int64) (diskTypeStr string) {
	switch diskType {
	case 4:
		diskTypeStr = "CLOUD_SSD"
	case 5:
		diskTypeStr = "CLOUD_PREMIUM"
	case 6:
		diskTypeStr = "CLOUD_HSSD"
	}
	return
}
func fetchMultiDisks(v *emr.NodeHardwareInfo, r *emr.OutterResource) (multiDisks []interface{}) {
	var inputDataDiskTag string
	if r.DiskType != nil && r.DiskSize != nil {
		inputDataDiskTag = fmt.Sprintf("%s-%d", *r.DiskType, *r.DiskSize)
	}
	for _, item := range v.MCMultiDisk {
		outputDataDiskTag := ""
		multiDisk := make(map[string]interface{})
		if item.Type != nil {
			var diskType string
			diskType = translateDiskType(*item.Type)
			multiDisk["disk_type"] = diskType
			outputDataDiskTag = diskType
		}
		if item.Volume != nil {
			volume := int(*item.Volume / 1024 / 1024 / 1024)
			multiDisk["volume"] = volume
			outputDataDiskTag = fmt.Sprintf("%s-%d", outputDataDiskTag, volume)
		}
		var count int
		if item.Count != nil {
			count = int(*item.Count)
			if count > 0 && inputDataDiskTag == outputDataDiskTag {
				count -= 1
			}
			multiDisk["count"] = count
		}

		if count != 0 {
			multiDisks = append(multiDisks, multiDisk)
		}
	}
	return
}
