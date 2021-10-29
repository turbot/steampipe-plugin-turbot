# Table: turbot_policy_setting

Policy settings in Turbot are policy definitions assigned to resources and then
applied throughout the hierarchy below (policy values).

It is recommended that queries to this table specify (usually in the `where` clause) at least one
of these columns: `id`, `resource_id`, `exception`, `orphan`, `policy_type_id`,
`policy_type_uri` or `filter`.

## Examples

### Find all policy settings that are exceptions to another policy

```sql
select
  policy_type_uri,
  resource_id,
  is_calculated,
  exception,
  value
from
  turbot_policy_setting
where
  exception;
```

### List policy settings with full resource and policy type information

```sql
select
  r.trunk_title as resource,
  pt.trunk_title as policy_type,
  ps.value,
  ps.is_calculated,
  ps.exception
from
  turbot_policy_setting as ps
  left join turbot_policy_type as pt on pt.id = ps.policy_type_id
  left join turbot_resource as r on r.id = ps.resource_id;
```

### All policy settings set on a given resource

```sql
select
  r.trunk_title as resource,
  ps.resource_id,
  pt.trunk_title as policy_type,
  ps.value,
  ps.is_calculated
from
  turbot_policy_setting as ps
  left join turbot_policy_type as pt on pt.id = ps.policy_type_id
  left join turbot_resource as r on r.id = ps.resource_id
where
  ps.resource_id = 173434983560398;
```

### All policy settings set on a given resource or below

```sql
select
  r.trunk_title as resource,
  ps.resource_id,
  pt.trunk_title as policy_type,
  ps.value,
  ps.is_calculated
from
  turbot_policy_setting as ps
  left join turbot_policy_type as pt on pt.id = ps.policy_type_id
  left join turbot_resource as r on r.id = ps.resource_id
where
  ps.filter = 'resourceId:173434983560398 level:self,descendant';
```

### All policy settings related to AWS > S3 > Bucket

```sql
select
  r.trunk_title as resource,
  ps.resource_id,
  pt.trunk_title as policy_type,
  ps.value,
  ps.is_calculated
from
  turbot_policy_setting as ps
  left join turbot_policy_type as pt on pt.id = ps.policy_type_id
  left join turbot_resource as r on r.id = ps.resource_id
where
  ps.filter = 'resourceTypeId:"tmod:@turbot/aws-s3#/resource/types/bucket"';
```
