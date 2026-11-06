package ags

import (
	"context"
	"fmt"
	"log"

	v20250920 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ags/v20250920"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

const (
	AGS_SANDBOX_TOOL_STATUS_CREATING = "CREATING"
	AGS_SANDBOX_TOOL_STATUS_ACTIVE   = "ACTIVE"
	AGS_SANDBOX_TOOL_STATUS_DELETING = "DELETING"
	AGS_SANDBOX_TOOL_STATUS_FAILED   = "FAILED"

	// AGS_DESCRIBE_SANDBOX_TOOL_LIMIT is the maximum page size accepted by
	// DescribeSandboxToolList (documented maximum is 100).
	AGS_DESCRIBE_SANDBOX_TOOL_LIMIT = 100
)

func NewAgsService(client *connectivity.TencentCloudClient) AgsService {
	return AgsService{client: client}
}

type AgsService struct {
	client *connectivity.TencentCloudClient
}

// CreateSandboxTool creates a sandbox tool and returns the new tool id.
func (me *AgsService) CreateSandboxTool(ctx context.Context, request *v20250920.CreateSandboxToolRequest) (toolId string, errRet error) {
	logId := tccommon.GetLogId(ctx)

	ratelimit.Check(request.GetAction())

	var response *v20250920.CreateSandboxToolResponse
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	}()

	response, errRet = me.client.UseAgsClient().CreateSandboxToolWithContext(ctx, request)
	if errRet != nil {
		return
	}
	if response == nil || response.Response == nil {
		errRet = fmt.Errorf("CreateSandboxTool returned an empty response")
		return
	}
	if response.Response.ToolId == nil || *response.Response.ToolId == "" {
		errRet = fmt.Errorf("CreateSandboxTool returned an empty ToolId")
		return
	}

	toolId = *response.Response.ToolId
	return
}

// DescribeSandboxToolById returns the sandbox tool with the given id, or
// (nil, nil) when the tool does not exist on the cloud side.
func (me *AgsService) DescribeSandboxToolById(ctx context.Context, toolId string) (sandboxTool *v20250920.SandboxTool, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := v20250920.NewDescribeSandboxToolListRequest()
	request.ToolIds = []*string{helper.String(toolId)}
	request.Limit = helper.Int64(AGS_DESCRIBE_SANDBOX_TOOL_LIMIT)

	ratelimit.Check(request.GetAction())

	var response *v20250920.DescribeSandboxToolListResponse
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	}()

	response, errRet = me.client.UseAgsClient().DescribeSandboxToolListWithContext(ctx, request)
	if errRet != nil {
		return
	}
	if response == nil || response.Response == nil {
		errRet = fmt.Errorf("DescribeSandboxToolList returned an empty response for tool %s", toolId)
		return
	}

	for _, item := range response.Response.SandboxToolSet {
		if item == nil || item.ToolId == nil {
			continue
		}
		if *item.ToolId == toolId {
			sandboxTool = item
			return
		}
	}

	return
}

// UpdateSandboxTool updates the mutable attributes of a sandbox tool.
func (me *AgsService) UpdateSandboxTool(ctx context.Context, request *v20250920.UpdateSandboxToolRequest) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	ratelimit.Check(request.GetAction())

	var response *v20250920.UpdateSandboxToolResponse
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	}()

	response, errRet = me.client.UseAgsClient().UpdateSandboxToolWithContext(ctx, request)
	if errRet != nil {
		return
	}
	if response == nil {
		errRet = fmt.Errorf("UpdateSandboxTool returned an empty response")
		return
	}

	return
}

// DeleteSandboxTool deletes a sandbox tool by id.
func (me *AgsService) DeleteSandboxTool(ctx context.Context, toolId string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := v20250920.NewDeleteSandboxToolRequest()
	request.ToolId = helper.String(toolId)

	ratelimit.Check(request.GetAction())

	var response *v20250920.DeleteSandboxToolResponse
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	}()

	response, errRet = me.client.UseAgsClient().DeleteSandboxToolWithContext(ctx, request)
	if errRet != nil {
		return
	}
	if response == nil {
		errRet = fmt.Errorf("DeleteSandboxTool returned an empty response")
		return
	}

	return
}
