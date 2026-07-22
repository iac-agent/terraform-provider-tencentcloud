package cynosdb

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cynosdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cynosdb/v20190107"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudCynosdbResourcePackageSaleSpecs() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCynosdbResourcePackageSaleSpecsRead,
		Schema: map[string]*schema.Schema{
			"instance_type": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "实例类型。值范围：cynosdb-serverless、cynosdb、cdb。",
			},
			"package_region": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "资源包使用地区 中国-中国大陆通用，海外-港澳台、海外通用。",
			},
			"package_type": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "资源包类型 CCU - 计算资源包 DISK - 存储资源包。",
			},
			"detail": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "资源包详情说明：该字段可能返回null，表示无法获取到有效值。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"package_region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "注意：该字段可能返回null，表示无法获取到有效值。",
						},
						"package_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "资源包类型 CCU - 计算资源包 DISK - 存储资源包 注意：该字段可能返回 null，表示无法获取有效值。",
						},
						"package_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "资源包版本base基础版、普通通用版、企业企业版 注：该字段可能返回null，表示无法获取有效值。",
						},
						"min_package_spec": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "当前版本资源包中的最小资源数量，以资源为单位计算；存储资源：GB 注：该字段可能返回null，表示无法获取有效值。",
						},
						"max_package_spec": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "当前版本资源包的最大资源数量，以资源为单位计算；存储资源：GB 注：该字段可能返回null，表示无法获取有效值。",
						},
						"expire_day": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "资源包有效期，单位为天。注意：该字段可能返回null，表示无法获取到有效值。",
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

func dataSourceTencentCloudCynosdbResourcePackageSaleSpecsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_cynosdb_resource_package_sale_specs.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId         = tccommon.GetLogId(tccommon.ContextNil)
		ctx           = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service       = CynosdbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		detail        []*cynosdb.SalePackageSpec
		instanceType  string
		packageRegion string
		packageType   string
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_type"); ok {
		paramMap["InstanceType"] = helper.String(v.(string))
		instanceType = v.(string)
	}

	if v, ok := d.GetOk("package_region"); ok {
		paramMap["PackageRegion"] = helper.String(v.(string))
		packageRegion = v.(string)
	}

	if v, ok := d.GetOk("package_type"); ok {
		paramMap["PackageType"] = helper.String(v.(string))
		packageType = v.(string)
	}

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeCynosdbResourcePackageSaleSpecsByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		detail = result
		return nil
	})

	if err != nil {
		return err
	}

	ids := make([]string, 0, len(detail))
	ids = append(ids, instanceType, packageRegion, packageType)
	tmpList := make([]map[string]interface{}, 0, len(detail))

	if detail != nil {
		for _, salePackageSpec := range detail {
			salePackageSpecMap := map[string]interface{}{}

			if salePackageSpec.PackageRegion != nil {
				salePackageSpecMap["package_region"] = salePackageSpec.PackageRegion
			}

			if salePackageSpec.PackageType != nil {
				salePackageSpecMap["package_type"] = salePackageSpec.PackageType
			}

			if salePackageSpec.PackageVersion != nil {
				salePackageSpecMap["package_version"] = salePackageSpec.PackageVersion
			}

			if salePackageSpec.MinPackageSpec != nil {
				salePackageSpecMap["min_package_spec"] = salePackageSpec.MinPackageSpec
			}

			if salePackageSpec.MaxPackageSpec != nil {
				salePackageSpecMap["max_package_spec"] = salePackageSpec.MaxPackageSpec
			}

			if salePackageSpec.ExpireDay != nil {
				salePackageSpecMap["expire_day"] = salePackageSpec.ExpireDay
			}

			tmpList = append(tmpList, salePackageSpecMap)
		}

		_ = d.Set("detail", tmpList)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), tmpList); e != nil {
			return e
		}
	}

	return nil
}
