## Context

The `tencentcloud_sms_sign` resource manages SMS sign lifecycle (CRUD) in the TencentCloud Terraform Provider. The resource currently reads sign information via the `DescribeSmsSignList` API (`DescribeSignListStatus` struct), but only exposes `sign_name` and `international` as computed attributes in the Read method.

The `DescribeSignListStatus` struct in the SMS SDK (`github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111`) already contains the `QualificationStatusCode` field (`*int64`), which indicates the domestic SMS qualification review status (0=pending, 1=approved, 2=rejected, 3=supplement required, 4=modified pending review, 5=modified rejected). For international SMS, this field defaults to 0.

## Goals / Non-Goals

**Goals:**
- Expose `QualificationStatusCode` as a computed attribute on the `tencentcloud_sms_sign` resource so users can track qualification review status in Terraform state

**Non-Goals:**
- Modifying any existing schema fields or breaking backward compatibility
- Adding write support for `QualificationStatusCode` (it is a read-only status from the API)
- Changing the resource ID format or CRUD logic

## Decisions

1. **Add `qualification_status_code` as a Computed attribute**: Since `QualificationStatusCode` is a status field returned by the API and cannot be set by the user, it is defined as `Computed: true` only. This maintains backward compatibility — existing Terraform configurations are unaffected.

2. **Read the field in `resourceTencentCloudSmsSignRead`**: The `DescribeSmsSign` service method already returns the `DescribeSignListStatus` struct which contains `QualificationStatusCode`. We only need to add a nil check and `d.Set()` call in the Read method, following the existing pattern for `sign_name` and `international`.

3. **Field type: `schema.TypeInt`**: The SDK field is `*int64`, which maps to `schema.TypeInt` in Terraform, consistent with the existing `sign_type`, `document_type`, `international`, and `sign_purpose` fields.

## Risks / Trade-offs

- [Nil field for international SMS] → The API documentation notes that for international SMS, `QualificationStatusCode` defaults to 0. We add a nil check before `d.Set()` to handle any potential nil response, consistent with the existing pattern for other computed fields.
