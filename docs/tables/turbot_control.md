# Table: turbot_control

Search all Turbot controls using a `filter`.

Notes:

- A `filter` must be provided in all queries to this table.
- Use a limit `filter = 'limit:50'` (max 5000) in the filter to limit results. If no limit is provided, then all matching resources will be returned.

## Examples

### Query the most recent 10 controls

```sql
select
  timestamp,
  state,
  reason,
  resource_id,
  control_type_uri
from
  turbot_control
where
  filter = 'limit:10'
order by
  timestamp desc
```

### List controls for AWS > IAM > Role > Approved

```sql
select
  timestamp,
  state,
  reason,
  resource_id,
  control_type_uri
from
  turbot_control
where
  filter = 'controlTypeId:"tmod:@turbot/aws-iam#/control/types/roleApproved" controlTypeLevel:self'
order by
  timestamp desc
```

### Control summary for AWS > IAM > Role > Approved

Simple table:

```sql
select
  state,
  count(*)
from
  turbot_control
where
  filter = 'controlTypeId:"tmod:@turbot/aws-iam#/control/types/roleApproved" controlTypeLevel:self'
group by
  state
order by
  count desc
```

Or, if you prefer a full view of all states:

```sql
select
  control_type_uri,
  sum(case when state = 'ok' then 1 else 0 end) as ok,
  sum(case when state = 'tbd' then 1 else 0 end) as tbd,
  sum(case when state = 'invalid' then 1 else 0 end) as invalid,
  sum(case when state = 'alarm' then 1 else 0 end) as alarm,
  sum(case when state = 'skipped' then 1 else 0 end) as skipped,
  sum(case when state = 'error' then 1 else 0 end) as error,
  sum(case when state in ('alarm', 'error', 'invalid') then 1 else 0 end) as alert,
  count(*) as total
from
  turbot_control as c
where
  filter = 'controlTypeId:"tmod:@turbot/aws-iam#/control/types/roleApproved" controlTypeLevel:self'
```

### Control & Resource data for for AWS > IAM > Role > Approved

```sql
select
  r.title,
  r.data ->> 'Arn' as arn,
  r.metadata -> 'aws' ->> 'accountId' as account_id,
  c.state,
  c.reason
from
  turbot_control as c,
  turbot_resource as r
where
  -- Filter to the control type
  c.filter = 'controlTypeId:"tmod:@turbot/aws-iam#/control/types/roleApproved"'
  -- Filter to the resource type as well, reducing the size of the join
  and r.filter = 'resourceTypeId:"tmod:@turbot/aws-iam#/resource/types/role" resourceTypeLevel:self'
  and r.id = c.resource_id
order by
  arn
```

### Controls with state for AWS > IAM > Role resources

```sql
select
  control_type_uri,
  state,
  count(*)
from
  turbot_control as c
where
  filter = 'resourceTypeId:"tmod:@turbot/aws-iam#/resource/types/role"'
group by
  control_type_uri,
  state
order by
  count desc
```

### Control state by control type for all AWS > IAM resources

```sql
select
  control_type_uri,
  sum(case when state = 'ok' then 1 else 0 end) as ok,
  sum(case when state = 'tbd' then 1 else 0 end) as tbd,
  sum(case when state = 'invalid' then 1 else 0 end) as invalid,
  sum(case when state = 'alarm' then 1 else 0 end) as alarm,
  sum(case when state = 'skipped' then 1 else 0 end) as skipped,
  sum(case when state = 'error' then 1 else 0 end) as error,
  sum(case when state in ('alarm', 'error', 'invalid') then 1 else 0 end) as alert,
  count(*) as total
from
  turbot_control as c
where
  filter = 'resourceTypeId:"tmod:@turbot/aws-iam#/resource/types/iam"'
group by
  control_type_uri
order by
  alert desc
```
