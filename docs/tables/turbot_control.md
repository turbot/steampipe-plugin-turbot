# Table: turbot_control

Controls in Turbot represent the state of a given check (control type) against
a resource. For example, is encryption at rest enabled for an AWS EBS Volume.

Queries to this table must specify (usually in the `where` clause) at least one
of these columns: `id`, `control_type_id`, `control_type_uri`,
`resource_type_id`, `resource_type_uri`, `state` or `filter`.

## Examples

### Control summary for AWS > IAM > Role > Approved

Simple table:

```sql
select
  state,
  count(*)
from
  turbot_control
where
  control_type_uri = 'tmod:@turbot/aws-iam#/control/types/roleApproved'
group by
  state
order by
  count desc;
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
  control_type_uri = 'tmod:@turbot/aws-iam#/control/types/roleApproved'
group by
  control_type_uri
order by
  total desc;
```

### Control summary for all AWS > IAM controls

```sql
select
  state,
  count(*)
from
  turbot_control
where
  filter = 'controlTypeId:"tmod:@turbot/aws-iam#/resource/types/iam"'
group by
  state
order by
  count desc;
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
  filter = 'controlTypeId:"tmod:@turbot/aws-iam#/resource/types/iam"'
group by
  control_type_uri
order by
  total desc;
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
  timestamp desc;
```

### Query the most recent 10 controls

Note: It's more efficient to have Turbot limit the results to the last 10
(`filter = 'limit:10'`), rather than using `limit 10` which will pull all rows
from Turbot and will then filter them afterwards on the Steampipe side.

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
  timestamp desc;
```

### Control & Resource data for for AWS > IAM > Role > Approved

```sql
select
  r.trunk_title,
  r.data ->> 'Arn' as arn,
  r.metadata -> 'aws' ->> 'accountId' as account_id,
  c.state,
  c.reason
from
  turbot_control as c,
  turbot_resource as r
where
  -- Filter to the control type
  c.control_type_uri = 'tmod:@turbot/aws-iam#/control/types/roleApproved'
  -- Filter to the resource type as well, reducing the size of the join
  and r.resource_type_uri = 'tmod:@turbot/aws-iam#/resource/types/role'
  and r.id = c.resource_id
order by
  r.trunk_title;
```

### Extract all controls from Turbot

WARNING - This is a large query and may take minutes to run. It is not recommended and may timeout.
It's included here as a reference for those who need to extract all data.

```sql
select
  *
from
  turbot_control
where
  filter = '';
```
