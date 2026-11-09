package teo

import (
	"context"
	"errors"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	teo "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

// teoObjectAsOptions is the options type accepted by types.Object.As.
type teoObjectAsOptions = basetypes.ObjectAsOptions

// teoTagModel is the intermediate Go model of a tags element.
type teoTagModel struct {
	TagKey   types.String `tfsdk:"tag_key"`
	TagValue types.String `tfsdk:"tag_value"`
}

// teoVanityNameServersModel is the intermediate Go model of
// vanity_name_servers.
type teoVanityNameServersModel struct {
	Switch  types.String `tfsdk:"switch"`
	Servers types.Set    `tfsdk:"servers"`
}

// Object type definitions used to build computed object attributes.
var (
	teoTagAttrTypes = map[string]attr.Type{
		"tag_key":   types.StringType,
		"tag_value": types.StringType,
	}

	teoDnsVerificationAttrTypes = map[string]attr.Type{
		"subdomain":    types.StringType,
		"record_type":  types.StringType,
		"record_value": types.StringType,
	}

	teoFileVerificationAttrTypes = map[string]attr.Type{
		"path":    types.StringType,
		"content": types.StringType,
	}

	teoNsVerificationAttrTypes = map[string]attr.Type{
		"name_servers": types.SetType{ElemType: types.StringType},
	}

	teoOwnershipVerificationAttrTypes = map[string]attr.Type{
		"dns_verification":  types.ObjectType{AttrTypes: teoDnsVerificationAttrTypes},
		"file_verification": types.ObjectType{AttrTypes: teoFileVerificationAttrTypes},
		"ns_verification":   types.ObjectType{AttrTypes: teoNsVerificationAttrTypes},
	}

	teoVanityNameServersAttrTypes = map[string]attr.Type{
		"switch":  types.StringType,
		"servers": types.SetType{ElemType: types.StringType},
	}
)

// teoNewCreateZoneRequest builds a CreateZone request (the constructor
// initializes the embedded BaseRequest; a struct literal would lose it).
func teoNewCreateZoneRequest() *teo.CreateZoneRequest {
	return teo.NewCreateZoneRequest()
}

// teoNewModifyZoneRequest builds a ModifyZone request.
func teoNewModifyZoneRequest() *teo.ModifyZoneRequest {
	return teo.NewModifyZoneRequest()
}

// teoNewDeleteZoneRequest builds a DeleteZone request.
func teoNewDeleteZoneRequest() *teo.DeleteZoneRequest {
	return teo.NewDeleteZoneRequest()
}

// teoStringPtrOrNull returns a *string for known non-null values, nil otherwise.
func teoStringPtrOrNull(v types.String) *string {
	if v.IsNull() || v.IsUnknown() {
		return nil
	}
	s := v.ValueString()
	return &s
}

// teoBoolPtrOrNull returns a *bool for known non-null values, nil otherwise.
func teoBoolPtrOrNull(v types.Bool) *bool {
	if v.IsNull() || v.IsUnknown() {
		return nil
	}
	b := v.ValueBool()
	return &b
}

// teoTagsFromModels converts tag models into SDK tags.
func teoTagsFromModels(tags []teoTagModel) []*teo.Tag {
	if len(tags) == 0 {
		return nil
	}
	out := make([]*teo.Tag, 0, len(tags))
	for _, t := range tags {
		out = append(out, &teo.Tag{
			TagKey:   teoStringPtrOrNull(t.TagKey),
			TagValue: teoStringPtrOrNull(t.TagValue),
		})
	}
	return out
}

// teoStringSetFromPtrs converts a []*string into a types.Set of strings; a
// nil/empty input yields an empty (not null) set to keep state stable.
func teoStringSetFromPtrs(ctx context.Context, in []*string) (types.Set, diag.Diagnostics) {
	vals := make([]types.String, 0, len(in))
	for _, v := range in {
		if v == nil {
			continue
		}
		vals = append(vals, types.StringValue(*v))
	}
	return types.SetValueFrom(ctx, types.StringType, vals)
}

// teoZone110901RefreshModel copies the cloud-owned fields into the model.
//
// Fields not returned by DescribeZones (allow_duplicates, jump_start) keep
// their existing model value so no diff is produced for them.
func teoZone110901RefreshModel(ctx context.Context, m *TeoZone110901ResourceModel, zone *teo.Zone) {
	if zone == nil {
		return
	}

	m.ZoneId = types.StringPointerValue(zone.ZoneId)
	m.ZoneName = types.StringPointerValue(zone.ZoneName)
	m.Type = types.StringPointerValue(zone.Type)
	m.Area = types.StringPointerValue(zone.Area)
	m.AliasZoneName = types.StringPointerValue(zone.AliasZoneName)
	m.Status = types.StringPointerValue(zone.Status)
	m.CnameStatus = types.StringPointerValue(zone.CnameStatus)
	m.ActiveStatus = types.StringPointerValue(zone.ActiveStatus)
	m.LockStatus = types.StringPointerValue(zone.LockStatus)
	m.Paused = types.BoolPointerValue(zone.Paused)
	m.CreatedOn = types.StringPointerValue(zone.CreatedOn)
	m.ModifiedOn = types.StringPointerValue(zone.ModifiedOn)

	ns, _ := teoStringSetFromPtrs(ctx, zone.NameServers)
	m.NameServers = ns

	ons, _ := teoStringSetFromPtrs(ctx, zone.OriginalNameServers)
	m.OriginalNameServers = ons

	tags, _ := teoZone110901TagsFromSdk(ctx, zone.Tags)
	m.Tags = tags

	m.PlanId = teoZone110901PlanIdFromResources(zone.Resources)

	ov, _ := teoZone110901OwnershipFromSdk(ctx, zone.OwnershipVerification)
	m.OwnershipVerification = ov

	vns, _ := teoZone110901VanityFromSdk(ctx, zone.VanityNameServers)
	m.VanityNameServers = vns
}

// teoZone110901TagsFromSdk converts SDK tags into a types.Set of tag objects.
// A nil/empty input yields an empty (not null) set to keep state stable.
func teoZone110901TagsFromSdk(ctx context.Context, in []*teo.Tag) (types.Set, diag.Diagnostics) {
	items := make([]teoTagModel, 0, len(in))
	for _, t := range in {
		if t == nil {
			continue
		}
		items = append(items, teoTagModel{
			TagKey:   types.StringPointerValue(t.TagKey),
			TagValue: types.StringPointerValue(t.TagValue),
		})
	}
	return types.SetValueFrom(ctx, types.ObjectType{AttrTypes: teoTagAttrTypes}, items)
}

// teoZone110901PlanIdFromResources extracts the plan id bound to the zone from
// the billing resources list (CreateZone.PlanId is read back at
// Zones.Resources[].PlanId).
func teoZone110901PlanIdFromResources(resources []*teo.Resource) types.String {
	for _, res := range resources {
		if res == nil {
			continue
		}
		if res.PlanId != nil && *res.PlanId != "" {
			return types.StringValue(*res.PlanId)
		}
	}
	return types.StringNull()
}

// teoZone110901OwnershipFromSdk maps the ownership verification structure into
// a types.Object matching the ownership_verification attribute.
func teoZone110901OwnershipFromSdk(ctx context.Context, in *teo.OwnershipVerification) (types.Object, diag.Diagnostics) {
	if in == nil {
		return types.ObjectNull(teoOwnershipVerificationAttrTypes), nil
	}

	attrs := map[string]attr.Value{
		"dns_verification":  types.ObjectNull(teoDnsVerificationAttrTypes),
		"file_verification": types.ObjectNull(teoFileVerificationAttrTypes),
		"ns_verification":   types.ObjectNull(teoNsVerificationAttrTypes),
	}

	if in.DnsVerification != nil {
		v, d := types.ObjectValue(teoDnsVerificationAttrTypes, map[string]attr.Value{
			"subdomain":    types.StringPointerValue(in.DnsVerification.Subdomain),
			"record_type":  types.StringPointerValue(in.DnsVerification.RecordType),
			"record_value": types.StringPointerValue(in.DnsVerification.RecordValue),
		})
		if d.HasError() {
			return types.ObjectNull(teoOwnershipVerificationAttrTypes), d
		}
		attrs["dns_verification"] = v
	}
	if in.FileVerification != nil {
		v, d := types.ObjectValue(teoFileVerificationAttrTypes, map[string]attr.Value{
			"path":    types.StringPointerValue(in.FileVerification.Path),
			"content": types.StringPointerValue(in.FileVerification.Content),
		})
		if d.HasError() {
			return types.ObjectNull(teoOwnershipVerificationAttrTypes), d
		}
		attrs["file_verification"] = v
	}
	if in.NsVerification != nil {
		servers, d := teoStringSetFromPtrs(ctx, in.NsVerification.NameServers)
		if d.HasError() {
			return types.ObjectNull(teoOwnershipVerificationAttrTypes), d
		}
		v, d := types.ObjectValue(teoNsVerificationAttrTypes, map[string]attr.Value{
			"name_servers": servers,
		})
		if d.HasError() {
			return types.ObjectNull(teoOwnershipVerificationAttrTypes), d
		}
		attrs["ns_verification"] = v
	}

	return types.ObjectValue(teoOwnershipVerificationAttrTypes, attrs)
}

// teoZone110901VanityFromSdk maps the vanity name servers structure into a
// types.Object matching the vanity_name_servers attribute.
func teoZone110901VanityFromSdk(ctx context.Context, in *teo.VanityNameServers) (types.Object, diag.Diagnostics) {
	if in == nil {
		return types.ObjectNull(teoVanityNameServersAttrTypes), nil
	}
	servers, diags := teoStringSetFromPtrs(ctx, in.Servers)
	if diags.HasError() {
		return types.ObjectNull(teoVanityNameServersAttrTypes), diags
	}
	return types.ObjectValue(teoVanityNameServersAttrTypes, map[string]attr.Value{
		"switch":  types.StringPointerValue(in.Switch),
		"servers": servers,
	})
}

// teoZone110901ResolveUnknown replaces any remaining unknown values with null
// so the state written after Create never contains unknowns. Only fields that
// are not read back (allow_duplicates, jump_start) can still be unknown here.
func teoZone110901ResolveUnknown(m *TeoZone110901ResourceModel) {
	if m.AllowDuplicates.IsUnknown() {
		m.AllowDuplicates = types.BoolNull()
	}
	if m.JumpStart.IsUnknown() {
		m.JumpStart = types.BoolNull()
	}
}

// teoZone110901ResourceId resolves the zone id preferring the dedicated
// zone_id attribute and falling back to id (both hold the zone id).
func teoZone110901ResourceId(m TeoZone110901ResourceModel) string {
	if !m.ZoneId.IsNull() && !m.ZoneId.IsUnknown() && m.ZoneId.ValueString() != "" {
		return m.ZoneId.ValueString()
	}
	if !m.Id.IsNull() && !m.Id.IsUnknown() {
		return m.Id.ValueString()
	}
	return ""
}

// teoDescribeZoneById fetches the zone by id via DescribeZones.
//
// Two different "empty" outcomes stay distinguishable:
//   - API failure -> (nil, err)
//   - query succeeded, no matching zone -> (nil, nil)
//
// The returned zone is additionally matched on ZoneId so an unrelated extra
// row in the response is never mistaken for this resource.
func teoDescribeZoneById(ctx context.Context, client *connectivity.TencentCloudClient, zoneId string) (*teo.Zone, error) {
	if client == nil {
		return nil, errors.New("provider client is not configured")
	}
	request := teo.NewDescribeZonesRequest()
	request.Filters = append(request.Filters, &teo.AdvancedFilter{
		Name:   helper.String("zone-id"),
		Values: []*string{&zoneId},
	})

	var offset int64 = 0
	var limit int64 = 20
	instances := make([]*teo.Zone, 0)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := client.UseTeoClient().DescribeZones(request)
		if err != nil {
			return nil, err
		}
		if response == nil || response.Response == nil {
			return nil, errors.New("DescribeZones returned an empty response")
		}
		if len(response.Response.Zones) < 1 {
			break
		}
		instances = append(instances, response.Response.Zones...)
		if int64(len(response.Response.Zones)) < limit {
			break
		}
		offset += limit
	}

	for _, z := range instances {
		if z != nil && z.ZoneId != nil && *z.ZoneId == zoneId {
			return z, nil
		}
	}
	return nil, nil
}
