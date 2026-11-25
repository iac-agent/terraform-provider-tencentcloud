package alb

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	albv20251030 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/alb/v20251030"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

// NewAlbService returns the ALB service bound to the shared provider client.
func NewAlbService(client *connectivity.TencentCloudClient) AlbService {
	return AlbService{client: client}
}

// AlbService wraps every ALB API call with the repo-standard retry, rate
// limiting and logging behaviour.
type AlbService struct {
	client *connectivity.TencentCloudClient
}

// CreateListener creates a listener and returns the id assigned by the cloud
// side. CreateListener is synchronous, so no task id is involved.
func (me *AlbService) CreateListener(ctx context.Context, request *albv20251030.CreateListenerRequest) (listenerId string, errRet error) {
	logId := tccommon.GetLogId(ctx)

	var response *albv20251030.CreateListenerResponse
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseAlbV20251030Client().CreateListener(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("create listener failed, Response is nil."))
		}

		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create listener failed, reason:%+v", logId, err)
		return
	}

	if response.Response.ListenerId == nil {
		errRet = fmt.Errorf("api %s returns empty listener id", request.GetAction())
		return
	}

	listenerId = *response.Response.ListenerId
	return
}

// DescribeListeners returns every listener matching the request filter. The
// API pages with MaxResults/NextToken, so the paging loop lives here and the
// caller gets the complete list (a unique-id filter terminates after one
// round, a full query keeps paging).
func (me *AlbService) DescribeListeners(ctx context.Context, request *albv20251030.DescribeListenersRequest) (listeners []*albv20251030.ListenerOutput, errRet error) {
	var (
		logId = tccommon.GetLogId(ctx)
		limit = uint64(100) // MaxResults range documented by the API: 1~100.
	)

	listeners = make([]*albv20251030.ListenerOutput, 0)

	for {
		request.MaxResults = &limit

		var (
			response *albv20251030.DescribeListenersResponse
			errRet   error
		)

		errRet = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(request.GetAction())
			result, e := me.client.UseAlbV20251030Client().DescribeListeners(request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			if result == nil || result.Response == nil {
				return resource.NonRetryableError(fmt.Errorf("describe listeners failed, Response is nil."))
			}

			response = result
			return nil
		})

		if errRet != nil {
			log.Printf("[CRITAL]%s describe listeners failed, reason:%+v", logId, errRet)
			return nil, errRet
		}

		listeners = append(listeners, response.Response.Listeners...)

		if len(response.Response.Listeners) < 1 {
			break
		}

		if response.Response.NextToken == nil || *response.Response.NextToken == "" {
			break
		}

		request.NextToken = response.Response.NextToken
	}

	return
}

// ModifyListenerAttributes updates the listener attributes accepted by the
// modify API. The call is synchronous.
func (me *AlbService) ModifyListenerAttributes(ctx context.Context, request *albv20251030.ModifyListenerAttributesRequest) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseAlbV20251030Client().ModifyListenerAttributes(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("modify listener attributes failed, Response is nil."))
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s modify listener attributes failed, reason:%+v", logId, err)
		errRet = err
	}

	return
}

// DeleteListener deletes the listener. The call is synchronous.
func (me *AlbService) DeleteListener(ctx context.Context, request *albv20251030.DeleteListenerRequest) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseAlbV20251030Client().DeleteListener(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("delete listener failed, Response is nil."))
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s delete listener failed, reason:%+v", logId, err)
		errRet = err
	}

	return
}
