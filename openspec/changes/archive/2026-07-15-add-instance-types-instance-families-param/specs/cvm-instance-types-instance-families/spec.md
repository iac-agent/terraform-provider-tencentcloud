## ADDED Requirements

### Requirement: instance_families parameter in cbs_filter
The `tencentcloud_instance_types` data source SHALL support an `instance_families` field (TypeList of TypeString) within the `cbs_filter` nested schema block. This field allows users to specify a list of instance family names to filter disk configuration quotas.

#### Scenario: User specifies instance_families in cbs_filter
- **WHEN** a user provides the `instance_families` parameter in the `cbs_filter` block with values such as `["S5", "M5"]`
- **THEN** the data source SHALL pass these instance families as the `InstanceFamilies` parameter to the `DescribeDiskConfigQuota` API request, instead of using the single `family` derived from instance type results

#### Scenario: User does not specify instance_families
- **WHEN** a user provides the `cbs_filter` block without the `instance_families` parameter
- **THEN** the data source SHALL continue to use the single `family` value from each instance type result as the `InstanceFamilies` parameter to the `DescribeDiskConfigQuota` API request, preserving backward compatibility

#### Scenario: instance_families specified as empty list
- **WHEN** a user provides the `instance_families` parameter as an empty list `[]`
- **THEN** the data source SHALL treat it as "not specified" and fall back to using the single `family` value from the instance type result
