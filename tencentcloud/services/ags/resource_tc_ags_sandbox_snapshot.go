package ags

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	ags "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ags/v20250920"
	sdkcommon "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/fw"
)

var (
	_ resource.ResourceWithConfigure = &AgsSandboxSnapshotResource{}
)

// NewAgsSandboxSnapshotResource is the factory referenced by
// tencentcloud/framework/registry.go to register this resource.
func NewAgsSandboxSnapshotResource() resource.Resource {
	return &AgsSandboxSnapshotResource{}
}

// AgsSandboxSnapshotResource implements resource.Resource for
// tencentcloud_ags_sandbox_snapshot. The AGS API provides no Modify
// interface, so every writable attribute forces replacement and the
// embedded WithNoUpdate rejects any Update call.
type AgsSandboxSnapshotResource struct {
	fw.ResourceWithConfigure
	fw.WithImportByID
	fw.WithNoUpdate
}

// AgsTagModel maps a tags object element.
type AgsTagModel struct {
	Key   types.String `tfsdk:"key"`
	Value types.String `tfsdk:"value"`
}

// AgsSandboxSnapshotResourceModel maps the schema attributes.
type AgsSandboxSnapshotResourceModel struct {
	Id           types.String  `tfsdk:"id"`
	InstanceId   types.String  `tfsdk:"instance_id"`
	SnapshotName types.String  `tfsdk:"snapshot_name"`
	Memory       types.Bool    `tfsdk:"memory"`
	Tags         []AgsTagModel `tfsdk:"tags"`
	ClientToken  types.String  `tfsdk:"client_token"`
	Status       types.String  `tfsdk:"status"`
	StatusReason types.String  `tfsdk:"status_reason"`
	CreateTime   types.String  `tfsdk:"create_time"`
	UpdateTime   types.String  `tfsdk:"update_time"`
}

func (r *AgsSandboxSnapshotResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "tencentcloud_ags_sandbox_snapshot"
}

func (r *AgsSandboxSnapshotResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a resource to create an AGS (TencentCloud Agent Sandbox) sandbox snapshot.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "ID of the sandbox snapshot, the same value as snapshot_id returned by the Create API.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"instance_id": schema.StringAttribute{
				Optional:    true,
				Description: "ID of the sandbox instance from which the snapshot is created.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"snapshot_name": schema.StringAttribute{
				Optional:    true,
				Description: "Name of the sandbox snapshot.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"memory": schema.BoolAttribute{
				Optional:    true,
				Description: "Whether the snapshot includes the memory state of the sandbox instance.",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
			},
			"tags": schema.ListNestedAttribute{
				Optional:    true,
				Description: "Tags of the sandbox snapshot.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"key": schema.StringAttribute{
							Optional:    true,
							Description: "Tag key.",
						},
						"value": schema.StringAttribute{
							Optional:    true,
							Description: "Tag value.",
						},
					},
				},
				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplace(),
				},
			},
			"client_token": schema.StringAttribute{
				Optional:    true,
				Description: "Idempotency token used when creating the snapshot.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"status": schema.StringAttribute{
				Computed:    true,
				Description: "Status of the sandbox snapshot. Valid values: ACTIVE, DELETING, DELETED.",
			},
			"status_reason": schema.StringAttribute{
				Computed:    true,
				Description: "Reason why the sandbox snapshot is in its current status.",
			},
			"create_time": schema.StringAttribute{
				Computed:    true,
				Description: "Creation time of the sandbox snapshot in ISO8601 format.",
			},
			"update_time": schema.StringAttribute{
				Computed:    true,
				Description: "Last update time of the sandbox snapshot in ISO8601 format.",
			},
		},
	}
}

// Create handles `terraform apply` for new resources.
func (r *AgsSandboxSnapshotResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan AgsSandboxSnapshotResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	request := ags.NewCreateSandboxSnapshotRequest()
	request.InstanceId = plan.InstanceId.ValueStringPointer()
	request.SnapshotName = plan.SnapshotName.ValueStringPointer()
	request.Memory = plan.Memory.ValueBoolPointer()
	request.ClientToken = plan.ClientToken.ValueStringPointer()
	if len(plan.Tags) > 0 {
		request.Tags = tagsToSdk(plan.Tags)
	}

	response, err := r.Client().UseAgsClient().CreateSandboxSnapshotWithContext(ctx, request)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating ags sandbox snapshot",
			fmt.Sprintf("CreateSandboxSnapshot failed: %v", err),
		)
		return
	}
	if response == nil || response.Response == nil || response.Response.Snapshot == nil ||
		response.Response.Snapshot.SnapshotId == nil || *response.Response.Snapshot.SnapshotId == "" {
		resp.Diagnostics.AddError(
			"Error creating ags sandbox snapshot",
			"CreateSandboxSnapshot returned an empty response or an empty snapshot id.",
		)
		return
	}

	// Persist the identity BEFORE any follow-up work: if the read-back fails,
	// the resource is still in state and `terraform destroy` can clean it up.
	plan.Id = types.StringValue(*response.Response.Snapshot.SnapshotId)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	item, err := r.read(ctx, plan.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading ags sandbox snapshot after create", err.Error())
		return
	}
	if item == nil {
		resp.Diagnostics.AddError(
			"ags sandbox snapshot not readable after create",
			fmt.Sprintf("snapshot %s was created but cannot be read back", plan.Id.ValueString()),
		)
		return
	}
	r.refreshModel(&plan, item)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Read handles `terraform refresh` and post-apply state refresh.
func (r *AgsSandboxSnapshotResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state AgsSandboxSnapshotResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	item, err := r.read(ctx, state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading ags sandbox snapshot", err.Error())
		return
	}
	if item == nil {
		resp.Diagnostics.AddWarning(
			"Resource not found",
			fmt.Sprintf("ags sandbox snapshot %s not found, removing from state.", state.Id.ValueString()),
		)
		resp.State.RemoveResource(ctx)
		return
	}
	r.refreshModel(&state, item)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Delete handles resource deletion.
func (r *AgsSandboxSnapshotResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state AgsSandboxSnapshotResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	request := ags.NewDeleteSandboxSnapshotRequest()
	request.SnapshotId = state.Id.ValueStringPointer()

	_, err := r.Client().UseAgsClient().DeleteSandboxSnapshotWithContext(ctx, request)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting ags sandbox snapshot",
			fmt.Sprintf("DeleteSandboxSnapshot failed: %v", err),
		)
		return
	}
}

// read fetches the snapshot by id via DescribeSandboxSnapshotList.
//
// Two different "empty" outcomes stay distinguishable:
//   - API failure (nil response / nil Response) -> error
//   - queried successfully but nothing matched  -> (nil, nil) => "does not exist"
func (r *AgsSandboxSnapshotResource) read(ctx context.Context, id string) (*ags.SandboxSnapshot, error) {
	request := ags.NewDescribeSandboxSnapshotListRequest()
	request.SnapshotIds = []*string{sdkcommon.StringPtr(id)}

	response, err := r.Client().UseAgsClient().DescribeSandboxSnapshotListWithContext(ctx, request)
	if err != nil {
		return nil, err
	}
	if response == nil || response.Response == nil {
		return nil, fmt.Errorf("DescribeSandboxSnapshotList returned an empty response for %s", id)
	}
	if len(response.Response.SnapshotSet) == 0 {
		return nil, nil
	}
	for _, snapshot := range response.Response.SnapshotSet {
		if snapshot != nil && snapshot.SnapshotId != nil && *snapshot.SnapshotId == id {
			return snapshot, nil
		}
	}
	return nil, nil
}

// refreshModel copies cloud-owned fields into the model.
//
// The Describe response echoes snapshot_name and tags but has no fields for
// instance_id, memory or client_token (SourceInstanceId and MemoryIncluded
// are name-similar candidates whose semantics are not confirmed), so those
// three keep their state values to avoid permanent diffs.
func (r *AgsSandboxSnapshotResource) refreshModel(m *AgsSandboxSnapshotResourceModel, item *ags.SandboxSnapshot) {
	if item == nil {
		return
	}
	if item.SnapshotId != nil && *item.SnapshotId != "" {
		m.Id = types.StringValue(*item.SnapshotId)
	}
	m.SnapshotName = types.StringPointerValue(item.SnapshotName)
	m.Tags = tagsFromSdk(item.Tags)
	m.Status = types.StringPointerValue(item.Status)
	m.StatusReason = types.StringPointerValue(item.StatusReason)
	m.CreateTime = types.StringPointerValue(item.CreateTime)
	m.UpdateTime = types.StringPointerValue(item.UpdateTime)
}

// tagsToSdk converts the model tag list into SDK Tag pointers.
func tagsToSdk(tags []AgsTagModel) []*ags.Tag {
	result := make([]*ags.Tag, 0, len(tags))
	for _, tag := range tags {
		result = append(result, &ags.Tag{
			Key:   tag.Key.ValueStringPointer(),
			Value: tag.Value.ValueStringPointer(),
		})
	}
	return result
}

// tagsFromSdk converts SDK Tag pointers into the model tag list.
func tagsFromSdk(tags []*ags.Tag) []AgsTagModel {
	if len(tags) == 0 {
		return nil
	}
	result := make([]AgsTagModel, 0, len(tags))
	for _, tag := range tags {
		if tag == nil {
			continue
		}
		result = append(result, AgsTagModel{
			Key:   types.StringPointerValue(tag.Key),
			Value: types.StringPointerValue(tag.Value),
		})
	}
	return result
}
