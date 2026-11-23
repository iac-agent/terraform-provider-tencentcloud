package alb

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	alb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/alb/v20251030"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

func NewAlbService(client *connectivity.TencentCloudClient) AlbService {
	return AlbService{client: client}
}

type AlbService struct {
	client *connectivity.TencentCloudClient
}

func (me *AlbService) DescribeListenerById(ctx context.Context, loadBalancerId, listenerId string) (listener *alb.ListenerOutput, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := alb.NewDescribeListenersRequest()
	request.LoadBalancerId = helper.String(loadBalancerId)
	request.ListenerIds = helper.Strings([]string{listenerId})

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		response, e := me.client.UseAlbV20251030Client().DescribeListeners(request)
		if e != nil {
			return tccommon.RetryError(e)
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
		if len(response.Response.Listeners) > 0 {
			listener = response.Response.Listeners[0]
		}
		return nil
	})
	if err != nil {
		errRet = err
		return
	}
	return
}
