package alb

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	albv20251030 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/alb/v20251030"

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

func (me *AlbService) CreateListener(ctx context.Context, request *albv20251030.CreateListenerRequest) (listenerId string, errRet error) {
	logId := tccommon.GetLogId(ctx)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseAlbV20251030Client().CreateListenerWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("create alb listener failed, Response is nil"))
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())

		if result.Response.ListenerId == nil || *result.Response.ListenerId == "" {
			return resource.NonRetryableError(fmt.Errorf("create alb listener failed, ListenerId is empty"))
		}

		listenerId = *result.Response.ListenerId
		return nil
	})

	if err != nil {
		errRet = err
		return
	}

	return
}

func (me *AlbService) DescribeListenerById(ctx context.Context, loadBalancerId, listenerId string) (ret *albv20251030.ListenerOutput, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := albv20251030.NewDescribeListenersRequest()
	request.LoadBalancerId = helper.String(loadBalancerId)
	request.ListenerIds = []*string{&listenerId}
	request.MaxResults = helper.Uint64(100)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for {
		var pageItems []*albv20251030.ListenerOutput
		nextToken := ""

		err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(request.GetAction())
			result, e := me.client.UseAlbV20251030Client().DescribeListenersWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			}

			if result == nil || result.Response == nil {
				return resource.NonRetryableError(fmt.Errorf("describe alb listeners failed, Response is nil"))
			}

			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())

			pageItems = result.Response.Listeners
			if result.Response.NextToken != nil {
				nextToken = *result.Response.NextToken
			}
			return nil
		})

		if err != nil {
			errRet = err
			return
		}

		for _, item := range pageItems {
			if item.ListenerId != nil && *item.ListenerId == listenerId {
				ret = item
				return
			}
		}

		if nextToken == "" {
			break
		}

		request.NextToken = helper.String(nextToken)
	}

	return
}

func (me *AlbService) ModifyListenerAttributes(ctx context.Context, request *albv20251030.ModifyListenerAttributesRequest) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseAlbV20251030Client().ModifyListenerAttributesWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("modify alb listener attributes failed, Response is nil"))
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		return nil
	})

	if err != nil {
		errRet = err
		return
	}

	return
}

func (me *AlbService) DeleteListener(ctx context.Context, request *albv20251030.DeleteListenerRequest) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseAlbV20251030Client().DeleteListenerWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("delete alb listener failed, Response is nil"))
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		return nil
	})

	if err != nil {
		errRet = err
		return
	}

	return
}
