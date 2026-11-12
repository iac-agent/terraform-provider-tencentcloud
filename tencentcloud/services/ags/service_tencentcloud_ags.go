package ags

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	agsv20250920 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ags/v20250920"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

const (
	// AGS_SANDBOX_TOOL_DESCRIBE_LIMIT is the page size used when locating a
	// single sandbox tool by id. The API accepts at most 100 ids per request.
	AGS_SANDBOX_TOOL_DESCRIBE_LIMIT = 100

	// Sandbox tool status values documented by DescribeSandboxToolList.
	AGS_SANDBOX_TOOL_STATUS_ACTIVE   = "ACTIVE"
	AGS_SANDBOX_TOOL_STATUS_FAILED   = "FAILED"
	AGS_SANDBOX_TOOL_STATUS_CREATING = "CREATING"
	AGS_SANDBOX_TOOL_STATUS_DELETING = "DELETING"

	// AGS_PRE_CACHE_IMAGE_TYPE_IMAGE is the default value of PreCacheType used
	// when a pre-cache image task is read back by id.
	AGS_PRE_CACHE_IMAGE_TYPE_IMAGE = "image"
)

// AgsPreCacheImageTask is the pre-cache image object located by
// DescribePreCacheImageTask.
//
// The describe API answers with a single object rather than a set, so the
// "not found" signal is an empty PreCacheImageId instead of an empty set.
type AgsPreCacheImageTask struct {
	PreCacheImageId      string
	Image                string
	ImageDigest          string
	ImageRegistryType    string
	Status               string
	Message              string
	CreateTime           string
	ReadyTime            string
	LastUsedTime         string
	SourceType           string
	CachedImageSizeBytes int64
}

func NewAgsService(client *connectivity.TencentCloudClient) AgsService {
	return AgsService{client: client}
}

type AgsService struct {
	client *connectivity.TencentCloudClient
}

// CreateSandboxTool creates a sandbox tool and returns the id assigned by the
// cloud side. The id is the only resource identity: RequestId is only a request
// trace id and must never be used as an identity.
func (me *AgsService) CreateSandboxTool(ctx context.Context, request *agsv20250920.CreateSandboxToolRequest) (toolId string, errRet error) {
	logId := tccommon.GetLogId(ctx)
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	var response *agsv20250920.CreateSandboxToolResponse
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, err := me.client.UseAgsClient().CreateSandboxToolWithContext(ctx, request)
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
		errRet = fmt.Errorf("CreateSandboxTool returned an empty response")
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response.Response.ToolId == nil || *response.Response.ToolId == "" {
		errRet = fmt.Errorf("CreateSandboxTool returned an empty tool id")
		return
	}

	toolId = *response.Response.ToolId
	return
}

// DescribeSandboxToolById returns the sandbox tool with the given id.
//
// A nil result with a nil error means "queried successfully, no such tool"; an
// API failure is always reported as an error so the caller never confuses a
// failed request with a deleted resource.
func (me *AgsService) DescribeSandboxToolById(ctx context.Context, toolId string) (sandboxTool *agsv20250920.SandboxTool, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := agsv20250920.NewDescribeSandboxToolListRequest()
	request.ToolIds = []*string{helper.String(toolId)}
	request.Limit = helper.Int64(AGS_SANDBOX_TOOL_DESCRIBE_LIMIT)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	var response *agsv20250920.DescribeSandboxToolListResponse
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, err := me.client.UseAgsClient().DescribeSandboxToolListWithContext(ctx, request)
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
		errRet = fmt.Errorf("DescribeSandboxToolList returned an empty response for sandbox tool %s", toolId)
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	for _, tool := range response.Response.SandboxToolSet {
		if tool != nil && tool.ToolId != nil && *tool.ToolId == toolId {
			sandboxTool = tool
			return
		}
	}

	return
}

// UpdateSandboxTool applies an in-place update to an existing sandbox tool.
func (me *AgsService) UpdateSandboxTool(ctx context.Context, request *agsv20250920.UpdateSandboxToolRequest) (errRet error) {
	logId := tccommon.GetLogId(ctx)
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	var response *agsv20250920.UpdateSandboxToolResponse
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, err := me.client.UseAgsClient().UpdateSandboxToolWithContext(ctx, request)
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
		errRet = fmt.Errorf("UpdateSandboxTool returned an empty response")
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

// CreatePreCacheImageTask creates a pre-cache image task and returns the id
// assigned by the cloud side. The id is the only resource identity: RequestId
// is only a request trace id and must never be used as an identity.
func (me *AgsService) CreatePreCacheImageTask(ctx context.Context, request *agsv20250920.CreatePreCacheImageTaskRequest) (preCacheImageId string, errRet error) {
	logId := tccommon.GetLogId(ctx)
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	var response *agsv20250920.CreatePreCacheImageTaskResponse
	// The API does not document an idempotency key for create, so the request
	// is not retried here: a retried create could leave a second pre-cache
	// image task behind that this resource never tracks.
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseAgsClient().CreatePreCacheImageTaskWithContext(ctx, request)
	if err != nil {
		errRet = err
		return
	}
	if response == nil || response.Response == nil {
		errRet = fmt.Errorf("CreatePreCacheImageTask returned an empty response")
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response.Response.PreCacheImageId == nil || *response.Response.PreCacheImageId == "" {
		errRet = fmt.Errorf("CreatePreCacheImageTask returned an empty pre cache image id")
		return
	}

	preCacheImageId = *response.Response.PreCacheImageId
	return
}

// DescribePreCacheImageTaskById returns the pre-cache image task with the given
// id.
//
// A nil result with a nil error means "queried successfully, no such
// pre-cache image task"; an API failure is always reported as an error so the
// caller never confuses a failed request with a deleted resource.
func (me *AgsService) DescribePreCacheImageTaskById(ctx context.Context, preCacheImageId string) (preCacheImageTask *AgsPreCacheImageTask, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := agsv20250920.NewDescribePreCacheImageTaskRequest()
	request.PreCacheImageId = helper.String(preCacheImageId)
	// PreCacheType defaults to "image" on the cloud side; it is sent
	// explicitly because this resource only manages image pre-caching.
	request.PreCacheType = helper.String(AGS_PRE_CACHE_IMAGE_TYPE_IMAGE)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	var response *agsv20250920.DescribePreCacheImageTaskResponse
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, err := me.client.UseAgsClient().DescribePreCacheImageTaskWithContext(ctx, request)
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
		errRet = fmt.Errorf("DescribePreCacheImageTask returned an empty response for pre cache image %s", preCacheImageId)
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	item := response.Response
	// An empty id means the requested pre-cache image is not reported by the
	// cloud side any more.
	if item.PreCacheImageId == nil || *item.PreCacheImageId == "" {
		return
	}

	preCacheImageTask = &AgsPreCacheImageTask{
		PreCacheImageId:   *item.PreCacheImageId,
		Image:             stringValueOf(item.Image),
		ImageDigest:       stringValueOf(item.ImageDigest),
		ImageRegistryType: stringValueOf(item.ImageRegistryType),
		Status:            stringValueOf(item.Status),
		Message:           stringValueOf(item.Message),
		CreateTime:        stringValueOf(item.CreateTime),
		ReadyTime:         stringValueOf(item.ReadyTime),
		LastUsedTime:      stringValueOf(item.LastUsedTime),
		SourceType:        stringValueOf(item.SourceType),
	}
	if item.CachedImageSizeBytes != nil {
		preCacheImageTask.CachedImageSizeBytes = *item.CachedImageSizeBytes
	}

	return
}

// DeletePreCacheImage deletes the pre-cache image identified by request.
func (me *AgsService) DeletePreCacheImage(ctx context.Context, request *agsv20250920.DeletePreCacheImageRequest) (errRet error) {
	logId := tccommon.GetLogId(ctx)
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	// The delete API is documented as "delete this pre-cache image"; it is
	// treated as idempotent, so a transient failure is retried.
	var response *agsv20250920.DeletePreCacheImageResponse
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, err := me.client.UseAgsClient().DeletePreCacheImageWithContext(ctx, request)
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
		errRet = fmt.Errorf("DeletePreCacheImage returned an empty response")
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

// stringValueOf dereferences an optional SDK string field.
func stringValueOf(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

// DeleteSandboxTool deletes the sandbox tool itself, identified by toolId.
func (me *AgsService) DeleteSandboxTool(ctx context.Context, toolId string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := agsv20250920.NewDeleteSandboxToolRequest()
	request.ToolId = helper.String(toolId)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	var response *agsv20250920.DeleteSandboxToolResponse
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, err := me.client.UseAgsClient().DeleteSandboxToolWithContext(ctx, request)
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
		errRet = fmt.Errorf("DeleteSandboxTool returned an empty response for sandbox tool %s", toolId)
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}
