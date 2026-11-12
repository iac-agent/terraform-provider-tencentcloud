package ags

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	sdkags "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ags/v20250920"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/fw"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

var (
	_ resource.ResourceWithConfigure   = &AgsRegistryResource{}
	_ resource.ResourceWithImportState = &AgsRegistryResource{}
)

// NewAgsRegistryResource is the factory referenced by
// tencentcloud/framework/registry.go to register this resource.
func NewAgsRegistryResource() resource.Resource {
	return &AgsRegistryResource{}
}

// AgsRegistryResource implements resource.Resource for
// tencentcloud_ags_registry.
//
// The cloud product exposes Create/Describe/Delete operations only; there is
// no Modify API, so every configurable attribute is immutable and marked with
// RequiresReplace. The Update method (embedded via fw.WithNoUpdate) rejects
// any call, which cannot happen because all changes force replacement.
type AgsRegistryResource struct {
	fw.ResourceWithConfigure
	fw.WithImportByID
	fw.WithNoUpdate
}

// AgsRegistryTagModel maps one element of the tags list attribute.
type AgsRegistryTagModel struct {
	Key   types.String `tfsdk:"key"`
	Value types.String `tfsdk:"value"`
}

// AgsRegistryInfoModel maps the computed registry single-nested attribute.
type AgsRegistryInfoModel struct {
	RegistryId            types.String `tfsdk:"registry_id"`
	Name                  types.String `tfsdk:"name"`
	Description           types.String `tfsdk:"description"`
	ApprovalMode          types.String `tfsdk:"approval_mode"`
	Region                types.String `tfsdk:"region"`
	Status                types.String `tfsdk:"status"`
	CreateTime            types.String `tfsdk:"create_time"`
	UpdateTime            types.String `tfsdk:"update_time"`
	RecordCount           types.Int64  `tfsdk:"record_count"`
	PublishedRecordCount  types.Int64  `tfsdk:"published_record_count"`
	AppId                 types.Int64  `tfsdk:"app_id"`
	CreatorUin            types.String `tfsdk:"creator_uin"`
	CreatorSubAccountUin  types.String `tfsdk:"creator_sub_account_uin"`
	Tags                  types.List   `tfsdk:"tags"`
}

// AgsRegistryResourceModel maps the schema attributes of the resource.
type AgsRegistryResourceModel struct {
	Id          types.String `tfsdk:"id"`
	RegistryId  types.String `tfsdk:"registry_id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	ApprovalMode types.String `tfsdk:"approval_mode"`
	Tags        types.List   `tfsdk:"tags"`
	Registry    types.Object `tfsdk:"registry"`
}

func (r *AgsRegistryResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "tencentcloud_ags_registry"
}

// tagsAttrTypes returns the attribute types of one tags element object.
func tagsAttrTypes() map[string]basetypes.AttrType {
	return map[string]basetypes.AttrType{
		"key":   basetypes.StringType{},
		"value": basetypes.StringType{},
	}
}

// tagsObjectType returns the element type used by both tags list attributes.
func tagsObjectType() basetypes.ObjectType {
	return basetypes.ObjectType{AttrTypes: tagsAttrTypes()}
}

// registryAttrTypes returns the attribute types of the computed registry object.
func registryAttrTypes() map[string]basetypes.AttrType {
	return map[string]basetypes.AttrType{
		"registry_id":              basetypes.StringType{},
		"name":                     basetypes.StringType{},
		"description":              basetypes.StringType{},
		"approval_mode":            basetypes.StringType{},
		"region":                   basetypes.StringType{},
		"status":                   basetypes.StringType{},
		"create_time":              basetypes.StringType{},
		"update_time":              basetypes.StringType{},
		"record_count":             basetypes.Int64Type{},
		"published_record_count":   basetypes.Int64Type{},
		"app_id":                   basetypes.Int64Type{},
		"creator_uin":              basetypes.StringType{},
		"creator_sub_account_uin":  basetypes.StringType{},
		"tags":                     basetypes.ListType{ElemType: tagsObjectType()},
	}
}

func (r *AgsRegistryResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a resource to create an AGS (Tencent Cloud Agent Gateway Service) registry.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "ID of the registry, same as registry_id.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"registry_id": schema.StringAttribute{
				Computed:    true,
				Description: "ID of the created registry, in the format `reg-` plus 8 lowercase letters/digits.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Name of the registry, unique within the same AppId and region, 1-255 characters long. Immutable after creation; changing it forces replacement.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "Description text, up to 4096 characters. Immutable after creation; changing it forces replacement.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"approval_mode": schema.StringAttribute{
				Optional:    true,
				Description: "Approval mode, decided at creation and cannot be changed afterwards; omitted defaults to AUTO. Valid values: `AUTO` (auto approve), `MANUAL` (manual approval). Case sensitive. Immutable after creation; changing it forces replacement.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"tags": schema.ListNestedAttribute{
				Optional:    true,
				Description: "Tencent Cloud custom tags bound at creation; keys must be unique, up to 10 tags. Immutable after creation; changing it forces replacement.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"key": schema.StringAttribute{
							Required:    true,
							Description: "Tag key. Reserved prefixes such as `qcs:`, `project` and project reserved prefixes are not allowed, and leading/trailing spaces are not allowed.",
						},
						"value": schema.StringAttribute{
							Required:    true,
							Description: "Tag value. Leading/trailing spaces are not allowed.",
						},
					},
				},
				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplace(),
				},
			},
			"registry": schema.SingleNestedAttribute{
				Computed:    true,
				Description: "Detailed information of the registry.",
				Attributes: map[string]schema.Attribute{
					"registry_id": schema.StringAttribute{
						Computed:    true,
						Description: "Registry ID, in the format `reg-` plus 8 lowercase letters/digits.",
					},
					"name": schema.StringAttribute{
						Computed:    true,
						Description: "Registry name, unique within the same AppId and region.",
					},
					"description": schema.StringAttribute{
						Computed:    true,
						Description: "Description of the registry.",
					},
					"approval_mode": schema.StringAttribute{
						Computed:    true,
						Description: "Approval mode. Valid values: `AUTO`, `MANUAL`. Decided at creation and cannot be changed.",
					},
					"region": schema.StringAttribute{
						Computed:    true,
						Description: "Tencent Cloud region where the registry resides, e.g. `ap-guangzhou`.",
					},
					"status": schema.StringAttribute{
						Computed:    true,
						Description: "Status of the registry. Valid values: `ACTIVE`, `ARCHIVED`.",
					},
					"create_time": schema.StringAttribute{
						Computed:    true,
						Description: "Creation time in ISO 8601 UTC, e.g. `2026-08-11T10:00:00Z`.",
					},
					"update_time": schema.StringAttribute{
						Computed:    true,
						Description: "Time of the most recent update, ISO 8601 UTC.",
					},
					"record_count": schema.Int64Attribute{
						Computed:    true,
						Description: "Total number of records under the registry.",
					},
					"published_record_count": schema.Int64Attribute{
						Computed:    true,
						Description: "Number of records bound to stable labels.",
					},
					"app_id": schema.Int64Attribute{
						Computed:    true,
						Description: "AppId of the owning account.",
					},
					"creator_uin": schema.StringAttribute{
						Computed:    true,
						Description: "UIN of the creator main account.",
					},
					"creator_sub_account_uin": schema.StringAttribute{
						Computed:    true,
						Description: "UIN of the creator sub account; empty string when created directly by the main account.",
					},
					"tags": schema.ListNestedAttribute{
						Computed:    true,
						Description: "Authoritative custom tags, stably sorted by key and value; an empty array is returned when there are no tags.",
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"key": schema.StringAttribute{
									Computed:    true,
									Description: "Tag key.",
								},
								"value": schema.StringAttribute{
									Computed:    true,
									Description: "Tag value.",
								},
							},
						},
					},
				},
			},
		},
	}
}

// Create handles `terraform apply` for new resources.
func (r *AgsRegistryResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan AgsRegistryResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	request := sdkags.NewCreateRegistryRequest()
	request.Name = plan.Name.ValueStringPointer()
	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		request.Description = plan.Description.ValueStringPointer()
	}
	if !plan.ApprovalMode.IsNull() && !plan.ApprovalMode.IsUnknown() {
		request.ApprovalMode = plan.ApprovalMode.ValueStringPointer()
	}
	if !plan.Tags.IsNull() && !plan.Tags.IsUnknown() {
		tags, diags := expandTags(ctx, plan.Tags)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		request.Tags = tags
	}

	response, err := r.Client().UseAgsClient().CreateRegistryWithContext(ctx, request)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating AGS registry",
			fmt.Sprintf("CreateRegistry failed: %v", err),
		)
		return
	}

	// A nil response (or a nil id field) is an API failure, not an empty result.
	// Never write an empty id into state: the resource would become unmanageable.
	if response == nil || response.Response == nil || response.Response.RegistryId == nil || *response.Response.RegistryId == "" {
		resp.Diagnostics.AddError(
			"Error creating AGS registry",
			"CreateRegistry returned an empty response or an empty registry id.",
		)
		return
	}

	plan.Id = types.StringValue(*response.Response.RegistryId)
	plan.RegistryId = types.StringValue(*response.Response.RegistryId)

	// Persist the identity BEFORE any follow-up read: if the read fails, the
	// resource is still in state and `terraform destroy` can clean it up.
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	item, err := r.read(ctx, plan.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading AGS registry after create", err.Error())
		return
	}
	if item == nil {
		resp.Diagnostics.AddError(
			"AGS registry not readable after create",
			fmt.Sprintf("registry %s was created but cannot be read back", plan.Id.ValueString()),
		)
		return
	}
	refreshRegistryModel(&plan, item)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Read handles `terraform refresh` and post-apply state refresh.
func (r *AgsRegistryResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state AgsRegistryResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	item, err := r.read(ctx, state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading AGS registry", err.Error())
		return
	}

	// Confirmed "not found" on the cloud side => drop it from state so Terraform
	// plans a recreation. A nil response is an API failure and must not land here.
	if item == nil {
		resp.Diagnostics.AddWarning(
			"Resource not found",
			fmt.Sprintf("AGS registry %s not found, removing from state.", state.Id.ValueString()),
		)
		resp.State.RemoveResource(ctx)
		return
	}

	refreshRegistryModel(&state, item)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Delete handles resource deletion.
func (r *AgsRegistryResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state AgsRegistryResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	request := sdkags.NewDeleteRegistryRequest()
	request.RegistryId = state.Id.ValueStringPointer()

	_, err := r.Client().UseAgsClient().DeleteRegistryWithContext(ctx, request)
	if err != nil {
		// No error code is documented as "already gone" for DeleteRegistry,
		// so every error must fail and keep the id in state for a retry.
		resp.Diagnostics.AddError(
			"Error deleting AGS registry",
			fmt.Sprintf("DeleteRegistry failed: %v", err),
		)
		return
	}
	// After a successful delete the Framework removes the state automatically.
}

// read fetches the registry by ID.
//
// Two different "empty" outcomes must stay distinguishable:
//   - API failure (nil response / nil Response)  -> error
//   - queried successfully but nothing matched  -> (nil, nil) => "does not exist"
func (r *AgsRegistryResource) read(ctx context.Context, id string) (*sdkags.CloudRegistry, error) {
	request := sdkags.NewDescribeRegistryRequest()
	request.RegistryId = helper.String(id)

	response, err := r.Client().UseAgsClient().DescribeRegistryWithContext(ctx, request)
	if err != nil {
		return nil, err
	}
	if response == nil || response.Response == nil {
		return nil, fmt.Errorf("DescribeRegistry returned an empty response for %s", id)
	}
	if response.Response.Registry == nil {
		return nil, nil
	}
	return response.Response.Registry, nil
}

// expandTags converts the tags list attribute into SDK CloudTag pointers.
func expandTags(ctx context.Context, list types.List) ([]*sdkags.CloudTag, diag.Diagnostics) {
	var diags diag.Diagnostics
	if list.IsNull() || list.IsUnknown() {
		return nil, diags
	}
	elements := list.Elements()
	tags := make([]*sdkags.CloudTag, 0, len(elements))
	for _, element := range elements {
		obj, ok := element.(basetypes.ObjectValue)
		if !ok {
			diags.AddError(
				"Invalid tags element",
				fmt.Sprintf("expected an object element in tags, got %T", element),
			)
			return nil, diags
		}
		var model AgsRegistryTagModel
		objDiags := obj.As(ctx, &model, basetypes.ObjectAsOptions{})
		diags.Append(objDiags...)
		if diags.HasError() {
			return nil, diags
		}
		tags = append(tags, &sdkags.CloudTag{
			Key:   model.Key.ValueStringPointer(),
			Value: model.Value.ValueStringPointer(),
		})
	}
	return tags, diags
}

// refreshRegistryModel copies cloud-owned fields (including the computed
// registry object) into the model. Configurable fields are echoed back by the
// API and are immutable, so overwriting them keeps state aligned with the cloud.
func refreshRegistryModel(m *AgsRegistryResourceModel, item *sdkags.CloudRegistry) {
	if item == nil {
		return
	}

	info := AgsRegistryInfoModel{
		RegistryId:           stringValue(item.RegistryId),
		Name:                 stringValue(item.Name),
		Description:          stringValue(item.Description),
		ApprovalMode:         stringValue(item.ApprovalMode),
		Region:               stringValue(item.Region),
		Status:               stringValue(item.Status),
		CreateTime:           stringValue(item.CreateTime),
		UpdateTime:           stringValue(item.UpdateTime),
		RecordCount:          int64Value(item.RecordCount),
		PublishedRecordCount: int64Value(item.PublishedRecordCount),
		AppId:                int64Value(item.AppId),
		CreatorUin:           stringValue(item.CreatorUin),
		CreatorSubAccountUin: stringValue(item.CreatorSubAccountUin),
		Tags:                 tagsListValue(item.Tags),
	}

	m.Registry = registryObjectValue(info)
	m.Tags = info.Tags
	m.Name = info.Name
	m.Description = info.Description
	m.ApprovalMode = info.ApprovalMode
}

// stringValue converts a nullable SDK string into a types.String.
func stringValue(v *string) types.String {
	if v == nil {
		return types.StringNull()
	}
	return types.StringValue(*v)
}

// int64Value converts a nullable SDK int64 into a types.Int64.
func int64Value(v *int64) types.Int64 {
	if v == nil {
		return types.Int64Null()
	}
	return types.Int64Value(*v)
}

// tagsListValue converts SDK CloudTag pointers into a types.List of objects.
func tagsListValue(tags []*sdkags.CloudTag) types.List {
	if len(tags) == 0 {
		return types.ListNull(tagsObjectType())
	}
	elements := make([]attr.Value, 0, len(tags))
	for _, tag := range tags {
		attributes := map[string]attr.Value{
			"key":   stringValue(tag.Key),
			"value": stringValue(tag.Value),
		}
		obj, diags := types.ObjectValue(tagsAttrTypes(), attributes)
		if diags.HasError() {
			// ObjectValue on a fixed attribute set cannot fail; fall back to a
			// null element to keep the conversion total.
			obj = types.ObjectNull(tagsAttrTypes())
		}
		elements = append(elements, obj)
	}
	list, diags := types.ListValue(tagsObjectType(), elements)
	if diags.HasError() {
		return types.ListNull(tagsObjectType())
	}
	return list
}

// registryObjectValue builds the computed registry object from the info model.
func registryObjectValue(info AgsRegistryInfoModel) types.Object {
	obj, diags := types.ObjectValueFrom(context.Background(), registryAttrTypes(), info)
	if diags.HasError() {
		return types.ObjectNull(registryAttrTypes())
	}
	return obj
}
