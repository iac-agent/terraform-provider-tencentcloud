package ags

import (
	"context"
	"fmt"
	"log"

	agsv20250920 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ags/v20250920"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

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
)

const (
	// AGS_GPU_RESOURCE_GROUP_FILTER_LIMIT is the maximum number of filters
	// accepted by DescribeGPUResourceGroups.
	AGS_GPU_RESOURCE_GROUP_FILTER_LIMIT = 10
)

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

// ModifyGPUResourceGroup renames an existing GPU resource group.
func (me *AgsService) ModifyGPUResourceGroup(ctx context.Context, request *agsv20250920.ModifyGPUResourceGroupRequest) (errRet error) {
	logId := tccommon.GetLogId(ctx)
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	var response *agsv20250920.ModifyGPUResourceGroupResponse
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, err := me.client.UseAgsClient().ModifyGPUResourceGroupWithContext(ctx, request)
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
		errRet = fmt.Errorf("ModifyGPUResourceGroup returned an empty response")
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

// DescribeGPUResourceGroups queries GPU resource groups with the given filters
// and returns the matched groups together with the total match count.
func (me *AgsService) DescribeGPUResourceGroups(ctx context.Context, request *agsv20250920.DescribeGPUResourceGroupsRequest) (groups []*agsv20250920.GPUResourceGroup, totalCount int64, errRet error) {
	logId := tccommon.GetLogId(ctx)
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	var response *agsv20250920.DescribeGPUResourceGroupsResponse
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, err := me.client.UseAgsClient().DescribeGPUResourceGroupsWithContext(ctx, request)
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
		errRet = fmt.Errorf("DescribeGPUResourceGroups returned an empty response")
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response.Response.ResourceGroups != nil {
		groups = response.Response.ResourceGroups
	}
	if response.Response.TotalCount != nil {
		totalCount = *response.Response.TotalCount
	}

	return
}

// DescribeGPUResourceGroupById returns the GPU resource group with the given
// id. A nil result with a nil error means "queried successfully, no such
// group"; an API failure is always reported as an error so the caller never
// confuses a failed request with a released resource.
func (me *AgsService) DescribeGPUResourceGroupById(ctx context.Context, resourceGroupId string) (group *agsv20250920.GPUResourceGroup, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := agsv20250920.NewDescribeGPUResourceGroupsRequest()
	filter := &agsv20250920.Filter{}
	filter.Name = helper.String("resource-group-id")
	filter.Values = []*string{helper.String(resourceGroupId)}
	request.Filters = []*agsv20250920.Filter{filter}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	groups, _, err := me.DescribeGPUResourceGroups(ctx, request)
	if err != nil {
		errRet = err
		return
	}

	for _, item := range groups {
		if item != nil && item.ResourceGroupId != nil && *item.ResourceGroupId == resourceGroupId {
			group = item
			return
		}
	}

	return
}

// DeleteGPUResourceGroup releases an empty GPU resource group identified by
// resourceGroupId. The API rejects groups that still contain valid GPU
// instances, occupations or bound sandbox tools.
func (me *AgsService) DeleteGPUResourceGroup(ctx context.Context, resourceGroupId string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := agsv20250920.NewDeleteGPUResourceGroupRequest()
	request.ResourceGroupId = helper.String(resourceGroupId)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	var response *agsv20250920.DeleteGPUResourceGroupResponse
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, err := me.client.UseAgsClient().DeleteGPUResourceGroupWithContext(ctx, request)
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
		errRet = fmt.Errorf("DeleteGPUResourceGroup returned an empty response for GPU resource group %s", resourceGroupId)
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}
