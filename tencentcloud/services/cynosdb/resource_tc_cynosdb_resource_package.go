package cynosdb

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cynosdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cynosdb/v20190107"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudCynosdbResourcePackage() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCynosdbResourcePackageCreate,
		Read:   resourceTencentCloudCynosdbResourcePackageRead,
		Update: resourceTencentCloudCynosdbResourcePackageUpdate,
		Delete: resourceTencentCloudCynosdbResourcePackageDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_type": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "实例类型。",
			},

			"package_region": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "资源包使用地区 中国-中国大陆通用，海外-港澳台、海外通用。",
			},

			"package_type": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "资源包类型：CCU计算资源包、DISK存储资源包。",
			},

			"package_version": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "资源包版本基础基础版、普通通用版、企业企业版。",
			},

			"package_spec": {
				Required:    true,
				Type:        schema.TypeFloat,
				Description: "资源包大小，以10000个单位计算；存储资源：GB。",
			},

			"expire_day": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "资源包有效期，单位为天。",
			},

			"package_count": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "购买的资源包数量。",
			},

			"package_name": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "资源包名称。",
			},
		},
	}
}

func resourceTencentCloudCynosdbResourcePackageCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cynosdb_resource_package.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request = cynosdb.NewCreateResourcePackageRequest()
		// response  = cynosdb.NewCreateResourcePackageResponse()
		// packageId string
	)
	if v, ok := d.GetOk("instance_type"); ok {
		request.InstanceType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("package_region"); ok {
		request.PackageRegion = helper.String(v.(string))
	}

	if v, ok := d.GetOk("package_type"); ok {
		request.PackageType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("package_version"); ok {
		request.PackageVersion = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("package_spec"); ok {
		request.PackageSpec = helper.Float64(v.(float64))
	}

	if v, ok := d.GetOkExists("expire_day"); ok {
		request.ExpireDay = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("package_count"); ok {
		request.PackageCount = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("package_name"); ok {
		request.PackageName = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCynosdbClient().CreateResourcePackage(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		// response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create cynosdb resourcePackage failed, reason:%+v", logId, err)
		return err
	}

	// packageId = *response.Response.PackageId
	// d.SetId(helper.String(packageId))

	return resourceTencentCloudCynosdbResourcePackageRead(d, meta)
}

func resourceTencentCloudCynosdbResourcePackageRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cynosdb_resource_package.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := CynosdbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	packageId := d.Id()
	resourcePackage, err := service.DescribeCynosdbResourcePackageById(ctx, packageId)
	if err != nil {
		return err
	}

	if resourcePackage == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `CynosdbResourcePackage` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	// if resourcePackage.InstanceType != nil {
	// 	_ = d.Set("instance_type", resourcePackage.InstanceType)
	// }

	// if resourcePackage.PackageRegion != nil {
	// 	_ = d.Set("package_region", resourcePackage.PackageRegion)
	// }

	// if resourcePackage.PackageType != nil {
	// 	_ = d.Set("package_type", resourcePackage.PackageType)
	// }

	// if resourcePackage.PackageVersion != nil {
	// 	_ = d.Set("package_version", resourcePackage.PackageVersion)
	// }

	// if resourcePackage.PackageSpec != nil {
	// 	_ = d.Set("package_spec", resourcePackage.PackageSpec)
	// }

	// if resourcePackage.ExpireDay != nil {
	// 	_ = d.Set("expire_day", resourcePackage.ExpireDay)
	// }

	// if resourcePackage.PackageCount != nil {
	// 	_ = d.Set("package_count", resourcePackage.PackageCount)
	// }

	// if resourcePackage.PackageName != nil {
	// 	_ = d.Set("package_name", resourcePackage.PackageName)
	// }

	return nil
}

func resourceTencentCloudCynosdbResourcePackageUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cynosdb_resource_package.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	request := cynosdb.NewModifyResourcePackageNameRequest()

	packageId := d.Id()

	request.PackageId = &packageId

	immutableArgs := []string{"instance_type", "package_region", "package_type", "package_version", "package_spec", "expire_day", "package_count", "package_name"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("package_name") {
		if v, ok := d.GetOk("package_name"); ok {
			request.PackageName = helper.String(v.(string))
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCynosdbClient().ModifyResourcePackageName(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update cynosdb resourcePackage failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudCynosdbResourcePackageRead(d, meta)
}

func resourceTencentCloudCynosdbResourcePackageDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cynosdb_resource_package.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := CynosdbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	packageId := d.Id()

	if err := service.DeleteCynosdbResourcePackageById(ctx, packageId); err != nil {
		return err
	}

	return nil
}
