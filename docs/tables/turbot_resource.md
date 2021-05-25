# Table: turbot_resource

Search the full Turbot resource collection using a `filter`.

Notes:

- A `filter` must be provided in all queries to this table.
- Use a limit `filter = 'limit:50'` (max 5000) in the filter to limit results. If no limit is provided, then all matching resources will be returned.

## Examples

### Query the most recent 10 resources

```sql
select
  create_timestamp,
  title,
  metadata,
  data
from
  turbot_resource
where
  filter = 'limit:10'
order by
  create_timestamp desc
```

### List all AWS IAM Role resources

```sql
select
  create_timestamp,
  title,
  metadata,
  data
from
  turbot_resource
where
  filter = 'resourceTypeId:"tmod:@turbot/aws-iam#/resource/types/role" resourceTypeLevel:self'
order by
  title
```

### Get the full hierarchy for a resource

TODO - None of these work :-(

```sql
select
  id,
  path,
  create_timestamp,
  title,
  metadata,
  data
from
  turbot_resource
where
  filter = 'resourceTypeId:"tmod:@turbot/aws-iam#/resource/types/role" resourceTypeLevel:self'
order by
  title
```

```sql
select
  p
from
  turbot_resource as r,
  jsonb_array_elements_text(r.path) as p
where
  r.filter = 'resourceId:208149797219788 level:self'
```

```sql
with hierarchy as (
select array(
  select
    p::bigint as id
  from
    turbot_resource as r,
    jsonb_array_elements_text(r.path) as p
  where
    r.filter = 'resourceId:208149797219788 level:self'
) as items)
select
  id,
  path,
  create_timestamp,
  title,
  metadata,
  data
from
  turbot_resource as r
where
  r.id in array[select id from hierarchy]
```

```sql
with hierarchy as (
  select
    p::bigint as id
  from
    turbot_resource as r,
    jsonb_array_elements_text(r.path) as p
  where
    r.filter = 'resourceId:208149797219788 level:self'
)
select
  r.id,
  r.path,
  r.create_timestamp,
  r.title,
  r.metadata,
  r.data
from
  turbot_resource as r,
  hierarchy as h
where
  r.filter = ('level:self resourceId:' || h.id)
```

```sql
select
  create_timestamp,
  title,
  metadata,
  data
from
  turbot_resource as r,
  jsonb_array_elements_text(r.path) as i
where
  r.filter = 'resourceId:"my-bucket"'
  and
  filter = 'resourceTypeId:"tmod:@turbot/aws-iam#/resource/types/role" resourceTypeLevel:self'
order by
  title
```
