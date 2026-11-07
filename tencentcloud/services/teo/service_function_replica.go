package teo

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

// CreateTeoFunctionReplica creates a TEO edge function replica. The
// CreateFunctionReplica API returns only a RequestId, so the composite
// identity is assembled by the caller from the request itself.
func (me *TeoService) CreateTeoFunctionReplica(ctx context.Context, request *teov20220901.CreateFunctionReplicaRequest) (errRet error) {
	if request == nil {
		errRet = fmt.Errorf("CreateTeoFunctionReplica: request is nil")
		return
	}

	var (
		logId = tccommon.GetLogId(ctx)
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseTeoV20220901Client().CreateFunctionReplicaWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		}
		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("CreateFunctionReplica returned an empty response"))
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		return nil
	})
	if err != nil {
		errRet = err
		return
	}

	return
}

// DescribeTeoFunctionReplicaByName queries the single function replica
// matching replicaName under zoneId/functionId. A nil return (with nil
// error) means the replica is confirmed absent from the listing result.
func (me *TeoService) DescribeTeoFunctionReplicaByName(ctx context.Context, zoneId, functionId, replicaName string) (replica *teov20220901.FunctionReplica, errRet error) {
	if zoneId == "" || functionId == "" || replicaName == "" {
		errRet = fmt.Errorf("DescribeTeoFunctionReplicaByName: zoneId, functionId and replicaName are all required to locate a replica")
		return
	}

	var (
		logId   = tccommon.GetLogId(ctx)
		request = teov20220901.NewDescribeFunctionReplicasRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	request.ZoneId = helper.String(zoneId)
	request.FunctionId = helper.String(functionId)
	request.Filters = []*teov20220901.AdvancedFilter{
		{
			Name:   helper.String("replica-name"),
			Values: []*string{helper.String(replicaName)},
		},
	}
	request.Limit = helper.Int64(200)

	var response *teov20220901.DescribeFunctionReplicasResponse
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseTeoV20220901Client().DescribeFunctionReplicasWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		}
		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("DescribeFunctionReplicas returned an empty response"))
		}
		response = result
		return nil
	})
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	// An empty listing is a valid "not found" outcome, not an API failure.
	for _, item := range response.Response.FunctionReplicas {
		if item != nil && item.ReplicaName != nil && *item.ReplicaName == replicaName {
			replica = item
			break
		}
	}

	return
}

// ModifyTeoFunctionReplica updates a TEO edge function replica with the
// fields populated on the request. Callers must fill only the fields whose
// value actually changed; an unset field is not sent to the API.
func (me *TeoService) ModifyTeoFunctionReplica(ctx context.Context, request *teov20220901.ModifyFunctionReplicaRequest) (errRet error) {
	if request == nil {
		errRet = fmt.Errorf("ModifyTeoFunctionReplica: request is nil")
		return
	}

	var (
		logId = tccommon.GetLogId(ctx)
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseTeoV20220901Client().ModifyFunctionReplicaWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		}
		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("ModifyFunctionReplica returned an empty response"))
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		return nil
	})
	if err != nil {
		errRet = err
		return
	}

	return
}

// DeleteTeoFunctionReplicaByName deletes the single function replica named
// replicaName under zoneId/functionId. The DeleteFunctionReplica API takes
// a list, but a managed resource owns exactly one replica, so only a single
// element is passed to avoid deleting anything beyond the managed object.
func (me *TeoService) DeleteTeoFunctionReplicaByName(ctx context.Context, zoneId, functionId, replicaName string) (errRet error) {
	if zoneId == "" || functionId == "" || replicaName == "" {
		errRet = fmt.Errorf("DeleteTeoFunctionReplicaByName: zoneId, functionId and replicaName are all required to delete a replica")
		return
	}

	var (
		logId   = tccommon.GetLogId(ctx)
		request = teov20220901.NewDeleteFunctionReplicaRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	request.ZoneId = helper.String(zoneId)
	request.FunctionId = helper.String(functionId)
	request.ReplicaNames = []*string{helper.String(replicaName)}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseTeoV20220901Client().DeleteFunctionReplicaWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		}
		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("DeleteFunctionReplica returned an empty response"))
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		return nil
	})
	if err != nil {
		errRet = err
		return
	}

	return
}

// splitTeoFunctionReplicaV8Id parses the composite resource id
// `<zone_id>#<function_id>#<replica_name>` built with tccommon.FILED_SP.
func splitTeoFunctionReplicaV8Id(id string) (zoneId, functionId, replicaName string, ok bool) {
	parts := strings.Split(id, tccommon.FILED_SP)
	if len(parts) != 3 {
		return "", "", "", false
	}
	for _, p := range parts {
		if p == "" {
			return "", "", "", false
		}
	}
	return parts[0], parts[1], parts[2], true
}
