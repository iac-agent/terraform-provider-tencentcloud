### Requirement: qualification_status_code computed attribute
The `tencentcloud_sms_sign` resource SHALL expose a computed attribute `qualification_status_code` of type integer that reflects the `QualificationStatusCode` field from the `DescribeSmsSignList` API response (`DescribeSignListStatus.QualificationStatusCode`).

#### Scenario: Read returns qualification status code for domestic SMS sign
- **WHEN** the resource Read method is called for a domestic SMS sign (international=0)
- **THEN** the `qualification_status_code` attribute SHALL be set to the value returned by the API (0=pending, 1=approved, 2=rejected, 3=supplement required, 4=modified pending review, 5=modified rejected)

#### Scenario: Read handles nil QualificationStatusCode
- **WHEN** the API returns a `DescribeSignListStatus` with `QualificationStatusCode` set to nil
- **THEN** the Read method SHALL skip setting the `qualification_status_code` attribute and not return an error

#### Scenario: Read returns qualification status code for international SMS sign
- **WHEN** the resource Read method is called for an international SMS sign (international=1)
- **THEN** the `qualification_status_code` attribute SHALL be set to 0 (the default value for international SMS)
