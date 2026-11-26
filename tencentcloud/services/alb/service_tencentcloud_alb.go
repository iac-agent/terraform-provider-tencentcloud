package alb

import (
	"context"
	"fmt"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	albv20251030 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/alb/v20251030"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

// NewAlbService builds the service used to call the ALB APIs.
func NewAlbService(client *connectivity.TencentCloudClient) AlbService {
	return AlbService{client: client}
}

// AlbService wraps ALB API calls with retry and logging.
type AlbService struct {
	client *connectivity.TencentCloudClient
}

// DescribeListeners queries the listener list, following NextToken until all
// pages have been fetched. The total count returned by the API is an internal
// paging detail and is not exposed to callers.
func (me *AlbService) DescribeListeners(ctx context.Context, request *albv20251030.DescribeListenersRequest) (listeners []*albv20251030.ListenerOutput, err error) {
	var (
		logId         = tccommon.GetLogId(ctx)
		limit         = uint64(100)
		nextToken     *string
		pageListeners []*albv20251030.ListenerOutput
	)

	request.MaxResults = &limit
	listeners = make([]*albv20251030.ListenerOutput, 0)

	for {
		errRet := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(request.GetAction())
			result, e := me.client.UseAlbV20251030Client().DescribeListeners(request)
			if e != nil {
				return tccommon.RetryError(e)
			}

			if result == nil || result.Response == nil {
				return resource.NonRetryableError(fmt.Errorf("describe listeners failed, Response is nil"))
			}

			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())

			nextToken = result.Response.NextToken
			pageListeners = result.Response.Listeners
			return nil
		})

		if errRet != nil {
			log.Printf("[CRITAL]%s describe listeners failed, reason:%+v", logId, errRet)
			return nil, errRet
		}

		listeners = append(listeners, pageListeners...)

		if nextToken == nil || *nextToken == "" {
			break
		}

		request.NextToken = nextToken
	}

	return
}
