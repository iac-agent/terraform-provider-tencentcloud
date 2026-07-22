package cvm

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudCvmChcNetworkMode() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCvmChcNetworkModeCreate,
		Read:   resourceTencentCloudCvmChcNetworkModeRead,
		Update: resourceTencentCloudCvmChcNetworkModeUpdate,
		Delete: resourceTencentCloudCvmChcNetworkModeDelete,
		Schema: map[string]*schema.Schema{
			"chc_ids": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "CHC physical server ID list.",
			},

			"network_mode": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Network mode to switch to. Valid values: DEPLOY (deploy network mode), BUSINESS (business network mode).",
			},
		},
	}
}

func resourceTencentCloudCvmChcNetworkModeCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cvm_chc_network_mode.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request = cvm.NewModifyChcNetworkModeRequest()
		chcIds  []string
	)

	if v, ok := d.GetOk("chc_ids"); ok {
		chcIds = helper.InterfacesStrings(v.([]interface{}))
		for _, chcId := range chcIds {
			request.ChcIds = append(request.ChcIds, &chcId)
		}
	}

	if v, ok := d.GetOk("network_mode"); ok {
		request.NetworkMode = helper.String(v.(string))
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCvmClient().ModifyChcNetworkMode(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create chc_network_mode failed, Response is nil."))
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create chc_network_mode failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	d.SetId(strings.Join(chcIds, tccommon.FILED_SP))

	return resourceTencentCloudCvmChcNetworkModeRead(d, meta)
}

func resourceTencentCloudCvmChcNetworkModeRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cvm_chc_network_mode.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := CvmService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	params := map[string]interface{}{
		"chc_ids": idSplit,
	}
	chcHosts, err := service.DescribeCvmChcHostsByFilter(ctx, params)
	if err != nil {
		return err
	}

	if len(chcHosts) < 1 {
		log.Printf("[CRUD] chc_network_mode id=%s", d.Id())
		d.SetId("")
		return nil
	}

	chcIds := make([]string, 0)
	for _, chcHost := range chcHosts {
		if chcHost.ChcId != nil {
			chcIds = append(chcIds, *chcHost.ChcId)
		}
	}
	_ = d.Set("chc_ids", chcIds)

	// network_mode is not returned by DescribeChcHosts API, keep from state
	// The ChcHost structure does not contain a NetworkMode field

	return nil
}

func resourceTencentCloudCvmChcNetworkModeUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cvm_chc_network_mode.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	if d.HasChange("network_mode") {
		request := cvm.NewModifyChcNetworkModeRequest()

		idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
		for _, chcId := range idSplit {
			request.ChcIds = append(request.ChcIds, &chcId)
		}

		if v, ok := d.GetOk("network_mode"); ok {
			request.NetworkMode = helper.String(v.(string))
		}

		reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCvmClient().ModifyChcNetworkMode(request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			if result == nil || result.Response == nil {
				return resource.NonRetryableError(fmt.Errorf("Update chc_network_mode failed, Response is nil."))
			}

			return nil
		})

		if reqErr != nil {
			log.Printf("[CRITAL]%s update chc_network_mode failed, reason:%+v", logId, reqErr)
			return reqErr
		}
	}

	return resourceTencentCloudCvmChcNetworkModeRead(d, meta)
}

func resourceTencentCloudCvmChcNetworkModeDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cvm_chc_network_mode.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	log.Printf("[DEBUG]%s resource`tencentcloud_cvm_chc_network_mode` delete, only remove from state, no api call\n", tccommon.GetLogId(tccommon.ContextNil))
	d.SetId("")

	return nil
}
