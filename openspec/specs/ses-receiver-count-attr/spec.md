## ADDED Requirements

### Requirement: SES receiver count computed attribute
The `tencentcloud_ses_receiver` resource SHALL expose a computed `count` attribute of type integer that reflects the total number of recipient email addresses in the recipient group, as returned by the `ListReceivers` API response (`ReceiverData.Count`).

#### Scenario: Read populates count from ListReceivers API
- **WHEN** the `resourceTencentCloudSesReceiverRead` function is called for an existing SES receiver
- **THEN** the `count` attribute SHALL be set to the value of `ReceiverData.Count` from the `ListReceivers` API response, provided the value is not nil

#### Scenario: Count is nil in API response
- **WHEN** the `ListReceivers` API response returns a nil `Count` field for the receiver
- **THEN** the `count` attribute SHALL NOT be set (skipped), avoiding a nil pointer dereference error
