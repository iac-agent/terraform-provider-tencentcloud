package alb

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	albv20251030 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/alb/v20251030"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

const (
	// ALB_LISTENER_DESCRIBE_LIMIT is the page size used by DescribeListeners.
	// The API accepts 1~100 entries per page.
	ALB_LISTENER_DESCRIBE_LIMIT = 100

	// Listener status values documented by DescribeListeners.
	ALB_LISTENER_STATUS_ACTIVE       = "Active"
	ALB_LISTENER_STATUS_PROVISIONING = "Provisioning"
	ALB_LISTENER_STATUS_CONFIGURING  = "Configuring"
	ALB_LISTENER_STATUS_FAILED       = "ProvisionFailed"
)

func NewAlbService(client *connectivity.TencentCloudClient) AlbService {
	return AlbService{client: client}
}

type AlbService struct {
	client *connectivity.TencentCloudClient
}

// CreateListener creates an ALB listener and returns the id assigned by the
// cloud side. RequestId is only a request trace id and must never be used as
// an identity.
func (me *AlbService) CreateListener(ctx context.Context, request *albv20251030.CreateListenerRequest) (listenerId string, errRet error) {
	logId := tccommon.GetLogId(ctx)
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	var response *albv20251030.CreateListenerResponse
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, err := me.client.UseAlbV20251030Client().CreateListenerWithContext(ctx, request)
		if err != nil {
			return tccommon.RetryError(err)
		}
		response = result
		return nil
	})
	if err != nil {
		errRet = err
		return
	}
	if response == nil || response.Response == nil {
		errRet = fmt.Errorf("CreateListener returned an empty response")
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response.Response.ListenerId == nil || *response.Response.ListenerId == "" {
		errRet = fmt.Errorf("CreateListener returned an empty listener id")
		return
	}

	listenerId = *response.Response.ListenerId
	return
}

// DescribeListeners returns the full listener list matching the request. The
// caller builds the request (load balancer id and/or listener ids) and picks
// its own item from the returned slice, so resource reads and data source
// queries can share this one method.
//
// The API is paginated through NextToken, so the whole result set is fetched
// here without caller-side branching.
func (me *AlbService) DescribeListeners(ctx context.Context, request *albv20251030.DescribeListenersRequest) (listeners []*albv20251030.ListenerOutput, totalCount uint64, errRet error) {
	logId := tccommon.GetLogId(ctx)
	listeners = make([]*albv20251030.ListenerOutput, 0)

	request.MaxResults = helper.Uint64(ALB_LISTENER_DESCRIBE_LIMIT)
	request.NextToken = nil

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for {
		var response *albv20251030.DescribeListenersResponse
		err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(request.GetAction())
			result, err := me.client.UseAlbV20251030Client().DescribeListenersWithContext(ctx, request)
			if err != nil {
				return tccommon.RetryError(err)
			}
			if result == nil || result.Response == nil {
				return resource.NonRetryableError(fmt.Errorf("DescribeListeners returned an empty response"))
			}

			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())

			response = result
			return nil
		})
		if err != nil {
			errRet = err
			return
		}

		if response.Response.TotalCount != nil {
			totalCount = *response.Response.TotalCount
		}
		listeners = append(listeners, response.Response.Listeners...)

		if response.Response.NextToken == nil || *response.Response.NextToken == "" {
			break
		}
		request.NextToken = response.Response.NextToken
	}

	return
}

// DescribeListenerDetailById returns the full detail of one listener.
//
// DescribeListeners does not echo back several configured attributes
// (DefaultActions, CaCertificateIds and CertificateIds); this call is the only
// read path that reports them. A nil detail with a nil error means the listener
// does not exist (any more).
func (me *AlbService) DescribeListenerDetailById(ctx context.Context, loadBalancerId, listenerId string) (detail *albv20251030.DescribeListenerDetailResponseParams, errRet error) {
	request := albv20251030.NewDescribeListenerDetailRequest()
	request.LoadBalancerId = helper.String(loadBalancerId)
	request.ListenerId = helper.String(listenerId)

	detail, err := me.DescribeListenerDetail(ctx, request)
	if err != nil {
		if tccommon.IsExpectError(err, []string{"ResourceNotFound.Listener", "ResourceNotFound"}) {
			return nil, nil
		}
		errRet = err
	}

	return
}

// DescribeListenerDetail returns the full detail of a single listener.
//
// DescribeListeners does not echo back several configured attributes
// (DefaultActions, CaCertificateIds and CertificateIds); this call is the only
// read path that reports them, so callers use it to refresh those attributes
// instead of trusting the value they wrote.
func (me *AlbService) DescribeListenerDetail(ctx context.Context, request *albv20251030.DescribeListenerDetailRequest) (detail *albv20251030.DescribeListenerDetailResponseParams, errRet error) {
	logId := tccommon.GetLogId(ctx)
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, err := me.client.UseAlbV20251030Client().DescribeListenerDetailWithContext(ctx, request)
		if err != nil {
			return tccommon.RetryError(err)
		}
		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("DescribeListenerDetail returned an empty response"))
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())

		detail = result.Response
		return nil
	})
	if err != nil {
		errRet = err
	}

	return
}

// ModifyListenerAttributes modifies the mutable attributes of a listener.
func (me *AlbService) ModifyListenerAttributes(ctx context.Context, request *albv20251030.ModifyListenerAttributesRequest) (errRet error) {
	logId := tccommon.GetLogId(ctx)
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, err := me.client.UseAlbV20251030Client().ModifyListenerAttributesWithContext(ctx, request)
		if err != nil {
			return tccommon.RetryError(err)
		}
		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("ModifyListenerAttributes returned an empty response"))
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		return nil
	})
	if err != nil {
		errRet = err
	}

	return
}

// DeleteListener deletes a listener.
func (me *AlbService) DeleteListener(ctx context.Context, request *albv20251030.DeleteListenerRequest) (errRet error) {
	logId := tccommon.GetLogId(ctx)
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, err := me.client.UseAlbV20251030Client().DeleteListenerWithContext(ctx, request)
		if err != nil {
			return tccommon.RetryError(err)
		}
		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("DeleteListener returned an empty response"))
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		return nil
	})
	if err != nil {
		errRet = err
	}

	return
}

// WaitListenerSettled polls the listener until the cloud reports a settled
// status (anything but "Provisioning" / "Configuring"), which is required
// before the state can be refreshed: writing back an attribute snapshot taken
// while a change is still in flight makes Terraform report "provider produced
// inconsistent result after apply".
//
// A listener that is not readable yet is treated as "still provisioning" so
// the wait also covers the eventual consistency window right after creation.
func (me *AlbService) WaitListenerSettled(ctx context.Context, loadBalancerId, listenerId string, timeout time.Duration) (errRet error) {
	logId := tccommon.GetLogId(ctx)
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s wait listener[%s] of load balancer[%s] settled failed, reason[%s]\n",
				logId, listenerId, loadBalancerId, errRet.Error())
		}
	}()

	err := resource.RetryContext(ctx, timeout, func() *resource.RetryError {
		request := albv20251030.NewDescribeListenersRequest()
		request.LoadBalancerId = helper.String(loadBalancerId)
		request.ListenerIds = []*string{helper.String(listenerId)}

		listeners, _, err := me.DescribeListeners(ctx, request)
		if err != nil {
			return resource.NonRetryableError(err)
		}

		var item *albv20251030.ListenerOutput
		for _, listener := range listeners {
			if listener != nil && listener.ListenerId != nil && *listener.ListenerId == listenerId {
				item = listener
				break
			}
		}
		if item == nil {
			return resource.RetryableError(fmt.Errorf("listener %s is not readable yet", listenerId))
		}

		status := ""
		if item.ListenerStatus != nil {
			status = *item.ListenerStatus
		}

		switch status {
		case ALB_LISTENER_STATUS_FAILED:
			return resource.NonRetryableError(fmt.Errorf("listener %s is in status %q", listenerId, status))
		case ALB_LISTENER_STATUS_ACTIVE:
			return nil
		default:
			return resource.RetryableError(fmt.Errorf("listener %s is in status %q, waiting for %q", listenerId, status, ALB_LISTENER_STATUS_ACTIVE))
		}
	})

	if err != nil {
		errRet = err
	}

	return
}
