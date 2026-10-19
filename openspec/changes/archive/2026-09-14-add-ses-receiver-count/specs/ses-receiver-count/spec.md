## ADDED Requirements

### Requirement: SES receiver resource exposes count attribute
The `tencentcloud_ses_receiver` resource SHALL expose a computed `count` attribute of type integer that reflects the total number of recipient email addresses in the receiver list, as returned by the `ListReceivers` API (`response.Response.Data.Count`).

#### Scenario: Count is set on resource read
- **WHEN** a `tencentcloud_ses_receiver` resource is read and the `ListReceivers` API returns a non-nil `Count` value
- **THEN** the `count` attribute SHALL be set to the integer value of `Count` from the API response

#### Scenario: Count is nil in API response
- **WHEN** a `tencentcloud_ses_receiver` resource is read and the `ListReceivers` API returns a nil `Count` value
- **THEN** the `count` attribute SHALL not be set (Terraform will use the zero value)

#### Scenario: Count is not user-settable
- **WHEN** a user writes a Terraform configuration for `tencentcloud_ses_receiver`
- **THEN** the `count` attribute SHALL NOT be settable by the user (computed only)

#### Scenario: Count attribute type
- **WHEN** the `count` attribute is defined in the resource schema
- **THEN** the attribute SHALL be of type `schema.TypeInt` with `Computed: true`
