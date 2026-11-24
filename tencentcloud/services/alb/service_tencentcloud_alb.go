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
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

// AlbJobStatusTerminal marks the terminal status values reported by
// DescribeAsyncJobs, which every ALB write API (CreateRules,
// ModifyRulesAttributes, DeleteRules) relies on as an asynchronous interface.
const AlbJobStatusTerminal = "Succeeded"

// AlbJobStatusFailed marks the terminal failure status of an ALB async job.
const AlbJobStatusFailed = "Failed"

func NewAlbService(client *connectivity.TencentCloudClient) AlbService {
	return AlbService{client: client}
}

type AlbService struct {
	client *connectivity.TencentCloudClient
}

// CreateRules creates forwarding rules under one listener. The API is
// asynchronous: the returned request id is waited on through DescribeAsyncJobs
// before the call returns.
func (me *AlbService) CreateRules(ctx context.Context, request *albv20251030.CreateRulesRequest) (ruleIds []string, err error) {
	var (
		logId    = tccommon.GetLogId(ctx)
		response *albv20251030.CreateRulesResponse
	)

	defer func() {
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), err.Error())
		}
	}()

	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseAlbV20251030Client().CreateRulesWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("create alb rules failed, Response is nil"))
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())

		if len(result.Response.RuleIds) == 0 {
			return resource.NonRetryableError(fmt.Errorf("create alb rules failed, RuleIds is empty"))
		}

		response = result
		return nil
	})

	if err != nil {
		return
	}

	for _, id := range response.Response.RuleIds {
		if id != nil {
			ruleIds = append(ruleIds, *id)
		}
	}

	if err = me.WaitJob(ctx, *response.Response.RequestId); err != nil {
		return
	}

	return
}

// ModifyRulesAttributes updates forwarding rule attributes. The API is
// asynchronous: the returned request id is waited on through DescribeAsyncJobs
// before the call returns.
func (me *AlbService) ModifyRulesAttributes(ctx context.Context, request *albv20251030.ModifyRulesAttributesRequest) (err error) {
	var (
		logId    = tccommon.GetLogId(ctx)
		response *albv20251030.ModifyRulesAttributesResponse
	)

	defer func() {
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), err.Error())
		}
	}()

	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseAlbV20251030Client().ModifyRulesAttributesWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("modify alb rule attributes failed, Response is nil"))
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())

		response = result
		return nil
	})

	if err != nil {
		return
	}

	if err = me.WaitJob(ctx, *response.Response.RequestId); err != nil {
		return
	}

	return
}

// DeleteRules deletes forwarding rules. The API is asynchronous: the returned
// request id is waited on through DescribeAsyncJobs before the call returns.
func (me *AlbService) DeleteRules(ctx context.Context, request *albv20251030.DeleteRulesRequest) (err error) {
	var (
		logId    = tccommon.GetLogId(ctx)
		response *albv20251030.DeleteRulesResponse
	)

	defer func() {
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), err.Error())
		}
	}()

	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseAlbV20251030Client().DeleteRulesWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("delete alb rules failed, Response is nil"))
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())

		response = result
		return nil
	})

	if err != nil {
		return
	}

	if err = me.WaitJob(ctx, *response.Response.RequestId); err != nil {
		return
	}

	return
}

// DescribeRules returns the forwarding rules matching the request. The caller
// picks its own items from the returned list.
func (me *AlbService) DescribeRules(ctx context.Context, request *albv20251030.DescribeRulesRequest) (rules []*albv20251030.RuleOutput, nextToken string, err error) {
	var (
		logId    = tccommon.GetLogId(ctx)
		response *albv20251030.DescribeRulesResponse
	)

	defer func() {
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), err.Error())
		}
	}()

	err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseAlbV20251030Client().DescribeRulesWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("describe alb rules failed, Response is nil"))
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())

		response = result
		return nil
	})

	if err != nil {
		return
	}

	rules = response.Response.Rules
	if response.Response.NextToken != nil {
		nextToken = *response.Response.NextToken
	}

	return
}

// WaitJob polls DescribeAsyncJobs until the asynchronous job of the given
// request id reaches the Succeeded status.
func (me *AlbService) WaitJob(ctx context.Context, requestId string) (err error) {
	var (
		logId = tccommon.GetLogId(ctx)
	)

	defer func() {
		if err != nil {
			log.Printf("[CRITAL]%s wait alb job[%s] fail, reason[%s]\n", logId, requestId, err.Error())
		}
	}()

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		request := albv20251030.NewDescribeAsyncJobsRequest()
		request.RequestIds = []*string{&requestId}
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseAlbV20251030Client().DescribeAsyncJobsWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("describe alb async jobs failed, Response is nil"))
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())

		if len(result.Response.Jobs) == 0 {
			return resource.RetryableError(fmt.Errorf("alb job %s is not queryable yet", requestId))
		}

		job := result.Response.Jobs[0]
		if job.Status != nil && *job.Status == AlbJobStatusTerminal {
			return nil
		}

		if job.Status != nil && *job.Status == AlbJobStatusFailed {
			return resource.NonRetryableError(fmt.Errorf("alb job %s failed", requestId))
		}

		return resource.RetryableError(fmt.Errorf("alb job %s status is %s", requestId, *job.Status))
	})

	return
}
